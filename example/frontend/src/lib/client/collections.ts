import { publicApiErrorMessage } from "$lib/client/api-errors";

type ClientFetch = typeof fetch;

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

export type CollectionListQuery = {
	cursor?: string | null;
	size?: number | string;
	sort?: "asc" | "desc" | string;
};

export type CollectionInput = {
	name: string;
	schema_ids: string[];
};

export async function fetchCollections(
	fetchFn: ClientFetch,
	query: CollectionListQuery = {}
): Promise<CollectionListResponse> {
	const params = new URLSearchParams();
	appendQueryParam(params, "cursor", query.cursor);
	appendQueryParam(params, "size", query.size ?? 100);
	appendQueryParam(params, "sort", query.sort);

	const queryString = params.toString();
	const response = await fetchFn(`/api/collections${queryString ? `?${queryString}` : ""}`, {
		method: "GET",
	});
	const json = await readResponseJSON(response);

	if (!response.ok) {
		throw new Error(publicApiErrorMessage(response.status, json, "Failed to load collections"));
	}
	if (!isCollectionListResponse(json)) {
		throw new Error("Invalid collection list response");
	}

	return json;
}

function appendQueryParam(
	params: URLSearchParams,
	name: string,
	value: number | string | null | undefined
) {
	if (value === null || value === undefined) return;
	if (typeof value === "string" && !value.trim()) return;

	params.append(name, String(value));
}

export function createCollection(
	fetchFn: ClientFetch,
	input: CollectionInput
): Promise<CollectionResponse> {
	return collectionRequest(
		fetchFn,
		"/api/collections",
		{
			method: "POST",
			headers: jsonHeaders(),
			body: JSON.stringify(input),
		},
		"Failed to save collection"
	);
}

export function updateCollection(
	fetchFn: ClientFetch,
	id: string,
	input: CollectionInput
): Promise<CollectionResponse> {
	return collectionRequest(
		fetchFn,
		`/api/collections/${encodeURIComponent(id)}`,
		{
			method: "PUT",
			headers: jsonHeaders(),
			body: JSON.stringify(input),
		},
		"Failed to save collection"
	);
}

export async function deleteCollection(
	fetchFn: ClientFetch,
	id: string
): Promise<{ deleted_id: string }> {
	const response = await fetchFn(`/api/collections/${encodeURIComponent(id)}`, {
		method: "DELETE",
	});
	const json = await readResponseJSON(response);

	if (!response.ok) {
		throw new Error(publicApiErrorMessage(response.status, json, "Failed to delete collection"));
	}
	if (response.status === 204 && json === null) {
		return { deleted_id: id };
	}
	if (!isDeleteCollectionResponse(json)) {
		throw new Error("Invalid collection delete response");
	}

	return json;
}

async function collectionRequest(
	fetchFn: ClientFetch,
	url: string,
	init: RequestInit,
	fallbackMessage: string
): Promise<CollectionResponse> {
	const response = await fetchFn(url, init);
	const json = await readResponseJSON(response);

	if (!response.ok) {
		throw new Error(publicApiErrorMessage(response.status, json, fallbackMessage));
	}
	if (!isCollectionResponse(json)) {
		throw new Error("Invalid collection response");
	}

	return json;
}

function jsonHeaders() {
	return new Headers({ "content-type": "application/json" });
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

function isCollectionListResponse(value: unknown): value is CollectionListResponse {
	if (!isRecord(value)) return false;
	if (!Array.isArray(value.collections)) return false;
	if (!(typeof value.next_cursor === "string" || value.next_cursor === null)) return false;

	return value.collections.every(isCollectionResponse);
}

function isCollectionResponse(value: unknown): value is CollectionResponse {
	if (!isRecord(value)) return false;

	return (
		typeof value.id === "string" &&
		typeof value.created_at === "string" &&
		typeof value.updated_at === "string" &&
		typeof value.user_id === "string" &&
		typeof value.name === "string" &&
		Array.isArray(value.schema_ids) &&
		value.schema_ids.every((id) => typeof id === "string") &&
		typeof value.schema_count === "number" &&
		Number.isFinite(value.schema_count) &&
		typeof value.document_count === "number" &&
		Number.isFinite(value.document_count)
	);
}

function isDeleteCollectionResponse(value: unknown): value is { deleted_id: string } {
	return isRecord(value) && typeof value.deleted_id === "string";
}

function isRecord(value: unknown): value is Record<string, unknown> {
	return typeof value === "object" && value !== null && !Array.isArray(value);
}
