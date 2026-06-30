import { beforeEach, describe, expect, it, vi } from 'vitest';

const {
	clearGitHubOAuthStateCookieMock,
	setPreferredLanguageCookieMock,
	setSessionCookieMock,
	signInGitHubOAuthMock
} = vi.hoisted(() => ({
	clearGitHubOAuthStateCookieMock: vi.fn(),
	setPreferredLanguageCookieMock: vi.fn(),
	setSessionCookieMock: vi.fn(),
	signInGitHubOAuthMock: vi.fn()
}));

vi.mock('$lib/server/auth', () => ({
	GITHUB_OAUTH_STATE_COOKIE_NAME: 'auth.github_oauth_state',
	clearGitHubOAuthStateCookie: clearGitHubOAuthStateCookieMock,
	setPreferredLanguageCookie: setPreferredLanguageCookieMock,
	setSessionCookie: setSessionCookieMock,
	signInGitHubOAuth: signInGitHubOAuthMock
}));

import { GET } from './+server';

describe('github oauth callback route', () => {
	beforeEach(() => {
		vi.clearAllMocks();
	});

	it('passes the Better Auth-compatible callback URI to the internal auth API', async () => {
		vi.stubEnv('SYNCRA_APP_ORIGIN', 'https://app.example.com');
		const auth = authPayload();
		signInGitHubOAuthMock.mockResolvedValue(auth);
		const event = callbackEvent('/api/auth/callback/github?code=code-1&state=state-1');

		await expect(GET(event as never)).rejects.toMatchObject({
			status: 303,
			location: '/app'
		});

		expect(signInGitHubOAuthMock).toHaveBeenCalledWith(event.fetch, {
			code: 'code-1',
			state: 'state-1',
			redirectURI: 'https://app.example.com/api/auth/callback/github'
		});
		expect(setSessionCookieMock).toHaveBeenCalledWith(event.cookies, auth.session, true);
		expect(setPreferredLanguageCookieMock).toHaveBeenCalledWith(event.cookies, 'en');
		expect(clearGitHubOAuthStateCookieMock).toHaveBeenCalledWith(event.cookies);
	});

	it('redirects invalid state to login', async () => {
		const event = callbackEvent('/api/auth/callback/github?code=code-1&state=bad-state');

		await expect(GET(event as never)).rejects.toMatchObject({
			status: 303,
			location: '/login?oauth_error=invalid'
		});
		expect(signInGitHubOAuthMock).not.toHaveBeenCalled();
		expect(clearGitHubOAuthStateCookieMock).toHaveBeenCalledWith(event.cookies);
	});
});

function callbackEvent(path: string) {
	return {
		cookies: {
			get: vi.fn((name: string) => (name === 'auth.github_oauth_state' ? 'state-1' : undefined))
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
