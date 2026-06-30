import { json } from "@sveltejs/kit";

import { getCreditBalance, isBillingApiError } from "$lib/server/billing";
import { jsonPublicErrorResponse } from "$lib/server/public-errors";
import type { RequestHandler } from "./$types";

export const GET: RequestHandler = async ({ fetch, locals }) => {
	if (!locals.user) {
		return json({ error: "authentication required" }, { status: 401 });
	}

	try {
		const balance = await getCreditBalance(fetch, locals.user.id);
		return json(balance);
	} catch (error) {
		if (isBillingApiError(error)) {
			return jsonPublicErrorResponse(error.status, error.message);
		}

		throw error;
	}
};
