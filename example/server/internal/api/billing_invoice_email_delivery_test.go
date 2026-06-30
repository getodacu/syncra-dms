package api

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"ai.ro/syncra/internal/billing"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestClaimBillingInvoiceEmailDeliveryRequiresInternalToken(t *testing.T) {
	router, _ := testBillingPDFRouter(t, "", t.TempDir())

	w := billingJSON(
		t,
		router,
		http.MethodPost,
		"/api/billing/invoices/"+uuid.NewString()+"/email-delivery/claim",
		`{"user_id":"`+uuid.NewString()+`"}`,
		"",
	)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestClaimBillingInvoiceEmailDeliveryRejectsInvalidRequests(t *testing.T) {
	router, _ := testBillingPDFRouter(t, "", t.TempDir())
	validUserID := uuid.NewString()

	for _, tt := range []struct {
		name string
		path string
		body string
	}{
		{name: "invalid invoice id", path: "/api/billing/invoices/not-a-uuid/email-delivery/claim", body: `{"user_id":"` + validUserID + `"}`},
		{name: "nil invoice id", path: "/api/billing/invoices/" + uuid.Nil.String() + "/email-delivery/claim", body: `{"user_id":"` + validUserID + `"}`},
		{name: "invalid json", path: "/api/billing/invoices/" + uuid.NewString() + "/email-delivery/claim", body: `{`},
		{name: "missing user", path: "/api/billing/invoices/" + uuid.NewString() + "/email-delivery/claim", body: `{}`},
		{name: "invalid user", path: "/api/billing/invoices/" + uuid.NewString() + "/email-delivery/claim", body: `{"user_id":"not-a-uuid"}`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			w := billingJSON(t, router, http.MethodPost, tt.path, tt.body, testInternalAPIToken)
			if w.Code != http.StatusBadRequest {
				t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
			}
		})
	}
}

func TestClaimBillingInvoiceEmailDeliveryReturnsNotReadyWithoutStoredPDF(t *testing.T) {
	router, db := testBillingPDFRouter(t, "", t.TempDir())
	userID := uuid.NewString()
	invoice := createBillingPDFTestInvoiceForUser(t, db, userID, 50)

	w := billingJSON(
		t,
		router,
		http.MethodPost,
		"/api/billing/invoices/"+invoice.ID.String()+"/email-delivery/claim",
		`{"user_id":"`+userID+`"}`,
		testInternalAPIToken,
	)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeBillingResponse[BillingInvoiceEmailDeliveryClaimResponse](t, w)
	if got.Status != billingInvoiceEmailDeliveryStatusNotReady || got.Invoice != nil {
		t.Fatalf("claim response = %#v, want not_ready without invoice", got)
	}
}

func TestClaimBillingInvoiceEmailDeliveryReturnsAlreadySent(t *testing.T) {
	outputDir := t.TempDir()
	router, db := testBillingPDFRouter(t, "", outputDir)
	userID := uuid.NewString()
	invoice := createBillingPDFTestInvoiceForUser(t, db, userID, 51)
	writeBillingInvoiceEmailDeliveryTestPDF(t, db, outputDir, invoice.ID)
	sentAt := time.Now().UTC()
	if err := db.Model(&billing.BillingInvoice{}).
		Where("id = ?", invoice.ID).
		Update("email_sent_at", sentAt).Error; err != nil {
		t.Fatalf("update email_sent_at: %v", err)
	}

	w := billingJSON(
		t,
		router,
		http.MethodPost,
		"/api/billing/invoices/"+invoice.ID.String()+"/email-delivery/claim",
		`{"user_id":"`+userID+`"}`,
		testInternalAPIToken,
	)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeBillingResponse[BillingInvoiceEmailDeliveryClaimResponse](t, w)
	if got.Status != billingInvoiceEmailDeliveryStatusAlreadySent || got.Invoice != nil {
		t.Fatalf("claim response = %#v, want already_sent without invoice", got)
	}
}

func TestClaimBillingInvoiceEmailDeliveryClaimsOnceDuringLease(t *testing.T) {
	outputDir := t.TempDir()
	router, db := testBillingPDFRouter(t, "", outputDir)
	userID := uuid.NewString()
	invoice := createBillingPDFTestInvoiceForUser(t, db, userID, 52)
	writeBillingInvoiceEmailDeliveryTestPDF(t, db, outputDir, invoice.ID)
	path := "/api/billing/invoices/" + invoice.ID.String() + "/email-delivery/claim"
	body := `{"user_id":"` + userID + `"}`

	first := billingJSON(t, router, http.MethodPost, path, body, testInternalAPIToken)
	if first.Code != http.StatusOK {
		t.Fatalf("first status = %d body=%s", first.Code, first.Body.String())
	}
	firstClaim := decodeBillingResponse[BillingInvoiceEmailDeliveryClaimResponse](t, first)
	if firstClaim.Status != billingInvoiceEmailDeliveryStatusClaimed ||
		firstClaim.Invoice == nil ||
		firstClaim.Invoice.ID != invoice.ID ||
		firstClaim.Invoice.BillingEmail != "billing@example.com" ||
		firstClaim.Invoice.PDFPath == nil ||
		firstClaim.Invoice.EmailDeliveryClaimedAt == nil {
		t.Fatalf("first claim = %#v, want claimed invoice with delivery timestamp", firstClaim)
	}

	second := billingJSON(t, router, http.MethodPost, path, body, testInternalAPIToken)
	if second.Code != http.StatusOK {
		t.Fatalf("second status = %d body=%s", second.Code, second.Body.String())
	}
	secondClaim := decodeBillingResponse[BillingInvoiceEmailDeliveryClaimResponse](t, second)
	if secondClaim.Status != billingInvoiceEmailDeliveryStatusClaimActive || secondClaim.Invoice != nil {
		t.Fatalf("second claim = %#v, want claim_active without invoice", secondClaim)
	}
}

func TestMarkBillingInvoiceEmailSentRecordsSentAt(t *testing.T) {
	outputDir := t.TempDir()
	router, db := testBillingPDFRouter(t, "", outputDir)
	userID := uuid.NewString()
	invoice := createBillingPDFTestInvoiceForUser(t, db, userID, 53)
	writeBillingInvoiceEmailDeliveryTestPDF(t, db, outputDir, invoice.ID)

	w := billingJSON(
		t,
		router,
		http.MethodPost,
		"/api/billing/invoices/"+invoice.ID.String()+"/email-delivery/sent",
		`{"user_id":"`+userID+`"}`,
		testInternalAPIToken,
	)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeBillingResponse[BillingInvoiceEmailDeliverySentResponse](t, w)
	if got.Status != billingInvoiceEmailDeliveryStatusSent || got.Invoice.EmailSentAt == nil {
		t.Fatalf("sent response = %#v, want sent timestamp", got)
	}

	var stored billing.BillingInvoice
	if err := db.First(&stored, "id = ?", invoice.ID).Error; err != nil {
		t.Fatalf("load stored invoice: %v", err)
	}
	if stored.EmailSentAt == nil {
		t.Fatal("stored EmailSentAt = nil, want timestamp")
	}
}

func writeBillingInvoiceEmailDeliveryTestPDF(t *testing.T, db *gorm.DB, outputDir string, invoiceID uuid.UUID) {
	t.Helper()
	invoiceDir := filepath.Join(outputDir, "invoices")
	if err := os.MkdirAll(invoiceDir, 0o755); err != nil {
		t.Fatalf("create output dir: %v", err)
	}
	pdfPath := filepath.Join(invoiceDir, invoiceID.String()+".pdf")
	if err := os.WriteFile(pdfPath, []byte("%PDF-email"), 0o644); err != nil {
		t.Fatalf("write pdf: %v", err)
	}
	if err := db.Model(&billing.BillingInvoice{}).
		Where("id = ?", invoiceID).
		Update("pdf_path", pdfPath).Error; err != nil {
		t.Fatalf("update pdf path: %v", err)
	}
}
