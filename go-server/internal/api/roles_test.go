package api

import (
	"net/http"
	"testing"

	"ai.ro/syncra/dms/internal/rbac"
	"gorm.io/gorm"
)

func TestRoleAPIRequiresRoleViewPermission(t *testing.T) {
	router, db := newAuthTestRouter(t)
	user := createVerifiedUser(t, db, "user@example.com", "password123")
	userToken := loginUser(t, router, user.Email, "password123")

	forbidden := authJSON(t, router, http.MethodGet, "/api/roles", "", authCookieHeaders(userToken))
	if forbidden.Code != http.StatusForbidden {
		t.Fatalf("status = %d body=%s, want forbidden", forbidden.Code, forbidden.Body.String())
	}

	adminToken := loginSeededAdmin(t, router, db, "admin@example.com")
	allowed := authJSON(t, router, http.MethodGet, "/api/roles", "", authCookieHeaders(adminToken))
	if allowed.Code != http.StatusOK {
		t.Fatalf("admin status = %d body=%s, want ok", allowed.Code, allowed.Body.String())
	}
	var body roleListTestResponse
	decodeJSON(t, allowed, &body)
	if len(body.Roles) == 0 {
		t.Fatal("roles list was empty")
	}
}

func TestRoleAPIAdminCanCreateCustomRoleAndRejectsDuplicateCode(t *testing.T) {
	router, db := newAuthTestRouter(t)
	adminToken := loginSeededAdmin(t, router, db, "admin@example.com")

	create := authJSON(t, router, http.MethodPost, "/api/roles", `{
		"name":"Finance Reviewer",
		"code":" Finance Reviewer ",
		"description":"Reviews finance documents"
	}`, authCookieHeaders(adminToken))
	if create.Code != http.StatusCreated {
		t.Fatalf("create status = %d body=%s, want created", create.Code, create.Body.String())
	}
	var created roleTestResponse
	decodeJSON(t, create, &created)
	if created.ID == "" || created.Code != "finance_reviewer" || created.Name != "Finance Reviewer" || created.IsSystem || !created.IsActive {
		t.Fatalf("created role = %#v", created)
	}

	duplicate := authJSON(t, router, http.MethodPost, "/api/roles", `{
		"name":"Finance Reviewer Duplicate",
		"code":"finance reviewer"
	}`, authCookieHeaders(adminToken))
	if duplicate.Code != http.StatusConflict {
		t.Fatalf("duplicate status = %d body=%s, want conflict", duplicate.Code, duplicate.Body.String())
	}
}

func TestRoleAPISystemRoleCodeAndDeleteAreProtected(t *testing.T) {
	router, db := newAuthTestRouter(t)
	adminToken := loginSeededAdmin(t, router, db, "admin@example.com")
	systemRole := loadRoleByCode(t, db, rbac.SystemAdministratorRoleCode)

	updateCode := authJSON(t, router, http.MethodPatch, "/api/roles/"+systemRole.ID, `{
		"name":"System Administrator",
		"code":"renamed_system_admin"
	}`, authCookieHeaders(adminToken))
	if updateCode.Code != http.StatusForbidden {
		t.Fatalf("system code update status = %d body=%s, want forbidden", updateCode.Code, updateCode.Body.String())
	}

	deleteSystem := authJSON(t, router, http.MethodDelete, "/api/roles/"+systemRole.ID, "", authCookieHeaders(adminToken))
	if deleteSystem.Code != http.StatusForbidden {
		t.Fatalf("system delete status = %d body=%s, want forbidden", deleteSystem.Code, deleteSystem.Body.String())
	}
}

func TestRoleAPIAdminCanAssignAndRemovePermissions(t *testing.T) {
	router, db := newAuthTestRouter(t)
	adminToken := loginSeededAdmin(t, router, db, "admin@example.com")
	role := createRoleViaAPI(t, router, adminToken, "Document Reviewer", "document_reviewer")
	permission := loadPermissionByCode(t, db, "user.view")

	assign := authJSON(t, router, http.MethodPost, "/api/roles/"+role.ID+"/permissions", `{
		"permissionId":"`+permission.ID+`"
	}`, authCookieHeaders(adminToken))
	if assign.Code != http.StatusCreated {
		t.Fatalf("assign permission status = %d body=%s, want created", assign.Code, assign.Body.String())
	}

	list := authJSON(t, router, http.MethodGet, "/api/roles/"+role.ID+"/permissions", "", authCookieHeaders(adminToken))
	if list.Code != http.StatusOK {
		t.Fatalf("list permissions status = %d body=%s, want ok", list.Code, list.Body.String())
	}
	var listBody permissionListTestResponse
	decodeJSON(t, list, &listBody)
	if len(listBody.Permissions) != 1 || listBody.Permissions[0].ID != permission.ID {
		t.Fatalf("role permissions = %#v, want assigned permission", listBody)
	}

	remove := authJSON(t, router, http.MethodDelete, "/api/roles/"+role.ID+"/permissions/"+permission.ID, "", authCookieHeaders(adminToken))
	if remove.Code != http.StatusOK {
		t.Fatalf("remove permission status = %d body=%s, want ok", remove.Code, remove.Body.String())
	}
	var count int64
	if err := db.Model(&rbac.RolePermission{}).Where("role_id = ? AND permission_id = ?", role.ID, permission.ID).Count(&count).Error; err != nil {
		t.Fatalf("count role permissions: %v", err)
	}
	if count != 0 {
		t.Fatalf("role permission count = %d, want 0", count)
	}
}

func createRoleViaAPI(t *testing.T, router http.Handler, token string, name string, code string) roleTestResponse {
	t.Helper()
	response := authJSON(t, router, http.MethodPost, "/api/roles", `{
		"name":"`+name+`",
		"code":"`+code+`"
	}`, authCookieHeaders(token))
	if response.Code != http.StatusCreated {
		t.Fatalf("create role status = %d body=%s, want created", response.Code, response.Body.String())
	}
	var role roleTestResponse
	decodeJSON(t, response, &role)
	return role
}

func loadPermissionByCode(t *testing.T, db *gorm.DB, code string) rbac.Permission {
	t.Helper()
	var permission rbac.Permission
	if err := db.First(&permission, "code = ?", code).Error; err != nil {
		t.Fatalf("load permission %s: %v", code, err)
	}
	return permission
}

type roleListTestResponse struct {
	Roles []roleTestResponse `json:"roles"`
}

type roleTestResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Code        string  `json:"code"`
	Description *string `json:"description"`
	IsSystem    bool    `json:"isSystem"`
	IsActive    bool    `json:"isActive"`
}
