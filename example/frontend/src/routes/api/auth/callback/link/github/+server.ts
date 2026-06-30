import { redirect } from "@sveltejs/kit";

import {
	GITHUB_OAUTH_LINK_STATE_COOKIE_NAME,
	clearGitHubOAuthLinkStateCookie,
	isAuthApiError,
	linkGitHubAccount
} from "$lib/server/auth";
import {
	accountLinkRedirect,
	appOrigin,
	cookieHeader
} from "../../../auth-proxy-utils";
import type { RequestHandler } from "./$types";

function githubLinkCallbackURI(url: URL) {
	return `${appOrigin(url)}/api/auth/callback/link/github`;
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
		clearGitHubOAuthLinkStateCookie(cookies);
		redirect(
			303,
			accountLinkRedirect("github", providerError === "access_denied" ? "denied" : "failed")
		);
	}

	if (!locals.user) {
		clearGitHubOAuthLinkStateCookie(cookies);
		redirect(303, "/login?oauth_error=invalid");
	}

	const code = url.searchParams.get("code")?.trim() ?? "";
	const state = url.searchParams.get("state")?.trim() ?? "";
	const expectedState = cookies.get(GITHUB_OAUTH_LINK_STATE_COOKIE_NAME)?.trim() ?? "";
	if (!code || !state || state !== expectedState) {
		clearGitHubOAuthLinkStateCookie(cookies);
		redirect(303, accountLinkRedirect("github", "invalid"));
	}

	try {
		await linkGitHubAccount(fetch, cookieHeader(request), {
			code,
			state,
			redirectURI: githubLinkCallbackURI(url)
		});
		clearGitHubOAuthLinkStateCookie(cookies);
		redirect(303, accountLinkRedirect("github", "linked"));
	} catch (error) {
		clearGitHubOAuthLinkStateCookie(cookies);
		redirect(303, accountLinkRedirect("github", linkCallbackStatus(error)));
	}
};
