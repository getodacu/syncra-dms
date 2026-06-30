CREATE INDEX "idx_credit_ledger_entries_transactions" ON "credit_ledger_entries" ("user_id", "created_at", "id") WHERE "entry_type" = 'purchase' OR "entry_type" = 'debit';
