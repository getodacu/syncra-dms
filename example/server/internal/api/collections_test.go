package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"ai.ro/syncra/internal/ocr"
)

func TestCreateCollectionRequiresUserID(t *testing.T) {
	router, _ := testRouter(t)
	body := []byte(`{"name":"Invoices"}`)

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodPost, "/api/collection", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Error != "user_id is required" {
		t.Fatalf("error = %q, want user_id is required", got.Error)
	}
}

func TestCreateCollectionRejectsInvalidSchemaIDs(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "collection-owner@example.com")
	body := []byte(`{"name":"Invoices","user_id":"` + user.ID + `","schema_ids":["not-a-uuid"]}`)

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodPost, "/api/collection", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Error != "invalid schema_ids" {
		t.Fatalf("error = %q, want invalid schema_ids", got.Error)
	}
}

func TestCreateCollectionPersistsSchemas(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "collection-create@example.com")
	schemaA := createCollectionSchema(t, db, user.ID, "invoice-a", uuid.MustParse("00000000-0000-0000-0000-00000000000a"))
	schemaB := createCollectionSchema(t, db, user.ID, "invoice-b", uuid.MustParse("00000000-0000-0000-0000-00000000000b"))
	body := []byte(`{"name":" Invoices ","user_id":"` + user.ID + `","schema_ids":["` + schemaA.ID.String() + `","` + schemaB.ID.String() + `"]}`)

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodPost, "/api/collection", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got CollectionResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Name != "Invoices" || string(got.UserID) != user.ID || got.SchemaCount != 2 || got.DocumentCount != 0 {
		t.Fatalf("unexpected response: %#v", got)
	}
	if !reflect.DeepEqual(got.SchemaIDs, []uuid.UUID{schemaA.ID, schemaB.ID}) {
		t.Fatalf("schema_ids = %#v, want schema IDs in request order", got.SchemaIDs)
	}
	assertCount(t, db, &ocr.Collection{}, "id = ?", got.ID, 1)
	assertCount(t, db, &ocr.CollectionSchema{}, "collection_id = ?", got.ID, 2)
}

func TestCreateCollectionAllowsEmptySchemaIDs(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "collection-empty@example.com")
	body := []byte(`{"name":"Empty","user_id":"` + user.ID + `","schema_ids":[]}`)

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodPost, "/api/collection", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got CollectionResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Name != "Empty" || string(got.UserID) != user.ID || got.SchemaCount != 0 || len(got.SchemaIDs) != 0 {
		t.Fatalf("unexpected response: %#v", got)
	}
	assertCount(t, db, &ocr.CollectionSchema{}, "collection_id = ?", got.ID, 0)
}

func TestCreateCollectionRejectsOtherUserSchema(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "collection-owner-schemas@example.com")
	other := createTestUser(t, db, "collection-other-schemas@example.com")
	otherSchema := createCollectionSchema(t, db, other.ID, "other-schema", uuid.MustParse("00000000-0000-0000-0000-00000000000c"))
	body := []byte(`{"name":"Invoices","user_id":"` + user.ID + `","schema_ids":["` + otherSchema.ID.String() + `"]}`)

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodPost, "/api/collection", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Error != "invalid schema_ids" {
		t.Fatalf("error = %q, want invalid schema_ids", got.Error)
	}
	assertCount(t, db, &ocr.Collection{}, "user_id = ?", user.ID, 0)
}

func TestGetCollectionScopesByUserID(t *testing.T) {
	router, db := testRouter(t)
	owner := createTestUser(t, db, "collection-get-owner@example.com")
	other := createTestUser(t, db, "collection-get-other@example.com")
	schemaA := createCollectionSchema(t, db, owner.ID, "schema-a", uuid.MustParse("00000000-0000-0000-0000-00000000000d"))
	schemaB := createCollectionSchema(t, db, owner.ID, "schema-b", uuid.MustParse("00000000-0000-0000-0000-00000000000e"))
	collection := createCollection(t, db, owner.ID, "Invoices", time.Now())
	createCollectionSchemas(t, db, collection.ID, schemaB.ID, schemaA.ID)

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodGet, "/api/collections/"+collection.ID.String()+"?user_id="+url.QueryEscape(other.ID), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("wrong-owner status = %d body=%s", w.Code, w.Body.String())
	}

	w = httptest.NewRecorder()
	req = newTestRequest(http.MethodGet, "/api/collections/"+collection.ID.String()+"?user_id="+url.QueryEscape(owner.ID), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("owner status = %d body=%s", w.Code, w.Body.String())
	}
	var got CollectionResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ID != collection.ID || string(got.UserID) != owner.ID || got.Name != "Invoices" || got.SchemaCount != 2 {
		t.Fatalf("unexpected response: %#v", got)
	}
	if !reflect.DeepEqual(got.SchemaIDs, []uuid.UUID{schemaA.ID, schemaB.ID}) {
		t.Fatalf("schema_ids = %#v, want schema IDs sorted ascending", got.SchemaIDs)
	}
}

func TestListCollectionsReturnsUserCollectionsOnly(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "collection-list-owner@example.com")
	other := createTestUser(t, db, "collection-list-other@example.com")
	base := time.Date(2026, 5, 31, 9, 0, 0, 0, time.UTC)
	older := createCollection(t, db, user.ID, "Older", base)
	newer := createCollection(t, db, user.ID, "Newer", base.Add(time.Hour))
	_ = createCollection(t, db, other.ID, "Other", base.Add(2*time.Hour))

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodGet, "/api/collections?user_id="+url.QueryEscape(user.ID), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got CollectionListResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.NextCursor != nil {
		t.Fatalf("next_cursor = %q, want nil", *got.NextCursor)
	}
	if len(got.Collections) != 2 || got.Collections[0].ID != newer.ID || got.Collections[1].ID != older.ID {
		t.Fatalf("collections = %#v, want only user collections newest first", got.Collections)
	}
}

func TestListCollectionsCursorPaginationAndAscendingSort(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "collection-list-cursor@example.com")
	base := time.Date(2026, 5, 31, 10, 0, 0, 0, time.UTC)
	oldest := createCollection(t, db, user.ID, "Oldest", base)
	middle := createCollection(t, db, user.ID, "Middle", base.Add(time.Hour))
	newest := createCollection(t, db, user.ID, "Newest", base.Add(2*time.Hour))

	query := url.Values{}
	query.Set("user_id", user.ID)
	query.Set("size", "2")
	req := newTestRequest(http.MethodGet, "/api/collections?"+query.Encode(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("desc first status = %d body=%s", w.Code, w.Body.String())
	}
	first := decodeCollectionListResponse(t, w)
	assertCollectionListIDs(t, first.Collections, newest.ID, middle.ID)
	if first.NextCursor == nil || *first.NextCursor == "" {
		t.Fatalf("desc next_cursor = %#v, want non-empty", first.NextCursor)
	}

	query.Set("cursor", *first.NextCursor)
	req = newTestRequest(http.MethodGet, "/api/collections?"+query.Encode(), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("desc second status = %d body=%s", w.Code, w.Body.String())
	}
	second := decodeCollectionListResponse(t, w)
	assertCollectionListIDs(t, second.Collections, oldest.ID)
	if second.NextCursor != nil {
		t.Fatalf("desc second next_cursor = %q, want nil", *second.NextCursor)
	}

	query = url.Values{}
	query.Set("user_id", user.ID)
	query.Set("sort", "asc")
	query.Set("size", "2")
	req = newTestRequest(http.MethodGet, "/api/collections?"+query.Encode(), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("asc first status = %d body=%s", w.Code, w.Body.String())
	}
	first = decodeCollectionListResponse(t, w)
	assertCollectionListIDs(t, first.Collections, oldest.ID, middle.ID)
	if first.NextCursor == nil || *first.NextCursor == "" {
		t.Fatalf("asc next_cursor = %#v, want non-empty", first.NextCursor)
	}

	query.Set("cursor", *first.NextCursor)
	req = newTestRequest(http.MethodGet, "/api/collections?"+query.Encode(), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("asc second status = %d body=%s", w.Code, w.Body.String())
	}
	second = decodeCollectionListResponse(t, w)
	assertCollectionListIDs(t, second.Collections, newest.ID)
	if second.NextCursor != nil {
		t.Fatalf("asc second next_cursor = %q, want nil", *second.NextCursor)
	}
}

func TestListCollectionsRejectsInvalidAndMismatchedCursor(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "collection-list-bad-cursor@example.com")
	collection := createCollection(t, db, user.ID, "Cursor", time.Now())
	descCursor, err := encodeCollectionListCursor(collection, "desc")
	if err != nil {
		t.Fatalf("encode cursor: %v", err)
	}

	cases := []struct {
		name      string
		query     string
		wantError string
	}{
		{name: "invalid", query: "cursor=not-base64", wantError: "invalid cursor"},
		{name: "mismatched sort", query: "sort=asc&cursor=" + url.QueryEscape(descCursor), wantError: "cursor sort does not match sort"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := newTestRequest(http.MethodGet, "/api/collections?user_id="+url.QueryEscape(user.ID)+"&"+tc.query, nil)
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

func TestListCollectionsBatchLoadsResponseData(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "collection-list-batch@example.com")
	schemaA := createCollectionSchema(t, db, user.ID, "batch-a", uuid.MustParse("00000000-0000-0000-0000-000000000013"))
	schemaB := createCollectionSchema(t, db, user.ID, "batch-b", uuid.MustParse("00000000-0000-0000-0000-000000000014"))
	base := time.Date(2026, 5, 31, 11, 0, 0, 0, time.UTC)
	oldest := createCollection(t, db, user.ID, "Oldest", base)
	middle := createCollection(t, db, user.ID, "Middle", base.Add(time.Hour))
	newest := createCollection(t, db, user.ID, "Newest", base.Add(2*time.Hour))
	createCollectionSchemas(t, db, oldest.ID, schemaB.ID, schemaA.ID)
	docA := createStoredOCRDocument(t, db, ocr.OCRDocument{UserID: &user.ID, SchemaID: &schemaA.ID})
	docB := createStoredOCRDocument(t, db, ocr.OCRDocument{UserID: &user.ID, SchemaID: &schemaB.ID})
	if err := db.Create(&[]ocr.CollectionDocument{
		{CollectionID: oldest.ID, DocumentID: docA.ID},
		{CollectionID: oldest.ID, DocumentID: docB.ID},
		{CollectionID: middle.ID, DocumentID: docA.ID},
	}).Error; err != nil {
		t.Fatalf("create collection documents: %v", err)
	}

	queryCount := 0
	callbackName := "test:count_collection_list_queries:" + uuid.NewString()
	if err := db.Callback().Query().Before("gorm:query").Register(callbackName, func(_ *gorm.DB) {
		queryCount++
	}); err != nil {
		t.Fatalf("register query counter: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Callback().Query().Remove(callbackName); err != nil {
			t.Fatalf("remove query counter: %v", err)
		}
	})

	req := newTestRequest(http.MethodGet, "/api/collections?user_id="+url.QueryEscape(user.ID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeCollectionListResponse(t, w)
	assertCollectionListIDs(t, got.Collections, newest.ID, middle.ID, oldest.ID)
	if got.Collections[1].DocumentCount != 1 || got.Collections[2].DocumentCount != 2 {
		t.Fatalf("document counts = %#v, want middle=1 oldest=2", got.Collections)
	}
	if !reflect.DeepEqual(got.Collections[2].SchemaIDs, []uuid.UUID{schemaA.ID, schemaB.ID}) {
		t.Fatalf("oldest schema_ids = %#v, want sorted schema IDs", got.Collections[2].SchemaIDs)
	}
	if queryCount > 3 {
		t.Fatalf("list query count = %d, want at most 3 batched queries", queryCount)
	}
}

func TestUpdateCollectionReplacesSchemaSet(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "collection-update-owner@example.com")
	schemaA := createCollectionSchema(t, db, user.ID, "schema-a", uuid.MustParse("00000000-0000-0000-0000-00000000000f"))
	schemaB := createCollectionSchema(t, db, user.ID, "schema-b", uuid.MustParse("00000000-0000-0000-0000-000000000010"))
	schemaC := createCollectionSchema(t, db, user.ID, "schema-c", uuid.MustParse("00000000-0000-0000-0000-000000000011"))
	collection := createCollection(t, db, user.ID, "Old", time.Now())
	createCollectionSchemas(t, db, collection.ID, schemaA.ID, schemaB.ID)
	body := []byte(`{"name":" Updated ","schema_ids":["` + schemaC.ID.String() + `","` + schemaA.ID.String() + `","` + schemaC.ID.String() + `"]}`)

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodPut, "/api/collection/"+collection.ID.String()+"?user_id="+url.QueryEscape(user.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got CollectionResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ID != collection.ID || got.Name != "Updated" || string(got.UserID) != user.ID || got.SchemaCount != 2 {
		t.Fatalf("unexpected response: %#v", got)
	}
	if !reflect.DeepEqual(got.SchemaIDs, []uuid.UUID{schemaC.ID, schemaA.ID}) {
		t.Fatalf("schema_ids = %#v, want deduped schema IDs in request order", got.SchemaIDs)
	}
	assertCount(t, db, &ocr.CollectionSchema{}, "collection_id = ?", collection.ID, 2)
	assertCount(t, db, &ocr.CollectionSchema{}, "collection_id = ? AND schema_id = ?", collection.ID, schemaA.ID, 1)
	assertCount(t, db, &ocr.CollectionSchema{}, "collection_id = ? AND schema_id = ?", collection.ID, schemaB.ID, 0)
	assertCount(t, db, &ocr.CollectionSchema{}, "collection_id = ? AND schema_id = ?", collection.ID, schemaC.ID, 1)
}

func TestUpdateCollectionWrongUserReturnsNotFoundWithoutChangingRows(t *testing.T) {
	router, db := testRouter(t)
	owner := createTestUser(t, db, "collection-update-wrong-owner@example.com")
	other := createTestUser(t, db, "collection-update-wrong-other@example.com")
	schema := createCollectionSchema(t, db, owner.ID, "wrong-update-schema", uuid.MustParse("00000000-0000-0000-0000-000000000015"))
	collection := createCollection(t, db, owner.ID, "Original", time.Now())
	createCollectionSchemas(t, db, collection.ID, schema.ID)
	body := []byte(`{"name":"Changed","schema_ids":[]}`)

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodPut, "/api/collection/"+collection.ID.String()+"?user_id="+url.QueryEscape(other.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var stored ocr.Collection
	if err := db.First(&stored, "id = ?", collection.ID).Error; err != nil {
		t.Fatalf("load collection: %v", err)
	}
	if stored.Name != "Original" || stored.UserID != owner.ID {
		t.Fatalf("stored collection = %#v, want unchanged owner collection", stored)
	}
	assertCount(t, db, &ocr.CollectionSchema{}, "collection_id = ?", collection.ID, 1)
}

func TestUpdateCollectionDoesNotSucceedWhenRowDisappearsBeforeMutation(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "collection-update-race@example.com")
	collection := createCollection(t, db, user.ID, "Original", time.Now())
	body := []byte(`{"name":"Changed","schema_ids":[]}`)

	deletedDuringUpdate := false
	callbackName := "test:delete_collection_before_update:" + uuid.NewString()
	if err := db.Callback().Update().Before("gorm:update").Register(callbackName, func(tx *gorm.DB) {
		if deletedDuringUpdate || tx.Statement.Schema == nil || tx.Statement.Schema.Table != "collections" {
			return
		}
		deletedDuringUpdate = true
		if err := tx.Session(&gorm.Session{NewDB: true}).Exec("DELETE FROM collections WHERE id = ?", collection.ID).Error; err != nil {
			t.Errorf("delete collection during update callback: %v", err)
		}
	}); err != nil {
		t.Fatalf("register update callback: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Callback().Update().Remove(callbackName); err != nil {
			t.Fatalf("remove update callback: %v", err)
		}
	})

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodPut, "/api/collection/"+collection.ID.String()+"?user_id="+url.QueryEscape(user.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if !deletedDuringUpdate {
		t.Fatal("update callback did not run")
	}
	var stored ocr.Collection
	if err := db.First(&stored, "id = ?", collection.ID).Error; err != nil {
		t.Fatalf("load collection: %v", err)
	}
	if stored.Name != "Original" {
		t.Fatalf("stored name = %q, want unchanged original name", stored.Name)
	}
}

func TestDeleteCollectionDeletesAssociationsButKeepsDocuments(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "collection-delete-owner@example.com")
	schema := createCollectionSchema(t, db, user.ID, "delete-schema", uuid.MustParse("00000000-0000-0000-0000-000000000012"))
	collection := createCollection(t, db, user.ID, "Delete", time.Now())
	createCollectionSchemas(t, db, collection.ID, schema.ID)
	doc := createStoredOCRDocument(t, db, ocr.OCRDocument{UserID: &user.ID, SchemaID: &schema.ID})
	if err := db.Create(&ocr.CollectionDocument{CollectionID: collection.ID, DocumentID: doc.ID}).Error; err != nil {
		t.Fatalf("create collection document: %v", err)
	}

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodDelete, "/api/collection/"+collection.ID.String()+"?user_id="+url.QueryEscape(user.ID), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	assertCount(t, db, &ocr.Collection{}, "id = ?", collection.ID, 0)
	assertCount(t, db, &ocr.CollectionSchema{}, "collection_id = ?", collection.ID, 0)
	assertCount(t, db, &ocr.CollectionDocument{}, "collection_id = ?", collection.ID, 0)
	assertCount(t, db, &ocr.OCRDocument{}, "id = ?", doc.ID, 1)
}

func TestDeleteCollectionWrongUserReturnsNotFoundWithoutDeletingRows(t *testing.T) {
	router, db := testRouter(t)
	owner := createTestUser(t, db, "collection-delete-wrong-owner@example.com")
	other := createTestUser(t, db, "collection-delete-wrong-other@example.com")
	schema := createCollectionSchema(t, db, owner.ID, "wrong-delete-schema", uuid.MustParse("00000000-0000-0000-0000-000000000016"))
	collection := createCollection(t, db, owner.ID, "Delete", time.Now())
	createCollectionSchemas(t, db, collection.ID, schema.ID)

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodDelete, "/api/collection/"+collection.ID.String()+"?user_id="+url.QueryEscape(other.ID), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	assertCount(t, db, &ocr.Collection{}, "id = ?", collection.ID, 1)
	assertCount(t, db, &ocr.CollectionSchema{}, "collection_id = ?", collection.ID, 1)
}

func decodeCollectionListResponse(t *testing.T, w *httptest.ResponseRecorder) CollectionListResponse {
	t.Helper()
	var got CollectionListResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode collection list response: %v body=%s", err, w.Body.String())
	}
	return got
}

func assertCollectionListIDs(t *testing.T, collections []CollectionResponse, want ...uuid.UUID) {
	t.Helper()
	if len(collections) != len(want) {
		t.Fatalf("collection count = %d, want %d: %#v", len(collections), len(want), collections)
	}
	for i, id := range want {
		if collections[i].ID != id {
			t.Fatalf("collection[%d].id = %s, want %s; collections=%#v", i, collections[i].ID, id, collections)
		}
	}
}

func createCollectionSchema(t *testing.T, db *gorm.DB, userID string, name string, id uuid.UUID) ocr.ExtractionSchema {
	t.Helper()
	schema := ocr.ExtractionSchema{
		ID:         id,
		UserID:     &userID,
		Name:       name,
		SchemaJSON: datatypes.JSON([]byte(`{"type":"object","title":"` + name + `"}`)),
		Strict:     true,
	}
	if err := db.Create(&schema).Error; err != nil {
		t.Fatalf("create collection schema fixture: %v", err)
	}
	return schema
}

func createCollection(t *testing.T, db *gorm.DB, userID string, name string, createdAt time.Time) ocr.Collection {
	t.Helper()
	collection := ocr.Collection{
		UserID:    userID,
		Name:      name,
		CreatedAt: createdAt.UTC(),
		UpdatedAt: createdAt.UTC(),
	}
	if err := db.Create(&collection).Error; err != nil {
		t.Fatalf("create collection fixture: %v", err)
	}
	return collection
}

func createCollectionSchemas(t *testing.T, db *gorm.DB, collectionID uuid.UUID, schemaIDs ...uuid.UUID) {
	t.Helper()
	for _, schemaID := range schemaIDs {
		if err := db.Create(&ocr.CollectionSchema{CollectionID: collectionID, SchemaID: schemaID}).Error; err != nil {
			t.Fatalf("create collection schema link: %v", err)
		}
	}
}

func assertCount(t *testing.T, db *gorm.DB, model any, query string, args ...any) {
	t.Helper()
	if len(args) == 0 {
		t.Fatal("assertCount requires a want count")
	}
	var want int64
	switch value := args[len(args)-1].(type) {
	case int:
		want = int64(value)
	case int64:
		want = value
	default:
		t.Fatalf("assertCount want has type %T, want int or int64", value)
	}
	queryArgs := args[:len(args)-1]
	var count int64
	if err := db.Model(model).Where(query, queryArgs...).Count(&count).Error; err != nil {
		t.Fatalf("count %T: %v", model, err)
	}
	if count != want {
		t.Fatalf("count %T where %q = %d, want %d", model, query, count, want)
	}
}
