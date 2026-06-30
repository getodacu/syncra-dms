<script lang="ts">
	import ArrowRightIcon from "@lucide/svelte/icons/arrow-right";
	import BarChart3Icon from "@lucide/svelte/icons/bar-chart-3";

	import { Button } from "$lib/components/ui/button/index.js";
	import * as Card from "$lib/components/ui/card/index.js";
	import { m } from "$lib/paraglide/messages.js";
	import type { DashboardSummaryResponse } from "./api";
	import { formatInteger } from "./page-state";

	let { summary }: { summary: DashboardSummaryResponse } = $props();

	function documentsProcessedLabel(count: number) {
		return count === 1
			? m.dashboard_documents_processed_one({ count })
			: m.dashboard_documents_processed_other({ count });
	}
</script>

<Card.Root size="sm" class="min-h-80">
	<Card.Header>
		<div>
			<Card.Title>{m.dashboard_schema_throughput_title()}</Card.Title>
			<Card.Description>{m.dashboard_schema_throughput_description()}</Card.Description>
		</div>
		<Card.Action>
			<Button href="/app/documents" variant="ghost" size="sm" class="gap-1">
				{m.dashboard_view()}
				<ArrowRightIcon class="size-4" aria-hidden="true" />
			</Button>
		</Card.Action>
	</Card.Header>
	<Card.Content>
		<div class="flex flex-col gap-3">
			{#each summary.schema_throughput as item (item.schema_id ?? item.schema_name)}
				<svelte:element
					this={item.schema_id ? "a" : "div"}
					href={item.schema_id ? `/app/documents?schema=${encodeURIComponent(item.schema_id)}` : undefined}
					class={item.schema_id
						? "group rounded-md border border-transparent p-2 transition-colors hover:border-border hover:bg-muted/40"
						: "rounded-md border border-transparent p-2"}
				>
					<div class="flex items-center gap-3">
						<div class="rounded-md bg-muted p-2 text-muted-foreground">
							<BarChart3Icon class="size-4" aria-hidden="true" />
						</div>
						<div class="min-w-0 flex-1">
							<div class={item.schema_id ? "truncate text-sm font-medium group-hover:text-primary" : "truncate text-sm font-medium"}>
								{item.schema_name}
							</div>
							<div class="mt-1 text-xs text-muted-foreground">
								{documentsProcessedLabel(item.documents_processed)}
							</div>
						</div>
					</div>
				</svelte:element>
			{:else}
				<div class="rounded-md border border-dashed p-4 text-sm text-muted-foreground">
					{m.dashboard_no_schema_throughput()}
				</div>
			{/each}
		</div>
	</Card.Content>
</Card.Root>
