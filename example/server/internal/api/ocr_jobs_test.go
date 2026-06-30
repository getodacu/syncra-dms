package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"ai.ro/syncra/internal/auth"
	"ai.ro/syncra/internal/billing"
	"ai.ro/syncra/internal/ocr"
)

func testRouterWithOCRJobs(t *testing.T) (*gin.Engine, *gorm.DB, string, *bool) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	db := apiPostgresTx(t)
	storageDir := t.TempDir()
	fileDir := filepath.Join(storageDir, "ocr-files")
	called := false
	h := &Handler{
		DB: db,
		OCR: func(ctx context.Context, input OCRProcessInput) (*MistralOCRResponse, []byte, error) {
			called = true
			return nil, nil, nil
		},
		MaxUploadBytes:   20 << 20,
		StorageDir:       storageDir,
		InternalAPIToken: testInternalAPIToken,
	}
	return NewRouter(h), db, fileDir, &called
}

func ocrJobTestModels() []any {
	return []any{
		&auth.User{},
		&auth.APIKey{},
		&billing.BillingOrder{},
		&billing.CreditBucket{},
		&billing.CreditLedgerEntry{},
		&ocr.ExtractionSchema{},
		&ocr.OCRDocument{},
		&ocr.OCRJob{},
	}
}

func assertOCRJobResponseKeys(t *testing.T, body []byte) map[string]any {
	t.Helper()
	var raw map[string]any
	if err := json.Unmarshal(body, &raw); err != nil {
		t.Fatalf("decode raw response: %v", err)
	}
	want := map[string]struct{}{
		"id":                {},
		"created_at":        {},
		"original_filename": {},
		"mime_type":         {},
		"status":            {},
		"file_size":         {},
		"page_count":        {},
		"document_id":       {},
		"has_inline_schema": {},
	}
	if len(raw) != len(want) {
		t.Fatalf("response keys = %#v, want only %#v", raw, want)
	}
	for key := range want {
		if _, ok := raw[key]; !ok {
			t.Fatalf("response missing key %q: %s", key, string(body))
		}
	}
	return raw
}

func assertRawObjectKeys(t *testing.T, raw map[string]any, want ...string) {
	t.Helper()
	if len(raw) != len(want) {
		t.Fatalf("response keys = %#v, want only %#v", raw, want)
	}
	for _, key := range want {
		if _, ok := raw[key]; !ok {
			t.Fatalf("response missing key %q: %#v", key, raw)
		}
	}
}

func assertRawJSONValue(t *testing.T, value any, want string) {
	t.Helper()
	raw, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("encode raw JSON value: %v", err)
	}
	assertJSONEqual(t, raw, want)
}

func TestPublicOCRJobDocumentDataDefaultsMissingRawFields(t *testing.T) {
	pages, annotation := publicOCRJobDocumentData([]byte(`{"model":"mistral-ocr-test"}`))
	if len(pages) != 0 {
		t.Fatalf("pages = %#v, want empty array", pages)
	}
	assertJSONEqual(t, annotation, `null`)

	pages, annotation = publicOCRJobDocumentData([]byte(`{"pages":null,"document_annotation":{"total":10}}`))
	if len(pages) != 0 {
		t.Fatalf("pages = %#v, want empty array", pages)
	}
	assertJSONEqual(t, annotation, `{"total":10}`)
}

func createStoredOCRJob(t *testing.T, db *gorm.DB, job ocr.OCRJob) ocr.OCRJob {
	t.Helper()
	if job.OriginalFilename == "" {
		job.OriginalFilename = "scan.png"
	}
	if job.MimeType == "" {
		job.MimeType = "image/png"
	}
	if job.FileSize == 0 {
		job.FileSize = int64(len(validPNGBytes()))
	}
	if job.PageCount == 0 {
		job.PageCount = 1
	}
	if job.FilePath == "" {
		job.FilePath = filepath.Join(t.TempDir(), uuid.NewString()+".png")
		if err := os.WriteFile(job.FilePath, validPNGBytes(), 0o600); err != nil {
			t.Fatalf("write existing job file: %v", err)
		}
	}
	if job.DocumentHash == "" {
		hash, err := computeDocumentHash(validPNGBytes(), nil, false)
		if err != nil {
			t.Fatalf("compute document hash: %v", err)
		}
		job.DocumentHash = hash
	}
	if job.Status == "" {
		job.Status = ocr.OCRJobStatusQueued
	}
	if err := db.Create(&job).Error; err != nil {
		t.Fatalf("create stored OCR job: %v", err)
	}
	return job
}

func createStoredJobSchema(t *testing.T, db *gorm.DB, userID string, name string) ocr.ExtractionSchema {
	t.Helper()
	schema := ocr.ExtractionSchema{
		UserID:      &userID,
		Name:        name,
		Description: "Test schema",
		SchemaJSON:  datatypes.JSON([]byte(`{"type":"object"}`)),
		Strict:      true,
	}
	if err := db.Create(&schema).Error; err != nil {
		t.Fatalf("create stored job schema: %v", err)
	}
	return schema
}

func grantTestCredits(t *testing.T, db *gorm.DB, userID string, credits int) {
	t.Helper()
	bucket := billing.CreditBucket{
		UserID:           userID,
		SourceType:       billing.CreditSourceAdjustment,
		CreditsGranted:   credits,
		CreditsRemaining: credits,
		ValidFrom:        time.Now().UTC().Add(-time.Hour),
	}
	if err := db.Create(&bucket).Error; err != nil {
		t.Fatalf("create credit bucket: %v", err)
	}
	entry := billing.CreditLedgerEntry{
		UserID:         userID,
		BucketID:       &bucket.ID,
		EntryType:      billing.CreditLedgerEntryAdjustment,
		CreditsDelta:   credits,
		IdempotencyKey: "test_credit:" + uuid.NewString(),
	}
	if err := db.Create(&entry).Error; err != nil {
		t.Fatalf("create credit ledger entry: %v", err)
	}
}

func createTestPublicAPIKey(t *testing.T, db *gorm.DB, userID string, expiresAt *time.Time) string {
	t.Helper()
	secret := "public-api-key-" + uuid.NewString()
	key := auth.APIKey{
		UserID:    userID,
		Name:      "Public API key",
		KeyHash:   auth.HashAPIKey(secret),
		KeyPrefix: secret[:8],
		ExpiresAt: expiresAt,
	}
	if err := db.Create(&key).Error; err != nil {
		t.Fatalf("create public API key: %v", err)
	}
	return secret
}

func authorizePublicAPIRequest(req *http.Request, secret string) *http.Request {
	req.Header.Set("Authorization", "Bearer "+secret)
	return req
}

func creditedOCRJobFields(t *testing.T, db *gorm.DB, credits int, fields map[string]string) map[string]string {
	t.Helper()
	user := createTestUser(t, db, "ocr-job-"+uuid.NewString()+"@example.com")
	grantTestCredits(t, db, user.ID, credits)
	return ocrJobFieldsForUser(user.ID, fields)
}

func ocrJobFieldsForUser(userID string, fields map[string]string) map[string]string {
	out := make(map[string]string, len(fields)+1)
	out["user_id"] = userID
	for key, value := range fields {
		out[key] = value
	}
	return out
}

func setStoredOCRJobCreatedAt(t *testing.T, db *gorm.DB, job ocr.OCRJob, createdAt time.Time) ocr.OCRJob {
	t.Helper()
	createdAt = createdAt.UTC()
	if err := db.Model(&job).UpdateColumn("created_at", createdAt).Error; err != nil {
		t.Fatalf("set OCR job created_at: %v", err)
	}
	job.CreatedAt = createdAt
	return job
}

func decodeOCRJobListResponse(t *testing.T, w *httptest.ResponseRecorder) OCRJobListResponse {
	t.Helper()
	var got OCRJobListResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode OCR job list response: %v body=%s", err, w.Body.String())
	}
	return got
}

func decodeDeleteOCRJobsResponse(t *testing.T, w *httptest.ResponseRecorder) DeleteOCRJobsResponse {
	t.Helper()
	var got DeleteOCRJobsResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode OCR job delete response: %v body=%s", err, w.Body.String())
	}
	return got
}

func assertOCRJobListIDs(t *testing.T, jobs []OCRJobListItemResponse, want ...uuid.UUID) {
	t.Helper()
	if len(jobs) != len(want) {
		t.Fatalf("job count = %d, want %d: %#v", len(jobs), len(want), jobs)
	}
	for i, id := range want {
		if jobs[i].ID != id {
			t.Fatalf("job[%d].id = %s, want %s; jobs=%#v", i, jobs[i].ID, id, jobs)
		}
	}
}

func assertOCRJobCount(t *testing.T, db *gorm.DB, want int64) {
	t.Helper()
	var count int64
	if err := db.Model(&ocr.OCRJob{}).Count(&count).Error; err != nil {
		t.Fatalf("count OCR jobs: %v", err)
	}
	if count != want {
		t.Fatalf("OCR job count = %d, want %d", count, want)
	}
}

func assertPersistedQueuedOCRJob(t *testing.T, db *gorm.DB, got OCRJobResponse, hash, fileDir string) ocr.OCRJob {
	t.Helper()
	if got.Status != string(ocr.OCRJobStatusQueued) {
		t.Fatalf("response status = %q, want %q", got.Status, ocr.OCRJobStatusQueued)
	}
	var job ocr.OCRJob
	if err := db.First(&job, "id = ?", got.ID).Error; err != nil {
		t.Fatalf("load returned OCR job %s: %v", got.ID, err)
	}
	if job.Status != ocr.OCRJobStatusQueued {
		t.Fatalf("stored status = %s, want %s", job.Status, ocr.OCRJobStatusQueued)
	}
	if job.DocumentID != nil {
		t.Fatalf("stored document_id = %#v, want nil", job.DocumentID)
	}
	if job.DocumentHash != hash {
		t.Fatalf("stored document_hash = %q, want %q", job.DocumentHash, hash)
	}
	if job.FilePath == "" {
		t.Fatal("stored file path is empty")
	}
	ext, ok := ocrJobFileExtension(got.MimeType)
	if !ok {
		t.Fatalf("unsupported response mime type %q", got.MimeType)
	}
	wantFilePath := filepath.Join(fileDir, got.ID.String()+ext)
	if job.FilePath != wantFilePath {
		t.Fatalf("stored file path = %q, want %q", job.FilePath, wantFilePath)
	}
	return job
}

func assertDirEmpty(t *testing.T, dir string) {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		t.Fatalf("read dir: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("dir entries = %v, want none", entryNames(entries))
	}
}

func TestCreateOCRJobQueuesUploadWithoutProcessing(t *testing.T) {
	router, db, _, called := testRouterWithOCRJobs(t)
	file := validPNGBytes()
	fields := creditedOCRJobFields(t, db, 5, nil)
	req := multipartRequestForPath(t, "/api/ocr/jobs", fields, "scan.png", file)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if *called {
		t.Fatal("OCR processor was called")
	}

	var got OCRJobResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ID == uuid.Nil || got.Status != string(ocr.OCRJobStatusQueued) {
		t.Fatalf("unexpected response: %#v", got)
	}
	if got.FileSize != int64(len(file)) {
		t.Fatalf("file_size = %d, want %d", got.FileSize, len(file))
	}
	if got.PageCount != 1 {
		t.Fatalf("page_count = %d, want 1", got.PageCount)
	}
	raw := assertOCRJobResponseKeys(t, w.Body.Bytes())
	if raw["document_id"] != nil {
		t.Fatalf("document_id = %#v, want null", raw["document_id"])
	}

	var job ocr.OCRJob
	if err := db.First(&job, "id = ?", got.ID).Error; err != nil {
		t.Fatalf("load OCR job: %v", err)
	}
	if job.FilePath == "" {
		t.Fatal("stored file path is empty")
	}
	stored, err := os.ReadFile(job.FilePath)
	if err != nil {
		t.Fatalf("read stored file: %v", err)
	}
	if !bytes.Equal(stored, file) {
		t.Fatalf("stored file = %v, want %v", stored, file)
	}
	if job.Status != ocr.OCRJobStatusQueued || job.OriginalFilename != "scan.png" || job.MimeType != "image/png" || job.FileSize != int64(len(file)) {
		t.Fatalf("unexpected stored job: %#v", job)
	}
	if job.PageCount != 1 {
		t.Fatalf("stored page_count = %d, want 1", job.PageCount)
	}
	if job.DocumentHash == "" {
		t.Fatal("stored document_hash is empty")
	}
	if len(job.InlineSchemaJSON) > 0 {
		t.Fatal("stored inline schema is not empty")
	}
	assertDirEntryNames(t, filepath.Dir(job.FilePath), job.ID.String()+".claim", job.ID.String()+".png")
}

func TestCreatePublicOCRJobQueuesUploadForAPIKeyUser(t *testing.T) {
	router, db, fileDir, called := testRouterWithOCRJobs(t)
	user := createTestUser(t, db, "public-ocr-job-owner@example.com")
	grantTestCredits(t, db, user.ID, 5)
	secret := createTestPublicAPIKey(t, db, user.ID, nil)
	file := validPNGBytes()
	req := multipartRequestForPath(t, "/v1/ocr/jobs", nil, "scan.png", file)
	authorizePublicAPIRequest(req, secret)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if *called {
		t.Fatal("OCR processor was called")
	}
	var got OCRJobResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	hash, err := computeDocumentHash(file, nil, false)
	if err != nil {
		t.Fatalf("compute document hash: %v", err)
	}
	job := assertPersistedQueuedOCRJob(t, db, got, hash, fileDir)
	if job.UserID == nil || *job.UserID != user.ID {
		t.Fatalf("stored user_id = %#v, want %s", job.UserID, user.ID)
	}
	assertOCRJobResponseKeys(t, w.Body.Bytes())
}

func TestCreatePublicOCRJobValidationErrors(t *testing.T) {
	router, db, _, _ := testRouterWithOCRJobs(t)
	user := createTestUser(t, db, "public-ocr-job-validation@example.com")
	secret := createTestPublicAPIKey(t, db, user.ID, nil)

	for _, tt := range []struct {
		name     string
		fields   map[string]string
		filename string
		content  []byte
		wantCode int
	}{
		{name: "missing file", wantCode: http.StatusBadRequest},
		{name: "unsupported file", filename: "notes.txt", content: []byte("hello"), wantCode: http.StatusBadRequest},
		{name: "invalid schema", fields: map[string]string{"schema": "not-json"}, filename: "scan.png", content: validPNGBytes(), wantCode: http.StatusBadRequest},
		{name: "insufficient credits", filename: "scan.png", content: validPNGBytes(), wantCode: http.StatusPaymentRequired},
		{name: "form user id", fields: map[string]string{"user_id": user.ID}, filename: "scan.png", content: validPNGBytes(), wantCode: http.StatusBadRequest},
	} {
		t.Run(tt.name, func(t *testing.T) {
			req := multipartRequestForPath(t, "/v1/ocr/jobs", tt.fields, tt.filename, tt.content)
			authorizePublicAPIRequest(req, secret)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.wantCode {
				t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
			}
		})
	}
}

func TestCreatePublicOCRJobRejectsQueryUserID(t *testing.T) {
	router, db, _, _ := testRouterWithOCRJobs(t)
	user := createTestUser(t, db, "public-ocr-job-query-user@example.com")
	secret := createTestPublicAPIKey(t, db, user.ID, nil)
	req := multipartRequestForPath(t, "/v1/ocr/jobs?user_id="+url.QueryEscape(user.ID), nil, "scan.png", validPNGBytes())
	authorizePublicAPIRequest(req, secret)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var count int64
	if err := db.Model(&ocr.OCRJob{}).Count(&count).Error; err != nil {
		t.Fatalf("count OCR jobs: %v", err)
	}
	if count != 0 {
		t.Fatalf("OCR job count = %d, want 0", count)
	}
}

func TestCreatePublicOCRJobRejectsOtherUserSchema(t *testing.T) {
	router, db, _, _ := testRouterWithOCRJobs(t)
	user := createTestUser(t, db, "public-ocr-job-schema-owner@example.com")
	other := createTestUser(t, db, "public-ocr-job-schema-other@example.com")
	secret := createTestPublicAPIKey(t, db, user.ID, nil)
	schema := createStoredJobSchema(t, db, other.ID, "Other schema")
	req := multipartRequestForPath(t, "/v1/ocr/jobs", map[string]string{"schema_id": schema.ID.String()}, "scan.png", validPNGBytes())
	authorizePublicAPIRequest(req, secret)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	assertOCRJobCount(t, db, 0)
}

func TestPublicOCRJobAPIRejectsInvalidAPIKeys(t *testing.T) {
	router, db, _, _ := testRouterWithOCRJobs(t)
	user := createTestUser(t, db, "public-ocr-job-auth@example.com")
	expiredAt := time.Now().UTC().Add(-time.Minute)
	expiredSecret := createTestPublicAPIKey(t, db, user.ID, &expiredAt)

	for _, tt := range []struct {
		name   string
		header string
	}{
		{name: "missing"},
		{name: "malformed", header: "Basic " + expiredSecret},
		{name: "unknown", header: "Bearer public-api-key-" + uuid.NewString()},
		{name: "expired", header: "Bearer " + expiredSecret},
	} {
		t.Run(tt.name, func(t *testing.T) {
			req := newTestRequest(http.MethodGet, "/v1/ocr/jobs/"+uuid.NewString(), nil)
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != http.StatusUnauthorized {
				t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
			}
		})
	}
}

func TestPublicOCRJobAPIAcceptsRawAPIKeyAuthorization(t *testing.T) {
	router, db, _, _ := testRouterWithOCRJobs(t)
	user := createTestUser(t, db, "public-ocr-job-raw-auth@example.com")
	secret := createTestPublicAPIKey(t, db, user.ID, nil)
	job := createStoredOCRJob(t, db, ocr.OCRJob{
		UserID:       &user.ID,
		DocumentHash: "public-raw-auth-hash",
		Status:       ocr.OCRJobStatusQueued,
	})
	req := newTestRequest(http.MethodGet, "/v1/ocr/jobs/"+job.ID.String(), nil)
	req.Header.Set("Authorization", secret)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got OCRJobResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ID != job.ID {
		t.Fatalf("job id = %s, want %s", got.ID, job.ID)
	}
}

func TestCreateOCRJobCreatesNewJobForCompletedDuplicate(t *testing.T) {
	router, db, fileDir, _ := testRouterWithOCRJobs(t)
	user := createTestUser(t, db, "completed-duplicate-owner@example.com")
	grantTestCredits(t, db, user.ID, 5)
	file := validPNGBytes()
	hash, err := computeDocumentHash(file, nil, false)
	if err != nil {
		t.Fatalf("compute document hash: %v", err)
	}
	doc := ocr.OCRDocument{
		UserID:           &user.ID,
		OriginalFilename: "scan.png",
		MimeType:         "image/png",
		FileSize:         int64(len(file)),
		DocumentHash:     hash,
		Markdown:         "# Existing",
		RawResponseJSON:  datatypes.JSON([]byte(`{"pages":[{"index":0,"markdown":"# Existing"}]}`)),
	}
	if err := db.Create(&doc).Error; err != nil {
		t.Fatalf("create OCR document: %v", err)
	}
	existing := createStoredOCRJob(t, db, ocr.OCRJob{
		UserID:       &user.ID,
		DocumentHash: hash,
		DocumentID:   &doc.ID,
		Status:       ocr.OCRJobStatusCompleted,
	})

	req := multipartRequestForPath(t, "/api/ocr/jobs", ocrJobFieldsForUser(user.ID, nil), "scan.png", file)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got OCRJobResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ID == existing.ID {
		t.Fatalf("job id = existing %s, want new job", got.ID)
	}
	if got.DocumentID != nil {
		t.Fatalf("document_id = %#v, want nil for queued new job", got.DocumentID)
	}
	assertPersistedQueuedOCRJob(t, db, got, hash, fileDir)
	assertOCRJobCount(t, db, 2)
	assertDirEntryNames(t, fileDir, got.ID.String()+".claim", got.ID.String()+".png")
}

func TestCreateOCRJobCreatesNewJobWhenDuplicateDocumentWasDeleted(t *testing.T) {
	router, db, fileDir, _ := testRouterWithOCRJobs(t)
	user := createTestUser(t, db, "deleted-duplicate-owner@example.com")
	grantTestCredits(t, db, user.ID, 5)
	file := validPNGBytes()
	hash, err := computeDocumentHash(file, nil, false)
	if err != nil {
		t.Fatalf("compute document hash: %v", err)
	}
	doc := ocr.OCRDocument{
		UserID:           &user.ID,
		OriginalFilename: "scan.png",
		MimeType:         "image/png",
		FileSize:         int64(len(file)),
		DocumentHash:     hash,
		Markdown:         "# Deleted",
		RawResponseJSON:  datatypes.JSON([]byte(`{"pages":[{"index":0,"markdown":"# Deleted"}]}`)),
	}
	if err := db.Create(&doc).Error; err != nil {
		t.Fatalf("create OCR document: %v", err)
	}
	existing := createStoredOCRJob(t, db, ocr.OCRJob{
		UserID:       &user.ID,
		DocumentHash: hash,
		DocumentID:   &doc.ID,
		Status:       ocr.OCRJobStatusCompleted,
	})
	if err := db.Delete(&doc).Error; err != nil {
		t.Fatalf("soft delete OCR document: %v", err)
	}
	var deletedDoc ocr.OCRDocument
	if err := db.Unscoped().First(&deletedDoc, "id = ?", doc.ID).Error; err != nil {
		t.Fatalf("load soft-deleted OCR document: %v", err)
	}
	if !deletedDoc.DeletedAt.Valid {
		t.Fatal("deleted_at is not valid after soft delete")
	}

	req := multipartRequestForPath(t, "/api/ocr/jobs", ocrJobFieldsForUser(user.ID, nil), "scan.png", file)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got OCRJobResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ID == existing.ID {
		t.Fatalf("job id = existing %s, want new job", got.ID)
	}
	if got.DocumentID != nil {
		t.Fatalf("document_id = %#v, want nil for queued new job", got.DocumentID)
	}
	assertPersistedQueuedOCRJob(t, db, got, hash, fileDir)
	assertOCRJobCount(t, db, 2)
	assertDirEntryNames(t, fileDir, got.ID.String()+".claim", got.ID.String()+".png")
}

func TestCreateOCRJobCreatesNewJobForActiveDuplicate(t *testing.T) {
	cases := []struct {
		name   string
		status ocr.OCRJobStatus
	}{
		{name: "queued", status: ocr.OCRJobStatusQueued},
		{name: "processing", status: ocr.OCRJobStatusProcessing},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			router, db, fileDir, _ := testRouterWithOCRJobs(t)
			user := createTestUser(t, db, "active-duplicate-owner-"+tc.name+"@example.com")
			grantTestCredits(t, db, user.ID, 10)
			file := validPNGBytes()
			hash, err := computeDocumentHash(file, nil, false)
			if err != nil {
				t.Fatalf("compute document hash: %v", err)
			}
			existing := createStoredOCRJob(t, db, ocr.OCRJob{
				UserID:       &user.ID,
				DocumentHash: hash,
				Status:       tc.status,
			})

			req := multipartRequestForPath(t, "/api/ocr/jobs", ocrJobFieldsForUser(user.ID, nil), "scan.png", file)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusAccepted {
				t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
			}
			var got OCRJobResponse
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Fatalf("decode: %v", err)
			}
			if got.ID == existing.ID {
				t.Fatalf("job id = existing %s, want new job", got.ID)
			}
			if got.DocumentID != nil {
				t.Fatalf("document_id = %#v, want nil", got.DocumentID)
			}
			assertPersistedQueuedOCRJob(t, db, got, hash, fileDir)
			assertOCRJobCount(t, db, 2)
			assertDirEntryNames(t, fileDir, got.ID.String()+".claim", got.ID.String()+".png")
		})
	}
}

func TestCreateOCRJobIgnoresFailedDuplicate(t *testing.T) {
	router, db, fileDir, _ := testRouterWithOCRJobs(t)
	user := createTestUser(t, db, "failed-duplicate-owner@example.com")
	grantTestCredits(t, db, user.ID, 5)
	file := validPNGBytes()
	hash, err := computeDocumentHash(file, nil, false)
	if err != nil {
		t.Fatalf("compute document hash: %v", err)
	}
	failed := createStoredOCRJob(t, db, ocr.OCRJob{
		UserID:       &user.ID,
		DocumentHash: hash,
		Status:       ocr.OCRJobStatusFailed,
	})

	req := multipartRequestForPath(t, "/api/ocr/jobs", ocrJobFieldsForUser(user.ID, nil), "scan.png", file)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got OCRJobResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ID == failed.ID {
		t.Fatalf("job id = failed duplicate %s, want new job", got.ID)
	}
	if got.DocumentID != nil {
		t.Fatalf("document_id = %#v, want nil", got.DocumentID)
	}
	assertPersistedQueuedOCRJob(t, db, got, hash, fileDir)
	assertOCRJobCount(t, db, 2)
	assertDirEntryNames(t, fileDir, got.ID.String()+".claim", got.ID.String()+".png")
}

func TestCreateOCRJobDoesNotReuseExpiredDuplicate(t *testing.T) {
	router, db, fileDir, _ := testRouterWithOCRJobs(t)
	user := createTestUser(t, db, "expired-duplicate-owner@example.com")
	grantTestCredits(t, db, user.ID, 5)
	file := validPNGBytes()
	hash, err := computeDocumentHash(file, nil, false)
	if err != nil {
		t.Fatalf("compute document hash: %v", err)
	}
	existing := createStoredOCRJob(t, db, ocr.OCRJob{
		UserID:       &user.ID,
		DocumentHash: hash,
		Status:       ocr.OCRJobStatusQueued,
	})
	if err := db.Model(&existing).Update("created_at", time.Now().Add(-25*time.Hour)).Error; err != nil {
		t.Fatalf("age existing job: %v", err)
	}

	req := multipartRequestForPath(t, "/api/ocr/jobs", ocrJobFieldsForUser(user.ID, nil), "scan.png", file)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got OCRJobResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ID == existing.ID {
		t.Fatalf("job id = expired duplicate %s, want new job", got.ID)
	}
	if got.DocumentID != nil {
		t.Fatalf("document_id = %#v, want nil", got.DocumentID)
	}
	assertPersistedQueuedOCRJob(t, db, got, hash, fileDir)
	assertOCRJobCount(t, db, 2)
	assertDirEntryNames(t, fileDir, got.ID.String()+".claim", got.ID.String()+".png")
}

func TestCreateOCRJobDoesNotReuseDifferentUserScope(t *testing.T) {
	router, db, fileDir, _ := testRouterWithOCRJobs(t)
	user := createTestUser(t, db, "dedupe-owner@example.com")
	requestUser := createTestUser(t, db, "dedupe-requester@example.com")
	grantTestCredits(t, db, requestUser.ID, 5)
	file := validPNGBytes()
	hash, err := computeDocumentHash(file, nil, false)
	if err != nil {
		t.Fatalf("compute document hash: %v", err)
	}
	existing := createStoredOCRJob(t, db, ocr.OCRJob{
		UserID:       &user.ID,
		DocumentHash: hash,
		Status:       ocr.OCRJobStatusQueued,
	})

	req := multipartRequestForPath(t, "/api/ocr/jobs", ocrJobFieldsForUser(requestUser.ID, nil), "scan.png", file)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got OCRJobResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ID == existing.ID {
		t.Fatalf("job id = different user duplicate %s, want new job", got.ID)
	}
	if got.DocumentID != nil {
		t.Fatalf("document_id = %#v, want nil", got.DocumentID)
	}
	assertPersistedQueuedOCRJob(t, db, got, hash, fileDir)
	assertOCRJobCount(t, db, 2)
	assertDirEntryNames(t, fileDir, got.ID.String()+".claim", got.ID.String()+".png")
}

func TestCreateOCRJobDoesNotReuseDifferentSchema(t *testing.T) {
	router, db, fileDir, _ := testRouterWithOCRJobs(t)
	user := createTestUser(t, db, "different-schema-owner@example.com")
	grantTestCredits(t, db, user.ID, 5)
	file := validPNGBytes()
	noSchemaHash, err := computeDocumentHash(file, nil, false)
	if err != nil {
		t.Fatalf("compute no-schema hash: %v", err)
	}
	createStoredOCRJob(t, db, ocr.OCRJob{
		UserID:       &user.ID,
		DocumentHash: noSchemaHash,
		Status:       ocr.OCRJobStatusQueued,
	})
	schema := ocr.ExtractionSchema{Name: "invoice", SchemaJSON: []byte(`{"type":"object"}`), Strict: true}
	if err := db.Create(&schema).Error; err != nil {
		t.Fatalf("create schema: %v", err)
	}
	hash, err := computeDocumentHash(file, json.RawMessage(schema.SchemaJSON), schema.Strict)
	if err != nil {
		t.Fatalf("compute schema hash: %v", err)
	}

	req := multipartRequestForPath(t, "/api/ocr/jobs", ocrJobFieldsForUser(user.ID, map[string]string{"schema_id": schema.ID.String()}), "scan.png", file)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got OCRJobResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.DocumentID != nil {
		t.Fatalf("document_id = %#v, want nil", got.DocumentID)
	}
	assertPersistedQueuedOCRJob(t, db, got, hash, fileDir)
	assertOCRJobCount(t, db, 2)
	assertDirEntryNames(t, fileDir, got.ID.String()+".claim", got.ID.String()+".png")
}

func TestCreateOCRJobCountsPDFPages(t *testing.T) {
	router, db, _, _ := testRouterWithOCRJobs(t)
	file := twoPagePDFBytes()
	fields := creditedOCRJobFields(t, db, 5, nil)
	req := multipartRequestForPath(t, "/api/ocr/jobs", fields, "invoice.pdf", file)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got OCRJobResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.PageCount != 2 {
		t.Fatalf("page_count = %d, want 2", got.PageCount)
	}

	var job ocr.OCRJob
	if err := db.First(&job, "id = ?", got.ID).Error; err != nil {
		t.Fatalf("load OCR job: %v", err)
	}
	if job.PageCount != 2 {
		t.Fatalf("stored page_count = %d, want 2", job.PageCount)
	}
}

func TestCreateOCRJobRejectsPageLimitWithoutSchema(t *testing.T) {
	router, db, fileDir, _ := testRouterWithOCRJobs(t)
	file := pageCountPDFBytes(1001)
	fields := creditedOCRJobFields(t, db, 1001, nil)
	req := multipartRequestForPath(t, "/api/ocr/jobs", fields, "large.pdf", file)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if got.Error != "document must have at most 1000 pages without a schema" {
		t.Fatalf("error = %q", got.Error)
	}
	assertOCRJobCount(t, db, 0)
	assertDirEmpty(t, fileDir)
}

func TestCreateOCRJobAllowsPageLimitBoundaryWithoutSchema(t *testing.T) {
	router, db, fileDir, _ := testRouterWithOCRJobs(t)
	file := pageCountPDFBytes(1000)
	fields := creditedOCRJobFields(t, db, 1000, nil)
	hash, err := computeDocumentHash(file, nil, false)
	if err != nil {
		t.Fatalf("compute document hash: %v", err)
	}
	req := multipartRequestForPath(t, "/api/ocr/jobs", fields, "large.pdf", file)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got OCRJobResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.PageCount != 1000 {
		t.Fatalf("page_count = %d, want 1000", got.PageCount)
	}
	assertPersistedQueuedOCRJob(t, db, got, hash, fileDir)
	assertOCRJobCount(t, db, 1)
}

func TestCreateOCRJobRejectsPageLimitWithInlineSchema(t *testing.T) {
	router, db, fileDir, _ := testRouterWithOCRJobs(t)
	file := pageCountPDFBytes(151)
	fields := creditedOCRJobFields(t, db, 151, map[string]string{"schema": `{"type":"object"}`})
	req := multipartRequestForPath(t, "/api/ocr/jobs", fields, "large.pdf", file)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if got.Error != "document must have at most 150 pages with a schema" {
		t.Fatalf("error = %q", got.Error)
	}
	assertOCRJobCount(t, db, 0)
	assertDirEmpty(t, fileDir)
}

func TestCreateOCRJobAllowsPageLimitBoundaryWithInlineSchema(t *testing.T) {
	router, db, fileDir, _ := testRouterWithOCRJobs(t)
	inlineSchema := json.RawMessage(`{"type":"object"}`)
	file := pageCountPDFBytes(150)
	fields := creditedOCRJobFields(t, db, 150, map[string]string{"schema": string(inlineSchema)})
	hash, err := computeDocumentHash(file, inlineSchema, true)
	if err != nil {
		t.Fatalf("compute document hash: %v", err)
	}
	req := multipartRequestForPath(t, "/api/ocr/jobs", fields, "large.pdf", file)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got OCRJobResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.PageCount != 150 {
		t.Fatalf("page_count = %d, want 150", got.PageCount)
	}
	job := assertPersistedQueuedOCRJob(t, db, got, hash, fileDir)
	if len(job.InlineSchemaJSON) == 0 {
		t.Fatal("stored inline schema is empty")
	}
	assertOCRJobCount(t, db, 1)
}

func TestCreatePublicOCRJobRejectsPageLimitWithSavedSchema(t *testing.T) {
	router, db, fileDir, _ := testRouterWithOCRJobs(t)
	user := createTestUser(t, db, "public-ocr-job-page-limit@example.com")
	grantTestCredits(t, db, user.ID, 151)
	secret := createTestPublicAPIKey(t, db, user.ID, nil)
	schema := createStoredJobSchema(t, db, user.ID, "Large schema job")
	file := pageCountPDFBytes(151)
	req := multipartRequestForPath(t, "/v1/ocr/jobs", map[string]string{"schema_id": schema.ID.String()}, "large.pdf", file)
	authorizePublicAPIRequest(req, secret)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if got.Error != "document must have at most 150 pages with a schema" {
		t.Fatalf("error = %q", got.Error)
	}
	assertOCRJobCount(t, db, 0)
	assertDirEmpty(t, fileDir)
}

func TestCreateOCRJobRejectsInsufficientCreditsForPageCount(t *testing.T) {
	router, db, fileDir, _ := testRouterWithOCRJobs(t)
	user := createTestUser(t, db, "ocr-job-low-balance@example.com")
	grantTestCredits(t, db, user.ID, 1)

	req := multipartRequestForPath(t, "/api/ocr/jobs", map[string]string{"user_id": user.ID}, "invoice.pdf", twoPagePDFBytes())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusPaymentRequired {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeOCRJobPaymentRequiredResponse(t, w)
	if got.Error != "insufficient credits" || got.RequiredCredits != 2 || got.AvailableCredits != 1 {
		t.Fatalf("payment required response = %#v", got)
	}
	assertOCRJobCount(t, db, 0)
	assertDirEmpty(t, fileDir)
}

func TestCreateOCRJobRejectsInsufficientEffectiveCreditsWithActiveJobs(t *testing.T) {
	router, db, fileDir, _ := testRouterWithOCRJobs(t)
	user := createTestUser(t, db, "ocr-job-active-balance@example.com")
	grantTestCredits(t, db, user.ID, 10)
	createStoredOCRJob(t, db, ocr.OCRJob{UserID: &user.ID, Status: ocr.OCRJobStatusQueued, PageCount: 5})
	createStoredOCRJob(t, db, ocr.OCRJob{UserID: &user.ID, Status: ocr.OCRJobStatusProcessing, PageCount: 3})

	req := multipartRequestForPath(t, "/api/ocr/jobs", map[string]string{"user_id": user.ID}, "invoice.pdf", threePagePDFBytes())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusPaymentRequired {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeOCRJobPaymentRequiredResponse(t, w)
	if got.Error != "insufficient credits" || got.RequiredCredits != 3 || got.AvailableCredits != 2 {
		t.Fatalf("payment required response = %#v", got)
	}
	assertOCRJobCount(t, db, 2)
	assertDirEmpty(t, fileDir)
}

func TestCreateOCRJobAllowsSufficientEffectiveCredits(t *testing.T) {
	router, db, fileDir, _ := testRouterWithOCRJobs(t)
	user := createTestUser(t, db, "ocr-job-enough-balance@example.com")
	grantTestCredits(t, db, user.ID, 10)
	createStoredOCRJob(t, db, ocr.OCRJob{UserID: &user.ID, Status: ocr.OCRJobStatusQueued, PageCount: 4})
	createStoredOCRJob(t, db, ocr.OCRJob{UserID: &user.ID, Status: ocr.OCRJobStatusProcessing, PageCount: 3})
	file := threePagePDFBytes()
	hash, err := computeDocumentHash(file, nil, false)
	if err != nil {
		t.Fatalf("compute document hash: %v", err)
	}

	req := multipartRequestForPath(t, "/api/ocr/jobs", map[string]string{"user_id": user.ID}, "invoice.pdf", file)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got OCRJobResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.PageCount != 3 {
		t.Fatalf("page_count = %d, want 3", got.PageCount)
	}
	assertPersistedQueuedOCRJob(t, db, got, hash, fileDir)
	assertOCRJobCount(t, db, 3)
}

func TestCreateOCRJobCountsPDFPagesWithRelaxedWhitespace(t *testing.T) {
	router, db, _, _ := testRouterWithOCRJobs(t)
	file := relaxedWhitespacePDFBytes()
	fields := creditedOCRJobFields(t, db, 5, nil)
	req := multipartRequestForPath(t, "/api/ocr/jobs", fields, "invoice.pdf", file)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got OCRJobResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.PageCount != 1 {
		t.Fatalf("page_count = %d, want 1", got.PageCount)
	}

	var job ocr.OCRJob
	if err := db.First(&job, "id = ?", got.ID).Error; err != nil {
		t.Fatalf("load OCR job: %v", err)
	}
	if job.PageCount != 1 {
		t.Fatalf("stored page_count = %d, want 1", job.PageCount)
	}
}

func TestCreateOCRJobRejectsPDFWithUnknownPageCount(t *testing.T) {
	router, db, fileDir, _ := testRouterWithOCRJobs(t)
	fields := creditedOCRJobFields(t, db, 5, nil)
	req := multipartRequestForPath(t, "/api/ocr/jobs", fields, "invoice.pdf", validPDFBytes())
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if got.Error != "failed to count document pages" {
		t.Fatalf("error = %q", got.Error)
	}
	var count int64
	if err := db.Model(&ocr.OCRJob{}).Count(&count).Error; err != nil {
		t.Fatalf("count OCR jobs: %v", err)
	}
	if count != 0 {
		t.Fatalf("OCR job count = %d, want 0", count)
	}
	assertDirEmpty(t, fileDir)
}

func TestCreateOCRJobNotifiesAfterQueueInsert(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := apiPostgresTx(t)
	storageDir := t.TempDir()
	var notified uuid.UUID

	router := NewRouter(&Handler{
		DB:               db,
		MaxUploadBytes:   20 << 20,
		StorageDir:       storageDir,
		InternalAPIToken: testInternalAPIToken,
		OCRJobNotifier: func(ctx context.Context, gotDB *gorm.DB, id uuid.UUID) error {
			if gotDB != db {
				t.Fatal("notifier received unexpected DB")
			}
			notified = id
			return nil
		},
	})

	fields := creditedOCRJobFields(t, db, 5, nil)
	req := multipartRequestForPath(t, "/api/ocr/jobs", fields, "scan.png", validPNGBytes())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got OCRJobResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if notified != got.ID {
		t.Fatalf("notified id = %s, want %s", notified, got.ID)
	}
}

func TestCreateOCRJobIgnoresNotifyFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := apiPostgresTx(t)
	router := NewRouter(&Handler{
		DB:               db,
		MaxUploadBytes:   20 << 20,
		StorageDir:       t.TempDir(),
		InternalAPIToken: testInternalAPIToken,
		OCRJobNotifier: func(context.Context, *gorm.DB, uuid.UUID) error {
			return errors.New("notify failed")
		},
	})

	fields := creditedOCRJobFields(t, db, 5, nil)
	req := multipartRequestForPath(t, "/api/ocr/jobs", fields, "scan.png", validPNGBytes())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestCreateOCRJobCleansUpStoredFileWhenDatabaseCreateFails(t *testing.T) {
	router, db, fileDir, _ := testRouterWithOCRJobs(t)
	createErr := errors.New("forced OCR job create failure")
	const callbackName = "test:fail_ocr_job_create"
	if err := db.Callback().Create().Before("gorm:create").Register(callbackName, func(tx *gorm.DB) {
		if tx.Statement != nil && tx.Statement.Schema != nil && tx.Statement.Schema.Table == "ocr_jobs" {
			tx.AddError(createErr)
		}
	}); err != nil {
		t.Fatalf("register OCR job create failure callback: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Callback().Create().Remove(callbackName); err != nil {
			t.Fatalf("remove OCR job create failure callback: %v", err)
		}
	})

	fields := creditedOCRJobFields(t, db, 5, nil)
	req := multipartRequestForPath(t, "/api/ocr/jobs", fields, "scan.png", validPNGBytes())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if got.Error != "failed to save OCR job" {
		t.Fatalf("error = %q", got.Error)
	}
	assertDirEmpty(t, fileDir)
}

func TestCreateOCRJobStoresInlineSchema(t *testing.T) {
	router, db, _, _ := testRouterWithOCRJobs(t)
	inlineSchema := `{"type":"object","properties":{"total":{"type":"number"}}}`
	fields := creditedOCRJobFields(t, db, 5, map[string]string{"schema": inlineSchema})
	req := multipartRequestForPath(t, "/api/ocr/jobs", fields, "scan.png", validPNGBytes())
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got OCRJobResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	var job ocr.OCRJob
	if err := db.First(&job, "id = ?", got.ID).Error; err != nil {
		t.Fatalf("load OCR job: %v", err)
	}
	if len(job.InlineSchemaJSON) == 0 {
		t.Fatal("stored inline schema is empty")
	}
	if job.SchemaID != nil {
		t.Fatalf("stored schema_id = %v, want nil", job.SchemaID)
	}
	assertJSONEqual(t, json.RawMessage(job.InlineSchemaJSON), inlineSchema)
}

func TestCreateOCRJobStoresSavedSchemaID(t *testing.T) {
	router, db, _, _ := testRouterWithOCRJobs(t)
	fields := creditedOCRJobFields(t, db, 5, nil)
	schema := ocr.ExtractionSchema{Name: "invoice", SchemaJSON: []byte(`{"type":"object"}`), Strict: true}
	if err := db.Create(&schema).Error; err != nil {
		t.Fatalf("create schema: %v", err)
	}
	fields["schema_id"] = schema.ID.String()
	req := multipartRequestForPath(t, "/api/ocr/jobs", fields, "scan.png", validPNGBytes())
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got OCRJobResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	var job ocr.OCRJob
	if err := db.First(&job, "id = ?", got.ID).Error; err != nil {
		t.Fatalf("load OCR job: %v", err)
	}
	if job.SchemaID == nil || *job.SchemaID != schema.ID {
		t.Fatalf("stored schema_id = %#v, want %s", job.SchemaID, schema.ID)
	}
}

func TestCreateOCRJobStoresUserID(t *testing.T) {
	router, db, _, _ := testRouterWithOCRJobs(t)
	user := createTestUser(t, db, "ocr-job-owner@example.com")
	grantTestCredits(t, db, user.ID, 1)
	req := multipartRequestForPath(t, "/api/ocr/jobs", map[string]string{"user_id": user.ID}, "scan.png", validPNGBytes())
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got OCRJobResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	var job ocr.OCRJob
	if err := db.First(&job, "id = ?", got.ID).Error; err != nil {
		t.Fatalf("load OCR job: %v", err)
	}
	if job.UserID == nil || *job.UserID != user.ID {
		t.Fatalf("stored user_id = %#v, want %s", job.UserID, user.ID)
	}
}

func decodeOCRJobPaymentRequiredResponse(t *testing.T, w *httptest.ResponseRecorder) PaymentRequiredResponse {
	t.Helper()
	var got PaymentRequiredResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode payment required response: %v body=%s", err, w.Body.String())
	}
	return got
}

func TestCreateOCRJobValidationErrors(t *testing.T) {
	cases := []struct {
		name       string
		fields     map[string]string
		filename   string
		content    []byte
		wantStatus int
		wantError  string
		setup      func(t *testing.T, db *gorm.DB) map[string]string
	}{
		{
			name:       "missing file",
			wantStatus: http.StatusBadRequest,
			wantError:  "file is required",
		},
		{
			name:       "unsupported mime",
			filename:   "notes.txt",
			content:    []byte("hello"),
			wantStatus: http.StatusBadRequest,
			wantError:  "unsupported file type",
		},
		{
			name:       "missing user",
			filename:   "scan.png",
			content:    validPNGBytes(),
			wantStatus: http.StatusBadRequest,
			wantError:  "user_id is required",
		},
		{
			name:       "invalid user",
			fields:     map[string]string{"user_id": "not-a-uuid"},
			filename:   "scan.png",
			content:    validPNGBytes(),
			wantStatus: http.StatusBadRequest,
			wantError:  "invalid user_id",
		},
		{
			name:       "unknown user",
			fields:     map[string]string{"user_id": uuid.NewString()},
			filename:   "scan.png",
			content:    validPNGBytes(),
			wantStatus: http.StatusBadRequest,
			wantError:  "invalid user_id",
		},
		{
			name:       "both schema sources",
			filename:   "scan.png",
			content:    validPNGBytes(),
			wantStatus: http.StatusBadRequest,
			wantError:  "provide either schema or schema_id, not both",
			setup: func(t *testing.T, db *gorm.DB) map[string]string {
				t.Helper()
				schema := ocr.ExtractionSchema{Name: "invoice", SchemaJSON: []byte(`{"type":"object"}`), Strict: true}
				if err := db.Create(&schema).Error; err != nil {
					t.Fatalf("create schema: %v", err)
				}
				return creditedOCRJobFields(t, db, 5, map[string]string{"schema": `{"type":"object"}`, "schema_id": schema.ID.String()})
			},
		},
		{
			name:       "invalid schema id",
			filename:   "scan.png",
			content:    validPNGBytes(),
			wantStatus: http.StatusBadRequest,
			wantError:  "invalid schema_id",
			setup: func(t *testing.T, db *gorm.DB) map[string]string {
				t.Helper()
				return creditedOCRJobFields(t, db, 5, map[string]string{"schema_id": "not-a-uuid"})
			},
		},
		{
			name:       "missing schema id",
			filename:   "scan.png",
			content:    validPNGBytes(),
			wantStatus: http.StatusNotFound,
			wantError:  "schema not found",
			setup: func(t *testing.T, db *gorm.DB) map[string]string {
				t.Helper()
				return creditedOCRJobFields(t, db, 5, map[string]string{"schema_id": uuid.NewString()})
			},
		},
		{
			name:       "invalid inline schema",
			filename:   "scan.png",
			content:    validPNGBytes(),
			wantStatus: http.StatusBadRequest,
			wantError:  "schema must be a JSON object",
			setup: func(t *testing.T, db *gorm.DB) map[string]string {
				t.Helper()
				return creditedOCRJobFields(t, db, 5, map[string]string{"schema": `["not-object"]`})
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			router, db, _, _ := testRouterWithOCRJobs(t)
			fields := tc.fields
			if tc.setup != nil {
				fields = tc.setup(t, db)
			}
			req := multipartRequestForPath(t, "/api/ocr/jobs", fields, tc.filename, tc.content)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tc.wantStatus {
				t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
			}
			var got ErrorResponse
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Fatalf("decode error: %v", err)
			}
			if got.Error != tc.wantError {
				t.Fatalf("error = %q, want %q", got.Error, tc.wantError)
			}
			var count int64
			if err := db.Model(&ocr.OCRJob{}).Count(&count).Error; err != nil {
				t.Fatalf("count OCR jobs: %v", err)
			}
			if count != 0 {
				t.Fatalf("OCR job count = %d, want 0", count)
			}
		})
	}
}

func TestCreateOCRJobRejectsLongFilenameBeforeWritingFile(t *testing.T) {
	router, db, fileDir, _ := testRouterWithOCRJobs(t)
	req := multipartRequestForPath(t, "/api/ocr/jobs", nil, strings.Repeat("a", 256)+".png", validPNGBytes())
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if got.Error != "filename must be at most 255 characters" {
		t.Fatalf("error = %q", got.Error)
	}
	var count int64
	if err := db.Model(&ocr.OCRJob{}).Count(&count).Error; err != nil {
		t.Fatalf("count OCR jobs: %v", err)
	}
	if count != 0 {
		t.Fatalf("OCR job count = %d, want 0", count)
	}
	assertDirEmpty(t, fileDir)
}

func TestListOCRJobsScopesByUserIDAndSystemDefault(t *testing.T) {
	router, db, _, _ := testRouterWithOCRJobs(t)
	user := createTestUser(t, db, "job-list-owner@example.com")
	other := createTestUser(t, db, "job-list-other@example.com")
	base := time.Date(2026, 5, 28, 10, 0, 0, 0, time.UTC)
	systemJob := setStoredOCRJobCreatedAt(t, db, createStoredOCRJob(t, db, ocr.OCRJob{
		DocumentHash: "system-list-hash",
	}), base)
	userJob := setStoredOCRJobCreatedAt(t, db, createStoredOCRJob(t, db, ocr.OCRJob{
		UserID:       &user.ID,
		DocumentHash: "user-list-hash",
	}), base.Add(time.Minute))
	setStoredOCRJobCreatedAt(t, db, createStoredOCRJob(t, db, ocr.OCRJob{
		UserID:       &other.ID,
		DocumentHash: "other-list-hash",
	}), base.Add(2*time.Minute))

	req := newTestRequest(http.MethodGet, "/api/ocr/jobs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("system list status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeOCRJobListResponse(t, w)
	assertOCRJobListIDs(t, got.Jobs, systemJob.ID)
	if got.Jobs[0].CreatedAt.IsZero() {
		t.Fatal("created_at is zero")
	}
	if got.NextCursor != nil {
		t.Fatalf("next_cursor = %q, want nil", *got.NextCursor)
	}

	req = newTestRequest(http.MethodGet, "/api/ocr/jobs?user_id="+user.ID, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("user list status = %d body=%s", w.Code, w.Body.String())
	}
	got = decodeOCRJobListResponse(t, w)
	assertOCRJobListIDs(t, got.Jobs, userJob.ID)
}

func TestListOCRJobsIncludesFileAndSchemaMetadata(t *testing.T) {
	router, db, _, _ := testRouterWithOCRJobs(t)
	user := createTestUser(t, db, "job-list-metadata@example.com")
	schema := createStoredJobSchema(t, db, user.ID, "Invoice")
	base := time.Date(2026, 5, 28, 10, 0, 0, 0, time.UTC)
	job := setStoredOCRJobCreatedAt(t, db, createStoredOCRJob(t, db, ocr.OCRJob{
		UserID:           &user.ID,
		OriginalFilename: "invoice.pdf",
		MimeType:         "application/pdf",
		FileSize:         2048,
		PageCount:        4,
		DocumentHash:     "job-metadata-hash",
		SchemaID:         &schema.ID,
		Status:           ocr.OCRJobStatusProcessing,
	}), base)

	req := newTestRequest(http.MethodGet, "/api/ocr/jobs?user_id="+user.ID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeOCRJobListResponse(t, w)
	if len(got.Jobs) != 1 {
		t.Fatalf("jobs = %#v, want one job", got.Jobs)
	}
	row := got.Jobs[0]
	if row.ID != job.ID || row.CreatedAt.IsZero() {
		t.Fatalf("job identity = %#v, want id %s and non-zero created_at", row, job.ID)
	}
	if row.OriginalFilename != "invoice.pdf" || row.MimeType != "application/pdf" {
		t.Fatalf("file metadata = %q/%q", row.OriginalFilename, row.MimeType)
	}
	if row.FileSize != 2048 || row.PageCount != 4 {
		t.Fatalf("size/page metadata = %d/%d", row.FileSize, row.PageCount)
	}
	if row.SchemaID == nil || *row.SchemaID != schema.ID {
		t.Fatalf("schema_id = %#v, want %s", row.SchemaID, schema.ID)
	}
	if row.SchemaName == nil || *row.SchemaName != "Invoice" {
		t.Fatalf("schema_name = %#v, want Invoice", row.SchemaName)
	}
	if row.HasInlineSchema {
		t.Fatal("has_inline_schema = true, want false for saved schema job")
	}
}

func TestListOCRJobsFiltersByStatusAndCreatedRange(t *testing.T) {
	router, db, _, _ := testRouterWithOCRJobs(t)
	user := createTestUser(t, db, "job-list-filter@example.com")
	base := time.Date(2026, 5, 28, 12, 0, 0, 0, time.UTC)
	setStoredOCRJobCreatedAt(t, db, createStoredOCRJob(t, db, ocr.OCRJob{
		UserID:       &user.ID,
		DocumentHash: "filter-too-old",
		Status:       ocr.OCRJobStatusCompleted,
	}), base.Add(-time.Hour))
	first := setStoredOCRJobCreatedAt(t, db, createStoredOCRJob(t, db, ocr.OCRJob{
		UserID:       &user.ID,
		DocumentHash: "filter-first",
		Status:       ocr.OCRJobStatusCompleted,
	}), base)
	second := setStoredOCRJobCreatedAt(t, db, createStoredOCRJob(t, db, ocr.OCRJob{
		UserID:       &user.ID,
		DocumentHash: "filter-second",
		Status:       ocr.OCRJobStatusCompleted,
	}), base.Add(time.Hour))
	setStoredOCRJobCreatedAt(t, db, createStoredOCRJob(t, db, ocr.OCRJob{
		UserID:       &user.ID,
		DocumentHash: "filter-wrong-status",
		Status:       ocr.OCRJobStatusQueued,
	}), base.Add(30*time.Minute))
	setStoredOCRJobCreatedAt(t, db, createStoredOCRJob(t, db, ocr.OCRJob{
		UserID:       &user.ID,
		DocumentHash: "filter-too-new",
		Status:       ocr.OCRJobStatusCompleted,
	}), base.Add(2*time.Hour))

	query := url.Values{}
	query.Set("user_id", user.ID)
	query.Set("status", string(ocr.OCRJobStatusCompleted))
	query.Set("created_from", base.Format(time.RFC3339Nano))
	query.Set("created_to", base.Add(time.Hour).Format(time.RFC3339Nano))
	query.Set("sort", "asc")
	req := newTestRequest(http.MethodGet, "/api/ocr/jobs?"+query.Encode(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeOCRJobListResponse(t, w)
	assertOCRJobListIDs(t, got.Jobs, first.ID, second.ID)
}

func TestListOCRJobsSortsByCreatedAtAndID(t *testing.T) {
	router, db, _, _ := testRouterWithOCRJobs(t)
	user := createTestUser(t, db, "job-list-sort@example.com")
	base := time.Date(2026, 5, 28, 14, 0, 0, 0, time.UTC)
	old := setStoredOCRJobCreatedAt(t, db, createStoredOCRJob(t, db, ocr.OCRJob{
		ID:           uuid.MustParse("00000000-0000-0000-0000-000000000003"),
		UserID:       &user.ID,
		DocumentHash: "sort-old",
	}), base.Add(-time.Hour))
	lowID := setStoredOCRJobCreatedAt(t, db, createStoredOCRJob(t, db, ocr.OCRJob{
		ID:           uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		UserID:       &user.ID,
		DocumentHash: "sort-low",
	}), base)
	highID := setStoredOCRJobCreatedAt(t, db, createStoredOCRJob(t, db, ocr.OCRJob{
		ID:           uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		UserID:       &user.ID,
		DocumentHash: "sort-high",
	}), base)

	req := newTestRequest(http.MethodGet, "/api/ocr/jobs?user_id="+user.ID+"&sort=asc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("asc status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeOCRJobListResponse(t, w)
	assertOCRJobListIDs(t, got.Jobs, old.ID, lowID.ID, highID.ID)

	req = newTestRequest(http.MethodGet, "/api/ocr/jobs?user_id="+user.ID+"&sort=desc", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("desc status = %d body=%s", w.Code, w.Body.String())
	}
	got = decodeOCRJobListResponse(t, w)
	assertOCRJobListIDs(t, got.Jobs, highID.ID, lowID.ID, old.ID)
}

func TestListOCRJobsPaginatesWithCursor(t *testing.T) {
	router, db, _, _ := testRouterWithOCRJobs(t)
	user := createTestUser(t, db, "job-list-cursor@example.com")
	base := time.Date(2026, 5, 28, 16, 0, 0, 0, time.UTC)
	oldest := setStoredOCRJobCreatedAt(t, db, createStoredOCRJob(t, db, ocr.OCRJob{
		UserID:       &user.ID,
		DocumentHash: "cursor-oldest",
	}), base)
	middle := setStoredOCRJobCreatedAt(t, db, createStoredOCRJob(t, db, ocr.OCRJob{
		UserID:       &user.ID,
		DocumentHash: "cursor-middle",
	}), base.Add(time.Minute))
	newest := setStoredOCRJobCreatedAt(t, db, createStoredOCRJob(t, db, ocr.OCRJob{
		UserID:       &user.ID,
		DocumentHash: "cursor-newest",
	}), base.Add(2*time.Minute))

	req := newTestRequest(http.MethodGet, "/api/ocr/jobs?user_id="+user.ID+"&size=2", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("first page status = %d body=%s", w.Code, w.Body.String())
	}
	firstPage := decodeOCRJobListResponse(t, w)
	assertOCRJobListIDs(t, firstPage.Jobs, newest.ID, middle.ID)
	if firstPage.NextCursor == nil || *firstPage.NextCursor == "" {
		t.Fatalf("next_cursor = %#v, want non-empty", firstPage.NextCursor)
	}

	query := url.Values{}
	query.Set("user_id", user.ID)
	query.Set("size", "2")
	query.Set("cursor", *firstPage.NextCursor)
	req = newTestRequest(http.MethodGet, "/api/ocr/jobs?"+query.Encode(), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("second page status = %d body=%s", w.Code, w.Body.String())
	}
	secondPage := decodeOCRJobListResponse(t, w)
	assertOCRJobListIDs(t, secondPage.Jobs, oldest.ID)
	if secondPage.NextCursor != nil {
		t.Fatalf("next_cursor = %q, want nil", *secondPage.NextCursor)
	}
}

func TestListOCRJobsRejectsInvalidQueries(t *testing.T) {
	router, _, _, _ := testRouterWithOCRJobs(t)
	cursor, err := encodeOCRJobListCursor(ocr.OCRJob{
		ID:        uuid.New(),
		CreatedAt: time.Date(2026, 5, 28, 18, 0, 0, 0, time.UTC),
	}, "asc")
	if err != nil {
		t.Fatalf("encode cursor: %v", err)
	}

	cases := []struct {
		name      string
		query     string
		wantError string
	}{
		{name: "invalid user", query: "user_id=not-a-uuid", wantError: "invalid user_id"},
		{name: "invalid status", query: "status=done", wantError: "invalid status"},
		{name: "invalid created from", query: "created_from=not-a-date", wantError: "invalid created_from"},
		{name: "invalid created to", query: "created_to=not-a-date", wantError: "invalid created_to"},
		{
			name: "inverted range",
			query: url.Values{
				"created_from": []string{time.Date(2026, 5, 29, 0, 0, 0, 0, time.UTC).Format(time.RFC3339Nano)},
				"created_to":   []string{time.Date(2026, 5, 28, 0, 0, 0, 0, time.UTC).Format(time.RFC3339Nano)},
			}.Encode(),
			wantError: "created_from must be before or equal to created_to",
		},
		{name: "invalid sort", query: "sort=newest", wantError: "sort must be asc or desc"},
		{name: "invalid size text", query: "size=large", wantError: "size must be between 1 and 100"},
		{name: "invalid size zero", query: "size=0", wantError: "size must be between 1 and 100"},
		{name: "invalid size too large", query: "size=101", wantError: "size must be between 1 and 100"},
		{name: "invalid cursor", query: "cursor=not-base64", wantError: "invalid cursor"},
		{
			name:      "cursor sort mismatch",
			query:     "sort=desc&cursor=" + url.QueryEscape(cursor),
			wantError: "cursor sort does not match sort",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := newTestRequest(http.MethodGet, "/api/ocr/jobs?"+tc.query, nil)
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

func TestGetOCRJobReturnsQueuedJob(t *testing.T) {
	router, db, _, _ := testRouterWithOCRJobs(t)
	job := ocr.OCRJob{
		OriginalFilename: "queued.png",
		MimeType:         "image/png",
		FileSize:         int64(len(validPNGBytes())),
		PageCount:        1,
		DocumentHash:     "queued-hash",
		FilePath:         "/tmp/queued.png",
		Status:           ocr.OCRJobStatusQueued,
	}
	if err := db.Create(&job).Error; err != nil {
		t.Fatalf("create OCR job: %v", err)
	}

	req := newTestRequest(http.MethodGet, "/api/ocr/jobs/"+job.ID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got OCRJobResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ID != job.ID || got.Status != string(ocr.OCRJobStatusQueued) {
		t.Fatalf("unexpected response: %#v", got)
	}
	if got.OriginalFilename != "queued.png" || got.MimeType != "image/png" {
		t.Fatalf("file metadata = %q/%q", got.OriginalFilename, got.MimeType)
	}
	if got.DocumentID != nil {
		t.Fatalf("document_id = %#v, want nil", got.DocumentID)
	}
	raw := assertOCRJobResponseKeys(t, w.Body.Bytes())
	if raw["document_id"] != nil {
		t.Fatalf("document_id = %#v, want null", raw["document_id"])
	}
}

func TestGetPublicOCRJobReturnsOwnedJobStatuses(t *testing.T) {
	router, db, _, _ := testRouterWithOCRJobs(t)
	user := createTestUser(t, db, "public-ocr-job-status-owner@example.com")
	secret := createTestPublicAPIKey(t, db, user.ID, nil)
	rawResponse := `{"pages":[{"index":0,"markdown":"# Done"}],"document_annotation":{"total":10},"model":"mistral-ocr-test"}`
	doc := ocr.OCRDocument{
		UserID:           &user.ID,
		OriginalFilename: "stored-done.pdf",
		MimeType:         "application/pdf",
		FileSize:         12345,
		DocumentHash:     "public-status-doc-hash",
		Markdown:         "# Done",
		RawResponseJSON:  datatypes.JSON([]byte(rawResponse)),
	}
	if err := db.Create(&doc).Error; err != nil {
		t.Fatalf("create OCR document: %v", err)
	}
	jobs := []ocr.OCRJob{
		createStoredOCRJob(t, db, ocr.OCRJob{
			UserID:           &user.ID,
			OriginalFilename: "queued.png",
			DocumentHash:     "public-status-queued-hash",
			Status:           ocr.OCRJobStatusQueued,
		}),
		createStoredOCRJob(t, db, ocr.OCRJob{
			UserID:           &user.ID,
			OriginalFilename: "completed.png",
			DocumentHash:     "public-status-completed-hash",
			DocumentID:       &doc.ID,
			Status:           ocr.OCRJobStatusCompleted,
		}),
		createStoredOCRJob(t, db, ocr.OCRJob{
			UserID:           &user.ID,
			OriginalFilename: "failed.png",
			DocumentHash:     "public-status-failed-hash",
			Status:           ocr.OCRJobStatusFailed,
			ErrorMessage:     "mistral OCR failed with status 503",
		}),
	}

	for _, job := range jobs {
		t.Run(string(job.Status), func(t *testing.T) {
			req := newTestRequest(http.MethodGet, "/v1/ocr/jobs/"+job.ID.String(), nil)
			authorizePublicAPIRequest(req, secret)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
			}
			var got PublicOCRJobResponse
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Fatalf("decode: %v", err)
			}
			if got.ID != job.ID || got.Status != string(job.Status) {
				t.Fatalf("unexpected response: %#v", got)
			}
			var raw map[string]any
			if err := json.Unmarshal(w.Body.Bytes(), &raw); err != nil {
				t.Fatalf("decode raw response: %v", err)
			}
			wantTopLevelKeys := []string{"id", "created_at", "status", "original_filename", "has_inline_schema"}
			if job.Status == ocr.OCRJobStatusCompleted {
				wantTopLevelKeys = append(wantTopLevelKeys, "document")
			}
			assertRawObjectKeys(t, raw, wantTopLevelKeys...)
			if raw["id"] != job.ID.String() {
				t.Fatalf("id = %#v, want %s", raw["id"], job.ID)
			}
			if raw["status"] != string(job.Status) {
				t.Fatalf("status = %#v, want %s", raw["status"], job.Status)
			}
			if raw["original_filename"] != job.OriginalFilename {
				t.Fatalf("original_filename = %#v, want %s", raw["original_filename"], job.OriginalFilename)
			}
			if raw["has_inline_schema"] != false {
				t.Fatalf("has_inline_schema = %#v, want false", raw["has_inline_schema"])
			}
			if job.Status == ocr.OCRJobStatusCompleted {
				if got.Document == nil {
					t.Fatalf("document = nil, body=%s", w.Body.String())
				}
				document, ok := raw["document"].(map[string]any)
				if !ok {
					t.Fatalf("document = %#v, want object", raw["document"])
				}
				assertRawObjectKeys(t, document, "document_id", "file_size", "page_count", "pages", "document_annotation")
				if document["document_id"] != doc.ID.String() {
					t.Fatalf("document.document_id = %#v, want %s", document["document_id"], doc.ID)
				}
				if document["file_size"] != float64(doc.FileSize) {
					t.Fatalf("document.file_size = %#v, want %d", document["file_size"], doc.FileSize)
				}
				if document["page_count"] != float64(1) {
					t.Fatalf("document.page_count = %#v, want 1", document["page_count"])
				}
				if got.Document.PageCount != 1 {
					t.Fatalf("document.page_count = %d, want 1", got.Document.PageCount)
				}
				assertRawJSONValue(t, document["pages"], `[{"index":0,"markdown":"# Done"}]`)
				assertRawJSONValue(t, document["document_annotation"], `{"total":10}`)
			} else if got.Document != nil {
				t.Fatalf("document = %#v, want nil for %s", got.Document, job.Status)
			}
			if job.Status != ocr.OCRJobStatusCompleted {
				if _, ok := raw["document"]; ok {
					t.Fatalf("document present for %s: %#v", job.Status, raw["document"])
				}
			}
		})
	}
}

func TestGetPublicOCRJobReturnsNullDocumentWhenLinkedDocumentUnavailable(t *testing.T) {
	router, db, _, _ := testRouterWithOCRJobs(t)
	user := createTestUser(t, db, "public-ocr-job-missing-doc-owner@example.com")
	secret := createTestPublicAPIKey(t, db, user.ID, nil)
	doc := ocr.OCRDocument{
		UserID:           &user.ID,
		OriginalFilename: "deleted.pdf",
		MimeType:         "application/pdf",
		FileSize:         99,
		DocumentHash:     "public-status-deleted-doc-hash",
		Markdown:         "# Deleted",
		RawResponseJSON:  datatypes.JSON([]byte(`{"pages":[]}`)),
	}
	if err := db.Create(&doc).Error; err != nil {
		t.Fatalf("create OCR document: %v", err)
	}
	job := createStoredOCRJob(t, db, ocr.OCRJob{
		UserID:       &user.ID,
		DocumentHash: "public-status-missing-doc-job-hash",
		DocumentID:   &doc.ID,
		Status:       ocr.OCRJobStatusCompleted,
	})
	if err := db.Delete(&doc).Error; err != nil {
		t.Fatalf("soft-delete OCR document: %v", err)
	}
	req := newTestRequest(http.MethodGet, "/v1/ocr/jobs/"+job.ID.String(), nil)
	authorizePublicAPIRequest(req, secret)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got PublicOCRJobResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Document != nil {
		t.Fatalf("document = %#v, want nil", got.Document)
	}
	var raw map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &raw); err != nil {
		t.Fatalf("decode raw response: %v", err)
	}
	assertRawObjectKeys(t, raw, "id", "created_at", "status", "original_filename", "has_inline_schema")
	if _, ok := raw["document"]; ok {
		t.Fatalf("document present for missing linked document: %#v", raw["document"])
	}
}

func TestGetPublicOCRJobScopesToAPIKeyUser(t *testing.T) {
	router, db, _, _ := testRouterWithOCRJobs(t)
	user := createTestUser(t, db, "public-ocr-job-scope-owner@example.com")
	other := createTestUser(t, db, "public-ocr-job-scope-other@example.com")
	secret := createTestPublicAPIKey(t, db, user.ID, nil)
	job := createStoredOCRJob(t, db, ocr.OCRJob{
		UserID:       &other.ID,
		DocumentHash: "public-scope-other-hash",
		Status:       ocr.OCRJobStatusQueued,
	})
	req := newTestRequest(http.MethodGet, "/v1/ocr/jobs/"+job.ID.String(), nil)
	authorizePublicAPIRequest(req, secret)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestGetPublicOCRJobRejectsInvalidID(t *testing.T) {
	router, db, _, _ := testRouterWithOCRJobs(t)
	user := createTestUser(t, db, "public-ocr-job-invalid-id@example.com")
	secret := createTestPublicAPIKey(t, db, user.ID, nil)
	req := newTestRequest(http.MethodGet, "/v1/ocr/jobs/not-a-uuid", nil)
	authorizePublicAPIRequest(req, secret)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestGetPublicOCRJobRejectsUserIDQuery(t *testing.T) {
	router, db, _, _ := testRouterWithOCRJobs(t)
	user := createTestUser(t, db, "public-ocr-job-query-status@example.com")
	secret := createTestPublicAPIKey(t, db, user.ID, nil)
	job := createStoredOCRJob(t, db, ocr.OCRJob{
		UserID:       &user.ID,
		DocumentHash: "public-query-status-hash",
		Status:       ocr.OCRJobStatusQueued,
	})
	req := newTestRequest(http.MethodGet, "/v1/ocr/jobs/"+job.ID.String()+"?user_id="+url.QueryEscape(user.ID), nil)
	authorizePublicAPIRequest(req, secret)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestGetOCRJobReturnsSchemaMetadataAndInlineSchema(t *testing.T) {
	router, db, _, _ := testRouterWithOCRJobs(t)
	job := createStoredOCRJob(t, db, ocr.OCRJob{
		OriginalFilename: "inline.png",
		DocumentHash:     "inline-schema-job-hash",
		Status:           ocr.OCRJobStatusQueued,
		InlineSchemaJSON: datatypes.JSON([]byte(`{"type":"object","properties":{"total":{"type":"number"}}}`)),
	})

	req := newTestRequest(http.MethodGet, "/api/ocr/jobs/"+job.ID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got OCRJobResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !got.HasInlineSchema {
		t.Fatal("has_inline_schema = false, want true")
	}
	if got.InlineSchema == nil {
		t.Fatalf("inline_schema = nil, body=%s", w.Body.String())
	}
	if got.SchemaName != nil || got.SchemaID != nil {
		t.Fatalf("saved schema metadata = %#v/%#v, want nil for inline schema job", got.SchemaID, got.SchemaName)
	}
}

func TestDeleteOCRJobsSoftDeletesOnlyScopedJobsAndKeepsDocuments(t *testing.T) {
	router, db, _, _ := testRouterWithOCRJobs(t)
	owner := createTestUser(t, db, "job-delete-owner@example.com")
	other := createTestUser(t, db, "job-delete-other@example.com")
	doc := ocr.OCRDocument{
		UserID:           &owner.ID,
		OriginalFilename: "completed.png",
		MimeType:         "image/png",
		FileSize:         int64(len(validPNGBytes())),
		DocumentHash:     "completed-doc-hash",
		Markdown:         "# Completed",
		RawResponseJSON:  datatypes.JSON([]byte(`{"pages":[]}`)),
	}
	if err := db.Create(&doc).Error; err != nil {
		t.Fatalf("create OCR document: %v", err)
	}
	owned := createStoredOCRJob(t, db, ocr.OCRJob{
		UserID:       &owner.ID,
		DocumentHash: "delete-owned-hash",
		DocumentID:   &doc.ID,
		Status:       ocr.OCRJobStatusCompleted,
	})
	wrongOwner := createStoredOCRJob(t, db, ocr.OCRJob{
		UserID:       &other.ID,
		DocumentHash: "delete-other-hash",
	})
	body := `{"ids":["` + owned.ID.String() + `","` + wrongOwner.ID.String() + `","` + uuid.NewString() + `"]}`

	req := newTestRequest(http.MethodDelete, "/api/ocr/jobs?user_id="+url.QueryEscape(owner.ID), strings.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeDeleteOCRJobsResponse(t, w)
	if got.DeletedCount != 1 || len(got.DeletedIDs) != 1 || got.DeletedIDs[0] != owned.ID {
		t.Fatalf("delete response = %#v, want only owned job", got)
	}
	var deleted ocr.OCRJob
	if err := db.Unscoped().First(&deleted, "id = ?", owned.ID).Error; err != nil {
		t.Fatalf("load soft-deleted job: %v", err)
	}
	if !deleted.DeletedAt.Valid {
		t.Fatal("deleted_at is not set on deleted job")
	}
	var keptDoc ocr.OCRDocument
	if err := db.First(&keptDoc, "id = ?", doc.ID).Error; err != nil {
		t.Fatalf("completed document was deleted with job: %v", err)
	}
	req = newTestRequest(http.MethodGet, "/api/ocr/jobs/"+owned.ID.String()+"?user_id="+url.QueryEscape(owner.ID), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("deleted job detail status = %d body=%s", w.Code, w.Body.String())
	}
	req = newTestRequest(http.MethodGet, "/api/ocr/jobs?user_id="+url.QueryEscape(owner.ID), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("list after delete status = %d body=%s", w.Code, w.Body.String())
	}
	list := decodeOCRJobListResponse(t, w)
	if len(list.Jobs) != 0 {
		t.Fatalf("list after delete = %#v, want empty", list.Jobs)
	}
}

func TestDeleteOCRJobsRejectsInvalidBulkBodies(t *testing.T) {
	router, _, _, _ := testRouterWithOCRJobs(t)

	for _, body := range []string{
		``,
		`{"ids":"job-1"}`,
		`{"ids":["not-a-uuid"]}`,
	} {
		req := newTestRequest(http.MethodDelete, "/api/ocr/jobs", strings.NewReader(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("body %q status = %d body=%s", body, w.Code, w.Body.String())
		}
	}
}

func TestGetOCRJobReturnsCompletedDocumentID(t *testing.T) {
	router, db, _, _ := testRouterWithOCRJobs(t)
	doc := ocr.OCRDocument{
		OriginalFilename: "done.png",
		MimeType:         "image/png",
		FileSize:         int64(len(validPNGBytes())),
		DocumentHash:     "done-hash",
		Markdown:         "# Done",
		RawResponseJSON:  datatypes.JSON([]byte(`{"pages":[]}`)),
	}
	if err := db.Create(&doc).Error; err != nil {
		t.Fatalf("create OCR document: %v", err)
	}
	job := ocr.OCRJob{
		OriginalFilename: "done.png",
		MimeType:         "image/png",
		FileSize:         int64(len(validPNGBytes())),
		PageCount:        1,
		DocumentHash:     "done-hash",
		FilePath:         "/tmp/done.png",
		DocumentID:       &doc.ID,
		Status:           ocr.OCRJobStatusCompleted,
	}
	if err := db.Create(&job).Error; err != nil {
		t.Fatalf("create OCR job: %v", err)
	}

	req := newTestRequest(http.MethodGet, "/api/ocr/jobs/"+job.ID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got OCRJobResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.DocumentID == nil || *got.DocumentID != doc.ID {
		t.Fatalf("document_id = %#v, want %s", got.DocumentID, doc.ID)
	}
	if got.PageCount != 1 {
		t.Fatalf("page_count = %d, want 1", got.PageCount)
	}
	raw := assertOCRJobResponseKeys(t, w.Body.Bytes())
	if _, ok := raw["document"]; ok {
		t.Fatalf("private OCR job response includes public document payload: %s", w.Body.String())
	}
}

func TestGetOCRJobReturnsFailedErrorMessage(t *testing.T) {
	router, db, _, _ := testRouterWithOCRJobs(t)
	job := ocr.OCRJob{
		OriginalFilename: "failed.png",
		MimeType:         "image/png",
		FileSize:         int64(len(validPNGBytes())),
		PageCount:        1,
		DocumentHash:     "failed-hash",
		FilePath:         "/tmp/failed.png",
		Status:           ocr.OCRJobStatusFailed,
		ErrorMessage:     "mistral OCR failed with status 503",
	}
	if err := db.Create(&job).Error; err != nil {
		t.Fatalf("create OCR job: %v", err)
	}

	req := newTestRequest(http.MethodGet, "/api/ocr/jobs/"+job.ID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var raw map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &raw); err != nil {
		t.Fatalf("decode raw response: %v", err)
	}
	if raw["error_message"] != job.ErrorMessage {
		t.Fatalf("error_message = %#v, want %q; body=%s", raw["error_message"], job.ErrorMessage, w.Body.String())
	}
}

func TestGetOCRJobScopesByUserIDQuery(t *testing.T) {
	router, db, _, _ := testRouterWithOCRJobs(t)
	user := createTestUser(t, db, "job-scope@example.com")
	other := createTestUser(t, db, "job-scope-other@example.com")
	job := createStoredOCRJob(t, db, ocr.OCRJob{
		UserID:       &user.ID,
		DocumentHash: "scope-hash",
		Status:       ocr.OCRJobStatusQueued,
	})

	req := newTestRequest(http.MethodGet, "/api/ocr/jobs/"+job.ID.String()+"?user_id="+other.ID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestGetOCRJobRejectsInvalidUserIDQuery(t *testing.T) {
	router, db, _, _ := testRouterWithOCRJobs(t)
	job := createStoredOCRJob(t, db, ocr.OCRJob{DocumentHash: "invalid-user-query-hash"})

	req := newTestRequest(http.MethodGet, "/api/ocr/jobs/"+job.ID.String()+"?user_id=not-a-uuid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestGetOCRJobRejectsInvalidID(t *testing.T) {
	router, _, _, _ := testRouterWithOCRJobs(t)
	req := newTestRequest(http.MethodGet, "/api/ocr/jobs/not-a-uuid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Error != "invalid OCR job id" {
		t.Fatalf("error = %q", got.Error)
	}
}

func TestGetOCRJobReturnsNotFound(t *testing.T) {
	router, _, _, _ := testRouterWithOCRJobs(t)
	req := newTestRequest(http.MethodGet, "/api/ocr/jobs/"+uuid.NewString(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var got ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Error != "OCR job not found" {
		t.Fatalf("error = %q", got.Error)
	}
}
