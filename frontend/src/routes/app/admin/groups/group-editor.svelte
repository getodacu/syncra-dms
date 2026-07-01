<script lang="ts">
	import SaveIcon from '@lucide/svelte/icons/save';
	import Trash2Icon from '@lucide/svelte/icons/trash-2';
	import UserRoundCogIcon from '@lucide/svelte/icons/user-round-cog';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import type { CreateGroupInput, Group, UpdateGroupInput } from './api';

	type Props = {
		selectedGroup: Group | null;
		canCreateGroups: boolean;
		canUpdateGroups: boolean;
		canDeleteGroups: boolean;
		isMutationPending: boolean;
		onCreate: (input: CreateGroupInput) => Promise<void>;
		onUpdate: (id: string, input: UpdateGroupInput) => Promise<void>;
		onDelete: (id: string) => Promise<void>;
	};

	let {
		selectedGroup,
		canCreateGroups,
		canUpdateGroups,
		canDeleteGroups,
		isMutationPending,
		onCreate,
		onUpdate,
		onDelete
	}: Props = $props();

	const isEditing = $derived(Boolean(selectedGroup));
	const canSubmit = $derived(isEditing ? canUpdateGroups : canCreateGroups);

	async function submit(event: SubmitEvent) {
		event.preventDefault();
		if (!canSubmit || isMutationPending) return;
		const form = event.currentTarget as HTMLFormElement;
		const data = new FormData(form);
		const input = {
			name: text(data, 'name'),
			code: text(data, 'code'),
			description: nullableText(data, 'description'),
			organizationUnitId: nullableText(data, 'organizationUnitId'),
			isActive: data.get('isActive') === 'on'
		};
		if (selectedGroup) {
			await onUpdate(selectedGroup.id, input);
			return;
		}
		await onCreate(input);
		form.reset();
	}

	function text(data: FormData, key: string) {
		const value = data.get(key);
		return typeof value === 'string' ? value : '';
	}

	function nullableText(data: FormData, key: string) {
		const value = text(data, key).trim();
		return value ? value : null;
	}
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>{selectedGroup ? 'Group details' : 'Create group'}</Card.Title>
		<Card.Description>{selectedGroup?.code ?? 'Assignment group'}</Card.Description>
	</Card.Header>
	<Card.Content>
		<form class="grid gap-3" onsubmit={submit}>
			<label class="grid gap-1 text-sm">
				<span class="font-medium">Name</span>
				<Input name="name" value={selectedGroup?.name ?? ''} disabled={!canSubmit || isMutationPending} required />
			</label>
			<label class="grid gap-1 text-sm">
				<span class="font-medium">Code</span>
				<Input
					name="code"
					value={selectedGroup?.code ?? ''}
					disabled={isMutationPending || isEditing}
					required={!isEditing}
				/>
			</label>
			<label class="grid gap-1 text-sm">
				<span class="font-medium">Organization unit ID</span>
				<Input
					name="organizationUnitId"
					value={selectedGroup?.organizationUnitId ?? ''}
					disabled={!canSubmit || isMutationPending}
					placeholder="UUID"
				/>
			</label>
			<label class="grid gap-1 text-sm">
				<span class="font-medium">Description</span>
				<textarea
					name="description"
					class="min-h-20 rounded-md border bg-background px-3 py-2 text-sm"
					disabled={!canSubmit || isMutationPending}
				>{selectedGroup?.description ?? ''}</textarea>
			</label>
			<label class="flex items-center gap-2 text-sm">
				<input
					name="isActive"
					type="checkbox"
					checked={selectedGroup?.isActive ?? true}
					disabled={!canSubmit || isMutationPending}
				/>
				<span>Active</span>
			</label>
			<div class="flex gap-2">
				<Button type="submit" class="gap-2" disabled={!canSubmit || isMutationPending}>
					{#if selectedGroup}
						<SaveIcon class="size-4" />
						Save
					{:else}
						<UserRoundCogIcon class="size-4" />
						Create
					{/if}
				</Button>
				{#if selectedGroup}
					<Button
						type="button"
						variant="destructive"
						class="gap-2"
						disabled={!canDeleteGroups || isMutationPending}
						onclick={() => onDelete(selectedGroup.id)}
					>
						<Trash2Icon class="size-4" />
						Delete
					</Button>
				{/if}
			</div>
		</form>
	</Card.Content>
</Card.Root>
