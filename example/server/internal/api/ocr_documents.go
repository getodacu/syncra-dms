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
	"gorm.io/gorm/clause"

	"ai.ro/syncra/internal/ocr"
)

type ocrDocumentListCursor struct {
	CreatedAt time.Time `json:"created_at"`
	ID        uuid.UUID `json:"id"`
	Sort      string    `json:"sort"`
}

const maxOCRDocumentMoveCollectionIDs = 100
const maxOCRDocumentMoveRequestBytes int64 = 1 << 20
const maxOCRDocumentUpdateRequestBytes int64 = 16 << 10

func (h *Handler) ListOCRDocuments(c *gin.Context) {
	userID, err := parseOptionalUserID(c.Query("user_id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	var collectionID uuid.UUID
	rawCollectionID, hasCollectionFilter := c.GetQuery("collection_id")
	if hasCollectionFilter {
		if userID == nil {
			writeError(c, http.StatusBadRequest, "user_id is required")
			return
		}
		collectionID, err = parseCollectionID(rawCollectionID)
		if err != nil {
			writeError(c, http.StatusBadRequest, err.Error())
			return
		}
	}

	var schemaID uuid.UUID
	rawSchemaID, hasSchemaFilter := c.GetQuery("schema_id")
	if hasSchemaFilter {
		schemaID, err = uuid.Parse(strings.TrimSpace(rawSchemaID))
		if err != nil || schemaID == uuid.Nil {
			writeError(c, http.StatusBadRequest, "invalid schema_id")
			return
		}
	}

	createdFrom, err := parseOCRJobTimeQuery(c, "created_from")
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid created_from")
		return
	}
	createdTo, err := parseOCRJobTimeQuery(c, "created_to")
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid created_to")
		return
	}
	if createdFrom != nil && createdTo != nil && createdFrom.After(*createdTo) {
		writeError(c, http.StatusBadRequest, "created_from must be before or equal to created_to")
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
	cursor, err := parseOCRDocumentListCursor(c.Query("cursor"), sortDirection)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	db := h.DB.WithContext(c.Request.Context())
	if hasCollectionFilter {
		var count int64
		if err := db.Model(&ocr.Collection{}).
			Where("id = ? AND user_id = ?", collectionID, *userID).
			Count(&count).Error; err != nil {
			writeError(c, http.StatusInternalServerError, "failed to validate collection")
			return
		}
		if count == 0 {
			writeError(c, http.StatusNotFound, "collection not found")
			return
		}
	}

	query := scopeOCRDocumentsByUserID(db.Model(&ocr.OCRDocument{}), userID)
	if hasCollectionFilter {
		query = query.Joins(
			"JOIN collection_documents ON collection_documents.document_id = ocr_documents.id AND collection_documents.collection_id = ?",
			collectionID,
		)
	}
	if hasSchemaFilter {
		query = query.Where("ocr_documents.schema_id = ?", schemaID)
	}
	filename := strings.TrimSpace(c.Query("filename"))
	if filename != "" {
		query = query.Where("lower(ocr_documents.original_filename) LIKE ? ESCAPE '\\'", "%"+escapeSQLLike(strings.ToLower(filename))+"%")
	}
	if createdFrom != nil {
		query = query.Where("ocr_documents.created_at >= ?", *createdFrom)
	}
	if createdTo != nil {
		query = query.Where("ocr_documents.created_at <= ?", *createdTo)
	}
	if cursor != nil {
		operator := "<"
		if sortDirection == "asc" {
			operator = ">"
		}
		query = query.Where("(ocr_documents.created_at, ocr_documents.id) "+operator+" (?, ?)", cursor.CreatedAt, cursor.ID)
	}

	order := "ocr_documents.created_at desc, ocr_documents.id desc"
	if sortDirection == "asc" {
		order = "ocr_documents.created_at asc, ocr_documents.id asc"
	}
	var documents []ocr.OCRDocument
	if err := query.Order(order).Limit(size + 1).Find(&documents).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to list OCR documents")
		return
	}

	var nextCursor *string
	if len(documents) > size {
		documents = documents[:size]
		if len(documents) > 0 {
			encoded, err := encodeOCRDocumentListCursor(documents[len(documents)-1], sortDirection)
			if err != nil {
				writeError(c, http.StatusInternalServerError, "failed to encode next cursor")
				return
			}
			nextCursor = &encoded
		}
	}

	collectionsByDocument, err := h.documentCollectionSummaries(c.Request.Context(), documents)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to load OCR document collections")
		return
	}

	out := make([]OCRDocumentListItemResponse, 0, len(documents))
	for _, document := range documents {
		out = append(out, ocrDocumentListItemResponse(document, collectionsByDocument[document.ID]))
	}
	loggerFromGin(c).Debug("ocr.documents_listed", "result_count", len(out), "has_next_cursor", nextCursor != nil)
	c.JSON(http.StatusOK, OCRDocumentListResponse{Documents: out, NextCursor: nextCursor})
}

func (h *Handler) documentCollectionSummaries(ctx context.Context, documents []ocr.OCRDocument) (map[uuid.UUID][]OCRDocumentCollectionSummaryResponse, error) {
	summaries := make(map[uuid.UUID][]OCRDocumentCollectionSummaryResponse, len(documents))
	if len(documents) == 0 {
		return summaries, nil
	}

	documentIDs := make([]uuid.UUID, 0, len(documents))
	for _, document := range documents {
		documentIDs = append(documentIDs, document.ID)
	}

	type collectionRow struct {
		DocumentID   uuid.UUID `gorm:"column:document_id"`
		CollectionID uuid.UUID `gorm:"column:collection_id"`
		Name         string    `gorm:"column:name"`
	}
	var rows []collectionRow
	if err := h.DB.WithContext(ctx).
		Table("collection_documents").
		Select("collection_documents.document_id, collections.id AS collection_id, collections.name").
		Joins("JOIN collections ON collections.id = collection_documents.collection_id").
		Where("collection_documents.document_id IN ?", documentIDs).
		Order("collection_documents.document_id asc, collections.name asc, collections.id asc").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		summaries[row.DocumentID] = append(summaries[row.DocumentID], OCRDocumentCollectionSummaryResponse{
			ID:   row.CollectionID,
			Name: row.Name,
		})
	}
	return summaries, nil
}

func scopeOCRDocumentsByUserID(db *gorm.DB, userID *string) *gorm.DB {
	if userID == nil {
		return db.Where("ocr_documents.user_id IS NULL")
	}
	return db.Where("ocr_documents.user_id = ?", *userID)
}

func (h *Handler) UpdateOCRDocument(c *gin.Context) {
	id, err := parseOCRDocumentID(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid OCR document id")
		return
	}

	userID, err := parseOptionalUserID(c.Query("user_id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxOCRDocumentUpdateRequestBytes)
	var req UpdateOCRDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid OCR document update request")
		return
	}

	filename, err := validateOCRDocumentFilename(req.OriginalFilename)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	var doc ocr.OCRDocument
	if err := h.DB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if err := scopeOCRDocumentsByUserID(tx.Where("ocr_documents.id = ?", id), userID).First(&doc).Error; err != nil {
			return err
		}
		if err := tx.Model(&doc).Update("original_filename", filename).Error; err != nil {
			return err
		}
		return tx.First(&doc, "id = ?", id).Error
	}); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusNotFound, "OCR document not found")
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to update OCR document")
		return
	}

	loggerFromGin(c).Info("ocr.document_updated", "document_id", doc.ID.String())
	c.JSON(http.StatusOK, ocrDocumentResponse(doc, false))
}

func validateOCRDocumentFilename(raw string) (string, error) {
	filename := strings.TrimSpace(raw)
	if filename == "" {
		return "", errors.New("original_filename is required")
	}
	if utf8.RuneCountInString(filename) > maxOriginalFilenameCharacters {
		return "", errors.New("filename must be at most 255 characters")
	}
	return filename, nil
}

func (h *Handler) DeleteOCRDocument(c *gin.Context) {
	id, err := parseOCRDocumentID(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid OCR document id")
		return
	}

	var deletedIDs []uuid.UUID
	if rawUserID, ok := c.GetQuery("user_id"); ok {
		userID, err := parseOptionalUserID(rawUserID)
		if err != nil {
			writeError(c, http.StatusBadRequest, err.Error())
			return
		}
		deletedIDs, err = h.deleteOCRDocuments(c.Request.Context(), []uuid.UUID{id}, userID)
	} else {
		deletedIDs, err = h.deleteOCRDocumentsMatching(c.Request.Context(), []uuid.UUID{id}, nil, false)
	}
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to delete OCR document")
		return
	}
	if len(deletedIDs) == 0 {
		writeError(c, http.StatusNotFound, "OCR document not found")
		return
	}

	loggerFromGin(c).Info("ocr.document_deleted", "document_id", id.String())
	c.Status(http.StatusNoContent)
}

func (h *Handler) DeleteOCRDocuments(c *gin.Context) {
	userID, err := parseOptionalUserID(c.Query("user_id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	var req DeleteOCRDocumentsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid OCR document delete request")
		return
	}

	ids, err := parseOCRDocumentIDs(req.IDs)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	deletedIDs, err := h.deleteOCRDocuments(c.Request.Context(), ids, userID)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to delete OCR documents")
		return
	}

	c.JSON(http.StatusOK, DeleteOCRDocumentsResponse{
		DeletedIDs:   deletedIDs,
		DeletedCount: len(deletedIDs),
	})
	loggerFromGin(c).Info("ocr.documents_deleted", "requested_count", len(ids), "deleted_count", len(deletedIDs))
}

func (h *Handler) MoveOCRDocumentsToCollections(c *gin.Context) {
	userID, err := parseRequiredUserID(c.Query("user_id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxOCRDocumentMoveRequestBytes)
	var req MoveOCRDocumentsToCollectionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid OCR document move request")
		return
	}

	ids, err := parseOCRDocumentIDs(req.IDs)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	collectionIDs, err := parseOCRDocumentMoveCollectionIDs(req.CollectionIDs)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	movedIDs, err := h.moveOCRDocumentsToCollections(c.Request.Context(), ids, collectionIDs, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusNotFound, "collection not found")
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to move OCR documents")
		return
	}

	c.JSON(http.StatusOK, MoveOCRDocumentsToCollectionsResponse{
		MovedIDs:      movedIDs,
		MovedCount:    len(movedIDs),
		CollectionIDs: collectionIDs,
	})
	loggerFromGin(c).Info("ocr.documents_moved_to_collections",
		"requested_count", len(ids),
		"moved_count", len(movedIDs),
		"collection_count", len(collectionIDs),
	)
}

func parseOCRDocumentID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(strings.TrimSpace(raw))
	if err != nil || id == uuid.Nil {
		return uuid.Nil, errors.New("invalid OCR document id")
	}
	return id, nil
}

func parseOCRDocumentIDs(rawIDs []string) ([]uuid.UUID, error) {
	if len(rawIDs) == 0 {
		return nil, errors.New("ids is required")
	}

	ids := make([]uuid.UUID, 0, len(rawIDs))
	seen := make(map[uuid.UUID]struct{}, len(rawIDs))
	for _, rawID := range rawIDs {
		id, err := parseOCRDocumentID(rawID)
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

func parseOCRDocumentMoveCollectionIDs(rawIDs []string) ([]uuid.UUID, error) {
	if len(rawIDs) > maxOCRDocumentMoveCollectionIDs {
		return nil, errors.New("invalid collection id")
	}

	ids := make([]uuid.UUID, 0, len(rawIDs))
	seen := make(map[uuid.UUID]struct{}, len(rawIDs))
	for _, rawID := range rawIDs {
		id, err := parseCollectionID(rawID)
		if err != nil {
			return nil, errors.New("invalid collection id")
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	return ids, nil
}

func (h *Handler) deleteOCRDocuments(ctx context.Context, ids []uuid.UUID, userID *string) ([]uuid.UUID, error) {
	return h.deleteOCRDocumentsMatching(ctx, ids, userID, true)
}

func (h *Handler) moveOCRDocumentsToCollections(ctx context.Context, ids []uuid.UUID, collectionIDs []uuid.UUID, userID string) ([]uuid.UUID, error) {
	if len(ids) == 0 {
		return []uuid.UUID{}, nil
	}

	var lockedIDs []uuid.UUID
	if err := h.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if len(collectionIDs) > 0 {
			var count int64
			if err := tx.Model(&ocr.Collection{}).
				Where("id IN ? AND user_id = ?", collectionIDs, userID).
				Count(&count).Error; err != nil {
				return err
			}
			if count != int64(len(collectionIDs)) {
				return gorm.ErrRecordNotFound
			}
		}

		if err := tx.Model(&ocr.OCRDocument{}).
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id IN ? AND user_id = ?", ids, userID).
			Pluck("id", &lockedIDs).Error; err != nil {
			return err
		}
		if len(lockedIDs) == 0 {
			return nil
		}

		if err := tx.Where("document_id IN ?", lockedIDs).Delete(&ocr.CollectionDocument{}).Error; err != nil {
			return err
		}
		if len(collectionIDs) == 0 {
			return nil
		}

		links := make([]ocr.CollectionDocument, 0, len(lockedIDs)*len(collectionIDs))
		for _, documentID := range lockedIDs {
			for _, collectionID := range collectionIDs {
				links = append(links, ocr.CollectionDocument{
					CollectionID: collectionID,
					DocumentID:   documentID,
				})
			}
		}
		return tx.Create(&links).Error
	}); err != nil {
		return nil, err
	}

	matched := make(map[uuid.UUID]struct{}, len(lockedIDs))
	for _, id := range lockedIDs {
		matched[id] = struct{}{}
	}

	movedIDs := make([]uuid.UUID, 0, len(lockedIDs))
	for _, id := range ids {
		if _, ok := matched[id]; ok {
			movedIDs = append(movedIDs, id)
		}
	}
	return movedIDs, nil
}

func (h *Handler) deleteOCRDocumentsMatching(ctx context.Context, ids []uuid.UUID, userID *string, scoped bool) ([]uuid.UUID, error) {
	if len(ids) == 0 {
		return []uuid.UUID{}, nil
	}

	query := h.DB.WithContext(ctx).Model(&ocr.OCRDocument{}).Where("id IN ?", ids)
	if scoped {
		query = scopeByUserID(query, userID)
	}

	var matchedIDs []uuid.UUID
	if err := query.Pluck("id", &matchedIDs).Error; err != nil {
		return nil, err
	}
	if len(matchedIDs) == 0 {
		return []uuid.UUID{}, nil
	}

	var lockedIDs []uuid.UUID
	if err := h.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		lockQuery := tx.Model(&ocr.OCRDocument{}).
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id IN ?", matchedIDs)
		if scoped {
			lockQuery = scopeByUserID(lockQuery, userID)
		}
		if err := lockQuery.Pluck("id", &lockedIDs).Error; err != nil {
			return err
		}
		if len(lockedIDs) == 0 {
			return nil
		}

		if err := tx.Where("document_id IN ?", lockedIDs).Delete(&ocr.CollectionDocument{}).Error; err != nil {
			return err
		}
		deleteQuery := tx.Where("id IN ?", lockedIDs)
		if scoped {
			deleteQuery = scopeByUserID(deleteQuery, userID)
		}
		return deleteQuery.Delete(&ocr.OCRDocument{}).Error
	}); err != nil {
		return nil, err
	}

	matched := make(map[uuid.UUID]struct{}, len(lockedIDs))
	for _, id := range lockedIDs {
		matched[id] = struct{}{}
	}

	deletedIDs := make([]uuid.UUID, 0, len(lockedIDs))
	for _, id := range ids {
		if _, ok := matched[id]; ok {
			deletedIDs = append(deletedIDs, id)
		}
	}

	return deletedIDs, nil
}

func escapeSQLLike(value string) string {
	replacer := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return replacer.Replace(value)
}

func parseOCRDocumentListCursor(raw string, sortDirection string) (*ocrDocumentListCursor, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	decoded, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil {
		return nil, errors.New("invalid cursor")
	}
	var cursor ocrDocumentListCursor
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

func encodeOCRDocumentListCursor(doc ocr.OCRDocument, sortDirection string) (string, error) {
	raw, err := json.Marshal(ocrDocumentListCursor{
		CreatedAt: doc.CreatedAt.UTC(),
		ID:        doc.ID,
		Sort:      sortDirection,
	})
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}
