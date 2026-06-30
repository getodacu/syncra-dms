package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
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
