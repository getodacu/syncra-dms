package database

import "testing"

func TestApplicationModelsStartsEmpty(t *testing.T) {
	if got := ApplicationModels(); len(got) != 0 {
		t.Fatalf("ApplicationModels() length = %d, want 0 for lean scaffold", len(got))
	}
}

func TestOpenPostgresRequiresDSN(t *testing.T) {
	db, err := OpenPostgres("")
	if err == nil {
		t.Fatalf("OpenPostgres(\"\") error = nil, db = %v", db)
	}
}
