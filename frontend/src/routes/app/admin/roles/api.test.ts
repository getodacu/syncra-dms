import { describe, expect, it, vi } from 'vitest';

import {
	ROLES_QUERY_KEY,
	assignRolePermission,
	createRole,
	deleteRole,
	fetchPermissionCategories,
	fetchPermissions,
	fetchRolePermissions,
	fetchRoles,
	groupPermissionsByCategory,
	removeRolePermission,
	updateRole
} from './api';
import type { Permission, RoleListResponse } from './api';

function jsonResponse(body: unknown, init?: ResponseInit) {
	return new Response(JSON.stringify(body), {
		headers: { 'content-type': 'application/json' },
		...init
	});
}

function rolesResponse(): RoleListResponse {
	return {
		roles: [
			{
				id: 'role-id',
				name: 'Unit Manager',
				code: 'unit_manager',
				description: null,
				isSystem: false,
				isActive: true,
				createdAt: '2026-06-30T00:00:00Z',
				updatedAt: '2026-06-30T00:00:00Z'
			}
		]
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
		createdAt: '2026-06-30T00:00:00Z',
		updatedAt: '2026-06-30T00:00:00Z',
		...overrides
	};
}

describe('admin roles browser API', () => {
	it('fetches roles through the Svelte API wrapper', async () => {
		const fetchMock = vi.fn().mockResolvedValue(jsonResponse(rolesResponse(), { status: 200 }));

		const result = await fetchRoles(fetchMock);

		expect(fetchMock).toHaveBeenCalledWith('/api/roles', {
			method: 'GET',
			headers: undefined,
			body: undefined
		});
		expect(result.roles[0].code).toBe('unit_manager');
		expect(ROLES_QUERY_KEY).toEqual(['admin-roles']);
	});

	it('creates, updates, and deletes roles with encoded ids', async () => {
		const role = rolesResponse().roles[0];
		const fetchMock = vi
			.fn()
			.mockResolvedValueOnce(jsonResponse(role, { status: 201 }))
			.mockResolvedValueOnce(jsonResponse({ ...role, name: 'Legal Reviewer' }, { status: 200 }))
			.mockResolvedValueOnce(jsonResponse({ ok: true }, { status: 200 }));

		await createRole(fetchMock, {
			name: 'Legal Reviewer',
			code: 'legal_reviewer',
			description: 'Reviews legal documents',
			isActive: true
		});
		await updateRole(fetchMock, {
			id: 'role/id',
			input: { name: 'Legal Reviewer', description: null, isActive: false }
		});
		await deleteRole(fetchMock, { id: 'role/id' });

		expect(fetchMock).toHaveBeenNthCalledWith(1, '/api/roles', {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({
				name: 'Legal Reviewer',
				code: 'legal_reviewer',
				description: 'Reviews legal documents',
				isActive: true
			})
		});
		expect(fetchMock).toHaveBeenNthCalledWith(2, '/api/roles/role%2Fid', {
			method: 'PATCH',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({ name: 'Legal Reviewer', description: null, isActive: false })
		});
		expect(fetchMock).toHaveBeenNthCalledWith(3, '/api/roles/role%2Fid', {
			method: 'DELETE',
			headers: undefined,
			body: undefined
		});
	});

	it('loads permissions and manages role permission assignments', async () => {
		const fetchMock = vi
			.fn()
			.mockResolvedValueOnce(jsonResponse({ permissions: [permission()] }, { status: 200 }))
			.mockResolvedValueOnce(jsonResponse({ categories: ['User Management'] }, { status: 200 }))
			.mockResolvedValueOnce(jsonResponse({ permissions: [permission()] }, { status: 200 }))
			.mockResolvedValueOnce(jsonResponse(permission(), { status: 201 }))
			.mockResolvedValueOnce(jsonResponse({ ok: true }, { status: 200 }));

		await fetchPermissions(fetchMock);
		await fetchPermissionCategories(fetchMock);
		await fetchRolePermissions(fetchMock, { id: 'role/id' });
		await assignRolePermission(fetchMock, { id: 'role/id', permissionId: 'permission/id' });
		await removeRolePermission(fetchMock, { id: 'role/id', permissionId: 'permission/id' });

		expect(fetchMock).toHaveBeenNthCalledWith(
			1,
			'/api/permissions',
			expect.objectContaining({ method: 'GET' })
		);
		expect(fetchMock).toHaveBeenNthCalledWith(
			2,
			'/api/permissions/categories',
			expect.objectContaining({ method: 'GET' })
		);
		expect(fetchMock).toHaveBeenNthCalledWith(
			3,
			'/api/roles/role%2Fid/permissions',
			expect.objectContaining({ method: 'GET' })
		);
		expect(fetchMock).toHaveBeenNthCalledWith(4, '/api/roles/role%2Fid/permissions', {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({ permissionId: 'permission/id' })
		});
		expect(fetchMock).toHaveBeenNthCalledWith(5, '/api/roles/role%2Fid/permissions/permission%2Fid', {
			method: 'DELETE',
			headers: undefined,
			body: undefined
		});
	});

	it('groups permissions by category in display order', () => {
		const grouped = groupPermissionsByCategory([
			permission({ code: 'role.view', category: 'Role Management' }),
			permission({ code: 'user.view', category: 'User Management' }),
			permission({ code: 'user.create', category: 'User Management' })
		]);

		expect(grouped).toEqual([
			{
				category: 'Role Management',
				permissions: [permission({ code: 'role.view', category: 'Role Management' })]
			},
			{
				category: 'User Management',
				permissions: [
					permission({ code: 'user.create', category: 'User Management' }),
					permission({ code: 'user.view', category: 'User Management' })
				]
			}
		]);
	});

	it('uses public-safe messages and rejects malformed payloads', async () => {
		const serverErrorFetch = vi
			.fn()
			.mockResolvedValue(jsonResponse({ error: 'database unavailable' }, { status: 503 }));
		await expect(fetchRoles(serverErrorFetch)).rejects.toThrow('Failed to load roles');

		const malformedFetch = vi.fn().mockResolvedValue(jsonResponse({ ok: true }, { status: 200 }));
		await expect(fetchRoles(malformedFetch)).rejects.toThrow('Invalid role response');
	});
});
