import type { OCRRecipeResponse } from "./api";
import type { RecipeSummary } from "./recipe-summary";

export const ALL_CATEGORY_KEY = "__all__";
export const OTHERS_CATEGORY_KEY = "__others__";

export type RecipeCard = {
	recipe: OCRRecipeResponse;
	summary: RecipeSummary;
};

export type RecipeGroup = {
	key: string;
	title: string;
	isOthers: boolean;
	items: RecipeCard[];
};

export type CategoryFilterOption = {
	key: string;
	title: string;
	count: number;
	isAll: boolean;
	isOthers: boolean;
};

type CategoryTitle = (recipe: OCRRecipeResponse) => string;

export function categoryKey(recipe: OCRRecipeResponse) {
	return recipe.category_id ?? OTHERS_CATEGORY_KEY;
}

export function filterRecipesByCategory(items: RecipeCard[], selectedCategoryKey: string) {
	if (selectedCategoryKey === ALL_CATEGORY_KEY) return items;
	return items.filter((item) => categoryKey(item.recipe) === selectedCategoryKey);
}

export function buildCategoryFilterOptions(
	items: RecipeCard[],
	titleForRecipe: CategoryTitle,
	allTitle: string,
	locale?: string
) {
	const options = new Map<string, CategoryFilterOption>();
	for (const item of items) {
		const key = categoryKey(item.recipe);
		const existing = options.get(key);
		if (existing) {
			existing.count += 1;
			continue;
		}
		options.set(key, {
			key,
			title: titleForRecipe(item.recipe),
			count: 1,
			isAll: false,
			isOthers: key === OTHERS_CATEGORY_KEY
		});
	}

	return [
		{
			key: ALL_CATEGORY_KEY,
			title: allTitle,
			count: items.length,
			isAll: true,
			isOthers: false
		},
		...[...options.values()].sort((left, right) => compareCategoryOptions(left, right, locale))
	];
}

export function groupRecipesByCategory(items: RecipeCard[], titleForRecipe: CategoryTitle, locale?: string) {
	const groups = new Map<string, RecipeGroup>();
	for (const item of items) {
		const key = categoryKey(item.recipe);
		const existing = groups.get(key);
		if (existing) {
			existing.items.push(item);
			continue;
		}
		groups.set(key, {
			key,
			title: titleForRecipe(item.recipe),
			isOthers: key === OTHERS_CATEGORY_KEY,
			items: [item]
		});
	}
	return [...groups.values()].sort((left, right) => compareCategoryOptions(left, right, locale));
}

function compareCategoryOptions(
	left: { title: string; isOthers: boolean },
	right: { title: string; isOthers: boolean },
	locale?: string
) {
	if (left.isOthers && !right.isOthers) return 1;
	if (!left.isOthers && right.isOthers) return -1;
	return left.title.localeCompare(right.title, locale);
}
