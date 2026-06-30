package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"ai.ro/syncra/internal/auth"
	"ai.ro/syncra/internal/ocr"
)

func testRouter(t *testing.T) (*gin.Engine, *gorm.DB) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	db := apiPostgresTx(t)
	h := &Handler{
		DB:               db,
		OCR:              successfulTestOCRProcessor(),
		MaxUploadBytes:   20 << 20,
		InternalAPIToken: testInternalAPIToken,
	}
	return NewRouter(h), db
}

func createTestUser(t *testing.T, db *gorm.DB, email string) auth.User {
	t.Helper()
	user := auth.User{Name: "Test User", Email: email}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create test user: %v", err)
	}
	return user
}

func successfulTestOCRProcessor() OCRProcessor {
	return func(ctx context.Context, input OCRProcessInput) (*MistralOCRResponse, []byte, error) {
		resp := &MistralOCRResponse{
			Pages: []MistralOCRPage{{Index: 0, Markdown: "# Test"}},
			Model: "mistral-ocr-latest",
		}
		raw := []byte(`{"pages":[{"index":0,"markdown":"# Test"}],"model":"mistral-ocr-latest"}`)
		if len(input.Schema) > 0 {
			resp.DocumentAnnotation = ptrString(`{"ok":true}`)
			raw = []byte(`{"pages":[{"index":0,"markdown":"# Test"}],"document_annotation":"{\"ok\":true}","model":"mistral-ocr-latest"}`)
		}
		return resp, raw, nil
	}
}

func assertJSONEqual(t *testing.T, got json.RawMessage, want string) {
	t.Helper()
	var gotValue any
	if err := json.Unmarshal(got, &gotValue); err != nil {
		t.Fatalf("decode got JSON: %v", err)
	}
	var wantValue any
	if err := json.Unmarshal([]byte(want), &wantValue); err != nil {
		t.Fatalf("decode want JSON: %v", err)
	}
	if !reflect.DeepEqual(gotValue, wantValue) {
		t.Fatalf("json = %s, want %s", string(got), want)
	}
}

func decodeSchemaListResponse(t *testing.T, w *httptest.ResponseRecorder) SchemaListResponse {
	t.Helper()
	var got SchemaListResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode schema list: %v", err)
	}
	return got
}

func assertSchemaListIDs(t *testing.T, schemas []SchemaResponse, want ...uuid.UUID) {
	t.Helper()
	if len(schemas) != len(want) {
		t.Fatalf("schema count = %d, want %d: %#v", len(schemas), len(want), schemas)
	}
	for i, wantID := range want {
		if schemas[i].ID != wantID {
			t.Fatalf("schema[%d].id = %s, want %s", i, schemas[i].ID, wantID)
		}
	}
}

func decodeDeleteSchemasResponse(t *testing.T, w *httptest.ResponseRecorder) struct {
	DeletedIDs   []uuid.UUID `json:"deleted_ids"`
	DeletedCount int         `json:"deleted_count"`
} {
	t.Helper()
	var got struct {
		DeletedIDs   []uuid.UUID `json:"deleted_ids"`
		DeletedCount int         `json:"deleted_count"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode delete schemas response: %v body=%s", err, w.Body.String())
	}
	return got
}

func TestCreateSchema(t *testing.T) {
	router, db := testRouter(t)
	wantSchema := `{"type":"object","properties":{"Furnizor":{"type":"string"}}}`
	body := []byte(`{"name":" invoice ","description":"Invoice fields","schema":` + wantSchema + `,"strict":true}`)

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodPost, "/api/ocr/schemas", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got SchemaResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ID == uuid.Nil || got.Name != "invoice" || got.Description != "Invoice fields" || !got.Strict {
		t.Fatalf("unexpected response: %#v", got)
	}
	assertJSONEqual(t, got.Schema, wantSchema)

	var stored ocr.ExtractionSchema
	if err := db.First(&stored, "id = ?", got.ID).Error; err != nil {
		t.Fatalf("load schema: %v", err)
	}
	if stored.Name != "invoice" || stored.Description != "Invoice fields" || !stored.Strict {
		t.Fatalf("unexpected stored schema: %#v", stored)
	}
	assertJSONEqual(t, json.RawMessage(stored.SchemaJSON), wantSchema)
}

func TestCreateSchemaPersistsFlexibleStrictMode(t *testing.T) {
	router, db := testRouter(t)
	body := []byte(`{"name":"flexible","schema":{"type":"object"},"strict":false}`)

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodPost, "/api/ocr/schemas", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got SchemaResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Strict {
		t.Fatalf("response strict = true, want false")
	}

	var stored ocr.ExtractionSchema
	if err := db.First(&stored, "id = ?", got.ID).Error; err != nil {
		t.Fatalf("load schema: %v", err)
	}
	if stored.Strict {
		t.Fatalf("stored strict = true, want false")
	}
}

func TestCreateSchemaWithUserID(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "schema-owner@example.com")
	wantSchema := `{"type":"object","properties":{"Furnizor":{"type":"string"}}}`
	body := []byte(`{"name":"invoice","schema":` + wantSchema + `,"user_id":"` + user.ID + `"}`)

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodPost, "/api/ocr/schemas", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got SchemaResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.UserID == nil || string(*got.UserID) != user.ID {
		t.Fatalf("user_id = %v, want %q", got.UserID, user.ID)
	}

	var stored ocr.ExtractionSchema
	if err := db.First(&stored, "id = ?", got.ID).Error; err != nil {
		t.Fatalf("load schema: %v", err)
	}
	if stored.UserID == nil || *stored.UserID != user.ID {
		t.Fatalf("stored user_id = %v, want %q", stored.UserID, user.ID)
	}
}

func TestCreateSchemaTreatsNullAndEmptyUserIDAsSystemWide(t *testing.T) {
	cases := []struct {
		name   string
		userID string
	}{
		{name: "null", userID: "null"},
		{name: "empty", userID: `""`},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			router, db := testRouter(t)
			body := []byte(`{"name":"system-` + tt.name + `","schema":{"type":"object"},"user_id":` + tt.userID + `}`)

			w := httptest.NewRecorder()
			req := newTestRequest(http.MethodPost, "/api/ocr/schemas", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			if w.Code != http.StatusCreated {
				t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
			}
			var got SchemaResponse
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Fatalf("decode: %v", err)
			}
			if got.UserID != nil {
				t.Fatalf("user_id = %v, want nil", got.UserID)
			}

			var stored ocr.ExtractionSchema
			if err := db.First(&stored, "id = ?", got.ID).Error; err != nil {
				t.Fatalf("load schema: %v", err)
			}
			if stored.UserID != nil {
				t.Fatalf("stored user_id = %v, want nil", stored.UserID)
			}
		})
	}
}

func TestCreateSchemaRejectsMalformedUserID(t *testing.T) {
	router, _ := testRouter(t)
	body := []byte(`{"name":"bad-owner","schema":{"type":"object"},"user_id":"not-a-uuid"}`)

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodPost, "/api/ocr/schemas", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Error != "invalid user_id" {
		t.Fatalf("error = %q", got.Error)
	}
}

func TestCreateSchemaRejectsNonStringUserID(t *testing.T) {
	router, _ := testRouter(t)
	body := []byte(`{"name":"bad-owner","schema":{"type":"object"},"user_id":123}`)

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodPost, "/api/ocr/schemas", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Error != "invalid user_id" {
		t.Fatalf("error = %q", got.Error)
	}
}

func TestCreateSchemaRejectsUnknownUserID(t *testing.T) {
	router, _ := testRouter(t)
	body := []byte(`{"name":"bad-owner","schema":{"type":"object"},"user_id":"` + uuid.NewString() + `"}`)

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodPost, "/api/ocr/schemas", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Error != "invalid user_id" {
		t.Fatalf("error = %q", got.Error)
	}
}

func TestCreateSchemaRejectsNonObjectSchema(t *testing.T) {
	router, _ := testRouter(t)
	body := []byte(`{"name":"bad","schema":["not","object"]}`)

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodPost, "/api/ocr/schemas", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestCreateSchemaRejectsInvalidJSONSchema(t *testing.T) {
	router, db := testRouter(t)
	body := []byte(`{"name":"bad","schema":{"type":"strung"}}`)

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodPost, "/api/ocr/schemas", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Error != "schema must be a valid JSON Schema" {
		t.Fatalf("error = %q", got.Error)
	}
	var count int64
	if err := db.Model(&ocr.ExtractionSchema{}).Count(&count).Error; err != nil {
		t.Fatalf("count schemas: %v", err)
	}
	if count != 0 {
		t.Fatalf("schema count = %d, want 0", count)
	}
}

func TestCreateSchemaRejectsRequestBodyOverLimit(t *testing.T) {
	router, _ := testRouter(t)
	body := []byte(`{"name":"too-big","schema":{"type":"object","padding":"` + strings.Repeat("a", 2<<20) + `"}}`)

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodPost, "/api/ocr/schemas", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Error != "request body too large" {
		t.Fatalf("error = %q", got.Error)
	}
}

func TestCreateSchemaRejectsNullSchema(t *testing.T) {
	router, _ := testRouter(t)
	body := []byte(`{"name":"bad","schema":null}`)

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodPost, "/api/ocr/schemas", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestCreateSchemaRejectsWhitespaceName(t *testing.T) {
	router, _ := testRouter(t)
	body := []byte(`{"name":"   ","schema":{"type":"object"}}`)

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodPost, "/api/ocr/schemas", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestCreateSchemaRejectsNameOver160Characters(t *testing.T) {
	router, _ := testRouter(t)
	body := []byte(`{"name":"` + strings.Repeat("a", 161) + `","schema":{"type":"object"}}`)

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodPost, "/api/ocr/schemas", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Error != "name must be at most 160 characters" {
		t.Fatalf("error = %q", got.Error)
	}
}

func TestListAndGetSchema(t *testing.T) {
	router, db := testRouter(t)
	older := ocr.ExtractionSchema{
		CreatedAt:   time.Now().Add(-time.Hour),
		UpdatedAt:   time.Now().Add(-time.Hour),
		Name:        "invoice",
		Description: "Invoice fields",
		SchemaJSON:  datatypes.JSON([]byte(`{"type":"object","title":"invoice"}`)),
		Strict:      true,
	}
	newer := ocr.ExtractionSchema{
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Name:        "receipt",
		Description: "Receipt fields",
		SchemaJSON:  datatypes.JSON([]byte(`{"type":"object","title":"receipt"}`)),
		Strict:      false,
	}
	if err := db.Create(&older).Error; err != nil {
		t.Fatalf("create older schema: %v", err)
	}
	if err := db.Create(&newer).Error; err != nil {
		t.Fatalf("create newer schema: %v", err)
	}
	if err := db.Model(&newer).Update("strict", false).Error; err != nil {
		t.Fatalf("update newer strict: %v", err)
	}
	newer.Strict = false

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodGet, "/api/ocr/schemas", nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("list status = %d body=%s", w.Code, w.Body.String())
	}
	listResponse := decodeSchemaListResponse(t, w)
	list := listResponse.Schemas
	if listResponse.NextCursor != nil {
		t.Fatalf("next_cursor = %q, want nil", *listResponse.NextCursor)
	}
	if len(list) != 2 {
		t.Fatalf("list length = %d, want 2: %#v", len(list), list)
	}
	if list[0].ID != newer.ID || list[0].Name != "receipt" || list[0].Description != "Receipt fields" || list[0].Strict {
		t.Fatalf("unexpected first list item: %#v", list[0])
	}
	assertJSONEqual(t, list[0].Schema, `{"type":"object","title":"receipt"}`)
	if list[1].ID != older.ID || list[1].Name != "invoice" || list[1].Description != "Invoice fields" || !list[1].Strict {
		t.Fatalf("unexpected second list item: %#v", list[1])
	}
	assertJSONEqual(t, list[1].Schema, `{"type":"object","title":"invoice"}`)

	w = httptest.NewRecorder()
	req = newTestRequest(http.MethodGet, fmt.Sprintf("/api/ocr/schemas/%s", older.ID), nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("get status = %d body=%s", w.Code, w.Body.String())
	}
	var got SchemaResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode get: %v", err)
	}
	if got.ID != older.ID || got.Name != "invoice" || got.Description != "Invoice fields" || !got.Strict {
		t.Fatalf("unexpected get response: %#v", got)
	}
	assertJSONEqual(t, got.Schema, `{"type":"object","title":"invoice"}`)
}

func TestListSchemasScopesToSystemWideByDefault(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "scoped-default@example.com")
	systemSchema := ocr.ExtractionSchema{
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Name:       "system",
		SchemaJSON: datatypes.JSON([]byte(`{"type":"object","title":"system"}`)),
		Strict:     true,
	}
	userSchema := ocr.ExtractionSchema{
		CreatedAt:  time.Now().Add(time.Second),
		UpdatedAt:  time.Now().Add(time.Second),
		Name:       "user",
		SchemaJSON: datatypes.JSON([]byte(`{"type":"object","title":"user"}`)),
		Strict:     true,
		UserID:     &user.ID,
	}
	if err := db.Create(&systemSchema).Error; err != nil {
		t.Fatalf("create system schema: %v", err)
	}
	if err := db.Create(&userSchema).Error; err != nil {
		t.Fatalf("create user schema: %v", err)
	}

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodGet, "/api/ocr/schemas", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeSchemaListResponse(t, w)
	if len(got.Schemas) != 1 || got.Schemas[0].ID != systemSchema.ID || got.Schemas[0].UserID != nil {
		t.Fatalf("schemas = %#v, want only system schema", got.Schemas)
	}
}

func TestListSchemasScopesToUserID(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "scoped-user@example.com")
	other := createTestUser(t, db, "scoped-other@example.com")
	systemSchema := ocr.ExtractionSchema{
		Name:       "system",
		SchemaJSON: datatypes.JSON([]byte(`{"type":"object","title":"system"}`)),
		Strict:     true,
	}
	userSchema := ocr.ExtractionSchema{
		Name:       "user",
		SchemaJSON: datatypes.JSON([]byte(`{"type":"object","title":"user"}`)),
		Strict:     true,
		UserID:     &user.ID,
	}
	otherSchema := ocr.ExtractionSchema{
		Name:       "other",
		SchemaJSON: datatypes.JSON([]byte(`{"type":"object","title":"other"}`)),
		Strict:     true,
		UserID:     &other.ID,
	}
	if err := db.Create(&systemSchema).Error; err != nil {
		t.Fatalf("create system schema: %v", err)
	}
	if err := db.Create(&userSchema).Error; err != nil {
		t.Fatalf("create user schema: %v", err)
	}
	if err := db.Create(&otherSchema).Error; err != nil {
		t.Fatalf("create other schema: %v", err)
	}

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodGet, "/api/ocr/schemas?user_id="+user.ID, nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeSchemaListResponse(t, w)
	if len(got.Schemas) != 1 || got.Schemas[0].ID != userSchema.ID || got.Schemas[0].UserID == nil || string(*got.Schemas[0].UserID) != user.ID {
		t.Fatalf("schemas = %#v, want only user's schema", got.Schemas)
	}
}

func TestListSchemasPaginatesWithCursor(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "schema-list-cursor@example.com")
	base := time.Date(2026, 5, 29, 10, 0, 0, 0, time.UTC)
	oldest := ocr.ExtractionSchema{
		CreatedAt:  base,
		UpdatedAt:  base,
		Name:       "oldest",
		SchemaJSON: datatypes.JSON([]byte(`{"type":"object","title":"oldest"}`)),
		Strict:     true,
		UserID:     &user.ID,
	}
	middle := ocr.ExtractionSchema{
		CreatedAt:  base.Add(time.Minute),
		UpdatedAt:  base.Add(time.Minute),
		Name:       "middle",
		SchemaJSON: datatypes.JSON([]byte(`{"type":"object","title":"middle"}`)),
		Strict:     true,
		UserID:     &user.ID,
	}
	newest := ocr.ExtractionSchema{
		CreatedAt:  base.Add(2 * time.Minute),
		UpdatedAt:  base.Add(2 * time.Minute),
		Name:       "newest",
		SchemaJSON: datatypes.JSON([]byte(`{"type":"object","title":"newest"}`)),
		Strict:     true,
		UserID:     &user.ID,
	}
	if err := db.Create(&oldest).Error; err != nil {
		t.Fatalf("create oldest schema: %v", err)
	}
	if err := db.Create(&middle).Error; err != nil {
		t.Fatalf("create middle schema: %v", err)
	}
	if err := db.Create(&newest).Error; err != nil {
		t.Fatalf("create newest schema: %v", err)
	}

	req := newTestRequest(http.MethodGet, "/api/ocr/schemas?user_id="+user.ID+"&size=2&sort=desc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("first page status = %d body=%s", w.Code, w.Body.String())
	}
	firstPage := decodeSchemaListResponse(t, w)
	assertSchemaListIDs(t, firstPage.Schemas, newest.ID, middle.ID)
	if firstPage.NextCursor == nil || *firstPage.NextCursor == "" {
		t.Fatalf("next_cursor = %#v, want non-empty", firstPage.NextCursor)
	}

	query := url.Values{}
	query.Set("user_id", user.ID)
	query.Set("size", "2")
	query.Set("sort", "desc")
	query.Set("cursor", *firstPage.NextCursor)
	req = newTestRequest(http.MethodGet, "/api/ocr/schemas?"+query.Encode(), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("second page status = %d body=%s", w.Code, w.Body.String())
	}
	secondPage := decodeSchemaListResponse(t, w)
	assertSchemaListIDs(t, secondPage.Schemas, oldest.ID)
	if secondPage.NextCursor != nil {
		t.Fatalf("next_cursor = %q, want nil", *secondPage.NextCursor)
	}
}

func TestListSchemasRejectsInvalidPaginationParameters(t *testing.T) {
	router, _ := testRouter(t)

	cases := []struct {
		name      string
		query     string
		wantError string
	}{
		{name: "invalid cursor", query: "cursor=not-base64", wantError: "invalid cursor"},
		{name: "size zero", query: "size=0", wantError: "size must be between 1 and 100"},
		{name: "size too large", query: "size=101", wantError: "size must be between 1 and 100"},
		{name: "invalid sort", query: "sort=sideways", wantError: "sort must be asc or desc"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := newTestRequest(http.MethodGet, "/api/ocr/schemas?"+tc.query, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
			}
			var got ErrorResponse
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Fatalf("decode error: %v", err)
			}
			if got.Error != tc.wantError {
				t.Fatalf("error = %q, want %q", got.Error, tc.wantError)
			}
		})
	}
}

func TestListSchemasRejectsMalformedUserID(t *testing.T) {
	router, _ := testRouter(t)

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodGet, "/api/ocr/schemas?user_id=not-a-uuid", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Error != "invalid user_id" {
		t.Fatalf("error = %q", got.Error)
	}
}

func TestGetSchemaScopesByQueryUserIDWhenProvided(t *testing.T) {
	router, db := testRouter(t)
	owner := createTestUser(t, db, "schema-get-owner@example.com")
	other := createTestUser(t, db, "schema-get-other@example.com")
	schema := ocr.ExtractionSchema{
		Name:       "owned",
		SchemaJSON: datatypes.JSON([]byte(`{"type":"object","title":"owned"}`)),
		Strict:     true,
		UserID:     &owner.ID,
	}
	if err := db.Create(&schema).Error; err != nil {
		t.Fatalf("create schema: %v", err)
	}

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodGet, fmt.Sprintf("/api/ocr/schemas/%s?user_id=%s", schema.ID, other.ID), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var errResponse ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &errResponse); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if errResponse.Error != "schema not found" {
		t.Fatalf("error = %q, want schema not found", errResponse.Error)
	}

	w = httptest.NewRecorder()
	req = newTestRequest(http.MethodGet, fmt.Sprintf("/api/ocr/schemas/%s?user_id=%s", schema.ID, owner.ID), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("owner status = %d body=%s", w.Code, w.Body.String())
	}
	var got SchemaResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ID != schema.ID || got.UserID == nil || string(*got.UserID) != owner.ID {
		t.Fatalf("schema = %#v, want schema by id with owner user_id", got)
	}
}

func TestGetSchemaPreservesIDOnlyLookupWithoutQueryUserID(t *testing.T) {
	router, db := testRouter(t)
	owner := createTestUser(t, db, "schema-get-unscoped-owner@example.com")
	schema := ocr.ExtractionSchema{
		Name:       "owned",
		SchemaJSON: datatypes.JSON([]byte(`{"type":"object","title":"owned"}`)),
		Strict:     true,
		UserID:     &owner.ID,
	}
	if err := db.Create(&schema).Error; err != nil {
		t.Fatalf("create schema: %v", err)
	}

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodGet, fmt.Sprintf("/api/ocr/schemas/%s", schema.ID), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got SchemaResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ID != schema.ID || got.UserID == nil || string(*got.UserID) != owner.ID {
		t.Fatalf("schema = %#v, want schema by id with owner user_id", got)
	}
}

func TestUpdateSchemaScopesToUserID(t *testing.T) {
	router, db := testRouter(t)
	owner := createTestUser(t, db, "schema-update-owner@example.com")
	other := createTestUser(t, db, "schema-update-other@example.com")
	schema := ocr.ExtractionSchema{
		Name:        "old",
		Description: "Old description",
		SchemaJSON:  datatypes.JSON([]byte(`{"type":"object","title":"old"}`)),
		Strict:      true,
		UserID:      &owner.ID,
	}
	if err := db.Create(&schema).Error; err != nil {
		t.Fatalf("create schema: %v", err)
	}
	body := []byte(`{"name":" updated ","description":"Updated description","schema":{"type":"object","title":"updated"},"strict":false}`)

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodPut, fmt.Sprintf("/api/ocr/schemas/%s?user_id=%s", schema.ID, other.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("wrong-owner status = %d body=%s", w.Code, w.Body.String())
	}

	w = httptest.NewRecorder()
	req = newTestRequest(http.MethodPut, fmt.Sprintf("/api/ocr/schemas/%s?user_id=%s", schema.ID, owner.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("owner status = %d body=%s", w.Code, w.Body.String())
	}
	var got SchemaResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ID != schema.ID || got.Name != "updated" || got.Description != "Updated description" || got.Strict {
		t.Fatalf("schema = %#v, want updated schema", got)
	}
	assertJSONEqual(t, got.Schema, `{"type":"object","title":"updated"}`)

	var stored ocr.ExtractionSchema
	if err := db.First(&stored, "id = ?", schema.ID).Error; err != nil {
		t.Fatalf("load stored schema: %v", err)
	}
	if stored.Name != "updated" || stored.Description != "Updated description" || stored.Strict || stored.UserID == nil || *stored.UserID != owner.ID {
		t.Fatalf("stored schema = %#v, want updated owned schema", stored)
	}
	assertJSONEqual(t, json.RawMessage(stored.SchemaJSON), `{"type":"object","title":"updated"}`)
}

func TestUpdateSchemaRejectsInvalidInputs(t *testing.T) {
	router, db := testRouter(t)
	schema := ocr.ExtractionSchema{
		Name:       "old",
		SchemaJSON: datatypes.JSON([]byte(`{"type":"object","title":"old"}`)),
		Strict:     true,
	}
	if err := db.Create(&schema).Error; err != nil {
		t.Fatalf("create schema: %v", err)
	}

	cases := []struct {
		name      string
		target    string
		body      string
		wantError string
	}{
		{name: "invalid id", target: "/api/ocr/schemas/not-a-uuid", body: `{"name":"ok","schema":{"type":"object"}}`, wantError: "invalid schema id"},
		{name: "malformed user id", target: fmt.Sprintf("/api/ocr/schemas/%s?user_id=not-a-uuid", schema.ID), body: `{"name":"ok","schema":{"type":"object"}}`, wantError: "invalid user_id"},
		{name: "malformed json", target: "/api/ocr/schemas/" + schema.ID.String(), body: `{`, wantError: "invalid JSON body"},
		{name: "missing name", target: "/api/ocr/schemas/" + schema.ID.String(), body: `{"schema":{"type":"object"}}`, wantError: "name is required"},
		{name: "long name", target: "/api/ocr/schemas/" + schema.ID.String(), body: `{"name":"` + strings.Repeat("a", 161) + `","schema":{"type":"object"}}`, wantError: "name must be at most 160 characters"},
		{name: "non-object schema", target: "/api/ocr/schemas/" + schema.ID.String(), body: `{"name":"bad","schema":[]}`, wantError: "schema must be a JSON object"},
		{name: "invalid JSON schema", target: "/api/ocr/schemas/" + schema.ID.String(), body: `{"name":"bad","schema":{"type":"strung"}}`, wantError: "schema must be a valid JSON Schema"},
		{name: "non-string description", target: "/api/ocr/schemas/" + schema.ID.String(), body: `{"name":"bad","description":1,"schema":{"type":"object"}}`, wantError: "invalid JSON body"},
		{name: "non-boolean strict", target: "/api/ocr/schemas/" + schema.ID.String(), body: `{"name":"bad","schema":{"type":"object"},"strict":"true"}`, wantError: "invalid JSON body"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := newTestRequest(http.MethodPut, tc.target, strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
			}
			var got ErrorResponse
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Fatalf("decode error: %v", err)
			}
			if got.Error != tc.wantError {
				t.Fatalf("error = %q, want %q", got.Error, tc.wantError)
			}
		})
	}
}

func TestUpdateSchemaIgnoresBodyUserID(t *testing.T) {
	router, db := testRouter(t)
	owner := createTestUser(t, db, "schema-update-body-user-id-owner@example.com")
	schema := ocr.ExtractionSchema{
		Name:       "old",
		SchemaJSON: datatypes.JSON([]byte(`{"type":"object","title":"old"}`)),
		Strict:     true,
		UserID:     &owner.ID,
	}
	if err := db.Create(&schema).Error; err != nil {
		t.Fatalf("create schema: %v", err)
	}
	body := []byte(`{"name":"updated","description":"Updated","schema":{"type":"object","title":"updated"},"strict":false,"user_id":123}`)

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodPut, fmt.Sprintf("/api/ocr/schemas/%s?user_id=%s", schema.ID, owner.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got SchemaResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.UserID == nil || string(*got.UserID) != owner.ID || got.Name != "updated" || got.Strict {
		t.Fatalf("schema = %#v, want updated schema retaining query owner", got)
	}
}

func TestDeleteSchemaScopesToUserID(t *testing.T) {
	router, db := testRouter(t)
	owner := createTestUser(t, db, "schema-delete-owner@example.com")
	other := createTestUser(t, db, "schema-delete-other@example.com")
	schema := ocr.ExtractionSchema{
		Name:       "owned",
		SchemaJSON: datatypes.JSON([]byte(`{"type":"object","title":"owned"}`)),
		Strict:     true,
		UserID:     &owner.ID,
	}
	if err := db.Create(&schema).Error; err != nil {
		t.Fatalf("create schema: %v", err)
	}

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodDelete, fmt.Sprintf("/api/ocr/schemas/%s?user_id=%s", schema.ID, other.ID), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("wrong-owner status = %d body=%s", w.Code, w.Body.String())
	}

	var count int64
	if err := db.Model(&ocr.ExtractionSchema{}).Where("id = ?", schema.ID).Count(&count).Error; err != nil {
		t.Fatalf("count schema: %v", err)
	}
	if count != 1 {
		t.Fatalf("schema count after wrong-owner delete = %d, want 1", count)
	}

	w = httptest.NewRecorder()
	req = newTestRequest(http.MethodDelete, fmt.Sprintf("/api/ocr/schemas/%s?user_id=%s", schema.ID, owner.ID), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("owner status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeDeleteSchemasResponse(t, w)
	if got.DeletedCount != 1 || len(got.DeletedIDs) != 1 || got.DeletedIDs[0] != schema.ID {
		t.Fatalf("delete response = %#v, want deleted schema %s", got, schema.ID)
	}
	if err := db.Model(&ocr.ExtractionSchema{}).Where("id = ?", schema.ID).Count(&count).Error; err != nil {
		t.Fatalf("count schema after delete: %v", err)
	}
	if count != 0 {
		t.Fatalf("schema count after delete = %d, want 0", count)
	}
}

func TestDeleteSchemaRejectsInvalidIDAndMissingSchema(t *testing.T) {
	router, _ := testRouter(t)

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodDelete, "/api/ocr/schemas/not-a-uuid", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("invalid-id status = %d body=%s", w.Code, w.Body.String())
	}
	var got ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode invalid-id error: %v", err)
	}
	if got.Error != "invalid schema id" {
		t.Fatalf("invalid-id error = %q, want invalid schema id", got.Error)
	}

	w = httptest.NewRecorder()
	req = newTestRequest(http.MethodDelete, "/api/ocr/schemas/"+uuid.NewString(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("missing status = %d body=%s", w.Code, w.Body.String())
	}
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode missing error: %v", err)
	}
	if got.Error != "schema not found" {
		t.Fatalf("missing error = %q, want schema not found", got.Error)
	}
}

func TestDeleteSchemasScopesToUserID(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "schema-bulk-delete@example.com")
	other := createTestUser(t, db, "schema-bulk-delete-other@example.com")
	owned1 := ocr.ExtractionSchema{Name: "owned-1", SchemaJSON: datatypes.JSON([]byte(`{"type":"object","title":"owned-1"}`)), Strict: true, UserID: &user.ID}
	owned2 := ocr.ExtractionSchema{Name: "owned-2", SchemaJSON: datatypes.JSON([]byte(`{"type":"object","title":"owned-2"}`)), Strict: true, UserID: &user.ID}
	otherSchema := ocr.ExtractionSchema{Name: "other", SchemaJSON: datatypes.JSON([]byte(`{"type":"object","title":"other"}`)), Strict: true, UserID: &other.ID}
	systemSchema := ocr.ExtractionSchema{Name: "system", SchemaJSON: datatypes.JSON([]byte(`{"type":"object","title":"system"}`)), Strict: true}
	for _, schema := range []*ocr.ExtractionSchema{&owned1, &owned2, &otherSchema, &systemSchema} {
		if err := db.Create(schema).Error; err != nil {
			t.Fatalf("create schema %s: %v", schema.Name, err)
		}
	}
	missingID := uuid.New()
	body := `{"ids":["` + owned1.ID.String() + `","` + otherSchema.ID.String() + `","` + missingID.String() + `","` + owned1.ID.String() + `","` + owned2.ID.String() + `"]}`

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodDelete, "/api/ocr/schemas?user_id="+url.QueryEscape(user.ID), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeDeleteSchemasResponse(t, w)
	if got.DeletedCount != 2 || len(got.DeletedIDs) != 2 || got.DeletedIDs[0] != owned1.ID || got.DeletedIDs[1] != owned2.ID {
		t.Fatalf("delete response = %#v, want owned schemas in request order", got)
	}

	var count int64
	if err := db.Model(&ocr.ExtractionSchema{}).Where("id IN ?", []uuid.UUID{owned1.ID, owned2.ID}).Count(&count).Error; err != nil {
		t.Fatalf("count owned schemas: %v", err)
	}
	if count != 0 {
		t.Fatalf("owned schema count = %d, want 0", count)
	}
	if err := db.Model(&ocr.ExtractionSchema{}).Where("id IN ?", []uuid.UUID{otherSchema.ID, systemSchema.ID}).Count(&count).Error; err != nil {
		t.Fatalf("count non-owned schemas: %v", err)
	}
	if count != 2 {
		t.Fatalf("non-owned schema count = %d, want 2", count)
	}
}

func TestDeleteSchemasRejectsInvalidBulkBodies(t *testing.T) {
	router, _ := testRouter(t)
	tooManyIDs := make([]string, 0, 101)
	for range 101 {
		tooManyIDs = append(tooManyIDs, uuid.NewString())
	}
	tooManyIDsBody, err := json.Marshal(map[string][]string{"ids": tooManyIDs})
	if err != nil {
		t.Fatalf("marshal too many ids body: %v", err)
	}

	cases := []struct {
		name      string
		query     string
		body      string
		wantError string
	}{
		{name: "malformed user id", query: "?user_id=not-a-uuid", body: `{"ids":["` + uuid.NewString() + `"]}`, wantError: "invalid user_id"},
		{name: "malformed json", body: `{`, wantError: "invalid schema delete request"},
		{name: "missing ids", body: `{}`, wantError: "ids is required"},
		{name: "empty ids", body: `{"ids":[]}`, wantError: "ids is required"},
		{name: "invalid id", body: `{"ids":["not-a-uuid"]}`, wantError: "invalid schema id"},
		{name: "too many ids", body: string(tooManyIDsBody), wantError: "ids must contain at most 100 values"},
		{name: "oversized body", body: `{"ids":["` + uuid.NewString() + `"],"padding":"` + strings.Repeat("a", 2<<20) + `"}`, wantError: "request body too large"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := newTestRequest(http.MethodDelete, "/api/ocr/schemas"+tc.query, strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
			}
			var got ErrorResponse
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Fatalf("decode error: %v", err)
			}
			if got.Error != tc.wantError {
				t.Fatalf("error = %q, want %q", got.Error, tc.wantError)
			}
		})
	}
}
