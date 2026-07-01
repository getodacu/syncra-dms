package rbac

import (
	"sort"
	"strings"
	"testing"
	"time"

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

func TestDocumentPermissionsAreSeeded(t *testing.T) {
	db := newSeedTestDB(t)
	if err := SeedDefaults(db); err != nil {
		t.Fatalf("SeedDefaults error = %v", err)
	}
	for _, code := range []string{"document.view", "document.create", "document.update", "document.delete", "document.download"} {
		var permission Permission
		if err := db.First(&permission, "code = ?", code).Error; err != nil {
			t.Fatalf("permission %s was not seeded: %v", code, err)
		}
		if permission.Category != "Document Repository" {
			t.Fatalf("%s category = %q", code, permission.Category)
		}
	}
}

func TestSeedDefaultsReconcilesDefaultRolePermissions(t *testing.T) {
	db := newSeedTestDB(t)
	if err := SeedDefaults(db); err != nil {
		t.Fatalf("SeedDefaults first: %v", err)
	}

	var orgAdmin Role
	if err := db.First(&orgAdmin, "code = ?", OrganizationAdministratorRoleCode).Error; err != nil {
		t.Fatalf("load organization administrator role: %v", err)
	}
	var systemAdmin Permission
	if err := db.First(&systemAdmin, "code = ?", "system.admin").Error; err != nil {
		t.Fatalf("load system admin permission: %v", err)
	}
	extra := RolePermission{
		RoleID:       orgAdmin.ID,
		PermissionID: systemAdmin.ID,
		CreatedAt:    time.Now().UTC(),
	}
	if err := db.Create(&extra).Error; err != nil {
		t.Fatalf("create extra role permission: %v", err)
	}

	if err := SeedDefaults(db); err != nil {
		t.Fatalf("SeedDefaults second: %v", err)
	}

	assertRolePermissionCodes(t, db, OrganizationAdministratorRoleCode, organizationAdministratorExpectedPermissionCodes())
}

func TestSeedDefaultsFailsForUnknownRolePermissionCode(t *testing.T) {
	originalRegistry := PermissionRegistry
	PermissionRegistry = []PermissionDefinition{
		{Code: "system.admin", Name: "System administration", Category: "System"},
		{Code: "organization_unit.view", Name: "View organization units", Category: "Organization Unit Management"},
		{Code: "organization_unit.update", Name: "Update organization units", Category: "Organization Unit Management"},
	}
	t.Cleanup(func() {
		PermissionRegistry = originalRegistry
	})

	err := SeedDefaults(newSeedTestDB(t))
	if err == nil {
		t.Fatal("SeedDefaults error = nil, want unknown permission code error")
	}
	if !strings.Contains(err.Error(), "unknown permission code") {
		t.Fatalf("SeedDefaults error = %q, want unknown permission code details", err)
	}
}

func TestSeedDefaultsDoesNotGrantFuturePermissionsToOrganizationAdministrator(t *testing.T) {
	originalRegistry := PermissionRegistry
	PermissionRegistry = append(append([]PermissionDefinition(nil), PermissionRegistry...), PermissionDefinition{
		Code:     "workflow.approve",
		Name:     "Approve workflows",
		Category: "Workflow Management",
	})
	t.Cleanup(func() {
		PermissionRegistry = originalRegistry
	})

	db := newSeedTestDB(t)
	if err := SeedDefaults(db); err != nil {
		t.Fatalf("SeedDefaults: %v", err)
	}

	assertRolePermissionCodes(t, db, OrganizationAdministratorRoleCode, organizationAdministratorExpectedPermissionCodes())
	assertRolePermissionCodes(t, db, SystemAdministratorRoleCode, PermissionCodes())
}

func TestSeedDefaultsAssignsExactDefaultRolePermissions(t *testing.T) {
	db := newSeedTestDB(t)
	if err := SeedDefaults(db); err != nil {
		t.Fatalf("SeedDefaults: %v", err)
	}

	tests := []struct {
		name     string
		roleCode string
		want     []string
	}{
		{
			name:     "system administrator",
			roleCode: SystemAdministratorRoleCode,
			want:     PermissionCodes(),
		},
		{
			name:     "organization administrator",
			roleCode: OrganizationAdministratorRoleCode,
			want:     organizationAdministratorExpectedPermissionCodes(),
		},
		{
			name:     "unit manager",
			roleCode: UnitManagerRoleCode,
			want: []string{
				"organization_unit.view",
				"organization_unit.update",
				"organization_unit.manage_users",
			},
		},
		{
			name:     "viewer",
			roleCode: ViewerRoleCode,
			want: []string{
				"organization_unit.view",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assertRolePermissionCodes(t, db, tc.roleCode, tc.want)
		})
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

func organizationAdministratorExpectedPermissionCodes() []string {
	return []string{
		"user.view",
		"user.create",
		"user.update",
		"user.delete",
		"user.activate",
		"user.suspend",
		"user.assign_role",
		"user.assign_group",
		"user.assign_unit",
		"role.view",
		"role.create",
		"role.update",
		"role.delete",
		"role.assign_permissions",
		"role.assign_users",
		"group.view",
		"group.create",
		"group.update",
		"group.delete",
		"group.manage_users",
		"group.assign_roles",
		"organization_unit.view",
		"organization_unit.create",
		"organization_unit.update",
		"organization_unit.delete",
		"organization_unit.manage_users",
		"organization_unit.manage_roles",
		"organization_unit.manage_permissions",
		"organization_unit.manage_hierarchy",
		"organization_unit.view_audit",
		"document.view",
		"document.create",
		"document.update",
		"document.delete",
		"document.download",
	}
}

func assertRolePermissionCodes(t *testing.T, db *gorm.DB, roleCode string, want []string) {
	t.Helper()
	var got []string
	if err := db.Table("role_permissions").
		Select("permissions.code").
		Joins("JOIN roles ON roles.id = role_permissions.role_id").
		Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
		Where("roles.code = ?", roleCode).
		Order("permissions.code").
		Scan(&got).Error; err != nil {
		t.Fatalf("load role permission codes: %v", err)
	}

	want = append([]string(nil), want...)
	sort.Strings(want)
	if !sameStrings(got, want) {
		t.Fatalf("%s permissions = %v, want %v", roleCode, got, want)
	}
}

func sameStrings(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}
