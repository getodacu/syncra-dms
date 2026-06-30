package auth

import (
	"encoding/base64"
	"encoding/hex"
	"strings"
	"testing"
)

func TestGenerateSessionToken(t *testing.T) {
	token, err := GenerateSessionToken()
	if err != nil {
		t.Fatalf("GenerateSessionToken: %v", err)
	}
	if token == "" || strings.ContainsAny(token, "+/=") {
		t.Fatalf("unsafe session token = %q", token)
	}
	raw, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		t.Fatalf("session token is not raw URL base64: %v", err)
	}
	if len(raw) != 32 {
		t.Fatalf("session token raw length = %d", len(raw))
	}
}

func TestGenerateAPIKey(t *testing.T) {
	key, err := GenerateAPIKey()
	if err != nil {
		t.Fatalf("GenerateAPIKey: %v", err)
	}
	if len(key) != 64 {
		t.Fatalf("api key length = %d, want 64", len(key))
	}
	if strings.ContainsAny(key, "+/=") {
		t.Fatalf("unsafe api key = %q", key)
	}
	raw, err := base64.RawURLEncoding.DecodeString(key)
	if err != nil {
		t.Fatalf("api key is not raw URL base64: %v", err)
	}
	if len(raw) != 48 {
		t.Fatalf("api key raw length = %d, want 48", len(raw))
	}
}

func TestHashAPIKey(t *testing.T) {
	hash := HashAPIKey("test-api-key")
	if len(hash) != 64 {
		t.Fatalf("hash length = %d, want 64", len(hash))
	}
	if _, err := hex.DecodeString(hash); err != nil {
		t.Fatalf("hash is not hex: %v", err)
	}
	if strings.Contains(hash, "test-api-key") {
		t.Fatalf("hash contains plaintext key: %q", hash)
	}
	if hash != HashAPIKey("test-api-key") {
		t.Fatal("HashAPIKey is not deterministic")
	}
	if hash == HashAPIKey("other-api-key") {
		t.Fatal("different API keys produced the same hash")
	}
}

func TestGenerateNumericCode(t *testing.T) {
	code, err := GenerateNumericCode(6)
	if err != nil {
		t.Fatalf("GenerateNumericCode: %v", err)
	}
	if len(code) != 6 {
		t.Fatalf("code length = %d", len(code))
	}
	for _, r := range code {
		if r < '0' || r > '9' {
			t.Fatalf("code contains non-digit: %q", code)
		}
	}
}

func TestGenerateNumericCodeRejectsNonPositiveLength(t *testing.T) {
	if _, err := GenerateNumericCode(0); err == nil {
		t.Fatal("GenerateNumericCode accepted zero length")
	}
	if _, err := GenerateNumericCode(-1); err == nil {
		t.Fatal("GenerateNumericCode accepted negative length")
	}
}

func TestHashAndVerifyCode(t *testing.T) {
	secret := "test-secret"
	hash := HashCode(secret, "email-verification:ada@example.com", "123456")
	if hash == "" || strings.Contains(hash, "123456") {
		t.Fatalf("unsafe code hash = %q", hash)
	}
	if !VerifyCode(secret, "email-verification:ada@example.com", "123456", hash) {
		t.Fatal("VerifyCode returned false for correct code")
	}
	if VerifyCode(secret, "email-verification:ada@example.com", "000000", hash) {
		t.Fatal("VerifyCode returned true for wrong code")
	}
}
