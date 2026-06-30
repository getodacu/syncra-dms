<script lang="ts">
	import ReceiptTextIcon from "@lucide/svelte/icons/receipt-text";
	import { createMutation } from "@tanstack/svelte-query";
	import { toast } from "svelte-sonner";

	import { Button } from "$lib/components/ui/button/index.js";
	import { Spinner } from "$lib/components/ui/spinner/index.js";
	import {
		generateAdminBillingOrderInvoice,
		type AdminBillingOrderResponse
	} from "./api";

	let {
		order,
		onGenerated
	}: {
		order: AdminBillingOrderResponse;
		onGenerated?: () => void | Promise<void>;
	} = $props();

	const disabledReason = $derived.by(() => {
		if (order.status !== "paid") return "Only paid orders can be invoiced";
		return null;
	});

	const invoiceMutation = createMutation(() => ({
		mutationFn: () => generateAdminBillingOrderInvoice(fetch, order.id),
		onSuccess: async (invoice) => {
			toast.success(`Invoice ${invoice.invoice_serie}-${invoice.invoice_number} generated.`);
			await onGenerated?.();
		},
		onError: (error) => {
			toast.error(error instanceof Error ? error.message : "Failed to generate invoice.");
		}
	}));

	function generateInvoice() {
		if (order.invoice || disabledReason || invoiceMutation.isPending) return;
		invoiceMutation.mutate();
	}
</script>

{#if !order.invoice && order.status === "paid"}
	<Button
		type="button"
		variant="default"
		size="sm"
		class="h-8 gap-1.5 whitespace-nowrap bg-emerald-600 text-xs text-white hover:bg-emerald-700 focus-visible:border-emerald-600 focus-visible:ring-emerald-500/30 dark:bg-emerald-500 dark:text-emerald-950 dark:hover:bg-emerald-400"
		disabled={Boolean(disabledReason) || invoiceMutation.isPending}
		title={disabledReason ?? "Generate Invoice"}
		onclick={generateInvoice}
	>
		{#if invoiceMutation.isPending}
			<Spinner class="size-3.5" />
		{:else}
			<ReceiptTextIcon class="size-3.5" aria-hidden="true" />
		{/if}
		Generate Invoice
	</Button>
{/if}
