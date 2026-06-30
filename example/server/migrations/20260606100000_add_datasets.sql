CREATE TABLE "datasets" (
  "id" uuid,
  "user_id" uuid NOT NULL,
  "schema_id" uuid NOT NULL,
  "created_at" timestamptz,
  "updated_at" timestamptz,
  "name" varchar(160) NOT NULL,
  "selected_fields" jsonb NOT NULL,
  PRIMARY KEY ("id")
);

CREATE INDEX "idx_datasets_schema_id" ON "datasets" ("schema_id");
CREATE INDEX "idx_datasets_user_id" ON "datasets" ("user_id");
CREATE INDEX "idx_datasets_user_created_id" ON "datasets" ("user_id", "created_at", "id");

ALTER TABLE "datasets" ADD CONSTRAINT "fk_datasets_user" FOREIGN KEY ("user_id") REFERENCES "user"("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "datasets" ADD CONSTRAINT "fk_datasets_schema" FOREIGN KEY ("schema_id") REFERENCES "extraction_schemas"("id") ON DELETE CASCADE ON UPDATE CASCADE;
