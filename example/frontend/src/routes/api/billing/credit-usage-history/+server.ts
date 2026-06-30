import { json } from "@sveltejs/kit";

import {
	isBillingApiError,
	listCreditUsageHistory,
	type ListCreditUsageHistoryOptions
} from "$lib/server/billing";
import { jsonPublicErrorResponse } from "$lib/server/public-errors";
import type { RequestHandler } from "./$types";

/**
 * Extracts a trimmed query string parameter from URL search params.
 * Returns undefined if the key doesn't exist or is empty.
 * 
 * @param url - The request URL
 * @param key - The query parameter key
 */
function optionalQuery(url: URL, key: string) {
	const value = url.searchParams.get(key);
	return value && value.trim() ? value.trim() : undefined;
}

/**
 * GET /api/billing/credit-usage-history
 * Lists credit usage history (purchases and debits) for the authenticated user.
 * Supports pagination filters, cursor, size, and sorting direction.
 */
export const GET: RequestHandler = async ({ url, fetch, locals }) => {
	if (!locals.user) {
		return json({ error: "authentication required" }, { status: 401 });
	}

	const options: ListCreditUsageHistoryOptions = {
		userId: locals.user.id,
		type: optionalQuery(url, "type") as ListCreditUsageHistoryOptions["type"],
		createdFrom: optionalQuery(url, "created_from"),
		createdTo: optionalQuery(url, "created_to"),
		cursor: optionalQuery(url, "cursor"),
		size: optionalQuery(url, "size"),
		sort: optionalQuery(url, "sort") as ListCreditUsageHistoryOptions["sort"]
	};

	try {
		const result = await listCreditUsageHistory(fetch, options);
		return json(result);
	} catch (error) {
		if (isBillingApiError(error)) {
			return jsonPublicErrorResponse(error.status, error.message);
		}

		throw error;
	}
};
