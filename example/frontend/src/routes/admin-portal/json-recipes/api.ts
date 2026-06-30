import { publicApiErrorMessage } from "$lib/client/api-errors";

export type AdminJSONRecipeResponse = {
	id: string;
	title: string;
	description: string;
	json: Record<string, unknown>;
	counter: number;
	category_id: string | null;
	category: AdminJSONRecipeCategoryResponse | null;
	created_at: string;
	updated_at: string;
};

export type AdminJSONRecipeCategoryTitle = {
	en: string;
	ro: string;
};

export type AdminJSONRecipeCategoryResponse = {
	id: string;
	title: AdminJSONRecipeCategoryTitle;
	created_at: string;
	updated_at: string;
};

export type AdminJSONRecipeCategoryListResponse = {
	categories: AdminJSONRecipeCategoryResponse[];
};

export type AdminJSONRecipeListResponse = {
	recipes: AdminJSONRecipeResponse[];
	next_cursor: string | null;
};

export type AdminJSONRecipeInput = {
	title: string;
	description: string;
	json: Record<string, unknown>;
	category_id?: string | null;
};

export type AdminJSONRecipeCategoryInput = {
	title: AdminJSONRecipeCategoryTitle;
};

export type AdminJSONRecipesQuery = {
	cursor?: string | null;
	size?: string | number;
	sort?: "asc" | "desc";
};

type ClientFetch = typeof fetch;

export function buildAdminJSONRecipesPath(query: AdminJSONRecipesQuery = {}) {
	const params = new URLSearchParams();
	setNonEmptyParam(params, "cursor", query.cursor);
	setNonEmptyParam(params, "size", query.size);
	setNonEmptyParam(params, "sort", query.sort);

	const queryString = params.toString();
	return queryString ? `/api/admin/json-recipes?${queryString}` : "/api/admin/json-recipes";
}

export function buildAdminJSONRecipeDetailPath(id: string) {
	return `/api/admin/json-recipes/${encodeURIComponent(id)}`;
}

export function buildAdminJSONRecipeCategoriesPath() {
	return "/api/admin/json-recipe-categories";
}

export function buildAdminJSONRecipeCategoryDetailPath(id: string) {
	return `/api/admin/json-recipe-categories/${encodeURIComponent(id)}`;
}

export async function fetchAdminJSONRecipes(
	fetchFn: ClientFetch,
	query: AdminJSONRecipesQuery
): Promise<AdminJSONRecipeListResponse> {
	const response = await fetchFn(buildAdminJSONRecipesPath(query), { method: "GET" });
	const body = await readResponseJSON(response);

	if (!response.ok) {
		throw new Error(publicApiErrorMessage(response.status, body, "Failed to load JSON recipes"));
	}
	if (!isAdminJSONRecipeListResponse(body)) {
		throw new Error("Invalid admin JSON recipes response");
	}
	return body;
}

export async function fetchAdminJSONRecipeCategories(
	fetchFn: ClientFetch
): Promise<AdminJSONRecipeCategoryListResponse> {
	const response = await fetchFn(buildAdminJSONRecipeCategoriesPath(), { method: "GET" });
	const body = await readResponseJSON(response);

	if (!response.ok) {
		throw new Error(publicApiErrorMessage(response.status, body, "Failed to load JSON recipe categories"));
	}
	if (!isAdminJSONRecipeCategoryListResponse(body)) {
		throw new Error("Invalid admin JSON recipe categories response");
	}
	return body;
}

export async function fetchAdminJSONRecipe(
	fetchFn: ClientFetch,
	id: string
): Promise<AdminJSONRecipeResponse> {
	const response = await fetchFn(buildAdminJSONRecipeDetailPath(id), { method: "GET" });
	const body = await readResponseJSON(response);

	if (!response.ok) {
		throw new Error(publicApiErrorMessage(response.status, body, "Failed to load JSON recipe"));
	}
	if (!isAdminJSONRecipeResponse(body)) {
		throw new Error("Invalid admin JSON recipe response");
	}
	return body;
}

export async function fetchAdminJSONRecipeCategory(
	fetchFn: ClientFetch,
	id: string
): Promise<AdminJSONRecipeCategoryResponse> {
	const response = await fetchFn(buildAdminJSONRecipeCategoryDetailPath(id), { method: "GET" });
	const body = await readResponseJSON(response);

	if (!response.ok) {
		throw new Error(publicApiErrorMessage(response.status, body, "Failed to load JSON recipe category"));
	}
	if (!isAdminJSONRecipeCategoryResponse(body)) {
		throw new Error("Invalid admin JSON recipe category response");
	}
	return body;
}

export async function createAdminJSONRecipe(
	fetchFn: ClientFetch,
	input: AdminJSONRecipeInput
): Promise<AdminJSONRecipeResponse> {
	const response = await fetchFn("/api/admin/json-recipes", {
		method: "POST",
		headers: { "content-type": "application/json" },
		body: JSON.stringify(input)
	});
	const body = await readResponseJSON(response);

	if (!response.ok) {
		throw new Error(publicApiErrorMessage(response.status, body, "Failed to create JSON recipe"));
	}
	if (!isAdminJSONRecipeResponse(body)) {
		throw new Error("Invalid admin JSON recipe response");
	}
	return body;
}

export async function createAdminJSONRecipeCategory(
	fetchFn: ClientFetch,
	input: AdminJSONRecipeCategoryInput
): Promise<AdminJSONRecipeCategoryResponse> {
	const response = await fetchFn(buildAdminJSONRecipeCategoriesPath(), {
		method: "POST",
		headers: { "content-type": "application/json" },
		body: JSON.stringify(input)
	});
	const body = await readResponseJSON(response);

	if (!response.ok) {
		throw new Error(publicApiErrorMessage(response.status, body, "Failed to create JSON recipe category"));
	}
	if (!isAdminJSONRecipeCategoryResponse(body)) {
		throw new Error("Invalid admin JSON recipe category response");
	}
	return body;
}

export async function updateAdminJSONRecipe(
	fetchFn: ClientFetch,
	id: string,
	input: AdminJSONRecipeInput
): Promise<AdminJSONRecipeResponse> {
	const response = await fetchFn(buildAdminJSONRecipeDetailPath(id), {
		method: "PUT",
		headers: { "content-type": "application/json" },
		body: JSON.stringify(input)
	});
	const body = await readResponseJSON(response);

	if (!response.ok) {
		throw new Error(publicApiErrorMessage(response.status, body, "Failed to update JSON recipe"));
	}
	if (!isAdminJSONRecipeResponse(body)) {
		throw new Error("Invalid admin JSON recipe response");
	}
	return body;
}

export async function updateAdminJSONRecipeCategory(
	fetchFn: ClientFetch,
	id: string,
	input: AdminJSONRecipeCategoryInput
): Promise<AdminJSONRecipeCategoryResponse> {
	const response = await fetchFn(buildAdminJSONRecipeCategoryDetailPath(id), {
		method: "PUT",
		headers: { "content-type": "application/json" },
		body: JSON.stringify(input)
	});
	const body = await readResponseJSON(response);

	if (!response.ok) {
		throw new Error(publicApiErrorMessage(response.status, body, "Failed to update JSON recipe category"));
	}
	if (!isAdminJSONRecipeCategoryResponse(body)) {
		throw new Error("Invalid admin JSON recipe category response");
	}
	return body;
}

export async function deleteAdminJSONRecipe(fetchFn: ClientFetch, id: string): Promise<void> {
	const response = await fetchFn(buildAdminJSONRecipeDetailPath(id), { method: "DELETE" });
	const body = await readResponseJSON(response);

	if (!response.ok) {
		throw new Error(publicApiErrorMessage(response.status, body, "Failed to delete JSON recipe"));
	}
}

export async function deleteAdminJSONRecipeCategory(fetchFn: ClientFetch, id: string): Promise<void> {
	const response = await fetchFn(buildAdminJSONRecipeCategoryDetailPath(id), { method: "DELETE" });
	const body = await readResponseJSON(response);

	if (!response.ok) {
		throw new Error(publicApiErrorMessage(response.status, body, "Failed to delete JSON recipe category"));
	}
}

async function readResponseJSON(response: Response): Promise<unknown> {
	const text = await response.text();
	if (!text.trim()) return null;
	try {
		return JSON.parse(text);
	} catch {
		return null;
	}
}

function setNonEmptyParam(params: URLSearchParams, name: string, value: string | number | null | undefined) {
	if (value === undefined || value === null) return;
	const text = String(value).trim();
	if (text) params.set(name, text);
}

function isRecord(value: unknown): value is Record<string, unknown> {
	return typeof value === "object" && value !== null && !Array.isArray(value);
}

function isAdminJSONRecipeCategoryTitle(value: unknown): value is AdminJSONRecipeCategoryTitle {
	return isRecord(value) && typeof value.en === "string" && typeof value.ro === "string";
}

function isAdminJSONRecipeCategoryResponse(value: unknown): value is AdminJSONRecipeCategoryResponse {
	return (
		isRecord(value) &&
		typeof value.id === "string" &&
		isAdminJSONRecipeCategoryTitle(value.title) &&
		typeof value.created_at === "string" &&
		typeof value.updated_at === "string"
	);
}

function isAdminJSONRecipeResponse(value: unknown): value is AdminJSONRecipeResponse {
	return (
		isRecord(value) &&
		typeof value.id === "string" &&
		typeof value.title === "string" &&
		typeof value.description === "string" &&
		isRecord(value.json) &&
		typeof value.counter === "number" &&
		Number.isFinite(value.counter) &&
		(typeof value.category_id === "string" || value.category_id === null) &&
		(value.category === null || isAdminJSONRecipeCategoryResponse(value.category)) &&
		typeof value.created_at === "string" &&
		typeof value.updated_at === "string"
	);
}

function isAdminJSONRecipeListResponse(value: unknown): value is AdminJSONRecipeListResponse {
	return (
		isRecord(value) &&
		Array.isArray(value.recipes) &&
		value.recipes.every(isAdminJSONRecipeResponse) &&
		(typeof value.next_cursor === "string" || value.next_cursor === null)
	);
}

function isAdminJSONRecipeCategoryListResponse(value: unknown): value is AdminJSONRecipeCategoryListResponse {
	return (
		isRecord(value) &&
		Array.isArray(value.categories) &&
		value.categories.every(isAdminJSONRecipeCategoryResponse)
	);
}
