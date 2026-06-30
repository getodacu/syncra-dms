package api

import (
	"encoding/json"
	"testing"
)

func TestValidateJSONSchemaAcceptsValidSchemas(t *testing.T) {
	tests := []struct {
		name   string
		schema string
	}{
		{
			name:   "draft 2020-12 default",
			schema: `{"type":"array","prefixItems":[{"type":"string"}]}`,
		},
		{
			name:   "declared draft 7",
			schema: `{"$schema":"http://json-schema.org/draft-07/schema#","type":"array","items":[{"type":"string"}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateJSONSchema(json.RawMessage(tt.schema)); err != nil {
				t.Fatalf("validateJSONSchema returned error: %v", err)
			}
		})
	}
}

func TestValidateJSONSchemaRejectsInvalidSchemas(t *testing.T) {
	tests := []struct {
		name   string
		schema string
	}{
		{
			name:   "invalid type keyword",
			schema: `{"type":"strung"}`,
		},
		{
			name:   "invalid required keyword",
			schema: `{"type":"object","required":"email"}`,
		},
		{
			name:   "invalid nested property schema",
			schema: `{"type":"object","properties":{"email":{"type":"strung"}}}`,
		},
		{
			name:   "unsupported schema dialect",
			schema: `{"$schema":"https://example.test/schema","type":"object"}`,
		},
		{
			name:   "unresolved external reference",
			schema: `{"$ref":"https://example.test/external-schema.json"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateJSONSchema(json.RawMessage(tt.schema)); err == nil {
				t.Fatal("validateJSONSchema returned nil error")
			}
		})
	}
}
