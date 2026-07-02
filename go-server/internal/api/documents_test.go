package api

import (
	"bytes"
	"errors"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"ai.ro/syncra/dms/internal/documents"
	"ai.ro/syncra/dms/internal/orgunits"
	"ai.ro/syncra/dms/internal/rbac"
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

func TestDocumentUploadRevalidatesFolderAfterStorageSave(t *testing.T) {
	storageRoot := t.TempDir()
	router, db := newAuthTestRouterWithOptions(t, RouterOptions{DocumentStorageRoot: storageRoot})
	token := loginSeededAdmin(t, router, db, "admin@example.com")
	unitID := createUnitViaAPI(t, router, token, `{"name":"Finance"}`)
	folderID := createFolderViaAPI(t, router, token, `{"organizationUnitId":"`+unitID+`","name":"Invoices"}`)

	callbackName := documentUploadCallbackName(t, "archive_folder_before_revalidation")
	folderQueries := 0
	archived := false
	if err := db.Callback().Query().Before("gorm:query").Register(callbackName, func(tx *gorm.DB) {
		if archived || tx.Statement.Table != "document_folders" {
			return
		}
		folderQueries++
		if folderQueries != 2 {
			return
		}
		archived = true
		if err := tx.Session(&gorm.Session{NewDB: true}).Model(&documents.Folder{}).Where("id = ?", folderID).Update("deleted_at", time.Now().UTC()).Error; err != nil {
			t.Fatalf("archive folder during upload revalidation: %v", err)
		}
	}); err != nil {
		t.Fatalf("register folder archive callback: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Callback().Query().Remove(callbackName); err != nil {
			t.Fatalf("remove folder archive callback: %v", err)
		}
	})

	response := uploadDocument(t, router, token, folderID, "invoice.pdf", []byte("%PDF stale folder"))
	if response.Code != http.StatusNotFound {
		t.Fatalf("upload after folder archive status = %d body=%s, want not found", response.Code, response.Body.String())
	}
	if !archived {
		t.Fatal("folder archive callback did not run")
	}
	assertDocumentCount(t, db, folderID, 0)
	if storedFiles := countStoredDocumentFiles(t, storageRoot); storedFiles != 0 {
		t.Fatalf("stored document files = %d, want 0 after folder revalidation cleanup", storedFiles)
	}
}

func TestDocumentUploadRevalidatesOrganizationUnitAfterStorageSave(t *testing.T) {
	storageRoot := t.TempDir()
	router, db := newAuthTestRouterWithOptions(t, RouterOptions{DocumentStorageRoot: storageRoot})
	token := loginSeededAdmin(t, router, db, "admin@example.com")
	unitID := createUnitViaAPI(t, router, token, `{"name":"Finance"}`)
	folderID := createFolderViaAPI(t, router, token, `{"organizationUnitId":"`+unitID+`","name":"Invoices"}`)

	callbackName := documentUploadCallbackName(t, "archive_unit_before_revalidation")
	archived := false
	if err := db.Callback().Query().Before("gorm:query").Register(callbackName, func(tx *gorm.DB) {
		if archived || tx.Statement.Table != "organization_units" {
			return
		}
		archived = true
		if err := tx.Session(&gorm.Session{NewDB: true}).Model(&orgunits.Unit{}).Where("id = ?", unitID).Update("archived_at", time.Now().UTC()).Error; err != nil {
			t.Fatalf("archive organization unit during upload revalidation: %v", err)
		}
	}); err != nil {
		t.Fatalf("register organization unit archive callback: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Callback().Query().Remove(callbackName); err != nil {
			t.Fatalf("remove organization unit archive callback: %v", err)
		}
	})

	response := uploadDocument(t, router, token, folderID, "invoice.pdf", []byte("%PDF stale unit"))
	if response.Code != http.StatusNotFound {
		t.Fatalf("upload after organization unit archive status = %d body=%s, want not found", response.Code, response.Body.String())
	}
	if !archived {
		t.Fatal("organization unit archive callback did not run")
	}
	assertDocumentCount(t, db, folderID, 0)
	if storedFiles := countStoredDocumentFiles(t, storageRoot); storedFiles != 0 {
		t.Fatalf("stored document files = %d, want 0 after organization unit revalidation cleanup", storedFiles)
	}
}

func TestDocumentUploadReturnsServerErrorWhenDuplicateCleanupFails(t *testing.T) {
	storageRoot := t.TempDir()
	router, db := newAuthTestRouterWithOptions(t, RouterOptions{DocumentStorageRoot: storageRoot})
	token := loginSeededAdmin(t, router, db, "admin@example.com")
	unitID := createUnitViaAPI(t, router, token, `{"name":"Finance"}`)
	folderID := createFolderViaAPI(t, router, token, `{"organizationUnitId":"`+unitID+`","name":"Invoices"}`)

	first := uploadDocument(t, router, token, folderID, "invoice.pdf", []byte("%PDF duplicate cleanup"))
	if first.Code != http.StatusCreated {
		t.Fatalf("first upload status = %d body=%s", first.Code, first.Body.String())
	}
	originalFiles := storedDocumentFileSet(t, storageRoot)

	callbackName := documentUploadCallbackName(t, "corrupt_duplicate_cleanup_file")
	corrupted := false
	if err := db.Callback().Query().Before("gorm:query").Register(callbackName, func(tx *gorm.DB) {
		if corrupted || tx.Statement.Table != "documents" {
			return
		}
		corrupted = true
		replaceNewStoredFileWithSymlink(t, storageRoot, originalFiles)
	}); err != nil {
		t.Fatalf("register duplicate cleanup callback: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Callback().Query().Remove(callbackName); err != nil {
			t.Fatalf("remove duplicate cleanup callback: %v", err)
		}
	})

	duplicate := uploadDocument(t, router, token, folderID, "copy.pdf", []byte("%PDF duplicate cleanup"))
	if duplicate.Code != http.StatusInternalServerError {
		t.Fatalf("duplicate cleanup failure status = %d body=%s, want internal server error", duplicate.Code, duplicate.Body.String())
	}
	if !corrupted {
		t.Fatal("duplicate cleanup callback did not run")
	}
	assertDocumentCount(t, db, folderID, 1)
}

func TestDocumentUploadDeleteStoredUploadReturnsCleanupErrors(t *testing.T) {
	handler := &documentHandler{storage: documents.NewLocalStorage(t.TempDir(), 1024, nil)}

	if err := handler.deleteStoredUpload("../outside.txt"); !errors.Is(err, documents.ErrInvalidStorageKey) {
		t.Fatalf("deleteStoredUpload error = %v, want ErrInvalidStorageKey", err)
	}
}

func TestDocumentUploadRequiresTrustedInternalRequest(t *testing.T) {
	router, db := newAuthTestRouterWithOptions(t, RouterOptions{DocumentStorageRoot: t.TempDir()})
	token := loginSeededAdmin(t, router, db, "admin@example.com")
	unitID := createUnitViaAPI(t, router, token, `{"name":"Finance"}`)
	folderID := createFolderViaAPI(t, router, token, `{"organizationUnitId":"`+unitID+`","name":"Invoices"}`)

	response := uploadDocumentWithHeaders(t, router, map[string]string{
		"Cookie": authSessionCookieName + "=" + token,
	}, folderID, "invoice.pdf", []byte("%PDF test"))
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("missing trusted internal token status = %d body=%s, want unauthorized", response.Code, response.Body.String())
	}
	assertDocumentCount(t, db, folderID, 0)
}

func TestDocumentUploadRequiresAuthenticatedSession(t *testing.T) {
	router, db := newAuthTestRouterWithOptions(t, RouterOptions{DocumentStorageRoot: t.TempDir()})
	token := loginSeededAdmin(t, router, db, "admin@example.com")
	unitID := createUnitViaAPI(t, router, token, `{"name":"Finance"}`)
	folderID := createFolderViaAPI(t, router, token, `{"organizationUnitId":"`+unitID+`","name":"Invoices"}`)

	response := uploadDocumentWithHeaders(t, router, map[string]string{
		internalAPIHeader: testInternalToken,
	}, folderID, "invoice.pdf", []byte("%PDF test"))
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("missing session status = %d body=%s, want unauthorized", response.Code, response.Body.String())
	}
	assertDocumentCount(t, db, folderID, 0)
}

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
	var getBody documentMetadataResponse
	decodeJSON(t, get, &getBody)
	if getBody.ID != docID || getBody.DisplayName != "invoice.pdf" || getBody.OriginalFileName != "invoice.pdf" || getBody.OrganizationUnitID != unitID || getBody.FolderID != folderID {
		t.Fatalf("get body = %#v", getBody)
	}
	assertResponseOmitsStorageKey(t, get)

	rename := authJSON(t, router, http.MethodPatch, "/api/documents/"+docID, `{"displayName":"Renamed.pdf"}`, authCookieHeaders(token))
	if rename.Code != http.StatusOK {
		t.Fatalf("rename status = %d body=%s", rename.Code, rename.Body.String())
	}
	var renameBody documentMetadataResponse
	decodeJSON(t, rename, &renameBody)
	if renameBody.DisplayName != "Renamed.pdf" || renameBody.ID != docID {
		t.Fatalf("rename body = %#v", renameBody)
	}
	assertResponseOmitsStorageKey(t, rename)

	var renamed documents.Document
	if err := db.First(&renamed, "id = ?", docID).Error; err != nil {
		t.Fatalf("load renamed document: %v", err)
	}
	if renamed.DisplayName != "Renamed.pdf" || renamed.UpdatedByUserID == nil {
		t.Fatalf("renamed document = %#v", renamed)
	}

	download := authJSON(t, router, http.MethodGet, "/api/documents/"+docID+"/download", "", authCookieHeaders(token))
	if download.Code != http.StatusOK || download.Body.String() != "hello" {
		t.Fatalf("download status = %d body=%q", download.Code, download.Body.String())
	}
	if got := download.Header().Get("Content-Type"); got != "text/plain; charset=utf-8" {
		t.Fatalf("download content type = %q, want text/plain; charset=utf-8", got)
	}
	if got := download.Header().Get("Content-Length"); got != "5" {
		t.Fatalf("download content length = %q, want 5", got)
	}
	mediaType, params, err := mime.ParseMediaType(download.Header().Get("Content-Disposition"))
	if err != nil {
		t.Fatalf("parse content disposition: %v", err)
	}
	if mediaType != "attachment" || params["filename"] != "invoice.pdf" {
		t.Fatalf("content disposition = %q params=%v", mediaType, params)
	}

	archive := authJSON(t, router, http.MethodPost, "/api/documents/"+docID+"/archive", `{}`, authCookieHeaders(token))
	if archive.Code != http.StatusOK {
		t.Fatalf("archive status = %d body=%s", archive.Code, archive.Body.String())
	}
	var archived documents.Document
	if err := db.First(&archived, "id = ?", docID).Error; err != nil {
		t.Fatalf("load archived document: %v", err)
	}
	if archived.DeletedAt == nil || archived.UpdatedByUserID == nil {
		t.Fatalf("archived document = %#v", archived)
	}
}

func TestDocumentMetadataOperationsRequireObjectScopedPermissions(t *testing.T) {
	router, db := newAuthTestRouterWithOptions(t, RouterOptions{DocumentStorageRoot: t.TempDir()})
	adminToken := loginSeededAdmin(t, router, db, "admin@example.com")
	unitID := createUnitViaAPI(t, router, adminToken, `{"name":"Finance"}`)
	folderID := createFolderViaAPI(t, router, adminToken, `{"organizationUnitId":"`+unitID+`","name":"Invoices"}`)
	docID := uploadDocumentID(t, router, adminToken, folderID, "invoice.pdf", []byte("hello"))

	for _, tc := range []struct {
		name              string
		missingPermission string
		method            string
		path              string
		body              string
	}{
		{name: "view", missingPermission: "document.view", method: http.MethodGet, path: "/api/documents/" + docID},
		{name: "update", missingPermission: "document.update", method: http.MethodPatch, path: "/api/documents/" + docID, body: `{"displayName":"Denied.pdf"}`},
		{name: "delete", missingPermission: "document.delete", method: http.MethodPost, path: "/api/documents/" + docID + "/archive", body: `{}`},
		{name: "download", missingPermission: "document.download", method: http.MethodGet, path: "/api/documents/" + docID + "/download"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			user := createVerifiedUser(t, db, tc.name+"-metadata-denied@example.com", "password123")
			assignDocumentOperationRoleWithoutPermission(t, db, user.ID, unitID, tc.missingPermission)
			token := loginUser(t, router, user.Email, "password123")

			response := authJSON(t, router, tc.method, tc.path, tc.body, authCookieHeaders(token))
			if response.Code != http.StatusNotFound {
				t.Fatalf("%s status = %d body=%s, want not found", tc.name, response.Code, response.Body.String())
			}
		})
	}
}

func TestDocumentMetadataOperationsHideArchivedDocuments(t *testing.T) {
	router, db := newAuthTestRouterWithOptions(t, RouterOptions{DocumentStorageRoot: t.TempDir()})
	token := loginSeededAdmin(t, router, db, "admin@example.com")
	unitID := createUnitViaAPI(t, router, token, `{"name":"Finance"}`)
	folderID := createFolderViaAPI(t, router, token, `{"organizationUnitId":"`+unitID+`","name":"Invoices"}`)
	docID := uploadDocumentID(t, router, token, folderID, "invoice.pdf", []byte("hello"))
	if err := db.Model(&documents.Document{}).Where("id = ?", docID).Update("deleted_at", time.Now().UTC()).Error; err != nil {
		t.Fatalf("archive document directly: %v", err)
	}

	assertDocumentMetadataOperationsNotFound(t, router, token, docID)
}

func TestDocumentMetadataOperationsHideDocumentsInArchivedFolderOrOrganizationUnit(t *testing.T) {
	for _, tc := range []struct {
		name    string
		archive func(t *testing.T, db *gorm.DB, unitID string, folderID string)
	}{
		{
			name: "folder",
			archive: func(t *testing.T, db *gorm.DB, _ string, folderID string) {
				t.Helper()
				if err := db.Model(&documents.Folder{}).Where("id = ?", folderID).Update("deleted_at", time.Now().UTC()).Error; err != nil {
					t.Fatalf("archive folder directly: %v", err)
				}
			},
		},
		{
			name: "organization unit",
			archive: func(t *testing.T, db *gorm.DB, unitID string, _ string) {
				t.Helper()
				if err := db.Model(&orgunits.Unit{}).Where("id = ?", unitID).Update("archived_at", time.Now().UTC()).Error; err != nil {
					t.Fatalf("archive organization unit directly: %v", err)
				}
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			router, db := newAuthTestRouterWithOptions(t, RouterOptions{DocumentStorageRoot: t.TempDir()})
			token := loginSeededAdmin(t, router, db, "admin@example.com")
			unitID := createUnitViaAPI(t, router, token, `{"name":"Finance"}`)
			folderID := createFolderViaAPI(t, router, token, `{"organizationUnitId":"`+unitID+`","name":"Invoices"}`)
			docID := uploadDocumentID(t, router, token, folderID, "invoice.pdf", []byte("hello"))
			tc.archive(t, db, unitID, folderID)

			assertDocumentMetadataOperationsNotFound(t, router, token, docID)
		})
	}
}

func TestDocumentMetadataDownloadReturnsNotFoundWhenStorageFileMissing(t *testing.T) {
	storageRoot := t.TempDir()
	router, db := newAuthTestRouterWithOptions(t, RouterOptions{DocumentStorageRoot: storageRoot})
	token := loginSeededAdmin(t, router, db, "admin@example.com")
	unitID := createUnitViaAPI(t, router, token, `{"name":"Finance"}`)
	folderID := createFolderViaAPI(t, router, token, `{"organizationUnitId":"`+unitID+`","name":"Invoices"}`)
	docID := uploadDocumentID(t, router, token, folderID, "invoice.pdf", []byte("hello"))
	stored := loadDocumentByID(t, db, docID)
	if err := os.Remove(filepath.Join(storageRoot, filepath.FromSlash(stored.StorageKey))); err != nil {
		t.Fatalf("remove stored file: %v", err)
	}

	download := authJSON(t, router, http.MethodGet, "/api/documents/"+docID+"/download", "", authCookieHeaders(token))
	if download.Code != http.StatusNotFound {
		t.Fatalf("download missing storage status = %d body=%s, want not found", download.Code, download.Body.String())
	}
}

func TestDocumentMetadataDownloadReturnsServerErrorForInvalidStorageKey(t *testing.T) {
	router, db := newAuthTestRouterWithOptions(t, RouterOptions{DocumentStorageRoot: t.TempDir()})
	token := loginSeededAdmin(t, router, db, "admin@example.com")
	unitID := createUnitViaAPI(t, router, token, `{"name":"Finance"}`)
	folderID := createFolderViaAPI(t, router, token, `{"organizationUnitId":"`+unitID+`","name":"Invoices"}`)
	docID := uploadDocumentID(t, router, token, folderID, "invoice.pdf", []byte("hello"))
	if err := db.Model(&documents.Document{}).Where("id = ?", docID).Update("storage_key", "../outside.txt").Error; err != nil {
		t.Fatalf("corrupt storage key: %v", err)
	}

	download := authJSON(t, router, http.MethodGet, "/api/documents/"+docID+"/download", "", authCookieHeaders(token))
	if download.Code != http.StatusInternalServerError {
		t.Fatalf("download invalid storage status = %d body=%s, want internal server error", download.Code, download.Body.String())
	}
}

func TestDocumentMetadataDownloadReturnsServerErrorForMismatchedStorageSize(t *testing.T) {
	storageRoot := t.TempDir()
	router, db := newAuthTestRouterWithOptions(t, RouterOptions{DocumentStorageRoot: storageRoot})
	token := loginSeededAdmin(t, router, db, "admin@example.com")
	unitID := createUnitViaAPI(t, router, token, `{"name":"Finance"}`)
	folderID := createFolderViaAPI(t, router, token, `{"organizationUnitId":"`+unitID+`","name":"Invoices"}`)
	docID := uploadDocumentID(t, router, token, folderID, "invoice.pdf", []byte("hello"))
	stored := loadDocumentByID(t, db, docID)
	storedPath := filepath.Join(storageRoot, filepath.FromSlash(stored.StorageKey))
	if err := os.WriteFile(storedPath, []byte("corrupted"), 0o600); err != nil {
		t.Fatalf("overwrite stored file: %v", err)
	}

	download := authJSON(t, router, http.MethodGet, "/api/documents/"+docID+"/download", "", authCookieHeaders(token))
	if download.Code != http.StatusInternalServerError {
		t.Fatalf("download mismatched storage size status = %d body=%s, want internal server error", download.Code, download.Body.String())
	}
}

func TestDocumentMetadataRenameRejectsInvalidDisplayNames(t *testing.T) {
	router, db := newAuthTestRouterWithOptions(t, RouterOptions{DocumentStorageRoot: t.TempDir()})
	token := loginSeededAdmin(t, router, db, "admin@example.com")
	unitID := createUnitViaAPI(t, router, token, `{"name":"Finance"}`)
	folderID := createFolderViaAPI(t, router, token, `{"organizationUnitId":"`+unitID+`","name":"Invoices"}`)
	docID := uploadDocumentID(t, router, token, folderID, "invoice.pdf", []byte("hello"))

	for _, tc := range []struct {
		name string
		body string
	}{
		{name: "empty", body: `{"displayName":"   "}`},
		{name: "too long", body: `{"displayName":"` + strings.Repeat("a", documents.MaxDocumentDisplayNameCharacters+1) + `"}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			response := authJSON(t, router, http.MethodPatch, "/api/documents/"+docID, tc.body, authCookieHeaders(token))
			if response.Code != http.StatusBadRequest {
				t.Fatalf("%s status = %d body=%s, want bad request", tc.name, response.Code, response.Body.String())
			}
			stored := loadDocumentByID(t, db, docID)
			if stored.DisplayName != "invoice.pdf" {
				t.Fatalf("display name after invalid rename = %q, want invoice.pdf", stored.DisplayName)
			}
		})
	}
}

func uploadDocument(t *testing.T, router http.Handler, token string, folderID string, fileName string, content []byte) *httptest.ResponseRecorder {
	t.Helper()
	return uploadDocumentWithHeaders(t, router, authCookieHeaders(token), folderID, fileName, content)
}

func uploadDocumentID(t *testing.T, router http.Handler, token string, folderID string, fileName string, content []byte) string {
	t.Helper()
	response := uploadDocument(t, router, token, folderID, fileName, content)
	if response.Code != http.StatusCreated {
		t.Fatalf("upload status = %d body=%s", response.Code, response.Body.String())
	}
	var body documentMetadataResponse
	decodeJSON(t, response, &body)
	if body.ID == "" {
		t.Fatal("uploaded document id was empty")
	}
	return body.ID
}

func uploadDocumentWithHeaders(t *testing.T, router http.Handler, headers map[string]string, folderID string, fileName string, content []byte) *httptest.ResponseRecorder {
	t.Helper()
	return uploadDocumentFormWithHeaders(t, router, headers, func(writer *multipart.Writer) {
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
	return uploadDocumentFormWithHeaders(t, router, authCookieHeaders(token), writeForm)
}

func uploadDocumentFormWithHeaders(t *testing.T, router http.Handler, headers map[string]string, writeForm func(*multipart.Writer)) *httptest.ResponseRecorder {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	writeForm(writer)
	if err := writer.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/documents/upload", &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	for key, value := range headers {
		request.Header.Set(key, value)
	}
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

func documentUploadCallbackName(t *testing.T, suffix string) string {
	t.Helper()
	name := strings.NewReplacer("/", "_", " ", "_").Replace(t.Name())
	return name + "_" + suffix
}

func assertDocumentMetadataOperationsNotFound(t *testing.T, router http.Handler, token string, docID string) {
	t.Helper()
	for _, tc := range []struct {
		name   string
		method string
		path   string
		body   string
	}{
		{name: "get", method: http.MethodGet, path: "/api/documents/" + docID},
		{name: "rename", method: http.MethodPatch, path: "/api/documents/" + docID, body: `{"displayName":"Hidden.pdf"}`},
		{name: "archive", method: http.MethodPost, path: "/api/documents/" + docID + "/archive", body: `{}`},
		{name: "download", method: http.MethodGet, path: "/api/documents/" + docID + "/download"},
	} {
		response := authJSON(t, router, tc.method, tc.path, tc.body, authCookieHeaders(token))
		if response.Code != http.StatusNotFound {
			t.Fatalf("%s status = %d body=%s, want not found", tc.name, response.Code, response.Body.String())
		}
	}
}

func loadDocumentByID(t *testing.T, db *gorm.DB, documentID string) documents.Document {
	t.Helper()
	var stored documents.Document
	if err := db.First(&stored, "id = ?", documentID).Error; err != nil {
		t.Fatalf("load document %s: %v", documentID, err)
	}
	return stored
}

func assignDocumentOperationRoleWithoutPermission(t *testing.T, db *gorm.DB, userID string, organizationUnitID string, excludedPermission string) {
	t.Helper()
	now := time.Now().UTC()
	role := rbac.Role{
		Name:      "Document operation test " + strings.ReplaceAll(excludedPermission, ".", " "),
		Code:      "doc_test_" + strings.ReplaceAll(excludedPermission, ".", "_") + "_" + strings.ReplaceAll(userID[:8], "-", ""),
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("create document operation test role: %v", err)
	}

	for _, permissionCode := range []string{"document.view", "document.create", "document.update", "document.delete", "document.download"} {
		if permissionCode == excludedPermission {
			continue
		}
		permission := loadPermissionByCode(t, db, permissionCode)
		rolePermission := rbac.RolePermission{
			RoleID:       role.ID,
			PermissionID: permission.ID,
			CreatedAt:    now,
		}
		if err := db.Create(&rolePermission).Error; err != nil {
			t.Fatalf("assign permission %s to document operation test role: %v", permissionCode, err)
		}
	}

	assignment := rbac.UserRole{
		UserID:             userID,
		RoleID:             role.ID,
		ScopeType:          rbac.ScopeOrganizationUnit,
		OrganizationUnitID: &organizationUnitID,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	if err := db.Create(&assignment).Error; err != nil {
		t.Fatalf("assign document operation test role: %v", err)
	}
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

func storedDocumentFileSet(t *testing.T, storageRoot string) map[string]bool {
	t.Helper()
	out := map[string]bool{}
	documentsRoot := filepath.Join(storageRoot, "documents")
	if _, err := os.Stat(documentsRoot); os.IsNotExist(err) {
		return out
	} else if err != nil {
		t.Fatalf("stat documents storage root: %v", err)
	}
	if err := filepath.WalkDir(documentsRoot, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Type().IsRegular() {
			out[path] = true
		}
		return nil
	}); err != nil {
		t.Fatalf("walk documents storage: %v", err)
	}
	return out
}

func replaceNewStoredFileWithSymlink(t *testing.T, storageRoot string, originalFiles map[string]bool) {
	t.Helper()
	documentsRoot := filepath.Join(storageRoot, "documents")
	var target string
	if err := filepath.WalkDir(documentsRoot, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Type().IsRegular() && !originalFiles[path] && target == "" {
			target = path
		}
		return nil
	}); err != nil {
		t.Fatalf("walk documents storage for cleanup corruption: %v", err)
	}
	if target == "" {
		t.Fatal("new stored document file was not found")
	}

	outside := filepath.Join(t.TempDir(), "outside.bin")
	if err := os.WriteFile(outside, []byte("outside"), 0o600); err != nil {
		t.Fatalf("write outside symlink target: %v", err)
	}
	if err := os.Remove(target); err != nil {
		t.Fatalf("remove stored file before symlink replacement: %v", err)
	}
	if err := os.Symlink(outside, target); err != nil {
		t.Fatalf("create stored file symlink: %v", err)
	}
}
