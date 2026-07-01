package api

import (
	"net/http"
	"strings"

	"ai.ro/syncra/dms/internal/rbac"
	"github.com/gin-gonic/gin"
)

type meHandler struct {
	auth *authHandler
}

type mePermissionResponse struct {
	Code               string  `json:"code"`
	ScopeType          string  `json:"scopeType"`
	OrganizationUnitID *string `json:"organizationUnitId"`
	Source             string  `json:"source"`
}

type mePermissionsResponse struct {
	Permissions []mePermissionResponse `json:"permissions"`
}

type checkPermissionRequest struct {
	Permission         string  `json:"permission"`
	OrganizationUnitID *string `json:"organizationUnitId"`
}

type checkPermissionResponse struct {
	Allowed bool `json:"allowed"`
}

func newMeHandler(_ RouterOptions, auth *authHandler) *meHandler {
	return &meHandler{auth: auth}
}

func (h *meHandler) getMe(c *gin.Context) {
	user, ok := requireAuthenticatedUser(c, h.auth)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, userResponseFromModel(user))
}

func (h *meHandler) getPermissions(c *gin.Context) {
	user, ok := requireAuthenticatedUser(c, h.auth)
	if !ok {
		return
	}
	grants, err := rbac.NewResolver(h.auth.db).EffectiveGrants(c.Request.Context(), user.ID)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to load permissions")
		return
	}
	out := make([]mePermissionResponse, 0, len(grants))
	for _, grant := range grants {
		out = append(out, mePermissionResponseFromGrant(grant))
	}
	c.JSON(http.StatusOK, mePermissionsResponse{Permissions: out})
}

func (h *meHandler) checkPermission(c *gin.Context) {
	user, ok := requireAuthenticatedUser(c, h.auth)
	if !ok {
		return
	}
	var req checkPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	permission := strings.TrimSpace(req.Permission)
	if permission == "" {
		writeError(c, http.StatusBadRequest, "permission is required")
		return
	}
	var organizationUnitID *string
	if req.OrganizationUnitID != nil && strings.TrimSpace(*req.OrganizationUnitID) != "" {
		id, ok := parseUUIDValue(c, *req.OrganizationUnitID, "invalid organization unit id")
		if !ok {
			return
		}
		organizationUnitID = &id
	}
	allowed, err := rbac.NewResolver(h.auth.db).Can(c.Request.Context(), rbac.Check{
		UserID:             user.ID,
		Permission:         permission,
		OrganizationUnitID: organizationUnitID,
	})
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to check permission")
		return
	}
	c.JSON(http.StatusOK, checkPermissionResponse{Allowed: allowed})
}

func mePermissionResponseFromGrant(grant rbac.Grant) mePermissionResponse {
	return mePermissionResponse{
		Code:               grant.PermissionCode,
		ScopeType:          string(grant.ScopeType),
		OrganizationUnitID: grant.OrganizationUnitID,
		Source:             grant.Source,
	}
}
