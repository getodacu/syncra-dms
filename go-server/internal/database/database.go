package database

import (
	"errors"
	"strings"

	"ai.ro/syncra/dms/internal/auth"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenPostgres(dsn string) (*gorm.DB, error) {
	if strings.TrimSpace(dsn) == "" {
		return nil, errors.New("DSN is required")
	}
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func ApplicationModels() []any {
	return []any{
		&auth.User{},
		&auth.AuthAccount{},
		&auth.Session{},
		&auth.Verification{},
	}
}
