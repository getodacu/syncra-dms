import { publicApiErrorMessage } from '$lib/client/api-errors';

type ClientFetch = typeof fetch;
type JsonObject = Record<string, unknown>;

export const GROUPS_QUERY_KEY = ['admin-groups'] as const;

export type Group = {
	id: string;
	name: string;
	code: string;
	description?: string | null;
	organizationUnitId?: string | null;
	isActive: boolean;
	createdAt: string;
	updatedAt: string;
};

export type GroupListResponse = {
	groups: Group[];
};

export type CreateGroupInput = {
	name: string;
	code: string;
	description?: string | null;
	organizationUnitId?: string | null;
	isActive?: boolean;
};

export type UpdateGroupInput = {
	name?: string;
	code?: string;
	description?: string | null;
	organizationUnitId?: string | null;
	isActive?: boolean;
};

export type UpdateGroupVariables = {
	id: string;
	input: UpdateGroupInput;
};

export type GroupIDVariables = {
	id: string;
};

export type GroupUserVariables = {
	id: string;
	userId: string;
};

export type ScopedGroupRoleInput = {
	roleId: string;
	scopeType: 'global' | 'organization_unit' | 'organization_unit_and_children';
	organizationUnitId?: string | null;
};

export type AssignGroupRoleVariables = {
	id: string;
	input: ScopedGroupRoleInput;
};

export type RemoveGroupRoleVariables = {
	id: string;
	assignmentId: string;
};

export async function fetchGroups(fetchFn: ClientFetch): Promise<GroupListResponse> {
	return groupRequest(
		fetchFn,
		'/api/groups',
		{ method: 'GET' },
		validateGroupListResponse,
		'Failed to load groups'
	);
}

export async function createGroup(fetchFn: ClientFetch, input: CreateGroupInput) {
	return groupRequest(
		fetchFn,
		'/api/groups',
		{ method: 'POST', body: input },
		validateGroup,
		'Failed to create group'
	);
}

export async function updateGroup(fetchFn: ClientFetch, { id, input }: UpdateGroupVariables) {
	return groupRequest(
		fetchFn,
		`/api/groups/${encodeURIComponent(id)}`,
		{ method: 'PATCH', body: input },
		validateGroup,
		'Failed to update group'
	);
}

export async function deleteGroup(fetchFn: ClientFetch, { id }: GroupIDVariables) {
	return groupRequest(
		fetchFn,
		`/api/groups/${encodeURIComponent(id)}`,
		{ method: 'DELETE' },
		validateOKResponse,
		'Failed to delete group'
	);
}

export async function addGroupUser(fetchFn: ClientFetch, { id, userId }: GroupUserVariables) {
	return groupRequest(
		fetchFn,
		`/api/groups/${encodeURIComponent(id)}/users`,
		{ method: 'POST', body: { userId } },
		validateRecordResponse,
		'Failed to add group user'
	);
}

export async function removeGroupUser(fetchFn: ClientFetch, { id, userId }: GroupUserVariables) {
	return groupRequest(
		fetchFn,
		`/api/groups/${encodeURIComponent(id)}/users/${encodeURIComponent(userId)}`,
		{ method: 'DELETE' },
		validateOKResponse,
		'Failed to remove group user'
	);
}

export async function assignGroupRole(
	fetchFn: ClientFetch,
	{ id, input }: AssignGroupRoleVariables
) {
	return groupRequest(
		fetchFn,
		`/api/groups/${encodeURIComponent(id)}/roles`,
		{ method: 'POST', body: input },
		validateRecordResponse,
		'Failed to assign group role'
	);
}

export async function removeGroupRole(
	fetchFn: ClientFetch,
	{ id, assignmentId }: RemoveGroupRoleVariables
) {
	return groupRequest(
		fetchFn,
		`/api/groups/${encodeURIComponent(id)}/roles/${encodeURIComponent(assignmentId)}`,
		{ method: 'DELETE' },
		validateOKResponse,
		'Failed to remove group role'
	);
}

async function groupRequest<T>(
	fetchFn: ClientFetch,
	path: string,
	options: {
		method: 'DELETE' | 'GET' | 'PATCH' | 'POST';
		body?: unknown;
	},
	validate: (data: unknown) => T,
	fallback: string
) {
	const response = await fetchFn(path, {
		method: options.method,
		headers: options.body === undefined ? undefined : { 'content-type': 'application/json' },
		body: options.body === undefined ? undefined : JSON.stringify(options.body)
	});
	const data = await readResponseJSON(response);

	if (!response.ok) {
		throw new Error(publicApiErrorMessage(response.status, data, fallback));
	}

	return validate(data);
}

async function readResponseJSON(response: Response): Promise<unknown> {
	let text: string;
	try {
		text = await response.text();
	} catch {
		return null;
	}
	if (!text.trim()) return null;
	try {
		return JSON.parse(text) as unknown;
	} catch {
		return null;
	}
}

function validateGroupListResponse(data: unknown): GroupListResponse {
	if (!isJsonObject(data) || !Array.isArray(data.groups)) invalidGroupResponse();
	return { groups: data.groups.map(validateGroup) };
}

function validateGroup(data: unknown): Group {
	if (!isJsonObject(data)) invalidGroupResponse();
	const output: Group = {
		id: requiredString(data.id),
		name: requiredString(data.name),
		code: requiredString(data.code),
		isActive: requiredBoolean(data.isActive),
		createdAt: requiredString(data.createdAt),
		updatedAt: requiredString(data.updatedAt)
	};
	assignOptionalString(output, 'description', data.description);
	assignOptionalString(output, 'organizationUnitId', data.organizationUnitId);
	return output;
}

function validateRecordResponse(data: unknown) {
	if (!isJsonObject(data)) invalidGroupResponse();
	return data;
}

function validateOKResponse(data: unknown) {
	if (!isJsonObject(data) || typeof data.ok !== 'boolean') invalidGroupResponse();
	return { ok: data.ok };
}

function requiredString(value: unknown) {
	if (typeof value !== 'string' || value.trim() === '') invalidGroupResponse();
	return value;
}

function requiredBoolean(value: unknown) {
	if (typeof value !== 'boolean') invalidGroupResponse();
	return value;
}

function optionalString(value: unknown) {
	if (value === undefined || value === null) return value;
	if (typeof value === 'string') return value;
	invalidGroupResponse();
}

function assignOptionalString<T extends object>(output: T, key: keyof T & string, value: unknown) {
	const parsed = optionalString(value);
	if (parsed !== undefined) (output as Record<string, unknown>)[key] = parsed;
}

function isJsonObject(value: unknown): value is JsonObject {
	return typeof value === 'object' && value !== null && !Array.isArray(value);
}

function invalidGroupResponse(): never {
	throw new Error('Invalid group response');
}
