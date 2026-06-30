package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"ai.ro/syncra/internal/api"
	"ai.ro/syncra/internal/config"
	"ai.ro/syncra/internal/database"
	"ai.ro/syncra/internal/logging"
	"ai.ro/syncra/internal/ocr"
	"ai.ro/syncra/internal/webhooks"
)

func RunAPI(cfg config.Config) error {
	logger := logging.ConfigureDefault(cfg.Debug, os.Stdout).With(
		"service", "syncra",
	)
	return runAPI(cfg, realAPIDeps(logger))
}

type apiDeps struct {
	logger         *slog.Logger
	openPostgres   func(string) (*gorm.DB, error)
	closeDB        func(*gorm.DB)
	runOCRExecutor func(context.Context, ocr.ExecutorConfig) func()
	newRouter      func(*api.Handler) http.Handler
	newServer      func(string, http.Handler) listenAndServeServer
}

type listenAndServeServer interface {
	ListenAndServe() error
}

func realAPIDeps(logger *slog.Logger) apiDeps {
	if logger == nil {
		logger = logging.Nop()
	}
	return apiDeps{
		logger:       logger,
		openPostgres: database.OpenPostgres,
		closeDB:      closeGormDB,
		runOCRExecutor: func(ctx context.Context, cfg ocr.ExecutorConfig) func() {
			executor := ocr.NewExecutor(cfg)
			done := make(chan struct{})
			go func() {
				defer close(done)
				if err := executor.Run(ctx); err != nil && ctx.Err() == nil {
					logger.Error("ocr.executor_stopped_unexpectedly", "error", err)
				}
			}()
			return func() {
				<-done
			}
		},
		newRouter: func(handler *api.Handler) http.Handler { return api.NewRouter(handler) },
		newServer: func(addr string, handler http.Handler) listenAndServeServer {
			return NewHTTPServer(addr, handler)
		},
	}
}

func runAPI(cfg config.Config, deps apiDeps) error {
	ApplyRuntimeConfig(cfg)
	logger := deps.logger
	if logger == nil {
		logger = logging.Nop()
	}

	db, err := deps.openPostgres(cfg.DSN)
	if err != nil {
		logger.Error("api.database_open_failed", "error", err)
		return err
	}

	lifecycleCtx, cancelLifecycle := context.WithCancel(context.Background())

	processor := ocr.NewMistralProcessor(ocr.MistralConfig{
		APIKey:  cfg.MistralAPIKey,
		BaseURL: cfg.MistralBaseURL,
		Model:   cfg.MistralOCRModel,
	})
	webhookDispatcher := webhooks.NewDispatcher(webhooks.DispatcherConfig{
		DB:         db,
		PrivateKey: cfg.AppPrivateKey,
	})

	handler := &api.Handler{
		DB:                  db,
		OCR:                 processor,
		StorageDir:          cfg.StorageDir,
		MaxUploadBytes:      cfg.MaxUploadBytes,
		GotenbergAPIURL:     cfg.GotenbergAPIURL,
		MistralAPIKey:       cfg.MistralAPIKey,
		MistralBaseURL:      cfg.MistralBaseURL,
		MistralModel:        cfg.MistralOCRModel,
		BetterAuthSecret:    cfg.BetterAuthSecret,
		AppPrivateKey:       cfg.AppPrivateKey,
		AuthDeliveryToken:   cfg.AuthDeliveryToken,
		InternalAPIToken:    cfg.InternalAPIToken,
		AuthSessionTTL:      time.Duration(cfg.AuthSessionTTL) * time.Second,
		AuthVerificationTTL: time.Duration(cfg.AuthVerificationTTL) * time.Second,
		AuthCookieSecure:    cfg.AuthCookieSecure,
		OnboardingCredits:   cfg.OnboardingCredits,
		GoogleClientID:      cfg.GoogleClientID,
		GoogleClientSecret:  cfg.GoogleClientSecret,
		GitHubClientID:      cfg.GitHubClientID,
		GitHubClientSecret:  cfg.GitHubClientSecret,
		SwaggerHost:         cfg.SwaggerHost,
		SwaggerSchemes:      cfg.SwaggerSchemes,
		Logger:              logger,
	}
	waitOCRExecutor := deps.runOCRExecutor(lifecycleCtx, ocr.ExecutorConfig{
		DB:                db,
		DSN:               cfg.DSN,
		Processor:         processor,
		WorkerCount:       cfg.OCRExecutorWorkers,
		QueueBuffer:       cfg.OCRExecutorQueueBuffer,
		PollInterval:      time.Duration(cfg.OCRExecutorPollIntervalSeconds) * time.Second,
		Logger:            logger,
		WebhookDispatcher: webhookDispatcher,
	})
	defer func() {
		cancelLifecycle()
		if waitOCRExecutor != nil {
			waitOCRExecutor()
		}
		deps.closeDB(db)
	}()
	router := deps.newRouter(handler)
	server := deps.newServer(cfg.ServerHostPort, router)
	logger.Info("api.server_starting", "component", "api", "addr", cfg.ServerHostPort, "debug", cfg.Debug)
	return server.ListenAndServe()
}

func closeGormDB(db *gorm.DB) {
	if sqlDB, err := db.DB(); err == nil {
		_ = sqlDB.Close()
	}
}

func ApplyRuntimeConfig(cfg config.Config) {
	if cfg.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
}

func NewHTTPServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      180 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}
