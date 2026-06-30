package api

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"ai.ro/syncra/internal/ocr"
)

func TestCreateDatasetRequiresUserID(t *testing.T) {
	router, _ := testDatasetRouter(t)
	body := []byte(`{"name":"Invoices"}`)

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodPost, "/api/datasets", bytes.NewReader(body))
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

func TestCreateDatasetRejectsInvalidSchemaID(t *testing.T) {
	router, db := testDatasetRouter(t)
	user := createTestUser(t, db, "dataset-invalid-schema@example.com")
	body := []byte(`{
		"name":"Invoices",
		"user_id":"` + user.ID + `",
		"schema_id":"not-a-uuid",
		"selected_fields":[{"path":"/total","key":"total","label":"Total"}]
	}`)

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodPost, "/api/datasets", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Error != "invalid schema_id" {
		t.Fatalf("error = %q, want invalid schema_id", got.Error)
	}
}

func TestCreateDatasetRejectsOtherUserSchema(t *testing.T) {
	router, db := testDatasetRouter(t)
	user := createTestUser(t, db, "dataset-owner-schema@example.com")
	other := createTestUser(t, db, "dataset-other-schema@example.com")
	otherSchema := createDatasetSchema(t, db, other.ID, "Other", `{
		"type":"object",
		"properties":{"total":{"type":"number"}}
	}`)
	body := []byte(`{
		"name":"Invoices",
		"user_id":"` + user.ID + `",
		"schema_id":"` + otherSchema.ID.String() + `",
		"selected_fields":[{"path":"/total","key":"total","label":"Total"}]
	}`)

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodPost, "/api/datasets", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Error != "invalid schema_id" {
		t.Fatalf("error = %q, want invalid schema_id", got.Error)
	}
	assertCount(t, db, &ocr.Dataset{}, "user_id = ?", user.ID, 0)
}

func TestCreateDatasetRejectsSystemSchema(t *testing.T) {
	router, db := testDatasetRouter(t)
	user := createTestUser(t, db, "dataset-system-schema@example.com")
	systemSchema := createSystemDatasetSchema(t, db, "System", `{
		"type":"object",
		"properties":{"total":{"type":"number"}}
	}`)
	body := []byte(`{
		"name":"Invoices",
		"user_id":"` + user.ID + `",
		"schema_id":"` + systemSchema.ID.String() + `",
		"selected_fields":[{"path":"/total","key":"total","label":"Total"}]
	}`)

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodPost, "/api/datasets", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Error != "invalid schema_id" {
		t.Fatalf("error = %q, want invalid schema_id", got.Error)
	}
	assertCount(t, db, &ocr.Dataset{}, "user_id = ?", user.ID, 0)
}

func TestCreateDatasetRejectsInvalidSelectedFields(t *testing.T) {
	router, db := testDatasetRouter(t)
	user := createTestUser(t, db, "dataset-create-fields@example.com")
	schema := createDatasetSchema(t, db, user.ID, "Invoice", `{
		"type":"object",
		"properties":{"supplier":{"type":"string"},"total":{"type":"number"}}
	}`)

	cases := []struct {
		name   string
		fields []ocr.DatasetField
	}{
		{name: "empty", fields: []ocr.DatasetField{}},
		{name: "invalid path", fields: []ocr.DatasetField{{Path: "/missing", Key: "missing", Label: "Missing"}}},
		{name: "duplicate key", fields: []ocr.DatasetField{
			{Path: "/supplier", Key: "value", Label: "Supplier"},
			{Path: "/total", Key: "value", Label: "Total"},
		}},
		{name: "too many fields", fields: manyDatasetFields(maxDatasetFieldCount + 1)},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			body, err := json.Marshal(createDatasetRequest{
				Name:           "Invoices",
				UserID:         user.ID,
				SchemaID:       schema.ID.String(),
				SelectedFields: tc.fields,
			})
			if err != nil {
				t.Fatalf("marshal request: %v", err)
			}

			w := httptest.NewRecorder()
			req := newTestRequest(http.MethodPost, "/api/datasets", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
			}
			var got ErrorResponse
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Fatalf("decode: %v", err)
			}
			if got.Error != "invalid selected_fields" {
				t.Fatalf("error = %q, want invalid selected_fields", got.Error)
			}
			assertCount(t, db, &ocr.Dataset{}, "user_id = ?", user.ID, 0)
		})
	}
}

func TestCreateDatasetPersistsFields(t *testing.T) {
	router, db := testDatasetRouter(t)
	user := createTestUser(t, db, "dataset-create@example.com")
	schema := createDatasetSchema(t, db, user.ID, "Invoice", `{
		"type":"object",
		"properties":{"supplier":{"type":"object","properties":{"name":{"type":"string"}}},"total":{"type":"number"}}
	}`)
	body := []byte(`{
		"name":" Invoices ",
		"user_id":"` + user.ID + `",
		"schema_id":"` + schema.ID.String() + `",
		"selected_fields":[
			{"path":"/supplier/name","key":"supplier_name","label":"Supplier name"},
			{"path":"/total","key":"total","label":"Total"}
		]
	}`)

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodPost, "/api/datasets", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got DatasetResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Name != "Invoices" || got.SchemaID != schema.ID || got.FieldCount != 2 {
		t.Fatalf("unexpected response: %#v", got)
	}
	wantFields := []DatasetFieldResponse{
		{Path: "/supplier/name", Key: "supplier_name", Label: "Supplier name"},
		{Path: "/total", Key: "total", Label: "Total"},
	}
	if !reflect.DeepEqual(got.SelectedFields, wantFields) {
		t.Fatalf("selected_fields = %#v, want %#v", got.SelectedFields, wantFields)
	}
	assertCount(t, db, &ocr.Dataset{}, "id = ?", got.ID, 1)

	var stored ocr.Dataset
	if err := db.First(&stored, "id = ?", got.ID).Error; err != nil {
		t.Fatalf("load dataset: %v", err)
	}
	assertJSONEqual(t, json.RawMessage(stored.SelectedFields), `[
		{"path":"/supplier/name","key":"supplier_name","label":"Supplier name"},
		{"path":"/total","key":"total","label":"Total"}
	]`)
}

func TestListDatasetsReturnsUserDatasetsOnly(t *testing.T) {
	router, db := testDatasetRouter(t)
	user := createTestUser(t, db, "dataset-list-owner@example.com")
	other := createTestUser(t, db, "dataset-list-other@example.com")
	userSchema := createDatasetSchema(t, db, user.ID, "Invoice", `{"type":"object","properties":{"total":{"type":"number"}}}`)
	otherSchema := createDatasetSchema(t, db, other.ID, "Receipt", `{"type":"object","properties":{"total":{"type":"number"}}}`)
	fields := []ocr.DatasetField{{Path: "/total", Key: "total", Label: "Total"}}
	base := time.Date(2026, 6, 1, 9, 0, 0, 0, time.UTC)
	older := createDataset(t, db, user.ID, userSchema.ID, "Older", fields, base)
	newer := createDataset(t, db, user.ID, userSchema.ID, "Newer", fields, base.Add(time.Hour))
	_ = createDataset(t, db, other.ID, otherSchema.ID, "Other", fields, base.Add(2*time.Hour))

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodGet, "/api/datasets?user_id="+url.QueryEscape(user.ID), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeDatasetListResponse(t, w)
	if got.NextCursor != nil {
		t.Fatalf("next_cursor = %q, want nil", *got.NextCursor)
	}
	assertDatasetListIDs(t, got.Datasets, newer.ID, older.ID)
}

func TestListDatasetsCursorPaginationAndAscendingSort(t *testing.T) {
	router, db := testDatasetRouter(t)
	user := createTestUser(t, db, "dataset-list-cursor@example.com")
	schema := createDatasetSchema(t, db, user.ID, "Invoice", `{"type":"object","properties":{"total":{"type":"number"}}}`)
	fields := []ocr.DatasetField{{Path: "/total", Key: "total", Label: "Total"}}
	base := time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC)
	oldest := createDataset(t, db, user.ID, schema.ID, "Oldest", fields, base)
	middle := createDataset(t, db, user.ID, schema.ID, "Middle", fields, base.Add(time.Hour))
	newest := createDataset(t, db, user.ID, schema.ID, "Newest", fields, base.Add(2*time.Hour))

	query := url.Values{}
	query.Set("user_id", user.ID)
	query.Set("size", "2")
	req := newTestRequest(http.MethodGet, "/api/datasets?"+query.Encode(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("desc first status = %d body=%s", w.Code, w.Body.String())
	}
	first := decodeDatasetListResponse(t, w)
	assertDatasetListIDs(t, first.Datasets, newest.ID, middle.ID)
	if first.NextCursor == nil || *first.NextCursor == "" {
		t.Fatalf("desc next_cursor = %#v, want non-empty", first.NextCursor)
	}

	query.Set("cursor", *first.NextCursor)
	req = newTestRequest(http.MethodGet, "/api/datasets?"+query.Encode(), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("desc second status = %d body=%s", w.Code, w.Body.String())
	}
	second := decodeDatasetListResponse(t, w)
	assertDatasetListIDs(t, second.Datasets, oldest.ID)
	if second.NextCursor != nil {
		t.Fatalf("desc second next_cursor = %q, want nil", *second.NextCursor)
	}

	query = url.Values{}
	query.Set("user_id", user.ID)
	query.Set("sort", "asc")
	query.Set("size", "2")
	req = newTestRequest(http.MethodGet, "/api/datasets?"+query.Encode(), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("asc first status = %d body=%s", w.Code, w.Body.String())
	}
	first = decodeDatasetListResponse(t, w)
	assertDatasetListIDs(t, first.Datasets, oldest.ID, middle.ID)
	if first.NextCursor == nil || *first.NextCursor == "" {
		t.Fatalf("asc next_cursor = %#v, want non-empty", first.NextCursor)
	}

	query.Set("cursor", *first.NextCursor)
	req = newTestRequest(http.MethodGet, "/api/datasets?"+query.Encode(), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("asc second status = %d body=%s", w.Code, w.Body.String())
	}
	second = decodeDatasetListResponse(t, w)
	assertDatasetListIDs(t, second.Datasets, newest.ID)
	if second.NextCursor != nil {
		t.Fatalf("asc second next_cursor = %q, want nil", *second.NextCursor)
	}
}

func TestListDatasetsRejectsInvalidAndMismatchedCursor(t *testing.T) {
	router, db := testDatasetRouter(t)
	user := createTestUser(t, db, "dataset-list-bad-cursor@example.com")
	schema := createDatasetSchema(t, db, user.ID, "Invoice", `{"type":"object","properties":{"total":{"type":"number"}}}`)
	fields := []ocr.DatasetField{{Path: "/total", Key: "total", Label: "Total"}}
	dataset := createDataset(t, db, user.ID, schema.ID, "Cursor", fields, time.Now())
	descCursor, err := encodeDatasetListCursor(dataset, "desc")
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
			req := newTestRequest(http.MethodGet, "/api/datasets?user_id="+url.QueryEscape(user.ID)+"&"+tc.query, nil)
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

func TestListDatasetsDoesNotExposeOtherUserSchemaName(t *testing.T) {
	router, db := testDatasetRouter(t)
	owner := createTestUser(t, db, "dataset-list-corrupt-owner@example.com")
	other := createTestUser(t, db, "dataset-list-corrupt-other@example.com")
	otherSchema := createDatasetSchema(t, db, other.ID, "Other Secret Schema", `{"type":"object","properties":{"total":{"type":"number"}}}`)
	fields := []ocr.DatasetField{{Path: "/total", Key: "total", Label: "Total"}}
	_ = createDataset(t, db, owner.ID, otherSchema.ID, "Corrupt", fields, time.Now())

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodGet, "/api/datasets?user_id="+url.QueryEscape(owner.ID), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if strings.Contains(w.Body.String(), otherSchema.Name) {
		t.Fatalf("response leaked other user's schema name: %s", w.Body.String())
	}
}

func TestGetDatasetScopesByUserID(t *testing.T) {
	router, db := testDatasetRouter(t)
	owner := createTestUser(t, db, "dataset-get-owner@example.com")
	other := createTestUser(t, db, "dataset-get-other@example.com")
	schema := createDatasetSchema(t, db, owner.ID, "Invoice", `{"type":"object","properties":{"total":{"type":"number"}}}`)
	fields := []ocr.DatasetField{{Path: "/total", Key: "total", Label: "Total"}}
	dataset := createDataset(t, db, owner.ID, schema.ID, "Invoices", fields, time.Now())

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodGet, "/api/datasets/"+dataset.ID.String()+"?user_id="+url.QueryEscape(other.ID), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("wrong-owner status = %d body=%s", w.Code, w.Body.String())
	}

	w = httptest.NewRecorder()
	req = newTestRequest(http.MethodGet, "/api/datasets/"+dataset.ID.String()+"?user_id="+url.QueryEscape(owner.ID), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("owner status = %d body=%s", w.Code, w.Body.String())
	}
	var got DatasetResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ID != dataset.ID || string(got.UserID) != owner.ID || got.Name != "Invoices" || got.SchemaName != "Invoice" || got.FieldCount != 1 {
		t.Fatalf("unexpected response: %#v", got)
	}
}

func TestGetDatasetRowsProjectsMatchingSchemaDocuments(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "dataset-rows@example.com")
	schema := createDatasetSchema(t, db, user.ID, "Invoice", `{"type":"object","properties":{"total":{"type":"number"},"line_items":{"type":"array"}}}`)
	dataset := createDatasetFixture(t, db, user.ID, schema.ID, "Invoices", `[{"path":"/total","key":"total","label":"Total"},{"path":"/line_items","key":"line_items","label":"Line items"}]`)
	doc := createDatasetDocument(t, db, user.ID, schema.ID, "invoice.pdf", `{"total":42,"line_items":[{"description":"Work"}]}`)
	otherSchema := createDatasetSchema(t, db, user.ID, "Receipt", `{"type":"object","properties":{"total":{"type":"number"}}}`)
	_ = createDatasetDocument(t, db, user.ID, otherSchema.ID, "receipt.pdf", `{"total":99}`)

	req := newTestRequest(http.MethodGet, "/api/datasets/"+dataset.ID.String()+"/rows?user_id="+url.QueryEscape(user.ID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got DatasetRowsResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(got.Rows) != 1 || got.Rows[0].DocumentID != doc.ID {
		t.Fatalf("rows = %#v, want one matching document", got.Rows)
	}
	if got.Rows[0].Values["line_items"] != `[{"description":"Work"}]` {
		t.Fatalf("line_items = %#v", got.Rows[0].Values["line_items"])
	}
}

func TestGetDatasetRowsScopesDeletesAndPaginates(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "dataset-rows-scope@example.com")
	other := createTestUser(t, db, "dataset-rows-other@example.com")
	schema := createDatasetSchema(t, db, user.ID, "Invoice", `{"type":"object","properties":{"total":{"type":"number"}}}`)
	dataset := createDatasetFixture(t, db, user.ID, schema.ID, "Invoices", `[{"path":"/total","key":"total","label":"Total"}]`)
	base := time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC)
	oldest := setListOCRDocumentCreatedAt(t, db, createDatasetDocument(t, db, user.ID, schema.ID, "oldest.pdf", `{"total":1}`), base)
	middle := setListOCRDocumentCreatedAt(t, db, createDatasetDocument(t, db, user.ID, schema.ID, "middle.pdf", `{"total":2}`), base.Add(time.Hour))
	newest := setListOCRDocumentCreatedAt(t, db, createDatasetDocument(t, db, user.ID, schema.ID, "newest.pdf", `{"total":3}`), base.Add(2*time.Hour))
	deleted := setListOCRDocumentCreatedAt(t, db, createDatasetDocument(t, db, user.ID, schema.ID, "deleted.pdf", `{"total":4}`), base.Add(3*time.Hour))
	if err := db.Delete(&deleted).Error; err != nil {
		t.Fatalf("soft-delete document: %v", err)
	}
	_ = setListOCRDocumentCreatedAt(t, db, createDatasetDocument(t, db, other.ID, schema.ID, "other.pdf", `{"total":5}`), base.Add(4*time.Hour))

	query := url.Values{}
	query.Set("user_id", user.ID)
	query.Set("size", "2")
	req := newTestRequest(http.MethodGet, "/api/datasets/"+dataset.ID.String()+"/rows?"+query.Encode(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("first page status = %d body=%s", w.Code, w.Body.String())
	}
	var first DatasetRowsResponse
	if err := json.Unmarshal(w.Body.Bytes(), &first); err != nil {
		t.Fatalf("decode first page: %v", err)
	}
	assertDatasetRowIDs(t, first.Rows, newest.ID, middle.ID)
	if first.NextCursor == nil || *first.NextCursor == "" {
		t.Fatalf("first next_cursor = %#v, want non-empty", first.NextCursor)
	}

	query.Set("cursor", *first.NextCursor)
	req = newTestRequest(http.MethodGet, "/api/datasets/"+dataset.ID.String()+"/rows?"+query.Encode(), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("second page status = %d body=%s", w.Code, w.Body.String())
	}
	var second DatasetRowsResponse
	if err := json.Unmarshal(w.Body.Bytes(), &second); err != nil {
		t.Fatalf("decode second page: %v", err)
	}
	assertDatasetRowIDs(t, second.Rows, oldest.ID)
	if second.NextCursor != nil {
		t.Fatalf("second next_cursor = %q, want nil", *second.NextCursor)
	}
}

func TestGetDatasetRowsFiltersByCreatedAtRange(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "dataset-rows-date-range@example.com")
	schema := createDatasetSchema(t, db, user.ID, "Invoice", `{"type":"object","properties":{"total":{"type":"number"}}}`)
	dataset := createDatasetFixture(t, db, user.ID, schema.ID, "Invoices", `[{"path":"/total","key":"total","label":"Total"}]`)
	base := time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC)
	_ = setListOCRDocumentCreatedAt(t, db, createDatasetDocument(t, db, user.ID, schema.ID, "before.pdf", `{"total":1}`), base.Add(-time.Hour))
	start := setListOCRDocumentCreatedAt(t, db, createDatasetDocument(t, db, user.ID, schema.ID, "start.pdf", `{"total":2}`), base)
	end := setListOCRDocumentCreatedAt(t, db, createDatasetDocument(t, db, user.ID, schema.ID, "end.pdf", `{"total":3}`), base.Add(time.Hour))
	_ = setListOCRDocumentCreatedAt(t, db, createDatasetDocument(t, db, user.ID, schema.ID, "after.pdf", `{"total":4}`), base.Add(2*time.Hour))

	query := url.Values{}
	query.Set("user_id", user.ID)
	query.Set("created_from", base.Format(time.RFC3339Nano))
	query.Set("created_to", base.Add(time.Hour).Format(time.RFC3339Nano))
	req := newTestRequest(http.MethodGet, "/api/datasets/"+dataset.ID.String()+"/rows?"+query.Encode(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got DatasetRowsResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	assertDatasetRowIDs(t, got.Rows, end.ID, start.ID)
}

func TestGetDatasetRowsRejectsInvalidDateBounds(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "dataset-rows-invalid-date@example.com")
	schema := createDatasetSchema(t, db, user.ID, "Invoice", `{"type":"object","properties":{"total":{"type":"number"}}}`)
	dataset := createDatasetFixture(t, db, user.ID, schema.ID, "Invoices", `[{"path":"/total","key":"total","label":"Total"}]`)

	tests := []struct {
		name      string
		query     string
		wantError string
	}{
		{name: "invalid from", query: "created_from=not-a-date", wantError: "invalid created_from"},
		{name: "invalid to", query: "created_to=not-a-date", wantError: "invalid created_to"},
		{name: "backwards", query: "created_from=2026-06-02T00%3A00%3A00Z&created_to=2026-06-01T00%3A00%3A00Z", wantError: "created_from must be before or equal to created_to"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := newTestRequest(http.MethodGet, "/api/datasets/"+dataset.ID.String()+"/rows?user_id="+url.QueryEscape(user.ID)+"&"+tc.query, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("status = %d body=%s, want 400", w.Code, w.Body.String())
			}
			var got ErrorResponse
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Fatalf("decode error: %v body=%s", err, w.Body.String())
			}
			if got.Error != tc.wantError {
				t.Fatalf("error = %q, want %q", got.Error, tc.wantError)
			}
		})
	}
}

func TestExportDatasetCSV(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "dataset-csv@example.com")
	schema := createDatasetSchema(t, db, user.ID, "Invoice", `{"type":"object","properties":{"total":{"type":"number"}}}`)
	dataset := createDatasetFixture(t, db, user.ID, schema.ID, "Invoices", `[{"path":"/total","key":"total","label":"Total"}]`)
	_ = createDatasetDocument(t, db, user.ID, schema.ID, "invoice,one.pdf", `{"total":42}`)

	req := newTestRequest(http.MethodGet, "/api/datasets/"+dataset.ID.String()+"/export?user_id="+url.QueryEscape(user.ID)+"&format=csv", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if got := w.Header().Get("Content-Type"); !strings.Contains(got, "text/csv") {
		t.Fatalf("content-type = %q, want text/csv", got)
	}
	if !strings.Contains(w.Body.String(), "document_id,filename,created_at,Total") {
		t.Fatalf("csv header missing: %s", w.Body.String())
	}
	if !strings.Contains(w.Body.String(), `"invoice,one.pdf"`) {
		t.Fatalf("csv escaping missing: %s", w.Body.String())
	}
}

func TestExportDatasetCSVFiltersByCreatedAtRange(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "dataset-csv-date-range@example.com")
	schema := createDatasetSchema(t, db, user.ID, "Invoice", `{"type":"object","properties":{"total":{"type":"number"}}}`)
	dataset := createDatasetFixture(t, db, user.ID, schema.ID, "Invoices", `[{"path":"/total","key":"total","label":"Total"}]`)
	base := time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC)
	_ = setListOCRDocumentCreatedAt(t, db, createDatasetDocument(t, db, user.ID, schema.ID, "outside.pdf", `{"total":1}`), base.Add(-time.Hour))
	inside := setListOCRDocumentCreatedAt(t, db, createDatasetDocument(t, db, user.ID, schema.ID, "inside.pdf", `{"total":2}`), base)

	query := url.Values{}
	query.Set("user_id", user.ID)
	query.Set("format", "csv")
	query.Set("created_from", base.Format(time.RFC3339Nano))
	query.Set("created_to", base.Format(time.RFC3339Nano))
	req := newTestRequest(http.MethodGet, "/api/datasets/"+dataset.ID.String()+"/export?"+query.Encode(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	records, err := csv.NewReader(strings.NewReader(w.Body.String())).ReadAll()
	if err != nil {
		t.Fatalf("read csv: %v body=%s", err, w.Body.String())
	}
	if len(records) != 2 {
		t.Fatalf("records = %#v, want header and one filtered row", records)
	}
	if records[1][0] != inside.ID.String() || records[1][1] != "inside.pdf" {
		t.Fatalf("filtered csv row = %#v, want inside document", records[1])
	}
}

func TestExportDatasetCSVSanitizesSpreadsheetInjection(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "dataset-csv-injection@example.com")
	schema := createDatasetSchema(t, db, user.ID, "Invoice", `{"type":"object","properties":{"formula":{"type":"string"}}}`)
	dataset := createDatasetFixture(t, db, user.ID, schema.ID, "Invoices", `[{"path":"/formula","key":"formula","label":"=Formula"}]`)
	_ = createDatasetDocument(t, db, user.ID, schema.ID, "+invoice.pdf", `{"formula":"=SUM(A1:A2)"}`)

	req := newTestRequest(http.MethodGet, "/api/datasets/"+dataset.ID.String()+"/export?user_id="+url.QueryEscape(user.ID)+"&format=csv", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	records, err := csv.NewReader(strings.NewReader(w.Body.String())).ReadAll()
	if err != nil {
		t.Fatalf("read csv: %v body=%s", err, w.Body.String())
	}
	if len(records) != 2 {
		t.Fatalf("records = %#v, want header and one row", records)
	}
	if records[0][3] != "'=Formula" {
		t.Fatalf("sanitized header = %q, want quoted formula header", records[0][3])
	}
	if records[1][1] != "'+invoice.pdf" {
		t.Fatalf("sanitized filename = %q, want quoted filename", records[1][1])
	}
	if records[1][3] != "'=SUM(A1:A2)" {
		t.Fatalf("sanitized value = %q, want quoted formula value", records[1][3])
	}
}

func TestExportDatasetCSVReturnsAPIErrorBeforeBodyForBadAnnotation(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "dataset-csv-bad-annotation@example.com")
	schema := createDatasetSchema(t, db, user.ID, "Invoice", `{"type":"object","properties":{"total":{"type":"number"}}}`)
	dataset := createDatasetFixture(t, db, user.ID, schema.ID, "Invoices", `[{"path":"/total","key":"total","label":"Total"}]`)
	_ = createStoredOCRDocument(t, db, ocr.OCRDocument{
		UserID:           &user.ID,
		SchemaID:         &schema.ID,
		OriginalFilename: "bad.pdf",
	})

	req := newTestRequest(http.MethodGet, "/api/datasets/"+dataset.ID.String()+"/export?user_id="+url.QueryEscape(user.ID)+"&format=csv", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if got := w.Header().Get("Content-Type"); strings.Contains(got, "text/csv") {
		t.Fatalf("content-type = %q, want API error response", got)
	}
	if strings.Contains(w.Body.String(), "document_id,filename,created_at,Total") {
		t.Fatalf("partial csv body leaked: %s", w.Body.String())
	}
	var got ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode error: %v body=%s", err, w.Body.String())
	}
	if got.Error != "failed to export dataset" {
		t.Fatalf("error = %q, want failed to export dataset", got.Error)
	}
}

func TestExportDatasetXLSX(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "dataset-xlsx@example.com")
	schema := createDatasetSchema(t, db, user.ID, "Invoice", `{"type":"object","properties":{"total":{"type":"number"}}}`)
	dataset := createDatasetFixture(t, db, user.ID, schema.ID, "Invoices", `[{"path":"/total","key":"total","label":"Total"}]`)
	doc := createDatasetDocument(t, db, user.ID, schema.ID, "invoice.pdf", `{"total":42}`)

	req := newTestRequest(http.MethodGet, "/api/datasets/"+dataset.ID.String()+"/export?user_id="+url.QueryEscape(user.ID)+"&format=xlsx", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if got := w.Header().Get("Content-Type"); !strings.Contains(got, "spreadsheetml.sheet") {
		t.Fatalf("content-type = %q, want xlsx", got)
	}
	if w.Body.Len() == 0 {
		t.Fatal("empty xlsx body")
	}
	file, err := excelize.OpenReader(bytes.NewReader(w.Body.Bytes()))
	if err != nil {
		t.Fatalf("open xlsx: %v", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			t.Fatalf("close xlsx: %v", err)
		}
	}()
	if got := file.GetSheetList(); !reflect.DeepEqual(got, []string{"Invoices"}) {
		t.Fatalf("sheets = %#v, want Invoices", got)
	}
	assertXLSXCell(t, file, "Invoices", "A1", "document_id")
	assertXLSXCell(t, file, "Invoices", "B1", "filename")
	assertXLSXCell(t, file, "Invoices", "C1", "created_at")
	assertXLSXCell(t, file, "Invoices", "D1", "Total")
	assertXLSXCell(t, file, "Invoices", "A2", doc.ID.String())
	assertXLSXCell(t, file, "Invoices", "B2", "invoice.pdf")
	assertXLSXCell(t, file, "Invoices", "D2", "42")
}

func TestExportDatasetXLSXFiltersByCreatedAtRange(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "dataset-xlsx-date-range@example.com")
	schema := createDatasetSchema(t, db, user.ID, "Invoice", `{"type":"object","properties":{"total":{"type":"number"}}}`)
	dataset := createDatasetFixture(t, db, user.ID, schema.ID, "Invoices", `[{"path":"/total","key":"total","label":"Total"}]`)
	base := time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC)
	_ = setListOCRDocumentCreatedAt(t, db, createDatasetDocument(t, db, user.ID, schema.ID, "outside.pdf", `{"total":1}`), base.Add(-time.Hour))
	inside := setListOCRDocumentCreatedAt(t, db, createDatasetDocument(t, db, user.ID, schema.ID, "inside.pdf", `{"total":2}`), base)

	query := url.Values{}
	query.Set("user_id", user.ID)
	query.Set("format", "xlsx")
	query.Set("created_from", base.Format(time.RFC3339Nano))
	query.Set("created_to", base.Format(time.RFC3339Nano))
	req := newTestRequest(http.MethodGet, "/api/datasets/"+dataset.ID.String()+"/export?"+query.Encode(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	file, err := excelize.OpenReader(bytes.NewReader(w.Body.Bytes()))
	if err != nil {
		t.Fatalf("open xlsx: %v", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			t.Fatalf("close xlsx: %v", err)
		}
	}()
	rows, err := file.GetRows("Invoices")
	if err != nil {
		t.Fatalf("read xlsx rows: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("xlsx rows = %#v, want header and one filtered row", rows)
	}
	assertXLSXCell(t, file, "Invoices", "A2", inside.ID.String())
	assertXLSXCell(t, file, "Invoices", "B2", "inside.pdf")
}

func TestExportDatasetRejectsInvalidFormat(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "dataset-export-format@example.com")
	schema := createDatasetSchema(t, db, user.ID, "Invoice", `{"type":"object","properties":{"total":{"type":"number"}}}`)
	dataset := createDatasetFixture(t, db, user.ID, schema.ID, "Invoices", `[{"path":"/total","key":"total","label":"Total"}]`)

	req := newTestRequest(http.MethodGet, "/api/datasets/"+dataset.ID.String()+"/export?user_id="+url.QueryEscape(user.ID)+"&format=pdf", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if got.Error != "format must be csv or xlsx" {
		t.Fatalf("error = %q, want invalid format error", got.Error)
	}
}

func TestExportDatasetRejectsInvalidDateBounds(t *testing.T) {
	router, db := testRouter(t)
	user := createTestUser(t, db, "dataset-export-invalid-date@example.com")
	schema := createDatasetSchema(t, db, user.ID, "Invoice", `{"type":"object","properties":{"total":{"type":"number"}}}`)
	dataset := createDatasetFixture(t, db, user.ID, schema.ID, "Invoices", `[{"path":"/total","key":"total","label":"Total"}]`)

	tests := []struct {
		name      string
		query     string
		wantError string
	}{
		{name: "invalid from", query: "created_from=not-a-date", wantError: "invalid created_from"},
		{name: "invalid to", query: "created_to=not-a-date", wantError: "invalid created_to"},
		{name: "backwards", query: "created_from=2026-06-02T00%3A00%3A00Z&created_to=2026-06-01T00%3A00%3A00Z", wantError: "created_from must be before or equal to created_to"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := newTestRequest(http.MethodGet, "/api/datasets/"+dataset.ID.String()+"/export?user_id="+url.QueryEscape(user.ID)+"&format=csv&"+tc.query, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("status = %d body=%s, want 400", w.Code, w.Body.String())
			}
			var got ErrorResponse
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Fatalf("decode error: %v body=%s", err, w.Body.String())
			}
			if got.Error != tc.wantError {
				t.Fatalf("error = %q, want %q", got.Error, tc.wantError)
			}
		})
	}
}

func TestSanitizeDatasetCSVCellForSpreadsheetInjection(t *testing.T) {
	for _, raw := range []string{"=SUM(A1:A2)", "+1", "-1", "@cmd", "\tindent", "\rreturn"} {
		if got := sanitizeDatasetCSVCell(raw); got != "'"+raw {
			t.Fatalf("sanitizeDatasetCSVCell(%q) = %q, want %q", raw, got, "'"+raw)
		}
	}
	if got := sanitizeDatasetCSVCell("plain"); got != "plain" {
		t.Fatalf("plain cell = %q", got)
	}
}

func TestSanitizeExcelSheetNameTrimsTrailingApostropheAfterTruncation(t *testing.T) {
	raw := strings.Repeat("A", 30) + "'ignored"
	if got, want := sanitizeExcelSheetName(raw), strings.Repeat("A", 30); got != want {
		t.Fatalf("sheet name = %q, want %q", got, want)
	}
	if got := sanitizeExcelSheetName(" 'Budget' "); got != "Budget" {
		t.Fatalf("quoted sheet name = %q, want Budget", got)
	}
	if got := sanitizeExcelSheetName(strings.Repeat("'", 40)); got != "Dataset" {
		t.Fatalf("apostrophe-only sheet name = %q, want Dataset", got)
	}
}

func TestGetDatasetDoesNotExposeOtherUserSchemaName(t *testing.T) {
	router, db := testDatasetRouter(t)
	owner := createTestUser(t, db, "dataset-get-corrupt-owner@example.com")
	other := createTestUser(t, db, "dataset-get-corrupt-other@example.com")
	otherSchema := createDatasetSchema(t, db, other.ID, "Other Secret Schema", `{"type":"object","properties":{"total":{"type":"number"}}}`)
	fields := []ocr.DatasetField{{Path: "/total", Key: "total", Label: "Total"}}
	dataset := createDataset(t, db, owner.ID, otherSchema.ID, "Corrupt", fields, time.Now())

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodGet, "/api/datasets/"+dataset.ID.String()+"?user_id="+url.QueryEscape(owner.ID), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if strings.Contains(w.Body.String(), otherSchema.Name) {
		t.Fatalf("response leaked other user's schema name: %s", w.Body.String())
	}
}

func TestUpdateDatasetReplacesNameSchemaAndFields(t *testing.T) {
	router, db := testDatasetRouter(t)
	user := createTestUser(t, db, "dataset-update-owner@example.com")
	schemaA := createDatasetSchema(t, db, user.ID, "Invoice", `{"type":"object","properties":{"total":{"type":"number"}}}`)
	schemaB := createDatasetSchema(t, db, user.ID, "Receipt", `{
		"type":"object",
		"properties":{"merchant":{"type":"string"},"paid":{"type":"number"}}
	}`)
	oldFields := []ocr.DatasetField{{Path: "/total", Key: "total", Label: "Total"}}
	dataset := createDataset(t, db, user.ID, schemaA.ID, "Old", oldFields, time.Now())
	body := []byte(`{
		"name":" Updated ",
		"schema_id":"` + schemaB.ID.String() + `",
		"selected_fields":[
			{"path":"/merchant","key":"merchant","label":"Merchant"},
			{"path":"/paid","key":"paid","label":"Paid"}
		]
	}`)

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodPut, "/api/datasets/"+dataset.ID.String()+"?user_id="+url.QueryEscape(user.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got DatasetResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	wantFields := []DatasetFieldResponse{
		{Path: "/merchant", Key: "merchant", Label: "Merchant"},
		{Path: "/paid", Key: "paid", Label: "Paid"},
	}
	if got.ID != dataset.ID || got.Name != "Updated" || got.SchemaID != schemaB.ID || got.SchemaName != "Receipt" || got.FieldCount != 2 {
		t.Fatalf("unexpected response: %#v", got)
	}
	if !reflect.DeepEqual(got.SelectedFields, wantFields) {
		t.Fatalf("selected_fields = %#v, want %#v", got.SelectedFields, wantFields)
	}

	var stored ocr.Dataset
	if err := db.First(&stored, "id = ?", dataset.ID).Error; err != nil {
		t.Fatalf("load dataset: %v", err)
	}
	if stored.Name != "Updated" || stored.SchemaID != schemaB.ID {
		t.Fatalf("stored dataset = %#v, want updated name and schema", stored)
	}
	assertJSONEqual(t, json.RawMessage(stored.SelectedFields), `[
		{"path":"/merchant","key":"merchant","label":"Merchant"},
		{"path":"/paid","key":"paid","label":"Paid"}
	]`)
}

func TestUpdateDatasetRejectsOtherUserAndSystemSchemasWithoutMutation(t *testing.T) {
	router, db := testDatasetRouter(t)
	user := createTestUser(t, db, "dataset-update-schema-owner@example.com")
	other := createTestUser(t, db, "dataset-update-schema-other@example.com")
	ownerSchema := createDatasetSchema(t, db, user.ID, "Invoice", `{"type":"object","properties":{"total":{"type":"number"}}}`)
	otherSchema := createDatasetSchema(t, db, other.ID, "Other", `{"type":"object","properties":{"total":{"type":"number"}}}`)
	systemSchema := createSystemDatasetSchema(t, db, "System", `{"type":"object","properties":{"total":{"type":"number"}}}`)
	fields := []ocr.DatasetField{{Path: "/total", Key: "total", Label: "Total"}}

	for _, tc := range []struct {
		name     string
		schemaID uuid.UUID
	}{
		{name: "other user", schemaID: otherSchema.ID},
		{name: "system", schemaID: systemSchema.ID},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dataset := createDataset(t, db, user.ID, ownerSchema.ID, "Original "+tc.name, fields, time.Now())
			body := []byte(`{
				"name":"Changed",
				"schema_id":"` + tc.schemaID.String() + `",
				"selected_fields":[{"path":"/total","key":"total","label":"Total"}]
			}`)

			w := httptest.NewRecorder()
			req := newTestRequest(http.MethodPut, "/api/datasets/"+dataset.ID.String()+"?user_id="+url.QueryEscape(user.ID), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
			}
			var got ErrorResponse
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Fatalf("decode: %v", err)
			}
			if got.Error != "invalid schema_id" {
				t.Fatalf("error = %q, want invalid schema_id", got.Error)
			}

			var stored ocr.Dataset
			if err := db.First(&stored, "id = ?", dataset.ID).Error; err != nil {
				t.Fatalf("load dataset: %v", err)
			}
			if stored.Name != dataset.Name || stored.SchemaID != ownerSchema.ID {
				t.Fatalf("stored dataset = %#v, want unchanged", stored)
			}
			assertJSONEqual(t, json.RawMessage(stored.SelectedFields), `[{"path":"/total","key":"total","label":"Total"}]`)
		})
	}
}

func TestUpdateDatasetRejectsInvalidSelectedFieldsWithoutMutation(t *testing.T) {
	router, db := testDatasetRouter(t)
	user := createTestUser(t, db, "dataset-update-fields@example.com")
	schema := createDatasetSchema(t, db, user.ID, "Invoice", `{
		"type":"object",
		"properties":{"supplier":{"type":"string"},"total":{"type":"number"}}
	}`)
	oldFields := []ocr.DatasetField{{Path: "/total", Key: "total", Label: "Total"}}

	cases := []struct {
		name   string
		fields []ocr.DatasetField
	}{
		{name: "empty", fields: []ocr.DatasetField{}},
		{name: "invalid path", fields: []ocr.DatasetField{{Path: "/missing", Key: "missing", Label: "Missing"}}},
		{name: "duplicate key", fields: []ocr.DatasetField{
			{Path: "/supplier", Key: "value", Label: "Supplier"},
			{Path: "/total", Key: "value", Label: "Total"},
		}},
		{name: "too many fields", fields: manyDatasetFields(maxDatasetFieldCount + 1)},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dataset := createDataset(t, db, user.ID, schema.ID, "Original "+tc.name, oldFields, time.Now())
			body, err := json.Marshal(updateDatasetRequest{
				Name:           "Changed",
				SchemaID:       schema.ID.String(),
				SelectedFields: tc.fields,
			})
			if err != nil {
				t.Fatalf("marshal request: %v", err)
			}

			w := httptest.NewRecorder()
			req := newTestRequest(http.MethodPut, "/api/datasets/"+dataset.ID.String()+"?user_id="+url.QueryEscape(user.ID), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
			}
			var got ErrorResponse
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Fatalf("decode: %v", err)
			}
			if got.Error != "invalid selected_fields" {
				t.Fatalf("error = %q, want invalid selected_fields", got.Error)
			}

			var stored ocr.Dataset
			if err := db.First(&stored, "id = ?", dataset.ID).Error; err != nil {
				t.Fatalf("load dataset: %v", err)
			}
			if stored.Name != dataset.Name || stored.SchemaID != schema.ID {
				t.Fatalf("stored dataset = %#v, want unchanged", stored)
			}
			assertJSONEqual(t, json.RawMessage(stored.SelectedFields), `[{"path":"/total","key":"total","label":"Total"}]`)
		})
	}
}

func TestUpdateDatasetDoesNotSucceedWhenRowDisappearsBeforeMutation(t *testing.T) {
	router, db := testDatasetRouter(t)
	user := createTestUser(t, db, "dataset-update-race@example.com")
	schema := createDatasetSchema(t, db, user.ID, "Invoice", `{"type":"object","properties":{"total":{"type":"number"}}}`)
	fields := []ocr.DatasetField{{Path: "/total", Key: "total", Label: "Total"}}
	dataset := createDataset(t, db, user.ID, schema.ID, "Original", fields, time.Now())
	body := []byte(`{
		"name":"Changed",
		"schema_id":"` + schema.ID.String() + `",
		"selected_fields":[{"path":"/total","key":"total","label":"Total"}]
	}`)

	deletedDuringUpdate := false
	callbackName := "test:delete_dataset_before_update:" + uuid.NewString()
	if err := db.Callback().Update().Before("gorm:update").Register(callbackName, func(tx *gorm.DB) {
		if deletedDuringUpdate || tx.Statement.Schema == nil || tx.Statement.Schema.Table != "datasets" {
			return
		}
		deletedDuringUpdate = true
		if err := tx.Session(&gorm.Session{NewDB: true}).Exec("DELETE FROM datasets WHERE id = ?", dataset.ID).Error; err != nil {
			t.Errorf("delete dataset during update callback: %v", err)
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
	req := newTestRequest(http.MethodPut, "/api/datasets/"+dataset.ID.String()+"?user_id="+url.QueryEscape(user.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if !deletedDuringUpdate {
		t.Fatal("update callback did not run")
	}
	var stored ocr.Dataset
	if err := db.First(&stored, "id = ?", dataset.ID).Error; err != nil {
		t.Fatalf("load dataset: %v", err)
	}
	if stored.Name != "Original" || stored.SchemaID != schema.ID {
		t.Fatalf("stored dataset = %#v, want unchanged original dataset", stored)
	}
	assertJSONEqual(t, json.RawMessage(stored.SelectedFields), `[{"path":"/total","key":"total","label":"Total"}]`)
}

func TestDeleteDatasetRemovesOnlyDataset(t *testing.T) {
	router, db := testDatasetRouter(t)
	user := createTestUser(t, db, "dataset-delete-owner@example.com")
	schema := createDatasetSchema(t, db, user.ID, "Invoice", `{"type":"object","properties":{"total":{"type":"number"}}}`)
	fields := []ocr.DatasetField{{Path: "/total", Key: "total", Label: "Total"}}
	dataset := createDataset(t, db, user.ID, schema.ID, "Delete", fields, time.Now())
	otherDataset := createDataset(t, db, user.ID, schema.ID, "Keep", fields, time.Now().Add(time.Second))
	doc := createStoredOCRDocument(t, db, ocr.OCRDocument{UserID: &user.ID, SchemaID: &schema.ID})

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodDelete, "/api/datasets/"+dataset.ID.String()+"?user_id="+url.QueryEscape(user.ID), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	assertCount(t, db, &ocr.Dataset{}, "id = ?", dataset.ID, 0)
	assertCount(t, db, &ocr.Dataset{}, "id = ?", otherDataset.ID, 1)
	assertCount(t, db, &ocr.ExtractionSchema{}, "id = ?", schema.ID, 1)
	assertCount(t, db, &ocr.OCRDocument{}, "id = ?", doc.ID, 1)
}

func testDatasetRouter(t *testing.T) (*gin.Engine, *gorm.DB) {
	t.Helper()
	router, db := testRouter(t)
	if err := db.AutoMigrate(&ocr.Dataset{}); err != nil {
		t.Fatalf("migrate dataset: %v", err)
	}
	return router, db
}

func decodeDatasetListResponse(t *testing.T, w *httptest.ResponseRecorder) DatasetListResponse {
	t.Helper()
	var got DatasetListResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode dataset list response: %v body=%s", err, w.Body.String())
	}
	return got
}

func assertDatasetListIDs(t *testing.T, datasets []DatasetResponse, want ...uuid.UUID) {
	t.Helper()
	if len(datasets) != len(want) {
		t.Fatalf("dataset count = %d, want %d: %#v", len(datasets), len(want), datasets)
	}
	for i, id := range want {
		if datasets[i].ID != id {
			t.Fatalf("dataset[%d].id = %s, want %s; datasets=%#v", i, datasets[i].ID, id, datasets)
		}
	}
}

func assertDatasetRowIDs(t *testing.T, rows []DatasetRowResponse, want ...uuid.UUID) {
	t.Helper()
	if len(rows) != len(want) {
		t.Fatalf("dataset row count = %d, want %d: %#v", len(rows), len(want), rows)
	}
	for i, id := range want {
		if rows[i].DocumentID != id {
			t.Fatalf("row[%d].document_id = %s, want %s; rows=%#v", i, rows[i].DocumentID, id, rows)
		}
	}
}

func assertXLSXCell(t *testing.T, file *excelize.File, sheet string, cell string, want string) {
	t.Helper()
	got, err := file.GetCellValue(sheet, cell)
	if err != nil {
		t.Fatalf("get %s!%s: %v", sheet, cell, err)
	}
	if got != want {
		t.Fatalf("%s!%s = %q, want %q", sheet, cell, got, want)
	}
}

func manyDatasetFields(count int) []ocr.DatasetField {
	fields := make([]ocr.DatasetField, 0, count)
	for i := 0; i < count; i++ {
		fields = append(fields, ocr.DatasetField{
			Path:  "/total",
			Key:   "total_" + uuid.NewString(),
			Label: "Total",
		})
	}
	return fields
}

func createDatasetSchema(t *testing.T, db *gorm.DB, userID, name, schemaJSON string) ocr.ExtractionSchema {
	t.Helper()
	schema := ocr.ExtractionSchema{
		UserID:     &userID,
		Name:       name,
		SchemaJSON: datatypes.JSON([]byte(schemaJSON)),
		Strict:     true,
	}
	if err := db.Create(&schema).Error; err != nil {
		t.Fatalf("create dataset schema: %v", err)
	}
	return schema
}

func createSystemDatasetSchema(t *testing.T, db *gorm.DB, name, schemaJSON string) ocr.ExtractionSchema {
	t.Helper()
	schema := ocr.ExtractionSchema{
		Name:       name,
		SchemaJSON: datatypes.JSON([]byte(schemaJSON)),
		Strict:     true,
	}
	if err := db.Create(&schema).Error; err != nil {
		t.Fatalf("create system dataset schema: %v", err)
	}
	return schema
}

func createDataset(t *testing.T, db *gorm.DB, userID string, schemaID uuid.UUID, name string, fields []ocr.DatasetField, createdAt time.Time) ocr.Dataset {
	t.Helper()
	rawFields, err := json.Marshal(fields)
	if err != nil {
		t.Fatalf("marshal dataset fields: %v", err)
	}
	dataset := ocr.Dataset{
		UserID:         userID,
		SchemaID:       schemaID,
		Name:           name,
		SelectedFields: datatypes.JSON(rawFields),
		CreatedAt:      createdAt.UTC(),
		UpdatedAt:      createdAt.UTC(),
	}
	if err := db.Create(&dataset).Error; err != nil {
		t.Fatalf("create dataset fixture: %v", err)
	}
	return dataset
}

func createDatasetFixture(t *testing.T, db *gorm.DB, userID string, schemaID uuid.UUID, name string, selectedFieldsJSON string) ocr.Dataset {
	t.Helper()
	if err := db.AutoMigrate(&ocr.Dataset{}); err != nil {
		t.Fatalf("migrate dataset: %v", err)
	}
	dataset := ocr.Dataset{
		UserID:         userID,
		SchemaID:       schemaID,
		Name:           name,
		SelectedFields: datatypes.JSON([]byte(selectedFieldsJSON)),
	}
	if err := db.Create(&dataset).Error; err != nil {
		t.Fatalf("create dataset fixture: %v", err)
	}
	return dataset
}

func createDatasetDocument(t *testing.T, db *gorm.DB, userID string, schemaID uuid.UUID, filename string, annotationJSON string) ocr.OCRDocument {
	t.Helper()
	return createStoredOCRDocument(t, db, ocr.OCRDocument{
		UserID:           &userID,
		SchemaID:         &schemaID,
		OriginalFilename: filename,
		AnnotationJSON:   datatypes.JSON([]byte(annotationJSON)),
	})
}
