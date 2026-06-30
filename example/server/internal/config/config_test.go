package config

import (
	"os"
	"slices"
	"testing"
)

func TestLoadFromEnv(t *testing.T) {
	t.Setenv("SERVER_HOST_PORT", "127.0.0.1:9090")
	t.Setenv("DEBUG", "false")
	t.Setenv("DSN", "host=localhost user=test dbname=syncra_test port=5432 sslmode=disable")
	t.Setenv("SWAGGER_HOST", "127.0.0.1:9090")
	t.Setenv("SWAGGER_SCHEMES", "http, https")
	t.Setenv("MISTRAL_API_KEY", "test-key")
	t.Setenv("MISTRAL_API_BASE_URL", "https://api.mistral.ai")
	t.Setenv("MISTRAL_OCR_MODEL", "mistral-ocr-latest")
	t.Setenv("MAX_UPLOAD_BYTES", "1048576")
	t.Setenv("STORAGE_DIR", "  /var/lib/syncra  ")
	t.Setenv("BETTER_AUTH_SECRET", "vN8qR4zT7mK2pL9xC5wY3hD6sF0aJ1uB")
	t.Setenv("APP_PRIVATE_KEY", "test-private-key-32-byte-material")
	t.Setenv("AUTH_DELIVERY_TOKEN", "trusted-delivery-token")
	t.Setenv("SYNCRA_INTERNAL_API_TOKEN", " trusted-internal-token ")
	t.Setenv("AUTH_SESSION_TTL_SECONDS", "604800")
	t.Setenv("AUTH_VERIFICATION_TTL_SECONDS", "300")
	t.Setenv("AUTH_COOKIE_SECURE", "false")
	t.Setenv("ONBOARDING_CREDITS", "250")
	t.Setenv("OCR_EXECUTOR_WORKERS", "4")
	t.Setenv("OCR_EXECUTOR_POLL_INTERVAL_SECONDS", "7")
	t.Setenv("OCR_EXECUTOR_QUEUE_BUFFER", "16")
	t.Setenv("GOTENBERG_API_URL", " gotenberg.local:3000/forms/chromium/convert/url ")
	t.Setenv("GOOGLE_CLIENT_ID", " google-client-id ")
	t.Setenv("GOOGLE_CLIENT_SECRET", " google-client-secret ")
	t.Setenv("GITHUB_CLIENT_ID", " github-client-id ")
	t.Setenv("GITHUB_CLIENT_SECRET", " github-client-secret ")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.ServerHostPort != "127.0.0.1:9090" {
		t.Fatalf("ServerHostPort = %q", cfg.ServerHostPort)
	}
	if cfg.Debug {
		t.Fatal("Debug should be false")
	}
	if cfg.MaxUploadBytes != 1048576 {
		t.Fatalf("MaxUploadBytes = %d", cfg.MaxUploadBytes)
	}
	if cfg.StorageDir != "/var/lib/syncra" {
		t.Fatalf("StorageDir = %q", cfg.StorageDir)
	}
	if cfg.MistralOCRModel != "mistral-ocr-latest" {
		t.Fatalf("MistralOCRModel = %q", cfg.MistralOCRModel)
	}
	if !slices.Equal(cfg.SwaggerSchemes, []string{"http", "https"}) {
		t.Fatalf("SwaggerSchemes = %#v", cfg.SwaggerSchemes)
	}
	if cfg.BetterAuthSecret != "vN8qR4zT7mK2pL9xC5wY3hD6sF0aJ1uB" {
		t.Fatalf("BetterAuthSecret = %q", cfg.BetterAuthSecret)
	}
	if cfg.AppPrivateKey != "test-private-key-32-byte-material" {
		t.Fatalf("AppPrivateKey = %q", cfg.AppPrivateKey)
	}
	if cfg.AuthDeliveryToken != "trusted-delivery-token" {
		t.Fatalf("AuthDeliveryToken = %q", cfg.AuthDeliveryToken)
	}
	if cfg.InternalAPIToken != "trusted-internal-token" {
		t.Fatalf("InternalAPIToken = %q", cfg.InternalAPIToken)
	}
	if cfg.AuthSessionTTL != 604800 {
		t.Fatalf("AuthSessionTTL = %d", cfg.AuthSessionTTL)
	}
	if cfg.AuthVerificationTTL != 300 {
		t.Fatalf("AuthVerificationTTL = %d", cfg.AuthVerificationTTL)
	}
	if cfg.AuthCookieSecure {
		t.Fatal("AuthCookieSecure should be false")
	}
	if cfg.OnboardingCredits != 250 {
		t.Fatalf("OnboardingCredits = %d", cfg.OnboardingCredits)
	}
	if cfg.OCRExecutorWorkers != 4 {
		t.Fatalf("OCRExecutorWorkers = %d", cfg.OCRExecutorWorkers)
	}
	if cfg.OCRExecutorPollIntervalSeconds != 7 {
		t.Fatalf("OCRExecutorPollIntervalSeconds = %d", cfg.OCRExecutorPollIntervalSeconds)
	}
	if cfg.OCRExecutorQueueBuffer != 16 {
		t.Fatalf("OCRExecutorQueueBuffer = %d", cfg.OCRExecutorQueueBuffer)
	}
	if cfg.GotenbergAPIURL != "http://gotenberg.local:3000/forms/chromium/convert/html" {
		t.Fatalf("GotenbergAPIURL = %q", cfg.GotenbergAPIURL)
	}
	if cfg.GoogleClientID != "google-client-id" {
		t.Fatalf("GoogleClientID = %q", cfg.GoogleClientID)
	}
	if cfg.GoogleClientSecret != "google-client-secret" {
		t.Fatalf("GoogleClientSecret = %q", cfg.GoogleClientSecret)
	}
	if cfg.GitHubClientID != "github-client-id" {
		t.Fatalf("GitHubClientID = %q", cfg.GitHubClientID)
	}
	if cfg.GitHubClientSecret != "github-client-secret" {
		t.Fatalf("GitHubClientSecret = %q", cfg.GitHubClientSecret)
	}
}

func TestLoadTrimsDSN(t *testing.T) {
	t.Setenv("DSN", "  host=localhost user=test dbname=syncra_test  ")
	t.Setenv("MISTRAL_API_KEY", "test-key")
	t.Setenv("BETTER_AUTH_SECRET", "vN8qR4zT7mK2pL9xC5wY3hD6sF0aJ1uB")
	t.Setenv("APP_PRIVATE_KEY", "test-private-key-32-byte-material")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.DSN != "host=localhost user=test dbname=syncra_test" {
		t.Fatalf("DSN = %q", cfg.DSN)
	}
}

func TestLoadRequiresDSN(t *testing.T) {
	t.Setenv("DSN", "")
	t.Setenv("MISTRAL_API_KEY", "test-key")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() expected DSN error")
	}
}

func TestLoadRequiresDSNAfterTrim(t *testing.T) {
	t.Setenv("DSN", " \t\n ")
	t.Setenv("MISTRAL_API_KEY", "test-key")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() expected DSN error")
	}
}

func TestLoadUsesMistralDefaults(t *testing.T) {
	t.Setenv("DSN", "host=localhost user=test dbname=syncra_test")
	t.Setenv("MISTRAL_API_KEY", "test-key")
	t.Setenv("BETTER_AUTH_SECRET", "vN8qR4zT7mK2pL9xC5wY3hD6sF0aJ1uB")
	t.Setenv("APP_PRIVATE_KEY", "test-private-key-32-byte-material")
	unsetenv(t, "MISTRAL_API_BASE_URL")
	unsetenv(t, "MISTRAL_OCR_MODEL")
	unsetenv(t, "MAX_UPLOAD_BYTES")
	unsetenv(t, "STORAGE_DIR")
	t.Setenv("OCR_JOB_FILE_DIR", "/legacy/ocr-jobs")
	t.Setenv("BILLING_INVOICE_PDF_DIR", "/legacy/invoices")
	unsetenv(t, "OCR_EXECUTOR_WORKERS")
	unsetenv(t, "OCR_EXECUTOR_POLL_INTERVAL_SECONDS")
	unsetenv(t, "OCR_EXECUTOR_QUEUE_BUFFER")
	unsetenv(t, "GOTENBERG_API_URL")
	unsetenv(t, "ONBOARDING_CREDITS")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.MistralBaseURL != "https://api.mistral.ai" {
		t.Fatalf("MistralBaseURL = %q", cfg.MistralBaseURL)
	}
	if cfg.MistralOCRModel != "mistral-ocr-latest" {
		t.Fatalf("MistralOCRModel = %q", cfg.MistralOCRModel)
	}
	if cfg.MaxUploadBytes != 20<<20 {
		t.Fatalf("MaxUploadBytes = %d", cfg.MaxUploadBytes)
	}
	if cfg.StorageDir != "" {
		t.Fatalf("StorageDir = %q, want empty", cfg.StorageDir)
	}
	if cfg.OCRExecutorWorkers != 2 {
		t.Fatalf("OCRExecutorWorkers = %d", cfg.OCRExecutorWorkers)
	}
	if cfg.OCRExecutorPollIntervalSeconds != 10 {
		t.Fatalf("OCRExecutorPollIntervalSeconds = %d", cfg.OCRExecutorPollIntervalSeconds)
	}
	if cfg.OCRExecutorQueueBuffer != 8 {
		t.Fatalf("OCRExecutorQueueBuffer = %d", cfg.OCRExecutorQueueBuffer)
	}
	if cfg.GotenbergAPIURL != "" {
		t.Fatalf("GotenbergAPIURL = %q, want empty", cfg.GotenbergAPIURL)
	}
	if cfg.OnboardingCredits != 100 {
		t.Fatalf("OnboardingCredits = %d, want 100", cfg.OnboardingCredits)
	}
}

func TestLoadTrimsMistralBaseURLTrailingSlash(t *testing.T) {
	t.Setenv("DSN", "host=localhost user=test dbname=syncra_test")
	t.Setenv("MISTRAL_API_KEY", "test-key")
	t.Setenv("BETTER_AUTH_SECRET", "vN8qR4zT7mK2pL9xC5wY3hD6sF0aJ1uB")
	t.Setenv("APP_PRIVATE_KEY", "test-private-key-32-byte-material")
	t.Setenv("MISTRAL_API_BASE_URL", "https://api.mistral.ai///")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.MistralBaseURL != "https://api.mistral.ai" {
		t.Fatalf("MistralBaseURL = %q", cfg.MistralBaseURL)
	}
}

func TestLoadFallsBackForInvalidDebugAndMaxUploadBytes(t *testing.T) {
	tests := []struct {
		name           string
		debug          string
		maxUploadBytes string
	}{
		{
			name:           "invalid values",
			debug:          "not-bool",
			maxUploadBytes: "not-int",
		},
		{
			name:           "nonpositive upload limit",
			debug:          "not-bool",
			maxUploadBytes: "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("DSN", "host=localhost user=test dbname=syncra_test")
			t.Setenv("MISTRAL_API_KEY", "test-key")
			t.Setenv("BETTER_AUTH_SECRET", "vN8qR4zT7mK2pL9xC5wY3hD6sF0aJ1uB")
			t.Setenv("APP_PRIVATE_KEY", "test-private-key-32-byte-material")
			t.Setenv("DEBUG", tt.debug)
			t.Setenv("MAX_UPLOAD_BYTES", tt.maxUploadBytes)

			cfg, err := Load()
			if err != nil {
				t.Fatalf("Load() error = %v", err)
			}
			if cfg.Debug {
				t.Fatal("Debug should fall back to false")
			}
			if cfg.MaxUploadBytes != 20<<20 {
				t.Fatalf("MaxUploadBytes = %d", cfg.MaxUploadBytes)
			}
		})
	}
}

func TestLoadFallsBackForInvalidOCRExecutorValues(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value string
	}{
		{
			name:  "invalid workers",
			key:   "OCR_EXECUTOR_WORKERS",
			value: "not-int",
		},
		{
			name:  "nonpositive workers",
			key:   "OCR_EXECUTOR_WORKERS",
			value: "0",
		},
		{
			name:  "invalid queue buffer",
			key:   "OCR_EXECUTOR_QUEUE_BUFFER",
			value: "not-int",
		},
		{
			name:  "nonpositive queue buffer",
			key:   "OCR_EXECUTOR_QUEUE_BUFFER",
			value: "0",
		},
		{
			name:  "invalid poll interval",
			key:   "OCR_EXECUTOR_POLL_INTERVAL_SECONDS",
			value: "not-int",
		},
		{
			name:  "nonpositive poll interval",
			key:   "OCR_EXECUTOR_POLL_INTERVAL_SECONDS",
			value: "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("DSN", "host=localhost user=test dbname=syncra_test")
			t.Setenv("MISTRAL_API_KEY", "test-key")
			t.Setenv("BETTER_AUTH_SECRET", "vN8qR4zT7mK2pL9xC5wY3hD6sF0aJ1uB")
			t.Setenv("APP_PRIVATE_KEY", "test-private-key-32-byte-material")
			t.Setenv(tt.key, tt.value)

			cfg, err := Load()
			if err != nil {
				t.Fatalf("Load() error = %v", err)
			}
			if cfg.OCRExecutorWorkers != 2 {
				t.Fatalf("OCRExecutorWorkers = %d", cfg.OCRExecutorWorkers)
			}
			if cfg.OCRExecutorPollIntervalSeconds != 10 {
				t.Fatalf("OCRExecutorPollIntervalSeconds = %d", cfg.OCRExecutorPollIntervalSeconds)
			}
			if cfg.OCRExecutorQueueBuffer != 8 {
				t.Fatalf("OCRExecutorQueueBuffer = %d", cfg.OCRExecutorQueueBuffer)
			}
		})
	}
}

func TestLoadRejectsInvalidOnboardingCredits(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{name: "invalid", value: "not-int"},
		{name: "zero", value: "0"},
		{name: "negative", value: "-1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("DSN", "host=localhost user=test dbname=syncra_test")
			t.Setenv("MISTRAL_API_KEY", "test-key")
			t.Setenv("BETTER_AUTH_SECRET", "vN8qR4zT7mK2pL9xC5wY3hD6sF0aJ1uB")
			t.Setenv("APP_PRIVATE_KEY", "test-private-key-32-byte-material")
			t.Setenv("ONBOARDING_CREDITS", tt.value)

			_, err := Load()
			if err == nil {
				t.Fatal("Load() expected ONBOARDING_CREDITS error")
			}
			if err.Error() != "ONBOARDING_CREDITS must be a positive integer" {
				t.Fatalf("error = %q", err.Error())
			}
		})
	}
}

func TestLoadRequiresAppPrivateKey(t *testing.T) {
	t.Setenv("DSN", "host=localhost user=test dbname=syncra_test")
	t.Setenv("MISTRAL_API_KEY", "test-key")
	t.Setenv("BETTER_AUTH_SECRET", "vN8qR4zT7mK2pL9xC5wY3hD6sF0aJ1uB")
	t.Setenv("APP_PRIVATE_KEY", "")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() expected APP_PRIVATE_KEY error")
	}
	if err.Error() != "APP_PRIVATE_KEY is required" {
		t.Fatalf("error = %q", err.Error())
	}
}

func TestLoadRejectsShortAppPrivateKey(t *testing.T) {
	t.Setenv("DSN", "host=localhost user=test dbname=syncra_test")
	t.Setenv("MISTRAL_API_KEY", "test-key")
	t.Setenv("BETTER_AUTH_SECRET", "vN8qR4zT7mK2pL9xC5wY3hD6sF0aJ1uB")
	t.Setenv("APP_PRIVATE_KEY", "short-private-key")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() expected APP_PRIVATE_KEY length error")
	}
	if err.Error() != "APP_PRIVATE_KEY must be at least 32 characters" {
		t.Fatalf("error = %q", err.Error())
	}
}

func TestLoadRequiresBetterAuthSecret(t *testing.T) {
	t.Setenv("DSN", "host=localhost user=test dbname=syncra_test port=5432 sslmode=disable")
	t.Setenv("MISTRAL_API_KEY", "test-key")
	t.Setenv("BETTER_AUTH_SECRET", "")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() expected BETTER_AUTH_SECRET error")
	}
}

func TestLoadRejectsWeakBetterAuthSecret(t *testing.T) {
	tests := []struct {
		name   string
		secret string
	}{
		{
			name:   "too short",
			secret: "short-value",
		},
		{
			name:   "repeated character",
			secret: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		},
		{
			name:   "too few distinct characters",
			secret: "aaaabbbbccccddddeeeeffffggggaaaa",
		},
		{
			name:   "placeholder",
			secret: "change-me-before-production-123456",
		},
		{
			name:   "common secret label",
			secret: "prod-secret-value-that-is-not-random",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("DSN", "host=localhost user=test dbname=syncra_test port=5432 sslmode=disable")
			t.Setenv("MISTRAL_API_KEY", "test-key")
			t.Setenv("BETTER_AUTH_SECRET", tt.secret)

			_, err := Load()
			if err == nil {
				t.Fatal("Load() expected weak BETTER_AUTH_SECRET error")
			}
		})
	}
}

func TestLoadRejectsAuthTTLsUnsafeForDurationConversion(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value string
	}{
		{
			name:  "session ttl exceeds duration seconds",
			key:   "AUTH_SESSION_TTL_SECONDS",
			value: "9223372037",
		},
		{
			name:  "verification ttl exceeds duration seconds",
			key:   "AUTH_VERIFICATION_TTL_SECONDS",
			value: "9223372037",
		},
		{
			name:  "session ttl exceeds int64",
			key:   "AUTH_SESSION_TTL_SECONDS",
			value: "9223372036854775808",
		},
		{
			name:  "verification ttl exceeds int64",
			key:   "AUTH_VERIFICATION_TTL_SECONDS",
			value: "9223372036854775808",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("DSN", "host=localhost user=test dbname=syncra_test port=5432 sslmode=disable")
			t.Setenv("MISTRAL_API_KEY", "test-key")
			t.Setenv("BETTER_AUTH_SECRET", "vN8qR4zT7mK2pL9xC5wY3hD6sF0aJ1uB")
			t.Setenv(tt.key, tt.value)

			_, err := Load()
			if err == nil {
				t.Fatalf("Load() expected %s overflow error", tt.key)
			}
		})
	}
}

func TestLoadRequiresMistralAPIKey(t *testing.T) {
	t.Setenv("DSN", "host=localhost user=test dbname=syncra_test")
	t.Setenv("MISTRAL_API_KEY", "")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() expected MISTRAL_API_KEY error")
	}
	if err.Error() != "MISTRAL_API_KEY is required" {
		t.Fatalf("error = %q", err.Error())
	}
}

func TestLoadDatabaseRequiresOnlyDSN(t *testing.T) {
	t.Setenv("DSN", "  host=localhost user=test dbname=syncra_test  ")
	t.Setenv("MISTRAL_API_KEY", "")
	t.Setenv("BETTER_AUTH_SECRET", "")

	cfg, err := LoadDatabase()
	if err != nil {
		t.Fatalf("LoadDatabase() error = %v", err)
	}
	if cfg.DSN != "host=localhost user=test dbname=syncra_test" {
		t.Fatalf("DSN = %q", cfg.DSN)
	}
}

func TestLoadDatabaseRequiresDSN(t *testing.T) {
	t.Setenv("DSN", " \t\n ")

	_, err := LoadDatabase()
	if err == nil {
		t.Fatal("LoadDatabase() expected DSN error")
	}
	if err.Error() != "DSN is required" {
		t.Fatalf("error = %q", err.Error())
	}
}

func TestLoadMigrationRequiresAtlasDatabaseURL(t *testing.T) {
	t.Setenv("ATLAS_DATABASE_URL", " \t\n ")

	_, err := LoadMigration()
	if err == nil {
		t.Fatal("LoadMigration() expected ATLAS_DATABASE_URL error")
	}
	if err.Error() != "ATLAS_DATABASE_URL is required" {
		t.Fatalf("error = %q", err.Error())
	}
}

func TestLoadMigrationTrimsAtlasDatabaseURL(t *testing.T) {
	t.Setenv("ATLAS_DATABASE_URL", "  postgres://postgres:pass@localhost:5432/syncra_dev?search_path=public&sslmode=disable  ")

	cfg, err := LoadMigration()
	if err != nil {
		t.Fatalf("LoadMigration() error = %v", err)
	}
	if cfg.AtlasDatabaseURL != "postgres://postgres:pass@localhost:5432/syncra_dev?search_path=public&sslmode=disable" {
		t.Fatalf("AtlasDatabaseURL = %q", cfg.AtlasDatabaseURL)
	}
}

func unsetenv(t *testing.T, key string) {
	t.Helper()

	oldValue, wasSet := os.LookupEnv(key)
	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("Unsetenv(%q) error = %v", key, err)
	}
	t.Cleanup(func() {
		if wasSet {
			if err := os.Setenv(key, oldValue); err != nil {
				t.Fatalf("Setenv(%q) error = %v", key, err)
			}
			return
		}
		if err := os.Unsetenv(key); err != nil {
			t.Fatalf("Unsetenv(%q) error = %v", key, err)
		}
	})
}
