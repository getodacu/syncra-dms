<script lang="ts">
	import AlertCircleIcon from '@lucide/svelte/icons/alert-circle';
	import LoaderCircleIcon from '@lucide/svelte/icons/loader-circle';
	import LockIcon from '@lucide/svelte/icons/lock';
	import PlusIcon from '@lucide/svelte/icons/plus';
	import SaveIcon from '@lucide/svelte/icons/save';
	import ShieldCheckIcon from '@lucide/svelte/icons/shield-check';
	import Trash2Icon from '@lucide/svelte/icons/trash-2';
	import { createMutation, createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import type { PageProps } from './$types';
	import {
		PERMISSIONS_QUERY_KEY,
		ROLES_QUERY_KEY,
		assignRolePermission,
		createRole,
		deleteRole,
		fetchPermissions,
		fetchRolePermissions,
		fetchRoles,
		removeRolePermission,
		updateRole,
		type CreateRoleInput,
		type Permission,
		type PermissionListResponse,
		type Role,
		type RoleListResponse,
		type RolePermissionVariables,
		type UpdateRoleVariables
	} from './api';
	import PermissionMatrix from './permission-matrix.svelte';

	type RolesPageData = {
		canViewRoles: boolean;
		canManageRoles: boolean;
		canCreateRoles: boolean;
		canUpdateRoles: boolean;
		canDeleteRoles: boolean;
		canAssignRolePermissions: boolean;
		selectedRoleId: string | null;
	};

	let { data }: PageProps = $props();

	const pageData = $derived(data as RolesPageData);
	const queryClient = useQueryClient();

	let selectedOverride = $state<string | null>(null);
	let localError = $state('');

	const rolesQuery = createQuery<RoleListResponse, Error>(() => ({
		queryKey: ROLES_QUERY_KEY,
		queryFn: () => fetchRoles(fetch),
		enabled: pageData.canViewRoles
	}));
	const permissionsQuery = createQuery<PermissionListResponse, Error>(() => ({
		queryKey: PERMISSIONS_QUERY_KEY,
		queryFn: () => fetchPermissions(fetch),
		enabled: pageData.canViewRoles
	}));
	const createMutationState = createMutation<Role, Error, CreateRoleInput>(() => ({
		mutationKey: ['admin-roles', 'create'],
		mutationFn: (input) => createRole(fetch, input),
		onSuccess: async (role) => {
			selectedOverride = role.id;
			await queryClient.invalidateQueries({ queryKey: ROLES_QUERY_KEY });
		}
	}));
	const updateMutationState = createMutation<Role, Error, UpdateRoleVariables>(() => ({
		mutationKey: ['admin-roles', 'update'],
		mutationFn: (variables) => updateRole(fetch, variables),
		onSuccess: async (role) => {
			selectedOverride = role.id;
			await queryClient.invalidateQueries({ queryKey: ROLES_QUERY_KEY });
		}
	}));
	const deleteMutationState = createMutation<unknown, Error, { id: string }>(() => ({
		mutationKey: ['admin-roles', 'delete'],
		mutationFn: (variables) => deleteRole(fetch, variables),
		onSuccess: async () => {
			selectedOverride = null;
			await queryClient.invalidateQueries({ queryKey: ROLES_QUERY_KEY });
		}
	}));
	const roles = $derived(rolesQuery.data?.roles ?? []);
	const selectedRoleId = $derived(selectedOverride ?? pageData.selectedRoleId ?? roles[0]?.id ?? null);
	const rolePermissionsQuery = createQuery<PermissionListResponse, Error>(() => ({
		queryKey: ['admin-roles', selectedRoleId, 'permissions'],
		queryFn: () => fetchRolePermissions(fetch, { id: selectedRoleId ?? '' }),
		enabled: Boolean(selectedRoleId) && pageData.canViewRoles
	}));
	const assignPermissionMutationState = createMutation<Permission, Error, RolePermissionVariables>(
		() => ({
			mutationKey: ['admin-roles', 'assign-permission'],
			mutationFn: (variables) => assignRolePermission(fetch, variables),
			onSuccess: async () => {
				await queryClient.invalidateQueries({ queryKey: ['admin-roles', selectedRoleId, 'permissions'] });
			}
		})
	);
	const removePermissionMutationState = createMutation<unknown, Error, RolePermissionVariables>(
		() => ({
			mutationKey: ['admin-roles', 'remove-permission'],
			mutationFn: (variables) => removeRolePermission(fetch, variables),
			onSuccess: async () => {
				await queryClient.invalidateQueries({ queryKey: ['admin-roles', selectedRoleId, 'permissions'] });
			}
		})
	);

	const selectedRole = $derived(roles.find((role) => role.id === selectedRoleId) ?? null);
	const permissions = $derived(permissionsQuery.data?.permissions ?? []);
	const assignedPermissions = $derived(rolePermissionsQuery.data?.permissions ?? []);
	const isMutationPending = $derived(
		createMutationState.isPending ||
			updateMutationState.isPending ||
			deleteMutationState.isPending ||
			assignPermissionMutationState.isPending ||
			removePermissionMutationState.isPending
	);
	const mutationError = $derived.by(
		() =>
			localError ||
			createMutationState.error?.message ||
			updateMutationState.error?.message ||
			deleteMutationState.error?.message ||
			assignPermissionMutationState.error?.message ||
			removePermissionMutationState.error?.message ||
			''
	);

	async function submitRole(event: SubmitEvent) {
		event.preventDefault();
		resetMutationErrors();
		const form = event.currentTarget as HTMLFormElement;
		const formData = new FormData(form);
		const input = {
			name: textField(formData, 'name'),
			code: textField(formData, 'code'),
			description: nullableTextField(formData, 'description'),
			isActive: formData.get('isActive') === 'on'
		};
		if (!input.name.trim()) {
			localError = 'Role name is required';
			return;
		}
		if (selectedRole) {
			if (selectedRole.isSystem || !pageData.canUpdateRoles) return;
			await updateMutationState.mutateAsync({
				id: selectedRole.id,
				input: {
					name: input.name,
					description: input.description,
					isActive: input.isActive
				}
			});
			return;
		}
		if (!pageData.canCreateRoles) return;
		await createMutationState.mutateAsync(input);
		form.reset();
	}

	async function runDeleteRole(role: Role) {
		resetMutationErrors();
		if (role.isSystem || !pageData.canDeleteRoles || !confirm('Delete this role?')) return;
		await deleteMutationState.mutateAsync({ id: role.id });
	}

	async function runAssignPermission(permissionId: string) {
		if (!selectedRoleId) return;
		resetMutationErrors();
		await assignPermissionMutationState.mutateAsync({ id: selectedRoleId, permissionId });
	}

	async function runRemovePermission(permissionId: string) {
		if (!selectedRoleId) return;
		resetMutationErrors();
		await removePermissionMutationState.mutateAsync({ id: selectedRoleId, permissionId });
	}

	function textField(data: FormData, key: string) {
		const value = data.get(key);
		return typeof value === 'string' ? value : '';
	}

	function nullableTextField(data: FormData, key: string) {
		const value = textField(data, key).trim();
		return value ? value : null;
	}

	function resetMutationErrors() {
		localError = '';
		createMutationState.reset();
		updateMutationState.reset();
		deleteMutationState.reset();
		assignPermissionMutationState.reset();
		removePermissionMutationState.reset();
	}
</script>

<svelte:head>
	<title>Roles | Syncra DMS</title>
</svelte:head>

<div class="flex flex-1 flex-col gap-4 p-4 lg:p-6">
	<div class="flex items-center gap-2">
		<ShieldCheckIcon class="size-5 text-primary" />
		<h2 class="text-xl font-semibold tracking-normal">Roles</h2>
	</div>

	{#if !pageData.canViewRoles}
		<div class="rounded-md border border-destructive/30 bg-destructive/10 p-3 text-sm text-destructive">
			Role access is unavailable.
		</div>
	{:else}
		{#if rolesQuery.isError || permissionsQuery.isError || rolePermissionsQuery.isError}
			<div class="flex items-center gap-2 rounded-md border border-destructive/30 bg-destructive/10 p-3 text-sm text-destructive">
				<AlertCircleIcon class="size-4" />
				<span>{rolesQuery.error?.message ?? permissionsQuery.error?.message ?? rolePermissionsQuery.error?.message}</span>
			</div>
		{/if}

		{#if rolesQuery.isLoading || permissionsQuery.isLoading}
			<div class="flex items-center gap-2 rounded-md border bg-muted/30 p-3 text-sm text-muted-foreground">
				<LoaderCircleIcon class="size-4 animate-spin" />
				<span>Loading roles</span>
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
					<Card.Title>Role list</Card.Title>
					<Card.Description>{roles.length} roles</Card.Description>
				</Card.Header>
				<Card.Content class="grid gap-2">
					{#each roles as role (role.id)}
						<button
							type="button"
							class="rounded-md border p-3 text-left text-sm {selectedRole?.id === role.id ? 'bg-muted' : ''}"
							onclick={() => (selectedOverride = role.id)}
						>
							<span class="flex items-center gap-2 font-medium">
								{role.name}
								{#if role.isSystem}
									<LockIcon class="size-3 text-muted-foreground" />
								{/if}
							</span>
							<span class="block font-mono text-xs text-muted-foreground">{role.code}</span>
						</button>
					{/each}
				</Card.Content>
			</Card.Root>

			<Card.Root>
				<Card.Header>
					<Card.Title>{selectedRole ? 'Role permissions' : 'Permission matrix'}</Card.Title>
					<Card.Description>{selectedRole?.name ?? 'Select a role'}</Card.Description>
				</Card.Header>
				<Card.Content>
					<PermissionMatrix
						{permissions}
						{assignedPermissions}
						canAssignRolePermissions={pageData.canAssignRolePermissions && Boolean(selectedRole)}
						{isMutationPending}
						onAssign={runAssignPermission}
						onRemove={runRemovePermission}
					/>
				</Card.Content>
			</Card.Root>

			<Card.Root>
				<Card.Header>
					<Card.Title>{selectedRole ? 'Edit role' : 'Create role'}</Card.Title>
					<Card.Description>
						{selectedRole?.isSystem ? 'System role locked' : 'Custom role settings'}
					</Card.Description>
				</Card.Header>
				<Card.Content>
					<form class="grid gap-3" onsubmit={submitRole}>
						<label class="grid gap-1 text-sm">
							<span class="font-medium">Name</span>
							<Input
								name="name"
								value={selectedRole?.name ?? ''}
								disabled={isMutationPending || (selectedRole ? selectedRole.isSystem || !pageData.canUpdateRoles : !pageData.canCreateRoles)}
								required
							/>
						</label>
						<label class="grid gap-1 text-sm">
							<span class="font-medium">Code</span>
							<Input
								name="code"
								value={selectedRole?.code ?? ''}
								disabled={isMutationPending || Boolean(selectedRole)}
								required={!selectedRole}
							/>
						</label>
						<label class="grid gap-1 text-sm">
							<span class="font-medium">Description</span>
							<textarea
								name="description"
								class="min-h-20 rounded-md border bg-background px-3 py-2 text-sm"
								disabled={isMutationPending || (selectedRole ? selectedRole.isSystem || !pageData.canUpdateRoles : !pageData.canCreateRoles)}
							>{selectedRole?.description ?? ''}</textarea>
						</label>
						<label class="flex items-center gap-2 text-sm">
							<input
								name="isActive"
								type="checkbox"
								checked={selectedRole?.isActive ?? true}
								disabled={isMutationPending || (selectedRole ? selectedRole.isSystem || !pageData.canUpdateRoles : !pageData.canCreateRoles)}
							/>
							<span>Active</span>
						</label>
						<div class="flex gap-2">
							<Button type="submit" class="gap-2" disabled={isMutationPending}>
								{#if selectedRole}
									<SaveIcon class="size-4" />
									Save
								{:else}
									<PlusIcon class="size-4" />
									Create
								{/if}
							</Button>
							{#if selectedRole}
								<Button
									type="button"
									variant="destructive"
									class="gap-2"
									disabled={selectedRole.isSystem || !pageData.canDeleteRoles || isMutationPending}
									onclick={() => runDeleteRole(selectedRole)}
								>
									<Trash2Icon class="size-4" />
									Delete
								</Button>
							{/if}
						</div>
					</form>
				</Card.Content>
			</Card.Root>
		</div>
	{/if}
</div>
