import JSZip from "jszip";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { OCRApiError } from "$lib/server/ocr";
import { POST } from "./+server";
import type { RequestEvent } from "./$types";

const { getOCRDocumentMock, OCRApiErrorMock } = vi.hoisted(() => {
	class MockOCRApiError extends Error {
		status: number;

		constructor(status: number, message: string) {
			super(message);
			this.name = "OCRApiError";
			this.status = status;
		}
	}

	return {
		getOCRDocumentMock: vi.fn(),
		OCRApiErrorMock: MockOCRApiError
	};
});

vi.mock("$lib/server/ocr", () => ({
	getOCRDocument: getOCRDocumentMock,
	OCRApiError: OCRApiErrorMock,
	isOCRApiError: (error: unknown) => error instanceof OCRApiErrorMock
}));

function createEvent(body: unknown, user: unknown = { id: "user-1" }) {
	return {
		request: new Request("http://localhost/api/ocr/documents/download", {
			method: "POST",
			headers: { "content-type": "application/json" },
			body: JSON.stringify(body)
		}),
		fetch: vi.fn(),
		locals: { user }
	} as unknown as RequestEvent;
}

async function responseJson(response: Response) {
	return (await response.json()) as unknown;
}

function validDocument(overrides: Record<string, unknown> = {}) {
	return {
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
		markdown: "# Invoice",
		annotation_json: { total: 10 },
		cached: false,
		...overrides
	};
}

describe("OCR document download endpoint", () => {
	beforeEach(() => {
		getOCRDocumentMock.mockReset();
	});

	it("returns 401 for unauthenticated download requests", async () => {
		const response = await POST(createEvent({ ids: ["document-1"], format: "markdown" }, null));

		expect(response.status).toBe(401);
		expect(await responseJson(response)).toEqual({ error: "authentication required" });
		expect(getOCRDocumentMock).not.toHaveBeenCalled();
	});

	it("rejects invalid download request bodies", async () => {
		const response = await POST(createEvent({ ids: [], format: "markdown" }));

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({
			error: "invalid OCR document download request"
		});
		expect(getOCRDocumentMock).not.toHaveBeenCalled();
	});

	it("rejects invalid download formats", async () => {
		const response = await POST(createEvent({ ids: ["document-1"], format: "pdf" }));

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({
			error: "invalid OCR document download request"
		});
		expect(getOCRDocumentMock).not.toHaveBeenCalled();
	});

	it("returns a single Markdown file without zipping", async () => {
		getOCRDocumentMock.mockResolvedValue(validDocument());
		const event = createEvent({ ids: ["document-1"], format: "markdown" });

		const response = await POST(event);

		expect(getOCRDocumentMock).toHaveBeenCalledWith(event.fetch, "document-1", {
			userId: "user-1"
		});
		expect(response.status).toBe(200);
		expect(response.headers.get("content-type")).toBe("text/markdown; charset=utf-8");
		expect(response.headers.get("content-disposition")).toContain('filename="invoice.md"');
		expect(await response.text()).toBe("# Invoice");
	});

	it("returns a single standalone HTML file without zipping", async () => {
		getOCRDocumentMock.mockResolvedValue(validDocument());

		const response = await POST(createEvent({ ids: ["document-1"], format: "html" }));

		expect(response.status).toBe(200);
		expect(response.headers.get("content-type")).toBe("text/html; charset=utf-8");
		expect(response.headers.get("content-disposition")).toContain('filename="invoice.html"');
		const html = await response.text();
		expect(html).toContain("<!doctype html>");
		expect(html).toContain("<h1>Invoice</h1>");
	});

	it("returns a single JSON file without zipping", async () => {
		getOCRDocumentMock.mockResolvedValue(validDocument());

		const response = await POST(createEvent({ ids: ["document-1"], format: "json" }));

		expect(response.status).toBe(200);
		expect(response.headers.get("content-type")).toBe("application/json; charset=utf-8");
		expect(response.headers.get("content-disposition")).toContain('filename="invoice.json"');
		expect(await response.text()).toBe(JSON.stringify({ total: 10 }, null, 2));
	});

	it("returns a server-side zip for multiple selected documents", async () => {
		getOCRDocumentMock
			.mockResolvedValueOnce(validDocument({ id: "document-1", original_filename: "invoice.pdf" }))
			.mockResolvedValueOnce(validDocument({ id: "document-2", original_filename: "invoice.png" }));

		const response = await POST(
			createEvent({ ids: ["document-1", "document-2"], format: "markdown" })
		);

		expect(response.status).toBe(200);
		expect(response.headers.get("content-type")).toBe("application/zip");
		expect(response.headers.get("content-disposition")).toMatch(
			/attachment; filename="syncra-\d{4}-\d{2}-\d{2}T\d{2}-\d{2}-\d{2}Z\.zip"/
		);

		const zip = await JSZip.loadAsync(await response.arrayBuffer());
		expect(Object.keys(zip.files).sort()).toEqual(["invoice-2.md", "invoice.md"]);
	});

	it("skips non-schema documents in mixed JSON zip downloads", async () => {
		getOCRDocumentMock
			.mockResolvedValueOnce(validDocument({ id: "document-1", original_filename: "invoice.pdf" }))
			.mockResolvedValueOnce(
				validDocument({
					id: "document-2",
					original_filename: "notes.pdf",
					schema_id: undefined,
					has_inline_schema: false,
					annotation_json: { skipped: true }
				})
			)
			.mockResolvedValueOnce(
				validDocument({
					id: "document-3",
					original_filename: "empty.pdf",
					annotation_json: undefined
				})
			);

		const response = await POST(
			createEvent({ ids: ["document-1", "document-2", "document-3"], format: "json" })
		);

		expect(response.status).toBe(200);
		const zip = await JSZip.loadAsync(await response.arrayBuffer());
		expect(Object.keys(zip.files)).toEqual(["invoice.json"]);
		await expect(zip.file("invoice.json")?.async("string")).resolves.toBe(
			JSON.stringify({ total: 10 }, null, 2)
		);
	});

	it("rejects JSON downloads when no selected document has schema-backed JSON", async () => {
		getOCRDocumentMock.mockResolvedValue(
			validDocument({
				schema_id: undefined,
				has_inline_schema: false,
				annotation_json: undefined
			})
		);

		const response = await POST(createEvent({ ids: ["document-1"], format: "json" }));

		expect(response.status).toBe(400);
		expect(await responseJson(response)).toEqual({
			error: "No schema-backed JSON is available for the selected documents"
		});
	});

	it("normalizes OCR service errors", async () => {
		getOCRDocumentMock.mockRejectedValue(new OCRApiError(503, "OCR service unavailable"));

		const response = await POST(createEvent({ ids: ["document-1"], format: "markdown" }));

		expect(response.status).toBe(502);
		expect(await responseJson(response)).toEqual({ error: "A server error occurred. Please try again." });
	});
});
