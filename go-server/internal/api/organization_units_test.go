package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"ai.ro/syncra/dms/internal/auth"
	"ai.ro/syncra/dms/internal/orgunits"
	"ai.ro/syncra/dms/internal/rbac"
	"gorm.io/gorm"
)

func TestOrganizationUnitRoutesRequireSessionAndAdminForMutations(t *testing.T) {
	router, db := newAuthTestRouter(t)

	missingInternalToken := orgUnitJSON(t, router, http.MethodGet, "/api/organization-units/tree", "", nil)
	if missingInternalToken.Code != http.StatusUnauthorized {
		t.Fatalf("missing internal token status = %d body=%s", missingInternalToken.Code, missingInternalToken.Body.String())
	}

	internalOnly := orgUnitJSON(t, router, http.MethodGet, "/api/organization-units/tree", "", map[string]string{
		internalAPIHeader: testInternalToken,
	})
	if internalOnly.Code != http.StatusUnauthorized {
		t.Fatalf("internal-only status = %d body=%s", internalOnly.Code, internalOnly.Body.String())
	}

	user := createVerifiedUser(t, db, "user@example.com", "password123")
	userToken := loginUser(t, router, user.Email, "password123")

	sessionOnly := orgUnitJSON(t, router, http.MethodGet, "/api/organization-units/tree", "", map[string]string{
		"Cookie": authSessionCookieName + "=" + userToken,
	})
	if sessionOnly.Code != http.StatusUnauthorized {
		t.Fatalf("session-only status = %d body=%s", sessionOnly.Code, sessionOnly.Body.String())
	}

	forbiddenTree := orgUnitJSON(t, router, http.MethodGet, "/api/organization-units/tree", "", authCookieHeaders(userToken))
	if forbiddenTree.Code != http.StatusForbidden {
		t.Fatalf("user list status = %d body=%s, want forbidden", forbiddenTree.Code, forbiddenTree.Body.String())
	}

	viewer := createVerifiedUser(t, db, "viewer@example.com", "password123")
	assignGlobalRoleByCode(t, db, viewer.ID, rbac.ViewerRoleCode)
	viewerToken := loginUser(t, router, viewer.Email, "password123")

	emptyTree := orgUnitJSON(t, router, http.MethodGet, "/api/organization-units/tree", "", authCookieHeaders(viewerToken))
	if emptyTree.Code != http.StatusOK {
		t.Fatalf("viewer list status = %d body=%s", emptyTree.Code, emptyTree.Body.String())
	}
	var emptyTreeBody organizationUnitTestListResponse
	decodeJSON(t, emptyTree, &emptyTreeBody)
	if len(emptyTreeBody.Units) != 0 {
		t.Fatalf("empty tree units = %#v, want none", emptyTreeBody.Units)
	}

	forbiddenCreate := orgUnitJSON(t, router, http.MethodPost, "/api/organization-units", `{"name":"Company"}`, authCookieHeaders(viewerToken))
	if forbiddenCreate.Code != http.StatusForbidden {
		t.Fatalf("viewer create status = %d body=%s", forbiddenCreate.Code, forbiddenCreate.Body.String())
	}

	admin := createAdminUser(t, db, "admin@example.com", "password123")
	adminToken := loginUser(t, router, admin.Email, "password123")
	legacyBeforeBootstrap := orgUnitJSON(t, router, http.MethodPost, "/api/organization-units", `{"name":"Legacy Before Bootstrap"}`, authCookieHeaders(adminToken))
	if legacyBeforeBootstrap.Code != http.StatusForbidden {
		t.Fatalf("legacy admin before bootstrap status = %d body=%s, want forbidden", legacyBeforeBootstrap.Code, legacyBeforeBootstrap.Body.String())
	}
	if err := rbac.BootstrapLegacyAdmins(db); err != nil {
		t.Fatalf("bootstrap legacy admin: %v", err)
	}
	created := orgUnitJSON(t, router, http.MethodPost, "/api/organization-units", `{"name":"Company","code":" root "}`, authCookieHeaders(adminToken))
	if created.Code != http.StatusCreated {
		t.Fatalf("admin create status = %d body=%s", created.Code, created.Body.String())
	}
	var createdBody struct {
		ID   string  `json:"id"`
		Code *string `json:"code"`
	}
	decodeJSON(t, created, &createdBody)
	if createdBody.ID == "" {
		t.Fatal("created unit id was empty")
	}
	if createdBody.Code == nil || *createdBody.Code != "ROOT" {
		t.Fatalf("created code = %#v, want ROOT", createdBody.Code)
	}

	list := orgUnitJSON(t, router, http.MethodGet, "/api/organization-units/tree", "", authCookieHeaders(viewerToken))
	if list.Code != http.StatusOK {
		t.Fatalf("viewer list status = %d body=%s", list.Code, list.Body.String())
	}
	var listBody struct {
		Units []struct {
			ID   string  `json:"id"`
			Name string  `json:"name"`
			Code *string `json:"code"`
		} `json:"units"`
	}
	decodeJSON(t, list, &listBody)
	if len(listBody.Units) != 1 || listBody.Units[0].ID != createdBody.ID || listBody.Units[0].Name != "Company" || listBody.Units[0].Code == nil || *listBody.Units[0].Code != "ROOT" {
		t.Fatalf("tree body = %#v", listBody)
	}

	forbiddenUpdate := orgUnitJSON(t, router, http.MethodPatch, "/api/organization-units/"+createdBody.ID, `{"name":"Company Updated"}`, authCookieHeaders(viewerToken))
	if forbiddenUpdate.Code != http.StatusForbidden {
		t.Fatalf("viewer update status = %d body=%s", forbiddenUpdate.Code, forbiddenUpdate.Body.String())
	}

	forbiddenMove := orgUnitJSON(t, router, http.MethodPatch, "/api/organization-units/"+createdBody.ID+"/parent", `{"parentId":null}`, authCookieHeaders(viewerToken))
	if forbiddenMove.Code != http.StatusForbidden {
		t.Fatalf("viewer move status = %d body=%s", forbiddenMove.Code, forbiddenMove.Body.String())
	}

	forbiddenArchive := orgUnitJSON(t, router, http.MethodPost, "/api/organization-units/"+createdBody.ID+"/archive", `{}`, authCookieHeaders(viewerToken))
	if forbiddenArchive.Code != http.StatusForbidden {
		t.Fatalf("viewer archive status = %d body=%s", forbiddenArchive.Code, forbiddenArchive.Body.String())
	}

	forbiddenArchivedList := orgUnitJSON(t, router, http.MethodGet, "/api/organization-units/archived", "", authCookieHeaders(viewerToken))
	if forbiddenArchivedList.Code != http.StatusForbidden {
		t.Fatalf("viewer archived list status = %d body=%s", forbiddenArchivedList.Code, forbiddenArchivedList.Body.String())
	}
}

func TestOrganizationUnitMoveRejectsCyclesAndArchiveCascades(t *testing.T) {
	router, db := newAuthTestRouter(t)
	token := loginSeededAdmin(t, router, db, "admin@example.com")

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

	tree := orgUnitJSON(t, router, http.MethodGet, "/api/organization-units/tree", "", authCookieHeaders(token))
	if tree.Code != http.StatusOK {
		t.Fatalf("tree after archive status = %d body=%s", tree.Code, tree.Body.String())
	}
	var treeBody organizationUnitTestListResponse
	decodeJSON(t, tree, &treeBody)
	if len(treeBody.Units) != 1 || len(treeBody.Units[0].Children) != 0 {
		t.Fatalf("tree after archive = %#v, want only root without archived child", treeBody)
	}

	archived := orgUnitJSON(t, router, http.MethodGet, "/api/organization-units/archived", "", authCookieHeaders(token))
	if archived.Code != http.StatusOK {
		t.Fatalf("archived status = %d body=%s", archived.Code, archived.Body.String())
	}
	var archivedBody organizationUnitTestListResponse
	decodeJSON(t, archived, &archivedBody)
	if len(archivedBody.Units) != 2 {
		t.Fatalf("archived units = %#v, want archived child and grandchild", archivedBody)
	}
	for _, unit := range archivedBody.Units {
		if unit.ArchivedAt == nil {
			t.Fatalf("archived unit %#v missing archivedAt", unit)
		}
	}
}

func TestOrganizationUnitTreeNestsOrdersAndExcludesArchived(t *testing.T) {
	router, db := newAuthTestRouter(t)
	token := loginSeededAdmin(t, router, db, "admin@example.com")

	companyID := createUnitViaAPI(t, router, token, `{"name":"Company"}`)
	operationsID := createUnitViaAPI(t, router, token, `{"name":"Operations","parentId":"`+companyID+`"}`)
	financeID := createUnitViaAPI(t, router, token, `{"name":"Finance","parentId":"`+companyID+`"}`)
	administrationID := createUnitViaAPI(t, router, token, `{"name":"Administration"}`)

	tree := orgUnitJSON(t, router, http.MethodGet, "/api/organization-units/tree", "", authCookieHeaders(token))
	if tree.Code != http.StatusOK {
		t.Fatalf("tree status = %d body=%s", tree.Code, tree.Body.String())
	}
	if !strings.Contains(tree.Body.String(), `"children":[]`) {
		t.Fatalf("tree body = %s, want leaf children arrays", tree.Body.String())
	}
	var body organizationUnitTestListResponse
	decodeJSON(t, tree, &body)
	if len(body.Units) != 2 {
		t.Fatalf("root count = %d body=%#v, want 2", len(body.Units), body)
	}
	if body.Units[0].ID != administrationID || body.Units[1].ID != companyID {
		t.Fatalf("root order = [%s %s], want administration then company", body.Units[0].ID, body.Units[1].ID)
	}
	if len(body.Units[1].Children) != 2 {
		t.Fatalf("company child count = %d body=%#v, want 2", len(body.Units[1].Children), body.Units[1])
	}
	if body.Units[1].Children[0].ID != financeID || body.Units[1].Children[1].ID != operationsID {
		t.Fatalf("company child order = [%s %s], want finance then operations", body.Units[1].Children[0].ID, body.Units[1].Children[1].ID)
	}

	archive := orgUnitJSON(t, router, http.MethodPost, "/api/organization-units/"+operationsID+"/archive", `{}`, authCookieHeaders(token))
	if archive.Code != http.StatusOK {
		t.Fatalf("archive status = %d body=%s", archive.Code, archive.Body.String())
	}

	tree = orgUnitJSON(t, router, http.MethodGet, "/api/organization-units/tree", "", authCookieHeaders(token))
	if tree.Code != http.StatusOK {
		t.Fatalf("tree after archive status = %d body=%s", tree.Code, tree.Body.String())
	}
	decodeJSON(t, tree, &body)
	if len(body.Units[1].Children) != 1 || body.Units[1].Children[0].ID != financeID {
		t.Fatalf("company children after archive = %#v, want only finance", body.Units[1].Children)
	}
}

func TestOrganizationUnitCreateAndUpdateValidation(t *testing.T) {
	router, db := newAuthTestRouter(t)
	token := loginSeededAdmin(t, router, db, "admin@example.com")

	blankName := orgUnitJSON(t, router, http.MethodPost, "/api/organization-units", `{"name":" "}`, authCookieHeaders(token))
	if blankName.Code != http.StatusBadRequest {
		t.Fatalf("blank name status = %d body=%s", blankName.Code, blankName.Body.String())
	}
	invalidCode := orgUnitJSON(t, router, http.MethodPost, "/api/organization-units", `{"name":"Finance","code":"FÎN"}`, authCookieHeaders(token))
	if invalidCode.Code != http.StatusBadRequest {
		t.Fatalf("invalid code status = %d body=%s", invalidCode.Code, invalidCode.Body.String())
	}
	invalidParent := orgUnitJSON(t, router, http.MethodPost, "/api/organization-units", `{"name":"Child","parentId":"not-a-uuid"}`, authCookieHeaders(token))
	if invalidParent.Code != http.StatusBadRequest {
		t.Fatalf("invalid parent status = %d body=%s", invalidParent.Code, invalidParent.Body.String())
	}
	missingParent := orgUnitJSON(t, router, http.MethodPost, "/api/organization-units", `{"name":"Child","parentId":"00000000-0000-0000-0000-000000000123"}`, authCookieHeaders(token))
	if missingParent.Code != http.StatusNotFound {
		t.Fatalf("missing parent status = %d body=%s", missingParent.Code, missingParent.Body.String())
	}

	financeID := createUnitViaAPI(t, router, token, `{"name":"Finance","code":"fin"}`)
	legalID := createUnitViaAPI(t, router, token, `{"name":"Legal","code":"legal"}`)
	duplicateCode := orgUnitJSON(t, router, http.MethodPost, "/api/organization-units", `{"name":"Finance Duplicate","code":" FIN "}`, authCookieHeaders(token))
	if duplicateCode.Code != http.StatusConflict {
		t.Fatalf("duplicate code status = %d body=%s", duplicateCode.Code, duplicateCode.Body.String())
	}
	updateDuplicateCode := orgUnitJSON(t, router, http.MethodPatch, "/api/organization-units/"+legalID, `{"name":"Legal","code":" FIN "}`, authCookieHeaders(token))
	if updateDuplicateCode.Code != http.StatusConflict {
		t.Fatalf("update duplicate code status = %d body=%s", updateDuplicateCode.Code, updateDuplicateCode.Body.String())
	}

	update := orgUnitJSON(t, router, http.MethodPatch, "/api/organization-units/"+financeID, `{"name":"Finance Operations","code":"fin_ops","description":" Accounts "}`, authCookieHeaders(token))
	if update.Code != http.StatusOK {
		t.Fatalf("update status = %d body=%s", update.Code, update.Body.String())
	}
	var updateBody organizationUnitTestResponse
	decodeJSON(t, update, &updateBody)
	if updateBody.Name != "Finance Operations" || updateBody.Code == nil || *updateBody.Code != "FIN_OPS" || updateBody.Description == nil || *updateBody.Description != "Accounts" {
		t.Fatalf("update body = %#v, want normalized values", updateBody)
	}
	preserveOptional := orgUnitJSON(t, router, http.MethodPatch, "/api/organization-units/"+financeID, `{"name":"Finance Renamed"}`, authCookieHeaders(token))
	if preserveOptional.Code != http.StatusOK {
		t.Fatalf("preserve optional status = %d body=%s", preserveOptional.Code, preserveOptional.Body.String())
	}
	updateBody = organizationUnitTestResponse{}
	decodeJSON(t, preserveOptional, &updateBody)
	if updateBody.Code == nil || *updateBody.Code != "FIN_OPS" || updateBody.Description == nil || *updateBody.Description != "Accounts" {
		t.Fatalf("preserve optional body = %#v, want existing code and description", updateBody)
	}
	clearOptional := orgUnitJSON(t, router, http.MethodPatch, "/api/organization-units/"+financeID, `{"name":"Finance Operations","code":" ","description":" "}`, authCookieHeaders(token))
	if clearOptional.Code != http.StatusOK {
		t.Fatalf("clear optional status = %d body=%s", clearOptional.Code, clearOptional.Body.String())
	}
	updateBody = organizationUnitTestResponse{}
	decodeJSON(t, clearOptional, &updateBody)
	if updateBody.Code != nil || updateBody.Description != nil {
		t.Fatalf("clear optional body = %#v, want nil code and description", updateBody)
	}

	updateBlank := orgUnitJSON(t, router, http.MethodPatch, "/api/organization-units/"+financeID, `{"name":" "}`, authCookieHeaders(token))
	if updateBlank.Code != http.StatusBadRequest {
		t.Fatalf("update blank status = %d body=%s", updateBlank.Code, updateBlank.Body.String())
	}
	updateMissing := orgUnitJSON(t, router, http.MethodPatch, "/api/organization-units/00000000-0000-0000-0000-000000000123", `{"name":"Missing"}`, authCookieHeaders(token))
	if updateMissing.Code != http.StatusNotFound {
		t.Fatalf("update missing status = %d body=%s", updateMissing.Code, updateMissing.Body.String())
	}

	archivedParentID := createUnitViaAPI(t, router, token, `{"name":"Archived Parent"}`)
	archiveParent := orgUnitJSON(t, router, http.MethodPost, "/api/organization-units/"+archivedParentID+"/archive", `{}`, authCookieHeaders(token))
	if archiveParent.Code != http.StatusOK {
		t.Fatalf("archive parent status = %d body=%s", archiveParent.Code, archiveParent.Body.String())
	}
	createUnderArchived := orgUnitJSON(t, router, http.MethodPost, "/api/organization-units", `{"name":"Child","parentId":"`+archivedParentID+`"}`, authCookieHeaders(token))
	if createUnderArchived.Code != http.StatusNotFound {
		t.Fatalf("create under archived parent status = %d body=%s", createUnderArchived.Code, createUnderArchived.Body.String())
	}
	updateArchived := orgUnitJSON(t, router, http.MethodPatch, "/api/organization-units/"+archivedParentID, `{"name":"Still Archived"}`, authCookieHeaders(token))
	if updateArchived.Code != http.StatusNotFound {
		t.Fatalf("update archived status = %d body=%s", updateArchived.Code, updateArchived.Body.String())
	}
}

func TestOrganizationUnitMoveValidationAndRootMove(t *testing.T) {
	router, db := newAuthTestRouter(t)
	token := loginSeededAdmin(t, router, db, "admin@example.com")

	rootID := createUnitViaAPI(t, router, token, `{"name":"Company"}`)
	childID := createUnitViaAPI(t, router, token, `{"name":"Finance","parentId":"`+rootID+`"}`)
	targetParentID := createUnitViaAPI(t, router, token, `{"name":"Operations"}`)

	moveToParent := orgUnitJSON(t, router, http.MethodPatch, "/api/organization-units/"+childID+"/parent", `{"parentId":"`+targetParentID+`"}`, authCookieHeaders(token))
	if moveToParent.Code != http.StatusOK {
		t.Fatalf("move to parent status = %d body=%s", moveToParent.Code, moveToParent.Body.String())
	}
	var moveBody organizationUnitTestResponse
	decodeJSON(t, moveToParent, &moveBody)
	if moveBody.ParentID == nil || *moveBody.ParentID != targetParentID {
		t.Fatalf("move body parentId = %#v, want %s", moveBody.ParentID, targetParentID)
	}

	moveToRoot := orgUnitJSON(t, router, http.MethodPatch, "/api/organization-units/"+childID+"/parent", `{"parentId":null}`, authCookieHeaders(token))
	if moveToRoot.Code != http.StatusOK {
		t.Fatalf("move to root status = %d body=%s", moveToRoot.Code, moveToRoot.Body.String())
	}
	moveBody = organizationUnitTestResponse{}
	decodeJSON(t, moveToRoot, &moveBody)
	if moveBody.ParentID != nil {
		t.Fatalf("move to root parentId = %#v, want nil", moveBody.ParentID)
	}

	moveUnderSelf := orgUnitJSON(t, router, http.MethodPatch, "/api/organization-units/"+targetParentID+"/parent", `{"parentId":"`+targetParentID+`"}`, authCookieHeaders(token))
	if moveUnderSelf.Code != http.StatusConflict {
		t.Fatalf("move under self status = %d body=%s", moveUnderSelf.Code, moveUnderSelf.Body.String())
	}
	moveWithMissingParentKey := orgUnitJSON(t, router, http.MethodPatch, "/api/organization-units/"+childID+"/parent", `{}`, authCookieHeaders(token))
	if moveWithMissingParentKey.Code != http.StatusBadRequest {
		t.Fatalf("move missing parent key status = %d body=%s", moveWithMissingParentKey.Code, moveWithMissingParentKey.Body.String())
	}
	moveUnderMissing := orgUnitJSON(t, router, http.MethodPatch, "/api/organization-units/"+childID+"/parent", `{"parentId":"00000000-0000-0000-0000-000000000123"}`, authCookieHeaders(token))
	if moveUnderMissing.Code != http.StatusNotFound {
		t.Fatalf("move under missing status = %d body=%s", moveUnderMissing.Code, moveUnderMissing.Body.String())
	}
	moveMissing := orgUnitJSON(t, router, http.MethodPatch, "/api/organization-units/00000000-0000-0000-0000-000000000123/parent", `{"parentId":null}`, authCookieHeaders(token))
	if moveMissing.Code != http.StatusNotFound {
		t.Fatalf("move missing status = %d body=%s", moveMissing.Code, moveMissing.Body.String())
	}

	archiveParent := orgUnitJSON(t, router, http.MethodPost, "/api/organization-units/"+targetParentID+"/archive", `{}`, authCookieHeaders(token))
	if archiveParent.Code != http.StatusOK {
		t.Fatalf("archive parent status = %d body=%s", archiveParent.Code, archiveParent.Body.String())
	}
	moveUnderArchived := orgUnitJSON(t, router, http.MethodPatch, "/api/organization-units/"+childID+"/parent", `{"parentId":"`+targetParentID+`"}`, authCookieHeaders(token))
	if moveUnderArchived.Code != http.StatusNotFound {
		t.Fatalf("move under archived status = %d body=%s", moveUnderArchived.Code, moveUnderArchived.Body.String())
	}
	moveArchived := orgUnitJSON(t, router, http.MethodPatch, "/api/organization-units/"+targetParentID+"/parent", `{"parentId":null}`, authCookieHeaders(token))
	if moveArchived.Code != http.StatusNotFound {
		t.Fatalf("move archived status = %d body=%s", moveArchived.Code, moveArchived.Body.String())
	}
	archiveMissing := orgUnitJSON(t, router, http.MethodPost, "/api/organization-units/00000000-0000-0000-0000-000000000123/archive", `{}`, authCookieHeaders(token))
	if archiveMissing.Code != http.StatusNotFound {
		t.Fatalf("archive missing status = %d body=%s", archiveMissing.Code, archiveMissing.Body.String())
	}
	archiveArchived := orgUnitJSON(t, router, http.MethodPost, "/api/organization-units/"+targetParentID+"/archive", `{}`, authCookieHeaders(token))
	if archiveArchived.Code != http.StatusNotFound {
		t.Fatalf("archive archived status = %d body=%s", archiveArchived.Code, archiveArchived.Body.String())
	}
}

type organizationUnitTestListResponse struct {
	Units []organizationUnitTestResponse `json:"units"`
}

type organizationUnitTestResponse struct {
	ID          string                         `json:"id"`
	ParentID    *string                        `json:"parentId"`
	Name        string                         `json:"name"`
	Code        *string                        `json:"code"`
	Description *string                        `json:"description"`
	ArchivedAt  *string                        `json:"archivedAt"`
	Children    []organizationUnitTestResponse `json:"children"`
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

func assignGlobalRoleByCode(t *testing.T, db *gorm.DB, userID string, roleCode string) {
	t.Helper()
	role := loadRoleByCode(t, db, roleCode)
	now := time.Now().UTC()
	assignment := rbac.UserRole{
		UserID:    userID,
		RoleID:    role.ID,
		ScopeType: rbac.ScopeGlobal,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := db.Create(&assignment).Error; err != nil {
		t.Fatalf("assign role %s: %v", roleCode, err)
	}
}

func authCookieHeaders(token string) map[string]string {
	return map[string]string{
		internalAPIHeader: testInternalToken,
		"Cookie":          authSessionCookieName + "=" + token,
	}
}

func orgUnitJSON(t *testing.T, router http.Handler, method string, path string, body string, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	return authJSON(t, router, method, path, body, headers)
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
