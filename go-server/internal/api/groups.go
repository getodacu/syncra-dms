package api

import (
	"errors"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"ai.ro/syncra/dms/internal/auth"
	"ai.ro/syncra/dms/internal/orgunits"
	"ai.ro/syncra/dms/internal/rbac"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type groupHandler struct {
	db   *gorm.DB
	auth *authHandler
}

type createGroupRequest struct {
	Name               string  `json:"name"`
	Code               string  `json:"code"`
	Description        *string `json:"description"`
	OrganizationUnitID *string `json:"organizationUnitId"`
	IsActive           *bool   `json:"isActive"`
}

type updateGroupRequest struct {
	Name               *string `json:"name"`
	Code               *string `json:"code"`
	Description        *string `json:"description"`
	OrganizationUnitID *string `json:"organizationUnitId"`
	IsActive           *bool   `json:"isActive"`
}

type assignGroupUserRequest struct {
	UserID string `json:"userId"`
}

type assignGroupRoleRequest struct {
	RoleID             string  `json:"roleId"`
	ScopeType          string  `json:"scopeType"`
	OrganizationUnitID *string `json:"organizationUnitId"`
}

type groupResponse struct {
	ID                 string  `json:"id"`
	Name               string  `json:"name"`
	Code               string  `json:"code"`
	Description        *string `json:"description,omitempty"`
	OrganizationUnitID *string `json:"organizationUnitId,omitempty"`
	IsActive           bool    `json:"isActive"`
	CreatedAt          string  `json:"createdAt"`
	UpdatedAt          string  `json:"updatedAt"`
}

type groupListResponse struct {
	Groups []groupResponse `json:"groups"`
}

type groupRoleAssignmentResponse struct {
	ID                 string  `json:"id"`
	GroupID            string  `json:"groupId"`
	RoleID             string  `json:"roleId"`
	ScopeType          string  `json:"scopeType"`
	OrganizationUnitID *string `json:"organizationUnitId,omitempty"`
	CreatedAt          string  `json:"createdAt"`
}

type groupUserAssignmentResponse struct {
	GroupID   string `json:"groupId"`
	UserID    string `json:"userId"`
	CreatedAt string `json:"createdAt"`
}

func newGroupHandler(options RouterOptions, auth *authHandler) *groupHandler {
	return &groupHandler{db: options.DB, auth: auth}
}

func (h *groupHandler) list(c *gin.Context) {
	if _, ok := requirePermission(c, h.auth, "group.view", nil); !ok {
		return
	}
	var groups []rbac.Group
	if err := h.db.WithContext(c.Request.Context()).
		Order("name asc, id asc").
		Find(&groups).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to list groups")
		return
	}
	out := make([]groupResponse, 0, len(groups))
	for _, group := range groups {
		out = append(out, groupResponseFromModel(group))
	}
	c.JSON(http.StatusOK, groupListResponse{Groups: out})
}

func (h *groupHandler) get(c *gin.Context) {
	if _, ok := requirePermission(c, h.auth, "group.view", nil); !ok {
		return
	}
	groupID, ok := parseUUIDValue(c, c.Param("id"), "invalid group id")
	if !ok {
		return
	}
	group, ok := h.loadGroup(c, groupID)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, groupResponseFromModel(group))
}

func (h *groupHandler) create(c *gin.Context) {
	if _, ok := requirePermission(c, h.auth, "group.create", nil); !ok {
		return
	}
	var req createGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	name, ok := normalizeGroupName(c, req.Name)
	if !ok {
		return
	}
	code, err := rbac.NormalizeCode(req.Code)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	organizationUnitID, ok := h.validOptionalOrganizationUnitID(c, req.OrganizationUnitID)
	if !ok {
		return
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	now := time.Now().UTC()
	group := rbac.Group{
		Name:               name,
		Code:               code,
		Description:        normalizeOptionalUserText(req.Description),
		OrganizationUnitID: organizationUnitID,
		IsActive:           isActive,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	if err := h.db.WithContext(c.Request.Context()).Create(&group).Error; err != nil {
		writeGroupMutationError(c, err, "failed to create group")
		return
	}
	c.JSON(http.StatusCreated, groupResponseFromModel(group))
}

func (h *groupHandler) update(c *gin.Context) {
	if _, ok := requirePermission(c, h.auth, "group.update", nil); !ok {
		return
	}
	groupID, ok := parseUUIDValue(c, c.Param("id"), "invalid group id")
	if !ok {
		return
	}
	var req updateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	updates := map[string]any{"updated_at": time.Now().UTC()}
	if req.Name != nil {
		name, ok := normalizeGroupName(c, *req.Name)
		if !ok {
			return
		}
		updates["name"] = name
	}
	if req.Code != nil {
		code, err := rbac.NormalizeCode(*req.Code)
		if err != nil {
			writeError(c, http.StatusBadRequest, err.Error())
			return
		}
		updates["code"] = code
	}
	if req.Description != nil {
		updates["description"] = nullableStringValue(normalizeOptionalUserText(req.Description))
	}
	if req.OrganizationUnitID != nil {
		organizationUnitID, ok := h.validOptionalOrganizationUnitID(c, req.OrganizationUnitID)
		if !ok {
			return
		}
		updates["organization_unit_id"] = nullableStringValue(organizationUnitID)
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	var group rbac.Group
	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&rbac.Group{}).Where("id = ?", groupID).Updates(updates)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return tx.First(&group, "id = ?", groupID).Error
	})
	if err != nil {
		writeGroupMutationError(c, err, "failed to update group")
		return
	}
	c.JSON(http.StatusOK, groupResponseFromModel(group))
}

func (h *groupHandler) delete(c *gin.Context) {
	if _, ok := requirePermission(c, h.auth, "group.delete", nil); !ok {
		return
	}
	groupID, ok := parseUUIDValue(c, c.Param("id"), "invalid group id")
	if !ok {
		return
	}
	if _, ok := h.loadGroup(c, groupID); !ok {
		return
	}
	inUse, err := h.groupInUse(c, groupID)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to validate group usage")
		return
	}
	if inUse {
		writeError(c, http.StatusConflict, "group is assigned")
		return
	}
	if err := h.db.WithContext(c.Request.Context()).Delete(&rbac.Group{}, "id = ?", groupID).Error; err != nil {
		writeGroupMutationError(c, err, "failed to delete group")
		return
	}
	c.JSON(http.StatusOK, okResponse{OK: true})
}

func (h *groupHandler) addUser(c *gin.Context) {
	if _, ok := requirePermission(c, h.auth, "group.manage_users", nil); !ok {
		return
	}
	groupID, ok := parseUUIDValue(c, c.Param("id"), "invalid group id")
	if !ok {
		return
	}
	if _, ok := h.loadGroup(c, groupID); !ok {
		return
	}
	var req assignGroupUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	userID, ok := parseUUIDValue(c, req.UserID, "invalid user id")
	if !ok {
		return
	}
	if !h.userExists(c, userID) {
		return
	}
	now := time.Now().UTC()
	assignment := rbac.GroupUser{GroupID: groupID, UserID: userID, CreatedAt: now}
	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&assignment).Error; err != nil {
			return err
		}
		return tx.Where("group_id = ? AND user_id = ?", groupID, userID).First(&assignment).Error
	})
	if err != nil {
		writeGroupMutationError(c, err, "failed to add group user")
		return
	}
	c.JSON(http.StatusCreated, groupUserAssignmentResponseFromModel(assignment))
}

func (h *groupHandler) removeUser(c *gin.Context) {
	if _, ok := requirePermission(c, h.auth, "group.manage_users", nil); !ok {
		return
	}
	groupID, ok := parseUUIDValue(c, c.Param("id"), "invalid group id")
	if !ok {
		return
	}
	userID, ok := parseUUIDValue(c, c.Param("userId"), "invalid user id")
	if !ok {
		return
	}
	result := h.db.WithContext(c.Request.Context()).Where("group_id = ? AND user_id = ?", groupID, userID).Delete(&rbac.GroupUser{})
	if result.Error != nil {
		writeError(c, http.StatusInternalServerError, "failed to remove group user")
		return
	}
	if result.RowsAffected == 0 {
		writeError(c, http.StatusNotFound, "group user assignment not found")
		return
	}
	c.JSON(http.StatusOK, okResponse{OK: true})
}

func (h *groupHandler) assignRole(c *gin.Context) {
	if _, ok := requirePermission(c, h.auth, "group.assign_roles", nil); !ok {
		return
	}
	groupID, ok := parseUUIDValue(c, c.Param("id"), "invalid group id")
	if !ok {
		return
	}
	if _, ok := h.loadGroup(c, groupID); !ok {
		return
	}
	var req assignGroupRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	roleID, ok := parseUUIDValue(c, req.RoleID, "invalid role id")
	if !ok {
		return
	}
	if !h.roleExists(c, roleID) {
		return
	}
	scope := rbac.ScopeType(strings.TrimSpace(req.ScopeType))
	if scope == "" {
		scope = rbac.ScopeGlobal
	}
	organizationUnitID, ok := h.validScopedOrganizationUnitID(c, scope, req.OrganizationUnitID)
	if !ok {
		return
	}
	now := time.Now().UTC()
	assignment := rbac.GroupRole{
		GroupID:            groupID,
		RoleID:             roleID,
		ScopeType:          scope,
		OrganizationUnitID: organizationUnitID,
		CreatedAt:          now,
	}
	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&assignment).Error; err != nil {
			return err
		}
		query := tx.Where("group_id = ? AND role_id = ? AND scope_type = ?", groupID, roleID, scope)
		if organizationUnitID == nil {
			query = query.Where("organization_unit_id IS NULL")
		} else {
			query = query.Where("organization_unit_id = ?", *organizationUnitID)
		}
		return query.First(&assignment).Error
	})
	if err != nil {
		writeGroupMutationError(c, err, "failed to assign group role")
		return
	}
	c.JSON(http.StatusCreated, groupRoleAssignmentResponseFromModel(assignment))
}

func (h *groupHandler) removeRole(c *gin.Context) {
	if _, ok := requirePermission(c, h.auth, "group.assign_roles", nil); !ok {
		return
	}
	groupID, ok := parseUUIDValue(c, c.Param("id"), "invalid group id")
	if !ok {
		return
	}
	assignmentID, ok := parseUUIDValue(c, c.Param("assignmentId"), "invalid role assignment id")
	if !ok {
		return
	}
	result := h.db.WithContext(c.Request.Context()).Where("id = ? AND group_id = ?", assignmentID, groupID).Delete(&rbac.GroupRole{})
	if result.Error != nil {
		writeError(c, http.StatusInternalServerError, "failed to remove group role")
		return
	}
	if result.RowsAffected == 0 {
		writeError(c, http.StatusNotFound, "group role assignment not found")
		return
	}
	c.JSON(http.StatusOK, okResponse{OK: true})
}

func (h *groupHandler) loadGroup(c *gin.Context, groupID string) (rbac.Group, bool) {
	var group rbac.Group
	if err := h.db.WithContext(c.Request.Context()).First(&group, "id = ?", groupID).Error; err != nil {
		writeGroupMutationError(c, err, "failed to load group")
		return rbac.Group{}, false
	}
	return group, true
}

func (h *groupHandler) userExists(c *gin.Context, userID string) bool {
	var count int64
	if err := h.db.WithContext(c.Request.Context()).Model(&auth.User{}).Where("id = ? AND deleted_at IS NULL", userID).Count(&count).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to validate user")
		return false
	}
	if count == 0 {
		writeError(c, http.StatusNotFound, "user not found")
		return false
	}
	return true
}

func (h *groupHandler) roleExists(c *gin.Context, roleID string) bool {
	var count int64
	if err := h.db.WithContext(c.Request.Context()).Model(&rbac.Role{}).Where("id = ? AND is_active = ?", roleID, true).Count(&count).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to validate role")
		return false
	}
	if count == 0 {
		writeError(c, http.StatusNotFound, "role not found")
		return false
	}
	return true
}

func (h *groupHandler) validOptionalOrganizationUnitID(c *gin.Context, raw *string) (*string, bool) {
	if raw == nil || strings.TrimSpace(*raw) == "" {
		return nil, true
	}
	id, ok := parseUUIDValue(c, *raw, "invalid organization unit id")
	if !ok {
		return nil, false
	}
	var count int64
	if err := h.db.WithContext(c.Request.Context()).Model(&orgunits.Unit{}).Where("id = ? AND archived_at IS NULL", id).Count(&count).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to validate organization unit")
		return nil, false
	}
	if count == 0 {
		writeError(c, http.StatusNotFound, "organization unit not found")
		return nil, false
	}
	return &id, true
}

func (h *groupHandler) validScopedOrganizationUnitID(c *gin.Context, scope rbac.ScopeType, raw *string) (*string, bool) {
	if !scope.Valid() {
		writeError(c, http.StatusBadRequest, "scope type is invalid")
		return nil, false
	}
	if scope == rbac.ScopeGlobal {
		if raw != nil && strings.TrimSpace(*raw) != "" {
			writeError(c, http.StatusBadRequest, "global scope must not include an organization unit")
			return nil, false
		}
		return nil, true
	}
	if raw == nil || strings.TrimSpace(*raw) == "" {
		writeError(c, http.StatusBadRequest, "organization unit id is required")
		return nil, false
	}
	return h.validOptionalOrganizationUnitID(c, raw)
}

func (h *groupHandler) groupInUse(c *gin.Context, groupID string) (bool, error) {
	tables := []string{"group_users", "group_roles"}
	for _, table := range tables {
		var count int64
		if err := h.db.WithContext(c.Request.Context()).Table(table).Where("group_id = ?", groupID).Count(&count).Error; err != nil {
			return false, err
		}
		if count > 0 {
			return true, nil
		}
	}
	return false, nil
}

func normalizeGroupName(c *gin.Context, raw string) (string, bool) {
	name := strings.TrimSpace(raw)
	if name == "" {
		writeError(c, http.StatusBadRequest, "name is required")
		return "", false
	}
	if utf8.RuneCountInString(name) > 160 {
		writeError(c, http.StatusBadRequest, "name must be at most 160 characters")
		return "", false
	}
	return name, true
}

func groupResponseFromModel(group rbac.Group) groupResponse {
	return groupResponse{
		ID:                 group.ID,
		Name:               group.Name,
		Code:               group.Code,
		Description:        group.Description,
		OrganizationUnitID: group.OrganizationUnitID,
		IsActive:           group.IsActive,
		CreatedAt:          group.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:          group.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
}

func groupRoleAssignmentResponseFromModel(assignment rbac.GroupRole) groupRoleAssignmentResponse {
	return groupRoleAssignmentResponse{
		ID:                 assignment.ID,
		GroupID:            assignment.GroupID,
		RoleID:             assignment.RoleID,
		ScopeType:          string(assignment.ScopeType),
		OrganizationUnitID: assignment.OrganizationUnitID,
		CreatedAt:          assignment.CreatedAt.UTC().Format(time.RFC3339Nano),
	}
}

func groupUserAssignmentResponseFromModel(assignment rbac.GroupUser) groupUserAssignmentResponse {
	return groupUserAssignmentResponse{
		GroupID:   assignment.GroupID,
		UserID:    assignment.UserID,
		CreatedAt: assignment.CreatedAt.UTC().Format(time.RFC3339Nano),
	}
}

func writeGroupMutationError(c *gin.Context, err error, fallback string) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		writeError(c, http.StatusNotFound, "group not found")
		return
	}
	if isUniqueConstraintError(err) {
		writeError(c, http.StatusConflict, "group code already exists")
		return
	}
	writeError(c, http.StatusInternalServerError, fallback)
}
