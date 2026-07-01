import { readFileSync } from 'node:fs';
import { describe, expect, it } from 'vitest';

import { load } from './+page.server';

describe('admin permissions page server load', () => {
	it('is read-only and requires role.view', () => {
		expect(load(loadEvent(['role.view']) as never)).toEqual({
			canViewPermissions: true,
			canManagePermissions: false
		});
		expect(load(loadEvent([]) as never)).toEqual({
			canViewPermissions: false,
			canManagePermissions: false
		});
	});

	it('allows system administrators to view permissions without mutation flags', () => {
		expect(load(loadEvent(['system.admin']) as never)).toEqual({
			canViewPermissions: true,
			canManagePermissions: false
		});
	});
});

describe('admin permissions page source', () => {
	it('renders grouped permissions through browser APIs only', () => {
		const source = readFileSync(new URL('./+page.svelte', import.meta.url), 'utf8');

		expect(source).toContain("import { createQuery }");
		expect(source).toContain('fetchPermissions(fetch)');
		expect(source).toContain('groupPermissionsByCategory');
		expect(source).not.toContain('createMutation');
		expect(source).not.toContain('$lib/server');
	});
});

function loadEvent(permissions: string[]) {
	return {
		locals: { user: { role: 'user' }, permissions },
		url: new URL('http://localhost/app/admin/permissions')
	};
}
