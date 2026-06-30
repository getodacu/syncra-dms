<script lang="ts">
	import AlertCircleIcon from "@lucide/svelte/icons/alert-circle";
	import EllipsisIcon from "@lucide/svelte/icons/ellipsis";
	import FilesIcon from "@lucide/svelte/icons/files";
	import FolderIcon from "@lucide/svelte/icons/folder";
	import PencilIcon from "@lucide/svelte/icons/pencil";
	import PlusIcon from "@lucide/svelte/icons/plus";
	import Trash2Icon from "@lucide/svelte/icons/trash-2";
	import { goto } from "$app/navigation";
	import { page } from "$app/state";
	import { createMutation, createQuery, useQueryClient } from "@tanstack/svelte-query";

	import {
		createCollection,
		deleteCollection,
		fetchCollections,
		updateCollection,
		type CollectionInput,
		type CollectionListResponse,
		type CollectionResponse,
	} from "$lib/client/collections";
	import {
		fetchPersonalSchemaOptions,
		PERSONAL_SCHEMA_OPTIONS_QUERY_KEY,
		type PersonalSchemaOption,
	} from "$lib/client/schemas";
	import { confirmDelete } from "$lib/components/ui/confirm-delete-dialog/index.js";
	import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index.js";
	import * as Sidebar from "$lib/components/ui/sidebar/index.js";
	import { m } from "$lib/paraglide/messages.js";
	import CollectionDialog from "./collection-dialog.svelte";
	import {
		COLLECTION_LOADING_ROWS,
		COLLECTIONS_QUERY_KEY,
		collectionDeleteSuccessEffects,
		collectionDialogError,
		collectionDialogInitialValue,
		collectionDialogPending,
		collectionHref,
		collectionListOverflows,
		collectionListStatus,
		collectionSubmitAction,
		collectionUpdateSuccessInvalidationKeys,
		composeCollectionMenuTriggerClick,
		closeCollectionDialogState,
		isAllDocumentsActive,
		isCollectionActive as isCollectionActiveByRoute,
		openCreateCollectionDialogState,
		openEditCollectionDialogState,
		retryCollections,
		runCollectionMenuAction,
		type CollectionDialogMode,
		type CollectionDialogState,
		type CollectionQueryKey,
	} from "./nav-collections-utils";

	type UpdateCollectionVariables = { id: string; input: CollectionInput };
	type DeleteCollectionVariables = { collection: CollectionResponse };

	const queryClient = useQueryClient();
	const sidebar = Sidebar.useSidebar();

	let dialogOpen = $state(false);
	let dialogMode = $state<CollectionDialogMode>("create");
	let editingCollection = $state<CollectionResponse | null>(null);

	const selectedCollectionId = $derived(page.url.searchParams.get("collection"));
	const currentPathname = $derived(page.url.pathname);
	const collectionsQuery = createQuery<CollectionListResponse, Error>(() => ({
		queryKey: ["collections"],
		queryFn: () => fetchCollections(fetch),
	}));
	const schemasQuery = createQuery<PersonalSchemaOption[], Error>(() => ({
		queryKey: PERSONAL_SCHEMA_OPTIONS_QUERY_KEY,
		queryFn: () => fetchPersonalSchemaOptions(fetch),
	}));
	const createCollectionMutation = createMutation<CollectionResponse, Error, CollectionInput>(() => ({
		mutationKey: ["collections", "create"],
		mutationFn: (input) => createCollection(fetch, input),
		onSuccess: () => {
			invalidateQueryKeys([COLLECTIONS_QUERY_KEY]);
			closeDialog();
		},
	}));
	const updateCollectionMutation = createMutation<
		CollectionResponse,
		Error,
		UpdateCollectionVariables
	>(() => ({
		mutationKey: ["collections", "update"],
		mutationFn: ({ id, input }) => updateCollection(fetch, id, input),
		onSuccess: (_result, variables) => {
			invalidateQueryKeys(
				collectionUpdateSuccessInvalidationKeys({
					pathname: currentPathname,
					selectedCollectionId,
					updatedCollectionId: variables.id,
				})
			);
			closeDialog();
		},
	}));
	const deleteCollectionMutation = createMutation<
		{ deleted_id: string },
		Error,
		DeleteCollectionVariables
	>(() => ({
		mutationKey: ["collections", "delete"],
		mutationFn: ({ collection }) => deleteCollection(fetch, collection.id),
		onSuccess: (result) => {
			const effects = collectionDeleteSuccessEffects({
				pathname: currentPathname,
				selectedCollectionId,
				deletedCollectionId: result.deleted_id,
			});
			invalidateQueryKeys(effects.invalidateQueryKeys);
			if (effects.navigateTo) void goto(effects.navigateTo);
		},
	}));

	const collections = $derived(collectionsQuery.data?.collections ?? []);
	const schemas = $derived(schemasQuery.data ?? []);
	const dialogInitialValue = $derived(collectionDialogInitialValue(editingCollection));
	const dialogPending = $derived(
		collectionDialogPending(createCollectionMutation.isPending, updateCollectionMutation.isPending)
	);
	const dialogError = $derived(
		collectionDialogError(
			dialogMode,
			createCollectionMutation.error ?? null,
			updateCollectionMutation.error ?? null
		)
	);
	const collectionStatus = $derived(
		collectionListStatus({
			isLoading: collectionsQuery.isLoading,
			isError: collectionsQuery.isError,
			collectionCount: collections.length,
		})
	);
	const collectionsOverflow = $derived(collectionListOverflows(collections.length));
	const allDocumentsActive = $derived(isAllDocumentsActive(currentPathname, selectedCollectionId));

	function invalidateQueryKeys(queryKeys: CollectionQueryKey[]) {
		for (const queryKey of queryKeys) {
			void queryClient.invalidateQueries({ queryKey: [...queryKey] });
		}
	}

	function applyDialogState(state: CollectionDialogState) {
		dialogOpen = state.open;
		dialogMode = state.mode;
		editingCollection = state.editingCollection;
	}

	function openCreateDialog() {
		createCollectionMutation.reset();
		updateCollectionMutation.reset();
		applyDialogState(openCreateCollectionDialogState());
	}

	function openEditDialog(collection: CollectionResponse) {
		createCollectionMutation.reset();
		updateCollectionMutation.reset();
		applyDialogState(openEditCollectionDialogState(collection));
	}

	function closeDialog() {
		applyDialogState(closeCollectionDialogState(dialogMode));
	}

	function submitCollection(value: CollectionInput) {
		const action = collectionSubmitAction(dialogMode, editingCollection, value);
		if (action.type === "create") createCollectionMutation.mutate(action.input);
		if (action.type === "update") {
			updateCollectionMutation.mutate({ id: action.id, input: action.input });
		}
	}

	async function runDelete(collection: CollectionResponse) {
		try {
			await deleteCollectionMutation.mutateAsync({ collection });
		} catch {
			// The mutation owns the error state; the menu stays usable after the dialog closes.
		}
	}

	function confirmCollectionDelete(collection: CollectionResponse) {
		deleteCollectionMutation.reset();
		confirmDelete({
			title: m.documents_delete_collection_title(),
			description: m.documents_delete_collection_description({ name: collection.name }),
			confirm: { text: m.documents_delete() },
			onConfirm: () => runDelete(collection),
		});
	}

	function isCollectionActive(collection: CollectionResponse) {
		return isCollectionActiveByRoute(currentPathname, selectedCollectionId, collection.id);
	}

</script>

<Sidebar.Group class="group-data-[collapsible=icon]:hidden">
	<Sidebar.GroupLabel>{m.documents_collections_nav_label()}</Sidebar.GroupLabel>
	<Sidebar.GroupAction
		type="button"
		aria-label={m.documents_add_collection()}
		title={m.documents_add_collection()}
		onclick={openCreateDialog}
	>
		<PlusIcon />
		<span class="sr-only">{m.documents_add_collection()}</span>
	</Sidebar.GroupAction>
	<Sidebar.GroupContent>
		<Sidebar.Menu>
			<Sidebar.MenuItem>
				<Sidebar.MenuButton isActive={allDocumentsActive} tooltipContent={m.documents_all_documents()}>
					{#snippet child({ props })}
						<a {...props} href="/app/documents">
							<FilesIcon />
							<span>{m.documents_all_documents()}</span>
						</a>
					{/snippet}
				</Sidebar.MenuButton>
			</Sidebar.MenuItem>
		</Sidebar.Menu>

		<Sidebar.Menu
			aria-label={m.documents_collections_nav_label()}
			class={["mt-1", collectionsOverflow && "max-h-[22.25rem] overflow-y-auto pr-1"]}
		>
			{#if collectionStatus === "loading"}
				{#each COLLECTION_LOADING_ROWS as row (row)}
					<Sidebar.MenuSkeleton showIcon />
				{/each}
			{:else if collectionStatus === "error"}
				<Sidebar.MenuItem>
					<Sidebar.MenuButton
						class="text-destructive hover:text-destructive"
						onclick={() => void retryCollections(() => collectionsQuery.refetch())}
					>
						<AlertCircleIcon />
						<span>{m.documents_retry_collections()}</span>
					</Sidebar.MenuButton>
				</Sidebar.MenuItem>
			{:else if collectionStatus === "empty"}
				<Sidebar.MenuItem>
					<Sidebar.MenuButton aria-disabled={true} class="text-sidebar-foreground/70">
						<FolderIcon class="text-sidebar-foreground/70" />
						<span>{m.documents_no_collections()}</span>
					</Sidebar.MenuButton>
				</Sidebar.MenuItem>
			{:else}
				{#each collections as collection (collection.id)}
					<Sidebar.MenuItem>
						<Sidebar.MenuButton
							isActive={isCollectionActive(collection)}
							tooltipContent={collection.name}
						>
							{#snippet child({ props })}
								<a {...props} href={collectionHref(collection.id)}>
									<FolderIcon />
									<span>{collection.name}</span>
								</a>
							{/snippet}
						</Sidebar.MenuButton>
						<DropdownMenu.Root>
							<DropdownMenu.Trigger>
								{#snippet child({ props })}
									<Sidebar.MenuAction
										{...props}
										type="button"
										showOnHover
										class="data-[state=open]:bg-accent rounded-sm"
										onclick={composeCollectionMenuTriggerClick(props.onclick)}
									>
										<EllipsisIcon />
										<span class="sr-only">{m.documents_collection_actions()}</span>
									</Sidebar.MenuAction>
								{/snippet}
							</DropdownMenu.Trigger>
							<DropdownMenu.Content
								class="w-28 rounded-lg"
								side={sidebar.isMobile ? "bottom" : "right"}
								align={sidebar.isMobile ? "end" : "start"}
							>
								<DropdownMenu.Item
									onclick={(event) => runCollectionMenuAction(event, () => openEditDialog(collection))}
								>
									<PencilIcon />
									<span>{m.documents_edit()}</span>
								</DropdownMenu.Item>
								<DropdownMenu.Item
									variant="destructive"
									onclick={(event) =>
										runCollectionMenuAction(event, () => confirmCollectionDelete(collection))}
								>
									<Trash2Icon />
									<span>{m.documents_delete()}</span>
								</DropdownMenu.Item>
							</DropdownMenu.Content>
						</DropdownMenu.Root>
					</Sidebar.MenuItem>
				{/each}
			{/if}
		</Sidebar.Menu>

		{#if deleteCollectionMutation.isError}
			<Sidebar.Menu class="mt-1">
				<Sidebar.MenuItem>
					<Sidebar.MenuButton aria-disabled={true} class="text-destructive">
						<AlertCircleIcon />
						<span>{m.documents_delete_failed()}</span>
					</Sidebar.MenuButton>
				</Sidebar.MenuItem>
			</Sidebar.Menu>
		{/if}
	</Sidebar.GroupContent>
</Sidebar.Group>

<CollectionDialog
	bind:open={dialogOpen}
	mode={dialogMode}
	initialValue={dialogInitialValue}
	{schemas}
	schemasLoading={schemasQuery.isLoading}
	schemasError={schemasQuery.error ?? null}
	pending={dialogPending}
	error={dialogError}
	onSubmit={submitCollection}
/>
