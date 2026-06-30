package config

import (
	"strings"
	"testing"
)

func TestLoadReadsEnvironmentAndDefaults(t *testing.T) {
	t.Setenv("DSN", "host=localhost user=syncra password=syncra dbname=syncra_dev port=5432 sslmode=disable TimeZone=Europe/Bucharest")
	t.Setenv("DSN_DEV", "host=localhost user=syncra password=syncra dbname=syncra_dev port=5432 sslmode=disable TimeZone=Europe/Bucharest")
	t.Setenv("ATLAS_DATABASE_URL", "postgres://syncra:syncra@localhost:5432/syncra_dev?sslmode=disable")
	t.Setenv("ATLAS_DEV_DATABASE_URL", "postgres://syncra:syncra@localhost:5432/syncra_atlas?sslmode=disable")
	t.Setenv("SERVER_HOST_PORT", "127.0.0.1:9090")
	t.Setenv("DEBUG", "true")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.ServerHostPort != "127.0.0.1:9090" {
		t.Fatalf("ServerHostPort = %q", cfg.ServerHostPort)
	}
	if !cfg.Debug {
		t.Fatal("Debug = false, want true")
	}
	if cfg.DSN == "" || cfg.DSNDev == "" {
		t.Fatal("database DSNs must be loaded")
	}
	if cfg.AtlasDatabaseURL == "" || cfg.AtlasDevDatabaseURL == "" {
		t.Fatal("Atlas URLs must be loaded")
	}
	if cfg.Version != "dev" {
		t.Fatalf("Version = %q, want dev", cfg.Version)
	}
}

func TestLoadRequiresDSN(t *testing.T) {
	t.Setenv("DSN_DEV", "host=localhost dbname=syncra_dev")
	t.Setenv("ATLAS_DATABASE_URL", "postgres://syncra:syncra@localhost:5432/syncra_dev?sslmode=disable")
	t.Setenv("ATLAS_DEV_DATABASE_URL", "postgres://syncra:syncra@localhost:5432/syncra_atlas?sslmode=disable")

	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "DSN is required") {
		t.Fatalf("Load() error = %v, want DSN required", err)
	}
}

func TestLoadRejectsAtlasDevDatabaseURLPointingAtAppDatabase(t *testing.T) {
	t.Setenv("DSN", "host=localhost dbname=syncra_dev")
	t.Setenv("DSN_DEV", "host=localhost dbname=syncra_dev")
	t.Setenv("ATLAS_DATABASE_URL", "postgres://syncra:syncra@localhost:5432/syncra_dev?sslmode=disable")
	t.Setenv("ATLAS_DEV_DATABASE_URL", "postgres://syncra:syncra@localhost:5432/syncra_dev?sslmode=disable")

	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "ATLAS_DEV_DATABASE_URL must not target syncra_dev") {
		t.Fatalf("Load() error = %v, want Atlas dev database safety error", err)
	}
}

func TestDatabaseNameFromDSN(t *testing.T) {
	tests := []struct {
		name string
		dsn  string
		want string
	}{
		{
			name: "postgres keyword dsn",
			dsn:  "host=localhost user=syncra dbname=syncra_dev sslmode=disable",
			want: "syncra_dev",
		},
		{
			name: "postgres url",
			dsn:  "postgres://syncra:syncra@localhost:5432/syncra_atlas?sslmode=disable",
			want: "syncra_atlas",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DatabaseNameFromDSN(tt.dsn)
			if err != nil {
				t.Fatalf("DatabaseNameFromDSN() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("DatabaseNameFromDSN() = %q, want %q", got, tt.want)
			}
		})
	}
}
