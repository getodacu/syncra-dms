package dbmigrate_test

import (
	"testing"

	"gorm.io/gorm"

	"ai.ro/syncra/internal/auth"
	"ai.ro/syncra/internal/dbmigrate"
	"ai.ro/syncra/internal/ocr"
	"ai.ro/syncra/internal/testsupport"
)

func TestMigrateOCRDocumentJobForeignKeyReplacesWrongReferencedColumn(t *testing.T) {
	db := testsupport.OpenPostgresTx(t, &auth.User{}, &ocr.ExtractionSchema{}, &ocr.OCRDocument{}, &ocr.OCRJob{})

	if err := db.Exec(`ALTER TABLE "ocr_documents" DROP CONSTRAINT "fk_ocr_documents_job"`).Error; err != nil {
		t.Fatalf("drop current OCR document job FK: %v", err)
	}
	if err := db.Exec(`ALTER TABLE "ocr_jobs" ADD COLUMN "wrong_job_id" uuid UNIQUE`).Error; err != nil {
		t.Fatalf("add wrong OCR job target column: %v", err)
	}
	if err := db.Exec(`
ALTER TABLE "ocr_documents"
ADD CONSTRAINT "fk_ocr_documents_job"
FOREIGN KEY ("job_id") REFERENCES "ocr_jobs"("wrong_job_id")
ON UPDATE CASCADE ON DELETE SET NULL
`).Error; err != nil {
		t.Fatalf("add stale OCR document job FK: %v", err)
	}

	if err := dbmigrate.MigrateOCRDocumentJobForeignKey(db); err != nil {
		t.Fatalf("migrate OCR document job FK: %v", err)
	}

	assertOCRDocumentJobForeignKeyTarget(t, db, "job_id", "id", "n", "c", 1)
}

func TestMigrateOCRDocumentJobForeignKeyReplacesWrongLocalColumn(t *testing.T) {
	db := testsupport.OpenPostgresTx(t, &auth.User{}, &ocr.ExtractionSchema{}, &ocr.OCRDocument{}, &ocr.OCRJob{})

	if err := db.Exec(`ALTER TABLE "ocr_documents" DROP CONSTRAINT "fk_ocr_documents_job"`).Error; err != nil {
		t.Fatalf("drop current OCR document job FK: %v", err)
	}
	if err := db.Exec(`ALTER TABLE "ocr_documents" ADD COLUMN "wrong_job_id" uuid`).Error; err != nil {
		t.Fatalf("add wrong OCR document local column: %v", err)
	}
	if err := db.Exec(`
ALTER TABLE "ocr_documents"
ADD CONSTRAINT "fk_ocr_documents_job"
FOREIGN KEY ("wrong_job_id") REFERENCES "ocr_jobs"("id")
ON UPDATE CASCADE ON DELETE SET NULL
`).Error; err != nil {
		t.Fatalf("add stale OCR document job FK: %v", err)
	}

	if err := dbmigrate.MigrateOCRDocumentJobForeignKey(db); err != nil {
		t.Fatalf("migrate OCR document job FK: %v", err)
	}

	assertOCRDocumentJobForeignKeyTarget(t, db, "job_id", "id", "n", "c", 1)
}

func TestMigrateOCRDocumentJobForeignKeyReplacesWrongReferencedTable(t *testing.T) {
	db := testsupport.OpenPostgresTx(t, &auth.User{}, &ocr.ExtractionSchema{}, &ocr.OCRDocument{}, &ocr.OCRJob{})

	if err := db.Exec(`ALTER TABLE "ocr_documents" DROP CONSTRAINT "fk_ocr_documents_job"`).Error; err != nil {
		t.Fatalf("drop current OCR document job FK: %v", err)
	}
	if err := db.Exec(`CREATE TABLE "wrong_ocr_jobs" ("id" uuid PRIMARY KEY)`).Error; err != nil {
		t.Fatalf("create wrong OCR jobs table: %v", err)
	}
	if err := db.Exec(`
ALTER TABLE "ocr_documents"
ADD CONSTRAINT "fk_ocr_documents_wrong_job"
FOREIGN KEY ("job_id") REFERENCES "wrong_ocr_jobs"("id")
ON UPDATE CASCADE ON DELETE SET NULL
`).Error; err != nil {
		t.Fatalf("add stale OCR document job FK: %v", err)
	}

	if err := dbmigrate.MigrateOCRDocumentJobForeignKey(db); err != nil {
		t.Fatalf("migrate OCR document job FK: %v", err)
	}

	assertOCRDocumentJobForeignKeyTarget(t, db, "job_id", "id", "n", "c", 1)
	assertForeignKeyAbsent(t, db, "fk_ocr_documents_wrong_job")
}

func TestMigrateOCRDocumentJobForeignKeyReplacesCompositeWithJobIDSecond(t *testing.T) {
	db := testsupport.OpenPostgresTx(t, &auth.User{}, &ocr.ExtractionSchema{}, &ocr.OCRDocument{}, &ocr.OCRJob{})

	if err := db.Exec(`ALTER TABLE "ocr_documents" DROP CONSTRAINT "fk_ocr_documents_job"`).Error; err != nil {
		t.Fatalf("drop current OCR document job FK: %v", err)
	}
	if err := db.Exec(`ALTER TABLE "ocr_documents" ADD COLUMN "wrong_job_id" uuid`).Error; err != nil {
		t.Fatalf("add wrong OCR document local column: %v", err)
	}
	if err := db.Exec(`ALTER TABLE "ocr_jobs" ADD COLUMN "wrong_job_id" uuid`).Error; err != nil {
		t.Fatalf("add wrong OCR job target column: %v", err)
	}
	if err := db.Exec(`ALTER TABLE "ocr_jobs" ADD CONSTRAINT "uq_ocr_jobs_wrong_job_id_id" UNIQUE ("wrong_job_id", "id")`).Error; err != nil {
		t.Fatalf("add composite OCR job unique constraint: %v", err)
	}
	if err := db.Exec(`
ALTER TABLE "ocr_documents"
ADD CONSTRAINT "fk_ocr_documents_composite_wrong_job"
FOREIGN KEY ("wrong_job_id", "job_id") REFERENCES "ocr_jobs"("wrong_job_id", "id")
ON UPDATE CASCADE ON DELETE SET NULL
`).Error; err != nil {
		t.Fatalf("add stale composite OCR document job FK: %v", err)
	}

	if err := dbmigrate.MigrateOCRDocumentJobForeignKey(db); err != nil {
		t.Fatalf("migrate OCR document job FK: %v", err)
	}

	assertOCRDocumentJobForeignKeyTarget(t, db, "job_id", "id", "n", "c", 1)
	assertForeignKeyAbsent(t, db, "fk_ocr_documents_composite_wrong_job")
}

func TestMigrateOCRDocumentJobForeignKeyKeepsCurrentConstraint(t *testing.T) {
	db := testsupport.OpenPostgresTx(t, &auth.User{}, &ocr.ExtractionSchema{}, &ocr.OCRDocument{}, &ocr.OCRJob{})

	beforeOID := ocrDocumentJobForeignKeyOID(t, db)
	if err := dbmigrate.MigrateOCRDocumentJobForeignKey(db); err != nil {
		t.Fatalf("migrate OCR document job FK: %v", err)
	}
	afterOID := ocrDocumentJobForeignKeyOID(t, db)
	if afterOID != beforeOID {
		t.Fatalf("OCR document job FK oid changed after idempotent migrate: before %s, after %s", beforeOID, afterOID)
	}
}

func assertOCRDocumentJobForeignKeyTarget(t *testing.T, db *gorm.DB, wantLocalColumn string, wantReferencedColumn string, wantDelete string, wantUpdate string, wantColumns int) {
	t.Helper()
	var constraints []struct {
		LocalColumn      string `gorm:"column:local_column"`
		ReferencedColumn string `gorm:"column:referenced_column"`
		DeleteAction     string `gorm:"column:delete_action"`
		UpdateAction     string `gorm:"column:update_action"`
		ColumnCount      int    `gorm:"column:column_count"`
	}
	if err := db.Raw(`
SELECT local_att.attname AS local_column,
	ref_att.attname AS referenced_column,
	con.confdeltype AS delete_action,
	con.confupdtype AS update_action,
	cardinality(con.conkey) AS column_count
FROM pg_constraint con
JOIN pg_class rel ON rel.oid = con.conrelid
JOIN pg_namespace nsp ON nsp.oid = rel.relnamespace
JOIN pg_attribute local_att ON local_att.attrelid = con.conrelid AND local_att.attnum = con.conkey[1]
JOIN pg_attribute ref_att ON ref_att.attrelid = con.confrelid AND ref_att.attnum = con.confkey[1]
WHERE con.contype = 'f'
	AND nsp.nspname = current_schema()
	AND rel.relname = 'ocr_documents'
	AND con.conname = 'fk_ocr_documents_job'
	AND con.confrelid = to_regclass(format('%I.%I', current_schema(), 'ocr_jobs'))
`).Scan(&constraints).Error; err != nil {
		t.Fatalf("query OCR document job FK target: %v", err)
	}
	if len(constraints) != 1 {
		t.Fatalf("OCR document job FK count = %d, want 1", len(constraints))
	}
	got := constraints[0]
	if got.LocalColumn != wantLocalColumn {
		t.Fatalf("local column = %q, want %q", got.LocalColumn, wantLocalColumn)
	}
	if got.ReferencedColumn != wantReferencedColumn {
		t.Fatalf("referenced column = %q, want %q", got.ReferencedColumn, wantReferencedColumn)
	}
	if got.DeleteAction != wantDelete {
		t.Fatalf("delete action = %q, want %q", got.DeleteAction, wantDelete)
	}
	if got.UpdateAction != wantUpdate {
		t.Fatalf("update action = %q, want %q", got.UpdateAction, wantUpdate)
	}
	if got.ColumnCount != wantColumns {
		t.Fatalf("column count = %d, want %d", got.ColumnCount, wantColumns)
	}
}

func assertForeignKeyAbsent(t *testing.T, db *gorm.DB, constraintName string) {
	t.Helper()
	var count int64
	if err := db.Raw(`
SELECT COUNT(*)
FROM pg_constraint con
JOIN pg_class rel ON rel.oid = con.conrelid
JOIN pg_namespace nsp ON nsp.oid = rel.relnamespace
WHERE con.contype = 'f'
	AND nsp.nspname = current_schema()
	AND rel.relname = 'ocr_documents'
	AND con.conname = ?
`, constraintName).Scan(&count).Error; err != nil {
		t.Fatalf("query stale OCR document job FK: %v", err)
	}
	if count != 0 {
		t.Fatalf("foreign key %q count = %d, want 0", constraintName, count)
	}
}

func ocrDocumentJobForeignKeyOID(t *testing.T, db *gorm.DB) string {
	t.Helper()
	var oid string
	if err := db.Raw(`
SELECT con.oid::text
FROM pg_constraint con
JOIN pg_class rel ON rel.oid = con.conrelid
JOIN pg_namespace nsp ON nsp.oid = rel.relnamespace
WHERE con.contype = 'f'
	AND nsp.nspname = current_schema()
	AND rel.relname = 'ocr_documents'
	AND con.conname = 'fk_ocr_documents_job'
`).Scan(&oid).Error; err != nil {
		t.Fatalf("query OCR document job FK oid: %v", err)
	}
	if oid == "" {
		t.Fatal("OCR document job FK oid is empty")
	}
	return oid
}
