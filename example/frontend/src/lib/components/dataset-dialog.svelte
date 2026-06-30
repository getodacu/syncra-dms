<script lang="ts">
	import ChevronRightIcon from "@lucide/svelte/icons/chevron-right";
	import ChevronsUpDownIcon from "@lucide/svelte/icons/chevrons-up-down";
	import LoaderIcon from "@lucide/svelte/icons/loader-circle";
	import SearchIcon from "@lucide/svelte/icons/search";

	import type { DatasetField } from "$lib/client/datasets";
	import { Badge } from "$lib/components/ui/badge/index.js";
	import { Button, buttonVariants } from "$lib/components/ui/button/index.js";
	import { Checkbox } from "$lib/components/ui/checkbox/index.js";
	import * as Command from "$lib/components/ui/command/index.js";
	import * as Dialog from "$lib/components/ui/dialog/index.js";
	import { Field, FieldError, FieldLabel } from "$lib/components/ui/field/index.js";
	import { Input } from "$lib/components/ui/input/index.js";
	import * as Popover from "$lib/components/ui/popover/index.js";
	import * as ScrollArea from "$lib/components/ui/scroll-area/index.js";
	import { m } from "$lib/paraglide/messages.js";
	import { cn } from "$lib/utils.js";
	import {
		buildDatasetFieldTree,
		type DatasetFieldTreeNode,
	} from "../../routes/app/datasets/field-tree";
	import {
		canSubmitDatasetDialog,
		datasetFieldNodePathMap,
		datasetFieldsAfterSchemaChange,
		validDatasetSelectedFields,
	} from "./nav-datasets-utils";

	type DatasetDialogMode = "create" | "edit";
	type DatasetDialogValue = { name: string; schema_id: string; selected_fields: DatasetField[] };
	type SchemaOption = {
		id: string;
		name: string;
		description: string;
		schema: unknown;
	};

	type Props = {
		open?: boolean;
		mode: DatasetDialogMode;
		initialValue?: DatasetDialogValue;
		schemas: SchemaOption[];
		schemasLoading?: boolean;
		schemasError?: Error | null;
		pending?: boolean;
		error?: Error | null;
		onSubmit: (value: DatasetDialogValue) => void;
	};

	let {
		open = $bindable(false),
		mode,
		initialValue,
		schemas = [],
		schemasLoading = false,
		schemasError = null,
		pending = false,
		error = null,
		onSubmit,
	}: Props = $props();

	let name = $state("");
	let selectedSchemaId = $state("");
	let selectedFields = $state<DatasetField[]>([]);
	let schemaPopoverOpen = $state(false);
	let expandedFieldIds = $state<string[]>([]);

	const initialValueKey = $derived(
		`${initialValue?.name ?? ""}:${initialValue?.schema_id ?? ""}:${JSON.stringify(
			initialValue?.selected_fields ?? []
		)}`
	);
	const title = $derived(
		mode === "edit" ? m.datasets_dialog_title_edit() : m.datasets_dialog_title_new()
	);
	const submitLabel = $derived(
		mode === "edit" ? m.datasets_save_changes() : m.datasets_create_dataset()
	);
	const selectedSchema = $derived(
		schemas.find((schema) => schema.id === selectedSchemaId) ?? null
	);
	const selectedSchemaLabel = $derived(
		selectedSchema?.name ??
			(selectedSchemaId
				? m.datasets_selected_schema()
				: schemasLoading
					? m.datasets_loading_schemas()
					: m.datasets_select_schema())
	);
	const fieldTree = $derived(buildDatasetFieldTree(selectedSchema?.schema));
	const validFieldNodesByPath = $derived(datasetFieldNodePathMap(fieldTree));
	const validSelectedFields = $derived(
		validDatasetSelectedFields(selectedFields, validFieldNodesByPath)
	);
	const selectedFieldPaths = $derived(new Set(validSelectedFields.map((field) => field.path)));
	const selectedFieldCountLabel = $derived.by(() => {
		if (validSelectedFields.length === 0) return m.datasets_no_fields_selected();
		if (validSelectedFields.length === 1) return m.datasets_one_field_selected();

		return m.datasets_fields_selected({ count: validSelectedFields.length });
	});
	const submitDisabled = $derived(
		!canSubmitDatasetDialog({
			pending,
			name,
			selectedSchemaExists: selectedSchema !== null,
			fieldTreeHasFields: fieldTree.length > 0,
			validSelectedFieldCount: validSelectedFields.length,
		})
	);

	$effect(() => {
		void initialValueKey;

		if (!open) {
			schemaPopoverOpen = false;
			return;
		}

		resetForm(initialValue);
	});

	function resetForm(value: DatasetDialogValue | undefined) {
		name = value?.name ?? "";
		selectedSchemaId = value?.schema_id ?? "";
		selectedFields = (value?.selected_fields ?? []).map((field) => ({ ...field }));
		schemaPopoverOpen = false;
		expandedFieldIds = [];
	}

	function selectSchema(id: string) {
		if (pending) return;

		const nextFields = datasetFieldsAfterSchemaChange(selectedSchemaId, id, selectedFields);
		const schemaChanged = id !== selectedSchemaId;
		selectedSchemaId = id;
		selectedFields = [...nextFields];
		if (schemaChanged) expandedFieldIds = [];
		schemaPopoverOpen = false;
	}

	function isFieldSelected(path: string) {
		return selectedFieldPaths.has(path);
	}

	function isFieldExpanded(id: string) {
		return expandedFieldIds.includes(id);
	}

	function toggleFieldExpanded(id: string) {
		expandedFieldIds = isFieldExpanded(id)
			? expandedFieldIds.filter((fieldId) => fieldId !== id)
			: [...expandedFieldIds, id];
	}

	function setFieldSelected(node: DatasetFieldTreeNode, checked: boolean | "indeterminate") {
		if (pending) return;

		const selected = checked === true;
		selectedFields = selected
			? [
					...selectedFields.filter((field) => field.path !== node.path),
					{ path: node.path, key: node.key, label: node.label },
				]
			: selectedFields.filter((field) => field.path !== node.path);
	}

	function handleSubmit(event: SubmitEvent) {
		event.preventDefault();

		const trimmedName = name.trim();
		if (submitDisabled || !trimmedName) return;

		onSubmit({
			name: trimmedName,
			schema_id: selectedSchemaId,
			selected_fields: validSelectedFields.map((field) => ({ ...field })),
		});
	}
</script>

{#snippet fieldRow(node: DatasetFieldTreeNode, depth: number)}
	{@const hasChildren = node.children.length > 0}
	{@const expanded = isFieldExpanded(node.id)}
	<div>
		<div
			class="flex min-h-8 min-w-0 items-center gap-2 rounded-md px-1.5 py-1 text-sm hover:bg-muted/60"
			style:padding-left={`${depth * 0.875 + 0.375}rem`}
		>
			{#if hasChildren}
				<button
					type="button"
					class="flex size-5 shrink-0 items-center justify-center rounded-sm text-muted-foreground hover:bg-muted hover:text-foreground"
					aria-label={expanded
						? m.datasets_collapse_field({ label: node.label })
						: m.datasets_expand_field({ label: node.label })}
					aria-expanded={expanded}
					onclick={() => toggleFieldExpanded(node.id)}
				>
					<ChevronRightIcon
						class={cn("size-4 transition-transform", expanded && "rotate-90")}
					/>
				</button>
			{:else}
				<span class="size-5 shrink-0" aria-hidden="true"></span>
			{/if}
			<Checkbox
				aria-label={m.datasets_select_field({ label: node.label })}
				disabled={pending}
				bind:checked={() => isFieldSelected(node.path), (checked) => setFieldSelected(node, checked)}
			/>
			<div class="flex min-w-0 flex-1 items-center gap-2">
				<span class="truncate">{node.label}</span>
				{#if node.jsonCell}
					<Badge variant="outline" class="h-4 rounded px-1 text-[10px] leading-none">
						{m.datasets_json_badge()}
					</Badge>
				{/if}
				<span class="ml-auto shrink-0 text-[11px] text-muted-foreground">{node.type}</span>
			</div>
		</div>
		{#if hasChildren && expanded}
			{#each node.children as child (child.id)}
				{@render fieldRow(child, depth + 1)}
			{/each}
		{/if}
	</div>
{/snippet}

<Dialog.Root bind:open>
	<Dialog.Content class="w-full gap-5 sm:max-w-2xl">
		<form class="grid gap-5" onsubmit={handleSubmit}>
			<Dialog.Header class="min-w-0">
				<Dialog.Title>{title}</Dialog.Title>
			</Dialog.Header>

			<div class="grid gap-4">
				<Field>
					<FieldLabel for="dataset-name">{m.datasets_name_column()}</FieldLabel>
					<Input
						id="dataset-name"
						bind:value={name}
						placeholder={m.datasets_name_placeholder()}
						disabled={pending}
						aria-invalid={Boolean(error)}
					/>
				</Field>

				<Field>
					<FieldLabel>{m.datasets_schema_column()}</FieldLabel>
					<Popover.Root bind:open={schemaPopoverOpen}>
						<Popover.Trigger
							type="button"
							role="combobox"
							aria-expanded={schemaPopoverOpen}
							disabled={pending}
							class={cn(
								buttonVariants({ variant: "outline" }),
								"h-9 w-full justify-between px-2.5 text-left"
							)}
						>
							<span class="flex min-w-0 items-center gap-2">
								<SearchIcon class="size-4 shrink-0 text-muted-foreground" />
								<span class="truncate">{selectedSchemaLabel}</span>
							</span>
							<ChevronsUpDownIcon class="size-4 shrink-0 text-muted-foreground" />
						</Popover.Trigger>
						<Popover.Content class="w-[min(calc(100vw-2rem),32rem)] p-0" align="start">
							<Command.Root>
								<Command.Input placeholder={m.datasets_search_schemas()} />
								<Command.List>
									{#if schemasLoading}
										<div class="px-3 py-6 text-center text-sm text-muted-foreground">
											{m.datasets_loading_schemas()}
										</div>
									{:else if schemasError}
										<div class="px-3 py-6 text-center text-sm text-destructive">
											{schemasError.message}
										</div>
									{:else}
										<Command.Empty>{m.datasets_no_schemas_found()}</Command.Empty>
										{#each schemas as schema (schema.id)}
											<Command.Item
												value={schema.id}
												keywords={[schema.name, schema.description, schema.id]}
												data-checked={selectedSchemaId === schema.id}
												onSelect={() => selectSchema(schema.id)}
											>
												<div class="flex min-w-0 flex-col">
													<span class="truncate">{schema.name}</span>
													{#if schema.description}
														<span class="truncate text-xs text-muted-foreground">
															{schema.description}
														</span>
													{/if}
												</div>
											</Command.Item>
										{/each}
									{/if}
								</Command.List>
							</Command.Root>
						</Popover.Content>
					</Popover.Root>
				</Field>

				<Field>
					<div class="flex items-center justify-between gap-3">
						<FieldLabel>{m.datasets_fields_column()}</FieldLabel>
						<span class="text-xs text-muted-foreground">{selectedFieldCountLabel}</span>
					</div>
					<div class="overflow-hidden rounded-md border bg-background">
						<ScrollArea.Root class="h-72">
							{#if !selectedSchemaId}
								<div class="px-3 py-8 text-center text-sm text-muted-foreground">
									{m.datasets_select_schema()}
								</div>
							{:else if fieldTree.length === 0}
								<div class="px-3 py-8 text-center text-sm text-muted-foreground">
									{m.datasets_no_fields()}
								</div>
							{:else}
								<div class="py-1">
									{#each fieldTree as node (node.id)}
										{@render fieldRow(node, 0)}
									{/each}
								</div>
							{/if}
						</ScrollArea.Root>
					</div>
				</Field>

				{#if error}
					<FieldError>{error.message}</FieldError>
				{/if}
			</div>

			<Dialog.Footer>
				<Button type="button" variant="outline" disabled={pending} onclick={() => (open = false)}>
					{m.datasets_cancel()}
				</Button>
				<Button type="submit" disabled={submitDisabled}>
					{#if pending}
						<LoaderIcon class="size-4 animate-spin" />
					{/if}
					{submitLabel}
				</Button>
			</Dialog.Footer>
		</form>
	</Dialog.Content>
</Dialog.Root>
