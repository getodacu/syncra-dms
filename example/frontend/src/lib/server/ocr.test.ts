import { beforeEach, describe, expect, it, vi } from "vitest";

type FetchInput = Parameters<typeof fetch>[0];
type FetchInit = Parameters<typeof fetch>[1];

const INTERNAL_API_HEADER = "X-Syncra-Internal-Token";

const { privateEnv } = vi.hoisted(() => ({
	privateEnv: {} as Record<string, string | undefined>
}));

vi.mock("$env/dynamic/private", () => ({ env: privateEnv }));

function requestHeaders(init: FetchInit | undefined) {
	return new Headers(init?.headers);
}

function validOCRResponse() {
	return {
		id: "document-1",
		created_at: "2026-05-27T00:00:00Z",
		updated_at: "2026-05-27T00:00:00Z",
		user_id: "user-1",
		original_filename: "invoice.pdf",
		mime_type: "application/pdf",
		file_size: 12,
		page_count: 2,
		document_hash: "abcdef",
		schema_id: "schema-1",
		has_inline_schema: false,
		markdown: "# Invoice",
		annotation_json: { total: 10 },
		cached: false
	};
}

function validOCRDocumentListResponse() {
	return {
		documents: [
			{
				id: "document-1",
				created_at: "2026-05-27T00:00:00Z",
				updated_at: "2026-05-27T00:01:00Z",
				user_id: "user-1",
				original_filename: "invoice.pdf",
				mime_type: "application/pdf",
				file_size: 12,
				page_count: 2,
				document_hash: "abcdef",
				schema_id: "schema-1",
				has_inline_schema: false,
				collections: [{ id: "collection-1", name: "Invoices" }]
			}
		],
		next_cursor: "cursor-1"
	};
}

function validOCRJobResponse() {
	return {
		id: "job-1",
		created_at: "2026-05-27T00:00:00Z",
		original_filename: "invoice.pdf",
		mime_type: "application/pdf",
		file_size: 12,
		page_count: 2,
		schema_id: "schema-1",
		schema_name: "Invoice",
		has_inline_schema: false,
		document_id: "document-1",
		status: "completed"
	};
}

function validOCRJobListResponse() {
	return {
		jobs: [validOCRJobResponse()],
		next_cursor: "cursor-1"
	};
}

function validDeleteOCRDocumentsResponse() {
	return {
		deleted_ids: ["document-1"],
		deleted_count: 1
	};
}

function validDeleteOCRJobsResponse() {
	return {
		deleted_ids: ["job-1"],
		deleted_count: 1
	};
}

describe("frontend OCR server helper", () => {
	beforeEach(() => {
		vi.resetModules();
		for (const key of Object.keys(privateEnv)) delete privateEnv[key];
		process.env.SYNCRA_API_BASE_URL = "http://ocr-api.test/";
		process.env.SYNCRA_INTERNAL_API_TOKEN = "internal-token";
		process.env.NODE_ENV = "test";
	});

	it("gets OCR documents through the backend", async () => {
		const { getOCRDocument } = await import("./ocr");
		const responseBody = validOCRResponse();
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(responseBody), { status: 200 });
		});

		const result = await getOCRDocument(fetchMock, "document-1", { userId: "user-1" });

		expect(result).toEqual(responseBody);
		expect(fetchMock).toHaveBeenCalledWith(
			"http://ocr-api.test/api/ocr/document/document-1?user_id=user-1",
			expect.objectContaining({ method: "GET" })
		);
		expect(requestHeaders(fetchMock.mock.calls[0][1]).get(INTERNAL_API_HEADER)).toBe("internal-token");
	});

	it("throws typed OCR API errors from document lookup backend error responses", async () => {
		const { getOCRDocument, isOCRApiError } = await import("./ocr");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify({ error: "OCR document not found" }), { status: 404 });
		});

		await expect(getOCRDocument(fetchMock, "missing-document", { userId: "user-1" }))
			.rejects.toMatchObject({
				status: 404,
				message: "OCR document not found"
			});

		try {
			await getOCRDocument(fetchMock, "missing-document", { userId: "user-1" });
		} catch (error) {
			expect(isOCRApiError(error)).toBe(true);
		}
	});

	it("rejects invalid successful OCR document lookup responses", async () => {
		const { getOCRDocument } = await import("./ocr");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify({ id: "document-1" }), { status: 200 });
		});

		await expect(getOCRDocument(fetchMock, "document-1")).rejects.toMatchObject({
			status: 502,
			message: "Invalid OCR response"
		});
	});

	it("deletes an OCR document through the backend", async () => {
		const { deleteOCRDocument } = await import("./ocr");
		const responseBody = validDeleteOCRDocumentsResponse();
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(null, { status: 204 });
		});

		const result = await deleteOCRDocument(fetchMock, "document-1", { userId: "user-1" });

		expect(result).toEqual(responseBody);
		expect(fetchMock).toHaveBeenCalledWith(
			"http://ocr-api.test/api/ocr/documents/document-1?user_id=user-1",
			expect.objectContaining({ method: "DELETE" })
		);
	});

	it("throws typed OCR API errors from document delete backend error responses", async () => {
		const { deleteOCRDocument, isOCRApiError } = await import("./ocr");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify({ error: "OCR document not found" }), { status: 404 });
		});

		await expect(deleteOCRDocument(fetchMock, "missing-document", { userId: "user-1" }))
			.rejects.toMatchObject({
				status: 404,
				message: "OCR document not found"
			});

		try {
			await deleteOCRDocument(fetchMock, "missing-document", { userId: "user-1" });
		} catch (error) {
			expect(isOCRApiError(error)).toBe(true);
		}
	});

	it("rejects invalid successful OCR document delete responses", async () => {
		const { deleteOCRDocument } = await import("./ocr");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify({ deleted_ids: [1], deleted_count: "1" }), { status: 200 });
		});

		await expect(deleteOCRDocument(fetchMock, "document-1")).rejects.toMatchObject({
			status: 502,
			message: "Invalid OCR document delete response"
		});
	});

	it("updates OCR document metadata through the backend", async () => {
		const { updateOCRDocument } = await import("./ocr");
		const responseBody = { ...validOCRResponse(), original_filename: "renamed.pdf" };
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(responseBody), { status: 200 });
		});

		const result = await updateOCRDocument(
			fetchMock,
			"document-1",
			{ originalFilename: "renamed.pdf" },
			{ userId: "user-1" }
		);

		expect(result).toEqual(responseBody);
		expect(fetchMock).toHaveBeenCalledWith(
			"http://ocr-api.test/api/ocr/documents/document-1?user_id=user-1",
			expect.objectContaining({
				method: "PATCH",
				body: JSON.stringify({ original_filename: "renamed.pdf" })
			})
		);
		expect(requestHeaders(fetchMock.mock.calls[0][1]).get("content-type")).toBe("application/json");
		expect(requestHeaders(fetchMock.mock.calls[0][1]).get(INTERNAL_API_HEADER)).toBe("internal-token");
	});

	it("throws typed OCR API errors from document update backend error responses", async () => {
		const { updateOCRDocument } = await import("./ocr");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify({ error: "OCR document not found" }), { status: 404 });
		});

		await expect(
			updateOCRDocument(
				fetchMock,
				"missing-document",
				{ originalFilename: "renamed.pdf" },
				{ userId: "user-1" }
			)
		).rejects.toMatchObject({
			status: 404,
			message: "OCR document not found"
		});
	});

	it("rejects invalid successful OCR document update responses", async () => {
		const { updateOCRDocument } = await import("./ocr");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify({ id: "document-1" }), { status: 200 });
		});

		await expect(
			updateOCRDocument(fetchMock, "document-1", { originalFilename: "renamed.pdf" })
		).rejects.toMatchObject({
			status: 502,
			message: "Invalid OCR document update response"
		});
	});

	it("lists OCR documents through the backend with filters", async () => {
		const { listOCRDocuments } = await import("./ocr");
		const responseBody = validOCRDocumentListResponse();
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(responseBody), { status: 200 });
		});

		const result = await listOCRDocuments(fetchMock, {
			userId: "user-1",
			filename: "invoice",
			createdFrom: "2026-05-27T00:00:00Z",
			createdTo: "2026-05-28T00:00:00Z",
			cursor: "cursor-0",
			size: 50,
			sort: "asc"
		});

		expect(result).toEqual(responseBody);
		expect(fetchMock).toHaveBeenCalledWith(
			"http://ocr-api.test/api/ocr/documents?user_id=user-1&filename=invoice&created_from=2026-05-27T00%3A00%3A00Z&created_to=2026-05-28T00%3A00%3A00Z&cursor=cursor-0&size=50&sort=asc",
			expect.objectContaining({ method: "GET" })
		);
	});

	it("lists OCR documents through the backend with collection filters", async () => {
		const { listOCRDocuments } = await import("./ocr");
		const responseBody = validOCRDocumentListResponse();
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(responseBody), { status: 200 });
		});

		await listOCRDocuments(fetchMock, { userId: "user-1", collectionId: "collection-1" });

		expect(fetchMock).toHaveBeenCalledWith(
			"http://ocr-api.test/api/ocr/documents?user_id=user-1&collection_id=collection-1",
			expect.objectContaining({ method: "GET" })
		);
		expect(requestHeaders(fetchMock.mock.calls[0][1]).get(INTERNAL_API_HEADER)).toBe("internal-token");
	});

	it("throws typed OCR API errors from document list backend error responses", async () => {
		const { listOCRDocuments, isOCRApiError } = await import("./ocr");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify({ error: "invalid cursor" }), { status: 400 });
		});

		await expect(listOCRDocuments(fetchMock, { userId: "user-1", cursor: "bad" }))
			.rejects.toMatchObject({
				status: 400,
				message: "invalid cursor"
			});

		try {
			await listOCRDocuments(fetchMock, { userId: "user-1", cursor: "bad" });
		} catch (error) {
			expect(isOCRApiError(error)).toBe(true);
		}
	});

	it("rejects invalid successful OCR document list responses", async () => {
		const { listOCRDocuments } = await import("./ocr");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify({ documents: [{ id: "document-1" }], next_cursor: null }), {
				status: 200
			});
		});

		await expect(listOCRDocuments(fetchMock)).rejects.toMatchObject({
			status: 502,
			message: "Invalid OCR document list response"
		});
	});

	it("bulk deletes OCR documents through the backend", async () => {
		const { deleteOCRDocuments } = await import("./ocr");
		const responseBody = {
			deleted_ids: ["document-1", "document-2"],
			deleted_count: 2
		};
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(responseBody), { status: 200 });
		});

		const result = await deleteOCRDocuments(fetchMock, ["document-1", "document-2"], {
			userId: "user-1"
		});

		expect(result).toEqual(responseBody);
		expect(fetchMock).toHaveBeenCalledWith(
			"http://ocr-api.test/api/ocr/documents?user_id=user-1",
			expect.objectContaining({
				method: "DELETE",
				body: JSON.stringify({ ids: ["document-1", "document-2"] })
			})
		);
		expect(requestHeaders(fetchMock.mock.calls[0][1]).get("content-type")).toBe("application/json");
		expect(requestHeaders(fetchMock.mock.calls[0][1]).get(INTERNAL_API_HEADER)).toBe("internal-token");
	});

	it("moves OCR documents to collections through the backend", async () => {
		const { moveOCRDocumentsToCollections } = await import("./ocr");
		const responseBody = {
			moved_ids: ["document-1", "document-2"],
			moved_count: 2,
			collection_ids: ["collection-1"]
		};
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(responseBody), { status: 200 });
		});

		const result = await moveOCRDocumentsToCollections(
			fetchMock,
			["document-1", "document-2"],
			["collection-1"],
			{ userId: "user-1" }
		);

		expect(result).toEqual(responseBody);
		expect(fetchMock).toHaveBeenCalledWith(
			"http://ocr-api.test/api/ocr/documents/collections?user_id=user-1",
			expect.objectContaining({
				method: "PUT",
				body: JSON.stringify({
					ids: ["document-1", "document-2"],
					collection_ids: ["collection-1"]
				})
			})
		);
		expect(requestHeaders(fetchMock.mock.calls[0][1]).get("content-type")).toBe("application/json");
		expect(requestHeaders(fetchMock.mock.calls[0][1]).get(INTERNAL_API_HEADER)).toBe("internal-token");
	});

	it("throws typed OCR API errors from document move backend error responses", async () => {
		const { moveOCRDocumentsToCollections } = await import("./ocr");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify({ error: "collection not found" }), { status: 404 });
		});

		await expect(
			moveOCRDocumentsToCollections(fetchMock, ["document-1"], ["collection-1"], {
				userId: "user-1"
			})
		).rejects.toMatchObject({
			status: 404,
			message: "collection not found"
		});
	});

	it("rejects invalid successful OCR document move responses", async () => {
		const { moveOCRDocumentsToCollections } = await import("./ocr");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify({ moved_ids: [1], moved_count: "1" }), { status: 200 });
		});

		await expect(moveOCRDocumentsToCollections(fetchMock, ["document-1"])).rejects.toMatchObject({
			status: 502,
			message: "Invalid OCR document move response"
		});
	});

	it("throws typed OCR API errors from bulk document delete backend error responses", async () => {
		const { deleteOCRDocuments } = await import("./ocr");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify({ error: "ids is required" }), { status: 400 });
		});

		await expect(deleteOCRDocuments(fetchMock, [], { userId: "user-1" })).rejects.toMatchObject({
			status: 400,
			message: "ids is required"
		});
	});

	it("creates async OCR jobs through the backend with multipart form data", async () => {
		const { createOCRJob } = await import("./ocr");
		const responseBody = { ...validOCRJobResponse(), document_id: null, status: "queued" };
		const file = new File([new Uint8Array([1, 2, 3])], "invoice.pdf", {
			type: "application/pdf"
		});
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(responseBody), { status: 202 });
		});

		const result = await createOCRJob(fetchMock, {
			file,
			schemaId: "schema-1",
			userId: "user-1"
		});

		expect(result).toEqual(responseBody);
		expect(fetchMock).toHaveBeenCalledWith(
			"http://ocr-api.test/api/ocr/jobs",
			expect.objectContaining({ method: "POST" })
		);
		const body = fetchMock.mock.calls[0][1]?.body as FormData;
		expect(body.get("file")).toBe(file);
		expect(body.get("schema_id")).toBe("schema-1");
		expect(body.get("user_id")).toBe("user-1");
		const headers = requestHeaders(fetchMock.mock.calls[0][1]);
		expect(headers.get(INTERNAL_API_HEADER)).toBe("internal-token");
		expect(headers.get("content-type")).toBeNull();
	});

	it("omits schema_id when creating async OCR jobs without a schema", async () => {
		const { createOCRJob } = await import("./ocr");
		const responseBody = { ...validOCRJobResponse(), document_id: null, status: "queued" };
		const file = new File([new Uint8Array([1, 2, 3])], "invoice.pdf", {
			type: "application/pdf"
		});
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(responseBody), { status: 202 });
		});

		const result = await createOCRJob(fetchMock, {
			file,
			userId: "user-1"
		});

		expect(result).toEqual(responseBody);
		const body = fetchMock.mock.calls[0][1]?.body as FormData;
		expect(body.get("file")).toBe(file);
		expect(body.get("schema_id")).toBeNull();
		expect(body.get("user_id")).toBe("user-1");
	});

	it("gets async OCR job status through the backend", async () => {
		const { getOCRJob } = await import("./ocr");
		const responseBody = validOCRJobResponse();
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(responseBody), { status: 200 });
		});

		const result = await getOCRJob(fetchMock, "job-1", { userId: "user-1" });

		expect(result).toEqual(responseBody);
		expect(fetchMock).toHaveBeenCalledWith(
			"http://ocr-api.test/api/ocr/jobs/job-1?user_id=user-1",
			expect.objectContaining({ method: "GET" })
		);
	});

	it("lists async OCR jobs through the backend", async () => {
		const { listOCRJobs } = await import("./ocr");
		const responseBody = validOCRJobListResponse();
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(responseBody), { status: 200 });
		});

		const result = await listOCRJobs(fetchMock, {
			userId: "user-1",
			cursor: "cursor-0",
			size: 50,
			sort: "desc",
			status: "processing"
		});

		expect(result).toEqual(responseBody);
		expect(fetchMock).toHaveBeenCalledWith(
			"http://ocr-api.test/api/ocr/jobs?user_id=user-1&status=processing&cursor=cursor-0&size=50&sort=desc",
			expect.objectContaining({ method: "GET" })
		);
	});

	it("deletes async OCR jobs through the backend", async () => {
		const { deleteOCRJobs } = await import("./ocr");
		const responseBody = validDeleteOCRJobsResponse();
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(responseBody), { status: 200 });
		});

		const result = await deleteOCRJobs(fetchMock, ["job-1"], { userId: "user-1" });

		expect(result).toEqual(responseBody);
		expect(fetchMock).toHaveBeenCalledWith(
			"http://ocr-api.test/api/ocr/jobs?user_id=user-1",
			expect.objectContaining({
				method: "DELETE",
				body: JSON.stringify({ ids: ["job-1"] })
			})
		);
		expect(requestHeaders(fetchMock.mock.calls[0][1]).get("content-type")).toBe("application/json");
		expect(requestHeaders(fetchMock.mock.calls[0][1]).get(INTERNAL_API_HEADER)).toBe("internal-token");
	});

	it("returns async OCR job failure messages when present", async () => {
		const { getOCRJob } = await import("./ocr");
		const responseBody = {
			...validOCRJobResponse(),
			status: "failed",
			document_id: null,
			error_message: "mistral OCR failed with status 503"
		};
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify(responseBody), { status: 200 });
		});

		const result = await getOCRJob(fetchMock, "job-1");

		expect(result.error_message).toBe("mistral OCR failed with status 503");
	});

	it("throws typed OCR API errors from async job backend error responses", async () => {
		const { createOCRJob, isOCRApiError } = await import("./ocr");
		const file = new File([new Uint8Array([1])], "bad.txt", { type: "text/plain" });
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify({ error: "unsupported file type" }), { status: 400 });
		});

		await expect(createOCRJob(fetchMock, { file, schemaId: "schema-1", userId: "user-1" }))
			.rejects.toMatchObject({
				status: 400,
				message: "unsupported file type"
			});

		try {
			await createOCRJob(fetchMock, { file, schemaId: "schema-1", userId: "user-1" });
		} catch (error) {
			expect(isOCRApiError(error)).toBe(true);
		}
	});

	it("rejects invalid successful async OCR job responses", async () => {
		const { getOCRJob } = await import("./ocr");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify({ id: "job-1" }), { status: 200 });
		});

		await expect(getOCRJob(fetchMock, "job-1")).rejects.toMatchObject({
			status: 502,
			message: "Invalid OCR job response"
		});
	});

	it("rejects invalid successful async OCR job list responses", async () => {
		const { listOCRJobs } = await import("./ocr");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify({ jobs: [{ id: "job-1" }], next_cursor: null }), {
				status: 200
			});
		});

		await expect(listOCRJobs(fetchMock)).rejects.toMatchObject({
			status: 502,
			message: "Invalid OCR job list response"
		});
	});

	it("rejects OCR requests before calling Go when the internal API token is missing", async () => {
		const { getOCRDocument } = await import("./ocr");
		delete process.env.SYNCRA_INTERNAL_API_TOKEN;
		const fetchMock = vi.fn();

		await expect(getOCRDocument(fetchMock, "document-1")).rejects.toMatchObject({
			status: 500,
			message: "OCR service is not configured"
		});
		expect(fetchMock).not.toHaveBeenCalled();
	});

	it("rejects invalid successful async OCR job delete responses", async () => {
		const { deleteOCRJobs } = await import("./ocr");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify({ deleted_ids: [1], deleted_count: "1" }), {
				status: 200
			});
		});

		await expect(deleteOCRJobs(fetchMock, ["job-1"])).rejects.toMatchObject({
			status: 502,
			message: "Invalid OCR job delete response"
		});
	});

	it("rejects async OCR job responses with invalid error messages", async () => {
		const { getOCRJob } = await import("./ocr");
		const fetchMock = vi.fn(async (_input: FetchInput, _init?: FetchInit) => {
			return new Response(JSON.stringify({ ...validOCRJobResponse(), error_message: 42 }), {
				status: 200
			});
		});

		await expect(getOCRJob(fetchMock, "job-1")).rejects.toMatchObject({
			status: 502,
			message: "Invalid OCR job response"
		});
	});
});
