ALTER TABLE "session" ADD COLUMN IF NOT EXISTS "impersonated_user_id" uuid;
ALTER TABLE "session" ADD COLUMN IF NOT EXISTS "impersonation_started_at" timestamptz;

CREATE INDEX IF NOT EXISTS "idx_session_impersonated_user_id" ON "session" ("impersonated_user_id");

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_constraint WHERE conname = 'fk_session_impersonated_user'
  ) THEN
    ALTER TABLE "session" ADD CONSTRAINT "fk_session_impersonated_user" FOREIGN KEY ("impersonated_user_id") REFERENCES "user"("id") ON DELETE SET NULL ON UPDATE CASCADE;
  END IF;
END $$;

CREATE TABLE IF NOT EXISTS "admin_impersonation_events" (
  "id" uuid NOT NULL,
  "event_type" varchar(20) NOT NULL,
  "session_id" uuid NOT NULL,
  "admin_user_id" uuid NOT NULL,
  "admin_user_email" varchar(320) NOT NULL,
  "target_user_id" uuid NOT NULL,
  "target_user_email" varchar(320) NOT NULL,
  "ip_address" text,
  "user_agent" text,
  "created_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "chk_admin_impersonation_event_type" CHECK ("event_type" IN ('start','stop'))
);

CREATE INDEX IF NOT EXISTS "idx_admin_impersonation_events_event_type" ON "admin_impersonation_events" ("event_type");
CREATE INDEX IF NOT EXISTS "idx_admin_impersonation_events_session_id" ON "admin_impersonation_events" ("session_id");
CREATE INDEX IF NOT EXISTS "idx_admin_impersonation_events_admin_user_id" ON "admin_impersonation_events" ("admin_user_id");
CREATE INDEX IF NOT EXISTS "idx_admin_impersonation_events_target_user_id" ON "admin_impersonation_events" ("target_user_id");
CREATE INDEX IF NOT EXISTS "idx_admin_impersonation_events_created_at" ON "admin_impersonation_events" ("created_at");
