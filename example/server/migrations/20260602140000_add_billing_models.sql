CREATE TABLE "billing_plans" (
  "id" varchar(40),
  "created_at" timestamptz,
  "updated_at" timestamptz,
  "name" varchar(120) NOT NULL,
  "monthly_price_cents" bigint NOT NULL DEFAULT 0,
  "monthly_credits" bigint NOT NULL DEFAULT 0,
  "topup_block_credits" bigint NOT NULL DEFAULT 0,
  "topup_block_price_cents" bigint NOT NULL DEFAULT 0,
  "currency" varchar(3) NOT NULL,
  "active" boolean NOT NULL DEFAULT true,
  PRIMARY KEY ("id"),
  CONSTRAINT "chk_billing_plans_id" CHECK ("id" IN ('credit_only', 'starter', 'pro')),
  CONSTRAINT "chk_billing_plans_monthly_price_cents" CHECK ("monthly_price_cents" >= 0),
  CONSTRAINT "chk_billing_plans_monthly_credits" CHECK ("monthly_credits" >= 0),
  CONSTRAINT "chk_billing_plans_topup_block_credits" CHECK ("topup_block_credits" > 0),
  CONSTRAINT "chk_billing_plans_topup_block_price_cents" CHECK ("topup_block_price_cents" >= 0)
);

CREATE INDEX "idx_billing_plans_active" ON "billing_plans" ("active");

CREATE TABLE "billing_subscriptions" (
  "id" uuid,
  "user_id" uuid NOT NULL,
  "plan_id" varchar(40) NOT NULL,
  "status" varchar(40) NOT NULL,
  "current_period_start" timestamptz,
  "current_period_end" timestamptz,
  "provider" varchar(80),
  "provider_customer_id" varchar(255),
  "provider_subscription_id" varchar(255),
  "canceled_at" timestamptz,
  "created_at" timestamptz,
  "updated_at" timestamptz,
  PRIMARY KEY ("id"),
  CONSTRAINT "chk_billing_subscriptions_status" CHECK ("status" IN ('active', 'past_due', 'canceled'))
);

CREATE INDEX "idx_billing_subscriptions_status" ON "billing_subscriptions" ("status");
CREATE INDEX "idx_billing_subscriptions_plan_id" ON "billing_subscriptions" ("plan_id");
CREATE INDEX "idx_billing_subscriptions_provider_customer_id" ON "billing_subscriptions" ("provider_customer_id");
CREATE UNIQUE INDEX "idx_billing_subscriptions_provider_subscription_id" ON "billing_subscriptions" ("provider_subscription_id");
CREATE UNIQUE INDEX "idx_billing_subscriptions_user_id" ON "billing_subscriptions" ("user_id");

CREATE TABLE "billing_orders" (
  "id" uuid,
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
  PRIMARY KEY ("id"),
  CONSTRAINT "chk_billing_orders_order_type" CHECK ("order_type" IN ('credit_topup')),
  CONSTRAINT "chk_billing_orders_status" CHECK ("status" IN ('pending', 'paid', 'failed', 'refunded', 'canceled')),
  CONSTRAINT "chk_billing_orders_plan_id_at_purchase" CHECK ("plan_id_at_purchase" IN ('credit_only', 'starter', 'pro')),
  CONSTRAINT "chk_billing_orders_credit_blocks" CHECK ("credit_blocks" > 0),
  CONSTRAINT "chk_billing_orders_credits" CHECK ("credits" > 0),
  CONSTRAINT "chk_billing_orders_amount_cents" CHECK ("amount_cents" >= 0)
);

CREATE INDEX "idx_billing_orders_plan_id_at_purchase" ON "billing_orders" ("plan_id_at_purchase");
CREATE INDEX "idx_billing_orders_status" ON "billing_orders" ("status");
CREATE INDEX "idx_billing_orders_order_type" ON "billing_orders" ("order_type");
CREATE UNIQUE INDEX "idx_billing_orders_provider_checkout_session_id" ON "billing_orders" ("provider_checkout_session_id") WHERE "provider_checkout_session_id" IS NOT NULL;
CREATE UNIQUE INDEX "idx_billing_orders_provider_payment_intent_id" ON "billing_orders" ("provider_payment_intent_id") WHERE "provider_payment_intent_id" IS NOT NULL;
CREATE INDEX "idx_billing_orders_user_id" ON "billing_orders" ("user_id");

CREATE TABLE "credit_buckets" (
  "id" uuid,
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
  PRIMARY KEY ("id"),
  CONSTRAINT "chk_credit_buckets_source_type" CHECK ("source_type" IN ('signup_bonus', 'monthly_allowance', 'topup_purchase', 'refund', 'adjustment')),
  CONSTRAINT "chk_credit_buckets_credits_granted" CHECK ("credits_granted" > 0),
  CONSTRAINT "chk_credit_buckets_credits_remaining" CHECK ("credits_remaining" >= 0 AND "credits_remaining" <= "credits_granted")
);

CREATE INDEX "idx_credit_buckets_voided_at" ON "credit_buckets" ("voided_at");
CREATE INDEX "idx_credit_buckets_expires_at" ON "credit_buckets" ("expires_at");
CREATE INDEX "idx_credit_buckets_valid_from" ON "credit_buckets" ("valid_from");
CREATE INDEX "idx_credit_buckets_order_id" ON "credit_buckets" ("order_id");
CREATE INDEX "idx_credit_buckets_plan_id" ON "credit_buckets" ("plan_id");
CREATE INDEX "idx_credit_buckets_source_type" ON "credit_buckets" ("source_type");
CREATE UNIQUE INDEX "idx_credit_buckets_one_signup_bonus" ON "credit_buckets" ("user_id") WHERE "source_type" = 'signup_bonus';
CREATE INDEX "idx_credit_buckets_user_available" ON "credit_buckets" ("user_id", "expires_at", "created_at") WHERE "credits_remaining" > 0 AND "voided_at" IS NULL;
CREATE INDEX "idx_credit_buckets_user_id" ON "credit_buckets" ("user_id");

CREATE TABLE "credit_ledger_entries" (
  "id" uuid,
  "user_id" uuid NOT NULL,
  "bucket_id" uuid,
  "entry_type" varchar(40) NOT NULL,
  "credits_delta" bigint NOT NULL,
  "related_job_id" uuid,
  "related_order_id" uuid,
  "idempotency_key" varchar(255) NOT NULL,
  "metadata" jsonb,
  "created_at" timestamptz,
  PRIMARY KEY ("id"),
  CONSTRAINT "chk_credit_ledger_entries_entry_type" CHECK ("entry_type" IN ('grant', 'purchase', 'debit', 'refund', 'expiry', 'adjustment')),
  CONSTRAINT "chk_credit_ledger_entries_credits_delta" CHECK ("credits_delta" <> 0)
);

CREATE UNIQUE INDEX "idx_credit_ledger_entries_idempotency_key" ON "credit_ledger_entries" ("idempotency_key");
CREATE INDEX "idx_credit_ledger_entries_related_order_id" ON "credit_ledger_entries" ("related_order_id");
CREATE INDEX "idx_credit_ledger_entries_related_job_id" ON "credit_ledger_entries" ("related_job_id");
CREATE INDEX "idx_credit_ledger_entries_entry_type" ON "credit_ledger_entries" ("entry_type");
CREATE INDEX "idx_credit_ledger_entries_bucket_id" ON "credit_ledger_entries" ("bucket_id");
CREATE INDEX "idx_credit_ledger_entries_user_id" ON "credit_ledger_entries" ("user_id");

ALTER TABLE "billing_subscriptions" ADD CONSTRAINT "fk_billing_subscriptions_plan" FOREIGN KEY ("plan_id") REFERENCES "billing_plans"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "billing_subscriptions" ADD CONSTRAINT "fk_billing_subscriptions_user" FOREIGN KEY ("user_id") REFERENCES "user"("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "billing_orders" ADD CONSTRAINT "fk_billing_orders_user" FOREIGN KEY ("user_id") REFERENCES "user"("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "credit_buckets" ADD CONSTRAINT "fk_credit_buckets_order" FOREIGN KEY ("order_id") REFERENCES "billing_orders"("id") ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE "credit_buckets" ADD CONSTRAINT "fk_credit_buckets_plan" FOREIGN KEY ("plan_id") REFERENCES "billing_plans"("id") ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE "credit_buckets" ADD CONSTRAINT "fk_credit_buckets_user" FOREIGN KEY ("user_id") REFERENCES "user"("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "credit_ledger_entries" ADD CONSTRAINT "fk_credit_ledger_entries_bucket" FOREIGN KEY ("bucket_id") REFERENCES "credit_buckets"("id") ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE "credit_ledger_entries" ADD CONSTRAINT "fk_credit_ledger_entries_related_order" FOREIGN KEY ("related_order_id") REFERENCES "billing_orders"("id") ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE "credit_ledger_entries" ADD CONSTRAINT "fk_credit_ledger_entries_user" FOREIGN KEY ("user_id") REFERENCES "user"("id") ON DELETE CASCADE ON UPDATE CASCADE;

CREATE FUNCTION "prevent_credit_ledger_entry_update"() RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
  RAISE EXCEPTION 'credit_ledger_entries rows are append-only';
END;
$$;

CREATE TRIGGER "trg_prevent_credit_ledger_entry_update"
BEFORE UPDATE ON "credit_ledger_entries"
FOR EACH ROW
EXECUTE FUNCTION "prevent_credit_ledger_entry_update"();
