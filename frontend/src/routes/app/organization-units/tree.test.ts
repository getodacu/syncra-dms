import { describe, expect, it } from 'vitest';

import {
	collectMoveTargets,
	flattenUnitTree,
	countUnits,
	findUnit,
	selectInitialUnit,
	type OrganizationUnitNode
} from './tree';

const units: OrganizationUnitNode[] = [
	{
		id: 'root-a',
		name: 'Root A',
		code: 'A',
		description: 'First root',
		createdAt: '2026-01-01T00:00:00Z',
		updatedAt: '2026-01-02T00:00:00Z',
		children: [
			{
				id: 'child-a1',
				parentId: 'root-a',
				name: 'Child A1',
				createdAt: '2026-01-03T00:00:00Z',
				updatedAt: '2026-01-04T00:00:00Z',
				children: [
					{
						id: 'grandchild-a1a',
						parentId: 'child-a1',
						name: 'Grandchild A1A',
						createdAt: '2026-01-05T00:00:00Z',
						updatedAt: '2026-01-06T00:00:00Z',
						children: []
					}
				]
			},
			{
				id: 'child-a2',
				parentId: 'root-a',
				name: 'Child A2',
				createdAt: '2026-01-07T00:00:00Z',
				updatedAt: '2026-01-08T00:00:00Z',
				children: []
			}
		]
	},
	{
		id: 'root-b',
		name: 'Root B',
		code: 'B',
		createdAt: '2026-01-09T00:00:00Z',
		updatedAt: '2026-01-10T00:00:00Z',
		children: [
			{
				id: 'child-b1',
				parentId: 'root-b',
				name: 'Child B1',
				createdAt: '2026-01-11T00:00:00Z',
				updatedAt: '2026-01-12T00:00:00Z',
				children: []
			}
		]
	}
];

describe('organization unit tree utilities', () => {
	it('flattens units in pre-order with depth values', () => {
		const rows = flattenUnitTree(units);

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
			code: 'A',
			description: 'First root',
			depth: 0
		});
		expect('children' in rows[0]).toBe(false);
	});

	it('selects the requested unit before falling back to the first root', () => {
		expect(selectInitialUnit(units, 'grandchild-a1a')?.id).toBe('grandchild-a1a');
		expect(selectInitialUnit(units, 'missing')?.id).toBe('root-a');
		expect(selectInitialUnit(units)?.id).toBe('root-a');
		expect(selectInitialUnit([], 'root-a')).toBeNull();
	});

	it('finds nested units and returns null for nil or unknown ids', () => {
		expect(findUnit(units, 'child-b1')?.name).toBe('Child B1');
		expect(findUnit(units, null)).toBeNull();
		expect(findUnit(units, undefined)).toBeNull();
		expect(findUnit(units, 'missing')).toBeNull();
	});

	it('counts every root and nested unit', () => {
		expect(countUnits(units)).toBe(6);
		expect(countUnits([])).toBe(0);
	});

	it('collects move targets without the selected node or its descendants', () => {
		const rows = collectMoveTargets(units, 'child-a1');

		expect(rows.map(({ id, depth }) => [id, depth])).toEqual([
			['root-a', 0],
			['child-a2', 1],
			['root-b', 0],
			['child-b1', 1]
		]);
	});

	it('collects move targets without archived branches', () => {
		const rows = collectMoveTargets(
			[
				{
					id: 'active-root',
					name: 'Active Root',
					createdAt: '2026-01-01T00:00:00Z',
					updatedAt: '2026-01-02T00:00:00Z',
					children: [
						{
							id: 'archived-child',
							parentId: 'active-root',
							name: 'Archived Child',
							isArchived: true,
							createdAt: '2026-01-03T00:00:00Z',
							updatedAt: '2026-01-04T00:00:00Z',
							children: [
								{
									id: 'archived-grandchild',
									parentId: 'archived-child',
									name: 'Archived Grandchild',
									createdAt: '2026-01-05T00:00:00Z',
									updatedAt: '2026-01-06T00:00:00Z',
									children: []
								}
							]
						}
					]
				},
				{
					id: 'archived-root',
					name: 'Archived Root',
					isArchived: true,
					createdAt: '2026-01-07T00:00:00Z',
					updatedAt: '2026-01-08T00:00:00Z',
					children: []
				},
				{
					id: 'other-root',
					name: 'Other Root',
					createdAt: '2026-01-09T00:00:00Z',
					updatedAt: '2026-01-10T00:00:00Z',
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
});
