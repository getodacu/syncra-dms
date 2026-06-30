import { beforeEach, describe, expect, it, vi } from "vitest";

import { DatasetApiError } from "$lib/server/datasets";
import { GET } from "./+server";
import type { RequestEvent } from "./$types";

const { listDatasetRowsMock, DatasetApiErrorMock } = vi.hoisted(() => {
	class MockDatasetApiError extends Error {
		status: number;

		constructor(status: number, message: string) {
			super(message);
			this.name = "DatasetApiError";
			this.status = status;
		}
	}

	return {
		listDatasetRowsMock: vi.fn(),
		DatasetApiErrorMock: MockDatasetApiError
	};
});

vi.mock("$lib/server/datasets", () => ({
	listDatasetRows: listDatasetRowsMock,
	DatasetApiError: DatasetApiErrorMock,
	isDatasetApiError: (error: unknown) => error instanceof DatasetApiErrorMock
}));

function fieldFixture() {
	return { path: "/total", key: "total", label: "Total" };
}

function rowsFixture() {
	return {
		dataset: {
			id: "dataset-1",
			created_at: "2026-06-01T00:00:00Z",
			updated_at: "2026-06-01T00:01:00Z",
			user_id: "user-1",
			name: "Invoices",
			schema_id: "schema-1",
			schema_name: "Invoice",
			selected_fields: [fieldFixture()],
			field_count: 1
		},
		columns: [fieldFixture()],
		rows: [
			{
				document_id: "document-1",
				filename: "invoice.pdf",
				created_at: "2026-06-01T00:02:00Z",
				values: { total: 42 }
			}
		],
		next_cursor: "next-page"
	};
}

function createGetEvent(url: string, user: unknown = { id: "user-1" }) {
	return {
		url: new URL(url),
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

describe("dataset rows API endpoint", () => {
	beforeEach(() => {
		listDatasetRowsMock.mockReset();
	});

	it("returns 401 for unauthenticated GET requests", async () => {
		const response = await GET(
			createGetEvent("http://localhost/api/datasets/dataset-1/rows", null)
		);

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(listDatasetRowsMock).not.toHaveBeenCalled();
	});

	it("calls the dataset service with current user id and pagination parameters", async () => {
		const rows = rowsFixture();
		listDatasetRowsMock.mockResolvedValue(rows);
		const event = createGetEvent(
			"http://localhost/api/datasets/dataset-1/rows?user_id=attacker&size=10&sort=asc&cursor=abc"
		);

		const response = await GET(event);

		expect(listDatasetRowsMock).toHaveBeenCalledWith(event.fetch, "dataset-1", {
			userId: "user-1",
			cursor: "abc",
			size: "10",
			sort: "asc"
		});
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(rows);
	});

	it("forwards created_at date bounds with current user id", async () => {
		const rows = rowsFixture();
		listDatasetRowsMock.mockResolvedValue(rows);
		const event = createGetEvent(
			"http://localhost/api/datasets/dataset-1/rows?user_id=attacker&created_from=2026-06-01T00%3A00%3A00.000Z&created_to=2026-06-07T23%3A59%3A59.999Z&size=10&sort=desc"
		);

		const response = await GET(event);

		expect(listDatasetRowsMock).toHaveBeenCalledWith(event.fetch, "dataset-1", {
			userId: "user-1",
			createdFrom: "2026-06-01T00:00:00.000Z",
			createdTo: "2026-06-07T23:59:59.999Z",
			size: "10",
			sort: "desc"
		});
		expect(response.status).toBe(200);
		expect(await responseJson(response)).toEqual(rows);
	});

	it("preserves dataset service client errors", async () => {
		listDatasetRowsMock.mockRejectedValue(new DatasetApiError(404, "dataset not found"));

		const response = await GET(createGetEvent("http://localhost/api/datasets/dataset-1/rows"));

		expect(response.status).toBe(404);
		expect(await responseJson(response)).toEqual({ error: "dataset not found" });
	});

	it("normalizes dataset service server errors to 502", async () => {
		listDatasetRowsMock.mockRejectedValue(
			new DatasetApiError(503, "Dataset service unavailable")
		);

		const response = await GET(createGetEvent("http://localhost/api/datasets/dataset-1/rows"));

		expect(response.status).toBe(502);
		expect(await responseJson(response)).toEqual({ error: "A server error occurred. Please try again." });
	});
});
