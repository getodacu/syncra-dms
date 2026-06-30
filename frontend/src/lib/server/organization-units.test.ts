import { afterEach, describe, expect, it, vi } from 'vitest';

import {
	OrganizationUnitApiError,
	archiveOrganizationUnit,
	createOrganizationUnit,
	getOrganizationUnitTree,
	isOrganizationUnitApiError,
	moveOrganizationUnit,
	updateOrganizationUnit
} from './organization-units';
import type { OrganizationUnit } from './organization-units';

describe('server organization unit client', () => {
	afterEach(() => {
		vi.unstubAllEnvs();
	});

	it('sends the internal token and forwarded cookie when reading the tree', async () => {
		vi.stubEnv('SYNCRA_API_BASE_URL', 'http://api.test');
		vi.stubEnv('SYNCRA_INTERNAL_API_TOKEN', 'internal-token');
		const unit = organizationUnit();
		const fetch = vi.fn(async () => jsonResponse({ units: [unit] }));

		const result = await getOrganizationUnitTree(fetch, 'auth.session_token=abc');

		expect(result.units).toEqual([unit]);
		expect(fetch).toHaveBeenCalledWith(
			'http://api.test/api/organization-units/tree',
			expect.objectContaining({
				method: 'GET',
				headers: expect.any(Headers)
			})
		);
		const firstCall = fetch.mock.calls[0] as unknown as [string, RequestInit];
		const headers = firstCall[1].headers as Headers;
		expect(headers.get('X-Syncra-Internal-Token')).toBe('internal-token');
		expect(headers.get('cookie')).toBe('auth.session_token=abc');
	});

	it('posts the expected JSON body when creating an organization unit', async () => {
		vi.stubEnv('SYNCRA_API_BASE_URL', 'http://api.test');
		vi.stubEnv('SYNCRA_INTERNAL_API_TOKEN', 'internal-token');
		const unit = organizationUnit({ id: 'new-unit', children: undefined });
		const fetch = vi.fn(async () => jsonResponse(unit));
		const input = {
			parentId: 'parent-id',
			name: 'Finance',
			code: 'FIN',
			description: 'Finance department'
		};

		const result = await createOrganizationUnit(fetch, 'auth.session_token=abc', input);

		expect(fetch).toHaveBeenCalledWith('http://api.test/api/organization-units', {
			method: 'POST',
			headers: expect.any(Headers),
			body: JSON.stringify(input)
		});
		expect(result.children).toBeUndefined();
		const firstCall = fetch.mock.calls[0] as unknown as [string, RequestInit];
		const headers = firstCall[1].headers as Headers;
		expect(headers.get('content-type')).toBe('application/json');
	});

	it('throws OrganizationUnitApiError with the public backend message for failures', async () => {
		vi.stubEnv('SYNCRA_API_BASE_URL', 'http://api.test');
		vi.stubEnv('SYNCRA_INTERNAL_API_TOKEN', 'internal-token');
		const fetch = vi.fn(async () => jsonResponse({ error: 'admin role required' }, 403));

		try {
			await createOrganizationUnit(fetch, 'auth.session_token=abc', {
				parentId: null,
				name: 'Finance',
				code: 'FIN',
				description: 'Finance department'
			});
			throw new Error('Expected createOrganizationUnit to reject');
		} catch (error) {
			expect(error).toMatchObject(new OrganizationUnitApiError(403, 'admin role required'));
			expect(isOrganizationUnitApiError(error)).toBe(true);
		}
	});

	it('maps server errors through the public error contract', async () => {
		vi.stubEnv('SYNCRA_API_BASE_URL', 'http://api.test');
		vi.stubEnv('SYNCRA_INTERNAL_API_TOKEN', 'internal-token');
		const fetch = vi.fn(async () => jsonResponse({ error: 'database exploded' }, 500));

		await expect(getOrganizationUnitTree(fetch, 'auth.session_token=abc')).rejects.toMatchObject(
			new OrganizationUnitApiError(502, 'Organization Unit request failed')
		);
	});

	it('throws clear boundary errors for missing config, network failures, and invalid success responses', async () => {
		vi.stubEnv('SYNCRA_API_BASE_URL', 'http://api.test');
		vi.stubEnv('SYNCRA_INTERNAL_API_TOKEN', ' ');
		await expect(getOrganizationUnitTree(vi.fn(), 'auth.session_token=abc')).rejects.toMatchObject(
			new OrganizationUnitApiError(500, 'Organization Unit service is not configured')
		);

		vi.stubEnv('SYNCRA_INTERNAL_API_TOKEN', 'internal-token');
		const failingFetch = vi.fn(async () => {
			throw new Error('network down');
		});
		await expect(getOrganizationUnitTree(failingFetch, 'auth.session_token=abc')).rejects.toMatchObject(
			new OrganizationUnitApiError(503, 'Organization Unit service unavailable')
		);

		const invalidFetch = vi.fn(async () => new Response('<html></html>', { status: 200 }));
		await expect(getOrganizationUnitTree(invalidFetch, 'auth.session_token=abc')).rejects.toMatchObject(
			new OrganizationUnitApiError(502, 'Invalid Organization Unit response')
		);
	});

	it('uses encoded ids and JSON payloads for update, move, and archive requests', async () => {
		vi.stubEnv('SYNCRA_API_BASE_URL', 'http://api.test');
		vi.stubEnv('SYNCRA_INTERNAL_API_TOKEN', 'internal-token');
		const fetch = vi.fn(async (...[, init]: Parameters<typeof globalThis.fetch>) => {
			if (init?.method === 'POST') return jsonResponse({ ok: true });
			return jsonResponse(organizationUnit());
		});
		const input = {
			parentId: null,
			name: 'Legal',
			code: 'LEGAL',
			description: 'Legal team'
		};

		await updateOrganizationUnit(fetch, null, 'unit/id', input);
		await moveOrganizationUnit(fetch, null, 'unit/id', 'parent/id');
		const archive = await archiveOrganizationUnit(fetch, null, 'unit/id');

		expect(fetch).toHaveBeenNthCalledWith(1, 'http://api.test/api/organization-units/unit%2Fid', {
			method: 'PATCH',
			headers: expect.any(Headers),
			body: JSON.stringify(input)
		});
		expect(fetch).toHaveBeenNthCalledWith(
			2,
			'http://api.test/api/organization-units/unit%2Fid/parent',
			{
				method: 'PATCH',
				headers: expect.any(Headers),
				body: JSON.stringify({ parentId: 'parent/id' })
			}
		);
		expect(fetch).toHaveBeenNthCalledWith(
			3,
			'http://api.test/api/organization-units/unit%2Fid/archive',
			{
				method: 'POST',
				headers: expect.any(Headers),
				body: JSON.stringify({})
			}
		);
		expect(archive).toEqual({ ok: true });
	});
});

function organizationUnit(overrides: Partial<OrganizationUnit> = {}): OrganizationUnit {
	return {
		id: 'unit-id',
		parentId: null,
		name: 'Operations',
		code: 'OPS',
		description: 'Operations team',
		archivedAt: null,
		createdAt: '2026-06-30T10:00:00Z',
		updatedAt: '2026-06-30T10:00:00Z',
		children: [],
		...overrides
	};
}

function jsonResponse(body: unknown, status = 200) {
	return new Response(JSON.stringify(body), {
		status,
		headers: { 'content-type': 'application/json' }
	});
}
