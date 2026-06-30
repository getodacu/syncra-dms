import { json } from "@sveltejs/kit";

import { fetchAdminBillingInvoicePDF, generateAdminBillingInvoicePDF } from "$lib/server/admin";
import type { RequestHandler } from "./$types";
import { adminApiErrorResponse, adminAuthError } from "../../../../admin-route-utils";

export const POST: RequestHandler = async ({ request, params, fetch, locals }) => {
	const authError = adminAuthError(locals);
	if (authError) return authError;

	try {
		const result = await generateAdminBillingInvoicePDF(
			fetch,
			request.headers.get("cookie"),
			params.id
		);
		return json(result);
	} catch (error) {
		return adminApiErrorResponse(error);
	}
};

export const GET: RequestHandler = async ({ request, url, params, fetch, locals }) => {
	const authError = adminAuthError(locals);
	if (authError) return authError;

	try {
		const result = await fetchAdminBillingInvoicePDF(
			fetch,
			request.headers.get("cookie"),
			params.id,
			{ download: url.searchParams.get("download") === "1" }
		);
		return new Response(result.body, {
			status: result.status,
			headers: result.headers
		});
	} catch (error) {
		return adminApiErrorResponse(error);
	}
};
