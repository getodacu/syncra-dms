package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"ai.ro/syncra/internal/ocr"
)

func TestOCRSuccessPersistsDocument(t *testing.T) {
	_, db := testRouter(t)
	file := validPNGBytes()
	inlineSchema := `{"type":"object","properties":{"Furnizor":{"type":"string"}}}`
	raw := []byte(`{"pages":[{"index":0,"markdown":"# Hello\n\n![img-0.jpeg](img-0.jpeg)","images":[{"id":"img-0.jpeg","image_base64":"data:image/jpeg;base64,ZmFrZQ=="}]},{"index":1,"markdown":"Second page"}],"document_annotation":"{\"Furnizor\":\"Acme\"}","model":"mistral-ocr-latest"}`)
	handler := routerForFakeOCR(t, db, func(ctx context.Context, input OCRProcessInput) (*MistralOCRResponse, []byte, error) {
		if input.MimeType != "image/png" {
			t.Fatalf("MimeType = %q", input.MimeType)
		}
		if !strings.HasPrefix(input.DataURL, "data:image/png;base64,") {
			t.Fatalf("DataURL prefix = %q", input.DataURL)
		}
		resp := &MistralOCRResponse{
			Pages: []MistralOCRPage{
				{
					Index:    0,
					Markdown: "# Hello\n\n![img-0.jpeg](img-0.jpeg)",
					Images:   json.RawMessage(`[{"id":"img-0.jpeg","image_base64":"data:image/jpeg;base64,ZmFrZQ=="}]`),
				},
				{Index: 1, Markdown: "Second page"},
			},
			Model:              "mistral-ocr-latest",
			DocumentAnnotation: ptrString(`{"Furnizor":"Acme"}`),
		}
		return resp, raw, nil
	})

	req := multipartRequest(t, map[string]string{"schema": inlineSchema}, "scan.png", file)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got OCRDocumentResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	var rawResponseFields map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &rawResponseFields); err != nil {
		t.Fatalf("decode raw response: %v", err)
	}
	if _, ok := rawResponseFields["status"]; ok {
		t.Fatalf("response includes document status: %s", w.Body.String())
	}
	if _, ok := rawResponseFields["error_message"]; ok {
		t.Fatalf("response includes document error_message: %s", w.Body.String())
	}
	wantMarkdown := "# Hello\n\n![img-0.jpeg](data:image/jpeg;base64,ZmFrZQ==)\n\nSecond page"
	if got.ID == uuid.Nil || got.Markdown != wantMarkdown {
		t.Fatalf("unexpected response: %#v", got)
	}
	if got.PageCount != 2 {
		t.Fatalf("response page_count = %d, want 2", got.PageCount)
	}
	assertOCRFileMetadata(t, OCRValidationResponse{
		OriginalFilename: got.OriginalFilename,
		MimeType:         got.MimeType,
		FileSize:         got.FileSize,
		DocumentHash:     got.DocumentHash,
	}, "scan.png", "image/png", file)
	if got.Cached {
		t.Fatal("Cached = true, want false")
	}
	if !got.HasInlineSchema {
		t.Fatalf("HasInlineSchema = false, want true")
	}
	if got.SchemaID != nil {
		t.Fatalf("SchemaID = %v, want nil", *got.SchemaID)
	}
	assertJSONEqual(t, got.AnnotationJSON, `{"Furnizor":"Acme"}`)

	var count int64
	if err := db.Model(&ocr.OCRDocument{}).Count(&count).Error; err != nil {
		t.Fatalf("count: %v", err)
	}
	if count != 1 {
		t.Fatalf("document count = %d", count)
	}

	var doc ocr.OCRDocument
	if err := db.First(&doc, "id = ?", got.ID).Error; err != nil {
		t.Fatalf("load document: %v", err)
	}
	if doc.OriginalFilename != "scan.png" || doc.MimeType != "image/png" || doc.FileSize != int64(len(file)) {
		t.Fatalf("unexpected saved metadata: %#v", doc)
	}
	if doc.PageCount != 2 {
		t.Fatalf("stored page_count = %d, want 2", doc.PageCount)
	}
	wantHash, err := computeDocumentHash(file, json.RawMessage(inlineSchema), true)
	if err != nil {
		t.Fatalf("compute document hash: %v", err)
	}
	if doc.DocumentHash != wantHash || got.DocumentHash != wantHash {
		t.Fatalf("DocumentHash response=%q stored=%q want %q", got.DocumentHash, doc.DocumentHash, wantHash)
	}
	if doc.SchemaID != nil {
		t.Fatalf("saved SchemaID = %v, want nil", *doc.SchemaID)
	}
	assertJSONEqual(t, json.RawMessage(doc.InlineSchemaJSON), inlineSchema)
	assertJSONEqual(t, json.RawMessage(doc.AnnotationJSON), `{"Furnizor":"Acme"}`)
	assertJSONEqual(t, json.RawMessage(doc.RawResponseJSON), string(raw))
	if doc.Markdown != wantMarkdown {
		t.Fatalf("unexpected saved OCR markdown: %q", doc.Markdown)
	}
}

func TestOCRSuccessPersistsUserID(t *testing.T) {
	_, db := testRouter(t)
	user := createTestUser(t, db, "ocr-owner@example.com")
	file := validPNGBytes()
	router := routerForFakeOCR(t, db, successfulTestOCRProcessor())

	req := multipartRequest(t, map[string]string{"user_id": user.ID}, "scan.png", file)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got OCRDocumentResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.UserID == nil || string(*got.UserID) != user.ID {
		t.Fatalf("response user_id = %#v, want %s", got.UserID, user.ID)
	}

	var doc ocr.OCRDocument
	if err := db.First(&doc, "id = ?", got.ID).Error; err != nil {
		t.Fatalf("load document: %v", err)
	}
	if doc.UserID == nil || *doc.UserID != user.ID {
		t.Fatalf("stored user_id = %#v, want %s", doc.UserID, user.ID)
	}
}

func TestOCRWithSavedSchemaLinksDocumentToCollection(t *testing.T) {
	_, db := testRouter(t)
	user := createTestUser(t, db, "ocr-saved-schema-collection@example.com")
	schema := ocr.ExtractionSchema{
		UserID:     &user.ID,
		Name:       "Invoice",
		SchemaJSON: datatypes.JSON([]byte(`{"type":"object","properties":{"ok":{"type":"boolean"}}}`)),
		Strict:     true,
	}
	if err := db.Create(&schema).Error; err != nil {
		t.Fatalf("create schema: %v", err)
	}
	collection := createCollection(t, db, user.ID, "Invoices", time.Now())
	createCollectionSchemas(t, db, collection.ID, schema.ID)
	router := routerForFakeOCR(t, db, successfulTestOCRProcessor())

	req := multipartRequest(t, map[string]string{"user_id": user.ID, "schema_id": schema.ID.String()}, "scan.png", validPNGBytes())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got OCRDocumentResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.SchemaID == nil || *got.SchemaID != schema.ID {
		t.Fatalf("schema_id = %#v, want %s", got.SchemaID, schema.ID)
	}
	assertCollectionDocumentLinkCount(t, db, collection.ID, got.ID, 1)
}

func TestOCRWithInlineSchemaDoesNotLinkDocumentToCollection(t *testing.T) {
	_, db := testRouter(t)
	user := createTestUser(t, db, "ocr-inline-schema-collection@example.com")
	schema := ocr.ExtractionSchema{
		UserID:     &user.ID,
		Name:       "Invoice",
		SchemaJSON: datatypes.JSON([]byte(`{"type":"object","properties":{"ok":{"type":"boolean"}}}`)),
		Strict:     true,
	}
	if err := db.Create(&schema).Error; err != nil {
		t.Fatalf("create schema: %v", err)
	}
	collection := createCollection(t, db, user.ID, "Invoices", time.Now())
	createCollectionSchemas(t, db, collection.ID, schema.ID)
	router := routerForFakeOCR(t, db, successfulTestOCRProcessor())

	req := multipartRequest(t, map[string]string{
		"user_id": user.ID,
		"schema":  string(schema.SchemaJSON),
	}, "scan.png", validPNGBytes())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got OCRDocumentResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.SchemaID != nil {
		t.Fatalf("schema_id = %#v, want nil", got.SchemaID)
	}
	if !got.HasInlineSchema {
		t.Fatal("has_inline_schema = false, want true")
	}
	assertCollectionDocumentLinkCount(t, db, collection.ID, got.ID, 0)
}

func TestOCRCachedDocumentBackfillsCollectionLink(t *testing.T) {
	_, db := testRouter(t)
	user := createTestUser(t, db, "ocr-cached-schema-collection@example.com")
	schema := ocr.ExtractionSchema{
		UserID:     &user.ID,
		Name:       "Invoice",
		SchemaJSON: datatypes.JSON([]byte(`{"type":"object","properties":{"ok":{"type":"boolean"}}}`)),
		Strict:     true,
	}
	if err := db.Create(&schema).Error; err != nil {
		t.Fatalf("create schema: %v", err)
	}
	collection := createCollection(t, db, user.ID, "Invoices", time.Now())
	createCollectionSchemas(t, db, collection.ID, schema.ID)
	file := validPNGBytes()
	documentHash, err := computeDocumentHash(file, json.RawMessage(schema.SchemaJSON), schema.Strict)
	if err != nil {
		t.Fatalf("compute document hash: %v", err)
	}
	cached := createStoredOCRDocument(t, db, ocr.OCRDocument{
		UserID:       &user.ID,
		SchemaID:     &schema.ID,
		DocumentHash: documentHash,
	})
	called := false
	router := routerForFakeOCR(t, db, func(ctx context.Context, input OCRProcessInput) (*MistralOCRResponse, []byte, error) {
		called = true
		return nil, nil, nil
	})

	req := multipartRequest(t, map[string]string{"user_id": user.ID, "schema_id": schema.ID.String()}, "scan.png", file)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if called {
		t.Fatal("OCR processor was called")
	}
	var got OCRDocumentResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ID != cached.ID {
		t.Fatalf("id = %s, want cached document %s", got.ID, cached.ID)
	}
	if !got.Cached {
		t.Fatal("cached = false, want true")
	}
	assertCollectionDocumentLinkCount(t, db, collection.ID, cached.ID, 1)
}

func TestOCRCachedInlineSchemaDoesNotBackfillCollectionLink(t *testing.T) {
	_, db := testRouter(t)
	user := createTestUser(t, db, "ocr-cached-inline-schema-collection@example.com")
	schema := ocr.ExtractionSchema{
		UserID:     &user.ID,
		Name:       "Invoice",
		SchemaJSON: datatypes.JSON([]byte(`{"type":"object","properties":{"ok":{"type":"boolean"}}}`)),
		Strict:     true,
	}
	if err := db.Create(&schema).Error; err != nil {
		t.Fatalf("create schema: %v", err)
	}
	collection := createCollection(t, db, user.ID, "Invoices", time.Now())
	createCollectionSchemas(t, db, collection.ID, schema.ID)
	file := validPNGBytes()
	documentHash, err := computeDocumentHash(file, json.RawMessage(schema.SchemaJSON), schema.Strict)
	if err != nil {
		t.Fatalf("compute document hash: %v", err)
	}
	cached := createStoredOCRDocument(t, db, ocr.OCRDocument{
		UserID:       &user.ID,
		SchemaID:     &schema.ID,
		DocumentHash: documentHash,
	})
	called := false
	router := routerForFakeOCR(t, db, func(ctx context.Context, input OCRProcessInput) (*MistralOCRResponse, []byte, error) {
		called = true
		return nil, nil, nil
	})

	req := multipartRequest(t, map[string]string{
		"user_id": user.ID,
		"schema":  string(schema.SchemaJSON),
	}, "scan.png", file)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if called {
		t.Fatal("OCR processor was called")
	}
	var got OCRDocumentResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ID != cached.ID {
		t.Fatalf("id = %s, want cached document %s", got.ID, cached.ID)
	}
	if !got.Cached {
		t.Fatal("cached = false, want true")
	}
	assertCollectionDocumentLinkCount(t, db, collection.ID, cached.ID, 0)
}

func TestOCRCachedSavedSchemaDoesNotBackfillInlineOriginDocument(t *testing.T) {
	_, db := testRouter(t)
	user := createTestUser(t, db, "ocr-cached-saved-schema-inline-origin@example.com")
	schema := ocr.ExtractionSchema{
		UserID:     &user.ID,
		Name:       "Invoice",
		SchemaJSON: datatypes.JSON([]byte(`{"type":"object","properties":{"ok":{"type":"boolean"}}}`)),
		Strict:     true,
	}
	if err := db.Create(&schema).Error; err != nil {
		t.Fatalf("create schema: %v", err)
	}
	collection := createCollection(t, db, user.ID, "Invoices", time.Now())
	createCollectionSchemas(t, db, collection.ID, schema.ID)
	file := validPNGBytes()
	documentHash, err := computeDocumentHash(file, json.RawMessage(schema.SchemaJSON), schema.Strict)
	if err != nil {
		t.Fatalf("compute document hash: %v", err)
	}
	cached := createStoredOCRDocument(t, db, ocr.OCRDocument{
		UserID:           &user.ID,
		InlineSchemaJSON: schema.SchemaJSON,
		DocumentHash:     documentHash,
	})
	called := false
	router := routerForFakeOCR(t, db, func(ctx context.Context, input OCRProcessInput) (*MistralOCRResponse, []byte, error) {
		called = true
		return nil, nil, nil
	})

	req := multipartRequest(t, map[string]string{"user_id": user.ID, "schema_id": schema.ID.String()}, "scan.png", file)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if called {
		t.Fatal("OCR processor was called")
	}
	var got OCRDocumentResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ID != cached.ID {
		t.Fatalf("id = %s, want cached document %s", got.ID, cached.ID)
	}
	if !got.Cached {
		t.Fatal("cached = false, want true")
	}
	assertCollectionDocumentLinkCount(t, db, collection.ID, cached.ID, 0)
}

func TestOCRCachedSavedSchemaBackfillsRequestedSchemaCollection(t *testing.T) {
	_, db := testRouter(t)
	user := createTestUser(t, db, "ocr-cached-requested-schema-collection@example.com")
	schemaJSON := datatypes.JSON([]byte(`{"type":"object","properties":{"ok":{"type":"boolean"}}}`))
	schemaA := ocr.ExtractionSchema{
		UserID:     &user.ID,
		Name:       "Invoice A",
		SchemaJSON: schemaJSON,
		Strict:     true,
	}
	if err := db.Create(&schemaA).Error; err != nil {
		t.Fatalf("create schema A: %v", err)
	}
	schemaB := ocr.ExtractionSchema{
		UserID:     &user.ID,
		Name:       "Invoice B",
		SchemaJSON: schemaJSON,
		Strict:     true,
	}
	if err := db.Create(&schemaB).Error; err != nil {
		t.Fatalf("create schema B: %v", err)
	}
	collectionA := createCollection(t, db, user.ID, "Invoices A", time.Now())
	createCollectionSchemas(t, db, collectionA.ID, schemaA.ID)
	collectionB := createCollection(t, db, user.ID, "Invoices B", time.Now())
	createCollectionSchemas(t, db, collectionB.ID, schemaB.ID)
	file := validPNGBytes()
	documentHash, err := computeDocumentHash(file, json.RawMessage(schemaJSON), true)
	if err != nil {
		t.Fatalf("compute document hash: %v", err)
	}
	cached := createStoredOCRDocument(t, db, ocr.OCRDocument{
		UserID:       &user.ID,
		SchemaID:     &schemaA.ID,
		DocumentHash: documentHash,
	})
	called := false
	router := routerForFakeOCR(t, db, func(ctx context.Context, input OCRProcessInput) (*MistralOCRResponse, []byte, error) {
		called = true
		return nil, nil, nil
	})

	req := multipartRequest(t, map[string]string{"user_id": user.ID, "schema_id": schemaB.ID.String()}, "scan.png", file)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if called {
		t.Fatal("OCR processor was called")
	}
	var got OCRDocumentResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ID != cached.ID {
		t.Fatalf("id = %s, want cached document %s", got.ID, cached.ID)
	}
	if !got.Cached {
		t.Fatal("cached = false, want true")
	}
	assertCollectionDocumentLinkCount(t, db, collectionA.ID, cached.ID, 0)
	assertCollectionDocumentLinkCount(t, db, collectionB.ID, cached.ID, 1)
}

func TestGetOCRDocumentReturnsStoredDocument(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "ocr-document-owner@example.com")
	schema := ocr.ExtractionSchema{
		UserID:     &user.ID,
		Name:       "Invoice",
		SchemaJSON: datatypes.JSON([]byte(`{"type":"object"}`)),
		Strict:     true,
	}
	if err := db.Create(&schema).Error; err != nil {
		t.Fatalf("create schema: %v", err)
	}
	doc := createStoredOCRDocument(t, db, ocr.OCRDocument{
		UserID:           &user.ID,
		OriginalFilename: "stored.png",
		MimeType:         "image/png",
		FileSize:         123,
		DocumentHash:     "stored-document-hash",
		SchemaID:         &schema.ID,
		Markdown:         "# Stored",
		AnnotationJSON:   datatypes.JSON([]byte(`{"total":10}`)),
	})

	req := newTestRequest(http.MethodGet, "/api/ocr/document/"+doc.ID.String()+"?user_id="+user.ID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got OCRDocumentResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ID != doc.ID || got.OriginalFilename != "stored.png" || got.Markdown != "# Stored" {
		t.Fatalf("unexpected response: %#v", got)
	}
	if got.UserID == nil || string(*got.UserID) != user.ID {
		t.Fatalf("user_id = %#v, want %s", got.UserID, user.ID)
	}
	if got.SchemaID == nil || *got.SchemaID != schema.ID {
		t.Fatalf("schema_id = %#v, want %s", got.SchemaID, schema.ID)
	}
	if got.MimeType != "image/png" || got.FileSize != 123 || got.DocumentHash != "stored-document-hash" {
		t.Fatalf("unexpected metadata: %#v", got)
	}
	if got.PageCount != 1 {
		t.Fatalf("page_count = %d, want 1", got.PageCount)
	}
	if got.Cached {
		t.Fatal("cached = true, want false")
	}
	if got.HasInlineSchema {
		t.Fatal("has_inline_schema = true, want false")
	}
	assertJSONEqual(t, got.AnnotationJSON, `{"total":10}`)
}

func TestGetOCRDocumentScopesByUserIDQuery(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "ocr-document-scope@example.com")
	other := createTestUser(t, db, "ocr-document-scope-other@example.com")
	doc := createStoredOCRDocument(t, db, ocr.OCRDocument{UserID: &user.ID})

	req := newTestRequest(http.MethodGet, "/api/ocr/document/"+doc.ID.String()+"?user_id="+other.ID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestGetOCRDocumentRejectsInvalidUserIDQuery(t *testing.T) {
	router, db := testRouter(t)
	doc := createStoredOCRDocument(t, db, ocr.OCRDocument{})

	req := newTestRequest(http.MethodGet, "/api/ocr/document/"+doc.ID.String()+"?user_id=not-a-uuid", nil)
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

func TestGetOCRDocumentRejectsInvalidID(t *testing.T) {
	router, _ := testRouter(t)
	req := newTestRequest(http.MethodGet, "/api/ocr/document/not-a-uuid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Error != "invalid OCR document id" {
		t.Fatalf("error = %q", got.Error)
	}
}

func TestGetOCRDocumentReturnsNotFound(t *testing.T) {
	router, _ := testRouter(t)
	req := newTestRequest(http.MethodGet, "/api/ocr/document/"+uuid.NewString(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Error != "OCR document not found" {
		t.Fatalf("error = %q", got.Error)
	}
}

func TestOCRCacheHitReturnsRecentDocumentWithoutProcessing(t *testing.T) {
	_, db := testRouter(t)
	file := validPNGBytes()
	documentHash, err := computeDocumentHash(file, nil, false)
	if err != nil {
		t.Fatalf("compute document hash: %v", err)
	}
	cached := ocr.OCRDocument{
		OriginalFilename: "cached.png",
		MimeType:         "image/png",
		FileSize:         int64(len(file)),
		DocumentHash:     documentHash,
		Markdown:         "# Cached",
		RawResponseJSON:  datatypes.JSON([]byte(`{"pages":[{"index":0,"markdown":"# Cached"}]}`)),
	}
	if err := db.Create(&cached).Error; err != nil {
		t.Fatalf("create cached document: %v", err)
	}
	called := false
	router := routerForFakeOCR(t, db, func(ctx context.Context, input OCRProcessInput) (*MistralOCRResponse, []byte, error) {
		called = true
		return nil, nil, nil
	})

	req := multipartRequest(t, nil, "scan.png", file)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if called {
		t.Fatal("OCR processor was called")
	}
	var got OCRDocumentResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ID != cached.ID || got.Markdown != "# Cached" || !got.Cached {
		t.Fatalf("unexpected cached response: %#v", got)
	}

	var count int64
	if err := db.Model(&ocr.OCRDocument{}).Count(&count).Error; err != nil {
		t.Fatalf("count: %v", err)
	}
	if count != 1 {
		t.Fatalf("document count = %d, want 1", count)
	}
}

func TestOCRCacheHitIsScopedByUserID(t *testing.T) {
	_, db := testRouter(t)
	userOne := createTestUser(t, db, "ocr-cache-one@example.com")
	userTwo := createTestUser(t, db, "ocr-cache-two@example.com")
	file := validPNGBytes()
	documentHash, err := computeDocumentHash(file, nil, false)
	if err != nil {
		t.Fatalf("compute document hash: %v", err)
	}
	olderRequestedCache := time.Now().Add(-1 * time.Hour)
	newerWrongUserCache := time.Now()
	cachedOne := ocr.OCRDocument{
		UserID:           &userOne.ID,
		CreatedAt:        newerWrongUserCache,
		UpdatedAt:        newerWrongUserCache,
		OriginalFilename: "user-one.png",
		MimeType:         "image/png",
		FileSize:         int64(len(file)),
		DocumentHash:     documentHash,
		Markdown:         "# User One",
		RawResponseJSON:  datatypes.JSON([]byte(`{"pages":[{"index":0,"markdown":"# User One"}]}`)),
	}
	if err := db.Create(&cachedOne).Error; err != nil {
		t.Fatalf("create first cached document: %v", err)
	}
	cachedTwo := ocr.OCRDocument{
		UserID:           &userTwo.ID,
		CreatedAt:        olderRequestedCache,
		UpdatedAt:        olderRequestedCache,
		OriginalFilename: "user-two.png",
		MimeType:         "image/png",
		FileSize:         int64(len(file)),
		DocumentHash:     documentHash,
		Markdown:         "# User Two",
		RawResponseJSON:  datatypes.JSON([]byte(`{"pages":[{"index":0,"markdown":"# User Two"}]}`)),
	}
	if err := db.Create(&cachedTwo).Error; err != nil {
		t.Fatalf("create second cached document: %v", err)
	}
	called := false
	router := routerForFakeOCR(t, db, func(ctx context.Context, input OCRProcessInput) (*MistralOCRResponse, []byte, error) {
		called = true
		return nil, nil, nil
	})

	req := multipartRequest(t, map[string]string{"user_id": userTwo.ID}, "scan.png", file)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if called {
		t.Fatal("OCR processor was called")
	}
	var got OCRDocumentResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ID != cachedTwo.ID || got.Markdown != "# User Two" || !got.Cached {
		t.Fatalf("unexpected cached response: %#v", got)
	}
	if got.UserID == nil || string(*got.UserID) != userTwo.ID {
		t.Fatalf("response user_id = %#v, want %s", got.UserID, userTwo.ID)
	}
}

func TestOCRSystemWideCacheIgnoresUserOwnedDocuments(t *testing.T) {
	_, db := testRouter(t)
	user := createTestUser(t, db, "ocr-cache-owner@example.com")
	file := validPNGBytes()
	documentHash, err := computeDocumentHash(file, nil, false)
	if err != nil {
		t.Fatalf("compute document hash: %v", err)
	}
	owned := ocr.OCRDocument{
		UserID:           &user.ID,
		OriginalFilename: "owned.png",
		MimeType:         "image/png",
		FileSize:         int64(len(file)),
		DocumentHash:     documentHash,
		Markdown:         "# Owned",
		RawResponseJSON:  datatypes.JSON([]byte(`{"pages":[{"index":0,"markdown":"# Owned"}]}`)),
	}
	if err := db.Create(&owned).Error; err != nil {
		t.Fatalf("create owned cached document: %v", err)
	}
	calls := 0
	router := routerForFakeOCR(t, db, func(ctx context.Context, input OCRProcessInput) (*MistralOCRResponse, []byte, error) {
		calls++
		return &MistralOCRResponse{
			Pages: []MistralOCRPage{{Index: 0, Markdown: "# System Wide"}},
			Model: "mistral-ocr-latest",
		}, []byte(`{"pages":[{"index":0,"markdown":"# System Wide"}],"model":"mistral-ocr-latest"}`), nil
	})

	req := multipartRequest(t, nil, "scan.png", file)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if calls != 1 {
		t.Fatalf("OCR calls = %d, want 1", calls)
	}
	var got OCRDocumentResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Cached {
		t.Fatal("Cached = true, want false")
	}
	if got.UserID != nil {
		t.Fatalf("response user_id = %#v, want nil", got.UserID)
	}

	var doc ocr.OCRDocument
	if err := db.First(&doc, "id = ?", got.ID).Error; err != nil {
		t.Fatalf("load document: %v", err)
	}
	if doc.UserID != nil {
		t.Fatalf("stored user_id = %#v, want nil", doc.UserID)
	}
	var count int64
	if err := db.Model(&ocr.OCRDocument{}).Count(&count).Error; err != nil {
		t.Fatalf("count: %v", err)
	}
	if count != 2 {
		t.Fatalf("document count = %d, want 2", count)
	}
}

func TestOCRCacheMissesExpiredDocument(t *testing.T) {
	_, db := testRouter(t)
	file := validPNGBytes()
	documentHash, err := computeDocumentHash(file, nil, false)
	if err != nil {
		t.Fatalf("compute document hash: %v", err)
	}
	expired := ocr.OCRDocument{
		CreatedAt:        time.Now().Add(-25 * time.Hour),
		UpdatedAt:        time.Now().Add(-25 * time.Hour),
		OriginalFilename: "expired.png",
		MimeType:         "image/png",
		FileSize:         int64(len(file)),
		DocumentHash:     documentHash,
		Markdown:         "# Expired",
		RawResponseJSON:  datatypes.JSON([]byte(`{"pages":[{"index":0,"markdown":"# Expired"}]}`)),
	}
	if err := db.Create(&expired).Error; err != nil {
		t.Fatalf("create expired document: %v", err)
	}
	calls := 0
	router := routerForFakeOCR(t, db, func(ctx context.Context, input OCRProcessInput) (*MistralOCRResponse, []byte, error) {
		calls++
		return &MistralOCRResponse{
			Pages: []MistralOCRPage{{Index: 0, Markdown: "# Fresh"}},
			Model: "mistral-ocr-latest",
		}, []byte(`{"pages":[{"index":0,"markdown":"# Fresh"}],"model":"mistral-ocr-latest"}`), nil
	})

	req := multipartRequest(t, nil, "scan.png", file)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if calls != 1 {
		t.Fatalf("OCR calls = %d, want 1", calls)
	}
	var got OCRDocumentResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Cached || got.Markdown != "# Fresh" {
		t.Fatalf("unexpected response: %#v", got)
	}

	var count int64
	if err := db.Model(&ocr.OCRDocument{}).Count(&count).Error; err != nil {
		t.Fatalf("count: %v", err)
	}
	if count != 2 {
		t.Fatalf("document count = %d, want 2", count)
	}
}

func TestDocumentHashUsesCanonicalSchemaContentAndStrictness(t *testing.T) {
	file := validPNGBytes()
	inline, err := computeDocumentHash(file, json.RawMessage(`{"properties":{"Furnizor":{"type":"string"}},"type":"object"}`), true)
	if err != nil {
		t.Fatalf("compute inline hash: %v", err)
	}
	saved, err := computeDocumentHash(file, json.RawMessage(`{
		"type": "object",
		"properties": {"Furnizor": {"type": "string"}}
	}`), true)
	if err != nil {
		t.Fatalf("compute saved hash: %v", err)
	}
	if inline != saved {
		t.Fatalf("canonical schema hashes differ: inline=%q saved=%q", inline, saved)
	}
	nonStrict, err := computeDocumentHash(file, json.RawMessage(`{"type":"object","properties":{"Furnizor":{"type":"string"}}}`), false)
	if err != nil {
		t.Fatalf("compute non-strict hash: %v", err)
	}
	if nonStrict == inline {
		t.Fatal("strictness did not affect document hash")
	}
	differentSchema, err := computeDocumentHash(file, json.RawMessage(`{"type":"object","properties":{"Total":{"type":"number"}}}`), true)
	if err != nil {
		t.Fatalf("compute different schema hash: %v", err)
	}
	if differentSchema == inline {
		t.Fatal("schema content did not affect document hash")
	}
	noSchema, err := computeDocumentHash(file, nil, false)
	if err != nil {
		t.Fatalf("compute no-schema hash: %v", err)
	}
	if noSchema == inline {
		t.Fatal("no-schema hash matches schema-backed hash")
	}
}

func TestOCRDefaultProcessorSendsMistralPayload(t *testing.T) {
	cases := []struct {
		name         string
		mimeType     string
		wantDocType  string
		wantURLField string
	}{
		{name: "image", mimeType: "image/png", wantDocType: "image_url", wantURLField: "image_url"},
		{name: "pdf", mimeType: "application/pdf", wantDocType: "document_url", wantURLField: "document_url"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			wantRaw := []byte(`{"pages":[{"index":0,"markdown":"# From Mistral"}],"document_annotation":"{\"Furnizor\":\"Acme\"}","model":"mistral-ocr-test"}`)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Fatalf("method = %s", r.Method)
				}
				if r.URL.Path != "/v1/ocr" {
					t.Fatalf("path = %s", r.URL.Path)
				}
				if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
					t.Fatalf("Authorization = %q", got)
				}
				if got := r.Header.Get("Content-Type"); got != "application/json" {
					t.Fatalf("Content-Type = %q", got)
				}

				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatalf("read body: %v", err)
				}
				var payload map[string]any
				if err := json.Unmarshal(body, &payload); err != nil {
					t.Fatalf("decode payload: %v", err)
				}
				if payload["model"] != "mistral-ocr-test" {
					t.Fatalf("model = %v", payload["model"])
				}
				if payload["include_image_base64"] != true {
					t.Fatalf("include_image_base64 = %v", payload["include_image_base64"])
				}
				document, ok := payload["document"].(map[string]any)
				if !ok {
					t.Fatalf("document = %#v", payload["document"])
				}
				if document["type"] != tc.wantDocType {
					t.Fatalf("document type = %v", document["type"])
				}
				if document[tc.wantURLField] != "data:"+tc.mimeType+";base64,ZmFrZQ==" {
					t.Fatalf("%s = %v", tc.wantURLField, document[tc.wantURLField])
				}

				format, ok := payload["document_annotation_format"].(map[string]any)
				if !ok {
					t.Fatalf("document_annotation_format = %#v", payload["document_annotation_format"])
				}
				if format["type"] != "json_schema" {
					t.Fatalf("format type = %v", format["type"])
				}
				schemaPayload, ok := format["json_schema"].(map[string]any)
				if !ok {
					t.Fatalf("json_schema = %#v", format["json_schema"])
				}
				if schemaPayload["name"] != "response_schema" {
					t.Fatalf("schema name = %v", schemaPayload["name"])
				}
				if schemaPayload["strict"] != false {
					t.Fatalf("strict = %v", schemaPayload["strict"])
				}
				assertJSONObjectField(t, schemaPayload["schema"], `{"type":"object","properties":{"Furnizor":{"type":"string"}}}`)

				w.Header().Set("Content-Type", "application/json")
				if _, err := w.Write(wantRaw); err != nil {
					t.Fatalf("write response: %v", err)
				}
			}))
			defer server.Close()

			h := &Handler{
				MistralAPIKey:  "test-key",
				MistralBaseURL: server.URL,
				MistralModel:   "mistral-ocr-test",
			}
			got, raw, err := h.defaultOCRProcessor()(context.Background(), OCRProcessInput{
				MimeType: tc.mimeType,
				DataURL:  "data:" + tc.mimeType + ";base64,ZmFrZQ==",
				Schema:   json.RawMessage(`{"type":"object","properties":{"Furnizor":{"type":"string"}}}`),
				Strict:   false,
			})
			if err != nil {
				t.Fatalf("process OCR: %v", err)
			}
			if !bytes.Equal(raw, wantRaw) {
				t.Fatalf("raw = %s, want %s", raw, wantRaw)
			}
			if got.Model != "mistral-ocr-test" || len(got.Pages) != 1 || got.Pages[0].Markdown != "# From Mistral" {
				t.Fatalf("unexpected response: %#v", got)
			}
			if got.DocumentAnnotation == nil || *got.DocumentAnnotation != `{"Furnizor":"Acme"}` {
				t.Fatalf("DocumentAnnotation = %v", got.DocumentAnnotation)
			}
		})
	}
}

func TestOCRDefaultProcessorSanitizesTransportErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	baseURL := server.URL
	server.Close()

	h := &Handler{
		MistralAPIKey:  "test-key",
		MistralBaseURL: baseURL,
		MistralModel:   "mistral-ocr-test",
	}
	_, _, err := h.defaultOCRProcessor()(context.Background(), OCRProcessInput{
		MimeType: "image/png",
		DataURL:  "data:image/png;base64,ZmFrZQ==",
	})
	if err == nil {
		t.Fatal("err = nil")
	}
	var upstream upstreamError
	if !errors.As(err, &upstream) {
		t.Fatalf("err type = %T, want upstreamError", err)
	}
	if err.Error() != "mistral OCR request failed" {
		t.Fatalf("err = %q", err.Error())
	}
}

func TestOCRDefaultProcessorRejectsOversizedMistralResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"pages":[],"padding":"`))
		_, _ = w.Write([]byte(strings.Repeat("a", 12<<20)))
		_, _ = w.Write([]byte(`"}`))
	}))
	defer server.Close()

	h := &Handler{
		MistralAPIKey:  "test-key",
		MistralBaseURL: server.URL,
		MistralModel:   "mistral-ocr-test",
	}
	_, _, err := h.defaultOCRProcessor()(context.Background(), OCRProcessInput{
		MimeType: "image/png",
		DataURL:  "data:image/png;base64,ZmFrZQ==",
	})
	if err == nil {
		t.Fatal("err = nil")
	}
	var upstream upstreamError
	if !errors.As(err, &upstream) {
		t.Fatalf("err type = %T, want upstreamError", err)
	}
	if err.Error() != "Mistral OCR response too large" {
		t.Fatalf("err = %q", err.Error())
	}
}

func TestOCRErrorDoesNotPersistDocument(t *testing.T) {
	_, db := testRouter(t)
	router := routerForFakeOCR(t, db, func(ctx context.Context, input OCRProcessInput) (*MistralOCRResponse, []byte, error) {
		return nil, nil, errUpstream("mistral failed")
	})

	req := multipartRequest(t, nil, "scan.png", validPNGBytes())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadGateway {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}

	var count int64
	if err := db.Model(&ocr.OCRDocument{}).Count(&count).Error; err != nil {
		t.Fatalf("count: %v", err)
	}
	if count != 0 {
		t.Fatalf("document count = %d", count)
	}
}

func TestOCRInvalidAnnotationDoesNotPersistDocument(t *testing.T) {
	_, db := testRouter(t)
	router := routerForFakeOCR(t, db, func(ctx context.Context, input OCRProcessInput) (*MistralOCRResponse, []byte, error) {
		resp := &MistralOCRResponse{
			Pages:              []MistralOCRPage{{Index: 0, Markdown: "# Hello"}},
			Model:              "mistral-ocr-latest",
			DocumentAnnotation: ptrString(`{"Furnizor":`),
		}
		raw := []byte(`{"pages":[{"index":0,"markdown":"# Hello"}],"document_annotation":"{\"Furnizor\":"}`)
		return resp, raw, nil
	})

	req := multipartRequest(t, nil, "scan.png", validPNGBytes())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadGateway {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}

	var got ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Error != "invalid Mistral document annotation JSON" {
		t.Fatalf("error = %q", got.Error)
	}

	var count int64
	if err := db.Model(&ocr.OCRDocument{}).Count(&count).Error; err != nil {
		t.Fatalf("count: %v", err)
	}
	if count != 0 {
		t.Fatalf("document count = %d", count)
	}
}

func TestOCRSchemaBackedMissingAnnotationDoesNotPersistDocument(t *testing.T) {
	cases := []struct {
		name       string
		annotation *string
	}{
		{name: "nil", annotation: nil},
		{name: "blank", annotation: ptrString(" \t\n ")},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, db := testRouter(t)
			router := routerForFakeOCR(t, db, func(ctx context.Context, input OCRProcessInput) (*MistralOCRResponse, []byte, error) {
				resp := &MistralOCRResponse{
					Pages:              []MistralOCRPage{{Index: 0, Markdown: "# Hello"}},
					Model:              "mistral-ocr-latest",
					DocumentAnnotation: tc.annotation,
				}
				raw := []byte(`{"pages":[{"index":0,"markdown":"# Hello"}],"model":"mistral-ocr-latest"}`)
				return resp, raw, nil
			})

			req := multipartRequest(t, map[string]string{"schema": `{"type":"object"}`}, "scan.png", validPNGBytes())
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			if w.Code != http.StatusBadGateway {
				t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
			}

			var got ErrorResponse
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Fatalf("decode: %v", err)
			}
			if got.Error != "missing Mistral document annotation JSON" {
				t.Fatalf("error = %q", got.Error)
			}

			var count int64
			if err := db.Model(&ocr.OCRDocument{}).Count(&count).Error; err != nil {
				t.Fatalf("count: %v", err)
			}
			if count != 0 {
				t.Fatalf("document count = %d", count)
			}
		})
	}
}

func routerForFakeOCR(t *testing.T, db *gorm.DB, fake OCRProcessor) *gin.Engine {
	t.Helper()
	h := &Handler{DB: db, OCR: fake, MaxUploadBytes: 20 << 20, InternalAPIToken: testInternalAPIToken}
	return NewRouter(h)
}

func createStoredOCRDocument(t *testing.T, db *gorm.DB, doc ocr.OCRDocument) ocr.OCRDocument {
	t.Helper()
	if doc.OriginalFilename == "" {
		doc.OriginalFilename = "scan.png"
	}
	if doc.MimeType == "" {
		doc.MimeType = "image/png"
	}
	if doc.FileSize == 0 {
		doc.FileSize = int64(len(validPNGBytes()))
	}
	if doc.DocumentHash == "" {
		hash, err := computeDocumentHash(validPNGBytes(), nil, false)
		if err != nil {
			t.Fatalf("compute document hash: %v", err)
		}
		doc.DocumentHash = hash
	}
	if doc.Markdown == "" {
		doc.Markdown = "# Stored"
	}
	if len(doc.RawResponseJSON) == 0 {
		doc.RawResponseJSON = datatypes.JSON([]byte(`{"pages":[{"index":0,"markdown":"# Stored"}]}`))
	}
	if err := db.Create(&doc).Error; err != nil {
		t.Fatalf("create OCR document: %v", err)
	}
	return doc
}

func assertCollectionDocumentLinkCount(t *testing.T, db *gorm.DB, collectionID uuid.UUID, documentID uuid.UUID, want int64) {
	t.Helper()
	var count int64
	if err := db.Model(&ocr.CollectionDocument{}).
		Where("collection_id = ? AND document_id = ?", collectionID, documentID).
		Count(&count).Error; err != nil {
		t.Fatalf("count collection document links: %v", err)
	}
	if count != want {
		t.Fatalf("collection document links = %d, want %d", count, want)
	}
}

func assertJSONObjectField(t *testing.T, got any, want string) {
	t.Helper()
	gotBytes, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("marshal JSON field: %v", err)
	}
	assertJSONEqual(t, gotBytes, want)
}

func ptrString(value string) *string {
	return &value
}
