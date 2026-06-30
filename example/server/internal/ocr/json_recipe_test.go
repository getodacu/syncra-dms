package ocr

import (
	"testing"

	"gorm.io/datatypes"

	"ai.ro/syncra/internal/testsupport"
)

func TestJSONRecipeAutoMigrateAndPersist(t *testing.T) {
	db := testsupport.OpenPostgresTx(t, &JSONRecipeCategory{}, &JSONRecipe{})

	category := JSONRecipeCategory{
		TitleEn: "Invoices",
		TitleRo: "Facturi",
	}
	if err := db.Create(&category).Error; err != nil {
		t.Fatalf("create JSON recipe category: %v", err)
	}

	recipe := JSONRecipe{
		Title:       "Invoice",
		Description: "Invoice extraction fields",
		JSON:        datatypes.JSON([]byte(`{"type":"object","properties":{"number":{"type":"string"}}}`)),
		Counter:     3,
		CategoryID:  &category.ID,
	}
	if err := db.Create(&recipe).Error; err != nil {
		t.Fatalf("create JSON recipe: %v", err)
	}
	if recipe.ID.String() == "00000000-0000-0000-0000-000000000000" {
		t.Fatal("recipe ID was not generated")
	}

	var got JSONRecipe
	if err := db.First(&got, "id = ?", recipe.ID).Error; err != nil {
		t.Fatalf("load JSON recipe: %v", err)
	}
	if got.Title != recipe.Title || got.Description != recipe.Description || got.Counter != 3 {
		t.Fatalf("unexpected recipe: %#v", got)
	}
	if got.CategoryID == nil || *got.CategoryID != category.ID {
		t.Fatalf("category_id = %v, want %s", got.CategoryID, category.ID)
	}
	assertJSONEqualLocal(t, "recipe JSON", recipe.JSON, got.JSON)
}

func TestJSONRecipeCategoryAutoMigrateAndPersist(t *testing.T) {
	db := testsupport.OpenPostgresTx(t, &JSONRecipeCategory{})

	category := JSONRecipeCategory{
		TitleEn: "Receipts",
		TitleRo: "Bonuri",
	}
	if err := db.Create(&category).Error; err != nil {
		t.Fatalf("create JSON recipe category: %v", err)
	}
	if category.ID.String() == "00000000-0000-0000-0000-000000000000" {
		t.Fatal("category ID was not generated")
	}

	var got JSONRecipeCategory
	if err := db.First(&got, "id = ?", category.ID).Error; err != nil {
		t.Fatalf("load JSON recipe category: %v", err)
	}
	if got.TitleEn != category.TitleEn || got.TitleRo != category.TitleRo {
		t.Fatalf("unexpected category: %#v", got)
	}
}

func TestJSONRecipeTableName(t *testing.T) {
	if got := (JSONRecipe{}).TableName(); got != "json_recipes" {
		t.Fatalf("table name = %q, want json_recipes", got)
	}
}

func TestJSONRecipeCategoryTableName(t *testing.T) {
	if got := (JSONRecipeCategory{}).TableName(); got != "json_recipe_categories" {
		t.Fatalf("table name = %q, want json_recipe_categories", got)
	}
}
