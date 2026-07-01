package api

import (
	"net/http"
	"testing"
)

func TestPermissionAPIRequiresRoleViewPermission(t *testing.T) {
	router, db := newAuthTestRouter(t)
	user := createVerifiedUser(t, db, "user@example.com", "password123")
	userToken := loginUser(t, router, user.Email, "password123")

	forbidden := authJSON(t, router, http.MethodGet, "/api/permissions", "", authCookieHeaders(userToken))
	if forbidden.Code != http.StatusForbidden {
		t.Fatalf("status = %d body=%s, want forbidden", forbidden.Code, forbidden.Body.String())
	}

	adminToken := loginSeededAdmin(t, router, db, "admin@example.com")
	allowed := authJSON(t, router, http.MethodGet, "/api/permissions", "", authCookieHeaders(adminToken))
	if allowed.Code != http.StatusOK {
		t.Fatalf("admin status = %d body=%s, want ok", allowed.Code, allowed.Body.String())
	}
	var body permissionListTestResponse
	decodeJSON(t, allowed, &body)
	if len(body.Permissions) == 0 {
		t.Fatal("permissions list was empty")
	}
}

func TestPermissionAPICategories(t *testing.T) {
	router, db := newAuthTestRouter(t)
	adminToken := loginSeededAdmin(t, router, db, "admin@example.com")

	response := authJSON(t, router, http.MethodGet, "/api/permissions/categories", "", authCookieHeaders(adminToken))
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s, want ok", response.Code, response.Body.String())
	}
	var body struct {
		Categories []string `json:"categories"`
	}
	decodeJSON(t, response, &body)
	if len(body.Categories) == 0 {
		t.Fatal("categories list was empty")
	}
	seen := map[string]bool{}
	for _, category := range body.Categories {
		if category == "" {
			t.Fatal("category was empty")
		}
		if seen[category] {
			t.Fatalf("duplicate category %q in %#v", category, body.Categories)
		}
		seen[category] = true
	}
}

type permissionListTestResponse struct {
	Permissions []permissionTestResponse `json:"permissions"`
}

type permissionTestResponse struct {
	ID       string `json:"id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	Category string `json:"category"`
}
