import { describe, expect, it, vi } from "vitest";

import { deleteOCRJobs, fetchOCRJob, fetchOCRJobs } from "./api";

function jsonResponse(body: unknown, init?: ResponseInit) {
	return new Response(JSON.stringify(body), {
		headers: { "content-type": "application/json" },
		...init,
	});
}

function jobResponse() {
	return {
		id: "job-1",
		created_at: "2026-05-27T00:00:00Z",
		original_filename: "invoice.pdf",
		mime_type: "application/pdf",
		status: "processing",
		file_size: 1536,
		page_count: 2,
		schema_id: "schema-1",
		schema_name: "Invoice",
		has_inline_schema: false,
		document_id: null,
	};
}

describe("jobs api client", () => {
	it("fetches OCR jobs through the SvelteKit proxy", async () => {
		const body = { jobs: [jobResponse()], next_cursor: "cursor-1" };
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body));

		await expect(fetchOCRJobs(fetchFn, { cursor: "cursor-0", size: 50, sort: "desc" })).resolves.toEqual(
			body
		);
		expect(fetchFn).toHaveBeenCalledWith("/api/ocr/jobs?cursor=cursor-0&size=50&sort=desc", {
			method: "GET",
		});
	});

	it("fetches an OCR job detail through the SvelteKit proxy", async () => {
		const body = {
			...jobResponse(),
			status: "queued",
			inline_schema: { type: "object" },
			has_inline_schema: true,
		};
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body));

		await expect(fetchOCRJob(fetchFn, "job-1")).resolves.toEqual(body);
		expect(fetchFn).toHaveBeenCalledWith("/api/ocr/jobs/job-1", { method: "GET" });
	});

	it("deletes OCR jobs through the SvelteKit proxy", async () => {
		const body = { deleted_ids: ["job-1"], deleted_count: 1 };
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse(body));

		await expect(deleteOCRJobs(fetchFn, ["job-1"])).resolves.toEqual(body);
		expect(fetchFn).toHaveBeenCalledWith("/api/ocr/jobs", {
			method: "DELETE",
			headers: { "content-type": "application/json" },
			body: JSON.stringify({ ids: ["job-1"] }),
		});
	});

	it("throws backend JSON messages for list errors", async () => {
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse({ error: "invalid cursor" }, { status: 400 }));

		await expect(fetchOCRJobs(fetchFn, {})).rejects.toThrow("invalid cursor");
	});

	it("rejects invalid job list responses", async () => {
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse({ jobs: [{ id: "job-1" }], next_cursor: null }));

		await expect(fetchOCRJobs(fetchFn, {})).rejects.toThrow("Invalid OCR job list response");
	});

	it("rejects invalid job detail responses", async () => {
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse({ id: "job-1" }));

		await expect(fetchOCRJob(fetchFn, "job-1")).rejects.toThrow("Invalid OCR job response");
	});

	it("rejects invalid job delete responses", async () => {
		const fetchFn = vi.fn().mockResolvedValue(jsonResponse({ deleted_ids: [1], deleted_count: "1" }));

		await expect(deleteOCRJobs(fetchFn, ["job-1"])).rejects.toThrow(
			"Invalid OCR job delete response"
		);
	});
});
