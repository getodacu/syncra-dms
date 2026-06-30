import { json } from "@sveltejs/kit";

import {
	createCollection,
	isCollectionApiError,
	listCollectionsPage
} from "$lib/server/collections";
import { jsonPublicErrorResponse } from "$lib/server/public-errors";
import { readCollectionRequest } from "./collection-request";
import type { RequestHandler } from "./$types";

function collectionErrorResponse(error: unknown) {
	if (isCollectionApiError(error)) {
		return jsonPublicErrorResponse(error.status, error.message);
	}

	throw error;
}

export const GET: RequestHandler = async ({ url, fetch, locals }) => {
	if (!locals.user) {
		return json({ error: "authentication required" }, { status: 401 });
	}

	try {
		const page = await listCollectionsPage(fetch, {
			userId: locals.user.id,
			...(url.searchParams.has("cursor")
				? { cursor: url.searchParams.get("cursor") ?? "" }
				: {}),
			...(url.searchParams.has("size") ? { size: url.searchParams.get("size") ?? "" } : {}),
			...(url.searchParams.has("sort") ? { sort: url.searchParams.get("sort") ?? "" } : {})
		});

		return json(page);
	} catch (error) {
		return collectionErrorResponse(error);
	}
};

export const POST: RequestHandler = async ({ request, fetch, locals }) => {
	if (!locals.user) {
		return json({ error: "authentication required" }, { status: 401 });
	}

	const bodyResult = await readCollectionRequest(request);
	if (!bodyResult.ok) {
		return json({ error: bodyResult.error }, { status: 400 });
	}

	try {
		const collection = await createCollection(fetch, bodyResult.input, {
			userId: locals.user.id
		});

		return json(collection, { status: 201 });
	} catch (error) {
		return collectionErrorResponse(error);
	}
};
