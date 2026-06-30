package auth

import (
	"strings"
	"testing"
)

func TestGenerateNumericCodeReturnsDigits(t *testing.T) {
	code, err := GenerateNumericCode(6)
	if err != nil {
		t.Fatalf("GenerateNumericCode() error = %v", err)
	}
	if len(code) != 6 {
		t.Fatalf("code length = %d, want 6", len(code))
	}
	if strings.Trim(code, "0123456789") != "" {
		t.Fatalf("code = %q, want digits only", code)
	}
}

func TestHashCodeVerifiesIdentifierScopedValue(t *testing.T) {
	hash := HashCode("secret", "email-verification:ada@example.com", "123456")

	if !VerifyCode("secret", "email-verification:ada@example.com", "123456", hash) {
		t.Fatal("VerifyCode rejected the matching code")
	}
	if VerifyCode("secret", "email-verification:ada@example.com", "000000", hash) {
		t.Fatal("VerifyCode accepted a different code")
	}
	if VerifyCode("secret", "password-reset:ada@example.com", "123456", hash) {
		t.Fatal("VerifyCode accepted a different identifier")
	}
}

func TestGenerateSessionTokenReturnsURLSafeSecret(t *testing.T) {
	token, err := GenerateSessionToken()
	if err != nil {
		t.Fatalf("GenerateSessionToken() error = %v", err)
	}
	if len(token) < 40 {
		t.Fatalf("token length = %d, want at least 40", len(token))
	}
	if strings.ContainsAny(token, "+/=") {
		t.Fatalf("token = %q, want raw URL-safe base64", token)
	}
}
