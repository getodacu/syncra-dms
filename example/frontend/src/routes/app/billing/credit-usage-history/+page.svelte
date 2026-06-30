<script lang="ts">
	import CalendarIcon from "@lucide/svelte/icons/calendar";
	import ChevronLeftIcon from "@lucide/svelte/icons/chevron-left";
	import ChevronRightIcon from "@lucide/svelte/icons/chevron-right";
	import CreditCardIcon from "@lucide/svelte/icons/credit-card";
	import ReceiptTextIcon from "@lucide/svelte/icons/receipt-text";
	import SearchIcon from "@lucide/svelte/icons/search";
	import HistoryIcon from "@lucide/svelte/icons/history";
	import { createQuery } from "@tanstack/svelte-query";
	import { getCoreRowModel, type ColumnDef } from "@tanstack/table-core";
	import {
		endOfMonth,
		endOfWeek,
		getLocalTimeZone,
		startOfMonth,
		startOfWeek,
		today,
	} from "@internationalized/date";
	import type { ComponentProps } from "svelte";

	import { Button } from "$lib/components/ui/button/index.js";
	import { FlexRender, renderComponent } from "$lib/components/ui/data-table/index.js";
	import { createSvelteTable } from "$lib/components/ui/data-table/data-table.svelte.js";
	import * as Popover from "$lib/components/ui/popover/index.js";
	import * as RangeCalendar from "$lib/components/ui/range-calendar/index.js";
	import * as Select from "$lib/components/ui/select/index.js";
	import { Spinner } from "$lib/components/ui/spinner/index.js";
	import * as Table from "$lib/components/ui/table/index.js";
	import { m } from "$lib/paraglide/messages.js";
	import {
		fetchCreditUsageHistory,
		type CreditUsageHistoryEntryResponse,
		type CreditUsageHistoryEntryType,
		type CreditUsageHistoryListResponse,
	} from "./api";
	import CreatedDateHeader from "./created-date-header.svelte";
	import CreditsCell from "./credits-cell.svelte";
	import {
		cursorNextState,
		cursorPreviousState,
		dateRangeToQueryBounds,
		formatCreditUsageHistoryDate,
		resetCursorState,
		type CursorState,
		type DateRangeValue,
		type SortDirection,
	} from "./table-utils";
	import TypeCell from "./type-cell.svelte";

	const PAGE_SIZE_OPTIONS = [10, 20, 50, 100];
	type RangeCalendarValue = ComponentProps<typeof RangeCalendar.RangeCalendar>["value"];

	let typeFilter = $state<CreditUsageHistoryEntryType | undefined>();
	let appliedDateRange = $state<DateRangeValue | undefined>();
	let pendingDateRange = $state<RangeCalendarValue>();
	let datePopoverOpen = $state(false);
	let sortDirection = $state<SortDirection>("desc");
	let pageSize = $state(20);
	let cursorState = $state<CursorState>(resetCursorState());

	const dateBounds = $derived(dateRangeToQueryBounds(appliedDateRange));
	const creditUsageHistoryQuery = createQuery<CreditUsageHistoryListResponse, Error>(() => ({
		queryKey: [
			"credit-usage-history",
			{
				type: typeFilter,
				createdFrom: dateBounds.createdFrom,
				createdTo: dateBounds.createdTo,
				cursor: cursorState.currentCursor,
				size: pageSize,
				sort: sortDirection,
			},
		],
		queryFn: () =>
			fetchCreditUsageHistory(fetch, {
				type: typeFilter,
				createdFrom: dateBounds.createdFrom,
				createdTo: dateBounds.createdTo,
				cursor: cursorState.currentCursor,
				size: pageSize,
				sort: sortDirection,
			}),
	}));

	const creditUsageHistory = $derived(
		creditUsageHistoryQuery.data?.credit_usage_history ?? []
	);
	const nextCursor = $derived(creditUsageHistoryQuery.data?.next_cursor ?? null);
	const filtersActive = $derived(Boolean(typeFilter || appliedDateRange?.start || appliedDateRange?.end));
	const dateRangeLabel = $derived.by(() => {
		if (!appliedDateRange?.start && !appliedDateRange?.end) return m.credit_usage_date_range_filter();

		const start = appliedDateRange.start ? appliedDateRange.start.toString() : m.common_any();
		const end = appliedDateRange.end ? appliedDateRange.end.toString() : m.common_any();

		return `${start} - ${end}`;
	});
	const activePreset = $derived.by(() => {
		if (!appliedDateRange?.start || !appliedDateRange?.end) return null;

		const t = today(getLocalTimeZone());
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

	function resetPagination() {
		cursorState = resetCursorState();
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
		const t = today(getLocalTimeZone());
		pendingDateRange = { start: t, end: t };
		appliedDateRange = pendingDateRange;
		datePopoverOpen = false;
		resetPagination();
	}

	function setThisWeekPreset() {
		const t = today(getLocalTimeZone());
		pendingDateRange = {
			start: startOfWeek(t, "en-US"),
			end: endOfWeek(t, "en-US"),
		};
		appliedDateRange = pendingDateRange;
		datePopoverOpen = false;
		resetPagination();
	}

	function setThisMonthPreset() {
		const t = today(getLocalTimeZone());
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

	function setTypeFilter(value: string) {
		typeFilter = value === "purchase" || value === "debit" ? value : undefined;
		resetPagination();
	}

	function setPageSize(value: string) {
		const nextPageSize = Number(value);
		if (!PAGE_SIZE_OPTIONS.includes(nextPageSize)) return;

		pageSize = nextPageSize;
		resetPagination();
	}

	function toggleSort() {
		sortDirection = sortDirection === "desc" ? "asc" : "desc";
		resetPagination();
	}

	function goNext() {
		cursorState = cursorNextState(cursorState, nextCursor);
	}

	function goPrevious() {
		cursorState = cursorPreviousState(cursorState);
	}

	function clearFilters() {
		typeFilter = undefined;
		clearDateRange();
	}

	function relatedId(entry: CreditUsageHistoryEntryResponse) {
		return entry.related_order_id ?? entry.related_job_id ?? "-";
	}

	function typeFilterLabel(type?: CreditUsageHistoryEntryType) {
		if (type === "purchase") return m.credit_usage_type_purchase();
		if (type === "debit") return m.credit_usage_type_debit();
		return m.credit_usage_all_activity();
	}

	const columns: ColumnDef<CreditUsageHistoryEntryResponse>[] = [
		{
			accessorKey: "created_at",
			header: () =>
				renderComponent(CreatedDateHeader, {
					sortDirection,
					onToggle: toggleSort,
					}),
			cell: ({ row }) => formatCreditUsageHistoryDate(row.original.created_at),
		},
		{
			accessorKey: "entry_type",
			header: m.credit_usage_type_column(),
			cell: ({ row }) => renderComponent(TypeCell, { type: row.original.entry_type }),
		},
		{
			accessorKey: "credits_delta",
			header: m.credit_usage_credits_column(),
			cell: ({ row }) => renderComponent(CreditsCell, { credits: row.original.credits_delta }),
		},
		{
			id: "related_id",
			header: m.credit_usage_related_id_column(),
			cell: ({ row }) => relatedId(row.original),
		},
	];

	const table = createSvelteTable({
		get data() {
			return creditUsageHistory;
		},
		columns,
		getRowId: (row) => row.id,
		getCoreRowModel: getCoreRowModel(),
	});
</script>

<svelte:head>
	<title>{m.credit_usage_page_title()} | Syncra</title>
</svelte:head>

<div class="@container/main flex flex-1 flex-col gap-5 p-4 lg:p-6">
	<div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-end">
			<div class="flex items-center gap-2">
				<Button href="/app/billing/orders" variant="outline" size="sm" class="gap-2 font-medium cursor-pointer">
					<HistoryIcon class="size-4 text-muted-foreground" aria-hidden="true" />
					<span>{m.nav_billing_orders()}</span>
				</Button>
				<Button href="/app/billing" size="sm" variant="outline" class="gap-2 font-medium cursor-pointer">
					<CreditCardIcon class="size-4 text-muted-foreground" aria-hidden="true" />
					<span>{m.billing_buy_credits()}</span>
				</Button>
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
							class="h-9 w-full min-w-[160px] justify-start bg-background/50 text-xs sm:w-auto font-medium"
							{...props}
						>
							<CalendarIcon class="mr-2 size-4 shrink-0 text-muted-foreground" aria-hidden="true" />
							<span class="truncate">{dateRangeLabel}</span>
						</Button>
					{/snippet}
				</Popover.Trigger>
				<Popover.Content align="start" class="w-auto p-0">
					<div class="flex flex-col sm:flex-row">
						<div class="flex min-w-[130px] flex-row gap-1.5 border-b border-border bg-muted/5 p-3 sm:flex-col sm:border-r sm:border-b-0">
							<span class="hidden select-none px-2 py-1 text-[10px] font-bold tracking-wider text-muted-foreground/60 uppercase sm:inline-block">{m.billing_orders_presets()}</span>
							<Button
								type="button"
								variant={activePreset === "today" ? "secondary" : "ghost"}
								class="h-8.5 flex-1 justify-center px-3 text-xs font-medium sm:flex-none sm:justify-start"
								onclick={setTodayPreset}
							>
								{m.common_today()}
							</Button>
							<Button
								type="button"
								variant={activePreset === "week" ? "secondary" : "ghost"}
								class="h-8.5 flex-1 justify-center px-3 text-xs font-medium sm:flex-none sm:justify-start"
								onclick={setThisWeekPreset}
							>
								{m.common_this_week()}
							</Button>
							<Button
								type="button"
								variant={activePreset === "month" ? "secondary" : "ghost"}
								class="h-8.5 flex-1 justify-center px-3 text-xs font-medium sm:flex-none sm:justify-start"
								onclick={setThisMonthPreset}
							>
								{m.common_this_month()}
							</Button>
						</div>
						<div class="flex flex-col">
							<RangeCalendar.RangeCalendar bind:value={pendingDateRange} numberOfMonths={2} />
							<div class="flex justify-end gap-2 border-t bg-muted/20 p-3">
								<Button type="button" variant="ghost" size="sm" onclick={clearDateRange}>{m.common_clear()}</Button>
								<Button type="button" size="sm" onclick={applyDateRange}>{m.common_apply()}</Button>
							</div>
						</div>
					</div>
				</Popover.Content>
			</Popover.Root>

			<Select.Root type="single" value={typeFilter ?? "all"} onValueChange={setTypeFilter}>
				<Select.Trigger class="h-9 w-full justify-between bg-background/50 text-xs sm:w-40" aria-label={m.credit_usage_filter_type()}>
					<div class="flex min-w-0 items-center">
						<ReceiptTextIcon class="mr-2 size-3.5 shrink-0 text-muted-foreground" aria-hidden="true" />
						<span class="truncate text-left font-medium">
							{typeFilterLabel(typeFilter)}
						</span>
					</div>
				</Select.Trigger>
				<Select.Content side="bottom">
					<Select.Item value="all" class="text-xs">{m.credit_usage_all_activity()}</Select.Item>
					<Select.Item value="purchase" class="text-xs">{m.credit_usage_type_purchase()}</Select.Item>
					<Select.Item value="debit" class="text-xs">{m.credit_usage_type_debit()}</Select.Item>
				</Select.Content>
			</Select.Root>

			{#if filtersActive}
				<Button
					type="button"
					variant="ghost"
					size="sm"
					class="h-9 text-xs text-muted-foreground hover:text-foreground font-medium cursor-pointer"
					onclick={clearFilters}
				>
					{m.billing_orders_clear_filters()}
				</Button>
			{/if}
		</div>
	</div>

	<div class="overflow-x-auto rounded-xl border bg-background shadow-xs">
		<Table.Root class="min-w-[720px]">
			<Table.Header class="sticky top-0 z-10 border-b bg-muted/40">
				{#each table.getHeaderGroups() as headerGroup (headerGroup.id)}
					<Table.Row class="hover:bg-transparent">
						{#each headerGroup.headers as header (header.id)}
							<Table.Head colspan={header.colSpan} class="h-10 py-2.5 text-xs font-semibold tracking-wider text-muted-foreground/90 uppercase">
								{#if !header.isPlaceholder}
									<FlexRender content={header.column.columnDef.header} context={header.getContext()} />
								{/if}
							</Table.Head>
						{/each}
					</Table.Row>
				{/each}
			</Table.Header>
			<Table.Body>
				{#if creditUsageHistoryQuery.isLoading}
					<Table.Row>
						<Table.Cell colspan={columns.length} class="h-56 p-0">
							<div class="flex h-56 w-full items-center justify-center">
								<Spinner class="size-18 text-foreground dark:text-blue-500" />
							</div>
						</Table.Cell>
					</Table.Row>
				{:else if creditUsageHistoryQuery.isError}
					<Table.Row>
						<Table.Cell colspan={columns.length} class="h-24">
							<div class="mx-auto flex max-w-xl flex-wrap items-center justify-center gap-3 text-center text-sm text-destructive">
								<span>{creditUsageHistoryQuery.error.message}</span>
								<Button
									type="button"
									variant="outline"
									size="sm"
									onclick={() => creditUsageHistoryQuery.refetch()}
								>
									{m.common_retry()}
								</Button>
							</div>
						</Table.Cell>
					</Table.Row>
				{:else if table.getRowModel().rows.length}
					{#each table.getRowModel().rows as row (row.id)}
						<Table.Row class="transition-colors duration-150 hover:bg-muted/40">
							{#each row.getVisibleCells() as cell (cell.id)}
								<Table.Cell class="px-4 py-3">
									{#if cell.column.id === "related_id"}
										<span class="block max-w-[18rem] truncate font-mono text-xs text-muted-foreground" title={relatedId(row.original)}>
											<FlexRender content={cell.column.columnDef.cell} context={cell.getContext()} />
										</span>
									{:else}
										<FlexRender content={cell.column.columnDef.cell} context={cell.getContext()} />
									{/if}
								</Table.Cell>
							{/each}
						</Table.Row>
					{/each}
				{:else if creditUsageHistoryQuery.isSuccess && creditUsageHistory.length === 0}
					<Table.Row>
						<Table.Cell colspan={columns.length} class="h-56 text-center">
							<div class="mx-auto flex max-w-[340px] flex-col items-center justify-center gap-2 p-6">
								<div class="rounded-full bg-muted/50 p-3.5 text-muted-foreground/80">
									<SearchIcon class="size-6" aria-hidden="true" />
								</div>
								<h3 class="mt-2 text-sm font-semibold text-foreground">
									{filtersActive ? m.credit_usage_no_usage_found() : m.credit_usage_no_usage_yet()}
								</h3>
								<p class="px-2 text-xs leading-normal text-muted-foreground">
									{#if filtersActive}
										{m.credit_usage_no_usage_match()}
									{:else}
										{m.credit_usage_empty_body()}
									{/if}
								</p>
								{#if filtersActive}
									<Button
										type="button"
										variant="outline"
										size="sm"
										class="mt-3.5 h-8 text-xs font-medium"
										onclick={clearFilters}
									>
										{m.billing_orders_clear_filters_action()}
									</Button>
								{/if}
							</div>
						</Table.Cell>
					</Table.Row>
				{/if}
			</Table.Body>
		</Table.Root>
	</div>

	<div class="flex flex-col gap-3 px-1.5 text-xs text-muted-foreground sm:flex-row sm:items-center sm:justify-between">
		<div>
			{#if creditUsageHistory.length > 0}
				{creditUsageHistory.length === 1
					? m.credit_usage_showing_one({ count: creditUsageHistory.length })
					: m.credit_usage_showing_other({ count: creditUsageHistory.length })}
			{:else}
				{m.credit_usage_none_to_show()}
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
					disabled={cursorState.history.length === 0 || creditUsageHistoryQuery.isFetching}
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
					disabled={!nextCursor || creditUsageHistoryQuery.isFetching}
				>
					{m.common_next()}
					<ChevronRightIcon class="ml-1 size-4" aria-hidden="true" />
				</Button>
			</div>
		</div>
	</div>
</div>
