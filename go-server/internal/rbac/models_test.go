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

func TestNormalizeCode(t *testing.T) {
	code, err := NormalizeCode(" System Administrator ")
	if err != nil {
		t.Fatalf("NormalizeCode() error = %v", err)
	}
	if code != "system_administrator" {
		t.Fatalf("NormalizeCode() = %q, want system_administrator", code)
	}
}

func TestScopeTypeValidation(t *testing.T) {
	valid := []ScopeType{ScopeGlobal, ScopeOrganizationUnit, ScopeOrganizationUnitAndChildren}
	for _, scope := range valid {
		if !scope.Valid() {
			t.Fatalf("scope %q should be valid", scope)
		}
	}
	if ScopeType("bad").Valid() {
		t.Fatal("bad scope should be invalid")
	}
}

func TestUserStatusValidation(t *testing.T) {
	if !UserStatusActive.Valid() || !UserStatusSuspended.Valid() {
		t.Fatal("expected active and suspended to be valid")
	}
	if UserStatus("bad").Valid() {
		t.Fatal("bad status should be invalid")
	}
}

func TestUserRoleScopeOrganizationUnitValidation(t *testing.T) {
	db := openRBACModelTestDB(t)
	fixtures := createRBACFixtures(t, db)
	unitID := fixtures.unit.ID

	invalid := []struct {
		name string
		role UserRole
	}{
		{
			name: "global with organization unit",
			role: UserRole{UserID: fixtures.user.ID, RoleID: fixtures.role.ID, ScopeType: ScopeGlobal, OrganizationUnitID: &unitID},
		},
		{
			name: "organization unit without organization unit",
			role: UserRole{UserID: fixtures.user.ID, RoleID: fixtures.role.ID, ScopeType: ScopeOrganizationUnit},
		},
		{
			name: "children without organization unit",
			role: UserRole{UserID: fixtures.user.ID, RoleID: fixtures.role.ID, ScopeType: ScopeOrganizationUnitAndChildren},
		},
		{
			name: "invalid scope type",
			role: UserRole{UserID: fixtures.user.ID, RoleID: fixtures.role.ID, ScopeType: ScopeType("bad")},
		},
	}
	for _, tc := range invalid {
		t.Run(tc.name, func(t *testing.T) {
			tc.role.CreatedAt = time.Now().UTC()
			tc.role.UpdatedAt = tc.role.CreatedAt
			if err := db.Create(&tc.role).Error; err == nil {
				t.Fatal("Create() error = nil, want scope validation error")
			}
		})
	}

	valid := UserRole{
		UserID:    fixtures.user.ID,
		RoleID:    fixtures.role.ID,
		ScopeType: ScopeGlobal,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	if err := db.Create(&valid).Error; err != nil {
		t.Fatalf("create valid global user role: %v", err)
	}
	valid.OrganizationUnitID = &unitID
	if err := db.Save(&valid).Error; err == nil {
		t.Fatal("Save() error = nil, want scope validation error")
	}
}

func TestGroupRoleScopeOrganizationUnitValidation(t *testing.T) {
	db := openRBACModelTestDB(t)
	fixtures := createRBACFixtures(t, db)
	unitID := fixtures.unit.ID

	invalid := []struct {
		name string
		role GroupRole
	}{
		{
			name: "global with organization unit",
			role: GroupRole{GroupID: fixtures.group.ID, RoleID: fixtures.role.ID, ScopeType: ScopeGlobal, OrganizationUnitID: &unitID},
		},
		{
			name: "organization unit without organization unit",
			role: GroupRole{GroupID: fixtures.group.ID, RoleID: fixtures.role.ID, ScopeType: ScopeOrganizationUnit},
		},
		{
			name: "children without organization unit",
			role: GroupRole{GroupID: fixtures.group.ID, RoleID: fixtures.role.ID, ScopeType: ScopeOrganizationUnitAndChildren},
		},
		{
			name: "invalid scope type",
			role: GroupRole{GroupID: fixtures.group.ID, RoleID: fixtures.role.ID, ScopeType: ScopeType("bad")},
		},
	}
	for _, tc := range invalid {
		t.Run(tc.name, func(t *testing.T) {
			tc.role.CreatedAt = time.Now().UTC()
			if err := db.Create(&tc.role).Error; err == nil {
				t.Fatal("Create() error = nil, want scope validation error")
			}
		})
	}

	valid := GroupRole{
		GroupID:            fixtures.group.ID,
		RoleID:             fixtures.role.ID,
		ScopeType:          ScopeOrganizationUnit,
		OrganizationUnitID: &unitID,
		CreatedAt:          time.Now().UTC(),
	}
	if err := db.Create(&valid).Error; err != nil {
		t.Fatalf("create valid scoped group role: %v", err)
	}
	valid.OrganizationUnitID = nil
	if err := db.Save(&valid).Error; err == nil {
		t.Fatal("Save() error = nil, want scope validation error")
	}
}

func TestOrganizationUnitRoleScopeValidation(t *testing.T) {
	db := openRBACModelTestDB(t)
	fixtures := createRBACFixtures(t, db)

	assignment := OrganizationUnitRole{
		OrganizationUnitID: fixtures.unit.ID,
		RoleID:             fixtures.role.ID,
		ScopeType:          ScopeType("bad"),
		CreatedAt:          time.Now().UTC(),
	}
	if err := db.Create(&assignment).Error; err == nil {
		t.Fatal("Create() error = nil, want scope validation error")
	}
}

func TestUserRoleScopeDatabaseCheckRejectsContradictoryRows(t *testing.T) {
	db := openRBACModelTestDB(t)
	fixtures := createRBACFixtures(t, db)
	now := time.Now().UTC()

	err := db.Exec(
		`INSERT INTO user_roles (id, user_id, role_id, scope_type, organization_unit_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		uuid.NewString(),
		fixtures.user.ID,
		fixtures.role.ID,
		string(ScopeGlobal),
		fixtures.unit.ID,
		now,
		now,
	).Error
	if err == nil {
		t.Fatal("raw insert error = nil, want database check error")
	}
}

func TestGroupRoleScopeDatabaseCheckRejectsContradictoryRows(t *testing.T) {
	db := openRBACModelTestDB(t)
	fixtures := createRBACFixtures(t, db)
	now := time.Now().UTC()

	err := db.Exec(
		`INSERT INTO group_roles (id, group_id, role_id, scope_type, organization_unit_id, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
		uuid.NewString(),
		fixtures.group.ID,
		fixtures.role.ID,
		string(ScopeOrganizationUnit),
		nil,
		now,
	).Error
	if err == nil {
		t.Fatal("raw insert error = nil, want database check error")
	}
}

func TestOrganizationUnitRoleScopeDatabaseCheckRejectsInvalidScope(t *testing.T) {
	db := openRBACModelTestDB(t)
	fixtures := createRBACFixtures(t, db)

	err := db.Exec(
		`INSERT INTO organization_unit_roles (id, organization_unit_id, role_id, scope_type, created_at) VALUES (?, ?, ?, ?, ?)`,
		uuid.NewString(),
		fixtures.unit.ID,
		fixtures.role.ID,
		"bad",
		time.Now().UTC(),
	).Error
	if err == nil {
		t.Fatal("raw insert error = nil, want database check error")
	}
}

type rbacFixtures struct {
	user  auth.User
	unit  orgunits.Unit
	role  Role
	group Group
}

func openRBACModelTestDB(t *testing.T) *gorm.DB {
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

func createRBACFixtures(t *testing.T, db *gorm.DB) rbacFixtures {
	t.Helper()
	now := time.Now().UTC()
	unitCode := "ENG"
	unit := orgunits.Unit{
		Name:      "Engineering",
		Code:      &unitCode,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := db.Create(&unit).Error; err != nil {
		t.Fatalf("create organization unit: %v", err)
	}
	user := auth.User{
		Name:          "Ada Lovelace",
		Email:         "ada@example.com",
		EmailVerified: true,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	role := Role{
		Name:      "Administrator",
		Code:      "administrator",
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}
	group := Group{
		Name:      "Admins",
		Code:      "admins",
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := db.Create(&group).Error; err != nil {
		t.Fatalf("create group: %v", err)
	}
	return rbacFixtures{user: user, unit: unit, role: role, group: group}
}
