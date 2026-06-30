package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
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

func TestSendVerificationOTPDoesNotRotateForNonActiveUsers(t *testing.T) {
	tests := []struct {
		name       string
		status     string
		softDelete bool
	}{
		{name: "inactive", status: "inactive"},
		{name: "suspended", status: "suspended"},
		{name: "deleted status", status: "deleted"},
		{name: "soft deleted active", status: "active", softDelete: true},
	}

	for _, tc := range tests {
		t.Run(tc.name+" does not create", func(t *testing.T) {
			router, db := newAuthTestRouter(t)
			createStatusUser(t, db, "ada@example.com", tc.status, false, tc.softDelete)

			response := authJSON(t, router, http.MethodPost, "/api/auth/email-otp/send-verification-otp", `{
				"email":"ada@example.com",
				"type":"email-verification"
			}`, trustedAuthHeaders())
			if response.Code != http.StatusOK {
				t.Fatalf("send verification status = %d body=%s, want ok", response.Code, response.Body.String())
			}
			var body map[string]any
			decodeJSON(t, response, &body)
			if _, ok := body["verificationCode"]; ok {
				t.Fatalf("response leaked verificationCode for non-active user: %s", response.Body.String())
			}
			assertVerificationCount(t, db, verificationIdentifier("ada@example.com"), 0)
		})

		t.Run(tc.name+" does not update", func(t *testing.T) {
			router, db := newAuthTestRouter(t)
			createStatusUser(t, db, "ada@example.com", tc.status, false, tc.softDelete)
			createEmailVerification(t, db, "ada@example.com")
			before := captureVerificationState(t, db, verificationIdentifier("ada@example.com"))

			response := authJSON(t, router, http.MethodPost, "/api/auth/email-otp/send-verification-otp", `{
				"email":"ada@example.com",
				"type":"email-verification"
			}`, trustedAuthHeaders())
			if response.Code != http.StatusOK {
				t.Fatalf("send verification status = %d body=%s, want ok", response.Code, response.Body.String())
			}
			var body map[string]any
			decodeJSON(t, response, &body)
			if _, ok := body["verificationCode"]; ok {
				t.Fatalf("response leaked verificationCode for non-active user: %s", response.Body.String())
			}
			assertVerificationCount(t, db, verificationIdentifier("ada@example.com"), 1)
			after := captureVerificationState(t, db, verificationIdentifier("ada@example.com"))
			assertVerificationStateUnchanged(t, before, after)
		})
	}
}

func TestSendVerificationOTPRotatesForInvitedUser(t *testing.T) {
	router, db := newAuthTestRouter(t)
	createInvitedUser(t, db, "ada@example.com")

	response := authJSON(t, router, http.MethodPost, "/api/auth/email-otp/send-verification-otp", `{
		"email":"ada@example.com",
		"type":"email-verification"
	}`, trustedAuthHeaders())
	if response.Code != http.StatusOK {
		t.Fatalf("send verification status = %d body=%s, want ok", response.Code, response.Body.String())
	}
	var body struct {
		VerificationCode      string `json:"verificationCode"`
		VerificationExpiresAt string `json:"verificationExpiresAt"`
	}
	decodeJSON(t, response, &body)
	if len(body.VerificationCode) != 6 || body.VerificationExpiresAt == "" {
		t.Fatalf("unexpected verification response: %#v", body)
	}
	assertVerificationCount(t, db, verificationIdentifier("ada@example.com"), 1)
}

func TestDuplicateSignUpDoesNotRotateVerificationForNonActiveUsers(t *testing.T) {
	tests := []struct {
		name       string
		status     string
		softDelete bool
	}{
		{name: "inactive", status: "inactive"},
		{name: "suspended", status: "suspended"},
		{name: "deleted status", status: "deleted"},
		{name: "soft deleted active", status: "active", softDelete: true},
	}

	for _, tc := range tests {
		t.Run(tc.name+" does not create", func(t *testing.T) {
			router, db := newAuthTestRouter(t)
			createStatusUser(t, db, "ada@example.com", tc.status, false, tc.softDelete)

			response := authJSON(t, router, http.MethodPost, "/api/auth/sign-up/email", `{
				"name":"Ada Lovelace",
				"email":"ada@example.com",
				"password":"password123"
			}`, trustedAuthHeaders())
			if response.Code != http.StatusOK {
				t.Fatalf("signup status = %d body=%s, want ok", response.Code, response.Body.String())
			}
			assertNoVerificationCode(t, response)
			assertVerificationCount(t, db, verificationIdentifier("ada@example.com"), 0)
		})

		t.Run(tc.name+" does not update", func(t *testing.T) {
			router, db := newAuthTestRouter(t)
			createStatusUser(t, db, "ada@example.com", tc.status, false, tc.softDelete)
			createEmailVerification(t, db, "ada@example.com")
			before := captureVerificationState(t, db, verificationIdentifier("ada@example.com"))

			response := authJSON(t, router, http.MethodPost, "/api/auth/sign-up/email", `{
				"name":"Ada Lovelace",
				"email":"ada@example.com",
				"password":"password123"
			}`, trustedAuthHeaders())
			if response.Code != http.StatusOK {
				t.Fatalf("signup status = %d body=%s, want ok", response.Code, response.Body.String())
			}
			assertNoVerificationCode(t, response)
			assertVerificationCount(t, db, verificationIdentifier("ada@example.com"), 1)
			after := captureVerificationState(t, db, verificationIdentifier("ada@example.com"))
			assertVerificationStateUnchanged(t, before, after)
		})
	}
}

func TestDuplicateSignUpRotatesVerificationForInvitedAndActiveUsers(t *testing.T) {
	tests := []struct {
		name       string
		createUser func(*testing.T, *gorm.DB)
	}{
		{
			name: "invited",
			createUser: func(t *testing.T, db *gorm.DB) {
				createInvitedUser(t, db, "ada@example.com")
			},
		},
		{
			name: "active",
			createUser: func(t *testing.T, db *gorm.DB) {
				createStatusUser(t, db, "ada@example.com", "active", false, false)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			router, db := newAuthTestRouter(t)
			tc.createUser(t, db)

			response := authJSON(t, router, http.MethodPost, "/api/auth/sign-up/email", `{
				"name":"Ada Lovelace",
				"email":"ada@example.com",
				"password":"password123"
			}`, trustedAuthHeaders())
			if response.Code != http.StatusOK {
				t.Fatalf("signup status = %d body=%s, want ok", response.Code, response.Body.String())
			}
			var body struct {
				VerificationCode      string `json:"verificationCode"`
				VerificationExpiresAt string `json:"verificationExpiresAt"`
			}
			decodeJSON(t, response, &body)
			if len(body.VerificationCode) != 6 || body.VerificationExpiresAt == "" {
				t.Fatalf("unexpected signup verification response: %#v", body)
			}
			assertVerificationCount(t, db, verificationIdentifier("ada@example.com"), 1)
		})
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

func TestInactiveAndSuspendedUsersCannotSignIn(t *testing.T) {
	tests := []struct {
		name   string
		email  string
		status string
	}{
		{name: "inactive", email: "inactive@example.com", status: "inactive"},
		{name: "suspended", email: "suspended@example.com", status: "suspended"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			router, db := newAuthTestRouter(t)
			user := createVerifiedUser(t, db, tc.email, "password123")
			if err := db.Model(&auth.User{}).Where("id = ?", user.ID).Update("status", tc.status).Error; err != nil {
				t.Fatalf("set %s: %v", tc.status, err)
			}

			response := authJSON(t, router, http.MethodPost, "/api/auth/sign-in/email", `{
				"email":"`+tc.email+`",
				"password":"password123"
			}`, map[string]string{"X-Syncra-Internal-Token": testInternalToken})
			if response.Code != http.StatusForbidden {
				t.Fatalf("status = %d body=%s, want forbidden", response.Code, response.Body.String())
			}
		})
	}
}

func TestDeletedUserCannotSignIn(t *testing.T) {
	router, db := newAuthTestRouter(t)
	user := createVerifiedUser(t, db, "deleted@example.com", "password123")
	now := time.Now().UTC()
	if err := db.Model(&auth.User{}).Where("id = ?", user.ID).Updates(map[string]any{
		"status":     "deleted",
		"deleted_at": now,
	}).Error; err != nil {
		t.Fatalf("delete user: %v", err)
	}

	response := authJSON(t, router, http.MethodPost, "/api/auth/sign-in/email", `{
		"email":"deleted@example.com",
		"password":"password123"
	}`, map[string]string{"X-Syncra-Internal-Token": testInternalToken})
	if response.Code != http.StatusForbidden {
		t.Fatalf("status = %d body=%s, want forbidden", response.Code, response.Body.String())
	}
}

func TestSoftDeletedActiveUserCannotSignIn(t *testing.T) {
	router, db := newAuthTestRouter(t)
	user := createVerifiedUser(t, db, "soft-deleted@example.com", "password123")
	softDeleteUser(t, db, user.ID)

	response := authJSON(t, router, http.MethodPost, "/api/auth/sign-in/email", `{
		"email":"soft-deleted@example.com",
		"password":"password123"
	}`, map[string]string{"X-Syncra-Internal-Token": testInternalToken})
	if response.Code != http.StatusForbidden {
		t.Fatalf("status = %d body=%s, want forbidden", response.Code, response.Body.String())
	}
}

func TestEmptyStatusUserCanSignInUnlessSoftDeleted(t *testing.T) {
	router, db := newAuthTestRouter(t)
	user := createVerifiedUser(t, db, "empty-status@example.com", "password123")
	if err := db.Exec("PRAGMA ignore_check_constraints = ON").Error; err != nil {
		t.Fatalf("disable sqlite check constraints: %v", err)
	}
	if err := db.Model(&auth.User{}).Where("id = ?", user.ID).Update("status", "").Error; err != nil {
		t.Fatalf("clear status: %v", err)
	}

	active := authJSON(t, router, http.MethodPost, "/api/auth/sign-in/email", `{
		"email":"empty-status@example.com",
		"password":"password123"
	}`, map[string]string{"X-Syncra-Internal-Token": testInternalToken})
	if active.Code != http.StatusOK {
		t.Fatalf("active status = %d body=%s, want ok", active.Code, active.Body.String())
	}

	now := time.Now().UTC()
	if err := db.Model(&auth.User{}).Where("id = ?", user.ID).Update("deleted_at", now).Error; err != nil {
		t.Fatalf("soft delete user: %v", err)
	}
	softDeleted := authJSON(t, router, http.MethodPost, "/api/auth/sign-in/email", `{
		"email":"empty-status@example.com",
		"password":"password123"
	}`, map[string]string{"X-Syncra-Internal-Token": testInternalToken})
	if softDeleted.Code != http.StatusForbidden {
		t.Fatalf("soft deleted status = %d body=%s, want forbidden", softDeleted.Code, softDeleted.Body.String())
	}
}

func TestSuspendedUserSessionIsRejected(t *testing.T) {
	router, db := newAuthTestRouter(t)
	user := createVerifiedUser(t, db, "ada@example.com", "password123")
	token := loginUser(t, router, user.Email, "password123")
	if err := db.Model(&auth.User{}).Where("id = ?", user.ID).Update("status", "suspended").Error; err != nil {
		t.Fatalf("suspend user: %v", err)
	}

	response := authJSON(t, router, http.MethodGet, "/api/auth/get-session", "", map[string]string{
		"X-Syncra-Internal-Token": testInternalToken,
		"Cookie":                  "auth.session_token=" + token,
	})
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", response.Code, response.Body.String())
	}
	if strings.TrimSpace(response.Body.String()) != "null" {
		t.Fatalf("body = %s, want null", response.Body.String())
	}

	var sessionCount int64
	if err := db.Model(&auth.Session{}).Where("token = ?", token).Count(&sessionCount).Error; err != nil {
		t.Fatalf("count session: %v", err)
	}
	if sessionCount != 0 {
		t.Fatalf("session count = %d, want 0", sessionCount)
	}
}

func TestSoftDeletedActiveUserSessionIsRejected(t *testing.T) {
	router, db := newAuthTestRouter(t)
	user := createVerifiedUser(t, db, "ada@example.com", "password123")
	token := loginUser(t, router, user.Email, "password123")
	softDeleteUser(t, db, user.ID)

	response := authJSON(t, router, http.MethodGet, "/api/auth/get-session", "", map[string]string{
		"X-Syncra-Internal-Token": testInternalToken,
		"Cookie":                  "auth.session_token=" + token,
	})
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", response.Code, response.Body.String())
	}
	if strings.TrimSpace(response.Body.String()) != "null" {
		t.Fatalf("body = %s, want null", response.Body.String())
	}

	var sessionCount int64
	if err := db.Model(&auth.Session{}).Where("token = ?", token).Count(&sessionCount).Error; err != nil {
		t.Fatalf("count session: %v", err)
	}
	if sessionCount != 0 {
		t.Fatalf("session count = %d, want 0", sessionCount)
	}
}

func TestOAuthCallbackRejectsInvitedUserThatRemainsInvited(t *testing.T) {
	router, db := newAuthTestRouterWithOptions(t, RouterOptions{
		GoogleClientID:     "google-client",
		GoogleClientSecret: "google-secret",
		OAuthProfileFetcher: func(_ context.Context, providerID string, code string, redirectURI string) (OAuthProfile, error) {
			return OAuthProfile{
				ProviderID: providerID,
				AccountID:  "google-invited-account",
				Email:      "ada@example.com",
				Name:       "Ada Lovelace",
				Verified:   false,
			}, nil
		},
	})
	createInvitedUser(t, db, "ada@example.com")

	response := authJSON(t, router, http.MethodPost, "/api/auth/oauth/google/callback", `{
		"code":"oauth-code",
		"state":"state",
		"redirectURI":"http://localhost:5173/api/auth/google/callback"
	}`, map[string]string{"X-Syncra-Internal-Token": testInternalToken})
	if response.Code != http.StatusForbidden {
		t.Fatalf("oauth callback status = %d body=%s, want forbidden", response.Code, response.Body.String())
	}

	var sessionCount int64
	if err := db.Model(&auth.Session{}).Count(&sessionCount).Error; err != nil {
		t.Fatalf("count sessions: %v", err)
	}
	if sessionCount != 0 {
		t.Fatalf("session count = %d, want 0", sessionCount)
	}

	var user auth.User
	if err := db.First(&user, "email = ?", "ada@example.com").Error; err != nil {
		t.Fatalf("load invited user: %v", err)
	}
	if user.Status != "invited" {
		t.Fatalf("user status = %q, want invited", user.Status)
	}
	assertOAuthAccountCount(t, db, auth.GoogleProviderID, "google-invited-account", 0)
}

func TestOAuthCallbackRejectsNewUnverifiedUserWithoutSideEffects(t *testing.T) {
	router, db := newAuthTestRouterWithOptions(t, RouterOptions{
		GoogleClientID:     "google-client",
		GoogleClientSecret: "google-secret",
		OAuthProfileFetcher: func(_ context.Context, providerID string, code string, redirectURI string) (OAuthProfile, error) {
			return OAuthProfile{
				ProviderID: providerID,
				AccountID:  "google-new-unverified",
				Email:      "new-user@example.com",
				Name:       "New User",
				Verified:   false,
			}, nil
		},
	})

	response := authJSON(t, router, http.MethodPost, "/api/auth/oauth/google/callback", `{
		"code":"oauth-code",
		"state":"state",
		"redirectURI":"http://localhost:5173/api/auth/google/callback"
	}`, map[string]string{"X-Syncra-Internal-Token": testInternalToken})
	if response.Code != http.StatusForbidden {
		t.Fatalf("oauth callback status = %d body=%s, want forbidden", response.Code, response.Body.String())
	}
	assertSessionCount(t, db, 0)
	assertOAuthAccountCount(t, db, auth.GoogleProviderID, "google-new-unverified", 0)

	var userCount int64
	if err := db.Model(&auth.User{}).Where("email = ?", "new-user@example.com").Count(&userCount).Error; err != nil {
		t.Fatalf("count users: %v", err)
	}
	if userCount != 0 {
		t.Fatalf("user count = %d, want 0", userCount)
	}
}

func TestOAuthCallbackRejectsUnverifiedEmailMatchWithoutSideEffects(t *testing.T) {
	router, db := newAuthTestRouterWithOptions(t, RouterOptions{
		GoogleClientID:     "google-client",
		GoogleClientSecret: "google-secret",
		OAuthProfileFetcher: func(_ context.Context, providerID string, code string, redirectURI string) (OAuthProfile, error) {
			return OAuthProfile{
				ProviderID: providerID,
				AccountID:  "google-unverified-email-match",
				Email:      "ada@example.com",
				Name:       "Updated Ada",
				Verified:   false,
			}, nil
		},
	})
	user := createStatusUser(t, db, "ada@example.com", "active", true, false)
	before := captureAuthUserState(t, db, user.ID)

	response := authJSON(t, router, http.MethodPost, "/api/auth/oauth/google/callback", `{
		"code":"oauth-code",
		"state":"state",
		"redirectURI":"http://localhost:5173/api/auth/google/callback"
	}`, map[string]string{"X-Syncra-Internal-Token": testInternalToken})
	if response.Code != http.StatusForbidden {
		t.Fatalf("oauth callback status = %d body=%s, want forbidden", response.Code, response.Body.String())
	}
	assertSessionCount(t, db, 0)
	assertOAuthAccountCount(t, db, auth.GoogleProviderID, "google-unverified-email-match", 0)
	after := captureAuthUserState(t, db, user.ID)
	assertAuthUserStateUnchanged(t, before, after)
}

func TestActiveUserCanSignIn(t *testing.T) {
	router, db := newAuthTestRouter(t)
	createVerifiedUser(t, db, "active@example.com", "password123")

	response := authJSON(t, router, http.MethodPost, "/api/auth/sign-in/email", `{
		"email":"active@example.com",
		"password":"password123"
	}`, map[string]string{"X-Syncra-Internal-Token": testInternalToken})
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s, want ok", response.Code, response.Body.String())
	}
}

func TestOAuthCallbackRejectsNonActiveExistingUsersWithoutSideEffects(t *testing.T) {
	tests := []struct {
		name       string
		status     string
		softDelete bool
	}{
		{name: "inactive", status: "inactive"},
		{name: "suspended", status: "suspended"},
		{name: "deleted status", status: "deleted"},
		{name: "soft deleted active", status: "active", softDelete: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			accountID := "google-" + strings.ReplaceAll(tc.name, " ", "-")
			router, db := newAuthTestRouterWithOptions(t, RouterOptions{
				GoogleClientID:     "google-client",
				GoogleClientSecret: "google-secret",
				OAuthProfileFetcher: func(_ context.Context, providerID string, code string, redirectURI string) (OAuthProfile, error) {
					return OAuthProfile{
						ProviderID: providerID,
						AccountID:  accountID,
						Email:      "ada@example.com",
						Name:       "Updated Ada",
						Verified:   true,
					}, nil
				},
			})
			user := createStatusUser(t, db, "ada@example.com", tc.status, false, tc.softDelete)
			before := captureAuthUserState(t, db, user.ID)

			response := authJSON(t, router, http.MethodPost, "/api/auth/oauth/google/callback", `{
				"code":"oauth-code",
				"state":"state",
				"redirectURI":"http://localhost:5173/api/auth/google/callback"
			}`, map[string]string{"X-Syncra-Internal-Token": testInternalToken})
			if response.Code != http.StatusForbidden {
				t.Fatalf("oauth callback status = %d body=%s, want forbidden", response.Code, response.Body.String())
			}
			assertSessionCount(t, db, 0)
			assertOAuthAccountCount(t, db, auth.GoogleProviderID, accountID, 0)
			after := captureAuthUserState(t, db, user.ID)
			assertAuthUserStateUnchanged(t, before, after)
		})
	}
}

func TestOAuthCallbackRejectsLinkedNonActiveUserWithoutMutation(t *testing.T) {
	router, db := newAuthTestRouterWithOptions(t, RouterOptions{
		GoogleClientID:     "google-client",
		GoogleClientSecret: "google-secret",
		OAuthProfileFetcher: func(_ context.Context, providerID string, code string, redirectURI string) (OAuthProfile, error) {
			return OAuthProfile{
				ProviderID: providerID,
				AccountID:  "google-linked-inactive",
				Email:      "ada@example.com",
				Name:       "Updated Ada",
				Verified:   true,
			}, nil
		},
	})
	user := createStatusUser(t, db, "ada@example.com", "inactive", false, false)
	now := time.Now().UTC()
	account := auth.AuthAccount{
		AccountID:  "google-linked-inactive",
		ProviderID: auth.GoogleProviderID,
		UserID:     user.ID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := db.Create(&account).Error; err != nil {
		t.Fatalf("create linked oauth account: %v", err)
	}
	before := captureAuthUserState(t, db, user.ID)

	response := authJSON(t, router, http.MethodPost, "/api/auth/oauth/google/callback", `{
		"code":"oauth-code",
		"state":"state",
		"redirectURI":"http://localhost:5173/api/auth/google/callback"
	}`, map[string]string{"X-Syncra-Internal-Token": testInternalToken})
	if response.Code != http.StatusForbidden {
		t.Fatalf("oauth callback status = %d body=%s, want forbidden", response.Code, response.Body.String())
	}
	assertSessionCount(t, db, 0)
	assertOAuthAccountCount(t, db, auth.GoogleProviderID, "google-linked-inactive", 1)
	after := captureAuthUserState(t, db, user.ID)
	assertAuthUserStateUnchanged(t, before, after)
}

func TestOAuthCallbackAllowsLinkedUnverifiedProfileWithoutPromotion(t *testing.T) {
	router, db := newAuthTestRouterWithOptions(t, RouterOptions{
		GoogleClientID:     "google-client",
		GoogleClientSecret: "google-secret",
		OAuthProfileFetcher: func(_ context.Context, providerID string, code string, redirectURI string) (OAuthProfile, error) {
			return OAuthProfile{
				ProviderID: providerID,
				AccountID:  "google-linked-unverified-profile",
				Email:      "ada@example.com",
				Name:       "Updated Ada",
				Verified:   false,
			}, nil
		},
	})
	user := createStatusUser(t, db, "ada@example.com", "active", false, false)
	now := time.Now().UTC()
	account := auth.AuthAccount{
		AccountID:  "google-linked-unverified-profile",
		ProviderID: auth.GoogleProviderID,
		UserID:     user.ID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := db.Create(&account).Error; err != nil {
		t.Fatalf("create linked oauth account: %v", err)
	}

	response := authJSON(t, router, http.MethodPost, "/api/auth/oauth/google/callback", `{
		"code":"oauth-code",
		"state":"state",
		"redirectURI":"http://localhost:5173/api/auth/google/callback"
	}`, map[string]string{"X-Syncra-Internal-Token": testInternalToken})
	if response.Code != http.StatusOK {
		t.Fatalf("oauth callback status = %d body=%s, want ok", response.Code, response.Body.String())
	}
	assertSessionCount(t, db, 1)
	assertOAuthAccountCount(t, db, auth.GoogleProviderID, "google-linked-unverified-profile", 1)

	var after auth.User
	if err := db.First(&after, "id = ?", user.ID).Error; err != nil {
		t.Fatalf("load user after oauth: %v", err)
	}
	if after.EmailVerified {
		t.Fatal("linked oauth callback promoted email verification for unverified profile")
	}
	if after.Status != "active" {
		t.Fatalf("status = %q, want active", after.Status)
	}
}

func TestOAuthCallbackPromotesLinkedInvitedOnlyWhenVerifiedEmailMatches(t *testing.T) {
	tests := []struct {
		name          string
		profileEmail  string
		wantStatus    int
		wantPromoted  bool
		wantSession   int64
		wantUnchanged bool
	}{
		{
			name:         "same email",
			profileEmail: "ada@example.com",
			wantStatus:   http.StatusOK,
			wantPromoted: true,
			wantSession:  1,
		},
		{
			name:          "different email",
			profileEmail:  "other@example.com",
			wantStatus:    http.StatusForbidden,
			wantSession:   0,
			wantUnchanged: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			accountID := "google-linked-email-" + strings.ReplaceAll(tc.name, " ", "-")
			router, db := newAuthTestRouterWithOptions(t, RouterOptions{
				GoogleClientID:     "google-client",
				GoogleClientSecret: "google-secret",
				OAuthProfileFetcher: func(_ context.Context, providerID string, code string, redirectURI string) (OAuthProfile, error) {
					return OAuthProfile{
						ProviderID: providerID,
						AccountID:  accountID,
						Email:      tc.profileEmail,
						Name:       "Ada Lovelace",
						Verified:   true,
					}, nil
				},
			})
			user := createInvitedUser(t, db, "ada@example.com")
			now := time.Now().UTC()
			account := auth.AuthAccount{
				AccountID:  accountID,
				ProviderID: auth.GoogleProviderID,
				UserID:     user.ID,
				CreatedAt:  now,
				UpdatedAt:  now,
			}
			if err := db.Create(&account).Error; err != nil {
				t.Fatalf("create linked oauth account: %v", err)
			}
			before := captureAuthUserState(t, db, user.ID)

			response := authJSON(t, router, http.MethodPost, "/api/auth/oauth/google/callback", `{
				"code":"oauth-code",
				"state":"state",
				"redirectURI":"http://localhost:5173/api/auth/google/callback"
			}`, map[string]string{"X-Syncra-Internal-Token": testInternalToken})
			if response.Code != tc.wantStatus {
				t.Fatalf("oauth callback status = %d body=%s, want %d", response.Code, response.Body.String(), tc.wantStatus)
			}
			assertSessionCount(t, db, tc.wantSession)
			assertOAuthAccountCount(t, db, auth.GoogleProviderID, accountID, 1)

			after := captureAuthUserState(t, db, user.ID)
			if tc.wantUnchanged {
				assertAuthUserStateUnchanged(t, before, after)
			}
			if tc.wantPromoted && (!after.EmailVerified || after.Status != "active") {
				t.Fatalf("after = %#v, want verified active user", after)
			}
		})
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

func TestPasswordResetRequestDoesNotRotateForNonActiveUsers(t *testing.T) {
	tests := []struct {
		name       string
		status     string
		softDelete bool
	}{
		{name: "inactive", status: "inactive"},
		{name: "suspended", status: "suspended"},
		{name: "deleted status", status: "deleted"},
		{name: "soft deleted active", status: "active", softDelete: true},
	}

	for _, tc := range tests {
		t.Run(tc.name+" does not create", func(t *testing.T) {
			router, db := newAuthTestRouter(t)
			user := createStatusUser(t, db, "ada@example.com", tc.status, true, tc.softDelete)
			before := captureAuthUserState(t, db, user.ID)

			response := authJSON(t, router, http.MethodPost, "/api/auth/password-reset/request", `{
				"email":"ada@example.com"
			}`, trustedAuthHeaders())
			if response.Code != http.StatusOK {
				t.Fatalf("reset request status = %d body=%s, want ok", response.Code, response.Body.String())
			}
			var body map[string]any
			decodeJSON(t, response, &body)
			if _, ok := body["resetToken"]; ok {
				t.Fatalf("response leaked resetToken for non-active user: %s", response.Body.String())
			}
			assertVerificationCount(t, db, passwordResetIdentifier("ada@example.com"), 0)
			after := captureAuthUserState(t, db, user.ID)
			assertAuthUserStateUnchanged(t, before, after)
		})

		t.Run(tc.name+" does not update", func(t *testing.T) {
			router, db := newAuthTestRouter(t)
			user := createStatusUser(t, db, "ada@example.com", tc.status, true, tc.softDelete)
			createPasswordResetVerification(t, db, "ada@example.com")
			beforeUser := captureAuthUserState(t, db, user.ID)
			beforeVerification := captureVerificationState(t, db, passwordResetIdentifier("ada@example.com"))

			response := authJSON(t, router, http.MethodPost, "/api/auth/password-reset/request", `{
				"email":"ada@example.com"
			}`, trustedAuthHeaders())
			if response.Code != http.StatusOK {
				t.Fatalf("reset request status = %d body=%s, want ok", response.Code, response.Body.String())
			}
			var body map[string]any
			decodeJSON(t, response, &body)
			if _, ok := body["resetToken"]; ok {
				t.Fatalf("response leaked resetToken for non-active user: %s", response.Body.String())
			}
			assertVerificationCount(t, db, passwordResetIdentifier("ada@example.com"), 1)
			afterUser := captureAuthUserState(t, db, user.ID)
			assertAuthUserStateUnchanged(t, beforeUser, afterUser)
			afterVerification := captureVerificationState(t, db, passwordResetIdentifier("ada@example.com"))
			assertVerificationStateUnchanged(t, beforeVerification, afterVerification)
		})
	}
}

func TestConfirmPasswordResetInvalidTokenDoesNotEnumerateLifecycle(t *testing.T) {
	tests := []struct {
		name       string
		email      string
		createUser func(*testing.T, *gorm.DB) (auth.User, bool)
	}{
		{
			name:  "active",
			email: "active@example.com",
			createUser: func(t *testing.T, db *gorm.DB) (auth.User, bool) {
				return createVerifiedUser(t, db, "active@example.com", "old-password"), true
			},
		},
		{
			name:  "missing",
			email: "missing@example.com",
			createUser: func(t *testing.T, db *gorm.DB) (auth.User, bool) {
				return auth.User{}, false
			},
		},
		{
			name:  "inactive",
			email: "inactive@example.com",
			createUser: func(t *testing.T, db *gorm.DB) (auth.User, bool) {
				user := createVerifiedUser(t, db, "inactive@example.com", "old-password")
				setUserLifecycleState(t, db, user.ID, "inactive", false)
				return user, true
			},
		},
		{
			name:  "suspended",
			email: "suspended@example.com",
			createUser: func(t *testing.T, db *gorm.DB) (auth.User, bool) {
				user := createVerifiedUser(t, db, "suspended@example.com", "old-password")
				setUserLifecycleState(t, db, user.ID, "suspended", false)
				return user, true
			},
		},
		{
			name:  "deleted status",
			email: "deleted@example.com",
			createUser: func(t *testing.T, db *gorm.DB) (auth.User, bool) {
				user := createVerifiedUser(t, db, "deleted@example.com", "old-password")
				setUserLifecycleState(t, db, user.ID, "deleted", false)
				return user, true
			},
		},
		{
			name:  "soft deleted active",
			email: "soft-deleted@example.com",
			createUser: func(t *testing.T, db *gorm.DB) (auth.User, bool) {
				user := createVerifiedUser(t, db, "soft-deleted@example.com", "old-password")
				setUserLifecycleState(t, db, user.ID, "active", true)
				return user, true
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			router, db := newAuthTestRouter(t)
			user, hasUser := tc.createUser(t, db)
			token := createPasswordResetVerification(t, db, tc.email)
			if token == "wrong-token" {
				t.Fatal("test setup used wrong token")
			}
			var before authUserState
			if hasUser {
				before = captureAuthUserState(t, db, user.ID)
			}

			response := authJSON(t, router, http.MethodPost, "/api/auth/password-reset/confirm", `{
				"email":"`+tc.email+`",
				"token":"wrong-token",
				"password":"new-password"
			}`, map[string]string{"X-Syncra-Internal-Token": testInternalToken})
			assertErrorResponse(t, response, http.StatusBadRequest, "invalid password reset token")
			assertVerificationCount(t, db, passwordResetIdentifier(tc.email), 1)
			if hasUser {
				assertCredentialPassword(t, db, user.ID, "old-password")
				after := captureAuthUserState(t, db, user.ID)
				assertAuthUserStateUnchanged(t, before, after)
			}
		})
	}
}

func TestPasswordResetConfirmRejectsNonActiveUsersWithoutMutation(t *testing.T) {
	tests := []struct {
		name       string
		status     string
		softDelete bool
	}{
		{name: "inactive", status: "inactive"},
		{name: "suspended", status: "suspended"},
		{name: "deleted status", status: "deleted"},
		{name: "soft deleted active", status: "active", softDelete: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			router, db := newAuthTestRouter(t)
			user := createVerifiedUser(t, db, "ada@example.com", "old-password")
			setUserLifecycleState(t, db, user.ID, tc.status, tc.softDelete)
			token := createPasswordResetVerification(t, db, "ada@example.com")
			before := captureAuthUserState(t, db, user.ID)

			response := authJSON(t, router, http.MethodPost, "/api/auth/password-reset/confirm", `{
				"email":"ada@example.com",
				"token":"`+token+`",
				"password":"new-password"
			}`, map[string]string{"X-Syncra-Internal-Token": testInternalToken})
			if response.Code != http.StatusForbidden {
				t.Fatalf("reset confirm status = %d body=%s, want forbidden", response.Code, response.Body.String())
			}
			assertVerificationCount(t, db, passwordResetIdentifier("ada@example.com"), 1)
			assertCredentialPassword(t, db, user.ID, "old-password")
			after := captureAuthUserState(t, db, user.ID)
			assertAuthUserStateUnchanged(t, before, after)
		})
	}
}

func TestPasswordResetPromotesInvitedUser(t *testing.T) {
	router, db := newAuthTestRouter(t)
	user := createInvitedUser(t, db, "ada@example.com")

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

	var promoted auth.User
	if err := db.First(&promoted, "id = ?", user.ID).Error; err != nil {
		t.Fatalf("load promoted user: %v", err)
	}
	if !promoted.EmailVerified {
		t.Fatal("password reset did not mark user email verified")
	}
	if promoted.Status != "active" {
		t.Fatalf("password reset promoted user status = %q, want active", promoted.Status)
	}
}

func TestVerifyEmailOTPInvalidTokenDoesNotEnumerateLifecycle(t *testing.T) {
	tests := []struct {
		name       string
		email      string
		createUser func(*testing.T, *gorm.DB) (auth.User, bool)
	}{
		{
			name:  "active",
			email: "active@example.com",
			createUser: func(t *testing.T, db *gorm.DB) (auth.User, bool) {
				return createStatusUser(t, db, "active@example.com", "active", false, false), true
			},
		},
		{
			name:  "missing",
			email: "missing@example.com",
			createUser: func(t *testing.T, db *gorm.DB) (auth.User, bool) {
				return auth.User{}, false
			},
		},
		{
			name:  "inactive",
			email: "inactive@example.com",
			createUser: func(t *testing.T, db *gorm.DB) (auth.User, bool) {
				return createStatusUser(t, db, "inactive@example.com", "inactive", false, false), true
			},
		},
		{
			name:  "suspended",
			email: "suspended@example.com",
			createUser: func(t *testing.T, db *gorm.DB) (auth.User, bool) {
				return createStatusUser(t, db, "suspended@example.com", "suspended", false, false), true
			},
		},
		{
			name:  "deleted status",
			email: "deleted@example.com",
			createUser: func(t *testing.T, db *gorm.DB) (auth.User, bool) {
				return createStatusUser(t, db, "deleted@example.com", "deleted", false, false), true
			},
		},
		{
			name:  "soft deleted active",
			email: "soft-deleted@example.com",
			createUser: func(t *testing.T, db *gorm.DB) (auth.User, bool) {
				return createStatusUser(t, db, "soft-deleted@example.com", "active", false, true), true
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			router, db := newAuthTestRouter(t)
			user, hasUser := tc.createUser(t, db)
			code := createEmailVerification(t, db, tc.email)
			if code == "wrong-code" {
				t.Fatal("test setup used wrong code")
			}
			var before authUserState
			if hasUser {
				before = captureAuthUserState(t, db, user.ID)
			}

			response := authJSON(t, router, http.MethodPost, "/api/auth/email-otp/verify-email", `{
				"email":"`+tc.email+`",
				"otp":"wrong-code"
			}`, map[string]string{"X-Syncra-Internal-Token": testInternalToken})
			assertErrorResponse(t, response, http.StatusBadRequest, "invalid verification code")
			assertVerificationCount(t, db, verificationIdentifier(tc.email), 1)
			if hasUser {
				after := captureAuthUserState(t, db, user.ID)
				assertAuthUserStateUnchanged(t, before, after)
			}
		})
	}
}

func TestVerifyEmailOTPRejectsNonActiveUsersWithoutMutation(t *testing.T) {
	tests := []struct {
		name       string
		status     string
		softDelete bool
	}{
		{name: "inactive", status: "inactive"},
		{name: "suspended", status: "suspended"},
		{name: "deleted status", status: "deleted"},
		{name: "soft deleted active", status: "active", softDelete: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			router, db := newAuthTestRouter(t)
			user := createStatusUser(t, db, "ada@example.com", tc.status, false, tc.softDelete)
			code := createEmailVerification(t, db, "ada@example.com")
			before := captureAuthUserState(t, db, user.ID)

			response := authJSON(t, router, http.MethodPost, "/api/auth/email-otp/verify-email", `{
				"email":"ada@example.com",
				"otp":"`+code+`"
			}`, map[string]string{"X-Syncra-Internal-Token": testInternalToken})
			if response.Code != http.StatusForbidden {
				t.Fatalf("verify status = %d body=%s, want forbidden", response.Code, response.Body.String())
			}
			assertVerificationCount(t, db, verificationIdentifier("ada@example.com"), 1)
			after := captureAuthUserState(t, db, user.ID)
			assertAuthUserStateUnchanged(t, before, after)
		})
	}
}

func TestVerifyEmailOTPPreservesInvitedAndActivePromotion(t *testing.T) {
	tests := []struct {
		name       string
		createUser func(*testing.T, *gorm.DB) auth.User
	}{
		{
			name: "invited",
			createUser: func(t *testing.T, db *gorm.DB) auth.User {
				return createInvitedUser(t, db, "ada@example.com")
			},
		},
		{
			name: "active",
			createUser: func(t *testing.T, db *gorm.DB) auth.User {
				now := time.Now().UTC()
				user := auth.User{
					Name:          "Ada Lovelace",
					Email:         "ada@example.com",
					EmailVerified: false,
					Status:        "active",
					CreatedAt:     now,
					UpdatedAt:     now,
				}
				if err := db.Create(&user).Error; err != nil {
					t.Fatalf("create active user: %v", err)
				}
				return user
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			router, db := newAuthTestRouter(t)
			user := tc.createUser(t, db)
			code := createEmailVerification(t, db, "ada@example.com")

			response := authJSON(t, router, http.MethodPost, "/api/auth/email-otp/verify-email", `{
				"email":"ada@example.com",
				"otp":"`+code+`"
			}`, map[string]string{"X-Syncra-Internal-Token": testInternalToken})
			if response.Code != http.StatusOK {
				t.Fatalf("verify status = %d body=%s", response.Code, response.Body.String())
			}

			var verified auth.User
			if err := db.First(&verified, "id = ?", user.ID).Error; err != nil {
				t.Fatalf("load verified user: %v", err)
			}
			if !verified.EmailVerified {
				t.Fatal("email verified = false, want true")
			}
			if verified.Status != "active" {
				t.Fatalf("status = %q, want active", verified.Status)
			}
		})
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

func TestFetchGitHubProfileUsesVerifiedPrimaryEmail(t *testing.T) {
	handler := newGitHubProfileTestHandler(t, `[
		{"email":"secondary@example.com","primary":false,"verified":true},
		{"email":"primary@example.com","primary":true,"verified":true}
	]`)

	profile, err := handler.fetchGitHubProfile(context.Background(), "oauth-code", "http://localhost/callback")
	if err != nil {
		t.Fatalf("fetchGitHubProfile() error = %v", err)
	}
	if profile.ProviderID != auth.GitHubProviderID || profile.AccountID != "123" || profile.Email != "primary@example.com" || !profile.Verified {
		t.Fatalf("profile = %#v, want verified primary email profile", profile)
	}
	if profile.Name != "Ada Lovelace" {
		t.Fatalf("profile name = %q, want Ada Lovelace", profile.Name)
	}
}

func TestFetchGitHubProfileRequiresVerifiedPrimaryEmail(t *testing.T) {
	tests := []struct {
		name       string
		emailsJSON string
	}{
		{name: "missing emails", emailsJSON: `[]`},
		{name: "primary unverified", emailsJSON: `[{"email":"primary@example.com","primary":true,"verified":false}]`},
		{name: "verified non primary", emailsJSON: `[{"email":"secondary@example.com","primary":false,"verified":true}]`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := newGitHubProfileTestHandler(t, tc.emailsJSON)

			profile, err := handler.fetchGitHubProfile(context.Background(), "oauth-code", "http://localhost/callback")
			if err == nil {
				t.Fatalf("fetchGitHubProfile() error = nil profile=%#v, want verified primary email error", profile)
			}
		})
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

func TestOAuthCallbackPromotesExistingVerifiedInvitedEmailUser(t *testing.T) {
	router, db := newAuthTestRouterWithOptions(t, RouterOptions{
		GoogleClientID:     "google-client",
		GoogleClientSecret: "google-secret",
		OAuthProfileFetcher: func(_ context.Context, providerID string, code string, redirectURI string) (OAuthProfile, error) {
			return OAuthProfile{
				ProviderID: providerID,
				AccountID:  "google-account-verified-invited",
				Email:      "ada@example.com",
				Name:       "Ada Lovelace",
				Verified:   true,
			}, nil
		},
	})
	user := createVerifiedInvitedUser(t, db, "ada@example.com")

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
		t.Fatal("oauth callback changed user email verified to false")
	}
	if promoted.Status != "active" {
		t.Fatalf("oauth promoted user status = %q, want active", promoted.Status)
	}
}

func TestOAuthCallbackPromotesExistingLinkedInvitedVerifiedEmailUser(t *testing.T) {
	router, db := newAuthTestRouterWithOptions(t, RouterOptions{
		GoogleClientID:     "google-client",
		GoogleClientSecret: "google-secret",
		OAuthProfileFetcher: func(_ context.Context, providerID string, code string, redirectURI string) (OAuthProfile, error) {
			return OAuthProfile{
				ProviderID: providerID,
				AccountID:  "google-linked-account",
				Email:      "ada@example.com",
				Name:       "Ada Lovelace",
				Verified:   true,
			}, nil
		},
	})
	user := createInvitedUser(t, db, "ada@example.com")
	now := time.Now().UTC()
	account := auth.AuthAccount{
		AccountID:  "google-linked-account",
		ProviderID: auth.GoogleProviderID,
		UserID:     user.ID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := db.Create(&account).Error; err != nil {
		t.Fatalf("create linked oauth account: %v", err)
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
		t.Fatal("linked oauth callback did not mark user email verified")
	}
	if promoted.Status != "active" {
		t.Fatalf("linked oauth promoted user status = %q, want active", promoted.Status)
	}
}

func TestOAuthCallbackPromotesExistingLinkedVerifiedInvitedEmailUser(t *testing.T) {
	router, db := newAuthTestRouterWithOptions(t, RouterOptions{
		GoogleClientID:     "google-client",
		GoogleClientSecret: "google-secret",
		OAuthProfileFetcher: func(_ context.Context, providerID string, code string, redirectURI string) (OAuthProfile, error) {
			return OAuthProfile{
				ProviderID: providerID,
				AccountID:  "google-linked-verified-invited",
				Email:      "ada@example.com",
				Name:       "Ada Lovelace",
				Verified:   true,
			}, nil
		},
	})
	user := createVerifiedInvitedUser(t, db, "ada@example.com")
	now := time.Now().UTC()
	account := auth.AuthAccount{
		AccountID:  "google-linked-verified-invited",
		ProviderID: auth.GoogleProviderID,
		UserID:     user.ID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := db.Create(&account).Error; err != nil {
		t.Fatalf("create linked oauth account: %v", err)
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
		t.Fatal("linked oauth callback changed user email verified to false")
	}
	if promoted.Status != "active" {
		t.Fatalf("linked oauth promoted user status = %q, want active", promoted.Status)
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

func newGitHubProfileTestHandler(t *testing.T, emailsJSON string) *authHandler {
	t.Helper()
	client := &http.Client{
		Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
			switch request.URL.String() {
			case "https://github.com/login/oauth/access_token":
				return jsonHTTPResponse(http.StatusOK, `{"access_token":"github-token"}`), nil
			case "https://api.github.com/user":
				if got := request.Header.Get("Authorization"); got != "Bearer github-token" {
					t.Fatalf("github user authorization = %q, want bearer token", got)
				}
				return jsonHTTPResponse(http.StatusOK, `{"id":123,"login":"ada","name":"Ada Lovelace","email":"","avatar_url":"https://example.com/avatar.png"}`), nil
			case "https://api.github.com/user/emails":
				if got := request.Header.Get("Authorization"); got != "Bearer github-token" {
					t.Fatalf("github emails authorization = %q, want bearer token", got)
				}
				return jsonHTTPResponse(http.StatusOK, emailsJSON), nil
			default:
				t.Fatalf("unexpected request URL: %s", request.URL.String())
				return nil, nil
			}
		}),
	}
	return &authHandler{
		gitHubClientID:     "github-client",
		gitHubClientSecret: "github-secret",
		httpClient:         client,
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func jsonHTTPResponse(status int, body string) *http.Response {
	header := make(http.Header)
	header.Set("Content-Type", "application/json")
	return &http.Response{
		StatusCode: status,
		Header:     header,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
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

func assertErrorResponse(t *testing.T, response *httptest.ResponseRecorder, status int, message string) {
	t.Helper()
	if response.Code != status {
		t.Fatalf("status = %d body=%s, want %d", response.Code, response.Body.String(), status)
	}
	var body struct {
		Error string `json:"error"`
	}
	decodeJSON(t, response, &body)
	if body.Error != message {
		t.Fatalf("error = %q, want %q; body=%s", body.Error, message, response.Body.String())
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

func createInvitedUser(t *testing.T, db *gorm.DB, email string) auth.User {
	t.Helper()
	now := time.Now().UTC()
	user := auth.User{
		Name:          "Ada Lovelace",
		Email:         email,
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
	return user
}

func createVerifiedInvitedUser(t *testing.T, db *gorm.DB, email string) auth.User {
	t.Helper()
	now := time.Now().UTC()
	user := auth.User{
		Name:          "Ada Lovelace",
		Email:         email,
		EmailVerified: true,
		Status:        "invited",
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create verified invited user: %v", err)
	}
	if !user.EmailVerified {
		t.Fatal("created user email verified = false, want true")
	}
	if user.Status != "invited" {
		t.Fatalf("created user status = %q, want invited", user.Status)
	}
	return user
}

func createStatusUser(t *testing.T, db *gorm.DB, email string, status string, emailVerified bool, softDelete bool) auth.User {
	t.Helper()
	now := time.Now().UTC()
	user := auth.User{
		Name:          "Ada Lovelace",
		Email:         email,
		EmailVerified: emailVerified,
		Status:        status,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if softDelete {
		user.DeletedAt = &now
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create status user: %v", err)
	}
	return user
}

func softDeleteUser(t *testing.T, db *gorm.DB, userID string) {
	t.Helper()
	setUserLifecycleState(t, db, userID, "active", true)
}

func setUserLifecycleState(t *testing.T, db *gorm.DB, userID string, status string, softDelete bool) {
	t.Helper()
	updates := map[string]any{"status": status}
	if softDelete {
		now := time.Now().UTC()
		updates["deleted_at"] = now
	} else {
		updates["deleted_at"] = nil
	}
	if err := db.Model(&auth.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		t.Fatalf("set user lifecycle state: %v", err)
	}
}

type authUserState struct {
	Name          string
	Email         string
	EmailVerified bool
	Status        string
	DeletedAt     *time.Time
	LastLoginAt   *time.Time
	UpdatedAt     time.Time
}

func captureAuthUserState(t *testing.T, db *gorm.DB, userID string) authUserState {
	t.Helper()
	var user auth.User
	if err := db.First(&user, "id = ?", userID).Error; err != nil {
		t.Fatalf("load user state: %v", err)
	}
	return authUserState{
		Name:          user.Name,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		Status:        user.Status,
		DeletedAt:     cloneTimePtr(user.DeletedAt),
		LastLoginAt:   cloneTimePtr(user.LastLoginAt),
		UpdatedAt:     user.UpdatedAt,
	}
}

func assertAuthUserStateUnchanged(t *testing.T, before authUserState, after authUserState) {
	t.Helper()
	if before.Name != after.Name ||
		before.Email != after.Email ||
		before.EmailVerified != after.EmailVerified ||
		before.Status != after.Status ||
		!timePtrEqual(before.DeletedAt, after.DeletedAt) ||
		!timePtrEqual(before.LastLoginAt, after.LastLoginAt) ||
		!before.UpdatedAt.Equal(after.UpdatedAt) {
		t.Fatalf("user state changed\nbefore=%#v\nafter=%#v", before, after)
	}
}

func cloneTimePtr(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	clone := *value
	return &clone
}

func timePtrEqual(left *time.Time, right *time.Time) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return left.Equal(*right)
}

func createEmailVerification(t *testing.T, db *gorm.DB, email string) string {
	t.Helper()
	code := "123456"
	createVerification(t, db, verificationIdentifier(email), code)
	return code
}

func createPasswordResetVerification(t *testing.T, db *gorm.DB, email string) string {
	t.Helper()
	token := "reset-token"
	createVerification(t, db, passwordResetIdentifier(email), token)
	return token
}

func createVerification(t *testing.T, db *gorm.DB, identifier string, plaintext string) {
	t.Helper()
	now := time.Now().UTC()
	verification := auth.Verification{
		Identifier: identifier,
		Value:      auth.HashCode(testBetterAuthSecret, identifier, plaintext),
		ExpiresAt:  now.Add(5 * time.Minute),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := db.Create(&verification).Error; err != nil {
		t.Fatalf("create verification: %v", err)
	}
}

type verificationState struct {
	Value     string
	ExpiresAt time.Time
	UpdatedAt time.Time
}

func captureVerificationState(t *testing.T, db *gorm.DB, identifier string) verificationState {
	t.Helper()
	var verification auth.Verification
	if err := db.First(&verification, "identifier = ?", identifier).Error; err != nil {
		t.Fatalf("load verification state: %v", err)
	}
	return verificationState{
		Value:     verification.Value,
		ExpiresAt: verification.ExpiresAt,
		UpdatedAt: verification.UpdatedAt,
	}
}

func assertVerificationStateUnchanged(t *testing.T, before verificationState, after verificationState) {
	t.Helper()
	if before.Value != after.Value || !before.ExpiresAt.Equal(after.ExpiresAt) || !before.UpdatedAt.Equal(after.UpdatedAt) {
		t.Fatalf("verification state changed\nbefore=%#v\nafter=%#v", before, after)
	}
}

func assertNoVerificationCode(t *testing.T, response *httptest.ResponseRecorder) {
	t.Helper()
	var body map[string]any
	decodeJSON(t, response, &body)
	if _, ok := body["verificationCode"]; ok {
		t.Fatalf("response leaked verificationCode: %s", response.Body.String())
	}
	if _, ok := body["verificationExpiresAt"]; ok {
		t.Fatalf("response leaked verificationExpiresAt: %s", response.Body.String())
	}
}

func assertSessionCount(t *testing.T, db *gorm.DB, want int64) {
	t.Helper()
	var count int64
	if err := db.Model(&auth.Session{}).Count(&count).Error; err != nil {
		t.Fatalf("count sessions: %v", err)
	}
	if count != want {
		t.Fatalf("session count = %d, want %d", count, want)
	}
}

func assertOAuthAccountCount(t *testing.T, db *gorm.DB, providerID string, accountID string, want int64) {
	t.Helper()
	var count int64
	if err := db.Model(&auth.AuthAccount{}).Where("provider_id = ? AND account_id = ?", providerID, accountID).Count(&count).Error; err != nil {
		t.Fatalf("count oauth accounts: %v", err)
	}
	if count != want {
		t.Fatalf("oauth account count = %d, want %d", count, want)
	}
}

func assertVerificationCount(t *testing.T, db *gorm.DB, identifier string, want int64) {
	t.Helper()
	var count int64
	if err := db.Model(&auth.Verification{}).Where("identifier = ?", identifier).Count(&count).Error; err != nil {
		t.Fatalf("count verifications: %v", err)
	}
	if count != want {
		t.Fatalf("verification count = %d, want %d", count, want)
	}
}

func assertCredentialPassword(t *testing.T, db *gorm.DB, userID string, password string) {
	t.Helper()
	var account auth.AuthAccount
	if err := db.First(&account, "user_id = ? AND provider_id = ?", userID, auth.CredentialProviderID).Error; err != nil {
		t.Fatalf("load credential account: %v", err)
	}
	if !auth.VerifyPassword(password, account.Password) {
		t.Fatal("credential password was changed")
	}
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
