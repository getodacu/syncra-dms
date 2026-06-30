<script lang="ts">
	import AlertCircleIcon from '@lucide/svelte/icons/alert-circle';
	import Building2Icon from '@lucide/svelte/icons/building-2';
	import LoaderCircleIcon from '@lucide/svelte/icons/loader-circle';
	import PlusIcon from '@lucide/svelte/icons/plus';
	import { createMutation, createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import type { PageProps } from './$types';
	import {
		archiveOrganizationUnit,
		createOrganizationUnit,
		fetchOrganizationUnitTree,
		moveOrganizationUnit,
		ORGANIZATION_UNITS_QUERY_KEY,
		updateOrganizationUnit,
		type ArchiveOrganizationUnitResponse,
		type ArchiveOrganizationUnitVariables,
		type MoveOrganizationUnitVariables,
		type OrganizationUnitInput,
		type OrganizationUnitListResponse,
		type UpdateOrganizationUnitVariables
	} from './api';
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
		canManageOrganizationUnits: boolean;
		selectedId: string | null;
	};

	let { data }: PageProps = $props();

	const pageData = $derived(data as OrganizationUnitsPageData);
	const queryClient = useQueryClient();

	let selectedOverride = $state<string | null>(null);
	let rootName = $state('');
	let rootCode = $state('');
	let localError = $state('');

	const organizationUnitsQuery = createQuery<OrganizationUnitListResponse, Error>(() => ({
		queryKey: ORGANIZATION_UNITS_QUERY_KEY,
		queryFn: () => fetchOrganizationUnitTree(fetch)
	}));
	const createMutationState = createMutation<OrganizationUnitNode, Error, OrganizationUnitInput>(
		() => ({
			mutationKey: ['organization-units', 'create'],
			mutationFn: (input) => createOrganizationUnit(fetch, input),
			onSuccess: async (unit) => {
				selectedOverride = unit.id;
				await queryClient.invalidateQueries({ queryKey: ORGANIZATION_UNITS_QUERY_KEY });
			}
		})
	);
	const updateMutationState = createMutation<
		OrganizationUnitNode,
		Error,
		UpdateOrganizationUnitVariables
	>(() => ({
		mutationKey: ['organization-units', 'update'],
		mutationFn: (variables) => updateOrganizationUnit(fetch, variables),
		onSuccess: async (unit) => {
			selectedOverride = unit.id;
			await queryClient.invalidateQueries({ queryKey: ORGANIZATION_UNITS_QUERY_KEY });
		}
	}));
	const moveMutationState = createMutation<
		OrganizationUnitNode,
		Error,
		MoveOrganizationUnitVariables
	>(() => ({
		mutationKey: ['organization-units', 'move'],
		mutationFn: (variables) => moveOrganizationUnit(fetch, variables),
		onSuccess: async (unit) => {
			selectedOverride = unit.id;
			await queryClient.invalidateQueries({ queryKey: ORGANIZATION_UNITS_QUERY_KEY });
		}
	}));
	const archiveMutationState = createMutation<
		ArchiveOrganizationUnitResponse,
		Error,
		ArchiveOrganizationUnitVariables
	>(() => ({
		mutationKey: ['organization-units', 'archive'],
		mutationFn: (variables) => archiveOrganizationUnit(fetch, variables),
		onSuccess: async () => {
			selectedOverride = null;
			await queryClient.invalidateQueries({ queryKey: ORGANIZATION_UNITS_QUERY_KEY });
		}
	}));

	const units = $derived(organizationUnitsQuery.data?.units ?? []);
	const requestedSelectedId = $derived(selectedOverride ?? pageData.selectedId);
	const flatUnits = $derived(flattenUnitTree(units));
	const selectedId = $derived(
		requestedSelectedId && findUnit(units, requestedSelectedId)
			? requestedSelectedId
			: (selectInitialUnit(units)?.id ?? null)
	);
	const selected = $derived(findUnit(units, selectedId));
	const parentOptions = $derived(selectedId ? collectMoveTargets(units, selectedId) : flatUnits);
	const unitCount = $derived(countUnits(units));
	const isMutationPending = $derived(
		createMutationState.isPending ||
			updateMutationState.isPending ||
			moveMutationState.isPending ||
			archiveMutationState.isPending
	);
	const mutationError = $derived.by(
		() =>
			localError ||
			createMutationState.error?.message ||
			updateMutationState.error?.message ||
			moveMutationState.error?.message ||
			archiveMutationState.error?.message ||
			''
	);

	async function submitRoot(event: SubmitEvent) {
		event.preventDefault();

		try {
			await runCreate({
				parentId: null,
				name: rootName,
				code: rootCode,
				description: ''
			});
			rootName = '';
			rootCode = '';
		} catch {
		}
	}

	async function runCreate(input: OrganizationUnitInput) {
		const normalized = normalizeNamedInput(input);
		await createMutationState.mutateAsync(normalized);
	}

	async function runUpdate(id: string, input: OrganizationUnitInput) {
		const normalized = normalizeNamedInput(input);
		await updateMutationState.mutateAsync({ id, input: normalized });
	}

	async function runMove(id: string, parentId: string | null) {
		resetMutationErrors();
		await moveMutationState.mutateAsync({ id, parentId });
	}

	async function runArchive(id: string) {
		resetMutationErrors();
		await archiveMutationState.mutateAsync({ id });
	}

	function normalizeNamedInput(input: OrganizationUnitInput): OrganizationUnitInput {
		resetMutationErrors();
		if (!input.name.trim()) {
			localError = 'Organization unit name is required';
			throw new Error(localError);
		}

		return {
			parentId: input.parentId ?? null,
			name: input.name,
			code: input.code ?? '',
			description: input.description ?? ''
		};
	}

	function resetMutationErrors() {
		localError = '';
		createMutationState.reset();
		updateMutationState.reset();
		moveMutationState.reset();
		archiveMutationState.reset();
	}
</script>

<svelte:head>
	<title>Organization Units | Syncra DMS</title>
</svelte:head>

<div class="flex flex-1 flex-col gap-4 p-4 lg:p-6">
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
				class="grid w-full gap-2 sm:w-auto sm:grid-cols-[minmax(12rem,1fr)_7rem_auto]"
				onsubmit={submitRoot}
			>
				<label class="sr-only" for="root-unit-name">Root unit name</label>
				<input
					id="root-unit-name"
					class="h-9 rounded-md border bg-background px-3 text-sm"
					bind:value={rootName}
					placeholder="Root name"
					disabled={isMutationPending}
					required
				/>
				<label class="sr-only" for="root-unit-code">Root unit code</label>
				<input
					id="root-unit-code"
					class="h-9 rounded-md border bg-background px-3 text-sm"
					bind:value={rootCode}
					placeholder="Code"
					disabled={isMutationPending}
				/>
				<Button type="submit" size="sm" class="gap-2" disabled={isMutationPending}>
					<PlusIcon class="size-4" />
					New root
				</Button>
			</form>
		{/if}
	</div>

	{#if organizationUnitsQuery.isError}
		<div
			class="flex flex-wrap items-center gap-2 rounded-md border border-destructive/30 bg-destructive/10 p-3 text-sm text-destructive"
			role="alert"
		>
			<AlertCircleIcon class="size-4 shrink-0" />
			<span>{organizationUnitsQuery.error.message}</span>
			<Button
				type="button"
				variant="outline"
				size="xs"
				class="ms-auto"
				onclick={() => organizationUnitsQuery.refetch()}
			>
				Retry
			</Button>
		</div>
	{/if}

	{#if organizationUnitsQuery.isLoading || organizationUnitsQuery.isFetching}
		<div
			class="flex items-center gap-2 rounded-md border bg-muted/30 p-3 text-sm text-muted-foreground"
			role="status"
			aria-live="polite"
		>
			<LoaderCircleIcon class="size-4 shrink-0 animate-spin" />
			<span>Loading organization units</span>
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

	<div class="grid gap-4 lg:grid-cols-[20rem_minmax(0,1fr)]">
		<Card.Root>
			<Card.Header>
				<Card.Title>Active tree</Card.Title>
				<Card.Description>{flatUnits.length} visible</Card.Description>
			</Card.Header>
			<Card.Content>
				<UnitTree units={flatUnits} {selectedId} onSelect={(id) => (selectedOverride = id)} />
			</Card.Content>
		</Card.Root>

		<UnitDetailsPanel
			{selected}
			canManage={pageData.canManageOrganizationUnits}
			{parentOptions}
			isPending={isMutationPending}
			onCreateChild={runCreate}
			onUpdate={runUpdate}
			onMove={runMove}
			onArchive={runArchive}
		/>
	</div>
</div>
