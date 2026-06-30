import { json } from "@sveltejs/kit";

import { createOCRJob, deleteOCRJobs, isOCRApiError, listOCRJobs } from "$lib/server/ocr";
import { jsonPublicErrorResponse } from "$lib/server/public-errors";
import { isSchemaApiError, listSchemas } from "$lib/server/schemas";
import type { RequestHandler } from "./$types";

const MAX_OCR_FILE_BYTES = 20 << 20;
const MULTIPART_OVERHEAD_ALLOWANCE_BYTES = 64 * 1024;
const MAX_OCR_REQUEST_BYTES = MAX_OCR_FILE_BYTES + MULTIPART_OVERHEAD_ALLOWANCE_BYTES;

function optionalText(data: FormData, key: string) {
	const value = data.get(key);
	return typeof value === "string" && value.trim() ? value.trim() : undefined;
}

function optionalQuery(url: URL, key: string) {
	const value = url.searchParams.get(key);
	return value && value.trim() ? value.trim() : undefined;
}

function isDeleteOCRJobsRequest(value: unknown): value is { ids: string[] } {
	return (
		typeof value === "object" &&
		value !== null &&
		"ids" in value &&
		Array.isArray(value.ids) &&
		value.ids.every((id) => typeof id === "string")
	);
}

function ocrErrorResponse(error: unknown) {
	if (isOCRApiError(error)) {
		return jsonPublicErrorResponse(error.status, error.message);
	}

	throw error;
}

function requestExceedsUploadLimit(request: Request) {
	const contentLength = request.headers.get("content-length");
	if (!contentLength) return false;

	const size = Number(contentLength);
	return Number.isFinite(size) && size > MAX_OCR_REQUEST_BYTES;
}

async function readLimitedBody(request: Request) {
	if (requestExceedsUploadLimit(request)) return null;
	if (!request.body) return new Uint8Array();

	const reader = request.body.getReader();
	const chunks: Uint8Array[] = [];
	let bytes = 0;

	try {
		while (true) {
			const { done, value } = await reader.read();
			if (done) break;

			bytes += value.byteLength;
			if (bytes > MAX_OCR_REQUEST_BYTES) {
				await reader.cancel();
				return null;
			}

			chunks.push(value);
		}
	} finally {
		reader.releaseLock();
	}

	const body = new Uint8Array(bytes);
	let offset = 0;
	for (const chunk of chunks) {
		body.set(chunk, offset);
		offset += chunk.byteLength;
	}

	return body;
}

export const GET: RequestHandler = async ({ url, fetch, locals }) => {
	if (!locals.user) {
		return json({ error: "authentication required" }, { status: 401 });
	}

	try {
		const result = await listOCRJobs(fetch, {
			userId: locals.user.id,
			status: optionalQuery(url, "status"),
			cursor: optionalQuery(url, "cursor"),
			size: optionalQuery(url, "size"),
			sort: optionalQuery(url, "sort")
		});
		return json(result);
	} catch (error) {
		return ocrErrorResponse(error);
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
		return json({ error: "invalid OCR job delete request" }, { status: 400 });
	}

	if (!isDeleteOCRJobsRequest(body)) {
		return json({ error: "invalid OCR job delete request" }, { status: 400 });
	}

	try {
		const result = await deleteOCRJobs(fetch, body.ids, { userId: locals.user.id });
		return json(result);
	} catch (error) {
		return ocrErrorResponse(error);
	}
};

export const POST: RequestHandler = async ({ request, fetch, locals }) => {
	if (!locals.user) {
		return json({ error: "authentication required" }, { status: 401 });
	}

	const limitedBody = await readLimitedBody(request);
	if (limitedBody === null) {
		return json({ error: "request body too large" }, { status: 400 });
	}

	let data: FormData;
	try {
		data = await new Request(request.url, {
			method: "POST",
			headers: request.headers,
			body: limitedBody
		}).formData();
	} catch {
		return json({ error: "invalid form data" }, { status: 400 });
	}

	const file = data.get("file");
	if (!(file instanceof File)) {
		return json({ error: "file is required" }, { status: 400 });
	}
	if (file.size > MAX_OCR_FILE_BYTES) {
		return json({ error: "file exceeds max upload size" }, { status: 400 });
	}

	const schemaId = optionalText(data, "schema_id");

	try {
		if (schemaId) {
			const userSchemas = await listSchemas(fetch, { userId: locals.user.id });
			if (!userSchemas.some((schema) => schema.id === schemaId)) {
				return json({ error: "invalid schema_id" }, { status: 400 });
			}
		}

		const result = await createOCRJob(fetch, {
			file,
			schemaId,
			userId: locals.user.id
		});

		return json(result, { status: 202 });
	} catch (error) {
		if (isSchemaApiError(error)) {
			return jsonPublicErrorResponse(error.status, error.message);
		}

		if (isOCRApiError(error)) {
			return jsonPublicErrorResponse(error.status, error.message);
		}

		throw error;
	}
};
