package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"ai.ro/syncra/internal/auth"
)

type createAPIKeyRequest struct {
	UserID    string                       `json:"user_id"`
	Name      string                       `json:"name"`
	ExpiresAt optionalAPIKeyExpirationTime `json:"expires_at"`
}

type optionalAPIKeyExpirationTime struct {
	Value *time.Time
}

const maxAPIKeyNameCharacters = 255
const maxAPIKeyRequestBytes int64 = 64 << 10
const maxAPIKeyGenerateAttempts = 10

var errInvalidAPIKeyExpiresAt = errors.New("expires_at must be RFC3339")

func (t *optionalAPIKeyExpirationTime) UnmarshalJSON(raw []byte) error {
	var value *string
	if err := json.Unmarshal(raw, &value); err != nil {
		return errInvalidAPIKeyExpiresAt
	}
	if value == nil {
		t.Value = nil
		return nil
	}
	parsed, err := time.Parse(time.RFC3339Nano, strings.TrimSpace(*value))
	if err != nil {
		return errInvalidAPIKeyExpiresAt
	}
	utc := parsed.UTC()
	t.Value = &utc
	return nil
}

func (h *Handler) CreateAPIKey(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxAPIKeyRequestBytes)

	var req createAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, apiKeyBindErrorMessage(err))
		return
	}

	userID, err := parseRequiredUserID(req.UserID)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	name, err := validateAPIKeyName(req.Name)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	db := h.DB.WithContext(c.Request.Context())
	if err := validateAPIKeyUserExists(db, userID); err != nil {
		if errors.Is(err, errInvalidUserID) {
			writeError(c, http.StatusBadRequest, err.Error())
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to validate user")
		return
	}

	apiKey, secret, ok := h.createUniqueAPIKey(c, db, userID, name, req.ExpiresAt.Value)
	if !ok {
		return
	}
	c.JSON(http.StatusCreated, apiKeyResponse(apiKey, secret))
}

func (h *Handler) ListAPIKeys(c *gin.Context) {
	userID, err := parseRequiredUserID(c.Param("user_id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	db := h.DB.WithContext(c.Request.Context())
	if err := validateAPIKeyUserExists(db, userID); err != nil {
		if errors.Is(err, errInvalidUserID) {
			writeError(c, http.StatusBadRequest, err.Error())
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to validate user")
		return
	}

	var keys []auth.APIKey
	if err := db.Where("user_id = ?", userID).Order("created_at desc, id desc").Find(&keys).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to list api keys")
		return
	}

	out := make([]APIKeyResponse, 0, len(keys))
	for _, key := range keys {
		out = append(out, apiKeyResponse(key, ""))
	}
	c.JSON(http.StatusOK, APIKeyListResponse{APIKeys: out})
}

func (h *Handler) DeleteAPIKey(c *gin.Context) {
	userID, err := parseRequiredUserID(c.Query("user_id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	apiKeyID, err := parseAPIKeyID(c.Query("api_key_id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	db := h.DB.WithContext(c.Request.Context())
	if err := validateAPIKeyUserExists(db, userID); err != nil {
		if errors.Is(err, errInvalidUserID) {
			writeError(c, http.StatusBadRequest, err.Error())
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to validate user")
		return
	}

	result := db.Where("id = ? AND user_id = ?", apiKeyID, userID).Delete(&auth.APIKey{})
	if result.Error != nil {
		writeError(c, http.StatusInternalServerError, "failed to delete api key")
		return
	}
	if result.RowsAffected == 0 {
		writeError(c, http.StatusNotFound, "api key not found")
		return
	}
	c.JSON(http.StatusOK, DeleteAPIKeyResponse{DeletedID: apiKeyID, DeletedCount: 1})
}

func (h *Handler) createUniqueAPIKey(c *gin.Context, db *gorm.DB, userID string, name string, expiresAt *time.Time) (auth.APIKey, string, bool) {
	generate := h.APIKeyGenerator
	if generate == nil {
		generate = auth.GenerateAPIKey
	}

	for attempt := 0; attempt < maxAPIKeyGenerateAttempts; attempt++ {
		secret, err := generate()
		if err != nil {
			writeError(c, http.StatusInternalServerError, "failed to generate api key")
			return auth.APIKey{}, "", false
		}
		if len(secret) != 64 {
			writeError(c, http.StatusInternalServerError, "failed to generate api key")
			return auth.APIKey{}, "", false
		}
		hash := auth.HashAPIKey(secret)
		var count int64
		if err := db.Model(&auth.APIKey{}).Where("key_hash = ?", hash).Count(&count).Error; err != nil {
			writeError(c, http.StatusInternalServerError, "failed to validate api key")
			return auth.APIKey{}, "", false
		}
		if count > 0 {
			continue
		}

		apiKey := auth.APIKey{
			UserID:    userID,
			Name:      name,
			KeyHash:   hash,
			KeyPrefix: secret[:8],
			ExpiresAt: expiresAt,
		}
		if err := db.Create(&apiKey).Error; err != nil {
			if isUniqueConstraintError(err) {
				continue
			}
			writeError(c, http.StatusInternalServerError, "failed to save api key")
			return auth.APIKey{}, "", false
		}
		return apiKey, secret, true
	}

	writeError(c, http.StatusInternalServerError, "failed to generate unique api key")
	return auth.APIKey{}, "", false
}

func validateAPIKeyName(raw string) (string, error) {
	name := strings.TrimSpace(raw)
	if name == "" {
		return "", errors.New("name is required")
	}
	if utf8.RuneCountInString(name) > maxAPIKeyNameCharacters {
		return "", errors.New("name must be at most 255 characters")
	}
	return name, nil
}

func parseAPIKeyID(raw string) (uuid.UUID, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return uuid.Nil, errors.New("api_key_id is required")
	}
	id, err := uuid.Parse(raw)
	if err != nil || id == uuid.Nil {
		return uuid.Nil, errors.New("invalid api_key_id")
	}
	return id, nil
}

func validateAPIKeyUserExists(db *gorm.DB, userID string) error {
	var count int64
	if err := db.Model(&auth.User{}).Where("id = ?", userID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return errInvalidUserID
	}
	return nil
}

func apiKeyBindErrorMessage(err error) string {
	if strings.Contains(err.Error(), "request body too large") {
		return "request body too large"
	}
	if strings.Contains(err.Error(), errInvalidAPIKeyExpiresAt.Error()) {
		return errInvalidAPIKeyExpiresAt.Error()
	}
	return "invalid JSON body"
}

func apiKeyResponse(key auth.APIKey, secret string) APIKeyResponse {
	return APIKeyResponse{
		ID:        key.ID,
		UserID:    userIDString(key.UserID),
		Name:      key.Name,
		KeyPrefix: key.KeyPrefix,
		APIKey:    secret,
		ExpiresAt: key.ExpiresAt,
		CreatedAt: key.CreatedAt,
		UpdatedAt: key.UpdatedAt,
	}
}
