package database

import (
	"testing"

	"ai.ro/syncra/dms/internal/orgunits"
)

func TestApplicationModelsIncludesDomainModels(t *testing.T) {
	got := ApplicationModels()
	if len(got) < 5 {
		t.Fatalf("ApplicationModels() length = %d, want at least 5 models", len(got))
	}
	foundOrganizationUnit := false
	for _, model := range got {
		if _, ok := model.(*orgunits.Unit); ok {
			foundOrganizationUnit = true
		}
	}
	if !foundOrganizationUnit {
		t.Fatal("ApplicationModels() does not include orgunits.Unit")
	}
}

func TestOpenPostgresRequiresDSN(t *testing.T) {
	db, err := OpenPostgres("")
	if err == nil {
		t.Fatalf("OpenPostgres(\"\") error = nil, db = %v", db)
	}
}
