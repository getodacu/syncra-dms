import { apiBaseUrl, internalAPIHeaders } from "./internal-api";

type ServerFetch = typeof fetch;

export type JsonValue =
	| null
	| boolean
	| number
	| string
	| JsonValue[]
	| { [key: string]: JsonValue };

export type OCRDocumentResponse = {
	id: string;
	created_at: string;
	updated_at: string;
	user_id?: string;
	original_filename: string;
	mime_type: string;
	file_size: number;
	page_count: number;
	document_hash: string;
	schema_id?: string;
	has_inline_schema: boolean;
	markdown: string;
	annotation_json?: JsonValue;
	cached: boolean;
};

export type OCRDocumentCollectionResponse = {
	id: string;
	name: string;
};

export type OCRDocumentListItemResponse = {
	id: string;
	created_at: string;
	updated_at: string;
	user_id?: string;
	original_filename: string;
	mime_type: string;
	file_size: number;
	page_count: number;
	document_hash: string;
	schema_id?: string;
	has_inline_schema: boolean;
	collections: OCRDocumentCollectionResponse[];
};

export type OCRDocumentListResponse = {
	documents: OCRDocumentListItemResponse[];
	next_cursor: string | null;
};

export type UpdateOCRDocumentInput = {
	originalFilename: string;
};

export type DeleteOCRDocumentsResponse = {
	deleted_ids: string[];
	deleted_count: number;
};

export type DeleteOCRJobsResponse = {
	deleted_ids: string[];
	deleted_count: number;
};

export type MoveOCRDocumentsToCollectionsResponse = {
	moved_ids: string[];
	moved_count: number;
	collection_ids: string[];
};

export type ListOCRDocumentsOptions = {
	userId?: string;
	collectionId?: string;
	schemaId?: string;
	filename?: string;
	createdFrom?: string;
	createdTo?: string;
	cursor?: string;
	size?: string | number;
	sort?: "asc" | "desc" | string;
};

export type OCRJobStatus = "queued" | "processing" | "completed" | "failed" | string;

export type CreateOCRJobInput = {
	file: File;
	schemaId?: string;
	userId: string;
};

export type OCRJobResponse = {
	id: string;
	created_at: string;
	original_filename: string;
	mime_type: string;
	file_size: number;
	page_count: number;
	schema_id?: string;
	schema_name?: string;
	has_inline_schema: boolean;
	inline_schema?: JsonValue;
	document_id: string | null;
	status: OCRJobStatus;
	error_message?: string;
};

export type OCRJobListItemResponse = Omit<OCRJobResponse, "inline_schema">;

export type OCRJobListResponse = {
	jobs: OCRJobListItemResponse[];
	next_cursor: string | null;
};

export type ListOCRJobsOptions = {
	userId?: string;
	status?: string;
	cursor?: string;
	size?: string | number;
	sort?: "asc" | "desc" | string;
};

export class OCRApiError extends Error {
	status: number;

	constructor(status: number, message: string) {
		super(message);
		this.name = "OCRApiError";
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
		throw new OCRApiError(503, "OCR service unavailable");
	}

	return parseResponseJSON(text);
}

function isJsonValue(value: unknown): value is JsonValue {
	if (
		value === null ||
		typeof value === "boolean" ||
		typeof value === "number" ||
		typeof value === "string"
	) {
		return true;
	}

	if (Array.isArray(value)) return value.every(isJsonValue);

	if (typeof value === "object" && value !== null) {
		return Object.values(value).every(isJsonValue);
	}

	return false;
}

function isJsonObject(value: unknown): value is Record<string, unknown> {
	return typeof value === "object" && value !== null && !Array.isArray(value);
}

function isOCRDocumentResponse(value: unknown): value is OCRDocumentResponse {
	return (
		isJsonObject(value) &&
		typeof value.id === "string" &&
		typeof value.created_at === "string" &&
		typeof value.updated_at === "string" &&
		!("user_id" in value && typeof value.user_id !== "string") &&
		typeof value.original_filename === "string" &&
		typeof value.mime_type === "string" &&
		typeof value.file_size === "number" &&
		typeof value.page_count === "number" &&
		typeof value.document_hash === "string" &&
		!("schema_id" in value && typeof value.schema_id !== "string") &&
		typeof value.has_inline_schema === "boolean" &&
		typeof value.markdown === "string" &&
		!("annotation_json" in value && !isJsonValue(value.annotation_json)) &&
		typeof value.cached === "boolean"
	);
}

function isOCRDocumentListItemResponse(value: unknown): value is OCRDocumentListItemResponse {
	return (
		isJsonObject(value) &&
		typeof value.id === "string" &&
		typeof value.created_at === "string" &&
		typeof value.updated_at === "string" &&
		!("user_id" in value && typeof value.user_id !== "string") &&
		typeof value.original_filename === "string" &&
		typeof value.mime_type === "string" &&
		typeof value.file_size === "number" &&
		typeof value.page_count === "number" &&
		typeof value.document_hash === "string" &&
		!("schema_id" in value && typeof value.schema_id !== "string") &&
		typeof value.has_inline_schema === "boolean" &&
		Array.isArray(value.collections) &&
		value.collections.every(isOCRDocumentCollectionResponse)
	);
}

function isOCRDocumentCollectionResponse(value: unknown): value is OCRDocumentCollectionResponse {
	return isJsonObject(value) && typeof value.id === "string" && typeof value.name === "string";
}

function isOCRDocumentListResponse(value: unknown): value is OCRDocumentListResponse {
	return (
		isJsonObject(value) &&
		Array.isArray(value.documents) &&
		value.documents.every(isOCRDocumentListItemResponse) &&
		"next_cursor" in value &&
		(value.next_cursor === null || typeof value.next_cursor === "string")
	);
}

function isDeleteOCRDocumentsResponse(value: unknown): value is DeleteOCRDocumentsResponse {
	return (
		isJsonObject(value) &&
		Array.isArray(value.deleted_ids) &&
		value.deleted_ids.every((id) => typeof id === "string") &&
		typeof value.deleted_count === "number" &&
		Number.isFinite(value.deleted_count)
	);
}

function isMoveOCRDocumentsToCollectionsResponse(value: unknown): value is MoveOCRDocumentsToCollectionsResponse {
	return (
		isJsonObject(value) &&
		Array.isArray(value.moved_ids) &&
		value.moved_ids.every((id) => typeof id === "string") &&
		typeof value.moved_count === "number" &&
		Number.isFinite(value.moved_count) &&
		Array.isArray(value.collection_ids) &&
		value.collection_ids.every((id) => typeof id === "string")
	);
}

function isOCRJobResponse(value: unknown): value is OCRJobResponse {
	return (
		isJsonObject(value) &&
		typeof value.id === "string" &&
		typeof value.created_at === "string" &&
		typeof value.original_filename === "string" &&
		typeof value.mime_type === "string" &&
		typeof value.file_size === "number" &&
		typeof value.page_count === "number" &&
		!("schema_id" in value && typeof value.schema_id !== "string") &&
		!("schema_name" in value && typeof value.schema_name !== "string") &&
		typeof value.has_inline_schema === "boolean" &&
		!("inline_schema" in value && !isJsonValue(value.inline_schema)) &&
		"document_id" in value &&
		(value.document_id === null || typeof value.document_id === "string") &&
		typeof value.status === "string" &&
		!("error_message" in value && typeof value.error_message !== "string")
	);
}

function isOCRJobListItemResponse(value: unknown): value is OCRJobListItemResponse {
	return isOCRJobResponse(value) && !("inline_schema" in value);
}

function isOCRJobListResponse(value: unknown): value is OCRJobListResponse {
	return (
		isJsonObject(value) &&
		Array.isArray(value.jobs) &&
		value.jobs.every(isOCRJobListItemResponse) &&
		"next_cursor" in value &&
		(value.next_cursor === null || typeof value.next_cursor === "string")
	);
}

function isDeleteOCRJobsResponse(value: unknown): value is DeleteOCRJobsResponse {
	return isDeleteOCRDocumentsResponse(value);
}

export function isOCRApiError(error: unknown): error is OCRApiError {
	return error instanceof OCRApiError;
}

function internalHeaders(headers: HeadersInit = {}) {
	const output = internalAPIHeaders(headers);
	if (!output) {
		throw new OCRApiError(500, "OCR service is not configured");
	}
	return output;
}

function jsonHeaders() {
	return internalHeaders({ "content-type": "application/json" });
}

export async function getOCRDocument(
	fetchFn: ServerFetch,
	id: string,
	options: { userId?: string } = {}
) {
	const url = new URL(`${apiBaseUrl()}/api/ocr/document/${encodeURIComponent(id)}`);
	if (options.userId) url.searchParams.set("user_id", options.userId);

	const headers = internalHeaders();
	let response: Response;
	try {
		response = await fetchFn(url.toString(), {
			method: "GET",
			headers
		});
	} catch {
		throw new OCRApiError(503, "OCR service unavailable");
	}

	const data = await readResponseJSON(response);

	if (!response.ok) {
		const message =
			data && typeof data === "object" && "error" in data && typeof data.error === "string"
				? data.error
				: "OCR request failed";
		throw new OCRApiError(response.status, message);
	}

	if (!isOCRDocumentResponse(data)) {
		throw new OCRApiError(502, "Invalid OCR response");
	}

	return data;
}

export async function deleteOCRDocument(
	fetchFn: ServerFetch,
	id: string,
	options: { userId?: string } = {}
) {
	const url = new URL(`${apiBaseUrl()}/api/ocr/documents/${encodeURIComponent(id)}`);
	if (options.userId) url.searchParams.set("user_id", options.userId);

	const headers = internalHeaders();
	let response: Response;
	try {
		response = await fetchFn(url.toString(), {
			method: "DELETE",
			headers
		});
	} catch {
		throw new OCRApiError(503, "OCR service unavailable");
	}

	const data = await readResponseJSON(response);

	if (!response.ok) {
		const message =
			data && typeof data === "object" && "error" in data && typeof data.error === "string"
				? data.error
				: "OCR request failed";
		throw new OCRApiError(response.status, message);
	}

	if (response.status === 204) {
		return { deleted_ids: [id], deleted_count: 1 };
	}

	if (!isDeleteOCRDocumentsResponse(data)) {
		throw new OCRApiError(502, "Invalid OCR document delete response");
	}

	return data;
}

export async function updateOCRDocument(
	fetchFn: ServerFetch,
	id: string,
	input: UpdateOCRDocumentInput,
	options: { userId?: string } = {}
) {
	const url = new URL(`${apiBaseUrl()}/api/ocr/documents/${encodeURIComponent(id)}`);
	if (options.userId) url.searchParams.set("user_id", options.userId);

	const headers = jsonHeaders();
	let response: Response;
	try {
		response = await fetchFn(url.toString(), {
			method: "PATCH",
			headers,
			body: JSON.stringify({ original_filename: input.originalFilename })
		});
	} catch {
		throw new OCRApiError(503, "OCR service unavailable");
	}

	const data = await readResponseJSON(response);

	if (!response.ok) {
		const message =
			data && typeof data === "object" && "error" in data && typeof data.error === "string"
				? data.error
				: "OCR request failed";
		throw new OCRApiError(response.status, message);
	}

	if (!isOCRDocumentResponse(data)) {
		throw new OCRApiError(502, "Invalid OCR document update response");
	}

	return data;
}

function setOptionalSearchParam(url: URL, name: string, value: string | number | undefined) {
	if (value === undefined) return;
	const text = String(value).trim();
	if (text) url.searchParams.set(name, text);
}

export async function listOCRDocuments(
	fetchFn: ServerFetch,
	options: ListOCRDocumentsOptions = {}
) {
	const url = new URL(`${apiBaseUrl()}/api/ocr/documents`);
	setOptionalSearchParam(url, "user_id", options.userId);
	setOptionalSearchParam(url, "collection_id", options.collectionId);
	setOptionalSearchParam(url, "schema_id", options.schemaId);
	setOptionalSearchParam(url, "filename", options.filename);
	setOptionalSearchParam(url, "created_from", options.createdFrom);
	setOptionalSearchParam(url, "created_to", options.createdTo);
	setOptionalSearchParam(url, "cursor", options.cursor);
	setOptionalSearchParam(url, "size", options.size);
	setOptionalSearchParam(url, "sort", options.sort);

	const headers = internalHeaders();
	let response: Response;
	try {
		response = await fetchFn(url.toString(), {
			method: "GET",
			headers
		});
	} catch {
		throw new OCRApiError(503, "OCR service unavailable");
	}

	const data = await readResponseJSON(response);

	if (!response.ok) {
		const message =
			data && typeof data === "object" && "error" in data && typeof data.error === "string"
				? data.error
				: "OCR request failed";
		throw new OCRApiError(response.status, message);
	}

	if (!isOCRDocumentListResponse(data)) {
		throw new OCRApiError(502, "Invalid OCR document list response");
	}

	return data;
}

export async function deleteOCRDocuments(
	fetchFn: ServerFetch,
	ids: string[],
	options: { userId?: string } = {}
) {
	const url = new URL(`${apiBaseUrl()}/api/ocr/documents`);
	setOptionalSearchParam(url, "user_id", options.userId);

	const headers = jsonHeaders();
	let response: Response;
	try {
		response = await fetchFn(url.toString(), {
			method: "DELETE",
			headers,
			body: JSON.stringify({ ids })
		});
	} catch {
		throw new OCRApiError(503, "OCR service unavailable");
	}

	const data = await readResponseJSON(response);

	if (!response.ok) {
		const message =
			data && typeof data === "object" && "error" in data && typeof data.error === "string"
				? data.error
				: "OCR request failed";
		throw new OCRApiError(response.status, message);
	}

	if (!isDeleteOCRDocumentsResponse(data)) {
		throw new OCRApiError(502, "Invalid OCR document delete response");
	}

	return data;
}

export async function moveOCRDocumentsToCollections(
	fetchFn: ServerFetch,
	ids: string[],
	collectionIds: string[] = [],
	options: { userId?: string } = {}
) {
	const url = new URL(`${apiBaseUrl()}/api/ocr/documents/collections`);
	setOptionalSearchParam(url, "user_id", options.userId);

	const headers = jsonHeaders();
	let response: Response;
	try {
		response = await fetchFn(url.toString(), {
			method: "PUT",
			headers,
			body: JSON.stringify({ ids, collection_ids: collectionIds })
		});
	} catch {
		throw new OCRApiError(503, "OCR service unavailable");
	}

	const data = await readResponseJSON(response);

	if (!response.ok) {
		const message =
			data && typeof data === "object" && "error" in data && typeof data.error === "string"
				? data.error
				: "OCR request failed";
		throw new OCRApiError(response.status, message);
	}

	if (!isMoveOCRDocumentsToCollectionsResponse(data)) {
		throw new OCRApiError(502, "Invalid OCR document move response");
	}

	return data;
}

export async function createOCRJob(fetchFn: ServerFetch, input: CreateOCRJobInput) {
	const formData = new FormData();
	formData.set("file", input.file);
	if (input.schemaId) formData.set("schema_id", input.schemaId);
	formData.set("user_id", input.userId);

	const headers = internalHeaders();
	let response: Response;
	try {
		response = await fetchFn(`${apiBaseUrl()}/api/ocr/jobs`, {
			method: "POST",
			headers,
			body: formData
		});
	} catch {
		throw new OCRApiError(503, "OCR service unavailable");
	}

	const data = await readResponseJSON(response);

	if (!response.ok) {
		const message =
			data && typeof data === "object" && "error" in data && typeof data.error === "string"
				? data.error
				: "OCR request failed";
		throw new OCRApiError(response.status, message);
	}

	if (!isOCRJobResponse(data)) {
		throw new OCRApiError(502, "Invalid OCR job response");
	}

	return data;
}

export async function listOCRJobs(
	fetchFn: ServerFetch,
	options: ListOCRJobsOptions = {}
) {
	const url = new URL(`${apiBaseUrl()}/api/ocr/jobs`);
	setOptionalSearchParam(url, "user_id", options.userId);
	setOptionalSearchParam(url, "status", options.status);
	setOptionalSearchParam(url, "cursor", options.cursor);
	setOptionalSearchParam(url, "size", options.size);
	setOptionalSearchParam(url, "sort", options.sort);

	const headers = internalHeaders();
	let response: Response;
	try {
		response = await fetchFn(url.toString(), {
			method: "GET",
			headers
		});
	} catch {
		throw new OCRApiError(503, "OCR service unavailable");
	}

	const data = await readResponseJSON(response);

	if (!response.ok) {
		const message =
			data && typeof data === "object" && "error" in data && typeof data.error === "string"
				? data.error
				: "OCR request failed";
		throw new OCRApiError(response.status, message);
	}

	if (!isOCRJobListResponse(data)) {
		throw new OCRApiError(502, "Invalid OCR job list response");
	}

	return data;
}

export async function getOCRJob(
	fetchFn: ServerFetch,
	id: string,
	options: { userId?: string } = {}
) {
	const url = new URL(`${apiBaseUrl()}/api/ocr/jobs/${encodeURIComponent(id)}`);
	if (options.userId) url.searchParams.set("user_id", options.userId);

	const headers = internalHeaders();
	let response: Response;
	try {
		response = await fetchFn(url.toString(), {
			method: "GET",
			headers
		});
	} catch {
		throw new OCRApiError(503, "OCR service unavailable");
	}

	const data = await readResponseJSON(response);

	if (!response.ok) {
		const message =
			data && typeof data === "object" && "error" in data && typeof data.error === "string"
				? data.error
				: "OCR request failed";
		throw new OCRApiError(response.status, message);
	}

	if (!isOCRJobResponse(data)) {
		throw new OCRApiError(502, "Invalid OCR job response");
	}

	return data;
}

export async function deleteOCRJobs(
	fetchFn: ServerFetch,
	ids: string[],
	options: { userId?: string } = {}
) {
	const url = new URL(`${apiBaseUrl()}/api/ocr/jobs`);
	setOptionalSearchParam(url, "user_id", options.userId);

	const headers = jsonHeaders();
	let response: Response;
	try {
		response = await fetchFn(url.toString(), {
			method: "DELETE",
			headers,
			body: JSON.stringify({ ids })
		});
	} catch {
		throw new OCRApiError(503, "OCR service unavailable");
	}

	const data = await readResponseJSON(response);

	if (!response.ok) {
		const message =
			data && typeof data === "object" && "error" in data && typeof data.error === "string"
				? data.error
				: "OCR request failed";
		throw new OCRApiError(response.status, message);
	}

	if (!isDeleteOCRJobsResponse(data)) {
		throw new OCRApiError(502, "Invalid OCR job delete response");
	}

	return data;
}
