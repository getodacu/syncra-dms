package dbmigrate_test

import (
	"testing"

	"gorm.io/gorm"

	"ai.ro/syncra/internal/auth"
	"ai.ro/syncra/internal/dbmigrate"
	"ai.ro/syncra/internal/ocr"
	"ai.ro/syncra/internal/testsupport"
)

func TestMigrateOCRDocumentLifecycleDropsDocumentStatusAndErrorMessage(t *testing.T) {
	db := testsupport.OpenPostgresTx(t, &auth.User{}, &ocr.ExtractionSchema{}, &ocr.OCRDocument{}, &ocr.OCRJob{})

	for _, column := range []string{"status", "error_message"} {
		exists := ocrDocumentColumnExists(t, db, column)
		if !exists {
			if err := db.Exec(`ALTER TABLE "ocr_documents" ADD COLUMN "` + column + `" text`).Error; err != nil {
				t.Fatalf("add OCR document %s column: %v", column, err)
			}
		}
	}

	if err := dbmigrate.MigrateOCRDocumentLifecycle(db); err != nil {
		t.Fatalf("migrate OCR document lifecycle: %v", err)
	}

	for _, column := range []string{"status", "error_message"} {
		if ocrDocumentColumnExists(t, db, column) {
			t.Fatalf("ocr_documents.%s still exists", column)
		}
	}
}

func ocrDocumentColumnExists(t *testing.T, db *gorm.DB, column string) bool {
	t.Helper()

	var count int64
	if err := db.Raw(`
SELECT COUNT(*)
FROM information_schema.columns
WHERE table_schema = current_schema()
	AND table_name = 'ocr_documents'
	AND column_name = ?
`, column).Scan(&count).Error; err != nil {
		t.Fatalf("check OCR document column %s: %v", column, err)
	}
	return count > 0
}
