<script lang="ts">
	import AlertTriangleIcon from "@lucide/svelte/icons/alert-triangle";
	import EditIcon from "@lucide/svelte/icons/square-pen";
	import PlusIcon from "@lucide/svelte/icons/plus";
	import RefreshCwIcon from "@lucide/svelte/icons/refresh-cw";
	import TrashIcon from "@lucide/svelte/icons/trash";
	import { createMutation, createQuery, useQueryClient } from "@tanstack/svelte-query";

	import { Button } from "$lib/components/ui/button/index.js";
	import * as Dialog from "$lib/components/ui/dialog/index.js";
	import { Input } from "$lib/components/ui/input/index.js";
	import { Spinner } from "$lib/components/ui/spinner/index.js";
	import * as Table from "$lib/components/ui/table/index.js";
	import { m } from "$lib/paraglide/messages.js";
	import {
		createAdminJSONRecipeCategory,
		deleteAdminJSONRecipeCategory,
		fetchAdminJSONRecipeCategories,
		type AdminJSONRecipeCategoryResponse,
		updateAdminJSONRecipeCategory
	} from "../api";

	const queryClient = useQueryClient();

	let titleEn = $state("");
	let titleRo = $state("");
	let editTitleEn = $state("");
	let editTitleRo = $state("");
	let editingCategory = $state<AdminJSONRecipeCategoryResponse | null>(null);
	let categoryPendingDelete = $state<AdminJSONRecipeCategoryResponse | null>(null);
	let feedback = $state<string | null>(null);
	let formError = $state<string | null>(null);

	const categoriesQuery = createQuery(() => ({
		queryKey: ["admin-json-recipe-categories"],
		queryFn: () => fetchAdminJSONRecipeCategories(fetch)
	}));
	const createMutationState = createMutation(() => ({
		mutationFn: () =>
			createAdminJSONRecipeCategory(fetch, {
				title: { en: titleEn.trim(), ro: titleRo.trim() }
			}),
		onSuccess: async (category) => {
			titleEn = "";
			titleRo = "";
			formError = null;
			feedback = m.json_recipe_categories_created_success({ name: category.title.en });
			await queryClient.invalidateQueries({ queryKey: ["admin-json-recipe-categories"] });
		}
	}));
	const updateMutation = createMutation(() => ({
		mutationFn: (category: AdminJSONRecipeCategoryResponse) =>
			updateAdminJSONRecipeCategory(fetch, category.id, {
				title: { en: editTitleEn.trim(), ro: editTitleRo.trim() }
			}),
		onSuccess: async (category) => {
			editingCategory = null;
			feedback = m.json_recipe_categories_saved_success({ name: category.title.en });
			await queryClient.invalidateQueries({ queryKey: ["admin-json-recipe-categories"] });
		}
	}));
	const deleteMutation = createMutation(() => ({
		mutationFn: (category: AdminJSONRecipeCategoryResponse) =>
			deleteAdminJSONRecipeCategory(fetch, category.id),
		onSuccess: async (_result, category) => {
			categoryPendingDelete = null;
			feedback = m.json_recipe_categories_deleted_success({ name: category.title.en });
			await queryClient.invalidateQueries({ queryKey: ["admin-json-recipe-categories"] });
		}
	}));

	const categories = $derived(categoriesQuery.data?.categories ?? []);

	function validateTitles(en: string, ro: string) {
		if (!en.trim() || !ro.trim()) return m.json_recipe_categories_validation_titles_required();
		if (Array.from(en.trim()).length > 160 || Array.from(ro.trim()).length > 160) {
			return m.json_recipe_categories_validation_titles_too_long();
		}
		return null;
	}

	async function createCategory() {
		formError = validateTitles(titleEn, titleRo);
		feedback = null;
		if (formError) return;
		try {
			await createMutationState.mutateAsync();
		} catch {
			// The mutation exposes the error below.
		}
	}

	function openEditDialog(category: AdminJSONRecipeCategoryResponse) {
		editTitleEn = category.title.en;
		editTitleRo = category.title.ro;
		editingCategory = category;
		updateMutation.reset();
	}

	function setEditDialogOpen(open: boolean) {
		if (!open && !updateMutation.isPending) editingCategory = null;
	}

	async function saveEditedCategory() {
		const category = editingCategory;
		if (!category) return;
		formError = validateTitles(editTitleEn, editTitleRo);
		feedback = null;
		if (formError) return;
		try {
			await updateMutation.mutateAsync(category);
		} catch {
			// The mutation exposes the error below.
		}
	}

	function openDeleteDialog(category: AdminJSONRecipeCategoryResponse) {
		deleteMutation.reset();
		categoryPendingDelete = category;
	}

	function setDeleteDialogOpen(open: boolean) {
		if (!open && !deleteMutation.isPending) categoryPendingDelete = null;
	}

	async function deleteSelectedCategory() {
		const category = categoryPendingDelete;
		if (!category) return;
		feedback = null;
		try {
			await deleteMutation.mutateAsync(category);
		} catch {
			// The mutation exposes the error below.
		}
	}
</script>

<div class="flex flex-1 flex-col gap-6 p-4 lg:p-6">
	<div class="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
		<div class="space-y-1">
			<h1 class="text-2xl font-semibold tracking-tight">{m.json_recipe_categories_title()}</h1>
			<p class="text-sm text-muted-foreground">{m.json_recipe_categories_description()}</p>
		</div>
		<div class="flex items-center gap-2">
			<Button variant="outline" size="icon" onclick={() => categoriesQuery.refetch()} aria-label={m.common_retry()}>
				<RefreshCwIcon class="size-4" />
			</Button>
			<Button href="/admin-portal/json-recipes" variant="outline">{m.json_recipes_view_recipes()}</Button>
		</div>
	</div>

	{#if feedback}
		<div role="status" class="rounded-lg border border-emerald-500/20 bg-emerald-500/10 px-4 py-3 text-sm font-medium text-emerald-700 dark:text-emerald-300">
			{feedback}
		</div>
	{/if}

	<div class="rounded-xl border bg-background p-4">
		<div class="grid gap-3 md:grid-cols-[minmax(0,1fr)_minmax(0,1fr)_auto] md:items-end">
			<label class="grid gap-2 text-sm font-medium">
				<span>{m.json_recipe_categories_title_en_label()}</span>
				<Input value={titleEn} oninput={(event) => (titleEn = event.currentTarget.value)} />
			</label>
			<label class="grid gap-2 text-sm font-medium">
				<span>{m.json_recipe_categories_title_ro_label()}</span>
				<Input value={titleRo} oninput={(event) => (titleRo = event.currentTarget.value)} />
			</label>
			<Button type="button" disabled={createMutationState.isPending} onclick={createCategory}>
				{#if createMutationState.isPending}
					<Spinner class="size-4" />
				{:else}
					<PlusIcon class="size-4" />
				{/if}
				{m.json_recipe_categories_create_category()}
			</Button>
		</div>
		{#if formError}
			<p class="mt-3 text-sm text-destructive">{formError}</p>
		{/if}
		{#if createMutationState.isError}
			<p class="mt-3 text-sm text-destructive">{createMutationState.error.message}</p>
		{/if}
	</div>

	{#if categoriesQuery.isLoading}
		<div class="flex min-h-[240px] items-center justify-center rounded-xl border bg-background">
			<div class="flex items-center gap-2 text-sm text-muted-foreground">
				<Spinner class="size-4" />
				{m.json_recipe_categories_loading()}
			</div>
		</div>
	{:else if categoriesQuery.isError}
		<div class="rounded-xl border border-destructive/20 bg-destructive/5 p-5 text-sm text-destructive">
			<div class="mb-3 flex items-center gap-2 font-semibold">
				<AlertTriangleIcon class="size-4" />
				{m.json_recipe_categories_could_not_load()}
			</div>
			<p>{categoriesQuery.error?.message}</p>
		</div>
	{:else if categories.length === 0}
		<div class="flex min-h-[240px] flex-col items-center justify-center rounded-xl border bg-background px-6 text-center">
			<h2 class="text-base font-semibold">{m.json_recipe_categories_empty_title()}</h2>
			<p class="mt-1 max-w-md text-sm text-muted-foreground">{m.json_recipe_categories_empty_body()}</p>
		</div>
	{:else}
		<div class="overflow-hidden rounded-xl border bg-background">
			<Table.Root>
				<Table.Header>
					<Table.Row>
						<Table.Head>{m.json_recipe_categories_title_en_label()}</Table.Head>
						<Table.Head>{m.json_recipe_categories_title_ro_label()}</Table.Head>
						<Table.Head class="w-[128px] text-right">{m.common_actions()}</Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					{#each categories as category (category.id)}
						<Table.Row>
							<Table.Cell class="font-medium">{category.title.en}</Table.Cell>
							<Table.Cell>{category.title.ro}</Table.Cell>
							<Table.Cell>
								<div class="flex justify-end gap-1">
									<Button type="button" variant="ghost" size="icon" onclick={() => openEditDialog(category)} aria-label={m.json_recipe_categories_edit_aria({ name: category.title.en })}>
										<EditIcon class="size-4" />
									</Button>
									<Button type="button" variant="ghost" size="icon" onclick={() => openDeleteDialog(category)} aria-label={m.json_recipe_categories_delete_aria({ name: category.title.en })}>
										<TrashIcon class="size-4 text-destructive" />
									</Button>
								</div>
							</Table.Cell>
						</Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
		</div>
	{/if}
</div>

<Dialog.Root bind:open={() => Boolean(editingCategory), setEditDialogOpen}>
	<Dialog.Content class="sm:max-w-md">
		{#if editingCategory}
			<Dialog.Header>
				<Dialog.Title>{m.json_recipe_categories_edit_title()}</Dialog.Title>
			</Dialog.Header>
			<div class="grid gap-3">
				<label class="grid gap-2 text-sm font-medium">
					<span>{m.json_recipe_categories_title_en_label()}</span>
					<Input value={editTitleEn} oninput={(event) => (editTitleEn = event.currentTarget.value)} />
				</label>
				<label class="grid gap-2 text-sm font-medium">
					<span>{m.json_recipe_categories_title_ro_label()}</span>
					<Input value={editTitleRo} oninput={(event) => (editTitleRo = event.currentTarget.value)} />
				</label>
				{#if formError}
					<p class="text-sm text-destructive">{formError}</p>
				{/if}
				{#if updateMutation.isError}
					<p class="text-sm text-destructive">{updateMutation.error.message}</p>
				{/if}
			</div>
			<Dialog.Footer>
				<Button type="button" variant="outline" disabled={updateMutation.isPending} onclick={() => setEditDialogOpen(false)}>
					{m.common_cancel()}
				</Button>
				<Button type="button" disabled={updateMutation.isPending} onclick={saveEditedCategory}>
					{#if updateMutation.isPending}
						<Spinner class="size-4" />
					{/if}
					{m.json_recipe_categories_save_category()}
				</Button>
			</Dialog.Footer>
		{/if}
	</Dialog.Content>
</Dialog.Root>

<Dialog.Root bind:open={() => Boolean(categoryPendingDelete), setDeleteDialogOpen}>
	<Dialog.Content class="sm:max-w-md">
		{#if categoryPendingDelete}
			<Dialog.Header>
				<Dialog.Title>{m.json_recipe_categories_delete_confirm()}</Dialog.Title>
				<Dialog.Description>
					<span class="font-medium text-foreground">{categoryPendingDelete.title.en}</span>
				</Dialog.Description>
			</Dialog.Header>
			{#if deleteMutation.isError}
				<div role="alert" class="rounded-lg border border-destructive/20 bg-destructive/5 px-3 py-2 text-sm text-destructive">
					{deleteMutation.error.message}
				</div>
			{/if}
			<Dialog.Footer>
				<Button type="button" variant="outline" disabled={deleteMutation.isPending} onclick={() => setDeleteDialogOpen(false)}>
					{m.common_cancel()}
				</Button>
				<Button type="button" variant="destructive" disabled={deleteMutation.isPending} onclick={deleteSelectedCategory}>
					{#if deleteMutation.isPending}
						<Spinner class="size-4" />
					{/if}
					{m.common_delete()}
				</Button>
			</Dialog.Footer>
		{/if}
	</Dialog.Content>
</Dialog.Root>
