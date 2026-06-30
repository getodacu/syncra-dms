<script lang="ts">
	import ChevronsUpDownIcon from "@lucide/svelte/icons/chevrons-up-down";
	import FolderIcon from "@lucide/svelte/icons/folder";
	import LoaderIcon from "@lucide/svelte/icons/loader-circle";
	import SearchIcon from "@lucide/svelte/icons/search";

	import { Button, buttonVariants } from "$lib/components/ui/button/index.js";
	import * as Command from "$lib/components/ui/command/index.js";
	import * as Dialog from "$lib/components/ui/dialog/index.js";
	import { Field, FieldError, FieldLabel } from "$lib/components/ui/field/index.js";
	import * as Popover from "$lib/components/ui/popover/index.js";
	import type { CollectionResponse } from "$lib/client/collections";
	import { m } from "$lib/paraglide/messages.js";
	import { cn } from "$lib/utils.js";

	type Props = {
		open?: boolean;
		collections: CollectionResponse[];
		collectionsLoading?: boolean;
		collectionsError?: Error | null;
		selectedCount: number;
		pending?: boolean;
		error?: Error | null;
		onSubmit: (collectionIds: string[]) => void;
	};

	let {
		open = $bindable(false),
		collections = [],
		collectionsLoading = false,
		collectionsError = null,
		selectedCount,
		pending = false,
		error = null,
		onSubmit,
	}: Props = $props();

	let selectedCollectionIds = $state<string[]>([]);
	let collectionPopoverOpen = $state(false);

	const selectedCollectionCountLabel = $derived.by(() => {
		if (selectedCollectionIds.length === 0) return m.documents_no_collections_selected();
		if (selectedCollectionIds.length === 1) return m.documents_one_collection_selected();

		return m.documents_collections_selected({ count: selectedCollectionIds.length });
	});
	const submitLabel = $derived(
		selectedCollectionIds.length === 0 ? m.documents_remove_from_all() : m.documents_move()
	);
	const dialogDescription = $derived(
		selectedCount === 1
			? m.documents_move_description_one()
			: m.documents_move_description_other({ count: selectedCount })
	);

	$effect(() => {
		if (open) {
			selectedCollectionIds = [];
			collectionPopoverOpen = false;
		} else {
			collectionPopoverOpen = false;
		}
	});

	function isCollectionSelected(id: string) {
		return selectedCollectionIds.includes(id);
	}

	function toggleCollection(id: string) {
		if (pending) return;

		selectedCollectionIds = isCollectionSelected(id)
			? selectedCollectionIds.filter((collectionId) => collectionId !== id)
			: [...selectedCollectionIds, id];
	}

	function handleSubmit(event: SubmitEvent) {
		event.preventDefault();
		if (pending) return;

		onSubmit(selectedCollectionIds);
	}
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="w-full gap-5 sm:max-w-lg">
		<form class="grid gap-5" onsubmit={handleSubmit}>
			<Dialog.Header class="min-w-0">
				<Dialog.Title>{m.documents_move_documents()}</Dialog.Title>
				<Dialog.Description class="text-sm">{dialogDescription}</Dialog.Description>
			</Dialog.Header>

			<div class="grid gap-4">
				<Field>
					<FieldLabel>{m.documents_collections_label()}</FieldLabel>
					<Popover.Root bind:open={collectionPopoverOpen}>
						<Popover.Trigger
							type="button"
							role="combobox"
							aria-expanded={collectionPopoverOpen}
							disabled={pending}
							class={cn(
								buttonVariants({ variant: "outline" }),
								"h-9 w-full justify-between px-2.5 text-left"
							)}
						>
							<span class="flex min-w-0 items-center gap-2">
								<SearchIcon class="size-4 shrink-0 text-muted-foreground" />
								<span class="truncate">{selectedCollectionCountLabel}</span>
							</span>
							<ChevronsUpDownIcon class="size-4 shrink-0 text-muted-foreground" />
						</Popover.Trigger>
						<Popover.Content class="w-[min(calc(100vw-2rem),28rem)] p-0" align="start">
							<Command.Root>
								<Command.Input placeholder={m.documents_search_collections()} />
								<Command.List>
									{#if collectionsLoading}
										<div class="px-3 py-6 text-center text-sm text-muted-foreground">
											{m.documents_loading_collections()}
										</div>
									{:else if collectionsError}
										<div class="px-3 py-6 text-center text-sm text-destructive">
											{collectionsError.message}
										</div>
									{:else}
										<Command.Empty>{m.documents_no_collections_found()}</Command.Empty>
										{#each collections as collection (collection.id)}
											<Command.Item
												value={collection.id}
												keywords={[collection.name, collection.id]}
												data-checked={isCollectionSelected(collection.id)}
												onSelect={() => toggleCollection(collection.id)}
											>
												<FolderIcon class="size-4 text-muted-foreground" />
												<span class="min-w-0 truncate">{collection.name}</span>
											</Command.Item>
										{/each}
									{/if}
								</Command.List>
							</Command.Root>
						</Popover.Content>
					</Popover.Root>
				</Field>

				{#if error}
					<FieldError>{error.message}</FieldError>
				{/if}
			</div>

			<Dialog.Footer>
				<Button type="button" variant="outline" disabled={pending} onclick={() => (open = false)}>
					{m.documents_cancel()}
				</Button>
				<Button type="submit" disabled={pending}>
					{#if pending}
						<LoaderIcon class="size-4 animate-spin" />
					{/if}
					{submitLabel}
				</Button>
			</Dialog.Footer>
		</form>
	</Dialog.Content>
</Dialog.Root>
