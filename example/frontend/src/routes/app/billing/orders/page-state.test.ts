import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";

function pageSource() {
	return readFileSync(new URL("./+page.svelte", import.meta.url), "utf8");
}

function invoiceCellSource() {
	return readFileSync(new URL("./invoice-pdf-cell.svelte", import.meta.url), "utf8");
}

function statusCellSource() {
	return readFileSync(new URL("./status-cell.svelte", import.meta.url), "utf8");
}

function orderDateHeaderSource() {
	return readFileSync(new URL("./order-date-header.svelte", import.meta.url), "utf8");
}

function normalizeSource(source: string) {
	return source.replace(/\s+/g, " ");
}

describe("billing orders invoice PDF preview", () => {
	it("wires an invoice column to the PDF preview cell", () => {
		const source = normalizeSource(pageSource());

		expect(source).toContain('import InvoicePDFCell from "./invoice-pdf-cell.svelte";');
		expect(source).toContain("header: m.billing_orders_invoice_column()");
		expect(source).toContain("renderComponent(InvoicePDFCell");
		expect(source).toContain('class="min-w-[1080px]"');
	});

	it("renders only paid invoice PDFs with preview and download controls", () => {
		const source = normalizeSource(invoiceCellSource());

		expect(source).toContain('order.status === "paid"');
		expect(source).toContain("Boolean(order.invoice?.pdf_path)");
		expect(source).toContain("<Dialog.Root bind:open={previewOpen}>");
		expect(source).toContain("<iframe");
		expect(source).toContain("m.billing_orders_invoice_preview_title({ invoice: invoiceLabel })");
		expect(source).toContain("m.billing_orders_download_invoice()");
		expect(source).toContain("buildBillingInvoicePDFPath(order.invoice.id, { download: true })");
		expect(source).toContain("padStart(5, \"0\")");
	});

	it("uses Paraglide messages for billing order labels, filters, and helper cells", () => {
		const page = pageSource();
		const invoiceCell = invoiceCellSource();
		const statusCell = statusCellSource();
		const orderDateHeader = orderDateHeaderSource();

		for (const source of [page, invoiceCell, statusCell, orderDateHeader]) {
			expect(source).toContain('import { m } from "$lib/paraglide/messages.js";');
		}

		for (const messageCall of [
			"m.billing_orders_all_orders()",
			"m.billing_orders_order_date_filter()",
			"m.billing_orders_amount_column()",
			"m.billing_orders_status_column()",
			"m.billing_orders_payment_datetime_column()",
			"m.billing_orders_no_orders_found()",
			"m.billing_orders_showing_one({ count: orders.length })",
			"m.common_rows_per_page()"
		]) {
			expect(page).toContain(messageCall);
		}

		expect(statusCell).toContain("m.billing_order_status_paid()");
		expect(invoiceCell).toContain("m.billing_orders_invoice_pdf_title({ invoice: invoiceLabel })");
		expect(orderDateHeader).toContain("m.billing_orders_sort_order_date_ascending()");
	});
});
