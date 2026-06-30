package dbmigrate

import "gorm.io/gorm"

// MigrateOCRDocumentLifecycle removes legacy lifecycle fields from successful
// OCR result rows. Queue state and failure details belong to ocr_jobs.
func MigrateOCRDocumentLifecycle(db *gorm.DB) error {
	_, tableExists, err := columnDataType(db, "ocr_documents", "id")
	if err != nil || !tableExists {
		return err
	}

	_, hasStatus, err := columnDataType(db, "ocr_documents", "status")
	if err != nil {
		return err
	}
	_, hasErrorMessage, err := columnDataType(db, "ocr_documents", "error_message")
	if err != nil {
		return err
	}
	if !hasStatus && !hasErrorMessage {
		return nil
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if hasStatus {
			if err := tx.Exec(`ALTER TABLE "ocr_documents" DROP COLUMN "status"`).Error; err != nil {
				return err
			}
		}
		if hasErrorMessage {
			if err := tx.Exec(`ALTER TABLE "ocr_documents" DROP COLUMN "error_message"`).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
