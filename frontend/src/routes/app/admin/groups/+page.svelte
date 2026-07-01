<script lang="ts">
	import AlertCircleIcon from '@lucide/svelte/icons/alert-circle';
	import LoaderCircleIcon from '@lucide/svelte/icons/loader-circle';
	import UserRoundCogIcon from '@lucide/svelte/icons/user-round-cog';
	import { createMutation, createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { Badge } from '$lib/components/ui/badge';
	import * as Card from '$lib/components/ui/card';
	import type { PageProps } from './$types';
	import {
		GROUPS_QUERY_KEY,
		addGroupUser,
		assignGroupRole,
		createGroup,
		deleteGroup,
		fetchGroups,
		removeGroupRole,
		removeGroupUser,
		updateGroup,
		type AssignGroupRoleVariables,
		type CreateGroupInput,
		type Group,
		type GroupListResponse,
		type GroupUserVariables,
		type RemoveGroupRoleVariables,
		type UpdateGroupVariables
	} from './api';
	import GroupEditor from './group-editor.svelte';
	import GroupMembers from './group-members.svelte';
	import GroupRoleAssignments from './group-role-assignments.svelte';

	type GroupsPageData = {
		canViewGroups: boolean;
		canManageGroups: boolean;
		canCreateGroups: boolean;
		canUpdateGroups: boolean;
		canDeleteGroups: boolean;
		canManageGroupUsers: boolean;
		canAssignGroupRoles: boolean;
		selectedGroupId: string | null;
	};

	type MemberMutationVariables = GroupUserVariables & { action: 'add' | 'remove' };
	type RoleMutationVariables =
		| (AssignGroupRoleVariables & { action: 'assign' })
		| (RemoveGroupRoleVariables & { action: 'remove' });

	let { data }: PageProps = $props();

	const pageData = $derived(data as GroupsPageData);
	const queryClient = useQueryClient();

	let selectedOverride = $state<string | null>(null);
	let localError = $state('');

	const groupsQuery = createQuery<GroupListResponse, Error>(() => ({
		queryKey: GROUPS_QUERY_KEY,
		queryFn: () => fetchGroups(fetch),
		enabled: pageData.canViewGroups
	}));
	const createMutationState = createMutation<Group, Error, CreateGroupInput>(() => ({
		mutationKey: ['admin-groups', 'create'],
		mutationFn: (input) => createGroup(fetch, input),
		onSuccess: async (group) => {
			selectedOverride = group.id;
			await queryClient.invalidateQueries({ queryKey: GROUPS_QUERY_KEY });
		}
	}));
	const updateMutationState = createMutation<Group, Error, UpdateGroupVariables>(() => ({
		mutationKey: ['admin-groups', 'update'],
		mutationFn: (variables) => updateGroup(fetch, variables),
		onSuccess: async (group) => {
			selectedOverride = group.id;
			await queryClient.invalidateQueries({ queryKey: GROUPS_QUERY_KEY });
		}
	}));
	const deleteMutationState = createMutation<unknown, Error, { id: string }>(() => ({
		mutationKey: ['admin-groups', 'delete'],
		mutationFn: (variables) => deleteGroup(fetch, variables),
		onSuccess: async () => {
			selectedOverride = null;
			await queryClient.invalidateQueries({ queryKey: GROUPS_QUERY_KEY });
		}
	}));
	const memberMutationState = createMutation<unknown, Error, MemberMutationVariables>(() => ({
		mutationKey: ['admin-groups', 'members'],
		mutationFn: ({ action, id, userId }) =>
			action === 'add' ? addGroupUser(fetch, { id, userId }) : removeGroupUser(fetch, { id, userId }),
		onSuccess: async () => {
			await queryClient.invalidateQueries({ queryKey: GROUPS_QUERY_KEY });
		}
	}));
	const roleMutationState = createMutation<unknown, Error, RoleMutationVariables>(() => ({
		mutationKey: ['admin-groups', 'roles'],
		mutationFn: (variables) =>
			variables.action === 'assign'
				? assignGroupRole(fetch, variables)
				: removeGroupRole(fetch, variables),
		onSuccess: async () => {
			await queryClient.invalidateQueries({ queryKey: GROUPS_QUERY_KEY });
		}
	}));

	const groups = $derived(groupsQuery.data?.groups ?? []);
	const selectedGroupId = $derived(selectedOverride ?? pageData.selectedGroupId ?? groups[0]?.id ?? null);
	const selectedGroup = $derived(groups.find((group) => group.id === selectedGroupId) ?? null);
	const activeGroups = $derived(groups.filter((group) => group.isActive));
	const isMutationPending = $derived(
		createMutationState.isPending ||
			updateMutationState.isPending ||
			deleteMutationState.isPending ||
			memberMutationState.isPending ||
			roleMutationState.isPending
	);
	const mutationError = $derived.by(
		() =>
			localError ||
			createMutationState.error?.message ||
			updateMutationState.error?.message ||
			deleteMutationState.error?.message ||
			memberMutationState.error?.message ||
			roleMutationState.error?.message ||
			''
	);

	async function runCreate(input: CreateGroupInput) {
		resetMutationErrors();
		if (!input.name.trim()) {
			localError = 'Group name is required';
			return;
		}
		await createMutationState.mutateAsync(input);
	}

	async function runUpdate(id: string, input: UpdateGroupVariables['input']) {
		resetMutationErrors();
		await updateMutationState.mutateAsync({ id, input });
	}

	async function runDelete(id: string) {
		resetMutationErrors();
		if (!confirm('Delete this group?')) return;
		await deleteMutationState.mutateAsync({ id });
	}

	async function runAddUser(id: string, userId: string) {
		resetMutationErrors();
		await memberMutationState.mutateAsync({ id, userId, action: 'add' });
	}

	async function runRemoveUser(id: string, userId: string) {
		resetMutationErrors();
		await memberMutationState.mutateAsync({ id, userId, action: 'remove' });
	}

	async function runAssignRole(id: string, input: AssignGroupRoleVariables['input']) {
		resetMutationErrors();
		await roleMutationState.mutateAsync({ id, input, action: 'assign' });
	}

	async function runRemoveRole(id: string, assignmentId: string) {
		resetMutationErrors();
		await roleMutationState.mutateAsync({ id, assignmentId, action: 'remove' });
	}

	function resetMutationErrors() {
		localError = '';
		createMutationState.reset();
		updateMutationState.reset();
		deleteMutationState.reset();
		memberMutationState.reset();
		roleMutationState.reset();
	}
</script>

<svelte:head>
	<title>Groups | Syncra DMS</title>
</svelte:head>

<div class="flex flex-1 flex-col gap-4 p-4 lg:p-6">
	<div class="flex items-center gap-2">
		<UserRoundCogIcon class="size-5 text-primary" />
		<h2 class="text-xl font-semibold tracking-normal">Groups</h2>
	</div>

	{#if !pageData.canViewGroups}
		<div class="rounded-md border border-destructive/30 bg-destructive/10 p-3 text-sm text-destructive">
			Group access is unavailable.
		</div>
	{:else}
		{#if groupsQuery.isError}
			<div class="flex items-center gap-2 rounded-md border border-destructive/30 bg-destructive/10 p-3 text-sm text-destructive">
				<AlertCircleIcon class="size-4" />
				<span>{groupsQuery.error.message}</span>
			</div>
		{/if}
		{#if groupsQuery.isLoading || groupsQuery.isFetching}
			<div class="flex items-center gap-2 rounded-md border bg-muted/30 p-3 text-sm text-muted-foreground">
				<LoaderCircleIcon class="size-4 animate-spin" />
				<span>Loading groups</span>
			</div>
		{/if}
		{#if mutationError}
			<div class="flex items-center gap-2 rounded-md border border-destructive/30 bg-destructive/10 p-3 text-sm text-destructive">
				<AlertCircleIcon class="size-4" />
				<span>{mutationError}</span>
			</div>
		{/if}

		<div class="grid gap-4 xl:grid-cols-[18rem_minmax(0,1fr)_24rem]">
			<Card.Root>
				<Card.Header>
					<Card.Title>Group list</Card.Title>
					<Card.Description>{activeGroups.length} active</Card.Description>
				</Card.Header>
				<Card.Content class="grid gap-2">
					{#each groups as group (group.id)}
						<button
							type="button"
							class="rounded-md border p-3 text-left text-sm {selectedGroup?.id === group.id ? 'bg-muted' : ''}"
							onclick={() => (selectedOverride = group.id)}
						>
							<span class="flex items-center justify-between gap-2 font-medium">
								{group.name}
								<Badge variant={group.isActive ? 'secondary' : 'outline'}>
									{group.isActive ? 'Active' : 'Inactive'}
								</Badge>
							</span>
							<span class="block font-mono text-xs text-muted-foreground">{group.code}</span>
						</button>
					{/each}
				</Card.Content>
			</Card.Root>

			<div class="grid gap-4">
				<GroupMembers
					{selectedGroup}
					canManageGroupUsers={pageData.canManageGroupUsers}
					{isMutationPending}
					onAddUser={runAddUser}
					onRemoveUser={runRemoveUser}
				/>
				<GroupRoleAssignments
					{selectedGroup}
					canAssignGroupRoles={pageData.canAssignGroupRoles}
					{isMutationPending}
					onAssignRole={runAssignRole}
					onRemoveRole={runRemoveRole}
				/>
			</div>

			<GroupEditor
				{selectedGroup}
				canCreateGroups={pageData.canCreateGroups}
				canUpdateGroups={pageData.canUpdateGroups}
				canDeleteGroups={pageData.canDeleteGroups}
				{isMutationPending}
				onCreate={runCreate}
				onUpdate={runUpdate}
				onDelete={runDelete}
			/>
		</div>
	{/if}
</div>
