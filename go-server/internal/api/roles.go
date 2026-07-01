package api

import (
	"errors"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"ai.ro/syncra/dms/internal/rbac"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type roleHandler struct {
	db   *gorm.DB
	auth *authHandler
}

type createRoleRequest struct {
	Name        string  `json:"name"`
	Code        string  `json:"code"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"isActive"`
}

type updateRoleRequest struct {
	Name        *string `json:"name"`
	Code        *string `json:"code"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"isActive"`
}

type assignRolePermissionRequest struct {
	PermissionID string `json:"permissionId"`
}

type roleResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Code        string  `json:"code"`
	Description *string `json:"description,omitempty"`
	IsSystem    bool    `json:"isSystem"`
	IsActive    bool    `json:"isActive"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

type roleListResponse struct {
	Roles []roleResponse `json:"roles"`
}

func newRoleHandler(options RouterOptions, auth *authHandler) *roleHandler {
	return &roleHandler{db: options.DB, auth: auth}
}

func (h *roleHandler) list(c *gin.Context) {
	if _, ok := requirePermission(c, h.auth, "role.view", nil); !ok {
		return
	}
	var roles []rbac.Role
	if err := h.db.WithContext(c.Request.Context()).
		Order("name asc, id asc").
		Find(&roles).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to list roles")
		return
	}
	out := make([]roleResponse, 0, len(roles))
	for _, role := range roles {
		out = append(out, roleResponseFromModel(role))
	}
	c.JSON(http.StatusOK, roleListResponse{Roles: out})
}

func (h *roleHandler) get(c *gin.Context) {
	if _, ok := requirePermission(c, h.auth, "role.view", nil); !ok {
		return
	}
	roleID, ok := parseUUIDValue(c, c.Param("id"), "invalid role id")
	if !ok {
		return
	}
	role, ok := h.loadRole(c, roleID)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, roleResponseFromModel(role))
}

func (h *roleHandler) create(c *gin.Context) {
	if _, ok := requirePermission(c, h.auth, "role.create", nil); !ok {
		return
	}
	var req createRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	name, ok := normalizeRoleName(c, req.Name)
	if !ok {
		return
	}
	code, err := rbac.NormalizeCode(req.Code)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	now := time.Now().UTC()
	role := rbac.Role{
		Name:        name,
		Code:        code,
		Description: normalizeOptionalUserText(req.Description),
		IsSystem:    false,
		IsActive:    isActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := h.db.WithContext(c.Request.Context()).Create(&role).Error; err != nil {
		writeRoleMutationError(c, err, "failed to create role")
		return
	}
	c.JSON(http.StatusCreated, roleResponseFromModel(role))
}

func (h *roleHandler) update(c *gin.Context) {
	if _, ok := requirePermission(c, h.auth, "role.update", nil); !ok {
		return
	}
	roleID, ok := parseUUIDValue(c, c.Param("id"), "invalid role id")
	if !ok {
		return
	}
	var req updateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	role, ok := h.loadRole(c, roleID)
	if !ok {
		return
	}
	updates := map[string]any{"updated_at": time.Now().UTC()}
	if req.Name != nil {
		name, ok := normalizeRoleName(c, *req.Name)
		if !ok {
			return
		}
		updates["name"] = name
	}
	if req.Code != nil {
		if role.IsSystem {
			writeError(c, http.StatusForbidden, "system role code cannot be changed")
			return
		}
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
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	result := h.db.WithContext(c.Request.Context()).Model(&rbac.Role{}).Where("id = ?", roleID).Updates(updates)
	if result.Error != nil {
		writeRoleMutationError(c, result.Error, "failed to update role")
		return
	}
	if result.RowsAffected == 0 {
		writeError(c, http.StatusNotFound, "role not found")
		return
	}
	if err := h.db.WithContext(c.Request.Context()).First(&role, "id = ?", roleID).Error; err != nil {
		writeRoleMutationError(c, err, "failed to load role")
		return
	}
	c.JSON(http.StatusOK, roleResponseFromModel(role))
}

func (h *roleHandler) delete(c *gin.Context) {
	if _, ok := requirePermission(c, h.auth, "role.delete", nil); !ok {
		return
	}
	roleID, ok := parseUUIDValue(c, c.Param("id"), "invalid role id")
	if !ok {
		return
	}
	role, ok := h.loadRole(c, roleID)
	if !ok {
		return
	}
	if role.IsSystem {
		writeError(c, http.StatusForbidden, "system role cannot be deleted")
		return
	}
	inUse, err := h.roleInUse(c, roleID)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to validate role usage")
		return
	}
	if inUse {
		writeError(c, http.StatusConflict, "role is assigned")
		return
	}
	if err := h.db.WithContext(c.Request.Context()).Delete(&rbac.Role{}, "id = ?", roleID).Error; err != nil {
		writeRoleMutationError(c, err, "failed to delete role")
		return
	}
	c.JSON(http.StatusOK, okResponse{OK: true})
}

func (h *roleHandler) listPermissions(c *gin.Context) {
	if _, ok := requirePermission(c, h.auth, "role.view", nil); !ok {
		return
	}
	roleID, ok := parseUUIDValue(c, c.Param("id"), "invalid role id")
	if !ok {
		return
	}
	if _, ok := h.loadRole(c, roleID); !ok {
		return
	}
	permissions, ok := h.permissionsForRole(c, roleID)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, permissionListResponse{Permissions: permissions})
}

func (h *roleHandler) assignPermission(c *gin.Context) {
	if _, ok := requirePermission(c, h.auth, "role.assign_permissions", nil); !ok {
		return
	}
	roleID, ok := parseUUIDValue(c, c.Param("id"), "invalid role id")
	if !ok {
		return
	}
	if _, ok := h.loadRole(c, roleID); !ok {
		return
	}
	var req assignRolePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	permissionID, ok := parseUUIDValue(c, req.PermissionID, "invalid permission id")
	if !ok {
		return
	}
	permission, ok := h.loadPermission(c, permissionID)
	if !ok {
		return
	}
	link := rbac.RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
		CreatedAt:    time.Now().UTC(),
	}
	if err := h.db.WithContext(c.Request.Context()).Clauses(clause.OnConflict{DoNothing: true}).Create(&link).Error; err != nil {
		writeRoleMutationError(c, err, "failed to assign permission")
		return
	}
	c.JSON(http.StatusCreated, permissionResponseFromModel(permission))
}

func (h *roleHandler) removePermission(c *gin.Context) {
	if _, ok := requirePermission(c, h.auth, "role.assign_permissions", nil); !ok {
		return
	}
	roleID, ok := parseUUIDValue(c, c.Param("id"), "invalid role id")
	if !ok {
		return
	}
	permissionID, ok := parseUUIDValue(c, c.Param("permissionId"), "invalid permission id")
	if !ok {
		return
	}
	result := h.db.WithContext(c.Request.Context()).Where("role_id = ? AND permission_id = ?", roleID, permissionID).Delete(&rbac.RolePermission{})
	if result.Error != nil {
		writeError(c, http.StatusInternalServerError, "failed to remove permission")
		return
	}
	if result.RowsAffected == 0 {
		writeError(c, http.StatusNotFound, "role permission not found")
		return
	}
	c.JSON(http.StatusOK, okResponse{OK: true})
}

func (h *roleHandler) loadRole(c *gin.Context, roleID string) (rbac.Role, bool) {
	var role rbac.Role
	if err := h.db.WithContext(c.Request.Context()).First(&role, "id = ?", roleID).Error; err != nil {
		writeRoleMutationError(c, err, "failed to load role")
		return rbac.Role{}, false
	}
	return role, true
}

func (h *roleHandler) loadPermission(c *gin.Context, permissionID string) (rbac.Permission, bool) {
	var permission rbac.Permission
	if err := h.db.WithContext(c.Request.Context()).First(&permission, "id = ?", permissionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusNotFound, "permission not found")
			return rbac.Permission{}, false
		}
		writeError(c, http.StatusInternalServerError, "failed to load permission")
		return rbac.Permission{}, false
	}
	return permission, true
}

func (h *roleHandler) permissionsForRole(c *gin.Context, roleID string) ([]permissionResponse, bool) {
	var permissions []rbac.Permission
	if err := h.db.WithContext(c.Request.Context()).
		Table("permissions").
		Select("permissions.*").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Order("permissions.category asc, permissions.code asc").
		Scan(&permissions).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to list role permissions")
		return nil, false
	}
	out := make([]permissionResponse, 0, len(permissions))
	for _, permission := range permissions {
		out = append(out, permissionResponseFromModel(permission))
	}
	return out, true
}

func (h *roleHandler) roleInUse(c *gin.Context, roleID string) (bool, error) {
	tables := []string{"user_roles", "group_roles", "organization_unit_roles"}
	for _, table := range tables {
		var count int64
		if err := h.db.WithContext(c.Request.Context()).Table(table).Where("role_id = ?", roleID).Count(&count).Error; err != nil {
			return false, err
		}
		if count > 0 {
			return true, nil
		}
	}
	return false, nil
}

func normalizeRoleName(c *gin.Context, raw string) (string, bool) {
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

func roleResponseFromModel(role rbac.Role) roleResponse {
	return roleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Code:        role.Code,
		Description: role.Description,
		IsSystem:    role.IsSystem,
		IsActive:    role.IsActive,
		CreatedAt:   role.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:   role.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
}

func writeRoleMutationError(c *gin.Context, err error, fallback string) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		writeError(c, http.StatusNotFound, "role not found")
		return
	}
	if isUniqueConstraintError(err) {
		writeError(c, http.StatusConflict, "role code already exists")
		return
	}
	writeError(c, http.StatusInternalServerError, fallback)
}
