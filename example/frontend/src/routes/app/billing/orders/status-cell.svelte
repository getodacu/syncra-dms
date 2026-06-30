<script lang="ts">
	import { Badge } from "$lib/components/ui/badge/index.js";
	import { m } from "$lib/paraglide/messages.js";
	import type { BillingOrderStatus } from "./api";

	let { status }: { status: BillingOrderStatus } = $props();

	const label = $derived.by(() => {
		if (status === "pending") return m.billing_order_status_pending();
		if (status === "paid") return m.billing_order_status_paid();
		if (status === "failed") return m.billing_order_status_failed();
		if (status === "refunded") return m.billing_order_status_refunded();
		return m.billing_order_status_canceled();
	});
	const className = $derived.by(() => {
		switch (status) {
			case "paid":
				return "border-emerald-500/20 bg-emerald-500/10 text-emerald-700 dark:border-emerald-500/30 dark:bg-emerald-500/15 dark:text-emerald-300";
			case "pending":
				return "border-amber-500/20 bg-amber-500/10 text-amber-700 dark:border-amber-500/30 dark:bg-amber-500/15 dark:text-amber-300";
			case "failed":
				return "border-rose-500/20 bg-rose-500/10 text-rose-700 dark:border-rose-500/30 dark:bg-rose-500/15 dark:text-rose-300";
			case "refunded":
				return "border-sky-500/20 bg-sky-500/10 text-sky-700 dark:border-sky-500/30 dark:bg-sky-500/15 dark:text-sky-300";
			case "canceled":
				return "border-muted-foreground/20 bg-muted text-muted-foreground";
		}
	});
</script>

<Badge variant="outline" class={className}>
	{label}
</Badge>
