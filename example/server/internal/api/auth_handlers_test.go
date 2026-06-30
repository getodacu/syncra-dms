package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"ai.ro/syncra/internal/auth"
	authcrypto "ai.ro/syncra/internal/auth"
	"ai.ro/syncra/internal/billing"
)

const testAuthSecret = "J7mN2qR9vT4xY8bC6pL3wS5zK1hD0fG2"
const testDeliveryToken = "test-delivery-token"
const testOnboardingCredits = 125

func testAuthRouter(t *testing.T) (*gin.Engine, *gorm.DB) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	db := apiPostgresTx(t)
	h := &Handler{
		DB:                  db,
		BetterAuthSecret:    testAuthSecret,
		AuthDeliveryToken:   testDeliveryToken,
		InternalAPIToken:    testInternalAPIToken,
		AuthSessionTTL:      7 * 24 * time.Hour,
		AuthVerificationTTL: 5 * time.Minute,
		AuthCookieSecure:    false,
		OnboardingCredits:   testOnboardingCredits,
	}
	return NewRouter(h), db
}

func testGoogleAuthRouter(
	t *testing.T,
	exchange func(context.Context, string, string) (googleOAuthToken, error),
	validate func(context.Context, string, string) (googleIDTokenPayload, error),
) (*gin.Engine, *gorm.DB) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	db := apiPostgresTx(t)
	h := &Handler{
		DB:                    db,
		BetterAuthSecret:      testAuthSecret,
		AuthDeliveryToken:     testDeliveryToken,
		InternalAPIToken:      testInternalAPIToken,
		AuthSessionTTL:        7 * 24 * time.Hour,
		AuthVerificationTTL:   5 * time.Minute,
		AuthCookieSecure:      false,
		OnboardingCredits:     testOnboardingCredits,
		GoogleClientID:        "google-client-id",
		GoogleClientSecret:    "google-client-secret",
		GoogleOAuthExchange:   exchange,
		GoogleIDTokenValidate: validate,
	}
	return NewRouter(h), db
}

func testGitHubAuthRouter(
	t *testing.T,
	exchange func(context.Context, string, string) (githubOAuthToken, error),
	fetch func(context.Context, string) (githubOAuthProfile, error),
) (*gin.Engine, *gorm.DB) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	db := apiPostgresTx(t)
	h := &Handler{
		DB:                  db,
		BetterAuthSecret:    testAuthSecret,
		AuthDeliveryToken:   testDeliveryToken,
		InternalAPIToken:    testInternalAPIToken,
		AuthSessionTTL:      7 * 24 * time.Hour,
		AuthVerificationTTL: 5 * time.Minute,
		AuthCookieSecure:    false,
		OnboardingCredits:   testOnboardingCredits,
		GitHubClientID:      "github-client-id",
		GitHubClientSecret:  "github-client-secret",
		GitHubOAuthExchange: exchange,
		GitHubProfileFetch:  fetch,
	}
	return NewRouter(h), db
}

type apiKeyTestResponse struct {
	ID        uuid.UUID  `json:"id"`
	UserID    string     `json:"user_id"`
	Name      string     `json:"name"`
	KeyPrefix string     `json:"key_prefix"`
	APIKey    string     `json:"api_key"`
	ExpiresAt *time.Time `json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type apiKeyListTestResponse struct {
	APIKeys []apiKeyTestResponse `json:"api_keys"`
}

type deleteAPIKeyTestResponse struct {
	DeletedID    uuid.UUID `json:"deleted_id"`
	DeletedCount int       `json:"deleted_count"`
}

func authJSON(t *testing.T, router http.Handler, method string, path string, body string, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	req := newTestRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func decodeAuthResponse[T any](t *testing.T, w *httptest.ResponseRecorder) T {
	t.Helper()
	var out T
	if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode response: %v body=%s", err, w.Body.String())
	}
	return out
}

func TestAuthCookieName(t *testing.T) {
	if authSessionCookieName != "auth.session_token" {
		t.Fatalf("authSessionCookieName = %q", authSessionCookieName)
	}
}

func TestCreateAPIKeyStoresHashAndReturnsPlaintextOnce(t *testing.T) {
	router, db := testAuthRouter(t)
	user := createTestUser(t, db, "api-key-create@example.com")
	expiresAt := "2026-12-31T00:00:00Z"

	w := authJSON(t, router, http.MethodPost, "/api/auth/apikeys", `{
		"user_id":"`+user.ID+`",
		"name":" CLI key ",
		"expires_at":"`+expiresAt+`"
	}`, nil)

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAuthResponse[apiKeyTestResponse](t, w)
	if got.ID == uuid.Nil || got.UserID != user.ID || got.Name != "CLI key" {
		t.Fatalf("unexpected api key response: %#v", got)
	}
	if len(got.APIKey) != 64 {
		t.Fatalf("api_key length = %d, want 64", len(got.APIKey))
	}
	if got.KeyPrefix != got.APIKey[:8] {
		t.Fatalf("key_prefix = %q, want first 8 chars of api_key %q", got.KeyPrefix, got.APIKey[:8])
	}
	if got.ExpiresAt == nil || got.ExpiresAt.UTC().Format(time.RFC3339) != expiresAt {
		t.Fatalf("expires_at = %v, want %s", got.ExpiresAt, expiresAt)
	}

	var stored auth.APIKey
	if err := db.First(&stored, "id = ?", got.ID).Error; err != nil {
		t.Fatalf("load api key: %v", err)
	}
	if stored.UserID != user.ID || stored.Name != "CLI key" || stored.KeyPrefix != got.KeyPrefix {
		t.Fatalf("unexpected stored api key: %#v", stored)
	}
	if stored.KeyHash != auth.HashAPIKey(got.APIKey) {
		t.Fatalf("stored key hash = %q, want hash of returned api key", stored.KeyHash)
	}
	if stored.KeyHash == got.APIKey {
		t.Fatal("api key stored in plaintext")
	}
}

func TestCreateAPIKeyAllowsNullExpiration(t *testing.T) {
	router, db := testAuthRouter(t)
	user := createTestUser(t, db, "api-key-null-expiration@example.com")

	w := authJSON(t, router, http.MethodPost, "/api/auth/apikeys", `{
		"user_id":"`+user.ID+`",
		"name":"No expiry",
		"expires_at":null
	}`, nil)

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAuthResponse[apiKeyTestResponse](t, w)
	if got.ExpiresAt != nil {
		t.Fatalf("response expires_at = %v, want nil", got.ExpiresAt)
	}
	var stored auth.APIKey
	if err := db.First(&stored, "id = ?", got.ID).Error; err != nil {
		t.Fatalf("load api key: %v", err)
	}
	if stored.ExpiresAt != nil {
		t.Fatalf("stored expires_at = %v, want nil", stored.ExpiresAt)
	}
}

func TestListAPIKeysScopesToUserAndHidesPlaintextKeys(t *testing.T) {
	router, db := testAuthRouter(t)
	user := createTestUser(t, db, "api-key-list@example.com")
	other := createTestUser(t, db, "api-key-list-other@example.com")
	base := time.Date(2026, 6, 9, 10, 0, 0, 0, time.UTC)
	older := auth.APIKey{
		UserID:    user.ID,
		Name:      "older",
		KeyHash:   auth.HashAPIKey("older-api-key"),
		KeyPrefix: "older-ap",
		CreatedAt: base,
		UpdatedAt: base,
	}
	newer := auth.APIKey{
		UserID:    user.ID,
		Name:      "newer",
		KeyHash:   auth.HashAPIKey("newer-api-key"),
		KeyPrefix: "newer-ap",
		CreatedAt: base.Add(time.Minute),
		UpdatedAt: base.Add(time.Minute),
	}
	otherKey := auth.APIKey{
		UserID:    other.ID,
		Name:      "other",
		KeyHash:   auth.HashAPIKey("other-api-key"),
		KeyPrefix: "other-ap",
		CreatedAt: base.Add(2 * time.Minute),
		UpdatedAt: base.Add(2 * time.Minute),
	}
	for _, key := range []*auth.APIKey{&older, &newer, &otherKey} {
		if err := db.Create(key).Error; err != nil {
			t.Fatalf("create api key %s: %v", key.Name, err)
		}
	}

	w := authJSON(t, router, http.MethodGet, "/api/auth/apikeys/"+user.ID, "", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAuthResponse[apiKeyListTestResponse](t, w)
	if len(got.APIKeys) != 2 {
		t.Fatalf("api key count = %d, want 2: %#v", len(got.APIKeys), got.APIKeys)
	}
	if got.APIKeys[0].ID != newer.ID || got.APIKeys[1].ID != older.ID {
		t.Fatalf("api keys order = %#v, want newer then older", got.APIKeys)
	}
	for _, item := range got.APIKeys {
		if item.APIKey != "" {
			t.Fatalf("list response exposed plaintext api key: %#v", item)
		}
		if item.UserID != user.ID {
			t.Fatalf("list response included other user key: %#v", item)
		}
	}
	var raw map[string][]map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &raw); err != nil {
		t.Fatalf("decode raw list response: %v", err)
	}
	for _, item := range raw["api_keys"] {
		if _, ok := item["api_key"]; ok {
			t.Fatalf("list item includes api_key field: %#v", item)
		}
	}
}

func TestDeleteAPIKeyScopesToUserAndKeyID(t *testing.T) {
	router, db := testAuthRouter(t)
	user := createTestUser(t, db, "api-key-delete@example.com")
	other := createTestUser(t, db, "api-key-delete-other@example.com")
	key := auth.APIKey{UserID: user.ID, Name: "delete me", KeyHash: auth.HashAPIKey("delete-api-key"), KeyPrefix: "delete-a"}
	otherKey := auth.APIKey{UserID: other.ID, Name: "keep me", KeyHash: auth.HashAPIKey("keep-api-key"), KeyPrefix: "keep-api"}
	if err := db.Create(&key).Error; err != nil {
		t.Fatalf("create api key: %v", err)
	}
	if err := db.Create(&otherKey).Error; err != nil {
		t.Fatalf("create other api key: %v", err)
	}

	wrongUser := authJSON(t, router, http.MethodDelete, "/api/auth/apikeys?user_id="+other.ID+"&api_key_id="+key.ID.String(), "", nil)
	if wrongUser.Code != http.StatusNotFound {
		t.Fatalf("wrong-user status = %d body=%s", wrongUser.Code, wrongUser.Body.String())
	}

	w := authJSON(t, router, http.MethodDelete, "/api/auth/apikeys?user_id="+user.ID+"&api_key_id="+key.ID.String(), "", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAuthResponse[deleteAPIKeyTestResponse](t, w)
	if got.DeletedID != key.ID || got.DeletedCount != 1 {
		t.Fatalf("delete response = %#v, want key id and count 1", got)
	}
	var deletedCount int64
	if err := db.Model(&auth.APIKey{}).Where("id = ?", key.ID).Count(&deletedCount).Error; err != nil {
		t.Fatalf("count deleted key: %v", err)
	}
	if deletedCount != 0 {
		t.Fatalf("deleted key count = %d, want 0", deletedCount)
	}
	var otherCount int64
	if err := db.Model(&auth.APIKey{}).Where("id = ?", otherKey.ID).Count(&otherCount).Error; err != nil {
		t.Fatalf("count other key: %v", err)
	}
	if otherCount != 1 {
		t.Fatalf("other key count = %d, want 1", otherCount)
	}
}

func TestCreateAPIKeyRetriesGeneratedKeyCollision(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := apiPostgresTx(t)
	user := createTestUser(t, db, "api-key-collision@example.com")
	collidingKey := strings.Repeat("a", 64)
	nextKey := strings.Repeat("b", 64)
	existing := auth.APIKey{
		UserID:    user.ID,
		Name:      "existing",
		KeyHash:   auth.HashAPIKey(collidingKey),
		KeyPrefix: collidingKey[:8],
	}
	if err := db.Create(&existing).Error; err != nil {
		t.Fatalf("create existing api key: %v", err)
	}
	generated := []string{collidingKey, nextKey}
	h := &Handler{
		DB:               db,
		InternalAPIToken: testInternalAPIToken,
		APIKeyGenerator: func() (string, error) {
			if len(generated) == 0 {
				t.Fatal("api key generator called too many times")
			}
			key := generated[0]
			generated = generated[1:]
			return key, nil
		},
	}
	router := NewRouter(h)

	w := authJSON(t, router, http.MethodPost, "/api/auth/apikeys", `{
		"user_id":"`+user.ID+`",
		"name":"retry"
	}`, nil)

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAuthResponse[apiKeyTestResponse](t, w)
	if got.APIKey != nextKey || got.KeyPrefix != nextKey[:8] {
		t.Fatalf("returned key = %#v, want second generated key", got)
	}
}

func TestAPIKeyValidation(t *testing.T) {
	router, db := testAuthRouter(t)
	user := createTestUser(t, db, "api-key-validation@example.com")

	tests := []struct {
		name   string
		method string
		path   string
		body   string
		want   string
	}{
		{name: "create missing user", method: http.MethodPost, path: "/api/auth/apikeys", body: `{"name":"key"}`, want: "user_id is required"},
		{name: "create unknown user", method: http.MethodPost, path: "/api/auth/apikeys", body: `{"user_id":"` + uuid.NewString() + `","name":"key"}`, want: "invalid user_id"},
		{name: "create empty name", method: http.MethodPost, path: "/api/auth/apikeys", body: `{"user_id":"` + user.ID + `","name":"   "}`, want: "name is required"},
		{name: "create long name", method: http.MethodPost, path: "/api/auth/apikeys", body: `{"user_id":"` + user.ID + `","name":"` + strings.Repeat("a", 256) + `"}`, want: "name must be at most 255 characters"},
		{name: "create invalid expiration", method: http.MethodPost, path: "/api/auth/apikeys", body: `{"user_id":"` + user.ID + `","name":"key","expires_at":"tomorrow"}`, want: "expires_at must be RFC3339"},
		{name: "list invalid user", method: http.MethodGet, path: "/api/auth/apikeys/not-a-uuid", want: "invalid user_id"},
		{name: "list unknown user", method: http.MethodGet, path: "/api/auth/apikeys/" + uuid.NewString(), want: "invalid user_id"},
		{name: "delete missing key id", method: http.MethodDelete, path: "/api/auth/apikeys?user_id=" + user.ID, want: "api_key_id is required"},
		{name: "delete invalid key id", method: http.MethodDelete, path: "/api/auth/apikeys?user_id=" + user.ID + "&api_key_id=not-a-uuid", want: "invalid api_key_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := authJSON(t, router, tt.method, tt.path, tt.body, nil)
			if w.Code != http.StatusBadRequest {
				t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
			}
			got := decodeAuthResponse[ErrorResponse](t, w)
			if got.Error != tt.want {
				t.Fatalf("error = %q, want %q", got.Error, tt.want)
			}
		})
	}
}

func TestSetSessionCookieRememberTrueSetsPersistentSecureCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	h := &Handler{AuthCookieSecure: true}
	expiresAt := time.Now().UTC().Add(time.Hour)

	h.setSessionCookie(c, "session-token", expiresAt, true)

	cookie := authCookieFromRecorder(t, w)
	if cookie.Value != "session-token" {
		t.Fatalf("cookie value = %q", cookie.Value)
	}
	if !cookie.HttpOnly {
		t.Fatal("cookie HttpOnly = false")
	}
	if cookie.SameSite != http.SameSiteLaxMode {
		t.Fatalf("cookie SameSite = %v", cookie.SameSite)
	}
	if cookie.Path != "/" {
		t.Fatalf("cookie Path = %q", cookie.Path)
	}
	if !cookie.Secure {
		t.Fatal("cookie Secure = false")
	}
	if cookie.Expires.IsZero() {
		t.Fatal("cookie Expires is zero")
	}
	if cookie.MaxAge <= 0 {
		t.Fatalf("cookie MaxAge = %d", cookie.MaxAge)
	}
}

func TestSetSessionCookieRememberFalseSetsBrowserSessionCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	h := &Handler{AuthCookieSecure: false}

	h.setSessionCookie(c, "session-token", time.Now().UTC().Add(time.Hour), false)

	cookie := authCookieFromRecorder(t, w)
	if !cookie.HttpOnly {
		t.Fatal("cookie HttpOnly = false")
	}
	if cookie.SameSite != http.SameSiteLaxMode {
		t.Fatalf("cookie SameSite = %v", cookie.SameSite)
	}
	if cookie.Path != "/" {
		t.Fatalf("cookie Path = %q", cookie.Path)
	}
	if cookie.Secure {
		t.Fatal("cookie Secure = true")
	}
	if !cookie.Expires.IsZero() {
		t.Fatalf("cookie Expires = %v, want zero", cookie.Expires)
	}
	if cookie.MaxAge != 0 {
		t.Fatalf("cookie MaxAge = %d, want 0", cookie.MaxAge)
	}
	setCookie := w.Result().Header.Get("Set-Cookie")
	if strings.Contains(strings.ToLower(setCookie), "expires=") || strings.Contains(strings.ToLower(setCookie), "max-age=") {
		t.Fatalf("session cookie should omit persistent attributes: %s", setCookie)
	}
}

func TestClearSessionCookieExpiresImmediately(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	h := &Handler{AuthCookieSecure: true}

	h.clearSessionCookie(c)

	cookie := authCookieFromRecorder(t, w)
	if cookie.Value != "" {
		t.Fatalf("cookie value = %q", cookie.Value)
	}
	if cookie.MaxAge != -1 {
		t.Fatalf("cookie MaxAge = %d, want -1", cookie.MaxAge)
	}
	if !cookie.Expires.Equal(time.Unix(0, 0)) {
		t.Fatalf("cookie Expires = %v, want Unix epoch", cookie.Expires)
	}
	if !cookie.HttpOnly {
		t.Fatal("cookie HttpOnly = false")
	}
	if cookie.SameSite != http.SameSiteLaxMode {
		t.Fatalf("cookie SameSite = %v", cookie.SameSite)
	}
	if cookie.Path != "/" {
		t.Fatalf("cookie Path = %q", cookie.Path)
	}
	if !cookie.Secure {
		t.Fatal("cookie Secure = false")
	}
}

func TestSignUpEmailCreatesUnverifiedUserCredentialAccountAndVerification(t *testing.T) {
	router, db := testAuthRouter(t)

	w := authJSON(t, router, http.MethodPost, "/api/auth/sign-up/email", `{
		"name":" Ada ",
		"email":"ADA@EXAMPLE.COM",
		"password":"password1234"
	}`, map[string]string{authDeliveryHeader: testDeliveryToken})
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}

	var got struct {
		User                  authUserResponse `json:"user"`
		VerificationRequired  bool             `json:"verificationRequired"`
		VerificationCode      string           `json:"verificationCode"`
		VerificationExpiresAt string           `json:"verificationExpiresAt"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.User.Email != "ada@example.com" || got.User.Name != "Ada" || got.User.EmailVerified {
		t.Fatalf("unexpected user response: %#v", got.User)
	}
	if !got.VerificationRequired || got.VerificationCode == "" || got.VerificationExpiresAt == "" {
		t.Fatalf("unexpected verification response: %#v", got)
	}

	var user auth.User
	if err := db.First(&user, "email = ?", "ada@example.com").Error; err != nil {
		t.Fatalf("load user: %v", err)
	}
	if user.EmailVerified {
		t.Fatal("new user should be unverified")
	}

	var account auth.AuthAccount
	if err := db.First(&account, "user_id = ? AND provider_id = ?", user.ID, auth.CredentialProviderID).Error; err != nil {
		t.Fatalf("load account: %v", err)
	}
	if account.AccountID != user.ID || account.Password == "" || account.Password == "password1234" {
		t.Fatalf("unsafe account: %#v", account)
	}

	var verification auth.Verification
	if err := db.First(&verification, "identifier = ?", verificationIdentifier("ada@example.com")).Error; err != nil {
		t.Fatalf("load verification: %v", err)
	}
	if verification.Value == got.VerificationCode {
		t.Fatal("verification code stored in plaintext")
	}
}

func TestSignUpEmailDoesNotLeakCodeWithoutTrustedHeader(t *testing.T) {
	router, _ := testAuthRouter(t)

	w := authJSON(t, router, http.MethodPost, "/api/auth/sign-up/email", `{
		"name":"Ada",
		"email":"ada@example.com",
		"password":"password1234"
	}`, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if bytes.Contains(w.Body.Bytes(), []byte("verificationCode")) {
		t.Fatalf("untrusted response leaked verification code: %s", w.Body.String())
	}
}

func TestSignUpEmailRotatesVerificationForExistingUnverifiedUser(t *testing.T) {
	router, db := testAuthRouter(t)

	first := authJSON(t, router, http.MethodPost, "/api/auth/sign-up/email", `{
		"name":"Ada",
		"email":"ada@example.com",
		"password":"password1234"
	}`, map[string]string{authDeliveryHeader: testDeliveryToken})
	if first.Code != http.StatusOK {
		t.Fatalf("first status = %d body=%s", first.Code, first.Body.String())
	}
	firstResponse := decodeAuthResponse[struct {
		VerificationCode string `json:"verificationCode"`
	}](t, first)

	var secondResponse struct {
		User             authUserResponse `json:"user"`
		VerificationCode string           `json:"verificationCode"`
	}
	for i := 0; i < 5; i++ {
		second := authJSON(t, router, http.MethodPost, "/api/auth/sign-up/email", `{
			"name":"Ada Lovelace",
			"email":"ADA@EXAMPLE.COM",
			"password":"password1234"
		}`, map[string]string{authDeliveryHeader: testDeliveryToken})
		if second.Code != http.StatusOK {
			t.Fatalf("second status = %d body=%s", second.Code, second.Body.String())
		}
		secondResponse = decodeAuthResponse[struct {
			User             authUserResponse `json:"user"`
			VerificationCode string           `json:"verificationCode"`
		}](t, second)
		if secondResponse.VerificationCode != firstResponse.VerificationCode {
			break
		}
	}
	if secondResponse.User.Name != "Ada Lovelace" || secondResponse.User.ID != "" || secondResponse.User.EmailVerified {
		t.Fatalf("duplicate signup should return generic unverified user: %#v", secondResponse.User)
	}
	if secondResponse.VerificationCode == "" || secondResponse.VerificationCode == firstResponse.VerificationCode {
		t.Fatalf("verification code was not rotated: first=%q second=%q", firstResponse.VerificationCode, secondResponse.VerificationCode)
	}

	identifier := verificationIdentifier("ada@example.com")
	var count int64
	if err := db.Model(&auth.Verification{}).Where("identifier = ?", identifier).Count(&count).Error; err != nil {
		t.Fatalf("count verifications: %v", err)
	}
	if count != 1 {
		t.Fatalf("verification count = %d, want 1", count)
	}
	var verification auth.Verification
	if err := db.First(&verification, "identifier = ?", identifier).Error; err != nil {
		t.Fatalf("load verification: %v", err)
	}
	if authcrypto.VerifyCode(testAuthSecret, identifier, firstResponse.VerificationCode, verification.Value) {
		t.Fatal("old verification code still verifies after rotation")
	}
	if !authcrypto.VerifyCode(testAuthSecret, identifier, secondResponse.VerificationCode, verification.Value) {
		t.Fatal("new verification code does not verify")
	}
}

func TestSignUpEmailVerifiedDuplicateReturnsGenericUserWithoutCreatingRowsOrCode(t *testing.T) {
	router, db := testAuthRouter(t)
	image := "https://example.com/real.png"
	stored := auth.User{
		Name:          "Stored Real Name",
		Email:         "ada@example.com",
		EmailVerified: true,
		Image:         &image,
	}
	if err := db.Create(&stored).Error; err != nil {
		t.Fatalf("create stored user: %v", err)
	}
	account := auth.AuthAccount{
		AccountID:  stored.ID,
		ProviderID: auth.CredentialProviderID,
		UserID:     stored.ID,
		Password:   "scrypt$v=1$stored",
	}
	if err := db.Create(&account).Error; err != nil {
		t.Fatalf("create stored account: %v", err)
	}

	w := authJSON(t, router, http.MethodPost, "/api/auth/sign-up/email", `{
		"name":" Supplied Name ",
		"email":"ADA@EXAMPLE.COM",
		"password":"password1234",
		"image":"https://example.com/supplied.png"
	}`, map[string]string{authDeliveryHeader: testDeliveryToken})
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAuthResponse[signUpEmailResponse](t, w)
	if !got.VerificationRequired {
		t.Fatalf("verificationRequired = false, want true: %#v", got)
	}
	if got.VerificationCode != "" || bytes.Contains(w.Body.Bytes(), []byte("verificationCode")) {
		t.Fatalf("verified duplicate response leaked verification code: %s", w.Body.String())
	}
	if got.User.ID == stored.ID || got.User.Name == stored.Name || got.User.Image != nil || got.User.EmailVerified {
		t.Fatalf("verified duplicate response leaked stored user metadata: %#v", got.User)
	}
	if got.User.ID != "" || got.User.Name != "Supplied Name" || got.User.Email != "ada@example.com" {
		t.Fatalf("verified duplicate response should use generic request user: %#v", got.User)
	}

	var userCount int64
	if err := db.Model(&auth.User{}).Where("email = ?", "ada@example.com").Count(&userCount).Error; err != nil {
		t.Fatalf("count users: %v", err)
	}
	if userCount != 1 {
		t.Fatalf("user count = %d, want 1", userCount)
	}
	var accountCount int64
	if err := db.Model(&auth.AuthAccount{}).Where("user_id = ?", stored.ID).Count(&accountCount).Error; err != nil {
		t.Fatalf("count accounts: %v", err)
	}
	if accountCount != 1 {
		t.Fatalf("account count = %d, want 1", accountCount)
	}
	var verificationCount int64
	if err := db.Model(&auth.Verification{}).Where("identifier = ?", verificationIdentifier("ada@example.com")).Count(&verificationCount).Error; err != nil {
		t.Fatalf("count verifications: %v", err)
	}
	if verificationCount != 0 {
		t.Fatalf("verification count = %d, want 0", verificationCount)
	}
}

func TestSignUpEmailRejectsInvalidInput(t *testing.T) {
	router, _ := testAuthRouter(t)
	overlongEmail := strings.Repeat("a", 309) + "@example.com"

	cases := []string{
		`{"name":"","email":"ada@example.com","password":"password1234"}`,
		`{"name":"Ada","email":"not-email","password":"password1234"}`,
		`{"name":"Ada","email":"ada@example.com","password":"short"}`,
		`{"name":"Ada","email":"` + overlongEmail + `","password":"password1234"}`,
	}
	for _, body := range cases {
		w := authJSON(t, router, http.MethodPost, "/api/auth/sign-up/email", body, nil)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("body %s status = %d want 400 response=%s", body, w.Code, w.Body.String())
		}
	}
}

func TestSignInEmailRejectsUnverifiedUser(t *testing.T) {
	router, _ := testAuthRouter(t)
	signupVerificationCode(t, router, "ADA@EXAMPLE.COM")

	w := authJSON(t, router, http.MethodPost, "/api/auth/sign-in/email", `{
		"email":"ada@example.com",
		"password":"password1234"
	}`, nil)
	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d want 403 response=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "email is not verified") {
		t.Fatalf("expected unverified error, got %s", w.Body.String())
	}
}

func TestSignInEmailRejectsUnverifiedUserWithWrongPasswordAsInvalidCredentials(t *testing.T) {
	router, _ := testAuthRouter(t)
	signupVerificationCode(t, router, "ADA@EXAMPLE.COM")

	w := authJSON(t, router, http.MethodPost, "/api/auth/sign-in/email", `{
		"email":"ada@example.com",
		"password":"wrong-password"
	}`, nil)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d want 401 response=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "invalid email or password") {
		t.Fatalf("expected generic credential error, got %s", w.Body.String())
	}
}

func TestSignInEmailCreatesSessionAndCookie(t *testing.T) {
	router, db := testAuthRouter(t)
	code := signupVerificationCode(t, router, "ADA@EXAMPLE.COM")
	verifySignupEmail(t, router, "ada@example.com", code)

	w := authJSON(t, router, http.MethodPost, "/api/auth/sign-in/email", `{
		"email":" ADA@EXAMPLE.COM ",
		"password":"password1234",
		"rememberMe":true
	}`, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAuthResponse[authSessionPayload](t, w)
	if got.Session.Token == "" {
		t.Fatalf("session token is empty: %#v", got.Session)
	}
	if got.User.Email != "ada@example.com" || !got.User.EmailVerified {
		t.Fatalf("unexpected user response: %#v", got.User)
	}

	cookie := authCookieFromRecorder(t, w)
	if cookie.Value != got.Session.Token {
		t.Fatalf("cookie value = %q, want session token %q", cookie.Value, got.Session.Token)
	}
	if !cookie.HttpOnly {
		t.Fatal("cookie HttpOnly = false")
	}

	var user auth.User
	if err := db.First(&user, "email = ?", "ada@example.com").Error; err != nil {
		t.Fatalf("load user: %v", err)
	}
	var session auth.Session
	if err := db.First(&session, "token = ? AND user_id = ?", got.Session.Token, user.ID).Error; err != nil {
		t.Fatalf("load session: %v", err)
	}
	if session.Token == "" || session.UserID != user.ID {
		t.Fatalf("unexpected stored session: %#v", session)
	}
}

func TestSignInEmailUpdatesLastLoginAt(t *testing.T) {
	router, db := testAuthRouter(t)
	code := signupVerificationCode(t, router, "LASTLOGIN@EXAMPLE.COM")
	verifySignupEmail(t, router, "lastlogin@example.com", code)

	var before sql.NullTime
	if err := db.Raw(`SELECT last_login_at FROM "user" WHERE email = ?`, "lastlogin@example.com").Scan(&before).Error; err != nil {
		t.Fatalf("load last_login_at before sign-in: %v", err)
	}
	if before.Valid {
		t.Fatalf("last_login_at before sign-in = %v, want null", before.Time)
	}

	w := authJSON(t, router, http.MethodPost, "/api/auth/sign-in/email", `{
		"email":"lastlogin@example.com",
		"password":"password1234"
	}`, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}

	var after sql.NullTime
	if err := db.Raw(`SELECT last_login_at FROM "user" WHERE email = ?`, "lastlogin@example.com").Scan(&after).Error; err != nil {
		t.Fatalf("load last_login_at after sign-in: %v", err)
	}
	if !after.Valid || after.Time.IsZero() {
		t.Fatalf("last_login_at after sign-in = %#v, want non-zero timestamp", after)
	}
}

func TestStartGoogleOAuthCreatesStateAndAuthorizationURL(t *testing.T) {
	redirectURI := "http://localhost:5173/api/auth/callback/google"
	router, db := testGoogleAuthRouter(t, nil, nil)

	w := authJSON(t, router, http.MethodPost, "/api/auth/oauth/google/start", `{
		"redirectURI":"`+redirectURI+`"
	}`, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAuthResponse[googleOAuthStartResponse](t, w)
	if got.State == "" || got.StateExpiresAt == "" {
		t.Fatalf("missing state fields: %#v", got)
	}
	authorizationURL, err := url.Parse(got.AuthorizationURL)
	if err != nil {
		t.Fatalf("parse authorization URL: %v", err)
	}
	query := authorizationURL.Query()
	if query.Get("client_id") != "google-client-id" {
		t.Fatalf("client_id = %q", query.Get("client_id"))
	}
	if query.Get("redirect_uri") != redirectURI {
		t.Fatalf("redirect_uri = %q", query.Get("redirect_uri"))
	}
	if query.Get("state") != got.State {
		t.Fatalf("state in URL = %q, want %q", query.Get("state"), got.State)
	}
	if query.Get("response_type") != "code" {
		t.Fatalf("response_type = %q", query.Get("response_type"))
	}
	scope := " " + query.Get("scope") + " "
	for _, want := range []string{"openid", "email", "profile"} {
		if !strings.Contains(scope, " "+want+" ") {
			t.Fatalf("scope %q missing %q", query.Get("scope"), want)
		}
	}

	var verification auth.Verification
	if err := db.First(&verification, "identifier = ?", googleOAuthStateIdentifier(got.State)).Error; err != nil {
		t.Fatalf("load oauth state verification: %v", err)
	}
	if !auth.VerifyCode(testAuthSecret, verification.Identifier, redirectURI, verification.Value) {
		t.Fatal("stored oauth state hash does not verify redirect URI")
	}
}

func TestStartGoogleOAuthRequiresConfiguration(t *testing.T) {
	router, _ := testAuthRouter(t)

	w := authJSON(t, router, http.MethodPost, "/api/auth/oauth/google/start", `{
		"redirectURI":"http://localhost:5173/api/auth/callback/google"
	}`, nil)
	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d body=%s, want 503", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "google oauth is not configured") {
		t.Fatalf("unexpected response body: %s", w.Body.String())
	}
}

func TestSignInGoogleOAuthCreatesUserAccountSessionAndSignupBonus(t *testing.T) {
	redirectURI := "http://localhost:5173/api/auth/callback/google"
	accessToken := "access-token"
	refreshToken := "refresh-token"
	scope := "openid email profile"
	expiresAt := time.Now().UTC().Add(time.Hour)
	router, db := testGoogleAuthRouter(t,
		func(_ context.Context, code string, gotRedirectURI string) (googleOAuthToken, error) {
			if code != "code-1" {
				t.Fatalf("code = %q, want code-1", code)
			}
			if gotRedirectURI != redirectURI {
				t.Fatalf("redirectURI = %q, want %q", gotRedirectURI, redirectURI)
			}
			return googleOAuthToken{
				AccessToken:          &accessToken,
				RefreshToken:         &refreshToken,
				IDToken:              "id-token-1",
				AccessTokenExpiresAt: &expiresAt,
				Scope:                &scope,
			}, nil
		},
		func(_ context.Context, idToken string, audience string) (googleIDTokenPayload, error) {
			if idToken != "id-token-1" {
				t.Fatalf("idToken = %q, want id-token-1", idToken)
			}
			if audience != "google-client-id" {
				t.Fatalf("audience = %q, want google-client-id", audience)
			}
			return googleIDTokenPayload{
				Subject:       "google-sub-1",
				Email:         "ADA@EXAMPLE.COM",
				EmailVerified: true,
				Name:          "Ada Google",
				Picture:       "https://lh3.googleusercontent.com/avatar",
			}, nil
		},
	)
	state := startGoogleOAuthState(t, router, redirectURI)

	w := authJSON(t, router, http.MethodPost, "/api/auth/oauth/google/callback", `{
		"code":"code-1",
		"state":"`+state+`",
		"redirectURI":"`+redirectURI+`"
	}`, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAuthResponse[authSessionPayload](t, w)
	if got.User.Email != "ada@example.com" || got.User.Name != "Ada Google" || !got.User.EmailVerified {
		t.Fatalf("unexpected user: %#v", got.User)
	}
	if got.Session.Token == "" || got.Session.UserID != got.User.ID {
		t.Fatalf("unexpected session: %#v", got.Session)
	}
	cookie := authCookieFromRecorder(t, w)
	if cookie.Value != got.Session.Token {
		t.Fatalf("cookie value = %q, want session token %q", cookie.Value, got.Session.Token)
	}

	var account auth.AuthAccount
	if err := db.First(&account, "provider_id = ? AND account_id = ?", auth.GoogleProviderID, "google-sub-1").Error; err != nil {
		t.Fatalf("load google account: %v", err)
	}
	if account.UserID != got.User.ID || account.AccessToken == nil || *account.AccessToken != accessToken || account.RefreshToken == nil || *account.RefreshToken != refreshToken || account.Scope == nil || *account.Scope != scope {
		t.Fatalf("unexpected google account: %#v", account)
	}

	var session auth.Session
	if err := db.First(&session, "token = ? AND user_id = ?", got.Session.Token, got.User.ID).Error; err != nil {
		t.Fatalf("load stored session: %v", err)
	}
	bucket, err := billing.BucketForLedgerIdempotencyKey(t.Context(), db, "signup_bonus:"+got.User.ID)
	if err != nil {
		t.Fatalf("load signup bonus bucket: %v", err)
	}
	if bucket == nil {
		t.Fatal("missing signup bonus bucket")
	}
	if bucket.CreditsGranted != testOnboardingCredits || bucket.CreditsRemaining != testOnboardingCredits {
		t.Fatalf("signup bonus credits = granted %d remaining %d, want %d/%d", bucket.CreditsGranted, bucket.CreditsRemaining, testOnboardingCredits, testOnboardingCredits)
	}
	assertGoogleOAuthStateConsumed(t, db, state)
}

func TestSignInGoogleOAuthLinksExistingUnverifiedUserByEmail(t *testing.T) {
	redirectURI := "http://localhost:5173/api/auth/callback/google"
	router, db := testGoogleAuthRouter(t, successfulGoogleExchange(t), googleProfileValidator(googleIDTokenPayload{
		Subject:       "google-sub-link",
		Email:         "link@example.com",
		EmailVerified: true,
		Name:          "Linked Google Name",
	}))
	user := createTestUser(t, db, "link@example.com")
	state := startGoogleOAuthState(t, router, redirectURI)

	w := authJSON(t, router, http.MethodPost, "/api/auth/oauth/google/callback", `{
		"code":"code-1",
		"state":"`+state+`",
		"redirectURI":"`+redirectURI+`"
	}`, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAuthResponse[authSessionPayload](t, w)
	if got.User.ID != user.ID || got.User.Name != user.Name || !got.User.EmailVerified {
		t.Fatalf("unexpected linked user: %#v want id=%s name=%q verified=true", got.User, user.ID, user.Name)
	}
	var account auth.AuthAccount
	if err := db.First(&account, "provider_id = ? AND account_id = ?", auth.GoogleProviderID, "google-sub-link").Error; err != nil {
		t.Fatalf("load linked google account: %v", err)
	}
	if account.UserID != user.ID {
		t.Fatalf("linked account user_id = %s, want %s", account.UserID, user.ID)
	}
	bucket, err := billing.BucketForLedgerIdempotencyKey(t.Context(), db, "signup_bonus:"+user.ID)
	if err != nil {
		t.Fatalf("load signup bonus bucket: %v", err)
	}
	if bucket == nil {
		t.Fatal("missing signup bonus bucket for newly verified linked user")
	}
	if bucket.CreditsGranted != testOnboardingCredits || bucket.CreditsRemaining != testOnboardingCredits {
		t.Fatalf("signup bonus credits = granted %d remaining %d, want %d/%d", bucket.CreditsGranted, bucket.CreditsRemaining, testOnboardingCredits, testOnboardingCredits)
	}
}

func TestSignInGoogleOAuthUsesExistingLinkedAccountAndPreservesRefreshToken(t *testing.T) {
	redirectURI := "http://localhost:5173/api/auth/callback/google"
	newAccessToken := "new-access-token"
	oldRefreshToken := "old-refresh-token"
	router, db := testGoogleAuthRouter(t,
		func(context.Context, string, string) (googleOAuthToken, error) {
			return googleOAuthToken{AccessToken: &newAccessToken, IDToken: "id-token-1"}, nil
		},
		googleProfileValidator(googleIDTokenPayload{
			Subject:       "google-sub-existing",
			Email:         "changed-google-email@example.com",
			EmailVerified: true,
			Name:          "Changed Google Name",
		}),
	)
	user := auth.User{Name: "Existing User", Email: "existing@example.com", EmailVerified: true}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create existing user: %v", err)
	}
	account := auth.AuthAccount{
		AccountID:    "google-sub-existing",
		ProviderID:   auth.GoogleProviderID,
		UserID:       user.ID,
		RefreshToken: &oldRefreshToken,
	}
	if err := db.Create(&account).Error; err != nil {
		t.Fatalf("create existing google account: %v", err)
	}
	state := startGoogleOAuthState(t, router, redirectURI)

	w := authJSON(t, router, http.MethodPost, "/api/auth/oauth/google/callback", `{
		"code":"code-1",
		"state":"`+state+`",
		"redirectURI":"`+redirectURI+`"
	}`, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAuthResponse[authSessionPayload](t, w)
	if got.User.ID != user.ID || got.User.Email != "existing@example.com" || got.User.Name != "Existing User" {
		t.Fatalf("unexpected existing linked user response: %#v", got.User)
	}

	var updated auth.AuthAccount
	if err := db.First(&updated, "id = ?", account.ID).Error; err != nil {
		t.Fatalf("load updated google account: %v", err)
	}
	if updated.AccessToken == nil || *updated.AccessToken != newAccessToken {
		t.Fatalf("access token = %v, want %q", updated.AccessToken, newAccessToken)
	}
	if updated.RefreshToken == nil || *updated.RefreshToken != oldRefreshToken {
		t.Fatalf("refresh token = %v, want preserved %q", updated.RefreshToken, oldRefreshToken)
	}
}

func TestSignInGoogleOAuthRejectsInvalidStateWithoutExchange(t *testing.T) {
	redirectURI := "http://localhost:5173/api/auth/callback/google"
	exchanged := false
	router, _ := testGoogleAuthRouter(t,
		func(context.Context, string, string) (googleOAuthToken, error) {
			exchanged = true
			return googleOAuthToken{}, nil
		},
		googleProfileValidator(googleIDTokenPayload{
			Subject:       "google-sub-1",
			Email:         "ada@example.com",
			EmailVerified: true,
		}),
	)

	w := authJSON(t, router, http.MethodPost, "/api/auth/oauth/google/callback", `{
		"code":"code-1",
		"state":"missing-state",
		"redirectURI":"`+redirectURI+`"
	}`, nil)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s, want 400", w.Code, w.Body.String())
	}
	if exchanged {
		t.Fatal("oauth code exchange ran for invalid state")
	}
}

func TestSignInGoogleOAuthRejectsInvalidGoogleIdentityAndConsumesState(t *testing.T) {
	redirectURI := "http://localhost:5173/api/auth/callback/google"
	router, db := testGoogleAuthRouter(t, successfulGoogleExchange(t), googleProfileValidator(googleIDTokenPayload{
		Subject:       "google-sub-invalid",
		Email:         "invalid@example.com",
		EmailVerified: false,
	}))
	state := startGoogleOAuthState(t, router, redirectURI)

	w := authJSON(t, router, http.MethodPost, "/api/auth/oauth/google/callback", `{
		"code":"code-1",
		"state":"`+state+`",
		"redirectURI":"`+redirectURI+`"
	}`, nil)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d body=%s, want 401", w.Code, w.Body.String())
	}
	assertGoogleOAuthStateConsumed(t, db, state)
}

func TestStartGitHubOAuthCreatesStateAndAuthorizationURL(t *testing.T) {
	redirectURI := "http://localhost:5173/api/auth/callback/github"
	router, db := testGitHubAuthRouter(t, nil, nil)

	w := authJSON(t, router, http.MethodPost, "/api/auth/oauth/github/start", `{
		"redirectURI":"`+redirectURI+`"
	}`, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAuthResponse[githubOAuthStartResponse](t, w)
	if got.State == "" || got.StateExpiresAt == "" {
		t.Fatalf("missing state fields: %#v", got)
	}
	authorizationURL, err := url.Parse(got.AuthorizationURL)
	if err != nil {
		t.Fatalf("parse authorization URL: %v", err)
	}
	query := authorizationURL.Query()
	if authorizationURL.Host != "github.com" {
		t.Fatalf("authorization host = %q, want github.com", authorizationURL.Host)
	}
	if query.Get("client_id") != "github-client-id" {
		t.Fatalf("client_id = %q", query.Get("client_id"))
	}
	if query.Get("redirect_uri") != redirectURI {
		t.Fatalf("redirect_uri = %q", query.Get("redirect_uri"))
	}
	if query.Get("state") != got.State {
		t.Fatalf("state in URL = %q, want %q", query.Get("state"), got.State)
	}
	if query.Get("response_type") != "code" {
		t.Fatalf("response_type = %q", query.Get("response_type"))
	}
	scope := " " + query.Get("scope") + " "
	for _, want := range []string{"read:user", "user:email"} {
		if !strings.Contains(scope, " "+want+" ") {
			t.Fatalf("scope %q missing %q", query.Get("scope"), want)
		}
	}

	var verification auth.Verification
	if err := db.First(&verification, "identifier = ?", githubOAuthStateIdentifier(got.State)).Error; err != nil {
		t.Fatalf("load oauth state verification: %v", err)
	}
	if !auth.VerifyCode(testAuthSecret, verification.Identifier, redirectURI, verification.Value) {
		t.Fatal("stored oauth state hash does not verify redirect URI")
	}
}

func TestStartGitHubOAuthRequiresConfiguration(t *testing.T) {
	router, _ := testAuthRouter(t)

	w := authJSON(t, router, http.MethodPost, "/api/auth/oauth/github/start", `{
		"redirectURI":"http://localhost:5173/api/auth/callback/github"
	}`, nil)
	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d body=%s, want 503", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "github oauth is not configured") {
		t.Fatalf("unexpected response body: %s", w.Body.String())
	}
}

func TestSignInGitHubOAuthCreatesUserAccountSessionAndSignupBonus(t *testing.T) {
	redirectURI := "http://localhost:5173/api/auth/callback/github"
	accessToken := "github-access-token"
	refreshToken := "github-refresh-token"
	scope := "read:user user:email"
	expiresAt := time.Now().UTC().Add(time.Hour)
	router, db := testGitHubAuthRouter(t,
		func(_ context.Context, code string, gotRedirectURI string) (githubOAuthToken, error) {
			if code != "code-1" {
				t.Fatalf("code = %q, want code-1", code)
			}
			if gotRedirectURI != redirectURI {
				t.Fatalf("redirectURI = %q, want %q", gotRedirectURI, redirectURI)
			}
			return githubOAuthToken{
				AccessToken:          &accessToken,
				RefreshToken:         &refreshToken,
				AccessTokenExpiresAt: &expiresAt,
				Scope:                &scope,
			}, nil
		},
		func(_ context.Context, gotAccessToken string) (githubOAuthProfile, error) {
			if gotAccessToken != accessToken {
				t.Fatalf("accessToken = %q, want %q", gotAccessToken, accessToken)
			}
			return githubOAuthProfile{
				Subject:       "12345",
				Email:         "ADA@EXAMPLE.COM",
				EmailVerified: true,
				Name:          "Ada GitHub",
				Login:         "adagithub",
				Picture:       "https://avatars.githubusercontent.com/u/12345",
			}, nil
		},
	)
	state := startGitHubOAuthState(t, router, redirectURI)

	w := authJSON(t, router, http.MethodPost, "/api/auth/oauth/github/callback", `{
		"code":"code-1",
		"state":"`+state+`",
		"redirectURI":"`+redirectURI+`"
	}`, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAuthResponse[authSessionPayload](t, w)
	if got.User.Email != "ada@example.com" || got.User.Name != "Ada GitHub" || !got.User.EmailVerified {
		t.Fatalf("unexpected user: %#v", got.User)
	}
	if got.Session.Token == "" || got.Session.UserID != got.User.ID {
		t.Fatalf("unexpected session: %#v", got.Session)
	}
	cookie := authCookieFromRecorder(t, w)
	if cookie.Value != got.Session.Token {
		t.Fatalf("cookie value = %q, want session token %q", cookie.Value, got.Session.Token)
	}

	var account auth.AuthAccount
	if err := db.First(&account, "provider_id = ? AND account_id = ?", auth.GitHubProviderID, "12345").Error; err != nil {
		t.Fatalf("load github account: %v", err)
	}
	if account.UserID != got.User.ID || account.AccessToken == nil || *account.AccessToken != accessToken || account.RefreshToken == nil || *account.RefreshToken != refreshToken || account.Scope == nil || *account.Scope != scope {
		t.Fatalf("unexpected github account: %#v", account)
	}

	var session auth.Session
	if err := db.First(&session, "token = ? AND user_id = ?", got.Session.Token, got.User.ID).Error; err != nil {
		t.Fatalf("load stored session: %v", err)
	}
	bucket, err := billing.BucketForLedgerIdempotencyKey(t.Context(), db, "signup_bonus:"+got.User.ID)
	if err != nil {
		t.Fatalf("load signup bonus bucket: %v", err)
	}
	if bucket == nil {
		t.Fatal("missing signup bonus bucket")
	}
	if bucket.CreditsGranted != testOnboardingCredits || bucket.CreditsRemaining != testOnboardingCredits {
		t.Fatalf("signup bonus credits = granted %d remaining %d, want %d/%d", bucket.CreditsGranted, bucket.CreditsRemaining, testOnboardingCredits, testOnboardingCredits)
	}
	assertGitHubOAuthStateConsumed(t, db, state)
}

func TestSignInGitHubOAuthLinksExistingUnverifiedUserByEmail(t *testing.T) {
	redirectURI := "http://localhost:5173/api/auth/callback/github"
	router, db := testGitHubAuthRouter(t, successfulGitHubExchange(t), githubProfileFetcher(githubOAuthProfile{
		Subject:       "67890",
		Email:         "link-github@example.com",
		EmailVerified: true,
		Name:          "Linked GitHub Name",
		Login:         "linkedgithub",
	}))
	user := createTestUser(t, db, "link-github@example.com")
	state := startGitHubOAuthState(t, router, redirectURI)

	w := authJSON(t, router, http.MethodPost, "/api/auth/oauth/github/callback", `{
		"code":"code-1",
		"state":"`+state+`",
		"redirectURI":"`+redirectURI+`"
	}`, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAuthResponse[authSessionPayload](t, w)
	if got.User.ID != user.ID || got.User.Name != user.Name || !got.User.EmailVerified {
		t.Fatalf("unexpected linked user: %#v want id=%s name=%q verified=true", got.User, user.ID, user.Name)
	}
	var account auth.AuthAccount
	if err := db.First(&account, "provider_id = ? AND account_id = ?", auth.GitHubProviderID, "67890").Error; err != nil {
		t.Fatalf("load linked github account: %v", err)
	}
	if account.UserID != user.ID {
		t.Fatalf("linked account user_id = %s, want %s", account.UserID, user.ID)
	}
	bucket, err := billing.BucketForLedgerIdempotencyKey(t.Context(), db, "signup_bonus:"+user.ID)
	if err != nil {
		t.Fatalf("load signup bonus bucket: %v", err)
	}
	if bucket == nil {
		t.Fatal("missing signup bonus bucket for newly verified linked user")
	}
	if bucket.CreditsGranted != testOnboardingCredits || bucket.CreditsRemaining != testOnboardingCredits {
		t.Fatalf("signup bonus credits = granted %d remaining %d, want %d/%d", bucket.CreditsGranted, bucket.CreditsRemaining, testOnboardingCredits, testOnboardingCredits)
	}
}

func TestSignInGitHubOAuthUsesExistingLinkedAccountAndPreservesRefreshToken(t *testing.T) {
	redirectURI := "http://localhost:5173/api/auth/callback/github"
	newAccessToken := "new-github-access-token"
	oldRefreshToken := "old-github-refresh-token"
	router, db := testGitHubAuthRouter(t,
		func(context.Context, string, string) (githubOAuthToken, error) {
			return githubOAuthToken{AccessToken: &newAccessToken}, nil
		},
		githubProfileFetcher(githubOAuthProfile{
			Subject:       "github-existing",
			Email:         "changed-github-email@example.com",
			EmailVerified: true,
			Name:          "Changed GitHub Name",
			Login:         "changedgithub",
		}),
	)
	user := auth.User{Name: "Existing User", Email: "existing-github@example.com", EmailVerified: true}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create existing user: %v", err)
	}
	account := auth.AuthAccount{
		AccountID:    "github-existing",
		ProviderID:   auth.GitHubProviderID,
		UserID:       user.ID,
		RefreshToken: &oldRefreshToken,
	}
	if err := db.Create(&account).Error; err != nil {
		t.Fatalf("create existing github account: %v", err)
	}
	state := startGitHubOAuthState(t, router, redirectURI)

	w := authJSON(t, router, http.MethodPost, "/api/auth/oauth/github/callback", `{
		"code":"code-1",
		"state":"`+state+`",
		"redirectURI":"`+redirectURI+`"
	}`, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAuthResponse[authSessionPayload](t, w)
	if got.User.ID != user.ID || got.User.Email != "existing-github@example.com" || got.User.Name != "Existing User" {
		t.Fatalf("unexpected existing linked user response: %#v", got.User)
	}

	var updated auth.AuthAccount
	if err := db.First(&updated, "id = ?", account.ID).Error; err != nil {
		t.Fatalf("load updated github account: %v", err)
	}
	if updated.AccessToken == nil || *updated.AccessToken != newAccessToken {
		t.Fatalf("access token = %v, want %q", updated.AccessToken, newAccessToken)
	}
	if updated.RefreshToken == nil || *updated.RefreshToken != oldRefreshToken {
		t.Fatalf("refresh token = %v, want preserved %q", updated.RefreshToken, oldRefreshToken)
	}
}

func TestSignInGitHubOAuthRejectsInvalidStateWithoutExchange(t *testing.T) {
	redirectURI := "http://localhost:5173/api/auth/callback/github"
	exchanged := false
	router, _ := testGitHubAuthRouter(t,
		func(context.Context, string, string) (githubOAuthToken, error) {
			exchanged = true
			return githubOAuthToken{}, nil
		},
		githubProfileFetcher(githubOAuthProfile{
			Subject:       "12345",
			Email:         "ada@example.com",
			EmailVerified: true,
		}),
	)

	w := authJSON(t, router, http.MethodPost, "/api/auth/oauth/github/callback", `{
		"code":"code-1",
		"state":"missing-state",
		"redirectURI":"`+redirectURI+`"
	}`, nil)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s, want 400", w.Code, w.Body.String())
	}
	if exchanged {
		t.Fatal("oauth code exchange ran for invalid state")
	}
}

func TestSignInGitHubOAuthRejectsMissingPrimaryVerifiedEmailAndConsumesState(t *testing.T) {
	redirectURI := "http://localhost:5173/api/auth/callback/github"
	router, db := testGitHubAuthRouter(t, successfulGitHubExchange(t), githubProfileFetcher(githubOAuthProfile{
		Subject:       "github-invalid",
		Email:         "",
		EmailVerified: false,
		Login:         "invalidgithub",
	}))
	state := startGitHubOAuthState(t, router, redirectURI)

	w := authJSON(t, router, http.MethodPost, "/api/auth/oauth/github/callback", `{
		"code":"code-1",
		"state":"`+state+`",
		"redirectURI":"`+redirectURI+`"
	}`, nil)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d body=%s, want 401", w.Code, w.Body.String())
	}
	assertGitHubOAuthStateConsumed(t, db, state)
}

func TestGetSessionReturnsCurrentSession(t *testing.T) {
	router, _ := testAuthRouter(t)
	cookie := signUpVerifyAndSignIn(t, router)

	req := newTestRequest(http.MethodGet, "/api/auth/get-session", nil)
	req.AddCookie(cookie)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAuthResponse[authSessionPayload](t, w)
	if got.User.Email != "ada@example.com" {
		t.Fatalf("user email = %q, want ada@example.com", got.User.Email)
	}
	if got.Session.Token != cookie.Value {
		t.Fatalf("session token = %q, want cookie token %q", got.Session.Token, cookie.Value)
	}
}

func TestPatchAuthUserRequiresSession(t *testing.T) {
	router, _ := testAuthRouter(t)

	w := authJSON(t, router, http.MethodPatch, "/api/auth/user", `{"name":"Ada"}`, nil)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d want 401 response=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "authentication required") {
		t.Fatalf("expected auth error, got %s", w.Body.String())
	}
}

func TestPatchAuthUserUpdatesProfileFields(t *testing.T) {
	router, db := testAuthRouter(t)
	cookie := signUpVerifyAndSignIn(t, router)
	image := "data:image/png;base64,iVBORw0KGgo="

	req := newTestRequest(http.MethodPatch, "/api/auth/user", strings.NewReader(`{
		"name":" Ada Updated ",
		"email":" ADA.UPDATED@EXAMPLE.COM ",
		"image":"`+image+`"
	}`))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAuthResponse[authUserResponse](t, w)
	if got.Name != "Ada Updated" || got.Email != "ada.updated@example.com" || got.EmailVerified {
		t.Fatalf("unexpected user response: %#v", got)
	}
	if got.Image == nil || *got.Image != image {
		t.Fatalf("image = %#v, want %q", got.Image, image)
	}

	var user auth.User
	if err := db.First(&user, "id = ?", got.ID).Error; err != nil {
		t.Fatalf("load user: %v", err)
	}
	if user.Name != "Ada Updated" || user.Email != "ada.updated@example.com" || user.EmailVerified {
		t.Fatalf("unexpected stored user: %#v", user)
	}
}

func TestPatchAuthUserUpdatesPreferredLanguage(t *testing.T) {
	router, db := testAuthRouter(t)
	cookie := signUpVerifyAndSignIn(t, router)

	req := newTestRequest(http.MethodPatch, "/api/auth/user", strings.NewReader(`{"preferredLanguage":"ro"}`))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAuthResponse[authUserResponse](t, w)
	if got.PreferredLanguage != "ro" {
		t.Fatalf("preferredLanguage = %q, want ro", got.PreferredLanguage)
	}

	var user auth.User
	if err := db.First(&user, "id = ?", got.ID).Error; err != nil {
		t.Fatalf("load user: %v", err)
	}
	if user.PreferredLanguage != "ro" {
		t.Fatalf("stored preferred language = %q, want ro", user.PreferredLanguage)
	}
}

func TestPatchAuthUserRejectsInvalidPreferredLanguage(t *testing.T) {
	router, _ := testAuthRouter(t)
	cookie := signUpVerifyAndSignIn(t, router)

	req := newTestRequest(http.MethodPatch, "/api/auth/user", strings.NewReader(`{"preferredLanguage":"de"}`))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d want 400 response=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "preferredLanguage must be en or ro") {
		t.Fatalf("expected preferred language error, got %s", w.Body.String())
	}
}

func TestPatchAuthUserRejectsDuplicateEmail(t *testing.T) {
	router, db := testAuthRouter(t)
	cookie := signUpVerifyAndSignIn(t, router)
	other := auth.User{Name: "Other", Email: "other@example.com", EmailVerified: true}
	if err := db.Create(&other).Error; err != nil {
		t.Fatalf("create other user: %v", err)
	}

	req := newTestRequest(http.MethodPatch, "/api/auth/user", strings.NewReader(`{"email":"other@example.com"}`))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("status = %d want 409 response=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "email is already in use") {
		t.Fatalf("expected duplicate email error, got %s", w.Body.String())
	}
}

func TestPatchAuthUserSetsPasswordAndCreatesCredentialAccount(t *testing.T) {
	router, db := testAuthRouter(t)
	user := auth.User{Name: "Social User", Email: "social@example.com", EmailVerified: true}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	session := auth.Session{
		Token:     "session-token",
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().Add(time.Hour),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	if err := db.Create(&session).Error; err != nil {
		t.Fatalf("create session: %v", err)
	}

	req := newTestRequest(http.MethodPatch, "/api/auth/user", strings.NewReader(`{"password":"newpassword123"}`))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: authSessionCookieName, Value: session.Token})
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var account auth.AuthAccount
	if err := db.First(&account, "user_id = ? AND provider_id = ?", user.ID, auth.CredentialProviderID).Error; err != nil {
		t.Fatalf("load account: %v", err)
	}
	if account.Password == "" || account.Password == "newpassword123" {
		t.Fatalf("unsafe password hash: %#v", account)
	}
	if !auth.VerifyPassword("newpassword123", account.Password) {
		t.Fatal("stored password hash does not verify")
	}
}

func TestPatchAuthUserUpdatesExistingCredentialPassword(t *testing.T) {
	router, db := testAuthRouter(t)
	cookie := signUpVerifyAndSignIn(t, router)

	var user auth.User
	if err := db.First(&user, "email = ?", "ada@example.com").Error; err != nil {
		t.Fatalf("load user: %v", err)
	}
	var before auth.AuthAccount
	if err := db.First(&before, "user_id = ? AND provider_id = ?", user.ID, auth.CredentialProviderID).Error; err != nil {
		t.Fatalf("load account: %v", err)
	}

	req := newTestRequest(http.MethodPatch, "/api/auth/user", strings.NewReader(`{"password":"updatedpassword123"}`))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var after auth.AuthAccount
	if err := db.First(&after, "user_id = ? AND provider_id = ?", user.ID, auth.CredentialProviderID).Error; err != nil {
		t.Fatalf("load updated account: %v", err)
	}
	if after.ID != before.ID {
		t.Fatalf("credential account ID changed from %q to %q", before.ID, after.ID)
	}
	if after.Password == "" || after.Password == before.Password {
		t.Fatalf("password hash was not updated: before=%q after=%q", before.Password, after.Password)
	}
	if !auth.VerifyPassword("updatedpassword123", after.Password) {
		t.Fatal("updated password hash does not verify")
	}
	var count int64
	if err := db.Model(&auth.AuthAccount{}).Where("user_id = ? AND provider_id = ?", user.ID, auth.CredentialProviderID).Count(&count).Error; err != nil {
		t.Fatalf("count accounts: %v", err)
	}
	if count != 1 {
		t.Fatalf("credential account count = %d, want 1", count)
	}
}

func TestPatchAuthUserRejectsPasswordUnderEightCharacters(t *testing.T) {
	router, _ := testAuthRouter(t)
	cookie := signUpVerifyAndSignIn(t, router)
	password := strings.Repeat(`\u754c`, 7)

	req := newTestRequest(http.MethodPatch, "/api/auth/user", strings.NewReader(`{"password":"`+password+`"}`))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d want 400 response=%s", w.Code, w.Body.String())
	}
}

func TestPatchAuthUserClearsAvatarWithNull(t *testing.T) {
	router, db := testAuthRouter(t)
	cookie := signUpVerifyAndSignIn(t, router)
	image := "data:image/png;base64,iVBORw0KGgo="
	if err := db.Model(&auth.User{}).Where("email = ?", "ada@example.com").Update("image", image).Error; err != nil {
		t.Fatalf("set image: %v", err)
	}

	req := newTestRequest(http.MethodPatch, "/api/auth/user", strings.NewReader(`{"image":null}`))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAuthResponse[authUserResponse](t, w)
	if got.Image != nil {
		t.Fatalf("image = %#v, want nil", got.Image)
	}
	var user auth.User
	if err := db.First(&user, "id = ?", got.ID).Error; err != nil {
		t.Fatalf("load user: %v", err)
	}
	if user.Image != nil {
		t.Fatalf("stored image = %#v, want nil", user.Image)
	}
}

func TestPatchAuthUserRejectsInvalidAvatarDataURL(t *testing.T) {
	router, _ := testAuthRouter(t)
	cookie := signUpVerifyAndSignIn(t, router)
	req := newTestRequest(http.MethodPatch, "/api/auth/user", strings.NewReader(`{"image":"data:text/plain;base64,SGk="}`))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d want 400 response=%s", w.Code, w.Body.String())
	}
}

func TestPatchAuthUserInvalidCookieClearsSessionCookie(t *testing.T) {
	router, _ := testAuthRouter(t)

	req := newTestRequest(http.MethodPatch, "/api/auth/user", strings.NewReader(`{"name":"Ada"}`))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: authSessionCookieName, Value: "missing-session-token"})
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d want 401 response=%s", w.Code, w.Body.String())
	}
	clearedCookie := authCookieFromRecorder(t, w)
	if clearedCookie.Value != "" || clearedCookie.MaxAge != -1 {
		t.Fatalf("cookie was not cleared: %#v", clearedCookie)
	}
}

func TestPatchAuthUserRejectsOversizedBody(t *testing.T) {
	router, _ := testAuthRouter(t)
	cookie := signUpVerifyAndSignIn(t, router)
	body := `{"image":"` + strings.Repeat("A", maxPatchAuthUserBodyBytes) + `"}`

	req := newTestRequest(http.MethodPatch, "/api/auth/user", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d want 400 response=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "request body too large") {
		t.Fatalf("expected body size error, got %s", w.Body.String())
	}
}

func TestGetSessionWithoutCookieReturnsNull(t *testing.T) {
	router, _ := testAuthRouter(t)

	w := authJSON(t, router, http.MethodGet, "/api/auth/get-session", "", nil)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if strings.TrimSpace(w.Body.String()) != "null" {
		t.Fatalf("body = %q, want JSON null", w.Body.String())
	}
	if setCookie := w.Result().Header.Values("Set-Cookie"); len(setCookie) != 0 {
		t.Fatalf("unexpected Set-Cookie header: %v", setCookie)
	}
}

func TestGetSessionBlankCookieReturnsNullAndClearsCookie(t *testing.T) {
	router, _ := testAuthRouter(t)

	tests := []struct {
		name         string
		cookieHeader string
	}{
		{name: "empty", cookieHeader: authSessionCookieName + "="},
		{name: "whitespace", cookieHeader: authSessionCookieName + "=   "},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := newTestRequest(http.MethodGet, "/api/auth/get-session", nil)
			req.Header.Set("Cookie", tt.cookieHeader)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
			}
			if strings.TrimSpace(w.Body.String()) != "null" {
				t.Fatalf("body = %q, want JSON null", w.Body.String())
			}
			clearedCookie := authCookieFromRecorder(t, w)
			if clearedCookie.Value != "" || clearedCookie.MaxAge != -1 {
				t.Fatalf("cookie was not cleared: %#v", clearedCookie)
			}
		})
	}
}

func TestGetSessionInvalidCookieReturnsNull(t *testing.T) {
	router, _ := testAuthRouter(t)

	req := newTestRequest(http.MethodGet, "/api/auth/get-session", nil)
	req.AddCookie(&http.Cookie{Name: authSessionCookieName, Value: "missing-session-token"})
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if strings.TrimSpace(w.Body.String()) != "null" {
		t.Fatalf("body = %q, want JSON null", w.Body.String())
	}
	clearedCookie := authCookieFromRecorder(t, w)
	if clearedCookie.Value != "" || clearedCookie.MaxAge != -1 {
		t.Fatalf("cookie was not cleared: %#v", clearedCookie)
	}
}

func TestGetSessionExpiredSessionDeletesRowClearsCookieAndReturnsNull(t *testing.T) {
	router, db := testAuthRouter(t)
	cookie := signUpVerifyAndSignIn(t, router)
	if err := db.Model(&auth.Session{}).
		Where("token = ?", cookie.Value).
		Update("expires_at", time.Now().UTC().Add(-time.Minute)).Error; err != nil {
		t.Fatalf("expire session: %v", err)
	}

	req := newTestRequest(http.MethodGet, "/api/auth/get-session", nil)
	req.AddCookie(cookie)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if strings.TrimSpace(w.Body.String()) != "null" {
		t.Fatalf("body = %q, want JSON null", w.Body.String())
	}
	clearedCookie := authCookieFromRecorder(t, w)
	if clearedCookie.Value != "" || clearedCookie.MaxAge != -1 {
		t.Fatalf("cookie was not cleared: %#v", clearedCookie)
	}
	var count int64
	if err := db.Model(&auth.Session{}).Where("token = ?", cookie.Value).Count(&count).Error; err != nil {
		t.Fatalf("count sessions: %v", err)
	}
	if count != 0 {
		t.Fatalf("session count = %d, want 0", count)
	}
}

func TestGetSessionStaleSessionWithInvalidUserDeletesRowClearsCookieAndReturnsNull(t *testing.T) {
	router, db := testAuthRouter(t)
	now := time.Now().UTC()
	user := auth.User{
		Name:          "Stale",
		Email:         "",
		EmailVerified: true,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create invalid user: %v", err)
	}
	session := auth.Session{
		Token:     "stale-session-token",
		UserID:    user.ID,
		ExpiresAt: now.Add(time.Hour),
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := db.Create(&session).Error; err != nil {
		t.Fatalf("create stale session: %v", err)
	}

	req := newTestRequest(http.MethodGet, "/api/auth/get-session", nil)
	req.AddCookie(&http.Cookie{Name: authSessionCookieName, Value: session.Token})
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if strings.TrimSpace(w.Body.String()) != "null" {
		t.Fatalf("body = %q, want JSON null", w.Body.String())
	}
	clearedCookie := authCookieFromRecorder(t, w)
	if clearedCookie.Value != "" || clearedCookie.MaxAge != -1 {
		t.Fatalf("cookie was not cleared: %#v", clearedCookie)
	}
	var count int64
	if err := db.Model(&auth.Session{}).Where("token = ?", session.Token).Count(&count).Error; err != nil {
		t.Fatalf("count sessions: %v", err)
	}
	if count != 0 {
		t.Fatalf("session count = %d, want 0", count)
	}
}

func TestSignOutDeletesSessionAndClearsCookie(t *testing.T) {
	router, db := testAuthRouter(t)
	cookie := signUpVerifyAndSignIn(t, router)

	req := newTestRequest(http.MethodPost, "/api/auth/sign-out", nil)
	req.AddCookie(cookie)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAuthResponse[struct {
		Success bool `json:"success"`
	}](t, w)
	if !got.Success {
		t.Fatalf("success = false response=%s", w.Body.String())
	}
	clearedCookie := authCookieFromRecorder(t, w)
	if clearedCookie.Value != "" || clearedCookie.MaxAge != -1 {
		t.Fatalf("cookie was not cleared: %#v", clearedCookie)
	}
	var count int64
	if err := db.Model(&auth.Session{}).Where("token = ?", cookie.Value).Count(&count).Error; err != nil {
		t.Fatalf("count sessions: %v", err)
	}
	if count != 0 {
		t.Fatalf("session count = %d, want 0", count)
	}
}

func TestListAuthSessionsReturnsSafeCurrentUserSessions(t *testing.T) {
	router, db := testAuthRouter(t)
	cookie := signUpVerifyAndSignIn(t, router)

	var current auth.Session
	if err := db.First(&current, "token = ?", cookie.Value).Error; err != nil {
		t.Fatalf("load current session: %v", err)
	}
	now := time.Now().UTC()
	other := auth.Session{
		Token:     "other-session-token",
		UserID:    current.UserID,
		ExpiresAt: now.Add(2 * time.Hour),
		IPAddress: "203.0.113.10",
		UserAgent: "Other Browser",
		CreatedAt: now.Add(time.Minute),
		UpdatedAt: now.Add(time.Minute),
	}
	if err := db.Create(&other).Error; err != nil {
		t.Fatalf("create other session: %v", err)
	}
	expired := auth.Session{
		Token:     "expired-session-token",
		UserID:    current.UserID,
		ExpiresAt: now.Add(-time.Minute),
		CreatedAt: now.Add(2 * time.Minute),
		UpdatedAt: now.Add(2 * time.Minute),
	}
	if err := db.Create(&expired).Error; err != nil {
		t.Fatalf("create expired session: %v", err)
	}
	otherUser := auth.User{Name: "Other", Email: "sessions-other@example.com", EmailVerified: true}
	if err := db.Create(&otherUser).Error; err != nil {
		t.Fatalf("create other user: %v", err)
	}
	if err := db.Create(&auth.Session{
		Token:     "different-user-session-token",
		UserID:    otherUser.ID,
		ExpiresAt: now.Add(time.Hour),
		CreatedAt: now.Add(3 * time.Minute),
		UpdatedAt: now.Add(3 * time.Minute),
	}).Error; err != nil {
		t.Fatalf("create different user session: %v", err)
	}

	req := newTestRequest(http.MethodGet, "/api/auth/sessions", nil)
	req.AddCookie(cookie)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if bytes.Contains(w.Body.Bytes(), []byte("token")) || bytes.Contains(w.Body.Bytes(), []byte(cookie.Value)) {
		t.Fatalf("session response leaked token: %s", w.Body.String())
	}
	got := decodeAuthResponse[authSessionListResponse](t, w)
	if len(got.Sessions) != 2 {
		t.Fatalf("session count = %d want 2 response=%#v", len(got.Sessions), got)
	}
	if got.Sessions[0].ID != other.ID {
		t.Fatalf("first session id = %q want newest other session %q", got.Sessions[0].ID, other.ID)
	}
	currentCount := 0
	for _, session := range got.Sessions {
		if session.UserID != current.UserID {
			t.Fatalf("session user_id = %q want %q", session.UserID, current.UserID)
		}
		if session.ID == current.ID && session.Current {
			currentCount++
		}
		if session.ID == other.ID && (session.Current || session.IPAddress != "203.0.113.10" || session.UserAgent != "Other Browser") {
			t.Fatalf("unexpected other session response: %#v", session)
		}
	}
	if currentCount != 1 {
		t.Fatalf("current session marker count = %d want 1 response=%#v", currentCount, got.Sessions)
	}
}

func TestRevokeAuthSessionDeletesOnlyOwnedNonCurrentSession(t *testing.T) {
	router, db := testAuthRouter(t)
	cookie := signUpVerifyAndSignIn(t, router)

	var current auth.Session
	if err := db.First(&current, "token = ?", cookie.Value).Error; err != nil {
		t.Fatalf("load current session: %v", err)
	}
	other := auth.Session{
		Token:     "revoke-other-session-token",
		UserID:    current.UserID,
		ExpiresAt: time.Now().UTC().Add(time.Hour),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	if err := db.Create(&other).Error; err != nil {
		t.Fatalf("create other session: %v", err)
	}

	req := newTestRequest(http.MethodDelete, "/api/auth/sessions/"+other.ID, nil)
	req.AddCookie(cookie)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAuthResponse[deleteAuthSessionResponse](t, w)
	if got.DeletedID != other.ID || got.DeletedCount != 1 {
		t.Fatalf("unexpected delete response: %#v", got)
	}
	var count int64
	if err := db.Model(&auth.Session{}).Where("id = ?", other.ID).Count(&count).Error; err != nil {
		t.Fatalf("count revoked session: %v", err)
	}
	if count != 0 {
		t.Fatalf("revoked session count = %d want 0", count)
	}
	if err := db.Model(&auth.Session{}).Where("id = ?", current.ID).Count(&count).Error; err != nil {
		t.Fatalf("count current session: %v", err)
	}
	if count != 1 {
		t.Fatalf("current session count = %d want 1", count)
	}
}

func TestRevokeAuthSessionRejectsCurrentAndNonOwnedSessions(t *testing.T) {
	router, db := testAuthRouter(t)
	cookie := signUpVerifyAndSignIn(t, router)

	var current auth.Session
	if err := db.First(&current, "token = ?", cookie.Value).Error; err != nil {
		t.Fatalf("load current session: %v", err)
	}
	req := newTestRequest(http.MethodDelete, "/api/auth/sessions/"+current.ID, nil)
	req.AddCookie(cookie)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("current revoke status = %d want 400 body=%s", w.Code, w.Body.String())
	}

	otherUser := auth.User{Name: "Other", Email: "revoke-other@example.com", EmailVerified: true}
	if err := db.Create(&otherUser).Error; err != nil {
		t.Fatalf("create other user: %v", err)
	}
	otherSession := auth.Session{
		Token:     "non-owned-session-token",
		UserID:    otherUser.ID,
		ExpiresAt: time.Now().UTC().Add(time.Hour),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	if err := db.Create(&otherSession).Error; err != nil {
		t.Fatalf("create non-owned session: %v", err)
	}
	req = newTestRequest(http.MethodDelete, "/api/auth/sessions/"+otherSession.ID, nil)
	req.AddCookie(cookie)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("non-owned revoke status = %d want 404 body=%s", w.Code, w.Body.String())
	}
}

func TestListAndUnlinkAuthAccountsUseSafePayloadsAndSafeguards(t *testing.T) {
	router, db := testAuthRouter(t)
	cookie := signUpVerifyAndSignIn(t, router)

	var user auth.User
	if err := db.First(&user, "email = ?", "ada@example.com").Error; err != nil {
		t.Fatalf("load user: %v", err)
	}
	accessToken := "secret-google-access-token"
	googleAccount := auth.AuthAccount{
		AccountID:   "google-subject",
		ProviderID:  auth.GoogleProviderID,
		UserID:      user.ID,
		AccessToken: &accessToken,
	}
	if err := db.Create(&googleAccount).Error; err != nil {
		t.Fatalf("create google account: %v", err)
	}

	req := newTestRequest(http.MethodGet, "/api/auth/accounts", nil)
	req.AddCookie(cookie)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if bytes.Contains(w.Body.Bytes(), []byte("secret-google-access-token")) || bytes.Contains(w.Body.Bytes(), []byte("google-subject")) {
		t.Fatalf("account response leaked sensitive data: %s", w.Body.String())
	}
	got := decodeAuthResponse[authAccountListResponse](t, w)
	if len(got.Accounts) != 2 {
		t.Fatalf("account count = %d want credential and google response=%#v", len(got.Accounts), got)
	}

	req = newTestRequest(http.MethodDelete, "/api/auth/accounts/credential", nil)
	req.AddCookie(cookie)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("credential unlink status = %d want 400 body=%s", w.Code, w.Body.String())
	}

	req = newTestRequest(http.MethodDelete, "/api/auth/accounts/google", nil)
	req.AddCookie(cookie)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("google unlink status = %d body=%s", w.Code, w.Body.String())
	}
	deleted := decodeAuthResponse[deleteAuthAccountResponse](t, w)
	if deleted.DeletedProviderID != auth.GoogleProviderID || deleted.DeletedCount != 1 {
		t.Fatalf("unexpected unlink response: %#v", deleted)
	}
}

func TestUnlinkAuthAccountRejectsLastSignInMethod(t *testing.T) {
	router, db := testAuthRouter(t)
	now := time.Now().UTC()
	user := auth.User{Name: "Social", Email: "social-last@example.com", EmailVerified: true}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	session := auth.Session{
		Token:     "social-session-token",
		UserID:    user.ID,
		ExpiresAt: now.Add(time.Hour),
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := db.Create(&session).Error; err != nil {
		t.Fatalf("create session: %v", err)
	}
	if err := db.Create(&auth.AuthAccount{
		AccountID:  "github-only",
		ProviderID: auth.GitHubProviderID,
		UserID:     user.ID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}).Error; err != nil {
		t.Fatalf("create github account: %v", err)
	}

	req := newTestRequest(http.MethodDelete, "/api/auth/accounts/github", nil)
	req.AddCookie(&http.Cookie{Name: authSessionCookieName, Value: session.Token})
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d want 400 body=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), errLastSignInMethod.Error()) {
		t.Fatalf("expected last sign-in method error, got %s", w.Body.String())
	}
}

func TestLinkGoogleAccountCreatesAccountWithoutCreatingSession(t *testing.T) {
	redirectURI := "http://localhost:5173/api/auth/callback/link/google"
	router, db := testGoogleAuthRouter(t, successfulGoogleExchange(t), googleProfileValidator(googleIDTokenPayload{
		Subject:       "google-link-subject",
		Email:         "google-link@example.com",
		EmailVerified: true,
		Name:          "Google Link",
	}))
	cookie := signUpVerifyAndSignIn(t, router)
	state := startGoogleAccountLinkState(t, router, redirectURI, cookie)

	var sessionCountBefore int64
	if err := db.Model(&auth.Session{}).Count(&sessionCountBefore).Error; err != nil {
		t.Fatalf("count sessions before: %v", err)
	}
	w := authJSONWithCookie(t, router, http.MethodPost, "/api/auth/accounts/google/callback", `{
		"code":"code-1",
		"state":"`+state+`",
		"redirectURI":"`+redirectURI+`"
	}`, cookie)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAuthResponse[authAccountListItemResponse](t, w)
	if got.ProviderID != auth.GoogleProviderID || got.ID == "" {
		t.Fatalf("unexpected linked account response: %#v", got)
	}
	if bytes.Contains(w.Body.Bytes(), []byte("google-link-subject")) || bytes.Contains(w.Body.Bytes(), []byte("access-token")) {
		t.Fatalf("link response leaked sensitive data: %s", w.Body.String())
	}
	assertGoogleOAuthLinkStateConsumed(t, db, state)
	var sessionCountAfter int64
	if err := db.Model(&auth.Session{}).Count(&sessionCountAfter).Error; err != nil {
		t.Fatalf("count sessions after: %v", err)
	}
	if sessionCountAfter != sessionCountBefore {
		t.Fatalf("session count after link = %d want %d", sessionCountAfter, sessionCountBefore)
	}
}

func TestLinkGoogleAccountRejectsIdentityLinkedToAnotherUser(t *testing.T) {
	redirectURI := "http://localhost:5173/api/auth/callback/link/google"
	router, db := testGoogleAuthRouter(t, successfulGoogleExchange(t), googleProfileValidator(googleIDTokenPayload{
		Subject:       "google-owned-subject",
		Email:         "owned-google@example.com",
		EmailVerified: true,
	}))
	cookie := signUpVerifyAndSignIn(t, router)
	other := auth.User{Name: "Other", Email: "google-owner@example.com", EmailVerified: true}
	if err := db.Create(&other).Error; err != nil {
		t.Fatalf("create other user: %v", err)
	}
	if err := db.Create(&auth.AuthAccount{
		AccountID:  "google-owned-subject",
		ProviderID: auth.GoogleProviderID,
		UserID:     other.ID,
	}).Error; err != nil {
		t.Fatalf("create existing google account: %v", err)
	}
	state := startGoogleAccountLinkState(t, router, redirectURI, cookie)

	w := authJSONWithCookie(t, router, http.MethodPost, "/api/auth/accounts/google/callback", `{
		"code":"code-1",
		"state":"`+state+`",
		"redirectURI":"`+redirectURI+`"
	}`, cookie)
	if w.Code != http.StatusConflict {
		t.Fatalf("status = %d want 409 body=%s", w.Code, w.Body.String())
	}
	assertGoogleOAuthLinkStateConsumed(t, db, state)
}

func TestLinkGitHubAccountCreatesAccountWithoutCreatingSession(t *testing.T) {
	redirectURI := "http://localhost:5173/api/auth/callback/link/github"
	router, db := testGitHubAuthRouter(t, successfulGitHubExchange(t), githubProfileFetcher(githubOAuthProfile{
		Subject:       "github-link-subject",
		Email:         "github-link@example.com",
		EmailVerified: true,
		Login:         "linkedgithub",
	}))
	cookie := signUpVerifyAndSignIn(t, router)
	state := startGitHubAccountLinkState(t, router, redirectURI, cookie)

	var sessionCountBefore int64
	if err := db.Model(&auth.Session{}).Count(&sessionCountBefore).Error; err != nil {
		t.Fatalf("count sessions before: %v", err)
	}
	w := authJSONWithCookie(t, router, http.MethodPost, "/api/auth/accounts/github/callback", `{
		"code":"code-1",
		"state":"`+state+`",
		"redirectURI":"`+redirectURI+`"
	}`, cookie)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAuthResponse[authAccountListItemResponse](t, w)
	if got.ProviderID != auth.GitHubProviderID || got.ID == "" {
		t.Fatalf("unexpected linked account response: %#v", got)
	}
	assertGitHubOAuthLinkStateConsumed(t, db, state)
	var sessionCountAfter int64
	if err := db.Model(&auth.Session{}).Count(&sessionCountAfter).Error; err != nil {
		t.Fatalf("count sessions after: %v", err)
	}
	if sessionCountAfter != sessionCountBefore {
		t.Fatalf("session count after link = %d want %d", sessionCountAfter, sessionCountBefore)
	}
}

func TestSignInEmailRejectsWrongPassword(t *testing.T) {
	router, _ := testAuthRouter(t)
	code := signupVerificationCode(t, router, "ada@example.com")
	verifySignupEmail(t, router, "ada@example.com", code)

	w := authJSON(t, router, http.MethodPost, "/api/auth/sign-in/email", `{
		"email":"ada@example.com",
		"password":"wrong-password"
	}`, nil)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d want 401 response=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "invalid email or password") {
		t.Fatalf("expected generic credential error, got %s", w.Body.String())
	}
}

func TestSignInEmailMissingUserAndWrongPasswordReturnSameUnauthorizedBody(t *testing.T) {
	router, _ := testAuthRouter(t)
	code := signupVerificationCode(t, router, "ada@example.com")
	verifySignupEmail(t, router, "ada@example.com", code)

	missing := authJSON(t, router, http.MethodPost, "/api/auth/sign-in/email", `{
		"email":"missing@example.com",
		"password":"password1234"
	}`, nil)
	if missing.Code != http.StatusUnauthorized {
		t.Fatalf("missing status = %d want 401 response=%s", missing.Code, missing.Body.String())
	}

	wrongPassword := authJSON(t, router, http.MethodPost, "/api/auth/sign-in/email", `{
		"email":"ada@example.com",
		"password":"wrong-password"
	}`, nil)
	if wrongPassword.Code != http.StatusUnauthorized {
		t.Fatalf("wrong password status = %d want 401 response=%s", wrongPassword.Code, wrongPassword.Body.String())
	}
	if missing.Body.String() != wrongPassword.Body.String() {
		t.Fatalf("401 bodies differ: missing=%s wrongPassword=%s", missing.Body.String(), wrongPassword.Body.String())
	}
}

func TestSignInEmailRejectsInvalidInput(t *testing.T) {
	router, _ := testAuthRouter(t)
	overlongEmail := strings.Repeat("a", 309) + "@example.com"

	cases := []string{
		`{"email":"","password":"password1234"}`,
		`{"email":"not-email","password":"password1234"}`,
		`{"email":"` + overlongEmail + `","password":"password1234"}`,
		`{"email":"ada@example.com","password":""}`,
	}
	for _, body := range cases {
		w := authJSON(t, router, http.MethodPost, "/api/auth/sign-in/email", body, nil)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("body %s status = %d want 400 response=%s", body, w.Code, w.Body.String())
		}
	}
}

func TestSendVerificationOTPRejectsInvalidEmail(t *testing.T) {
	router, _ := testAuthRouter(t)

	w := authJSON(t, router, http.MethodPost, "/api/auth/email-otp/send-verification-otp", `{
		"email":"not-email",
		"type":"email-verification"
	}`, nil)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d want 400 response=%s", w.Code, w.Body.String())
	}
}

func TestSendVerificationOTPRejectsOverlongEmail(t *testing.T) {
	router, _ := testAuthRouter(t)
	overlongEmail := strings.Repeat("a", 309) + "@example.com"

	w := authJSON(t, router, http.MethodPost, "/api/auth/email-otp/send-verification-otp", `{
		"email":"`+overlongEmail+`",
		"type":"email-verification"
	}`, nil)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d want 400 response=%s", w.Code, w.Body.String())
	}
}

func TestSendVerificationOTPRejectsInvalidType(t *testing.T) {
	router, _ := testAuthRouter(t)

	w := authJSON(t, router, http.MethodPost, "/api/auth/email-otp/send-verification-otp", `{
		"email":"ada@example.com",
		"type":"password-reset"
	}`, nil)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d want 400 response=%s", w.Code, w.Body.String())
	}
}

func TestSendVerificationOTPMissingUserReturnsOKWithoutCode(t *testing.T) {
	router, _ := testAuthRouter(t)

	w := authJSON(t, router, http.MethodPost, "/api/auth/email-otp/send-verification-otp", `{
		"email":"missing@example.com"
	}`, map[string]string{authDeliveryHeader: testDeliveryToken})
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAuthResponse[sendVerificationOTPResponse](t, w)
	if !got.OK {
		t.Fatalf("ok = %v, want true", got.OK)
	}
	if bytes.Contains(w.Body.Bytes(), []byte("verificationCode")) {
		t.Fatalf("missing user response leaked verification code: %s", w.Body.String())
	}
}

func TestSendVerificationOTPRotatesUnverifiedUserCodeForTrustedRequest(t *testing.T) {
	router, db := testAuthRouter(t)

	signup := authJSON(t, router, http.MethodPost, "/api/auth/sign-up/email", `{
		"name":"Ada",
		"email":"ADA@EXAMPLE.COM",
		"password":"password1234"
	}`, map[string]string{authDeliveryHeader: testDeliveryToken})
	if signup.Code != http.StatusOK {
		t.Fatalf("signup status = %d body=%s", signup.Code, signup.Body.String())
	}
	signupResponse := decodeAuthResponse[struct {
		VerificationCode string `json:"verificationCode"`
	}](t, signup)

	identifier := verificationIdentifier("ada@example.com")
	var firstVerification auth.Verification
	if err := db.First(&firstVerification, "identifier = ?", identifier).Error; err != nil {
		t.Fatalf("load first verification: %v", err)
	}

	var otpResponse sendVerificationOTPResponse
	var secondVerification auth.Verification
	for i := 0; i < 5; i++ {
		otp := authJSON(t, router, http.MethodPost, "/api/auth/email-otp/send-verification-otp", `{
			"email":"ADA@EXAMPLE.COM",
			"type":"email-verification"
		}`, map[string]string{authDeliveryHeader: testDeliveryToken})
		if otp.Code != http.StatusOK {
			t.Fatalf("otp status = %d body=%s", otp.Code, otp.Body.String())
		}
		otpResponse = decodeAuthResponse[sendVerificationOTPResponse](t, otp)
		if !otpResponse.OK || otpResponse.VerificationCode == "" || otpResponse.VerificationExpiresAt == "" {
			t.Fatalf("unexpected otp response: %#v", otpResponse)
		}
		if err := db.First(&secondVerification, "identifier = ?", identifier).Error; err != nil {
			t.Fatalf("load second verification: %v", err)
		}
		if secondVerification.Value != firstVerification.Value {
			break
		}
	}

	if secondVerification.Value == firstVerification.Value {
		t.Fatal("verification hash was not replaced")
	}
	var count int64
	if err := db.Model(&auth.Verification{}).Where("identifier = ?", identifier).Count(&count).Error; err != nil {
		t.Fatalf("count verifications: %v", err)
	}
	if count != 1 {
		t.Fatalf("verification count = %d, want 1", count)
	}
	if authcrypto.VerifyCode(testAuthSecret, identifier, signupResponse.VerificationCode, secondVerification.Value) {
		t.Fatal("old signup code still verifies after OTP rotation")
	}
	if !authcrypto.VerifyCode(testAuthSecret, identifier, otpResponse.VerificationCode, secondVerification.Value) {
		t.Fatal("new OTP code does not verify")
	}
}

func TestSendVerificationOTPDoesNotLeakCodeWithoutTrustedHeader(t *testing.T) {
	router, _ := testAuthRouter(t)

	signup := authJSON(t, router, http.MethodPost, "/api/auth/sign-up/email", `{
		"name":"Ada",
		"email":"ada@example.com",
		"password":"password1234"
	}`, map[string]string{authDeliveryHeader: testDeliveryToken})
	if signup.Code != http.StatusOK {
		t.Fatalf("signup status = %d body=%s", signup.Code, signup.Body.String())
	}

	w := authJSON(t, router, http.MethodPost, "/api/auth/email-otp/send-verification-otp", `{
		"email":"ada@example.com",
		"type":"email-verification"
	}`, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAuthResponse[sendVerificationOTPResponse](t, w)
	if !got.OK {
		t.Fatalf("ok = %v, want true", got.OK)
	}
	if bytes.Contains(w.Body.Bytes(), []byte("verificationCode")) {
		t.Fatalf("untrusted response leaked verification code: %s", w.Body.String())
	}
}

func TestSendVerificationOTPVerifiedUserReturnsOKWithoutCode(t *testing.T) {
	router, db := testAuthRouter(t)
	if err := db.Create(&auth.User{
		Name:          "Ada",
		Email:         "ada@example.com",
		EmailVerified: true,
	}).Error; err != nil {
		t.Fatalf("create verified user: %v", err)
	}

	w := authJSON(t, router, http.MethodPost, "/api/auth/email-otp/send-verification-otp", `{
		"email":"ada@example.com",
		"type":"email-verification"
	}`, map[string]string{authDeliveryHeader: testDeliveryToken})
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAuthResponse[sendVerificationOTPResponse](t, w)
	if !got.OK {
		t.Fatalf("ok = %v, want true", got.OK)
	}
	if bytes.Contains(w.Body.Bytes(), []byte("verificationCode")) {
		t.Fatalf("verified user response leaked verification code: %s", w.Body.String())
	}
}

func TestVerifyEmailOTPMarksUserVerified(t *testing.T) {
	router, db := testAuthRouter(t)

	code := signupVerificationCode(t, router, "ADA@EXAMPLE.COM")

	w := authJSON(t, router, http.MethodPost, "/api/auth/email-otp/verify-email", `{
		"email":" ada@example.com ",
		"otp":" `+code+` "
	}`, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAuthResponse[struct {
		OK   bool             `json:"ok"`
		User authUserResponse `json:"user"`
	}](t, w)
	if !got.OK {
		t.Fatalf("ok = %v, want true", got.OK)
	}
	if got.User.Email != "ada@example.com" || !got.User.EmailVerified {
		t.Fatalf("unexpected user response: %#v", got.User)
	}

	var user auth.User
	if err := db.First(&user, "email = ?", "ada@example.com").Error; err != nil {
		t.Fatalf("load user: %v", err)
	}
	if !user.EmailVerified {
		t.Fatal("user email_verified = false, want true")
	}

	var verificationCount int64
	if err := db.Model(&auth.Verification{}).Where("identifier = ?", verificationIdentifier("ada@example.com")).Count(&verificationCount).Error; err != nil {
		t.Fatalf("count verifications: %v", err)
	}
	if verificationCount != 0 {
		t.Fatalf("verification count = %d, want 0", verificationCount)
	}

	var bucket billing.CreditBucket
	if err := db.First(&bucket, "user_id = ? AND source_type = ?", user.ID, billing.CreditSourceSignupBonus).Error; err != nil {
		t.Fatalf("load signup bonus bucket: %v", err)
	}
	if bucket.CreditsGranted != testOnboardingCredits || bucket.CreditsRemaining != testOnboardingCredits {
		t.Fatalf("signup bonus credits = granted %d remaining %d, want %d/%d", bucket.CreditsGranted, bucket.CreditsRemaining, testOnboardingCredits, testOnboardingCredits)
	}
	if bucket.ExpiresAt != nil {
		t.Fatalf("signup bonus expires_at = %v, want nil", bucket.ExpiresAt)
	}

	var entry billing.CreditLedgerEntry
	if err := db.First(&entry, "user_id = ? AND idempotency_key = ?", user.ID, "signup_bonus:"+user.ID).Error; err != nil {
		t.Fatalf("load signup bonus ledger entry: %v", err)
	}
	if entry.BucketID == nil || *entry.BucketID != bucket.ID || entry.EntryType != billing.CreditLedgerEntryGrant || entry.CreditsDelta != testOnboardingCredits {
		t.Fatalf("unexpected signup bonus ledger entry: %#v", entry)
	}
}

func TestVerifyEmailOTPDoesNotDuplicateExistingSignupBonus(t *testing.T) {
	router, db := testAuthRouter(t)

	code := signupVerificationCode(t, router, "ada@example.com")
	var user auth.User
	if err := db.First(&user, "email = ?", "ada@example.com").Error; err != nil {
		t.Fatalf("load user: %v", err)
	}
	existing, err := billing.GrantSignupBonus(t.Context(), db, billing.GrantSignupBonusInput{
		UserID:  user.ID,
		Credits: testOnboardingCredits,
		Now:     time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("grant existing signup bonus: %v", err)
	}

	w := authJSON(t, router, http.MethodPost, "/api/auth/email-otp/verify-email", `{
		"email":"ada@example.com",
		"otp":"`+code+`"
	}`, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}

	var bucketCount int64
	if err := db.Model(&billing.CreditBucket{}).Where("user_id = ? AND source_type = ?", user.ID, billing.CreditSourceSignupBonus).Count(&bucketCount).Error; err != nil {
		t.Fatalf("count signup bonus buckets: %v", err)
	}
	if bucketCount != 1 {
		t.Fatalf("signup bonus bucket count = %d, want 1", bucketCount)
	}
	var bucket billing.CreditBucket
	if err := db.First(&bucket, "user_id = ? AND source_type = ?", user.ID, billing.CreditSourceSignupBonus).Error; err != nil {
		t.Fatalf("load signup bonus bucket: %v", err)
	}
	if bucket.ID != existing.ID {
		t.Fatalf("signup bonus bucket id = %s, want existing %s", bucket.ID, existing.ID)
	}
	var ledgerCount int64
	if err := db.Model(&billing.CreditLedgerEntry{}).Where("user_id = ? AND idempotency_key = ?", user.ID, "signup_bonus:"+user.ID).Count(&ledgerCount).Error; err != nil {
		t.Fatalf("count signup bonus ledger entries: %v", err)
	}
	if ledgerCount != 1 {
		t.Fatalf("signup bonus ledger count = %d, want 1", ledgerCount)
	}
}

func TestVerifyEmailOTPRejectsInvalidCode(t *testing.T) {
	router, db := testAuthRouter(t)
	signupVerificationCode(t, router, "ada@example.com")

	w := authJSON(t, router, http.MethodPost, "/api/auth/email-otp/verify-email", `{
		"email":"ada@example.com",
		"otp":"000000"
	}`, nil)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d want 400 response=%s", w.Code, w.Body.String())
	}

	var user auth.User
	if err := db.First(&user, "email = ?", "ada@example.com").Error; err != nil {
		t.Fatalf("load user: %v", err)
	}
	if user.EmailVerified {
		t.Fatal("user should remain unverified")
	}

	var bucketCount int64
	if err := db.Model(&billing.CreditBucket{}).Where("user_id = ?", user.ID).Count(&bucketCount).Error; err != nil {
		t.Fatalf("count credit buckets: %v", err)
	}
	if bucketCount != 0 {
		t.Fatalf("credit bucket count = %d, want 0", bucketCount)
	}
}

func TestVerifyEmailOTPRejectsMissingUserAsInvalidCode(t *testing.T) {
	router, _ := testAuthRouter(t)

	w := authJSON(t, router, http.MethodPost, "/api/auth/email-otp/verify-email", `{
		"email":"missing@example.com",
		"otp":"123456"
	}`, nil)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d want 400 response=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "invalid verification code") {
		t.Fatalf("expected invalid verification code error, got %s", w.Body.String())
	}
}

func TestVerifyEmailOTPRejectsExpiredCode(t *testing.T) {
	router, db := testAuthRouter(t)
	code := signupVerificationCode(t, router, "ada@example.com")
	identifier := verificationIdentifier("ada@example.com")
	if err := db.Model(&auth.Verification{}).
		Where("identifier = ?", identifier).
		Update("expires_at", time.Now().UTC().Add(-time.Minute)).Error; err != nil {
		t.Fatalf("expire verification: %v", err)
	}

	w := authJSON(t, router, http.MethodPost, "/api/auth/email-otp/verify-email", `{
		"email":"ada@example.com",
		"otp":"`+code+`"
	}`, nil)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d want 400 response=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "invalid verification code") {
		t.Fatalf("expected invalid verification code error, got %s", w.Body.String())
	}

	var verificationCount int64
	if err := db.Model(&auth.Verification{}).Where("identifier = ?", identifier).Count(&verificationCount).Error; err != nil {
		t.Fatalf("count verifications: %v", err)
	}
	if verificationCount != 0 {
		t.Fatalf("verification count = %d, want 0", verificationCount)
	}

	var user auth.User
	if err := db.First(&user, "email = ?", "ada@example.com").Error; err != nil {
		t.Fatalf("load user: %v", err)
	}
	var bucketCount int64
	if err := db.Model(&billing.CreditBucket{}).Where("user_id = ?", user.ID).Count(&bucketCount).Error; err != nil {
		t.Fatalf("count credit buckets: %v", err)
	}
	if bucketCount != 0 {
		t.Fatalf("credit bucket count = %d, want 0", bucketCount)
	}
}

func TestVerifyEmailOTPRejectsInvalidInput(t *testing.T) {
	router, _ := testAuthRouter(t)
	overlongEmail := strings.Repeat("a", 309) + "@example.com"

	cases := []string{
		`{"email":"","otp":"123456"}`,
		`{"email":"not-email","otp":"123456"}`,
		`{"email":"` + overlongEmail + `","otp":"123456"}`,
		`{"email":"ada@example.com","otp":""}`,
		`{"email":"ada@example.com","otp":"   "}`,
	}
	for _, body := range cases {
		w := authJSON(t, router, http.MethodPost, "/api/auth/email-otp/verify-email", body, nil)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("body %s status = %d want 400 response=%s", body, w.Code, w.Body.String())
		}
	}
}

func TestVerifyEmailOTPConsumeRequiresCurrentHash(t *testing.T) {
	router, db := testAuthRouter(t)
	signupVerificationCode(t, router, "ada@example.com")
	identifier := verificationIdentifier("ada@example.com")

	var staleVerification auth.Verification
	if err := db.First(&staleVerification, "identifier = ?", identifier).Error; err != nil {
		t.Fatalf("load verification: %v", err)
	}
	rotatedHash := authcrypto.HashCode(testAuthSecret, identifier, "111111")
	if err := db.Model(&auth.Verification{}).
		Where("id = ?", staleVerification.ID).
		Update("value", rotatedHash).Error; err != nil {
		t.Fatalf("rotate verification hash: %v", err)
	}

	var consumed bool
	if err := db.Transaction(func(tx *gorm.DB) error {
		var err error
		consumed, err = consumeEmailVerification(tx, staleVerification, identifier)
		return err
	}); err != nil {
		t.Fatalf("consume verification: %v", err)
	}
	if consumed {
		t.Fatal("stale verification hash consumed rotated row")
	}

	var verification auth.Verification
	if err := db.First(&verification, "id = ?", staleVerification.ID).Error; err != nil {
		t.Fatalf("rotated verification should remain: %v", err)
	}
	if verification.Value != rotatedHash {
		t.Fatalf("verification value = %q, want rotated hash", verification.Value)
	}
}

func TestVerifyEmailOTPExpiredCleanupRequiresCurrentHash(t *testing.T) {
	router, db := testAuthRouter(t)
	signupVerificationCode(t, router, "ada@example.com")
	identifier := verificationIdentifier("ada@example.com")

	var staleVerification auth.Verification
	if err := db.First(&staleVerification, "identifier = ?", identifier).Error; err != nil {
		t.Fatalf("load verification: %v", err)
	}
	staleVerification.ExpiresAt = time.Now().UTC().Add(-time.Minute)
	if err := db.Model(&auth.Verification{}).
		Where("id = ?", staleVerification.ID).
		Updates(map[string]any{
			"value":      authcrypto.HashCode(testAuthSecret, identifier, "222222"),
			"expires_at": time.Now().UTC().Add(time.Minute),
		}).Error; err != nil {
		t.Fatalf("rotate verification: %v", err)
	}

	var consumed bool
	if err := db.Transaction(func(tx *gorm.DB) error {
		var err error
		consumed, err = consumeEmailVerification(tx, staleVerification, identifier)
		return err
	}); err != nil {
		t.Fatalf("cleanup verification: %v", err)
	}
	if consumed {
		t.Fatal("stale expired verification consumed rotated row")
	}

	var verification auth.Verification
	if err := db.First(&verification, "id = ?", staleVerification.ID).Error; err != nil {
		t.Fatalf("rotated verification should remain: %v", err)
	}
	if !verification.ExpiresAt.After(time.Now().UTC()) {
		t.Fatalf("verification expires_at = %v, want future", verification.ExpiresAt)
	}
}

func TestRequestPasswordResetCreatesTrustedTokenForExistingUser(t *testing.T) {
	router, db := testAuthRouter(t)
	hash, err := auth.HashPassword("oldpassword123")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	user := auth.User{Name: "Ada", Email: "ada@example.com", EmailVerified: true}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	if err := db.Create(&auth.AuthAccount{
		AccountID:  user.ID,
		ProviderID: auth.CredentialProviderID,
		UserID:     user.ID,
		Password:   hash,
	}).Error; err != nil {
		t.Fatalf("create account: %v", err)
	}

	w := authJSON(t, router, http.MethodPost, "/api/auth/password-reset/request", `{
		"email":" ADA@EXAMPLE.COM "
	}`, map[string]string{authDeliveryHeader: testDeliveryToken})
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAuthResponse[struct {
		OK             bool   `json:"ok"`
		ResetToken     string `json:"resetToken"`
		ResetExpiresAt string `json:"resetExpiresAt"`
	}](t, w)
	if !got.OK || got.ResetToken == "" || got.ResetExpiresAt == "" {
		t.Fatalf("unexpected reset response: %#v", got)
	}
	if strings.ContainsAny(got.ResetToken, "+/=") {
		t.Fatalf("reset token is not URL-safe: %q", got.ResetToken)
	}

	identifier := "password-reset:ada@example.com"
	var verification auth.Verification
	if err := db.First(&verification, "identifier = ?", identifier).Error; err != nil {
		t.Fatalf("load password reset verification: %v", err)
	}
	if verification.Value == got.ResetToken || strings.Contains(verification.Value, got.ResetToken) {
		t.Fatalf("reset token stored unsafely: %#v", verification)
	}
	if !authcrypto.VerifyCode(testAuthSecret, identifier, got.ResetToken, verification.Value) {
		t.Fatal("stored reset token hash does not verify")
	}
}

func TestRequestPasswordResetMissingUserReturnsOKWithoutTokenOrRow(t *testing.T) {
	router, db := testAuthRouter(t)

	w := authJSON(t, router, http.MethodPost, "/api/auth/password-reset/request", `{
		"email":"missing@example.com"
	}`, map[string]string{authDeliveryHeader: testDeliveryToken})
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAuthResponse[struct {
		OK bool `json:"ok"`
	}](t, w)
	if !got.OK {
		t.Fatalf("ok = %v, want true", got.OK)
	}
	if bytes.Contains(w.Body.Bytes(), []byte("resetToken")) {
		t.Fatalf("missing user response leaked reset token: %s", w.Body.String())
	}
	var count int64
	if err := db.Model(&auth.Verification{}).Where("identifier = ?", "password-reset:missing@example.com").Count(&count).Error; err != nil {
		t.Fatalf("count reset verifications: %v", err)
	}
	if count != 0 {
		t.Fatalf("reset verification count = %d, want 0", count)
	}
}

func TestRequestPasswordResetDoesNotLeakTokenWithoutTrustedHeader(t *testing.T) {
	router, db := testAuthRouter(t)
	user := auth.User{Name: "Ada", Email: "ada@example.com", EmailVerified: true}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}

	w := authJSON(t, router, http.MethodPost, "/api/auth/password-reset/request", `{
		"email":"ada@example.com"
	}`, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAuthResponse[struct {
		OK bool `json:"ok"`
	}](t, w)
	if !got.OK {
		t.Fatalf("ok = %v, want true", got.OK)
	}
	if bytes.Contains(w.Body.Bytes(), []byte("resetToken")) {
		t.Fatalf("untrusted response leaked reset token: %s", w.Body.String())
	}
}

func TestConfirmPasswordResetResetsPasswordConsumesTokenAndInvalidatesSessions(t *testing.T) {
	router, db := testAuthRouter(t)
	cookie := signUpVerifyAndSignIn(t, router)

	request := authJSON(t, router, http.MethodPost, "/api/auth/password-reset/request", `{
		"email":"ada@example.com"
	}`, map[string]string{authDeliveryHeader: testDeliveryToken})
	if request.Code != http.StatusOK {
		t.Fatalf("request status = %d body=%s", request.Code, request.Body.String())
	}
	requestResponse := decodeAuthResponse[struct {
		ResetToken string `json:"resetToken"`
	}](t, request)
	if requestResponse.ResetToken == "" {
		t.Fatalf("missing reset token: %s", request.Body.String())
	}

	w := authJSON(t, router, http.MethodPost, "/api/auth/password-reset/confirm", `{
		"email":" ADA@EXAMPLE.COM ",
		"token":" `+requestResponse.ResetToken+` ",
		"password":"newpassword123"
	}`, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAuthResponse[struct {
		OK bool `json:"ok"`
	}](t, w)
	if !got.OK {
		t.Fatalf("ok = %v, want true", got.OK)
	}

	var user auth.User
	if err := db.First(&user, "email = ?", "ada@example.com").Error; err != nil {
		t.Fatalf("load user: %v", err)
	}
	if !user.EmailVerified {
		t.Fatal("user email_verified = false, want true")
	}
	var account auth.AuthAccount
	if err := db.First(&account, "user_id = ? AND provider_id = ?", user.ID, auth.CredentialProviderID).Error; err != nil {
		t.Fatalf("load account: %v", err)
	}
	if !auth.VerifyPassword("newpassword123", account.Password) {
		t.Fatal("new password does not verify")
	}
	if auth.VerifyPassword("password1234", account.Password) {
		t.Fatal("old password still verifies")
	}
	var verificationCount int64
	if err := db.Model(&auth.Verification{}).Where("identifier = ?", "password-reset:ada@example.com").Count(&verificationCount).Error; err != nil {
		t.Fatalf("count reset verifications: %v", err)
	}
	if verificationCount != 0 {
		t.Fatalf("reset verification count = %d, want 0", verificationCount)
	}
	var sessionCount int64
	if err := db.Model(&auth.Session{}).Where("token = ?", cookie.Value).Count(&sessionCount).Error; err != nil {
		t.Fatalf("count sessions: %v", err)
	}
	if sessionCount != 0 {
		t.Fatalf("session count = %d, want 0", sessionCount)
	}

	reuse := authJSON(t, router, http.MethodPost, "/api/auth/password-reset/confirm", `{
		"email":"ada@example.com",
		"token":"`+requestResponse.ResetToken+`",
		"password":"anotherpassword123"
	}`, nil)
	if reuse.Code != http.StatusBadRequest {
		t.Fatalf("reuse status = %d want 400 response=%s", reuse.Code, reuse.Body.String())
	}
}

func TestConfirmPasswordResetRejectsExpiredToken(t *testing.T) {
	router, db := testAuthRouter(t)
	signupVerificationCode(t, router, "ada@example.com")

	request := authJSON(t, router, http.MethodPost, "/api/auth/password-reset/request", `{
		"email":"ada@example.com"
	}`, map[string]string{authDeliveryHeader: testDeliveryToken})
	if request.Code != http.StatusOK {
		t.Fatalf("request status = %d body=%s", request.Code, request.Body.String())
	}
	requestResponse := decodeAuthResponse[struct {
		ResetToken string `json:"resetToken"`
	}](t, request)
	identifier := "password-reset:ada@example.com"
	if err := db.Model(&auth.Verification{}).
		Where("identifier = ?", identifier).
		Update("expires_at", time.Now().UTC().Add(-time.Minute)).Error; err != nil {
		t.Fatalf("expire reset token: %v", err)
	}

	w := authJSON(t, router, http.MethodPost, "/api/auth/password-reset/confirm", `{
		"email":"ada@example.com",
		"token":"`+requestResponse.ResetToken+`",
		"password":"newpassword123"
	}`, nil)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d want 400 response=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "invalid password reset token") {
		t.Fatalf("expected invalid reset token error, got %s", w.Body.String())
	}
}

func TestConfirmPasswordResetRejectsInvalidPasswordWithoutConsumingToken(t *testing.T) {
	router, db := testAuthRouter(t)
	signupVerificationCode(t, router, "ada@example.com")

	request := authJSON(t, router, http.MethodPost, "/api/auth/password-reset/request", `{
		"email":"ada@example.com"
	}`, map[string]string{authDeliveryHeader: testDeliveryToken})
	if request.Code != http.StatusOK {
		t.Fatalf("request status = %d body=%s", request.Code, request.Body.String())
	}
	requestResponse := decodeAuthResponse[struct {
		ResetToken string `json:"resetToken"`
	}](t, request)

	w := authJSON(t, router, http.MethodPost, "/api/auth/password-reset/confirm", `{
		"email":"ada@example.com",
		"token":"`+requestResponse.ResetToken+`",
		"password":"short"
	}`, nil)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d want 400 response=%s", w.Code, w.Body.String())
	}
	var verificationCount int64
	if err := db.Model(&auth.Verification{}).Where("identifier = ?", "password-reset:ada@example.com").Count(&verificationCount).Error; err != nil {
		t.Fatalf("count reset verifications: %v", err)
	}
	if verificationCount != 1 {
		t.Fatalf("reset verification count = %d, want 1", verificationCount)
	}
}

func TestTrustedAuthDeliveryRequestRequiresConfiguredExactHeader(t *testing.T) {
	tests := []struct {
		name        string
		configToken string
		headerToken string
		want        bool
	}{
		{name: "configured exact match", configToken: testDeliveryToken, headerToken: testDeliveryToken, want: true},
		{name: "configured missing header", configToken: testDeliveryToken, headerToken: "", want: false},
		{name: "configured mismatched header", configToken: testDeliveryToken, headerToken: testDeliveryToken + " ", want: false},
		{name: "empty config rejects matching empty header", configToken: "", headerToken: "", want: false},
		{name: "empty config rejects non-empty header", configToken: "", headerToken: testDeliveryToken, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = newTestRequest(http.MethodPost, "/api/auth/test", nil)
			if tt.headerToken != "" {
				c.Request.Header.Set(authDeliveryHeader, tt.headerToken)
			}

			got := (&Handler{AuthDeliveryToken: tt.configToken}).trustedAuthDeliveryRequest(c)
			if got != tt.want {
				t.Fatalf("trustedAuthDeliveryRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func signupVerificationCode(t *testing.T, router http.Handler, email string) string {
	t.Helper()
	signup := authJSON(t, router, http.MethodPost, "/api/auth/sign-up/email", `{
		"name":"Ada",
		"email":"`+email+`",
		"password":"password1234"
	}`, map[string]string{authDeliveryHeader: testDeliveryToken})
	if signup.Code != http.StatusOK {
		t.Fatalf("signup status = %d body=%s", signup.Code, signup.Body.String())
	}
	got := decodeAuthResponse[struct {
		VerificationCode string `json:"verificationCode"`
	}](t, signup)
	if got.VerificationCode == "" {
		t.Fatalf("signup did not return verification code: %s", signup.Body.String())
	}
	return got.VerificationCode
}

func verifySignupEmail(t *testing.T, router http.Handler, email string, code string) {
	t.Helper()
	w := authJSON(t, router, http.MethodPost, "/api/auth/email-otp/verify-email", `{
		"email":"`+email+`",
		"otp":"`+code+`"
	}`, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("verify status = %d body=%s", w.Code, w.Body.String())
	}
}

func signUpVerifyAndSignIn(t *testing.T, router http.Handler) *http.Cookie {
	t.Helper()
	code := signupVerificationCode(t, router, "ada@example.com")
	verifySignupEmail(t, router, "ada@example.com", code)

	w := authJSON(t, router, http.MethodPost, "/api/auth/sign-in/email", `{
		"email":"ada@example.com",
		"password":"password1234",
		"rememberMe":true
	}`, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("sign-in status = %d body=%s", w.Code, w.Body.String())
	}
	return authCookieFromRecorder(t, w)
}

func startGoogleOAuthState(t *testing.T, router http.Handler, redirectURI string) string {
	t.Helper()
	w := authJSON(t, router, http.MethodPost, "/api/auth/oauth/google/start", `{
		"redirectURI":"`+redirectURI+`"
	}`, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("google oauth start status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAuthResponse[googleOAuthStartResponse](t, w)
	if got.State == "" {
		t.Fatalf("google oauth start returned empty state: %#v", got)
	}
	return got.State
}

func startGoogleAccountLinkState(t *testing.T, router http.Handler, redirectURI string, cookie *http.Cookie) string {
	t.Helper()
	w := authJSONWithCookie(t, router, http.MethodPost, "/api/auth/accounts/google/start", `{
		"redirectURI":"`+redirectURI+`"
	}`, cookie)
	if w.Code != http.StatusOK {
		t.Fatalf("google account link start status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAuthResponse[googleOAuthStartResponse](t, w)
	if got.State == "" {
		t.Fatalf("google account link start returned empty state: %#v", got)
	}
	return got.State
}

func successfulGoogleExchange(t *testing.T) func(context.Context, string, string) (googleOAuthToken, error) {
	t.Helper()
	accessToken := "access-token"
	return func(context.Context, string, string) (googleOAuthToken, error) {
		return googleOAuthToken{AccessToken: &accessToken, IDToken: "id-token-1"}, nil
	}
}

func googleProfileValidator(profile googleIDTokenPayload) func(context.Context, string, string) (googleIDTokenPayload, error) {
	return func(context.Context, string, string) (googleIDTokenPayload, error) {
		return profile, nil
	}
}

func assertGoogleOAuthStateConsumed(t *testing.T, db *gorm.DB, state string) {
	t.Helper()
	var count int64
	if err := db.Model(&auth.Verification{}).Where("identifier = ?", googleOAuthStateIdentifier(state)).Count(&count).Error; err != nil {
		t.Fatalf("count oauth state verification: %v", err)
	}
	if count != 0 {
		t.Fatalf("oauth state verification count = %d, want 0", count)
	}
}

func assertGoogleOAuthLinkStateConsumed(t *testing.T, db *gorm.DB, state string) {
	t.Helper()
	var count int64
	if err := db.Model(&auth.Verification{}).Where("identifier = ?", googleOAuthLinkStateIdentifier(state)).Count(&count).Error; err != nil {
		t.Fatalf("count oauth link state verification: %v", err)
	}
	if count != 0 {
		t.Fatalf("oauth link state verification count = %d, want 0", count)
	}
}

func startGitHubOAuthState(t *testing.T, router http.Handler, redirectURI string) string {
	t.Helper()
	w := authJSON(t, router, http.MethodPost, "/api/auth/oauth/github/start", `{
		"redirectURI":"`+redirectURI+`"
	}`, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("github oauth start status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAuthResponse[githubOAuthStartResponse](t, w)
	if got.State == "" {
		t.Fatalf("github oauth start returned empty state: %#v", got)
	}
	return got.State
}

func startGitHubAccountLinkState(t *testing.T, router http.Handler, redirectURI string, cookie *http.Cookie) string {
	t.Helper()
	w := authJSONWithCookie(t, router, http.MethodPost, "/api/auth/accounts/github/start", `{
		"redirectURI":"`+redirectURI+`"
	}`, cookie)
	if w.Code != http.StatusOK {
		t.Fatalf("github account link start status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAuthResponse[githubOAuthStartResponse](t, w)
	if got.State == "" {
		t.Fatalf("github account link start returned empty state: %#v", got)
	}
	return got.State
}

func successfulGitHubExchange(t *testing.T) func(context.Context, string, string) (githubOAuthToken, error) {
	t.Helper()
	accessToken := "github-access-token"
	return func(context.Context, string, string) (githubOAuthToken, error) {
		return githubOAuthToken{AccessToken: &accessToken}, nil
	}
}

func githubProfileFetcher(profile githubOAuthProfile) func(context.Context, string) (githubOAuthProfile, error) {
	return func(context.Context, string) (githubOAuthProfile, error) {
		return profile, nil
	}
}

func assertGitHubOAuthStateConsumed(t *testing.T, db *gorm.DB, state string) {
	t.Helper()
	var count int64
	if err := db.Model(&auth.Verification{}).Where("identifier = ?", githubOAuthStateIdentifier(state)).Count(&count).Error; err != nil {
		t.Fatalf("count oauth state verification: %v", err)
	}
	if count != 0 {
		t.Fatalf("oauth state verification count = %d, want 0", count)
	}
}

func assertGitHubOAuthLinkStateConsumed(t *testing.T, db *gorm.DB, state string) {
	t.Helper()
	var count int64
	if err := db.Model(&auth.Verification{}).Where("identifier = ?", githubOAuthLinkStateIdentifier(state)).Count(&count).Error; err != nil {
		t.Fatalf("count oauth link state verification: %v", err)
	}
	if count != 0 {
		t.Fatalf("oauth link state verification count = %d, want 0", count)
	}
}

func authJSONWithCookie(t *testing.T, router http.Handler, method string, path string, body string, cookie *http.Cookie) *httptest.ResponseRecorder {
	t.Helper()
	req := newTestRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func authCookieFromRecorder(t *testing.T, w *httptest.ResponseRecorder) *http.Cookie {
	t.Helper()
	for _, cookie := range w.Result().Cookies() {
		if cookie.Name == authSessionCookieName {
			return cookie
		}
	}
	t.Fatalf("missing %s cookie in %v", authSessionCookieName, w.Result().Header.Values("Set-Cookie"))
	return nil
}
