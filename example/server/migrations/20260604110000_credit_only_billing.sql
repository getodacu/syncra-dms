ALTER TABLE "credit_buckets" DROP CONSTRAINT IF EXISTS "fk_credit_buckets_plan";
DROP INDEX IF EXISTS "idx_credit_buckets_plan_id";

UPDATE "credit_buckets"
SET
  "source_type" = 'adjustment'
WHERE "source_type" = 'monthly_allowance';

ALTER TABLE "credit_buckets" DROP CONSTRAINT IF EXISTS "chk_credit_buckets_source_type";
ALTER TABLE "credit_buckets" DROP COLUMN IF EXISTS "plan_id";
ALTER TABLE "credit_buckets" ADD CONSTRAINT "chk_credit_buckets_source_type" CHECK ("source_type" IN ('signup_bonus', 'topup_purchase', 'refund', 'adjustment'));

DROP INDEX IF EXISTS "idx_billing_orders_plan_id_at_purchase";
ALTER TABLE "billing_orders" DROP CONSTRAINT IF EXISTS "chk_billing_orders_plan_id_at_purchase";
ALTER TABLE "billing_orders" DROP CONSTRAINT IF EXISTS "chk_billing_orders_credit_blocks";
ALTER TABLE "billing_orders" DROP COLUMN IF EXISTS "plan_id_at_purchase";

ALTER TABLE "billing_orders" ADD COLUMN "provider" varchar(40);
ALTER TABLE "billing_orders" ADD COLUMN "pricing_tier" varchar(40);
ALTER TABLE "billing_orders" ADD COLUMN "unit_amount_cents" bigint;

UPDATE "billing_orders"
SET
  "provider" = 'stripe',
  "pricing_tier" = CASE
    WHEN "credits" >= 20000 THEN 'tier_4'
    WHEN "credits" >= 10000 THEN 'tier_3'
    WHEN "credits" >= 5000 THEN 'tier_2'
    ELSE 'tier_1'
  END,
  "unit_amount_cents" = CASE
    WHEN "amount_cents" > 0 THEN ("amount_cents" * 1000) / "credits"
    WHEN "credits" >= 20000 THEN 850
    WHEN "credits" >= 10000 THEN 900
    WHEN "credits" >= 5000 THEN 950
    ELSE 1000
  END
WHERE "provider" IS NULL
  OR "pricing_tier" IS NULL
  OR "unit_amount_cents" IS NULL;

ALTER TABLE "billing_orders" DROP COLUMN IF EXISTS "credit_blocks";

ALTER TABLE "billing_orders" ALTER COLUMN "provider" SET NOT NULL;
ALTER TABLE "billing_orders" ALTER COLUMN "pricing_tier" SET NOT NULL;
ALTER TABLE "billing_orders" ALTER COLUMN "unit_amount_cents" SET NOT NULL;
ALTER TABLE "billing_orders" ADD CONSTRAINT "chk_billing_orders_provider" CHECK ("provider" IN ('stripe'));
ALTER TABLE "billing_orders" ADD CONSTRAINT "chk_billing_orders_pricing_tier" CHECK ("pricing_tier" IN ('tier_1', 'tier_2', 'tier_3', 'tier_4'));
ALTER TABLE "billing_orders" ADD CONSTRAINT "chk_billing_orders_unit_amount_cents" CHECK ("unit_amount_cents" > 0);

DROP TABLE IF EXISTS "billing_subscriptions";
DROP TABLE IF EXISTS "billing_plans";
