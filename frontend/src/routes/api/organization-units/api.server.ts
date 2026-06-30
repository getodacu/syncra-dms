import { json } from '@sveltejs/kit';
import {
	isOrganizationUnitApiError,
	type OrganizationUnitInput
} from '$lib/server/organization-units';
import { publicErrorMessage, publicErrorStatus } from '$lib/server/public-errors';

export function requireAuthenticatedUser(locals: App.Locals) {
	if (locals.user) return null;
	return jsonError(401, 'Authentication required');
}

export function requireAdminUser(locals: App.Locals) {
	if (!locals.user) return jsonError(401, 'Authentication required');
	if (locals.user.role !== 'admin') return jsonError(403, 'admin role required');
	return null;
}

export function cookieHeader(request: Request) {
	return request.headers.get('cookie');
}

export async function readOrganizationUnitInput(
	request: Request
): Promise<OrganizationUnitInput | Response> {
	const data = await readJSONRecord(request);
	if (data instanceof Response) return data;

	return {
		parentId: parentIdValue(data.parentId),
		name: textValue(data.name),
		code: textValue(data.code),
		description: textValue(data.description)
	};
}

export async function readMoveParentId(request: Request): Promise<string | null | Response> {
	const data = await readJSONRecord(request);
	if (data instanceof Response) return data;

	if (!('parentId' in data)) {
		return jsonError(400, 'parentId is required');
	}

	return parentIdValue(data.parentId);
}

export function organizationUnitAPIErrorResponse(error: unknown, fallback: string) {
	if (isOrganizationUnitApiError(error)) {
		return jsonError(
			publicErrorStatus(error.status),
			publicErrorMessage(error.status, error.message, fallback)
		);
	}

	throw error;
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

function parentIdValue(value: unknown) {
	if (value === null || value === undefined) return null;
	if (typeof value !== 'string') return null;

	const trimmed = value.trim();
	return trimmed ? trimmed : null;
}

function isRecord(value: unknown): value is Record<string, unknown> {
	return typeof value === 'object' && value !== null && !Array.isArray(value);
}
