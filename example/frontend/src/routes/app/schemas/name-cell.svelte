<script lang="ts">
	import BracesIcon from "@lucide/svelte/icons/braces";
	import { m } from "$lib/paraglide/messages.js";
	import type { SchemaListItemResponse } from "./api";

	let { schema }: { schema: SchemaListItemResponse } = $props();

	const description = $derived(schema.description.trim());
</script>

<div class="flex min-w-0 items-start gap-2.5" title={schema.name}>
	<BracesIcon class="mt-0.5 size-4 shrink-0 text-muted-foreground" aria-hidden="true" />
	<div class="min-w-0 space-y-0.5">
		<div class="truncate text-sm font-medium text-foreground">{schema.name}</div>
		{#if description}
			<div class="truncate text-xs text-muted-foreground" title={description}>{description}</div>
		{:else}
			<div class="text-xs text-muted-foreground">{m.schemas_no_description()}</div>
		{/if}
	</div>
</div>
