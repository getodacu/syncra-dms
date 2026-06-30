import { describe, expect, it, vi } from "vitest";

import { downloadDatasetExport, fetchDatasetRows, fetchDatasets } from "./api";

function jsonResponse(body: unknown, init?: ResponseInit) {
	return new Response(JSON.stringify(body), {
		headers: { "content-type": "application/json" },
		...init,
	});
}

function datasetResponse(id = "dataset-1") {
	return {
		id,
		created_at: "2026-06-01T00:00:00Z",
		updated_at: "2026-06-01T00:01:00Z",
		user_id: "user-1",
		name: "Invoice dataset",
		schema_id: "schema-1",
		schema_name: "Invoice",
		selected_fields: [{ path: "/total", key: "total", label: "Total" }],
		field_count: 1,
	};
}

describe("datasets app api wrapper", () => {
	it("delegates dataset list requests to the browser client helper", async () => {
		const body = { datasets: [datasetResponse()], next_cursor: null };
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body));

		await expect(fetchDatasets(fetchFn, { sort: "desc" })).resolves.toEqual(body);
		expect(fetchFn).toHaveBeenCalledWith("/api/datasets?size=100&sort=desc", {
			method: "GET",
		});
	});

	it("delegates dataset row requests to the browser client helper", async () => {
		const body = {
			dataset: datasetResponse(),
			columns: [{ path: "/total", key: "total", label: "Total" }],
			rows: [
				{
					document_id: "document-1",
					filename: "invoice.pdf",
					created_at: "2026-06-02T00:00:00Z",
					values: { total: 100 },
				},
			],
			next_cursor: null,
		};
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body));

		await expect(fetchDatasetRows(fetchFn, "dataset/1")).resolves.toEqual(body);
		expect(fetchFn).toHaveBeenCalledWith("/api/datasets/dataset%2F1/rows?size=100", {
			method: "GET",
		});
	});

	it("delegates dataset row requests with created_at date bounds", async () => {
		const body = {
			dataset: datasetResponse(),
			columns: [{ path: "/total", key: "total", label: "Total" }],
			rows: [],
			next_cursor: null,
		};
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body));

		await expect(
			fetchDatasetRows(fetchFn, "dataset-1", {
				createdFrom: "2026-06-01T00:00:00.000Z",
				createdTo: "2026-06-07T23:59:59.999Z",
				sort: "desc",
			})
		).resolves.toEqual(body);
		expect(fetchFn).toHaveBeenCalledWith(
			"/api/datasets/dataset-1/rows?created_from=2026-06-01T00%3A00%3A00.000Z&created_to=2026-06-07T23%3A59%3A59.999Z&size=100&sort=desc",
			{ method: "GET" }
		);
	});

	it("downloads dataset export binary data with response metadata", async () => {
		const fetchFn = vi.fn().mockResolvedValue(
			new Response("invoice,total\n1,100\n", {
				headers: {
					"content-type": "text/csv",
					"content-disposition": 'attachment; filename="invoices.csv"',
				},
			})
		);

		const result = await downloadDatasetExport(fetchFn, "dataset/1", "csv");

		expect(fetchFn).toHaveBeenCalledWith("/api/datasets/dataset%2F1/export?format=csv", {
			method: "GET",
		});
		expect(result.filename).toBe("invoices.csv");
		expect(result.contentType).toBe("text/csv");
		await expect(result.blob.text()).resolves.toBe("invoice,total\n1,100\n");
	});

	it("preserves dataset export options alongside the selected format", async () => {
		const fetchFn = vi.fn().mockResolvedValue(new Response("id,total\n1,100\n"));

		await downloadDatasetExport(fetchFn, "dataset-1", "csv", {
			createdFrom: "2026-06-01T00:00:00.000Z",
			createdTo: "2026-06-07T23:59:59.999Z",
			sort: "asc",
		});

		expect(fetchFn).toHaveBeenCalledWith(
			"/api/datasets/dataset-1/export?format=csv&created_from=2026-06-01T00%3A00%3A00.000Z&created_to=2026-06-07T23%3A59%3A59.999Z&sort=asc",
			{
				method: "GET",
			}
		);
	});

	it("uses fallback messages for dataset export server errors", async () => {
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse({ error: "Export unavailable" }, { status: 500 }));

		await expect(downloadDatasetExport(fetchFn, "dataset-1", "xlsx")).rejects.toThrow(
			"Failed to export dataset"
		);
		expect(fetchFn).toHaveBeenCalledWith("/api/datasets/dataset-1/export?format=xlsx", {
			method: "GET",
		});
	});
});
