package main

import (
	"strings"
	"testing"
)

func TestLoadSchemaIncludesApplicationTables(t *testing.T) {
	schema, err := loadSchema()
	if err != nil {
		t.Fatalf("loadSchema() error = %v", err)
	}

	for _, table := range []string{
		`CREATE TABLE "user"`,
		`CREATE TABLE "ocr_documents"`,
		`CREATE TABLE "ocr_jobs"`,
		`CREATE TABLE "collections"`,
	} {
		if !strings.Contains(schema, table) {
			t.Fatalf("schema does not contain %q:\n%s", table, schema)
		}
	}
}
