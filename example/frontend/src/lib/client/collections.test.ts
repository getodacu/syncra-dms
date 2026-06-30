import { describe, expect, it, vi } from "vitest";

import {
	createCollection,
	deleteCollection,
	fetchCollections,
	updateCollection,
} from "./collections";

function jsonResponse(body: unknown, init?: ResponseInit) {
	return new Response(JSON.stringify(body), {
		headers: { "content-type": "application/json" },
		...init,
	});
}

function collectionResponse(id = "collection-1") {
	return {
		id,
		created_at: "2026-06-01T00:00:00Z",
		updated_at: "2026-06-01T00:01:00Z",
		user_id: "user-1",
		name: "Invoices",
		schema_ids: ["schema-1"],
		schema_count: 1,
		document_count: 2,
	};
}

function expectJsonContentType(init: RequestInit | undefined) {
	const headers = init?.headers as Headers;
	expect(headers.get("content-type")).toBe("application/json");
}

describe("collections api client", () => {
	it("fetches collections with default pagination through the SvelteKit proxy", async () => {
		const body = {
			collections: [collectionResponse()],
			next_cursor: null,
		};
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body));

		await expect(fetchCollections(fetchFn)).resolves.toEqual(body);
		expect(fetchFn).toHaveBeenCalledWith("/api/collections?size=100", { method: "GET" });
	});

	it("fetches collections with cursor pagination options", async () => {
		const body = {
			collections: [collectionResponse()],
			next_cursor: "cursor-2",
		};
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body));

		await expect(
			fetchCollections(fetchFn, { cursor: "cursor-1", size: 50, sort: "asc" })
		).resolves.toEqual(body);
		expect(fetchFn).toHaveBeenCalledWith(
			"/api/collections?cursor=cursor-1&size=50&sort=asc",
			{ method: "GET" }
		);
	});

	it("omits blank collection list options while keeping the default size", async () => {
		const body = {
			collections: [collectionResponse()],
			next_cursor: null,
		};
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body));

		await expect(fetchCollections(fetchFn, { cursor: "", sort: "   " })).resolves.toEqual(body);
		expect(fetchFn).toHaveBeenCalledWith("/api/collections?size=100", { method: "GET" });
	});

	it("creates collections through the SvelteKit proxy with JSON", async () => {
		const saved = collectionResponse();
		const input = { name: "Invoices", schema_ids: ["schema-1"] };
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(saved, { status: 201 }));

		await expect(createCollection(fetchFn, input)).resolves.toEqual(saved);
		expect(fetchFn).toHaveBeenCalledWith(
			"/api/collections",
			expect.objectContaining({
				method: "POST",
				body: JSON.stringify(input),
			})
		);
		expectJsonContentType(fetchFn.mock.calls[0]?.[1]);
	});

	it("updates collections through the SvelteKit proxy with encoded ids and JSON", async () => {
		const saved = { ...collectionResponse("collection/1"), name: "Receipts" };
		const input = { name: "Receipts", schema_ids: [] };
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(saved));

		await expect(updateCollection(fetchFn, "collection/1", input)).resolves.toEqual(saved);
		expect(fetchFn).toHaveBeenCalledWith(
			"/api/collections/collection%2F1",
			expect.objectContaining({
				method: "PUT",
				body: JSON.stringify(input),
			})
		);
		expectJsonContentType(fetchFn.mock.calls[0]?.[1]);
	});

	it("deletes collections with encoded ids and accepts 204 responses", async () => {
		const fetchFn = vi.fn().mockResolvedValue(new Response(null, { status: 204 }));

		await expect(deleteCollection(fetchFn, "collection/1")).resolves.toEqual({
			deleted_id: "collection/1",
		});
		expect(fetchFn).toHaveBeenCalledWith("/api/collections/collection%2F1", {
			method: "DELETE",
		});
	});

	it("throws backend JSON messages for list errors", async () => {
		const fetchFn = vi
			.fn()
			.mockResolvedValue(jsonResponse({ error: "invalid pagination size" }, { status: 400 }));

		await expect(fetchCollections(fetchFn)).rejects.toThrow("invalid pagination size");
	});

	it("throws backend JSON messages for mutation errors", async () => {
		const fetchFn = vi
			.fn()
			.mockResolvedValue(jsonResponse({ error: "collection name already exists" }, { status: 409 }));

		await expect(
			createCollection(fetchFn, { name: "Invoices", schema_ids: [] })
		).rejects.toThrow("collection name already exists");
	});

	it("rejects invalid list responses", async () => {
		const fetchFn = vi
			.fn()
			.mockResolvedValue(jsonResponse({ collections: [{ id: "collection-1" }], next_cursor: null }));

		await expect(fetchCollections(fetchFn)).rejects.toThrow("Invalid collection list response");
	});

	it("rejects invalid collection responses", async () => {
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse({ id: "collection-1" }));

		await expect(
			updateCollection(fetchFn, "collection-1", { name: "Invoices", schema_ids: [] })
		).rejects.toThrow("Invalid collection response");
	});

	it("rejects invalid delete responses", async () => {
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse({ deleted_id: 1 }));

		await expect(deleteCollection(fetchFn, "collection-1")).rejects.toThrow(
			"Invalid collection delete response"
		);
	});
});
