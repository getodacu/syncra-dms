import { beforeEach, describe, expect, it, vi } from "vitest";

import { OCRApiError } from "$lib/server/ocr";
import { DELETE, PATCH } from "./+server";
import type { RequestEvent } from "./$types";

const { deleteOCRDocumentMock, updateOCRDocumentMock, OCRApiErrorMock } = vi.hoisted(() => {
	class MockOCRApiError extends Error {
		status: number;

		constructor(status: number, message: string) {
			super(message);
			this.name = "OCRApiError";
			this.status = status;
		}
	}

	return {
		deleteOCRDocumentMock: vi.fn(),
		updateOCRDocumentMock: vi.fn(),
		OCRApiErrorMock: MockOCRApiError
	};
});

vi.mock("$lib/server/ocr", () => ({
	deleteOCRDocument: deleteOCRDocumentMock,
	updateOCRDocument: updateOCRDocumentMock,
	OCRApiError: OCRApiErrorMock,
	isOCRApiError: (error: unknown) => error instanceof OCRApiErrorMock
}));

function createEvent(id = "document-1", user: unknown = { id: "user-1" }) {
	return {
		params: { id },
		fetch: vi.fn(),
		locals: { user }
	} as unknown as RequestEvent;
}

function createPatchEvent(
	body: unknown = { original_filename: "renamed.pdf" },
	id = "document-1",
	user: unknown = { id: "user-1" }
) {
	return {
		params: { id },
		request: new Request("http://localhost/api/ocr/documents/" + id, {
			method: "PATCH",
			headers: { "content-type": "application/json" },
			body: typeof body === "string" ? body : JSON.stringify(body)
		}),
		fetch: vi.fn(),
		locals: { user }
	} as unknown as RequestEvent;
}

async function responseJson(response: Response) {
	return (await response.json()) as unknown;
}

describe("OCR documents single-delete API endpoint", () => {
	beforeEach(() => {
		deleteOCRDocumentMock.mockReset();
		updateOCRDocumentMock.mockReset();
	});

	it("returns 401 for unauthenticated OCR document update requests", async () => {
		const response = await PATCH(createPatchEvent({ original_filename: "renamed.pdf" }, "document-1", null));

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(updateOCRDocumentMock).not.toHaveBeenCalled();
	});

	it("rejects invalid OCR document update payloads", async () => {
		const response = await PATCH(createPatchEvent({ name: "renamed.pdf" }));

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({
			error: "invalid OCR document update request"
		});
		expect(updateOCRDocumentMock).not.toHaveBeenCalled();
	});

	it("rejects malformed OCR document update JSON", async () => {
		const response = await PATCH(createPatchEvent("{"));

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({
			error: "invalid OCR document update request"
		});
		expect(updateOCRDocumentMock).not.toHaveBeenCalled();
	});

	it("updates an owned document", async () => {
		const result = {
			id: "document-1",
			created_at: "2026-05-27T00:00:00Z",
			updated_at: "2026-05-27T00:01:00Z",
			original_filename: "renamed.pdf",
			mime_type: "application/pdf",
			file_size: 1536,
			page_count: 2,
			document_hash: "abcdef",
			has_inline_schema: false,
			markdown: "# OCR",
			cached: false
		};
		updateOCRDocumentMock.mockResolvedValue(result);
		const event = createPatchEvent({ original_filename: "renamed.pdf" }, "document-1", { id: "user-1" });

		const response = await PATCH(event);

		expect(updateOCRDocumentMock).toHaveBeenCalledWith(
			event.fetch,
			"document-1",
			{ originalFilename: "renamed.pdf" },
			{ userId: "user-1" }
		);
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(result);
	});

	it("preserves backend document update 404 errors", async () => {
		updateOCRDocumentMock.mockRejectedValue(new OCRApiError(404, "OCR document not found"));

		const response = await PATCH(createPatchEvent());

		expect(response.status).toBe(404);
		expect(await responseJson(response)).toEqual({ error: "OCR document not found" });
	});

	it("normalizes backend document update server errors", async () => {
		updateOCRDocumentMock.mockRejectedValue(new OCRApiError(503, "OCR service unavailable"));

		const response = await PATCH(createPatchEvent());

		expect(response.status).toBe(502);
		expect(await responseJson(response)).toEqual({ error: "A server error occurred. Please try again." });
	});

	it("returns 401 for unauthenticated OCR document delete requests", async () => {
		const response = await DELETE(createEvent("document-1", null));

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(deleteOCRDocumentMock).not.toHaveBeenCalled();
	});

	it("deletes an owned document", async () => {
		const result = { deleted_ids: ["document-1"], deleted_count: 1 };
		deleteOCRDocumentMock.mockResolvedValue(result);
		const event = createEvent("document-1", { id: "user-1" });

		const response = await DELETE(event);

		expect(deleteOCRDocumentMock).toHaveBeenCalledWith(event.fetch, "document-1", {
			userId: "user-1"
		});
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(result);
	});

	it("preserves backend document delete 404 errors", async () => {
		deleteOCRDocumentMock.mockRejectedValue(new OCRApiError(404, "OCR document not found"));

		const response = await DELETE(createEvent());

		expect(response.status).toBe(404);
		expect(await responseJson(response)).toEqual({ error: "OCR document not found" });
	});

	it("normalizes backend document delete server errors", async () => {
		deleteOCRDocumentMock.mockRejectedValue(new OCRApiError(503, "OCR service unavailable"));

		const response = await DELETE(createEvent());

		expect(response.status).toBe(502);
		expect(await responseJson(response)).toEqual({ error: "A server error occurred. Please try again." });
	});
});
