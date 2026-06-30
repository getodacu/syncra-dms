package database

import "testing"

func TestApplicationModelsIncludesAuthModels(t *testing.T) {
	if got := ApplicationModels(); len(got) != 4 {
		t.Fatalf("ApplicationModels() length = %d, want 4 auth models", len(got))
	}
}

func TestOpenPostgresRequiresDSN(t *testing.T) {
	db, err := OpenPostgres("")
	if err == nil {
		t.Fatalf("OpenPostgres(\"\") error = nil, db = %v", db)
	}
}
