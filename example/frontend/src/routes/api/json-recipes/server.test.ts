import { beforeEach, describe, expect, it, vi } from "vitest";

import { GET } from "./+server";
import type { RequestEvent } from "./$types";

const { listJsonRecipesMock, schemaApiErrorCtor } = vi.hoisted(() => {
	class MockSchemaApiError extends Error {
		status: number;

		constructor(status: number, message: string) {
			super(message);
			this.name = "SchemaApiError";
			this.status = status;
		}
	}

	return {
		listJsonRecipesMock: vi.fn(),
		schemaApiErrorCtor: MockSchemaApiError
	};
});

vi.mock("$lib/server/schemas", () => ({
	SchemaApiError: schemaApiErrorCtor,
	isSchemaApiError: (error: unknown) => error instanceof schemaApiErrorCtor,
	listJsonRecipes: listJsonRecipesMock
}));

function recipe() {
	const category = {
		id: "category-1",
		title: { en: "Invoices", ro: "Facturi" },
		created_at: "2026-06-19T00:00:00Z",
		updated_at: "2026-06-19T00:00:00Z"
	};
	return {
		id: "recipe-1",
		title: "Invoice",
		description: "Invoice fields",
		json: { type: "object" },
		counter: 0,
		category_id: category.id,
		category,
		created_at: "2026-06-20T00:00:00Z",
		updated_at: "2026-06-20T00:00:00Z"
	};
}

function createEvent(url = "http://localhost/api/json-recipes?cursor=cursor-0&size=50&sort=asc") {
	return {
		url: new URL(url),
		request: new Request(url),
		fetch: vi.fn()
	} as unknown as RequestEvent;
}

async function responseJson(response: Response) {
	return (await response.json()) as unknown;
}

describe("public JSON recipes API endpoint", () => {
	beforeEach(() => {
		listJsonRecipesMock.mockReset();
	});

	it("forwards list pagination query params", async () => {
		const result = { recipes: [recipe()], next_cursor: "cursor-1" };
		listJsonRecipesMock.mockResolvedValue(result);
		const event = createEvent();

		const response = await GET(event);

		expect(listJsonRecipesMock).toHaveBeenCalledWith(event.fetch, {
			cursor: "cursor-0",
			size: "50",
			sort: "asc"
		});
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(result);
	});

	it("sanitizes upstream 5xx errors", async () => {
		listJsonRecipesMock.mockRejectedValue(new schemaApiErrorCtor(503, "database password leaked"));

		const response = await GET(createEvent("http://localhost/api/json-recipes"));

		expect(response.status).toBe(502);
		expect(await responseJson(response)).toEqual({ error: "Failed to load OCR recipes" });
	});

	it("passes through safe client validation errors", async () => {
		listJsonRecipesMock.mockRejectedValue(new schemaApiErrorCtor(400, "invalid cursor"));

		const response = await GET(createEvent("http://localhost/api/json-recipes?cursor=bad"));

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "invalid cursor" });
	});
});
