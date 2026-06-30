<script lang="ts">
	import Building2Icon from '@lucide/svelte/icons/building-2';
	import { Badge } from '$lib/components/ui/badge';
	import type { FlatOrganizationUnitNode } from './tree';

	let {
		units,
		selectedId,
		onSelect
	}: {
		units: FlatOrganizationUnitNode[];
		selectedId: string | null;
		onSelect: (id: string) => void;
	} = $props();
</script>

<div class="grid gap-1" role="tree" aria-label="Organization units">
	{#each units as unit (unit.id)}
		<button
			type="button"
			role="treeitem"
			class="flex h-9 min-w-0 items-center gap-2 rounded-md px-2 text-left text-sm transition-colors hover:bg-muted data-[selected=true]:bg-secondary data-[selected=true]:text-secondary-foreground"
			style={`padding-left: ${0.5 + unit.depth * 1.25}rem`}
			data-selected={unit.id === selectedId}
			aria-level={unit.depth + 1}
			aria-selected={unit.id === selectedId}
			onclick={() => onSelect(unit.id)}
		>
			<Building2Icon class="size-3.5 shrink-0 text-muted-foreground" />
			<span class="min-w-0 truncate font-medium" title={unit.name}>{unit.name}</span>
			{#if unit.code}
				<Badge variant="outline" class="ms-auto max-w-24 truncate text-muted-foreground">
					{unit.code}
				</Badge>
			{/if}
		</button>
	{:else}
		<div
			class="flex h-16 items-center justify-center rounded-md border border-dashed bg-muted/20 px-3 text-sm text-muted-foreground"
		>
			No active units
		</div>
	{/each}
</div>
