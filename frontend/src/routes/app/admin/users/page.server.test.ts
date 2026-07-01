import { readFileSync } from 'node:fs';
import { describe, expect, it } from 'vitest';

import { load } from './+page.server';

describe('admin users page server load', () => {
	it('exposes user management permission flags and selected user id', () => {
		const result = load(
			loadEvent([
				'user.view',
				'user.create',
				'user.update',
				'user.delete',
				'user.activate',
				'user.suspend',
				'user.assign_role',
				'user.assign_group',
				'user.assign_unit'
			], 'user-id') as never
		);

		expect(result).toEqual({
			canViewUsers: true,
			canManageUsers: true,
			canCreateUsers: true,
			canUpdateUsers: true,
			canDeleteUsers: true,
			canActivateUsers: true,
			canSuspendUsers: true,
			canAssignUserRoles: true,
			canAssignUserGroups: true,
			canAssignUserUnits: true,
			selectedUserId: 'user-id'
		});
	});

	it('allows system administrators to manage users', () => {
		const result = load(loadEvent(['system.admin']) as never) as Record<string, unknown>;

		expect(result.canViewUsers).toBe(true);
		expect(result.canManageUsers).toBe(true);
		expect(result.canAssignUserRoles).toBe(true);
	});

	it('does not infer access from auth role alone', () => {
		const result = load(loadEvent([], undefined, { role: 'admin' }) as never) as Record<
			string,
			unknown
		>;

		expect(result.canViewUsers).toBe(false);
		expect(result.canManageUsers).toBe(false);
	});
});

describe('admin users page source', () => {
	it('uses TanStack query and mutations against the Svelte API wrapper', () => {
		const source = readFileSync(new URL('./+page.svelte', import.meta.url), 'utf8');

		expect(source).toContain("import { createMutation, createQuery, useQueryClient }");
		expect(source).toContain('const usersQuery = createQuery<UserListResponse, Error>');
		expect(source).toContain('const createMutationState = createMutation<User, Error, CreateUserInput>');
		expect(source).toContain('const updateMutationState = createMutation<');
		expect(source).toContain('const statusMutationState = createMutation<');
		expect(source).toContain('queryClient.invalidateQueries({ queryKey: USERS_QUERY_KEY })');
		expect(source).toContain('<UserEditor');
		expect(source).toContain('<UserRoleAssignments');
		expect(source).not.toContain('$lib/server');
	});
});

function loadEvent(
	permissions: string[],
	selectedUserId?: string,
	user: { role: string } | null = { role: 'user' }
) {
	const url = new URL('http://localhost/app/admin/users');
	if (selectedUserId) url.searchParams.set('selectedUserId', selectedUserId);
	return {
		locals: { user, permissions },
		url
	};
}
