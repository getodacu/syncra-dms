CREATE TABLE "billing_profiles" (
  "id" uuid,
  "user_id" uuid NOT NULL,
  "entity_type" varchar(40) NOT NULL,
  "billing_name" varchar(255) NOT NULL,
  "billing_email" varchar(320) NOT NULL,
  "country_code" varchar(2) NOT NULL,
  "address_line1" varchar(255) NOT NULL,
  "address_line2" varchar(255),
  "city" varchar(160) NOT NULL,
  "region" varchar(160),
  "postal_code" varchar(40) NOT NULL,
  "fiscal_code" varchar(80),
  "registration_number" varchar(120),
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "chk_billing_profiles_entity_type" CHECK ("entity_type" IN ('individual','company'))
);

CREATE UNIQUE INDEX "idx_billing_profiles_user_id" ON "billing_profiles" ("user_id");

ALTER TABLE "billing_profiles" ADD CONSTRAINT "fk_billing_profiles_user" FOREIGN KEY ("user_id") REFERENCES "user"("id") ON DELETE CASCADE ON UPDATE CASCADE;
