<script lang="ts">
	import ArchiveIcon from '@lucide/svelte/icons/archive';
	import GitBranchIcon from '@lucide/svelte/icons/git-branch';
	import PlusIcon from '@lucide/svelte/icons/plus';
	import SaveIcon from '@lucide/svelte/icons/save';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { confirmDelete } from '$lib/components/ui/confirm-delete-dialog/index.js';
	import { Separator } from '$lib/components/ui/separator';
	import type { OrganizationUnitInput } from './api';
	import type { FlatOrganizationUnitNode, OrganizationUnitNode } from './tree';

	let {
		selected,
		canManage,
		parentOptions,
		isPending,
		onCreateChild,
		onUpdate,
		onMove,
		onArchive
	}: {
		selected: OrganizationUnitNode | null;
		canManage: boolean;
		parentOptions: FlatOrganizationUnitNode[];
		isPending: boolean;
		onCreateChild: (input: OrganizationUnitInput) => Promise<void>;
		onUpdate: (id: string, input: OrganizationUnitInput) => Promise<void>;
		onMove: (id: string, parentId: string | null) => Promise<void>;
		onArchive: (id: string) => Promise<void>;
	} = $props();

	const fieldClass =
		'h-9 rounded-md border bg-background px-3 text-sm disabled:cursor-not-allowed disabled:opacity-60';
	const textareaClass =
		'min-h-24 rounded-md border bg-background px-3 py-2 text-sm disabled:cursor-not-allowed disabled:opacity-60';

	let editName = $derived(selected?.name ?? '');
	let editCode = $derived(selected?.code ?? '');
	let editDescription = $derived(selected?.description ?? '');
	let childName = $derived.by(() => {
		selected?.id;
		return '';
	});
	let childCode = $derived.by(() => {
		selected?.id;
		return '';
	});
	let childDescription = $derived.by(() => {
		selected?.id;
		return '';
	});
	let moveParentId = $derived(selected?.parentId ?? '');

	async function submitUpdate(event: SubmitEvent) {
		event.preventDefault();
		if (!selected || isPending) return;

		try {
			await onUpdate(selected.id, {
				parentId: selected.parentId ?? null,
				name: editName,
				code: editCode,
				description: editDescription
			});
		} catch {
		}
	}

	async function submitChild(event: SubmitEvent) {
		event.preventDefault();
		if (!selected || isPending) return;

		try {
			await onCreateChild({
				parentId: selected.id,
				name: childName,
				code: childCode,
				description: childDescription
			});
			childName = '';
			childCode = '';
			childDescription = '';
		} catch {
		}
	}

	async function submitMove(event: SubmitEvent) {
		event.preventDefault();
		if (!selected || isPending) return;

		try {
			await onMove(selected.id, moveParentId || null);
		} catch {
		}
	}

	function confirmArchive() {
		if (!selected || isPending) return;

		const unit = selected;
		confirmDelete({
			title: `Archive ${unit.name}?`,
			description: 'This archives the selected unit and all descendants.',
			confirm: { text: 'Archive' },
			onConfirm: async () => {
				try {
					await onArchive(unit.id);
				} catch {
				}
			}
		});
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
				<form class="grid gap-3" onsubmit={submitUpdate}>
					<label class="grid gap-1.5 text-sm font-medium">
						Name
						<input class={fieldClass} bind:value={editName} disabled={isPending} required />
					</label>
					<label class="grid gap-1.5 text-sm font-medium">
						Code
						<input class={fieldClass} bind:value={editCode} disabled={isPending} />
					</label>
					<label class="grid gap-1.5 text-sm font-medium">
						Description
						<textarea
							class={textareaClass}
							bind:value={editDescription}
							disabled={isPending}
						></textarea>
					</label>
					<Button type="submit" size="sm" class="w-fit gap-2" disabled={isPending}>
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

				<form class="grid gap-3" onsubmit={submitChild}>
					<div class="grid gap-3 md:grid-cols-[minmax(0,1fr)_8rem]">
						<label class="grid gap-1.5 text-sm font-medium">
							Child name
							<input class={fieldClass} bind:value={childName} disabled={isPending} required />
						</label>
						<label class="grid gap-1.5 text-sm font-medium">
							Code
							<input class={fieldClass} bind:value={childCode} disabled={isPending} />
						</label>
					</div>
					<label class="grid gap-1.5 text-sm font-medium">
						Description
						<input class={fieldClass} bind:value={childDescription} disabled={isPending} />
					</label>
					<Button type="submit" variant="outline" size="sm" class="w-fit gap-2" disabled={isPending}>
						<PlusIcon class="size-4" />
						Create child
					</Button>
				</form>

				<Separator />

				<form class="grid gap-3" onsubmit={submitMove}>
					<label class="grid gap-1.5 text-sm font-medium">
						Parent
						<select class={fieldClass} bind:value={moveParentId} disabled={isPending}>
							<option value="">Root</option>
							{#each parentOptions as option (option.id)}
								<option value={option.id}>
									{`${'- '.repeat(option.depth)}${option.name}${option.code ? ` (${option.code})` : ''}`}
								</option>
							{/each}
						</select>
					</label>
					<Button type="submit" variant="outline" size="sm" class="w-fit gap-2" disabled={isPending}>
						<GitBranchIcon class="size-4" />
						Move
					</Button>
				</form>

				<Separator />

				<Button
					type="button"
					variant="destructive"
					size="sm"
					class="w-fit gap-2"
					disabled={isPending}
					onclick={confirmArchive}
				>
					<ArchiveIcon class="size-4" />
					Archive
				</Button>
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
