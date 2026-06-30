import { beforeEach, describe, expect, it, vi } from "vitest";

import { DELETE, GET, PUT } from "./+server";
import type { RequestEvent } from "./$types";

const { deleteAdminJSONRecipeMock, getAdminJSONRecipeMock, updateAdminJSONRecipeMock } = vi.hoisted(() => ({
	deleteAdminJSONRecipeMock: vi.fn(),
	getAdminJSONRecipeMock: vi.fn(),
	updateAdminJSONRecipeMock: vi.fn()
}));

vi.mock("$lib/server/admin", () => ({
	deleteAdminJSONRecipe: deleteAdminJSONRecipeMock,
	getAdminJSONRecipe: getAdminJSONRecipeMock,
	updateAdminJSONRecipe: updateAdminJSONRecipeMock,
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
	body?: unknown,
	user: unknown = { id: "admin-1", role: "admin" },
	method = body === undefined ? "GET" : "PUT"
) {
	const url = "http://localhost/api/admin/json-recipes/recipe-1";
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
		params: { id: "recipe-1" },
		fetch: vi.fn(),
		locals: { user, adminUser: user }
	} as unknown as RequestEvent;
}

async function responseJson(response: Response) {
	return (await response.json()) as unknown;
}

describe("admin JSON recipe detail API endpoint", () => {
	beforeEach(() => {
		deleteAdminJSONRecipeMock.mockReset();
		getAdminJSONRecipeMock.mockReset();
		updateAdminJSONRecipeMock.mockReset();
	});

	it("loads recipe details for admins", async () => {
		const result = recipe();
		getAdminJSONRecipeMock.mockResolvedValue(result);
		const event = createEvent();

		const response = await GET(event);

		expect(getAdminJSONRecipeMock).toHaveBeenCalledWith(event.fetch, "auth.session_token=session-1", "recipe-1");
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(result);
	});

	it("updates recipes with only editable fields", async () => {
		const result = { ...recipe(), title: "Receipt" };
		updateAdminJSONRecipeMock.mockResolvedValue(result);
		const event = createEvent({
			title: "Receipt",
			description: "Receipt fields",
			json: { type: "object" },
			category_id: null
		});

		const response = await PUT(event);

		expect(updateAdminJSONRecipeMock).toHaveBeenCalledWith(event.fetch, "auth.session_token=session-1", "recipe-1", {
			title: "Receipt",
			description: "Receipt fields",
			json: { type: "object" },
			category_id: null
		});
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(result);
	});

	it("rejects counter updates before forwarding", async () => {
		const response = await PUT(
			createEvent({
				title: "Receipt",
				description: "Receipt fields",
				json: { type: "object" },
				counter: 99
			})
		);

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "invalid JSON recipe payload" });
		expect(updateAdminJSONRecipeMock).not.toHaveBeenCalled();
	});

	it("deletes recipes for admins", async () => {
		deleteAdminJSONRecipeMock.mockResolvedValue(undefined);
		const event = createEvent(undefined, { id: "admin-1", role: "admin" }, "DELETE");

		const response = await DELETE(event);

		expect(deleteAdminJSONRecipeMock).toHaveBeenCalledWith(event.fetch, "auth.session_token=session-1", "recipe-1");
		expect(response.status).toBe(204);
		expect(await response.text()).toBe("");
	});

	it("returns 403 for non-admin users", async () => {
		const response = await GET(createEvent(undefined, { id: "user-1", role: "user" }));

		expect(response.status).toBe(403);
		expect(await responseJson(response)).toEqual({ error: "admin access required" });
		expect(getAdminJSONRecipeMock).not.toHaveBeenCalled();
	});
});
