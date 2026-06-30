package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"ai.ro/syncra/internal/auth"
	"ai.ro/syncra/internal/billing"
)

type createBillingOrderRequest struct {
	UserID  string `json:"user_id"`
	Credits int    `json:"credits"`
}

type upsertBillingProfileRequest struct {
	UserID             string  `json:"user_id"`
	EntityType         string  `json:"entity_type"`
	BillingName        string  `json:"billing_name"`
	BillingEmail       string  `json:"billing_email"`
	CountryCode        string  `json:"country_code"`
	AddressLine1       string  `json:"address_line1"`
	AddressLine2       *string `json:"address_line2"`
	City               string  `json:"city"`
	Region             *string `json:"region"`
	PostalCode         string  `json:"postal_code"`
	FiscalCode         *string `json:"fiscal_code"`
	RegistrationNumber *string `json:"registration_number"`
}

type createBillingInvoiceRequest struct {
	UserID       string                            `json:"user_id"`
	InvoiceSerie string                            `json:"invoice_serie"`
	InvoiceDate  *string                           `json:"invoice_date"`
	Lines        []createBillingInvoiceLineRequest `json:"lines"`
}

type createBillingInvoiceLineRequest struct {
	Name          string `json:"name"`
	Quantity      int    `json:"quantity"`
	UnitPrice     string `json:"unit_price"`
	VATPercentage string `json:"vat_percentage"`
}

type attachCheckoutSessionRequest struct {
	CheckoutSessionID string `json:"checkout_session_id"`
}

type markBillingOrderPaidRequest struct {
	CheckoutSessionID *string `json:"checkout_session_id"`
	PaymentIntentID   *string `json:"payment_intent_id"`
	PaidAt            string  `json:"paid_at"`
}

type creditUsageHistoryCursor struct {
	CreatedAt time.Time `json:"created_at"`
	ID        uuid.UUID `json:"id"`
	Sort      string    `json:"sort"`
}

type billingOrderCursor struct {
	CreatedAt time.Time `json:"created_at"`
	ID        uuid.UUID `json:"id"`
	Sort      string    `json:"sort"`
}

type billingInvoiceCursor struct {
	CreatedAt time.Time `json:"created_at"`
	ID        uuid.UUID `json:"id"`
	Sort      string    `json:"sort"`
}

func (h *Handler) GetCreditBalance(c *gin.Context) {
	userID, err := parseRequiredUserID(c.Query("user_id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	balance, err := billing.AvailableCredits(c.Request.Context(), h.DB, userID, time.Now().UTC())
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to load credit balance")
		return
	}
	c.JSON(http.StatusOK, CreditBalanceResponse{
		UserID:           userIDString(userID),
		AvailableCredits: balance.Available,
	})
}

func (h *Handler) GetPublicCreditBalance(c *gin.Context) {
	if rejectPublicQueryUserID(c) {
		return
	}
	userID, ok := publicAPIUserID(c)
	if !ok {
		writeError(c, http.StatusInternalServerError, "authenticated user not found")
		return
	}

	balance, err := billing.AvailableCredits(c.Request.Context(), h.DB, userID, time.Now().UTC())
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to load credit balance")
		return
	}
	c.JSON(http.StatusOK, CreditBalanceResponse{
		UserID:           userIDString(userID),
		AvailableCredits: balance.Available,
	})
}

func (h *Handler) GetBillingProfile(c *gin.Context) {
	if !h.trustedInternalRequest(c) {
		writeError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	userID, err := parseRequiredUserID(c.Query("user_id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	profile, err := billing.GetBillingProfile(c.Request.Context(), h.DB, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, BillingProfileEnvelopeResponse{Profile: nil})
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to load billing profile")
		return
	}
	out := billingProfileResponse(profile)
	c.JSON(http.StatusOK, BillingProfileEnvelopeResponse{Profile: &out})
}

func (h *Handler) UpsertBillingProfile(c *gin.Context) {
	if !h.trustedInternalRequest(c) {
		writeError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req upsertBillingProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	userID, err := parseRequiredUserID(req.UserID)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := validateBillingUserExists(h.DB, userID); err != nil {
		if errors.Is(err, errInvalidUserID) {
			writeError(c, http.StatusBadRequest, err.Error())
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to validate user")
		return
	}

	profile, err := billing.UpsertBillingProfile(c.Request.Context(), h.DB, billing.UpsertBillingProfileInput{
		UserID:             userID,
		EntityType:         billing.BillingEntityType(req.EntityType),
		BillingName:        req.BillingName,
		BillingEmail:       req.BillingEmail,
		CountryCode:        req.CountryCode,
		AddressLine1:       req.AddressLine1,
		AddressLine2:       req.AddressLine2,
		City:               req.City,
		Region:             req.Region,
		PostalCode:         req.PostalCode,
		FiscalCode:         req.FiscalCode,
		RegistrationNumber: req.RegistrationNumber,
	})
	if err != nil {
		if message, ok := billingProfileValidationError(err); ok {
			writeError(c, http.StatusBadRequest, message)
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to upsert billing profile")
		return
	}
	loggerFromGin(c).Info("billing.profile_upserted", "domain", "billing", "user_id", userID, "profile_id", profile.ID.String())
	c.JSON(http.StatusOK, billingProfileResponse(profile))
}

func (h *Handler) CreateBillingInvoice(c *gin.Context) {
	if !h.trustedInternalRequest(c) {
		writeError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req createBillingInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	userID, err := parseRequiredUserID(req.UserID)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := validateBillingUserExists(h.DB, userID); err != nil {
		if errors.Is(err, errInvalidUserID) {
			writeError(c, http.StatusBadRequest, err.Error())
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to validate user")
		return
	}
	invoiceDate, err := parseBillingInvoiceDate(req.InvoiceDate)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	lines := make([]billing.CreateBillingInvoiceLineInput, 0, len(req.Lines))
	for _, line := range req.Lines {
		lines = append(lines, billing.CreateBillingInvoiceLineInput{
			Name:          line.Name,
			Quantity:      line.Quantity,
			UnitPrice:     line.UnitPrice,
			VATPercentage: line.VATPercentage,
		})
	}
	invoice, err := billing.CreateBillingInvoice(c.Request.Context(), h.DB, billing.CreateBillingInvoiceInput{
		UserID:       userID,
		InvoiceSerie: req.InvoiceSerie,
		InvoiceDate:  invoiceDate,
		Lines:        lines,
	})
	if err != nil {
		if message, ok := billingInvoiceValidationError(err); ok {
			writeError(c, http.StatusBadRequest, message)
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to create billing invoice")
		return
	}
	loggerFromGin(c).Info("billing.invoice_created", "domain", "billing", "user_id", userID, "invoice_id", invoice.ID.String())
	c.JSON(http.StatusCreated, billingInvoiceResponse(invoice))
}

func (h *Handler) ListCreditUsageHistory(c *gin.Context) {
	if !h.trustedInternalRequest(c) {
		writeError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	userID, err := parseRequiredUserID(c.Query("user_id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	entryType, err := parseCreditUsageHistoryType(c.Query("type"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	createdFrom, err := parseOCRJobTimeQuery(c, "created_from")
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid created_from")
		return
	}
	createdTo, err := parseOCRJobTimeQuery(c, "created_to")
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid created_to")
		return
	}
	if createdFrom != nil && createdTo != nil && createdFrom.After(*createdTo) {
		writeError(c, http.StatusBadRequest, "created_from must be before or equal to created_to")
		return
	}
	sortDirection, err := parseOCRJobListSort(c.Query("sort"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	size, err := parseOCRJobListSize(c.Query("size"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	cursor, err := parseCreditUsageHistoryCursor(c.Query("cursor"), sortDirection)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	page, err := billing.ListCreditLedgerTransactions(c.Request.Context(), h.DB, billing.ListCreditLedgerTransactionsInput{
		UserID:      userID,
		EntryType:   entryType,
		CreatedFrom: createdFrom,
		CreatedTo:   createdTo,
		Cursor:      cursor,
		Size:        size,
		Sort:        sortDirection,
	})
	if err != nil {
		if message, ok := creditUsageHistoryValidationError(err); ok {
			writeError(c, http.StatusBadRequest, message)
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to list credit usage history")
		return
	}

	var nextCursor *string
	if page.NextCursor != nil {
		encoded, err := encodeCreditUsageHistoryCursor(*page.NextCursor)
		if err != nil {
			writeError(c, http.StatusInternalServerError, "failed to encode next cursor")
			return
		}
		nextCursor = &encoded
	}

	out := make([]CreditUsageHistoryEntryResponse, 0, len(page.Entries))
	for _, entry := range page.Entries {
		out = append(out, creditUsageHistoryEntryResponse(entry))
	}
	c.JSON(http.StatusOK, CreditUsageHistoryListResponse{
		CreditUsageHistory: out,
		NextCursor:         nextCursor,
	})
}

func (h *Handler) ListBillingOrders(c *gin.Context) {
	if !h.trustedInternalRequest(c) {
		writeError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	userID, err := parseRequiredUserID(c.Query("user_id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	status, err := parseBillingOrderStatus(c.Query("status"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	createdFrom, err := parseOCRJobTimeQuery(c, "created_from")
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid created_from")
		return
	}
	createdTo, err := parseOCRJobTimeQuery(c, "created_to")
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid created_to")
		return
	}
	if createdFrom != nil && createdTo != nil && createdFrom.After(*createdTo) {
		writeError(c, http.StatusBadRequest, "created_from must be before or equal to created_to")
		return
	}
	sortDirection, err := parseOCRJobListSort(c.Query("sort"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	size, err := parseOCRJobListSize(c.Query("size"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	cursor, err := parseBillingOrderCursor(c.Query("cursor"), sortDirection)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	page, err := billing.ListBillingOrders(c.Request.Context(), h.DB, billing.ListBillingOrdersInput{
		UserID:      userID,
		Status:      status,
		CreatedFrom: createdFrom,
		CreatedTo:   createdTo,
		Cursor:      cursor,
		Size:        size,
		Sort:        sortDirection,
	})
	if err != nil {
		if message, ok := billingOrderListValidationError(err); ok {
			writeError(c, http.StatusBadRequest, message)
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to list billing orders")
		return
	}

	var nextCursor *string
	if page.NextCursor != nil {
		encoded, err := encodeBillingOrderCursor(*page.NextCursor)
		if err != nil {
			writeError(c, http.StatusInternalServerError, "failed to encode next cursor")
			return
		}
		nextCursor = &encoded
	}

	out := make([]BillingOrderResponse, 0, len(page.Orders))
	for _, order := range page.Orders {
		out = append(out, billingOrderResponse(order))
	}
	c.JSON(http.StatusOK, BillingOrderListResponse{
		Orders:     out,
		NextCursor: nextCursor,
	})
}

func (h *Handler) CreateBillingOrder(c *gin.Context) {
	if !h.trustedInternalRequest(c) {
		writeError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req createBillingOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	userID, err := parseRequiredUserID(req.UserID)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := validateBillingUserExists(h.DB, userID); err != nil {
		if errors.Is(err, errInvalidUserID) {
			writeError(c, http.StatusBadRequest, err.Error())
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to validate user")
		return
	}

	order, err := billing.CreateCreditOrder(c.Request.Context(), h.DB, billing.CreateCreditOrderInput{
		UserID:  userID,
		Credits: req.Credits,
	})
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	loggerFromGin(c).Info("billing.order_created",
		"domain", "billing",
		"user_id", userID,
		"order_id", order.ID.String(),
		"credits", order.Credits,
		"status", string(order.Status),
	)
	c.JSON(http.StatusCreated, billingOrderResponse(order))
}

func (h *Handler) AttachBillingOrderCheckoutSession(c *gin.Context) {
	if !h.trustedInternalRequest(c) {
		writeError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	orderID, ok := parseBillingOrderIDParam(c)
	if !ok {
		return
	}
	var req attachCheckoutSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	checkoutSessionID := strings.TrimSpace(req.CheckoutSessionID)
	if checkoutSessionID == "" {
		writeError(c, http.StatusBadRequest, "checkout_session_id is required")
		return
	}
	if err := billing.AttachCreditOrderCheckoutSession(c.Request.Context(), h.DB, orderID, checkoutSessionID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusNotFound, "billing order not found")
			return
		}
		if errors.Is(err, billing.ErrProviderMetadataConflict) {
			writeError(c, http.StatusConflict, "payment metadata belongs to another billing order")
			return
		}
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	loggerFromGin(c).Info("billing.order_checkout_session_attached", "domain", "billing", "order_id", orderID.String())
	c.Status(http.StatusNoContent)
}

func (h *Handler) MarkBillingOrderPaid(c *gin.Context) {
	if !h.trustedInternalRequest(c) {
		writeError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	orderID, ok := parseBillingOrderIDParam(c)
	if !ok {
		return
	}
	var req markBillingOrderPaidRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	paidAt, err := time.Parse(time.RFC3339, strings.TrimSpace(req.PaidAt))
	if err != nil || paidAt.IsZero() {
		writeError(c, http.StatusBadRequest, "invalid paid_at")
		return
	}

	checkoutSessionID := nonEmptyStringPtr(req.CheckoutSessionID)
	paymentIntentID := nonEmptyStringPtr(req.PaymentIntentID)
	if _, err := billing.MarkCreditOrderPaidAndGrantCredits(c.Request.Context(), h.DB, billing.MarkCreditOrderPaidInput{
		OrderID:                   orderID,
		ProviderCheckoutSessionID: checkoutSessionID,
		ProviderPaymentIntentID:   paymentIntentID,
		PaidAt:                    paidAt.UTC(),
	}); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusNotFound, "billing order not found")
			return
		}
		if errors.Is(err, billing.ErrProviderMetadataConflict) {
			writeError(c, http.StatusConflict, "payment metadata belongs to another billing order")
			return
		}
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	loggerFromGin(c).Info("billing.order_marked_paid", "domain", "billing", "order_id", orderID.String())
	h.ensurePaidBillingOrderInvoiceArtifacts(c, orderID)
	order, ok := h.loadBillingOrderForResponse(c, orderID)
	if !ok {
		return
	}
	if !billingOrderMatchesPaidRequest(order, checkoutSessionID, paymentIntentID) {
		writeError(c, http.StatusConflict, "payment metadata belongs to another billing order")
		return
	}
	c.JSON(http.StatusOK, billingOrderResponse(order))
}

func (h *Handler) ensurePaidBillingOrderInvoiceArtifacts(c *gin.Context, orderID uuid.UUID) {
	logger := loggerFromGin(c).With("domain", "billing", "order_id", orderID.String())
	invoiceDate, err := currentBillingInvoiceDate(h.now())
	if err != nil {
		logger.Error("billing.paid_order_invoice_date_failed", "error", err)
		return
	}

	invoice, err := billing.CreateBillingInvoiceForPaidOrder(c.Request.Context(), h.DB, orderID, invoiceDate)
	if errors.Is(err, billing.ErrBillingInvoiceExists) {
		invoice, err = h.loadBillingInvoiceForOrder(c, orderID)
	}
	if err != nil {
		logger.Error("billing.paid_order_invoice_create_failed", "error", err)
		return
	}
	logger = logger.With("invoice_id", invoice.ID.String())

	needsPDF, err := h.billingInvoicePDFNeedsGeneration(invoice)
	if err != nil {
		logger.Error("billing.paid_order_invoice_pdf_inspect_failed", "error", err)
		return
	}
	if !needsPDF {
		logger.Debug("billing.paid_order_invoice_pdf_current")
		return
	}
	if _, err := h.generateAndStoreBillingInvoicePDF(c.Request.Context(), invoice); err != nil {
		logger.Error("billing.paid_order_invoice_pdf_generate_failed", "error", err)
		return
	}
	logger.Info("billing.paid_order_invoice_pdf_generated")
}

func (h *Handler) loadBillingInvoiceForOrder(c *gin.Context, orderID uuid.UUID) (billing.BillingInvoice, error) {
	var invoice billing.BillingInvoice
	if err := h.DB.WithContext(c.Request.Context()).First(&invoice, "order_id = ?", orderID).Error; err != nil {
		return billing.BillingInvoice{}, err
	}
	return invoice, nil
}

func (h *Handler) MarkBillingOrderFailed(c *gin.Context) {
	if !h.trustedInternalRequest(c) {
		writeError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	orderID, ok := parseBillingOrderIDParam(c)
	if !ok {
		return
	}
	if err := billing.MarkCreditOrderFailed(c.Request.Context(), h.DB, orderID, time.Now().UTC()); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusNotFound, "billing order not found")
			return
		}
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	loggerFromGin(c).Info("billing.order_marked_failed", "domain", "billing", "order_id", orderID.String())
	h.respondWithBillingOrder(c, orderID)
}

func (h *Handler) respondWithBillingOrder(c *gin.Context, orderID uuid.UUID) {
	order, ok := h.loadBillingOrderForResponse(c, orderID)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, billingOrderResponse(order))
}

func (h *Handler) loadBillingOrderForResponse(c *gin.Context, orderID uuid.UUID) (billing.BillingOrder, bool) {
	var order billing.BillingOrder
	if err := h.DB.WithContext(c.Request.Context()).Preload("Invoice").First(&order, "id = ?", orderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusNotFound, "billing order not found")
			return billing.BillingOrder{}, false
		}
		writeError(c, http.StatusInternalServerError, "failed to load billing order")
		return billing.BillingOrder{}, false
	}
	return order, true
}

func billingOrderMatchesPaidRequest(order billing.BillingOrder, checkoutSessionID *string, paymentIntentID *string) bool {
	if order.Status != billing.OrderStatusPaid {
		return false
	}
	if checkoutSessionID != nil && (order.ProviderCheckoutSessionID == nil || *order.ProviderCheckoutSessionID != *checkoutSessionID) {
		return false
	}
	if paymentIntentID != nil && (order.ProviderPaymentIntentID == nil || *order.ProviderPaymentIntentID != *paymentIntentID) {
		return false
	}
	return true
}

func parseBillingOrderIDParam(c *gin.Context) (uuid.UUID, bool) {
	orderID, err := uuid.Parse(strings.TrimSpace(c.Param("id")))
	if err != nil || orderID == uuid.Nil {
		writeError(c, http.StatusBadRequest, "invalid billing order id")
		return uuid.Nil, false
	}
	return orderID, true
}

func validateBillingUserExists(db *gorm.DB, userID string) error {
	var count int64
	if err := db.Model(&auth.User{}).Where("id = ?", userID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return errInvalidUserID
	}
	return nil
}

func creditUsageHistoryValidationError(err error) (string, bool) {
	message := strings.TrimPrefix(err.Error(), "billing: ")
	switch message {
	case "user id is required",
		"size must be between 1 and 100",
		"sort must be asc or desc",
		"entry type must be purchase or debit",
		"created_from must be before or equal to created_to",
		"invalid cursor":
		return message, true
	default:
		return "", false
	}
}

func billingOrderListValidationError(err error) (string, bool) {
	message := strings.TrimPrefix(err.Error(), "billing: ")
	switch message {
	case "user id is required",
		"size must be between 1 and 100",
		"sort must be asc or desc",
		"created_from must be before or equal to created_to",
		"invalid cursor":
		return message, true
	default:
		if strings.HasPrefix(message, "invalid order status") {
			return message, true
		}
		return "", false
	}
}

func billingInvoiceListValidationError(err error) (string, bool) {
	message := strings.TrimPrefix(err.Error(), "billing: ")
	switch message {
	case "user id is required",
		"size must be between 1 and 100",
		"sort must be asc or desc",
		"created_from must be before or equal to created_to",
		"invalid cursor":
		return message, true
	default:
		return "", false
	}
}

func billingProfileValidationError(err error) (string, bool) {
	message := strings.TrimPrefix(err.Error(), "billing: ")
	switch message {
	case "user id is required",
		"billing name is required",
		"billing email is invalid",
		"country is required",
		"country code is invalid",
		"address line 1 is required",
		"city is required",
		"postal code is required",
		"fiscal code is required for Romanian companies":
		return message, true
	default:
		if strings.HasSuffix(message, "characters or fewer") {
			return message, true
		}
		if strings.HasPrefix(message, "invalid billing entity type") {
			return message, true
		}
		return "", false
	}
}

func billingInvoiceValidationError(err error) (string, bool) {
	if errors.Is(err, billing.ErrBillingProfileRequired) {
		return "billing profile is required", true
	}
	message := strings.TrimPrefix(err.Error(), "billing: ")
	switch message {
	case "user id is required",
		"invoice serie is required",
		"invoice serie contains invalid characters",
		"invoice lines are required":
		return message, true
	default:
		if strings.HasPrefix(message, "invoice line ") {
			return message, true
		}
		if strings.HasPrefix(message, "invoice serie must be ") {
			return message, true
		}
		return "", false
	}
}

func parseBillingInvoiceDate(raw *string) (time.Time, error) {
	if raw == nil || strings.TrimSpace(*raw) == "" {
		return time.Time{}, nil
	}
	parsed, err := time.Parse("2006-01-02", strings.TrimSpace(*raw))
	if err != nil {
		return time.Time{}, errors.New("invalid invoice_date")
	}
	return parsed, nil
}

func parseCreditUsageHistoryType(raw string) (*billing.CreditLedgerEntryType, error) {
	raw = strings.ToLower(strings.TrimSpace(raw))
	if raw == "" {
		return nil, nil
	}
	entryType := billing.CreditLedgerEntryType(raw)
	switch entryType {
	case billing.CreditLedgerEntryPurchase, billing.CreditLedgerEntryDebit:
		return &entryType, nil
	default:
		return nil, errors.New("entry type must be purchase or debit")
	}
}

func parseBillingOrderStatus(raw string) (*billing.OrderStatus, error) {
	raw = strings.ToLower(strings.TrimSpace(raw))
	if raw == "" {
		return nil, nil
	}
	status := billing.OrderStatus(raw)
	switch status {
	case billing.OrderStatusPending,
		billing.OrderStatusPaid,
		billing.OrderStatusFailed,
		billing.OrderStatusRefunded,
		billing.OrderStatusCanceled:
		return &status, nil
	default:
		return nil, errors.New("invalid order status")
	}
}

func parseCreditUsageHistoryCursor(raw string, sortDirection string) (*billing.CreditLedgerTransactionCursor, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	decoded, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil {
		return nil, errors.New("invalid cursor")
	}
	var cursor creditUsageHistoryCursor
	if err := json.Unmarshal(decoded, &cursor); err != nil {
		return nil, errors.New("invalid cursor")
	}
	if cursor.ID == uuid.Nil || cursor.CreatedAt.IsZero() || (cursor.Sort != "asc" && cursor.Sort != "desc") {
		return nil, errors.New("invalid cursor")
	}
	if cursor.Sort != sortDirection {
		return nil, errors.New("cursor sort does not match sort")
	}
	return &billing.CreditLedgerTransactionCursor{
		CreatedAt: cursor.CreatedAt,
		ID:        cursor.ID,
		Sort:      cursor.Sort,
	}, nil
}

func parseBillingOrderCursor(raw string, sortDirection string) (*billing.BillingOrderCursor, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	decoded, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil {
		return nil, errors.New("invalid cursor")
	}
	var cursor billingOrderCursor
	if err := json.Unmarshal(decoded, &cursor); err != nil {
		return nil, errors.New("invalid cursor")
	}
	if cursor.ID == uuid.Nil || cursor.CreatedAt.IsZero() || (cursor.Sort != "asc" && cursor.Sort != "desc") {
		return nil, errors.New("invalid cursor")
	}
	if cursor.Sort != sortDirection {
		return nil, errors.New("cursor sort does not match sort")
	}
	return &billing.BillingOrderCursor{
		CreatedAt: cursor.CreatedAt,
		ID:        cursor.ID,
		Sort:      cursor.Sort,
	}, nil
}

func parseBillingInvoiceCursor(raw string, sortDirection string) (*billing.BillingInvoiceCursor, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	decoded, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil {
		return nil, errors.New("invalid cursor")
	}
	var cursor billingInvoiceCursor
	if err := json.Unmarshal(decoded, &cursor); err != nil {
		return nil, errors.New("invalid cursor")
	}
	if cursor.ID == uuid.Nil || cursor.CreatedAt.IsZero() || (cursor.Sort != "asc" && cursor.Sort != "desc") {
		return nil, errors.New("invalid cursor")
	}
	if cursor.Sort != sortDirection {
		return nil, errors.New("cursor sort does not match sort")
	}
	return &billing.BillingInvoiceCursor{
		CreatedAt: cursor.CreatedAt,
		ID:        cursor.ID,
		Sort:      cursor.Sort,
	}, nil
}

func encodeCreditUsageHistoryCursor(cursor billing.CreditLedgerTransactionCursor) (string, error) {
	raw, err := json.Marshal(creditUsageHistoryCursor{
		CreatedAt: cursor.CreatedAt.UTC(),
		ID:        cursor.ID,
		Sort:      cursor.Sort,
	})
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}

func encodeBillingOrderCursor(cursor billing.BillingOrderCursor) (string, error) {
	raw, err := json.Marshal(billingOrderCursor{
		CreatedAt: cursor.CreatedAt.UTC(),
		ID:        cursor.ID,
		Sort:      cursor.Sort,
	})
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}

func encodeBillingInvoiceCursor(cursor billing.BillingInvoiceCursor) (string, error) {
	raw, err := json.Marshal(billingInvoiceCursor{
		CreatedAt: cursor.CreatedAt.UTC(),
		ID:        cursor.ID,
		Sort:      cursor.Sort,
	})
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}

func nonEmptyStringPtr(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
