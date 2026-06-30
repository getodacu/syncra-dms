import { json } from "@sveltejs/kit";

import { listAdminBillingOrders } from "$lib/server/admin";
import type { RequestHandler } from "./$types";
import {
	adminApiErrorResponse,
	adminAuthError,
	optionalBooleanQuery,
	optionalQuery
} from "../../admin-route-utils";

export const GET: RequestHandler = async ({ url, request, fetch, locals }) => {
	const authError = adminAuthError(locals);
	if (authError) return authError;

	try {
		const withoutInvoice = optionalBooleanQuery(url, "without_invoice");
		if (!withoutInvoice.ok) return withoutInvoice.error;

		const result = await listAdminBillingOrders(fetch, request.headers.get("cookie"), {
			userId: optionalQuery(url, "user_id"),
			status: optionalQuery(url, "status") as "pending" | "paid" | "failed" | "refunded" | "canceled" | undefined,
			withoutInvoice: withoutInvoice.value,
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
