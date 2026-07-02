import { afterEach, describe, expect, it, vi } from 'vitest';

import {
	filesToUploadItems,
	markUploadFailed,
	markUploadUploaded,
	markUploadUploading
} from './upload-queue';

describe('document upload queue utilities', () => {
	afterEach(() => {
		vi.unstubAllGlobals();
	});

	it('creates queued upload items with crypto random ids when available', () => {
		const randomUUID = vi
			.fn()
			.mockReturnValueOnce('00000000-0000-4000-8000-000000000001')
			.mockReturnValueOnce('00000000-0000-4000-8000-000000000002');
		vi.stubGlobal('crypto', { randomUUID });
		const first = file('invoice.pdf', 'invoice');
		const second = file('receipt.pdf', 'receipt');

		const items = filesToUploadItems([first, second]);

		expect(items).toEqual([
			{
				id: '00000000-0000-4000-8000-000000000001',
				file: first,
				status: 'queued'
			},
			{
				id: '00000000-0000-4000-8000-000000000002',
				file: second,
				status: 'queued'
			}
		]);
		expect(randomUUID).toHaveBeenCalledTimes(2);
	});

	it('uses deterministic fallback ids when crypto random ids are unavailable', () => {
		vi.stubGlobal('crypto', {});
		const first = file('invoice.pdf', 'invoice', 1700000000000);
		const second = file('invoice.pdf', 'invoice', 1700000000000);

		const items = filesToUploadItems([first, second]);

		expect(items.map((item) => item.id)).toEqual([
			'upload-0-invoice-pdf-7-1700000000000',
			'upload-1-invoice-pdf-7-1700000000000'
		]);
	});

	it('transitions upload item status while preserving item ids', () => {
		vi.stubGlobal('crypto', {
			randomUUID: vi
				.fn()
				.mockReturnValueOnce('00000000-0000-4000-8000-000000000001')
				.mockReturnValueOnce('00000000-0000-4000-8000-000000000002')
		});
		const items = filesToUploadItems([file('invoice.pdf', 'invoice'), file('receipt.pdf', 'receipt')]);

		const uploading = markUploadUploading(items, items[0].id);
		const failed = markUploadFailed(uploading, items[0].id, 'Network error');
		const uploaded = markUploadUploaded(failed, items[0].id, 'document-id');

		expect(uploading[0]).toMatchObject({ id: items[0].id, status: 'uploading' });
		expect(uploading[1]).toBe(items[1]);
		expect(failed[0]).toMatchObject({
			id: items[0].id,
			status: 'failed',
			error: 'Network error'
		});
		expect(uploaded[0]).toEqual({
			id: items[0].id,
			file: items[0].file,
			status: 'uploaded',
			documentId: 'document-id'
		});
		expect(uploaded[1]).toBe(items[1]);
	});
});

function file(name: string, contents: string, lastModified = 1700000000001) {
	return new File([contents], name, {
		type: 'application/pdf',
		lastModified
	});
}
