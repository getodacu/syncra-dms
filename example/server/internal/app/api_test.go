package app

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"ai.ro/syncra/internal/api"
	"ai.ro/syncra/internal/auth"
	"ai.ro/syncra/internal/config"
	"ai.ro/syncra/internal/ocr"
	"ai.ro/syncra/internal/testsupport"
)

type staticServer struct {
	err error
}

func (s staticServer) ListenAndServe() error {
	return s.err
}

func TestRunAPIDoesNotMigrateDatabase(t *testing.T) {
	db := testsupport.OpenPostgresTx(t)
	modelsToMigrate := []any{
		&ocr.ExtractionSchema{},
		&ocr.OCRDocument{},
		&auth.User{},
		&auth.AuthAccount{},
		&auth.Session{},
		&auth.Verification{},
		&auth.AdminImpersonationEvent{},
	}
	for _, model := range modelsToMigrate {
		if err := db.Migrator().DropTable(model); err != nil {
			t.Fatalf("drop table for %T: %v", model, err)
		}
	}

	listenErr := errors.New("server stopped")
	openCalled := false
	serverCreated := false
	err := runAPI(config.Config{
		DSN:            "postgres://example",
		ServerHostPort: ":0",
	}, apiDeps{
		openPostgres: func(dsn string) (*gorm.DB, error) {
			openCalled = true
			if dsn != "postgres://example" {
				t.Fatalf("openPostgres DSN = %q, want postgres://example", dsn)
			}
			return db, nil
		},
		closeDB: func(*gorm.DB) {},
		runOCRExecutor: func(context.Context, ocr.ExecutorConfig) func() {
			return func() {}
		},
		newRouter: func(handler *api.Handler) http.Handler {
			if handler == nil {
				t.Fatal("api handler is nil")
			}
			return http.NewServeMux()
		},
		newServer: func(addr string, handler http.Handler) listenAndServeServer {
			serverCreated = true
			if addr != ":0" {
				t.Fatalf("server addr = %q, want :0", addr)
			}
			if handler == nil {
				t.Fatal("server handler is nil")
			}
			return staticServer{err: listenErr}
		},
	})

	if !errors.Is(err, listenErr) {
		t.Fatalf("runAPI() error = %v, want %v", err, listenErr)
	}
	if !openCalled {
		t.Fatal("openPostgres was not called")
	}
	if !serverCreated {
		t.Fatal("server was not created")
	}
	for _, model := range modelsToMigrate {
		if db.Migrator().HasTable(model) {
			t.Fatalf("API startup migrated table for %T", model)
		}
	}
}

func TestRunAPIPassesStorageDirToHandler(t *testing.T) {
	listenErr := errors.New("server stopped")
	wantDir := "/var/lib/syncra"
	gotDir := ""

	err := runAPI(config.Config{
		DSN:            "postgres://example",
		ServerHostPort: ":0",
		StorageDir:     wantDir,
	}, apiDeps{
		openPostgres: func(dsn string) (*gorm.DB, error) {
			if dsn != "postgres://example" {
				t.Fatalf("openPostgres DSN = %q, want postgres://example", dsn)
			}
			return &gorm.DB{}, nil
		},
		closeDB: func(*gorm.DB) {},
		runOCRExecutor: func(context.Context, ocr.ExecutorConfig) func() {
			return func() {}
		},
		newRouter: func(handler *api.Handler) http.Handler {
			if handler == nil {
				t.Fatal("api handler is nil")
			}
			gotDir = handler.StorageDir
			return http.NewServeMux()
		},
		newServer: func(addr string, handler http.Handler) listenAndServeServer {
			if addr != ":0" {
				t.Fatalf("server addr = %q, want :0", addr)
			}
			if handler == nil {
				t.Fatal("server handler is nil")
			}
			return staticServer{err: listenErr}
		},
	})

	if !errors.Is(err, listenErr) {
		t.Fatalf("runAPI() error = %v, want %v", err, listenErr)
	}
	if gotDir != wantDir {
		t.Fatalf("handler StorageDir = %q, want %q", gotDir, wantDir)
	}
}

func TestRunAPIPassesInvoicePDFConfigToHandler(t *testing.T) {
	listenErr := errors.New("server stopped")
	wantGotenbergURL := "http://gotenberg.example/forms/chromium/convert/html"
	wantStorageDir := "/var/lib/syncra"
	var gotGotenbergURL string
	var gotStorageDir string

	err := runAPI(config.Config{
		DSN:             "postgres://example",
		ServerHostPort:  ":0",
		GotenbergAPIURL: wantGotenbergURL,
		StorageDir:      wantStorageDir,
	}, apiDeps{
		openPostgres: func(dsn string) (*gorm.DB, error) {
			if dsn != "postgres://example" {
				t.Fatalf("openPostgres DSN = %q, want postgres://example", dsn)
			}
			return &gorm.DB{}, nil
		},
		closeDB: func(*gorm.DB) {},
		runOCRExecutor: func(context.Context, ocr.ExecutorConfig) func() {
			return func() {}
		},
		newRouter: func(handler *api.Handler) http.Handler {
			if handler == nil {
				t.Fatal("api handler is nil")
			}
			gotGotenbergURL = handler.GotenbergAPIURL
			gotStorageDir = handler.StorageDir
			return http.NewServeMux()
		},
		newServer: func(addr string, handler http.Handler) listenAndServeServer {
			if addr != ":0" {
				t.Fatalf("server addr = %q, want :0", addr)
			}
			if handler == nil {
				t.Fatal("server handler is nil")
			}
			return staticServer{err: listenErr}
		},
	})

	if !errors.Is(err, listenErr) {
		t.Fatalf("runAPI() error = %v, want %v", err, listenErr)
	}
	if gotGotenbergURL != wantGotenbergURL {
		t.Fatalf("handler GotenbergAPIURL = %q, want %q", gotGotenbergURL, wantGotenbergURL)
	}
	if gotStorageDir != wantStorageDir {
		t.Fatalf("handler StorageDir = %q, want %q", gotStorageDir, wantStorageDir)
	}
}

func TestRunAPIStartsOCRExecutor(t *testing.T) {
	listenErr := errors.New("server stopped")
	db := &gorm.DB{}
	var gotCtx context.Context
	var gotExecutorCfg ocr.ExecutorConfig
	var gotHandler *api.Handler
	executorStarted := false
	executorExited := false
	closeCalled := false

	err := runAPI(config.Config{
		DSN:                            "postgres://executor",
		ServerHostPort:                 ":0",
		MistralAPIKey:                  "test-key",
		MistralBaseURL:                 "https://mistral.example",
		MistralOCRModel:                "test-ocr-model",
		AppPrivateKey:                  "test-private-key-32-byte-material",
		InternalAPIToken:               "trusted-internal-token",
		GoogleClientID:                 "google-client-id",
		GoogleClientSecret:             "google-client-secret",
		GitHubClientID:                 "github-client-id",
		GitHubClientSecret:             "github-client-secret",
		OnboardingCredits:              250,
		OCRExecutorWorkers:             5,
		OCRExecutorQueueBuffer:         11,
		OCRExecutorPollIntervalSeconds: 7,
	}, apiDeps{
		openPostgres: func(dsn string) (*gorm.DB, error) {
			if dsn != "postgres://executor" {
				t.Fatalf("openPostgres DSN = %q, want postgres://executor", dsn)
			}
			return db, nil
		},
		closeDB: func(*gorm.DB) {
			closeCalled = true
			if !executorExited {
				t.Fatal("closeDB called before executor wait completed")
			}
		},
		runOCRExecutor: func(ctx context.Context, cfg ocr.ExecutorConfig) func() {
			executorStarted = true
			gotCtx = ctx
			gotExecutorCfg = cfg
			if ctx.Err() != nil {
				t.Fatalf("executor context already canceled: %v", ctx.Err())
			}
			return func() {
				if ctx.Err() != context.Canceled {
					t.Fatalf("executor wait called before lifecycle context cancellation: %v", ctx.Err())
				}
				executorExited = true
			}
		},
		newRouter: func(handler *api.Handler) http.Handler {
			gotHandler = handler
			return http.NewServeMux()
		},
		newServer: func(addr string, handler http.Handler) listenAndServeServer {
			if !executorStarted {
				t.Fatal("executor was not started before server creation")
			}
			return staticServer{err: listenErr}
		},
	})

	if !errors.Is(err, listenErr) {
		t.Fatalf("runAPI() error = %v, want %v", err, listenErr)
	}
	if !executorStarted {
		t.Fatal("OCR executor was not started")
	}
	if !executorExited {
		t.Fatal("executor wait was not called")
	}
	if !closeCalled {
		t.Fatal("closeDB was not called")
	}
	if gotCtx == nil {
		t.Fatal("executor context was nil")
	}
	if gotCtx.Err() != context.Canceled {
		t.Fatalf("executor context error after runAPI return = %v, want context canceled", gotCtx.Err())
	}
	if gotExecutorCfg.DB != db {
		t.Fatalf("executor DB = %#v, want opened DB", gotExecutorCfg.DB)
	}
	if gotExecutorCfg.DSN != "postgres://executor" {
		t.Fatalf("executor DSN = %q, want postgres://executor", gotExecutorCfg.DSN)
	}
	if gotExecutorCfg.Processor == nil {
		t.Fatal("executor processor is nil")
	}
	if gotExecutorCfg.WorkerCount != 5 {
		t.Fatalf("executor worker count = %d, want 5", gotExecutorCfg.WorkerCount)
	}
	if gotExecutorCfg.QueueBuffer != 11 {
		t.Fatalf("executor queue buffer = %d, want 11", gotExecutorCfg.QueueBuffer)
	}
	if gotExecutorCfg.PollInterval != 7*time.Second {
		t.Fatalf("executor poll interval = %s, want 7s", gotExecutorCfg.PollInterval)
	}
	if gotExecutorCfg.Logger == nil {
		t.Fatal("executor Logger is nil")
	}
	if gotExecutorCfg.WebhookDispatcher == nil {
		t.Fatal("executor WebhookDispatcher is nil")
	}
	if gotHandler == nil {
		t.Fatal("api handler was nil")
	}
	if gotHandler.Logger == nil {
		t.Fatal("handler Logger is nil")
	}
	if gotHandler.OCR == nil {
		t.Fatal("handler OCR processor is nil")
	}
	if gotHandler.InternalAPIToken != "trusted-internal-token" {
		t.Fatalf("handler InternalAPIToken = %q", gotHandler.InternalAPIToken)
	}
	if gotHandler.AppPrivateKey != "test-private-key-32-byte-material" {
		t.Fatalf("handler AppPrivateKey = %q", gotHandler.AppPrivateKey)
	}
	if gotHandler.GoogleClientID != "google-client-id" {
		t.Fatalf("handler GoogleClientID = %q", gotHandler.GoogleClientID)
	}
	if gotHandler.GoogleClientSecret != "google-client-secret" {
		t.Fatalf("handler GoogleClientSecret = %q", gotHandler.GoogleClientSecret)
	}
	if gotHandler.GitHubClientID != "github-client-id" {
		t.Fatalf("handler GitHubClientID = %q", gotHandler.GitHubClientID)
	}
	if gotHandler.GitHubClientSecret != "github-client-secret" {
		t.Fatalf("handler GitHubClientSecret = %q", gotHandler.GitHubClientSecret)
	}
	if gotHandler.OnboardingCredits != 250 {
		t.Fatalf("handler OnboardingCredits = %d, want 250", gotHandler.OnboardingCredits)
	}
	if reflect.ValueOf(gotHandler.OCR).Pointer() != reflect.ValueOf(gotExecutorCfg.Processor).Pointer() {
		t.Fatal("handler and executor did not receive the same OCR processor function")
	}
}

func TestApplyRuntimeConfigUsesDebugSetting(t *testing.T) {
	originalMode := gin.Mode()
	t.Cleanup(func() {
		gin.SetMode(originalMode)
	})
	gin.SetMode(gin.DebugMode)

	ApplyRuntimeConfig(config.Config{Debug: false})

	if got := gin.Mode(); got != gin.ReleaseMode {
		t.Fatalf("gin mode = %q, want %q", got, gin.ReleaseMode)
	}

	gin.SetMode(gin.ReleaseMode)

	ApplyRuntimeConfig(config.Config{Debug: true})

	if got := gin.Mode(); got != gin.DebugMode {
		t.Fatalf("gin mode = %q, want %q", got, gin.DebugMode)
	}
}

func TestNewHTTPServerConfiguresTimeouts(t *testing.T) {
	handler := http.NewServeMux()

	server := NewHTTPServer(":0", handler)

	if server.Addr != ":0" {
		t.Fatalf("Addr = %q, want :0", server.Addr)
	}
	if server.Handler != handler {
		t.Fatalf("Handler = %#v, want provided handler", server.Handler)
	}
	if server.ReadHeaderTimeout != 5*time.Second {
		t.Fatalf("ReadHeaderTimeout = %s, want 5s", server.ReadHeaderTimeout)
	}
	if server.ReadTimeout != 30*time.Second {
		t.Fatalf("ReadTimeout = %s, want 30s", server.ReadTimeout)
	}
	if server.WriteTimeout != 180*time.Second {
		t.Fatalf("WriteTimeout = %s, want 180s", server.WriteTimeout)
	}
	if server.IdleTimeout != 60*time.Second {
		t.Fatalf("IdleTimeout = %s, want 60s", server.IdleTimeout)
	}
}
