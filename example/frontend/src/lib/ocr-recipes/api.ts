import { publicApiErrorMessage } from "$lib/client/api-errors";

export type JsonSchemaObject = Record<string, unknown>;

export type OCRRecipeResponse = {
	id: string;
	title: string;
	description: string;
	json: JsonSchemaObject;
	counter: number;
	category_id: string | null;
	category: OCRRecipeCategoryResponse | null;
	created_at: string;
	updated_at: string;
};

export type OCRRecipeCategoryTitle = {
	en: string;
	ro: string;
};

export type OCRRecipeCategoryResponse = {
	id: string;
	title: OCRRecipeCategoryTitle;
	created_at: string;
	updated_at: string;
};

export type OCRRecipeListResponse = {
	recipes: OCRRecipeResponse[];
	next_cursor: string | null;
};

export type OCRRecipeDeployResponse = {
	recipe: OCRRecipeResponse;
	schema: {
		id: string;
		created_at: string;
		updated_at: string;
		user_id?: string | null;
		name: string;
		description: string;
		strict: boolean;
		schema: JsonSchemaObject;
	};
};

export type OCRRecipeListQuery = {
	cursor?: string | null;
	size?: string | number;
	sort?: "asc" | "desc";
};

type ClientFetch = typeof fetch;

export function buildOCRRecipesPath(query: OCRRecipeListQuery = {}) {
	const params = new URLSearchParams();
	setNonEmptyParam(params, "cursor", query.cursor);
	setNonEmptyParam(params, "size", query.size);
	setNonEmptyParam(params, "sort", query.sort);

	const queryString = params.toString();
	return queryString ? `/api/json-recipes?${queryString}` : "/api/json-recipes";
}

export async function fetchOCRRecipes(
	fetchFn: ClientFetch,
	query: OCRRecipeListQuery = {}
): Promise<OCRRecipeListResponse> {
	const response = await fetchFn(buildOCRRecipesPath(query), { method: "GET" });
	const body = await readResponseJSON(response);

	if (!response.ok) {
		throw new Error(publicApiErrorMessage(response.status, body, "Failed to load OCR recipes"));
	}
	if (!isOCRRecipeListResponse(body)) {
		throw new Error("Invalid OCR recipe list response");
	}
	return body;
}

export async function deployOCRRecipe(
	fetchFn: ClientFetch,
	recipeId: string,
	userId: string
): Promise<OCRRecipeDeployResponse> {
	const response = await fetchFn(`/api/json-recipes/${encodeURIComponent(recipeId)}/deploy`, {
		method: "POST",
		headers: { "content-type": "application/json" },
		body: JSON.stringify({ user_id: userId })
	});
	const body = await readResponseJSON(response);

	if (!response.ok) {
		throw new Error(publicApiErrorMessage(response.status, body, "Failed to clone OCR recipe"));
	}
	if (!isOCRRecipeDeployResponse(body)) {
		throw new Error("Invalid OCR recipe deploy response");
	}
	return body;
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

function isOCRRecipeResponse(value: unknown): value is OCRRecipeResponse {
	return (
		isRecord(value) &&
		typeof value.id === "string" &&
		typeof value.title === "string" &&
		typeof value.description === "string" &&
		isRecord(value.json) &&
		typeof value.counter === "number" &&
		Number.isFinite(value.counter) &&
		(typeof value.category_id === "string" || value.category_id === null) &&
		(value.category === null || isOCRRecipeCategoryResponse(value.category)) &&
		typeof value.created_at === "string" &&
		typeof value.updated_at === "string"
	);
}

function isOCRRecipeCategoryTitle(value: unknown): value is OCRRecipeCategoryTitle {
	return isRecord(value) && typeof value.en === "string" && typeof value.ro === "string";
}

function isOCRRecipeCategoryResponse(value: unknown): value is OCRRecipeCategoryResponse {
	return (
		isRecord(value) &&
		typeof value.id === "string" &&
		isOCRRecipeCategoryTitle(value.title) &&
		typeof value.created_at === "string" &&
		typeof value.updated_at === "string"
	);
}

function isOCRRecipeListResponse(value: unknown): value is OCRRecipeListResponse {
	return (
		isRecord(value) &&
		Array.isArray(value.recipes) &&
		value.recipes.every(isOCRRecipeResponse) &&
		(typeof value.next_cursor === "string" || value.next_cursor === null)
	);
}

function isSchemaResponse(value: unknown): value is OCRRecipeDeployResponse["schema"] {
	return (
		isRecord(value) &&
		typeof value.id === "string" &&
		typeof value.created_at === "string" &&
		typeof value.updated_at === "string" &&
		(value.user_id === undefined || typeof value.user_id === "string" || value.user_id === null) &&
		typeof value.name === "string" &&
		typeof value.description === "string" &&
		typeof value.strict === "boolean" &&
		isRecord(value.schema)
	);
}

function isOCRRecipeDeployResponse(value: unknown): value is OCRRecipeDeployResponse {
	return isRecord(value) && isOCRRecipeResponse(value.recipe) && isSchemaResponse(value.schema);
}
