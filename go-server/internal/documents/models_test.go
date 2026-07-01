package documents

import (
	"strings"
	"testing"
	"time"

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
