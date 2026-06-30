import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";

const siteHeaderSource = () => readFileSync(new URL("./site-header.svelte", import.meta.url), "utf8");
const normalizeSource = (value: string) => value.replace(/\s+/g, " ").trim();

describe("site header behavior", () => {
	it("shows the developer settings title on the developer settings route", () => {
		const source = normalizeSource(siteHeaderSource());

		expect(source).toContain('import { m } from "$lib/paraglide/messages.js";');
		expect(source).toContain('pathname === "/app/developer-settings"');
		expect(source).toContain("m.nav_developer_settings()");
	});

	it("uses Paraglide messages for route titles and header controls", () => {
		const source = normalizeSource(siteHeaderSource());

		for (const messageCall of [
			"m.nav_dashboard()",
			"m.nav_schemas()",
			"m.schemas_library()",
			"m.nav_new_schema()",
			"m.nav_edit_schema()",
			"m.nav_new_job()",
			"m.nav_jobs()",
			"m.nav_billing()",
			"m.nav_billing_orders()",
			"m.nav_credit_usage_history()",
			"m.nav_developer_settings()",
			"m.datasets_page_title()",
			"m.documents_page_title()",
			"m.header_credits_unavailable()",
			"m.header_credits({ count: creditBalanceQuery.data.available_credits.toLocaleString(getLocale()) })",
			"m.header_credit_balance_unavailable({ message: errorMessage })",
			"m.common_toggle_theme()"
		]) {
			expect(source).toContain(messageCall);
		}
		expect(source).toContain('pathname === "/app/schemas/library"');
	});

	it("renders the credit balance link beside the theme toggle", () => {
		const source = normalizeSource(siteHeaderSource());
		const balanceLinkIndex = source.indexOf('href="/app/billing"');
		const themeToggleIndex = source.indexOf("onclick={toggleMode}");

		expect(source).toContain('import { createQuery } from "@tanstack/svelte-query";');
		expect(source).toContain("CREDIT_BALANCE_QUERY_KEY");
		expect(source).toContain("fetchCreditBalance");
		expect(source).toContain("initialCreditBalance: CreditBalanceResponse | null");
		expect(source).toContain("initialCreditBalanceError: string | null");
		expect(source).toContain("const creditBalanceQuery = createQuery<CreditBalanceResponse, Error>");
		expect(source).toContain("queryKey: CREDIT_BALANCE_QUERY_KEY");
		expect(source).toContain("queryFn: () => fetchCreditBalance(fetch)");
		expect(source).toContain("initialData: initialCreditBalance ?? undefined");
		expect(source).toContain('import CoinsIcon from "@lucide/svelte/icons/coins";');
		expect(source).toContain('import { getLocale } from "$lib/paraglide/runtime.js";');
		expect(balanceLinkIndex).toBeGreaterThan(-1);
		expect(themeToggleIndex).toBeGreaterThan(-1);
		expect(balanceLinkIndex).toBeLessThan(themeToggleIndex);
	});
});
