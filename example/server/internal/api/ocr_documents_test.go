package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"ai.ro/syncra/internal/ocr"
)

func createListOCRDocument(t *testing.T, db *gorm.DB, doc ocr.OCRDocument) ocr.OCRDocument {
	t.Helper()
	if doc.OriginalFilename == "" {
		doc.OriginalFilename = "scan.png"
	}
	if doc.MimeType == "" {
		doc.MimeType = "image/png"
	}
	if doc.FileSize == 0 {
		doc.FileSize = 12
	}
	if doc.DocumentHash == "" {
		doc.DocumentHash = uuid.NewString()
	}
	if doc.Markdown == "" {
		doc.Markdown = "# OCR"
	}
	if len(doc.RawResponseJSON) == 0 {
		doc.RawResponseJSON = datatypes.JSON([]byte(`{"pages":[{"index":0}]}`))
	}
	if err := db.Create(&doc).Error; err != nil {
		t.Fatalf("create OCR document: %v", err)
	}
	return doc
}

func setListOCRDocumentCreatedAt(t *testing.T, db *gorm.DB, doc ocr.OCRDocument, createdAt time.Time) ocr.OCRDocument {
	t.Helper()
	createdAt = createdAt.UTC()
	if err := db.Model(&doc).UpdateColumn("created_at", createdAt).Error; err != nil {
		t.Fatalf("set OCR document created_at: %v", err)
	}
	doc.CreatedAt = createdAt
	return doc
}

func decodeOCRDocumentListResponse(t *testing.T, w *httptest.ResponseRecorder) OCRDocumentListResponse {
	t.Helper()
	var got OCRDocumentListResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode OCR document list response: %v body=%s", err, w.Body.String())
	}
	return got
}

func decodeDeleteOCRDocumentsResponse(t *testing.T, w *httptest.ResponseRecorder) DeleteOCRDocumentsResponse {
	t.Helper()
	var got DeleteOCRDocumentsResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode delete OCR documents response: %v body=%s", err, w.Body.String())
	}
	return got
}

func decodeMoveOCRDocumentsToCollectionsResponse(t *testing.T, w *httptest.ResponseRecorder) MoveOCRDocumentsToCollectionsResponse {
	t.Helper()
	var got MoveOCRDocumentsToCollectionsResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode move OCR documents response: %v body=%s", err, w.Body.String())
	}
	return got
}

func assertOCRDocumentListIDs(t *testing.T, documents []OCRDocumentListItemResponse, want ...uuid.UUID) {
	t.Helper()
	if len(documents) != len(want) {
		t.Fatalf("document count = %d, want %d: %#v", len(documents), len(want), documents)
	}
	for i, id := range want {
		if documents[i].ID != id {
			t.Fatalf("document[%d].id = %s, want %s; documents=%#v", i, documents[i].ID, id, documents)
		}
	}
}

func assertOCRDocumentExists(t *testing.T, db *gorm.DB, id uuid.UUID, want bool) {
	t.Helper()
	var count int64
	if err := db.Model(&ocr.OCRDocument{}).Where("id = ?", id).Count(&count).Error; err != nil {
		t.Fatalf("count OCR document %s: %v", id, err)
	}
	if got := count > 0; got != want {
		t.Fatalf("OCR document %s exists = %t, want %t", id, got, want)
	}
}

func TestUpdateOCRDocumentRenamesOwnedDocument(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "document-update-owner@example.com")
	doc := createListOCRDocument(t, db, ocr.OCRDocument{UserID: &user.ID, OriginalFilename: "old-name.pdf"})

	req := newTestRequest(
		http.MethodPatch,
		"/api/ocr/documents/"+doc.ID.String()+"?user_id="+url.QueryEscape(user.ID),
		strings.NewReader(`{"original_filename":"  renamed invoice.pdf  "}`),
	)
	req.Header.Set("content-type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got OCRDocumentResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode OCR document update response: %v body=%s", err, w.Body.String())
	}
	if got.ID != doc.ID || got.OriginalFilename != "renamed invoice.pdf" {
		t.Fatalf("update response = %#v, want renamed document", got)
	}

	var updated ocr.OCRDocument
	if err := db.First(&updated, "id = ?", doc.ID).Error; err != nil {
		t.Fatalf("load updated OCR document: %v", err)
	}
	if updated.OriginalFilename != "renamed invoice.pdf" {
		t.Fatalf("original_filename = %q, want renamed invoice.pdf", updated.OriginalFilename)
	}
}

func TestUpdateOCRDocumentReturnsNotFoundForWrongOwner(t *testing.T) {
	router, db := testRouter(t)
	owner := createTestUser(t, db, "document-update-real-owner@example.com")
	other := createTestUser(t, db, "document-update-other-owner@example.com")
	doc := createListOCRDocument(t, db, ocr.OCRDocument{UserID: &owner.ID, OriginalFilename: "owned.pdf"})

	req := newTestRequest(
		http.MethodPatch,
		"/api/ocr/documents/"+doc.ID.String()+"?user_id="+url.QueryEscape(other.ID),
		strings.NewReader(`{"original_filename":"other.pdf"}`),
	)
	req.Header.Set("content-type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if got.Error != "OCR document not found" {
		t.Fatalf("error = %q, want OCR document not found", got.Error)
	}

	var unchanged ocr.OCRDocument
	if err := db.First(&unchanged, "id = ?", doc.ID).Error; err != nil {
		t.Fatalf("load OCR document: %v", err)
	}
	if unchanged.OriginalFilename != "owned.pdf" {
		t.Fatalf("original_filename = %q, want owned.pdf", unchanged.OriginalFilename)
	}
}

func TestUpdateOCRDocumentRejectsInvalidRequests(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "document-update-invalid@example.com")
	doc := createListOCRDocument(t, db, ocr.OCRDocument{UserID: &user.ID})

	cases := []struct {
		name      string
		path      string
		body      string
		wantError string
	}{
		{
			name:      "invalid id",
			path:      "/api/ocr/documents/not-a-uuid?user_id=" + url.QueryEscape(user.ID),
			body:      `{"original_filename":"renamed.pdf"}`,
			wantError: "invalid OCR document id",
		},
		{
			name:      "malformed json",
			path:      "/api/ocr/documents/" + doc.ID.String() + "?user_id=" + url.QueryEscape(user.ID),
			body:      `{`,
			wantError: "invalid OCR document update request",
		},
		{
			name:      "empty filename",
			path:      "/api/ocr/documents/" + doc.ID.String() + "?user_id=" + url.QueryEscape(user.ID),
			body:      `{"original_filename":"   "}`,
			wantError: "original_filename is required",
		},
		{
			name:      "too long filename",
			path:      "/api/ocr/documents/" + doc.ID.String() + "?user_id=" + url.QueryEscape(user.ID),
			body:      `{"original_filename":"` + strings.Repeat("a", maxOriginalFilenameCharacters+1) + `"}`,
			wantError: "filename must be at most 255 characters",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := newTestRequest(http.MethodPatch, tc.path, strings.NewReader(tc.body))
			req.Header.Set("content-type", "application/json")
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

func TestListOCRDocumentsDefaultsToSystemDocumentsByNewestFirst(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "document-list-owner@example.com")
	base := time.Date(2026, 5, 27, 12, 0, 0, 0, time.UTC)
	older := createListOCRDocument(t, db, ocr.OCRDocument{OriginalFilename: "older.pdf"})
	older = setListOCRDocumentCreatedAt(t, db, older, base.Add(-time.Hour))
	newer := createListOCRDocument(t, db, ocr.OCRDocument{OriginalFilename: "newer.pdf"})
	newer = setListOCRDocumentCreatedAt(t, db, newer, base)
	owned := createListOCRDocument(t, db, ocr.OCRDocument{UserID: &user.ID, OriginalFilename: "owned.pdf"})
	_ = setListOCRDocumentCreatedAt(t, db, owned, base.Add(time.Hour))

	req := newTestRequest(http.MethodGet, "/api/ocr/documents", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeOCRDocumentListResponse(t, w)
	assertOCRDocumentListIDs(t, got.Documents, newer.ID, older.ID)
	if got.NextCursor != nil {
		t.Fatalf("next_cursor = %#v, want nil", got.NextCursor)
	}

	var raw struct {
		Documents []map[string]any `json:"documents"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &raw); err != nil {
		t.Fatalf("decode raw response: %v", err)
	}
	for _, forbidden := range []string{"markdown", "annotation_json", "raw_response_json", "cached"} {
		if _, ok := raw.Documents[0][forbidden]; ok {
			t.Fatalf("list item includes %q: %s", forbidden, w.Body.String())
		}
	}
}

func TestDeleteOCRDocumentSoftDeletesOwnedDocument(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "document-delete-owner@example.com")
	doc := createListOCRDocument(t, db, ocr.OCRDocument{UserID: &user.ID, OriginalFilename: "delete-me.pdf"})
	job := createStoredOCRJob(t, db, ocr.OCRJob{
		UserID:     &user.ID,
		DocumentID: &doc.ID,
		Status:     ocr.OCRJobStatusCompleted,
	})

	req := newTestRequest(
		http.MethodDelete,
		"/api/ocr/documents/"+doc.ID.String()+"?user_id="+url.QueryEscape(user.ID),
		nil,
	)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if w.Body.Len() != 0 {
		t.Fatalf("body = %q, want empty", w.Body.String())
	}
	assertOCRDocumentExists(t, db, doc.ID, false)

	var deleted ocr.OCRDocument
	if err := db.Unscoped().First(&deleted, "id = ?", doc.ID).Error; err != nil {
		t.Fatalf("load soft-deleted OCR document: %v", err)
	}
	if !deleted.DeletedAt.Valid {
		t.Fatal("deleted_at is not set")
	}

	var gotJob ocr.OCRJob
	if err := db.First(&gotJob, "id = ?", job.ID).Error; err != nil {
		t.Fatalf("load OCR job: %v", err)
	}
	if gotJob.DocumentID == nil || *gotJob.DocumentID != doc.ID {
		t.Fatalf("job document_id = %#v, want %s after soft delete", gotJob.DocumentID, doc.ID)
	}
}

func TestDeleteOCRDocumentRemovesCollectionAssociations(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "document-delete-collection-owner@example.com")
	collection := createCollection(t, db, user.ID, "Linked Documents", time.Now())
	doc := createListOCRDocument(t, db, ocr.OCRDocument{UserID: &user.ID, OriginalFilename: "linked.pdf"})
	if err := db.Create(&ocr.CollectionDocument{CollectionID: collection.ID, DocumentID: doc.ID}).Error; err != nil {
		t.Fatalf("create collection document: %v", err)
	}

	req := newTestRequest(
		http.MethodDelete,
		"/api/ocr/documents/"+doc.ID.String()+"?user_id="+url.QueryEscape(user.ID),
		nil,
	)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if w.Body.Len() != 0 {
		t.Fatalf("body = %q, want empty", w.Body.String())
	}
	assertCount(t, db, &ocr.CollectionDocument{}, "collection_id = ? AND document_id = ?", collection.ID, doc.ID, 0)
	assertCount(t, db, &ocr.OCRDocument{}, "id = ?", doc.ID, 0)
	assertCount(t, db.Unscoped(), &ocr.OCRDocument{}, "id = ?", doc.ID, 1)
	assertCount(t, db, &ocr.Collection{}, "id = ?", collection.ID, 1)
}

func TestDeleteOCRDocumentLegacySingularRouteReturnsNotFound(t *testing.T) {
	router, _ := testRouter(t)

	req := newTestRequest(http.MethodDelete, "/api/ocr/document/"+uuid.NewString(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestDeleteOCRDocumentReturnsNotFoundForWrongOwner(t *testing.T) {
	router, db := testRouter(t)
	owner := createTestUser(t, db, "document-delete-real-owner@example.com")
	other := createTestUser(t, db, "document-delete-other-owner@example.com")
	doc := createListOCRDocument(t, db, ocr.OCRDocument{UserID: &owner.ID, OriginalFilename: "owned.pdf"})

	req := newTestRequest(
		http.MethodDelete,
		"/api/ocr/documents/"+doc.ID.String()+"?user_id="+url.QueryEscape(other.ID),
		nil,
	)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if got.Error != "OCR document not found" {
		t.Fatalf("error = %q, want OCR document not found", got.Error)
	}
	assertOCRDocumentExists(t, db, doc.ID, true)
}

func TestDeleteOCRDocumentsDeletesOnlyScopedDocuments(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "document-bulk-delete@example.com")
	other := createTestUser(t, db, "document-bulk-delete-other@example.com")
	owned1 := createListOCRDocument(t, db, ocr.OCRDocument{UserID: &user.ID, OriginalFilename: "owned-1.pdf"})
	owned2 := createListOCRDocument(t, db, ocr.OCRDocument{UserID: &user.ID, OriginalFilename: "owned-2.pdf"})
	otherDoc := createListOCRDocument(t, db, ocr.OCRDocument{UserID: &other.ID, OriginalFilename: "other.pdf"})
	systemDoc := createListOCRDocument(t, db, ocr.OCRDocument{OriginalFilename: "system.pdf"})
	missingID := uuid.New()
	body := `{"ids":["` + owned1.ID.String() + `","` + otherDoc.ID.String() + `","` + missingID.String() + `","` + owned2.ID.String() + `"]}`

	req := newTestRequest(
		http.MethodDelete,
		"/api/ocr/documents?user_id="+url.QueryEscape(user.ID),
		strings.NewReader(body),
	)
	req.Header.Set("content-type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeDeleteOCRDocumentsResponse(t, w)
	if got.DeletedCount != 2 || len(got.DeletedIDs) != 2 || got.DeletedIDs[0] != owned1.ID || got.DeletedIDs[1] != owned2.ID {
		t.Fatalf("delete response = %#v, want owned docs in request order", got)
	}
	assertOCRDocumentExists(t, db, owned1.ID, false)
	assertOCRDocumentExists(t, db, owned2.ID, false)
	assertOCRDocumentExists(t, db, otherDoc.ID, true)
	assertOCRDocumentExists(t, db, systemDoc.ID, true)
}

func TestDeleteOCRDocumentsRemovesCollectionAssociations(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "document-bulk-delete-collection-owner@example.com")
	collection := createCollection(t, db, user.ID, "Linked Documents", time.Now())
	doc := createListOCRDocument(t, db, ocr.OCRDocument{UserID: &user.ID, OriginalFilename: "linked.pdf"})
	if err := db.Create(&ocr.CollectionDocument{CollectionID: collection.ID, DocumentID: doc.ID}).Error; err != nil {
		t.Fatalf("create collection document: %v", err)
	}

	req := newTestRequest(
		http.MethodDelete,
		"/api/ocr/documents?user_id="+url.QueryEscape(user.ID),
		strings.NewReader(`{"ids":["`+doc.ID.String()+`"]}`),
	)
	req.Header.Set("content-type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	assertCount(t, db, &ocr.CollectionDocument{}, "collection_id = ? AND document_id = ?", collection.ID, doc.ID, 0)
	assertCount(t, db, &ocr.OCRDocument{}, "id = ?", doc.ID, 0)
	assertCount(t, db.Unscoped(), &ocr.OCRDocument{}, "id = ?", doc.ID, 1)
	assertCount(t, db, &ocr.Collection{}, "id = ?", collection.ID, 1)
}

func TestDeleteOCRDocumentsRejectsInvalidBulkBodies(t *testing.T) {
	router, _ := testRouter(t)

	cases := []struct {
		name      string
		body      string
		wantError string
	}{
		{name: "malformed json", body: `{`, wantError: "invalid OCR document delete request"},
		{name: "missing ids", body: `{}`, wantError: "ids is required"},
		{name: "empty ids", body: `{"ids":[]}`, wantError: "ids is required"},
		{name: "invalid id", body: `{"ids":["not-a-uuid"]}`, wantError: "invalid OCR document id"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := newTestRequest(http.MethodDelete, "/api/ocr/documents", strings.NewReader(tc.body))
			req.Header.Set("content-type", "application/json")
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

func TestMoveOCRDocumentsToCollectionsRequiresUserID(t *testing.T) {
	router, _ := testRouter(t)

	req := newTestRequest(
		http.MethodPut,
		"/api/ocr/documents/collections",
		strings.NewReader(`{"ids":["`+uuid.NewString()+`"],"collection_ids":[]}`),
	)
	req.Header.Set("content-type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if got.Error != "user_id is required" {
		t.Fatalf("error = %q, want user_id is required", got.Error)
	}
}

func TestMoveOCRDocumentsToCollectionsReplacesAssociations(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "document-move-owner@example.com")
	other := createTestUser(t, db, "document-move-other@example.com")
	oldCollection := createCollection(t, db, user.ID, "Old", time.Now())
	targetB := createCollection(t, db, user.ID, "Target B", time.Now())
	targetA := createCollection(t, db, user.ID, "Target A", time.Now())
	otherCollection := createCollection(t, db, other.ID, "Other", time.Now())
	owned1 := createListOCRDocument(t, db, ocr.OCRDocument{UserID: &user.ID, OriginalFilename: "owned-1.pdf"})
	owned2 := createListOCRDocument(t, db, ocr.OCRDocument{UserID: &user.ID, OriginalFilename: "owned-2.pdf"})
	otherDoc := createListOCRDocument(t, db, ocr.OCRDocument{UserID: &other.ID, OriginalFilename: "other.pdf"})
	for _, doc := range []ocr.OCRDocument{owned1, owned2} {
		if err := db.Create(&ocr.CollectionDocument{CollectionID: oldCollection.ID, DocumentID: doc.ID}).Error; err != nil {
			t.Fatalf("create old collection document: %v", err)
		}
	}
	if err := db.Create(&ocr.CollectionDocument{CollectionID: otherCollection.ID, DocumentID: otherDoc.ID}).Error; err != nil {
		t.Fatalf("create other collection document: %v", err)
	}
	missingID := uuid.New()
	body := `{"ids":["` + owned1.ID.String() + `","` + otherDoc.ID.String() + `","` + missingID.String() + `","` + owned2.ID.String() + `"],"collection_ids":["` + targetB.ID.String() + `","` + targetA.ID.String() + `","` + targetA.ID.String() + `"]}`

	req := newTestRequest(
		http.MethodPut,
		"/api/ocr/documents/collections?user_id="+url.QueryEscape(user.ID),
		strings.NewReader(body),
	)
	req.Header.Set("content-type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeMoveOCRDocumentsToCollectionsResponse(t, w)
	if got.MovedCount != 2 || len(got.MovedIDs) != 2 || got.MovedIDs[0] != owned1.ID || got.MovedIDs[1] != owned2.ID {
		t.Fatalf("move response = %#v, want owned docs in request order", got)
	}
	if len(got.CollectionIDs) != 2 || got.CollectionIDs[0] != targetB.ID || got.CollectionIDs[1] != targetA.ID {
		t.Fatalf("collection_ids = %#v, want deduped ids in request order", got.CollectionIDs)
	}
	assertCount(t, db, &ocr.CollectionDocument{}, "document_id = ? AND collection_id = ?", owned1.ID, oldCollection.ID, 0)
	assertCount(t, db, &ocr.CollectionDocument{}, "document_id = ? AND collection_id = ?", owned2.ID, oldCollection.ID, 0)
	for _, doc := range []ocr.OCRDocument{owned1, owned2} {
		assertCount(t, db, &ocr.CollectionDocument{}, "document_id = ? AND collection_id = ?", doc.ID, targetA.ID, 1)
		assertCount(t, db, &ocr.CollectionDocument{}, "document_id = ? AND collection_id = ?", doc.ID, targetB.ID, 1)
	}
	assertCount(t, db, &ocr.CollectionDocument{}, "document_id = ? AND collection_id = ?", otherDoc.ID, otherCollection.ID, 1)
}

func TestMoveOCRDocumentsToCollectionsAllowsEmptyTargets(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "document-move-empty@example.com")
	collection := createCollection(t, db, user.ID, "Old", time.Now())
	doc := createListOCRDocument(t, db, ocr.OCRDocument{UserID: &user.ID, OriginalFilename: "linked.pdf"})
	if err := db.Create(&ocr.CollectionDocument{CollectionID: collection.ID, DocumentID: doc.ID}).Error; err != nil {
		t.Fatalf("create collection document: %v", err)
	}

	req := newTestRequest(
		http.MethodPut,
		"/api/ocr/documents/collections?user_id="+url.QueryEscape(user.ID),
		strings.NewReader(`{"ids":["`+doc.ID.String()+`"]}`),
	)
	req.Header.Set("content-type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeMoveOCRDocumentsToCollectionsResponse(t, w)
	if got.MovedCount != 1 || len(got.MovedIDs) != 1 || got.MovedIDs[0] != doc.ID || len(got.CollectionIDs) != 0 {
		t.Fatalf("move response = %#v, want document moved to no collections", got)
	}
	assertCount(t, db, &ocr.CollectionDocument{}, "document_id = ?", doc.ID, 0)
}

func TestMoveOCRDocumentsToCollectionsRejectsWrongOwnerCollection(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "document-move-wrong-owner@example.com")
	other := createTestUser(t, db, "document-move-wrong-owner-other@example.com")
	oldCollection := createCollection(t, db, user.ID, "Old", time.Now())
	otherCollection := createCollection(t, db, other.ID, "Other", time.Now())
	doc := createListOCRDocument(t, db, ocr.OCRDocument{UserID: &user.ID, OriginalFilename: "linked.pdf"})
	if err := db.Create(&ocr.CollectionDocument{CollectionID: oldCollection.ID, DocumentID: doc.ID}).Error; err != nil {
		t.Fatalf("create collection document: %v", err)
	}

	req := newTestRequest(
		http.MethodPut,
		"/api/ocr/documents/collections?user_id="+url.QueryEscape(user.ID),
		strings.NewReader(`{"ids":["`+doc.ID.String()+`"],"collection_ids":["`+otherCollection.ID.String()+`"]}`),
	)
	req.Header.Set("content-type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if got.Error != "collection not found" {
		t.Fatalf("error = %q, want collection not found", got.Error)
	}
	assertCount(t, db, &ocr.CollectionDocument{}, "document_id = ? AND collection_id = ?", doc.ID, oldCollection.ID, 1)
}

func TestMoveOCRDocumentsToCollectionsRejectsInvalidBodies(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "document-move-invalid@example.com")

	cases := []struct {
		name      string
		body      string
		wantError string
	}{
		{name: "malformed json", body: `{`, wantError: "invalid OCR document move request"},
		{name: "missing ids", body: `{}`, wantError: "ids is required"},
		{name: "empty ids", body: `{"ids":[]}`, wantError: "ids is required"},
		{name: "invalid document id", body: `{"ids":["not-a-uuid"]}`, wantError: "invalid OCR document id"},
		{name: "invalid collection id", body: `{"ids":["` + uuid.NewString() + `"],"collection_ids":["not-a-uuid"]}`, wantError: "invalid collection id"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := newTestRequest(
				http.MethodPut,
				"/api/ocr/documents/collections?user_id="+url.QueryEscape(user.ID),
				strings.NewReader(tc.body),
			)
			req.Header.Set("content-type", "application/json")
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

func TestListOCRDocumentsScopesByUserID(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "document-list-scope@example.com")
	other := createTestUser(t, db, "document-list-scope-other@example.com")
	owned := createListOCRDocument(t, db, ocr.OCRDocument{UserID: &user.ID, OriginalFilename: "owned.pdf"})
	createListOCRDocument(t, db, ocr.OCRDocument{UserID: &other.ID, OriginalFilename: "other.pdf"})
	createListOCRDocument(t, db, ocr.OCRDocument{OriginalFilename: "system.pdf"})

	req := newTestRequest(http.MethodGet, "/api/ocr/documents?user_id="+url.QueryEscape(user.ID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeOCRDocumentListResponse(t, w)
	assertOCRDocumentListIDs(t, got.Documents, owned.ID)
	if got.Documents[0].UserID == nil || string(*got.Documents[0].UserID) != user.ID {
		t.Fatalf("user_id = %#v, want %s", got.Documents[0].UserID, user.ID)
	}
}

func TestListOCRDocumentsFiltersByCollection(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "document-list-collection@example.com")
	collection := createCollection(t, db, user.ID, "Invoices", time.Now())
	base := time.Date(2026, 5, 27, 12, 0, 0, 0, time.UTC)
	older := createListOCRDocument(t, db, ocr.OCRDocument{UserID: &user.ID, OriginalFilename: "old-invoice.pdf"})
	older = setListOCRDocumentCreatedAt(t, db, older, base)
	newer := createListOCRDocument(t, db, ocr.OCRDocument{UserID: &user.ID, OriginalFilename: "new-invoice.pdf"})
	newer = setListOCRDocumentCreatedAt(t, db, newer, base.Add(time.Hour))
	unlinked := createListOCRDocument(t, db, ocr.OCRDocument{UserID: &user.ID, OriginalFilename: "other.pdf"})
	if err := db.Create(&ocr.CollectionDocument{CollectionID: collection.ID, DocumentID: older.ID}).Error; err != nil {
		t.Fatalf("create older collection document: %v", err)
	}
	if err := db.Create(&ocr.CollectionDocument{CollectionID: collection.ID, DocumentID: newer.ID}).Error; err != nil {
		t.Fatalf("create newer collection document: %v", err)
	}
	_ = unlinked

	req := newTestRequest(
		http.MethodGet,
		"/api/ocr/documents?user_id="+url.QueryEscape(user.ID)+"&collection_id="+collection.ID.String(),
		nil,
	)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeOCRDocumentListResponse(t, w)
	assertOCRDocumentListIDs(t, got.Documents, newer.ID, older.ID)
}

func TestListOCRDocumentsFiltersBySchema(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "document-list-schema@example.com")

	schema := ocr.ExtractionSchema{
		UserID:     &user.ID,
		Name:       "Invoices",
		SchemaJSON: []byte(`{"type":"object"}`),
	}
	if err := db.Create(&schema).Error; err != nil {
		t.Fatalf("create schema: %v", err)
	}

	base := time.Date(2026, 5, 27, 12, 0, 0, 0, time.UTC)
	matched := createListOCRDocument(t, db, ocr.OCRDocument{
		UserID:           &user.ID,
		OriginalFilename: "matched-invoice.pdf",
		SchemaID:         &schema.ID,
	})
	matched = setListOCRDocumentCreatedAt(t, db, matched, base)

	unmatched := createListOCRDocument(t, db, ocr.OCRDocument{
		UserID:           &user.ID,
		OriginalFilename: "other.pdf",
	})
	unmatched = setListOCRDocumentCreatedAt(t, db, unmatched, base.Add(time.Hour))

	req := newTestRequest(
		http.MethodGet,
		"/api/ocr/documents?user_id="+url.QueryEscape(user.ID)+"&schema_id="+schema.ID.String(),
		nil,
	)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeOCRDocumentListResponse(t, w)
	assertOCRDocumentListIDs(t, got.Documents, matched.ID)
}

func TestListOCRDocumentsIncludesCollectionSummaries(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "document-list-collections-summary@example.com")
	receipts := createCollection(t, db, user.ID, "Receipts", time.Now())
	invoices := createCollection(t, db, user.ID, "Invoices", time.Now())
	doc := createListOCRDocument(t, db, ocr.OCRDocument{UserID: &user.ID, OriginalFilename: "linked.pdf"})
	for _, collection := range []ocr.Collection{receipts, invoices} {
		if err := db.Create(&ocr.CollectionDocument{CollectionID: collection.ID, DocumentID: doc.ID}).Error; err != nil {
			t.Fatalf("create collection document: %v", err)
		}
	}

	req := newTestRequest(http.MethodGet, "/api/ocr/documents?user_id="+url.QueryEscape(user.ID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeOCRDocumentListResponse(t, w)
	assertOCRDocumentListIDs(t, got.Documents, doc.ID)
	collections := got.Documents[0].Collections
	if len(collections) != 2 {
		t.Fatalf("collections = %#v, want 2 collection summaries", collections)
	}
	if collections[0].ID != invoices.ID || collections[0].Name != "Invoices" {
		t.Fatalf("collections[0] = %#v, want Invoices sorted by name", collections[0])
	}
	if collections[1].ID != receipts.ID || collections[1].Name != "Receipts" {
		t.Fatalf("collections[1] = %#v, want Receipts sorted by name", collections[1])
	}
}

func TestListOCRDocumentsRejectsWrongOwnerCollection(t *testing.T) {
	router, db := testRouter(t)
	owner := createTestUser(t, db, "document-list-collection-owner@example.com")
	other := createTestUser(t, db, "document-list-collection-other@example.com")
	collection := createCollection(t, db, owner.ID, "Private", time.Now())

	req := newTestRequest(
		http.MethodGet,
		"/api/ocr/documents?user_id="+url.QueryEscape(other.ID)+"&collection_id="+collection.ID.String(),
		nil,
	)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if got.Error != "collection not found" {
		t.Fatalf("error = %q, want collection not found", got.Error)
	}
}

func TestListOCRDocumentsRejectsInvalidCollectionFilter(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "document-list-collection-invalid@example.com")

	req := newTestRequest(
		http.MethodGet,
		"/api/ocr/documents?user_id="+url.QueryEscape(user.ID)+"&collection_id=not-a-uuid",
		nil,
	)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if got.Error != "invalid collection id" {
		t.Fatalf("error = %q, want invalid collection id", got.Error)
	}
}

func TestListOCRDocumentsCollectionFilterComposesWithPagination(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "document-list-collection-page@example.com")
	collection := createCollection(t, db, user.ID, "Paged", time.Now())
	base := time.Date(2026, 5, 27, 12, 0, 0, 0, time.UTC)
	first := setListOCRDocumentCreatedAt(t, db, createListOCRDocument(t, db, ocr.OCRDocument{UserID: &user.ID, OriginalFilename: "invoice-a.pdf"}), base)
	second := setListOCRDocumentCreatedAt(t, db, createListOCRDocument(t, db, ocr.OCRDocument{UserID: &user.ID, OriginalFilename: "invoice-b.pdf"}), base.Add(time.Hour))
	third := setListOCRDocumentCreatedAt(t, db, createListOCRDocument(t, db, ocr.OCRDocument{UserID: &user.ID, OriginalFilename: "invoice-c.pdf"}), base.Add(2*time.Hour))
	for _, doc := range []ocr.OCRDocument{first, second, third} {
		if err := db.Create(&ocr.CollectionDocument{CollectionID: collection.ID, DocumentID: doc.ID}).Error; err != nil {
			t.Fatalf("create collection document: %v", err)
		}
	}

	path := "/api/ocr/documents?user_id=" + url.QueryEscape(user.ID) + "&collection_id=" + collection.ID.String() + "&filename=invoice&size=2"
	req := newTestRequest(http.MethodGet, path, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("first status = %d body=%s", w.Code, w.Body.String())
	}
	page := decodeOCRDocumentListResponse(t, w)
	assertOCRDocumentListIDs(t, page.Documents, third.ID, second.ID)
	if page.NextCursor == nil {
		t.Fatal("next_cursor = nil, want cursor")
	}

	req = newTestRequest(http.MethodGet, path+"&cursor="+url.QueryEscape(*page.NextCursor), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("next status = %d body=%s", w.Code, w.Body.String())
	}
	page = decodeOCRDocumentListResponse(t, w)
	assertOCRDocumentListIDs(t, page.Documents, first.ID)
}

func TestListOCRDocumentsFiltersFilenameCaseInsensitiveContains(t *testing.T) {
	router, db := testRouter(t)
	report := createListOCRDocument(t, db, ocr.OCRDocument{OriginalFilename: "Quarterly-Report.PDF"})
	createListOCRDocument(t, db, ocr.OCRDocument{OriginalFilename: "invoice.pdf"})

	req := newTestRequest(http.MethodGet, "/api/ocr/documents?filename=REPORT", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeOCRDocumentListResponse(t, w)
	assertOCRDocumentListIDs(t, got.Documents, report.ID)
}

func TestListOCRDocumentsEscapesFilenameWildcards(t *testing.T) {
	router, db := testRouter(t)
	percent := createListOCRDocument(t, db, ocr.OCRDocument{OriginalFilename: "literal-percent%.pdf"})
	underscore := createListOCRDocument(t, db, ocr.OCRDocument{OriginalFilename: "literal_underscore.pdf"})
	createListOCRDocument(t, db, ocr.OCRDocument{OriginalFilename: "plain.pdf"})

	req := newTestRequest(http.MethodGet, "/api/ocr/documents?filename="+url.QueryEscape("%"), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("percent status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeOCRDocumentListResponse(t, w)
	assertOCRDocumentListIDs(t, got.Documents, percent.ID)

	req = newTestRequest(http.MethodGet, "/api/ocr/documents?filename="+url.QueryEscape("_"), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("underscore status = %d body=%s", w.Code, w.Body.String())
	}
	got = decodeOCRDocumentListResponse(t, w)
	assertOCRDocumentListIDs(t, got.Documents, underscore.ID)
}

func TestListOCRDocumentsFiltersCreatedAtRange(t *testing.T) {
	router, db := testRouter(t)
	base := time.Date(2026, 5, 27, 12, 0, 0, 0, time.UTC)
	before := createListOCRDocument(t, db, ocr.OCRDocument{OriginalFilename: "before.pdf"})
	_ = setListOCRDocumentCreatedAt(t, db, before, base.Add(-time.Second))
	from := createListOCRDocument(t, db, ocr.OCRDocument{OriginalFilename: "from.pdf"})
	from = setListOCRDocumentCreatedAt(t, db, from, base)
	to := createListOCRDocument(t, db, ocr.OCRDocument{OriginalFilename: "to.pdf"})
	to = setListOCRDocumentCreatedAt(t, db, to, base.Add(time.Hour))
	after := createListOCRDocument(t, db, ocr.OCRDocument{OriginalFilename: "after.pdf"})
	_ = setListOCRDocumentCreatedAt(t, db, after, base.Add(time.Hour+time.Second))

	values := url.Values{}
	values.Set("created_from", base.Format(time.RFC3339Nano))
	values.Set("created_to", base.Add(time.Hour).Format(time.RFC3339Nano))
	req := newTestRequest(http.MethodGet, "/api/ocr/documents?"+values.Encode(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeOCRDocumentListResponse(t, w)
	assertOCRDocumentListIDs(t, got.Documents, to.ID, from.ID)
}

func TestListOCRDocumentsRejectsInvalidCreatedAtRange(t *testing.T) {
	router, _ := testRouter(t)

	cases := []struct {
		name      string
		query     string
		wantError string
	}{
		{name: "invalid from", query: "created_from=not-a-time", wantError: "invalid created_from"},
		{name: "invalid to", query: "created_to=not-a-time", wantError: "invalid created_to"},
		{name: "inverted", query: "created_from=2026-05-28T00%3A00%3A00Z&created_to=2026-05-27T00%3A00%3A00Z", wantError: "created_from must be before or equal to created_to"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := newTestRequest(http.MethodGet, "/api/ocr/documents?"+tc.query, nil)
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

func TestListOCRDocumentsCursorPaginationByCreatedAt(t *testing.T) {
	router, db := testRouter(t)
	base := time.Date(2026, 5, 27, 12, 0, 0, 0, time.UTC)
	oldest := createListOCRDocument(t, db, ocr.OCRDocument{OriginalFilename: "oldest.pdf"})
	oldest = setListOCRDocumentCreatedAt(t, db, oldest, base)
	middle := createListOCRDocument(t, db, ocr.OCRDocument{OriginalFilename: "middle.pdf"})
	middle = setListOCRDocumentCreatedAt(t, db, middle, base.Add(time.Hour))
	newest := createListOCRDocument(t, db, ocr.OCRDocument{OriginalFilename: "newest.pdf"})
	newest = setListOCRDocumentCreatedAt(t, db, newest, base.Add(2*time.Hour))

	req := newTestRequest(http.MethodGet, "/api/ocr/documents?size=2", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("desc first status = %d body=%s", w.Code, w.Body.String())
	}
	first := decodeOCRDocumentListResponse(t, w)
	assertOCRDocumentListIDs(t, first.Documents, newest.ID, middle.ID)
	if first.NextCursor == nil {
		t.Fatal("desc next_cursor = nil, want cursor")
	}

	req = newTestRequest(http.MethodGet, "/api/ocr/documents?cursor="+*first.NextCursor, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("desc next status = %d body=%s", w.Code, w.Body.String())
	}
	next := decodeOCRDocumentListResponse(t, w)
	assertOCRDocumentListIDs(t, next.Documents, oldest.ID)

	req = newTestRequest(http.MethodGet, "/api/ocr/documents?sort=asc&size=2", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("asc first status = %d body=%s", w.Code, w.Body.String())
	}
	first = decodeOCRDocumentListResponse(t, w)
	assertOCRDocumentListIDs(t, first.Documents, oldest.ID, middle.ID)
	if first.NextCursor == nil {
		t.Fatal("asc next_cursor = nil, want cursor")
	}

	req = newTestRequest(http.MethodGet, "/api/ocr/documents?sort=asc&cursor="+*first.NextCursor, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("asc next status = %d body=%s", w.Code, w.Body.String())
	}
	next = decodeOCRDocumentListResponse(t, w)
	assertOCRDocumentListIDs(t, next.Documents, newest.ID)
}

func TestListOCRDocumentsRejectsInvalidPaginationParameters(t *testing.T) {
	router, _ := testRouter(t)

	cases := []struct {
		name      string
		query     string
		wantError string
	}{
		{name: "sort", query: "sort=sideways", wantError: "sort must be asc or desc"},
		{name: "size", query: "size=0", wantError: "size must be between 1 and 100"},
		{name: "cursor", query: "cursor=not-a-cursor", wantError: "invalid cursor"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := newTestRequest(http.MethodGet, "/api/ocr/documents?"+tc.query, nil)
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

func TestDeleteOCRDocumentSoftDeletesAndHidesDocument(t *testing.T) {
	router, db := testRouter(t)
	doc := createListOCRDocument(t, db, ocr.OCRDocument{OriginalFilename: "delete-me.pdf"})

	req := newTestRequest(http.MethodDelete, "/api/ocr/documents/"+doc.ID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if w.Body.Len() != 0 {
		t.Fatalf("body = %q, want empty", w.Body.String())
	}

	var visibleCount int64
	if err := db.Model(&ocr.OCRDocument{}).Where("id = ?", doc.ID).Count(&visibleCount).Error; err != nil {
		t.Fatalf("count visible OCR documents: %v", err)
	}
	if visibleCount != 0 {
		t.Fatalf("visible document count = %d, want 0", visibleCount)
	}

	var deleted ocr.OCRDocument
	if err := db.Unscoped().First(&deleted, "id = ?", doc.ID).Error; err != nil {
		t.Fatalf("load soft-deleted OCR document: %v", err)
	}
	if !deleted.DeletedAt.Valid {
		t.Fatal("deleted_at is not set")
	}

	req = newTestRequest(http.MethodGet, "/api/ocr/document/"+doc.ID.String(), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("get status = %d body=%s", w.Code, w.Body.String())
	}

	req = newTestRequest(http.MethodGet, "/api/ocr/documents", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("list status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeOCRDocumentListResponse(t, w)
	assertOCRDocumentListIDs(t, got.Documents)
}

func TestDeleteOCRDocumentRejectsInvalidID(t *testing.T) {
	router, _ := testRouter(t)

	req := newTestRequest(http.MethodDelete, "/api/ocr/documents/not-a-uuid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if got.Error != "invalid OCR document id" {
		t.Fatalf("error = %q, want invalid OCR document id", got.Error)
	}
}

func TestDeleteOCRDocumentReturnsNotFound(t *testing.T) {
	router, _ := testRouter(t)

	req := newTestRequest(http.MethodDelete, "/api/ocr/documents/"+uuid.NewString(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if got.Error != "OCR document not found" {
		t.Fatalf("error = %q, want OCR document not found", got.Error)
	}
}

func TestDeleteOCRDocumentReturnsNotFoundWhenAlreadyDeleted(t *testing.T) {
	router, db := testRouter(t)
	doc := createListOCRDocument(t, db, ocr.OCRDocument{})

	req := newTestRequest(http.MethodDelete, "/api/ocr/documents/"+doc.ID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("first status = %d body=%s", w.Code, w.Body.String())
	}

	req = newTestRequest(http.MethodDelete, "/api/ocr/documents/"+doc.ID.String(), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("second status = %d body=%s", w.Code, w.Body.String())
	}
	var got ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if got.Error != "OCR document not found" {
		t.Fatalf("error = %q, want OCR document not found", got.Error)
	}
}

func TestDeleteOCRDocumentRemovesDocumentFromOCRCache(t *testing.T) {
	_, db := testRouter(t)
	file := validPNGBytes()
	documentHash, err := computeDocumentHash(file, nil, false)
	if err != nil {
		t.Fatalf("compute document hash: %v", err)
	}
	cached := createListOCRDocument(t, db, ocr.OCRDocument{
		OriginalFilename: "cached.png",
		MimeType:         "image/png",
		FileSize:         int64(len(file)),
		DocumentHash:     documentHash,
		Markdown:         "# Cached",
		RawResponseJSON:  datatypes.JSON([]byte(`{"pages":[{"index":0,"markdown":"# Cached"}]}`)),
	})
	if err := db.Delete(&cached).Error; err != nil {
		t.Fatalf("soft delete cached OCR document: %v", err)
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
		t.Fatalf("decode response: %v", err)
	}
	if got.ID == cached.ID || got.Cached || got.Markdown != "# Fresh" {
		t.Fatalf("unexpected OCR response: %#v", got)
	}

	var visibleCount int64
	if err := db.Model(&ocr.OCRDocument{}).Count(&visibleCount).Error; err != nil {
		t.Fatalf("count visible OCR documents: %v", err)
	}
	if visibleCount != 1 {
		t.Fatalf("visible document count = %d, want 1", visibleCount)
	}
	var totalCount int64
	if err := db.Unscoped().Model(&ocr.OCRDocument{}).Count(&totalCount).Error; err != nil {
		t.Fatalf("count all OCR documents: %v", err)
	}
	if totalCount != 2 {
		t.Fatalf("all document count = %d, want 2", totalCount)
	}
}

func TestDeleteOCRDocumentKeepsCompletedJobDocumentID(t *testing.T) {
	router, db, _, _ := testRouterWithOCRJobs(t)
	doc := createListOCRDocument(t, db, ocr.OCRDocument{})
	job := createStoredOCRJob(t, db, ocr.OCRJob{
		DocumentID: &doc.ID,
		Status:     ocr.OCRJobStatusCompleted,
	})

	req := newTestRequest(http.MethodDelete, "/api/ocr/documents/"+doc.ID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}

	var got ocr.OCRJob
	if err := db.First(&got, "id = ?", job.ID).Error; err != nil {
		t.Fatalf("load OCR job: %v", err)
	}
	if got.DocumentID == nil || *got.DocumentID != doc.ID {
		t.Fatalf("job document_id = %#v, want %s", got.DocumentID, doc.ID)
	}
}
