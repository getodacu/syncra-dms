package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ai.ro/syncra/dms/internal/config"
)

func TestRunDispatchesAPI(t *testing.T) {
	ctx := context.Background()
	wantCfg := config.Config{
		DSN:            "host=localhost dbname=syncra_dms",
		DSNDev:         "host=localhost dbname=syncra_dms_dev",
		ServerHostPort: "localhost:8080",
	}
	var loaded bool
	var ran bool
	deps := testCommandDeps()
	deps.loadConfig = func() (config.Config, error) {
		loaded = true
		return wantCfg, nil
	}
	deps.runAPI = func(got config.Config) error {
		ran = true
		if got.ServerHostPort != wantCfg.ServerHostPort {
			t.Fatalf("ServerHostPort = %q, want %q", got.ServerHostPort, wantCfg.ServerHostPort)
		}
		return nil
	}

	if err := run(ctx, []string{"api"}, deps); err != nil {
		t.Fatalf("run returned error: %v", err)
	}
	if !loaded {
		t.Fatal("loadConfig was not called")
	}
	if !ran {
		t.Fatal("runAPI was not called")
	}
}

func TestRunDispatchesAPIWithPortFlag(t *testing.T) {
	deps := testCommandDeps()
	deps.loadConfig = func() (config.Config, error) {
		return config.Config{ServerHostPort: "127.0.0.1:9090"}, nil
	}
	deps.runAPI = func(got config.Config) error {
		if got.ServerHostPort != "127.0.0.1:8090" {
			t.Fatalf("ServerHostPort = %q, want 127.0.0.1:8090", got.ServerHostPort)
		}
		return nil
	}

	if err := run(context.Background(), []string{"api", "--port", "8090"}, deps); err != nil {
		t.Fatalf("run returned error: %v", err)
	}
}

func TestRunRejectsInvalidAPIPortFlag(t *testing.T) {
	for _, port := range []string{"0", "65536", "not-a-port"} {
		t.Run(port, func(t *testing.T) {
			deps := testCommandDeps()
			deps.runAPI = func(config.Config) error {
				t.Fatal("runAPI should not be called")
				return nil
			}

			err := run(context.Background(), []string{"api", "--port", port}, deps)
			if err == nil {
				t.Fatal("run returned nil error")
			}
			if !strings.Contains(err.Error(), "invalid --port") {
				t.Fatalf("error = %q, want invalid --port", err.Error())
			}
		})
	}
}

func TestRunDispatchesMigrateWithAtlasApply(t *testing.T) {
	ctx := context.Background()
	wantCfg := config.MigrationConfig{AtlasDatabaseURL: "postgres://syncra:syncra@localhost:5432/syncra_dms?sslmode=disable"}
	var loaded bool
	var ran bool
	deps := testCommandDeps()
	deps.loadConfig = func() (config.Config, error) {
		t.Fatal("loadConfig should not be called for migrate")
		return config.Config{}, nil
	}
	deps.loadMigration = func() (config.MigrationConfig, error) {
		loaded = true
		return wantCfg, nil
	}
	deps.runMigrate = func(got context.Context, cfg config.MigrationConfig) error {
		ran = true
		if got != ctx {
			t.Fatal("runMigrate received a different context")
		}
		if cfg.AtlasDatabaseURL != wantCfg.AtlasDatabaseURL {
			t.Fatalf("AtlasDatabaseURL = %q, want %q", cfg.AtlasDatabaseURL, wantCfg.AtlasDatabaseURL)
		}
		return nil
	}

	if err := run(ctx, []string{"migrate"}, deps); err != nil {
		t.Fatalf("run returned error: %v", err)
	}
	if !loaded {
		t.Fatal("loadConfig was not called")
	}
	if !ran {
		t.Fatal("runMigrate was not called")
	}
}

func TestRunDispatchesSwagger(t *testing.T) {
	ctx := context.Background()
	var ran bool
	deps := testCommandDeps()
	deps.runSwagger = func(got context.Context) error {
		ran = true
		if got != ctx {
			t.Fatal("runSwagger received a different context")
		}
		return nil
	}

	if err := run(ctx, []string{"swagger"}, deps); err != nil {
		t.Fatalf("run returned error: %v", err)
	}
	if !ran {
		t.Fatal("runSwagger was not called")
	}
}

func TestRunRejectsUnknownCommand(t *testing.T) {
	var stderr bytes.Buffer
	deps := testCommandDeps()
	deps.stderr = &stderr

	err := run(context.Background(), []string{"wat"}, deps)
	if err == nil {
		t.Fatal("run returned nil error")
	}
	if !strings.Contains(err.Error(), `unknown command "wat"`) {
		t.Fatalf("error = %q, want unknown command", err.Error())
	}
	if !strings.Contains(stderr.String(), "Usage: syncra <api|migrate|swagger>") {
		t.Fatalf("stderr = %q, want usage", stderr.String())
	}
}

func TestRunMigrateUsesAtlasApply(t *testing.T) {
	root := t.TempDir()
	writeServerRoot(t, root)
	binDir := t.TempDir()
	logPath := filepath.Join(root, "atlas.log")
	atlasPath := filepath.Join(binDir, "atlas")
	script := fmt.Sprintf("#!/bin/sh\nprintf '%%s|%%s\\n' \"$PWD\" \"$*\" >> %q\n", logPath)
	if err := os.WriteFile(atlasPath, []byte(script), 0o755); err != nil {
		t.Fatalf("write fake atlas: %v", err)
	}
	t.Setenv("PATH", binDir)

	cfg := config.MigrationConfig{AtlasDatabaseURL: "postgres://syncra:syncra@localhost:5432/syncra_dms?sslmode=disable"}
	if err := runMigrate(context.Background(), root, cfg, bytes.NewBuffer(nil), bytes.NewBuffer(nil)); err != nil {
		t.Fatalf("runMigrate returned error: %v", err)
	}

	gotBytes, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read atlas log: %v", err)
	}
	want := root + "|migrate apply --env local --url " + cfg.AtlasDatabaseURL
	if got := strings.TrimSpace(string(gotBytes)); got != want {
		t.Fatalf("atlas call = %q, want %q", got, want)
	}
}

func TestRunMigrateReportsMissingAtlasBinary(t *testing.T) {
	root := t.TempDir()
	writeServerRoot(t, root)
	t.Setenv("PATH", t.TempDir())

	err := runMigrate(context.Background(), root, config.MigrationConfig{AtlasDatabaseURL: "postgres://db"}, bytes.NewBuffer(nil), bytes.NewBuffer(nil))
	if err == nil {
		t.Fatal("runMigrate returned nil error")
	}
	if !strings.Contains(err.Error(), "atlas binary not found") {
		t.Fatalf("error = %q, want missing atlas binary", err.Error())
	}
}

func TestGenerateSwaggerUsesInstalledBinary(t *testing.T) {
	root := t.TempDir()
	writeServerRoot(t, root)
	binDir := t.TempDir()
	logPath := filepath.Join(root, "swagger.log")
	swaggerPath := filepath.Join(binDir, "swagger")
	script := fmt.Sprintf("#!/bin/sh\nprintf '%%s|%%s\\n' \"$PWD\" \"$*\" >> %q\n", logPath)
	if err := os.WriteFile(swaggerPath, []byte(script), 0o755); err != nil {
		t.Fatalf("write fake swagger: %v", err)
	}
	t.Setenv("PATH", binDir)

	if err := generateSwagger(context.Background(), root, bytes.NewBuffer(nil), bytes.NewBuffer(nil)); err != nil {
		t.Fatalf("generateSwagger returned error: %v", err)
	}

	gotBytes, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read swagger log: %v", err)
	}
	got := strings.Split(strings.TrimSpace(string(gotBytes)), "\n")
	want := []string{
		root + "|generate spec --work-dir cmd/syncra --scan-models -o docs/swagger.json",
		root + "|validate docs/swagger.json",
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("swagger calls = %q, want %q", got, want)
	}
}

func TestGenerateSwaggerCreatesDocsDirectory(t *testing.T) {
	root := t.TempDir()
	writeServerRoot(t, root)
	if err := os.RemoveAll(filepath.Join(root, "docs")); err != nil {
		t.Fatalf("remove docs directory: %v", err)
	}
	binDir := t.TempDir()
	swaggerPath := filepath.Join(binDir, "swagger")
	script := `#!/bin/sh
case "$1" in
  generate)
    test -d docs || exit 7
    : > docs/swagger.json
    ;;
esac
exit 0
`
	if err := os.WriteFile(swaggerPath, []byte(script), 0o755); err != nil {
		t.Fatalf("write fake swagger: %v", err)
	}
	t.Setenv("PATH", binDir)

	if err := generateSwagger(context.Background(), root, bytes.NewBuffer(nil), bytes.NewBuffer(nil)); err != nil {
		t.Fatalf("generateSwagger returned error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "docs", "swagger.json")); err != nil {
		t.Fatalf("swagger.json was not generated: %v", err)
	}
}

func TestGenerateSwaggerReportsMissingBinary(t *testing.T) {
	root := t.TempDir()
	writeServerRoot(t, root)
	t.Setenv("PATH", t.TempDir())

	err := generateSwagger(context.Background(), root, bytes.NewBuffer(nil), bytes.NewBuffer(nil))
	if err == nil {
		t.Fatal("generateSwagger returned nil error")
	}
	if !strings.Contains(err.Error(), "swagger binary not found") {
		t.Fatalf("error = %q, want missing swagger binary", err.Error())
	}
}

func TestRunReturnsSubcommandError(t *testing.T) {
	wantErr := errors.New("boom")
	deps := testCommandDeps()
	deps.loadConfig = func() (config.Config, error) {
		return config.Config{}, nil
	}
	deps.runAPI = func(config.Config) error {
		return wantErr
	}

	err := run(context.Background(), []string{"api"}, deps)
	if !errors.Is(err, wantErr) {
		t.Fatalf("run error = %v, want %v", err, wantErr)
	}
}

func writeServerRoot(t *testing.T, root string) {
	t.Helper()
	for _, dir := range []string{
		filepath.Join(root, "cmd", "syncra"),
		filepath.Join(root, "internal", "api"),
		filepath.Join(root, "migrations"),
		filepath.Join(root, "docs"),
	} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("create %s: %v", dir, err)
		}
	}
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module ai.ro/syncra/dms\n"), 0o644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "internal", "api", "swagger_doc.go"), []byte("package api\n"), 0o644); err != nil {
		t.Fatalf("write swagger_doc.go: %v", err)
	}
}

func testCommandDeps() commandDeps {
	return commandDeps{
		stdout: bytes.NewBuffer(nil),
		stderr: bytes.NewBuffer(nil),
		loadConfig: func() (config.Config, error) {
			return config.Config{}, nil
		},
		loadMigration: func() (config.MigrationConfig, error) {
			return config.MigrationConfig{}, nil
		},
		runAPI: func(config.Config) error {
			return nil
		},
		runMigrate: func(context.Context, config.MigrationConfig) error {
			return nil
		},
		runSwagger: func(context.Context) error {
			return nil
		},
	}
}
