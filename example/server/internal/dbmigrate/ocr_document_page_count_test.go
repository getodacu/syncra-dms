package dbmigrate_test

import (
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"ai.ro/syncra/internal/auth"
	"ai.ro/syncra/internal/dbmigrate"
	"ai.ro/syncra/internal/ocr"
	"ai.ro/syncra/internal/testsupport"
)

func TestMigrateOCRDocumentPageCountAddsAndBackfillsColumn(t *testing.T) {
	db := testsupport.OpenPostgresTx(t, &auth.User{}, &ocr.ExtractionSchema{}, &ocr.OCRDocument{}, &ocr.OCRJob{})

	if err := db.Exec(`ALTER TABLE "ocr_documents" DROP COLUMN IF EXISTS "page_count"`).Error; err != nil {
		t.Fatalf("drop OCR document page_count column: %v", err)
	}

	withPagesID := insertOCRDocumentWithoutPageCount(t, db, `{"pages":[{"index":0},{"index":1}]}`)
	missingPagesID := insertOCRDocumentWithoutPageCount(t, db, `{"model":"mistral-ocr-latest"}`)
	nonArrayPagesID := insertOCRDocumentWithoutPageCount(t, db, `{"pages":{"index":0}}`)

	if err := dbmigrate.MigrateOCRDocumentPageCount(db); err != nil {
		t.Fatalf("migrate OCR document page count: %v", err)
	}
	if err := dbmigrate.MigrateOCRDocumentPageCount(db); err != nil {
		t.Fatalf("repeat OCR document page count migration: %v", err)
	}

	if !ocrDocumentColumnExists(t, db, "page_count") {
		t.Fatal("ocr_documents.page_count does not exist")
	}
	assertOCRDocumentPageCount(t, db, withPagesID, 2)
	assertOCRDocumentPageCount(t, db, missingPagesID, 0)
	assertOCRDocumentPageCount(t, db, nonArrayPagesID, 0)
}

func insertOCRDocumentWithoutPageCount(t *testing.T, db *gorm.DB, rawResponseJSON string) uuid.UUID {
	t.Helper()

	id := uuid.New()
	if err := db.Exec(`
INSERT INTO "ocr_documents" (
	"id",
	"created_at",
	"updated_at",
	"original_filename",
	"mime_type",
	"file_size",
	"document_hash",
	"markdown",
	"raw_response_json"
) VALUES (?, now(), now(), ?, ?, ?, ?, ?, ?::jsonb)
`, id, id.String()+".pdf", "application/pdf", 42, documentHashForID(id), "# OCR", rawResponseJSON).Error; err != nil {
		t.Fatalf("insert OCR document without page_count: %v", err)
	}
	return id
}

func documentHashForID(id uuid.UUID) string {
	hash := id.String() + id.String()
	if len(hash) > 64 {
		return hash[:64]
	}
	return hash
}

func assertOCRDocumentPageCount(t *testing.T, db *gorm.DB, id uuid.UUID, want int) {
	t.Helper()

	var got int
	if err := db.Raw(`SELECT "page_count" FROM "ocr_documents" WHERE "id" = ?`, id).Scan(&got).Error; err != nil {
		t.Fatalf("load OCR document page_count: %v", err)
	}
	if got != want {
		t.Fatalf("page_count for %s = %d, want %d", id, got, want)
	}
}
