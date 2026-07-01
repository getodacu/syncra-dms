package api

import (
	"net/http"
	"testing"

	"ai.ro/syncra/dms/internal/rbac"
	"gorm.io/gorm"
)

func TestGroupAPIRequiresGroupViewPermission(t *testing.T) {
	router, db := newAuthTestRouter(t)
	user := createVerifiedUser(t, db, "user@example.com", "password123")
	userToken := loginUser(t, router, user.Email, "password123")

	forbidden := authJSON(t, router, http.MethodGet, "/api/groups", "", authCookieHeaders(userToken))
	if forbidden.Code != http.StatusForbidden {
		t.Fatalf("status = %d body=%s, want forbidden", forbidden.Code, forbidden.Body.String())
	}

	adminToken := loginSeededAdmin(t, router, db, "admin@example.com")
	allowed := authJSON(t, router, http.MethodGet, "/api/groups", "", authCookieHeaders(adminToken))
	if allowed.Code != http.StatusOK {
		t.Fatalf("admin status = %d body=%s, want ok", allowed.Code, allowed.Body.String())
	}
}

func TestGroupAPIAdminCanCreateUpdateAndRejectDuplicateCode(t *testing.T) {
	router, db := newAuthTestRouter(t)
	adminToken := loginSeededAdmin(t, router, db, "admin@example.com")

	create := authJSON(t, router, http.MethodPost, "/api/groups", `{
		"name":"Finance Approvers",
		"code":" Finance Approvers ",
		"description":"Approves finance documents"
	}`, authCookieHeaders(adminToken))
	if create.Code != http.StatusCreated {
		t.Fatalf("create status = %d body=%s, want created", create.Code, create.Body.String())
	}
	var created groupTestResponse
	decodeJSON(t, create, &created)
	if created.ID == "" || created.Code != "finance_approvers" || created.Name != "Finance Approvers" || !created.IsActive {
		t.Fatalf("created group = %#v", created)
	}

	duplicate := authJSON(t, router, http.MethodPost, "/api/groups", `{
		"name":"Finance Duplicate",
		"code":"finance approvers"
	}`, authCookieHeaders(adminToken))
	if duplicate.Code != http.StatusConflict {
		t.Fatalf("duplicate status = %d body=%s, want conflict", duplicate.Code, duplicate.Body.String())
	}

	update := authJSON(t, router, http.MethodPatch, "/api/groups/"+created.ID, `{
		"name":"Finance Reviewers",
		"isActive":false
	}`, authCookieHeaders(adminToken))
	if update.Code != http.StatusOK {
		t.Fatalf("update status = %d body=%s, want ok", update.Code, update.Body.String())
	}
	var updated groupTestResponse
	decodeJSON(t, update, &updated)
	if updated.Name != "Finance Reviewers" || updated.IsActive {
		t.Fatalf("updated group = %#v", updated)
	}
}

func TestGroupAPIAdminCanAddAndRemoveUsers(t *testing.T) {
	router, db := newAuthTestRouter(t)
	adminToken := loginSeededAdmin(t, router, db, "admin@example.com")
	group := createGroupViaAPI(t, router, adminToken, "Legal Reviewers", "legal_reviewers")
	user := createVerifiedUser(t, db, "member@example.com", "password123")

	add := authJSON(t, router, http.MethodPost, "/api/groups/"+group.ID+"/users", `{
		"userId":"`+user.ID+`"
	}`, authCookieHeaders(adminToken))
	if add.Code != http.StatusCreated {
		t.Fatalf("add user status = %d body=%s, want created", add.Code, add.Body.String())
	}

	remove := authJSON(t, router, http.MethodDelete, "/api/groups/"+group.ID+"/users/"+user.ID, "", authCookieHeaders(adminToken))
	if remove.Code != http.StatusOK {
		t.Fatalf("remove user status = %d body=%s, want ok", remove.Code, remove.Body.String())
	}
	var count int64
	if err := db.Model(&rbac.GroupUser{}).Where("group_id = ? AND user_id = ?", group.ID, user.ID).Count(&count).Error; err != nil {
		t.Fatalf("count group user: %v", err)
	}
	if count != 0 {
		t.Fatalf("group user count = %d, want 0", count)
	}
}

func TestGroupAPIAdminCanAssignAndRemoveRoles(t *testing.T) {
	router, db := newAuthTestRouter(t)
	adminToken := loginSeededAdmin(t, router, db, "admin@example.com")
	group := createGroupViaAPI(t, router, adminToken, "Finance Managers", "finance_managers")
	role := loadRoleByCode(t, db, rbac.ViewerRoleCode)
	unit := createUserAPITestUnit(t, db, "Finance")

	assign := authJSON(t, router, http.MethodPost, "/api/groups/"+group.ID+"/roles", `{
		"roleId":"`+role.ID+`",
		"scopeType":"organization_unit",
		"organizationUnitId":"`+unit.ID+`"
	}`, authCookieHeaders(adminToken))
	if assign.Code != http.StatusCreated {
		t.Fatalf("assign role status = %d body=%s, want created", assign.Code, assign.Body.String())
	}
	var assignment groupRoleAssignmentTestResponse
	decodeJSON(t, assign, &assignment)
	if assignment.ID == "" || assignment.GroupID != group.ID || assignment.RoleID != role.ID || assignment.OrganizationUnitID == nil || *assignment.OrganizationUnitID != unit.ID {
		t.Fatalf("assignment response = %#v", assignment)
	}

	remove := authJSON(t, router, http.MethodDelete, "/api/groups/"+group.ID+"/roles/"+assignment.ID, "", authCookieHeaders(adminToken))
	if remove.Code != http.StatusOK {
		t.Fatalf("remove role status = %d body=%s, want ok", remove.Code, remove.Body.String())
	}
	var count int64
	if err := db.Model(&rbac.GroupRole{}).Where("id = ?", assignment.ID).Count(&count).Error; err != nil {
		t.Fatalf("count group role: %v", err)
	}
	if count != 0 {
		t.Fatalf("group role count = %d, want 0", count)
	}
}

func TestGroupAPIDeleteRejectsGroupsWithMembersOrRoles(t *testing.T) {
	router, db := newAuthTestRouter(t)
	adminToken := loginSeededAdmin(t, router, db, "admin@example.com")
	memberGroup := createGroupViaAPI(t, router, adminToken, "Member Group", "member_group")
	roleGroup := createGroupViaAPI(t, router, adminToken, "Role Group", "role_group")
	emptyGroup := createGroupViaAPI(t, router, adminToken, "Empty Group", "empty_group")
	user := createVerifiedUser(t, db, "member@example.com", "password123")
	role := loadRoleByCode(t, db, rbac.ViewerRoleCode)
	if err := db.Create(&rbac.GroupUser{GroupID: memberGroup.ID, UserID: user.ID}).Error; err != nil {
		t.Fatalf("create group user: %v", err)
	}
	if err := db.Create(&rbac.GroupRole{GroupID: roleGroup.ID, RoleID: role.ID, ScopeType: rbac.ScopeGlobal}).Error; err != nil {
		t.Fatalf("create group role: %v", err)
	}

	deleteMemberGroup := authJSON(t, router, http.MethodDelete, "/api/groups/"+memberGroup.ID, "", authCookieHeaders(adminToken))
	if deleteMemberGroup.Code != http.StatusConflict {
		t.Fatalf("member group delete status = %d body=%s, want conflict", deleteMemberGroup.Code, deleteMemberGroup.Body.String())
	}
	deleteRoleGroup := authJSON(t, router, http.MethodDelete, "/api/groups/"+roleGroup.ID, "", authCookieHeaders(adminToken))
	if deleteRoleGroup.Code != http.StatusConflict {
		t.Fatalf("role group delete status = %d body=%s, want conflict", deleteRoleGroup.Code, deleteRoleGroup.Body.String())
	}
	deleteEmptyGroup := authJSON(t, router, http.MethodDelete, "/api/groups/"+emptyGroup.ID, "", authCookieHeaders(adminToken))
	if deleteEmptyGroup.Code != http.StatusOK {
		t.Fatalf("empty group delete status = %d body=%s, want ok", deleteEmptyGroup.Code, deleteEmptyGroup.Body.String())
	}
}

func createGroupViaAPI(t *testing.T, router http.Handler, token string, name string, code string) groupTestResponse {
	t.Helper()
	response := authJSON(t, router, http.MethodPost, "/api/groups", `{
		"name":"`+name+`",
		"code":"`+code+`"
	}`, authCookieHeaders(token))
	if response.Code != http.StatusCreated {
		t.Fatalf("create group status = %d body=%s, want created", response.Code, response.Body.String())
	}
	var group groupTestResponse
	decodeJSON(t, response, &group)
	return group
}

func loadGroupByCode(t *testing.T, db *gorm.DB, code string) rbac.Group {
	t.Helper()
	var group rbac.Group
	if err := db.First(&group, "code = ?", code).Error; err != nil {
		t.Fatalf("load group %s: %v", code, err)
	}
	return group
}

type groupTestResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Code        string  `json:"code"`
	Description *string `json:"description"`
	IsActive    bool    `json:"isActive"`
}

type groupRoleAssignmentTestResponse struct {
	ID                 string  `json:"id"`
	GroupID            string  `json:"groupId"`
	RoleID             string  `json:"roleId"`
	ScopeType          string  `json:"scopeType"`
	OrganizationUnitID *string `json:"organizationUnitId"`
}
