ALTER TABLE "user" ADD COLUMN IF NOT EXISTS "role" varchar(40) NOT NULL DEFAULT 'user';
ALTER TABLE "user" ADD COLUMN IF NOT EXISTS "last_login_at" timestamptz;

-- Bootstrap the first admin manually after applying this migration.
-- Replace the email with a known verified account:
-- UPDATE "user" SET "role" = 'admin' WHERE "email" = 'admin@example.com';

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM pg_constraint
    WHERE conname = 'chk_user_role'
  ) THEN
    ALTER TABLE "user" ADD CONSTRAINT "chk_user_role" CHECK ("role" IN ('user','admin'));
  END IF;
END $$;

CREATE INDEX IF NOT EXISTS "idx_user_role" ON "user" ("role");
CREATE INDEX IF NOT EXISTS "idx_user_created_id" ON "user" ("created_at", "id");
CREATE INDEX IF NOT EXISTS "idx_user_last_login_id" ON "user" ("last_login_at", "id");
CREATE INDEX IF NOT EXISTS "idx_user_name_trgm" ON "user" USING gin (lower("name") gin_trgm_ops);
CREATE INDEX IF NOT EXISTS "idx_user_email_trgm" ON "user" USING gin (lower("email") gin_trgm_ops);
