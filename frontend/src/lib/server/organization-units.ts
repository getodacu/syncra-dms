import { apiBaseUrl, internalAPIHeaders } from './internal-api';
import { publicErrorMessage, publicErrorStatus } from './public-errors';

type ServerFetch = typeof fetch;

export type OrganizationUnit = {
	id: string;
	parentId?: string | null;
	name: string;
	code?: string | null;
	description?: string | null;
	archivedAt?: string | null;
	createdAt: string;
	updatedAt: string;
	children?: OrganizationUnit[];
};

export type OrganizationUnitTreeNode = Omit<OrganizationUnit, 'children'> & {
	children: OrganizationUnitTreeNode[];
};

export type OrganizationUnitListResponse = {
	units: OrganizationUnitTreeNode[];
};

export type OrganizationUnitInput = {
	parentId?: string | null;
	name: string;
	code: string;
	description: string;
};

export class OrganizationUnitApiError extends Error {
	status: number;

	constructor(status: number, message: string) {
		super(message);
		this.name = 'OrganizationUnitApiError';
		this.status = status;
	}
}

export function isOrganizationUnitApiError(error: unknown): error is OrganizationUnitApiError {
	return error instanceof OrganizationUnitApiError;
}

async function organizationUnitRequest<T>(
	fetchFn: ServerFetch,
	path: string,
	options: {
		method?: 'GET' | 'POST' | 'PATCH';
		body?: unknown;
		cookieHeader?: string | null;
	} = {},
	validate: (data: unknown) => T
) {
	const headers = organizationUnitInternalHeaders();
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
		throw new OrganizationUnitApiError(503, 'Organization Unit service unavailable');
	}

	const text = await response.text();
	const data = parseResponseJSON(text);
	if (!response.ok) {
		const message =
			data && typeof data === 'object' && 'error' in data && typeof data.error === 'string'
				? data.error
				: 'Organization Unit request failed';
		throw new OrganizationUnitApiError(
			publicErrorStatus(response.status),
			publicErrorMessage(response.status, message, 'Organization Unit request failed')
		);
	}
	return validate(data);
}

function organizationUnitInternalHeaders() {
	const headers = internalAPIHeaders();
	if (!headers) throw new OrganizationUnitApiError(500, 'Organization Unit service is not configured');
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

export async function getOrganizationUnitTree(fetchFn: ServerFetch, cookieHeader: string | null) {
	return organizationUnitRequest<OrganizationUnitListResponse>(
		fetchFn,
		'/api/organization-units/tree',
		{
			cookieHeader
		},
		validateOrganizationUnitListResponse
	);
}

export async function createOrganizationUnit(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	input: OrganizationUnitInput
) {
	return organizationUnitRequest<OrganizationUnit>(
		fetchFn,
		'/api/organization-units',
		{
			method: 'POST',
			cookieHeader,
			body: input
		},
		validateOrganizationUnit
	);
}

export async function updateOrganizationUnit(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	id: string,
	input: OrganizationUnitInput
) {
	return organizationUnitRequest<OrganizationUnit>(
		fetchFn,
		`/api/organization-units/${encodeURIComponent(id)}`,
		{
			method: 'PATCH',
			cookieHeader,
			body: input
		},
		validateOrganizationUnit
	);
}

export async function moveOrganizationUnit(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	id: string,
	parentId: string | null
) {
	return organizationUnitRequest<OrganizationUnit>(
		fetchFn,
		`/api/organization-units/${encodeURIComponent(id)}/parent`,
		{
			method: 'PATCH',
			cookieHeader,
			body: { parentId }
		},
		validateOrganizationUnit
	);
}

export async function archiveOrganizationUnit(
	fetchFn: ServerFetch,
	cookieHeader: string | null,
	id: string
) {
	return organizationUnitRequest<{ ok: boolean }>(
		fetchFn,
		`/api/organization-units/${encodeURIComponent(id)}/archive`,
		{
			method: 'POST',
			cookieHeader,
			body: {}
		},
		validateOKResponse
	);
}

function validateOrganizationUnitListResponse(data: unknown): OrganizationUnitListResponse {
	if (!isRecord(data) || !Array.isArray(data.units)) invalidOrganizationUnitResponse();
	return {
		units: data.units.map(validateOrganizationUnitTreeNode)
	};
}

function validateOrganizationUnitTreeNode(data: unknown): OrganizationUnitTreeNode {
	const unit = validateOrganizationUnit(data);
	if (!isRecord(data) || !Array.isArray(data.children)) invalidOrganizationUnitResponse();
	return {
		...unit,
		children: data.children.map(validateOrganizationUnitTreeNode)
	};
}

function validateOrganizationUnit(data: unknown): OrganizationUnit {
	if (!isRecord(data)) invalidOrganizationUnitResponse();
	const id = requiredString(data.id);
	const name = requiredString(data.name);
	const createdAt = requiredString(data.createdAt);
	const updatedAt = requiredString(data.updatedAt);
	const parentId = optionalString(data.parentId);
	const code = optionalString(data.code);
	const description = optionalString(data.description);
	const archivedAt = optionalString(data.archivedAt);
	const output: OrganizationUnit = {
		id,
		name,
		createdAt,
		updatedAt
	};
	if (parentId !== undefined) output.parentId = parentId;
	if (code !== undefined) output.code = code;
	if (description !== undefined) output.description = description;
	if (archivedAt !== undefined) output.archivedAt = archivedAt;
	if ('children' in data) {
		if (!Array.isArray(data.children)) invalidOrganizationUnitResponse();
		output.children = data.children.map(validateOrganizationUnit);
	}
	return output;
}

function validateOKResponse(data: unknown): { ok: boolean } {
	if (!isRecord(data) || typeof data.ok !== 'boolean') invalidOrganizationUnitResponse();
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

function isRecord(value: unknown): value is Record<string, unknown> {
	return typeof value === 'object' && value !== null;
}

function invalidOrganizationUnitResponse(): never {
	throw new OrganizationUnitApiError(502, 'Invalid Organization Unit response');
}
