import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";

const source = () => readFileSync(new URL("./ocr-recipe-library.svelte", import.meta.url), "utf8");

function normalize(value: string) {
	return value.replace(/\s+/g, " ");
}

describe("OCR recipe library markup", () => {
	it("renders recipes as one responsive grid with category labels on each card", () => {
		const component = source();
		const normalized = normalize(component);

		expect(component).not.toContain("groupRecipesByCategory");
		expect(component).not.toContain("visibleRecipeGroups");
		expect(normalized).toContain('{#each visibleRecipes as item (item.recipe.id)}');
		expect(normalized).toContain(
			'<div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">'
		);
		expect(normalized).toContain("{categoryTitle(item.recipe)}");
		expect(normalized).not.toContain("{group.title}");
	});
});
