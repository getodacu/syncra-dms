import { beforeEach, describe, expect, it, vi } from "vitest";

import { OCRApiError } from "$lib/server/ocr";
import { DELETE, GET } from "./+server";
import type { RequestEvent } from "./$types";

const { deleteOCRDocumentsMock, listOCRDocumentsMock, OCRApiErrorMock } = vi.hoisted(() => {
	class MockOCRApiError extends Error {
		status: number;

		constructor(status: number, message: string) {
			super(message);
			this.name = "OCRApiError";
			this.status = status;
		}
	}

	return {
		deleteOCRDocumentsMock: vi.fn(),
		listOCRDocumentsMock: vi.fn(),
		OCRApiErrorMock: MockOCRApiError
	};
});

vi.mock("$lib/server/ocr", () => ({
	deleteOCRDocuments: deleteOCRDocumentsMock,
	listOCRDocuments: listOCRDocumentsMock,
	OCRApiError: OCRApiErrorMock,
	isOCRApiError: (error: unknown) => error instanceof OCRApiErrorMock
}));

function createEvent(url: string, user: unknown = { id: "user-1" }, requestInit?: RequestInit) {
	return {
		url: new URL(url),
		request: new Request(url, requestInit),
		fetch: vi.fn(),
		locals: { user }
	} as unknown as RequestEvent;
}

async function responseJson(response: Response) {
	return (await response.json()) as unknown;
}

function validDocumentList() {
	return {
		documents: [
			{
				id: "document-1",
				created_at: "2026-05-27T00:00:00Z",
				updated_at: "2026-05-27T00:01:00Z",
				user_id: "user-1",
				original_filename: "invoice.pdf",
				mime_type: "application/pdf",
				file_size: 12,
				page_count: 2,
				document_hash: "abcdef",
				schema_id: "schema-1",
				has_inline_schema: false
			}
		],
		next_cursor: null
	};
}

describe("OCR documents API endpoint", () => {
	beforeEach(() => {
		deleteOCRDocumentsMock.mockReset();
		listOCRDocumentsMock.mockReset();
	});

	it("returns 401 for unauthenticated document list requests", async () => {
		const response = await GET(createEvent("http://localhost/api/ocr/documents", null));

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(listOCRDocumentsMock).not.toHaveBeenCalled();
	});

	it("forwards filters and injects the authenticated user id", async () => {
		const result = validDocumentList();
		listOCRDocumentsMock.mockResolvedValue(result);
		const event = createEvent(
			"http://localhost/api/ocr/documents?user_id=attacker&collection=collection-1&filename=invoice&created_from=2026-05-27T00%3A00%3A00Z&created_to=2026-05-28T00%3A00%3A00Z&cursor=cursor-1&size=50&sort=asc",
			{ id: "user-1" }
		);

		const response = await GET(event);

		expect(listOCRDocumentsMock).toHaveBeenCalledWith(event.fetch, {
			userId: "user-1",
			collectionId: "collection-1",
			filename: "invoice",
			createdFrom: "2026-05-27T00:00:00Z",
			createdTo: "2026-05-28T00:00:00Z",
			cursor: "cursor-1",
			size: "50",
			sort: "asc"
		});
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(result);
	});

	it("preserves OCR document list client errors", async () => {
		listOCRDocumentsMock.mockRejectedValue(new OCRApiError(400, "invalid cursor"));

		const response = await GET(createEvent("http://localhost/api/ocr/documents?cursor=bad"));

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "invalid cursor" });
	});

	it("normalizes OCR document list server errors", async () => {
		listOCRDocumentsMock.mockRejectedValue(new OCRApiError(503, "OCR service unavailable"));

		const response = await GET(createEvent("http://localhost/api/ocr/documents"));

		expect(response.status).toBe(502);
		expect(await responseJson(response)).toEqual({ error: "A server error occurred. Please try again." });
	});

	it("returns 401 for unauthenticated bulk document delete requests", async () => {
		const response = await DELETE(
			createEvent("http://localhost/api/ocr/documents", null, {
				method: "DELETE",
				body: JSON.stringify({ ids: ["document-1"] })
			})
		);

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(deleteOCRDocumentsMock).not.toHaveBeenCalled();
	});

	it("bulk deletes documents and injects the authenticated user id", async () => {
		const result = { deleted_ids: ["document-1", "document-2"], deleted_count: 2 };
		deleteOCRDocumentsMock.mockResolvedValue(result);
		const event = createEvent("http://localhost/api/ocr/documents?user_id=attacker", { id: "user-1" }, {
			method: "DELETE",
			headers: { "content-type": "application/json" },
			body: JSON.stringify({ ids: ["document-1", "document-2"] })
		});

		const response = await DELETE(event);

		expect(deleteOCRDocumentsMock).toHaveBeenCalledWith(
			event.fetch,
			["document-1", "document-2"],
			{ userId: "user-1" }
		);
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(result);
	});

	it("rejects invalid bulk document delete request bodies", async () => {
		const response = await DELETE(
			createEvent("http://localhost/api/ocr/documents", { id: "user-1" }, {
				method: "DELETE",
				headers: { "content-type": "application/json" },
				body: JSON.stringify({ ids: [1] })
			})
		);

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({
			error: "invalid OCR document delete request"
		});
		expect(deleteOCRDocumentsMock).not.toHaveBeenCalled();
	});

	it("preserves OCR bulk document delete client errors", async () => {
		deleteOCRDocumentsMock.mockRejectedValue(new OCRApiError(400, "ids is required"));

		const response = await DELETE(
			createEvent("http://localhost/api/ocr/documents", { id: "user-1" }, {
				method: "DELETE",
				headers: { "content-type": "application/json" },
				body: JSON.stringify({ ids: [] })
			})
		);

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "ids is required" });
	});

	it("normalizes OCR bulk document delete server errors", async () => {
		deleteOCRDocumentsMock.mockRejectedValue(new OCRApiError(503, "OCR service unavailable"));

		const response = await DELETE(
			createEvent("http://localhost/api/ocr/documents", { id: "user-1" }, {
				method: "DELETE",
				headers: { "content-type": "application/json" },
				body: JSON.stringify({ ids: ["document-1"] })
			})
		);

		expect(response.status).toBe(502);
		expect(await responseJson(response)).toEqual({ error: "A server error occurred. Please try again." });
	});
});
