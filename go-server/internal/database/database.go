package database

import (
	"errors"
	"strings"

	"ai.ro/syncra/dms/internal/auth"
	"ai.ro/syncra/dms/internal/documents"
	"ai.ro/syncra/dms/internal/orgunits"
	"ai.ro/syncra/dms/internal/rbac"
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
		&orgunits.Unit{},
		&documents.Folder{},
		&documents.Document{},
		&rbac.Role{},
		&rbac.Permission{},
		&rbac.RolePermission{},
		&rbac.UserRole{},
		&rbac.Group{},
		&rbac.GroupUser{},
		&rbac.GroupRole{},
		&rbac.OrganizationUnitRole{},
		&rbac.BootstrapMarker{},
	}
}
