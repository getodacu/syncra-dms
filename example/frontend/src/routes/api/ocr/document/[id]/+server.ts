import { json } from "@sveltejs/kit";

import { getOCRDocument, isOCRApiError } from "$lib/server/ocr";
import { jsonPublicErrorResponse } from "$lib/server/public-errors";
import type { RequestHandler } from "./$types";

export const GET: RequestHandler = async ({ params, fetch, locals }) => {
	if (!locals.user) {
		return json({ error: "authentication required" }, { status: 401 });
	}

	try {
		const result = await getOCRDocument(fetch, params.id, { userId: locals.user.id });
		return json(result);
	} catch (error) {
		if (isOCRApiError(error)) {
			return jsonPublicErrorResponse(error.status, error.message);
		}

		throw error;
	}
};
