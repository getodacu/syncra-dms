import { isAuthApiError } from "$lib/server/auth";
import { privateEnv } from "$lib/server/internal-api";
import { jsonPublicErrorResponse } from "$lib/server/public-errors";

export function authErrorResponse(error: unknown) {
	if (isAuthApiError(error)) {
		return jsonPublicErrorResponse(error.status, error.message);
	}

	throw error;
}

export function cookieHeader(request: Request) {
	return request.headers.get("cookie");
}

export function appOrigin(url: URL) {
	const configured = privateEnv("SYNCRA_APP_ORIGIN")?.trim();
	if (!configured) return url.origin;
	try {
		return new URL(configured).origin;
	} catch {
		return url.origin;
	}
}

export function accountLinkRedirect(provider: "google" | "github", status: string) {
	const params = new URLSearchParams({
		account_settings: "linked",
		account_link_provider: provider,
		account_link_status: status
	});
	return `/app?${params.toString()}`;
}
