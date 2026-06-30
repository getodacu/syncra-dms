import { json } from "@sveltejs/kit";

import { deleteOCRDocument, isOCRApiError, updateOCRDocument } from "$lib/server/ocr";
import { jsonPublicErrorResponse } from "$lib/server/public-errors";
import type { RequestHandler } from "./$types";

function isUpdateOCRDocumentRequest(value: unknown): value is { original_filename: string } {
	return (
		typeof value === "object" &&
		value !== null &&
		!Array.isArray(value) &&
		"original_filename" in value &&
		typeof value.original_filename === "string"
	);
}

export const PATCH: RequestHandler = async ({ params, request, fetch, locals }) => {
	if (!locals.user) {
		return json({ error: "authentication required" }, { status: 401 });
	}

	let body: unknown;
	try {
		body = await request.json();
	} catch {
		return json({ error: "invalid OCR document update request" }, { status: 400 });
	}

	if (!isUpdateOCRDocumentRequest(body)) {
		return json({ error: "invalid OCR document update request" }, { status: 400 });
	}

	try {
		const result = await updateOCRDocument(
			fetch,
			params.id,
			{ originalFilename: body.original_filename },
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

export const DELETE: RequestHandler = async ({ params, fetch, locals }) => {
	if (!locals.user) {
		return json({ error: "authentication required" }, { status: 401 });
	}

	try {
		const result = await deleteOCRDocument(fetch, params.id, { userId: locals.user.id });
		return json(result);
	} catch (error) {
		if (isOCRApiError(error)) {
			return jsonPublicErrorResponse(error.status, error.message);
		}

		throw error;
	}
};
