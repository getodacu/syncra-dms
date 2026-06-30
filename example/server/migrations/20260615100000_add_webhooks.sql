CREATE TABLE "webhooks" (
  "id" uuid,
  "user_id" uuid NOT NULL,
  "url" text NOT NULL,
  "secret_key" text NOT NULL,
  "events_active" jsonb NOT NULL DEFAULT '[]'::jsonb,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "chk_webhooks_events_active_array" CHECK (jsonb_typeof("events_active") = 'array')
);

CREATE UNIQUE INDEX "idx_webhooks_user_id" ON "webhooks" ("user_id");

ALTER TABLE "webhooks" ADD CONSTRAINT "fk_webhooks_user" FOREIGN KEY ("user_id") REFERENCES "user"("id") ON DELETE CASCADE ON UPDATE CASCADE;
