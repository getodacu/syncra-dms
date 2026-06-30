import { describe, expect, it } from "vitest";

import type { OCRRecipeCategoryResponse, OCRRecipeResponse } from "./api";
import {
	ALL_CATEGORY_KEY,
	OTHERS_CATEGORY_KEY,
	buildCategoryFilterOptions,
	filterRecipesByCategory,
	groupRecipesByCategory,
	type RecipeCard
} from "./page-state";

const created_at = "2026-06-20T00:00:00Z";
const updated_at = "2026-06-20T00:00:00Z";

function category(id: string, en: string, ro: string): OCRRecipeCategoryResponse {
	return {
		id,
		title: { en, ro },
		created_at,
		updated_at
	};
}

function recipe(id: string, title: string, category: OCRRecipeCategoryResponse | null): RecipeCard {
	const response: OCRRecipeResponse = {
		id,
		title,
		description: "",
		json: { type: "object" },
		counter: 0,
		category_id: category?.id ?? null,
		category,
		created_at,
		updated_at
	};

	return {
		recipe: response,
		summary: {
			fieldCount: 0,
			requiredCount: 0,
			fields: [],
			searchText: title.toLocaleLowerCase(),
			prettyJson: "{}"
		}
	};
}

const invoices = category("invoices", "Invoices", "Facturi");
const identity = category("identity", "Identity", "Identitate");
const cards = [
	recipe("invoice-1", "Invoice", invoices),
	recipe("identity-1", "Carte de identitate", identity),
	recipe("invoice-2", "Receipt", invoices),
	recipe("custom-1", "Custom document", null)
];
const categoryTitle = (item: OCRRecipeResponse) => item.category?.title.en ?? "Others";

describe("OCR recipe library state", () => {
	it("builds category filter options from loaded recipes with All first and Others last", () => {
		const options = buildCategoryFilterOptions(cards, categoryTitle, "All categories", "en");

		expect(options.map((option) => [option.key, option.title, option.count])).toEqual([
			[ALL_CATEGORY_KEY, "All categories", 4],
			["identity", "Identity", 1],
			["invoices", "Invoices", 2],
			[OTHERS_CATEGORY_KEY, "Others", 1]
		]);
		expect(options[0]).toMatchObject({ isAll: true, isOthers: false });
		expect(options.at(-1)).toMatchObject({ isAll: false, isOthers: true });
	});

	it("filters recipes by one selected category", () => {
		expect(filterRecipesByCategory(cards, ALL_CATEGORY_KEY).map((item) => item.recipe.id)).toEqual([
			"invoice-1",
			"identity-1",
			"invoice-2",
			"custom-1"
		]);
		expect(filterRecipesByCategory(cards, "invoices").map((item) => item.recipe.id)).toEqual([
			"invoice-1",
			"invoice-2"
		]);
		expect(filterRecipesByCategory(cards, OTHERS_CATEGORY_KEY).map((item) => item.recipe.id)).toEqual([
			"custom-1"
		]);
	});

	it("groups recipes by localized category with Others last", () => {
		const groups = groupRecipesByCategory(cards, categoryTitle, "en");

		expect(groups.map((group) => [group.key, group.title, group.items.length])).toEqual([
			["identity", "Identity", 1],
			["invoices", "Invoices", 2],
			[OTHERS_CATEGORY_KEY, "Others", 1]
		]);
	});
});
