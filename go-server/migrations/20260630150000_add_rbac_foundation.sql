ALTER TABLE "user" ADD COLUMN "status" varchar(40);
ALTER TABLE "user" ADD COLUMN "primary_organization_unit_id" uuid;
ALTER TABLE "user" ADD COLUMN "manager_user_id" uuid;
ALTER TABLE "user" ADD COLUMN "job_title" varchar(160);
ALTER TABLE "user" ADD COLUMN "phone" varchar(80);
ALTER TABLE "user" ADD COLUMN "deleted_at" timestamptz;

UPDATE "user" SET "status" = CASE WHEN "email_verified" THEN 'active' ELSE 'invited' END WHERE "status" IS NULL;

ALTER TABLE "user" ALTER COLUMN "status" SET DEFAULT 'active';
ALTER TABLE "user" ALTER COLUMN "status" SET NOT NULL;
ALTER TABLE "user" ADD CONSTRAINT "chk_user_status" CHECK (status IN ('invited','active','inactive','suspended','deleted'));

CREATE INDEX "idx_user_status" ON "user" ("status");
CREATE INDEX "idx_user_primary_organization_unit_id" ON "user" ("primary_organization_unit_id");
CREATE INDEX "idx_user_manager_user_id" ON "user" ("manager_user_id");
CREATE INDEX "idx_user_deleted_at" ON "user" ("deleted_at");

CREATE TABLE "roles" (
  "id" uuid,
  "name" varchar(160) NOT NULL,
  "code" varchar(80) NOT NULL,
  "description" text,
  "is_system" boolean NOT NULL DEFAULT false,
  "is_active" boolean NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX "idx_roles_code" ON "roles" ("code");
CREATE INDEX "idx_roles_is_active" ON "roles" ("is_active");

CREATE TABLE "permissions" (
  "id" uuid,
  "code" varchar(120) NOT NULL,
  "name" varchar(160) NOT NULL,
  "description" text,
  "category" varchar(80) NOT NULL,
  "is_system" boolean NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX "idx_permissions_code" ON "permissions" ("code");
CREATE INDEX "idx_permissions_category" ON "permissions" ("category");

CREATE TABLE "role_permissions" (
  "id" uuid,
  "role_id" uuid NOT NULL,
  "permission_id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX "idx_role_permission_unique" ON "role_permissions" ("role_id", "permission_id");

CREATE TABLE "user_roles" (
  "id" uuid,
  "user_id" uuid NOT NULL,
  "role_id" uuid NOT NULL,
  "scope_type" varchar(40) NOT NULL,
  "organization_unit_id" uuid,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "chk_user_roles_scope_type" CHECK (scope_type IN ('global','organization_unit','organization_unit_and_children'))
);

CREATE INDEX "idx_user_roles_user_id" ON "user_roles" ("user_id");
CREATE INDEX "idx_user_roles_role_id" ON "user_roles" ("role_id");
CREATE UNIQUE INDEX "idx_user_role_scope_unique" ON "user_roles" ("user_id", "role_id", "scope_type", "organization_unit_id") NULLS NOT DISTINCT;

CREATE TABLE "groups" (
  "id" uuid,
  "name" varchar(160) NOT NULL,
  "code" varchar(80) NOT NULL,
  "description" text,
  "organization_unit_id" uuid,
  "is_active" boolean NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX "idx_groups_code" ON "groups" ("code");
CREATE INDEX "idx_groups_organization_unit_id" ON "groups" ("organization_unit_id");
CREATE INDEX "idx_groups_is_active" ON "groups" ("is_active");

CREATE TABLE "group_users" (
  "id" uuid,
  "group_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX "idx_group_user_unique" ON "group_users" ("group_id", "user_id");

CREATE TABLE "group_roles" (
  "id" uuid,
  "group_id" uuid NOT NULL,
  "role_id" uuid NOT NULL,
  "scope_type" varchar(40) NOT NULL,
  "organization_unit_id" uuid,
  "created_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "chk_group_roles_scope_type" CHECK (scope_type IN ('global','organization_unit','organization_unit_and_children'))
);

CREATE UNIQUE INDEX "idx_group_role_scope_unique" ON "group_roles" ("group_id", "role_id", "scope_type", "organization_unit_id") NULLS NOT DISTINCT;

CREATE TABLE "organization_unit_roles" (
  "id" uuid,
  "organization_unit_id" uuid NOT NULL,
  "role_id" uuid NOT NULL,
  "scope_type" varchar(40) NOT NULL,
  "created_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "chk_organization_unit_roles_scope_type" CHECK (scope_type IN ('global','organization_unit','organization_unit_and_children'))
);

CREATE UNIQUE INDEX "idx_organization_unit_role_unique" ON "organization_unit_roles" ("organization_unit_id", "role_id", "scope_type");

ALTER TABLE "user" ADD CONSTRAINT "fk_user_primary_organization_unit" FOREIGN KEY ("primary_organization_unit_id") REFERENCES "organization_units"("id") ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE "user" ADD CONSTRAINT "fk_user_manager_user" FOREIGN KEY ("manager_user_id") REFERENCES "user"("id") ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE "role_permissions" ADD CONSTRAINT "fk_role_permissions_role" FOREIGN KEY ("role_id") REFERENCES "roles"("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "role_permissions" ADD CONSTRAINT "fk_role_permissions_permission" FOREIGN KEY ("permission_id") REFERENCES "permissions"("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "user_roles" ADD CONSTRAINT "fk_user_roles_user" FOREIGN KEY ("user_id") REFERENCES "user"("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "user_roles" ADD CONSTRAINT "fk_user_roles_role" FOREIGN KEY ("role_id") REFERENCES "roles"("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "user_roles" ADD CONSTRAINT "fk_user_roles_organization_unit" FOREIGN KEY ("organization_unit_id") REFERENCES "organization_units"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "groups" ADD CONSTRAINT "fk_groups_organization_unit" FOREIGN KEY ("organization_unit_id") REFERENCES "organization_units"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "group_users" ADD CONSTRAINT "fk_group_users_group" FOREIGN KEY ("group_id") REFERENCES "groups"("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "group_users" ADD CONSTRAINT "fk_group_users_user" FOREIGN KEY ("user_id") REFERENCES "user"("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "group_roles" ADD CONSTRAINT "fk_group_roles_group" FOREIGN KEY ("group_id") REFERENCES "groups"("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "group_roles" ADD CONSTRAINT "fk_group_roles_role" FOREIGN KEY ("role_id") REFERENCES "roles"("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "group_roles" ADD CONSTRAINT "fk_group_roles_organization_unit" FOREIGN KEY ("organization_unit_id") REFERENCES "organization_units"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "organization_unit_roles" ADD CONSTRAINT "fk_organization_unit_roles_organization_unit" FOREIGN KEY ("organization_unit_id") REFERENCES "organization_units"("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "organization_unit_roles" ADD CONSTRAINT "fk_organization_unit_roles_role" FOREIGN KEY ("role_id") REFERENCES "roles"("id") ON DELETE CASCADE ON UPDATE CASCADE;
