package api

import (
	"net/http"
	"time"

	"ai.ro/syncra/dms/internal/rbac"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type permissionHandler struct {
	db   *gorm.DB
	auth *authHandler
}

type permissionResponse struct {
	ID          string  `json:"id"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Category    string  `json:"category"`
	IsSystem    bool    `json:"isSystem"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

type permissionListResponse struct {
	Permissions []permissionResponse `json:"permissions"`
}

type permissionCategoriesResponse struct {
	Categories []string `json:"categories"`
}

func newPermissionHandler(options RouterOptions, auth *authHandler) *permissionHandler {
	return &permissionHandler{db: options.DB, auth: auth}
}

func (h *permissionHandler) list(c *gin.Context) {
	if _, ok := requirePermission(c, h.auth, "role.view", nil); !ok {
		return
	}
	var permissions []rbac.Permission
	if err := h.db.WithContext(c.Request.Context()).
		Order("category asc, code asc").
		Find(&permissions).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to list permissions")
		return
	}
	out := make([]permissionResponse, 0, len(permissions))
	for _, permission := range permissions {
		out = append(out, permissionResponseFromModel(permission))
	}
	c.JSON(http.StatusOK, permissionListResponse{Permissions: out})
}

func (h *permissionHandler) categories(c *gin.Context) {
	if _, ok := requirePermission(c, h.auth, "role.view", nil); !ok {
		return
	}
	var categories []string
	if err := h.db.WithContext(c.Request.Context()).
		Model(&rbac.Permission{}).
		Distinct("category").
		Order("category asc").
		Pluck("category", &categories).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to list permission categories")
		return
	}
	c.JSON(http.StatusOK, permissionCategoriesResponse{Categories: categories})
}

func permissionResponseFromModel(permission rbac.Permission) permissionResponse {
	return permissionResponse{
		ID:          permission.ID,
		Code:        permission.Code,
		Name:        permission.Name,
		Description: permission.Description,
		Category:    permission.Category,
		IsSystem:    permission.IsSystem,
		CreatedAt:   permission.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:   permission.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
}
