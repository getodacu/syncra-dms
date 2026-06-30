<script lang="ts">
	import JsonSchemaBuilder from "$lib/components/json-schema-builder.svelte";
	import { Button } from "$lib/components/ui/button/index.js";
	import {
		Field,
		FieldContent,
		FieldDescription,
		FieldError,
		FieldGroup,
		FieldLabel,
	} from "$lib/components/ui/field/index.js";
	import { Input } from "$lib/components/ui/input/index.js";
	import { Switch } from "$lib/components/ui/switch/index.js";
	import { Textarea } from "$lib/components/ui/textarea/index.js";
	import { m } from "$lib/paraglide/messages.js";
	import CheckIcon from "@lucide/svelte/icons/check";
	import FileCodeIcon from "@lucide/svelte/icons/file-code";
	import SettingsIcon from "@lucide/svelte/icons/settings";
	import AlertTriangleIcon from "@lucide/svelte/icons/alert-triangle";
	import SchemaIdCopy from "./schema-id-copy.svelte";

	export type JsonSchemaValue = boolean | Record<string, unknown>;
	export type JsonSchemaObject = Record<string, unknown>;
	export type SchemaEditorSubmitInput = {
		name: string;
		description: string;
		strict: boolean;
		schema: JsonSchemaObject;
	};

	let {
		title,
		schemaId = null,
		descriptionText,
		submitLabel,
		pending = false,
		initial,
		serverError = null,
		successMessage = null,
		clonePending = false,
		onSubmit,
		onClone = null,
		onDirty = () => undefined,
	}: {
		title: string;
		schemaId?: string | null;
		descriptionText: string;
		submitLabel: string;
		pending?: boolean;
		initial?: Partial<SchemaEditorSubmitInput>;
		serverError?: string | null;
		successMessage?: string | null;
		clonePending?: boolean;
		onSubmit: (input: SchemaEditorSubmitInput) => void;
		onClone?: (() => void) | null;
		onDirty?: () => void;
	} = $props();

	// svelte-ignore state_referenced_locally
	let name = $state(initial?.name ?? "");
	// svelte-ignore state_referenced_locally
	let description = $state(initial?.description ?? "");
	// svelte-ignore state_referenced_locally
	let strict = $state(initial?.strict ?? true);
	// JSONJoy uses structuredClone internally, so keep the schema as a plain cloneable value.
	// svelte-ignore state_referenced_locally
	let schema = $state.raw<JsonSchemaValue>(initial?.schema ?? { type: "object", properties: {} });
	let validationError = $state<string | null>(null);

	const errorMessage = $derived(validationError ?? serverError);
	const nameError = $derived(
		validationError === m.schemas_validation_name_required() ||
			validationError === m.schemas_validation_name_too_long()
			? validationError
			: null
	);

	function isSchemaObject(value: JsonSchemaValue): value is JsonSchemaObject {
		return typeof value === "object" && value !== null && !Array.isArray(value);
	}

	function clearValidationError() {
		validationError = null;
	}

	function markDirty() {
		clearValidationError();
		onDirty();
	}

	function updateSchema(next: JsonSchemaValue) {
		schema = next;
		markDirty();
	}

	function updateName(next: string) {
		name = next;
		markDirty();
	}

	function updateDescription(next: string) {
		description = next;
		markDirty();
	}

	function updateStrict(next: boolean) {
		strict = next;
		markDirty();
	}

	function submit(event: SubmitEvent) {
		event.preventDefault();
		validationError = null;

		const trimmedName = name.trim();

		if (!trimmedName) {
			validationError = m.schemas_validation_name_required();
			return;
		}

		if (Array.from(trimmedName).length > 160) {
			validationError = m.schemas_validation_name_too_long();
			return;
		}

		if (!isSchemaObject(schema)) {
			validationError = m.schemas_validation_schema_object();
			return;
		}

		onSubmit({
			name: trimmedName,
			description,
			strict,
			schema,
		});
	}
</script>

<div class="@container/main flex min-h-0 flex-1 flex-col bg-muted/5">
	<form class="flex min-h-0 flex-1 flex-col gap-6 p-4 lg:p-6" onsubmit={submit}>
		<!-- Pinned Header Bar -->
		<div class="flex flex-col gap-4 border-b bg-background rounded-2xl border p-5 shadow-xs shrink-0 transition-all duration-300">
			<div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
				<div class="space-y-1.5 min-w-0">
					<div class="flex flex-wrap items-center gap-2">
						<span class="inline-flex h-5 items-center rounded-md bg-primary/10 px-2 text-[10px] font-bold text-primary ring-1 ring-primary/20">
							{m.schemas_editor_badge()}
						</span>
						{#if strict}
							<span class="inline-flex h-5 items-center gap-1 rounded-md bg-amber-500/10 px-2 text-[10px] font-bold text-amber-600 dark:bg-amber-500/20 dark:text-amber-400 ring-1 ring-amber-500/20">
								{m.schemas_strict_mode()}
							</span>
						{:else}
							<span class="inline-flex h-5 items-center gap-1 rounded-md bg-muted px-2 text-[10px] font-bold text-muted-foreground ring-1 ring-muted-foreground/10">
								{m.schemas_flexible_mode()}
							</span>
						{/if}
					</div>
					<h1 class="text-2xl font-extrabold tracking-tight text-foreground truncate">
						{name.trim() || title}
					</h1>
					{#if schemaId}
						<SchemaIdCopy schemaId={schemaId} showLabel />
					{/if}
					<p class="text-xs text-muted-foreground max-w-2xl truncate">
						{description.trim() || descriptionText}
					</p>
				</div>
				<div class="flex items-center gap-2 shrink-0">
					{#if onClone}
						<Button
							type="button"
							variant="outline"
							disabled={clonePending || pending}
							class="h-10 px-5 shadow-sm cursor-pointer"
							onclick={onClone}
						>
							{#if clonePending}
								<span class="mr-2 animate-spin size-4 border-2 border-foreground border-t-transparent rounded-full"></span>
								{m.schemas_cloning()}
							{:else}
								{m.schemas_clone()}
							{/if}
						</Button>
					{/if}
					<Button type="submit" disabled={pending} class="h-10 px-5 shadow-sm cursor-pointer">
						{#if pending}
							<span class="mr-2 animate-spin size-4 border-2 border-primary-foreground border-t-transparent rounded-full"></span>
							{m.schemas_saving()}
						{:else}
							{submitLabel}
						{/if}
					</Button>
				</div>
			</div>

			<!-- Feedback alerts -->
			{#if successMessage}
				<div
					role="status"
					aria-live="polite"
					class="mt-2 flex items-center gap-2 rounded-lg border border-emerald-500/20 bg-emerald-500/10 px-3.5 py-2.5 text-xs font-semibold text-emerald-600 dark:text-emerald-400"
				>
					<CheckIcon class="size-4 shrink-0 text-emerald-500" />
					<span>{successMessage}</span>
				</div>
			{/if}

			{#if errorMessage && !nameError}
				<div class="space-y-2 mt-2">
					
					{#if errorMessage && !nameError}
						<div
							class="flex items-center gap-2 rounded-lg border border-red-500/20 bg-red-500/10 px-3.5 py-2.5 text-xs font-semibold text-red-600 dark:text-red-400"
						>
							<AlertTriangleIcon class="size-4 shrink-0 text-red-500" />
							<span>{errorMessage}</span>
						</div>
					{/if}
				</div>
			{/if}
		</div>

		<!-- Top Section: Schema Settings Card -->
		<div class="rounded-2xl border border-border/80 bg-card p-5 shadow-3xs shrink-0">
			<div class="flex items-center gap-2 border-b pb-3.5 mb-5 shrink-0">
				<SettingsIcon class="size-4 text-indigo-500" />
				<h2 class="text-xs font-bold uppercase tracking-wider text-muted-foreground">{m.schemas_general_settings()}</h2>
			</div>

			<div class="grid grid-cols-1 md:grid-cols-[1fr_320px] gap-5">
				<div class="space-y-4">
					<!-- Name Field -->
					<Field>
						<div class="flex items-center justify-between">
							<FieldLabel for="schema-name" class="font-bold text-foreground text-xs uppercase tracking-wider text-muted-foreground/90">{m.schemas_schema_name_label()}</FieldLabel>
							<span class="text-[10px] text-red-500 font-semibold">{m.common_required()}</span>
						</div>
						<Input
							id="schema-name"
							value={name}
							oninput={(event) => updateName(event.currentTarget.value)}
							placeholder={m.schemas_schema_name_placeholder()}
							aria-invalid={Boolean(nameError) || undefined}
							aria-describedby={nameError ? "schema-name-error" : undefined}
							class="h-10 focus-visible:ring-primary/30 rounded-lg text-sm"
						/>
						{#if nameError}
							<FieldError id="schema-name-error">{nameError}</FieldError>
						{/if}
					</Field>

					<!-- Description Field -->
					<Field>
						<FieldLabel for="schema-description" class="font-bold text-foreground text-xs uppercase tracking-wider text-muted-foreground/90">{m.schemas_description_label()}</FieldLabel>
						<Textarea
							id="schema-description"
							value={description}
							oninput={(event) => updateDescription(event.currentTarget.value)}
							placeholder={m.schemas_description_placeholder()}
							rows={2}
							class="resize-none focus-visible:ring-primary/30 rounded-lg text-sm"
						/>
					</Field>
				</div>

				<div class="flex flex-col justify-start">
					<!-- Strict Mode Toggle -->
					<Field class="rounded-xl border border-border/80 bg-muted/10 p-4 transition-all hover:bg-muted/20 h-full flex flex-col justify-center">
						<div class="flex items-start justify-between gap-4">
							<div class="space-y-1">
								<FieldLabel for="schema-strict" class="font-bold text-foreground text-xs uppercase tracking-wider text-muted-foreground/90">{m.schemas_strict_mode()}</FieldLabel>
								<FieldDescription class="text-[11px] text-muted-foreground leading-normal">
									{m.schemas_strict_mode_description()}
								</FieldDescription>
							</div>
							<Switch id="schema-strict" checked={strict} onCheckedChange={updateStrict} class="cursor-pointer" />
						</div>
					</Field>
				</div>
			</div>
		</div>

		<!-- Bottom Section: Structure Builder Card -->
		<div class="flex min-h-[640px] flex-1 flex-col rounded-2xl border border-border/80 bg-background shadow-3xs overflow-hidden">
			<div class="flex items-center justify-between border-b px-5 py-3.5 bg-muted/20 shrink-0">
				<div class="flex items-center gap-2">
					<FileCodeIcon class="size-4 text-indigo-500" />
					<span class="text-xs font-bold uppercase tracking-wider text-muted-foreground">{m.schemas_structure_designer()}</span>
				</div>
				<span class="text-[10px] text-muted-foreground/80 font-medium">{m.schemas_visual_node_designer()}</span>
			</div>
			<div class="min-h-0 flex-1">
				<JsonSchemaBuilder value={schema} onChange={updateSchema} class="h-full min-h-0 border-none rounded-none" />
			</div>
		</div>
	</form>
</div>
