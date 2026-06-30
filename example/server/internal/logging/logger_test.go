package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
)

func TestNewJSONLoggerWritesStructuredInfo(t *testing.T) {
	var buf bytes.Buffer
	logger := NewJSONLogger(&buf, false)

	logger.Info("server.started", "component", "test", "port", 8080)

	var got map[string]any
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("decode log JSON: %v\n%s", err, buf.String())
	}
	if got["msg"] != "server.started" {
		t.Fatalf("msg = %#v, want server.started", got["msg"])
	}
	if got["level"] != "INFO" {
		t.Fatalf("level = %#v, want INFO", got["level"])
	}
	if got["component"] != "test" {
		t.Fatalf("component = %#v, want test", got["component"])
	}
	if got["port"] != float64(8080) {
		t.Fatalf("port = %#v, want 8080", got["port"])
	}
}

func TestNewJSONLoggerFiltersDebugByDefault(t *testing.T) {
	var buf bytes.Buffer
	logger := NewJSONLogger(&buf, false)

	logger.Debug("debug.hidden")
	logger.Info("info.visible")

	if strings.Contains(buf.String(), "debug.hidden") {
		t.Fatalf("debug log was emitted: %s", buf.String())
	}
	if !strings.Contains(buf.String(), "info.visible") {
		t.Fatalf("info log was not emitted: %s", buf.String())
	}
}

func TestNewJSONLoggerEnablesDebug(t *testing.T) {
	var buf bytes.Buffer
	logger := NewJSONLogger(&buf, true)

	logger.Debug("debug.visible")

	if !strings.Contains(buf.String(), "debug.visible") {
		t.Fatalf("debug log was not emitted: %s", buf.String())
	}
}

func TestContextLoggerRoundTrip(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(bytes.NewBuffer(nil), nil))
	ctx := WithContext(context.Background(), logger)

	if got := FromContext(ctx); got != logger {
		t.Fatalf("FromContext returned %#v, want %#v", got, logger)
	}
}
