package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"ai.ro/syncra/dms/internal/api"
	"ai.ro/syncra/dms/internal/config"
	"ai.ro/syncra/dms/internal/database"
	"ai.ro/syncra/dms/internal/logging"
)

func RunAPI(cfg config.Config) error {
	logger := logging.ConfigureDefault(cfg.Debug, os.Stdout).With("service", "syncra-dms")
	ApplyRuntimeConfig(cfg)

	db, err := database.OpenPostgres(cfg.DSN)
	if err != nil {
		logger.Error("api.database_open_failed", "error", err)
		return err
	}
	defer closeGormDB(db)

	router := api.NewRouter(api.RouterOptions{
		Version: api.VersionInfo{
			AppName: "Syncra DMS",
			Module:  "ai.ro/syncra/dms",
			Version: cfg.Version,
		},
		Ready: func(ctx context.Context) error {
			return Ping(ctx, db)
		},
		DB:                  db,
		BetterAuthSecret:    cfg.BetterAuthSecret,
		AuthDeliveryToken:   cfg.AuthDeliveryToken,
		InternalAPIToken:    cfg.InternalAPIToken,
		AuthSessionTTL:      time.Duration(cfg.AuthSessionTTLSeconds) * time.Second,
		AuthVerificationTTL: time.Duration(cfg.AuthVerificationTTLSeconds) * time.Second,
		AuthCookieSecure:    cfg.AuthCookieSecure,
		GoogleClientID:      cfg.GoogleClientID,
		GoogleClientSecret:  cfg.GoogleClientSecret,
		GitHubClientID:      cfg.GitHubClientID,
		GitHubClientSecret:  cfg.GitHubClientSecret,
	})
	server := NewHTTPServer(cfg.ServerHostPort, router)
	logger.Info("api.server_starting", "addr", cfg.ServerHostPort, "debug", cfg.Debug)
	return server.ListenAndServe()
}

func ApplyRuntimeConfig(cfg config.Config) {
	if cfg.Debug {
		gin.SetMode(gin.DebugMode)
		return
	}
	gin.SetMode(gin.ReleaseMode)
}

func NewHTTPServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}

func Ping(ctx context.Context, db *gorm.DB) error {
	if db == nil {
		return errors.New("database is not configured")
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

func closeGormDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		slog.Debug("database close skipped", "error", err)
		return
	}
	_ = sqlDB.Close()
}
