package api

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"ai.ro/syncra/internal/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const publicAPIUserIDContextKey = "syncra_public_api_user_id"

func (h *Handler) publicAPIAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		secret, ok := apiKeyFromAuthorization(c.GetHeader("Authorization"))
		if !ok {
			writeError(c, http.StatusUnauthorized, "invalid API key")
			c.Abort()
			return
		}
		if h.DB == nil {
			writeError(c, http.StatusInternalServerError, "authentication service unavailable")
			c.Abort()
			return
		}

		var key auth.APIKey
		err := h.DB.WithContext(c.Request.Context()).First(&key, "key_hash = ?", auth.HashAPIKey(secret)).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				writeError(c, http.StatusUnauthorized, "invalid API key")
				c.Abort()
				return
			}
			writeError(c, http.StatusInternalServerError, "failed to authenticate API key")
			c.Abort()
			return
		}
		if key.ExpiresAt != nil && !key.ExpiresAt.After(time.Now().UTC()) {
			writeError(c, http.StatusUnauthorized, "invalid API key")
			c.Abort()
			return
		}

		c.Set(publicAPIUserIDContextKey, key.UserID)
		c.Next()
	}
}

func apiKeyFromAuthorization(header string) (string, bool) {
	fields := strings.Fields(header)
	switch {
	case len(fields) == 1:
		return fields[0], true
	case len(fields) == 2 && strings.EqualFold(fields[0], "Bearer"):
		return fields[1], true
	default:
		return "", false
	}
}

func publicAPIUserID(c *gin.Context) (string, bool) {
	value, ok := c.Get(publicAPIUserIDContextKey)
	if !ok {
		return "", false
	}
	userID, ok := value.(string)
	if !ok || strings.TrimSpace(userID) == "" {
		return "", false
	}
	return userID, true
}
