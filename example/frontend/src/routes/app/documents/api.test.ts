import { describe, expect, it, vi } from "vitest";

import {
	deleteOCRDocument,
	deleteOCRDocuments,
	downloadOCRDocuments,
	fetchOCRDocumentPreview,
	fetchOCRDocuments,
	isOCRDocumentClientError,
	isOCRDocumentNotFoundError,
	moveOCRDocuments,
	OCR_DOCUMENT_QUERY_RETRY_LIMIT,
	OCRDocumentClientError,
	parseContentDispositionFilename,
	shouldRetryOCRDocumentsQuery,
	updateOCRDocument,
} from "./api";

function jsonResponse(body: unknown, init?: ResponseInit) {
	return new Response(JSON.stringify(body), {
		headers: { "content-type": "application/json" },
		...init,
	});
}

function fullPreviewResponse() {
	return {
		id: "document-1",
		created_at: "2026-05-27T00:00:00Z",
		updated_at: "2026-05-27T00:01:00Z",
		user_id: "user-1",
		original_filename: "invoice.pdf",
		mime_type: "application/pdf",
		file_size: 1536,
		page_count: 2,
		document_hash: "abcdef",
		schema_id: "schema-1",
		has_inline_schema: false,
		markdown: "# Invoice",
		annotation_json: { total: 10 },
		cached: true,
	};
}

describe("documents api client", () => {
	it("fetches OCR documents through the SvelteKit proxy", async () => {
		const body = {
			documents: [
				{
					id: "document-1",
					created_at: "2026-05-27T00:00:00Z",
					updated_at: "2026-05-27T00:01:00Z",
					user_id: "user-1",
					original_filename: "invoice.pdf",
					mime_type: "application/pdf",
					file_size: 1536,
					page_count: 2,
					document_hash: "abcdef",
					schema_id: "schema-1",
					has_inline_schema: false,
					collections: [{ id: "collection-1", name: "Invoices" }],
				},
			],
			next_cursor: "cursor-1",
		};
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body));

		await expect(
			fetchOCRDocuments(fetchFn, {
				filename: "invoice",
				cursor: "cursor-0",
				size: 50,
				sort: "asc",
			})
		).resolves.toEqual(body);
		expect(fetchFn).toHaveBeenCalledWith(
			"/api/ocr/documents?filename=invoice&cursor=cursor-0&size=50&sort=asc",
			{ method: "GET" }
		);
	});

	it("fetches OCR documents filtered by collection through the SvelteKit proxy", async () => {
		const body = {
			documents: [],
			next_cursor: null,
		};
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body));

		await expect(fetchOCRDocuments(fetchFn, { collectionId: "collection-1" })).resolves.toEqual(
			body
		);
		expect(fetchFn).toHaveBeenCalledWith("/api/ocr/documents?collection=collection-1", {
			method: "GET",
		});
	});

	it("throws backend JSON messages for list errors", async () => {
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse({ error: "invalid cursor" }, { status: 400 }));

		await expect(fetchOCRDocuments(fetchFn, {})).rejects.toThrow("invalid cursor");
	});

	it("preserves backend response status for list errors", async () => {
		const fetchFn = vi.fn(() =>
			Promise.resolve(jsonResponse({ error: "collection not found" }, { status: 404 }))
		);

		await expect(fetchOCRDocuments(fetchFn, { collectionId: "missing-collection" })).rejects.toMatchObject({
			name: "OCRDocumentClientError",
			message: "collection not found",
			status: 404,
		});

		let error: unknown;
		try {
			await fetchOCRDocuments(fetchFn, { collectionId: "missing-collection" });
		} catch (caught) {
			error = caught;
		}

		expect(error).toBeInstanceOf(OCRDocumentClientError);
		expect(isOCRDocumentClientError(error)).toBe(true);
		expect(isOCRDocumentNotFoundError(error)).toBe(true);
	});

	it("does not retry missing collection document queries", () => {
		expect(OCR_DOCUMENT_QUERY_RETRY_LIMIT).toBe(3);
		expect(
			shouldRetryOCRDocumentsQuery(0, new OCRDocumentClientError(404, "collection not found"), {
				collectionId: "missing-collection",
			})
		).toBe(false);
		expect(shouldRetryOCRDocumentsQuery(0, new OCRDocumentClientError(404, "not found"))).toBe(
			true
		);
		expect(
			shouldRetryOCRDocumentsQuery(0, new OCRDocumentClientError(500, "documents unavailable"), {
				collectionId: "collection-1",
			})
		).toBe(true);
		expect(shouldRetryOCRDocumentsQuery(3, new Error("network unavailable"))).toBe(false);
	});

	it("rejects invalid list responses", async () => {
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse({ documents: [{ id: "document-1" }] }));

		await expect(fetchOCRDocuments(fetchFn, {})).rejects.toThrow(
			"Invalid OCR document list response"
		);
	});

	it("accepts list items without optional user and schema ids", async () => {
		const body = {
			documents: [
				{
					id: "document-1",
					created_at: "2026-05-27T00:00:00Z",
					updated_at: "2026-05-27T00:01:00Z",
					original_filename: "invoice.pdf",
					mime_type: "application/pdf",
					file_size: 1536,
					page_count: 2,
					document_hash: "abcdef",
					has_inline_schema: false,
					collections: [],
				},
			],
			next_cursor: null,
		};
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body));

		await expect(fetchOCRDocuments(fetchFn, {})).resolves.toEqual(body);
	});

	it("rejects list items with null schema ids", async () => {
		const body = {
			documents: [
				{
					id: "document-1",
					created_at: "2026-05-27T00:00:00Z",
					updated_at: "2026-05-27T00:01:00Z",
					original_filename: "invoice.pdf",
					mime_type: "application/pdf",
					file_size: 1536,
					page_count: 2,
					document_hash: "abcdef",
					schema_id: null,
					has_inline_schema: false,
					collections: [],
				},
			],
			next_cursor: null,
		};
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body));

		await expect(fetchOCRDocuments(fetchFn, {})).rejects.toThrow(
			"Invalid OCR document list response"
		);
	});

	it("rejects list items with invalid collection summaries", async () => {
		const body = {
			documents: [
				{
					id: "document-1",
					created_at: "2026-05-27T00:00:00Z",
					updated_at: "2026-05-27T00:01:00Z",
					original_filename: "invoice.pdf",
					mime_type: "application/pdf",
					file_size: 1536,
					page_count: 2,
					document_hash: "abcdef",
					has_inline_schema: false,
					collections: [{ id: "collection-1" }],
				},
			],
			next_cursor: null,
		};
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body));

		await expect(fetchOCRDocuments(fetchFn, {})).rejects.toThrow(
			"Invalid OCR document list response"
		);
	});

	it("fetches OCR document previews through the SvelteKit proxy", async () => {
		const body = fullPreviewResponse();
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body));

		await expect(fetchOCRDocumentPreview(fetchFn, "document-1")).resolves.toEqual(body);
		expect(fetchFn).toHaveBeenCalledWith("/api/ocr/document/document-1", { method: "GET" });
	});

	it("throws backend JSON messages for preview errors", async () => {
		const fetchFn = vi
			.fn()
			.mockResolvedValue(jsonResponse({ error: "document not found" }, { status: 404 }));

		await expect(fetchOCRDocumentPreview(fetchFn, "document-1")).rejects.toThrow(
			"document not found"
		);
	});

	it("rejects invalid preview responses", async () => {
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse({ id: "document-1" }));

		await expect(fetchOCRDocumentPreview(fetchFn, "document-1")).rejects.toThrow(
			"Invalid OCR document response"
		);
	});

	it("accepts full preview responses without optional fields", async () => {
		const body = {
			id: "document-1",
			created_at: "2026-05-27T00:00:00Z",
			updated_at: "2026-05-27T00:01:00Z",
			original_filename: "invoice.pdf",
			mime_type: "application/pdf",
			file_size: 1536,
			page_count: 2,
			document_hash: "abcdef",
			has_inline_schema: false,
			markdown: "# Invoice",
			cached: true,
		};
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body));

		await expect(fetchOCRDocumentPreview(fetchFn, "document-1")).resolves.toEqual(body);
	});

	it("rejects preview responses with null schema ids", async () => {
		const body = {
			...fullPreviewResponse(),
			schema_id: null,
		};
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body));

		await expect(fetchOCRDocumentPreview(fetchFn, "document-1")).rejects.toThrow(
			"Invalid OCR document response"
		);
	});

	it("rejects truncated preview responses", async () => {
		const body = {
			id: "document-1",
			original_filename: "invoice.pdf",
			markdown: "# Invoice",
		};
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body));

		await expect(fetchOCRDocumentPreview(fetchFn, "document-1")).rejects.toThrow(
			"Invalid OCR document response"
		);
	});

	it("deletes one OCR document through the SvelteKit proxy", async () => {
		const body = { deleted_ids: ["document-1"], deleted_count: 1 };
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body));

		await expect(deleteOCRDocument(fetchFn, "document-1")).resolves.toEqual(body);
		expect(fetchFn).toHaveBeenCalledWith("/api/ocr/documents/document-1", { method: "DELETE" });
	});

	it("throws backend JSON messages for single delete errors", async () => {
		const fetchFn = vi
			.fn()
			.mockResolvedValue(jsonResponse({ error: "OCR document not found" }, { status: 404 }));

		await expect(deleteOCRDocument(fetchFn, "document-1")).rejects.toThrow(
			"OCR document not found"
		);
	});

	it("downloads OCR documents through the SvelteKit proxy", async () => {
		const blob = new Blob(["# Invoice"], { type: "text/markdown" });
		const fetchFn = vi.fn().mockResolvedValue(
			new Response(blob, {
				headers: {
					"content-disposition":
						"attachment; filename=\"invoice.md\"; filename*=UTF-8''invoice.md",
				},
			})
		);

		const result = await downloadOCRDocuments(fetchFn, ["document-1"], "markdown");

		expect(fetchFn).toHaveBeenCalledWith("/api/ocr/documents/download", {
			method: "POST",
			headers: { "content-type": "application/json" },
			body: JSON.stringify({ ids: ["document-1"], format: "markdown" }),
		});
		expect(result.filename).toBe("invoice.md");
		expect(await result.blob.text()).toBe("# Invoice");
	});

	it("throws backend JSON messages for download errors", async () => {
		const fetchFn = vi
			.fn()
			.mockResolvedValue(jsonResponse({ error: "No schema-backed JSON" }, { status: 400 }));

		await expect(downloadOCRDocuments(fetchFn, ["document-1"], "json")).rejects.toThrow(
			"No schema-backed JSON"
		);
	});

	it("parses encoded content-disposition filenames first", () => {
		expect(
			parseContentDispositionFilename(
				"attachment; filename=\"fallback.zip\"; filename*=UTF-8''syncra-%C8%99.zip"
			)
		).toBe("syncra-ș.zip");
	});

	it("parses quoted content-disposition filenames", () => {
		expect(parseContentDispositionFilename('attachment; filename="invoice.md"')).toBe(
			"invoice.md"
		);
	});

	it("updates one OCR document through the SvelteKit proxy", async () => {
		const body = {
			id: "document-1",
			created_at: "2026-05-27T00:00:00Z",
			updated_at: "2026-05-27T00:02:00Z",
			user_id: "user-1",
			original_filename: "renamed.pdf",
			mime_type: "application/pdf",
			file_size: 1536,
			page_count: 2,
			document_hash: "abcdef",
			schema_id: "schema-1",
			has_inline_schema: false,
		};
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body));

		await expect(updateOCRDocument(fetchFn, "document-1", "renamed.pdf")).resolves.toEqual(body);
		expect(fetchFn).toHaveBeenCalledWith("/api/ocr/documents/document-1", {
			method: "PATCH",
			headers: { "content-type": "application/json" },
			body: JSON.stringify({ original_filename: "renamed.pdf" }),
		});
	});

	it("throws backend JSON messages for update errors", async () => {
		const fetchFn = vi
			.fn()
			.mockResolvedValue(jsonResponse({ error: "original_filename is required" }, { status: 400 }));

		await expect(updateOCRDocument(fetchFn, "document-1", "")).rejects.toThrow(
			"original_filename is required"
		);
	});

	it("rejects invalid update responses", async () => {
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse({ id: "document-1" }));

		await expect(updateOCRDocument(fetchFn, "document-1", "renamed.pdf")).rejects.toThrow(
			"Invalid OCR document update response"
		);
	});

	it("bulk deletes OCR documents through the SvelteKit proxy", async () => {
		const body = { deleted_ids: ["document-1", "document-2"], deleted_count: 2 };
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body));

		await expect(deleteOCRDocuments(fetchFn, ["document-1", "document-2"])).resolves.toEqual(
			body
		);
		expect(fetchFn).toHaveBeenCalledWith("/api/ocr/documents", {
			method: "DELETE",
			headers: { "content-type": "application/json" },
			body: JSON.stringify({ ids: ["document-1", "document-2"] }),
		});
	});

	it("moves OCR documents through the SvelteKit proxy", async () => {
		const body = {
			moved_ids: ["document-1", "document-2"],
			moved_count: 2,
			collection_ids: ["collection-1"],
		};
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body));

		await expect(
			moveOCRDocuments(fetchFn, ["document-1", "document-2"], ["collection-1"])
		).resolves.toEqual(body);
		expect(fetchFn).toHaveBeenCalledWith("/api/ocr/documents/collections", {
			method: "PUT",
			headers: { "content-type": "application/json" },
			body: JSON.stringify({
				ids: ["document-1", "document-2"],
				collection_ids: ["collection-1"],
			}),
		});
	});

	it("throws backend JSON messages for move errors", async () => {
		const fetchFn = vi
			.fn()
			.mockResolvedValue(jsonResponse({ error: "collection not found" }, { status: 404 }));

		await expect(moveOCRDocuments(fetchFn, ["document-1"], ["collection-1"])).rejects.toThrow(
			"collection not found"
		);
	});

	it("rejects invalid move responses", async () => {
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse({ moved_ids: [1], moved_count: 1 }));

		await expect(moveOCRDocuments(fetchFn, ["document-1"], [])).rejects.toThrow(
			"Invalid OCR document move response"
		);
	});

	it("throws backend JSON messages for bulk delete errors", async () => {
		const fetchFn = vi
			.fn()
			.mockResolvedValue(jsonResponse({ error: "ids is required" }, { status: 400 }));

		await expect(deleteOCRDocuments(fetchFn, [])).rejects.toThrow("ids is required");
	});

	it("rejects invalid delete responses", async () => {
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse({ deleted_ids: [1], deleted_count: 1 }));

		await expect(deleteOCRDocument(fetchFn, "document-1")).rejects.toThrow(
			"Invalid OCR document delete response"
		);
	});
});
