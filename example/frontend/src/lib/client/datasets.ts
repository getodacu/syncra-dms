import { publicApiErrorMessage } from "$lib/client/api-errors";

type ClientFetch = typeof fetch;

export type JsonValue =
	| null
	| boolean
	| number
	| string
	| JsonValue[]
	| { [key: string]: JsonValue };

export type DatasetField = {
	path: string;
	key: string;
	label: string;
};

export type DatasetResponse = {
	id: string;
	created_at: string;
	updated_at: string;
	user_id: string;
	name: string;
	schema_id: string;
	schema_name: string;
	selected_fields: DatasetField[];
	field_count: number;
};

export type DatasetListResponse = {
	datasets: DatasetResponse[];
	next_cursor: string | null;
};

export type DatasetColumnResponse = {
	key: string;
	label: string;
	path: string;
};

export type DatasetRowResponse = {
	document_id: string;
	filename: string;
	created_at: string;
	values: Record<string, JsonValue>;
};

export type DatasetRowsResponse = {
	dataset: DatasetResponse;
	columns: DatasetColumnResponse[];
	rows: DatasetRowResponse[];
	next_cursor: string | null;
};

export type CreateDatasetInput = {
	name: string;
	schema_id: string;
	selected_fields: DatasetField[];
};

export type UpdateDatasetInput = {
	name: string;
	schema_id: string;
	selected_fields: DatasetField[];
};

export type DatasetListQuery = {
	cursor?: string | null;
	size?: number | string | null;
	sort?: "asc" | "desc" | string | null;
};

export type DatasetDateRangeQuery = {
	createdFrom?: string | null;
	createdTo?: string | null;
};

export type DatasetRowsQuery = DatasetListQuery & DatasetDateRangeQuery;

export type DatasetExportOptions = {
	format?: "csv" | "xlsx" | string | null;
	sort?: "asc" | "desc" | string | null;
} & DatasetDateRangeQuery;

export type DatasetExportResponse = {
	blob: Blob;
	headers: Headers;
	contentType: string | null;
	contentDisposition: string | null;
	filename: string | null;
};

const DEFAULT_DATASET_PAGE_SIZE = 100;

export class DatasetClientError extends Error {
	readonly status: number;

	constructor(status: number, message: string) {
		super(message);
		this.name = "DatasetClientError";
		this.status = status;
	}
}

export function isDatasetClientError(error: unknown): error is DatasetClientError {
	if (error instanceof DatasetClientError) return true;
	if (!isRecord(error)) return false;

	return error.name === "DatasetClientError" && typeof error.status === "number";
}

export function isDatasetNotFoundError(error: unknown) {
	return isDatasetClientError(error) && error.status === 404;
}

export async function fetchDatasets(
	fetchFn: ClientFetch,
	query: DatasetListQuery = {}
): Promise<DatasetListResponse> {
	const response = await fetchFn(datasetListPath(query), { method: "GET" });
	const json = await readResponseJSON(response);

	if (!response.ok) {
		throw datasetClientError(response, json, "Failed to load datasets");
	}
	if (!isDatasetListResponse(json)) {
		throw new Error("Invalid dataset list response");
	}

	return json;
}

export function createDataset(
	fetchFn: ClientFetch,
	input: CreateDatasetInput
): Promise<DatasetResponse> {
	return datasetRequest(
		fetchFn,
		"/api/datasets",
		{
			method: "POST",
			headers: jsonHeaders(),
			body: JSON.stringify(input),
		},
		"Failed to save dataset"
	);
}

export function getDataset(fetchFn: ClientFetch, id: string): Promise<DatasetResponse> {
	return datasetRequest(
		fetchFn,
		`/api/datasets/${encodeURIComponent(id)}`,
		{ method: "GET" },
		"Failed to load dataset"
	);
}

export function updateDataset(
	fetchFn: ClientFetch,
	id: string,
	input: UpdateDatasetInput
): Promise<DatasetResponse> {
	return datasetRequest(
		fetchFn,
		`/api/datasets/${encodeURIComponent(id)}`,
		{
			method: "PUT",
			headers: jsonHeaders(),
			body: JSON.stringify(input),
		},
		"Failed to save dataset"
	);
}

export async function deleteDataset(
	fetchFn: ClientFetch,
	id: string
): Promise<{ deleted_id: string }> {
	const response = await fetchFn(`/api/datasets/${encodeURIComponent(id)}`, {
		method: "DELETE",
	});
	const json = await readResponseJSON(response);

	if (!response.ok) {
		throw datasetClientError(response, json, "Failed to delete dataset");
	}
	if (response.status === 204 && json === null) {
		return { deleted_id: id };
	}
	if (!isDeleteDatasetResponse(json)) {
		throw new Error("Invalid dataset delete response");
	}

	return json;
}

export async function fetchDatasetRows(
	fetchFn: ClientFetch,
	id: string,
	query: DatasetRowsQuery = {}
): Promise<DatasetRowsResponse> {
	const response = await fetchFn(datasetRowsPath(id, query), { method: "GET" });
	const json = await readResponseJSON(response);

	if (!response.ok) {
		throw datasetClientError(response, json, "Failed to load dataset rows");
	}
	if (!isDatasetRowsResponse(json)) {
		throw new Error("Invalid dataset rows response");
	}

	return json;
}

export async function exportDataset(
	fetchFn: ClientFetch,
	id: string,
	options: DatasetExportOptions = {}
): Promise<DatasetExportResponse> {
	const response = await fetchFn(datasetExportPath(id, options), { method: "GET" });

	if (!response.ok) {
		const json = await readResponseJSON(response);
		throw datasetClientError(response, json, "Failed to export dataset");
	}

	const contentType = response.headers.get("content-type");
	const contentDisposition = response.headers.get("content-disposition");

	return {
		blob: await response.blob(),
		headers: new Headers(response.headers),
		contentType,
		contentDisposition,
		filename: parseContentDispositionFilename(contentDisposition),
	};
}

async function datasetRequest(
	fetchFn: ClientFetch,
	url: string,
	init: RequestInit,
	fallbackMessage: string
): Promise<DatasetResponse> {
	const response = await fetchFn(url, init);
	const json = await readResponseJSON(response);

	if (!response.ok) {
		throw datasetClientError(response, json, fallbackMessage);
	}
	if (!isDatasetResponse(json)) {
		throw new Error("Invalid dataset response");
	}

	return json;
}

function datasetListPath(query: DatasetListQuery = {}) {
	const params = new URLSearchParams();
	appendQueryParam(params, "cursor", query.cursor);
	appendQueryParam(params, "size", normalizedQueryValue(query.size) ?? DEFAULT_DATASET_PAGE_SIZE);
	appendQueryParam(params, "sort", query.sort);

	return `/api/datasets?${params.toString()}`;
}

function datasetRowsPath(id: string, query: DatasetRowsQuery = {}) {
	const params = new URLSearchParams();
	appendQueryParam(params, "created_from", query.createdFrom);
	appendQueryParam(params, "created_to", query.createdTo);
	appendQueryParam(params, "cursor", query.cursor);
	appendQueryParam(params, "size", normalizedQueryValue(query.size) ?? DEFAULT_DATASET_PAGE_SIZE);
	appendQueryParam(params, "sort", query.sort);

	return `/api/datasets/${encodeURIComponent(id)}/rows?${params.toString()}`;
}

function datasetExportPath(id: string, options: DatasetExportOptions = {}) {
	const params = new URLSearchParams();
	appendQueryParam(params, "format", options.format);
	appendQueryParam(params, "created_from", options.createdFrom);
	appendQueryParam(params, "created_to", options.createdTo);
	appendQueryParam(params, "sort", options.sort);

	const queryString = params.toString();
	return `/api/datasets/${encodeURIComponent(id)}/export${queryString ? `?${queryString}` : ""}`;
}

function datasetClientError(response: Response, json: unknown, fallbackMessage: string) {
	return new DatasetClientError(
		response.status,
		publicApiErrorMessage(response.status, json, fallbackMessage)
	);
}

function appendQueryParam(
	params: URLSearchParams,
	name: string,
	value: number | string | null | undefined
) {
	const normalizedValue = normalizedQueryValue(value);
	if (normalizedValue === undefined) return;

	params.append(name, normalizedValue);
}

function normalizedQueryValue(value: number | string | null | undefined) {
	if (value === null || value === undefined) return undefined;

	const stringValue = String(value).trim();
	return stringValue ? stringValue : undefined;
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

export function parseContentDispositionFilename(value: string | null) {
	if (!value) return null;

	const encodedMatch = value.match(/(?:^|;)\s*filename\*=([^']*)'[^']*'([^;]+)/i);
	if (encodedMatch) {
		const encodedFilename = encodedMatch[2].trim().replace(/^"|"$/g, "");
		try {
			return decodeURIComponent(encodedFilename);
		} catch {
			return encodedFilename;
		}
	}

	const quotedMatch = value.match(/(?:^|;)\s*filename="((?:\\.|[^"])*)"/i);
	if (quotedMatch) return quotedMatch[1].replace(/\\(.)/g, "$1");

	const plainMatch = value.match(/(?:^|;)\s*filename=([^;]+)/i);
	return plainMatch?.[1]?.trim() || null;
}

function isDatasetListResponse(value: unknown): value is DatasetListResponse {
	if (!isRecord(value)) return false;
	if (!Array.isArray(value.datasets)) return false;
	if (!isPageCursor(value.next_cursor)) return false;

	return value.datasets.every(isDatasetResponse);
}

function isDatasetResponse(value: unknown): value is DatasetResponse {
	if (!isRecord(value)) return false;

	return (
		typeof value.id === "string" &&
		typeof value.created_at === "string" &&
		typeof value.updated_at === "string" &&
		typeof value.user_id === "string" &&
		typeof value.name === "string" &&
		typeof value.schema_id === "string" &&
		typeof value.schema_name === "string" &&
		Array.isArray(value.selected_fields) &&
		value.selected_fields.every(isDatasetField) &&
		typeof value.field_count === "number" &&
		Number.isFinite(value.field_count)
	);
}

function isDatasetField(value: unknown): value is DatasetField {
	return (
		isRecord(value) &&
		typeof value.path === "string" &&
		typeof value.key === "string" &&
		typeof value.label === "string"
	);
}

function isDatasetColumnResponse(value: unknown): value is DatasetColumnResponse {
	return (
		isRecord(value) &&
		typeof value.key === "string" &&
		typeof value.label === "string" &&
		typeof value.path === "string"
	);
}

function isDatasetRowResponse(value: unknown): value is DatasetRowResponse {
	return (
		isRecord(value) &&
		typeof value.document_id === "string" &&
		typeof value.filename === "string" &&
		typeof value.created_at === "string" &&
		isJsonValueRecord(value.values)
	);
}

function isDatasetRowsResponse(value: unknown): value is DatasetRowsResponse {
	return (
		isRecord(value) &&
		isDatasetResponse(value.dataset) &&
		Array.isArray(value.columns) &&
		value.columns.every(isDatasetColumnResponse) &&
		Array.isArray(value.rows) &&
		value.rows.every(isDatasetRowResponse) &&
		isPageCursor(value.next_cursor)
	);
}

function isDeleteDatasetResponse(value: unknown): value is { deleted_id: string } {
	return isRecord(value) && typeof value.deleted_id === "string";
}

function isPageCursor(value: unknown): value is string | null {
	return value === null || (typeof value === "string" && value.trim() !== "");
}

function isJsonValueRecord(value: unknown): value is Record<string, JsonValue> {
	return isRecord(value) && Object.values(value).every(isJsonValue);
}

function isJsonValue(value: unknown): value is JsonValue {
	if (
		value === null ||
		typeof value === "boolean" ||
		typeof value === "string" ||
		(typeof value === "number" && Number.isFinite(value))
	) {
		return true;
	}

	if (Array.isArray(value)) return value.every(isJsonValue);
	if (isRecord(value)) return Object.values(value).every(isJsonValue);
	return false;
}

function isRecord(value: unknown): value is Record<string, unknown> {
	return typeof value === "object" && value !== null && !Array.isArray(value);
}
