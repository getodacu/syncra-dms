package api

import (
	"errors"
	"math"
	"net/http"
	"strings"
	"time"

	"ai.ro/syncra/dms/internal/documents"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	documentUploadMultipartMemoryBytes   int64 = 1 << 20
	documentUploadMultipartOverheadBytes int64 = 1 << 20
)

type documentHandler struct {
	db             *gorm.DB
	auth           *authHandler
	storage        *documents.LocalStorage
	maxUploadBytes int64
}

func newDocumentHandler(options RouterOptions, auth *authHandler) *documentHandler {
	return &documentHandler{
		db:             options.DB,
		auth:           auth,
		storage:        documents.NewLocalStorage(options.DocumentStorageRoot, options.DocumentMaxUploadBytes, options.DocumentAllowedMIMETypes),
		maxUploadBytes: options.DocumentMaxUploadBytes,
	}
}

func (h *documentHandler) upload(c *gin.Context) {
	user, ok := requireAuthenticatedUser(c, h.auth)
	if !ok {
		return
	}
	if h.db == nil {
		writeError(c, http.StatusServiceUnavailable, "database is not configured")
		return
	}
	if h.storage == nil {
		writeError(c, http.StatusServiceUnavailable, "document storage is not configured")
		return
	}

	if requestLimit := documentUploadRequestLimit(h.maxUploadBytes); requestLimit > 0 {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, requestLimit)
	}
	if err := c.Request.ParseMultipartForm(documentUploadMultipartMemoryBytes); err != nil {
		if isRequestBodyTooLarge(err) {
			writeError(c, http.StatusRequestEntityTooLarge, "document exceeds maximum upload size")
			return
		}
		writeError(c, http.StatusBadRequest, "invalid multipart form")
		return
	}
	if c.Request.MultipartForm != nil {
		defer c.Request.MultipartForm.RemoveAll()
	}

	folderID, ok := parseDocumentFolderID(c, c.Request.FormValue("folderId"))
	if !ok {
		return
	}
	uploadedFile, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			writeError(c, http.StatusBadRequest, "file is required")
			return
		}
		writeError(c, http.StatusBadRequest, "invalid file upload")
		return
	}
	defer uploadedFile.Close()

	folderHelper := documentFolderHandler{db: h.db, auth: h.auth}
	folder, ok := folderHelper.loadActiveFolderWithActiveOrganizationUnit(c, folderID)
	if !ok {
		return
	}
	if ok := requireDocumentFolderObjectPermissionForAuthenticatedUser(c, h.auth, user, "document.create", &folder.OrganizationUnitID); !ok {
		return
	}

	stored, err := h.storage.Save(c.Request.Context(), uploadedFile, fileHeader.Filename)
	if err != nil {
		writeDocumentUploadStorageError(c, err)
		return
	}

	var duplicateCount int64
	if err := h.db.WithContext(c.Request.Context()).
		Model(&documents.Document{}).
		Where("folder_id = ? AND sha256_hash = ? AND deleted_at IS NULL", folder.ID, stored.SHA256Hash).
		Count(&duplicateCount).Error; err != nil {
		h.deleteStoredUpload(stored.StorageKey)
		writeError(c, http.StatusInternalServerError, "failed to validate document upload")
		return
	}
	if duplicateCount != 0 {
		h.deleteStoredUpload(stored.StorageKey)
		writeError(c, http.StatusConflict, "active document already exists in folder")
		return
	}

	displayName, err := documents.NormalizeDisplayName(stored.OriginalFileName)
	if err != nil {
		h.deleteStoredUpload(stored.StorageKey)
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	now := time.Now().UTC()
	documentRow := documents.Document{
		FolderID:           folder.ID,
		OrganizationUnitID: folder.OrganizationUnitID,
		OriginalFileName:   stored.OriginalFileName,
		DisplayName:        displayName,
		MimeType:           stored.MimeType,
		Extension:          stored.Extension,
		SizeBytes:          stored.SizeBytes,
		SHA256Hash:         stored.SHA256Hash,
		StorageKey:         stored.StorageKey,
		CreatedByUserID:    user.ID,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	if err := h.db.WithContext(c.Request.Context()).Create(&documentRow).Error; err != nil {
		h.deleteStoredUpload(stored.StorageKey)
		if isDocumentDuplicateHashError(err) {
			writeError(c, http.StatusConflict, "active document already exists in folder")
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to create document")
		return
	}

	c.JSON(http.StatusCreated, documentMetadataResponseFromModel(documentRow))
}

func (h *documentHandler) get(c *gin.Context) {
	writeError(c, http.StatusNotImplemented, "document metadata endpoint is not implemented")
}

func (h *documentHandler) download(c *gin.Context) {
	writeError(c, http.StatusNotImplemented, "document download endpoint is not implemented")
}

func (h *documentHandler) update(c *gin.Context) {
	writeError(c, http.StatusNotImplemented, "document update endpoint is not implemented")
}

func (h *documentHandler) archive(c *gin.Context) {
	writeError(c, http.StatusNotImplemented, "document archive endpoint is not implemented")
}

func (h *documentHandler) deleteStoredUpload(storageKey string) {
	if h.storage == nil || storageKey == "" {
		return
	}
	_ = h.storage.Delete(storageKey)
}

func writeDocumentUploadStorageError(c *gin.Context, err error) {
	if errors.Is(err, documents.ErrUploadTooLarge) {
		writeError(c, http.StatusRequestEntityTooLarge, "document exceeds maximum upload size")
		return
	}
	if errors.Is(err, documents.ErrUnsupportedMIME) {
		writeError(c, http.StatusUnsupportedMediaType, "document MIME type is not allowed")
		return
	}
	writeError(c, http.StatusInternalServerError, "failed to store document")
}

func documentUploadRequestLimit(maxUploadBytes int64) int64 {
	if maxUploadBytes <= 0 {
		return 0
	}
	if maxUploadBytes > math.MaxInt64-documentUploadMultipartOverheadBytes {
		return math.MaxInt64
	}
	return maxUploadBytes + documentUploadMultipartOverheadBytes
}

func isRequestBodyTooLarge(err error) bool {
	var maxBytesError *http.MaxBytesError
	return errors.As(err, &maxBytesError) || strings.Contains(err.Error(), "request body too large")
}

func isDocumentDuplicateHashError(err error) bool {
	if isUniqueConstraintError(err) {
		return true
	}
	message := err.Error()
	return strings.Contains(message, "idx_documents_active_folder_hash_unique") ||
		strings.Contains(message, "documents.folder_id, documents.sha256_hash")
}
