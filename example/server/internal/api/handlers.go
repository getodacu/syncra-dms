package api

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"ai.ro/syncra/internal/auth"
	ocrsvc "ai.ro/syncra/internal/ocr"
)

// createSchemaRequest contains the fields used to create an extraction schema.
//
// swagger:model createSchemaRequest
type createSchemaRequest struct {
	// Schema name.
	// required: true
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	// Optional owner user id. Omit, empty string, or null stores a system-wide schema.
	// Extensions:
	// ---
	// x-nullable: true
	// ---
	UserID optionalUserID `json:"user_id,omitempty" format:"uuid"`
	// Valid extraction JSON Schema object.
	// required: true
	// type: object
	Schema json.RawMessage `json:"schema" validate:"required" swaggertype:"object"`
	Strict *bool           `json:"strict"`
}

// updateSchemaRequest contains the fields used to update an extraction schema.
//
// swagger:model updateSchemaRequest
type updateSchemaRequest struct {
	// Schema name.
	// required: true
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	// Valid extraction JSON Schema object.
	// required: true
	// type: object
	Schema json.RawMessage `json:"schema" validate:"required" swaggertype:"object"`
	Strict *bool           `json:"strict"`
}

type validatedSchemaFields struct {
	Name        string
	Description string
	Schema      json.RawMessage
	Strict      bool
}

type resolvedSchema struct {
	SchemaID *uuid.UUID
	Schema   json.RawMessage
	Strict   bool
	Inline   bool
}

type schemaListCursor struct {
	CreatedAt time.Time `json:"created_at"`
	ID        uuid.UUID `json:"id"`
	Sort      string    `json:"sort"`
}

// optionalUserID is a nullable UUID string for request bodies.
//
// swagger:strfmt uuid
type optionalUserID string

type uploadData struct {
	Filename string
	MimeType string
	Size     int64
	Bytes    []byte
}

type OCRProcessor = ocrsvc.Processor
type OCRProcessInput = ocrsvc.ProcessInput
type MistralOCRResponse = ocrsvc.MistralResponse
type MistralOCRPage = ocrsvc.MistralPage
type upstreamError = ocrsvc.UpstreamError

func errUpstream(message string) error {
	return ocrsvc.UpstreamError(message)
}

const multipartOverheadAllowanceBytes int64 = 64 << 10
const maxSchemaNameCharacters = 160
const maxOriginalFilenameCharacters = 255
const maxSchemaRequestBytes int64 = 1 << 20
const maxSchemaJSONBytes int64 = 1 << 20
const maxSchemaDeleteIDs = 100

var errInvalidUserID = errors.New("invalid user_id")

func (id *optionalUserID) UnmarshalJSON(raw []byte) error {
	var value *string
	if err := json.Unmarshal(raw, &value); err != nil {
		return errInvalidUserID
	}
	if value == nil {
		*id = ""
		return nil
	}
	*id = optionalUserID(*value)
	return nil
}

func writeError(c *gin.Context, status int, message string) {
	logger := loggerFromGin(c)
	attrs := []any{
		"status", status,
		"message", message,
		"route", routeForLog(c),
	}
	if status >= http.StatusInternalServerError {
		logger.Error("http.error_response", attrs...)
	} else {
		logger.Warn("http.error_response", attrs...)
	}
	c.JSON(status, ErrorResponse{Error: message})
}

func schemaResponse(schema ocrsvc.ExtractionSchema) SchemaResponse {
	return SchemaResponse{
		ID:          schema.ID,
		CreatedAt:   schema.CreatedAt,
		UpdatedAt:   schema.UpdatedAt,
		UserID:      optionalUserIDResponse(schema.UserID),
		Name:        schema.Name,
		Description: schema.Description,
		Schema:      json.RawMessage(schema.SchemaJSON),
		Strict:      schema.Strict,
	}
}

func parseOptionalUserID(raw string) (*string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	id, err := uuid.Parse(raw)
	if err != nil || id == uuid.Nil {
		return nil, errInvalidUserID
	}
	normalized := id.String()
	return &normalized, nil
}

func optionalUserIDResponse(raw *string) *optionalUserID {
	if raw == nil {
		return nil
	}
	value := optionalUserID(*raw)
	return &value
}

func scopeByUserID(db *gorm.DB, userID *string) *gorm.DB {
	if userID == nil {
		return db.Where("user_id IS NULL")
	}
	return db.Where("user_id = ?", *userID)
}

func schemaBindErrorMessage(err error) string {
	if strings.Contains(err.Error(), "request body too large") {
		return "request body too large"
	}
	if strings.Contains(err.Error(), errInvalidUserID.Error()) {
		return errInvalidUserID.Error()
	}
	return "invalid JSON body"
}

func validateSchemaFields(name string, description string, schema json.RawMessage, strict *bool) (validatedSchemaFields, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return validatedSchemaFields{}, errors.New("name is required")
	}
	if utf8.RuneCountInString(name) > maxSchemaNameCharacters {
		return validatedSchemaFields{}, errors.New("name must be at most 160 characters")
	}
	if !isJSONObject(schema) {
		return validatedSchemaFields{}, errors.New("schema must be a JSON object")
	}
	if int64(len(schema)) > maxSchemaJSONBytes {
		return validatedSchemaFields{}, errors.New("schema is too large")
	}
	if err := validateJSONSchema(schema); err != nil {
		return validatedSchemaFields{}, errors.New("schema must be a valid JSON Schema")
	}

	strictValue := true
	if strict != nil {
		strictValue = *strict
	}

	return validatedSchemaFields{
		Name:        name,
		Description: description,
		Schema:      schema,
		Strict:      strictValue,
	}, nil
}

func (h *Handler) CreateSchema(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSchemaRequestBytes)

	var req createSchemaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, schemaBindErrorMessage(err))
		return
	}
	fields, err := validateSchemaFields(req.Name, req.Description, req.Schema, req.Strict)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	userID, err := parseOptionalUserID(string(req.UserID))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	if userID != nil {
		var count int64
		if err := h.DB.Model(&auth.User{}).Where("id = ?", *userID).Count(&count).Error; err != nil {
			writeError(c, http.StatusInternalServerError, "failed to validate user")
			return
		}
		if count == 0 {
			writeError(c, http.StatusBadRequest, "invalid user_id")
			return
		}
	}

	schema := ocrsvc.ExtractionSchema{
		Name:        fields.Name,
		Description: fields.Description,
		UserID:      userID,
		SchemaJSON:  datatypes.JSON(fields.Schema),
		Strict:      fields.Strict,
	}
	if err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&schema).Error; err != nil {
			return err
		}
		if !fields.Strict {
			if err := tx.Model(&schema).UpdateColumn("strict", false).Error; err != nil {
				return err
			}
			schema.Strict = false
		}
		return nil
	}); err != nil {
		writeError(c, http.StatusInternalServerError, "failed to save schema")
		return
	}
	c.JSON(http.StatusCreated, schemaResponse(schema))
}

func (h *Handler) ListSchemas(c *gin.Context) {
	userID, err := parseOptionalUserID(c.Query("user_id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	sortDirection, err := parseOCRJobListSort(c.Query("sort"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	size, err := parseOCRJobListSize(c.Query("size"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	cursor, err := parseSchemaListCursor(c.Query("cursor"), sortDirection)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	query := scopeByUserID(h.DB.WithContext(c.Request.Context()), userID)
	if cursor != nil {
		operator := "<"
		if sortDirection == "asc" {
			operator = ">"
		}
		query = query.Where("(created_at, id) "+operator+" (?, ?)", cursor.CreatedAt, cursor.ID)
	}

	order := "created_at desc, id desc"
	if sortDirection == "asc" {
		order = "created_at asc, id asc"
	}
	var schemas []ocrsvc.ExtractionSchema
	if err := query.Order(order).Limit(size + 1).Find(&schemas).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to list schemas")
		return
	}

	var nextCursor *string
	if len(schemas) > size {
		schemas = schemas[:size]
		if len(schemas) > 0 {
			encoded, err := encodeSchemaListCursor(schemas[len(schemas)-1], sortDirection)
			if err != nil {
				writeError(c, http.StatusInternalServerError, "failed to encode next cursor")
				return
			}
			nextCursor = &encoded
		}
	}

	out := make([]SchemaResponse, 0, len(schemas))
	for _, schema := range schemas {
		out = append(out, schemaResponse(schema))
	}
	c.JSON(http.StatusOK, SchemaListResponse{Schemas: out, NextCursor: nextCursor})
}

func parseSchemaListCursor(raw string, sortDirection string) (*schemaListCursor, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	decoded, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil {
		return nil, errors.New("invalid cursor")
	}
	var cursor schemaListCursor
	if err := json.Unmarshal(decoded, &cursor); err != nil {
		return nil, errors.New("invalid cursor")
	}
	if cursor.ID == uuid.Nil || cursor.CreatedAt.IsZero() || (cursor.Sort != "asc" && cursor.Sort != "desc") {
		return nil, errors.New("invalid cursor")
	}
	if cursor.Sort != sortDirection {
		return nil, errors.New("cursor sort does not match sort")
	}
	return &cursor, nil
}

func encodeSchemaListCursor(schema ocrsvc.ExtractionSchema, sortDirection string) (string, error) {
	raw, err := json.Marshal(schemaListCursor{
		CreatedAt: schema.CreatedAt.UTC(),
		ID:        schema.ID,
		Sort:      sortDirection,
	})
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}

func (h *Handler) GetSchema(c *gin.Context) {
	id, err := parseSchemaID(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid schema id")
		return
	}

	rawUserID, hasUserID := c.GetQuery("user_id")
	userID, err := parseOptionalUserID(rawUserID)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	var schema ocrsvc.ExtractionSchema
	query := h.DB.WithContext(c.Request.Context()).Where("id = ?", id)
	if hasUserID {
		query = scopeByUserID(query, userID)
	}
	if err := query.First(&schema).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusNotFound, "schema not found")
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to load schema")
		return
	}
	c.JSON(http.StatusOK, schemaResponse(schema))
}

func (h *Handler) UpdateSchema(c *gin.Context) {
	id, err := parseSchemaID(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid schema id")
		return
	}

	userID, err := parseOptionalUserID(c.Query("user_id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSchemaRequestBytes)
	var req updateSchemaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, schemaBindErrorMessage(err))
		return
	}
	fields, err := validateSchemaFields(req.Name, req.Description, req.Schema, req.Strict)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	var schema ocrsvc.ExtractionSchema
	query := scopeByUserID(h.DB.WithContext(c.Request.Context()).Where("id = ?", id), userID)
	if err := query.First(&schema).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusNotFound, "schema not found")
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to load schema")
		return
	}

	schema.Name = fields.Name
	schema.Description = fields.Description
	schema.SchemaJSON = datatypes.JSON(fields.Schema)
	schema.Strict = fields.Strict
	if err := h.DB.WithContext(c.Request.Context()).Save(&schema).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to update schema")
		return
	}

	c.JSON(http.StatusOK, schemaResponse(schema))
}

func (h *Handler) DeleteSchema(c *gin.Context) {
	id, err := parseSchemaID(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid schema id")
		return
	}

	userID, err := parseOptionalUserID(c.Query("user_id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	deletedIDs, err := h.deleteSchemas(c.Request.Context(), []uuid.UUID{id}, userID)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to delete schema")
		return
	}
	if len(deletedIDs) == 0 {
		writeError(c, http.StatusNotFound, "schema not found")
		return
	}

	c.JSON(http.StatusOK, DeleteSchemasResponse{
		DeletedIDs:   deletedIDs,
		DeletedCount: len(deletedIDs),
	})
}

func (h *Handler) DeleteSchemas(c *gin.Context) {
	userID, err := parseOptionalUserID(c.Query("user_id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSchemaRequestBytes)
	var req DeleteSchemasRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if strings.Contains(err.Error(), "request body too large") {
			writeError(c, http.StatusBadRequest, "request body too large")
			return
		}
		writeError(c, http.StatusBadRequest, "invalid schema delete request")
		return
	}

	rawIDs := make([]string, 0, len(req.IDs))
	for _, rawID := range req.IDs {
		rawIDs = append(rawIDs, string(rawID))
	}
	ids, err := parseSchemaIDs(rawIDs)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	deletedIDs, err := h.deleteSchemas(c.Request.Context(), ids, userID)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to delete schemas")
		return
	}

	c.JSON(http.StatusOK, DeleteSchemasResponse{
		DeletedIDs:   deletedIDs,
		DeletedCount: len(deletedIDs),
	})
}

func parseSchemaID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(strings.TrimSpace(raw))
	if err != nil || id == uuid.Nil {
		return uuid.Nil, errors.New("invalid schema id")
	}
	return id, nil
}

func parseSchemaIDs(rawIDs []string) ([]uuid.UUID, error) {
	if len(rawIDs) == 0 {
		return nil, errors.New("ids is required")
	}
	if len(rawIDs) > maxSchemaDeleteIDs {
		return nil, errors.New("ids must contain at most 100 values")
	}

	ids := make([]uuid.UUID, 0, len(rawIDs))
	seen := make(map[uuid.UUID]struct{}, len(rawIDs))
	for _, rawID := range rawIDs {
		id, err := parseSchemaID(rawID)
		if err != nil {
			return nil, err
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	return ids, nil
}

func (h *Handler) deleteSchemas(ctx context.Context, ids []uuid.UUID, userID *string) ([]uuid.UUID, error) {
	if len(ids) == 0 {
		return []uuid.UUID{}, nil
	}

	query := scopeByUserID(h.DB.WithContext(ctx).Model(&ocrsvc.ExtractionSchema{}).Where("id IN ?", ids), userID)
	var matchedIDs []uuid.UUID
	if err := query.Pluck("id", &matchedIDs).Error; err != nil {
		return nil, err
	}
	if len(matchedIDs) == 0 {
		return []uuid.UUID{}, nil
	}

	if err := scopeByUserID(h.DB.WithContext(ctx).Where("id IN ?", matchedIDs), userID).
		Delete(&ocrsvc.ExtractionSchema{}).Error; err != nil {
		return nil, err
	}

	matched := make(map[uuid.UUID]struct{}, len(matchedIDs))
	for _, id := range matchedIDs {
		matched[id] = struct{}{}
	}

	deletedIDs := make([]uuid.UUID, 0, len(matchedIDs))
	for _, id := range ids {
		if _, ok := matched[id]; ok {
			deletedIDs = append(deletedIDs, id)
		}
	}
	return deletedIDs, nil
}

func isJSONObject(raw json.RawMessage) bool {
	if len(raw) == 0 {
		return false
	}
	var value map[string]any
	return json.Unmarshal(raw, &value) == nil && value != nil
}

func (h *Handler) maxUploadBytes() int64 {
	if h.MaxUploadBytes > 0 {
		return h.MaxUploadBytes
	}
	return 20 << 20
}

func (h *Handler) readUpload(fileHeader *multipart.FileHeader) (uploadData, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return uploadData{}, err
	}
	defer file.Close()

	limit := h.maxUploadBytes()
	data, err := io.ReadAll(io.LimitReader(file, limit+1))
	if err != nil {
		return uploadData{}, err
	}
	if int64(len(data)) > limit {
		return uploadData{}, errors.New("file exceeds max upload size")
	}

	mimeType, ok := detectSupportedMIME(data)
	if !ok {
		return uploadData{}, errors.New("unsupported file type")
	}

	return uploadData{
		Filename: fileHeader.Filename,
		MimeType: mimeType,
		Size:     int64(len(data)),
		Bytes:    data,
	}, nil
}

func detectSupportedMIME(data []byte) (string, bool) {
	detected := http.DetectContentType(data)
	switch detected {
	case "application/pdf", "image/png", "image/jpeg":
		return detected, true
	}
	if len(data) >= 5 && string(data[:5]) == "%PDF-" {
		return "application/pdf", true
	}
	return "", false
}

func (h *Handler) defaultOCRProcessor() OCRProcessor {
	return ocrsvc.NewMistralProcessor(ocrsvc.MistralConfig{
		APIKey:  h.MistralAPIKey,
		BaseURL: h.MistralBaseURL,
		Model:   h.MistralModel,
	})
}

func computeDocumentHash(file []byte, schema json.RawMessage, strict bool) (string, error) {
	canonicalSchema, hasSchema, err := canonicalSchemaJSON(schema)
	if err != nil {
		return "", err
	}

	h := sha256.New()
	writeHashPart(h, []byte("syncra-ocr-document-hash-v1"))
	writeHashPart(h, file)
	if hasSchema {
		writeHashPart(h, []byte("schema"))
		writeHashPart(h, canonicalSchema)
		if strict {
			writeHashPart(h, []byte("strict:true"))
		} else {
			writeHashPart(h, []byte("strict:false"))
		}
	} else {
		writeHashPart(h, []byte("schema:none"))
		writeHashPart(h, []byte("strict:false"))
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func writeHashPart(w io.Writer, data []byte) {
	var length [8]byte
	binary.BigEndian.PutUint64(length[:], uint64(len(data)))
	_, _ = w.Write(length[:])
	_, _ = w.Write(data)
}

func canonicalSchemaJSON(raw json.RawMessage) ([]byte, bool, error) {
	if len(raw) == 0 {
		return nil, false, nil
	}
	var value any
	if err := json.Unmarshal(raw, &value); err != nil {
		return nil, false, err
	}
	out, err := json.Marshal(value)
	if err != nil {
		return nil, false, err
	}
	return out, true, nil
}

func (h *Handler) resolveSchema(c *gin.Context) (resolvedSchema, bool) {
	return h.resolveSchemaScoped(c, nil)
}

func (h *Handler) resolveSchemaScoped(c *gin.Context, scopeUserID *string) (resolvedSchema, bool) {
	inline := strings.TrimSpace(c.PostForm("schema"))
	schemaIDRaw := strings.TrimSpace(c.PostForm("schema_id"))
	if inline != "" && schemaIDRaw != "" {
		writeError(c, http.StatusBadRequest, "provide either schema or schema_id, not both")
		return resolvedSchema{}, false
	}
	if inline != "" {
		raw := json.RawMessage(inline)
		if !isJSONObject(raw) {
			writeError(c, http.StatusBadRequest, "schema must be a JSON object")
			return resolvedSchema{}, false
		}
		return resolvedSchema{Schema: raw, Strict: true, Inline: true}, true
	}
	if schemaIDRaw == "" {
		return resolvedSchema{}, true
	}
	id, err := uuid.Parse(schemaIDRaw)
	if err != nil || id == uuid.Nil {
		writeError(c, http.StatusBadRequest, "invalid schema_id")
		return resolvedSchema{}, false
	}
	var schema ocrsvc.ExtractionSchema
	query := h.DB.WithContext(c.Request.Context()).Where("id = ?", id)
	if scopeUserID != nil {
		query = query.Where("user_id = ?", *scopeUserID)
	}
	if err := query.First(&schema).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusNotFound, "schema not found")
			return resolvedSchema{}, false
		}
		writeError(c, http.StatusInternalServerError, "failed to load schema")
		return resolvedSchema{}, false
	}
	return resolvedSchema{SchemaID: &schema.ID, Schema: json.RawMessage(schema.SchemaJSON), Strict: schema.Strict}, true
}

func (h *Handler) CreateOCRDocument(c *gin.Context) {
	logger := loggerFromGin(c).With("domain", "ocr")
	limit := h.maxUploadBytes()
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limit+multipartOverheadAllowanceBytes)

	fileHeader, err := c.FormFile("file")
	if err != nil {
		if strings.Contains(err.Error(), "request body too large") {
			writeError(c, http.StatusBadRequest, "request body too large")
			return
		}
		writeError(c, http.StatusBadRequest, "file is required")
		return
	}
	upload, err := h.readUpload(fileHeader)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	logger.Debug("ocr.document_upload_accepted",
		"mime_type", upload.MimeType,
		"file_size", upload.Size,
		"max_upload_bytes", limit,
	)
	if utf8.RuneCountInString(upload.Filename) > maxOriginalFilenameCharacters {
		writeError(c, http.StatusBadRequest, "filename must be at most 255 characters")
		return
	}
	userID, err := parseOptionalUserID(c.PostForm("user_id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	if userID != nil {
		logger = logger.With("user_id", *userID)
		var count int64
		if err := h.DB.Model(&auth.User{}).Where("id = ?", *userID).Count(&count).Error; err != nil {
			writeError(c, http.StatusInternalServerError, "failed to validate user")
			return
		}
		if count == 0 {
			writeError(c, http.StatusBadRequest, "invalid user_id")
			return
		}
	}
	schema, ok := h.resolveSchema(c)
	if !ok {
		return
	}
	documentHash, err := computeDocumentHash(upload.Bytes, schema.Schema, schema.Strict)
	if err != nil {
		logger.Warn("ocr.document_hash_failed", "error", safeLogError(err))
		writeError(c, http.StatusBadRequest, "invalid schema")
		return
	}
	logger.Debug("ocr.document_hash_computed", resolvedSchemaLogAttrs(schema)...)

	var cachedDoc ocrsvc.OCRDocument
	cacheQuery := h.DB.Where("document_hash = ? AND created_at >= ?", documentHash, time.Now().Add(-24*time.Hour))
	if err := scopeByUserID(cacheQuery, userID).
		Order("created_at desc").
		First(&cachedDoc).Error; err == nil {
		if schema.SchemaID != nil && cachedDoc.SchemaID != nil {
			if err := ocrsvc.LinkDocumentToMatchingCollectionsForSchema(c.Request.Context(), h.DB, cachedDoc.ID, cachedDoc.UserID, schema.SchemaID); err != nil {
				writeError(c, http.StatusInternalServerError, "failed to link OCR document collections")
				return
			}
		}
		logger.Info("ocr.document_cache_hit",
			"document_id", cachedDoc.ID.String(),
			"mime_type", cachedDoc.MimeType,
			"file_size", cachedDoc.FileSize,
			"page_count", cachedDoc.PageCount,
		)
		c.JSON(http.StatusOK, ocrDocumentResponse(cachedDoc, true))
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		writeError(c, http.StatusInternalServerError, "failed to load cached OCR document")
		return
	}

	processor := h.OCR
	if processor == nil {
		processor = h.defaultOCRProcessor()
	}
	processorStart := time.Now()
	ocrResponse, rawResponse, err := processor(c.Request.Context(), OCRProcessInput{
		Filename: upload.Filename,
		MimeType: upload.MimeType,
		DataURL:  ocrsvc.DataURL(upload.MimeType, upload.Bytes),
		Schema:   schema.Schema,
		Strict:   schema.Strict,
	})
	if err != nil {
		logger.Error("ocr.upstream_request_failed", "duration_ms", time.Since(processorStart).Milliseconds(), "error", safeLogError(err))
		var upstream ocrsvc.UpstreamError
		if errors.As(err, &upstream) {
			writeError(c, http.StatusBadGateway, err.Error())
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to process OCR")
		return
	}
	logger.Debug("ocr.upstream_request_completed", "duration_ms", time.Since(processorStart).Milliseconds(), "response_bytes", len(rawResponse))

	pageCount, err := ocrsvc.CountRawResponsePages(rawResponse)
	if err != nil {
		writeError(c, http.StatusBadGateway, "invalid Mistral OCR response")
		return
	}
	logger.Debug("ocr.raw_response_pages_counted", "page_count", pageCount)

	annotation, err := ocrsvc.ParseAnnotationJSON(ocrResponse.DocumentAnnotation, len(schema.Schema) > 0)
	if err != nil {
		if err.Error() == "missing Mistral document annotation JSON" {
			writeError(c, http.StatusBadGateway, err.Error())
			return
		}
		writeError(c, http.StatusBadGateway, "invalid Mistral document annotation JSON")
		return
	}
	logger.Debug("ocr.annotation_parsed", "has_annotation", len(annotation) > 0)

	doc := ocrsvc.OCRDocument{
		UserID:           userID,
		OriginalFilename: upload.Filename,
		MimeType:         upload.MimeType,
		FileSize:         upload.Size,
		PageCount:        pageCount,
		DocumentHash:     documentHash,
		SchemaID:         schema.SchemaID,
		Markdown:         ocrsvc.JoinMarkdown(ocrResponse.Pages),
		AnnotationJSON:   datatypes.JSON(annotation),
		RawResponseJSON:  datatypes.JSON(rawResponse),
	}
	if schema.Inline {
		doc.InlineSchemaJSON = datatypes.JSON(schema.Schema)
	}
	if err := h.DB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&doc).Error; err != nil {
			return err
		}
		return ocrsvc.LinkDocumentToMatchingCollections(c.Request.Context(), tx, doc)
	}); err != nil {
		writeError(c, http.StatusInternalServerError, "failed to save OCR document")
		return
	}
	logger.Info("ocr.document_created",
		"document_id", doc.ID.String(),
		"mime_type", doc.MimeType,
		"file_size", doc.FileSize,
		"page_count", doc.PageCount,
	)
	c.JSON(http.StatusCreated, ocrDocumentResponse(doc, false))
}

func (h *Handler) GetOCRDocument(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil || id == uuid.Nil {
		writeError(c, http.StatusBadRequest, "invalid OCR document id")
		return
	}

	rawUserID, hasUserID := c.GetQuery("user_id")
	userID, err := parseOptionalUserID(rawUserID)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	var doc ocrsvc.OCRDocument
	query := h.DB.Where("id = ?", id)
	if hasUserID {
		query = scopeByUserID(query, userID)
	}
	if err := query.First(&doc).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusNotFound, "OCR document not found")
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to load OCR document")
		return
	}
	c.JSON(http.StatusOK, ocrDocumentResponse(doc, false))
}
