import { json } from "@sveltejs/kit";

import { isOCRApiError, moveOCRDocumentsToCollections } from "$lib/server/ocr";
import { jsonPublicErrorResponse } from "$lib/server/public-errors";
import type { RequestHandler } from "./$types";

function isMoveOCRDocumentsRequest(value: unknown): value is {
	ids: string[];
	collection_ids?: string[];
} {
	return (
		typeof value === "object" &&
		value !== null &&
		"ids" in value &&
		Array.isArray(value.ids) &&
		value.ids.every((id) => typeof id === "string") &&
		(!("collection_ids" in value) ||
			(Array.isArray(value.collection_ids) &&
				value.collection_ids.every((id) => typeof id === "string")))
	);
}

export const PUT: RequestHandler = async ({ request, fetch, locals }) => {
	if (!locals.user) {
		return json({ error: "authentication required" }, { status: 401 });
	}

	let body: unknown;
	try {
		body = await request.json();
	} catch {
		return json({ error: "invalid OCR document move request" }, { status: 400 });
	}

	if (!isMoveOCRDocumentsRequest(body)) {
		return json({ error: "invalid OCR document move request" }, { status: 400 });
	}

	try {
		const result = await moveOCRDocumentsToCollections(
			fetch,
			body.ids,
			body.collection_ids ?? [],
			{ userId: locals.user.id }
		);
		return json(result);
	} catch (error) {
		if (isOCRApiError(error)) {
			return jsonPublicErrorResponse(error.status, error.message);
		}

		throw error;
	}
};
