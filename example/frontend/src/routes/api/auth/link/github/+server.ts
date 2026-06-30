import { redirect } from "@sveltejs/kit";

import {
	isAuthApiError,
	setGitHubOAuthLinkStateCookie,
	startGitHubAccountLink
} from "$lib/server/auth";
import {
	accountLinkRedirect,
	appOrigin,
	cookieHeader
} from "../../auth-proxy-utils";
import type { RequestHandler } from "./$types";

function githubLinkCallbackURI(url: URL) {
	return `${appOrigin(url)}/api/auth/callback/link/github`;
}

function linkStartStatus(error: unknown) {
	if (!isAuthApiError(error)) throw error;
	if (error.status === 503) return "configuration";
	if (error.status === 401) return "auth";
	return "failed";
}

export const GET: RequestHandler = async ({ cookies, request, fetch, locals, url }) => {
	if (!locals.user) redirect(303, "/login");

	try {
		const started = await startGitHubAccountLink(
			fetch,
			cookieHeader(request),
			githubLinkCallbackURI(url)
		);
		setGitHubOAuthLinkStateCookie(cookies, started.state, started.stateExpiresAt);
		redirect(303, started.authorizationUrl);
	} catch (error) {
		redirect(303, accountLinkRedirect("github", linkStartStatus(error)));
	}
};
