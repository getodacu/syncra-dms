export type UploadQueueItem = {
	id: string;
	file: File;
	status: 'queued' | 'uploading' | 'uploaded' | 'failed';
	error?: string;
	documentId?: string;
};

export function filesToUploadItems(files: Iterable<File> | ArrayLike<File>): UploadQueueItem[] {
	return Array.from(files, (file, index) => ({
		id: uploadItemId(file, index),
		file,
		status: 'queued' as const
	}));
}

export function markUploadUploading(items: UploadQueueItem[], id: string): UploadQueueItem[] {
	return updateUploadItem(items, id, (item) => ({
		id: item.id,
		file: item.file,
		status: 'uploading'
	}));
}

export function markUploadUploaded(
	items: UploadQueueItem[],
	id: string,
	documentId: string
): UploadQueueItem[] {
	return updateUploadItem(items, id, (item) => ({
		id: item.id,
		file: item.file,
		status: 'uploaded',
		documentId
	}));
}

export function markUploadFailed(
	items: UploadQueueItem[],
	id: string,
	error: string
): UploadQueueItem[] {
	return updateUploadItem(items, id, (item) => ({
		id: item.id,
		file: item.file,
		status: 'failed',
		error
	}));
}

function updateUploadItem(
	items: UploadQueueItem[],
	id: string,
	update: (item: UploadQueueItem) => UploadQueueItem
) {
	return items.map((item) => (item.id === id ? update(item) : item));
}

function uploadItemId(file: File, index: number) {
	const randomUUID = globalThis.crypto?.randomUUID;
	if (typeof randomUUID === 'function') {
		return randomUUID.call(globalThis.crypto);
	}

	return `upload-${index}-${fileNameToken(file.name)}-${file.size}-${file.lastModified}`;
}

function fileNameToken(name: string) {
	const token = name
		.trim()
		.toLowerCase()
		.replace(/[^a-z0-9]+/g, '-')
		.replace(/^-+|-+$/g, '');
	return token || 'file';
}
