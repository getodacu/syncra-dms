import { readFileSync } from 'node:fs';
import { describe, expect, it } from 'vitest';

import { load } from './+page.server';

describe('documents page server load', () => {
	it('exposes document permission flags from locals', () => {
		const result = load(loadEvent(['document.view', 'document.create'], 'unit-id') as never);

		expect(result).toEqual({
			canViewDocuments: true,
			canCreateDocuments: true,
			canUpdateDocuments: false,
			canDeleteDocuments: false,
			canDownloadDocuments: false,
			selectedOrganizationUnitId: 'unit-id'
		});
	});

	it('returns full document access for system administrators', () => {
		const result = load(loadEvent(['system.admin']) as never);

		expect(result).toEqual({
			canViewDocuments: true,
			canCreateDocuments: true,
			canUpdateDocuments: true,
			canDeleteDocuments: true,
			canDownloadDocuments: true,
			selectedOrganizationUnitId: null
		});
	});

	it('returns no document access without document permissions, regardless of auth role', () => {
		const result = load(loadEvent([], undefined, { role: 'admin' }) as never);

		expect(result).toEqual({
			canViewDocuments: false,
			canCreateDocuments: false,
			canUpdateDocuments: false,
			canDeleteDocuments: false,
			canDownloadDocuments: false,
			selectedOrganizationUnitId: null
		});
	});
});

describe('documents page source', () => {
	it('uses TanStack query and mutations against the document API wrapper', () => {
		const source = readFileSync(new URL('./+page.svelte', import.meta.url), 'utf8');

		expect(source).toContain("import type { PageProps } from './$types'");
		expect(source).toContain("import { createMutation, createQuery, useQueryClient }");
		expect(source).toContain('fetchDocumentFolderTree');
		expect(source).toContain('fetchDocumentFolderContents');
		expect(source).toContain('uploadDocument');
		expect(source).toContain('createDocumentFolder');
		expect(source).toContain('updateDocumentFolder');
		expect(source).toContain('moveDocumentFolder');
		expect(source).toContain('archiveDocumentFolder');
		expect(source).toContain('updateDocument');
		expect(source).toContain('archiveDocument');
		expect(source).toContain('documentFoldersQueryKey(selectedOrganizationUnitId)');
		expect(source).toContain('documentFolderContentsQueryKey(selectedFolderId)');
		expect(source).toContain('documentFoldersQueryKey(organizationUnitId)');
		expect(source).toContain('documentFolderContentsQueryKey(folderId)');
		expect(source).toContain(
			'queryClient.invalidateQueries({ queryKey: documentFoldersQueryKey(organizationUnitId) })'
		);
		expect(source).toContain(
			'queryClient.invalidateQueries({ queryKey: documentFolderContentsQueryKey(folderId) })'
		);
		expect(source).toContain(
			'queryClient.invalidateQueries({ queryKey: DOCUMENT_FOLDER_CONTENTS_QUERY_KEY })'
		);
		expect(source).toContain('filesToUploadItems');
		expect(source).toContain('markUploadUploading');
		expect(source).toContain('markUploadUploaded');
		expect(source).toContain('markUploadFailed');
		expect(source).toContain('canViewDocuments');
		expect(source).toContain('canCreateDocuments');
		expect(source).toContain('canUpdateDocuments');
		expect(source).toContain('canDeleteDocuments');
		expect(source).toContain('canDownloadDocuments');
		expect(source).toContain('selectedOrganizationUnitId');
		expect(source).toContain('fetchOrganizationUnitTree');
		expect(source).toContain('FolderTree');
		expect(source).toContain('RepositoryTable');
		expect(source).toContain('UploadPanel');
		expect(source).toContain('const activeRepositoryCount = $derived.by');
		expect(source).toContain('const activeRepositoryCountLabel = $derived.by');
		expect(source).toContain('{activeRepositoryCountLabel}');
		expect(source).toContain('active items');
		expect(source).toContain('active folders');
		expect(source).toContain(
			'{#if pageData.canCreateDocuments && folderTree.length === 0 && !folderTreeQuery.isLoading}'
		);
		expect(source).toContain('submitRootFolder');
		expect(source).toContain('Create root');
		expect(source).toContain('organizationUnitsQuery.isError');
		expect(source).toContain('folderTreeQuery.isError');
		expect(source).toContain('folderContentsQuery.isError');
		expect(source).toContain('Loading folders');
		expect(source).toContain('No document access');
		expect(source).toContain('Select an organization unit to open the repository.');
		expect(source).not.toContain('$lib/server/documents');
		expect(source).not.toContain("from '$lib/server");
	});

	it('implements document repository component contracts', () => {
		const folderTree = readFileSync(new URL('./folder-tree.svelte', import.meta.url), 'utf8');
		const repositoryTable = readFileSync(
			new URL('./repository-table.svelte', import.meta.url),
			'utf8'
		);
		const uploadPanel = readFileSync(new URL('./upload-panel.svelte', import.meta.url), 'utf8');

		expect(folderTree).toContain('folders: FlatDocumentFolderNode[]');
		expect(folderTree).toContain('selectedId: string | null');
		expect(folderTree).toContain('onSelect: (id: string) => void');
		expect(folderTree).toContain('FolderIcon');
		expect(folderTree).toContain('role="tree"');
		expect(folderTree).toContain('aria-level={folder.depth + 1}');
		expect(folderTree).toContain('padding-left');
		expect(folderTree).toContain('No folders');

		expect(repositoryTable).toContain('rows: RepositoryRow[]');
		expect(repositoryTable).toContain('canUpdate: boolean');
		expect(repositoryTable).toContain('canDelete: boolean');
		expect(repositoryTable).toContain('canDownload: boolean');
		expect(repositoryTable).toContain('onOpenFolder: (id: string) => void');
		expect(repositoryTable).toContain('onRenameFolder: (id: string, name: string) => Promise<void>');
		expect(repositoryTable).toContain('onArchiveFolder: (id: string) => Promise<void>');
		expect(repositoryTable).toContain(
			'onRenameDocument: (id: string, displayName: string) => Promise<void>'
		);
		expect(repositoryTable).toContain('onArchiveDocument: (id: string) => Promise<void>');
		expect(repositoryTable).toContain('documentDownloadHref(row.id)');
		expect(repositoryTable).toContain('confirm(');
		expect(repositoryTable).toContain('<table');
		expect(repositoryTable).toContain('DownloadIcon');
		expect(repositoryTable).toContain('PencilIcon');
		expect(repositoryTable).toContain('ArchiveIcon');

		expect(uploadPanel).toContain('canCreate: boolean');
		expect(uploadPanel).toContain('selectedFolderId: string | null');
		expect(uploadPanel).toContain('queue: UploadQueueItem[]');
		expect(uploadPanel).toContain('isUploading: boolean');
		expect(uploadPanel).toContain('onFilesSelected: (files: FileList) => void');
		expect(uploadPanel).toContain('onUpload: () => Promise<void>');
		expect(uploadPanel).toContain('multiple');
		expect(uploadPanel).toContain('queued');
		expect(uploadPanel).toContain('uploading');
		expect(uploadPanel).toContain('uploaded');
		expect(uploadPanel).toContain('failed');
	});
});

function loadEvent(
	permissions: string[],
	organizationUnitId?: string,
	user: { role: string } | null = { role: 'user' }
) {
	const url = new URL('http://localhost/app/documents');
	if (organizationUnitId) url.searchParams.set('organizationUnitId', organizationUnitId);
	return {
		locals: { user, permissions },
		url
	};
}
