<script lang="ts">
	import ArrowRightIcon from "@lucide/svelte/icons/arrow-right";
	import FileTextIcon from "@lucide/svelte/icons/file-text";

	import { Button } from "$lib/components/ui/button/index.js";
	import * as Card from "$lib/components/ui/card/index.js";
	import { m } from "$lib/paraglide/messages.js";
	import type { DashboardSummaryResponse } from "./api";
	import { formatDashboardDate, formatInteger } from "./page-state";

	let { summary }: { summary: DashboardSummaryResponse } = $props();

	function pageCountLabel(count: number) {
		return count === 1 ? m.dashboard_pages_one({ count }) : m.dashboard_pages_other({ count });
	}
</script>

<Card.Root size="sm" class="min-h-80">
	<Card.Header>
		<div>
			<Card.Title>{m.dashboard_recent_documents_title()}</Card.Title>
			<Card.Description>{m.dashboard_recent_documents_description()}</Card.Description>
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
			{#each summary.recent_documents as document (document.id)}
				<a
					href="/app/documents"
					class="group rounded-md border border-transparent p-2 transition-colors hover:border-border hover:bg-muted/40"
				>
					<div class="flex items-start gap-3">
						<div class="rounded-md bg-muted p-2 text-muted-foreground">
							<FileTextIcon class="size-4" aria-hidden="true" />
						</div>
						<div class="min-w-0 flex-1">
							<div class="truncate text-sm font-medium group-hover:text-primary">
								{document.original_filename}
							</div>
							<div class="mt-1 flex min-w-0 gap-x-2 text-xs text-muted-foreground">
								<span class="min-w-0 flex-1 truncate">{document.schema_name ?? m.dashboard_no_saved_schema()}</span>
								<span>{pageCountLabel(document.page_count)}</span>
							</div>
							<div class="mt-1 text-xs text-muted-foreground">
								{formatDashboardDate(document.created_at)}
							</div>
						</div>
					</div>
				</a>
			{:else}
				<div class="rounded-md border border-dashed p-4 text-sm text-muted-foreground">
					{m.dashboard_no_completed_documents()}
				</div>
			{/each}
		</div>
	</Card.Content>
</Card.Root>
