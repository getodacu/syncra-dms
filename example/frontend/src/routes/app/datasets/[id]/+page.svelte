<script lang="ts">
	import ArrowLeftIcon from "@lucide/svelte/icons/arrow-left";
	import ArrowUpDownIcon from "@lucide/svelte/icons/arrow-up-down";
	import CalendarIcon from "@lucide/svelte/icons/calendar";
	import ChevronLeftIcon from "@lucide/svelte/icons/chevron-left";
	import ChevronRightIcon from "@lucide/svelte/icons/chevron-right";
	import DownloadIcon from "@lucide/svelte/icons/download";
	import LoaderIcon from "@lucide/svelte/icons/loader-circle";
	import SearchIcon from "@lucide/svelte/icons/search";
	import { page } from "$app/state";
	import {
		endOfMonth,
		endOfWeek,
		getLocalTimeZone,
		startOfMonth,
		startOfWeek,
		today,
	} from "@internationalized/date";
	import { createQuery } from "@tanstack/svelte-query";
	import type { ComponentProps } from "svelte";
	import { toast } from "svelte-sonner";

	import { isDatasetNotFoundError } from "$lib/client/datasets";
	import * as Alert from "$lib/components/ui/alert/index.js";
	import { Button } from "$lib/components/ui/button/index.js";
	import { DocumentPreviewDialog } from "$lib/components/ui/document-preview/index.js";
	import * as Popover from "$lib/components/ui/popover/index.js";
	import * as RangeCalendar from "$lib/components/ui/range-calendar/index.js";
	import * as Select from "$lib/components/ui/select/index.js";
	import { Spinner } from "$lib/components/ui/spinner/index.js";
	import * as Table from "$lib/components/ui/table/index.js";
	import { m } from "$lib/paraglide/messages.js";
	import { fetchOCRDocumentPreview, type OCRDocumentPreview } from "../../documents/api";
	import {
		downloadDatasetExport,
		fetchDatasetRows,
		type DatasetExportFormat,
		type DatasetExportResponse,
		type DatasetRowsResponse,
	} from "../api";
	import {
		cursorNextState,
		cursorPreviousState,
		datasetCellText,
		datasetExportFilename,
		dateRangeToQueryBounds,
		formatDatasetDate,
		resetCursorState,
		shouldRetryDatasetRowsQuery,
		type CursorState,
		type DateRangeValue,
		type SortDirection,
	} from "../table-utils";

	const PAGE_SIZE_OPTIONS = [10, 20, 50, 100];
	const BASE_ROW_COLUMN_COUNT = 3;

	type RangeCalendarValue = ComponentProps<typeof RangeCalendar.RangeCalendar>["value"];
	type DatasetCursorState = CursorState & { datasetId: string };
	type DatasetExportRequestState = {
		requestId: number;
		datasetId: string;
		format: DatasetExportFormat;
	};
	type DatasetExportPendingState = DatasetExportRequestState;

	const datasetId = $derived(page.params.id ?? "");
	let nextExportRequestId = 0;
	let sortDirection = $state<SortDirection>("desc");
	let pageSize = $state(20);
	let cursorState = $state<DatasetCursorState>(cursorStateForCurrentDataset());
	let appliedDateRange = $state<DateRangeValue | undefined>();
	let pendingDateRange = $state<RangeCalendarValue>();
	let datePopoverOpen = $state(false);
	let exportPendingState = $state<DatasetExportPendingState | null>(null);
	let previewOpen = $state(false);
	let previewDocumentId = $state<string | null>(null);
	let previewFilename = $state<string | null>(null);
	const activeCursorState = $derived(
		cursorState.datasetId === datasetId ? cursorState : cursorStateForCurrentDataset()
	);
	const effectiveCursor = $derived(activeCursorState.currentCursor);
	const dateBounds = $derived(dateRangeToQueryBounds(appliedDateRange));
	const activeExportFormatPending = $derived(
		exportPendingState?.datasetId === datasetId ? exportPendingState.format : null
	);
	const dateRangeLabel = $derived.by(() => {
		if (!appliedDateRange?.start && !appliedDateRange?.end) return m.datasets_date_range();

		const start = appliedDateRange.start ? appliedDateRange.start.toString() : m.datasets_any_date();
		const end = appliedDateRange.end ? appliedDateRange.end.toString() : m.datasets_any_date();

		return m.datasets_date_range_value({ start, end });
	});
	const activePreset = $derived.by(() => {
		if (!appliedDateRange?.start || !appliedDateRange?.end) return null;

		const tz = getLocalTimeZone();
		const t = today(tz);

		const todayStr = t.toString();
		if (appliedDateRange.start.toString() === todayStr && appliedDateRange.end.toString() === todayStr) {
			return "today";
		}

		const weekStart = startOfWeek(t, "en-US").toString();
		const weekEnd = endOfWeek(t, "en-US").toString();
		if (appliedDateRange.start.toString() === weekStart && appliedDateRange.end.toString() === weekEnd) {
			return "week";
		}

		const monthStart = startOfMonth(t).toString();
		const monthEnd = endOfMonth(t).toString();
		if (appliedDateRange.start.toString() === monthStart && appliedDateRange.end.toString() === monthEnd) {
			return "month";
		}

		return null;
	});

	const rowsQuery = createQuery<DatasetRowsResponse, Error>(() => ({
		queryKey: [
			"dataset-rows",
			datasetId,
			{
				createdFrom: dateBounds.createdFrom,
				createdTo: dateBounds.createdTo,
				cursor: effectiveCursor,
				size: pageSize,
				sort: sortDirection,
			},
		],
		queryFn: () =>
			fetchDatasetRows(fetch, datasetId, {
				createdFrom: dateBounds.createdFrom,
				createdTo: dateBounds.createdTo,
				cursor: effectiveCursor,
				size: pageSize,
				sort: sortDirection,
			}),
		enabled: Boolean(datasetId),
		retry: shouldRetryDatasetRowsQuery,
	}));
	const previewQuery = createQuery<OCRDocumentPreview, Error>(() => ({
		queryKey: ["ocr-document", previewDocumentId],
		queryFn: () => {
			const id = previewDocumentId;
			if (!id) throw new Error(m.datasets_missing_document_id());

			return fetchOCRDocumentPreview(fetch, id);
		},
		enabled: previewOpen && previewDocumentId !== null,
		staleTime: 1000 * 60 * 2,
	}));

	const dataset = $derived(rowsQuery.data?.dataset ?? null);
	const rowColumns = $derived(rowsQuery.data?.columns ?? []);
	const rows = $derived(rowsQuery.data?.rows ?? []);
	const nextCursor = $derived(rowsQuery.data?.next_cursor ?? null);
	const rowColumnCount = $derived(BASE_ROW_COLUMN_COUNT + rowColumns.length);
	const datasetNotFound = $derived(isDatasetNotFoundError(rowsQuery.error));

	function resetPagination() {
		cursorState = cursorStateForCurrentDataset();
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

	function rangeCalendarValue(range?: DateRangeValue): RangeCalendarValue {
		if (!range?.start && !range?.end) return undefined;

		return {
			start: range.start,
			end: range.end,
		};
	}

	function setDatePopoverOpen(open: boolean) {
		datePopoverOpen = open;
		if (open) pendingDateRange = rangeCalendarValue(appliedDateRange);
	}

	function setTodayPreset() {
		const tz = getLocalTimeZone();
		const t = today(tz);
		pendingDateRange = {
			start: t,
			end: t,
		};
		appliedDateRange = pendingDateRange;
		datePopoverOpen = false;
		resetPagination();
	}

	function setThisWeekPreset() {
		const tz = getLocalTimeZone();
		const t = today(tz);
		pendingDateRange = {
			start: startOfWeek(t, "en-US"),
			end: endOfWeek(t, "en-US"),
		};
		appliedDateRange = pendingDateRange;
		datePopoverOpen = false;
		resetPagination();
	}

	function setThisMonthPreset() {
		const tz = getLocalTimeZone();
		const t = today(tz);
		pendingDateRange = {
			start: startOfMonth(t),
			end: endOfMonth(t),
		};
		appliedDateRange = pendingDateRange;
		datePopoverOpen = false;
		resetPagination();
	}

	function applyDateRange() {
		appliedDateRange = pendingDateRange;
		datePopoverOpen = false;
		resetPagination();
	}

	function clearDateRange() {
		pendingDateRange = undefined;
		appliedDateRange = undefined;
		datePopoverOpen = false;
		resetPagination();
	}

	function goNext() {
		cursorState = cursorStateForCurrentDataset(cursorNextState(activeCursorState, nextCursor));
	}

	function goPrevious() {
		cursorState = cursorStateForCurrentDataset(cursorPreviousState(activeCursorState));
	}

	function openDocumentPreview(documentId: string, filename: string) {
		previewDocumentId = documentId;
		previewFilename = filename;
		previewOpen = true;
	}

	async function downloadDataset(format: DatasetExportFormat) {
		const currentDataset = dataset;
		const currentDatasetId = datasetId;
		if (!currentDataset || activeExportFormatPending) return;

		const currentExportRequest = {
			requestId: ++nextExportRequestId,
			datasetId: currentDatasetId,
			format,
		};

		exportPendingState = currentExportRequest;

		try {
			const result = await downloadDatasetExport(fetch, currentDatasetId, format, {
				createdFrom: dateBounds.createdFrom,
				createdTo: dateBounds.createdTo,
				sort: sortDirection,
			});
			triggerBrowserDownload(result, currentDataset.name, format);
		} catch (error) {
			if (isActiveExportRequest(currentExportRequest)) {
				toast.error(error instanceof Error ? error.message : m.datasets_failed_export());
			}
		} finally {
			if (isSameExportRequest(exportPendingState, currentExportRequest)) {
				exportPendingState = null;
			}
		}
	}

	function cursorStateForCurrentDataset(state: CursorState = resetCursorState()): DatasetCursorState {
		return {
			datasetId,
			currentCursor: state.currentCursor,
			history: state.history,
		};
	}

	function isSameExportRequest(
		state: DatasetExportRequestState | null,
		request: DatasetExportRequestState
	) {
		return (
			state?.requestId === request.requestId &&
			state.datasetId === request.datasetId &&
			state.format === request.format
		);
	}

	function isActiveExportRequest(request: DatasetExportRequestState) {
		return datasetId === request.datasetId && isSameExportRequest(exportPendingState, request);
	}

	function triggerBrowserDownload(
		result: DatasetExportResponse,
		datasetName: string,
		format: DatasetExportFormat
	) {
		const filename = result.filename || datasetExportFilename(datasetName, format, new Date());
		const url = URL.createObjectURL(result.blob);
		const link = document.createElement("a");
		link.href = url;
		link.download = filename;
		document.body.append(link);
		link.click();
		link.remove();

		setTimeout(() => URL.revokeObjectURL(url), 1000);
	}

	function rowCountSummary(count: number) {
		return count === 1
			? m.datasets_showing_rows_one({ count })
			: m.datasets_showing_rows_other({ count });
	}

	function sortCreatedDateLabel() {
		return sortDirection === "desc"
			? m.datasets_sort_created_ascending()
			: m.datasets_sort_created_descending();
	}
</script>

<svelte:head>
	<title>{m.datasets_detail_page_title()} | Syncra</title>
</svelte:head>

<div class="@container/main flex flex-1 flex-col gap-4 p-4 lg:p-6">
	<div class="flex flex-col gap-2 border-b pb-5">
		<div class="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
			<div class="min-w-0 space-y-1">
				<Button href="/app/datasets" variant="ghost" size="sm" class="-ml-2 h-8 px-2 text-xs">
					<ArrowLeftIcon class="size-3.5" aria-hidden="true" />
					{m.datasets_page_title()}
				</Button>
				<h1 class="truncate text-2xl font-semibold tracking-tight">
					{dataset?.name ?? m.datasets_detail_page_title()}
				</h1>
			</div>
			<div class="mt-1 flex flex-wrap items-center gap-2 sm:mt-0">
				<Button
					type="button"
					variant="outline"
					size="sm"
					class="h-9"
					disabled={!dataset || Boolean(activeExportFormatPending)}
					onclick={() => downloadDataset("csv")}
				>
					{#if activeExportFormatPending === "csv"}
						<LoaderIcon class="size-4 animate-spin" aria-hidden="true" />
					{:else}
						<DownloadIcon class="size-4" aria-hidden="true" />
					{/if}
					{m.datasets_export_csv()}
				</Button>
				<Button
					type="button"
					variant="outline"
					size="sm"
					class="h-9"
					disabled={!dataset || Boolean(activeExportFormatPending)}
					onclick={() => downloadDataset("xlsx")}
				>
					{#if activeExportFormatPending === "xlsx"}
						<LoaderIcon class="size-4 animate-spin" aria-hidden="true" />
					{:else}
						<DownloadIcon class="size-4" aria-hidden="true" />
					{/if}
					{m.datasets_export_xlsx()}
				</Button>
			</div>
		</div>
	</div>

	<div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
		<div class="flex min-w-0 flex-1 flex-wrap items-center gap-2">
			<Popover.Root bind:open={() => datePopoverOpen, setDatePopoverOpen}>
				<Popover.Trigger>
					{#snippet child({ props })}
						<Button
							type="button"
							variant="outline"
							class="h-9 w-full min-w-[160px] justify-start bg-background/50 text-xs sm:w-auto"
							{...props}
						>
							<CalendarIcon class="mr-2 size-4 text-muted-foreground" aria-hidden="true" />
							<span class="truncate">{dateRangeLabel}</span>
						</Button>
					{/snippet}
				</Popover.Trigger>
				<Popover.Content align="start" class="w-auto p-0">
					<div class="flex flex-col sm:flex-row">
						<div
							class="flex min-w-[130px] flex-row gap-1.5 border-b border-border bg-muted/5 p-3 sm:flex-col sm:border-r sm:border-b-0"
						>
							<span
								class="hidden select-none px-2 py-1 text-[10px] font-bold uppercase tracking-wider text-muted-foreground/60 sm:inline-block"
							>
								{m.datasets_presets()}
							</span>
							<Button
								type="button"
								variant={activePreset === "today" ? "secondary" : "ghost"}
								class="h-8.5 flex-1 justify-center rounded-md px-3 text-xs font-medium transition-all duration-150 active:scale-[0.98] sm:flex-none sm:justify-start"
								onclick={setTodayPreset}
							>
								{m.datasets_today()}
							</Button>
							<Button
								type="button"
								variant={activePreset === "week" ? "secondary" : "ghost"}
								class="h-8.5 flex-1 justify-center rounded-md px-3 text-xs font-medium transition-all duration-150 active:scale-[0.98] sm:flex-none sm:justify-start"
								onclick={setThisWeekPreset}
							>
								{m.datasets_this_week()}
							</Button>
							<Button
								type="button"
								variant={activePreset === "month" ? "secondary" : "ghost"}
								class="h-8.5 flex-1 justify-center rounded-md px-3 text-xs font-medium transition-all duration-150 active:scale-[0.98] sm:flex-none sm:justify-start"
								onclick={setThisMonthPreset}
							>
								{m.datasets_this_month()}
							</Button>
						</div>
						<div class="flex flex-col">
							<RangeCalendar.RangeCalendar bind:value={pendingDateRange} numberOfMonths={2} />
							<div class="flex justify-end gap-2 border-t bg-muted/20 p-3">
								<Button type="button" variant="ghost" size="sm" onclick={clearDateRange}>
									{m.datasets_clear()}
								</Button>
								<Button type="button" size="sm" onclick={applyDateRange}>
									{m.datasets_apply()}
								</Button>
							</div>
						</div>
					</div>
				</Popover.Content>
			</Popover.Root>
		</div>
	</div>

	{#if rowsQuery.isError && !datasetNotFound}
		<Alert.Root variant="destructive">
			<Alert.Description>{rowsQuery.error.message}</Alert.Description>
		</Alert.Root>
	{/if}

	<div class="overflow-x-auto rounded-xl border bg-background shadow-xs">
		<Table.Root>
			<Table.Header class="sticky top-0 z-10 border-b bg-muted/40">
				<Table.Row class="hover:bg-transparent">
					<Table.Head
						class="h-10 min-w-[180px] py-2.5 text-xs font-semibold uppercase tracking-wider text-muted-foreground/90"
					>
						{m.datasets_document_id_column()}
					</Table.Head>
					<Table.Head
						class="h-10 min-w-[240px] py-2.5 text-xs font-semibold uppercase tracking-wider text-muted-foreground/90"
					>
						{m.datasets_filename_column()}
					</Table.Head>
					<Table.Head
						class="h-10 w-[150px] py-2.5 text-xs font-semibold uppercase tracking-wider text-muted-foreground/90"
					>
						<Button
							type="button"
							variant="ghost"
							size="sm"
							class="-ml-2 h-7 px-2 text-xs font-semibold uppercase tracking-wider text-muted-foreground/90"
							aria-label={sortCreatedDateLabel()}
							onclick={toggleSort}
						>
							{m.datasets_created_column()}
							<ArrowUpDownIcon class="size-3.5" aria-hidden="true" />
						</Button>
					</Table.Head>
					{#each rowColumns as column (column.key)}
						<Table.Head
							class="h-10 min-w-[180px] py-2.5 text-xs font-semibold uppercase tracking-wider text-muted-foreground/90"
							title={column.path}
						>
							<span class="block max-w-[260px] truncate">{column.label}</span>
						</Table.Head>
					{/each}
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#if rowsQuery.isLoading}
					<Table.Row>
						<Table.Cell colspan={rowColumnCount} class="h-56 p-0">
							<div class="flex h-56 w-full items-center justify-center">
								<Spinner class="size-18 text-foreground dark:text-blue-500" />
							</div>
						</Table.Cell>
					</Table.Row>
				{:else if datasetNotFound}
					<Table.Row>
						<Table.Cell colspan={rowColumnCount} class="h-56 text-center">
							<div class="mx-auto flex max-w-[340px] flex-col items-center justify-center gap-2 p-6">
								<div class="rounded-full bg-muted/50 p-3.5 text-muted-foreground/80">
									<SearchIcon class="size-6" aria-hidden="true" />
								</div>
								<h3 class="mt-2 text-sm font-semibold text-foreground">
									{m.datasets_not_found_title()}
								</h3>
								<p class="text-sm text-muted-foreground">{m.datasets_not_found_body()}</p>
								<Button href="/app/datasets" variant="outline" size="sm" class="mt-2">
									{m.datasets_view_datasets()}
								</Button>
							</div>
						</Table.Cell>
					</Table.Row>
				{:else if rowsQuery.isError}
					<Table.Row>
						<Table.Cell colspan={rowColumnCount} class="h-24">
							<div
								class="mx-auto flex max-w-xl flex-wrap items-center justify-center gap-3 text-center text-sm text-destructive"
							>
								<span>{rowsQuery.error.message}</span>
								<Button type="button" variant="outline" size="sm" onclick={() => rowsQuery.refetch()}>
									{m.datasets_retry()}
								</Button>
							</div>
						</Table.Cell>
					</Table.Row>
				{:else if rows.length > 0}
					{#each rows as row (row.document_id)}
						<Table.Row class="transition-colors duration-150 hover:bg-muted/40">
							<Table.Cell class="max-w-[220px] px-4 py-3 text-xs text-muted-foreground">
								<Button
									type="button"
									variant="link"
									class="h-auto w-full min-w-0 cursor-pointer justify-start p-0 text-left text-xs font-mono"
									aria-label={m.datasets_preview_document({ documentId: row.document_id })}
									onclick={() => openDocumentPreview(row.document_id, row.filename)}
								>
									<span class="min-w-0 truncate" title={row.document_id}>{row.document_id}</span>
								</Button>
							</Table.Cell>
							<Table.Cell class="max-w-[320px] px-4 py-3 text-sm font-medium">
								<span class="block truncate" title={row.filename}>{row.filename}</span>
							</Table.Cell>
							<Table.Cell class="whitespace-nowrap px-4 py-3 text-sm text-muted-foreground">
								{formatDatasetDate(row.created_at, m.datasets_invalid_date())}
							</Table.Cell>
							{#each rowColumns as column (column.key)}
								{@const cellText = datasetCellText(row.values[column.key])}
								<Table.Cell class="max-w-[280px] px-4 py-3 text-sm" title={cellText}>
									<span class="block truncate">{cellText}</span>
								</Table.Cell>
							{/each}
						</Table.Row>
					{/each}
				{:else if rowsQuery.isSuccess && rows.length === 0}
					<Table.Row>
						<Table.Cell colspan={rowColumnCount} class="h-56 text-center">
							<div class="mx-auto flex max-w-[340px] flex-col items-center justify-center gap-2 p-6">
								<div class="rounded-full bg-muted/50 p-3.5 text-muted-foreground/80">
									<SearchIcon class="size-6" aria-hidden="true" />
								</div>
								<h3 class="mt-2 text-sm font-semibold text-foreground">
									{m.datasets_no_documents_extracted()}
								</h3>
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
			{#if rows.length > 0}
				{rowCountSummary(rows.length)}
			{:else}
				{m.datasets_no_rows_to_show()}
			{/if}
		</div>
		<div class="flex flex-wrap items-center gap-3">
			<Select.Root type="single" bind:value={() => String(pageSize), setPageSize}>
				<Select.Trigger
					size="sm"
					class="h-8 w-24 bg-background/50 text-xs"
					aria-label={m.datasets_rows_per_page()}
				>
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
					disabled={activeCursorState.history.length === 0 || rowsQuery.isFetching}
					onclick={goPrevious}
					aria-label={m.datasets_previous_page()}
				>
					<ChevronLeftIcon class="size-4" aria-hidden="true" />
				</Button>
				<Button
					type="button"
					variant="outline"
					size="icon-sm"
					disabled={!nextCursor || rowsQuery.isFetching}
					onclick={goNext}
					aria-label={m.datasets_next_page()}
				>
					<ChevronRightIcon class="size-4" aria-hidden="true" />
				</Button>
			</div>
		</div>
	</div>

	<DocumentPreviewDialog
		bind:open={previewOpen}
		filename={previewFilename}
		markdown={previewQuery.data?.markdown}
		annotationJson={previewQuery.data?.annotation_json}
		isLoading={previewQuery.isLoading}
		error={previewQuery.error}
		onRetry={() => previewQuery.refetch()}
	/>
</div>
