<script lang="ts">
	import ArrowRightIcon from "@lucide/svelte/icons/arrow-right";
	import DatabaseIcon from "@lucide/svelte/icons/database";

	import { Button } from "$lib/components/ui/button/index.js";
	import * as Card from "$lib/components/ui/card/index.js";
	import { m } from "$lib/paraglide/messages.js";
	import type { DashboardSummaryResponse } from "./api";
	import { formatDashboardDate, formatInteger } from "./page-state";

	let { summary }: { summary: DashboardSummaryResponse } = $props();

	function datasetTotalLabel(count: number) {
		return count === 1
			? m.dashboard_total_datasets_one({ count })
			: m.dashboard_total_datasets_other({ count });
	}

	function fieldCountLabel(count: number) {
		return count === 1 ? m.dashboard_fields_one({ count }) : m.dashboard_fields_other({ count });
	}
</script>

<Card.Root size="sm" class="min-h-80">
	<Card.Header>
		<div>
			<Card.Title>{m.dashboard_datasets_title()}</Card.Title>
			<Card.Description>{datasetTotalLabel(summary.dataset_summary.total_count)}</Card.Description>
		</div>
		<Card.Action>
			<Button href="/app/datasets" variant="ghost" size="sm" class="gap-1">
				{m.dashboard_view()}
				<ArrowRightIcon class="size-4" aria-hidden="true" />
			</Button>
		</Card.Action>
	</Card.Header>
	<Card.Content>
		<div class="flex flex-col gap-3">
			{#each summary.dataset_summary.recent as dataset (dataset.id)}
				<a
					href={`/app/datasets/${encodeURIComponent(dataset.id)}`}
					class="group rounded-md border border-transparent p-2 transition-colors hover:border-border hover:bg-muted/40"
				>
					<div class="flex items-start gap-3">
						<div class="rounded-md bg-muted p-2 text-muted-foreground">
							<DatabaseIcon class="size-4" aria-hidden="true" />
						</div>
						<div class="min-w-0 flex-1">
							<div class="truncate text-sm font-medium group-hover:text-primary">{dataset.name}</div>
							<div class="mt-1 flex min-w-0 gap-x-2 text-xs text-muted-foreground">
								<span class="min-w-0 flex-1 truncate">{dataset.schema_name || m.dashboard_no_saved_schema()}</span>
								<span>{fieldCountLabel(dataset.field_count)}</span>
							</div>
							<div class="mt-1 text-xs text-muted-foreground">
								{formatDashboardDate(dataset.created_at)}
							</div>
						</div>
					</div>
				</a>
			{:else}
				<div class="rounded-md border border-dashed p-4 text-sm text-muted-foreground">
					{m.dashboard_no_datasets()}
				</div>
			{/each}
		</div>
	</Card.Content>
</Card.Root>
