import { redirect } from "@sveltejs/kit";

import {
	isAuthApiError,
	setGitHubOAuthStateCookie,
	startGitHubOAuth
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

export const GET: RequestHandler = async ({ cookies, fetch, locals, url }) => {
	if (locals.user) redirect(303, "/app");

	try {
		const started = await startGitHubOAuth(fetch, githubCallbackURI(url));
		setGitHubOAuthStateCookie(cookies, started.state, started.stateExpiresAt);
		redirect(303, started.authorizationUrl);
	} catch (error) {
		if (!isAuthApiError(error)) throw error;
		const reason = error.status === 503 ? "configuration" : "failed";
		redirect(303, `/login?oauth_error=${reason}`);
	}
};
