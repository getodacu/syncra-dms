import { beforeEach, describe, expect, it, vi } from "vitest";

import { DatasetApiError } from "$lib/server/datasets";
import { DELETE, GET, PUT } from "./+server";
import type { RequestEvent } from "./$types";

const { deleteDatasetMock, getDatasetMock, updateDatasetMock, DatasetApiErrorMock } = vi.hoisted(
	() => {
		class MockDatasetApiError extends Error {
			status: number;

			constructor(status: number, message: string) {
				super(message);
				this.name = "DatasetApiError";
				this.status = status;
			}
		}

		return {
			deleteDatasetMock: vi.fn(),
			getDatasetMock: vi.fn(),
			updateDatasetMock: vi.fn(),
			DatasetApiErrorMock: MockDatasetApiError
		};
	}
);

vi.mock("$lib/server/datasets", () => ({
	deleteDataset: deleteDatasetMock,
	getDataset: getDatasetMock,
	updateDataset: updateDatasetMock,
	DatasetApiError: DatasetApiErrorMock,
	isDatasetApiError: (error: unknown) => error instanceof DatasetApiErrorMock
}));

function fieldFixture() {
	return { path: "/total", key: "total", label: "Total" };
}

function datasetFixture() {
	return {
		id: "dataset-1",
		created_at: "2026-06-01T00:00:00Z",
		updated_at: "2026-06-01T00:01:00Z",
		user_id: "user-1",
		name: "Invoices",
		schema_id: "schema-1",
		schema_name: "Invoice",
		selected_fields: [fieldFixture()],
		field_count: 1
	};
}

function createEvent(
	method: "DELETE" | "GET" | "PUT",
	body: unknown = undefined,
	user: unknown = { id: "user-1" }
) {
	return createRequestEvent(
		new Request("http://localhost/api/datasets/dataset-1", {
			method,
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
		},
		params: {
			id: "dataset-1"
		}
	} as unknown as RequestEvent;
}

async function responseJson(response: Response) {
	return (await response.json()) as unknown;
}

describe("dataset item API endpoint", () => {
	beforeEach(() => {
		deleteDatasetMock.mockReset();
		getDatasetMock.mockReset();
		updateDatasetMock.mockReset();
	});

	it("returns 401 for unauthenticated GET requests", async () => {
		const response = await GET(createEvent("GET", undefined, null));

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(getDatasetMock).not.toHaveBeenCalled();
	});

	it("returns 401 for unauthenticated PUT requests", async () => {
		const response = await PUT(
			createEvent(
				"PUT",
				{
					name: "Invoices",
					schema_id: "schema-1",
					selected_fields: [fieldFixture()]
				},
				null
			)
		);

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(updateDatasetMock).not.toHaveBeenCalled();
	});

	it("returns 401 for unauthenticated DELETE requests", async () => {
		const response = await DELETE(createEvent("DELETE", undefined, null));

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(deleteDatasetMock).not.toHaveBeenCalled();
	});

	it("calls the dataset service with current user id for GET requests", async () => {
		const dataset = datasetFixture();
		getDatasetMock.mockResolvedValue(dataset);
		const event = createEvent("GET");

		const response = await GET(event);

		expect(getDatasetMock).toHaveBeenCalledWith(event.fetch, "dataset-1", {
			userId: "user-1"
		});
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(dataset);
	});

	it("returns 400 for invalid PUT bodies", async () => {
		const response = await PUT(
			createEvent("PUT", {
				name: "Invoices",
				selected_fields: "total"
			})
		);

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({
			error: "selected_fields must be an array of dataset fields"
		});
		expect(updateDatasetMock).not.toHaveBeenCalled();
	});

	it("calls the dataset service with id, current user id, and input for PUT requests", async () => {
		const dataset = { ...datasetFixture(), name: "Receipts" };
		updateDatasetMock.mockResolvedValue(dataset);
		const selectedField = { path: "/ total ", key: " total ", label: " Total " };
		const event = createEvent("PUT", {
			name: "  Receipts  ",
			schema_id: " schema-1 ",
			selected_fields: [selectedField]
		});

		const response = await PUT(event);

		expect(updateDatasetMock).toHaveBeenCalledWith(
			event.fetch,
			"dataset-1",
			{
				name: "Receipts",
				schema_id: "schema-1",
				selected_fields: [selectedField]
			},
			{ userId: "user-1" }
		);
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(dataset);
	});

	it("returns 400 when schema_id is missing for PUT requests", async () => {
		const response = await PUT(
			createEvent("PUT", {
				name: "Invoices",
				selected_fields: [fieldFixture()]
			})
		);

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "schema_id is required" });
		expect(updateDatasetMock).not.toHaveBeenCalled();
	});

	it("returns 400 when schema_id is empty for PUT requests", async () => {
		const event = createEvent("PUT", {
			name: "Invoices",
			schema_id: "   ",
			selected_fields: [fieldFixture()]
		});

		const response = await PUT(event);

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "schema_id is required" });
		expect(updateDatasetMock).not.toHaveBeenCalled();
	});

	it("calls the dataset service with id and current user id for DELETE requests", async () => {
		deleteDatasetMock.mockResolvedValue({ deleted_id: "dataset-1" });
		const event = createEvent("DELETE");

		const response = await DELETE(event);

		expect(deleteDatasetMock).toHaveBeenCalledWith(event.fetch, "dataset-1", {
			userId: "user-1"
		});
		expect(response.status).toBe(204);
		expect(await response.text()).toBe("");
	});

	it("preserves dataset service client errors", async () => {
		getDatasetMock.mockRejectedValue(new DatasetApiError(404, "dataset not found"));

		const response = await GET(createEvent("GET"));

		expect(response.status).toBe(404);
		expect(await responseJson(response)).toEqual({ error: "dataset not found" });
	});

	it("normalizes dataset service server errors to 502", async () => {
		deleteDatasetMock.mockRejectedValue(
			new DatasetApiError(503, "Dataset service unavailable")
		);

		const response = await DELETE(createEvent("DELETE"));

		expect(response.status).toBe(502);
		expect(await responseJson(response)).toEqual({ error: "A server error occurred. Please try again." });
	});
});
