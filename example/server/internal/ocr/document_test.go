package ocr

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"ai.ro/syncra/internal/auth"
	"ai.ro/syncra/internal/testsupport"
)

var ocrDocumentTestGroup *testsupport.PostgresGroup

func TestOCRDocuments(t *testing.T) {
	ocrDocumentTestGroup = testsupport.OpenPostgresGroup(t, &auth.User{}, &ExtractionSchema{}, &OCRDocument{}, &OCRJob{})
	defer func() { ocrDocumentTestGroup = nil }()

	for _, tt := range []struct {
		name string
		fn   func(*testing.T)
	}{
		{name: "AutoMigrateAndPersistJSONFields", fn: testAutoMigrateAndPersistJSONFields},
		{name: "BeforeCreateComputesPageCountFromRawResponseJSON", fn: testOCRDocumentBeforeCreateComputesPageCountFromRawResponseJSON},
		{name: "BeforeCreateDefaultsPageCountForMissingPages", fn: testOCRDocumentBeforeCreateDefaultsPageCountForMissingPages},
		{name: "BeforeCreateRejectsNonArrayPages", fn: testOCRDocumentBeforeCreateRejectsNonArrayPages},
		{name: "AutoMigrateAndPersistUserOwnedOCRFields", fn: testAutoMigrateAndPersistUserOwnedOCRFields},
		{name: "AutoMigrateAndPersistOCRDocumentJobID", fn: testAutoMigrateAndPersistOCRDocumentJobID},
	} {
		t.Run(tt.name, tt.fn)
	}
}

func ocrDocumentTx(t *testing.T) *gorm.DB {
	t.Helper()
	if ocrDocumentTestGroup != nil {
		return ocrDocumentTestGroup.Tx(t)
	}
	return testsupport.OpenPostgresTx(t, &auth.User{}, &ExtractionSchema{}, &OCRDocument{}, &OCRJob{})
}

func testAutoMigrateAndPersistJSONFields(t *testing.T) {
	db := ocrDocumentTx(t)

	expectedSchemaJSON := datatypes.JSON([]byte(`{"type":"object","properties":{"Furnizor":{"type":"string"}}}`))
	expectedAnnotationJSON := datatypes.JSON([]byte(`{"Furnizor":"Acme"}`))
	expectedRawResponseJSON := datatypes.JSON([]byte(`{"pages":[{"index":0},{"index":1}]}`))

	schema := ExtractionSchema{
		Name:        "invoice",
		Description: "Invoice extraction schema",
		SchemaJSON:  expectedSchemaJSON,
		Strict:      true,
	}
	if err := db.Create(&schema).Error; err != nil {
		t.Fatalf("create schema: %v", err)
	}

	doc := OCRDocument{
		OriginalFilename: "factura.pdf",
		MimeType:         "application/pdf",
		FileSize:         10,
		DocumentHash:     "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		SchemaID:         &schema.ID,
		Markdown:         "# Test",
		AnnotationJSON:   expectedAnnotationJSON,
		RawResponseJSON:  expectedRawResponseJSON,
	}
	if err := db.Create(&doc).Error; err != nil {
		t.Fatalf("create document: %v", err)
	}

	var got OCRDocument
	if err := db.First(&got, "id = ?", doc.ID).Error; err != nil {
		t.Fatalf("load document: %v", err)
	}
	assertJSONEqualLocal(t, "annotation JSON", expectedAnnotationJSON, got.AnnotationJSON)
	assertJSONEqualLocal(t, "raw response JSON", expectedRawResponseJSON, got.RawResponseJSON)
	if got.PageCount != 2 {
		t.Fatalf("page_count = %d, want 2", got.PageCount)
	}

	var gotSchema ExtractionSchema
	if err := db.First(&gotSchema, "id = ?", schema.ID).Error; err != nil {
		t.Fatalf("load schema: %v", err)
	}
	assertJSONEqualLocal(t, "schema JSON", expectedSchemaJSON, gotSchema.SchemaJSON)
}

func testOCRDocumentBeforeCreateComputesPageCountFromRawResponseJSON(t *testing.T) {
	db := ocrDocumentTx(t)
	doc := minimalOCRDocumentForPageCountTest(`{"pages":[{"index":0},{"index":1}]}`)

	if err := db.Create(&doc).Error; err != nil {
		t.Fatalf("create OCR document: %v", err)
	}

	var got OCRDocument
	if err := db.First(&got, "id = ?", doc.ID).Error; err != nil {
		t.Fatalf("load OCR document: %v", err)
	}
	if got.PageCount != 2 {
		t.Fatalf("page_count = %d, want 2", got.PageCount)
	}
}

func testOCRDocumentBeforeCreateDefaultsPageCountForMissingPages(t *testing.T) {
	tests := []struct {
		name string
		raw  string
	}{
		{name: "missing pages", raw: `{"model":"mistral-ocr-latest"}`},
		{name: "null pages", raw: `{"pages":null}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := ocrDocumentTx(t)
			doc := minimalOCRDocumentForPageCountTest(tt.raw)

			if err := db.Create(&doc).Error; err != nil {
				t.Fatalf("create OCR document: %v", err)
			}

			var got OCRDocument
			if err := db.First(&got, "id = ?", doc.ID).Error; err != nil {
				t.Fatalf("load OCR document: %v", err)
			}
			if got.PageCount != 0 {
				t.Fatalf("page_count = %d, want 0", got.PageCount)
			}
		})
	}
}

func testOCRDocumentBeforeCreateRejectsNonArrayPages(t *testing.T) {
	db := ocrDocumentTx(t)
	doc := minimalOCRDocumentForPageCountTest(`{"pages":{"index":0}}`)

	err := db.Create(&doc).Error
	if err == nil {
		t.Fatal("create OCR document error = nil, want error")
	}
	if !strings.Contains(err.Error(), "invalid OCR raw response pages") {
		t.Fatalf("error = %v, want invalid raw response pages", err)
	}
}

func testAutoMigrateAndPersistUserOwnedOCRFields(t *testing.T) {
	db := ocrDocumentTx(t)

	user := auth.User{
		Name:  "Test User",
		Email: "owned-ocr@example.com",
	}
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
		OriginalFilename: "factura.pdf",
		MimeType:         "application/pdf",
		FileSize:         10,
		DocumentHash:     "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		SchemaID:         &schema.ID,
		Markdown:         "# Test",
		RawResponseJSON:  datatypes.JSON([]byte(`{"pages":[]}`)),
	}
	if err := db.Create(&doc).Error; err != nil {
		t.Fatalf("create document: %v", err)
	}

	var gotSchema ExtractionSchema
	if err := db.First(&gotSchema, "id = ?", schema.ID).Error; err != nil {
		t.Fatalf("load schema: %v", err)
	}
	if gotSchema.UserID == nil {
		t.Fatal("schema user ID is nil")
	}
	if *gotSchema.UserID != user.ID {
		t.Fatalf("schema user ID = %q, want %q", *gotSchema.UserID, user.ID)
	}

	var gotDoc OCRDocument
	if err := db.First(&gotDoc, "id = ?", doc.ID).Error; err != nil {
		t.Fatalf("load document: %v", err)
	}
	if gotDoc.UserID == nil {
		t.Fatal("document user ID is nil")
	}
	if *gotDoc.UserID != user.ID {
		t.Fatalf("document user ID = %q, want %q", *gotDoc.UserID, user.ID)
	}
}

func testAutoMigrateAndPersistOCRDocumentJobID(t *testing.T) {
	db := ocrDocumentTx(t)

	job := OCRJob{
		OriginalFilename: "invoice.pdf",
		MimeType:         "application/pdf",
		FileSize:         42,
		PageCount:        2,
		DocumentHash:     "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		FilePath:         "/tmp/syncra-test/invoice.pdf",
		Status:           OCRJobStatusQueued,
	}
	if err := db.Create(&job).Error; err != nil {
		t.Fatalf("create OCR job: %v", err)
	}

	doc := OCRDocument{
		OriginalFilename: "invoice.pdf",
		MimeType:         "application/pdf",
		FileSize:         42,
		DocumentHash:     job.DocumentHash,
		JobID:            &job.ID,
		Markdown:         "# Invoice",
		RawResponseJSON:  datatypes.JSON([]byte(`{"pages":[]}`)),
	}
	if err := db.Create(&doc).Error; err != nil {
		t.Fatalf("create OCR document: %v", err)
	}

	var got OCRDocument
	if err := db.First(&got, "id = ?", doc.ID).Error; err != nil {
		t.Fatalf("load OCR document: %v", err)
	}
	if got.JobID == nil || *got.JobID != job.ID {
		t.Fatalf("job_id = %#v, want %s", got.JobID, job.ID)
	}
}

func minimalOCRDocumentForPageCountTest(rawResponseJSON string) OCRDocument {
	return OCRDocument{
		OriginalFilename: "scan.png",
		MimeType:         "image/png",
		FileSize:         10,
		PageCount:        99,
		DocumentHash:     "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		Markdown:         "# Test",
		RawResponseJSON:  datatypes.JSON([]byte(rawResponseJSON)),
	}
}

func assertJSONEqualLocal(t *testing.T, name string, want datatypes.JSON, got datatypes.JSON) {
	t.Helper()

	var wantValue any
	if err := json.Unmarshal(want, &wantValue); err != nil {
		t.Fatalf("unmarshal expected %s: %v", name, err)
	}

	var gotValue any
	if err := json.Unmarshal(got, &gotValue); err != nil {
		t.Fatalf("unmarshal actual %s: %v", name, err)
	}

	if !reflect.DeepEqual(wantValue, gotValue) {
		t.Fatalf("%s = %s, want %s", name, got, want)
	}
}
