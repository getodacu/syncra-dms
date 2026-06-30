import { json } from "@sveltejs/kit";

import { deleteDataset, getDataset, isDatasetApiError, updateDataset } from "$lib/server/datasets";
import { jsonPublicErrorResponse } from "$lib/server/public-errors";
import { readUpdateDatasetRequest } from "../dataset-request";
import type { RequestHandler } from "./$types";

function datasetErrorResponse(error: unknown) {
	if (isDatasetApiError(error)) {
		return jsonPublicErrorResponse(error.status, error.message);
	}

	throw error;
}

export const GET: RequestHandler = async ({ fetch, locals, params }) => {
	if (!locals.user) {
		return json({ error: "authentication required" }, { status: 401 });
	}

	try {
		const dataset = await getDataset(fetch, params.id, { userId: locals.user.id });
		return json(dataset);
	} catch (error) {
		return datasetErrorResponse(error);
	}
};

export const PUT: RequestHandler = async ({ request, fetch, locals, params }) => {
	if (!locals.user) {
		return json({ error: "authentication required" }, { status: 401 });
	}

	const bodyResult = await readUpdateDatasetRequest(request);
	if (!bodyResult.ok) {
		return json({ error: bodyResult.error }, { status: 400 });
	}

	try {
		const dataset = await updateDataset(fetch, params.id, bodyResult.input, {
			userId: locals.user.id
		});
		return json(dataset);
	} catch (error) {
		return datasetErrorResponse(error);
	}
};

export const DELETE: RequestHandler = async ({ fetch, locals, params }) => {
	if (!locals.user) {
		return json({ error: "authentication required" }, { status: 401 });
	}

	try {
		await deleteDataset(fetch, params.id, { userId: locals.user.id });
		return new Response(null, { status: 204 });
	} catch (error) {
		return datasetErrorResponse(error);
	}
};
