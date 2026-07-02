import { describe, expect, it } from 'vitest';

import {
	collectFolderMoveTargets,
	findFolder,
	flattenFolderTree,
	repositoryRows,
	selectInitialFolder,
	type DocumentFolderNode,
	type RepositoryDocument
} from './tree';

const folders: DocumentFolderNode[] = [
	{
		id: 'root-a',
		parentId: null,
		organizationUnitId: 'unit-id',
		name: 'Root A',
		description: 'First root',
		deletedAt: null,
		createdAt: '2026-01-01T00:00:00Z',
		updatedAt: '2026-01-02T00:00:00Z',
		children: [
			{
				id: 'child-a1',
				parentId: 'root-a',
				organizationUnitId: 'unit-id',
				name: 'Child A1',
				description: null,
				deletedAt: null,
				createdAt: '2026-01-03T00:00:00Z',
				updatedAt: '2026-01-04T00:00:00Z',
				children: [
					{
						id: 'grandchild-a1a',
						parentId: 'child-a1',
						organizationUnitId: 'unit-id',
						name: 'Grandchild A1A',
						description: null,
						deletedAt: null,
						createdAt: '2026-01-05T00:00:00Z',
						updatedAt: '2026-01-06T00:00:00Z',
						children: []
					}
				]
			},
			{
				id: 'child-a2',
				parentId: 'root-a',
				organizationUnitId: 'unit-id',
				name: 'Child A2',
				description: null,
				deletedAt: null,
				createdAt: '2026-01-07T00:00:00Z',
				updatedAt: '2026-01-08T00:00:00Z',
				children: []
			}
		]
	},
	{
		id: 'root-b',
		parentId: null,
		organizationUnitId: 'unit-id',
		name: 'Root B',
		description: null,
		deletedAt: null,
		createdAt: '2026-01-09T00:00:00Z',
		updatedAt: '2026-01-10T00:00:00Z',
		children: [
			{
				id: 'child-b1',
				parentId: 'root-b',
				organizationUnitId: 'unit-id',
				name: 'Child B1',
				description: null,
				deletedAt: null,
				createdAt: '2026-01-11T00:00:00Z',
				updatedAt: '2026-01-12T00:00:00Z',
				children: []
			}
		]
	}
];

function document(overrides: Partial<RepositoryDocument> = {}): RepositoryDocument {
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

describe('document folder tree utilities', () => {
	it('flattens folders in pre-order with depth values', () => {
		const rows = flattenFolderTree(folders);

		expect(rows.map(({ id, depth }) => [id, depth])).toEqual([
			['root-a', 0],
			['child-a1', 1],
			['grandchild-a1a', 2],
			['child-a2', 1],
			['root-b', 0],
			['child-b1', 1]
		]);
		expect(rows[0]).toMatchObject({
			id: 'root-a',
			name: 'Root A',
			description: 'First root',
			depth: 0
		});
		expect('children' in rows[0]).toBe(false);
	});

	it('finds nested folders and returns null for nil or unknown ids', () => {
		expect(findFolder(folders, 'child-b1')?.name).toBe('Child B1');
		expect(findFolder(folders, null)).toBeNull();
		expect(findFolder(folders, undefined)).toBeNull();
		expect(findFolder(folders, 'missing')).toBeNull();
	});

	it('selects the requested folder before falling back to the first root', () => {
		expect(selectInitialFolder(folders, 'grandchild-a1a')?.id).toBe('grandchild-a1a');
		expect(selectInitialFolder(folders, 'missing')?.id).toBe('root-a');
		expect(selectInitialFolder(folders)?.id).toBe('root-a');
		expect(selectInitialFolder([], 'root-a')).toBeNull();
	});

	it('collects move targets without the selected folder or its descendants', () => {
		const rows = collectFolderMoveTargets(folders, 'child-a1');

		expect(rows.map(({ id, depth }) => [id, depth])).toEqual([
			['root-a', 0],
			['child-a2', 1],
			['root-b', 0],
			['child-b1', 1]
		]);
	});

	it('collects move targets without deleted branches', () => {
		const rows = collectFolderMoveTargets(
			[
				{
					id: 'active-root',
					parentId: null,
					organizationUnitId: 'unit-id',
					name: 'Active Root',
					description: null,
					deletedAt: null,
					createdAt: '2026-01-01T00:00:00Z',
					updatedAt: '2026-01-02T00:00:00Z',
					children: [
						{
							id: 'deleted-child',
							parentId: 'active-root',
							organizationUnitId: 'unit-id',
							name: 'Deleted Child',
							description: null,
							deletedAt: '2026-01-03T00:00:00Z',
							createdAt: '2026-01-03T00:00:00Z',
							updatedAt: '2026-01-04T00:00:00Z',
							children: [
								{
									id: 'deleted-grandchild',
									parentId: 'deleted-child',
									organizationUnitId: 'unit-id',
									name: 'Deleted Grandchild',
									description: null,
									deletedAt: null,
									createdAt: '2026-01-05T00:00:00Z',
									updatedAt: '2026-01-06T00:00:00Z',
									children: []
								}
							]
						}
					]
				},
				{
					id: 'other-root',
					parentId: null,
					organizationUnitId: 'unit-id',
					name: 'Other Root',
					description: null,
					deletedAt: null,
					createdAt: '2026-01-07T00:00:00Z',
					updatedAt: '2026-01-08T00:00:00Z',
					children: []
				}
			],
			'selected-missing'
		);

		expect(rows.map(({ id, depth }) => [id, depth])).toEqual([
			['active-root', 0],
			['other-root', 0]
		]);
	});

	it('returns folder and document repository rows in a deterministic order', () => {
		const rows = repositoryRows(
			[
				folder('folder-b', 'Statements'),
				folder('folder-a2', 'Invoices'),
				folder('folder-a1', 'Invoices')
			],
			[
				document({ id: 'document-b', displayName: 'zeta.pdf' }),
				document({ id: 'document-a2', displayName: 'Invoice.pdf' }),
				document({ id: 'document-a1', displayName: 'Invoice.pdf' })
			]
		);

		expect(rows.map(({ type, id, name }) => [type, id, name])).toEqual([
			['folder', 'folder-a1', 'Invoices'],
			['folder', 'folder-a2', 'Invoices'],
			['folder', 'folder-b', 'Statements'],
			['document', 'document-a1', 'Invoice.pdf'],
			['document', 'document-a2', 'Invoice.pdf'],
			['document', 'document-b', 'zeta.pdf']
		]);
	});
});

function folder(id: string, name: string): DocumentFolderNode {
	return {
		id,
		parentId: 'root-a',
		organizationUnitId: 'unit-id',
		name,
		description: null,
		deletedAt: null,
		createdAt: '2026-07-01T10:00:00Z',
		updatedAt: '2026-07-01T10:00:00Z',
		children: []
	};
}
