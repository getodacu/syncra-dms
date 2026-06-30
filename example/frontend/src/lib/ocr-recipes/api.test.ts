import { describe, expect, it, vi } from "vitest";

import { buildOCRRecipesPath, deployOCRRecipe, fetchOCRRecipes } from "./api";

function jsonResponse(body: unknown, init?: ResponseInit) {
	return new Response(JSON.stringify(body), {
		headers: { "content-type": "application/json" },
		...init
	});
}

function recipeResponse() {
	const category = {
		id: "category-1",
		title: { en: "Invoices", ro: "Facturi" },
		created_at: "2026-06-19T00:00:00Z",
		updated_at: "2026-06-19T00:00:00Z"
	};

	return {
		id: "recipe-1",
		title: "Invoice",
		description: "Invoice fields",
		json: { type: "object" },
		counter: 2,
		category_id: category.id,
		category,
		created_at: "2026-06-20T00:00:00Z",
		updated_at: "2026-06-20T00:00:00Z"
	};
}

describe("OCR recipe API client", () => {
	it("builds public JSON recipe list paths with pagination query params", () => {
		expect(buildOCRRecipesPath()).toBe("/api/json-recipes");
		expect(buildOCRRecipesPath({ cursor: "cursor-1", size: 50, sort: "asc" })).toBe(
			"/api/json-recipes?cursor=cursor-1&size=50&sort=asc"
		);
	});

	it("fetches OCR recipes through the SvelteKit proxy", async () => {
		const body = { recipes: [recipeResponse()], next_cursor: "cursor-2" };
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body));

		await expect(fetchOCRRecipes(fetchFn, { size: 20, sort: "desc" })).resolves.toEqual(body);
		expect(fetchFn).toHaveBeenCalledWith("/api/json-recipes?size=20&sort=desc", {
			method: "GET"
		});
	});

	it("deploys OCR recipes through the SvelteKit proxy and redirects callers with schema ids", async () => {
		const body = {
			recipe: recipeResponse(),
			schema: {
				id: "schema-1",
				created_at: "2026-06-20T00:00:00Z",
				updated_at: "2026-06-20T00:00:00Z",
				user_id: "user-1",
				name: "Invoice",
				description: "Invoice fields",
				strict: true,
				schema: { type: "object" }
			}
		};
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body));

		await expect(deployOCRRecipe(fetchFn, "recipe/1", "user-1")).resolves.toEqual(body);
		expect(fetchFn).toHaveBeenCalledWith("/api/json-recipes/recipe%2F1/deploy", {
			method: "POST",
			headers: { "content-type": "application/json" },
			body: JSON.stringify({ user_id: "user-1" })
		});
	});
});
