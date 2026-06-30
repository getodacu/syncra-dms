ALTER TABLE billing_invoices
ADD COLUMN email_delivery_claimed_at TIMESTAMPTZ,
ADD COLUMN email_sent_at TIMESTAMPTZ;
