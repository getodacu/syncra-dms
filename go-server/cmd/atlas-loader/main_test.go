package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunWritesAuthSchema(t *testing.T) {
	var stdout bytes.Buffer

	if err := run(&stdout); err != nil {
		t.Fatalf("run() error = %v", err)
	}
	got := stdout.String()
	for _, want := range []string{
		`CREATE TABLE "user"`,
		`CREATE TABLE "account"`,
		`CREATE TABLE "session"`,
		`CREATE TABLE "verification"`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("run() output missing %q:\n%s", want, got)
		}
	}
}
