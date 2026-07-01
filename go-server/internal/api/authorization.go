package api

import (
	"net/http"

	"ai.ro/syncra/dms/internal/auth"
	"ai.ro/syncra/dms/internal/rbac"
	"github.com/gin-gonic/gin"
)

func requireAuthenticatedUser(c *gin.Context, h *authHandler) (auth.User, bool) {
	if h == nil {
		writeError(c, http.StatusServiceUnavailable, "authentication is not configured")
		return auth.User{}, false
	}
	if !h.authConfigured(c) {
		return auth.User{}, false
	}
	session, ok, err := h.loadAuthenticatedSession(c)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return auth.User{}, false
	}
	if !ok {
		writeError(c, http.StatusUnauthorized, "authenticated session required")
		return auth.User{}, false
	}
	return session.User, true
}

func requirePermission(c *gin.Context, h *authHandler, permission string, organizationUnitID *string) (auth.User, bool) {
	user, ok := requireAuthenticatedUser(c, h)
	if !ok {
		return auth.User{}, false
	}
	allowed, err := rbac.NewResolver(h.db).Can(c.Request.Context(), rbac.Check{
		UserID:             user.ID,
		Permission:         permission,
		OrganizationUnitID: organizationUnitID,
	})
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to check permission")
		return auth.User{}, false
	}
	if !allowed {
		writeError(c, http.StatusForbidden, "permission required")
		return auth.User{}, false
	}
	return user, true
}

func requireAnyPermission(c *gin.Context, h *authHandler, permissions []string, organizationUnitID *string) (auth.User, bool) {
	user, ok := requireAuthenticatedUser(c, h)
	if !ok {
		return auth.User{}, false
	}
	resolver := rbac.NewResolver(h.db)
	for _, permission := range permissions {
		allowed, err := resolver.Can(c.Request.Context(), rbac.Check{
			UserID:             user.ID,
			Permission:         permission,
			OrganizationUnitID: organizationUnitID,
		})
		if err != nil {
			writeError(c, http.StatusInternalServerError, "failed to check permission")
			return auth.User{}, false
		}
		if allowed {
			return user, true
		}
	}
	writeError(c, http.StatusForbidden, "permission required")
	return auth.User{}, false
}
