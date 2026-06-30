import { redirect } from "@sveltejs/kit";

import {
	GOOGLE_OAUTH_LINK_STATE_COOKIE_NAME,
	clearGoogleOAuthLinkStateCookie,
	isAuthApiError,
	linkGoogleAccount
} from "$lib/server/auth";
import {
	accountLinkRedirect,
	appOrigin,
	cookieHeader
} from "../../../auth-proxy-utils";
import type { RequestHandler } from "./$types";

function googleLinkCallbackURI(url: URL) {
	return `${appOrigin(url)}/api/auth/callback/link/google`;
}

function linkCallbackStatus(error: unknown) {
	if (!isAuthApiError(error)) throw error;
	if (error.status === 409) return "conflict";
	if (error.status === 503) return "configuration";
	if (error.status === 401) return "invalid";
	return "failed";
}

export const GET: RequestHandler = async ({ cookies, request, fetch, locals, url }) => {
	const providerError = url.searchParams.get("error");
	if (providerError) {
		clearGoogleOAuthLinkStateCookie(cookies);
		redirect(
			303,
			accountLinkRedirect("google", providerError === "access_denied" ? "denied" : "failed")
		);
	}

	if (!locals.user) {
		clearGoogleOAuthLinkStateCookie(cookies);
		redirect(303, "/login?oauth_error=invalid");
	}

	const code = url.searchParams.get("code")?.trim() ?? "";
	const state = url.searchParams.get("state")?.trim() ?? "";
	const expectedState = cookies.get(GOOGLE_OAUTH_LINK_STATE_COOKIE_NAME)?.trim() ?? "";
	if (!code || !state || state !== expectedState) {
		clearGoogleOAuthLinkStateCookie(cookies);
		redirect(303, accountLinkRedirect("google", "invalid"));
	}

	try {
		await linkGoogleAccount(fetch, cookieHeader(request), {
			code,
			state,
			redirectURI: googleLinkCallbackURI(url)
		});
		clearGoogleOAuthLinkStateCookie(cookies);
		redirect(303, accountLinkRedirect("google", "linked"));
	} catch (error) {
		clearGoogleOAuthLinkStateCookie(cookies);
		redirect(303, accountLinkRedirect("google", linkCallbackStatus(error)));
	}
};
