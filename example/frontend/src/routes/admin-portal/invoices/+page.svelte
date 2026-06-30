<script lang="ts">
	import CalendarIcon from "@lucide/svelte/icons/calendar";
	import CheckIcon from "@lucide/svelte/icons/check";
	import ChevronLeftIcon from "@lucide/svelte/icons/chevron-left";
	import ChevronRightIcon from "@lucide/svelte/icons/chevron-right";
	import ChevronsUpDownIcon from "@lucide/svelte/icons/chevrons-up-down";
	import DownloadIcon from "@lucide/svelte/icons/download";
	import SearchIcon from "@lucide/svelte/icons/search";
	import UserIcon from "@lucide/svelte/icons/user";
	import XIcon from "@lucide/svelte/icons/x";
	import { createQuery } from "@tanstack/svelte-query";
	import { getCoreRowModel, type ColumnDef } from "@tanstack/table-core";
	import {
		endOfMonth,
		endOfWeek,
		getLocalTimeZone,
		startOfMonth,
		startOfWeek,
		today
	} from "@internationalized/date";
	import type { ComponentProps } from "svelte";

	import { Button } from "$lib/components/ui/button/index.js";
	import * as Command from "$lib/components/ui/command/index.js";
	import { FlexRender, renderComponent } from "$lib/components/ui/data-table/index.js";
	import { createSvelteTable } from "$lib/components/ui/data-table/data-table.svelte.js";
	import * as Dialog from "$lib/components/ui/dialog/index.js";
	import { Input } from "$lib/components/ui/input/index.js";
	import * as Popover from "$lib/components/ui/popover/index.js";
	import * as RangeCalendar from "$lib/components/ui/range-calendar/index.js";
	import * as Select from "$lib/components/ui/select/index.js";
	import { Spinner } from "$lib/components/ui/spinner/index.js";
	import * as Table from "$lib/components/ui/table/index.js";
	import {
		cursorNextState,
		cursorPreviousState,
		dateRangeToQueryBounds,
		formatPaymentDateTime,
		resetCursorState,
		type CursorState,
		type DateRangeValue,
		type SortDirection
	} from "../../app/billing/orders/table-utils";
	import {
		fetchAdminUsers,
		type AdminUserListResponse,
		type AdminUserResponse
	} from "../users/api";
	import {
		buildAdminBillingInvoicePDFPath,
		fetchAdminBillingInvoices,
		type AdminBillingInvoiceResponse,
		type AdminBillingInvoicesListResponse
	} from "./api";
	import CreatedDateHeader from "./created-date-header.svelte";
	import InvoicePDFActionsCell from "./invoice-pdf-actions-cell.svelte";

	const PAGE_SIZE_OPTIONS = [10, 20, 50, 100];
	const USER_SEARCH_PAGE_SIZE = 20;
	type RangeCalendarValue = ComponentProps<typeof RangeCalendar.RangeCalendar>["value"];

	let pendingSearch = $state("");
	let search = $state("");
	let selectedUser = $state<AdminUserResponse | null>(null);
	let userSearch = $state("");
	let userPopoverOpen = $state(false);
	let appliedDateRange = $state<DateRangeValue | undefined>();
	let pendingDateRange = $state<RangeCalendarValue>();
	let datePopoverOpen = $state(false);
	let sortDirection = $state<SortDirection>("desc");
	let pageSize = $state(20);
	let cursorState = $state<CursorState>(resetCursorState());
	let previewInvoice = $state<AdminBillingInvoiceResponse | null>(null);

	const dateBounds = $derived(dateRangeToQueryBounds(appliedDateRange));
	const usersQuery = createQuery<AdminUserListResponse, Error>(() => ({
		queryKey: ["admin-invoice-user-search", { search: userSearch }],
		queryFn: () =>
			fetchAdminUsers(fetch, {
				search: userSearch,
				sort: "created_at",
				direction: "desc",
				size: USER_SEARCH_PAGE_SIZE
			}),
		enabled: userPopoverOpen
	}));
	const billingInvoicesQuery = createQuery<AdminBillingInvoicesListResponse, Error>(() => ({
		queryKey: [
			"admin-billing-invoices",
			{
				search,
				userId: selectedUser?.id,
				createdFrom: dateBounds.createdFrom,
				createdTo: dateBounds.createdTo,
				cursor: cursorState.currentCursor,
				size: pageSize,
				sort: sortDirection
			}
		],
		queryFn: () =>
			fetchAdminBillingInvoices(fetch, {
				search,
				userId: selectedUser?.id,
				createdFrom: dateBounds.createdFrom,
				createdTo: dateBounds.createdTo,
				cursor: cursorState.currentCursor,
				size: pageSize,
				sort: sortDirection
			})
	}));

	const userOptions = $derived(usersQuery.data?.users ?? []);
	const invoices = $derived(billingInvoicesQuery.data?.invoices ?? []);
	const nextCursor = $derived(billingInvoicesQuery.data?.next_cursor ?? null);
	const previewInvoiceLabel = $derived(previewInvoice ? formatInvoiceNumber(previewInvoice) : "");
	const previewURL = $derived(previewInvoice ? buildAdminBillingInvoicePDFPath(previewInvoice.id) : "");
	const previewDownloadURL = $derived(
		previewInvoice ? buildAdminBillingInvoicePDFPath(previewInvoice.id, { download: true }) : ""
	);
	const previewFilename = $derived(previewInvoice ? formatInvoiceFilename(previewInvoice) : "");
	const filtersActive = $derived(Boolean(search || selectedUser || appliedDateRange?.start || appliedDateRange?.end));
	const selectedUserLabel = $derived.by(() => {
		if (!selectedUser) return "All users";
		return `${selectedUser.name || "Unnamed user"} (${selectedUser.email})`;
	});
	const dateRangeLabel = $derived.by(() => {
		if (!appliedDateRange?.start && !appliedDateRange?.end) return "Created date";

		const start = appliedDateRange.start ? appliedDateRange.start.toString() : "Any";
		const end = appliedDateRange.end ? appliedDateRange.end.toString() : "Any";

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

	function applySearch() {
		search = pendingSearch.trim();
		resetPagination();
	}

	function rangeCalendarValue(range?: DateRangeValue): RangeCalendarValue {
		if (!range?.start && !range?.end) return undefined;

		return {
			start: range.start,
			end: range.end
		};
	}

	function setDatePopoverOpen(open: boolean) {
		datePopoverOpen = open;
		if (open) pendingDateRange = rangeCalendarValue(appliedDateRange);
	}

	function selectUser(user: AdminUserResponse) {
		selectedUser = user;
		userPopoverOpen = false;
		userSearch = "";
		resetPagination();
	}

	function clearUserFilter() {
		selectedUser = null;
		userSearch = "";
		userPopoverOpen = false;
		resetPagination();
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
			end: endOfWeek(t, "en-US")
		};
		appliedDateRange = pendingDateRange;
		datePopoverOpen = false;
		resetPagination();
	}

	function setThisMonthPreset() {
		const t = today(getLocalTimeZone());
		pendingDateRange = {
			start: startOfMonth(t),
			end: endOfMonth(t)
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
		pendingSearch = "";
		search = "";
		selectedUser = null;
		userSearch = "";
		userPopoverOpen = false;
		pendingDateRange = undefined;
		appliedDateRange = undefined;
		datePopoverOpen = false;
		resetPagination();
	}

	function formatInvoiceNumber(invoice: AdminBillingInvoiceResponse) {
		return `${invoice.invoice_serie}-${String(invoice.invoice_number).padStart(5, "0")}`;
	}

	function formatInvoiceFilename(invoice: AdminBillingInvoiceResponse) {
		return `${invoice.invoice_serie}_${String(invoice.invoice_number).padStart(5, "0")}_${formatInvoiceDateStamp(invoice.invoice_date)}.pdf`;
	}

	function formatInvoiceDateStamp(invoiceDate: string) {
		const match = invoiceDate.match(/^(\d{4})-(\d{2})-(\d{2})/);
		if (!match) return invoiceDate.replace(/\D/g, "").slice(2, 8);

		const [, year, month, day] = match;
		return `${year.slice(2)}${month}${day}`;
	}

	function formatInvoiceAmount(value: string) {
		const amount = Number(value);
		if (!Number.isFinite(amount)) return value;

		return new Intl.NumberFormat("en-US", {
			style: "currency",
			currency: "EUR"
		}).format(amount);
	}

	function openInvoicePreview(invoice: AdminBillingInvoiceResponse) {
		if (!invoice.pdf_path) return;
		previewInvoice = invoice;
	}

	function setPreviewOpen(open: boolean) {
		if (!open) previewInvoice = null;
	}

	const columns: ColumnDef<AdminBillingInvoiceResponse>[] = [
		{
			id: "invoice_number",
			header: "Invoice",
			cell: ({ row }) => formatInvoiceNumber(row.original)
		},
		{
			id: "client",
			header: "Client",
			cell: ({ row }) => row.original.billing_name
		},
		{
			accessorKey: "net_amount",
			header: "Net amount",
			cell: ({ row }) => formatInvoiceAmount(row.original.net_amount)
		},
		{
			accessorKey: "vat_amount",
			header: "VAT amount",
			cell: ({ row }) => formatInvoiceAmount(row.original.vat_amount)
		},
		{
			accessorKey: "total_amount",
			header: "Total amount",
			cell: ({ row }) => formatInvoiceAmount(row.original.total_amount)
		},
		{
			accessorKey: "created_at",
			header: () =>
				renderComponent(CreatedDateHeader, {
					sortDirection,
					onToggle: toggleSort
				}),
			cell: ({ row }) => formatPaymentDateTime(row.original.created_at)
		},
		{
			id: "actions",
			header: "",
			cell: ({ row }) =>
				renderComponent(InvoicePDFActionsCell, {
					invoice: row.original,
					onGenerated: async () => {
						await billingInvoicesQuery.refetch();
					}
				})
		}
	];

	const table = createSvelteTable({
		get data() {
			return invoices;
		},
		columns,
		getRowId: (row) => row.id,
		getCoreRowModel: getCoreRowModel()
	});
</script>

<svelte:head>
	<title>Billing Invoices | Syncra Admin</title>
</svelte:head>

<div class="@container/main flex flex-1 flex-col gap-5 p-4 lg:p-6">
	<div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
		<div class="flex min-w-0 flex-1 flex-wrap items-center gap-2">
			<form
				class="flex min-w-0 flex-1 gap-2 sm:min-w-[340px] sm:max-w-xl"
				onsubmit={(event) => {
					event.preventDefault();
					applySearch();
				}}
			>
				<div class="relative min-w-0 flex-1">
					<SearchIcon class="pointer-events-none absolute left-2.5 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" aria-hidden="true" />
					<Input
						class="h-9 pl-8 text-xs"
						placeholder="Search client or invoice number"
						value={pendingSearch}
						oninput={(event) => (pendingSearch = event.currentTarget.value)}
					/>
				</div>
				<Button type="submit" size="sm" class="h-9 px-4 text-xs">Search</Button>
			</form>

			<Popover.Root bind:open={userPopoverOpen}>
				<Popover.Trigger>
					{#snippet child({ props })}
						<Button
							type="button"
							variant="outline"
							class="h-9 w-full min-w-[240px] justify-between bg-background/50 text-xs font-medium sm:w-[320px]"
							role="combobox"
							aria-expanded={userPopoverOpen}
							{...props}
						>
							<span class="flex min-w-0 items-center gap-2">
								<UserIcon class="size-4 shrink-0 text-muted-foreground" aria-hidden="true" />
								<span class="truncate text-left">{selectedUserLabel}</span>
							</span>
							<ChevronsUpDownIcon class="size-4 shrink-0 text-muted-foreground" aria-hidden="true" />
						</Button>
					{/snippet}
				</Popover.Trigger>
				<Popover.Content class="w-[min(calc(100vw-2rem),28rem)] p-0" align="start">
					<Command.Root>
						<Command.Input bind:value={userSearch} placeholder="Search users by name or email" />
						<Command.List>
							{#if usersQuery.isLoading || usersQuery.isFetching}
								<div class="px-3 py-6 text-center text-sm text-muted-foreground">
									Loading users
								</div>
							{:else if usersQuery.isError}
								<div class="px-3 py-6 text-center text-sm text-destructive">
									{usersQuery.error.message}
								</div>
							{:else}
								<Command.Empty>No users found.</Command.Empty>
								{#each userOptions as user (user.id)}
									<Command.Item
										value={user.id}
										keywords={[user.name, user.email, user.id]}
										onSelect={() => selectUser(user)}
									>
										<div class="flex min-w-0 flex-col">
											<span class="truncate">{user.name || "Unnamed user"}</span>
											<span class="truncate text-xs text-muted-foreground">{user.email}</span>
										</div>
										{#if selectedUser?.id === user.id}
											<CheckIcon class="ml-auto size-4" aria-hidden="true" />
										{/if}
									</Command.Item>
								{/each}
							{/if}
						</Command.List>
					</Command.Root>
				</Popover.Content>
			</Popover.Root>

			<Popover.Root bind:open={() => datePopoverOpen, setDatePopoverOpen}>
				<Popover.Trigger>
					{#snippet child({ props })}
						<Button
							type="button"
							variant="outline"
							class="h-9 w-full min-w-[160px] justify-start bg-background/50 text-xs font-medium sm:w-auto"
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
							<span class="hidden select-none px-2 py-1 text-[10px] font-bold tracking-wider text-muted-foreground/60 uppercase sm:inline-block">Presets</span>
							<Button
								type="button"
								variant={activePreset === "today" ? "secondary" : "ghost"}
								class="h-8.5 flex-1 justify-center px-3 text-xs font-medium sm:flex-none sm:justify-start"
								onclick={setTodayPreset}
							>
								Today
							</Button>
							<Button
								type="button"
								variant={activePreset === "week" ? "secondary" : "ghost"}
								class="h-8.5 flex-1 justify-center px-3 text-xs font-medium sm:flex-none sm:justify-start"
								onclick={setThisWeekPreset}
							>
								This week
							</Button>
							<Button
								type="button"
								variant={activePreset === "month" ? "secondary" : "ghost"}
								class="h-8.5 flex-1 justify-center px-3 text-xs font-medium sm:flex-none sm:justify-start"
								onclick={setThisMonthPreset}
							>
								This month
							</Button>
						</div>
						<div class="flex flex-col">
							<RangeCalendar.RangeCalendar bind:value={pendingDateRange} numberOfMonths={2} />
							<div class="flex justify-end gap-2 border-t bg-muted/20 p-3">
								<Button type="button" variant="ghost" size="sm" onclick={clearDateRange}>Clear</Button>
								<Button type="button" size="sm" onclick={applyDateRange}>Apply</Button>
							</div>
						</div>
					</div>
				</Popover.Content>
			</Popover.Root>

			{#if selectedUser}
				<Button
					type="button"
					variant="ghost"
					size="sm"
					class="h-9 text-xs font-medium text-muted-foreground hover:text-foreground"
					onclick={clearUserFilter}
				>
					<XIcon class="mr-1 size-3.5" aria-hidden="true" />
					Clear user
				</Button>
			{/if}

			{#if filtersActive}
				<Button
					type="button"
					variant="ghost"
					size="sm"
					class="h-9 text-xs font-medium text-muted-foreground hover:text-foreground"
					onclick={clearFilters}
				>
					Clear filters
				</Button>
			{/if}
		</div>
	</div>

	<div class="overflow-x-auto rounded-xl border bg-background shadow-xs">
		<Table.Root class="min-w-[1280px]">
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
				{#if billingInvoicesQuery.isLoading}
					<Table.Row>
						<Table.Cell colspan={columns.length} class="h-56 p-0">
							<div class="flex h-56 w-full items-center justify-center">
								<Spinner class="size-18 text-foreground dark:text-blue-500" />
							</div>
						</Table.Cell>
					</Table.Row>
				{:else if billingInvoicesQuery.isError}
					<Table.Row>
						<Table.Cell colspan={columns.length} class="h-24">
							<div class="mx-auto flex max-w-xl flex-wrap items-center justify-center gap-3 text-center text-sm text-destructive">
								<span>{billingInvoicesQuery.error.message}</span>
								<Button
									type="button"
									variant="outline"
									size="sm"
									onclick={() => billingInvoicesQuery.refetch()}
								>
									Retry
								</Button>
							</div>
						</Table.Cell>
					</Table.Row>
				{:else if table.getRowModel().rows.length}
					{#each table.getRowModel().rows as row (row.id)}
						<Table.Row class="transition-colors duration-150 hover:bg-muted/40">
							{#each row.getVisibleCells() as cell (cell.id)}
								<Table.Cell class="px-4 py-3">
									{#if cell.column.id === "invoice_number"}
										{#if row.original.pdf_path}
											<button
												type="button"
												class="text-left text-sm font-semibold text-primary underline-offset-4 hover:underline focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:outline-none"
												title={`Preview ${formatInvoiceNumber(row.original)} PDF`}
												onclick={() => openInvoicePreview(row.original)}
											>
												{formatInvoiceNumber(row.original)}
											</button>
										{:else}
											<span class="text-sm font-semibold text-foreground">
												{formatInvoiceNumber(row.original)}
											</span>
										{/if}
									{:else if cell.column.id === "client"}
										{#if row.original.user_id}
											<a
												href={`/admin-portal/users/${row.original.user_id}`}
												class="flex min-w-0 items-center gap-3 hover:underline"
											>
												<div class="flex size-8 shrink-0 items-center justify-center rounded-md bg-muted text-muted-foreground">
													<UserIcon class="size-4" aria-hidden="true" />
												</div>
												<div class="min-w-0">
													<div class="truncate text-sm font-medium">{row.original.billing_name || "Unnamed client"}</div>
													<div class="truncate text-xs text-muted-foreground">{row.original.billing_email}</div>
												</div>
											</a>
										{:else}
											<div class="flex min-w-0 items-center gap-3">
												<div class="flex size-8 shrink-0 items-center justify-center rounded-md bg-muted text-muted-foreground">
													<UserIcon class="size-4" aria-hidden="true" />
												</div>
												<div class="min-w-0">
													<div class="truncate text-sm font-medium">{row.original.billing_name || "Unnamed client"}</div>
													<div class="truncate text-xs text-muted-foreground">{row.original.billing_email}</div>
												</div>
											</div>
										{/if}
									{:else}
										<FlexRender content={cell.column.columnDef.cell} context={cell.getContext()} />
									{/if}
								</Table.Cell>
							{/each}
						</Table.Row>
					{/each}
				{:else if billingInvoicesQuery.isSuccess && invoices.length === 0}
					<Table.Row>
						<Table.Cell colspan={columns.length} class="h-56 text-center">
							<div class="mx-auto flex max-w-[340px] flex-col items-center justify-center gap-2 p-6">
								<div class="rounded-full bg-muted/50 p-3.5 text-muted-foreground/80">
									<SearchIcon class="size-6" aria-hidden="true" />
								</div>
								<h3 class="mt-2 text-sm font-semibold text-foreground">
									{filtersActive ? "No billing invoices found" : "No billing invoices yet"}
								</h3>
								<p class="px-2 text-xs leading-normal text-muted-foreground">
									{#if filtersActive}
										No billing invoices match the selected filters.
									{:else}
										Billing invoices will appear here after invoices are generated.
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
										Clear filters
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
			{#if invoices.length > 0}
				Showing <span class="font-semibold text-foreground">{invoices.length}</span> invoice{invoices.length === 1 ? "" : "s"} on this page.
			{:else}
				No billing invoices to show.
			{/if}
		</div>
		<div class="flex flex-wrap items-center gap-3">
			<Select.Root type="single" bind:value={() => String(pageSize), setPageSize}>
				<Select.Trigger size="sm" class="h-8 w-24 bg-background/50 text-xs" aria-label="Rows per page">
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
					disabled={cursorState.history.length === 0 || billingInvoicesQuery.isFetching}
				>
					<ChevronLeftIcon class="mr-1 size-4" aria-hidden="true" />
					Previous
				</Button>
				<Button
					type="button"
					variant="outline"
					size="sm"
					class="h-8 text-xs"
					onclick={goNext}
					disabled={!nextCursor || billingInvoicesQuery.isFetching}
				>
					Next
					<ChevronRightIcon class="ml-1 size-4" aria-hidden="true" />
				</Button>
			</div>
		</div>
	</div>

	<Dialog.Root bind:open={() => Boolean(previewInvoice), setPreviewOpen}>
		<Dialog.Content
			class="flex h-[min(90vh,56rem)] w-full flex-col overflow-hidden p-0 sm:max-w-[min(94vw,72rem)]"
		>
			{#if previewInvoice}
				<Dialog.Header class="border-b px-5 py-4 pr-14">
					<div class="flex min-w-0 items-center justify-between gap-3">
						<div class="min-w-0">
							<Dialog.Title class="truncate text-lg">Invoice {previewInvoiceLabel}</Dialog.Title>
							<Dialog.Description>PDF preview</Dialog.Description>
						</div>
						<Button
							href={previewDownloadURL}
							download={previewFilename}
							variant="outline"
							size="sm"
							class="h-8 gap-1.5 text-xs"
						>
							<DownloadIcon class="size-3.5" aria-hidden="true" />
							Download
						</Button>
					</div>
				</Dialog.Header>

				<div class="min-h-0 flex-1 bg-muted/30">
					<iframe
						title={`Invoice ${previewInvoiceLabel} PDF preview`}
						src={previewURL}
						class="h-full w-full border-0 bg-background"
					></iframe>
				</div>
			{/if}
		</Dialog.Content>
	</Dialog.Root>

	{#if billingInvoicesQuery.isFetching && !billingInvoicesQuery.isLoading}
		<div class="pointer-events-none fixed right-4 bottom-4 flex items-center gap-2 rounded-md border bg-background/95 px-3 py-2 text-xs text-muted-foreground shadow-sm">
			<Spinner class="size-4" />
			Updating invoices
		</div>
	{/if}
</div>
