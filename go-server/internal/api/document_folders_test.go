package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"ai.ro/syncra/dms/internal/auth"
	"ai.ro/syncra/dms/internal/documents"
	"gorm.io/gorm"
)

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

	rootID := createFolderViaAPI(t, router, token, `{"organizationUnitId":"`+unitID+`","name":"Invoices","description":" Monthly invoices "}`)
	childID := createFolderViaAPI(t, router, token, `{"organizationUnitId":"`+unitID+`","parentId":"`+rootID+`","name":"2026"}`)

	update := folderJSON(t, router, http.MethodPatch, "/api/document-folders/"+rootID, `{"organizationUnitId":"`+unitID+`","name":"Invoices","description":"Paid invoices"}`, authCookieHeaders(token))
	if update.Code != http.StatusOK {
		t.Fatalf("update status = %d body=%s", update.Code, update.Body.String())
	}
	var updateBody documentFolderTestResponse
	decodeJSON(t, update, &updateBody)
	if updateBody.Description == nil || *updateBody.Description != "Paid invoices" {
		t.Fatalf("update description = %#v, want Paid invoices", updateBody.Description)
	}

	tree := folderJSON(t, router, http.MethodGet, "/api/document-folders/tree?organizationUnitId="+unitID, "", authCookieHeaders(token))
	if tree.Code != http.StatusOK {
		t.Fatalf("tree status = %d body=%s", tree.Code, tree.Body.String())
	}
	var treeBody documentFolderTestListResponse
	decodeJSON(t, tree, &treeBody)
	if len(treeBody.Folders) != 1 {
		t.Fatalf("tree roots = %#v, want one root", treeBody.Folders)
	}
	if treeBody.Folders[0].ID != rootID || treeBody.Folders[0].Name != "Invoices" {
		t.Fatalf("tree root = %#v, want Invoices", treeBody.Folders[0])
	}
	if len(treeBody.Folders[0].Children) != 1 || treeBody.Folders[0].Children[0].ID != childID {
		t.Fatalf("tree child nesting = %#v, want child under root", treeBody.Folders[0].Children)
	}

	move := folderJSON(t, router, http.MethodPatch, "/api/document-folders/"+childID+"/parent", `{"parentId":null}`, authCookieHeaders(token))
	if move.Code != http.StatusOK {
		t.Fatalf("move status = %d body=%s", move.Code, move.Body.String())
	}
	var moveBody documentFolderTestResponse
	decodeJSON(t, move, &moveBody)
	if moveBody.ParentID != nil {
		t.Fatalf("moved parentId = %#v, want nil", moveBody.ParentID)
	}

	contents := folderJSON(t, router, http.MethodGet, "/api/document-folders/"+childID+"/contents", "", authCookieHeaders(token))
	if contents.Code != http.StatusOK {
		t.Fatalf("contents status = %d body=%s", contents.Code, contents.Body.String())
	}
	var contentsBody documentFolderContentsTestResponse
	decodeJSON(t, contents, &contentsBody)
	if contentsBody.Folder.ID != childID || len(contentsBody.Folders) != 0 || len(contentsBody.Documents) != 0 {
		t.Fatalf("contents body = %#v, want active folder with empty folders/documents", contentsBody)
	}

	archive := folderJSON(t, router, http.MethodPost, "/api/document-folders/"+rootID+"/archive", `{}`, authCookieHeaders(token))
	if archive.Code != http.StatusOK {
		t.Fatalf("archive status = %d body=%s", archive.Code, archive.Body.String())
	}
}

func TestDocumentFolderValidationRejectsDuplicatesAndInvalidParents(t *testing.T) {
	router, db := newAuthTestRouter(t)
	token := loginSeededAdmin(t, router, db, "admin@example.com")
	unitID := createUnitViaAPI(t, router, token, `{"name":"Finance"}`)
	otherUnitID := createUnitViaAPI(t, router, token, `{"name":"Legal"}`)

	rootID := createFolderViaAPI(t, router, token, `{"organizationUnitId":"`+unitID+`","name":"Invoices"}`)
	duplicateRoot := folderJSON(t, router, http.MethodPost, "/api/document-folders", `{"organizationUnitId":"`+unitID+`","name":"Invoices"}`, authCookieHeaders(token))
	if duplicateRoot.Code != http.StatusConflict {
		t.Fatalf("duplicate root status = %d body=%s, want conflict", duplicateRoot.Code, duplicateRoot.Body.String())
	}

	createFolderViaAPI(t, router, token, `{"organizationUnitId":"`+unitID+`","parentId":"`+rootID+`","name":"2026"}`)
	duplicateChild := folderJSON(t, router, http.MethodPost, "/api/document-folders", `{"organizationUnitId":"`+unitID+`","parentId":"`+rootID+`","name":"2026"}`, authCookieHeaders(token))
	if duplicateChild.Code != http.StatusConflict {
		t.Fatalf("duplicate child status = %d body=%s, want conflict", duplicateChild.Code, duplicateChild.Body.String())
	}

	invalidUnitID := folderJSON(t, router, http.MethodGet, "/api/document-folders/tree?organizationUnitId=not-a-uuid", "", authCookieHeaders(token))
	if invalidUnitID.Code != http.StatusBadRequest {
		t.Fatalf("invalid unit id status = %d body=%s, want bad request", invalidUnitID.Code, invalidUnitID.Body.String())
	}
	invalidParentID := folderJSON(t, router, http.MethodPost, "/api/document-folders", `{"organizationUnitId":"`+unitID+`","parentId":"not-a-uuid","name":"Bad Parent"}`, authCookieHeaders(token))
	if invalidParentID.Code != http.StatusBadRequest {
		t.Fatalf("invalid parent status = %d body=%s, want bad request", invalidParentID.Code, invalidParentID.Body.String())
	}
	missingParent := folderJSON(t, router, http.MethodPost, "/api/document-folders", `{"organizationUnitId":"`+unitID+`","parentId":"00000000-0000-0000-0000-000000000123","name":"Missing Parent"}`, authCookieHeaders(token))
	if missingParent.Code != http.StatusNotFound {
		t.Fatalf("missing parent status = %d body=%s, want not found", missingParent.Code, missingParent.Body.String())
	}
	otherParentID := createFolderViaAPI(t, router, token, `{"organizationUnitId":"`+otherUnitID+`","name":"Cases"}`)
	crossUnitCreate := folderJSON(t, router, http.MethodPost, "/api/document-folders", `{"organizationUnitId":"`+unitID+`","parentId":"`+otherParentID+`","name":"Cross Unit"}`, authCookieHeaders(token))
	if crossUnitCreate.Code != http.StatusConflict {
		t.Fatalf("cross-unit create status = %d body=%s, want conflict", crossUnitCreate.Code, crossUnitCreate.Body.String())
	}

	archiveRoot := folderJSON(t, router, http.MethodPost, "/api/document-folders/"+rootID+"/archive", `{}`, authCookieHeaders(token))
	if archiveRoot.Code != http.StatusOK {
		t.Fatalf("archive root status = %d body=%s", archiveRoot.Code, archiveRoot.Body.String())
	}
	recreateRoot := folderJSON(t, router, http.MethodPost, "/api/document-folders", `{"organizationUnitId":"`+unitID+`","name":"Invoices"}`, authCookieHeaders(token))
	if recreateRoot.Code != http.StatusCreated {
		t.Fatalf("recreate archived root status = %d body=%s, want created", recreateRoot.Code, recreateRoot.Body.String())
	}
}

func TestDocumentFolderMoveRejectsCyclesAndCrossUnitParents(t *testing.T) {
	router, db := newAuthTestRouter(t)
	token := loginSeededAdmin(t, router, db, "admin@example.com")
	unitID := createUnitViaAPI(t, router, token, `{"name":"Finance"}`)
	otherUnitID := createUnitViaAPI(t, router, token, `{"name":"Legal"}`)

	rootID := createFolderViaAPI(t, router, token, `{"organizationUnitId":"`+unitID+`","name":"Invoices"}`)
	childID := createFolderViaAPI(t, router, token, `{"organizationUnitId":"`+unitID+`","parentId":"`+rootID+`","name":"2026"}`)
	grandchildID := createFolderViaAPI(t, router, token, `{"organizationUnitId":"`+unitID+`","parentId":"`+childID+`","name":"Q1"}`)
	otherParentID := createFolderViaAPI(t, router, token, `{"organizationUnitId":"`+otherUnitID+`","name":"Cases"}`)

	moveUnderSelf := folderJSON(t, router, http.MethodPatch, "/api/document-folders/"+rootID+"/parent", `{"parentId":"`+rootID+`"}`, authCookieHeaders(token))
	if moveUnderSelf.Code != http.StatusConflict {
		t.Fatalf("move under self status = %d body=%s, want conflict", moveUnderSelf.Code, moveUnderSelf.Body.String())
	}
	moveUnderDescendant := folderJSON(t, router, http.MethodPatch, "/api/document-folders/"+rootID+"/parent", `{"parentId":"`+grandchildID+`"}`, authCookieHeaders(token))
	if moveUnderDescendant.Code != http.StatusConflict {
		t.Fatalf("move under descendant status = %d body=%s, want conflict", moveUnderDescendant.Code, moveUnderDescendant.Body.String())
	}
	moveUnderOtherUnit := folderJSON(t, router, http.MethodPatch, "/api/document-folders/"+childID+"/parent", `{"parentId":"`+otherParentID+`"}`, authCookieHeaders(token))
	if moveUnderOtherUnit.Code != http.StatusConflict {
		t.Fatalf("move under other unit status = %d body=%s, want conflict", moveUnderOtherUnit.Code, moveUnderOtherUnit.Body.String())
	}
	moveMissingParentKey := folderJSON(t, router, http.MethodPatch, "/api/document-folders/"+childID+"/parent", `{}`, authCookieHeaders(token))
	if moveMissingParentKey.Code != http.StatusBadRequest {
		t.Fatalf("move missing parent key status = %d body=%s, want bad request", moveMissingParentKey.Code, moveMissingParentKey.Body.String())
	}
}

func TestDocumentFolderArchiveCascadesAndTreeExcludesArchived(t *testing.T) {
	router, db := newAuthTestRouter(t)
	token := loginSeededAdmin(t, router, db, "admin@example.com")
	admin := loadUserByEmail(t, db, "admin@example.com")
	unitID := createUnitViaAPI(t, router, token, `{"name":"Finance"}`)

	rootID := createFolderViaAPI(t, router, token, `{"organizationUnitId":"`+unitID+`","name":"Invoices"}`)
	childID := createFolderViaAPI(t, router, token, `{"organizationUnitId":"`+unitID+`","parentId":"`+rootID+`","name":"2026"}`)
	siblingID := createFolderViaAPI(t, router, token, `{"organizationUnitId":"`+unitID+`","name":"Receipts"}`)
	createDocumentInFolder(t, db, unitID, rootID, admin.ID, "root.pdf")
	createDocumentInFolder(t, db, unitID, childID, admin.ID, "child.pdf")
	activeDocumentID := createDocumentInFolder(t, db, unitID, siblingID, admin.ID, "sibling.pdf")

	archive := folderJSON(t, router, http.MethodPost, "/api/document-folders/"+rootID+"/archive", `{}`, authCookieHeaders(token))
	if archive.Code != http.StatusOK {
		t.Fatalf("archive status = %d body=%s", archive.Code, archive.Body.String())
	}

	var archivedFolders int64
	if err := db.Model(&documents.Folder{}).Where("id IN ?", []string{rootID, childID}).Where("deleted_at IS NOT NULL").Count(&archivedFolders).Error; err != nil {
		t.Fatalf("count archived folders: %v", err)
	}
	if archivedFolders != 2 {
		t.Fatalf("archived folder count = %d, want 2", archivedFolders)
	}
	var archivedDocuments int64
	if err := db.Model(&documents.Document{}).Where("folder_id IN ?", []string{rootID, childID}).Where("deleted_at IS NOT NULL").Count(&archivedDocuments).Error; err != nil {
		t.Fatalf("count archived documents: %v", err)
	}
	if archivedDocuments != 2 {
		t.Fatalf("archived document count = %d, want 2", archivedDocuments)
	}
	var activeDocuments int64
	if err := db.Model(&documents.Document{}).Where("id = ? AND deleted_at IS NULL", activeDocumentID).Count(&activeDocuments).Error; err != nil {
		t.Fatalf("count active sibling document: %v", err)
	}
	if activeDocuments != 1 {
		t.Fatalf("active sibling document count = %d, want 1", activeDocuments)
	}

	tree := folderJSON(t, router, http.MethodGet, "/api/document-folders/tree?organizationUnitId="+unitID, "", authCookieHeaders(token))
	if tree.Code != http.StatusOK {
		t.Fatalf("tree status = %d body=%s", tree.Code, tree.Body.String())
	}
	var treeBody documentFolderTestListResponse
	decodeJSON(t, tree, &treeBody)
	if len(treeBody.Folders) != 1 || treeBody.Folders[0].ID != siblingID {
		t.Fatalf("tree after archive = %#v, want only sibling folder", treeBody)
	}

	updateArchived := folderJSON(t, router, http.MethodPatch, "/api/document-folders/"+rootID, `{"organizationUnitId":"`+unitID+`","name":"Archived"}`, authCookieHeaders(token))
	if updateArchived.Code != http.StatusNotFound {
		t.Fatalf("update archived status = %d body=%s, want not found", updateArchived.Code, updateArchived.Body.String())
	}
}

type documentFolderTestListResponse struct {
	Folders []documentFolderTestResponse `json:"folders"`
}

type documentFolderTestResponse struct {
	ID                 string                       `json:"id"`
	ParentID           *string                      `json:"parentId"`
	OrganizationUnitID string                       `json:"organizationUnitId"`
	Name               string                       `json:"name"`
	Description        *string                      `json:"description"`
	DeletedAt          *string                      `json:"deletedAt"`
	Children           []documentFolderTestResponse `json:"children"`
}

type documentFolderContentsTestResponse struct {
	Folder    documentFolderTestResponse   `json:"folder"`
	Folders   []documentFolderTestResponse `json:"folders"`
	Documents []any                        `json:"documents"`
}

func folderJSON(t *testing.T, router http.Handler, method string, path string, body string, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	return authJSON(t, router, method, path, body, headers)
}

func createFolderViaAPI(t *testing.T, router http.Handler, token string, body string) string {
	t.Helper()
	response := folderJSON(t, router, http.MethodPost, "/api/document-folders", body, authCookieHeaders(token))
	if response.Code != http.StatusCreated {
		t.Fatalf("create folder status = %d body=%s", response.Code, response.Body.String())
	}
	var out struct {
		ID string `json:"id"`
	}
	decodeJSON(t, response, &out)
	if out.ID == "" {
		t.Fatal("created folder id was empty")
	}
	return out.ID
}

func createDocumentInFolder(t *testing.T, db *gorm.DB, organizationUnitID string, folderID string, creatorID string, fileName string) string {
	t.Helper()
	now := time.Now().UTC()
	hash := strings.Repeat("a", 63) + "1"
	doc := documents.Document{
		FolderID:           folderID,
		OrganizationUnitID: organizationUnitID,
		OriginalFileName:   fileName,
		DisplayName:        fileName,
		MimeType:           "application/pdf",
		SizeBytes:          1,
		SHA256Hash:         hash,
		StorageKey:         "test/" + fileName,
		CreatedByUserID:    creatorID,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	if err := db.Create(&doc).Error; err != nil {
		t.Fatalf("create document: %v", err)
	}
	return doc.ID
}

func loadUserByEmail(t *testing.T, db *gorm.DB, email string) auth.User {
	t.Helper()
	var user auth.User
	if err := db.First(&user, "email = ?", email).Error; err != nil {
		t.Fatalf("load user %s: %v", email, err)
	}
	return user
}
