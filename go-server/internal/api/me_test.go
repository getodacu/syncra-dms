package api

import (
	"net/http"
	"strings"
	"testing"

	"ai.ro/syncra/dms/internal/rbac"
)

func TestMePermissionsReturnsEffectivePermissions(t *testing.T) {
	router, db := newAuthTestRouter(t)
	admin := createAdminUser(t, db, "admin@example.com", "password123")
	if err := rbac.BootstrapLegacyAdmins(db); err != nil {
		t.Fatalf("bootstrap: %v", err)
	}
	token := loginUser(t, router, admin.Email, "password123")

	response := authJSON(t, router, http.MethodGet, "/api/me/permissions", "", authCookieHeaders(token))
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), "organization_unit.manage_hierarchy") {
		t.Fatalf("body = %s, want org unit permission", response.Body.String())
	}
}

func TestMeReturnsCurrentUser(t *testing.T) {
	router, db := newAuthTestRouter(t)
	user := createVerifiedUser(t, db, "ada@example.com", "password123")
	token := loginUser(t, router, user.Email, "password123")

	response := authJSON(t, router, http.MethodGet, "/api/me", "", authCookieHeaders(token))
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s, want ok", response.Code, response.Body.String())
	}
	var body userTestResponse
	decodeJSON(t, response, &body)
	if body.ID != user.ID || body.Email != user.Email {
		t.Fatalf("me response = %#v, want current user", body)
	}
}

func TestCheckPermissionReturnsDecision(t *testing.T) {
	router, db := newAuthTestRouter(t)
	adminToken := loginSeededAdmin(t, router, db, "admin@example.com")
	user := createVerifiedUser(t, db, "user@example.com", "password123")
	userToken := loginUser(t, router, user.Email, "password123")

	allowed := authJSON(t, router, http.MethodPost, "/api/auth/check-permission", `{
		"permission":"organization_unit.create"
	}`, authCookieHeaders(adminToken))
	if allowed.Code != http.StatusOK {
		t.Fatalf("allowed status = %d body=%s, want ok", allowed.Code, allowed.Body.String())
	}
	if !strings.Contains(allowed.Body.String(), `"allowed":true`) {
		t.Fatalf("allowed body = %s, want true", allowed.Body.String())
	}

	denied := authJSON(t, router, http.MethodPost, "/api/auth/check-permission", `{
		"permission":"organization_unit.create"
	}`, authCookieHeaders(userToken))
	if denied.Code != http.StatusOK {
		t.Fatalf("denied status = %d body=%s, want ok", denied.Code, denied.Body.String())
	}
	if !strings.Contains(denied.Body.String(), `"allowed":false`) {
		t.Fatalf("denied body = %s, want false", denied.Body.String())
	}
}
