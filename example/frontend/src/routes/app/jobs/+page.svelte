<script lang="ts">
	import ChevronLeftIcon from "@lucide/svelte/icons/chevron-left";
	import ChevronRightIcon from "@lucide/svelte/icons/chevron-right";
	import EyeIcon from "@lucide/svelte/icons/eye";
	import FileIcon from "@lucide/svelte/icons/file";
	import FileImageIcon from "@lucide/svelte/icons/file-image";
	import FileTextIcon from "@lucide/svelte/icons/file-text";
	import LoaderIcon from "@lucide/svelte/icons/loader-circle";
	import SearchIcon from "@lucide/svelte/icons/search";
	import Trash2Icon from "@lucide/svelte/icons/trash-2";
	import PlusIcon from "@lucide/svelte/icons/plus";
	import { createMutation, createQuery, useQueryClient } from "@tanstack/svelte-query";
	import { getCoreRowModel, type ColumnDef } from "@tanstack/table-core";

	import * as Alert from "$lib/components/ui/alert/index.js";
	import { Badge } from "$lib/components/ui/badge/index.js";
	import { Button } from "$lib/components/ui/button/index.js";
	import { confirmDelete } from "$lib/components/ui/confirm-delete-dialog/index.js";
	import { FlexRender, renderComponent } from "$lib/components/ui/data-table/index.js";
	import { createSvelteTable } from "$lib/components/ui/data-table/data-table.svelte.js";
	import * as Dialog from "$lib/components/ui/dialog/index.js";
	import * as Select from "$lib/components/ui/select/index.js";
	import { Spinner } from "$lib/components/ui/spinner/index.js";
	import * as Table from "$lib/components/ui/table/index.js";
	import * as Tooltip from "$lib/components/ui/tooltip/index.js";
	import { CREDIT_BALANCE_QUERY_KEY } from "$lib/client/billing";
	import { m } from "$lib/paraglide/messages.js";
	import { cn } from "$lib/utils.js";
	import { getSchema, type SchemaResponse } from "../schemas/api";
	import SelectionCheckbox from "../documents/selection-checkbox.svelte";
	import {
		deleteOCRJobs,
		fetchOCRJob,
		fetchOCRJobs,
		type DeleteOCRJobsResponse,
		type OCRJobListItemResponse,
		type OCRJobListResponse,
		type OCRJobResponse,
	} from "./api";
	import {
		cursorNextState,
		cursorPreviousState,
		formatCreatedDate,
		formatFileSize,
		headerSelectionState,
		resetCursorState,
		shouldPollJobStatus,
		togglePageSelection,
		toggleSelection,
		type CursorState,
		type SortDirection,
	} from "./table-utils";

	const PAGE_SIZE_OPTIONS = [10, 20, 50, 100];
	const POLL_INTERVAL_MS = 3000;

	type DeleteJobsVariables = { ids: string[] };
	type SchemaDialogSource = OCRJobListItemResponse | null;

	let sortDirection = $state<SortDirection>("desc");
	let pageSize = $state(20);
	let cursorState = $state<CursorState>(resetCursorState());
	let selectedIds = $state<Set<string>>(new Set());
	let schemaDialogOpen = $state(false);
	let schemaDialogJob = $state<SchemaDialogSource>(null);

	const queryClient = useQueryClient();

	const jobsQueryKey = $derived([
		"ocr-jobs",
		{
			cursor: cursorState.currentCursor,
			size: pageSize,
			sort: sortDirection,
		},
	]);
	const jobsQuery = createQuery<OCRJobListResponse, Error>(() => ({
		queryKey: jobsQueryKey,
		queryFn: () =>
			fetchOCRJobs(fetch, {
				cursor: cursorState.currentCursor,
				size: pageSize,
				sort: sortDirection,
			}),
	}));
	const schemaQuery = createQuery<SchemaResponse, Error>(() => ({
		queryKey: ["schema", schemaDialogJob?.schema_id],
		queryFn: () => {
			const schemaId = schemaDialogJob?.schema_id;
			if (!schemaId) throw new Error(m.jobs_missing_schema_id());

			return getSchema(fetch, schemaId);
		},
		enabled: schemaDialogOpen && Boolean(schemaDialogJob?.schema_id),
	}));
	const inlineJobQuery = createQuery<OCRJobResponse, Error>(() => ({
		queryKey: ["ocr-job", schemaDialogJob?.id],
		queryFn: () => {
			const jobId = schemaDialogJob?.id;
			if (!jobId) throw new Error(m.jobs_missing_job_id());

			return fetchOCRJob(fetch, jobId);
		},
		enabled:
			schemaDialogOpen &&
			Boolean(schemaDialogJob?.id) &&
			Boolean(schemaDialogJob?.has_inline_schema) &&
			!schemaDialogJob?.schema_id,
	}));
	const deleteMutation = createMutation<DeleteOCRJobsResponse, Error, DeleteJobsVariables>(() => ({
		mutationKey: ["ocr-jobs", "delete"],
		mutationFn: ({ ids }) => deleteOCRJobs(fetch, ids),
		onSuccess: (result, variables) => {
			clearDeletedJobs(result.deleted_ids, variables.ids);
			void queryClient.invalidateQueries({ queryKey: ["ocr-jobs"] });
		},
	}));

	const jobs = $derived(jobsQuery.data?.jobs ?? []);
	const nextCursor = $derived(jobsQuery.data?.next_cursor ?? null);
	const visibleIds = $derived(jobs.map((job) => job.id));
	const pageSelection = $derived(headerSelectionState(visibleIds, selectedIds));
	const activeJobs = $derived(jobs.filter((job) => shouldPollJobStatus(job.status)));
	const selectedCount = $derived(selectedIds.size);

	$effect(() => {
		const ids = activeJobs.map((job) => job.id).join(",");
		if (!ids) return;

		const interval = window.setInterval(() => {
			void pollVisibleJobs();
		}, POLL_INTERVAL_MS);

		return () => window.clearInterval(interval);
	});

	function resetPagination() {
		cursorState = resetCursorState();
		selectedIds = new Set();
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

	function toggleJob(id: string, checked: boolean) {
		selectedIds = toggleSelection(selectedIds, id, checked);
	}

	function goNext() {
		cursorState = cursorNextState(cursorState, nextCursor);
		selectedIds = new Set();
	}

	function goPrevious() {
		cursorState = cursorPreviousState(cursorState);
		selectedIds = new Set();
	}

	function openSchemaDialog(job: OCRJobListItemResponse) {
		schemaDialogJob = job;
		schemaDialogOpen = true;
	}

	function clearDeletedJobs(deletedIds: string[], submittedIds: string[]) {
		const submitted = new Set(submittedIds);
		const deleted = new Set(deletedIds);

		selectedIds = new Set([...selectedIds].filter((id) => !submitted.has(id)));
		for (const id of deleted) {
			queryClient.removeQueries({ queryKey: ["ocr-job", id] });
		}

		if (schemaDialogJob && deleted.has(schemaDialogJob.id)) {
			schemaDialogOpen = false;
			schemaDialogJob = null;
		}
	}

	async function runDelete(ids: string[]) {
		try {
			await deleteMutation.mutateAsync({ ids });
		} catch {
			// The mutation owns the error state; the page renders it below the toolbar.
		}
	}

	function confirmBulkDelete() {
		const ids = [...selectedIds];
		if (ids.length === 0) return;

		deleteMutation.reset();
		confirmDelete({
			title:
				ids.length === 1
					? m.jobs_delete_bulk_title_one({ count: ids.length })
					: m.jobs_delete_bulk_title_other({ count: ids.length }),
			description:
				ids.length === 1
					? m.jobs_delete_bulk_description_one({ count: ids.length })
					: m.jobs_delete_bulk_description_other({ count: ids.length }),
			confirm: { text: m.common_delete() },
			onConfirm: () => runDelete(ids),
		});
	}

	function confirmSingleDelete(job: OCRJobListItemResponse) {
		deleteMutation.reset();
		confirmDelete({
			title: m.jobs_delete_single_title(),
			description: m.jobs_delete_single_description({ name: job.original_filename }),
			confirm: { text: m.common_delete() },
			onConfirm: () => runDelete([job.id]),
		});
	}

	async function pollVisibleJobs() {
		const currentJobs = activeJobs;
		if (currentJobs.length === 0) return;
		let completedJob = false;

		await Promise.all(
			currentJobs.map(async (job) => {
				try {
					const latest = await fetchOCRJob(fetch, job.id);
					if (shouldPollJobStatus(job.status) && latest.status === "completed") {
						completedJob = true;
					}
					queryClient.setQueryData<OCRJobListResponse>(jobsQueryKey, (current) => {
						if (!current) return current;

						return {
							...current,
							jobs: current.jobs.map((item) =>
								item.id === latest.id ? toListItem(latest) : item
							),
						};
					});
					queryClient.setQueryData<OCRJobResponse>(["ocr-job", latest.id], latest);
					if (shouldPollJobStatus(job.status) && !shouldPollJobStatus(latest.status)) {
						void queryClient.invalidateQueries({ queryKey: ["ocr-jobs"] });
					}
				} catch {
					// Keep the current row visible; the next list refetch owns surfaced errors.
				}
			})
		);

		if (completedJob) {
			void queryClient.invalidateQueries({ queryKey: CREDIT_BALANCE_QUERY_KEY });
		}
	}

	function toListItem(job: OCRJobResponse): OCRJobListItemResponse {
		const { inline_schema, ...listItem } = job;
		void inline_schema;
		return listItem;
	}

	function fileIconKind(mimeType: string) {
		const normalizedMimeType = mimeType.toLowerCase();
		if (normalizedMimeType === "application/pdf") return "pdf";
		if (normalizedMimeType.startsWith("image/")) return "image";
		return "file";
	}

	function truncateFilename(filename: string) {
		return filename.length > 60 ? `${filename.slice(0, 60)}...` : filename;
	}

	function statusLabel(status: string) {
		if (status === "queued") return m.jobs_status_queued();
		if (status === "pending") return m.jobs_status_pending();
		if (status === "processing") return m.jobs_status_processing();
		if (status === "completed") return m.jobs_status_completed();
		if (status === "failed") return m.jobs_status_failed();
		return status || m.common_unknown();
	}

	function statusClass(status: string) {
		if (status === "completed") {
			return "bg-emerald-500/10 text-emerald-700 border-emerald-500/20 dark:text-emerald-400";
		}
		if (status === "failed") {
			return "bg-destructive/10 text-destructive border-destructive/20";
		}
		if (status === "processing") {
			return "bg-blue-500/10 text-blue-700 border-blue-500/20 dark:text-blue-400";
		}
		return "bg-amber-500/10 text-amber-700 border-amber-500/20 dark:text-amber-400";
	}

	function failedDescription(job: OCRJobListItemResponse) {
		if (job.status !== "failed") return "";
		return job.error_message?.trim() ?? "";
	}

	function schemaButtonLabel(job: OCRJobListItemResponse) {
		if (job.schema_name) return job.schema_name;
		if (job.has_inline_schema) return m.jobs_inline_schema();
		return null;
	}

	function schemaJSON() {
		if (schemaDialogJob?.schema_id) return schemaQuery.data?.schema;
		return inlineJobQuery.data?.inline_schema;
	}

	function schemaModalTitle() {
		if (schemaQuery.data?.name) return schemaQuery.data.name;
		if (schemaDialogJob?.schema_name) return schemaDialogJob.schema_name;
		if (schemaDialogJob?.has_inline_schema) return m.jobs_inline_schema();
		return m.jobs_schema();
	}

	function formatSchemaJSON(value: unknown) {
		if (value === undefined) return "";
		return JSON.stringify(value, null, 2);
	}

	const columns: ColumnDef<OCRJobListItemResponse>[] = [
		{
			id: "select",
			header: () =>
				renderComponent(SelectionCheckbox, {
					checked: pageSelection.checked,
					indeterminate: pageSelection.indeterminate,
					ariaLabel: m.jobs_select_all_on_page(),
					onCheckedChange: toggleAllVisible,
				}),
			cell: ({ row }) =>
				renderComponent(SelectionCheckbox, {
					checked: selectedIds.has(row.original.id),
					ariaLabel: m.jobs_select_job({ name: row.original.original_filename }),
					onCheckedChange: (checked) => toggleJob(row.original.id, checked),
				}),
			enableSorting: false,
			enableHiding: false,
		},
		{
			accessorKey: "original_filename",
			header: m.jobs_filename_column(),
		},
		{
			accessorKey: "status",
			header: m.jobs_status_column(),
		},
		{
			accessorKey: "created_at",
			header: m.jobs_created_column(),
			cell: ({ row }) => formatCreatedDate(row.original.created_at),
		},
		{
			accessorKey: "file_size",
			header: m.jobs_file_size_column(),
			cell: ({ row }) => formatFileSize(row.original.file_size),
		},
		{
			accessorKey: "page_count",
			header: m.jobs_pages_column(),
		},
		{
			id: "schema",
			header: m.jobs_schema(),
		},
		{
			id: "actions",
			header: () => null,
			enableSorting: false,
			enableHiding: false,
		},
	];

	const table = createSvelteTable({
		get data() {
			return jobs;
		},
		columns,
		getRowId: (row) => row.id,
		getCoreRowModel: getCoreRowModel(),
	});
</script>

<svelte:head>
	<title>{m.jobs_page_title()} | Syncra</title>
</svelte:head>

<div class="@container/main flex flex-1 flex-col gap-4 p-4 lg:p-6">
	<div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-end">
		
		<Button href="/app/new-job" size="sm" class="h-10 px-5 shadow-sm cursor-pointer">
			<PlusIcon class="size-4" aria-hidden="true" />
			{m.jobs_new_job()}
		</Button>
	</div>

	{#if deleteMutation.isError}
		<Alert.Root variant="destructive">
			<Alert.Description>{deleteMutation.error.message}</Alert.Description>
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
									<FlexRender
										content={header.column.columnDef.header}
										context={header.getContext()}
									/>
								{/if}
							</Table.Head>
						{/each}
					</Table.Row>
				{/each}
			</Table.Header>
			<Table.Body>
				{#if jobsQuery.isLoading}
					<Table.Row>
						<Table.Cell colspan={columns.length} class="h-56 p-0">
							<div class="flex h-56 w-full items-center justify-center">
								<Spinner class="size-18 text-foreground dark:text-blue-500" />
							</div>
						</Table.Cell>
					</Table.Row>
				{:else if jobsQuery.isError}
					<Table.Row>
						<Table.Cell colspan={columns.length} class="h-24">
							<div class="mx-auto flex max-w-xl flex-wrap items-center justify-center gap-3 text-center text-sm text-destructive">
								<span>{jobsQuery.error.message}</span>
								<Button type="button" variant="outline" size="sm" onclick={() => jobsQuery.refetch()}>
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
										cell.column.id === "select" && "w-10",
										cell.column.id === "original_filename" && "min-w-[260px]",
										cell.column.id === "status" && "w-[150px]",
										cell.column.id === "created_at" && "w-[140px] whitespace-nowrap text-sm",
										cell.column.id === "file_size" && "w-[110px] whitespace-nowrap text-sm",
										cell.column.id === "page_count" && "w-[80px] text-sm",
										cell.column.id === "schema" && "min-w-[150px]",
										cell.column.id === "actions" && "w-[92px]"
									)}
								>
									{#if cell.column.id === "original_filename"}
										{@const iconKind = fileIconKind(row.original.mime_type)}
										<div class="flex min-w-0 items-center gap-2" title={row.original.original_filename}>
											{#if iconKind === "pdf"}
												<FileTextIcon class="size-4 shrink-0 text-rose-500 dark:text-rose-400" aria-hidden="true" />
											{:else if iconKind === "image"}
												<FileImageIcon class="size-4 shrink-0 text-blue-500 dark:text-blue-400" aria-hidden="true" />
											{:else}
												<FileIcon class="size-4 shrink-0 text-muted-foreground" aria-hidden="true" />
											{/if}
											<span class="truncate text-sm font-medium">{truncateFilename(row.original.original_filename)}</span>
										</div>
									{:else if cell.column.id === "status"}
										{@const failureText = failedDescription(row.original)}
										{#if failureText}
											<Tooltip.Root>
												<Tooltip.Trigger>
													{#snippet child({ props })}
														<span {...props} class="inline-flex">
															<Badge variant="outline" class={cn("rounded-sm border text-xs", statusClass(row.original.status))}>
																{statusLabel(row.original.status)}
															</Badge>
														</span>
													{/snippet}
												</Tooltip.Trigger>
												<Tooltip.Content class="max-w-[260px] whitespace-normal leading-relaxed">
													{failureText}
												</Tooltip.Content>
											</Tooltip.Root>
										{:else}
											<Badge variant="outline" class={cn("rounded-sm border text-xs", statusClass(row.original.status))}>
												{#if shouldPollJobStatus(row.original.status)}
													<LoaderIcon class="size-3 animate-spin" aria-hidden="true" />
												{/if}
												{statusLabel(row.original.status)}
											</Badge>
										{/if}
									{:else if cell.column.id === "schema"}
										{@const schemaLabel = schemaButtonLabel(row.original)}
										{#if schemaLabel}
											<Button
												type="button"
												variant="ghost"
												size="sm"
												class="h-7 max-w-[180px] justify-start px-2 text-xs"
												onclick={() => openSchemaDialog(row.original)}
											>
												<span class="truncate">{schemaLabel}</span>
											</Button>
										{:else}
											<span class="text-xs text-muted-foreground">{m.jobs_no_schema()}</span>
										{/if}
									{:else if cell.column.id === "actions"}
										<div class="flex items-center justify-end gap-1">
											
											<Button
												type="button"
												variant="ghost"
												size="icon-sm"
												class="text-muted-foreground transition-all hover:bg-destructive/10 hover:text-destructive dark:hover:bg-destructive/20"
												disabled={deleteMutation.isPending}
												aria-label={m.jobs_delete_job({ name: row.original.original_filename })}
												onclick={() => confirmSingleDelete(row.original)}
											>
												<Trash2Icon class="size-4" aria-hidden="true" />
											</Button>
										</div>
									{:else}
										<FlexRender content={cell.column.columnDef.cell} context={cell.getContext()} />
									{/if}
								</Table.Cell>
							{/each}
						</Table.Row>
					{/each}
				{:else if jobsQuery.isSuccess && jobs.length === 0}
					<Table.Row>
						<Table.Cell colspan={columns.length} class="h-56 text-center">
							<div class="mx-auto flex max-w-[340px] flex-col items-center justify-center gap-2 p-6">
								<div class="rounded-full bg-muted/50 p-3.5 text-muted-foreground/80">
									<SearchIcon class="size-6" aria-hidden="true" />
								</div>
								<h3 class="mt-2 text-sm font-semibold text-foreground">{m.jobs_no_jobs_found()}</h3>
								<p class="px-2 text-xs leading-normal text-muted-foreground">
									{m.jobs_empty_body()}
								</p>
								<Button href="/app/new-job" size="sm" class="mt-3.5 h-8 text-xs font-medium">
									{m.jobs_new_job()}
								</Button>
							</div>
						</Table.Cell>
					</Table.Row>
				{/if}
			</Table.Body>
		</Table.Root>
	</div>

	<div class="flex flex-col gap-3 px-1.5 text-xs text-muted-foreground sm:flex-row sm:items-center sm:justify-between">
		<div>
			{#if jobs.length > 0}
				{jobs.length === 1
					? m.jobs_showing_jobs_one({ count: jobs.length })
					: m.jobs_showing_jobs_other({ count: jobs.length })}
			{:else}
				{m.jobs_no_jobs_to_show()}
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
			<div class="flex items-center gap-1">
				<Button
					type="button"
					variant="outline"
					size="icon-sm"
					disabled={cursorState.history.length === 0 || jobsQuery.isFetching}
					onclick={goPrevious}
					aria-label={m.common_previous()}
				>
					<ChevronLeftIcon class="size-4" aria-hidden="true" />
				</Button>
				<Button
					type="button"
					variant="outline"
					size="icon-sm"
					disabled={!nextCursor || jobsQuery.isFetching}
					onclick={goNext}
					aria-label={m.common_next()}
				>
					<ChevronRightIcon class="size-4" aria-hidden="true" />
				</Button>
			</div>
		</div>
	</div>

	<!-- Floating Bulk Actions Center Bar -->
	{#if selectedIds.size > 0}
		<div class="fixed bottom-6 left-1/2 z-50 -translate-x-1/2 flex items-center gap-3.5 px-4.5 py-2.5 rounded-full border border-border bg-background/95 shadow-xl backdrop-blur-md animate-in fade-in slide-in-from-bottom-4 duration-300">
			<span class="text-xs font-semibold text-foreground px-1 select-none whitespace-nowrap">
				{selectedIds.size === 1
					? m.jobs_selected_count_one({ count: selectedIds.size })
					: m.jobs_selected_count_other({ count: selectedIds.size })}
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
						{m.jobs_deleting()}
					{:else}
						<Trash2Icon class="size-3.5 mr-1.5" aria-hidden="true" />
						{m.common_delete()}
					{/if}
				</Button>
			</div>
		</div>
	{/if}
</div>

<Dialog.Root bind:open={schemaDialogOpen}>
	<Dialog.Content class="sm:max-w-4xl">
		<Dialog.Header>
			<Dialog.Title>{schemaModalTitle()}</Dialog.Title>
			<Dialog.Description>
				{#if schemaQuery.data}
					{schemaQuery.data.description || m.jobs_saved_extraction_schema()}
				{:else if schemaDialogJob?.has_inline_schema}
					{m.jobs_inline_schema_description()}
				{:else}
					{m.jobs_extraction_schema_details()}
				{/if}
			</Dialog.Description>
		</Dialog.Header>
		{#if schemaDialogJob?.schema_id && schemaQuery.isLoading}
			<div class="flex h-40 items-center justify-center">
				<Spinner class="size-8" />
			</div>
		{:else if !schemaDialogJob?.schema_id && inlineJobQuery.isLoading}
			<div class="flex h-40 items-center justify-center">
				<Spinner class="size-8" />
			</div>
		{:else if schemaQuery.isError}
			<Alert.Root variant="destructive">
				<Alert.Description>{schemaQuery.error.message}</Alert.Description>
			</Alert.Root>
		{:else if inlineJobQuery.isError}
			<Alert.Root variant="destructive">
				<Alert.Description>{inlineJobQuery.error.message}</Alert.Description>
			</Alert.Root>
		{:else}
			{#if schemaQuery.data}
				<div class="flex flex-wrap items-center gap-2 text-xs text-muted-foreground">
					<Badge variant="outline" class="rounded-sm">
						{schemaQuery.data.strict ? m.common_strict() : m.common_flexible()}
					</Badge>
				</div>
			{/if}
			<pre class="max-h-[55vh] overflow-auto rounded-md border bg-muted/30 p-3 text-xs leading-relaxed"><code>{formatSchemaJSON(schemaJSON())}</code></pre>
		{/if}
	</Dialog.Content>
</Dialog.Root>
