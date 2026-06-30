import { beforeEach, describe, expect, it, vi } from "vitest";

import { GET, POST } from "./+server";
import type { RequestEvent } from "./$types";

const { createAdminJSONRecipeCategoryMock, listAdminJSONRecipeCategoriesMock } = vi.hoisted(() => ({
	createAdminJSONRecipeCategoryMock: vi.fn(),
	listAdminJSONRecipeCategoriesMock: vi.fn()
}));

vi.mock("$lib/server/admin", () => ({
	createAdminJSONRecipeCategory: createAdminJSONRecipeCategoryMock,
	listAdminJSONRecipeCategories: listAdminJSONRecipeCategoriesMock,
	isAdminApiError: () => false
}));

function category() {
	return {
		id: "category-1",
		title: { en: "Invoices", ro: "Facturi" },
		created_at: "2026-06-19T00:00:00Z",
		updated_at: "2026-06-19T00:00:00Z"
	};
}

function createEvent(
	url = "http://localhost/api/admin/json-recipe-categories",
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

describe("admin JSON recipe categories API endpoint", () => {
	beforeEach(() => {
		createAdminJSONRecipeCategoryMock.mockReset();
		listAdminJSONRecipeCategoriesMock.mockReset();
	});

	it("requires admin auth", async () => {
		const response = await GET(createEvent("http://localhost/api/admin/json-recipe-categories", undefined, null));

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(listAdminJSONRecipeCategoriesMock).not.toHaveBeenCalled();
	});

	it("lists categories with cookies", async () => {
		const result = { categories: [category()] };
		listAdminJSONRecipeCategoriesMock.mockResolvedValue(result);
		const event = createEvent();

		const response = await GET(event);

		expect(listAdminJSONRecipeCategoriesMock).toHaveBeenCalledWith(event.fetch, "auth.session_token=session-1");
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(result);
	});

	it("rejects unknown create payload keys before forwarding", async () => {
		const response = await POST(
			createEvent("http://localhost/api/admin/json-recipe-categories", {
				title: { en: "Invoices", ro: "Facturi" },
				slug: "invoices"
			})
		);

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "invalid JSON recipe category payload" });
		expect(createAdminJSONRecipeCategoryMock).not.toHaveBeenCalled();
	});

	it("creates categories for admins", async () => {
		const result = category();
		createAdminJSONRecipeCategoryMock.mockResolvedValue(result);
		const event = createEvent("http://localhost/api/admin/json-recipe-categories", {
			title: { en: "Invoices", ro: "Facturi" }
		});

		const response = await POST(event);

		expect(createAdminJSONRecipeCategoryMock).toHaveBeenCalledWith(event.fetch, "auth.session_token=session-1", {
			title: { en: "Invoices", ro: "Facturi" }
		});
		expect(response.status).toBe(201);
		expect(await responseJson(response)).toEqual(result);
	});
});
