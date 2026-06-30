import { apiBaseUrl, internalAPIHeaders } from "./internal-api";

type ServerFetch = typeof fetch;
type JsonObject = Record<string, unknown>;

export type CollectionResponse = {
	id: string;
	created_at: string;
	updated_at: string;
	user_id: string;
	name: string;
	schema_ids: string[];
	schema_count: number;
	document_count: number;
};

export type CollectionListResponse = {
	collections: CollectionResponse[];
	next_cursor: string | null;
};

export type CollectionInput = {
	name: string;
	schema_ids: string[];
};

export type CollectionOptions = {
	userId?: string;
	cursor?: string;
	size?: string | number;
	sort?: "asc" | "desc" | string;
};

export class CollectionApiError extends Error {
	status: number;

	constructor(status: number, message: string) {
		super(message);
		this.name = "CollectionApiError";
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
		throw new CollectionApiError(503, "Collection service unavailable");
	}

	return parseResponseJSON(text);
}

function isJsonObject(value: unknown): value is JsonObject {
	return typeof value === "object" && value !== null && !Array.isArray(value);
}

function isStringArray(value: unknown): value is string[] {
	return Array.isArray(value) && value.every((item) => typeof item === "string");
}

function isFiniteNumber(value: unknown): value is number {
	return typeof value === "number" && Number.isFinite(value);
}

function isCollectionResponse(value: unknown): value is CollectionResponse {
	return (
		isJsonObject(value) &&
		typeof value.id === "string" &&
		typeof value.created_at === "string" &&
		typeof value.updated_at === "string" &&
		typeof value.user_id === "string" &&
		typeof value.name === "string" &&
		isStringArray(value.schema_ids) &&
		isFiniteNumber(value.schema_count) &&
		isFiniteNumber(value.document_count)
	);
}

function isCollectionListResponse(value: unknown): value is CollectionListResponse {
	return (
		isJsonObject(value) &&
		Array.isArray(value.collections) &&
		value.collections.every(isCollectionResponse) &&
		"next_cursor" in value &&
		(value.next_cursor === null ||
			(typeof value.next_cursor === "string" && value.next_cursor.trim() !== ""))
	);
}

function isDeleteCollectionResponse(value: unknown): value is { deleted_id: string } {
	return isJsonObject(value) && typeof value.deleted_id === "string";
}

export function isCollectionApiError(error: unknown): error is CollectionApiError {
	return error instanceof CollectionApiError;
}

function collectionListUrl(options: CollectionOptions = {}) {
	const url = new URL(`${apiBaseUrl()}/api/collections`);
	if (options.userId) url.searchParams.set("user_id", options.userId);
	if (options.size !== undefined) url.searchParams.set("size", String(options.size));
	if (options.cursor !== undefined) url.searchParams.set("cursor", options.cursor);
	if (options.sort !== undefined) url.searchParams.set("sort", options.sort);
	return url.toString();
}

function collectionDetailUrl(id: string, options: CollectionOptions = {}) {
	const url = new URL(`${apiBaseUrl()}/api/collection/${encodeURIComponent(id)}`);
	if (options.userId) url.searchParams.set("user_id", options.userId);
	return url.toString();
}

async function requestCollectionData(fetchFn: ServerFetch, url: string, init: RequestInit) {
	const headers = internalAPIHeaders(init.headers);
	if (!headers) {
		throw new CollectionApiError(500, "Collection service is not configured");
	}

	let response: Response;
	try {
		response = await fetchFn(url, { ...init, headers });
	} catch {
		throw new CollectionApiError(503, "Collection service unavailable");
	}

	const data = await readResponseJSON(response);

	if (!response.ok) {
		const message =
			data && typeof data === "object" && "error" in data && typeof data.error === "string"
				? data.error
				: "Collection request failed";
		throw new CollectionApiError(response.status, message);
	}

	return { data, status: response.status };
}

function jsonHeaders() {
	return new Headers({ "content-type": "application/json" });
}

export async function listCollectionsPage(
	fetchFn: ServerFetch,
	options: CollectionOptions = {}
) {
	const { data } = await requestCollectionData(fetchFn, collectionListUrl(options), {
		method: "GET"
	});

	if (!isCollectionListResponse(data)) {
		throw new CollectionApiError(502, "Invalid collection response");
	}

	return data;
}

export async function createCollection(
	fetchFn: ServerFetch,
	input: CollectionInput,
	options: CollectionOptions = {}
) {
	const body = options.userId ? { ...input, user_id: options.userId } : input;
	const { data } = await requestCollectionData(fetchFn, `${apiBaseUrl()}/api/collection`, {
		method: "POST",
		headers: jsonHeaders(),
		body: JSON.stringify(body)
	});

	if (!isCollectionResponse(data)) {
		throw new CollectionApiError(502, "Invalid collection response");
	}

	return data;
}

export async function updateCollection(
	fetchFn: ServerFetch,
	id: string,
	input: CollectionInput,
	options: CollectionOptions = {}
) {
	const { data } = await requestCollectionData(fetchFn, collectionDetailUrl(id, options), {
		method: "PUT",
		headers: jsonHeaders(),
		body: JSON.stringify(input)
	});

	if (!isCollectionResponse(data)) {
		throw new CollectionApiError(502, "Invalid collection response");
	}

	return data;
}

export async function deleteCollection(
	fetchFn: ServerFetch,
	id: string,
	options: CollectionOptions = {}
) {
	const { data, status } = await requestCollectionData(fetchFn, collectionDetailUrl(id, options), {
		method: "DELETE"
	});

	if (data === null && status === 204) {
		return { deleted_id: id };
	}

	if (!isDeleteCollectionResponse(data)) {
		throw new CollectionApiError(502, "Invalid collection delete response");
	}

	return data;
}
