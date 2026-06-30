<script lang="ts">
	import { goto } from "$app/navigation";
	import CopyPlusIcon from "@lucide/svelte/icons/copy-plus";
	import LoaderIcon from "@lucide/svelte/icons/loader-circle";
	import SearchIcon from "@lucide/svelte/icons/search";

	import * as Alert from "$lib/components/ui/alert/index.js";
	import { Button } from "$lib/components/ui/button/index.js";
	import * as Card from "$lib/components/ui/card/index.js";
	import { Input } from "$lib/components/ui/input/index.js";
	import { Spinner } from "$lib/components/ui/spinner/index.js";
	import { m } from "$lib/paraglide/messages.js";
	import { getLocale } from "$lib/paraglide/runtime.js";
	import { deployOCRRecipe, fetchOCRRecipes, type OCRRecipeResponse } from "./api";
	import {
		ALL_CATEGORY_KEY,
		buildCategoryFilterOptions,
		filterRecipesByCategory,
		type CategoryFilterOption
	} from "./page-state";
	import type { RecipeFieldSummary } from "./recipe-summary";
	import { normalizeSearchText, summarizeRecipe } from "./recipe-summary";

	type SortMode = "popular" | "newest" | "az";

	let {
		recipes: initialRecipes,
		loadError: initialLoadError,
		isLoggedIn,
		userId
	}: {
		recipes: OCRRecipeResponse[];
		loadError: string | null;
		isLoggedIn: boolean;
		userId: string | null;
	} = $props();

	let refreshedRecipes = $state<OCRRecipeResponse[] | null>(null);
	let loadErrorOverride = $state<string | null | undefined>(undefined);
	let search = $state("");
	let sortMode = $state<SortMode>("popular");
	let selectedCategoryKey = $state(ALL_CATEGORY_KEY);
	let cloningId = $state<string | null>(null);
	let cloneError = $state<string | null>(null);
	let refreshing = $state(false);

	const recipes = $derived(refreshedRecipes ?? initialRecipes);
	const loadError = $derived(loadErrorOverride === undefined ? initialLoadError : loadErrorOverride);
	const recipeCards = $derived(
		recipes.map((recipe) => ({
			recipe,
			summary: summarizeRecipe(recipe)
		}))
	);
	const categoryFilterOptions = $derived(
		buildCategoryFilterOptions(recipeCards, categoryTitle, m.ocr_recipes_all_categories(), getLocale())
	);
	const activeCategoryKey = $derived(
		categoryFilterOptions.some((option) => option.key === selectedCategoryKey)
			? selectedCategoryKey
			: ALL_CATEGORY_KEY
	);
	const visibleRecipes = $derived.by(() => {
		const query = normalizeSearchText(search);
		const categoryFiltered = filterRecipesByCategory(recipeCards, activeCategoryKey);
		const filtered = query
			? categoryFiltered.filter((item) => item.summary.searchText.includes(query))
			: categoryFiltered;

		return [...filtered].sort((left, right) => {
			if (sortMode === "az") return left.recipe.title.localeCompare(right.recipe.title);
			if (sortMode === "newest") {
				return dateValue(right.recipe.created_at) - dateValue(left.recipe.created_at);
			}
			return right.recipe.counter - left.recipe.counter || left.recipe.title.localeCompare(right.recipe.title);
		});
	});
	function dateValue(value: string) {
		const time = new Date(value).getTime();
		return Number.isFinite(time) ? time : 0;
	}

	function categoryTitle(recipe: OCRRecipeResponse) {
		const category = recipe.category;
		if (!category) return m.ocr_recipes_others();
		return getLocale() === "ro" ? category.title.ro : category.title.en;
	}

	function fieldCountLabel(count: number) {
		return count === 1
			? m.ocr_recipes_fields_one({ count })
			: m.ocr_recipes_fields_other({ count });
	}

	function requiredCountLabel(count: number) {
		return count === 1
			? m.ocr_recipes_required_one({ count })
			: m.ocr_recipes_required_other({ count });
	}

	function cloneCountLabel(count: number) {
		return count === 1
			? m.ocr_recipes_deploys_one({ count })
			: m.ocr_recipes_deploys_other({ count });
	}

	function showingLabel(count: number) {
		return count === 1
			? m.ocr_recipes_showing_one({ count })
			: m.ocr_recipes_showing_other({ count });
	}

	function categoryFilterClass(option: CategoryFilterOption) {
		const selected = activeCategoryKey === option.key;
		return [
			"inline-flex min-w-0 max-w-full items-center gap-2 rounded-md border px-3 py-1.5 text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-[3px] focus-visible:ring-ring/50 sm:max-w-[16rem]",
			selected
				? "border-foreground bg-foreground text-background"
				: "border-border bg-background text-muted-foreground hover:border-foreground/50 hover:text-foreground"
		].join(" ");
	}

	function categoryCountClass(option: CategoryFilterOption) {
		return activeCategoryKey === option.key
			? "shrink-0 rounded bg-background/20 px-1.5 py-0.5 text-xs tabular-nums text-background"
			: "shrink-0 rounded bg-muted px-1.5 py-0.5 text-xs tabular-nums text-muted-foreground";
	}

	function fieldBadgeClass(field: RecipeFieldSummary) {
		return [
			"inline-flex max-w-full items-center gap-1.5 rounded-md border px-2.5 py-1 text-[11px] font-medium transition-colors",
			field.required
				? "border-emerald-600/40 bg-emerald-500/10 text-emerald-700 dark:border-emerald-500/30 dark:bg-emerald-500/10 dark:text-emerald-300"
				: "border-border bg-muted/50 text-muted-foreground"
		].join(" ");
	}

	function fieldDotClass(field: RecipeFieldSummary) {
		return field.required
			? "size-1.5 shrink-0 rounded-full bg-emerald-600 dark:bg-emerald-400"
			: "size-1.5 shrink-0 rounded-full bg-muted-foreground/60";
	}

	function fieldTitle(field: RecipeFieldSummary) {
		return `${field.key} (${field.type})`;
	}

	function fieldAriaLabel(field: RecipeFieldSummary) {
		return field.required
			? `${field.key}, ${field.type}, ${m.ocr_recipes_required()}`
			: `${field.key}, ${field.type}`;
	}

	async function refreshRecipes() {
		refreshing = true;
		loadErrorOverride = null;
		try {
			const result = await fetchOCRRecipes(fetch, { size: 100, sort: "desc" });
			refreshedRecipes = result.recipes;
		} catch (error) {
			loadErrorOverride = error instanceof Error ? error.message : m.ocr_recipes_load_failed();
		} finally {
			refreshing = false;
		}
	}

	async function cloneRecipe(recipe: OCRRecipeResponse) {
		if (!userId || cloningId) return;
		cloningId = recipe.id;
		cloneError = null;
		try {
			const result = await deployOCRRecipe(fetch, recipe.id, userId);
			await goto(`/app/schemas/edit/${result.schema.id}`);
		} catch (error) {
			cloneError = error instanceof Error ? error.message : m.ocr_recipes_clone_failed();
		} finally {
			cloningId = null;
		}
	}
</script>

<div class="flex flex-col gap-6">
	<section class="grid gap-3 rounded-[8px] border border-border bg-muted/20 p-3 md:grid-cols-[1fr_180px_auto] md:items-end">
		<div class="min-w-0">
			<label for="recipe-search" class="text-xs font-semibold uppercase tracking-wide text-muted-foreground">
				{m.ocr_recipes_search_label()}
			</label>
			<div class="relative mt-2">
				<SearchIcon class="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
				<Input
					id="recipe-search"
					value={search}
					oninput={(event) => (search = event.currentTarget.value)}
					placeholder={m.ocr_recipes_search_placeholder()}
					class="pl-9"
				/>
			</div>
		</div>

		<div>
			<label for="recipe-sort" class="text-xs font-semibold uppercase tracking-wide text-muted-foreground">
				{m.ocr_recipes_sort_label()}
			</label>
			<select
				id="recipe-sort"
				value={sortMode}
				onchange={(event) => (sortMode = event.currentTarget.value as SortMode)}
				class="mt-2 h-10 w-full rounded-md border border-input bg-background px-3 text-sm shadow-xs outline-none transition-colors focus-visible:ring-[3px] focus-visible:ring-ring/50"
			>
				<option value="popular">{m.ocr_recipes_sort_popular()}</option>
				<option value="newest">{m.ocr_recipes_sort_newest()}</option>
				<option value="az">{m.ocr_recipes_sort_az()}</option>
			</select>
		</div>

		<Button
			type="button"
			variant="outline"
			class="h-10 cursor-pointer"
			disabled={refreshing}
			onclick={refreshRecipes}
		>
			{#if refreshing}
				<Spinner class="size-4" />
			{:else}
				{m.common_retry()}
			{/if}
		</Button>

		<div class="min-w-0 md:col-span-3">
			<p
				id="recipe-category-filter"
				class="text-xs font-semibold uppercase tracking-wide text-muted-foreground"
			>
				{m.ocr_recipes_category_filter()}
			</p>
			<div
				class="mt-2 flex flex-wrap items-center gap-2"
				aria-labelledby="recipe-category-filter"
			>
				{#each categoryFilterOptions as option (option.key)}
					<button
						type="button"
						class={categoryFilterClass(option)}
						aria-pressed={activeCategoryKey === option.key}
						onclick={() => (selectedCategoryKey = option.key)}
					>
						<span class="min-w-0 truncate">{option.title}</span>
						<span class={categoryCountClass(option)}>{option.count}</span>
					</button>
				{/each}
			</div>
		</div>
	</section>

	{#if loadError}
		<Alert.Root variant="destructive">
			<Alert.Description>{loadError}</Alert.Description>
		</Alert.Root>
	{/if}

	{#if cloneError}
		<Alert.Root variant="destructive">
			<Alert.Description>{cloneError}</Alert.Description>
		</Alert.Root>
	{/if}

	<div class="flex items-center justify-between gap-3 text-sm text-muted-foreground">
		<p>{showingLabel(visibleRecipes.length)}</p>
	</div>

	{#if visibleRecipes.length === 0}
		<section class="flex min-h-[240px] flex-col items-center justify-center rounded-[8px] border border-border bg-muted/20 p-8 text-center">
			<SearchIcon class="size-6 text-muted-foreground" aria-hidden="true" />
			<h2 class="mt-4 text-lg font-semibold">{m.ocr_recipes_no_matches_title()}</h2>
			<p class="mt-2 max-w-md text-sm leading-6 text-muted-foreground">
				{m.ocr_recipes_no_matches_body()}
			</p>
		</section>
	{:else}
		<div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
			{#each visibleRecipes as item (item.recipe.id)}
				<Card.Root class="flex h-full flex-col rounded-[8px] border-border shadow-xs">
					<Card.Header class="gap-4">
						<div class="flex flex-wrap items-center justify-between gap-2">
							<p class="min-w-0 max-w-full truncate rounded-md border border-border bg-background px-2 py-1 text-xs font-medium text-muted-foreground">
								{categoryTitle(item.recipe)}
							</p>
							<p class="shrink-0 rounded-md bg-muted px-2 py-1 text-xs font-medium text-muted-foreground">
								{cloneCountLabel(item.recipe.counter)}
							</p>
						</div>

						<div class="min-w-0 space-y-2">
							<Card.Title class="line-clamp-2 text-lg leading-6">{item.recipe.title}</Card.Title>
							<Card.Description class="line-clamp-3 min-h-[60px]">
								{item.recipe.description || m.schemas_no_description()}
							</Card.Description>
						</div>

						<dl class="grid grid-cols-2 gap-2 border-y border-border py-3 text-xs">
							<div>
								<dt class="text-muted-foreground">{m.ocr_recipes_json_fields()}</dt>
								<dd class="mt-1 font-semibold text-foreground">{fieldCountLabel(item.summary.fieldCount)}</dd>
							</div>
							<div>
								<dt class="text-muted-foreground">{m.ocr_recipes_required()}</dt>
								<dd class="mt-1 font-semibold text-foreground">{requiredCountLabel(item.summary.requiredCount)}</dd>
							</div>
						</dl>
					</Card.Header>

					<Card.Content class="flex flex-1 flex-col gap-3">
						<div class="flex items-center justify-between gap-2">
							<h2 class="text-xs font-semibold uppercase tracking-wide text-muted-foreground">
								{m.ocr_recipes_json_fields()}
							</h2>
							<span class="text-xs text-muted-foreground">{m.ocr_recipes_strict_schema()}</span>
						</div>

						{#if item.summary.fields.length > 0}
							<div class="flex flex-wrap gap-1.5" role="list">
								{#each item.summary.fields as field (field.key)}
									<span
										class={fieldBadgeClass(field)}
										title={fieldTitle(field)}
										role="listitem"
										aria-label={fieldAriaLabel(field)}
									>
										<span class={fieldDotClass(field)} aria-hidden="true"></span>
										<span class="min-w-0 truncate">{field.key}</span>
									</span>
								{/each}
							</div>
						{:else}
							<p class="rounded-[8px] border border-border p-3 text-sm text-muted-foreground">
								{m.ocr_recipes_no_fields()}
							</p>
						{/if}
					</Card.Content>

					<Card.Footer class="mt-auto">
						{#if isLoggedIn}
							<Button
								type="button"
								class="w-full cursor-pointer"
								disabled={cloningId !== null}
								aria-label={m.ocr_recipes_clone_aria({ name: item.recipe.title })}
								onclick={() => cloneRecipe(item.recipe)}
							>
								{#if cloningId === item.recipe.id}
									<LoaderIcon class="size-4 animate-spin" aria-hidden="true" />
									{m.schemas_cloning()}
								{:else}
									<CopyPlusIcon class="size-4" aria-hidden="true" />
									{m.ocr_recipes_clone_recipe()}
								{/if}
							</Button>
						{:else}
							<Button href="/login" variant="outline" class="w-full">
								{m.ocr_recipes_log_in_to_clone()}
							</Button>
						{/if}
					</Card.Footer>
				</Card.Root>
			{/each}
		</div>
	{/if}
</div>
