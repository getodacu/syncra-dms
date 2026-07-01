import { readFileSync } from 'node:fs';
import { describe, expect, it } from 'vitest';

import { load as layoutLoad } from './app/+layout.server';

const privateUser = {
	id: 'user-id',
	name: 'Ada Lovelace',
	email: 'ada@example.com',
	emailVerified: true,
	image: 'https://example.com/avatar.png',
	preferredLanguage: 'en',
	role: 'admin',
	status: 'active',
	lastLoginAt: '2026-06-29T00:00:00Z',
	createdAt: '2026-06-01T00:00:00Z',
	updatedAt: '2026-06-30T00:00:00Z'
};

describe('app layout permission data', () => {
	it('returns public permissions without exposing the session token', () => {
		const result = layoutLoad({
			locals: {
				user: privateUser,
				permissions: ['user.view', 'role.view'],
				session: {
					id: 'session-id',
					token: 'secret-token',
					userId: 'user-id',
					expiresAt: '2026-06-30T00:00:00Z',
					createdAt: '2026-06-01T00:00:00Z',
					updatedAt: '2026-06-30T00:00:00Z'
				}
			}
		} as never);

		expect(result).toEqual({
			user: {
				id: 'user-id',
				name: 'Ada Lovelace',
				email: 'ada@example.com',
				image: 'https://example.com/avatar.png',
				role: 'admin'
			},
			permissions: ['user.view', 'role.view'],
			session: { expiresAt: '2026-06-30T00:00:00Z' }
		});
		expect(JSON.stringify(result)).not.toContain('secret-token');
	});
});

describe('app sidebar permission gates', () => {
	it('passes layout permissions into the sidebar shell', () => {
		const source = readFileSync(new URL('./app/+layout.svelte', import.meta.url), 'utf8');

		expect(source).toContain(
			'<AppSidebar variant="inset" user={data.user} permissions={data.permissions}'
		);
	});

	it('gates admin navigation by RBAC permissions instead of auth role', () => {
		const source = readFileSync(
			new URL('../lib/components/app-sidebar.svelte', import.meta.url),
			'utf8'
		);

		expect(source).toContain('permissions.includes');
		expect(source).toContain("'system.admin'");
		expect(source).toContain("'user.view'");
		expect(source).toContain("'role.view'");
		expect(source).toContain("'group.view'");
		expect(source).not.toContain("user?.role === 'admin'");
	});
});
