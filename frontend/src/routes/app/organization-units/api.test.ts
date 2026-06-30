import { readFileSync } from 'node:fs';
import { describe, expect, it, vi } from 'vitest';

import {
	archiveOrganizationUnit,
	createOrganizationUnit,
	fetchOrganizationUnitTree,
	isOrganizationUnitListResponse,
	moveOrganizationUnit,
	updateOrganizationUnit
} from './api';
import type { OrganizationUnitListResponse } from './api';

function jsonResponse(body: unknown, init?: ResponseInit) {
	return new Response(JSON.stringify(body), {
		headers: { 'content-type': 'application/json' },
		...init
	});
}

function listResponse(): OrganizationUnitListResponse {
	return {
		units: [
			{
				id: 'root',
				name: 'Root',
				code: 'ROOT',
				description: 'Root unit',
				createdAt: '2026-06-30T00:00:00Z',
				updatedAt: '2026-06-30T00:00:00Z',
				children: []
			}
		]
	};
}

describe('organization unit browser API', () => {
	it('fetches the organization unit tree through the Svelte API wrapper', async () => {
		const fetchMock = vi.fn().mockResolvedValue(jsonResponse(listResponse(), { status: 200 }));

		const result = await fetchOrganizationUnitTree(fetchMock);

		expect(fetchMock).toHaveBeenCalledWith('/api/organization-units/tree', {
			method: 'GET',
			headers: undefined,
			body: undefined
		});
		expect(result.units[0].name).toBe('Root');
	});

	it('posts normalized create payloads', async () => {
		const fetchMock = vi
			.fn()
			.mockResolvedValue(jsonResponse({ ...listResponse().units[0], id: 'created-id' }, { status: 201 }));

		await createOrganizationUnit(fetchMock, {
			parentId: '   ',
			name: 'Operations',
			code: 'OPS',
			description: null
		});

		expect(fetchMock).toHaveBeenCalledWith('/api/organization-units', {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({
				parentId: null,
				name: 'Operations',
				code: 'OPS',
				description: ''
			})
		});
	});

	it('patches update and move requests with encoded ids', async () => {
		const unit = listResponse().units[0];
		const fetchMock = vi
			.fn()
			.mockResolvedValueOnce(jsonResponse(unit, { status: 200 }))
			.mockResolvedValueOnce(jsonResponse(unit, { status: 200 }));

		await updateOrganizationUnit(fetchMock, {
			id: 'unit/id',
			input: { parentId: null, name: 'Root', code: '', description: '' }
		});
		await moveOrganizationUnit(fetchMock, { id: 'unit/id', parentId: null });

		expect(fetchMock).toHaveBeenNthCalledWith(1, '/api/organization-units/unit%2Fid', {
			method: 'PATCH',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({
				parentId: null,
				name: 'Root',
				code: '',
				description: ''
			})
		});
		expect(fetchMock).toHaveBeenNthCalledWith(2, '/api/organization-units/unit%2Fid/parent', {
			method: 'PATCH',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({ parentId: null })
		});
	});

	it('posts archive requests with encoded ids', async () => {
		const fetchMock = vi.fn().mockResolvedValue(jsonResponse({ ok: true }, { status: 200 }));

		const result = await archiveOrganizationUnit(fetchMock, { id: 'unit/id' });

		expect(fetchMock).toHaveBeenCalledWith('/api/organization-units/unit%2Fid/archive', {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({})
		});
		expect(result).toEqual({ ok: true });
	});

	it('uses public-safe messages for failed requests', async () => {
		const fetchMock = vi
			.fn()
			.mockResolvedValue(jsonResponse({ error: 'database unavailable' }, { status: 503 }));

		await expect(fetchOrganizationUnitTree(fetchMock)).rejects.toThrow(
			'Failed to load organization units'
		);
	});

	it('rejects malformed successful payloads', async () => {
		const fetchMock = vi.fn().mockResolvedValue(jsonResponse({ ok: true }, { status: 200 }));

		await expect(fetchOrganizationUnitTree(fetchMock)).rejects.toThrow(
			'Invalid organization unit response'
		);
	});

	it('identifies valid organization unit list responses', () => {
		expect(isOrganizationUnitListResponse({ ok: true })).toBe(false);
		expect(isOrganizationUnitListResponse(listResponse())).toBe(true);
	});

	it('exports only the intended browser API surface', () => {
		const source = readFileSync(new URL('./api.ts', import.meta.url), 'utf8');
		const exportedNames = [
			...source.matchAll(/export\s+(?:async\s+)?(?:function|type|const|class)\s+(\w+)/g)
		]
			.map((match) => match[1])
			.sort();

		expect(exportedNames).toEqual([
			'ArchiveOrganizationUnitResponse',
			'ArchiveOrganizationUnitVariables',
			'MoveOrganizationUnitVariables',
			'ORGANIZATION_UNITS_QUERY_KEY',
			'OrganizationUnitInput',
			'OrganizationUnitListResponse',
			'UpdateOrganizationUnitVariables',
			'archiveOrganizationUnit',
			'createOrganizationUnit',
			'fetchOrganizationUnitTree',
			'isOrganizationUnitListResponse',
			'moveOrganizationUnit',
			'updateOrganizationUnit'
		]);
	});
});
