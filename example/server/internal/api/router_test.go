package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"ai.ro/syncra/docs"
	"ai.ro/syncra/internal/logging"

	"github.com/gin-gonic/gin"
)

func TestNewRouterServesSwaggerUI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter(&Handler{})

	req := httptest.NewRequest(http.MethodGet, "/swagger/index.html", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestNewRouterServesSwaggerSpec(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter(&Handler{})

	req := httptest.NewRequest(http.MethodGet, "/swagger/doc.json", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if got := w.Header().Get("Content-Type"); !strings.Contains(got, "application/json") {
		t.Fatalf("Content-Type = %q, want application/json", got)
	}
	var spec struct {
		Paths       map[string]json.RawMessage `json:"paths"`
		Definitions map[string]json.RawMessage `json:"definitions"`
	}
	if err := json.NewDecoder(w.Body).Decode(&spec); err != nil {
		t.Fatalf("decode swagger spec: %v", err)
	}
	expectedRoutes := documentedRouterRoutes(router, internalAPIDocRoute)
	specRoutes := swaggerSpecRoutes(spec.Paths)
	if !reflect.DeepEqual(specRoutes, expectedRoutes) {
		t.Fatalf("swagger spec routes = %#v, want %#v", specRoutes, expectedRoutes)
	}
	assertGetSessionResponseNullable(t, spec.Paths)
	assertCollectionListNextCursorNullable(t, spec.Definitions)
	assertSwaggerOperationTags(t, spec.Paths, internalAPIDocRoute)
	assertInternalAPISwaggerContract(t, spec.Paths)
	assertOCRJobSwaggerContract(t, spec.Paths)
	assertBillingSwaggerContract(t, spec.Paths, spec.Definitions)
	assertDashboardSwaggerContract(t, spec.Paths, spec.Definitions)
	assertSwaggerDefinitionMissing(t, spec.Definitions, "publicOCRJobResponse")
}

func TestNewRouterServesPublicSwaggerUI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter(&Handler{})

	req := httptest.NewRequest(http.MethodGet, "/swagger-public/index.html", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestNewRouterServesPublicSwaggerSpec(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter(&Handler{})

	req := httptest.NewRequest(http.MethodGet, "/swagger-public/doc.json", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if got := w.Header().Get("Content-Type"); !strings.Contains(got, "application/json") {
		t.Fatalf("Content-Type = %q, want application/json", got)
	}
	var spec struct {
		Paths       map[string]json.RawMessage `json:"paths"`
		Definitions map[string]json.RawMessage `json:"definitions"`
	}
	if err := json.NewDecoder(w.Body).Decode(&spec); err != nil {
		t.Fatalf("decode public swagger spec: %v", err)
	}
	expectedRoutes := documentedRouterRoutes(router, publicAPIDocRoute)
	specRoutes := swaggerSpecRoutes(spec.Paths)
	if !reflect.DeepEqual(specRoutes, expectedRoutes) {
		t.Fatalf("public swagger spec routes = %#v, want %#v", specRoutes, expectedRoutes)
	}
	assertPublicOCRJobSwaggerContract(t, spec.Paths, spec.Definitions)
	assertPublicBalanceSwaggerContract(t, spec.Paths, spec.Definitions)
	assertSwaggerOperationTags(t, spec.Paths, publicAPIDocRoute)
	assertSwaggerDefinitionMissing(t, spec.Definitions, "collectionResponse")
}

func TestNewRouterHandlesPublicAPICORSPreflight(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter(&Handler{})

	req := httptest.NewRequest(http.MethodOptions, "/v1/get-balance", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("Access-Control-Request-Method", http.MethodGet)
	req.Header.Set("Access-Control-Request-Headers", "authorization")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Fatalf("Access-Control-Allow-Origin = %q, want *", got)
	}
	if got := w.Header().Get("Access-Control-Allow-Methods"); !strings.Contains(got, http.MethodGet) {
		t.Fatalf("Access-Control-Allow-Methods = %q, want GET", got)
	}
	if got := w.Header().Get("Access-Control-Allow-Headers"); !strings.Contains(strings.ToLower(got), "authorization") {
		t.Fatalf("Access-Control-Allow-Headers = %q, want Authorization", got)
	}
}

func TestNewRouterAddsPublicAPICORSHeadersToResponses(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter(&Handler{})

	req := httptest.NewRequest(http.MethodGet, "/v1/get-balance", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Fatalf("Access-Control-Allow-Origin = %q, want *", got)
	}
}

func TestNewRouterLogsStructuredRequestMetadata(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var logs bytes.Buffer
	router := NewRouter(&Handler{Logger: logging.NewJSONLogger(&logs, true)})

	req := httptest.NewRequest(http.MethodGet, "/swagger/doc.json", nil)
	req.Header.Set(requestIDHeader, "request-123")
	req.Header.Set("Authorization", "Bearer secret-token")
	req.Header.Set("Cookie", "auth.session_token=secret-cookie")
	req.Header.Set("User-Agent", "test-client")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if got := w.Header().Get(requestIDHeader); got != "request-123" {
		t.Fatalf("request id header = %q, want request-123", got)
	}
	out := logs.String()
	for _, want := range []string{
		`"msg":"http.request_completed"`,
		`"request_id":"request-123"`,
		`"route":"/swagger/doc.json"`,
		`"status":200`,
		`"user_agent_class":"api-client"`,
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("logs missing %s in:\n%s", want, out)
		}
	}
	for _, forbidden := range []string{"secret-token", "secret-cookie", "Authorization", "Cookie"} {
		if strings.Contains(out, forbidden) {
			t.Fatalf("logs leaked %q in:\n%s", forbidden, out)
		}
	}
}

func TestWriteErrorLogsStructuredErrorResponseWithoutSecrets(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var logs bytes.Buffer
	router := NewRouter(&Handler{
		InternalAPIToken: "trusted-internal-token",
		Logger:           logging.NewJSONLogger(&logs, true),
	})

	req := httptest.NewRequest(http.MethodGet, "/api/ocr/jobs/not-a-uuid", nil)
	req.Header.Set(requestIDHeader, "error-request-123")
	req.Header.Set("Authorization", "Bearer secret-api-key")
	req.Header.Set(internalAPIHeader, "wrong-secret-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	out := logs.String()
	for _, want := range []string{
		`"msg":"http.error_response"`,
		`"request_id":"error-request-123"`,
		`"route":"/api/ocr/jobs/:id"`,
		`"status":401`,
		`"message":"unauthorized"`,
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("logs missing %s in:\n%s", want, out)
		}
	}
	for _, forbidden := range []string{"secret-api-key", "wrong-secret-token", "Authorization", internalAPIHeader} {
		if strings.Contains(out, forbidden) {
			t.Fatalf("logs leaked %q:\n%s", forbidden, out)
		}
	}
}

func TestNewRouterProtectsInternalAPIRoutesWithInternalToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter(&Handler{InternalAPIToken: "trusted-internal-token"})

	for _, tt := range []struct {
		name   string
		method string
		path   string
	}{
		{name: "auth", method: http.MethodGet, path: "/api/auth/get-session"},
		{name: "json recipes", method: http.MethodGet, path: "/api/json-recipes"},
		{name: "ocr", method: http.MethodGet, path: "/api/ocr/jobs/not-a-uuid"},
		{name: "collections", method: http.MethodGet, path: "/api/collections/not-a-uuid"},
		{name: "datasets", method: http.MethodGet, path: "/api/datasets/not-a-uuid"},
		{name: "billing", method: http.MethodGet, path: "/api/billing/balance?user_id=not-a-uuid"},
		{name: "dashboard", method: http.MethodGet, path: "/api/dashboard/summary?user_id=not-a-uuid"},
	} {
		t.Run(tt.name+" missing token", func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusUnauthorized {
				t.Fatalf("status = %d body=%s, want 401", w.Code, w.Body.String())
			}
			if !strings.Contains(w.Body.String(), "unauthorized") {
				t.Fatalf("body = %s, want unauthorized", w.Body.String())
			}
		})

		t.Run(tt.name+" wrong token", func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			req.Header.Set(internalAPIHeader, "wrong-token")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusUnauthorized {
				t.Fatalf("status = %d body=%s, want 401", w.Code, w.Body.String())
			}
			if !strings.Contains(w.Body.String(), "unauthorized") {
				t.Fatalf("body = %s, want unauthorized", w.Body.String())
			}
		})

		t.Run(tt.name+" accepted token reaches handler", func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			req.Header.Set(internalAPIHeader, "trusted-internal-token")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code == http.StatusUnauthorized {
				t.Fatalf("status = %d body=%s, want request past internal-token gate", w.Code, w.Body.String())
			}
		})
	}
}

func TestNewRouterLeavesPublicAndSwaggerRoutesOutsideInternalTokenGate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter(&Handler{InternalAPIToken: "trusted-internal-token"})

	for _, tt := range []struct {
		name       string
		path       string
		wantStatus int
	}{
		{name: "public api", path: "/v1/get-balance", wantStatus: http.StatusUnauthorized},
		{name: "internal swagger spec", path: "/swagger/doc.json", wantStatus: http.StatusOK},
		{name: "public swagger spec", path: "/swagger-public/doc.json", wantStatus: http.StatusOK},
	} {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Fatalf("status = %d body=%s, want %d", w.Code, w.Body.String(), tt.wantStatus)
			}
		})
	}
}

func TestSwaggerStaticSpecDocumentsRouterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter(&Handler{})

	var spec struct {
		Paths       map[string]json.RawMessage `json:"paths"`
		Definitions map[string]json.RawMessage `json:"definitions"`
	}
	if err := json.Unmarshal(docs.SwaggerJSON, &spec); err != nil {
		t.Fatalf("decode static swagger spec: %v", err)
	}
	expectedRoutes := documentedRouterRoutes(router, allAPIDocRoute)
	specRoutes := swaggerSpecRoutes(spec.Paths)
	if !reflect.DeepEqual(specRoutes, expectedRoutes) {
		t.Fatalf("static swagger spec routes = %#v, want %#v", specRoutes, expectedRoutes)
	}
	assertPublicOCRJobSwaggerContract(t, spec.Paths, spec.Definitions)
	assertPublicBalanceSwaggerContract(t, spec.Paths, spec.Definitions)
	assertBillingSwaggerContract(t, spec.Paths, spec.Definitions)
	assertDashboardSwaggerContract(t, spec.Paths, spec.Definitions)
	assertInternalAPISwaggerContract(t, spec.Paths)
	assertSwaggerOperationTags(t, spec.Paths, allAPIDocRoute)
}

func TestNewRouterDoesNotServeLegacyUnprefixedAPIRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter(&Handler{})

	for _, tt := range []struct {
		method string
		path   string
	}{
		{method: http.MethodPost, path: "/schemas"},
		{method: http.MethodGet, path: "/schemas"},
		{method: http.MethodGet, path: "/schemas/schema-1"},
		{method: http.MethodPost, path: "/ocr"},
		{method: http.MethodGet, path: "/ocr/documents"},
		{method: http.MethodGet, path: "/ocr/document/document-1"},
		{method: http.MethodPost, path: "/ocr/jobs"},
		{method: http.MethodDelete, path: "/ocr/jobs"},
		{method: http.MethodGet, path: "/ocr/jobs/job-1"},
	} {
		req := httptest.NewRequest(tt.method, tt.path, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("%s %s status = %d body=%s", tt.method, tt.path, w.Code, w.Body.String())
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/swagger/doc.json", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("/swagger/doc.json status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestNewRouterDoesNotRouteNestedSwaggerDocsToUI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter(&Handler{})

	for _, path := range []string{
		"/swagger/nested/doc.json",
		"/swagger/nested/doc.json/",
		"/swagger/nested/index.html",
		"/swagger-public/nested/doc.json",
		"/swagger-public/nested/doc.json/",
		"/swagger-public/nested/index.html",
	} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("%s status = %d body=%s", path, w.Code, w.Body.String())
		}
	}
}

func TestNewRouterAppliesRuntimeSwaggerHostAndSchemes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter(&Handler{
		SwaggerHost:    "api.example.test",
		SwaggerSchemes: []string{"https"},
	})

	req := httptest.NewRequest(http.MethodGet, "/swagger/doc.json", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var spec struct {
		Host    string   `json:"host"`
		Schemes []string `json:"schemes"`
	}
	if err := json.NewDecoder(w.Body).Decode(&spec); err != nil {
		t.Fatalf("decode swagger spec: %v", err)
	}
	if spec.Host != "api.example.test" {
		t.Fatalf("host = %q, want api.example.test", spec.Host)
	}
	if !reflect.DeepEqual(spec.Schemes, []string{"https"}) {
		t.Fatalf("schemes = %#v, want https", spec.Schemes)
	}
}

func TestNewRouterAppliesRuntimePublicSwaggerHostAndSchemes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter(&Handler{
		SwaggerHost:    "api.example.test",
		SwaggerSchemes: []string{"https"},
	})

	req := httptest.NewRequest(http.MethodGet, "/swagger-public/doc.json", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var spec struct {
		Host    string   `json:"host"`
		Schemes []string `json:"schemes"`
	}
	if err := json.NewDecoder(w.Body).Decode(&spec); err != nil {
		t.Fatalf("decode public swagger spec: %v", err)
	}
	if spec.Host != "api.example.test" {
		t.Fatalf("host = %q, want api.example.test", spec.Host)
	}
	if !reflect.DeepEqual(spec.Schemes, []string{"https"}) {
		t.Fatalf("schemes = %#v, want https", spec.Schemes)
	}
}

var ginParamPattern = regexp.MustCompile(`:([^/]+)`)

type docRoutePredicate func(path string) bool

func internalAPIDocRoute(path string) bool {
	return !strings.HasPrefix(path, "/v1/")
}

func publicAPIDocRoute(path string) bool {
	return strings.HasPrefix(path, "/v1/")
}

func allAPIDocRoute(path string) bool {
	return internalAPIDocRoute(path) || publicAPIDocRoute(path)
}

func documentedRouterRoutes(router *gin.Engine, include docRoutePredicate) map[string]map[string]struct{} {
	routes := make(map[string]map[string]struct{})
	for _, route := range router.Routes() {
		if strings.HasPrefix(route.Path, "/swagger/") || strings.HasPrefix(route.Path, "/swagger-public/") {
			continue
		}
		if route.Method == http.MethodOptions && route.Path == "/v1/*path" {
			continue
		}
		if !include(route.Path) {
			continue
		}
		path := ginParamPattern.ReplaceAllString(route.Path, `{$1}`)
		methods := routes[path]
		if methods == nil {
			methods = make(map[string]struct{})
			routes[path] = methods
		}
		methods[strings.ToLower(route.Method)] = struct{}{}
	}
	return routes
}

func swaggerSpecRoutes(paths map[string]json.RawMessage) map[string]map[string]struct{} {
	routes := make(map[string]map[string]struct{})
	for path, rawOperations := range paths {
		var operations map[string]json.RawMessage
		if err := json.Unmarshal(rawOperations, &operations); err != nil {
			routes[path] = map[string]struct{}{"<invalid>": {}}
			continue
		}
		methods := make(map[string]struct{})
		for method := range operations {
			methods[strings.ToLower(method)] = struct{}{}
		}
		routes[path] = methods
	}
	return routes
}

func assertGetSessionResponseNullable(t *testing.T, paths map[string]json.RawMessage) {
	t.Helper()

	var path map[string]struct {
		Responses map[string]struct {
			Schema map[string]any `json:"schema"`
		} `json:"responses"`
	}
	if err := json.Unmarshal(paths["/api/auth/get-session"], &path); err != nil {
		t.Fatalf("decode get-session path: %v", err)
	}
	schema := path["get"].Responses["200"].Schema
	if schema["x-nullable"] != true {
		t.Fatalf("get-session 200 schema x-nullable = %#v, want true", schema["x-nullable"])
	}
}

func assertCollectionListNextCursorNullable(t *testing.T, definitions map[string]json.RawMessage) {
	t.Helper()

	var definition struct {
		Properties map[string]map[string]any `json:"properties"`
	}
	raw, ok := definitions["collectionListResponse"]
	if !ok {
		t.Fatal("collectionListResponse definition missing")
	}
	if err := json.Unmarshal(raw, &definition); err != nil {
		t.Fatalf("decode collectionListResponse definition: %v", err)
	}
	nextCursor := definition.Properties["next_cursor"]
	if nextCursor["x-nullable"] != true {
		t.Fatalf("collectionListResponse.next_cursor x-nullable = %#v, want true", nextCursor["x-nullable"])
	}
}

func assertSwaggerOperationTags(t *testing.T, paths map[string]json.RawMessage, include docRoutePredicate) {
	t.Helper()

	expected := map[string]map[string]string{
		"/api/auth/apikeys":                               {"post": "auth", "delete": "auth"},
		"/api/auth/apikeys/{user_id}":                     {"get": "auth"},
		"/api/auth/email-otp/send-verification-otp":       {"post": "auth"},
		"/api/auth/email-otp/verify-email":                {"post": "auth"},
		"/api/auth/get-session":                           {"get": "auth"},
		"/api/auth/oauth/github/callback":                 {"post": "auth"},
		"/api/auth/oauth/github/start":                    {"post": "auth"},
		"/api/auth/oauth/google/callback":                 {"post": "auth"},
		"/api/auth/oauth/google/start":                    {"post": "auth"},
		"/api/auth/sign-in/email":                         {"post": "auth"},
		"/api/auth/sign-out":                              {"post": "auth"},
		"/api/auth/sign-up/email":                         {"post": "auth"},
		"/api/auth/user":                                  {"patch": "auth"},
		"/api/auth/webhook":                               {"post": "auth", "delete": "auth"},
		"/api/auth/webhook/{user_id}":                     {"get": "auth"},
		"/api/auth/webhook/{user_id}/secret":              {"patch": "auth"},
		"/api/admin/impersonation/stop":                   {"post": "admin"},
		"/api/admin/billing/invoices":                     {"get": "admin"},
		"/api/admin/billing/orders":                       {"get": "admin"},
		"/api/admin/json-recipe-categories":               {"get": "json-recipes", "post": "json-recipes"},
		"/api/admin/json-recipe-categories/{id}":          {"get": "json-recipes", "put": "json-recipes", "delete": "json-recipes"},
		"/api/admin/json-recipes":                         {"get": "json-recipes", "post": "json-recipes"},
		"/api/admin/json-recipes/{id}":                    {"get": "json-recipes", "put": "json-recipes", "delete": "json-recipes"},
		"/api/admin/users":                                {"get": "admin"},
		"/api/admin/users/{id}":                           {"get": "admin", "patch": "admin"},
		"/api/admin/users/{id}/billing-profile":           {"put": "admin"},
		"/api/admin/users/{id}/impersonation":             {"post": "admin"},
		"/api/admin/users/{id}/password":                  {"post": "admin"},
		"/api/billing/balance":                            {"get": "billing"},
		"/api/billing/profile":                            {"get": "billing", "put": "billing"},
		"/api/billing/credit-usage-history":               {"get": "billing"},
		"/api/billing/invoices/{id}/email-delivery/claim": {"post": "billing"},
		"/api/billing/invoices/{id}/email-delivery/sent":  {"post": "billing"},
		"/api/billing/orders":                             {"get": "billing", "post": "billing"},
		"/api/billing/orders/{id}/checkout-session":       {"post": "billing"},
		"/api/billing/orders/{id}/failed":                 {"post": "billing"},
		"/api/billing/orders/{id}/paid":                   {"post": "billing"},
		"/api/dashboard/summary":                          {"get": "dashboard"},
		"/api/collection":                                 {"post": "collections"},
		"/api/collection/{id}":                            {"put": "collections", "delete": "collections"},
		"/api/collections":                                {"get": "collections"},
		"/api/collections/{id}":                           {"get": "collections"},
		"/api/json-recipes":                               {"get": "json-recipes"},
		"/api/json-recipes/{id}/deploy":                   {"post": "json-recipes"},
		"/api/ocr":                                        {"post": "ocr"},
		"/api/ocr/document/{id}":                          {"get": "documents"},
		"/api/ocr/documents":                              {"get": "documents", "delete": "documents"},
		"/api/ocr/documents/{id}":                         {"patch": "documents", "delete": "documents"},
		"/api/ocr/jobs":                                   {"post": "jobs", "get": "jobs", "delete": "jobs"},
		"/api/ocr/jobs/{id}":                              {"get": "jobs"},
		"/api/ocr/schemas":                                {"post": "schemas", "get": "schemas", "delete": "schemas"},
		"/api/ocr/schemas/{id}":                           {"get": "schemas", "put": "schemas", "delete": "schemas"},
		"/v1/get-balance":                                 {"get": "public-billing"},
		"/v1/ocr/jobs":                                    {"post": "public-ocr"},
		"/v1/ocr/jobs/{id}":                               {"get": "public-ocr"},
	}

	for path, methods := range expected {
		if !include(path) {
			continue
		}
		var operations map[string]struct {
			Tags []string `json:"tags"`
		}
		if err := json.Unmarshal(paths[path], &operations); err != nil {
			t.Fatalf("decode swagger operations for %s: %v", path, err)
		}
		for method, want := range methods {
			got := operations[method].Tags
			if !reflect.DeepEqual(got, []string{want}) {
				t.Fatalf("%s %s tags = %#v, want [%q]", method, path, got, want)
			}
		}
	}
}

func assertBillingSwaggerContract(t *testing.T, paths map[string]json.RawMessage, definitions map[string]json.RawMessage) {
	t.Helper()

	for _, check := range []struct {
		path   string
		method string
	}{
		{path: "/api/billing/profile", method: "get"},
		{path: "/api/billing/profile", method: "put"},
		{path: "/api/billing/credit-usage-history", method: "get"},
		{path: "/api/billing/invoices/{id}/email-delivery/claim", method: "post"},
		{path: "/api/billing/invoices/{id}/email-delivery/sent", method: "post"},
		{path: "/api/billing/orders", method: "get"},
		{path: "/api/billing/orders", method: "post"},
		{path: "/api/billing/orders/{id}/checkout-session", method: "post"},
		{path: "/api/billing/orders/{id}/paid", method: "post"},
		{path: "/api/billing/orders/{id}/failed", method: "post"},
	} {
		operation := swaggerOperation(t, paths, check.path, check.method)
		if !swaggerOperationHasHeader(operation, internalAPIHeader) {
			t.Fatalf("%s %s missing %s header parameter", check.method, check.path, internalAPIHeader)
		}
	}

	paid := swaggerOperation(t, paths, "/api/billing/orders/{id}/paid", "post")
	if _, ok := paid.Responses["409"]; !ok {
		t.Fatal("billing paid operation missing 409 response")
	}

	assertSwaggerDefinitionProperty(t, definitions, "creditBalanceResponse", "user_id", "format", "uuid")
	assertSwaggerDefinitionProperty(t, definitions, "billingProfileResponse", "user_id", "format", "uuid")
	assertSwaggerDefinitionProperty(t, definitions, "billingProfileEnvelopeResponse", "profile", "x-nullable", true)
	assertSwaggerDefinitionProperty(t, definitions, "billingOrderResponse", "user_id", "format", "uuid")
	assertSwaggerDefinitionProperty(t, definitions, "billingOrderListResponse", "next_cursor", "x-nullable", true)
	assertSwaggerDefinitionProperty(t, definitions, "markBillingOrderPaidRequest", "paid_at", "format", "date-time")
}

func assertDashboardSwaggerContract(t *testing.T, paths map[string]json.RawMessage, definitions map[string]json.RawMessage) {
	t.Helper()

	getSummary := swaggerOperation(t, paths, "/api/dashboard/summary", "get")
	for _, status := range []string{"200", "400", "401", "500"} {
		if _, ok := getSummary.Responses[status]; !ok {
			t.Fatalf("dashboard summary operation missing %s response", status)
		}
	}
	userID := swaggerOperationParameter(t, getSummary, "query", "user_id")
	if !userID.Required {
		t.Fatal("dashboard summary user_id parameter required = false, want true")
	}
	if userID.Format != "uuid" {
		t.Fatalf("dashboard summary user_id format = %q, want uuid", userID.Format)
	}
	rangeParam := swaggerOperationParameter(t, getSummary, "query", "range")
	if !reflect.DeepEqual(rangeParam.Enum, []string{"7d", "30d", "90d"}) {
		t.Fatalf("dashboard summary range enum = %#v, want [7d 30d 90d]", rangeParam.Enum)
	}

	assertSwaggerDefinitionProperty(t, definitions, "dashboardRecentDocumentResponse", "id", "format", "uuid")
	assertSwaggerDefinitionProperty(t, definitions, "dashboardRecentDocumentResponse", "schema_id", "x-nullable", true)
	assertSwaggerDefinitionProperty(t, definitions, "dashboardRecentDocumentResponse", "schema_id", "format", "uuid")
	assertSwaggerDefinitionProperty(t, definitions, "dashboardRecentDocumentResponse", "schema_name", "x-nullable", true)
	assertSwaggerDefinitionProperty(t, definitions, "dashboardRecentDatasetResponse", "id", "format", "uuid")
	assertSwaggerDefinitionProperty(t, definitions, "dashboardSchemaThroughputResponse", "schema_id", "x-nullable", true)
	assertSwaggerDefinitionProperty(t, definitions, "dashboardSchemaThroughputResponse", "schema_id", "format", "uuid")
}

func assertInternalAPISwaggerContract(t *testing.T, paths map[string]json.RawMessage) {
	t.Helper()

	for path, raw := range paths {
		if !strings.HasPrefix(path, "/api/") {
			continue
		}

		var operations map[string]swaggerOperationDoc
		if err := json.Unmarshal(raw, &operations); err != nil {
			t.Fatalf("decode swagger operations for %s: %v", path, err)
		}

		for method, operation := range operations {
			header := swaggerOperationParameter(t, operation, "header", internalAPIHeader)
			if !header.Required {
				t.Fatalf("%s %s %s required = false, want true", method, path, internalAPIHeader)
			}
			if header.Type != "string" {
				t.Fatalf("%s %s %s type = %q, want string", method, path, internalAPIHeader, header.Type)
			}
		}
	}
}

func assertOCRJobSwaggerContract(t *testing.T, paths map[string]json.RawMessage) {
	t.Helper()

	createJob := swaggerOperation(t, paths, "/api/ocr/jobs", "post")
	userID := swaggerOperationParameter(t, createJob, "formData", "user_id")
	if !userID.Required {
		t.Fatal("create OCR job user_id parameter required = false, want true")
	}
	if userID.Format != "uuid" {
		t.Fatalf("create OCR job user_id format = %q, want uuid", userID.Format)
	}
	if userID.Nullable {
		t.Fatal("create OCR job user_id x-nullable = true, want false")
	}
}

func assertPublicOCRJobSwaggerContract(t *testing.T, paths map[string]json.RawMessage, definitions map[string]json.RawMessage) {
	t.Helper()

	createJob := swaggerOperation(t, paths, "/v1/ocr/jobs", "post")
	if !swaggerOperationHasHeader(createJob, "Authorization") {
		t.Fatal("public create OCR job operation missing Authorization header parameter")
	}
	assertContains(t, createJob.Description, "multipart/form-data", "public create OCR job operation description")
	assertContains(t, createJob.Description, "raw API key", "public create OCR job operation description")
	assertContains(t, createJob.Description, "Bearer <api_key>", "public create OCR job operation description")
	assertContains(t, createJob.Description, "PDF, PNG, or JPEG", "public create OCR job operation description")
	assertContains(t, createJob.Description, "schema and schema_id are mutually exclusive", "public create OCR job operation description")
	assertContains(t, createJob.Description, "user_id is not accepted", "public create OCR job operation description")
	assertContains(t, createJob.Description, "Poll GET /v1/ocr/jobs/{id}", "public create OCR job operation description")
	assertContains(t, createJob.Description, "Accepted job response example", "public create OCR job operation description")
	createAuthHeader := swaggerOperationParameter(t, createJob, "header", "Authorization")
	assertContains(t, createAuthHeader.Description, "raw API key", "public create OCR job Authorization description")
	assertContains(t, createAuthHeader.Description, "Bearer <api_key>", "public create OCR job Authorization description")
	fileParam := swaggerOperationParameter(t, createJob, "formData", "file")
	if fileParam.Type != "file" {
		t.Fatalf("public create OCR job file parameter type = %q, want file", fileParam.Type)
	}
	assertContains(t, fileParam.Description, "PDF, PNG, or JPEG", "public create OCR job file description")
	assertContains(t, fileParam.Description, "255 characters", "public create OCR job file description")
	schemaParam := swaggerOperationParameter(t, createJob, "formData", "schema")
	assertContains(t, schemaParam.Description, "Inline JSON Schema", "public create OCR job schema description")
	assertContains(t, schemaParam.Description, "mutually exclusive", "public create OCR job schema description")
	schemaIDParam := swaggerOperationParameter(t, createJob, "formData", "schema_id")
	if schemaIDParam.Format != "uuid" {
		t.Fatalf("public create OCR job schema_id format = %q, want uuid", schemaIDParam.Format)
	}
	assertContains(t, schemaIDParam.Description, "saved schema", "public create OCR job schema_id description")
	assertContains(t, schemaIDParam.Description, "same API key owner", "public create OCR job schema_id description")
	if swaggerOperationHasParameter(createJob, "formData", "user_id") {
		t.Fatal("public create OCR job operation documents forbidden user_id form parameter")
	}
	if got := swaggerOperationResponseRef(t, createJob, "202"); got != "#/definitions/ocrJobResponse" {
		t.Fatalf("public create OCR job 202 response ref = %q, want ocrJobResponse", got)
	}
	assertContains(t, swaggerOperationResponseDescription(t, createJob, "202"), "queued", "public create OCR job 202 description")
	assertContains(t, swaggerOperationResponseDescription(t, createJob, "202"), "Poll GET /v1/ocr/jobs/{id}", "public create OCR job 202 description")
	assertContains(t, swaggerOperationResponseDescription(t, createJob, "400"), "missing file", "public create OCR job 400 description")
	assertContains(t, swaggerOperationResponseDescription(t, createJob, "400"), "schema and schema_id", "public create OCR job 400 description")
	assertContains(t, swaggerOperationResponseDescription(t, createJob, "401"), "API key", "public create OCR job 401 description")
	assertContains(t, swaggerOperationResponseDescription(t, createJob, "402"), "available credits", "public create OCR job 402 description")
	assertContains(t, swaggerOperationResponseDescription(t, createJob, "404"), "saved schema", "public create OCR job 404 description")
	assertContains(t, swaggerOperationResponseDescription(t, createJob, "500"), "unexpected", "public create OCR job 500 description")
	assertSwaggerResponseExampleString(t, createJob, "202", "id", "f2415ed5-5d0b-4c8f-8da4-46fd5cf8f7cb")
	assertSwaggerResponseExampleString(t, createJob, "202", "status", "queued")
	assertSwaggerResponseExampleString(t, createJob, "202", "original_filename", "invoice.pdf")
	assertSwaggerDefinitionRequiredProperties(t, definitions, "ocrJobResponse", "id", "created_at", "original_filename", "mime_type", "status", "file_size", "page_count", "has_inline_schema")
	assertSwaggerDefinitionProperty(t, definitions, "ocrJobResponse", "status", "description", "Current job status. Values are queued, processing, completed, or failed.")
	assertSwaggerDefinitionProperty(t, definitions, "ocrJobResponse", "document_id", "x-nullable", true)
	assertSwaggerDefinitionProperty(t, definitions, "paymentRequiredResponse", "available_credits", "description", "Credits currently available to the API key owner.")
	assertSwaggerDefinitionRequiredProperties(t, definitions, "paymentRequiredResponse", "error", "required_credits", "available_credits")

	getJob := swaggerOperation(t, paths, "/v1/ocr/jobs/{id}", "get")
	if !swaggerOperationHasHeader(getJob, "Authorization") {
		t.Fatal("public get OCR job operation missing Authorization header parameter")
	}
	assertContains(t, getJob.Description, "Poll this endpoint", "public get OCR job operation description")
	assertContains(t, getJob.Description, "queued", "public get OCR job operation description")
	assertContains(t, getJob.Description, "processing", "public get OCR job operation description")
	assertContains(t, getJob.Description, "completed", "public get OCR job operation description")
	assertContains(t, getJob.Description, "failed", "public get OCR job operation description")
	assertContains(t, getJob.Description, "document", "public get OCR job operation description")
	assertContains(t, getJob.Description, "Completed job response example", "public get OCR job operation description")
	assertContains(t, getJob.Description, "Queued job response example", "public get OCR job operation description")
	authHeader := swaggerOperationParameter(t, getJob, "header", "Authorization")
	assertContains(t, authHeader.Description, "raw API key", "public get OCR job Authorization description")
	assertContains(t, authHeader.Description, "Bearer <api_key>", "public get OCR job Authorization description")
	idPath := swaggerOperationParameter(t, getJob, "path", "id")
	assertContains(t, idPath.Description, "POST /v1/ocr/jobs", "public get OCR job id description")
	if swaggerOperationHasParameter(getJob, "query", "user_id") {
		t.Fatal("public get OCR job operation documents forbidden user_id query parameter")
	}
	if got := swaggerOperationResponseRef(t, getJob, "200"); got != "#/definitions/publicOCRJobResponse" {
		t.Fatalf("public get OCR job 200 response ref = %q, want publicOCRJobResponse", got)
	}
	assertContains(t, swaggerOperationResponseDescription(t, getJob, "200"), "document is included only", "public get OCR job 200 description")
	assertContains(t, swaggerOperationResponseDescription(t, getJob, "400"), "invalid", "public get OCR job 400 description")
	assertContains(t, swaggerOperationResponseDescription(t, getJob, "401"), "API key", "public get OCR job 401 description")
	assertContains(t, swaggerOperationResponseDescription(t, getJob, "404"), "owned by the API key", "public get OCR job 404 description")
	assertContains(t, swaggerOperationResponseDescription(t, getJob, "500"), "unexpected", "public get OCR job 500 description")
	assertSwaggerResponseExampleString(t, getJob, "200", "id", "f2415ed5-5d0b-4c8f-8da4-46fd5cf8f7cb")
	assertSwaggerResponseExampleString(t, getJob, "200", "status", "completed")
	assertSwaggerDefinitionProperties(t, definitions, "publicOCRJobResponse", "created_at", "document", "has_inline_schema", "id", "original_filename", "status")
	assertSwaggerDefinitionProperties(t, definitions, "publicOCRJobDocumentResponse", "document_annotation", "document_id", "file_size", "page_count", "pages")
	assertSwaggerDefinitionProperty(t, definitions, "publicOCRJobResponse", "document", "$ref", "#/definitions/publicOCRJobDocumentResponse")
	assertSwaggerDefinitionRequiredProperties(t, definitions, "publicOCRJobResponse", "id", "created_at", "original_filename", "status", "has_inline_schema")
	assertSwaggerDefinitionRequiredProperties(t, definitions, "publicOCRJobDocumentResponse", "document_id", "file_size", "page_count", "pages", "document_annotation")
	assertSwaggerDefinitionProperty(t, definitions, "publicOCRJobDocumentResponse", "document_annotation", "x-nullable", true)
}

func assertPublicBalanceSwaggerContract(t *testing.T, paths map[string]json.RawMessage, definitions map[string]json.RawMessage) {
	t.Helper()

	getBalance := swaggerOperation(t, paths, "/v1/get-balance", "get")
	if !swaggerOperationHasHeader(getBalance, "Authorization") {
		t.Fatal("public get balance operation missing Authorization header parameter")
	}
	assertContains(t, getBalance.Description, "available OCR credits", "public get balance operation description")
	assertContains(t, getBalance.Description, "raw API key", "public get balance operation description")
	assertContains(t, getBalance.Description, "Bearer <api_key>", "public get balance operation description")
	assertContains(t, getBalance.Description, "user_id is not accepted", "public get balance operation description")
	assertContains(t, getBalance.Description, "Balance response example", "public get balance operation description")
	authHeader := swaggerOperationParameter(t, getBalance, "header", "Authorization")
	assertContains(t, authHeader.Description, "raw API key", "public get balance Authorization description")
	assertContains(t, authHeader.Description, "Bearer <api_key>", "public get balance Authorization description")
	if swaggerOperationHasParameter(getBalance, "query", "user_id") {
		t.Fatal("public get balance operation documents forbidden user_id query parameter")
	}
	if got := swaggerOperationResponseRef(t, getBalance, "200"); got != "#/definitions/creditBalanceResponse" {
		t.Fatalf("public get balance 200 response ref = %q, want creditBalanceResponse", got)
	}
	assertContains(t, swaggerOperationResponseDescription(t, getBalance, "200"), "available credits", "public get balance 200 description")
	assertContains(t, swaggerOperationResponseDescription(t, getBalance, "400"), "user_id", "public get balance 400 description")
	assertContains(t, swaggerOperationResponseDescription(t, getBalance, "401"), "API key", "public get balance 401 description")
	assertContains(t, swaggerOperationResponseDescription(t, getBalance, "500"), "unexpected", "public get balance 500 description")
	assertSwaggerResponseExampleString(t, getBalance, "200", "user_id", "550e8400-e29b-41d4-a716-446655440000")
	assertSwaggerResponseExampleNumber(t, getBalance, "200", "available_credits", 750)
	assertSwaggerDefinitionRequiredProperties(t, definitions, "creditBalanceResponse", "user_id", "available_credits")
	assertSwaggerDefinitionProperty(t, definitions, "creditBalanceResponse", "user_id", "description", "User id owned by the supplied public API key.")
	assertSwaggerDefinitionProperty(t, definitions, "creditBalanceResponse", "available_credits", "description", "Available OCR credits for the API key owner.")
}

type swaggerOperationDoc struct {
	Description string                     `json:"description"`
	Parameters  []swaggerParameterDoc      `json:"parameters"`
	Responses   map[string]json.RawMessage `json:"responses"`
}

type swaggerParameterDoc struct {
	Name        string   `json:"name"`
	In          string   `json:"in"`
	Description string   `json:"description"`
	Required    bool     `json:"required"`
	Type        string   `json:"type"`
	Format      string   `json:"format"`
	Enum        []string `json:"enum"`
	Nullable    bool     `json:"x-nullable"`
}

func swaggerOperation(t *testing.T, paths map[string]json.RawMessage, path string, method string) swaggerOperationDoc {
	t.Helper()
	var operations map[string]swaggerOperationDoc
	if err := json.Unmarshal(paths[path], &operations); err != nil {
		t.Fatalf("decode swagger operations for %s: %v", path, err)
	}
	operation, ok := operations[method]
	if !ok {
		t.Fatalf("%s %s operation missing", method, path)
	}
	return operation
}

func swaggerOperationParameter(t *testing.T, operation swaggerOperationDoc, in string, name string) swaggerParameterDoc {
	t.Helper()
	for _, parameter := range operation.Parameters {
		if parameter.In == in && parameter.Name == name {
			return parameter
		}
	}
	t.Fatalf("swagger parameter %s %s missing", in, name)
	return swaggerParameterDoc{}
}

func swaggerOperationHasParameter(operation swaggerOperationDoc, in string, name string) bool {
	for _, parameter := range operation.Parameters {
		if parameter.In == in && parameter.Name == name {
			return true
		}
	}
	return false
}

func swaggerOperationHasHeader(operation swaggerOperationDoc, name string) bool {
	for _, parameter := range operation.Parameters {
		if parameter.In == "header" && parameter.Name == name {
			return true
		}
	}
	return false
}

func swaggerOperationResponseRef(t *testing.T, operation swaggerOperationDoc, status string) string {
	t.Helper()
	raw, ok := operation.Responses[status]
	if !ok {
		t.Fatalf("swagger response %s missing", status)
	}
	var response struct {
		Schema map[string]any `json:"schema"`
	}
	if err := json.Unmarshal(raw, &response); err != nil {
		t.Fatalf("decode swagger response %s: %v", status, err)
	}
	ref, ok := response.Schema["$ref"].(string)
	if !ok {
		t.Fatalf("swagger response %s schema ref missing: %#v", status, response.Schema)
	}
	return ref
}

func swaggerOperationResponseDescription(t *testing.T, operation swaggerOperationDoc, status string) string {
	t.Helper()
	raw, ok := operation.Responses[status]
	if !ok {
		t.Fatalf("swagger response %s missing", status)
	}
	var response struct {
		Description string `json:"description"`
	}
	if err := json.Unmarshal(raw, &response); err != nil {
		t.Fatalf("decode swagger response %s: %v", status, err)
	}
	return response.Description
}

func assertSwaggerResponseExampleString(t *testing.T, operation swaggerOperationDoc, status string, key string, want string) {
	t.Helper()
	got := swaggerResponseExampleField(t, operation, status, key)
	if got != want {
		t.Fatalf("swagger response %s example %s = %#v, want %#v", status, key, got, want)
	}
}

func assertSwaggerResponseExampleNumber(t *testing.T, operation swaggerOperationDoc, status string, key string, want float64) {
	t.Helper()
	got := swaggerResponseExampleField(t, operation, status, key)
	if got != want {
		t.Fatalf("swagger response %s example %s = %#v, want %#v", status, key, got, want)
	}
}

func swaggerResponseExampleField(t *testing.T, operation swaggerOperationDoc, status string, key string) any {
	t.Helper()
	raw, ok := operation.Responses[status]
	if !ok {
		t.Fatalf("swagger response %s missing", status)
	}
	var response struct {
		Examples map[string]map[string]any `json:"examples"`
	}
	if err := json.Unmarshal(raw, &response); err != nil {
		t.Fatalf("decode swagger response %s: %v", status, err)
	}
	example, ok := response.Examples["application/json"]
	if !ok {
		t.Fatalf("swagger response %s application/json example missing: %#v", status, response.Examples)
	}
	got, ok := example[key]
	if !ok {
		t.Fatalf("swagger response %s example %s missing: %#v", status, key, example)
	}
	return got
}

func assertSwaggerDefinitionProperty(t *testing.T, definitions map[string]json.RawMessage, definitionName string, propertyName string, key string, want any) {
	t.Helper()
	var definition struct {
		Properties map[string]map[string]any `json:"properties"`
	}
	raw, ok := definitions[definitionName]
	if !ok {
		t.Fatalf("%s definition missing", definitionName)
	}
	if err := json.Unmarshal(raw, &definition); err != nil {
		t.Fatalf("decode %s definition: %v", definitionName, err)
	}
	property, ok := definition.Properties[propertyName]
	if !ok {
		t.Fatalf("%s.%s property missing", definitionName, propertyName)
	}
	if got := property[key]; got != want {
		t.Fatalf("%s.%s %s = %#v, want %#v", definitionName, propertyName, key, got, want)
	}
}

func assertSwaggerDefinitionRequiredProperties(t *testing.T, definitions map[string]json.RawMessage, definitionName string, want ...string) {
	t.Helper()
	var definition struct {
		Required []string `json:"required"`
	}
	raw, ok := definitions[definitionName]
	if !ok {
		t.Fatalf("%s definition missing", definitionName)
	}
	if err := json.Unmarshal(raw, &definition); err != nil {
		t.Fatalf("decode %s definition: %v", definitionName, err)
	}
	if !reflect.DeepEqual(definition.Required, want) {
		t.Fatalf("%s required = %#v, want %#v", definitionName, definition.Required, want)
	}
}

func assertSwaggerDefinitionProperties(t *testing.T, definitions map[string]json.RawMessage, definitionName string, want ...string) {
	t.Helper()
	var definition struct {
		Properties map[string]json.RawMessage `json:"properties"`
	}
	raw, ok := definitions[definitionName]
	if !ok {
		t.Fatalf("%s definition missing", definitionName)
	}
	if err := json.Unmarshal(raw, &definition); err != nil {
		t.Fatalf("decode %s definition: %v", definitionName, err)
	}
	if len(definition.Properties) != len(want) {
		t.Fatalf("%s properties = %#v, want only %#v", definitionName, definition.Properties, want)
	}
	for _, propertyName := range want {
		if _, ok := definition.Properties[propertyName]; !ok {
			t.Fatalf("%s.%s property missing", definitionName, propertyName)
		}
	}
}

func assertContains(t *testing.T, got string, want string, label string) {
	t.Helper()
	if !strings.Contains(got, want) {
		t.Fatalf("%s = %q, want substring %q", label, got, want)
	}
}

func assertSwaggerDefinitionMissing(t *testing.T, definitions map[string]json.RawMessage, definitionName string) {
	t.Helper()
	if _, ok := definitions[definitionName]; ok {
		t.Fatalf("%s definition present, want missing", definitionName)
	}
}
