package dbmigrate

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type ownerForeignKey struct {
	tableName      string
	constraintName string
}

var ownerForeignKeys = []ownerForeignKey{
	{tableName: "extraction_schemas", constraintName: "fk_extraction_schemas_user"},
	{tableName: "ocr_documents", constraintName: "fk_ocr_documents_user"},
	{tableName: "ocr_jobs", constraintName: "fk_ocr_jobs_user"},
	{tableName: "collections", constraintName: "fk_collections_user"},
}

// MigrateOwnerForeignKeyCascades ensures user-owned OCR rows are deleted with
// their owner instead of being converted to system-wide rows by stale FKs.
func MigrateOwnerForeignKeyCascades(db *gorm.DB) error {
	for _, fk := range ownerForeignKeys {
		if err := migrateOwnerForeignKeyCascade(db, fk); err != nil {
			return err
		}
	}
	return nil
}

func migrateOwnerForeignKeyCascade(db *gorm.DB, fk ownerForeignKey) error {
	_, tableExists, err := columnDataType(db, fk.tableName, "user_id")
	if err != nil || !tableExists {
		return err
	}

	var constraintNames []string
	if err := db.Raw(`
SELECT con.conname
FROM pg_constraint con
JOIN pg_class rel ON rel.oid = con.conrelid
JOIN pg_namespace nsp ON nsp.oid = rel.relnamespace
JOIN pg_attribute att ON att.attrelid = con.conrelid AND att.attnum = ANY(con.conkey)
WHERE con.contype = 'f'
	AND nsp.nspname = current_schema()
	AND rel.relname = ?
	AND att.attname = 'user_id'
	AND con.confrelid = to_regclass(format('%I.%I', current_schema(), 'user'))
`, fk.tableName).Scan(&constraintNames).Error; err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for _, constraintName := range constraintNames {
			if err := tx.Exec(fmt.Sprintf(
				`ALTER TABLE %s DROP CONSTRAINT %s`,
				quoteIdentifier(fk.tableName),
				quoteIdentifier(constraintName),
			)).Error; err != nil {
				return err
			}
		}

		return tx.Exec(fmt.Sprintf(
			`ALTER TABLE %s ADD CONSTRAINT %s FOREIGN KEY ("user_id") REFERENCES "user"("id") ON UPDATE CASCADE ON DELETE CASCADE`,
			quoteIdentifier(fk.tableName),
			quoteIdentifier(fk.constraintName),
		)).Error
	})
}

func quoteIdentifier(identifier string) string {
	return `"` + strings.ReplaceAll(identifier, `"`, `""`) + `"`
}
