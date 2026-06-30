<script lang="ts">
	import DownloadIcon from "@lucide/svelte/icons/download";
	import EyeIcon from "@lucide/svelte/icons/eye";

	import { Button } from "$lib/components/ui/button/index.js";
	import * as Dialog from "$lib/components/ui/dialog/index.js";
	import { m } from "$lib/paraglide/messages.js";
	import { buildBillingInvoicePDFPath, type BillingOrderResponse } from "./api";

	let { order }: { order: BillingOrderResponse } = $props();

	let previewOpen = $state(false);

	const hasPDF = $derived(order.status === "paid" && Boolean(order.invoice?.pdf_path));
	const invoiceLabel = $derived(order.invoice ? formatInvoiceLabel(order.invoice) : "");
	const invoiceFilename = $derived(order.invoice ? formatInvoiceFilename(order.invoice) : "");
	const previewURL = $derived(order.invoice ? buildBillingInvoicePDFPath(order.invoice.id) : "");
	const downloadURL = $derived(
		order.invoice ? buildBillingInvoicePDFPath(order.invoice.id, { download: true }) : ""
	);

	function formatInvoiceLabel(value: NonNullable<BillingOrderResponse["invoice"]>) {
		return `${value.invoice_serie}-${String(value.invoice_number).padStart(5, "0")}`;
	}

	function formatInvoiceFilename(value: NonNullable<BillingOrderResponse["invoice"]>) {
		return `${value.invoice_serie}_${String(value.invoice_number).padStart(5, "0")}_${formatInvoiceDateStamp(value.invoice_date)}.pdf`;
	}

	function formatInvoiceDateStamp(invoiceDate: string) {
		const match = invoiceDate.match(/^(\d{4})-(\d{2})-(\d{2})/);
		if (!match) return invoiceDate.replace(/\D/g, "").slice(2, 8);

		const [, year, month, day] = match;
		return `${year.slice(2)}${month}${day}`;
	}

	function openPreview() {
		if (!hasPDF) return;
		previewOpen = true;
	}
</script>

{#if hasPDF && order.invoice}
	<Button
		type="button"
		variant="link"
		size="sm"
		class="h-8 justify-start gap-1.5 px-0 text-xs"
		title={m.billing_orders_invoice_pdf_title({ invoice: invoiceLabel })}
		onclick={openPreview}
	>
		<EyeIcon class="size-3.5" aria-hidden="true" />
		<span>{invoiceLabel}</span>
	</Button>
{:else}
	<span class="text-sm text-muted-foreground">-</span>
{/if}

<Dialog.Root bind:open={previewOpen}>
	<Dialog.Content
		class="flex h-[min(90vh,56rem)] w-full flex-col overflow-hidden p-0 sm:max-w-[min(94vw,72rem)]"
	>
		<Dialog.Header class="border-b px-5 py-4 pr-14">
			<div class="flex min-w-0 items-center justify-between gap-3">
				<div class="min-w-0">
					<Dialog.Title class="truncate text-lg">{m.billing_orders_invoice_preview_title({ invoice: invoiceLabel })}</Dialog.Title>
					<Dialog.Description>{m.billing_orders_invoice_preview_description()}</Dialog.Description>
				</div>
				<Button
					href={downloadURL}
					download={invoiceFilename}
					variant="outline"
					size="sm"
					class="h-8 gap-1.5 text-xs"
				>
					<DownloadIcon class="size-3.5" aria-hidden="true" />
					{m.billing_orders_download_invoice()}
				</Button>
			</div>
		</Dialog.Header>

		<div class="min-h-0 flex-1 bg-muted/30">
			<iframe
				title={m.billing_orders_invoice_iframe_title({ invoice: invoiceLabel })}
				src={previewURL}
				class="h-full w-full border-0 bg-background"
			></iframe>
		</div>
	</Dialog.Content>
</Dialog.Root>
