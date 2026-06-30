package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ai.ro/syncra/fixtures"
	"ai.ro/syncra/internal/config"
	"ai.ro/syncra/internal/ocr"

	"gorm.io/datatypes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestRunDispatchesAPI(t *testing.T) {
	ctx := context.Background()
	wantCfg := config.Config{
		DSN:            "postgres://api",
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
		if got.DSN != wantCfg.DSN {
			t.Fatalf("DSN = %q, want %q", got.DSN, wantCfg.DSN)
		}
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
	ctx := context.Background()
	var ran bool
	deps := testCommandDeps()
	deps.loadConfig = func() (config.Config, error) {
		return config.Config{
			DSN:            "postgres://api",
			ServerHostPort: "127.0.0.1:9090",
		}, nil
	}
	deps.runAPI = func(got config.Config) error {
		ran = true
		if got.ServerHostPort != "127.0.0.1:8081" {
			t.Fatalf("ServerHostPort = %q, want 127.0.0.1:8081", got.ServerHostPort)
		}
		return nil
	}

	if err := run(ctx, []string{"api", "--port", "8081"}, deps); err != nil {
		t.Fatalf("run returned error: %v", err)
	}
	if !ran {
		t.Fatal("runAPI was not called")
	}
}

func TestRunDispatchesAPIWithEqualsPortFlag(t *testing.T) {
	ctx := context.Background()
	var ran bool
	deps := testCommandDeps()
	deps.loadConfig = func() (config.Config, error) {
		return config.Config{
			DSN:            "postgres://api",
			ServerHostPort: "localhost:9090",
		}, nil
	}
	deps.runAPI = func(got config.Config) error {
		ran = true
		if got.ServerHostPort != "localhost:8082" {
			t.Fatalf("ServerHostPort = %q, want localhost:8082", got.ServerHostPort)
		}
		return nil
	}

	if err := run(ctx, []string{"api", "--port=8082"}, deps); err != nil {
		t.Fatalf("run returned error: %v", err)
	}
	if !ran {
		t.Fatal("runAPI was not called")
	}
}

func TestRunRejectsInvalidAPIPortFlag(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "zero",
			args: []string{"api", "--port", "0"},
		},
		{
			name: "too high",
			args: []string{"api", "--port", "65536"},
		},
		{
			name: "not numeric",
			args: []string{"api", "--port", "not-a-number"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ran bool
			deps := testCommandDeps()
			deps.loadConfig = func() (config.Config, error) {
				return config.Config{
					DSN:            "postgres://api",
					ServerHostPort: "localhost:8080",
				}, nil
			}
			deps.runAPI = func(config.Config) error {
				ran = true
				return nil
			}

			err := run(context.Background(), tt.args, deps)
			if err == nil {
				t.Fatal("run returned nil error")
			}
			if !strings.Contains(err.Error(), "invalid --port") {
				t.Fatalf("error = %q, want invalid --port", err.Error())
			}
			if ran {
				t.Fatal("runAPI was called")
			}
		})
	}
}

func TestRunDispatchesMigrateWithMigrationConfig(t *testing.T) {
	ctx := context.Background()
	wantCfg := config.MigrationConfig{AtlasDatabaseURL: "postgres://migrate"}
	var loaded bool
	var ran bool
	deps := testCommandDeps()
	deps.loadMigration = func() (config.MigrationConfig, error) {
		loaded = true
		return wantCfg, nil
	}
	deps.runMigrate = func(got context.Context, cfg config.MigrationConfig) error {
		ran = true
		if got != ctx {
			t.Fatal("runMigrate received a different context")
		}
		if cfg != wantCfg {
			t.Fatalf("runMigrate config = %+v, want %+v", cfg, wantCfg)
		}
		return nil
	}

	if err := run(ctx, []string{"migrate"}, deps); err != nil {
		t.Fatalf("run returned error: %v", err)
	}
	if !loaded {
		t.Fatal("loadMigration was not called")
	}
	if !ran {
		t.Fatal("runMigrate was not called")
	}
}

func TestRunDispatchesDBSeedWithDatabaseConfig(t *testing.T) {
	ctx := context.Background()
	wantCfg := config.DatabaseConfig{DSN: "postgres://seed"}
	var loaded bool
	var ran bool
	deps := testCommandDeps()
	deps.loadDatabase = func() (config.DatabaseConfig, error) {
		loaded = true
		return wantCfg, nil
	}
	deps.runDBSeed = func(got context.Context, cfg config.DatabaseConfig, opts dbSeedOptions) error {
		ran = true
		if got != ctx {
			t.Fatal("runDBSeed received a different context")
		}
		if cfg != wantCfg {
			t.Fatalf("runDBSeed config = %+v, want %+v", cfg, wantCfg)
		}
		if !opts.empty() {
			t.Fatalf("dbseed options = %+v, want empty", opts)
		}
		return nil
	}

	if err := run(ctx, []string{"dbseed"}, deps); err != nil {
		t.Fatalf("run returned error: %v", err)
	}
	if !loaded {
		t.Fatal("loadDatabase was not called")
	}
	if !ran {
		t.Fatal("runDBSeed was not called")
	}
}

func TestRunDispatchesDBSeedWithRecipeCategoriesFlag(t *testing.T) {
	ctx := context.Background()
	wantCfg := config.DatabaseConfig{DSN: "postgres://seed"}
	var ran bool
	deps := testCommandDeps()
	deps.loadDatabase = func() (config.DatabaseConfig, error) {
		return wantCfg, nil
	}
	deps.runDBSeed = func(got context.Context, cfg config.DatabaseConfig, opts dbSeedOptions) error {
		ran = true
		if got != ctx {
			t.Fatal("runDBSeed received a different context")
		}
		if cfg != wantCfg {
			t.Fatalf("runDBSeed config = %+v, want %+v", cfg, wantCfg)
		}
		if !opts.recipeCategories {
			t.Fatalf("recipeCategories = false, want true")
		}
		return nil
	}

	if err := run(ctx, []string{"dbseed", "--recipe_categories"}, deps); err != nil {
		t.Fatalf("run returned error: %v", err)
	}
	if !ran {
		t.Fatal("runDBSeed was not called")
	}
}

func TestRunDispatchesDBSeedWithRecipesFlag(t *testing.T) {
	ctx := context.Background()
	wantCfg := config.DatabaseConfig{DSN: "postgres://seed"}
	var ran bool
	deps := testCommandDeps()
	deps.loadDatabase = func() (config.DatabaseConfig, error) {
		return wantCfg, nil
	}
	deps.runDBSeed = func(got context.Context, cfg config.DatabaseConfig, opts dbSeedOptions) error {
		ran = true
		if got != ctx {
			t.Fatal("runDBSeed received a different context")
		}
		if cfg != wantCfg {
			t.Fatalf("runDBSeed config = %+v, want %+v", cfg, wantCfg)
		}
		if !opts.recipes {
			t.Fatalf("recipes = false, want true")
		}
		return nil
	}

	if err := run(ctx, []string{"dbseed", "--recipes"}, deps); err != nil {
		t.Fatalf("run returned error: %v", err)
	}
	if !ran {
		t.Fatal("runDBSeed was not called")
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

	err := run(context.Background(), []string{"unknown"}, deps)
	if err == nil {
		t.Fatal("run returned nil error")
	}
	if !strings.Contains(err.Error(), `unknown command "unknown"`) {
		t.Fatalf("error = %q, want unknown command", err.Error())
	}
	if !strings.Contains(stderr.String(), "Usage: syncra <api|dbseed|migrate|swagger>") {
		t.Fatalf("stderr = %q, want usage", stderr.String())
	}
}

func TestRunRejectsMissingCommand(t *testing.T) {
	var stderr bytes.Buffer
	deps := testCommandDeps()
	deps.stderr = &stderr

	err := run(context.Background(), nil, deps)
	if err == nil {
		t.Fatal("run returned nil error")
	}
	if !strings.Contains(err.Error(), "command is required") {
		t.Fatalf("error = %q, want missing command", err.Error())
	}
	if !strings.Contains(stderr.String(), "Commands:") {
		t.Fatalf("stderr = %q, want commands", stderr.String())
	}
}

func TestRunRejectsUnexpectedArgsForNonAPICommand(t *testing.T) {
	var stderr bytes.Buffer
	deps := testCommandDeps()
	deps.stderr = &stderr

	err := run(context.Background(), []string{"migrate", "--port", "8081"}, deps)
	if err == nil {
		t.Fatal("run returned nil error")
	}
	if !strings.Contains(err.Error(), `unexpected arguments for "migrate"`) {
		t.Fatalf("error = %q, want unexpected arguments", err.Error())
	}
	if !strings.Contains(stderr.String(), "Usage: syncra <api|dbseed|migrate|swagger>") {
		t.Fatalf("stderr = %q, want usage", stderr.String())
	}
}

func TestRunRejectsUnexpectedDBSeedArgs(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "unknown flag",
			args: []string{"dbseed", "--recipez"},
			want: "flag provided but not defined",
		},
		{
			name: "positional argument",
			args: []string{"dbseed", "recipe_categories"},
			want: `unexpected arguments for "dbseed"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stderr bytes.Buffer
			deps := testCommandDeps()
			deps.stderr = &stderr
			deps.loadDatabase = func() (config.DatabaseConfig, error) {
				t.Fatal("loadDatabase should not be called")
				return config.DatabaseConfig{}, nil
			}

			err := run(context.Background(), tt.args, deps)
			if err == nil {
				t.Fatal("run returned nil error")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error = %q, want %q", err.Error(), tt.want)
			}
			if !strings.Contains(stderr.String(), "Usage: syncra <api|dbseed|migrate|swagger>") {
				t.Fatalf("stderr = %q, want usage", stderr.String())
			}
		})
	}
}

func TestRunReturnsFailureOnCommandError(t *testing.T) {
	wantErr := errors.New("api failed")
	deps := testCommandDeps()
	deps.loadConfig = func() (config.Config, error) {
		return config.Config{DSN: "postgres://api"}, nil
	}
	deps.runAPI = func(config.Config) error {
		return wantErr
	}

	err := run(context.Background(), []string{"api"}, deps)
	if !errors.Is(err, wantErr) {
		t.Fatalf("run error = %v, want %v", err, wantErr)
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

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if err := generateSwagger(context.Background(), root, &stdout, &stderr); err != nil {
		t.Fatalf("generateSwagger returned error: %v; stderr: %s", err, stderr.String())
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
	if len(got) != len(want) {
		t.Fatalf("swagger calls = %q, want %q", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("swagger call %d = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestGenerateSwaggerFindsServerRootFromSubdirectory(t *testing.T) {
	root := t.TempDir()
	writeServerRoot(t, root)
	subdir := filepath.Join(root, "cmd", "syncra")
	binDir := t.TempDir()
	logPath := filepath.Join(root, "swagger.log")
	swaggerPath := filepath.Join(binDir, "swagger")
	script := fmt.Sprintf("#!/bin/sh\nprintf '%%s|%%s\\n' \"$PWD\" \"$*\" >> %q\n", logPath)
	if err := os.WriteFile(swaggerPath, []byte(script), 0o755); err != nil {
		t.Fatalf("write fake swagger: %v", err)
	}
	t.Setenv("PATH", binDir)

	if err := generateSwagger(context.Background(), subdir, bytes.NewBuffer(nil), bytes.NewBuffer(nil)); err != nil {
		t.Fatalf("generateSwagger returned error: %v", err)
	}

	gotBytes, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read swagger log: %v", err)
	}
	if got := strings.Split(strings.TrimSpace(string(gotBytes)), "\n")[0]; !strings.HasPrefix(got, root+"|") {
		t.Fatalf("first swagger call = %q, want root %q", got, root)
	}
}

func TestGenerateSwaggerRejectsOutsideServerRoot(t *testing.T) {
	root := t.TempDir()
	writeServerRoot(t, root)

	err := generateSwagger(context.Background(), filepath.Dir(root), bytes.NewBuffer(nil), bytes.NewBuffer(nil))
	if err == nil {
		t.Fatal("generateSwagger returned nil error")
	}
	if !strings.Contains(err.Error(), "server module root") {
		t.Fatalf("error = %q, want server module root", err.Error())
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

func TestGenerateSwaggerSkipsValidateAfterGenerateFailure(t *testing.T) {
	root := t.TempDir()
	writeServerRoot(t, root)
	binDir := t.TempDir()
	logPath := filepath.Join(root, "swagger.log")
	swaggerPath := filepath.Join(binDir, "swagger")
	script := fmt.Sprintf(`#!/bin/sh
printf '%%s\n' "$*" >> %q
case "$1 $2" in
  "generate spec") exit 7 ;;
esac
exit 0
`, logPath)
	if err := os.WriteFile(swaggerPath, []byte(script), 0o755); err != nil {
		t.Fatalf("write fake swagger: %v", err)
	}
	t.Setenv("PATH", binDir)

	err := generateSwagger(context.Background(), root, bytes.NewBuffer(nil), bytes.NewBuffer(nil))
	if err == nil {
		t.Fatal("generateSwagger returned nil error")
	}
	if !strings.Contains(err.Error(), "swagger generate spec") {
		t.Fatalf("error = %q, want generate failure", err.Error())
	}
	gotBytes, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read swagger log: %v", err)
	}
	if got := string(gotBytes); strings.Contains(got, "validate") {
		t.Fatalf("swagger calls = %q, validate should not run after generate failure", got)
	}
}

func TestGenerateSwaggerReportsValidateFailure(t *testing.T) {
	root := t.TempDir()
	writeServerRoot(t, root)
	binDir := t.TempDir()
	swaggerPath := filepath.Join(binDir, "swagger")
	script := `#!/bin/sh
case "$1" in
  validate) exit 9 ;;
esac
exit 0
`
	if err := os.WriteFile(swaggerPath, []byte(script), 0o755); err != nil {
		t.Fatalf("write fake swagger: %v", err)
	}
	t.Setenv("PATH", binDir)

	err := generateSwagger(context.Background(), root, bytes.NewBuffer(nil), bytes.NewBuffer(nil))
	if err == nil {
		t.Fatal("generateSwagger returned nil error")
	}
	if !strings.Contains(err.Error(), "swagger validate docs/swagger.json") {
		t.Fatalf("error = %q, want validate failure", err.Error())
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

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cfg := config.MigrationConfig{AtlasDatabaseURL: "postgres://postgres:pass@localhost:5432/syncra_dev?search_path=public&sslmode=disable"}
	if err := runMigrate(context.Background(), root, cfg, &stdout, &stderr); err != nil {
		t.Fatalf("runMigrate returned error: %v; stderr: %s", err, stderr.String())
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

func TestRunMigrateReportsAtlasFailure(t *testing.T) {
	root := t.TempDir()
	writeServerRoot(t, root)
	binDir := t.TempDir()
	atlasPath := filepath.Join(binDir, "atlas")
	if err := os.WriteFile(atlasPath, []byte("#!/bin/sh\nexit 11\n"), 0o755); err != nil {
		t.Fatalf("write fake atlas: %v", err)
	}
	t.Setenv("PATH", binDir)

	err := runMigrate(context.Background(), root, config.MigrationConfig{AtlasDatabaseURL: "postgres://db"}, bytes.NewBuffer(nil), bytes.NewBuffer(nil))
	if err == nil {
		t.Fatal("runMigrate returned nil error")
	}
	if !strings.Contains(err.Error(), "atlas migrate apply") {
		t.Fatalf("error = %q, want atlas migrate apply failure", err.Error())
	}
}

func TestSeedDatabaseWithNoOptionsHasNoReferenceData(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	result, err := seedDatabase(context.Background(), db, dbSeedOptions{})
	if err != nil {
		t.Fatalf("seedDatabase() error = %v", err)
	}
	if result != (dbSeedResult{}) {
		t.Fatalf("seed result = %+v, want empty", result)
	}

	var tableCount int64
	if err := db.Raw(`SELECT count(*) FROM sqlite_master WHERE type = 'table'`).Scan(&tableCount).Error; err != nil {
		t.Fatalf("count sqlite tables: %v", err)
	}
	if tableCount != 0 {
		t.Fatalf("table count = %d, want 0", tableCount)
	}
}

func TestSeedDatabaseRecipeCategoriesSeedsFixture(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&ocr.JSONRecipeCategory{}); err != nil {
		t.Fatalf("automigrate category: %v", err)
	}

	fixtureCategories, err := fixtures.RecipeCategories()
	if err != nil {
		t.Fatalf("load fixture categories: %v", err)
	}
	opts := dbSeedOptions{recipeCategories: true}
	result, err := seedDatabase(context.Background(), db, opts)
	if err != nil {
		t.Fatalf("seedDatabase() error = %v", err)
	}
	if result.RecipeCategoriesInserted != int64(len(fixtureCategories)) || result.RecipeCategoriesExisting != 0 {
		t.Fatalf("seed result = %+v, want %d inserted and 0 existing", result, len(fixtureCategories))
	}

	var count int64
	if err := db.Model(&ocr.JSONRecipeCategory{}).Count(&count).Error; err != nil {
		t.Fatalf("count categories: %v", err)
	}
	if count != int64(len(fixtureCategories)) {
		t.Fatalf("category count = %d, want %d", count, len(fixtureCategories))
	}

	result, err = seedDatabase(context.Background(), db, opts)
	if err != nil {
		t.Fatalf("seedDatabase() rerun error = %v", err)
	}
	if result.RecipeCategoriesInserted != 0 || result.RecipeCategoriesExisting != int64(len(fixtureCategories)) {
		t.Fatalf("rerun seed result = %+v, want 0 inserted and %d existing", result, len(fixtureCategories))
	}
	if err := db.Model(&ocr.JSONRecipeCategory{}).Count(&count).Error; err != nil {
		t.Fatalf("count categories after rerun: %v", err)
	}
	if count != int64(len(fixtureCategories)) {
		t.Fatalf("category count after rerun = %d, want %d", count, len(fixtureCategories))
	}
}

func TestSeedDatabaseRecipeCategoriesPreservesExistingTitle(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&ocr.JSONRecipeCategory{}); err != nil {
		t.Fatalf("automigrate category: %v", err)
	}

	existing := ocr.JSONRecipeCategory{
		TitleEn: "Finance and Accounting",
		TitleRo: "Custom Romanian Title",
	}
	if err := db.Create(&existing).Error; err != nil {
		t.Fatalf("create existing category: %v", err)
	}
	fixtureCategories, err := fixtures.RecipeCategories()
	if err != nil {
		t.Fatalf("load fixture categories: %v", err)
	}

	result, err := seedDatabase(context.Background(), db, dbSeedOptions{recipeCategories: true})
	if err != nil {
		t.Fatalf("seedDatabase() error = %v", err)
	}
	if result.RecipeCategoriesInserted != int64(len(fixtureCategories)-1) || result.RecipeCategoriesExisting != 1 {
		t.Fatalf("seed result = %+v, want %d inserted and 1 existing", result, len(fixtureCategories)-1)
	}

	var got ocr.JSONRecipeCategory
	if err := db.First(&got, "title_en = ?", existing.TitleEn).Error; err != nil {
		t.Fatalf("load existing category: %v", err)
	}
	if got.ID != existing.ID {
		t.Fatalf("existing category id = %s, want %s", got.ID, existing.ID)
	}
	if got.TitleRo != existing.TitleRo {
		t.Fatalf("existing category TitleRo = %q, want %q", got.TitleRo, existing.TitleRo)
	}

	var duplicateCount int64
	if err := db.Model(&ocr.JSONRecipeCategory{}).Where("title_en = ?", existing.TitleEn).Count(&duplicateCount).Error; err != nil {
		t.Fatalf("count duplicate title: %v", err)
	}
	if duplicateCount != 1 {
		t.Fatalf("duplicate count = %d, want 1", duplicateCount)
	}
}

func TestSeedDatabaseRecipesSeedsFixtureAndReferencedCategories(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&ocr.JSONRecipeCategory{}, &ocr.JSONRecipe{}); err != nil {
		t.Fatalf("automigrate recipes: %v", err)
	}

	fixtureRecipes, err := fixtures.Recipes()
	if err != nil {
		t.Fatalf("load fixture recipes: %v", err)
	}
	referencedCategories := referencedRecipeCategoryCount(fixtureRecipes)

	opts := dbSeedOptions{recipes: true}
	result, err := seedDatabase(context.Background(), db, opts)
	if err != nil {
		t.Fatalf("seedDatabase() error = %v", err)
	}
	if result.RecipesInserted != int64(len(fixtureRecipes)) || result.RecipesExisting != 0 {
		t.Fatalf("seed result = %+v, want %d recipes inserted and 0 existing", result, len(fixtureRecipes))
	}
	if result.RecipeCategoriesInserted != int64(referencedCategories) || result.RecipeCategoriesExisting != 0 {
		t.Fatalf("seed result = %+v, want %d categories inserted and 0 existing", result, referencedCategories)
	}

	var recipeCount int64
	if err := db.Model(&ocr.JSONRecipe{}).Count(&recipeCount).Error; err != nil {
		t.Fatalf("count recipes: %v", err)
	}
	if recipeCount != int64(len(fixtureRecipes)) {
		t.Fatalf("recipe count = %d, want %d", recipeCount, len(fixtureRecipes))
	}

	var categoryCount int64
	if err := db.Model(&ocr.JSONRecipeCategory{}).Count(&categoryCount).Error; err != nil {
		t.Fatalf("count categories: %v", err)
	}
	if categoryCount != int64(referencedCategories) {
		t.Fatalf("category count = %d, want %d", categoryCount, referencedCategories)
	}

	var uncategorizedCount int64
	if err := db.Model(&ocr.JSONRecipe{}).Where("category_id IS NULL").Count(&uncategorizedCount).Error; err != nil {
		t.Fatalf("count uncategorized recipes: %v", err)
	}
	if uncategorizedCount != 0 {
		t.Fatalf("uncategorized recipe count = %d, want 0", uncategorizedCount)
	}

	result, err = seedDatabase(context.Background(), db, opts)
	if err != nil {
		t.Fatalf("seedDatabase() rerun error = %v", err)
	}
	if result.RecipesInserted != 0 || result.RecipesExisting != int64(len(fixtureRecipes)) {
		t.Fatalf("rerun seed result = %+v, want 0 recipes inserted and %d existing", result, len(fixtureRecipes))
	}
	if result.RecipeCategoriesInserted != 0 || result.RecipeCategoriesExisting != int64(referencedCategories) {
		t.Fatalf("rerun seed result = %+v, want 0 categories inserted and %d existing", result, referencedCategories)
	}
	if err := db.Model(&ocr.JSONRecipe{}).Count(&recipeCount).Error; err != nil {
		t.Fatalf("count recipes after rerun: %v", err)
	}
	if recipeCount != int64(len(fixtureRecipes)) {
		t.Fatalf("recipe count after rerun = %d, want %d", recipeCount, len(fixtureRecipes))
	}
}

func TestSeedDatabaseRecipesPreservesExistingTitle(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&ocr.JSONRecipeCategory{}, &ocr.JSONRecipe{}); err != nil {
		t.Fatalf("automigrate recipes: %v", err)
	}

	fixtureRecipes, err := fixtures.Recipes()
	if err != nil {
		t.Fatalf("load fixture recipes: %v", err)
	}
	existingJSON := datatypes.JSON([]byte(`{"type":"object","properties":{"custom":{"type":"string"}}}`))
	existing := ocr.JSONRecipe{
		Title:       fixtureRecipes[0].Title,
		Description: "Custom description",
		JSON:        existingJSON,
		Counter:     7,
	}
	if err := db.Create(&existing).Error; err != nil {
		t.Fatalf("create existing recipe: %v", err)
	}

	result, err := seedDatabase(context.Background(), db, dbSeedOptions{recipes: true})
	if err != nil {
		t.Fatalf("seedDatabase() error = %v", err)
	}
	if result.RecipesInserted != int64(len(fixtureRecipes)-1) || result.RecipesExisting != 1 {
		t.Fatalf("seed result = %+v, want %d recipes inserted and 1 existing", result, len(fixtureRecipes)-1)
	}

	var got ocr.JSONRecipe
	if err := db.First(&got, "id = ?", existing.ID).Error; err != nil {
		t.Fatalf("load existing recipe: %v", err)
	}
	if got.Description != existing.Description {
		t.Fatalf("existing recipe description = %q, want %q", got.Description, existing.Description)
	}
	if got.Counter != existing.Counter {
		t.Fatalf("existing recipe counter = %d, want %d", got.Counter, existing.Counter)
	}
	if got.CategoryID != nil {
		t.Fatalf("existing recipe category_id = %v, want nil", got.CategoryID)
	}
	if string(got.JSON) != string(existingJSON) {
		t.Fatalf("existing recipe JSON = %s, want %s", got.JSON, existingJSON)
	}

	var duplicateCount int64
	if err := db.Model(&ocr.JSONRecipe{}).Where("title = ?", existing.Title).Count(&duplicateCount).Error; err != nil {
		t.Fatalf("count duplicate title: %v", err)
	}
	if duplicateCount != 1 {
		t.Fatalf("duplicate count = %d, want 1", duplicateCount)
	}
}

func TestSeedRecipesRejectsMissingCategory(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&ocr.JSONRecipeCategory{}, &ocr.JSONRecipe{}); err != nil {
		t.Fatalf("automigrate recipes: %v", err)
	}

	_, _, err = seedRecipes(context.Background(), db, []fixtures.Recipe{
		{
			Title:       "Missing category recipe",
			Description: "Description",
			Category:    "Missing",
			JSON:        json.RawMessage(`{"type":"object"}`),
		},
	})
	if err == nil {
		t.Fatal("seedRecipes() returned nil error")
	}
	if !strings.Contains(err.Error(), `category "Missing" is not registered`) {
		t.Fatalf("error = %q, want missing category", err.Error())
	}
}

func referencedRecipeCategoryCount(recipes []fixtures.Recipe) int {
	categories := make(map[string]struct{}, len(recipes))
	for _, recipe := range recipes {
		categories[recipe.Category] = struct{}{}
	}
	return len(categories)
}

func writeServerRoot(t *testing.T, root string) {
	t.Helper()
	for _, dir := range []string{
		filepath.Join(root, "cmd", "syncra"),
		filepath.Join(root, "internal", "api"),
	} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("create %s: %v", dir, err)
		}
	}
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module ai.ro/syncra\n"), 0o644); err != nil {
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
		loadDatabase: func() (config.DatabaseConfig, error) {
			return config.DatabaseConfig{}, nil
		},
		loadMigration: func() (config.MigrationConfig, error) {
			return config.MigrationConfig{}, nil
		},
		runAPI: func(config.Config) error {
			return nil
		},
		runDBSeed: func(context.Context, config.DatabaseConfig, dbSeedOptions) error {
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
