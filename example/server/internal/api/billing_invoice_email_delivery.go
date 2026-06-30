package api

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"ai.ro/syncra/internal/billing"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const billingInvoiceEmailDeliveryClaimLease = 15 * time.Minute

const (
	billingInvoiceEmailDeliveryStatusClaimed     = "claimed"
	billingInvoiceEmailDeliveryStatusClaimActive = "claim_active"
	billingInvoiceEmailDeliveryStatusAlreadySent = "already_sent"
	billingInvoiceEmailDeliveryStatusNotReady    = "not_ready"
	billingInvoiceEmailDeliveryStatusSent        = "sent"
)

type billingInvoiceEmailDeliveryRequest struct {
	UserID string `json:"user_id"`
}

func (h *Handler) ClaimBillingInvoiceEmailDelivery(c *gin.Context) {
	if !h.trustedInternalRequest(c) {
		writeError(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	invoiceID, ok := parseBillingInvoiceIDParam(c)
	if !ok {
		return
	}
	userID, ok := bindBillingInvoiceEmailDeliveryUser(c)
	if !ok {
		return
	}
	if h.DB == nil {
		writeError(c, http.StatusInternalServerError, "database is not configured")
		return
	}

	var response BillingInvoiceEmailDeliveryClaimResponse
	err := h.DB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		invoice, err := lockedUserBillingInvoice(tx, invoiceID, userID)
		if err != nil {
			return err
		}

		now := h.now().UTC()
		switch {
		case invoice.EmailSentAt != nil:
			response.Status = billingInvoiceEmailDeliveryStatusAlreadySent
			return nil
		case !h.billingInvoicePDFFileExists(invoice):
			response.Status = billingInvoiceEmailDeliveryStatusNotReady
			return nil
		case invoice.EmailDeliveryClaimedAt != nil &&
			invoice.EmailDeliveryClaimedAt.After(now.Add(-billingInvoiceEmailDeliveryClaimLease)):
			response.Status = billingInvoiceEmailDeliveryStatusClaimActive
			return nil
		}

		if err := tx.Model(&billing.BillingInvoice{}).
			Where("id = ?", invoice.ID).
			Update("email_delivery_claimed_at", now).Error; err != nil {
			return err
		}
		invoice.EmailDeliveryClaimedAt = &now
		out := billingInvoiceResponse(invoice)
		response.Status = billingInvoiceEmailDeliveryStatusClaimed
		response.Invoice = &out
		return nil
	})
	if err != nil {
		writeBillingInvoiceEmailDeliveryError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) MarkBillingInvoiceEmailSent(c *gin.Context) {
	if !h.trustedInternalRequest(c) {
		writeError(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	invoiceID, ok := parseBillingInvoiceIDParam(c)
	if !ok {
		return
	}
	userID, ok := bindBillingInvoiceEmailDeliveryUser(c)
	if !ok {
		return
	}
	if h.DB == nil {
		writeError(c, http.StatusInternalServerError, "database is not configured")
		return
	}

	var invoice billing.BillingInvoice
	err := h.DB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		locked, err := lockedUserBillingInvoice(tx, invoiceID, userID)
		if err != nil {
			return err
		}
		if locked.EmailSentAt == nil {
			now := h.now().UTC()
			if err := tx.Model(&billing.BillingInvoice{}).
				Where("id = ?", locked.ID).
				Update("email_sent_at", now).Error; err != nil {
				return err
			}
			locked.EmailSentAt = &now
		}
		invoice = locked
		return nil
	})
	if err != nil {
		writeBillingInvoiceEmailDeliveryError(c, err)
		return
	}

	c.JSON(http.StatusOK, BillingInvoiceEmailDeliverySentResponse{
		Status:  billingInvoiceEmailDeliveryStatusSent,
		Invoice: billingInvoiceResponse(invoice),
	})
}

func parseBillingInvoiceIDParam(c *gin.Context) (uuid.UUID, bool) {
	invoiceID, err := uuid.Parse(strings.TrimSpace(c.Param("id")))
	if err != nil || invoiceID == uuid.Nil {
		writeError(c, http.StatusBadRequest, "invalid billing invoice id")
		return uuid.Nil, false
	}
	return invoiceID, true
}

func bindBillingInvoiceEmailDeliveryUser(c *gin.Context) (string, bool) {
	var req billingInvoiceEmailDeliveryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return "", false
	}
	userID, err := parseRequiredUserID(req.UserID)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return "", false
	}
	return userID, true
}

func lockedUserBillingInvoice(tx *gorm.DB, invoiceID uuid.UUID, userID string) (billing.BillingInvoice, error) {
	var invoice billing.BillingInvoice
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ? AND user_id = ?", invoiceID, userID).
		First(&invoice).Error
	if err != nil {
		return billing.BillingInvoice{}, err
	}
	return invoice, nil
}

func (h *Handler) billingInvoicePDFFileExists(invoice billing.BillingInvoice) bool {
	if invoice.PDFPath == nil || strings.TrimSpace(*invoice.PDFPath) == "" {
		return false
	}
	pdfDir, err := h.billingInvoicePDFDir()
	if err != nil {
		return false
	}
	info, err := os.Stat(filepath.Join(pdfDir, invoice.ID.String()+".pdf"))
	return err == nil && !info.IsDir()
}

func writeBillingInvoiceEmailDeliveryError(c *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		writeError(c, http.StatusNotFound, "billing invoice not found")
		return
	}
	writeError(c, http.StatusInternalServerError, "failed to update billing invoice email delivery")
}
