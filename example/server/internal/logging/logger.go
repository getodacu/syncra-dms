package logging

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

type contextKey struct{}

func NewJSONLogger(w io.Writer, debug bool) *slog.Logger {
	if w == nil {
		w = io.Discard
	}
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}
	return slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: level,
	}))
}

func ConfigureDefault(debug bool, w io.Writer) *slog.Logger {
	logger := NewJSONLogger(w, debug)
	slog.SetDefault(logger)
	return logger
}

func Nop() *slog.Logger {
	return NewJSONLogger(io.Discard, true)
}

func WithContext(ctx context.Context, logger *slog.Logger) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if logger == nil {
		return ctx
	}
	return context.WithValue(ctx, contextKey{}, logger)
}

func FromContext(ctx context.Context) *slog.Logger {
	if ctx != nil {
		if logger, ok := ctx.Value(contextKey{}).(*slog.Logger); ok && logger != nil {
			return logger
		}
	}
	return slog.Default()
}

func DebugFromEnv() bool {
	value := strings.TrimSpace(os.Getenv("DEBUG"))
	if value == "" {
		return false
	}
	parsed, err := strconv.ParseBool(value)
	return err == nil && parsed
}
