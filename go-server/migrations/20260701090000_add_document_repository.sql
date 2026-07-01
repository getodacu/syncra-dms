CREATE TABLE "document_folders" (
  "id" uuid,
  "parent_id" uuid,
  "organization_unit_id" uuid NOT NULL,
  "name" varchar(160) NOT NULL,
  "description" text,
  "created_by_user_id" uuid NOT NULL,
  "updated_by_user_id" uuid,
  "deleted_at" timestamptz,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);

CREATE INDEX "idx_document_folders_deleted_at" ON "document_folders" ("deleted_at");
CREATE INDEX "idx_document_folders_updated_by_user_id" ON "document_folders" ("updated_by_user_id");
CREATE INDEX "idx_document_folders_created_by_user_id" ON "document_folders" ("created_by_user_id");
CREATE UNIQUE INDEX "idx_document_folders_root_name_unique" ON "document_folders" ("organization_unit_id", "name") WHERE "parent_id" IS NULL AND "deleted_at" IS NULL;
CREATE UNIQUE INDEX "idx_document_folders_child_name_unique" ON "document_folders" ("organization_unit_id", "parent_id", "name") WHERE "parent_id" IS NOT NULL AND "deleted_at" IS NULL;
CREATE UNIQUE INDEX "idx_document_folders_id_organization_unit_unique" ON "document_folders" ("id", "organization_unit_id");
CREATE INDEX "idx_document_folders_parent_name_id" ON "document_folders" ("organization_unit_id", "parent_id", "name", "id");

CREATE TABLE "documents" (
  "id" uuid,
  "folder_id" uuid NOT NULL,
  "organization_unit_id" uuid NOT NULL,
  "original_file_name" varchar(255) NOT NULL,
  "display_name" varchar(255) NOT NULL,
  "mime_type" varchar(255) NOT NULL,
  "extension" varchar(32),
  "size_bytes" bigint NOT NULL,
  "sha256_hash" char(64) NOT NULL,
  "storage_key" text NOT NULL,
  "created_by_user_id" uuid NOT NULL,
  "updated_by_user_id" uuid,
  "deleted_at" timestamptz,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "chk_documents_sha256_hash_lower_hex" CHECK (length(sha256_hash) = 64 AND replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(sha256_hash,'0',''),'1',''),'2',''),'3',''),'4',''),'5',''),'6',''),'7',''),'8',''),'9',''),'a',''),'b',''),'c',''),'d',''),'e',''),'f','') = ''),
  CONSTRAINT "chk_documents_size_bytes_non_negative" CHECK (size_bytes >= 0)
);

CREATE INDEX "idx_documents_deleted_at" ON "documents" ("deleted_at");
CREATE INDEX "idx_documents_updated_by_user_id" ON "documents" ("updated_by_user_id");
CREATE INDEX "idx_documents_created_by_user_id" ON "documents" ("created_by_user_id");
CREATE INDEX "idx_documents_organization_unit_id" ON "documents" ("organization_unit_id");
CREATE UNIQUE INDEX "idx_documents_active_folder_hash_unique" ON "documents" ("folder_id", "sha256_hash") WHERE "deleted_at" IS NULL;
CREATE INDEX "idx_documents_folder_id" ON "documents" ("folder_id");

ALTER TABLE "document_folders" ADD CONSTRAINT "fk_document_folders_created_by_user" FOREIGN KEY ("created_by_user_id") REFERENCES "user"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "document_folders" ADD CONSTRAINT "fk_document_folders_organization_unit" FOREIGN KEY ("organization_unit_id") REFERENCES "organization_units"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "document_folders" ADD CONSTRAINT "fk_document_folders_parent" FOREIGN KEY ("parent_id", "organization_unit_id") REFERENCES "document_folders"("id", "organization_unit_id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "document_folders" ADD CONSTRAINT "fk_document_folders_updated_by_user" FOREIGN KEY ("updated_by_user_id") REFERENCES "user"("id") ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE "documents" ADD CONSTRAINT "fk_documents_created_by_user" FOREIGN KEY ("created_by_user_id") REFERENCES "user"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "documents" ADD CONSTRAINT "fk_documents_folder" FOREIGN KEY ("folder_id", "organization_unit_id") REFERENCES "document_folders"("id", "organization_unit_id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "documents" ADD CONSTRAINT "fk_documents_organization_unit" FOREIGN KEY ("organization_unit_id") REFERENCES "organization_units"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "documents" ADD CONSTRAINT "fk_documents_updated_by_user" FOREIGN KEY ("updated_by_user_id") REFERENCES "user"("id") ON DELETE SET NULL ON UPDATE CASCADE;
