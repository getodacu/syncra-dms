package api

import (
	"errors"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"ai.ro/syncra/dms/internal/auth"
	"ai.ro/syncra/dms/internal/orgunits"
	"ai.ro/syncra/dms/internal/rbac"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userHandler struct {
	db   *gorm.DB
	auth *authHandler
}

type createUserRequest struct {
	Name                      string  `json:"name"`
	Email                     string  `json:"email"`
	Status                    *string `json:"status"`
	PrimaryOrganizationUnitID *string `json:"primaryOrganizationUnitId"`
	ManagerUserID             *string `json:"managerUserId"`
	JobTitle                  *string `json:"jobTitle"`
	Phone                     *string `json:"phone"`
}

type updateUserRequest struct {
	Name          *string `json:"name"`
	ManagerUserID *string `json:"managerUserId"`
	JobTitle      *string `json:"jobTitle"`
	Phone         *string `json:"phone"`
}

type setPrimaryOrganizationUnitRequest struct {
	OrganizationUnitID *string `json:"organizationUnitId"`
}

type assignUserRoleRequest struct {
	RoleID             string  `json:"roleId"`
	ScopeType          string  `json:"scopeType"`
	OrganizationUnitID *string `json:"organizationUnitId"`
}

type assignUserGroupRequest struct {
	GroupID string `json:"groupId"`
}

type userResponse struct {
	ID                        string  `json:"id"`
	Name                      string  `json:"name"`
	Email                     string  `json:"email"`
	EmailVerified             bool    `json:"emailVerified"`
	Image                     *string `json:"image,omitempty"`
	PreferredLanguage         string  `json:"preferredLanguage"`
	Role                      string  `json:"role"`
	Status                    string  `json:"status"`
	PrimaryOrganizationUnitID *string `json:"primaryOrganizationUnitId,omitempty"`
	ManagerUserID             *string `json:"managerUserId,omitempty"`
	JobTitle                  *string `json:"jobTitle,omitempty"`
	Phone                     *string `json:"phone,omitempty"`
	LastLoginAt               *string `json:"lastLoginAt,omitempty"`
	DeletedAt                 *string `json:"deletedAt,omitempty"`
	CreatedAt                 string  `json:"createdAt"`
	UpdatedAt                 string  `json:"updatedAt"`
}

type userListResponse struct {
	Users []userResponse `json:"users"`
}

type userRoleAssignmentResponse struct {
	ID                 string  `json:"id"`
	UserID             string  `json:"userId"`
	RoleID             string  `json:"roleId"`
	ScopeType          string  `json:"scopeType"`
	OrganizationUnitID *string `json:"organizationUnitId,omitempty"`
	CreatedAt          string  `json:"createdAt"`
	UpdatedAt          string  `json:"updatedAt"`
}

type userGroupAssignmentResponse struct {
	UserID    string `json:"userId"`
	GroupID   string `json:"groupId"`
	CreatedAt string `json:"createdAt"`
}

func newUserHandler(options RouterOptions, auth *authHandler) *userHandler {
	return &userHandler{db: options.DB, auth: auth}
}

func (h *userHandler) list(c *gin.Context) {
	if _, ok := requirePermission(c, h.auth, "user.view", nil); !ok {
		return
	}
	var users []auth.User
	if err := h.db.WithContext(c.Request.Context()).
		Where("deleted_at IS NULL").
		Order("email asc, id asc").
		Find(&users).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to list users")
		return
	}
	out := make([]userResponse, 0, len(users))
	for _, user := range users {
		out = append(out, userResponseFromModel(user))
	}
	c.JSON(http.StatusOK, userListResponse{Users: out})
}

func (h *userHandler) get(c *gin.Context) {
	if _, ok := requirePermission(c, h.auth, "user.view", nil); !ok {
		return
	}
	id, ok := parseUserID(c, c.Param("id"))
	if !ok {
		return
	}
	user, ok := h.loadVisibleUser(c, id)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, userResponseFromModel(user))
}

func (h *userHandler) create(c *gin.Context) {
	if _, ok := requirePermission(c, h.auth, "user.create", nil); !ok {
		return
	}
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	name, ok := normalizeUserName(c, req.Name)
	if !ok {
		return
	}
	email := normalizeEmail(req.Email)
	if !validEmail(email) {
		writeError(c, http.StatusBadRequest, "valid email is required")
		return
	}
	status := string(rbac.UserStatusInvited)
	if req.Status != nil {
		status = strings.TrimSpace(*req.Status)
	}
	if !validCreatableUserStatus(status) {
		writeError(c, http.StatusBadRequest, "status is invalid")
		return
	}
	primaryUnitID, ok := h.validOptionalOrganizationUnitID(c, req.PrimaryOrganizationUnitID)
	if !ok {
		return
	}
	managerID, ok := h.validOptionalManagerID(c, req.ManagerUserID, "")
	if !ok {
		return
	}
	now := time.Now().UTC()
	user := auth.User{
		Name:                      name,
		Email:                     email,
		EmailVerified:             status == string(rbac.UserStatusActive),
		Role:                      auth.UserRoleUser,
		Status:                    status,
		PrimaryOrganizationUnitID: primaryUnitID,
		ManagerUserID:             managerID,
		JobTitle:                  normalizeOptionalUserText(req.JobTitle),
		Phone:                     normalizeOptionalUserText(req.Phone),
		CreatedAt:                 now,
		UpdatedAt:                 now,
	}
	if err := h.db.WithContext(c.Request.Context()).Create(&user).Error; err != nil {
		writeUserMutationError(c, err, "failed to create user")
		return
	}
	c.JSON(http.StatusCreated, userResponseFromModel(user))
}

func (h *userHandler) update(c *gin.Context) {
	if _, ok := requirePermission(c, h.auth, "user.update", nil); !ok {
		return
	}
	id, ok := parseUserID(c, c.Param("id"))
	if !ok {
		return
	}
	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	updates := map[string]any{"updated_at": time.Now().UTC()}
	if req.Name != nil {
		name, ok := normalizeUserName(c, *req.Name)
		if !ok {
			return
		}
		updates["name"] = name
	}
	if req.ManagerUserID != nil {
		managerID, ok := h.validOptionalManagerID(c, req.ManagerUserID, id)
		if !ok {
			return
		}
		updates["manager_user_id"] = nullableStringValue(managerID)
	}
	if req.JobTitle != nil {
		updates["job_title"] = nullableStringValue(normalizeOptionalUserText(req.JobTitle))
	}
	if req.Phone != nil {
		updates["phone"] = nullableStringValue(normalizeOptionalUserText(req.Phone))
	}
	user, ok := h.updateUser(c, id, updates)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, userResponseFromModel(user))
}

func (h *userHandler) activate(c *gin.Context) {
	if _, ok := requirePermission(c, h.auth, "user.activate", nil); !ok {
		return
	}
	h.setUserStatus(c, c.Param("id"), string(rbac.UserStatusActive), false, false)
}

func (h *userHandler) deactivate(c *gin.Context) {
	if _, ok := requirePermission(c, h.auth, "user.update", nil); !ok {
		return
	}
	h.setUserStatus(c, c.Param("id"), string(rbac.UserStatusInactive), false, false)
}

func (h *userHandler) suspend(c *gin.Context) {
	if _, ok := requirePermission(c, h.auth, "user.suspend", nil); !ok {
		return
	}
	h.setUserStatus(c, c.Param("id"), string(rbac.UserStatusSuspended), false, true)
}

func (h *userHandler) softDelete(c *gin.Context) {
	if _, ok := requirePermission(c, h.auth, "user.delete", nil); !ok {
		return
	}
	h.setUserStatus(c, c.Param("id"), string(rbac.UserStatusDeleted), true, true)
}

func (h *userHandler) setPrimaryOrganizationUnit(c *gin.Context) {
	if _, ok := requirePermission(c, h.auth, "user.assign_unit", nil); !ok {
		return
	}
	id, ok := parseUserID(c, c.Param("id"))
	if !ok {
		return
	}
	var req setPrimaryOrganizationUnitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	organizationUnitID, ok := h.validOptionalOrganizationUnitID(c, req.OrganizationUnitID)
	if !ok {
		return
	}
	user, ok := h.updateUser(c, id, map[string]any{
		"primary_organization_unit_id": nullableStringValue(organizationUnitID),
		"updated_at":                   time.Now().UTC(),
	})
	if !ok {
		return
	}
	c.JSON(http.StatusOK, userResponseFromModel(user))
}

func (h *userHandler) assignRole(c *gin.Context) {
	if _, ok := requirePermission(c, h.auth, "user.assign_role", nil); !ok {
		return
	}
	userID, ok := parseUserID(c, c.Param("id"))
	if !ok {
		return
	}
	if _, ok := h.loadVisibleUser(c, userID); !ok {
		return
	}
	var req assignUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	roleID, ok := parseUUIDValue(c, req.RoleID, "invalid role id")
	if !ok {
		return
	}
	if !h.roleExists(c, roleID) {
		return
	}
	scope := rbac.ScopeType(strings.TrimSpace(req.ScopeType))
	if scope == "" {
		scope = rbac.ScopeGlobal
	}
	organizationUnitID, ok := h.validScopedOrganizationUnitID(c, scope, req.OrganizationUnitID)
	if !ok {
		return
	}
	now := time.Now().UTC()
	assignment := rbac.UserRole{
		UserID:             userID,
		RoleID:             roleID,
		ScopeType:          scope,
		OrganizationUnitID: organizationUnitID,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&assignment).Error; err != nil {
			return err
		}
		query := tx.Where("user_id = ? AND role_id = ? AND scope_type = ?", userID, roleID, scope)
		if organizationUnitID == nil {
			query = query.Where("organization_unit_id IS NULL")
		} else {
			query = query.Where("organization_unit_id = ?", *organizationUnitID)
		}
		return query.First(&assignment).Error
	})
	if err != nil {
		writeUserMutationError(c, err, "failed to assign user role")
		return
	}
	c.JSON(http.StatusCreated, userRoleAssignmentResponseFromModel(assignment))
}

func (h *userHandler) removeRole(c *gin.Context) {
	if _, ok := requirePermission(c, h.auth, "user.assign_role", nil); !ok {
		return
	}
	userID, ok := parseUserID(c, c.Param("id"))
	if !ok {
		return
	}
	assignmentID, ok := parseUUIDValue(c, c.Param("assignmentId"), "invalid role assignment id")
	if !ok {
		return
	}
	result := h.db.WithContext(c.Request.Context()).Where("id = ? AND user_id = ?", assignmentID, userID).Delete(&rbac.UserRole{})
	if result.Error != nil {
		writeError(c, http.StatusInternalServerError, "failed to remove user role")
		return
	}
	if result.RowsAffected == 0 {
		writeError(c, http.StatusNotFound, "user role assignment not found")
		return
	}
	c.JSON(http.StatusOK, okResponse{OK: true})
}

func (h *userHandler) addGroup(c *gin.Context) {
	if _, ok := requirePermission(c, h.auth, "user.assign_group", nil); !ok {
		return
	}
	userID, ok := parseUserID(c, c.Param("id"))
	if !ok {
		return
	}
	if _, ok := h.loadVisibleUser(c, userID); !ok {
		return
	}
	var req assignUserGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	groupID, ok := parseUUIDValue(c, req.GroupID, "invalid group id")
	if !ok {
		return
	}
	if !h.groupExists(c, groupID) {
		return
	}
	now := time.Now().UTC()
	assignment := rbac.GroupUser{
		UserID:    userID,
		GroupID:   groupID,
		CreatedAt: now,
	}
	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&assignment).Error; err != nil {
			return err
		}
		return tx.Where("user_id = ? AND group_id = ?", userID, groupID).First(&assignment).Error
	})
	if err != nil {
		writeUserMutationError(c, err, "failed to assign user group")
		return
	}
	c.JSON(http.StatusCreated, userGroupAssignmentResponseFromModel(assignment))
}

func (h *userHandler) removeGroup(c *gin.Context) {
	if _, ok := requirePermission(c, h.auth, "user.assign_group", nil); !ok {
		return
	}
	userID, ok := parseUserID(c, c.Param("id"))
	if !ok {
		return
	}
	groupID, ok := parseUUIDValue(c, c.Param("groupId"), "invalid group id")
	if !ok {
		return
	}
	result := h.db.WithContext(c.Request.Context()).Where("user_id = ? AND group_id = ?", userID, groupID).Delete(&rbac.GroupUser{})
	if result.Error != nil {
		writeError(c, http.StatusInternalServerError, "failed to remove user group")
		return
	}
	if result.RowsAffected == 0 {
		writeError(c, http.StatusNotFound, "user group assignment not found")
		return
	}
	c.JSON(http.StatusOK, okResponse{OK: true})
}

func (h *userHandler) setUserStatus(c *gin.Context, rawID string, status string, deleted bool, deleteSessions bool) {
	id, ok := parseUserID(c, rawID)
	if !ok {
		return
	}
	now := time.Now().UTC()
	updates := map[string]any{
		"status":     status,
		"updated_at": now,
	}
	if deleted {
		updates["deleted_at"] = now
	}
	var user auth.User
	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&auth.User{}).Where("id = ? AND deleted_at IS NULL", id).Updates(updates)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		if deleteSessions {
			if err := tx.Where("user_id = ?", id).Delete(&auth.Session{}).Error; err != nil {
				return err
			}
		}
		return tx.First(&user, "id = ?", id).Error
	})
	if err != nil {
		writeUserMutationError(c, err, "failed to update user status")
		return
	}
	c.JSON(http.StatusOK, userResponseFromModel(user))
}

func (h *userHandler) updateUser(c *gin.Context, id string, updates map[string]any) (auth.User, bool) {
	var user auth.User
	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&auth.User{}).Where("id = ? AND deleted_at IS NULL", id).Updates(updates)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return tx.First(&user, "id = ?", id).Error
	})
	if err != nil {
		writeUserMutationError(c, err, "failed to update user")
		return auth.User{}, false
	}
	return user, true
}

func (h *userHandler) loadVisibleUser(c *gin.Context, id string) (auth.User, bool) {
	var user auth.User
	err := h.db.WithContext(c.Request.Context()).Where("id = ? AND deleted_at IS NULL", id).First(&user).Error
	if err != nil {
		writeUserMutationError(c, err, "failed to load user")
		return auth.User{}, false
	}
	return user, true
}

func (h *userHandler) validOptionalOrganizationUnitID(c *gin.Context, raw *string) (*string, bool) {
	if raw == nil || strings.TrimSpace(*raw) == "" {
		return nil, true
	}
	id, ok := parseUUIDValue(c, *raw, "invalid organization unit id")
	if !ok {
		return nil, false
	}
	var count int64
	if err := h.db.WithContext(c.Request.Context()).Model(&orgunits.Unit{}).Where("id = ? AND archived_at IS NULL", id).Count(&count).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to validate organization unit")
		return nil, false
	}
	if count == 0 {
		writeError(c, http.StatusNotFound, "organization unit not found")
		return nil, false
	}
	return &id, true
}

func (h *userHandler) validOptionalManagerID(c *gin.Context, raw *string, userID string) (*string, bool) {
	if raw == nil || strings.TrimSpace(*raw) == "" {
		return nil, true
	}
	id, ok := parseUUIDValue(c, *raw, "invalid manager user id")
	if !ok {
		return nil, false
	}
	if id == userID {
		writeError(c, http.StatusConflict, "user cannot manage themselves")
		return nil, false
	}
	var count int64
	if err := h.db.WithContext(c.Request.Context()).Model(&auth.User{}).Where("id = ? AND deleted_at IS NULL", id).Count(&count).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to validate manager user")
		return nil, false
	}
	if count == 0 {
		writeError(c, http.StatusNotFound, "manager user not found")
		return nil, false
	}
	return &id, true
}

func (h *userHandler) validScopedOrganizationUnitID(c *gin.Context, scope rbac.ScopeType, raw *string) (*string, bool) {
	if !scope.Valid() {
		writeError(c, http.StatusBadRequest, "scope type is invalid")
		return nil, false
	}
	if scope == rbac.ScopeGlobal {
		if raw != nil && strings.TrimSpace(*raw) != "" {
			writeError(c, http.StatusBadRequest, "global scope must not include an organization unit")
			return nil, false
		}
		return nil, true
	}
	if raw == nil || strings.TrimSpace(*raw) == "" {
		writeError(c, http.StatusBadRequest, "organization unit id is required")
		return nil, false
	}
	return h.validOptionalOrganizationUnitID(c, raw)
}

func (h *userHandler) roleExists(c *gin.Context, roleID string) bool {
	var count int64
	if err := h.db.WithContext(c.Request.Context()).Model(&rbac.Role{}).Where("id = ? AND is_active = ?", roleID, true).Count(&count).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to validate role")
		return false
	}
	if count == 0 {
		writeError(c, http.StatusNotFound, "role not found")
		return false
	}
	return true
}

func (h *userHandler) groupExists(c *gin.Context, groupID string) bool {
	var count int64
	if err := h.db.WithContext(c.Request.Context()).Model(&rbac.Group{}).Where("id = ? AND is_active = ?", groupID, true).Count(&count).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to validate group")
		return false
	}
	if count == 0 {
		writeError(c, http.StatusNotFound, "group not found")
		return false
	}
	return true
}

func parseUserID(c *gin.Context, raw string) (string, bool) {
	return parseUUIDValue(c, raw, "invalid user id")
}

func parseUUIDValue(c *gin.Context, raw string, message string) (string, bool) {
	id, err := uuid.Parse(strings.TrimSpace(raw))
	if err != nil || id == uuid.Nil {
		writeError(c, http.StatusBadRequest, message)
		return "", false
	}
	return id.String(), true
}

func normalizeUserName(c *gin.Context, raw string) (string, bool) {
	name := strings.TrimSpace(raw)
	if name == "" {
		writeError(c, http.StatusBadRequest, "name is required")
		return "", false
	}
	if utf8.RuneCountInString(name) > maxAuthNameCharacters {
		writeError(c, http.StatusBadRequest, "name must be at most 255 characters")
		return "", false
	}
	return name, true
}

func normalizeOptionalUserText(raw *string) *string {
	if raw == nil {
		return nil
	}
	value := strings.TrimSpace(*raw)
	if value == "" {
		return nil
	}
	return &value
}

func validCreatableUserStatus(status string) bool {
	switch rbac.UserStatus(status) {
	case rbac.UserStatusInvited, rbac.UserStatusActive, rbac.UserStatusInactive, rbac.UserStatusSuspended:
		return true
	default:
		return false
	}
}

func userResponseFromModel(user auth.User) userResponse {
	var lastLoginAt *string
	if user.LastLoginAt != nil {
		value := user.LastLoginAt.UTC().Format(time.RFC3339Nano)
		lastLoginAt = &value
	}
	var deletedAt *string
	if user.DeletedAt != nil {
		value := user.DeletedAt.UTC().Format(time.RFC3339Nano)
		deletedAt = &value
	}
	return userResponse{
		ID:                        user.ID,
		Name:                      user.Name,
		Email:                     user.Email,
		EmailVerified:             user.EmailVerified,
		Image:                     user.Image,
		PreferredLanguage:         user.PreferredLanguage,
		Role:                      string(user.Role),
		Status:                    user.Status,
		PrimaryOrganizationUnitID: user.PrimaryOrganizationUnitID,
		ManagerUserID:             user.ManagerUserID,
		JobTitle:                  user.JobTitle,
		Phone:                     user.Phone,
		LastLoginAt:               lastLoginAt,
		DeletedAt:                 deletedAt,
		CreatedAt:                 user.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:                 user.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
}

func userRoleAssignmentResponseFromModel(assignment rbac.UserRole) userRoleAssignmentResponse {
	return userRoleAssignmentResponse{
		ID:                 assignment.ID,
		UserID:             assignment.UserID,
		RoleID:             assignment.RoleID,
		ScopeType:          string(assignment.ScopeType),
		OrganizationUnitID: assignment.OrganizationUnitID,
		CreatedAt:          assignment.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:          assignment.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
}

func userGroupAssignmentResponseFromModel(assignment rbac.GroupUser) userGroupAssignmentResponse {
	return userGroupAssignmentResponse{
		UserID:    assignment.UserID,
		GroupID:   assignment.GroupID,
		CreatedAt: assignment.CreatedAt.UTC().Format(time.RFC3339Nano),
	}
}

func writeUserMutationError(c *gin.Context, err error, fallback string) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		writeError(c, http.StatusNotFound, "user not found")
		return
	}
	if isUniqueConstraintError(err) {
		writeError(c, http.StatusConflict, "user already exists")
		return
	}
	writeError(c, http.StatusInternalServerError, fallback)
}
