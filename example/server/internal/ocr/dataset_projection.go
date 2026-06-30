package ocr

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

const MaxDatasetFieldKeyCharacters = 120
const MaxDatasetFieldLabelCharacters = 160

type DatasetField struct {
	Path  string `json:"path"`
	Key   string `json:"key"`
	Label string `json:"label"`
}

func ValidateDatasetFields(schema json.RawMessage, fields []DatasetField) error {
	if len(fields) == 0 {
		return errors.New("selected_fields is required")
	}
	var schemaValue any
	if err := json.Unmarshal(schema, &schemaValue); err != nil {
		return fmt.Errorf("invalid schema JSON: %w", err)
	}
	seenKeys := make(map[string]struct{}, len(fields))
	for _, field := range fields {
		if err := validateDatasetFieldShape(field); err != nil {
			return err
		}
		if _, ok := seenKeys[field.Key]; ok {
			return errors.New("duplicate field key")
		}
		seenKeys[field.Key] = struct{}{}
		if !datasetSchemaPathExists(schemaValue, field.Path) {
			return fmt.Errorf("field path does not exist: %s", field.Path)
		}
	}
	return nil
}

func validateDatasetFieldShape(field DatasetField) error {
	if strings.TrimSpace(field.Path) == "" {
		return errors.New("field path is required")
	}
	if _, err := parseJSONPointer(field.Path); err != nil {
		return fmt.Errorf("invalid field path: %w", err)
	}
	if strings.TrimSpace(field.Key) == "" {
		return errors.New("field key is required")
	}
	if utf8.RuneCountInString(field.Key) > MaxDatasetFieldKeyCharacters {
		return fmt.Errorf("field key exceeds %d characters", MaxDatasetFieldKeyCharacters)
	}
	if strings.TrimSpace(field.Label) == "" {
		return errors.New("field label is required")
	}
	if utf8.RuneCountInString(field.Label) > MaxDatasetFieldLabelCharacters {
		return fmt.Errorf("field label exceeds %d characters", MaxDatasetFieldLabelCharacters)
	}
	return nil
}

func datasetSchemaPathExists(schemaValue any, path string) bool {
	parts, err := parseJSONPointer(path)
	if err != nil || len(parts) == 0 {
		return false
	}

	current := schemaValue
	for _, part := range parts {
		schemaObject, ok := current.(map[string]any)
		if !ok {
			return false
		}
		if schemaObject["type"] == "array" {
			return false
		}
		properties, ok := schemaObject["properties"].(map[string]any)
		if !ok {
			return false
		}
		next, ok := properties[part]
		if !ok {
			return false
		}
		current = next
	}
	return true
}

func parseJSONPointer(path string) ([]string, error) {
	if path == "" {
		return nil, nil
	}
	if !strings.HasPrefix(path, "/") {
		return nil, errors.New("JSON pointer must start with /")
	}

	rawParts := strings.Split(path[1:], "/")
	parts := make([]string, len(rawParts))
	for i, rawPart := range rawParts {
		var builder strings.Builder
		for j := 0; j < len(rawPart); j++ {
			if rawPart[j] != '~' {
				builder.WriteByte(rawPart[j])
				continue
			}
			if j+1 >= len(rawPart) {
				return nil, errors.New("invalid JSON pointer escape")
			}
			switch rawPart[j+1] {
			case '0':
				builder.WriteByte('~')
			case '1':
				builder.WriteByte('/')
			default:
				return nil, errors.New("invalid JSON pointer escape")
			}
			j++
		}
		parts[i] = builder.String()
	}
	return parts, nil
}

func jsonPointerLookup(value any, path string) (any, bool, error) {
	parts, err := parseJSONPointer(path)
	if err != nil {
		return nil, false, err
	}

	current := value
	for _, part := range parts {
		switch typed := current.(type) {
		case map[string]any:
			next, ok := typed[part]
			if !ok {
				return nil, false, nil
			}
			current = next
		case []any:
			index, err := strconv.Atoi(part)
			if err != nil || index < 0 || index >= len(typed) {
				return nil, false, nil
			}
			current = typed[index]
		default:
			return nil, false, nil
		}
	}
	return current, true, nil
}

func ProjectDatasetValues(annotation json.RawMessage, fields []DatasetField) (map[string]any, error) {
	var annotationValue any
	if err := json.Unmarshal(annotation, &annotationValue); err != nil {
		return nil, fmt.Errorf("invalid annotation JSON: %w", err)
	}

	values := make(map[string]any, len(fields))
	for _, field := range fields {
		if err := validateDatasetFieldShape(field); err != nil {
			return nil, err
		}
		value, ok, err := jsonPointerLookup(annotationValue, field.Path)
		if err != nil {
			return nil, fmt.Errorf("invalid field path: %w", err)
		}
		cellValue, err := formatDatasetCellValue(value, ok)
		if err != nil {
			return nil, fmt.Errorf("format field %s: %w", field.Key, err)
		}
		values[field.Key] = cellValue
	}
	return values, nil
}

func formatDatasetCellValue(value any, ok bool) (any, error) {
	if !ok || value == nil {
		return "", nil
	}

	switch value.(type) {
	case map[string]any, []any:
		encoded, err := json.Marshal(value)
		if err != nil {
			return nil, err
		}
		return string(encoded), nil
	default:
		return value, nil
	}
}
