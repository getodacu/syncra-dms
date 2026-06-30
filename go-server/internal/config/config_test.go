package config

import (
	"strings"
	"testing"
)

func TestLoadReadsEnvironmentAndDefaults(t *testing.T) {
	t.Setenv("DSN", "host=localhost user=syncra password=syncra dbname=syncra_dms port=5432 sslmode=disable TimeZone=Europe/Bucharest")
	t.Setenv("DSN_DEV", "host=localhost user=syncra password=syncra dbname=syncra_dms_dev port=5432 sslmode=disable TimeZone=Europe/Bucharest")
	t.Setenv("ATLAS_DATABASE_URL", "postgres://syncra:syncra@localhost:5432/syncra_dms?sslmode=disable")
	t.Setenv("ATLAS_DEV_DATABASE_URL", "postgres://syncra:syncra@localhost:5432/syncra_atlas?sslmode=disable")
	t.Setenv("SERVER_HOST_PORT", "127.0.0.1:9090")
	t.Setenv("DEBUG", "true")
	t.Setenv("BETTER_AUTH_SECRET", "better-auth-secret-from-env")
	t.Setenv("AUTH_DELIVERY_TOKEN", "delivery-token")
	t.Setenv("AUTH_SESSION_TTL_SECONDS", "3600")
	t.Setenv("AUTH_VERIFICATION_TTL_SECONDS", "900")
	t.Setenv("AUTH_COOKIE_SECURE", "true")
	t.Setenv("GOOGLE_CLIENT_ID", "google-client")
	t.Setenv("GOOGLE_CLIENT_SECRET", "google-secret")
	t.Setenv("GITHUB_CLIENT_ID", "github-client")
	t.Setenv("GITHUB_CLIENT_SECRET", "github-secret")

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
	if cfg.BetterAuthSecret != "better-auth-secret-from-env" || cfg.AuthDeliveryToken != "delivery-token" {
		t.Fatalf("auth secret/token were not loaded")
	}
	if cfg.AuthSessionTTLSeconds != 3600 || cfg.AuthVerificationTTLSeconds != 900 || !cfg.AuthCookieSecure {
		t.Fatalf("auth TTL/cookie config = %d/%d/%v", cfg.AuthSessionTTLSeconds, cfg.AuthVerificationTTLSeconds, cfg.AuthCookieSecure)
	}
	if cfg.GoogleClientID != "google-client" || cfg.GoogleClientSecret != "google-secret" {
		t.Fatalf("google oauth env was not loaded")
	}
	if cfg.GitHubClientID != "github-client" || cfg.GitHubClientSecret != "github-secret" {
		t.Fatalf("github oauth env was not loaded")
	}
}

func TestLoadRequiresDSN(t *testing.T) {
	t.Setenv("DSN_DEV", "host=localhost dbname=syncra_dms_dev")

	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "DSN is required") {
		t.Fatalf("Load() error = %v, want DSN required", err)
	}
}

func TestLoadDoesNotValidateAtlasDevDatabaseURLDuringAPIRuntime(t *testing.T) {
	t.Setenv("DSN", "host=localhost dbname=syncra_dms")
	t.Setenv("DSN_DEV", "host=localhost dbname=syncra_dms_dev")
	t.Setenv("ATLAS_DATABASE_URL", "postgres://syncra:syncra@localhost:5432/syncra_dms?sslmode=disable")
	t.Setenv("ATLAS_DEV_DATABASE_URL", "postgres://syncra:syncra@localhost:5432/syncra_dms_dev?sslmode=disable")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if !strings.Contains(cfg.AtlasDevDatabaseURL, "syncra_dms_dev") {
		t.Fatalf("AtlasDevDatabaseURL = %q, want configured value loaded", cfg.AtlasDevDatabaseURL)
	}
}

func TestLoadDoesNotRequireAtlasURLsForAPIRuntime(t *testing.T) {
	t.Setenv("DSN", "host=localhost dbname=syncra_dms")
	t.Setenv("DSN_DEV", "host=localhost dbname=syncra_dms_dev")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.AtlasDatabaseURL != "" || cfg.AtlasDevDatabaseURL != "" {
		t.Fatalf("Atlas URLs = %q/%q, want empty when unset", cfg.AtlasDatabaseURL, cfg.AtlasDevDatabaseURL)
	}
}

func TestLoadMigrationRequiresOnlyAtlasDatabaseURL(t *testing.T) {
	t.Setenv("ATLAS_DATABASE_URL", "  postgres://syncra:syncra@localhost:5432/syncra_dms?sslmode=disable  ")

	cfg, err := LoadMigration()
	if err != nil {
		t.Fatalf("LoadMigration() error = %v", err)
	}
	if cfg.AtlasDatabaseURL != "postgres://syncra:syncra@localhost:5432/syncra_dms?sslmode=disable" {
		t.Fatalf("AtlasDatabaseURL = %q", cfg.AtlasDatabaseURL)
	}
}

func TestLoadMigrationRequiresAtlasDatabaseURL(t *testing.T) {
	_, err := LoadMigration()
	if err == nil || !strings.Contains(err.Error(), "ATLAS_DATABASE_URL is required") {
		t.Fatalf("LoadMigration() error = %v, want ATLAS_DATABASE_URL required", err)
	}
}

func TestLoadRejectsDSNDevPointingAtAppDatabase(t *testing.T) {
	t.Setenv("DSN", "host=localhost dbname=syncra_dms")
	t.Setenv("DSN_DEV", "host=localhost dbname=syncra_dms")
	t.Setenv("ATLAS_DATABASE_URL", "postgres://syncra:syncra@localhost:5432/syncra_dms?sslmode=disable")
	t.Setenv("ATLAS_DEV_DATABASE_URL", "postgres://syncra:syncra@localhost:5432/syncra_atlas?sslmode=disable")

	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "DSN_DEV must target syncra_dms_dev") {
		t.Fatalf("Load() error = %v, want DSN_DEV database safety error", err)
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
			dsn:  "host=localhost user=syncra dbname=syncra_dms sslmode=disable",
			want: "syncra_dms",
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
