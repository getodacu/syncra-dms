package dbschema

import (
	"ariga.io/atlas-provider-gorm/gormschema"

	"ai.ro/syncra/dms/internal/database"
)

func LoadPostgresSchema() (string, error) {
	return gormschema.New("postgres").Load(database.ApplicationModels()...)
}
