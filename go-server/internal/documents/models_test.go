package documents

import (
	"strings"
	"testing"
	"time"
	"unicode"

	"ai.ro/syncra/dms/internal/auth"
	"ai.ro/syncra/dms/internal/orgunits"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestDocumentModelsAssignIDs(t *testing.T) {
	db := sqliteMemoryDB(t)
	if err := db.AutoMigrate(&Folder{}, &Document{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	now := time.Now().UTC()
	folder := Folder{
		Name:               "Invoices",
		OrganizationUnitID: uuid.NewString(),
		CreatedByUserID:    uuid.NewString(),
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	if err := db.Create(&folder).Error; err != nil {
		t.Fatalf("create folder: %v", err)
	}
	if folder.ID == "" {
		t.Fatal("folder ID was empty")
	}

	document := Document{
		FolderID:           folder.ID,
		OrganizationUnitID: folder.OrganizationUnitID,
		OriginalFileName:   "invoice.pdf",
		DisplayName:        "Invoice",
		MimeType:           "application/pdf",
		SizeBytes:          42,
		SHA256Hash:         strings.Repeat("a", 64),
		StorageKey:         "org/invoice.pdf",
		CreatedByUserID:    folder.CreatedByUserID,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	if err := db.Create(&document).Error; err != nil {
		t.Fatalf("create document: %v", err)
	}
	if document.ID == "" {
		t.Fatal("document ID was empty")
	}
}

func TestDocumentModelsUseExpectedTableNames(t *testing.T) {
	if got := (Folder{}).TableName(); got != "document_folders" {
		t.Fatalf("Folder.TableName() = %q, want document_folders", got)
	}
	if got := (Document{}).TableName(); got != "documents" {
		t.Fatalf("Document.TableName() = %q, want documents", got)
	}
}

func TestFolderModelCreatesRootAndChildNameUniqueIndexes(t *testing.T) {
	db := sqliteMemoryDB(t)
	if err := db.AutoMigrate(&Folder{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	assertIndexSQLContains(t, db, "idx_document_folders_root_name_unique",
		"CREATE UNIQUE INDEX",
		"`organization_unit_id`,`name`",
		"parent_id IS NULL AND deleted_at IS NULL",
	)
	assertIndexSQLContains(t, db, "idx_document_folders_child_name_unique",
		"CREATE UNIQUE INDEX",
		"`organization_unit_id`,`parent_id`,`name`",
		"parent_id IS NOT NULL AND deleted_at IS NULL",
	)
}

func TestDocumentModelsCreateOrganizationIntegrityConstraints(t *testing.T) {
	db := sqliteMemoryDB(t)
	if err := db.AutoMigrate(&Folder{}, &Document{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	assertIndexSQLContains(t, db, "idx_document_folders_id_organization_unit_unique",
		"CREATE UNIQUE INDEX",
		"`id`,`organization_unit_id`",
	)
	assertTableSQLContains(t, db, "document_folders",
		"FOREIGN KEY (`parent_id`,`organization_unit_id`) REFERENCES `document_folders`(`id`,`organization_unit_id`)",
	)
	assertTableSQLContains(t, db, "documents",
		"FOREIGN KEY (`folder_id`,`organization_unit_id`) REFERENCES `document_folders`(`id`,`organization_unit_id`)",
	)
}

func TestDocumentModelsCreateMetadataChecks(t *testing.T) {
	db := sqliteMemoryDB(t)
	if err := db.AutoMigrate(&Document{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	assertTableSQLContains(t, db, "documents",
		"CONSTRAINT `chk_documents_size_bytes_non_negative` CHECK (size_bytes >= 0)",
		"CONSTRAINT `chk_documents_sha256_hash_lower_hex` CHECK",
		"length(sha256_hash) = 64",
		"replace(sha256_hash,'0','')",
	)
}

func TestDocumentModelsRejectCrossOrganizationRelationships(t *testing.T) {
	db := sqliteMemoryDBWithForeignKeys(t)
	if err := db.AutoMigrate(&orgunits.Unit{}, &auth.User{}, &Folder{}, &Document{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	orgA := createTestOrganizationUnit(t, db, "Finance")
	orgB := createTestOrganizationUnit(t, db, "Legal")
	user := createTestUser(t, db)
	now := time.Now().UTC()

	parent := Folder{
		Name:               "Invoices",
		OrganizationUnitID: orgA.ID,
		CreatedByUserID:    user.ID,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	if err := db.Create(&parent).Error; err != nil {
		t.Fatalf("create parent folder: %v", err)
	}

	sameOrgChild := Folder{
		ParentID:           &parent.ID,
		Name:               "Same org child",
		OrganizationUnitID: orgA.ID,
		CreatedByUserID:    user.ID,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	if err := db.Create(&sameOrgChild).Error; err != nil {
		t.Fatalf("create same-organization child folder: %v", err)
	}

	crossOrgChild := Folder{
		ParentID:           &parent.ID,
		Name:               "Cross org child",
		OrganizationUnitID: orgB.ID,
		CreatedByUserID:    user.ID,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	if err := db.Create(&crossOrgChild).Error; err == nil {
		t.Fatal("cross-organization child folder was accepted")
	}

	sameOrgDocument := validDocument(parent.ID, orgA.ID, user.ID)
	if err := db.Create(&sameOrgDocument).Error; err != nil {
		t.Fatalf("create same-organization document: %v", err)
	}

	crossOrgDocument := validDocument(parent.ID, orgB.ID, user.ID)
	crossOrgDocument.SHA256Hash = strings.Repeat("b", 64)
	if err := db.Create(&crossOrgDocument).Error; err == nil {
		t.Fatal("cross-organization document was accepted")
	}
}

func TestDocumentModelsRejectInvalidMetadata(t *testing.T) {
	db := sqliteMemoryDBWithForeignKeys(t)
	if err := db.AutoMigrate(&orgunits.Unit{}, &auth.User{}, &Folder{}, &Document{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	org := createTestOrganizationUnit(t, db, "Finance")
	user := createTestUser(t, db)
	now := time.Now().UTC()

	folder := Folder{
		Name:               "Invoices",
		OrganizationUnitID: org.ID,
		CreatedByUserID:    user.ID,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	if err := db.Create(&folder).Error; err != nil {
		t.Fatalf("create folder: %v", err)
	}

	tests := []struct {
		name       string
		sizeBytes  int64
		sha256Hash string
	}{
		{name: "negative size", sizeBytes: -1, sha256Hash: strings.Repeat("a", 64)},
		{name: "short hash", sizeBytes: 1, sha256Hash: strings.Repeat("a", 63)},
		{name: "uppercase hash", sizeBytes: 1, sha256Hash: strings.Repeat("A", 64)},
		{name: "non hex hash", sizeBytes: 1, sha256Hash: strings.Repeat("g", 64)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			document := validDocument(folder.ID, org.ID, user.ID)
			document.SizeBytes = tt.sizeBytes
			document.SHA256Hash = tt.sha256Hash
			if err := db.Create(&document).Error; err == nil {
				t.Fatal("invalid document metadata was accepted")
			}
		})
	}
}

func TestNormalizeFolderName(t *testing.T) {
	got, err := NormalizeFolderName(" Invoices ")
	if err != nil {
		t.Fatalf("NormalizeFolderName() error = %v", err)
	}
	if got != "Invoices" {
		t.Fatalf("NormalizeFolderName() = %q, want Invoices", got)
	}

	if got, err := NormalizeFolderName(" \t\n "); err == nil {
		t.Fatalf("NormalizeFolderName() error = nil, name = %q", got)
	}

	if got, err := NormalizeFolderName(strings.Repeat("\u0103", MaxFolderNameCharacters+1)); err == nil {
		t.Fatalf("NormalizeFolderName() error = nil, name = %q", got)
	}
}

func TestNormalizeDescription(t *testing.T) {
	got := NormalizeDescription(" Invoice folder ")
	if got == nil || *got != "Invoice folder" {
		t.Fatalf("NormalizeDescription() = %v, want Invoice folder", got)
	}

	if got := NormalizeDescription(" \t\n "); got != nil {
		t.Fatalf("NormalizeDescription() = %v, want nil", *got)
	}
}

func TestNormalizeDisplayName(t *testing.T) {
	got, err := NormalizeDisplayName(" Invoice.pdf ")
	if err != nil {
		t.Fatalf("NormalizeDisplayName() error = %v", err)
	}
	if got != "Invoice.pdf" {
		t.Fatalf("NormalizeDisplayName() = %q, want Invoice.pdf", got)
	}

	if got, err := NormalizeDisplayName(" \t\n "); err == nil {
		t.Fatalf("NormalizeDisplayName() error = nil, name = %q", got)
	}

	if got, err := NormalizeDisplayName(strings.Repeat("\u0103", MaxDocumentDisplayNameCharacters+1)); err == nil {
		t.Fatalf("NormalizeDisplayName() error = nil, name = %q", got)
	}
}

func TestSafeOriginalFileName(t *testing.T) {
	if got := SafeOriginalFileName("../uploads/invoice.pdf"); got != "invoice.pdf" {
		t.Fatalf("SafeOriginalFileName() = %q, want invoice.pdf", got)
	}
	if got := SafeOriginalFileName("C:\\fakepath\\invoice.pdf"); got != "invoice.pdf" {
		t.Fatalf("SafeOriginalFileName() = %q, want invoice.pdf", got)
	}
	if got := SafeOriginalFileName("..\\invoice.pdf"); got != "invoice.pdf" {
		t.Fatalf("SafeOriginalFileName() = %q, want invoice.pdf", got)
	}
	if got := SafeOriginalFileName("invoice\n\t.pdf"); strings.IndexFunc(got, unicode.IsControl) >= 0 {
		t.Fatalf("SafeOriginalFileName() = %q, want no control runes", got)
	}
	if got := SafeOriginalFileName(" \t\n "); got != "upload" {
		t.Fatalf("SafeOriginalFileName() = %q, want upload", got)
	}
	if got := SafeOriginalFileName(strings.Repeat("\u0103", 256)); len([]rune(got)) != 255 {
		t.Fatalf("SafeOriginalFileName() length = %d, want 255", len([]rune(got)))
	}
}

func sqliteMemoryDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	return db
}

func sqliteMemoryDBWithForeignKeys(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+uuid.NewString()+"?mode=memory&cache=shared&_foreign_keys=on"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		t.Fatalf("enable foreign keys: %v", err)
	}
	return db
}

func createTestOrganizationUnit(t *testing.T, db *gorm.DB, name string) orgunits.Unit {
	t.Helper()
	now := time.Now().UTC()
	unit := orgunits.Unit{
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := db.Create(&unit).Error; err != nil {
		t.Fatalf("create organization unit: %v", err)
	}
	return unit
}

func createTestUser(t *testing.T, db *gorm.DB) auth.User {
	t.Helper()
	now := time.Now().UTC()
	user := auth.User{
		Name:          "Document Tester",
		Email:         uuid.NewString() + "@example.test",
		EmailVerified: true,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	return user
}

func validDocument(folderID string, organizationUnitID string, userID string) Document {
	now := time.Now().UTC()
	return Document{
		FolderID:           folderID,
		OrganizationUnitID: organizationUnitID,
		OriginalFileName:   "invoice.pdf",
		DisplayName:        "Invoice",
		MimeType:           "application/pdf",
		SizeBytes:          42,
		SHA256Hash:         strings.Repeat("a", 64),
		StorageKey:         uuid.NewString() + "/invoice.pdf",
		CreatedByUserID:    userID,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
}

func assertIndexSQLContains(t *testing.T, db *gorm.DB, indexName string, wantParts ...string) {
	t.Helper()
	if !db.Migrator().HasIndex(&Folder{}, indexName) {
		t.Fatalf("missing index %s", indexName)
	}

	var sql string
	if err := db.Raw("SELECT sql FROM sqlite_master WHERE type = 'index' AND name = ?", indexName).Scan(&sql).Error; err != nil {
		t.Fatalf("load index %s SQL: %v", indexName, err)
	}
	for _, want := range wantParts {
		if !strings.Contains(sql, want) {
			t.Fatalf("index %s SQL = %q, want to contain %q", indexName, sql, want)
		}
	}
}

func assertTableSQLContains(t *testing.T, db *gorm.DB, tableName string, wantParts ...string) {
	t.Helper()
	var sql string
	if err := db.Raw("SELECT sql FROM sqlite_master WHERE type = 'table' AND name = ?", tableName).Scan(&sql).Error; err != nil {
		t.Fatalf("load table %s SQL: %v", tableName, err)
	}
	if sql == "" {
		t.Fatalf("table %s SQL was empty", tableName)
	}
	for _, want := range wantParts {
		if !strings.Contains(sql, want) {
			t.Fatalf("table %s SQL = %q, want to contain %q", tableName, sql, want)
		}
	}
}
