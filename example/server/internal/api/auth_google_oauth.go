package api

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleidtoken "google.golang.org/api/idtoken"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	auth "ai.ro/syncra/internal/auth"
	"ai.ro/syncra/internal/billing"
)

const googleOAuthStatePrefix = "oauth-state:google:"
const googleOAuthLinkStatePrefix = "oauth-link-state:google:"

var (
	errGoogleOAuthNotConfigured = errors.New("google oauth is not configured")
	errInvalidGoogleOAuthState  = errors.New("invalid google oauth state")
	errInvalidGoogleIdentity    = errors.New("invalid google identity")
)

// swagger:model googleOAuthStartRequest
type googleOAuthStartRequest struct {
	RedirectURI string `json:"redirectURI"`
}

// swagger:model googleOAuthCallbackRequest
type googleOAuthCallbackRequest struct {
	Code        string `json:"code"`
	State       string `json:"state"`
	RedirectURI string `json:"redirectURI"`
}

// swagger:model googleOAuthStartResponse
type googleOAuthStartResponse struct {
	AuthorizationURL string `json:"authorizationUrl"`
	State            string `json:"state"`
	StateExpiresAt   string `json:"stateExpiresAt"`
}

type googleOAuthToken struct {
	AccessToken           *string
	RefreshToken          *string
	IDToken               string
	AccessTokenExpiresAt  *time.Time
	RefreshTokenExpiresAt *time.Time
	Scope                 *string
}

type googleIDTokenPayload struct {
	Subject       string
	Email         string
	EmailVerified bool
	Name          string
	Picture       string
}

func (h *Handler) StartGoogleOAuth(c *gin.Context) {
	var req googleOAuthStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	redirectURI, err := normalizeOAuthRedirectURI(req.RedirectURI)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	if !h.googleOAuthConfigured() {
		writeError(c, http.StatusServiceUnavailable, errGoogleOAuthNotConfigured.Error())
		return
	}

	var state string
	var expiresAt time.Time
	if err := h.DB.Transaction(func(tx *gorm.DB) error {
		return rotateGoogleOAuthState(tx, h.BetterAuthSecret, redirectURI, h.authVerificationTTL(), &state, &expiresAt)
	}); err != nil {
		writeError(c, http.StatusInternalServerError, "failed to create google oauth state")
		return
	}

	c.JSON(http.StatusOK, googleOAuthStartResponse{
		AuthorizationURL: h.googleOAuthConfig(redirectURI).AuthCodeURL(state),
		State:            state,
		StateExpiresAt:   expiresAt.UTC().Format(time.RFC3339Nano),
	})
}

func (h *Handler) SignInGoogleOAuth(c *gin.Context) {
	var req googleOAuthCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	code := strings.TrimSpace(req.Code)
	state := strings.TrimSpace(req.State)
	redirectURI, err := normalizeOAuthRedirectURI(req.RedirectURI)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	if code == "" {
		writeError(c, http.StatusBadRequest, "authorization code is required")
		return
	}
	if state == "" {
		writeError(c, http.StatusBadRequest, "state is required")
		return
	}
	if !h.googleOAuthConfigured() {
		writeError(c, http.StatusServiceUnavailable, errGoogleOAuthNotConfigured.Error())
		return
	}

	if err := h.consumeGoogleOAuthState(state, redirectURI); err != nil {
		if errors.Is(err, errInvalidGoogleOAuthState) {
			writeError(c, http.StatusBadRequest, errInvalidGoogleOAuthState.Error())
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to verify google oauth state")
		return
	}

	token, err := h.exchangeGoogleOAuthCode(c.Request.Context(), code, redirectURI)
	if err != nil {
		writeError(c, http.StatusBadGateway, "failed to exchange google authorization code")
		return
	}
	profile, err := h.validateGoogleIDToken(c.Request.Context(), token.IDToken)
	if err != nil {
		writeError(c, http.StatusUnauthorized, errInvalidGoogleIdentity.Error())
		return
	}

	session, user, err := h.signInGoogleProfile(c, profile, token)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to sign in with google")
		return
	}

	h.setSessionCookie(c, session.Token, session.ExpiresAt, true)
	c.JSON(http.StatusOK, authSessionPayload{
		Session: authSessionJSON(session),
		User:    authUserJSON(user),
	})
}

func (h *Handler) StartGoogleAccountLink(c *gin.Context) {
	if _, ok, err := h.loadAuthenticatedSession(c); err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	} else if !ok {
		writeError(c, http.StatusUnauthorized, "authentication required")
		return
	}

	var req googleOAuthStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	redirectURI, err := normalizeOAuthRedirectURI(req.RedirectURI)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	if !h.googleOAuthConfigured() {
		writeError(c, http.StatusServiceUnavailable, errGoogleOAuthNotConfigured.Error())
		return
	}

	var state string
	var expiresAt time.Time
	if err := h.DB.Transaction(func(tx *gorm.DB) error {
		return rotateGoogleOAuthLinkState(tx, h.BetterAuthSecret, redirectURI, h.authVerificationTTL(), &state, &expiresAt)
	}); err != nil {
		writeError(c, http.StatusInternalServerError, "failed to create google oauth state")
		return
	}

	c.JSON(http.StatusOK, googleOAuthStartResponse{
		AuthorizationURL: h.googleOAuthConfig(redirectURI).AuthCodeURL(state),
		State:            state,
		StateExpiresAt:   expiresAt.UTC().Format(time.RFC3339Nano),
	})
}

func (h *Handler) LinkGoogleAccount(c *gin.Context) {
	session, ok, err := h.loadAuthenticatedSession(c)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}
	if !ok {
		writeError(c, http.StatusUnauthorized, "authentication required")
		return
	}

	var req googleOAuthCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	code := strings.TrimSpace(req.Code)
	state := strings.TrimSpace(req.State)
	redirectURI, err := normalizeOAuthRedirectURI(req.RedirectURI)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	if code == "" {
		writeError(c, http.StatusBadRequest, "authorization code is required")
		return
	}
	if state == "" {
		writeError(c, http.StatusBadRequest, "state is required")
		return
	}
	if !h.googleOAuthConfigured() {
		writeError(c, http.StatusServiceUnavailable, errGoogleOAuthNotConfigured.Error())
		return
	}

	if err := h.consumeGoogleOAuthLinkState(state, redirectURI); err != nil {
		if errors.Is(err, errInvalidGoogleOAuthState) {
			writeError(c, http.StatusBadRequest, errInvalidGoogleOAuthState.Error())
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to verify google oauth state")
		return
	}

	token, err := h.exchangeGoogleOAuthCode(c.Request.Context(), code, redirectURI)
	if err != nil {
		writeError(c, http.StatusBadGateway, "failed to exchange google authorization code")
		return
	}
	profile, err := h.validateGoogleIDToken(c.Request.Context(), token.IDToken)
	if err != nil {
		writeError(c, http.StatusUnauthorized, errInvalidGoogleIdentity.Error())
		return
	}

	account, err := h.linkGoogleProfile(c, session.UserID, profile, token)
	if err != nil {
		if errors.Is(err, errOAuthAccountLinked) || errors.Is(err, errProviderAlreadyLinked) {
			writeError(c, http.StatusConflict, err.Error())
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to link google account")
		return
	}

	c.JSON(http.StatusOK, authAccountListItemJSON(account))
}

func (h *Handler) googleOAuthConfigured() bool {
	return strings.TrimSpace(h.GoogleClientID) != "" && strings.TrimSpace(h.GoogleClientSecret) != ""
}

func (h *Handler) googleOAuthConfig(redirectURI string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     strings.TrimSpace(h.GoogleClientID),
		ClientSecret: strings.TrimSpace(h.GoogleClientSecret),
		RedirectURL:  redirectURI,
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}
}

func (h *Handler) exchangeGoogleOAuthCode(ctx context.Context, code string, redirectURI string) (googleOAuthToken, error) {
	if h.GoogleOAuthExchange != nil {
		return h.GoogleOAuthExchange(ctx, code, redirectURI)
	}

	oauthToken, err := h.googleOAuthConfig(redirectURI).Exchange(ctx, code)
	if err != nil {
		return googleOAuthToken{}, err
	}
	idToken, ok := oauthToken.Extra("id_token").(string)
	if !ok || strings.TrimSpace(idToken) == "" {
		return googleOAuthToken{}, errors.New("google token response missing id_token")
	}

	return googleOAuthToken{
		AccessToken:           optionalString(oauthToken.AccessToken),
		RefreshToken:          optionalString(oauthToken.RefreshToken),
		IDToken:               idToken,
		AccessTokenExpiresAt:  optionalTime(oauthToken.Expiry),
		RefreshTokenExpiresAt: optionalTime(time.Time{}),
		Scope:                 optionalString(tokenExtraString(oauthToken, "scope")),
	}, nil
}

func (h *Handler) validateGoogleIDToken(ctx context.Context, rawIDToken string) (googleIDTokenPayload, error) {
	if h.GoogleIDTokenValidate != nil {
		profile, err := h.GoogleIDTokenValidate(ctx, rawIDToken, strings.TrimSpace(h.GoogleClientID))
		if err != nil {
			return googleIDTokenPayload{}, err
		}
		return validateGoogleProfile(profile)
	}
	payload, err := googleidtoken.Validate(ctx, rawIDToken, strings.TrimSpace(h.GoogleClientID))
	if err != nil {
		return googleIDTokenPayload{}, err
	}

	return googlePayloadFromClaims(payload)
}

func (h *Handler) consumeGoogleOAuthState(state string, redirectURI string) error {
	identifier := googleOAuthStateIdentifier(state)
	var verification auth.Verification
	if err := h.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("identifier = ?", identifier).Limit(1).Find(&verification)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errInvalidGoogleOAuthState
		}
		if !verification.ExpiresAt.After(h.now()) {
			if _, err := consumeGoogleOAuthStateVerification(tx, verification, identifier); err != nil {
				return err
			}
			return errInvalidGoogleOAuthState
		}
		if !auth.VerifyCode(h.BetterAuthSecret, identifier, redirectURI, verification.Value) {
			return errInvalidGoogleOAuthState
		}
		consumed, err := consumeGoogleOAuthStateVerification(tx, verification, identifier)
		if err != nil {
			return err
		}
		if !consumed {
			return errInvalidGoogleOAuthState
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (h *Handler) consumeGoogleOAuthLinkState(state string, redirectURI string) error {
	identifier := googleOAuthLinkStateIdentifier(state)
	var verification auth.Verification
	if err := h.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("identifier = ?", identifier).Limit(1).Find(&verification)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errInvalidGoogleOAuthState
		}
		if !verification.ExpiresAt.After(h.now()) {
			if _, err := consumeGoogleOAuthStateVerification(tx, verification, identifier); err != nil {
				return err
			}
			return errInvalidGoogleOAuthState
		}
		if !auth.VerifyCode(h.BetterAuthSecret, identifier, redirectURI, verification.Value) {
			return errInvalidGoogleOAuthState
		}
		consumed, err := consumeGoogleOAuthStateVerification(tx, verification, identifier)
		if err != nil {
			return err
		}
		if !consumed {
			return errInvalidGoogleOAuthState
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (h *Handler) signInGoogleProfile(c *gin.Context, profile googleIDTokenPayload, token googleOAuthToken) (auth.Session, auth.User, error) {
	now := h.now()
	var session auth.Session
	var user auth.User

	err := h.DB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		account, found, err := googleAccountForSubject(tx, profile.Subject)
		if err != nil {
			return err
		}
		newlyVerified := false
		if found {
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&user, "id = ?", account.UserID).Error; err != nil {
				return err
			}
			if err := updateGoogleAccountTokens(tx, account, token, now); err != nil {
				return err
			}
		} else {
			result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("email = ?", profile.Email).Limit(1).Find(&user)
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				image := optionalString(profile.Picture)
				user = auth.User{
					Name:          profileName(profile),
					Email:         profile.Email,
					EmailVerified: true,
					Image:         image,
					CreatedAt:     now,
					UpdatedAt:     now,
				}
				if err := tx.Create(&user).Error; err != nil {
					return err
				}
				newlyVerified = true
			} else if !user.EmailVerified {
				if err := tx.Model(&auth.User{}).Where("id = ?", user.ID).Updates(map[string]any{
					"email_verified": true,
					"updated_at":     now,
				}).Error; err != nil {
					return err
				}
				user.EmailVerified = true
				user.UpdatedAt = now
				newlyVerified = true
			}
			account := auth.AuthAccount{
				AccountID:             profile.Subject,
				ProviderID:            auth.GoogleProviderID,
				UserID:                user.ID,
				AccessToken:           token.AccessToken,
				RefreshToken:          token.RefreshToken,
				IDToken:               optionalString(token.IDToken),
				AccessTokenExpiresAt:  token.AccessTokenExpiresAt,
				RefreshTokenExpiresAt: token.RefreshTokenExpiresAt,
				Scope:                 token.Scope,
				CreatedAt:             now,
				UpdatedAt:             now,
			}
			if err := tx.Create(&account).Error; err != nil {
				return err
			}
		}

		if err := tx.Model(&auth.User{}).Where("id = ?", user.ID).Update("last_login_at", now).Error; err != nil {
			return err
		}
		user.LastLoginAt = &now
		if newlyVerified {
			if _, err := billing.GrantSignupBonus(c.Request.Context(), tx, billing.GrantSignupBonusInput{
				UserID:  user.ID,
				Credits: h.onboardingCredits(),
				Now:     now,
			}); err != nil {
				return err
			}
		}

		token, err := auth.GenerateSessionToken()
		if err != nil {
			return err
		}
		session = auth.Session{
			Token:     token,
			UserID:    user.ID,
			ExpiresAt: now.Add(h.authSessionTTL()),
			IPAddress: c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
			CreatedAt: now,
			UpdatedAt: now,
		}
		return tx.Create(&session).Error
	})
	return session, user, err
}

func (h *Handler) linkGoogleProfile(c *gin.Context, userID string, profile googleIDTokenPayload, token googleOAuthToken) (auth.AuthAccount, error) {
	now := h.now()
	var linked auth.AuthAccount

	err := h.DB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		account, found, err := googleAccountForSubject(tx, profile.Subject)
		if err != nil {
			return err
		}
		if found {
			if account.UserID != userID {
				return errOAuthAccountLinked
			}
			if err := updateGoogleAccountTokens(tx, account, token, now); err != nil {
				return err
			}
			return tx.First(&linked, "id = ?", account.ID).Error
		}

		var existing auth.AuthAccount
		result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ? AND provider_id = ?", userID, auth.GoogleProviderID).
			Limit(1).
			Find(&existing)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected > 0 {
			return errProviderAlreadyLinked
		}

		linked = auth.AuthAccount{
			AccountID:             profile.Subject,
			ProviderID:            auth.GoogleProviderID,
			UserID:                userID,
			AccessToken:           token.AccessToken,
			RefreshToken:          token.RefreshToken,
			IDToken:               optionalString(token.IDToken),
			AccessTokenExpiresAt:  token.AccessTokenExpiresAt,
			RefreshTokenExpiresAt: token.RefreshTokenExpiresAt,
			Scope:                 token.Scope,
			CreatedAt:             now,
			UpdatedAt:             now,
		}
		return tx.Create(&linked).Error
	})
	return linked, err
}

func googleAccountForSubject(tx *gorm.DB, subject string) (auth.AuthAccount, bool, error) {
	var account auth.AuthAccount
	result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("provider_id = ? AND account_id = ?", auth.GoogleProviderID, subject).
		Limit(1).
		Find(&account)
	if result.Error != nil {
		return auth.AuthAccount{}, false, result.Error
	}
	return account, result.RowsAffected > 0, nil
}

func updateGoogleAccountTokens(tx *gorm.DB, account auth.AuthAccount, token googleOAuthToken, now time.Time) error {
	updates := map[string]any{
		"access_token":            token.AccessToken,
		"id_token":                optionalString(token.IDToken),
		"access_token_expires_at": token.AccessTokenExpiresAt,
		"scope":                   token.Scope,
		"updated_at":              now,
	}
	if token.RefreshToken != nil {
		updates["refresh_token"] = token.RefreshToken
		updates["refresh_token_expires_at"] = token.RefreshTokenExpiresAt
	}
	return tx.Model(&auth.AuthAccount{}).Where("id = ?", account.ID).Updates(updates).Error
}

func rotateGoogleOAuthState(tx *gorm.DB, secret string, redirectURI string, ttl time.Duration, stateOut *string, expiresOut *time.Time) error {
	state, err := auth.GenerateSessionToken()
	if err != nil {
		return err
	}
	identifier := googleOAuthStateIdentifier(state)
	now := time.Now().UTC()
	expiresAt := now.Add(ttl)
	value := auth.HashCode(secret, identifier, redirectURI)
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
	*stateOut = state
	*expiresOut = expiresAt
	return nil
}

func rotateGoogleOAuthLinkState(tx *gorm.DB, secret string, redirectURI string, ttl time.Duration, stateOut *string, expiresOut *time.Time) error {
	state, err := auth.GenerateSessionToken()
	if err != nil {
		return err
	}
	identifier := googleOAuthLinkStateIdentifier(state)
	now := time.Now().UTC()
	expiresAt := now.Add(ttl)
	value := auth.HashCode(secret, identifier, redirectURI)
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
	*stateOut = state
	*expiresOut = expiresAt
	return nil
}

func consumeGoogleOAuthStateVerification(tx *gorm.DB, verification auth.Verification, identifier string) (bool, error) {
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

func googleOAuthStateIdentifier(state string) string {
	return googleOAuthStatePrefix + strings.TrimSpace(state)
}

func googleOAuthLinkStateIdentifier(state string) string {
	return googleOAuthLinkStatePrefix + strings.TrimSpace(state)
}

func googlePayloadFromClaims(payload *googleidtoken.Payload) (googleIDTokenPayload, error) {
	if payload == nil {
		return googleIDTokenPayload{}, errInvalidGoogleIdentity
	}
	profile := googleIDTokenPayload{
		Subject: strings.TrimSpace(payload.Subject),
		Email:   normalizeEmail(claimString(payload.Claims, "email")),
		Name:    strings.TrimSpace(claimString(payload.Claims, "name")),
		Picture: strings.TrimSpace(claimString(payload.Claims, "picture")),
	}
	verified, ok := claimBool(payload.Claims, "email_verified")
	if !ok || !verified {
		return googleIDTokenPayload{}, errInvalidGoogleIdentity
	}
	profile.EmailVerified = true
	return validateGoogleProfile(profile)
}

func validateGoogleProfile(profile googleIDTokenPayload) (googleIDTokenPayload, error) {
	profile.Subject = strings.TrimSpace(profile.Subject)
	profile.Email = normalizeEmail(profile.Email)
	profile.Name = strings.TrimSpace(profile.Name)
	profile.Picture = strings.TrimSpace(profile.Picture)
	if profile.Subject == "" || !validEmail(profile.Email) || !profile.EmailVerified {
		return googleIDTokenPayload{}, errInvalidGoogleIdentity
	}
	return profile, nil
}

func normalizeOAuthRedirectURI(raw string) (string, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return "", errors.New("redirectURI is required")
	}
	parsed, err := url.Parse(value)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return "", errors.New("redirectURI must be an absolute URL")
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", errors.New("redirectURI must use http or https")
	}
	parsed.Fragment = ""
	return parsed.String(), nil
}

func profileName(profile googleIDTokenPayload) string {
	name := strings.TrimSpace(profile.Name)
	if name != "" {
		return name
	}
	if local, _, ok := strings.Cut(profile.Email, "@"); ok && strings.TrimSpace(local) != "" {
		return local
	}
	return profile.Email
}

func tokenExtraString(token *oauth2.Token, key string) string {
	value, ok := token.Extra(key).(string)
	if !ok {
		return ""
	}
	return strings.TrimSpace(value)
}

func claimString(claims map[string]interface{}, key string) string {
	value, ok := claims[key].(string)
	if !ok {
		return ""
	}
	return value
}

func claimBool(claims map[string]interface{}, key string) (bool, bool) {
	switch value := claims[key].(type) {
	case bool:
		return value, true
	case string:
		switch strings.ToLower(strings.TrimSpace(value)) {
		case "true":
			return true, true
		case "false":
			return false, true
		}
	}
	return false, false
}

func optionalString(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}

func optionalTime(value time.Time) *time.Time {
	if value.IsZero() {
		return nil
	}
	value = value.UTC()
	return &value
}
