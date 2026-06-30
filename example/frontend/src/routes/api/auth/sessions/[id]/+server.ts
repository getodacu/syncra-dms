import { json } from "@sveltejs/kit";

import { revokeAuthSession } from "$lib/server/auth";
import { authErrorResponse, cookieHeader } from "../../auth-proxy-utils";
import type { RequestHandler } from "./$types";

export const DELETE: RequestHandler = async ({ params, request, fetch, locals }) => {
	if (!locals.user) return json({ error: "authentication required" }, { status: 401 });

	const sessionId = params.id?.trim() ?? "";
	if (!sessionId) return json({ error: "session id is required" }, { status: 400 });

	try {
		return json(await revokeAuthSession(fetch, cookieHeader(request), sessionId));
	} catch (error) {
		return authErrorResponse(error);
	}
};
