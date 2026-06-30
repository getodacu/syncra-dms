import { redirect, type RequestHandler } from '@sveltejs/kit';
import { setGitHubOAuthStateCookie, startGitHubOAuth } from '$lib/server/auth';
import { privateEnv } from '$lib/server/internal-api';

export const GET: RequestHandler = async ({ cookies, fetch, url }) => {
	const redirectURI = `${appOrigin(url)}/api/auth/github/callback`;
	const result = await startGitHubOAuth(fetch, redirectURI);
	setGitHubOAuthStateCookie(cookies, result.state, result.stateExpiresAt);
	redirect(303, result.authorizationUrl);
};

function appOrigin(url: URL) {
	return (privateEnv('SYNCRA_APP_ORIGIN') || url.origin).replace(/\/+$/, '');
}
