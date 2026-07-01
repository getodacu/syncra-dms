<script lang="ts">
	import LinkIcon from '@lucide/svelte/icons/link';
	import UnlinkIcon from '@lucide/svelte/icons/unlink';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import type { Group, ScopedGroupRoleInput } from './api';

	type Props = {
		selectedGroup: Group | null;
		canAssignGroupRoles: boolean;
		isMutationPending: boolean;
		onAssignRole: (id: string, input: ScopedGroupRoleInput) => Promise<void>;
		onRemoveRole: (id: string, assignmentId: string) => Promise<void>;
	};

	let { selectedGroup, canAssignGroupRoles, isMutationPending, onAssignRole, onRemoveRole }: Props =
		$props();

	let roleId = $state('');
	let scopeType = $state<ScopedGroupRoleInput['scopeType']>('global');
	let organizationUnitId = $state('');
	let assignmentId = $state('');

	async function assignRole(event: SubmitEvent) {
		event.preventDefault();
		if (!selectedGroup || !canAssignGroupRoles || !roleId.trim()) return;
		await onAssignRole(selectedGroup.id, {
			roleId,
			scopeType,
			organizationUnitId: blankToNull(organizationUnitId)
		});
		roleId = '';
		organizationUnitId = '';
	}

	async function removeRole(event: SubmitEvent) {
		event.preventDefault();
		if (!selectedGroup || !canAssignGroupRoles || !assignmentId.trim()) return;
		await onRemoveRole(selectedGroup.id, assignmentId);
		assignmentId = '';
	}

	function blankToNull(value: string) {
		const trimmed = value.trim();
		return trimmed ? trimmed : null;
	}
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>Role assignments</Card.Title>
		<Card.Description>{selectedGroup?.name ?? 'Select a group'}</Card.Description>
	</Card.Header>
	<Card.Content class="grid gap-4">
		<form class="grid gap-2" onsubmit={assignRole}>
			<label class="grid gap-1 text-sm">
				<span class="font-medium">Role ID</span>
				<Input bind:value={roleId} disabled={!selectedGroup || !canAssignGroupRoles || isMutationPending} />
			</label>
			<label class="grid gap-1 text-sm">
				<span class="font-medium">Scope</span>
				<select
					class="h-9 rounded-md border bg-background px-3 text-sm"
					bind:value={scopeType}
					disabled={!selectedGroup || !canAssignGroupRoles || isMutationPending}
				>
					<option value="global">Global</option>
					<option value="organization_unit">Unit</option>
					<option value="organization_unit_and_children">Unit tree</option>
				</select>
			</label>
			<label class="grid gap-1 text-sm">
				<span class="font-medium">Scope unit ID</span>
				<Input
					bind:value={organizationUnitId}
					disabled={!selectedGroup || !canAssignGroupRoles || isMutationPending}
				/>
			</label>
			<Button type="submit" variant="outline" class="gap-2" disabled={!selectedGroup || !canAssignGroupRoles || isMutationPending}>
				<LinkIcon class="size-4" />
				Assign role
			</Button>
		</form>
		<form class="grid gap-2" onsubmit={removeRole}>
			<label class="grid gap-1 text-sm">
				<span class="font-medium">Assignment ID</span>
				<Input bind:value={assignmentId} disabled={!selectedGroup || !canAssignGroupRoles || isMutationPending} />
			</label>
			<Button type="submit" variant="ghost" class="gap-2" disabled={!selectedGroup || !canAssignGroupRoles || isMutationPending}>
				<UnlinkIcon class="size-4" />
				Remove role
			</Button>
		</form>
	</Card.Content>
</Card.Root>
