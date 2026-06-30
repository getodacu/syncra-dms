package auth

import (
	"database/sql"
	"testing"
	"time"

	"ai.ro/syncra/internal/testsupport"
	"gorm.io/gorm"
)

var authModelTestGroup *testsupport.PostgresGroup

func TestAuthModels(t *testing.T) {
	authModelTestGroup = testsupport.OpenPostgresGroup(t, &User{}, &AuthAccount{}, &Session{}, &Verification{}, &APIKey{}, &AdminImpersonationEvent{})
	defer func() { authModelTestGroup = nil }()

	for _, tt := range []struct {
		name string
		fn   func(*testing.T)
	}{
		{name: "AutoMigrateAndPersist", fn: testAuthModelsAutoMigrateAndPersist},
		{name: "UserModelIncludesRoleAndLastLoginColumns", fn: testUserModelIncludesRoleAndLastLoginColumns},
		{name: "UserModelIncludesPreferredLanguageDefaultAndConstraint", fn: testUserModelIncludesPreferredLanguageDefaultAndConstraint},
		{name: "UserRoleRejectsInvalidValues", fn: testUserRoleRejectsInvalidValues},
		{name: "VerificationIdentifierAllowsPrefixedMaxLengthEmail", fn: testVerificationIdentifierAllowsPrefixedMaxLengthEmail},
		{name: "AuthModelConstraints", fn: testAuthModelConstraints},
		{name: "APIKeyRejectsDuplicateKeyHash", fn: testAPIKeyRejectsDuplicateKeyHash},
		{name: "AuthAccountRejectsDuplicateProviderAccountIdentity", fn: testAuthAccountRejectsDuplicateProviderAccountIdentity},
		{name: "SessionRejectsDuplicateToken", fn: testSessionRejectsDuplicateToken},
		{name: "VerificationRejectsDuplicateIdentifier", fn: testVerificationRejectsDuplicateIdentifier},
		{name: "CascadeDeleteUser", fn: testAuthModelsCascadeDeleteUser},
	} {
		t.Run(tt.name, tt.fn)
	}
}

func authModelTx(t *testing.T) *gorm.DB {
	t.Helper()
	if authModelTestGroup != nil {
		return authModelTestGroup.Tx(t)
	}
	return testsupport.OpenPostgresTx(t, &User{}, &AuthAccount{}, &Session{}, &Verification{}, &APIKey{}, &AdminImpersonationEvent{})
}

func testAuthModelsAutoMigrateAndPersist(t *testing.T) {
	db := authModelTx(t)

	user := User{
		Name:          "Ada Lovelace",
		Email:         "ada@example.com",
		EmailVerified: false,
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	if user.ID == "" {
		t.Fatal("user id was not generated")
	}

	account := AuthAccount{
		AccountID:  user.ID,
		ProviderID: CredentialProviderID,
		UserID:     user.ID,
		Password:   "scrypt$v=1$test",
	}
	if err := db.Create(&account).Error; err != nil {
		t.Fatalf("create account: %v", err)
	}

	session := Session{
		Token:     "session-token",
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour),
		UserAgent: "test-agent",
		IPAddress: "127.0.0.1",
	}
	if err := db.Create(&session).Error; err != nil {
		t.Fatalf("create session: %v", err)
	}

	verification := Verification{
		Identifier: "email-verification:ada@example.com",
		Value:      "hmac-value",
		ExpiresAt:  time.Now().Add(5 * time.Minute),
	}
	if err := db.Create(&verification).Error; err != nil {
		t.Fatalf("create verification: %v", err)
	}

	apiKey := APIKey{
		UserID:    user.ID,
		Name:      "CLI",
		KeyHash:   HashAPIKey("test-api-key"),
		KeyPrefix: "test-api",
	}
	if err := db.Create(&apiKey).Error; err != nil {
		t.Fatalf("create api key: %v", err)
	}

	var gotUser User
	if err := db.First(&gotUser, "email = ?", "ada@example.com").Error; err != nil {
		t.Fatalf("load user: %v", err)
	}
	if gotUser.ID != user.ID || gotUser.EmailVerified {
		t.Fatalf("unexpected user: %#v", gotUser)
	}

	var gotAccount AuthAccount
	if err := db.First(&gotAccount, "id = ?", account.ID).Error; err != nil {
		t.Fatalf("load account: %v", err)
	}
	if gotAccount.ProviderID != CredentialProviderID || gotAccount.AccountID != user.ID || gotAccount.UserID != user.ID {
		t.Fatalf("unexpected account: %#v", gotAccount)
	}

	var gotSession Session
	if err := db.First(&gotSession, "token = ?", "session-token").Error; err != nil {
		t.Fatalf("load session: %v", err)
	}
	if gotSession.UserID != user.ID || gotSession.UserAgent != "test-agent" || gotSession.IPAddress != "127.0.0.1" || gotSession.ExpiresAt.IsZero() {
		t.Fatalf("unexpected session: %#v", gotSession)
	}

	var gotVerification Verification
	if err := db.First(&gotVerification, "identifier = ?", "email-verification:ada@example.com").Error; err != nil {
		t.Fatalf("load verification: %v", err)
	}
	if gotVerification.Value != "hmac-value" || gotVerification.ExpiresAt.IsZero() {
		t.Fatalf("unexpected verification: %#v", gotVerification)
	}

	var gotAPIKey APIKey
	if err := db.First(&gotAPIKey, "id = ?", apiKey.ID).Error; err != nil {
		t.Fatalf("load api key: %v", err)
	}
	if gotAPIKey.UserID != user.ID || gotAPIKey.Name != "CLI" || gotAPIKey.KeyHash == "" || gotAPIKey.KeyHash == "test-api-key" || gotAPIKey.KeyPrefix != "test-api" {
		t.Fatalf("unexpected api key: %#v", gotAPIKey)
	}
	if gotAPIKey.ExpiresAt != nil {
		t.Fatalf("expires_at = %v, want nil", gotAPIKey.ExpiresAt)
	}
}

func testUserModelIncludesRoleAndLastLoginColumns(t *testing.T) {
	db := authModelTx(t)

	user := User{Name: "Default Role", Email: "default-role@example.com"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}

	var got struct {
		Role        string       `gorm:"column:role"`
		LastLoginAt sql.NullTime `gorm:"column:last_login_at"`
	}
	if err := db.Raw(`SELECT role, last_login_at FROM "user" WHERE id = ?`, user.ID).Scan(&got).Error; err != nil {
		t.Fatalf("load user role metadata: %v", err)
	}
	if got.Role != "user" {
		t.Fatalf("default role = %q, want user", got.Role)
	}
	if got.LastLoginAt.Valid {
		t.Fatalf("last_login_at valid = true, want null: %v", got.LastLoginAt.Time)
	}

	if err := db.Exec(`UPDATE "user" SET role = 'admin', last_login_at = ? WHERE id = ?`, time.Now().UTC(), user.ID).Error; err != nil {
		t.Fatalf("set admin role and last login: %v", err)
	}
}

func testUserRoleRejectsInvalidValues(t *testing.T) {
	db := authModelTx(t)

	user := User{Name: "Invalid Role", Email: "invalid-role@example.com"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	if err := db.Exec(`UPDATE "user" SET role = 'owner' WHERE id = ?`, user.ID).Error; err == nil {
		t.Fatal("invalid role update succeeded, want check constraint failure")
	}
}

func testUserModelIncludesPreferredLanguageDefaultAndConstraint(t *testing.T) {
	db := authModelTx(t)

	user := User{Name: "Language Default", Email: "language-default@example.com"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}

	var got struct {
		PreferredLanguage string `gorm:"column:preferred_language"`
	}
	if err := db.Raw(`SELECT preferred_language FROM "user" WHERE id = ?`, user.ID).Scan(&got).Error; err != nil {
		t.Fatalf("load user preferred language: %v", err)
	}
	if got.PreferredLanguage != "en" {
		t.Fatalf("preferred_language = %q, want en", got.PreferredLanguage)
	}

	if err := db.Exec(`UPDATE "user" SET preferred_language = 'de' WHERE id = ?`, user.ID).Error; err == nil {
		t.Fatal("invalid preferred language update succeeded, want check constraint failure")
	}
}

func testVerificationIdentifierAllowsPrefixedMaxLengthEmail(t *testing.T) {
	db := authModelTx(t)

	columns, err := db.Migrator().ColumnTypes(&Verification{})
	if err != nil {
		t.Fatalf("load verification column types: %v", err)
	}

	for _, column := range columns {
		if column.Name() != "identifier" {
			continue
		}
		length, ok := column.Length()
		if !ok {
			t.Fatal("verification identifier column length unavailable")
		}
		if length < 512 {
			t.Fatalf("verification identifier column length = %d, want at least 512", length)
		}
		return
	}
	t.Fatal("verification identifier column not found")
}

func testAuthModelConstraints(t *testing.T) {
	db := authModelTx(t)

	first := User{Name: "First", Email: "same@example.com"}
	second := User{Name: "Second", Email: "same@example.com"}
	if err := db.Create(&first).Error; err != nil {
		t.Fatalf("create first user: %v", err)
	}
	if err := db.Create(&second).Error; err == nil {
		t.Fatal("create duplicate email succeeded, want unique constraint failure")
	}
}

func testAPIKeyRejectsDuplicateKeyHash(t *testing.T) {
	db := authModelTx(t)

	user := User{Name: "API Key Owner", Email: "api-key-owner@example.com"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}

	hash := HashAPIKey("duplicate-api-key")
	first := APIKey{UserID: user.ID, Name: "first", KeyHash: hash, KeyPrefix: "duplicat"}
	second := APIKey{UserID: user.ID, Name: "second", KeyHash: hash, KeyPrefix: "duplicat"}
	if err := db.Create(&first).Error; err != nil {
		t.Fatalf("create first api key: %v", err)
	}
	if err := db.Create(&second).Error; err == nil {
		t.Fatal("create duplicate api key hash succeeded, want unique constraint failure")
	}
}

func testAuthAccountRejectsDuplicateProviderAccountIdentity(t *testing.T) {
	db := authModelTx(t)

	user := User{Name: "Account Owner", Email: "account-owner@example.com"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}

	first := AuthAccount{
		AccountID:  user.ID,
		ProviderID: CredentialProviderID,
		UserID:     user.ID,
		Password:   "scrypt$v=1$first",
	}
	second := AuthAccount{
		AccountID:  user.ID,
		ProviderID: CredentialProviderID,
		UserID:     user.ID,
		Password:   "scrypt$v=1$second",
	}
	if err := db.Create(&first).Error; err != nil {
		t.Fatalf("create first account: %v", err)
	}
	if err := db.Create(&second).Error; err == nil {
		t.Fatal("create duplicate provider/account identity succeeded, want unique constraint failure")
	}
}

func testSessionRejectsDuplicateToken(t *testing.T) {
	db := authModelTx(t)

	user := User{Name: "Session Owner", Email: "session-owner@example.com"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}

	first := Session{
		Token:     "duplicate-session-token",
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	second := Session{
		Token:     "duplicate-session-token",
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(2 * time.Hour),
	}
	if err := db.Create(&first).Error; err != nil {
		t.Fatalf("create first session: %v", err)
	}
	if err := db.Create(&second).Error; err == nil {
		t.Fatal("create duplicate session token succeeded, want unique constraint failure")
	}
}

func testVerificationRejectsDuplicateIdentifier(t *testing.T) {
	db := authModelTx(t)

	first := Verification{
		Identifier: "email-verification:duplicate@example.com",
		Value:      "first",
		ExpiresAt:  time.Now().Add(5 * time.Minute),
	}
	second := Verification{
		Identifier: "email-verification:duplicate@example.com",
		Value:      "second",
		ExpiresAt:  time.Now().Add(10 * time.Minute),
	}
	if err := db.Create(&first).Error; err != nil {
		t.Fatalf("create first verification: %v", err)
	}
	if err := db.Create(&second).Error; err == nil {
		t.Fatal("create duplicate verification identifier succeeded, want unique constraint failure")
	}
}

func testAuthModelsCascadeDeleteUser(t *testing.T) {
	db := authModelTx(t)

	user := User{Name: "Cascade Owner", Email: "cascade-owner@example.com"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	account := AuthAccount{
		AccountID:  user.ID,
		ProviderID: CredentialProviderID,
		UserID:     user.ID,
		Password:   "scrypt$v=1$cascade",
	}
	if err := db.Create(&account).Error; err != nil {
		t.Fatalf("create account: %v", err)
	}
	session := Session{
		Token:     "cascade-session-token",
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	if err := db.Create(&session).Error; err != nil {
		t.Fatalf("create session: %v", err)
	}
	apiKey := APIKey{
		UserID:    user.ID,
		Name:      "cascade key",
		KeyHash:   HashAPIKey("cascade-api-key"),
		KeyPrefix: "cascade",
	}
	if err := db.Create(&apiKey).Error; err != nil {
		t.Fatalf("create api key: %v", err)
	}

	if err := db.Delete(&user).Error; err != nil {
		t.Fatalf("delete user: %v", err)
	}

	var accountCount int64
	if err := db.Model(&AuthAccount{}).Where("user_id = ?", user.ID).Count(&accountCount).Error; err != nil {
		t.Fatalf("count accounts: %v", err)
	}
	if accountCount != 0 {
		t.Fatalf("account count after user delete = %d, want 0", accountCount)
	}

	var sessionCount int64
	if err := db.Model(&Session{}).Where("user_id = ?", user.ID).Count(&sessionCount).Error; err != nil {
		t.Fatalf("count sessions: %v", err)
	}
	if sessionCount != 0 {
		t.Fatalf("session count after user delete = %d, want 0", sessionCount)
	}

	var apiKeyCount int64
	if err := db.Model(&APIKey{}).Where("user_id = ?", user.ID).Count(&apiKeyCount).Error; err != nil {
		t.Fatalf("count api keys: %v", err)
	}
	if apiKeyCount != 0 {
		t.Fatalf("api key count after user delete = %d, want 0", apiKeyCount)
	}
}
