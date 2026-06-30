import { json } from "@sveltejs/kit";

import { isSchemaApiError, listJsonRecipes } from "$lib/server/schemas";
import { jsonPublicErrorResponse } from "$lib/server/public-errors";
import type { RequestHandler } from "./$types";

function optionalQuery(url: URL, key: string) {
	const value = url.searchParams.get(key);
	return value && value.trim() ? value.trim() : undefined;
}

export const GET: RequestHandler = async ({ url, fetch }) => {
	try {
		const result = await listJsonRecipes(fetch, {
			cursor: optionalQuery(url, "cursor"),
			size: optionalQuery(url, "size"),
			sort: optionalQuery(url, "sort")
		});
		return json(result);
	} catch (error) {
		if (isSchemaApiError(error)) {
			return jsonPublicErrorResponse(error.status, error.message, "Failed to load OCR recipes");
		}
		throw error;
	}
};
