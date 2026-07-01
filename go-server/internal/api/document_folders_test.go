package api

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"ai.ro/syncra/dms/internal/auth"
	"ai.ro/syncra/dms/internal/documents"
	"ai.ro/syncra/dms/internal/orgunits"
	"ai.ro/syncra/dms/internal/rbac"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestDocumentFolderRoutesRequirePermissions(t *testing.T) {
	router, db := newAuthTestRouter(t)
	adminToken := loginSeededAdmin(t, router, db, "admin@example.com")
	unitID := createUnitViaAPI(t, router, adminToken, `{"name":"Finance"}`)
	rootID := createFolderViaAPI(t, router, adminToken, `{"organizationUnitId":"`+unitID+`","name":"Invoices"}`)
	childID := createFolderViaAPI(t, router, adminToken, `{"organizationUnitId":"`+unitID+`","parentId":"`+rootID+`","name":"2026"}`)
	user := createVerifiedUser(t, db, "viewer@example.com", "password123")
	token := loginUser(t, router, user.Email, "password123")

	for _, tc := range []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
	}{
		{name: "tree", method: http.MethodGet, path: "/api/document-folders/tree?organizationUnitId=" + unitID, wantStatus: http.StatusForbidden},
		{name: "create", method: http.MethodPost, path: "/api/document-folders", body: `{"organizationUnitId":"` + unitID + `","name":"Receipts"}`, wantStatus: http.StatusForbidden},
		{name: "update", method: http.MethodPatch, path: "/api/document-folders/" + rootID, body: `{"organizationUnitId":"` + unitID + `","name":"Invoices Updated"}`, wantStatus: http.StatusNotFound},
		{name: "move", method: http.MethodPatch, path: "/api/document-folders/" + childID + "/parent", body: `{"parentId":null}`, wantStatus: http.StatusNotFound},
		{name: "archive", method: http.MethodPost, path: "/api/document-folders/" + rootID + "/archive", body: `{}`, wantStatus: http.StatusNotFound},
		{name: "contents", method: http.MethodGet, path: "/api/document-folders/" + rootID + "/contents", wantStatus: http.StatusNotFound},
	} {
		response := folderJSON(t, router, tc.method, tc.path, tc.body, authCookieHeaders(token))
		if response.Code != tc.wantStatus {
			t.Fatalf("%s status = %d body=%s, want %d", tc.name, response.Code, response.Body.String(), tc.wantStatus)
		}
	}
}

func TestDocumentFolderRoutesAuthenticateBeforeExistenceChecks(t *testing.T) {
	router, db := newAuthTestRouter(t)
	adminToken := loginSeededAdmin(t, router, db, "admin@example.com")
	unitID := createUnitViaAPI(t, router, adminToken, `{"name":"Finance"}`)
	rootID := createFolderViaAPI(t, router, adminToken, `{"organizationUnitId":"`+unitID+`","name":"Invoices"}`)
	missingID := "00000000-0000-0000-0000-000000000123"
	internalOnly := map[string]string{internalAPIHeader: testInternalToken}

	for _, tc := range []struct {
		name   string
		method string
		path   string
		body   string
	}{
		{name: "tree missing unit", method: http.MethodGet, path: "/api/document-folders/tree?organizationUnitId=" + missingID},
		{name: "create missing unit", method: http.MethodPost, path: "/api/document-folders", body: `{"organizationUnitId":"` + missingID + `","name":"Invoices"}`},
		{name: "update missing folder", method: http.MethodPatch, path: "/api/document-folders/" + missingID, body: `{"organizationUnitId":"` + unitID + `","name":"Invoices Updated"}`},
		{name: "update existing folder", method: http.MethodPatch, path: "/api/document-folders/" + rootID, body: `{"organizationUnitId":"` + unitID + `","name":"Invoices Updated"}`},
		{name: "move missing folder", method: http.MethodPatch, path: "/api/document-folders/" + missingID + "/parent", body: `{"parentId":null}`},
		{name: "move existing folder", method: http.MethodPatch, path: "/api/document-folders/" + rootID + "/parent", body: `{"parentId":null}`},
		{name: "archive missing folder", method: http.MethodPost, path: "/api/document-folders/" + missingID + "/archive", body: `{}`},
		{name: "archive existing folder", method: http.MethodPost, path: "/api/document-folders/" + rootID + "/archive", body: `{}`},
		{name: "contents missing folder", method: http.MethodGet, path: "/api/document-folders/" + missingID + "/contents"},
		{name: "contents existing folder", method: http.MethodGet, path: "/api/document-folders/" + rootID + "/contents"},
	} {
		response := folderJSON(t, router, tc.method, tc.path, tc.body, internalOnly)
		if response.Code != http.StatusUnauthorized {
			t.Fatalf("%s status = %d body=%s, want unauthorized", tc.name, response.Code, response.Body.String())
		}
	}
}

func TestDocumentFolderRoutesDenyCrossUnitScopedDocumentPermissions(t *testing.T) {
	router, db := newAuthTestRouter(t)
	adminToken := loginSeededAdmin(t, router, db, "admin@example.com")
	financeID := createUnitViaAPI(t, router, adminToken, `{"name":"Finance"}`)
	legalID := createUnitViaAPI(t, router, adminToken, `{"name":"Legal"}`)
	financeRootID := createFolderViaAPI(t, router, adminToken, `{"organizationUnitId":"`+financeID+`","name":"Invoices"}`)
	financeChildID := createFolderViaAPI(t, router, adminToken, `{"organizationUnitId":"`+financeID+`","parentId":"`+financeRootID+`","name":"2026"}`)
	legalFolderID := createFolderViaAPI(t, router, adminToken, `{"organizationUnitId":"`+legalID+`","name":"Cases"}`)
	user := createVerifiedUser(t, db, "finance-docs@example.com", "password123")
	assignOrganizationUnitRoleByCode(t, db, user.ID, rbac.OrganizationAdministratorRoleCode, financeID)
	token := loginUser(t, router, user.Email, "password123")
	missingID := "00000000-0000-0000-0000-000000000123"

	for _, tc := range []struct {
		name   string
		method string
		path   string
		body   string
	}{
		{name: "tree", method: http.MethodGet, path: "/api/document-folders/tree?organizationUnitId=" + legalID},
		{name: "create", method: http.MethodPost, path: "/api/document-folders", body: `{"organizationUnitId":"` + legalID + `","name":"Contracts"}`},
	} {
		response := folderJSON(t, router, tc.method, tc.path, tc.body, authCookieHeaders(token))
		if response.Code != http.StatusForbidden {
			t.Fatalf("%s status = %d body=%s, want forbidden", tc.name, response.Code, response.Body.String())
		}
	}

	for _, tc := range []struct {
		name         string
		method       string
		existingPath string
		missingPath  string
		body         string
	}{
		{name: "update", method: http.MethodPatch, existingPath: "/api/document-folders/" + legalFolderID, missingPath: "/api/document-folders/" + missingID, body: `{"organizationUnitId":"` + legalID + `","name":"Cases Updated"}`},
		{name: "move", method: http.MethodPatch, existingPath: "/api/document-folders/" + legalFolderID + "/parent", missingPath: "/api/document-folders/" + missingID + "/parent", body: `{"parentId":null}`},
		{name: "archive", method: http.MethodPost, existingPath: "/api/document-folders/" + legalFolderID + "/archive", missingPath: "/api/document-folders/" + missingID + "/archive", body: `{}`},
		{name: "contents", method: http.MethodGet, existingPath: "/api/document-folders/" + legalFolderID + "/contents", missingPath: "/api/document-folders/" + missingID + "/contents"},
	} {
		existing := folderJSON(t, router, tc.method, tc.existingPath, tc.body, authCookieHeaders(token))
		if existing.Code != http.StatusNotFound {
			t.Fatalf("%s existing cross-unit status = %d body=%s, want not found", tc.name, existing.Code, existing.Body.String())
		}
		missing := folderJSON(t, router, tc.method, tc.missingPath, tc.body, authCookieHeaders(token))
		if missing.Code != http.StatusNotFound {
			t.Fatalf("%s missing status = %d body=%s, want not found", tc.name, missing.Code, missing.Body.String())
		}
	}

	update := folderJSON(t, router, http.MethodPatch, "/api/document-folders/"+financeRootID, `{"organizationUnitId":"`+financeID+`","name":"Invoices Updated"}`, authCookieHeaders(token))
	if update.Code != http.StatusOK {
		t.Fatalf("update scoped folder status = %d body=%s, want ok", update.Code, update.Body.String())
	}
	move := folderJSON(t, router, http.MethodPatch, "/api/document-folders/"+financeChildID+"/parent", `{"parentId":null}`, authCookieHeaders(token))
	if move.Code != http.StatusOK {
		t.Fatalf("move scoped folder status = %d body=%s, want ok", move.Code, move.Body.String())
	}
	contents := folderJSON(t, router, http.MethodGet, "/api/document-folders/"+financeRootID+"/contents", "", authCookieHeaders(token))
	if contents.Code != http.StatusOK {
		t.Fatalf("contents scoped folder status = %d body=%s, want ok", contents.Code, contents.Body.String())
	}
	archive := folderJSON(t, router, http.MethodPost, "/api/document-folders/"+financeRootID+"/archive", `{}`, authCookieHeaders(token))
	if archive.Code != http.StatusOK {
		t.Fatalf("archive scoped folder status = %d body=%s, want ok", archive.Code, archive.Body.String())
	}
}

func TestDocumentFolderParentValidationHidesInaccessibleCrossUnitParentIDs(t *testing.T) {
	router, db := newAuthTestRouter(t)
	adminToken := loginSeededAdmin(t, router, db, "admin@example.com")
	financeID := createUnitViaAPI(t, router, adminToken, `{"name":"Finance"}`)
	legalID := createUnitViaAPI(t, router, adminToken, `{"name":"Legal"}`)
	financeRootID := createFolderViaAPI(t, router, adminToken, `{"organizationUnitId":"`+financeID+`","name":"Invoices"}`)
	financeChildID := createFolderViaAPI(t, router, adminToken, `{"organizationUnitId":"`+financeID+`","parentId":"`+financeRootID+`","name":"2026"}`)
	legalFolderID := createFolderViaAPI(t, router, adminToken, `{"organizationUnitId":"`+legalID+`","name":"Cases"}`)
	user := createVerifiedUser(t, db, "finance-parent-probe@example.com", "password123")
	assignOrganizationUnitRoleByCode(t, db, user.ID, rbac.OrganizationAdministratorRoleCode, financeID)
	token := loginUser(t, router, user.Email, "password123")
	missingID := "00000000-0000-0000-0000-000000000123"

	for _, tc := range []struct {
		name             string
		method           string
		path             string
		bodyWithExisting string
		bodyWithMissing  string
	}{
		{
			name:             "create",
			method:           http.MethodPost,
			path:             "/api/document-folders",
			bodyWithExisting: `{"organizationUnitId":"` + financeID + `","parentId":"` + legalFolderID + `","name":"Cross Unit Create"}`,
			bodyWithMissing:  `{"organizationUnitId":"` + financeID + `","parentId":"` + missingID + `","name":"Missing Parent Create"}`,
		},
		{
			name:             "move",
			method:           http.MethodPatch,
			path:             "/api/document-folders/" + financeChildID + "/parent",
			bodyWithExisting: `{"parentId":"` + legalFolderID + `"}`,
			bodyWithMissing:  `{"parentId":"` + missingID + `"}`,
		},
	} {
		existing := folderJSON(t, router, tc.method, tc.path, tc.bodyWithExisting, authCookieHeaders(token))
		if existing.Code != http.StatusNotFound {
			t.Fatalf("%s inaccessible parent status = %d body=%s, want not found", tc.name, existing.Code, existing.Body.String())
		}
		missing := folderJSON(t, router, tc.method, tc.path, tc.bodyWithMissing, authCookieHeaders(token))
		if missing.Code != http.StatusNotFound {
			t.Fatalf("%s missing parent status = %d body=%s, want not found", tc.name, missing.Code, missing.Body.String())
		}
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
		t.Fatalf("contents status = %d body=%s, want ok", contents.Code, contents.Body.String())
	}

	archive := folderJSON(t, router, http.MethodPost, "/api/document-folders/"+rootID+"/archive", `{}`, authCookieHeaders(token))
	if archive.Code != http.StatusOK {
		t.Fatalf("archive status = %d body=%s", archive.Code, archive.Body.String())
	}
}

func TestDocumentFolderContentsListsActiveChildrenAndDocuments(t *testing.T) {
	router, db := newAuthTestRouter(t)
	token := loginSeededAdmin(t, router, db, "admin@example.com")
	admin := loadUserByEmail(t, db, "admin@example.com")
	unitID := createUnitViaAPI(t, router, token, `{"name":"Finance"}`)

	rootID := createFolderViaAPI(t, router, token, `{"organizationUnitId":"`+unitID+`","name":"Invoices"}`)
	childID := createFolderViaAPI(t, router, token, `{"organizationUnitId":"`+unitID+`","parentId":"`+rootID+`","name":"2026"}`)
	archivedChildID := createFolderViaAPI(t, router, token, `{"organizationUnitId":"`+unitID+`","parentId":"`+rootID+`","name":"Archived"}`)
	archivedAt := time.Now().UTC()
	if err := db.Model(&documents.Folder{}).Where("id = ?", archivedChildID).Update("deleted_at", archivedAt).Error; err != nil {
		t.Fatalf("archive child folder: %v", err)
	}

	secondDocID := createDocumentInFolder(t, db, unitID, rootID, admin.ID, "b-invoice.pdf")
	firstDocID := createDocumentInFolder(t, db, unitID, rootID, admin.ID, "a-invoice.pdf")
	archivedDocID := createDocumentInFolder(t, db, unitID, rootID, admin.ID, "archived.pdf")
	if err := db.Model(&documents.Document{}).Where("id = ?", archivedDocID).Update("deleted_at", archivedAt).Error; err != nil {
		t.Fatalf("archive document: %v", err)
	}

	response := folderJSON(t, router, http.MethodGet, "/api/document-folders/"+rootID+"/contents", "", authCookieHeaders(token))
	if response.Code != http.StatusOK {
		t.Fatalf("contents status = %d body=%s, want ok", response.Code, response.Body.String())
	}
	var body documentFolderContentsTestResponse
	decodeJSON(t, response, &body)

	if body.Folder.ID != rootID || body.Folder.Name != "Invoices" {
		t.Fatalf("contents folder = %#v, want root Invoices", body.Folder)
	}
	if len(body.Folders) != 1 {
		t.Fatalf("contents folders = %#v, want one active child", body.Folders)
	}
	if body.Folders[0].ID != childID || body.Folders[0].Name != "2026" {
		t.Fatalf("contents child folder = %#v, want 2026", body.Folders[0])
	}

	if len(body.Documents) != 2 {
		t.Fatalf("contents documents = %#v, want two active documents", body.Documents)
	}
	expectedDocuments := []struct {
		id          string
		displayName string
	}{
		{id: firstDocID, displayName: "a-invoice.pdf"},
		{id: secondDocID, displayName: "b-invoice.pdf"},
	}
	for i, want := range expectedDocuments {
		got := body.Documents[i]
		if got.ID != want.id || got.DisplayName != want.displayName {
			t.Fatalf("contents document[%d] = %#v, want %s %s", i, got, want.id, want.displayName)
		}
		if got.FolderID != rootID || got.OrganizationUnitID != unitID {
			t.Fatalf("contents document[%d] scope = %#v, want folder/unit", i, got)
		}
		if got.OriginalFileName != want.displayName || got.MimeType != "application/pdf" || got.SizeBytes != 1 || got.SHA256Hash == "" {
			t.Fatalf("contents document[%d] metadata = %#v, want file metadata", i, got)
		}
		if got.StorageKey != nil {
			t.Fatalf("contents document[%d] storageKey = %q, want omitted", i, *got.StorageKey)
		}
		if got.DeletedAt != nil {
			t.Fatalf("contents document[%d] deletedAt = %q, want omitted for active document", i, *got.DeletedAt)
		}
		if got.CreatedAt == "" || got.UpdatedAt == "" {
			t.Fatalf("contents document[%d] timestamps = %#v, want populated", i, got)
		}
	}
	for _, folder := range body.Folders {
		if folder.ID == archivedChildID {
			t.Fatalf("contents included archived child folder %s", archivedChildID)
		}
	}
	for _, doc := range body.Documents {
		if doc.ID == archivedDocID {
			t.Fatalf("contents included archived document %s", archivedDocID)
		}
	}
}

func TestDocumentFolderContentsRequiresDocumentViewForFolderOrganizationUnit(t *testing.T) {
	router, db := newAuthTestRouter(t)
	adminToken := loginSeededAdmin(t, router, db, "admin@example.com")
	financeID := createUnitViaAPI(t, router, adminToken, `{"name":"Finance"}`)
	legalID := createUnitViaAPI(t, router, adminToken, `{"name":"Legal"}`)
	financeFolderID := createFolderViaAPI(t, router, adminToken, `{"organizationUnitId":"`+financeID+`","name":"Invoices"}`)
	legalFolderID := createFolderViaAPI(t, router, adminToken, `{"organizationUnitId":"`+legalID+`","name":"Cases"}`)
	user := createVerifiedUser(t, db, "legal-docs@example.com", "password123")
	assignOrganizationUnitRoleByCode(t, db, user.ID, rbac.OrganizationAdministratorRoleCode, legalID)
	token := loginUser(t, router, user.Email, "password123")

	denied := folderJSON(t, router, http.MethodGet, "/api/document-folders/"+financeFolderID+"/contents", "", authCookieHeaders(token))
	if denied.Code != http.StatusNotFound {
		t.Fatalf("contents without folder unit document.view status = %d body=%s, want not found", denied.Code, denied.Body.String())
	}
	allowed := folderJSON(t, router, http.MethodGet, "/api/document-folders/"+legalFolderID+"/contents", "", authCookieHeaders(token))
	if allowed.Code != http.StatusOK {
		t.Fatalf("contents with folder unit document.view status = %d body=%s, want ok", allowed.Code, allowed.Body.String())
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
	crossUnitUpdate := folderJSON(t, router, http.MethodPatch, "/api/document-folders/"+rootID, `{"organizationUnitId":"`+otherUnitID+`","name":"Invoices"}`, authCookieHeaders(token))
	if crossUnitUpdate.Code != http.StatusConflict {
		t.Fatalf("cross-unit update status = %d body=%s, want conflict", crossUnitUpdate.Code, crossUnitUpdate.Body.String())
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

func TestDocumentFolderByIDRoutesReturnNotFoundForArchivedOrganizationUnit(t *testing.T) {
	router, db := newAuthTestRouter(t)
	token := loginSeededAdmin(t, router, db, "admin@example.com")
	unitID := createUnitViaAPI(t, router, token, `{"name":"Finance"}`)
	rootID := createFolderViaAPI(t, router, token, `{"organizationUnitId":"`+unitID+`","name":"Invoices"}`)
	childID := createFolderViaAPI(t, router, token, `{"organizationUnitId":"`+unitID+`","parentId":"`+rootID+`","name":"2026"}`)

	archiveUnit := orgUnitJSON(t, router, http.MethodPost, "/api/organization-units/"+unitID+"/archive", `{}`, authCookieHeaders(token))
	if archiveUnit.Code != http.StatusOK {
		t.Fatalf("archive organization unit status = %d body=%s", archiveUnit.Code, archiveUnit.Body.String())
	}

	for _, tc := range []struct {
		name   string
		method string
		path   string
		body   string
	}{
		{name: "update", method: http.MethodPatch, path: "/api/document-folders/" + rootID, body: `{"organizationUnitId":"` + unitID + `","name":"Invoices Updated"}`},
		{name: "move", method: http.MethodPatch, path: "/api/document-folders/" + childID + "/parent", body: `{"parentId":null}`},
		{name: "contents", method: http.MethodGet, path: "/api/document-folders/" + rootID + "/contents"},
		{name: "archive", method: http.MethodPost, path: "/api/document-folders/" + rootID + "/archive", body: `{}`},
	} {
		response := folderJSON(t, router, tc.method, tc.path, tc.body, authCookieHeaders(token))
		if response.Code != http.StatusNotFound {
			t.Errorf("%s status = %d body=%s, want not found", tc.name, response.Code, response.Body.String())
		}
	}
}

func TestDocumentFolderUpdateRevalidatesOrganizationUnitAtWriteTime(t *testing.T) {
	router, db := newAuthTestRouter(t)
	token := loginSeededAdmin(t, router, db, "admin@example.com")
	unitID := createUnitViaAPI(t, router, token, `{"name":"Finance"}`)
	rootID := createFolderViaAPI(t, router, token, `{"organizationUnitId":"`+unitID+`","name":"Invoices"}`)

	callbackName := "document_folder_update_test_archive_unit_before_folder_update"
	archived := false
	if err := db.Callback().Update().Before("gorm:update").Register(callbackName, func(tx *gorm.DB) {
		if archived || tx.Statement.Table != "document_folders" {
			return
		}
		archived = true
		if err := tx.Session(&gorm.Session{NewDB: true}).Model(&orgunits.Unit{}).Where("id = ?", unitID).Update("archived_at", time.Now().UTC()).Error; err != nil {
			t.Fatalf("archive organization unit during update callback: %v", err)
		}
	}); err != nil {
		t.Fatalf("register update callback: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Callback().Update().Remove(callbackName); err != nil {
			t.Fatalf("remove update callback: %v", err)
		}
	})

	response := folderJSON(t, router, http.MethodPatch, "/api/document-folders/"+rootID, `{"organizationUnitId":"`+unitID+`","name":"Renamed"}`, authCookieHeaders(token))
	if response.Code != http.StatusNotFound {
		t.Fatalf("update status = %d body=%s, want not found", response.Code, response.Body.String())
	}
	if !archived {
		t.Fatal("organization unit archive callback did not run")
	}
	var folder documents.Folder
	if err := db.First(&folder, "id = ?", rootID).Error; err != nil {
		t.Fatalf("load folder after update: %v", err)
	}
	if folder.Name != "Invoices" {
		t.Fatalf("folder name = %q, want Invoices", folder.Name)
	}
}

func TestDocumentFolderActiveOrganizationUnitLockUsesPostgresRowLock(t *testing.T) {
	postgresDB, err := gorm.Open(postgres.New(postgres.Config{
		DSN: "postgres://syncra:syncra@localhost/syncra_dms?sslmode=disable",
	}), &gorm.Config{DryRun: true, DisableAutomaticPing: true})
	if err != nil {
		t.Fatalf("open postgres dry-run db: %v", err)
	}

	postgresStatement := activeOrganizationUnitLockQuery(t.Context(), postgresDB, "00000000-0000-0000-0000-000000000123").First(&orgunits.Unit{}).Statement
	postgresSQL := postgresStatement.SQL.String()
	if !strings.Contains(postgresSQL, "FOR UPDATE") {
		t.Fatalf("postgres active organization unit lock SQL = %q, want FOR UPDATE", postgresSQL)
	}
	if !strings.Contains(postgresSQL, "archived_at IS NULL") {
		t.Fatalf("postgres active organization unit lock SQL = %q, want active organization unit predicate", postgresSQL)
	}

	sqliteDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{DryRun: true})
	if err != nil {
		t.Fatalf("open sqlite dry-run db: %v", err)
	}
	sqliteStatement := activeOrganizationUnitLockQuery(t.Context(), sqliteDB, "00000000-0000-0000-0000-000000000123").First(&orgunits.Unit{}).Statement
	sqliteSQL := sqliteStatement.SQL.String()
	if strings.Contains(sqliteSQL, "FOR UPDATE") {
		t.Fatalf("sqlite active organization unit lock SQL = %q, want no FOR UPDATE", sqliteSQL)
	}
	if !strings.Contains(sqliteSQL, "archived_at IS NULL") {
		t.Fatalf("sqlite active organization unit lock SQL = %q, want active organization unit predicate", sqliteSQL)
	}
}

func TestDocumentFolderArchiveSuppressesPostLockValidationSentinel(t *testing.T) {
	router, db := newAuthTestRouter(t)
	token := loginSeededAdmin(t, router, db, "admin@example.com")
	unitID := createUnitViaAPI(t, router, token, `{"name":"Finance"}`)
	rootID := createFolderViaAPI(t, router, token, `{"organizationUnitId":"`+unitID+`","name":"Invoices"}`)

	callbackName := "document_folder_archive_test_delete_after_first_folder_load"
	deleted := false
	if err := db.Callback().Query().After("gorm:query").Register(callbackName, func(tx *gorm.DB) {
		if deleted || tx.Statement.Table != "document_folders" {
			return
		}
		deleted = true
		if err := tx.Session(&gorm.Session{NewDB: true}).Model(&documents.Folder{}).Where("id = ?", rootID).Update("deleted_at", time.Now().UTC()).Error; err != nil {
			t.Fatalf("delete folder during archive callback: %v", err)
		}
	}); err != nil {
		t.Fatalf("register archive callback: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Callback().Query().Remove(callbackName); err != nil {
			t.Fatalf("remove archive callback: %v", err)
		}
	})

	response := folderJSON(t, router, http.MethodPost, "/api/document-folders/"+rootID+"/archive", `{}`, authCookieHeaders(token))
	if response.Code != http.StatusNotFound {
		t.Fatalf("archive status = %d body=%s, want not found", response.Code, response.Body.String())
	}
	if strings.Contains(response.Body.String(), "failed to archive document folder") {
		t.Fatalf("archive body = %s, want only post-lock validation error", response.Body.String())
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
	Folder    documentFolderTestResponse     `json:"folder"`
	Folders   []documentFolderTestResponse   `json:"folders"`
	Documents []documentMetadataTestResponse `json:"documents"`
}

type documentMetadataTestResponse struct {
	ID                 string  `json:"id"`
	FolderID           string  `json:"folderId"`
	OrganizationUnitID string  `json:"organizationUnitId"`
	OriginalFileName   string  `json:"originalFileName"`
	DisplayName        string  `json:"displayName"`
	MimeType           string  `json:"mimeType"`
	Extension          *string `json:"extension"`
	SizeBytes          int64   `json:"sizeBytes"`
	SHA256Hash         string  `json:"sha256Hash"`
	StorageKey         *string `json:"storageKey"`
	DeletedAt          *string `json:"deletedAt"`
	CreatedAt          string  `json:"createdAt"`
	UpdatedAt          string  `json:"updatedAt"`
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
	sum := sha256.Sum256([]byte(folderID + ":" + fileName))
	hash := hex.EncodeToString(sum[:])
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

func assignOrganizationUnitRoleByCode(t *testing.T, db *gorm.DB, userID string, roleCode string, organizationUnitID string) {
	t.Helper()
	role := loadRoleByCode(t, db, roleCode)
	now := time.Now().UTC()
	assignment := rbac.UserRole{
		UserID:             userID,
		RoleID:             role.ID,
		ScopeType:          rbac.ScopeOrganizationUnit,
		OrganizationUnitID: &organizationUnitID,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	if err := db.Create(&assignment).Error; err != nil {
		t.Fatalf("assign organization unit role %s: %v", roleCode, err)
	}
}
