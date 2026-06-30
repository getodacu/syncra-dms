import { afterEach, describe, expect, it, vi } from 'vitest';

import {
	AUTH_SESSION_COOKIE_NAME,
	GOOGLE_OAUTH_STATE_COOKIE_NAME,
	AuthApiError,
	clearGitHubOAuthStateCookie,
	getSession,
	requestPasswordReset,
	setGoogleOAuthStateCookie,
	setSessionCookie,
	signUpEmail
} from './auth';

describe('server auth client', () => {
	afterEach(() => {
		vi.unstubAllEnvs();
	});

	it('sends internal and delivery tokens only from server helpers', async () => {
		vi.stubEnv('SYNCRA_API_BASE_URL', 'http://api.test');
		vi.stubEnv('SYNCRA_INTERNAL_API_TOKEN', 'internal-token');
		vi.stubEnv('AUTH_DELIVERY_TOKEN', 'delivery-token');
		const fetch = vi.fn(async () => jsonResponse({ ok: true, verificationCode: '123456' }));

		await signUpEmail(fetch, {
			name: 'Ada Lovelace',
			email: 'ada@example.com',
			password: 'password123'
		});

		expect(fetch).toHaveBeenCalledWith('http://api.test/api/auth/sign-up/email', {
			method: 'POST',
			headers: expect.any(Headers),
			body: JSON.stringify({
				name: 'Ada Lovelace',
				email: 'ada@example.com',
				password: 'password123'
			})
		});
		const firstCall = fetch.mock.calls[0] as unknown as [string, RequestInit];
		const headers = firstCall[1].headers as Headers;
		expect(headers.get('X-Syncra-Internal-Token')).toBe('internal-token');
		expect(headers.get('X-Syncra-Auth-Delivery-Token')).toBe('delivery-token');
	});

	it('forwards cookies for session lookups', async () => {
		vi.stubEnv('SYNCRA_API_BASE_URL', 'http://api.test');
		vi.stubEnv('SYNCRA_INTERNAL_API_TOKEN', 'internal-token');
		const fetch = vi.fn(async () => jsonResponse(null));

		await getSession(fetch, 'auth.session_token=abc');

		const firstCall = fetch.mock.calls[0] as unknown as [string, RequestInit];
		const headers = firstCall[1].headers as Headers;
		expect(headers.get('cookie')).toBe('auth.session_token=abc');
	});

	it('throws AuthApiError for backend failures', async () => {
		vi.stubEnv('SYNCRA_API_BASE_URL', 'http://api.test');
		vi.stubEnv('SYNCRA_INTERNAL_API_TOKEN', 'internal-token');
		const fetch = vi.fn(async () => jsonResponse({ error: 'invalid email' }, 400));

		await expect(requestPasswordReset(fetch, 'bad')).rejects.toMatchObject(
			new AuthApiError(400, 'invalid email')
		);
	});

	it('sets the Better Auth-compatible session cookie', () => {
		vi.stubEnv('AUTH_COOKIE_SECURE', 'true');
		const cookies = {
			set: vi.fn()
		};

		setSessionCookie(
			cookies as never,
			{
				id: 'session-id',
				token: 'session-token',
				userId: 'user-id',
				expiresAt: new Date(Date.now() + 60_000).toISOString(),
				createdAt: new Date().toISOString(),
				updatedAt: new Date().toISOString()
			},
			true
		);

		expect(cookies.set).toHaveBeenCalledWith(
			AUTH_SESSION_COOKIE_NAME,
			'session-token',
			expect.objectContaining({
				path: '/',
				httpOnly: true,
				sameSite: 'lax',
				secure: true
			})
		);
	});

	it('sets and clears OAuth state cookies server-side', () => {
		vi.stubEnv('AUTH_COOKIE_SECURE', 'false');
		const cookies = {
			set: vi.fn(),
			delete: vi.fn()
		};

		setGoogleOAuthStateCookie(
			cookies as never,
			'oauth-state',
			new Date(Date.now() + 60_000).toISOString()
		);
		clearGitHubOAuthStateCookie(cookies as never);

		expect(cookies.set).toHaveBeenCalledWith(
			GOOGLE_OAUTH_STATE_COOKIE_NAME,
			'oauth-state',
			expect.objectContaining({
				path: '/',
				httpOnly: true,
				sameSite: 'lax',
				secure: false
			})
		);
		expect(cookies.delete).toHaveBeenCalledWith('auth.github_oauth_state', { path: '/' });
	});
});

function jsonResponse(body: unknown, status = 200) {
	return new Response(JSON.stringify(body), {
		status,
		headers: { 'content-type': 'application/json' }
	});
}
