import { json } from "@sveltejs/kit";

import { isAuthApiError, regenerateWebhookSecret } from "$lib/server/auth";
import { jsonPublicErrorResponse } from "$lib/server/public-errors";
import type { RequestHandler } from "./$types";

function authErrorResponse(error: unknown) {
	if (isAuthApiError(error)) {
		return jsonPublicErrorResponse(error.status, error.message);
	}

	throw error;
}

function cookieHeader(request: Request) {
	return request.headers.get("cookie");
}

function userId(locals: App.Locals) {
	return locals.user?.id ?? "";
}

export const POST: RequestHandler = async ({ request, fetch, locals }) => {
	if (!locals.user) return json({ error: "authentication required" }, { status: 401 });

	try {
		return json(await regenerateWebhookSecret(fetch, cookieHeader(request), userId(locals)));
	} catch (error) {
		return authErrorResponse(error);
	}
};
