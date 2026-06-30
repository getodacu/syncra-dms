package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunWritesEmptySchemaForLeanScaffold(t *testing.T) {
	var stdout bytes.Buffer

	if err := run(&stdout); err != nil {
		t.Fatalf("run() error = %v", err)
	}
	if got := strings.TrimSpace(stdout.String()); got != "" {
		t.Fatalf("run() output = %q, want empty schema before domain models are added", got)
	}
}
