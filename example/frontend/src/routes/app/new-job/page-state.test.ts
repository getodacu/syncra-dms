import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";

function pageSource() {
	return readFileSync(new URL("./+page.svelte", import.meta.url), "utf8");
}

function normalizeSource(source: string) {
	return source.replace(/\s+/g, " ");
}

describe("new job polling state", () => {
	it("invalidates the shared credit balance query only after a polled job completes", () => {
		const source = normalizeSource(pageSource());

		expect(source).toContain(
			'import { createQuery, useQueryClient } from "@tanstack/svelte-query";'
		);
		expect(source).toContain('import { CREDIT_BALANCE_QUERY_KEY } from "$lib/client/billing";');
		expect(source).toContain("const queryClient = useQueryClient();");
		expect(source).toContain("let completedJob = false;");
		expect(source).toContain('patch.jobStatus === "completed"');
		expect(source).toContain(
			"queryClient.invalidateQueries({ queryKey: CREDIT_BALANCE_QUERY_KEY })"
		);
	});

	it("uses Paraglide messages for new-job workflow labels and feedback", () => {
		const source = pageSource();

		expect(source).toContain('import { m } from "$lib/paraglide/messages.js";');
		for (const messageCall of [
			"m.new_job_select_schema()",
			"m.new_job_upload_documents()",
			"m.new_job_run_monitor()",
			"m.new_job_search_schemas()",
			"m.new_job_no_schema_ocr_only()",
			"m.new_job_ocr_only_mode_active()",
			"m.new_job_pending_upload_queue({ count: unsubmittedRows.length })",
			"m.new_job_active_batch_status()",
			"m.new_job_progress({ progress: overallProgress })",
			"m.new_job_run_extraction_one({ count: unsubmittedRows.length })",
			"m.new_job_preview_unavailable()",
			"m.new_job_queueing_documents()"
		]) {
			expect(source).toContain(messageCall);
		}
	});
});
