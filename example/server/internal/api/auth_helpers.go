package api

import (
	"net/http"
	"strings"
	"time"

	auth "ai.ro/syncra/internal/auth"
	"github.com/gin-gonic/gin"
)

const authSessionCookieName = "auth.session_token"
const authDeliveryHeader = "X-Syncra-Auth-Delivery-Token"
const defaultOnboardingCredits = 100

// authUserResponse is a public auth user payload.
//
// swagger:model authUserResponse
type authUserResponse struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"emailVerified"`
	// Optional profile image URL.
	//
	// Extensions:
	// ---
	// x-nullable: true
	// ---
	Image *string `json:"image"`
	// Preferred UI language. Supported values: en, ro.
	PreferredLanguage string `json:"preferredLanguage"`
	Role              string `json:"role"`
	// Last successful email/password login timestamp.
	//
	// Extensions:
	// ---
	// x-nullable: true
	// ---
	LastLoginAt *string `json:"lastLoginAt"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

// authSessionResponse is a public auth session payload.
//
// swagger:model authSessionResponse
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

// authSessionListItemResponse is a safe session payload without the session token.
//
// swagger:model authSessionListItemResponse
type authSessionListItemResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"userId"`
	ExpiresAt string `json:"expiresAt"`
	IPAddress string `json:"ipAddress,omitempty"`
	UserAgent string `json:"userAgent,omitempty"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Current   bool   `json:"current"`
}

// authSessionListResponse lists active sessions for the current user.
//
// swagger:model authSessionListResponse
type authSessionListResponse struct {
	Sessions []authSessionListItemResponse `json:"sessions"`
}

// deleteAuthSessionResponse confirms session revocation.
//
// swagger:model deleteAuthSessionResponse
type deleteAuthSessionResponse struct {
	DeletedID    string `json:"deleted_id"`
	DeletedCount int    `json:"deleted_count"`
}

// authAccountListItemResponse is a safe linked account payload.
//
// swagger:model authAccountListItemResponse
type authAccountListItemResponse struct {
	ID         string `json:"id"`
	ProviderID string `json:"providerId"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
}

// authAccountListResponse lists linked sign-in methods for the current user.
//
// swagger:model authAccountListResponse
type authAccountListResponse struct {
	Accounts []authAccountListItemResponse `json:"accounts"`
}

// deleteAuthAccountResponse confirms linked account removal.
//
// swagger:model deleteAuthAccountResponse
type deleteAuthAccountResponse struct {
	DeletedProviderID string `json:"deleted_provider_id"`
	DeletedCount      int    `json:"deleted_count"`
}

// signInEmailResponse is returned after email sign-in.
//
// swagger:model signInEmailResponse
type signInEmailResponse struct {
	Session authSessionResponse `json:"session"`
	User    authUserResponse    `json:"user"`
}

// authImpersonationResponse describes an active admin impersonation.
//
// swagger:model authImpersonationResponse
type authImpersonationResponse struct {
	AdminUser  authUserResponse `json:"adminUser"`
	TargetUser authUserResponse `json:"targetUser"`
	StartedAt  string           `json:"startedAt"`
}

// getSessionResponse is returned for an active email session.
//
// swagger:model getSessionResponse
type getSessionResponse struct {
	Session       authSessionResponse        `json:"session"`
	User          authUserResponse           `json:"user"`
	Impersonation *authImpersonationResponse `json:"impersonation"`
}

type authSessionPayload = signInEmailResponse

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func verificationIdentifier(email string) string {
	return "email-verification:" + normalizeEmail(email)
}

func passwordResetIdentifier(email string) string {
	return "password-reset:" + normalizeEmail(email)
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

func authSessionListItemJSON(session auth.Session, currentSessionID string) authSessionListItemResponse {
	return authSessionListItemResponse{
		ID:        session.ID,
		UserID:    session.UserID,
		ExpiresAt: session.ExpiresAt.UTC().Format(time.RFC3339Nano),
		IPAddress: session.IPAddress,
		UserAgent: session.UserAgent,
		CreatedAt: session.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt: session.UpdatedAt.UTC().Format(time.RFC3339Nano),
		Current:   session.ID == currentSessionID,
	}
}

func authAccountListItemJSON(account auth.AuthAccount) authAccountListItemResponse {
	return authAccountListItemResponse{
		ID:         account.ID,
		ProviderID: account.ProviderID,
		CreatedAt:  account.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:  account.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
}

func effectiveAuthUser(session auth.Session) auth.User {
	if session.ImpersonatedUserID != nil && session.ImpersonatedUser != nil && session.ImpersonatedUser.ID != "" {
		return *session.ImpersonatedUser
	}
	return session.User
}

func authImpersonationJSON(session auth.Session) *authImpersonationResponse {
	if session.ImpersonatedUserID == nil || session.ImpersonationStartedAt == nil || session.ImpersonatedUser == nil || session.ImpersonatedUser.ID == "" {
		return nil
	}
	return &authImpersonationResponse{
		AdminUser:  authUserJSON(session.User),
		TargetUser: authUserJSON(*session.ImpersonatedUser),
		StartedAt:  session.ImpersonationStartedAt.UTC().Format(time.RFC3339Nano),
	}
}

func authSessionPayloadJSON(session auth.Session) getSessionResponse {
	return getSessionResponse{
		Session:       authSessionJSON(session),
		User:          authUserJSON(effectiveAuthUser(session)),
		Impersonation: authImpersonationJSON(session),
	}
}

func (h *Handler) authSessionTTL() time.Duration {
	if h.AuthSessionTTL > 0 {
		return h.AuthSessionTTL
	}
	return 7 * 24 * time.Hour
}

func (h *Handler) authVerificationTTL() time.Duration {
	if h.AuthVerificationTTL > 0 {
		return h.AuthVerificationTTL
	}
	return 5 * time.Minute
}

func (h *Handler) onboardingCredits() int {
	if h.OnboardingCredits > 0 {
		return h.OnboardingCredits
	}
	return defaultOnboardingCredits
}

func (h *Handler) trustedAuthDeliveryRequest(c *gin.Context) bool {
	return h.AuthDeliveryToken != "" && c.GetHeader(authDeliveryHeader) == h.AuthDeliveryToken
}

func (h *Handler) setSessionCookie(c *gin.Context, token string, expiresAt time.Time, remember bool) {
	cookie := &http.Cookie{
		Name:     authSessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   h.AuthCookieSecure,
		SameSite: http.SameSiteLaxMode,
	}
	if remember {
		cookie.MaxAge = int(time.Until(expiresAt).Seconds())
		if cookie.MaxAge < 0 {
			cookie.MaxAge = 0
		}
		cookie.Expires = expiresAt
	}
	http.SetCookie(c.Writer, cookie)
}

func (h *Handler) clearSessionCookie(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     authSessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   h.AuthCookieSecure,
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
