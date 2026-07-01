import { readFileSync } from 'node:fs';
import { describe, expect, it } from 'vitest';

import { load } from './+page.server';

describe('admin groups page server load', () => {
	it('exposes group management permission flags', () => {
		const result = load(
			loadEvent(['group.view', 'group.create', 'group.update', 'group.delete', 'group.manage_users', 'group.assign_roles']) as never
		);

		expect(result).toEqual({
			canViewGroups: true,
			canManageGroups: true,
			canCreateGroups: true,
			canUpdateGroups: true,
			canDeleteGroups: true,
			canManageGroupUsers: true,
			canAssignGroupRoles: true,
			selectedGroupId: null
		});
	});

	it('allows system administrators and ignores auth role alone', () => {
		const systemResult = load(loadEvent(['system.admin']) as never) as Record<string, unknown>;
		const roleOnlyResult = load(loadEvent([], { role: 'admin' }) as never) as Record<
			string,
			unknown
		>;

		expect(systemResult.canViewGroups).toBe(true);
		expect(systemResult.canAssignGroupRoles).toBe(true);
		expect(roleOnlyResult.canViewGroups).toBe(false);
		expect(roleOnlyResult.canManageGroups).toBe(false);
	});
});

describe('admin groups page source', () => {
	it('uses TanStack query and mutations against browser APIs', () => {
		const source = readFileSync(new URL('./+page.svelte', import.meta.url), 'utf8');

		expect(source).toContain("import { createMutation, createQuery, useQueryClient }");
		expect(source).toContain('const groupsQuery = createQuery<GroupListResponse, Error>');
		expect(source).toContain('const createMutationState = createMutation<Group, Error, CreateGroupInput>');
		expect(source).toContain('const memberMutationState = createMutation<');
		expect(source).toContain('const roleMutationState = createMutation<');
		expect(source).toContain('<GroupEditor');
		expect(source).toContain('<GroupMembers');
		expect(source).toContain('<GroupRoleAssignments');
		expect(source).not.toContain('$lib/server');
	});
});

function loadEvent(permissions: string[], user: { role: string } | null = { role: 'user' }) {
	return {
		locals: { user, permissions },
		url: new URL('http://localhost/app/admin/groups')
	};
}
