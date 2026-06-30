import JSZip from "jszip";
import { describe, expect, it } from "vitest";

import type { OCRDocumentResponse } from "$lib/server/ocr";
import {
	buildDownloadFiles,
	buildZip,
	contentDispositionAttachment,
	renderStandaloneHTML,
	sanitizeDownloadBasename,
	stripFileExtension,
	zipDownloadFilename
} from "./export-utils";

function document(overrides: Partial<OCRDocumentResponse> = {}): OCRDocumentResponse {
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
		...overrides
	};
}

describe("document download export helpers", () => {
	it("strips only the final file extension from document filenames", () => {
		expect(stripFileExtension("invoice.final.pdf")).toBe("invoice.final");
		expect(stripFileExtension(".env")).toBe(".env");
	});

	it("sanitizes download basenames for direct files and zip entries", () => {
		expect(sanitizeDownloadBasename('../unsafe:"invoice"?.pdf')).toBe("unsafe-invoice-");
		expect(sanitizeDownloadBasename("   .pdf")).toBe("pdf");
		expect(sanitizeDownloadBasename("///")).toBe("document");
	});

	it("renders standalone HTML from markdown", () => {
		const html = renderStandaloneHTML(document({ original_filename: "invoice<script>.pdf" }));

		expect(html).toContain("<!doctype html>");
		expect(html).toContain('<meta charset="utf-8">');
		expect(html).toContain("<title>invoice-script-</title>");
		expect(html).toContain("<h1>Invoice</h1>");
	});

	it("creates unique filenames for duplicate document basenames", () => {
		const files = buildDownloadFiles(
			[
				document({ id: "document-1", original_filename: "invoice.pdf" }),
				document({ id: "document-2", original_filename: "invoice.png" })
			],
			"markdown"
		);

		expect(files.map((file) => file.filename)).toEqual(["invoice.md", "invoice-2.md"]);
	});

	it("skips non-schema documents for JSON exports", () => {
		const files = buildDownloadFiles(
			[
				document({ id: "document-1", original_filename: "invoice.pdf" }),
				document({
					id: "document-2",
					original_filename: "notes.pdf",
					schema_id: undefined,
					has_inline_schema: false,
					annotation_json: { skipped: true }
				}),
				document({
					id: "document-3",
					original_filename: "empty.pdf",
					annotation_json: undefined
				})
			],
			"json"
		);

		expect(files).toHaveLength(1);
		expect(files[0]).toMatchObject({
			filename: "invoice.json",
			content: JSON.stringify({ total: 10 }, null, 2)
		});
	});

	it("generates server-side zip bytes with the expected entries", async () => {
		const files = buildDownloadFiles(
			[
				document({ id: "document-1", original_filename: "invoice.pdf" }),
				document({ id: "document-2", original_filename: "receipt.pdf", markdown: "# Receipt" })
			],
			"markdown"
		);

		const zip = await JSZip.loadAsync(await buildZip(files));

		expect(Object.keys(zip.files).sort()).toEqual(["invoice.md", "receipt.md"]);
		await expect(zip.file("receipt.md")?.async("string")).resolves.toBe("# Receipt");
	});

	it("formats zip filenames with the syncra prefix and UTC timestamp", () => {
		expect(zipDownloadFilename(new Date("2026-06-06T15:04:05.123Z"))).toBe(
			"syncra-2026-06-06T15-04-05Z.zip"
		);
	});

	it("builds attachment headers with ASCII fallback and encoded filename", () => {
		expect(contentDispositionAttachment("factura-ș.pdf")).toBe(
			'attachment; filename="factura-_.pdf"; filename*=UTF-8\'\'factura-%C8%99.pdf'
		);
	});
});
