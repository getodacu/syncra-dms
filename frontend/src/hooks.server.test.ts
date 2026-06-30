import { describe, expect, it, vi } from 'vitest';

const { getSessionMock } = vi.hoisted(() => ({
	getSessionMock: vi.fn()
}));

vi.mock('$lib/server/auth', () => ({
	clearSessionCookie: vi.fn(),
	getSession: getSessionMock,
	hasSessionCookie: vi.fn(() => false),
	setPreferredLanguageCookie: vi.fn()
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
	return {
		request,
		url: new URL(url),
		route: { id: new URL(url).pathname },
		fetch: vi.fn(),
		cookies: {},
		locals: {}
	};
}
