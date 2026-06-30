package rbac

import (
	"strings"
	"testing"
	"time"

	"ai.ro/syncra/dms/internal/auth"
	"ai.ro/syncra/dms/internal/orgunits"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestBootstrapLegacyAdminsAssignsSystemAdministratorRole(t *testing.T) {
	db := newBootstrapTestDB(t)
	if err := SeedDefaults(db); err != nil {
		t.Fatalf("seed: %v", err)
	}
	admin := auth.User{Name: "Admin", Email: "admin@example.com", EmailVerified: true, Role: auth.UserRoleAdmin, Status: string(UserStatusActive), CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := db.Create(&admin).Error; err != nil {
		t.Fatalf("create admin: %v", err)
	}
	var adminRole Role
	if err := db.First(&adminRole, "code = ?", SystemAdministratorRoleCode).Error; err != nil {
		t.Fatalf("load system administrator role: %v", err)
	}

	if err := BootstrapLegacyAdmins(db); err != nil {
		t.Fatalf("bootstrap: %v", err)
	}
	if err := BootstrapLegacyAdmins(db); err != nil {
		t.Fatalf("bootstrap second: %v", err)
	}

	var count int64
	if err := db.Model(&UserRole{}).Where("user_id = ? AND role_id = ? AND scope_type = ?", admin.ID, adminRole.ID, ScopeGlobal).Count(&count).Error; err != nil {
		t.Fatalf("count user roles: %v", err)
	}
	if count != 1 {
		t.Fatalf("user role count = %d, want 1", count)
	}
}

func TestBootstrapLegacyAdminsAssignsEmptyStatusLegacyAdmin(t *testing.T) {
	db := newBootstrapTestDB(t)
	if err := SeedDefaults(db); err != nil {
		t.Fatalf("seed: %v", err)
	}
	admin := auth.User{Name: "Legacy Admin", Email: "legacy-admin@example.com", EmailVerified: true, Role: auth.UserRoleAdmin, Status: string(UserStatusActive), CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := db.Create(&admin).Error; err != nil {
		t.Fatalf("create admin: %v", err)
	}
	if err := db.Exec("PRAGMA ignore_check_constraints = ON").Error; err != nil {
		t.Fatalf("disable sqlite check constraints: %v", err)
	}
	if err := db.Model(&auth.User{}).Where("id = ?", admin.ID).Update("status", "").Error; err != nil {
		t.Fatalf("clear admin status: %v", err)
	}
	var adminRole Role
	if err := db.First(&adminRole, "code = ?", SystemAdministratorRoleCode).Error; err != nil {
		t.Fatalf("load system administrator role: %v", err)
	}

	if err := BootstrapLegacyAdmins(db); err != nil {
		t.Fatalf("bootstrap: %v", err)
	}

	var count int64
	if err := db.Model(&UserRole{}).Where("user_id = ? AND role_id = ? AND scope_type = ?", admin.ID, adminRole.ID, ScopeGlobal).Count(&count).Error; err != nil {
		t.Fatalf("count user roles: %v", err)
	}
	if count != 1 {
		t.Fatalf("user role count = %d, want 1", count)
	}
}

func TestBootstrapLegacyAdminsDoesNotAssignNonAdminUsers(t *testing.T) {
	db := newBootstrapTestDB(t)
	if err := SeedDefaults(db); err != nil {
		t.Fatalf("seed: %v", err)
	}
	now := time.Now().UTC()
	user := auth.User{
		Name:          "User",
		Email:         "user@example.com",
		EmailVerified: true,
		Role:          auth.UserRoleUser,
		Status:        string(UserStatusActive),
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}

	if err := BootstrapLegacyAdmins(db); err != nil {
		t.Fatalf("bootstrap: %v", err)
	}

	var count int64
	if err := db.Model(&UserRole{}).Where("user_id = ?", user.ID).Count(&count).Error; err != nil {
		t.Fatalf("count user roles: %v", err)
	}
	if count != 0 {
		t.Fatalf("user role count = %d, want 0", count)
	}
}

func TestBootstrapLegacyAdminsSkipsInactiveSuspendedAndDeletedAdmins(t *testing.T) {
	db := newBootstrapTestDB(t)
	if err := SeedDefaults(db); err != nil {
		t.Fatalf("seed: %v", err)
	}
	now := time.Now().UTC()
	statuses := []UserStatus{UserStatusInactive, UserStatusSuspended, UserStatusDeleted}
	for _, status := range statuses {
		user := auth.User{
			Name:          "Admin " + string(status),
			Email:         string(status) + "@example.com",
			EmailVerified: true,
			Role:          auth.UserRoleAdmin,
			Status:        string(status),
			CreatedAt:     now,
			UpdatedAt:     now,
		}
		if status == UserStatusDeleted {
			user.DeletedAt = &now
		}
		if err := db.Create(&user).Error; err != nil {
			t.Fatalf("create %s admin: %v", status, err)
		}
	}

	if err := BootstrapLegacyAdmins(db); err != nil {
		t.Fatalf("bootstrap: %v", err)
	}

	var count int64
	if err := db.Model(&UserRole{}).Count(&count).Error; err != nil {
		t.Fatalf("count user roles: %v", err)
	}
	if count != 0 {
		t.Fatalf("user role count = %d, want 0", count)
	}
}

func TestBootstrapLegacyAdminsSkipsSoftDeletedActiveAdmin(t *testing.T) {
	db := newBootstrapTestDB(t)
	if err := SeedDefaults(db); err != nil {
		t.Fatalf("seed: %v", err)
	}
	now := time.Now().UTC()
	admin := auth.User{
		Name:          "Soft Deleted Admin",
		Email:         "soft-deleted-admin@example.com",
		EmailVerified: true,
		Role:          auth.UserRoleAdmin,
		Status:        string(UserStatusActive),
		DeletedAt:     &now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := db.Create(&admin).Error; err != nil {
		t.Fatalf("create soft deleted admin: %v", err)
	}

	if err := BootstrapLegacyAdmins(db); err != nil {
		t.Fatalf("bootstrap: %v", err)
	}

	var count int64
	if err := db.Model(&UserRole{}).Where("user_id = ?", admin.ID).Count(&count).Error; err != nil {
		t.Fatalf("count user roles: %v", err)
	}
	if count != 0 {
		t.Fatalf("user role count = %d, want 0", count)
	}
}

func TestBootstrapLegacyAdminsReportsMissingSystemAdministratorRole(t *testing.T) {
	err := BootstrapLegacyAdmins(newBootstrapTestDB(t))
	if err == nil {
		t.Fatal("bootstrap error = nil, want missing system administrator role error")
	}
	if !strings.Contains(err.Error(), SystemAdministratorRoleCode) {
		t.Fatalf("bootstrap error = %q, want role code %q", err, SystemAdministratorRoleCode)
	}
}

func newBootstrapTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	name := strings.NewReplacer("/", "_", " ", "_").Replace(t.Name())
	db, err := gorm.Open(sqlite.Open("file:"+name+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(
		&orgunits.Unit{},
		&auth.User{},
		&Role{},
		&Permission{},
		&RolePermission{},
		&UserRole{},
	); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	return db
}
