CREATE TABLE "billing_invoice_counters" (
  "invoice_serie" varchar(40),
  "last_number" bigint NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("invoice_serie"),
  CONSTRAINT "chk_billing_invoice_counters_last_number" CHECK ("last_number" >= 0)
);

CREATE TABLE "billing_invoices" (
  "id" uuid,
  "user_id" uuid,
  "billing_profile_id" uuid,
  "billing_name" varchar(255) NOT NULL,
  "billing_email" varchar(320) NOT NULL,
  "billing_fiscal_code" varchar(80),
  "billing_profile_snapshot" jsonb NOT NULL,
  "lines" jsonb NOT NULL,
  "net_amount" numeric(20,2) NOT NULL,
  "vat_amount" numeric(20,2) NOT NULL,
  "total_amount" numeric(20,2) NOT NULL,
  "invoice_date" date NOT NULL,
  "invoice_serie" varchar(40) NOT NULL,
  "invoice_number" bigint NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "chk_billing_invoices_net_amount" CHECK ("net_amount" >= 0),
  CONSTRAINT "chk_billing_invoices_vat_amount" CHECK ("vat_amount" >= 0),
  CONSTRAINT "chk_billing_invoices_total_amount" CHECK ("total_amount" >= 0 AND "total_amount" = "net_amount" + "vat_amount"),
  CONSTRAINT "chk_billing_invoices_invoice_number" CHECK ("invoice_number" > 0),
  CONSTRAINT "chk_billing_invoices_billing_profile_snapshot" CHECK (jsonb_typeof("billing_profile_snapshot") = 'object'),
  CONSTRAINT "chk_billing_invoices_lines" CHECK (jsonb_typeof("lines") = 'array')
);

CREATE UNIQUE INDEX "idx_billing_invoices_serie_number" ON "billing_invoices" ("invoice_serie", "invoice_number");
CREATE INDEX "idx_billing_invoices_user_id" ON "billing_invoices" ("user_id");
CREATE INDEX "idx_billing_invoices_billing_profile_id" ON "billing_invoices" ("billing_profile_id");
CREATE INDEX "idx_billing_invoices_invoice_date" ON "billing_invoices" ("invoice_date");
CREATE INDEX "idx_billing_invoices_billing_fiscal_code" ON "billing_invoices" ("billing_fiscal_code") WHERE "billing_fiscal_code" IS NOT NULL;
CREATE INDEX "idx_billing_invoices_billing_name_trgm" ON "billing_invoices" USING gin (lower("billing_name") gin_trgm_ops);
CREATE INDEX "idx_billing_invoices_billing_email_trgm" ON "billing_invoices" USING gin (lower("billing_email") gin_trgm_ops);
CREATE INDEX "idx_billing_invoices_billing_fiscal_code_trgm" ON "billing_invoices" USING gin (lower("billing_fiscal_code") gin_trgm_ops) WHERE "billing_fiscal_code" IS NOT NULL;

ALTER TABLE "billing_invoices" ADD CONSTRAINT "fk_billing_invoices_user" FOREIGN KEY ("user_id") REFERENCES "user"("id") ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE "billing_invoices" ADD CONSTRAINT "fk_billing_invoices_billing_profile" FOREIGN KEY ("billing_profile_id") REFERENCES "billing_profiles"("id") ON DELETE SET NULL ON UPDATE CASCADE;
