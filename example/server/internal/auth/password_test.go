package auth

import "testing"

func TestHashPasswordAndVerify(t *testing.T) {
	hash, err := HashPassword("password1234")
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	if hash == "" || hash == "password1234" {
		t.Fatalf("unsafe hash = %q", hash)
	}
	if ok := VerifyPassword("password1234", hash); !ok {
		t.Fatal("VerifyPassword returned false for correct password")
	}
	if ok := VerifyPassword("wrong-password", hash); ok {
		t.Fatal("VerifyPassword returned true for wrong password")
	}
}

func TestVerifyPasswordRejectsMalformedHash(t *testing.T) {
	if ok := VerifyPassword("password1234", "not-a-valid-hash"); ok {
		t.Fatal("malformed hash verified")
	}
}
