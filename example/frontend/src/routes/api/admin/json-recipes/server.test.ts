import { beforeEach, describe, expect, it, vi } from "vitest";

import { GET, POST } from "./+server";
import type { RequestEvent } from "./$types";

const { createAdminJSONRecipeMock, listAdminJSONRecipesMock } = vi.hoisted(() => ({
	createAdminJSONRecipeMock: vi.fn(),
	listAdminJSONRecipesMock: vi.fn()
}));

vi.mock("$lib/server/admin", () => ({
	createAdminJSONRecipe: createAdminJSONRecipeMock,
	listAdminJSONRecipes: listAdminJSONRecipesMock,
	isAdminApiError: () => false
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

function createEvent(
	url = "http://localhost/api/admin/json-recipes?cursor=cursor-0&size=50&sort=asc",
	body?: unknown,
	user: unknown = { id: "admin-1", role: "admin" }
) {
	const init: RequestInit =
		body === undefined
			? { headers: { cookie: "auth.session_token=session-1" } }
			: {
					method: "POST",
					headers: {
						"content-type": "application/json",
						cookie: "auth.session_token=session-1"
					},
					body: JSON.stringify(body)
				};
	return {
		url: new URL(url),
		request: new Request(url, init),
		fetch: vi.fn(),
		locals: { user, adminUser: user }
	} as unknown as RequestEvent;
}

async function responseJson(response: Response) {
	return (await response.json()) as unknown;
}

describe("admin JSON recipes API endpoint", () => {
	beforeEach(() => {
		createAdminJSONRecipeMock.mockReset();
		listAdminJSONRecipesMock.mockReset();
	});

	it("requires admin auth", async () => {
		const response = await GET(createEvent("http://localhost/api/admin/json-recipes", undefined, null));

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(listAdminJSONRecipesMock).not.toHaveBeenCalled();
	});

	it("forwards list pagination with cookies", async () => {
		const result = { recipes: [recipe()], next_cursor: "cursor-1" };
		listAdminJSONRecipesMock.mockResolvedValue(result);
		const event = createEvent();

		const response = await GET(event);

		expect(listAdminJSONRecipesMock).toHaveBeenCalledWith(event.fetch, "auth.session_token=session-1", {
			cursor: "cursor-0",
			size: "50",
			sort: "asc"
		});
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(result);
	});

	it("rejects unknown create payload keys before forwarding", async () => {
		const response = await POST(
			createEvent("http://localhost/api/admin/json-recipes", {
				title: "Invoice",
				description: "Invoice fields",
				json: { type: "object" },
				counter: 10
			})
		);

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "invalid JSON recipe payload" });
		expect(createAdminJSONRecipeMock).not.toHaveBeenCalled();
	});

	it("creates recipes for admins", async () => {
		const result = recipe();
		createAdminJSONRecipeMock.mockResolvedValue(result);
		const event = createEvent("http://localhost/api/admin/json-recipes", {
			title: "Invoice",
			description: "Invoice fields",
			json: { type: "object" },
			category_id: "category-1"
		});

		const response = await POST(event);

		expect(createAdminJSONRecipeMock).toHaveBeenCalledWith(event.fetch, "auth.session_token=session-1", {
			title: "Invoice",
			description: "Invoice fields",
			json: { type: "object" },
			category_id: "category-1"
		});
		expect(response.status).toBe(201);
		expect(await responseJson(response)).toEqual(result);
	});
});
