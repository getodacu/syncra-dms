CREATE TABLE "json_recipe_categories" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "title_en" varchar(160) NOT NULL,
  "title_ro" varchar(160) NOT NULL,
  PRIMARY KEY ("id")
);

ALTER TABLE "json_recipes" ADD COLUMN "category_id" uuid NULL;

CREATE INDEX "idx_json_recipes_category_id" ON "json_recipes" ("category_id");

ALTER TABLE "json_recipes" ADD CONSTRAINT "fk_json_recipes_category" FOREIGN KEY ("category_id") REFERENCES "json_recipe_categories"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
