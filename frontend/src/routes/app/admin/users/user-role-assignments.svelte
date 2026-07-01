<script lang="ts">
	import LinkIcon from '@lucide/svelte/icons/link';
	import UnlinkIcon from '@lucide/svelte/icons/unlink';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import type { ScopedRoleAssignmentInput, User } from './api';

	type Props = {
		selectedUser: User | null;
		canAssignUserRoles: boolean;
		canAssignUserGroups: boolean;
		isMutationPending: boolean;
		onAssignRole: (id: string, input: ScopedRoleAssignmentInput) => Promise<void>;
		onRemoveRole: (id: string, assignmentId: string) => Promise<void>;
		onAssignGroup: (id: string, groupId: string) => Promise<void>;
		onRemoveGroup: (id: string, groupId: string) => Promise<void>;
	};

	let {
		selectedUser,
		canAssignUserRoles,
		canAssignUserGroups,
		isMutationPending,
		onAssignRole,
		onRemoveRole,
		onAssignGroup,
		onRemoveGroup
	}: Props = $props();

	let roleId = $state('');
	let roleScopeType = $state<ScopedRoleAssignmentInput['scopeType']>('global');
	let roleOrganizationUnitId = $state('');
	let roleAssignmentId = $state('');
	let groupId = $state('');
	let removeGroupId = $state('');

	async function assignRole(event: SubmitEvent) {
		event.preventDefault();
		if (!selectedUser || !canAssignUserRoles || !roleId.trim()) return;
		await onAssignRole(selectedUser.id, {
			roleId,
			scopeType: roleScopeType,
			organizationUnitId: blankToNull(roleOrganizationUnitId)
		});
		roleId = '';
		roleOrganizationUnitId = '';
	}

	async function removeRole(event: SubmitEvent) {
		event.preventDefault();
		if (!selectedUser || !canAssignUserRoles || !roleAssignmentId.trim()) return;
		await onRemoveRole(selectedUser.id, roleAssignmentId);
		roleAssignmentId = '';
	}

	async function assignGroup(event: SubmitEvent) {
		event.preventDefault();
		if (!selectedUser || !canAssignUserGroups || !groupId.trim()) return;
		await onAssignGroup(selectedUser.id, groupId);
		groupId = '';
	}

	async function removeGroup(event: SubmitEvent) {
		event.preventDefault();
		if (!selectedUser || !canAssignUserGroups || !removeGroupId.trim()) return;
		await onRemoveGroup(selectedUser.id, removeGroupId);
		removeGroupId = '';
	}

	function blankToNull(value: string) {
		const trimmed = value.trim();
		return trimmed ? trimmed : null;
	}
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>Assignments</Card.Title>
		<Card.Description>{selectedUser?.name ?? 'Select a user'}</Card.Description>
	</Card.Header>
	<Card.Content class="grid gap-4">
		<form class="grid gap-2" onsubmit={assignRole}>
			<div class="grid gap-2 sm:grid-cols-[minmax(0,1fr)_9rem]">
				<label class="grid gap-1 text-sm">
					<span class="font-medium">Role ID</span>
					<Input
						bind:value={roleId}
						disabled={!selectedUser || !canAssignUserRoles || isMutationPending}
						placeholder="UUID"
					/>
				</label>
				<label class="grid gap-1 text-sm">
					<span class="font-medium">Scope</span>
					<select
						class="h-9 rounded-md border bg-background px-3 text-sm"
						bind:value={roleScopeType}
						disabled={!selectedUser || !canAssignUserRoles || isMutationPending}
					>
						<option value="global">Global</option>
						<option value="organization_unit">Unit</option>
						<option value="organization_unit_and_children">Unit tree</option>
					</select>
				</label>
			</div>
			<label class="grid gap-1 text-sm">
				<span class="font-medium">Scope unit ID</span>
				<Input
					bind:value={roleOrganizationUnitId}
					disabled={!selectedUser || !canAssignUserRoles || isMutationPending}
					placeholder="UUID"
				/>
			</label>
			<Button
				type="submit"
				variant="outline"
				class="gap-2"
				disabled={!selectedUser || !canAssignUserRoles || isMutationPending}
			>
				<LinkIcon class="size-4" />
				Assign role
			</Button>
		</form>

		<form class="grid gap-2" onsubmit={removeRole}>
			<label class="grid gap-1 text-sm">
				<span class="font-medium">Role assignment ID</span>
				<Input
					bind:value={roleAssignmentId}
					disabled={!selectedUser || !canAssignUserRoles || isMutationPending}
					placeholder="UUID"
				/>
			</label>
			<Button
				type="submit"
				variant="ghost"
				class="gap-2"
				disabled={!selectedUser || !canAssignUserRoles || isMutationPending}
			>
				<UnlinkIcon class="size-4" />
				Remove role
			</Button>
		</form>

		<form class="grid gap-2" onsubmit={assignGroup}>
			<label class="grid gap-1 text-sm">
				<span class="font-medium">Group ID</span>
				<Input
					bind:value={groupId}
					disabled={!selectedUser || !canAssignUserGroups || isMutationPending}
					placeholder="UUID"
				/>
			</label>
			<Button
				type="submit"
				variant="outline"
				class="gap-2"
				disabled={!selectedUser || !canAssignUserGroups || isMutationPending}
			>
				<LinkIcon class="size-4" />
				Assign group
			</Button>
		</form>

		<form class="grid gap-2" onsubmit={removeGroup}>
			<label class="grid gap-1 text-sm">
				<span class="font-medium">Group ID</span>
				<Input
					bind:value={removeGroupId}
					disabled={!selectedUser || !canAssignUserGroups || isMutationPending}
					placeholder="UUID"
				/>
			</label>
			<Button
				type="submit"
				variant="ghost"
				class="gap-2"
				disabled={!selectedUser || !canAssignUserGroups || isMutationPending}
			>
				<UnlinkIcon class="size-4" />
				Remove group
			</Button>
		</form>
	</Card.Content>
</Card.Root>
