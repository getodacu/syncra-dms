import { describe, expect, it, vi } from "vitest";

import {
	createDataset,
	DatasetClientError,
	deleteDataset,
	exportDataset,
	fetchDatasetRows,
	fetchDatasets,
	getDataset,
	isDatasetClientError,
	isDatasetNotFoundError,
	parseContentDispositionFilename,
	updateDataset,
} from "./datasets";

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
		selected_fields: [
			{ path: "/invoice_number", key: "invoice_number", label: "Invoice number" },
			{ path: "/line_items", key: "line_items", label: "Line items" },
		],
		field_count: 2,
	};
}

function datasetInput() {
	return {
		name: "Invoice dataset",
		schema_id: "schema-1",
		selected_fields: [
			{ path: "/invoice_number", key: "invoice_number", label: "Invoice number" },
			{ path: "/line_items", key: "line_items", label: "Line items" },
		],
	};
}

function expectJsonContentType(init: RequestInit | undefined) {
	const headers = init?.headers as Headers;
	expect(headers.get("content-type")).toBe("application/json");
}

describe("datasets api client", () => {
	it("fetches datasets with default pagination through the SvelteKit proxy", async () => {
		const body = { datasets: [datasetResponse()], next_cursor: null };
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body));

		await expect(fetchDatasets(fetchFn)).resolves.toEqual(body);
		expect(fetchFn).toHaveBeenCalledWith("/api/datasets?size=100", { method: "GET" });
	});

	it("fetches datasets with cursor pagination options", async () => {
		const body = { datasets: [datasetResponse()], next_cursor: "cursor-2" };
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body));

		await expect(
			fetchDatasets(fetchFn, { cursor: "cursor-1", size: 50, sort: "asc" })
		).resolves.toEqual(body);
		expect(fetchFn).toHaveBeenCalledWith("/api/datasets?cursor=cursor-1&size=50&sort=asc", {
			method: "GET",
		});
	});

	it("creates datasets through the SvelteKit proxy with JSON", async () => {
		const saved = datasetResponse();
		const input = datasetInput();
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(saved, { status: 201 }));

		await expect(createDataset(fetchFn, input)).resolves.toEqual(saved);
		expect(fetchFn).toHaveBeenCalledWith(
			"/api/datasets",
			expect.objectContaining({
				method: "POST",
				body: JSON.stringify(input),
			})
		);
		expectJsonContentType(fetchFn.mock.calls[0]?.[1]);
	});

	it("updates datasets through the SvelteKit proxy with encoded ids and JSON", async () => {
		const saved = { ...datasetResponse("dataset/1"), name: "Updated dataset" };
		const input = { ...datasetInput(), name: "Updated dataset" };
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(saved));

		await expect(updateDataset(fetchFn, "dataset/1", input)).resolves.toEqual(saved);
		expect(fetchFn).toHaveBeenCalledWith(
			"/api/datasets/dataset%2F1",
			expect.objectContaining({
				method: "PUT",
				body: JSON.stringify(input),
			})
		);
		expectJsonContentType(fetchFn.mock.calls[0]?.[1]);
	});

	it("gets and deletes datasets with encoded ids", async () => {
		const dataset = datasetResponse("dataset/1");
		const fetchFn = vi
			.fn()
			.mockResolvedValueOnce(jsonResponse(dataset))
			.mockResolvedValueOnce(new Response(null, { status: 204 }));

		await expect(getDataset(fetchFn, "dataset/1")).resolves.toEqual(dataset);
		await expect(deleteDataset(fetchFn, "dataset/1")).resolves.toEqual({
			deleted_id: "dataset/1",
		});
		expect(fetchFn).toHaveBeenNthCalledWith(1, "/api/datasets/dataset%2F1", {
			method: "GET",
		});
		expect(fetchFn).toHaveBeenNthCalledWith(2, "/api/datasets/dataset%2F1", {
			method: "DELETE",
		});
	});

	it("fetches dataset rows with pagination and accepts nested JSON cell values", async () => {
		const body = {
			dataset: datasetResponse(),
			columns: [
				{ path: "/invoice_number", key: "invoice_number", label: "Invoice number" },
				{ path: "/line_items", key: "line_items", label: "Line items" },
				{ path: "/metadata", key: "metadata", label: "Metadata" },
			],
			rows: [
				{
					document_id: "document-1",
					filename: "invoice.pdf",
					created_at: "2026-06-02T00:00:00Z",
					values: {
						invoice_number: "INV-001",
						total: 42.5,
						paid: false,
						due_date: null,
						line_items: [
							{ description: "Service", amount: 40 },
							{ description: "Tax", amount: 2.5 },
						],
						metadata: { currency: "USD", tags: ["ocr", "reviewed"] },
					},
				},
			],
			next_cursor: "rows-2",
		};
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body));

		await expect(
			fetchDatasetRows(fetchFn, "dataset/1", { cursor: "rows-1", size: 25, sort: "desc" })
		).resolves.toEqual(body);
		expect(fetchFn).toHaveBeenCalledWith(
			"/api/datasets/dataset%2F1/rows?cursor=rows-1&size=25&sort=desc",
			{ method: "GET" }
		);
	});

	it("fetches dataset rows with created_at date bounds", async () => {
		const body = {
			dataset: datasetResponse(),
			columns: [{ path: "/invoice_number", key: "invoice_number", label: "Invoice number" }],
			rows: [],
			next_cursor: null,
		};
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body));

		await expect(
			fetchDatasetRows(fetchFn, "dataset/1", {
				createdFrom: "2026-06-01T00:00:00.000Z",
				createdTo: "2026-06-07T23:59:59.999Z",
				size: 25,
				sort: "desc",
			})
		).resolves.toEqual(body);
		expect(fetchFn).toHaveBeenCalledWith(
			"/api/datasets/dataset%2F1/rows?created_from=2026-06-01T00%3A00%3A00.000Z&created_to=2026-06-07T23%3A59%3A59.999Z&size=25&sort=desc",
			{ method: "GET" }
		);
	});

	it("exports dataset binaries while preserving useful headers and filename", async () => {
		const bytes = new Uint8Array([1, 2, 3]);
		const fetchFn = vi.fn().mockResolvedValue(
			new Response(bytes, {
				headers: {
					"content-type": "text/csv",
					"content-disposition": "attachment; filename*=UTF-8''invoice%20dataset.csv",
				},
			})
		);

		const result = await exportDataset(fetchFn, "dataset/1", { format: "csv", sort: "desc" });

		expect(fetchFn).toHaveBeenCalledWith(
			"/api/datasets/dataset%2F1/export?format=csv&sort=desc",
			{ method: "GET" }
		);
		expect(result.contentType).toBe("text/csv");
		expect(result.contentDisposition).toBe(
			"attachment; filename*=UTF-8''invoice%20dataset.csv"
		);
		expect(result.filename).toBe("invoice dataset.csv");
		expect(Array.from(new Uint8Array(await result.blob.arrayBuffer()))).toEqual([1, 2, 3]);
	});

	it("exports dataset binaries with created_at date bounds", async () => {
		const fetchFn = vi.fn().mockResolvedValue(new Response("document_id,filename\n"));

		await exportDataset(fetchFn, "dataset/1", {
			format: "csv",
			createdFrom: "2026-06-01T00:00:00.000Z",
			createdTo: "2026-06-07T23:59:59.999Z",
			sort: "asc",
		});

		expect(fetchFn).toHaveBeenCalledWith(
			"/api/datasets/dataset%2F1/export?format=csv&created_from=2026-06-01T00%3A00%3A00.000Z&created_to=2026-06-07T23%3A59%3A59.999Z&sort=asc",
			{ method: "GET" }
		);
	});

	it("parses encoded and escaped content disposition filenames", () => {
		expect(
			parseContentDispositionFilename("attachment; filename*=UTF-8'en'invoice%20dataset.csv")
		).toBe("invoice dataset.csv");
		expect(
			parseContentDispositionFilename('attachment; filename="invoice \\"final\\".csv"')
		).toBe('invoice "final".csv');
	});

	it("throws backend JSON messages for request errors", async () => {
		const fetchFn = vi
			.fn()
			.mockResolvedValue(jsonResponse({ error: "dataset name already exists" }, { status: 409 }));

		await expect(createDataset(fetchFn, datasetInput())).rejects.toThrow(
			"dataset name already exists"
		);
	});

	it("preserves backend response status on request errors", async () => {
		const fetchFn = vi.fn(() =>
			Promise.resolve(jsonResponse({ error: "dataset not found" }, { status: 404 }))
		);

		await expect(fetchDatasetRows(fetchFn, "missing-dataset")).rejects.toMatchObject({
			name: "DatasetClientError",
			message: "dataset not found",
			status: 404,
		});

		let error: unknown;
		try {
			await fetchDatasetRows(fetchFn, "missing-dataset");
		} catch (caught) {
			error = caught;
		}

		expect(error).toBeInstanceOf(DatasetClientError);
		expect(isDatasetClientError(error)).toBe(true);
		expect(isDatasetNotFoundError(error)).toBe(true);
	});

	it("uses fallback messages for export server errors", async () => {
		const fetchFn = vi
			.fn()
			.mockResolvedValue(jsonResponse({ error: "failed to export dataset" }, { status: 500 }));

		await expect(exportDataset(fetchFn, "dataset-1", { format: "csv" })).rejects.toThrow(
			"Failed to export dataset"
		);
	});

	it("rejects invalid list, dataset, row, and delete responses", async () => {
		await expect(
			fetchDatasets(vi.fn().mockResolvedValue(jsonResponse({ datasets: [{ id: "dataset-1" }] })))
		).rejects.toThrow("Invalid dataset list response");

		await expect(
			getDataset(vi.fn().mockResolvedValue(jsonResponse({ id: "dataset-1" })), "dataset-1")
		).rejects.toThrow("Invalid dataset response");

		await expect(
			fetchDatasetRows(
				vi.fn().mockResolvedValue(
					jsonResponse({
						dataset: datasetResponse(),
						columns: [{ key: "total", label: "Total", path: "/total" }],
						rows: [{ document_id: "document-1", filename: "invoice.pdf", values: [] }],
						next_cursor: null,
					})
				),
				"dataset-1"
			)
		).rejects.toThrow("Invalid dataset rows response");

		await expect(
			deleteDataset(vi.fn().mockResolvedValue(jsonResponse({ deleted_id: 1 })), "dataset-1")
		).rejects.toThrow("Invalid dataset delete response");
	});
});
