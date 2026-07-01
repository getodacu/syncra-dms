import { beforeEach, describe, expect, it, vi } from 'vitest';

const rbacMocks = vi.hoisted(() => {
	class MockRbacApiError extends Error {
		status: number;

		constructor(status: number, message: string) {
			super(message);
			this.name = 'RbacApiError';
			this.status = status;
		}
	}

	return {
		RbacApiError: MockRbacApiError,
		addGroupUser: vi.fn(),
		assignGroupRole: vi.fn(),
		assignRolePermission: vi.fn(),
		assignUserRole: vi.fn(),
		createRole: vi.fn(),
		createUser: vi.fn(),
		getMyPermissions: vi.fn(),
		listUsers: vi.fn(),
		updateUser: vi.fn()
	};
});

vi.mock('$lib/server/rbac', () => ({
	RbacApiError: rbacMocks.RbacApiError,
	addGroupUser: rbacMocks.addGroupUser,
	assignGroupRole: rbacMocks.assignGroupRole,
	assignRolePermission: rbacMocks.assignRolePermission,
	assignUserRole: rbacMocks.assignUserRole,
	createRole: rbacMocks.createRole,
	createUser: rbacMocks.createUser,
	getMyPermissions: rbacMocks.getMyPermissions,
	isRbacApiError: (error: unknown) => error instanceof rbacMocks.RbacApiError,
	listUsers: rbacMocks.listUsers,
	updateUser: rbacMocks.updateUser
}));

import { GET as getMyPermissionsRoute } from '../me/permissions/+server';
import { POST as createRoleRoute } from '../roles/+server';
import { POST as assignRolePermissionRoute } from '../roles/[id]/permissions/+server';
import { POST as addGroupUserRoute } from '../groups/[id]/users/+server';
import { POST as assignGroupRoleRoute } from '../groups/[id]/roles/+server';
import { GET as listUsersRoute, POST as createUserRoute } from '../users/+server';
import { PATCH as updateUserRoute } from '../users/[id]/+server';
import { POST as assignUserRoleRoute } from '../users/[id]/roles/+server';

const cookieHeader = 'auth.session_token=token';
const user = {
	id: 'user-id',
	name: 'Ada Lovelace',
	email: 'ada@example.com',
	status: 'active',
	createdAt: '2026-06-30T00:00:00Z',
	updatedAt: '2026-06-30T00:00:00Z'
};

describe('RBAC Svelte API routes', () => {
	beforeEach(() => {
		vi.clearAllMocks();
	});

	it('requires authenticated users for protected reads', async () => {
		const response = await listUsersRoute(
			event({ locals: { user: null, permissions: [] }, path: '/api/users' }) as never
		);

		expect(response.status).toBe(401);
		expect(await response.json()).toEqual({ error: 'Authentication required' });
		expect(rbacMocks.listUsers).not.toHaveBeenCalled();
	});

	it('requires local permissions before proxying', async () => {
		const response = await listUsersRoute(
			event({
				locals: { user: { role: 'user' }, permissions: [] },
				path: '/api/users'
			}) as never
		);

		expect(response.status).toBe(403);
		expect(await response.json()).toEqual({ error: 'permission required' });
		expect(rbacMocks.listUsers).not.toHaveBeenCalled();
	});

	it('proxies allowed reads with fetch and cookie context', async () => {
		rbacMocks.listUsers.mockResolvedValue({ users: [user] });
		const fetchMock = vi.fn();

		const response = await listUsersRoute(
			event({
				fetch: fetchMock,
				locals: { user: { role: 'user' }, permissions: ['user.view'] },
				path: '/api/users'
			}) as never
		);

		expect(response.status).toBe(200);
		expect(await response.json()).toEqual({ users: [user] });
		expect(rbacMocks.listUsers).toHaveBeenCalledWith(fetchMock, cookieHeader);
	});

	it('allows system administrators through local permission checks', async () => {
		rbacMocks.createRole.mockResolvedValue({ id: 'role-id', name: 'Unit Manager' });
		const fetchMock = vi.fn();
		const input = { name: 'Unit Manager', code: 'unit_manager', isActive: true };

		const response = await createRoleRoute(
			event({
				fetch: fetchMock,
				locals: { user: { role: 'user' }, permissions: ['system.admin'] },
				method: 'POST',
				path: '/api/roles',
				body: input
			}) as never
		);

		expect(response.status).toBe(201);
		expect(await response.json()).toEqual({ id: 'role-id', name: 'Unit Manager' });
		expect(rbacMocks.createRole).toHaveBeenCalledWith(fetchMock, cookieHeader, input);
	});

	it('proxies representative mutation routes with parsed request bodies', async () => {
		rbacMocks.createUser.mockResolvedValue(user);
		rbacMocks.updateUser.mockResolvedValue({ ...user, name: 'Grace Hopper' });
		rbacMocks.assignUserRole.mockResolvedValue({ id: 'assignment-id' });
		rbacMocks.assignRolePermission.mockResolvedValue({ id: 'permission-id' });
		rbacMocks.addGroupUser.mockResolvedValue({ groupId: 'group-id', userId: 'user-id' });
		rbacMocks.assignGroupRole.mockResolvedValue({ id: 'group-role-id' });
		const fetchMock = vi.fn();

		await createUserRoute(
			event({
				fetch: fetchMock,
				locals: { user: { role: 'user' }, permissions: ['user.create'] },
				method: 'POST',
				path: '/api/users',
				body: {
					name: 'Ada Lovelace',
					email: 'ada@example.com',
					status: 'invited',
					primaryOrganizationUnitId: 'unit-id',
					managerUserId: null,
					jobTitle: 'Engineer',
					phone: null
				}
			}) as never
		);
		await updateUserRoute(
			event({
				fetch: fetchMock,
				locals: { user: { role: 'user' }, permissions: ['user.update'] },
				method: 'PATCH',
				path: '/api/users/user-id',
				params: { id: 'user-id' },
				body: { name: 'Grace Hopper', jobTitle: 'Admiral' }
			}) as never
		);
		await assignUserRoleRoute(
			event({
				fetch: fetchMock,
				locals: { user: { role: 'user' }, permissions: ['user.assign_role'] },
				method: 'POST',
				path: '/api/users/user-id/roles',
				params: { id: 'user-id' },
				body: {
					roleId: 'role-id',
					scopeType: 'organization_unit',
					organizationUnitId: 'unit-id'
				}
			}) as never
		);
		await assignRolePermissionRoute(
			event({
				fetch: fetchMock,
				locals: { user: { role: 'user' }, permissions: ['role.assign_permissions'] },
				method: 'POST',
				path: '/api/roles/role-id/permissions',
				params: { id: 'role-id' },
				body: { permissionId: 'permission-id' }
			}) as never
		);
		await addGroupUserRoute(
			event({
				fetch: fetchMock,
				locals: { user: { role: 'user' }, permissions: ['group.manage_users'] },
				method: 'POST',
				path: '/api/groups/group-id/users',
				params: { id: 'group-id' },
				body: { userId: 'user-id' }
			}) as never
		);
		await assignGroupRoleRoute(
			event({
				fetch: fetchMock,
				locals: { user: { role: 'user' }, permissions: ['group.assign_roles'] },
				method: 'POST',
				path: '/api/groups/group-id/roles',
				params: { id: 'group-id' },
				body: { roleId: 'role-id', scopeType: 'global', organizationUnitId: null }
			}) as never
		);

		expect(rbacMocks.createUser).toHaveBeenCalledWith(fetchMock, cookieHeader, {
			name: 'Ada Lovelace',
			email: 'ada@example.com',
			status: 'invited',
			primaryOrganizationUnitId: 'unit-id',
			managerUserId: null,
			jobTitle: 'Engineer',
			phone: null
		});
		expect(rbacMocks.updateUser).toHaveBeenCalledWith(fetchMock, cookieHeader, 'user-id', {
			name: 'Grace Hopper',
			jobTitle: 'Admiral'
		});
		expect(rbacMocks.assignUserRole).toHaveBeenCalledWith(fetchMock, cookieHeader, 'user-id', {
			roleId: 'role-id',
			scopeType: 'organization_unit',
			organizationUnitId: 'unit-id'
		});
		expect(rbacMocks.assignRolePermission).toHaveBeenCalledWith(
			fetchMock,
			cookieHeader,
			'role-id',
			'permission-id'
		);
		expect(rbacMocks.addGroupUser).toHaveBeenCalledWith(
			fetchMock,
			cookieHeader,
			'group-id',
			'user-id'
		);
		expect(rbacMocks.assignGroupRole).toHaveBeenCalledWith(fetchMock, cookieHeader, 'group-id', {
			roleId: 'role-id',
			scopeType: 'global',
			organizationUnitId: null
		});
	});

	it('allows authenticated users to read their effective permissions', async () => {
		rbacMocks.getMyPermissions.mockResolvedValue({
			permissions: [{ code: 'user.view', scopeType: 'global', source: 'user_role' }]
		});
		const fetchMock = vi.fn();

		const response = await getMyPermissionsRoute(
			event({
				fetch: fetchMock,
				locals: { user: { role: 'user' }, permissions: [] },
				path: '/api/me/permissions'
			}) as never
		);

		expect(response.status).toBe(200);
		expect(await response.json()).toEqual({
			permissions: [{ code: 'user.view', scopeType: 'global', source: 'user_role' }]
		});
		expect(rbacMocks.getMyPermissions).toHaveBeenCalledWith(fetchMock, cookieHeader);
	});

	it('maps known backend errors to public-safe JSON responses', async () => {
		rbacMocks.listUsers.mockRejectedValue(
			new rbacMocks.RbacApiError(503, 'database unavailable')
		);

		const response = await listUsersRoute(
			event({
				locals: { user: { role: 'user' }, permissions: ['user.view'] },
				path: '/api/users'
			}) as never
		);

		expect(response.status).toBe(502);
		expect(await response.json()).toEqual({ error: 'Failed to load users' });
	});
});

function event({
	fetch = vi.fn(),
	locals,
	method = 'GET',
	path,
	params = {},
	body
}: {
	fetch?: typeof globalThis.fetch;
	locals: { user: { role: string } | null; permissions?: string[] };
	method?: string;
	path: string;
	params?: Record<string, string>;
	body?: unknown;
}) {
	return {
		fetch,
		locals,
		params,
		request: new Request(`http://localhost${path}`, {
			method,
			headers: { cookie: cookieHeader, ...(body === undefined ? {} : { 'content-type': 'application/json' }) },
			body: body === undefined ? undefined : JSON.stringify(body)
		})
	};
}
