import { json } from "@sveltejs/kit";

import { exportDataset, isDatasetApiError } from "$lib/server/datasets";
import { jsonPublicErrorResponse } from "$lib/server/public-errors";
import type { RequestHandler } from "./$types";

function datasetErrorResponse(error: unknown) {
	if (isDatasetApiError(error)) {
		return jsonPublicErrorResponse(error.status, error.message);
	}

	throw error;
}

function normalizeExportFormat(format: string | null) {
	return format === null ? null : format.trim().toLowerCase();
}

function validExportFormat(format: string | null) {
	return format === null || format === "csv" || format === "xlsx";
}

function optionalQuery(url: URL, key: string) {
	return url.searchParams.has(key) ? (url.searchParams.get(key) ?? "") : undefined;
}

export const GET: RequestHandler = async ({ url, fetch, locals, params }) => {
	if (!locals.user) {
		return json({ error: "authentication required" }, { status: 401 });
	}

	const format = normalizeExportFormat(url.searchParams.get("format"));
	if (!validExportFormat(format)) {
		return json({ error: "format must be csv or xlsx" }, { status: 400 });
	}

	try {
		const createdFrom = optionalQuery(url, "created_from");
		const createdTo = optionalQuery(url, "created_to");
		const sort = optionalQuery(url, "sort");
		const exported = await exportDataset(fetch, params.id, {
			userId: locals.user.id,
			...(format ? { format } : {}),
			...(createdFrom !== undefined ? { createdFrom } : {}),
			...(createdTo !== undefined ? { createdTo } : {}),
			...(sort !== undefined ? { sort } : {})
		});

		return new Response(exported.body, {
			status: exported.status,
			headers: exported.headers
		});
	} catch (error) {
		return datasetErrorResponse(error);
	}
};
