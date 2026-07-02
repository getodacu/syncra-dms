import { afterEach, describe, expect, it, vi } from 'vitest';

import {
	DocumentApiError,
	archiveDocument,
	archiveDocumentFolder,
	createDocumentFolder,
	downloadDocument,
	getDocument,
	getDocumentFolderContents,
	getDocumentFolderTree,
	isDocumentApiError,
	moveDocumentFolder,
	updateDocument,
	updateDocumentFolder,
	uploadDocument
} from './documents';
import type { DocumentFolder, DocumentMetadata } from './documents';

describe('server document client', () => {
	afterEach(() => {
		vi.unstubAllEnvs();
	});

	it('fetches folder contents with internal and cookie headers', async () => {
		stubAPIEnv();
		const fetch = vi.fn(async () =>
			jsonResponse({
				folder: folder(),
				folders: [],
				documents: []
			})
		);

		const result = await getDocumentFolderContents(fetch, 'session=token', 'folder/id');

		expect(fetch).toHaveBeenCalledWith(
			'http://api.test/api/document-folders/folder%2Fid/contents',
			expect.objectContaining({
				method: 'GET',
				headers: expect.any(Headers)
			})
		);
		expect(result.folder.id).toBe('folder-id');
		const firstCall = fetch.mock.calls[0] as unknown as [string, RequestInit];
		const headers = firstCall[1].headers as Headers;
		expect(headers.get('X-Syncra-Internal-Token')).toBe('internal-token');
		expect(headers.get('cookie')).toBe('session=token');
	});

	it('fetches a document folder tree with an encoded organization unit query', async () => {
		stubAPIEnv();
		const root = folder({ children: [folder({ id: 'child-id', parentId: 'folder-id' })] });
		const fetch = vi.fn(async () => jsonResponse({ folders: [root] }));

		const result = await getDocumentFolderTree(fetch, 'session=token', 'unit/id');

		expect(fetch).toHaveBeenCalledWith(
			'http://api.test/api/document-folders/tree?organizationUnitId=unit%2Fid',
			expect.objectContaining({ method: 'GET' })
		);
		expect(result.folders[0].children?.[0].id).toBe('child-id');
	});

	it('uses encoded ids and JSON payloads for folder create, update, move, and archive', async () => {
		stubAPIEnv();
		const responses: unknown[] = [
			folder({ id: 'created-folder' }),
			folder({ id: 'folder/id', name: 'Invoices Updated' }),
			folder({ id: 'folder/id', parentId: 'parent/id' }),
			{ ok: true }
		];
		const fetch = vi.fn(async () => jsonResponse(responses.shift()));
		const input = {
			organizationUnitId: 'unit-id',
			parentId: null,
			name: 'Invoices',
			description: 'Monthly invoices'
		};

		await createDocumentFolder(fetch, null, input);
		await updateDocumentFolder(fetch, null, 'folder/id', {
			organizationUnitId: 'unit-id',
			name: 'Invoices Updated',
			description: null
		});
		await moveDocumentFolder(fetch, null, 'folder/id', 'parent/id');
		const archive = await archiveDocumentFolder(fetch, null, 'folder/id');

		expect(fetch).toHaveBeenNthCalledWith(1, 'http://api.test/api/document-folders', {
			method: 'POST',
			headers: expect.any(Headers),
			body: JSON.stringify(input)
		});
		expect(fetch).toHaveBeenNthCalledWith(
			2,
			'http://api.test/api/document-folders/folder%2Fid',
			{
				method: 'PATCH',
				headers: expect.any(Headers),
				body: JSON.stringify({
					organizationUnitId: 'unit-id',
					name: 'Invoices Updated',
					description: null
				})
			}
		);
		expect(fetch).toHaveBeenNthCalledWith(
			3,
			'http://api.test/api/document-folders/folder%2Fid/parent',
			{
				method: 'PATCH',
				headers: expect.any(Headers),
				body: JSON.stringify({ parentId: 'parent/id' })
			}
		);
		expect(fetch).toHaveBeenNthCalledWith(
			4,
			'http://api.test/api/document-folders/folder%2Fid/archive',
			{
				method: 'POST',
				headers: expect.any(Headers),
				body: JSON.stringify({})
			}
		);
		expect(archive).toEqual({ ok: true });
		const headers = (fetch.mock.calls[0] as unknown as [string, RequestInit])[1].headers as Headers;
		expect(headers.get('content-type')).toBe('application/json');
	});

	it('uploads a document by forwarding FormData without a manual content-type', async () => {
		stubAPIEnv();
		const form = new FormData();
		form.set('folderId', 'folder-id');
		form.set('file', new Blob(['invoice'], { type: 'application/pdf' }), 'invoice.pdf');
		const fetch = vi.fn(async () => jsonResponse(documentMetadata()));

		const result = await uploadDocument(fetch, 'session=token', form);

		expect(result.id).toBe('document-id');
		expect(fetch).toHaveBeenCalledWith('http://api.test/api/documents/upload', {
			method: 'POST',
			headers: expect.any(Headers),
			body: form
		});
		const headers = (fetch.mock.calls[0] as unknown as [string, RequestInit])[1].headers as Headers;
		expect(headers.get('X-Syncra-Internal-Token')).toBe('internal-token');
		expect(headers.get('cookie')).toBe('session=token');
		expect(headers.has('content-type')).toBe(false);
	});

	it('gets, updates, and archives document metadata with expected request shapes', async () => {
		stubAPIEnv();
		const responses: unknown[] = [
			documentMetadata(),
			documentMetadata({ displayName: 'Invoice renamed.pdf' }),
			{ ok: true }
		];
		const fetch = vi.fn(async () => jsonResponse(responses.shift()));

		const metadata = await getDocument(fetch, null, 'document/id');
		const updated = await updateDocument(fetch, null, 'document/id', {
			displayName: 'Invoice renamed.pdf'
		});
		const archived = await archiveDocument(fetch, null, 'document/id');

		expect(metadata.id).toBe('document-id');
		expect(updated.displayName).toBe('Invoice renamed.pdf');
		expect(archived).toEqual({ ok: true });
		expect(fetch).toHaveBeenNthCalledWith(
			1,
			'http://api.test/api/documents/document%2Fid',
			expect.objectContaining({ method: 'GET' })
		);
		expect(fetch).toHaveBeenNthCalledWith(2, 'http://api.test/api/documents/document%2Fid', {
			method: 'PATCH',
			headers: expect.any(Headers),
			body: JSON.stringify({ displayName: 'Invoice renamed.pdf' })
		});
		expect(fetch).toHaveBeenNthCalledWith(
			3,
			'http://api.test/api/documents/document%2Fid/archive',
			{
				method: 'POST',
				headers: expect.any(Headers),
				body: JSON.stringify({})
			}
		);
	});

	it('returns the raw upstream Response for downloads after checking ok', async () => {
		stubAPIEnv();
		const upstream = new Response('download-body', {
			status: 200,
			headers: {
				'content-type': 'application/pdf',
				'content-disposition': 'attachment; filename="invoice.pdf"'
			}
		});
		const fetch = vi.fn(async () => upstream);

		const result = await downloadDocument(fetch, 'session=token', 'document/id');

		expect(result).toBe(upstream);
		expect(fetch).toHaveBeenCalledWith(
			'http://api.test/api/documents/document%2Fid/download',
			expect.objectContaining({ method: 'GET', headers: expect.any(Headers) })
		);
		expect(await result.text()).toBe('download-body');
	});

	it('rejects invalid success payloads at the boundary', async () => {
		stubAPIEnv();
		const invalidFetch = vi.fn(async () => jsonResponse({ folders: [{ id: 'folder-id' }] }));

		await expect(getDocumentFolderTree(invalidFetch, 'session=token', 'unit-id')).rejects.toMatchObject(
			new DocumentApiError(502, 'Invalid Document response')
		);
	});

	it('rejects document metadata sizeBytes values that are not safe byte counts', async () => {
		stubAPIEnv();

		for (const sizeBytes of [1.5, -1, Number.MAX_SAFE_INTEGER + 1]) {
			const fetch = vi.fn(async () => jsonResponse(documentMetadata({ sizeBytes })));

			await expect(getDocument(fetch, 'session=token', 'document-id')).rejects.toMatchObject(
				new DocumentApiError(502, 'Invalid Document response')
			);
		}
	});

	it('maps public backend errors and supports DocumentApiError narrowing', async () => {
		stubAPIEnv();
		const forbiddenFetch = vi.fn(async () => jsonResponse({ error: 'document.view required' }, 403));

		try {
			await getDocumentFolderContents(forbiddenFetch, 'session=token', 'folder-id');
			throw new Error('Expected getDocumentFolderContents to reject');
		} catch (error) {
			expect(error).toMatchObject(new DocumentApiError(403, 'document.view required'));
			expect(isDocumentApiError(error)).toBe(true);
		}

		const serverErrorFetch = vi.fn(async () => jsonResponse({ error: 'database exploded' }, 500));
		await expect(getDocument(serverErrorFetch, 'session=token', 'document-id')).rejects.toMatchObject(
			new DocumentApiError(502, 'Document request failed')
		);
	});

	it('maps download failures without consuming successful download bodies', async () => {
		stubAPIEnv();
		const fetch = vi.fn(async () => jsonResponse({ error: 'document file not found' }, 404));

		await expect(downloadDocument(fetch, 'session=token', 'document-id')).rejects.toMatchObject(
			new DocumentApiError(404, 'document file not found')
		);
	});
});

function stubAPIEnv() {
	vi.stubEnv('SYNCRA_API_BASE_URL', 'http://api.test');
	vi.stubEnv('SYNCRA_INTERNAL_API_TOKEN', 'internal-token');
}

function folder(overrides: Partial<DocumentFolder> = {}): DocumentFolder {
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

function documentMetadata(overrides: Partial<DocumentMetadata> = {}): DocumentMetadata {
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

function jsonResponse(body: unknown, status = 200) {
	return new Response(JSON.stringify(body), {
		status,
		headers: { 'content-type': 'application/json' }
	});
}
