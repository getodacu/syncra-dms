import { json } from "@sveltejs/kit";

import {
	isBillingApiError,
	listBillingOrders,
	type ListBillingOrdersOptions
} from "$lib/server/billing";
import { jsonPublicErrorResponse } from "$lib/server/public-errors";
import type { RequestHandler } from "./$types";

function optionalQuery(url: URL, key: string) {
	const value = url.searchParams.get(key);
	return value && value.trim() ? value.trim() : undefined;
}

export const GET: RequestHandler = async ({ url, fetch, locals }) => {
	if (!locals.user) {
		return json({ error: "authentication required" }, { status: 401 });
	}

	const options: ListBillingOrdersOptions = {
		userId: locals.user.id,
		status: optionalQuery(url, "status") as ListBillingOrdersOptions["status"],
		createdFrom: optionalQuery(url, "created_from"),
		createdTo: optionalQuery(url, "created_to"),
		cursor: optionalQuery(url, "cursor"),
		size: optionalQuery(url, "size"),
		sort: optionalQuery(url, "sort") as ListBillingOrdersOptions["sort"]
	};

	try {
		const result = await listBillingOrders(fetch, options);
		return json(result);
	} catch (error) {
		if (isBillingApiError(error)) {
			return jsonPublicErrorResponse(error.status, error.message);
		}

		throw error;
	}
};
