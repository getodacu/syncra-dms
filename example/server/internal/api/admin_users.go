package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"ai.ro/syncra/internal/auth"
	"ai.ro/syncra/internal/billing"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const maxAdminRequestBytes int64 = 1 << 20
const maxAdminUserPageSize = 100

type adminContextKey string

const adminSessionKey adminContextKey = "admin_session"

var (
	errLoadAdminBillingProfile = errors.New("load admin billing profile")
	errLoadAdminCreditBalance  = errors.New("load admin credit balance")
)

type adminUserListCursor struct {
	Offset    int    `json:"offset"`
	Search    string `json:"search"`
	Sort      string `json:"sort"`
	Direction string `json:"direction"`
}

type adminUserResponse struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	Email         string     `json:"email"`
	EmailVerified bool       `json:"email_verified"`
	Role          string     `json:"role"`
	Image         *string    `json:"image,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	LastLoginAt   *time.Time `json:"last_login_at"`
}

type adminUserDetailResponse struct {
	adminUserResponse
	AvailableCredits int                     `json:"available_credits"`
	BillingProfile   *BillingProfileResponse `json:"billing_profile"`
}

type adminUserListResponse struct {
	Users      []adminUserResponse `json:"users"`
	NextCursor *string             `json:"next_cursor"`
}

type patchAdminUserRequest struct {
	Name  optionalStringField `json:"name"`
	Email optionalStringField `json:"email"`
}

type adminSetPasswordRequest struct {
	Password string `json:"password"`
}

type adminSetPasswordResponse struct {
	OK bool `json:"ok"`
}

type adminAdjustUserBalanceRequest struct {
	CreditsDelta int `json:"credits_delta"`
}

type adminUpsertBillingProfileRequest struct {
	EntityType         string  `json:"entity_type"`
	BillingName        string  `json:"billing_name"`
	BillingEmail       string  `json:"billing_email"`
	CountryCode        string  `json:"country_code"`
	AddressLine1       string  `json:"address_line1"`
	AddressLine2       *string `json:"address_line2"`
	City               string  `json:"city"`
	Region             *string `json:"region"`
	PostalCode         string  `json:"postal_code"`
	FiscalCode         *string `json:"fiscal_code"`
	RegistrationNumber *string `json:"registration_number"`
}

func (h *Handler) requireAdminSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, ok, err := h.loadAuthenticatedSession(c)
		if err != nil {
			writeError(c, http.StatusInternalServerError, err.Error())
			c.Abort()
			return
		}
		if !ok {
			writeError(c, http.StatusUnauthorized, "authentication required")
			c.Abort()
			return
		}
		if session.User.Role != auth.UserRoleAdmin {
			writeError(c, http.StatusForbidden, "admin access required")
			c.Abort()
			return
		}
		c.Set(string(adminSessionKey), session)
		c.Next()
	}
}

func (h *Handler) ListAdminUsers(c *gin.Context) {
	search := strings.TrimSpace(c.Query("search"))
	sortField, err := parseAdminUserSort(c.Query("sort"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	direction, err := parseAdminUserDirection(c.Query("direction"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	size, err := parseAdminUserPageSize(c.Query("size"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	offset, err := parseAdminUserCursor(c.Query("cursor"), search, sortField, direction)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	query := h.DB.WithContext(c.Request.Context()).Model(&auth.User{})
	if search != "" {
		term := "%" + strings.ToLower(search) + "%"
		query = query.Where("lower(name) LIKE ? OR lower(email) LIKE ?", term, term)
	}

	var users []auth.User
	if err := query.
		Order(adminUserOrder(sortField, direction)).
		Limit(size + 1).
		Offset(offset).
		Find(&users).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to list users")
		return
	}

	var nextCursor *string
	if len(users) > size {
		users = users[:size]
		encoded, err := encodeAdminUserCursor(adminUserListCursor{
			Offset:    offset + size,
			Search:    search,
			Sort:      sortField,
			Direction: direction,
		})
		if err != nil {
			writeError(c, http.StatusInternalServerError, "failed to encode cursor")
			return
		}
		nextCursor = &encoded
	}

	out := make([]adminUserResponse, 0, len(users))
	for _, user := range users {
		out = append(out, adminUserJSON(user))
	}
	c.JSON(http.StatusOK, adminUserListResponse{Users: out, NextCursor: nextCursor})
}

func (h *Handler) GetAdminUser(c *gin.Context) {
	userID, ok := parseAdminUserIDParam(c)
	if !ok {
		return
	}

	response, err := h.loadAdminUserDetail(c, userID)
	if err != nil {
		h.writeAdminUserDetailError(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h *Handler) PatchAdminUser(c *gin.Context) {
	userID, ok := parseAdminUserIDParam(c)
	if !ok {
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxAdminRequestBytes)
	var req patchAdminUserRequest
	if !decodeStrictAdminJSON(c, &req, "invalid user update payload") {
		return
	}

	updates := map[string]any{}
	if req.Name.Set {
		if req.Name.Value == nil {
			writeError(c, http.StatusBadRequest, "name is required")
			return
		}
		name := strings.TrimSpace(*req.Name.Value)
		if name == "" {
			writeError(c, http.StatusBadRequest, "name is required")
			return
		}
		if utf8.RuneCountInString(name) > maxAuthNameCharacters {
			writeError(c, http.StatusBadRequest, "name must be at most 255 characters")
			return
		}
		updates["name"] = name
	}
	if req.Email.Set {
		if req.Email.Value == nil {
			writeError(c, http.StatusBadRequest, "valid email is required")
			return
		}
		email := normalizeEmail(*req.Email.Value)
		if !validEmail(email) {
			writeError(c, http.StatusBadRequest, "valid email is required")
			return
		}
		updates["email"] = email
	}

	var user auth.User
	if err := h.DB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		result := tx.First(&user, "id = ?", userID)
		if result.Error != nil {
			return result.Error
		}
		if req.Email.Set && req.Email.Value != nil {
			email := updates["email"].(string)
			var count int64
			if err := tx.Model(&auth.User{}).Where("email = ? AND id <> ?", email, userID).Count(&count).Error; err != nil {
				return err
			}
			if count > 0 {
				return errEmailAlreadyExists
			}
		}
		if len(updates) > 0 {
			if err := tx.Model(&auth.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
				return err
			}
		}
		return tx.First(&user, "id = ?", userID).Error
	}); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusNotFound, "user not found")
			return
		}
		if errors.Is(err, errEmailAlreadyExists) || isUniqueConstraintError(err) {
			writeError(c, http.StatusConflict, errEmailAlreadyExists.Error())
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to update user")
		return
	}

	c.JSON(http.StatusOK, adminUserJSON(user))
}

func (h *Handler) SetAdminUserPassword(c *gin.Context) {
	userID, ok := parseAdminUserIDParam(c)
	if !ok {
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxAdminRequestBytes)
	var req adminSetPasswordRequest
	if !decodeStrictAdminJSON(c, &req, "invalid password reset payload") {
		return
	}
	passwordCharacters := utf8.RuneCountInString(req.Password)
	if passwordCharacters < minPasswordLength || passwordCharacters > maxPasswordLength {
		writeError(c, http.StatusBadRequest, "password must be between 8 and 128 characters")
		return
	}
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to hash password")
		return
	}

	if err := h.DB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&auth.User{}).Where("id = ?", userID).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return gorm.ErrRecordNotFound
		}

		now := time.Now().UTC()
		account := auth.AuthAccount{
			AccountID:  userID,
			ProviderID: auth.CredentialProviderID,
			UserID:     userID,
			Password:   passwordHash,
			CreatedAt:  now,
			UpdatedAt:  now,
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "provider_id"},
				{Name: "account_id"},
			},
			DoUpdates: clause.Assignments(map[string]any{
				"password":   passwordHash,
				"updated_at": now,
			}),
		}).Create(&account).Error; err != nil {
			return err
		}
		return tx.Where("user_id = ?", userID).Delete(&auth.Session{}).Error
	}); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusNotFound, "user not found")
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to reset password")
		return
	}

	c.JSON(http.StatusOK, adminSetPasswordResponse{OK: true})
}

func (h *Handler) AdjustAdminUserBalance(c *gin.Context) {
	userID, ok := parseAdminUserIDParam(c)
	if !ok {
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxAdminRequestBytes)
	var req adminAdjustUserBalanceRequest
	if !decodeStrictAdminJSON(c, &req, "invalid balance adjustment payload") {
		return
	}
	if req.CreditsDelta == 0 {
		writeError(c, http.StatusBadRequest, "credits_delta must be non-zero")
		return
	}
	if err := validateBillingUserExists(h.DB.WithContext(c.Request.Context()), userID); err != nil {
		if errors.Is(err, errInvalidUserID) {
			writeError(c, http.StatusNotFound, "user not found")
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to validate user")
		return
	}

	session, _ := adminSessionFromContext(c)
	adminID := session.UserID
	if adminID == "" {
		adminID = "unknown"
	}
	_, err := billing.AdjustCredits(c.Request.Context(), h.DB, billing.AdjustCreditsInput{
		UserID:         userID,
		Delta:          req.CreditsDelta,
		Now:            time.Now().UTC(),
		IdempotencyKey: fmt.Sprintf("admin_adjustment:%s:%s:%s", userID, adminID, uuid.NewString()),
	})
	if err != nil {
		if errors.Is(err, billing.ErrInsufficientCredits) {
			writeError(c, http.StatusBadRequest, "insufficient credits")
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to adjust balance")
		return
	}

	response, err := h.loadAdminUserDetail(c, userID)
	if err != nil {
		h.writeAdminUserDetailError(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h *Handler) UpsertAdminUserBillingProfile(c *gin.Context) {
	userID, ok := parseAdminUserIDParam(c)
	if !ok {
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxAdminRequestBytes)
	var req adminUpsertBillingProfileRequest
	if !decodeStrictAdminJSON(c, &req, "invalid billing profile request") {
		return
	}
	if err := validateBillingUserExists(h.DB.WithContext(c.Request.Context()), userID); err != nil {
		if errors.Is(err, errInvalidUserID) {
			writeError(c, http.StatusNotFound, "user not found")
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to validate user")
		return
	}

	profile, err := billing.UpsertBillingProfile(c.Request.Context(), h.DB, billing.UpsertBillingProfileInput{
		UserID:             userID,
		EntityType:         billing.BillingEntityType(req.EntityType),
		BillingName:        req.BillingName,
		BillingEmail:       req.BillingEmail,
		CountryCode:        req.CountryCode,
		AddressLine1:       req.AddressLine1,
		AddressLine2:       req.AddressLine2,
		City:               req.City,
		Region:             req.Region,
		PostalCode:         req.PostalCode,
		FiscalCode:         req.FiscalCode,
		RegistrationNumber: req.RegistrationNumber,
	})
	if err != nil {
		if message, ok := billingProfileValidationError(err); ok {
			writeError(c, http.StatusBadRequest, message)
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to upsert billing profile")
		return
	}
	c.JSON(http.StatusOK, billingProfileResponse(profile))
}

func (h *Handler) loadAdminUserDetail(c *gin.Context, userID string) (adminUserDetailResponse, error) {
	var user auth.User
	if err := h.DB.WithContext(c.Request.Context()).First(&user, "id = ?", userID).Error; err != nil {
		return adminUserDetailResponse{}, err
	}

	response := adminUserDetailResponse{adminUserResponse: adminUserJSON(user)}
	var profile billing.BillingProfile
	if err := h.DB.WithContext(c.Request.Context()).First(&profile, "user_id = ?", userID).Error; err == nil {
		out := billingProfileResponse(profile)
		response.BillingProfile = &out
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return adminUserDetailResponse{}, fmt.Errorf("%w: %v", errLoadAdminBillingProfile, err)
	}

	balance, err := billing.AvailableCredits(c.Request.Context(), h.DB, userID, time.Now().UTC())
	if err != nil {
		return adminUserDetailResponse{}, fmt.Errorf("%w: %v", errLoadAdminCreditBalance, err)
	}
	response.AvailableCredits = balance.Available
	return response, nil
}

func (h *Handler) writeAdminUserDetailError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		writeError(c, http.StatusNotFound, "user not found")
	case errors.Is(err, errLoadAdminBillingProfile):
		writeError(c, http.StatusInternalServerError, "failed to load billing profile")
	case errors.Is(err, errLoadAdminCreditBalance):
		writeError(c, http.StatusInternalServerError, "failed to load credit balance")
	default:
		writeError(c, http.StatusInternalServerError, "failed to load user")
	}
}

func adminSessionFromContext(c *gin.Context) (auth.Session, bool) {
	raw, ok := c.Get(string(adminSessionKey))
	if !ok {
		return auth.Session{}, false
	}
	session, ok := raw.(auth.Session)
	return session, ok
}

func adminUserJSON(user auth.User) adminUserResponse {
	role := string(user.Role)
	if role == "" {
		role = string(auth.UserRoleUser)
	}
	return adminUserResponse{
		ID:            user.ID,
		Name:          user.Name,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		Role:          role,
		Image:         user.Image,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		LastLoginAt:   user.LastLoginAt,
	}
}

func parseAdminUserIDParam(c *gin.Context) (string, bool) {
	id, err := uuid.Parse(strings.TrimSpace(c.Param("id")))
	if err != nil || id == uuid.Nil {
		writeError(c, http.StatusBadRequest, "invalid user id")
		return "", false
	}
	return id.String(), true
}

func parseAdminUserSort(raw string) (string, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return "created_at", nil
	}
	switch value {
	case "created_at", "last_login_at":
		return value, nil
	default:
		return "", errors.New("sort must be created_at or last_login_at")
	}
}

func parseAdminUserDirection(raw string) (string, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return "desc", nil
	}
	switch value {
	case "asc", "desc":
		return value, nil
	default:
		return "", errors.New("direction must be asc or desc")
	}
}

func parseAdminUserPageSize(raw string) (int, error) {
	if strings.TrimSpace(raw) == "" {
		return 20, nil
	}
	var size int
	if err := json.Unmarshal([]byte(strings.TrimSpace(raw)), &size); err != nil || size < 1 || size > maxAdminUserPageSize {
		return 0, errors.New("size must be between 1 and 100")
	}
	return size, nil
}

func adminUserOrder(sortField string, direction string) string {
	if sortField == "last_login_at" {
		return "last_login_at IS NULL asc, last_login_at " + direction + ", id " + direction
	}
	return "created_at " + direction + ", id " + direction
}

func parseAdminUserCursor(raw string, search string, sortField string, direction string) (int, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, nil
	}
	decoded, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil {
		return 0, errors.New("invalid cursor")
	}
	var cursor adminUserListCursor
	if err := json.Unmarshal(decoded, &cursor); err != nil {
		return 0, errors.New("invalid cursor")
	}
	if cursor.Offset < 0 || cursor.Search != search || cursor.Sort != sortField || cursor.Direction != direction {
		return 0, errors.New("invalid cursor")
	}
	return cursor.Offset, nil
}

func encodeAdminUserCursor(cursor adminUserListCursor) (string, error) {
	raw, err := json.Marshal(cursor)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}

func decodeStrictAdminJSON(c *gin.Context, out any, fallback string) bool {
	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(out); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			writeError(c, http.StatusBadRequest, "request body too large")
			return false
		}
		writeError(c, http.StatusBadRequest, fallback)
		return false
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		writeError(c, http.StatusBadRequest, fallback)
		return false
	}
	return true
}
