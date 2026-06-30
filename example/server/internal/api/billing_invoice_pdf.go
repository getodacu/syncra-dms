package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"ai.ro/syncra/assets"
	"ai.ro/syncra/internal/billing"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const maxGotenbergPDFBytes int64 = 50 << 20

var errGotenbergConversion = errors.New("gotenberg conversion failed")
var errRenderBillingInvoicePDF = errors.New("render billing invoice PDF failed")
var errConvertBillingInvoicePDF = errors.New("convert billing invoice PDF failed")
var errResolveBillingInvoicePDFStorage = errors.New("resolve billing invoice PDF directory failed")
var errSaveBillingInvoicePDF = errors.New("save billing invoice PDF failed")
var errUpdateBillingInvoicePDFPath = errors.New("update billing invoice PDF path failed")
var errReloadBillingInvoicePDF = errors.New("reload billing invoice PDF failed")

type generateInvoicePDFRequest struct {
	InvoiceID string `json:"invoice_id"`
}

type invoicePDFTemplateData struct {
	Seller        invoicePDFSeller
	Buyer         invoicePDFBuyer
	InvoiceID     string
	InvoiceLabel  string
	InvoiceDate   string
	InvoiceSerie  string
	InvoiceNumber int64
	Lines         []billing.BillingInvoiceLine
	NetAmount     string
	VATAmount     string
	TotalAmount   string
	GeneratedAt   string
}

type invoicePDFSeller struct {
	Name               string
	AddressLine1       string
	AddressLine2       string
	City               string
	CountryCode        string
	BillingEmail       string
	FiscalCode         string
	RegistrationNumber string
	BankAccount        string
}

type invoicePDFBuyer struct {
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

// GenerateBillingInvoicePDF generates and stores a PDF rendering for an existing invoice.
func (h *Handler) GenerateBillingInvoicePDF(c *gin.Context) {
	if !h.trustedInternalRequest(c) {
		writeError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req generateInvoicePDFRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	invoiceID, err := parseGenerateInvoicePDFID(req.InvoiceID)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	if h.DB == nil {
		writeError(c, http.StatusInternalServerError, "database is not configured")
		return
	}

	var invoice billing.BillingInvoice
	if err := h.DB.WithContext(c.Request.Context()).First(&invoice, "id = ?", invoiceID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusNotFound, "billing invoice not found")
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to load billing invoice")
		return
	}

	invoice, err = h.generateAndStoreBillingInvoicePDF(c.Request.Context(), invoice)
	if err != nil {
		switch {
		case errors.Is(err, errRenderBillingInvoicePDF):
			writeError(c, http.StatusInternalServerError, "failed to render billing invoice")
		case errors.Is(err, errGotenbergConversion):
			writeError(c, http.StatusBadGateway, "failed to convert invoice PDF")
		case errors.Is(err, errConvertBillingInvoicePDF):
			writeError(c, http.StatusInternalServerError, "failed to convert invoice PDF")
		case errors.Is(err, errResolveBillingInvoicePDFStorage):
			writeError(c, http.StatusInternalServerError, "failed to resolve billing invoice PDF directory")
		case errors.Is(err, errSaveBillingInvoicePDF):
			writeError(c, http.StatusInternalServerError, "failed to save invoice PDF")
		case errors.Is(err, errUpdateBillingInvoicePDFPath):
			writeError(c, http.StatusInternalServerError, "failed to update billing invoice PDF path")
		case errors.Is(err, errReloadBillingInvoicePDF):
			writeError(c, http.StatusInternalServerError, "failed to reload billing invoice")
		default:
			writeError(c, http.StatusInternalServerError, "failed to generate invoice PDF")
		}
		return
	}

	c.JSON(http.StatusOK, billingInvoiceResponse(invoice))
}

func (h *Handler) ServeBillingInvoicePDF(c *gin.Context) {
	if !h.trustedInternalRequest(c) {
		writeError(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	filename := strings.TrimSpace(c.Param("invoice_pdf"))
	invoiceID, err := parseBillingInvoicePDFFilename(filename)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	h.serveBillingInvoicePDF(c, invoiceID, "")
}

func (h *Handler) ServeUserBillingInvoicePDF(c *gin.Context) {
	if !h.trustedInternalRequest(c) {
		writeError(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	invoiceID, err := uuid.Parse(strings.TrimSpace(c.Param("id")))
	if err != nil || invoiceID == uuid.Nil {
		writeError(c, http.StatusBadRequest, "invalid invoice id")
		return
	}
	userID, err := parseRequiredUserID(c.Query("user_id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	h.serveBillingInvoicePDF(c, invoiceID, userID)
}

func (h *Handler) serveBillingInvoicePDF(c *gin.Context, invoiceID uuid.UUID, userID string) {
	if h.DB == nil {
		writeError(c, http.StatusInternalServerError, "database is not configured")
		return
	}

	var invoice billing.BillingInvoice
	query := h.DB.WithContext(c.Request.Context()).Where("id = ?", invoiceID)
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if err := query.First(&invoice).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusNotFound, "billing invoice PDF not found")
			return
		}
		writeError(c, http.StatusInternalServerError, "failed to load billing invoice")
		return
	}
	if invoice.PDFPath == nil || strings.TrimSpace(*invoice.PDFPath) == "" {
		writeError(c, http.StatusNotFound, "billing invoice PDF not found")
		return
	}
	pdfDir, err := h.billingInvoicePDFDir()
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to resolve billing invoice PDF directory")
		return
	}
	filename := invoice.ID.String() + ".pdf"
	pdfPath := filepath.Join(pdfDir, filename)
	info, err := os.Stat(pdfPath)
	if err != nil || info.IsDir() {
		writeError(c, http.StatusNotFound, "billing invoice PDF not found")
		return
	}

	disposition := "inline"
	if c.Query("download") == "1" {
		disposition = "attachment"
	}
	servedFilename := billingInvoicePDFServedFilename(invoice)
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf(`%s; filename="%s"`, disposition, servedFilename))
	c.File(pdfPath)
}

func parseGenerateInvoicePDFID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(strings.TrimSpace(raw))
	if err != nil || id == uuid.Nil {
		return uuid.Nil, errors.New("invalid invoice_id")
	}
	return id, nil
}

func parseBillingInvoicePDFFilename(filename string) (uuid.UUID, error) {
	if !strings.HasSuffix(filename, ".pdf") {
		return uuid.Nil, errors.New("invalid invoice PDF filename")
	}
	id, err := uuid.Parse(strings.TrimSuffix(filename, ".pdf"))
	if err != nil || id == uuid.Nil {
		return uuid.Nil, errors.New("invalid invoice PDF filename")
	}
	return id, nil
}

func billingInvoicePDFServedFilename(invoice billing.BillingInvoice) string {
	return fmt.Sprintf(
		"%s_%05d_%s.pdf",
		invoice.InvoiceSerie,
		invoice.InvoiceNumber,
		invoice.InvoiceDate.Format("060102"),
	)
}

func (h *Handler) generateAndStoreBillingInvoicePDF(ctx context.Context, invoice billing.BillingInvoice) (billing.BillingInvoice, error) {
	html, err := h.renderBillingInvoiceHTML(invoice)
	if err != nil {
		return billing.BillingInvoice{}, fmt.Errorf("%w: %w", errRenderBillingInvoicePDF, err)
	}
	pdf, err := h.convertHTMLToPDF(ctx, invoice.ID, html)
	if err != nil {
		return billing.BillingInvoice{}, fmt.Errorf("%w: %w", errConvertBillingInvoicePDF, err)
	}
	pdfDir, err := h.billingInvoicePDFDir()
	if err != nil {
		return billing.BillingInvoice{}, fmt.Errorf("%w: %w", errResolveBillingInvoicePDFStorage, err)
	}
	pdfPath, err := writeBillingInvoicePDF(pdfDir, invoice.ID, pdf)
	if err != nil {
		return billing.BillingInvoice{}, fmt.Errorf("%w: %w", errSaveBillingInvoicePDF, err)
	}
	if err := h.DB.WithContext(ctx).
		Model(&billing.BillingInvoice{}).
		Where("id = ?", invoice.ID).
		Update("pdf_path", pdfPath).Error; err != nil {
		return billing.BillingInvoice{}, fmt.Errorf("%w: %w", errUpdateBillingInvoicePDFPath, err)
	}
	if err := h.DB.WithContext(ctx).First(&invoice, "id = ?", invoice.ID).Error; err != nil {
		return billing.BillingInvoice{}, fmt.Errorf("%w: %w", errReloadBillingInvoicePDF, err)
	}
	return invoice, nil
}

func (h *Handler) billingInvoicePDFNeedsGeneration(invoice billing.BillingInvoice) (bool, error) {
	if invoice.PDFPath == nil || strings.TrimSpace(*invoice.PDFPath) == "" {
		return true, nil
	}
	pdfDir, err := h.billingInvoicePDFDir()
	if err != nil {
		return false, err
	}
	pdfPath := filepath.Join(pdfDir, invoice.ID.String()+".pdf")
	info, err := os.Stat(pdfPath)
	if err == nil {
		return info.IsDir(), nil
	}
	if os.IsNotExist(err) {
		return true, nil
	}
	return false, err
}

func (h *Handler) renderBillingInvoiceHTML(invoice billing.BillingInvoice) ([]byte, error) {
	tmpl, err := template.ParseFS(assets.FS, assets.InvoiceTemplatePath)
	if err != nil {
		return nil, err
	}
	data, err := billingInvoicePDFTemplateData(invoice)
	if err != nil {
		return nil, err
	}
	var out bytes.Buffer
	if err := tmpl.Execute(&out, data); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func billingInvoicePDFTemplateData(invoice billing.BillingInvoice) (invoicePDFTemplateData, error) {
	var lines []billing.BillingInvoiceLine
	if err := json.Unmarshal(invoice.Lines, &lines); err != nil {
		return invoicePDFTemplateData{}, err
	}
	var buyer invoicePDFBuyer
	if len(invoice.BillingProfileSnapshot) > 0 {
		_ = json.Unmarshal(invoice.BillingProfileSnapshot, &buyer)
	}
	if buyer.BillingName == "" {
		buyer.BillingName = invoice.BillingName
	}
	if buyer.BillingEmail == "" {
		buyer.BillingEmail = invoice.BillingEmail
	}
	if buyer.FiscalCode == nil {
		buyer.FiscalCode = invoice.BillingFiscalCode
	}
	return invoicePDFTemplateData{
		Seller:        defaultInvoicePDFSeller(),
		Buyer:         buyer,
		InvoiceID:     invoice.ID.String(),
		InvoiceLabel:  fmt.Sprintf("%s-%05d", invoice.InvoiceSerie, invoice.InvoiceNumber),
		InvoiceDate:   invoice.InvoiceDate.Format("2006-01-02"),
		InvoiceSerie:  invoice.InvoiceSerie,
		InvoiceNumber: invoice.InvoiceNumber,
		Lines:         lines,
		NetAmount:     invoice.NetAmount.StringFixed(2),
		VATAmount:     invoice.VATAmount.StringFixed(2),
		TotalAmount:   invoice.TotalAmount.StringFixed(2),
		GeneratedAt:   time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func defaultInvoicePDFSeller() invoicePDFSeller {
	return invoicePDFSeller{
		Name:               "Syncra SRL (placeholder)",
		AddressLine1:       "1 Automation Avenue",
		AddressLine2:       "Suite 100",
		City:               "Bucharest",
		CountryCode:        "RO",
		BillingEmail:       "billing@syncra.example",
		FiscalCode:         "RO00000000",
		RegistrationNumber: "J40/0000/2026",
		BankAccount:        "RO00BANK0000000000000000",
	}
}

func (h *Handler) convertHTMLToPDF(ctx context.Context, invoiceID uuid.UUID, html []byte) ([]byte, error) {
	endpoint, err := gotenbergConvertHTMLURL(h.GotenbergAPIURL)
	if err != nil {
		return nil, err
	}
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("files", "index.html")
	if err != nil {
		return nil, err
	}
	if _, err := part.Write(html); err != nil {
		return nil, err
	}
	if err := writer.WriteField("printBackground", "true"); err != nil {
		return nil, err
	}
	if err := writer.WriteField("preferCssPageSize", "true"); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, &body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Gotenberg-Output-Filename", invoiceID.String())

	client := &http.Client{Timeout: 60 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errGotenbergConversion, err)
	}
	defer res.Body.Close()
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		_, _ = io.Copy(io.Discard, io.LimitReader(res.Body, 4<<10))
		return nil, fmt.Errorf("%w: status %d", errGotenbergConversion, res.StatusCode)
	}
	mediaType, _, err := mime.ParseMediaType(res.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/pdf" {
		_, _ = io.Copy(io.Discard, io.LimitReader(res.Body, 4<<10))
		return nil, fmt.Errorf("%w: unexpected content type %q", errGotenbergConversion, res.Header.Get("Content-Type"))
	}
	pdf, err := readLimited(res.Body, maxGotenbergPDFBytes)
	if err != nil {
		if errors.Is(err, errReaderTooLarge) {
			return nil, fmt.Errorf("%w: PDF response too large", errGotenbergConversion)
		}
		return nil, err
	}
	return pdf, nil
}

func gotenbergConvertHTMLURL(raw string) (string, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return "", errors.New("GOTENBERG_API_URL is required")
	}
	if !strings.Contains(value, "://") {
		value = "http://" + value
	}
	parsed, err := url.Parse(value)
	if err != nil || parsed.Host == "" {
		return "", errors.New("invalid GOTENBERG_API_URL")
	}
	parsed.Path = "/forms/chromium/convert/html"
	parsed.RawQuery = ""
	parsed.Fragment = ""
	return parsed.String(), nil
}

func writeBillingInvoicePDF(dir string, invoiceID uuid.UUID, data []byte) (string, error) {
	if invoiceID == uuid.Nil {
		return "", errors.New("billing invoice id is required")
	}
	dir, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	dir = filepath.Clean(dir)
	if err := ensureDirDurable(dir); err != nil {
		return "", err
	}
	finalPath := filepath.Join(dir, invoiceID.String()+".pdf")
	tmp, err := os.CreateTemp(dir, invoiceID.String()+"-*.tmp")
	if err != nil {
		return "", err
	}
	tmpPath := tmp.Name()
	success := false
	defer func() {
		if !success {
			_ = os.Remove(tmpPath)
		}
	}()
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return "", err
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return "", err
	}
	if err := tmp.Close(); err != nil {
		return "", err
	}
	if err := os.Rename(tmpPath, finalPath); err != nil {
		return "", err
	}
	if err := syncPathFunc(dir); err != nil {
		return "", err
	}
	success = true
	return finalPath, nil
}

var errReaderTooLarge = errors.New("reader too large")

func readLimited(r io.Reader, limit int64) ([]byte, error) {
	data, err := io.ReadAll(io.LimitReader(r, limit+1))
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > limit {
		return nil, errReaderTooLarge
	}
	return data, nil
}
