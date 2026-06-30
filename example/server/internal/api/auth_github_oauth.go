package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	auth "ai.ro/syncra/internal/auth"
	"ai.ro/syncra/internal/billing"
)

const githubOAuthStatePrefix = "oauth-state:github:"
const githubOAuthLinkStatePrefix = "oauth-link-state:github:"

var (
	errGitHubOAuthNotConfigured = errors.New("github oauth is not configured")
	errInvalidGitHubOAuthState  = errors.New("invalid github oauth state")
	errInvalidGitHubIdentity    = errors.New("invalid github identity")
)

// swagger:model githubOAuthStartRequest
type githubOAuthStartRequest struct {
	RedirectURI string `json:"redirectURI"`
}

// swagger:model githubOAuthCallbackRequest
type githubOAuthCallbackRequest struct {
	Code        string `json:"code"`
	State       string `json:"state"`
	RedirectURI string `json:"redirectURI"`
}

// swagger:model githubOAuthStartResponse
type githubOAuthStartResponse struct {
	AuthorizationURL string `json:"authorizationUrl"`
	State            string `json:"state"`
	StateExpiresAt   string `json:"stateExpiresAt"`
}

type githubOAuthToken struct {
	AccessToken           *string
	RefreshToken          *string
	AccessTokenExpiresAt  *time.Time
	RefreshTokenExpiresAt *time.Time
	Scope                 *string
}

type githubOAuthProfile struct {
	Subject       string
	Email         string
	EmailVerified bool
	Name          string
	Login         string
	Picture       string
}

type githubUserResponse struct {
	ID        int64   `json:"id"`
	Login     string  `json:"login"`
	Name      string  `json:"name"`
	AvatarURL string  `json:"avatar_url"`
	Email     *string `json:"email"`
}

type githubEmailResponse struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

func (h *Handler) StartGitHubOAuth(c *gin.Context) {
	var req githubOAuthStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	redirectURI, err := normalizeOAuthRedirectURI(req.RedirectURI)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	if !h.githubOAuthConfigured() {
		writeError(c, http.StatusServiceUnavailable, errGitHubOAuthNotConfigured.Error())
		return
	}

	var state string
	var expiresAt time.Time
	if err := h.DB.Transaction(func(tx *gorm.DB) error {
		return rotateGitHubOAuthState(tx, h.BetterAuthSecret, redirectURI, h.authVerificationTTL(), &state, &expiresAt)
	}); err != nil {
		writeError(c, http.StatusInternalServerError, "failed to create github oauth state")
		return
	}

	c.JSON(http.StatusOK, githubOAuthStartResponse{
		AuthorizationURL: h.githubOAuthConfig(redirectURI).AuthCodeURL(state),
		State:            state,
		StateExpiresAt:   expiresAt.UTC().Format(time.RFC3339Nano),
	})
}

func (h *Handler) SignInGitHubOAuth(c *gin.Context) {
	var req githubOAuthCallbackRequest
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
	if !h.githubOAuthConfigured() {
		writeError(c, http.StatusServiceUnavailable, errGitHubOAuthNotConfigured.Error())
		return
	}

	if err := h.consumeGitHubOAuthState(state, redirectURI); err != nil {
		if errors.Is(err, errInvalidGitHubOAuthState) {
			writeError(c, http.StatusBadRequest, errInvalidGitHubOAuthState.Error())
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to verify github oauth state")
		return
	}

	token, err := h.exchangeGitHubOAuthCode(c.Request.Context(), code, redirectURI)
	if err != nil {
		writeError(c, http.StatusBadGateway, "failed to exchange github authorization code")
		return
	}
	if token.AccessToken == nil {
		writeError(c, http.StatusBadGateway, "failed to exchange github authorization code")
		return
	}
	profile, err := h.fetchGitHubOAuthProfile(c.Request.Context(), *token.AccessToken)
	if err != nil {
		writeError(c, http.StatusUnauthorized, errInvalidGitHubIdentity.Error())
		return
	}

	session, user, err := h.signInGitHubProfile(c, profile, token)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to sign in with github")
		return
	}

	h.setSessionCookie(c, session.Token, session.ExpiresAt, true)
	c.JSON(http.StatusOK, authSessionPayload{
		Session: authSessionJSON(session),
		User:    authUserJSON(user),
	})
}

func (h *Handler) StartGitHubAccountLink(c *gin.Context) {
	if _, ok, err := h.loadAuthenticatedSession(c); err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	} else if !ok {
		writeError(c, http.StatusUnauthorized, "authentication required")
		return
	}

	var req githubOAuthStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	redirectURI, err := normalizeOAuthRedirectURI(req.RedirectURI)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	if !h.githubOAuthConfigured() {
		writeError(c, http.StatusServiceUnavailable, errGitHubOAuthNotConfigured.Error())
		return
	}

	var state string
	var expiresAt time.Time
	if err := h.DB.Transaction(func(tx *gorm.DB) error {
		return rotateGitHubOAuthLinkState(tx, h.BetterAuthSecret, redirectURI, h.authVerificationTTL(), &state, &expiresAt)
	}); err != nil {
		writeError(c, http.StatusInternalServerError, "failed to create github oauth state")
		return
	}

	c.JSON(http.StatusOK, githubOAuthStartResponse{
		AuthorizationURL: h.githubOAuthConfig(redirectURI).AuthCodeURL(state),
		State:            state,
		StateExpiresAt:   expiresAt.UTC().Format(time.RFC3339Nano),
	})
}

func (h *Handler) LinkGitHubAccount(c *gin.Context) {
	session, ok, err := h.loadAuthenticatedSession(c)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}
	if !ok {
		writeError(c, http.StatusUnauthorized, "authentication required")
		return
	}

	var req githubOAuthCallbackRequest
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
	if !h.githubOAuthConfigured() {
		writeError(c, http.StatusServiceUnavailable, errGitHubOAuthNotConfigured.Error())
		return
	}

	if err := h.consumeGitHubOAuthLinkState(state, redirectURI); err != nil {
		if errors.Is(err, errInvalidGitHubOAuthState) {
			writeError(c, http.StatusBadRequest, errInvalidGitHubOAuthState.Error())
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to verify github oauth state")
		return
	}

	token, err := h.exchangeGitHubOAuthCode(c.Request.Context(), code, redirectURI)
	if err != nil {
		writeError(c, http.StatusBadGateway, "failed to exchange github authorization code")
		return
	}
	if token.AccessToken == nil {
		writeError(c, http.StatusBadGateway, "failed to exchange github authorization code")
		return
	}
	profile, err := h.fetchGitHubOAuthProfile(c.Request.Context(), *token.AccessToken)
	if err != nil {
		writeError(c, http.StatusUnauthorized, errInvalidGitHubIdentity.Error())
		return
	}

	account, err := h.linkGitHubProfile(c, session.UserID, profile, token)
	if err != nil {
		if errors.Is(err, errOAuthAccountLinked) || errors.Is(err, errProviderAlreadyLinked) {
			writeError(c, http.StatusConflict, err.Error())
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to link github account")
		return
	}

	c.JSON(http.StatusOK, authAccountListItemJSON(account))
}

func (h *Handler) githubOAuthConfigured() bool {
	return strings.TrimSpace(h.GitHubClientID) != "" && strings.TrimSpace(h.GitHubClientSecret) != ""
}

func (h *Handler) githubOAuthConfig(redirectURI string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     strings.TrimSpace(h.GitHubClientID),
		ClientSecret: strings.TrimSpace(h.GitHubClientSecret),
		RedirectURL:  redirectURI,
		Scopes:       []string{"read:user", "user:email"},
		Endpoint:     githuboauth.Endpoint,
	}
}

func (h *Handler) exchangeGitHubOAuthCode(ctx context.Context, code string, redirectURI string) (githubOAuthToken, error) {
	if h.GitHubOAuthExchange != nil {
		return h.GitHubOAuthExchange(ctx, code, redirectURI)
	}

	oauthToken, err := h.githubOAuthConfig(redirectURI).Exchange(ctx, code)
	if err != nil {
		return githubOAuthToken{}, err
	}
	accessToken := optionalString(oauthToken.AccessToken)
	if accessToken == nil {
		return githubOAuthToken{}, errors.New("github token response missing access_token")
	}

	return githubOAuthToken{
		AccessToken:           accessToken,
		RefreshToken:          optionalString(oauthToken.RefreshToken),
		AccessTokenExpiresAt:  optionalTime(oauthToken.Expiry),
		RefreshTokenExpiresAt: optionalTime(time.Time{}),
		Scope:                 optionalString(tokenExtraString(oauthToken, "scope")),
	}, nil
}

func (h *Handler) fetchGitHubOAuthProfile(ctx context.Context, accessToken string) (githubOAuthProfile, error) {
	if h.GitHubProfileFetch != nil {
		profile, err := h.GitHubProfileFetch(ctx, accessToken)
		if err != nil {
			return githubOAuthProfile{}, err
		}
		return validateGitHubProfile(profile)
	}

	accessToken = strings.TrimSpace(accessToken)
	if accessToken == "" {
		return githubOAuthProfile{}, errInvalidGitHubIdentity
	}
	client := &http.Client{Timeout: 10 * time.Second}
	var user githubUserResponse
	if err := githubAPIGetJSON(ctx, client, accessToken, "https://api.github.com/user", &user); err != nil {
		return githubOAuthProfile{}, err
	}
	var emails []githubEmailResponse
	if err := githubAPIGetJSON(ctx, client, accessToken, "https://api.github.com/user/emails", &emails); err != nil {
		return githubOAuthProfile{}, err
	}

	var primaryVerifiedEmail string
	for _, email := range emails {
		if email.Primary && email.Verified {
			primaryVerifiedEmail = email.Email
			break
		}
	}
	if primaryVerifiedEmail == "" {
		return githubOAuthProfile{}, errInvalidGitHubIdentity
	}

	return validateGitHubProfile(githubOAuthProfile{
		Subject:       strconv.FormatInt(user.ID, 10),
		Email:         primaryVerifiedEmail,
		EmailVerified: true,
		Name:          user.Name,
		Login:         user.Login,
		Picture:       user.AvatarURL,
	})
}

func githubAPIGetJSON(ctx context.Context, client *http.Client, accessToken string, url string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("github api %s returned %s", url, resp.Status)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func (h *Handler) consumeGitHubOAuthState(state string, redirectURI string) error {
	identifier := githubOAuthStateIdentifier(state)
	var verification auth.Verification
	if err := h.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("identifier = ?", identifier).Limit(1).Find(&verification)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errInvalidGitHubOAuthState
		}
		if !verification.ExpiresAt.After(h.now()) {
			if _, err := consumeGitHubOAuthStateVerification(tx, verification, identifier); err != nil {
				return err
			}
			return errInvalidGitHubOAuthState
		}
		if !auth.VerifyCode(h.BetterAuthSecret, identifier, redirectURI, verification.Value) {
			return errInvalidGitHubOAuthState
		}
		consumed, err := consumeGitHubOAuthStateVerification(tx, verification, identifier)
		if err != nil {
			return err
		}
		if !consumed {
			return errInvalidGitHubOAuthState
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (h *Handler) consumeGitHubOAuthLinkState(state string, redirectURI string) error {
	identifier := githubOAuthLinkStateIdentifier(state)
	var verification auth.Verification
	if err := h.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("identifier = ?", identifier).Limit(1).Find(&verification)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errInvalidGitHubOAuthState
		}
		if !verification.ExpiresAt.After(h.now()) {
			if _, err := consumeGitHubOAuthStateVerification(tx, verification, identifier); err != nil {
				return err
			}
			return errInvalidGitHubOAuthState
		}
		if !auth.VerifyCode(h.BetterAuthSecret, identifier, redirectURI, verification.Value) {
			return errInvalidGitHubOAuthState
		}
		consumed, err := consumeGitHubOAuthStateVerification(tx, verification, identifier)
		if err != nil {
			return err
		}
		if !consumed {
			return errInvalidGitHubOAuthState
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (h *Handler) signInGitHubProfile(c *gin.Context, profile githubOAuthProfile, token githubOAuthToken) (auth.Session, auth.User, error) {
	now := h.now()
	var session auth.Session
	var user auth.User

	err := h.DB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		account, found, err := githubAccountForSubject(tx, profile.Subject)
		if err != nil {
			return err
		}
		newlyVerified := false
		if found {
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&user, "id = ?", account.UserID).Error; err != nil {
				return err
			}
			if err := updateGitHubAccountTokens(tx, account, token, now); err != nil {
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
					Name:          githubProfileName(profile),
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
				ProviderID:            auth.GitHubProviderID,
				UserID:                user.ID,
				AccessToken:           token.AccessToken,
				RefreshToken:          token.RefreshToken,
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

func (h *Handler) linkGitHubProfile(c *gin.Context, userID string, profile githubOAuthProfile, token githubOAuthToken) (auth.AuthAccount, error) {
	now := h.now()
	var linked auth.AuthAccount

	err := h.DB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		account, found, err := githubAccountForSubject(tx, profile.Subject)
		if err != nil {
			return err
		}
		if found {
			if account.UserID != userID {
				return errOAuthAccountLinked
			}
			if err := updateGitHubAccountTokens(tx, account, token, now); err != nil {
				return err
			}
			return tx.First(&linked, "id = ?", account.ID).Error
		}

		var existing auth.AuthAccount
		result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ? AND provider_id = ?", userID, auth.GitHubProviderID).
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
			ProviderID:            auth.GitHubProviderID,
			UserID:                userID,
			AccessToken:           token.AccessToken,
			RefreshToken:          token.RefreshToken,
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

func githubAccountForSubject(tx *gorm.DB, subject string) (auth.AuthAccount, bool, error) {
	var account auth.AuthAccount
	result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("provider_id = ? AND account_id = ?", auth.GitHubProviderID, subject).
		Limit(1).
		Find(&account)
	if result.Error != nil {
		return auth.AuthAccount{}, false, result.Error
	}
	return account, result.RowsAffected > 0, nil
}

func updateGitHubAccountTokens(tx *gorm.DB, account auth.AuthAccount, token githubOAuthToken, now time.Time) error {
	updates := map[string]any{
		"access_token":            token.AccessToken,
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

func rotateGitHubOAuthState(tx *gorm.DB, secret string, redirectURI string, ttl time.Duration, stateOut *string, expiresOut *time.Time) error {
	state, err := auth.GenerateSessionToken()
	if err != nil {
		return err
	}
	identifier := githubOAuthStateIdentifier(state)
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

func rotateGitHubOAuthLinkState(tx *gorm.DB, secret string, redirectURI string, ttl time.Duration, stateOut *string, expiresOut *time.Time) error {
	state, err := auth.GenerateSessionToken()
	if err != nil {
		return err
	}
	identifier := githubOAuthLinkStateIdentifier(state)
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

func consumeGitHubOAuthStateVerification(tx *gorm.DB, verification auth.Verification, identifier string) (bool, error) {
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

func githubOAuthStateIdentifier(state string) string {
	return githubOAuthStatePrefix + strings.TrimSpace(state)
}

func githubOAuthLinkStateIdentifier(state string) string {
	return githubOAuthLinkStatePrefix + strings.TrimSpace(state)
}

func validateGitHubProfile(profile githubOAuthProfile) (githubOAuthProfile, error) {
	profile.Subject = strings.TrimSpace(profile.Subject)
	profile.Email = normalizeEmail(profile.Email)
	profile.Name = strings.TrimSpace(profile.Name)
	profile.Login = strings.TrimSpace(profile.Login)
	profile.Picture = strings.TrimSpace(profile.Picture)
	if profile.Subject == "" || !validEmail(profile.Email) || !profile.EmailVerified {
		return githubOAuthProfile{}, errInvalidGitHubIdentity
	}
	return profile, nil
}

func githubProfileName(profile githubOAuthProfile) string {
	name := strings.TrimSpace(profile.Name)
	if name != "" {
		return name
	}
	login := strings.TrimSpace(profile.Login)
	if login != "" {
		return login
	}
	if local, _, ok := strings.Cut(profile.Email, "@"); ok && strings.TrimSpace(local) != "" {
		return local
	}
	return profile.Email
}
