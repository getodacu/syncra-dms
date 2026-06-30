package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	auth "ai.ro/syncra/internal/auth"
	"ai.ro/syncra/internal/billing"
)

const minPasswordLength = 8
const maxPasswordLength = 128
const maxAuthNameCharacters = 255
const maxAuthEmailBytes = 320
const maxAuthAvatarImageBytes = 5 * 1024 * 1024
const maxPatchAuthUserBodyBytes = ((maxAuthAvatarImageBytes+2)/3)*4 + 4096

// signUpEmailRequest contains credentials for email sign-up.
//
// swagger:model signUpEmailRequest
type signUpEmailRequest struct {
	// User display name.
	// required: true
	Name string `json:"name"`
	// User email address.
	// required: true
	// format: email
	Email string `json:"email"`
	// User password.
	// required: true
	Password string `json:"password"`
	// Optional profile image URL.
	//
	// Extensions:
	// ---
	// x-nullable: true
	// ---
	Image *string `json:"image"`
}

// signInEmailRequest contains credentials for email sign-in.
//
// swagger:model signInEmailRequest
type signInEmailRequest struct {
	// User email address.
	// required: true
	// format: email
	Email string `json:"email"`
	// User password.
	// required: true
	Password   string `json:"password"`
	RememberMe *bool  `json:"rememberMe"`
}

// swagger:type string
type optionalStringField struct {
	Set   bool
	Value *string
}

func (f *optionalStringField) UnmarshalJSON(raw []byte) error {
	f.Set = true
	if string(raw) == "null" {
		f.Value = nil
		return nil
	}
	var value string
	if err := json.Unmarshal(raw, &value); err != nil {
		return err
	}
	f.Value = &value
	return nil
}

// patchAuthUserRequest contains optional authenticated user profile updates.
//
// swagger:model patchAuthUserRequest
type patchAuthUserRequest struct {
	// Optional user display name. Omit to leave unchanged; null is rejected.
	// max length: 255
	Name optionalStringField `json:"name"`
	// Optional user email address. Omit to leave unchanged; null is rejected.
	//
	// swagger:strfmt email
	Email optionalStringField `json:"email"`
	// Optional profile image data URL. Omit to leave unchanged; send null to clear the image. Data URLs must be base64-encoded PNG, JPEG, GIF, AVIF, APNG, SVG, or WebP images up to 5 MiB decoded.
	//
	// Extensions:
	// ---
	// x-nullable: true
	// ---
	Image optionalStringField `json:"image"`
	// Optional preferred UI language. Omit to leave unchanged. Must be en or ro.
	PreferredLanguage optionalStringField `json:"preferredLanguage"`
	// Optional new password. Omit to leave unchanged; null is rejected. Treat as password input.
	// min length: 8
	// max length: 128
	Password optionalStringField `json:"password"`
}

// sendVerificationOTPRequest requests an email verification OTP.
//
// swagger:model sendVerificationOTPRequest
type sendVerificationOTPRequest struct {
	// User email address.
	// required: true
	// format: email
	Email string `json:"email"`
	// Verification type. Must be email-verification when set.
	Type string `json:"type"`
}

// verifyEmailOTPRequest verifies an email OTP.
//
// swagger:model verifyEmailOTPRequest
type verifyEmailOTPRequest struct {
	// User email address.
	// required: true
	// format: email
	Email string `json:"email"`
	// One-time verification code.
	// required: true
	OTP string `json:"otp"`
}

// requestPasswordResetRequest requests a password reset link.
//
// swagger:model requestPasswordResetRequest
type requestPasswordResetRequest struct {
	// User email address.
	// required: true
	// format: email
	Email string `json:"email"`
}

// confirmPasswordResetRequest confirms a password reset token and sets a new password.
//
// swagger:model confirmPasswordResetRequest
type confirmPasswordResetRequest struct {
	// User email address.
	// required: true
	// format: email
	Email string `json:"email"`
	// One-time password reset token.
	// required: true
	Token string `json:"token"`
	// New password.
	// required: true
	// min length: 8
	// max length: 128
	Password string `json:"password"`
}

// signUpEmailResponse is returned after email sign-up.
//
// swagger:model signUpEmailResponse
type signUpEmailResponse struct {
	User                  authUserResponse `json:"user"`
	VerificationRequired  bool             `json:"verificationRequired"`
	VerificationCode      string           `json:"verificationCode,omitempty"`
	VerificationExpiresAt string           `json:"verificationExpiresAt,omitempty"`
}

// signOutResponse confirms a sign-out request.
//
// swagger:model signOutResponse
type signOutResponse struct {
	Success bool `json:"success"`
}

// sendVerificationOTPResponse reports email verification OTP status.
//
// swagger:model sendVerificationOTPResponse
type sendVerificationOTPResponse struct {
	OK                    bool   `json:"ok"`
	VerificationCode      string `json:"verificationCode,omitempty"`
	VerificationExpiresAt string `json:"verificationExpiresAt,omitempty"`
}

// verifyEmailOTPResponse reports email verification status.
//
// swagger:model verifyEmailOTPResponse
type verifyEmailOTPResponse struct {
	OK   bool             `json:"ok"`
	User authUserResponse `json:"user"`
}

// requestPasswordResetResponse reports password reset request status.
//
// swagger:model requestPasswordResetResponse
type requestPasswordResetResponse struct {
	OK             bool   `json:"ok"`
	ResetToken     string `json:"resetToken,omitempty"`
	ResetExpiresAt string `json:"resetExpiresAt,omitempty"`
}

// confirmPasswordResetResponse reports password reset completion.
//
// swagger:model confirmPasswordResetResponse
type confirmPasswordResetResponse struct {
	OK bool `json:"ok"`
}

var errInvalidVerificationCode = errors.New("invalid verification code")
var errInvalidPasswordResetToken = errors.New("invalid password reset token")
var errEmailAlreadyExists = errors.New("email is already in use")
var errInvalidAvatarImage = errors.New("invalid avatar image")

var dummyCredentialPasswordHash = mustHashDummyCredentialPassword()

func (h *Handler) SignUpEmail(c *gin.Context) {
	var req signUpEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	name := strings.TrimSpace(req.Name)
	email := normalizeEmail(req.Email)
	if err := validateSignupInput(name, email, req.Password); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	var user auth.User
	responseUser := authUserResponse{}
	verificationRequired := true
	var verificationCode string
	var verificationExpiresAt time.Time
	if err := h.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("email = ?", email).Limit(1).Find(&user)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected > 0 {
			responseUser = genericAuthUserJSON(name, email)
			if user.EmailVerified {
				return nil
			}
			return rotateEmailVerification(tx, h.BetterAuthSecret, email, h.authVerificationTTL(), &verificationCode, &verificationExpiresAt)
		}

		hash, err := auth.HashPassword(req.Password)
		if err != nil {
			return err
		}
		user = auth.User{Name: name, Email: email, Image: req.Image, EmailVerified: false}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		account := auth.AuthAccount{
			AccountID:  user.ID,
			ProviderID: auth.CredentialProviderID,
			UserID:     user.ID,
			Password:   hash,
		}
		if err := tx.Create(&account).Error; err != nil {
			return err
		}
		responseUser = authUserJSON(user)
		return rotateEmailVerification(tx, h.BetterAuthSecret, email, h.authVerificationTTL(), &verificationCode, &verificationExpiresAt)
	}); err != nil {
		writeError(c, http.StatusInternalServerError, "failed to sign up")
		return
	}

	out := signUpEmailResponse{
		User:                 responseUser,
		VerificationRequired: verificationRequired,
	}
	if h.trustedAuthDeliveryRequest(c) && verificationCode != "" {
		out.VerificationCode = verificationCode
		out.VerificationExpiresAt = verificationExpiresAt.UTC().Format(time.RFC3339Nano)
	}
	c.JSON(http.StatusOK, out)
}

func (h *Handler) SignInEmail(c *gin.Context) {
	var req signInEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	email := normalizeEmail(req.Email)
	if !validEmail(email) {
		writeError(c, http.StatusBadRequest, "valid email is required")
		return
	}
	if req.Password == "" {
		writeError(c, http.StatusBadRequest, "password is required")
		return
	}

	var user auth.User
	result := h.DB.Where("email = ?", email).Limit(1).Find(&user)
	if result.Error != nil {
		writeError(c, http.StatusInternalServerError, "failed to load user")
		return
	}
	if result.RowsAffected == 0 {
		consumeDummyCredentialPassword(req.Password)
		writeError(c, http.StatusUnauthorized, "invalid email or password")
		return
	}

	var account auth.AuthAccount
	result = h.DB.Where("user_id = ? AND provider_id = ?", user.ID, auth.CredentialProviderID).Limit(1).Find(&account)
	if result.Error != nil {
		writeError(c, http.StatusInternalServerError, "failed to load auth account")
		return
	}
	if result.RowsAffected == 0 {
		consumeDummyCredentialPassword(req.Password)
		writeError(c, http.StatusUnauthorized, "invalid email or password")
		return
	}
	if !auth.VerifyPassword(req.Password, account.Password) {
		writeError(c, http.StatusUnauthorized, "invalid email or password")
		return
	}
	if !user.EmailVerified {
		writeError(c, http.StatusForbidden, "email is not verified")
		return
	}

	token, err := auth.GenerateSessionToken()
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to create session")
		return
	}
	now := time.Now().UTC()
	session := auth.Session{
		Token:     token,
		UserID:    user.ID,
		ExpiresAt: now.Add(h.authSessionTTL()),
		IPAddress: c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&auth.User{}).Where("id = ?", user.ID).Update("last_login_at", now).Error; err != nil {
			return err
		}
		user.LastLoginAt = &now
		return tx.Create(&session).Error
	}); err != nil {
		writeError(c, http.StatusInternalServerError, "failed to create session")
		return
	}

	remember := true
	if req.RememberMe != nil {
		remember = *req.RememberMe
	}
	h.setSessionCookie(c, session.Token, session.ExpiresAt, remember)
	c.JSON(http.StatusOK, signInEmailResponse{
		Session: authSessionJSON(session),
		User:    authUserJSON(user),
	})
}

func (h *Handler) GetSession(c *gin.Context) {
	session, ok, err := h.loadAuthenticatedSession(c)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}
	if !ok {
		c.JSON(http.StatusOK, nil)
		return
	}

	c.JSON(http.StatusOK, authSessionPayloadJSON(session))
}

func (h *Handler) PatchAuthUser(c *gin.Context) {
	session, ok, err := h.loadAuthenticatedSession(c)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}
	if !ok {
		writeError(c, http.StatusUnauthorized, "authentication required")
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxPatchAuthUserBodyBytes)
	var req patchAuthUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			writeError(c, http.StatusBadRequest, "request body too large")
			return
		}
		writeError(c, http.StatusBadRequest, "invalid JSON body")
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
		effectiveUser := effectiveAuthUser(session)
		if email != effectiveUser.Email {
			updates["email_verified"] = false
		}
	}
	if req.Image.Set {
		if req.Image.Value == nil {
			updates["image"] = nil
		} else {
			image := *req.Image.Value
			if err := validateAvatarImage(image); err != nil {
				writeError(c, http.StatusBadRequest, err.Error())
				return
			}
			updates["image"] = image
		}
	}
	if req.PreferredLanguage.Set {
		if req.PreferredLanguage.Value == nil || !validPreferredLanguage(*req.PreferredLanguage.Value) {
			writeError(c, http.StatusBadRequest, "preferredLanguage must be en or ro")
			return
		}
		updates["preferred_language"] = *req.PreferredLanguage.Value
	}

	var passwordHash string
	if req.Password.Set {
		passwordCharacters := 0
		if req.Password.Value != nil {
			passwordCharacters = utf8.RuneCountInString(*req.Password.Value)
		}
		if req.Password.Value == nil || passwordCharacters < minPasswordLength || passwordCharacters > maxPasswordLength {
			writeError(c, http.StatusBadRequest, "password must be between 8 and 128 characters")
			return
		}
		passwordHash, err = auth.HashPassword(*req.Password.Value)
		if err != nil {
			writeError(c, http.StatusInternalServerError, "failed to hash password")
			return
		}
	}

	var user auth.User
	if err := h.DB.Transaction(func(tx *gorm.DB) error {
		if req.Email.Set && req.Email.Value != nil {
			email := updates["email"].(string)
			effectiveUser := effectiveAuthUser(session)
			var count int64
			if err := tx.Model(&auth.User{}).Where("email = ? AND id <> ?", email, effectiveUser.ID).Count(&count).Error; err != nil {
				return err
			}
			if count > 0 {
				return errEmailAlreadyExists
			}
		}

		effectiveUser := effectiveAuthUser(session)
		if len(updates) > 0 {
			if err := tx.Model(&auth.User{}).Where("id = ?", effectiveUser.ID).Updates(updates).Error; err != nil {
				return err
			}
		}

		if passwordHash != "" {
			now := time.Now().UTC()
			account := auth.AuthAccount{
				AccountID:  effectiveUser.ID,
				ProviderID: auth.CredentialProviderID,
				UserID:     effectiveUser.ID,
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
		}

		return tx.First(&user, "id = ?", effectiveUser.ID).Error
	}); err != nil {
		if errors.Is(err, errEmailAlreadyExists) || (req.Email.Set && isUniqueConstraintError(err)) {
			writeError(c, http.StatusConflict, errEmailAlreadyExists.Error())
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to update user")
		return
	}

	c.JSON(http.StatusOK, authUserJSON(user))
}

func (h *Handler) SignOut(c *gin.Context) {
	token := sessionTokenFromRequest(c)
	if token != "" {
		if err := h.DB.Where("token = ?", token).Delete(&auth.Session{}).Error; err != nil {
			writeError(c, http.StatusInternalServerError, "failed to delete session")
			return
		}
	}
	h.clearSessionCookie(c)
	c.JSON(http.StatusOK, signOutResponse{Success: true})
}

func mustHashDummyCredentialPassword() string {
	hash, err := auth.HashPassword("syncra-dummy-credential-password")
	if err != nil {
		panic(err)
	}
	return hash
}

func consumeDummyCredentialPassword(password string) {
	_ = auth.VerifyPassword(password, dummyCredentialPasswordHash)
}

func (h *Handler) SendVerificationOTP(c *gin.Context) {
	var req sendVerificationOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.Type != "" && req.Type != "email-verification" {
		writeError(c, http.StatusBadRequest, "type must be email-verification")
		return
	}
	email := normalizeEmail(req.Email)
	if !validEmail(email) {
		writeError(c, http.StatusBadRequest, "valid email is required")
		return
	}

	var user auth.User
	result := h.DB.Where("email = ?", email).Limit(1).Find(&user)
	if result.Error != nil {
		writeError(c, http.StatusInternalServerError, "failed to load user")
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, sendVerificationOTPResponse{OK: true})
		return
	}
	if user.EmailVerified {
		c.JSON(http.StatusOK, sendVerificationOTPResponse{OK: true})
		return
	}

	var code string
	var expiresAt time.Time
	if err := h.DB.Transaction(func(tx *gorm.DB) error {
		return rotateEmailVerification(tx, h.BetterAuthSecret, email, h.authVerificationTTL(), &code, &expiresAt)
	}); err != nil {
		writeError(c, http.StatusInternalServerError, "failed to create verification code")
		return
	}

	out := sendVerificationOTPResponse{OK: true}
	if h.trustedAuthDeliveryRequest(c) {
		out.VerificationCode = code
		out.VerificationExpiresAt = expiresAt.UTC().Format(time.RFC3339Nano)
	}
	c.JSON(http.StatusOK, out)
}

func (h *Handler) VerifyEmailOTP(c *gin.Context) {
	var req verifyEmailOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	email := normalizeEmail(req.Email)
	otp := strings.TrimSpace(req.OTP)
	if !validEmail(email) {
		writeError(c, http.StatusBadRequest, "valid email is required")
		return
	}
	if otp == "" {
		writeError(c, http.StatusBadRequest, "otp is required")
		return
	}

	var user auth.User
	invalidVerificationCode := false
	if err := h.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("email = ?", email).Limit(1).Find(&user)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			invalidVerificationCode = true
			return nil
		}

		identifier := verificationIdentifier(email)
		var verification auth.Verification
		result = tx.Where("identifier = ?", identifier).Limit(1).Find(&verification)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			invalidVerificationCode = true
			return nil
		}
		if !verification.ExpiresAt.After(time.Now().UTC()) {
			if _, err := consumeEmailVerification(tx, verification, identifier); err != nil {
				return err
			}
			invalidVerificationCode = true
			return nil
		}
		if !auth.VerifyCode(h.BetterAuthSecret, identifier, otp, verification.Value) {
			invalidVerificationCode = true
			return nil
		}

		consumed, err := consumeEmailVerification(tx, verification, identifier)
		if err != nil {
			return err
		}
		if !consumed {
			invalidVerificationCode = true
			return nil
		}
		if err := tx.Model(&auth.User{}).Where("id = ?", user.ID).Update("email_verified", true).Error; err != nil {
			return err
		}
		if _, err := billing.GrantSignupBonus(c.Request.Context(), tx, billing.GrantSignupBonusInput{
			UserID:  user.ID,
			Credits: h.onboardingCredits(),
			Now:     time.Now().UTC(),
		}); err != nil {
			return err
		}
		return tx.First(&user, "id = ?", user.ID).Error
	}); err != nil {
		writeError(c, http.StatusInternalServerError, "failed to verify email")
		return
	}
	if invalidVerificationCode {
		writeError(c, http.StatusBadRequest, errInvalidVerificationCode.Error())
		return
	}

	c.JSON(http.StatusOK, verifyEmailOTPResponse{
		OK:   true,
		User: authUserJSON(user),
	})
}

func (h *Handler) RequestPasswordReset(c *gin.Context) {
	var req requestPasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	email := normalizeEmail(req.Email)
	if !validEmail(email) {
		writeError(c, http.StatusBadRequest, "valid email is required")
		return
	}

	var user auth.User
	result := h.DB.Where("email = ?", email).Limit(1).Find(&user)
	if result.Error != nil {
		writeError(c, http.StatusInternalServerError, "failed to load user")
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, requestPasswordResetResponse{OK: true})
		return
	}

	var token string
	var expiresAt time.Time
	if err := h.DB.Transaction(func(tx *gorm.DB) error {
		return rotatePasswordResetVerification(tx, h.BetterAuthSecret, email, h.authVerificationTTL(), &token, &expiresAt)
	}); err != nil {
		writeError(c, http.StatusInternalServerError, "failed to create password reset token")
		return
	}

	out := requestPasswordResetResponse{OK: true}
	if h.trustedAuthDeliveryRequest(c) {
		out.ResetToken = token
		out.ResetExpiresAt = expiresAt.UTC().Format(time.RFC3339Nano)
	}
	c.JSON(http.StatusOK, out)
}

func (h *Handler) ConfirmPasswordReset(c *gin.Context) {
	var req confirmPasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	email := normalizeEmail(req.Email)
	token := strings.TrimSpace(req.Token)
	if !validEmail(email) {
		writeError(c, http.StatusBadRequest, "valid email is required")
		return
	}
	if token == "" {
		writeError(c, http.StatusBadRequest, "password reset token is required")
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

	invalidResetToken := false
	if err := h.DB.Transaction(func(tx *gorm.DB) error {
		var user auth.User
		result := tx.Where("email = ?", email).Limit(1).Find(&user)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			invalidResetToken = true
			return nil
		}

		identifier := passwordResetIdentifier(email)
		var verification auth.Verification
		result = tx.Where("identifier = ?", identifier).Limit(1).Find(&verification)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			invalidResetToken = true
			return nil
		}
		if !verification.ExpiresAt.After(time.Now().UTC()) {
			if _, err := consumePasswordResetVerification(tx, verification, identifier); err != nil {
				return err
			}
			invalidResetToken = true
			return nil
		}
		if !auth.VerifyCode(h.BetterAuthSecret, identifier, token, verification.Value) {
			invalidResetToken = true
			return nil
		}

		consumed, err := consumePasswordResetVerification(tx, verification, identifier)
		if err != nil {
			return err
		}
		if !consumed {
			invalidResetToken = true
			return nil
		}

		now := time.Now().UTC()
		account := auth.AuthAccount{
			AccountID:  user.ID,
			ProviderID: auth.CredentialProviderID,
			UserID:     user.ID,
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
		if err := tx.Model(&auth.User{}).Where("id = ?", user.ID).Update("email_verified", true).Error; err != nil {
			return err
		}
		return tx.Where("user_id = ?", user.ID).Delete(&auth.Session{}).Error
	}); err != nil {
		writeError(c, http.StatusInternalServerError, "failed to reset password")
		return
	}
	if invalidResetToken {
		writeError(c, http.StatusBadRequest, errInvalidPasswordResetToken.Error())
		return
	}

	c.JSON(http.StatusOK, confirmPasswordResetResponse{OK: true})
}

func consumeEmailVerification(tx *gorm.DB, verification auth.Verification, identifier string) (bool, error) {
	result := tx.Where(
		"id = ? AND identifier = ? AND value = ? AND expires_at = ?",
		verification.ID,
		identifier,
		verification.Value,
		verification.ExpiresAt,
	).Delete(&auth.Verification{})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected == 1, nil
}

func consumePasswordResetVerification(tx *gorm.DB, verification auth.Verification, identifier string) (bool, error) {
	result := tx.Where(
		"id = ? AND identifier = ? AND value = ? AND expires_at = ?",
		verification.ID,
		identifier,
		verification.Value,
		verification.ExpiresAt,
	).Delete(&auth.Verification{})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected == 1, nil
}

func validateSignupInput(name string, email string, password string) error {
	if name == "" {
		return errors.New("name is required")
	}
	if utf8.RuneCountInString(name) > maxAuthNameCharacters {
		return errors.New("name must be at most 255 characters")
	}
	if !validEmail(email) {
		return errors.New("valid email is required")
	}
	if len(password) < minPasswordLength || len(password) > maxPasswordLength {
		return errors.New("password must be between 8 and 128 characters")
	}
	return nil
}

func (h *Handler) loadAuthenticatedSession(c *gin.Context) (auth.Session, bool, error) {
	token, tokenPresent := sessionTokenFromRequestWithPresence(c)
	if token == "" {
		if tokenPresent {
			h.clearSessionCookie(c)
		}
		return auth.Session{}, false, nil
	}

	var session auth.Session
	result := h.DB.Preload("User").Preload("ImpersonatedUser").Where("token = ?", token).Limit(1).Find(&session)
	if result.Error != nil {
		return auth.Session{}, false, errors.New("failed to load session")
	}
	if result.RowsAffected == 0 {
		h.clearSessionCookie(c)
		return auth.Session{}, false, nil
	}
	if !session.ExpiresAt.After(time.Now().UTC()) {
		if err := h.DB.Delete(&session).Error; err != nil {
			return auth.Session{}, false, errors.New("failed to delete expired session")
		}
		h.clearSessionCookie(c)
		return auth.Session{}, false, nil
	}
	if session.User.ID == "" || session.User.Email == "" {
		if err := h.DB.Where("token = ?", token).Delete(&auth.Session{}).Error; err != nil {
			return auth.Session{}, false, errors.New("failed to delete stale session")
		}
		h.clearSessionCookie(c)
		return auth.Session{}, false, nil
	}
	if session.ImpersonatedUserID == nil && session.ImpersonationStartedAt != nil {
		if err := clearSessionImpersonation(c.Request.Context(), h.DB, session.ID); err != nil {
			return auth.Session{}, false, errors.New("failed to clear stale impersonation")
		}
		session.ImpersonationStartedAt = nil
	}
	if session.ImpersonatedUserID != nil && (session.ImpersonatedUser == nil || session.ImpersonatedUser.ID == "" || session.ImpersonatedUser.Role != auth.UserRoleUser) {
		if err := clearSessionImpersonation(c.Request.Context(), h.DB, session.ID); err != nil {
			return auth.Session{}, false, errors.New("failed to clear stale impersonation")
		}
		session.ImpersonatedUserID = nil
		session.ImpersonatedUser = nil
		session.ImpersonationStartedAt = nil
	}

	return session, true, nil
}

func clearSessionImpersonation(ctx context.Context, db *gorm.DB, sessionID string) error {
	return db.WithContext(ctx).
		Model(&auth.Session{}).
		Where("id = ?", sessionID).
		Updates(map[string]any{
			"impersonated_user_id":     nil,
			"impersonation_started_at": nil,
		}).Error
}

func validateAvatarImage(image string) error {
	if image == "" {
		return nil
	}
	header, payload, ok := strings.Cut(image, ",")
	if !ok || !strings.HasPrefix(header, "data:") {
		return errInvalidAvatarImage
	}

	metadata := strings.TrimPrefix(header, "data:")
	parts := strings.Split(metadata, ";")
	if len(parts) < 2 {
		return errInvalidAvatarImage
	}
	mimeType := strings.ToLower(parts[0])
	if _, ok := supportedAvatarImageMIMETypes[mimeType]; !ok {
		return errInvalidAvatarImage
	}
	hasBase64 := false
	for _, part := range parts[1:] {
		if strings.EqualFold(part, "base64") {
			hasBase64 = true
			break
		}
	}
	if !hasBase64 {
		return errInvalidAvatarImage
	}
	if base64.StdEncoding.DecodedLen(len(payload)) > maxAuthAvatarImageBytes+2 {
		return errInvalidAvatarImage
	}
	raw, err := base64.StdEncoding.DecodeString(payload)
	if err != nil || len(raw) > maxAuthAvatarImageBytes {
		return errInvalidAvatarImage
	}
	return nil
}

func isUniqueConstraintError(err error) bool {
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

var supportedAvatarImageMIMETypes = map[string]struct{}{
	"image/png":     {},
	"image/jpeg":    {},
	"image/jpg":     {},
	"image/gif":     {},
	"image/avif":    {},
	"image/apng":    {},
	"image/svg+xml": {},
	"image/webp":    {},
}

func genericAuthUserJSON(name string, email string) authUserResponse {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	return authUserResponse{
		Name:              name,
		Email:             email,
		EmailVerified:     false,
		PreferredLanguage: "en",
		Role:              string(auth.UserRoleUser),
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

func validPreferredLanguage(value string) bool {
	return value == "en" || value == "ro"
}

func validEmail(email string) bool {
	if len(email) > maxAuthEmailBytes {
		return false
	}
	parsed, err := mail.ParseAddress(email)
	return err == nil && parsed.Address == email
}

func rotateEmailVerification(tx *gorm.DB, secret string, email string, ttl time.Duration, codeOut *string, expiresOut *time.Time) error {
	code, err := auth.GenerateNumericCode(6)
	if err != nil {
		return err
	}
	identifier := verificationIdentifier(email)
	now := time.Now().UTC()
	expiresAt := now.Add(ttl)
	value := auth.HashCode(secret, identifier, code)
	if err := tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "identifier"}},
		DoUpdates: clause.Assignments(map[string]any{
			"value":      value,
			"expires_at": expiresAt,
			"updated_at": now,
		}),
	}).Create(&auth.Verification{
		Identifier: identifier,
		Value:      value,
		ExpiresAt:  expiresAt,
		CreatedAt:  now,
		UpdatedAt:  now,
	}).Error; err != nil {
		return err
	}
	*codeOut = code
	*expiresOut = expiresAt

	return nil
}

func rotatePasswordResetVerification(tx *gorm.DB, secret string, email string, ttl time.Duration, tokenOut *string, expiresOut *time.Time) error {
	token, err := auth.GenerateSessionToken()
	if err != nil {
		return err
	}
	identifier := passwordResetIdentifier(email)
	now := time.Now().UTC()
	expiresAt := now.Add(ttl)
	value := auth.HashCode(secret, identifier, token)
	if err := tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "identifier"}},
		DoUpdates: clause.Assignments(map[string]any{
			"value":      value,
			"expires_at": expiresAt,
			"updated_at": now,
		}),
	}).Create(&auth.Verification{
		Identifier: identifier,
		Value:      value,
		ExpiresAt:  expiresAt,
		CreatedAt:  now,
		UpdatedAt:  now,
	}).Error; err != nil {
		return err
	}
	*tokenOut = token
	*expiresOut = expiresAt

	return nil
}
