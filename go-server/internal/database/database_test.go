package database

import (
	"testing"

	"ai.ro/syncra/dms/internal/documents"
	"ai.ro/syncra/dms/internal/orgunits"
)

func TestApplicationModelsIncludesDomainModels(t *testing.T) {
	got := ApplicationModels()
	if len(got) < 5 {
		t.Fatalf("ApplicationModels() length = %d, want at least 5 models", len(got))
	}
	foundOrganizationUnit := false
	foundDocumentFolder := false
	foundDocument := false
	for _, model := range got {
		if _, ok := model.(*orgunits.Unit); ok {
			foundOrganizationUnit = true
		}
		if _, ok := model.(*documents.Folder); ok {
			foundDocumentFolder = true
		}
		if _, ok := model.(*documents.Document); ok {
			foundDocument = true
		}
	}
	if !foundOrganizationUnit {
		t.Fatal("ApplicationModels() does not include orgunits.Unit")
	}
	if !foundDocumentFolder {
		t.Fatal("ApplicationModels() does not include documents.Folder")
	}
	if !foundDocument {
		t.Fatal("ApplicationModels() does not include documents.Document")
	}
}

func TestOpenPostgresRequiresDSN(t *testing.T) {
	db, err := OpenPostgres("")
	if err == nil {
		t.Fatalf("OpenPostgres(\"\") error = nil, db = %v", db)
	}
}
