import { json } from "@sveltejs/kit";

import { listAdminUsers } from "$lib/server/admin";
import type { RequestHandler } from "./$types";
import { adminApiErrorResponse, adminAuthError, optionalQuery } from "../admin-route-utils";

export const GET: RequestHandler = async ({ url, request, fetch, locals }) => {
	const authError = adminAuthError(locals);
	if (authError) return authError;

	try {
		const result = await listAdminUsers(fetch, request.headers.get("cookie"), {
			search: optionalQuery(url, "search"),
			sort: optionalQuery(url, "sort") as "created_at" | "last_login_at" | undefined,
			direction: optionalQuery(url, "direction") as "asc" | "desc" | undefined,
			cursor: optionalQuery(url, "cursor"),
			size: optionalQuery(url, "size")
		});
		return json(result);
	} catch (error) {
		return adminApiErrorResponse(error);
	}
};
