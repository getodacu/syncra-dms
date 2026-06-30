import { beforeEach, describe, expect, it, vi } from 'vitest';

const {
	clearGoogleOAuthStateCookieMock,
	setPreferredLanguageCookieMock,
	setSessionCookieMock,
	signInGoogleOAuthMock
} = vi.hoisted(() => ({
	clearGoogleOAuthStateCookieMock: vi.fn(),
	setPreferredLanguageCookieMock: vi.fn(),
	setSessionCookieMock: vi.fn(),
	signInGoogleOAuthMock: vi.fn()
}));

vi.mock('$lib/server/auth', () => ({
	GOOGLE_OAUTH_STATE_COOKIE_NAME: 'auth.google_oauth_state',
	clearGoogleOAuthStateCookie: clearGoogleOAuthStateCookieMock,
	setPreferredLanguageCookie: setPreferredLanguageCookieMock,
	setSessionCookie: setSessionCookieMock,
	signInGoogleOAuth: signInGoogleOAuthMock
}));

import { GET } from './+server';

describe('google oauth callback route', () => {
	beforeEach(() => {
		vi.clearAllMocks();
	});

	it('passes the Better Auth-compatible callback URI to the internal auth API', async () => {
		vi.stubEnv('SYNCRA_APP_ORIGIN', 'https://app.example.com');
		const auth = authPayload();
		signInGoogleOAuthMock.mockResolvedValue(auth);
		const event = callbackEvent('/api/auth/callback/google?code=code-1&state=state-1');

		await expect(GET(event as never)).rejects.toMatchObject({
			status: 303,
			location: '/app'
		});

		expect(signInGoogleOAuthMock).toHaveBeenCalledWith(event.fetch, {
			code: 'code-1',
			state: 'state-1',
			redirectURI: 'https://app.example.com/api/auth/callback/google'
		});
		expect(setSessionCookieMock).toHaveBeenCalledWith(event.cookies, auth.session, true);
		expect(setPreferredLanguageCookieMock).toHaveBeenCalledWith(event.cookies, 'en');
		expect(clearGoogleOAuthStateCookieMock).toHaveBeenCalledWith(event.cookies);
	});

	it('redirects invalid state to login', async () => {
		const event = callbackEvent('/api/auth/callback/google?code=code-1&state=bad-state');

		await expect(GET(event as never)).rejects.toMatchObject({
			status: 303,
			location: '/login?oauth_error=invalid'
		});
		expect(signInGoogleOAuthMock).not.toHaveBeenCalled();
		expect(clearGoogleOAuthStateCookieMock).toHaveBeenCalledWith(event.cookies);
	});
});

function callbackEvent(path: string) {
	return {
		cookies: {
			get: vi.fn((name: string) => (name === 'auth.google_oauth_state' ? 'state-1' : undefined))
		},
		fetch: vi.fn(),
		locals: {},
		url: new URL(`https://app.example.com${path}`)
	};
}

function authPayload() {
	return {
		session: {
			id: 'session-id',
			token: 'session-token',
			userId: 'user-id',
			expiresAt: new Date(Date.now() + 60_000).toISOString(),
			createdAt: new Date().toISOString(),
			updatedAt: new Date().toISOString()
		},
		user: {
			id: 'user-id',
			name: 'Ada',
			email: 'ada@example.com',
			emailVerified: true,
			image: null,
			preferredLanguage: 'en',
			role: 'user',
			lastLoginAt: null,
			createdAt: new Date().toISOString(),
			updatedAt: new Date().toISOString()
		}
	};
}
