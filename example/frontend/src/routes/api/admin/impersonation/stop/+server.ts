import { json } from "@sveltejs/kit";

import { stopAdminImpersonation } from "$lib/server/admin";
import type { RequestHandler } from "./$types";
import { adminApiErrorResponse, adminAuthError } from "../../admin-route-utils";

export const POST: RequestHandler = async ({ request, fetch, locals }) => {
	const authError = adminAuthError(locals);
	if (authError) return authError;

	try {
		const result = await stopAdminImpersonation(fetch, request.headers.get("cookie"));
		return json(result);
	} catch (error) {
		return adminApiErrorResponse(error);
	}
};
