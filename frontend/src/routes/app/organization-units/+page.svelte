<script lang="ts">
	import { navigating } from '$app/state';
	import AlertCircleIcon from '@lucide/svelte/icons/alert-circle';
	import Building2Icon from '@lucide/svelte/icons/building-2';
	import PlusIcon from '@lucide/svelte/icons/plus';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import type { PageProps } from './$types';
	import {
		collectMoveTargets,
		countUnits,
		findUnit,
		flattenUnitTree,
		selectInitialUnit,
		type OrganizationUnitNode
	} from './tree';
	import UnitDetailsPanel from './unit-details-panel.svelte';
	import UnitTree from './unit-tree.svelte';

	type OrganizationUnitsPageData = {
		units: OrganizationUnitNode[];
		loadError: string | null;
		canManageOrganizationUnits: boolean;
		selectedId: string | null;
	};

	type OrganizationUnitsActionData = {
		error?: string;
		action?: string;
		selectedId?: string | null;
		success?: boolean;
		values?: {
			id?: string;
			parentId?: string | null;
			name?: string;
			code?: string;
			description?: string;
		};
	};

	let { data, form }: PageProps = $props();

	const pageData = $derived(data as OrganizationUnitsPageData);
	const actionData = $derived(form as OrganizationUnitsActionData | null | undefined);
	let selectedOverride = $state<string | null>(null);
	const requestedSelectedId = $derived(
		selectedOverride ?? actionData?.selectedId ?? pageData.selectedId
	);

	const flatUnits = $derived(flattenUnitTree(pageData.units));
	const selectedId = $derived(
		requestedSelectedId && findUnit(pageData.units, requestedSelectedId)
			? requestedSelectedId
			: (selectInitialUnit(pageData.units)?.id ?? null)
	);
	const selected = $derived(findUnit(pageData.units, selectedId));
	const parentOptions = $derived(
		selectedId ? collectMoveTargets(pageData.units, selectedId) : flatUnits
	);
	const unitCount = $derived(countUnits(pageData.units));
	const selectedFormValue = $derived(selectedId ?? '');
	const createRootValues = $derived(
		actionData?.action === 'create' && actionData.values?.parentId === null
			? actionData.values
			: null
	);
	const isRouteLoading = $derived(navigating.to?.url.pathname === '/app/organization-units');
</script>

<svelte:head>
	<title>Organization Units | Syncra DMS</title>
</svelte:head>

<section class="mx-auto grid max-w-6xl gap-4 px-4 py-6">
	<div class="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
		<div class="min-w-0">
			<div class="flex items-center gap-2">
				<Building2Icon class="size-5 text-primary" />
				<h2 class="truncate text-xl font-semibold tracking-normal">Organization Units</h2>
			</div>
			<p class="mt-1 text-sm text-muted-foreground">{unitCount} active units</p>
		</div>

		{#if pageData.canManageOrganizationUnits}
			<form
				method="POST"
				action="?/create"
				class="grid w-full gap-2 sm:w-auto sm:grid-cols-[minmax(12rem,1fr)_7rem_auto]"
			>
				<label class="sr-only" for="root-unit-name">Root unit name</label>
				<input
					id="root-unit-name"
					class="h-9 rounded-md border bg-background px-3 text-sm"
					name="name"
					placeholder="Root name"
					value={createRootValues?.name ?? ''}
					required
				/>
				<label class="sr-only" for="root-unit-code">Root unit code</label>
				<input
					id="root-unit-code"
					class="h-9 rounded-md border bg-background px-3 text-sm"
					name="code"
					placeholder="Code"
					value={createRootValues?.code ?? ''}
				/>
				<input type="hidden" name="description" value={createRootValues?.description ?? ''} />
				<input type="hidden" name="parentId" value="" />
				<input type="hidden" name="selectedId" value={selectedFormValue} />
				<Button type="submit" size="sm" class="gap-2">
					<PlusIcon class="size-4" />
					New root
				</Button>
			</form>
		{/if}
	</div>

	{#if pageData.loadError}
		<div
			class="flex flex-wrap items-center gap-2 rounded-md border border-destructive/30 bg-destructive/10 p-3 text-sm text-destructive"
			role="alert"
		>
			<AlertCircleIcon class="size-4 shrink-0" />
			<span>{pageData.loadError}</span>
			<Button href="/app/organization-units" variant="outline" size="xs" class="ms-auto">
				Retry
			</Button>
		</div>
	{/if}

	{#if isRouteLoading}
		<div
			class="flex items-center gap-2 rounded-md border bg-muted/30 p-3 text-sm text-muted-foreground"
			role="status"
			aria-live="polite"
		>
			<Building2Icon class="size-4 shrink-0" />
			<span>Loading organization units</span>
		</div>
	{/if}

	{#if actionData?.error}
		<div
			class="flex items-center gap-2 rounded-md border border-destructive/30 bg-destructive/10 p-3 text-sm text-destructive"
			role="alert"
		>
			<AlertCircleIcon class="size-4 shrink-0" />
			<span>{actionData.error}</span>
		</div>
	{/if}

	<div class="grid gap-4 lg:grid-cols-[20rem_minmax(0,1fr)]">
		<Card.Root>
			<Card.Header>
				<Card.Title>Active tree</Card.Title>
				<Card.Description>{flatUnits.length} visible</Card.Description>
			</Card.Header>
			<Card.Content>
				<UnitTree
					units={flatUnits}
					{selectedId}
					onSelect={(id) => (selectedOverride = id)}
				/>
			</Card.Content>
		</Card.Root>

		<UnitDetailsPanel
			{selected}
			canManage={pageData.canManageOrganizationUnits}
			{parentOptions}
			selectedId={selectedFormValue}
			{actionData}
		/>
	</div>
</section>
