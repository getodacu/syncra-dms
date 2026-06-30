package dbmigrate

import (
	"fmt"

	"gorm.io/gorm"
)

const ocrDocumentJobForeignKeyConstraint = "fk_ocr_documents_job"

// MigrateOCRDocumentJobForeignKey ensures the reverse OCR document-to-job link
// has a database constraint after AutoMigrate has created both OCR tables.
func MigrateOCRDocumentJobForeignKey(db *gorm.DB) error {
	_, documentJobIDExists, err := columnDataType(db, "ocr_documents", "job_id")
	if err != nil || !documentJobIDExists {
		return err
	}
	_, jobIDExists, err := columnDataType(db, "ocr_jobs", "id")
	if err != nil || !jobIDExists {
		return err
	}

	constraints, err := ocrDocumentJobForeignKeyConstraints(db)
	if err != nil {
		return err
	}

	if len(constraints) == 1 &&
		constraints[0].Name == ocrDocumentJobForeignKeyConstraint &&
		constraints[0].DeleteAction == "n" &&
		constraints[0].UpdateAction == "c" &&
		constraints[0].LocalColumn == "job_id" &&
		constraints[0].ReferencedTable == "ocr_jobs" &&
		constraints[0].ReferencedColumn == "id" &&
		constraints[0].LocalColumnCount == 1 &&
		constraints[0].ReferencedColumnCount == 1 {
		return nil
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for _, constraint := range constraints {
			if err := tx.Exec(fmt.Sprintf(
				`ALTER TABLE %s DROP CONSTRAINT %s`,
				quoteIdentifier("ocr_documents"),
				quoteIdentifier(constraint.Name),
			)).Error; err != nil {
				return err
			}
		}

		return tx.Exec(fmt.Sprintf(
			`ALTER TABLE %s ADD CONSTRAINT %s FOREIGN KEY ("job_id") REFERENCES %s("id") ON UPDATE CASCADE ON DELETE SET NULL`,
			quoteIdentifier("ocr_documents"),
			quoteIdentifier(ocrDocumentJobForeignKeyConstraint),
			quoteIdentifier("ocr_jobs"),
		)).Error
	})
}

type ocrDocumentJobForeignKeyConstraintInfo struct {
	Name                  string `gorm:"column:name"`
	DeleteAction          string `gorm:"column:delete_action"`
	UpdateAction          string `gorm:"column:update_action"`
	LocalColumn           string `gorm:"column:local_column"`
	ReferencedTable       string `gorm:"column:referenced_table"`
	ReferencedColumn      string `gorm:"column:referenced_column"`
	LocalColumnCount      int    `gorm:"column:local_column_count"`
	ReferencedColumnCount int    `gorm:"column:referenced_column_count"`
}

func ocrDocumentJobForeignKeyConstraints(db *gorm.DB) ([]ocrDocumentJobForeignKeyConstraintInfo, error) {
	var constraints []ocrDocumentJobForeignKeyConstraintInfo
	err := db.Raw(`
SELECT con.conname AS name,
	con.confdeltype AS delete_action,
	con.confupdtype AS update_action,
	local_att.attname AS local_column,
	ref_rel.relname AS referenced_table,
	ref_att.attname AS referenced_column,
	cardinality(con.conkey) AS local_column_count,
	cardinality(con.confkey) AS referenced_column_count
FROM pg_constraint con
JOIN pg_class rel ON rel.oid = con.conrelid
JOIN pg_namespace nsp ON nsp.oid = rel.relnamespace
JOIN pg_class ref_rel ON ref_rel.oid = con.confrelid
JOIN pg_namespace ref_nsp ON ref_nsp.oid = ref_rel.relnamespace
JOIN pg_attribute local_att ON local_att.attrelid = con.conrelid AND local_att.attnum = con.conkey[1]
JOIN pg_attribute ref_att ON ref_att.attrelid = con.confrelid AND ref_att.attnum = con.confkey[1]
WHERE con.contype = 'f'
	AND nsp.nspname = current_schema()
	AND rel.relname = 'ocr_documents'
	AND ref_nsp.nspname = current_schema()
	AND (
		con.conname = ?
		OR EXISTS (
			SELECT 1
			FROM pg_attribute any_local_att
			WHERE any_local_att.attrelid = con.conrelid
				AND any_local_att.attnum = ANY(con.conkey)
				AND any_local_att.attname = 'job_id'
		)
	)
`, ocrDocumentJobForeignKeyConstraint).Scan(&constraints).Error
	if err != nil {
		return nil, err
	}
	return dedupeOCRDocumentJobForeignKeyConstraints(constraints), nil
}

func dedupeOCRDocumentJobForeignKeyConstraints(constraints []ocrDocumentJobForeignKeyConstraintInfo) []ocrDocumentJobForeignKeyConstraintInfo {
	seen := make(map[string]struct{}, len(constraints))
	deduped := make([]ocrDocumentJobForeignKeyConstraintInfo, 0, len(constraints))
	for _, constraint := range constraints {
		if _, ok := seen[constraint.Name]; ok {
			continue
		}
		seen[constraint.Name] = struct{}{}
		deduped = append(deduped, constraint)
	}
	return deduped
}
