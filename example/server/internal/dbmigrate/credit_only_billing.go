package dbmigrate

import "gorm.io/gorm"

const (
	creditBucketSourceTypeConstraint      = "chk_credit_buckets_source_type"
	billingOrderProviderConstraint        = "chk_billing_orders_provider"
	billingOrderPricingTierConstraint     = "chk_billing_orders_pricing_tier"
	billingOrderUnitAmountCentsConstraint = "chk_billing_orders_unit_amount_cents"
)

// MigrateCreditOnlyBilling removes legacy billing-plan schema left behind by
// reused databases before GORM AutoMigrate runs the current credit-only models.
func MigrateCreditOnlyBilling(db *gorm.DB) error {
	creditBucketsExists, err := tableExists(db, "credit_buckets")
	if err != nil {
		return err
	}
	_, creditBucketsSourceTypeExists, err := columnDataType(db, "credit_buckets", "source_type")
	if err != nil {
		return err
	}
	_, creditBucketsPlanIDExists, err := columnDataType(db, "credit_buckets", "plan_id")
	if err != nil {
		return err
	}

	billingOrdersExists, err := tableExists(db, "billing_orders")
	if err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if creditBucketsExists {
			if err := migrateCreditBucketsToCreditOnly(tx, creditBucketsSourceTypeExists, creditBucketsPlanIDExists); err != nil {
				return err
			}
		}
		if billingOrdersExists {
			if err := migrateBillingOrdersToCreditOnly(tx); err != nil {
				return err
			}
		}
		if err := tx.Exec(`DROP TABLE IF EXISTS "billing_subscriptions"`).Error; err != nil {
			return err
		}
		return tx.Exec(`DROP TABLE IF EXISTS "billing_plans"`).Error
	})
}

func migrateCreditBucketsToCreditOnly(db *gorm.DB, sourceTypeExists bool, planIDExists bool) error {
	if err := db.Exec(`ALTER TABLE "credit_buckets" DROP CONSTRAINT IF EXISTS "fk_credit_buckets_plan"`).Error; err != nil {
		return err
	}
	if err := db.Exec(`DROP INDEX IF EXISTS "idx_credit_buckets_plan_id"`).Error; err != nil {
		return err
	}
	if sourceTypeExists {
		if err := db.Exec(`ALTER TABLE "credit_buckets" DROP CONSTRAINT IF EXISTS "` + creditBucketSourceTypeConstraint + `"`).Error; err != nil {
			return err
		}
		if err := db.Exec(`UPDATE "credit_buckets" SET "source_type" = 'adjustment' WHERE "source_type" = 'monthly_allowance'`).Error; err != nil {
			return err
		}
	}
	if planIDExists {
		if err := db.Exec(`ALTER TABLE "credit_buckets" DROP COLUMN "plan_id"`).Error; err != nil {
			return err
		}
	}
	if sourceTypeExists {
		return db.Exec(`ALTER TABLE "credit_buckets" ADD CONSTRAINT "` + creditBucketSourceTypeConstraint + `" CHECK ("source_type" IN ('signup_bonus', 'topup_purchase', 'refund', 'adjustment'))`).Error
	}
	return nil
}

func migrateBillingOrdersToCreditOnly(db *gorm.DB) error {
	if err := db.Exec(`DROP INDEX IF EXISTS "idx_billing_orders_plan_id_at_purchase"`).Error; err != nil {
		return err
	}
	for _, constraintName := range []string{
		"chk_billing_orders_plan_id_at_purchase",
		"chk_billing_orders_credit_blocks",
		billingOrderProviderConstraint,
		billingOrderPricingTierConstraint,
		billingOrderUnitAmountCentsConstraint,
	} {
		if err := db.Exec(`ALTER TABLE "billing_orders" DROP CONSTRAINT IF EXISTS "` + constraintName + `"`).Error; err != nil {
			return err
		}
	}
	if err := db.Exec(`ALTER TABLE "billing_orders" DROP COLUMN IF EXISTS "plan_id_at_purchase"`).Error; err != nil {
		return err
	}
	if err := db.Exec(`ALTER TABLE "billing_orders" ADD COLUMN IF NOT EXISTS "provider" varchar(40)`).Error; err != nil {
		return err
	}
	if err := db.Exec(`ALTER TABLE "billing_orders" ADD COLUMN IF NOT EXISTS "pricing_tier" varchar(40)`).Error; err != nil {
		return err
	}
	if err := db.Exec(`ALTER TABLE "billing_orders" ADD COLUMN IF NOT EXISTS "unit_amount_cents" bigint`).Error; err != nil {
		return err
	}
	if err := db.Exec(`
UPDATE "billing_orders"
SET
  "provider" = COALESCE("provider", 'stripe'),
  "pricing_tier" = COALESCE("pricing_tier", CASE
    WHEN "credits" >= 20000 THEN 'tier_4'
    WHEN "credits" >= 10000 THEN 'tier_3'
    WHEN "credits" >= 5000 THEN 'tier_2'
    ELSE 'tier_1'
  END),
  "unit_amount_cents" = COALESCE("unit_amount_cents", CASE
    WHEN "amount_cents" > 0 AND "credits" > 0 THEN ("amount_cents" * 1000) / "credits"
    WHEN "credits" >= 20000 THEN 850
    WHEN "credits" >= 10000 THEN 900
    WHEN "credits" >= 5000 THEN 950
    ELSE 1000
  END)
WHERE "provider" IS NULL
  OR "pricing_tier" IS NULL
  OR "unit_amount_cents" IS NULL
`).Error; err != nil {
		return err
	}
	if err := db.Exec(`ALTER TABLE "billing_orders" DROP COLUMN IF EXISTS "credit_blocks"`).Error; err != nil {
		return err
	}
	for _, columnName := range []string{"provider", "pricing_tier", "unit_amount_cents"} {
		if err := db.Exec(`ALTER TABLE "billing_orders" ALTER COLUMN "` + columnName + `" SET NOT NULL`).Error; err != nil {
			return err
		}
	}
	if err := db.Exec(`ALTER TABLE "billing_orders" ADD CONSTRAINT "` + billingOrderProviderConstraint + `" CHECK ("provider" IN ('stripe'))`).Error; err != nil {
		return err
	}
	if err := db.Exec(`ALTER TABLE "billing_orders" ADD CONSTRAINT "` + billingOrderPricingTierConstraint + `" CHECK ("pricing_tier" IN ('tier_1', 'tier_2', 'tier_3', 'tier_4'))`).Error; err != nil {
		return err
	}
	return db.Exec(`ALTER TABLE "billing_orders" ADD CONSTRAINT "` + billingOrderUnitAmountCentsConstraint + `" CHECK ("unit_amount_cents" > 0)`).Error
}

func tableExists(db *gorm.DB, tableName string) (bool, error) {
	var exists bool
	err := db.Raw(`
SELECT EXISTS (
	SELECT 1
	FROM information_schema.tables
	WHERE table_schema = current_schema()
		AND table_name = ?
)
`, tableName).Scan(&exists).Error
	return exists, err
}
