package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
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

const maxDocumentFolderRequestBytes int64 = 1 << 20

var errDocumentFolderResponseWritten = errors.New("document folder response already written")

type documentFolderHandler struct {
	db   *gorm.DB
	auth *authHandler
}

type documentFolderRequest struct {
	OrganizationUnitID string  `json:"organizationUnitId"`
	ParentID           *string `json:"parentId"`
	Name               string  `json:"name"`
	Description        *string `json:"description"`
}

type moveDocumentFolderRequest struct {
	ParentID *string
}

func (r *moveDocumentFolderRequest) UnmarshalJSON(data []byte) error {
	type rawRequest struct {
		ParentID *string `json:"parentId"`
	}
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if _, ok := raw["parentId"]; !ok {
		return errors.New("parentId is required")
	}
	var out rawRequest
	if err := json.Unmarshal(data, &out); err != nil {
		return err
	}
	r.ParentID = out.ParentID
	return nil
}

type documentFolderResponse struct {
	ID                 string                   `json:"id"`
	ParentID           *string                  `json:"parentId,omitempty"`
	OrganizationUnitID string                   `json:"organizationUnitId"`
	Name               string                   `json:"name"`
	Description        *string                  `json:"description,omitempty"`
	DeletedAt          *string                  `json:"deletedAt,omitempty"`
	CreatedAt          string                   `json:"createdAt"`
	UpdatedAt          string                   `json:"updatedAt"`
	Children           []documentFolderResponse `json:"children"`
}

type documentFolderListResponse struct {
	Folders []documentFolderResponse `json:"folders"`
}

type documentMetadataResponse struct {
	ID                 string  `json:"id"`
	FolderID           string  `json:"folderId"`
	OrganizationUnitID string  `json:"organizationUnitId"`
	OriginalFileName   string  `json:"originalFileName"`
	DisplayName        string  `json:"displayName"`
	MimeType           string  `json:"mimeType"`
	Extension          *string `json:"extension,omitempty"`
	SizeBytes          int64   `json:"sizeBytes"`
	SHA256Hash         string  `json:"sha256Hash"`
	DeletedAt          *string `json:"deletedAt,omitempty"`
	CreatedAt          string  `json:"createdAt"`
	UpdatedAt          string  `json:"updatedAt"`
}

type documentFolderContentsResponse struct {
	Folder    documentFolderResponse     `json:"folder"`
	Folders   []documentFolderResponse   `json:"folders"`
	Documents []documentMetadataResponse `json:"documents"`
}

type normalizedDocumentFolderInput struct {
	Name                string
	Description         *string
	DescriptionProvided bool
}

func newDocumentFolderHandler(options RouterOptions, auth *authHandler) *documentFolderHandler {
	return &documentFolderHandler{db: options.DB, auth: auth}
}

func (h *documentFolderHandler) listTree(c *gin.Context) {
	organizationUnitID, ok := parseDocumentFolderOrganizationUnitID(c, c.Query("organizationUnitId"))
	if !ok {
		return
	}
	if _, ok := requirePermission(c, h.auth, "document.view", &organizationUnitID); !ok {
		return
	}
	if ok := h.activeOrganizationUnitExists(c, organizationUnitID); !ok {
		return
	}
	var folders []documents.Folder
	if err := h.db.WithContext(c.Request.Context()).
		Where("organization_unit_id = ? AND deleted_at IS NULL", organizationUnitID).
		Order("name asc, id asc").
		Find(&folders).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to list document folders")
		return
	}
	c.JSON(http.StatusOK, documentFolderListResponse{Folders: documentFolderTreeResponse(organizationUnitID, documents.BuildFolderTree(folders))})
}

func (h *documentFolderHandler) create(c *gin.Context) {
	req, ok := bindDocumentFolderRequest(c)
	if !ok {
		return
	}
	organizationUnitID, ok := parseDocumentFolderOrganizationUnitID(c, req.OrganizationUnitID)
	if !ok {
		return
	}
	user, ok := requirePermission(c, h.auth, "document.create", &organizationUnitID)
	if !ok {
		return
	}
	if ok := h.activeOrganizationUnitExists(c, organizationUnitID); !ok {
		return
	}
	input, ok := normalizeDocumentFolderInput(c, req)
	if !ok {
		return
	}

	now := time.Now().UTC()
	updaterID := user.ID
	var folder documents.Folder
	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if ok := h.lockFolderHierarchyWithDB(c, tx, organizationUnitID); !ok {
			return errDocumentFolderResponseWritten
		}
		if ok := h.lockActiveOrganizationUnitWithDB(c, tx, organizationUnitID, "organization unit not found"); !ok {
			return errDocumentFolderResponseWritten
		}
		parentID, ok := h.validateOptionalActiveParentWithDB(c, tx, req.ParentID, organizationUnitID, user, []string{"document.create", "document.view"})
		if !ok {
			return errDocumentFolderResponseWritten
		}
		folder = documents.Folder{
			ParentID:           parentID,
			OrganizationUnitID: organizationUnitID,
			Name:               input.Name,
			Description:        input.Description,
			CreatedByUserID:    user.ID,
			UpdatedByUserID:    &updaterID,
			CreatedAt:          now,
			UpdatedAt:          now,
		}
		return tx.Create(&folder).Error
	})
	if err != nil {
		if errors.Is(err, errDocumentFolderResponseWritten) {
			return
		}
		writeDocumentFolderMutationError(c, err, "failed to create document folder")
		return
	}
	c.JSON(http.StatusCreated, documentFolderResponseFromModel(folder))
}

func (h *documentFolderHandler) update(c *gin.Context) {
	user, ok := requireAuthenticatedUser(c, h.auth)
	if !ok {
		return
	}
	id, ok := parseDocumentFolderID(c, c.Param("id"))
	if !ok {
		return
	}
	folder, ok := h.loadActiveFolderWithActiveOrganizationUnit(c, id)
	if !ok {
		return
	}
	if ok := requireDocumentFolderObjectPermissionForAuthenticatedUser(c, h.auth, user, "document.update", &folder.OrganizationUnitID); !ok {
		return
	}
	req, ok := bindDocumentFolderRequest(c)
	if !ok {
		return
	}
	if strings.TrimSpace(req.OrganizationUnitID) != "" {
		organizationUnitID, ok := parseDocumentFolderOrganizationUnitID(c, req.OrganizationUnitID)
		if !ok {
			return
		}
		if organizationUnitID != folder.OrganizationUnitID {
			writeError(c, http.StatusConflict, "document folder organization unit cannot be changed")
			return
		}
	}
	input, ok := normalizeDocumentFolderInput(c, req)
	if !ok {
		return
	}

	updates := map[string]any{
		"name":               input.Name,
		"updated_by_user_id": user.ID,
		"updated_at":         time.Now().UTC(),
	}
	if input.DescriptionProvided {
		updates["description"] = nullableStringValue(input.Description)
	}
	if err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if ok := h.lockActiveOrganizationUnitWithDB(c, tx, folder.OrganizationUnitID, "document folder not found"); !ok {
			return errDocumentFolderResponseWritten
		}
		result := tx.Model(&documents.Folder{}).
			Where("id = ? AND deleted_at IS NULL", id).
			Where("EXISTS (SELECT 1 FROM organization_units WHERE organization_units.id = document_folders.organization_unit_id AND organization_units.archived_at IS NULL)").
			Updates(updates)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return tx.First(&folder, "id = ?", id).Error
	}); err != nil {
		if errors.Is(err, errDocumentFolderResponseWritten) {
			return
		}
		writeDocumentFolderMutationError(c, err, "failed to update document folder")
		return
	}
	c.JSON(http.StatusOK, documentFolderResponseFromModel(folder))
}

func (h *documentFolderHandler) move(c *gin.Context) {
	user, ok := requireAuthenticatedUser(c, h.auth)
	if !ok {
		return
	}
	id, ok := parseDocumentFolderID(c, c.Param("id"))
	if !ok {
		return
	}
	folder, ok := h.loadActiveFolderWithActiveOrganizationUnit(c, id)
	if !ok {
		return
	}
	if ok := requireDocumentFolderObjectPermissionForAuthenticatedUser(c, h.auth, user, "document.update", &folder.OrganizationUnitID); !ok {
		return
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxDocumentFolderRequestBytes)
	var req moveDocumentFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return
	}

	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if ok := h.lockFolderHierarchyWithDB(c, tx, folder.OrganizationUnitID); !ok {
			return errDocumentFolderResponseWritten
		}
		if ok := h.lockActiveOrganizationUnitWithDB(c, tx, folder.OrganizationUnitID, "document folder not found"); !ok {
			return errDocumentFolderResponseWritten
		}
		lockedFolder, ok := h.loadActiveFolderWithActiveOrganizationUnitWithDB(c, tx, id)
		if !ok {
			return errDocumentFolderResponseWritten
		}
		folder = lockedFolder
		parentID, ok := h.validateOptionalActiveParentWithDB(c, tx, req.ParentID, folder.OrganizationUnitID, user, []string{"document.update", "document.view"})
		if !ok {
			return errDocumentFolderResponseWritten
		}
		if parentID != nil && *parentID == id {
			writeError(c, http.StatusConflict, "document folder cannot be moved under itself")
			return errDocumentFolderResponseWritten
		}
		if parentID != nil {
			wouldCreateCycle, ok := h.parentWouldCreateCycleWithDB(c, tx, folder.OrganizationUnitID, id, *parentID)
			if !ok {
				return errDocumentFolderResponseWritten
			}
			if wouldCreateCycle {
				writeError(c, http.StatusConflict, "document folder cannot be moved under its descendant")
				return errDocumentFolderResponseWritten
			}
		}
		result := tx.Model(&documents.Folder{}).Where("id = ? AND deleted_at IS NULL", id).Updates(map[string]any{
			"parent_id":          nullableStringValue(parentID),
			"updated_by_user_id": user.ID,
			"updated_at":         time.Now().UTC(),
		})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return tx.First(&folder, "id = ?", id).Error
	})
	if err != nil {
		if errors.Is(err, errDocumentFolderResponseWritten) {
			return
		}
		writeDocumentFolderMutationError(c, err, "failed to move document folder")
		return
	}
	c.JSON(http.StatusOK, documentFolderResponseFromModel(folder))
}

func (h *documentFolderHandler) archive(c *gin.Context) {
	user, ok := requireAuthenticatedUser(c, h.auth)
	if !ok {
		return
	}
	id, ok := parseDocumentFolderID(c, c.Param("id"))
	if !ok {
		return
	}
	folder, ok := h.loadActiveFolderWithActiveOrganizationUnit(c, id)
	if !ok {
		return
	}
	if ok := requireDocumentFolderObjectPermissionForAuthenticatedUser(c, h.auth, user, "document.delete", &folder.OrganizationUnitID); !ok {
		return
	}

	now := time.Now().UTC()
	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if ok := h.lockFolderHierarchyWithDB(c, tx, folder.OrganizationUnitID); !ok {
			return errDocumentFolderResponseWritten
		}
		if ok := h.lockActiveOrganizationUnitWithDB(c, tx, folder.OrganizationUnitID, "document folder not found"); !ok {
			return errDocumentFolderResponseWritten
		}
		lockedFolder, ok := h.loadActiveFolderWithActiveOrganizationUnitWithDB(c, tx, id)
		if !ok {
			return errDocumentFolderResponseWritten
		}
		folder = lockedFolder
		var folders []documents.Folder
		if err := tx.Where("organization_unit_id = ? AND deleted_at IS NULL", folder.OrganizationUnitID).Find(&folders).Error; err != nil {
			return err
		}
		found := false
		for _, activeFolder := range folders {
			if activeFolder.ID == id {
				found = true
				break
			}
		}
		if !found {
			return gorm.ErrRecordNotFound
		}
		folderIDs := []string{id}
		for descendantID := range documents.DescendantFolderIDs(id, folders) {
			folderIDs = append(folderIDs, descendantID)
		}
		updates := map[string]any{
			"deleted_at":         now,
			"updated_at":         now,
			"updated_by_user_id": user.ID,
		}
		if err := tx.Model(&documents.Folder{}).Where("id IN ? AND deleted_at IS NULL", folderIDs).Updates(updates).Error; err != nil {
			return err
		}
		return tx.Model(&documents.Document{}).Where("folder_id IN ? AND deleted_at IS NULL", folderIDs).Updates(updates).Error
	})
	if err != nil {
		if errors.Is(err, errDocumentFolderResponseWritten) {
			return
		}
		writeDocumentFolderMutationError(c, err, "failed to archive document folder")
		return
	}
	c.JSON(http.StatusOK, okResponse{OK: true})
}

func (h *documentFolderHandler) contents(c *gin.Context) {
	user, ok := requireAuthenticatedUser(c, h.auth)
	if !ok {
		return
	}
	id, ok := parseDocumentFolderID(c, c.Param("id"))
	if !ok {
		return
	}
	folder, ok := h.loadActiveFolderWithActiveOrganizationUnit(c, id)
	if !ok {
		return
	}
	if ok := requireDocumentFolderObjectPermissionForAuthenticatedUser(c, h.auth, user, "document.view", &folder.OrganizationUnitID); !ok {
		return
	}
	var childFolders []documents.Folder
	if err := h.db.WithContext(c.Request.Context()).
		Where("organization_unit_id = ? AND parent_id = ? AND deleted_at IS NULL", folder.OrganizationUnitID, id).
		Order("name asc, id asc").
		Find(&childFolders).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to list document folder contents")
		return
	}
	var documentRows []documents.Document
	if err := h.db.WithContext(c.Request.Context()).
		Where("organization_unit_id = ? AND folder_id = ? AND deleted_at IS NULL", folder.OrganizationUnitID, id).
		Order("display_name asc, id asc").
		Find(&documentRows).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to list document folder contents")
		return
	}
	c.JSON(http.StatusOK, documentFolderContentsResponse{
		Folder:    documentFolderResponseFromModel(folder),
		Folders:   documentFolderResponsesFromModels(childFolders),
		Documents: documentMetadataResponsesFromModels(documentRows),
	})
}

func bindDocumentFolderRequest(c *gin.Context) (documentFolderRequest, bool) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxDocumentFolderRequestBytes)
	var req documentFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return documentFolderRequest{}, false
	}
	return req, true
}

func normalizeDocumentFolderInput(c *gin.Context, req documentFolderRequest) (normalizedDocumentFolderInput, bool) {
	name, err := documents.NormalizeFolderName(req.Name)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return normalizedDocumentFolderInput{}, false
	}
	descriptionRaw := ""
	if req.Description != nil {
		descriptionRaw = *req.Description
	}
	return normalizedDocumentFolderInput{
		Name:                name,
		Description:         documents.NormalizeDescription(descriptionRaw),
		DescriptionProvided: req.Description != nil,
	}, true
}

func parseDocumentFolderID(c *gin.Context, raw string) (string, bool) {
	return parseUUIDValue(c, raw, "invalid document folder id")
}

func parseDocumentFolderOrganizationUnitID(c *gin.Context, raw string) (string, bool) {
	return parseUUIDValue(c, raw, "invalid organization unit id")
}

func (h *documentFolderHandler) activeOrganizationUnitExists(c *gin.Context, organizationUnitID string) bool {
	return h.activeOrganizationUnitExistsWithDB(c, h.db, organizationUnitID)
}

func (h *documentFolderHandler) activeOrganizationUnitExistsWithDB(c *gin.Context, db *gorm.DB, organizationUnitID string) bool {
	var count int64
	if err := db.WithContext(c.Request.Context()).Model(&orgunits.Unit{}).Where("id = ? AND archived_at IS NULL", organizationUnitID).Count(&count).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to validate organization unit")
		return false
	}
	if count == 0 {
		writeError(c, http.StatusNotFound, "organization unit not found")
		return false
	}
	return true
}

func (h *documentFolderHandler) lockActiveOrganizationUnitWithDB(c *gin.Context, db *gorm.DB, organizationUnitID string, notFoundMessage string) bool {
	var unit orgunits.Unit
	if err := activeOrganizationUnitLockQuery(c.Request.Context(), db, organizationUnitID).First(&unit).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusNotFound, notFoundMessage)
			return false
		}
		writeError(c, http.StatusInternalServerError, "failed to validate organization unit")
		return false
	}
	return true
}

func activeOrganizationUnitLockQuery(ctx context.Context, db *gorm.DB, organizationUnitID string) *gorm.DB {
	query := db.WithContext(ctx).
		Model(&orgunits.Unit{}).
		Select("id").
		Where("id = ? AND archived_at IS NULL", organizationUnitID)
	if db.Dialector.Name() == "postgres" {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	return query
}

func (h *documentFolderHandler) loadActiveFolderWithActiveOrganizationUnit(c *gin.Context, folderID string) (documents.Folder, bool) {
	return h.loadActiveFolderWithActiveOrganizationUnitWithDB(c, h.db, folderID)
}

func (h *documentFolderHandler) loadActiveFolderWithActiveOrganizationUnitWithDB(c *gin.Context, db *gorm.DB, folderID string) (documents.Folder, bool) {
	var folder documents.Folder
	if err := db.WithContext(c.Request.Context()).
		Joins("JOIN organization_units ON organization_units.id = document_folders.organization_unit_id").
		Where("document_folders.id = ? AND document_folders.deleted_at IS NULL AND organization_units.archived_at IS NULL", folderID).
		First(&folder).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusNotFound, "document folder not found")
			return documents.Folder{}, false
		}
		writeError(c, http.StatusInternalServerError, "failed to load document folder")
		return documents.Folder{}, false
	}
	return folder, true
}

func (h *documentFolderHandler) lockFolderHierarchyWithDB(c *gin.Context, db *gorm.DB, organizationUnitID string) bool {
	switch db.Dialector.Name() {
	case "postgres":
		key := "document_folders:" + organizationUnitID
		if err := db.WithContext(c.Request.Context()).Exec("SELECT pg_advisory_xact_lock(hashtext(?))", key).Error; err != nil {
			writeError(c, http.StatusInternalServerError, "failed to lock document folder hierarchy")
			return false
		}
	case "sqlite":
		return true
	default:
		return true
	}
	return true
}

func requireDocumentFolderObjectPermissionForAuthenticatedUser(c *gin.Context, h *authHandler, user auth.User, permission string, organizationUnitID *string) bool {
	allowed, ok := checkDocumentFolderPermissionForAuthenticatedUser(c, h, user, permission, organizationUnitID)
	if !ok {
		return false
	}
	if !allowed {
		writeError(c, http.StatusNotFound, "document folder not found")
		return false
	}
	return true
}

func checkDocumentFolderPermissionForAuthenticatedUser(c *gin.Context, h *authHandler, user auth.User, permission string, organizationUnitID *string) (bool, bool) {
	allowed, err := rbac.NewResolver(h.db).Can(c.Request.Context(), rbac.Check{
		UserID:             user.ID,
		Permission:         permission,
		OrganizationUnitID: organizationUnitID,
	})
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to check permission")
		return false, false
	}
	return allowed, true
}

func (h *documentFolderHandler) validateOptionalActiveParentWithDB(c *gin.Context, db *gorm.DB, raw *string, organizationUnitID string, user auth.User, parentAccessPermissions []string) (*string, bool) {
	if raw == nil || strings.TrimSpace(*raw) == "" {
		return nil, true
	}
	parentID, ok := parseDocumentFolderID(c, *raw)
	if !ok {
		return nil, false
	}
	var parent documents.Folder
	if err := db.WithContext(c.Request.Context()).
		Joins("JOIN organization_units ON organization_units.id = document_folders.organization_unit_id").
		Where("document_folders.id = ? AND document_folders.deleted_at IS NULL AND organization_units.archived_at IS NULL", parentID).
		First(&parent).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusNotFound, "parent document folder not found")
			return nil, false
		}
		writeError(c, http.StatusInternalServerError, "failed to validate parent document folder")
		return nil, false
	}
	if parent.OrganizationUnitID != organizationUnitID {
		canAccessParent, ok := h.userHasAnyDocumentFolderPermissionWithDB(c, db, user, parentAccessPermissions, parent.OrganizationUnitID)
		if !ok {
			return nil, false
		}
		if !canAccessParent {
			writeError(c, http.StatusNotFound, "parent document folder not found")
			return nil, false
		}
		writeError(c, http.StatusConflict, "parent document folder belongs to another organization unit")
		return nil, false
	}
	return &parentID, true
}

func (h *documentFolderHandler) userHasAnyDocumentFolderPermissionWithDB(c *gin.Context, db *gorm.DB, user auth.User, permissions []string, organizationUnitID string) (bool, bool) {
	resolver := rbac.NewResolver(db)
	for _, permission := range permissions {
		allowed, err := resolver.Can(c.Request.Context(), rbac.Check{
			UserID:             user.ID,
			Permission:         permission,
			OrganizationUnitID: &organizationUnitID,
		})
		if err != nil {
			writeError(c, http.StatusInternalServerError, "failed to check permission")
			return false, false
		}
		if allowed {
			return true, true
		}
	}
	return false, true
}

func (h *documentFolderHandler) parentWouldCreateCycleWithDB(c *gin.Context, db *gorm.DB, organizationUnitID string, folderID string, parentID string) (bool, bool) {
	var folders []documents.Folder
	if err := db.WithContext(c.Request.Context()).Where("organization_unit_id = ? AND deleted_at IS NULL", organizationUnitID).Find(&folders).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to validate document folder move")
		return false, false
	}
	descendants := documents.DescendantFolderIDs(folderID, folders)
	return descendants[parentID], true
}

func documentFolderTreeResponse(organizationUnitID string, nodes []documents.FolderTreeNode) []documentFolderResponse {
	out := make([]documentFolderResponse, 0, len(nodes))
	for _, node := range nodes {
		out = append(out, documentFolderResponse{
			ID:                 node.ID,
			ParentID:           node.ParentID,
			OrganizationUnitID: organizationUnitID,
			Name:               node.Name,
			Description:        node.Description,
			CreatedAt:          node.CreatedAt,
			UpdatedAt:          node.UpdatedAt,
			Children:           documentFolderTreeResponse(organizationUnitID, node.Children),
		})
	}
	return out
}

func documentFolderResponseFromModel(folder documents.Folder) documentFolderResponse {
	var deletedAt *string
	if folder.DeletedAt != nil {
		value := folder.DeletedAt.UTC().Format(time.RFC3339Nano)
		deletedAt = &value
	}
	return documentFolderResponse{
		ID:                 folder.ID,
		ParentID:           folder.ParentID,
		OrganizationUnitID: folder.OrganizationUnitID,
		Name:               folder.Name,
		Description:        folder.Description,
		DeletedAt:          deletedAt,
		CreatedAt:          folder.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:          folder.UpdatedAt.UTC().Format(time.RFC3339Nano),
		Children:           []documentFolderResponse{},
	}
}

func documentFolderResponsesFromModels(folders []documents.Folder) []documentFolderResponse {
	out := make([]documentFolderResponse, 0, len(folders))
	for _, folder := range folders {
		out = append(out, documentFolderResponseFromModel(folder))
	}
	return out
}

func documentMetadataResponsesFromModels(documentRows []documents.Document) []documentMetadataResponse {
	out := make([]documentMetadataResponse, 0, len(documentRows))
	for _, documentRow := range documentRows {
		out = append(out, documentMetadataResponseFromModel(documentRow))
	}
	return out
}

func documentMetadataResponseFromModel(documentRow documents.Document) documentMetadataResponse {
	var deletedAt *string
	if documentRow.DeletedAt != nil {
		value := documentRow.DeletedAt.UTC().Format(time.RFC3339Nano)
		deletedAt = &value
	}
	return documentMetadataResponse{
		ID:                 documentRow.ID,
		FolderID:           documentRow.FolderID,
		OrganizationUnitID: documentRow.OrganizationUnitID,
		OriginalFileName:   documentRow.OriginalFileName,
		DisplayName:        documentRow.DisplayName,
		MimeType:           documentRow.MimeType,
		Extension:          documentRow.Extension,
		SizeBytes:          documentRow.SizeBytes,
		SHA256Hash:         documentRow.SHA256Hash,
		DeletedAt:          deletedAt,
		CreatedAt:          documentRow.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:          documentRow.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
}

func writeDocumentFolderMutationError(c *gin.Context, err error, fallback string) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		writeError(c, http.StatusNotFound, "document folder not found")
		return
	}
	if isUniqueConstraintError(err) || isDocumentFolderUniqueConstraintError(err) {
		writeError(c, http.StatusConflict, "active document folder name already exists")
		return
	}
	writeError(c, http.StatusInternalServerError, fallback)
}

func isDocumentFolderUniqueConstraintError(err error) bool {
	message := err.Error()
	return strings.Contains(message, "idx_document_folders_root_name_unique") ||
		strings.Contains(message, "idx_document_folders_child_name_unique") ||
		strings.Contains(message, "document_folders.organization_unit_id, document_folders.name") ||
		strings.Contains(message, "document_folders.organization_unit_id, document_folders.parent_id, document_folders.name")
}
