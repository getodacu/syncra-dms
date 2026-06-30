package fixtures

import (
	"strings"
	"testing"
)

func TestRecipeCategoriesLoadsFixture(t *testing.T) {
	categories, err := RecipeCategories()
	if err != nil {
		t.Fatalf("RecipeCategories() error = %v", err)
	}
	if len(categories) != 10 {
		t.Fatalf("category count = %d, want 10", len(categories))
	}
	if got := categories[0]; got.TitleEn != "Finance and Accounting" || got.TitleRo != "Finanțe și Contabilitate" {
		t.Fatalf("first category = %#v, want finance category", got)
	}
}

func TestParseRecipeCategoriesValidatesMalformedFixture(t *testing.T) {
	longTitle := strings.Repeat("x", 161)
	tests := []struct {
		name string
		data string
		want string
	}{
		{
			name: "invalid json",
			data: `{`,
			want: "parse recipe categories fixture",
		},
		{
			name: "wrong row length",
			data: `[["Finance"]]`,
			want: "must contain english and romanian titles",
		},
		{
			name: "empty english title",
			data: `[[ " ", "Finanțe" ]]`,
			want: "empty english title",
		},
		{
			name: "empty romanian title",
			data: `[[ "Finance", " " ]]`,
			want: "empty romanian title",
		},
		{
			name: "long english title",
			data: `[[ "` + longTitle + `", "Finanțe" ]]`,
			want: "english title exceeds",
		},
		{
			name: "long romanian title",
			data: `[[ "Finance", "` + longTitle + `" ]]`,
			want: "romanian title exceeds",
		},
		{
			name: "duplicate english title",
			data: `[[ "Finance", "Finanțe" ], [ "Finance", "Contabilitate" ]]`,
			want: "duplicates english title",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseRecipeCategories([]byte(tt.data))
			if err == nil {
				t.Fatal("ParseRecipeCategories() returned nil error")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error = %q, want %q", err.Error(), tt.want)
			}
		})
	}
}
