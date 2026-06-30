<script lang="ts">
	import ReceiptTextIcon from "@lucide/svelte/icons/receipt-text";
	import { createMutation } from "@tanstack/svelte-query";
	import { toast } from "svelte-sonner";

	import { Button } from "$lib/components/ui/button/index.js";
	import { Spinner } from "$lib/components/ui/spinner/index.js";
	import {
		generateAdminBillingInvoicePDF,
		type AdminBillingInvoiceResponse
	} from "./api";

	let {
		invoice,
		onGenerated
	}: {
		invoice: AdminBillingInvoiceResponse;
		onGenerated?: () => void | Promise<void>;
	} = $props();

	const hasPDF = $derived(Boolean(invoice.pdf_path));
	const invoiceLabel = $derived(formatInvoiceLabel(invoice));

	const pdfMutation = createMutation(() => ({
		mutationFn: () => generateAdminBillingInvoicePDF(fetch, invoice.id),
		onSuccess: async () => {
			toast.success(`Invoice ${invoiceLabel} PDF generated.`);
			await onGenerated?.();
		},
		onError: (error) => {
			toast.error(error instanceof Error ? error.message : "Failed to generate invoice PDF.");
		}
	}));

	function formatInvoiceLabel(value: AdminBillingInvoiceResponse) {
		return `${value.invoice_serie}-${String(value.invoice_number).padStart(5, "0")}`;
	}

	function generatePDF() {
		if (pdfMutation.isPending) return;
		pdfMutation.mutate();
	}
</script>

<div class="flex items-center justify-end gap-2">
	<Button
		type="button"
		variant="outline"
		size="sm"
		class="h-8 gap-1.5 whitespace-nowrap text-xs"
		disabled={pdfMutation.isPending}
		title={hasPDF ? `Regenerate ${invoiceLabel} PDF` : `Generate ${invoiceLabel} PDF`}
		onclick={generatePDF}
	>
		{#if pdfMutation.isPending}
			<Spinner class="size-3.5" />
		{:else}
			<ReceiptTextIcon class="size-3.5" aria-hidden="true" />
		{/if}
		{hasPDF ? "Regenerate" : "Generate PDF"}
	</Button>
</div>
