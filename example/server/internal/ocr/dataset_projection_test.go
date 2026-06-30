package ocr

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestValidateDatasetFieldsAgainstSchema(t *testing.T) {
	schema := json.RawMessage(`{
		"type":"object",
		"properties":{
			"supplier":{"type":"object","properties":{"name":{"type":"string"}}},
			"total":{"type":"number"},
			"line_items":{"type":"array","items":{"type":"object","properties":{"description":{"type":"string"}}}},
			"a/b":{"type":"string"},
			"c~d":{"type":"string"}
		}
	}`)
	fields := []DatasetField{
		{Path: "/supplier/name", Key: "supplier_name", Label: "Supplier name"},
		{Path: "/total", Key: "total", Label: "Total"},
		{Path: "/line_items", Key: "line_items", Label: "Line items"},
		{Path: "/a~1b", Key: "slash", Label: "Slash"},
		{Path: "/c~0d", Key: "tilde", Label: "Tilde"},
	}

	if err := ValidateDatasetFields(schema, fields); err != nil {
		t.Fatalf("ValidateDatasetFields: %v", err)
	}
}

func TestValidateDatasetFieldsRejectsInvalidPathAndDuplicateKey(t *testing.T) {
	schema := json.RawMessage(`{"type":"object","properties":{"total":{"type":"number"}}}`)

	if err := ValidateDatasetFields(schema, []DatasetField{{Path: "/missing", Key: "missing", Label: "Missing"}}); err == nil {
		t.Fatal("ValidateDatasetFields invalid path error = nil")
	}
	if err := ValidateDatasetFields(schema, []DatasetField{{Path: "/bad~2path", Key: "bad", Label: "Bad"}}); err == nil {
		t.Fatal("ValidateDatasetFields invalid escape error = nil")
	}
	if err := ValidateDatasetFields(schema, []DatasetField{
		{Path: "/total", Key: "total", Label: "Total"},
		{Path: "/total", Key: "total", Label: "Total duplicate"},
	}); err == nil {
		t.Fatal("ValidateDatasetFields duplicate key error = nil")
	}
}

func TestProjectDatasetRowFormatsCells(t *testing.T) {
	annotation := json.RawMessage(`{
		"supplier":{"name":"Acme"},
		"total":123.45,
		"paid":true,
		"line_items":[{"description":"Work"}],
		"missing_null":null
	}`)
	fields := []DatasetField{
		{Path: "/supplier/name", Key: "supplier_name", Label: "Supplier name"},
		{Path: "/supplier", Key: "supplier", Label: "Supplier"},
		{Path: "/total", Key: "total", Label: "Total"},
		{Path: "/paid", Key: "paid", Label: "Paid"},
		{Path: "/line_items", Key: "line_items", Label: "Line items"},
		{Path: "/missing", Key: "missing", Label: "Missing"},
		{Path: "/missing_null", Key: "missing_null", Label: "Missing null"},
	}

	got, err := ProjectDatasetValues(annotation, fields)
	if err != nil {
		t.Fatalf("ProjectDatasetValues: %v", err)
	}
	want := map[string]any{
		"supplier_name": "Acme",
		"supplier":      `{"name":"Acme"}`,
		"total":         float64(123.45),
		"paid":          true,
		"line_items":    `[{"description":"Work"}]`,
		"missing":       "",
		"missing_null":  "",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("values = %#v, want %#v", got, want)
	}
}
