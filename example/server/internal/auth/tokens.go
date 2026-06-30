package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
)

func GenerateSessionToken() (string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}

func GenerateAPIKey() (string, error) {
	raw := make([]byte, 48)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}

func HashAPIKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}

func GenerateNumericCode(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be positive")
	}

	out := make([]byte, length)
	for i := range out {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		out[i] = byte('0' + n.Int64())
	}
	return string(out), nil
}

func HashCode(secret string, identifier string, code string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(identifier))
	_, _ = mac.Write([]byte{0})
	_, _ = mac.Write([]byte(code))
	return base64.RawStdEncoding.EncodeToString(mac.Sum(nil))
}

func VerifyCode(secret string, identifier string, code string, expectedHash string) bool {
	got := HashCode(secret, identifier, code)
	return hmac.Equal([]byte(got), []byte(expectedHash))
}
