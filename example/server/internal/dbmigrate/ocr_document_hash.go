package dbmigrate

import "gorm.io/gorm"

// ResetLegacyIntegerIDTables is kept for older callers. It now preserves rows
// by migrating legacy integer IDs to UUIDs in place.
func ResetLegacyIntegerIDTables(db *gorm.DB) error {
	return MigrateLegacyIntegerIDTables(db)
}

// MigrateLegacyIntegerIDTables converts pre-UUID OCR table IDs in place instead
// of dropping tables. Existing numeric IDs are kept in _syncra_legacy_id columns
// so old-to-new mappings remain inspectable after the migration.
func MigrateLegacyIntegerIDTables(db *gorm.DB) error {
	legacy, err := hasLegacyIntegerIDColumns(db)
	if err != nil || !legacy {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := migrateLegacyExtractionSchemaIDs(tx); err != nil {
			return err
		}
		return migrateLegacyOCRDocumentIDs(tx)
	})
}

func hasLegacyIntegerIDColumns(db *gorm.DB) (bool, error) {
	checks := []struct {
		table  string
		column string
	}{
		{table: "extraction_schemas", column: "id"},
		{table: "ocr_documents", column: "id"},
		{table: "ocr_documents", column: "schema_id"},
	}

	for _, check := range checks {
		dataType, exists, err := columnDataType(db, check.table, check.column)
		if err != nil {
			return false, err
		}
		if exists && dataType != "uuid" {
			return true, nil
		}
	}
	return false, nil
}

func migrateLegacyExtractionSchemaIDs(db *gorm.DB) error {
	dataType, exists, err := columnDataType(db, "extraction_schemas", "id")
	if err != nil || !exists || dataType == "uuid" {
		return err
	}

	if err := dropReferencingConstraints(db, "extraction_schemas"); err != nil {
		return err
	}
	if err := dropTableConstraints(db, "extraction_schemas", true); err != nil {
		return err
	}

	for _, stmt := range []string{
		`ALTER TABLE "extraction_schemas" ADD COLUMN IF NOT EXISTS "_syncra_legacy_id" bigint`,
		`UPDATE "extraction_schemas" SET "_syncra_legacy_id" = "id" WHERE "_syncra_legacy_id" IS NULL`,
		`ALTER TABLE "extraction_schemas" ADD COLUMN IF NOT EXISTS "_syncra_uuid_id" uuid`,
	} {
		if err := db.Exec(stmt).Error; err != nil {
			return err
		}
	}
	if err := assignLegacyUUID(db, "extraction_schemas", "_syncra_uuid_id", "id", "syncra-extraction-schema"); err != nil {
		return err
	}
	for _, stmt := range []string{
		`ALTER TABLE "extraction_schemas" DROP COLUMN "id"`,
		`ALTER TABLE "extraction_schemas" RENAME COLUMN "_syncra_uuid_id" TO "id"`,
		`ALTER TABLE "extraction_schemas" ALTER COLUMN "id" SET NOT NULL`,
		`ALTER TABLE "extraction_schemas" ADD PRIMARY KEY ("id")`,
	} {
		if err := db.Exec(stmt).Error; err != nil {
			return err
		}
	}
	return nil
}

func migrateLegacyOCRDocumentIDs(db *gorm.DB) error {
	idDataType, idExists, err := columnDataType(db, "ocr_documents", "id")
	if err != nil || !idExists {
		return err
	}
	schemaIDDataType, schemaIDExists, err := columnDataType(db, "ocr_documents", "schema_id")
	if err != nil {
		return err
	}

	legacyID := idDataType != "uuid"
	legacySchemaID := schemaIDExists && schemaIDDataType != "uuid"
	if !legacyID && !legacySchemaID {
		return nil
	}

	if legacyID {
		if err := dropReferencingConstraints(db, "ocr_documents"); err != nil {
			return err
		}
	}
	if err := dropTableConstraints(db, "ocr_documents", legacyID); err != nil {
		return err
	}

	if legacyID {
		for _, stmt := range []string{
			`ALTER TABLE "ocr_documents" ADD COLUMN IF NOT EXISTS "_syncra_legacy_id" bigint`,
			`UPDATE "ocr_documents" SET "_syncra_legacy_id" = "id" WHERE "_syncra_legacy_id" IS NULL`,
			`ALTER TABLE "ocr_documents" ADD COLUMN IF NOT EXISTS "_syncra_uuid_id" uuid`,
		} {
			if err := db.Exec(stmt).Error; err != nil {
				return err
			}
		}
		if err := assignLegacyUUID(db, "ocr_documents", "_syncra_uuid_id", "id", "syncra-ocr-document"); err != nil {
			return err
		}
	}

	if legacySchemaID {
		if err := db.Exec(`ALTER TABLE "ocr_documents" ADD COLUMN IF NOT EXISTS "_syncra_uuid_schema_id" uuid`).Error; err != nil {
			return err
		}
		_, hasLegacySchemaID, err := columnDataType(db, "extraction_schemas", "_syncra_legacy_id")
		if err != nil {
			return err
		}
		if hasLegacySchemaID {
			if err := db.Exec(`
	UPDATE "ocr_documents" AS doc
	SET "_syncra_uuid_schema_id" = schema."id"
	FROM "extraction_schemas" AS schema
	WHERE doc."schema_id" = schema."_syncra_legacy_id"
	`).Error; err != nil {
				return err
			}
		}
		if err := db.Exec(`ALTER TABLE "ocr_documents" DROP COLUMN "schema_id"`).Error; err != nil {
			return err
		}
		if err := db.Exec(`ALTER TABLE "ocr_documents" RENAME COLUMN "_syncra_uuid_schema_id" TO "schema_id"`).Error; err != nil {
			return err
		}
	}

	if legacyID {
		for _, stmt := range []string{
			`ALTER TABLE "ocr_documents" DROP COLUMN "id"`,
			`ALTER TABLE "ocr_documents" RENAME COLUMN "_syncra_uuid_id" TO "id"`,
			`ALTER TABLE "ocr_documents" ALTER COLUMN "id" SET NOT NULL`,
			`ALTER TABLE "ocr_documents" ADD PRIMARY KEY ("id")`,
		} {
			if err := db.Exec(stmt).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func assignLegacyUUID(db *gorm.DB, tableName string, targetColumn string, sourceColumn string, namespace string) error {
	table := quoteIdentifier(tableName)
	target := quoteIdentifier(targetColumn)
	source := quoteIdentifier(sourceColumn)
	return db.Exec(`
	UPDATE `+table+` AS item
	SET `+target+` = (
		substr(mapped.hash, 1, 8) || '-' ||
		substr(mapped.hash, 9, 4) || '-' ||
		substr(mapped.hash, 13, 4) || '-' ||
		substr(mapped.hash, 17, 4) || '-' ||
		substr(mapped.hash, 21, 12)
	)::uuid
	FROM (
		SELECT `+source+` AS legacy_id,
			md5(? || ':' || `+source+`::text) AS hash
		FROM `+table+`
	) AS mapped
	WHERE item.`+source+` = mapped.legacy_id
		AND item.`+target+` IS NULL
	`, namespace).Error
}

type tableConstraint struct {
	TableName      string `gorm:"column:table_name"`
	ConstraintName string `gorm:"column:constraint_name"`
}

func dropTableConstraints(db *gorm.DB, tableName string, includePrimary bool) error {
	constraintTypes := []string{"f"}
	if includePrimary {
		constraintTypes = append(constraintTypes, "p")
	}

	var constraintNames []string
	if err := db.Raw(`
	SELECT con.conname
	FROM pg_constraint con
	JOIN pg_class rel ON rel.oid = con.conrelid
	JOIN pg_namespace nsp ON nsp.oid = rel.relnamespace
	WHERE nsp.nspname = current_schema()
		AND rel.relname = ?
		AND con.contype IN ?
	`, tableName, constraintTypes).Scan(&constraintNames).Error; err != nil {
		return err
	}
	for _, constraintName := range constraintNames {
		if err := db.Exec(`ALTER TABLE ` + quoteIdentifier(tableName) + ` DROP CONSTRAINT IF EXISTS ` + quoteIdentifier(constraintName)).Error; err != nil {
			return err
		}
	}
	return nil
}

func dropReferencingConstraints(db *gorm.DB, referencedTable string) error {
	var constraints []tableConstraint
	if err := db.Raw(`
	SELECT rel.relname AS table_name,
		con.conname AS constraint_name
	FROM pg_constraint con
	JOIN pg_class rel ON rel.oid = con.conrelid
	JOIN pg_namespace nsp ON nsp.oid = rel.relnamespace
	JOIN pg_class ref_rel ON ref_rel.oid = con.confrelid
	JOIN pg_namespace ref_nsp ON ref_nsp.oid = ref_rel.relnamespace
	WHERE con.contype = 'f'
		AND nsp.nspname = current_schema()
		AND ref_nsp.nspname = current_schema()
		AND ref_rel.relname = ?
	`, referencedTable).Scan(&constraints).Error; err != nil {
		return err
	}
	for _, constraint := range constraints {
		if err := db.Exec(`ALTER TABLE ` + quoteIdentifier(constraint.TableName) + ` DROP CONSTRAINT IF EXISTS ` + quoteIdentifier(constraint.ConstraintName)).Error; err != nil {
			return err
		}
	}
	return nil
}

// MigrateOCRDocumentHash migrates legacy file-only OCR hashes without deleting
// OCR rows. When no legacy hash exists, it assigns a stable per-row placeholder
// that cannot be used for future content dedupe but keeps historical rows valid.
func MigrateOCRDocumentHash(db *gorm.DB) error {
	_, tableExists, err := columnDataType(db, "ocr_documents", "id")
	if err != nil || !tableExists {
		return err
	}

	_, hasFileSHA, err := columnDataType(db, "ocr_documents", "file_sha256")
	if err != nil {
		return err
	}
	_, hasDocumentHash, err := columnDataType(db, "ocr_documents", "document_hash")
	if err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		switch {
		case hasFileSHA && !hasDocumentHash:
			if err := tx.Exec(`ALTER TABLE "ocr_documents" RENAME COLUMN "file_sha256" TO "document_hash"`).Error; err != nil {
				return err
			}
		case hasFileSHA && hasDocumentHash:
			if err := tx.Exec(`
	UPDATE "ocr_documents"
	SET "document_hash" = "file_sha256"
	WHERE ("document_hash" IS NULL OR "document_hash" = '')
		AND "file_sha256" IS NOT NULL
		AND "file_sha256" <> ''
	`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`ALTER TABLE "ocr_documents" DROP COLUMN "file_sha256"`).Error; err != nil {
				return err
			}
		case !hasFileSHA && !hasDocumentHash:
			if err := tx.Exec(`ALTER TABLE "ocr_documents" ADD COLUMN "document_hash" text`).Error; err != nil {
				return err
			}
		}

		if err := tx.Exec(`
	UPDATE "ocr_documents"
	SET "document_hash" = md5('syncra-legacy-ocr-document:' || "id"::text) ||
		md5('syncra-legacy-ocr-document-v2:' || "id"::text)
	WHERE "document_hash" IS NULL OR "document_hash" = ''
	`).Error; err != nil {
			return err
		}
		return tx.Exec(`ALTER TABLE "ocr_documents" ALTER COLUMN "document_hash" SET NOT NULL`).Error
	})
}

func columnDataType(db *gorm.DB, tableName string, columnName string) (string, bool, error) {
	var dataType string
	err := db.Raw(`
SELECT data_type
FROM information_schema.columns
WHERE table_schema = current_schema()
	AND table_name = ?
	AND column_name = ?
`, tableName, columnName).Scan(&dataType).Error
	if err != nil {
		return "", false, err
	}
	return dataType, dataType != "", nil
}
