package api

import (
	"bytes"
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"ai.ro/syncra/internal/auth"
	"ai.ro/syncra/internal/billing"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGenerateBillingInvoicePDFRequiresInternalToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter(&Handler{InternalAPIToken: testInternalAPIToken})

	w := billingJSON(t, router, http.MethodPost, "/api/billing/generate-invoice-pdf", `{"invoice_id":"`+uuid.NewString()+`"}`, "")

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestGenerateBillingInvoicePDFRejectsInvalidRequests(t *testing.T) {
	router, _ := testBillingPDFRouter(t, "", t.TempDir())

	for _, tt := range []struct {
		name string
		body string
	}{
		{name: "invalid json", body: `{`},
		{name: "missing invoice id", body: `{}`},
		{name: "invalid invoice id", body: `{"invoice_id":"not-a-uuid"}`},
		{name: "nil invoice id", body: `{"invoice_id":"00000000-0000-0000-0000-000000000000"}`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			w := billingJSON(t, router, http.MethodPost, "/api/billing/generate-invoice-pdf", tt.body, testInternalAPIToken)
			if w.Code != http.StatusBadRequest {
				t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
			}
		})
	}
}

func TestGenerateBillingInvoicePDFReturnsNotFoundForUnknownInvoice(t *testing.T) {
	router, _ := testBillingPDFRouter(t, "", t.TempDir())

	w := billingJSON(t, router, http.MethodPost, "/api/billing/generate-invoice-pdf", `{"invoice_id":"`+uuid.NewString()+`"}`, testInternalAPIToken)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestGenerateBillingInvoicePDFSendsHTMLToGotenbergAndOverwritesStoredPDF(t *testing.T) {
	t.Chdir(t.TempDir())

	var requests int32
	gotenberg := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/forms/chromium/convert/html" {
			http.Error(w, "unexpected route", http.StatusNotFound)
			return
		}
		if mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type")); err != nil || mediaType != "multipart/form-data" {
			http.Error(w, "expected multipart form", http.StatusBadRequest)
			return
		}
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			http.Error(w, "parse multipart", http.StatusBadRequest)
			return
		}
		if r.FormValue("printBackground") != "true" || r.FormValue("preferCssPageSize") != "true" {
			http.Error(w, "missing render fields", http.StatusBadRequest)
			return
		}
		file, header, err := r.FormFile("files")
		if err != nil {
			http.Error(w, "missing html file", http.StatusBadRequest)
			return
		}
		defer file.Close()
		if header.Filename != "index.html" {
			http.Error(w, "wrong html filename", http.StatusBadRequest)
			return
		}
		html, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "read html", http.StatusBadRequest)
			return
		}
		if !bytes.Contains(html, []byte("Invoice SYN-00042")) || !bytes.Contains(html, []byte("ICI Bucuresti")) {
			http.Error(w, "html missing invoice data", http.StatusBadRequest)
			return
		}
		next := atomic.AddInt32(&requests, 1)
		w.Header().Set("Content-Type", "application/pdf")
		_, _ = w.Write([]byte("%PDF-" + strconv.Itoa(int(next))))
	}))
	defer gotenberg.Close()

	outputDir := t.TempDir()
	router, db := testBillingPDFRouter(t, gotenberg.URL+"/forms/chromium/convert/url", outputDir)
	invoiceDir := filepath.Join(outputDir, "invoices")
	invoice := createBillingPDFTestInvoice(t, db)
	body := `{"invoice_id":"` + invoice.ID.String() + `"}`

	w := billingJSON(t, router, http.MethodPost, "/api/billing/generate-invoice-pdf", body, testInternalAPIToken)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeBillingResponse[BillingInvoiceResponse](t, w)
	wantPath := filepath.Join(invoiceDir, invoice.ID.String()+".pdf")
	if got.PDFPath == nil || *got.PDFPath != wantPath {
		t.Fatalf("pdf_path = %#v, want %q", got.PDFPath, wantPath)
	}
	assertFileBytes(t, wantPath, []byte("%PDF-1"))

	w = billingJSON(t, router, http.MethodPost, "/api/billing/generate-invoice-pdf", body, testInternalAPIToken)
	if w.Code != http.StatusOK {
		t.Fatalf("second status = %d body=%s", w.Code, w.Body.String())
	}
	got = decodeBillingResponse[BillingInvoiceResponse](t, w)
	if got.PDFPath == nil || *got.PDFPath != wantPath {
		t.Fatalf("second pdf_path = %#v, want %q", got.PDFPath, wantPath)
	}
	assertFileBytes(t, wantPath, []byte("%PDF-2"))

	var stored billing.BillingInvoice
	if err := db.First(&stored, "id = ?", invoice.ID).Error; err != nil {
		t.Fatalf("load stored invoice: %v", err)
	}
	if stored.PDFPath == nil || *stored.PDFPath != wantPath {
		t.Fatalf("stored PDFPath = %#v, want %q", stored.PDFPath, wantPath)
	}
	if atomic.LoadInt32(&requests) != 2 {
		t.Fatalf("gotenberg requests = %d, want 2", requests)
	}
}

func TestGenerateBillingInvoicePDFDoesNotUpdatePathWhenGotenbergFails(t *testing.T) {
	gotenberg := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("not a pdf"))
	}))
	defer gotenberg.Close()

	outputDir := t.TempDir()
	router, db := testBillingPDFRouter(t, gotenberg.URL, outputDir)
	invoiceDir := filepath.Join(outputDir, "invoices")
	invoice := createBillingPDFTestInvoice(t, db)

	w := billingJSON(t, router, http.MethodPost, "/api/billing/generate-invoice-pdf", `{"invoice_id":"`+invoice.ID.String()+`"}`, testInternalAPIToken)
	if w.Code != http.StatusBadGateway {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}

	var stored billing.BillingInvoice
	if err := db.First(&stored, "id = ?", invoice.ID).Error; err != nil {
		t.Fatalf("load stored invoice: %v", err)
	}
	if stored.PDFPath != nil {
		t.Fatalf("stored PDFPath = %#v, want nil", stored.PDFPath)
	}
	if _, err := os.Stat(filepath.Join(invoiceDir, invoice.ID.String()+".pdf")); !os.IsNotExist(err) {
		t.Fatalf("pdf file stat err = %v, want not exist", err)
	}
}

func TestServeBillingInvoicePDFRequiresInternalToken(t *testing.T) {
	router, _ := testBillingPDFRouter(t, "", t.TempDir())
	invoiceID := uuid.NewString()

	req := httptest.NewRequest(http.MethodGet, "/static/invoice/"+invoiceID+".pdf", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestServeBillingInvoicePDFRejectsInvalidFilename(t *testing.T) {
	router, _ := testBillingPDFRouter(t, "", t.TempDir())

	for _, filename := range []string{
		"not-a-uuid.pdf",
		uuid.NewString() + ".txt",
		uuid.Nil.String() + ".pdf",
	} {
		t.Run(filename, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/static/invoice/"+filename, nil)
			req.Header.Set(internalAPIHeader, testInternalAPIToken)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
			}
		})
	}
}

func TestServeBillingInvoicePDFReturnsNotFound(t *testing.T) {
	outputDir := t.TempDir()
	router, db := testBillingPDFRouter(t, "", outputDir)
	invoiceDir := filepath.Join(outputDir, "invoices")

	missingInvoice := uuid.NewString()
	assertBillingInvoicePDFNotFound(t, router, missingInvoice)

	withoutPDFPath := createBillingPDFTestInvoiceWithNumber(t, db, 43)
	assertBillingInvoicePDFNotFound(t, router, withoutPDFPath.ID.String())

	missingPath := filepath.Join(invoiceDir, uuid.NewString()+".pdf")
	withMissingFile := createBillingPDFTestInvoiceWithNumber(t, db, 44)
	if err := db.Model(&billing.BillingInvoice{}).
		Where("id = ?", withMissingFile.ID).
		Update("pdf_path", missingPath).Error; err != nil {
		t.Fatalf("update pdf path: %v", err)
	}
	assertBillingInvoicePDFNotFound(t, router, withMissingFile.ID.String())
}

func TestServeBillingInvoicePDFServesInlineAndAttachment(t *testing.T) {
	outputDir := t.TempDir()
	router, db := testBillingPDFRouter(t, "", outputDir)
	invoiceDir := filepath.Join(outputDir, "invoices")
	invoice := createBillingPDFTestInvoice(t, db)
	pdfPath := filepath.Join(invoiceDir, invoice.ID.String()+".pdf")
	if err := os.MkdirAll(invoiceDir, 0o755); err != nil {
		t.Fatalf("create output dir: %v", err)
	}
	if err := os.WriteFile(pdfPath, []byte("%PDF-test"), 0o644); err != nil {
		t.Fatalf("write pdf: %v", err)
	}
	if err := db.Model(&billing.BillingInvoice{}).
		Where("id = ?", invoice.ID).
		Update("pdf_path", pdfPath).Error; err != nil {
		t.Fatalf("update pdf path: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/static/invoice/"+invoice.ID.String()+".pdf", nil)
	req.Header.Set(internalAPIHeader, testInternalAPIToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if got := w.Header().Get("Content-Type"); !strings.Contains(got, "application/pdf") {
		t.Fatalf("Content-Type = %q, want application/pdf", got)
	}
	wantFilename := "SYN_00042_260610.pdf"
	if got := w.Header().Get("Content-Disposition"); got != `inline; filename="`+wantFilename+`"` {
		t.Fatalf("Content-Disposition = %q", got)
	}
	if got := w.Body.Bytes(); !bytes.Equal(got, []byte("%PDF-test")) {
		t.Fatalf("body = %q, want PDF bytes", got)
	}

	req = httptest.NewRequest(http.MethodGet, "/static/invoice/"+invoice.ID.String()+".pdf?download=1", nil)
	req.Header.Set(internalAPIHeader, testInternalAPIToken)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("download status = %d body=%s", w.Code, w.Body.String())
	}
	if got := w.Header().Get("Content-Disposition"); got != `attachment; filename="`+wantFilename+`"` {
		t.Fatalf("download Content-Disposition = %q", got)
	}
}

func TestServeUserBillingInvoicePDFRequiresInternalToken(t *testing.T) {
	router, _ := testBillingPDFRouter(t, "", t.TempDir())

	req := httptest.NewRequest(http.MethodGet, "/api/billing/invoices/"+uuid.NewString()+"/pdf?user_id="+uuid.NewString(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestServeUserBillingInvoicePDFRejectsInvalidRequests(t *testing.T) {
	router, _ := testBillingPDFRouter(t, "", t.TempDir())
	userID := uuid.NewString()

	for _, tt := range []struct {
		name string
		path string
	}{
		{name: "invalid invoice id", path: "/api/billing/invoices/not-a-uuid/pdf?user_id=" + userID},
		{name: "nil invoice id", path: "/api/billing/invoices/" + uuid.Nil.String() + "/pdf?user_id=" + userID},
		{name: "missing user id", path: "/api/billing/invoices/" + uuid.NewString() + "/pdf"},
		{name: "invalid user id", path: "/api/billing/invoices/" + uuid.NewString() + "/pdf?user_id=not-a-uuid"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			req.Header.Set(internalAPIHeader, testInternalAPIToken)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
			}
		})
	}
}

func TestServeUserBillingInvoicePDFReturnsNotFoundForUnownedOrPDFLessInvoice(t *testing.T) {
	outputDir := t.TempDir()
	router, db := testBillingPDFRouter(t, "", outputDir)
	invoiceDir := filepath.Join(outputDir, "invoices")
	user := createTestUser(t, db, "invoice-pdf-user@example.com")
	other := createTestUser(t, db, "invoice-pdf-other@example.com")

	otherInvoice := createBillingPDFTestInvoiceForUser(t, db, other.ID, 43)
	pdfPath := filepath.Join(invoiceDir, otherInvoice.ID.String()+".pdf")
	if err := os.MkdirAll(invoiceDir, 0o755); err != nil {
		t.Fatalf("create output dir: %v", err)
	}
	if err := os.WriteFile(pdfPath, []byte("%PDF-other"), 0o644); err != nil {
		t.Fatalf("write pdf: %v", err)
	}
	if err := db.Model(&billing.BillingInvoice{}).
		Where("id = ?", otherInvoice.ID).
		Update("pdf_path", pdfPath).Error; err != nil {
		t.Fatalf("update other invoice pdf path: %v", err)
	}

	pdfLess := createBillingPDFTestInvoiceForUser(t, db, user.ID, 44)

	for _, invoiceID := range []string{otherInvoice.ID.String(), pdfLess.ID.String(), uuid.NewString()} {
		req := httptest.NewRequest(http.MethodGet, "/api/billing/invoices/"+invoiceID+"/pdf?user_id="+user.ID, nil)
		req.Header.Set(internalAPIHeader, testInternalAPIToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("invoice %s status = %d body=%s", invoiceID, w.Code, w.Body.String())
		}
	}
}

func TestServeUserBillingInvoicePDFServesInlineAndAttachment(t *testing.T) {
	outputDir := t.TempDir()
	router, db := testBillingPDFRouter(t, "", outputDir)
	invoiceDir := filepath.Join(outputDir, "invoices")
	user := createTestUser(t, db, "invoice-pdf-owner@example.com")
	invoice := createBillingPDFTestInvoiceForUser(t, db, user.ID, 45)
	pdfPath := filepath.Join(invoiceDir, invoice.ID.String()+".pdf")
	if err := os.MkdirAll(invoiceDir, 0o755); err != nil {
		t.Fatalf("create output dir: %v", err)
	}
	if err := os.WriteFile(pdfPath, []byte("%PDF-user"), 0o644); err != nil {
		t.Fatalf("write pdf: %v", err)
	}
	if err := db.Model(&billing.BillingInvoice{}).
		Where("id = ?", invoice.ID).
		Update("pdf_path", pdfPath).Error; err != nil {
		t.Fatalf("update pdf path: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/billing/invoices/"+invoice.ID.String()+"/pdf?user_id="+user.ID, nil)
	req.Header.Set(internalAPIHeader, testInternalAPIToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if got := w.Header().Get("Content-Type"); !strings.Contains(got, "application/pdf") {
		t.Fatalf("Content-Type = %q, want application/pdf", got)
	}
	wantFilename := "SYN_00045_260610.pdf"
	if got := w.Header().Get("Content-Disposition"); got != `inline; filename="`+wantFilename+`"` {
		t.Fatalf("Content-Disposition = %q", got)
	}
	if got := w.Body.Bytes(); !bytes.Equal(got, []byte("%PDF-user")) {
		t.Fatalf("body = %q, want PDF bytes", got)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/billing/invoices/"+invoice.ID.String()+"/pdf?user_id="+user.ID+"&download=1", nil)
	req.Header.Set(internalAPIHeader, testInternalAPIToken)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("download status = %d body=%s", w.Code, w.Body.String())
	}
	if got := w.Header().Get("Content-Disposition"); got != `attachment; filename="`+wantFilename+`"` {
		t.Fatalf("download Content-Disposition = %q", got)
	}
}

func assertBillingInvoicePDFNotFound(t *testing.T, router http.Handler, invoiceID string) {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/static/invoice/"+invoiceID+".pdf", nil)
	req.Header.Set(internalAPIHeader, testInternalAPIToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func testBillingPDFRouter(t *testing.T, gotenbergURL string, outputDir string) (*gin.Engine, *gorm.DB) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("sqlite db handle: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	if err := db.AutoMigrate(
		&auth.User{},
		&billing.BillingProfile{},
		&billing.BillingOrder{},
		&billing.BillingInvoiceCounter{},
		&billing.BillingInvoice{},
	); err != nil {
		t.Fatalf("auto migrate sqlite: %v", err)
	}
	router := NewRouter(&Handler{
		DB:               db,
		InternalAPIToken: testInternalAPIToken,
		GotenbergAPIURL:  gotenbergURL,
		StorageDir:       outputDir,
	})
	return router, db
}

func createBillingPDFTestInvoice(t *testing.T, db *gorm.DB) billing.BillingInvoice {
	t.Helper()
	return createBillingPDFTestInvoiceWithNumber(t, db, 42)
}

func createBillingPDFTestInvoiceForUser(t *testing.T, db *gorm.DB, userID string, invoiceNumber int64) billing.BillingInvoice {
	t.Helper()
	invoice := createBillingPDFTestInvoiceWithNumber(t, db, invoiceNumber)
	if err := db.Model(&billing.BillingInvoice{}).
		Where("id = ?", invoice.ID).
		Update("user_id", userID).Error; err != nil {
		t.Fatalf("update invoice user id: %v", err)
	}
	var out billing.BillingInvoice
	if err := db.First(&out, "id = ?", invoice.ID).Error; err != nil {
		t.Fatalf("reload user invoice: %v", err)
	}
	return out
}

func createBillingPDFTestInvoiceWithNumber(t *testing.T, db *gorm.DB, invoiceNumber int64) billing.BillingInvoice {
	t.Helper()
	invoice := billing.BillingInvoice{
		BillingName:            "ICI Bucuresti",
		BillingEmail:           "billing@example.com",
		BillingFiscalCode:      ptrString("RO2785503"),
		BillingProfileSnapshot: datatypes.JSON([]byte(`{"billing_name":"ICI Bucuresti","address_line1":"Maresal Averescu 8-10","city":"Bucuresti","country_code":"RO"}`)),
		Lines:                  datatypes.JSON([]byte(`[{"name":"OCR credits","quantity":1,"unit_price":"10.00","vat_percentage":"19.00","total_vat_amount":"1.90","total_amount":"11.90"}]`)),
		NetAmount:              decimal.RequireFromString("10.00"),
		VATAmount:              decimal.RequireFromString("1.90"),
		TotalAmount:            decimal.RequireFromString("11.90"),
		InvoiceDate:            time.Date(2026, 6, 10, 0, 0, 0, 0, time.UTC),
		InvoiceSerie:           "SYN",
		InvoiceNumber:          invoiceNumber,
	}
	if err := db.Create(&invoice).Error; err != nil {
		t.Fatalf("create invoice: %v", err)
	}
	return invoice
}

func assertFileBytes(t *testing.T, path string, want []byte) {
	t.Helper()
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("file %s = %q, want %q", path, got, want)
	}
}

func TestBillingInvoicePDFResponseIncludesPath(t *testing.T) {
	pdfPath := "/tmp/invoice.pdf"
	invoice := billing.BillingInvoice{
		ID:                     uuid.New(),
		BillingName:            "ICI Bucuresti",
		BillingEmail:           "billing@example.com",
		BillingProfileSnapshot: datatypes.JSON([]byte(`{"billing_name":"ICI Bucuresti"}`)),
		Lines:                  datatypes.JSON([]byte(`[]`)),
		NetAmount:              decimal.Zero,
		VATAmount:              decimal.Zero,
		TotalAmount:            decimal.Zero,
		InvoiceDate:            time.Date(2026, 6, 10, 0, 0, 0, 0, time.UTC),
		InvoiceSerie:           "SYN",
		InvoiceNumber:          42,
		PDFPath:                &pdfPath,
	}

	payload, err := json.Marshal(billingInvoiceResponse(invoice))
	if err != nil {
		t.Fatalf("marshal response: %v", err)
	}
	if !strings.Contains(string(payload), `"pdf_path":"/tmp/invoice.pdf"`) {
		t.Fatalf("response JSON = %s, want pdf_path", payload)
	}
}
