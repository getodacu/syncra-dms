import { beforeEach, describe, expect, it, vi } from "vitest";

import { CollectionApiError } from "$lib/server/collections";
import { GET, POST } from "./+server";
import type { RequestEvent } from "./$types";

const MAX_COLLECTION_REQUEST_BYTES = 1 << 20;

const { createCollectionMock, listCollectionsPageMock, CollectionApiErrorMock } = vi.hoisted(
	() => {
		class MockCollectionApiError extends Error {
			status: number;

			constructor(status: number, message: string) {
				super(message);
				this.name = "CollectionApiError";
				this.status = status;
			}
		}

		return {
			createCollectionMock: vi.fn(),
			listCollectionsPageMock: vi.fn(),
			CollectionApiErrorMock: MockCollectionApiError
		};
	}
);

vi.mock("$lib/server/collections", () => ({
	createCollection: createCollectionMock,
	listCollectionsPage: listCollectionsPageMock,
	CollectionApiError: CollectionApiErrorMock,
	isCollectionApiError: (error: unknown) => error instanceof CollectionApiErrorMock
}));

function createPostEvent(body: unknown, user: unknown = { id: "user-1" }) {
	return createPostRequestEvent(
		new Request("http://localhost/api/collections", {
			method: "POST",
			body: body === undefined ? undefined : JSON.stringify(body)
		}),
		user
	);
}

function createPostRequestEvent(request: Request, user: unknown = { id: "user-1" }) {
	return {
		request,
		fetch: vi.fn(),
		locals: {
			user
		}
	} as unknown as RequestEvent;
}

function createGetEvent(url: string, user: unknown = { id: "user-1" }) {
	return {
		url: new URL(url),
		fetch: vi.fn(),
		locals: {
			user
		}
	} as unknown as RequestEvent;
}

async function responseJson(response: Response) {
	return (await response.json()) as unknown;
}

describe("collection API endpoint", () => {
	beforeEach(() => {
		createCollectionMock.mockReset();
		listCollectionsPageMock.mockReset();
	});

	it("returns 401 for unauthenticated GET requests", async () => {
		const response = await GET(createGetEvent("http://localhost/api/collections", null));

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(listCollectionsPageMock).not.toHaveBeenCalled();
	});

	it("returns 401 for unauthenticated POST requests", async () => {
		const response = await POST(
			createPostEvent(
				{
					name: "Invoices",
					schema_ids: ["schema-1"]
				},
				null
			)
		);

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(createCollectionMock).not.toHaveBeenCalled();
	});

	it("calls the collection service with current user id and pagination parameters", async () => {
		const page = {
			collections: [
				{
					id: "collection-1",
					created_at: "2026-05-30T00:00:00Z",
					updated_at: "2026-05-30T00:00:00Z",
					user_id: "user-1",
					name: "Invoices",
					schema_ids: ["schema-1"],
					schema_count: 1,
					document_count: 0
				}
			],
			next_cursor: "next-page"
		};
		listCollectionsPageMock.mockResolvedValue(page);
		const event = createGetEvent(
			"http://localhost/api/collections?size=25&sort=desc&cursor=abc"
		);

		const response = await GET(event);

		expect(listCollectionsPageMock).toHaveBeenCalledWith(event.fetch, {
			userId: "user-1",
			cursor: "abc",
			size: "25",
			sort: "desc"
		});
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(page);
	});

	it("returns 400 for invalid POST bodies", async () => {
		const response = await POST(
			createPostEvent({
				name: "Invoices",
				schema_ids: ["schema-1", 42]
			})
		);

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({
			error: "schema_ids must be an array of strings"
		});
		expect(createCollectionMock).not.toHaveBeenCalled();
	});

	it("returns 400 when POST content-length exceeds the request limit", async () => {
		const response = await POST(
			createPostRequestEvent(
				new Request("http://localhost/api/collections", {
					method: "POST",
					headers: { "content-length": String(MAX_COLLECTION_REQUEST_BYTES + 1) },
					body: "{}"
				})
			)
		);

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "request body too large" });
		expect(createCollectionMock).not.toHaveBeenCalled();
	});

	it("calls the collection service with event fetch, current user id, and input", async () => {
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
		createCollectionMock.mockResolvedValue(collection);
		const event = createPostEvent({
			name: "  Invoices  ",
			schema_ids: ["schema-1", "schema-2"]
		});

		const response = await POST(event);

		expect(createCollectionMock).toHaveBeenCalledWith(
			event.fetch,
			{
				name: "Invoices",
				schema_ids: ["schema-1", "schema-2"]
			},
			{ userId: "user-1" }
		);
		expect(response.status).toBe(201);
		expect(await responseJson(response)).toEqual(collection);
	});

	it("preserves collection service client errors", async () => {
		listCollectionsPageMock.mockRejectedValue(
			new CollectionApiError(404, "collection not found")
		);

		const response = await GET(createGetEvent("http://localhost/api/collections"));

		expect(response.status).toBe(404);
		expect(await responseJson(response)).toEqual({ error: "collection not found" });
	});

	it("normalizes collection service server errors to 502", async () => {
		createCollectionMock.mockRejectedValue(
			new CollectionApiError(503, "Collection service unavailable")
		);

		const response = await POST(
			createPostEvent({
				name: "Invoices",
				schema_ids: ["schema-1"]
			})
		);

		expect(response.status).toBe(502);
		expect(await responseJson(response)).toEqual({
			error: "A server error occurred. Please try again."
		});
	});
});
