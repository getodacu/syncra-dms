package ocr

import (
	"context"
	"strings"
	"testing"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"ai.ro/syncra/internal/auth"
	"ai.ro/syncra/internal/testsupport"
)

var ocrCollectionTestGroup *testsupport.PostgresGroup

func TestOCRCollections(t *testing.T) {
	ocrCollectionTestGroup = testsupport.OpenPostgresGroup(t, &auth.User{}, &ExtractionSchema{}, &OCRDocument{}, &OCRJob{}, &Collection{}, &CollectionSchema{}, &CollectionDocument{})
	defer func() { ocrCollectionTestGroup = nil }()

	for _, tt := range []struct {
		name string
		fn   func(*testing.T)
	}{
		{name: "LinkDocumentToMatchingCollectionsLinksSavedSchemaDocument", fn: testLinkDocumentToMatchingCollectionsLinksSavedSchemaDocument},
		{name: "LinkDocumentToMatchingCollectionsSkipsDocumentWithoutSchema", fn: testLinkDocumentToMatchingCollectionsSkipsDocumentWithoutSchema},
		{name: "LinkDocumentToMatchingCollectionsSkipsDeletedDocument", fn: testLinkDocumentToMatchingCollectionsSkipsDeletedDocument},
		{name: "LinkDocumentToMatchingCollectionsSkipsDocumentWithoutUser", fn: testLinkDocumentToMatchingCollectionsSkipsDocumentWithoutUser},
		{name: "LinkDocumentToMatchingCollectionsSkipsCollectionsOwnedByDifferentUser", fn: testLinkDocumentToMatchingCollectionsSkipsCollectionsOwnedByDifferentUser},
	} {
		t.Run(tt.name, tt.fn)
	}
}

func ocrCollectionTx(t *testing.T) *gorm.DB {
	t.Helper()
	if ocrCollectionTestGroup != nil {
		return ocrCollectionTestGroup.Tx(t)
	}
	return testsupport.OpenPostgresTx(t, &auth.User{}, &ExtractionSchema{}, &OCRDocument{}, &OCRJob{}, &Collection{}, &CollectionSchema{}, &CollectionDocument{})
}

func testLinkDocumentToMatchingCollectionsLinksSavedSchemaDocument(t *testing.T) {
	db := ocrCollectionTx(t)
	user := createExecutorTestUser(t, db)
	schema := ExtractionSchema{UserID: &user.ID, Name: "Invoice", SchemaJSON: datatypes.JSON([]byte(`{"type":"object"}`)), Strict: true}
	if err := db.Create(&schema).Error; err != nil {
		t.Fatalf("create schema: %v", err)
	}
	collection := Collection{UserID: user.ID, Name: "Invoices"}
	if err := db.Create(&collection).Error; err != nil {
		t.Fatalf("create collection: %v", err)
	}
	if err := db.Create(&CollectionSchema{CollectionID: collection.ID, SchemaID: schema.ID}).Error; err != nil {
		t.Fatalf("create collection schema: %v", err)
	}
	doc := OCRDocument{
		UserID: &user.ID, OriginalFilename: "scan.png", MimeType: "image/png",
		FileSize: 10, DocumentHash: strings.Repeat("b", 64), SchemaID: &schema.ID,
		Markdown: "# OCR", RawResponseJSON: datatypes.JSON([]byte(`{"pages":[]}`)),
	}
	if err := db.Create(&doc).Error; err != nil {
		t.Fatalf("create document: %v", err)
	}

	if err := LinkDocumentToMatchingCollections(context.Background(), db, doc); err != nil {
		t.Fatalf("link document: %v", err)
	}
	if err := LinkDocumentToMatchingCollections(context.Background(), db, doc); err != nil {
		t.Fatalf("link document twice: %v", err)
	}

	var count int64
	if err := db.Model(&CollectionDocument{}).Where("collection_id = ? AND document_id = ?", collection.ID, doc.ID).Count(&count).Error; err != nil {
		t.Fatalf("count links: %v", err)
	}
	if count != 1 {
		t.Fatalf("collection document count = %d, want 1", count)
	}
}

func testLinkDocumentToMatchingCollectionsSkipsDocumentWithoutSchema(t *testing.T) {
	db := ocrCollectionTx(t)
	user := createExecutorTestUser(t, db)
	collection := Collection{UserID: user.ID, Name: "Invoices"}
	if err := db.Create(&collection).Error; err != nil {
		t.Fatalf("create collection: %v", err)
	}
	doc := OCRDocument{
		UserID: &user.ID, OriginalFilename: "scan.png", MimeType: "image/png",
		FileSize: 10, DocumentHash: strings.Repeat("c", 64),
		Markdown: "# OCR", RawResponseJSON: datatypes.JSON([]byte(`{"pages":[]}`)),
	}
	if err := db.Create(&doc).Error; err != nil {
		t.Fatalf("create document: %v", err)
	}

	if err := LinkDocumentToMatchingCollections(context.Background(), db, doc); err != nil {
		t.Fatalf("link document: %v", err)
	}

	assertCollectionDocumentCount(t, db, "document_id = ?", 0, doc.ID)
}

func testLinkDocumentToMatchingCollectionsSkipsDeletedDocument(t *testing.T) {
	db := ocrCollectionTx(t)
	user := createExecutorTestUser(t, db)
	schema := ExtractionSchema{UserID: &user.ID, Name: "Invoice", SchemaJSON: datatypes.JSON([]byte(`{"type":"object"}`)), Strict: true}
	if err := db.Create(&schema).Error; err != nil {
		t.Fatalf("create schema: %v", err)
	}
	collection := Collection{UserID: user.ID, Name: "Invoices"}
	if err := db.Create(&collection).Error; err != nil {
		t.Fatalf("create collection: %v", err)
	}
	if err := db.Create(&CollectionSchema{CollectionID: collection.ID, SchemaID: schema.ID}).Error; err != nil {
		t.Fatalf("create collection schema: %v", err)
	}
	doc := OCRDocument{
		UserID: &user.ID, OriginalFilename: "scan.png", MimeType: "image/png",
		FileSize: 10, DocumentHash: strings.Repeat("f", 64), SchemaID: &schema.ID,
		Markdown: "# OCR", RawResponseJSON: datatypes.JSON([]byte(`{"pages":[]}`)),
	}
	if err := db.Create(&doc).Error; err != nil {
		t.Fatalf("create document: %v", err)
	}
	if err := db.Delete(&doc).Error; err != nil {
		t.Fatalf("delete document: %v", err)
	}

	if err := LinkDocumentToMatchingCollections(context.Background(), db, doc); err != nil {
		t.Fatalf("link document: %v", err)
	}

	assertCollectionDocumentCount(t, db, "collection_id = ? AND document_id = ?", 0, collection.ID, doc.ID)
}

func testLinkDocumentToMatchingCollectionsSkipsDocumentWithoutUser(t *testing.T) {
	db := ocrCollectionTx(t)
	user := createExecutorTestUser(t, db)
	schema := ExtractionSchema{UserID: &user.ID, Name: "Invoice", SchemaJSON: datatypes.JSON([]byte(`{"type":"object"}`)), Strict: true}
	if err := db.Create(&schema).Error; err != nil {
		t.Fatalf("create schema: %v", err)
	}
	collection := Collection{UserID: user.ID, Name: "Invoices"}
	if err := db.Create(&collection).Error; err != nil {
		t.Fatalf("create collection: %v", err)
	}
	if err := db.Create(&CollectionSchema{CollectionID: collection.ID, SchemaID: schema.ID}).Error; err != nil {
		t.Fatalf("create collection schema: %v", err)
	}
	doc := OCRDocument{
		OriginalFilename: "scan.png", MimeType: "image/png",
		FileSize: 10, DocumentHash: strings.Repeat("d", 64), SchemaID: &schema.ID,
		Markdown: "# OCR", RawResponseJSON: datatypes.JSON([]byte(`{"pages":[]}`)),
	}
	if err := db.Create(&doc).Error; err != nil {
		t.Fatalf("create document: %v", err)
	}

	if err := LinkDocumentToMatchingCollections(context.Background(), db, doc); err != nil {
		t.Fatalf("link document: %v", err)
	}

	assertCollectionDocumentCount(t, db, "document_id = ?", 0, doc.ID)
}

func testLinkDocumentToMatchingCollectionsSkipsCollectionsOwnedByDifferentUser(t *testing.T) {
	db := ocrCollectionTx(t)
	user := createExecutorTestUser(t, db)
	otherUser := createExecutorTestUser(t, db)
	schema := ExtractionSchema{UserID: &user.ID, Name: "Invoice", SchemaJSON: datatypes.JSON([]byte(`{"type":"object"}`)), Strict: true}
	if err := db.Create(&schema).Error; err != nil {
		t.Fatalf("create schema: %v", err)
	}
	collection := Collection{UserID: otherUser.ID, Name: "Other invoices"}
	if err := db.Create(&collection).Error; err != nil {
		t.Fatalf("create collection: %v", err)
	}
	if err := db.Create(&CollectionSchema{CollectionID: collection.ID, SchemaID: schema.ID}).Error; err != nil {
		t.Fatalf("create collection schema: %v", err)
	}
	doc := OCRDocument{
		UserID: &user.ID, OriginalFilename: "scan.png", MimeType: "image/png",
		FileSize: 10, DocumentHash: strings.Repeat("e", 64), SchemaID: &schema.ID,
		Markdown: "# OCR", RawResponseJSON: datatypes.JSON([]byte(`{"pages":[]}`)),
	}
	if err := db.Create(&doc).Error; err != nil {
		t.Fatalf("create document: %v", err)
	}

	if err := LinkDocumentToMatchingCollections(context.Background(), db, doc); err != nil {
		t.Fatalf("link document: %v", err)
	}

	assertCollectionDocumentCount(t, db, "collection_id = ? AND document_id = ?", 0, collection.ID, doc.ID)
}

func assertCollectionDocumentCount(t *testing.T, db *gorm.DB, where string, want int64, args ...any) {
	t.Helper()

	var count int64
	if err := db.Model(&CollectionDocument{}).Where(where, args...).Count(&count).Error; err != nil {
		t.Fatalf("count collection documents: %v", err)
	}
	if count != want {
		t.Fatalf("collection document count = %d, want %d", count, want)
	}
}
