import { beforeEach, describe, expect, it, vi } from 'vitest';

const organizationUnitMocks = vi.hoisted(() => {
	class MockOrganizationUnitApiError extends Error {
		status: number;

		constructor(status: number, message: string) {
			super(message);
			this.name = 'OrganizationUnitApiError';
			this.status = status;
		}
	}

	return {
		OrganizationUnitApiError: MockOrganizationUnitApiError,
		archiveOrganizationUnit: vi.fn(),
		createOrganizationUnit: vi.fn(),
		getOrganizationUnitTree: vi.fn(),
		moveOrganizationUnit: vi.fn(),
		updateOrganizationUnit: vi.fn()
	};
});

vi.mock('$lib/server/organization-units', () => ({
	OrganizationUnitApiError: organizationUnitMocks.OrganizationUnitApiError,
	archiveOrganizationUnit: organizationUnitMocks.archiveOrganizationUnit,
	createOrganizationUnit: organizationUnitMocks.createOrganizationUnit,
	getOrganizationUnitTree: organizationUnitMocks.getOrganizationUnitTree,
	isOrganizationUnitApiError: (error: unknown) =>
		error instanceof organizationUnitMocks.OrganizationUnitApiError,
	moveOrganizationUnit: organizationUnitMocks.moveOrganizationUnit,
	updateOrganizationUnit: organizationUnitMocks.updateOrganizationUnit
}));

import { POST as createOrganizationUnitRoute } from './+server';
import { PATCH as updateOrganizationUnitRoute } from './[id]/+server';
import { POST as archiveOrganizationUnitRoute } from './[id]/archive/+server';
import { PATCH as moveOrganizationUnitRoute } from './[id]/parent/+server';
import { GET as getOrganizationUnitTreeRoute } from './tree/+server';

const cookieHeader = 'auth.session_token=token';
const unit = {
	id: 'unit-id',
	name: 'Operations',
	code: 'OPS',
	description: 'Handles operations',
	createdAt: '2026-06-30T00:00:00Z',
	updatedAt: '2026-06-30T00:00:00Z',
	children: []
};

describe('organization unit Svelte API routes', () => {
	beforeEach(() => {
		vi.clearAllMocks();
	});

	it('requires authenticated users for tree reads', async () => {
		const response = await getOrganizationUnitTreeRoute(
			event({ locals: { user: null }, path: '/api/organization-units/tree' }) as never
		);

		expect(response.status).toBe(401);
		expect(await response.json()).toEqual({ error: 'Authentication required' });
		expect(organizationUnitMocks.getOrganizationUnitTree).not.toHaveBeenCalled();
	});

	it('proxies tree reads with cookie context', async () => {
		organizationUnitMocks.getOrganizationUnitTree.mockResolvedValue({ units: [unit] });
		const fetchMock = vi.fn();

		const response = await getOrganizationUnitTreeRoute(
			event({
				fetch: fetchMock,
				locals: { user: { role: 'user' } },
				path: '/api/organization-units/tree'
			}) as never
		);

		expect(response.status).toBe(200);
		expect(await response.json()).toEqual({ units: [unit] });
		expect(organizationUnitMocks.getOrganizationUnitTree).toHaveBeenCalledWith(
			fetchMock,
			cookieHeader
		);
	});

	it('rejects non-admin mutations before proxying', async () => {
		const response = await createOrganizationUnitRoute(
			event({
				locals: { user: { role: 'user' } },
				method: 'POST',
				path: '/api/organization-units',
				body: { name: 'Operations' }
			}) as never
		);

		expect(response.status).toBe(403);
		expect(await response.json()).toEqual({ error: 'admin role required' });
		expect(organizationUnitMocks.createOrganizationUnit).not.toHaveBeenCalled();
	});

	it('proxies admin create requests with normalized input and cookie context', async () => {
		organizationUnitMocks.createOrganizationUnit.mockResolvedValue(unit);
		const fetchMock = vi.fn();

		const response = await createOrganizationUnitRoute(
			event({
				fetch: fetchMock,
				locals: { user: { role: 'admin' } },
				method: 'POST',
				path: '/api/organization-units',
				body: {
					parentId: '   ',
					name: 'Operations',
					code: 'OPS',
					description: 'Handles operations'
				}
			}) as never
		);

		expect(response.status).toBe(201);
		expect(await response.json()).toEqual(unit);
		expect(organizationUnitMocks.createOrganizationUnit).toHaveBeenCalledWith(
			fetchMock,
			cookieHeader,
			{
				parentId: null,
				name: 'Operations',
				code: 'OPS',
				description: 'Handles operations'
			}
		);
	});

	it('proxies admin update, move, and archive requests', async () => {
		organizationUnitMocks.updateOrganizationUnit.mockResolvedValue(unit);
		organizationUnitMocks.moveOrganizationUnit.mockResolvedValue(unit);
		organizationUnitMocks.archiveOrganizationUnit.mockResolvedValue({ ok: true });
		const fetchMock = vi.fn();

		await updateOrganizationUnitRoute(
			event({
				fetch: fetchMock,
				locals: { user: { role: 'admin' } },
				method: 'PATCH',
				path: '/api/organization-units/unit-id',
				params: { id: 'unit-id' },
				body: {
					parentId: 'parent-id',
					name: 'Operations',
					code: '',
					description: ''
				}
			}) as never
		);
		await moveOrganizationUnitRoute(
			event({
				fetch: fetchMock,
				locals: { user: { role: 'admin' } },
				method: 'PATCH',
				path: '/api/organization-units/unit-id/parent',
				params: { id: 'unit-id' },
				body: { parentId: '  ' }
			}) as never
		);
		const archiveResponse = await archiveOrganizationUnitRoute(
			event({
				fetch: fetchMock,
				locals: { user: { role: 'admin' } },
				method: 'POST',
				path: '/api/organization-units/unit-id/archive',
				params: { id: 'unit-id' },
				body: {}
			}) as never
		);

		expect(organizationUnitMocks.updateOrganizationUnit).toHaveBeenCalledWith(
			fetchMock,
			cookieHeader,
			'unit-id',
			{
				parentId: 'parent-id',
				name: 'Operations',
				code: '',
				description: ''
			}
		);
		expect(organizationUnitMocks.moveOrganizationUnit).toHaveBeenCalledWith(
			fetchMock,
			cookieHeader,
			'unit-id',
			null
		);
		expect(organizationUnitMocks.archiveOrganizationUnit).toHaveBeenCalledWith(
			fetchMock,
			cookieHeader,
			'unit-id'
		);
		expect(await archiveResponse.json()).toEqual({ ok: true });
	});

	it('maps known backend errors to public-safe JSON responses', async () => {
		organizationUnitMocks.getOrganizationUnitTree.mockRejectedValue(
			new organizationUnitMocks.OrganizationUnitApiError(503, 'database unavailable')
		);

		const response = await getOrganizationUnitTreeRoute(
			event({
				locals: { user: { role: 'admin' } },
				path: '/api/organization-units/tree'
			}) as never
		);

		expect(response.status).toBe(502);
		expect(await response.json()).toEqual({ error: 'Failed to load organization units' });
	});
});

function event({
	fetch = vi.fn(),
	locals,
	method = 'GET',
	path,
	params = {},
	body
}: {
	fetch?: typeof globalThis.fetch;
	locals: { user: { role: string } | null };
	method?: string;
	path: string;
	params?: Record<string, string>;
	body?: unknown;
}) {
	return {
		fetch,
		locals,
		params,
		request: new Request(`http://localhost${path}`, {
			method,
			headers: { cookie: cookieHeader, ...(body === undefined ? {} : { 'content-type': 'application/json' }) },
			body: body === undefined ? undefined : JSON.stringify(body)
		})
	};
}
