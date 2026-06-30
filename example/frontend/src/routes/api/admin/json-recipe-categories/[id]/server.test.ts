import { beforeEach, describe, expect, it, vi } from "vitest";

import { DELETE, GET, PUT } from "./+server";
import type { RequestEvent } from "./$types";

const {
	deleteAdminJSONRecipeCategoryMock,
	getAdminJSONRecipeCategoryMock,
	updateAdminJSONRecipeCategoryMock
} = vi.hoisted(() => ({
	deleteAdminJSONRecipeCategoryMock: vi.fn(),
	getAdminJSONRecipeCategoryMock: vi.fn(),
	updateAdminJSONRecipeCategoryMock: vi.fn()
}));

vi.mock("$lib/server/admin", () => ({
	deleteAdminJSONRecipeCategory: deleteAdminJSONRecipeCategoryMock,
	getAdminJSONRecipeCategory: getAdminJSONRecipeCategoryMock,
	updateAdminJSONRecipeCategory: updateAdminJSONRecipeCategoryMock,
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
	body?: unknown,
	user: unknown = { id: "admin-1", role: "admin" },
	method = body === undefined ? "GET" : "PUT"
) {
	const url = "http://localhost/api/admin/json-recipe-categories/category-1";
	const init: RequestInit =
		body === undefined
			? { method, headers: { cookie: "auth.session_token=session-1" } }
			: {
					method,
					headers: {
						"content-type": "application/json",
						cookie: "auth.session_token=session-1"
					},
					body: JSON.stringify(body)
				};
	return {
		url: new URL(url),
		request: new Request(url, init),
		params: { id: "category-1" },
		fetch: vi.fn(),
		locals: { user, adminUser: user }
	} as unknown as RequestEvent;
}

async function responseJson(response: Response) {
	return (await response.json()) as unknown;
}

describe("admin JSON recipe category detail API endpoint", () => {
	beforeEach(() => {
		deleteAdminJSONRecipeCategoryMock.mockReset();
		getAdminJSONRecipeCategoryMock.mockReset();
		updateAdminJSONRecipeCategoryMock.mockReset();
	});

	it("loads category details for admins", async () => {
		const result = category();
		getAdminJSONRecipeCategoryMock.mockResolvedValue(result);
		const event = createEvent();

		const response = await GET(event);

		expect(getAdminJSONRecipeCategoryMock).toHaveBeenCalledWith(
			event.fetch,
			"auth.session_token=session-1",
			"category-1"
		);
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(result);
	});

	it("updates categories with localized titles", async () => {
		const result = { ...category(), title: { en: "Receipts", ro: "Bonuri" } };
		updateAdminJSONRecipeCategoryMock.mockResolvedValue(result);
		const event = createEvent({
			title: { en: "Receipts", ro: "Bonuri" }
		});

		const response = await PUT(event);

		expect(updateAdminJSONRecipeCategoryMock).toHaveBeenCalledWith(
			event.fetch,
			"auth.session_token=session-1",
			"category-1",
			{ title: { en: "Receipts", ro: "Bonuri" } }
		);
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(result);
	});

	it("rejects malformed title payloads before forwarding", async () => {
		const response = await PUT(
			createEvent({
				title: { en: "Receipts" }
			})
		);

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "invalid JSON recipe category payload" });
		expect(updateAdminJSONRecipeCategoryMock).not.toHaveBeenCalled();
	});

	it("deletes categories for admins", async () => {
		deleteAdminJSONRecipeCategoryMock.mockResolvedValue(undefined);
		const event = createEvent(undefined, { id: "admin-1", role: "admin" }, "DELETE");

		const response = await DELETE(event);

		expect(deleteAdminJSONRecipeCategoryMock).toHaveBeenCalledWith(
			event.fetch,
			"auth.session_token=session-1",
			"category-1"
		);
		expect(response.status).toBe(204);
		expect(await response.text()).toBe("");
	});
});
