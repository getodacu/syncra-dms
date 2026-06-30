import { beforeEach, describe, expect, it, vi } from "vitest";

import { DatasetApiError } from "$lib/server/datasets";
import { GET, POST } from "./+server";
import type { RequestEvent } from "./$types";

const { createDatasetMock, listDatasetsPageMock, DatasetApiErrorMock } = vi.hoisted(() => {
	class MockDatasetApiError extends Error {
		status: number;

		constructor(status: number, message: string) {
			super(message);
			this.name = "DatasetApiError";
			this.status = status;
		}
	}

	return {
		createDatasetMock: vi.fn(),
		listDatasetsPageMock: vi.fn(),
		DatasetApiErrorMock: MockDatasetApiError
	};
});

vi.mock("$lib/server/datasets", () => ({
	createDataset: createDatasetMock,
	listDatasetsPage: listDatasetsPageMock,
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

function createPostEvent(body: unknown, user: unknown = { id: "user-1" }) {
	return createPostRequestEvent(
		new Request("http://localhost/api/datasets", {
			method: "POST",
			body: body === undefined ? undefined : JSON.stringify(body)
		}),
		user
	);
}

function createPostRequestEvent(request: Request, user: unknown = { id: "user-1" }) {
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

async function responseJson(response: Response) {
	return (await response.json()) as unknown;
}

describe("dataset API endpoint", () => {
	beforeEach(() => {
		createDatasetMock.mockReset();
		listDatasetsPageMock.mockReset();
	});

	it("returns 401 for unauthenticated GET requests", async () => {
		const response = await GET(createGetEvent("http://localhost/api/datasets", null));

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(listDatasetsPageMock).not.toHaveBeenCalled();
	});

	it("returns 401 for unauthenticated POST requests", async () => {
		const response = await POST(
			createPostEvent(
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
		expect(createDatasetMock).not.toHaveBeenCalled();
	});

	it("calls the dataset service with current user id and pagination parameters", async () => {
		const page = { datasets: [datasetFixture()], next_cursor: "next-page" };
		listDatasetsPageMock.mockResolvedValue(page);
		const event = createGetEvent(
			"http://localhost/api/datasets?user_id=attacker&size=25&sort=desc&cursor=abc"
		);

		const response = await GET(event);

		expect(listDatasetsPageMock).toHaveBeenCalledWith(event.fetch, {
			userId: "user-1",
			cursor: "abc",
			size: "25",
			sort: "desc"
		});
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(page);
	});

	it("returns 400 for invalid POST bodies", async () => {
		const response = await POST(
			createPostEvent({
				name: "Invoices",
				schema_id: "schema-1",
				selected_fields: [{ path: "/total", key: 42, label: "Total" }]
			})
		);

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({
			error: "selected_fields must contain path, key, and label strings"
		});
		expect(createDatasetMock).not.toHaveBeenCalled();
	});

	it("returns 400 when schema_id is missing for create requests", async () => {
		const response = await POST(
			createPostEvent({
				name: "Invoices",
				selected_fields: [fieldFixture()]
			})
		);

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "schema_id is required" });
		expect(createDatasetMock).not.toHaveBeenCalled();
	});

	it("calls the dataset service with normalized ids and unchanged selected fields", async () => {
		const dataset = datasetFixture();
		createDatasetMock.mockResolvedValue(dataset);
		const selectedField = { path: "/ total ", key: " total ", label: " Total " };
		const event = createPostEvent({
			user_id: "attacker",
			name: "  Invoices  ",
			schema_id: "  schema-1  ",
			selected_fields: [selectedField]
		});

		const response = await POST(event);

		expect(createDatasetMock).toHaveBeenCalledWith(
			event.fetch,
			{
				name: "Invoices",
				schema_id: "schema-1",
				selected_fields: [selectedField]
			},
			{ userId: "user-1" }
		);
		expect(response.status).toBe(201);
		expect(await responseJson(response)).toEqual(dataset);
	});

	it("preserves dataset service client errors", async () => {
		listDatasetsPageMock.mockRejectedValue(new DatasetApiError(400, "invalid cursor"));

		const response = await GET(createGetEvent("http://localhost/api/datasets"));

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "invalid cursor" });
	});

	it("normalizes dataset service server errors to 502", async () => {
		createDatasetMock.mockRejectedValue(
			new DatasetApiError(503, "Dataset service unavailable")
		);

		const response = await POST(
			createPostEvent({
				name: "Invoices",
				schema_id: "schema-1",
				selected_fields: [fieldFixture()]
			})
		);

		expect(response.status).toBe(502);
		expect(await responseJson(response)).toEqual({ error: "A server error occurred. Please try again." });
	});
});
