<script lang="ts">
	import JsonSchemaBuilder from "$lib/components/json-schema-builder.svelte";
	import { Button } from "$lib/components/ui/button/index.js";
	import {
		Field,
		FieldError,
		FieldLabel,
	} from "$lib/components/ui/field/index.js";
	import { Input } from "$lib/components/ui/input/index.js";
	import { Textarea } from "$lib/components/ui/textarea/index.js";
	import { m } from "$lib/paraglide/messages.js";
	import AlertTriangleIcon from "@lucide/svelte/icons/alert-triangle";
	import CheckIcon from "@lucide/svelte/icons/check";
	import FileCodeIcon from "@lucide/svelte/icons/file-code";
	import SettingsIcon from "@lucide/svelte/icons/settings";
	import TrashIcon from "@lucide/svelte/icons/trash";
	import type { AdminJSONRecipeCategoryResponse } from "./api";

	export type JsonRecipeValue = boolean | Record<string, unknown>;
	export type JsonRecipeObject = Record<string, unknown>;
	export type JsonRecipeEditorSubmitInput = {
		title: string;
		description: string;
		json: JsonRecipeObject;
		category_id: string | null;
	};

	let {
		title,
		descriptionText,
		submitLabel,
		pending = false,
		deletePending = false,
		initial,
		categories = [],
		serverError = null,
		successMessage = null,
		onSubmit,
		onDelete = null,
		onDirty = () => undefined,
	}: {
		title: string;
		descriptionText: string;
		submitLabel: string;
		pending?: boolean;
		deletePending?: boolean;
		initial?: Partial<JsonRecipeEditorSubmitInput>;
		categories?: AdminJSONRecipeCategoryResponse[];
		serverError?: string | null;
		successMessage?: string | null;
		onSubmit: (input: JsonRecipeEditorSubmitInput) => void;
		onDelete?: (() => void) | null;
		onDirty?: () => void;
	} = $props();

	// svelte-ignore state_referenced_locally
	let recipeTitle = $state(initial?.title ?? "");
	// svelte-ignore state_referenced_locally
	let description = $state(initial?.description ?? "");
	// svelte-ignore state_referenced_locally
	let selectedCategoryId = $state(initial?.category_id ?? "");
	// JSONJoy expects cloneable plain values.
	// svelte-ignore state_referenced_locally
	let recipeJson = $state.raw<JsonRecipeValue>(initial?.json ?? { type: "object", properties: {} });
	let validationError = $state<string | null>(null);

	const errorMessage = $derived(validationError ?? serverError);
	const titleError = $derived(
		validationError === m.json_recipes_validation_title_required() ||
			validationError === m.json_recipes_validation_title_too_long()
			? validationError
			: null
	);

	function isJsonObject(value: JsonRecipeValue): value is JsonRecipeObject {
		return typeof value === "object" && value !== null && !Array.isArray(value);
	}

	function markDirty() {
		validationError = null;
		onDirty();
	}

	function updateTitle(next: string) {
		recipeTitle = next;
		markDirty();
	}

	function updateDescription(next: string) {
		description = next;
		markDirty();
	}

	function updateCategory(next: string) {
		selectedCategoryId = next;
		markDirty();
	}

	function updateJson(next: JsonRecipeValue) {
		recipeJson = next;
		markDirty();
	}

	function submit(event: SubmitEvent) {
		event.preventDefault();
		validationError = null;

		const trimmedTitle = recipeTitle.trim();
		if (!trimmedTitle) {
			validationError = m.json_recipes_validation_title_required();
			return;
		}
		if (Array.from(trimmedTitle).length > 160) {
			validationError = m.json_recipes_validation_title_too_long();
			return;
		}
		if (!isJsonObject(recipeJson)) {
			validationError = m.json_recipes_validation_json_object();
			return;
		}

		onSubmit({
			title: trimmedTitle,
			description,
			json: recipeJson,
			category_id: selectedCategoryId || null,
		});
	}
</script>

<div class="@container/main flex min-h-0 flex-1 flex-col bg-muted/5">
	<form class="flex min-h-0 flex-1 flex-col gap-6 p-4 lg:p-6" onsubmit={submit}>
		<div class="flex shrink-0 flex-col gap-4 rounded-xl border bg-background p-5 shadow-xs">
			<div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
				<div class="min-w-0 space-y-1.5">
					<div class="flex flex-wrap items-center gap-2">
						<span class="inline-flex h-5 items-center rounded-md bg-primary/10 px-2 text-[10px] font-bold text-primary ring-1 ring-primary/20">
							{m.json_recipes_editor_badge()}
						</span>
					</div>
					<h1 class="truncate text-2xl font-extrabold text-foreground">
						{recipeTitle.trim() || title}
					</h1>
					<p class="max-w-2xl truncate text-xs text-muted-foreground">
						{description.trim() || descriptionText}
					</p>
				</div>
				<div class="flex shrink-0 items-center gap-2">
					{#if onDelete}
						<Button
							type="button"
							variant="outline"
							disabled={deletePending || pending}
							class="h-10 px-4 text-destructive"
							onclick={onDelete}
						>
							<TrashIcon class="mr-2 size-4" />
							{#if deletePending}
								{m.json_recipes_deleting()}
							{:else}
								{m.common_delete()}
							{/if}
						</Button>
					{/if}
					<Button type="submit" disabled={pending || deletePending} class="h-10 px-5">
						{#if pending}
							<span class="mr-2 size-4 animate-spin rounded-full border-2 border-primary-foreground border-t-transparent"></span>
							{m.json_recipes_saving()}
						{:else}
							{submitLabel}
						{/if}
					</Button>
				</div>
			</div>

			{#if successMessage}
				<div
					role="status"
					aria-live="polite"
					class="flex items-center gap-2 rounded-lg border border-emerald-500/20 bg-emerald-500/10 px-3.5 py-2.5 text-xs font-semibold text-emerald-600 dark:text-emerald-400"
				>
					<CheckIcon class="size-4 shrink-0 text-emerald-500" />
					<span>{successMessage}</span>
				</div>
			{/if}

			{#if errorMessage && !titleError}
				<div
					class="flex items-center gap-2 rounded-lg border border-red-500/20 bg-red-500/10 px-3.5 py-2.5 text-xs font-semibold text-red-600 dark:text-red-400"
				>
					<AlertTriangleIcon class="size-4 shrink-0 text-red-500" />
					<span>{errorMessage}</span>
				</div>
			{/if}
		</div>

		<div class="shrink-0 rounded-xl border border-border/80 bg-card p-5 shadow-3xs">
			<div class="mb-5 flex items-center gap-2 border-b pb-3.5">
				<SettingsIcon class="size-4 text-indigo-500" />
				<h2 class="text-xs font-bold uppercase tracking-wider text-muted-foreground">{m.json_recipes_general_settings()}</h2>
			</div>

			<div class="grid grid-cols-1 gap-5 md:grid-cols-[minmax(0,1fr)_minmax(220px,320px)_minmax(260px,420px)]">
				<Field>
					<div class="flex items-center justify-between">
						<FieldLabel for="json-recipe-title" class="text-xs font-bold uppercase tracking-wider text-muted-foreground/90">
							{m.json_recipes_title_label()}
						</FieldLabel>
						<span class="text-[10px] font-semibold text-red-500">{m.common_required()}</span>
					</div>
					<Input
						id="json-recipe-title"
						value={recipeTitle}
						oninput={(event) => updateTitle(event.currentTarget.value)}
						placeholder={m.json_recipes_title_placeholder()}
						aria-invalid={Boolean(titleError) || undefined}
						aria-describedby={titleError ? "json-recipe-title-error" : undefined}
						class="h-10 rounded-lg text-sm"
					/>
					{#if titleError}
						<FieldError id="json-recipe-title-error">{titleError}</FieldError>
					{/if}
				</Field>

				<Field>
					<FieldLabel for="json-recipe-category" class="text-xs font-bold uppercase tracking-wider text-muted-foreground/90">
						{m.json_recipes_category_label()}
					</FieldLabel>
					<select
						id="json-recipe-category"
						value={selectedCategoryId}
						onchange={(event) => updateCategory(event.currentTarget.value)}
						class="h-10 rounded-lg border border-input bg-background px-3 text-sm shadow-xs outline-none transition-colors focus-visible:ring-[3px] focus-visible:ring-ring/50"
					>
						<option value="">{m.json_recipes_others()}</option>
						{#each categories as category (category.id)}
							<option value={category.id}>{category.title.en} / {category.title.ro}</option>
						{/each}
					</select>
				</Field>

				<Field>
					<FieldLabel for="json-recipe-description" class="text-xs font-bold uppercase tracking-wider text-muted-foreground/90">
						{m.json_recipes_description_label()}
					</FieldLabel>
					<Textarea
						id="json-recipe-description"
						value={description}
						oninput={(event) => updateDescription(event.currentTarget.value)}
						placeholder={m.json_recipes_description_placeholder()}
						rows={2}
						class="resize-none rounded-lg text-sm"
					/>
				</Field>
			</div>
		</div>

		<div class="flex min-h-[640px] flex-1 flex-col overflow-hidden rounded-xl border border-border/80 bg-background shadow-3xs">
			<div class="flex shrink-0 items-center justify-between border-b bg-muted/20 px-5 py-3.5">
				<div class="flex items-center gap-2">
					<FileCodeIcon class="size-4 text-indigo-500" />
					<span class="text-xs font-bold uppercase tracking-wider text-muted-foreground">{m.json_recipes_structure_designer()}</span>
				</div>
				<span class="text-[10px] font-medium text-muted-foreground/80">{m.json_recipes_visual_node_designer()}</span>
			</div>
			<div class="min-h-0 flex-1">
				<JsonSchemaBuilder value={recipeJson} onChange={updateJson} class="h-full min-h-0 rounded-none border-none" />
			</div>
		</div>
	</form>
</div>
