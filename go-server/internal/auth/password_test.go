package auth

import "testing"

func TestHashPasswordVerifiesOnlyMatchingPassword(t *testing.T) {
	hash, err := HashPassword("correct horse battery staple")
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	if hash == "correct horse battery staple" {
		t.Fatal("HashPassword stored plaintext")
	}
	if !VerifyPassword("correct horse battery staple", hash) {
		t.Fatal("VerifyPassword rejected the matching password")
	}
	if VerifyPassword("wrong password", hash) {
		t.Fatal("VerifyPassword accepted the wrong password")
	}
}
