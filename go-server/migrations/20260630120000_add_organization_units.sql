CREATE TABLE "organization_units" (
  "id" uuid,
  "parent_id" uuid,
  "name" varchar(160) NOT NULL,
  "code" varchar(40),
  "description" text,
  "archived_at" timestamptz,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);

CREATE INDEX "idx_organization_units_parent_name_id" ON "organization_units" ("parent_id", "name", "id");
CREATE INDEX "idx_organization_units_archived_at" ON "organization_units" ("archived_at");
CREATE INDEX "idx_organization_units_code" ON "organization_units" ("code");
CREATE UNIQUE INDEX "idx_organization_units_active_code_unique" ON "organization_units" ("code") WHERE "code" IS NOT NULL AND "archived_at" IS NULL;

ALTER TABLE "organization_units" ADD CONSTRAINT "fk_organization_units_parent" FOREIGN KEY ("parent_id") REFERENCES "organization_units"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
