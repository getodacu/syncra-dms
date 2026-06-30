package fixtures

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf8"
)

const maxRecipeCategoryTitleRunes = 160

//go:embed recipe_categories.json
var recipeCategoriesJSON []byte

type RecipeCategory struct {
	TitleEn string
	TitleRo string
}

func RecipeCategories() ([]RecipeCategory, error) {
	return ParseRecipeCategories(recipeCategoriesJSON)
}

func ParseRecipeCategories(data []byte) ([]RecipeCategory, error) {
	var raw [][]string
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse recipe categories fixture: %w", err)
	}

	categories := make([]RecipeCategory, 0, len(raw))
	seenTitleEn := make(map[string]struct{}, len(raw))
	for i, row := range raw {
		if len(row) != 2 {
			return nil, fmt.Errorf("recipe category fixture row %d must contain english and romanian titles", i)
		}

		titleEn := strings.TrimSpace(row[0])
		titleRo := strings.TrimSpace(row[1])
		if titleEn == "" {
			return nil, fmt.Errorf("recipe category fixture row %d has empty english title", i)
		}
		if titleRo == "" {
			return nil, fmt.Errorf("recipe category fixture row %d has empty romanian title", i)
		}
		if utf8.RuneCountInString(titleEn) > maxRecipeCategoryTitleRunes {
			return nil, fmt.Errorf("recipe category fixture row %d english title exceeds %d characters", i, maxRecipeCategoryTitleRunes)
		}
		if utf8.RuneCountInString(titleRo) > maxRecipeCategoryTitleRunes {
			return nil, fmt.Errorf("recipe category fixture row %d romanian title exceeds %d characters", i, maxRecipeCategoryTitleRunes)
		}
		if _, ok := seenTitleEn[titleEn]; ok {
			return nil, fmt.Errorf("recipe category fixture row %d duplicates english title %q", i, titleEn)
		}
		seenTitleEn[titleEn] = struct{}{}

		categories = append(categories, RecipeCategory{
			TitleEn: titleEn,
			TitleRo: titleRo,
		})
	}

	return categories, nil
}
