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
		"POST /api/documents/upload",
		"GET /api/documents/{id}",
		"GET /api/documents/{id}/download",
		"PATCH /api/documents/{id}",
		"POST /api/documents/{id}/archive",
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

func TestServedSwaggerDocDeclaresDocumentRoutes(t *testing.T) {
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
		{method: "post", path: "/api/documents/upload"},
		{method: "get", path: "/api/documents/{id}"},
		{method: "get", path: "/api/documents/{id}/download"},
		{method: "patch", path: "/api/documents/{id}"},
		{method: "post", path: "/api/documents/{id}/archive"},
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

func TestServedSwaggerDocDocumentsHiddenDocumentFolderIDs(t *testing.T) {
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
		{method: "patch", path: "/api/document-folders/{id}"},
		{method: "patch", path: "/api/document-folders/{id}/parent"},
		{method: "post", path: "/api/document-folders/{id}/archive"},
		{method: "get", path: "/api/document-folders/{id}/contents"},
	} {
		pathDoc, ok := paths[route.path].(map[string]any)
		if !ok {
			t.Fatalf("/swagger/doc.json missing path %q", route.path)
		}
		operation, ok := pathDoc[route.method].(map[string]any)
		if !ok {
			t.Fatalf("/swagger/doc.json missing %s %s", strings.ToUpper(route.method), route.path)
		}
		responses, ok := operation["responses"].(map[string]any)
		if !ok {
			t.Fatalf("/swagger/doc.json missing responses for %s %s", strings.ToUpper(route.method), route.path)
		}
		if _, ok := responses["403"]; ok {
			t.Fatalf("/swagger/doc.json documents 403 for %s %s, want hidden as 404", strings.ToUpper(route.method), route.path)
		}
		notFound, ok := responses["404"].(map[string]any)
		if !ok {
			t.Fatalf("/swagger/doc.json missing 404 for %s %s", strings.ToUpper(route.method), route.path)
		}
		description, _ := notFound["description"].(string)
		if !strings.Contains(strings.ToLower(description), "inaccessible") {
			t.Fatalf("/swagger/doc.json 404 description for %s %s = %q, want inaccessible behavior documented", strings.ToUpper(route.method), route.path, description)
		}
	}
}

func TestServedSwaggerDocDocumentsDocumentFolderContentsResponse(t *testing.T) {
	router := NewRouter(RouterOptions{})
	spec := assertJSONStatus(t, router, "http://localhost:8090/swagger/doc.json", http.StatusOK, "swagger", "2.0")

	paths, ok := spec["paths"].(map[string]any)
	if !ok {
		t.Fatal("/swagger/doc.json missing paths object")
	}
	pathDoc, ok := paths["/api/document-folders/{id}/contents"].(map[string]any)
	if !ok {
		t.Fatal("/swagger/doc.json missing path /api/document-folders/{id}/contents")
	}
	operation, ok := pathDoc["get"].(map[string]any)
	if !ok {
		t.Fatal("/swagger/doc.json missing GET /api/document-folders/{id}/contents")
	}
	description, _ := operation["description"].(string)
	if strings.Contains(strings.ToLower(description), "not implemented") {
		t.Fatalf("/swagger/doc.json contents description = %q, want implemented listing contract", description)
	}
	responses, ok := operation["responses"].(map[string]any)
	if !ok {
		t.Fatal("/swagger/doc.json missing responses for GET /api/document-folders/{id}/contents")
	}
	okResponse, ok := responses["200"].(map[string]any)
	if !ok {
		t.Fatal("/swagger/doc.json missing 200 for GET /api/document-folders/{id}/contents")
	}
	okDescription, _ := okResponse["description"].(string)
	for _, want := range []string{"folders", "document metadata"} {
		if !strings.Contains(strings.ToLower(okDescription), want) {
			t.Fatalf("/swagger/doc.json 200 description for GET /api/document-folders/{id}/contents = %q, want %q", okDescription, want)
		}
	}
	if _, ok := responses["501"]; ok {
		t.Fatal("/swagger/doc.json documents 501 for implemented folder contents endpoint")
	}
}

func TestServedSwaggerDocDocumentsDocumentFolderUpdateConflicts(t *testing.T) {
	router := NewRouter(RouterOptions{})
	spec := assertJSONStatus(t, router, "http://localhost:8090/swagger/doc.json", http.StatusOK, "swagger", "2.0")

	paths, ok := spec["paths"].(map[string]any)
	if !ok {
		t.Fatal("/swagger/doc.json missing paths object")
	}
	pathDoc, ok := paths["/api/document-folders/{id}"].(map[string]any)
	if !ok {
		t.Fatal("/swagger/doc.json missing path /api/document-folders/{id}")
	}
	operation, ok := pathDoc["patch"].(map[string]any)
	if !ok {
		t.Fatal("/swagger/doc.json missing PATCH /api/document-folders/{id}")
	}
	responses, ok := operation["responses"].(map[string]any)
	if !ok {
		t.Fatal("/swagger/doc.json missing responses for PATCH /api/document-folders/{id}")
	}
	conflict, ok := responses["409"].(map[string]any)
	if !ok {
		t.Fatal("/swagger/doc.json missing 409 for PATCH /api/document-folders/{id}")
	}
	description, _ := conflict["description"].(string)
	normalized := strings.ToLower(description)
	for _, want := range []string{"name", "organization unit"} {
		if !strings.Contains(normalized, want) {
			t.Fatalf("/swagger/doc.json 409 description for PATCH /api/document-folders/{id} = %q, want %q", description, want)
		}
	}
}
