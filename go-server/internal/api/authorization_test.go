package api

import (
	"net/http"
	"testing"
	"time"

	"ai.ro/syncra/dms/internal/rbac"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TestRequireAuthenticatedUserRejectsMissingSession(t *testing.T) {
	_, db := newAuthTestRouter(t)
	router := newAuthorizationTestRouter(db, "organization_unit.create", nil)

	response := authJSON(t, router, http.MethodPost, "/test/permission", `{}`, map[string]string{
		internalAPIHeader: testInternalToken,
	})
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d body=%s, want unauthorized", response.Code, response.Body.String())
	}
}

func TestRequirePermissionRejectsUserWithoutPermission(t *testing.T) {
	authRouter, db := newAuthTestRouter(t)
	router := newAuthorizationTestRouter(db, "organization_unit.create", nil)
	user := createVerifiedUser(t, db, "user@example.com", "password123")
	token := loginUser(t, authRouter, user.Email, "password123")

	response := authJSON(t, router, http.MethodPost, "/test/permission", `{}`, authCookieHeaders(token))
	if response.Code != http.StatusForbidden {
		t.Fatalf("status = %d body=%s, want forbidden", response.Code, response.Body.String())
	}
}

func TestRequirePermissionAllowsSeededLegacyAdmin(t *testing.T) {
	authRouter, db := newAuthTestRouter(t)
	router := newAuthorizationTestRouter(db, "organization_unit.create", nil)
	admin := createAdminUser(t, db, "admin@example.com", "password123")
	if err := rbac.BootstrapLegacyAdmins(db); err != nil {
		t.Fatalf("bootstrap: %v", err)
	}
	token := loginUser(t, authRouter, admin.Email, "password123")

	response := authJSON(t, router, http.MethodPost, "/test/permission", `{}`, authCookieHeaders(token))
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s, want ok", response.Code, response.Body.String())
	}
}

func newAuthorizationTestRouter(db *gorm.DB, permission string, organizationUnitID *string) http.Handler {
	authHandler := newAuthHandler(RouterOptions{
		DB:                  db,
		BetterAuthSecret:    testBetterAuthSecret,
		InternalAPIToken:    testInternalToken,
		AuthSessionTTL:      7 * 24 * time.Hour,
		AuthVerificationTTL: 5 * time.Minute,
	})
	router := gin.New()
	router.Use(authHandler.requireTrustedInternalRequest())
	router.POST("/test/permission", func(c *gin.Context) {
		user, ok := requirePermission(c, authHandler, permission, organizationUnitID)
		if !ok {
			return
		}
		c.JSON(http.StatusOK, gin.H{"userId": user.ID})
	})
	return router
}
