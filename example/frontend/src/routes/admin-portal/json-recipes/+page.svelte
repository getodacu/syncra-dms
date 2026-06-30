<script lang="ts">
	import AlertTriangleIcon from "@lucide/svelte/icons/alert-triangle";
	import ArrowDownIcon from "@lucide/svelte/icons/arrow-down";
	import ArrowUpIcon from "@lucide/svelte/icons/arrow-up";
	import EditIcon from "@lucide/svelte/icons/square-pen";
	import FileJsonIcon from "@lucide/svelte/icons/file-json";
	import PlusIcon from "@lucide/svelte/icons/plus";
	import RefreshCwIcon from "@lucide/svelte/icons/refresh-cw";
	import TrashIcon from "@lucide/svelte/icons/trash";
	import { createMutation, createQuery, useQueryClient } from "@tanstack/svelte-query";

	import { Button } from "$lib/components/ui/button/index.js";
	import * as Dialog from "$lib/components/ui/dialog/index.js";
	import { Spinner } from "$lib/components/ui/spinner/index.js";
	import * as Table from "$lib/components/ui/table/index.js";
	import { m } from "$lib/paraglide/messages.js";
	import {
		cursorNextState,
		cursorPreviousState,
		resetCursorState,
		type CursorState,
		type SortDirection
	} from "../../app/billing/orders/table-utils";
	import {
		deleteAdminJSONRecipe,
		fetchAdminJSONRecipeCategories,
		fetchAdminJSONRecipes,
		type AdminJSONRecipeResponse
	} from "./api";

	const PAGE_SIZE_OPTIONS = [10, 20, 50, 100];

	let sortDirection = $state<SortDirection>("desc");
	let pageSize = $state(20);
	let cursorState = $state<CursorState>(resetCursorState());
	let feedback = $state<string | null>(null);
	let recipePendingDelete = $state<AdminJSONRecipeResponse | null>(null);

	const queryClient = useQueryClient();
	const recipesQuery = createQuery(() => ({
		queryKey: [
			"admin-json-recipes",
			{
				cursor: cursorState.currentCursor,
				size: pageSize,
				sort: sortDirection
			}
		],
		queryFn: () =>
			fetchAdminJSONRecipes(fetch, {
				cursor: cursorState.currentCursor,
				size: pageSize,
				sort: sortDirection
			})
	}));
	const categoriesQuery = createQuery(() => ({
		queryKey: ["admin-json-recipe-categories"],
		queryFn: () => fetchAdminJSONRecipeCategories(fetch)
	}));
	const deleteMutation = createMutation(() => ({
		mutationFn: (recipe: AdminJSONRecipeResponse) => deleteAdminJSONRecipe(fetch, recipe.id),
		onSuccess: async (_result, recipe) => {
			feedback = m.json_recipes_deleted_success({ name: recipe.title });
			await queryClient.invalidateQueries({ queryKey: ["admin-json-recipes"] });
		}
	}));

	const recipes = $derived(recipesQuery.data?.recipes ?? []);
	const categories = $derived(categoriesQuery.data?.categories ?? []);
	const categoryMap = $derived(new Map(categories.map((category) => [category.id, category])));
	const nextCursor = $derived(recipesQuery.data?.next_cursor ?? null);
	const deletingId = $derived(deleteMutation.variables?.id ?? null);

	function resetPagination() {
		cursorState = resetCursorState();
	}

	function setSortDirection(next: SortDirection) {
		sortDirection = next;
		resetPagination();
	}

	function setPageSize(next: number) {
		pageSize = next;
		resetPagination();
	}

	function nextPage() {
		if (!nextCursor) return;
		cursorState = cursorNextState(cursorState, nextCursor);
	}

	function previousPage() {
		cursorState = cursorPreviousState(cursorState);
	}

	function formatDate(value: string) {
		const date = new Date(value);
		if (Number.isNaN(date.getTime())) return value;
		return new Intl.DateTimeFormat(undefined, {
			dateStyle: "medium",
			timeStyle: "short"
		}).format(date);
	}

	function fieldCount(recipe: AdminJSONRecipeResponse) {
		const properties = recipe.json.properties;
		if (properties && typeof properties === "object" && !Array.isArray(properties)) {
			return Object.keys(properties).length;
		}
		return 0;
	}

	function categoryLabel(recipe: AdminJSONRecipeResponse) {
		const category = recipe.category ?? (recipe.category_id ? categoryMap.get(recipe.category_id) : null);
		return category ? `${category.title.en} / ${category.title.ro}` : m.json_recipes_others();
	}

	function setDeleteDialogOpen(open: boolean) {
		if (!open && !deleteMutation.isPending) {
			recipePendingDelete = null;
		}
	}

	function openDeleteDialog(recipe: AdminJSONRecipeResponse) {
		deleteMutation.reset();
		recipePendingDelete = recipe;
	}

	async function deleteSelectedRecipe() {
		const recipe = recipePendingDelete;
		if (!recipe) return;

		feedback = null;
		try {
			await deleteMutation.mutateAsync(recipe);
			recipePendingDelete = null;
		} catch {
			// The mutation stores the error for the dialog to render.
		}
	}
</script>

<div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
	<div class="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
		<div class="space-y-1">
			<h1 class="text-2xl font-semibold tracking-tight">{m.json_recipes_title()}</h1>
			<p class="text-sm text-muted-foreground">{m.json_recipes_description()}</p>
		</div>
		<div class="flex items-center gap-2">
			<Button variant="outline" size="icon" onclick={() => recipesQuery.refetch()} aria-label={m.common_retry()}>
				<RefreshCwIcon class="size-4" />
			</Button>
			<Button href="/admin-portal/json-recipes/categories" variant="outline">
				{m.json_recipes_manage_categories()}
			</Button>
			<Button href="/admin-portal/json-recipes/new">
				<PlusIcon class="mr-2 size-4" />
				{m.json_recipes_new_recipe()}
			</Button>
		</div>
	</div>

	{#if feedback}
		<div role="status" class="rounded-lg border border-emerald-500/20 bg-emerald-500/10 px-4 py-3 text-sm font-medium text-emerald-700 dark:text-emerald-300">
			{feedback}
		</div>
	{/if}

	{#if recipesQuery.isLoading}
		<div class="flex min-h-[320px] items-center justify-center rounded-xl border bg-background">
			<div class="flex items-center gap-2 text-sm text-muted-foreground">
				<Spinner class="size-4" />
				{m.json_recipes_loading()}
			</div>
		</div>
	{:else if recipesQuery.isError}
		<div class="rounded-xl border border-destructive/20 bg-destructive/5 p-5 text-sm text-destructive">
			<div class="mb-3 flex items-center gap-2 font-semibold">
				<AlertTriangleIcon class="size-4" />
				{m.json_recipes_could_not_load()}
			</div>
			<p>{recipesQuery.error?.message}</p>
		</div>
	{:else if recipes.length === 0}
		<div class="flex min-h-[320px] flex-col items-center justify-center rounded-xl border bg-background px-6 text-center">
			<FileJsonIcon class="mb-4 size-10 text-muted-foreground" />
			<h2 class="text-base font-semibold">{m.json_recipes_no_recipes_found()}</h2>
			<p class="mt-1 max-w-md text-sm text-muted-foreground">{m.json_recipes_empty_body()}</p>
			<Button href="/admin-portal/json-recipes/new" class="mt-5">
				<PlusIcon class="mr-2 size-4" />
				{m.json_recipes_new_recipe()}
			</Button>
		</div>
	{:else}
		<div class="overflow-hidden rounded-xl border bg-background">
			<Table.Root>
				<Table.Header>
					<Table.Row>
						<Table.Head>{m.schemas_name_column()}</Table.Head>
						<Table.Head>{m.schemas_description_label()}</Table.Head>
						<Table.Head class="w-[180px]">{m.json_recipes_category_label()}</Table.Head>
						<Table.Head class="w-[110px] text-right">{m.json_recipes_json_fields_column()}</Table.Head>
						<Table.Head class="w-[110px] text-right">{m.json_recipes_counter_column()}</Table.Head>
						<Table.Head class="w-[190px]">
							<Button
								type="button"
								variant="ghost"
								class="h-8 px-2"
								onclick={() => setSortDirection(sortDirection === "asc" ? "desc" : "asc")}
								aria-label={sortDirection === "asc" ? m.json_recipes_sort_created_descending() : m.json_recipes_sort_created_ascending()}
							>
								{m.json_recipes_created_column()}
								{#if sortDirection === "asc"}
									<ArrowUpIcon class="ml-2 size-3.5" />
								{:else}
									<ArrowDownIcon class="ml-2 size-3.5" />
								{/if}
							</Button>
						</Table.Head>
						<Table.Head class="w-[190px]">{m.json_recipes_updated_column()}</Table.Head>
						<Table.Head class="w-[128px] text-right">{m.common_actions()}</Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					{#each recipes as recipe (recipe.id)}
						<Table.Row>
							<Table.Cell class="font-medium">
								<a class="hover:underline" href="/admin-portal/json-recipes/{recipe.id}">{recipe.title}</a>
							</Table.Cell>
							<Table.Cell class="max-w-[420px] truncate text-muted-foreground">
								{recipe.description || m.schemas_no_description()}
							</Table.Cell>
							<Table.Cell class="max-w-[180px] truncate text-muted-foreground">
								{categoryLabel(recipe)}
							</Table.Cell>
							<Table.Cell class="text-right tabular-nums">{fieldCount(recipe)}</Table.Cell>
							<Table.Cell class="text-right tabular-nums">{recipe.counter}</Table.Cell>
							<Table.Cell>{formatDate(recipe.created_at)}</Table.Cell>
							<Table.Cell>{formatDate(recipe.updated_at)}</Table.Cell>
							<Table.Cell>
								<div class="flex justify-end gap-1">
									<Button
										href="/admin-portal/json-recipes/{recipe.id}"
										variant="ghost"
										size="icon"
										aria-label={m.json_recipes_edit_aria({ name: recipe.title })}
									>
										<EditIcon class="size-4" />
									</Button>
									<Button
										type="button"
										variant="ghost"
										size="icon"
										disabled={deleteMutation.isPending && deletingId === recipe.id}
										aria-label={m.json_recipes_delete_aria({ name: recipe.title })}
										onclick={() => openDeleteDialog(recipe)}
									>
										{#if deleteMutation.isPending && deletingId === recipe.id}
											<Spinner class="size-4" />
										{:else}
											<TrashIcon class="size-4 text-destructive" />
										{/if}
									</Button>
								</div>
							</Table.Cell>
						</Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
		</div>

		<div class="flex flex-col gap-3 text-sm text-muted-foreground sm:flex-row sm:items-center sm:justify-between">
			<div>
				{#if recipes.length === 0}
					{m.json_recipes_no_recipes_to_show()}
				{:else}
					{recipes.length === 1
						? m.json_recipes_showing_one({ count: recipes.length })
						: m.json_recipes_showing_other({ count: recipes.length })}
				{/if}
			</div>
			<div class="flex flex-wrap items-center gap-2">
				<span>{m.common_rows_per_page()}</span>
				<select
					class="h-9 rounded-md border bg-background px-2 text-sm"
					value={pageSize}
					onchange={(event) => setPageSize(Number(event.currentTarget.value))}
				>
					{#each PAGE_SIZE_OPTIONS as option (option)}
						<option value={option}>{option}</option>
					{/each}
				</select>
				<Button variant="outline" onclick={previousPage} disabled={cursorState.history.length === 0}>
					{m.common_previous()}
				</Button>
				<Button variant="outline" onclick={nextPage} disabled={!nextCursor}>
					{m.common_next()}
				</Button>
			</div>
		</div>
	{/if}
</div>

<Dialog.Root bind:open={() => Boolean(recipePendingDelete), setDeleteDialogOpen}>
	<Dialog.Content class="sm:max-w-md">
		{#if recipePendingDelete}
			<Dialog.Header>
				<Dialog.Title>{m.json_recipes_delete_confirm()}</Dialog.Title>
				<Dialog.Description>
					<span class="font-medium text-foreground">{recipePendingDelete.title}</span>
				</Dialog.Description>
			</Dialog.Header>

			{#if deleteMutation.isError}
				<div
					role="alert"
					class="flex items-center gap-2 rounded-lg border border-destructive/20 bg-destructive/5 px-3 py-2 text-sm text-destructive"
				>
					<AlertTriangleIcon class="size-4 shrink-0" />
					<span>{deleteMutation.error.message}</span>
				</div>
			{/if}

			<Dialog.Footer>
				<Button
					type="button"
					variant="outline"
					disabled={deleteMutation.isPending}
					onclick={() => setDeleteDialogOpen(false)}
				>
					{m.common_cancel()}
				</Button>
				<Button
					type="button"
					variant="destructive"
					disabled={deleteMutation.isPending}
					onclick={deleteSelectedRecipe}
				>
					{#if deleteMutation.isPending}
						<Spinner class="size-4" />
						{m.json_recipes_deleting()}
					{:else}
						{m.common_delete()}
					{/if}
				</Button>
			</Dialog.Footer>
		{/if}
	</Dialog.Content>
</Dialog.Root>
