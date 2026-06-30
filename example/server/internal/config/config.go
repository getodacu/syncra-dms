package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

const (
	minBetterAuthSecretLength = 32
	minBetterAuthSecretRunes  = 8
	minAppPrivateKeyLength    = 32
	defaultOnboardingCredits  = 100
	maxAuthTTLSeconds         = int64(1<<63-1) / int64(time.Second)
)

var weakBetterAuthSecretSubstrings = []string{
	"change-me",
	"changeme",
	"placeholder",
	"secret",
	"password",
	"example",
	"test-secret",
}

type Config struct {
	ServerHostPort                 string
	Debug                          bool
	DSN                            string
	SwaggerHost                    string
	SwaggerSchemes                 []string
	MistralAPIKey                  string
	MistralBaseURL                 string
	MistralOCRModel                string
	MaxUploadBytes                 int64
	StorageDir                     string
	OCRExecutorWorkers             int
	OCRExecutorPollIntervalSeconds int64
	OCRExecutorQueueBuffer         int
	GotenbergAPIURL                string
	BetterAuthSecret               string
	AppPrivateKey                  string
	AuthDeliveryToken              string
	InternalAPIToken               string
	AuthSessionTTL                 int64
	AuthVerificationTTL            int64
	AuthCookieSecure               bool
	OnboardingCredits              int
	GoogleClientID                 string
	GoogleClientSecret             string
	GitHubClientID                 string
	GitHubClientSecret             string
}

type DatabaseConfig struct {
	DSN string
}

type MigrationConfig struct {
	AtlasDatabaseURL string
}

func LoadDatabase() (DatabaseConfig, error) {
	_ = godotenv.Load()

	dsn := strings.TrimSpace(os.Getenv("DSN"))
	if dsn == "" {
		return DatabaseConfig{}, errors.New("DSN is required")
	}
	return DatabaseConfig{DSN: dsn}, nil
}

func LoadMigration() (MigrationConfig, error) {
	_ = godotenv.Load()

	atlasDatabaseURL := strings.TrimSpace(os.Getenv("ATLAS_DATABASE_URL"))
	if atlasDatabaseURL == "" {
		return MigrationConfig{}, errors.New("ATLAS_DATABASE_URL is required")
	}
	return MigrationConfig{AtlasDatabaseURL: atlasDatabaseURL}, nil
}

func Load() (Config, error) {
	dbCfg, err := LoadDatabase()
	if err != nil {
		return Config{}, err
	}

	authSessionTTL, err := getenvAuthTTLSeconds("AUTH_SESSION_TTL_SECONDS", 604800)
	if err != nil {
		return Config{}, err
	}
	authVerificationTTL, err := getenvAuthTTLSeconds("AUTH_VERIFICATION_TTL_SECONDS", 300)
	if err != nil {
		return Config{}, err
	}
	onboardingCredits, err := getenvPositiveInt("ONBOARDING_CREDITS", defaultOnboardingCredits)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		ServerHostPort:                 getenv("SERVER_HOST_PORT", "localhost:8080"),
		Debug:                          getenvBool("DEBUG", false),
		DSN:                            dbCfg.DSN,
		SwaggerHost:                    getenv("SWAGGER_HOST", "localhost:8080"),
		SwaggerSchemes:                 getenvCSV("SWAGGER_SCHEMES", "http"),
		MistralAPIKey:                  strings.TrimSpace(os.Getenv("MISTRAL_API_KEY")),
		MistralBaseURL:                 strings.TrimRight(getenv("MISTRAL_API_BASE_URL", "https://api.mistral.ai"), "/"),
		MistralOCRModel:                getenv("MISTRAL_OCR_MODEL", "mistral-ocr-latest"),
		MaxUploadBytes:                 getenvInt64("MAX_UPLOAD_BYTES", 20<<20),
		StorageDir:                     strings.TrimSpace(os.Getenv("STORAGE_DIR")),
		OCRExecutorWorkers:             getenvInt("OCR_EXECUTOR_WORKERS", 2),
		OCRExecutorPollIntervalSeconds: getenvInt64("OCR_EXECUTOR_POLL_INTERVAL_SECONDS", 10),
		OCRExecutorQueueBuffer:         getenvInt("OCR_EXECUTOR_QUEUE_BUFFER", 8),
		GotenbergAPIURL:                normalizeGotenbergAPIURL(os.Getenv("GOTENBERG_API_URL")),
		BetterAuthSecret:               strings.TrimSpace(os.Getenv("BETTER_AUTH_SECRET")),
		AppPrivateKey:                  strings.TrimSpace(os.Getenv("APP_PRIVATE_KEY")),
		AuthDeliveryToken:              strings.TrimSpace(os.Getenv("AUTH_DELIVERY_TOKEN")),
		InternalAPIToken:               strings.TrimSpace(os.Getenv("SYNCRA_INTERNAL_API_TOKEN")),
		AuthSessionTTL:                 authSessionTTL,
		AuthVerificationTTL:            authVerificationTTL,
		AuthCookieSecure:               getenvBool("AUTH_COOKIE_SECURE", false),
		OnboardingCredits:              onboardingCredits,
		GoogleClientID:                 strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_ID")),
		GoogleClientSecret:             strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_SECRET")),
		GitHubClientID:                 strings.TrimSpace(os.Getenv("GITHUB_CLIENT_ID")),
		GitHubClientSecret:             strings.TrimSpace(os.Getenv("GITHUB_CLIENT_SECRET")),
	}

	if cfg.MistralAPIKey == "" {
		return Config{}, errors.New("MISTRAL_API_KEY is required")
	}
	if cfg.BetterAuthSecret == "" {
		return Config{}, errors.New("BETTER_AUTH_SECRET is required")
	}
	if err := validateBetterAuthSecret(cfg.BetterAuthSecret); err != nil {
		return Config{}, err
	}
	if cfg.AppPrivateKey == "" {
		return Config{}, errors.New("APP_PRIVATE_KEY is required")
	}
	if err := validateAppPrivateKey(cfg.AppPrivateKey); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func getenv(key string, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func getenvBool(key string, fallback bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func getenvInt(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}

func getenvInt64(key string, fallback int64) int64 {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}

func getenvPositiveInt(key string, fallback int) (int, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return 0, fmt.Errorf("%s must be a positive integer", key)
	}
	return parsed, nil
}

func getenvCSV(key string, fallback string) []string {
	raw := getenv(key, fallback)
	parts := strings.Split(raw, ",")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			values = append(values, part)
		}
	}
	return values
}

func normalizeGotenbergAPIURL(raw string) string {
	value := strings.TrimSpace(raw)
	if value == "" {
		return ""
	}
	if !strings.Contains(value, "://") {
		value = "http://" + value
	}
	parsed, err := url.Parse(value)
	if err != nil || parsed.Host == "" {
		return strings.TrimRight(value, "/")
	}
	parsed.Path = "/forms/chromium/convert/html"
	parsed.RawQuery = ""
	parsed.Fragment = ""
	return parsed.String()
}

func validateAuthTTL(key string, seconds int64) error {
	if seconds > maxAuthTTLSeconds {
		return fmt.Errorf("%s must be <= %d", key, maxAuthTTLSeconds)
	}
	return nil
}

func validateBetterAuthSecret(secret string) error {
	if len(secret) < minBetterAuthSecretLength {
		return fmt.Errorf("BETTER_AUTH_SECRET must be at least %d characters", minBetterAuthSecretLength)
	}

	distinct := make(map[rune]struct{}, len(secret))
	for _, r := range secret {
		distinct[r] = struct{}{}
	}
	if len(distinct) < minBetterAuthSecretRunes {
		return fmt.Errorf("BETTER_AUTH_SECRET must contain at least %d distinct characters", minBetterAuthSecretRunes)
	}

	lower := strings.ToLower(strings.TrimSpace(secret))
	for _, weak := range weakBetterAuthSecretSubstrings {
		if strings.Contains(lower, weak) {
			return errors.New("BETTER_AUTH_SECRET must not contain placeholder or common secret values")
		}
	}

	return nil
}

func validateAppPrivateKey(privateKey string) error {
	if len(privateKey) < minAppPrivateKeyLength {
		return fmt.Errorf("APP_PRIVATE_KEY must be at least %d characters", minAppPrivateKeyLength)
	}
	return nil
}

func getenvAuthTTLSeconds(key string, fallback int64) (int64, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback, nil
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		if numErr, ok := err.(*strconv.NumError); ok && numErr.Err == strconv.ErrRange {
			return 0, fmt.Errorf("%s must be <= %d", key, maxAuthTTLSeconds)
		}
		return fallback, nil
	}
	if parsed <= 0 {
		return fallback, nil
	}
	if err := validateAuthTTL(key, parsed); err != nil {
		return 0, err
	}
	return parsed, nil
}
