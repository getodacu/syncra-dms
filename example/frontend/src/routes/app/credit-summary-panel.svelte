<script lang="ts">
	import ArrowRightIcon from "@lucide/svelte/icons/arrow-right";
	import CoinsIcon from "@lucide/svelte/icons/coins";

	import { Badge } from "$lib/components/ui/badge/index.js";
	import { Button } from "$lib/components/ui/button/index.js";
	import * as Card from "$lib/components/ui/card/index.js";
	import { m } from "$lib/paraglide/messages.js";
	import type { DashboardSummaryResponse } from "./api";
	import { formatInteger } from "./page-state";

	let { summary }: { summary: DashboardSummaryResponse } = $props();
</script>

<Card.Root size="sm" class="min-h-80">
	<Card.Header>
		<div>
			<Card.Title>{m.dashboard_credits_title()}</Card.Title>
			<Card.Description>{m.dashboard_credits_description()}</Card.Description>
		</div>
		<Card.Action>
			{#if summary.credit_summary.low_credit}
				<Badge variant="destructive">{m.dashboard_low_credit()}</Badge>
			{/if}
		</Card.Action>
	</Card.Header>
	<Card.Content class="flex flex-1 flex-col justify-between gap-6">
		<div class="space-y-4">
			<div class="flex items-center gap-3">
				<div class="rounded-md bg-muted p-2 text-muted-foreground">
					<CoinsIcon class="size-4" aria-hidden="true" />
				</div>
				<div>
					<div class="text-2xl font-semibold tabular-nums">
						{formatInteger(summary.credit_summary.available_credits)}
					</div>
					<div class="text-sm text-muted-foreground">{m.dashboard_available_credits()}</div>
				</div>
			</div>
			<div class="rounded-md bg-muted/40 p-3">
				<div class="text-sm font-medium tabular-nums">
					{formatInteger(summary.credit_summary.credits_spent)}
				</div>
				<div class="text-xs text-muted-foreground">{m.dashboard_credits_spent_in_range()}</div>
			</div>
		</div>
		<Button href="/app/billing" variant="outline" class="w-full gap-2">
			{m.dashboard_billing()}
			<ArrowRightIcon class="size-4" aria-hidden="true" />
		</Button>
	</Card.Content>
</Card.Root>
