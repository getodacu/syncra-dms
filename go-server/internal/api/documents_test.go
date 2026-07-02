package api

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"ai.ro/syncra/dms/internal/documents"
	"gorm.io/gorm"
)

func TestDocumentUploadStoresMetadataAndFile(t *testing.T) {
	storageRoot := t.TempDir()
	router, db := newAuthTestRouterWithOptions(t, RouterOptions{DocumentStorageRoot: storageRoot})
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
	assertResponseOmitsStorageKey(t, response)

	var stored documents.Document
	if err := db.First(&stored, "id = ?", body.ID).Error; err != nil {
		t.Fatalf("load stored document: %v", err)
	}
	if stored.StorageKey == "" {
		t.Fatal("stored document storage key was empty")
	}
	if _, err := os.Stat(filepath.Join(storageRoot, filepath.FromSlash(stored.StorageKey))); err != nil {
		t.Fatalf("stored file stat: %v", err)
	}

	var count int64
	if err := db.Model(&documents.Document{}).Where("folder_id = ?", folderID).Count(&count).Error; err != nil {
		t.Fatalf("count documents: %v", err)
	}
	if count != 1 {
		t.Fatalf("document count = %d", count)
	}
}

func TestDocumentUploadRejectsMissingFile(t *testing.T) {
	router, db := newAuthTestRouterWithOptions(t, RouterOptions{DocumentStorageRoot: t.TempDir()})
	token := loginSeededAdmin(t, router, db, "admin@example.com")
	unitID := createUnitViaAPI(t, router, token, `{"name":"Finance"}`)
	folderID := createFolderViaAPI(t, router, token, `{"organizationUnitId":"`+unitID+`","name":"Invoices"}`)

	response := uploadDocumentForm(t, router, token, func(writer *multipart.Writer) {
		writeMultipartField(t, writer, "folderId", folderID)
	})
	if response.Code != http.StatusBadRequest {
		t.Fatalf("missing file status = %d body=%s, want bad request", response.Code, response.Body.String())
	}
	assertDocumentCount(t, db, folderID, 0)
}

func TestDocumentUploadRejectsMissingFolder(t *testing.T) {
	router, db := newAuthTestRouterWithOptions(t, RouterOptions{DocumentStorageRoot: t.TempDir()})
	token := loginSeededAdmin(t, router, db, "admin@example.com")

	response := uploadDocument(t, router, token, "00000000-0000-0000-0000-000000000123", "invoice.pdf", []byte("%PDF test"))
	if response.Code != http.StatusNotFound {
		t.Fatalf("missing folder status = %d body=%s, want not found", response.Code, response.Body.String())
	}
}

func TestDocumentUploadRejectsArchivedFolder(t *testing.T) {
	router, db := newAuthTestRouterWithOptions(t, RouterOptions{DocumentStorageRoot: t.TempDir()})
	token := loginSeededAdmin(t, router, db, "admin@example.com")
	unitID := createUnitViaAPI(t, router, token, `{"name":"Finance"}`)
	folderID := createFolderViaAPI(t, router, token, `{"organizationUnitId":"`+unitID+`","name":"Invoices"}`)
	archivedAt := time.Now().UTC()
	if err := db.Model(&documents.Folder{}).Where("id = ?", folderID).Update("deleted_at", archivedAt).Error; err != nil {
		t.Fatalf("archive folder: %v", err)
	}

	response := uploadDocument(t, router, token, folderID, "invoice.pdf", []byte("%PDF test"))
	if response.Code != http.StatusNotFound {
		t.Fatalf("archived folder status = %d body=%s, want not found", response.Code, response.Body.String())
	}
	assertDocumentCount(t, db, folderID, 0)
}

func TestDocumentUploadRequiresDocumentCreateForFolderOrganizationUnit(t *testing.T) {
	router, db := newAuthTestRouterWithOptions(t, RouterOptions{DocumentStorageRoot: t.TempDir()})
	adminToken := loginSeededAdmin(t, router, db, "admin@example.com")
	unitID := createUnitViaAPI(t, router, adminToken, `{"name":"Finance"}`)
	folderID := createFolderViaAPI(t, router, adminToken, `{"organizationUnitId":"`+unitID+`","name":"Invoices"}`)
	user := createVerifiedUser(t, db, "viewer@example.com", "password123")
	token := loginUser(t, router, user.Email, "password123")

	response := uploadDocument(t, router, token, folderID, "invoice.pdf", []byte("%PDF test"))
	if response.Code != http.StatusNotFound {
		t.Fatalf("missing document.create status = %d body=%s, want not found", response.Code, response.Body.String())
	}
	assertDocumentCount(t, db, folderID, 0)
}

func TestDocumentUploadRejectsDuplicateHashInSameFolder(t *testing.T) {
	storageRoot := t.TempDir()
	router, db := newAuthTestRouterWithOptions(t, RouterOptions{DocumentStorageRoot: storageRoot})
	token := loginSeededAdmin(t, router, db, "admin@example.com")
	unitID := createUnitViaAPI(t, router, token, `{"name":"Finance"}`)
	folderID := createFolderViaAPI(t, router, token, `{"organizationUnitId":"`+unitID+`","name":"Invoices"}`)

	first := uploadDocument(t, router, token, folderID, "invoice.pdf", []byte("%PDF duplicate"))
	if first.Code != http.StatusCreated {
		t.Fatalf("first upload status = %d body=%s", first.Code, first.Body.String())
	}
	duplicate := uploadDocument(t, router, token, folderID, "copy.pdf", []byte("%PDF duplicate"))
	if duplicate.Code != http.StatusConflict {
		t.Fatalf("duplicate upload status = %d body=%s, want conflict", duplicate.Code, duplicate.Body.String())
	}
	assertDocumentCount(t, db, folderID, 1)
	if storedFiles := countStoredDocumentFiles(t, storageRoot); storedFiles != 1 {
		t.Fatalf("stored document files = %d, want 1 after duplicate cleanup", storedFiles)
	}
}

func TestDocumentUploadRejectsMaxUploadSize(t *testing.T) {
	router, db := newAuthTestRouterWithOptions(t, RouterOptions{
		DocumentStorageRoot:    t.TempDir(),
		DocumentMaxUploadBytes: 4,
	})
	token := loginSeededAdmin(t, router, db, "admin@example.com")
	unitID := createUnitViaAPI(t, router, token, `{"name":"Finance"}`)
	folderID := createFolderViaAPI(t, router, token, `{"organizationUnitId":"`+unitID+`","name":"Invoices"}`)

	response := uploadDocument(t, router, token, folderID, "invoice.txt", []byte("hello"))
	if response.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("oversize upload status = %d body=%s, want request entity too large", response.Code, response.Body.String())
	}
	assertDocumentCount(t, db, folderID, 0)
}

func TestDocumentUploadRejectsUnsupportedMIMEType(t *testing.T) {
	router, db := newAuthTestRouterWithOptions(t, RouterOptions{
		DocumentStorageRoot:      t.TempDir(),
		DocumentAllowedMIMETypes: []string{"application/pdf"},
	})
	token := loginSeededAdmin(t, router, db, "admin@example.com")
	unitID := createUnitViaAPI(t, router, token, `{"name":"Finance"}`)
	folderID := createFolderViaAPI(t, router, token, `{"organizationUnitId":"`+unitID+`","name":"Invoices"}`)

	response := uploadDocument(t, router, token, folderID, "invoice.txt", []byte("hello"))
	if response.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("unsupported MIME status = %d body=%s, want unsupported media type", response.Code, response.Body.String())
	}
	assertDocumentCount(t, db, folderID, 0)
}

func uploadDocument(t *testing.T, router http.Handler, token string, folderID string, fileName string, content []byte) *httptest.ResponseRecorder {
	t.Helper()
	return uploadDocumentForm(t, router, token, func(writer *multipart.Writer) {
		writeMultipartField(t, writer, "folderId", folderID)
		part, err := writer.CreateFormFile("file", fileName)
		if err != nil {
			t.Fatalf("create multipart file: %v", err)
		}
		if _, err := part.Write(content); err != nil {
			t.Fatalf("write multipart file: %v", err)
		}
	})
}

func uploadDocumentForm(t *testing.T, router http.Handler, token string, writeForm func(*multipart.Writer)) *httptest.ResponseRecorder {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	writeForm(writer)
	if err := writer.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/documents/upload", &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	for key, value := range authCookieHeaders(token) {
		request.Header.Set(key, value)
	}
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

func writeMultipartField(t *testing.T, writer *multipart.Writer, name string, value string) {
	t.Helper()
	if err := writer.WriteField(name, value); err != nil {
		t.Fatalf("write multipart field %s: %v", name, err)
	}
}

func assertDocumentCount(t *testing.T, db *gorm.DB, folderID string, want int64) {
	t.Helper()
	var count int64
	if err := db.Model(&documents.Document{}).Where("folder_id = ?", folderID).Count(&count).Error; err != nil {
		t.Fatalf("count documents: %v", err)
	}
	if count != want {
		t.Fatalf("document count = %d, want %d", count, want)
	}
}

func assertResponseOmitsStorageKey(t *testing.T, response *httptest.ResponseRecorder) {
	t.Helper()
	var raw map[string]any
	decodeJSON(t, response, &raw)
	if _, ok := raw["storageKey"]; ok {
		t.Fatalf("response leaked storageKey: %s", response.Body.String())
	}
}

func countStoredDocumentFiles(t *testing.T, storageRoot string) int {
	t.Helper()
	documentsRoot := filepath.Join(storageRoot, "documents")
	if _, err := os.Stat(documentsRoot); os.IsNotExist(err) {
		return 0
	} else if err != nil {
		t.Fatalf("stat documents storage root: %v", err)
	}

	count := 0
	if err := filepath.WalkDir(documentsRoot, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Type().IsRegular() {
			count++
		}
		return nil
	}); err != nil {
		t.Fatalf("walk documents storage: %v", err)
	}
	return count
}
