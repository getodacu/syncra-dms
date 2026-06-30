package dbmigrate_test

import (
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"ai.ro/syncra/internal/dbmigrate"
	"ai.ro/syncra/internal/testsupport"
)

func TestMigrateOCRDocumentHashPreservesLegacyFileSHA256Rows(t *testing.T) {
	db := testsupport.OpenPostgresTx(t)
	if err := db.Migrator().DropTable("ocr_documents"); err != nil {
		t.Fatalf("drop OCR documents table: %v", err)
	}
	if err := db.Exec(`
	CREATE TABLE "ocr_documents" (
		"id" uuid PRIMARY KEY,
		"created_at" timestamptz NOT NULL,
		"updated_at" timestamptz NOT NULL,
		"original_filename" varchar(255) NOT NULL,
		"mime_type" varchar(120) NOT NULL,
		"file_size" bigint NOT NULL,
		"file_sha256" varchar(64) NOT NULL,
		"markdown" text NOT NULL,
		"raw_response_json" jsonb NOT NULL
	)
	`).Error; err != nil {
		t.Fatalf("create legacy OCR documents table: %v", err)
	}

	id := uuid.New()
	const fileSHA = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	if err := db.Exec(`
	INSERT INTO "ocr_documents" (
		"id",
		"created_at",
		"updated_at",
		"original_filename",
		"mime_type",
		"file_size",
		"file_sha256",
		"markdown",
		"raw_response_json"
	) VALUES (?, now(), now(), ?, ?, ?, ?, ?, ?::jsonb)
	`, id, "legacy.pdf", "application/pdf", 42, fileSHA, "# OCR", `{"pages":[]}`).Error; err != nil {
		t.Fatalf("insert legacy OCR document: %v", err)
	}

	if err := dbmigrate.MigrateOCRDocumentHash(db); err != nil {
		t.Fatalf("MigrateOCRDocumentHash() error = %v", err)
	}

	var gotHash string
	if err := db.Raw(`SELECT "document_hash" FROM "ocr_documents" WHERE "id" = ?`, id).Scan(&gotHash).Error; err != nil {
		t.Fatalf("select migrated document hash: %v", err)
	}
	if gotHash != fileSHA {
		t.Fatalf("document_hash = %q, want %q", gotHash, fileSHA)
	}
	if ocrDocumentColumnExists(t, db, "file_sha256") {
		t.Fatal("ocr_documents.file_sha256 still exists")
	}
	assertOCRDocumentRowCount(t, db, 1)
}

func assertOCRDocumentRowCount(t *testing.T, db *gorm.DB, want int64) {
	t.Helper()

	var count int64
	if err := db.Raw(`SELECT COUNT(*) FROM "ocr_documents"`).Scan(&count).Error; err != nil {
		t.Fatalf("count OCR documents: %v", err)
	}
	if count != want {
		t.Fatalf("OCR document count = %d, want %d", count, want)
	}
}
