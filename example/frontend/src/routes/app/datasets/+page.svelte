<script lang="ts">
	import ArrowUpDownIcon from "@lucide/svelte/icons/arrow-up-down";
	import ChevronLeftIcon from "@lucide/svelte/icons/chevron-left";
	import ChevronRightIcon from "@lucide/svelte/icons/chevron-right";
	import DatabaseIcon from "@lucide/svelte/icons/database";
	import ExternalLinkIcon from "@lucide/svelte/icons/external-link";
	import SearchIcon from "@lucide/svelte/icons/search";
	import { createQuery } from "@tanstack/svelte-query";

	import { datasetHref } from "$lib/components/nav-datasets-utils";
	import * as Alert from "$lib/components/ui/alert/index.js";
	import { Button } from "$lib/components/ui/button/index.js";
	import * as Select from "$lib/components/ui/select/index.js";
	import { Spinner } from "$lib/components/ui/spinner/index.js";
	import * as Table from "$lib/components/ui/table/index.js";
	import { m } from "$lib/paraglide/messages.js";
	import { fetchDatasets, type DatasetListResponse } from "./api";
	import {
		cursorNextState,
		cursorPreviousState,
		formatDatasetDate,
		resetCursorState,
		type CursorState,
		type SortDirection,
	} from "./table-utils";

	const PAGE_SIZE_OPTIONS = [10, 20, 50, 100];
	const DATASET_LIST_COLUMN_COUNT = 5;

	let sortDirection = $state<SortDirection>("desc");
	let pageSize = $state(20);
	let cursorState = $state<CursorState>(resetCursorState());

	const datasetsQuery = createQuery<DatasetListResponse, Error>(() => ({
		queryKey: [
			"datasets",
			{
				cursor: cursorState.currentCursor,
				size: pageSize,
				sort: sortDirection,
			},
		],
		queryFn: () =>
			fetchDatasets(fetch, {
				cursor: cursorState.currentCursor,
				size: pageSize,
				sort: sortDirection,
			}),
	}));

	const datasets = $derived(datasetsQuery.data?.datasets ?? []);
	const nextCursor = $derived(datasetsQuery.data?.next_cursor ?? null);

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

	function goNext() {
		cursorState = cursorNextState(cursorState, nextCursor);
	}

	function goPrevious() {
		cursorState = cursorPreviousState(cursorState);
	}

	function fieldCountLabel(count: number) {
		return count === 1
			? m.datasets_field_count_one({ count })
			: m.datasets_field_count_other({ count });
	}

	function datasetCountSummary(count: number) {
		return count === 1
			? m.datasets_showing_datasets_one({ count })
			: m.datasets_showing_datasets_other({ count });
	}

	function sortCreatedDateLabel() {
		return sortDirection === "desc"
			? m.datasets_sort_created_ascending()
			: m.datasets_sort_created_descending();
	}
</script>

<svelte:head>
	<title>{m.datasets_page_title()} | Syncra</title>
</svelte:head>

<div class="@container/main flex flex-1 flex-col gap-4 p-4 lg:p-6">

	{#if datasetsQuery.isError}
		<Alert.Root variant="destructive">
			<Alert.Description>{datasetsQuery.error.message}</Alert.Description>
		</Alert.Root>
	{/if}

	<div class="overflow-x-auto rounded-xl border bg-background shadow-xs">
		<Table.Root>
			<Table.Header class="sticky top-0 z-10 border-b bg-muted/40">
				<Table.Row class="hover:bg-transparent">
					<Table.Head
						class="h-10 min-w-[260px] py-2.5 text-xs font-semibold uppercase tracking-wider text-muted-foreground/90"
					>
						{m.datasets_name_column()}
					</Table.Head>
					<Table.Head
						class="h-10 min-w-[180px] py-2.5 text-xs font-semibold uppercase tracking-wider text-muted-foreground/90"
					>
						{m.datasets_schema_column()}
					</Table.Head>
					<Table.Head
						class="h-10 w-[120px] py-2.5 text-xs font-semibold uppercase tracking-wider text-muted-foreground/90"
					>
						{m.datasets_fields_column()}
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
					<Table.Head
						class="h-10 w-[100px] py-2.5 text-right text-xs font-semibold uppercase tracking-wider text-muted-foreground/90"
					>
						{m.datasets_actions_column()}
					</Table.Head>
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#if datasetsQuery.isLoading}
					<Table.Row>
						<Table.Cell colspan={DATASET_LIST_COLUMN_COUNT} class="h-56 p-0">
							<div class="flex h-56 w-full items-center justify-center">
								<Spinner class="size-18 text-foreground dark:text-blue-500" />
							</div>
						</Table.Cell>
					</Table.Row>
				{:else if datasetsQuery.isError}
					<Table.Row>
						<Table.Cell colspan={DATASET_LIST_COLUMN_COUNT} class="h-24">
							<div
								class="mx-auto flex max-w-xl flex-wrap items-center justify-center gap-3 text-center text-sm text-destructive"
							>
								<span>{datasetsQuery.error.message}</span>
								<Button
									type="button"
									variant="outline"
									size="sm"
									onclick={() => datasetsQuery.refetch()}
								>
									{m.datasets_retry()}
								</Button>
							</div>
						</Table.Cell>
					</Table.Row>
				{:else if datasets.length > 0}
					{#each datasets as dataset (dataset.id)}
						<Table.Row class="transition-colors duration-150 hover:bg-muted/40">
							<Table.Cell class="max-w-[460px] px-4 py-3">
								<a
									href={datasetHref(dataset.id)}
									class="flex min-w-0 items-center gap-2 text-sm font-medium hover:underline"
									title={dataset.name}
								>
									<DatabaseIcon class="size-4 shrink-0 text-muted-foreground" aria-hidden="true" />
									<span class="truncate">{dataset.name}</span>
								</a>
							</Table.Cell>
							<Table.Cell class="max-w-[280px] px-4 py-3 text-sm" title={dataset.schema_name}>
								<span class="block truncate">{dataset.schema_name}</span>
							</Table.Cell>
							<Table.Cell class="px-4 py-3 text-sm text-muted-foreground">
								{fieldCountLabel(dataset.field_count)}
							</Table.Cell>
							<Table.Cell class="whitespace-nowrap px-4 py-3 text-sm text-muted-foreground">
								{formatDatasetDate(dataset.created_at, m.datasets_invalid_date())}
							</Table.Cell>
							<Table.Cell class="px-4 py-3">
								<div class="flex justify-end">
									<Button href={datasetHref(dataset.id)} variant="ghost" size="sm" class="h-8 text-xs">
										{m.datasets_open()}
										<ExternalLinkIcon class="size-3.5" aria-hidden="true" />
									</Button>
								</div>
							</Table.Cell>
						</Table.Row>
					{/each}
				{:else if datasetsQuery.isSuccess && datasets.length === 0}
					<Table.Row>
						<Table.Cell colspan={DATASET_LIST_COLUMN_COUNT} class="h-56 text-center">
							<div class="mx-auto flex max-w-[340px] flex-col items-center justify-center gap-2 p-6">
								<div class="rounded-full bg-muted/50 p-3.5 text-muted-foreground/80">
									<SearchIcon class="size-6" aria-hidden="true" />
								</div>
								<h3 class="mt-2 text-sm font-semibold text-foreground">
									{m.datasets_no_datasets_found()}
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
			{#if datasets.length > 0}
				{datasetCountSummary(datasets.length)}
			{:else}
				{m.datasets_no_datasets_to_show()}
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
					disabled={cursorState.history.length === 0 || datasetsQuery.isFetching}
					onclick={goPrevious}
					aria-label={m.datasets_previous_page()}
				>
					<ChevronLeftIcon class="size-4" aria-hidden="true" />
				</Button>
				<Button
					type="button"
					variant="outline"
					size="icon-sm"
					disabled={!nextCursor || datasetsQuery.isFetching}
					onclick={goNext}
					aria-label={m.datasets_next_page()}
				>
					<ChevronRightIcon class="size-4" aria-hidden="true" />
				</Button>
			</div>
		</div>
	</div>
</div>
