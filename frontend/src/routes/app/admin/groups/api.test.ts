import { describe, expect, it, vi } from 'vitest';

import {
	GROUPS_QUERY_KEY,
	addGroupUser,
	assignGroupRole,
	createGroup,
	deleteGroup,
	fetchGroups,
	removeGroupRole,
	removeGroupUser,
	updateGroup
} from './api';
import type { GroupListResponse } from './api';

function jsonResponse(body: unknown, init?: ResponseInit) {
	return new Response(JSON.stringify(body), {
		headers: { 'content-type': 'application/json' },
		...init
	});
}

function groupsResponse(): GroupListResponse {
	return {
		groups: [
			{
				id: 'group-id',
				name: 'Legal Reviewers',
				code: 'legal_reviewers',
				description: null,
				organizationUnitId: null,
				isActive: true,
				createdAt: '2026-06-30T00:00:00Z',
				updatedAt: '2026-06-30T00:00:00Z'
			}
		]
	};
}

describe('admin groups browser API', () => {
	it('fetches groups through the Svelte API wrapper', async () => {
		const fetchMock = vi.fn().mockResolvedValue(jsonResponse(groupsResponse(), { status: 200 }));

		const result = await fetchGroups(fetchMock);

		expect(fetchMock).toHaveBeenCalledWith('/api/groups', {
			method: 'GET',
			headers: undefined,
			body: undefined
		});
		expect(result.groups[0].code).toBe('legal_reviewers');
		expect(GROUPS_QUERY_KEY).toEqual(['admin-groups']);
	});

	it('creates, updates, and deletes groups with encoded ids', async () => {
		const group = groupsResponse().groups[0];
		const fetchMock = vi
			.fn()
			.mockResolvedValueOnce(jsonResponse(group, { status: 201 }))
			.mockResolvedValueOnce(jsonResponse({ ...group, name: 'Finance Approvers' }, { status: 200 }))
			.mockResolvedValueOnce(jsonResponse({ ok: true }, { status: 200 }));

		await createGroup(fetchMock, {
			name: 'Finance Approvers',
			code: 'finance_approvers',
			description: 'Approves finance documents',
			organizationUnitId: 'unit-id',
			isActive: true
		});
		await updateGroup(fetchMock, {
			id: 'group/id',
			input: { name: 'Finance Approvers', description: null, isActive: false }
		});
		await deleteGroup(fetchMock, { id: 'group/id' });

		expect(fetchMock).toHaveBeenNthCalledWith(1, '/api/groups', {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({
				name: 'Finance Approvers',
				code: 'finance_approvers',
				description: 'Approves finance documents',
				organizationUnitId: 'unit-id',
				isActive: true
			})
		});
		expect(fetchMock).toHaveBeenNthCalledWith(2, '/api/groups/group%2Fid', {
			method: 'PATCH',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({ name: 'Finance Approvers', description: null, isActive: false })
		});
		expect(fetchMock).toHaveBeenNthCalledWith(3, '/api/groups/group%2Fid', {
			method: 'DELETE',
			headers: undefined,
			body: undefined
		});
	});

	it('adds and removes group users', async () => {
		const fetchMock = vi
			.fn()
			.mockResolvedValueOnce(jsonResponse({ groupId: 'group-id', userId: 'user-id' }, { status: 201 }))
			.mockResolvedValueOnce(jsonResponse({ ok: true }, { status: 200 }));

		await addGroupUser(fetchMock, { id: 'group/id', userId: 'user/id' });
		await removeGroupUser(fetchMock, { id: 'group/id', userId: 'user/id' });

		expect(fetchMock).toHaveBeenNthCalledWith(1, '/api/groups/group%2Fid/users', {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({ userId: 'user/id' })
		});
		expect(fetchMock).toHaveBeenNthCalledWith(2, '/api/groups/group%2Fid/users/user%2Fid', {
			method: 'DELETE',
			headers: undefined,
			body: undefined
		});
	});

	it('assigns and removes scoped group roles', async () => {
		const fetchMock = vi
			.fn()
			.mockResolvedValueOnce(jsonResponse({ id: 'assignment-id' }, { status: 201 }))
			.mockResolvedValueOnce(jsonResponse({ ok: true }, { status: 200 }));

		await assignGroupRole(fetchMock, {
			id: 'group/id',
			input: { roleId: 'role-id', scopeType: 'organization_unit', organizationUnitId: 'unit/id' }
		});
		await removeGroupRole(fetchMock, { id: 'group/id', assignmentId: 'assignment/id' });

		expect(fetchMock).toHaveBeenNthCalledWith(1, '/api/groups/group%2Fid/roles', {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({
				roleId: 'role-id',
				scopeType: 'organization_unit',
				organizationUnitId: 'unit/id'
			})
		});
		expect(fetchMock).toHaveBeenNthCalledWith(2, '/api/groups/group%2Fid/roles/assignment%2Fid', {
			method: 'DELETE',
			headers: undefined,
			body: undefined
		});
	});

	it('uses public-safe messages and rejects malformed payloads', async () => {
		const serverErrorFetch = vi
			.fn()
			.mockResolvedValue(jsonResponse({ error: 'database unavailable' }, { status: 503 }));
		await expect(fetchGroups(serverErrorFetch)).rejects.toThrow('Failed to load groups');

		const malformedFetch = vi.fn().mockResolvedValue(jsonResponse({ ok: true }, { status: 200 }));
		await expect(fetchGroups(malformedFetch)).rejects.toThrow('Invalid group response');
	});
});
