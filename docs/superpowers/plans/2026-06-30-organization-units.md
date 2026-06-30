# Organization Units Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a shared, admin-managed Organization Units hierarchy in the Go API and SvelteKit frontend.

**Architecture:** Add a focused `orgunits` backend domain package, expose trusted internal API endpoints under `/api/organization-units`, and consume them from server-only SvelteKit helpers. The frontend adds one protected `/app/organization-units` page with a tree + details split; admins see mutation forms, regular users see read-only details.

**Tech Stack:** Go, Gin, GORM, PostgreSQL, Atlas migrations, SvelteKit, Svelte 5 runes, TypeScript, Vitest, Tailwind CSS, local shadcn-style UI primitives.

---

## Scope Check

This plan implements one coherent feature: Organization Units. It does not add per-unit membership, document linking, restore flows, import/export, OCR, billing, email, PDF, or a separate admin portal.

## File Structure

- Create `go-server/internal/orgunits/unit.go`: GORM model, code normalization, validation constants.
- Create `go-server/internal/orgunits/tree.go`: pure tree-building and descendant helpers.
- Create `go-server/internal/orgunits/tree_test.go`: behavior tests for normalization, ordering, tree shape, and descendant exclusion.
- Modify `go-server/internal/database/database.go`: register `orgunits.Unit` for schema generation.
- Create `go-server/migrations/20260630120000_add_organization_units.sql`: table, self-reference, indexes, partial unique active code index.
- Modify `go-server/migrations/atlas.sum`: refresh checksum after adding the migration.
- Create `go-server/internal/api/organization_units_test.go`: API behavior tests.
- Create `go-server/internal/api/organization_units.go`: request/response structs, auth/role checks, handlers.
- Modify `go-server/internal/api/router.go`: wire Organization Unit routes.
- Modify `go-server/internal/api/swagger_doc.go`: document API operations.
- Create `frontend/src/lib/server/organization-units.ts`: server-only Go API client and response validation.
- Create `frontend/src/lib/server/organization-units.test.ts`: API client tests.
- Create `frontend/src/routes/app/+layout.server.ts`: app layout data from `locals`.
- Create `frontend/src/routes/app/+layout.svelte`: shared protected app chrome and navigation.
- Modify `frontend/src/routes/app/+page.svelte`: remove duplicate shell chrome and add an Organization Units entry point.
- Create `frontend/src/routes/app/organization-units/tree.ts`: frontend tree utilities.
- Create `frontend/src/routes/app/organization-units/tree.test.ts`: frontend tree utility tests.
- Create `frontend/src/routes/app/organization-units/+page.server.ts`: page load and server actions.
- Create `frontend/src/routes/app/organization-units/page.server.test.ts`: page load/action tests.
- Create `frontend/src/routes/app/organization-units/unit-tree.svelte`: recursive-ish flat tree renderer using depth metadata.
- Create `frontend/src/routes/app/organization-units/unit-details-panel.svelte`: details, edit, move, and archive forms.
- Create `frontend/src/routes/app/organization-units/+page.svelte`: page composition and selection state.

## Task 1: Backend Domain Model, Tree Helpers, And Migration

**Files:**
- Create: `go-server/internal/orgunits/unit.go`
- Create: `go-server/internal/orgunits/tree.go`
- Create: `go-server/internal/orgunits/tree_test.go`
- Modify: `go-server/internal/database/database.go`
- Create: `go-server/migrations/20260630120000_add_organization_units.sql`
- Modify: `go-server/migrations/atlas.sum`

- [ ] **Step 1: Write failing domain tests**

Create `go-server/internal/orgunits/tree_test.go`:

```go
package orgunits

import (
	"testing"
	"time"
)

func TestNormalizeCode(t *testing.T) {
	code, err := NormalizeCode(" fin-ap ")
	if err != nil {
		t.Fatalf("NormalizeCode() error = %v", err)
	}
	if code == nil || *code != "FIN-AP" {
		t.Fatalf("NormalizeCode() = %#v, want FIN-AP", code)
	}

	empty, err := NormalizeCode("   ")
	if err != nil {
		t.Fatalf("NormalizeCode(empty) error = %v", err)
	}
	if empty != nil {
		t.Fatalf("NormalizeCode(empty) = %#v, want nil", empty)
	}

	if _, err := NormalizeCode("THIS-CODE-IS-MORE-THAN-FORTY-CHARACTERS-LONG"); err == nil {
		t.Fatal("NormalizeCode(long) error = nil, want error")
	}
}

func TestBuildTreeOrdersSiblingsAndNestsChildren(t *testing.T) {
	now := time.Date(2026, 6, 30, 12, 0, 0, 0, time.UTC)
	rootB := Unit{ID: "00000000-0000-0000-0000-000000000002", Name: "Operations", CreatedAt: now, UpdatedAt: now}
	rootA := Unit{ID: "00000000-0000-0000-0000-000000000001", Name: "Company", CreatedAt: now, UpdatedAt: now}
	childB := Unit{ID: "00000000-0000-0000-0000-000000000004", ParentID: &rootA.ID, Name: "Operations", CreatedAt: now, UpdatedAt: now}
	childA := Unit{ID: "00000000-0000-0000-0000-000000000003", ParentID: &rootA.ID, Name: "Finance", CreatedAt: now, UpdatedAt: now}

	tree := BuildTree([]Unit{rootB, childB, rootA, childA})

	if len(tree) != 2 {
		t.Fatalf("root count = %d, want 2", len(tree))
	}
	if tree[0].Name != "Company" || tree[1].Name != "Operations" {
		t.Fatalf("root order = %q, %q", tree[0].Name, tree[1].Name)
	}
	if len(tree[0].Children) != 2 {
		t.Fatalf("Company child count = %d, want 2", len(tree[0].Children))
	}
	if tree[0].Children[0].Name != "Finance" || tree[0].Children[1].Name != "Operations" {
		t.Fatalf("child order = %q, %q", tree[0].Children[0].Name, tree[0].Children[1].Name)
	}
}

func TestDescendantIDs(t *testing.T) {
	rootID := "00000000-0000-0000-0000-000000000001"
	childID := "00000000-0000-0000-0000-000000000002"
	grandchildID := "00000000-0000-0000-0000-000000000003"
	otherID := "00000000-0000-0000-0000-000000000004"

	ids := DescendantIDs(rootID, []Unit{
		{ID: rootID, Name: "Root"},
		{ID: childID, ParentID: &rootID, Name: "Child"},
		{ID: grandchildID, ParentID: &childID, Name: "Grandchild"},
		{ID: otherID, Name: "Other"},
	})

	if !ids[childID] || !ids[grandchildID] {
		t.Fatalf("descendant ids = %#v, want child and grandchild", ids)
	}
	if ids[rootID] || ids[otherID] {
		t.Fatalf("descendant ids = %#v, did not expect root or other", ids)
	}
}
```

- [ ] **Step 2: Run domain tests to verify failure**

Run:

```sh
cd go-server
rtk go test ./internal/orgunits
```

Expected: FAIL because package `internal/orgunits` or identifiers `Unit`, `NormalizeCode`, `BuildTree`, and `DescendantIDs` do not exist.

- [ ] **Step 3: Implement the domain model**

Create `go-server/internal/orgunits/unit.go`:

```go
package orgunits

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	MaxNameCharacters = 160
	MaxCodeCharacters = 40
)

type Unit struct {
	ID          string     `gorm:"type:uuid;primaryKey" json:"id"`
	ParentID    *string    `gorm:"column:parent_id;type:uuid;index:idx_organization_units_parent_name_id,priority:1" json:"parent_id,omitempty"`
	Parent      *Unit      `gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
	Name        string     `gorm:"not null;size:160;index:idx_organization_units_parent_name_id,priority:2" json:"name"`
	Code        *string    `gorm:"size:40;index" json:"code,omitempty"`
	Description *string    `gorm:"type:text" json:"description,omitempty"`
	ArchivedAt  *time.Time `gorm:"column:archived_at;index" json:"archived_at,omitempty"`
	CreatedAt   time.Time  `gorm:"column:created_at;not null;index:idx_organization_units_parent_name_id,priority:3" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;not null" json:"updated_at"`
}

func (u *Unit) BeforeCreate(_ *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	return nil
}

func (Unit) TableName() string {
	return "organization_units"
}

func NormalizeName(raw string) (string, error) {
	name := strings.TrimSpace(raw)
	if name == "" {
		return "", errors.New("name is required")
	}
	if utf8.RuneCountInString(name) > MaxNameCharacters {
		return "", errors.New("name must be at most 160 characters")
	}
	return name, nil
}

func NormalizeCode(raw string) (*string, error) {
	code := strings.ToUpper(strings.TrimSpace(raw))
	if code == "" {
		return nil, nil
	}
	if utf8.RuneCountInString(code) > MaxCodeCharacters {
		return nil, errors.New("code must be at most 40 characters")
	}
	return &code, nil
}

func NormalizeDescription(raw string) *string {
	description := strings.TrimSpace(raw)
	if description == "" {
		return nil
	}
	return &description
}
```

- [ ] **Step 4: Implement tree helpers**

Create `go-server/internal/orgunits/tree.go`:

```go
package orgunits

import (
	"sort"
	"time"
)

type TreeNode struct {
	ID          string     `json:"id"`
	ParentID    *string    `json:"parentId,omitempty"`
	Name        string     `json:"name"`
	Code        *string    `json:"code,omitempty"`
	Description *string    `json:"description,omitempty"`
	CreatedAt   string     `json:"createdAt"`
	UpdatedAt   string     `json:"updatedAt"`
	Children    []TreeNode `json:"children"`
}

func BuildTree(units []Unit) []TreeNode {
	childrenByParent := map[string][]Unit{}
	var roots []Unit
	for _, unit := range units {
		if unit.ParentID == nil || *unit.ParentID == "" {
			roots = append(roots, unit)
			continue
		}
		childrenByParent[*unit.ParentID] = append(childrenByParent[*unit.ParentID], unit)
	}
	sortUnits(roots)
	for parentID := range childrenByParent {
		sortUnits(childrenByParent[parentID])
	}
	return buildNodes(roots, childrenByParent)
}

func DescendantIDs(rootID string, units []Unit) map[string]bool {
	childrenByParent := map[string][]Unit{}
	for _, unit := range units {
		if unit.ParentID != nil {
			childrenByParent[*unit.ParentID] = append(childrenByParent[*unit.ParentID], unit)
		}
	}
	out := map[string]bool{}
	var visit func(string)
	visit = func(parentID string) {
		for _, child := range childrenByParent[parentID] {
			if out[child.ID] {
				continue
			}
			out[child.ID] = true
			visit(child.ID)
		}
	}
	visit(rootID)
	return out
}

func buildNodes(units []Unit, childrenByParent map[string][]Unit) []TreeNode {
	nodes := make([]TreeNode, 0, len(units))
	for _, unit := range units {
		nodes = append(nodes, TreeNode{
			ID:          unit.ID,
			ParentID:    unit.ParentID,
			Name:        unit.Name,
			Code:        unit.Code,
			Description: unit.Description,
			CreatedAt:   unit.CreatedAt.UTC().Format(time.RFC3339Nano),
			UpdatedAt:   unit.UpdatedAt.UTC().Format(time.RFC3339Nano),
			Children:    buildNodes(childrenByParent[unit.ID], childrenByParent),
		})
	}
	return nodes
}

func sortUnits(units []Unit) {
	sort.SliceStable(units, func(i, j int) bool {
		if units[i].Name == units[j].Name {
			return units[i].ID < units[j].ID
		}
		return units[i].Name < units[j].Name
	})
}
```

- [ ] **Step 5: Register the model**

Modify `go-server/internal/database/database.go`:

```go
import (
	"errors"
	"strings"

	"ai.ro/syncra/dms/internal/auth"
	"ai.ro/syncra/dms/internal/orgunits"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)
```

and add `&orgunits.Unit{}` to `ApplicationModels()`:

```go
func ApplicationModels() []any {
	return []any{
		&auth.User{},
		&auth.AuthAccount{},
		&auth.Session{},
		&auth.Verification{},
		&orgunits.Unit{},
	}
}
```

- [ ] **Step 6: Add migration SQL**

Create `go-server/migrations/20260630120000_add_organization_units.sql`:

```sql
CREATE TABLE "organization_units" (
  "id" uuid,
  "parent_id" uuid,
  "name" varchar(160) NOT NULL,
  "code" varchar(40),
  "description" text,
  "archived_at" timestamptz,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);

CREATE INDEX "idx_organization_units_parent_name_id" ON "organization_units" ("parent_id", "name", "id");
CREATE INDEX "idx_organization_units_archived_at" ON "organization_units" ("archived_at");
CREATE INDEX "idx_organization_units_code" ON "organization_units" ("code");
CREATE UNIQUE INDEX "idx_organization_units_active_code_unique" ON "organization_units" ("code") WHERE "code" IS NOT NULL AND "archived_at" IS NULL;

ALTER TABLE "organization_units"
  ADD CONSTRAINT "fk_organization_units_parent"
  FOREIGN KEY ("parent_id") REFERENCES "organization_units"("id")
  ON DELETE RESTRICT ON UPDATE CASCADE;
```

- [ ] **Step 7: Refresh migration checksum**

Run:

```sh
cd go-server
rtk atlas migrate hash --dir file://migrations
```

Expected: `migrations/atlas.sum` changes and includes the new migration.

- [ ] **Step 8: Run domain tests and migration validation**

Run:

```sh
cd go-server
rtk go test ./internal/orgunits ./internal/database
rtk atlas migrate validate --dir file://migrations
```

Expected: tests PASS and Atlas validation PASS.

- [ ] **Step 9: Commit backend domain and migration**

Run:

```sh
rtk git add go-server/internal/orgunits go-server/internal/database/database.go go-server/migrations/20260630120000_add_organization_units.sql go-server/migrations/atlas.sum
rtk git commit -m "Add organization units model"
```

## Task 2: Backend API Behavior Tests

**Files:**
- Create: `go-server/internal/api/organization_units_test.go`
- Modify: `go-server/internal/api/auth_handlers_test.go`

- [ ] **Step 1: Extend the test router migration**

Modify `newAuthTestRouterWithOptions` in `go-server/internal/api/auth_handlers_test.go` so SQLite test routers migrate Organization Units:

```go
if err := db.AutoMigrate(&auth.User{}, &auth.AuthAccount{}, &auth.Session{}, &auth.Verification{}, &orgunits.Unit{}); err != nil {
	t.Fatalf("auto migrate: %v", err)
}
```

Add this import:

```go
import "ai.ro/syncra/dms/internal/orgunits"
```

- [ ] **Step 2: Write failing API tests**

Create `go-server/internal/api/organization_units_test.go`:

```go
package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"ai.ro/syncra/dms/internal/auth"
	"ai.ro/syncra/dms/internal/orgunits"
	"gorm.io/gorm"
)

func TestOrganizationUnitRoutesRequireSessionAndAdminForMutations(t *testing.T) {
	router, db := newAuthTestRouter(t)

	unauthenticated := orgUnitJSON(t, router, http.MethodGet, "/api/organization-units/tree", "", map[string]string{
		internalAPIHeader: testInternalToken,
	})
	if unauthenticated.Code != http.StatusUnauthorized {
		t.Fatalf("unauthenticated status = %d body=%s", unauthenticated.Code, unauthenticated.Body.String())
	}

	user := createVerifiedUser(t, db, "user@example.com", "password123")
	userToken := loginUser(t, router, user.Email, "password123")
	forbidden := orgUnitJSON(t, router, http.MethodPost, "/api/organization-units", `{"name":"Company"}`, authCookieHeaders(userToken))
	if forbidden.Code != http.StatusForbidden {
		t.Fatalf("user create status = %d body=%s", forbidden.Code, forbidden.Body.String())
	}

	admin := createAdminUser(t, db, "admin@example.com", "password123")
	adminToken := loginUser(t, router, admin.Email, "password123")
	created := orgUnitJSON(t, router, http.MethodPost, "/api/organization-units", `{"name":"Company","code":" root "}`, authCookieHeaders(adminToken))
	if created.Code != http.StatusCreated {
		t.Fatalf("admin create status = %d body=%s", created.Code, created.Body.String())
	}

	list := orgUnitJSON(t, router, http.MethodGet, "/api/organization-units/tree", "", authCookieHeaders(userToken))
	if list.Code != http.StatusOK {
		t.Fatalf("user list status = %d body=%s", list.Code, list.Body.String())
	}
	var body struct {
		Units []struct {
			Name string  `json:"name"`
			Code *string `json:"code"`
		} `json:"units"`
	}
	decodeJSON(t, list, &body)
	if len(body.Units) != 1 || body.Units[0].Name != "Company" || body.Units[0].Code == nil || *body.Units[0].Code != "ROOT" {
		t.Fatalf("tree body = %#v", body)
	}
}

func TestOrganizationUnitMoveRejectsCyclesAndArchiveCascades(t *testing.T) {
	router, db := newAuthTestRouter(t)
	admin := createAdminUser(t, db, "admin@example.com", "password123")
	token := loginUser(t, router, admin.Email, "password123")

	rootID := createUnitViaAPI(t, router, token, `{"name":"Company"}`)
	childID := createUnitViaAPI(t, router, token, `{"name":"Finance","parentId":"`+rootID+`"}`)
	grandchildID := createUnitViaAPI(t, router, token, `{"name":"Accounts Payable","parentId":"`+childID+`"}`)

	cycle := orgUnitJSON(t, router, http.MethodPatch, "/api/organization-units/"+rootID+"/parent", `{"parentId":"`+grandchildID+`"}`, authCookieHeaders(token))
	if cycle.Code != http.StatusConflict {
		t.Fatalf("cycle status = %d body=%s", cycle.Code, cycle.Body.String())
	}

	archive := orgUnitJSON(t, router, http.MethodPost, "/api/organization-units/"+childID+"/archive", `{}`, authCookieHeaders(token))
	if archive.Code != http.StatusOK {
		t.Fatalf("archive status = %d body=%s", archive.Code, archive.Body.String())
	}

	var archivedCount int64
	if err := db.Model(&orgunits.Unit{}).Where("id IN ?", []string{childID, grandchildID}).Where("archived_at IS NOT NULL").Count(&archivedCount).Error; err != nil {
		t.Fatalf("count archived units: %v", err)
	}
	if archivedCount != 2 {
		t.Fatalf("archived count = %d, want 2", archivedCount)
	}
}

func createAdminUser(t *testing.T, db *gorm.DB, email string, password string) auth.User {
	t.Helper()
	user := createVerifiedUser(t, db, email, password)
	if err := db.Model(&auth.User{}).Where("id = ?", user.ID).Update("role", auth.UserRoleAdmin).Error; err != nil {
		t.Fatalf("promote admin: %v", err)
	}
	user.Role = auth.UserRoleAdmin
	return user
}

func authCookieHeaders(token string) map[string]string {
	return map[string]string{
		internalAPIHeader: testInternalToken,
		"Cookie":          authSessionCookieName + "=" + token,
	}
}

func createUnitViaAPI(t *testing.T, router http.Handler, token string, body string) string {
	t.Helper()
	response := orgUnitJSON(t, router, http.MethodPost, "/api/organization-units", body, authCookieHeaders(token))
	if response.Code != http.StatusCreated {
		t.Fatalf("create unit status = %d body=%s", response.Code, response.Body.String())
	}
	var out struct {
		ID string `json:"id"`
	}
	decodeJSON(t, response, &out)
	if out.ID == "" {
		t.Fatal("created unit id was empty")
	}
	return out.ID
}
```

Also add `orgUnitJSON`:

```go
func orgUnitJSON(t *testing.T, router http.Handler, method string, path string, body string, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	return authJSON(t, router, method, path, body, headers)
}
```

- [ ] **Step 3: Run API tests to verify failure**

Run:

```sh
cd go-server
rtk go test ./internal/api -run OrganizationUnit
```

Expected: FAIL because Organization Unit routes and handlers are not wired.

## Task 3: Backend API Implementation

**Files:**
- Create: `go-server/internal/api/organization_units.go`
- Modify: `go-server/internal/api/router.go`
- Modify: `go-server/internal/api/swagger_doc.go`

- [ ] **Step 1: Wire routes in the router**

Modify `go-server/internal/api/router.go` after auth routes:

```go
	orgUnits := newOrganizationUnitHandler(options, auth)
	orgUnitAPI := router.Group("/api/organization-units")
	orgUnitAPI.Use(auth.requireTrustedInternalRequest())
	orgUnitAPI.GET("/tree", orgUnits.listTree)
	orgUnitAPI.GET("/archived", orgUnits.listArchived)
	orgUnitAPI.POST("", orgUnits.create)
	orgUnitAPI.PATCH("/:id", orgUnits.update)
	orgUnitAPI.PATCH("/:id/parent", orgUnits.move)
	orgUnitAPI.POST("/:id/archive", orgUnits.archive)
```

- [ ] **Step 2: Implement request and response types**

Create `go-server/internal/api/organization_units.go` with these top-level types:

```go
package api

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"ai.ro/syncra/dms/internal/auth"
	"ai.ro/syncra/dms/internal/orgunits"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const maxOrganizationUnitRequestBytes int64 = 1 << 20

type organizationUnitHandler struct {
	db   *gorm.DB
	auth *authHandler
}

type organizationUnitRequest struct {
	ParentID    *string `json:"parentId"`
	Name        string  `json:"name"`
	Code        string  `json:"code"`
	Description string  `json:"description"`
}

type moveOrganizationUnitRequest struct {
	ParentID *string `json:"parentId"`
}

type organizationUnitResponse struct {
	ID          string                      `json:"id"`
	ParentID    *string                     `json:"parentId,omitempty"`
	Name        string                      `json:"name"`
	Code        *string                     `json:"code,omitempty"`
	Description *string                     `json:"description,omitempty"`
	ArchivedAt  *string                     `json:"archivedAt,omitempty"`
	CreatedAt   string                      `json:"createdAt"`
	UpdatedAt   string                      `json:"updatedAt"`
	Children    []organizationUnitResponse `json:"children,omitempty"`
}

type organizationUnitListResponse struct {
	Units []organizationUnitResponse `json:"units"`
}

func newOrganizationUnitHandler(options RouterOptions, auth *authHandler) *organizationUnitHandler {
	return &organizationUnitHandler{db: options.DB, auth: auth}
}
```

- [ ] **Step 3: Implement auth helpers**

Add these helpers in the same file:

```go
func (h *organizationUnitHandler) requireUser(c *gin.Context) (auth.User, bool) {
	if h.db == nil {
		writeError(c, http.StatusServiceUnavailable, "database is not configured")
		return auth.User{}, false
	}
	if h.auth == nil || !h.auth.authConfigured(c) {
		return auth.User{}, false
	}
	session, ok, err := h.auth.loadAuthenticatedSession(c)
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

func (h *organizationUnitHandler) requireAdmin(c *gin.Context) (auth.User, bool) {
	user, ok := h.requireUser(c)
	if !ok {
		return auth.User{}, false
	}
	if user.Role != auth.UserRoleAdmin {
		writeError(c, http.StatusForbidden, "admin role required")
		return auth.User{}, false
	}
	return user, true
}
```

- [ ] **Step 4: Implement list and response mapping**

Add:

```go
func (h *organizationUnitHandler) listTree(c *gin.Context) {
	if _, ok := h.requireUser(c); !ok {
		return
	}
	var units []orgunits.Unit
	if err := h.db.WithContext(c.Request.Context()).
		Where("archived_at IS NULL").
		Order("name asc, id asc").
		Find(&units).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to list organization units")
		return
	}
	c.JSON(http.StatusOK, organizationUnitListResponse{Units: organizationUnitTreeResponse(orgunits.BuildTree(units))})
}

func (h *organizationUnitHandler) listArchived(c *gin.Context) {
	if _, ok := h.requireAdmin(c); !ok {
		return
	}
	var units []orgunits.Unit
	if err := h.db.WithContext(c.Request.Context()).
		Where("archived_at IS NOT NULL").
		Order("archived_at desc, name asc, id asc").
		Find(&units).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to list archived organization units")
		return
	}
	out := make([]organizationUnitResponse, 0, len(units))
	for _, unit := range units {
		out = append(out, organizationUnitResponseFromUnit(unit))
	}
	c.JSON(http.StatusOK, organizationUnitListResponse{Units: out})
}
```

Implement mapping helpers:

```go
func organizationUnitTreeResponse(nodes []orgunits.TreeNode) []organizationUnitResponse {
	out := make([]organizationUnitResponse, 0, len(nodes))
	for _, node := range nodes {
		out = append(out, organizationUnitResponse{
			ID:          node.ID,
			ParentID:    node.ParentID,
			Name:        node.Name,
			Code:        node.Code,
			Description: node.Description,
			CreatedAt:   node.CreatedAt,
			UpdatedAt:   node.UpdatedAt,
			Children:    organizationUnitTreeResponse(node.Children),
		})
	}
	return out
}

func organizationUnitResponseFromUnit(unit orgunits.Unit) organizationUnitResponse {
	var archivedAt *string
	if unit.ArchivedAt != nil {
		value := unit.ArchivedAt.UTC().Format(time.RFC3339Nano)
		archivedAt = &value
	}
	return organizationUnitResponse{
		ID:          unit.ID,
		ParentID:    unit.ParentID,
		Name:        unit.Name,
		Code:        unit.Code,
		Description: unit.Description,
		ArchivedAt:  archivedAt,
		CreatedAt:   unit.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:   unit.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
}
```

- [ ] **Step 5: Implement create/update/move/archive handlers**

Implement these handlers and helper rules:

```go
func (h *organizationUnitHandler) create(c *gin.Context) {
	if _, ok := h.requireAdmin(c); !ok {
		return
	}
	req, ok := bindOrganizationUnitRequest(c)
	if !ok {
		return
	}
	name, code, description, ok := normalizeOrganizationUnitInput(c, req.Name, req.Code, req.Description)
	if !ok {
		return
	}
	parentID, ok := h.validateOptionalActiveParent(c, req.ParentID)
	if !ok {
		return
	}
	if ok := h.ensureActiveCodeAvailable(c, code, ""); !ok {
		return
	}
	now := time.Now().UTC()
	unit := orgunits.Unit{
		ParentID:    parentID,
		Name:        name,
		Code:        code,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := h.db.WithContext(c.Request.Context()).Create(&unit).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to create organization unit")
		return
	}
	c.JSON(http.StatusCreated, organizationUnitResponseFromUnit(unit))
}

func (h *organizationUnitHandler) update(c *gin.Context) {
	if _, ok := h.requireAdmin(c); !ok {
		return
	}
	id, ok := parseOrganizationUnitID(c, c.Param("id"))
	if !ok {
		return
	}
	req, ok := bindOrganizationUnitRequest(c)
	if !ok {
		return
	}
	name, code, description, ok := normalizeOrganizationUnitInput(c, req.Name, req.Code, req.Description)
	if !ok {
		return
	}
	if ok := h.ensureActiveCodeAvailable(c, code, id); !ok {
		return
	}
	var unit orgunits.Unit
	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&orgunits.Unit{}).Where("id = ? AND archived_at IS NULL", id).Updates(map[string]any{
			"name":        name,
			"code":        code,
			"description": description,
			"updated_at":   time.Now().UTC(),
		})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return tx.First(&unit, "id = ?", id).Error
	})
	if err != nil {
		writeOrganizationUnitMutationError(c, err, "failed to update organization unit")
		return
	}
	c.JSON(http.StatusOK, organizationUnitResponseFromUnit(unit))
}
```

Use the same style for `move` and `archive`:

```go
func (h *organizationUnitHandler) move(c *gin.Context) {
	if _, ok := h.requireAdmin(c); !ok {
		return
	}
	id, ok := parseOrganizationUnitID(c, c.Param("id"))
	if !ok {
		return
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxOrganizationUnitRequestBytes)
	var req moveOrganizationUnitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	parentID, ok := h.validateOptionalActiveParent(c, req.ParentID)
	if !ok {
		return
	}
	if parentID != nil && *parentID == id {
		writeError(c, http.StatusConflict, "organization unit cannot be moved under itself")
		return
	}
	if parentID != nil && h.parentWouldCreateCycle(c, id, *parentID) {
		writeError(c, http.StatusConflict, "organization unit cannot be moved under its descendant")
		return
	}
	var unit orgunits.Unit
	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&orgunits.Unit{}).Where("id = ? AND archived_at IS NULL", id).Updates(map[string]any{
			"parent_id":  parentID,
			"updated_at": time.Now().UTC(),
		})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return tx.First(&unit, "id = ?", id).Error
	})
	if err != nil {
		writeOrganizationUnitMutationError(c, err, "failed to move organization unit")
		return
	}
	c.JSON(http.StatusOK, organizationUnitResponseFromUnit(unit))
}

func (h *organizationUnitHandler) archive(c *gin.Context) {
	if _, ok := h.requireAdmin(c); !ok {
		return
	}
	id, ok := parseOrganizationUnitID(c, c.Param("id"))
	if !ok {
		return
	}
	now := time.Now().UTC()
	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var units []orgunits.Unit
		if err := tx.Where("archived_at IS NULL").Find(&units).Error; err != nil {
			return err
		}
		found := false
		for _, unit := range units {
			if unit.ID == id {
				found = true
				break
			}
		}
		if !found {
			return gorm.ErrRecordNotFound
		}
		ids := []string{id}
		for descendantID := range orgunits.DescendantIDs(id, units) {
			ids = append(ids, descendantID)
		}
		return tx.Model(&orgunits.Unit{}).Where("id IN ?", ids).Updates(map[string]any{
			"archived_at": now,
			"updated_at":  now,
		}).Error
	})
	if err != nil {
		writeOrganizationUnitMutationError(c, err, "failed to archive organization unit")
		return
	}
	c.JSON(http.StatusOK, okResponse{OK: true})
}
```

- [ ] **Step 6: Implement handler helper functions**

Add helper functions for ids, binding, parent validation, duplicate code, and cycle detection:

```go
func bindOrganizationUnitRequest(c *gin.Context) (organizationUnitRequest, bool) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxOrganizationUnitRequestBytes)
	var req organizationUnitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return organizationUnitRequest{}, false
	}
	return req, true
}

func normalizeOrganizationUnitInput(c *gin.Context, rawName string, rawCode string, rawDescription string) (string, *string, *string, bool) {
	name, err := orgunits.NormalizeName(rawName)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return "", nil, nil, false
	}
	code, err := orgunits.NormalizeCode(rawCode)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return "", nil, nil, false
	}
	return name, code, orgunits.NormalizeDescription(rawDescription), true
}

func parseOrganizationUnitID(c *gin.Context, raw string) (string, bool) {
	id, err := uuid.Parse(strings.TrimSpace(raw))
	if err != nil || id == uuid.Nil {
		writeError(c, http.StatusBadRequest, "invalid organization unit id")
		return "", false
	}
	return id.String(), true
}

func (h *organizationUnitHandler) validateOptionalActiveParent(c *gin.Context, raw *string) (*string, bool) {
	if raw == nil || strings.TrimSpace(*raw) == "" {
		return nil, true
	}
	parentID, ok := parseOrganizationUnitID(c, *raw)
	if !ok {
		return nil, false
	}
	var count int64
	if err := h.db.WithContext(c.Request.Context()).Model(&orgunits.Unit{}).Where("id = ? AND archived_at IS NULL", parentID).Count(&count).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to validate parent organization unit")
		return nil, false
	}
	if count == 0 {
		writeError(c, http.StatusNotFound, "parent organization unit not found")
		return nil, false
	}
	return &parentID, true
}

func (h *organizationUnitHandler) ensureActiveCodeAvailable(c *gin.Context, code *string, currentID string) bool {
	if code == nil {
		return true
	}
	query := h.db.WithContext(c.Request.Context()).Model(&orgunits.Unit{}).Where("code = ? AND archived_at IS NULL", *code)
	if currentID != "" {
		query = query.Where("id <> ?", currentID)
	}
	var count int64
	if err := query.Count(&count).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to validate organization unit code")
		return false
	}
	if count > 0 {
		writeError(c, http.StatusConflict, "organization unit code already exists")
		return false
	}
	return true
}

func (h *organizationUnitHandler) parentWouldCreateCycle(c *gin.Context, unitID string, parentID string) bool {
	var units []orgunits.Unit
	if err := h.db.WithContext(c.Request.Context()).Where("archived_at IS NULL").Find(&units).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to validate organization unit move")
		return true
	}
	descendants := orgunits.DescendantIDs(unitID, units)
	return descendants[parentID]
}

func writeOrganizationUnitMutationError(c *gin.Context, err error, fallback string) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		writeError(c, http.StatusNotFound, "organization unit not found")
		return
	}
	writeError(c, http.StatusInternalServerError, fallback)
}
```

- [ ] **Step 7: Add Swagger operations**

Modify `go-server/internal/api/swagger_doc.go` and add these operation comments inside `swaggerOperations()`:

```go
	// swagger:operation GET /api/organization-units/tree organizationUnits listOrganizationUnits
	//
	// List active organization units as a hierarchy.
	//
	// Trusted SvelteKit server endpoint. Requires an authenticated session.
	//
	// ---
	// responses:
	//   "200":
	//     description: Active organization unit hierarchy.
	//   "401":
	//     description: Authenticated session or trusted internal request required.

	// swagger:operation GET /api/organization-units/archived organizationUnits listArchivedOrganizationUnits
	//
	// List archived organization units.
	//
	// Trusted SvelteKit server endpoint. Requires an admin session.
	//
	// ---
	// responses:
	//   "200":
	//     description: Archived organization units.
	//   "401":
	//     description: Authenticated session or trusted internal request required.
	//   "403":
	//     description: Admin role required.

	// swagger:operation POST /api/organization-units organizationUnits createOrganizationUnit
	//
	// Create a root or child organization unit.
	//
	// Trusted SvelteKit server endpoint. Requires an admin session.
	//
	// ---
	// responses:
	//   "201":
	//     description: Organization unit created.
	//   "400":
	//     description: Invalid organization unit request.
	//   "401":
	//     description: Authenticated session or trusted internal request required.
	//   "403":
	//     description: Admin role required.
	//   "404":
	//     description: Parent organization unit not found.
	//   "409":
	//     description: Active organization unit code already exists.

	// swagger:operation PATCH /api/organization-units/{id} organizationUnits updateOrganizationUnit
	//
	// Update an active organization unit's details.
	//
	// Trusted SvelteKit server endpoint. Requires an admin session.
	//
	// ---
	// responses:
	//   "200":
	//     description: Organization unit updated.
	//   "400":
	//     description: Invalid organization unit request.
	//   "401":
	//     description: Authenticated session or trusted internal request required.
	//   "403":
	//     description: Admin role required.
	//   "404":
	//     description: Organization unit not found.
	//   "409":
	//     description: Active organization unit code already exists.

	// swagger:operation PATCH /api/organization-units/{id}/parent organizationUnits moveOrganizationUnit
	//
	// Move an active organization unit to another parent or root.
	//
	// Trusted SvelteKit server endpoint. Requires an admin session.
	//
	// ---
	// responses:
	//   "200":
	//     description: Organization unit moved.
	//   "400":
	//     description: Invalid move request.
	//   "401":
	//     description: Authenticated session or trusted internal request required.
	//   "403":
	//     description: Admin role required.
	//   "404":
	//     description: Organization unit or parent not found.
	//   "409":
	//     description: Move would create a cycle.

	// swagger:operation POST /api/organization-units/{id}/archive organizationUnits archiveOrganizationUnit
	//
	// Archive an active organization unit and its descendants.
	//
	// Trusted SvelteKit server endpoint. Requires an admin session.
	//
	// ---
	// responses:
	//   "200":
	//     description: Organization unit subtree archived.
	//   "400":
	//     description: Invalid organization unit id.
	//   "401":
	//     description: Authenticated session or trusted internal request required.
	//   "403":
	//     description: Admin role required.
	//   "404":
	//     description: Organization unit not found.
```

- [ ] **Step 8: Run backend API tests**

Run:

```sh
cd go-server
rtk go test ./internal/api -run OrganizationUnit
rtk go test ./...
```

Expected: PASS.

- [ ] **Step 9: Generate Swagger**

Run:

```sh
cd go-server
rtk go run ./cmd/syncra swagger
```

Expected: `go-server/docs/swagger.json` and `go-server/docs/docs.go` are updated. If the command reports that the `swagger` binary is missing, keep the source annotations and report the missing tool in the task handoff.

- [ ] **Step 10: Commit backend API**

Run:

```sh
rtk git add go-server/internal/api go-server/docs go-server/internal/api/swagger_doc.go
rtk git commit -m "Add organization units API"
```

## Task 4: Frontend Server API Client

**Files:**
- Create: `frontend/src/lib/server/organization-units.ts`
- Create: `frontend/src/lib/server/organization-units.test.ts`

- [ ] **Step 1: Write failing server-client tests**

Create `frontend/src/lib/server/organization-units.test.ts`:

```ts
import { afterEach, describe, expect, it, vi } from 'vitest';
import {
	OrganizationUnitApiError,
	createOrganizationUnit,
	getOrganizationUnitTree
} from './organization-units';

describe('organization unit server client', () => {
	afterEach(() => vi.unstubAllEnvs());

	it('loads tree with internal token and forwarded cookie', async () => {
		vi.stubEnv('SYNCRA_API_BASE_URL', 'http://api.test');
		vi.stubEnv('SYNCRA_INTERNAL_API_TOKEN', 'internal-token');
		const fetch = vi.fn(async () =>
			jsonResponse({
				units: [{ id: 'unit-id', name: 'Company', createdAt: '2026-06-30T12:00:00Z', updatedAt: '2026-06-30T12:00:00Z', children: [] }]
			})
		);

		const tree = await getOrganizationUnitTree(fetch, 'auth.session_token=session-token');

		expect(tree.units[0].name).toBe('Company');
		const call = fetch.mock.calls[0] as unknown as [string, RequestInit];
		expect(call[0]).toBe('http://api.test/api/organization-units/tree');
		const headers = call[1].headers as Headers;
		expect(headers.get('X-Syncra-Internal-Token')).toBe('internal-token');
		expect(headers.get('cookie')).toBe('auth.session_token=session-token');
	});

	it('posts create payloads', async () => {
		vi.stubEnv('SYNCRA_API_BASE_URL', 'http://api.test');
		vi.stubEnv('SYNCRA_INTERNAL_API_TOKEN', 'internal-token');
		const fetch = vi.fn(async () =>
			jsonResponse({ id: 'unit-id', name: 'Finance', createdAt: '2026-06-30T12:00:00Z', updatedAt: '2026-06-30T12:00:00Z' }, 201)
		);

		await createOrganizationUnit(fetch, 'auth.session_token=session-token', {
			name: 'Finance',
			code: 'FIN',
			description: 'Finance team',
			parentId: null
		});

		const call = fetch.mock.calls[0] as unknown as [string, RequestInit];
		expect(call[0]).toBe('http://api.test/api/organization-units');
		expect(call[1].method).toBe('POST');
		expect(call[1].body).toBe(JSON.stringify({ name: 'Finance', code: 'FIN', description: 'Finance team', parentId: null }));
	});

	it('throws public errors', async () => {
		vi.stubEnv('SYNCRA_API_BASE_URL', 'http://api.test');
		vi.stubEnv('SYNCRA_INTERNAL_API_TOKEN', 'internal-token');
		const fetch = vi.fn(async () => jsonResponse({ error: 'admin role required' }, 403));

		await expect(getOrganizationUnitTree(fetch, 'auth.session_token=session-token')).rejects.toMatchObject(
			new OrganizationUnitApiError(403, 'admin role required')
		);
	});
});

function jsonResponse(body: unknown, status = 200) {
	return new Response(JSON.stringify(body), {
		status,
		headers: { 'content-type': 'application/json' }
	});
}
```

- [ ] **Step 2: Run tests to verify failure**

Run:

```sh
cd frontend
rtk pnpm test -- src/lib/server/organization-units.test.ts
```

Expected: FAIL because `organization-units.ts` does not exist.

- [ ] **Step 3: Implement the server API client**

Create `frontend/src/lib/server/organization-units.ts`:

```ts
import { apiBaseUrl, internalAPIHeaders } from './internal-api';
import { publicErrorMessage, publicErrorStatus } from './public-errors';

type ServerFetch = typeof fetch;

export type OrganizationUnit = {
	id: string;
	parentId?: string | null;
	name: string;
	code?: string | null;
	description?: string | null;
	archivedAt?: string | null;
	createdAt: string;
	updatedAt: string;
	children: OrganizationUnit[];
};

export type OrganizationUnitListResponse = {
	units: OrganizationUnit[];
};

export type OrganizationUnitInput = {
	parentId?: string | null;
	name: string;
	code: string;
	description: string;
};

export class OrganizationUnitApiError extends Error {
	status: number;

	constructor(status: number, message: string) {
		super(message);
		this.name = 'OrganizationUnitApiError';
		this.status = status;
	}
}

export function isOrganizationUnitApiError(error: unknown): error is OrganizationUnitApiError {
	return error instanceof OrganizationUnitApiError;
}

export function getOrganizationUnitTree(fetchFn: ServerFetch, cookieHeader: string | null) {
	return organizationUnitRequest<OrganizationUnitListResponse>(fetchFn, '/api/organization-units/tree', {
		cookieHeader
	});
}

export function createOrganizationUnit(fetchFn: ServerFetch, cookieHeader: string | null, input: OrganizationUnitInput) {
	return organizationUnitRequest<OrganizationUnit>(fetchFn, '/api/organization-units', {
		method: 'POST',
		cookieHeader,
		body: input
	});
}

export function updateOrganizationUnit(fetchFn: ServerFetch, cookieHeader: string | null, id: string, input: OrganizationUnitInput) {
	return organizationUnitRequest<OrganizationUnit>(fetchFn, `/api/organization-units/${encodeURIComponent(id)}`, {
		method: 'PATCH',
		cookieHeader,
		body: input
	});
}

export function moveOrganizationUnit(fetchFn: ServerFetch, cookieHeader: string | null, id: string, parentId: string | null) {
	return organizationUnitRequest<OrganizationUnit>(fetchFn, `/api/organization-units/${encodeURIComponent(id)}/parent`, {
		method: 'PATCH',
		cookieHeader,
		body: { parentId }
	});
}

export function archiveOrganizationUnit(fetchFn: ServerFetch, cookieHeader: string | null, id: string) {
	return organizationUnitRequest<{ ok: boolean }>(fetchFn, `/api/organization-units/${encodeURIComponent(id)}/archive`, {
		method: 'POST',
		cookieHeader,
		body: {}
	});
}

async function organizationUnitRequest<T>(
	fetchFn: ServerFetch,
	path: string,
	options: { method?: 'GET' | 'POST' | 'PATCH'; cookieHeader?: string | null; body?: unknown } = {}
) {
	const headers = internalAPIHeaders();
	if (!headers) throw new OrganizationUnitApiError(500, 'Organization Unit service is not configured');
	if (options.cookieHeader) headers.set('cookie', options.cookieHeader);
	if (options.body !== undefined) headers.set('content-type', 'application/json');

	let response: Response;
	try {
		response = await fetchFn(`${apiBaseUrl()}${path}`, {
			method: options.method ?? 'GET',
			headers,
			body: options.body === undefined ? undefined : JSON.stringify(options.body)
		});
	} catch {
		throw new OrganizationUnitApiError(503, 'Organization Unit service unavailable');
	}

	const json = await readResponseJSON(response);
	if (!response.ok) {
		const message =
			json && typeof json === 'object' && 'error' in json && typeof json.error === 'string'
				? json.error
				: 'Organization Unit request failed';
		throw new OrganizationUnitApiError(publicErrorStatus(response.status), publicErrorMessage(response.status, message, 'Organization Unit request failed'));
	}
	return json as T;
}

async function readResponseJSON(response: Response) {
	const text = await response.text();
	if (!text.trim()) return null;
	try {
		return JSON.parse(text) as unknown;
	} catch {
		return null;
	}
}
```

- [ ] **Step 4: Run client tests**

Run:

```sh
cd frontend
rtk pnpm test -- src/lib/server/organization-units.test.ts
```

Expected: PASS.

- [ ] **Step 5: Commit frontend server client**

Run:

```sh
rtk git add frontend/src/lib/server/organization-units.ts frontend/src/lib/server/organization-units.test.ts
rtk git commit -m "Add organization units server client"
```

## Task 5: Frontend Tree Utilities

**Files:**
- Create: `frontend/src/routes/app/organization-units/tree.ts`
- Create: `frontend/src/routes/app/organization-units/tree.test.ts`

- [ ] **Step 1: Write failing tree utility tests**

Create `frontend/src/routes/app/organization-units/tree.test.ts`:

```ts
import { describe, expect, it } from 'vitest';
import { descendantIds, findUnitById, flattenOrganizationUnitTree, selectableParentOptions } from './tree';
import type { OrganizationUnit } from '$lib/server/organization-units';

const tree: OrganizationUnit[] = [
	{
		id: 'root',
		name: 'Company',
		createdAt: '2026-06-30T12:00:00Z',
		updatedAt: '2026-06-30T12:00:00Z',
		children: [
			{
				id: 'finance',
				parentId: 'root',
				name: 'Finance',
				createdAt: '2026-06-30T12:00:00Z',
				updatedAt: '2026-06-30T12:00:00Z',
				children: [
					{
						id: 'ap',
						parentId: 'finance',
						name: 'Accounts Payable',
						createdAt: '2026-06-30T12:00:00Z',
						updatedAt: '2026-06-30T12:00:00Z',
						children: []
					}
				]
			}
		]
	}
];

describe('organization unit tree utilities', () => {
	it('flattens units with depth and path labels', () => {
		expect(flattenOrganizationUnitTree(tree).map((item) => [item.id, item.depth, item.path])).toEqual([
			['root', 0, 'Company'],
			['finance', 1, 'Company / Finance'],
			['ap', 2, 'Company / Finance / Accounts Payable']
		]);
	});

	it('finds units by id', () => {
		expect(findUnitById(tree, 'finance')?.name).toBe('Finance');
		expect(findUnitById(tree, 'missing')).toBeNull();
	});

	it('excludes selected unit and descendants from parent options', () => {
		expect(descendantIds(tree, 'finance')).toEqual(new Set(['ap']));
		expect(selectableParentOptions(tree, 'finance').map((option) => option.id)).toEqual(['root']);
	});
});
```

- [ ] **Step 2: Run tests to verify failure**

Run:

```sh
cd frontend
rtk pnpm test -- src/routes/app/organization-units/tree.test.ts
```

Expected: FAIL because `tree.ts` does not exist.

- [ ] **Step 3: Implement tree utilities**

Create `frontend/src/routes/app/organization-units/tree.ts`:

```ts
import type { OrganizationUnit } from '$lib/server/organization-units';

export type FlatOrganizationUnit = OrganizationUnit & {
	depth: number;
	path: string;
};

export type ParentOption = {
	id: string;
	label: string;
};

export function flattenOrganizationUnitTree(units: OrganizationUnit[], depth = 0, ancestors: string[] = []): FlatOrganizationUnit[] {
	return units.flatMap((unit) => {
		const pathParts = [...ancestors, unit.name];
		return [
			{ ...unit, depth, path: pathParts.join(' / ') },
			...flattenOrganizationUnitTree(unit.children ?? [], depth + 1, pathParts)
		];
	});
}

export function findUnitById(units: OrganizationUnit[], id: string | null): OrganizationUnit | null {
	if (!id) return null;
	for (const unit of units) {
		if (unit.id === id) return unit;
		const child = findUnitById(unit.children ?? [], id);
		if (child) return child;
	}
	return null;
}

export function descendantIds(units: OrganizationUnit[], id: string): Set<string> {
	const selected = findUnitById(units, id);
	const ids = new Set<string>();
	function visit(unit: OrganizationUnit) {
		for (const child of unit.children ?? []) {
			ids.add(child.id);
			visit(child);
		}
	}
	if (selected) visit(selected);
	return ids;
}

export function selectableParentOptions(units: OrganizationUnit[], selectedId: string | null): ParentOption[] {
	const excluded = selectedId ? descendantIds(units, selectedId) : new Set<string>();
	if (selectedId) excluded.add(selectedId);
	return flattenOrganizationUnitTree(units)
		.filter((unit) => !excluded.has(unit.id))
		.map((unit) => ({ id: unit.id, label: unit.path }));
}
```

- [ ] **Step 4: Run utility tests**

Run:

```sh
cd frontend
rtk pnpm test -- src/routes/app/organization-units/tree.test.ts
```

Expected: PASS.

- [ ] **Step 5: Commit tree utilities**

Run:

```sh
rtk git add frontend/src/routes/app/organization-units/tree.ts frontend/src/routes/app/organization-units/tree.test.ts
rtk git commit -m "Add organization unit tree utilities"
```

## Task 6: SvelteKit Route Load And Server Actions

**Files:**
- Create: `frontend/src/routes/app/+layout.server.ts`
- Create: `frontend/src/routes/app/organization-units/+page.server.ts`
- Create: `frontend/src/routes/app/organization-units/page.server.test.ts`

- [ ] **Step 1: Add app layout server data**

Create `frontend/src/routes/app/+layout.server.ts`:

```ts
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = ({ locals }) => ({
	user: locals.user,
	session: locals.session
});
```

- [ ] **Step 2: Write failing page server tests**

Create `frontend/src/routes/app/organization-units/page.server.test.ts`:

```ts
import { describe, expect, it, vi } from 'vitest';

const {
	getOrganizationUnitTreeMock,
	createOrganizationUnitMock,
	updateOrganizationUnitMock,
	moveOrganizationUnitMock,
	archiveOrganizationUnitMock
} = vi.hoisted(() => ({
	getOrganizationUnitTreeMock: vi.fn(),
	createOrganizationUnitMock: vi.fn(),
	updateOrganizationUnitMock: vi.fn(),
	moveOrganizationUnitMock: vi.fn(),
	archiveOrganizationUnitMock: vi.fn()
}));

vi.mock('$lib/server/organization-units', () => ({
	OrganizationUnitApiError: class OrganizationUnitApiError extends Error {
		status: number;
		constructor(status: number, message: string) {
			super(message);
			this.status = status;
		}
	},
	archiveOrganizationUnit: archiveOrganizationUnitMock,
	createOrganizationUnit: createOrganizationUnitMock,
	getOrganizationUnitTree: getOrganizationUnitTreeMock,
	isOrganizationUnitApiError: (error: unknown) => error instanceof Error && 'status' in error,
	moveOrganizationUnit: moveOrganizationUnitMock,
	updateOrganizationUnit: updateOrganizationUnitMock
}));

import { actions, load } from './+page.server';

describe('organization units page server', () => {
	it('loads tree and role data', async () => {
		getOrganizationUnitTreeMock.mockResolvedValue({ units: [] });
		const event = pageEvent({ role: 'admin' });

		const data = await load(event as never);

		expect(getOrganizationUnitTreeMock).toHaveBeenCalledWith(event.fetch, 'auth.session_token=token');
		expect(data.canManageOrganizationUnits).toBe(true);
		expect(data.units).toEqual([]);
		expect(data.loadError).toBeNull();
	});

	it('creates units through server action', async () => {
		createOrganizationUnitMock.mockResolvedValue({ id: 'unit-id' });
		const formData = new FormData();
		formData.set('name', 'Finance');
		formData.set('code', 'FIN');
		formData.set('description', 'Finance team');
		formData.set('parentId', '');
		const event = actionEvent(formData);

		const result = await actions.create(event as never);

		expect(result).toEqual({ success: true });
		expect(createOrganizationUnitMock).toHaveBeenCalledWith(event.fetch, 'auth.session_token=token', {
			name: 'Finance',
			code: 'FIN',
			description: 'Finance team',
			parentId: null
		});
	});
});

function pageEvent(user = { role: 'user' }) {
	return {
		fetch: vi.fn(),
		locals: { user },
		request: new Request('http://localhost/app/organization-units', {
			headers: { cookie: 'auth.session_token=token' }
		})
	};
}

function actionEvent(formData: FormData) {
	return {
		fetch: vi.fn(),
		locals: { user: { role: 'admin' } },
		request: new Request('http://localhost/app/organization-units', {
			method: 'POST',
			body: formData,
			headers: { cookie: 'auth.session_token=token' }
		})
	};
}
```

- [ ] **Step 3: Run tests to verify failure**

Run:

```sh
cd frontend
rtk pnpm test -- src/routes/app/organization-units/page.server.test.ts
```

Expected: FAIL because `+page.server.ts` does not exist.

- [ ] **Step 4: Implement page load and actions**

Create `frontend/src/routes/app/organization-units/+page.server.ts`:

```ts
import { fail, type Actions } from '@sveltejs/kit';
import {
	archiveOrganizationUnit,
	createOrganizationUnit,
	getOrganizationUnitTree,
	isOrganizationUnitApiError,
	moveOrganizationUnit,
	updateOrganizationUnit
} from '$lib/server/organization-units';
import { publicErrorMessage } from '$lib/server/public-errors';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch, locals, request }) => {
	const cookieHeader = request.headers.get('cookie');
	try {
		const tree = await getOrganizationUnitTree(fetch, cookieHeader);
		return {
			units: tree.units,
			loadError: null,
			canManageOrganizationUnits: locals.user?.role === 'admin'
		};
	} catch (error) {
		if (isOrganizationUnitApiError(error)) {
			return {
				units: [],
				loadError: publicErrorMessage(error.status, error.message, 'Failed to load organization units'),
				canManageOrganizationUnits: locals.user?.role === 'admin'
			};
		}
		throw error;
	}
};

export const actions: Actions = {
	create: async (event) => {
		const input = await organizationUnitInput(event.request);
		try {
			await createOrganizationUnit(event.fetch, event.request.headers.get('cookie'), input);
			return { success: true };
		} catch (error) {
			return actionFailure(error, 'Failed to create organization unit');
		}
	},
	update: async (event) => {
		const formData = await event.request.formData();
		const id = String(formData.get('id') ?? '');
		const input = inputFromFormData(formData);
		try {
			await updateOrganizationUnit(event.fetch, event.request.headers.get('cookie'), id, input);
			return { success: true };
		} catch (error) {
			return actionFailure(error, 'Failed to update organization unit');
		}
	},
	move: async (event) => {
		const formData = await event.request.formData();
		const id = String(formData.get('id') ?? '');
		const parentId = optionalString(formData.get('parentId'));
		try {
			await moveOrganizationUnit(event.fetch, event.request.headers.get('cookie'), id, parentId);
			return { success: true };
		} catch (error) {
			return actionFailure(error, 'Failed to move organization unit');
		}
	},
	archive: async (event) => {
		const formData = await event.request.formData();
		const id = String(formData.get('id') ?? '');
		try {
			await archiveOrganizationUnit(event.fetch, event.request.headers.get('cookie'), id);
			return { success: true };
		} catch (error) {
			return actionFailure(error, 'Failed to archive organization unit');
		}
	}
};

async function organizationUnitInput(request: Request) {
	return inputFromFormData(await request.formData());
}

function inputFromFormData(formData: FormData) {
	return {
		name: String(formData.get('name') ?? ''),
		code: String(formData.get('code') ?? ''),
		description: String(formData.get('description') ?? ''),
		parentId: optionalString(formData.get('parentId'))
	};
}

function optionalString(value: FormDataEntryValue | null) {
	const text = String(value ?? '').trim();
	return text === '' ? null : text;
}

function actionFailure(error: unknown, fallback: string) {
	if (isOrganizationUnitApiError(error)) {
		return fail(error.status, {
			error: publicErrorMessage(error.status, error.message, fallback)
		});
	}
	throw error;
}
```

- [ ] **Step 5: Run page server tests**

Run:

```sh
cd frontend
rtk pnpm test -- src/routes/app/organization-units/page.server.test.ts
```

Expected: PASS.

- [ ] **Step 6: Commit route server logic**

Run:

```sh
rtk git add frontend/src/routes/app/+layout.server.ts frontend/src/routes/app/organization-units/+page.server.ts frontend/src/routes/app/organization-units/page.server.test.ts
rtk git commit -m "Add organization units route server logic"
```

## Task 7: App Layout And Organization Units UI

**Files:**
- Create: `frontend/src/routes/app/+layout.svelte`
- Modify: `frontend/src/routes/app/+page.svelte`
- Create: `frontend/src/routes/app/organization-units/unit-tree.svelte`
- Create: `frontend/src/routes/app/organization-units/unit-details-panel.svelte`
- Create: `frontend/src/routes/app/organization-units/+page.svelte`

- [ ] **Step 1: Add shared app layout**

Create `frontend/src/routes/app/+layout.svelte` with a compact shell:

```svelte
<script lang="ts">
	import Building2Icon from '@lucide/svelte/icons/building-2';
	import HomeIcon from '@lucide/svelte/icons/home';
	import LogOutIcon from '@lucide/svelte/icons/log-out';
	import type { Snippet } from 'svelte';
	import { Button } from '$lib/components/ui/button';
	import type { LayoutProps } from './$types';

	let { data, children }: LayoutProps & { children: Snippet } = $props();
</script>

<main class="min-h-screen bg-background text-foreground">
	<header class="border-b bg-card">
		<div class="mx-auto flex max-w-6xl flex-col gap-3 px-4 py-4 sm:flex-row sm:items-center sm:justify-between">
			<div class="min-w-0">
				<p class="text-sm text-muted-foreground">Syncra DMS</p>
				<h1 class="truncate text-xl font-semibold">App</h1>
			</div>
			<div class="flex flex-wrap items-center gap-2">
				<Button href="/app" variant="ghost" size="sm" class="gap-2">
					<HomeIcon class="size-4" />
					Dashboard
				</Button>
				<Button href="/app/organization-units" variant="ghost" size="sm" class="gap-2">
					<Building2Icon class="size-4" />
					Organization Units
				</Button>
				<form method="POST" action="/logout">
					<Button type="submit" variant="outline" size="sm" class="gap-2">
						<LogOutIcon class="size-4" />
						Logout
					</Button>
				</form>
			</div>
		</div>
	</header>
	{@render children()}
</main>
```

- [ ] **Step 2: Slim the existing app dashboard page**

Modify `frontend/src/routes/app/+page.svelte` so it no longer renders its own outer `<main>` or duplicate logout header. Keep the account card and add an Organization Units link:

```svelte
<script lang="ts">
	import Building2Icon from '@lucide/svelte/icons/building-2';
	import UserIcon from '@lucide/svelte/icons/user';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

<section class="mx-auto grid max-w-5xl gap-4 px-4 py-6 md:grid-cols-2">
	<Card.Root>
		<Card.Header>
			<Card.Title class="flex items-center gap-2"><UserIcon class="size-4 text-primary" />Account</Card.Title>
			<Card.Description>Authenticated session</Card.Description>
		</Card.Header>
		<Card.Content class="grid gap-2 text-sm">
			<p><span class="font-medium">Name:</span> {data.user?.name}</p>
			<p><span class="font-medium">Email:</span> {data.user?.email}</p>
			<p><span class="font-medium">Session expires:</span> {data.session?.expiresAt}</p>
		</Card.Content>
	</Card.Root>

	<Card.Root>
		<Card.Header>
			<Card.Title class="flex items-center gap-2"><Building2Icon class="size-4 text-primary" />Organization Units</Card.Title>
			<Card.Description>Company structure</Card.Description>
		</Card.Header>
		<Card.Content class="text-sm text-muted-foreground">
			Browse the shared organization hierarchy.
		</Card.Content>
		<Card.Footer>
			<Button href="/app/organization-units" variant="outline" size="sm">Open Organization Units</Button>
		</Card.Footer>
	</Card.Root>
</section>
```

- [ ] **Step 3: Add the tree component**

Create `frontend/src/routes/app/organization-units/unit-tree.svelte`:

```svelte
<script lang="ts">
	import ChevronRightIcon from '@lucide/svelte/icons/chevron-right';
	import type { FlatOrganizationUnit } from './tree';

	let {
		units,
		selectedId,
		onSelect
	}: {
		units: FlatOrganizationUnit[];
		selectedId: string | null;
		onSelect: (id: string) => void;
	} = $props();
</script>

<div class="grid gap-1">
	{#each units as unit (unit.id)}
		<button
			type="button"
			class="flex h-9 min-w-0 items-center gap-2 rounded-md px-2 text-left text-sm hover:bg-muted data-[selected=true]:bg-secondary"
			style={`padding-left: ${0.5 + unit.depth * 1.25}rem`}
			data-selected={unit.id === selectedId}
			onclick={() => onSelect(unit.id)}
		>
			<ChevronRightIcon class="size-3.5 shrink-0 text-muted-foreground" />
			<span class="truncate font-medium">{unit.name}</span>
			{#if unit.code}
				<span class="ms-auto rounded border px-1.5 py-0.5 text-xs text-muted-foreground">{unit.code}</span>
			{/if}
		</button>
	{:else}
		<div class="rounded-md border border-dashed p-6 text-center text-sm text-muted-foreground">
			No organization units yet.
		</div>
	{/each}
</div>
```

- [ ] **Step 4: Add the details panel component**

Create `frontend/src/routes/app/organization-units/unit-details-panel.svelte`:

```svelte
<script lang="ts">
	import ArchiveIcon from '@lucide/svelte/icons/archive';
	import PlusIcon from '@lucide/svelte/icons/plus';
	import SaveIcon from '@lucide/svelte/icons/save';
	import type { OrganizationUnit } from '$lib/server/organization-units';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import type { ParentOption } from './tree';

	let {
		selected,
		canManage,
		parentOptions
	}: {
		selected: OrganizationUnit | null;
		canManage: boolean;
		parentOptions: ParentOption[];
	} = $props();
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>{selected ? selected.name : 'Select a unit'}</Card.Title>
		<Card.Description>
			{#if canManage}
				Admins can update details, create children, move units, and archive subtrees.
			{:else}
				Read-only organization unit details.
			{/if}
		</Card.Description>
	</Card.Header>
	<Card.Content class="grid gap-4">
		{#if selected}
			<form method="POST" action="?/update" class="grid gap-3">
				<input type="hidden" name="id" value={selected.id} />
				<label class="grid gap-1 text-sm">
					<span class="font-medium">Name</span>
					<input name="name" value={selected.name} disabled={!canManage} class="h-9 rounded-md border bg-background px-3" />
				</label>
				<label class="grid gap-1 text-sm">
					<span class="font-medium">Code</span>
					<input name="code" value={selected.code ?? ''} disabled={!canManage} class="h-9 rounded-md border bg-background px-3" />
				</label>
				<label class="grid gap-1 text-sm">
					<span class="font-medium">Description</span>
					<textarea name="description" disabled={!canManage} class="min-h-24 rounded-md border bg-background px-3 py-2">{selected.description ?? ''}</textarea>
				</label>
				{#if canManage}
					<Button type="submit" size="sm" class="w-fit gap-2"><SaveIcon class="size-4" />Save changes</Button>
				{/if}
			</form>

			{#if canManage}
				<form method="POST" action="?/create" class="grid gap-3 rounded-md border p-3">
					<input type="hidden" name="parentId" value={selected.id} />
					<input name="name" placeholder="New child unit name" class="h-9 rounded-md border bg-background px-3" />
					<input name="code" placeholder="Optional code" class="h-9 rounded-md border bg-background px-3" />
					<input name="description" placeholder="Optional description" class="h-9 rounded-md border bg-background px-3" />
					<Button type="submit" variant="outline" size="sm" class="w-fit gap-2"><PlusIcon class="size-4" />Create child</Button>
				</form>

				<form method="POST" action="?/move" class="grid gap-3 rounded-md border p-3">
					<input type="hidden" name="id" value={selected.id} />
					<select name="parentId" class="h-9 rounded-md border bg-background px-3">
						<option value="">Root level</option>
						{#each parentOptions as option (option.id)}
							<option value={option.id}>{option.label}</option>
						{/each}
					</select>
					<Button type="submit" variant="outline" size="sm" class="w-fit">Move unit</Button>
				</form>

				<form method="POST" action="?/archive">
					<input type="hidden" name="id" value={selected.id} />
					<Button type="submit" variant="destructive" size="sm" class="gap-2">
						<ArchiveIcon class="size-4" />Archive subtree
					</Button>
				</form>
			{/if}
		{:else}
			<p class="text-sm text-muted-foreground">Choose a unit from the tree to view details.</p>
		{/if}
	</Card.Content>
</Card.Root>
```

- [ ] **Step 5: Add the page component**

Create `frontend/src/routes/app/organization-units/+page.svelte`:

```svelte
<script lang="ts">
	import AlertCircleIcon from '@lucide/svelte/icons/alert-circle';
	import PlusIcon from '@lucide/svelte/icons/plus';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import type { PageProps } from './$types';
	import UnitDetailsPanel from './unit-details-panel.svelte';
	import UnitTree from './unit-tree.svelte';
	import { findUnitById, flattenOrganizationUnitTree, selectableParentOptions } from './tree';

	let { data, form }: PageProps = $props();

	let selectedId = $state<string | null>(data.units[0]?.id ?? null);
	const flatUnits = $derived(flattenOrganizationUnitTree(data.units));
	const selected = $derived(findUnitById(data.units, selectedId));
	const parentOptions = $derived(selectableParentOptions(data.units, selectedId));
</script>

<svelte:head>
	<title>Organization Units | Syncra DMS</title>
</svelte:head>

<section class="mx-auto grid max-w-6xl gap-4 px-4 py-6">
	<div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
		<div>
			<h2 class="text-xl font-semibold">Organization Units</h2>
			<p class="text-sm text-muted-foreground">Shared company hierarchy.</p>
		</div>
		{#if data.canManageOrganizationUnits}
			<form method="POST" action="?/create" class="flex flex-wrap gap-2">
				<input name="name" placeholder="Root unit name" class="h-9 rounded-md border bg-background px-3 text-sm" />
				<input name="code" placeholder="Code" class="h-9 w-28 rounded-md border bg-background px-3 text-sm" />
				<input type="hidden" name="description" value="" />
				<input type="hidden" name="parentId" value="" />
				<Button type="submit" size="sm" class="gap-2"><PlusIcon class="size-4" />New root unit</Button>
			</form>
		{/if}
	</div>

	{#if data.loadError}
		<div class="flex items-center gap-2 rounded-md border border-destructive/30 bg-destructive/10 p-3 text-sm text-destructive">
			<AlertCircleIcon class="size-4" />
			{data.loadError}
		</div>
	{/if}

	{#if form?.error}
		<div class="flex items-center gap-2 rounded-md border border-destructive/30 bg-destructive/10 p-3 text-sm text-destructive">
			<AlertCircleIcon class="size-4" />
			{form.error}
		</div>
	{/if}

	<div class="grid gap-4 lg:grid-cols-[22rem_1fr]">
		<Card.Root>
			<Card.Header>
				<Card.Title>Active tree</Card.Title>
				<Card.Description>{flatUnits.length} active units</Card.Description>
			</Card.Header>
			<Card.Content>
				<UnitTree units={flatUnits} {selectedId} onSelect={(id) => (selectedId = id)} />
			</Card.Content>
		</Card.Root>

		<UnitDetailsPanel {selected} canManage={data.canManageOrganizationUnits} {parentOptions} />
	</div>
</section>
```

- [ ] **Step 6: Run Svelte autofixer on new components**

Run the Svelte MCP `svelte_autofixer` on:

- `+layout.svelte`
- `unit-tree.svelte`
- `unit-details-panel.svelte`
- `+page.svelte`

Apply any concrete syntax fixes it returns, then run the autofixer again on changed components until it reports no issues.

- [ ] **Step 7: Run frontend checks**

Run:

```sh
cd frontend
rtk pnpm check
rtk pnpm test -- src/routes/app/organization-units src/lib/server/organization-units.test.ts
```

Expected: PASS.

- [ ] **Step 8: Commit frontend UI**

Run:

```sh
rtk git add frontend/src/routes/app
rtk git commit -m "Add organization units app page"
```

## Task 8: End-To-End Verification And Cleanup

**Files:**
- Verify all changed backend and frontend files.
- No new source files unless verification reveals a focused defect.

- [ ] **Step 1: Run backend verification**

Run:

```sh
cd go-server
rtk go test ./...
rtk atlas migrate validate --dir file://migrations
```

Expected: PASS.

- [ ] **Step 2: Run frontend verification**

Run:

```sh
cd frontend
rtk pnpm test
rtk pnpm check
rtk pnpm build
```

Expected: PASS.

- [ ] **Step 3: Start local services for browser inspection**

Use two terminals or background sessions:

```sh
cd go-server
rtk go run ./cmd/syncra api
```

```sh
cd frontend
rtk pnpm dev -- --host 127.0.0.1
```

Expected: Go API serves on `localhost:8080`; SvelteKit prints a local Vite URL.

- [ ] **Step 4: Inspect the Organization Units page in browser**

Open the frontend URL and sign in with an admin test user. Navigate to `/app/organization-units`.

Verify:

- admin can create a root unit;
- admin can create a child unit;
- code displays uppercase after reload;
- move selector excludes the selected unit and its descendants;
- archive removes the selected subtree from the active tree;
- normal user can view the tree but does not see mutation controls.

- [ ] **Step 5: Check git status and avoid unrelated changes**

Run:

```sh
rtk git status --short
```

Expected: only Organization Unit implementation files and generated Swagger/migration checksum changes are present. Existing unrelated auth-route changes must remain untouched.

- [ ] **Step 6: Final commit if verification fixes were needed**

If Step 1 through Step 4 required a focused fix, commit that fix:

```sh
rtk git add go-server frontend
rtk git commit -m "Fix organization units verification issues"
```

If no fixes were needed, do not create an empty commit.

## Self-Review Notes

- Spec coverage: The plan covers global hierarchy, admin mutations, read-only user access, archive-only deletion, name/code/description fields, API routes, frontend tree/details layout, error mapping, tests, migration, Swagger, and browser inspection.
- Placeholder scan: The plan uses concrete files, commands, status expectations, and code snippets. It intentionally does not use open-ended implementation placeholders.
- Type consistency: Backend JSON uses `parentId`, `createdAt`, `updatedAt`, and `children`; frontend types and route components use the same names. Backend role enforcement uses existing `auth.UserRoleAdmin`; frontend role checks use `locals.user?.role === 'admin'`.
