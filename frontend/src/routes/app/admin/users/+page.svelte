<script lang="ts">
	import AlertCircleIcon from '@lucide/svelte/icons/alert-circle';
	import LoaderCircleIcon from '@lucide/svelte/icons/loader-circle';
	import SearchIcon from '@lucide/svelte/icons/search';
	import ShieldCheckIcon from '@lucide/svelte/icons/shield-check';
	import Trash2Icon from '@lucide/svelte/icons/trash-2';
	import UserCheckIcon from '@lucide/svelte/icons/user-check';
	import UserMinusIcon from '@lucide/svelte/icons/user-minus';
	import UsersIcon from '@lucide/svelte/icons/users';
	import { createMutation, createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import type { PageProps } from './$types';
	import {
		USERS_QUERY_KEY,
		activateUser,
		assignUserGroup,
		assignUserRole,
		createUser,
		deactivateUser,
		fetchUsers,
		removeUserGroup,
		removeUserRole,
		setPrimaryOrganizationUnit,
		softDeleteUser,
		suspendUser,
		updateUser,
		type AssignUserGroupVariables,
		type AssignUserRoleVariables,
		type CreateUserInput,
		type RemoveUserGroupVariables,
		type RemoveUserRoleVariables,
		type SetPrimaryOrganizationUnitVariables,
		type UpdateUserVariables,
		type User,
		type UserListResponse,
		type UserStatusVariables
	} from './api';
	import UserEditor from './user-editor.svelte';
	import UserRoleAssignments from './user-role-assignments.svelte';

	type UsersPageData = {
		canViewUsers: boolean;
		canManageUsers: boolean;
		canCreateUsers: boolean;
		canUpdateUsers: boolean;
		canDeleteUsers: boolean;
		canActivateUsers: boolean;
		canSuspendUsers: boolean;
		canAssignUserRoles: boolean;
		canAssignUserGroups: boolean;
		canAssignUserUnits: boolean;
		selectedUserId: string | null;
	};

	type StatusMutationVariables = UserStatusVariables & {
		action: 'activate' | 'deactivate' | 'suspend' | 'delete';
	};

	let { data }: PageProps = $props();

	const pageData = $derived(data as UsersPageData);
	const queryClient = useQueryClient();

	let selectedOverride = $state<string | null>(null);
	let filter = $state('');
	let localError = $state('');

	const usersQuery = createQuery<UserListResponse, Error>(() => ({
		queryKey: USERS_QUERY_KEY,
		queryFn: () => fetchUsers(fetch),
		enabled: pageData.canViewUsers
	}));
	const createMutationState = createMutation<User, Error, CreateUserInput>(() => ({
		mutationKey: ['admin-users', 'create'],
		mutationFn: (input) => createUser(fetch, input),
		onSuccess: async (user) => {
			selectedOverride = user.id;
			await queryClient.invalidateQueries({ queryKey: USERS_QUERY_KEY });
		}
	}));
	const updateMutationState = createMutation<User, Error, UpdateUserVariables>(() => ({
		mutationKey: ['admin-users', 'update'],
		mutationFn: (variables) => updateUser(fetch, variables),
		onSuccess: async (user) => {
			selectedOverride = user.id;
			await queryClient.invalidateQueries({ queryKey: USERS_QUERY_KEY });
		}
	}));
	const statusMutationState = createMutation<User, Error, StatusMutationVariables>(() => ({
		mutationKey: ['admin-users', 'status'],
		mutationFn: ({ action, id }) => {
			if (action === 'activate') return activateUser(fetch, { id });
			if (action === 'deactivate') return deactivateUser(fetch, { id });
			if (action === 'suspend') return suspendUser(fetch, { id });
			return softDeleteUser(fetch, { id });
		},
		onSuccess: async (user) => {
			selectedOverride = user.id;
			await queryClient.invalidateQueries({ queryKey: USERS_QUERY_KEY });
		}
	}));
	const primaryUnitMutationState = createMutation<
		User,
		Error,
		SetPrimaryOrganizationUnitVariables
	>(() => ({
		mutationKey: ['admin-users', 'primary-organization-unit'],
		mutationFn: (variables) => setPrimaryOrganizationUnit(fetch, variables),
		onSuccess: async (user) => {
			selectedOverride = user.id;
			await queryClient.invalidateQueries({ queryKey: USERS_QUERY_KEY });
		}
	}));
	const assignRoleMutationState = createMutation<unknown, Error, AssignUserRoleVariables>(() => ({
		mutationKey: ['admin-users', 'assign-role'],
		mutationFn: (variables) => assignUserRole(fetch, variables),
		onSuccess: async () => {
			await queryClient.invalidateQueries({ queryKey: USERS_QUERY_KEY });
		}
	}));
	const removeRoleMutationState = createMutation<unknown, Error, RemoveUserRoleVariables>(() => ({
		mutationKey: ['admin-users', 'remove-role'],
		mutationFn: (variables) => removeUserRole(fetch, variables),
		onSuccess: async () => {
			await queryClient.invalidateQueries({ queryKey: USERS_QUERY_KEY });
		}
	}));
	const assignGroupMutationState = createMutation<unknown, Error, AssignUserGroupVariables>(() => ({
		mutationKey: ['admin-users', 'assign-group'],
		mutationFn: (variables) => assignUserGroup(fetch, variables),
		onSuccess: async () => {
			await queryClient.invalidateQueries({ queryKey: USERS_QUERY_KEY });
		}
	}));
	const removeGroupMutationState = createMutation<unknown, Error, RemoveUserGroupVariables>(() => ({
		mutationKey: ['admin-users', 'remove-group'],
		mutationFn: (variables) => removeUserGroup(fetch, variables),
		onSuccess: async () => {
			await queryClient.invalidateQueries({ queryKey: USERS_QUERY_KEY });
		}
	}));

	const users = $derived(usersQuery.data?.users ?? []);
	const normalizedFilter = $derived(filter.trim().toLowerCase());
	const filteredUsers = $derived(
		normalizedFilter
			? users.filter((user) =>
					[user.name, user.email, user.status, user.jobTitle ?? '']
						.join(' ')
						.toLowerCase()
						.includes(normalizedFilter)
				)
			: users
	);
	const requestedSelectedId = $derived(selectedOverride ?? pageData.selectedUserId);
	const selectedUser = $derived(
		(users.find((user) => user.id === requestedSelectedId) ?? filteredUsers[0]) || null
	);
	const isMutationPending = $derived(
		createMutationState.isPending ||
			updateMutationState.isPending ||
			statusMutationState.isPending ||
			primaryUnitMutationState.isPending ||
			assignRoleMutationState.isPending ||
			removeRoleMutationState.isPending ||
			assignGroupMutationState.isPending ||
			removeGroupMutationState.isPending
	);
	const mutationError = $derived.by(
		() =>
			localError ||
			createMutationState.error?.message ||
			updateMutationState.error?.message ||
			statusMutationState.error?.message ||
			primaryUnitMutationState.error?.message ||
			assignRoleMutationState.error?.message ||
			removeRoleMutationState.error?.message ||
			assignGroupMutationState.error?.message ||
			removeGroupMutationState.error?.message ||
			''
	);

	function statusVariant(status: User['status']) {
		if (status === 'active') return 'default';
		if (status === 'invited') return 'secondary';
		if (status === 'deleted' || status === 'suspended') return 'destructive';
		return 'outline';
	}

	async function runCreate(input: CreateUserInput) {
		resetMutationErrors();
		await createMutationState.mutateAsync(input);
	}

	async function runUpdate(id: string, input: UpdateUserVariables['input']) {
		resetMutationErrors();
		await updateMutationState.mutateAsync({ id, input });
	}

	async function runSetPrimaryOrganizationUnit(id: string, organizationUnitId: string | null) {
		resetMutationErrors();
		await primaryUnitMutationState.mutateAsync({ id, organizationUnitId });
	}

	async function runStatus(id: string, action: StatusMutationVariables['action']) {
		resetMutationErrors();
		if (action === 'suspend' && !confirm('Suspend this user?')) return;
		if (action === 'delete' && !confirm('Delete this user?')) return;
		await statusMutationState.mutateAsync({ id, action });
	}

	async function runAssignRole(id: string, input: AssignUserRoleVariables['input']) {
		resetMutationErrors();
		await assignRoleMutationState.mutateAsync({ id, input });
	}

	async function runRemoveRole(id: string, assignmentId: string) {
		resetMutationErrors();
		await removeRoleMutationState.mutateAsync({ id, assignmentId });
	}

	async function runAssignGroup(id: string, groupId: string) {
		resetMutationErrors();
		await assignGroupMutationState.mutateAsync({ id, groupId });
	}

	async function runRemoveGroup(id: string, groupId: string) {
		resetMutationErrors();
		await removeGroupMutationState.mutateAsync({ id, groupId });
	}

	function resetMutationErrors() {
		localError = '';
		createMutationState.reset();
		updateMutationState.reset();
		statusMutationState.reset();
		primaryUnitMutationState.reset();
		assignRoleMutationState.reset();
		removeRoleMutationState.reset();
		assignGroupMutationState.reset();
		removeGroupMutationState.reset();
	}
</script>

<svelte:head>
	<title>Users | Syncra DMS</title>
</svelte:head>

<div class="flex flex-1 flex-col gap-4 p-4 lg:p-6">
	<div class="flex flex-col gap-3 lg:flex-row lg:items-start lg:justify-between">
		<div class="min-w-0">
			<div class="flex items-center gap-2">
				<UsersIcon class="size-5 text-primary" />
				<h2 class="truncate text-xl font-semibold tracking-normal">Users</h2>
			</div>
			<p class="mt-1 text-sm text-muted-foreground">
				{filteredUsers.length} of {users.length} visible
			</p>
		</div>

		<label class="relative w-full lg:w-80">
			<span class="sr-only">Filter users</span>
			<SearchIcon class="pointer-events-none absolute left-3 top-2.5 size-4 text-muted-foreground" />
			<input
				class="h-9 w-full rounded-md border bg-background pl-9 pr-3 text-sm"
				bind:value={filter}
				placeholder="Search users"
			/>
		</label>
	</div>

	{#if !pageData.canViewUsers}
		<div class="rounded-md border border-destructive/30 bg-destructive/10 p-3 text-sm text-destructive">
			User access is unavailable.
		</div>
	{:else}
		{#if usersQuery.isError}
			<div
				class="flex flex-wrap items-center gap-2 rounded-md border border-destructive/30 bg-destructive/10 p-3 text-sm text-destructive"
				role="alert"
			>
				<AlertCircleIcon class="size-4 shrink-0" />
				<span>{usersQuery.error.message}</span>
				<Button type="button" variant="outline" size="xs" class="ms-auto" onclick={() => usersQuery.refetch()}>
					Retry
				</Button>
			</div>
		{/if}

		{#if usersQuery.isLoading || usersQuery.isFetching}
			<div
				class="flex items-center gap-2 rounded-md border bg-muted/30 p-3 text-sm text-muted-foreground"
				role="status"
				aria-live="polite"
			>
				<LoaderCircleIcon class="size-4 shrink-0 animate-spin" />
				<span>Loading users</span>
			</div>
		{/if}

		{#if mutationError}
			<div
				class="flex items-center gap-2 rounded-md border border-destructive/30 bg-destructive/10 p-3 text-sm text-destructive"
				role="alert"
			>
				<AlertCircleIcon class="size-4 shrink-0" />
				<span>{mutationError}</span>
			</div>
		{/if}

		<div class="grid min-h-0 gap-4 xl:grid-cols-[minmax(0,1fr)_24rem]">
			<Card.Root>
				<Card.Header>
					<Card.Title>User directory</Card.Title>
					<Card.Description>{pageData.canManageUsers ? 'Manage accounts' : 'Read only'}</Card.Description>
				</Card.Header>
				<Card.Content>
					<div class="overflow-x-auto">
						<table class="w-full min-w-[48rem] text-sm">
							<thead class="border-b text-left text-xs uppercase text-muted-foreground">
								<tr>
									<th class="px-2 py-2 font-medium">Name</th>
									<th class="px-2 py-2 font-medium">Status</th>
									<th class="px-2 py-2 font-medium">Unit</th>
									<th class="px-2 py-2 font-medium">Last login</th>
									<th class="px-2 py-2 text-right font-medium">Actions</th>
								</tr>
							</thead>
							<tbody>
								{#each filteredUsers as user (user.id)}
									<tr
										class="border-b last:border-0 {selectedUser?.id === user.id ? 'bg-muted/60' : ''}"
									>
										<td class="px-2 py-2">
											<button
												type="button"
												class="min-w-0 text-left"
												onclick={() => (selectedOverride = user.id)}
											>
												<span class="block truncate font-medium">{user.name}</span>
												<span class="block truncate text-xs text-muted-foreground">{user.email}</span>
											</button>
										</td>
										<td class="px-2 py-2">
											<Badge variant={statusVariant(user.status)}>{user.status}</Badge>
										</td>
										<td class="px-2 py-2 text-muted-foreground">
											{user.primaryOrganizationUnitId ?? 'None'}
										</td>
										<td class="px-2 py-2 text-muted-foreground">
											{user.lastLoginAt ?? 'Never'}
										</td>
										<td class="px-2 py-2">
											<div class="flex justify-end gap-1">
												<Button
													type="button"
													variant="ghost"
													size="icon-sm"
													title="Activate"
													disabled={!pageData.canActivateUsers || isMutationPending}
													onclick={() => runStatus(user.id, 'activate')}
												>
													<UserCheckIcon class="size-4" />
												</Button>
												<Button
													type="button"
													variant="ghost"
													size="icon-sm"
													title="Deactivate"
													disabled={!pageData.canUpdateUsers || isMutationPending}
													onclick={() => runStatus(user.id, 'deactivate')}
												>
													<UserMinusIcon class="size-4" />
												</Button>
												<Button
													type="button"
													variant="ghost"
													size="icon-sm"
													title="Suspend"
													disabled={!pageData.canSuspendUsers || isMutationPending}
													onclick={() => runStatus(user.id, 'suspend')}
												>
													<ShieldCheckIcon class="size-4" />
												</Button>
												<Button
													type="button"
													variant="ghost"
													size="icon-sm"
													title="Delete"
													disabled={!pageData.canDeleteUsers || isMutationPending}
													onclick={() => runStatus(user.id, 'delete')}
												>
													<Trash2Icon class="size-4" />
												</Button>
											</div>
										</td>
									</tr>
								{/each}
							</tbody>
						</table>
					</div>
				</Card.Content>
			</Card.Root>

			<div class="grid gap-4">
				<UserEditor
					selectedUser={selectedUser}
					canCreateUsers={pageData.canCreateUsers}
					canUpdateUsers={pageData.canUpdateUsers}
					canAssignUserUnits={pageData.canAssignUserUnits}
					{isMutationPending}
					onCreate={runCreate}
					onUpdate={runUpdate}
					onSetPrimaryOrganizationUnit={runSetPrimaryOrganizationUnit}
				/>
				<UserRoleAssignments
					selectedUser={selectedUser}
					canAssignUserRoles={pageData.canAssignUserRoles}
					canAssignUserGroups={pageData.canAssignUserGroups}
					{isMutationPending}
					onAssignRole={runAssignRole}
					onRemoveRole={runRemoveRole}
					onAssignGroup={runAssignGroup}
					onRemoveGroup={runRemoveGroup}
				/>
			</div>
		</div>
	{/if}
</div>
