import { beforeEach, describe, expect, it, vi } from "vitest";

const INTERNAL_API_HEADER = "X-Syncra-Internal-Token";

const privateEnv = vi.hoisted(() => ({}) as Record<string, string | undefined>);

vi.mock("$env/dynamic/private", () => ({ env: privateEnv }));

function datasetFieldFixture() {
	return { path: "/total", key: "total", label: "Total" };
}

function datasetFixture(id = "dataset-1") {
	return {
		id,
		created_at: "2026-06-01T00:00:00Z",
		updated_at: "2026-06-01T00:01:00Z",
		user_id: "user-1",
		name: "Invoices",
		schema_id: "schema-1",
		schema_name: "Invoice",
		selected_fields: [datasetFieldFixture()],
		field_count: 1
	};
}

function rowsFixture() {
	return {
		dataset: datasetFixture(),
		columns: [datasetFieldFixture()],
		rows: [
			{
				document_id: "document-1",
				filename: "invoice.pdf",
				created_at: "2026-06-01T00:02:00Z",
				values: { total: 42 }
			}
		],
		next_cursor: null
	};
}

function jsonResponse(body: unknown, init?: ResponseInit) {
	return new Response(JSON.stringify(body), {
		headers: { "content-type": "application/json" },
		...init
	});
}

function expectJsonContentType(init: RequestInit | undefined) {
	expect(new Headers(init?.headers).get("content-type")).toBe("application/json");
}

function expectInternalToken(init: RequestInit | undefined) {
	expect(new Headers(init?.headers).get(INTERNAL_API_HEADER)).toBe("internal-token");
}

describe("frontend dataset server helper", () => {
	beforeEach(() => {
		privateEnv.SYNCRA_API_BASE_URL = "http://dataset-api.test/";
		privateEnv.SYNCRA_INTERNAL_API_TOKEN = "internal-token";
	});

	it("lists user datasets through the backend with pagination", async () => {
		const { listDatasetsPage } = await import("./datasets");
		const page = { datasets: [datasetFixture()], next_cursor: "next-page" };
		const fetchMock = vi.fn().mockResolvedValue(jsonResponse(page));

		await expect(
			listDatasetsPage(fetchMock, {
				userId: "user-1",
				cursor: "cursor-1",
				size: 20,
				sort: "desc"
			})
		).resolves.toEqual(page);
		expect(fetchMock).toHaveBeenCalledWith(
			"http://dataset-api.test/api/datasets?user_id=user-1&cursor=cursor-1&size=20&sort=desc",
			expect.objectContaining({ method: "GET" })
		);
		expectInternalToken(fetchMock.mock.calls[0]?.[1]);
	});

	it("creates datasets with injected user id", async () => {
		const { createDataset } = await import("./datasets");
		const saved = datasetFixture();
		const fetchMock = vi.fn().mockResolvedValue(jsonResponse(saved, { status: 201 }));

		await expect(
			createDataset(
				fetchMock,
				{
					name: "Invoices",
					schema_id: "schema-1",
					selected_fields: [datasetFieldFixture()]
				},
				{ userId: "user-1" }
			)
		).resolves.toEqual(saved);
		expect(fetchMock).toHaveBeenCalledWith(
			"http://dataset-api.test/api/datasets",
			expect.objectContaining({
				method: "POST",
				body: JSON.stringify({
					name: "Invoices",
					schema_id: "schema-1",
					selected_fields: [datasetFieldFixture()],
					user_id: "user-1"
				})
			})
		);
		expectJsonContentType(fetchMock.mock.calls[0]?.[1]);
		expectInternalToken(fetchMock.mock.calls[0]?.[1]);
	});

	it("gets datasets by id with user scope", async () => {
		const { getDataset } = await import("./datasets");
		const dataset = datasetFixture("dataset/1");
		const fetchMock = vi.fn().mockResolvedValue(jsonResponse(dataset));

		await expect(getDataset(fetchMock, "dataset/1", { userId: "user-1" })).resolves.toEqual(
			dataset
		);
		expect(fetchMock).toHaveBeenCalledWith(
			"http://dataset-api.test/api/datasets/dataset%2F1?user_id=user-1",
			expect.objectContaining({ method: "GET" })
		);
		expectInternalToken(fetchMock.mock.calls[0]?.[1]);
	});

	it("updates datasets by id with user scope", async () => {
		const { updateDataset } = await import("./datasets");
		const saved = { ...datasetFixture(), name: "Receipts" };
		const fetchMock = vi.fn().mockResolvedValue(jsonResponse(saved));

		await expect(
			updateDataset(
				fetchMock,
				"dataset/1",
				{
					name: "Receipts",
					schema_id: "schema-1",
					selected_fields: [datasetFieldFixture()]
				},
				{ userId: "user-1" }
			)
		).resolves.toEqual(saved);
		expect(fetchMock).toHaveBeenCalledWith(
			"http://dataset-api.test/api/datasets/dataset%2F1?user_id=user-1",
			expect.objectContaining({
				method: "PUT",
				body: JSON.stringify({
					name: "Receipts",
					schema_id: "schema-1",
					selected_fields: [datasetFieldFixture()]
				})
			})
		);
		expectJsonContentType(fetchMock.mock.calls[0]?.[1]);
		expectInternalToken(fetchMock.mock.calls[0]?.[1]);
	});

	it("deletes datasets and accepts 204 responses", async () => {
		const { deleteDataset } = await import("./datasets");
		const fetchMock = vi.fn().mockResolvedValue(new Response(null, { status: 204 }));

		await expect(deleteDataset(fetchMock, "dataset-1", { userId: "user-1" })).resolves.toEqual({
			deleted_id: "dataset-1"
		});
		expect(fetchMock).toHaveBeenCalledWith(
			"http://dataset-api.test/api/datasets/dataset-1?user_id=user-1",
			expect.objectContaining({ method: "DELETE" })
		);
		expectInternalToken(fetchMock.mock.calls[0]?.[1]);
	});

	it("lists dataset rows with user scope and pagination", async () => {
		const { listDatasetRows } = await import("./datasets");
		const rows = rowsFixture();
		const fetchMock = vi.fn().mockResolvedValue(jsonResponse(rows));

		await expect(
			listDatasetRows(fetchMock, "dataset-1", {
				userId: "user-1",
				cursor: "cursor-1",
				size: "10",
				sort: "asc"
			})
		).resolves.toEqual(rows);
		expect(fetchMock).toHaveBeenCalledWith(
			"http://dataset-api.test/api/datasets/dataset-1/rows?user_id=user-1&cursor=cursor-1&size=10&sort=asc",
			expect.objectContaining({ method: "GET" })
		);
		expectInternalToken(fetchMock.mock.calls[0]?.[1]);
	});

	it("lists dataset rows with created_at date bounds", async () => {
		const { listDatasetRows } = await import("./datasets");
		const rows = rowsFixture();
		const fetchMock = vi.fn().mockResolvedValue(jsonResponse(rows));

		await expect(
			listDatasetRows(fetchMock, "dataset-1", {
				userId: "user-1",
				createdFrom: "2026-06-01T00:00:00.000Z",
				createdTo: "2026-06-07T23:59:59.999Z",
				size: "10",
				sort: "asc"
			})
		).resolves.toEqual(rows);
		expect(fetchMock).toHaveBeenCalledWith(
			"http://dataset-api.test/api/datasets/dataset-1/rows?user_id=user-1&created_from=2026-06-01T00%3A00%3A00.000Z&created_to=2026-06-07T23%3A59%3A59.999Z&size=10&sort=asc",
			expect.objectContaining({ method: "GET" })
		);
		expectInternalToken(fetchMock.mock.calls[0]?.[1]);
	});

	it("accepts nested JSON values in dataset rows", async () => {
		const { listDatasetRows } = await import("./datasets");
		const rows = {
			...rowsFixture(),
			rows: [
				{
					document_id: "document-1",
					filename: "invoice.pdf",
					created_at: "2026-06-01T00:02:00Z",
					values: {
						customer: { name: "Ada" },
						line_items: [{ description: "Service", quantity: 1 }],
						total: 42,
						notes: null
					}
				}
			]
		};
		const fetchMock = vi.fn().mockResolvedValue(jsonResponse(rows));

		await expect(listDatasetRows(fetchMock, "dataset-1", { userId: "user-1" })).resolves.toEqual(
			rows
		);
	});

	it("exports datasets without parsing the upstream body", async () => {
		const { exportDataset } = await import("./datasets");
		const fetchMock = vi.fn().mockResolvedValue(
			new Response("document_id,filename\n", {
				headers: {
					"content-type": "text/csv; charset=utf-8",
					"content-disposition": 'attachment; filename="invoices.csv"',
					"x-ignored": "ignored"
				}
			})
		);

		const exported = await exportDataset(fetchMock, "dataset-1", {
			userId: "user-1",
			format: "csv",
			sort: "desc"
		});

		expect(fetchMock).toHaveBeenCalledWith(
			"http://dataset-api.test/api/datasets/dataset-1/export?user_id=user-1&format=csv&sort=desc",
			expect.objectContaining({ method: "GET" })
		);
		expectInternalToken(fetchMock.mock.calls[0]?.[1]);
		expect(exported.status).toBe(200);
		expect(exported.headers.get("content-type")).toBe("text/csv; charset=utf-8");
		expect(exported.headers.get("content-disposition")).toBe(
			'attachment; filename="invoices.csv"'
		);
		expect(exported.headers.has("x-ignored")).toBe(false);
		await expect(new Response(exported.body).text()).resolves.toBe("document_id,filename\n");
	});

	it("exports datasets with created_at date bounds", async () => {
		const { exportDataset } = await import("./datasets");
		const fetchMock = vi.fn().mockResolvedValue(new Response("document_id,filename\n"));

		await exportDataset(fetchMock, "dataset-1", {
			userId: "user-1",
			format: "csv",
			createdFrom: "2026-06-01T00:00:00.000Z",
			createdTo: "2026-06-07T23:59:59.999Z",
			sort: "desc"
		});

		expect(fetchMock).toHaveBeenCalledWith(
			"http://dataset-api.test/api/datasets/dataset-1/export?user_id=user-1&format=csv&created_from=2026-06-01T00%3A00%3A00.000Z&created_to=2026-06-07T23%3A59%3A59.999Z&sort=desc",
			expect.objectContaining({ method: "GET" })
		);
	});

	it("rejects dataset requests before calling Go when the internal API token is missing", async () => {
		const { listDatasetsPage } = await import("./datasets");
		delete privateEnv.SYNCRA_INTERNAL_API_TOKEN;
		const fetchMock = vi.fn();

		await expect(listDatasetsPage(fetchMock)).rejects.toMatchObject({
			status: 500,
			message: "Dataset service is not configured"
		});
		expect(fetchMock).not.toHaveBeenCalled();
	});

	it("throws typed errors for successful invalid list payloads", async () => {
		const { listDatasetsPage } = await import("./datasets");
		const fetchMock = vi
			.fn()
			.mockResolvedValue(jsonResponse({ datasets: [{ id: "dataset-1" }], next_cursor: null }));

		await expect(listDatasetsPage(fetchMock, { userId: "user-1" })).rejects.toMatchObject({
			status: 502,
			message: "Invalid dataset response"
		});
	});

	it("throws typed errors for invalid row payloads", async () => {
		const { listDatasetRows } = await import("./datasets");
		const fetchMock = vi
			.fn()
			.mockResolvedValue(jsonResponse({ ...rowsFixture(), rows: [{ document_id: "document-1" }] }));

		await expect(listDatasetRows(fetchMock, "dataset-1", { userId: "user-1" })).rejects.toMatchObject({
			status: 502,
			message: "Invalid dataset rows response"
		});
	});

	it("throws typed dataset API errors from backend error responses", async () => {
		const { createDataset, isDatasetApiError } = await import("./datasets");
		const fetchMock = vi
			.fn()
			.mockResolvedValue(jsonResponse({ error: "invalid selected_fields" }, { status: 400 }));

		await expect(
			createDataset(
				fetchMock,
				{ name: "Invoices", schema_id: "schema-1", selected_fields: [datasetFieldFixture()] },
				{ userId: "user-1" }
			)
		).rejects.toMatchObject({ status: 400, message: "invalid selected_fields" });
		await createDataset(
			fetchMock,
			{ name: "Invoices", schema_id: "schema-1", selected_fields: [datasetFieldFixture()] },
			{ userId: "user-1" }
		).catch((error) => {
			expect(isDatasetApiError(error)).toBe(true);
		});
	});

	it("throws typed errors from export error responses", async () => {
		const { exportDataset } = await import("./datasets");
		const fetchMock = vi
			.fn()
			.mockResolvedValue(jsonResponse({ error: "failed to export dataset" }, { status: 500 }));

		await expect(
			exportDataset(fetchMock, "dataset-1", { userId: "user-1", format: "csv" })
		).rejects.toMatchObject({
			status: 500,
			message: "failed to export dataset"
		});
	});
});
