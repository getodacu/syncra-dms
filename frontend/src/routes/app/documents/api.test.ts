import { readFileSync } from 'node:fs';
import { describe, expect, it, vi } from 'vitest';

import {
	DOCUMENT_FOLDER_CONTENTS_QUERY_KEY,
	DOCUMENT_FOLDERS_QUERY_KEY,
	archiveDocument,
	archiveDocumentFolder,
	createDocumentFolder,
	documentDownloadHref,
	fetchDocumentFolderContents,
	fetchDocumentFolderTree,
	moveDocumentFolder,
	updateDocument,
	updateDocumentFolder,
	uploadDocument
} from './api';
import type { DocumentFolderContentsResponse, DocumentFolderTreeResponse } from './api';
import type { DocumentFolderNode, RepositoryDocument } from './tree';

function jsonResponse(body: unknown, init?: ResponseInit) {
	return new Response(JSON.stringify(body), {
		headers: { 'content-type': 'application/json' },
		...init
	});
}

function folderNode(overrides: Partial<DocumentFolderNode> = {}): DocumentFolderNode {
	return {
		id: 'folder-id',
		parentId: null,
		organizationUnitId: 'unit-id',
		name: 'Invoices',
		description: null,
		deletedAt: null,
		createdAt: '2026-07-01T10:00:00Z',
		updatedAt: '2026-07-01T10:00:00Z',
		children: [],
		...overrides
	};
}

function documentMetadata(overrides: Partial<RepositoryDocument> = {}): RepositoryDocument {
	return {
		id: 'document-id',
		folderId: 'folder-id',
		organizationUnitId: 'unit-id',
		originalFileName: 'invoice.pdf',
		displayName: 'invoice.pdf',
		mimeType: 'application/pdf',
		extension: '.pdf',
		sizeBytes: 7,
		sha256Hash: 'hash',
		deletedAt: null,
		createdAt: '2026-07-01T10:00:00Z',
		updatedAt: '2026-07-01T10:00:00Z',
		...overrides
	};
}

function treeResponse(): DocumentFolderTreeResponse {
	return { folders: [folderNode()] };
}

function contentsResponse(): DocumentFolderContentsResponse {
	return {
		folder: folderNode(),
		folders: [folderNode({ id: 'child-folder', parentId: 'folder-id', name: 'Receipts' })],
		documents: [documentMetadata()]
	};
}

describe('document repository browser API', () => {
	it('exports document query keys', () => {
		expect(DOCUMENT_FOLDERS_QUERY_KEY).toEqual(['document-folders']);
		expect(DOCUMENT_FOLDER_CONTENTS_QUERY_KEY).toEqual(['document-folder-contents']);
	});

	it('fetches folder tree through the Svelte API wrapper with an encoded query id', async () => {
		const fetchMock = vi.fn().mockResolvedValue(jsonResponse(treeResponse(), { status: 200 }));

		const result = await fetchDocumentFolderTree(fetchMock, 'unit/id');

		expect(fetchMock).toHaveBeenCalledWith(
			'/api/document-folders/tree?organizationUnitId=unit%2Fid',
			{
				method: 'GET',
				headers: undefined,
				body: undefined
			}
		);
		expect(result.folders[0].name).toBe('Invoices');
	});

	it('posts normalized create folder payloads', async () => {
		const fetchMock = vi.fn().mockResolvedValue(jsonResponse(folderNode(), { status: 201 }));

		await createDocumentFolder(fetchMock, {
			organizationUnitId: 'unit-id',
			parentId: '   ',
			name: 'Invoices',
			description: undefined
		});

		expect(fetchMock).toHaveBeenCalledWith('/api/document-folders', {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({
				organizationUnitId: 'unit-id',
				parentId: null,
				name: 'Invoices',
				description: null
			})
		});
	});

	it('patches update and move folder requests with encoded ids', async () => {
		const fetchMock = vi
			.fn()
			.mockResolvedValueOnce(jsonResponse(folderNode(), { status: 200 }))
			.mockResolvedValueOnce(jsonResponse(folderNode({ parentId: null }), { status: 200 }));

		await updateDocumentFolder(fetchMock, {
			id: 'folder/id',
			input: {
				organizationUnitId: 'unit-id',
				parentId: 'parent/id',
				name: 'Invoices Updated',
				description: ''
			}
		});
		await moveDocumentFolder(fetchMock, { id: 'folder/id', parentId: '   ' });

		expect(fetchMock).toHaveBeenNthCalledWith(1, '/api/document-folders/folder%2Fid', {
			method: 'PATCH',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({
				organizationUnitId: 'unit-id',
				parentId: 'parent/id',
				name: 'Invoices Updated',
				description: ''
			})
		});
		expect(fetchMock).toHaveBeenNthCalledWith(2, '/api/document-folders/folder%2Fid/parent', {
			method: 'PATCH',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({ parentId: null })
		});
	});

	it('posts archive folder requests and fetches folder contents with encoded ids', async () => {
		const fetchMock = vi
			.fn()
			.mockResolvedValueOnce(jsonResponse({ ok: true }, { status: 200 }))
			.mockResolvedValueOnce(jsonResponse(contentsResponse(), { status: 200 }));

		const archiveResult = await archiveDocumentFolder(fetchMock, { id: 'folder/id' });
		const contents = await fetchDocumentFolderContents(fetchMock, 'folder/id');

		expect(fetchMock).toHaveBeenNthCalledWith(1, '/api/document-folders/folder%2Fid/archive', {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({})
		});
		expect(fetchMock).toHaveBeenNthCalledWith(2, '/api/document-folders/folder%2Fid/contents', {
			method: 'GET',
			headers: undefined,
			body: undefined
		});
		expect(archiveResult).toEqual({ ok: true });
		expect(contents.documents[0].displayName).toBe('invoice.pdf');
	});

	it('uploads documents with FormData and no manual content-type', async () => {
		const fetchMock = vi.fn().mockResolvedValue(jsonResponse(documentMetadata(), { status: 201 }));
		const file = new File(['invoice'], 'invoice.pdf', { type: 'application/pdf' });

		const result = await uploadDocument(fetchMock, { folderId: 'folder/id', file });

		expect(fetchMock).toHaveBeenCalledWith('/api/documents/upload', {
			method: 'POST',
			headers: undefined,
			body: expect.any(FormData)
		});
		const body = fetchMock.mock.calls[0][1].body as FormData;
		expect(body.get('folderId')).toBe('folder/id');
		expect(body.get('file')).toBe(file);
		expect(result.id).toBe('document-id');
	});

	it('patches document metadata, archives documents, and builds download hrefs', async () => {
		const fetchMock = vi
			.fn()
			.mockResolvedValueOnce(
				jsonResponse(documentMetadata({ displayName: 'Invoice renamed.pdf' }), { status: 200 })
			)
			.mockResolvedValueOnce(jsonResponse({ ok: true }, { status: 200 }));

		const updated = await updateDocument(fetchMock, {
			id: 'document/id',
			input: { displayName: 'Invoice renamed.pdf' }
		});
		const archived = await archiveDocument(fetchMock, { id: 'document/id' });

		expect(fetchMock).toHaveBeenNthCalledWith(1, '/api/documents/document%2Fid', {
			method: 'PATCH',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({ displayName: 'Invoice renamed.pdf' })
		});
		expect(fetchMock).toHaveBeenNthCalledWith(2, '/api/documents/document%2Fid/archive', {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({})
		});
		expect(updated.displayName).toBe('Invoice renamed.pdf');
		expect(archived).toEqual({ ok: true });
		expect(documentDownloadHref('document/id')).toBe('/api/documents/document%2Fid/download');
	});

	it('uses public-safe fallback messages for server failures', async () => {
		const fetchMock = vi
			.fn()
			.mockResolvedValue(jsonResponse({ error: 'database unavailable' }, { status: 503 }));

		await expect(fetchDocumentFolderTree(fetchMock, 'unit-id')).rejects.toThrow(
			'Failed to load document folders'
		);
	});

	it('uses public error messages for client failures', async () => {
		const fetchMock = vi
			.fn()
			.mockResolvedValue(jsonResponse({ error: 'Display name is required' }, { status: 400 }));

		await expect(
			updateDocument(fetchMock, { id: 'document-id', input: { displayName: '' } })
		).rejects.toThrow('Display name is required');
	});

	it('rejects malformed successful folder payloads', async () => {
		const fetchMock = vi
			.fn()
			.mockResolvedValue(jsonResponse({ folders: [folderNode({ name: '' })] }, { status: 200 }));

		await expect(fetchDocumentFolderTree(fetchMock, 'unit-id')).rejects.toThrow(
			'Invalid document response'
		);
	});

	it('rejects documents with unsafe sizeBytes values', async () => {
		const fetchMock = vi.fn().mockResolvedValue(
			jsonResponse(
				contentsResponseWithDocument({ sizeBytes: Number.MAX_SAFE_INTEGER + 1 }),
				{ status: 200 }
			)
		);

		await expect(fetchDocumentFolderContents(fetchMock, 'folder-id')).rejects.toThrow(
			'Invalid document response'
		);
	});

	it('exports only the intended browser API surface', () => {
		const source = readFileSync(new URL('./api.ts', import.meta.url), 'utf8');
		const exportedNames = [
			...source.matchAll(/export\s+(?:async\s+)?(?:function|type|const|class)\s+(\w+)/g)
		]
			.map((match) => match[1])
			.sort();

		expect(exportedNames).toEqual([
			'ArchiveDocumentFolderResponse',
			'ArchiveDocumentFolderVariables',
			'ArchiveDocumentResponse',
			'ArchiveDocumentVariables',
			'DOCUMENT_FOLDERS_QUERY_KEY',
			'DOCUMENT_FOLDER_CONTENTS_QUERY_KEY',
			'DocumentFolderContentsResponse',
			'DocumentFolderInput',
			'DocumentFolderTreeResponse',
			'DocumentUpdateInput',
			'MoveDocumentFolderVariables',
			'UpdateDocumentFolderVariables',
			'UpdateDocumentVariables',
			'UploadDocumentVariables',
			'archiveDocument',
			'archiveDocumentFolder',
			'createDocumentFolder',
			'documentDownloadHref',
			'fetchDocumentFolderContents',
			'fetchDocumentFolderTree',
			'moveDocumentFolder',
			'updateDocument',
			'updateDocumentFolder',
			'uploadDocument'
		]);
	});
});

function contentsResponseWithDocument(
	overrides: Partial<RepositoryDocument>
): DocumentFolderContentsResponse {
	return {
		...contentsResponse(),
		documents: [documentMetadata(overrides)]
	};
}
