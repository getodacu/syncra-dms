package api

import (
	"errors"
	"io"
	"math"
	"mime"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"ai.ro/syncra/dms/internal/auth"
	"ai.ro/syncra/dms/internal/documents"
	"ai.ro/syncra/dms/internal/orgunits"
	"ai.ro/syncra/dms/internal/rbac"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	documentUploadMultipartMemoryBytes   int64 = 1 << 20
	documentUploadMultipartOverheadBytes int64 = 1 << 20
	maxDocumentMetadataRequestBytes      int64 = 1 << 20
)

var (
	errDocumentUploadDuplicate          = errors.New("document upload duplicate")
	errDocumentUploadResponseWritten    = errors.New("document upload response already written")
	errDocumentUploadScopeNotFound      = errors.New("document upload scope not found")
	errDocumentOperationResponseWritten = errors.New("document operation response already written")
)

type documentHandler struct {
	db             *gorm.DB
	auth           *authHandler
	storage        *documents.LocalStorage
	maxUploadBytes int64
}

type documentUpdateRequest struct {
	DisplayName string `json:"displayName"`
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

	displayName, err := documents.NormalizeDisplayName(stored.OriginalFileName)
	if err != nil {
		if !h.cleanupStoredUploadOrWriteError(c, stored.StorageKey) {
			return
		}
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	var documentRow documents.Document
	err = h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if ok := folderHelper.lockFolderHierarchyWithDB(c, tx, folder.OrganizationUnitID); !ok {
			return errDocumentUploadResponseWritten
		}

		lockedFolder, err := h.revalidateActiveUploadScopeWithDB(c, tx, folder.ID, folder.OrganizationUnitID)
		if err != nil {
			return err
		}

		var duplicateCount int64
		if err := tx.WithContext(c.Request.Context()).
			Model(&documents.Document{}).
			Where("folder_id = ? AND sha256_hash = ? AND deleted_at IS NULL", lockedFolder.ID, stored.SHA256Hash).
			Count(&duplicateCount).Error; err != nil {
			return err
		}
		if duplicateCount != 0 {
			return errDocumentUploadDuplicate
		}

		now := time.Now().UTC()
		documentRow = documents.Document{
			FolderID:           lockedFolder.ID,
			OrganizationUnitID: lockedFolder.OrganizationUnitID,
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
		return tx.WithContext(c.Request.Context()).Create(&documentRow).Error
	})
	if err != nil {
		if !h.cleanupStoredUploadOrWriteError(c, stored.StorageKey) {
			return
		}
		switch {
		case errors.Is(err, errDocumentUploadResponseWritten):
			return
		case errors.Is(err, errDocumentUploadScopeNotFound):
			writeError(c, http.StatusNotFound, "document folder not found")
			return
		case errors.Is(err, errDocumentUploadDuplicate), isDocumentDuplicateHashError(err):
			writeError(c, http.StatusConflict, "active document already exists in folder")
			return
		default:
			writeError(c, http.StatusInternalServerError, "failed to create document")
			return
		}
	}

	c.JSON(http.StatusCreated, documentMetadataResponseFromModel(documentRow))
}

func (h *documentHandler) get(c *gin.Context) {
	user, ok := requireAuthenticatedUser(c, h.auth)
	if !ok {
		return
	}
	if h.db == nil {
		writeError(c, http.StatusServiceUnavailable, "database is not configured")
		return
	}
	id, ok := parseDocumentID(c, c.Param("id"))
	if !ok {
		return
	}
	documentRow, ok := h.loadActiveDocumentWithActiveScope(c, id)
	if !ok {
		return
	}
	if ok := requireDocumentObjectPermissionForAuthenticatedUser(c, h.auth, user, "document.view", documentRow.OrganizationUnitID); !ok {
		return
	}

	c.JSON(http.StatusOK, documentMetadataResponseFromModel(documentRow))
}

func (h *documentHandler) download(c *gin.Context) {
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
	id, ok := parseDocumentID(c, c.Param("id"))
	if !ok {
		return
	}
	documentRow, ok := h.loadActiveDocumentWithActiveScope(c, id)
	if !ok {
		return
	}
	if ok := requireDocumentObjectPermissionForAuthenticatedUser(c, h.auth, user, "document.download", documentRow.OrganizationUnitID); !ok {
		return
	}

	reader, err := h.storage.Open(documentRow.StorageKey)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			writeError(c, http.StatusNotFound, "document file not found")
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to open document file")
		return
	}
	defer reader.Close()

	info, err := reader.Stat()
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to open document file")
		return
	}
	if info.Size() != documentRow.SizeBytes {
		writeError(c, http.StatusInternalServerError, "document file storage is invalid")
		return
	}

	c.Header("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": documentRow.OriginalFileName}))
	c.Header("Content-Type", documentRow.MimeType)
	c.Header("Content-Length", strconv.FormatInt(documentRow.SizeBytes, 10))
	c.Status(http.StatusOK)
	if _, err := io.Copy(c.Writer, reader); err != nil {
		_ = c.Error(err)
	}
}

func (h *documentHandler) update(c *gin.Context) {
	user, ok := requireAuthenticatedUser(c, h.auth)
	if !ok {
		return
	}
	if h.db == nil {
		writeError(c, http.StatusServiceUnavailable, "database is not configured")
		return
	}
	id, ok := parseDocumentID(c, c.Param("id"))
	if !ok {
		return
	}
	documentRow, ok := h.loadActiveDocumentWithActiveScope(c, id)
	if !ok {
		return
	}
	if ok := requireDocumentObjectPermissionForAuthenticatedUser(c, h.auth, user, "document.update", documentRow.OrganizationUnitID); !ok {
		return
	}
	req, ok := bindDocumentUpdateRequest(c)
	if !ok {
		return
	}
	displayName, err := documents.NormalizeDisplayName(req.DisplayName)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	now := time.Now().UTC()
	updaterID := user.ID
	err = h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		folderHelper := documentFolderHandler{db: tx, auth: h.auth}
		if ok := folderHelper.lockFolderHierarchyWithDB(c, tx, documentRow.OrganizationUnitID); !ok {
			return errDocumentOperationResponseWritten
		}
		if ok := folderHelper.lockActiveOrganizationUnitWithDB(c, tx, documentRow.OrganizationUnitID, "document not found"); !ok {
			return errDocumentOperationResponseWritten
		}
		if _, ok := h.loadActiveDocumentWithActiveScopeWithDB(c, tx, id, true); !ok {
			return errDocumentOperationResponseWritten
		}
		result := tx.WithContext(c.Request.Context()).Model(&documents.Document{}).
			Where("id = ? AND deleted_at IS NULL", id).
			Updates(map[string]any{
				"display_name":       displayName,
				"updated_by_user_id": updaterID,
				"updated_at":         now,
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return tx.WithContext(c.Request.Context()).First(&documentRow, "id = ?", id).Error
	})
	if err != nil {
		if errors.Is(err, errDocumentOperationResponseWritten) {
			return
		}
		writeDocumentMutationError(c, err, "failed to update document")
		return
	}

	c.JSON(http.StatusOK, documentMetadataResponseFromModel(documentRow))
}

func (h *documentHandler) archive(c *gin.Context) {
	user, ok := requireAuthenticatedUser(c, h.auth)
	if !ok {
		return
	}
	if h.db == nil {
		writeError(c, http.StatusServiceUnavailable, "database is not configured")
		return
	}
	id, ok := parseDocumentID(c, c.Param("id"))
	if !ok {
		return
	}
	documentRow, ok := h.loadActiveDocumentWithActiveScope(c, id)
	if !ok {
		return
	}
	if ok := requireDocumentObjectPermissionForAuthenticatedUser(c, h.auth, user, "document.delete", documentRow.OrganizationUnitID); !ok {
		return
	}

	now := time.Now().UTC()
	updaterID := user.ID
	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		folderHelper := documentFolderHandler{db: tx, auth: h.auth}
		if ok := folderHelper.lockFolderHierarchyWithDB(c, tx, documentRow.OrganizationUnitID); !ok {
			return errDocumentOperationResponseWritten
		}
		if ok := folderHelper.lockActiveOrganizationUnitWithDB(c, tx, documentRow.OrganizationUnitID, "document not found"); !ok {
			return errDocumentOperationResponseWritten
		}
		if _, ok := h.loadActiveDocumentWithActiveScopeWithDB(c, tx, id, true); !ok {
			return errDocumentOperationResponseWritten
		}
		result := tx.WithContext(c.Request.Context()).Model(&documents.Document{}).
			Where("id = ? AND deleted_at IS NULL", id).
			Updates(map[string]any{
				"deleted_at":         now,
				"updated_by_user_id": updaterID,
				"updated_at":         now,
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, errDocumentOperationResponseWritten) {
			return
		}
		writeDocumentMutationError(c, err, "failed to archive document")
		return
	}

	c.JSON(http.StatusOK, okResponse{OK: true})
}

func (h *documentHandler) deleteStoredUpload(storageKey string) error {
	if h.storage == nil || storageKey == "" {
		return nil
	}
	return h.storage.Delete(storageKey)
}

func (h *documentHandler) cleanupStoredUploadOrWriteError(c *gin.Context, storageKey string) bool {
	if err := h.deleteStoredUpload(storageKey); err != nil {
		writeError(c, http.StatusInternalServerError, "failed to clean up document upload")
		return false
	}
	return true
}

func (h *documentHandler) revalidateActiveUploadScopeWithDB(c *gin.Context, db *gorm.DB, folderID string, organizationUnitID string) (documents.Folder, error) {
	var unit orgunits.Unit
	if err := activeOrganizationUnitLockQuery(c.Request.Context(), db, organizationUnitID).First(&unit).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return documents.Folder{}, errDocumentUploadScopeNotFound
		}
		return documents.Folder{}, err
	}

	var folder documents.Folder
	if err := db.WithContext(c.Request.Context()).
		Joins("JOIN organization_units ON organization_units.id = document_folders.organization_unit_id").
		Where("document_folders.id = ? AND document_folders.organization_unit_id = ? AND document_folders.deleted_at IS NULL AND organization_units.archived_at IS NULL", folderID, organizationUnitID).
		First(&folder).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return documents.Folder{}, errDocumentUploadScopeNotFound
		}
		return documents.Folder{}, err
	}
	return folder, nil
}

func bindDocumentUpdateRequest(c *gin.Context) (documentUpdateRequest, bool) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxDocumentMetadataRequestBytes)
	var req documentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return documentUpdateRequest{}, false
	}
	return req, true
}

func parseDocumentID(c *gin.Context, raw string) (string, bool) {
	return parseUUIDValue(c, raw, "invalid document id")
}

func (h *documentHandler) loadActiveDocumentWithActiveScope(c *gin.Context, documentID string) (documents.Document, bool) {
	return h.loadActiveDocumentWithActiveScopeWithDB(c, h.db, documentID, false)
}

func (h *documentHandler) loadActiveDocumentWithActiveScopeWithDB(c *gin.Context, db *gorm.DB, documentID string, lockDocument bool) (documents.Document, bool) {
	var documentRow documents.Document
	query := db.WithContext(c.Request.Context()).
		Joins("JOIN document_folders ON document_folders.id = documents.folder_id AND document_folders.organization_unit_id = documents.organization_unit_id").
		Joins("JOIN organization_units ON organization_units.id = documents.organization_unit_id").
		Where("documents.id = ? AND documents.deleted_at IS NULL AND document_folders.deleted_at IS NULL AND organization_units.archived_at IS NULL", documentID)
	if lockDocument && db.Dialector.Name() == "postgres" {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	if err := query.First(&documentRow).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusNotFound, "document not found")
			return documents.Document{}, false
		}
		writeError(c, http.StatusInternalServerError, "failed to load document")
		return documents.Document{}, false
	}
	return documentRow, true
}

func requireDocumentObjectPermissionForAuthenticatedUser(c *gin.Context, h *authHandler, user auth.User, permission string, organizationUnitID string) bool {
	allowed, err := rbac.NewResolver(h.db).Can(c.Request.Context(), rbac.Check{
		UserID:             user.ID,
		Permission:         permission,
		OrganizationUnitID: &organizationUnitID,
	})
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to check permission")
		return false
	}
	if !allowed {
		writeError(c, http.StatusNotFound, "document not found")
		return false
	}
	return true
}

func writeDocumentMutationError(c *gin.Context, err error, fallback string) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		writeError(c, http.StatusNotFound, "document not found")
		return
	}
	writeError(c, http.StatusInternalServerError, fallback)
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
