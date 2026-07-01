import { json } from '@sveltejs/kit';
import { isRbacApiError } from '$lib/server/rbac';
import type {
	CreateGroupInput,
	CreateRoleInput,
	CreateUserInput,
	ScopeType,
	ScopedRoleAssignmentInput,
	UpdateGroupInput,
	UpdateRoleInput,
	UpdateUserInput,
	UserStatus
} from '$lib/server/rbac';
import { publicErrorMessage, publicErrorStatus } from '$lib/server/public-errors';

export function requireAuthenticatedUser(locals: App.Locals) {
	if (locals.user) return null;
	return jsonError(401, 'Authentication required');
}

export function requireLocalPermission(locals: App.Locals, permission: string) {
	const authError = requireAuthenticatedUser(locals);
	if (authError) return authError;

	const permissions = locals.permissions ?? locals.user?.permissions ?? [];
	if (!permissions.includes(permission) && !permissions.includes('system.admin')) {
		return jsonError(403, 'permission required');
	}

	return null;
}

export function cookieHeader(request: Request) {
	return request.headers.get('cookie');
}

export async function authenticatedRbacJSON<T>(
	locals: App.Locals,
	fallback: string,
	operation: () => Promise<T>,
	status = 200
) {
	const authError = requireAuthenticatedUser(locals);
	if (authError) return authError;

	try {
		return json(await operation(), { status });
	} catch (error) {
		return rbacAPIErrorResponse(error, fallback);
	}
}

export async function permittedRbacJSON<T>(
	locals: App.Locals,
	permission: string,
	fallback: string,
	operation: () => Promise<T>,
	status = 200
) {
	const authError = requireLocalPermission(locals, permission);
	if (authError) return authError;

	try {
		return json(await operation(), { status });
	} catch (error) {
		return rbacAPIErrorResponse(error, fallback);
	}
}

export function rbacAPIErrorResponse(error: unknown, fallback: string) {
	if (isRbacApiError(error)) {
		return jsonError(
			publicErrorStatus(error.status),
			publicErrorMessage(error.status, error.message, fallback)
		);
	}

	throw error;
}

export async function readCreateUserInput(request: Request): Promise<CreateUserInput | Response> {
	const data = await readJSONRecord(request);
	if (data instanceof Response) return data;

	const status = optionalUserStatus(data.status);
	if (status instanceof Response) return status;

	return {
		name: textValue(data.name),
		email: textValue(data.email),
		status,
		primaryOrganizationUnitId: optionalNullableString(data.primaryOrganizationUnitId),
		managerUserId: optionalNullableString(data.managerUserId),
		jobTitle: optionalNullableString(data.jobTitle),
		phone: optionalNullableString(data.phone)
	};
}

export async function readUpdateUserInput(request: Request): Promise<UpdateUserInput | Response> {
	const data = await readJSONRecord(request);
	if (data instanceof Response) return data;

	const input: UpdateUserInput = {};
	assignOptionalText(input, 'name', data);
	assignOptionalNullableText(input, 'managerUserId', data);
	assignOptionalNullableText(input, 'jobTitle', data);
	assignOptionalNullableText(input, 'phone', data);
	return input;
}

export async function readPrimaryOrganizationUnitInput(
	request: Request
): Promise<string | null | Response> {
	const data = await readJSONRecord(request);
	if (data instanceof Response) return data;

	if (!('organizationUnitId' in data)) {
		return jsonError(400, 'organizationUnitId is required');
	}

	return nullableTrimmedString(data.organizationUnitId);
}

export async function readScopedRoleAssignmentInput(
	request: Request
): Promise<ScopedRoleAssignmentInput | Response> {
	const data = await readJSONRecord(request);
	if (data instanceof Response) return data;

	const roleId = requiredText(data, 'roleId');
	if (roleId instanceof Response) return roleId;

	const scopeType = scopeTypeValue(data.scopeType ?? 'global');
	if (scopeType instanceof Response) return scopeType;

	return {
		roleId,
		scopeType,
		organizationUnitId: optionalNullableString(data.organizationUnitId)
	};
}

export async function readGroupIdInput(request: Request): Promise<string | Response> {
	const data = await readJSONRecord(request);
	if (data instanceof Response) return data;
	return requiredText(data, 'groupId');
}

export async function readUserIdInput(request: Request): Promise<string | Response> {
	const data = await readJSONRecord(request);
	if (data instanceof Response) return data;
	return requiredText(data, 'userId');
}

export async function readPermissionIdInput(request: Request): Promise<string | Response> {
	const data = await readJSONRecord(request);
	if (data instanceof Response) return data;
	return requiredText(data, 'permissionId');
}

export async function readCreateRoleInput(request: Request): Promise<CreateRoleInput | Response> {
	const data = await readJSONRecord(request);
	if (data instanceof Response) return data;

	return {
		name: textValue(data.name),
		code: textValue(data.code),
		description: optionalNullableString(data.description),
		isActive: optionalBoolean(data.isActive)
	};
}

export async function readUpdateRoleInput(request: Request): Promise<UpdateRoleInput | Response> {
	const data = await readJSONRecord(request);
	if (data instanceof Response) return data;

	const input: UpdateRoleInput = {};
	assignOptionalText(input, 'name', data);
	assignOptionalText(input, 'code', data);
	assignOptionalNullableText(input, 'description', data);
	assignOptionalBoolean(input, 'isActive', data);
	return input;
}

export async function readCreateGroupInput(request: Request): Promise<CreateGroupInput | Response> {
	const data = await readJSONRecord(request);
	if (data instanceof Response) return data;

	return {
		name: textValue(data.name),
		code: textValue(data.code),
		description: optionalNullableString(data.description),
		organizationUnitId: optionalNullableString(data.organizationUnitId),
		isActive: optionalBoolean(data.isActive)
	};
}

export async function readUpdateGroupInput(request: Request): Promise<UpdateGroupInput | Response> {
	const data = await readJSONRecord(request);
	if (data instanceof Response) return data;

	const input: UpdateGroupInput = {};
	assignOptionalText(input, 'name', data);
	assignOptionalText(input, 'code', data);
	assignOptionalNullableText(input, 'description', data);
	assignOptionalNullableText(input, 'organizationUnitId', data);
	assignOptionalBoolean(input, 'isActive', data);
	return input;
}

function jsonError(status: number, message: string) {
	return json({ error: message }, { status });
}

async function readJSONRecord(request: Request): Promise<Record<string, unknown> | Response> {
	let data: unknown;
	try {
		data = await request.json();
	} catch {
		return jsonError(400, 'invalid JSON body');
	}

	if (!isRecord(data)) {
		return jsonError(400, 'invalid JSON body');
	}

	return data;
}

function textValue(value: unknown) {
	return typeof value === 'string' ? value : '';
}

function optionalNullableString(value: unknown) {
	if (value === null) return null;
	if (typeof value === 'string') return value;
	return undefined;
}

function nullableTrimmedString(value: unknown) {
	if (value === null || value === undefined) return null;
	if (typeof value !== 'string') return null;

	const trimmed = value.trim();
	return trimmed ? trimmed : null;
}

function optionalBoolean(value: unknown) {
	return typeof value === 'boolean' ? value : undefined;
}

function optionalUserStatus(value: unknown): UserStatus | undefined | Response {
	if (value === undefined) return undefined;
	if (
		value === 'invited' ||
		value === 'active' ||
		value === 'inactive' ||
		value === 'suspended' ||
		value === 'deleted'
	) {
		return value;
	}
	return jsonError(400, 'status is invalid');
}

function scopeTypeValue(value: unknown): ScopeType | Response {
	if (
		value === 'global' ||
		value === 'organization_unit' ||
		value === 'organization_unit_and_children'
	) {
		return value;
	}
	return jsonError(400, 'scopeType is invalid');
}

function requiredText(data: Record<string, unknown>, field: string) {
	const value = data[field];
	if (typeof value === 'string' && value.trim() !== '') return value;
	return jsonError(400, `${field} is required`);
}

function assignOptionalText<T extends object>(output: T, field: keyof T & string, data: Record<string, unknown>) {
	if (!(field in data)) return;
	(output as Record<string, unknown>)[field] = textValue(data[field]);
}

function assignOptionalNullableText<T extends object>(
	output: T,
	field: keyof T & string,
	data: Record<string, unknown>
) {
	if (!(field in data)) return;
	(output as Record<string, unknown>)[field] = optionalNullableString(data[field]);
}

function assignOptionalBoolean<T extends object>(output: T, field: keyof T & string, data: Record<string, unknown>) {
	if (!(field in data)) return;
	(output as Record<string, unknown>)[field] = optionalBoolean(data[field]);
}

function isRecord(value: unknown): value is Record<string, unknown> {
	return typeof value === 'object' && value !== null && !Array.isArray(value);
}
