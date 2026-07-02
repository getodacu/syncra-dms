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

		expect(items[0].id).toMatch(/^upload-\d+-0-invoice-pdf-7-1700000000000$/);
		expect(items[1].id).toMatch(/^upload-\d+-1-invoice-pdf-7-1700000000000$/);
		expect(new Set(items.map((item) => item.id)).size).toBe(2);
		expect(batchId(items[0].id)).toBe(batchId(items[1].id));
	});

	it('keeps fallback ids unique across separate file selections', () => {
		vi.stubGlobal('crypto', {});
		const firstSelection = filesToUploadItems([
			file('invoice.pdf', 'invoice', 1700000000000)
		]);
		const secondSelection = filesToUploadItems([
			file('invoice.pdf', 'invoice', 1700000000000)
		]);
		const items = [...firstSelection, ...secondSelection];

		expect(items[0].id).not.toBe(items[1].id);
		expect(new Set(items.map((item) => item.id)).size).toBe(2);

		const uploading = markUploadUploading(items, items[1].id);

		expect(uploading[0]).toEqual(items[0]);
		expect(uploading[1]).toMatchObject({
			id: items[1].id,
			status: 'uploading'
		});
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

function batchId(id: string) {
	return id.split('-')[1];
}
