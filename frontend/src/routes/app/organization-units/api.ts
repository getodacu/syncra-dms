import { publicApiErrorMessage } from '$lib/client/api-errors';
import type { OrganizationUnitNode } from './tree';

type ClientFetch = typeof fetch;
type JsonObject = Record<string, unknown>;

export const ORGANIZATION_UNITS_QUERY_KEY = ['organization-units', 'tree'] as const;

export type OrganizationUnitListResponse = {
	units: OrganizationUnitNode[];
};

export type OrganizationUnitInput = {
	parentId?: string | null;
	name: string;
	code?: string | null;
	description?: string | null;
};

export type UpdateOrganizationUnitVariables = {
	id: string;
	input: OrganizationUnitInput;
};

export type MoveOrganizationUnitVariables = {
	id: string;
	parentId: string | null;
};

export type ArchiveOrganizationUnitVariables = {
	id: string;
};

export type ArchiveOrganizationUnitResponse = {
	ok: boolean;
};

export async function fetchOrganizationUnitTree(
	fetchFn: ClientFetch
): Promise<OrganizationUnitListResponse> {
	return organizationUnitRequest(
		fetchFn,
		'/api/organization-units/tree',
		{ method: 'GET' },
		validateOrganizationUnitListResponse,
		'Failed to load organization units'
	);
}

export async function createOrganizationUnit(fetchFn: ClientFetch, input: OrganizationUnitInput) {
	return organizationUnitRequest(
		fetchFn,
		'/api/organization-units',
		{
			method: 'POST',
			body: organizationUnitPayload(input)
		},
		validateOrganizationUnitNode,
		'Failed to create organization unit'
	);
}

export async function updateOrganizationUnit(
	fetchFn: ClientFetch,
	{ id, input }: UpdateOrganizationUnitVariables
) {
	return organizationUnitRequest(
		fetchFn,
		`/api/organization-units/${encodeURIComponent(id)}`,
		{
			method: 'PATCH',
			body: organizationUnitPayload(input)
		},
		validateOrganizationUnitNode,
		'Failed to update organization unit'
	);
}

export async function moveOrganizationUnit(
	fetchFn: ClientFetch,
	{ id, parentId }: MoveOrganizationUnitVariables
) {
	return organizationUnitRequest(
		fetchFn,
		`/api/organization-units/${encodeURIComponent(id)}/parent`,
		{
			method: 'PATCH',
			body: { parentId }
		},
		validateOrganizationUnitNode,
		'Failed to move organization unit'
	);
}

export async function archiveOrganizationUnit(
	fetchFn: ClientFetch,
	{ id }: ArchiveOrganizationUnitVariables
) {
	return organizationUnitRequest(
		fetchFn,
		`/api/organization-units/${encodeURIComponent(id)}/archive`,
		{
			method: 'POST',
			body: {}
		},
		validateArchiveOrganizationUnitResponse,
		'Failed to archive organization unit'
	);
}

export function isOrganizationUnitListResponse(
	value: unknown
): value is OrganizationUnitListResponse {
	try {
		validateOrganizationUnitListResponse(value);
		return true;
	} catch {
		return false;
	}
}

async function organizationUnitRequest<T>(
	fetchFn: ClientFetch,
	path: string,
	options: {
		method: 'GET' | 'PATCH' | 'POST';
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
	const json = await readResponseJSON(response);

	if (!response.ok) {
		throw new Error(publicApiErrorMessage(response.status, json, fallback));
	}

	return validate(json);
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

function organizationUnitPayload(input: OrganizationUnitInput) {
	return {
		parentId: normalizeParentId(input.parentId),
		name: input.name,
		code: input.code ?? '',
		description: input.description ?? ''
	};
}

function normalizeParentId(parentId: string | null | undefined) {
	if (!parentId) return null;
	const trimmed = parentId.trim();
	return trimmed ? trimmed : null;
}

function validateOrganizationUnitListResponse(data: unknown): OrganizationUnitListResponse {
	if (!isJsonObject(data) || !Array.isArray(data.units)) invalidOrganizationUnitResponse();
	return {
		units: data.units.map(validateOrganizationUnitNode)
	};
}

function validateOrganizationUnitNode(data: unknown): OrganizationUnitNode {
	if (!isJsonObject(data)) invalidOrganizationUnitResponse();
	const id = requiredString(data.id);
	const name = requiredString(data.name);
	const createdAt = requiredString(data.createdAt);
	const updatedAt = requiredString(data.updatedAt);
	const children = Array.isArray(data.children)
		? data.children.map(validateOrganizationUnitNode)
		: [];
	const output: OrganizationUnitNode = {
		id,
		name,
		createdAt,
		updatedAt,
		children
	};
	const parentId = optionalString(data.parentId);
	const code = optionalString(data.code);
	const description = optionalString(data.description);
	const archivedAt = optionalString(data.archivedAt);

	if (parentId !== undefined) output.parentId = parentId;
	if (code !== undefined) output.code = code;
	if (description !== undefined) output.description = description;
	if (archivedAt !== undefined) output.archivedAt = archivedAt;

	return output;
}

function validateArchiveOrganizationUnitResponse(data: unknown): ArchiveOrganizationUnitResponse {
	if (!isJsonObject(data) || typeof data.ok !== 'boolean') invalidOrganizationUnitResponse();
	return { ok: data.ok };
}

function requiredString(value: unknown) {
	if (typeof value !== 'string' || value.trim() === '') invalidOrganizationUnitResponse();
	return value;
}

function optionalString(value: unknown) {
	if (value === undefined || value === null) return value;
	if (typeof value === 'string') return value;
	invalidOrganizationUnitResponse();
}

function isJsonObject(value: unknown): value is JsonObject {
	return typeof value === 'object' && value !== null && !Array.isArray(value);
}

function invalidOrganizationUnitResponse(): never {
	throw new Error('Invalid organization unit response');
}
