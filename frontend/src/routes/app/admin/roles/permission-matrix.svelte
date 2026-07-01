<script lang="ts">
	import LinkIcon from '@lucide/svelte/icons/link';
	import UnlinkIcon from '@lucide/svelte/icons/unlink';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { groupPermissionsByCategory, type Permission } from './api';

	type Props = {
		permissions: Permission[];
		assignedPermissions: Permission[];
		canAssignRolePermissions: boolean;
		isMutationPending: boolean;
		onAssign: (permissionId: string) => Promise<void>;
		onRemove: (permissionId: string) => Promise<void>;
	};

	let {
		permissions,
		assignedPermissions,
		canAssignRolePermissions,
		isMutationPending,
		onAssign,
		onRemove
	}: Props = $props();

	const assignedPermissionIds = $derived(new Set(assignedPermissions.map((permission) => permission.id)));
	const groupedPermissions = $derived(groupPermissionsByCategory(permissions));
</script>

<div class="grid gap-4">
	{#each groupedPermissions as group (group.category)}
		<section class="grid gap-2 rounded-md border p-3">
			<div class="flex items-center justify-between gap-3">
				<h3 class="text-sm font-semibold">{group.category}</h3>
				<Badge variant="secondary">{group.permissions.length}</Badge>
			</div>
			<div class="grid gap-2">
				{#each group.permissions as permission (permission.id)}
					{@const assigned = assignedPermissionIds.has(permission.id)}
					<div class="grid gap-2 rounded-md bg-muted/30 p-2 sm:grid-cols-[minmax(0,1fr)_auto]">
						<div class="min-w-0">
							<p class="truncate text-sm font-medium">{permission.name}</p>
							<p class="truncate font-mono text-xs text-muted-foreground">{permission.code}</p>
							{#if permission.description}
								<p class="mt-1 text-xs text-muted-foreground">{permission.description}</p>
							{/if}
						</div>
						<Button
							type="button"
							variant={assigned ? 'ghost' : 'outline'}
							size="sm"
							class="gap-2"
							disabled={!canAssignRolePermissions || isMutationPending}
							onclick={() => (assigned ? onRemove(permission.id) : onAssign(permission.id))}
						>
							{#if assigned}
								<UnlinkIcon class="size-4" />
								Remove
							{:else}
								<LinkIcon class="size-4" />
								Assign
							{/if}
						</Button>
					</div>
				{/each}
			</div>
		</section>
	{/each}
</div>
