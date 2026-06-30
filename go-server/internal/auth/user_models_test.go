package auth

import (
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestAuthModelsPersistBetterAuthCompatibleTables(t *testing.T) {
	db := openAuthModelTestDB(t)
	now := time.Now().UTC()
	user := User{
		Name:          "Ada Lovelace",
		Email:         "ada@example.com",
		EmailVerified: false,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	if user.ID == "" {
		t.Fatal("user ID was not generated")
	}
	if user.Role != UserRoleUser {
		t.Fatalf("user role = %q, want %q", user.Role, UserRoleUser)
	}
	if user.PreferredLanguage != "en" {
		t.Fatalf("preferred language = %q, want en", user.PreferredLanguage)
	}

	account := AuthAccount{
		AccountID:  user.ID,
		ProviderID: CredentialProviderID,
		UserID:     user.ID,
		Password:   "hash",
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := db.Create(&account).Error; err != nil {
		t.Fatalf("create account: %v", err)
	}

	session := Session{
		Token:     "session-token",
		UserID:    user.ID,
		ExpiresAt: now.Add(time.Hour),
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := db.Create(&session).Error; err != nil {
		t.Fatalf("create session: %v", err)
	}

	verification := Verification{
		Identifier: "email-verification:ada@example.com",
		Value:      "hmac-value",
		ExpiresAt:  now.Add(time.Minute),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := db.Create(&verification).Error; err != nil {
		t.Fatalf("create verification: %v", err)
	}

	assertTableHasRows(t, db, "user", 1)
	assertTableHasRows(t, db, "account", 1)
	assertTableHasRows(t, db, "session", 1)
	assertTableHasRows(t, db, "verification", 1)
}

func openAuthModelTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&User{}, &AuthAccount{}, &Session{}, &Verification{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	return db
}

func assertTableHasRows(t *testing.T, db *gorm.DB, table string, want int64) {
	t.Helper()
	var got int64
	if err := db.Table(table).Count(&got).Error; err != nil {
		t.Fatalf("count %s: %v", table, err)
	}
	if got != want {
		t.Fatalf("%s rows = %d, want %d", table, got, want)
	}
}
