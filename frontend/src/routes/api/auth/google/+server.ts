import { redirect, type RequestHandler } from '@sveltejs/kit';
import { privateEnv } from '$lib/server/internal-api';
import { setGoogleOAuthStateCookie, startGoogleOAuth } from '$lib/server/auth';

export const GET: RequestHandler = async ({ cookies, fetch, url }) => {
	const redirectURI = `${appOrigin(url)}/api/auth/callback/google`;
	const result = await startGoogleOAuth(fetch, redirectURI);
	setGoogleOAuthStateCookie(cookies, result.state, result.stateExpiresAt);
	redirect(303, result.authorizationUrl);
};

function appOrigin(url: URL) {
	return (privateEnv('SYNCRA_APP_ORIGIN') || url.origin).replace(/\/+$/, '');
}
