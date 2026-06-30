package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"ai.ro/syncra/dms/internal/app"
	"ai.ro/syncra/dms/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "configuration error: %v\n", err)
		os.Exit(1)
	}
	if err := app.RunAPI(cfg); err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Fprintf(os.Stderr, "api server error: %v\n", err)
		os.Exit(1)
	}
}
