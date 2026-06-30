package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"ai.ro/syncra/internal/billing"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type adminBillingOrderUserTestResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type adminBillingOrderTestResponse struct {
	ID          string                                `json:"id"`
	UserID      string                                `json:"user_id"`
	User        adminBillingOrderUserTestResponse     `json:"user"`
	Invoice     *adminBillingOrderInvoiceTestResponse `json:"invoice"`
	Status      string                                `json:"status"`
	AmountCents int                                   `json:"amount_cents"`
	Credits     int                                   `json:"credits"`
	PaidAt      *time.Time                            `json:"paid_at"`
}

type adminBillingOrderInvoiceTestResponse struct {
	ID            string `json:"id"`
	InvoiceSerie  string `json:"invoice_serie"`
	InvoiceNumber int64  `json:"invoice_number"`
	InvoiceDate   string `json:"invoice_date"`
}

type adminBillingOrderListTestResponse struct {
	Orders     []adminBillingOrderTestResponse `json:"orders"`
	NextCursor *string                         `json:"next_cursor"`
}

type adminBillingInvoiceListTestResponse struct {
	Invoices   []BillingInvoiceResponse `json:"invoices"`
	NextCursor *string                  `json:"next_cursor"`
}

func TestListAdminBillingOrdersRequiresAdminSession(t *testing.T) {
	router, db := testAdminRouter(t)
	normal := createAdminTestUser(t, db, "normal-admin-orders@example.com", "user")
	normalCookie := createAdminTestSession(t, db, normal, "normal-admin-orders-session")

	w := adminJSON(t, router, http.MethodGet, "/api/admin/billing/orders", "", nil)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("missing session status = %d body=%s", w.Code, w.Body.String())
	}

	w = adminJSON(t, router, http.MethodGet, "/api/admin/billing/orders", "", normalCookie)
	if w.Code != http.StatusForbidden {
		t.Fatalf("non-admin status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestListAdminBillingInvoicesRequiresAdminSession(t *testing.T) {
	router, db := testAdminRouter(t)
	normal := createAdminTestUser(t, db, "normal-admin-invoices@example.com", "user")
	normalCookie := createAdminTestSession(t, db, normal, "normal-admin-invoices-session")

	w := adminJSON(t, router, http.MethodGet, "/api/admin/billing/invoices", "", nil)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("missing session status = %d body=%s", w.Code, w.Body.String())
	}

	w = adminJSON(t, router, http.MethodGet, "/api/admin/billing/invoices", "", normalCookie)
	if w.Code != http.StatusForbidden {
		t.Fatalf("non-admin status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestListAdminBillingOrdersListsAllUsersAndFiltersByUser(t *testing.T) {
	router, db := testAdminRouter(t)
	admin := createAdminTestUser(t, db, "admin-orders@example.com", "admin")
	cookie := createAdminTestSession(t, db, admin, "admin-orders-session")
	user := createAdminTestUser(t, db, "target-orders@example.com", "user")
	other := createAdminTestUser(t, db, "other-orders@example.com", "user")
	now := time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC)
	paidAt := now.Add(5 * time.Minute)

	older := createBillingOrderForList(t, db, user.ID, 5000, billing.OrderStatusPaid, now.Add(-20*time.Minute), &paidAt)
	newer := createBillingOrderForList(t, db, other.ID, 1000, billing.OrderStatusPending, now.Add(-10*time.Minute), nil)
	if _, err := billing.UpsertBillingProfile(t.Context(), db, billing.UpsertBillingProfileInput{
		UserID:       user.ID,
		EntityType:   billing.BillingEntityIndividual,
		BillingName:  "Target User",
		BillingEmail: "target-orders@example.com",
		CountryCode:  "RO",
		AddressLine1: "Maresal Averescu 8-10",
		City:         "Bucuresti",
		PostalCode:   "011455",
	}); err != nil {
		t.Fatalf("seed billing profile: %v", err)
	}
	invoice, err := billing.CreateBillingInvoiceForPaidOrder(t.Context(), db, older.ID, time.Date(2026, 6, 5, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("seed billing invoice: %v", err)
	}

	w := adminJSON(t, router, http.MethodGet, "/api/admin/billing/orders?size=20", "", cookie)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAdminResponse[adminBillingOrderListTestResponse](t, w)
	if len(got.Orders) != 2 || got.Orders[0].ID != newer.ID.String() || got.Orders[1].ID != older.ID.String() {
		t.Fatalf("orders = %#v, want newer then older", got.Orders)
	}
	if got.Orders[0].User.ID != other.ID ||
		got.Orders[0].User.Name != other.Name ||
		got.Orders[0].User.Email != other.Email {
		t.Fatalf("newer order user = %#v, want %s/%s", got.Orders[0].User, other.Name, other.Email)
	}
	if got.Orders[1].User.ID != user.ID ||
		got.Orders[1].User.Name != user.Name ||
		got.Orders[1].User.Email != user.Email ||
		got.Orders[1].PaidAt == nil {
		t.Fatalf("older order response = %#v, want user metadata and paid_at", got.Orders[1])
	}
	if got.Orders[0].Invoice != nil {
		t.Fatalf("newer invoice = %#v, want nil", got.Orders[0].Invoice)
	}
	if got.Orders[1].Invoice == nil ||
		got.Orders[1].Invoice.ID != invoice.ID.String() ||
		got.Orders[1].Invoice.InvoiceSerie != "SYN" ||
		got.Orders[1].Invoice.InvoiceNumber != 1 ||
		got.Orders[1].Invoice.InvoiceDate != "2026-06-05" {
		t.Fatalf("older invoice = %#v, want seeded invoice metadata", got.Orders[1].Invoice)
	}

	w = adminJSON(t, router, http.MethodGet, "/api/admin/billing/orders?user_id="+user.ID+"&size=20", "", cookie)
	if w.Code != http.StatusOK {
		t.Fatalf("filtered status = %d body=%s", w.Code, w.Body.String())
	}
	got = decodeAdminResponse[adminBillingOrderListTestResponse](t, w)
	if len(got.Orders) != 1 || got.Orders[0].ID != older.ID.String() {
		t.Fatalf("filtered orders = %#v, want older only", got.Orders)
	}
}

func TestListAdminBillingOrdersFiltersPaidOrdersWithoutInvoices(t *testing.T) {
	router, db := testAdminRouter(t)
	admin := createAdminTestUser(t, db, "admin-orders-without-invoices@example.com", "admin")
	cookie := createAdminTestSession(t, db, admin, "admin-orders-without-invoices-session")
	user := createAdminTestUser(t, db, "orders-without-invoices@example.com", "user")
	now := time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC)

	paidNoInvoiceAt := now.Add(-9 * time.Minute)
	paidNoInvoice := createBillingOrderForList(t, db, user.ID, 10000, billing.OrderStatusPaid, now.Add(-10*time.Minute), &paidNoInvoiceAt)
	paidWithInvoiceAt := now.Add(-19 * time.Minute)
	paidWithInvoice := createBillingOrderForList(t, db, user.ID, 5000, billing.OrderStatusPaid, now.Add(-20*time.Minute), &paidWithInvoiceAt)
	createBillingOrderForList(t, db, user.ID, 1000, billing.OrderStatusPending, now.Add(-5*time.Minute), nil)
	if _, err := billing.CreateBillingInvoiceForPaidOrder(t.Context(), db, paidWithInvoice.ID, time.Date(2026, 6, 5, 0, 0, 0, 0, time.UTC)); err != nil {
		t.Fatalf("seed billing invoice: %v", err)
	}

	w := adminJSON(t, router, http.MethodGet, "/api/admin/billing/orders?status=paid&without_invoice=true&size=20", "", cookie)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAdminResponse[adminBillingOrderListTestResponse](t, w)
	if len(got.Orders) != 1 || got.Orders[0].ID != paidNoInvoice.ID.String() {
		t.Fatalf("orders = %#v, want paid order without invoice only", got.Orders)
	}
	if got.Orders[0].Invoice != nil || got.Orders[0].Status != string(billing.OrderStatusPaid) {
		t.Fatalf("order response = %#v, want paid order with nil invoice", got.Orders[0])
	}
}

func TestListAdminBillingInvoicesListsAllUsersAndFilters(t *testing.T) {
	router, db := testAdminRouter(t)
	admin := createAdminTestUser(t, db, "admin-invoices@example.com", "admin")
	cookie := createAdminTestSession(t, db, admin, "admin-invoices-session")
	user := createAdminTestUser(t, db, "target-invoices@example.com", "user")
	other := createAdminTestUser(t, db, "other-invoices@example.com", "user")
	now := time.Date(2026, 6, 11, 12, 0, 0, 0, time.UTC)

	older := createAdminBillingInvoiceForList(t, db, user.ID, "SYN-ADMIN", now.Add(-30*time.Minute))
	middle := createAdminBillingInvoiceForList(t, db, other.ID, "SYN-ADMIN", now.Add(-20*time.Minute))
	newer := createAdminBillingInvoiceForList(t, db, user.ID, "SYN-ADMIN", now.Add(-10*time.Minute))

	w := adminJSON(t, router, http.MethodGet, "/api/admin/billing/invoices?size=2", "", cookie)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAdminResponse[adminBillingInvoiceListTestResponse](t, w)
	if len(got.Invoices) != 2 ||
		got.Invoices[0].ID != newer.ID ||
		got.Invoices[1].ID != middle.ID ||
		got.NextCursor == nil {
		t.Fatalf("first invoices page = %#v cursor=%v, want newer,middle and cursor", got.Invoices, got.NextCursor)
	}
	if got.Invoices[0].BillingName != user.Name ||
		got.Invoices[0].BillingEmail != user.Email ||
		got.Invoices[0].NetAmount != "10.00" ||
		got.Invoices[0].VATAmount != "1.90" ||
		got.Invoices[0].TotalAmount != "11.90" ||
		got.Invoices[0].InvoiceNumber != 3 {
		t.Fatalf("newer invoice response = %#v, want billing fields and totals", got.Invoices[0])
	}

	w = adminJSON(t, router, http.MethodGet, "/api/admin/billing/invoices?size=2&cursor="+*got.NextCursor, "", cookie)
	if w.Code != http.StatusOK {
		t.Fatalf("second status = %d body=%s", w.Code, w.Body.String())
	}
	got = decodeAdminResponse[adminBillingInvoiceListTestResponse](t, w)
	if len(got.Invoices) != 1 || got.Invoices[0].ID != older.ID || got.NextCursor != nil {
		t.Fatalf("second invoices page = %#v cursor=%v, want older and no cursor", got.Invoices, got.NextCursor)
	}

	createdFrom := now.Add(-25 * time.Minute).Format(time.RFC3339Nano)
	createdTo := now.Add(-15 * time.Minute).Format(time.RFC3339Nano)
	path := "/api/admin/billing/invoices?user_id=" + other.ID + "&created_from=" + createdFrom + "&created_to=" + createdTo + "&size=20"
	w = adminJSON(t, router, http.MethodGet, path, "", cookie)
	if w.Code != http.StatusOK {
		t.Fatalf("filtered status = %d body=%s", w.Code, w.Body.String())
	}
	got = decodeAdminResponse[adminBillingInvoiceListTestResponse](t, w)
	if len(got.Invoices) != 1 || got.Invoices[0].ID != middle.ID {
		t.Fatalf("filtered invoices = %#v, want middle only", got.Invoices)
	}
}

func TestListAdminBillingInvoicesSearchesClientAndInvoiceNumber(t *testing.T) {
	router, db := testAdminRouter(t)
	admin := createAdminTestUser(t, db, "admin-invoices-search@example.com", "admin")
	cookie := createAdminTestSession(t, db, admin, "admin-invoices-search-session")
	user := createAdminTestUser(t, db, "target-invoices-search@example.com", "user")
	now := time.Date(2026, 6, 11, 12, 0, 0, 0, time.UTC)

	ada := updateAdminBillingInvoiceClient(
		t,
		db,
		createAdminBillingInvoiceForList(t, db, user.ID, "SYN-ADMIN-SEARCH", now.Add(-30*time.Minute)).ID,
		"Ada Lovelace",
		"ada@example.com",
	)
	grace := updateAdminBillingInvoiceClient(
		t,
		db,
		createAdminBillingInvoiceForList(t, db, user.ID, "SYN-ADMIN-SEARCH", now.Add(-20*time.Minute)).ID,
		"Grace Hopper",
		"grace@example.com",
	)
	katherine := updateAdminBillingInvoiceClient(
		t,
		db,
		createAdminBillingInvoiceForList(t, db, user.ID, "SYN-ADMIN-SEARCH", now.Add(-10*time.Minute)).ID,
		"Katherine Johnson",
		"katherine@example.com",
	)

	tests := []struct {
		name   string
		search string
		want   uuid.UUID
	}{
		{name: "billing name", search: "lovelace", want: ada.ID},
		{name: "billing email", search: "GRACE@EXAMPLE.COM", want: grace.ID},
		{
			name:   "formatted invoice label",
			search: fmt.Sprintf("%s-%05d", ada.InvoiceSerie, ada.InvoiceNumber),
			want:   ada.ID,
		},
		{
			name:   "unpadded invoice label",
			search: fmt.Sprintf("%s-%d", katherine.InvoiceSerie, katherine.InvoiceNumber),
			want:   katherine.ID,
		},
		{name: "raw invoice number fragment", search: "2", want: grace.ID},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := "/api/admin/billing/invoices?search=" + url.QueryEscape(tt.search) + "&size=20"
			w := adminJSON(t, router, http.MethodGet, path, "", cookie)
			if w.Code != http.StatusOK {
				t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
			}
			got := decodeAdminResponse[adminBillingInvoiceListTestResponse](t, w)
			if len(got.Invoices) != 1 || got.Invoices[0].ID != tt.want {
				t.Fatalf("search %q invoices = %#v, want %s", tt.search, got.Invoices, tt.want)
			}
		})
	}
}

func TestListAdminBillingInvoicesRejectsInvalidUserID(t *testing.T) {
	router, db := testAdminRouter(t)
	admin := createAdminTestUser(t, db, "admin-invoices-invalid@example.com", "admin")
	cookie := createAdminTestSession(t, db, admin, "admin-invoices-invalid-session")

	w := adminJSON(t, router, http.MethodGet, "/api/admin/billing/invoices?user_id=nope", "", cookie)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "invalid user_id") {
		t.Fatalf("body = %s, want invalid user_id", w.Body.String())
	}
}

func TestListAdminBillingOrdersRejectsInvalidUserID(t *testing.T) {
	router, db := testAdminRouter(t)
	admin := createAdminTestUser(t, db, "admin-orders-invalid@example.com", "admin")
	cookie := createAdminTestSession(t, db, admin, "admin-orders-invalid-session")

	w := adminJSON(t, router, http.MethodGet, "/api/admin/billing/orders?user_id=nope", "", cookie)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "invalid user_id") {
		t.Fatalf("body = %s, want invalid user_id", w.Body.String())
	}
}

func TestListAdminBillingOrdersRejectsInvalidWithoutInvoice(t *testing.T) {
	router, db := testAdminRouter(t)
	admin := createAdminTestUser(t, db, "admin-orders-invalid-without-invoice@example.com", "admin")
	cookie := createAdminTestSession(t, db, admin, "admin-orders-invalid-without-invoice-session")

	w := adminJSON(t, router, http.MethodGet, "/api/admin/billing/orders?without_invoice=nope", "", cookie)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "invalid without_invoice") {
		t.Fatalf("body = %s, want invalid without_invoice", w.Body.String())
	}
}

func TestCreateAdminBillingOrderInvoiceCreatesInvoiceWithBucharestDate(t *testing.T) {
	router, db := testAdminRouterWithNow(t, func() time.Time {
		return time.Date(2026, 6, 10, 21, 30, 0, 0, time.UTC)
	})
	admin := createAdminTestUser(t, db, "admin-invoice@example.com", "admin")
	cookie := createAdminTestSession(t, db, admin, "admin-invoice-session")
	user := createAdminTestUser(t, db, "invoice-target@example.com", "user")
	seedAdminBillingProfile(t, db, user.ID)
	now := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)
	order := createBillingOrderForList(t, db, user.ID, 5000, billing.OrderStatusPaid, now, timePtr(now.Add(time.Minute)))

	w := adminJSON(t, router, http.MethodPost, "/api/admin/billing/orders/"+order.ID.String()+"/invoice", "", cookie)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAdminResponse[BillingInvoiceResponse](t, w)
	if got.OrderID == nil ||
		*got.OrderID != order.ID ||
		got.InvoiceSerie != "SYN" ||
		got.InvoiceNumber != 1 ||
		got.InvoiceDate != "2026-06-11" ||
		got.NetAmount != "47.50" ||
		got.VATAmount != "0.00" ||
		got.TotalAmount != "47.50" {
		t.Fatalf("invoice response = %#v", got)
	}
	if len(got.Lines) != 1 ||
		got.Lines[0].Name != "SYNCRA SaaS 5000 credits" ||
		got.Lines[0].Quantity != 1 ||
		got.Lines[0].UnitPrice != "47.50" ||
		got.Lines[0].VATPercentage != "0.00" ||
		got.Lines[0].TotalVATAmount != "0.00" ||
		got.Lines[0].TotalAmount != "47.50" {
		t.Fatalf("invoice lines = %#v", got.Lines)
	}
}

func TestCreateAdminBillingOrderInvoiceRequiresAdminSession(t *testing.T) {
	router, db := testAdminRouter(t)
	normal := createAdminTestUser(t, db, "normal-invoice@example.com", "user")
	normalCookie := createAdminTestSession(t, db, normal, "normal-invoice-session")
	orderID := uuid.New()

	w := adminJSON(t, router, http.MethodPost, "/api/admin/billing/orders/"+orderID.String()+"/invoice", "", nil)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("missing session status = %d body=%s", w.Code, w.Body.String())
	}

	w = adminJSON(t, router, http.MethodPost, "/api/admin/billing/orders/"+orderID.String()+"/invoice", "", normalCookie)
	if w.Code != http.StatusForbidden {
		t.Fatalf("non-admin status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestCreateAdminBillingOrderInvoiceHandlesDuplicateUnpaidMissingAndNoProfile(t *testing.T) {
	router, db := testAdminRouter(t)
	admin := createAdminTestUser(t, db, "admin-invoice-errors@example.com", "admin")
	cookie := createAdminTestSession(t, db, admin, "admin-invoice-errors-session")
	user := createAdminTestUser(t, db, "invoice-errors@example.com", "user")
	seedAdminBillingProfile(t, db, user.ID)
	now := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)
	paid := createBillingOrderForList(t, db, user.ID, 1000, billing.OrderStatusPaid, now, timePtr(now.Add(time.Minute)))
	pending := createBillingOrderForList(t, db, user.ID, 1000, billing.OrderStatusPending, now.Add(time.Minute), nil)
	withoutProfile := createAdminTestUser(t, db, "invoice-no-profile@example.com", "user")
	noProfileOrder := createBillingOrderForList(t, db, withoutProfile.ID, 1000, billing.OrderStatusPaid, now.Add(2*time.Minute), timePtr(now.Add(3*time.Minute)))

	first := adminJSON(t, router, http.MethodPost, "/api/admin/billing/orders/"+paid.ID.String()+"/invoice", "", cookie)
	if first.Code != http.StatusCreated {
		t.Fatalf("first status = %d body=%s", first.Code, first.Body.String())
	}

	duplicate := adminJSON(t, router, http.MethodPost, "/api/admin/billing/orders/"+paid.ID.String()+"/invoice", "", cookie)
	if duplicate.Code != http.StatusConflict || !strings.Contains(duplicate.Body.String(), "already has an invoice") {
		t.Fatalf("duplicate status = %d body=%s", duplicate.Code, duplicate.Body.String())
	}

	unpaid := adminJSON(t, router, http.MethodPost, "/api/admin/billing/orders/"+pending.ID.String()+"/invoice", "", cookie)
	if unpaid.Code != http.StatusBadRequest || !strings.Contains(unpaid.Body.String(), "not paid") {
		t.Fatalf("unpaid status = %d body=%s", unpaid.Code, unpaid.Body.String())
	}

	missing := adminJSON(t, router, http.MethodPost, "/api/admin/billing/orders/"+uuid.NewString()+"/invoice", "", cookie)
	if missing.Code != http.StatusNotFound || !strings.Contains(missing.Body.String(), "not found") {
		t.Fatalf("missing status = %d body=%s", missing.Code, missing.Body.String())
	}

	noProfile := adminJSON(t, router, http.MethodPost, "/api/admin/billing/orders/"+noProfileOrder.ID.String()+"/invoice", "", cookie)
	if noProfile.Code != http.StatusCreated {
		t.Fatalf("no profile status = %d body=%s", noProfile.Code, noProfile.Body.String())
	}
	got := decodeAdminResponse[BillingInvoiceResponse](t, noProfile)
	if got.UserID == nil ||
		string(*got.UserID) != withoutProfile.ID ||
		got.OrderID == nil ||
		*got.OrderID != noProfileOrder.ID ||
		got.BillingProfileID != nil ||
		got.BillingName != withoutProfile.Name ||
		got.BillingEmail != withoutProfile.Email ||
		got.BillingFiscalCode != nil {
		t.Fatalf("no profile invoice response = %#v, want user fallback buyer fields", got)
	}
	var snapshot struct {
		Source       string `json:"source"`
		UserID       string `json:"user_id"`
		BillingName  string `json:"billing_name"`
		BillingEmail string `json:"billing_email"`
	}
	if err := json.Unmarshal(got.BillingProfileSnapshot, &snapshot); err != nil {
		t.Fatalf("unmarshal no profile snapshot: %v", err)
	}
	if snapshot.Source != "user" ||
		snapshot.UserID != withoutProfile.ID ||
		snapshot.BillingName != withoutProfile.Name ||
		snapshot.BillingEmail != withoutProfile.Email {
		t.Fatalf("no profile snapshot = %#v, want user buyer data", snapshot)
	}
}

func seedAdminBillingProfile(t *testing.T, db *gorm.DB, userID string) {
	t.Helper()
	if _, err := billing.UpsertBillingProfile(t.Context(), db, billing.UpsertBillingProfileInput{
		UserID:       userID,
		EntityType:   billing.BillingEntityIndividual,
		BillingName:  "Invoice Target",
		BillingEmail: "invoice-target@example.com",
		CountryCode:  "RO",
		AddressLine1: "Maresal Averescu 8-10",
		City:         "Bucuresti",
		PostalCode:   "011455",
	}); err != nil {
		t.Fatalf("seed billing profile: %v", err)
	}
}

func createAdminBillingInvoiceForList(t *testing.T, db *gorm.DB, userID string, invoiceSerie string, createdAt time.Time) billing.BillingInvoice {
	t.Helper()
	invoice, err := billing.CreateBillingInvoice(t.Context(), db, billing.CreateBillingInvoiceInput{
		UserID:       userID,
		InvoiceSerie: invoiceSerie,
		InvoiceDate:  time.Date(2026, 6, 11, 0, 0, 0, 0, time.UTC),
		Lines: []billing.CreateBillingInvoiceLineInput{
			{Name: "OCR credits", Quantity: 1, UnitPrice: "10.00", VATPercentage: "19.00"},
		},
	})
	if err != nil {
		t.Fatalf("create billing invoice: %v", err)
	}
	if err := db.Model(&billing.BillingInvoice{}).Where("id = ?", invoice.ID).Update("created_at", createdAt.UTC()).Error; err != nil {
		t.Fatalf("update billing invoice created_at: %v", err)
	}
	var out billing.BillingInvoice
	if err := db.First(&out, "id = ?", invoice.ID).Error; err != nil {
		t.Fatalf("reload billing invoice: %v", err)
	}
	return out
}

func updateAdminBillingInvoiceClient(t *testing.T, db *gorm.DB, invoiceID uuid.UUID, billingName string, billingEmail string) billing.BillingInvoice {
	t.Helper()
	if err := db.Model(&billing.BillingInvoice{}).Where("id = ?", invoiceID).Updates(map[string]any{
		"billing_name":  billingName,
		"billing_email": billingEmail,
	}).Error; err != nil {
		t.Fatalf("update admin billing invoice client: %v", err)
	}
	var out billing.BillingInvoice
	if err := db.First(&out, "id = ?", invoiceID).Error; err != nil {
		t.Fatalf("reload admin billing invoice client: %v", err)
	}
	return out
}
