import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";

function pageSource() {
	return readFileSync(new URL("./+page.svelte", import.meta.url), "utf8");
}

function actionsCellSource() {
	return readFileSync(new URL("./invoice-pdf-actions-cell.svelte", import.meta.url), "utf8");
}

function normalizeSource(source: string) {
	return source.replace(/\s+/g, " ");
}

describe("admin invoice page PDF actions", () => {
	it("adds an applied search filter for invoice client and number searches", () => {
		const source = normalizeSource(pageSource());

		expect(source).toContain('import { Input } from "$lib/components/ui/input/index.js";');
		expect(source).toContain('let pendingSearch = $state("");');
		expect(source).toContain('let search = $state("");');
		expect(source).toContain("search = pendingSearch.trim();");
		expect(source).toContain('placeholder="Search client or invoice number"');
		expect(source).toContain("search,");
		expect(source).toContain('pendingSearch = ""; search = "";');
	});

	it("renders invoice number first and merges billing name with email in the client column", () => {
		const source = normalizeSource(pageSource());
		const invoiceColumnIndex = source.indexOf('id: "invoice_number", header: "Invoice"');
		const clientColumnIndex = source.indexOf('id: "client", header: "Client"');

		expect(invoiceColumnIndex).toBeGreaterThan(-1);
		expect(clientColumnIndex).toBeGreaterThan(invoiceColumnIndex);
		expect(source).not.toContain('header: "Billing email"');
		expect(source).toContain('cell.column.id === "client"');
		expect(source).toContain('{row.original.billing_name || "Unnamed client"}');
		expect(source).toContain("{row.original.billing_email}");
	});

	it("opens invoice PDF previews from the invoice number column", () => {
		const source = normalizeSource(pageSource());

		expect(source).toContain("let previewInvoice = $state<AdminBillingInvoiceResponse | null>(null);");
		expect(source).toContain("function openInvoicePreview(invoice: AdminBillingInvoiceResponse)");
		expect(source).toContain('cell.column.id === "invoice_number"');
		expect(source).toContain("onclick={() => openInvoicePreview(row.original)}");
		expect(source).toContain("<Dialog.Root bind:open={() => Boolean(previewInvoice), setPreviewOpen}>");
		expect(source).toContain("<iframe");
		expect(source).toContain("buildAdminBillingInvoicePDFPath(previewInvoice.id)");
		expect(source).toContain("buildAdminBillingInvoicePDFPath(previewInvoice.id, { download: true })");
		expect(source).toContain("download={previewFilename}");
	});

	it("wires an actions column to the invoice PDF action cell", () => {
		const source = normalizeSource(pageSource());

		expect(source).toContain('import InvoicePDFActionsCell from "./invoice-pdf-actions-cell.svelte";');
		expect(source).toContain('id: "actions"');
		expect(source).toContain("renderComponent(InvoicePDFActionsCell");
		expect(source).toContain("await billingInvoicesQuery.refetch();");
	});

	it("renders generate and regenerate controls without an actions preview button", () => {
		const source = normalizeSource(actionsCellSource());

		expect(source).toContain("Generate PDF");
		expect(source).toContain("Regenerate");
		expect(source).not.toContain("Preview");
		expect(source).not.toContain("<Dialog.Root");
		expect(source).not.toContain("<iframe");
		expect(source).not.toContain("EyeIcon");
	});

	it("uses padded invoice labels in action feedback", () => {
		const source = normalizeSource(actionsCellSource());

		expect(source).toContain("formatInvoiceLabel");
		expect(source).toContain("padStart(5, \"0\")");
		expect(source).toContain("Invoice ${invoiceLabel} PDF generated.");
	});
});
