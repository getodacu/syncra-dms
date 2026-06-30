import { redirect, type RequestHandler } from '@sveltejs/kit';
import {
	GITHUB_OAUTH_STATE_COOKIE_NAME,
	clearGitHubOAuthStateCookie,
	setPreferredLanguageCookie,
	setSessionCookie,
	signInGitHubOAuth
} from '$lib/server/auth';
import { privateEnv } from '$lib/server/internal-api';

export const GET: RequestHandler = async ({ cookies, fetch, locals, url }) => {
	const code = url.searchParams.get('code') ?? '';
	const state = url.searchParams.get('state') ?? '';
	if (!code || !state || cookies.get(GITHUB_OAUTH_STATE_COOKIE_NAME) !== state) {
		clearGitHubOAuthStateCookie(cookies);
		redirect(303, '/login?oauth_error=invalid');
	}
	try {
		const auth = await signInGitHubOAuth(fetch, {
			code,
			state,
			redirectURI: `${appOrigin(url)}/api/auth/github/callback`
		});
		setSessionCookie(cookies, auth.session, true);
		setPreferredLanguageCookie(cookies, auth.user.preferredLanguage);
		locals.session = auth.session;
		locals.user = auth.user;
		clearGitHubOAuthStateCookie(cookies);
	} catch {
		clearGitHubOAuthStateCookie(cookies);
		redirect(303, '/login?oauth_error=failed');
	}
	redirect(303, '/app');
};

function appOrigin(url: URL) {
	return (privateEnv('SYNCRA_APP_ORIGIN') || url.origin).replace(/\/+$/, '');
}
