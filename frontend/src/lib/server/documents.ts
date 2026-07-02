import { apiBaseUrl, internalAPIHeaders } from './internal-api';
import { publicErrorMessage, publicErrorStatus } from './public-errors';

type ServerFetch = typeof fetch;
type HTTPMethod = 'GET' | 'PATCH' | 'POST';

export type DocumentFolder = {
	id: string;
	parentId?: string | null;
	organizationUnitId: string;
	name: string;
	description?: string | null;
	deletedAt?: string | null;
	createdAt: string;
	updatedAt: string;
	children?: DocumentFolder[];
};

export type DocumentMetadata = {
	id: string;
	folderId: string;
	organizationUnitId: string;
	originalFileName: string;
	displayName: string;
	mimeType: string;
	extension?: string | null;
	sizeBytes: number;
	sha256Hash: string;
	deletedAt?: string | null;
	createdAt: string;
	updatedAt: string;
};

export type DocumentFolderTreeResponse = {
	folders: DocumentFolder[];
};

export type DocumentFolderContentsResponse = {
	folder: DocumentFolder;
	folders: DocumentFolder[];
	documents: DocumentMetadata[];
};

export type DocumentFolderInput = {
	organizationUnitId: string;
	parentId?: string | null;
	name: string;
	description?: string | null;
};

export type DocumentUpdateInput = {
	displayName: string;
};

export type OKResponse = {
	ok: boolean;
};

export class DocumentApiError extends Error {
	status: number;

	constructor(status: number, message: string) {
		super(message);
		this.name = 'DocumentApiError';
		this.status = status;
	}
}

export function isDocumentApiError(error: unknown): error is DocumentApiError {
	return error instanceof DocumentApiError;
}

export async function getDocumentFolderTree(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	organizationUnitId: string
) {
	const params = new URLSearchParams({ organizationUnitId });
	return documentJSONRequest<DocumentFolderTreeResponse>(
		fetchFn,
		`/api/document-folders/tree?${params.toString()}`,
		{
			cookieHeader
		},
		validateDocumentFolderTreeResponse
	);
}

export async function createDocumentFolder(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	input: DocumentFolderInput
) {
	return documentJSONRequest<DocumentFolder>(
		fetchFn,
		'/api/document-folders',
		{
			method: 'POST',
			cookieHeader,
			body: input
		},
		validateDocumentFolder
	);
}

export async function updateDocumentFolder(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	id: string,
	input: DocumentFolderInput
) {
	return documentJSONRequest<DocumentFolder>(
		fetchFn,
		`/api/document-folders/${pathID(id)}`,
		{
			method: 'PATCH',
			cookieHeader,
			body: input
		},
		validateDocumentFolder
	);
}

export async function moveDocumentFolder(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	id: string,
	parentId: string | null
) {
	return documentJSONRequest<DocumentFolder>(
		fetchFn,
		`/api/document-folders/${pathID(id)}/parent`,
		{
			method: 'PATCH',
			cookieHeader,
			body: { parentId }
		},
		validateDocumentFolder
	);
}

export async function archiveDocumentFolder(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	id: string
) {
	return documentJSONRequest<OKResponse>(
		fetchFn,
		`/api/document-folders/${pathID(id)}/archive`,
		{
			method: 'POST',
			cookieHeader,
			body: {}
		},
		validateOKResponse
	);
}

export async function getDocumentFolderContents(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	id: string
) {
	return documentJSONRequest<DocumentFolderContentsResponse>(
		fetchFn,
		`/api/document-folders/${pathID(id)}/contents`,
		{
			cookieHeader
		},
		validateDocumentFolderContentsResponse
	);
}

export async function uploadDocument(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	formData: FormData
) {
	return documentFormRequest<DocumentMetadata>(
		fetchFn,
		'/api/documents/upload',
		{
			cookieHeader,
			formData
		},
		validateDocumentMetadata
	);
}

export async function getDocument(fetchFn: ServerFetch, cookieHeader: string | null, id: string) {
	return documentJSONRequest<DocumentMetadata>(
		fetchFn,
		`/api/documents/${pathID(id)}`,
		{
			cookieHeader
		},
		validateDocumentMetadata
	);
}

export async function updateDocument(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	id: string,
	input: DocumentUpdateInput
) {
	return documentJSONRequest<DocumentMetadata>(
		fetchFn,
		`/api/documents/${pathID(id)}`,
		{
			method: 'PATCH',
			cookieHeader,
			body: input
		},
		validateDocumentMetadata
	);
}

export async function archiveDocument(fetchFn: ServerFetch, cookieHeader: string | null, id: string) {
	return documentJSONRequest<OKResponse>(
		fetchFn,
		`/api/documents/${pathID(id)}/archive`,
		{
			method: 'POST',
			cookieHeader,
			body: {}
		},
		validateOKResponse
	);
}

export async function downloadDocument(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	id: string
) {
	const headers = documentInternalHeaders();
	if (cookieHeader) headers.set('cookie', cookieHeader);

	let response: Response;
	try {
		response = await fetchFn(`${apiBaseUrl()}/api/documents/${pathID(id)}/download`, {
			method: 'GET',
			headers
		});
	} catch {
		throw new DocumentApiError(503, 'Document service unavailable');
	}

	if (!response.ok) {
		await throwDocumentResponseError(response);
	}
	return response;
}

async function documentJSONRequest<T>(
	fetchFn: ServerFetch,
	path: string,
	options: {
		method?: HTTPMethod;
		body?: unknown;
		cookieHeader?: string | null;
	} = {},
	validate: (data: unknown) => T
) {
	const headers = documentInternalHeaders();
	if (options.body !== undefined) headers.set('content-type', 'application/json');
	if (options.cookieHeader) headers.set('cookie', options.cookieHeader);

	let response: Response;
	try {
		response = await fetchFn(`${apiBaseUrl()}${path}`, {
			method: options.method ?? 'GET',
			headers,
			body: options.body === undefined ? undefined : JSON.stringify(options.body)
		});
	} catch {
		throw new DocumentApiError(503, 'Document service unavailable');
	}

	const text = await response.text();
	const data = parseResponseJSON(text);
	if (!response.ok) {
		throwDocumentError(response.status, data);
	}
	return validate(data);
}

async function documentFormRequest<T>(
	fetchFn: ServerFetch,
	path: string,
	options: {
		formData: FormData;
		cookieHeader?: string | null;
	},
	validate: (data: unknown) => T
) {
	const headers = documentInternalHeaders();
	if (options.cookieHeader) headers.set('cookie', options.cookieHeader);

	let response: Response;
	try {
		response = await fetchFn(`${apiBaseUrl()}${path}`, {
			method: 'POST',
			headers,
			body: options.formData
		});
	} catch {
		throw new DocumentApiError(503, 'Document service unavailable');
	}

	const text = await response.text();
	const data = parseResponseJSON(text);
	if (!response.ok) {
		throwDocumentError(response.status, data);
	}
	return validate(data);
}

function documentInternalHeaders() {
	const headers = internalAPIHeaders();
	if (!headers) throw new DocumentApiError(500, 'Document service is not configured');
	return headers;
}

async function throwDocumentResponseError(response: Response): Promise<never> {
	const text = await response.text();
	throwDocumentError(response.status, parseResponseJSON(text));
}

function throwDocumentError(status: number, data: unknown): never {
	const message =
		data && typeof data === 'object' && 'error' in data && typeof data.error === 'string'
			? data.error
			: 'Document request failed';
	throw new DocumentApiError(
		publicErrorStatus(status),
		publicErrorMessage(status, message, 'Document request failed')
	);
}

function parseResponseJSON(text: string) {
	if (!text) return null;
	try {
		return JSON.parse(text) as unknown;
	} catch {
		return undefined;
	}
}

function validateDocumentFolderTreeResponse(data: unknown): DocumentFolderTreeResponse {
	if (!isRecord(data) || !Array.isArray(data.folders)) invalidDocumentResponse();
	return {
		folders: data.folders.map(validateDocumentFolderTreeNode)
	};
}

function validateDocumentFolderTreeNode(data: unknown): DocumentFolder {
	const folder = validateDocumentFolder(data);
	if (!isRecord(data) || !Array.isArray(data.children)) invalidDocumentResponse();
	return {
		...folder,
		children: data.children.map(validateDocumentFolderTreeNode)
	};
}

function validateDocumentFolderContentsResponse(data: unknown): DocumentFolderContentsResponse {
	if (
		!isRecord(data) ||
		!Array.isArray(data.folders) ||
		!Array.isArray(data.documents) ||
		!('folder' in data)
	) {
		invalidDocumentResponse();
	}
	return {
		folder: validateDocumentFolder(data.folder),
		folders: data.folders.map(validateDocumentFolder),
		documents: data.documents.map(validateDocumentMetadata)
	};
}

function validateDocumentFolder(data: unknown): DocumentFolder {
	if (!isRecord(data)) invalidDocumentResponse();
	const output: DocumentFolder = {
		id: requiredString(data.id),
		organizationUnitId: requiredString(data.organizationUnitId),
		name: requiredString(data.name),
		createdAt: requiredString(data.createdAt),
		updatedAt: requiredString(data.updatedAt)
	};
	assignOptionalString(output, 'parentId', data.parentId);
	assignOptionalString(output, 'description', data.description);
	assignOptionalString(output, 'deletedAt', data.deletedAt);
	if ('children' in data) {
		if (!Array.isArray(data.children)) invalidDocumentResponse();
		output.children = data.children.map(validateDocumentFolder);
	}
	return output;
}

function validateDocumentMetadata(data: unknown): DocumentMetadata {
	if (!isRecord(data)) invalidDocumentResponse();
	const output: DocumentMetadata = {
		id: requiredString(data.id),
		folderId: requiredString(data.folderId),
		organizationUnitId: requiredString(data.organizationUnitId),
		originalFileName: requiredString(data.originalFileName),
		displayName: requiredString(data.displayName),
		mimeType: requiredString(data.mimeType),
		sizeBytes: requiredSizeBytes(data.sizeBytes),
		sha256Hash: requiredString(data.sha256Hash),
		createdAt: requiredString(data.createdAt),
		updatedAt: requiredString(data.updatedAt)
	};
	assignOptionalString(output, 'extension', data.extension);
	assignOptionalString(output, 'deletedAt', data.deletedAt);
	return output;
}

function validateOKResponse(data: unknown): OKResponse {
	if (!isRecord(data) || typeof data.ok !== 'boolean') invalidDocumentResponse();
	return { ok: data.ok };
}

function requiredString(value: unknown) {
	if (typeof value !== 'string' || value.trim() === '') invalidDocumentResponse();
	return value;
}

function requiredSizeBytes(value: unknown) {
	if (typeof value !== 'number' || !Number.isFinite(value) || value < 0) invalidDocumentResponse();
	return value;
}

function optionalString(value: unknown) {
	if (value === undefined || value === null) return value;
	if (typeof value === 'string') return value;
	invalidDocumentResponse();
}

function assignOptionalString<T extends Record<string, unknown>, K extends keyof T>(
	output: T,
	key: K,
	value: unknown
) {
	const parsed = optionalString(value);
	if (parsed !== undefined) output[key] = parsed as T[K];
}

function isRecord(value: unknown): value is Record<string, unknown> {
	return typeof value === 'object' && value !== null;
}

function pathID(id: string) {
	return encodeURIComponent(id);
}

function invalidDocumentResponse(): never {
	throw new DocumentApiError(502, 'Invalid Document response');
}
