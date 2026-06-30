package webhooks

import (
	"encoding/base64"
	"strings"
	"testing"
)

const (
	testPrivateKey      = "0123456789abcdef0123456789abcdef"
	testWrongPrivateKey = "abcdef0123456789abcdef0123456789"
)

func TestGenerateSecretReturnsURLSafeString(t *testing.T) {
	secret, err := GenerateSecret()
	if err != nil {
		t.Fatalf("GenerateSecret() error = %v", err)
	}
	if secret == "" {
		t.Fatal("GenerateSecret() returned empty secret")
	}
	if strings.ContainsAny(secret, "+/=") {
		t.Fatalf("GenerateSecret() = %q, want raw URL-safe base64", secret)
	}
	if _, err := base64.RawURLEncoding.DecodeString(secret); err != nil {
		t.Fatalf("GenerateSecret() returned invalid raw URL-safe base64: %v", err)
	}
}

func TestEncryptDecryptSecret(t *testing.T) {
	ciphertext, err := EncryptSecret(testPrivateKey, "webhook-secret")
	if err != nil {
		t.Fatalf("EncryptSecret() error = %v", err)
	}
	if strings.Contains(ciphertext, "webhook-secret") {
		t.Fatal("ciphertext contains plaintext secret")
	}
	if !strings.HasPrefix(ciphertext, "v1:") {
		t.Fatalf("ciphertext = %q, want v1 prefix", ciphertext)
	}

	plaintext, err := DecryptSecret(testPrivateKey, ciphertext)
	if err != nil {
		t.Fatalf("DecryptSecret() error = %v", err)
	}
	if plaintext != "webhook-secret" {
		t.Fatalf("plaintext = %q", plaintext)
	}

	if _, err := DecryptSecret(testWrongPrivateKey, ciphertext); err == nil {
		t.Fatal("DecryptSecret() expected wrong key error")
	}
}

func TestEncryptDecryptSecretRequiresPrivateKey(t *testing.T) {
	if _, err := EncryptSecret(" ", "webhook-secret"); err == nil || err.Error() != "APP_PRIVATE_KEY is required" {
		t.Fatalf("EncryptSecret() error = %v, want APP_PRIVATE_KEY is required", err)
	}
	if _, err := DecryptSecret(" ", "v1:nonce:ciphertext"); err == nil || err.Error() != "APP_PRIVATE_KEY is required" {
		t.Fatalf("DecryptSecret() error = %v, want APP_PRIVATE_KEY is required", err)
	}
}

func TestEncryptDecryptSecretRejectsShortPrivateKey(t *testing.T) {
	if _, err := EncryptSecret("short-private-key", "webhook-secret"); err == nil || err.Error() != "APP_PRIVATE_KEY must be at least 32 characters" {
		t.Fatalf("EncryptSecret() error = %v, want APP_PRIVATE_KEY must be at least 32 characters", err)
	}
	if _, err := DecryptSecret("short-private-key", "v1:nonce:ciphertext"); err == nil || err.Error() != "APP_PRIVATE_KEY must be at least 32 characters" {
		t.Fatalf("DecryptSecret() error = %v, want APP_PRIVATE_KEY must be at least 32 characters", err)
	}
}
