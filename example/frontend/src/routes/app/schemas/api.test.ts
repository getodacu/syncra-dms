import { describe, expect, it, vi } from "vitest";

import {
	SCHEMA_QUERY_RETRY_LIMIT,
	SchemaClientError,
	cloneSchema,
	deleteSchema,
	deleteSchemas,
	fetchSchemas,
	getSchema,
	isSchemaClientError,
	isSchemaNotFoundError,
	shouldRetrySchemaQuery,
} from "./api";

function jsonResponse(body: unknown, init?: ResponseInit) {
	return new Response(JSON.stringify(body), {
		headers: { "content-type": "application/json" },
		...init,
	});
}

function schemaResponse() {
	return {
		id: "schema-1",
		created_at: "2026-05-27T00:00:00Z",
		updated_at: "2026-05-27T00:01:00Z",
		user_id: "user-1",
		name: "Invoice",
		description: "Invoice extraction schema",
		strict: true,
		schema: { type: "object" },
	};
}

describe("schemas api client", () => {
	it("fetches schemas with default pagination through the SvelteKit proxy", async () => {
		const body = {
			schemas: [schemaResponse()],
			next_cursor: "cursor-1",
		};
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body));

		await expect(fetchSchemas(fetchFn, {})).resolves.toEqual(body);
		expect(fetchFn).toHaveBeenCalledWith("/api/schemas?scope=mine&size=20", { method: "GET" });
	});

	it("throws backend JSON messages for list errors", async () => {
		const fetchFn = vi
			.fn()
			.mockResolvedValue(jsonResponse({ error: "invalid cursor" }, { status: 400 }));

		await expect(fetchSchemas(fetchFn, {})).rejects.toThrow("invalid cursor");
	});

	it("preserves backend response status for schema detail errors", async () => {
		const fetchFn = vi.fn(() =>
			Promise.resolve(jsonResponse({ error: "schema not found" }, { status: 404 }))
		);

		await expect(getSchema(fetchFn, "missing-schema")).rejects.toMatchObject({
			name: "SchemaClientError",
			message: "schema not found",
			status: 404,
		});

		let error: unknown;
		try {
			await getSchema(fetchFn, "missing-schema");
		} catch (caught) {
			error = caught;
		}

		expect(error).toBeInstanceOf(SchemaClientError);
		expect(isSchemaClientError(error)).toBe(true);
		expect(isSchemaNotFoundError(error)).toBe(true);
	});

	it("does not retry missing schema detail queries", () => {
		expect(SCHEMA_QUERY_RETRY_LIMIT).toBe(3);
		expect(shouldRetrySchemaQuery(0, new SchemaClientError(404, "schema not found"))).toBe(
			false
		);
		expect(shouldRetrySchemaQuery(0, new SchemaClientError(500, "schema unavailable"))).toBe(
			true
		);
		expect(shouldRetrySchemaQuery(3, new Error("network unavailable"))).toBe(false);
	});

	it("rejects invalid list responses", async () => {
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse({ schemas: [schemaResponse()] }));

		await expect(fetchSchemas(fetchFn, {})).rejects.toThrow("Invalid schema list response");
	});

	it("rejects blank pagination cursors", async () => {
		const fetchFn = vi
			.fn()
			.mockResolvedValue(jsonResponse({ schemas: [schemaResponse()], next_cursor: " " }));

		await expect(fetchSchemas(fetchFn, {})).rejects.toThrow("Invalid schema list response");
	});

	it("clones schemas through the SvelteKit proxy with copied schema data", async () => {
		const source = {
			...schemaResponse(),
			description: "Detailed invoice schema",
			strict: false,
			schema: { type: "object", properties: { total: { type: "number" } } },
		};
		const cloned = { ...source, id: "schema-2", name: "Clone of Invoice" };
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(cloned));

		await expect(cloneSchema(fetchFn, source)).resolves.toEqual(cloned);
		expect(fetchFn).toHaveBeenCalledWith(
			"/api/schemas",
			expect.objectContaining({
				method: "POST",
				body: JSON.stringify({
					name: "Clone of Invoice",
					description: source.description,
					strict: false,
					schema: source.schema,
				}),
			})
		);
		const headers = fetchFn.mock.calls[0][1]?.headers as Headers;
		expect(headers.get("content-type")).toBe("application/json");
	});

	it("truncates cloned schema names to 160 Unicode code points", async () => {
		const source = {
			...schemaResponse(),
			name: "😀".repeat(200),
		};
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse({ ...source, id: "schema-2" }));

		await cloneSchema(fetchFn, source);

		const requestBody = JSON.parse(String(fetchFn.mock.calls[0][1]?.body)) as { name: string };
		expect(Array.from(requestBody.name)).toHaveLength(160);
		expect(requestBody.name).toMatch(/^Clone of /);
	});

	it("preserves backend JSON messages for clone errors", async () => {
		const fetchFn = vi
			.fn()
			.mockResolvedValue(jsonResponse({ error: "schema name already exists" }, { status: 409 }));

		await expect(cloneSchema(fetchFn, schemaResponse())).rejects.toThrow(
			"schema name already exists"
		);
	});

	it("path-encodes schema ids for single deletes", async () => {
		const body = { deleted_ids: ["schema/1"], deleted_count: 1 };
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body));

		await expect(deleteSchema(fetchFn, "schema/1")).resolves.toEqual(body);
		expect(fetchFn).toHaveBeenCalledWith("/api/schemas/schema%2F1", { method: "DELETE" });
	});

	it("bulk deletes schemas through the SvelteKit proxy", async () => {
		const body = { deleted_ids: ["schema-1", "schema-2"], deleted_count: 2 };
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body));

		await expect(deleteSchemas(fetchFn, ["schema-1", "schema-2"])).resolves.toEqual(body);
		expect(fetchFn).toHaveBeenCalledWith("/api/schemas", {
			method: "DELETE",
			headers: { "content-type": "application/json" },
			body: JSON.stringify({ ids: ["schema-1", "schema-2"] }),
		});
	});

	it("rejects invalid delete responses", async () => {
		const fetchFn = vi
			.fn()
			.mockResolvedValue(jsonResponse({ deleted_ids: ["schema-1"], deleted_count: 2 }));

		await expect(deleteSchemas(fetchFn, ["schema-1"])).rejects.toThrow(
			"Invalid schema delete response"
		);
	});
});
