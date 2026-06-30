CREATE EXTENSION IF NOT EXISTS "pg_trgm";

CREATE TABLE "user" (
  "id" uuid,
  "name" varchar(255) NOT NULL,
  "email" varchar(320) NOT NULL,
  "email_verified" boolean NOT NULL DEFAULT false,
  "image" text,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX "idx_user_email" ON "user" ("email");

CREATE TABLE "extraction_schemas" (
  "id" uuid,
  "user_id" uuid,
  "created_at" timestamptz,
  "updated_at" timestamptz,
  "name" varchar(160) NOT NULL,
  "description" text,
  "schema_json" jsonb NOT NULL,
  "strict" boolean NOT NULL DEFAULT true,
  PRIMARY KEY ("id")
);

CREATE INDEX "idx_extraction_schemas_user_id" ON "extraction_schemas" ("user_id");

CREATE TABLE "ocr_documents" (
  "id" uuid,
  "user_id" uuid,
  "created_at" timestamptz,
  "updated_at" timestamptz,
  "deleted_at" timestamptz,
  "original_filename" varchar(255) NOT NULL,
  "mime_type" varchar(120) NOT NULL,
  "file_size" bigint NOT NULL,
  "page_count" bigint NOT NULL DEFAULT 0,
  "document_hash" varchar(64) NOT NULL,
  "job_id" uuid,
  "schema_id" uuid,
  "inline_schema_json" jsonb,
  "markdown" text NOT NULL,
  "annotation_json" jsonb,
  "raw_response_json" jsonb NOT NULL,
  PRIMARY KEY ("id")
);

CREATE INDEX "idx_ocr_documents_schema_id" ON "ocr_documents" ("schema_id");
CREATE INDEX "idx_ocr_documents_job_id" ON "ocr_documents" ("job_id");
CREATE INDEX "idx_ocr_documents_document_hash" ON "ocr_documents" ("document_hash");
CREATE INDEX "idx_ocr_documents_deleted_at" ON "ocr_documents" ("deleted_at");
CREATE INDEX "idx_ocr_documents_user_id" ON "ocr_documents" ("user_id");
CREATE INDEX "idx_ocr_documents_original_filename_trgm" ON "ocr_documents" USING gin (lower("original_filename") gin_trgm_ops);

CREATE TABLE "ocr_jobs" (
  "id" uuid,
  "user_id" uuid,
  "created_at" timestamptz,
  "updated_at" timestamptz,
  "original_filename" varchar(255) NOT NULL,
  "mime_type" varchar(120) NOT NULL,
  "file_size" bigint NOT NULL,
  "page_count" bigint NOT NULL DEFAULT 0,
  "document_hash" varchar(64) NOT NULL,
  "file_path" text NOT NULL,
  "schema_id" uuid,
  "inline_schema_json" jsonb,
  "document_id" uuid,
  "status" varchar(40) NOT NULL DEFAULT 'queued',
  "error_message" text,
  PRIMARY KEY ("id"),
  CONSTRAINT "chk_ocr_jobs_status" CHECK ("status" IN ('queued', 'processing', 'completed', 'failed'))
);

CREATE INDEX "idx_ocr_jobs_status" ON "ocr_jobs" ("status");
CREATE INDEX "idx_ocr_jobs_document_id" ON "ocr_jobs" ("document_id");
CREATE INDEX "idx_ocr_jobs_schema_id" ON "ocr_jobs" ("schema_id");
CREATE INDEX "idx_ocr_jobs_document_hash" ON "ocr_jobs" ("document_hash");
CREATE INDEX "idx_ocr_jobs_user_id" ON "ocr_jobs" ("user_id");
CREATE INDEX "idx_ocr_jobs_user_status_created_id" ON "ocr_jobs" ("user_id", "status", "created_at", "id");
CREATE INDEX "idx_ocr_jobs_user_created_id" ON "ocr_jobs" ("user_id", "created_at", "id");

CREATE TABLE "collections" (
  "id" uuid,
  "user_id" uuid NOT NULL,
  "created_at" timestamptz,
  "updated_at" timestamptz,
  "name" varchar(160) NOT NULL,
  PRIMARY KEY ("id")
);

CREATE INDEX "idx_collections_user_id" ON "collections" ("user_id");
CREATE INDEX "idx_collections_user_created_id" ON "collections" ("user_id", "created_at", "id");

CREATE TABLE "collection_schemas" (
  "collection_id" uuid,
  "schema_id" uuid,
  "created_at" timestamptz,
  PRIMARY KEY ("collection_id", "schema_id")
);

CREATE INDEX "idx_collection_schemas_schema_id" ON "collection_schemas" ("schema_id");
CREATE INDEX "idx_collection_schemas_collection_id" ON "collection_schemas" ("collection_id");

CREATE TABLE "collection_documents" (
  "collection_id" uuid,
  "document_id" uuid,
  "created_at" timestamptz,
  PRIMARY KEY ("collection_id", "document_id")
);

CREATE INDEX "idx_collection_documents_document_id" ON "collection_documents" ("document_id");
CREATE INDEX "idx_collection_documents_collection_id" ON "collection_documents" ("collection_id");

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
CREATE UNIQUE INDEX "idx_account_provider_account" ON "account" ("provider_id", "account_id");
CREATE INDEX "idx_account_account_id" ON "account" ("account_id");

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

CREATE INDEX "idx_verification_expires_at" ON "verification" ("expires_at");
CREATE UNIQUE INDEX "idx_verification_identifier_unique" ON "verification" ("identifier");

ALTER TABLE "extraction_schemas" ADD CONSTRAINT "fk_extraction_schemas_user" FOREIGN KEY ("user_id") REFERENCES "user"("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "ocr_documents" ADD CONSTRAINT "fk_ocr_documents_schema" FOREIGN KEY ("schema_id") REFERENCES "extraction_schemas"("id") ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE "ocr_documents" ADD CONSTRAINT "fk_ocr_documents_user" FOREIGN KEY ("user_id") REFERENCES "user"("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "ocr_documents" ADD CONSTRAINT "fk_ocr_documents_job" FOREIGN KEY ("job_id") REFERENCES "ocr_jobs"("id") ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE "ocr_jobs" ADD CONSTRAINT "fk_ocr_jobs_document" FOREIGN KEY ("document_id") REFERENCES "ocr_documents"("id") ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE "ocr_jobs" ADD CONSTRAINT "fk_ocr_jobs_schema" FOREIGN KEY ("schema_id") REFERENCES "extraction_schemas"("id") ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE "ocr_jobs" ADD CONSTRAINT "fk_ocr_jobs_user" FOREIGN KEY ("user_id") REFERENCES "user"("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "collections" ADD CONSTRAINT "fk_collections_user" FOREIGN KEY ("user_id") REFERENCES "user"("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "collection_schemas" ADD CONSTRAINT "fk_collection_schemas_collection" FOREIGN KEY ("collection_id") REFERENCES "collections"("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "collection_schemas" ADD CONSTRAINT "fk_collection_schemas_schema" FOREIGN KEY ("schema_id") REFERENCES "extraction_schemas"("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "collection_documents" ADD CONSTRAINT "fk_collection_documents_collection" FOREIGN KEY ("collection_id") REFERENCES "collections"("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "collection_documents" ADD CONSTRAINT "fk_collection_documents_document" FOREIGN KEY ("document_id") REFERENCES "ocr_documents"("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "account" ADD CONSTRAINT "fk_account_user" FOREIGN KEY ("user_id") REFERENCES "user"("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "session" ADD CONSTRAINT "fk_session_user" FOREIGN KEY ("user_id") REFERENCES "user"("id") ON DELETE CASCADE ON UPDATE CASCADE;
