import { publicApiErrorMessage } from '$lib/client/api-errors';

type ClientFetch = typeof fetch;
type JsonObject = Record<string, unknown>;

export const ROLES_QUERY_KEY = ['admin-roles'] as const;
export const PERMISSIONS_QUERY_KEY = ['admin-permissions'] as const;

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

export type UpdateRoleVariables = {
	id: string;
	input: UpdateRoleInput;
};

export type DeleteRoleVariables = {
	id: string;
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

export type RolePermissionVariables = {
	id: string;
	permissionId: string;
};

export type GroupedPermissions = {
	category: string;
	permissions: Permission[];
};

export async function fetchRoles(fetchFn: ClientFetch): Promise<RoleListResponse> {
	return roleRequest(
		fetchFn,
		'/api/roles',
		{ method: 'GET' },
		validateRoleListResponse,
		'Failed to load roles'
	);
}

export async function createRole(fetchFn: ClientFetch, input: CreateRoleInput) {
	return roleRequest(
		fetchFn,
		'/api/roles',
		{ method: 'POST', body: input },
		validateRole,
		'Failed to create role'
	);
}

export async function updateRole(fetchFn: ClientFetch, { id, input }: UpdateRoleVariables) {
	return roleRequest(
		fetchFn,
		`/api/roles/${encodeURIComponent(id)}`,
		{ method: 'PATCH', body: input },
		validateRole,
		'Failed to update role'
	);
}

export async function deleteRole(fetchFn: ClientFetch, { id }: DeleteRoleVariables) {
	return roleRequest(
		fetchFn,
		`/api/roles/${encodeURIComponent(id)}`,
		{ method: 'DELETE' },
		validateOKResponse,
		'Failed to delete role'
	);
}

export async function fetchPermissions(fetchFn: ClientFetch): Promise<PermissionListResponse> {
	return roleRequest(
		fetchFn,
		'/api/permissions',
		{ method: 'GET' },
		validatePermissionListResponse,
		'Failed to load permissions'
	);
}

export async function fetchPermissionCategories(
	fetchFn: ClientFetch
): Promise<PermissionCategoriesResponse> {
	return roleRequest(
		fetchFn,
		'/api/permissions/categories',
		{ method: 'GET' },
		validatePermissionCategoriesResponse,
		'Failed to load permission categories'
	);
}

export async function fetchRolePermissions(
	fetchFn: ClientFetch,
	{ id }: Pick<RolePermissionVariables, 'id'>
) {
	return roleRequest(
		fetchFn,
		`/api/roles/${encodeURIComponent(id)}/permissions`,
		{ method: 'GET' },
		validatePermissionListResponse,
		'Failed to load role permissions'
	);
}

export async function assignRolePermission(
	fetchFn: ClientFetch,
	{ id, permissionId }: RolePermissionVariables
) {
	return roleRequest(
		fetchFn,
		`/api/roles/${encodeURIComponent(id)}/permissions`,
		{ method: 'POST', body: { permissionId } },
		validatePermission,
		'Failed to assign role permission'
	);
}

export async function removeRolePermission(
	fetchFn: ClientFetch,
	{ id, permissionId }: RolePermissionVariables
) {
	return roleRequest(
		fetchFn,
		`/api/roles/${encodeURIComponent(id)}/permissions/${encodeURIComponent(permissionId)}`,
		{ method: 'DELETE' },
		validateOKResponse,
		'Failed to remove role permission'
	);
}

export function groupPermissionsByCategory(permissions: Permission[]): GroupedPermissions[] {
	const groups = new Map<string, Permission[]>();
	for (const permission of permissions) {
		const group = groups.get(permission.category) ?? [];
		group.push(permission);
		groups.set(permission.category, group);
	}

	return [...groups.entries()]
		.sort(([left], [right]) => left.localeCompare(right))
		.map(([category, groupedPermissions]) => ({
			category,
			permissions: [...groupedPermissions].sort((left, right) => left.code.localeCompare(right.code))
		}));
}

async function roleRequest<T>(
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

function validateRoleListResponse(data: unknown): RoleListResponse {
	if (!isJsonObject(data) || !Array.isArray(data.roles)) invalidRoleResponse();
	return { roles: data.roles.map(validateRole) };
}

function validateRole(data: unknown): Role {
	if (!isJsonObject(data)) invalidRoleResponse();
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
	if (!isJsonObject(data) || !Array.isArray(data.permissions)) invalidRoleResponse();
	return { permissions: data.permissions.map(validatePermission) };
}

function validatePermission(data: unknown): Permission {
	if (!isJsonObject(data)) invalidRoleResponse();
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
	if (!isJsonObject(data) || !Array.isArray(data.categories)) invalidRoleResponse();
	return { categories: data.categories.map(requiredString) };
}

function validateOKResponse(data: unknown) {
	if (!isJsonObject(data) || typeof data.ok !== 'boolean') invalidRoleResponse();
	return { ok: data.ok };
}

function requiredString(value: unknown) {
	if (typeof value !== 'string' || value.trim() === '') invalidRoleResponse();
	return value;
}

function requiredBoolean(value: unknown) {
	if (typeof value !== 'boolean') invalidRoleResponse();
	return value;
}

function optionalString(value: unknown) {
	if (value === undefined || value === null) return value;
	if (typeof value === 'string') return value;
	invalidRoleResponse();
}

function assignOptionalString<T extends object>(output: T, key: keyof T & string, value: unknown) {
	const parsed = optionalString(value);
	if (parsed !== undefined) (output as Record<string, unknown>)[key] = parsed;
}

function isJsonObject(value: unknown): value is JsonObject {
	return typeof value === 'object' && value !== null && !Array.isArray(value);
}

function invalidRoleResponse(): never {
	throw new Error('Invalid role response');
}
