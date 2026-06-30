import { redirect } from "@sveltejs/kit";

import {
	GOOGLE_OAUTH_STATE_COOKIE_NAME,
	clearGoogleOAuthStateCookie,
	isAuthApiError,
	setPreferredLanguageCookie,
	setSessionCookie,
	signInGoogleOAuth
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

function googleCallbackURI(url: URL) {
	return `${appOrigin(url)}/api/auth/callback/google`;
}

function loginRedirect(reason: string): never {
	redirect(303, `/login?oauth_error=${reason}`);
}

export const GET: RequestHandler = async ({ cookies, fetch, locals, url }) => {
	const providerError = url.searchParams.get("error");
	if (providerError) {
		clearGoogleOAuthStateCookie(cookies);
		loginRedirect(providerError === "access_denied" ? "denied" : "failed");
	}

	const code = url.searchParams.get("code")?.trim() ?? "";
	const state = url.searchParams.get("state")?.trim() ?? "";
	const stateCookie = cookies.get(GOOGLE_OAUTH_STATE_COOKIE_NAME)?.trim() ?? "";
	if (!code || !state || !stateCookie || state !== stateCookie) {
		clearGoogleOAuthStateCookie(cookies);
		loginRedirect("invalid");
	}

	try {
		const auth = await signInGoogleOAuth(fetch, {
			code,
			state,
			redirectURI: googleCallbackURI(url)
		});
		setSessionCookie(cookies, auth.session, true);
		setPreferredLanguageCookie(cookies, auth.user.preferredLanguage);
		clearGoogleOAuthStateCookie(cookies);
		locals.session = auth.session;
		locals.user = auth.user;
		redirect(303, "/app");
	} catch (error) {
		clearGoogleOAuthStateCookie(cookies);
		if (!isAuthApiError(error)) throw error;
		const reason = error.status === 503 ? "configuration" : "failed";
		loginRedirect(reason);
	}
};
