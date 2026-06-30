package api

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"ai.ro/syncra/internal/ocr"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func multipartRequest(t *testing.T, fields map[string]string, filename string, content []byte) *http.Request {
	t.Helper()
	return multipartRequestForPath(t, "/api/ocr", fields, filename, content)
}

func multipartRequestForPath(t *testing.T, path string, fields map[string]string, filename string, content []byte) *http.Request {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			t.Fatalf("write field: %v", err)
		}
	}
	if filename != "" {
		part, err := writer.CreateFormFile("file", filename)
		if err != nil {
			t.Fatalf("create form file: %v", err)
		}
		if _, err := part.Write(content); err != nil {
			t.Fatalf("write file: %v", err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}
	req := newTestRequest(http.MethodPost, path, &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func testRouterWithMaxUploadBytes(t *testing.T, maxUploadBytes int64) (*gin.Engine, *gorm.DB) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	db := apiPostgresTx(t)
	h := &Handler{DB: db, MaxUploadBytes: maxUploadBytes, InternalAPIToken: testInternalAPIToken}
	return NewRouter(h), db
}

func TestOCRRequiresFile(t *testing.T) {
	router, _ := testRouter(t)
	req := multipartRequest(t, nil, "", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestOCRRejectsUnsupportedMime(t *testing.T) {
	router, _ := testRouter(t)
	req := multipartRequest(t, nil, "notes.txt", []byte("hello"))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestOCRRejectsFilenameOver255CharactersBeforeProcessing(t *testing.T) {
	_, db := testRouter(t)
	called := false
	router := routerForFakeOCR(t, db, func(ctx context.Context, input OCRProcessInput) (*MistralOCRResponse, []byte, error) {
		called = true
		return nil, nil, nil
	})

	req := multipartRequest(t, nil, strings.Repeat("a", 256)+".png", validPNGBytes())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Error != "filename must be at most 255 characters" {
		t.Fatalf("error = %q", got.Error)
	}
	if called {
		t.Fatal("OCR processor was called")
	}
}

func TestOCRRejectsInvalidUserID(t *testing.T) {
	router, _ := testRouter(t)
	req := multipartRequest(t, map[string]string{"user_id": "not-a-uuid"}, "scan.png", validPNGBytes())
	w := httptest.NewRecorder()
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

func TestOCRRejectsUnknownUserID(t *testing.T) {
	router, _ := testRouter(t)
	req := multipartRequest(t, map[string]string{"user_id": uuid.NewString()}, "scan.png", validPNGBytes())
	w := httptest.NewRecorder()
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

func TestOCRRejectsBothSchemaSources(t *testing.T) {
	router, db := testRouter(t)
	schema := ocr.ExtractionSchema{Name: "invoice", SchemaJSON: []byte(`{"type":"object"}`), Strict: true}
	if err := db.Create(&schema).Error; err != nil {
		t.Fatalf("create schema: %v", err)
	}
	fields := map[string]string{
		"schema_id": schema.ID.String(),
		"schema":    `{"type":"object"}`,
	}
	req := multipartRequest(t, fields, "scan.png", validPNGBytes())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestOCRRejectsMissingSchemaID(t *testing.T) {
	router, _ := testRouter(t)
	fields := map[string]string{"schema_id": uuid.NewString()}
	req := multipartRequest(t, fields, "scan.png", validPNGBytes())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestOCRRejectsInvalidInlineSchema(t *testing.T) {
	router, _ := testRouter(t)
	fields := map[string]string{"schema": `["not-object"]`}
	req := multipartRequest(t, fields, "scan.png", validPNGBytes())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestOCRRejectsRequestBodyOverLimitBeforeMultipartParsing(t *testing.T) {
	file := validPNGBytes()
	router, _ := testRouterWithMaxUploadBytes(t, int64(len(file)))
	fields := map[string]string{"padding": strings.Repeat("x", 70<<10)}
	req := multipartRequest(t, fields, "scan.png", file)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestOCRAcceptsValidPNGWithoutSchema(t *testing.T) {
	router, _ := testRouter(t)
	file := validPNGBytes()
	req := multipartRequest(t, nil, "scan.png", file)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeOCRValidationResponse(t, w)
	assertOCRFileMetadata(t, got, "scan.png", "image/png", file)
	if got.HasInlineSchema {
		t.Fatalf("has_inline_schema = true, want false")
	}
	if got.SchemaID != nil {
		t.Fatalf("schema_id = %v, want nil", *got.SchemaID)
	}
}

func TestOCRAcceptsValidPDFWithInlineSchema(t *testing.T) {
	router, _ := testRouter(t)
	file := validPDFBytes()
	fields := map[string]string{"schema": `{"type":"object","properties":{"total":{"type":"number"}}}`}
	req := multipartRequest(t, fields, "invoice.pdf", file)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeOCRValidationResponse(t, w)
	assertOCRFileMetadata(t, got, "invoice.pdf", "application/pdf", file)
	if !got.HasInlineSchema {
		t.Fatalf("has_inline_schema = false, want true")
	}
	if got.SchemaID != nil {
		t.Fatalf("schema_id = %v, want nil", *got.SchemaID)
	}
}

func TestOCRAcceptsValidPNGWithSavedSchemaID(t *testing.T) {
	router, db := testRouter(t)
	schema := ocr.ExtractionSchema{Name: "invoice", SchemaJSON: []byte(`{"type":"object"}`), Strict: true}
	if err := db.Create(&schema).Error; err != nil {
		t.Fatalf("create schema: %v", err)
	}
	file := validPNGBytes()
	fields := map[string]string{"schema_id": schema.ID.String()}
	req := multipartRequest(t, fields, "scan.png", file)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeOCRValidationResponse(t, w)
	assertOCRFileMetadata(t, got, "scan.png", "image/png", file)
	if got.HasInlineSchema {
		t.Fatalf("has_inline_schema = true, want false")
	}
	if got.SchemaID == nil || *got.SchemaID != schema.ID {
		t.Fatalf("schema_id = %v, want %s", got.SchemaID, schema.ID)
	}
}

func decodeOCRValidationResponse(t *testing.T, w *httptest.ResponseRecorder) OCRValidationResponse {
	t.Helper()
	var got OCRValidationResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return got
}

func assertOCRFileMetadata(t *testing.T, got OCRValidationResponse, filename string, mimeType string, content []byte) {
	t.Helper()
	if got.OriginalFilename != filename {
		t.Fatalf("original_filename = %q, want %q", got.OriginalFilename, filename)
	}
	if got.MimeType != mimeType {
		t.Fatalf("mime_type = %q, want %q", got.MimeType, mimeType)
	}
	if got.FileSize != int64(len(content)) {
		t.Fatalf("file_size = %d, want %d", got.FileSize, len(content))
	}
	if got.DocumentHash == "" {
		t.Fatal("document_hash is empty")
	}
}

func validPNGBytes() []byte {
	return []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n', 0, 0, 0, 0, 'I', 'H', 'D', 'R'}
}

func validPDFBytes() []byte {
	return []byte("%PDF-1.4\n1 0 obj\n<<>>\nendobj\ntrailer\n<<>>\n%%EOF\n")
}
