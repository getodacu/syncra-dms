import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";

const pageSource = () => readFileSync(new URL("./+page.svelte", import.meta.url), "utf8");
const typeCellSource = () => readFileSync(new URL("./type-cell.svelte", import.meta.url), "utf8");
const createdDateHeaderSource = () =>
	readFileSync(new URL("./created-date-header.svelte", import.meta.url), "utf8");

describe("credit usage history page i18n", () => {
	it("uses Paraglide messages for credit usage labels, filters, and helper cells", () => {
		const page = pageSource();
		const typeCell = typeCellSource();
		const createdDateHeader = createdDateHeaderSource();

		for (const source of [page, typeCell, createdDateHeader]) {
			expect(source).toContain('import { m } from "$lib/paraglide/messages.js";');
		}

		for (const messageCall of [
			"m.credit_usage_date_range_filter()",
			"m.credit_usage_type_column()",
			"m.credit_usage_related_id_column()",
			"m.credit_usage_all_activity()",
			"m.credit_usage_no_usage_found()",
			"m.credit_usage_showing_one({ count: creditUsageHistory.length })",
			"m.common_rows_per_page()"
		]) {
			expect(page).toContain(messageCall);
		}

		expect(typeCell).toContain("m.credit_usage_type_purchase()");
		expect(typeCell).toContain("m.credit_usage_type_debit()");
		expect(createdDateHeader).toContain("m.credit_usage_sort_created_ascending()");
	});
});
