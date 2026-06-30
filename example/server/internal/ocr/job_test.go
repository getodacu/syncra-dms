package ocr

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"ai.ro/syncra/internal/auth"
	"ai.ro/syncra/internal/testsupport"
)

var ocrJobTestGroup *testsupport.PostgresGroup

func TestOCRJobs(t *testing.T) {
	ocrJobTestGroup = testsupport.OpenPostgresGroup(t, &auth.User{}, &ExtractionSchema{}, &OCRJob{})
	defer func() { ocrJobTestGroup = nil }()

	for _, tt := range []struct {
		name string
		fn   func(*testing.T)
	}{
		{name: "AutoMigrateAndPersist", fn: testOCRJobAutoMigrateAndPersist},
		{name: "RejectsInvalidStatus", fn: testOCRJobRejectsInvalidStatus},
		{name: "OpenPostgresTxAppliesOCRJobStatusDBCheck", fn: testOpenPostgresTxAppliesOCRJobStatusDBCheck},
		{name: "PartialStatusUpdatesValidateIncomingStatus", fn: testOCRJobPartialStatusUpdatesValidateIncomingStatus},
		{name: "JSONRedactsFilePath", fn: testOCRJobJSONRedactsFilePath},
	} {
		t.Run(tt.name, tt.fn)
	}
}

func ocrJobTx(t *testing.T) *gorm.DB {
	t.Helper()
	if ocrJobTestGroup != nil {
		return ocrJobTestGroup.Tx(t)
	}
	return testsupport.OpenPostgresTx(t, &auth.User{}, &ExtractionSchema{}, &OCRJob{})
}

func testOCRJobAutoMigrateAndPersist(t *testing.T) {
	db := ocrJobTx(t)

	user := auth.User{Name: "Queue Owner", Email: "queue-owner@example.com"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}

	schema := ExtractionSchema{
		UserID:     &user.ID,
		Name:       "invoice",
		SchemaJSON: datatypes.JSON([]byte(`{"type":"object"}`)),
		Strict:     true,
	}
	if err := db.Create(&schema).Error; err != nil {
		t.Fatalf("create schema: %v", err)
	}

	doc := OCRDocument{
		UserID:           &user.ID,
		OriginalFilename: "invoice.pdf",
		MimeType:         "application/pdf",
		FileSize:         42,
		DocumentHash:     "result-hash",
		Markdown:         "# Invoice",
		RawResponseJSON:  datatypes.JSON([]byte(`{"pages":[]}`)),
	}
	if err := db.Create(&doc).Error; err != nil {
		t.Fatalf("create OCR document: %v", err)
	}

	inlineSchema := datatypes.JSON([]byte(`{"type":"object","properties":{"total":{"type":"number"}}}`))
	job := OCRJob{
		UserID:           &user.ID,
		OriginalFilename: "invoice.pdf",
		MimeType:         "application/pdf",
		FileSize:         42,
		PageCount:        3,
		DocumentHash:     "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		FilePath:         "/tmp/syncra-test/invoice.pdf",
		SchemaID:         &schema.ID,
		InlineSchemaJSON: inlineSchema,
		DocumentID:       &doc.ID,
		Status:           OCRJobStatusQueued,
	}
	if err := db.Create(&job).Error; err != nil {
		t.Fatalf("create OCR job: %v", err)
	}
	if job.ID == uuid.Nil {
		t.Fatal("job ID was not assigned")
	}

	var got OCRJob
	if err := db.First(&got, "id = ?", job.ID).Error; err != nil {
		t.Fatalf("load OCR job: %v", err)
	}
	if got.UserID == nil || *got.UserID != user.ID {
		t.Fatalf("user_id = %#v, want %s", got.UserID, user.ID)
	}
	if got.SchemaID == nil || *got.SchemaID != schema.ID {
		t.Fatalf("schema_id = %#v, want %s", got.SchemaID, schema.ID)
	}
	if got.DocumentID == nil || *got.DocumentID != doc.ID {
		t.Fatalf("document_id = %#v, want %s", got.DocumentID, doc.ID)
	}
	if got.OriginalFilename != "invoice.pdf" || got.MimeType != "application/pdf" || got.FileSize != 42 {
		t.Fatalf("unexpected metadata: %#v", got)
	}
	if got.PageCount != 3 {
		t.Fatalf("page_count = %d, want 3", got.PageCount)
	}
	if got.DocumentHash != job.DocumentHash || got.FilePath != job.FilePath || got.Status != OCRJobStatusQueued {
		t.Fatalf("unexpected queue fields: %#v", got)
	}
	assertJSONEqual(t, got.InlineSchemaJSON, inlineSchema)
}

func testOCRJobRejectsInvalidStatus(t *testing.T) {
	db := ocrJobTx(t)

	t.Run("empty", func(t *testing.T) {
		job := OCRJob{
			OriginalFilename: "invoice.pdf",
			MimeType:         "application/pdf",
			FileSize:         42,
			PageCount:        1,
			DocumentHash:     "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			FilePath:         "/tmp/syncra-test/invoice.pdf",
			Status:           "",
		}

		assertInvalidStatusError(t, db.Create(&job).Error)
	})

	t.Run("unknown", func(t *testing.T) {
		job := OCRJob{
			OriginalFilename: "invoice.pdf",
			MimeType:         "application/pdf",
			FileSize:         42,
			PageCount:        1,
			DocumentHash:     "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			FilePath:         "/tmp/syncra-test/invoice.pdf",
			Status:           "done",
		}

		assertInvalidStatusError(t, db.Create(&job).Error)
	})
}

func testOpenPostgresTxAppliesOCRJobStatusDBCheck(t *testing.T) {
	db := ocrJobTx(t)

	if err := db.Exec(`SAVEPOINT ocr_job_raw_invalid_status`).Error; err != nil {
		t.Fatalf("create invalid status savepoint: %v", err)
	}
	invalidErr := db.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)}).Exec(`
INSERT INTO ocr_jobs (id, original_filename, mime_type, file_size, page_count, document_hash, file_path, status)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
`,
		uuid.New(),
		"invalid-raw-job.png",
		"image/png",
		int64(8),
		1,
		"invalid-raw-job-hash",
		"/tmp/invalid-raw-job.png",
		"done",
	).Error
	if err := db.Exec(`ROLLBACK TO SAVEPOINT ocr_job_raw_invalid_status`).Error; err != nil {
		t.Fatalf("rollback invalid status savepoint: %v", err)
	}
	if err := db.Exec(`RELEASE SAVEPOINT ocr_job_raw_invalid_status`).Error; err != nil {
		t.Fatalf("release invalid status savepoint: %v", err)
	}
	if invalidErr == nil {
		t.Fatal("raw SQL insert with invalid OCR job status succeeded, want database check error")
	}
}

func testOCRJobPartialStatusUpdatesValidateIncomingStatus(t *testing.T) {
	db := ocrJobTx(t)

	t.Run("valid status update succeeds with zero model receiver", func(t *testing.T) {
		job := createValidOCRJob(t, db)

		err := db.Model(&OCRJob{}).
			Where("id = ?", job.ID).
			Update("status", OCRJobStatusProcessing).
			Error
		if err != nil {
			t.Fatalf("update status to processing: %v", err)
		}

		var got OCRJob
		if err := db.First(&got, "id = ?", job.ID).Error; err != nil {
			t.Fatalf("load OCR job: %v", err)
		}
		if got.Status != OCRJobStatusProcessing {
			t.Fatalf("status = %q, want %q", got.Status, OCRJobStatusProcessing)
		}
	})

	t.Run("invalid status update fails with loaded model receiver", func(t *testing.T) {
		job := createValidOCRJob(t, db)

		var loaded OCRJob
		if err := db.First(&loaded, "id = ?", job.ID).Error; err != nil {
			t.Fatalf("load OCR job: %v", err)
		}

		err := db.Model(&loaded).Update("status", "done").Error
		assertInvalidStatusError(t, err)
	})

	t.Run("invalid selected zero status update fails", func(t *testing.T) {
		job := createValidOCRJob(t, db)

		err := db.Model(&OCRJob{}).
			Where("id = ?", job.ID).
			Select("status").
			Updates(OCRJob{Status: ""}).
			Error
		assertInvalidStatusError(t, err)
		assertStoredStatus(t, db, job.ID, OCRJobStatusQueued)
	})

	t.Run("invalid save status fails", func(t *testing.T) {
		job := createValidOCRJob(t, db)
		job.Status = "done"

		err := db.Save(&job).Error
		assertInvalidStatusError(t, err)
		assertStoredStatus(t, db, job.ID, OCRJobStatusQueued)
	})
}

func testOCRJobJSONRedactsFilePath(t *testing.T) {
	job := OCRJob{
		ID:               uuid.New(),
		OriginalFilename: "invoice.pdf",
		MimeType:         "application/pdf",
		FileSize:         42,
		PageCount:        1,
		DocumentHash:     "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		FilePath:         "/tmp/syncra-test/invoice.pdf",
		Status:           OCRJobStatusQueued,
	}

	payload, err := json.Marshal(job)
	if err != nil {
		t.Fatalf("marshal OCR job: %v", err)
	}
	if strings.Contains(string(payload), "file_path") {
		t.Fatalf("json contains file_path key: %s", string(payload))
	}
	if strings.Contains(string(payload), job.FilePath) {
		t.Fatalf("json contains file path value: %s", string(payload))
	}
}

func assertJSONEqual(t *testing.T, got datatypes.JSON, want datatypes.JSON) {
	t.Helper()
	var gotValue any
	if err := json.Unmarshal(got, &gotValue); err != nil {
		t.Fatalf("decode got JSON: %v", err)
	}
	var wantValue any
	if err := json.Unmarshal(want, &wantValue); err != nil {
		t.Fatalf("decode want JSON: %v", err)
	}
	if !reflect.DeepEqual(gotValue, wantValue) {
		t.Fatalf("json = %s, want %s", string(got), string(want))
	}
}

func assertInvalidStatusError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("operation succeeded, want invalid status error")
	}
	if !strings.Contains(err.Error(), "invalid OCR job status") {
		t.Fatalf("error = %q, want invalid status error", err.Error())
	}
}

func createValidOCRJob(t *testing.T, db *gorm.DB) OCRJob {
	t.Helper()

	job := OCRJob{
		OriginalFilename: "invoice.pdf",
		MimeType:         "application/pdf",
		FileSize:         42,
		PageCount:        1,
		DocumentHash:     "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		FilePath:         "/tmp/syncra-test/invoice.pdf",
		Status:           OCRJobStatusQueued,
	}
	if err := db.Create(&job).Error; err != nil {
		t.Fatalf("create OCR job: %v", err)
	}
	return job
}

func assertStoredStatus(t *testing.T, db *gorm.DB, id uuid.UUID, want OCRJobStatus) {
	t.Helper()

	var got OCRJob
	if err := db.First(&got, "id = ?", id).Error; err != nil {
		t.Fatalf("load OCR job: %v", err)
	}
	if got.Status != want {
		t.Fatalf("status = %q, want %q", got.Status, want)
	}
}
