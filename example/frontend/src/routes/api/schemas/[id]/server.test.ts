import { beforeEach, describe, expect, it, vi } from "vitest";

import { SchemaApiError } from "$lib/server/schemas";
import { DELETE, GET, PUT } from "./+server";
import type { RequestEvent } from "./$types";

const { deleteSchemaMock, getSchemaMock, updateSchemaMock, SchemaApiErrorMock } = vi.hoisted(
	() => {
		class MockSchemaApiError extends Error {
			status: number;

			constructor(status: number, message: string) {
				super(message);
				this.name = "SchemaApiError";
				this.status = status;
			}
		}

		return {
			deleteSchemaMock: vi.fn(),
			getSchemaMock: vi.fn(),
			updateSchemaMock: vi.fn(),
			SchemaApiErrorMock: MockSchemaApiError
		};
	}
);

vi.mock("$lib/server/schemas", () => ({
	deleteSchema: deleteSchemaMock,
	getSchema: getSchemaMock,
	updateSchema: updateSchemaMock,
	SchemaApiError: SchemaApiErrorMock,
	isSchemaApiError: (error: unknown) => error instanceof SchemaApiErrorMock
}));

function createEvent(
	method: "DELETE" | "GET" | "PUT",
	body: unknown = undefined,
	user: unknown = { id: "user-1" }
) {
	return {
		request: new Request("http://localhost/api/schemas/schema-1", {
			method,
			body: body === undefined ? undefined : JSON.stringify(body)
		}),
		fetch: vi.fn(),
		locals: {
			user
		},
		params: {
			id: "schema-1"
		}
	} as unknown as RequestEvent;
}

async function responseJson(response: Response) {
	return (await response.json()) as unknown;
}

describe("schema item API endpoint", () => {
	beforeEach(() => {
		deleteSchemaMock.mockReset();
		getSchemaMock.mockReset();
		updateSchemaMock.mockReset();
	});

	it("returns 401 for unauthenticated GET requests", async () => {
		const response = await GET(createEvent("GET", undefined, null));

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(getSchemaMock).not.toHaveBeenCalled();
	});

	it("calls the schema service with current user id for GET requests", async () => {
		const schema = {
			id: "schema-1",
			created_at: "2026-05-26T00:00:00Z",
			updated_at: "2026-05-26T00:00:00Z",
			user_id: "user-1",
			name: "Customer",
			description: "Customer profile payload",
			strict: true,
			schema: { type: "object" }
		};
		getSchemaMock.mockResolvedValue(schema);
		const event = createEvent("GET");

		const response = await GET(event);

		expect(getSchemaMock).toHaveBeenCalledWith(event.fetch, "schema-1", { userId: "user-1" });
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(schema);
	});

	it("returns 400 when the PUT schema name is empty", async () => {
		const response = await PUT(
			createEvent("PUT", {
				name: "   ",
				description: "Customer profile payload",
				strict: true,
				schema: { type: "object" }
			})
		);

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "name is required" });
		expect(updateSchemaMock).not.toHaveBeenCalled();
	});

	it("calls the schema service with a normalized payload for PUT requests", async () => {
		const schema = {
			id: "schema-1",
			created_at: "2026-05-26T00:00:00Z",
			updated_at: "2026-05-26T00:00:00Z",
			user_id: "user-1",
			name: "Customer",
			description: "",
			strict: true,
			schema: { type: "object" }
		};
		updateSchemaMock.mockResolvedValue(schema);
		const event = createEvent("PUT", {
			name: "  Customer  ",
			schema: { type: "object" }
		});

		const response = await PUT(event);

		expect(updateSchemaMock).toHaveBeenCalledWith(
			event.fetch,
			"schema-1",
			{
				name: "Customer",
				description: "",
				strict: true,
				schema: { type: "object" }
			},
			{ userId: "user-1" }
		);
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(schema);
	});

	it("calls the schema service with current user id for DELETE requests", async () => {
		const deleted = { deleted_ids: ["schema-1"], deleted_count: 1 };
		deleteSchemaMock.mockResolvedValue(deleted);
		const event = createEvent("DELETE");

		const response = await DELETE(event);

		expect(deleteSchemaMock).toHaveBeenCalledWith(event.fetch, "schema-1", { userId: "user-1" });
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(deleted);
	});

	it("preserves schema service client errors", async () => {
		getSchemaMock.mockRejectedValue(new SchemaApiError(404, "schema not found"));

		const response = await GET(createEvent("GET"));

		expect(response.status).toBe(404);
		expect(await responseJson(response)).toEqual({ error: "schema not found" });
	});

	it("normalizes schema service server errors", async () => {
		updateSchemaMock.mockRejectedValue(new SchemaApiError(503, "Schema service unavailable"));

		const response = await PUT(
			createEvent("PUT", { name: "Customer", schema: { type: "object" } })
		);

		expect(response.status).toBe(502);
		expect(await responseJson(response)).toEqual({ error: "A server error occurred. Please try again." });
	});

	it("rethrows unknown schema service errors", async () => {
		const error = new Error("unexpected failure");
		deleteSchemaMock.mockRejectedValue(error);

		await expect(DELETE(createEvent("DELETE"))).rejects.toBe(error);
	});
});
