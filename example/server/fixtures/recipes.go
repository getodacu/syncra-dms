package fixtures

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"strings"
	"unicode/utf8"
)

const maxRecipeTitleRunes = 160

//go:embed recipes/*.json
var recipeFiles embed.FS

type Recipe struct {
	Title       string
	Description string
	Category    string
	JSON        json.RawMessage
}

type rawRecipe struct {
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Category    string          `json:"category"`
	JSON        json.RawMessage `json:"json"`
}

func Recipes() ([]Recipe, error) {
	categories, err := RecipeCategories()
	if err != nil {
		return nil, err
	}
	return ParseRecipes(recipeFiles, categories)
}

func ParseRecipes(fsys fs.FS, categories []RecipeCategory) ([]Recipe, error) {
	paths, err := fs.Glob(fsys, "recipes/*.json")
	if err != nil {
		return nil, fmt.Errorf("list recipe fixtures: %w", err)
	}

	knownCategories := make(map[string]struct{}, len(categories))
	for _, category := range categories {
		knownCategories[category.TitleEn] = struct{}{}
	}

	recipes := make([]Recipe, 0, len(paths))
	seenTitles := make(map[string]struct{}, len(paths))
	for _, path := range paths {
		data, err := fs.ReadFile(fsys, path)
		if err != nil {
			return nil, fmt.Errorf("read recipe fixture %s: %w", path, err)
		}

		recipe, err := parseRecipe(path, data, knownCategories)
		if err != nil {
			return nil, err
		}
		if _, ok := seenTitles[recipe.Title]; ok {
			return nil, fmt.Errorf("recipe fixture %s duplicates title %q", path, recipe.Title)
		}
		seenTitles[recipe.Title] = struct{}{}
		recipes = append(recipes, recipe)
	}

	return recipes, nil
}

func parseRecipe(path string, data []byte, knownCategories map[string]struct{}) (Recipe, error) {
	var raw rawRecipe
	if err := json.Unmarshal(data, &raw); err != nil {
		return Recipe{}, fmt.Errorf("parse recipe fixture %s: %w", path, err)
	}

	title := strings.TrimSpace(raw.Title)
	description := strings.TrimSpace(raw.Description)
	category := strings.TrimSpace(raw.Category)
	if title == "" {
		return Recipe{}, fmt.Errorf("recipe fixture %s title is required", path)
	}
	if utf8.RuneCountInString(title) > maxRecipeTitleRunes {
		return Recipe{}, fmt.Errorf("recipe fixture %s title exceeds %d characters", path, maxRecipeTitleRunes)
	}
	if description == "" {
		return Recipe{}, fmt.Errorf("recipe fixture %s description is required", path)
	}
	if category == "" {
		return Recipe{}, fmt.Errorf("recipe fixture %s category is required", path)
	}
	if _, ok := knownCategories[category]; !ok {
		return Recipe{}, fmt.Errorf("recipe fixture %s references unknown category %q", path, category)
	}
	if !isRecipeJSONObject(raw.JSON) {
		return Recipe{}, fmt.Errorf("recipe fixture %s json must be a JSON object", path)
	}

	return Recipe{
		Title:       title,
		Description: description,
		Category:    category,
		JSON:        raw.JSON,
	}, nil
}

func isRecipeJSONObject(raw json.RawMessage) bool {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 || trimmed[0] != '{' {
		return false
	}
	var obj map[string]any
	return json.Unmarshal(trimmed, &obj) == nil
}
