import { publicApiErrorMessage } from "$lib/client/api-errors";

import { buildDocumentsQueryPath, type DocumentListQuery } from "./table-utils";

export type JsonValue =
	| null
	| string
	| number
	| boolean
	| JsonValue[]
	| { [key: string]: JsonValue };

export type OCRDocumentBaseResponse = {
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
};

export type OCRDocumentListItemResponse = OCRDocumentBaseResponse & {
	collections: OCRDocumentCollectionResponse[];
};

export type OCRDocumentCollectionResponse = {
	id: string;
	name: string;
};

export type OCRDocumentListResponse = {
	documents: OCRDocumentListItemResponse[];
	next_cursor: string | null;
};

export type OCRDocumentPreview = OCRDocumentBaseResponse & {
	markdown: string;
	annotation_json?: JsonValue;
	cached: boolean;
};

export type DeleteOCRDocumentsResponse = {
	deleted_ids: string[];
	deleted_count: number;
};

export type MoveOCRDocumentsResponse = {
	moved_ids: string[];
	moved_count: number;
	collection_ids: string[];
};

export type UpdateOCRDocumentResponse = OCRDocumentBaseResponse;
export type DownloadFormat = "markdown" | "html" | "json";
export type DownloadOCRDocumentsResponse = {
	blob: Blob;
	filename: string;
};

type ClientFetch = typeof fetch;

export const OCR_DOCUMENT_QUERY_RETRY_LIMIT = 3;

export class OCRDocumentClientError extends Error {
	readonly status: number;

	constructor(status: number, message: string) {
		super(message);
		this.name = "OCRDocumentClientError";
		this.status = status;
	}
}

export function isOCRDocumentClientError(error: unknown): error is OCRDocumentClientError {
	if (error instanceof OCRDocumentClientError) return true;
	if (!isRecord(error)) return false;

	return error.name === "OCRDocumentClientError" && typeof error.status === "number";
}

export function isOCRDocumentNotFoundError(error: unknown) {
	return isOCRDocumentClientError(error) && error.status === 404;
}

export function shouldRetryOCRDocumentsQuery(
	failureCount: number,
	error: unknown,
	query: Pick<DocumentListQuery, "collectionId"> = {}
) {
	return (
		!(query.collectionId && isOCRDocumentNotFoundError(error)) &&
		failureCount < OCR_DOCUMENT_QUERY_RETRY_LIMIT
	);
}

export async function fetchOCRDocuments(
	fetchFn: ClientFetch,
	query: DocumentListQuery,
	fallbackMessage = "Failed to load documents"
): Promise<OCRDocumentListResponse> {
	const response = await fetchFn(buildDocumentsQueryPath(query), { method: "GET" });
	const json = await readResponseJSON(response);

	if (!response.ok) {
		throw documentClientError(response, json, fallbackMessage);
	}
	if (!isOCRDocumentListResponse(json)) {
		throw new Error("Invalid OCR document list response");
	}

	return json;
}

export async function fetchOCRDocumentPreview(
	fetchFn: ClientFetch,
	id: string,
	fallbackMessage = "Failed to load document"
): Promise<OCRDocumentPreview> {
	const response = await fetchFn(`/api/ocr/document/${encodeURIComponent(id)}`, { method: "GET" });
	const json = await readResponseJSON(response);

	if (!response.ok) {
		throw documentClientError(response, json, fallbackMessage);
	}
	if (!isOCRDocumentPreview(json)) {
		throw new Error("Invalid OCR document response");
	}

	return json;
}

export async function deleteOCRDocument(
	fetchFn: ClientFetch,
	id: string,
	fallbackMessage = "Failed to delete document"
): Promise<DeleteOCRDocumentsResponse> {
	const response = await fetchFn(`/api/ocr/documents/${encodeURIComponent(id)}`, {
		method: "DELETE",
	});
	const json = await readResponseJSON(response);

	if (!response.ok) {
		throw documentClientError(response, json, fallbackMessage);
	}
	if (!isDeleteOCRDocumentsResponse(json)) {
		throw new Error("Invalid OCR document delete response");
	}

	return json;
}

export async function updateOCRDocument(
	fetchFn: ClientFetch,
	id: string,
	originalFilename: string,
	fallbackMessage = "Failed to update document"
): Promise<UpdateOCRDocumentResponse> {
	const response = await fetchFn(`/api/ocr/documents/${encodeURIComponent(id)}`, {
		method: "PATCH",
		headers: { "content-type": "application/json" },
		body: JSON.stringify({ original_filename: originalFilename }),
	});
	const json = await readResponseJSON(response);

	if (!response.ok) {
		throw documentClientError(response, json, fallbackMessage);
	}
	if (!isOCRDocumentBaseResponse(json)) {
		throw new Error("Invalid OCR document update response");
	}

	return json;
}

export async function deleteOCRDocuments(
	fetchFn: ClientFetch,
	ids: string[],
	fallbackMessage = "Failed to delete documents"
): Promise<DeleteOCRDocumentsResponse> {
	const response = await fetchFn("/api/ocr/documents", {
		method: "DELETE",
		headers: { "content-type": "application/json" },
		body: JSON.stringify({ ids }),
	});
	const json = await readResponseJSON(response);

	if (!response.ok) {
		throw documentClientError(response, json, fallbackMessage);
	}
	if (!isDeleteOCRDocumentsResponse(json)) {
		throw new Error("Invalid OCR document delete response");
	}

	return json;
}

export async function moveOCRDocuments(
	fetchFn: ClientFetch,
	ids: string[],
	collectionIds: string[],
	fallbackMessage = "Failed to move documents"
): Promise<MoveOCRDocumentsResponse> {
	const response = await fetchFn("/api/ocr/documents/collections", {
		method: "PUT",
		headers: { "content-type": "application/json" },
		body: JSON.stringify({ ids, collection_ids: collectionIds }),
	});
	const json = await readResponseJSON(response);

	if (!response.ok) {
		throw documentClientError(response, json, fallbackMessage);
	}
	if (!isMoveOCRDocumentsResponse(json)) {
		throw new Error("Invalid OCR document move response");
	}

	return json;
}

export async function downloadOCRDocuments(
	fetchFn: ClientFetch,
	ids: string[],
	format: DownloadFormat,
	fallbackMessage = "Failed to download documents"
): Promise<DownloadOCRDocumentsResponse> {
	const response = await fetchFn("/api/ocr/documents/download", {
		method: "POST",
		headers: { "content-type": "application/json" },
		body: JSON.stringify({ ids, format }),
	});

	if (!response.ok) {
		const json = await readResponseJSON(response);
		throw documentClientError(response, json, fallbackMessage);
	}

	return {
		blob: await response.blob(),
		filename: parseContentDispositionFilename(response.headers.get("content-disposition")) ?? "download",
	};
}

function documentClientError(response: Response, json: unknown, fallbackMessage: string) {
	return new OCRDocumentClientError(
		response.status,
		publicApiErrorMessage(response.status, json, fallbackMessage)
	);
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

	const encodedMatch = value.match(/(?:^|;)\s*filename\*=UTF-8''([^;]+)/i);
	if (encodedMatch) {
		try {
			return decodeURIComponent(encodedMatch[1]);
		} catch {
			return encodedMatch[1];
		}
	}

	const quotedMatch = value.match(/(?:^|;)\s*filename="([^"]*)"/i);
	if (quotedMatch) return quotedMatch[1].replace(/\\"/g, '"');

	const plainMatch = value.match(/(?:^|;)\s*filename=([^;]+)/i);
	return plainMatch?.[1]?.trim() || null;
}

function isOCRDocumentListResponse(value: unknown): value is OCRDocumentListResponse {
	if (!isRecord(value)) return false;
	if (!Array.isArray(value.documents)) return false;
	if (!(typeof value.next_cursor === "string" || value.next_cursor === null)) return false;

	return value.documents.every(isOCRDocumentListItemResponse);
}

function isDeleteOCRDocumentsResponse(value: unknown): value is DeleteOCRDocumentsResponse {
	if (!isRecord(value)) return false;
	if (!Array.isArray(value.deleted_ids)) return false;

	return (
		value.deleted_ids.every((id) => typeof id === "string") &&
		typeof value.deleted_count === "number" &&
		Number.isFinite(value.deleted_count)
	);
}

function isMoveOCRDocumentsResponse(value: unknown): value is MoveOCRDocumentsResponse {
	if (!isRecord(value)) return false;
	if (!Array.isArray(value.moved_ids)) return false;
	if (!Array.isArray(value.collection_ids)) return false;

	return (
		value.moved_ids.every((id) => typeof id === "string") &&
		typeof value.moved_count === "number" &&
		Number.isFinite(value.moved_count) &&
		value.collection_ids.every((id) => typeof id === "string")
	);
}

function isOCRDocumentListItemResponse(value: unknown): value is OCRDocumentListItemResponse {
	if (!isOCRDocumentBaseResponse(value)) return false;
	const record = value as Record<string, unknown>;

	return (
		Array.isArray(record.collections) &&
		record.collections.every(isOCRDocumentCollectionResponse)
	);
}

function isOCRDocumentBaseResponse(value: unknown): value is OCRDocumentBaseResponse {
	if (!isRecord(value)) return false;

	return (
		typeof value.id === "string" &&
		typeof value.created_at === "string" &&
		typeof value.updated_at === "string" &&
		(value.user_id === undefined || typeof value.user_id === "string") &&
		typeof value.original_filename === "string" &&
		typeof value.mime_type === "string" &&
		typeof value.file_size === "number" &&
		typeof value.page_count === "number" &&
		typeof value.document_hash === "string" &&
		(value.schema_id === undefined || typeof value.schema_id === "string") &&
		typeof value.has_inline_schema === "boolean"
	);
}

function isOCRDocumentCollectionResponse(value: unknown): value is OCRDocumentCollectionResponse {
	return isRecord(value) && typeof value.id === "string" && typeof value.name === "string";
}

function isOCRDocumentPreview(value: unknown): value is OCRDocumentPreview {
	if (!isRecord(value)) return false;
	if (!isOCRDocumentBaseResponse(value)) return false;

	const record = value as Record<string, unknown>;

	return (
		typeof record.markdown === "string" &&
		(record.annotation_json === undefined || isJsonValue(record.annotation_json)) &&
		typeof record.cached === "boolean"
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
	if (Array.isArray(value)) {
		return value.every(isJsonValue);
	}
	if (isRecord(value)) {
		return Object.values(value).every(isJsonValue);
	}

	return false;
}
