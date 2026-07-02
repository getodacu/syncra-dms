import { publicApiErrorMessage } from '$lib/client/api-errors';
import type { DocumentFolderNode, RepositoryDocument } from './tree';

type ClientFetch = typeof fetch;
type JsonObject = Record<string, unknown>;

export const DOCUMENT_FOLDERS_QUERY_KEY = ['document-folders'] as const;
export const DOCUMENT_FOLDER_CONTENTS_QUERY_KEY = ['document-folder-contents'] as const;

export type DocumentFolderTreeResponse = {
	folders: DocumentFolderNode[];
};

export type DocumentFolderContentsResponse = {
	folder: DocumentFolderNode;
	folders: DocumentFolderNode[];
	documents: RepositoryDocument[];
};

export type DocumentFolderInput = {
	organizationUnitId: string;
	parentId?: string | null;
	name: string;
	description?: string | null;
};

export type UpdateDocumentFolderVariables = {
	id: string;
	input: DocumentFolderInput;
};

export type MoveDocumentFolderVariables = {
	id: string;
	parentId: string | null;
};

export type ArchiveDocumentFolderVariables = {
	id: string;
};

export type ArchiveDocumentFolderResponse = {
	ok: boolean;
};

export type UploadDocumentVariables = {
	folderId: string;
	file: File;
};

export type DocumentUpdateInput = {
	displayName: string;
};

export type UpdateDocumentVariables = {
	id: string;
	input: DocumentUpdateInput;
};

export type ArchiveDocumentVariables = {
	id: string;
};

export type ArchiveDocumentResponse = {
	ok: boolean;
};

export async function fetchDocumentFolderTree(
	fetchFn: ClientFetch,
	organizationUnitId: string
): Promise<DocumentFolderTreeResponse> {
	const params = new URLSearchParams({ organizationUnitId });
	return documentJSONRequest(
		fetchFn,
		`/api/document-folders/tree?${params.toString()}`,
		{ method: 'GET' },
		validateDocumentFolderTreeResponse,
		'Failed to load document folders'
	);
}

export async function createDocumentFolder(fetchFn: ClientFetch, input: DocumentFolderInput) {
	return documentJSONRequest(
		fetchFn,
		'/api/document-folders',
		{
			method: 'POST',
			body: documentFolderPayload(input)
		},
		validateDocumentFolderNode,
		'Failed to create document folder'
	);
}

export async function updateDocumentFolder(
	fetchFn: ClientFetch,
	{ id, input }: UpdateDocumentFolderVariables
) {
	return documentJSONRequest(
		fetchFn,
		`/api/document-folders/${pathId(id)}`,
		{
			method: 'PATCH',
			body: documentFolderPayload(input)
		},
		validateDocumentFolderNode,
		'Failed to update document folder'
	);
}

export async function moveDocumentFolder(
	fetchFn: ClientFetch,
	{ id, parentId }: MoveDocumentFolderVariables
) {
	return documentJSONRequest(
		fetchFn,
		`/api/document-folders/${pathId(id)}/parent`,
		{
			method: 'PATCH',
			body: { parentId: normalizeParentId(parentId) }
		},
		validateDocumentFolderNode,
		'Failed to move document folder'
	);
}

export async function archiveDocumentFolder(
	fetchFn: ClientFetch,
	{ id }: ArchiveDocumentFolderVariables
) {
	return documentJSONRequest(
		fetchFn,
		`/api/document-folders/${pathId(id)}/archive`,
		{
			method: 'POST',
			body: {}
		},
		validateArchiveDocumentFolderResponse,
		'Failed to archive document folder'
	);
}

export async function fetchDocumentFolderContents(
	fetchFn: ClientFetch,
	folderId: string
): Promise<DocumentFolderContentsResponse> {
	return documentJSONRequest(
		fetchFn,
		`/api/document-folders/${pathId(folderId)}/contents`,
		{ method: 'GET' },
		validateDocumentFolderContentsResponse,
		'Failed to load document folder contents'
	);
}

export async function uploadDocument(
	fetchFn: ClientFetch,
	{ folderId, file }: UploadDocumentVariables
) {
	const formData = new FormData();
	formData.set('folderId', folderId);
	formData.set('file', file);

	return documentFormRequest(
		fetchFn,
		'/api/documents/upload',
		formData,
		validateRepositoryDocument,
		'Failed to upload document'
	);
}

export async function updateDocument(
	fetchFn: ClientFetch,
	{ id, input }: UpdateDocumentVariables
) {
	return documentJSONRequest(
		fetchFn,
		`/api/documents/${pathId(id)}`,
		{
			method: 'PATCH',
			body: { displayName: input.displayName }
		},
		validateRepositoryDocument,
		'Failed to update document'
	);
}

export async function archiveDocument(fetchFn: ClientFetch, { id }: ArchiveDocumentVariables) {
	return documentJSONRequest(
		fetchFn,
		`/api/documents/${pathId(id)}/archive`,
		{
			method: 'POST',
			body: {}
		},
		validateArchiveDocumentResponse,
		'Failed to archive document'
	);
}

export function documentDownloadHref(id: string) {
	return `/api/documents/${pathId(id)}/download`;
}

async function documentJSONRequest<T>(
	fetchFn: ClientFetch,
	path: string,
	options: {
		method: 'GET' | 'PATCH' | 'POST';
		body?: unknown;
	},
	validate: (data: unknown) => T,
	fallback: string
) {
	const response = await fetchFn(path, {
		method: options.method,
		headers: options.body === undefined ? undefined : { 'content-type': 'application/json' },
		body: options.body === undefined ? undefined : JSON.stringify(options.body)
	});
	const json = await readResponseJSON(response);

	if (!response.ok) {
		throw new Error(publicApiErrorMessage(response.status, json, fallback));
	}

	return validate(json);
}

async function documentFormRequest<T>(
	fetchFn: ClientFetch,
	path: string,
	formData: FormData,
	validate: (data: unknown) => T,
	fallback: string
) {
	const response = await fetchFn(path, {
		method: 'POST',
		headers: undefined,
		body: formData
	});
	const json = await readResponseJSON(response);

	if (!response.ok) {
		throw new Error(publicApiErrorMessage(response.status, json, fallback));
	}

	return validate(json);
}

async function readResponseJSON(response: Response): Promise<unknown> {
	let text: string;
	try {
		text = await response.text();
	} catch {
		return null;
	}

	if (!text.trim()) return null;

	try {
		return JSON.parse(text) as unknown;
	} catch {
		return null;
	}
}

function documentFolderPayload(input: DocumentFolderInput) {
	return {
		organizationUnitId: input.organizationUnitId,
		parentId: normalizeParentId(input.parentId),
		name: input.name,
		description: input.description ?? null
	};
}

function normalizeParentId(parentId: string | null | undefined) {
	if (!parentId) return null;
	const trimmed = parentId.trim();
	return trimmed ? trimmed : null;
}

function validateDocumentFolderTreeResponse(data: unknown): DocumentFolderTreeResponse {
	if (!isJsonObject(data) || !Array.isArray(data.folders)) invalidDocumentResponse();
	return {
		folders: data.folders.map(validateDocumentFolderNode)
	};
}

function validateDocumentFolderContentsResponse(data: unknown): DocumentFolderContentsResponse {
	if (
		!isJsonObject(data) ||
		!('folder' in data) ||
		!Array.isArray(data.folders) ||
		!Array.isArray(data.documents)
	) {
		invalidDocumentResponse();
	}

	return {
		folder: validateDocumentFolderNode(data.folder),
		folders: data.folders.map(validateDocumentFolderNode),
		documents: data.documents.map(validateRepositoryDocument)
	};
}

function validateDocumentFolderNode(data: unknown): DocumentFolderNode {
	if (!isJsonObject(data)) invalidDocumentResponse();
	if (data.children !== undefined && !Array.isArray(data.children)) invalidDocumentResponse();

	const output: DocumentFolderNode = {
		id: requiredString(data.id),
		organizationUnitId: requiredString(data.organizationUnitId),
		name: requiredString(data.name),
		createdAt: requiredString(data.createdAt),
		updatedAt: requiredString(data.updatedAt),
		children: Array.isArray(data.children) ? data.children.map(validateDocumentFolderNode) : []
	};
	const parentId = optionalString(data.parentId);
	const description = optionalString(data.description);
	const deletedAt = optionalString(data.deletedAt);

	if (parentId !== undefined) output.parentId = parentId;
	if (description !== undefined) output.description = description;
	if (deletedAt !== undefined) output.deletedAt = deletedAt;

	return output;
}

function validateRepositoryDocument(data: unknown): RepositoryDocument {
	if (!isJsonObject(data)) invalidDocumentResponse();

	const output: RepositoryDocument = {
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
	const extension = optionalString(data.extension);
	const deletedAt = optionalString(data.deletedAt);

	if (extension !== undefined) output.extension = extension;
	if (deletedAt !== undefined) output.deletedAt = deletedAt;

	return output;
}

function validateArchiveDocumentFolderResponse(data: unknown): ArchiveDocumentFolderResponse {
	if (!isJsonObject(data) || typeof data.ok !== 'boolean') invalidDocumentResponse();
	return { ok: data.ok };
}

function validateArchiveDocumentResponse(data: unknown): ArchiveDocumentResponse {
	if (!isJsonObject(data) || typeof data.ok !== 'boolean') invalidDocumentResponse();
	return { ok: data.ok };
}

function requiredString(value: unknown) {
	if (typeof value !== 'string' || value.trim() === '') invalidDocumentResponse();
	return value;
}

function requiredSizeBytes(value: unknown) {
	if (typeof value !== 'number' || !Number.isSafeInteger(value) || value < 0) {
		invalidDocumentResponse();
	}
	return value;
}

function optionalString(value: unknown) {
	if (value === undefined || value === null) return value;
	if (typeof value === 'string') return value;
	invalidDocumentResponse();
}

function isJsonObject(value: unknown): value is JsonObject {
	return typeof value === 'object' && value !== null && !Array.isArray(value);
}

function pathId(id: string) {
	return encodeURIComponent(id);
}

function invalidDocumentResponse(): never {
	throw new Error('Invalid document response');
}
