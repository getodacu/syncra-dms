# Document Repository MVP Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build the first document repository slice: organization-unit-scoped folders, document metadata, local one-file uploads, downloads, and soft archive behavior.

**Architecture:** Add a small Go `internal/documents` domain package for models, normalization, tree helpers, and local storage. Expose trusted Gin endpoints through the existing SvelteKit server proxy pattern, then build `/app/documents` as a dense repository workspace using TanStack Query and the current app shell.

**Tech Stack:** Go, Gin, GORM, PostgreSQL/Atlas migrations, SQLite-backed Go tests, SvelteKit, Svelte 5, TypeScript, TanStack Svelte Query, Vitest, Tailwind, shadcn-style local components.

---

## Context Notes

- The repo currently has unrelated modified files. Before every task, run `rtk git status --short` and only stage files from that task.
- Prefix every shell command with `rtk`.
- Keep `example/` read-only.
- Do not expose `SYNCRA_DOCUMENT_STORAGE_ROOT` or absolute storage paths to browser code.
- Use the existing Organization Units and RBAC handlers as style references.
- The design doc is `docs/plans/2026-07-01-document-repository-mvp-design.md`.

## Task 1: Backend Document Storage Configuration

**Files:**
- Modify: `go-server/internal/config/config.go`
- Modify: `go-server/internal/config/config_test.go`
- Modify: `go-server/internal/app/api.go`
- Modify: `go-server/internal/api/router.go`
- Modify: `go-server/internal/api/auth_handlers_test.go`
- Modify: `go-server/.env.example`

**Step 1: Write failing config tests**

Add tests in `go-server/internal/config/config_test.go`:

```go
func TestLoadRequiresDocumentStorageRoot(t *testing.T) {
	t.Setenv("DSN", `host=localhost dbname=syncra_dms`)
	t.Setenv("DSN_DEV", `host=localhost dbname=syncra_dms_dev`)
	t.Setenv("SYNCRA_DOCUMENT_STORAGE_ROOT", "")

	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "SYNCRA_DOCUMENT_STORAGE_ROOT is required") {
		t.Fatalf("Load error = %v, want document storage root required", err)
	}
}

func TestLoadReadsDocumentStorageSettings(t *testing.T) {
	t.Setenv("DSN", `host=localhost dbname=syncra_dms`)
	t.Setenv("DSN_DEV", `host=localhost dbname=syncra_dms_dev`)
	t.Setenv("SYNCRA_DOCUMENT_STORAGE_ROOT", "/var/lib/syncra/documents")
	t.Setenv("SYNCRA_DOCUMENT_MAX_UPLOAD_BYTES", "1048576")
	t.Setenv("SYNCRA_DOCUMENT_ALLOWED_MIME_TYPES", "application/pdf,image/png")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load error = %v", err)
	}
	if cfg.DocumentStorageRoot != "/var/lib/syncra/documents" {
		t.Fatalf("DocumentStorageRoot = %q", cfg.DocumentStorageRoot)
	}
	if cfg.DocumentMaxUploadBytes != 1048576 {
		t.Fatalf("DocumentMaxUploadBytes = %d", cfg.DocumentMaxUploadBytes)
	}
	if !reflect.DeepEqual(cfg.DocumentAllowedMIMETypes, []string{"application/pdf", "image/png"}) {
		t.Fatalf("DocumentAllowedMIMETypes = %#v", cfg.DocumentAllowedMIMETypes)
	}
}
```

**Step 2: Run tests to verify failure**

Run:

```bash
cd go-server
rtk go test ./internal/config
```

Expected: FAIL because `Config` lacks document storage fields.

**Step 3: Implement config fields and router options**

Add to `config.Config`:

```go
DocumentStorageRoot       string
DocumentMaxUploadBytes    int64
DocumentAllowedMIMETypes  []string
```

Add a parser:

```go
func getenvCSV(key string) []string {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value != "" {
			values = append(values, value)
		}
	}
	return values
}
```

In `Load`, read:

```go
DocumentStorageRoot:      strings.TrimSpace(os.Getenv("SYNCRA_DOCUMENT_STORAGE_ROOT")),
DocumentMaxUploadBytes:   getenvInt64("SYNCRA_DOCUMENT_MAX_UPLOAD_BYTES", 25*1024*1024),
DocumentAllowedMIMETypes: getenvCSV("SYNCRA_DOCUMENT_ALLOWED_MIME_TYPES"),
```

After DSN checks, require:

```go
if cfg.DocumentStorageRoot == "" {
	return Config{}, errors.New("SYNCRA_DOCUMENT_STORAGE_ROOT is required")
}
```

Add to `api.RouterOptions`:

```go
DocumentStorageRoot      string
DocumentMaxUploadBytes   int64
DocumentAllowedMIMETypes []string
```

Pass those from `app.RunAPI` into `api.NewRouter`.

In `newAuthTestRouterWithOptions`, default test options:

```go
DocumentStorageRoot:    t.TempDir(),
DocumentMaxUploadBytes: 25 * 1024 * 1024,
```

Also pass through overrides from `options` when non-zero/non-empty.

**Step 4: Update environment example**

Add to `go-server/.env.example`:

```dotenv
# Local filesystem storage for uploaded document bytes. Use a path outside the repo.
SYNCRA_DOCUMENT_STORAGE_ROOT="/var/lib/syncra-dms/documents"

# Optional upload limits. Defaults to 25 MiB when omitted.
SYNCRA_DOCUMENT_MAX_UPLOAD_BYTES=26214400
SYNCRA_DOCUMENT_ALLOWED_MIME_TYPES=""
```

**Step 5: Run tests**

Run:

```bash
cd go-server
rtk go test ./internal/config ./internal/app ./internal/api
```

Expected: PASS.

**Step 6: Commit**

```bash
rtk git add go-server/internal/config/config.go go-server/internal/config/config_test.go go-server/internal/app/api.go go-server/internal/api/router.go go-server/internal/api/auth_handlers_test.go go-server/.env.example
rtk git commit -m "Add document storage configuration"
```

## Task 2: Document Domain Models And Tree Helpers

**Files:**
- Create: `go-server/internal/documents/models.go`
- Create: `go-server/internal/documents/models_test.go`
- Create: `go-server/internal/documents/tree.go`
- Create: `go-server/internal/documents/tree_test.go`
- Modify: `go-server/internal/database/database.go`
- Modify: `go-server/internal/api/auth_handlers_test.go`

**Step 1: Write failing model and tree tests**

In `models_test.go`, test IDs and table names:

```go
func TestDocumentModelsAssignIDs(t *testing.T) {
	db := sqliteMemoryDB(t)
	if err := db.AutoMigrate(&Folder{}, &Document{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	folder := Folder{Name: "Invoices", OrganizationUnitID: uuid.NewString(), CreatedByUserID: uuid.NewString(), CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := db.Create(&folder).Error; err != nil {
		t.Fatalf("create folder: %v", err)
	}
	if folder.ID == "" {
		t.Fatal("folder ID was empty")
	}
}
```

In `tree_test.go`:

```go
func TestBuildFolderTreeOrdersAndNestsFolders(t *testing.T) {
	rootID := "root"
	childID := "child"
	nodes := BuildFolderTree([]Folder{
		{ID: childID, ParentID: &rootID, Name: "Beta", CreatedAt: time.Unix(2, 0), UpdatedAt: time.Unix(2, 0)},
		{ID: rootID, Name: "Root", CreatedAt: time.Unix(1, 0), UpdatedAt: time.Unix(1, 0)},
		{ID: "alpha", ParentID: &rootID, Name: "Alpha", CreatedAt: time.Unix(3, 0), UpdatedAt: time.Unix(3, 0)},
	})
	if len(nodes) != 1 || len(nodes[0].Children) != 2 || nodes[0].Children[0].ID != "alpha" || nodes[0].Children[1].ID != childID {
		t.Fatalf("tree = %#v", nodes)
	}
}
```

**Step 2: Run tests to verify failure**

```bash
cd go-server
rtk go test ./internal/documents
```

Expected: FAIL because the package does not exist.

**Step 3: Implement models**

Create `models.go`:

```go
package documents

import (
	"errors"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"ai.ro/syncra/dms/internal/auth"
	"ai.ro/syncra/dms/internal/orgunits"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const MaxFolderNameCharacters = 160
const MaxDocumentDisplayNameCharacters = 255

type Folder struct {
	ID                 string         `gorm:"type:uuid;primaryKey;index:idx_document_folders_parent_name_id,priority:4" json:"id"`
	ParentID           *string        `gorm:"column:parent_id;type:uuid;index:idx_document_folders_parent_name_id,priority:2;uniqueIndex:idx_document_folders_active_name_unique,priority:2,where:deleted_at IS NULL" json:"parentId,omitempty"`
	Parent             *Folder        `gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
	OrganizationUnitID string         `gorm:"column:organization_unit_id;type:uuid;not null;index:idx_document_folders_parent_name_id,priority:1;uniqueIndex:idx_document_folders_active_name_unique,priority:1,where:deleted_at IS NULL" json:"organizationUnitId"`
	OrganizationUnit   orgunits.Unit  `gorm:"foreignKey:OrganizationUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
	Name               string         `gorm:"not null;size:160;index:idx_document_folders_parent_name_id,priority:3;uniqueIndex:idx_document_folders_active_name_unique,priority:3,where:deleted_at IS NULL" json:"name"`
	Description        *string        `gorm:"type:text" json:"description,omitempty"`
	CreatedByUserID    string         `gorm:"column:created_by_user_id;type:uuid;not null;index" json:"createdByUserId"`
	CreatedByUser      auth.User      `gorm:"foreignKey:CreatedByUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
	UpdatedByUserID    *string        `gorm:"column:updated_by_user_id;type:uuid;index" json:"updatedByUserId,omitempty"`
	UpdatedByUser      *auth.User     `gorm:"foreignKey:UpdatedByUserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
	DeletedAt          *time.Time     `gorm:"column:deleted_at;index" json:"deletedAt,omitempty"`
	CreatedAt          time.Time      `gorm:"column:created_at;not null" json:"createdAt"`
	UpdatedAt          time.Time      `gorm:"column:updated_at;not null" json:"updatedAt"`
}

func (f *Folder) BeforeCreate(_ *gorm.DB) error {
	if f.ID == "" {
		f.ID = uuid.NewString()
	}
	return nil
}

func (Folder) TableName() string { return "document_folders" }

type Document struct {
	ID                 string     `gorm:"type:uuid;primaryKey" json:"id"`
	FolderID           string     `gorm:"column:folder_id;type:uuid;not null;index;uniqueIndex:idx_documents_active_folder_hash_unique,priority:1,where:deleted_at IS NULL" json:"folderId"`
	Folder             Folder     `gorm:"foreignKey:FolderID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
	OrganizationUnitID string     `gorm:"column:organization_unit_id;type:uuid;not null;index" json:"organizationUnitId"`
	OrganizationUnit   orgunits.Unit `gorm:"foreignKey:OrganizationUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
	OriginalFileName   string     `gorm:"column:original_file_name;not null;size:255" json:"originalFileName"`
	DisplayName        string     `gorm:"column:display_name;not null;size:255" json:"displayName"`
	MimeType           string     `gorm:"column:mime_type;not null;size:255" json:"mimeType"`
	Extension          *string    `gorm:"size:32" json:"extension,omitempty"`
	SizeBytes          int64      `gorm:"column:size_bytes;not null" json:"sizeBytes"`
	SHA256Hash         string     `gorm:"column:sha256_hash;not null;size:64;uniqueIndex:idx_documents_active_folder_hash_unique,priority:2,where:deleted_at IS NULL" json:"sha256Hash"`
	StorageKey         string     `gorm:"column:storage_key;not null;type:text" json:"-"`
	CreatedByUserID    string     `gorm:"column:created_by_user_id;type:uuid;not null;index" json:"createdByUserId"`
	CreatedByUser      auth.User  `gorm:"foreignKey:CreatedByUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
	UpdatedByUserID    *string    `gorm:"column:updated_by_user_id;type:uuid;index" json:"updatedByUserId,omitempty"`
	UpdatedByUser      *auth.User `gorm:"foreignKey:UpdatedByUserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
	DeletedAt          *time.Time `gorm:"column:deleted_at;index" json:"deletedAt,omitempty"`
	CreatedAt          time.Time  `gorm:"column:created_at;not null" json:"createdAt"`
	UpdatedAt          time.Time  `gorm:"column:updated_at;not null" json:"updatedAt"`
}

func (d *Document) BeforeCreate(_ *gorm.DB) error {
	if d.ID == "" {
		d.ID = uuid.NewString()
	}
	return nil
}

func (Document) TableName() string { return "documents" }
```

Add normalization helpers:

```go
func NormalizeFolderName(raw string) (string, error) {
	name := strings.TrimSpace(raw)
	if name == "" {
		return "", errors.New("folder name is required")
	}
	if utf8.RuneCountInString(name) > MaxFolderNameCharacters {
		return "", errors.New("folder name must be 160 characters or fewer")
	}
	return name, nil
}

func NormalizeDescription(raw string) *string {
	value := strings.TrimSpace(raw)
	if value == "" {
		return nil
	}
	return &value
}

func NormalizeDisplayName(raw string) (string, error) {
	name := strings.TrimSpace(raw)
	if name == "" {
		return "", errors.New("document name is required")
	}
	if utf8.RuneCountInString(name) > MaxDocumentDisplayNameCharacters {
		return "", errors.New("document name must be 255 characters or fewer")
	}
	return name, nil
}

func SafeOriginalFileName(raw string) string {
	name := strings.TrimSpace(filepath.Base(raw))
	if name == "." || name == "/" || name == "" {
		return "upload"
	}
	if utf8.RuneCountInString(name) > 255 {
		return string([]rune(name)[:255])
	}
	return name
}
```

**Step 4: Implement tree helpers**

Create `tree.go` with `FolderTreeNode`, `BuildFolderTree`, and `DescendantFolderIDs` mirroring `internal/orgunits/tree.go`, sorting by `Name` then `ID`.

**Step 5: Register models**

Add `&documents.Folder{}` and `&documents.Document{}` to `database.ApplicationModels()` and the test `AutoMigrate` list in `auth_handlers_test.go`.

**Step 6: Run tests**

```bash
cd go-server
rtk go test ./internal/documents ./internal/database ./internal/api
```

Expected: PASS.

**Step 7: Commit**

```bash
rtk git add go-server/internal/documents go-server/internal/database/database.go go-server/internal/api/auth_handlers_test.go
rtk git commit -m "Add document repository domain models"
```

## Task 3: Atlas Migration For Repository Tables

**Files:**
- Create: `go-server/migrations/20260701090000_add_document_repository.sql`
- Modify: `go-server/migrations/atlas.sum`
- Modify: `go-server/internal/dbschema/schema.go` only if model registration needs adjustment from Task 2

**Step 1: Write migration SQL**

Create the migration:

```sql
CREATE TABLE "document_folders" (
  "id" uuid NOT NULL,
  "parent_id" uuid NULL,
  "organization_unit_id" uuid NOT NULL,
  "name" varchar(160) NOT NULL,
  "description" text NULL,
  "created_by_user_id" uuid NOT NULL,
  "updated_by_user_id" uuid NULL,
  "deleted_at" timestamptz NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_document_folders_parent" FOREIGN KEY ("parent_id") REFERENCES "document_folders" ("id") ON UPDATE CASCADE ON DELETE RESTRICT,
  CONSTRAINT "fk_document_folders_organization_unit" FOREIGN KEY ("organization_unit_id") REFERENCES "organization_units" ("id") ON UPDATE CASCADE ON DELETE RESTRICT,
  CONSTRAINT "fk_document_folders_created_by_user" FOREIGN KEY ("created_by_user_id") REFERENCES "user" ("id") ON UPDATE CASCADE ON DELETE RESTRICT,
  CONSTRAINT "fk_document_folders_updated_by_user" FOREIGN KEY ("updated_by_user_id") REFERENCES "user" ("id") ON UPDATE CASCADE ON DELETE SET NULL
);

CREATE INDEX "idx_document_folders_parent_name_id" ON "document_folders" ("organization_unit_id", "parent_id", "name", "id");
CREATE INDEX "idx_document_folders_deleted_at" ON "document_folders" ("deleted_at");
CREATE UNIQUE INDEX "idx_document_folders_active_name_unique" ON "document_folders" ("organization_unit_id", "parent_id", "name") WHERE "deleted_at" IS NULL;

CREATE TABLE "documents" (
  "id" uuid NOT NULL,
  "folder_id" uuid NOT NULL,
  "organization_unit_id" uuid NOT NULL,
  "original_file_name" varchar(255) NOT NULL,
  "display_name" varchar(255) NOT NULL,
  "mime_type" varchar(255) NOT NULL,
  "extension" varchar(32) NULL,
  "size_bytes" bigint NOT NULL,
  "sha256_hash" char(64) NOT NULL,
  "storage_key" text NOT NULL,
  "created_by_user_id" uuid NOT NULL,
  "updated_by_user_id" uuid NULL,
  "deleted_at" timestamptz NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_documents_folder" FOREIGN KEY ("folder_id") REFERENCES "document_folders" ("id") ON UPDATE CASCADE ON DELETE RESTRICT,
  CONSTRAINT "fk_documents_organization_unit" FOREIGN KEY ("organization_unit_id") REFERENCES "organization_units" ("id") ON UPDATE CASCADE ON DELETE RESTRICT,
  CONSTRAINT "fk_documents_created_by_user" FOREIGN KEY ("created_by_user_id") REFERENCES "user" ("id") ON UPDATE CASCADE ON DELETE RESTRICT,
  CONSTRAINT "fk_documents_updated_by_user" FOREIGN KEY ("updated_by_user_id") REFERENCES "user" ("id") ON UPDATE CASCADE ON DELETE SET NULL
);

CREATE INDEX "idx_documents_folder_id" ON "documents" ("folder_id");
CREATE INDEX "idx_documents_organization_unit_id" ON "documents" ("organization_unit_id");
CREATE INDEX "idx_documents_deleted_at" ON "documents" ("deleted_at");
CREATE UNIQUE INDEX "idx_documents_active_folder_hash_unique" ON "documents" ("folder_id", "sha256_hash") WHERE "deleted_at" IS NULL;
```

**Step 2: Validate and update Atlas checksum**

Run:

```bash
cd go-server
rtk atlas migrate hash --dir file://migrations
rtk atlas migrate validate --dir file://migrations
```

Expected: PASS and `atlas.sum` updated.

**Step 3: Run schema-related tests**

```bash
cd go-server
rtk go test ./internal/dbschema ./internal/database
```

Expected: PASS.

**Step 4: Commit**

```bash
rtk git add go-server/migrations/20260701090000_add_document_repository.sql go-server/migrations/atlas.sum go-server/internal/dbschema/schema.go
rtk git commit -m "Add document repository migration"
```

## Task 4: Seed Document Permissions

**Files:**
- Modify: `go-server/internal/rbac/registry.go`
- Modify: `go-server/internal/rbac/seed_test.go`
- Modify: `frontend/src/routes/layout-permissions.test.ts` later in Task 12 for nav checks

**Step 1: Write failing RBAC tests**

In `seed_test.go`, extend the registry/default-role assertions:

```go
func TestDocumentPermissionsAreSeeded(t *testing.T) {
	db := rbacTestDB(t)
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
```

Extend the organization administrator expected permissions to include all five document permissions.

**Step 2: Run tests to verify failure**

```bash
cd go-server
rtk go test ./internal/rbac
```

Expected: FAIL because document permissions are missing.

**Step 3: Add permission definitions**

In `registry.go`, append:

```go
{Code: "document.view", Name: "View documents", Category: "Document Repository"},
{Code: "document.create", Name: "Create document folders and uploads", Category: "Document Repository"},
{Code: "document.update", Name: "Update document folders and metadata", Category: "Document Repository"},
{Code: "document.delete", Name: "Archive document folders and documents", Category: "Document Repository"},
{Code: "document.download", Name: "Download documents", Category: "Document Repository"},
```

Add those codes to `organizationAdministratorPermissionCodes`. Do not add them to the generic Viewer role unless a product decision explicitly says viewers should see all documents.

**Step 4: Run tests**

```bash
cd go-server
rtk go test ./internal/rbac ./internal/api
```

Expected: PASS.

**Step 5: Commit**

```bash
rtk git add go-server/internal/rbac/registry.go go-server/internal/rbac/seed_test.go
rtk git commit -m "Seed document repository permissions"
```

## Task 5: Folder API

**Files:**
- Create: `go-server/internal/api/document_folders.go`
- Create: `go-server/internal/api/document_folders_test.go`
- Modify: `go-server/internal/api/router.go`
- Modify: `go-server/internal/api/swagger_doc.go`
- Modify: `go-server/internal/api/swagger_doc_test.go`

**Step 1: Write failing route and behavior tests**

Create tests covering:

```go
func TestDocumentFolderRoutesRequirePermissions(t *testing.T) {
	router, db := newAuthTestRouter(t)
	unitID := createUnitViaAPI(t, router, loginSeededAdmin(t, router, db, "admin@example.com"), `{"name":"Finance"}`)
	user := createVerifiedUser(t, db, "viewer@example.com", "password123")
	token := loginUser(t, router, user.Email, "password123")

	response := folderJSON(t, router, http.MethodGet, "/api/document-folders/tree?organizationUnitId="+unitID, "", authCookieHeaders(token))
	if response.Code != http.StatusForbidden {
		t.Fatalf("tree status = %d body=%s, want forbidden", response.Code, response.Body.String())
	}
}

func TestDocumentFolderLifecycle(t *testing.T) {
	router, db := newAuthTestRouter(t)
	token := loginSeededAdmin(t, router, db, "admin@example.com")
	unitID := createUnitViaAPI(t, router, token, `{"name":"Finance"}`)

	rootID := createFolderViaAPI(t, router, token, `{"organizationUnitId":"`+unitID+`","name":"Invoices"}`)
	childID := createFolderViaAPI(t, router, token, `{"organizationUnitId":"`+unitID+`","parentId":"`+rootID+`","name":"2026"}`)

	tree := folderJSON(t, router, http.MethodGet, "/api/document-folders/tree?organizationUnitId="+unitID, "", authCookieHeaders(token))
	if tree.Code != http.StatusOK {
		t.Fatalf("tree status = %d body=%s", tree.Code, tree.Body.String())
	}
	// Decode and assert root/child nesting.

	move := folderJSON(t, router, http.MethodPatch, "/api/document-folders/"+childID+"/parent", `{"parentId":null}`, authCookieHeaders(token))
	if move.Code != http.StatusOK {
		t.Fatalf("move status = %d body=%s", move.Code, move.Body.String())
	}

	archive := folderJSON(t, router, http.MethodPost, "/api/document-folders/"+rootID+"/archive", `{}`, authCookieHeaders(token))
	if archive.Code != http.StatusOK {
		t.Fatalf("archive status = %d body=%s", archive.Code, archive.Body.String())
	}
}
```

Also test duplicate active folder names, invalid parent IDs, move cycles, cross-unit moves, and archived-folder exclusion.

**Step 2: Run tests to verify failure**

```bash
cd go-server
rtk go test ./internal/api -run DocumentFolder
```

Expected: FAIL because routes do not exist.

**Step 3: Implement handler**

Mirror `organization_units.go` but use organization-unit-scoped permission checks:

```go
type documentFolderHandler struct {
	db   *gorm.DB
	auth *authHandler
}

type documentFolderRequest struct {
	OrganizationUnitID string  `json:"organizationUnitId"`
	ParentID           *string `json:"parentId"`
	Name               string  `json:"name"`
	Description        *string `json:"description"`
}

type moveDocumentFolderRequest struct {
	ParentID *string
}
```

Implement:

- `listTree`: parse `organizationUnitId`, require `document.view`, return active folder tree for that unit.
- `create`: validate unit/parent, require `document.create`, set creator/updater, create folder.
- `update`: load active folder first, require `document.update` for its unit, update name/description.
- `move`: load active folder, require `document.update`, validate parent in same unit, reject self/descendant.
- `archive`: load active folder, require `document.delete`, set `deleted_at` recursively for folder subtree and active documents under those folders.

Response shape:

```go
type documentFolderResponse struct {
	ID                 string                   `json:"id"`
	ParentID           *string                  `json:"parentId,omitempty"`
	OrganizationUnitID string                   `json:"organizationUnitId"`
	Name               string                   `json:"name"`
	Description        *string                  `json:"description,omitempty"`
	DeletedAt          *string                  `json:"deletedAt,omitempty"`
	CreatedAt          string                   `json:"createdAt"`
	UpdatedAt          string                   `json:"updatedAt"`
	Children           []documentFolderResponse `json:"children"`
}
```

Use `documents.NormalizeFolderName`, `documents.NormalizeDescription`, `documents.BuildFolderTree`, and `documents.DescendantFolderIDs`.

**Step 4: Register routes**

In `router.go`:

```go
documentFolders := newDocumentFolderHandler(options, auth)
folderAPI := router.Group("/api/document-folders")
folderAPI.Use(auth.requireTrustedInternalRequest())
folderAPI.GET("/tree", documentFolders.listTree)
folderAPI.POST("", documentFolders.create)
folderAPI.PATCH("/:id", documentFolders.update)
folderAPI.PATCH("/:id/parent", documentFolders.move)
folderAPI.POST("/:id/archive", documentFolders.archive)
folderAPI.GET("/:id/contents", documentFolders.contents) // stub can return empty arrays until Task 6
```

**Step 5: Add Swagger docs**

Add route entries in `swagger_doc.go` and update `swagger_doc_test.go` expected route list for all folder endpoints.

**Step 6: Run tests**

```bash
cd go-server
rtk go test ./internal/api ./internal/documents
```

Expected: PASS.

**Step 7: Commit**

```bash
rtk git add go-server/internal/api/document_folders.go go-server/internal/api/document_folders_test.go go-server/internal/api/router.go go-server/internal/api/swagger_doc.go go-server/internal/api/swagger_doc_test.go
rtk git commit -m "Add document folder API"
```

## Task 6: Folder Contents API

**Files:**
- Modify: `go-server/internal/api/document_folders.go`
- Modify: `go-server/internal/api/document_folders_test.go`

**Step 1: Write failing contents tests**

Add a test that creates a folder, a child folder, and two document rows, then calls:

```go
GET /api/document-folders/:id/contents
```

Expected JSON:

```json
{
  "folder": { "id": "folder-id", "name": "Invoices" },
  "folders": [{ "id": "child-id", "name": "2026" }],
  "documents": [{ "id": "doc-id", "displayName": "invoice.pdf" }]
}
```

Also assert archived child folders/documents are excluded and `document.view` is required for the folder organization unit.

**Step 2: Run test to verify failure**

```bash
cd go-server
rtk go test ./internal/api -run DocumentFolderContents
```

Expected: FAIL because contents returns a stub or missing response.

**Step 3: Implement contents**

Add response types:

```go
type documentFolderContentsResponse struct {
	Folder    documentFolderResponse    `json:"folder"`
	Folders   []documentFolderResponse  `json:"folders"`
	Documents []documentMetadataResponse `json:"documents"`
}
```

Load active folder by ID, require `document.view`, query active child folders and active documents ordered by `name/display_name asc, id asc`.

Document metadata response must include:

```go
type documentMetadataResponse struct {
	ID                 string  `json:"id"`
	FolderID           string  `json:"folderId"`
	OrganizationUnitID string  `json:"organizationUnitId"`
	OriginalFileName   string  `json:"originalFileName"`
	DisplayName        string  `json:"displayName"`
	MimeType           string  `json:"mimeType"`
	Extension          *string `json:"extension,omitempty"`
	SizeBytes          int64   `json:"sizeBytes"`
	SHA256Hash         string  `json:"sha256Hash"`
	DeletedAt          *string `json:"deletedAt,omitempty"`
	CreatedAt          string  `json:"createdAt"`
	UpdatedAt          string  `json:"updatedAt"`
}
```

Do not include `StorageKey`.

**Step 4: Run tests**

```bash
cd go-server
rtk go test ./internal/api -run DocumentFolder
```

Expected: PASS.

**Step 5: Commit**

```bash
rtk git add go-server/internal/api/document_folders.go go-server/internal/api/document_folders_test.go
rtk git commit -m "Add document folder contents API"
```

## Task 7: Local Document Storage Service

**Files:**
- Create: `go-server/internal/documents/storage.go`
- Create: `go-server/internal/documents/storage_test.go`

**Step 1: Write failing storage tests**

Test that storage writes bytes under a generated key, computes SHA-256, detects MIME type, prevents traversal, and opens stored content:

```go
func TestLocalStorageStoresAndOpensDocument(t *testing.T) {
	root := t.TempDir()
	store := NewLocalStorage(root, 1024, nil)
	result, err := store.Save(context.Background(), strings.NewReader("hello"), "invoice.pdf")
	if err != nil {
		t.Fatalf("Save error = %v", err)
	}
	if result.SizeBytes != 5 || result.SHA256Hash == "" || result.StorageKey == "" {
		t.Fatalf("result = %#v", result)
	}
	reader, err := store.Open(result.StorageKey)
	if err != nil {
		t.Fatalf("Open error = %v", err)
	}
	defer reader.Close()
	body, _ := io.ReadAll(reader)
	if string(body) != "hello" {
		t.Fatalf("body = %q", body)
	}
}
```

Add tests for max size and allowed MIME type rejection.

**Step 2: Run tests to verify failure**

```bash
cd go-server
rtk go test ./internal/documents -run LocalStorage
```

Expected: FAIL because storage does not exist.

**Step 3: Implement local storage**

Create:

```go
type LocalStorage struct {
	root             string
	maxUploadBytes   int64
	allowedMIMETypes map[string]bool
}

type StoredFile struct {
	OriginalFileName string
	MimeType         string
	Extension        *string
	SizeBytes        int64
	SHA256Hash       string
	StorageKey       string
}
```

`Save` should:

- create `root/tmp`;
- stream into a temp file while hashing and counting bytes;
- reject if size exceeds `maxUploadBytes`;
- sniff MIME from the first 512 bytes with `http.DetectContentType`;
- reject MIME if an allow-list is configured and the sniffed type is absent;
- move the temp file to `root/documents/<sha-prefix>/<uuid>.bin`;
- return a storage key relative to root, using forward slashes.

`Open(storageKey string)` should clean and join the key with root, reject traversal, and return `*os.File`.

**Step 4: Run tests**

```bash
cd go-server
rtk go test ./internal/documents
```

Expected: PASS.

**Step 5: Commit**

```bash
rtk git add go-server/internal/documents/storage.go go-server/internal/documents/storage_test.go
rtk git commit -m "Add local document storage"
```

## Task 8: Document Upload API

**Files:**
- Create: `go-server/internal/api/documents.go`
- Create: `go-server/internal/api/documents_test.go`
- Modify: `go-server/internal/api/router.go`
- Modify: `go-server/internal/api/swagger_doc.go`
- Modify: `go-server/internal/api/swagger_doc_test.go`

**Step 1: Write failing upload tests**

Use `multipart.NewWriter` in `documents_test.go`:

```go
func TestDocumentUploadStoresMetadataAndFile(t *testing.T) {
	router, db := newAuthTestRouterWithOptions(t, RouterOptions{DocumentStorageRoot: t.TempDir()})
	token := loginSeededAdmin(t, router, db, "admin@example.com")
	unitID := createUnitViaAPI(t, router, token, `{"name":"Finance"}`)
	folderID := createFolderViaAPI(t, router, token, `{"organizationUnitId":"`+unitID+`","name":"Invoices"}`)

	response := uploadDocument(t, router, token, folderID, "invoice.pdf", []byte("%PDF test"))
	if response.Code != http.StatusCreated {
		t.Fatalf("upload status = %d body=%s", response.Code, response.Body.String())
	}
	var body documentMetadataResponse
	decodeJSON(t, response, &body)
	if body.FolderID != folderID || body.OrganizationUnitID != unitID || body.SHA256Hash == "" || body.SizeBytes == 0 {
		t.Fatalf("upload body = %#v", body)
	}
	var count int64
	if err := db.Model(&documents.Document{}).Where("folder_id = ?", folderID).Count(&count).Error; err != nil {
		t.Fatalf("count documents: %v", err)
	}
	if count != 1 {
		t.Fatalf("document count = %d", count)
	}
}
```

Add tests for missing file, missing folder, archived folder, missing `document.create`, duplicate hash in same folder, max upload size, and unsupported MIME type.

**Step 2: Run test to verify failure**

```bash
cd go-server
rtk go test ./internal/api -run DocumentUpload
```

Expected: FAIL because upload endpoint does not exist.

**Step 3: Implement document handler**

In `documents.go`:

```go
type documentHandler struct {
	db      *gorm.DB
	auth    *authHandler
	storage *documents.LocalStorage
}
```

Build storage in `newDocumentHandler(options, auth)` from `RouterOptions`.

Implement `upload`:

- require trusted internal request and session;
- limit request with `http.MaxBytesReader`;
- parse multipart form;
- read `folderId` form value and `file` file part;
- load active folder;
- require `document.create` for folder organization unit;
- call `storage.Save`;
- check duplicate active document by `folder_id` and `sha256_hash`;
- create `documents.Document` with folder unit, user ID, metadata, and storage key;
- if duplicate after storage save, remove stored temp/final file if a delete helper exists, then return `409`.

Register:

```go
documents := newDocumentHandler(options, auth)
documentAPI := router.Group("/api/documents")
documentAPI.Use(auth.requireTrustedInternalRequest())
documentAPI.POST("/upload", documents.upload)
documentAPI.GET("/:id", documents.get)
documentAPI.GET("/:id/download", documents.download)
documentAPI.PATCH("/:id", documents.update)
documentAPI.POST("/:id/archive", documents.archive)
```

For this task, `get/download/update/archive` may return `501` only if tests do not touch them. Prefer implementing stubs that compile and are replaced in Task 9.

**Step 4: Update Swagger docs**

Add all document endpoint route names to `swagger_doc.go` and `swagger_doc_test.go`.

**Step 5: Run tests**

```bash
cd go-server
rtk go test ./internal/api ./internal/documents
```

Expected: PASS.

**Step 6: Commit**

```bash
rtk git add go-server/internal/api/documents.go go-server/internal/api/documents_test.go go-server/internal/api/router.go go-server/internal/api/swagger_doc.go go-server/internal/api/swagger_doc_test.go
rtk git commit -m "Add document upload API"
```

## Task 9: Document Metadata, Archive, Rename, And Download API

**Files:**
- Modify: `go-server/internal/api/documents.go`
- Modify: `go-server/internal/api/documents_test.go`
- Modify: `go-server/internal/documents/storage.go`
- Modify: `go-server/internal/documents/storage_test.go`

**Step 1: Write failing document operation tests**

Add tests:

```go
func TestDocumentMetadataRenameArchiveAndDownload(t *testing.T) {
	router, db := newAuthTestRouterWithOptions(t, RouterOptions{DocumentStorageRoot: t.TempDir()})
	token := loginSeededAdmin(t, router, db, "admin@example.com")
	unitID := createUnitViaAPI(t, router, token, `{"name":"Finance"}`)
	folderID := createFolderViaAPI(t, router, token, `{"organizationUnitId":"`+unitID+`","name":"Invoices"}`)
	docID := uploadDocumentID(t, router, token, folderID, "invoice.pdf", []byte("hello"))

	get := authJSON(t, router, http.MethodGet, "/api/documents/"+docID, "", authCookieHeaders(token))
	if get.Code != http.StatusOK {
		t.Fatalf("get status = %d body=%s", get.Code, get.Body.String())
	}

	rename := authJSON(t, router, http.MethodPatch, "/api/documents/"+docID, `{"displayName":"Renamed.pdf"}`, authCookieHeaders(token))
	if rename.Code != http.StatusOK {
		t.Fatalf("rename status = %d body=%s", rename.Code, rename.Body.String())
	}

	download := authJSON(t, router, http.MethodGet, "/api/documents/"+docID+"/download", "", authCookieHeaders(token))
	if download.Code != http.StatusOK || download.Body.String() != "hello" {
		t.Fatalf("download status = %d body=%q", download.Code, download.Body.String())
	}

	archive := authJSON(t, router, http.MethodPost, "/api/documents/"+docID+"/archive", `{}`, authCookieHeaders(token))
	if archive.Code != http.StatusOK {
		t.Fatalf("archive status = %d body=%s", archive.Code, archive.Body.String())
	}
}
```

Add negative tests for missing `document.view`, `document.update`, `document.delete`, `document.download`, archived documents, missing storage file, and invalid display names.

**Step 2: Run tests to verify failure**

```bash
cd go-server
rtk go test ./internal/api -run DocumentMetadata
```

Expected: FAIL because stubs or missing logic remain.

**Step 3: Implement operations**

Implement:

- `get`: load active document, require `document.view`, return metadata response.
- `update`: load active document, require `document.update`, normalize `displayName`, set updater/timestamp.
- `archive`: load active document, require `document.delete`, set `deleted_at` and updater/timestamp.
- `download`: load active document, require `document.download`, open storage key, stream with:

```go
c.Header("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": doc.OriginalFileName}))
c.Header("Content-Type", doc.MimeType)
c.Header("Content-Length", strconv.FormatInt(doc.SizeBytes, 10))
c.FileFromFS(storagePath, http.FS(...)) // or io.Copy from storage.Open
```

Prefer `io.Copy(c.Writer, reader)` so storage key resolution stays inside `LocalStorage`.

**Step 4: Run tests**

```bash
cd go-server
rtk go test ./internal/api ./internal/documents
```

Expected: PASS.

**Step 5: Commit**

```bash
rtk git add go-server/internal/api/documents.go go-server/internal/api/documents_test.go go-server/internal/documents/storage.go go-server/internal/documents/storage_test.go
rtk git commit -m "Add document metadata and download API"
```

## Task 10: Backend Verification

**Files:**
- Modify as needed from failures only.

**Step 1: Run full backend tests**

```bash
cd go-server
rtk go test ./...
```

Expected: PASS.

**Step 2: Validate migrations**

```bash
cd go-server
rtk atlas migrate validate --dir file://migrations
```

Expected: PASS.

**Step 3: Fix failures with focused tests**

For each failure:

1. Write or adjust the smallest failing test proving the intended behavior.
2. Implement the smallest fix.
3. Rerun the focused package.
4. Rerun `rtk go test ./...`.

**Step 4: Commit fixes if any**

```bash
rtk git add <changed-backend-files>
rtk git commit -m "Stabilize document repository backend"
```

Skip the commit if there were no changes.

## Task 11: SvelteKit Server Document Client

**Files:**
- Create: `frontend/src/lib/server/documents.ts`
- Create: `frontend/src/lib/server/documents.test.ts`

**Step 1: Write failing server-client tests**

In `documents.test.ts`:

```ts
it('fetches folder contents with internal and cookie headers', async () => {
	const fetch = vi.fn(async () =>
		jsonResponse({
			folder: folder(),
			folders: [],
			documents: []
		})
	);

	const result = await getDocumentFolderContents(fetch, 'session=token', 'folder/id');

	expect(fetch).toHaveBeenCalledWith(
		'http://api.test/api/document-folders/folder%2Fid/contents',
		expect.objectContaining({
			method: 'GET',
			headers: expect.any(Headers)
		})
	);
	expect(result.folder.id).toBe('folder-id');
});
```

Add tests for tree, create/update/move/archive folder, upload using `FormData`, metadata get/update/archive, download returning a `Response`, invalid payload rejection, and public error mapping.

**Step 2: Run tests to verify failure**

```bash
cd frontend
rtk pnpm test -- src/lib/server/documents.test.ts
```

Expected: FAIL because module does not exist.

**Step 3: Implement server client**

Mirror `organization-units.ts`. Export types:

```ts
export type DocumentFolder = {
	id: string;
	parentId?: string | null;
	organizationUnitId: string;
	name: string;
	description?: string | null;
	deletedAt?: string | null;
	createdAt: string;
	updatedAt: string;
	children?: DocumentFolder[];
};

export type DocumentMetadata = {
	id: string;
	folderId: string;
	organizationUnitId: string;
	originalFileName: string;
	displayName: string;
	mimeType: string;
	extension?: string | null;
	sizeBytes: number;
	sha256Hash: string;
	deletedAt?: string | null;
	createdAt: string;
	updatedAt: string;
};
```

Implement request helpers that:

- add internal API headers;
- forward cookie header;
- JSON-encode normal requests;
- forward multipart `FormData` for upload without setting manual `content-type`;
- return the raw `Response` for downloads after checking `ok`.

**Step 4: Run tests**

```bash
cd frontend
rtk pnpm test -- src/lib/server/documents.test.ts
```

Expected: PASS.

**Step 5: Commit**

```bash
rtk git add frontend/src/lib/server/documents.ts frontend/src/lib/server/documents.test.ts
rtk git commit -m "Add SvelteKit document server client"
```

## Task 12: SvelteKit API Proxy Routes

**Files:**
- Create: `frontend/src/routes/api/documents/api.server.ts`
- Create: `frontend/src/routes/api/document-folders/+server.ts`
- Create: `frontend/src/routes/api/document-folders/tree/+server.ts`
- Create: `frontend/src/routes/api/document-folders/[id]/+server.ts`
- Create: `frontend/src/routes/api/document-folders/[id]/parent/+server.ts`
- Create: `frontend/src/routes/api/document-folders/[id]/archive/+server.ts`
- Create: `frontend/src/routes/api/document-folders/[id]/contents/+server.ts`
- Create: `frontend/src/routes/api/documents/upload/+server.ts`
- Create: `frontend/src/routes/api/documents/[id]/+server.ts`
- Create: `frontend/src/routes/api/documents/[id]/archive/+server.ts`
- Create: `frontend/src/routes/api/documents/[id]/download/+server.ts`
- Create tests near these routes if the existing route-test pattern is present; otherwise cover through `src/lib/server/documents.test.ts` and page API tests.

**Step 1: Write failing proxy tests**

Add route tests that call handlers directly and assert:

- unauthenticated requests return `401`;
- JSON folder payloads are normalized;
- upload forwards `FormData`;
- download returns the upstream body and preserves `content-type` and `content-disposition`;
- backend errors map to public messages.

**Step 2: Run tests to verify failure**

```bash
cd frontend
rtk pnpm test -- src/routes/api
```

Expected: FAIL because routes do not exist.

**Step 3: Implement shared helpers**

In `api.server.ts`:

```ts
export function requireAuthenticatedUser(locals: App.Locals) {
	if (locals.user) return null;
	return json({ error: 'Authentication required' }, { status: 401 });
}

export function hasAnyPermission(permissions: string[], required: string[]) {
	return permissions.includes('system.admin') || required.some((permission) => permissions.includes(permission));
}

export function cookieHeader(request: Request) {
	return request.headers.get('cookie');
}
```

Do not require `locals.user.role === 'admin'`; repository actions are RBAC-driven.

**Step 4: Implement route handlers**

Each JSON route should:

1. call `requireAuthenticatedUser`;
2. read/validate request body or query;
3. call the matching `$lib/server/documents` function with `fetch` and `cookieHeader(request)`;
4. return `json(...)`;
5. map `DocumentApiError` through `publicErrorStatus/publicErrorMessage`.

Upload route should call `await request.formData()` and pass the form data through.

Download route should return:

```ts
const upstream = await downloadDocument(fetch, cookieHeader(request), params.id);
return new Response(upstream.body, {
	status: upstream.status,
	headers: {
		'content-type': upstream.headers.get('content-type') ?? 'application/octet-stream',
		'content-disposition': upstream.headers.get('content-disposition') ?? 'attachment'
	}
});
```

**Step 5: Run tests**

```bash
cd frontend
rtk pnpm test -- src/lib/server/documents.test.ts src/routes/api
```

Expected: PASS.

**Step 6: Commit**

```bash
rtk git add frontend/src/routes/api/documents frontend/src/routes/api/document-folders frontend/src/lib/server/documents.ts frontend/src/lib/server/documents.test.ts
rtk git commit -m "Add document API proxy routes"
```

## Task 13: Documents Page Server Data And Navigation

**Files:**
- Create: `frontend/src/routes/app/documents/+page.server.ts`
- Create: `frontend/src/routes/app/documents/page.server.test.ts`
- Modify: `frontend/src/lib/components/app-sidebar.svelte`
- Modify: `frontend/src/lib/components/site-header.svelte`
- Modify: `frontend/src/routes/layout-permissions.test.ts`
- Modify: `frontend/src/routes/layout-source.test.ts`

**Step 1: Write failing tests**

In `page.server.test.ts`:

```ts
it('exposes document permission flags from locals', () => {
	const result = load(loadEvent(['document.view', 'document.create'], 'unit-id') as never);
	expect(result).toEqual({
		canViewDocuments: true,
		canCreateDocuments: true,
		canUpdateDocuments: false,
		canDeleteDocuments: false,
		canDownloadDocuments: false,
		selectedOrganizationUnitId: 'unit-id'
	});
});
```

Extend `layout-permissions.test.ts` to assert `documentNavPermissions`, `document.view`, and `/app/documents` exist in the sidebar source.

Extend `layout-source.test.ts` or site-header tests to assert `/app/documents` maps to `Documents`.

**Step 2: Run tests to verify failure**

```bash
cd frontend
rtk pnpm test -- src/routes/app/documents/page.server.test.ts src/routes/layout-permissions.test.ts src/routes/layout-source.test.ts
```

Expected: FAIL because route/nav entries do not exist.

**Step 3: Implement page server load**

Create:

```ts
const documentPermissions = {
	view: ['system.admin', 'document.view'],
	create: ['system.admin', 'document.create'],
	update: ['system.admin', 'document.update'],
	delete: ['system.admin', 'document.delete'],
	download: ['system.admin', 'document.download']
};

function hasAny(permissions: string[], required: string[]) {
	return required.some((permission) => permissions.includes(permission));
}

export const load: PageServerLoad = ({ locals, url }) => {
	const permissions = locals.permissions ?? [];
	return {
		canViewDocuments: hasAny(permissions, documentPermissions.view),
		canCreateDocuments: hasAny(permissions, documentPermissions.create),
		canUpdateDocuments: hasAny(permissions, documentPermissions.update),
		canDeleteDocuments: hasAny(permissions, documentPermissions.delete),
		canDownloadDocuments: hasAny(permissions, documentPermissions.download),
		selectedOrganizationUnitId: url.searchParams.get('organizationUnitId')
	};
};
```

**Step 4: Update navigation**

In `app-sidebar.svelte`, add `FileTextIcon` from lucide and:

```ts
const documentNavPermissions = [
	'document.view',
	'document.create',
	'document.update',
	'document.delete',
	'document.download'
];
```

Add an Operations nav item:

```ts
...(hasAnyPermission(documentNavPermissions)
	? [{ title: 'Documents', url: '/app/documents', icon: FileTextIcon }]
	: [])
```

Update `site-header.svelte` route title map:

```ts
'/app/documents': 'Documents'
```

**Step 5: Run tests**

```bash
cd frontend
rtk pnpm test -- src/routes/app/documents/page.server.test.ts src/routes/layout-permissions.test.ts src/routes/layout-source.test.ts
```

Expected: PASS.

**Step 6: Commit**

```bash
rtk git add frontend/src/routes/app/documents/+page.server.ts frontend/src/routes/app/documents/page.server.test.ts frontend/src/lib/components/app-sidebar.svelte frontend/src/lib/components/site-header.svelte frontend/src/routes/layout-permissions.test.ts frontend/src/routes/layout-source.test.ts
rtk git commit -m "Add document repository navigation"
```

## Task 14: Frontend Browser API And Repository Utilities

**Files:**
- Create: `frontend/src/routes/app/documents/api.ts`
- Create: `frontend/src/routes/app/documents/api.test.ts`
- Create: `frontend/src/routes/app/documents/tree.ts`
- Create: `frontend/src/routes/app/documents/tree.test.ts`
- Create: `frontend/src/routes/app/documents/upload-queue.ts`
- Create: `frontend/src/routes/app/documents/upload-queue.test.ts`

**Step 1: Write failing browser API tests**

Mirror `organization-units/api.test.ts`:

```ts
it('fetches folder tree through the Svelte API wrapper', async () => {
	const fetchMock = vi.fn().mockResolvedValue(jsonResponse({ folders: [folderNode()] }, { status: 200 }));
	const result = await fetchDocumentFolderTree(fetchMock, 'unit-id');
	expect(fetchMock).toHaveBeenCalledWith('/api/document-folders/tree?organizationUnitId=unit-id', {
		method: 'GET',
		headers: undefined,
		body: undefined
	});
	expect(result.folders[0].name).toBe('Invoices');
});
```

Add tests for create/update/move/archive folders, folder contents, upload `FormData`, document rename/archive/download URL building, public error messages, and malformed payloads.

**Step 2: Write failing utility tests**

In `tree.test.ts`, test flattening, selected-folder lookup, move-target exclusion, folder/document row sorting, and empty state selection.

In `upload-queue.test.ts`, test that `filesToUploadItems(files)` creates stable per-file queue items and that status transitions preserve item IDs.

**Step 3: Run tests to verify failure**

```bash
cd frontend
rtk pnpm test -- src/routes/app/documents/api.test.ts src/routes/app/documents/tree.test.ts src/routes/app/documents/upload-queue.test.ts
```

Expected: FAIL because files do not exist.

**Step 4: Implement browser API**

Export:

```ts
export const DOCUMENT_FOLDERS_QUERY_KEY = ['document-folders'] as const;
export const DOCUMENT_FOLDER_CONTENTS_QUERY_KEY = ['document-folder-contents'] as const;
```

Implement functions:

- `fetchDocumentFolderTree(fetchFn, organizationUnitId)`
- `createDocumentFolder(fetchFn, input)`
- `updateDocumentFolder(fetchFn, variables)`
- `moveDocumentFolder(fetchFn, variables)`
- `archiveDocumentFolder(fetchFn, variables)`
- `fetchDocumentFolderContents(fetchFn, folderId)`
- `uploadDocument(fetchFn, { folderId, file })`
- `updateDocument(fetchFn, variables)`
- `archiveDocument(fetchFn, variables)`
- `documentDownloadHref(id)`

Use runtime validators and `publicApiErrorMessage` as in `organization-units/api.ts`.

**Step 5: Implement utilities**

`tree.ts` should export types and functions:

- `DocumentFolderNode`
- `FlatDocumentFolderNode`
- `RepositoryDocument`
- `RepositoryRow`
- `flattenFolderTree`
- `findFolder`
- `selectInitialFolder`
- `collectFolderMoveTargets`
- `repositoryRows(folders, documents)`

`upload-queue.ts` should export:

```ts
export type UploadQueueItem = {
	id: string;
	file: File;
	status: 'queued' | 'uploading' | 'uploaded' | 'failed';
	error?: string;
	documentId?: string;
};
```

Use `crypto.randomUUID()` when available and a deterministic fallback for tests.

**Step 6: Run tests**

```bash
cd frontend
rtk pnpm test -- src/routes/app/documents/api.test.ts src/routes/app/documents/tree.test.ts src/routes/app/documents/upload-queue.test.ts
```

Expected: PASS.

**Step 7: Commit**

```bash
rtk git add frontend/src/routes/app/documents/api.ts frontend/src/routes/app/documents/api.test.ts frontend/src/routes/app/documents/tree.ts frontend/src/routes/app/documents/tree.test.ts frontend/src/routes/app/documents/upload-queue.ts frontend/src/routes/app/documents/upload-queue.test.ts
rtk git commit -m "Add document repository frontend utilities"
```

## Task 15: Documents Page UI

**Files:**
- Create: `frontend/src/routes/app/documents/+page.svelte`
- Create: `frontend/src/routes/app/documents/folder-tree.svelte`
- Create: `frontend/src/routes/app/documents/repository-table.svelte`
- Create: `frontend/src/routes/app/documents/upload-panel.svelte`
- Create or modify tests as needed under `frontend/src/routes/app/documents/`

**Step 1: Write failing source/component tests**

Add a source test similar to Organization Units:

```ts
it('uses TanStack query and mutations against the document API wrapper', () => {
	const source = readFileSync(new URL('./+page.svelte', import.meta.url), 'utf8');
	expect(source).toContain("import { createMutation, createQuery, useQueryClient }");
	expect(source).toContain('fetchDocumentFolderTree');
	expect(source).toContain('fetchDocumentFolderContents');
	expect(source).toContain('uploadDocument');
	expect(source).not.toContain('$lib/server/documents');
});
```

Add tests for component source exports if the repo does not have Svelte component testing set up.

**Step 2: Run tests to verify failure**

```bash
cd frontend
rtk pnpm test -- src/routes/app/documents
```

Expected: FAIL because UI files are missing.

**Step 3: Implement `folder-tree.svelte`**

Props:

```ts
type Props = {
	folders: FlatDocumentFolderNode[];
	selectedId: string | null;
	onSelect: (id: string) => void;
};
```

Render stable-height buttons with folder icon, indentation from `depth`, selected state, and empty state. Use lucide `FolderIcon`.

**Step 4: Implement `repository-table.svelte`**

Props:

```ts
type Props = {
	rows: RepositoryRow[];
	canUpdate: boolean;
	canDelete: boolean;
	canDownload: boolean;
	isPending: boolean;
	onOpenFolder: (id: string) => void;
	onRenameFolder: (id: string, name: string) => Promise<void>;
	onArchiveFolder: (id: string) => Promise<void>;
	onRenameDocument: (id: string, displayName: string) => Promise<void>;
	onArchiveDocument: (id: string) => Promise<void>;
};
```

Use compact rows, icon buttons with lucide icons, `confirm(...)` for archive, and `<a href={documentDownloadHref(row.id)}>` for downloads when permitted.

**Step 5: Implement `upload-panel.svelte`**

Props:

```ts
type Props = {
	canCreate: boolean;
	selectedFolderId: string | null;
	queue: UploadQueueItem[];
	isUploading: boolean;
	onFilesSelected: (files: FileList) => void;
	onUpload: () => Promise<void>;
};
```

Use a file input with `multiple`, per-file queue rows, and no in-app instructions beyond labels/status.

**Step 6: Implement `+page.svelte`**

Use:

- `fetchDocumentFolderTree(fetch, selectedOrganizationUnitId)`
- `fetchDocumentFolderContents(fetch, selectedFolderId)`
- mutations for create folder, update folder, move folder, archive folder, upload document, update document, archive document
- `queryClient.invalidateQueries` after each successful mutation

Layout:

- header with `FileTextIcon`, title `Documents`, active count;
- organization-unit selector initially can be a text/select placeholder only if accessible units are not loaded yet; prefer loading existing Organization Unit tree through the frontend organization-unit API when `canViewDocuments` is true;
- left folder tree card;
- right contents table and upload panel.

If no folder exists and `canCreateDocuments`, show an inline root-folder create form. If no permission, show a forbidden state.

**Step 7: Run tests**

```bash
cd frontend
rtk pnpm test -- src/routes/app/documents
rtk pnpm check
```

Expected: PASS.

**Step 8: Commit**

```bash
rtk git add frontend/src/routes/app/documents
rtk git commit -m "Add document repository UI"
```

## Task 16: Final Full Verification

**Files:**
- Modify only files needed to fix verification failures.
- Update: `README.md` if setup commands or env documentation changed beyond `.env.example`.

**Step 1: Backend verification**

```bash
cd go-server
rtk go test ./...
rtk atlas migrate validate --dir file://migrations
```

Expected: PASS.

**Step 2: Frontend verification**

```bash
cd frontend
rtk pnpm test
rtk pnpm check
rtk pnpm build
```

Expected: PASS.

**Step 3: Manual smoke test**

Start services if local DB/env are available:

```bash
cd go-server
rtk go run ./cmd/api
```

In another terminal:

```bash
cd frontend
rtk pnpm dev
```

Visit `/app/documents` as a seeded/system admin user. Verify:

- Documents nav is visible.
- A root folder can be created.
- A child folder can be created and selected.
- A small file uploads.
- The uploaded document appears in contents.
- Download returns the same file bytes.
- Archive hides the document or folder from active views.

If local env is unavailable, record that manual smoke testing was skipped and why.

**Step 4: Update README if needed**

If `SYNCRA_DOCUMENT_STORAGE_ROOT` is required by `config.Load`, add it to README backend setup:

```md
- `SYNCRA_DOCUMENT_STORAGE_ROOT` points to a local directory outside the repo for uploaded file bytes.
```

**Step 5: Final commit**

```bash
rtk git add README.md go-server frontend
rtk git commit -m "Complete document repository MVP"
```

Skip this commit if all previous task commits already contain all changes and no verification/doc fixes were needed.

## Execution Handoff

Use `superpowers:executing-plans` to execute this plan. Work one task at a time, run the exact verification commands in each task, and commit after each task. If a task uncovers a mismatch with the approved design, stop and update the design before continuing.
