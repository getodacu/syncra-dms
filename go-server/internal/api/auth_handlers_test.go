package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"ai.ro/syncra/dms/internal/auth"
	"ai.ro/syncra/dms/internal/orgunits"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const testInternalToken = "internal-token"
const testAuthDeliveryToken = "delivery-token"
const testBetterAuthSecret = "J7mN2qR9vT4xY8bC6pL3wS5zK1hD0fG2"

func TestAuthRoutesRequireTrustedInternalToken(t *testing.T) {
	router, _ := newAuthTestRouter(t)

	response := authJSON(t, router, http.MethodPost, "/api/auth/sign-up/email", `{
		"name":"Ada",
		"email":"ada@example.com",
		"password":"password123"
	}`, nil)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusUnauthorized, response.Body.String())
	}
}

func TestSignUpReturnsVerificationCodeOnlyForTrustedDelivery(t *testing.T) {
	router, _ := newAuthTestRouter(t)
	internalHeaders := map[string]string{"X-Syncra-Internal-Token": testInternalToken}

	untrusted := authJSON(t, router, http.MethodPost, "/api/auth/sign-up/email", `{
		"name":"Ada",
		"email":"ada@example.com",
		"password":"password123"
	}`, internalHeaders)
	if untrusted.Code != http.StatusOK {
		t.Fatalf("untrusted status = %d body=%s", untrusted.Code, untrusted.Body.String())
	}
	var untrustedBody map[string]any
	decodeJSON(t, untrusted, &untrustedBody)
	if _, ok := untrustedBody["verificationCode"]; ok {
		t.Fatal("untrusted signup response leaked verificationCode")
	}

	trustedHeaders := map[string]string{
		"X-Syncra-Internal-Token":      testInternalToken,
		"X-Syncra-Auth-Delivery-Token": testAuthDeliveryToken,
	}
	trusted := authJSON(t, router, http.MethodPost, "/api/auth/sign-up/email", `{
		"name":"Grace",
		"email":"grace@example.com",
		"password":"password123"
	}`, trustedHeaders)
	if trusted.Code != http.StatusOK {
		t.Fatalf("trusted status = %d body=%s", trusted.Code, trusted.Body.String())
	}
	var trustedBody map[string]any
	decodeJSON(t, trusted, &trustedBody)
	if code, ok := trustedBody["verificationCode"].(string); !ok || len(code) != 6 {
		t.Fatalf("verificationCode = %#v, want 6-digit string", trustedBody["verificationCode"])
	}
}

func TestEmailSignupVerifyLoginSessionAndSignOut(t *testing.T) {
	router, db := newAuthTestRouter(t)
	headers := trustedAuthHeaders()

	signup := authJSON(t, router, http.MethodPost, "/api/auth/sign-up/email", `{
		"name":"Ada Lovelace",
		"email":"ada@example.com",
		"password":"password123"
	}`, headers)
	if signup.Code != http.StatusOK {
		t.Fatalf("signup status = %d body=%s", signup.Code, signup.Body.String())
	}
	var signupBody struct {
		VerificationCode string `json:"verificationCode"`
	}
	decodeJSON(t, signup, &signupBody)
	if signupBody.VerificationCode == "" {
		t.Fatal("signup did not return verificationCode for trusted delivery")
	}

	verify := authJSON(t, router, http.MethodPost, "/api/auth/email-otp/verify-email", `{
		"email":"ada@example.com",
		"otp":"`+signupBody.VerificationCode+`"
	}`, map[string]string{"X-Syncra-Internal-Token": testInternalToken})
	if verify.Code != http.StatusOK {
		t.Fatalf("verify status = %d body=%s", verify.Code, verify.Body.String())
	}
	var verifiedUser auth.User
	if err := db.First(&verifiedUser, "email = ?", "ada@example.com").Error; err != nil {
		t.Fatalf("load verified user: %v", err)
	}
	if verifiedUser.Status != "active" {
		t.Fatalf("verified user status = %q, want active", verifiedUser.Status)
	}

	login := authJSON(t, router, http.MethodPost, "/api/auth/sign-in/email", `{
		"email":"ada@example.com",
		"password":"password123",
		"rememberMe":true
	}`, map[string]string{"X-Syncra-Internal-Token": testInternalToken})
	if login.Code != http.StatusOK {
		t.Fatalf("login status = %d body=%s", login.Code, login.Body.String())
	}
	var loginBody struct {
		Session struct {
			Token string `json:"token"`
		} `json:"session"`
		User struct {
			Email         string `json:"email"`
			EmailVerified bool   `json:"emailVerified"`
		} `json:"user"`
	}
	decodeJSON(t, login, &loginBody)
	if loginBody.Session.Token == "" || loginBody.User.Email != "ada@example.com" || !loginBody.User.EmailVerified {
		t.Fatalf("unexpected login body: %#v", loginBody)
	}
	if !strings.Contains(login.Header().Get("Set-Cookie"), "auth.session_token=") {
		t.Fatalf("login Set-Cookie = %q, want auth.session_token", login.Header().Get("Set-Cookie"))
	}

	session := authJSON(t, router, http.MethodGet, "/api/auth/get-session", "", map[string]string{
		"X-Syncra-Internal-Token": testInternalToken,
		"Cookie":                  "auth.session_token=" + loginBody.Session.Token,
	})
	if session.Code != http.StatusOK {
		t.Fatalf("session status = %d body=%s", session.Code, session.Body.String())
	}
	var sessionBody struct {
		User struct {
			Email string `json:"email"`
		} `json:"user"`
	}
	decodeJSON(t, session, &sessionBody)
	if sessionBody.User.Email != "ada@example.com" {
		t.Fatalf("session user email = %q, want ada@example.com", sessionBody.User.Email)
	}

	signout := authJSON(t, router, http.MethodPost, "/api/auth/sign-out", `{}`, map[string]string{
		"X-Syncra-Internal-Token": testInternalToken,
		"Cookie":                  "auth.session_token=" + loginBody.Session.Token,
	})
	if signout.Code != http.StatusOK {
		t.Fatalf("signout status = %d body=%s", signout.Code, signout.Body.String())
	}
	if !strings.Contains(signout.Header().Get("Set-Cookie"), "Max-Age=0") {
		t.Fatalf("signout Set-Cookie = %q, want deletion cookie", signout.Header().Get("Set-Cookie"))
	}
}

func TestPasswordResetConsumesTokenAndRevokesSessions(t *testing.T) {
	router, db := newAuthTestRouter(t)
	createVerifiedUser(t, db, "ada@example.com", "old-password")
	loginUser(t, router, "ada@example.com", "old-password")

	resetRequest := authJSON(t, router, http.MethodPost, "/api/auth/password-reset/request", `{
		"email":"ada@example.com"
	}`, trustedAuthHeaders())
	if resetRequest.Code != http.StatusOK {
		t.Fatalf("reset request status = %d body=%s", resetRequest.Code, resetRequest.Body.String())
	}
	var resetRequestBody struct {
		ResetToken string `json:"resetToken"`
	}
	decodeJSON(t, resetRequest, &resetRequestBody)
	if resetRequestBody.ResetToken == "" {
		t.Fatal("trusted reset request did not return reset token")
	}

	reset := authJSON(t, router, http.MethodPost, "/api/auth/password-reset/confirm", `{
		"email":"ada@example.com",
		"token":"`+resetRequestBody.ResetToken+`",
		"password":"new-password"
	}`, map[string]string{"X-Syncra-Internal-Token": testInternalToken})
	if reset.Code != http.StatusOK {
		t.Fatalf("reset status = %d body=%s", reset.Code, reset.Body.String())
	}

	var sessionCount int64
	if err := db.Model(&auth.Session{}).Count(&sessionCount).Error; err != nil {
		t.Fatalf("count sessions: %v", err)
	}
	if sessionCount != 0 {
		t.Fatalf("session count after reset = %d, want 0", sessionCount)
	}

	oldLogin := authJSON(t, router, http.MethodPost, "/api/auth/sign-in/email", `{
		"email":"ada@example.com",
		"password":"old-password"
	}`, map[string]string{"X-Syncra-Internal-Token": testInternalToken})
	if oldLogin.Code != http.StatusUnauthorized {
		t.Fatalf("old password status = %d, want %d", oldLogin.Code, http.StatusUnauthorized)
	}

	newLogin := authJSON(t, router, http.MethodPost, "/api/auth/sign-in/email", `{
		"email":"ada@example.com",
		"password":"new-password"
	}`, map[string]string{"X-Syncra-Internal-Token": testInternalToken})
	if newLogin.Code != http.StatusOK {
		t.Fatalf("new password status = %d body=%s", newLogin.Code, newLogin.Body.String())
	}
}

func TestPasswordResetTokenOnlyReturnedForTrustedDelivery(t *testing.T) {
	router, db := newAuthTestRouter(t)
	createVerifiedUser(t, db, "ada@example.com", "old-password")

	response := authJSON(t, router, http.MethodPost, "/api/auth/password-reset/request", `{
		"email":"ada@example.com"
	}`, map[string]string{"X-Syncra-Internal-Token": testInternalToken})
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", response.Code, response.Body.String())
	}
	var body map[string]any
	decodeJSON(t, response, &body)
	if _, ok := body["resetToken"]; ok {
		t.Fatal("untrusted reset response leaked resetToken")
	}
}

func TestOAuthStartRequiresConfiguredProvider(t *testing.T) {
	router, _ := newAuthTestRouter(t)

	missingConfig := authJSON(t, router, http.MethodPost, "/api/auth/oauth/google/start", `{
		"redirectURI":"http://localhost:5173/api/auth/google/callback"
	}`, map[string]string{"X-Syncra-Internal-Token": testInternalToken})
	if missingConfig.Code != http.StatusServiceUnavailable {
		t.Fatalf("missing config status = %d, want %d", missingConfig.Code, http.StatusServiceUnavailable)
	}

	configured, _ := newAuthTestRouterWithOptions(t, RouterOptions{
		GoogleClientID:     "google-client",
		GoogleClientSecret: "google-secret",
	})
	response := authJSON(t, configured, http.MethodPost, "/api/auth/oauth/google/start", `{
		"redirectURI":"http://localhost:5173/api/auth/google/callback"
	}`, map[string]string{"X-Syncra-Internal-Token": testInternalToken})
	if response.Code != http.StatusOK {
		t.Fatalf("configured status = %d body=%s", response.Code, response.Body.String())
	}
	var body struct {
		AuthorizationURL string `json:"authorizationUrl"`
		State            string `json:"state"`
		StateExpiresAt   string `json:"stateExpiresAt"`
	}
	decodeJSON(t, response, &body)
	if !strings.Contains(body.AuthorizationURL, "client_id=google-client") || body.State == "" || body.StateExpiresAt == "" {
		t.Fatalf("unexpected oauth start body: %#v", body)
	}
}

func TestOAuthCallbackReturnsUnauthorizedOnProviderFailure(t *testing.T) {
	router, _ := newAuthTestRouterWithOptions(t, RouterOptions{
		GoogleClientID:     "google-client",
		GoogleClientSecret: "google-secret",
		OAuthProfileFetcher: func(context.Context, string, string, string) (OAuthProfile, error) {
			return OAuthProfile{}, context.Canceled
		},
	})

	response := authJSON(t, router, http.MethodPost, "/api/auth/oauth/google/callback", `{
		"code":"oauth-code",
		"state":"state",
		"redirectURI":"http://localhost:5173/api/auth/google/callback"
	}`, map[string]string{"X-Syncra-Internal-Token": testInternalToken})
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("oauth callback status = %d, want %d; body=%s", response.Code, http.StatusUnauthorized, response.Body.String())
	}
}

func TestOAuthCallbackCreatesSessionAndProviderAccount(t *testing.T) {
	router, db := newAuthTestRouterWithOptions(t, RouterOptions{
		GoogleClientID:     "google-client",
		GoogleClientSecret: "google-secret",
		OAuthProfileFetcher: func(_ context.Context, providerID string, code string, redirectURI string) (OAuthProfile, error) {
			if providerID != auth.GoogleProviderID || code != "oauth-code" || redirectURI == "" {
				t.Fatalf("unexpected oauth fetch input: %s %s %s", providerID, code, redirectURI)
			}
			return OAuthProfile{
				ProviderID: auth.GoogleProviderID,
				AccountID:  "google-account-1",
				Email:      "ada@example.com",
				Name:       "Ada Lovelace",
				Verified:   true,
			}, nil
		},
	})

	response := authJSON(t, router, http.MethodPost, "/api/auth/oauth/google/callback", `{
		"code":"oauth-code",
		"state":"state",
		"redirectURI":"http://localhost:5173/api/auth/google/callback"
	}`, map[string]string{"X-Syncra-Internal-Token": testInternalToken})
	if response.Code != http.StatusOK {
		t.Fatalf("oauth callback status = %d body=%s", response.Code, response.Body.String())
	}
	var body struct {
		Session struct {
			Token string `json:"token"`
		} `json:"session"`
		User struct {
			Email         string `json:"email"`
			EmailVerified bool   `json:"emailVerified"`
		} `json:"user"`
	}
	decodeJSON(t, response, &body)
	if body.Session.Token == "" || body.User.Email != "ada@example.com" || !body.User.EmailVerified {
		t.Fatalf("unexpected oauth callback body: %#v", body)
	}

	var account auth.AuthAccount
	if err := db.First(&account, "provider_id = ? AND account_id = ?", auth.GoogleProviderID, "google-account-1").Error; err != nil {
		t.Fatalf("load oauth account: %v", err)
	}
	if account.UserID == "" {
		t.Fatal("oauth account user id was empty")
	}
}

func TestOAuthCallbackPromotesExistingInvitedVerifiedEmailUser(t *testing.T) {
	router, db := newAuthTestRouterWithOptions(t, RouterOptions{
		GoogleClientID:     "google-client",
		GoogleClientSecret: "google-secret",
		OAuthProfileFetcher: func(_ context.Context, providerID string, code string, redirectURI string) (OAuthProfile, error) {
			return OAuthProfile{
				ProviderID: providerID,
				AccountID:  "google-account-2",
				Email:      "ada@example.com",
				Name:       "Ada Lovelace",
				Verified:   true,
			}, nil
		},
	})
	now := time.Now().UTC()
	user := auth.User{
		Name:          "Ada Lovelace",
		Email:         "ada@example.com",
		EmailVerified: false,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create invited user: %v", err)
	}
	if user.Status != "invited" {
		t.Fatalf("created user status = %q, want invited", user.Status)
	}

	response := authJSON(t, router, http.MethodPost, "/api/auth/oauth/google/callback", `{
		"code":"oauth-code",
		"state":"state",
		"redirectURI":"http://localhost:5173/api/auth/google/callback"
	}`, map[string]string{"X-Syncra-Internal-Token": testInternalToken})
	if response.Code != http.StatusOK {
		t.Fatalf("oauth callback status = %d body=%s", response.Code, response.Body.String())
	}

	var promoted auth.User
	if err := db.First(&promoted, "id = ?", user.ID).Error; err != nil {
		t.Fatalf("load promoted user: %v", err)
	}
	if !promoted.EmailVerified {
		t.Fatal("oauth callback did not mark user email verified")
	}
	if promoted.Status != "active" {
		t.Fatalf("oauth promoted user status = %q, want active", promoted.Status)
	}
}

func newAuthTestRouter(t *testing.T) (http.Handler, *gorm.DB) {
	t.Helper()
	return newAuthTestRouterWithOptions(t, RouterOptions{})
}

func newAuthTestRouterWithOptions(t *testing.T, options RouterOptions) (http.Handler, *gorm.DB) {
	t.Helper()
	name := strings.NewReplacer("/", "_", " ", "_").Replace(t.Name())
	db, err := gorm.Open(sqlite.Open("file:"+name+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&auth.User{}, &auth.AuthAccount{}, &auth.Session{}, &auth.Verification{}, &orgunits.Unit{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	base := RouterOptions{
		DB:                  db,
		BetterAuthSecret:    testBetterAuthSecret,
		AuthDeliveryToken:   testAuthDeliveryToken,
		InternalAPIToken:    testInternalToken,
		AuthSessionTTL:      7 * 24 * time.Hour,
		AuthVerificationTTL: 5 * time.Minute,
		AuthCookieSecure:    false,
		GoogleClientID:      options.GoogleClientID,
		GoogleClientSecret:  options.GoogleClientSecret,
		GitHubClientID:      options.GitHubClientID,
		GitHubClientSecret:  options.GitHubClientSecret,
		OAuthProfileFetcher: options.OAuthProfileFetcher,
	}
	return NewRouter(base), db
}

func trustedAuthHeaders() map[string]string {
	return map[string]string{
		"X-Syncra-Internal-Token":      testInternalToken,
		"X-Syncra-Auth-Delivery-Token": testAuthDeliveryToken,
	}
}

func authJSON(t *testing.T, router http.Handler, method string, path string, body string, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	if body != "" {
		request.Header.Set("Content-Type", "application/json")
	}
	for key, value := range headers {
		request.Header.Set(key, value)
	}
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

func decodeJSON(t *testing.T, response *httptest.ResponseRecorder, out any) {
	t.Helper()
	if err := json.Unmarshal(response.Body.Bytes(), out); err != nil {
		t.Fatalf("decode response: %v body=%s", err, response.Body.String())
	}
}

func createVerifiedUser(t *testing.T, db *gorm.DB, email string, password string) auth.User {
	t.Helper()
	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	now := time.Now().UTC()
	user := auth.User{
		Name:          "Ada Lovelace",
		Email:         email,
		EmailVerified: true,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	account := auth.AuthAccount{
		AccountID:  user.ID,
		ProviderID: auth.CredentialProviderID,
		UserID:     user.ID,
		Password:   hash,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := db.Create(&account).Error; err != nil {
		t.Fatalf("create account: %v", err)
	}
	return user
}

func loginUser(t *testing.T, router http.Handler, email string, password string) string {
	t.Helper()
	login := authJSON(t, router, http.MethodPost, "/api/auth/sign-in/email", `{
		"email":"`+email+`",
		"password":"`+password+`"
	}`, map[string]string{"X-Syncra-Internal-Token": testInternalToken})
	if login.Code != http.StatusOK {
		t.Fatalf("login status = %d body=%s", login.Code, login.Body.String())
	}
	var body struct {
		Session struct {
			Token string `json:"token"`
		} `json:"session"`
	}
	decodeJSON(t, login, &body)
	return body.Session.Token
}
