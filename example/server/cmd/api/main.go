package main

import (
	"log/slog"
	"os"

	"ai.ro/syncra/internal/app"
	"ai.ro/syncra/internal/config"
	"ai.ro/syncra/internal/logging"
)

func main() {
	logger := logging.ConfigureDefault(logging.DebugFromEnv(), os.Stdout).With(
		"service", "syncra",
		"component", "cmd_api",
	)
	cfg, err := config.Load()
	if err != nil {
		logger.Error("api.config_load_failed", "error", err)
		os.Exit(1)
	}
	if err := app.RunAPI(cfg); err != nil {
		slog.Default().Error("api.server_failed", "error", err)
		os.Exit(1)
	}
}
