package dbmigrate

import "gorm.io/gorm"

// MigrateOCRDocumentPageCount ensures OCR documents expose the number of pages
// returned by the OCR provider response.
func MigrateOCRDocumentPageCount(db *gorm.DB) error {
	_, tableExists, err := columnDataType(db, "ocr_documents", "id")
	if err != nil || !tableExists {
		return err
	}

	_, hasPageCount, err := columnDataType(db, "ocr_documents", "page_count")
	if err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if !hasPageCount {
			if err := tx.Exec(`ALTER TABLE "ocr_documents" ADD COLUMN "page_count" bigint NOT NULL DEFAULT 0`).Error; err != nil {
				return err
			}
		}
		if err := tx.Exec(`
UPDATE "ocr_documents"
SET "page_count" = CASE
	WHEN jsonb_typeof("raw_response_json"->'pages') = 'array'
	THEN jsonb_array_length("raw_response_json"->'pages')
	ELSE 0
END
`).Error; err != nil {
			return err
		}
		if err := tx.Exec(`ALTER TABLE "ocr_documents" ALTER COLUMN "page_count" SET DEFAULT 0`).Error; err != nil {
			return err
		}
		return tx.Exec(`ALTER TABLE "ocr_documents" ALTER COLUMN "page_count" SET NOT NULL`).Error
	})
}
