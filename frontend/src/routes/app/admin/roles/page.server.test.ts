import { readFileSync } from 'node:fs';
import { describe, expect, it } from 'vitest';

import { load } from './+page.server';

describe('admin roles page server load', () => {
	it('exposes role management flags from permissions', () => {
		const result = load(
			loadEvent(['role.view', 'role.create', 'role.update', 'role.delete', 'role.assign_permissions']) as never
		);

		expect(result).toEqual({
			canViewRoles: true,
			canManageRoles: true,
			canCreateRoles: true,
			canUpdateRoles: true,
			canDeleteRoles: true,
			canAssignRolePermissions: true,
			selectedRoleId: null
		});
	});

	it('requires role.view for read access and ignores auth role alone', () => {
		const result = load(loadEvent([], { role: 'admin' }) as never) as Record<string, unknown>;

		expect(result.canViewRoles).toBe(false);
		expect(result.canManageRoles).toBe(false);
	});

	it('allows system administrators to manage roles', () => {
		const result = load(loadEvent(['system.admin']) as never) as Record<string, unknown>;

		expect(result.canViewRoles).toBe(true);
		expect(result.canAssignRolePermissions).toBe(true);
	});
});

describe('admin roles page source', () => {
	it('uses TanStack query and the permission matrix against browser APIs', () => {
		const source = readFileSync(new URL('./+page.svelte', import.meta.url), 'utf8');

		expect(source).toContain("import { createMutation, createQuery, useQueryClient }");
		expect(source).toContain('const rolesQuery = createQuery<RoleListResponse, Error>');
		expect(source).toContain('const permissionsQuery = createQuery<PermissionListResponse, Error>');
		expect(source).toContain('const createMutationState = createMutation<Role, Error, CreateRoleInput>');
		expect(source).toContain('queryClient.invalidateQueries({ queryKey: ROLES_QUERY_KEY })');
		expect(source).toContain('<PermissionMatrix');
		expect(source).not.toContain('$lib/server');
	});
});

function loadEvent(permissions: string[], user: { role: string } | null = { role: 'user' }) {
	return {
		locals: { user, permissions },
		url: new URL('http://localhost/app/admin/roles')
	};
}
