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
