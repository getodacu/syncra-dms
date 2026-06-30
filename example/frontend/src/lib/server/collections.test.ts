import { beforeEach, describe, expect, it, vi } from "vitest";

const INTERNAL_API_HEADER = "X-Syncra-Internal-Token";

const privateEnv = vi.hoisted(() => ({}) as Record<string, string | undefined>);

vi.mock("$env/dynamic/private", () => ({ env: privateEnv }));

function collectionFixture(id = "collection-1") {
	return {
		id,
		created_at: "2026-06-01T00:00:00Z",
		updated_at: "2026-06-01T00:01:00Z",
		user_id: "user-1",
		name: "Invoices",
		schema_ids: ["schema-1"],
		schema_count: 1,
		document_count: 2
	};
}

function jsonResponse(body: unknown, init?: ResponseInit) {
	return new Response(JSON.stringify(body), {
		headers: { "content-type": "application/json" },
		...init
	});
}

function expectJsonContentType(init: RequestInit | undefined) {
	expect(new Headers(init?.headers).get("content-type")).toBe("application/json");
}

function expectInternalToken(init: RequestInit | undefined) {
	expect(new Headers(init?.headers).get(INTERNAL_API_HEADER)).toBe("internal-token");
}

describe("frontend collection server helper", () => {
	beforeEach(() => {
		privateEnv.SYNCRA_API_BASE_URL = "http://collection-api.test/";
		privateEnv.SYNCRA_INTERNAL_API_TOKEN = "internal-token";
	});

	it("lists user collections through the backend", async () => {
		const { listCollectionsPage } = await import("./collections");
		const page = { collections: [collectionFixture()], next_cursor: null };
		const fetchMock = vi.fn().mockResolvedValue(jsonResponse(page));

		await expect(listCollectionsPage(fetchMock, { userId: "user-1", size: 20 })).resolves.toEqual(
			page
		);
		expect(fetchMock).toHaveBeenCalledWith(
			"http://collection-api.test/api/collections?user_id=user-1&size=20",
			expect.objectContaining({ method: "GET" })
		);
		expectInternalToken(fetchMock.mock.calls[0]?.[1]);
	});

	it("creates collections with user id and schema ids", async () => {
		const { createCollection } = await import("./collections");
		const saved = collectionFixture();
		const fetchMock = vi.fn().mockResolvedValue(jsonResponse(saved, { status: 201 }));

		await expect(
			createCollection(fetchMock, { name: "Invoices", schema_ids: ["schema-1"] }, { userId: "user-1" })
		).resolves.toEqual(saved);
		expect(fetchMock).toHaveBeenCalledWith(
			"http://collection-api.test/api/collection",
			expect.objectContaining({
				method: "POST",
				body: JSON.stringify({ name: "Invoices", schema_ids: ["schema-1"], user_id: "user-1" })
			})
		);
		expectJsonContentType(fetchMock.mock.calls[0]?.[1]);
		expectInternalToken(fetchMock.mock.calls[0]?.[1]);
	});

	it("updates collections by id with user scope", async () => {
		const { updateCollection } = await import("./collections");
		const saved = { ...collectionFixture(), name: "Receipts" };
		const fetchMock = vi.fn().mockResolvedValue(jsonResponse(saved));

		await expect(
			updateCollection(
				fetchMock,
				"collection/1",
				{ name: "Receipts", schema_ids: [] },
				{ userId: "user-1" }
			)
		).resolves.toEqual(saved);
		expect(fetchMock).toHaveBeenCalledWith(
			"http://collection-api.test/api/collection/collection%2F1?user_id=user-1",
			expect.objectContaining({ method: "PUT" })
		);
		expectJsonContentType(fetchMock.mock.calls[0]?.[1]);
		expectInternalToken(fetchMock.mock.calls[0]?.[1]);
	});

	it("deletes collections and accepts 204 responses", async () => {
		const { deleteCollection } = await import("./collections");
		const fetchMock = vi.fn().mockResolvedValue(new Response(null, { status: 204 }));

		await expect(deleteCollection(fetchMock, "collection-1", { userId: "user-1" })).resolves.toEqual({
			deleted_id: "collection-1"
		});
		expect(fetchMock).toHaveBeenCalledWith(
			"http://collection-api.test/api/collection/collection-1?user_id=user-1",
			expect.objectContaining({ method: "DELETE" })
		);
		expectInternalToken(fetchMock.mock.calls[0]?.[1]);
	});

	it("rejects collection requests before calling Go when the internal API token is missing", async () => {
		const { listCollectionsPage } = await import("./collections");
		delete privateEnv.SYNCRA_INTERNAL_API_TOKEN;
		const fetchMock = vi.fn();

		await expect(listCollectionsPage(fetchMock)).rejects.toMatchObject({
			status: 500,
			message: "Collection service is not configured"
		});
		expect(fetchMock).not.toHaveBeenCalled();
	});

	it("throws typed errors for successful invalid list payloads", async () => {
		const { listCollectionsPage } = await import("./collections");
		const fetchMock = vi
			.fn()
			.mockResolvedValue(jsonResponse({ collections: [{ id: "collection-1" }], next_cursor: null }));

		await expect(listCollectionsPage(fetchMock, { userId: "user-1" })).rejects.toMatchObject({
			status: 502,
			message: "Invalid collection response"
		});
	});

	it("throws typed errors for empty non-204 delete responses", async () => {
		const { deleteCollection } = await import("./collections");
		const fetchMock = vi.fn().mockResolvedValue(new Response(null, { status: 200 }));

		await expect(deleteCollection(fetchMock, "collection-1")).rejects.toMatchObject({
			status: 502,
			message: "Invalid collection delete response"
		});
	});

	it("throws typed collection API errors from backend error responses", async () => {
		const { createCollection, isCollectionApiError } = await import("./collections");
		const fetchMock = vi
			.fn()
			.mockResolvedValue(jsonResponse({ error: "name is required" }, { status: 400 }));

		await expect(
			createCollection(fetchMock, { name: "", schema_ids: [] }, { userId: "user-1" })
		).rejects.toMatchObject({ status: 400, message: "name is required" });
		await createCollection(fetchMock, { name: "", schema_ids: [] }, { userId: "user-1" }).catch(
			(error) => {
				expect(isCollectionApiError(error)).toBe(true);
			}
		);
	});
});
