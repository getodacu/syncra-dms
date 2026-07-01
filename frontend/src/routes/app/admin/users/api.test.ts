import { readFileSync } from 'node:fs';
import { describe, expect, it, vi } from 'vitest';

import {
	USERS_QUERY_KEY,
	activateUser,
	assignUserGroup,
	assignUserRole,
	createUser,
	deactivateUser,
	fetchUsers,
	removeUserGroup,
	removeUserRole,
	setPrimaryOrganizationUnit,
	softDeleteUser,
	suspendUser,
	updateUser
} from './api';
import type { UserListResponse } from './api';

function jsonResponse(body: unknown, init?: ResponseInit) {
	return new Response(JSON.stringify(body), {
		headers: { 'content-type': 'application/json' },
		...init
	});
}

function usersResponse(): UserListResponse {
	return {
		users: [
			{
				id: 'user-id',
				name: 'Ada Lovelace',
				email: 'ada@example.com',
				emailVerified: true,
				image: null,
				preferredLanguage: 'en',
				role: 'user',
				status: 'active',
				primaryOrganizationUnitId: null,
				managerUserId: null,
				jobTitle: 'Engineer',
				phone: null,
				lastLoginAt: null,
				createdAt: '2026-06-30T00:00:00Z',
				updatedAt: '2026-06-30T00:00:00Z'
			}
		]
	};
}

describe('admin users browser API', () => {
	it('fetches users through the Svelte API wrapper', async () => {
		const fetchMock = vi.fn().mockResolvedValue(jsonResponse(usersResponse(), { status: 200 }));

		const result = await fetchUsers(fetchMock);

		expect(fetchMock).toHaveBeenCalledWith('/api/users', {
			method: 'GET',
			headers: undefined,
			body: undefined
		});
		expect(result.users[0].email).toBe('ada@example.com');
		expect(USERS_QUERY_KEY).toEqual(['admin-users']);
	});

	it('creates and updates users with JSON payloads and encoded ids', async () => {
		const user = usersResponse().users[0];
		const fetchMock = vi
			.fn()
			.mockResolvedValueOnce(jsonResponse(user, { status: 201 }))
			.mockResolvedValueOnce(jsonResponse({ ...user, name: 'Grace Hopper' }, { status: 200 }));

		await createUser(fetchMock, {
			name: 'Ada Lovelace',
			email: 'ada@example.com',
			status: 'invited',
			primaryOrganizationUnitId: 'unit-id',
			managerUserId: null,
			jobTitle: 'Engineer',
			phone: null
		});
		await updateUser(fetchMock, {
			id: 'user/id',
			input: {
				name: 'Grace Hopper',
				managerUserId: null,
				jobTitle: 'Admiral',
				phone: '555-0100'
			}
		});

		expect(fetchMock).toHaveBeenNthCalledWith(1, '/api/users', {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
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
		expect(fetchMock).toHaveBeenNthCalledWith(2, '/api/users/user%2Fid', {
			method: 'PATCH',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({
				name: 'Grace Hopper',
				managerUserId: null,
				jobTitle: 'Admiral',
				phone: '555-0100'
			})
		});
	});

	it('sends status and primary organization unit updates to the expected endpoints', async () => {
		const user = usersResponse().users[0];
		const fetchMock = vi.fn(async () => jsonResponse(user, { status: 200 }));

		await activateUser(fetchMock, { id: 'user/id' });
		await deactivateUser(fetchMock, { id: 'user/id' });
		await suspendUser(fetchMock, { id: 'user/id' });
		await setPrimaryOrganizationUnit(fetchMock, {
			id: 'user/id',
			organizationUnitId: null
		});
		await softDeleteUser(fetchMock, { id: 'user/id' });

		expect(fetchMock).toHaveBeenNthCalledWith(1, '/api/users/user%2Fid/activate', {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({})
		});
		expect(fetchMock).toHaveBeenNthCalledWith(2, '/api/users/user%2Fid/deactivate', {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({})
		});
		expect(fetchMock).toHaveBeenNthCalledWith(3, '/api/users/user%2Fid/suspend', {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({})
		});
		expect(fetchMock).toHaveBeenNthCalledWith(4, '/api/users/user%2Fid/primary-organization-unit', {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({ organizationUnitId: null })
		});
		expect(fetchMock).toHaveBeenNthCalledWith(5, '/api/users/user%2Fid', {
			method: 'DELETE',
			headers: undefined,
			body: undefined
		});
	});

	it('assigns and removes scoped roles and groups', async () => {
		const fetchMock = vi
			.fn()
			.mockResolvedValueOnce(jsonResponse({ id: 'assignment-id' }, { status: 201 }))
			.mockResolvedValueOnce(jsonResponse({ ok: true }, { status: 200 }))
			.mockResolvedValueOnce(jsonResponse({ userId: 'user-id', groupId: 'group-id' }, { status: 201 }))
			.mockResolvedValueOnce(jsonResponse({ ok: true }, { status: 200 }));

		await assignUserRole(fetchMock, {
			id: 'user/id',
			input: { roleId: 'role-id', scopeType: 'global', organizationUnitId: null }
		});
		await removeUserRole(fetchMock, { id: 'user/id', assignmentId: 'assignment/id' });
		await assignUserGroup(fetchMock, { id: 'user/id', groupId: 'group/id' });
		await removeUserGroup(fetchMock, { id: 'user/id', groupId: 'group/id' });

		expect(fetchMock).toHaveBeenNthCalledWith(1, '/api/users/user%2Fid/roles', {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({ roleId: 'role-id', scopeType: 'global', organizationUnitId: null })
		});
		expect(fetchMock).toHaveBeenNthCalledWith(2, '/api/users/user%2Fid/roles/assignment%2Fid', {
			method: 'DELETE',
			headers: undefined,
			body: undefined
		});
		expect(fetchMock).toHaveBeenNthCalledWith(3, '/api/users/user%2Fid/groups', {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({ groupId: 'group/id' })
		});
		expect(fetchMock).toHaveBeenNthCalledWith(4, '/api/users/user%2Fid/groups/group%2Fid', {
			method: 'DELETE',
			headers: undefined,
			body: undefined
		});
	});

	it('uses public-safe messages and rejects malformed payloads', async () => {
		const serverErrorFetch = vi
			.fn()
			.mockResolvedValue(jsonResponse({ error: 'database unavailable' }, { status: 503 }));
		await expect(fetchUsers(serverErrorFetch)).rejects.toThrow('Failed to load users');

		const malformedFetch = vi.fn().mockResolvedValue(jsonResponse({ ok: true }, { status: 200 }));
		await expect(fetchUsers(malformedFetch)).rejects.toThrow('Invalid user response');
	});

	it('exports only browser-safe helpers', () => {
		const source = readFileSync(new URL('./api.ts', import.meta.url), 'utf8');
		const exportedNames = [
			...source.matchAll(/export\s+(?:async\s+)?(?:function|type|const|class)\s+(\w+)/g)
		]
			.map((match) => match[1])
			.sort();

		expect(exportedNames).toEqual([
			'AssignUserGroupVariables',
			'AssignUserRoleVariables',
			'CreateUserInput',
			'RemoveUserGroupVariables',
			'RemoveUserRoleVariables',
			'ScopedRoleAssignmentInput',
			'SetPrimaryOrganizationUnitVariables',
			'USERS_QUERY_KEY',
			'UpdateUserInput',
			'UpdateUserVariables',
			'User',
			'UserListResponse',
			'UserStatusVariables',
			'activateUser',
			'assignUserGroup',
			'assignUserRole',
			'createUser',
			'deactivateUser',
			'fetchUsers',
			'removeUserGroup',
			'removeUserRole',
			'setPrimaryOrganizationUnit',
			'softDeleteUser',
			'suspendUser',
			'updateUser'
		]);
		expect(source).not.toContain('$lib/server');
	});
});
