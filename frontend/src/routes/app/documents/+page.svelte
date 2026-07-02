<script lang="ts">
	import AlertCircleIcon from '@lucide/svelte/icons/alert-circle';
	import ArchiveIcon from '@lucide/svelte/icons/archive';
	import FileTextIcon from '@lucide/svelte/icons/file-text';
	import FolderPlusIcon from '@lucide/svelte/icons/folder-plus';
	import GitBranchIcon from '@lucide/svelte/icons/git-branch';
	import LoaderCircleIcon from '@lucide/svelte/icons/loader-circle';
	import PlusIcon from '@lucide/svelte/icons/plus';
	import SaveIcon from '@lucide/svelte/icons/save';
	import { createMutation, createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Separator } from '$lib/components/ui/separator';
	import type { PageProps } from './$types';
	import {
		DOCUMENT_FOLDER_CONTENTS_QUERY_KEY,
		archiveDocument,
		archiveDocumentFolder,
		createDocumentFolder,
		documentFolderContentsQueryKey,
		documentFoldersQueryKey,
		fetchDocumentFolderContents,
		fetchDocumentFolderTree,
		moveDocumentFolder,
		updateDocument,
		updateDocumentFolder,
		uploadDocument,
		type ArchiveDocumentFolderResponse,
		type ArchiveDocumentFolderVariables,
		type ArchiveDocumentResponse,
		type ArchiveDocumentVariables,
		type DocumentFolderContentsResponse,
		type DocumentFolderInput,
		type DocumentFolderTreeResponse,
		type MoveDocumentFolderVariables,
		type UpdateDocumentFolderVariables,
		type UpdateDocumentVariables,
		type UploadDocumentVariables
	} from './api';
	import FolderTree from './folder-tree.svelte';
	import RepositoryTable from './repository-table.svelte';
	import {
		collectFolderMoveTargets,
		findFolder,
		flattenFolderTree,
		repositoryRows,
		selectInitialFolder,
		type DocumentFolderNode,
		type RepositoryDocument
	} from './tree';
	import UploadPanel from './upload-panel.svelte';
	import {
		filesToUploadItems,
		markUploadFailed,
		markUploadUploaded,
		markUploadUploading,
		type UploadQueueItem
	} from './upload-queue';
	import {
		ORGANIZATION_UNITS_QUERY_KEY,
		fetchOrganizationUnitTree,
		type OrganizationUnitListResponse
	} from '../organization-units/api';
	import { flattenUnitTree } from '../organization-units/tree';

	type DocumentsPageData = {
		canViewDocuments: boolean;
		canCreateDocuments: boolean;
		canUpdateDocuments: boolean;
		canDeleteDocuments: boolean;
		canDownloadDocuments: boolean;
		selectedOrganizationUnitId: string | null;
	};

	type PageArchiveDocumentFolderVariables = ArchiveDocumentFolderVariables & {
		organizationUnitId: string;
		wasSelected: boolean;
	};
	type PageArchiveDocumentVariables = ArchiveDocumentVariables & { folderId: string };

	let { data }: PageProps = $props();

	const pageData = $derived(data as DocumentsPageData);
	const queryClient = useQueryClient();

	let organizationUnitIdOverride = $state<string | null>(null);
	let selectedFolderOverride = $state<string | null>(null);
	let rootFolderName = $state('');
	let rootFolderDescription = $state('');
	let childFolderName = $state('');
	let childFolderDescription = $state('');
	let localError = $state('');
	let uploadQueue = $state<UploadQueueItem[]>([]);
	let isUploadingQueue = $state(false);

	const fieldClass =
		'h-9 rounded-md border bg-background px-3 text-sm disabled:cursor-not-allowed disabled:opacity-60';

	const selectedOrganizationUnitId = $derived(
		(organizationUnitIdOverride ?? pageData.selectedOrganizationUnitId ?? '').trim()
	);
	const hasDocumentAccess = $derived(pageData.canViewDocuments);

	const organizationUnitsQuery = createQuery<OrganizationUnitListResponse, Error>(() => ({
		queryKey: ORGANIZATION_UNITS_QUERY_KEY,
		queryFn: () => fetchOrganizationUnitTree(fetch),
		enabled: pageData.canViewDocuments
	}));
	const folderTreeQuery = createQuery<DocumentFolderTreeResponse, Error>(() => ({
		queryKey: documentFoldersQueryKey(selectedOrganizationUnitId),
		queryFn: () => fetchDocumentFolderTree(fetch, selectedOrganizationUnitId),
		enabled: pageData.canViewDocuments && selectedOrganizationUnitId.length > 0
	}));
	const organizationUnitOptions = $derived(flattenUnitTree(organizationUnitsQuery.data?.units ?? []));
	const folderTree = $derived(folderTreeQuery.data?.folders ?? []);
	const flatFolders = $derived(flattenFolderTree(folderTree));
	const selectedFolder = $derived(selectInitialFolder(folderTree, selectedFolderOverride));
	const selectedFolderId = $derived(selectedFolder?.id ?? null);
	const folderContentsQuery = createQuery<DocumentFolderContentsResponse, Error>(() => ({
		queryKey: documentFolderContentsQueryKey(selectedFolderId),
		queryFn: () => fetchDocumentFolderContents(fetch, selectedFolderId ?? ''),
		enabled: pageData.canViewDocuments && selectedFolderId !== null
	}));

	const createFolderMutationState = createMutation<DocumentFolderNode, Error, DocumentFolderInput>(
		() => ({
			mutationKey: ['document-folders', 'create'],
			mutationFn: (input) => createDocumentFolder(fetch, input),
			onSuccess: async (folder, variables) => {
				selectFolder(folder.id);
				await Promise.all([
					invalidateFolderTree(variables.organizationUnitId),
					invalidateFolderContents(variables.parentId)
				]);
			}
		})
	);
	const updateFolderMutationState = createMutation<
		DocumentFolderNode,
		Error,
		UpdateDocumentFolderVariables
	>(() => ({
		mutationKey: ['document-folders', 'update'],
		mutationFn: (variables) => updateDocumentFolder(fetch, variables),
		onSuccess: async (folder) => {
			selectFolder(folder.id);
			await Promise.all([
				invalidateFolderTree(folder.organizationUnitId),
				invalidateFolderContents(folder.id),
				invalidateFolderContents(folder.parentId)
			]);
		}
	}));
	const moveFolderMutationState = createMutation<
		DocumentFolderNode,
		Error,
		MoveDocumentFolderVariables
	>(() => ({
		mutationKey: ['document-folders', 'move'],
		mutationFn: (variables) => moveDocumentFolder(fetch, variables),
		onSuccess: async (folder) => {
			selectFolder(folder.id);
			await Promise.all([invalidateFolderTree(folder.organizationUnitId), invalidateAllFolderContents()]);
		}
	}));
	const archiveFolderMutationState = createMutation<
		ArchiveDocumentFolderResponse,
		Error,
		PageArchiveDocumentFolderVariables
	>(() => ({
		mutationKey: ['document-folders', 'archive'],
		mutationFn: (variables) => archiveDocumentFolder(fetch, variables),
		onSuccess: async (_response, variables) => {
			if (variables.wasSelected) {
				selectFolder(null);
			}
			await Promise.all([
				invalidateFolderTree(variables.organizationUnitId),
				invalidateAllFolderContents()
			]);
		}
	}));
	const uploadMutationState = createMutation<RepositoryDocument, Error, UploadDocumentVariables>(
		() => ({
			mutationKey: ['documents', 'upload'],
			mutationFn: (variables) => uploadDocument(fetch, variables),
			onSuccess: async (_document, variables) => {
				await invalidateFolderContents(variables.folderId);
			}
		})
	);
	const updateDocumentMutationState = createMutation<
		RepositoryDocument,
		Error,
		UpdateDocumentVariables
	>(() => ({
		mutationKey: ['documents', 'update'],
		mutationFn: (variables) => updateDocument(fetch, variables),
		onSuccess: async (document) => {
			await invalidateFolderContents(document.folderId);
		}
	}));
	const archiveDocumentMutationState = createMutation<
		ArchiveDocumentResponse,
		Error,
		PageArchiveDocumentVariables
	>(() => ({
		mutationKey: ['documents', 'archive'],
		mutationFn: (variables) => archiveDocument(fetch, variables),
		onSuccess: async (_response, variables) => {
			await invalidateFolderContents(variables.folderId);
		}
	}));

	const moveTargets = $derived(
		selectedFolderId ? collectFolderMoveTargets(folderTree, selectedFolderId) : flatFolders
	);
	const contentRows = $derived(
		repositoryRows(folderContentsQuery.data?.folders ?? [], folderContentsQuery.data?.documents ?? [])
	);
	const activeRepositoryCount = $derived.by(() => {
		if (!hasDocumentAccess || !selectedOrganizationUnitId) return 0;
		if (selectedFolderId) return contentRows.length;
		return flatFolders.length;
	});
	const activeRepositoryCountLabel = $derived.by(() => {
		if (!hasDocumentAccess || !selectedOrganizationUnitId || selectedFolderId) {
			return `${activeRepositoryCount} active items`;
		}

		return `${activeRepositoryCount} active folders`;
	});
	let selectedFolderName = $derived(selectedFolder?.name ?? '');
	let selectedFolderDescription = $derived(selectedFolder?.description ?? '');
	let selectedFolderParentId = $derived(selectedFolder?.parentId ?? '');
	const isMutationPending = $derived(
		createFolderMutationState.isPending ||
			updateFolderMutationState.isPending ||
			moveFolderMutationState.isPending ||
			archiveFolderMutationState.isPending ||
			uploadMutationState.isPending ||
			updateDocumentMutationState.isPending ||
			archiveDocumentMutationState.isPending ||
			isUploadingQueue
	);
	const mutationError = $derived.by(
		() =>
			localError ||
			createFolderMutationState.error?.message ||
			updateFolderMutationState.error?.message ||
			moveFolderMutationState.error?.message ||
			archiveFolderMutationState.error?.message ||
			uploadMutationState.error?.message ||
			updateDocumentMutationState.error?.message ||
			archiveDocumentMutationState.error?.message ||
			''
	);

	function handleOrganizationUnitInput(event: Event) {
		const input = event.currentTarget as HTMLInputElement;
		setOrganizationUnit(input.value);
	}

	function handleOrganizationUnitSelect(event: Event) {
		const select = event.currentTarget as HTMLSelectElement;
		setOrganizationUnit(select.value);
	}

	function setOrganizationUnit(id: string) {
		organizationUnitIdOverride = id;
		selectFolder(null);
		resetMutationErrors();
	}

	function selectFolder(id: string | null) {
		if (id !== selectedFolderId) {
			uploadQueue = [];
		}
		selectedFolderOverride = id;
	}

	async function submitRootFolder(event: SubmitEvent) {
		event.preventDefault();

		try {
			await runCreateFolder(null, rootFolderName, rootFolderDescription);
			rootFolderName = '';
			rootFolderDescription = '';
		} catch {
		}
	}

	async function submitChildFolder(event: SubmitEvent) {
		event.preventDefault();
		if (!selectedFolderId) return;

		try {
			await runCreateFolder(selectedFolderId, childFolderName, childFolderDescription);
			childFolderName = '';
			childFolderDescription = '';
		} catch {
		}
	}

	async function submitSelectedFolderUpdate(event: SubmitEvent) {
		event.preventDefault();
		if (!selectedFolder) return;

		try {
			await runRenameFolder(selectedFolder.id, selectedFolderName, selectedFolderDescription);
		} catch {
		}
	}

	async function submitSelectedFolderMove(event: SubmitEvent) {
		event.preventDefault();
		if (!selectedFolder) return;

		try {
			await runMoveFolder(selectedFolder.id, selectedFolderParentId || null);
		} catch {
		}
	}

	async function confirmArchiveSelectedFolder() {
		if (!selectedFolder || isMutationPending) return;
		if (!confirm(`Archive folder "${selectedFolder.name}"?`)) return;

		try {
			await runArchiveFolder(selectedFolder.id);
		} catch {
		}
	}

	async function runCreateFolder(parentId: string | null, name: string, description: string) {
		resetMutationErrors();
		const organizationUnitId = selectedOrganizationUnitId;
		if (!organizationUnitId) {
			localError = 'Organization unit is required';
			throw new Error(localError);
		}

		await createFolderMutationState.mutateAsync({
			organizationUnitId,
			parentId,
			name: requireName(name, 'Folder name is required'),
			description: description.trim() || null
		});
	}

	async function runRenameFolder(id: string, name: string, description?: string | null) {
		resetMutationErrors();
		const folder = findFolder(folderTree, id);
		if (!folder) {
			localError = 'Folder was not found';
			throw new Error(localError);
		}

		await updateFolderMutationState.mutateAsync({
			id,
			input: {
				organizationUnitId: folder.organizationUnitId,
				parentId: folder.parentId ?? null,
				name: requireName(name, 'Folder name is required'),
				description: description ?? folder.description ?? null
			}
		});
	}

	async function runMoveFolder(id: string, parentId: string | null) {
		resetMutationErrors();
		await moveFolderMutationState.mutateAsync({ id, parentId });
	}

	async function runArchiveFolder(id: string) {
		resetMutationErrors();
		const folder = findFolder(folderTree, id);
		if (!folder) {
			localError = 'Folder was not found';
			throw new Error(localError);
		}

		await archiveFolderMutationState.mutateAsync({
			id,
			organizationUnitId: folder.organizationUnitId,
			wasSelected: selectedFolderId === id
		});
	}

	async function runRenameDocument(id: string, displayName: string) {
		resetMutationErrors();
		await updateDocumentMutationState.mutateAsync({
			id,
			input: { displayName: requireName(displayName, 'Document name is required') }
		});
	}

	async function runArchiveDocument(id: string) {
		resetMutationErrors();
		const row = contentRows.find((row) => row.type === 'document' && row.id === id);
		if (!row || row.type !== 'document') {
			localError = 'Document was not found';
			throw new Error(localError);
		}

		await archiveDocumentMutationState.mutateAsync({ id, folderId: row.document.folderId });
	}

	function handleFilesSelected(files: FileList) {
		uploadQueue = [...uploadQueue, ...filesToUploadItems(files)];
	}

	async function runUploadQueue() {
		if (isMutationPending || isUploadingQueue) return;

		resetMutationErrors();
		const folderId = selectedFolderId;
		if (!folderId) return;

		const pendingItems = uploadQueue.filter(
			(item) => item.status === 'queued' || item.status === 'failed'
		);
		if (pendingItems.length === 0) return;

		isUploadingQueue = true;
		try {
			for (const item of pendingItems) {
				uploadQueue = markUploadUploading(uploadQueue, item.id);
				try {
					const document = await uploadMutationState.mutateAsync({ folderId, file: item.file });
					uploadQueue = markUploadUploaded(uploadQueue, item.id, document.id);
				} catch (error) {
					uploadQueue = markUploadFailed(uploadQueue, item.id, errorMessage(error, 'Upload failed'));
				}
			}
		} finally {
			isUploadingQueue = false;
		}
	}

	function requireName(value: string, message: string) {
		const trimmed = value.trim();
		if (!trimmed) {
			localError = message;
			throw new Error(message);
		}
		return trimmed;
	}

	function resetMutationErrors() {
		localError = '';
		createFolderMutationState.reset();
		updateFolderMutationState.reset();
		moveFolderMutationState.reset();
		archiveFolderMutationState.reset();
		uploadMutationState.reset();
		updateDocumentMutationState.reset();
		archiveDocumentMutationState.reset();
	}

	async function invalidateFolderTree(organizationUnitId: string) {
		if (!organizationUnitId.trim()) return;
		await queryClient.invalidateQueries({ queryKey: documentFoldersQueryKey(organizationUnitId) });
	}

	async function invalidateFolderContents(folderId: string | null | undefined) {
		if (!folderId) return;
		await queryClient.invalidateQueries({ queryKey: documentFolderContentsQueryKey(folderId) });
	}

	async function invalidateAllFolderContents() {
		await queryClient.invalidateQueries({ queryKey: DOCUMENT_FOLDER_CONTENTS_QUERY_KEY });
	}

	function errorMessage(error: unknown, fallback: string) {
		return error instanceof Error ? error.message : fallback;
	}
</script>

<svelte:head>
	<title>Documents | Syncra DMS</title>
</svelte:head>

<div class="flex flex-1 flex-col gap-4 p-4 lg:p-6">
	<div class="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
		<div class="min-w-0">
			<div class="flex items-center gap-2">
				<FileTextIcon class="size-5 text-primary" />
				<h2 class="truncate text-xl font-semibold tracking-normal">Documents</h2>
			</div>
			<div class="mt-1 flex flex-wrap items-center gap-x-3 gap-y-1 text-sm text-muted-foreground">
				<span>{activeRepositoryCountLabel}</span>
				<span>
					{#if selectedOrganizationUnitId}
						Organization unit {selectedOrganizationUnitId}
					{:else}
						No organization unit selected
					{/if}
				</span>
			</div>
		</div>
		<Badge variant={hasDocumentAccess ? 'secondary' : 'outline'}>
			{hasDocumentAccess ? 'Repository' : 'No access'}
		</Badge>
	</div>

	{#if !hasDocumentAccess}
		<section class="rounded-md border border-dashed p-4" role="status">
			<h3 class="text-sm font-medium">No document access</h3>
			<p class="mt-1 text-sm text-muted-foreground">
				This account does not have permission to view documents.
			</p>
		</section>
	{:else}
		<Card.Root>
			<Card.Content class="grid gap-3 p-4 md:grid-cols-[minmax(0,1fr)_minmax(14rem,18rem)]">
				<label class="grid gap-1.5 text-sm font-medium">
					Organization unit id
					<input
						class={fieldClass}
						value={selectedOrganizationUnitId}
						placeholder="Enter organization unit id"
						oninput={handleOrganizationUnitInput}
					/>
				</label>
				<label class="grid gap-1.5 text-sm font-medium">
					Unit selector
					<select
						class={fieldClass}
						value={selectedOrganizationUnitId}
						disabled={organizationUnitsQuery.isLoading}
						onchange={handleOrganizationUnitSelect}
					>
						<option value="">Select unit</option>
						{#each organizationUnitOptions as unit (unit.id)}
							<option value={unit.id}>
								{`${'- '.repeat(unit.depth)}${unit.name}${unit.code ? ` (${unit.code})` : ''}`}
							</option>
						{/each}
					</select>
				</label>
			</Card.Content>
		</Card.Root>

		{#if organizationUnitsQuery.isError}
			<div
				class="flex flex-wrap items-center gap-2 rounded-md border border-destructive/30 bg-destructive/10 p-3 text-sm text-destructive"
				role="alert"
			>
				<AlertCircleIcon class="size-4 shrink-0" />
				<span>{organizationUnitsQuery.error.message}</span>
				<Button
					type="button"
					variant="outline"
					size="xs"
					class="ms-auto"
					onclick={() => organizationUnitsQuery.refetch()}
				>
					Retry
				</Button>
			</div>
		{/if}

		{#if mutationError}
			<div
				class="flex items-center gap-2 rounded-md border border-destructive/30 bg-destructive/10 p-3 text-sm text-destructive"
				role="alert"
			>
				<AlertCircleIcon class="size-4 shrink-0" />
				<span>{mutationError}</span>
			</div>
		{/if}

		{#if !selectedOrganizationUnitId}
			<section
				class="flex min-h-32 items-center justify-center rounded-md border border-dashed bg-muted/20 px-4 text-sm text-muted-foreground"
				role="status"
			>
				Select an organization unit to open the repository.
			</section>
		{:else}
			{#if folderTreeQuery.isError}
				<div
					class="flex flex-wrap items-center gap-2 rounded-md border border-destructive/30 bg-destructive/10 p-3 text-sm text-destructive"
					role="alert"
				>
					<AlertCircleIcon class="size-4 shrink-0" />
					<span>{folderTreeQuery.error.message}</span>
					<Button
						type="button"
						variant="outline"
						size="xs"
						class="ms-auto"
						onclick={() => folderTreeQuery.refetch()}
					>
						Retry
					</Button>
				</div>
			{/if}

			{#if folderTreeQuery.isLoading || folderTreeQuery.isFetching}
				<div
					class="flex items-center gap-2 rounded-md border bg-muted/30 p-3 text-sm text-muted-foreground"
					role="status"
					aria-live="polite"
				>
					<LoaderCircleIcon class="size-4 shrink-0 animate-spin" />
					<span>Loading folders</span>
				</div>
			{/if}

			<div class="grid gap-4 xl:grid-cols-[20rem_minmax(0,1fr)]">
				<Card.Root>
					<Card.Header>
						<Card.Title>Folders</Card.Title>
						<Card.Description>{flatFolders.length} visible</Card.Description>
					</Card.Header>
					<Card.Content class="grid gap-4">
						<FolderTree
							folders={flatFolders}
							selectedId={selectedFolderId}
							onSelect={selectFolder}
						/>

						{#if pageData.canCreateDocuments && folderTree.length === 0 && !folderTreeQuery.isLoading}
							<form class="grid gap-3" onsubmit={submitRootFolder}>
								<Separator />
								<label class="grid gap-1.5 text-sm font-medium">
									Root folder
									<input
										class={fieldClass}
										bind:value={rootFolderName}
										placeholder="Folder name"
										disabled={isMutationPending}
										required
									/>
								</label>
								<label class="grid gap-1.5 text-sm font-medium">
									Description
									<input
										class={fieldClass}
										bind:value={rootFolderDescription}
										placeholder="Optional"
										disabled={isMutationPending}
									/>
								</label>
								<Button type="submit" size="sm" class="w-fit gap-2" disabled={isMutationPending}>
									<FolderPlusIcon class="size-4" />
									Create root
								</Button>
							</form>
						{/if}
					</Card.Content>
				</Card.Root>

				<div class="grid min-w-0 gap-4">
					<Card.Root>
						<Card.Header>
							<div class="flex min-w-0 items-start justify-between gap-3">
								<div class="min-w-0">
									<Card.Title class="truncate">
										{selectedFolder ? selectedFolder.name : 'No folder selected'}
									</Card.Title>
									<Card.Description>Folder controls</Card.Description>
								</div>
								<Badge variant={pageData.canUpdateDocuments ? 'secondary' : 'outline'}>
									{pageData.canUpdateDocuments ? 'Editable' : 'Read only'}
								</Badge>
							</div>
						</Card.Header>
						<Card.Content class="grid gap-4">
							{#if selectedFolder}
								{#if pageData.canUpdateDocuments}
									<form class="grid gap-3" onsubmit={submitSelectedFolderUpdate}>
										<div class="grid gap-3 md:grid-cols-[minmax(0,1fr)_minmax(0,1fr)_auto] md:items-end">
											<label class="grid gap-1.5 text-sm font-medium">
												Name
												<input
													class={fieldClass}
													bind:value={selectedFolderName}
													disabled={isMutationPending}
													required
												/>
											</label>
											<label class="grid gap-1.5 text-sm font-medium">
												Description
												<input
													class={fieldClass}
													bind:value={selectedFolderDescription}
													disabled={isMutationPending}
												/>
											</label>
											<Button type="submit" size="sm" class="gap-2" disabled={isMutationPending}>
												<SaveIcon class="size-4" />
												Save
											</Button>
										</div>
									</form>
								{/if}

								<dl class="grid gap-2 text-xs text-muted-foreground sm:grid-cols-3">
									<div>
										<dt class="font-medium text-foreground">Created</dt>
										<dd>{selectedFolder.createdAt}</dd>
									</div>
									<div>
										<dt class="font-medium text-foreground">Updated</dt>
										<dd>{selectedFolder.updatedAt}</dd>
									</div>
									<div>
										<dt class="font-medium text-foreground">Folder id</dt>
										<dd class="truncate" title={selectedFolder.id}>{selectedFolder.id}</dd>
									</div>
								</dl>

								{#if pageData.canCreateDocuments || pageData.canUpdateDocuments || pageData.canDeleteDocuments}
									<Separator />
								{/if}

								{#if pageData.canCreateDocuments}
									<form class="grid gap-3" onsubmit={submitChildFolder}>
										<div class="grid gap-3 md:grid-cols-[minmax(0,1fr)_minmax(0,1fr)_auto] md:items-end">
											<label class="grid gap-1.5 text-sm font-medium">
												New folder
												<input
													class={fieldClass}
													bind:value={childFolderName}
													placeholder="Folder name"
													disabled={isMutationPending}
													required
												/>
											</label>
											<label class="grid gap-1.5 text-sm font-medium">
												Description
												<input
													class={fieldClass}
													bind:value={childFolderDescription}
													placeholder="Optional"
													disabled={isMutationPending}
												/>
											</label>
											<Button
												type="submit"
												variant="outline"
												size="sm"
												class="gap-2"
												disabled={isMutationPending}
											>
												<PlusIcon class="size-4" />
												Create
											</Button>
										</div>
									</form>
								{/if}

								{#if pageData.canUpdateDocuments}
									<form class="grid gap-3" onsubmit={submitSelectedFolderMove}>
										<div class="grid gap-3 md:grid-cols-[minmax(0,1fr)_auto] md:items-end">
											<label class="grid gap-1.5 text-sm font-medium">
												Parent
												<select
													class={fieldClass}
													bind:value={selectedFolderParentId}
													disabled={isMutationPending}
												>
													<option value="">Root</option>
													{#each moveTargets as option (option.id)}
														<option value={option.id}>
															{`${'- '.repeat(option.depth)}${option.name}`}
														</option>
													{/each}
												</select>
											</label>
											<Button
												type="submit"
												variant="outline"
												size="sm"
												class="gap-2"
												disabled={isMutationPending}
											>
												<GitBranchIcon class="size-4" />
												Move
											</Button>
										</div>
									</form>
								{/if}

								{#if pageData.canDeleteDocuments}
									<Button
										type="button"
										variant="destructive"
										size="sm"
										class="w-fit gap-2"
										disabled={isMutationPending}
										onclick={confirmArchiveSelectedFolder}
									>
										<ArchiveIcon class="size-4" />
										Archive folder
									</Button>
								{/if}
							{:else}
								<div
									class="flex min-h-24 items-center justify-center rounded-md border border-dashed bg-muted/20 px-4 text-sm text-muted-foreground"
								>
									No folder selected
								</div>
							{/if}
						</Card.Content>
					</Card.Root>

					<Card.Root>
						<Card.Header>
							<div class="flex min-w-0 items-start justify-between gap-3">
								<div class="min-w-0">
									<Card.Title>Repository</Card.Title>
									<Card.Description>
										{#if selectedFolder}
											{contentRows.length} items
										{:else}
											Select a folder
										{/if}
									</Card.Description>
								</div>
								{#if folderContentsQuery.isFetching}
									<LoaderCircleIcon class="size-4 animate-spin text-muted-foreground" />
								{/if}
							</div>
						</Card.Header>
						<Card.Content class="grid gap-3">
							{#if folderContentsQuery.isError}
								<div
									class="flex flex-wrap items-center gap-2 rounded-md border border-destructive/30 bg-destructive/10 p-3 text-sm text-destructive"
									role="alert"
								>
									<AlertCircleIcon class="size-4 shrink-0" />
									<span>{folderContentsQuery.error.message}</span>
									<Button
										type="button"
										variant="outline"
										size="xs"
										class="ms-auto"
										onclick={() => folderContentsQuery.refetch()}
									>
										Retry
									</Button>
								</div>
							{/if}

							<RepositoryTable
								rows={contentRows}
								canUpdate={pageData.canUpdateDocuments}
								canDelete={pageData.canDeleteDocuments}
								canDownload={pageData.canDownloadDocuments}
								isPending={isMutationPending || folderContentsQuery.isLoading}
								onOpenFolder={selectFolder}
								onRenameFolder={runRenameFolder}
								onArchiveFolder={runArchiveFolder}
								onRenameDocument={runRenameDocument}
								onArchiveDocument={runArchiveDocument}
							/>
						</Card.Content>
					</Card.Root>

					<UploadPanel
						canCreate={pageData.canCreateDocuments}
						{selectedFolderId}
						queue={uploadQueue}
						isUploading={isUploadingQueue}
						isPending={isMutationPending}
						onFilesSelected={handleFilesSelected}
						onUpload={runUploadQueue}
					/>
				</div>
			</div>
		{/if}
	{/if}
</div>
