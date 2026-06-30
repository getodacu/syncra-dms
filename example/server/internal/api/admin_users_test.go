package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"ai.ro/syncra/internal/auth"
	"ai.ro/syncra/internal/billing"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type adminUserTestResponse struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	Email         string     `json:"email"`
	EmailVerified bool       `json:"email_verified"`
	Role          string     `json:"role"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	LastLoginAt   *time.Time `json:"last_login_at"`
}

type adminUserListTestResponse struct {
	Users      []adminUserTestResponse `json:"users"`
	NextCursor *string                 `json:"next_cursor"`
}

type adminUserDetailTestResponse struct {
	adminUserTestResponse
	AvailableCredits int                     `json:"available_credits"`
	BillingProfile   *BillingProfileResponse `json:"billing_profile"`
}

func testAdminRouter(t *testing.T) (*gin.Engine, *gorm.DB) {
	return testAdminRouterWithNow(t, nil)
}

func testAdminRouterWithNow(t *testing.T, now func() time.Time) (*gin.Engine, *gorm.DB) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	db := apiPostgresTx(t)
	h := &Handler{
		DB:               db,
		BetterAuthSecret: testAuthSecret,
		InternalAPIToken: testInternalAPIToken,
		AuthCookieSecure: false,
		Now:              now,
	}
	return NewRouter(h), db
}

func createAdminTestUser(t *testing.T, db *gorm.DB, email string, role string) auth.User {
	t.Helper()
	user := auth.User{Name: strings.TrimSuffix(email, "@example.com"), Email: email, EmailVerified: true}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create admin test user: %v", err)
	}
	if role != "" {
		if err := db.Exec(`UPDATE "user" SET role = ? WHERE id = ?`, role, user.ID).Error; err != nil {
			t.Fatalf("set user role: %v", err)
		}
	}
	return user
}

func createAdminTestSession(t *testing.T, db *gorm.DB, user auth.User, token string) *http.Cookie {
	t.Helper()
	session := auth.Session{
		Token:     token,
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().Add(time.Hour),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	if err := db.Create(&session).Error; err != nil {
		t.Fatalf("create admin test session: %v", err)
	}
	return &http.Cookie{Name: authSessionCookieName, Value: token}
}

func adminJSON(t *testing.T, router http.Handler, method string, path string, body string, cookie *http.Cookie) *httptest.ResponseRecorder {
	t.Helper()
	req := newTestRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	if cookie != nil {
		req.AddCookie(cookie)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func decodeAdminResponse[T any](t *testing.T, w *httptest.ResponseRecorder) T {
	t.Helper()
	var out T
	if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode admin response: %v body=%s", err, w.Body.String())
	}
	return out
}

func TestAdminUsersRequireInternalTokenAndAdminSession(t *testing.T) {
	router, db := testAdminRouter(t)
	admin := createAdminTestUser(t, db, "admin-authz@example.com", "admin")
	normal := createAdminTestUser(t, db, "normal-authz@example.com", "user")
	adminCookie := createAdminTestSession(t, db, admin, "admin-authz-session")
	normalCookie := createAdminTestSession(t, db, normal, "normal-authz-session")

	missingTokenReq := httptest.NewRequest(http.MethodGet, "/api/admin/users", nil)
	missingTokenReq.AddCookie(adminCookie)
	missingToken := httptest.NewRecorder()
	router.ServeHTTP(missingToken, missingTokenReq)
	if missingToken.Code != http.StatusUnauthorized {
		t.Fatalf("missing token status = %d body=%s, want 401", missingToken.Code, missingToken.Body.String())
	}

	noSession := adminJSON(t, router, http.MethodGet, "/api/admin/users", "", nil)
	if noSession.Code != http.StatusUnauthorized {
		t.Fatalf("no session status = %d body=%s, want 401", noSession.Code, noSession.Body.String())
	}

	nonAdmin := adminJSON(t, router, http.MethodGet, "/api/admin/users", "", normalCookie)
	if nonAdmin.Code != http.StatusForbidden {
		t.Fatalf("non-admin status = %d body=%s, want 403", nonAdmin.Code, nonAdmin.Body.String())
	}

	ok := adminJSON(t, router, http.MethodGet, "/api/admin/users", "", adminCookie)
	if ok.Code != http.StatusOK {
		t.Fatalf("admin status = %d body=%s, want 200", ok.Code, ok.Body.String())
	}
}

func TestAdminListUsersSearchSortsAndPaginates(t *testing.T) {
	router, db := testAdminRouter(t)
	admin := createAdminTestUser(t, db, "admin-list@example.com", "admin")
	cookie := createAdminTestSession(t, db, admin, "admin-list-session")
	early := createAdminTestUser(t, db, "ada-early@example.com", "user")
	late := createAdminTestUser(t, db, "ada-late@example.com", "user")
	noLogin := createAdminTestUser(t, db, "ada-none@example.com", "user")
	_ = createAdminTestUser(t, db, "grace@example.com", "user")

	base := time.Date(2026, 6, 11, 8, 0, 0, 0, time.UTC)
	if err := db.Model(&early).UpdateColumn("name", "Ada Early").Error; err != nil {
		t.Fatalf("set early name: %v", err)
	}
	if err := db.Model(&late).UpdateColumn("name", "Ada Late").Error; err != nil {
		t.Fatalf("set late name: %v", err)
	}
	if err := db.Model(&noLogin).UpdateColumn("name", "Ada No Login").Error; err != nil {
		t.Fatalf("set no-login name: %v", err)
	}
	if err := db.Exec(`UPDATE "user" SET last_login_at = ? WHERE id = ?`, base, early.ID).Error; err != nil {
		t.Fatalf("set early last login: %v", err)
	}
	if err := db.Exec(`UPDATE "user" SET last_login_at = ? WHERE id = ?`, base.Add(time.Hour), late.ID).Error; err != nil {
		t.Fatalf("set late last login: %v", err)
	}

	w := adminJSON(t, router, http.MethodGet, "/api/admin/users?search=ada&sort=last_login_at&direction=desc&size=2", "", cookie)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAdminResponse[adminUserListTestResponse](t, w)
	if len(got.Users) != 2 {
		t.Fatalf("user count = %d, want 2: %#v", len(got.Users), got.Users)
	}
	if got.Users[0].ID != late.ID || got.Users[1].ID != early.ID {
		t.Fatalf("users order = %#v, want late then early", got.Users)
	}
	if got.NextCursor == nil || *got.NextCursor == "" {
		t.Fatalf("next_cursor = %#v, want non-empty cursor", got.NextCursor)
	}

	next := adminJSON(t, router, http.MethodGet, "/api/admin/users?search=ada&sort=last_login_at&direction=desc&size=2&cursor="+*got.NextCursor, "", cookie)
	if next.Code != http.StatusOK {
		t.Fatalf("next status = %d body=%s", next.Code, next.Body.String())
	}
	gotNext := decodeAdminResponse[adminUserListTestResponse](t, next)
	if len(gotNext.Users) != 1 || gotNext.Users[0].ID != noLogin.ID {
		t.Fatalf("next page users = %#v, want no-login user only", gotNext.Users)
	}
	if gotNext.Users[0].LastLoginAt != nil {
		t.Fatalf("no-login last_login_at = %v, want nil", gotNext.Users[0].LastLoginAt)
	}
}

func TestAdminGetAndPatchUserWithoutRoleEscalation(t *testing.T) {
	router, db := testAdminRouter(t)
	admin := createAdminTestUser(t, db, "admin-patch@example.com", "admin")
	cookie := createAdminTestSession(t, db, admin, "admin-patch-session")
	target := createAdminTestUser(t, db, "target-patch@example.com", "user")
	other := createAdminTestUser(t, db, "target-other@example.com", "user")

	rejected := adminJSON(t, router, http.MethodPatch, "/api/admin/users/"+target.ID, `{"role":"admin"}`, cookie)
	if rejected.Code != http.StatusBadRequest {
		t.Fatalf("role patch status = %d body=%s, want 400", rejected.Code, rejected.Body.String())
	}

	valid := adminJSON(t, router, http.MethodPatch, "/api/admin/users/"+target.ID, `{
		"name":" Target Updated ",
		"email":" TARGET.UPDATED@EXAMPLE.COM "
	}`, cookie)
	if valid.Code != http.StatusOK {
		t.Fatalf("valid patch status = %d body=%s", valid.Code, valid.Body.String())
	}
	got := decodeAdminResponse[adminUserTestResponse](t, valid)
	if got.Name != "Target Updated" || got.Email != "target.updated@example.com" || got.Role != "user" || !got.EmailVerified {
		t.Fatalf("unexpected patched user: %#v", got)
	}

	duplicate := adminJSON(t, router, http.MethodPatch, "/api/admin/users/"+target.ID, `{"email":"`+other.Email+`"}`, cookie)
	if duplicate.Code != http.StatusConflict {
		t.Fatalf("duplicate email status = %d body=%s, want 409", duplicate.Code, duplicate.Body.String())
	}

	var stored struct {
		Role          string `gorm:"column:role"`
		EmailVerified bool   `gorm:"column:email_verified"`
	}
	if err := db.Raw(`SELECT role, email_verified FROM "user" WHERE id = ?`, target.ID).Scan(&stored).Error; err != nil {
		t.Fatalf("load stored target: %v", err)
	}
	if stored.Role != "user" || !stored.EmailVerified {
		t.Fatalf("stored role/verification = %#v, want user and verified", stored)
	}
}

func TestAdminGetUserIncludesAvailableCredits(t *testing.T) {
	router, db := testAdminRouter(t)
	admin := createAdminTestUser(t, db, "admin-balance-get@example.com", "admin")
	cookie := createAdminTestSession(t, db, admin, "admin-balance-get-session")
	target := createAdminTestUser(t, db, "target-balance-get@example.com", "user")
	now := time.Now().UTC().Add(-time.Hour)
	bucket := billing.CreditBucket{
		UserID:           target.ID,
		SourceType:       billing.CreditSourceAdjustment,
		CreditsGranted:   375,
		CreditsRemaining: 375,
		ValidFrom:        now,
	}
	if err := db.Create(&bucket).Error; err != nil {
		t.Fatalf("create target credit bucket: %v", err)
	}

	w := adminJSON(t, router, http.MethodGet, "/api/admin/users/"+target.ID, "", cookie)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s, want 200", w.Code, w.Body.String())
	}
	got := decodeAdminResponse[adminUserDetailTestResponse](t, w)
	if got.ID != target.ID || got.AvailableCredits != 375 {
		t.Fatalf("unexpected user detail response: %#v", got)
	}
}

func TestAdminAdjustUserBalanceAddsAndSubtractsCredits(t *testing.T) {
	router, db := testAdminRouter(t)
	admin := createAdminTestUser(t, db, "admin-balance-adjust@example.com", "admin")
	cookie := createAdminTestSession(t, db, admin, "admin-balance-adjust-session")
	target := createAdminTestUser(t, db, "target-balance-adjust@example.com", "user")

	add := adminJSON(t, router, http.MethodPost, "/api/admin/users/"+target.ID+"/balance-adjustment", `{"credits_delta":250}`, cookie)
	if add.Code != http.StatusOK {
		t.Fatalf("add status = %d body=%s, want 200", add.Code, add.Body.String())
	}
	added := decodeAdminResponse[adminUserDetailTestResponse](t, add)
	if added.AvailableCredits != 250 {
		t.Fatalf("available credits after add = %d, want 250", added.AvailableCredits)
	}

	subtract := adminJSON(t, router, http.MethodPost, "/api/admin/users/"+target.ID+"/balance-adjustment", `{"credits_delta":-75}`, cookie)
	if subtract.Code != http.StatusOK {
		t.Fatalf("subtract status = %d body=%s, want 200", subtract.Code, subtract.Body.String())
	}
	subtracted := decodeAdminResponse[adminUserDetailTestResponse](t, subtract)
	if subtracted.AvailableCredits != 175 {
		t.Fatalf("available credits after subtract = %d, want 175", subtracted.AvailableCredits)
	}

	var entries []billing.CreditLedgerEntry
	if err := db.Where("user_id = ? AND entry_type = ?", target.ID, billing.CreditLedgerEntryAdjustment).
		Order("credits_delta desc").
		Find(&entries).Error; err != nil {
		t.Fatalf("load adjustment entries: %v", err)
	}
	if len(entries) != 2 || entries[0].CreditsDelta != 250 || entries[1].CreditsDelta != -75 {
		t.Fatalf("adjustment entries = %#v, want +250 and -75", entries)
	}
}

func TestAdminAdjustUserBalanceRejectsInsufficientCredits(t *testing.T) {
	router, db := testAdminRouter(t)
	admin := createAdminTestUser(t, db, "admin-balance-insufficient@example.com", "admin")
	cookie := createAdminTestSession(t, db, admin, "admin-balance-insufficient-session")
	target := createAdminTestUser(t, db, "target-balance-insufficient@example.com", "user")
	now := time.Now().UTC().Add(-time.Hour)
	bucket := billing.CreditBucket{
		UserID:           target.ID,
		SourceType:       billing.CreditSourceAdjustment,
		CreditsGranted:   50,
		CreditsRemaining: 50,
		ValidFrom:        now,
	}
	if err := db.Create(&bucket).Error; err != nil {
		t.Fatalf("create target credit bucket: %v", err)
	}

	w := adminJSON(t, router, http.MethodPost, "/api/admin/users/"+target.ID+"/balance-adjustment", `{"credits_delta":-51}`, cookie)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s, want 400", w.Code, w.Body.String())
	}
	var reloaded billing.CreditBucket
	if err := db.First(&reloaded, "id = ?", bucket.ID).Error; err != nil {
		t.Fatalf("reload credit bucket: %v", err)
	}
	if reloaded.CreditsRemaining != 50 {
		t.Fatalf("credits remaining = %d, want 50", reloaded.CreditsRemaining)
	}
}

func TestAdminSetPasswordHashesPasswordAndInvalidatesTargetSessions(t *testing.T) {
	router, db := testAdminRouter(t)
	admin := createAdminTestUser(t, db, "admin-password@example.com", "admin")
	adminCookie := createAdminTestSession(t, db, admin, "admin-password-session")
	target := createAdminTestUser(t, db, "target-password@example.com", "user")
	oldHash, err := auth.HashPassword("oldpassword123")
	if err != nil {
		t.Fatalf("hash old password: %v", err)
	}
	account := auth.AuthAccount{
		AccountID:  target.ID,
		ProviderID: auth.CredentialProviderID,
		UserID:     target.ID,
		Password:   oldHash,
	}
	if err := db.Create(&account).Error; err != nil {
		t.Fatalf("create target account: %v", err)
	}
	createAdminTestSession(t, db, target, "target-password-session-1")
	createAdminTestSession(t, db, target, "target-password-session-2")

	w := adminJSON(t, router, http.MethodPost, "/api/admin/users/"+target.ID+"/password", `{"password":"newpassword123"}`, adminCookie)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s, want 200", w.Code, w.Body.String())
	}

	var updated auth.AuthAccount
	if err := db.First(&updated, "user_id = ? AND provider_id = ?", target.ID, auth.CredentialProviderID).Error; err != nil {
		t.Fatalf("load updated account: %v", err)
	}
	if updated.Password == "" || updated.Password == "newpassword123" || updated.Password == oldHash {
		t.Fatalf("unsafe or unchanged password hash: %#v", updated)
	}
	if !auth.VerifyPassword("newpassword123", updated.Password) {
		t.Fatal("new password does not verify")
	}
	if auth.VerifyPassword("oldpassword123", updated.Password) {
		t.Fatal("old password still verifies after reset")
	}

	var targetSessions int64
	if err := db.Model(&auth.Session{}).Where("user_id = ?", target.ID).Count(&targetSessions).Error; err != nil {
		t.Fatalf("count target sessions: %v", err)
	}
	if targetSessions != 0 {
		t.Fatalf("target sessions = %d, want 0", targetSessions)
	}
	var adminSessions int64
	if err := db.Model(&auth.Session{}).Where("user_id = ?", admin.ID).Count(&adminSessions).Error; err != nil {
		t.Fatalf("count admin sessions: %v", err)
	}
	if adminSessions != 1 {
		t.Fatalf("admin sessions = %d, want 1", adminSessions)
	}
}

func TestAdminUpsertsBillingProfileForTargetUser(t *testing.T) {
	router, db := testAdminRouter(t)
	admin := createAdminTestUser(t, db, "admin-billing@example.com", "admin")
	cookie := createAdminTestSession(t, db, admin, "admin-billing-session")
	target := createAdminTestUser(t, db, "target-billing@example.com", "user")

	body := `{
		"entity_type":"company",
		"billing_name":"Syncra Test SRL",
		"billing_email":"billing@example.com",
		"country_code":"RO",
		"address_line1":"Str Test 1",
		"address_line2":"Suite 2",
		"city":"Bucharest",
		"region":"B",
		"postal_code":"010101",
		"fiscal_code":"RO123",
		"registration_number":"J40/1/2026"
	}`
	w := adminJSON(t, router, http.MethodPut, "/api/admin/users/"+target.ID+"/billing-profile", body, cookie)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s, want 200", w.Code, w.Body.String())
	}
	got := decodeAdminResponse[BillingProfileResponse](t, w)
	if string(got.UserID) != target.ID || got.BillingName != "Syncra Test SRL" || got.CountryCode != "RO" {
		t.Fatalf("unexpected billing profile response: %#v", got)
	}

	var profile billing.BillingProfile
	if err := db.First(&profile, "user_id = ?", target.ID).Error; err != nil {
		t.Fatalf("load billing profile: %v", err)
	}
	if profile.UserID != target.ID || profile.BillingName != "Syncra Test SRL" {
		t.Fatalf("unexpected stored billing profile: %#v", profile)
	}

	invalid := adminJSON(t, router, http.MethodPut, "/api/admin/users/"+target.ID+"/billing-profile", `{"entity_type":"company"}`, cookie)
	if invalid.Code != http.StatusBadRequest {
		t.Fatalf("invalid profile status = %d body=%s, want 400", invalid.Code, invalid.Body.String())
	}
}

func TestAdminStartImpersonationSetsEffectiveUserAndAudit(t *testing.T) {
	now := time.Date(2026, 6, 14, 10, 0, 0, 0, time.UTC)
	router, db := testAdminRouterWithNow(t, func() time.Time { return now })
	admin := createAdminTestUser(t, db, "admin-impersonate@example.com", "admin")
	target := createAdminTestUser(t, db, "target-impersonate@example.com", "user")
	cookie := createAdminTestSession(t, db, admin, "admin-impersonate-session")

	w := adminJSON(t, router, http.MethodPost, "/api/admin/users/"+target.ID+"/impersonation", "", cookie)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s, want 200", w.Code, w.Body.String())
	}
	got := decodeAdminResponse[getSessionResponse](t, w)
	if got.Session.UserID != admin.ID || got.User.ID != target.ID || got.User.Email != target.Email {
		t.Fatalf("unexpected impersonation response: %#v", got)
	}
	if got.Impersonation == nil || got.Impersonation.AdminUser.ID != admin.ID || got.Impersonation.TargetUser.ID != target.ID || got.Impersonation.StartedAt != now.Format(time.RFC3339Nano) {
		t.Fatalf("unexpected impersonation metadata: %#v", got.Impersonation)
	}

	var session auth.Session
	if err := db.First(&session, "token = ?", cookie.Value).Error; err != nil {
		t.Fatalf("load session: %v", err)
	}
	if session.ImpersonatedUserID == nil || *session.ImpersonatedUserID != target.ID {
		t.Fatalf("impersonated_user_id = %#v, want %s", session.ImpersonatedUserID, target.ID)
	}
	if session.ImpersonationStartedAt == nil || !session.ImpersonationStartedAt.Equal(now) {
		t.Fatalf("impersonation_started_at = %v, want %v", session.ImpersonationStartedAt, now)
	}

	var events []auth.AdminImpersonationEvent
	if err := db.Find(&events).Error; err != nil {
		t.Fatalf("load impersonation events: %v", err)
	}
	if len(events) != 1 || events[0].EventType != auth.AdminImpersonationEventStart || events[0].AdminUserID != admin.ID || events[0].TargetUserID != target.ID {
		t.Fatalf("unexpected audit events: %#v", events)
	}
}

func TestAdminStartImpersonationRejectsInvalidRequests(t *testing.T) {
	router, db := testAdminRouter(t)
	admin := createAdminTestUser(t, db, "admin-impersonate-reject@example.com", "admin")
	normal := createAdminTestUser(t, db, "normal-impersonate-reject@example.com", "user")
	target := createAdminTestUser(t, db, "target-impersonate-reject@example.com", "user")
	otherTarget := createAdminTestUser(t, db, "target-other-impersonate-reject@example.com", "user")
	adminTarget := createAdminTestUser(t, db, "admin-target-impersonate-reject@example.com", "admin")
	adminCookie := createAdminTestSession(t, db, admin, "admin-impersonate-reject-session")
	normalCookie := createAdminTestSession(t, db, normal, "normal-impersonate-reject-session")

	cases := []struct {
		name   string
		userID string
		cookie *http.Cookie
		want   int
	}{
		{name: "no session", userID: target.ID, cookie: nil, want: http.StatusUnauthorized},
		{name: "non admin", userID: target.ID, cookie: normalCookie, want: http.StatusForbidden},
		{name: "self", userID: admin.ID, cookie: adminCookie, want: http.StatusBadRequest},
		{name: "admin target", userID: adminTarget.ID, cookie: adminCookie, want: http.StatusBadRequest},
		{name: "missing target", userID: "00000000-0000-0000-0000-000000000001", cookie: adminCookie, want: http.StatusNotFound},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := adminJSON(t, router, http.MethodPost, "/api/admin/users/"+tc.userID+"/impersonation", "", tc.cookie)
			if w.Code != tc.want {
				t.Fatalf("status = %d body=%s, want %d", w.Code, w.Body.String(), tc.want)
			}
		})
	}

	start := adminJSON(t, router, http.MethodPost, "/api/admin/users/"+target.ID+"/impersonation", "", adminCookie)
	if start.Code != http.StatusOK {
		t.Fatalf("start status = %d body=%s, want 200", start.Code, start.Body.String())
	}
	switching := adminJSON(t, router, http.MethodPost, "/api/admin/users/"+otherTarget.ID+"/impersonation", "", adminCookie)
	if switching.Code != http.StatusConflict {
		t.Fatalf("switching status = %d body=%s, want 409", switching.Code, switching.Body.String())
	}
}

func TestGetSessionAndPatchAuthUserUseImpersonatedEffectiveUser(t *testing.T) {
	router, db := testAdminRouter(t)
	admin := createAdminTestUser(t, db, "admin-effective@example.com", "admin")
	target := createAdminTestUser(t, db, "target-effective@example.com", "user")
	cookie := createAdminTestSession(t, db, admin, "admin-effective-session")

	start := adminJSON(t, router, http.MethodPost, "/api/admin/users/"+target.ID+"/impersonation", "", cookie)
	if start.Code != http.StatusOK {
		t.Fatalf("start status = %d body=%s, want 200", start.Code, start.Body.String())
	}

	req := newTestRequest(http.MethodGet, "/api/auth/get-session", nil)
	req.AddCookie(cookie)
	sessionResponse := httptest.NewRecorder()
	router.ServeHTTP(sessionResponse, req)
	if sessionResponse.Code != http.StatusOK {
		t.Fatalf("get-session status = %d body=%s, want 200", sessionResponse.Code, sessionResponse.Body.String())
	}
	got := decodeAdminResponse[getSessionResponse](t, sessionResponse)
	if got.Session.UserID != admin.ID || got.User.ID != target.ID || got.Impersonation == nil {
		t.Fatalf("unexpected get-session response: %#v", got)
	}

	patch := newTestRequest(http.MethodPatch, "/api/auth/user", strings.NewReader(`{"name":"Target Updated","email":"target.effective.updated@example.com","preferredLanguage":"ro"}`))
	patch.Header.Set("Content-Type", "application/json")
	patch.AddCookie(cookie)
	patchResponse := httptest.NewRecorder()
	router.ServeHTTP(patchResponse, patch)
	if patchResponse.Code != http.StatusOK {
		t.Fatalf("patch status = %d body=%s, want 200", patchResponse.Code, patchResponse.Body.String())
	}
	patched := decodeAdminResponse[authUserResponse](t, patchResponse)
	if patched.ID != target.ID || patched.Name != "Target Updated" || patched.Email != "target.effective.updated@example.com" || patched.PreferredLanguage != "ro" {
		t.Fatalf("unexpected patched user: %#v", patched)
	}

	var reloadedAdmin auth.User
	if err := db.First(&reloadedAdmin, "id = ?", admin.ID).Error; err != nil {
		t.Fatalf("load admin: %v", err)
	}
	if reloadedAdmin.Email != admin.Email || reloadedAdmin.Name != admin.Name || reloadedAdmin.PreferredLanguage != admin.PreferredLanguage {
		t.Fatalf("admin was unexpectedly changed: %#v", reloadedAdmin)
	}

	var reloadedTarget auth.User
	if err := db.First(&reloadedTarget, "id = ?", target.ID).Error; err != nil {
		t.Fatalf("load target: %v", err)
	}
	if reloadedTarget.PreferredLanguage != "ro" {
		t.Fatalf("target preferred language = %q, want ro", reloadedTarget.PreferredLanguage)
	}
}

func TestAdminStopImpersonationClearsSessionWritesAuditAndIsIdempotent(t *testing.T) {
	now := time.Date(2026, 6, 14, 11, 0, 0, 0, time.UTC)
	router, db := testAdminRouterWithNow(t, func() time.Time { return now })
	admin := createAdminTestUser(t, db, "admin-stop-impersonation@example.com", "admin")
	target := createAdminTestUser(t, db, "target-stop-impersonation@example.com", "user")
	cookie := createAdminTestSession(t, db, admin, "admin-stop-impersonation-session")

	start := adminJSON(t, router, http.MethodPost, "/api/admin/users/"+target.ID+"/impersonation", "", cookie)
	if start.Code != http.StatusOK {
		t.Fatalf("start status = %d body=%s, want 200", start.Code, start.Body.String())
	}
	stop := adminJSON(t, router, http.MethodPost, "/api/admin/impersonation/stop", "", cookie)
	if stop.Code != http.StatusOK {
		t.Fatalf("stop status = %d body=%s, want 200", stop.Code, stop.Body.String())
	}
	got := decodeAdminResponse[getSessionResponse](t, stop)
	if got.User.ID != admin.ID || got.Impersonation != nil {
		t.Fatalf("unexpected stop response: %#v", got)
	}

	var session auth.Session
	if err := db.First(&session, "token = ?", cookie.Value).Error; err != nil {
		t.Fatalf("load session: %v", err)
	}
	if session.ImpersonatedUserID != nil || session.ImpersonationStartedAt != nil {
		t.Fatalf("session still impersonating: %#v", session)
	}
	var eventCount int64
	if err := db.Model(&auth.AdminImpersonationEvent{}).Count(&eventCount).Error; err != nil {
		t.Fatalf("count events: %v", err)
	}
	if eventCount != 2 {
		t.Fatalf("event count = %d, want start and stop", eventCount)
	}

	secondStop := adminJSON(t, router, http.MethodPost, "/api/admin/impersonation/stop", "", cookie)
	if secondStop.Code != http.StatusOK {
		t.Fatalf("second stop status = %d body=%s, want 200", secondStop.Code, secondStop.Body.String())
	}
	if err := db.Model(&auth.AdminImpersonationEvent{}).Count(&eventCount).Error; err != nil {
		t.Fatalf("count events after idempotent stop: %v", err)
	}
	if eventCount != 2 {
		t.Fatalf("event count after idempotent stop = %d, want 2", eventCount)
	}
}

func TestGetSessionClearsImpersonationWhenTargetIsNoLongerNormalUser(t *testing.T) {
	router, db := testAdminRouter(t)
	admin := createAdminTestUser(t, db, "admin-stale-impersonation@example.com", "admin")
	target := createAdminTestUser(t, db, "target-stale-impersonation@example.com", "user")
	cookie := createAdminTestSession(t, db, admin, "admin-stale-impersonation-session")
	start := adminJSON(t, router, http.MethodPost, "/api/admin/users/"+target.ID+"/impersonation", "", cookie)
	if start.Code != http.StatusOK {
		t.Fatalf("start status = %d body=%s, want 200", start.Code, start.Body.String())
	}
	if err := db.Exec(`UPDATE "user" SET role = 'admin' WHERE id = ?`, target.ID).Error; err != nil {
		t.Fatalf("promote target: %v", err)
	}

	req := newTestRequest(http.MethodGet, "/api/auth/get-session", nil)
	req.AddCookie(cookie)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s, want 200", w.Code, w.Body.String())
	}
	got := decodeAdminResponse[getSessionResponse](t, w)
	if got.User.ID != admin.ID || got.Impersonation != nil {
		t.Fatalf("unexpected stale impersonation response: %#v", got)
	}
	var session auth.Session
	if err := db.First(&session, "token = ?", cookie.Value).Error; err != nil {
		t.Fatalf("load session: %v", err)
	}
	if session.ImpersonatedUserID != nil || session.ImpersonationStartedAt != nil {
		t.Fatalf("stale impersonation was not cleared: %#v", session)
	}
}
