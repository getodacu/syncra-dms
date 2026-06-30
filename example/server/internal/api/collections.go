package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"ai.ro/syncra/internal/auth"
	ocrsvc "ai.ro/syncra/internal/ocr"
)

type createCollectionRequest struct {
	Name      string   `json:"name"`
	UserID    string   `json:"user_id"`
	SchemaIDs []string `json:"schema_ids"`
}

type updateCollectionRequest struct {
	Name      string   `json:"name"`
	SchemaIDs []string `json:"schema_ids"`
}

type collectionListCursor struct {
	CreatedAt time.Time `json:"created_at"`
	ID        uuid.UUID `json:"id"`
	Sort      string    `json:"sort"`
}

const maxCollectionNameCharacters = 160
const maxCollectionSchemaIDs = 100
const maxCollectionRequestBytes int64 = 1 << 20

var errInvalidCollectionSchemaIDs = errors.New("invalid schema_ids")

func parseRequiredUserID(raw string) (string, error) {
	userID, err := parseOptionalUserID(raw)
	if err != nil {
		return "", err
	}
	if userID == nil {
		return "", errors.New("user_id is required")
	}
	return *userID, nil
}

func validateCollectionName(raw string) (string, error) {
	name := strings.TrimSpace(raw)
	if name == "" {
		return "", errors.New("name is required")
	}
	if utf8.RuneCountInString(name) > maxCollectionNameCharacters {
		return "", errors.New("name must be at most 160 characters")
	}
	return name, nil
}

func parseCollectionID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(strings.TrimSpace(raw))
	if err != nil || id == uuid.Nil {
		return uuid.Nil, errors.New("invalid collection id")
	}
	return id, nil
}

func parseCollectionSchemaIDs(rawIDs []string) ([]uuid.UUID, error) {
	if len(rawIDs) > maxCollectionSchemaIDs {
		return nil, errInvalidCollectionSchemaIDs
	}

	ids := make([]uuid.UUID, 0, len(rawIDs))
	seen := make(map[uuid.UUID]struct{}, len(rawIDs))
	for _, rawID := range rawIDs {
		id, err := uuid.Parse(strings.TrimSpace(rawID))
		if err != nil || id == uuid.Nil {
			return nil, errInvalidCollectionSchemaIDs
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	return ids, nil
}

func validateCollectionUserExists(db *gorm.DB, userID string) error {
	var count int64
	if err := db.Model(&auth.User{}).Where("id = ?", userID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return errInvalidUserID
	}
	return nil
}

func validateCollectionSchemasBelongToUser(db *gorm.DB, userID string, schemaIDs []uuid.UUID) error {
	if len(schemaIDs) == 0 {
		return nil
	}

	var count int64
	if err := db.Model(&ocrsvc.ExtractionSchema{}).
		Where("id IN ? AND user_id = ?", schemaIDs, userID).
		Count(&count).Error; err != nil {
		return err
	}
	if count != int64(len(schemaIDs)) {
		return errInvalidCollectionSchemaIDs
	}
	return nil
}

func collectionBindErrorMessage(err error) string {
	if strings.Contains(err.Error(), "request body too large") {
		return "request body too large"
	}
	return "invalid JSON body"
}

func (h *Handler) CreateCollection(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxCollectionRequestBytes)

	var req createCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, collectionBindErrorMessage(err))
		return
	}
	userID, err := parseRequiredUserID(req.UserID)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	name, err := validateCollectionName(req.Name)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	db := h.DB.WithContext(c.Request.Context())
	if err := validateCollectionUserExists(db, userID); err != nil {
		if errors.Is(err, errInvalidUserID) {
			writeError(c, http.StatusBadRequest, err.Error())
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to validate user")
		return
	}
	schemaIDs, err := parseCollectionSchemaIDs(req.SchemaIDs)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := validateCollectionSchemasBelongToUser(db, userID, schemaIDs); err != nil {
		if errors.Is(err, errInvalidCollectionSchemaIDs) {
			writeError(c, http.StatusBadRequest, err.Error())
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to validate schemas")
		return
	}

	collection := ocrsvc.Collection{
		UserID: userID,
		Name:   name,
	}
	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&collection).Error; err != nil {
			return err
		}
		return h.replaceCollectionSchemas(tx, collection.ID, schemaIDs)
	}); err != nil {
		writeError(c, http.StatusInternalServerError, "failed to save collection")
		return
	}

	resp, err := h.collectionResponse(c.Request.Context(), collection)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to load collection")
		return
	}
	resp.SchemaIDs = make([]uuid.UUID, len(schemaIDs))
	copy(resp.SchemaIDs, schemaIDs)
	resp.SchemaCount = len(schemaIDs)
	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) ListCollections(c *gin.Context) {
	userID, err := parseRequiredUserID(c.Query("user_id"))
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
	cursor, err := parseCollectionListCursor(c.Query("cursor"), sortDirection)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	query := h.DB.WithContext(c.Request.Context()).Where("user_id = ?", userID)
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
	var collections []ocrsvc.Collection
	if err := query.Order(order).Limit(size + 1).Find(&collections).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to list collections")
		return
	}

	var nextCursor *string
	if len(collections) > size {
		collections = collections[:size]
		if len(collections) > 0 {
			encoded, err := encodeCollectionListCursor(collections[len(collections)-1], sortDirection)
			if err != nil {
				writeError(c, http.StatusInternalServerError, "failed to encode next cursor")
				return
			}
			nextCursor = &encoded
		}
	}

	out, err := h.collectionResponses(c.Request.Context(), collections)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to load collections")
		return
	}
	c.JSON(http.StatusOK, CollectionListResponse{Collections: out, NextCursor: nextCursor})
}

func (h *Handler) GetCollection(c *gin.Context) {
	id, err := parseCollectionID(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	userID, err := parseRequiredUserID(c.Query("user_id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	var collection ocrsvc.Collection
	if err := h.DB.WithContext(c.Request.Context()).
		Where("id = ? AND user_id = ?", id, userID).
		First(&collection).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusNotFound, "collection not found")
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to load collection")
		return
	}

	resp, err := h.collectionResponse(c.Request.Context(), collection)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to load collection")
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateCollection(c *gin.Context) {
	id, err := parseCollectionID(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	userID, err := parseRequiredUserID(c.Query("user_id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxCollectionRequestBytes)
	var req updateCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, collectionBindErrorMessage(err))
		return
	}
	name, err := validateCollectionName(req.Name)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	schemaIDs, err := parseCollectionSchemaIDs(req.SchemaIDs)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	var collection ocrsvc.Collection
	db := h.DB.WithContext(c.Request.Context())
	if err := db.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&ocrsvc.Collection{}).
			Where("id = ? AND user_id = ?", id, userID).
			Update("name", name)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		if err := validateCollectionSchemasBelongToUser(tx, userID, schemaIDs); err != nil {
			return err
		}
		if err := h.replaceCollectionSchemas(tx, id, schemaIDs); err != nil {
			return err
		}
		return tx.Where("id = ? AND user_id = ?", id, userID).First(&collection).Error
	}); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusNotFound, "collection not found")
			return
		}
		if errors.Is(err, errInvalidCollectionSchemaIDs) {
			writeError(c, http.StatusBadRequest, err.Error())
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to update collection")
		return
	}

	resp, err := h.collectionResponse(c.Request.Context(), collection)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to load collection")
		return
	}
	resp.SchemaIDs = make([]uuid.UUID, len(schemaIDs))
	copy(resp.SchemaIDs, schemaIDs)
	resp.SchemaCount = len(schemaIDs)
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) DeleteCollection(c *gin.Context) {
	id, err := parseCollectionID(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	userID, err := parseRequiredUserID(c.Query("user_id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.DB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var collection ocrsvc.Collection
		if err := tx.Where("id = ? AND user_id = ?", id, userID).First(&collection).Error; err != nil {
			return err
		}
		if err := tx.Where("collection_id = ?", collection.ID).Delete(&ocrsvc.CollectionDocument{}).Error; err != nil {
			return err
		}
		if err := tx.Where("collection_id = ?", collection.ID).Delete(&ocrsvc.CollectionSchema{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ? AND user_id = ?", collection.ID, userID).Delete(&ocrsvc.Collection{}).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusNotFound, "collection not found")
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to delete collection")
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) collectionResponse(ctx context.Context, collection ocrsvc.Collection) (CollectionResponse, error) {
	schemaIDs := make([]uuid.UUID, 0)
	if err := h.DB.WithContext(ctx).
		Model(&ocrsvc.CollectionSchema{}).
		Where("collection_id = ?", collection.ID).
		Order("schema_id asc").
		Pluck("schema_id", &schemaIDs).Error; err != nil {
		return CollectionResponse{}, err
	}

	var documentCount int64
	if err := h.DB.WithContext(ctx).
		Model(&ocrsvc.CollectionDocument{}).
		Where("collection_id = ?", collection.ID).
		Count(&documentCount).Error; err != nil {
		return CollectionResponse{}, err
	}

	return makeCollectionResponse(collection, schemaIDs, documentCount), nil
}

func (h *Handler) collectionResponses(ctx context.Context, collections []ocrsvc.Collection) ([]CollectionResponse, error) {
	if len(collections) == 0 {
		return []CollectionResponse{}, nil
	}

	collectionIDs := make([]uuid.UUID, 0, len(collections))
	for _, collection := range collections {
		collectionIDs = append(collectionIDs, collection.ID)
	}

	type collectionSchemaRow struct {
		CollectionID uuid.UUID `gorm:"column:collection_id"`
		SchemaID     uuid.UUID `gorm:"column:schema_id"`
	}
	var schemaRows []collectionSchemaRow
	if err := h.DB.WithContext(ctx).
		Model(&ocrsvc.CollectionSchema{}).
		Select("collection_id, schema_id").
		Where("collection_id IN ?", collectionIDs).
		Order("collection_id asc, schema_id asc").
		Scan(&schemaRows).Error; err != nil {
		return nil, err
	}
	schemaIDsByCollection := make(map[uuid.UUID][]uuid.UUID, len(collections))
	for _, row := range schemaRows {
		schemaIDsByCollection[row.CollectionID] = append(schemaIDsByCollection[row.CollectionID], row.SchemaID)
	}

	type collectionDocumentCountRow struct {
		CollectionID uuid.UUID `gorm:"column:collection_id"`
		Count        int64     `gorm:"column:count"`
	}
	var documentCountRows []collectionDocumentCountRow
	if err := h.DB.WithContext(ctx).
		Model(&ocrsvc.CollectionDocument{}).
		Select("collection_id, count(*) AS count").
		Where("collection_id IN ?", collectionIDs).
		Group("collection_id").
		Scan(&documentCountRows).Error; err != nil {
		return nil, err
	}
	documentCountsByCollection := make(map[uuid.UUID]int64, len(documentCountRows))
	for _, row := range documentCountRows {
		documentCountsByCollection[row.CollectionID] = row.Count
	}

	responses := make([]CollectionResponse, 0, len(collections))
	for _, collection := range collections {
		responses = append(responses, makeCollectionResponse(
			collection,
			schemaIDsByCollection[collection.ID],
			documentCountsByCollection[collection.ID],
		))
	}
	return responses, nil
}

func makeCollectionResponse(collection ocrsvc.Collection, schemaIDs []uuid.UUID, documentCount int64) CollectionResponse {
	if schemaIDs == nil {
		schemaIDs = []uuid.UUID{}
	}
	return CollectionResponse{
		ID:            collection.ID,
		CreatedAt:     collection.CreatedAt,
		UpdatedAt:     collection.UpdatedAt,
		UserID:        collectionUserIDString(collection.UserID),
		Name:          collection.Name,
		SchemaIDs:     schemaIDs,
		SchemaCount:   len(schemaIDs),
		DocumentCount: documentCount,
	}
}

func (h *Handler) replaceCollectionSchemas(tx *gorm.DB, collectionID uuid.UUID, schemaIDs []uuid.UUID) error {
	if err := tx.Where("collection_id = ?", collectionID).Delete(&ocrsvc.CollectionSchema{}).Error; err != nil {
		return err
	}

	links := make([]ocrsvc.CollectionSchema, 0, len(schemaIDs))
	seen := make(map[uuid.UUID]struct{}, len(schemaIDs))
	for _, schemaID := range schemaIDs {
		if _, ok := seen[schemaID]; ok {
			continue
		}
		seen[schemaID] = struct{}{}
		links = append(links, ocrsvc.CollectionSchema{
			CollectionID: collectionID,
			SchemaID:     schemaID,
		})
	}
	if len(links) == 0 {
		return nil
	}
	return tx.Create(&links).Error
}

func parseCollectionListCursor(raw string, sortDirection string) (*collectionListCursor, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	decoded, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil {
		return nil, errors.New("invalid cursor")
	}
	var cursor collectionListCursor
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

func encodeCollectionListCursor(collection ocrsvc.Collection, sortDirection string) (string, error) {
	raw, err := json.Marshal(collectionListCursor{
		CreatedAt: collection.CreatedAt.UTC(),
		ID:        collection.ID,
		Sort:      sortDirection,
	})
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}
