<script lang="ts">
	import ChevronsUpDownIcon from "@lucide/svelte/icons/chevrons-up-down";
	import LoaderIcon from "@lucide/svelte/icons/loader-circle";
	import SearchIcon from "@lucide/svelte/icons/search";

	import { Button, buttonVariants } from "$lib/components/ui/button/index.js";
	import * as Command from "$lib/components/ui/command/index.js";
	import * as Dialog from "$lib/components/ui/dialog/index.js";
	import {
		Field,
		FieldDescription,
		FieldError,
		FieldLabel,
	} from "$lib/components/ui/field/index.js";
	import { Input } from "$lib/components/ui/input/index.js";
	import * as Popover from "$lib/components/ui/popover/index.js";
	import { m } from "$lib/paraglide/messages.js";
	import { cn } from "$lib/utils.js";

	type CollectionDialogMode = "create" | "edit";
	type CollectionDialogValue = { name: string; schema_ids: string[] };
	type SchemaOption = { id: string; name: string; description: string };

	type Props = {
		open?: boolean;
		mode: CollectionDialogMode;
		initialValue?: CollectionDialogValue;
		schemas: SchemaOption[];
		schemasLoading?: boolean;
		schemasError?: Error | null;
		pending?: boolean;
		error?: Error | null;
		onSubmit: (value: CollectionDialogValue) => void;
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
	let selectedSchemaIds = $state<string[]>([]);
	let schemaPopoverOpen = $state(false);

	const initialValueKey = $derived(
		`${initialValue?.name ?? ""}:${(initialValue?.schema_ids ?? []).join("\u001f")}`
	);
	const title = $derived(
		mode === "edit"
			? m.documents_collection_dialog_title_edit()
			: m.documents_collection_dialog_title_new()
	);
	const description = $derived(
		mode === "edit"
			? m.documents_collection_dialog_description_edit()
			: m.documents_collection_dialog_description_new()
	);
	const submitLabel = $derived(
		mode === "edit" ? m.documents_save_changes() : m.documents_create_collection()
	);
	const selectedSchemaCountLabel = $derived.by(() => {
		if (selectedSchemaIds.length === 0) return m.documents_no_schemas_selected();
		if (selectedSchemaIds.length === 1) return m.documents_one_schema_selected();

		return m.documents_schemas_selected({ count: selectedSchemaIds.length });
	});
	const submitDisabled = $derived(name.trim().length === 0 || pending);

	$effect(() => {
		void initialValueKey;

		if (!open) {
			schemaPopoverOpen = false;
			return;
		}

		resetForm(initialValue);
	});

	function resetForm(value: CollectionDialogValue | undefined) {
		name = value?.name ?? "";
		selectedSchemaIds = [...(value?.schema_ids ?? [])];
		schemaPopoverOpen = false;
	}

	function isSchemaSelected(id: string) {
		return selectedSchemaIds.includes(id);
	}

	function toggleSchema(id: string) {
		if (pending) return;

		selectedSchemaIds = isSchemaSelected(id)
			? selectedSchemaIds.filter((schemaId) => schemaId !== id)
			: [...selectedSchemaIds, id];
	}

	function handleSubmit(event: SubmitEvent) {
		event.preventDefault();

		const trimmedName = name.trim();
		if (!trimmedName || pending) return;

		onSubmit({
			name: trimmedName,
			schema_ids: selectedSchemaIds,
		});
	}
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="w-full gap-5 sm:max-w-lg">
		<form class="grid gap-5" onsubmit={handleSubmit}>
			<Dialog.Header class="min-w-0">
				<Dialog.Title>{title}</Dialog.Title>
				<Dialog.Description class="text-sm">{description}</Dialog.Description>
			</Dialog.Header>

			<div class="grid gap-4">
				<Field>
					<FieldLabel for="collection-name">{m.documents_name_column()}</FieldLabel>
					<Input
						id="collection-name"
						bind:value={name}
						placeholder={m.documents_collection_name_placeholder()}
						disabled={pending}
						aria-invalid={Boolean(error)}
					/>
				</Field>

				<Field>
					<FieldLabel>{m.documents_schemas_label()}</FieldLabel>
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
								<span class="truncate">{selectedSchemaCountLabel}</span>
							</span>
							<ChevronsUpDownIcon class="size-4 shrink-0 text-muted-foreground" />
						</Popover.Trigger>
						<Popover.Content class="w-[min(calc(100vw-2rem),28rem)] p-0" align="start">
							<Command.Root>
								<Command.Input placeholder={m.documents_search_schemas()} />
								<Command.List>
									{#if schemasLoading}
										<div class="px-3 py-6 text-center text-sm text-muted-foreground">
											{m.documents_loading_schemas()}
										</div>
									{:else if schemasError}
										<div class="px-3 py-6 text-center text-sm text-destructive">
											{schemasError.message}
										</div>
									{:else}
										<Command.Empty>{m.documents_no_schemas_found()}</Command.Empty>
										{#each schemas as schema (schema.id)}
											<Command.Item
												value={schema.id}
												keywords={[schema.name, schema.description, schema.id]}
												data-checked={isSchemaSelected(schema.id)}
												onSelect={() => toggleSchema(schema.id)}
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
					<FieldDescription class="text-xs">
						{m.documents_collection_schema_hint()}
					</FieldDescription>
				</Field>

				{#if error}
					<FieldError>{error.message}</FieldError>
				{/if}
			</div>

			<Dialog.Footer>
				<Button type="button" variant="outline" disabled={pending} onclick={() => (open = false)}>
					{m.documents_cancel()}
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
