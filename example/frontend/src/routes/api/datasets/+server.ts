import { json } from "@sveltejs/kit";

import { createDataset, isDatasetApiError, listDatasetsPage } from "$lib/server/datasets";
import { jsonPublicErrorResponse } from "$lib/server/public-errors";
import { readCreateDatasetRequest } from "./dataset-request";
import type { RequestHandler } from "./$types";

function datasetErrorResponse(error: unknown) {
	if (isDatasetApiError(error)) {
		return jsonPublicErrorResponse(error.status, error.message);
	}

	throw error;
}

export const GET: RequestHandler = async ({ url, fetch, locals }) => {
	if (!locals.user) {
		return json({ error: "authentication required" }, { status: 401 });
	}

	try {
		const page = await listDatasetsPage(fetch, {
			userId: locals.user.id,
			...(url.searchParams.has("cursor")
				? { cursor: url.searchParams.get("cursor") ?? "" }
				: {}),
			...(url.searchParams.has("size") ? { size: url.searchParams.get("size") ?? "" } : {}),
			...(url.searchParams.has("sort") ? { sort: url.searchParams.get("sort") ?? "" } : {})
		});

		return json(page);
	} catch (error) {
		return datasetErrorResponse(error);
	}
};

export const POST: RequestHandler = async ({ request, fetch, locals }) => {
	if (!locals.user) {
		return json({ error: "authentication required" }, { status: 401 });
	}

	const bodyResult = await readCreateDatasetRequest(request);
	if (!bodyResult.ok) {
		return json({ error: bodyResult.error }, { status: 400 });
	}

	try {
		const dataset = await createDataset(fetch, bodyResult.input, {
			userId: locals.user.id
		});

		return json(dataset, { status: 201 });
	} catch (error) {
		return datasetErrorResponse(error);
	}
};
