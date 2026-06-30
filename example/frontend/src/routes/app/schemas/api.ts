import { publicApiErrorMessage } from "$lib/client/api-errors";

import { buildSchemasQueryPath, type SchemaListQuery } from "./table-utils";

export type JsonSchemaObject = Record<string, unknown>;

export type SchemaEditorSubmitInput = {
	name: string;
	description: string;
	strict: boolean;
	schema: JsonSchemaObject;
};

export type SchemaResponse = {
	id: string;
	created_at: string;
	updated_at: string;
	user_id?: string | null;
	name: string;
	description: string;
	strict: boolean;
	schema: JsonSchemaObject;
};

export type SchemaListItemResponse = SchemaResponse;

export type SchemaListResponse = {
	schemas: SchemaListItemResponse[];
	next_cursor: string | null;
};

export type DeleteSchemasResponse = {
	deleted_ids: string[];
	deleted_count: number;
};

type ClientFetch = typeof fetch;
const MAX_SCHEMA_NAME_LENGTH = 160;
export const SCHEMA_QUERY_RETRY_LIMIT = 3;

export class SchemaClientError extends Error {
	readonly status: number;

	constructor(status: number, message: string) {
		super(message);
		this.name = "SchemaClientError";
		this.status = status;
	}
}

export function isSchemaClientError(error: unknown): error is SchemaClientError {
	if (error instanceof SchemaClientError) return true;
	if (!isRecord(error)) return false;

	return error.name === "SchemaClientError" && typeof error.status === "number";
}

export function isSchemaNotFoundError(error: unknown) {
	return isSchemaClientError(error) && error.status === 404;
}

export function shouldRetrySchemaQuery(failureCount: number, error: unknown) {
	return !isSchemaNotFoundError(error) && failureCount < SCHEMA_QUERY_RETRY_LIMIT;
}

export function createSchema(
	fetchFn: ClientFetch,
	input: SchemaEditorSubmitInput
): Promise<SchemaResponse> {
	return schemaRequest(
		fetchFn,
		"/api/schemas",
		{
			method: "POST",
			headers: jsonHeaders(),
			body: JSON.stringify(input),
		},
		"Failed to save schema"
	);
}

export function cloneSchema(fetchFn: ClientFetch, source: SchemaResponse): Promise<SchemaResponse> {
	return createSchema(fetchFn, {
		name: cloneSchemaName(source.name),
		description: source.description,
		strict: source.strict,
		schema: source.schema,
	});
}

export async function fetchSchemas(
	fetchFn: ClientFetch,
	query: SchemaListQuery
): Promise<SchemaListResponse> {
	const response = await fetchFn(buildSchemasQueryPath(query), { method: "GET" });
	const json = await readResponseJSON(response);

	if (!response.ok) {
		throw schemaClientError(response, json, "Failed to load schemas");
	}
	if (!isSchemaListResponse(json)) {
		throw new Error("Invalid schema list response");
	}

	return json;
}

export function getSchema(fetchFn: ClientFetch, id: string): Promise<SchemaResponse> {
	return schemaRequest(
		fetchFn,
		`/api/schemas/${encodeURIComponent(id)}`,
		{ method: "GET" },
		"Failed to load schema"
	);
}

export function updateSchema(
	fetchFn: ClientFetch,
	id: string,
	input: SchemaEditorSubmitInput
): Promise<SchemaResponse> {
	return schemaRequest(
		fetchFn,
		`/api/schemas/${encodeURIComponent(id)}`,
		{
			method: "PUT",
			headers: jsonHeaders(),
			body: JSON.stringify(input),
		},
		"Failed to save schema"
	);
}

export async function deleteSchema(
	fetchFn: ClientFetch,
	id: string
): Promise<DeleteSchemasResponse> {
	const response = await fetchFn(`/api/schemas/${encodeURIComponent(id)}`, { method: "DELETE" });
	const json = await readResponseJSON(response);

	if (!response.ok) {
		throw schemaClientError(response, json, "Failed to delete schema");
	}
	if (!isDeleteSchemasResponse(json)) {
		throw new Error("Invalid schema delete response");
	}

	return json;
}

export async function deleteSchemas(
	fetchFn: ClientFetch,
	ids: string[]
): Promise<DeleteSchemasResponse> {
	const response = await fetchFn("/api/schemas", {
		method: "DELETE",
		headers: { "content-type": "application/json" },
		body: JSON.stringify({ ids }),
	});
	const json = await readResponseJSON(response);

	if (!response.ok) {
		throw schemaClientError(response, json, "Failed to delete schemas");
	}
	if (!isDeleteSchemasResponse(json)) {
		throw new Error("Invalid schema delete response");
	}

	return json;
}

async function schemaRequest(
	fetchFn: ClientFetch,
	url: string,
	init: RequestInit,
	fallbackMessage: string
): Promise<SchemaResponse> {
	const response = await fetchFn(url, init);
	const json = await readResponseJSON(response);

	if (!response.ok) {
		throw schemaClientError(response, json, fallbackMessage);
	}

	if (!isSchemaResponse(json)) {
		throw new Error("Invalid schema response");
	}

	return json;
}

function schemaClientError(response: Response, json: unknown, fallbackMessage: string) {
	return new SchemaClientError(
		response.status,
		publicApiErrorMessage(response.status, json, fallbackMessage)
	);
}

function jsonHeaders() {
	return new Headers({ "content-type": "application/json" });
}

function cloneSchemaName(sourceName: string) {
	const baseName = sourceName.trim() || "Untitled schema";
	return Array.from(`Clone of ${baseName}`).slice(0, MAX_SCHEMA_NAME_LENGTH).join("");
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

function isSchemaListResponse(value: unknown): value is SchemaListResponse {
	if (!isRecord(value)) return false;
	if (!Array.isArray(value.schemas)) return false;
	if (
		!(
			(typeof value.next_cursor === "string" && value.next_cursor.trim()) ||
			value.next_cursor === null
		)
	) {
		return false;
	}

	return value.schemas.every(isSchemaListItemResponse);
}

function isSchemaResponse(value: unknown): value is SchemaResponse {
	if (!isRecord(value)) return false;

	return (
		typeof value.id === "string" &&
		typeof value.created_at === "string" &&
		typeof value.updated_at === "string" &&
		(value.user_id === undefined ||
			typeof value.user_id === "string" ||
			value.user_id === null) &&
		typeof value.name === "string" &&
		typeof value.description === "string" &&
		typeof value.strict === "boolean" &&
		isJsonObject(value.schema)
	);
}

function isSchemaListItemResponse(value: unknown): value is SchemaListItemResponse {
	return isSchemaResponse(value);
}

function isDeleteSchemasResponse(value: unknown): value is DeleteSchemasResponse {
	if (!isRecord(value)) return false;
	if (!Array.isArray(value.deleted_ids)) return false;

	return (
		value.deleted_ids.every((id) => typeof id === "string") &&
		typeof value.deleted_count === "number" &&
		Number.isInteger(value.deleted_count) &&
		value.deleted_count >= 0 &&
		value.deleted_count === value.deleted_ids.length
	);
}

function isRecord(value: unknown): value is Record<string, unknown> {
	return typeof value === "object" && value !== null && !Array.isArray(value);
}

function isJsonObject(value: unknown): value is JsonSchemaObject {
	return isRecord(value);
}
