import { json } from "@sveltejs/kit";

import { fetchBillingInvoicePDF, isBillingApiError } from "$lib/server/billing";
import { jsonPublicErrorResponse } from "$lib/server/public-errors";
import type { RequestHandler } from "./$types";

export const GET: RequestHandler = async ({ url, params, fetch, locals }) => {
	if (!locals.user) {
		return json({ error: "authentication required" }, { status: 401 });
	}

	try {
		const result = await fetchBillingInvoicePDF(fetch, locals.user.id, params.id, {
			download: url.searchParams.get("download") === "1"
		});
		return new Response(result.body, {
			status: result.status,
			headers: result.headers
		});
	} catch (error) {
		if (isBillingApiError(error)) {
			return jsonPublicErrorResponse(error.status, error.message);
		}

		throw error;
	}
};
