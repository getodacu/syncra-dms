import { json } from "@sveltejs/kit";

import {
	deleteCollection,
	isCollectionApiError,
	updateCollection
} from "$lib/server/collections";
import { jsonPublicErrorResponse } from "$lib/server/public-errors";
import { readCollectionRequest } from "../collection-request";
import type { RequestHandler } from "./$types";

function collectionErrorResponse(error: unknown) {
	if (isCollectionApiError(error)) {
		return jsonPublicErrorResponse(error.status, error.message);
	}

	throw error;
}

export const PUT: RequestHandler = async ({ request, fetch, locals, params }) => {
	if (!locals.user) {
		return json({ error: "authentication required" }, { status: 401 });
	}

	const bodyResult = await readCollectionRequest(request);
	if (!bodyResult.ok) {
		return json({ error: bodyResult.error }, { status: 400 });
	}

	try {
		const collection = await updateCollection(fetch, params.id, bodyResult.input, {
			userId: locals.user.id
		});
		return json(collection);
	} catch (error) {
		return collectionErrorResponse(error);
	}
};

export const DELETE: RequestHandler = async ({ fetch, locals, params }) => {
	if (!locals.user) {
		return json({ error: "authentication required" }, { status: 401 });
	}

	try {
		await deleteCollection(fetch, params.id, { userId: locals.user.id });
		return new Response(null, { status: 204 });
	} catch (error) {
		return collectionErrorResponse(error);
	}
};
