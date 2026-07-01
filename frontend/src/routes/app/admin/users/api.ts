import { publicApiErrorMessage } from '$lib/client/api-errors';

type ClientFetch = typeof fetch;
type JsonObject = Record<string, unknown>;

export const USERS_QUERY_KEY = ['admin-users'] as const;

export type User = {
	id: string;
	name: string;
	email: string;
	emailVerified: boolean;
	image?: string | null;
	preferredLanguage: 'en' | 'ro';
	role: 'admin' | 'user';
	status: 'invited' | 'active' | 'inactive' | 'suspended' | 'deleted';
	primaryOrganizationUnitId?: string | null;
	managerUserId?: string | null;
	jobTitle?: string | null;
	phone?: string | null;
	lastLoginAt?: string | null;
	deletedAt?: string | null;
	createdAt: string;
	updatedAt: string;
};

export type UserListResponse = {
	users: User[];
};

export type CreateUserInput = {
	name: string;
	email: string;
	status?: User['status'];
	primaryOrganizationUnitId?: string | null;
	managerUserId?: string | null;
	jobTitle?: string | null;
	phone?: string | null;
};

export type UpdateUserInput = {
	name?: string;
	managerUserId?: string | null;
	jobTitle?: string | null;
	phone?: string | null;
};

export type UpdateUserVariables = {
	id: string;
	input: UpdateUserInput;
};

export type UserStatusVariables = {
	id: string;
};

export type SetPrimaryOrganizationUnitVariables = {
	id: string;
	organizationUnitId: string | null;
};

export type ScopedRoleAssignmentInput = {
	roleId: string;
	scopeType: 'global' | 'organization_unit' | 'organization_unit_and_children';
	organizationUnitId?: string | null;
};

export type AssignUserRoleVariables = {
	id: string;
	input: ScopedRoleAssignmentInput;
};

export type RemoveUserRoleVariables = {
	id: string;
	assignmentId: string;
};

export type AssignUserGroupVariables = {
	id: string;
	groupId: string;
};

export type RemoveUserGroupVariables = {
	id: string;
	groupId: string;
};

export async function fetchUsers(fetchFn: ClientFetch): Promise<UserListResponse> {
	return usersRequest(
		fetchFn,
		'/api/users',
		{ method: 'GET' },
		validateUserListResponse,
		'Failed to load users'
	);
}

export async function createUser(fetchFn: ClientFetch, input: CreateUserInput) {
	return usersRequest(
		fetchFn,
		'/api/users',
		{ method: 'POST', body: input },
		validateUser,
		'Failed to create user'
	);
}

export async function updateUser(fetchFn: ClientFetch, { id, input }: UpdateUserVariables) {
	return usersRequest(
		fetchFn,
		`/api/users/${encodeURIComponent(id)}`,
		{ method: 'PATCH', body: input },
		validateUser,
		'Failed to update user'
	);
}

export async function activateUser(fetchFn: ClientFetch, { id }: UserStatusVariables) {
	return userStatusRequest(fetchFn, id, 'activate', 'Failed to activate user');
}

export async function deactivateUser(fetchFn: ClientFetch, { id }: UserStatusVariables) {
	return userStatusRequest(fetchFn, id, 'deactivate', 'Failed to deactivate user');
}

export async function suspendUser(fetchFn: ClientFetch, { id }: UserStatusVariables) {
	return userStatusRequest(fetchFn, id, 'suspend', 'Failed to suspend user');
}

export async function softDeleteUser(fetchFn: ClientFetch, { id }: UserStatusVariables) {
	return usersRequest(
		fetchFn,
		`/api/users/${encodeURIComponent(id)}`,
		{ method: 'DELETE' },
		validateUser,
		'Failed to delete user'
	);
}

export async function setPrimaryOrganizationUnit(
	fetchFn: ClientFetch,
	{ id, organizationUnitId }: SetPrimaryOrganizationUnitVariables
) {
	return usersRequest(
		fetchFn,
		`/api/users/${encodeURIComponent(id)}/primary-organization-unit`,
		{ method: 'POST', body: { organizationUnitId } },
		validateUser,
		'Failed to update primary organization unit'
	);
}

export async function assignUserRole(
	fetchFn: ClientFetch,
	{ id, input }: AssignUserRoleVariables
) {
	return usersRequest(
		fetchFn,
		`/api/users/${encodeURIComponent(id)}/roles`,
		{ method: 'POST', body: input },
		validateRecordResponse,
		'Failed to assign user role'
	);
}

export async function removeUserRole(
	fetchFn: ClientFetch,
	{ id, assignmentId }: RemoveUserRoleVariables
) {
	return usersRequest(
		fetchFn,
		`/api/users/${encodeURIComponent(id)}/roles/${encodeURIComponent(assignmentId)}`,
		{ method: 'DELETE' },
		validateOKResponse,
		'Failed to remove user role'
	);
}

export async function assignUserGroup(
	fetchFn: ClientFetch,
	{ id, groupId }: AssignUserGroupVariables
) {
	return usersRequest(
		fetchFn,
		`/api/users/${encodeURIComponent(id)}/groups`,
		{ method: 'POST', body: { groupId } },
		validateRecordResponse,
		'Failed to assign user group'
	);
}

export async function removeUserGroup(
	fetchFn: ClientFetch,
	{ id, groupId }: RemoveUserGroupVariables
) {
	return usersRequest(
		fetchFn,
		`/api/users/${encodeURIComponent(id)}/groups/${encodeURIComponent(groupId)}`,
		{ method: 'DELETE' },
		validateOKResponse,
		'Failed to remove user group'
	);
}

async function userStatusRequest(
	fetchFn: ClientFetch,
	id: string,
	action: 'activate' | 'deactivate' | 'suspend',
	fallback: string
) {
	return usersRequest(
		fetchFn,
		`/api/users/${encodeURIComponent(id)}/${action}`,
		{ method: 'POST', body: {} },
		validateUser,
		fallback
	);
}

async function usersRequest<T>(
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

function validateUserListResponse(data: unknown): UserListResponse {
	if (!isJsonObject(data) || !Array.isArray(data.users)) invalidUserResponse();
	return { users: data.users.map(validateUser) };
}

function validateUser(data: unknown): User {
	if (!isJsonObject(data)) invalidUserResponse();
	const output: User = {
		id: requiredString(data.id),
		name: requiredString(data.name),
		email: requiredString(data.email),
		emailVerified: requiredBoolean(data.emailVerified),
		preferredLanguage: requiredPreferredLanguage(data.preferredLanguage),
		role: requiredUserRole(data.role),
		status: requiredUserStatus(data.status),
		createdAt: requiredString(data.createdAt),
		updatedAt: requiredString(data.updatedAt)
	};
	assignOptionalString(output, 'image', data.image);
	assignOptionalString(output, 'primaryOrganizationUnitId', data.primaryOrganizationUnitId);
	assignOptionalString(output, 'managerUserId', data.managerUserId);
	assignOptionalString(output, 'jobTitle', data.jobTitle);
	assignOptionalString(output, 'phone', data.phone);
	assignOptionalString(output, 'lastLoginAt', data.lastLoginAt);
	assignOptionalString(output, 'deletedAt', data.deletedAt);
	return output;
}

function validateRecordResponse(data: unknown) {
	if (!isJsonObject(data)) invalidUserResponse();
	return data;
}

function validateOKResponse(data: unknown) {
	if (!isJsonObject(data) || typeof data.ok !== 'boolean') invalidUserResponse();
	return { ok: data.ok };
}

function requiredString(value: unknown) {
	if (typeof value !== 'string' || value.trim() === '') invalidUserResponse();
	return value;
}

function requiredBoolean(value: unknown) {
	if (typeof value !== 'boolean') invalidUserResponse();
	return value;
}

function optionalString(value: unknown) {
	if (value === undefined || value === null) return value;
	if (typeof value === 'string') return value;
	invalidUserResponse();
}

function assignOptionalString<T extends object>(output: T, key: keyof T & string, value: unknown) {
	const parsed = optionalString(value);
	if (parsed !== undefined) (output as Record<string, unknown>)[key] = parsed;
}

function requiredPreferredLanguage(value: unknown): User['preferredLanguage'] {
	if (value === 'en' || value === 'ro') return value;
	invalidUserResponse();
}

function requiredUserRole(value: unknown): User['role'] {
	if (value === 'admin' || value === 'user') return value;
	invalidUserResponse();
}

function requiredUserStatus(value: unknown): User['status'] {
	if (
		value === 'invited' ||
		value === 'active' ||
		value === 'inactive' ||
		value === 'suspended' ||
		value === 'deleted'
	) {
		return value;
	}
	invalidUserResponse();
}

function isJsonObject(value: unknown): value is JsonObject {
	return typeof value === 'object' && value !== null && !Array.isArray(value);
}

function invalidUserResponse(): never {
	throw new Error('Invalid user response');
}
