<script lang="ts">
	import { goto } from "$app/navigation";
	import { createQuery } from "@tanstack/svelte-query";
	import { m } from "$lib/paraglide/messages.js";
	import { createAdminJSONRecipe, fetchAdminJSONRecipeCategories } from "../api";
	import JsonRecipeEditor, {
		type JsonRecipeEditorSubmitInput
	} from "../json-recipe-editor.svelte";

	let pending = $state(false);
	let serverError = $state<string | null>(null);
	let successMessage = $state<string | null>(null);

	const categoriesQuery = createQuery(() => ({
		queryKey: ["admin-json-recipe-categories"],
		queryFn: () => fetchAdminJSONRecipeCategories(fetch)
	}));
	const categories = $derived(categoriesQuery.data?.categories ?? []);

	async function submit(input: JsonRecipeEditorSubmitInput) {
		pending = true;
		serverError = null;
		successMessage = null;
		try {
			const recipe = await createAdminJSONRecipe(fetch, input);
			successMessage = m.json_recipes_created_success({ name: recipe.title });
			await goto(`/admin-portal/json-recipes/${recipe.id}`);
		} catch (error) {
			serverError = error instanceof Error ? error.message : m.json_recipes_could_not_load();
		} finally {
			pending = false;
		}
	}
</script>

<JsonRecipeEditor
	title={m.json_recipes_new_title()}
	descriptionText={m.json_recipes_new_description()}
	submitLabel={m.json_recipes_save_recipe()}
	{pending}
	{serverError}
	{successMessage}
	{categories}
	onSubmit={submit}
/>
