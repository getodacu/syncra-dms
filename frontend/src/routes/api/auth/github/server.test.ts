import { describe, expect, it, vi } from 'vitest';

const { setGitHubOAuthStateCookieMock, startGitHubOAuthMock } = vi.hoisted(() => ({
	setGitHubOAuthStateCookieMock: vi.fn(),
	startGitHubOAuthMock: vi.fn()
}));

vi.mock('$lib/server/auth', () => ({
	setGitHubOAuthStateCookie: setGitHubOAuthStateCookieMock,
	startGitHubOAuth: startGitHubOAuthMock
}));

import { GET } from './+server';

describe('github oauth start route', () => {
	it('uses the Better Auth-compatible callback route as redirect URI', async () => {
		vi.stubEnv('SYNCRA_APP_ORIGIN', 'https://app.example.com');
		startGitHubOAuthMock.mockResolvedValue({
			authorizationUrl: 'https://github.com/login/oauth/authorize',
			state: 'state-1',
			stateExpiresAt: '2026-06-30T12:00:00Z'
		});
		const fetchMock = vi.fn();
		const cookies = {};

		await expect(
			GET({
				cookies,
				fetch: fetchMock,
				url: new URL('http://localhost/api/auth/github')
			} as never)
		).rejects.toMatchObject({
			status: 303,
			location: 'https://github.com/login/oauth/authorize'
		});

		expect(startGitHubOAuthMock).toHaveBeenCalledWith(
			fetchMock,
			'https://app.example.com/api/auth/callback/github'
		);
		expect(setGitHubOAuthStateCookieMock).toHaveBeenCalledWith(
			cookies,
			'state-1',
			'2026-06-30T12:00:00Z'
		);
	});
});
