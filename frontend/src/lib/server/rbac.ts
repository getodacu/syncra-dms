import { apiBaseUrl, internalAPIHeaders } from './internal-api';
import { publicErrorMessage, publicErrorStatus } from './public-errors';

type ServerFetch = typeof fetch;
type HTTPMethod = 'DELETE' | 'GET' | 'PATCH' | 'POST';

export type PreferredLanguage = 'en' | 'ro';
export type UserRole = 'admin' | 'user';
export type UserStatus = 'invited' | 'active' | 'inactive' | 'suspended' | 'deleted';
export type ScopeType = 'global' | 'organization_unit' | 'organization_unit_and_children';

export type User = {
	id: string;
	name: string;
	email: string;
	emailVerified: boolean;
	image?: string | null;
	preferredLanguage: PreferredLanguage;
	role: UserRole;
	status: UserStatus;
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
	status?: UserStatus;
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

export type Role = {
	id: string;
	name: string;
	code: string;
	description?: string | null;
	isSystem: boolean;
	isActive: boolean;
	createdAt: string;
	updatedAt: string;
};

export type RoleListResponse = {
	roles: Role[];
};

export type CreateRoleInput = {
	name: string;
	code: string;
	description?: string | null;
	isActive?: boolean;
};

export type UpdateRoleInput = {
	name?: string;
	code?: string;
	description?: string | null;
	isActive?: boolean;
};

export type Permission = {
	id: string;
	code: string;
	name: string;
	description?: string | null;
	category: string;
	isSystem: boolean;
	createdAt: string;
	updatedAt: string;
};

export type PermissionListResponse = {
	permissions: Permission[];
};

export type PermissionCategoriesResponse = {
	categories: string[];
};

export type PermissionGrantSource = 'user_role' | 'group_role' | 'organization_unit_role';

export type PermissionGrant = {
	code: string;
	scopeType: ScopeType;
	organizationUnitId?: string | null;
	source: PermissionGrantSource;
};

export type MyPermissionsResponse = {
	permissions: PermissionGrant[];
};

export type CheckPermissionInput = {
	permission: string;
	organizationUnitId?: string | null;
};

export type CheckPermissionResponse = {
	allowed: boolean;
};

export type ScopedRoleAssignmentInput = {
	roleId: string;
	scopeType: ScopeType;
	organizationUnitId?: string | null;
};

export type UserRoleAssignment = {
	id: string;
	userId: string;
	roleId: string;
	scopeType: ScopeType;
	organizationUnitId?: string | null;
	createdAt: string;
	updatedAt: string;
};

export type UserGroupAssignment = {
	userId: string;
	groupId: string;
	createdAt: string;
};

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

export type GroupRoleAssignment = {
	id: string;
	groupId: string;
	roleId: string;
	scopeType: ScopeType;
	organizationUnitId?: string | null;
	createdAt: string;
};

export type GroupUserAssignment = {
	groupId: string;
	userId: string;
	createdAt: string;
};

export type OKResponse = {
	ok: boolean;
};

export class RbacApiError extends Error {
	status: number;

	constructor(status: number, message: string) {
		super(message);
		this.name = 'RbacApiError';
		this.status = status;
	}
}

export function isRbacApiError(error: unknown): error is RbacApiError {
	return error instanceof RbacApiError;
}

export async function getMe(fetchFn: ServerFetch, cookieHeader: string | null) {
	return rbacRequest<User>(
		fetchFn,
		'/api/me',
		{
			cookieHeader
		},
		validateUser
	);
}

export async function getMyPermissions(fetchFn: ServerFetch, cookieHeader: string | null) {
	return rbacRequest<MyPermissionsResponse>(
		fetchFn,
		'/api/me/permissions',
		{
			cookieHeader
		},
		validateMyPermissionsResponse
	);
}

export async function checkPermission(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	input: CheckPermissionInput
) {
	return rbacRequest<CheckPermissionResponse>(
		fetchFn,
		'/api/auth/check-permission',
		{
			method: 'POST',
			cookieHeader,
			body: input
		},
		validateCheckPermissionResponse
	);
}

export async function listUsers(fetchFn: ServerFetch, cookieHeader: string | null) {
	return rbacRequest<UserListResponse>(
		fetchFn,
		'/api/users',
		{
			cookieHeader
		},
		validateUserListResponse
	);
}

export async function getUser(fetchFn: ServerFetch, cookieHeader: string | null, id: string) {
	return rbacRequest<User>(
		fetchFn,
		`/api/users/${pathID(id)}`,
		{
			cookieHeader
		},
		validateUser
	);
}

export async function createUser(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	input: CreateUserInput
) {
	return rbacRequest<User>(
		fetchFn,
		'/api/users',
		{
			method: 'POST',
			cookieHeader,
			body: input
		},
		validateUser
	);
}

export async function updateUser(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	id: string,
	input: UpdateUserInput
) {
	return rbacRequest<User>(
		fetchFn,
		`/api/users/${pathID(id)}`,
		{
			method: 'PATCH',
			cookieHeader,
			body: input
		},
		validateUser
	);
}

export async function activateUser(fetchFn: ServerFetch, cookieHeader: string | null, id: string) {
	return userStatusRequest(fetchFn, cookieHeader, id, 'activate');
}

export async function deactivateUser(fetchFn: ServerFetch, cookieHeader: string | null, id: string) {
	return userStatusRequest(fetchFn, cookieHeader, id, 'deactivate');
}

export async function suspendUser(fetchFn: ServerFetch, cookieHeader: string | null, id: string) {
	return userStatusRequest(fetchFn, cookieHeader, id, 'suspend');
}

export async function deleteUser(fetchFn: ServerFetch, cookieHeader: string | null, id: string) {
	return rbacRequest<User>(
		fetchFn,
		`/api/users/${pathID(id)}`,
		{
			method: 'DELETE',
			cookieHeader
		},
		validateUser
	);
}

export async function setUserPrimaryOrganizationUnit(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	id: string,
	organizationUnitId: string | null
) {
	return rbacRequest<User>(
		fetchFn,
		`/api/users/${pathID(id)}/primary-organization-unit`,
		{
			method: 'POST',
			cookieHeader,
			body: { organizationUnitId }
		},
		validateUser
	);
}

export async function assignUserRole(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	id: string,
	input: ScopedRoleAssignmentInput
) {
	return rbacRequest<UserRoleAssignment>(
		fetchFn,
		`/api/users/${pathID(id)}/roles`,
		{
			method: 'POST',
			cookieHeader,
			body: input
		},
		validateUserRoleAssignment
	);
}

export async function removeUserRole(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	id: string,
	assignmentId: string
) {
	return rbacRequest<OKResponse>(
		fetchFn,
		`/api/users/${pathID(id)}/roles/${pathID(assignmentId)}`,
		{
			method: 'DELETE',
			cookieHeader
		},
		validateOKResponse
	);
}

export async function assignUserGroup(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	id: string,
	groupId: string
) {
	return rbacRequest<UserGroupAssignment>(
		fetchFn,
		`/api/users/${pathID(id)}/groups`,
		{
			method: 'POST',
			cookieHeader,
			body: { groupId }
		},
		validateUserGroupAssignment
	);
}

export async function removeUserGroup(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	id: string,
	groupId: string
) {
	return rbacRequest<OKResponse>(
		fetchFn,
		`/api/users/${pathID(id)}/groups/${pathID(groupId)}`,
		{
			method: 'DELETE',
			cookieHeader
		},
		validateOKResponse
	);
}

export async function listRoles(fetchFn: ServerFetch, cookieHeader: string | null) {
	return rbacRequest<RoleListResponse>(
		fetchFn,
		'/api/roles',
		{
			cookieHeader
		},
		validateRoleListResponse
	);
}

export async function getRole(fetchFn: ServerFetch, cookieHeader: string | null, id: string) {
	return rbacRequest<Role>(
		fetchFn,
		`/api/roles/${pathID(id)}`,
		{
			cookieHeader
		},
		validateRole
	);
}

export async function createRole(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	input: CreateRoleInput
) {
	return rbacRequest<Role>(
		fetchFn,
		'/api/roles',
		{
			method: 'POST',
			cookieHeader,
			body: input
		},
		validateRole
	);
}

export async function updateRole(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	id: string,
	input: UpdateRoleInput
) {
	return rbacRequest<Role>(
		fetchFn,
		`/api/roles/${pathID(id)}`,
		{
			method: 'PATCH',
			cookieHeader,
			body: input
		},
		validateRole
	);
}

export async function deleteRole(fetchFn: ServerFetch, cookieHeader: string | null, id: string) {
	return rbacRequest<OKResponse>(
		fetchFn,
		`/api/roles/${pathID(id)}`,
		{
			method: 'DELETE',
			cookieHeader
		},
		validateOKResponse
	);
}

export async function listRolePermissions(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	id: string
) {
	return rbacRequest<PermissionListResponse>(
		fetchFn,
		`/api/roles/${pathID(id)}/permissions`,
		{
			cookieHeader
		},
		validatePermissionListResponse
	);
}

export async function assignRolePermission(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	id: string,
	permissionId: string
) {
	return rbacRequest<Permission>(
		fetchFn,
		`/api/roles/${pathID(id)}/permissions`,
		{
			method: 'POST',
			cookieHeader,
			body: { permissionId }
		},
		validatePermission
	);
}

export async function removeRolePermission(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	id: string,
	permissionId: string
) {
	return rbacRequest<OKResponse>(
		fetchFn,
		`/api/roles/${pathID(id)}/permissions/${pathID(permissionId)}`,
		{
			method: 'DELETE',
			cookieHeader
		},
		validateOKResponse
	);
}

export async function listPermissions(fetchFn: ServerFetch, cookieHeader: string | null) {
	return rbacRequest<PermissionListResponse>(
		fetchFn,
		'/api/permissions',
		{
			cookieHeader
		},
		validatePermissionListResponse
	);
}

export async function listPermissionCategories(fetchFn: ServerFetch, cookieHeader: string | null) {
	return rbacRequest<PermissionCategoriesResponse>(
		fetchFn,
		'/api/permissions/categories',
		{
			cookieHeader
		},
		validatePermissionCategoriesResponse
	);
}

export async function listGroups(fetchFn: ServerFetch, cookieHeader: string | null) {
	return rbacRequest<GroupListResponse>(
		fetchFn,
		'/api/groups',
		{
			cookieHeader
		},
		validateGroupListResponse
	);
}

export async function getGroup(fetchFn: ServerFetch, cookieHeader: string | null, id: string) {
	return rbacRequest<Group>(
		fetchFn,
		`/api/groups/${pathID(id)}`,
		{
			cookieHeader
		},
		validateGroup
	);
}

export async function createGroup(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	input: CreateGroupInput
) {
	return rbacRequest<Group>(
		fetchFn,
		'/api/groups',
		{
			method: 'POST',
			cookieHeader,
			body: input
		},
		validateGroup
	);
}

export async function updateGroup(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	id: string,
	input: UpdateGroupInput
) {
	return rbacRequest<Group>(
		fetchFn,
		`/api/groups/${pathID(id)}`,
		{
			method: 'PATCH',
			cookieHeader,
			body: input
		},
		validateGroup
	);
}

export async function deleteGroup(fetchFn: ServerFetch, cookieHeader: string | null, id: string) {
	return rbacRequest<OKResponse>(
		fetchFn,
		`/api/groups/${pathID(id)}`,
		{
			method: 'DELETE',
			cookieHeader
		},
		validateOKResponse
	);
}

export async function addGroupUser(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	id: string,
	userId: string
) {
	return rbacRequest<GroupUserAssignment>(
		fetchFn,
		`/api/groups/${pathID(id)}/users`,
		{
			method: 'POST',
			cookieHeader,
			body: { userId }
		},
		validateGroupUserAssignment
	);
}

export async function removeGroupUser(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	id: string,
	userId: string
) {
	return rbacRequest<OKResponse>(
		fetchFn,
		`/api/groups/${pathID(id)}/users/${pathID(userId)}`,
		{
			method: 'DELETE',
			cookieHeader
		},
		validateOKResponse
	);
}

export async function assignGroupRole(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	id: string,
	input: ScopedRoleAssignmentInput
) {
	return rbacRequest<GroupRoleAssignment>(
		fetchFn,
		`/api/groups/${pathID(id)}/roles`,
		{
			method: 'POST',
			cookieHeader,
			body: input
		},
		validateGroupRoleAssignment
	);
}

export async function removeGroupRole(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	id: string,
	assignmentId: string
) {
	return rbacRequest<OKResponse>(
		fetchFn,
		`/api/groups/${pathID(id)}/roles/${pathID(assignmentId)}`,
		{
			method: 'DELETE',
			cookieHeader
		},
		validateOKResponse
	);
}

async function userStatusRequest(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	id: string,
	action: 'activate' | 'deactivate' | 'suspend'
) {
	return rbacRequest<User>(
		fetchFn,
		`/api/users/${pathID(id)}/${action}`,
		{
			method: 'POST',
			cookieHeader,
			body: {}
		},
		validateUser
	);
}

async function rbacRequest<T>(
	fetchFn: ServerFetch,
	path: string,
	options: {
		method?: HTTPMethod;
		body?: unknown;
		cookieHeader?: string | null;
	} = {},
	validate: (data: unknown) => T
) {
	const headers = rbacInternalHeaders();
	if (options.body !== undefined) headers.set('content-type', 'application/json');
	if (options.cookieHeader) headers.set('cookie', options.cookieHeader);

	let response: Response;
	try {
		response = await fetchFn(`${apiBaseUrl()}${path}`, {
			method: options.method ?? 'GET',
			headers,
			body: options.body === undefined ? undefined : JSON.stringify(options.body)
		});
	} catch {
		throw new RbacApiError(503, 'RBAC service unavailable');
	}

	const text = await response.text();
	const data = parseResponseJSON(text);
	if (!response.ok) {
		const message =
			data && typeof data === 'object' && 'error' in data && typeof data.error === 'string'
				? data.error
				: 'RBAC request failed';
		throw new RbacApiError(
			publicErrorStatus(response.status),
			publicErrorMessage(response.status, message, 'RBAC request failed')
		);
	}
	return validate(data);
}

function rbacInternalHeaders() {
	const headers = internalAPIHeaders();
	if (!headers) throw new RbacApiError(500, 'RBAC service is not configured');
	return headers;
}

function parseResponseJSON(text: string) {
	if (!text) return null;
	try {
		return JSON.parse(text) as unknown;
	} catch {
		return undefined;
	}
}

function validateUserListResponse(data: unknown): UserListResponse {
	if (!isRecord(data) || !Array.isArray(data.users)) invalidRbacResponse();
	return { users: data.users.map(validateUser) };
}

function validateUser(data: unknown): User {
	if (!isRecord(data)) invalidRbacResponse();
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

function validateRoleListResponse(data: unknown): RoleListResponse {
	if (!isRecord(data) || !Array.isArray(data.roles)) invalidRbacResponse();
	return { roles: data.roles.map(validateRole) };
}

function validateRole(data: unknown): Role {
	if (!isRecord(data)) invalidRbacResponse();
	const output: Role = {
		id: requiredString(data.id),
		name: requiredString(data.name),
		code: requiredString(data.code),
		isSystem: requiredBoolean(data.isSystem),
		isActive: requiredBoolean(data.isActive),
		createdAt: requiredString(data.createdAt),
		updatedAt: requiredString(data.updatedAt)
	};
	assignOptionalString(output, 'description', data.description);
	return output;
}

function validatePermissionListResponse(data: unknown): PermissionListResponse {
	if (!isRecord(data) || !Array.isArray(data.permissions)) invalidRbacResponse();
	return { permissions: data.permissions.map(validatePermission) };
}

function validatePermission(data: unknown): Permission {
	if (!isRecord(data)) invalidRbacResponse();
	const output: Permission = {
		id: requiredString(data.id),
		code: requiredString(data.code),
		name: requiredString(data.name),
		category: requiredString(data.category),
		isSystem: requiredBoolean(data.isSystem),
		createdAt: requiredString(data.createdAt),
		updatedAt: requiredString(data.updatedAt)
	};
	assignOptionalString(output, 'description', data.description);
	return output;
}

function validatePermissionCategoriesResponse(data: unknown): PermissionCategoriesResponse {
	if (!isRecord(data) || !Array.isArray(data.categories)) invalidRbacResponse();
	return { categories: data.categories.map(requiredString) };
}

function validateMyPermissionsResponse(data: unknown): MyPermissionsResponse {
	if (!isRecord(data) || !Array.isArray(data.permissions)) invalidRbacResponse();
	return { permissions: data.permissions.map(validatePermissionGrant) };
}

function validatePermissionGrant(data: unknown): PermissionGrant {
	if (!isRecord(data)) invalidRbacResponse();
	const output: PermissionGrant = {
		code: requiredString(data.code),
		scopeType: requiredScopeType(data.scopeType),
		source: requiredPermissionGrantSource(data.source)
	};
	assignOptionalString(output, 'organizationUnitId', data.organizationUnitId);
	return output;
}

function validateCheckPermissionResponse(data: unknown): CheckPermissionResponse {
	if (!isRecord(data) || typeof data.allowed !== 'boolean') invalidRbacResponse();
	return { allowed: data.allowed };
}

function validateUserRoleAssignment(data: unknown): UserRoleAssignment {
	if (!isRecord(data)) invalidRbacResponse();
	const output: UserRoleAssignment = {
		id: requiredString(data.id),
		userId: requiredString(data.userId),
		roleId: requiredString(data.roleId),
		scopeType: requiredScopeType(data.scopeType),
		createdAt: requiredString(data.createdAt),
		updatedAt: requiredString(data.updatedAt)
	};
	assignOptionalString(output, 'organizationUnitId', data.organizationUnitId);
	return output;
}

function validateUserGroupAssignment(data: unknown): UserGroupAssignment {
	if (!isRecord(data)) invalidRbacResponse();
	return {
		userId: requiredString(data.userId),
		groupId: requiredString(data.groupId),
		createdAt: requiredString(data.createdAt)
	};
}

function validateGroupListResponse(data: unknown): GroupListResponse {
	if (!isRecord(data) || !Array.isArray(data.groups)) invalidRbacResponse();
	return { groups: data.groups.map(validateGroup) };
}

function validateGroup(data: unknown): Group {
	if (!isRecord(data)) invalidRbacResponse();
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

function validateGroupRoleAssignment(data: unknown): GroupRoleAssignment {
	if (!isRecord(data)) invalidRbacResponse();
	const output: GroupRoleAssignment = {
		id: requiredString(data.id),
		groupId: requiredString(data.groupId),
		roleId: requiredString(data.roleId),
		scopeType: requiredScopeType(data.scopeType),
		createdAt: requiredString(data.createdAt)
	};
	assignOptionalString(output, 'organizationUnitId', data.organizationUnitId);
	return output;
}

function validateGroupUserAssignment(data: unknown): GroupUserAssignment {
	if (!isRecord(data)) invalidRbacResponse();
	return {
		groupId: requiredString(data.groupId),
		userId: requiredString(data.userId),
		createdAt: requiredString(data.createdAt)
	};
}

function validateOKResponse(data: unknown): OKResponse {
	if (!isRecord(data) || typeof data.ok !== 'boolean') invalidRbacResponse();
	return { ok: data.ok };
}

function requiredString(value: unknown) {
	if (typeof value !== 'string' || value.trim() === '') invalidRbacResponse();
	return value;
}

function requiredBoolean(value: unknown) {
	if (typeof value !== 'boolean') invalidRbacResponse();
	return value;
}

function optionalString(value: unknown) {
	if (value === undefined || value === null) return value;
	if (typeof value === 'string') return value;
	invalidRbacResponse();
}

function assignOptionalString<T extends Record<string, unknown>, K extends keyof T>(
	output: T,
	key: K,
	value: unknown
) {
	const parsed = optionalString(value);
	if (parsed !== undefined) output[key] = parsed as T[K];
}

function requiredPreferredLanguage(value: unknown): PreferredLanguage {
	if (value === 'en' || value === 'ro') return value;
	invalidRbacResponse();
}

function requiredUserRole(value: unknown): UserRole {
	if (value === 'admin' || value === 'user') return value;
	invalidRbacResponse();
}

function requiredUserStatus(value: unknown): UserStatus {
	if (
		value === 'invited' ||
		value === 'active' ||
		value === 'inactive' ||
		value === 'suspended' ||
		value === 'deleted'
	) {
		return value;
	}
	invalidRbacResponse();
}

function requiredScopeType(value: unknown): ScopeType {
	if (
		value === 'global' ||
		value === 'organization_unit' ||
		value === 'organization_unit_and_children'
	) {
		return value;
	}
	invalidRbacResponse();
}

function requiredPermissionGrantSource(value: unknown): PermissionGrantSource {
	if (value === 'user_role' || value === 'group_role' || value === 'organization_unit_role') {
		return value;
	}
	invalidRbacResponse();
}

function isRecord(value: unknown): value is Record<string, unknown> {
	return typeof value === 'object' && value !== null;
}

function pathID(id: string) {
	return encodeURIComponent(id);
}

function invalidRbacResponse(): never {
	throw new RbacApiError(502, 'Invalid RBAC response');
}
