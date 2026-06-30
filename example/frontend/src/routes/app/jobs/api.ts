import { publicApiErrorMessage } from "$lib/client/api-errors";

import { buildJobsQueryPath, type JobListQuery } from "./table-utils";

export type JsonValue =
	| null
	| string
	| number
	| boolean
	| JsonValue[]
	| { [key: string]: JsonValue };

export type OCRJobListItemResponse = {
	id: string;
	created_at: string;
	original_filename: string;
	mime_type: string;
	status: string;
	file_size: number;
	page_count: number;
	schema_id?: string;
	schema_name?: string;
	has_inline_schema: boolean;
	document_id: string | null;
	error_message?: string;
};

export type OCRJobResponse = OCRJobListItemResponse & {
	inline_schema?: JsonValue;
};

export type OCRJobListResponse = {
	jobs: OCRJobListItemResponse[];
	next_cursor: string | null;
};

export type DeleteOCRJobsResponse = {
	deleted_ids: string[];
	deleted_count: number;
};

type ClientFetch = typeof fetch;

export async function fetchOCRJobs(
	fetchFn: ClientFetch,
	query: JobListQuery
): Promise<OCRJobListResponse> {
	const response = await fetchFn(buildJobsQueryPath(query), { method: "GET" });
	const json = await readResponseJSON(response);

	if (!response.ok) {
		throw new Error(publicApiErrorMessage(response.status, json, "Failed to load jobs"));
	}
	if (!isOCRJobListResponse(json)) {
		throw new Error("Invalid OCR job list response");
	}

	return json;
}

export async function fetchOCRJob(fetchFn: ClientFetch, id: string): Promise<OCRJobResponse> {
	const response = await fetchFn(`/api/ocr/jobs/${encodeURIComponent(id)}`, { method: "GET" });
	const json = await readResponseJSON(response);

	if (!response.ok) {
		throw new Error(publicApiErrorMessage(response.status, json, "Failed to load job"));
	}
	if (!isOCRJobResponse(json)) {
		throw new Error("Invalid OCR job response");
	}

	return json;
}

export async function deleteOCRJobs(
	fetchFn: ClientFetch,
	ids: string[]
): Promise<DeleteOCRJobsResponse> {
	const response = await fetchFn("/api/ocr/jobs", {
		method: "DELETE",
		headers: { "content-type": "application/json" },
		body: JSON.stringify({ ids }),
	});
	const json = await readResponseJSON(response);

	if (!response.ok) {
		throw new Error(publicApiErrorMessage(response.status, json, "Failed to delete jobs"));
	}
	if (!isDeleteOCRJobsResponse(json)) {
		throw new Error("Invalid OCR job delete response");
	}

	return json;
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

function isOCRJobListResponse(value: unknown): value is OCRJobListResponse {
	if (!isRecord(value)) return false;
	if (!Array.isArray(value.jobs)) return false;
	if (!(typeof value.next_cursor === "string" || value.next_cursor === null)) return false;

	return value.jobs.every(isOCRJobListItemResponse);
}

function isOCRJobListItemResponse(value: unknown): value is OCRJobListItemResponse {
	if (!isRecord(value)) return false;

	return (
		typeof value.id === "string" &&
		typeof value.created_at === "string" &&
		typeof value.original_filename === "string" &&
		typeof value.mime_type === "string" &&
		typeof value.status === "string" &&
		typeof value.file_size === "number" &&
		typeof value.page_count === "number" &&
		(value.schema_id === undefined || typeof value.schema_id === "string") &&
		(value.schema_name === undefined || typeof value.schema_name === "string") &&
		typeof value.has_inline_schema === "boolean" &&
		"document_id" in value &&
		(value.document_id === null || typeof value.document_id === "string") &&
		(value.error_message === undefined || typeof value.error_message === "string")
	);
}

function isOCRJobResponse(value: unknown): value is OCRJobResponse {
	if (!isOCRJobListItemResponse(value) || !isRecord(value)) return false;
	const record = value as Record<string, unknown>;

	return (
		record.inline_schema === undefined || isJsonValue(record.inline_schema)
	);
}

function isDeleteOCRJobsResponse(value: unknown): value is DeleteOCRJobsResponse {
	if (!isRecord(value)) return false;
	if (!Array.isArray(value.deleted_ids)) return false;

	return (
		value.deleted_ids.every((id) => typeof id === "string") &&
		typeof value.deleted_count === "number" &&
		Number.isFinite(value.deleted_count)
	);
}

function isRecord(value: unknown): value is Record<string, unknown> {
	return typeof value === "object" && value !== null && !Array.isArray(value);
}

function isJsonValue(value: unknown): value is JsonValue {
	if (
		value === null ||
		typeof value === "string" ||
		typeof value === "boolean" ||
		(typeof value === "number" && Number.isFinite(value))
	) {
		return true;
	}

	if (Array.isArray(value)) return value.every(isJsonValue);

	if (isRecord(value)) {
		return Object.values(value).every(isJsonValue);
	}

	return false;
}
