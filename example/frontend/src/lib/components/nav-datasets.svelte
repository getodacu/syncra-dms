<script lang="ts">
	import AlertCircleIcon from "@lucide/svelte/icons/alert-circle";
	import DatabaseIcon from "@lucide/svelte/icons/database";
	import EllipsisIcon from "@lucide/svelte/icons/ellipsis";
	import PencilIcon from "@lucide/svelte/icons/pencil";
	import PlusIcon from "@lucide/svelte/icons/plus";
	import Table2Icon from "@lucide/svelte/icons/table-2";
	import Trash2Icon from "@lucide/svelte/icons/trash-2";
	import { goto } from "$app/navigation";
	import { page } from "$app/state";
	import { createMutation, createQuery, useQueryClient } from "@tanstack/svelte-query";

	import {
		createDataset,
		deleteDataset,
		fetchDatasets,
		updateDataset,
		type CreateDatasetInput,
		type DatasetListResponse,
		type DatasetResponse,
	} from "$lib/client/datasets";
	import { confirmDelete } from "$lib/components/ui/confirm-delete-dialog/index.js";
	import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index.js";
	import * as Sidebar from "$lib/components/ui/sidebar/index.js";
	import { m } from "$lib/paraglide/messages.js";
	import { fetchSchemas, type SchemaListItemResponse } from "../../routes/app/schemas/api";
	import DatasetDialog from "./dataset-dialog.svelte";
	import {
		DATASET_DIALOG_SCHEMA_PAGE_SIZE,
		DATASET_LOADING_ROWS,
		DATASETS_QUERY_KEY,
		closeDatasetDialogState,
		composeDatasetMenuTriggerClick,
		datasetCreateSuccessInvalidationKeys,
		datasetDeleteSuccessEffects,
		datasetDialogError,
		datasetDialogInitialValue,
		datasetDialogPending,
		fetchAllDatasetDialogSchemas,
		datasetHref,
		datasetListOverflows,
		datasetListStatus,
		datasetSubmitAction,
		datasetUpdateSuccessInvalidationKeys,
		isAllDatasetsActive,
		isDatasetActive as isDatasetActiveByRoute,
		openCreateDatasetDialogState,
		openEditDatasetDialogState,
		retryDatasets,
		runDatasetMenuAction,
		type DatasetDialogMode,
		type DatasetDialogState,
		type DatasetQueryKey,
	} from "./nav-datasets-utils";

	type UpdateDatasetVariables = { id: string; input: CreateDatasetInput };
	type DeleteDatasetVariables = { dataset: DatasetResponse };

	const queryClient = useQueryClient();
	const sidebar = Sidebar.useSidebar();

	let dialogOpen = $state(false);
	let dialogMode = $state<DatasetDialogMode>("create");
	let editingDataset = $state<DatasetResponse | null>(null);

	const currentPathname = $derived(page.url.pathname);
	const datasetsQuery = createQuery<DatasetListResponse, Error>(() => ({
		queryKey: DATASETS_QUERY_KEY,
		queryFn: () => fetchDatasets(fetch),
	}));
	const schemasQuery = createQuery<SchemaListItemResponse[], Error>(() => ({
		queryKey: ["schemas", "mine", "datasets", "all", { size: DATASET_DIALOG_SCHEMA_PAGE_SIZE }],
		queryFn: () =>
			fetchAllDatasetDialogSchemas((cursor) =>
				fetchSchemas(fetch, { cursor, size: DATASET_DIALOG_SCHEMA_PAGE_SIZE })
			),
	}));
	const createDatasetMutation = createMutation<DatasetResponse, Error, CreateDatasetInput>(() => ({
		mutationKey: ["datasets", "create"],
		mutationFn: (input) => createDataset(fetch, input),
		onSuccess: () => {
			invalidateQueryKeys(datasetCreateSuccessInvalidationKeys());
			closeDialog();
		},
	}));
	const updateDatasetMutation = createMutation<DatasetResponse, Error, UpdateDatasetVariables>(
		() => ({
			mutationKey: ["datasets", "update"],
			mutationFn: ({ id, input }) => updateDataset(fetch, id, input),
			onSuccess: (_result, variables) => {
				invalidateQueryKeys(
					datasetUpdateSuccessInvalidationKeys({
						pathname: currentPathname,
						updatedDatasetId: variables.id,
					})
				);
				closeDialog();
			},
		})
	);
	const deleteDatasetMutation = createMutation<
		{ deleted_id: string },
		Error,
		DeleteDatasetVariables
	>(() => ({
		mutationKey: ["datasets", "delete"],
		mutationFn: ({ dataset }) => deleteDataset(fetch, dataset.id),
		onSuccess: (result) => {
			const effects = datasetDeleteSuccessEffects({
				pathname: currentPathname,
				deletedDatasetId: result.deleted_id,
			});
			invalidateQueryKeys(effects.invalidateQueryKeys);
			if (effects.navigateTo) void goto(effects.navigateTo);
		},
	}));

	const datasets = $derived(datasetsQuery.data?.datasets ?? []);
	const schemas = $derived(schemasQuery.data ?? []);
	const dialogInitialValue = $derived(datasetDialogInitialValue(editingDataset));
	const dialogPending = $derived(
		datasetDialogPending(createDatasetMutation.isPending, updateDatasetMutation.isPending)
	);
	const dialogError = $derived(
		datasetDialogError(
			dialogMode,
			createDatasetMutation.error ?? null,
			updateDatasetMutation.error ?? null
		)
	);
	const datasetStatus = $derived(
		datasetListStatus({
			isLoading: datasetsQuery.isLoading,
			isError: datasetsQuery.isError,
			datasetCount: datasets.length,
		})
	);
	const datasetsOverflow = $derived(datasetListOverflows(datasets.length));
	const allDatasetsActive = $derived(isAllDatasetsActive(currentPathname));

	function invalidateQueryKeys(queryKeys: DatasetQueryKey[]) {
		for (const queryKey of queryKeys) {
			void queryClient.invalidateQueries({ queryKey: [...queryKey] });
		}
	}

	function applyDialogState(state: DatasetDialogState) {
		dialogOpen = state.open;
		dialogMode = state.mode;
		editingDataset = state.editingDataset;
	}

	function openCreateDialog() {
		createDatasetMutation.reset();
		updateDatasetMutation.reset();
		applyDialogState(openCreateDatasetDialogState());
	}

	function openEditDialog(dataset: DatasetResponse) {
		createDatasetMutation.reset();
		updateDatasetMutation.reset();
		applyDialogState(openEditDatasetDialogState(dataset));
	}

	function closeDialog() {
		applyDialogState(closeDatasetDialogState(dialogMode));
	}

	function submitDataset(value: CreateDatasetInput) {
		const action = datasetSubmitAction(dialogMode, editingDataset, value);
		if (action.type === "create") createDatasetMutation.mutate(action.input);
		if (action.type === "update") {
			updateDatasetMutation.mutate({ id: action.id, input: action.input });
		}
	}

	async function runDelete(dataset: DatasetResponse) {
		try {
			await deleteDatasetMutation.mutateAsync({ dataset });
		} catch {
			// The mutation owns the error state; the menu stays usable after the dialog closes.
		}
	}

	function confirmDatasetDelete(dataset: DatasetResponse) {
		deleteDatasetMutation.reset();
		confirmDelete({
			title: m.datasets_delete_confirm_title(),
			description: m.datasets_delete_confirm_description({ name: dataset.name }),
			confirm: { text: m.datasets_delete() },
			onConfirm: () => runDelete(dataset),
		});
	}

	function isDatasetActive(dataset: DatasetResponse) {
		return isDatasetActiveByRoute(currentPathname, dataset.id);
	}
</script>

<Sidebar.Group class="group-data-[collapsible=icon]:hidden">
	<Sidebar.GroupLabel>{m.datasets_page_title()}</Sidebar.GroupLabel>
	<Sidebar.GroupAction
		type="button"
		aria-label={m.datasets_add_dataset()}
		title={m.datasets_add_dataset()}
		onclick={openCreateDialog}
	>
		<PlusIcon />
		<span class="sr-only">{m.datasets_add_dataset()}</span>
	</Sidebar.GroupAction>
	<Sidebar.GroupContent>
		<Sidebar.Menu>
			<Sidebar.MenuItem>
				<Sidebar.MenuButton isActive={allDatasetsActive} tooltipContent={m.datasets_all_datasets()}>
					{#snippet child({ props })}
						<a {...props} href="/app/datasets">
							<DatabaseIcon />
							<span>{m.datasets_all_datasets()}</span>
						</a>
					{/snippet}
				</Sidebar.MenuButton>
			</Sidebar.MenuItem>
		</Sidebar.Menu>

		<Sidebar.Menu
			aria-label={m.datasets_page_title()}
			class={["mt-1", datasetsOverflow && "max-h-[22.25rem] overflow-y-auto pr-1"]}
		>
			{#if datasetStatus === "loading"}
				{#each DATASET_LOADING_ROWS as row (row)}
					<Sidebar.MenuSkeleton showIcon />
				{/each}
			{:else if datasetStatus === "error"}
				<Sidebar.MenuItem>
					<Sidebar.MenuButton
						class="text-destructive hover:text-destructive"
						onclick={() => void retryDatasets(() => datasetsQuery.refetch())}
					>
						<AlertCircleIcon />
						<span>{m.datasets_retry_datasets()}</span>
					</Sidebar.MenuButton>
				</Sidebar.MenuItem>
			{:else if datasetStatus === "empty"}
				<Sidebar.MenuItem>
					<Sidebar.MenuButton aria-disabled={true} class="text-sidebar-foreground/70">
						<Table2Icon class="text-sidebar-foreground/70" />
						<span>{m.datasets_no_datasets()}</span>
					</Sidebar.MenuButton>
				</Sidebar.MenuItem>
			{:else}
				{#each datasets as dataset (dataset.id)}
					<Sidebar.MenuItem>
						<Sidebar.MenuButton isActive={isDatasetActive(dataset)} tooltipContent={dataset.name}>
							{#snippet child({ props })}
								<a {...props} href={datasetHref(dataset.id)}>
									<Table2Icon />
									<span>{dataset.name}</span>
								</a>
							{/snippet}
						</Sidebar.MenuButton>
						<DropdownMenu.Root>
							<DropdownMenu.Trigger>
								{#snippet child({ props })}
									<Sidebar.MenuAction
										{...props}
										type="button"
										showOnHover
										class="data-[state=open]:bg-accent rounded-sm"
										onclick={composeDatasetMenuTriggerClick(props.onclick)}
									>
										<EllipsisIcon />
										<span class="sr-only">{m.datasets_dataset_actions()}</span>
									</Sidebar.MenuAction>
								{/snippet}
							</DropdownMenu.Trigger>
							<DropdownMenu.Content
								class="w-28 rounded-lg"
								side={sidebar.isMobile ? "bottom" : "right"}
								align={sidebar.isMobile ? "end" : "start"}
							>
								<DropdownMenu.Item
									onclick={(event) => runDatasetMenuAction(event, () => openEditDialog(dataset))}
								>
									<PencilIcon />
									<span>{m.datasets_edit()}</span>
								</DropdownMenu.Item>
								<DropdownMenu.Item
									variant="destructive"
									onclick={(event) =>
										runDatasetMenuAction(event, () => confirmDatasetDelete(dataset))}
								>
									<Trash2Icon />
									<span>{m.datasets_delete()}</span>
								</DropdownMenu.Item>
							</DropdownMenu.Content>
						</DropdownMenu.Root>
					</Sidebar.MenuItem>
				{/each}
			{/if}
		</Sidebar.Menu>

		{#if deleteDatasetMutation.isError}
			<Sidebar.Menu class="mt-1">
				<Sidebar.MenuItem>
					<Sidebar.MenuButton aria-disabled={true} class="text-destructive">
						<AlertCircleIcon />
						<span>{m.datasets_delete_failed()}</span>
					</Sidebar.MenuButton>
				</Sidebar.MenuItem>
			</Sidebar.Menu>
		{/if}
	</Sidebar.GroupContent>
</Sidebar.Group>

<DatasetDialog
	bind:open={dialogOpen}
	mode={dialogMode}
	initialValue={dialogInitialValue}
	{schemas}
	schemasLoading={schemasQuery.isLoading}
	schemasError={schemasQuery.error ?? null}
	pending={dialogPending}
	error={dialogError}
	onSubmit={submitDataset}
/>
