package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/scrypt"
)

const (
	scryptN      = 32768
	scryptR      = 8
	scryptP      = 1
	scryptKeyLen = 32
	saltBytes    = 16
)

func HashPassword(password string) (string, error) {
	salt := make([]byte, saltBytes)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	key, err := scrypt.Key([]byte(password), salt, scryptN, scryptR, scryptP, scryptKeyLen)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"scrypt$v=1$N=%d$r=%d$p=%d$%s$%s",
		scryptN,
		scryptR,
		scryptP,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(key),
	), nil
}

func VerifyPassword(password string, encoded string) bool {
	parts := strings.Split(encoded, "$")
	if len(parts) != 7 || parts[0] != "scrypt" || parts[1] != "v=1" {
		return false
	}

	n, err := parseScryptParam(parts[2], "N")
	if err != nil || n != scryptN {
		return false
	}
	r, err := parseScryptParam(parts[3], "r")
	if err != nil || r != scryptR {
		return false
	}
	p, err := parseScryptParam(parts[4], "p")
	if err != nil || p != scryptP {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil || len(salt) != saltBytes {
		return false
	}
	want, err := base64.RawStdEncoding.DecodeString(parts[6])
	if err != nil || len(want) != scryptKeyLen {
		return false
	}

	got, err := scrypt.Key([]byte(password), salt, n, r, p, len(want))
	if err != nil {
		return false
	}
	return subtle.ConstantTimeCompare(got, want) == 1
}

func parseScryptParam(raw string, key string) (int, error) {
	prefix := key + "="
	if !strings.HasPrefix(raw, prefix) {
		return 0, fmt.Errorf("missing %s", key)
	}
	return strconv.Atoi(strings.TrimPrefix(raw, prefix))
}
