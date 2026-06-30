import { readFileSync } from 'node:fs';
import { describe, expect, it } from 'vitest';

import { load as layoutLoad } from '../+layout.server';
import { load as appPageLoad } from '../+page.server';
import { load } from './+page.server';

const privateUser = {
	id: 'user-id',
	name: 'Ada Lovelace',
	email: 'ada@example.com',
	emailVerified: true,
	image: 'https://example.com/avatar.png',
	preferredLanguage: 'en',
	role: 'admin',
	lastLoginAt: '2026-06-29T00:00:00Z',
	createdAt: '2026-06-01T00:00:00Z',
	updatedAt: '2026-06-30T00:00:00Z'
};
const layoutPublicUser = {
	id: 'user-id',
	name: 'Ada Lovelace',
	email: 'ada@example.com',
	image: 'https://example.com/avatar.png',
	role: 'admin'
};
const pagePublicUser = {
	name: 'Ada Lovelace',
	email: 'ada@example.com',
	role: 'admin'
};

describe('app layout server load', () => {
	it('returns public user and session shapes from locals', () => {
		const locals = {
			user: privateUser,
			session: { id: 'session-id', token: 'secret-token', expiresAt: '2026-06-30T00:00:00Z' }
		};

		expect(layoutLoad({ locals } as never)).toEqual({
			user: layoutPublicUser,
			session: { expiresAt: '2026-06-30T00:00:00Z' }
		});
	});
});

describe('app dashboard server load', () => {
	it('returns a public session shape without the session token', () => {
		const locals = {
			user: privateUser,
			session: { id: 'session-id', token: 'secret-token', expiresAt: '2026-06-30T00:00:00Z' }
		};

		expect(appPageLoad({ locals } as never)).toEqual({
			user: pagePublicUser,
			session: { expiresAt: '2026-06-30T00:00:00Z' }
		});
	});
});

describe('organization units page server load', () => {
	it('returns admin manage access and selectedId without loading the tree', () => {
		const result = load(loadEvent({ role: 'admin' }, 'selected-unit') as never);

		expect(result).toEqual({
			canManageOrganizationUnits: true,
			selectedId: 'selected-unit'
		});
	});

	it('returns read-only access for regular users', () => {
		const result = load(loadEvent({ role: 'user' }) as never);

		expect(result).toEqual({
			canManageOrganizationUnits: false,
			selectedId: null
		});
	});
});

describe('organization units page source', () => {
	it('uses TanStack query and mutations against the Svelte API wrapper', () => {
		const source = readFileSync(new URL('./+page.svelte', import.meta.url), 'utf8');

		expect(source).toContain("import { createMutation, createQuery, useQueryClient }");
		expect(source).toContain('const organizationUnitsQuery = createQuery<OrganizationUnitListResponse, Error>');
		expect(source).toContain('const createMutationState = createMutation<OrganizationUnitNode, Error, OrganizationUnitInput>');
		expect(source).toContain('const updateMutationState = createMutation<');
		expect(source).toContain('const moveMutationState = createMutation<');
		expect(source).toContain('const archiveMutationState = createMutation<');
		expect(source).toContain('queryClient.invalidateQueries({ queryKey: ORGANIZATION_UNITS_QUERY_KEY })');
		expect(source).not.toContain('$lib/server/organization-units');
		expect(source).not.toContain('method="POST"');
	});
});

function loadEvent(user: { role: string } | null, selectedId?: string) {
	const url = new URL('http://localhost/app/organization-units');
	if (selectedId) url.searchParams.set('selectedId', selectedId);
	return {
		locals: { user },
		url
	};
}
