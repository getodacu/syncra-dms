<script lang="ts">
	import ArchiveIcon from '@lucide/svelte/icons/archive';
	import GitBranchIcon from '@lucide/svelte/icons/git-branch';
	import PlusIcon from '@lucide/svelte/icons/plus';
	import SaveIcon from '@lucide/svelte/icons/save';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Separator } from '$lib/components/ui/separator';
	import type { FlatOrganizationUnitNode, OrganizationUnitNode } from './tree';

	type OrganizationUnitsActionData = {
		action?: string;
		values?: {
			id?: string;
			parentId?: string | null;
			name?: string;
			code?: string;
			description?: string;
		};
	};

	let {
		selected,
		canManage,
		parentOptions,
		selectedId,
		actionData
	}: {
		selected: OrganizationUnitNode | null;
		canManage: boolean;
		parentOptions: FlatOrganizationUnitNode[];
		selectedId: string;
		actionData: OrganizationUnitsActionData | null | undefined;
	} = $props();

	const fieldClass =
		'h-9 rounded-md border bg-background px-3 text-sm disabled:cursor-not-allowed disabled:opacity-60';
	const textareaClass =
		'min-h-24 rounded-md border bg-background px-3 py-2 text-sm disabled:cursor-not-allowed disabled:opacity-60';
	const selectedActionValues = $derived.by(() => {
		const values = actionData?.values;
		return values?.id === selected?.id ? values : null;
	});
	const updateValues = $derived(actionData?.action === 'update' ? selectedActionValues : null);
	const moveValues = $derived(actionData?.action === 'move' ? selectedActionValues : null);
	const childCreateValues = $derived(
		actionData?.action === 'create' && actionData.values?.parentId === selected?.id
			? actionData.values
			: null
	);

	function confirmArchive(event: SubmitEvent) {
		if (!selected) return;
		if (!confirm(`Archive ${selected.name} and all descendants?`)) {
			event.preventDefault();
		}
	}
</script>

<Card.Root>
	<Card.Header>
		<div class="flex min-w-0 items-start justify-between gap-3">
			<div class="min-w-0">
				<Card.Title class="truncate">{selected ? selected.name : 'No selection'}</Card.Title>
				<Card.Description>
					{#if selected?.code}
						{selected.code}
					{:else}
						Details
					{/if}
				</Card.Description>
			</div>
			<Badge variant={canManage ? 'secondary' : 'outline'}>{canManage ? 'Admin' : 'Read only'}</Badge>
		</div>
	</Card.Header>

	<Card.Content class="grid gap-4">
		{#if selected}
			{#if canManage}
				<form method="POST" action="?/update" class="grid gap-3">
					<input type="hidden" name="id" value={selected.id} />
					<input type="hidden" name="parentId" value={selected.parentId ?? ''} />
					<input type="hidden" name="selectedId" value={selectedId} />
					<label class="grid gap-1.5 text-sm font-medium">
						Name
						<input class={fieldClass} name="name" value={updateValues?.name ?? selected.name} required />
					</label>
					<label class="grid gap-1.5 text-sm font-medium">
						Code
						<input class={fieldClass} name="code" value={updateValues?.code ?? selected.code ?? ''} />
					</label>
					<label class="grid gap-1.5 text-sm font-medium">
						Description
						<textarea
							class={textareaClass}
							name="description"
							value={updateValues?.description ?? selected.description ?? ''}
						></textarea>
					</label>
					<Button type="submit" size="sm" class="w-fit gap-2">
						<SaveIcon class="size-4" />
						Save
					</Button>
				</form>
			{:else}
				<div class="grid gap-3">
					<label class="grid gap-1.5 text-sm font-medium">
						Name
						<input class={fieldClass} value={selected.name} disabled />
					</label>
					<label class="grid gap-1.5 text-sm font-medium">
						Code
						<input class={fieldClass} value={selected.code ?? ''} disabled />
					</label>
					<label class="grid gap-1.5 text-sm font-medium">
						Description
						<textarea class={textareaClass} value={selected.description ?? ''} disabled></textarea>
					</label>
				</div>
			{/if}

			<dl class="grid gap-2 text-xs text-muted-foreground sm:grid-cols-2">
				<div>
					<dt class="font-medium text-foreground">Created</dt>
					<dd>{selected.createdAt}</dd>
				</div>
				<div>
					<dt class="font-medium text-foreground">Updated</dt>
					<dd>{selected.updatedAt}</dd>
				</div>
			</dl>

			{#if canManage}
				<Separator />

				<form method="POST" action="?/create" class="grid gap-3">
					<input type="hidden" name="parentId" value={selected.id} />
					<input type="hidden" name="selectedId" value={selectedId} />
					<div class="grid gap-3 md:grid-cols-[minmax(0,1fr)_8rem]">
						<label class="grid gap-1.5 text-sm font-medium">
							Child name
							<input class={fieldClass} name="name" value={childCreateValues?.name ?? ''} required />
						</label>
						<label class="grid gap-1.5 text-sm font-medium">
							Code
							<input class={fieldClass} name="code" value={childCreateValues?.code ?? ''} />
						</label>
					</div>
					<label class="grid gap-1.5 text-sm font-medium">
						Description
						<input class={fieldClass} name="description" value={childCreateValues?.description ?? ''} />
					</label>
					<Button type="submit" variant="outline" size="sm" class="w-fit gap-2">
						<PlusIcon class="size-4" />
						Create child
					</Button>
				</form>

				<Separator />

				<form method="POST" action="?/move" class="grid gap-3">
					<input type="hidden" name="id" value={selected.id} />
					<input type="hidden" name="selectedId" value={selectedId} />
					<label class="grid gap-1.5 text-sm font-medium">
						Parent
						<select
							class={fieldClass}
							name="parentId"
							value={moveValues ? (moveValues.parentId ?? '') : (selected.parentId ?? '')}
						>
							<option value="">Root</option>
							{#each parentOptions as option (option.id)}
								<option value={option.id}>
									{`${'— '.repeat(option.depth)}${option.name}${option.code ? ` (${option.code})` : ''}`}
								</option>
							{/each}
						</select>
					</label>
					<Button type="submit" variant="outline" size="sm" class="w-fit gap-2">
						<GitBranchIcon class="size-4" />
						Move
					</Button>
				</form>

				<Separator />

				<form method="POST" action="?/archive" onsubmit={confirmArchive}>
					<input type="hidden" name="id" value={selected.id} />
					<input type="hidden" name="selectedId" value={selectedId} />
					<Button type="submit" variant="destructive" size="sm" class="gap-2">
						<ArchiveIcon class="size-4" />
						Archive
					</Button>
				</form>
			{/if}
		{:else}
			<div
				class="flex min-h-32 items-center justify-center rounded-md border border-dashed bg-muted/20 px-4 text-sm text-muted-foreground"
			>
				No unit selected
			</div>
		{/if}
	</Card.Content>
</Card.Root>
