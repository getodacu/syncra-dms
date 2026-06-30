package fixtures

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"testing/fstest"
)

func TestRecipesLoadsFixture(t *testing.T) {
	recipes, err := Recipes()
	if err != nil {
		t.Fatalf("Recipes() error = %v", err)
	}
	if len(recipes) != 4 {
		t.Fatalf("recipe count = %d, want 4", len(recipes))
	}

	var foundBonFiscal bool
	for _, recipe := range recipes {
		if !json.Valid(recipe.JSON) {
			t.Fatalf("recipe %q JSON is invalid", recipe.Title)
		}
		if recipe.Title == "Bon fiscal" {
			foundBonFiscal = true
			if recipe.Category != "Finance and Accounting" {
				t.Fatalf("Bon fiscal category = %q, want Finance and Accounting", recipe.Category)
			}
		}
	}
	if !foundBonFiscal {
		t.Fatal("Bon fiscal fixture was not loaded")
	}
}

func TestParseRecipesTrimsFieldsAndPreservesJSON(t *testing.T) {
	recipes, err := ParseRecipes(recipeTestFS(map[string]string{
		"recipes/one.json": recipeFixture(" Invoice ", " Invoice fields ", " Finance and Accounting ", `{"type":"object"}`),
	}), testRecipeCategories())
	if err != nil {
		t.Fatalf("ParseRecipes() error = %v", err)
	}
	if len(recipes) != 1 {
		t.Fatalf("recipe count = %d, want 1", len(recipes))
	}
	got := recipes[0]
	if got.Title != "Invoice" || got.Description != "Invoice fields" || got.Category != "Finance and Accounting" {
		t.Fatalf("recipe = %#v, want trimmed fields", got)
	}
	if string(got.JSON) != `{"type":"object"}` {
		t.Fatalf("recipe JSON = %s, want original object bytes", got.JSON)
	}
}

func TestParseRecipesValidatesMalformedFixture(t *testing.T) {
	longTitle := strings.Repeat("x", 161)
	tests := []struct {
		name  string
		files map[string]string
		want  string
	}{
		{
			name: "invalid json",
			files: map[string]string{
				"recipes/one.json": `{`,
			},
			want: "parse recipe fixture",
		},
		{
			name: "empty title",
			files: map[string]string{
				"recipes/one.json": recipeFixture(" ", "Description", "Finance and Accounting", `{"type":"object"}`),
			},
			want: "title is required",
		},
		{
			name: "long title",
			files: map[string]string{
				"recipes/one.json": recipeFixture(longTitle, "Description", "Finance and Accounting", `{"type":"object"}`),
			},
			want: "title exceeds",
		},
		{
			name: "empty description",
			files: map[string]string{
				"recipes/one.json": recipeFixture("Invoice", " ", "Finance and Accounting", `{"type":"object"}`),
			},
			want: "description is required",
		},
		{
			name: "empty category",
			files: map[string]string{
				"recipes/one.json": recipeFixture("Invoice", "Description", " ", `{"type":"object"}`),
			},
			want: "category is required",
		},
		{
			name: "unknown category",
			files: map[string]string{
				"recipes/one.json": recipeFixture("Invoice", "Description", "Unknown", `{"type":"object"}`),
			},
			want: "unknown category",
		},
		{
			name: "non object json",
			files: map[string]string{
				"recipes/one.json": recipeFixture("Invoice", "Description", "Finance and Accounting", `[]`),
			},
			want: "json must be a JSON object",
		},
		{
			name: "duplicate title",
			files: map[string]string{
				"recipes/one.json": recipeFixture("Invoice", "Description", "Finance and Accounting", `{"type":"object"}`),
				"recipes/two.json": recipeFixture("Invoice", "Other description", "Finance and Accounting", `{"type":"object"}`),
			},
			want: "duplicates title",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseRecipes(recipeTestFS(tt.files), testRecipeCategories())
			if err == nil {
				t.Fatal("ParseRecipes() returned nil error")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error = %q, want %q", err.Error(), tt.want)
			}
		})
	}
}

func recipeFixture(title string, description string, category string, rawJSON string) string {
	return fmt.Sprintf(`{"title":%q,"description":%q,"category":%q,"json":%s}`, title, description, category, rawJSON)
}

func recipeTestFS(files map[string]string) fstest.MapFS {
	fsys := make(fstest.MapFS, len(files))
	for path, data := range files {
		fsys[path] = &fstest.MapFile{Data: []byte(data)}
	}
	return fsys
}

func testRecipeCategories() []RecipeCategory {
	return []RecipeCategory{
		{TitleEn: "Finance and Accounting", TitleRo: "Finance"},
	}
}
