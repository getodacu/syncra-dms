import { json } from "@sveltejs/kit";

import { generateAdminBillingOrderInvoice } from "$lib/server/admin";
import type { RequestHandler } from "./$types";
import { adminApiErrorResponse, adminAuthError } from "../../../../admin-route-utils";

export const POST: RequestHandler = async ({ request, params, fetch, locals }) => {
	const authError = adminAuthError(locals);
	if (authError) return authError;

	try {
		const result = await generateAdminBillingOrderInvoice(
			fetch,
			request.headers.get("cookie"),
			params.id
		);
		return json(result, { status: 201 });
	} catch (error) {
		return adminApiErrorResponse(error);
	}
};
