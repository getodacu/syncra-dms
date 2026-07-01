import { afterEach, describe, expect, it, vi } from 'vitest';

import {
	RbacApiError,
	addGroupUser,
	assignGroupRole,
	assignRolePermission,
	assignUserRole,
	checkPermission,
	createGroup,
	createRole,
	createUser,
	deleteRole,
	getMyPermissions,
	isRbacApiError,
	listPermissionCategories,
	listPermissions,
	listUsers,
	removeRolePermission,
	updateUser
} from './rbac';
import type {
	Group,
	GroupRoleAssignment,
	GroupUserAssignment,
	Permission,
	PermissionGrant,
	Role,
	User,
	UserRoleAssignment
} from './rbac';

describe('server RBAC client', () => {
	afterEach(() => {
		vi.unstubAllEnvs();
	});

	it('fetches current permissions with internal and cookie headers', async () => {
		vi.stubEnv('SYNCRA_API_BASE_URL', 'http://api.test');
		vi.stubEnv('SYNCRA_INTERNAL_API_TOKEN', 'internal-token');
		const grant = permissionGrant();
		const fetch = vi.fn(async () => jsonResponse({ permissions: [grant] }));

		const result = await getMyPermissions(fetch, 'auth.session_token=abc');

		expect(result.permissions).toEqual([grant]);
		expect(fetch).toHaveBeenCalledWith(
			'http://api.test/api/me/permissions',
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

	it('checks a permission with the expected JSON payload', async () => {
		vi.stubEnv('SYNCRA_API_BASE_URL', 'http://api.test');
		vi.stubEnv('SYNCRA_INTERNAL_API_TOKEN', 'internal-token');
		const input = { permission: 'document.view', organizationUnitId: 'unit-id' };
		const fetch = vi.fn(async () => jsonResponse({ allowed: true }));

		await expect(checkPermission(fetch, 'auth.session_token=abc', input)).resolves.toEqual({
			allowed: true
		});

		expect(fetch).toHaveBeenCalledWith('http://api.test/api/auth/check-permission', {
			method: 'POST',
			headers: expect.any(Headers),
			body: JSON.stringify(input)
		});
		const firstCall = fetch.mock.calls[0] as unknown as [string, RequestInit];
		const headers = firstCall[1].headers as Headers;
		expect(headers.get('content-type')).toBe('application/json');
	});

	it('throws RbacApiError with public backend messages for failures', async () => {
		vi.stubEnv('SYNCRA_API_BASE_URL', 'http://api.test');
		vi.stubEnv('SYNCRA_INTERNAL_API_TOKEN', 'internal-token');
		const fetch = vi.fn(async () => jsonResponse({ error: 'role.view required' }, 403));

		try {
			await listUsers(fetch, 'auth.session_token=abc');
			throw new Error('Expected listUsers to reject');
		} catch (error) {
			expect(error).toMatchObject(new RbacApiError(403, 'role.view required'));
			expect(isRbacApiError(error)).toBe(true);
		}
	});

	it('maps server errors and invalid success responses to boundary errors', async () => {
		vi.stubEnv('SYNCRA_API_BASE_URL', 'http://api.test');
		vi.stubEnv('SYNCRA_INTERNAL_API_TOKEN', 'internal-token');
		const serverErrorFetch = vi.fn(async () => jsonResponse({ error: 'database exploded' }, 500));

		await expect(getMyPermissions(serverErrorFetch, 'auth.session_token=abc')).rejects.toMatchObject(
			new RbacApiError(502, 'RBAC request failed')
		);

		const invalidFetch = vi.fn(async () => jsonResponse({ permissions: [{ code: 'document.view' }] }));
		await expect(getMyPermissions(invalidFetch, 'auth.session_token=abc')).rejects.toMatchObject(
			new RbacApiError(502, 'Invalid RBAC response')
		);
	});

	it('uses the expected admin API request shapes', async () => {
		vi.stubEnv('SYNCRA_API_BASE_URL', 'http://api.test');
		vi.stubEnv('SYNCRA_INTERNAL_API_TOKEN', 'internal-token');
		const responses: unknown[] = [
			userList([user()]),
			user({ id: 'created-user' }),
			user({ id: 'user/id', name: 'Grace Hopper' }),
			{ permissions: [permission()] },
			{ categories: ['User Management'] },
			userRoleAssignment(),
			role(),
			permission(),
			group(),
			groupUserAssignment(),
			groupRoleAssignment(),
			{ ok: true },
			{ ok: true }
		];
		const fetch = vi.fn(async () => jsonResponse(responses.shift()));

		await expect(listUsers(fetch, null)).resolves.toEqual(userList([user()]));
		await createUser(fetch, null, {
			name: 'Ada Lovelace',
			email: 'ada@example.com',
			status: 'invited',
			primaryOrganizationUnitId: 'unit-id',
			managerUserId: null,
			jobTitle: 'Engineer',
			phone: null
		});
		await updateUser(fetch, null, 'user/id', {
			name: 'Grace Hopper',
			managerUserId: null,
			jobTitle: 'Admiral',
			phone: '555-0100'
		});
		await listPermissions(fetch, null);
		await listPermissionCategories(fetch, null);
		await assignUserRole(fetch, null, 'user/id', {
			roleId: 'role-id',
			scopeType: 'organization_unit',
			organizationUnitId: 'unit/id'
		});
		await createRole(fetch, null, {
			name: 'Unit Manager',
			code: 'unit_manager',
			description: 'Manages a unit',
			isActive: true
		});
		await assignRolePermission(fetch, null, 'role/id', 'permission/id');
		await createGroup(fetch, null, {
			name: 'Legal Reviewers',
			code: 'legal_reviewers',
			description: null,
			organizationUnitId: 'unit-id',
			isActive: true
		});
		await addGroupUser(fetch, null, 'group/id', 'user/id');
		await assignGroupRole(fetch, null, 'group/id', {
			roleId: 'role-id',
			scopeType: 'global',
			organizationUnitId: null
		});
		await removeRolePermission(fetch, null, 'role/id', 'permission/id');
		await deleteRole(fetch, null, 'role/id');

		expect(fetch).toHaveBeenNthCalledWith(
			1,
			'http://api.test/api/users',
			expect.objectContaining({ method: 'GET' })
		);
		expect(fetch).toHaveBeenNthCalledWith(2, 'http://api.test/api/users', {
			method: 'POST',
			headers: expect.any(Headers),
			body: JSON.stringify({
				name: 'Ada Lovelace',
				email: 'ada@example.com',
				status: 'invited',
				primaryOrganizationUnitId: 'unit-id',
				managerUserId: null,
				jobTitle: 'Engineer',
				phone: null
			})
		});
		expect(fetch).toHaveBeenNthCalledWith(3, 'http://api.test/api/users/user%2Fid', {
			method: 'PATCH',
			headers: expect.any(Headers),
			body: JSON.stringify({
				name: 'Grace Hopper',
				managerUserId: null,
				jobTitle: 'Admiral',
				phone: '555-0100'
			})
		});
		expect(fetch).toHaveBeenNthCalledWith(
			4,
			'http://api.test/api/permissions',
			expect.objectContaining({ method: 'GET' })
		);
		expect(fetch).toHaveBeenNthCalledWith(
			5,
			'http://api.test/api/permissions/categories',
			expect.objectContaining({ method: 'GET' })
		);
		expect(fetch).toHaveBeenNthCalledWith(6, 'http://api.test/api/users/user%2Fid/roles', {
			method: 'POST',
			headers: expect.any(Headers),
			body: JSON.stringify({
				roleId: 'role-id',
				scopeType: 'organization_unit',
				organizationUnitId: 'unit/id'
			})
		});
		expect(fetch).toHaveBeenNthCalledWith(7, 'http://api.test/api/roles', {
			method: 'POST',
			headers: expect.any(Headers),
			body: JSON.stringify({
				name: 'Unit Manager',
				code: 'unit_manager',
				description: 'Manages a unit',
				isActive: true
			})
		});
		expect(fetch).toHaveBeenNthCalledWith(8, 'http://api.test/api/roles/role%2Fid/permissions', {
			method: 'POST',
			headers: expect.any(Headers),
			body: JSON.stringify({ permissionId: 'permission/id' })
		});
		expect(fetch).toHaveBeenNthCalledWith(9, 'http://api.test/api/groups', {
			method: 'POST',
			headers: expect.any(Headers),
			body: JSON.stringify({
				name: 'Legal Reviewers',
				code: 'legal_reviewers',
				description: null,
				organizationUnitId: 'unit-id',
				isActive: true
			})
		});
		expect(fetch).toHaveBeenNthCalledWith(10, 'http://api.test/api/groups/group%2Fid/users', {
			method: 'POST',
			headers: expect.any(Headers),
			body: JSON.stringify({ userId: 'user/id' })
		});
		expect(fetch).toHaveBeenNthCalledWith(11, 'http://api.test/api/groups/group%2Fid/roles', {
			method: 'POST',
			headers: expect.any(Headers),
			body: JSON.stringify({
				roleId: 'role-id',
				scopeType: 'global',
				organizationUnitId: null
			})
		});
		expect(fetch).toHaveBeenNthCalledWith(
			12,
			'http://api.test/api/roles/role%2Fid/permissions/permission%2Fid',
			expect.objectContaining({ method: 'DELETE' })
		);
		expect(fetch).toHaveBeenNthCalledWith(
			13,
			'http://api.test/api/roles/role%2Fid',
			expect.objectContaining({ method: 'DELETE' })
		);
	});
});

function user(overrides: Partial<User> = {}): User {
	return {
		id: 'user-id',
		name: 'Ada Lovelace',
		email: 'ada@example.com',
		emailVerified: true,
		image: null,
		preferredLanguage: 'en',
		role: 'admin',
		status: 'active',
		primaryOrganizationUnitId: null,
		managerUserId: null,
		jobTitle: null,
		phone: null,
		lastLoginAt: null,
		createdAt: '2026-06-30T10:00:00Z',
		updatedAt: '2026-06-30T10:00:00Z',
		...overrides
	};
}

function userList(users: User[]) {
	return { users };
}

function permissionGrant(overrides: Partial<PermissionGrant> = {}): PermissionGrant {
	return {
		code: 'document.view',
		scopeType: 'organization_unit',
		organizationUnitId: 'unit-id',
		source: 'user_role',
		...overrides
	};
}

function permission(overrides: Partial<Permission> = {}): Permission {
	return {
		id: 'permission-id',
		code: 'user.view',
		name: 'View users',
		description: null,
		category: 'User Management',
		isSystem: true,
		createdAt: '2026-06-30T10:00:00Z',
		updatedAt: '2026-06-30T10:00:00Z',
		...overrides
	};
}

function role(overrides: Partial<Role> = {}): Role {
	return {
		id: 'role-id',
		name: 'Unit Manager',
		code: 'unit_manager',
		description: null,
		isSystem: false,
		isActive: true,
		createdAt: '2026-06-30T10:00:00Z',
		updatedAt: '2026-06-30T10:00:00Z',
		...overrides
	};
}

function group(overrides: Partial<Group> = {}): Group {
	return {
		id: 'group-id',
		name: 'Legal Reviewers',
		code: 'legal_reviewers',
		description: null,
		organizationUnitId: null,
		isActive: true,
		createdAt: '2026-06-30T10:00:00Z',
		updatedAt: '2026-06-30T10:00:00Z',
		...overrides
	};
}

function userRoleAssignment(overrides: Partial<UserRoleAssignment> = {}): UserRoleAssignment {
	return {
		id: 'assignment-id',
		userId: 'user-id',
		roleId: 'role-id',
		scopeType: 'global',
		organizationUnitId: null,
		createdAt: '2026-06-30T10:00:00Z',
		updatedAt: '2026-06-30T10:00:00Z',
		...overrides
	};
}

function groupRoleAssignment(overrides: Partial<GroupRoleAssignment> = {}): GroupRoleAssignment {
	return {
		id: 'group-role-id',
		groupId: 'group-id',
		roleId: 'role-id',
		scopeType: 'global',
		organizationUnitId: null,
		createdAt: '2026-06-30T10:00:00Z',
		...overrides
	};
}

function groupUserAssignment(overrides: Partial<GroupUserAssignment> = {}): GroupUserAssignment {
	return {
		groupId: 'group-id',
		userId: 'user-id',
		createdAt: '2026-06-30T10:00:00Z',
		...overrides
	};
}

function jsonResponse(body: unknown, status = 200) {
	return new Response(JSON.stringify(body), {
		status,
		headers: { 'content-type': 'application/json' }
	});
}
