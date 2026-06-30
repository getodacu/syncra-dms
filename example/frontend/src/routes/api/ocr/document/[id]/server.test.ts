import { beforeEach, describe, expect, it, vi } from "vitest";

import { OCRApiError } from "$lib/server/ocr";
import * as endpoint from "./+server";
import type { RequestEvent } from "./$types";

const { getOCRDocumentMock, OCRApiErrorMock } = vi.hoisted(() => {
	class MockOCRApiError extends Error {
		status: number;

		constructor(status: number, message: string) {
			super(message);
			this.name = "OCRApiError";
			this.status = status;
		}
	}

	return {
		getOCRDocumentMock: vi.fn(),
		OCRApiErrorMock: MockOCRApiError
	};
});

vi.mock("$lib/server/ocr", () => ({
	getOCRDocument: getOCRDocumentMock,
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

async function responseJson(response: Response) {
	return (await response.json()) as unknown;
}

function validDocument(overrides: Record<string, unknown> = {}) {
	return {
		id: "document-1",
		created_at: "2026-05-27T00:00:00Z",
		updated_at: "2026-05-27T00:00:00Z",
		user_id: "user-1",
		original_filename: "invoice.pdf",
		mime_type: "application/pdf",
		file_size: 12,
		document_hash: "abcdef",
		schema_id: "schema-1",
		has_inline_schema: false,
		markdown: "# Invoice",
		annotation_json: { total: 10 },
		cached: false,
		...overrides
	};
}

describe("OCR document API endpoint", () => {
	beforeEach(() => {
		getOCRDocumentMock.mockReset();
	});

	it("returns 401 for unauthenticated OCR document requests", async () => {
		const response = await endpoint.GET(createEvent("document-1", null));

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(getOCRDocumentMock).not.toHaveBeenCalled();
	});

	it("returns an owned document", async () => {
		const result = validDocument();
		getOCRDocumentMock.mockResolvedValue(result);
		const event = createEvent("document-1", { id: "user-1" });

		const response = await endpoint.GET(event);

		expect(getOCRDocumentMock).toHaveBeenCalledWith(event.fetch, "document-1", {
			userId: "user-1"
		});
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(result);
	});

	it("preserves backend 404 errors", async () => {
		getOCRDocumentMock.mockRejectedValue(new OCRApiError(404, "OCR document not found"));

		const response = await endpoint.GET(createEvent());

		expect(response.status).toBe(404);
		expect(await responseJson(response)).toEqual({ error: "OCR document not found" });
	});

	it("normalizes backend server errors", async () => {
		getOCRDocumentMock.mockRejectedValue(new OCRApiError(503, "OCR service unavailable"));

		const response = await endpoint.GET(createEvent());

		expect(response.status).toBe(502);
		expect(await responseJson(response)).toEqual({ error: "A server error occurred. Please try again." });
	});

	it("does not expose a single-delete handler on the legacy singular route", () => {
		expect("DELETE" in endpoint).toBe(false);
	});
});
