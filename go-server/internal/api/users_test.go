package api

import (
	"net/http"
	"testing"
	"time"

	"ai.ro/syncra/dms/internal/auth"
	"ai.ro/syncra/dms/internal/orgunits"
	"ai.ro/syncra/dms/internal/rbac"
	"gorm.io/gorm"
)

func TestUserAPIRequiresUserViewPermission(t *testing.T) {
	router, db := newAuthTestRouter(t)

	unauthenticated := authJSON(t, router, http.MethodGet, "/api/users", "", map[string]string{
		internalAPIHeader: testInternalToken,
	})
	if unauthenticated.Code != http.StatusUnauthorized {
		t.Fatalf("unauthenticated status = %d body=%s, want unauthorized", unauthenticated.Code, unauthenticated.Body.String())
	}

	user := createVerifiedUser(t, db, "user@example.com", "password123")
	token := loginUser(t, router, user.Email, "password123")

	response := authJSON(t, router, http.MethodGet, "/api/users", "", authCookieHeaders(token))
	if response.Code != http.StatusForbidden {
		t.Fatalf("status = %d body=%s, want forbidden", response.Code, response.Body.String())
	}
}

func TestUserAPIAdminCanListAndCreateInvitedUser(t *testing.T) {
	router, db := newAuthTestRouter(t)
	token := loginSeededAdmin(t, router, db, "admin@example.com")

	list := authJSON(t, router, http.MethodGet, "/api/users", "", authCookieHeaders(token))
	if list.Code != http.StatusOK {
		t.Fatalf("list status = %d body=%s, want ok", list.Code, list.Body.String())
	}
	var listBody userListTestResponse
	decodeJSON(t, list, &listBody)
	if len(listBody.Users) != 1 || listBody.Users[0].Email != "admin@example.com" {
		t.Fatalf("list body = %#v, want admin user", listBody)
	}

	create := authJSON(t, router, http.MethodPost, "/api/users", `{
		"name":"Ada Lovelace",
		"email":"ada@example.com",
		"status":"invited"
	}`, authCookieHeaders(token))
	if create.Code != http.StatusCreated {
		t.Fatalf("create status = %d body=%s, want created", create.Code, create.Body.String())
	}
	var created userTestResponse
	decodeJSON(t, create, &created)
	if created.ID == "" || created.Name != "Ada Lovelace" || created.Email != "ada@example.com" || created.Status != "invited" || created.EmailVerified {
		t.Fatalf("created user = %#v", created)
	}
}

func TestUserAPIAdminCanUpdateProfileAndStatuses(t *testing.T) {
	router, db := newAuthTestRouter(t)
	adminToken := loginSeededAdmin(t, router, db, "admin@example.com")
	target := createVerifiedUser(t, db, "target@example.com", "password123")
	targetToken := loginUser(t, router, target.Email, "password123")

	update := authJSON(t, router, http.MethodPatch, "/api/users/"+target.ID, `{
		"name":"Grace Hopper",
		"jobTitle":"Admiral",
		"phone":"+1 555 0100"
	}`, authCookieHeaders(adminToken))
	if update.Code != http.StatusOK {
		t.Fatalf("update status = %d body=%s, want ok", update.Code, update.Body.String())
	}
	var updated userTestResponse
	decodeJSON(t, update, &updated)
	if updated.Name != "Grace Hopper" || updated.JobTitle == nil || *updated.JobTitle != "Admiral" || updated.Phone == nil || *updated.Phone != "+1 555 0100" {
		t.Fatalf("updated user = %#v", updated)
	}

	deactivate := authJSON(t, router, http.MethodPost, "/api/users/"+target.ID+"/deactivate", `{}`, authCookieHeaders(adminToken))
	if deactivate.Code != http.StatusOK {
		t.Fatalf("deactivate status = %d body=%s, want ok", deactivate.Code, deactivate.Body.String())
	}
	assertUserStatus(t, db, target.ID, "inactive")

	activate := authJSON(t, router, http.MethodPost, "/api/users/"+target.ID+"/activate", `{}`, authCookieHeaders(adminToken))
	if activate.Code != http.StatusOK {
		t.Fatalf("activate status = %d body=%s, want ok", activate.Code, activate.Body.String())
	}
	assertUserStatus(t, db, target.ID, "active")

	suspend := authJSON(t, router, http.MethodPost, "/api/users/"+target.ID+"/suspend", `{}`, authCookieHeaders(adminToken))
	if suspend.Code != http.StatusOK {
		t.Fatalf("suspend status = %d body=%s, want ok", suspend.Code, suspend.Body.String())
	}
	assertUserStatus(t, db, target.ID, "suspended")
	var sessionCount int64
	if err := db.Model(&auth.Session{}).Where("token = ?", targetToken).Count(&sessionCount).Error; err != nil {
		t.Fatalf("count target sessions: %v", err)
	}
	if sessionCount != 0 {
		t.Fatalf("target session count = %d, want 0", sessionCount)
	}
}

func TestUserAPIAdminCanAssignUnitsRolesAndGroups(t *testing.T) {
	router, db := newAuthTestRouter(t)
	adminToken := loginSeededAdmin(t, router, db, "admin@example.com")
	target := createVerifiedUser(t, db, "target@example.com", "password123")
	unit := createUserAPITestUnit(t, db, "Finance")
	group := createUserAPITestGroup(t, db, "Finance Approvers", "finance_approvers")
	role := loadRoleByCode(t, db, rbac.ViewerRoleCode)

	primaryUnit := authJSON(t, router, http.MethodPost, "/api/users/"+target.ID+"/primary-organization-unit", `{
		"organizationUnitId":"`+unit.ID+`"
	}`, authCookieHeaders(adminToken))
	if primaryUnit.Code != http.StatusOK {
		t.Fatalf("primary unit status = %d body=%s, want ok", primaryUnit.Code, primaryUnit.Body.String())
	}
	var primaryBody userTestResponse
	decodeJSON(t, primaryUnit, &primaryBody)
	if primaryBody.PrimaryOrganizationUnitID == nil || *primaryBody.PrimaryOrganizationUnitID != unit.ID {
		t.Fatalf("primary unit response = %#v", primaryBody)
	}

	assignRole := authJSON(t, router, http.MethodPost, "/api/users/"+target.ID+"/roles", `{
		"roleId":"`+role.ID+`",
		"scopeType":"organization_unit",
		"organizationUnitId":"`+unit.ID+`"
	}`, authCookieHeaders(adminToken))
	if assignRole.Code != http.StatusCreated {
		t.Fatalf("assign role status = %d body=%s, want created", assignRole.Code, assignRole.Body.String())
	}
	var roleBody userRoleAssignmentTestResponse
	decodeJSON(t, assignRole, &roleBody)
	if roleBody.ID == "" || roleBody.UserID != target.ID || roleBody.RoleID != role.ID || roleBody.ScopeType != string(rbac.ScopeOrganizationUnit) || roleBody.OrganizationUnitID == nil || *roleBody.OrganizationUnitID != unit.ID {
		t.Fatalf("role assignment response = %#v", roleBody)
	}

	removeRole := authJSON(t, router, http.MethodDelete, "/api/users/"+target.ID+"/roles/"+roleBody.ID, "", authCookieHeaders(adminToken))
	if removeRole.Code != http.StatusOK {
		t.Fatalf("remove role status = %d body=%s, want ok", removeRole.Code, removeRole.Body.String())
	}
	var roleCount int64
	if err := db.Model(&rbac.UserRole{}).Where("id = ?", roleBody.ID).Count(&roleCount).Error; err != nil {
		t.Fatalf("count removed user role: %v", err)
	}
	if roleCount != 0 {
		t.Fatalf("role assignment count = %d, want 0", roleCount)
	}

	assignGroup := authJSON(t, router, http.MethodPost, "/api/users/"+target.ID+"/groups", `{
		"groupId":"`+group.ID+`"
	}`, authCookieHeaders(adminToken))
	if assignGroup.Code != http.StatusCreated {
		t.Fatalf("assign group status = %d body=%s, want created", assignGroup.Code, assignGroup.Body.String())
	}

	removeGroup := authJSON(t, router, http.MethodDelete, "/api/users/"+target.ID+"/groups/"+group.ID, "", authCookieHeaders(adminToken))
	if removeGroup.Code != http.StatusOK {
		t.Fatalf("remove group status = %d body=%s, want ok", removeGroup.Code, removeGroup.Body.String())
	}
	var groupCount int64
	if err := db.Model(&rbac.GroupUser{}).Where("user_id = ? AND group_id = ?", target.ID, group.ID).Count(&groupCount).Error; err != nil {
		t.Fatalf("count removed group user: %v", err)
	}
	if groupCount != 0 {
		t.Fatalf("group assignment count = %d, want 0", groupCount)
	}
}

func loginSeededAdmin(t *testing.T, router http.Handler, db *gorm.DB, email string) string {
	t.Helper()
	admin := createAdminUser(t, db, email, "password123")
	if err := rbac.BootstrapLegacyAdmins(db); err != nil {
		t.Fatalf("bootstrap legacy admin: %v", err)
	}
	return loginUser(t, router, admin.Email, "password123")
}

func createUserAPITestUnit(t *testing.T, db *gorm.DB, name string) orgunits.Unit {
	t.Helper()
	now := time.Now().UTC()
	unit := orgunits.Unit{Name: name, CreatedAt: now, UpdatedAt: now}
	if err := db.Create(&unit).Error; err != nil {
		t.Fatalf("create organization unit: %v", err)
	}
	return unit
}

func createUserAPITestGroup(t *testing.T, db *gorm.DB, name string, code string) rbac.Group {
	t.Helper()
	now := time.Now().UTC()
	group := rbac.Group{Name: name, Code: code, IsActive: true, CreatedAt: now, UpdatedAt: now}
	if err := db.Create(&group).Error; err != nil {
		t.Fatalf("create group: %v", err)
	}
	return group
}

func loadRoleByCode(t *testing.T, db *gorm.DB, code string) rbac.Role {
	t.Helper()
	var role rbac.Role
	if err := db.First(&role, "code = ?", code).Error; err != nil {
		t.Fatalf("load role %s: %v", code, err)
	}
	return role
}

func assertUserStatus(t *testing.T, db *gorm.DB, userID string, status string) {
	t.Helper()
	var user auth.User
	if err := db.First(&user, "id = ?", userID).Error; err != nil {
		t.Fatalf("load user: %v", err)
	}
	if user.Status != status {
		t.Fatalf("user status = %q, want %q", user.Status, status)
	}
}

type userListTestResponse struct {
	Users []userTestResponse `json:"users"`
}

type userTestResponse struct {
	ID                        string  `json:"id"`
	Name                      string  `json:"name"`
	Email                     string  `json:"email"`
	EmailVerified             bool    `json:"emailVerified"`
	Status                    string  `json:"status"`
	Role                      string  `json:"role"`
	PrimaryOrganizationUnitID *string `json:"primaryOrganizationUnitId"`
	JobTitle                  *string `json:"jobTitle"`
	Phone                     *string `json:"phone"`
}

type userRoleAssignmentTestResponse struct {
	ID                 string  `json:"id"`
	UserID             string  `json:"userId"`
	RoleID             string  `json:"roleId"`
	ScopeType          string  `json:"scopeType"`
	OrganizationUnitID *string `json:"organizationUnitId"`
}
