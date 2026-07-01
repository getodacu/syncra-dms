package api

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSwaggerDocDeclaresCurrentRoutes(t *testing.T) {
	content, err := os.ReadFile(filepath.Join("swagger_doc.go"))
	if err != nil {
		t.Fatalf("read swagger_doc.go: %v", err)
	}
	got := string(content)
	for _, want := range []string{
		"swagger:meta",
		"GET /healthz",
		"GET /readyz",
		"GET /version",
		"POST /api/auth/sign-up/email",
		"POST /api/auth/sign-in/email",
		"GET /api/auth/get-session",
		"POST /api/auth/sign-out",
		"POST /api/auth/email-otp/send-verification-otp",
		"POST /api/auth/email-otp/verify-email",
		"POST /api/auth/password-reset/request",
		"POST /api/auth/password-reset/confirm",
		"POST /api/auth/oauth/google/start",
		"POST /api/auth/oauth/google/callback",
		"POST /api/auth/oauth/github/start",
		"POST /api/auth/oauth/github/callback",
		"GET /api/organization-units/tree",
		"GET /api/organization-units/archived",
		"POST /api/organization-units",
		"GET /api/document-folders/tree",
		"POST /api/document-folders",
		"PATCH /api/document-folders/{id}",
		"PATCH /api/document-folders/{id}/parent",
		"POST /api/document-folders/{id}/archive",
		"GET /api/document-folders/{id}/contents",
		"GET /api/users",
		"POST /api/users",
		"GET /api/roles",
		"POST /api/roles",
		"GET /api/permissions",
		"GET /api/permissions/categories",
		"GET /api/groups",
		"POST /api/groups",
		"GET /api/me",
		"GET /api/me/permissions",
		"POST /api/auth/check-permission",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("swagger_doc.go missing %q", want)
		}
	}
}

func TestServedSwaggerDocDeclaresDocumentFolderRoutes(t *testing.T) {
	router := NewRouter(RouterOptions{})
	spec := assertJSONStatus(t, router, "http://localhost:8090/swagger/doc.json", http.StatusOK, "swagger", "2.0")

	paths, ok := spec["paths"].(map[string]any)
	if !ok {
		t.Fatal("/swagger/doc.json missing paths object")
	}

	for _, route := range []struct {
		method string
		path   string
	}{
		{method: "get", path: "/api/document-folders/tree"},
		{method: "post", path: "/api/document-folders"},
		{method: "patch", path: "/api/document-folders/{id}"},
		{method: "patch", path: "/api/document-folders/{id}/parent"},
		{method: "post", path: "/api/document-folders/{id}/archive"},
		{method: "get", path: "/api/document-folders/{id}/contents"},
	} {
		pathDoc, ok := paths[route.path].(map[string]any)
		if !ok {
			t.Fatalf("/swagger/doc.json missing path %q", route.path)
		}
		if _, ok := pathDoc[route.method]; !ok {
			t.Fatalf("/swagger/doc.json missing %s %s", strings.ToUpper(route.method), route.path)
		}
	}
}
