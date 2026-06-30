import { describe, expect, it, vi } from 'vitest';

const { setGoogleOAuthStateCookieMock, startGoogleOAuthMock } = vi.hoisted(() => ({
	setGoogleOAuthStateCookieMock: vi.fn(),
	startGoogleOAuthMock: vi.fn()
}));

vi.mock('$lib/server/auth', () => ({
	setGoogleOAuthStateCookie: setGoogleOAuthStateCookieMock,
	startGoogleOAuth: startGoogleOAuthMock
}));

import { GET } from './+server';

describe('google oauth start route', () => {
	it('uses the Better Auth-compatible callback route as redirect URI', async () => {
		vi.stubEnv('SYNCRA_APP_ORIGIN', 'https://app.example.com');
		startGoogleOAuthMock.mockResolvedValue({
			authorizationUrl: 'https://accounts.google.com/auth',
			state: 'state-1',
			stateExpiresAt: '2026-06-30T12:00:00Z'
		});
		const fetchMock = vi.fn();
		const cookies = {};

		await expect(
			GET({
				cookies,
				fetch: fetchMock,
				url: new URL('http://localhost/api/auth/google')
			} as never)
		).rejects.toMatchObject({
			status: 303,
			location: 'https://accounts.google.com/auth'
		});

		expect(startGoogleOAuthMock).toHaveBeenCalledWith(
			fetchMock,
			'https://app.example.com/api/auth/callback/google'
		);
		expect(setGoogleOAuthStateCookieMock).toHaveBeenCalledWith(
			cookies,
			'state-1',
			'2026-06-30T12:00:00Z'
		);
	});
});
