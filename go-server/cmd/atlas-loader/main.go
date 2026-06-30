package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"

	"ai.ro/syncra/dms/internal/database"
)

func main() {
	if err := run(os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
}

func run(stdout io.Writer) error {
	schema, err := loadSchema()
	if err != nil {
		return err
	}
	_, err = io.WriteString(stdout, schema)
	return err
}

func loadSchema() (string, error) {
	return gormschema.New("postgres").Load(database.ApplicationModels()...)
}
