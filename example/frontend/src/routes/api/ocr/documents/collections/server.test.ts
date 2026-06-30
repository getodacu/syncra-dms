import { beforeEach, describe, expect, it, vi } from "vitest";

import { OCRApiError } from "$lib/server/ocr";
import { PUT } from "./+server";
import type { RequestEvent } from "./$types";

const { moveOCRDocumentsToCollectionsMock, OCRApiErrorMock } = vi.hoisted(() => {
	class MockOCRApiError extends Error {
		status: number;

		constructor(status: number, message: string) {
			super(message);
			this.name = "OCRApiError";
			this.status = status;
		}
	}

	return {
		moveOCRDocumentsToCollectionsMock: vi.fn(),
		OCRApiErrorMock: MockOCRApiError
	};
});

vi.mock("$lib/server/ocr", () => ({
	moveOCRDocumentsToCollections: moveOCRDocumentsToCollectionsMock,
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

describe("OCR document collection move API endpoint", () => {
	beforeEach(() => {
		moveOCRDocumentsToCollectionsMock.mockReset();
	});

	it("returns 401 for unauthenticated document move requests", async () => {
		const response = await PUT(
			createEvent("http://localhost/api/ocr/documents/collections", null, {
				method: "PUT",
				body: JSON.stringify({ ids: ["document-1"], collection_ids: ["collection-1"] })
			})
		);

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(moveOCRDocumentsToCollectionsMock).not.toHaveBeenCalled();
	});

	it("moves documents and injects the authenticated user id", async () => {
		const result = {
			moved_ids: ["document-1", "document-2"],
			moved_count: 2,
			collection_ids: ["collection-1"]
		};
		moveOCRDocumentsToCollectionsMock.mockResolvedValue(result);
		const event = createEvent(
			"http://localhost/api/ocr/documents/collections?user_id=attacker",
			{ id: "user-1" },
			{
				method: "PUT",
				headers: { "content-type": "application/json" },
				body: JSON.stringify({
					ids: ["document-1", "document-2"],
					collection_ids: ["collection-1"]
				})
			}
		);

		const response = await PUT(event);

		expect(moveOCRDocumentsToCollectionsMock).toHaveBeenCalledWith(
			event.fetch,
			["document-1", "document-2"],
			["collection-1"],
			{ userId: "user-1" }
		);
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(result);
	});

	it("defaults omitted collection ids to an empty move target", async () => {
		const result = { moved_ids: ["document-1"], moved_count: 1, collection_ids: [] };
		moveOCRDocumentsToCollectionsMock.mockResolvedValue(result);
		const event = createEvent("http://localhost/api/ocr/documents/collections", { id: "user-1" }, {
			method: "PUT",
			headers: { "content-type": "application/json" },
			body: JSON.stringify({ ids: ["document-1"] })
		});

		const response = await PUT(event);

		expect(moveOCRDocumentsToCollectionsMock).toHaveBeenCalledWith(
			event.fetch,
			["document-1"],
			[],
			{ userId: "user-1" }
		);
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(result);
	});

	it("rejects invalid document move request bodies", async () => {
		const response = await PUT(
			createEvent("http://localhost/api/ocr/documents/collections", { id: "user-1" }, {
				method: "PUT",
				headers: { "content-type": "application/json" },
				body: JSON.stringify({ ids: ["document-1"], collection_ids: [1] })
			})
		);

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({
			error: "invalid OCR document move request"
		});
		expect(moveOCRDocumentsToCollectionsMock).not.toHaveBeenCalled();
	});

	it("preserves OCR document move client errors", async () => {
		moveOCRDocumentsToCollectionsMock.mockRejectedValue(
			new OCRApiError(404, "collection not found")
		);

		const response = await PUT(
			createEvent("http://localhost/api/ocr/documents/collections", { id: "user-1" }, {
				method: "PUT",
				headers: { "content-type": "application/json" },
				body: JSON.stringify({ ids: ["document-1"], collection_ids: ["collection-1"] })
			})
		);

		expect(response.status).toBe(404);
		expect(await responseJson(response)).toEqual({ error: "collection not found" });
	});

	it("normalizes OCR document move server errors", async () => {
		moveOCRDocumentsToCollectionsMock.mockRejectedValue(
			new OCRApiError(503, "OCR service unavailable")
		);

		const response = await PUT(
			createEvent("http://localhost/api/ocr/documents/collections", { id: "user-1" }, {
				method: "PUT",
				headers: { "content-type": "application/json" },
				body: JSON.stringify({ ids: ["document-1"] })
			})
		);

		expect(response.status).toBe(502);
		expect(await responseJson(response)).toEqual({ error: "A server error occurred. Please try again." });
	});
});
