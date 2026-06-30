package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

const (
	defaultServerHostPort = "localhost:8080"
	defaultVersion        = "dev"
	appDatabaseName       = "syncra_dms"
	testDatabaseName      = "syncra_dms_dev"
)

type Config struct {
	ServerHostPort             string
	Debug                      bool
	DSN                        string
	DSNDev                     string
	AtlasDatabaseURL           string
	AtlasDevDatabaseURL        string
	InternalAPIToken           string
	BetterAuthSecret           string
	AuthDeliveryToken          string
	AuthSessionTTLSeconds      int64
	AuthVerificationTTLSeconds int64
	AuthCookieSecure           bool
	GoogleClientID             string
	GoogleClientSecret         string
	GitHubClientID             string
	GitHubClientSecret         string
	Version                    string
}

func Load() (Config, error) {
	_ = godotenv.Load()

	cfg := Config{
		ServerHostPort:             getenv("SERVER_HOST_PORT", defaultServerHostPort),
		Debug:                      getenvBool("DEBUG", false),
		DSN:                        strings.TrimSpace(os.Getenv("DSN")),
		DSNDev:                     strings.TrimSpace(os.Getenv("DSN_DEV")),
		AtlasDatabaseURL:           strings.TrimSpace(os.Getenv("ATLAS_DATABASE_URL")),
		AtlasDevDatabaseURL:        strings.TrimSpace(os.Getenv("ATLAS_DEV_DATABASE_URL")),
		InternalAPIToken:           strings.TrimSpace(os.Getenv("SYNCRA_INTERNAL_API_TOKEN")),
		BetterAuthSecret:           strings.TrimSpace(os.Getenv("BETTER_AUTH_SECRET")),
		AuthDeliveryToken:          strings.TrimSpace(os.Getenv("AUTH_DELIVERY_TOKEN")),
		AuthSessionTTLSeconds:      getenvInt64("AUTH_SESSION_TTL_SECONDS", 604800),
		AuthVerificationTTLSeconds: getenvInt64("AUTH_VERIFICATION_TTL_SECONDS", 300),
		AuthCookieSecure:           getenvBool("AUTH_COOKIE_SECURE", false),
		GoogleClientID:             strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_ID")),
		GoogleClientSecret:         strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_SECRET")),
		GitHubClientID:             strings.TrimSpace(os.Getenv("GITHUB_CLIENT_ID")),
		GitHubClientSecret:         strings.TrimSpace(os.Getenv("GITHUB_CLIENT_SECRET")),
		Version:                    getenv("APP_VERSION", defaultVersion),
	}

	if cfg.DSN == "" {
		return Config{}, errors.New("DSN is required")
	}
	if cfg.DSNDev == "" {
		return Config{}, errors.New("DSN_DEV is required")
	}
	if err := requireDatabaseName("DSN", cfg.DSN, appDatabaseName); err != nil {
		return Config{}, err
	}
	if err := requireDatabaseName("DSN_DEV", cfg.DSNDev, testDatabaseName); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func DatabaseNameFromDSN(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", errors.New("database DSN is empty")
	}

	if strings.HasPrefix(raw, "postgres://") || strings.HasPrefix(raw, "postgresql://") {
		parsed, err := url.Parse(raw)
		if err != nil {
			return "", err
		}
		name := strings.Trim(strings.TrimSpace(parsed.Path), "/")
		if name == "" {
			return "", errors.New("database name is missing")
		}
		return name, nil
	}

	for _, field := range strings.Fields(raw) {
		key, value, found := strings.Cut(field, "=")
		if !found {
			continue
		}
		if key == "dbname" {
			name := strings.Trim(value, `"'`)
			if name == "" {
				return "", errors.New("database name is missing")
			}
			return name, nil
		}
	}
	return "", errors.New("dbname is required")
}

func requireDatabaseName(label string, raw string, want string) error {
	got, err := DatabaseNameFromDSN(raw)
	if err != nil {
		return fmt.Errorf("%s: %w", label, err)
	}
	if got != want {
		return fmt.Errorf("%s must target %s, got %s", label, want, got)
	}
	return nil
}

func getenv(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func getenvBool(key string, fallback bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value == "1" || strings.EqualFold(value, "true") || strings.EqualFold(value, "yes")
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
