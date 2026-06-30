import { json } from "@sveltejs/kit";

import { getDashboardSummary, isDashboardApiError } from "$lib/server/dashboard";
import { jsonPublicErrorResponse } from "$lib/server/public-errors";
import type { RequestHandler } from "./$types";

const ALLOWED_RANGES = new Set(["7d", "30d", "90d"]);

export const GET: RequestHandler = async ({ fetch, locals, url }) => {
	if (!locals.user) {
		return json({ error: "authentication required" }, { status: 401 });
	}

	const range = url.searchParams.get("range") ?? "";
	if (range && !ALLOWED_RANGES.has(range)) {
		return json({ error: "invalid range" }, { status: 400 });
	}

	try {
		const summary = await getDashboardSummary(fetch, {
			userId: locals.user.id,
			...(range ? { range } : {})
		});
		return json(summary);
	} catch (error) {
		if (isDashboardApiError(error)) {
			return jsonPublicErrorResponse(error.status, error.message);
		}

		throw error;
	}
};
