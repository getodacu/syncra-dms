import { apiBaseUrl, internalAPIHeaders } from "./internal-api";

type ServerFetch = typeof fetch;
type JsonObject = Record<string, unknown>;

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

export type DatasetOptions = {
	userId?: string;
	createdFrom?: string;
	createdTo?: string;
	cursor?: string;
	size?: string | number;
	sort?: "asc" | "desc" | string;
};

export type DatasetExportOptions = {
	userId?: string;
	format?: "csv" | "xlsx" | string;
	createdFrom?: string;
	createdTo?: string;
	sort?: "asc" | "desc" | string;
};

export type DatasetExportResponse = {
	body: ReadableStream<Uint8Array> | null;
	headers: Headers;
	status: number;
};

export class DatasetApiError extends Error {
	status: number;

	constructor(status: number, message: string) {
		super(message);
		this.name = "DatasetApiError";
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
		throw new DatasetApiError(503, "Dataset service unavailable");
	}

	return parseResponseJSON(text);
}

function isJsonObject(value: unknown): value is JsonObject {
	return typeof value === "object" && value !== null && !Array.isArray(value);
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
	if (isJsonObject(value)) return Object.values(value).every(isJsonValue);
	return false;
}

function isJsonValueRecord(value: unknown): value is Record<string, JsonValue> {
	return isJsonObject(value) && Object.values(value).every(isJsonValue);
}

function isFiniteNumber(value: unknown): value is number {
	return typeof value === "number" && Number.isFinite(value);
}

function isDatasetField(value: unknown): value is DatasetField {
	return (
		isJsonObject(value) &&
		typeof value.path === "string" &&
		typeof value.key === "string" &&
		typeof value.label === "string"
	);
}

function isDatasetResponse(value: unknown): value is DatasetResponse {
	return (
		isJsonObject(value) &&
		typeof value.id === "string" &&
		typeof value.created_at === "string" &&
		typeof value.updated_at === "string" &&
		typeof value.user_id === "string" &&
		typeof value.name === "string" &&
		typeof value.schema_id === "string" &&
		typeof value.schema_name === "string" &&
		Array.isArray(value.selected_fields) &&
		value.selected_fields.every(isDatasetField) &&
		isFiniteNumber(value.field_count)
	);
}

function isPageCursor(value: unknown): value is string | null {
	return value === null || (typeof value === "string" && value.trim() !== "");
}

function isDatasetListResponse(value: unknown): value is DatasetListResponse {
	return (
		isJsonObject(value) &&
		Array.isArray(value.datasets) &&
		value.datasets.every(isDatasetResponse) &&
		"next_cursor" in value &&
		isPageCursor(value.next_cursor)
	);
}

function isDatasetColumnResponse(value: unknown): value is DatasetColumnResponse {
	return (
		isJsonObject(value) &&
		typeof value.key === "string" &&
		typeof value.label === "string" &&
		typeof value.path === "string"
	);
}

function isDatasetRowResponse(value: unknown): value is DatasetRowResponse {
	return (
		isJsonObject(value) &&
		typeof value.document_id === "string" &&
		typeof value.filename === "string" &&
		typeof value.created_at === "string" &&
		isJsonValueRecord(value.values)
	);
}

function isDatasetRowsResponse(value: unknown): value is DatasetRowsResponse {
	return (
		isJsonObject(value) &&
		isDatasetResponse(value.dataset) &&
		Array.isArray(value.columns) &&
		value.columns.every(isDatasetColumnResponse) &&
		Array.isArray(value.rows) &&
		value.rows.every(isDatasetRowResponse) &&
		"next_cursor" in value &&
		isPageCursor(value.next_cursor)
	);
}

function isDeleteDatasetResponse(value: unknown): value is { deleted_id: string } {
	return isJsonObject(value) && typeof value.deleted_id === "string";
}

export function isDatasetApiError(error: unknown): error is DatasetApiError {
	return error instanceof DatasetApiError;
}

function appendUserId(url: URL, options: { userId?: string }) {
	if (options.userId) url.searchParams.set("user_id", options.userId);
}

function appendPagination(url: URL, options: DatasetOptions) {
	if (options.cursor !== undefined) url.searchParams.set("cursor", options.cursor);
	if (options.size !== undefined) url.searchParams.set("size", String(options.size));
	if (options.sort !== undefined) url.searchParams.set("sort", options.sort);
}

function appendDateBounds(url: URL, options: { createdFrom?: string; createdTo?: string }) {
	if (options.createdFrom !== undefined) url.searchParams.set("created_from", options.createdFrom);
	if (options.createdTo !== undefined) url.searchParams.set("created_to", options.createdTo);
}

function datasetListUrl(options: DatasetOptions = {}) {
	const url = new URL(`${apiBaseUrl()}/api/datasets`);
	appendUserId(url, options);
	appendPagination(url, options);
	return url.toString();
}

function datasetDetailUrl(id: string, options: DatasetOptions = {}) {
	const url = new URL(`${apiBaseUrl()}/api/datasets/${encodeURIComponent(id)}`);
	appendUserId(url, options);
	return url.toString();
}

function datasetRowsUrl(id: string, options: DatasetOptions = {}) {
	const url = new URL(`${apiBaseUrl()}/api/datasets/${encodeURIComponent(id)}/rows`);
	appendUserId(url, options);
	appendDateBounds(url, options);
	appendPagination(url, options);
	return url.toString();
}

function datasetExportUrl(id: string, options: DatasetExportOptions = {}) {
	const url = new URL(`${apiBaseUrl()}/api/datasets/${encodeURIComponent(id)}/export`);
	appendUserId(url, options);
	if (options.format !== undefined) url.searchParams.set("format", options.format);
	appendDateBounds(url, options);
	if (options.sort !== undefined) url.searchParams.set("sort", options.sort);
	return url.toString();
}

async function requestDatasetData(fetchFn: ServerFetch, url: string, init: RequestInit) {
	const headers = internalAPIHeaders(init.headers);
	if (!headers) {
		throw new DatasetApiError(500, "Dataset service is not configured");
	}

	let response: Response;
	try {
		response = await fetchFn(url, { ...init, headers });
	} catch {
		throw new DatasetApiError(503, "Dataset service unavailable");
	}

	const data = await readResponseJSON(response);

	if (!response.ok) {
		const message =
			data && typeof data === "object" && "error" in data && typeof data.error === "string"
				? data.error
				: "Dataset request failed";
		throw new DatasetApiError(response.status, message);
	}

	return { data, status: response.status };
}

function jsonHeaders() {
	return new Headers({ "content-type": "application/json" });
}

export async function listDatasetsPage(fetchFn: ServerFetch, options: DatasetOptions = {}) {
	const { data } = await requestDatasetData(fetchFn, datasetListUrl(options), {
		method: "GET"
	});

	if (!isDatasetListResponse(data)) {
		throw new DatasetApiError(502, "Invalid dataset response");
	}

	return data;
}

export async function createDataset(
	fetchFn: ServerFetch,
	input: CreateDatasetInput,
	options: DatasetOptions = {}
) {
	const body = options.userId ? { ...input, user_id: options.userId } : input;
	const { data } = await requestDatasetData(fetchFn, `${apiBaseUrl()}/api/datasets`, {
		method: "POST",
		headers: jsonHeaders(),
		body: JSON.stringify(body)
	});

	if (!isDatasetResponse(data)) {
		throw new DatasetApiError(502, "Invalid dataset response");
	}

	return data;
}

export async function getDataset(
	fetchFn: ServerFetch,
	id: string,
	options: DatasetOptions = {}
) {
	const { data } = await requestDatasetData(fetchFn, datasetDetailUrl(id, options), {
		method: "GET"
	});

	if (!isDatasetResponse(data)) {
		throw new DatasetApiError(502, "Invalid dataset response");
	}

	return data;
}

export async function updateDataset(
	fetchFn: ServerFetch,
	id: string,
	input: UpdateDatasetInput,
	options: DatasetOptions = {}
) {
	const { data } = await requestDatasetData(fetchFn, datasetDetailUrl(id, options), {
		method: "PUT",
		headers: jsonHeaders(),
		body: JSON.stringify(input)
	});

	if (!isDatasetResponse(data)) {
		throw new DatasetApiError(502, "Invalid dataset response");
	}

	return data;
}

export async function deleteDataset(
	fetchFn: ServerFetch,
	id: string,
	options: DatasetOptions = {}
) {
	const { data, status } = await requestDatasetData(fetchFn, datasetDetailUrl(id, options), {
		method: "DELETE"
	});

	if (data === null && status === 204) {
		return { deleted_id: id };
	}

	if (!isDeleteDatasetResponse(data)) {
		throw new DatasetApiError(502, "Invalid dataset delete response");
	}

	return data;
}

export async function listDatasetRows(
	fetchFn: ServerFetch,
	id: string,
	options: DatasetOptions = {}
) {
	const { data } = await requestDatasetData(fetchFn, datasetRowsUrl(id, options), {
		method: "GET"
	});

	if (!isDatasetRowsResponse(data)) {
		throw new DatasetApiError(502, "Invalid dataset rows response");
	}

	return data;
}

export async function exportDataset(
	fetchFn: ServerFetch,
	id: string,
	options: DatasetExportOptions = {}
): Promise<DatasetExportResponse> {
	const requestHeaders = internalAPIHeaders();
	if (!requestHeaders) {
		throw new DatasetApiError(500, "Dataset service is not configured");
	}

	let response: Response;
	try {
		response = await fetchFn(datasetExportUrl(id, options), {
			method: "GET",
			headers: requestHeaders
		});
	} catch {
		throw new DatasetApiError(503, "Dataset service unavailable");
	}

	if (!response.ok) {
		const data = await readResponseJSON(response);
		const message =
			data && typeof data === "object" && "error" in data && typeof data.error === "string"
				? data.error
				: "Dataset export failed";
		throw new DatasetApiError(response.status, message);
	}

	const headers = new Headers();
	for (const header of ["content-type", "content-disposition"]) {
		const value = response.headers.get(header);
		if (value) headers.set(header, value);
	}

	return {
		body: response.body,
		headers,
		status: response.status
	};
}
