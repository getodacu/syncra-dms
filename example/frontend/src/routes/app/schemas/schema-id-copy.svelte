<script lang="ts">
	import CopyIcon from "@lucide/svelte/icons/copy";
	import { toast } from "svelte-sonner";

	import { Button } from "$lib/components/ui/button/index.js";
	import * as Tooltip from "$lib/components/ui/tooltip/index.js";
	import { m } from "$lib/paraglide/messages.js";
	import { cn } from "$lib/utils.js";

	let {
		schemaId,
		showLabel = false,
		compact = false,
	}: {
		schemaId: string;
		showLabel?: boolean;
		compact?: boolean;
	} = $props();

	async function copySchemaId() {
		try {
			await navigator.clipboard.writeText(schemaId);
			toast.success(m.schemas_copy_id_success());
		} catch {
			toast.error(m.schemas_copy_id_error());
		}
	}
</script>

<div class={cn("flex min-w-0 items-center gap-1.5", showLabel && "flex-wrap")}>
	{#if showLabel}
		<span class="text-[10px] font-bold uppercase tracking-wider text-muted-foreground">
			{m.schemas_id_label()}
		</span>
	{/if}
	<code
		class={cn(
			"block min-w-0 truncate rounded-md border bg-muted/40 px-2 py-1 font-mono text-xs text-muted-foreground",
			compact ? "max-w-[220px]" : "max-w-full text-foreground"
		)}
		title={schemaId}
	>
		{schemaId}
	</code>
	<Tooltip.Root>
		<Tooltip.Trigger>
			{#snippet child({ props })}
				<span {...props} class="inline-flex shrink-0">
					<Button
						type="button"
						variant="ghost"
						size="icon-xs"
						class="text-muted-foreground hover:text-foreground"
						title={m.schemas_copy_id()}
						aria-label={m.schemas_copy_id_aria({ id: schemaId })}
						onclick={() => void copySchemaId()}
					>
						<CopyIcon class="size-3.5" aria-hidden="true" />
					</Button>
				</span>
			{/snippet}
		</Tooltip.Trigger>
		<Tooltip.Content>{m.schemas_copy_id()}</Tooltip.Content>
	</Tooltip.Root>
</div>
