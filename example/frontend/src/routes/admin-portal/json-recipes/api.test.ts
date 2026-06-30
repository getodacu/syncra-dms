import { describe, expect, it, vi } from "vitest";

import {
	buildAdminJSONRecipeCategoriesPath,
	buildAdminJSONRecipeCategoryDetailPath,
	buildAdminJSONRecipeDetailPath,
	buildAdminJSONRecipesPath,
	createAdminJSONRecipeCategory,
	createAdminJSONRecipe,
	deleteAdminJSONRecipeCategory,
	deleteAdminJSONRecipe,
	fetchAdminJSONRecipeCategories,
	fetchAdminJSONRecipe,
	fetchAdminJSONRecipes,
	updateAdminJSONRecipeCategory,
	updateAdminJSONRecipe
} from "./api";

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
		counter: 3,
		category_id: category.id,
		category,
		created_at: "2026-06-20T00:00:00Z",
		updated_at: "2026-06-21T00:00:00Z"
	};
}

function categoryResponse() {
	return {
		id: "category-1",
		title: { en: "Invoices", ro: "Facturi" },
		created_at: "2026-06-19T00:00:00Z",
		updated_at: "2026-06-19T00:00:00Z"
	};
}

describe("admin JSON recipes client API", () => {
	it("builds list and detail paths", () => {
		expect(buildAdminJSONRecipesPath()).toBe("/api/admin/json-recipes");
		expect(buildAdminJSONRecipesPath({ cursor: "cursor-1", size: 50, sort: "asc" })).toBe(
			"/api/admin/json-recipes?cursor=cursor-1&size=50&sort=asc"
		);
		expect(buildAdminJSONRecipeDetailPath("recipe 1")).toBe("/api/admin/json-recipes/recipe%201");
		expect(buildAdminJSONRecipeCategoriesPath()).toBe("/api/admin/json-recipe-categories");
		expect(buildAdminJSONRecipeCategoryDetailPath("category 1")).toBe(
			"/api/admin/json-recipe-categories/category%201"
		);
	});

	it("fetches recipe lists", async () => {
		const result = { recipes: [recipeResponse()], next_cursor: null };
		const fetchMock = vi.fn().mockResolvedValue(jsonResponse(result));

		await expect(fetchAdminJSONRecipes(fetchMock, { size: 20, sort: "desc" })).resolves.toEqual(result);
		expect(fetchMock).toHaveBeenCalledWith("/api/admin/json-recipes?size=20&sort=desc", {
			method: "GET"
		});
	});

	it("creates, updates, gets, and deletes recipes", async () => {
		const recipe = recipeResponse();
		const input = {
			title: "Invoice",
			description: "Invoice fields",
			json: { type: "object" },
			category_id: "category-1"
		};
		const fetchMock = vi.fn(async (_input: Parameters<typeof fetch>[0], init?: Parameters<typeof fetch>[1]) => {
			if (init?.method === "DELETE") return new Response(null, { status: 204 });
			return jsonResponse(recipe, { status: init?.method === "POST" ? 201 : 200 });
		});

		await expect(fetchAdminJSONRecipe(fetchMock, "recipe-1")).resolves.toEqual(recipe);
		await expect(createAdminJSONRecipe(fetchMock, input)).resolves.toEqual(recipe);
		await expect(updateAdminJSONRecipe(fetchMock, "recipe-1", input)).resolves.toEqual(recipe);
		await expect(deleteAdminJSONRecipe(fetchMock, "recipe-1")).resolves.toBeUndefined();

		expect(fetchMock).toHaveBeenNthCalledWith(1, "/api/admin/json-recipes/recipe-1", { method: "GET" });
		expect(fetchMock).toHaveBeenNthCalledWith(
			2,
			"/api/admin/json-recipes",
			expect.objectContaining({
				method: "POST",
				body: JSON.stringify(input)
			})
		);
		expect(fetchMock).toHaveBeenNthCalledWith(
			3,
			"/api/admin/json-recipes/recipe-1",
			expect.objectContaining({
				method: "PUT",
				body: JSON.stringify(input)
			})
		);
		expect(fetchMock).toHaveBeenNthCalledWith(4, "/api/admin/json-recipes/recipe-1", {
			method: "DELETE"
		});
	});

	it("creates, updates, gets, lists, and deletes categories", async () => {
		const category = categoryResponse();
		const list = { categories: [category] };
		const input = { title: { en: "Invoices", ro: "Facturi" } };
		const fetchMock = vi.fn(async (request: Parameters<typeof fetch>[0], init?: Parameters<typeof fetch>[1]) => {
			const url = String(request);
			if (init?.method === "GET" && url.endsWith("/json-recipe-categories")) return jsonResponse(list);
			if (init?.method === "DELETE") return new Response(null, { status: 204 });
			return jsonResponse(category, { status: init?.method === "POST" ? 201 : 200 });
		});

		await expect(fetchAdminJSONRecipeCategories(fetchMock)).resolves.toEqual(list);
		await expect(createAdminJSONRecipeCategory(fetchMock, input)).resolves.toEqual(category);
		await expect(fetchAdminJSONRecipeCategories(fetchMock)).resolves.toEqual(list);
		await expect(updateAdminJSONRecipeCategory(fetchMock, "category-1", input)).resolves.toEqual(category);
		await expect(deleteAdminJSONRecipeCategory(fetchMock, "category-1")).resolves.toBeUndefined();

		expect(fetchMock).toHaveBeenNthCalledWith(1, "/api/admin/json-recipe-categories", {
			method: "GET"
		});
		expect(fetchMock).toHaveBeenNthCalledWith(
			2,
			"/api/admin/json-recipe-categories",
			expect.objectContaining({ method: "POST", body: JSON.stringify(input) })
		);
		expect(fetchMock).toHaveBeenNthCalledWith(3, "/api/admin/json-recipe-categories", {
			method: "GET"
		});
		expect(fetchMock).toHaveBeenNthCalledWith(
			4,
			"/api/admin/json-recipe-categories/category-1",
			expect.objectContaining({ method: "PUT", body: JSON.stringify(input) })
		);
		expect(fetchMock).toHaveBeenNthCalledWith(5, "/api/admin/json-recipe-categories/category-1", {
			method: "DELETE"
		});
	});

	it("throws backend JSON errors and rejects invalid responses", async () => {
		const backendErrorFetch = vi.fn().mockResolvedValue(jsonResponse({ error: "title is required" }, { status: 400 }));
		const invalidFetch = vi.fn().mockResolvedValue(jsonResponse({ recipes: [{ ...recipeResponse(), counter: "3" }], next_cursor: null }));
		const invalidCategoryFetch = vi.fn().mockResolvedValue(jsonResponse({ categories: [{ ...categoryResponse(), title: "Invoices" }] }));

		await expect(fetchAdminJSONRecipes(backendErrorFetch, {})).rejects.toThrow("title is required");
		await expect(fetchAdminJSONRecipes(invalidFetch, {})).rejects.toThrow(
			"Invalid admin JSON recipes response"
		);
		await expect(fetchAdminJSONRecipeCategories(invalidCategoryFetch)).rejects.toThrow(
			"Invalid admin JSON recipe categories response"
		);
	});
});
