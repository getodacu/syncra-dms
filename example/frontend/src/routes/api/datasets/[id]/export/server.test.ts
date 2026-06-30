import { beforeEach, describe, expect, it, vi } from "vitest";

import { DatasetApiError } from "$lib/server/datasets";
import { GET } from "./+server";
import type { RequestEvent } from "./$types";

const { exportDatasetMock, DatasetApiErrorMock } = vi.hoisted(() => {
	class MockDatasetApiError extends Error {
		status: number;

		constructor(status: number, message: string) {
			super(message);
			this.name = "DatasetApiError";
			this.status = status;
		}
	}

	return {
		exportDatasetMock: vi.fn(),
		DatasetApiErrorMock: MockDatasetApiError
	};
});

vi.mock("$lib/server/datasets", () => ({
	exportDataset: exportDatasetMock,
	DatasetApiError: DatasetApiErrorMock,
	isDatasetApiError: (error: unknown) => error instanceof DatasetApiErrorMock
}));

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

describe("dataset export API endpoint", () => {
	beforeEach(() => {
		exportDatasetMock.mockReset();
	});

	it("returns 401 for unauthenticated GET requests", async () => {
		const response = await GET(
			createGetEvent("http://localhost/api/datasets/dataset-1/export", null)
		);

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(exportDatasetMock).not.toHaveBeenCalled();
	});

	it("rejects invalid export formats", async () => {
		const response = await GET(
			createGetEvent("http://localhost/api/datasets/dataset-1/export?format=pdf")
		);

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({ error: "format must be csv or xlsx" });
		expect(exportDatasetMock).not.toHaveBeenCalled();
	});

	it("returns the upstream export body and useful headers", async () => {
		exportDatasetMock.mockResolvedValue({
			body: new Response("document_id,filename\n").body,
			status: 200,
			headers: new Headers({
				"content-type": "text/csv; charset=utf-8",
				"content-disposition": 'attachment; filename="invoices.csv"'
			})
		});
		const event = createGetEvent(
			"http://localhost/api/datasets/dataset-1/export?user_id=attacker&format=csv&sort=desc"
		);

		const response = await GET(event);

		expect(exportDatasetMock).toHaveBeenCalledWith(event.fetch, "dataset-1", {
			userId: "user-1",
			format: "csv",
			sort: "desc"
		});
		expect(response.status).toBe(200);
		expect(response.headers.get("content-type")).toBe("text/csv; charset=utf-8");
		expect(response.headers.get("content-disposition")).toBe(
			'attachment; filename="invoices.csv"'
		);
		expect(await response.text()).toBe("document_id,filename\n");
	});

	it("forwards created_at date bounds with current user id", async () => {
		exportDatasetMock.mockResolvedValue({
			body: new Response("document_id,filename\n").body,
			status: 200,
			headers: new Headers({
				"content-type": "text/csv; charset=utf-8"
			})
		});
		const event = createGetEvent(
			"http://localhost/api/datasets/dataset-1/export?user_id=attacker&format=csv&created_from=2026-06-01T00%3A00%3A00.000Z&created_to=2026-06-07T23%3A59%3A59.999Z&sort=desc"
		);

		const response = await GET(event);

		expect(exportDatasetMock).toHaveBeenCalledWith(event.fetch, "dataset-1", {
			userId: "user-1",
			format: "csv",
			createdFrom: "2026-06-01T00:00:00.000Z",
			createdTo: "2026-06-07T23:59:59.999Z",
			sort: "desc"
		});
		expect(response.status).toBe(200);
		expect(await response.text()).toBe("document_id,filename\n");
	});

	it("normalizes export format before forwarding", async () => {
		exportDatasetMock.mockResolvedValue({
			body: new Response("xlsx").body,
			status: 200,
			headers: new Headers({
				"content-type": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
			})
		});
		const event = createGetEvent(
			"http://localhost/api/datasets/dataset-1/export?format=%20XLSX%20"
		);

		const response = await GET(event);

		expect(exportDatasetMock).toHaveBeenCalledWith(event.fetch, "dataset-1", {
			userId: "user-1",
			format: "xlsx"
		});
		expect(response.status).toBe(200);
		expect(await response.text()).toBe("xlsx");
	});

	it("defaults the export format to the backend", async () => {
		exportDatasetMock.mockResolvedValue({
			body: new Response("xlsx").body,
			status: 200,
			headers: new Headers({
				"content-type": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
			})
		});
		const event = createGetEvent("http://localhost/api/datasets/dataset-1/export?sort=asc");

		const response = await GET(event);

		expect(exportDatasetMock).toHaveBeenCalledWith(event.fetch, "dataset-1", {
			userId: "user-1",
			sort: "asc"
		});
		expect(response.status).toBe(200);
		expect(await response.text()).toBe("xlsx");
	});

	it("preserves dataset service client errors", async () => {
		exportDatasetMock.mockRejectedValue(new DatasetApiError(404, "dataset not found"));

		const response = await GET(
			createGetEvent("http://localhost/api/datasets/dataset-1/export?format=csv")
		);

		expect(response.status).toBe(404);
		expect(await responseJson(response)).toEqual({ error: "dataset not found" });
	});

	it("normalizes dataset service server errors to 502", async () => {
		exportDatasetMock.mockRejectedValue(
			new DatasetApiError(500, "failed to export dataset")
		);

		const response = await GET(
			createGetEvent("http://localhost/api/datasets/dataset-1/export?format=csv")
		);

		expect(response.status).toBe(502);
		expect(await responseJson(response)).toEqual({ error: "A server error occurred. Please try again." });
	});
});
