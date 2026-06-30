import { beforeEach, describe, expect, it, vi } from "vitest";

import { CollectionApiError } from "$lib/server/collections";
import { DELETE, PUT } from "./+server";
import type { RequestEvent } from "./$types";

const MAX_COLLECTION_REQUEST_BYTES = 1 << 20;

const { deleteCollectionMock, updateCollectionMock, CollectionApiErrorMock } = vi.hoisted(() => {
	class MockCollectionApiError extends Error {
		status: number;

		constructor(status: number, message: string) {
			super(message);
			this.name = "CollectionApiError";
			this.status = status;
		}
	}

	return {
		deleteCollectionMock: vi.fn(),
		updateCollectionMock: vi.fn(),
		CollectionApiErrorMock: MockCollectionApiError
	};
});

vi.mock("$lib/server/collections", () => ({
	deleteCollection: deleteCollectionMock,
	updateCollection: updateCollectionMock,
	CollectionApiError: CollectionApiErrorMock,
	isCollectionApiError: (error: unknown) => error instanceof CollectionApiErrorMock
}));

function createEvent(
	method: "DELETE" | "PUT",
	body: unknown = undefined,
	user: unknown = { id: "user-1" }
) {
	return createRequestEvent(
		new Request("http://localhost/api/collections/collection-1", {
			method,
			body: body === undefined ? undefined : JSON.stringify(body)
		}),
		user
	);
}

function createRequestEvent(request: Request, user: unknown = { id: "user-1" }) {
	return {
		request,
		fetch: vi.fn(),
		locals: {
			user
		},
		params: {
			id: "collection-1"
		}
	} as unknown as RequestEvent;
}

async function responseJson(response: Response) {
	return (await response.json()) as unknown;
}

describe("collection item API endpoint", () => {
	beforeEach(() => {
		deleteCollectionMock.mockReset();
		updateCollectionMock.mockReset();
	});

	it("returns 401 for unauthenticated PUT requests", async () => {
		const response = await PUT(
			createEvent(
				"PUT",
				{
					name: "Invoices",
					schema_ids: ["schema-1"]
				},
				null
			)
		);

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(updateCollectionMock).not.toHaveBeenCalled();
	});

	it("returns 401 for unauthenticated DELETE requests", async () => {
		const response = await DELETE(createEvent("DELETE", undefined, null));

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(deleteCollectionMock).not.toHaveBeenCalled();
	});

	it("returns 400 for invalid PUT bodies", async () => {
		const response = await PUT(
			createEvent("PUT", {
				name: "Invoices",
				schema_ids: ["schema-1", 42]
			})
		);

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({
			error: "schema_ids must be an array of strings"
		});
		expect(updateCollectionMock).not.toHaveBeenCalled();
	});

	it("returns 400 when PUT content-length exceeds the request limit", async () => {
		const response = await PUT(
			createRequestEvent(
				new Request("http://localhost/api/collections/collection-1", {
					method: "PUT",
					headers: { "content-length": String(MAX_COLLECTION_REQUEST_BYTES + 1) },
					body: "{}"
				})
			)
		);

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "request body too large" });
		expect(updateCollectionMock).not.toHaveBeenCalled();
	});

	it("calls the collection service with id, current user id, and input for PUT requests", async () => {
		const collection = {
			id: "collection-1",
			created_at: "2026-05-30T00:00:00Z",
			updated_at: "2026-05-30T00:00:00Z",
			user_id: "user-1",
			name: "Invoices",
			schema_ids: ["schema-1", "schema-2"],
			schema_count: 2,
			document_count: 0
		};
		updateCollectionMock.mockResolvedValue(collection);
		const event = createEvent("PUT", {
			name: "  Invoices  ",
			schema_ids: ["schema-1", "schema-2"]
		});

		const response = await PUT(event);

		expect(updateCollectionMock).toHaveBeenCalledWith(
			event.fetch,
			"collection-1",
			{
				name: "Invoices",
				schema_ids: ["schema-1", "schema-2"]
			},
			{ userId: "user-1" }
		);
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(collection);
	});

	it("calls the collection service with id and current user id for DELETE requests", async () => {
		deleteCollectionMock.mockResolvedValue({ deleted_id: "collection-1" });
		const event = createEvent("DELETE");

		const response = await DELETE(event);

		expect(deleteCollectionMock).toHaveBeenCalledWith(event.fetch, "collection-1", {
			userId: "user-1"
		});
		expect(response.status).toBe(204);
		expect(await response.text()).toBe("");
	});

	it("preserves collection service client errors", async () => {
		updateCollectionMock.mockRejectedValue(
			new CollectionApiError(404, "collection not found")
		);

		const response = await PUT(
			createEvent("PUT", {
				name: "Invoices",
				schema_ids: ["schema-1"]
			})
		);

		expect(response.status).toBe(404);
		expect(await responseJson(response)).toEqual({ error: "collection not found" });
	});

	it("normalizes collection service server errors to 502", async () => {
		deleteCollectionMock.mockRejectedValue(
			new CollectionApiError(503, "Collection service unavailable")
		);

		const response = await DELETE(createEvent("DELETE"));

		expect(response.status).toBe(502);
		expect(await responseJson(response)).toEqual({
			error: "A server error occurred. Please try again."
		});
	});
});
