package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"ai.ro/syncra/internal/webhooks"
)

type upsertWebhookRequest struct {
	UserID       string           `json:"user_id"`
	URL          string           `json:"url"`
	EventsActive []webhooks.Event `json:"events_active"`
}

const maxWebhookRequestBytes int64 = 64 << 10

func (h *Handler) GetWebhook(c *gin.Context) {
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

	var hook webhooks.Webhook
	result := db.Where("user_id = ?", userID).Limit(1).Find(&hook)
	if result.Error != nil {
		writeError(c, http.StatusInternalServerError, "failed to load webhook")
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, WebhookEnvelopeResponse{Webhook: nil})
		return
	}

	out := webhookResponse(hook, "")
	c.JSON(http.StatusOK, WebhookEnvelopeResponse{Webhook: &out})
}

func (h *Handler) UpsertWebhook(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxWebhookRequestBytes)

	var req upsertWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, webhookBindErrorMessage(err))
		return
	}

	userID, webhookURL, events, ok := validateUpsertWebhookRequest(c, req)
	if !ok {
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

	var existing webhooks.Webhook
	result := db.Where("user_id = ?", userID).Limit(1).Find(&existing)
	if result.Error != nil {
		writeError(c, http.StatusInternalServerError, "failed to load webhook")
		return
	}
	if result.RowsAffected == 0 {
		h.createWebhook(c, db, userID, webhookURL, events)
		return
	}

	if err := updateWebhookColumns(db, &existing, map[string]any{
		"url":           webhookURL,
		"events_active": events,
	}); err != nil {
		writeWebhookMutationError(c, err, "failed to update webhook")
		return
	}
	c.JSON(http.StatusOK, webhookResponse(existing, ""))
}

func (h *Handler) RegenerateWebhookSecret(c *gin.Context) {
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

	var hook webhooks.Webhook
	result := db.Where("user_id = ?", userID).Limit(1).Find(&hook)
	if result.Error != nil {
		writeError(c, http.StatusInternalServerError, "failed to load webhook")
		return
	}
	if result.RowsAffected == 0 {
		writeError(c, http.StatusNotFound, "webhook not found")
		return
	}

	secret, encryptedSecret, ok := h.generateEncryptedWebhookSecret(c)
	if !ok {
		return
	}
	if err := updateWebhookColumns(db, &hook, map[string]any{
		"secret_key": encryptedSecret,
	}); err != nil {
		writeWebhookMutationError(c, err, "failed to update webhook secret")
		return
	}
	c.JSON(http.StatusOK, webhookResponse(hook, secret))
}

func (h *Handler) DeleteWebhook(c *gin.Context) {
	userID, err := parseRequiredUserID(c.Query("user_id"))
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

	var hook webhooks.Webhook
	result := db.Where("user_id = ?", userID).Limit(1).Find(&hook)
	if result.Error != nil {
		writeError(c, http.StatusInternalServerError, "failed to load webhook")
		return
	}
	if result.RowsAffected == 0 {
		writeError(c, http.StatusNotFound, "webhook not found")
		return
	}
	if err := db.Delete(&hook).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to delete webhook")
		return
	}
	c.JSON(http.StatusOK, DeleteWebhookResponse{DeletedID: hook.ID, DeletedCount: 1})
}

func validateUpsertWebhookRequest(c *gin.Context, req upsertWebhookRequest) (string, string, datatypes.JSON, bool) {
	userID, err := parseRequiredUserID(req.UserID)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return "", "", nil, false
	}
	webhookURL, err := webhooks.ValidateURL(req.URL)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return "", "", nil, false
	}
	events, err := webhooks.EncodeEvents(req.EventsActive)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return "", "", nil, false
	}
	return userID, webhookURL, events, true
}

func (h *Handler) createWebhook(c *gin.Context, db *gorm.DB, userID string, webhookURL string, events datatypes.JSON) {
	secret, encryptedSecret, ok := h.generateEncryptedWebhookSecret(c)
	if !ok {
		return
	}

	hook := webhooks.Webhook{
		UserID:       userID,
		URL:          webhookURL,
		SecretKey:    encryptedSecret,
		EventsActive: events,
	}
	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoNothing: true,
	}).Create(&hook)
	if result.Error != nil {
		if isUniqueConstraintError(result.Error) {
			h.updateWebhookAfterCreateConflict(c, db, userID, webhookURL, events)
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to save webhook")
		return
	}
	if result.RowsAffected == 0 {
		h.updateWebhookAfterCreateConflict(c, db, userID, webhookURL, events)
		return
	}
	c.JSON(http.StatusCreated, webhookResponse(hook, secret))
}

func (h *Handler) updateWebhookAfterCreateConflict(c *gin.Context, db *gorm.DB, userID string, webhookURL string, events datatypes.JSON) {
	var existing webhooks.Webhook
	result := db.Where("user_id = ?", userID).Limit(1).Find(&existing)
	if result.Error != nil {
		writeError(c, http.StatusInternalServerError, "failed to load webhook")
		return
	}
	if result.RowsAffected == 0 {
		writeError(c, http.StatusNotFound, "webhook not found")
		return
	}

	if err := updateWebhookColumns(db, &existing, map[string]any{
		"url":           webhookURL,
		"events_active": events,
	}); err != nil {
		writeWebhookMutationError(c, err, "failed to update webhook")
		return
	}
	c.JSON(http.StatusOK, webhookResponse(existing, ""))
}

func updateWebhookColumns(db *gorm.DB, hook *webhooks.Webhook, values map[string]any) error {
	result := db.Model(&webhooks.Webhook{}).
		Where("id = ? AND user_id = ?", hook.ID, hook.UserID).
		Updates(values)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return db.Where("id = ? AND user_id = ?", hook.ID, hook.UserID).First(hook).Error
}

func writeWebhookMutationError(c *gin.Context, err error, internalMessage string) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		writeError(c, http.StatusNotFound, "webhook not found")
		return
	}
	writeError(c, http.StatusInternalServerError, internalMessage)
}

func (h *Handler) generateEncryptedWebhookSecret(c *gin.Context) (string, string, bool) {
	secret, err := webhooks.GenerateSecret()
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to generate webhook secret")
		return "", "", false
	}
	encryptedSecret, err := webhooks.EncryptSecret(h.AppPrivateKey, secret)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to encrypt webhook secret")
		return "", "", false
	}
	return secret, encryptedSecret, true
}

func webhookBindErrorMessage(err error) string {
	if strings.Contains(err.Error(), "request body too large") {
		return "request body too large"
	}
	return "invalid JSON body"
}
