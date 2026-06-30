CREATE TABLE "user" (
  "id" uuid,
  "name" varchar(255) NOT NULL,
  "email" varchar(320) NOT NULL,
  "email_verified" boolean NOT NULL DEFAULT false,
  "image" text,
  "preferred_language" varchar(5) NOT NULL DEFAULT 'en',
  "role" varchar(40) NOT NULL DEFAULT 'user',
  "last_login_at" timestamptz,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "chk_user_preferred_language" CHECK (preferred_language IN ('en','ro')),
  CONSTRAINT "chk_user_role" CHECK (role IN ('user','admin'))
);

CREATE UNIQUE INDEX "idx_user_email" ON "user" ("email");
CREATE INDEX "idx_user_role" ON "user" ("role");
CREATE INDEX "idx_user_created_id" ON "user" ("created_at", "id");
CREATE INDEX "idx_user_last_login_id" ON "user" ("last_login_at", "id");

CREATE TABLE "account" (
  "id" uuid,
  "account_id" varchar(255) NOT NULL,
  "provider_id" varchar(120) NOT NULL,
  "user_id" uuid NOT NULL,
  "access_token" text,
  "refresh_token" text,
  "id_token" text,
  "access_token_expires_at" timestamptz,
  "refresh_token_expires_at" timestamptz,
  "scope" text,
  "password" text,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);

CREATE INDEX "idx_account_user_id" ON "account" ("user_id");
CREATE INDEX "idx_account_provider_id" ON "account" ("provider_id");
CREATE INDEX "idx_account_account_id" ON "account" ("account_id");
CREATE UNIQUE INDEX "idx_account_provider_account" ON "account" ("provider_id", "account_id");

CREATE TABLE "session" (
  "id" uuid,
  "expires_at" timestamptz NOT NULL,
  "token" varchar(255) NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "ip_address" text,
  "user_agent" text,
  "user_id" uuid NOT NULL,
  PRIMARY KEY ("id")
);

CREATE INDEX "idx_session_user_id" ON "session" ("user_id");
CREATE UNIQUE INDEX "idx_session_token" ON "session" ("token");
CREATE INDEX "idx_session_expires_at" ON "session" ("expires_at");

CREATE TABLE "verification" (
  "id" uuid,
  "identifier" varchar(512) NOT NULL,
  "value" text NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX "idx_verification_identifier_unique" ON "verification" ("identifier");
CREATE INDEX "idx_verification_expires_at" ON "verification" ("expires_at");

ALTER TABLE "account" ADD CONSTRAINT "fk_account_user" FOREIGN KEY ("user_id") REFERENCES "user"("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "session" ADD CONSTRAINT "fk_session_user" FOREIGN KEY ("user_id") REFERENCES "user"("id") ON DELETE CASCADE ON UPDATE CASCADE;
