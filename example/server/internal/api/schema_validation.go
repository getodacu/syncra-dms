package api

import (
	"bytes"
	"encoding/json"

	jsonschema "github.com/santhosh-tekuri/jsonschema/v6"
)

const schemaValidationResourceURL = "schema.json"

func validateJSONSchema(raw json.RawMessage) error {
	doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(raw))
	if err != nil {
		return err
	}

	compiler := jsonschema.NewCompiler()
	compiler.DefaultDraft(jsonschema.Draft2020)
	if err := compiler.AddResource(schemaValidationResourceURL, doc); err != nil {
		return err
	}
	_, err = compiler.Compile(schemaValidationResourceURL)
	return err
}
