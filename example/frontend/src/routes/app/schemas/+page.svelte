<script lang="ts">
	import ChevronLeftIcon from "@lucide/svelte/icons/chevron-left";
	import ChevronRightIcon from "@lucide/svelte/icons/chevron-right";
	import BookOpenIcon from "@lucide/svelte/icons/book-open";
	import LoaderIcon from "@lucide/svelte/icons/loader-circle";
	import PlusIcon from "@lucide/svelte/icons/plus";
	import SearchIcon from "@lucide/svelte/icons/search";
	import Trash2Icon from "@lucide/svelte/icons/trash-2";
	import { goto } from "$app/navigation";
	import { createMutation, createQuery, useQueryClient } from "@tanstack/svelte-query";
	import { getCoreRowModel, type ColumnDef } from "@tanstack/table-core";

	import * as Alert from "$lib/components/ui/alert/index.js";
	import { Badge } from "$lib/components/ui/badge/index.js";
	import { Button } from "$lib/components/ui/button/index.js";
	import { confirmDelete } from "$lib/components/ui/confirm-delete-dialog/index.js";
	import { FlexRender, renderComponent } from "$lib/components/ui/data-table/index.js";
	import { createSvelteTable } from "$lib/components/ui/data-table/data-table.svelte.js";
	import * as Select from "$lib/components/ui/select/index.js";
	import * as Table from "$lib/components/ui/table/index.js";
	import {
		PERSONAL_SCHEMA_OPTIONS_QUERY_KEY,
		removePersonalSchemaOptions,
		upsertPersonalSchemaOption,
		type PersonalSchemaOption,
	} from "$lib/client/schemas";
	import { m } from "$lib/paraglide/messages.js";
	import { cn } from "$lib/utils.js";
	import ActionsCell from "./actions-cell.svelte";
	import { Spinner } from "$lib/components/ui/spinner/index.js";
	import {
		cloneSchema,
		deleteSchema,
		deleteSchemas,
		fetchSchemas,
		type DeleteSchemasResponse,
		type SchemaListItemResponse,
		type SchemaListResponse,
		type SchemaResponse,
	} from "./api";
	import CreatedDateHeader from "./created-date-header.svelte";
	import NameCell from "./name-cell.svelte";
	import SchemaIdCopy from "./schema-id-copy.svelte";
	import SelectionCheckbox from "./selection-checkbox.svelte";
	import {
		cursorNextState,
		cursorPreviousState,
		formatDate,
		headerSelectionState,
		resetCursorState,
		togglePageSelection,
		toggleSelection,
		type CursorState,
		type SortDirection,
	} from "./table-utils";

	const PAGE_SIZE_OPTIONS = [10, 20, 50, 100];

	type DeleteSchemasVariables = { ids: string[] };
	type CloneSchemaVariables = { schema: SchemaListItemResponse };

	let sortDirection = $state<SortDirection>("desc");
	let pageSize = $state(20);
	let cursorState = $state<CursorState>(resetCursorState());
	let selectedIds = $state<Set<string>>(new Set());

	const queryClient = useQueryClient();

	const schemasQuery = createQuery<SchemaListResponse, Error>(() => ({
		queryKey: [
			"schemas",
			{
				cursor: cursorState.currentCursor,
				size: pageSize,
				sort: sortDirection,
			},
		],
		queryFn: () =>
			fetchSchemas(fetch, {
				cursor: cursorState.currentCursor,
				size: pageSize,
				sort: sortDirection,
			}),
	}));
	const deleteMutation = createMutation<DeleteSchemasResponse, Error, DeleteSchemasVariables>(
		() => ({
			mutationKey: ["schemas", "delete"],
			mutationFn: ({ ids }) => {
				if (ids.length === 1) return deleteSchema(fetch, ids[0]);

				return deleteSchemas(fetch, ids);
			},
			onSuccess: (result, variables) => {
				clearDeletedSchemas(result.deleted_ids, variables.ids);
				queryClient.setQueryData<PersonalSchemaOption[]>(
					PERSONAL_SCHEMA_OPTIONS_QUERY_KEY,
					(current) => removePersonalSchemaOptions(current, result.deleted_ids)
				);
				void queryClient.invalidateQueries({ queryKey: ["schemas"] });
			},
		})
	);
	const cloneMutation = createMutation<SchemaResponse, Error, CloneSchemaVariables>(() => ({
		mutationKey: ["schemas", "clone"],
		mutationFn: ({ schema }) => cloneSchema(fetch, schema),
		onSuccess: (result) => {
			queryClient.setQueryData<SchemaResponse>(["schema", result.id], result);
			queryClient.setQueryData<PersonalSchemaOption[]>(
				PERSONAL_SCHEMA_OPTIONS_QUERY_KEY,
				(current) => upsertPersonalSchemaOption(current, result)
			);
			void queryClient.invalidateQueries({ queryKey: ["schemas"] });
			void goto(`/app/schemas/edit/${result.id}`);
		},
	}));

	const schemas = $derived(schemasQuery.data?.schemas ?? []);
	const nextCursor = $derived(schemasQuery.data?.next_cursor ?? null);
	const visibleIds = $derived(schemas.map((schema) => schema.id));
	const pageSelection = $derived(headerSelectionState(visibleIds, selectedIds));

	function resetPagination() {
		cursorState = resetCursorState();
	}

	function toggleSort() {
		sortDirection = sortDirection === "desc" ? "asc" : "desc";
		resetPagination();
	}

	function setPageSize(value: string) {
		const nextPageSize = Number(value);
		if (!PAGE_SIZE_OPTIONS.includes(nextPageSize)) return;

		pageSize = nextPageSize;
		resetPagination();
	}

	function toggleAllVisible(checked: boolean) {
		selectedIds = togglePageSelection(visibleIds, selectedIds, checked);
	}

	function toggleSchema(id: string, checked: boolean) {
		selectedIds = toggleSelection(selectedIds, id, checked);
	}

	function goNext() {
		cursorState = cursorNextState(cursorState, nextCursor);
	}

	function goPrevious() {
		cursorState = cursorPreviousState(cursorState);
	}

	function clearDeletedSchemas(deletedIds: string[], submittedIds: string[]) {
		const submitted = new Set(submittedIds);
		const deleted = new Set(deletedIds);

		selectedIds = new Set([...selectedIds].filter((id) => !submitted.has(id)));
		for (const id of deleted) {
			queryClient.removeQueries({ queryKey: ["schema", id] });
		}
	}

	async function runDelete(ids: string[]) {
		try {
			await deleteMutation.mutateAsync({ ids });
		} catch {
			// The mutation owns the error state; the page renders it below the toolbar.
		}
	}

	function confirmSingleDelete(schema: SchemaListItemResponse) {
		deleteMutation.reset();
		confirmDelete({
			title: m.schemas_delete_single_title(),
			description: m.schemas_delete_single_description({ name: schema.name }),
			confirm: { text: m.common_delete() },
			onConfirm: () => runDelete([schema.id]),
		});
	}

	function cloneSingleSchema(schema: SchemaListItemResponse) {
		if (cloneMutation.isPending) return;

		cloneMutation.reset();
		cloneMutation.mutate({ schema });
	}

	function confirmBulkDelete() {
		const ids = [...selectedIds];
		if (ids.length === 0) return;

		deleteMutation.reset();
		confirmDelete({
			title:
				ids.length === 1
					? m.schemas_delete_bulk_title_one({ count: ids.length })
					: m.schemas_delete_bulk_title_other({ count: ids.length }),
			description:
				ids.length === 1
					? m.schemas_delete_bulk_description_one({ count: ids.length })
					: m.schemas_delete_bulk_description_other({ count: ids.length }),
			confirm: { text: m.common_delete() },
			onConfirm: () => runDelete(ids),
		});
	}

	const columns: ColumnDef<SchemaListItemResponse>[] = [
		{
			id: "select",
			header: () =>
				renderComponent(SelectionCheckbox, {
					checked: pageSelection.checked,
					indeterminate: pageSelection.indeterminate,
					ariaLabel: m.schemas_select_all_on_page(),
					onCheckedChange: toggleAllVisible,
				}),
			cell: ({ row }) =>
				renderComponent(SelectionCheckbox, {
					checked: selectedIds.has(row.original.id),
					ariaLabel: m.schemas_select_schema({ name: row.original.name }),
					onCheckedChange: (checked) => toggleSchema(row.original.id, checked),
				}),
			enableSorting: false,
			enableHiding: false,
		},
		{
			accessorKey: "name",
			header: m.schemas_name_column(),
			cell: ({ row }) => renderComponent(NameCell, { schema: row.original }),
		},
		{
			accessorKey: "id",
			header: m.schemas_id_column(),
			cell: ({ row }) =>
				renderComponent(SchemaIdCopy, { schemaId: row.original.id, compact: true }),
		},
		{
			accessorKey: "strict",
			header: m.schemas_strict_mode_column(),
			cell: ({ row }) => (row.original.strict ? m.common_strict() : m.common_flexible()),
		},
		{
			accessorKey: "created_at",
			header: () =>
				renderComponent(CreatedDateHeader, {
					sortDirection,
					onToggle: toggleSort,
				}),
			cell: ({ row }) => formatDate(row.original.created_at),
		},
		{
			accessorKey: "updated_at",
			header: m.schemas_updated_column(),
			cell: ({ row }) => formatDate(row.original.updated_at),
		},
		{
			id: "actions",
			header: () => null,
			cell: ({ row }) =>
				renderComponent(ActionsCell, {
					schema: row.original,
					onClone: cloneSingleSchema,
					onDelete: confirmSingleDelete,
					clonePending: cloneMutation.isPending,
					deletePending: deleteMutation.isPending,
				}),
			enableSorting: false,
			enableHiding: false,
		},
	];

	const table = createSvelteTable({
		get data() {
			return schemas;
		},
		columns,
		getRowId: (row) => row.id,
		getCoreRowModel: getCoreRowModel(),
	});
</script>
<svelte:head>
	<title>{m.nav_schemas()} | Syncra</title>
</svelte:head>
<div class="@container/main flex flex-1 flex-col gap-4 p-4 lg:p-6">
	<div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-end">
		<Button
			href="/app/schemas/library"
			variant="outline"
			size="sm"
			class="h-10 px-5 shadow-sm cursor-pointer"
		>
			<BookOpenIcon class="size-4" aria-hidden="true" />
			{m.schemas_library()}
		</Button>

		<Button href="/app/schemas/new" size="sm" class="h-10 px-5 shadow-sm cursor-pointer">
			<PlusIcon class="size-4" aria-hidden="true" />
			{m.schemas_new_schema()}
		</Button>
	</div>

	{#if deleteMutation.isError}
		<Alert.Root variant="destructive" class="mt-2">
			<Alert.Description>{deleteMutation.error.message}</Alert.Description>
		</Alert.Root>
	{/if}

	{#if cloneMutation.isError}
		<Alert.Root variant="destructive" class="mt-2">
			<Alert.Description>{cloneMutation.error.message}</Alert.Description>
		</Alert.Root>
	{/if}

	<div class="overflow-x-auto rounded-xl border bg-background shadow-xs">
		<Table.Root>
			<Table.Header class="sticky top-0 z-10 border-b bg-muted/40">
				{#each table.getHeaderGroups() as headerGroup (headerGroup.id)}
					<Table.Row class="hover:bg-transparent">
						{#each headerGroup.headers as header (header.id)}
							<Table.Head
								colspan={header.colSpan}
								class="h-10 py-2.5 text-xs font-semibold uppercase tracking-wider text-muted-foreground/90"
							>
								{#if !header.isPlaceholder}
									<FlexRender content={header.column.columnDef.header} context={header.getContext()} />
								{/if}
							</Table.Head>
						{/each}
					</Table.Row>
				{/each}
			</Table.Header>
			<Table.Body>
				{#if schemasQuery.isLoading}
					<Table.Row>
						<Table.Cell colspan={columns.length} class="h-56 p-0">
							<div class="flex h-56 w-full items-center justify-center">
								<Spinner class="size-18 text-foreground dark:text-blue-500" />
							</div>
						</Table.Cell>
					</Table.Row>
				{:else if schemasQuery.isError}
					<Table.Row>
						<Table.Cell colspan={columns.length} class="h-24">
							<div
								class="mx-auto flex max-w-xl flex-wrap items-center justify-center gap-3 text-center text-sm text-destructive"
							>
								<span>{schemasQuery.error.message}</span>
								<Button type="button" variant="outline" size="sm" onclick={() => schemasQuery.refetch()}>
									{m.common_retry()}
								</Button>
							</div>
						</Table.Cell>
					</Table.Row>
				{:else if table.getRowModel().rows.length}
					{#each table.getRowModel().rows as row (row.id)}
						<Table.Row
							data-state={selectedIds.has(row.original.id) && "selected"}
							class="transition-colors duration-150 hover:bg-muted/40"
						>
							{#each row.getVisibleCells() as cell (cell.id)}
								<Table.Cell
									class={cn(
										"px-4 py-3",
										cell.column.id === "select" && "w-12",
										cell.column.id === "name" && "min-w-[260px] max-w-[460px]",
										cell.column.id === "id" && "min-w-[220px] max-w-[320px]",
										cell.column.id === "strict" && "w-[140px]",
										cell.column.id === "created_at" && "w-[140px] whitespace-nowrap text-sm",
										cell.column.id === "updated_at" && "w-[140px] whitespace-nowrap text-sm",
										cell.column.id === "actions" && "w-40"
									)}
								>
									{#if cell.column.id === "strict"}
										<Badge
											variant={row.original.strict ? "default" : "secondary"}
											class={cn(!row.original.strict && "text-muted-foreground")}
										>
											{row.original.strict ? m.common_strict() : m.common_flexible()}
										</Badge>
									{:else}
										<FlexRender content={cell.column.columnDef.cell} context={cell.getContext()} />
									{/if}
								</Table.Cell>
							{/each}
						</Table.Row>
					{/each}
					{:else if schemasQuery.isSuccess && schemas.length === 0}
					<Table.Row>
						<Table.Cell colspan={columns.length} class="h-56 text-center">
							<div class="mx-auto flex max-w-[340px] flex-col items-center justify-center gap-2 p-6">
								<div class="rounded-full bg-muted/50 p-3.5 text-muted-foreground/80">
									<SearchIcon class="size-6" aria-hidden="true" />
								</div>
								<h3 class="mt-2 text-sm font-semibold text-foreground">{m.schemas_no_schemas_found()}</h3>
								<p class="px-2 text-xs leading-normal text-muted-foreground">
									{m.schemas_empty_body()}
								</p>
								<Button href="/app/schemas/new" size="sm" class="mt-3.5 h-8 gap-1.5 text-xs font-medium">
									<PlusIcon class="size-3.5" aria-hidden="true" />
									{m.schemas_create_schema()}
								</Button>
							</div>
						</Table.Cell>
					</Table.Row>
				{/if}
			</Table.Body>
		</Table.Root>
	</div>

	<div
		class="flex flex-col gap-3 px-1.5 text-xs text-muted-foreground sm:flex-row sm:items-center sm:justify-between"
	>
		<div>
			{#if schemas.length > 0}
				{schemas.length === 1
					? m.schemas_showing_schemas_one({ count: schemas.length })
					: m.schemas_showing_schemas_other({ count: schemas.length })}
			{:else}
				{m.schemas_no_schemas_to_show()}
			{/if}
		</div>
		<div class="flex flex-wrap items-center gap-3">
			<Select.Root type="single" bind:value={() => String(pageSize), setPageSize}>
				<Select.Trigger size="sm" class="h-8 w-24 bg-background/50 text-xs" aria-label={m.common_rows_per_page()}>
					{pageSize}
				</Select.Trigger>
				<Select.Content side="top">
					{#each PAGE_SIZE_OPTIONS as option (option)}
						<Select.Item value={String(option)} class="text-xs">{option}</Select.Item>
					{/each}
				</Select.Content>
			</Select.Root>

			<div class="flex items-center gap-2">
				<Button
					type="button"
					variant="outline"
					size="sm"
					class="h-8 text-xs"
					onclick={goPrevious}
					disabled={cursorState.history.length === 0 || schemasQuery.isFetching}
				>
					<ChevronLeftIcon class="mr-1 size-4" aria-hidden="true" />
					{m.common_previous()}
				</Button>
				<Button
					type="button"
					variant="outline"
					size="sm"
					class="h-8 text-xs"
					onclick={goNext}
					disabled={!nextCursor || schemasQuery.isFetching}
				>
					{m.common_next()}
					<ChevronRightIcon class="ml-1 size-4" aria-hidden="true" />
				</Button>
			</div>
		</div>
	</div>

	<!-- Floating Bulk Actions Center Bar -->
	{#if selectedIds.size > 0}
		<div class="fixed bottom-6 left-1/2 z-50 -translate-x-1/2 flex items-center gap-3.5 px-4.5 py-2.5 rounded-full border border-border bg-background/95 shadow-xl backdrop-blur-md animate-in fade-in slide-in-from-bottom-4 duration-300">
			<span class="text-xs font-semibold text-foreground px-1 select-none whitespace-nowrap">
				{selectedIds.size === 1
					? m.schemas_selected_count_one({ count: selectedIds.size })
					: m.schemas_selected_count_other({ count: selectedIds.size })}
			</span>
			<div class="h-4 w-[1px] bg-border"></div>
			<div class="flex items-center gap-1.5">
				<Button
					type="button"
					variant="destructive"
					size="sm"
					class="h-8 rounded-full px-4.5 text-xs font-medium shadow-xs"
					disabled={deleteMutation.isPending}
					onclick={confirmBulkDelete}
				>
					{#if deleteMutation.isPending}
						<LoaderIcon class="size-3.5 animate-spin mr-1.5" aria-hidden="true" />
						{m.schemas_deleting()}
					{:else}
						<Trash2Icon class="size-3.5 mr-1.5" aria-hidden="true" />
						{m.common_delete()}
					{/if}
				</Button>
			</div>
		</div>
	{/if}
</div>
