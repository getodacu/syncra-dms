import { describe, expect, it, vi, beforeEach } from 'vitest';

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

import { load as layoutLoad } from '../+layout.server';
import { load as appPageLoad } from '../+page.server';
import { actions, load } from './+page.server';

const cookieHeader = 'auth.session_token=token';
const privateUser = {
	id: 'user-id',
	name: 'Ada Lovelace',
	email: 'ada@example.com',
	emailVerified: true,
	image: 'https://example.com/avatar.png',
	preferredLanguage: 'en',
	role: 'admin',
	lastLoginAt: '2026-06-29T00:00:00Z',
	createdAt: '2026-06-01T00:00:00Z',
	updatedAt: '2026-06-30T00:00:00Z'
};
const publicUser = {
	name: 'Ada Lovelace',
	email: 'ada@example.com',
	role: 'admin'
};

describe('app layout server load', () => {
	it('returns public user and session shapes from locals', () => {
		const locals = {
			user: privateUser,
			session: { id: 'session-id', token: 'secret-token', expiresAt: '2026-06-30T00:00:00Z' }
		};

		expect(layoutLoad({ locals } as never)).toEqual({
			user: publicUser,
			session: { expiresAt: '2026-06-30T00:00:00Z' }
		});
	});
});

describe('app dashboard server load', () => {
	it('returns a public session shape without the session token', () => {
		const locals = {
			user: privateUser,
			session: { id: 'session-id', token: 'secret-token', expiresAt: '2026-06-30T00:00:00Z' }
		};

		expect(appPageLoad({ locals } as never)).toEqual({
			user: publicUser,
			session: { expiresAt: '2026-06-30T00:00:00Z' }
		});
	});
});

describe('organization units page server load', () => {
	beforeEach(() => {
		vi.clearAllMocks();
	});

	it('loads the tree with the forwarded cookie and returns admin manage access', async () => {
		const units = [
			{
				id: 'root',
				name: 'Root',
				code: 'ROOT',
				description: 'Root unit',
				createdAt: '2026-06-30T00:00:00Z',
				updatedAt: '2026-06-30T00:00:00Z',
				children: []
			}
		];
		organizationUnitMocks.getOrganizationUnitTree.mockResolvedValue({ units });
		const event = loadEvent({ role: 'admin' });

		const result = await load(event as never);

		expect(organizationUnitMocks.getOrganizationUnitTree).toHaveBeenCalledWith(
			event.fetch,
			cookieHeader
		);
		expect(result).toEqual({
			units,
			loadError: null,
			canManageOrganizationUnits: true,
			selectedId: null
		});
	});

	it('loads the tree for regular users without manage access', async () => {
		organizationUnitMocks.getOrganizationUnitTree.mockResolvedValue({ units: [] });
		const event = loadEvent({ role: 'user' });

		const result = await load(event as never);

		expect(result).toEqual({
			units: [],
			loadError: null,
			canManageOrganizationUnits: false,
			selectedId: null
		});
	});

	it('preserves selectedId from the query string', async () => {
		organizationUnitMocks.getOrganizationUnitTree.mockResolvedValue({ units: [] });
		const event = loadEvent({ role: 'admin' }, 'selected-unit');

		const result = await load(event as never);

		expect(result).toEqual({
			units: [],
			loadError: null,
			canManageOrganizationUnits: true,
			selectedId: 'selected-unit'
		});
	});

	it('maps known API errors to loadError and empty units', async () => {
		organizationUnitMocks.getOrganizationUnitTree.mockRejectedValue(
			new organizationUnitMocks.OrganizationUnitApiError(503, 'backend down')
		);
		const event = loadEvent({ role: 'user' });

		const result = await load(event as never);

		expect(result).toEqual({
			units: [],
			loadError: 'Failed to load organization units',
			canManageOrganizationUnits: false,
			selectedId: null
		});
	});
});

describe('organization units page actions', () => {
	beforeEach(() => {
		vi.clearAllMocks();
	});

	it('creates organization units with blank parentId normalized to null', async () => {
		organizationUnitMocks.createOrganizationUnit.mockResolvedValue({ id: 'created-id' });
		const form = organizationUnitForm({
			name: 'Operations',
			code: 'OPS',
			description: 'Handles operations',
			parentId: '   '
		});
		const event = actionEvent(form);

		const result = await actions.create(event as never);

		expect(organizationUnitMocks.createOrganizationUnit).toHaveBeenCalledWith(
			event.fetch,
			cookieHeader,
			{
				parentId: null,
				name: 'Operations',
				code: 'OPS',
				description: 'Handles operations'
			}
		);
		expect(result).toEqual({ success: true, selectedId: 'created-id' });
	});

	it('updates organization units with id, cookie, and input fields', async () => {
		organizationUnitMocks.updateOrganizationUnit.mockResolvedValue({});
		const form = organizationUnitForm({
			id: 'unit-id',
			name: 'Finance',
			code: 'FIN',
			description: 'Handles finance',
			parentId: 'parent-id'
		});
		const event = actionEvent(form);

		const result = await actions.update(event as never);

		expect(organizationUnitMocks.updateOrganizationUnit).toHaveBeenCalledWith(
			event.fetch,
			cookieHeader,
			'unit-id',
			{
				parentId: 'parent-id',
				name: 'Finance',
				code: 'FIN',
				description: 'Handles finance'
			}
		);
		expect(result).toEqual({ success: true, selectedId: 'unit-id' });
	});

	it('moves organization units with blank parentId normalized to null', async () => {
		organizationUnitMocks.moveOrganizationUnit.mockResolvedValue({});
		const form = new FormData();
		form.set('id', 'unit-id');
		form.set('parentId', '  ');
		const event = actionEvent(form);

		const result = await actions.move(event as never);

		expect(organizationUnitMocks.moveOrganizationUnit).toHaveBeenCalledWith(
			event.fetch,
			cookieHeader,
			'unit-id',
			null
		);
		expect(result).toEqual({ success: true, selectedId: 'unit-id' });
	});

	it('archives organization units with id and cookie', async () => {
		organizationUnitMocks.archiveOrganizationUnit.mockResolvedValue({});
		const form = new FormData();
		form.set('id', 'unit-id');
		const event = actionEvent(form);

		const result = await actions.archive(event as never);

		expect(organizationUnitMocks.archiveOrganizationUnit).toHaveBeenCalledWith(
			event.fetch,
			cookieHeader,
			'unit-id'
		);
		expect(result).toEqual({ success: true, selectedId: null });
	});

	it('maps known action API errors through public status and message', async () => {
		organizationUnitMocks.createOrganizationUnit.mockRejectedValue(
			new organizationUnitMocks.OrganizationUnitApiError(503, 'backend down')
		);
		const event = actionEvent(
			organizationUnitForm({
				name: 'Operations',
				code: 'OPS',
				description: 'Handles operations',
				parentId: ''
			})
		);

		const result = await actions.create(event as never);

		expect(result).toMatchObject({
			status: 502,
			data: {
				error: 'Failed to create organization unit',
				action: 'create',
				selectedId: null,
				values: {
					parentId: null,
					name: 'Operations',
					code: 'OPS',
					description: 'Handles operations'
				}
			}
		});
	});

	it('returns submitted update values on known update errors', async () => {
		organizationUnitMocks.updateOrganizationUnit.mockRejectedValue(
			new organizationUnitMocks.OrganizationUnitApiError(409, 'duplicate active code')
		);
		const event = actionEvent(
			organizationUnitForm({
				id: 'unit-id',
				name: 'Finance',
				code: 'FIN',
				description: 'Handles finance',
				parentId: 'parent-id'
			})
		);

		const result = await actions.update(event as never);

		expect(result).toMatchObject({
			status: 409,
			data: {
				error: 'duplicate active code',
				action: 'update',
				selectedId: 'unit-id',
				values: {
					id: 'unit-id',
					parentId: 'parent-id',
					name: 'Finance',
					code: 'FIN',
					description: 'Handles finance'
				}
			}
		});
	});

	it('returns submitted move values on known move errors', async () => {
		organizationUnitMocks.moveOrganizationUnit.mockRejectedValue(
			new organizationUnitMocks.OrganizationUnitApiError(409, 'invalid move target')
		);
		const form = new FormData();
		form.set('id', 'unit-id');
		form.set('parentId', 'parent-id');
		const event = actionEvent(form);

		const result = await actions.move(event as never);

		expect(result).toMatchObject({
			status: 409,
			data: {
				error: 'invalid move target',
				action: 'move',
				selectedId: 'unit-id',
				values: {
					id: 'unit-id',
					parentId: 'parent-id'
				}
			}
		});
	});

	it('keeps the selected unit on known archive errors', async () => {
		organizationUnitMocks.archiveOrganizationUnit.mockRejectedValue(
			new organizationUnitMocks.OrganizationUnitApiError(500, 'database unavailable')
		);
		const form = new FormData();
		form.set('id', 'unit-id');
		const event = actionEvent(form);

		const result = await actions.archive(event as never);

		expect(result).toMatchObject({
			status: 502,
			data: {
				error: 'Failed to archive organization unit',
				action: 'archive',
				selectedId: 'unit-id',
				values: {
					id: 'unit-id'
				}
			}
		});
	});
});

function loadEvent(user: { role: string } | null, selectedId?: string) {
	const url = new URL('http://localhost/app/organization-units');
	if (selectedId) url.searchParams.set('selectedId', selectedId);
	return {
		fetch: vi.fn(),
		locals: { user },
		request: new Request('http://localhost/app/organization-units', {
			headers: { cookie: cookieHeader }
		}),
		url
	};
}

function actionEvent(form: FormData) {
	return {
		fetch: vi.fn(),
		request: new Request('http://localhost/app/organization-units', {
			method: 'POST',
			headers: { cookie: cookieHeader },
			body: form
		})
	};
}

function organizationUnitForm(values: {
	id?: string;
	parentId: string;
	name: string;
	code: string;
	description: string;
}) {
	const form = new FormData();
	if (values.id) form.set('id', values.id);
	form.set('parentId', values.parentId);
	form.set('name', values.name);
	form.set('code', values.code);
	form.set('description', values.description);
	return form;
}
