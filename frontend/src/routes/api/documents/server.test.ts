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
		archiveDocument: vi.fn(),
		downloadDocument: vi.fn(),
		getDocument: vi.fn(),
		updateDocument: vi.fn(),
		uploadDocument: vi.fn()
	};
});

vi.mock('$lib/server/documents', () => ({
	DocumentApiError: documentMocks.DocumentApiError,
	archiveDocument: documentMocks.archiveDocument,
	downloadDocument: documentMocks.downloadDocument,
	getDocument: documentMocks.getDocument,
	isDocumentApiError: (error: unknown) => error instanceof documentMocks.DocumentApiError,
	updateDocument: documentMocks.updateDocument,
	uploadDocument: documentMocks.uploadDocument
}));

import { PATCH as updateDocumentRoute, GET as getDocumentRoute } from './[id]/+server';
import { POST as archiveDocumentRoute } from './[id]/archive/+server';
import { GET as downloadDocumentRoute } from './[id]/download/+server';
import { POST as uploadDocumentRoute } from './upload/+server';

const cookieHeader = 'auth.session_token=token';
const metadata = {
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

describe('document Svelte API routes', () => {
	beforeEach(() => {
		vi.clearAllMocks();
	});

	it('requires authenticated users before proxying uploads', async () => {
		const response = await uploadDocumentRoute(
			event({
				body: uploadForm(),
				locals: { user: null, permissions: [] },
				method: 'POST',
				path: '/api/documents/upload'
			}) as never
		);

		expect(response.status).toBe(401);
		expect(await response.json()).toEqual({ error: 'Authentication required' });
		expect(documentMocks.uploadDocument).not.toHaveBeenCalled();
	});

	it('forwards upload FormData with cookie context', async () => {
		documentMocks.uploadDocument.mockResolvedValue(metadata);
		const fetchMock = vi.fn();

		const response = await uploadDocumentRoute(
			event({
				body: uploadForm(),
				fetch: fetchMock,
				locals: { user: { role: 'user' }, permissions: ['document.create'] },
				method: 'POST',
				path: '/api/documents/upload'
			}) as never
		);

		expect(response.status).toBe(201);
		expect(await response.json()).toEqual(metadata);
		expect(documentMocks.uploadDocument).toHaveBeenCalledWith(
			fetchMock,
			cookieHeader,
			expect.any(FormData)
		);
		const form = documentMocks.uploadDocument.mock.calls[0][2] as FormData;
		expect(form.get('folderId')).toBe('folder-id');
		expect(form.get('file')).toBeInstanceOf(File);
	});

	it('proxies get, update, and archive document routes with cookie context', async () => {
		documentMocks.getDocument.mockResolvedValue(metadata);
		documentMocks.updateDocument.mockResolvedValue({ ...metadata, displayName: 'Invoice renamed.pdf' });
		documentMocks.archiveDocument.mockResolvedValue({ ok: true });
		const fetchMock = vi.fn();

		await getDocumentRoute(
			event({
				fetch: fetchMock,
				locals: { user: { role: 'user' }, permissions: ['document.view'] },
				path: '/api/documents/document-id',
				params: { id: 'document-id' }
			}) as never
		);
		const updateResponse = await updateDocumentRoute(
			event({
				body: { displayName: 'Invoice renamed.pdf' },
				fetch: fetchMock,
				locals: { user: { role: 'user' }, permissions: ['document.update'] },
				method: 'PATCH',
				path: '/api/documents/document-id',
				params: { id: 'document-id' }
			}) as never
		);
		const archiveResponse = await archiveDocumentRoute(
			event({
				fetch: fetchMock,
				locals: { user: { role: 'user' }, permissions: ['document.delete'] },
				method: 'POST',
				path: '/api/documents/document-id/archive',
				params: { id: 'document-id' }
			}) as never
		);

		expect(documentMocks.getDocument).toHaveBeenCalledWith(fetchMock, cookieHeader, 'document-id');
		expect(documentMocks.updateDocument).toHaveBeenCalledWith(
			fetchMock,
			cookieHeader,
			'document-id',
			{ displayName: 'Invoice renamed.pdf' }
		);
		expect(documentMocks.archiveDocument).toHaveBeenCalledWith(fetchMock, cookieHeader, 'document-id');
		expect(await updateResponse.json()).toEqual({ ...metadata, displayName: 'Invoice renamed.pdf' });
		expect(await archiveResponse.json()).toEqual({ ok: true });
	});

	it('normalizes non-string document display names before proxying updates', async () => {
		documentMocks.updateDocument.mockResolvedValue({ ...metadata, displayName: '' });
		const fetchMock = vi.fn();

		await updateDocumentRoute(
			event({
				body: { displayName: 42 },
				fetch: fetchMock,
				locals: { user: { role: 'user' }, permissions: ['document.update'] },
				method: 'PATCH',
				path: '/api/documents/document-id',
				params: { id: 'document-id' }
			}) as never
		);

		expect(documentMocks.updateDocument).toHaveBeenCalledWith(
			fetchMock,
			cookieHeader,
			'document-id',
			{ displayName: '' }
		);
	});

	it('returns download bodies and preserves upstream download headers', async () => {
		documentMocks.downloadDocument.mockResolvedValue(
			new Response('download-body', {
				status: 200,
				headers: {
					'content-type': 'application/pdf',
					'content-disposition': 'attachment; filename="invoice.pdf"'
				}
			})
		);
		const fetchMock = vi.fn();

		const response = await downloadDocumentRoute(
			event({
				fetch: fetchMock,
				locals: { user: { role: 'user' }, permissions: ['document.download'] },
				path: '/api/documents/document-id/download',
				params: { id: 'document-id' }
			}) as never
		);

		expect(documentMocks.downloadDocument).toHaveBeenCalledWith(fetchMock, cookieHeader, 'document-id');
		expect(response.status).toBe(200);
		expect(response.headers.get('content-type')).toBe('application/pdf');
		expect(response.headers.get('content-disposition')).toBe('attachment; filename="invoice.pdf"');
		expect(await response.text()).toBe('download-body');
	});

	it('uses safe download header defaults when upstream headers are absent', async () => {
		documentMocks.downloadDocument.mockResolvedValue(
			new Response(new Uint8Array([1, 2, 3]), { status: 200 })
		);

		const response = await downloadDocumentRoute(
			event({
				locals: { user: { role: 'user' }, permissions: ['document.download'] },
				path: '/api/documents/document-id/download',
				params: { id: 'document-id' }
			}) as never
		);

		expect(response.headers.get('content-type')).toBe('application/octet-stream');
		expect(response.headers.get('content-disposition')).toBe('attachment');
	});

	it('maps known backend errors to public-safe JSON responses', async () => {
		documentMocks.updateDocument.mockRejectedValue(
			new documentMocks.DocumentApiError(503, 'database unavailable')
		);

		const response = await updateDocumentRoute(
			event({
				body: { displayName: 'Invoice renamed.pdf' },
				locals: { user: { role: 'user' }, permissions: ['document.update'] },
				method: 'PATCH',
				path: '/api/documents/document-id',
				params: { id: 'document-id' }
			}) as never
		);

		expect(response.status).toBe(502);
		expect(await response.json()).toEqual({ error: 'Failed to update document' });
	});
});

function uploadForm() {
	const form = new FormData();
	form.set('folderId', 'folder-id');
	form.set('file', new Blob(['invoice'], { type: 'application/pdf' }), 'invoice.pdf');
	return form;
}

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
	const isFormData = body instanceof FormData;
	return {
		fetch,
		locals,
		params,
		request: new Request(`http://localhost${path}`, {
			method,
			headers: {
				cookie: cookieHeader,
				...(body === undefined || isFormData ? {} : { 'content-type': 'application/json' })
			},
			body:
				body === undefined
					? undefined
					: isFormData
						? body
						: JSON.stringify(body)
		})
	};
}
