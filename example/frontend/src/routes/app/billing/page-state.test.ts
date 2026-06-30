import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";

function pageSource() {
	return readFileSync(new URL("./+page.svelte", import.meta.url), "utf8");
}

function normalizeSource(source: string) {
	return source.replace(/\s+/g, " ");
}

describe("billing page credit balance state", () => {
	it("reads the shared TanStack credit balance query seeded from layout data", () => {
		const source = normalizeSource(pageSource());

		expect(source).toContain('import { createQuery } from "@tanstack/svelte-query";');
		expect(source).toContain("CREDIT_BALANCE_QUERY_KEY");
		expect(source).toContain("fetchCreditBalance");
		expect(source).toContain("const creditBalanceQuery = createQuery<CreditBalanceResponse, Error>");
		expect(source).toContain("queryKey: CREDIT_BALANCE_QUERY_KEY");
		expect(source).toContain("queryFn: () => fetchCreditBalance(fetch)");
		expect(source).toContain("initialData: data.initialCreditBalance ?? undefined");
		expect(source).toContain("!creditBalanceQuery.data && balanceError");
		expect(source).not.toContain("data.creditBalance)");
		expect(source).not.toContain("data.creditBalanceError");
	});

	it("uses Paraglide messages for billing page labels and checkout feedback", () => {
		const source = pageSource();

		expect(source).toContain('import { m } from "$lib/paraglide/messages.js";');
		for (const messageCall of [
			"m.billing_unavailable()",
			"m.billing_credit_blocks_error()",
			"m.billing_checkout_unavailable()",
			"m.billing_payment_received_title()",
			"m.nav_billing_orders()",
			"m.nav_credit_usage_history()",
			"m.billing_available_balance()",
			"m.billing_purchase_credits()",
			"m.billing_starting_checkout()",
			"m.billing_secure_checkout()"
		]) {
			expect(source).toContain(messageCall);
		}
	});
});
