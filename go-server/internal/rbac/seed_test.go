package rbac

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSeedDefaultsIsIdempotent(t *testing.T) {
	db := newSeedTestDB(t)

	if err := SeedDefaults(db); err != nil {
		t.Fatalf("SeedDefaults first: %v", err)
	}
	if err := SeedDefaults(db); err != nil {
		t.Fatalf("SeedDefaults second: %v", err)
	}

	var permissionCount int64
	if err := db.Model(&Permission{}).Count(&permissionCount).Error; err != nil {
		t.Fatalf("count permissions: %v", err)
	}
	if permissionCount < 1 {
		t.Fatal("expected seeded permissions")
	}

	var admin Role
	if err := db.First(&admin, "code = ?", SystemAdministratorRoleCode).Error; err != nil {
		t.Fatalf("load system administrator role: %v", err)
	}

	var adminPermissionCount int64
	if err := db.Model(&RolePermission{}).Where("role_id = ?", admin.ID).Count(&adminPermissionCount).Error; err != nil {
		t.Fatalf("count admin permissions: %v", err)
	}
	if adminPermissionCount != permissionCount {
		t.Fatalf("admin permission count = %d, want %d", adminPermissionCount, permissionCount)
	}
}

func newSeedTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&Role{}, &Permission{}, &RolePermission{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	return db
}
