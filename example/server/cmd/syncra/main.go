package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"ai.ro/syncra/fixtures"
	"ai.ro/syncra/internal/app"
	"ai.ro/syncra/internal/config"
	"ai.ro/syncra/internal/database"
	"ai.ro/syncra/internal/logging"
	"ai.ro/syncra/internal/ocr"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func main() {
	logger := logging.ConfigureDefault(logging.DebugFromEnv(), os.Stdout).With(
		"service", "syncra",
		"component", "cmd_syncra",
	)
	if err := run(context.Background(), os.Args[1:], realCommandDeps()); err != nil {
		logger.Error("syncra.command_failed", "error", err)
		os.Exit(1)
	}
	slog.Default().Debug("syncra.command_completed")
}

type commandDeps struct {
	stdout        io.Writer
	stderr        io.Writer
	loadConfig    func() (config.Config, error)
	loadDatabase  func() (config.DatabaseConfig, error)
	loadMigration func() (config.MigrationConfig, error)
	runAPI        func(config.Config) error
	runDBSeed     func(context.Context, config.DatabaseConfig, dbSeedOptions) error
	runMigrate    func(context.Context, config.MigrationConfig) error
	runSwagger    func(context.Context) error
}

func realCommandDeps() commandDeps {
	deps := commandDeps{
		stdout:        os.Stdout,
		stderr:        os.Stderr,
		loadConfig:    config.Load,
		loadDatabase:  config.LoadDatabase,
		loadMigration: config.LoadMigration,
		runAPI:        app.RunAPI,
	}
	deps.runDBSeed = func(ctx context.Context, cfg config.DatabaseConfig, opts dbSeedOptions) error {
		return runDBSeed(ctx, cfg, opts, deps.stdout)
	}
	deps.runMigrate = func(ctx context.Context, cfg config.MigrationConfig) error {
		root, err := serverRoot(".")
		if err != nil {
			return err
		}
		return runMigrate(ctx, root, cfg, deps.stdout, deps.stderr)
	}
	deps.runSwagger = func(ctx context.Context) error {
		root, err := serverRoot(".")
		if err != nil {
			return err
		}
		return generateSwagger(ctx, root, deps.stdout, deps.stderr)
	}
	return deps
}

func run(ctx context.Context, args []string, deps commandDeps) error {
	if len(args) == 0 {
		printUsage(deps.stderr)
		return errors.New("command is required")
	}

	command := args[0]
	commandArgs := args[1:]

	switch command {
	case "api":
		apiOpts, err := parseAPIOptions(commandArgs, deps.stderr)
		if err != nil {
			return err
		}
		cfg, err := deps.loadConfig()
		if err != nil {
			return err
		}
		cfg = applyAPIOptions(cfg, apiOpts)
		return deps.runAPI(cfg)
	case "migrate":
		if err := rejectUnexpectedArgs(command, commandArgs, deps.stderr); err != nil {
			return err
		}
		cfg, err := deps.loadMigration()
		if err != nil {
			return err
		}
		return deps.runMigrate(ctx, cfg)
	case "dbseed":
		seedOpts, err := parseDBSeedOptions(commandArgs, deps.stderr)
		if err != nil {
			return err
		}
		cfg, err := deps.loadDatabase()
		if err != nil {
			return err
		}
		return deps.runDBSeed(ctx, cfg, seedOpts)
	case "swagger":
		if err := rejectUnexpectedArgs(command, commandArgs, deps.stderr); err != nil {
			return err
		}
		return deps.runSwagger(ctx)
	default:
		if len(commandArgs) > 0 {
			printUsage(deps.stderr)
			return fmt.Errorf("unexpected arguments for %q: %v", command, commandArgs)
		}
		printUsage(deps.stderr)
		return fmt.Errorf("unknown command %q", command)
	}
}

func printUsage(w io.Writer) {
	if w == nil {
		return
	}
	fmt.Fprintln(w, "Usage: syncra <api|dbseed|migrate|swagger>")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Commands:")
	fmt.Fprintln(w, "  api [--port PORT]  Start the HTTP API")
	fmt.Fprintln(w, "  dbseed [--recipe_categories] [--recipes]")
	fmt.Fprintln(w, "                     Seed database reference data")
	fmt.Fprintln(w, "  migrate            Run database migrations")
	fmt.Fprintln(w, "  swagger            Generate and validate Swagger docs")
}

type apiOptions struct {
	port string
}

func parseAPIOptions(args []string, stderr io.Writer) (apiOptions, error) {
	flags := flag.NewFlagSet("api", flag.ContinueOnError)
	if stderr == nil {
		flags.SetOutput(io.Discard)
	} else {
		flags.SetOutput(stderr)
	}

	var opts apiOptions
	flags.StringVar(&opts.port, "port", "", "HTTP API port")
	if err := flags.Parse(args); err != nil {
		printUsage(stderr)
		return apiOptions{}, err
	}
	if remaining := flags.Args(); len(remaining) > 0 {
		printUsage(stderr)
		return apiOptions{}, fmt.Errorf("unexpected arguments for %q: %v", "api", remaining)
	}
	if opts.port != "" {
		if err := validateAPIPort(opts.port); err != nil {
			printUsage(stderr)
			return apiOptions{}, err
		}
	}
	return opts, nil
}

func validateAPIPort(port string) error {
	parsed, err := strconv.Atoi(port)
	if err != nil || parsed < 1 || parsed > 65535 {
		return fmt.Errorf("invalid --port %q: must be an integer from 1 to 65535", port)
	}
	return nil
}

func applyAPIOptions(cfg config.Config, opts apiOptions) config.Config {
	if opts.port == "" {
		return cfg
	}

	host, _, err := net.SplitHostPort(cfg.ServerHostPort)
	if err != nil {
		host = cfg.ServerHostPort
	}
	cfg.ServerHostPort = net.JoinHostPort(host, opts.port)
	return cfg
}

func rejectUnexpectedArgs(command string, args []string, stderr io.Writer) error {
	if len(args) == 0 {
		return nil
	}
	printUsage(stderr)
	return fmt.Errorf("unexpected arguments for %q: %v", command, args)
}

type dbSeedOptions struct {
	recipeCategories bool
	recipes          bool
}

func parseDBSeedOptions(args []string, stderr io.Writer) (dbSeedOptions, error) {
	flags := flag.NewFlagSet("dbseed", flag.ContinueOnError)
	if stderr == nil {
		flags.SetOutput(io.Discard)
	} else {
		flags.SetOutput(stderr)
	}

	var opts dbSeedOptions
	flags.BoolVar(&opts.recipeCategories, "recipe_categories", false, "Seed JSON recipe categories")
	flags.BoolVar(&opts.recipes, "recipes", false, "Seed JSON recipes")
	if err := flags.Parse(args); err != nil {
		printUsage(stderr)
		return dbSeedOptions{}, err
	}
	if remaining := flags.Args(); len(remaining) > 0 {
		printUsage(stderr)
		return dbSeedOptions{}, fmt.Errorf("unexpected arguments for %q: %v", "dbseed", remaining)
	}
	return opts, nil
}

func (opts dbSeedOptions) empty() bool {
	return !opts.recipeCategories && !opts.recipes
}

type dbSeedResult struct {
	RecipeCategoriesInserted int64
	RecipeCategoriesExisting int64
	RecipesInserted          int64
	RecipesExisting          int64
}

func runDBSeed(ctx context.Context, cfg config.DatabaseConfig, opts dbSeedOptions, stdout io.Writer) error {
	db, err := database.OpenPostgres(cfg.DSN)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("database handle: %w", err)
	}
	defer sqlDB.Close()

	result, err := seedDatabase(ctx, db, opts)
	if err != nil {
		return err
	}
	if stdout != nil {
		if opts.empty() {
			fmt.Fprintln(stdout, "No database reference data to seed")
		}
		if opts.recipeCategories || opts.recipes {
			fmt.Fprintf(stdout, "Seeded recipe categories: %d inserted, %d existing\n", result.RecipeCategoriesInserted, result.RecipeCategoriesExisting)
		}
		if opts.recipes {
			fmt.Fprintf(stdout, "Seeded recipes: %d inserted, %d existing\n", result.RecipesInserted, result.RecipesExisting)
		}
	}
	return nil
}

func seedDatabase(ctx context.Context, db *gorm.DB, opts dbSeedOptions) (dbSeedResult, error) {
	var result dbSeedResult
	if opts.empty() {
		return result, nil
	}

	var categories []fixtures.RecipeCategory
	var recipes []fixtures.Recipe

	if opts.recipeCategories || opts.recipes {
		var err error
		categories, err = fixtures.RecipeCategories()
		if err != nil {
			return result, err
		}
	}
	if opts.recipes {
		var err error
		recipes, err = fixtures.Recipes()
		if err != nil {
			return result, err
		}
	}

	if opts.recipeCategories || opts.recipes {
		categoriesToSeed := categories
		if opts.recipes && !opts.recipeCategories {
			var err error
			categoriesToSeed, err = recipeCategoriesForRecipes(categories, recipes)
			if err != nil {
				return result, err
			}
		}
		inserted, existing, err := seedRecipeCategories(ctx, db, categoriesToSeed)
		if err != nil {
			return result, err
		}
		result.RecipeCategoriesInserted = inserted
		result.RecipeCategoriesExisting = existing
	}

	if opts.recipes {
		inserted, existing, err := seedRecipes(ctx, db, recipes)
		if err != nil {
			return result, err
		}
		result.RecipesInserted = inserted
		result.RecipesExisting = existing
	}

	return result, nil
}

func recipeCategoriesForRecipes(categories []fixtures.RecipeCategory, recipes []fixtures.Recipe) ([]fixtures.RecipeCategory, error) {
	referenced := make(map[string]struct{}, len(recipes))
	for _, recipe := range recipes {
		referenced[recipe.Category] = struct{}{}
	}

	categoriesToSeed := make([]fixtures.RecipeCategory, 0, len(referenced))
	for _, category := range categories {
		if _, ok := referenced[category.TitleEn]; !ok {
			continue
		}
		categoriesToSeed = append(categoriesToSeed, category)
		delete(referenced, category.TitleEn)
	}
	for category := range referenced {
		return nil, fmt.Errorf("recipe category %q is not registered in fixtures", category)
	}
	return categoriesToSeed, nil
}

func seedRecipeCategories(ctx context.Context, db *gorm.DB, categories []fixtures.RecipeCategory) (int64, int64, error) {
	var inserted int64
	var existing int64

	if err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, fixtureCategory := range categories {
			var category ocr.JSONRecipeCategory
			result := tx.Where("title_en = ?", fixtureCategory.TitleEn).
				Attrs(ocr.JSONRecipeCategory{
					TitleEn: fixtureCategory.TitleEn,
					TitleRo: fixtureCategory.TitleRo,
				}).
				FirstOrCreate(&category)
			if result.Error != nil {
				return fmt.Errorf("seed recipe category %q: %w", fixtureCategory.TitleEn, result.Error)
			}
			if result.RowsAffected == 1 {
				inserted++
			} else {
				existing++
			}
		}
		return nil
	}); err != nil {
		return 0, 0, err
	}

	return inserted, existing, nil
}

func seedRecipes(ctx context.Context, db *gorm.DB, recipes []fixtures.Recipe) (int64, int64, error) {
	var inserted int64
	var existing int64

	if err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		categoryTitles := make([]string, 0, len(recipes))
		seenCategoryTitles := make(map[string]struct{}, len(recipes))
		for _, recipe := range recipes {
			if _, ok := seenCategoryTitles[recipe.Category]; ok {
				continue
			}
			seenCategoryTitles[recipe.Category] = struct{}{}
			categoryTitles = append(categoryTitles, recipe.Category)
		}

		categoriesByTitle := make(map[string]ocr.JSONRecipeCategory, len(categoryTitles))
		if len(categoryTitles) > 0 {
			var categories []ocr.JSONRecipeCategory
			if err := tx.Where("title_en IN ?", categoryTitles).Find(&categories).Error; err != nil {
				return fmt.Errorf("load recipe categories: %w", err)
			}
			for _, category := range categories {
				categoriesByTitle[category.TitleEn] = category
			}
		}

		for _, fixtureRecipe := range recipes {
			category, ok := categoriesByTitle[fixtureRecipe.Category]
			if !ok {
				return fmt.Errorf("seed recipe %q: category %q is not registered", fixtureRecipe.Title, fixtureRecipe.Category)
			}

			categoryID := category.ID
			var recipe ocr.JSONRecipe
			result := tx.Where("title = ?", fixtureRecipe.Title).
				Attrs(ocr.JSONRecipe{
					Title:       fixtureRecipe.Title,
					Description: fixtureRecipe.Description,
					JSON:        datatypes.JSON(fixtureRecipe.JSON),
					CategoryID:  &categoryID,
				}).
				FirstOrCreate(&recipe)
			if result.Error != nil {
				return fmt.Errorf("seed recipe %q: %w", fixtureRecipe.Title, result.Error)
			}
			if result.RowsAffected == 1 {
				inserted++
			} else {
				existing++
			}
		}
		return nil
	}); err != nil {
		return 0, 0, err
	}

	return inserted, existing, nil
}

func runMigrate(ctx context.Context, root string, cfg config.MigrationConfig, stdout, stderr io.Writer) error {
	root, err := serverRoot(root)
	if err != nil {
		return err
	}

	atlasPath, err := exec.LookPath("atlas")
	if err != nil {
		return fmt.Errorf("atlas binary not found in PATH: %w", err)
	}

	apply := exec.CommandContext(ctx, atlasPath,
		"migrate", "apply",
		"--env", "local",
		"--url", cfg.AtlasDatabaseURL,
	)
	apply.Dir = root
	apply.Stdout = stdout
	apply.Stderr = stderr
	if err := apply.Run(); err != nil {
		return fmt.Errorf("atlas migrate apply: %w", err)
	}
	return nil
}

func generateSwagger(ctx context.Context, root string, stdout, stderr io.Writer) error {
	root, err := serverRoot(root)
	if err != nil {
		return err
	}

	swaggerPath, err := exec.LookPath("swagger")
	if err != nil {
		return fmt.Errorf("swagger binary not found in PATH: %w", err)
	}

	generate := exec.CommandContext(ctx, swaggerPath,
		"generate", "spec",
		"--work-dir", "cmd/syncra",
		"--scan-models",
		"-o", "docs/swagger.json",
	)
	generate.Dir = root
	generate.Stdout = stdout
	generate.Stderr = stderr
	if err := generate.Run(); err != nil {
		return fmt.Errorf("swagger generate spec: %w", err)
	}

	validate := exec.CommandContext(ctx, swaggerPath, "validate", "docs/swagger.json")
	validate.Dir = root
	validate.Stdout = stdout
	validate.Stderr = stderr
	if err := validate.Run(); err != nil {
		return fmt.Errorf("swagger validate docs/swagger.json: %w", err)
	}

	return nil
}

func serverRoot(start string) (string, error) {
	abs, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(abs)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		abs = filepath.Dir(abs)
	}
	for {
		if hasServerRootFiles(abs) {
			return abs, nil
		}
		next := filepath.Dir(abs)
		if next == abs {
			return "", fmt.Errorf("syncra swagger must be run from the server module root or one of its subdirectories")
		}
		abs = next
	}
}

func hasServerRootFiles(dir string) bool {
	for _, path := range []string{
		filepath.Join(dir, "go.mod"),
		filepath.Join(dir, "cmd", "syncra"),
		filepath.Join(dir, "internal", "api", "swagger_doc.go"),
	} {
		if _, err := os.Stat(path); err != nil {
			return false
		}
	}
	return true
}
