import { beforeEach, describe, expect, it, vi } from 'vitest';

const { getMyPermissionsMock, getSessionMock } = vi.hoisted(() => ({
	getMyPermissionsMock: vi.fn(),
	getSessionMock: vi.fn()
}));

vi.mock('$lib/server/auth', () => ({
	clearSessionCookie: vi.fn(),
	getSession: getSessionMock,
	hasSessionCookie: vi.fn(() => false),
	setPreferredLanguageCookie: vi.fn()
}));

vi.mock('$lib/server/rbac', () => ({
	getMyPermissions: getMyPermissionsMock
}));

vi.mock('$lib/paraglide/server', () => ({
	paraglideMiddleware: (_request: Request, handler: (input: { request: Request; locale: string }) => unknown) =>
		handler({ request: _request, locale: 'en' })
}));

vi.mock('$lib/paraglide/runtime', () => ({
	getTextDirection: () => 'ltr'
}));

import { handle } from './hooks.server';

describe('server hooks auth guards', () => {
	beforeEach(() => {
		vi.clearAllMocks();
	});

	it('loads effective permission codes for authenticated sessions', async () => {
		getSessionMock.mockResolvedValue({
			session: { id: 'session-id' },
			user: { id: 'user-id', preferredLanguage: 'en' }
		});
		getMyPermissionsMock.mockResolvedValue({
			permissions: [
				{ code: 'user.view', scopeType: 'global', source: 'user_role' },
				{ code: 'role.view', scopeType: 'global', source: 'user_role' }
			]
		});
		const event = hookEvent('http://localhost/app');
		const resolve = vi.fn(async () => new Response('ok'));

		await handle({ event, resolve } as never);

		expect(getMyPermissionsMock).toHaveBeenCalledWith(event.fetch, null);
		expect(event.locals.permissions).toEqual(['user.view', 'role.view']);
		expect(resolve).toHaveBeenCalled();
	});

	it('keeps authenticated users when permission loading fails', async () => {
		getSessionMock.mockResolvedValue({
			session: { id: 'session-id' },
			user: { id: 'user-id', preferredLanguage: 'en' }
		});
		getMyPermissionsMock.mockRejectedValue(new Error('rbac unavailable'));
		const event = hookEvent('http://localhost/app');
		const resolve = vi.fn(async () => new Response('ok'));

		await handle({ event, resolve } as never);

		expect(event.locals.user).toEqual({ id: 'user-id', preferredLanguage: 'en' });
		expect(event.locals.permissions).toEqual([]);
		expect(resolve).toHaveBeenCalled();
	});

	it('redirects guests from protected app routes', async () => {
		getSessionMock.mockResolvedValue(null);
		await expect(
			handle({
				event: hookEvent('http://localhost/app'),
				resolve: vi.fn()
			} as never)
		).rejects.toMatchObject({ status: 303, location: '/login' });
	});

	it('redirects authenticated users away from guest auth routes', async () => {
		getSessionMock.mockResolvedValue({
			session: { id: 'session-id' },
			user: { id: 'user-id', preferredLanguage: 'en' }
		});
		await expect(
			handle({
				event: hookEvent('http://localhost/login'),
				resolve: vi.fn()
			} as never)
		).rejects.toMatchObject({ status: 303, location: '/app' });
	});
});

function hookEvent(url: string) {
	const request = new Request(url);
	const locals: {
		permissions?: string[];
		user?: unknown;
	} = {};
	return {
		request,
		url: new URL(url),
		route: { id: new URL(url).pathname },
		fetch: vi.fn(),
		cookies: {},
		locals
	};
}
