import { json } from "@sveltejs/kit";

import {
	deleteOCRDocuments,
	isOCRApiError,
	listOCRDocuments,
	type ListOCRDocumentsOptions
} from "$lib/server/ocr";
import { jsonPublicErrorResponse } from "$lib/server/public-errors";
import type { RequestHandler } from "./$types";

function optionalQuery(url: URL, key: string) {
	const value = url.searchParams.get(key);
	return value && value.trim() ? value.trim() : undefined;
}

function isDeleteOCRDocumentsRequest(value: unknown): value is { ids: string[] } {
	return (
		typeof value === "object" &&
		value !== null &&
		"ids" in value &&
		Array.isArray(value.ids) &&
		value.ids.every((id) => typeof id === "string")
	);
}

export const GET: RequestHandler = async ({ url, fetch, locals }) => {
	if (!locals.user) {
		return json({ error: "authentication required" }, { status: 401 });
	}

	const options: ListOCRDocumentsOptions = {
		userId: locals.user.id,
		collectionId: optionalQuery(url, "collection"),
		schemaId: optionalQuery(url, "schema_id"),
		filename: optionalQuery(url, "filename"),
		createdFrom: optionalQuery(url, "created_from"),
		createdTo: optionalQuery(url, "created_to"),
		cursor: optionalQuery(url, "cursor"),
		size: optionalQuery(url, "size"),
		sort: optionalQuery(url, "sort")
	};

	try {
		const result = await listOCRDocuments(fetch, options);
		return json(result);
	} catch (error) {
		if (isOCRApiError(error)) {
			return jsonPublicErrorResponse(error.status, error.message);
		}

		throw error;
	}
};

export const DELETE: RequestHandler = async ({ request, fetch, locals }) => {
	if (!locals.user) {
		return json({ error: "authentication required" }, { status: 401 });
	}

	let body: unknown;
	try {
		body = await request.json();
	} catch {
		return json({ error: "invalid OCR document delete request" }, { status: 400 });
	}

	if (!isDeleteOCRDocumentsRequest(body)) {
		return json({ error: "invalid OCR document delete request" }, { status: 400 });
	}

	try {
		const result = await deleteOCRDocuments(fetch, body.ids, { userId: locals.user.id });
		return json(result);
	} catch (error) {
		if (isOCRApiError(error)) {
			return jsonPublicErrorResponse(error.status, error.message);
		}

		throw error;
	}
};
