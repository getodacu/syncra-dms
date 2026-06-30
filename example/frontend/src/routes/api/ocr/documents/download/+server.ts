import { json } from "@sveltejs/kit";

import { getOCRDocument, isOCRApiError } from "$lib/server/ocr";
import { jsonPublicErrorResponse } from "$lib/server/public-errors";
import type { RequestHandler } from "./$types";
import {
	buildDownloadFiles,
	buildZip,
	contentDispositionAttachment,
	isDownloadFormat,
	zipDownloadFilename,
	type DownloadFormat
} from "./export-utils";

type DownloadRequest = {
	ids: string[];
	format: DownloadFormat;
};

function isDownloadRequest(value: unknown): value is DownloadRequest {
	if (typeof value !== "object" || value === null) return false;
	if (!("ids" in value) || !Array.isArray(value.ids)) return false;
	if (!value.ids.every((id) => typeof id === "string" && id.trim())) return false;

	return "format" in value && isDownloadFormat(value.format);
}

export const POST: RequestHandler = async ({ request, fetch, locals }) => {
	if (!locals.user) {
		return json({ error: "authentication required" }, { status: 401 });
	}
	const userId = locals.user.id;

	let body: unknown;
	try {
		body = await request.json();
	} catch {
		return json({ error: "invalid OCR document download request" }, { status: 400 });
	}

	if (!isDownloadRequest(body) || body.ids.length === 0) {
		return json({ error: "invalid OCR document download request" }, { status: 400 });
	}

	try {
		const documents = await Promise.all(
			body.ids.map((id) => getOCRDocument(fetch, id.trim(), { userId }))
		);
		const files = buildDownloadFiles(documents, body.format);

		if (files.length === 0) {
			return json({ error: "No schema-backed JSON is available for the selected documents" }, { status: 400 });
		}

		if (body.ids.length === 1) {
			const file = files[0];

			return new Response(file.content, {
				headers: {
					"content-type": file.contentType,
					"content-disposition": contentDispositionAttachment(file.filename)
				}
			});
		}

		const zipBytes = await buildZip(files);
		const filename = zipDownloadFilename();

		return new Response(uint8ArrayBody(zipBytes), {
			headers: {
				"content-type": "application/zip",
				"content-disposition": contentDispositionAttachment(filename)
			}
		});
	} catch (error) {
		if (isOCRApiError(error)) {
			return jsonPublicErrorResponse(error.status, error.message);
		}

		throw error;
	}
};

function uint8ArrayBody(bytes: Uint8Array): ArrayBuffer {
	const copy = new Uint8Array(bytes.byteLength);
	copy.set(bytes);

	return copy.buffer;
}
