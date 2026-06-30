# RBAC Foundation Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build the RBAC foundation for Syncra DMS: users, statuses, roles, permissions, groups, organization-unit scoped assignments, and backend permission enforcement for existing admin behavior.

**Architecture:** Add a focused Go `internal/rbac` domain with persistence models, seed registry, and a deny-by-default resolver. Keep the existing SvelteKit server-proxy pattern: browser code calls Svelte `/api/...` routes, which call server-only frontend clients, which call the Go API with the internal token and session cookie. Migrate existing admin checks to RBAC permissions while keeping the old `user.role` field only as a compatibility bridge.

**Tech Stack:** Go, Gin, GORM, PostgreSQL, Atlas, SQLite-backed handler tests, SvelteKit, Svelte 5, TypeScript, Vitest, TanStack Svelte Query, Tailwind CSS, local shadcn-style UI primitives.

---

## Scope Check

This plan implements the approved option A: RBAC admin foundation only.

It does not implement documents, workflows, document owner units, document visibility, workflow routing, audit logs, external identity sync, direct user permissions, deny rules, temporary access expiry, or field-level permissions.

## Relevant Docs And Constraints

- Approved design: `docs/plans/2026-06-30-rbac-foundation-design.md`
- Source spec: `specs.md`
- Repository instructions: `AGENTS.md`, `go-server/AGENT.md`, `frontend/AGENT.md`
- Prefix every shell command with `rtk`.
- Use `apply_patch` for manual file edits.
- Add behavior tests before production behavior.
- Keep private tokens in server-only modules.
- Do not copy product modules from `example/`.
- Update Swagger when Go API routes or shapes change.
- Before Svelte component implementation, use the available Svelte tooling/autofixer if exposed in the execution session.

## Permission Codes For This Slice

Seed these categories and permissions first. The registry can include document/workflow permissions from the spec, but the implementation should only enforce admin and organization-unit permissions in this slice.

User Management:
- `user.view`
- `user.create`
- `user.update`
- `user.delete`
- `user.activate`
- `user.suspend`
- `user.assign_role`
- `user.assign_group`
- `user.assign_unit`

Role Management:
- `role.view`
- `role.create`
- `role.update`
- `role.delete`
- `role.assign_permissions`
- `role.assign_users`

Group Management:
- `group.view`
- `group.create`
- `group.update`
- `group.delete`
- `group.manage_users`
- `group.assign_roles`

Organization Unit Management:
- `organization_unit.view`
- `organization_unit.create`
- `organization_unit.update`
- `organization_unit.delete`
- `organization_unit.manage_users`
- `organization_unit.manage_roles`
- `organization_unit.manage_permissions`
- `organization_unit.manage_hierarchy`
- `organization_unit.view_audit`

System:
- `system.admin`

Default roles:
- `system_administrator`: all permissions above, global scope.
- `organization_administrator`: user, role, group, and organization-unit management except `system.admin`.
- `unit_manager`: organization-unit view/update/manage_users within `organization_unit_and_children`.
- `viewer`: organization-unit view only.

## Task 1: Backend RBAC Domain Models And Migration

**Files:**
- Create: `go-server/internal/rbac/models.go`
- Create: `go-server/internal/rbac/models_test.go`
- Modify: `go-server/internal/auth/user.go`
- Modify: `go-server/internal/database/database.go`
- Create: `go-server/migrations/20260630150000_add_rbac_foundation.sql`
- Modify: `go-server/migrations/atlas.sum`

**Step 1: Write failing model tests**

Create `go-server/internal/rbac/models_test.go`:

```go
package rbac

import "testing"

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
```

**Step 2: Run tests to verify failure**

Run:

```sh
cd go-server
rtk go test ./internal/rbac
```

Expected: FAIL because `internal/rbac` does not exist.

**Step 3: Implement RBAC models**

Create `go-server/internal/rbac/models.go`:

```go
package rbac

import (
	"errors"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"ai.ro/syncra/dms/internal/auth"
	"ai.ro/syncra/dms/internal/orgunits"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScopeType string

const (
	ScopeGlobal                      ScopeType = "global"
	ScopeOrganizationUnit            ScopeType = "organization_unit"
	ScopeOrganizationUnitAndChildren ScopeType = "organization_unit_and_children"
)

func (s ScopeType) Valid() bool {
	return s == ScopeGlobal || s == ScopeOrganizationUnit || s == ScopeOrganizationUnitAndChildren
}

type UserStatus string

const (
	UserStatusInvited   UserStatus = "invited"
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusDeleted   UserStatus = "deleted"
)

func (s UserStatus) Valid() bool {
	return s == UserStatusInvited || s == UserStatusActive || s == UserStatusInactive || s == UserStatusSuspended || s == UserStatusDeleted
}

type Role struct {
	ID          string    `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string    `gorm:"not null;size:160" json:"name"`
	Code        string    `gorm:"not null;size:80;uniqueIndex" json:"code"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`
	IsSystem    bool      `gorm:"column:is_system;not null;default:false" json:"isSystem"`
	IsActive    bool      `gorm:"column:is_active;not null;default:true;index" json:"isActive"`
	CreatedAt   time.Time `gorm:"column:created_at;not null" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null" json:"updatedAt"`
}

func (r *Role) BeforeCreate(_ *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.NewString()
	}
	return nil
}

func (Role) TableName() string { return "roles" }

type Permission struct {
	ID          string    `gorm:"type:uuid;primaryKey" json:"id"`
	Code        string    `gorm:"not null;size:120;uniqueIndex" json:"code"`
	Name        string    `gorm:"not null;size:160" json:"name"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`
	Category    string    `gorm:"not null;size:80;index" json:"category"`
	IsSystem    bool      `gorm:"column:is_system;not null;default:true" json:"isSystem"`
	CreatedAt   time.Time `gorm:"column:created_at;not null" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null" json:"updatedAt"`
}

func (p *Permission) BeforeCreate(_ *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	return nil
}

func (Permission) TableName() string { return "permissions" }

type RolePermission struct {
	ID           string     `gorm:"type:uuid;primaryKey" json:"id"`
	RoleID       string     `gorm:"column:role_id;type:uuid;not null;uniqueIndex:idx_role_permission_unique,priority:1" json:"roleId"`
	Role         Role       `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	PermissionID string     `gorm:"column:permission_id;type:uuid;not null;uniqueIndex:idx_role_permission_unique,priority:2" json:"permissionId"`
	Permission   Permission `gorm:"foreignKey:PermissionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	CreatedAt    time.Time  `gorm:"column:created_at;not null" json:"createdAt"`
}

func (r *RolePermission) BeforeCreate(_ *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.NewString()
	}
	return nil
}

func (RolePermission) TableName() string { return "role_permissions" }

type UserRole struct {
	ID                 string         `gorm:"type:uuid;primaryKey" json:"id"`
	UserID             string         `gorm:"column:user_id;type:uuid;not null;index;uniqueIndex:idx_user_role_scope_unique,priority:1" json:"userId"`
	User               auth.User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	RoleID             string         `gorm:"column:role_id;type:uuid;not null;index;uniqueIndex:idx_user_role_scope_unique,priority:2" json:"roleId"`
	Role               Role           `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	ScopeType          ScopeType      `gorm:"column:scope_type;not null;size:40;uniqueIndex:idx_user_role_scope_unique,priority:3" json:"scopeType"`
	OrganizationUnitID *string        `gorm:"column:organization_unit_id;type:uuid;uniqueIndex:idx_user_role_scope_unique,priority:4" json:"organizationUnitId,omitempty"`
	OrganizationUnit   *orgunits.Unit `gorm:"foreignKey:OrganizationUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
	CreatedAt          time.Time      `gorm:"column:created_at;not null" json:"createdAt"`
	UpdatedAt          time.Time      `gorm:"column:updated_at;not null" json:"updatedAt"`
}

func (r *UserRole) BeforeCreate(_ *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.NewString()
	}
	return nil
}

func (UserRole) TableName() string { return "user_roles" }

type Group struct {
	ID                 string         `gorm:"type:uuid;primaryKey" json:"id"`
	Name               string         `gorm:"not null;size:160" json:"name"`
	Code               string         `gorm:"not null;size:80;uniqueIndex" json:"code"`
	Description        *string        `gorm:"type:text" json:"description,omitempty"`
	OrganizationUnitID *string        `gorm:"column:organization_unit_id;type:uuid;index" json:"organizationUnitId,omitempty"`
	OrganizationUnit   *orgunits.Unit `gorm:"foreignKey:OrganizationUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
	IsActive           bool           `gorm:"column:is_active;not null;default:true;index" json:"isActive"`
	CreatedAt          time.Time      `gorm:"column:created_at;not null" json:"createdAt"`
	UpdatedAt          time.Time      `gorm:"column:updated_at;not null" json:"updatedAt"`
}

func (g *Group) BeforeCreate(_ *gorm.DB) error {
	if g.ID == "" {
		g.ID = uuid.NewString()
	}
	return nil
}

func (Group) TableName() string { return "groups" }

type GroupUser struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	GroupID   string    `gorm:"column:group_id;type:uuid;not null;uniqueIndex:idx_group_user_unique,priority:1" json:"groupId"`
	Group     Group     `gorm:"foreignKey:GroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	UserID    string    `gorm:"column:user_id;type:uuid;not null;uniqueIndex:idx_group_user_unique,priority:2" json:"userId"`
	User      auth.User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"createdAt"`
}

func (g *GroupUser) BeforeCreate(_ *gorm.DB) error {
	if g.ID == "" {
		g.ID = uuid.NewString()
	}
	return nil
}

func (GroupUser) TableName() string { return "group_users" }

type GroupRole struct {
	ID                 string         `gorm:"type:uuid;primaryKey" json:"id"`
	GroupID            string         `gorm:"column:group_id;type:uuid;not null;uniqueIndex:idx_group_role_scope_unique,priority:1" json:"groupId"`
	Group              Group          `gorm:"foreignKey:GroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	RoleID             string         `gorm:"column:role_id;type:uuid;not null;uniqueIndex:idx_group_role_scope_unique,priority:2" json:"roleId"`
	Role               Role           `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	ScopeType          ScopeType      `gorm:"column:scope_type;not null;size:40;uniqueIndex:idx_group_role_scope_unique,priority:3" json:"scopeType"`
	OrganizationUnitID *string        `gorm:"column:organization_unit_id;type:uuid;uniqueIndex:idx_group_role_scope_unique,priority:4" json:"organizationUnitId,omitempty"`
	OrganizationUnit   *orgunits.Unit `gorm:"foreignKey:OrganizationUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
	CreatedAt          time.Time      `gorm:"column:created_at;not null" json:"createdAt"`
}

func (g *GroupRole) BeforeCreate(_ *gorm.DB) error {
	if g.ID == "" {
		g.ID = uuid.NewString()
	}
	return nil
}

func (GroupRole) TableName() string { return "group_roles" }

type OrganizationUnitRole struct {
	ID                 string        `gorm:"type:uuid;primaryKey" json:"id"`
	OrganizationUnitID string        `gorm:"column:organization_unit_id;type:uuid;not null;uniqueIndex:idx_organization_unit_role_unique,priority:1" json:"organizationUnitId"`
	OrganizationUnit   orgunits.Unit `gorm:"foreignKey:OrganizationUnitID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	RoleID             string        `gorm:"column:role_id;type:uuid;not null;uniqueIndex:idx_organization_unit_role_unique,priority:2" json:"roleId"`
	Role               Role          `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	ScopeType          ScopeType     `gorm:"column:scope_type;not null;size:40;uniqueIndex:idx_organization_unit_role_unique,priority:3" json:"scopeType"`
	CreatedAt          time.Time     `gorm:"column:created_at;not null" json:"createdAt"`
}

func (o *OrganizationUnitRole) BeforeCreate(_ *gorm.DB) error {
	if o.ID == "" {
		o.ID = uuid.NewString()
	}
	return nil
}

func (OrganizationUnitRole) TableName() string { return "organization_unit_roles" }

var codeNonAlnum = regexp.MustCompile(`[^a-z0-9]+`)

func NormalizeCode(raw string) (string, error) {
	code := strings.TrimSpace(strings.ToLower(raw))
	code = codeNonAlnum.ReplaceAllString(code, "_")
	code = strings.Trim(code, "_")
	if code == "" {
		return "", errors.New("code is required")
	}
	if utf8.RuneCountInString(code) > 80 {
		return "", errors.New("code must be at most 80 characters")
	}
	return code, nil
}
```

**Step 4: Extend the auth user model**

Modify `go-server/internal/auth/user.go`:

```go
// add fields to User
Status                    string     `gorm:"not null;size:40;default:active;index;check:chk_user_status,status IN ('invited','active','inactive','suspended','deleted')" json:"status"`
PrimaryOrganizationUnitID *string    `gorm:"column:primary_organization_unit_id;type:uuid;index" json:"primaryOrganizationUnitId,omitempty"`
ManagerUserID             *string    `gorm:"column:manager_user_id;type:uuid;index" json:"managerUserId,omitempty"`
JobTitle                  *string    `gorm:"column:job_title;size:160" json:"jobTitle,omitempty"`
Phone                     *string    `gorm:"size:80" json:"phone,omitempty"`
DeletedAt                 *time.Time `gorm:"column:deleted_at;index" json:"deletedAt,omitempty"`
```

Add defaults in `BeforeCreate`:

```go
if u.Status == "" {
	if u.EmailVerified {
		u.Status = "active"
	} else {
		u.Status = "invited"
	}
}
```

Keep `Role UserRole` for compatibility in this task.

**Step 5: Register models**

Modify `go-server/internal/database/database.go` imports:

```go
import "ai.ro/syncra/dms/internal/rbac"
```

Add models to `ApplicationModels()` after organization units:

```go
&rbac.Role{},
&rbac.Permission{},
&rbac.RolePermission{},
&rbac.UserRole{},
&rbac.Group{},
&rbac.GroupUser{},
&rbac.GroupRole{},
&rbac.OrganizationUnitRole{},
```

**Step 6: Add migration**

Create `go-server/migrations/20260630150000_add_rbac_foundation.sql`. Include:

- `ALTER TABLE "user"` additions for status and profile fields.
- `UPDATE "user" SET status = CASE WHEN email_verified THEN 'active' ELSE 'invited' END WHERE status IS NULL`.
- Check constraint for user status.
- Tables listed in the design.
- Unique indexes on codes and assignment uniqueness.
- Foreign keys to `"user"`, `roles`, `permissions`, `groups`, and `organization_units`.

Use explicit SQL rather than relying on AutoMigrate.

**Step 7: Refresh migration hash**

Run:

```sh
cd go-server
rtk atlas migrate hash --dir file://migrations
```

Expected: `go-server/migrations/atlas.sum` changes.

**Step 8: Run tests and validation**

Run:

```sh
cd go-server
rtk go test ./internal/rbac ./internal/database
rtk atlas migrate validate --dir file://migrations
```

Expected: PASS.

**Step 9: Commit**

```sh
rtk git add go-server/internal/rbac go-server/internal/auth/user.go go-server/internal/database/database.go go-server/migrations/20260630150000_add_rbac_foundation.sql go-server/migrations/atlas.sum
rtk git commit -m "Add RBAC foundation models"
```

## Task 2: Backend Permission Registry And Seeder

**Files:**
- Create: `go-server/internal/rbac/registry.go`
- Create: `go-server/internal/rbac/seed.go`
- Create: `go-server/internal/rbac/seed_test.go`

**Step 1: Write failing seed tests**

Create `go-server/internal/rbac/seed_test.go`:

```go
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
```

**Step 2: Run test to verify failure**

```sh
cd go-server
rtk go test ./internal/rbac -run SeedDefaults
```

Expected: FAIL because `SeedDefaults` and registry constants do not exist.

**Step 3: Implement registry**

Create `go-server/internal/rbac/registry.go` with typed definitions:

```go
package rbac

const (
	SystemAdministratorRoleCode    = "system_administrator"
	OrganizationAdministratorCode  = "organization_administrator"
	UnitManagerRoleCode            = "unit_manager"
	ViewerRoleCode                 = "viewer"
)

type PermissionDefinition struct {
	Code        string
	Name        string
	Description string
	Category    string
}

type RoleDefinition struct {
	Code            string
	Name            string
	Description     string
	PermissionCodes []string
}

var PermissionRegistry = []PermissionDefinition{
	{Code: "system.admin", Name: "System administration", Category: "System"},
	{Code: "user.view", Name: "View users", Category: "User Management"},
	{Code: "user.create", Name: "Create users", Category: "User Management"},
	{Code: "user.update", Name: "Update users", Category: "User Management"},
	{Code: "user.delete", Name: "Delete users", Category: "User Management"},
	{Code: "user.activate", Name: "Activate users", Category: "User Management"},
	{Code: "user.suspend", Name: "Suspend users", Category: "User Management"},
	{Code: "user.assign_role", Name: "Assign user roles", Category: "User Management"},
	{Code: "user.assign_group", Name: "Assign user groups", Category: "User Management"},
	{Code: "user.assign_unit", Name: "Assign user units", Category: "User Management"},
	{Code: "role.view", Name: "View roles", Category: "Role Management"},
	{Code: "role.create", Name: "Create roles", Category: "Role Management"},
	{Code: "role.update", Name: "Update roles", Category: "Role Management"},
	{Code: "role.delete", Name: "Delete roles", Category: "Role Management"},
	{Code: "role.assign_permissions", Name: "Assign role permissions", Category: "Role Management"},
	{Code: "role.assign_users", Name: "Assign roles to users", Category: "Role Management"},
	{Code: "group.view", Name: "View groups", Category: "Group Management"},
	{Code: "group.create", Name: "Create groups", Category: "Group Management"},
	{Code: "group.update", Name: "Update groups", Category: "Group Management"},
	{Code: "group.delete", Name: "Delete groups", Category: "Group Management"},
	{Code: "group.manage_users", Name: "Manage group users", Category: "Group Management"},
	{Code: "group.assign_roles", Name: "Assign group roles", Category: "Group Management"},
	{Code: "organization_unit.view", Name: "View organization units", Category: "Organization Unit Management"},
	{Code: "organization_unit.create", Name: "Create organization units", Category: "Organization Unit Management"},
	{Code: "organization_unit.update", Name: "Update organization units", Category: "Organization Unit Management"},
	{Code: "organization_unit.delete", Name: "Delete organization units", Category: "Organization Unit Management"},
	{Code: "organization_unit.manage_users", Name: "Manage organization unit users", Category: "Organization Unit Management"},
	{Code: "organization_unit.manage_roles", Name: "Manage organization unit roles", Category: "Organization Unit Management"},
	{Code: "organization_unit.manage_permissions", Name: "Manage organization unit permissions", Category: "Organization Unit Management"},
	{Code: "organization_unit.manage_hierarchy", Name: "Manage organization unit hierarchy", Category: "Organization Unit Management"},
	{Code: "organization_unit.view_audit", Name: "View organization unit audit", Category: "Organization Unit Management"},
}
```

Add helper functions:

```go
func PermissionCodes() []string
func PermissionByCode(code string) (PermissionDefinition, bool)
func DefaultRoles() []RoleDefinition
```

`DefaultRoles()` should grant all permissions to `system_administrator`, all non-system admin permissions to `organization_administrator`, org-unit management subset to `unit_manager`, and `organization_unit.view` to `viewer`.

**Step 4: Implement seeding**

Create `go-server/internal/rbac/seed.go`:

```go
package rbac

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedDefaults(db *gorm.DB) error {
	now := time.Now().UTC()
	return db.Transaction(func(tx *gorm.DB) error {
		permissionIDs := map[string]string{}
		for _, definition := range PermissionRegistry {
			permission := Permission{
				Code:      definition.Code,
				Name:      definition.Name,
				Category:  definition.Category,
				IsSystem:  true,
				CreatedAt: now,
				UpdatedAt: now,
			}
			if definition.Description != "" {
				permission.Description = &definition.Description
			}
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "code"}},
				DoUpdates: clause.AssignmentColumns([]string{"name", "description", "category", "is_system", "updated_at"}),
			}).Create(&permission).Error; err != nil {
				return err
			}
			var saved Permission
			if err := tx.First(&saved, "code = ?", definition.Code).Error; err != nil {
				return err
			}
			permissionIDs[definition.Code] = saved.ID
		}

		for _, definition := range DefaultRoles() {
			role := Role{
				Code:        definition.Code,
				Name:        definition.Name,
				IsSystem:    true,
				IsActive:    true,
				CreatedAt:   now,
				UpdatedAt:   now,
			}
			if definition.Description != "" {
				role.Description = &definition.Description
			}
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "code"}},
				DoUpdates: clause.AssignmentColumns([]string{"name", "description", "is_system", "is_active", "updated_at"}),
			}).Create(&role).Error; err != nil {
				return err
			}
			if err := tx.First(&role, "code = ?", definition.Code).Error; err != nil {
				return err
			}
			for _, permissionCode := range definition.PermissionCodes {
				permissionID := permissionIDs[permissionCode]
				if permissionID == "" {
					continue
				}
				link := RolePermission{RoleID: role.ID, PermissionID: permissionID, CreatedAt: now}
				if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&link).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}
```

**Step 5: Wire seeding into app startup**

Modify `go-server/internal/app/api.go` after opening the database:

```go
if err := rbac.SeedDefaults(db); err != nil {
	return nil, err
}
```

If `api.go` does not currently return errors from initialization cleanly, add the smallest local adjustment and test it in the next task.

**Step 6: Run tests**

```sh
cd go-server
rtk go test ./internal/rbac ./internal/app
```

Expected: PASS.

**Step 7: Commit**

```sh
rtk git add go-server/internal/rbac/registry.go go-server/internal/rbac/seed.go go-server/internal/rbac/seed_test.go go-server/internal/app/api.go
rtk git commit -m "Seed RBAC permissions and roles"
```

## Task 3: Permission Resolver

**Files:**
- Create: `go-server/internal/rbac/resolver.go`
- Create: `go-server/internal/rbac/resolver_test.go`

**Step 1: Write failing resolver tests**

Create `go-server/internal/rbac/resolver_test.go` with tests for:

- inactive user denied;
- direct user role grants global permission;
- group role grants permission through membership;
- organization-unit role grants permission through primary unit;
- `organization_unit_and_children` matches descendant unit;
- unmatched scope denies.

Use SQLite AutoMigrate with `auth.User`, `orgunits.Unit`, and all RBAC models.

Key test shape:

```go
func TestResolverAllowsUserRoleWithGlobalScope(t *testing.T) {
	db := newResolverTestDB(t)
	user := createResolverUser(t, db, "active", nil)
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
```

**Step 2: Run test to verify failure**

```sh
cd go-server
rtk go test ./internal/rbac -run Resolver
```

Expected: FAIL because resolver types do not exist.

**Step 3: Implement resolver types**

Create `go-server/internal/rbac/resolver.go`:

```go
package rbac

import (
	"context"
	"errors"

	"ai.ro/syncra/dms/internal/auth"
	"ai.ro/syncra/dms/internal/orgunits"
	"gorm.io/gorm"
)

type Resolver struct {
	db *gorm.DB
}

type Check struct {
	UserID             string
	Permission         string
	OrganizationUnitID *string
}

type Grant struct {
	PermissionCode     string
	ScopeType          ScopeType
	OrganizationUnitID *string
	Source             string
}

func NewResolver(db *gorm.DB) *Resolver {
	return &Resolver{db: db}
}

func (r *Resolver) Can(ctx context.Context, check Check) (bool, error) {
	grants, err := r.EffectiveGrants(ctx, check.UserID)
	if err != nil {
		return false, err
	}
	for _, grant := range grants {
		if grant.PermissionCode != check.Permission {
			continue
		}
		matches, err := r.scopeMatches(ctx, grant, check.OrganizationUnitID)
		if err != nil {
			return false, err
		}
		if matches {
			return true, nil
		}
	}
	return false, nil
}
```

Implement `EffectiveGrants` by querying:

- `user_roles` joined to active `roles`, `role_permissions`, and `permissions`;
- `group_users` -> active `groups` -> `group_roles` -> active `roles` -> `role_permissions` -> `permissions`;
- active user's `primary_organization_unit_id` -> `organization_unit_roles` -> active `roles` -> `role_permissions` -> `permissions`.

Implementation can use small row structs and SQL joins through GORM `Table`.

`scopeMatches` rules:

```go
case ScopeGlobal:
	return true, nil
case ScopeOrganizationUnit:
	return grant.OrganizationUnitID != nil && requested != nil && *grant.OrganizationUnitID == *requested, nil
case ScopeOrganizationUnitAndChildren:
	if grant.OrganizationUnitID == nil || requested == nil {
		return false, nil
	}
	if *grant.OrganizationUnitID == *requested {
		return true, nil
	}
	var units []orgunits.Unit
	if err := r.db.WithContext(ctx).Where("archived_at IS NULL").Find(&units).Error; err != nil {
		return false, err
	}
	return orgunits.DescendantIDs(*grant.OrganizationUnitID, units)[*requested], nil
default:
	return false, nil
}
```

Treat missing users and non-active users as denied, not errors. Return errors only for database failures.

**Step 4: Run resolver tests**

```sh
cd go-server
rtk go test ./internal/rbac -run Resolver
```

Expected: PASS.

**Step 5: Commit**

```sh
rtk git add go-server/internal/rbac/resolver.go go-server/internal/rbac/resolver_test.go
rtk git commit -m "Add RBAC permission resolver"
```

## Task 4: Auth Status Enforcement And Admin Bootstrap

**Files:**
- Modify: `go-server/internal/api/auth_handlers.go`
- Modify: `go-server/internal/api/auth_handlers_test.go`
- Create: `go-server/internal/rbac/bootstrap.go`
- Create: `go-server/internal/rbac/bootstrap_test.go`

**Step 1: Write failing auth tests**

Add tests to `go-server/internal/api/auth_handlers_test.go`:

```go
func TestInactiveAndSuspendedUsersCannotSignIn(t *testing.T) {
	router, db := newAuthTestRouter(t)
	inactive := createVerifiedUser(t, db, "inactive@example.com", "password123")
	if err := db.Model(&auth.User{}).Where("id = ?", inactive.ID).Update("status", "inactive").Error; err != nil {
		t.Fatalf("set inactive: %v", err)
	}

	response := authJSON(t, router, http.MethodPost, "/api/auth/sign-in/email", `{
		"email":"inactive@example.com",
		"password":"password123"
	}`, map[string]string{"X-Syncra-Internal-Token": testInternalToken})
	if response.Code != http.StatusForbidden {
		t.Fatalf("status = %d body=%s, want forbidden", response.Code, response.Body.String())
	}
}

func TestSuspendedUserSessionIsRejected(t *testing.T) {
	router, db := newAuthTestRouter(t)
	user := createVerifiedUser(t, db, "ada@example.com", "password123")
	token := loginUser(t, router, user.Email, "password123")
	if err := db.Model(&auth.User{}).Where("id = ?", user.ID).Update("status", "suspended").Error; err != nil {
		t.Fatalf("suspend user: %v", err)
	}

	response := authJSON(t, router, http.MethodGet, "/api/auth/get-session", "", map[string]string{
		"X-Syncra-Internal-Token": testInternalToken,
		"Cookie":                  "auth.session_token=" + token,
	})
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", response.Code, response.Body.String())
	}
	if strings.TrimSpace(response.Body.String()) != "null" {
		t.Fatalf("body = %s, want null", response.Body.String())
	}
}
```

**Step 2: Run tests to verify failure**

```sh
cd go-server
rtk go test ./internal/api -run 'Inactive|Suspended'
```

Expected: FAIL because status is not enforced.

**Step 3: Enforce status in auth handlers**

Modify `signInEmail` after email verification:

```go
if !authUserActive(user) {
	writeError(c, http.StatusForbidden, "user account is not active")
	return
}
```

Modify `loadAuthenticatedSession` after loading the user:

```go
if !authUserActive(session.User) {
	if err := h.db.Where("token = ?", token).Delete(&auth.Session{}).Error; err != nil {
		return auth.Session{}, false, errors.New("failed to delete inactive user session")
	}
	h.clearSessionCookie(c)
	return auth.Session{}, false, nil
}
```

Add helper:

```go
func authUserActive(user auth.User) bool {
	return user.Status == "" || user.Status == "active"
}
```

Accept empty status temporarily so older test fixtures stay valid until all helpers set status.

**Step 4: Add RBAC bootstrap tests**

Create `go-server/internal/rbac/bootstrap_test.go`:

```go
func TestBootstrapLegacyAdminsAssignsSystemAdministratorRole(t *testing.T) {
	db := newBootstrapTestDB(t)
	if err := SeedDefaults(db); err != nil {
		t.Fatalf("seed: %v", err)
	}
	admin := auth.User{Name: "Admin", Email: "admin@example.com", EmailVerified: true, Role: auth.UserRoleAdmin, Status: string(UserStatusActive), CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := db.Create(&admin).Error; err != nil {
		t.Fatalf("create admin: %v", err)
	}

	if err := BootstrapLegacyAdmins(db); err != nil {
		t.Fatalf("bootstrap: %v", err)
	}
	if err := BootstrapLegacyAdmins(db); err != nil {
		t.Fatalf("bootstrap second: %v", err)
	}

	var count int64
	if err := db.Model(&UserRole{}).Where("user_id = ? AND scope_type = ?", admin.ID, ScopeGlobal).Count(&count).Error; err != nil {
		t.Fatalf("count user roles: %v", err)
	}
	if count != 1 {
		t.Fatalf("user role count = %d, want 1", count)
	}
}
```

**Step 5: Implement bootstrap**

Create `go-server/internal/rbac/bootstrap.go`:

```go
package rbac

import (
	"time"

	"ai.ro/syncra/dms/internal/auth"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func BootstrapLegacyAdmins(db *gorm.DB) error {
	var adminRole Role
	if err := db.First(&adminRole, "code = ?", SystemAdministratorRoleCode).Error; err != nil {
		return err
	}
	var users []auth.User
	if err := db.Where("role = ?", auth.UserRoleAdmin).Find(&users).Error; err != nil {
		return err
	}
	now := time.Now().UTC()
	for _, user := range users {
		link := UserRole{
			UserID:    user.ID,
			RoleID:    adminRole.ID,
			ScopeType: ScopeGlobal,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&link).Error; err != nil {
			return err
		}
	}
	return nil
}
```

Wire after `SeedDefaults(db)` in app startup.

**Step 6: Run tests**

```sh
cd go-server
rtk go test ./internal/api ./internal/rbac
```

Expected: PASS.

**Step 7: Commit**

```sh
rtk git add go-server/internal/api/auth_handlers.go go-server/internal/api/auth_handlers_test.go go-server/internal/rbac/bootstrap.go go-server/internal/rbac/bootstrap_test.go go-server/internal/app/api.go
rtk git commit -m "Enforce user status for authentication"
```

## Task 5: Shared Backend Authorization Helpers

**Files:**
- Create: `go-server/internal/api/authorization.go`
- Create: `go-server/internal/api/authorization_test.go`
- Modify: `go-server/internal/api/auth_handlers_test.go`

**Step 1: Write failing helper tests**

Create `go-server/internal/api/authorization_test.go`:

```go
package api

import (
	"net/http"
	"testing"

	"ai.ro/syncra/dms/internal/rbac"
)

func TestRequirePermissionAllowsSeededLegacyAdmin(t *testing.T) {
	router, db := newAuthTestRouter(t)
	admin := createAdminUser(t, db, "admin@example.com", "password123")
	if err := rbac.SeedDefaults(db); err != nil {
		t.Fatalf("seed: %v", err)
	}
	if err := rbac.BootstrapLegacyAdmins(db); err != nil {
		t.Fatalf("bootstrap: %v", err)
	}
	token := loginUser(t, router, admin.Email, "password123")

	response := orgUnitJSON(t, router, http.MethodPost, "/api/organization-units", `{"name":"Company"}`, authCookieHeaders(token))
	if response.Code != http.StatusCreated {
		t.Fatalf("status = %d body=%s", response.Code, response.Body.String())
	}
}
```

This test initially depends on organization-unit routes still using old checks; it will become meaningful when Task 8 replaces those checks. If needed, add a small test-only route inside `authorization_test.go` that calls `requirePermission`.

**Step 2: Implement shared helper**

Create `go-server/internal/api/authorization.go`:

```go
package api

import (
	"net/http"

	"ai.ro/syncra/dms/internal/auth"
	"ai.ro/syncra/dms/internal/rbac"
	"github.com/gin-gonic/gin"
)

func requireAuthenticatedUser(c *gin.Context, h *authHandler) (auth.User, bool) {
	if h == nil || !h.authConfigured(c) {
		return auth.User{}, false
	}
	session, ok, err := h.loadAuthenticatedSession(c)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return auth.User{}, false
	}
	if !ok {
		writeError(c, http.StatusUnauthorized, "authenticated session required")
		return auth.User{}, false
	}
	return session.User, true
}

func requirePermission(c *gin.Context, h *authHandler, permission string, organizationUnitID *string) (auth.User, bool) {
	user, ok := requireAuthenticatedUser(c, h)
	if !ok {
		return auth.User{}, false
	}
	allowed, err := rbac.NewResolver(h.db).Can(c.Request.Context(), rbac.Check{
		UserID:             user.ID,
		Permission:         permission,
		OrganizationUnitID: organizationUnitID,
	})
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to check permission")
		return auth.User{}, false
	}
	if !allowed {
		writeError(c, http.StatusForbidden, "permission required")
		return auth.User{}, false
	}
	return user, true
}
```

**Step 3: Update test router setup**

Modify `newAuthTestRouterWithOptions` in `go-server/internal/api/auth_handlers_test.go` to AutoMigrate RBAC models and call `rbac.SeedDefaults(db)`.

Expected imports:

```go
"ai.ro/syncra/dms/internal/rbac"
```

**Step 4: Run tests**

```sh
cd go-server
rtk go test ./internal/api
```

Expected: PASS.

**Step 5: Commit**

```sh
rtk git add go-server/internal/api/authorization.go go-server/internal/api/authorization_test.go go-server/internal/api/auth_handlers_test.go
rtk git commit -m "Add API permission helpers"
```

## Task 6: User Management Backend API

**Files:**
- Create: `go-server/internal/api/users.go`
- Create: `go-server/internal/api/users_test.go`
- Modify: `go-server/internal/api/router.go`

**Step 1: Write failing user API tests**

Create `go-server/internal/api/users_test.go` covering:

- unauthenticated list returns `401`;
- user without `user.view` returns `403`;
- seeded admin can list users;
- admin can create a user with status `invited`;
- admin can update profile fields;
- admin can activate, deactivate, suspend;
- suspend deletes active sessions;
- admin can assign primary organization unit;
- admin can assign and remove scoped roles;
- admin can assign and remove groups.

Start with this first test:

```go
func TestUserAPIRequiresUserViewPermission(t *testing.T) {
	router, db := newAuthTestRouter(t)
	user := createVerifiedUser(t, db, "user@example.com", "password123")
	token := loginUser(t, router, user.Email, "password123")

	response := authJSON(t, router, http.MethodGet, "/api/users", "", authCookieHeaders(token))
	if response.Code != http.StatusForbidden {
		t.Fatalf("status = %d body=%s, want forbidden", response.Code, response.Body.String())
	}
}
```

**Step 2: Run tests to verify failure**

```sh
cd go-server
rtk go test ./internal/api -run UserAPI
```

Expected: FAIL because `/api/users` is not routed.

**Step 3: Implement users handler**

Create `go-server/internal/api/users.go` with:

- `userHandler` struct containing `db` and `auth`.
- typed request structs for create, update, status, primary unit, role assignment, group assignment.
- response struct omitting password hash and session data.
- helper `userResponseFromModel`.
- validation for email, name, status, scope type, and UUIDs.

Minimum routes:

```go
GET    /api/users
GET    /api/users/:id
POST   /api/users
PATCH  /api/users/:id
POST   /api/users/:id/activate
POST   /api/users/:id/deactivate
POST   /api/users/:id/suspend
DELETE /api/users/:id
POST   /api/users/:id/primary-organization-unit
POST   /api/users/:id/roles
DELETE /api/users/:id/roles/:assignmentId
POST   /api/users/:id/groups
DELETE /api/users/:id/groups/:groupId
```

Permission checks:

- read routes: `user.view`
- create: `user.create`
- update profile/status: `user.update`, plus `user.activate` or `user.suspend` for status routes
- soft delete: `user.delete`
- primary unit: `user.assign_unit`
- roles: `user.assign_role`
- groups: `user.assign_group`

Soft delete should set status `deleted`, set `deleted_at`, and delete sessions.

**Step 4: Wire routes**

Modify `go-server/internal/api/router.go`:

```go
users := newUserHandler(options, auth)
userAPI := router.Group("/api/users")
userAPI.Use(auth.requireTrustedInternalRequest())
userAPI.GET("", users.list)
userAPI.GET("/:id", users.get)
userAPI.POST("", users.create)
userAPI.PATCH("/:id", users.update)
userAPI.POST("/:id/activate", users.activate)
userAPI.POST("/:id/deactivate", users.deactivate)
userAPI.POST("/:id/suspend", users.suspend)
userAPI.DELETE("/:id", users.softDelete)
userAPI.POST("/:id/primary-organization-unit", users.setPrimaryOrganizationUnit)
userAPI.POST("/:id/roles", users.assignRole)
userAPI.DELETE("/:id/roles/:assignmentId", users.removeRole)
userAPI.POST("/:id/groups", users.addGroup)
userAPI.DELETE("/:id/groups/:groupId", users.removeGroup)
```

**Step 5: Run user API tests**

```sh
cd go-server
rtk go test ./internal/api -run UserAPI
```

Expected: PASS.

**Step 6: Commit**

```sh
rtk git add go-server/internal/api/users.go go-server/internal/api/users_test.go go-server/internal/api/router.go
rtk git commit -m "Add user management API"
```

## Task 7: Role And Permission Backend API

**Files:**
- Create: `go-server/internal/api/roles.go`
- Create: `go-server/internal/api/roles_test.go`
- Create: `go-server/internal/api/permissions.go`
- Create: `go-server/internal/api/permissions_test.go`
- Modify: `go-server/internal/api/router.go`

**Step 1: Write failing tests**

Create tests for:

- permissions list requires `role.view` or `system.admin`;
- permissions categories returns unique categories;
- role list requires `role.view`;
- admin can create custom role;
- admin cannot mutate system role code or delete system role;
- admin can assign and remove permissions;
- duplicate role code returns `409`.

**Step 2: Run tests to verify failure**

```sh
cd go-server
rtk go test ./internal/api -run 'RoleAPI|PermissionAPI'
```

Expected: FAIL because routes do not exist.

**Step 3: Implement permissions handler**

Create `go-server/internal/api/permissions.go`:

- `GET /api/permissions`
- `GET /api/permissions/categories`

Both require `role.view`. Return fixed registry rows from database ordered by category, code.

**Step 4: Implement roles handler**

Create `go-server/internal/api/roles.go`:

- `GET /api/roles`
- `GET /api/roles/:id`
- `POST /api/roles`
- `PATCH /api/roles/:id`
- `DELETE /api/roles/:id`
- `GET /api/roles/:id/permissions`
- `POST /api/roles/:id/permissions`
- `DELETE /api/roles/:id/permissions/:permissionId`

Permission checks:

- read: `role.view`
- create: `role.create`
- update/activate/deactivate: `role.update`
- delete: `role.delete`
- assign/remove permissions: `role.assign_permissions`

Deletion is allowed only for non-system roles with no user, group, or organization-unit assignments.

**Step 5: Wire routes**

Modify `router.go` with `/api/roles` and `/api/permissions` groups protected by trusted internal token.

**Step 6: Run tests**

```sh
cd go-server
rtk go test ./internal/api -run 'RoleAPI|PermissionAPI'
```

Expected: PASS.

**Step 7: Commit**

```sh
rtk git add go-server/internal/api/roles.go go-server/internal/api/roles_test.go go-server/internal/api/permissions.go go-server/internal/api/permissions_test.go go-server/internal/api/router.go
rtk git commit -m "Add role and permission APIs"
```

## Task 8: Group Backend API

**Files:**
- Create: `go-server/internal/api/groups.go`
- Create: `go-server/internal/api/groups_test.go`
- Modify: `go-server/internal/api/router.go`

**Step 1: Write failing group API tests**

Cover:

- list requires `group.view`;
- create requires `group.create`;
- update requires `group.update`;
- delete requires `group.delete`;
- add/remove users requires `group.manage_users`;
- assign/remove roles requires `group.assign_roles`;
- deleting a group with members or roles returns `409`;
- duplicate group code returns `409`.

**Step 2: Run tests to verify failure**

```sh
cd go-server
rtk go test ./internal/api -run GroupAPI
```

Expected: FAIL because routes do not exist.

**Step 3: Implement groups handler**

Create `go-server/internal/api/groups.go`:

Routes:

```go
GET    /api/groups
GET    /api/groups/:id
POST   /api/groups
PATCH  /api/groups/:id
DELETE /api/groups/:id
POST   /api/groups/:id/users
DELETE /api/groups/:id/users/:userId
POST   /api/groups/:id/roles
DELETE /api/groups/:id/roles/:assignmentId
```

Use the same scope validation helper as user role assignment.

**Step 4: Wire routes**

Modify `router.go` with `/api/groups` group.

**Step 5: Run tests**

```sh
cd go-server
rtk go test ./internal/api -run GroupAPI
```

Expected: PASS.

**Step 6: Commit**

```sh
rtk git add go-server/internal/api/groups.go go-server/internal/api/groups_test.go go-server/internal/api/router.go
rtk git commit -m "Add group management API"
```

## Task 9: Current User Authorization API And Organization Unit Enforcement

**Files:**
- Create: `go-server/internal/api/me.go`
- Create: `go-server/internal/api/me_test.go`
- Modify: `go-server/internal/api/organization_units.go`
- Modify: `go-server/internal/api/organization_units_test.go`
- Modify: `go-server/internal/api/router.go`

**Step 1: Write failing current-user API tests**

Create `go-server/internal/api/me_test.go`:

```go
func TestMePermissionsReturnsEffectivePermissions(t *testing.T) {
	router, db := newAuthTestRouter(t)
	admin := createAdminUser(t, db, "admin@example.com", "password123")
	if err := rbac.BootstrapLegacyAdmins(db); err != nil {
		t.Fatalf("bootstrap: %v", err)
	}
	token := loginUser(t, router, admin.Email, "password123")

	response := authJSON(t, router, http.MethodGet, "/api/me/permissions", "", authCookieHeaders(token))
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), "organization_unit.manage_hierarchy") {
		t.Fatalf("body = %s, want org unit permission", response.Body.String())
	}
}
```

**Step 2: Update organization-unit tests**

Change expectations so a legacy admin only succeeds after RBAC bootstrap. Add a regular active user assigned only `organization_unit.view` and assert:

- tree read returns `200`;
- create/update/move/archive return `403`.

**Step 3: Run tests to verify failure**

```sh
cd go-server
rtk go test ./internal/api -run 'Me|OrganizationUnit'
```

Expected: FAIL because `/api/me` does not exist and org-unit routes still use enum checks.

**Step 4: Implement me handler**

Create `go-server/internal/api/me.go` with:

- `GET /api/me`
- `GET /api/me/permissions`
- `POST /api/auth/check-permission`

Response for permissions:

```json
{
  "permissions": [
    {
      "code": "organization_unit.manage_hierarchy",
      "scopeType": "global",
      "organizationUnitId": null,
      "source": "user_role"
    }
  ]
}
```

**Step 5: Replace organization-unit enum checks**

Modify `organization_units.go`:

- `requireUser` should call shared `requirePermission` with `organization_unit.view` for tree reads, or `requireAuthenticatedUser` plus resolver check.
- `listTree`: require `organization_unit.view`.
- `listArchived`: require `organization_unit.view_audit` or `organization_unit.manage_hierarchy`.
- `create`: require `organization_unit.create` or `organization_unit.manage_hierarchy`.
- `update`: require `organization_unit.update` or `organization_unit.manage_hierarchy`.
- `move`: require `organization_unit.manage_hierarchy`.
- `archive`: require `organization_unit.delete` or `organization_unit.manage_hierarchy`.

If adding OR permission support is cleaner, add helper:

```go
func requireAnyPermission(c *gin.Context, h *authHandler, permissions []string, organizationUnitID *string) (auth.User, bool)
```

**Step 6: Wire routes**

Modify `router.go`:

```go
me := newMeHandler(options, auth)
meAPI := router.Group("/api")
meAPI.Use(auth.requireTrustedInternalRequest())
meAPI.GET("/me", me.getMe)
meAPI.GET("/me/permissions", me.getPermissions)
meAPI.POST("/auth/check-permission", me.checkPermission)
```

Do not conflict with existing `/api/auth` group.

**Step 7: Run tests**

```sh
cd go-server
rtk go test ./internal/api -run 'Me|OrganizationUnit'
```

Expected: PASS.

**Step 8: Commit**

```sh
rtk git add go-server/internal/api/me.go go-server/internal/api/me_test.go go-server/internal/api/organization_units.go go-server/internal/api/organization_units_test.go go-server/internal/api/router.go
rtk git commit -m "Enforce organization unit permissions"
```

## Task 10: Swagger Documentation

**Files:**
- Modify: `go-server/internal/api/swagger_doc.go`
- Modify generated files under `go-server/docs/`
- Modify: `go-server/internal/api/swagger_doc_test.go`

**Step 1: Update Swagger tests**

Extend `swagger_doc_test.go` to assert the generated doc includes:

- `/api/users`
- `/api/roles`
- `/api/permissions`
- `/api/groups`
- `/api/me`
- `/api/me/permissions`
- `/api/auth/check-permission`

**Step 2: Run Swagger test to verify failure**

```sh
cd go-server
rtk go test ./internal/api -run Swagger
```

Expected: FAIL because docs do not include new routes.

**Step 3: Add Swagger annotations/docs**

Modify `go-server/internal/api/swagger_doc.go` consistently with existing docs style.

**Step 4: Regenerate Swagger**

```sh
cd go-server
rtk go run ./cmd/syncra swagger
```

Expected: `go-server/docs/swagger.json` and related generated docs update.

**Step 5: Run tests**

```sh
cd go-server
rtk go test ./internal/api -run Swagger
```

Expected: PASS.

**Step 6: Commit**

```sh
rtk git add go-server/internal/api/swagger_doc.go go-server/internal/api/swagger_doc_test.go go-server/docs
rtk git commit -m "Document RBAC APIs"
```

## Task 11: Frontend Server RBAC Clients

**Files:**
- Create: `frontend/src/lib/server/rbac.ts`
- Create: `frontend/src/lib/server/rbac.test.ts`
- Modify: `frontend/src/lib/server/auth.ts`
- Modify: `frontend/src/app.d.ts` if locals types are defined there

**Step 1: Write failing server-client tests**

Create `frontend/src/lib/server/rbac.test.ts` covering:

- internal token header is attached;
- cookie header is forwarded;
- backend error maps to `RbacApiError`;
- permissions response is validated;
- invalid response throws `RbacApiError(502, ...)`.

Test shape:

```ts
import { describe, expect, it, vi } from 'vitest';
import { getMyPermissions } from './rbac';

describe('rbac server client', () => {
	it('fetches current permissions with internal and cookie headers', async () => {
		const fetchMock = vi.fn().mockResolvedValue(
			new Response(JSON.stringify({ permissions: [] }), { status: 200 })
		);

		await getMyPermissions(fetchMock, 'auth.session_token=token');

		expect(fetchMock).toHaveBeenCalledWith(
			expect.stringContaining('/api/me/permissions'),
			expect.objectContaining({
				headers: expect.any(Headers)
			})
		);
	});
});
```

**Step 2: Run test to verify failure**

```sh
cd frontend
rtk pnpm test -- src/lib/server/rbac.test.ts
```

Expected: FAIL because `rbac.ts` does not exist.

**Step 3: Implement server client**

Create `frontend/src/lib/server/rbac.ts`:

- `getMe(fetchFn, cookieHeader)`
- `getMyPermissions(fetchFn, cookieHeader)`
- `checkPermission(fetchFn, cookieHeader, input)`
- user CRUD functions
- role CRUD and permission assignment functions
- permission list/category functions
- group CRUD/member/role functions
- response validators
- `RbacApiError` class and `isRbacApiError`

Use the existing `apiBaseUrl`, `internalAPIHeaders`, and public error helpers.

**Step 4: Extend auth types**

Modify `frontend/src/lib/server/auth.ts` `AuthUser`:

```ts
status: 'invited' | 'active' | 'inactive' | 'suspended' | 'deleted';
primaryOrganizationUnitId?: string | null;
managerUserId?: string | null;
jobTitle?: string | null;
phone?: string | null;
permissions?: string[];
```

Only include fields actually returned by Go auth session JSON after backend is updated.

**Step 5: Run tests**

```sh
cd frontend
rtk pnpm test -- src/lib/server/rbac.test.ts
```

Expected: PASS.

**Step 6: Commit**

```sh
rtk git add frontend/src/lib/server/rbac.ts frontend/src/lib/server/rbac.test.ts frontend/src/lib/server/auth.ts frontend/src/app.d.ts
rtk git commit -m "Add frontend RBAC server client"
```

## Task 12: Frontend Svelte API Proxy Routes

**Files:**
- Create: `frontend/src/routes/api/rbac/api.server.ts`
- Create: `frontend/src/routes/api/rbac/server.test.ts`
- Create route files under:
  - `frontend/src/routes/api/users/`
  - `frontend/src/routes/api/roles/`
  - `frontend/src/routes/api/permissions/`
  - `frontend/src/routes/api/groups/`
  - `frontend/src/routes/api/me/`

**Step 1: Write failing proxy route tests**

Create `frontend/src/routes/api/rbac/server.test.ts` using the same mock style as `frontend/src/routes/api/organization-units/server.test.ts`.

Cover:

- unauthenticated requests return `401`;
- requests without required local permission return `403`;
- allowed requests call `$lib/server/rbac` with `fetch` and cookie header;
- known backend errors map to public-safe JSON.

**Step 2: Run tests to verify failure**

```sh
cd frontend
rtk pnpm test -- src/routes/api/rbac/server.test.ts
```

Expected: FAIL because routes do not exist.

**Step 3: Implement shared API helpers**

Create `frontend/src/routes/api/rbac/api.server.ts`:

```ts
import { json } from '@sveltejs/kit';
import { isRbacApiError } from '$lib/server/rbac';
import { publicErrorMessage, publicErrorStatus } from '$lib/server/public-errors';

export function requireAuthenticatedUser(locals: App.Locals) {
	if (locals.user) return null;
	return json({ error: 'Authentication required' }, { status: 401 });
}

export function requireLocalPermission(locals: App.Locals, permission: string) {
	const authError = requireAuthenticatedUser(locals);
	if (authError) return authError;
	if (!locals.permissions?.includes(permission) && !locals.permissions?.includes('system.admin')) {
		return json({ error: 'permission required' }, { status: 403 });
	}
	return null;
}

export function cookieHeader(request: Request) {
	return request.headers.get('cookie');
}

export function rbacAPIErrorResponse(error: unknown, fallback: string) {
	if (isRbacApiError(error)) {
		return json(
			{ error: publicErrorMessage(error.status, error.message, fallback) },
			{ status: publicErrorStatus(error.status) }
		);
	}
	throw error;
}
```

**Step 4: Implement route files**

Add route handlers that call the server client. Keep request parsing small and typed.

Examples:

- `frontend/src/routes/api/users/+server.ts`: `GET`, `POST`
- `frontend/src/routes/api/users/[id]/+server.ts`: `GET`, `PATCH`, `DELETE`
- `frontend/src/routes/api/users/[id]/roles/+server.ts`: `POST`
- `frontend/src/routes/api/users/[id]/roles/[assignmentId]/+server.ts`: `DELETE`
- `frontend/src/routes/api/roles/+server.ts`: `GET`, `POST`
- `frontend/src/routes/api/roles/[id]/permissions/+server.ts`: `GET`, `POST`
- `frontend/src/routes/api/permissions/+server.ts`: `GET`
- `frontend/src/routes/api/permissions/categories/+server.ts`: `GET`
- `frontend/src/routes/api/groups/+server.ts`: `GET`, `POST`
- `frontend/src/routes/api/me/+server.ts`: `GET`
- `frontend/src/routes/api/me/permissions/+server.ts`: `GET`

Use frontend-local permission checks for clean UX, but rely on Go for authority.

**Step 5: Run tests**

```sh
cd frontend
rtk pnpm test -- src/routes/api/rbac/server.test.ts
```

Expected: PASS.

**Step 6: Commit**

```sh
rtk git add frontend/src/routes/api/rbac frontend/src/routes/api/users frontend/src/routes/api/roles frontend/src/routes/api/permissions frontend/src/routes/api/groups frontend/src/routes/api/me
rtk git commit -m "Add RBAC Svelte API proxies"
```

## Task 13: Session Locals And App Navigation Permissions

**Files:**
- Modify: `frontend/src/hooks.server.ts`
- Modify: `frontend/src/app.d.ts`
- Modify: `frontend/src/routes/app/+layout.server.ts`
- Modify: `frontend/src/routes/app/+layout.svelte`
- Modify: `frontend/src/lib/components/app-sidebar.svelte`
- Modify: `frontend/src/routes/app/organization-units/+page.server.ts`
- Modify: `frontend/src/routes/app/organization-units/page.server.test.ts`
- Create: `frontend/src/routes/layout-permissions.test.ts`

**Step 1: Write failing tests**

Create or extend tests to assert:

- app layout load returns public permissions, not session token;
- admin nav section appears only when permissions include `system.admin`, `user.view`, `role.view`, or `group.view`;
- organization-units page uses `organization_unit.*` permissions instead of `role === 'admin'`.

**Step 2: Run tests to verify failure**

```sh
cd frontend
rtk pnpm test -- src/routes/layout-permissions.test.ts src/routes/app/organization-units/page.server.test.ts
```

Expected: FAIL because locals do not include permissions.

**Step 3: Load permissions in hook**

Modify `frontend/src/hooks.server.ts` after successful session load:

```ts
const permissions = auth ? await getMyPermissions(event.fetch, cookieHeader) : null;
event.locals.permissions = permissions?.permissions.map((permission) => permission.code) ?? [];
```

If permission loading fails for a protected route, redirect to login only if the session is invalid. For RBAC service failures, keep user but no permissions and let page/API routes show errors.

**Step 4: Update app types**

Modify `frontend/src/app.d.ts`:

```ts
declare global {
	namespace App {
		interface Locals {
			user: import('$lib/server/auth').AuthUser | null;
			session: import('$lib/server/auth').AuthSession | null;
			permissions: string[];
		}
	}
}
```

**Step 5: Update layout data and sidebar**

Add permissions to layout server public data:

```ts
permissions: locals.permissions ?? []
```

Update `AppSidebar` props and nav filtering. Add Admin section links only if relevant permissions exist:

- Users: `user.view`
- Roles: `role.view`
- Permissions: `role.view` or `system.admin`
- Groups: `group.view`

**Step 6: Update organization-units page load**

Replace:

```ts
canManageOrganizationUnits: locals.user?.role === 'admin'
```

With:

```ts
const permissions = locals.permissions ?? [];
canManageOrganizationUnits:
	permissions.includes('system.admin') ||
	permissions.includes('organization_unit.manage_hierarchy') ||
	permissions.includes('organization_unit.create') ||
	permissions.includes('organization_unit.update') ||
	permissions.includes('organization_unit.delete')
```

**Step 7: Run tests**

```sh
cd frontend
rtk pnpm test -- src/routes/layout-permissions.test.ts src/routes/app/organization-units/page.server.test.ts
```

Expected: PASS.

**Step 8: Commit**

```sh
rtk git add frontend/src/hooks.server.ts frontend/src/app.d.ts frontend/src/routes/app/+layout.server.ts frontend/src/routes/app/+layout.svelte frontend/src/lib/components/app-sidebar.svelte frontend/src/routes/app/organization-units/+page.server.ts frontend/src/routes/app/organization-units/page.server.test.ts frontend/src/routes/layout-permissions.test.ts
rtk git commit -m "Gate app navigation with RBAC permissions"
```

## Task 14: Admin Users Frontend

**Files:**
- Create: `frontend/src/routes/app/admin/users/api.ts`
- Create: `frontend/src/routes/app/admin/users/api.test.ts`
- Create: `frontend/src/routes/app/admin/users/+page.server.ts`
- Create: `frontend/src/routes/app/admin/users/page.server.test.ts`
- Create: `frontend/src/routes/app/admin/users/+page.svelte`
- Create: `frontend/src/routes/app/admin/users/user-editor.svelte`
- Create: `frontend/src/routes/app/admin/users/user-role-assignments.svelte`

**Step 1: Write failing frontend tests**

Cover:

- API client fetches `/api/users`;
- create/update/status calls use correct methods and payloads;
- page server load exposes `canManageUsers`;
- source imports TanStack query/mutations and does not import server-only modules.

**Step 2: Run tests to verify failure**

```sh
cd frontend
rtk pnpm test -- src/routes/app/admin/users
```

Expected: FAIL because files do not exist.

**Step 3: Implement browser API helper**

Create `api.ts` with fetch helpers against local Svelte routes:

- `fetchUsers`
- `createUser`
- `updateUser`
- `activateUser`
- `deactivateUser`
- `suspendUser`
- `softDeleteUser`
- `assignUserRole`
- `removeUserRole`
- `assignUserGroup`
- `removeUserGroup`
- `setPrimaryOrganizationUnit`

Use response validation like organization-unit client helpers.

**Step 4: Implement page server load**

Return only permission flags and any selected IDs. Do not load full user data server-side if the page uses TanStack query.

**Step 5: Implement Svelte components**

Use a dense table/list plus a right-side editor panel:

- filter by name/email/status;
- create user form;
- status controls with confirmation for suspend/delete;
- primary org-unit selector;
- scoped role assignment selector;
- group assignment selector.

Use existing UI primitives and lucide icons.

**Step 6: Run Svelte checks/autofixer**

Use Svelte tooling if available, then run:

```sh
cd frontend
rtk pnpm check
rtk pnpm test -- src/routes/app/admin/users
```

Expected: PASS.

**Step 7: Commit**

```sh
rtk git add frontend/src/routes/app/admin/users
rtk git commit -m "Add admin user management page"
```

## Task 15: Admin Roles And Permissions Frontend

**Files:**
- Create: `frontend/src/routes/app/admin/roles/api.ts`
- Create: `frontend/src/routes/app/admin/roles/api.test.ts`
- Create: `frontend/src/routes/app/admin/roles/+page.server.ts`
- Create: `frontend/src/routes/app/admin/roles/page.server.test.ts`
- Create: `frontend/src/routes/app/admin/roles/+page.svelte`
- Create: `frontend/src/routes/app/admin/roles/permission-matrix.svelte`
- Create: `frontend/src/routes/app/admin/permissions/+page.server.ts`
- Create: `frontend/src/routes/app/admin/permissions/+page.svelte`
- Create: `frontend/src/routes/app/admin/permissions/page.server.test.ts`

**Step 1: Write failing tests**

Cover:

- roles API helper calls local Svelte routes;
- permission matrix groups permissions by category;
- page load flags require `role.view`;
- permissions page is read-only.

**Step 2: Run tests to verify failure**

```sh
cd frontend
rtk pnpm test -- src/routes/app/admin/roles src/routes/app/admin/permissions
```

Expected: FAIL because files do not exist.

**Step 3: Implement role API helper and page**

Role page should support:

- list roles;
- create custom role;
- edit name/description/isActive for non-system roles;
- assign/remove permissions via matrix;
- show system roles as locked for destructive operations.

**Step 4: Implement read-only permissions page**

Group permissions by category and show code, name, description.

**Step 5: Run checks**

```sh
cd frontend
rtk pnpm check
rtk pnpm test -- src/routes/app/admin/roles src/routes/app/admin/permissions
```

Expected: PASS.

**Step 6: Commit**

```sh
rtk git add frontend/src/routes/app/admin/roles frontend/src/routes/app/admin/permissions
rtk git commit -m "Add admin role and permission pages"
```

## Task 16: Admin Groups Frontend

**Files:**
- Create: `frontend/src/routes/app/admin/groups/api.ts`
- Create: `frontend/src/routes/app/admin/groups/api.test.ts`
- Create: `frontend/src/routes/app/admin/groups/+page.server.ts`
- Create: `frontend/src/routes/app/admin/groups/page.server.test.ts`
- Create: `frontend/src/routes/app/admin/groups/+page.svelte`
- Create: `frontend/src/routes/app/admin/groups/group-editor.svelte`
- Create: `frontend/src/routes/app/admin/groups/group-members.svelte`
- Create: `frontend/src/routes/app/admin/groups/group-role-assignments.svelte`

**Step 1: Write failing tests**

Cover:

- fetch/create/update/delete groups request shapes;
- add/remove users request shapes;
- assign/remove roles request shapes;
- page load exposes permission flags.

**Step 2: Run tests to verify failure**

```sh
cd frontend
rtk pnpm test -- src/routes/app/admin/groups
```

Expected: FAIL because files do not exist.

**Step 3: Implement group page**

Group page should support:

- list active groups;
- create/edit group;
- optional org-unit selector;
- add/remove members;
- assign/remove scoped roles;
- prevent delete UI when group has members or roles if that count is available.

**Step 4: Run checks**

```sh
cd frontend
rtk pnpm check
rtk pnpm test -- src/routes/app/admin/groups
```

Expected: PASS.

**Step 5: Commit**

```sh
rtk git add frontend/src/routes/app/admin/groups
rtk git commit -m "Add admin group management page"
```

## Task 17: Full Verification And Browser Check

**Files:**
- Modify only focused defects found by verification.

**Step 1: Run backend verification**

```sh
cd go-server
rtk go test ./...
rtk atlas migrate validate --dir file://migrations
```

Expected: PASS.

**Step 2: Run frontend verification**

```sh
cd frontend
rtk pnpm test
rtk pnpm check
rtk pnpm build
```

Expected: PASS.

**Step 3: Regenerate Swagger if needed**

If any API annotations changed after Task 10:

```sh
cd go-server
rtk go run ./cmd/syncra swagger
rtk go test ./internal/api -run Swagger
```

Expected: PASS and generated docs are current.

**Step 4: Start local services**

Terminal 1:

```sh
cd go-server
rtk go run ./cmd/syncra api
```

Terminal 2:

```sh
cd frontend
rtk pnpm dev -- --host 127.0.0.1
```

Expected: Go API serves on `localhost:8080`; SvelteKit prints a local Vite URL.

**Step 5: Browser inspect key flows**

Verify:

- legacy admin can sign in and sees Admin nav;
- admin can create a role and assign permissions;
- admin can create a group, add a user, and assign a scoped role;
- admin can update user status and suspended users lose access;
- organization-unit mutation controls are permission-gated;
- regular viewer can browse allowed organization units but cannot mutate.

**Step 6: Check git status**

```sh
rtk git status --short
```

Expected: only intentional RBAC files remain changed.

**Step 7: Commit verification fixes if any**

If verification required fixes:

```sh
rtk git add go-server frontend
rtk git commit -m "Fix RBAC verification issues"
```

Do not create an empty commit.

## Execution Notes

- Keep commits small and task-scoped.
- If a task exposes a simpler local design, update this plan before implementing divergent behavior.
- Preserve unrelated user changes.
- Prefer shared validation helpers only after duplication appears in at least two handlers.
- Do not remove the legacy `auth.User.Role` field in this implementation pass.
