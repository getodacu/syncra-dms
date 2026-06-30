import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";

const pageSource = () => readFileSync(new URL("./+page.svelte", import.meta.url), "utf8");
const invoiceActionCellSource = () =>
	readFileSync(new URL("./invoice-action-cell.svelte", import.meta.url), "utf8");

describe("admin billing orders page state", () => {
	it("renders payment ids and invoice numbers in dedicated table columns", () => {
		const source = pageSource();

		expect(source).toContain('header: "Payment ID"');
		expect(source).toContain('row.original.provider_payment_intent_id ?? ""');
		expect(source).toContain('header: "Invoice number"');
		expect(source).toContain("row.original.invoice");
		expect(source).toContain('`${invoice.invoice_serie}-${String(invoice.invoice_number).padStart(5, "0")}`');
		expect(source).toContain('class="min-w-[1540px]"');
	});

	it("keeps generated invoice metadata out of the action cell", () => {
		const source = invoiceActionCellSource();

		expect(source).not.toContain("{order.invoice.invoice_serie}-{order.invoice.invoice_number}");
		expect(source).not.toContain("Invoice already generated");
		expect(source).toContain('{#if !order.invoice && order.status === "paid"}');
		expect(source).toContain("Generate Invoice");
		expect(source).toContain("bg-emerald-600");
	});

	it("adds a paid no-invoice shortcut that forces paid orders without changing the status dropdown", () => {
		const source = pageSource();

		expect(source).toContain("let paidNoInvoiceFilter = $state(false);");
		expect(source).toContain('const effectiveStatusFilter = $derived(paidNoInvoiceFilter ? "paid" : statusFilter);');
		expect(source).toContain("withoutInvoice: paidNoInvoiceFilter");
		expect(source).toContain("Paid/No Invoice");
		expect(source).toContain("aria-pressed={paidNoInvoiceFilter}");
	});
});
