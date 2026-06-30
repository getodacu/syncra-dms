import { apiBaseUrl, internalAPIHeaders } from "./internal-api";

type ServerFetch = typeof fetch;

export type JsonSchemaObject = Record<string, unknown>;

export type CreateSchemaInput = {
	name: string;
	description: string;
	strict: boolean;
	schema: JsonSchemaObject;
};

export type CreateSchemaOptions = {
	userId?: string;
};

export type ListSchemasOptions = {
	userId?: string;
};

export type ListSchemasPageOptions = ListSchemasOptions & {
	cursor?: string;
	size?: string | number;
	sort?: "asc" | "desc" | string;
};

export type ListJsonRecipesOptions = {
	cursor?: string | null;
	size?: string | number;
	sort?: "asc" | "desc" | string;
};

export type UpdateSchemaInput = CreateSchemaInput;

export type SchemaResponse = {
	id: string;
	created_at: string;
	updated_at: string;
	user_id?: string;
	name: string;
	description: string;
	schema: JsonSchemaObject;
	strict: boolean;
};

export type SchemaListResponse = {
	schemas: SchemaResponse[];
	next_cursor: string | null;
};

export type DeleteSchemasResponse = {
	deleted_ids: string[];
	deleted_count: number;
};

export type JsonRecipeResponse = {
	id: string;
	title: string;
	description: string;
	json: JsonSchemaObject;
	counter: number;
	category_id: string | null;
	category: JsonRecipeCategoryResponse | null;
	created_at: string;
	updated_at: string;
};

export type JsonRecipeCategoryTitle = {
	en: string;
	ro: string;
};

export type JsonRecipeCategoryResponse = {
	id: string;
	title: JsonRecipeCategoryTitle;
	created_at: string;
	updated_at: string;
};

export type JsonRecipeListResponse = {
	recipes: JsonRecipeResponse[];
	next_cursor: string | null;
};

export type JsonRecipeDeployResponse = {
	recipe: JsonRecipeResponse;
	schema: SchemaResponse;
};

export type DeployJsonRecipeOptions = {
	userId: string;
};

export class SchemaApiError extends Error {
	status: number;

	constructor(status: number, message: string) {
		super(message);
		this.name = "SchemaApiError";
		this.status = status;
	}
}

function parseResponseJSON(text: string) {
	if (!text) return null;
	try {
		return JSON.parse(text) as unknown;
	} catch {
		return undefined;
	}
}

async function readResponseJSON(response: Response) {
	let text: string;
	try {
		text = await response.text();
	} catch {
		throw new SchemaApiError(503, "Schema service unavailable");
	}

	return parseResponseJSON(text);
}

function isJsonObject(value: unknown): value is JsonSchemaObject {
	return typeof value === "object" && value !== null && !Array.isArray(value);
}

function isSchemaResponse(value: unknown): value is SchemaResponse {
	return (
		isJsonObject(value) &&
		typeof value.id === "string" &&
		typeof value.created_at === "string" &&
		typeof value.updated_at === "string" &&
		!("user_id" in value && typeof value.user_id !== "string") &&
		typeof value.name === "string" &&
		typeof value.description === "string" &&
		typeof value.strict === "boolean" &&
		isJsonObject(value.schema)
	);
}

function isJsonRecipeResponse(value: unknown): value is JsonRecipeResponse {
	return (
		isJsonObject(value) &&
		typeof value.id === "string" &&
		typeof value.title === "string" &&
		typeof value.description === "string" &&
		isJsonObject(value.json) &&
		typeof value.counter === "number" &&
		Number.isFinite(value.counter) &&
		(typeof value.category_id === "string" || value.category_id === null) &&
		(value.category === null || isJsonRecipeCategoryResponse(value.category)) &&
		typeof value.created_at === "string" &&
		typeof value.updated_at === "string"
	);
}

function isJsonRecipeCategoryTitle(value: unknown): value is JsonRecipeCategoryTitle {
	return isJsonObject(value) && typeof value.en === "string" && typeof value.ro === "string";
}

function isJsonRecipeCategoryResponse(value: unknown): value is JsonRecipeCategoryResponse {
	return (
		isJsonObject(value) &&
		typeof value.id === "string" &&
		isJsonRecipeCategoryTitle(value.title) &&
		typeof value.created_at === "string" &&
		typeof value.updated_at === "string"
	);
}

function isJsonRecipeDeployResponse(value: unknown): value is JsonRecipeDeployResponse {
	return isJsonObject(value) && isJsonRecipeResponse(value.recipe) && isSchemaResponse(value.schema);
}

function isJsonRecipeListResponse(value: unknown): value is JsonRecipeListResponse {
	return (
		isJsonObject(value) &&
		Array.isArray(value.recipes) &&
		value.recipes.every(isJsonRecipeResponse) &&
		"next_cursor" in value &&
		(value.next_cursor === null ||
			(typeof value.next_cursor === "string" && value.next_cursor.trim() !== ""))
	);
}

function isSchemaListResponse(value: unknown): value is SchemaListResponse {
	return (
		isJsonObject(value) &&
		Array.isArray(value.schemas) &&
		value.schemas.every(isSchemaResponse) &&
		"next_cursor" in value &&
		(value.next_cursor === null ||
			(typeof value.next_cursor === "string" && value.next_cursor.trim() !== ""))
	);
}

function isDeleteSchemasResponse(value: unknown): value is DeleteSchemasResponse {
	return (
		isJsonObject(value) &&
		Array.isArray(value.deleted_ids) &&
		value.deleted_ids.every((id) => typeof id === "string") &&
		typeof value.deleted_count === "number" &&
		Number.isFinite(value.deleted_count)
	);
}

export function isSchemaApiError(error: unknown): error is SchemaApiError {
	return error instanceof SchemaApiError;
}

function schemaListUrl(options: ListSchemasPageOptions = {}) {
	const url = new URL(`${apiBaseUrl()}/api/ocr/schemas`);
	if (options.userId) url.searchParams.set("user_id", options.userId);
	if (options.cursor !== undefined) url.searchParams.set("cursor", options.cursor);
	if (options.size !== undefined) url.searchParams.set("size", String(options.size));
	if (options.sort !== undefined) url.searchParams.set("sort", options.sort);
	return url.toString();
}

function schemaDetailUrl(id: string, options: CreateSchemaOptions = {}) {
	const url = new URL(`${apiBaseUrl()}/api/ocr/schemas/${encodeURIComponent(id)}`);
	if (options.userId) url.searchParams.set("user_id", options.userId);
	return url.toString();
}

function jsonRecipeListUrl(options: ListJsonRecipesOptions = {}) {
	const url = new URL(`${apiBaseUrl()}/api/json-recipes`);
	if (options.cursor) url.searchParams.set("cursor", options.cursor);
	if (options.size !== undefined) url.searchParams.set("size", String(options.size));
	if (options.sort !== undefined) url.searchParams.set("sort", options.sort);
	return url.toString();
}

async function requestSchemaData(fetchFn: ServerFetch, url: string, init: RequestInit) {
	const headers = internalAPIHeaders(init.headers);
	if (!headers) {
		throw new SchemaApiError(500, "Schema service is not configured");
	}

	let response: Response;
	try {
		response = await fetchFn(url, { ...init, headers });
	} catch {
		throw new SchemaApiError(503, "Schema service unavailable");
	}

	const data = await readResponseJSON(response);

	if (!response.ok) {
		const message =
			data && typeof data === "object" && "error" in data && typeof data.error === "string"
				? data.error
				: "Schema request failed";
		throw new SchemaApiError(response.status, message);
	}

	return data;
}

function jsonHeaders() {
	return new Headers({ "content-type": "application/json" });
}

export async function createSchema(
	fetchFn: ServerFetch,
	input: CreateSchemaInput,
	options: CreateSchemaOptions = {}
) {
	const body = options.userId ? { ...input, user_id: options.userId } : input;
	const data = await requestSchemaData(fetchFn, `${apiBaseUrl()}/api/ocr/schemas`, {
		method: "POST",
		headers: jsonHeaders(),
		body: JSON.stringify(body)
	});

	if (!isSchemaResponse(data)) {
		throw new SchemaApiError(502, "Invalid schema response");
	}

	return data;
}

export async function listSchemasPage(
	fetchFn: ServerFetch,
	options: ListSchemasPageOptions = {}
) {
	const data = await requestSchemaData(fetchFn, schemaListUrl(options), { method: "GET" });

	if (!isSchemaListResponse(data)) {
		throw new SchemaApiError(502, "Invalid schema response");
	}

	return data;
}

export async function listSchemas(fetchFn: ServerFetch, options: ListSchemasOptions = {}) {
	const schemas: SchemaResponse[] = [];
	const seenCursors = new Set<string>();
	let cursor: string | null = null;

	do {
		if (cursor !== null) {
			if (seenCursors.has(cursor)) {
				throw new SchemaApiError(502, "Invalid schema pagination response");
			}
			seenCursors.add(cursor);
		}
		const page = await listSchemasPage(fetchFn, {
			...options,
			...(cursor !== null ? { cursor } : {}),
			size: 100
		});
		schemas.push(...page.schemas);
		cursor = page.next_cursor;
	} while (cursor !== null);

	return schemas;
}

export async function getSchema(
	fetchFn: ServerFetch,
	id: string,
	options: CreateSchemaOptions = {}
) {
	const data = await requestSchemaData(fetchFn, schemaDetailUrl(id, options), { method: "GET" });

	if (!isSchemaResponse(data)) {
		throw new SchemaApiError(502, "Invalid schema response");
	}

	return data;
}

export async function updateSchema(
	fetchFn: ServerFetch,
	id: string,
	input: UpdateSchemaInput,
	options: CreateSchemaOptions = {}
) {
	const data = await requestSchemaData(fetchFn, schemaDetailUrl(id, options), {
		method: "PUT",
		headers: jsonHeaders(),
		body: JSON.stringify(input)
	});

	if (!isSchemaResponse(data)) {
		throw new SchemaApiError(502, "Invalid schema response");
	}

	return data;
}

export async function deleteSchema(
	fetchFn: ServerFetch,
	id: string,
	options: CreateSchemaOptions = {}
) {
	const data = await requestSchemaData(fetchFn, schemaDetailUrl(id, options), {
		method: "DELETE"
	});

	if (!isDeleteSchemasResponse(data)) {
		throw new SchemaApiError(502, "Invalid schema response");
	}

	return data;
}

export async function deleteSchemas(
	fetchFn: ServerFetch,
	ids: string[],
	options: ListSchemasOptions = {}
) {
	const data = await requestSchemaData(fetchFn, schemaListUrl(options), {
		method: "DELETE",
		headers: jsonHeaders(),
		body: JSON.stringify({ ids })
	});

	if (!isDeleteSchemasResponse(data)) {
		throw new SchemaApiError(502, "Invalid schema response");
	}

	return data;
}

export async function listJsonRecipes(
	fetchFn: ServerFetch,
	options: ListJsonRecipesOptions = {}
) {
	const data = await requestSchemaData(fetchFn, jsonRecipeListUrl(options), { method: "GET" });

	if (!isJsonRecipeListResponse(data)) {
		throw new SchemaApiError(502, "Invalid JSON recipe response");
	}

	return data;
}

export async function deployJsonRecipe(
	fetchFn: ServerFetch,
	recipeId: string,
	options: DeployJsonRecipeOptions
) {
	const data = await requestSchemaData(
		fetchFn,
		`${apiBaseUrl()}/api/json-recipes/${encodeURIComponent(recipeId)}/deploy`,
		{
			method: "POST",
			headers: jsonHeaders(),
			body: JSON.stringify({ user_id: options.userId })
		}
	);

	if (!isJsonRecipeDeployResponse(data)) {
		throw new SchemaApiError(502, "Invalid JSON recipe deploy response");
	}

	return data;
}
