package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"ai.ro/syncra/internal/billing"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const billingInvoiceBusinessLocation = "Europe/Bucharest"

type adminBillingOrderUserResponse struct {
	ID    userIDString `json:"id"`
	Name  string       `json:"name"`
	Email string       `json:"email"`
}

type adminBillingOrderResponse struct {
	BillingOrderResponse
	User    adminBillingOrderUserResponse     `json:"user"`
	Invoice *adminBillingOrderInvoiceResponse `json:"invoice"`
}

type adminBillingOrderInvoiceResponse struct {
	ID            uuid.UUID `json:"id"`
	InvoiceSerie  string    `json:"invoice_serie"`
	InvoiceNumber int64     `json:"invoice_number"`
	InvoiceDate   string    `json:"invoice_date"`
}

type adminBillingOrderListResponse struct {
	Orders     []adminBillingOrderResponse `json:"orders"`
	NextCursor *string                     `json:"next_cursor"`
}

type adminBillingInvoiceListResponse struct {
	Invoices   []BillingInvoiceResponse `json:"invoices"`
	NextCursor *string                  `json:"next_cursor"`
}

func (h *Handler) ListAdminBillingOrders(c *gin.Context) {
	userID, err := parseOptionalUserID(c.Query("user_id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	status, err := parseBillingOrderStatus(c.Query("status"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	withoutInvoice, err := parseOptionalBoolQuery(c.Query("without_invoice"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid without_invoice")
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

	page, err := billing.ListAdminBillingOrders(c.Request.Context(), h.DB, billing.ListAdminBillingOrdersInput{
		UserID:         userID,
		Status:         status,
		CreatedFrom:    createdFrom,
		CreatedTo:      createdTo,
		Cursor:         cursor,
		Size:           size,
		Sort:           sortDirection,
		WithoutInvoice: withoutInvoice,
	})
	if err != nil {
		if message, ok := billingOrderListValidationError(err); ok {
			writeError(c, http.StatusBadRequest, message)
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to list admin billing orders")
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

	out := make([]adminBillingOrderResponse, 0, len(page.Orders))
	for _, order := range page.Orders {
		out = append(out, adminBillingOrderJSON(order))
	}
	c.JSON(http.StatusOK, adminBillingOrderListResponse{
		Orders:     out,
		NextCursor: nextCursor,
	})
}

func (h *Handler) ListAdminBillingInvoices(c *gin.Context) {
	userID, err := parseOptionalUserID(c.Query("user_id"))
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
	cursor, err := parseBillingInvoiceCursor(c.Query("cursor"), sortDirection)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	page, err := billing.ListAdminBillingInvoices(c.Request.Context(), h.DB, billing.ListAdminBillingInvoicesInput{
		Search:      c.Query("search"),
		UserID:      userID,
		CreatedFrom: createdFrom,
		CreatedTo:   createdTo,
		Cursor:      cursor,
		Size:        size,
		Sort:        sortDirection,
	})
	if err != nil {
		if message, ok := billingInvoiceListValidationError(err); ok {
			writeError(c, http.StatusBadRequest, message)
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to list admin billing invoices")
		return
	}

	var nextCursor *string
	if page.NextCursor != nil {
		encoded, err := encodeBillingInvoiceCursor(*page.NextCursor)
		if err != nil {
			writeError(c, http.StatusInternalServerError, "failed to encode next cursor")
			return
		}
		nextCursor = &encoded
	}

	out := make([]BillingInvoiceResponse, 0, len(page.Invoices))
	for _, invoice := range page.Invoices {
		out = append(out, billingInvoiceResponse(invoice))
	}
	c.JSON(http.StatusOK, adminBillingInvoiceListResponse{
		Invoices:   out,
		NextCursor: nextCursor,
	})
}

func (h *Handler) CreateAdminBillingOrderInvoice(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid order id")
		return
	}
	invoiceDate, err := currentBillingInvoiceDate(h.now())
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to determine invoice date")
		return
	}

	invoice, err := billing.CreateBillingInvoiceForPaidOrder(c.Request.Context(), h.DB, orderID, invoiceDate)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			writeError(c, http.StatusNotFound, "billing order not found")
		case errors.Is(err, billing.ErrBillingOrderNotPaid):
			writeError(c, http.StatusBadRequest, "billing order is not paid")
		case errors.Is(err, billing.ErrBillingInvoiceExists):
			writeError(c, http.StatusConflict, "billing order already has an invoice")
		default:
			writeError(c, http.StatusInternalServerError, "failed to create billing invoice")
		}
		return
	}

	c.JSON(http.StatusCreated, billingInvoiceResponse(invoice))
}

func parseOptionalBoolQuery(raw string) (bool, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return false, nil
	}
	return strconv.ParseBool(value)
}

func adminBillingOrderJSON(order billing.BillingOrder) adminBillingOrderResponse {
	var invoice *adminBillingOrderInvoiceResponse
	if order.Invoice != nil {
		invoice = &adminBillingOrderInvoiceResponse{
			ID:            order.Invoice.ID,
			InvoiceSerie:  order.Invoice.InvoiceSerie,
			InvoiceNumber: order.Invoice.InvoiceNumber,
			InvoiceDate:   order.Invoice.InvoiceDate.Format("2006-01-02"),
		}
	}
	return adminBillingOrderResponse{
		BillingOrderResponse: billingOrderResponse(order),
		User: adminBillingOrderUserResponse{
			ID:    userIDString(order.UserID),
			Name:  order.User.Name,
			Email: order.User.Email,
		},
		Invoice: invoice,
	}
}

func (h *Handler) now() time.Time {
	if h.Now != nil {
		return h.Now()
	}
	return time.Now().UTC()
}

func currentBillingInvoiceDate(now time.Time) (time.Time, error) {
	location, err := time.LoadLocation(billingInvoiceBusinessLocation)
	if err != nil {
		return time.Time{}, err
	}
	year, month, day := now.In(location).Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC), nil
}
