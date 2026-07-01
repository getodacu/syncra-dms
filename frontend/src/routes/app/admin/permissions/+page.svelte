<script lang="ts">
	import AlertCircleIcon from '@lucide/svelte/icons/alert-circle';
	import KeyRoundIcon from '@lucide/svelte/icons/key-round';
	import LoaderCircleIcon from '@lucide/svelte/icons/loader-circle';
	import { createQuery } from '@tanstack/svelte-query';
	import { Badge } from '$lib/components/ui/badge';
	import * as Card from '$lib/components/ui/card';
	import type { PageProps } from './$types';
	import {
		PERMISSIONS_QUERY_KEY,
		fetchPermissions,
		groupPermissionsByCategory,
		type PermissionListResponse
	} from '../roles/api';

	type PermissionsPageData = {
		canViewPermissions: boolean;
		canManagePermissions: boolean;
	};

	let { data }: PageProps = $props();

	const pageData = $derived(data as PermissionsPageData);
	const permissionsQuery = createQuery<PermissionListResponse, Error>(() => ({
		queryKey: PERMISSIONS_QUERY_KEY,
		queryFn: () => fetchPermissions(fetch),
		enabled: pageData.canViewPermissions
	}));
	const groupedPermissions = $derived(
		groupPermissionsByCategory(permissionsQuery.data?.permissions ?? [])
	);
</script>

<svelte:head>
	<title>Permissions | Syncra DMS</title>
</svelte:head>

<div class="flex flex-1 flex-col gap-4 p-4 lg:p-6">
	<div class="flex items-center gap-2">
		<KeyRoundIcon class="size-5 text-primary" />
		<h2 class="text-xl font-semibold tracking-normal">Permissions</h2>
	</div>

	{#if !pageData.canViewPermissions}
		<div class="rounded-md border border-destructive/30 bg-destructive/10 p-3 text-sm text-destructive">
			Permission registry access is unavailable.
		</div>
	{:else}
		{#if permissionsQuery.isError}
			<div class="flex items-center gap-2 rounded-md border border-destructive/30 bg-destructive/10 p-3 text-sm text-destructive">
				<AlertCircleIcon class="size-4" />
				<span>{permissionsQuery.error.message}</span>
			</div>
		{/if}

		{#if permissionsQuery.isLoading || permissionsQuery.isFetching}
			<div class="flex items-center gap-2 rounded-md border bg-muted/30 p-3 text-sm text-muted-foreground">
				<LoaderCircleIcon class="size-4 animate-spin" />
				<span>Loading permissions</span>
			</div>
		{/if}

		<div class="grid gap-4 lg:grid-cols-2">
			{#each groupedPermissions as group (group.category)}
				<Card.Root>
					<Card.Header>
						<Card.Title>{group.category}</Card.Title>
						<Card.Description>{group.permissions.length} permissions</Card.Description>
					</Card.Header>
					<Card.Content class="grid gap-2">
						{#each group.permissions as permission (permission.id)}
							<div class="rounded-md border p-3">
								<div class="flex items-start justify-between gap-3">
									<div class="min-w-0">
										<p class="truncate text-sm font-medium">{permission.name}</p>
										<p class="truncate font-mono text-xs text-muted-foreground">{permission.code}</p>
									</div>
									{#if permission.isSystem}
										<Badge variant="secondary">System</Badge>
									{/if}
								</div>
								{#if permission.description}
									<p class="mt-2 text-xs text-muted-foreground">{permission.description}</p>
								{/if}
							</div>
						{/each}
					</Card.Content>
				</Card.Root>
			{/each}
		</div>
	{/if}
</div>
