package dbmigrate

import "gorm.io/gorm"

// MigrateOCRDocumentListIndexes adds indexes used by OCR document listing and
// filename search.
func MigrateOCRDocumentListIndexes(db *gorm.DB) error {
	if !db.Migrator().HasTable("ocr_documents") {
		return nil
	}

	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS pg_trgm`).Error; err != nil {
		return err
	}
	if err := db.Exec(`CREATE INDEX IF NOT EXISTS "idx_ocr_documents_user_created_id" ON "ocr_documents" ("user_id", "created_at", "id")`).Error; err != nil {
		return err
	}
	return db.Exec(`CREATE INDEX IF NOT EXISTS "idx_ocr_documents_original_filename_trgm" ON "ocr_documents" USING gin (lower("original_filename") gin_trgm_ops)`).Error
}
