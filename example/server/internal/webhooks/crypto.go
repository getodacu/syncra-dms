package webhooks

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strings"
)

const minAppPrivateKeyLength = 32

var (
	errAppPrivateKeyRequired = errors.New("APP_PRIVATE_KEY is required")
	errAppPrivateKeyTooShort = errors.New("APP_PRIVATE_KEY must be at least 32 characters")
)

func GenerateSecret() (string, error) {
	var data [32]byte
	if _, err := rand.Read(data[:]); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data[:]), nil
}

func EncryptSecret(privateKey string, plaintext string) (string, error) {
	trimmedKey, err := validateAppPrivateKey(privateKey)
	if err != nil {
		return "", err
	}
	key := sha256.Sum256([]byte(trimmedKey))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)
	return "v1:" + base64.RawURLEncoding.EncodeToString(nonce) + ":" + base64.RawURLEncoding.EncodeToString(ciphertext), nil
}

func DecryptSecret(privateKey string, encrypted string) (string, error) {
	trimmedKey, err := validateAppPrivateKey(privateKey)
	if err != nil {
		return "", err
	}

	parts := strings.Split(encrypted, ":")
	if len(parts) != 3 || parts[0] != "v1" {
		return "", errors.New("invalid encrypted secret format")
	}
	nonce, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", errors.New("invalid encrypted secret nonce")
	}
	ciphertext, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return "", errors.New("invalid encrypted secret ciphertext")
	}

	key := sha256.Sum256([]byte(trimmedKey))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(nonce) != gcm.NonceSize() {
		return "", errors.New("invalid encrypted secret nonce")
	}
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

func validateAppPrivateKey(privateKey string) (string, error) {
	trimmedKey := strings.TrimSpace(privateKey)
	if trimmedKey == "" {
		return "", errAppPrivateKeyRequired
	}
	if len(trimmedKey) < minAppPrivateKeyLength {
		return "", errAppPrivateKeyTooShort
	}
	return trimmedKey, nil
}
