package database_test

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"ai.ro/syncra/internal/auth"
	"ai.ro/syncra/internal/billing"
	"ai.ro/syncra/internal/database"
	"ai.ro/syncra/internal/dbmigrate"
	"ai.ro/syncra/internal/ocr"
	"ai.ro/syncra/internal/testsupport"
	"ai.ro/syncra/internal/webhooks"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func migrateApplicationModelsForTest(db *gorm.DB) error {
	if err := dbmigrate.MigrateLegacyIntegerIDTables(db); err != nil {
		return err
	}
	if err := dbmigrate.MigrateOCRDocumentHash(db); err != nil {
		return err
	}
	if err := dbmigrate.MigrateCreditOnlyBilling(db); err != nil {
		return err
	}
	if err := dbmigrate.ValidateOCRJobStatuses(db); err != nil {
		return err
	}
	if err := db.AutoMigrate(database.ApplicationModels()...); err != nil {
		return err
	}
	if err := dbmigrate.MigrateOCRDocumentJobForeignKey(db); err != nil {
		return err
	}
	if err := dbmigrate.MigrateOCRDocumentLifecycle(db); err != nil {
		return err
	}
	if err := dbmigrate.MigrateOCRDocumentPageCount(db); err != nil {
		return err
	}
	if err := dbmigrate.MigrateOCRDocumentListIndexes(db); err != nil {
		return err
	}
	if err := dbmigrate.MigrateOCRJobStatus(db); err != nil {
		return err
	}
	return dbmigrate.MigrateOwnerForeignKeyCascades(db)
}

func migrationTestModels() []any {
	return []any{
		&billing.CreditLedgerEntry{},
		&billing.CreditBucket{},
		&billing.BillingOrder{},
		&billing.BillingInvoice{},
		&billing.BillingInvoiceCounter{},
		&ocr.Dataset{},
		&ocr.CollectionSchema{},
		&ocr.CollectionDocument{},
		&ocr.Collection{},
		&ocr.OCRDocument{},
		&ocr.OCRJob{},
		&ocr.ExtractionSchema{},
		&ocr.JSONRecipeCategory{},
		&ocr.JSONRecipe{},
		&auth.APIKey{},
		&auth.AuthAccount{},
		&auth.Session{},
		&auth.Verification{},
		&auth.User{},
	}
}

func TestApplicationModelsIncludesDataset(t *testing.T) {
	models := database.ApplicationModels()
	for _, model := range models {
		if _, ok := model.(*ocr.Dataset); ok {
			return
		}
	}
	t.Fatal("ApplicationModels missing *ocr.Dataset")
}

func TestApplicationModelsIncludesAPIKey(t *testing.T) {
	models := database.ApplicationModels()
	for _, model := range models {
		if _, ok := model.(*auth.APIKey); ok {
			return
		}
	}
	t.Fatal("ApplicationModels missing *auth.APIKey")
}

func TestApplicationModelsIncludesBillingInvoices(t *testing.T) {
	models := database.ApplicationModels()
	hasInvoice := false
	hasCounter := false
	for _, model := range models {
		switch model.(type) {
		case *billing.BillingInvoice:
			hasInvoice = true
		case *billing.BillingInvoiceCounter:
			hasCounter = true
		}
	}
	if !hasInvoice {
		t.Fatal("ApplicationModels missing *billing.BillingInvoice")
	}
	if !hasCounter {
		t.Fatal("ApplicationModels missing *billing.BillingInvoiceCounter")
	}
}

func TestApplicationModelsIncludesWebhook(t *testing.T) {
	models := database.ApplicationModels()
	for _, model := range models {
		if _, ok := model.(*webhooks.Webhook); ok {
			return
		}
	}
	t.Fatal("ApplicationModels missing *webhooks.Webhook")
}

func TestApplicationModelsIncludesJSONRecipe(t *testing.T) {
	models := database.ApplicationModels()
	for _, model := range models {
		if _, ok := model.(*ocr.JSONRecipe); ok {
			return
		}
	}
	t.Fatal("ApplicationModels missing *ocr.JSONRecipe")
}

func TestApplicationModelsIncludesJSONRecipeCategory(t *testing.T) {
	models := database.ApplicationModels()
	for _, model := range models {
		if _, ok := model.(*ocr.JSONRecipeCategory); ok {
			return
		}
	}
	t.Fatal("ApplicationModels missing *ocr.JSONRecipeCategory")
}

func TestMigrateCreatesApplicationTables(t *testing.T) {
	db := testsupport.OpenPostgresTx(t)
	modelsToMigrate := migrationTestModels()

	for _, model := range modelsToMigrate {
		if err := db.Migrator().DropTable(model); err != nil {
			t.Fatalf("drop table for %T: %v", model, err)
		}
	}

	if err := migrateApplicationModelsForTest(db); err != nil {
		t.Fatalf("Migrate() error = %v", err)
	}

	for _, model := range modelsToMigrate {
		if !db.Migrator().HasTable(model) {
			t.Fatalf("expected migrated table for %T", model)
		}
	}
}

func TestMigratePreservesLegacyOCRTables(t *testing.T) {
	db := testsupport.OpenPostgresTx(t)
	for _, model := range migrationTestModels() {
		if err := db.Migrator().DropTable(model); err != nil {
			t.Fatalf("drop table for %T: %v", model, err)
		}
	}

	if err := db.Exec(`
	CREATE TABLE "extraction_schemas" (
		"id" bigserial PRIMARY KEY,
		"created_at" timestamptz NOT NULL DEFAULT now(),
		"updated_at" timestamptz NOT NULL DEFAULT now(),
		"name" varchar(160) NOT NULL,
		"description" text NOT NULL DEFAULT '',
		"schema_json" jsonb NOT NULL,
		"strict" boolean NOT NULL DEFAULT true
	)
	`).Error; err != nil {
		t.Fatalf("create legacy extraction schemas table: %v", err)
	}
	if err := db.Exec(`
	CREATE TABLE "ocr_documents" (
		"id" bigserial PRIMARY KEY,
		"created_at" timestamptz NOT NULL DEFAULT now(),
		"updated_at" timestamptz NOT NULL DEFAULT now(),
		"original_filename" varchar(255) NOT NULL,
		"mime_type" varchar(120) NOT NULL,
		"file_size" bigint NOT NULL,
		"file_sha256" varchar(64) NOT NULL,
		"schema_id" bigint REFERENCES "extraction_schemas"("id") ON UPDATE CASCADE ON DELETE SET NULL,
		"inline_schema_json" jsonb,
		"markdown" text NOT NULL,
		"annotation_json" jsonb,
		"raw_response_json" jsonb NOT NULL,
		"status" varchar(40) NOT NULL DEFAULT 'completed',
		"error_message" text NOT NULL DEFAULT ''
	)
	`).Error; err != nil {
		t.Fatalf("create legacy OCR documents table: %v", err)
	}

	var legacySchemaID int64
	if err := db.Raw(`
	INSERT INTO "extraction_schemas" ("name", "schema_json", "strict")
	VALUES (?, ?::jsonb, ?)
	RETURNING "id"
	`, "legacy schema", `{"type":"object"}`, true).Scan(&legacySchemaID).Error; err != nil {
		t.Fatalf("insert legacy extraction schema: %v", err)
	}
	const legacyFileSHA = "fedcba9876543210fedcba9876543210fedcba9876543210fedcba9876543210"
	if err := db.Exec(`
	INSERT INTO "ocr_documents" (
		"original_filename",
		"mime_type",
		"file_size",
		"file_sha256",
		"schema_id",
		"markdown",
		"raw_response_json"
	) VALUES (?, ?, ?, ?, ?, ?, ?::jsonb)
	`, "legacy.pdf", "application/pdf", 42, legacyFileSHA, legacySchemaID, "# Legacy", `{"pages":[{},{}]}`).Error; err != nil {
		t.Fatalf("insert legacy OCR document: %v", err)
	}

	if err := migrateApplicationModelsForTest(db); err != nil {
		t.Fatalf("Migrate() error = %v", err)
	}

	for _, check := range []struct {
		table  string
		column string
	}{
		{table: "extraction_schemas", column: "id"},
		{table: "ocr_documents", column: "id"},
		{table: "ocr_documents", column: "schema_id"},
	} {
		dataType, exists, err := columnExists(t, db, check.table, check.column)
		if err != nil {
			t.Fatalf("query %s.%s type: %v", check.table, check.column, err)
		}
		if !exists || dataType != "uuid" {
			t.Fatalf("%s.%s type = %q, exists %v, want uuid", check.table, check.column, dataType, exists)
		}
	}
	if _, exists, err := columnExists(t, db, "ocr_documents", "file_sha256"); err != nil {
		t.Fatalf("query ocr_documents.file_sha256: %v", err)
	} else if exists {
		t.Fatal("ocr_documents.file_sha256 exists after migration")
	}
	if _, exists, err := columnExists(t, db, "ocr_documents", "status"); err != nil {
		t.Fatalf("query ocr_documents.status: %v", err)
	} else if exists {
		t.Fatal("ocr_documents.status exists after migration")
	}

	var migratedSchemaID string
	if err := db.Raw(`SELECT "id"::text FROM "extraction_schemas"`).Scan(&migratedSchemaID).Error; err != nil {
		t.Fatalf("select migrated extraction schema id: %v", err)
	}
	var migratedDoc struct {
		SchemaID     string `gorm:"column:schema_id"`
		DocumentHash string `gorm:"column:document_hash"`
		PageCount    int    `gorm:"column:page_count"`
	}
	if err := db.Raw(`
	SELECT "schema_id"::text AS schema_id,
		"document_hash",
		"page_count"
	FROM "ocr_documents"
	`).Scan(&migratedDoc).Error; err != nil {
		t.Fatalf("select migrated OCR document: %v", err)
	}
	if migratedDoc.SchemaID != migratedSchemaID {
		t.Fatalf("migrated document schema_id = %q, want %q", migratedDoc.SchemaID, migratedSchemaID)
	}
	if migratedDoc.DocumentHash != legacyFileSHA {
		t.Fatalf("migrated document_hash = %q, want %q", migratedDoc.DocumentHash, legacyFileSHA)
	}
	if migratedDoc.PageCount != 2 {
		t.Fatalf("migrated page_count = %d, want 2", migratedDoc.PageCount)
	}
	assertCount(t, db, &ocr.ExtractionSchema{}, "name = ?", "legacy schema", 1)
	assertCount(t, db, &ocr.OCRDocument{}, "original_filename = ?", "legacy.pdf", 1)
}

func TestMigrateCreatesOCRListIndexes(t *testing.T) {
	db := testsupport.OpenPostgresTx(t)
	modelsToMigrate := migrationTestModels()

	for _, model := range modelsToMigrate {
		if err := db.Migrator().DropTable(model); err != nil {
			t.Fatalf("drop table for %T: %v", model, err)
		}
	}

	if err := migrateApplicationModelsForTest(db); err != nil {
		t.Fatalf("Migrate() error = %v", err)
	}

	for _, indexName := range []string{
		"idx_ocr_jobs_deleted_at",
		"idx_ocr_jobs_user_created_id",
		"idx_ocr_jobs_user_status_created_id",
	} {
		if !db.Migrator().HasIndex(&ocr.OCRJob{}, indexName) {
			t.Fatalf("expected OCR job list index %q", indexName)
		}
	}
	for _, indexName := range []string{
		"idx_ocr_documents_user_created_id",
		"idx_ocr_documents_original_filename_trgm",
	} {
		if !db.Migrator().HasIndex(&ocr.OCRDocument{}, indexName) {
			t.Fatalf("expected OCR document list index %q", indexName)
		}
	}
}

func TestMigrateCreatesBillingIndexes(t *testing.T) {
	db := testsupport.OpenPostgresTx(t)
	modelsToMigrate := migrationTestModels()

	for _, model := range modelsToMigrate {
		if err := db.Migrator().DropTable(model); err != nil {
			t.Fatalf("drop table for %T: %v", model, err)
		}
	}

	if err := migrateApplicationModelsForTest(db); err != nil {
		t.Fatalf("Migrate() error = %v", err)
	}

	for _, check := range []struct {
		model any
		index string
	}{
		{model: &billing.BillingOrder{}, index: "idx_billing_orders_user_id"},
		{model: &billing.BillingOrder{}, index: "idx_billing_orders_provider_checkout_session_id"},
		{model: &billing.BillingOrder{}, index: "idx_billing_orders_provider_payment_intent_id"},
		{model: &billing.CreditBucket{}, index: "idx_credit_buckets_user_available"},
		{model: &billing.CreditBucket{}, index: "idx_credit_buckets_one_signup_bonus"},
		{model: &billing.CreditLedgerEntry{}, index: "idx_credit_ledger_entries_idempotency_key"},
		{model: &billing.CreditLedgerEntry{}, index: "idx_credit_ledger_entries_related_job_id"},
		{model: &billing.CreditLedgerEntry{}, index: "idx_credit_ledger_entries_transactions"},
	} {
		if !db.Migrator().HasIndex(check.model, check.index) {
			t.Fatalf("expected billing index %q", check.index)
		}
	}

	signupBonusIndex := indexDefinition(t, db, "idx_credit_buckets_one_signup_bonus")
	if !strings.Contains(signupBonusIndex, "UNIQUE INDEX") ||
		!strings.Contains(signupBonusIndex, "user_id") ||
		!strings.Contains(signupBonusIndex, "source_type") ||
		!strings.Contains(signupBonusIndex, string(billing.CreditSourceSignupBonus)) {
		t.Fatalf("signup bonus index definition = %q, want unique partial user_id index", signupBonusIndex)
	}
	availableIndex := indexDefinition(t, db, "idx_credit_buckets_user_available")
	if !strings.Contains(availableIndex, "user_id") ||
		!strings.Contains(availableIndex, "expires_at") ||
		!strings.Contains(availableIndex, "created_at") ||
		!strings.Contains(availableIndex, "credits_remaining > 0") ||
		!strings.Contains(availableIndex, "voided_at IS NULL") {
		t.Fatalf("available bucket index definition = %q, want partial user availability index", availableIndex)
	}
	checkoutSessionIndex := indexDefinition(t, db, "idx_billing_orders_provider_checkout_session_id")
	if !strings.Contains(checkoutSessionIndex, "UNIQUE INDEX") ||
		!strings.Contains(checkoutSessionIndex, "provider_checkout_session_id") ||
		!strings.Contains(checkoutSessionIndex, "provider_checkout_session_id IS NOT NULL") {
		t.Fatalf("checkout session index definition = %q, want unique partial provider index", checkoutSessionIndex)
	}
	paymentIntentIndex := indexDefinition(t, db, "idx_billing_orders_provider_payment_intent_id")
	if !strings.Contains(paymentIntentIndex, "UNIQUE INDEX") ||
		!strings.Contains(paymentIntentIndex, "provider_payment_intent_id") ||
		!strings.Contains(paymentIntentIndex, "provider_payment_intent_id IS NOT NULL") {
		t.Fatalf("payment intent index definition = %q, want unique partial provider index", paymentIntentIndex)
	}
	transactionsIndex := indexDefinition(t, db, "idx_credit_ledger_entries_transactions")
	if got, want := strings.Join(indexColumns(t, db, "idx_credit_ledger_entries_transactions"), ","), "user_id,created_at,id"; got != want {
		t.Fatalf("ledger transactions index columns = %q, want %q", got, want)
	}
	if !strings.Contains(transactionsIndex, "entry_type") ||
		!strings.Contains(transactionsIndex, string(billing.CreditLedgerEntryPurchase)) ||
		!strings.Contains(transactionsIndex, string(billing.CreditLedgerEntryDebit)) {
		t.Fatalf("ledger transactions index definition = %q, want purchase/debit partial index", transactionsIndex)
	}
}

func TestCreditOnlyBillingMigrationUpgradesPlanSchema(t *testing.T) {
	db := testsupport.OpenPostgresTx(t)

	for _, tableName := range []string{
		"credit_ledger_entries",
		"credit_buckets",
		"billing_orders",
		"billing_subscriptions",
		"billing_plans",
	} {
		if err := db.Exec(`DROP TABLE IF EXISTS "` + tableName + `" CASCADE`).Error; err != nil {
			t.Fatalf("drop %s: %v", tableName, err)
		}
	}
	if err := db.AutoMigrate(&auth.User{}); err != nil {
		t.Fatalf("auto migrate users: %v", err)
	}
	user := auth.User{Name: "Billing Migration Owner", Email: "billing-migration@example.com"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}

	if err := db.Exec(`
CREATE TABLE "billing_plans" (
  "id" varchar(40) PRIMARY KEY,
  "name" varchar(120) NOT NULL,
  "monthly_price_cents" bigint NOT NULL DEFAULT 0,
  "monthly_credits" bigint NOT NULL DEFAULT 0,
  "topup_block_credits" bigint NOT NULL DEFAULT 0,
  "topup_block_price_cents" bigint NOT NULL DEFAULT 0,
  "currency" varchar(3) NOT NULL,
  "active" boolean NOT NULL DEFAULT true
);
CREATE TABLE "billing_subscriptions" (
  "id" uuid PRIMARY KEY,
  "user_id" uuid NOT NULL,
  "plan_id" varchar(40) NOT NULL,
  "status" varchar(40) NOT NULL
);
CREATE TABLE "billing_orders" (
  "id" uuid PRIMARY KEY,
  "user_id" uuid NOT NULL,
  "order_type" varchar(40) NOT NULL,
  "status" varchar(40) NOT NULL,
  "plan_id_at_purchase" varchar(40) NOT NULL,
  "credit_blocks" bigint NOT NULL DEFAULT 0,
  "credits" bigint NOT NULL DEFAULT 0,
  "amount_cents" bigint NOT NULL DEFAULT 0,
  "currency" varchar(3) NOT NULL,
  "provider_checkout_session_id" varchar(255),
  "provider_payment_intent_id" varchar(255),
  "created_at" timestamptz,
  "updated_at" timestamptz,
  "paid_at" timestamptz,
  "failed_at" timestamptz,
  "refunded_at" timestamptz,
  "canceled_at" timestamptz,
  CONSTRAINT "chk_billing_orders_plan_id_at_purchase" CHECK ("plan_id_at_purchase" IN ('credit_only', 'starter', 'pro')),
  CONSTRAINT "chk_billing_orders_credit_blocks" CHECK ("credit_blocks" > 0),
  CONSTRAINT "chk_billing_orders_credits" CHECK ("credits" > 0),
  CONSTRAINT "chk_billing_orders_amount_cents" CHECK ("amount_cents" >= 0),
  CONSTRAINT "chk_billing_orders_order_type" CHECK ("order_type" IN ('credit_topup')),
  CONSTRAINT "chk_billing_orders_status" CHECK ("status" IN ('pending', 'paid', 'failed', 'refunded', 'canceled'))
);
CREATE INDEX "idx_billing_orders_plan_id_at_purchase" ON "billing_orders" ("plan_id_at_purchase");
CREATE TABLE "credit_buckets" (
  "id" uuid PRIMARY KEY,
  "user_id" uuid NOT NULL,
  "source_type" varchar(40) NOT NULL,
  "plan_id" varchar(40),
  "order_id" uuid,
  "credits_granted" bigint NOT NULL,
  "credits_remaining" bigint NOT NULL,
  "valid_from" timestamptz NOT NULL,
  "expires_at" timestamptz,
  "voided_at" timestamptz,
  "created_at" timestamptz,
  "updated_at" timestamptz,
  CONSTRAINT "chk_credit_buckets_source_type" CHECK ("source_type" IN ('signup_bonus', 'monthly_allowance', 'topup_purchase', 'refund', 'adjustment')),
  CONSTRAINT "chk_credit_buckets_credits_granted" CHECK ("credits_granted" > 0),
  CONSTRAINT "chk_credit_buckets_credits_remaining" CHECK ("credits_remaining" >= 0 AND "credits_remaining" <= "credits_granted")
);
CREATE INDEX "idx_credit_buckets_plan_id" ON "credit_buckets" ("plan_id");
ALTER TABLE "credit_buckets" ADD CONSTRAINT "fk_credit_buckets_plan" FOREIGN KEY ("plan_id") REFERENCES "billing_plans"("id") ON DELETE SET NULL ON UPDATE CASCADE;
`).Error; err != nil {
		t.Fatalf("create legacy billing schema: %v", err)
	}

	planID := "starter"
	orderID := uuid.New()
	proOrderID := uuid.New()
	freeOrderID := uuid.New()
	bucketID := uuid.New()
	if err := db.Exec(`
INSERT INTO "billing_plans" ("id", "name", "topup_block_credits", "topup_block_price_cents", "currency")
VALUES (?, ?, ?, ?, ?)
`, planID, "Starter", 1000, 1000, "EUR").Error; err != nil {
		t.Fatalf("insert legacy billing plan: %v", err)
	}
	if err := db.Exec(`
INSERT INTO "billing_plans" ("id", "name", "topup_block_credits", "topup_block_price_cents", "currency")
VALUES (?, ?, ?, ?, ?)
`, "pro", "Pro", 1000, 800, "EUR").Error; err != nil {
		t.Fatalf("insert legacy pro billing plan: %v", err)
	}
	if err := db.Exec(`
INSERT INTO "billing_subscriptions" ("id", "user_id", "plan_id", "status")
VALUES (?, ?, ?, ?)
`, uuid.New(), user.ID, planID, "active").Error; err != nil {
		t.Fatalf("insert legacy billing subscription: %v", err)
	}
	if err := db.Exec(`
INSERT INTO "billing_orders" ("id", "user_id", "order_type", "status", "plan_id_at_purchase", "credit_blocks", "credits", "amount_cents", "currency")
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
`, orderID, user.ID, billing.OrderTypeCreditTopup, billing.OrderStatusPending, planID, 5, 5000, 5000, "EUR").Error; err != nil {
		t.Fatalf("insert legacy billing order: %v", err)
	}
	if err := db.Exec(`
INSERT INTO "billing_orders" ("id", "user_id", "order_type", "status", "plan_id_at_purchase", "credit_blocks", "credits", "amount_cents", "currency")
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
`, proOrderID, user.ID, billing.OrderTypeCreditTopup, billing.OrderStatusPaid, "pro", 2, 2000, 1600, "EUR").Error; err != nil {
		t.Fatalf("insert legacy pro billing order: %v", err)
	}
	if err := db.Exec(`
INSERT INTO "billing_orders" ("id", "user_id", "order_type", "status", "plan_id_at_purchase", "credit_blocks", "credits", "amount_cents", "currency")
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
`, freeOrderID, user.ID, billing.OrderTypeCreditTopup, billing.OrderStatusPaid, "credit_only", 1, 1000, 0, "EUR").Error; err != nil {
		t.Fatalf("insert legacy free billing order: %v", err)
	}
	if err := db.Exec(`
INSERT INTO "credit_buckets" ("id", "user_id", "source_type", "plan_id", "credits_granted", "credits_remaining", "valid_from")
VALUES (?, ?, ?, ?, ?, ?, now())
`, bucketID, user.ID, "monthly_allowance", planID, 1200, 1200).Error; err != nil {
		t.Fatalf("insert legacy credit bucket: %v", err)
	}

	if err := db.Exec(readMigrationSQL(t, "20260604110000_credit_only_billing.sql")).Error; err != nil {
		t.Fatalf("apply credit-only billing migration: %v", err)
	}

	if tableExists(t, db, "billing_plans") {
		t.Fatal("billing_plans exists after credit-only migration")
	}
	if tableExists(t, db, "billing_subscriptions") {
		t.Fatal("billing_subscriptions exists after credit-only migration")
	}
	for _, check := range []struct {
		table  string
		column string
	}{
		{table: "billing_orders", column: "plan_id_at_purchase"},
		{table: "billing_orders", column: "credit_blocks"},
		{table: "credit_buckets", column: "plan_id"},
	} {
		if _, exists, err := columnExists(t, db, check.table, check.column); err != nil {
			t.Fatalf("query %s.%s: %v", check.table, check.column, err)
		} else if exists {
			t.Fatalf("%s.%s exists after credit-only migration", check.table, check.column)
		}
	}
	for _, check := range []struct {
		table  string
		column string
	}{
		{table: "billing_orders", column: "provider"},
		{table: "billing_orders", column: "pricing_tier"},
		{table: "billing_orders", column: "unit_amount_cents"},
	} {
		if _, exists, err := columnExists(t, db, check.table, check.column); err != nil {
			t.Fatalf("query %s.%s: %v", check.table, check.column, err)
		} else if !exists {
			t.Fatalf("%s.%s missing after credit-only migration", check.table, check.column)
		}
	}

	var migratedOrder billing.BillingOrder
	if err := db.First(&migratedOrder, "id = ?", orderID).Error; err != nil {
		t.Fatalf("load migrated order: %v", err)
	}
	if migratedOrder.Provider != billing.BillingProviderStripe ||
		migratedOrder.PricingTier != billing.CreditPurchaseTier2 ||
		migratedOrder.UnitAmountCents != 1000 {
		t.Fatalf("migrated order billing fields = provider %q tier %q unit %d",
			migratedOrder.Provider, migratedOrder.PricingTier, migratedOrder.UnitAmountCents)
	}
	var migratedProOrder billing.BillingOrder
	if err := db.First(&migratedProOrder, "id = ?", proOrderID).Error; err != nil {
		t.Fatalf("load migrated pro order: %v", err)
	}
	if migratedProOrder.Provider != billing.BillingProviderStripe ||
		migratedProOrder.PricingTier != billing.CreditPurchaseTier1 ||
		migratedProOrder.UnitAmountCents != 800 {
		t.Fatalf("migrated pro order billing fields = provider %q tier %q unit %d",
			migratedProOrder.Provider, migratedProOrder.PricingTier, migratedProOrder.UnitAmountCents)
	}
	var migratedFreeOrder billing.BillingOrder
	if err := db.First(&migratedFreeOrder, "id = ?", freeOrderID).Error; err != nil {
		t.Fatalf("load migrated free order: %v", err)
	}
	if migratedFreeOrder.Provider != billing.BillingProviderStripe ||
		migratedFreeOrder.PricingTier != billing.CreditPurchaseTier1 ||
		migratedFreeOrder.UnitAmountCents != 1000 {
		t.Fatalf("migrated free order billing fields = provider %q tier %q unit %d",
			migratedFreeOrder.Provider, migratedFreeOrder.PricingTier, migratedFreeOrder.UnitAmountCents)
	}

	var migratedBucket billing.CreditBucket
	if err := db.First(&migratedBucket, "id = ?", bucketID).Error; err != nil {
		t.Fatalf("load migrated bucket: %v", err)
	}
	if migratedBucket.SourceType != billing.CreditSourceAdjustment ||
		migratedBucket.CreditsRemaining != 1200 ||
		migratedBucket.VoidedAt != nil {
		t.Fatalf("migrated monthly bucket = source %q remaining %d voided %v",
			migratedBucket.SourceType, migratedBucket.CreditsRemaining, migratedBucket.VoidedAt)
	}

	newOrder := billing.BillingOrder{
		UserID:          user.ID,
		OrderType:       billing.OrderTypeCreditTopup,
		Status:          billing.OrderStatusPending,
		Provider:        billing.BillingProviderStripe,
		PricingTier:     billing.CreditPurchaseTier1,
		UnitAmountCents: 1000,
		Credits:         1000,
		AmountCents:     1000,
		Currency:        "EUR",
	}
	if err := db.Create(&newOrder).Error; err != nil {
		t.Fatalf("create credit-only billing order after migration: %v", err)
	}
}

func TestDatasetOwnerAndSchemaForeignKeys(t *testing.T) {
	db := testsupport.OpenPostgresTx(t, &auth.User{}, &ocr.ExtractionSchema{}, &ocr.Dataset{})
	user := auth.User{Name: "Dataset Owner", Email: "dataset-fk@example.com"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	schema := ocr.ExtractionSchema{
		UserID:     &user.ID,
		Name:       "Invoice",
		SchemaJSON: datatypes.JSON([]byte(`{"type":"object","properties":{"total":{"type":"number"}}}`)),
		Strict:     true,
	}
	if err := db.Create(&schema).Error; err != nil {
		t.Fatalf("create schema: %v", err)
	}
	dataset := ocr.Dataset{
		UserID:         user.ID,
		SchemaID:       schema.ID,
		Name:           "Invoices",
		SelectedFields: datatypes.JSON([]byte(`[{"path":"/total","key":"total","label":"Total"}]`)),
	}
	if err := db.Create(&dataset).Error; err != nil {
		t.Fatalf("create dataset: %v", err)
	}
	if dataset.ID == uuid.Nil {
		t.Fatal("dataset ID was not populated")
	}

	assertCount(t, db, &ocr.Dataset{}, "id = ?", dataset.ID, 1)
	if err := db.Delete(&schema).Error; err != nil {
		t.Fatalf("delete schema: %v", err)
	}
	assertCount(t, db, &ocr.Dataset{}, "id = ?", dataset.ID, 0)
}

func TestMigrateOwnerForeignKeysCascadeOnUserDelete(t *testing.T) {
	db := testsupport.OpenPostgresTx(t)
	modelsToMigrate := migrationTestModels()

	for _, model := range modelsToMigrate {
		if err := db.Migrator().DropTable(model); err != nil {
			t.Fatalf("drop table for %T: %v", model, err)
		}
	}

	if err := migrateApplicationModelsForTest(db); err != nil {
		t.Fatalf("Migrate() error = %v", err)
	}

	assertOwnerForeignKeyDeleteAction(t, db, "extraction_schemas", "c")
	assertOwnerForeignKeyDeleteAction(t, db, "ocr_documents", "c")
	assertOwnerForeignKeyDeleteAction(t, db, "ocr_jobs", "c")
	assertOwnerForeignKeyDeleteAction(t, db, "collections", "c")
	assertOwnerForeignKeyDeleteAction(t, db, "datasets", "c")
	assertOCRDocumentJobForeignKeyActions(t, db, "n", "c")

	user := auth.User{Name: "Owner", Email: "owner-cascade@example.com"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}

	ownedSchema := ocr.ExtractionSchema{
		UserID:     &user.ID,
		Name:       "owned schema",
		SchemaJSON: datatypes.JSON([]byte(`{"type":"object"}`)),
		Strict:     true,
	}
	if err := db.Create(&ownedSchema).Error; err != nil {
		t.Fatalf("create owned schema: %v", err)
	}
	systemSchema := ocr.ExtractionSchema{
		Name:       "system schema",
		SchemaJSON: datatypes.JSON([]byte(`{"type":"object"}`)),
		Strict:     true,
	}
	if err := db.Create(&systemSchema).Error; err != nil {
		t.Fatalf("create system schema: %v", err)
	}
	dataset := ocr.Dataset{
		UserID:         user.ID,
		SchemaID:       systemSchema.ID,
		Name:           "Owner cascade dataset",
		SelectedFields: datatypes.JSON([]byte(`[{"path":"/total","key":"total","label":"Total"}]`)),
	}
	if err := db.Create(&dataset).Error; err != nil {
		t.Fatalf("create dataset: %v", err)
	}

	ownedDoc := ocr.OCRDocument{
		UserID:           &user.ID,
		OriginalFilename: "owned.png",
		MimeType:         "image/png",
		FileSize:         8,
		DocumentHash:     "owned-hash",
		Markdown:         "# Owned",
		RawResponseJSON:  datatypes.JSON([]byte(`{"pages":[]}`)),
	}
	if err := db.Create(&ownedDoc).Error; err != nil {
		t.Fatalf("create owned OCR document: %v", err)
	}
	collection := ocr.Collection{
		UserID: user.ID,
		Name:   "Invoices",
	}
	if err := db.Create(&collection).Error; err != nil {
		t.Fatalf("create collection: %v", err)
	}
	collectionSchema := ocr.CollectionSchema{
		CollectionID: collection.ID,
		SchemaID:     ownedSchema.ID,
	}
	if err := db.Create(&collectionSchema).Error; err != nil {
		t.Fatalf("create collection schema: %v", err)
	}
	collectionDocument := ocr.CollectionDocument{
		CollectionID: collection.ID,
		DocumentID:   ownedDoc.ID,
	}
	if err := db.Create(&collectionDocument).Error; err != nil {
		t.Fatalf("create collection document: %v", err)
	}
	systemDoc := ocr.OCRDocument{
		OriginalFilename: "system.png",
		MimeType:         "image/png",
		FileSize:         8,
		DocumentHash:     "system-hash",
		Markdown:         "# System",
		RawResponseJSON:  datatypes.JSON([]byte(`{"pages":[]}`)),
	}
	if err := db.Create(&systemDoc).Error; err != nil {
		t.Fatalf("create system OCR document: %v", err)
	}

	ownedJob := ocr.OCRJob{
		UserID:           &user.ID,
		OriginalFilename: "owned-job.png",
		MimeType:         "image/png",
		FileSize:         8,
		PageCount:        1,
		DocumentHash:     "owned-job-hash",
		FilePath:         "/tmp/owned-job.png",
		Status:           ocr.OCRJobStatusQueued,
	}
	if err := db.Create(&ownedJob).Error; err != nil {
		t.Fatalf("create owned OCR job: %v", err)
	}
	systemJob := ocr.OCRJob{
		OriginalFilename: "system-job.png",
		MimeType:         "image/png",
		FileSize:         8,
		PageCount:        1,
		DocumentHash:     "system-job-hash",
		FilePath:         "/tmp/system-job.png",
		Status:           ocr.OCRJobStatusQueued,
	}
	if err := db.Create(&systemJob).Error; err != nil {
		t.Fatalf("create system OCR job: %v", err)
	}

	if err := db.Delete(&user).Error; err != nil {
		t.Fatalf("delete user: %v", err)
	}

	assertCount(t, db, &ocr.ExtractionSchema{}, "id = ?", ownedSchema.ID, 0)
	assertCount(t, db, &ocr.OCRDocument{}, "id = ?", ownedDoc.ID, 0)
	assertCount(t, db, &ocr.ExtractionSchema{}, "id = ? AND user_id IS NULL", ownedSchema.ID, 0)
	assertCount(t, db, &ocr.OCRDocument{}, "id = ? AND user_id IS NULL", ownedDoc.ID, 0)
	assertCount(t, db, &ocr.Collection{}, "id = ?", collection.ID, 0)
	assertCount(t, db, &ocr.CollectionSchema{}, "collection_id = ?", collection.ID, 0)
	assertCount(t, db, &ocr.CollectionDocument{}, "collection_id = ?", collection.ID, 0)
	assertCount(t, db, &ocr.Dataset{}, "id = ?", dataset.ID, 0)
	assertCount(t, db, &ocr.OCRJob{}, "id = ?", ownedJob.ID, 0)
	assertCount(t, db, &ocr.OCRJob{}, "id = ? AND user_id IS NULL", ownedJob.ID, 0)
	assertCount(t, db, &ocr.ExtractionSchema{}, "id = ? AND user_id IS NULL", systemSchema.ID, 1)
	assertCount(t, db, &ocr.OCRDocument{}, "id = ? AND user_id IS NULL", systemDoc.ID, 1)
	assertCount(t, db, &ocr.OCRJob{}, "id = ? AND user_id IS NULL", systemJob.ID, 1)
}

func TestMigrateOCRJobStatusConstraints(t *testing.T) {
	db := testsupport.OpenPostgresTx(t)
	modelsToMigrate := migrationTestModels()

	for _, model := range modelsToMigrate {
		if err := db.Migrator().DropTable(model); err != nil {
			t.Fatalf("drop table for %T: %v", model, err)
		}
	}

	if err := migrateApplicationModelsForTest(db); err != nil {
		t.Fatalf("Migrate() error = %v", err)
	}

	var defaultExpression string
	if err := db.Raw(`
SELECT column_default
FROM information_schema.columns
WHERE table_schema = current_schema()
	AND table_name = 'ocr_jobs'
	AND column_name = 'status'
`).Scan(&defaultExpression).Error; err != nil {
		t.Fatalf("query OCR job status default: %v", err)
	}
	if !strings.Contains(defaultExpression, string(ocr.OCRJobStatusQueued)) {
		t.Fatalf("ocr_jobs.status default = %q, want queued", defaultExpression)
	}

	defaultJobID := uuid.New()
	if err := db.Exec(`
INSERT INTO ocr_jobs (id, original_filename, mime_type, file_size, page_count, document_hash, file_path)
VALUES (?, ?, ?, ?, ?, ?, ?)
`,
		defaultJobID,
		"default-job.png",
		"image/png",
		int64(8),
		1,
		"default-job-hash",
		"/tmp/default-job.png",
	).Error; err != nil {
		t.Fatalf("insert OCR job without status: %v", err)
	}

	var defaultStatus ocr.OCRJobStatus
	if err := db.Raw(`SELECT status FROM ocr_jobs WHERE id = ?`, defaultJobID).Scan(&defaultStatus).Error; err != nil {
		t.Fatalf("query default OCR job status: %v", err)
	}
	if defaultStatus != ocr.OCRJobStatusQueued {
		t.Fatalf("default OCR job status = %q, want %q", defaultStatus, ocr.OCRJobStatusQueued)
	}

	if err := db.Exec(`SAVEPOINT ocr_job_invalid_status`).Error; err != nil {
		t.Fatalf("create invalid status savepoint: %v", err)
	}
	invalidErr := db.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)}).Exec(`
INSERT INTO ocr_jobs (id, original_filename, mime_type, file_size, page_count, document_hash, file_path, status)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
`,
		uuid.New(),
		"invalid-job.png",
		"image/png",
		int64(8),
		1,
		"invalid-job-hash",
		"/tmp/invalid-job.png",
		"done",
	).Error
	if err := db.Exec(`ROLLBACK TO SAVEPOINT ocr_job_invalid_status`).Error; err != nil {
		t.Fatalf("rollback invalid status savepoint: %v", err)
	}
	if err := db.Exec(`RELEASE SAVEPOINT ocr_job_invalid_status`).Error; err != nil {
		t.Fatalf("release invalid status savepoint: %v", err)
	}
	if invalidErr == nil {
		t.Fatal("insert OCR job with invalid status succeeded, want error")
	}
}

func TestMigrateOCRJobStatusIsIdempotent(t *testing.T) {
	db := testsupport.OpenPostgresTx(t)
	modelsToMigrate := migrationTestModels()

	for _, model := range modelsToMigrate {
		if err := db.Migrator().DropTable(model); err != nil {
			t.Fatalf("drop table for %T: %v", model, err)
		}
	}

	if err := migrateApplicationModelsForTest(db); err != nil {
		t.Fatalf("first Migrate() error = %v", err)
	}
	beforeConstraint := ocrJobStatusConstraint(t, db)
	beforeDefault := ocrJobStatusDefault(t, db)
	if !strings.Contains(beforeDefault, string(ocr.OCRJobStatusQueued)) {
		t.Fatalf("ocr_jobs.status default = %q, want queued", beforeDefault)
	}

	if err := migrateApplicationModelsForTest(db); err != nil {
		t.Fatalf("second Migrate() error = %v", err)
	}
	afterConstraint := ocrJobStatusConstraint(t, db)
	afterDefault := ocrJobStatusDefault(t, db)

	if afterConstraint.OID != beforeConstraint.OID {
		t.Fatalf("chk_ocr_jobs_status oid changed after second migrate: before %s, after %s", beforeConstraint.OID, afterConstraint.OID)
	}
	if afterConstraint.Definition != beforeConstraint.Definition {
		t.Fatalf("chk_ocr_jobs_status definition changed after second migrate:\nbefore: %s\nafter:  %s", beforeConstraint.Definition, afterConstraint.Definition)
	}
	if afterDefault != beforeDefault {
		t.Fatalf("ocr_jobs.status default changed after second migrate: before %q, after %q", beforeDefault, afterDefault)
	}
}

func TestMigrateOCRJobStatusRejectsExistingInvalidRows(t *testing.T) {
	db := testsupport.OpenPostgresTx(t)
	if err := db.Migrator().DropTable(&ocr.OCRJob{}); err != nil {
		t.Fatalf("drop OCR jobs table: %v", err)
	}
	if err := db.AutoMigrate(&ocr.OCRJob{}); err != nil {
		t.Fatalf("auto migrate OCR jobs table: %v", err)
	}

	if err := db.Exec(`
INSERT INTO ocr_jobs (id, original_filename, mime_type, file_size, page_count, document_hash, file_path, status)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
`,
		uuid.New(),
		"invalid-existing-job.png",
		"image/png",
		int64(8),
		1,
		"invalid-existing-job-hash",
		"/tmp/invalid-existing-job.png",
		"done",
	).Error; err != nil {
		t.Fatalf("insert existing invalid OCR job status: %v", err)
	}

	err := dbmigrate.MigrateOCRJobStatus(db)
	if err == nil {
		t.Fatal("MigrateOCRJobStatus() error = nil, want invalid rows error")
	}
	if !strings.Contains(err.Error(), "invalid OCR job status rows") {
		t.Fatalf("MigrateOCRJobStatus() error = %q, want invalid OCR job status rows", err.Error())
	}
}

func TestMigrateOCRJobStatusSetsMissingDefault(t *testing.T) {
	db := testsupport.OpenPostgresTx(t)
	prepareOCRJobsTableWithoutStatusDefault(t, db)

	if err := dbmigrate.MigrateOCRJobStatus(db); err != nil {
		t.Fatalf("MigrateOCRJobStatus() error = %v", err)
	}

	defaultExpression := ocrJobStatusDefault(t, db)
	if !strings.Contains(defaultExpression, string(ocr.OCRJobStatusQueued)) {
		t.Fatalf("ocr_jobs.status default = %q, want queued", defaultExpression)
	}
}

func TestMigrateOCRJobStatusRejectsNullAfterNullableLegacyColumn(t *testing.T) {
	db := testsupport.OpenPostgresTx(t)
	prepareOCRJobsTableWithNullableStatus(t, db)

	if err := dbmigrate.MigrateOCRJobStatus(db); err != nil {
		t.Fatalf("MigrateOCRJobStatus() error = %v", err)
	}

	if ocrJobStatusNullable(t, db) {
		t.Fatal("ocr_jobs.status is nullable after migration, want NOT NULL")
	}

	if err := db.Exec(`SAVEPOINT ocr_job_null_status`).Error; err != nil {
		t.Fatalf("create null status savepoint: %v", err)
	}
	nullErr := db.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)}).Exec(`
INSERT INTO ocr_jobs (id, original_filename, mime_type, file_size, page_count, document_hash, file_path, status)
VALUES (?, ?, ?, ?, ?, ?, ?, NULL)
`,
		uuid.New(),
		"null-status-job.png",
		"image/png",
		int64(8),
		1,
		"null-status-job-hash",
		"/tmp/null-status-job.png",
	).Error
	if err := db.Exec(`ROLLBACK TO SAVEPOINT ocr_job_null_status`).Error; err != nil {
		t.Fatalf("rollback null status savepoint: %v", err)
	}
	if err := db.Exec(`RELEASE SAVEPOINT ocr_job_null_status`).Error; err != nil {
		t.Fatalf("release null status savepoint: %v", err)
	}
	if nullErr == nil {
		t.Fatal("raw SQL insert with NULL OCR job status succeeded, want database error")
	}
}

func TestMigrateOCRJobStatusRejectsInvalidRowsWithMissingDefault(t *testing.T) {
	db := testsupport.OpenPostgresTx(t)
	prepareOCRJobsTableWithoutStatusDefault(t, db)

	if err := db.Exec(`
INSERT INTO ocr_jobs (id, original_filename, mime_type, file_size, page_count, document_hash, file_path, status)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
`,
		uuid.New(),
		"invalid-existing-job.png",
		"image/png",
		int64(8),
		1,
		"invalid-existing-job-hash",
		"/tmp/invalid-existing-job.png",
		"done",
	).Error; err != nil {
		t.Fatalf("insert existing invalid OCR job status: %v", err)
	}

	err := dbmigrate.MigrateOCRJobStatus(db)
	if err == nil {
		t.Fatal("MigrateOCRJobStatus() error = nil, want invalid rows error")
	}
	if !strings.Contains(err.Error(), "invalid OCR job status rows") {
		t.Fatalf("MigrateOCRJobStatus() error = %q, want invalid OCR job status rows", err.Error())
	}
}

func TestOCRJobStatusMigrationsSkipTableWithoutStatusColumn(t *testing.T) {
	db := testsupport.OpenPostgresTx(t)
	if err := db.Migrator().DropTable(&ocr.OCRJob{}); err != nil {
		t.Fatalf("drop OCR jobs table: %v", err)
	}
	if err := db.Exec(`
CREATE TABLE ocr_jobs (
	id uuid PRIMARY KEY,
	original_filename varchar(255) NOT NULL
)
`).Error; err != nil {
		t.Fatalf("create OCR jobs table without status: %v", err)
	}

	if err := dbmigrate.ValidateOCRJobStatuses(db); err != nil {
		t.Fatalf("ValidateOCRJobStatuses() error = %v", err)
	}
	if err := dbmigrate.MigrateOCRJobStatus(db); err != nil {
		t.Fatalf("MigrateOCRJobStatus() error = %v", err)
	}
	if _, statusExists, err := columnExists(t, db, "ocr_jobs", "status"); err != nil {
		t.Fatalf("query OCR jobs status column: %v", err)
	} else if statusExists {
		t.Fatal("ocr_jobs.status exists after status-only migrations, want skipped missing column")
	}
}

func TestMigratePreflightsInvalidOCRJobStatusBeforeAutoMigrateDDL(t *testing.T) {
	db := testsupport.OpenPostgresTx(t)
	prepareOCRJobsTableWithoutStatusDefault(t, db)

	if err := db.Exec(`
INSERT INTO ocr_jobs (id, original_filename, mime_type, file_size, page_count, document_hash, file_path, status)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
`,
		uuid.New(),
		"invalid-top-level-job.png",
		"image/png",
		int64(8),
		1,
		"invalid-top-level-job-hash",
		"/tmp/invalid-top-level-job.png",
		"done",
	).Error; err != nil {
		t.Fatalf("insert existing invalid OCR job status: %v", err)
	}

	err := migrateApplicationModelsForTest(db)
	if err == nil {
		t.Fatal("Migrate() error = nil, want invalid rows error")
	}
	if !strings.Contains(err.Error(), "invalid OCR job status rows") {
		t.Fatalf("Migrate() error = %q, want invalid OCR job status rows", err.Error())
	}
	if ocrJobStatusHasDefault(t, db) {
		t.Fatal("ocr_jobs.status has a default after failed Migrate(), want no partial status default DDL")
	}
	if ocrJobStatusConstraintExists(t, db) {
		t.Fatal("chk_ocr_jobs_status exists after failed Migrate(), want no check constraint DDL")
	}
}

func assertOwnerForeignKeyDeleteAction(t *testing.T, db *gorm.DB, tableName string, want string) {
	t.Helper()
	var actions []string
	if err := db.Raw(`
SELECT con.confdeltype
FROM pg_constraint con
JOIN pg_class rel ON rel.oid = con.conrelid
JOIN pg_namespace nsp ON nsp.oid = rel.relnamespace
JOIN pg_attribute att ON att.attrelid = con.conrelid AND att.attnum = ANY(con.conkey)
WHERE con.contype = 'f'
	AND nsp.nspname = current_schema()
	AND rel.relname = ?
	AND att.attname = 'user_id'
	AND con.confrelid = to_regclass(format('%I.%I', current_schema(), 'user'))
`, tableName).Scan(&actions).Error; err != nil {
		t.Fatalf("query %s user_id foreign key: %v", tableName, err)
	}
	if len(actions) != 1 {
		t.Fatalf("%s user_id foreign key count = %d, want 1", tableName, len(actions))
	}
	if actions[0] != want {
		t.Fatalf("%s user_id foreign key delete action = %q, want %q", tableName, actions[0], want)
	}
}

func assertOCRDocumentJobForeignKeyActions(t *testing.T, db *gorm.DB, wantDelete string, wantUpdate string) {
	t.Helper()
	var actions []struct {
		DeleteAction string `gorm:"column:delete_action"`
		UpdateAction string `gorm:"column:update_action"`
	}
	if err := db.Raw(`
SELECT con.confdeltype AS delete_action, con.confupdtype AS update_action
FROM pg_constraint con
JOIN pg_class rel ON rel.oid = con.conrelid
JOIN pg_namespace nsp ON nsp.oid = rel.relnamespace
JOIN pg_attribute att ON att.attrelid = con.conrelid AND att.attnum = ANY(con.conkey)
WHERE con.contype = 'f'
	AND nsp.nspname = current_schema()
	AND rel.relname = 'ocr_documents'
	AND att.attname = 'job_id'
	AND con.confrelid = to_regclass(format('%I.%I', current_schema(), 'ocr_jobs'))
`).Scan(&actions).Error; err != nil {
		t.Fatalf("query ocr_documents job_id foreign key: %v", err)
	}
	if len(actions) != 1 {
		t.Fatalf("ocr_documents job_id foreign key count = %d, want 1", len(actions))
	}
	if actions[0].DeleteAction != wantDelete {
		t.Fatalf("ocr_documents job_id foreign key delete action = %q, want %q", actions[0].DeleteAction, wantDelete)
	}
	if actions[0].UpdateAction != wantUpdate {
		t.Fatalf("ocr_documents job_id foreign key update action = %q, want %q", actions[0].UpdateAction, wantUpdate)
	}
}

func assertCount(t *testing.T, db *gorm.DB, model any, query string, arg any, want int64) {
	t.Helper()
	var count int64
	if err := db.Model(model).Where(query, arg).Count(&count).Error; err != nil {
		t.Fatalf("count %T: %v", model, err)
	}
	if count != want {
		t.Fatalf("count %T where %q = %d, want %d", model, query, count, want)
	}
}

type statusConstraint struct {
	OID        string `gorm:"column:oid"`
	Definition string `gorm:"column:definition"`
}

func ocrJobStatusConstraint(t *testing.T, db *gorm.DB) statusConstraint {
	t.Helper()

	var constraints []statusConstraint
	if err := db.Raw(`
SELECT con.oid::text AS oid, pg_get_constraintdef(con.oid, true) AS definition
FROM pg_constraint con
JOIN pg_class rel ON rel.oid = con.conrelid
JOIN pg_namespace nsp ON nsp.oid = rel.relnamespace
WHERE con.contype = 'c'
	AND nsp.nspname = current_schema()
	AND rel.relname = 'ocr_jobs'
	AND con.conname = 'chk_ocr_jobs_status'
`).Scan(&constraints).Error; err != nil {
		t.Fatalf("query OCR job status check constraint: %v", err)
	}
	if len(constraints) != 1 {
		t.Fatalf("chk_ocr_jobs_status constraint count = %d, want 1", len(constraints))
	}
	return constraints[0]
}

func ocrJobStatusDefault(t *testing.T, db *gorm.DB) string {
	t.Helper()

	var defaultExpression string
	if err := db.Raw(`
SELECT column_default
FROM information_schema.columns
WHERE table_schema = current_schema()
	AND table_name = 'ocr_jobs'
	AND column_name = 'status'
`).Scan(&defaultExpression).Error; err != nil {
		t.Fatalf("query OCR job status default: %v", err)
	}
	return defaultExpression
}

func ocrJobStatusHasDefault(t *testing.T, db *gorm.DB) bool {
	t.Helper()

	var count int64
	if err := db.Raw(`
SELECT COUNT(*)
FROM information_schema.columns
WHERE table_schema = current_schema()
	AND table_name = 'ocr_jobs'
	AND column_name = 'status'
	AND column_default IS NOT NULL
`).Scan(&count).Error; err != nil {
		t.Fatalf("query OCR job status default presence: %v", err)
	}
	return count == 1
}

func ocrJobStatusConstraintExists(t *testing.T, db *gorm.DB) bool {
	t.Helper()

	var count int64
	if err := db.Raw(`
SELECT COUNT(*)
FROM pg_constraint con
JOIN pg_class rel ON rel.oid = con.conrelid
JOIN pg_namespace nsp ON nsp.oid = rel.relnamespace
WHERE con.contype = 'c'
	AND nsp.nspname = current_schema()
	AND rel.relname = 'ocr_jobs'
	AND con.conname = 'chk_ocr_jobs_status'
`).Scan(&count).Error; err != nil {
		t.Fatalf("query OCR job status check constraint presence: %v", err)
	}
	return count == 1
}

func ocrJobStatusNullable(t *testing.T, db *gorm.DB) bool {
	t.Helper()

	var isNullable string
	if err := db.Raw(`
SELECT is_nullable
FROM information_schema.columns
WHERE table_schema = current_schema()
	AND table_name = 'ocr_jobs'
	AND column_name = 'status'
`).Scan(&isNullable).Error; err != nil {
		t.Fatalf("query OCR job status nullability: %v", err)
	}
	return isNullable == "YES"
}

func indexDefinition(t *testing.T, db *gorm.DB, indexName string) string {
	t.Helper()

	var definition string
	if err := db.Raw(`
SELECT indexdef
FROM pg_indexes
WHERE schemaname = current_schema()
	AND indexname = ?
`, indexName).Scan(&definition).Error; err != nil {
		t.Fatalf("query index %s definition: %v", indexName, err)
	}
	return definition
}

func indexColumns(t *testing.T, db *gorm.DB, indexName string) []string {
	t.Helper()

	var columns []string
	if err := db.Raw(`
SELECT a.attname
FROM pg_class i
JOIN pg_index ix ON ix.indexrelid = i.oid
JOIN LATERAL unnest(ix.indkey) WITH ORDINALITY AS k(attnum, ord) ON true
JOIN pg_attribute a ON a.attrelid = ix.indrelid
	AND a.attnum = k.attnum
WHERE i.relname = ?
ORDER BY k.ord
`, indexName).Scan(&columns).Error; err != nil {
		t.Fatalf("query index %s columns: %v", indexName, err)
	}
	return columns
}

func tableExists(t *testing.T, db *gorm.DB, tableName string) bool {
	t.Helper()

	var exists bool
	if err := db.Raw(`
SELECT EXISTS (
	SELECT 1
	FROM information_schema.tables
	WHERE table_schema = current_schema()
		AND table_name = ?
)
`, tableName).Scan(&exists).Error; err != nil {
		t.Fatalf("query table %s existence: %v", tableName, err)
	}
	return exists
}

func columnExists(t *testing.T, db *gorm.DB, tableName string, columnName string) (string, bool, error) {
	t.Helper()

	var dataType string
	err := db.Raw(`
SELECT data_type
FROM information_schema.columns
WHERE table_schema = current_schema()
	AND table_name = ?
	AND column_name = ?
`, tableName, columnName).Scan(&dataType).Error
	return dataType, dataType != "", err
}

func readMigrationSQL(t *testing.T, filename string) string {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot resolve database test path")
	}
	path := filepath.Join(filepath.Dir(file), "..", "..", "migrations", filename)
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read migration %s: %v", filename, err)
	}
	return string(raw)
}

func prepareOCRJobsTableWithoutStatusDefault(t *testing.T, db *gorm.DB) {
	t.Helper()

	if err := db.Migrator().DropTable(&ocr.OCRJob{}); err != nil {
		t.Fatalf("drop OCR jobs table: %v", err)
	}
	if err := db.AutoMigrate(&ocr.OCRJob{}); err != nil {
		t.Fatalf("auto migrate OCR jobs table: %v", err)
	}
	if err := db.Exec(`ALTER TABLE "ocr_jobs" ALTER COLUMN "status" DROP DEFAULT`).Error; err != nil {
		t.Fatalf("drop OCR jobs status default: %v", err)
	}
}

func prepareOCRJobsTableWithNullableStatus(t *testing.T, db *gorm.DB) {
	t.Helper()

	if err := db.Migrator().DropTable(&ocr.OCRJob{}); err != nil {
		t.Fatalf("drop OCR jobs table: %v", err)
	}
	if err := db.AutoMigrate(&ocr.OCRJob{}); err != nil {
		t.Fatalf("auto migrate OCR jobs table: %v", err)
	}
	if err := db.Exec(`ALTER TABLE "ocr_jobs" ALTER COLUMN "status" DROP NOT NULL`).Error; err != nil {
		t.Fatalf("drop OCR jobs status not-null: %v", err)
	}
}
