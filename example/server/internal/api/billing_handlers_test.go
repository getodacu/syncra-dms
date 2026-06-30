package api

import (
	"bytes"
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"ai.ro/syncra/internal/auth"
	"ai.ro/syncra/internal/billing"
)

const testInternalAPIToken = "trusted-internal-token"

func testBillingRouter(t *testing.T) (*gin.Engine, *gorm.DB) {
	t.Helper()
	return testBillingRouterWithConfig(t, nil)
}

func testBillingRouterWithConfig(t *testing.T, configure func(*Handler)) (*gin.Engine, *gorm.DB) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	db := apiPostgresTx(t)
	handler := &Handler{
		DB:               db,
		InternalAPIToken: testInternalAPIToken,
	}
	if configure != nil {
		configure(handler)
	}
	return NewRouter(handler), db
}

func billingHandlerTestModels() []any {
	return []any{
		&auth.User{},
		&auth.APIKey{},
		&billing.BillingProfile{},
		&billing.BillingInvoiceCounter{},
		&billing.BillingInvoice{},
		&billing.BillingOrder{},
		&billing.CreditBucket{},
		&billing.CreditLedgerEntry{},
	}
}

func TestBillingMutationEndpointsRequireInternalToken(t *testing.T) {
	router, _ := testBillingRouter(t)
	orderID := uuid.New()
	for _, tt := range []struct {
		name   string
		method string
		path   string
		body   string
	}{
		{name: "get profile", method: http.MethodGet, path: "/api/billing/profile?user_id=" + uuid.NewString()},
		{name: "upsert profile", method: http.MethodPut, path: "/api/billing/profile", body: `{"user_id":"` + uuid.NewString() + `"}`},
		{name: "create invoice", method: http.MethodPost, path: "/api/billing/invoices", body: `{"user_id":"` + uuid.NewString() + `","invoice_serie":"SYNC","lines":[]}`},
		{name: "list orders", method: http.MethodGet, path: "/api/billing/orders?user_id=" + uuid.NewString()},
		{name: "create order", method: http.MethodPost, path: "/api/billing/orders", body: `{"user_id":"` + uuid.NewString() + `","credits":1000}`},
		{name: "attach checkout session", method: http.MethodPost, path: "/api/billing/orders/" + orderID.String() + "/checkout-session", body: `{"checkout_session_id":"cs_test_123"}`},
		{name: "mark paid", method: http.MethodPost, path: "/api/billing/orders/" + orderID.String() + "/paid", body: `{"paid_at":"2026-06-04T12:00:00Z"}`},
		{name: "mark failed", method: http.MethodPost, path: "/api/billing/orders/" + orderID.String() + "/failed", body: `{}`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			w := billingJSON(t, router, tt.method, tt.path, tt.body, "")
			if w.Code != http.StatusUnauthorized {
				t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
			}
		})
	}
}

func TestCreateBillingInvoiceGeneratesInvoice(t *testing.T) {
	router, db := testBillingRouter(t)
	user := createTestUser(t, db, "billing-invoice@example.com")
	fiscalCode := "RO2785503"
	if _, err := billing.UpsertBillingProfile(t.Context(), db, billing.UpsertBillingProfileInput{
		UserID:       user.ID,
		EntityType:   billing.BillingEntityCompany,
		BillingName:  "ICI Bucuresti",
		BillingEmail: "billing@example.com",
		CountryCode:  "RO",
		AddressLine1: "Maresal Averescu 8-10",
		City:         "Bucuresti",
		PostalCode:   "011455",
		FiscalCode:   &fiscalCode,
	}); err != nil {
		t.Fatalf("seed billing profile: %v", err)
	}

	body := `{"user_id":"` + user.ID + `","invoice_serie":" sync ","invoice_date":"2026-06-10","lines":[{"name":"OCR credits","quantity":2,"unit_price":"10.00","vat_percentage":"19"}]}`
	w := billingJSON(t, router, http.MethodPost, "/api/billing/invoices", body, testInternalAPIToken)

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeBillingResponse[BillingInvoiceResponse](t, w)
	if got.UserID == nil ||
		string(*got.UserID) != user.ID ||
		got.BillingName != "ICI Bucuresti" ||
		got.BillingEmail != "billing@example.com" ||
		got.BillingFiscalCode == nil ||
		*got.BillingFiscalCode != fiscalCode ||
		got.InvoiceSerie != "SYNC" ||
		got.InvoiceNumber != 1 ||
		got.InvoiceDate != "2026-06-10" ||
		got.NetAmount != "20.00" ||
		got.VATAmount != "3.80" ||
		got.TotalAmount != "23.80" {
		t.Fatalf("invoice response = %#v", got)
	}
	if len(got.Lines) != 1 ||
		got.Lines[0].Name != "OCR credits" ||
		got.Lines[0].Quantity != 2 ||
		got.Lines[0].UnitPrice != "10.00" ||
		got.Lines[0].VATPercentage != "19.00" ||
		got.Lines[0].TotalVATAmount != "3.80" ||
		got.Lines[0].TotalAmount != "23.80" {
		t.Fatalf("invoice lines = %#v", got.Lines)
	}
	if !strings.Contains(string(got.BillingProfileSnapshot), "ICI Bucuresti") {
		t.Fatalf("billing profile snapshot = %s", got.BillingProfileSnapshot)
	}

	var stored billing.BillingInvoice
	if err := db.First(&stored, "id = ?", got.ID).Error; err != nil {
		t.Fatalf("load stored invoice: %v", err)
	}
	if stored.InvoiceNumber != 1 || stored.InvoiceSerie != "SYNC" {
		t.Fatalf("stored invoice = %#v", stored)
	}
}

func TestCreateBillingInvoiceUsesUserFallbackWithoutBillingProfile(t *testing.T) {
	router, db := testBillingRouter(t)
	user := createTestUser(t, db, "billing-invoice-no-profile@example.com")

	body := `{"user_id":"` + user.ID + `","invoice_serie":" sync ","invoice_date":"2026-06-10","lines":[{"name":"OCR credits","quantity":1,"unit_price":"10.00","vat_percentage":"0"}]}`
	w := billingJSON(t, router, http.MethodPost, "/api/billing/invoices", body, testInternalAPIToken)

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeBillingResponse[BillingInvoiceResponse](t, w)
	if got.UserID == nil ||
		string(*got.UserID) != user.ID ||
		got.BillingProfileID != nil ||
		got.BillingName != user.Name ||
		got.BillingEmail != user.Email ||
		got.BillingFiscalCode != nil {
		t.Fatalf("invoice response = %#v, want user fallback buyer fields", got)
	}
	var snapshot struct {
		Source       string `json:"source"`
		UserID       string `json:"user_id"`
		BillingName  string `json:"billing_name"`
		BillingEmail string `json:"billing_email"`
	}
	if err := json.Unmarshal(got.BillingProfileSnapshot, &snapshot); err != nil {
		t.Fatalf("unmarshal fallback snapshot: %v", err)
	}
	if snapshot.Source != "user" ||
		snapshot.UserID != user.ID ||
		snapshot.BillingName != user.Name ||
		snapshot.BillingEmail != user.Email {
		t.Fatalf("fallback snapshot = %#v, want user buyer data", snapshot)
	}
}

func TestCreateBillingInvoiceRejectsInvalidRequests(t *testing.T) {
	router, db := testBillingRouter(t)
	user := createTestUser(t, db, "billing-invoice-invalid@example.com")
	if _, err := billing.UpsertBillingProfile(t.Context(), db, billing.UpsertBillingProfileInput{
		UserID:       user.ID,
		EntityType:   billing.BillingEntityIndividual,
		BillingName:  "Radu Boncea",
		BillingEmail: "radu@example.com",
		CountryCode:  "RO",
		AddressLine1: "Maresal Averescu 8-10",
		City:         "Bucuresti",
		PostalCode:   "011455",
	}); err != nil {
		t.Fatalf("seed billing profile: %v", err)
	}

	validBody := func() string {
		return `{"user_id":"` + user.ID + `","invoice_serie":"SYNC","invoice_date":"2026-06-10","lines":[{"name":"OCR credits","quantity":1,"unit_price":"10.00","vat_percentage":"19.00"}]}`
	}

	for _, tt := range []struct {
		name     string
		body     string
		wantBody string
	}{
		{name: "invalid json", body: `{`, wantBody: "invalid JSON body"},
		{name: "invalid user", body: strings.Replace(validBody(), user.ID, uuid.NewString(), 1), wantBody: "invalid user_id"},
		{name: "invalid invoice date", body: strings.Replace(validBody(), `"2026-06-10"`, `"06/10/2026"`, 1), wantBody: "invalid invoice_date"},
		{name: "empty lines", body: strings.Replace(validBody(), `[{"name":"OCR credits","quantity":1,"unit_price":"10.00","vat_percentage":"19.00"}]`, `[]`, 1), wantBody: "invoice lines are required"},
		{name: "invalid decimal", body: strings.Replace(validBody(), `"10.00"`, `"ten"`, 1), wantBody: "unit price is invalid"},
		{name: "invalid serie", body: strings.Replace(validBody(), `"SYNC"`, `"SYNC/2026"`, 1), wantBody: "invoice serie contains invalid characters"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			w := billingJSON(t, router, http.MethodPost, "/api/billing/invoices", tt.body, testInternalAPIToken)
			if w.Code != http.StatusBadRequest || !strings.Contains(w.Body.String(), tt.wantBody) {
				t.Fatalf("status = %d body=%s, want %q", w.Code, w.Body.String(), tt.wantBody)
			}
		})
	}
}

func TestGetCreditBalanceReturnsAvailableCredits(t *testing.T) {
	router, db := testBillingRouter(t)
	user := createTestUser(t, db, "billing-balance@example.com")
	now := time.Now().UTC().Add(-time.Hour)
	if _, err := billing.GrantSignupBonus(t.Context(), db, billing.GrantSignupBonusInput{
		UserID:  user.ID,
		Credits: 100,
		Now:     now,
	}); err != nil {
		t.Fatalf("grant signup bonus: %v", err)
	}

	req := newTestRequest(http.MethodGet, "/api/billing/balance?user_id="+user.ID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeBillingResponse[CreditBalanceResponse](t, w)
	if string(got.UserID) != user.ID || got.AvailableCredits != 100 {
		t.Fatalf("balance response = %#v, want user %s available 100", got, user.ID)
	}
}

func TestGetPublicCreditBalanceReturnsAPIKeyUserAvailableCredits(t *testing.T) {
	router, db := testBillingRouter(t)
	user := createTestUser(t, db, "public-billing-balance@example.com")
	other := createTestUser(t, db, "public-billing-balance-other@example.com")
	secret := createTestPublicAPIKey(t, db, user.ID, nil)
	grantTestCredits(t, db, user.ID, 750)
	grantTestCredits(t, db, other.ID, 900)

	req := httptest.NewRequest(http.MethodGet, "/v1/get-balance", nil)
	authorizePublicAPIRequest(req, secret)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeBillingResponse[CreditBalanceResponse](t, w)
	if string(got.UserID) != user.ID || got.AvailableCredits != 750 {
		t.Fatalf("balance response = %#v, want user %s available 750", got, user.ID)
	}
}

func TestGetPublicCreditBalanceReturnsZeroWithoutCredits(t *testing.T) {
	router, db := testBillingRouter(t)
	user := createTestUser(t, db, "public-billing-balance-empty@example.com")
	other := createTestUser(t, db, "public-billing-balance-empty-other@example.com")
	secret := createTestPublicAPIKey(t, db, user.ID, nil)
	grantTestCredits(t, db, other.ID, 500)

	req := httptest.NewRequest(http.MethodGet, "/v1/get-balance", nil)
	authorizePublicAPIRequest(req, secret)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeBillingResponse[CreditBalanceResponse](t, w)
	if string(got.UserID) != user.ID || got.AvailableCredits != 0 {
		t.Fatalf("balance response = %#v, want user %s available 0", got, user.ID)
	}
}

func TestGetPublicCreditBalanceRejectsUserIDQuery(t *testing.T) {
	router, db := testBillingRouter(t)
	user := createTestUser(t, db, "public-billing-balance-query@example.com")
	secret := createTestPublicAPIKey(t, db, user.ID, nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/get-balance?user_id="+url.QueryEscape(user.ID), nil)
	authorizePublicAPIRequest(req, secret)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestGetPublicCreditBalanceRejectsInvalidAPIKeys(t *testing.T) {
	router, db := testBillingRouter(t)
	user := createTestUser(t, db, "public-billing-balance-auth@example.com")
	expiredAt := time.Now().UTC().Add(-time.Minute)
	expiredSecret := createTestPublicAPIKey(t, db, user.ID, &expiredAt)

	for _, tt := range []struct {
		name          string
		authorization string
	}{
		{name: "missing"},
		{name: "malformed", authorization: "Bearer one two"},
		{name: "unknown", authorization: "Bearer public-api-key-unknown"},
		{name: "expired", authorization: "Bearer " + expiredSecret},
	} {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/v1/get-balance", nil)
			if tt.authorization != "" {
				req.Header.Set("Authorization", tt.authorization)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusUnauthorized {
				t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
			}
		})
	}
}

func TestGetBillingProfileReturnsNullWhenMissing(t *testing.T) {
	router, db := testBillingRouter(t)
	user := createTestUser(t, db, "billing-profile-missing@example.com")

	req := httptest.NewRequest(http.MethodGet, "/api/billing/profile?user_id="+user.ID, nil)
	req.Header.Set("X-Syncra-Internal-Token", testInternalAPIToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeBillingResponse[BillingProfileEnvelopeResponse](t, w)
	if got.Profile != nil {
		t.Fatalf("profile = %#v, want nil", got.Profile)
	}
}

func TestGetBillingProfileReturnsExistingProfile(t *testing.T) {
	router, db := testBillingRouter(t)
	user := createTestUser(t, db, "billing-profile-existing@example.com")
	fiscalCode := "RO2785503"
	if _, err := billing.UpsertBillingProfile(t.Context(), db, billing.UpsertBillingProfileInput{
		UserID:       user.ID,
		EntityType:   billing.BillingEntityCompany,
		BillingName:  "ICI Bucuresti",
		BillingEmail: "billing@example.com",
		CountryCode:  "RO",
		AddressLine1: "Maresal Averescu 8-10",
		City:         "Bucuresti",
		PostalCode:   "011455",
		FiscalCode:   &fiscalCode,
	}); err != nil {
		t.Fatalf("seed billing profile: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/billing/profile?user_id="+user.ID, nil)
	req.Header.Set(internalAPIHeader, testInternalAPIToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeBillingResponse[BillingProfileEnvelopeResponse](t, w)
	if got.Profile == nil {
		t.Fatal("profile = nil, want existing profile")
	}
	if string(got.Profile.UserID) != user.ID || got.Profile.EntityType != "company" || got.Profile.FiscalCode == nil || *got.Profile.FiscalCode != fiscalCode {
		t.Fatalf("profile = %#v", got.Profile)
	}
}

func TestPutBillingProfileUpsertsProfile(t *testing.T) {
	router, db := testBillingRouter(t)
	user := createTestUser(t, db, "billing-profile-upsert@example.com")

	body := `{"user_id":"` + user.ID + `","entity_type":"company","billing_name":"ICI Bucuresti","billing_email":"billing@example.com","country_code":"RO","address_line1":"Maresal Averescu 8-10","city":"Bucuresti","postal_code":"011455","fiscal_code":"RO2785503","registration_number":"J40/1234/1999"}`
	w := billingJSON(t, router, http.MethodPut, "/api/billing/profile", body, testInternalAPIToken)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeBillingResponse[BillingProfileResponse](t, w)
	if string(got.UserID) != user.ID || got.EntityType != "company" || got.FiscalCode == nil || *got.FiscalCode != "RO2785503" {
		t.Fatalf("profile response = %#v", got)
	}
}

func TestPutBillingProfileRequiresRomanianCompanyFiscalCode(t *testing.T) {
	router, db := testBillingRouter(t)
	user := createTestUser(t, db, "billing-profile-validation@example.com")

	body := `{"user_id":"` + user.ID + `","entity_type":"company","billing_name":"ICI Bucuresti","billing_email":"billing@example.com","country_code":"RO","address_line1":"Maresal Averescu 8-10","city":"Bucuresti","postal_code":"011455"}`
	w := billingJSON(t, router, http.MethodPut, "/api/billing/profile", body, testInternalAPIToken)

	if w.Code != http.StatusBadRequest || !strings.Contains(w.Body.String(), "fiscal code is required for Romanian companies") {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestPutBillingProfileReturnsValidationErrorForOverlongFields(t *testing.T) {
	router, db := testBillingRouter(t)
	user := createTestUser(t, db, "billing-profile-overlong@example.com")

	body := `{"user_id":"` + user.ID + `","entity_type":"individual","billing_name":"` + strings.Repeat("a", 256) + `","billing_email":"billing@example.com","country_code":"RO","address_line1":"Maresal Averescu 8-10","city":"Bucuresti","postal_code":"011455"}`
	w := billingJSON(t, router, http.MethodPut, "/api/billing/profile", body, testInternalAPIToken)

	if w.Code != http.StatusBadRequest || !strings.Contains(w.Body.String(), "billing name must be 255 characters or fewer") {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestCreateBillingOrderValidatesAndReturnsOrder(t *testing.T) {
	router, db := testBillingRouter(t)
	user := createTestUser(t, db, "billing-create-order@example.com")

	invalid := billingJSON(t, router, http.MethodPost, "/api/billing/orders", `{"user_id":"`+user.ID+`","credits":1500}`, testInternalAPIToken)
	if invalid.Code != http.StatusBadRequest {
		t.Fatalf("invalid status = %d body=%s", invalid.Code, invalid.Body.String())
	}

	w := billingJSON(t, router, http.MethodPost, "/api/billing/orders", `{"user_id":"`+user.ID+`","credits":5000}`, testInternalAPIToken)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeBillingResponse[BillingOrderResponse](t, w)
	if got.ID == uuid.Nil ||
		string(got.UserID) != user.ID ||
		got.Status != string(billing.OrderStatusPending) ||
		got.Provider != string(billing.BillingProviderStripe) ||
		got.PricingTier != string(billing.CreditPurchaseTier2) ||
		got.UnitAmountCents != 950 ||
		got.Credits != 5000 ||
		got.AmountCents != 4750 ||
		got.Currency != "EUR" {
		t.Fatalf("order response = %#v", got)
	}
}

func TestAttachBillingOrderCheckoutSessionStoresSessionID(t *testing.T) {
	router, db := testBillingRouter(t)
	user := createTestUser(t, db, "billing-checkout-session@example.com")
	order := createBillingTestOrder(t, db, user.ID, 1000)

	w := billingJSON(t, router, http.MethodPost, "/api/billing/orders/"+order.ID.String()+"/checkout-session", `{"checkout_session_id":"cs_test_attach"}`, testInternalAPIToken)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}

	var got billing.BillingOrder
	if err := db.First(&got, "id = ?", order.ID).Error; err != nil {
		t.Fatalf("load order: %v", err)
	}
	if got.ProviderCheckoutSessionID == nil || *got.ProviderCheckoutSessionID != "cs_test_attach" {
		t.Fatalf("checkout session = %#v, want cs_test_attach", got.ProviderCheckoutSessionID)
	}
}

func TestAttachBillingOrderCheckoutSessionRejectsSessionBelongingToAnotherPendingOrder(t *testing.T) {
	router, db := testBillingRouter(t)
	user := createTestUser(t, db, "billing-checkout-session-conflict@example.com")
	first := createBillingTestOrder(t, db, user.ID, 1000)
	second := createBillingTestOrder(t, db, user.ID, 1000)

	firstAttach := billingJSON(t, router, http.MethodPost, "/api/billing/orders/"+first.ID.String()+"/checkout-session", `{"checkout_session_id":"cs_test_pending_duplicate"}`, testInternalAPIToken)
	if firstAttach.Code != http.StatusNoContent {
		t.Fatalf("first attach status = %d body=%s", firstAttach.Code, firstAttach.Body.String())
	}

	conflict := billingJSON(t, router, http.MethodPost, "/api/billing/orders/"+second.ID.String()+"/checkout-session", `{"checkout_session_id":"cs_test_pending_duplicate"}`, testInternalAPIToken)
	if conflict.Code != http.StatusConflict {
		t.Fatalf("conflict status = %d body=%s", conflict.Code, conflict.Body.String())
	}

	var got billing.BillingOrder
	if err := db.First(&got, "id = ?", second.ID).Error; err != nil {
		t.Fatalf("load second order: %v", err)
	}
	if got.ProviderCheckoutSessionID != nil {
		t.Fatalf("second order checkout session = %#v, want nil", got.ProviderCheckoutSessionID)
	}
}

func TestMarkBillingOrderPaidGrantsCreditsOnce(t *testing.T) {
	router, db := testBillingRouter(t)
	user := createTestUser(t, db, "billing-paid@example.com")
	order := createBillingTestOrder(t, db, user.ID, 1000)
	body := `{"checkout_session_id":"cs_test_paid","payment_intent_id":"pi_test_paid","paid_at":"2026-06-04T12:00:00Z"}`

	for i := 0; i < 2; i++ {
		w := billingJSON(t, router, http.MethodPost, "/api/billing/orders/"+order.ID.String()+"/paid", body, testInternalAPIToken)
		if w.Code != http.StatusOK {
			t.Fatalf("attempt %d status = %d body=%s", i+1, w.Code, w.Body.String())
		}
		got := decodeBillingResponse[BillingOrderResponse](t, w)
		if got.Status != string(billing.OrderStatusPaid) || got.PaidAt == nil || got.ProviderPaymentIntentID == nil || *got.ProviderPaymentIntentID != "pi_test_paid" {
			t.Fatalf("attempt %d order response = %#v", i+1, got)
		}
	}

	balance, err := billing.AvailableCredits(t.Context(), db, user.ID, time.Date(2026, 6, 4, 12, 1, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("available credits: %v", err)
	}
	if balance.Available != 1000 {
		t.Fatalf("available = %d, want 1000", balance.Available)
	}
}

func TestMarkBillingOrderPaidCreatesInvoiceAndPDFOnce(t *testing.T) {
	var gotenbergRequests int32
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
		if r.FormValue("printBackground") != "true" ||
			r.FormValue("preferCssPageSize") != "true" ||
			!bytes.Contains(html, []byte("Invoice SYN-00001")) ||
			!bytes.Contains(html, []byte("Paid Customer")) {
			http.Error(w, "unexpected invoice render input", http.StatusBadRequest)
			return
		}
		atomic.AddInt32(&gotenbergRequests, 1)
		w.Header().Set("Content-Type", "application/pdf")
		_, _ = w.Write([]byte("%PDF-paid-order"))
	}))
	defer gotenberg.Close()

	outputDir := t.TempDir()
	router, db := testBillingRouterWithConfig(t, func(h *Handler) {
		h.GotenbergAPIURL = gotenberg.URL
		h.StorageDir = outputDir
		h.Now = func() time.Time {
			return time.Date(2026, 6, 10, 22, 30, 0, 0, time.UTC)
		}
	})
	user := createTestUser(t, db, "billing-paid-invoice@example.com")
	if _, err := billing.UpsertBillingProfile(t.Context(), db, billing.UpsertBillingProfileInput{
		UserID:       user.ID,
		EntityType:   billing.BillingEntityIndividual,
		BillingName:  "Paid Customer",
		BillingEmail: "paid@example.com",
		CountryCode:  "RO",
		AddressLine1: "Maresal Averescu 8-10",
		City:         "Bucuresti",
		PostalCode:   "011455",
	}); err != nil {
		t.Fatalf("seed billing profile: %v", err)
	}
	order := createBillingTestOrder(t, db, user.ID, 1000)
	body := `{"checkout_session_id":"cs_test_paid_invoice","payment_intent_id":"pi_test_paid_invoice","paid_at":"2026-06-04T12:00:00Z"}`

	first := billingJSON(t, router, http.MethodPost, "/api/billing/orders/"+order.ID.String()+"/paid", body, testInternalAPIToken)
	if first.Code != http.StatusOK {
		t.Fatalf("first status = %d body=%s", first.Code, first.Body.String())
	}
	firstOrder := decodeBillingResponse[BillingOrderResponse](t, first)
	if firstOrder.Invoice == nil ||
		firstOrder.Invoice.InvoiceSerie != "SYN" ||
		firstOrder.Invoice.InvoiceNumber != 1 ||
		firstOrder.Invoice.InvoiceDate != "2026-06-11" ||
		firstOrder.Invoice.PDFPath == nil {
		t.Fatalf("first paid order invoice = %#v", firstOrder.Invoice)
	}
	pdfPath := filepath.Join(outputDir, "invoices", firstOrder.Invoice.ID.String()+".pdf")
	if *firstOrder.Invoice.PDFPath != pdfPath {
		t.Fatalf("pdf_path = %q, want %q", *firstOrder.Invoice.PDFPath, pdfPath)
	}
	assertFileBytes(t, pdfPath, []byte("%PDF-paid-order"))

	second := billingJSON(t, router, http.MethodPost, "/api/billing/orders/"+order.ID.String()+"/paid", body, testInternalAPIToken)
	if second.Code != http.StatusOK {
		t.Fatalf("second status = %d body=%s", second.Code, second.Body.String())
	}
	secondOrder := decodeBillingResponse[BillingOrderResponse](t, second)
	if secondOrder.Invoice == nil || secondOrder.Invoice.ID != firstOrder.Invoice.ID || secondOrder.Invoice.PDFPath == nil || *secondOrder.Invoice.PDFPath != pdfPath {
		t.Fatalf("second paid order invoice = %#v, want existing invoice %s", secondOrder.Invoice, firstOrder.Invoice.ID)
	}

	var invoiceCount int64
	if err := db.Model(&billing.BillingInvoice{}).Where("order_id = ?", order.ID).Count(&invoiceCount).Error; err != nil {
		t.Fatalf("count invoices: %v", err)
	}
	if invoiceCount != 1 {
		t.Fatalf("invoice count = %d, want 1", invoiceCount)
	}
	var counter billing.BillingInvoiceCounter
	if err := db.First(&counter, "invoice_serie = ?", "SYN").Error; err != nil {
		t.Fatalf("load invoice counter: %v", err)
	}
	if counter.LastNumber != 1 {
		t.Fatalf("counter last number = %d, want 1", counter.LastNumber)
	}
	if got := atomic.LoadInt32(&gotenbergRequests); got != 1 {
		t.Fatalf("gotenberg requests = %d, want 1", got)
	}
	balance, err := billing.AvailableCredits(t.Context(), db, user.ID, time.Date(2026, 6, 4, 12, 1, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("available credits: %v", err)
	}
	if balance.Available != 1000 {
		t.Fatalf("available = %d, want 1000", balance.Available)
	}
}

func TestMarkBillingOrderPaidAcknowledgesPDFGenerationFailure(t *testing.T) {
	var gotenbergRequests int32
	gotenberg := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		atomic.AddInt32(&gotenbergRequests, 1)
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("not a pdf"))
	}))
	defer gotenberg.Close()

	outputDir := t.TempDir()
	router, db := testBillingRouterWithConfig(t, func(h *Handler) {
		h.GotenbergAPIURL = gotenberg.URL
		h.StorageDir = outputDir
		h.Now = func() time.Time {
			return time.Date(2026, 6, 11, 10, 0, 0, 0, time.UTC)
		}
	})
	user := createTestUser(t, db, "billing-paid-pdf-failure@example.com")
	order := createBillingTestOrder(t, db, user.ID, 1000)

	w := billingJSON(t, router, http.MethodPost, "/api/billing/orders/"+order.ID.String()+"/paid", `{"checkout_session_id":"cs_test_paid_pdf_failure","payment_intent_id":"pi_test_paid_pdf_failure","paid_at":"2026-06-04T12:00:00Z"}`, testInternalAPIToken)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeBillingResponse[BillingOrderResponse](t, w)
	if got.Status != string(billing.OrderStatusPaid) || got.Invoice == nil || got.Invoice.PDFPath != nil {
		t.Fatalf("paid order response = %#v", got)
	}
	var invoice billing.BillingInvoice
	if err := db.First(&invoice, "order_id = ?", order.ID).Error; err != nil {
		t.Fatalf("load invoice: %v", err)
	}
	if invoice.PDFPath != nil {
		t.Fatalf("stored PDFPath = %#v, want nil", invoice.PDFPath)
	}
	if _, err := os.Stat(filepath.Join(outputDir, "invoices", invoice.ID.String()+".pdf")); !os.IsNotExist(err) {
		t.Fatalf("pdf file stat err = %v, want not exist", err)
	}
	if got := atomic.LoadInt32(&gotenbergRequests); got != 1 {
		t.Fatalf("gotenberg requests = %d, want 1", got)
	}
	balance, err := billing.AvailableCredits(t.Context(), db, user.ID, time.Date(2026, 6, 4, 12, 1, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("available credits: %v", err)
	}
	if balance.Available != 1000 {
		t.Fatalf("available = %d, want 1000", balance.Available)
	}
}

func TestMarkBillingOrderPaidRejectsCheckoutSessionBelongingToAnotherPendingOrder(t *testing.T) {
	router, db := testBillingRouter(t)
	user := createTestUser(t, db, "billing-paid-pending-conflict@example.com")
	first := createBillingTestOrder(t, db, user.ID, 1000)
	second := createBillingTestOrder(t, db, user.ID, 1000)

	firstAttach := billingJSON(t, router, http.MethodPost, "/api/billing/orders/"+first.ID.String()+"/checkout-session", `{"checkout_session_id":"cs_test_pending_paid_duplicate"}`, testInternalAPIToken)
	if firstAttach.Code != http.StatusNoContent {
		t.Fatalf("first attach status = %d body=%s", firstAttach.Code, firstAttach.Body.String())
	}

	conflict := billingJSON(t, router, http.MethodPost, "/api/billing/orders/"+second.ID.String()+"/paid", `{"checkout_session_id":"cs_test_pending_paid_duplicate","paid_at":"2026-06-04T12:01:00Z"}`, testInternalAPIToken)
	if conflict.Code != http.StatusConflict {
		t.Fatalf("conflict status = %d body=%s", conflict.Code, conflict.Body.String())
	}

	var got billing.BillingOrder
	if err := db.First(&got, "id = ?", second.ID).Error; err != nil {
		t.Fatalf("load second order: %v", err)
	}
	if got.Status != billing.OrderStatusPending || got.ProviderCheckoutSessionID != nil {
		t.Fatalf("second order = %#v, want pending without checkout session", got)
	}
	balance, err := billing.AvailableCredits(t.Context(), db, user.ID, time.Date(2026, 6, 4, 12, 2, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("available credits: %v", err)
	}
	if balance.Available != 0 {
		t.Fatalf("available = %d, want 0", balance.Available)
	}
}

func TestMarkBillingOrderPaidRejectsProviderIDBelongingToAnotherOrder(t *testing.T) {
	router, db := testBillingRouter(t)
	user := createTestUser(t, db, "billing-paid-conflict@example.com")
	first := createBillingTestOrder(t, db, user.ID, 1000)
	second := createBillingTestOrder(t, db, user.ID, 1000)

	firstPaid := billingJSON(t, router, http.MethodPost, "/api/billing/orders/"+first.ID.String()+"/paid", `{"payment_intent_id":"pi_test_duplicate","paid_at":"2026-06-04T12:00:00Z"}`, testInternalAPIToken)
	if firstPaid.Code != http.StatusOK {
		t.Fatalf("first paid status = %d body=%s", firstPaid.Code, firstPaid.Body.String())
	}

	conflict := billingJSON(t, router, http.MethodPost, "/api/billing/orders/"+second.ID.String()+"/paid", `{"payment_intent_id":"pi_test_duplicate","paid_at":"2026-06-04T12:01:00Z"}`, testInternalAPIToken)
	if conflict.Code != http.StatusConflict {
		t.Fatalf("conflict status = %d body=%s", conflict.Code, conflict.Body.String())
	}

	var got billing.BillingOrder
	if err := db.First(&got, "id = ?", second.ID).Error; err != nil {
		t.Fatalf("load second order: %v", err)
	}
	if got.Status != billing.OrderStatusPending {
		t.Fatalf("second order status = %q, want pending", got.Status)
	}
	balance, err := billing.AvailableCredits(t.Context(), db, user.ID, time.Date(2026, 6, 4, 12, 2, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("available credits: %v", err)
	}
	if balance.Available != 1000 {
		t.Fatalf("available = %d, want only first order credits", balance.Available)
	}
}

func TestMarkBillingOrderFailedDoesNotGrantCredits(t *testing.T) {
	router, db := testBillingRouter(t)
	user := createTestUser(t, db, "billing-failed@example.com")
	order := createBillingTestOrder(t, db, user.ID, 1000)

	w := billingJSON(t, router, http.MethodPost, "/api/billing/orders/"+order.ID.String()+"/failed", `{}`, testInternalAPIToken)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeBillingResponse[BillingOrderResponse](t, w)
	if got.Status != string(billing.OrderStatusFailed) || got.FailedAt == nil {
		t.Fatalf("failed order response = %#v", got)
	}

	balance, err := billing.AvailableCredits(t.Context(), db, user.ID, time.Now().UTC())
	if err != nil {
		t.Fatalf("available credits: %v", err)
	}
	if balance.Available != 0 {
		t.Fatalf("available = %d, want 0", balance.Available)
	}
}

func TestListBillingOrdersFiltersScopesAndPaginates(t *testing.T) {
	router, db := testBillingRouter(t)
	user := createTestUser(t, db, "billing-orders-list@example.com")
	other := createTestUser(t, db, "billing-orders-list-other@example.com")
	now := time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC)
	paidAt := now.Add(5 * time.Minute)

	matching := createBillingOrderForList(t, db, user.ID, 5000, billing.OrderStatusPaid, now.Add(-20*time.Minute), &paidAt)
	createBillingOrderForList(t, db, user.ID, 1000, billing.OrderStatusPending, now.Add(-10*time.Minute), nil)
	createBillingOrderForList(t, db, user.ID, 10000, billing.OrderStatusPaid, now.Add(-2*time.Hour), &paidAt)
	createBillingOrderForList(t, db, other.ID, 5000, billing.OrderStatusPaid, now.Add(-15*time.Minute), &paidAt)

	path := "/api/billing/orders?user_id=" + user.ID +
		"&status=paid&created_from=2026-06-04T11:30:00Z&created_to=2026-06-04T12:00:00Z&size=1"
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.Header.Set(internalAPIHeader, testInternalAPIToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeBillingResponse[BillingOrderListResponse](t, w)
	if len(got.Orders) != 1 || got.Orders[0].ID != matching.ID {
		t.Fatalf("orders = %#v, want matching order %s", got.Orders, matching.ID)
	}
	if got.Orders[0].Status != "paid" ||
		got.Orders[0].AmountCents != matching.AmountCents ||
		got.Orders[0].Credits != matching.Credits ||
		got.Orders[0].PaidAt == nil {
		t.Fatalf("order response = %#v", got.Orders[0])
	}
	if got.NextCursor != nil {
		t.Fatalf("next_cursor = %#v, want nil for one matching row", got.NextCursor)
	}
}

func TestListBillingOrdersIncludesInvoiceMetadata(t *testing.T) {
	router, db := testBillingRouter(t)
	user := createTestUser(t, db, "billing-orders-invoice@example.com")
	now := time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC)
	paidAt := now.Add(5 * time.Minute)
	order := createBillingOrderForList(t, db, user.ID, 5000, billing.OrderStatusPaid, now, &paidAt)
	invoice, err := billing.CreateBillingInvoiceForPaidOrder(t.Context(), db, order.ID, time.Date(2026, 6, 11, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("create billing invoice for order: %v", err)
	}
	pdfPath := "/var/lib/syncra/invoices/" + invoice.ID.String() + ".pdf"
	if err := db.Model(&billing.BillingInvoice{}).
		Where("id = ?", invoice.ID).
		Update("pdf_path", pdfPath).Error; err != nil {
		t.Fatalf("update invoice pdf path: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/billing/orders?user_id="+user.ID, nil)
	req.Header.Set(internalAPIHeader, testInternalAPIToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeBillingResponse[struct {
		Orders []struct {
			ID      uuid.UUID `json:"id"`
			Invoice *struct {
				ID            uuid.UUID `json:"id"`
				InvoiceSerie  string    `json:"invoice_serie"`
				InvoiceNumber int64     `json:"invoice_number"`
				InvoiceDate   string    `json:"invoice_date"`
				PDFPath       *string   `json:"pdf_path,omitempty"`
			} `json:"invoice"`
		} `json:"orders"`
		NextCursor *string `json:"next_cursor"`
	}](t, w)
	if len(got.Orders) != 1 || got.Orders[0].ID != order.ID {
		t.Fatalf("orders = %#v, want order %s", got.Orders, order.ID)
	}
	if got.Orders[0].Invoice == nil ||
		got.Orders[0].Invoice.ID != invoice.ID ||
		got.Orders[0].Invoice.InvoiceSerie != "SYN" ||
		got.Orders[0].Invoice.InvoiceNumber != 1 ||
		got.Orders[0].Invoice.InvoiceDate != "2026-06-11" ||
		got.Orders[0].Invoice.PDFPath == nil ||
		*got.Orders[0].Invoice.PDFPath != pdfPath {
		t.Fatalf("invoice = %#v, want seeded invoice metadata", got.Orders[0].Invoice)
	}
}

func TestListBillingOrdersCursorPagination(t *testing.T) {
	router, db := testBillingRouter(t)
	user := createTestUser(t, db, "billing-orders-cursor@example.com")
	now := time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC)

	oldest := createBillingOrderForList(t, db, user.ID, 1000, billing.OrderStatusPending, now.Add(-30*time.Minute), nil)
	middle := createBillingOrderForList(t, db, user.ID, 5000, billing.OrderStatusPaid, now.Add(-20*time.Minute), timePtr(now.Add(-19*time.Minute)))
	newest := createBillingOrderForList(t, db, user.ID, 10000, billing.OrderStatusFailed, now.Add(-10*time.Minute), nil)

	req := httptest.NewRequest(http.MethodGet, "/api/billing/orders?user_id="+user.ID+"&size=2", nil)
	req.Header.Set(internalAPIHeader, testInternalAPIToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("first status = %d body=%s", w.Code, w.Body.String())
	}
	first := decodeBillingResponse[BillingOrderListResponse](t, w)
	if len(first.Orders) != 2 ||
		first.Orders[0].ID != newest.ID ||
		first.Orders[1].ID != middle.ID {
		t.Fatalf("first page = %#v, want newest,middle", first.Orders)
	}
	if first.NextCursor == nil {
		t.Fatal("first next_cursor = nil, want cursor")
	}

	req = httptest.NewRequest(http.MethodGet, "/api/billing/orders?user_id="+user.ID+"&size=2&cursor="+url.QueryEscape(*first.NextCursor), nil)
	req.Header.Set(internalAPIHeader, testInternalAPIToken)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("second status = %d body=%s", w.Code, w.Body.String())
	}
	second := decodeBillingResponse[BillingOrderListResponse](t, w)
	if len(second.Orders) != 1 || second.Orders[0].ID != oldest.ID {
		t.Fatalf("second page = %#v, want oldest", second.Orders)
	}
	if second.NextCursor != nil {
		t.Fatalf("second next_cursor = %#v, want nil", second.NextCursor)
	}
}

func TestListBillingOrdersRequiresInternalToken(t *testing.T) {
	router, db := testBillingRouter(t)
	user := createTestUser(t, db, "billing-orders-auth@example.com")
	for _, tt := range []struct {
		name  string
		token string
	}{
		{name: "missing token"},
		{name: "invalid token", token: "wrong-token"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/billing/orders?user_id="+user.ID, nil)
			if tt.token != "" {
				req.Header.Set(internalAPIHeader, tt.token)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			if w.Code != http.StatusUnauthorized {
				t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
			}
			if !strings.Contains(w.Body.String(), "unauthorized") {
				t.Fatalf("body = %s, want unauthorized", w.Body.String())
			}
		})
	}
}

func TestListBillingOrdersRepositoryFailureReturnsInternalServerError(t *testing.T) {
	router, db := testBillingRouter(t)
	if err := db.Exec("DROP TABLE billing_orders CASCADE").Error; err != nil {
		t.Fatalf("drop billing orders: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/billing/orders?user_id="+uuid.NewString(), nil)
	req.Header.Set(internalAPIHeader, testInternalAPIToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "failed to list billing orders") {
		t.Fatalf("body = %s, want failed to list billing orders", w.Body.String())
	}
}

func TestListBillingOrdersRejectsInvalidParameters(t *testing.T) {
	router, db := testBillingRouter(t)
	user := createTestUser(t, db, "billing-orders-invalid@example.com")
	mismatchCursor, err := encodeBillingOrderCursor(billing.BillingOrderCursor{
		CreatedAt: time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC),
		ID:        uuid.New(),
		Sort:      "desc",
	})
	if err != nil {
		t.Fatalf("encode mismatch cursor: %v", err)
	}
	for _, tt := range []struct {
		name      string
		query     string
		wantError string
	}{
		{name: "missing user", query: "", wantError: "user_id is required"},
		{name: "invalid status", query: "user_id=" + user.ID + "&status=settled", wantError: "invalid order status"},
		{name: "invalid created from", query: "user_id=" + user.ID + "&created_from=nope", wantError: "invalid created_from"},
		{name: "invalid created to", query: "user_id=" + user.ID + "&created_to=nope", wantError: "invalid created_to"},
		{name: "backwards date", query: "user_id=" + user.ID + "&created_from=2026-06-05T00:00:00Z&created_to=2026-06-04T00:00:00Z", wantError: "created_from must be before or equal to created_to"},
		{name: "bad base64 cursor", query: "user_id=" + user.ID + "&cursor=%25", wantError: "invalid cursor"},
		{name: "malformed cursor", query: "user_id=" + user.ID + "&cursor=e30", wantError: "invalid cursor"},
		{name: "cursor sort mismatch", query: "user_id=" + user.ID + "&sort=asc&cursor=" + url.QueryEscape(mismatchCursor), wantError: "cursor sort does not match sort"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/billing/orders?"+tt.query, nil)
			req.Header.Set(internalAPIHeader, testInternalAPIToken)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			if w.Code != http.StatusBadRequest {
				t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
			}
			if !strings.Contains(w.Body.String(), tt.wantError) {
				t.Fatalf("body = %s, want %q", w.Body.String(), tt.wantError)
			}
		})
	}
}

func TestListCreditUsageHistoryFiltersAndPaginates(t *testing.T) {
	router, db := testBillingRouter(t)
	user := createTestUser(t, db, "credit-usage-history@example.com")
	now := time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC)
	bucket := billing.CreditBucket{
		UserID: user.ID, SourceType: billing.CreditSourceAdjustment,
		CreditsGranted: 5000, CreditsRemaining: 5000, ValidFrom: now.Add(-time.Hour),
	}
	if err := db.Create(&bucket).Error; err != nil {
		t.Fatalf("create bucket: %v", err)
	}
	purchase := createBillingLedgerEntry(t, db, billing.CreditLedgerEntry{
		UserID: user.ID, BucketID: &bucket.ID, EntryType: billing.CreditLedgerEntryPurchase,
		CreditsDelta: 1000, IdempotencyKey: "api:purchase", CreatedAt: now.Add(-20 * time.Minute),
	})
	createBillingLedgerEntry(t, db, billing.CreditLedgerEntry{
		UserID: user.ID, BucketID: &bucket.ID, EntryType: billing.CreditLedgerEntryDebit,
		CreditsDelta: -3, IdempotencyKey: "api:debit", CreatedAt: now.Add(-10 * time.Minute),
	})

	path := "/api/billing/credit-usage-history?user_id=" + user.ID +
		"&type=purchase&created_from=2026-06-04T11:00:00Z&created_to=2026-06-04T12:00:00Z&size=1"
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.Header.Set(internalAPIHeader, testInternalAPIToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeBillingResponse[CreditUsageHistoryListResponse](t, w)
	if len(got.CreditUsageHistory) != 1 ||
		got.CreditUsageHistory[0].ID != purchase.ID ||
		got.CreditUsageHistory[0].EntryType != "purchase" {
		t.Fatalf("credit_usage_history = %#v, want purchase", got.CreditUsageHistory)
	}
	if got.CreditUsageHistory[0].CreditsDelta != 1000 || got.CreditUsageHistory[0].CreatedAt.IsZero() {
		t.Fatalf("credit usage history response = %#v", got.CreditUsageHistory[0])
	}
	if got.NextCursor != nil {
		t.Fatalf("next_cursor = %#v, want nil for one matching row", got.NextCursor)
	}
}

func TestListCreditUsageHistoryCursorPagination(t *testing.T) {
	router, db := testBillingRouter(t)
	user := createTestUser(t, db, "credit-usage-history-cursor@example.com")
	now := time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC)
	bucket := billing.CreditBucket{
		UserID: user.ID, SourceType: billing.CreditSourceAdjustment,
		CreditsGranted: 5000, CreditsRemaining: 5000, ValidFrom: now.Add(-time.Hour),
	}
	if err := db.Create(&bucket).Error; err != nil {
		t.Fatalf("create bucket: %v", err)
	}
	oldest := createBillingLedgerEntry(t, db, billing.CreditLedgerEntry{
		UserID: user.ID, BucketID: &bucket.ID, EntryType: billing.CreditLedgerEntryDebit,
		CreditsDelta: -1, IdempotencyKey: "api:cursor:oldest", CreatedAt: now.Add(-30 * time.Minute),
	})
	middle := createBillingLedgerEntry(t, db, billing.CreditLedgerEntry{
		UserID: user.ID, BucketID: &bucket.ID, EntryType: billing.CreditLedgerEntryDebit,
		CreditsDelta: -2, IdempotencyKey: "api:cursor:middle", CreatedAt: now.Add(-20 * time.Minute),
	})
	newest := createBillingLedgerEntry(t, db, billing.CreditLedgerEntry{
		UserID: user.ID, BucketID: &bucket.ID, EntryType: billing.CreditLedgerEntryDebit,
		CreditsDelta: -3, IdempotencyKey: "api:cursor:newest", CreatedAt: now.Add(-10 * time.Minute),
	})

	req := httptest.NewRequest(http.MethodGet, "/api/billing/credit-usage-history?user_id="+user.ID+"&size=2", nil)
	req.Header.Set(internalAPIHeader, testInternalAPIToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("first status = %d body=%s", w.Code, w.Body.String())
	}
	first := decodeBillingResponse[CreditUsageHistoryListResponse](t, w)
	if len(first.CreditUsageHistory) != 2 ||
		first.CreditUsageHistory[0].ID != newest.ID ||
		first.CreditUsageHistory[1].ID != middle.ID {
		t.Fatalf("first page = %#v, want newest,middle", first.CreditUsageHistory)
	}
	if first.NextCursor == nil {
		t.Fatal("first next_cursor = nil, want cursor")
	}

	req = httptest.NewRequest(http.MethodGet, "/api/billing/credit-usage-history?user_id="+user.ID+"&size=2&cursor="+url.QueryEscape(*first.NextCursor), nil)
	req.Header.Set(internalAPIHeader, testInternalAPIToken)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("second status = %d body=%s", w.Code, w.Body.String())
	}
	second := decodeBillingResponse[CreditUsageHistoryListResponse](t, w)
	if len(second.CreditUsageHistory) != 1 || second.CreditUsageHistory[0].ID != oldest.ID {
		t.Fatalf("second page = %#v, want oldest", second.CreditUsageHistory)
	}
	if second.NextCursor != nil {
		t.Fatalf("second next_cursor = %#v, want nil", second.NextCursor)
	}
}

func TestListCreditUsageHistoryRequiresInternalToken(t *testing.T) {
	router, db := testBillingRouter(t)
	user := createTestUser(t, db, "credit-usage-history-auth@example.com")
	for _, tt := range []struct {
		name  string
		token string
	}{
		{name: "missing token"},
		{name: "invalid token", token: "wrong-token"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/billing/credit-usage-history?user_id="+user.ID, nil)
			if tt.token != "" {
				req.Header.Set(internalAPIHeader, tt.token)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			if w.Code != http.StatusUnauthorized {
				t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
			}
			if !strings.Contains(w.Body.String(), "unauthorized") {
				t.Fatalf("body = %s, want unauthorized", w.Body.String())
			}
		})
	}
}

func TestListCreditUsageHistoryRepositoryFailureReturnsInternalServerError(t *testing.T) {
	router, db := testBillingRouter(t)
	if err := db.Exec("DROP TABLE credit_ledger_entries").Error; err != nil {
		t.Fatalf("drop credit ledger entries: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/billing/credit-usage-history?user_id="+uuid.NewString(), nil)
	req.Header.Set(internalAPIHeader, testInternalAPIToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "failed to list credit usage history") {
		t.Fatalf("body = %s, want failed to list credit usage history", w.Body.String())
	}
}

func TestListCreditUsageHistoryRejectsInvalidParameters(t *testing.T) {
	router, db := testBillingRouter(t)
	user := createTestUser(t, db, "credit-usage-history-invalid@example.com")
	mismatchCursor, err := encodeCreditUsageHistoryCursor(billing.CreditLedgerTransactionCursor{
		CreatedAt: time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC),
		ID:        uuid.New(),
		Sort:      "desc",
	})
	if err != nil {
		t.Fatalf("encode mismatch cursor: %v", err)
	}
	for _, tt := range []struct {
		name      string
		query     string
		wantError string
	}{
		{name: "missing user", query: "", wantError: "user_id is required"},
		{name: "invalid type", query: "user_id=" + user.ID + "&type=grant", wantError: "entry type must be purchase or debit"},
		{name: "invalid created from", query: "user_id=" + user.ID + "&created_from=nope", wantError: "invalid created_from"},
		{name: "invalid created to", query: "user_id=" + user.ID + "&created_to=nope", wantError: "invalid created_to"},
		{name: "backwards date", query: "user_id=" + user.ID + "&created_from=2026-06-05T00:00:00Z&created_to=2026-06-04T00:00:00Z", wantError: "created_from must be before or equal to created_to"},
		{name: "bad base64 cursor", query: "user_id=" + user.ID + "&cursor=%25", wantError: "invalid cursor"},
		{name: "malformed cursor", query: "user_id=" + user.ID + "&cursor=e30", wantError: "invalid cursor"},
		{name: "cursor sort mismatch", query: "user_id=" + user.ID + "&sort=asc&cursor=" + url.QueryEscape(mismatchCursor), wantError: "cursor sort does not match sort"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/billing/credit-usage-history?"+tt.query, nil)
			req.Header.Set(internalAPIHeader, testInternalAPIToken)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			if w.Code != http.StatusBadRequest {
				t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
			}
			if !strings.Contains(w.Body.String(), tt.wantError) {
				t.Fatalf("body = %s, want %q", w.Body.String(), tt.wantError)
			}
		})
	}
}

func createBillingTestOrder(t *testing.T, db *gorm.DB, userID string, credits int) billing.BillingOrder {
	t.Helper()
	order, err := billing.CreateCreditOrder(t.Context(), db, billing.CreateCreditOrderInput{
		UserID:  userID,
		Credits: credits,
	})
	if err != nil {
		t.Fatalf("create billing order: %v", err)
	}
	return order
}

func createBillingOrderForList(t *testing.T, db *gorm.DB, userID string, credits int, status billing.OrderStatus, createdAt time.Time, paidAt *time.Time) billing.BillingOrder {
	t.Helper()
	order := createBillingTestOrder(t, db, userID, credits)
	updates := map[string]any{
		"status":     status,
		"created_at": createdAt.UTC(),
	}
	if paidAt != nil {
		updates["paid_at"] = paidAt.UTC()
	}
	if err := db.Model(&billing.BillingOrder{}).Where("id = ?", order.ID).Updates(updates).Error; err != nil {
		t.Fatalf("update billing order for list: %v", err)
	}
	var out billing.BillingOrder
	if err := db.First(&out, "id = ?", order.ID).Error; err != nil {
		t.Fatalf("reload billing order for list: %v", err)
	}
	return out
}

func createBillingLedgerEntry(t *testing.T, db *gorm.DB, entry billing.CreditLedgerEntry) billing.CreditLedgerEntry {
	t.Helper()
	if err := db.Create(&entry).Error; err != nil {
		t.Fatalf("create billing ledger entry: %v", err)
	}
	return entry
}

func timePtr(value time.Time) *time.Time {
	return &value
}

func billingJSON(t *testing.T, router http.Handler, method string, path string, body string, internalToken string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	if internalToken != "" {
		req.Header.Set(internalAPIHeader, internalToken)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func decodeBillingResponse[T any](t *testing.T, w *httptest.ResponseRecorder) T {
	t.Helper()
	var out T
	if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode billing response: %v body=%s", err, w.Body.String())
	}
	return out
}
