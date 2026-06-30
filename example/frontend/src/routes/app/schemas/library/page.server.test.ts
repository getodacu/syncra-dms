import { beforeEach, describe, expect, it, vi } from "vitest";

const { listJsonRecipesMock, publicErrorMessageMock, schemaApiErrorCtor } = vi.hoisted(() => {
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
		publicErrorMessageMock: vi.fn(),
		schemaApiErrorCtor: MockSchemaApiError
	};
});

vi.mock("$lib/server/public-errors", () => ({
	publicErrorMessage: publicErrorMessageMock
}));

vi.mock("$lib/server/schemas", () => ({
	SchemaApiError: schemaApiErrorCtor,
	isSchemaApiError: (error: unknown) => error instanceof schemaApiErrorCtor,
	listJsonRecipes: listJsonRecipesMock
}));

import { load } from "./+page.server";

function recipe() {
	return {
		id: "recipe-1",
		title: "Invoice",
		description: "Invoice fields",
		json: { type: "object" },
		counter: 0,
		category_id: null,
		category: null,
		created_at: "2026-06-20T00:00:00Z",
		updated_at: "2026-06-20T00:00:00Z"
	};
}

describe("schema library page load", () => {
	beforeEach(() => {
		listJsonRecipesMock.mockReset();
		publicErrorMessageMock.mockReset();
		publicErrorMessageMock.mockReturnValue("Unable to load OCR recipes.");
	});

	it("loads the first page of system OCR recipes for the authenticated schema library", async () => {
		const result = { recipes: [recipe()], next_cursor: null };
		const fetchMock = vi.fn();
		listJsonRecipesMock.mockResolvedValue(result);

		await expect(
			load({
				fetch: fetchMock,
				locals: {
					user: {
						id: "user-1"
					}
				}
			} as never)
		).resolves.toEqual({
			isLoggedIn: true,
			userId: "user-1",
			recipes: result.recipes,
			nextCursor: null,
			loadError: null
		});
		expect(listJsonRecipesMock).toHaveBeenCalledWith(fetchMock, { size: 100, sort: "desc" });
	});

	it("converts schema API failures to safe public load errors", async () => {
		listJsonRecipesMock.mockRejectedValue(new schemaApiErrorCtor(503, "database password leaked"));

		await expect(
			load({
				fetch: vi.fn(),
				locals: {
					user: {
						id: "user-1"
					}
				}
			} as never)
		).resolves.toEqual({
			isLoggedIn: true,
			userId: "user-1",
			recipes: [],
			nextCursor: null,
			loadError: "Unable to load OCR recipes."
		});
		expect(publicErrorMessageMock).toHaveBeenCalledWith(
			503,
			"database password leaked",
			"Unable to load OCR recipes."
		);
	});
});
