package dbmigrate

import (
	"database/sql"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

const ocrJobStatusCheckConstraint = "chk_ocr_jobs_status"

var ocrJobStatuses = []string{"queued", "processing", "completed", "failed"}

// ValidateOCRJobStatuses checks existing OCR job rows before AutoMigrate can
// change status-column DDL.
func ValidateOCRJobStatuses(db *gorm.DB) error {
	_, statusExists, err := columnDataType(db, "ocr_jobs", "status")
	if err != nil || !statusExists {
		return err
	}
	return rejectInvalidOCRJobStatusRows(db)
}

// MigrateOCRJobStatus enforces the OCR job status contract for hook-skipping
// writes while staying safe to run before the OCR job table exists.
func MigrateOCRJobStatus(db *gorm.DB) error {
	_, statusExists, err := columnDataType(db, "ocr_jobs", "status")
	if err != nil || !statusExists {
		return err
	}

	defaultExpression, err := ocrJobStatusDefaultExpression(db)
	if err != nil {
		return err
	}
	constraintDefinition, err := ocrJobStatusConstraintDefinition(db)
	if err != nil {
		return err
	}
	statusNotNull, err := ocrJobStatusIsNotNull(db)
	if err != nil {
		return err
	}

	defaultCurrent := ocrJobStatusDefaultIsCurrent(defaultExpression)
	constraintCurrent := ocrJobStatusConstraintIsCurrent(constraintDefinition)
	if defaultCurrent && constraintCurrent && statusNotNull {
		return nil
	}

	if err := rejectInvalidOCRJobStatusRows(db); err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if !defaultCurrent {
			if err := tx.Exec(`ALTER TABLE "ocr_jobs" ALTER COLUMN "status" SET DEFAULT 'queued'`).Error; err != nil {
				return err
			}
		}
		if !statusNotNull {
			if err := tx.Exec(`ALTER TABLE "ocr_jobs" ALTER COLUMN "status" SET NOT NULL`).Error; err != nil {
				return err
			}
		}
		if !constraintCurrent {
			if err := tx.Exec(`ALTER TABLE "ocr_jobs" DROP CONSTRAINT IF EXISTS "` + ocrJobStatusCheckConstraint + `"`).Error; err != nil {
				return err
			}
			return tx.Exec(`ALTER TABLE "ocr_jobs" ADD CONSTRAINT "` + ocrJobStatusCheckConstraint + `" CHECK ("status" IN ('queued', 'processing', 'completed', 'failed'))`).Error
		}
		return nil
	})
}

func ocrJobStatusDefaultExpression(db *gorm.DB) (string, error) {
	var defaultExpression sql.NullString
	err := db.Raw(`
SELECT column_default
FROM information_schema.columns
WHERE table_schema = current_schema()
	AND table_name = 'ocr_jobs'
	AND column_name = 'status'
`).Scan(&defaultExpression).Error
	if err != nil || !defaultExpression.Valid {
		return "", err
	}
	return defaultExpression.String, nil
}

func ocrJobStatusConstraintDefinition(db *gorm.DB) (string, error) {
	var definition string
	err := db.Raw(`
SELECT pg_get_constraintdef(con.oid, true)
FROM pg_constraint con
JOIN pg_class rel ON rel.oid = con.conrelid
JOIN pg_namespace nsp ON nsp.oid = rel.relnamespace
WHERE con.contype = 'c'
	AND nsp.nspname = current_schema()
	AND rel.relname = 'ocr_jobs'
	AND con.conname = ?
`, ocrJobStatusCheckConstraint).Scan(&definition).Error
	return definition, err
}

func ocrJobStatusIsNotNull(db *gorm.DB) (bool, error) {
	var isNullable string
	err := db.Raw(`
SELECT is_nullable
FROM information_schema.columns
WHERE table_schema = current_schema()
	AND table_name = 'ocr_jobs'
	AND column_name = 'status'
`).Scan(&isNullable).Error
	return isNullable == "NO", err
}

func ocrJobStatusDefaultIsCurrent(defaultExpression string) bool {
	switch strings.ToLower(strings.TrimSpace(defaultExpression)) {
	case "'queued'", "'queued'::text", "'queued'::character varying":
		return true
	default:
		return false
	}
}

func ocrJobStatusConstraintIsCurrent(definition string) bool {
	normalized := normalizeConstraintDefinition(definition)
	expectedDefinitions := []string{
		`CHECK ("status" IN ('queued', 'processing', 'completed', 'failed'))`,
		`CHECK (status IN ('queued', 'processing', 'completed', 'failed'))`,
		`CHECK (status::text = ANY (ARRAY['queued'::character varying, 'processing'::character varying, 'completed'::character varying, 'failed'::character varying]::text[]))`,
		`CHECK (((status)::text = ANY ((ARRAY['queued'::character varying, 'processing'::character varying, 'completed'::character varying, 'failed'::character varying])::text[])))`,
	}
	for _, expected := range expectedDefinitions {
		if normalized == normalizeConstraintDefinition(expected) {
			return true
		}
	}
	return false
}

func normalizeConstraintDefinition(definition string) string {
	normalized := strings.ToLower(definition)
	normalized = strings.ReplaceAll(normalized, `"`, "")
	normalized = strings.Join(strings.Fields(normalized), "")
	return normalized
}

func rejectInvalidOCRJobStatusRows(db *gorm.DB) error {
	var invalidCount int64
	if err := db.Raw(`
SELECT COUNT(*)
FROM "ocr_jobs"
WHERE "status" IS NULL OR "status" NOT IN ?
`, ocrJobStatuses).Scan(&invalidCount).Error; err != nil {
		return err
	}
	if invalidCount > 0 {
		return fmt.Errorf("invalid OCR job status rows: %d row(s) have status outside queued, processing, completed, failed", invalidCount)
	}
	return nil
}
