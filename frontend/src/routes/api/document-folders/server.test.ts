import { beforeEach, describe, expect, it, vi } from 'vitest';

const documentMocks = vi.hoisted(() => {
	class MockDocumentApiError extends Error {
		status: number;

		constructor(status: number, message: string) {
			super(message);
			this.name = 'DocumentApiError';
			this.status = status;
		}
	}

	return {
		DocumentApiError: MockDocumentApiError,
		archiveDocumentFolder: vi.fn(),
		createDocumentFolder: vi.fn(),
		getDocumentFolderContents: vi.fn(),
		getDocumentFolderTree: vi.fn(),
		moveDocumentFolder: vi.fn(),
		updateDocumentFolder: vi.fn()
	};
});

vi.mock('$lib/server/documents', () => ({
	DocumentApiError: documentMocks.DocumentApiError,
	archiveDocumentFolder: documentMocks.archiveDocumentFolder,
	createDocumentFolder: documentMocks.createDocumentFolder,
	getDocumentFolderContents: documentMocks.getDocumentFolderContents,
	getDocumentFolderTree: documentMocks.getDocumentFolderTree,
	isDocumentApiError: (error: unknown) => error instanceof documentMocks.DocumentApiError,
	moveDocumentFolder: documentMocks.moveDocumentFolder,
	updateDocumentFolder: documentMocks.updateDocumentFolder
}));

import { POST as createDocumentFolderRoute } from './+server';
import { PATCH as updateDocumentFolderRoute } from './[id]/+server';
import { POST as archiveDocumentFolderRoute } from './[id]/archive/+server';
import { GET as getDocumentFolderContentsRoute } from './[id]/contents/+server';
import { PATCH as moveDocumentFolderRoute } from './[id]/parent/+server';
import { GET as getDocumentFolderTreeRoute } from './tree/+server';
import { hasAnyPermission } from '../documents/api.server';

const cookieHeader = 'auth.session_token=token';
const folder = {
	id: 'folder-id',
	parentId: null,
	organizationUnitId: 'unit-id',
	name: 'Invoices',
	description: null,
	deletedAt: null,
	createdAt: '2026-07-01T10:00:00Z',
	updatedAt: '2026-07-01T10:00:00Z',
	children: []
};
const document = {
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
	updatedAt: '2026-07-01T10:00:00Z'
};

describe('document folder Svelte API routes', () => {
	beforeEach(() => {
		vi.clearAllMocks();
	});

	it('requires authenticated users before proxying folder reads', async () => {
		const response = await getDocumentFolderTreeRoute(
			event({
				locals: { user: null, permissions: [] },
				path: '/api/document-folders/tree?organizationUnitId=unit-id'
			}) as never
		);

		expect(response.status).toBe(401);
		expect(await response.json()).toEqual({ error: 'Authentication required' });
		expect(documentMocks.getDocumentFolderTree).not.toHaveBeenCalled();
	});

	it('supports document permissions including system admins', () => {
		expect(hasAnyPermission(['document.view'], ['document.view'])).toBe(true);
		expect(hasAnyPermission(['system.admin'], ['document.delete'])).toBe(true);
		expect(hasAnyPermission(['document.view'], ['document.create'])).toBe(false);
	});

	it('proxies tree reads with organization unit query and cookie context', async () => {
		documentMocks.getDocumentFolderTree.mockResolvedValue({ folders: [folder] });
		const fetchMock = vi.fn();

		const response = await getDocumentFolderTreeRoute(
			event({
				fetch: fetchMock,
				locals: { user: { role: 'user' }, permissions: ['document.view'] },
				path: '/api/document-folders/tree?organizationUnitId=unit-id'
			}) as never
		);

		expect(response.status).toBe(200);
		expect(await response.json()).toEqual({ folders: [folder] });
		expect(documentMocks.getDocumentFolderTree).toHaveBeenCalledWith(
			fetchMock,
			cookieHeader,
			'unit-id'
		);
	});

	it('rejects missing organization unit tree queries before proxying', async () => {
		const response = await getDocumentFolderTreeRoute(
			event({
				locals: { user: { role: 'user' }, permissions: ['document.view'] },
				path: '/api/document-folders/tree?organizationUnitId=  '
			}) as never
		);

		expect(response.status).toBe(400);
		expect(await response.json()).toEqual({ error: 'organizationUnitId is required' });
		expect(documentMocks.getDocumentFolderTree).not.toHaveBeenCalled();
	});

	it('proxies create and update requests with normalized JSON and cookie context', async () => {
		documentMocks.createDocumentFolder.mockResolvedValue(folder);
		documentMocks.updateDocumentFolder.mockResolvedValue({ ...folder, name: 'Invoices Updated' });
		const fetchMock = vi.fn();

		await createDocumentFolderRoute(
			event({
				fetch: fetchMock,
				locals: { user: { role: 'user' }, permissions: ['document.create'] },
				method: 'POST',
				path: '/api/document-folders',
				body: {
					organizationUnitId: 'unit-id',
					name: 'Invoices',
					description: 12
				}
			}) as never
		);
		const updateResponse = await updateDocumentFolderRoute(
			event({
				fetch: fetchMock,
				locals: { user: { role: 'user' }, permissions: ['document.update'] },
				method: 'PATCH',
				path: '/api/document-folders/folder-id',
				params: { id: 'folder-id' },
				body: {
					organizationUnitId: 'unit-id',
					parentId: 'parent-id',
					name: 'Invoices Updated',
					description: ''
				}
			}) as never
		);

		expect(documentMocks.createDocumentFolder).toHaveBeenCalledWith(fetchMock, cookieHeader, {
			organizationUnitId: 'unit-id',
			parentId: null,
			name: 'Invoices',
			description: null
		});
		expect(documentMocks.updateDocumentFolder).toHaveBeenCalledWith(
			fetchMock,
			cookieHeader,
			'folder-id',
			{
				organizationUnitId: 'unit-id',
				parentId: 'parent-id',
				name: 'Invoices Updated',
				description: ''
			}
		);
		expect(await updateResponse.json()).toEqual({ ...folder, name: 'Invoices Updated' });
	});

	it('proxies move, archive, and contents requests with cookie context', async () => {
		documentMocks.moveDocumentFolder.mockResolvedValue({ ...folder, parentId: null });
		documentMocks.archiveDocumentFolder.mockResolvedValue({ ok: true });
		documentMocks.getDocumentFolderContents.mockResolvedValue({
			folder,
			folders: [],
			documents: [document]
		});
		const fetchMock = vi.fn();

		await moveDocumentFolderRoute(
			event({
				fetch: fetchMock,
				locals: { user: { role: 'user' }, permissions: ['document.update'] },
				method: 'PATCH',
				path: '/api/document-folders/folder-id/parent',
				params: { id: 'folder-id' },
				body: { parentId: '  ' }
			}) as never
		);
		await archiveDocumentFolderRoute(
			event({
				fetch: fetchMock,
				locals: { user: { role: 'user' }, permissions: ['document.delete'] },
				method: 'POST',
				path: '/api/document-folders/folder-id/archive',
				params: { id: 'folder-id' }
			}) as never
		);
		const contentsResponse = await getDocumentFolderContentsRoute(
			event({
				fetch: fetchMock,
				locals: { user: { role: 'user' }, permissions: ['document.view'] },
				path: '/api/document-folders/folder-id/contents',
				params: { id: 'folder-id' }
			}) as never
		);

		expect(documentMocks.moveDocumentFolder).toHaveBeenCalledWith(
			fetchMock,
			cookieHeader,
			'folder-id',
			null
		);
		expect(documentMocks.archiveDocumentFolder).toHaveBeenCalledWith(
			fetchMock,
			cookieHeader,
			'folder-id'
		);
		expect(documentMocks.getDocumentFolderContents).toHaveBeenCalledWith(
			fetchMock,
			cookieHeader,
			'folder-id'
		);
		expect(await contentsResponse.json()).toEqual({ folder, folders: [], documents: [document] });
	});

	it('rejects parent moves without a parentId field before proxying', async () => {
		const response = await moveDocumentFolderRoute(
			event({
				locals: { user: { role: 'user' }, permissions: ['document.update'] },
				method: 'PATCH',
				path: '/api/document-folders/folder-id/parent',
				params: { id: 'folder-id' },
				body: {}
			}) as never
		);

		expect(response.status).toBe(400);
		expect(await response.json()).toEqual({ error: 'parentId is required' });
		expect(documentMocks.moveDocumentFolder).not.toHaveBeenCalled();
	});

	it('maps known backend errors to public-safe JSON responses', async () => {
		documentMocks.getDocumentFolderTree.mockRejectedValue(
			new documentMocks.DocumentApiError(503, 'database unavailable')
		);

		const response = await getDocumentFolderTreeRoute(
			event({
				locals: { user: { role: 'user' }, permissions: ['document.view'] },
				path: '/api/document-folders/tree?organizationUnitId=unit-id'
			}) as never
		);

		expect(response.status).toBe(502);
		expect(await response.json()).toEqual({ error: 'Failed to load document folders' });
	});
});

function event({
	fetch = vi.fn(),
	locals,
	method = 'GET',
	path,
	params = {},
	body
}: {
	fetch?: typeof globalThis.fetch;
	locals: { user: { role: string } | null; permissions: string[] };
	method?: string;
	path: string;
	params?: Record<string, string>;
	body?: unknown;
}) {
	const url = new URL(`http://localhost${path}`);
	return {
		fetch,
		locals,
		params,
		request: new Request(url, {
			method,
			headers: {
				cookie: cookieHeader,
				...(body === undefined ? {} : { 'content-type': 'application/json' })
			},
			body: body === undefined ? undefined : JSON.stringify(body)
		}),
		url
	};
}
