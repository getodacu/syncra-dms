import { json } from '@sveltejs/kit';
import {
	isDocumentApiError,
	type DocumentFolderInput,
	type DocumentUpdateInput
} from '$lib/server/documents';
import { publicErrorMessage, publicErrorStatus } from '$lib/server/public-errors';

export function requireAuthenticatedUser(locals: App.Locals) {
	if (locals.user) return null;
	return jsonError(401, 'Authentication required');
}

export function hasAnyPermission(permissions: string[], required: string[]) {
	return (
		permissions.includes('system.admin') ||
		required.some((permission) => permissions.includes(permission))
	);
}

export function cookieHeader(request: Request) {
	return request.headers.get('cookie');
}

export async function readDocumentFolderInput(
	request: Request
): Promise<DocumentFolderInput | Response> {
	const data = await readJSONRecord(request);
	if (data instanceof Response) return data;

	return {
		organizationUnitId: textValue(data.organizationUnitId),
		parentId: parentIdValue(data.parentId),
		name: textValue(data.name),
		description: nullableTextValue(data.description)
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

export async function readDocumentUpdateInput(
	request: Request
): Promise<DocumentUpdateInput | Response> {
	const data = await readJSONRecord(request);
	if (data instanceof Response) return data;

	return {
		displayName: textValue(data.displayName)
	};
}

export function documentAPIErrorResponse(error: unknown, fallback: string) {
	if (isDocumentApiError(error)) {
		return jsonError(
			publicErrorStatus(error.status),
			publicErrorMessage(error.status, error.message, fallback)
		);
	}

	throw error;
}

export function jsonError(status: number, message: string) {
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

function nullableTextValue(value: unknown) {
	return typeof value === 'string' ? value : null;
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
