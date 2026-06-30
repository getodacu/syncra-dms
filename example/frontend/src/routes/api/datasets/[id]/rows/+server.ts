import { json } from "@sveltejs/kit";

import { isDatasetApiError, listDatasetRows } from "$lib/server/datasets";
import { jsonPublicErrorResponse } from "$lib/server/public-errors";
import type { RequestHandler } from "./$types";

function datasetErrorResponse(error: unknown) {
	if (isDatasetApiError(error)) {
		return jsonPublicErrorResponse(error.status, error.message);
	}

	throw error;
}

function optionalQuery(url: URL, key: string) {
	return url.searchParams.has(key) ? (url.searchParams.get(key) ?? "") : undefined;
}

export const GET: RequestHandler = async ({ url, fetch, locals, params }) => {
	if (!locals.user) {
		return json({ error: "authentication required" }, { status: 401 });
	}

	try {
		const createdFrom = optionalQuery(url, "created_from");
		const createdTo = optionalQuery(url, "created_to");
		const cursor = optionalQuery(url, "cursor");
		const size = optionalQuery(url, "size");
		const sort = optionalQuery(url, "sort");
		const page = await listDatasetRows(fetch, params.id, {
			userId: locals.user.id,
			...(createdFrom !== undefined ? { createdFrom } : {}),
			...(createdTo !== undefined ? { createdTo } : {}),
			...(cursor !== undefined ? { cursor } : {}),
			...(size !== undefined ? { size } : {}),
			...(sort !== undefined ? { sort } : {})
		});

		return json(page);
	} catch (error) {
		return datasetErrorResponse(error);
	}
};
