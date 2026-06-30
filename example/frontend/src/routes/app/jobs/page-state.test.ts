import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";

function pageSource() {
	return readFileSync(new URL("./+page.svelte", import.meta.url), "utf8");
}

function normalizeSource(source: string) {
	return source.replace(/\s+/g, " ");
}

describe("jobs page polling state", () => {
	it("invalidates the shared credit balance query when a visible job completes", () => {
		const source = normalizeSource(pageSource());

		expect(source).toContain('import { CREDIT_BALANCE_QUERY_KEY } from "$lib/client/billing";');
		expect(source).toContain("let completedJob = false;");
		expect(source).toContain('latest.status === "completed"');
		expect(source).toContain(
			"queryClient.invalidateQueries({ queryKey: CREDIT_BALANCE_QUERY_KEY })"
		);
	});

	it("uses Paraglide messages for jobs page labels and dialogs", () => {
		const source = pageSource();

		expect(source).toContain('import { m } from "$lib/paraglide/messages.js";');
		for (const messageCall of [
			"m.jobs_delete_bulk_title_one({ count: ids.length })",
			"m.jobs_delete_single_title()",
			"m.jobs_status_completed()",
			"m.jobs_inline_schema()",
			"m.jobs_filename_column()",
			"m.jobs_select_all_on_page()",
			"m.jobs_new_job()",
			"m.jobs_no_jobs_found()",
			"m.jobs_showing_jobs_other({ count: jobs.length })",
			"m.jobs_saved_extraction_schema()",
			"m.common_strict()"
		]) {
			expect(source).toContain(messageCall);
		}
	});
});
