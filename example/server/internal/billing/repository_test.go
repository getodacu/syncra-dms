package billing

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"ai.ro/syncra/internal/testsupport"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var billingRepositoryTestGroup *testsupport.PostgresGroup

func TestBillingRepository(t *testing.T) {
	t.Run("GrantCreditsUsesAdvisoryIdempotencyLock", testGrantCreditsUsesAdvisoryIdempotencyLock)

	t.Run("transactional", func(t *testing.T) {
		billingRepositoryTestGroup = testsupport.OpenPostgresGroup(t, billingTestModels()...)
		defer func() { billingRepositoryTestGroup = nil }()

		for _, tt := range []struct {
			name string
			fn   func(*testing.T)
		}{
			{name: "BillingRepositoryValidation", fn: testBillingRepositoryValidation},
			{name: "UpsertBillingProfileCreatesAndUpdatesProfile", fn: testUpsertBillingProfileCreatesAndUpdatesProfile},
			{name: "CreateBillingInvoiceSnapshotsProfileAndComputesTotals", fn: testCreateBillingInvoiceSnapshotsProfileAndComputesTotals},
			{name: "CreateBillingInvoiceUsesUserFallbackWithoutBillingProfile", fn: testCreateBillingInvoiceUsesUserFallbackWithoutBillingProfile},
			{name: "CreateBillingInvoiceMaintainsSeparateSeriesCounters", fn: testCreateBillingInvoiceMaintainsSeparateSeriesCounters},
			{name: "CreateBillingInvoiceForPaidOrderGeneratesSYNInvoice", fn: testCreateBillingInvoiceForPaidOrderGeneratesSYNInvoice},
			{name: "CreateBillingInvoiceForPaidOrderRejectsDuplicateWithoutAllocatingNumber", fn: testCreateBillingInvoiceForPaidOrderRejectsDuplicateWithoutAllocatingNumber},
			{name: "CreateBillingInvoiceForPaidOrderRejectsUnpaidOrders", fn: testCreateBillingInvoiceForPaidOrderRejectsUnpaidOrders},
			{name: "CreateBillingInvoiceForPaidOrderUsesUserFallbackWithoutBillingProfile", fn: testCreateBillingInvoiceForPaidOrderUsesUserFallbackWithoutBillingProfile},
			{name: "CreateBillingInvoiceValidation", fn: testCreateBillingInvoiceValidation},
			{name: "GetBillingProfile", fn: testGetBillingProfile},
			{name: "UpsertBillingProfileValidation", fn: testUpsertBillingProfileValidation},
			{name: "GrantSignupBonusCreatesNonExpiringBucket", fn: testGrantSignupBonusCreatesNonExpiringBucket},
			{name: "GrantSignupBonusIsIdempotent", fn: testGrantSignupBonusIsIdempotent},
			{name: "CreateCreditOrderUsesPurchaseQuote", fn: testCreateCreditOrderUsesPurchaseQuote},
			{name: "CreateCreditOrderRejectsInvalidCredits", fn: testCreateCreditOrderRejectsInvalidCredits},
			{name: "AttachCreditOrderCheckoutSession", fn: testAttachCreditOrderCheckoutSession},
			{name: "AttachCreditOrderCheckoutSessionRejectsConflictingSession", fn: testAttachCreditOrderCheckoutSessionRejectsConflictingSession},
			{name: "MarkCreditOrderPaidAndGrantCreditsCreatesNonExpiringBucket", fn: testMarkCreditOrderPaidAndGrantCreditsCreatesNonExpiringBucket},
			{name: "MarkCreditOrderPaidAndGrantCreditsIdempotent", fn: testMarkCreditOrderPaidAndGrantCreditsIdempotent},
			{name: "ProviderPaymentIntentCannotGrantTwoOrders", fn: testProviderPaymentIntentCannotGrantTwoOrders},
			{name: "MarkCreditOrderFailedDoesNotGrantCredits", fn: testMarkCreditOrderFailedDoesNotGrantCredits},
			{name: "MarkCreditOrderCanceledDoesNotGrantCredits", fn: testMarkCreditOrderCanceledDoesNotGrantCredits},
			{name: "AvailableCreditsIgnoresExpiredAndVoidedBuckets", fn: testAvailableCreditsIgnoresExpiredAndVoidedBuckets},
			{name: "AvailableCreditsIgnoresFutureBuckets", fn: testAvailableCreditsIgnoresFutureBuckets},
			{name: "AdjustCreditsAddsAdjustmentBucketAndLedgerEntry", fn: testAdjustCreditsAddsAdjustmentBucketAndLedgerEntry},
			{name: "AdjustCreditsSubtractsFromAvailableBuckets", fn: testAdjustCreditsSubtractsFromAvailableBuckets},
			{name: "AdjustCreditsRejectsInsufficientBalance", fn: testAdjustCreditsRejectsInsufficientBalance},
			{name: "ListCreditLedgerTransactionsFiltersAndSorts", fn: testListCreditLedgerTransactionsFiltersAndSorts},
			{name: "ListCreditLedgerTransactionsCursorPaginates", fn: testListCreditLedgerTransactionsCursorPaginates},
			{name: "ListCreditLedgerTransactionsValidation", fn: testListCreditLedgerTransactionsValidation},
			{name: "ListBillingOrdersFiltersAndSorts", fn: testListBillingOrdersFiltersAndSorts},
			{name: "ListBillingOrdersCursorPaginates", fn: testListBillingOrdersCursorPaginates},
			{name: "ListAdminBillingOrdersListsAllUsersAndFiltersByUser", fn: testListAdminBillingOrdersListsAllUsersAndFiltersByUser},
			{name: "ListAdminBillingOrdersFiltersPaidOrdersWithoutInvoices", fn: testListAdminBillingOrdersFiltersPaidOrdersWithoutInvoices},
			{name: "ListAdminBillingInvoicesListsAllUsersAndFilters", fn: testListAdminBillingInvoicesListsAllUsersAndFilters},
			{name: "ListAdminBillingInvoicesSearchesClientAndInvoiceNumber", fn: testListAdminBillingInvoicesSearchesClientAndInvoiceNumber},
			{name: "ListAdminBillingInvoicesPaginatesByCreatedAt", fn: testListAdminBillingInvoicesPaginatesByCreatedAt},
			{name: "ListBillingOrdersValidation", fn: testListBillingOrdersValidation},
			{name: "ListAdminBillingInvoicesValidation", fn: testListAdminBillingInvoicesValidation},
			{name: "DebitCreditsForJobConsumesBucketsInPriorityOrder", fn: testDebitCreditsForJobConsumesBucketsInPriorityOrder},
			{name: "DebitCreditsForJobTxUsesExistingTransaction", fn: testDebitCreditsForJobTxUsesExistingTransaction},
			{name: "DebitCreditsForJobInsufficientCreditsDoesNotMutate", fn: testDebitCreditsForJobInsufficientCreditsDoesNotMutate},
			{name: "DebitCreditsForJobRejectsNonPositiveCredits", fn: testDebitCreditsForJobRejectsNonPositiveCredits},
			{name: "DebitCreditsForJobIgnoresExpiredAndVoidedBuckets", fn: testDebitCreditsForJobIgnoresExpiredAndVoidedBuckets},
			{name: "DebitCreditsForJobIgnoresFutureBuckets", fn: testDebitCreditsForJobIgnoresFutureBuckets},
			{name: "DebitCreditsForJobIsIdempotent", fn: testDebitCreditsForJobIsIdempotent},
			{name: "DebitCreditsForJobIdempotencyKeyTreatsWildcardsLiterally", fn: testDebitCreditsForJobIdempotencyKeyTreatsWildcardsLiterally},
			{name: "RefundCreditsForJobRestoresActiveBucket", fn: testRefundCreditsForJobRestoresActiveBucket},
			{name: "RefundCreditsForJobCreatesRefundBucketForExpiredCredit", fn: testRefundCreditsForJobCreatesRefundBucketForExpiredCredit},
			{name: "RefundCreditsForJobUsesLaterOriginalExpiry", fn: testRefundCreditsForJobUsesLaterOriginalExpiry},
			{name: "RefundCreditsForJobCreatesRefundBucketForVoidedCredit", fn: testRefundCreditsForJobCreatesRefundBucketForVoidedCredit},
			{name: "RefundCreditsForJobCreatesNonExpiringRefundBucketForPurchasedCredit", fn: testRefundCreditsForJobCreatesNonExpiringRefundBucketForPurchasedCredit},
		} {
			t.Run(tt.name, tt.fn)
		}
	})

	t.Run("CreateBillingInvoiceConcurrentCreatesIncrementByOne", testCreateBillingInvoiceConcurrentCreatesIncrementByOne)
	t.Run("UpsertBillingProfileConcurrentFirstUpserts", testUpsertBillingProfileConcurrentFirstUpserts)
}

func billingRepositoryTx(t *testing.T) *gorm.DB {
	t.Helper()
	if billingRepositoryTestGroup != nil {
		return billingRepositoryTestGroup.Tx(t)
	}
	return testsupport.OpenPostgresTx(t, billingTestModels()...)
}

func testBillingRepositoryValidation(t *testing.T) {
	if _, err := CreateCreditOrder(context.Background(), nil, CreateCreditOrderInput{UserID: uuid.NewString(), Credits: 1000}); err == nil {
		t.Fatal("CreateCreditOrder nil db succeeded")
	}
	if _, err := CreateBillingInvoice(context.Background(), nil, CreateBillingInvoiceInput{UserID: uuid.NewString(), InvoiceSerie: "SYNC"}); err == nil {
		t.Fatal("CreateBillingInvoice nil db succeeded")
	}
	if _, err := CreateBillingInvoiceForPaidOrder(context.Background(), nil, uuid.New(), time.Now()); err == nil {
		t.Fatal("CreateBillingInvoiceForPaidOrder nil db succeeded")
	}

	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	if _, err := CreateCreditOrder(context.Background(), db, CreateCreditOrderInput{UserID: user.ID, Credits: 0}); err == nil {
		t.Fatal("CreateCreditOrder invalid credits succeeded")
	}
	if err := AttachCreditOrderCheckoutSession(context.Background(), db, uuid.Nil, "cs_test"); err == nil {
		t.Fatal("AttachCreditOrderCheckoutSession nil order id succeeded")
	}
	if err := AttachCreditOrderCheckoutSession(context.Background(), db, uuid.New(), ""); err == nil {
		t.Fatal("AttachCreditOrderCheckoutSession empty checkout session succeeded")
	}
	if _, err := MarkCreditOrderPaidAndGrantCredits(context.Background(), nil, MarkCreditOrderPaidInput{OrderID: uuid.New(), PaidAt: time.Now()}); err == nil {
		t.Fatal("MarkCreditOrderPaidAndGrantCredits nil db succeeded")
	}
	if _, err := MarkCreditOrderPaidAndGrantCredits(context.Background(), db, MarkCreditOrderPaidInput{OrderID: uuid.Nil, PaidAt: time.Now()}); err == nil {
		t.Fatal("MarkCreditOrderPaidAndGrantCredits nil order id succeeded")
	}
	if _, err := MarkCreditOrderPaidAndGrantCredits(context.Background(), db, MarkCreditOrderPaidInput{OrderID: uuid.New()}); err == nil {
		t.Fatal("MarkCreditOrderPaidAndGrantCredits zero paid time succeeded")
	}
	if _, err := CreateBillingInvoiceForPaidOrder(context.Background(), db, uuid.Nil, time.Now()); err == nil {
		t.Fatal("CreateBillingInvoiceForPaidOrder nil order id succeeded")
	}
	if err := MarkCreditOrderFailed(context.Background(), nil, uuid.New(), time.Now()); err == nil {
		t.Fatal("MarkCreditOrderFailed nil db succeeded")
	}
	if err := MarkCreditOrderCanceled(context.Background(), db, uuid.New(), time.Time{}); err == nil {
		t.Fatal("MarkCreditOrderCanceled zero canceled time succeeded")
	}
}

func testUpsertBillingProfileCreatesAndUpdatesProfile(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	input := validBillingProfileInput(user.ID)
	input.CountryCode = "ro"
	input.Region = stringPtr(" Bucuresti ")
	input.AddressLine2 = stringPtr(" ")

	created, err := UpsertBillingProfile(context.Background(), db, input)
	if err != nil {
		t.Fatalf("create UpsertBillingProfile() error = %v", err)
	}
	if created.ID == uuid.Nil {
		t.Fatal("created billing profile has nil ID")
	}
	if created.CountryCode != "RO" {
		t.Fatalf("created country code = %q, want RO", created.CountryCode)
	}
	if created.AddressLine2 != nil {
		t.Fatalf("created address line 2 = %#v, want nil for blank optional input", created.AddressLine2)
	}
	if created.Region == nil || *created.Region != "Bucuresti" {
		t.Fatalf("created region = %#v, want trimmed Bucuresti", created.Region)
	}

	fiscalCode := " RO2785503 "
	registrationNumber := " J40/1234/2026 "
	update := validBillingProfileInput(user.ID)
	update.EntityType = BillingEntityCompany
	update.BillingName = " ICI Bucuresti "
	update.FiscalCode = &fiscalCode
	update.RegistrationNumber = &registrationNumber

	updated, err := UpsertBillingProfile(context.Background(), db, update)
	if err != nil {
		t.Fatalf("update UpsertBillingProfile() error = %v", err)
	}
	if updated.ID != created.ID {
		t.Fatalf("updated profile ID = %s, want %s", updated.ID, created.ID)
	}
	if updated.EntityType != BillingEntityCompany ||
		updated.BillingName != "ICI Bucuresti" ||
		updated.FiscalCode == nil ||
		*updated.FiscalCode != "RO2785503" ||
		updated.RegistrationNumber == nil ||
		*updated.RegistrationNumber != "J40/1234/2026" {
		t.Fatalf("unexpected updated profile: %#v", updated)
	}
	assertBillingCount(t, db, &BillingProfile{}, 1, "user_id = ?", user.ID)
}

func testCreateBillingInvoiceSnapshotsProfileAndComputesTotals(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	fiscalCode := "RO2785503"
	registrationNumber := "J40/1234/2026"
	profile, err := UpsertBillingProfile(context.Background(), db, UpsertBillingProfileInput{
		UserID:             user.ID,
		EntityType:         BillingEntityCompany,
		BillingName:        "ICI Bucuresti",
		BillingEmail:       "billing@example.com",
		CountryCode:        "RO",
		AddressLine1:       "Maresal Averescu 8-10",
		City:               "Bucuresti",
		PostalCode:         "011455",
		FiscalCode:         &fiscalCode,
		RegistrationNumber: &registrationNumber,
	})
	if err != nil {
		t.Fatalf("create billing profile: %v", err)
	}

	invoiceDate := time.Date(2026, 6, 10, 14, 30, 0, 0, time.FixedZone("UTC+2", 2*60*60))
	invoice, err := CreateBillingInvoice(context.Background(), db, CreateBillingInvoiceInput{
		UserID:       " " + user.ID + " ",
		InvoiceSerie: " sync-2026 ",
		InvoiceDate:  invoiceDate,
		Lines: []CreateBillingInvoiceLineInput{
			{Name: " OCR credits ", Quantity: 2, UnitPrice: "10.00", VATPercentage: "19"},
			{Name: "Setup", Quantity: 1, UnitPrice: "5.55", VATPercentage: "0.00"},
		},
	})
	if err != nil {
		t.Fatalf("CreateBillingInvoice() error = %v", err)
	}

	if invoice.ID == uuid.Nil ||
		invoice.UserID == nil ||
		*invoice.UserID != user.ID ||
		invoice.BillingProfileID == nil ||
		*invoice.BillingProfileID != profile.ID ||
		invoice.InvoiceSerie != "SYNC-2026" ||
		invoice.InvoiceNumber != 1 {
		t.Fatalf("unexpected invoice identity fields: %#v", invoice)
	}
	if invoice.BillingName != "ICI Bucuresti" ||
		invoice.BillingEmail != "billing@example.com" ||
		invoice.BillingFiscalCode == nil ||
		*invoice.BillingFiscalCode != fiscalCode {
		t.Fatalf("unexpected searchable billing fields: %#v", invoice)
	}
	if invoice.NetAmount.StringFixed(2) != "25.55" ||
		invoice.VATAmount.StringFixed(2) != "3.80" ||
		invoice.TotalAmount.StringFixed(2) != "29.35" {
		t.Fatalf("invoice amounts = net %s vat %s total %s",
			invoice.NetAmount, invoice.VATAmount, invoice.TotalAmount)
	}
	if !invoice.InvoiceDate.Equal(time.Date(2026, 6, 10, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("invoice date = %s", invoice.InvoiceDate)
	}

	var snapshot BillingProfileSnapshot
	if err := json.Unmarshal(invoice.BillingProfileSnapshot, &snapshot); err != nil {
		t.Fatalf("unmarshal billing profile snapshot: %v", err)
	}
	if snapshot.ID != profile.ID ||
		snapshot.BillingName != "ICI Bucuresti" ||
		snapshot.FiscalCode == nil ||
		*snapshot.FiscalCode != fiscalCode {
		t.Fatalf("snapshot = %#v, want original profile", snapshot)
	}

	var lines []BillingInvoiceLine
	if err := json.Unmarshal(invoice.Lines, &lines); err != nil {
		t.Fatalf("unmarshal invoice lines: %v", err)
	}
	if len(lines) != 2 ||
		lines[0].Name != "OCR credits" ||
		lines[0].Quantity != 2 ||
		lines[0].UnitPrice != "10.00" ||
		lines[0].VATPercentage != "19.00" ||
		lines[0].TotalVATAmount != "3.80" ||
		lines[0].TotalAmount != "23.80" ||
		lines[1].TotalAmount != "5.55" {
		t.Fatalf("lines = %#v", lines)
	}

	updatedName := "Changed Billing Name"
	updatedProfile := validBillingProfileInput(user.ID)
	updatedProfile.BillingName = updatedName
	if _, err := UpsertBillingProfile(context.Background(), db, updatedProfile); err != nil {
		t.Fatalf("update billing profile: %v", err)
	}
	var reloaded BillingInvoice
	if err := db.First(&reloaded, "id = ?", invoice.ID).Error; err != nil {
		t.Fatalf("reload invoice: %v", err)
	}
	if reloaded.BillingName != "ICI Bucuresti" || strings.Contains(string(reloaded.BillingProfileSnapshot), updatedName) {
		t.Fatalf("invoice snapshot changed after profile update: %#v snapshot=%s", reloaded, reloaded.BillingProfileSnapshot)
	}
}

func testCreateBillingInvoiceUsesUserFallbackWithoutBillingProfile(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)

	invoice, err := CreateBillingInvoice(context.Background(), db, validInvoiceInput(user.ID, "SYNC"))
	if err != nil {
		t.Fatalf("CreateBillingInvoice() error = %v", err)
	}

	if invoice.UserID == nil ||
		*invoice.UserID != user.ID ||
		invoice.BillingProfileID != nil ||
		invoice.BillingName != user.Name ||
		invoice.BillingEmail != user.Email ||
		invoice.BillingFiscalCode != nil {
		t.Fatalf("unexpected fallback buyer fields: %#v", invoice)
	}
	var snapshot billingUserSnapshot
	if err := json.Unmarshal(invoice.BillingProfileSnapshot, &snapshot); err != nil {
		t.Fatalf("unmarshal fallback snapshot: %v", err)
	}
	if snapshot.Source != "user" ||
		snapshot.UserID != user.ID ||
		snapshot.BillingName != user.Name ||
		snapshot.BillingEmail != user.Email {
		t.Fatalf("fallback snapshot = %#v, want user buyer data", snapshot)
	}
}

func testCreateBillingInvoiceMaintainsSeparateSeriesCounters(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	if _, err := UpsertBillingProfile(context.Background(), db, validBillingProfileInput(user.ID)); err != nil {
		t.Fatalf("create billing profile: %v", err)
	}

	first, err := CreateBillingInvoice(context.Background(), db, validInvoiceInput(user.ID, "SYNC"))
	if err != nil {
		t.Fatalf("create first invoice: %v", err)
	}
	second, err := CreateBillingInvoice(context.Background(), db, validInvoiceInput(user.ID, "SYNC"))
	if err != nil {
		t.Fatalf("create second invoice: %v", err)
	}
	other, err := CreateBillingInvoice(context.Background(), db, validInvoiceInput(user.ID, "ALT"))
	if err != nil {
		t.Fatalf("create other series invoice: %v", err)
	}

	if first.InvoiceNumber != 1 || second.InvoiceNumber != 2 || other.InvoiceNumber != 1 {
		t.Fatalf("invoice numbers = first %d second %d other %d", first.InvoiceNumber, second.InvoiceNumber, other.InvoiceNumber)
	}
}

func testCreateBillingInvoiceForPaidOrderGeneratesSYNInvoice(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	if _, err := UpsertBillingProfile(context.Background(), db, validBillingProfileInput(user.ID)); err != nil {
		t.Fatalf("UpsertBillingProfile() error = %v", err)
	}
	now := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)
	order := createBillingOrderForList(t, db, user.ID, 5000, OrderStatusPaid, now, timePtr(now.Add(time.Minute)))
	invoiceDate := time.Date(2026, 6, 11, 0, 0, 0, 0, time.UTC)

	invoice, err := CreateBillingInvoiceForPaidOrder(context.Background(), db, order.ID, invoiceDate)
	if err != nil {
		t.Fatalf("CreateBillingInvoiceForPaidOrder() error = %v", err)
	}
	if invoice.OrderID == nil || *invoice.OrderID != order.ID {
		t.Fatalf("invoice order id = %#v, want %s", invoice.OrderID, order.ID)
	}
	if invoice.InvoiceSerie != "SYN" || invoice.InvoiceNumber != 1 || !invoice.InvoiceDate.Equal(invoiceDate) {
		t.Fatalf("invoice reference/date = %s-%d %s", invoice.InvoiceSerie, invoice.InvoiceNumber, invoice.InvoiceDate)
	}
	if invoice.NetAmount.StringFixed(2) != "47.50" ||
		invoice.VATAmount.StringFixed(2) != "0.00" ||
		invoice.TotalAmount.StringFixed(2) != "47.50" {
		t.Fatalf("invoice amounts = net %s vat %s total %s", invoice.NetAmount, invoice.VATAmount, invoice.TotalAmount)
	}
	var lines []BillingInvoiceLine
	if err := json.Unmarshal(invoice.Lines, &lines); err != nil {
		t.Fatalf("unmarshal lines: %v", err)
	}
	if len(lines) != 1 ||
		lines[0].Name != "SYNCRA SaaS 5000 credits" ||
		lines[0].Quantity != 1 ||
		lines[0].UnitPrice != "47.50" ||
		lines[0].VATPercentage != "0.00" ||
		lines[0].TotalVATAmount != "0.00" ||
		lines[0].TotalAmount != "47.50" {
		t.Fatalf("invoice lines = %#v", lines)
	}
}

func testCreateBillingInvoiceForPaidOrderRejectsDuplicateWithoutAllocatingNumber(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	if _, err := UpsertBillingProfile(context.Background(), db, validBillingProfileInput(user.ID)); err != nil {
		t.Fatalf("UpsertBillingProfile() error = %v", err)
	}
	now := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)
	order := createBillingOrderForList(t, db, user.ID, 1000, OrderStatusPaid, now, timePtr(now.Add(time.Minute)))
	invoiceDate := time.Date(2026, 6, 11, 0, 0, 0, 0, time.UTC)

	if _, err := CreateBillingInvoiceForPaidOrder(context.Background(), db, order.ID, invoiceDate); err != nil {
		t.Fatalf("first CreateBillingInvoiceForPaidOrder() error = %v", err)
	}
	if _, err := CreateBillingInvoiceForPaidOrder(context.Background(), db, order.ID, invoiceDate); !errors.Is(err, ErrBillingInvoiceExists) {
		t.Fatalf("duplicate CreateBillingInvoiceForPaidOrder() error = %v, want ErrBillingInvoiceExists", err)
	}
	assertBillingCount(t, db, &BillingInvoice{}, 1, "order_id = ?", order.ID)
	var counter BillingInvoiceCounter
	if err := db.First(&counter, "invoice_serie = ?", "SYN").Error; err != nil {
		t.Fatalf("load invoice counter: %v", err)
	}
	if counter.LastNumber != 1 {
		t.Fatalf("counter last number = %d, want 1", counter.LastNumber)
	}
}

func testCreateBillingInvoiceForPaidOrderRejectsUnpaidOrders(t *testing.T) {
	for _, status := range []OrderStatus{
		OrderStatusPending,
		OrderStatusFailed,
		OrderStatusRefunded,
		OrderStatusCanceled,
	} {
		t.Run(string(status), func(t *testing.T) {
			db := billingRepositoryTx(t)
			user := createBillingTestUser(t, db)
			if _, err := UpsertBillingProfile(context.Background(), db, validBillingProfileInput(user.ID)); err != nil {
				t.Fatalf("UpsertBillingProfile() error = %v", err)
			}
			now := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)
			order := createBillingOrderForList(t, db, user.ID, 1000, status, now, nil)

			_, err := CreateBillingInvoiceForPaidOrder(context.Background(), db, order.ID, now)
			if !errors.Is(err, ErrBillingOrderNotPaid) {
				t.Fatalf("CreateBillingInvoiceForPaidOrder() error = %v, want ErrBillingOrderNotPaid", err)
			}
		})
	}
}

func testCreateBillingInvoiceForPaidOrderUsesUserFallbackWithoutBillingProfile(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	now := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)
	order := createBillingOrderForList(t, db, user.ID, 1000, OrderStatusPaid, now, timePtr(now.Add(time.Minute)))

	invoice, err := CreateBillingInvoiceForPaidOrder(context.Background(), db, order.ID, now)
	if err != nil {
		t.Fatalf("CreateBillingInvoiceForPaidOrder() error = %v", err)
	}
	if invoice.OrderID == nil ||
		*invoice.OrderID != order.ID ||
		invoice.UserID == nil ||
		*invoice.UserID != user.ID ||
		invoice.BillingProfileID != nil ||
		invoice.BillingName != user.Name ||
		invoice.BillingEmail != user.Email ||
		invoice.BillingFiscalCode != nil {
		t.Fatalf("unexpected fallback paid-order invoice: %#v", invoice)
	}
	var snapshot billingUserSnapshot
	if err := json.Unmarshal(invoice.BillingProfileSnapshot, &snapshot); err != nil {
		t.Fatalf("unmarshal fallback snapshot: %v", err)
	}
	if snapshot.Source != "user" ||
		snapshot.UserID != user.ID ||
		snapshot.BillingName != user.Name ||
		snapshot.BillingEmail != user.Email {
		t.Fatalf("fallback snapshot = %#v, want user buyer data", snapshot)
	}
}

func testCreateBillingInvoiceConcurrentCreatesIncrementByOne(t *testing.T) {
	db := testsupport.OpenPostgresDB(t, billingTestModels()...)
	user := createBillingTestUser(t, db)
	if _, err := UpsertBillingProfile(context.Background(), db, validBillingProfileInput(user.ID)); err != nil {
		t.Fatalf("create billing profile: %v", err)
	}

	const workers = 12
	start := make(chan struct{})
	results := make(chan struct {
		number int64
		err    error
	}, workers)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			invoice, err := CreateBillingInvoice(context.Background(), db, validInvoiceInput(user.ID, "SYNC"))
			if err != nil {
				results <- struct {
					number int64
					err    error
				}{err: err}
				return
			}
			results <- struct {
				number int64
				err    error
			}{number: invoice.InvoiceNumber}
		}()
	}
	close(start)
	wg.Wait()
	close(results)

	seen := make(map[int64]bool, workers)
	for result := range results {
		if result.err != nil {
			t.Fatalf("concurrent CreateBillingInvoice() error = %v", result.err)
		}
		if result.number < 1 || result.number > workers {
			t.Fatalf("invoice number = %d, want 1..%d", result.number, workers)
		}
		if seen[result.number] {
			t.Fatalf("duplicate invoice number %d", result.number)
		}
		seen[result.number] = true
	}
	if len(seen) != workers {
		t.Fatalf("created %d invoice numbers, want %d", len(seen), workers)
	}
}

func testCreateBillingInvoiceValidation(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	if _, err := UpsertBillingProfile(context.Background(), db, validBillingProfileInput(user.ID)); err != nil {
		t.Fatalf("create billing profile: %v", err)
	}

	for _, tt := range []struct {
		name    string
		mutate  func(*CreateBillingInvoiceInput)
		wantErr string
	}{
		{
			name: "blank serie",
			mutate: func(input *CreateBillingInvoiceInput) {
				input.InvoiceSerie = " "
			},
			wantErr: "invoice serie is required",
		},
		{
			name: "invalid serie",
			mutate: func(input *CreateBillingInvoiceInput) {
				input.InvoiceSerie = "SYNC/2026"
			},
			wantErr: "invoice serie contains invalid characters",
		},
		{
			name: "no lines",
			mutate: func(input *CreateBillingInvoiceInput) {
				input.Lines = nil
			},
			wantErr: "invoice lines are required",
		},
		{
			name: "blank line name",
			mutate: func(input *CreateBillingInvoiceInput) {
				input.Lines[0].Name = " "
			},
			wantErr: "line name is required",
		},
		{
			name: "invalid quantity",
			mutate: func(input *CreateBillingInvoiceInput) {
				input.Lines[0].Quantity = 0
			},
			wantErr: "quantity must be positive",
		},
		{
			name: "invalid decimal",
			mutate: func(input *CreateBillingInvoiceInput) {
				input.Lines[0].UnitPrice = "ten"
			},
			wantErr: "unit price is invalid",
		},
		{
			name: "negative vat",
			mutate: func(input *CreateBillingInvoiceInput) {
				input.Lines[0].VATPercentage = "-1"
			},
			wantErr: "vat percentage must be non-negative",
		},
		{
			name: "vat too high",
			mutate: func(input *CreateBillingInvoiceInput) {
				input.Lines[0].VATPercentage = "101"
			},
			wantErr: "vat percentage must be between 0 and 100",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			input := validInvoiceInput(user.ID, "SYNC")
			tt.mutate(&input)
			_, err := CreateBillingInvoice(context.Background(), db, input)
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("CreateBillingInvoice() error = %v, want %q", err, tt.wantErr)
			}
		})
	}
}

func testUpsertBillingProfileConcurrentFirstUpserts(t *testing.T) {
	setupDB := testsupport.OpenPostgresDB(t, billingTestModels()...)
	user := createBillingTestUser(t, setupDB)
	firstDB := openBillingRepositoryTestDB(t)
	secondDB := openBillingRepositoryTestDB(t)
	ready := make(chan struct{}, 2)
	release := make(chan struct{})
	registerBillingProfileCreateBarrier(t, firstDB, ready, release)
	registerBillingProfileCreateBarrier(t, secondDB, ready, release)

	firstInput := validBillingProfileInput(user.ID)
	firstInput.BillingName = "First Billing Name"
	secondInput := validBillingProfileInput(user.ID)
	secondInput.BillingName = "Second Billing Name"

	type upsertResult struct {
		profile BillingProfile
		err     error
	}
	start := make(chan struct{})
	results := make(chan upsertResult, 2)
	var wg sync.WaitGroup
	for _, run := range []struct {
		db    *gorm.DB
		input UpsertBillingProfileInput
	}{
		{db: firstDB, input: firstInput},
		{db: secondDB, input: secondInput},
	} {
		wg.Add(1)
		go func(db *gorm.DB, input UpsertBillingProfileInput) {
			defer wg.Done()
			<-start
			profile, err := UpsertBillingProfile(context.Background(), db, input)
			results <- upsertResult{profile: profile, err: err}
		}(run.db, run.input)
	}

	close(start)
	released := false
	releaseBarrier := func() {
		if !released {
			close(release)
			released = true
		}
	}
	defer releaseBarrier()
	for i := 0; i < 2; i++ {
		select {
		case <-ready:
		case result := <-results:
			t.Fatalf("UpsertBillingProfile returned before insert barrier: profile=%#v err=%v", result.profile, result.err)
		case <-time.After(5 * time.Second):
			t.Fatal("timed out waiting for concurrent billing profile creates")
		}
	}
	releaseBarrier()
	wg.Wait()
	close(results)

	for result := range results {
		if result.err != nil {
			t.Fatalf("concurrent UpsertBillingProfile() error = %v", result.err)
		}
		if result.profile.ID == uuid.Nil || result.profile.UserID != user.ID {
			t.Fatalf("concurrent UpsertBillingProfile() returned invalid profile: %#v", result.profile)
		}
		if result.profile.BillingName != "First Billing Name" && result.profile.BillingName != "Second Billing Name" {
			t.Fatalf("concurrent UpsertBillingProfile() returned unexpected billing name: %#v", result.profile)
		}
	}
	assertBillingCount(t, setupDB, &BillingProfile{}, 1, "user_id = ?", user.ID)
}

func testGetBillingProfile(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)

	if _, err := GetBillingProfile(context.Background(), db, user.ID); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("missing GetBillingProfile() error = %v, want gorm.ErrRecordNotFound", err)
	}

	created, err := UpsertBillingProfile(context.Background(), db, validBillingProfileInput(user.ID))
	if err != nil {
		t.Fatalf("UpsertBillingProfile() error = %v", err)
	}
	got, err := GetBillingProfile(context.Background(), db, " "+user.ID+" ")
	if err != nil {
		t.Fatalf("GetBillingProfile() error = %v", err)
	}
	if got.ID != created.ID {
		t.Fatalf("got profile ID = %s, want %s", got.ID, created.ID)
	}
}

func testUpsertBillingProfileValidation(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)

	tests := []struct {
		name    string
		mutate  func(*UpsertBillingProfileInput)
		wantErr string
	}{
		{
			name: "blank billing name",
			mutate: func(input *UpsertBillingProfileInput) {
				input.BillingName = " "
			},
			wantErr: "billing name is required",
		},
		{
			name: "invalid email",
			mutate: func(input *UpsertBillingProfileInput) {
				input.BillingEmail = "not-an-email"
			},
			wantErr: "billing email is invalid",
		},
		{
			name: "Romanian company with no fiscal code",
			mutate: func(input *UpsertBillingProfileInput) {
				input.EntityType = BillingEntityCompany
				input.CountryCode = "ro"
				input.FiscalCode = stringPtr(" ")
			},
			wantErr: "fiscal code is required for Romanian companies",
		},
		{
			name: "invalid country code",
			mutate: func(input *UpsertBillingProfileInput) {
				input.CountryCode = "ROU"
			},
			wantErr: "country code is invalid",
		},
		{
			name: "unknown ISO country code",
			mutate: func(input *UpsertBillingProfileInput) {
				input.CountryCode = "ZZ"
			},
			wantErr: "country code is invalid",
		},
		{
			name: "overlong billing name",
			mutate: func(input *UpsertBillingProfileInput) {
				input.BillingName = strings.Repeat("a", 256)
			},
			wantErr: "billing name must be 255 characters or fewer",
		},
		{
			name: "overlong optional registration number",
			mutate: func(input *UpsertBillingProfileInput) {
				input.EntityType = BillingEntityCompany
				input.CountryCode = "DE"
				input.RegistrationNumber = stringPtr(strings.Repeat("a", 121))
			},
			wantErr: "registration number must be 120 characters or fewer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := validBillingProfileInput(user.ID)
			tt.mutate(&input)
			_, err := UpsertBillingProfile(context.Background(), db, input)
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("UpsertBillingProfile() error = %v, want %q", err, tt.wantErr)
			}
		})
	}
}

func testGrantSignupBonusCreatesNonExpiringBucket(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	now := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)

	bucket, err := GrantSignupBonus(context.Background(), db, GrantSignupBonusInput{
		UserID:  user.ID,
		Credits: 250,
		Now:     now,
	})
	if err != nil {
		t.Fatalf("GrantSignupBonus() error = %v", err)
	}
	if bucket.UserID != user.ID ||
		bucket.SourceType != CreditSourceSignupBonus ||
		bucket.CreditsGranted != 250 ||
		bucket.CreditsRemaining != 250 ||
		!bucket.ValidFrom.Equal(now) ||
		bucket.ExpiresAt != nil {
		t.Fatalf("unexpected signup bucket: %#v", bucket)
	}
	assertBillingCount(t, db, &CreditLedgerEntry{}, 1, "bucket_id = ? AND entry_type = ? AND credits_delta = ?", bucket.ID, CreditLedgerEntryGrant, 250)
}

func testGrantSignupBonusIsIdempotent(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	now := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)

	first, err := GrantSignupBonus(context.Background(), db, GrantSignupBonusInput{
		UserID:  user.ID,
		Credits: 250,
		Now:     now,
	})
	if err != nil {
		t.Fatalf("first GrantSignupBonus() error = %v", err)
	}
	second, err := GrantSignupBonus(context.Background(), db, GrantSignupBonusInput{
		UserID:  user.ID,
		Credits: 500,
		Now:     now.Add(time.Hour),
	})
	if err != nil {
		t.Fatalf("second GrantSignupBonus() error = %v", err)
	}
	if second.ID != first.ID {
		t.Fatalf("second bucket id = %s, want %s", second.ID, first.ID)
	}
	if second.CreditsGranted != 250 || second.CreditsRemaining != 250 {
		t.Fatalf("idempotent bucket credits = granted %d remaining %d, want original 250/250", second.CreditsGranted, second.CreditsRemaining)
	}
	assertBillingCount(t, db, &CreditBucket{}, 1, "user_id = ? AND source_type = ?", user.ID, CreditSourceSignupBonus)
	assertBillingCount(t, db, &CreditLedgerEntry{}, 1, "user_id = ? AND idempotency_key = ?", user.ID, "signup_bonus:"+user.ID)
}

func testGrantCreditsUsesAdvisoryIdempotencyLock(t *testing.T) {
	source, err := os.ReadFile("repository.go")
	if err != nil {
		t.Fatalf("read repository.go: %v", err)
	}
	if !strings.Contains(string(source), "pg_advisory_xact_lock(hashtext(?))") {
		t.Fatal("grant idempotency does not use a transaction-scoped advisory lock")
	}
}

func testCreateCreditOrderUsesPurchaseQuote(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)

	order, err := CreateCreditOrder(context.Background(), db, CreateCreditOrderInput{
		UserID:  user.ID,
		Credits: 5000,
	})
	if err != nil {
		t.Fatalf("CreateCreditOrder() error = %v", err)
	}
	if order.Credits != 5000 ||
		order.AmountCents != 4750 ||
		order.Currency != "EUR" ||
		order.PricingTier != CreditPurchaseTier2 ||
		order.UnitAmountCents != 950 ||
		order.Provider != BillingProviderStripe ||
		order.Status != OrderStatusPending {
		t.Fatalf("unexpected order: %#v", order)
	}
}

func testCreateCreditOrderRejectsInvalidCredits(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)

	_, err := CreateCreditOrder(context.Background(), db, CreateCreditOrderInput{
		UserID:  user.ID,
		Credits: 1500,
	})
	if err == nil {
		t.Fatal("CreateCreditOrder() error = nil, want invalid credit amount error")
	}
}

func testAttachCreditOrderCheckoutSession(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	order, err := CreateCreditOrder(context.Background(), db, CreateCreditOrderInput{
		UserID:  user.ID,
		Credits: 1000,
	})
	if err != nil {
		t.Fatalf("create order: %v", err)
	}

	if err := AttachCreditOrderCheckoutSession(context.Background(), db, order.ID, "cs_test_attach"); err != nil {
		t.Fatalf("AttachCreditOrderCheckoutSession() error = %v", err)
	}

	var got BillingOrder
	if err := db.First(&got, "id = ?", order.ID).Error; err != nil {
		t.Fatalf("load order: %v", err)
	}
	if got.ProviderCheckoutSessionID == nil || *got.ProviderCheckoutSessionID != "cs_test_attach" {
		t.Fatalf("checkout session = %#v, want cs_test_attach", got.ProviderCheckoutSessionID)
	}
}

func testAttachCreditOrderCheckoutSessionRejectsConflictingSession(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	order, err := CreateCreditOrder(context.Background(), db, CreateCreditOrderInput{
		UserID:  user.ID,
		Credits: 1000,
	})
	if err != nil {
		t.Fatalf("create order: %v", err)
	}
	if err := AttachCreditOrderCheckoutSession(context.Background(), db, order.ID, "cs_test_first"); err != nil {
		t.Fatalf("attach first checkout session: %v", err)
	}

	err = AttachCreditOrderCheckoutSession(context.Background(), db, order.ID, "cs_test_second")
	if err == nil {
		t.Error("AttachCreditOrderCheckoutSession() error = nil, want conflicting checkout session error")
	}

	var got BillingOrder
	if loadErr := db.First(&got, "id = ?", order.ID).Error; loadErr != nil {
		t.Fatalf("load order: %v", loadErr)
	}
	if got.ProviderCheckoutSessionID == nil || *got.ProviderCheckoutSessionID != "cs_test_first" {
		t.Fatalf("checkout session = %#v, want cs_test_first", got.ProviderCheckoutSessionID)
	}
}

func testMarkCreditOrderPaidAndGrantCreditsCreatesNonExpiringBucket(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	now := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)
	order, err := CreateCreditOrder(context.Background(), db, CreateCreditOrderInput{
		UserID:  user.ID,
		Credits: 2000,
	})
	if err != nil {
		t.Fatalf("create order: %v", err)
	}
	paidAt := now.Add(time.Hour)

	bucket, err := MarkCreditOrderPaidAndGrantCredits(context.Background(), db, MarkCreditOrderPaidInput{
		OrderID:                 order.ID,
		ProviderPaymentIntentID: stringPtr("pi_test"),
		PaidAt:                  paidAt,
	})
	if err != nil {
		t.Fatalf("MarkCreditOrderPaidAndGrantCredits() error = %v", err)
	}
	if bucket.UserID != user.ID ||
		bucket.SourceType != CreditSourceTopupPurchase ||
		bucket.OrderID == nil ||
		*bucket.OrderID != order.ID ||
		bucket.CreditsGranted != 2000 ||
		bucket.CreditsRemaining != 2000 ||
		bucket.ExpiresAt != nil {
		t.Fatalf("unexpected purchase bucket: %#v", bucket)
	}
	var gotOrder BillingOrder
	if err := db.First(&gotOrder, "id = ?", order.ID).Error; err != nil {
		t.Fatalf("load order: %v", err)
	}
	if gotOrder.Status != OrderStatusPaid ||
		gotOrder.PaidAt == nil ||
		!gotOrder.PaidAt.Equal(paidAt) ||
		gotOrder.ProviderPaymentIntentID == nil ||
		*gotOrder.ProviderPaymentIntentID != "pi_test" {
		t.Fatalf("unexpected paid order: %#v", gotOrder)
	}
	assertBillingCount(t, db, &CreditLedgerEntry{}, 1, "bucket_id = ? AND entry_type = ? AND credits_delta = ?", bucket.ID, CreditLedgerEntryPurchase, 2000)
}

func testMarkCreditOrderPaidAndGrantCreditsIdempotent(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	now := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)
	order, err := CreateCreditOrder(context.Background(), db, CreateCreditOrderInput{
		UserID:  user.ID,
		Credits: 10000,
	})
	if err != nil {
		t.Fatalf("create order: %v", err)
	}
	if err := AttachCreditOrderCheckoutSession(context.Background(), db, order.ID, "cs_test_123"); err != nil {
		t.Fatalf("attach checkout session: %v", err)
	}
	input := MarkCreditOrderPaidInput{
		OrderID:                   order.ID,
		ProviderCheckoutSessionID: stringPtr("cs_test_123"),
		ProviderPaymentIntentID:   stringPtr("pi_test_123"),
		PaidAt:                    now,
	}

	first, err := MarkCreditOrderPaidAndGrantCredits(context.Background(), db, input)
	if err != nil {
		t.Fatalf("first mark paid: %v", err)
	}
	input.PaidAt = now.Add(time.Hour)
	second, err := MarkCreditOrderPaidAndGrantCredits(context.Background(), db, input)
	if err != nil {
		t.Fatalf("second mark paid: %v", err)
	}
	if second.ID != first.ID {
		t.Fatalf("second bucket id = %s, want %s", second.ID, first.ID)
	}

	balance, err := AvailableCredits(context.Background(), db, user.ID, now)
	if err != nil {
		t.Fatalf("AvailableCredits() error = %v", err)
	}
	if balance.Available != 10000 {
		t.Fatalf("available credits = %d, want 10000", balance.Available)
	}
	assertBillingCount(t, db, &CreditBucket{}, 1, "order_id = ?", order.ID)
	assertBillingCount(t, db, &CreditLedgerEntry{}, 1, "idempotency_key = ?", "topup_paid:"+order.ID.String())
}

func testProviderPaymentIntentCannotGrantTwoOrders(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	now := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)
	firstOrder, err := CreateCreditOrder(context.Background(), db, CreateCreditOrderInput{UserID: user.ID, Credits: 1000})
	if err != nil {
		t.Fatalf("create first order: %v", err)
	}
	secondOrder, err := CreateCreditOrder(context.Background(), db, CreateCreditOrderInput{UserID: user.ID, Credits: 1000})
	if err != nil {
		t.Fatalf("create second order: %v", err)
	}
	paymentIntentID := stringPtr("pi_duplicate")

	if _, err := MarkCreditOrderPaidAndGrantCredits(context.Background(), db, MarkCreditOrderPaidInput{
		OrderID:                 firstOrder.ID,
		ProviderPaymentIntentID: paymentIntentID,
		PaidAt:                  now,
	}); err != nil {
		t.Fatalf("first mark paid: %v", err)
	}
	_, err = MarkCreditOrderPaidAndGrantCredits(context.Background(), db, MarkCreditOrderPaidInput{
		OrderID:                 secondOrder.ID,
		ProviderPaymentIntentID: paymentIntentID,
		PaidAt:                  now.Add(time.Hour),
	})
	if !errors.Is(err, ErrProviderMetadataConflict) {
		t.Fatalf("second mark paid error = %v, want ErrProviderMetadataConflict", err)
	}
	assertBillingCount(t, db, &CreditBucket{}, 1, "source_type = ?", CreditSourceTopupPurchase)
	assertBillingCount(t, db, &CreditLedgerEntry{}, 1, "entry_type = ? AND credits_delta = ?", CreditLedgerEntryPurchase, 1000)
	assertBillingCount(t, db, &BillingOrder{}, 1, "status = ? AND provider_payment_intent_id = ?", OrderStatusPaid, *paymentIntentID)

	var second BillingOrder
	if err := db.First(&second, "id = ?", secondOrder.ID).Error; err != nil {
		t.Fatalf("load second order: %v", err)
	}
	if second.Status != OrderStatusPending || second.ProviderPaymentIntentID != nil {
		t.Fatalf("second order was mutated by duplicate payment intent: %#v", second)
	}
}

func testMarkCreditOrderFailedDoesNotGrantCredits(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	failedAt := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)
	order, err := CreateCreditOrder(context.Background(), db, CreateCreditOrderInput{UserID: user.ID, Credits: 1000})
	if err != nil {
		t.Fatalf("create order: %v", err)
	}

	if err := MarkCreditOrderFailed(context.Background(), db, order.ID, failedAt); err != nil {
		t.Fatalf("MarkCreditOrderFailed() error = %v", err)
	}

	var got BillingOrder
	if err := db.First(&got, "id = ?", order.ID).Error; err != nil {
		t.Fatalf("load order: %v", err)
	}
	if got.Status != OrderStatusFailed || got.FailedAt == nil || !got.FailedAt.Equal(failedAt) {
		t.Fatalf("unexpected failed order: %#v", got)
	}
	assertBillingCount(t, db, &CreditBucket{}, 0, "order_id = ?", order.ID)
}

func testMarkCreditOrderCanceledDoesNotGrantCredits(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	canceledAt := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)
	order, err := CreateCreditOrder(context.Background(), db, CreateCreditOrderInput{UserID: user.ID, Credits: 1000})
	if err != nil {
		t.Fatalf("create order: %v", err)
	}

	if err := MarkCreditOrderCanceled(context.Background(), db, order.ID, canceledAt); err != nil {
		t.Fatalf("MarkCreditOrderCanceled() error = %v", err)
	}

	var got BillingOrder
	if err := db.First(&got, "id = ?", order.ID).Error; err != nil {
		t.Fatalf("load order: %v", err)
	}
	if got.Status != OrderStatusCanceled || got.CanceledAt == nil || !got.CanceledAt.Equal(canceledAt) {
		t.Fatalf("unexpected canceled order: %#v", got)
	}
	assertBillingCount(t, db, &CreditBucket{}, 0, "order_id = ?", order.ID)
}

func testAvailableCreditsIgnoresExpiredAndVoidedBuckets(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	now := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)
	createCreditBucket(t, db, user.ID, CreditSourceTopupPurchase, 300, now, nil)
	createCreditBucket(t, db, user.ID, CreditSourceSignupBonus, 200, now, ptrTime(now.Add(time.Hour)))
	createCreditBucket(t, db, user.ID, CreditSourceRefund, 500, now, ptrTime(now.Add(-time.Hour)))
	voided := createCreditBucket(t, db, user.ID, CreditSourceAdjustment, 700, now, ptrTime(now.Add(time.Hour)))
	if err := db.Model(&CreditBucket{}).Where("id = ?", voided.ID).Update("voided_at", now).Error; err != nil {
		t.Fatalf("void bucket: %v", err)
	}

	balance, err := AvailableCredits(context.Background(), db, user.ID, now)
	if err != nil {
		t.Fatalf("AvailableCredits() error = %v", err)
	}
	if balance.Available != 500 {
		t.Fatalf("available credits = %d, want 500", balance.Available)
	}
}

func testAvailableCreditsIgnoresFutureBuckets(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	now := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)
	createCreditBucket(t, db, user.ID, CreditSourceTopupPurchase, 300, now, nil)
	createCreditBucket(t, db, user.ID, CreditSourceRefund, 700, now.Add(time.Hour), nil)

	balance, err := AvailableCredits(context.Background(), db, user.ID, now)
	if err != nil {
		t.Fatalf("AvailableCredits() error = %v", err)
	}
	if balance.Available != 300 {
		t.Fatalf("available credits = %d, want 300", balance.Available)
	}
}

func testAdjustCreditsAddsAdjustmentBucketAndLedgerEntry(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	now := time.Date(2026, 6, 11, 12, 0, 0, 0, time.UTC)

	balance, err := AdjustCredits(context.Background(), db, AdjustCreditsInput{
		UserID:         user.ID,
		Delta:          250,
		Now:            now,
		IdempotencyKey: "admin_adjustment:add",
	})
	if err != nil {
		t.Fatalf("AdjustCredits() error = %v", err)
	}
	if balance.Available != 250 {
		t.Fatalf("available credits = %d, want 250", balance.Available)
	}

	var bucket CreditBucket
	if err := db.First(&bucket, "user_id = ? AND source_type = ?", user.ID, CreditSourceAdjustment).Error; err != nil {
		t.Fatalf("load adjustment bucket: %v", err)
	}
	if bucket.CreditsGranted != 250 || bucket.CreditsRemaining != 250 || !bucket.ValidFrom.Equal(now) {
		t.Fatalf("unexpected adjustment bucket: %#v", bucket)
	}
	assertBillingCount(t, db, &CreditLedgerEntry{}, 1,
		"user_id = ? AND bucket_id = ? AND entry_type = ? AND credits_delta = ? AND idempotency_key = ?",
		user.ID, bucket.ID, CreditLedgerEntryAdjustment, 250, "admin_adjustment:add")
}

func testAdjustCreditsSubtractsFromAvailableBuckets(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	now := time.Date(2026, 6, 11, 12, 0, 0, 0, time.UTC)
	bucket := createCreditBucket(t, db, user.ID, CreditSourceAdjustment, 500, now.Add(-time.Hour), nil)

	balance, err := AdjustCredits(context.Background(), db, AdjustCreditsInput{
		UserID:         user.ID,
		Delta:          -200,
		Now:            now,
		IdempotencyKey: "admin_adjustment:subtract",
	})
	if err != nil {
		t.Fatalf("AdjustCredits() error = %v", err)
	}
	if balance.Available != 300 {
		t.Fatalf("available credits = %d, want 300", balance.Available)
	}
	assertBucketRemaining(t, db, bucket.ID, 300)
	assertBillingCount(t, db, &CreditLedgerEntry{}, 1,
		"user_id = ? AND bucket_id = ? AND entry_type = ? AND credits_delta = ? AND idempotency_key = ?",
		user.ID, bucket.ID, CreditLedgerEntryAdjustment, -200, "admin_adjustment:subtract")
}

func testAdjustCreditsRejectsInsufficientBalance(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	now := time.Date(2026, 6, 11, 12, 0, 0, 0, time.UTC)
	bucket := createCreditBucket(t, db, user.ID, CreditSourceAdjustment, 100, now.Add(-time.Hour), nil)

	_, err := AdjustCredits(context.Background(), db, AdjustCreditsInput{
		UserID:         user.ID,
		Delta:          -101,
		Now:            now,
		IdempotencyKey: "admin_adjustment:too-much",
	})
	if !errors.Is(err, ErrInsufficientCredits) {
		t.Fatalf("AdjustCredits() error = %v, want ErrInsufficientCredits", err)
	}
	assertBucketRemaining(t, db, bucket.ID, 100)
	assertBillingCount(t, db, &CreditLedgerEntry{}, 0, "idempotency_key = ?", "admin_adjustment:too-much")
}

func testListCreditLedgerTransactionsFiltersAndSorts(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	other := createBillingTestUser(t, db)
	now := time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC)
	bucket := createCreditBucket(t, db, user.ID, CreditSourceAdjustment, 5000, now.Add(-time.Hour), nil)
	otherBucket := createCreditBucket(t, db, other.ID, CreditSourceAdjustment, 1000, now.Add(-time.Hour), nil)

	olderPurchase := createCreditLedgerEntry(t, db, CreditLedgerEntry{
		UserID: user.ID, BucketID: &bucket.ID, EntryType: CreditLedgerEntryPurchase,
		CreditsDelta: 1000, IdempotencyKey: "list:purchase:older", CreatedAt: now.Add(-30 * time.Minute),
	})
	debit := createCreditLedgerEntry(t, db, CreditLedgerEntry{
		UserID: user.ID, BucketID: &bucket.ID, EntryType: CreditLedgerEntryDebit,
		CreditsDelta: -25, RelatedJobID: uuidPtr(uuid.New()), IdempotencyKey: "list:debit",
		CreatedAt: now.Add(-20 * time.Minute),
	})
	newerPurchase := createCreditLedgerEntry(t, db, CreditLedgerEntry{
		UserID: user.ID, BucketID: &bucket.ID, EntryType: CreditLedgerEntryPurchase,
		CreditsDelta: 2000, IdempotencyKey: "list:purchase:newer", CreatedAt: now.Add(-10 * time.Minute),
	})
	createCreditLedgerEntry(t, db, CreditLedgerEntry{
		UserID: user.ID, BucketID: &bucket.ID, EntryType: CreditLedgerEntryGrant,
		CreditsDelta: 500, IdempotencyKey: "list:grant", CreatedAt: now.Add(-5 * time.Minute),
	})
	createCreditLedgerEntry(t, db, CreditLedgerEntry{
		UserID: other.ID, BucketID: &otherBucket.ID, EntryType: CreditLedgerEntryPurchase,
		CreditsDelta: 1000, IdempotencyKey: "list:other", CreatedAt: now,
	})

	entryType := CreditLedgerEntryPurchase
	page, err := ListCreditLedgerTransactions(context.Background(), db, ListCreditLedgerTransactionsInput{
		UserID:      user.ID,
		EntryType:   &entryType,
		CreatedFrom: timePtr(now.Add(-40 * time.Minute)),
		CreatedTo:   timePtr(now),
		Size:        20,
		Sort:        "desc",
	})
	if err != nil {
		t.Fatalf("ListCreditLedgerTransactions() error = %v", err)
	}
	if page.NextCursor != nil {
		t.Fatalf("NextCursor = %#v, want nil", page.NextCursor)
	}
	if got := ledgerEntryIDs(page.Entries); len(got) != 2 || got[0] != newerPurchase.ID || got[1] != olderPurchase.ID {
		t.Fatalf("entries = %v, want newest and older purchase only", got)
	}

	entryType = CreditLedgerEntryDebit
	page, err = ListCreditLedgerTransactions(context.Background(), db, ListCreditLedgerTransactionsInput{
		UserID: user.ID, EntryType: &entryType, Size: 20, Sort: "desc",
	})
	if err != nil {
		t.Fatalf("ListCreditLedgerTransactions(debit) error = %v", err)
	}
	if got := ledgerEntryIDs(page.Entries); len(got) != 1 || got[0] != debit.ID {
		t.Fatalf("debit entries = %v, want %s", got, debit.ID)
	}
}

func testListCreditLedgerTransactionsCursorPaginates(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	now := time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC)
	bucket := createCreditBucket(t, db, user.ID, CreditSourceAdjustment, 5000, now.Add(-time.Hour), nil)
	oldest := createCreditLedgerEntry(t, db, CreditLedgerEntry{
		UserID: user.ID, BucketID: &bucket.ID, EntryType: CreditLedgerEntryDebit,
		CreditsDelta: -1, IdempotencyKey: "cursor:oldest", CreatedAt: now.Add(-30 * time.Minute),
	})
	middle := createCreditLedgerEntry(t, db, CreditLedgerEntry{
		UserID: user.ID, BucketID: &bucket.ID, EntryType: CreditLedgerEntryDebit,
		CreditsDelta: -2, IdempotencyKey: "cursor:middle", CreatedAt: now.Add(-20 * time.Minute),
	})
	newest := createCreditLedgerEntry(t, db, CreditLedgerEntry{
		UserID: user.ID, BucketID: &bucket.ID, EntryType: CreditLedgerEntryDebit,
		CreditsDelta: -3, IdempotencyKey: "cursor:newest", CreatedAt: now.Add(-10 * time.Minute),
	})

	first, err := ListCreditLedgerTransactions(context.Background(), db, ListCreditLedgerTransactionsInput{
		UserID: user.ID, Size: 2, Sort: "desc",
	})
	if err != nil {
		t.Fatalf("first page error = %v", err)
	}
	if got := ledgerEntryIDs(first.Entries); len(got) != 2 || got[0] != newest.ID || got[1] != middle.ID {
		t.Fatalf("first entries = %v, want newest,middle", got)
	}
	if first.NextCursor == nil || first.NextCursor.ID != middle.ID || first.NextCursor.Sort != "desc" {
		t.Fatalf("first cursor = %#v, want middle desc cursor", first.NextCursor)
	}

	second, err := ListCreditLedgerTransactions(context.Background(), db, ListCreditLedgerTransactionsInput{
		UserID: user.ID, Size: 2, Sort: "desc", Cursor: first.NextCursor,
	})
	if err != nil {
		t.Fatalf("second page error = %v", err)
	}
	if got := ledgerEntryIDs(second.Entries); len(got) != 1 || got[0] != oldest.ID {
		t.Fatalf("second entries = %v, want oldest", got)
	}
	if second.NextCursor != nil {
		t.Fatalf("second cursor = %#v, want nil", second.NextCursor)
	}
}

func testListCreditLedgerTransactionsValidation(t *testing.T) {
	db := billingRepositoryTx(t)
	if _, err := ListCreditLedgerTransactions(context.Background(), nil, ListCreditLedgerTransactionsInput{UserID: uuid.NewString(), Size: 20, Sort: "desc"}); err == nil {
		t.Fatal("nil db succeeded")
	}
	if _, err := ListCreditLedgerTransactions(context.Background(), db, ListCreditLedgerTransactionsInput{Size: 20, Sort: "desc"}); err == nil {
		t.Fatal("empty user id succeeded")
	}
	grant := CreditLedgerEntryGrant
	if _, err := ListCreditLedgerTransactions(context.Background(), db, ListCreditLedgerTransactionsInput{UserID: uuid.NewString(), EntryType: &grant, Size: 20, Sort: "desc"}); err == nil {
		t.Fatal("unsupported entry type succeeded")
	}
	if _, err := ListCreditLedgerTransactions(context.Background(), db, ListCreditLedgerTransactionsInput{UserID: uuid.NewString(), Size: 0, Sort: "desc"}); err == nil {
		t.Fatal("invalid size succeeded")
	}
	if _, err := ListCreditLedgerTransactions(context.Background(), db, ListCreditLedgerTransactionsInput{UserID: uuid.NewString(), Size: 20, Sort: "newest"}); err == nil {
		t.Fatal("invalid sort succeeded")
	}
}

func testListBillingOrdersFiltersAndSorts(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	other := createBillingTestUser(t, db)
	now := time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC)

	olderPaid := createBillingOrderForList(t, db, user.ID, 1000, OrderStatusPaid, now.Add(-30*time.Minute), timePtr(now.Add(-29*time.Minute)))
	pending := createBillingOrderForList(t, db, user.ID, 5000, OrderStatusPending, now.Add(-20*time.Minute), nil)
	newerPaid := createBillingOrderForList(t, db, user.ID, 10000, OrderStatusPaid, now.Add(-10*time.Minute), timePtr(now.Add(-9*time.Minute)))
	createBillingOrderForList(t, db, other.ID, 5000, OrderStatusPaid, now.Add(-5*time.Minute), timePtr(now.Add(-4*time.Minute)))

	status := OrderStatusPaid
	page, err := ListBillingOrders(context.Background(), db, ListBillingOrdersInput{
		UserID:      user.ID,
		Status:      &status,
		CreatedFrom: timePtr(now.Add(-40 * time.Minute)),
		CreatedTo:   timePtr(now),
		Size:        20,
		Sort:        "desc",
	})
	if err != nil {
		t.Fatalf("ListBillingOrders() error = %v", err)
	}
	if page.NextCursor != nil {
		t.Fatalf("NextCursor = %#v, want nil", page.NextCursor)
	}
	if got := billingOrderIDs(page.Orders); len(got) != 2 || got[0] != newerPaid.ID || got[1] != olderPaid.ID {
		t.Fatalf("orders = %v, want newer and older paid only", got)
	}

	status = OrderStatusPending
	page, err = ListBillingOrders(context.Background(), db, ListBillingOrdersInput{
		UserID: user.ID, Status: &status, Size: 20, Sort: "desc",
	})
	if err != nil {
		t.Fatalf("ListBillingOrders(pending) error = %v", err)
	}
	if got := billingOrderIDs(page.Orders); len(got) != 1 || got[0] != pending.ID {
		t.Fatalf("pending orders = %v, want %s", got, pending.ID)
	}
}

func testListBillingOrdersCursorPaginates(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	now := time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC)
	oldest := createBillingOrderForList(t, db, user.ID, 1000, OrderStatusPending, now.Add(-30*time.Minute), nil)
	middle := createBillingOrderForList(t, db, user.ID, 5000, OrderStatusPaid, now.Add(-20*time.Minute), timePtr(now.Add(-19*time.Minute)))
	newest := createBillingOrderForList(t, db, user.ID, 10000, OrderStatusFailed, now.Add(-10*time.Minute), nil)

	first, err := ListBillingOrders(context.Background(), db, ListBillingOrdersInput{
		UserID: user.ID, Size: 2, Sort: "desc",
	})
	if err != nil {
		t.Fatalf("first page error = %v", err)
	}
	if got := billingOrderIDs(first.Orders); len(got) != 2 || got[0] != newest.ID || got[1] != middle.ID {
		t.Fatalf("first orders = %v, want newest,middle", got)
	}
	if first.NextCursor == nil || first.NextCursor.ID != middle.ID || first.NextCursor.Sort != "desc" {
		t.Fatalf("first cursor = %#v, want middle desc cursor", first.NextCursor)
	}

	second, err := ListBillingOrders(context.Background(), db, ListBillingOrdersInput{
		UserID: user.ID, Size: 2, Sort: "desc", Cursor: first.NextCursor,
	})
	if err != nil {
		t.Fatalf("second page error = %v", err)
	}
	if got := billingOrderIDs(second.Orders); len(got) != 1 || got[0] != oldest.ID {
		t.Fatalf("second orders = %v, want oldest", got)
	}
	if second.NextCursor != nil {
		t.Fatalf("second cursor = %#v, want nil", second.NextCursor)
	}
}

func testListAdminBillingOrdersListsAllUsersAndFiltersByUser(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	other := createBillingTestUser(t, db)
	now := time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC)

	older := createBillingOrderForList(t, db, user.ID, 1000, OrderStatusPaid, now.Add(-30*time.Minute), timePtr(now.Add(-29*time.Minute)))
	newer := createBillingOrderForList(t, db, other.ID, 5000, OrderStatusPending, now.Add(-10*time.Minute), nil)

	page, err := ListAdminBillingOrders(context.Background(), db, ListAdminBillingOrdersInput{
		Size: 20,
		Sort: "desc",
	})
	if err != nil {
		t.Fatalf("ListAdminBillingOrders() error = %v", err)
	}
	if got := billingOrderIDs(page.Orders); len(got) != 2 || got[0] != newer.ID || got[1] != older.ID {
		t.Fatalf("orders = %v, want newer,older across users", got)
	}
	if page.Orders[0].User.Email != other.Email || page.Orders[1].User.Email != user.Email {
		t.Fatalf("preloaded users = %#v, %#v; want order owners", page.Orders[0].User, page.Orders[1].User)
	}

	page, err = ListAdminBillingOrders(context.Background(), db, ListAdminBillingOrdersInput{
		UserID: &user.ID,
		Size:   20,
		Sort:   "desc",
	})
	if err != nil {
		t.Fatalf("ListAdminBillingOrders(user) error = %v", err)
	}
	if got := billingOrderIDs(page.Orders); len(got) != 1 || got[0] != older.ID {
		t.Fatalf("filtered orders = %v, want only %s", got, older.ID)
	}
}

func testListAdminBillingOrdersFiltersPaidOrdersWithoutInvoices(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	now := time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC)

	oldestPaidNoInvoice := createBillingOrderForList(t, db, user.ID, 1000, OrderStatusPaid, now.Add(-30*time.Minute), timePtr(now.Add(-29*time.Minute)))
	paidWithInvoice := createBillingOrderForList(t, db, user.ID, 5000, OrderStatusPaid, now.Add(-20*time.Minute), timePtr(now.Add(-19*time.Minute)))
	newestPaidNoInvoice := createBillingOrderForList(t, db, user.ID, 10000, OrderStatusPaid, now.Add(-10*time.Minute), timePtr(now.Add(-9*time.Minute)))
	createBillingOrderForList(t, db, user.ID, 5000, OrderStatusPending, now.Add(-5*time.Minute), nil)
	if _, err := CreateBillingInvoiceForPaidOrder(context.Background(), db, paidWithInvoice.ID, time.Date(2026, 6, 5, 0, 0, 0, 0, time.UTC)); err != nil {
		t.Fatalf("create billing invoice for paid order: %v", err)
	}

	status := OrderStatusPaid
	first, err := ListAdminBillingOrders(context.Background(), db, ListAdminBillingOrdersInput{
		Status:         &status,
		WithoutInvoice: true,
		Size:           1,
		Sort:           "desc",
	})
	if err != nil {
		t.Fatalf("ListAdminBillingOrders() error = %v", err)
	}
	if got := billingOrderIDs(first.Orders); len(got) != 1 || got[0] != newestPaidNoInvoice.ID {
		t.Fatalf("first orders = %v, want newest paid order without invoice", got)
	}
	if first.Orders[0].Invoice != nil {
		t.Fatalf("first invoice = %#v, want nil", first.Orders[0].Invoice)
	}
	if first.NextCursor == nil || first.NextCursor.ID != newestPaidNoInvoice.ID || first.NextCursor.Sort != "desc" {
		t.Fatalf("first cursor = %#v, want newest desc cursor", first.NextCursor)
	}

	second, err := ListAdminBillingOrders(context.Background(), db, ListAdminBillingOrdersInput{
		Status:         &status,
		WithoutInvoice: true,
		Size:           1,
		Sort:           "desc",
		Cursor:         first.NextCursor,
	})
	if err != nil {
		t.Fatalf("ListAdminBillingOrders(second) error = %v", err)
	}
	if got := billingOrderIDs(second.Orders); len(got) != 1 || got[0] != oldestPaidNoInvoice.ID {
		t.Fatalf("second orders = %v, want oldest paid order without invoice", got)
	}
	if second.NextCursor != nil {
		t.Fatalf("second cursor = %#v, want nil", second.NextCursor)
	}
}

func testListAdminBillingInvoicesListsAllUsersAndFilters(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	other := createBillingTestUser(t, db)
	now := time.Date(2026, 6, 11, 12, 0, 0, 0, time.UTC)

	older := createBillingInvoiceForList(t, db, user.ID, "SYN-LIST", now.Add(-30*time.Minute))
	middle := createBillingInvoiceForList(t, db, other.ID, "SYN-LIST", now.Add(-20*time.Minute))
	newer := createBillingInvoiceForList(t, db, user.ID, "SYN-LIST", now.Add(-10*time.Minute))

	page, err := ListAdminBillingInvoices(context.Background(), db, ListAdminBillingInvoicesInput{
		Size: 20,
		Sort: "desc",
	})
	if err != nil {
		t.Fatalf("ListAdminBillingInvoices() error = %v", err)
	}
	if got := billingInvoiceIDs(page.Invoices); len(got) != 3 || got[0] != newer.ID || got[1] != middle.ID || got[2] != older.ID {
		t.Fatalf("invoices = %v, want newer,middle,older across users", got)
	}

	page, err = ListAdminBillingInvoices(context.Background(), db, ListAdminBillingInvoicesInput{
		UserID: &user.ID,
		Size:   20,
		Sort:   "desc",
	})
	if err != nil {
		t.Fatalf("ListAdminBillingInvoices(user) error = %v", err)
	}
	if got := billingInvoiceIDs(page.Invoices); len(got) != 2 || got[0] != newer.ID || got[1] != older.ID {
		t.Fatalf("filtered invoices = %v, want newer and older for target user", got)
	}

	createdFrom := now.Add(-25 * time.Minute)
	createdTo := now.Add(-15 * time.Minute)
	page, err = ListAdminBillingInvoices(context.Background(), db, ListAdminBillingInvoicesInput{
		CreatedFrom: &createdFrom,
		CreatedTo:   &createdTo,
		Size:        20,
		Sort:        "desc",
	})
	if err != nil {
		t.Fatalf("ListAdminBillingInvoices(created range) error = %v", err)
	}
	if got := billingInvoiceIDs(page.Invoices); len(got) != 1 || got[0] != middle.ID {
		t.Fatalf("created range invoices = %v, want middle only", got)
	}
}

func testListAdminBillingInvoicesSearchesClientAndInvoiceNumber(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	now := time.Date(2026, 6, 11, 12, 0, 0, 0, time.UTC)

	ada := updateBillingInvoiceClient(
		t,
		db,
		createBillingInvoiceForList(t, db, user.ID, "SYN-SRCH", now.Add(-30*time.Minute)).ID,
		"Ada Lovelace",
		"ada@example.com",
	)
	grace := updateBillingInvoiceClient(
		t,
		db,
		createBillingInvoiceForList(t, db, user.ID, "SYN-SRCH", now.Add(-20*time.Minute)).ID,
		"Grace Hopper",
		"grace@example.com",
	)
	katherine := updateBillingInvoiceClient(
		t,
		db,
		createBillingInvoiceForList(t, db, user.ID, "SYN-SRCH", now.Add(-10*time.Minute)).ID,
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
			page, err := ListAdminBillingInvoices(context.Background(), db, ListAdminBillingInvoicesInput{
				Search: tt.search,
				Size:   20,
				Sort:   "desc",
			})
			if err != nil {
				t.Fatalf("ListAdminBillingInvoices(search=%q) error = %v", tt.search, err)
			}
			if got := billingInvoiceIDs(page.Invoices); len(got) != 1 || got[0] != tt.want {
				t.Fatalf("search %q invoices = %v, want %s", tt.search, got, tt.want)
			}
		})
	}
}

func testListAdminBillingInvoicesPaginatesByCreatedAt(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	now := time.Date(2026, 6, 11, 12, 0, 0, 0, time.UTC)

	oldest := createBillingInvoiceForList(t, db, user.ID, "SYN-PAGE", now.Add(-30*time.Minute))
	middle := createBillingInvoiceForList(t, db, user.ID, "SYN-PAGE", now.Add(-20*time.Minute))
	newest := createBillingInvoiceForList(t, db, user.ID, "SYN-PAGE", now.Add(-10*time.Minute))

	first, err := ListAdminBillingInvoices(context.Background(), db, ListAdminBillingInvoicesInput{
		UserID: &user.ID,
		Size:   2,
		Sort:   "desc",
	})
	if err != nil {
		t.Fatalf("first ListAdminBillingInvoices() error = %v", err)
	}
	if got := billingInvoiceIDs(first.Invoices); len(got) != 2 || got[0] != newest.ID || got[1] != middle.ID {
		t.Fatalf("first page invoices = %v, want newest,middle", got)
	}
	if first.NextCursor == nil {
		t.Fatal("first page next cursor = nil, want cursor")
	}

	second, err := ListAdminBillingInvoices(context.Background(), db, ListAdminBillingInvoicesInput{
		UserID: &user.ID,
		Cursor: first.NextCursor,
		Size:   2,
		Sort:   "desc",
	})
	if err != nil {
		t.Fatalf("second ListAdminBillingInvoices() error = %v", err)
	}
	if got := billingInvoiceIDs(second.Invoices); len(got) != 1 || got[0] != oldest.ID {
		t.Fatalf("second page invoices = %v, want oldest", got)
	}
	if second.NextCursor != nil {
		t.Fatalf("second page next cursor = %#v, want nil", second.NextCursor)
	}
}

func testListBillingOrdersValidation(t *testing.T) {
	db := billingRepositoryTx(t)
	if _, err := ListBillingOrders(context.Background(), nil, ListBillingOrdersInput{UserID: uuid.NewString(), Size: 20, Sort: "desc"}); err == nil {
		t.Fatal("nil db succeeded")
	}
	if _, err := ListBillingOrders(context.Background(), db, ListBillingOrdersInput{Size: 20, Sort: "desc"}); err == nil {
		t.Fatal("empty user id succeeded")
	}
	invalidStatus := OrderStatus("settled")
	if _, err := ListBillingOrders(context.Background(), db, ListBillingOrdersInput{UserID: uuid.NewString(), Status: &invalidStatus, Size: 20, Sort: "desc"}); err == nil {
		t.Fatal("invalid status succeeded")
	}
	if _, err := ListBillingOrders(context.Background(), db, ListBillingOrdersInput{UserID: uuid.NewString(), Size: 0, Sort: "desc"}); err == nil {
		t.Fatal("invalid size succeeded")
	}
	if _, err := ListBillingOrders(context.Background(), db, ListBillingOrdersInput{UserID: uuid.NewString(), Size: 20, Sort: "newest"}); err == nil {
		t.Fatal("invalid sort succeeded")
	}
	if _, err := ListBillingOrders(context.Background(), db, ListBillingOrdersInput{UserID: uuid.NewString(), Size: 20, Sort: "desc", Cursor: &BillingOrderCursor{ID: uuid.New(), CreatedAt: time.Now().UTC(), Sort: "asc"}}); err == nil {
		t.Fatal("cursor sort mismatch succeeded")
	}
}

func testListAdminBillingInvoicesValidation(t *testing.T) {
	db := billingRepositoryTx(t)
	if _, err := ListAdminBillingInvoices(context.Background(), nil, ListAdminBillingInvoicesInput{Size: 20, Sort: "desc"}); err == nil {
		t.Fatal("nil db succeeded")
	}
	blankUserID := " "
	if _, err := ListAdminBillingInvoices(context.Background(), db, ListAdminBillingInvoicesInput{UserID: &blankUserID, Size: 20, Sort: "desc"}); err == nil {
		t.Fatal("empty user id succeeded")
	}
	if _, err := ListAdminBillingInvoices(context.Background(), db, ListAdminBillingInvoicesInput{Size: 0, Sort: "desc"}); err == nil {
		t.Fatal("invalid size succeeded")
	}
	if _, err := ListAdminBillingInvoices(context.Background(), db, ListAdminBillingInvoicesInput{Size: 20, Sort: "newest"}); err == nil {
		t.Fatal("invalid sort succeeded")
	}
	createdFrom := time.Date(2026, 6, 11, 12, 0, 0, 0, time.UTC)
	createdTo := createdFrom.Add(-time.Minute)
	if _, err := ListAdminBillingInvoices(context.Background(), db, ListAdminBillingInvoicesInput{
		CreatedFrom: &createdFrom,
		CreatedTo:   &createdTo,
		Size:        20,
		Sort:        "desc",
	}); err == nil {
		t.Fatal("inverted created range succeeded")
	}
	if _, err := ListAdminBillingInvoices(context.Background(), db, ListAdminBillingInvoicesInput{
		Cursor: &BillingInvoiceCursor{CreatedAt: createdFrom, ID: uuid.Nil, Sort: "desc"},
		Size:   20,
		Sort:   "desc",
	}); err == nil {
		t.Fatal("invalid cursor succeeded")
	}
	if _, err := ListAdminBillingInvoices(context.Background(), db, ListAdminBillingInvoicesInput{
		Cursor: &BillingInvoiceCursor{CreatedAt: createdFrom, ID: uuid.New(), Sort: "asc"},
		Size:   20,
		Sort:   "desc",
	}); err == nil {
		t.Fatal("cursor sort mismatch succeeded")
	}
}

func testDebitCreditsForJobConsumesBucketsInPriorityOrder(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	now := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)
	signup := createCreditBucket(t, db, user.ID, CreditSourceSignupBonus, 5, now, ptrTime(now.AddDate(0, 0, 90)))
	purchased := createCreditBucket(t, db, user.ID, CreditSourceTopupPurchase, 5, now, nil)
	jobID := uuid.New()

	if err := DebitCreditsForJob(context.Background(), db, DebitCreditsInput{
		UserID:         user.ID,
		RelatedJobID:   jobID,
		Credits:        8,
		IdempotencyKey: "job-debit:" + jobID.String(),
		Now:            now,
	}); err != nil {
		t.Fatalf("DebitCreditsForJob() error = %v", err)
	}
	assertBucketRemaining(t, db, signup.ID, 0)
	assertBucketRemaining(t, db, purchased.ID, 2)
	assertBillingCount(t, db, &CreditLedgerEntry{}, 2, "related_job_id = ? AND entry_type = ?", jobID, CreditLedgerEntryDebit)
}

func testDebitCreditsForJobTxUsesExistingTransaction(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	now := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)
	bucket := createCreditBucket(t, db, user.ID, CreditSourceTopupPurchase, 10, now, nil)
	jobID := uuid.New()

	err := db.Transaction(func(tx *gorm.DB) error {
		return DebitCreditsForJobTx(context.Background(), tx, DebitCreditsInput{
			UserID:         user.ID,
			RelatedJobID:   jobID,
			Credits:        3,
			IdempotencyKey: "job-debit:" + jobID.String(),
			Now:            now,
		})
	})
	if err != nil {
		t.Fatalf("DebitCreditsForJobTx() error = %v", err)
	}
	assertBucketRemaining(t, db, bucket.ID, 7)
}

func testDebitCreditsForJobInsufficientCreditsDoesNotMutate(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	now := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)
	bucket := createCreditBucket(t, db, user.ID, CreditSourceTopupPurchase, 3, now, nil)
	jobID := uuid.New()

	err := DebitCreditsForJob(context.Background(), db, DebitCreditsInput{
		UserID:         user.ID,
		RelatedJobID:   jobID,
		Credits:        4,
		IdempotencyKey: "job-debit:" + jobID.String(),
		Now:            now,
	})
	if err == nil {
		t.Fatal("DebitCreditsForJob() error = nil, want insufficient credits")
	}
	assertBucketRemaining(t, db, bucket.ID, 3)
	assertBillingCount(t, db, &CreditLedgerEntry{}, 0, "related_job_id = ?", jobID)
}

func testDebitCreditsForJobRejectsNonPositiveCredits(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	now := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)

	err := DebitCreditsForJob(context.Background(), db, DebitCreditsInput{
		UserID:         user.ID,
		RelatedJobID:   uuid.New(),
		Credits:        0,
		IdempotencyKey: "job-debit:zero",
		Now:            now,
	})
	if err == nil {
		t.Fatal("DebitCreditsForJob() error = nil, want non-positive credits error")
	}
}

func testDebitCreditsForJobIgnoresExpiredAndVoidedBuckets(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	now := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)
	valid := createCreditBucket(t, db, user.ID, CreditSourceTopupPurchase, 5, now, nil)
	expired := createCreditBucket(t, db, user.ID, CreditSourceRefund, 5, now, ptrTime(now.Add(-time.Hour)))
	voided := createCreditBucket(t, db, user.ID, CreditSourceAdjustment, 5, now, ptrTime(now.Add(time.Hour)))
	if err := db.Model(&CreditBucket{}).Where("id = ?", voided.ID).Update("voided_at", now).Error; err != nil {
		t.Fatalf("void bucket: %v", err)
	}
	jobID := uuid.New()

	if err := DebitCreditsForJob(context.Background(), db, DebitCreditsInput{
		UserID:         user.ID,
		RelatedJobID:   jobID,
		Credits:        4,
		IdempotencyKey: "job-debit:" + jobID.String(),
		Now:            now,
	}); err != nil {
		t.Fatalf("DebitCreditsForJob() error = %v", err)
	}
	assertBucketRemaining(t, db, valid.ID, 1)
	assertBucketRemaining(t, db, expired.ID, 5)
	assertBucketRemaining(t, db, voided.ID, 5)
}

func testDebitCreditsForJobIgnoresFutureBuckets(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	now := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)
	valid := createCreditBucket(t, db, user.ID, CreditSourceTopupPurchase, 5, now, nil)
	future := createCreditBucket(t, db, user.ID, CreditSourceRefund, 5, now.Add(time.Hour), nil)
	jobID := uuid.New()

	err := DebitCreditsForJob(context.Background(), db, DebitCreditsInput{
		UserID:         user.ID,
		RelatedJobID:   jobID,
		Credits:        6,
		IdempotencyKey: "job-debit:" + jobID.String(),
		Now:            now,
	})
	if err == nil {
		t.Fatal("DebitCreditsForJob() error = nil, want insufficient credits without future bucket")
	}
	assertBucketRemaining(t, db, valid.ID, 5)
	assertBucketRemaining(t, db, future.ID, 5)
	assertBillingCount(t, db, &CreditLedgerEntry{}, 0, "related_job_id = ?", jobID)
}

func testDebitCreditsForJobIsIdempotent(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	now := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)
	bucket := createCreditBucket(t, db, user.ID, CreditSourceTopupPurchase, 10, now, nil)
	jobID := uuid.New()
	input := DebitCreditsInput{
		UserID:         user.ID,
		RelatedJobID:   jobID,
		Credits:        4,
		IdempotencyKey: "job-debit:" + jobID.String(),
		Now:            now,
	}
	if err := DebitCreditsForJob(context.Background(), db, input); err != nil {
		t.Fatalf("first debit: %v", err)
	}
	if err := DebitCreditsForJob(context.Background(), db, input); err != nil {
		t.Fatalf("second debit: %v", err)
	}
	assertBucketRemaining(t, db, bucket.ID, 6)
	assertBillingCount(t, db, &CreditLedgerEntry{}, 1, "related_job_id = ? AND entry_type = ?", jobID, CreditLedgerEntryDebit)
}

func testDebitCreditsForJobIdempotencyKeyTreatsWildcardsLiterally(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	now := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)
	bucket := createCreditBucket(t, db, user.ID, CreditSourceTopupPurchase, 10, now, nil)
	jobID := uuid.New()
	if err := DebitCreditsForJob(context.Background(), db, DebitCreditsInput{
		UserID:         user.ID,
		RelatedJobID:   jobID,
		Credits:        1,
		IdempotencyKey: "job-debit:" + jobID.String() + ":abc",
		Now:            now,
	}); err != nil {
		t.Fatalf("first debit: %v", err)
	}
	if err := DebitCreditsForJob(context.Background(), db, DebitCreditsInput{
		UserID:         user.ID,
		RelatedJobID:   jobID,
		Credits:        1,
		IdempotencyKey: "job-debit:" + jobID.String() + ":%",
		Now:            now,
	}); err != nil {
		t.Fatalf("wildcard debit: %v", err)
	}
	assertBucketRemaining(t, db, bucket.ID, 8)
	assertBillingCount(t, db, &CreditLedgerEntry{}, 2, "related_job_id = ? AND entry_type = ?", jobID, CreditLedgerEntryDebit)
}

func testRefundCreditsForJobRestoresActiveBucket(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	now := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)
	bucket := createCreditBucket(t, db, user.ID, CreditSourceTopupPurchase, 10, now, nil)
	jobID := uuid.New()
	if err := DebitCreditsForJob(context.Background(), db, DebitCreditsInput{
		UserID:         user.ID,
		RelatedJobID:   jobID,
		Credits:        4,
		IdempotencyKey: "job-debit:" + jobID.String(),
		Now:            now,
	}); err != nil {
		t.Fatalf("debit: %v", err)
	}

	if err := RefundCreditsForJob(context.Background(), db, user.ID, jobID, now); err != nil {
		t.Fatalf("refund: %v", err)
	}
	if err := RefundCreditsForJob(context.Background(), db, user.ID, jobID, now.Add(time.Hour)); err != nil {
		t.Fatalf("refund twice: %v", err)
	}
	assertBucketRemaining(t, db, bucket.ID, 10)
	assertBillingCount(t, db, &CreditLedgerEntry{}, 1, "related_job_id = ? AND entry_type = ?", jobID, CreditLedgerEntryRefund)
}

func testRefundCreditsForJobCreatesRefundBucketForExpiredCredit(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	now := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)
	originalExpiry := now.Add(time.Hour)
	bucket := createCreditBucket(t, db, user.ID, CreditSourceSignupBonus, 10, now, &originalExpiry)
	jobID := uuid.New()
	if err := DebitCreditsForJob(context.Background(), db, DebitCreditsInput{
		UserID:         user.ID,
		RelatedJobID:   jobID,
		Credits:        4,
		IdempotencyKey: "job-debit:" + jobID.String(),
		Now:            now,
	}); err != nil {
		t.Fatalf("debit: %v", err)
	}
	refundAt := now.Add(2 * time.Hour)

	if err := RefundCreditsForJob(context.Background(), db, user.ID, jobID, refundAt); err != nil {
		t.Fatalf("refund: %v", err)
	}
	assertBucketRemaining(t, db, bucket.ID, 6)
	var refund CreditBucket
	if err := db.Where("user_id = ? AND source_type = ?", user.ID, CreditSourceRefund).First(&refund).Error; err != nil {
		t.Fatalf("load refund bucket: %v", err)
	}
	if refund.CreditsGranted != 4 ||
		refund.CreditsRemaining != 4 ||
		refund.ExpiresAt == nil ||
		!refund.ExpiresAt.Equal(refundAt.AddDate(0, 0, 7)) {
		t.Fatalf("unexpected refund bucket: %#v", refund)
	}
}

func testRefundCreditsForJobUsesLaterOriginalExpiry(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	now := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)
	originalExpiry := now.AddDate(0, 0, 20)
	bucket := createCreditBucket(t, db, user.ID, CreditSourceSignupBonus, 10, now, &originalExpiry)
	jobID := uuid.New()
	if err := DebitCreditsForJob(context.Background(), db, DebitCreditsInput{
		UserID:         user.ID,
		RelatedJobID:   jobID,
		Credits:        4,
		IdempotencyKey: "job-debit:" + jobID.String(),
		Now:            now,
	}); err != nil {
		t.Fatalf("debit: %v", err)
	}
	if err := db.Model(&CreditBucket{}).Where("id = ?", bucket.ID).Update("voided_at", now.Add(time.Hour)).Error; err != nil {
		t.Fatalf("void bucket: %v", err)
	}
	refundAt := now.Add(2 * time.Hour)

	if err := RefundCreditsForJob(context.Background(), db, user.ID, jobID, refundAt); err != nil {
		t.Fatalf("refund: %v", err)
	}
	var refund CreditBucket
	if err := db.Where("user_id = ? AND source_type = ?", user.ID, CreditSourceRefund).First(&refund).Error; err != nil {
		t.Fatalf("load refund bucket: %v", err)
	}
	if refund.ExpiresAt == nil || !refund.ExpiresAt.Equal(originalExpiry) {
		t.Fatalf("refund expiry = %#v, want original expiry %v", refund.ExpiresAt, originalExpiry)
	}
}

func testRefundCreditsForJobCreatesRefundBucketForVoidedCredit(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	now := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)
	bucket := createCreditBucket(t, db, user.ID, CreditSourceAdjustment, 10, now, ptrTime(now.AddDate(0, 1, 0)))
	jobID := uuid.New()
	if err := DebitCreditsForJob(context.Background(), db, DebitCreditsInput{
		UserID:         user.ID,
		RelatedJobID:   jobID,
		Credits:        4,
		IdempotencyKey: "job-debit:" + jobID.String(),
		Now:            now,
	}); err != nil {
		t.Fatalf("debit: %v", err)
	}
	if err := db.Model(&CreditBucket{}).Where("id = ?", bucket.ID).Update("voided_at", now.Add(time.Hour)).Error; err != nil {
		t.Fatalf("void bucket: %v", err)
	}
	refundAt := now.Add(2 * time.Hour)

	if err := RefundCreditsForJob(context.Background(), db, user.ID, jobID, refundAt); err != nil {
		t.Fatalf("refund: %v", err)
	}
	assertBucketRemaining(t, db, bucket.ID, 6)
	var refund CreditBucket
	if err := db.Where("user_id = ? AND source_type = ?", user.ID, CreditSourceRefund).First(&refund).Error; err != nil {
		t.Fatalf("load refund bucket: %v", err)
	}
	if refund.CreditsGranted != 4 || refund.CreditsRemaining != 4 {
		t.Fatalf("unexpected refund bucket: %#v", refund)
	}
}

func testRefundCreditsForJobCreatesNonExpiringRefundBucketForPurchasedCredit(t *testing.T) {
	db := billingRepositoryTx(t)
	user := createBillingTestUser(t, db)
	now := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)
	bucket := createCreditBucket(t, db, user.ID, CreditSourceTopupPurchase, 10, now, nil)
	jobID := uuid.New()
	if err := DebitCreditsForJob(context.Background(), db, DebitCreditsInput{
		UserID:         user.ID,
		RelatedJobID:   jobID,
		Credits:        4,
		IdempotencyKey: "job-debit:" + jobID.String(),
		Now:            now,
	}); err != nil {
		t.Fatalf("debit: %v", err)
	}
	if err := db.Model(&CreditBucket{}).Where("id = ?", bucket.ID).Update("voided_at", now.Add(time.Hour)).Error; err != nil {
		t.Fatalf("void bucket: %v", err)
	}

	if err := RefundCreditsForJob(context.Background(), db, user.ID, jobID, now.Add(2*time.Hour)); err != nil {
		t.Fatalf("refund: %v", err)
	}
	var refund CreditBucket
	if err := db.Where("user_id = ? AND source_type = ?", user.ID, CreditSourceRefund).First(&refund).Error; err != nil {
		t.Fatalf("load refund bucket: %v", err)
	}
	if refund.ExpiresAt != nil {
		t.Fatalf("refund expires_at = %#v, want nil for purchased credit", refund.ExpiresAt)
	}
}

func createBillingOrderForList(t *testing.T, db *gorm.DB, userID string, credits int, status OrderStatus, createdAt time.Time, paidAt *time.Time) BillingOrder {
	t.Helper()
	order, err := CreateCreditOrder(context.Background(), db, CreateCreditOrderInput{
		UserID:  userID,
		Credits: credits,
	})
	if err != nil {
		t.Fatalf("CreateCreditOrder() error = %v", err)
	}
	updates := map[string]any{
		"status":     status,
		"created_at": createdAt.UTC(),
	}
	if paidAt != nil {
		updates["paid_at"] = paidAt.UTC()
	}
	if err := db.Model(&BillingOrder{}).Where("id = ?", order.ID).Updates(updates).Error; err != nil {
		t.Fatalf("update billing order for list: %v", err)
	}
	var out BillingOrder
	if err := db.First(&out, "id = ?", order.ID).Error; err != nil {
		t.Fatalf("reload billing order for list: %v", err)
	}
	return out
}

func createBillingInvoiceForList(t *testing.T, db *gorm.DB, userID string, invoiceSerie string, createdAt time.Time) BillingInvoice {
	t.Helper()
	invoice, err := CreateBillingInvoice(context.Background(), db, validInvoiceInput(userID, invoiceSerie))
	if err != nil {
		t.Fatalf("CreateBillingInvoice() error = %v", err)
	}
	if err := db.Model(&BillingInvoice{}).Where("id = ?", invoice.ID).Update("created_at", createdAt.UTC()).Error; err != nil {
		t.Fatalf("update billing invoice for list: %v", err)
	}
	var out BillingInvoice
	if err := db.First(&out, "id = ?", invoice.ID).Error; err != nil {
		t.Fatalf("reload billing invoice for list: %v", err)
	}
	return out
}

func updateBillingInvoiceClient(t *testing.T, db *gorm.DB, invoiceID uuid.UUID, billingName string, billingEmail string) BillingInvoice {
	t.Helper()
	if err := db.Model(&BillingInvoice{}).Where("id = ?", invoiceID).Updates(map[string]any{
		"billing_name":  billingName,
		"billing_email": billingEmail,
	}).Error; err != nil {
		t.Fatalf("update billing invoice client: %v", err)
	}
	var out BillingInvoice
	if err := db.First(&out, "id = ?", invoiceID).Error; err != nil {
		t.Fatalf("reload billing invoice client: %v", err)
	}
	return out
}

func createCreditLedgerEntry(t *testing.T, db *gorm.DB, entry CreditLedgerEntry) CreditLedgerEntry {
	t.Helper()
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = time.Now().UTC()
	}
	if err := db.Create(&entry).Error; err != nil {
		t.Fatalf("create credit ledger entry: %v", err)
	}
	return entry
}

func billingOrderIDs(orders []BillingOrder) []uuid.UUID {
	ids := make([]uuid.UUID, 0, len(orders))
	for _, order := range orders {
		ids = append(ids, order.ID)
	}
	return ids
}

func billingInvoiceIDs(invoices []BillingInvoice) []uuid.UUID {
	ids := make([]uuid.UUID, 0, len(invoices))
	for _, invoice := range invoices {
		ids = append(ids, invoice.ID)
	}
	return ids
}

func ledgerEntryIDs(entries []CreditLedgerEntry) []uuid.UUID {
	ids := make([]uuid.UUID, 0, len(entries))
	for _, entry := range entries {
		ids = append(ids, entry.ID)
	}
	return ids
}

func uuidPtr(value uuid.UUID) *uuid.UUID {
	return &value
}

func validBillingProfileInput(userID string) UpsertBillingProfileInput {
	return UpsertBillingProfileInput{
		UserID:       userID,
		EntityType:   BillingEntityIndividual,
		BillingName:  "Radu Boncea",
		BillingEmail: "radu@example.com",
		CountryCode:  "RO",
		AddressLine1: "Maresal Averescu 8-10",
		City:         "Bucuresti",
		PostalCode:   "011455",
	}
}

func validInvoiceInput(userID string, invoiceSerie string) CreateBillingInvoiceInput {
	return CreateBillingInvoiceInput{
		UserID:       userID,
		InvoiceSerie: invoiceSerie,
		InvoiceDate:  time.Date(2026, 6, 10, 0, 0, 0, 0, time.UTC),
		Lines: []CreateBillingInvoiceLineInput{
			{Name: "OCR credits", Quantity: 1, UnitPrice: "10.00", VATPercentage: "19.00"},
		},
	}
}

func openBillingRepositoryTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := testsupport.PostgresTestDSN(t)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open postgres: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get postgres sql db: %v", err)
	}
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})
	return db
}

func registerBillingProfileCreateBarrier(t *testing.T, db *gorm.DB, ready chan<- struct{}, release <-chan struct{}) {
	t.Helper()
	err := db.Callback().Create().Before("gorm:create").Register("billing_profile_create_barrier", func(tx *gorm.DB) {
		if tx.Statement == nil || tx.Statement.Schema == nil || tx.Statement.Schema.Table != (BillingProfile{}).TableName() {
			return
		}
		ready <- struct{}{}
		<-release
	})
	if err != nil {
		t.Fatalf("register billing profile create barrier: %v", err)
	}
}

func createCreditBucket(t *testing.T, db *gorm.DB, userID string, sourceType CreditSourceType, credits int, validFrom time.Time, expiresAt *time.Time) CreditBucket {
	t.Helper()
	bucket := CreditBucket{
		UserID:           userID,
		SourceType:       sourceType,
		CreditsGranted:   credits,
		CreditsRemaining: credits,
		ValidFrom:        validFrom,
		ExpiresAt:        expiresAt,
	}
	if err := db.Create(&bucket).Error; err != nil {
		t.Fatalf("create credit bucket: %v", err)
	}
	return bucket
}

func assertBucketRemaining(t *testing.T, db *gorm.DB, bucketID uuid.UUID, want int) {
	t.Helper()
	var bucket CreditBucket
	if err := db.First(&bucket, "id = ?", bucketID).Error; err != nil {
		t.Fatalf("load bucket %s: %v", bucketID, err)
	}
	if bucket.CreditsRemaining != want {
		t.Fatalf("bucket %s remaining = %d, want %d", bucketID, bucket.CreditsRemaining, want)
	}
}

func stringPtr(value string) *string {
	return &value
}
