import { redirect } from "@sveltejs/kit";

import {
	GITHUB_OAUTH_STATE_COOKIE_NAME,
	clearGitHubOAuthStateCookie,
	isAuthApiError,
	setPreferredLanguageCookie,
	setSessionCookie,
	signInGitHubOAuth
} from "$lib/server/auth";
import { privateEnv } from "$lib/server/internal-api";
import type { RequestHandler } from "./$types";

function appOrigin(url: URL) {
	const configured = privateEnv("SYNCRA_APP_ORIGIN")?.trim();
	if (!configured) return url.origin;
	try {
		return new URL(configured).origin;
	} catch {
		return url.origin;
	}
}

function githubCallbackURI(url: URL) {
	return `${appOrigin(url)}/api/auth/callback/github`;
}

function loginRedirect(reason: string): never {
	redirect(303, `/login?oauth_error=${reason}`);
}

export const GET: RequestHandler = async ({ cookies, fetch, locals, url }) => {
	const providerError = url.searchParams.get("error");
	if (providerError) {
		clearGitHubOAuthStateCookie(cookies);
		loginRedirect(providerError === "access_denied" ? "denied" : "failed");
	}

	const code = url.searchParams.get("code")?.trim() ?? "";
	const state = url.searchParams.get("state")?.trim() ?? "";
	const stateCookie = cookies.get(GITHUB_OAUTH_STATE_COOKIE_NAME)?.trim() ?? "";
	if (!code || !state || !stateCookie || state !== stateCookie) {
		clearGitHubOAuthStateCookie(cookies);
		loginRedirect("invalid");
	}

	try {
		const auth = await signInGitHubOAuth(fetch, {
			code,
			state,
			redirectURI: githubCallbackURI(url)
		});
		setSessionCookie(cookies, auth.session, true);
		setPreferredLanguageCookie(cookies, auth.user.preferredLanguage);
		clearGitHubOAuthStateCookie(cookies);
		locals.session = auth.session;
		locals.user = auth.user;
		redirect(303, "/app");
	} catch (error) {
		clearGitHubOAuthStateCookie(cookies);
		if (!isAuthApiError(error)) throw error;
		const reason = error.status === 503 ? "configuration" : "failed";
		loginRedirect(reason);
	}
};
