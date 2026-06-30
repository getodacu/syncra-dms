package dbmigrate_test

import (
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"ai.ro/syncra/internal/auth"
	"ai.ro/syncra/internal/billing"
	"ai.ro/syncra/internal/dbmigrate"
	"ai.ro/syncra/internal/testsupport"
)

func TestMigrateCreditOnlyBillingUpgradesLegacyPlanSchemaBeforeAutoMigrate(t *testing.T) {
	db := testsupport.OpenPostgresTx(t)
	resetBillingTables(t, db)
	if err := db.AutoMigrate(&auth.User{}); err != nil {
		t.Fatalf("auto migrate users: %v", err)
	}
	user := auth.User{Name: "Credit Only Migration", Email: "credit-only-" + uuid.NewString() + "@example.com"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	createLegacyBillingPlanSchema(t, db)

	orderID := uuid.New()
	bucketID := uuid.New()
	if err := db.Exec(`
INSERT INTO "billing_plans" ("id", "name", "currency")
VALUES (?, ?, ?)
`, "starter", "Starter", "EUR").Error; err != nil {
		t.Fatalf("insert legacy billing plan: %v", err)
	}
	if err := db.Exec(`
INSERT INTO "billing_subscriptions" ("id", "user_id", "plan_id", "status")
VALUES (?, ?, ?, ?)
`, uuid.New(), user.ID, "starter", "active").Error; err != nil {
		t.Fatalf("insert legacy billing subscription: %v", err)
	}
	if err := db.Exec(`
INSERT INTO "billing_orders" ("id", "user_id", "order_type", "status", "plan_id_at_purchase", "credit_blocks", "credits", "amount_cents", "currency")
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
`, orderID, user.ID, billing.OrderTypeCreditTopup, billing.OrderStatusPending, "starter", 5, 5000, 5000, "EUR").Error; err != nil {
		t.Fatalf("insert legacy billing order: %v", err)
	}
	if err := db.Exec(`
INSERT INTO "credit_buckets" ("id", "user_id", "source_type", "plan_id", "credits_granted", "credits_remaining", "valid_from")
VALUES (?, ?, ?, ?, ?, ?, now())
`, bucketID, user.ID, "monthly_allowance", "starter", 1200, 1200).Error; err != nil {
		t.Fatalf("insert legacy credit bucket: %v", err)
	}

	if err := dbmigrate.MigrateCreditOnlyBilling(db); err != nil {
		t.Fatalf("MigrateCreditOnlyBilling() error = %v", err)
	}
	if err := dbmigrate.MigrateCreditOnlyBilling(db); err != nil {
		t.Fatalf("repeat MigrateCreditOnlyBilling() error = %v", err)
	}

	if billingTableExists(t, db, "billing_plans") {
		t.Fatal("billing_plans exists after credit-only migration")
	}
	if billingTableExists(t, db, "billing_subscriptions") {
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
		if billingColumnExists(t, db, check.table, check.column) {
			t.Fatalf("%s.%s exists after credit-only migration", check.table, check.column)
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

	var migratedBucket billing.CreditBucket
	if err := db.First(&migratedBucket, "id = ?", bucketID).Error; err != nil {
		t.Fatalf("load migrated bucket: %v", err)
	}
	if migratedBucket.SourceType != billing.CreditSourceAdjustment {
		t.Fatalf("migrated bucket source_type = %q, want %q", migratedBucket.SourceType, billing.CreditSourceAdjustment)
	}

	if err := db.AutoMigrate(&auth.User{}, &billing.BillingOrder{}, &billing.CreditBucket{}, &billing.CreditLedgerEntry{}); err != nil {
		t.Fatalf("auto migrate credit-only models: %v", err)
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
		t.Fatalf("create credit-only order after migration and AutoMigrate: %v", err)
	}
}

func resetBillingTables(t *testing.T, db *gorm.DB) {
	t.Helper()
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
}

func createLegacyBillingPlanSchema(t *testing.T, db *gorm.DB) {
	t.Helper()
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
}

func billingTableExists(t *testing.T, db *gorm.DB, tableName string) bool {
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

func billingColumnExists(t *testing.T, db *gorm.DB, tableName string, columnName string) bool {
	t.Helper()
	var count int64
	if err := db.Raw(`
SELECT COUNT(*)
FROM information_schema.columns
WHERE table_schema = current_schema()
	AND table_name = ?
	AND column_name = ?
`, tableName, columnName).Scan(&count).Error; err != nil {
		t.Fatalf("query %s.%s existence: %v", tableName, columnName, err)
	}
	return count > 0
}
