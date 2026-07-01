<script lang="ts">
	import SaveIcon from '@lucide/svelte/icons/save';
	import UserPlusIcon from '@lucide/svelte/icons/user-plus';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import type { CreateUserInput, UpdateUserInput, User } from './api';

	type Props = {
		selectedUser: User | null;
		canCreateUsers: boolean;
		canUpdateUsers: boolean;
		canAssignUserUnits: boolean;
		isMutationPending: boolean;
		onCreate: (input: CreateUserInput) => Promise<void>;
		onUpdate: (id: string, input: UpdateUserInput) => Promise<void>;
		onSetPrimaryOrganizationUnit: (id: string, organizationUnitId: string | null) => Promise<void>;
	};

	let {
		selectedUser,
		canCreateUsers,
		canUpdateUsers,
		canAssignUserUnits,
		isMutationPending,
		onCreate,
		onUpdate,
		onSetPrimaryOrganizationUnit
	}: Props = $props();

	const isEditing = $derived(Boolean(selectedUser));
	const canSubmit = $derived(
		isEditing ? canUpdateUsers || canAssignUserUnits : canCreateUsers
	);

	async function submit(event: SubmitEvent) {
		event.preventDefault();
		if (!canSubmit || isMutationPending) return;

		const form = event.currentTarget as HTMLFormElement;
		const data = new FormData(form);
		const name = formText(data, 'name');
		const email = formText(data, 'email');
		const status = formStatus(data, 'status');
		const primaryOrganizationUnitId = blankToNull(formText(data, 'primaryOrganizationUnitId'));
		const managerUserId = blankToNull(formText(data, 'managerUserId'));
		const jobTitle = blankToNull(formText(data, 'jobTitle'));
		const phone = blankToNull(formText(data, 'phone'));

		if (selectedUser) {
			if (canUpdateUsers) {
				await onUpdate(selectedUser.id, {
					name,
					managerUserId,
					jobTitle,
					phone
				});
			}
			if (
				canAssignUserUnits &&
				primaryOrganizationUnitId !== (selectedUser.primaryOrganizationUnitId ?? null)
			) {
				await onSetPrimaryOrganizationUnit(selectedUser.id, primaryOrganizationUnitId);
			}
			return;
		}

		await onCreate({
			name,
			email,
			status,
			primaryOrganizationUnitId,
			managerUserId,
			jobTitle,
			phone
		});
		form.reset();
	}

	function formText(data: FormData, key: string) {
		const value = data.get(key);
		return typeof value === 'string' ? value : '';
	}

	function formStatus(data: FormData, key: string): User['status'] {
		const value = data.get(key);
		if (
			value === 'invited' ||
			value === 'active' ||
			value === 'inactive' ||
			value === 'suspended' ||
			value === 'deleted'
		) {
			return value;
		}
		return 'invited';
	}

	function blankToNull(value: string) {
		const trimmed = value.trim();
		return trimmed ? trimmed : null;
	}
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>{selectedUser ? 'User details' : 'Create user'}</Card.Title>
		<Card.Description>{selectedUser?.email ?? 'Invite a user into the directory'}</Card.Description>
	</Card.Header>
	<Card.Content>
		<form class="grid gap-3" onsubmit={submit}>
			<label class="grid gap-1 text-sm">
				<span class="font-medium">Name</span>
				<Input
					name="name"
					value={selectedUser?.name ?? ''}
					disabled={isMutationPending || !canSubmit}
					required
				/>
			</label>

			<label class="grid gap-1 text-sm">
				<span class="font-medium">Email</span>
				<Input
					name="email"
					type="email"
					value={selectedUser?.email ?? ''}
					disabled={isMutationPending || isEditing || !canSubmit}
					required={!isEditing}
				/>
			</label>

			<label class="grid gap-1 text-sm">
				<span class="font-medium">Status</span>
				<select
					name="status"
					class="h-9 rounded-md border bg-background px-3 text-sm"
					value={selectedUser?.status ?? 'invited'}
					disabled={isMutationPending || isEditing || !canCreateUsers}
				>
					<option value="invited">Invited</option>
					<option value="active">Active</option>
					<option value="inactive">Inactive</option>
					<option value="suspended">Suspended</option>
					<option value="deleted">Deleted</option>
				</select>
			</label>

			<label class="grid gap-1 text-sm">
				<span class="font-medium">Primary unit ID</span>
				<Input
					name="primaryOrganizationUnitId"
					value={selectedUser?.primaryOrganizationUnitId ?? ''}
					disabled={isMutationPending || !(canCreateUsers || canAssignUserUnits)}
					placeholder="UUID"
				/>
			</label>

			<label class="grid gap-1 text-sm">
				<span class="font-medium">Manager user ID</span>
				<Input
					name="managerUserId"
					value={selectedUser?.managerUserId ?? ''}
					disabled={isMutationPending || !canSubmit}
					placeholder="UUID"
				/>
			</label>

			<div class="grid gap-3 sm:grid-cols-2">
				<label class="grid gap-1 text-sm">
					<span class="font-medium">Job title</span>
					<Input
						name="jobTitle"
						value={selectedUser?.jobTitle ?? ''}
						disabled={isMutationPending || !canSubmit}
					/>
				</label>
				<label class="grid gap-1 text-sm">
					<span class="font-medium">Phone</span>
					<Input name="phone" value={selectedUser?.phone ?? ''} disabled={isMutationPending || !canSubmit} />
				</label>
			</div>

			<Button type="submit" class="gap-2" disabled={!canSubmit || isMutationPending}>
				{#if selectedUser}
					<SaveIcon class="size-4" />
					Save user
				{:else}
					<UserPlusIcon class="size-4" />
					Create user
				{/if}
			</Button>
		</form>
	</Card.Content>
</Card.Root>
