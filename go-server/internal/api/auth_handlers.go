package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"

	auth "ai.ro/syncra/dms/internal/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	authSessionCookieName = "auth.session_token"
	authDeliveryHeader    = "X-Syncra-Auth-Delivery-Token"
	internalAPIHeader     = "X-Syncra-Internal-Token"
	minPasswordLength     = 8
	maxPasswordLength     = 128
	maxAuthNameCharacters = 255
	maxAuthEmailBytes     = 320
	oauthStateTTL         = 10 * time.Minute
)

type authHandler struct {
	db                  *gorm.DB
	betterAuthSecret    string
	authDeliveryToken   string
	internalAPIToken    string
	authSessionTTL      time.Duration
	authVerificationTTL time.Duration
	authCookieSecure    bool
	googleClientID      string
	googleClientSecret  string
	gitHubClientID      string
	gitHubClientSecret  string
	oauthProfileFetcher func(context.Context, string, string, string) (OAuthProfile, error)
	httpClient          *http.Client
}

type signUpEmailRequest struct {
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Image    *string `json:"image"`
}

type signInEmailRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	RememberMe *bool  `json:"rememberMe"`
}

type sendVerificationOTPRequest struct {
	Email string `json:"email"`
	Type  string `json:"type"`
}

type verifyEmailOTPRequest struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

type requestPasswordResetRequest struct {
	Email string `json:"email"`
}

type confirmPasswordResetRequest struct {
	Email    string `json:"email"`
	Token    string `json:"token"`
	Password string `json:"password"`
}

type oauthStartRequest struct {
	RedirectURI string `json:"redirectURI"`
}

type oauthCallbackRequest struct {
	Code        string `json:"code"`
	State       string `json:"state"`
	RedirectURI string `json:"redirectURI"`
}

type authUserResponse struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	Email             string  `json:"email"`
	EmailVerified     bool    `json:"emailVerified"`
	Image             *string `json:"image"`
	PreferredLanguage string  `json:"preferredLanguage"`
	Role              string  `json:"role"`
	LastLoginAt       *string `json:"lastLoginAt"`
	CreatedAt         string  `json:"createdAt"`
	UpdatedAt         string  `json:"updatedAt"`
}

type authSessionResponse struct {
	ID        string `json:"id"`
	Token     string `json:"token"`
	UserID    string `json:"userId"`
	ExpiresAt string `json:"expiresAt"`
	IPAddress string `json:"ipAddress,omitempty"`
	UserAgent string `json:"userAgent,omitempty"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type authSessionPayload struct {
	Session authSessionResponse `json:"session"`
	User    authUserResponse    `json:"user"`
}

type signUpEmailResponse struct {
	User                  authUserResponse `json:"user"`
	VerificationRequired  bool             `json:"verificationRequired"`
	VerificationCode      string           `json:"verificationCode,omitempty"`
	VerificationExpiresAt string           `json:"verificationExpiresAt,omitempty"`
}

type oauthStartResponse struct {
	AuthorizationURL string `json:"authorizationUrl"`
	State            string `json:"state"`
	StateExpiresAt   string `json:"stateExpiresAt"`
}

type sendVerificationOTPResponse struct {
	OK                    bool   `json:"ok"`
	VerificationCode      string `json:"verificationCode,omitempty"`
	VerificationExpiresAt string `json:"verificationExpiresAt,omitempty"`
}

type verifyEmailOTPResponse struct {
	OK   bool             `json:"ok"`
	User authUserResponse `json:"user"`
}

type requestPasswordResetResponse struct {
	OK             bool   `json:"ok"`
	ResetToken     string `json:"resetToken,omitempty"`
	ResetExpiresAt string `json:"resetExpiresAt,omitempty"`
}

type okResponse struct {
	OK bool `json:"ok"`
}

type signOutResponse struct {
	Success bool `json:"success"`
}

type OAuthProfile struct {
	ProviderID string
	AccountID  string
	Email      string
	Name       string
	Image      *string
	Verified   bool
}

func newAuthHandler(options RouterOptions) *authHandler {
	return &authHandler{
		db:                  options.DB,
		betterAuthSecret:    options.BetterAuthSecret,
		authDeliveryToken:   options.AuthDeliveryToken,
		internalAPIToken:    options.InternalAPIToken,
		authSessionTTL:      options.AuthSessionTTL,
		authVerificationTTL: options.AuthVerificationTTL,
		authCookieSecure:    options.AuthCookieSecure,
		googleClientID:      options.GoogleClientID,
		googleClientSecret:  options.GoogleClientSecret,
		gitHubClientID:      options.GitHubClientID,
		gitHubClientSecret:  options.GitHubClientSecret,
		oauthProfileFetcher: options.OAuthProfileFetcher,
		httpClient:          http.DefaultClient,
	}
}

func (h *authHandler) requireTrustedInternalRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		if h.internalAPIToken == "" {
			writeError(c, http.StatusServiceUnavailable, "internal API token is not configured")
			c.Abort()
			return
		}
		if c.GetHeader(internalAPIHeader) != h.internalAPIToken {
			writeError(c, http.StatusUnauthorized, "trusted internal request required")
			c.Abort()
			return
		}
		c.Next()
	}
}

func (h *authHandler) signUpEmail(c *gin.Context) {
	if !h.authConfigured(c) {
		return
	}
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
	var verificationCode string
	var verificationExpiresAt time.Time
	if err := h.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("email = ?", email).Limit(1).Find(&user)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected > 0 {
			responseUser = genericAuthUserJSON(name, email)
			if user.EmailVerified {
				return nil
			}
			return h.rotateEmailVerification(tx, email, &verificationCode, &verificationExpiresAt)
		}

		hash, err := auth.HashPassword(req.Password)
		if err != nil {
			return err
		}
		now := time.Now().UTC()
		user = auth.User{
			Name:          name,
			Email:         email,
			Image:         req.Image,
			EmailVerified: false,
			CreatedAt:     now,
			UpdatedAt:     now,
		}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		account := auth.AuthAccount{
			AccountID:  user.ID,
			ProviderID: auth.CredentialProviderID,
			UserID:     user.ID,
			Password:   hash,
			CreatedAt:  now,
			UpdatedAt:  now,
		}
		if err := tx.Create(&account).Error; err != nil {
			return err
		}
		responseUser = authUserJSON(user)
		return h.rotateEmailVerification(tx, email, &verificationCode, &verificationExpiresAt)
	}); err != nil {
		writeError(c, http.StatusInternalServerError, "failed to sign up")
		return
	}

	out := signUpEmailResponse{
		User:                 responseUser,
		VerificationRequired: true,
	}
	if h.trustedAuthDeliveryRequest(c) && verificationCode != "" {
		out.VerificationCode = verificationCode
		out.VerificationExpiresAt = verificationExpiresAt.UTC().Format(time.RFC3339Nano)
	}
	c.JSON(http.StatusOK, out)
}

func (h *authHandler) signInEmail(c *gin.Context) {
	if !h.authConfigured(c) {
		return
	}
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
	result := h.db.Where("email = ?", email).Limit(1).Find(&user)
	if result.Error != nil {
		writeError(c, http.StatusInternalServerError, "failed to load user")
		return
	}
	if result.RowsAffected == 0 {
		writeError(c, http.StatusUnauthorized, "invalid email or password")
		return
	}

	var account auth.AuthAccount
	result = h.db.Where("user_id = ? AND provider_id = ?", user.ID, auth.CredentialProviderID).Limit(1).Find(&account)
	if result.Error != nil {
		writeError(c, http.StatusInternalServerError, "failed to load auth account")
		return
	}
	if result.RowsAffected == 0 || !auth.VerifyPassword(req.Password, account.Password) {
		writeError(c, http.StatusUnauthorized, "invalid email or password")
		return
	}
	if !user.EmailVerified {
		writeError(c, http.StatusForbidden, "email is not verified")
		return
	}

	session, err := h.createSession(c, user)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to create session")
		return
	}
	remember := true
	if req.RememberMe != nil {
		remember = *req.RememberMe
	}
	h.setSessionCookie(c, session.Token, session.ExpiresAt, remember)
	c.JSON(http.StatusOK, authSessionPayload{
		Session: authSessionJSON(session),
		User:    authUserJSON(user),
	})
}

func (h *authHandler) getSession(c *gin.Context) {
	if !h.authConfigured(c) {
		return
	}
	session, ok, err := h.loadAuthenticatedSession(c)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}
	if !ok {
		c.JSON(http.StatusOK, nil)
		return
	}
	c.JSON(http.StatusOK, authSessionPayload{
		Session: authSessionJSON(session),
		User:    authUserJSON(session.User),
	})
}

func (h *authHandler) signOut(c *gin.Context) {
	if !h.authConfigured(c) {
		return
	}
	token := sessionTokenFromRequest(c)
	if token != "" {
		if err := h.db.Where("token = ?", token).Delete(&auth.Session{}).Error; err != nil {
			writeError(c, http.StatusInternalServerError, "failed to delete session")
			return
		}
	}
	h.clearSessionCookie(c)
	c.JSON(http.StatusOK, signOutResponse{Success: true})
}

func (h *authHandler) sendVerificationOTP(c *gin.Context) {
	if !h.authConfigured(c) {
		return
	}
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
	result := h.db.Where("email = ?", email).Limit(1).Find(&user)
	if result.Error != nil {
		writeError(c, http.StatusInternalServerError, "failed to load user")
		return
	}
	if result.RowsAffected == 0 || user.EmailVerified {
		c.JSON(http.StatusOK, sendVerificationOTPResponse{OK: true})
		return
	}

	var code string
	var expiresAt time.Time
	if err := h.db.Transaction(func(tx *gorm.DB) error {
		return h.rotateEmailVerification(tx, email, &code, &expiresAt)
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

func (h *authHandler) verifyEmailOTP(c *gin.Context) {
	if !h.authConfigured(c) {
		return
	}
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
	invalid := false
	if err := h.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("email = ?", email).Limit(1).Find(&user)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			invalid = true
			return nil
		}
		verification, ok, err := h.loadVerification(tx, verificationIdentifier(email))
		if err != nil {
			return err
		}
		if !ok || !verification.ExpiresAt.After(time.Now().UTC()) || !auth.VerifyCode(h.betterAuthSecret, verification.Identifier, otp, verification.Value) {
			if ok && !verification.ExpiresAt.After(time.Now().UTC()) {
				if err := tx.Delete(&verification).Error; err != nil {
					return err
				}
			}
			invalid = true
			return nil
		}
		if err := tx.Delete(&verification).Error; err != nil {
			return err
		}
		if err := markEmailVerified(tx, user.ID, time.Now().UTC()); err != nil {
			return err
		}
		return tx.First(&user, "id = ?", user.ID).Error
	}); err != nil {
		writeError(c, http.StatusInternalServerError, "failed to verify email")
		return
	}
	if invalid {
		writeError(c, http.StatusBadRequest, "invalid verification code")
		return
	}
	c.JSON(http.StatusOK, verifyEmailOTPResponse{OK: true, User: authUserJSON(user)})
}

func (h *authHandler) requestPasswordReset(c *gin.Context) {
	if !h.authConfigured(c) {
		return
	}
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
	result := h.db.Where("email = ?", email).Limit(1).Find(&user)
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
	if err := h.db.Transaction(func(tx *gorm.DB) error {
		return h.rotatePasswordResetVerification(tx, email, &token, &expiresAt)
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

func (h *authHandler) confirmPasswordReset(c *gin.Context) {
	if !h.authConfigured(c) {
		return
	}
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
	if err := validatePassword(req.Password); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to hash password")
		return
	}

	invalid := false
	if err := h.db.Transaction(func(tx *gorm.DB) error {
		var user auth.User
		result := tx.Where("email = ?", email).Limit(1).Find(&user)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			invalid = true
			return nil
		}
		verification, ok, err := h.loadVerification(tx, passwordResetIdentifier(email))
		if err != nil {
			return err
		}
		if !ok || !verification.ExpiresAt.After(time.Now().UTC()) || !auth.VerifyCode(h.betterAuthSecret, verification.Identifier, token, verification.Value) {
			if ok && !verification.ExpiresAt.After(time.Now().UTC()) {
				if err := tx.Delete(&verification).Error; err != nil {
					return err
				}
			}
			invalid = true
			return nil
		}
		if err := tx.Delete(&verification).Error; err != nil {
			return err
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
		if err := upsertAccount(tx, account); err != nil {
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
	if invalid {
		writeError(c, http.StatusBadRequest, "invalid password reset token")
		return
	}
	c.JSON(http.StatusOK, okResponse{OK: true})
}

func (h *authHandler) startGoogleOAuth(c *gin.Context) {
	h.startOAuth(c, "google")
}

func (h *authHandler) startGitHubOAuth(c *gin.Context) {
	h.startOAuth(c, "github")
}

func (h *authHandler) startOAuth(c *gin.Context, provider string) {
	var req oauthStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	redirectURI := strings.TrimSpace(req.RedirectURI)
	if redirectURI == "" {
		writeError(c, http.StatusBadRequest, "redirectURI is required")
		return
	}
	state, err := auth.GenerateSessionToken()
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to create oauth state")
		return
	}
	expiresAt := time.Now().UTC().Add(oauthStateTTL)
	var authorizationURL string
	switch provider {
	case "google":
		if h.googleClientID == "" || h.googleClientSecret == "" {
			writeError(c, http.StatusServiceUnavailable, "google oauth is not configured")
			return
		}
		values := url.Values{
			"client_id":     {h.googleClientID},
			"redirect_uri":  {redirectURI},
			"response_type": {"code"},
			"scope":         {"openid email profile"},
			"state":         {state},
			"access_type":   {"offline"},
			"prompt":        {"select_account"},
		}
		authorizationURL = "https://accounts.google.com/o/oauth2/v2/auth?" + values.Encode()
	case "github":
		if h.gitHubClientID == "" || h.gitHubClientSecret == "" {
			writeError(c, http.StatusServiceUnavailable, "github oauth is not configured")
			return
		}
		values := url.Values{
			"client_id":    {h.gitHubClientID},
			"redirect_uri": {redirectURI},
			"scope":        {"read:user user:email"},
			"state":        {state},
		}
		authorizationURL = "https://github.com/login/oauth/authorize?" + values.Encode()
	default:
		writeError(c, http.StatusBadRequest, "unsupported oauth provider")
		return
	}
	c.JSON(http.StatusOK, oauthStartResponse{
		AuthorizationURL: authorizationURL,
		State:            state,
		StateExpiresAt:   expiresAt.Format(time.RFC3339Nano),
	})
}

func (h *authHandler) signInGoogleOAuth(c *gin.Context) {
	h.signInOAuth(c, auth.GoogleProviderID)
}

func (h *authHandler) signInGitHubOAuth(c *gin.Context) {
	h.signInOAuth(c, auth.GitHubProviderID)
}

func (h *authHandler) signInOAuth(c *gin.Context, providerID string) {
	if !h.authConfigured(c) {
		return
	}
	var req oauthCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if strings.TrimSpace(req.Code) == "" || strings.TrimSpace(req.RedirectURI) == "" {
		writeError(c, http.StatusBadRequest, "code and redirectURI are required")
		return
	}

	profile, err := h.fetchOAuthProfile(c.Request.Context(), providerID, req.Code, req.RedirectURI)
	if err != nil {
		writeError(c, http.StatusUnauthorized, "oauth sign-in failed")
		return
	}
	user, err := h.upsertOAuthUser(profile)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to save oauth user")
		return
	}
	session, err := h.createSession(c, user)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to create session")
		return
	}
	h.setSessionCookie(c, session.Token, session.ExpiresAt, true)
	c.JSON(http.StatusOK, authSessionPayload{
		Session: authSessionJSON(session),
		User:    authUserJSON(user),
	})
}

func (h *authHandler) authConfigured(c *gin.Context) bool {
	if h.db == nil {
		writeError(c, http.StatusServiceUnavailable, "authentication database is not configured")
		return false
	}
	if strings.TrimSpace(h.betterAuthSecret) == "" {
		writeError(c, http.StatusServiceUnavailable, "BETTER_AUTH_SECRET is not configured")
		return false
	}
	return true
}

func (h *authHandler) createSession(c *gin.Context, user auth.User) (auth.Session, error) {
	token, err := auth.GenerateSessionToken()
	if err != nil {
		return auth.Session{}, err
	}
	now := time.Now().UTC()
	session := auth.Session{
		Token:     token,
		UserID:    user.ID,
		ExpiresAt: now.Add(h.sessionTTL()),
		IPAddress: c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
		CreatedAt: now,
		UpdatedAt: now,
	}
	err = h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&auth.User{}).Where("id = ?", user.ID).Update("last_login_at", now).Error; err != nil {
			return err
		}
		user.LastLoginAt = &now
		return tx.Create(&session).Error
	})
	return session, err
}

func (h *authHandler) loadAuthenticatedSession(c *gin.Context) (auth.Session, bool, error) {
	token, tokenPresent := sessionTokenFromRequestWithPresence(c)
	if token == "" {
		if tokenPresent {
			h.clearSessionCookie(c)
		}
		return auth.Session{}, false, nil
	}
	var session auth.Session
	result := h.db.Preload("User").Where("token = ?", token).Limit(1).Find(&session)
	if result.Error != nil {
		return auth.Session{}, false, errors.New("failed to load session")
	}
	if result.RowsAffected == 0 {
		h.clearSessionCookie(c)
		return auth.Session{}, false, nil
	}
	if !session.ExpiresAt.After(time.Now().UTC()) {
		if err := h.db.Delete(&session).Error; err != nil {
			return auth.Session{}, false, errors.New("failed to delete expired session")
		}
		h.clearSessionCookie(c)
		return auth.Session{}, false, nil
	}
	if session.User.ID == "" || session.User.Email == "" {
		if err := h.db.Where("token = ?", token).Delete(&auth.Session{}).Error; err != nil {
			return auth.Session{}, false, errors.New("failed to delete stale session")
		}
		h.clearSessionCookie(c)
		return auth.Session{}, false, nil
	}
	return session, true, nil
}

func (h *authHandler) rotateEmailVerification(tx *gorm.DB, email string, codeOut *string, expiresOut *time.Time) error {
	code, err := auth.GenerateNumericCode(6)
	if err != nil {
		return err
	}
	identifier := verificationIdentifier(email)
	return h.rotateVerification(tx, identifier, code, codeOut, expiresOut)
}

func (h *authHandler) rotatePasswordResetVerification(tx *gorm.DB, email string, tokenOut *string, expiresOut *time.Time) error {
	token, err := auth.GenerateSessionToken()
	if err != nil {
		return err
	}
	identifier := passwordResetIdentifier(email)
	return h.rotateVerification(tx, identifier, token, tokenOut, expiresOut)
}

func (h *authHandler) rotateVerification(tx *gorm.DB, identifier string, plaintext string, plaintextOut *string, expiresOut *time.Time) error {
	now := time.Now().UTC()
	expiresAt := now.Add(h.verificationTTL())
	value := auth.HashCode(h.betterAuthSecret, identifier, plaintext)
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
	*plaintextOut = plaintext
	*expiresOut = expiresAt
	return nil
}

func (h *authHandler) loadVerification(tx *gorm.DB, identifier string) (auth.Verification, bool, error) {
	var verification auth.Verification
	result := tx.Where("identifier = ?", identifier).Limit(1).Find(&verification)
	if result.Error != nil {
		return auth.Verification{}, false, result.Error
	}
	return verification, result.RowsAffected > 0, nil
}

func (h *authHandler) upsertOAuthUser(profile OAuthProfile) (auth.User, error) {
	email := normalizeEmail(profile.Email)
	if !validEmail(email) || profile.AccountID == "" || profile.ProviderID == "" {
		return auth.User{}, errors.New("invalid oauth profile")
	}
	now := time.Now().UTC()
	var user auth.User
	err := h.db.Transaction(func(tx *gorm.DB) error {
		var account auth.AuthAccount
		result := tx.Preload("User").Where("provider_id = ? AND account_id = ?", profile.ProviderID, profile.AccountID).Limit(1).Find(&account)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected > 0 {
			user = account.User
			if user.ID == "" {
				return errors.New("oauth account has no user")
			}
			return nil
		}
		result = tx.Where("email = ?", email).Limit(1).Find(&user)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			name := strings.TrimSpace(profile.Name)
			if name == "" {
				name = email
			}
			user = auth.User{
				Name:          name,
				Email:         email,
				EmailVerified: profile.Verified,
				Image:         profile.Image,
				CreatedAt:     now,
				UpdatedAt:     now,
			}
			if err := tx.Create(&user).Error; err != nil {
				return err
			}
		} else if !user.EmailVerified && profile.Verified {
			if err := markEmailVerified(tx, user.ID, now); err != nil {
				return err
			}
			if err := tx.First(&user, "id = ?", user.ID).Error; err != nil {
				return err
			}
		}
		account = auth.AuthAccount{
			AccountID:  profile.AccountID,
			ProviderID: profile.ProviderID,
			UserID:     user.ID,
			CreatedAt:  now,
			UpdatedAt:  now,
		}
		return tx.Create(&account).Error
	})
	return user, err
}

func upsertAccount(tx *gorm.DB, account auth.AuthAccount) error {
	return tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "provider_id"},
			{Name: "account_id"},
		},
		DoUpdates: clause.Assignments(map[string]any{
			"password":   account.Password,
			"updated_at": account.UpdatedAt,
		}),
	}).Create(&account).Error
}

func markEmailVerified(tx *gorm.DB, userID string, now time.Time) error {
	return tx.Model(&auth.User{}).Where("id = ?", userID).Updates(map[string]any{
		"email_verified": true,
		"status":         gorm.Expr("CASE WHEN status = ? THEN ? ELSE status END", "invited", "active"),
		"updated_at":     now,
	}).Error
}

func (h *authHandler) fetchOAuthProfile(ctx context.Context, providerID string, code string, redirectURI string) (OAuthProfile, error) {
	if h.oauthProfileFetcher != nil {
		return h.oauthProfileFetcher(ctx, providerID, code, redirectURI)
	}
	switch providerID {
	case auth.GoogleProviderID:
		return h.fetchGoogleProfile(ctx, code, redirectURI)
	case auth.GitHubProviderID:
		return h.fetchGitHubProfile(ctx, code, redirectURI)
	default:
		return OAuthProfile{}, errors.New("unsupported oauth provider")
	}
}

func (h *authHandler) fetchGoogleProfile(ctx context.Context, code string, redirectURI string) (OAuthProfile, error) {
	if h.googleClientID == "" || h.googleClientSecret == "" {
		return OAuthProfile{}, errors.New("google oauth is not configured")
	}
	var tokenBody struct {
		AccessToken string `json:"access_token"`
	}
	if err := h.postFormJSON(ctx, "https://oauth2.googleapis.com/token", url.Values{
		"client_id":     {h.googleClientID},
		"client_secret": {h.googleClientSecret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {redirectURI},
	}, &tokenBody); err != nil {
		return OAuthProfile{}, err
	}
	var profile struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
	}
	if err := h.getBearerJSON(ctx, "https://www.googleapis.com/oauth2/v2/userinfo", tokenBody.AccessToken, &profile); err != nil {
		return OAuthProfile{}, err
	}
	return OAuthProfile{
		ProviderID: auth.GoogleProviderID,
		AccountID:  profile.ID,
		Email:      profile.Email,
		Name:       profile.Name,
		Image:      stringPtr(profile.Picture),
		Verified:   profile.VerifiedEmail,
	}, nil
}

func (h *authHandler) fetchGitHubProfile(ctx context.Context, code string, redirectURI string) (OAuthProfile, error) {
	if h.gitHubClientID == "" || h.gitHubClientSecret == "" {
		return OAuthProfile{}, errors.New("github oauth is not configured")
	}
	var tokenBody struct {
		AccessToken string `json:"access_token"`
	}
	if err := h.postFormJSON(ctx, "https://github.com/login/oauth/access_token", url.Values{
		"client_id":     {h.gitHubClientID},
		"client_secret": {h.gitHubClientSecret},
		"code":          {code},
		"redirect_uri":  {redirectURI},
	}, &tokenBody); err != nil {
		return OAuthProfile{}, err
	}
	var profile struct {
		ID        int64  `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := h.getBearerJSON(ctx, "https://api.github.com/user", tokenBody.AccessToken, &profile); err != nil {
		return OAuthProfile{}, err
	}
	email := profile.Email
	if email == "" {
		email = profile.Login + "@users.noreply.github.com"
	}
	name := profile.Name
	if strings.TrimSpace(name) == "" {
		name = profile.Login
	}
	return OAuthProfile{
		ProviderID: auth.GitHubProviderID,
		AccountID:  fmt.Sprint(profile.ID),
		Email:      email,
		Name:       name,
		Image:      stringPtr(profile.AvatarURL),
		Verified:   true,
	}, nil
}

func (h *authHandler) postFormJSON(ctx context.Context, endpoint string, values url.Values, out any) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(values.Encode()))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Accept", "application/json")
	return h.doJSON(request, out)
}

func (h *authHandler) getBearerJSON(ctx context.Context, endpoint string, token string, out any) error {
	if strings.TrimSpace(token) == "" {
		return errors.New("missing access token")
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Accept", "application/json")
	return h.doJSON(request, out)
}

func (h *authHandler) doJSON(request *http.Request, out any) error {
	client := h.httpClient
	if client == nil {
		client = http.DefaultClient
	}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return err
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("oauth response status %d", response.StatusCode)
	}
	return json.Unmarshal(body, out)
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func verificationIdentifier(email string) string {
	return "email-verification:" + normalizeEmail(email)
}

func passwordResetIdentifier(email string) string {
	return "password-reset:" + normalizeEmail(email)
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
	return validatePassword(password)
}

func validatePassword(password string) error {
	characters := utf8.RuneCountInString(password)
	if characters < minPasswordLength || characters > maxPasswordLength {
		return errors.New("password must be between 8 and 128 characters")
	}
	return nil
}

func validEmail(email string) bool {
	if len(email) > maxAuthEmailBytes {
		return false
	}
	parsed, err := mail.ParseAddress(email)
	return err == nil && parsed.Address == email
}

func authUserJSON(user auth.User) authUserResponse {
	var lastLoginAt *string
	if user.LastLoginAt != nil {
		value := user.LastLoginAt.UTC().Format(time.RFC3339Nano)
		lastLoginAt = &value
	}
	role := string(user.Role)
	if role == "" {
		role = string(auth.UserRoleUser)
	}
	preferredLanguage := user.PreferredLanguage
	if preferredLanguage == "" {
		preferredLanguage = "en"
	}
	return authUserResponse{
		ID:                user.ID,
		Name:              user.Name,
		Email:             user.Email,
		EmailVerified:     user.EmailVerified,
		Image:             user.Image,
		PreferredLanguage: preferredLanguage,
		Role:              role,
		LastLoginAt:       lastLoginAt,
		CreatedAt:         user.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:         user.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
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

func authSessionJSON(session auth.Session) authSessionResponse {
	return authSessionResponse{
		ID:        session.ID,
		Token:     session.Token,
		UserID:    session.UserID,
		ExpiresAt: session.ExpiresAt.UTC().Format(time.RFC3339Nano),
		IPAddress: session.IPAddress,
		UserAgent: session.UserAgent,
		CreatedAt: session.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt: session.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
}

func (h *authHandler) trustedAuthDeliveryRequest(c *gin.Context) bool {
	return h.authDeliveryToken != "" && c.GetHeader(authDeliveryHeader) == h.authDeliveryToken
}

func (h *authHandler) sessionTTL() time.Duration {
	if h.authSessionTTL > 0 {
		return h.authSessionTTL
	}
	return 7 * 24 * time.Hour
}

func (h *authHandler) verificationTTL() time.Duration {
	if h.authVerificationTTL > 0 {
		return h.authVerificationTTL
	}
	return 5 * time.Minute
}

func (h *authHandler) setSessionCookie(c *gin.Context, token string, expiresAt time.Time, remember bool) {
	cookie := &http.Cookie{
		Name:     authSessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   h.authCookieSecure,
		SameSite: http.SameSiteLaxMode,
	}
	if remember {
		cookie.MaxAge = max(0, int(time.Until(expiresAt).Seconds()))
		cookie.Expires = expiresAt
	}
	http.SetCookie(c.Writer, cookie)
}

func (h *authHandler) clearSessionCookie(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     authSessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   h.authCookieSecure,
		SameSite: http.SameSiteLaxMode,
	})
}

func sessionTokenFromRequest(c *gin.Context) string {
	token, _ := sessionTokenFromRequestWithPresence(c)
	return token
}

func sessionTokenFromRequestWithPresence(c *gin.Context) (string, bool) {
	cookie, err := c.Cookie(authSessionCookieName)
	if err != nil {
		return "", false
	}
	return strings.TrimSpace(cookie), true
}

func stringPtr(value string) *string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return &value
}
