ALTER TABLE "billing_invoices" ADD COLUMN "order_id" uuid;

CREATE UNIQUE INDEX "idx_billing_invoices_order_id" ON "billing_invoices" ("order_id") WHERE "order_id" IS NOT NULL;

ALTER TABLE "billing_invoices" ADD CONSTRAINT "fk_billing_invoices_order" FOREIGN KEY ("order_id") REFERENCES "billing_orders"("id") ON DELETE SET NULL ON UPDATE CASCADE;
