import { json } from "@sveltejs/kit";
import type { RequestHandler } from "@sveltejs/kit";

import { listAdminBillingInvoices } from "$lib/server/admin";
import { adminApiErrorResponse, adminAuthError, optionalQuery } from "../../admin-route-utils";

export const GET: RequestHandler = async ({ url, request, fetch, locals }) => {
	const authError = adminAuthError(locals);
	if (authError) return authError;

	try {
		const result = await listAdminBillingInvoices(fetch, request.headers.get("cookie"), {
			search: optionalQuery(url, "search"),
			userId: optionalQuery(url, "user_id"),
			createdFrom: optionalQuery(url, "created_from"),
			createdTo: optionalQuery(url, "created_to"),
			cursor: optionalQuery(url, "cursor"),
			size: optionalQuery(url, "size"),
			sort: optionalQuery(url, "sort") as "asc" | "desc" | undefined
		});
		return json(result);
	} catch (error) {
		return adminApiErrorResponse(error);
	}
};
