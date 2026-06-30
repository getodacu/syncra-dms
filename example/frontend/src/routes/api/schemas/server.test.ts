import { beforeEach, describe, expect, it, vi } from "vitest";

import { SchemaApiError } from "$lib/server/schemas";
import { DELETE, GET, POST } from "./+server";
import type { RequestEvent } from "./$types";

const MAX_SCHEMA_REQUEST_BYTES = 1 << 20;
const MAX_SCHEMA_JSON_BYTES = 1 << 20;

const {
	createSchemaMock,
	deleteSchemasMock,
	listSchemasMock,
	listSchemasPageMock,
	SchemaApiErrorMock
} = vi.hoisted(() => {
	class MockSchemaApiError extends Error {
		status: number;

		constructor(status: number, message: string) {
			super(message);
			this.name = "SchemaApiError";
			this.status = status;
		}
	}

	return {
		createSchemaMock: vi.fn(),
		deleteSchemasMock: vi.fn(),
		listSchemasMock: vi.fn(),
		listSchemasPageMock: vi.fn(),
		SchemaApiErrorMock: MockSchemaApiError
	};
});

vi.mock("$lib/server/schemas", () => ({
	createSchema: createSchemaMock,
	deleteSchemas: deleteSchemasMock,
	listSchemas: listSchemasMock,
	listSchemasPage: listSchemasPageMock,
	SchemaApiError: SchemaApiErrorMock,
	isSchemaApiError: (error: unknown) => error instanceof SchemaApiErrorMock
}));

function createEvent(body: unknown, user: unknown = { id: "user-1" }) {
	return createRequestEvent(
		new Request("http://localhost/api/schemas", {
			method: "POST",
			body: body === undefined ? undefined : JSON.stringify(body)
		}),
		user
	);
}

function createDeleteEvent(body: unknown, user: unknown = { id: "user-1" }) {
	return createRequestEvent(
		new Request("http://localhost/api/schemas", {
			method: "DELETE",
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

function streamRequest(bodySize: number) {
	return new Request("http://localhost/api/schemas", {
		method: "POST",
		body: new ReadableStream<Uint8Array>({
			start(controller) {
				controller.enqueue(new Uint8Array(bodySize).fill(123));
				controller.close();
			}
		}),
		duplex: "half"
	} as RequestInit & { duplex: "half" });
}

async function responseJson(response: Response) {
	return (await response.json()) as unknown;
}

function buildLargeSchemaBody() {
	const properties: string[] = [];
	const prefix = `{"name":"Large schema","schema":`;
	let schemaByteLength = 2;

	for (let index = 0; index < 100_000; index += 1) {
		const property = `"k${index.toString(36)}":1e100`;
		properties.push(property);
		schemaByteLength += property.length + (properties.length === 1 ? 0 : 1);

		if (
			prefix.length + schemaByteLength + 1 <= MAX_SCHEMA_REQUEST_BYTES &&
			schemaByteLength + properties.length > MAX_SCHEMA_JSON_BYTES
		) {
			return `${prefix}{${properties.join(",")}}}`;
		}
	}

	throw new Error("failed to construct oversized schema fixture");
}

describe("schema API endpoint", () => {
	beforeEach(() => {
		createSchemaMock.mockReset();
		deleteSchemasMock.mockReset();
		listSchemasMock.mockReset();
		listSchemasPageMock.mockReset();
	});

	it("returns 401 for unauthenticated GET requests", async () => {
		const response = await GET(createGetEvent("http://localhost/api/schemas", null));

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(listSchemasMock).not.toHaveBeenCalled();
	});

	it("defaults missing GET scope to system schemas", async () => {
		const schemas = [
			{
				id: "schema-1",
				created_at: "2026-05-26T00:00:00Z",
				updated_at: "2026-05-26T00:00:00Z",
				name: "Invoice",
				description: "Invoice extraction",
				strict: true,
				schema: { type: "object" },
				user_id: null
			}
		];
		listSchemasMock.mockResolvedValue(schemas);
		const event = createGetEvent("http://localhost/api/schemas");

		const response = await GET(event);

		expect(listSchemasMock).toHaveBeenCalledWith(event.fetch);
		expect(listSchemasPageMock).not.toHaveBeenCalled();
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(schemas);
	});

	it("calls the schema service without user id for system scope", async () => {
		const schemas = [
			{
				id: "schema-1",
				created_at: "2026-05-26T00:00:00Z",
				updated_at: "2026-05-26T00:00:00Z",
				name: "Receipt",
				description: "Receipt extraction",
				strict: false,
				schema: { type: "object" },
				user_id: null
			}
		];
		listSchemasMock.mockResolvedValue(schemas);
		const event = createGetEvent("http://localhost/api/schemas?scope=system");

		const response = await GET(event);

		expect(listSchemasMock).toHaveBeenCalledWith(event.fetch);
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(schemas);
	});

	it("returns a paginated schema response when pagination parameters are present", async () => {
		const page = {
			schemas: [
				{
					id: "schema-2",
					created_at: "2026-05-26T00:00:00Z",
					updated_at: "2026-05-26T00:00:00Z",
					name: "Customer",
					description: "Customer extraction",
					strict: true,
					schema: { type: "object" },
					user_id: "user-1"
				}
			],
			next_cursor: "next-page"
		};
		listSchemasPageMock.mockResolvedValue(page);
		const event = createGetEvent(
			"http://localhost/api/schemas?scope=mine&size=25&cursor=abc&sort=desc"
		);

		const response = await GET(event);

		expect(listSchemasPageMock).toHaveBeenCalledWith(event.fetch, {
			userId: "user-1",
			cursor: "abc",
			size: "25",
			sort: "desc"
		});
		expect(listSchemasMock).not.toHaveBeenCalled();
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(page);
	});

	it("calls the schema service with user id for mine scope", async () => {
		const schemas = [
			{
				id: "schema-2",
				created_at: "2026-05-26T00:00:00Z",
				updated_at: "2026-05-26T00:00:00Z",
				name: "Customer",
				description: "Customer extraction",
				strict: true,
				schema: { type: "object" },
				user_id: "user-1"
			}
		];
		listSchemasMock.mockResolvedValue(schemas);
		const event = createGetEvent("http://localhost/api/schemas?scope=mine");

		const response = await GET(event);

		expect(listSchemasMock).toHaveBeenCalledWith(event.fetch, { userId: "user-1" });
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(schemas);
	});

	it("returns 400 for invalid GET scope", async () => {
		const response = await GET(createGetEvent("http://localhost/api/schemas?scope=team"));

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "invalid schema scope" });
		expect(listSchemasMock).not.toHaveBeenCalled();
	});

	it("preserves schema service client errors for GET requests", async () => {
		listSchemasMock.mockRejectedValue(new SchemaApiError(409, "schema conflict"));

		const response = await GET(createGetEvent("http://localhost/api/schemas?scope=mine"));

		expect(response.status).toBe(409);
		expect(await responseJson(response)).toEqual({ error: "schema conflict" });
	});

	it("normalizes schema service server errors for GET requests", async () => {
		listSchemasMock.mockRejectedValue(new SchemaApiError(503, "Schema service unavailable"));

		const response = await GET(createGetEvent("http://localhost/api/schemas?scope=system"));

		expect(response.status).toBe(502);
		expect(await responseJson(response)).toEqual({ error: "A server error occurred. Please try again." });
	});

	it("rethrows unknown schema service errors for GET requests", async () => {
		const error = new Error("unexpected failure");
		listSchemasMock.mockRejectedValue(error);

		await expect(GET(createGetEvent("http://localhost/api/schemas?scope=mine"))).rejects.toBe(
			error
		);
	});

	it("returns 401 for unauthenticated requests", async () => {
		const response = await POST(createEvent({ name: "Customer" }, null));

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(createSchemaMock).not.toHaveBeenCalled();
	});

	it("returns 400 when the schema name is empty", async () => {
		const response = await POST(
			createEvent({
				name: "   ",
				description: "Customer profile payload",
				strict: true,
				schema: { type: "object" }
			})
		);

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "name is required" });
		expect(createSchemaMock).not.toHaveBeenCalled();
	});

	it("returns 400 when content-length exceeds the request limit", async () => {
		const response = await POST(
			createRequestEvent(
				new Request("http://localhost/api/schemas", {
					method: "POST",
					headers: { "content-length": String(MAX_SCHEMA_REQUEST_BYTES + 1) },
					body: "{}"
				})
			)
		);

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "request body too large" });
		expect(createSchemaMock).not.toHaveBeenCalled();
	});

	it("returns 400 when a streamed body exceeds the request limit", async () => {
		const response = await POST(createRequestEvent(streamRequest(MAX_SCHEMA_REQUEST_BYTES + 1)));

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "request body too large" });
		expect(createSchemaMock).not.toHaveBeenCalled();
	});

	it("returns 400 when the schema JSON exceeds the schema limit", async () => {
		const response = await POST(
			createRequestEvent(
				new Request("http://localhost/api/schemas", {
					method: "POST",
					body: buildLargeSchemaBody()
				})
			)
		);

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "schema is too large" });
		expect(createSchemaMock).not.toHaveBeenCalled();
	});

	it("calls the schema service with event fetch, current user id, and a normalized payload", async () => {
		const savedSchema = {
			id: "schema-1",
			created_at: "2026-05-26T00:00:00Z",
			updated_at: "2026-05-26T00:00:00Z",
			user_id: "user-1",
			name: "Customer",
			description: "Customer profile payload",
			strict: false,
			schema: { type: "object" }
		};
		createSchemaMock.mockResolvedValue(savedSchema);
		const event = createEvent({
			name: "  Customer  ",
			description: "Customer profile payload",
			strict: false,
			schema: { type: "object" }
		});

		const response = await POST(event);

		expect(createSchemaMock).toHaveBeenCalledWith(
			event.fetch,
			{
				name: "Customer",
				description: "Customer profile payload",
				strict: false,
				schema: { type: "object" }
			},
			{ userId: "user-1" }
		);
		expect(response.status).toBe(201);
		expect(await responseJson(response)).toEqual(savedSchema);
	});

	it("counts schema names by Unicode code point", async () => {
		createSchemaMock.mockResolvedValue({
			id: "schema-1",
			created_at: "2026-05-26T00:00:00Z",
			updated_at: "2026-05-26T00:00:00Z",
			name: "schema",
			description: "",
			strict: true,
			schema: { type: "object" }
		});
		const name = "😀".repeat(160);
		const event = createEvent({ name, schema: { type: "object" } });

		const response = await POST(event);

		expect(response.status).toBe(201);
		expect(createSchemaMock).toHaveBeenCalledWith(
			event.fetch,
			expect.objectContaining({ name }),
			{ userId: "user-1" }
		);
	});

	it("returns 400 when the schema service rejects with a client error", async () => {
		createSchemaMock.mockRejectedValue(new SchemaApiError(400, "invalid schema payload"));

		const response = await POST(
			createEvent({
				name: "Customer",
				schema: { type: "object" }
			})
		);

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "invalid schema payload" });
	});

	it("returns 502 when the schema service rejects with an upstream server error", async () => {
		createSchemaMock.mockRejectedValue(new SchemaApiError(503, "Schema service unavailable"));

		const response = await POST(
			createEvent({
				name: "Customer",
				schema: { type: "object" }
			})
		);

		expect(response.status).toBe(502);
		expect(await responseJson(response)).toEqual({ error: "A server error occurred. Please try again." });
	});

	it("rethrows unknown schema service errors", async () => {
		const error = new Error("unexpected failure");
		createSchemaMock.mockRejectedValue(error);

		await expect(
			POST(
				createEvent({
					name: "Customer",
					schema: { type: "object" }
				})
			)
		).rejects.toBe(error);
	});

	it("returns 401 for unauthenticated DELETE requests", async () => {
		const response = await DELETE(createDeleteEvent({ ids: ["schema-1"] }, null));

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(deleteSchemasMock).not.toHaveBeenCalled();
	});

	it("returns 400 when DELETE ids are missing or invalid", async () => {
		const response = await DELETE(createDeleteEvent({ ids: ["schema-1", 42] }));

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "ids must be an array of strings" });
		expect(deleteSchemasMock).not.toHaveBeenCalled();
	});

	it("calls the schema service with current user id for DELETE requests", async () => {
		const deleted = { deleted_ids: ["schema-1", "schema-2"], deleted_count: 2 };
		deleteSchemasMock.mockResolvedValue(deleted);
		const event = createDeleteEvent({ ids: ["schema-1", "schema-2"] });

		const response = await DELETE(event);

		expect(deleteSchemasMock).toHaveBeenCalledWith(
			event.fetch,
			["schema-1", "schema-2"],
			{ userId: "user-1" }
		);
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(deleted);
	});

	it("preserves schema service client errors for DELETE requests", async () => {
		deleteSchemasMock.mockRejectedValue(new SchemaApiError(404, "schema not found"));

		const response = await DELETE(createDeleteEvent({ ids: ["schema-1"] }));

		expect(response.status).toBe(404);
		expect(await responseJson(response)).toEqual({ error: "schema not found" });
	});

	it("normalizes schema service server errors for DELETE requests", async () => {
		deleteSchemasMock.mockRejectedValue(new SchemaApiError(503, "Schema service unavailable"));

		const response = await DELETE(createDeleteEvent({ ids: ["schema-1"] }));

		expect(response.status).toBe(502);
		expect(await responseJson(response)).toEqual({ error: "A server error occurred. Please try again." });
	});
});
