import { describe, expect, it, vi } from 'vitest';

const { signInEmailMock, setSessionCookieMock } = vi.hoisted(() => ({
	signInEmailMock: vi.fn(),
	setSessionCookieMock: vi.fn()
}));

vi.mock('$lib/server/auth', () => ({
	isAuthApiError: (error: unknown) => error instanceof Error && 'status' in error,
	setPreferredLanguageCookie: vi.fn(),
	setSessionCookie: setSessionCookieMock,
	signInEmail: signInEmailMock
}));

import { actions, load } from './+page.server';

describe('login page', () => {
	it('redirects authenticated users away from login', () => {
		expect(() =>
			load({ locals: { user: { id: 'user-id' } }, url: new URL('http://localhost/login') } as never)
		).toThrow();
	});

	it('signs in with email and sets the session cookie', async () => {
		const auth = {
			session: {
				id: 'session-id',
				token: 'token',
				userId: 'user-id',
				expiresAt: new Date(Date.now() + 60_000).toISOString(),
				createdAt: new Date().toISOString(),
				updatedAt: new Date().toISOString()
			},
			user: { id: 'user-id', email: 'ada@example.com', preferredLanguage: 'en' }
		};
		signInEmailMock.mockResolvedValue(auth);
		const formData = new FormData();
		formData.set('email', 'ADA@EXAMPLE.COM');
		formData.set('password', 'password123');
		const event = {
			request: new Request('http://localhost/login', { method: 'POST', body: formData }),
			fetch: vi.fn(),
			cookies: {},
			locals: {}
		};

		await expect(actions.default(event as never)).rejects.toMatchObject({
			status: 303,
			location: '/app'
		});
		expect(signInEmailMock).toHaveBeenCalledWith(event.fetch, {
			email: 'ada@example.com',
			password: 'password123',
			rememberMe: false
		});
		expect(setSessionCookieMock).toHaveBeenCalledWith(event.cookies, auth.session, false);
	});
});
