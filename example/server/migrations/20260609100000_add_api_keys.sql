CREATE TABLE "api_keys" (
  "id" uuid,
  "user_id" uuid NOT NULL,
  "name" varchar(255) NOT NULL,
  "key_hash" varchar(64) NOT NULL,
  "key_prefix" varchar(8) NOT NULL,
  "expires_at" timestamptz,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);

CREATE INDEX "idx_api_keys_expires_at" ON "api_keys" ("expires_at");
CREATE UNIQUE INDEX "idx_api_keys_key_hash" ON "api_keys" ("key_hash");
CREATE INDEX "idx_api_keys_user_id" ON "api_keys" ("user_id");

ALTER TABLE "api_keys" ADD CONSTRAINT "fk_api_keys_user" FOREIGN KEY ("user_id") REFERENCES "user"("id") ON DELETE CASCADE ON UPDATE CASCADE;
