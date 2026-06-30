import { beforeEach, describe, expect, it, vi } from "vitest";

import { POST } from "./+server";
import type { RequestEvent } from "./$types";

const { deployJsonRecipeMock } = vi.hoisted(() => ({
	deployJsonRecipeMock: vi.fn()
}));

vi.mock("$lib/server/schemas", () => ({
	deployJsonRecipe: deployJsonRecipeMock,
	isSchemaApiError: () => false
}));

function deployResponse() {
	return {
		recipe: {
			id: "recipe-1",
			title: "Invoice",
			description: "Invoice fields",
			json: { type: "object" },
			counter: 1,
			created_at: "2026-06-20T00:00:00Z",
			updated_at: "2026-06-20T00:00:00Z"
		},
		schema: {
			id: "schema-1",
			user_id: "user-1",
			name: "Invoice",
			description: "Invoice fields",
			schema: { type: "object" },
			strict: true,
			created_at: "2026-06-20T00:00:00Z",
			updated_at: "2026-06-20T00:00:00Z"
		}
	};
}

function createEvent(body: unknown, user: unknown = { id: "user-1", role: "user" }) {
	const url = "http://localhost/api/json-recipes/recipe-1/deploy";
	return {
		url: new URL(url),
		request: new Request(url, {
			method: "POST",
			headers: { "content-type": "application/json" },
			body: JSON.stringify(body)
		}),
		params: { id: "recipe-1" },
		fetch: vi.fn(),
		locals: { user, adminUser: user && (user as { role?: string }).role === "admin" ? user : null }
	} as unknown as RequestEvent;
}

async function responseJson(response: Response) {
	return (await response.json()) as unknown;
}

describe("JSON recipe deploy API endpoint", () => {
	beforeEach(() => {
		deployJsonRecipeMock.mockReset();
	});

	it("requires authentication", async () => {
		const response = await POST(createEvent({ user_id: "user-1" }, null));

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(deployJsonRecipeMock).not.toHaveBeenCalled();
	});

	it("requires ordinary users to deploy into their own account", async () => {
		const response = await POST(createEvent({ user_id: "other-user" }));

		expect(response.status).toBe(403);
		expect(await responseJson(response)).toEqual({ error: "cannot deploy recipe for another user" });
		expect(deployJsonRecipeMock).not.toHaveBeenCalled();
	});

	it("forwards valid deploy requests", async () => {
		const result = deployResponse();
		deployJsonRecipeMock.mockResolvedValue(result);
		const event = createEvent({ user_id: "user-1" });

		const response = await POST(event);

		expect(deployJsonRecipeMock).toHaveBeenCalledWith(event.fetch, "recipe-1", { userId: "user-1" });
		expect(response.status).toBe(201);
		expect(await responseJson(response)).toEqual(result);
	});

	it("allows admins to deploy for a specified user id", async () => {
		const result = deployResponse();
		deployJsonRecipeMock.mockResolvedValue(result);
		const event = createEvent({ user_id: "user-2" }, { id: "admin-1", role: "admin" });

		const response = await POST(event);

		expect(deployJsonRecipeMock).toHaveBeenCalledWith(event.fetch, "recipe-1", { userId: "user-2" });
		expect(response.status).toBe(201);
	});

	it("rejects invalid deploy payloads", async () => {
		const response = await POST(createEvent({ user_id: 42 }));

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "invalid JSON recipe deploy payload" });
		expect(deployJsonRecipeMock).not.toHaveBeenCalled();
	});
});
