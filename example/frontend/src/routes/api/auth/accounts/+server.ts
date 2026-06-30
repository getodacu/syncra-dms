import { json } from "@sveltejs/kit";

import { listAuthAccounts } from "$lib/server/auth";
import { authErrorResponse, cookieHeader } from "../auth-proxy-utils";
import type { RequestHandler } from "./$types";

export const GET: RequestHandler = async ({ request, fetch, locals }) => {
	if (!locals.user) return json({ error: "authentication required" }, { status: 401 });

	try {
		return json(await listAuthAccounts(fetch, cookieHeader(request)));
	} catch (error) {
		return authErrorResponse(error);
	}
};
