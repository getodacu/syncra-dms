package logging

import (
	"io"
	"log/slog"
)

func ConfigureDefault(debug bool, output io.Writer) *slog.Logger {
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}
	logger := slog.New(slog.NewTextHandler(output, &slog.HandlerOptions{Level: level}))
	slog.SetDefault(logger)
	return logger
}

func Nop() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
