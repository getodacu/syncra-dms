package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRouterExposesHealthVersionAndReadiness(t *testing.T) {
	router := NewRouter(RouterOptions{
		Version: VersionInfo{
			AppName: "Syncra DMS",
			Module:  "ai.ro/syncra/dms",
			Version: "test",
		},
		Ready: func(context.Context) error {
			return nil
		},
	})

	assertJSONStatus(t, router, "/healthz", http.StatusOK, "status", "ok")
	assertJSONStatus(t, router, "/readyz", http.StatusOK, "status", "ready")
	assertJSONStatus(t, router, "/version", http.StatusOK, "module", "ai.ro/syncra/dms")
}

func TestRouterReportsReadinessFailure(t *testing.T) {
	router := NewRouter(RouterOptions{
		Ready: func(context.Context) error {
			return errors.New("database unavailable")
		},
	})

	body := assertJSONStatus(t, router, "/readyz", http.StatusServiceUnavailable, "status", "not_ready")
	if body["error"] != "database unavailable" {
		t.Fatalf("readyz error = %v, want database unavailable", body["error"])
	}
}

func TestRouterExposesSwaggerDocs(t *testing.T) {
	router := NewRouter(RouterOptions{})

	spec := assertJSONStatus(t, router, "http://localhost:8090/swagger/doc.json", http.StatusOK, "swagger", "2.0")
	if spec["info"] == nil {
		t.Fatal("/swagger/doc.json missing info section")
	}
	if spec["host"] != "localhost:8090" {
		t.Fatalf("/swagger/doc.json host = %v, want localhost:8090", spec["host"])
	}

	request := httptest.NewRequest(http.MethodGet, "/swagger/index.html", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("/swagger/index.html status = %d, want %d; body = %s", response.Code, http.StatusOK, response.Body.String())
	}
	if contentType := response.Header().Get("Content-Type"); !strings.HasPrefix(contentType, "text/html") {
		t.Fatalf("/swagger/index.html Content-Type = %q, want text/html", contentType)
	}
	body := response.Body.String()
	for _, want := range []string{"SwaggerUIBundle", "/swagger/doc.json"} {
		if !strings.Contains(body, want) {
			t.Fatalf("/swagger/index.html missing %q", want)
		}
	}
}

func assertJSONStatus(t *testing.T, handler http.Handler, path string, wantStatus int, key string, wantValue any) map[string]any {
	t.Helper()

	request := httptest.NewRequest(http.MethodGet, path, nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != wantStatus {
		t.Fatalf("%s status = %d, want %d; body = %s", path, response.Code, wantStatus, response.Body.String())
	}

	var body map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("%s response is not JSON: %v", path, err)
	}
	if body[key] != wantValue {
		t.Fatalf("%s %q = %v, want %v", path, key, body[key], wantValue)
	}
	return body
}
