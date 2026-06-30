package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"ai.ro/syncra/dms/internal/app"
	"ai.ro/syncra/dms/internal/config"
	"ai.ro/syncra/dms/internal/logging"
)

func main() {
	logger := logging.ConfigureDefault(false, os.Stdout).With(
		"service", "syncra-dms",
		"component", "cmd_syncra",
	)
	if err := run(context.Background(), os.Args[1:], realCommandDeps()); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return
		}
		logger.Error("syncra.command_failed", "error", err)
		os.Exit(1)
	}
	slog.Default().Debug("syncra.command_completed")
}

type commandDeps struct {
	stdout        io.Writer
	stderr        io.Writer
	loadConfig    func() (config.Config, error)
	loadMigration func() (config.MigrationConfig, error)
	runAPI        func(config.Config) error
	runMigrate    func(context.Context, config.MigrationConfig) error
	runSwagger    func(context.Context) error
}

func realCommandDeps() commandDeps {
	deps := commandDeps{
		stdout:        os.Stdout,
		stderr:        os.Stderr,
		loadConfig:    config.Load,
		loadMigration: config.LoadMigration,
		runAPI:        app.RunAPI,
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
		return deps.runAPI(applyAPIOptions(cfg, apiOpts))
	case "migrate":
		if err := rejectUnexpectedArgs(command, commandArgs, deps.stderr); err != nil {
			return err
		}
		cfg, err := deps.loadMigration()
		if err != nil {
			return err
		}
		return deps.runMigrate(ctx, cfg)
	case "swagger":
		if err := rejectUnexpectedArgs(command, commandArgs, deps.stderr); err != nil {
			return err
		}
		return deps.runSwagger(ctx)
	default:
		printUsage(deps.stderr)
		return fmt.Errorf("unknown command %q", command)
	}
}

func printUsage(w io.Writer) {
	if w == nil {
		return
	}
	fmt.Fprintln(w, "Usage: syncra <api|migrate|swagger>")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Commands:")
	fmt.Fprintln(w, "  api [--port PORT]  Start the HTTP API")
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

func runMigrate(ctx context.Context, root string, cfg config.MigrationConfig, stdout io.Writer, stderr io.Writer) error {
	root, err := serverRoot(root)
	if err != nil {
		return err
	}
	if cfg.AtlasDatabaseURL == "" {
		return errors.New("ATLAS_DATABASE_URL is required")
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

func generateSwagger(ctx context.Context, root string, stdout io.Writer, stderr io.Writer) error {
	root, err := serverRoot(root)
	if err != nil {
		return err
	}
	swaggerPath, err := exec.LookPath("swagger")
	if err != nil {
		return fmt.Errorf("swagger binary not found in PATH: %w", err)
	}
	if err := os.MkdirAll(filepath.Join(root, "docs"), 0o755); err != nil {
		return fmt.Errorf("create docs directory: %w", err)
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
			return "", errors.New("syncra must be run from the server module root or one of its subdirectories")
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
