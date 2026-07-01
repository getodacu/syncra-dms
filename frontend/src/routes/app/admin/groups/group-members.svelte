<script lang="ts">
	import LinkIcon from '@lucide/svelte/icons/link';
	import UnlinkIcon from '@lucide/svelte/icons/unlink';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import type { Group } from './api';

	type Props = {
		selectedGroup: Group | null;
		canManageGroupUsers: boolean;
		isMutationPending: boolean;
		onAddUser: (id: string, userId: string) => Promise<void>;
		onRemoveUser: (id: string, userId: string) => Promise<void>;
	};

	let { selectedGroup, canManageGroupUsers, isMutationPending, onAddUser, onRemoveUser }: Props =
		$props();

	let addUserId = $state('');
	let removeUserId = $state('');

	async function addUser(event: SubmitEvent) {
		event.preventDefault();
		if (!selectedGroup || !canManageGroupUsers || !addUserId.trim()) return;
		await onAddUser(selectedGroup.id, addUserId);
		addUserId = '';
	}

	async function removeUser(event: SubmitEvent) {
		event.preventDefault();
		if (!selectedGroup || !canManageGroupUsers || !removeUserId.trim()) return;
		await onRemoveUser(selectedGroup.id, removeUserId);
		removeUserId = '';
	}
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>Members</Card.Title>
		<Card.Description>{selectedGroup?.name ?? 'Select a group'}</Card.Description>
	</Card.Header>
	<Card.Content class="grid gap-4">
		<form class="grid gap-2" onsubmit={addUser}>
			<label class="grid gap-1 text-sm">
				<span class="font-medium">User ID</span>
				<Input bind:value={addUserId} disabled={!selectedGroup || !canManageGroupUsers || isMutationPending} />
			</label>
			<Button type="submit" variant="outline" class="gap-2" disabled={!selectedGroup || !canManageGroupUsers || isMutationPending}>
				<LinkIcon class="size-4" />
				Add user
			</Button>
		</form>
		<form class="grid gap-2" onsubmit={removeUser}>
			<label class="grid gap-1 text-sm">
				<span class="font-medium">User ID</span>
				<Input bind:value={removeUserId} disabled={!selectedGroup || !canManageGroupUsers || isMutationPending} />
			</label>
			<Button type="submit" variant="ghost" class="gap-2" disabled={!selectedGroup || !canManageGroupUsers || isMutationPending}>
				<UnlinkIcon class="size-4" />
				Remove user
			</Button>
		</form>
	</Card.Content>
</Card.Root>
