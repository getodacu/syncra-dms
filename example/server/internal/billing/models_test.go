package billing

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"ai.ro/syncra/internal/auth"
	"ai.ro/syncra/internal/testsupport"
)

var billingModelTestGroup *testsupport.PostgresGroup

func TestBillingModels(t *testing.T) {
	billingModelTestGroup = testsupport.OpenPostgresGroup(t, billingTestModels()...)
	defer func() { billingModelTestGroup = nil }()

	for _, tt := range []struct {
		name string
		fn   func(*testing.T)
	}{
		{name: "AutoMigrateAndPersist", fn: testBillingModelsAutoMigrateAndPersist},
		{name: "Constraints", fn: testBillingModelConstraints},
		{name: "BeforeCreateAssignsIDAndValidatesEntityType", fn: testBillingProfileBeforeCreateAssignsIDAndValidatesEntityType},
		{name: "GORMUpdateValidation", fn: testBillingProfileGORMUpdateValidation},
		{name: "TableName", fn: testBillingProfileTableName},
		{name: "CascadeDeleteUser", fn: testBillingModelsCascadeDeleteUser},
	} {
		t.Run(tt.name, tt.fn)
	}
}

func billingModelTx(t *testing.T) *gorm.DB {
	t.Helper()
	if billingModelTestGroup != nil {
		return billingModelTestGroup.Tx(t)
	}
	return testsupport.OpenPostgresTx(t, billingTestModels()...)
}

func testBillingModelsAutoMigrateAndPersist(t *testing.T) {
	db := billingModelTx(t)
	now := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)
	user := createBillingTestUser(t, db)

	profile := validBillingProfile(user.ID)
	if err := db.Create(&profile).Error; err != nil {
		t.Fatalf("create billing profile: %v", err)
	}

	counter := BillingInvoiceCounter{InvoiceSerie: "SYNC"}
	if err := db.Create(&counter).Error; err != nil {
		t.Fatalf("create invoice counter: %v", err)
	}

	userID := user.ID
	profileID := profile.ID
	pdfPath := "/var/lib/syncra/invoices/invoice.pdf"
	invoice := BillingInvoice{
		UserID:                 &userID,
		BillingProfileID:       &profileID,
		BillingName:            profile.BillingName,
		BillingEmail:           profile.BillingEmail,
		BillingFiscalCode:      profile.FiscalCode,
		BillingProfileSnapshot: datatypes.JSON([]byte(`{"billing_name":"ICI Bucuresti"}`)),
		Lines:                  datatypes.JSON([]byte(`[{"name":"OCR credits","quantity":1,"unit_price":"10.00","vat_percentage":"19.00","total_vat_amount":"1.90","total_amount":"11.90"}]`)),
		NetAmount:              decimal.RequireFromString("10.00"),
		VATAmount:              decimal.RequireFromString("1.90"),
		TotalAmount:            decimal.RequireFromString("11.90"),
		InvoiceDate:            time.Date(2026, 6, 10, 0, 0, 0, 0, time.UTC),
		InvoiceSerie:           "SYNC",
		InvoiceNumber:          1,
		PDFPath:                &pdfPath,
	}
	if err := db.Create(&invoice).Error; err != nil {
		t.Fatalf("create invoice: %v", err)
	}

	order := BillingOrder{
		UserID: user.ID, OrderType: OrderTypeCreditTopup, Status: OrderStatusPending,
		Credits: 5000, AmountCents: 4750, Currency: "EUR",
		Provider: BillingProviderStripe, PricingTier: CreditPurchaseTier2,
		UnitAmountCents: 950,
	}
	if err := db.Create(&order).Error; err != nil {
		t.Fatalf("create order: %v", err)
	}

	bucket := CreditBucket{
		UserID: user.ID, SourceType: CreditSourceTopupPurchase, OrderID: &order.ID,
		CreditsGranted: 5000, CreditsRemaining: 5000,
		ValidFrom: now,
	}
	if err := db.Create(&bucket).Error; err != nil {
		t.Fatalf("create bucket: %v", err)
	}

	entry := CreditLedgerEntry{
		UserID: user.ID, BucketID: &bucket.ID, EntryType: CreditLedgerEntryDebit,
		CreditsDelta: -1, RelatedJobID: ptrUUID(uuid.New()),
		RelatedOrderID: &order.ID, IdempotencyKey: "test:debit:1",
		Metadata: datatypes.JSON([]byte(`{"page_count":1}`)),
	}
	if err := db.Create(&entry).Error; err != nil {
		t.Fatalf("create ledger entry: %v", err)
	}

	var got CreditBucket
	if err := db.First(&got, "id = ?", bucket.ID).Error; err != nil {
		t.Fatalf("load bucket: %v", err)
	}
	if got.UserID != user.ID || got.CreditsGranted != 5000 || got.CreditsRemaining != 5000 {
		t.Fatalf("unexpected bucket: %#v", got)
	}

	var gotProfile BillingProfile
	if err := db.First(&gotProfile, "id = ?", profile.ID).Error; err != nil {
		t.Fatalf("load billing profile: %v", err)
	}
	if gotProfile.UserID != user.ID || gotProfile.EntityType != BillingEntityCompany || gotProfile.BillingName != "ICI Bucuresti" {
		t.Fatalf("unexpected billing profile: %#v", gotProfile)
	}

	var gotInvoice BillingInvoice
	if err := db.First(&gotInvoice, "id = ?", invoice.ID).Error; err != nil {
		t.Fatalf("load invoice: %v", err)
	}
	if gotInvoice.InvoiceSerie != "SYNC" ||
		gotInvoice.InvoiceNumber != 1 ||
		!gotInvoice.NetAmount.Equal(decimal.RequireFromString("10.00")) ||
		!gotInvoice.VATAmount.Equal(decimal.RequireFromString("1.90")) ||
		!gotInvoice.TotalAmount.Equal(decimal.RequireFromString("11.90")) ||
		gotInvoice.PDFPath == nil ||
		*gotInvoice.PDFPath != pdfPath ||
		!strings.Contains(string(gotInvoice.BillingProfileSnapshot), "ICI Bucuresti") ||
		!strings.Contains(string(gotInvoice.Lines), "OCR credits") {
		t.Fatalf("unexpected invoice: %#v", gotInvoice)
	}
}

func testBillingModelConstraints(t *testing.T) {
	t.Run("duplicate signup bonus bucket fails unique constraint", func(t *testing.T) {
		db := billingModelTx(t)
		user := createBillingTestUser(t, db)
		now := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)

		first := CreditBucket{
			UserID: user.ID, SourceType: CreditSourceSignupBonus,
			CreditsGranted: 500, CreditsRemaining: 500,
			ValidFrom: now, ExpiresAt: ptrTime(now.AddDate(0, 0, 90)),
		}
		second := first
		if err := db.Create(&first).Error; err != nil {
			t.Fatalf("create first signup bucket: %v", err)
		}
		if err := db.Create(&second).Error; err == nil {
			t.Fatal("duplicate signup bonus bucket succeeded, want unique constraint failure")
		}
	})

	t.Run("duplicate ledger idempotency key fails unique constraint", func(t *testing.T) {
		db := billingModelTx(t)
		user := createBillingTestUser(t, db)
		now := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)

		bucket := CreditBucket{
			UserID: user.ID, SourceType: CreditSourceSignupBonus,
			CreditsGranted: 500, CreditsRemaining: 500,
			ValidFrom: now, ExpiresAt: ptrTime(now.AddDate(0, 0, 90)),
		}
		if err := db.Create(&bucket).Error; err != nil {
			t.Fatalf("create bucket: %v", err)
		}

		entry := CreditLedgerEntry{
			UserID: user.ID, BucketID: &bucket.ID, EntryType: CreditLedgerEntryGrant,
			CreditsDelta: 500, IdempotencyKey: "signup:" + user.ID,
		}
		if err := db.Create(&entry).Error; err != nil {
			t.Fatalf("create first ledger entry: %v", err)
		}
		duplicate := entry
		duplicate.ID = uuid.Nil
		if err := db.Create(&duplicate).Error; err == nil {
			t.Fatal("duplicate ledger idempotency key succeeded, want unique constraint failure")
		}
	})

	t.Run("monthly allowance source fails validation", func(t *testing.T) {
		db := billingModelTx(t)
		user := createBillingTestUser(t, db)
		now := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)

		monthlyAllowance := CreditBucket{
			UserID: user.ID, SourceType: CreditSourceType("monthly_allowance"),
			CreditsGranted: 500, CreditsRemaining: 500,
			ValidFrom: now,
		}
		if err := db.Create(&monthlyAllowance).Error; err == nil {
			t.Fatal("monthly allowance credit source succeeded, want validation failure")
		}
	})

	t.Run("duplicate billing profile per user fails unique constraint", func(t *testing.T) {
		db := billingModelTx(t)
		user := createBillingTestUser(t, db)

		first := validBillingProfile(user.ID)
		if err := db.Create(&first).Error; err != nil {
			t.Fatalf("create first billing profile: %v", err)
		}
		second := validBillingProfile(user.ID)
		second.BillingEmail = "billing-secondary@example.com"
		if err := db.Create(&second).Error; err == nil {
			t.Fatal("duplicate billing profile succeeded, want unique constraint failure")
		}
	})

	t.Run("duplicate invoice number per series fails unique constraint", func(t *testing.T) {
		db := billingModelTx(t)
		user := createBillingTestUser(t, db)
		profile := validBillingProfile(user.ID)
		if err := db.Create(&profile).Error; err != nil {
			t.Fatalf("create billing profile: %v", err)
		}
		first := validBillingInvoice(user.ID, profile.ID, 1, "SYNC")
		if err := db.Create(&first).Error; err != nil {
			t.Fatalf("create first invoice: %v", err)
		}
		otherSeries := validBillingInvoice(user.ID, profile.ID, 1, "ALT")
		if err := db.Create(&otherSeries).Error; err != nil {
			t.Fatalf("create same invoice number in another series: %v", err)
		}
		duplicate := validBillingInvoice(user.ID, profile.ID, 1, "SYNC")
		if err := db.Create(&duplicate).Error; err == nil {
			t.Fatal("duplicate invoice number in series succeeded, want unique constraint failure")
		}
	})
}

func testBillingProfileBeforeCreateAssignsIDAndValidatesEntityType(t *testing.T) {
	fiscalCode := "RO2785503"
	profile := BillingProfile{
		UserID:       uuid.NewString(),
		EntityType:   BillingEntityCompany,
		BillingName:  "ICI Bucuresti",
		BillingEmail: "billing@example.com",
		CountryCode:  "RO",
		AddressLine1: "Maresal Averescu 8-10",
		City:         "Bucuresti",
		PostalCode:   "011455",
		FiscalCode:   &fiscalCode,
	}

	if err := profile.BeforeCreate(nil); err != nil {
		t.Fatalf("BeforeCreate() error = %v", err)
	}
	if profile.ID == uuid.Nil {
		t.Fatal("BillingProfile.BeforeCreate() left nil ID")
	}

	invalid := profile
	invalid.ID = uuid.Nil
	invalid.EntityType = BillingEntityType("government")
	if err := invalid.BeforeCreate(nil); err == nil {
		t.Fatal("BeforeCreate() error = nil, want invalid entity type error")
	}
}

func testBillingProfileGORMUpdateValidation(t *testing.T) {
	db := billingModelTx(t)
	user := createBillingTestUser(t, db)
	profile := validBillingProfile(user.ID)
	if err := db.Create(&profile).Error; err != nil {
		t.Fatalf("create billing profile: %v", err)
	}

	if err := db.Model(&profile).Updates(BillingProfile{BillingName: "ICI Bucuresti Updated"}).Error; err != nil {
		t.Fatalf("update billing profile name: %v", err)
	}

	var got BillingProfile
	if err := db.First(&got, "id = ?", profile.ID).Error; err != nil {
		t.Fatalf("load updated billing profile: %v", err)
	}
	if got.BillingName != "ICI Bucuresti Updated" {
		t.Fatalf("BillingName = %q, want updated value", got.BillingName)
	}
	if got.EntityType != BillingEntityCompany {
		t.Fatalf("EntityType = %q, want %q", got.EntityType, BillingEntityCompany)
	}

	err := db.Model(&profile).Updates(map[string]any{"entity_type": "government"}).Error
	if err == nil {
		t.Fatal("invalid billing profile entity_type update succeeded, want validation failure")
	}
	if !strings.Contains(err.Error(), "invalid billing entity type") {
		t.Fatalf("invalid entity_type update error = %v, want billing entity validation error", err)
	}
}

func testBillingProfileTableName(t *testing.T) {
	if (BillingProfile{}).TableName() != "billing_profiles" {
		t.Fatalf("BillingProfile table = %q", (BillingProfile{}).TableName())
	}
	if (BillingInvoice{}).TableName() != "billing_invoices" {
		t.Fatalf("BillingInvoice table = %q", (BillingInvoice{}).TableName())
	}
	if (BillingInvoiceCounter{}).TableName() != "billing_invoice_counters" {
		t.Fatalf("BillingInvoiceCounter table = %q", (BillingInvoiceCounter{}).TableName())
	}
}

func testBillingModelsCascadeDeleteUser(t *testing.T) {
	db := billingModelTx(t)
	user := createBillingTestUser(t, db)
	now := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)

	profile := validBillingProfile(user.ID)
	if err := db.Create(&profile).Error; err != nil {
		t.Fatalf("create billing profile: %v", err)
	}

	invoice := validBillingInvoice(user.ID, profile.ID, 1, "SYNC")
	if err := db.Create(&invoice).Error; err != nil {
		t.Fatalf("create invoice: %v", err)
	}

	bucket := CreditBucket{
		UserID: user.ID, SourceType: CreditSourceSignupBonus,
		CreditsGranted: 500, CreditsRemaining: 500,
		ValidFrom: now, ExpiresAt: ptrTime(now.AddDate(0, 0, 90)),
	}
	if err := db.Create(&bucket).Error; err != nil {
		t.Fatalf("create bucket: %v", err)
	}
	if err := db.Create(&CreditLedgerEntry{
		UserID: user.ID, BucketID: &bucket.ID, EntryType: CreditLedgerEntryGrant,
		CreditsDelta: 500, IdempotencyKey: "signup:" + user.ID,
	}).Error; err != nil {
		t.Fatalf("create ledger entry: %v", err)
	}

	if err := db.Delete(&user).Error; err != nil {
		t.Fatalf("delete user: %v", err)
	}
	assertBillingCount(t, db, &CreditBucket{}, 0, "user_id = ?", user.ID)
	assertBillingCount(t, db, &CreditLedgerEntry{}, 0, "user_id = ?", user.ID)
	assertBillingCount(t, db, &BillingProfile{}, 0, "user_id = ?", user.ID)
	assertBillingCount(t, db, &BillingInvoice{}, 1, "id = ? AND user_id IS NULL AND billing_profile_id IS NULL", invoice.ID)
}

func billingTestModels() []any {
	return []any{
		&auth.User{},
		&BillingProfile{},
		&BillingInvoiceCounter{},
		&BillingInvoice{},
		&BillingOrder{},
		&CreditBucket{},
		&CreditLedgerEntry{},
	}
}

func createBillingTestUser(t *testing.T, db interface{ Create(value any) *gorm.DB }) auth.User {
	t.Helper()
	user := auth.User{Name: "Billing Owner", Email: uuid.NewString() + "@example.com"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	return user
}

func ptrTime(value time.Time) *time.Time {
	return &value
}

func ptrUUID(value uuid.UUID) *uuid.UUID {
	return &value
}

func validBillingProfile(userID string) BillingProfile {
	fiscalCode := "RO2785503"
	return BillingProfile{
		UserID:       userID,
		EntityType:   BillingEntityCompany,
		BillingName:  "ICI Bucuresti",
		BillingEmail: "billing@example.com",
		CountryCode:  "RO",
		AddressLine1: "Maresal Averescu 8-10",
		City:         "Bucuresti",
		PostalCode:   "011455",
		FiscalCode:   &fiscalCode,
	}
}

func validBillingInvoice(userID string, profileID uuid.UUID, invoiceNumber int64, invoiceSerie string) BillingInvoice {
	return BillingInvoice{
		UserID:                 &userID,
		BillingProfileID:       &profileID,
		BillingName:            "ICI Bucuresti",
		BillingEmail:           "billing@example.com",
		BillingFiscalCode:      stringPtr("RO2785503"),
		BillingProfileSnapshot: datatypes.JSON([]byte(`{"billing_name":"ICI Bucuresti"}`)),
		Lines:                  datatypes.JSON([]byte(`[{"name":"OCR credits","quantity":1,"unit_price":"10.00","vat_percentage":"19.00","total_vat_amount":"1.90","total_amount":"11.90"}]`)),
		NetAmount:              decimal.RequireFromString("10.00"),
		VATAmount:              decimal.RequireFromString("1.90"),
		TotalAmount:            decimal.RequireFromString("11.90"),
		InvoiceDate:            time.Date(2026, 6, 10, 0, 0, 0, 0, time.UTC),
		InvoiceSerie:           invoiceSerie,
		InvoiceNumber:          invoiceNumber,
	}
}

func assertBillingCount(t *testing.T, db *gorm.DB, model any, want int64, where string, args ...any) {
	t.Helper()
	var count int64
	if err := db.Model(model).Where(where, args...).Count(&count).Error; err != nil {
		t.Fatalf("count %T: %v", model, err)
	}
	if count != want {
		t.Fatalf("count %T = %d, want %d", model, count, want)
	}
}
