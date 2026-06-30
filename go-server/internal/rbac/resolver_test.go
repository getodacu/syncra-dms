package rbac

import (
	"strings"
	"testing"
	"time"

	"ai.ro/syncra/dms/internal/auth"
	"ai.ro/syncra/dms/internal/orgunits"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestResolverDeniesInactiveUser(t *testing.T) {
	db := newResolverTestDB(t)
	user := createResolverUser(t, db, string(UserStatusInactive), nil)
	role := createResolverRoleWithPermission(t, db, "role.view")
	assignUserRole(t, db, user.ID, role.ID, ScopeGlobal, nil)

	allowed, err := NewResolver(db).Can(t.Context(), Check{
		UserID:     user.ID,
		Permission: "role.view",
	})
	if err != nil {
		t.Fatalf("Can() error = %v", err)
	}
	if allowed {
		t.Fatal("Can() = true, want false")
	}
}

func TestResolverEffectiveGrantsDeniesSoftDeletedActiveUser(t *testing.T) {
	db := newResolverTestDB(t)
	user := createResolverUser(t, db, string(UserStatusActive), nil)
	role := createResolverRoleWithPermission(t, db, "role.view")
	assignUserRole(t, db, user.ID, role.ID, ScopeGlobal, nil)
	softDeleteResolverUser(t, db, user.ID)

	grants, err := NewResolver(db).EffectiveGrants(t.Context(), user.ID)
	if err != nil {
		t.Fatalf("EffectiveGrants() error = %v", err)
	}
	if len(grants) != 0 {
		t.Fatalf("EffectiveGrants() = %#v, want no grants", grants)
	}
}

func TestResolverAllowsEmptyStatusNonDeletedUser(t *testing.T) {
	db := newResolverTestDB(t)
	user := createResolverUser(t, db, string(UserStatusActive), nil)
	role := createResolverRoleWithPermission(t, db, "role.view")
	assignUserRole(t, db, user.ID, role.ID, ScopeGlobal, nil)
	setResolverUserStatus(t, db, user.ID, "")

	allowed, err := NewResolver(db).Can(t.Context(), Check{
		UserID:     user.ID,
		Permission: "role.view",
	})
	if err != nil {
		t.Fatalf("Can() error = %v", err)
	}
	if !allowed {
		t.Fatal("Can() = false, want true")
	}
}

func TestResolverAllowsUserRoleWithGlobalScope(t *testing.T) {
	db := newResolverTestDB(t)
	user := createResolverUser(t, db, string(UserStatusActive), nil)
	role := createResolverRoleWithPermission(t, db, "role.view")
	assignUserRole(t, db, user.ID, role.ID, ScopeGlobal, nil)

	allowed, err := NewResolver(db).Can(t.Context(), Check{
		UserID:     user.ID,
		Permission: "role.view",
	})
	if err != nil {
		t.Fatalf("Can() error = %v", err)
	}
	if !allowed {
		t.Fatal("Can() = false, want true")
	}
}

func TestResolverAllowsGroupRoleThroughMembership(t *testing.T) {
	db := newResolverTestDB(t)
	user := createResolverUser(t, db, string(UserStatusActive), nil)
	group := createResolverGroupWithUser(t, db, user.ID)
	role := createResolverRoleWithPermission(t, db, "group.view")
	assignGroupRole(t, db, group.ID, role.ID, ScopeGlobal, nil)

	allowed, err := NewResolver(db).Can(t.Context(), Check{
		UserID:     user.ID,
		Permission: "group.view",
	})
	if err != nil {
		t.Fatalf("Can() error = %v", err)
	}
	if !allowed {
		t.Fatal("Can() = false, want true")
	}
}

func TestResolverAllowsOrganizationUnitRoleThroughPrimaryUnit(t *testing.T) {
	db := newResolverTestDB(t)
	unit := createResolverUnit(t, db, "Finance", nil)
	user := createResolverUser(t, db, string(UserStatusActive), &unit.ID)
	role := createResolverRoleWithPermission(t, db, "organization_unit.view")
	assignOrganizationUnitRole(t, db, unit.ID, role.ID, ScopeOrganizationUnit)

	allowed, err := NewResolver(db).Can(t.Context(), Check{
		UserID:             user.ID,
		Permission:         "organization_unit.view",
		OrganizationUnitID: &unit.ID,
	})
	if err != nil {
		t.Fatalf("Can() error = %v", err)
	}
	if !allowed {
		t.Fatal("Can() = false, want true")
	}
}

func TestResolverDeniesOrganizationUnitRoleFromArchivedPrimaryUnit(t *testing.T) {
	db := newResolverTestDB(t)
	unit := createResolverUnit(t, db, "Finance", nil)
	user := createResolverUser(t, db, string(UserStatusActive), &unit.ID)
	role := createResolverRoleWithPermission(t, db, "organization_unit.view")
	assignOrganizationUnitRole(t, db, unit.ID, role.ID, ScopeGlobal)
	archiveResolverUnit(t, db, unit.ID)

	allowed, err := NewResolver(db).Can(t.Context(), Check{
		UserID:     user.ID,
		Permission: "organization_unit.view",
	})
	if err != nil {
		t.Fatalf("Can() error = %v", err)
	}
	if allowed {
		t.Fatal("Can() = true, want false")
	}
}

func TestResolverEffectiveGrantsUsesAssignmentSourceValues(t *testing.T) {
	db := newResolverTestDB(t)
	unit := createResolverUnit(t, db, "Finance", nil)
	user := createResolverUser(t, db, string(UserStatusActive), &unit.ID)

	userRole := createResolverRoleWithPermission(t, db, "source.user.view")
	assignUserRole(t, db, user.ID, userRole.ID, ScopeGlobal, nil)

	group := createResolverGroupWithUser(t, db, user.ID)
	groupRole := createResolverRoleWithPermission(t, db, "source.group.view")
	assignGroupRole(t, db, group.ID, groupRole.ID, ScopeGlobal, nil)

	organizationUnitRole := createResolverRoleWithPermission(t, db, "source.organization_unit.view")
	assignOrganizationUnitRole(t, db, unit.ID, organizationUnitRole.ID, ScopeOrganizationUnit)

	grants, err := NewResolver(db).EffectiveGrants(t.Context(), user.ID)
	if err != nil {
		t.Fatalf("EffectiveGrants() error = %v", err)
	}

	sourceByPermission := make(map[string]string, len(grants))
	for _, grant := range grants {
		sourceByPermission[grant.PermissionCode] = grant.Source
	}
	want := map[string]string{
		"source.user.view":              "user_role",
		"source.group.view":             "group_role",
		"source.organization_unit.view": "organization_unit_role",
	}
	for permissionCode, wantSource := range want {
		if sourceByPermission[permissionCode] != wantSource {
			t.Fatalf("EffectiveGrants() source for %s = %q, want %q", permissionCode, sourceByPermission[permissionCode], wantSource)
		}
	}
}

func TestResolverAllowsOrganizationUnitAndChildrenForDescendantUnit(t *testing.T) {
	db := newResolverTestDB(t)
	parent := createResolverUnit(t, db, "Finance", nil)
	child := createResolverUnit(t, db, "Accounts Payable", &parent.ID)
	user := createResolverUser(t, db, string(UserStatusActive), nil)
	role := createResolverRoleWithPermission(t, db, "document.view")
	assignUserRole(t, db, user.ID, role.ID, ScopeOrganizationUnitAndChildren, &parent.ID)

	allowed, err := NewResolver(db).Can(t.Context(), Check{
		UserID:             user.ID,
		Permission:         "document.view",
		OrganizationUnitID: &child.ID,
	})
	if err != nil {
		t.Fatalf("Can() error = %v", err)
	}
	if !allowed {
		t.Fatal("Can() = false, want true")
	}
}

func TestResolverDeniesOrganizationUnitAndChildrenWhenGrantUnitArchived(t *testing.T) {
	db := newResolverTestDB(t)
	parent := createResolverUnit(t, db, "Finance", nil)
	child := createResolverUnit(t, db, "Accounts Payable", &parent.ID)
	user := createResolverUser(t, db, string(UserStatusActive), nil)
	role := createResolverRoleWithPermission(t, db, "document.view")
	assignUserRole(t, db, user.ID, role.ID, ScopeOrganizationUnitAndChildren, &parent.ID)
	archiveResolverUnit(t, db, parent.ID)

	allowed, err := NewResolver(db).Can(t.Context(), Check{
		UserID:             user.ID,
		Permission:         "document.view",
		OrganizationUnitID: &child.ID,
	})
	if err != nil {
		t.Fatalf("Can() error = %v", err)
	}
	if allowed {
		t.Fatal("Can() = true, want false")
	}
}

func TestResolverDeniesUnmatchedScope(t *testing.T) {
	db := newResolverTestDB(t)
	grantedUnit := createResolverUnit(t, db, "Finance", nil)
	requestedUnit := createResolverUnit(t, db, "Legal", nil)
	user := createResolverUser(t, db, string(UserStatusActive), nil)
	role := createResolverRoleWithPermission(t, db, "document.view")
	assignUserRole(t, db, user.ID, role.ID, ScopeOrganizationUnit, &grantedUnit.ID)

	allowed, err := NewResolver(db).Can(t.Context(), Check{
		UserID:             user.ID,
		Permission:         "document.view",
		OrganizationUnitID: &requestedUnit.ID,
	})
	if err != nil {
		t.Fatalf("Can() error = %v", err)
	}
	if allowed {
		t.Fatal("Can() = true, want false")
	}
}

func TestResolverDeniesArchivedExactScopeOrganizationUnit(t *testing.T) {
	db := newResolverTestDB(t)
	unit := createResolverUnit(t, db, "Finance", nil)
	user := createResolverUser(t, db, string(UserStatusActive), nil)
	role := createResolverRoleWithPermission(t, db, "document.view")
	assignUserRole(t, db, user.ID, role.ID, ScopeOrganizationUnit, &unit.ID)
	archiveResolverUnit(t, db, unit.ID)

	allowed, err := NewResolver(db).Can(t.Context(), Check{
		UserID:             user.ID,
		Permission:         "document.view",
		OrganizationUnitID: &unit.ID,
	})
	if err != nil {
		t.Fatalf("Can() error = %v", err)
	}
	if allowed {
		t.Fatal("Can() = true, want false")
	}
}

func TestResolverDeniesMissingUserWithoutError(t *testing.T) {
	db := newResolverTestDB(t)

	allowed, err := NewResolver(db).Can(t.Context(), Check{
		UserID:     uuid.NewString(),
		Permission: "role.view",
	})
	if err != nil {
		t.Fatalf("Can() error = %v", err)
	}
	if allowed {
		t.Fatal("Can() = true, want false")
	}
}

func TestResolverDeniesInactiveUserRole(t *testing.T) {
	db := newResolverTestDB(t)
	user := createResolverUser(t, db, string(UserStatusActive), nil)
	role := createResolverRoleWithPermission(t, db, "role.view")
	setResolverRoleActive(t, db, role.ID, false)
	assignUserRole(t, db, user.ID, role.ID, ScopeGlobal, nil)

	allowed, err := NewResolver(db).Can(t.Context(), Check{
		UserID:     user.ID,
		Permission: "role.view",
	})
	if err != nil {
		t.Fatalf("Can() error = %v", err)
	}
	if allowed {
		t.Fatal("Can() = true, want false")
	}
}

func TestResolverDeniesInactiveGroupRole(t *testing.T) {
	db := newResolverTestDB(t)
	user := createResolverUser(t, db, string(UserStatusActive), nil)
	group := createResolverGroupWithUser(t, db, user.ID)
	setResolverGroupActive(t, db, group.ID, false)
	role := createResolverRoleWithPermission(t, db, "group.view")
	assignGroupRole(t, db, group.ID, role.ID, ScopeGlobal, nil)

	allowed, err := NewResolver(db).Can(t.Context(), Check{
		UserID:     user.ID,
		Permission: "group.view",
	})
	if err != nil {
		t.Fatalf("Can() error = %v", err)
	}
	if allowed {
		t.Fatal("Can() = true, want false")
	}
}

func TestResolverDeniesScopedGrantWithoutRequestedOrganizationUnit(t *testing.T) {
	db := newResolverTestDB(t)
	unit := createResolverUnit(t, db, "Finance", nil)
	user := createResolverUser(t, db, string(UserStatusActive), nil)
	role := createResolverRoleWithPermission(t, db, "document.view")
	assignUserRole(t, db, user.ID, role.ID, ScopeOrganizationUnit, &unit.ID)

	allowed, err := NewResolver(db).Can(t.Context(), Check{
		UserID:     user.ID,
		Permission: "document.view",
	})
	if err != nil {
		t.Fatalf("Can() error = %v", err)
	}
	if allowed {
		t.Fatal("Can() = true, want false")
	}
}

func TestResolverDeniesOrganizationUnitRoleWhenUserHasNoPrimaryUnit(t *testing.T) {
	db := newResolverTestDB(t)
	unit := createResolverUnit(t, db, "Finance", nil)
	user := createResolverUser(t, db, string(UserStatusActive), nil)
	role := createResolverRoleWithPermission(t, db, "organization_unit.view")
	assignOrganizationUnitRole(t, db, unit.ID, role.ID, ScopeOrganizationUnit)

	allowed, err := NewResolver(db).Can(t.Context(), Check{
		UserID:             user.ID,
		Permission:         "organization_unit.view",
		OrganizationUnitID: &unit.ID,
	})
	if err != nil {
		t.Fatalf("Can() error = %v", err)
	}
	if allowed {
		t.Fatal("Can() = true, want false")
	}
}

func TestResolverPropagatesDatabaseErrors(t *testing.T) {
	db := newResolverTestDB(t)
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("DB() error = %v", err)
	}
	if err := sqlDB.Close(); err != nil {
		t.Fatalf("close db: %v", err)
	}

	allowed, err := NewResolver(db).Can(t.Context(), Check{
		UserID:     uuid.NewString(),
		Permission: "role.view",
	})
	if err == nil {
		t.Fatal("Can() error = nil, want database error")
	}
	if allowed {
		t.Fatal("Can() = true, want false")
	}
}

func newResolverTestDB(t *testing.T) *gorm.DB {
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
		&Group{},
		&GroupUser{},
		&GroupRole{},
		&OrganizationUnitRole{},
	); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	return db
}

func createResolverUser(t *testing.T, db *gorm.DB, status string, primaryOrganizationUnitID *string) auth.User {
	t.Helper()
	now := time.Now().UTC()
	user := auth.User{
		Name:                      "Resolver User",
		Email:                     "resolver-" + uuid.NewString() + "@example.com",
		EmailVerified:             status == string(UserStatusActive),
		Status:                    status,
		PrimaryOrganizationUnitID: primaryOrganizationUnitID,
		CreatedAt:                 now,
		UpdatedAt:                 now,
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	return user
}

func setResolverUserStatus(t *testing.T, db *gorm.DB, userID string, status string) {
	t.Helper()
	if status == "" {
		if err := db.Exec("PRAGMA ignore_check_constraints = ON").Error; err != nil {
			t.Fatalf("disable sqlite check constraints: %v", err)
		}
	}
	if err := db.Model(&auth.User{}).Where("id = ?", userID).Update("status", status).Error; err != nil {
		t.Fatalf("set user status: %v", err)
	}
}

func softDeleteResolverUser(t *testing.T, db *gorm.DB, userID string) {
	t.Helper()
	if err := db.Model(&auth.User{}).Where("id = ?", userID).Update("deleted_at", time.Now().UTC()).Error; err != nil {
		t.Fatalf("soft delete user: %v", err)
	}
}

func createResolverRoleWithPermission(t *testing.T, db *gorm.DB, permissionCode string) Role {
	t.Helper()
	now := time.Now().UTC()
	role := Role{
		Name:      "Resolver Role " + uuid.NewString(),
		Code:      "resolver_role_" + uuid.NewString(),
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}

	permission := Permission{
		Code:      permissionCode,
		Name:      permissionCode,
		Category:  "Resolver",
		IsSystem:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := db.Create(&permission).Error; err != nil {
		t.Fatalf("create permission: %v", err)
	}

	rolePermission := RolePermission{
		RoleID:       role.ID,
		PermissionID: permission.ID,
		CreatedAt:    now,
	}
	if err := db.Create(&rolePermission).Error; err != nil {
		t.Fatalf("create role permission: %v", err)
	}
	return role
}

func assignUserRole(t *testing.T, db *gorm.DB, userID, roleID string, scope ScopeType, organizationUnitID *string) {
	t.Helper()
	now := time.Now().UTC()
	userRole := UserRole{
		UserID:             userID,
		RoleID:             roleID,
		ScopeType:          scope,
		OrganizationUnitID: organizationUnitID,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	if err := db.Create(&userRole).Error; err != nil {
		t.Fatalf("assign user role: %v", err)
	}
}

func setResolverRoleActive(t *testing.T, db *gorm.DB, roleID string, active bool) {
	t.Helper()
	if err := db.Model(&Role{}).Where("id = ?", roleID).Update("is_active", active).Error; err != nil {
		t.Fatalf("update role active flag: %v", err)
	}
}

func createResolverGroupWithUser(t *testing.T, db *gorm.DB, userID string) Group {
	t.Helper()
	now := time.Now().UTC()
	group := Group{
		Name:      "Resolver Group " + uuid.NewString(),
		Code:      "resolver_group_" + uuid.NewString(),
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := db.Create(&group).Error; err != nil {
		t.Fatalf("create group: %v", err)
	}

	groupUser := GroupUser{
		GroupID:   group.ID,
		UserID:    userID,
		CreatedAt: now,
	}
	if err := db.Create(&groupUser).Error; err != nil {
		t.Fatalf("create group user: %v", err)
	}
	return group
}

func setResolverGroupActive(t *testing.T, db *gorm.DB, groupID string, active bool) {
	t.Helper()
	if err := db.Model(&Group{}).Where("id = ?", groupID).Update("is_active", active).Error; err != nil {
		t.Fatalf("update group active flag: %v", err)
	}
}

func assignGroupRole(t *testing.T, db *gorm.DB, groupID, roleID string, scope ScopeType, organizationUnitID *string) {
	t.Helper()
	groupRole := GroupRole{
		GroupID:            groupID,
		RoleID:             roleID,
		ScopeType:          scope,
		OrganizationUnitID: organizationUnitID,
		CreatedAt:          time.Now().UTC(),
	}
	if err := db.Create(&groupRole).Error; err != nil {
		t.Fatalf("assign group role: %v", err)
	}
}

func assignOrganizationUnitRole(t *testing.T, db *gorm.DB, organizationUnitID, roleID string, scope ScopeType) {
	t.Helper()
	organizationUnitRole := OrganizationUnitRole{
		OrganizationUnitID: organizationUnitID,
		RoleID:             roleID,
		ScopeType:          scope,
		CreatedAt:          time.Now().UTC(),
	}
	if err := db.Create(&organizationUnitRole).Error; err != nil {
		t.Fatalf("assign organization unit role: %v", err)
	}
}

func createResolverUnit(t *testing.T, db *gorm.DB, name string, parentID *string) orgunits.Unit {
	t.Helper()
	now := time.Now().UTC()
	code := "UNIT_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	unit := orgunits.Unit{
		ParentID:  parentID,
		Name:      name,
		Code:      &code,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := db.Create(&unit).Error; err != nil {
		t.Fatalf("create organization unit: %v", err)
	}
	return unit
}

func archiveResolverUnit(t *testing.T, db *gorm.DB, unitID string) {
	t.Helper()
	if err := db.Model(&orgunits.Unit{}).Where("id = ?", unitID).Update("archived_at", time.Now().UTC()).Error; err != nil {
		t.Fatalf("archive organization unit: %v", err)
	}
}
