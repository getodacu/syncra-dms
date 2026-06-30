<script lang="ts">
	import { goto } from "$app/navigation";
	import AlertTriangleIcon from "@lucide/svelte/icons/alert-triangle";
	import FileJsonIcon from "@lucide/svelte/icons/file-json";
	import { createQuery } from "@tanstack/svelte-query";

	import { Button } from "$lib/components/ui/button/index.js";
	import * as Dialog from "$lib/components/ui/dialog/index.js";
	import { Spinner } from "$lib/components/ui/spinner/index.js";
	import { m } from "$lib/paraglide/messages.js";
	import {
		deleteAdminJSONRecipe,
		fetchAdminJSONRecipeCategories,
		fetchAdminJSONRecipe,
		updateAdminJSONRecipe
	} from "../api";
	import JsonRecipeEditor, {
		type JsonRecipeEditorSubmitInput
	} from "../json-recipe-editor.svelte";

	let { params } = $props<{ params: { id: string } }>();

	let pending = $state(false);
	let deletePending = $state(false);
	let deleteDialogOpen = $state(false);
	let serverError = $state<string | null>(null);
	let successMessage = $state<string | null>(null);

	const recipeQuery = createQuery(() => ({
		queryKey: ["admin-json-recipe", params.id],
		queryFn: () => fetchAdminJSONRecipe(fetch, params.id)
	}));
	const categoriesQuery = createQuery(() => ({
		queryKey: ["admin-json-recipe-categories"],
		queryFn: () => fetchAdminJSONRecipeCategories(fetch)
	}));

	const recipe = $derived(recipeQuery.data);
	const categories = $derived(categoriesQuery.data?.categories ?? []);

	function clearFeedback() {
		serverError = null;
		successMessage = null;
	}

	function setDeleteDialogOpen(open: boolean) {
		if (!open && deletePending) return;
		deleteDialogOpen = open;
	}

	async function submit(input: JsonRecipeEditorSubmitInput) {
		pending = true;
		clearFeedback();
		try {
			const updated = await updateAdminJSONRecipe(fetch, params.id, input);
			successMessage = m.json_recipes_saved_success({ name: updated.title });
			await recipeQuery.refetch();
		} catch (error) {
			serverError = error instanceof Error ? error.message : m.json_recipes_could_not_load();
		} finally {
			pending = false;
		}
	}

	async function deleteRecipe() {
		if (!recipe) return;
		deletePending = true;
		clearFeedback();
		try {
			await deleteAdminJSONRecipe(fetch, recipe.id);
			deleteDialogOpen = false;
			await goto("/admin-portal/json-recipes");
		} catch (error) {
			serverError = error instanceof Error ? error.message : m.json_recipes_could_not_load();
		} finally {
			deletePending = false;
		}
	}
</script>

{#if recipeQuery.isLoading}
	<div class="flex min-h-[420px] items-center justify-center p-6">
		<div class="flex items-center gap-2 text-sm text-muted-foreground">
			<Spinner class="size-4" />
			{m.json_recipes_loading_recipe()}
		</div>
	</div>
{:else if recipeQuery.isError}
	<div class="flex min-h-[420px] items-center justify-center p-6">
		<div class="max-w-md rounded-xl border border-destructive/20 bg-destructive/5 p-6 text-center">
			<AlertTriangleIcon class="mx-auto mb-4 size-8 text-destructive" />
			<h1 class="text-lg font-semibold">{m.json_recipes_not_found_title()}</h1>
			<p class="mt-2 text-sm text-muted-foreground">{m.json_recipes_not_found_body()}</p>
			<p class="mt-3 text-sm text-destructive">{m.json_recipes_could_not_load()}</p>
			<Button href="/admin-portal/json-recipes" variant="outline" class="mt-5">
				{m.json_recipes_view_recipes()}
			</Button>
		</div>
	</div>
{:else if recipe}
	<JsonRecipeEditor
		title={m.json_recipes_edit_title()}
		descriptionText={m.json_recipes_edit_description()}
		submitLabel={m.json_recipes_save_changes()}
		initial={{
			title: recipe.title,
			description: recipe.description,
			json: recipe.json,
			category_id: recipe.category_id
		}}
		{categories}
		{pending}
		{deletePending}
		{serverError}
		{successMessage}
		onSubmit={submit}
		onDelete={() => setDeleteDialogOpen(true)}
		onDirty={clearFeedback}
	/>

	<Dialog.Root bind:open={() => deleteDialogOpen, setDeleteDialogOpen}>
		<Dialog.Content class="sm:max-w-md">
			<Dialog.Header>
				<Dialog.Title>{m.json_recipes_delete_confirm()}</Dialog.Title>
				<Dialog.Description>
					<span class="font-medium text-foreground">{recipe.title}</span>
				</Dialog.Description>
			</Dialog.Header>

			{#if serverError}
				<div
					role="alert"
					class="flex items-center gap-2 rounded-lg border border-destructive/20 bg-destructive/5 px-3 py-2 text-sm text-destructive"
				>
					<AlertTriangleIcon class="size-4 shrink-0" />
					<span>{serverError}</span>
				</div>
			{/if}

			<Dialog.Footer>
				<Button
					type="button"
					variant="outline"
					disabled={deletePending}
					onclick={() => setDeleteDialogOpen(false)}
				>
					{m.common_cancel()}
				</Button>
				<Button
					type="button"
					variant="destructive"
					disabled={deletePending}
					onclick={deleteRecipe}
				>
					{#if deletePending}
						<Spinner class="size-4" />
						{m.json_recipes_deleting()}
					{:else}
						{m.common_delete()}
					{/if}
				</Button>
			</Dialog.Footer>
		</Dialog.Content>
	</Dialog.Root>
{:else}
	<div class="flex min-h-[420px] items-center justify-center p-6">
		<div class="max-w-md rounded-xl border bg-background p-6 text-center">
			<FileJsonIcon class="mx-auto mb-4 size-8 text-muted-foreground" />
			<h1 class="text-lg font-semibold">{m.json_recipes_not_found_title()}</h1>
			<p class="mt-2 text-sm text-muted-foreground">{m.json_recipes_not_found_body()}</p>
			<Button href="/admin-portal/json-recipes" variant="outline" class="mt-5">
				{m.json_recipes_view_recipes()}
			</Button>
		</div>
	</div>
{/if}
