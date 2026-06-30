<script lang="ts">
	import CalendarIcon from "@lucide/svelte/icons/calendar";
	import ChevronLeftIcon from "@lucide/svelte/icons/chevron-left";
	import ChevronRightIcon from "@lucide/svelte/icons/chevron-right";
	import DownloadIcon from "@lucide/svelte/icons/download";
	import FolderIcon from "@lucide/svelte/icons/folder";
	import LoaderIcon from "@lucide/svelte/icons/loader-circle";
	import SearchIcon from "@lucide/svelte/icons/search";
	import Trash2Icon from "@lucide/svelte/icons/trash-2";
	import UploadIcon from "@lucide/svelte/icons/upload";
	import { createMutation, createQuery, useQueryClient } from "@tanstack/svelte-query";
	import { getCoreRowModel, type ColumnDef } from "@tanstack/table-core";
	import type { ComponentProps } from "svelte";
	import { toast } from "svelte-sonner";

	import { goto } from "$app/navigation";
	import { page } from "$app/state";
	import * as Alert from "$lib/components/ui/alert/index.js";
	import { Button } from "$lib/components/ui/button/index.js";
	import { confirmDelete } from "$lib/components/ui/confirm-delete-dialog/index.js";
	import { FlexRender, renderComponent } from "$lib/components/ui/data-table/index.js";
	import { createSvelteTable } from "$lib/components/ui/data-table/data-table.svelte.js";
	import { DocumentPreviewDialog } from "$lib/components/ui/document-preview/index.js";
	import { Input } from "$lib/components/ui/input/index.js";
	import * as Popover from "$lib/components/ui/popover/index.js";
	import * as RangeCalendar from "$lib/components/ui/range-calendar/index.js";
	import * as Select from "$lib/components/ui/select/index.js";
	import * as Table from "$lib/components/ui/table/index.js";
	import { fetchCollections, type CollectionListResponse } from "$lib/client/collections";
	import PlusIcon from "@lucide/svelte/icons/plus";
	
	import {
		fetchPersonalSchemaOptions,
		PERSONAL_SCHEMA_OPTIONS_QUERY_KEY,
		type PersonalSchemaOption,
	} from "$lib/client/schemas";
	import ActionsCell from "./actions-cell.svelte";
	import CollectionsCell from "./collections-cell.svelte";
	import { Spinner } from "$lib/components/ui/spinner/index.js";
	import { m } from "$lib/paraglide/messages.js";

	import {
		deleteOCRDocument,
		deleteOCRDocuments,
		downloadOCRDocuments,
		fetchOCRDocumentPreview,
		fetchOCRDocuments,
		isOCRDocumentNotFoundError,
		moveOCRDocuments,
		shouldRetryOCRDocumentsQuery,
		updateOCRDocument,
		type DeleteOCRDocumentsResponse,
		type DownloadFormat,
		type MoveOCRDocumentsResponse,
		type OCRDocumentPreview,
		type OCRDocumentListItemResponse,
		type OCRDocumentListResponse,
		type UpdateOCRDocumentResponse,
	} from "./api";
	import CreatedDateHeader from "./created-date-header.svelte";
	import DownloadDialog from "./download-dialog.svelte";
	import FilenameCell from "./filename-cell.svelte";
	import MoveCollectionsDialog from "./move-collections-dialog.svelte";
	import SelectionCheckbox from "./selection-checkbox.svelte";
	import {
		cursorNextState,
		cursorPreviousState,
		dateRangeToQueryBounds,
		formatCreatedDate,
		formatFileSize,
		headerSelectionState,
		resetCursorState,
		togglePageSelection,
		toggleSelection,
		type CursorState,
		type DateRangeValue,
		type SortDirection,
	} from "./table-utils";
	import {
		today,
		getLocalTimeZone,
		startOfWeek,
		endOfWeek,
		startOfMonth,
		endOfMonth,
	} from "@internationalized/date";

	const PAGE_SIZE_OPTIONS = [10, 20, 50, 100];
	type RangeCalendarValue = ComponentProps<typeof RangeCalendar.RangeCalendar>["value"];
	type DeleteDocumentsVariables = { ids: string[] };
	type MoveDocumentsVariables = { ids: string[]; collectionIds: string[] };
	type UpdateDocumentVariables = { id: string; originalFilename: string };
	type SelectedDocumentsById = Record<string, OCRDocumentListItemResponse>;

	let filenameFilter = $state("");
	let debouncedFilenameFilter = $state("");
	let appliedDateRange = $state<DateRangeValue | undefined>();
	let pendingDateRange = $state<RangeCalendarValue>();
	let datePopoverOpen = $state(false);
	let sortDirection = $state<SortDirection>("desc");
	let pageSize = $state(20);
	let cursorState = $state<CursorState>(resetCursorState());
	let selectedIds = $state<Set<string>>(new Set());
	let selectedDocumentsById = $state.raw<SelectedDocumentsById>({});
	let previousCollectionId = $state<string | undefined>();
	let previousSchemaId = $state<string | undefined>();
	let previewOpen = $state(false);
	let previewDocumentId = $state<string | null>(null);
	let previewFilename = $state<string | null>(null);
	let moveDialogOpen = $state(false);
	let updatingDocumentId = $state<string | null>(null);
	let renamingDocumentId = $state<string | null>(null);
	let downloadDialogOpen = $state(false);
	let downloadTargets = $state<OCRDocumentListItemResponse[]>([]);
	let downloadPending = $state(false);
	let downloadError = $state<Error | null>(null);

	const queryClient = useQueryClient();

	const normalizedFilename = $derived(debouncedFilenameFilter.trim() || undefined);
	const dateBounds = $derived(dateRangeToQueryBounds(appliedDateRange));
	const collectionId = $derived(page.url.searchParams.get("collection")?.trim() || undefined);
	const schemaId = $derived(page.url.searchParams.get("schema")?.trim() || undefined);
	const effectiveCursor = $derived(
		previousCollectionId === collectionId && previousSchemaId === schemaId ? cursorState.currentCursor : null
	);

	const documentsQuery = createQuery<OCRDocumentListResponse, Error>(() => ({
		queryKey: [
			"ocr-documents",
			{
				collectionId,
				schemaId,
				filename: normalizedFilename,
				createdFrom: dateBounds.createdFrom,
				createdTo: dateBounds.createdTo,
				cursor: effectiveCursor,
				size: pageSize,
				sort: sortDirection,
			},
		],
		queryFn: () =>
			fetchOCRDocuments(
				fetch,
				{
					collectionId,
					schemaId,
					filename: normalizedFilename,
					createdFrom: dateBounds.createdFrom,
					createdTo: dateBounds.createdTo,
					cursor: effectiveCursor,
					size: pageSize,
					sort: sortDirection,
				},
				m.documents_failed_load_documents()
			),
		retry: (failureCount, error) =>
			shouldRetryOCRDocumentsQuery(failureCount, error, { collectionId }),
	}));
	const previewQuery = createQuery<OCRDocumentPreview, Error>(() => ({
		queryKey: ["ocr-document", previewDocumentId],
		queryFn: () => {
			const id = previewDocumentId;
			if (!id) throw new Error(m.documents_missing_document_id());

			return fetchOCRDocumentPreview(fetch, id, m.documents_failed_load_document());
		},
		enabled: previewOpen && previewDocumentId !== null,
		staleTime: 1000 * 60 * 2,
	}));
	const collectionsQuery = createQuery<CollectionListResponse, Error>(() => ({
		queryKey: ["collections"],
		queryFn: () => fetchCollections(fetch),
	}));
	const schemasQuery = createQuery<PersonalSchemaOption[], Error>(() => ({
		queryKey: PERSONAL_SCHEMA_OPTIONS_QUERY_KEY,
		queryFn: () => fetchPersonalSchemaOptions(fetch),
	}));
	const deleteMutation = createMutation<DeleteOCRDocumentsResponse, Error, DeleteDocumentsVariables>(
		() => ({
			mutationKey: ["ocr-documents", "delete"],
			mutationFn: ({ ids }) => {
				if (ids.length === 1) {
					return deleteOCRDocument(fetch, ids[0], m.documents_failed_delete_document());
				}

				return deleteOCRDocuments(fetch, ids, m.documents_failed_delete_documents());
			},
			onSuccess: (result, variables) => {
				clearDeletedDocuments(result.deleted_ids, variables.ids);
				void queryClient.invalidateQueries({ queryKey: ["ocr-documents"] });
			},
		})
	);
	const moveMutation = createMutation<MoveOCRDocumentsResponse, Error, MoveDocumentsVariables>(() => ({
		mutationKey: ["ocr-documents", "move"],
		mutationFn: ({ ids, collectionIds }) =>
			moveOCRDocuments(fetch, ids, collectionIds, m.documents_failed_move_documents()),
		onSuccess: () => {
			selectedIds = new Set();
			selectedDocumentsById = {};
			moveDialogOpen = false;
			void queryClient.invalidateQueries({ queryKey: ["ocr-documents"] });
			void queryClient.invalidateQueries({ queryKey: ["collections"] });
		},
	}));
	const updateMutation = createMutation<UpdateOCRDocumentResponse, Error, UpdateDocumentVariables>(
		() => ({
			mutationKey: ["ocr-documents", "update"],
			mutationFn: ({ id, originalFilename }) =>
				updateOCRDocument(fetch, id, originalFilename, m.documents_failed_update_document()),
			onSuccess: (updated) => {
				updateDocumentCaches(updated);
				void queryClient.invalidateQueries({ queryKey: ["ocr-documents"] });
			},
		})
	);

	const documents = $derived(documentsQuery.data?.documents ?? []);
	const collections = $derived(collectionsQuery.data?.collections ?? []);
	const schemas = $derived(schemasQuery.data ?? []);
	const schemasMap = $derived(new Map<string, PersonalSchemaOption>(schemas.map((s) => [s.id, s])));
	const selectedCollectionName = $derived(
		collectionId
			? (collections.find((collection) => collection.id === collectionId)?.name ??
					m.documents_unknown_collection())
			: m.documents_all_collections()
	);
	const selectedSchemaName = $derived(
		schemaId
			? (schemas.find((schema) => schema.id === schemaId)?.name ?? m.documents_all_schemas())
			: m.documents_all_schemas()
	);
	const collectionNotFound = $derived(
		Boolean(collectionId) && isOCRDocumentNotFoundError(documentsQuery.error)
	);
	const nextCursor = $derived(documentsQuery.data?.next_cursor ?? null);
	const visibleIds = $derived(documents.map((document) => document.id));
	const pageSelection = $derived(headerSelectionState(visibleIds, selectedIds));
	const selectedDocuments = $derived(
		[...selectedIds]
			.map((id) => selectedDocumentsById[id])
			.filter((document): document is OCRDocumentListItemResponse => Boolean(document))
	);
	const dateRangeLabel = $derived.by(() => {
		if (!appliedDateRange?.start && !appliedDateRange?.end) return m.documents_date_range();

		const start = appliedDateRange.start
			? appliedDateRange.start.toString()
			: m.documents_any_date();
		const end = appliedDateRange.end ? appliedDateRange.end.toString() : m.documents_any_date();

		return m.documents_date_range_value({ start, end });
	});
	const documentsPageCountLabel = $derived.by(() => {
		const count = documents.length;

		return count === 1
			? m.documents_showing_documents_one({ count })
			: m.documents_showing_documents_other({ count });
	});
	const selectedDocumentsCountLabel = $derived.by(() => {
		const count = selectedIds.size;

		return count === 1
			? m.documents_selected_count_one({ count })
			: m.documents_selected_count_other({ count });
	});

	function resetPagination() {
		cursorState = resetCursorState();
	}

	function setFilterParam(key: string, value: string | undefined) {
		const newParams = new URLSearchParams(page.url.searchParams);
		if (value) {
			newParams.set(key, value);
		} else {
			newParams.delete(key);
		}
		newParams.delete("cursor");
		resetPagination();
		goto(`?${newParams.toString()}`, { keepFocus: true, noScroll: true });
	}

	$effect(() => {
		if (previousCollectionId !== collectionId || previousSchemaId !== schemaId) {
			previousCollectionId = collectionId;
			previousSchemaId = schemaId;
			resetPagination();
			selectedIds = new Set();
			selectedDocumentsById = {};
			renamingDocumentId = null;
		}
	});

	$effect(() => {
		let changed = false;
		const nextSelectedDocumentsById = { ...selectedDocumentsById };

		for (const document of documents) {
			if (!selectedIds.has(document.id)) continue;

			const currentDocument = nextSelectedDocumentsById[document.id];
			if (
				currentDocument &&
				currentDocument.original_filename === document.original_filename &&
				currentDocument.updated_at === document.updated_at &&
				currentDocument.schema_id === document.schema_id &&
				currentDocument.has_inline_schema === document.has_inline_schema
			) continue;

			nextSelectedDocumentsById[document.id] = document;
			changed = true;
		}

		if (changed) {
			selectedDocumentsById = nextSelectedDocumentsById;
		}
	});

	$effect(() => {
		const nextValue = filenameFilter;
		const timeout = setTimeout(() => {
			if (debouncedFilenameFilter !== nextValue) {
				debouncedFilenameFilter = nextValue;
				resetPagination();
			}
		}, 1000);

		return () => clearTimeout(timeout);
	});

	$effect(() => {
		if (renamingDocumentId && !visibleIds.includes(renamingDocumentId)) {
			renamingDocumentId = null;
		}
	});

	function rangeCalendarValue(range?: DateRangeValue): RangeCalendarValue {
		if (!range?.start && !range?.end) return undefined;

		return {
			start: range.start,
			end: range.end,
		};
	}

	function setDatePopoverOpen(open: boolean) {
		datePopoverOpen = open;
		if (open) pendingDateRange = rangeCalendarValue(appliedDateRange);
	}

	const activePreset = $derived.by(() => {
		if (!appliedDateRange?.start || !appliedDateRange?.end) return null;

		const tz = getLocalTimeZone();
		const t = today(tz);

		const todayStr = t.toString();
		if (appliedDateRange.start.toString() === todayStr && appliedDateRange.end.toString() === todayStr) {
			return "today";
		}

		const weekStart = startOfWeek(t, "en-US").toString();
		const weekEnd = endOfWeek(t, "en-US").toString();
		if (appliedDateRange.start.toString() === weekStart && appliedDateRange.end.toString() === weekEnd) {
			return "week";
		}

		const monthStart = startOfMonth(t).toString();
		const monthEnd = endOfMonth(t).toString();
		if (appliedDateRange.start.toString() === monthStart && appliedDateRange.end.toString() === monthEnd) {
			return "month";
		}

		return null;
	});

	function setTodayPreset() {
		const tz = getLocalTimeZone();
		const t = today(tz);
		pendingDateRange = {
			start: t,
			end: t,
		};
		appliedDateRange = pendingDateRange;
		datePopoverOpen = false;
		resetPagination();
	}

	function setThisWeekPreset() {
		const tz = getLocalTimeZone();
		const t = today(tz);
		pendingDateRange = {
			start: startOfWeek(t, "en-US"),
			end: endOfWeek(t, "en-US"),
		};
		appliedDateRange = pendingDateRange;
		datePopoverOpen = false;
		resetPagination();
	}

	function setThisMonthPreset() {
		const tz = getLocalTimeZone();
		const t = today(tz);
		pendingDateRange = {
			start: startOfMonth(t),
			end: endOfMonth(t),
		};
		appliedDateRange = pendingDateRange;
		datePopoverOpen = false;
		resetPagination();
	}

	function applyDateRange() {
		appliedDateRange = pendingDateRange;
		datePopoverOpen = false;
		resetPagination();
	}

	function clearDateRange() {
		pendingDateRange = undefined;
		appliedDateRange = undefined;
		datePopoverOpen = false;
		resetPagination();
	}

	function toggleSort() {
		sortDirection = sortDirection === "desc" ? "asc" : "desc";
		resetPagination();
	}

	function setPageSize(value: string) {
		const nextPageSize = Number(value);
		if (!PAGE_SIZE_OPTIONS.includes(nextPageSize)) return;

		pageSize = nextPageSize;
		resetPagination();
	}

	function toggleAllVisible(checked: boolean) {
		selectedIds = togglePageSelection(visibleIds, selectedIds, checked);

		const nextSelectedDocumentsById = { ...selectedDocumentsById };
		for (const document of documents) {
			if (checked) {
				nextSelectedDocumentsById[document.id] = document;
			} else {
				delete nextSelectedDocumentsById[document.id];
			}
		}
		selectedDocumentsById = nextSelectedDocumentsById;
	}

	function toggleDocument(id: string, checked: boolean) {
		selectedIds = toggleSelection(selectedIds, id, checked);

		const nextSelectedDocumentsById = { ...selectedDocumentsById };
		if (checked) {
			const document = documents.find((document) => document.id === id);
			if (document) nextSelectedDocumentsById[id] = document;
		} else {
			delete nextSelectedDocumentsById[id];
		}
		selectedDocumentsById = nextSelectedDocumentsById;
	}

	function goNext() {
		cursorState = cursorNextState(cursorState, nextCursor);
	}

	function goPrevious() {
		cursorState = cursorPreviousState(cursorState);
	}

	function openPreview(document: OCRDocumentListItemResponse) {
		previewDocumentId = document.id;
		previewFilename = document.original_filename;
		previewOpen = true;
	}

	function openDownloadDialog(documents: OCRDocumentListItemResponse[]) {
		if (documents.length === 0) return;

		downloadTargets = documents;
		downloadError = null;
		downloadDialogOpen = true;
	}

	async function downloadSelectedFormat(format: DownloadFormat) {
		const ids = downloadTargets.map((document) => document.id);
		if (ids.length === 0) return;

		downloadPending = true;
		downloadError = null;

		try {
			const result = await downloadOCRDocuments(fetch, ids, format, m.documents_failed_download());
			triggerBrowserDownload(result.blob, result.filename);
			downloadDialogOpen = false;
		} catch (error) {
			const message = error instanceof Error ? error.message : m.documents_failed_download();
			downloadError = new Error(message);
			toast.error(message);
		} finally {
			downloadPending = false;
		}
	}

	function triggerBrowserDownload(blob: Blob, filename: string) {
		const url = URL.createObjectURL(blob);
		const link = document.createElement("a");
		link.href = url;
		link.download = filename;
		document.body.append(link);
		link.click();
		link.remove();

		setTimeout(() => URL.revokeObjectURL(url), 1000);
	}

	function updateDocumentCaches(updated: UpdateOCRDocumentResponse) {
		queryClient.setQueriesData<OCRDocumentListResponse>(
			{ queryKey: ["ocr-documents"] },
			(current) => {
				if (!current) return current;

				let changed = false;
				const nextDocuments = current.documents.map((document) => {
					if (document.id !== updated.id) return document;

					changed = true;
					return {
						...document,
						...updated,
						collections: document.collections,
					};
				});

				return changed ? { ...current, documents: nextDocuments } : current;
			}
		);
		queryClient.setQueryData<OCRDocumentPreview>(["ocr-document", updated.id], (current) =>
			current ? { ...current, ...updated } : current
		);

		if (selectedDocumentsById[updated.id]) {
			selectedDocumentsById = {
				...selectedDocumentsById,
				[updated.id]: {
					...selectedDocumentsById[updated.id],
					...updated,
				},
			};
		}

		if (previewDocumentId === updated.id) {
			previewFilename = updated.original_filename;
		}
	}

	async function updateDocumentFilename(
		id: string,
		currentFilename: string | null,
		originalFilename: string
	) {
		const nextFilename = originalFilename.trim();
		if (!nextFilename || nextFilename === currentFilename) return;

		updateMutation.reset();
		updatingDocumentId = id;
		try {
			await updateMutation.mutateAsync({ id, originalFilename: nextFilename });
		} finally {
			updatingDocumentId = null;
		}
	}

	async function renameDocument(document: OCRDocumentListItemResponse, originalFilename: string) {
		await updateDocumentFilename(document.id, document.original_filename, originalFilename);
	}

	async function renamePreviewDocument(originalFilename: string) {
		if (!previewDocumentId) return;

		await updateDocumentFilename(previewDocumentId, previewFilename, originalFilename);
	}

	function clearDeletedDocuments(deletedIds: string[], submittedIds: string[]) {
		const submitted = new Set(submittedIds);
		const deleted = new Set(deletedIds);

		selectedIds = new Set([...selectedIds].filter((id) => !submitted.has(id)));
		selectedDocumentsById = Object.fromEntries(
			Object.entries(selectedDocumentsById).filter(([id]) => !submitted.has(id))
		);
		for (const id of deleted) {
			queryClient.removeQueries({ queryKey: ["ocr-document", id] });
		}

		if (previewDocumentId && deleted.has(previewDocumentId)) {
			previewOpen = false;
			previewDocumentId = null;
			previewFilename = null;
		}

		if (renamingDocumentId && submitted.has(renamingDocumentId)) {
			renamingDocumentId = null;
		}
	}

	async function runDelete(ids: string[]) {
		try {
			await deleteMutation.mutateAsync({ ids });
		} catch {
			// The mutation owns the error state; the page renders it below the toolbar.
		}
	}

	function confirmSingleDelete(document: OCRDocumentListItemResponse) {
		deleteMutation.reset();
		confirmDelete({
			title: m.documents_delete_single_title(),
			description: m.documents_delete_single_description({ name: document.original_filename }),
			confirm: { text: m.documents_delete() },
			onConfirm: () => runDelete([document.id]),
		});
	}

	function confirmBulkDelete() {
		const ids = [...selectedIds];
		if (ids.length === 0) return;

		deleteMutation.reset();
		confirmDelete({
			title:
				ids.length === 1
					? m.documents_delete_bulk_title_one({ count: ids.length })
					: m.documents_delete_bulk_title_other({ count: ids.length }),
			description:
				ids.length === 1
					? m.documents_delete_bulk_description_one({ count: ids.length })
					: m.documents_delete_bulk_description_other({ count: ids.length }),
			confirm: { text: m.documents_delete() },
			onConfirm: () => runDelete(ids),
		});
	}

	function openMoveDialog() {
		moveMutation.reset();
		moveDialogOpen = true;
	}

	function moveSelectedDocuments(collectionIds: string[]) {
		const ids = [...selectedIds];
		if (ids.length === 0) return;

		moveMutation.mutate({ ids, collectionIds });
	}

	const columns: ColumnDef<OCRDocumentListItemResponse>[] = [
		{
			id: "select",
			header: () =>
				renderComponent(SelectionCheckbox, {
					checked: pageSelection.checked,
					indeterminate: pageSelection.indeterminate,
					ariaLabel: m.documents_select_all_on_page(),
					onCheckedChange: toggleAllVisible,
				}),
			cell: ({ row }) =>
				renderComponent(SelectionCheckbox, {
					checked: selectedIds.has(row.original.id),
					ariaLabel: m.documents_select_document({ name: row.original.original_filename }),
					onCheckedChange: (checked) => toggleDocument(row.original.id, checked),
				}),
			enableSorting: false,
			enableHiding: false,
		},
		{
			accessorKey: "original_filename",
			header: m.documents_filename_column(),
			cell: ({ row }) =>
				renderComponent(FilenameCell, {
					document: row.original,
					schemaName: row.original.schema_id ? schemasMap.get(row.original.schema_id)?.name : undefined,
					editing: () => renamingDocumentId === row.original.id,
					onEditingChange: (editing) => {
						renamingDocumentId = editing ? row.original.id : null;
					},
					onPreview: openPreview,
					onRename: (originalFilename) => renameDocument(row.original, originalFilename),
					renamePending: updatingDocumentId === row.original.id,
				}),
		},
		{
			id: "collections",
			header: m.documents_collections_column(),
			cell: ({ row }) => renderComponent(CollectionsCell, { document: row.original }),
			enableSorting: false,
		},
		{
			accessorKey: "page_count",
			header: m.documents_pages_column(),
		},
		{
			accessorKey: "created_at",
			header: () =>
				renderComponent(CreatedDateHeader, {
					sortDirection,
					onToggle: toggleSort,
				}),
			cell: ({ row }) => formatCreatedDate(row.original.created_at, m.documents_invalid_date()),
		},
		{
			accessorKey: "file_size",
			header: m.documents_file_size_column(),
			cell: ({ row }) => formatFileSize(row.original.file_size),
		},
		{
			id: "actions",
			header: () => null,
			cell: ({ row }) =>
				renderComponent(ActionsCell, {
					document: row.original,
					onPreview: openPreview,
					onRename: (document) => {
						renamingDocumentId = document.id;
					},
					onDownload: (document) => openDownloadDialog([document]),
					onDelete: confirmSingleDelete,
					renamePending: updatingDocumentId === row.original.id,
					downloadPending,
					deletePending: deleteMutation.isPending,
				}),
			enableSorting: false,
			enableHiding: false,
		},
	];

	const table = createSvelteTable({
		get data() {
			return documents;
		},
		columns,
		getRowId: (row) => row.id,
		getCoreRowModel: getCoreRowModel(),
	});
</script>

<svelte:head>
	<title>{m.documents_page_title()} | Syncra</title>
</svelte:head>

<div class="@container/main flex flex-1 flex-col gap-4 p-4 lg:p-6">
	<!-- Page Header -->
	<div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-end">
		
		<Button href="/app/new-job" size="sm" class="h-10 px-5 shadow-sm cursor-pointer">
			<PlusIcon class="size-4" aria-hidden="true" />
			{m.documents_new_ocr_job()}
		</Button>
	</div>
	<!-- Search & Filters Toolbar -->
	<div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
		<div class="flex min-w-0 flex-1 flex-wrap items-center gap-2">
			<div class="relative min-w-[240px] flex-1 max-w-sm">
				<SearchIcon
					class="text-muted-foreground pointer-events-none absolute top-1/2 left-2.5 size-4 -translate-y-1/2"
					aria-hidden="true"
				/>
				<Input
					value={filenameFilter}
					placeholder={m.documents_search_filename_placeholder()}
					class="pl-8 bg-background/50 focus-visible:ring-1"
					aria-label={m.documents_search_filename()}
					oninput={(event) => {
						filenameFilter = event.currentTarget.value;
					}}
				/>
			</div>

			<Popover.Root bind:open={() => datePopoverOpen, setDatePopoverOpen}>
				<Popover.Trigger>
					{#snippet child({ props })}
						<Button type="button" variant="outline" class="w-full justify-start sm:w-auto bg-background/50 min-w-[160px] h-9 text-xs" {...props}>
							<CalendarIcon class="size-4 mr-2 text-muted-foreground" aria-hidden="true" />
							<span class="truncate">{dateRangeLabel}</span>
						</Button>
					{/snippet}
				</Popover.Trigger>
				<Popover.Content align="start" class="w-auto p-0">
					<div class="flex flex-col sm:flex-row">
						<div class="flex flex-row sm:flex-col gap-1.5 p-3 border-b sm:border-b-0 sm:border-r border-border min-w-[130px] bg-muted/5">
							<span class="hidden sm:inline-block text-[10px] font-bold uppercase tracking-wider text-muted-foreground/60 px-2 py-1 select-none">{m.documents_presets()}</span>
							<Button
								type="button"
								variant={activePreset === "today" ? "secondary" : "ghost"}
								class="flex-1 sm:flex-none justify-center sm:justify-start text-xs font-medium h-8.5 px-3 rounded-md transition-all duration-150 active:scale-[0.98]"
								onclick={setTodayPreset}
							>
								{m.documents_today()}
							</Button>
							<Button
								type="button"
								variant={activePreset === "week" ? "secondary" : "ghost"}
								class="flex-1 sm:flex-none justify-center sm:justify-start text-xs font-medium h-8.5 px-3 rounded-md transition-all duration-150 active:scale-[0.98]"
								onclick={setThisWeekPreset}
							>
								{m.documents_this_week()}
							</Button>
							<Button
								type="button"
								variant={activePreset === "month" ? "secondary" : "ghost"}
								class="flex-1 sm:flex-none justify-center sm:justify-start text-xs font-medium h-8.5 px-3 rounded-md transition-all duration-150 active:scale-[0.98]"
								onclick={setThisMonthPreset}
							>
								{m.documents_this_month()}
							</Button>
						</div>
						<div class="flex flex-col">
							<RangeCalendar.RangeCalendar bind:value={pendingDateRange} numberOfMonths={2} />
							<div class="flex justify-end gap-2 border-t p-3 bg-muted/20">
								<Button type="button" variant="ghost" size="sm" onclick={clearDateRange}>{m.documents_clear()}</Button>
								<Button type="button" size="sm" onclick={applyDateRange}>{m.documents_apply()}</Button>
							</div>
						</div>
					</div>
				</Popover.Content>
			</Popover.Root>

			<!-- Filter by Collection Select -->
			<Select.Root type="single" value={collectionId ?? "all"} onValueChange={(val) => setFilterParam("collection", val === "all" ? undefined : val)}>
				<Select.Trigger class="w-48 bg-background/50 h-9 text-xs justify-between" aria-label={m.documents_filter_by_collection()}>
					<div class="flex items-center min-w-0">
						<FolderIcon class="size-3.5 mr-2 text-muted-foreground shrink-0" aria-hidden="true" />
						<span class="truncate text-left">
							{selectedCollectionName}
						</span>
					</div>
				</Select.Trigger>
				<Select.Content side="bottom" class="max-h-60 overflow-y-auto">
					<Select.Item value="all" class="text-xs">{m.documents_all_collections()}</Select.Item>
					{#each collections as col (col.id)}
						<Select.Item value={col.id} class="text-xs">{col.name}</Select.Item>
					{/each}
				</Select.Content>
			</Select.Root>

			<!-- Filter by Schema Select -->
			<Select.Root type="single" value={schemaId ?? "all"} onValueChange={(val) => setFilterParam("schema", val === "all" ? undefined : val)}>
				<Select.Trigger class="w-48 bg-background/50 h-9 text-xs justify-between" aria-label={m.documents_filter_by_schema()}>
					<div class="flex items-center min-w-0">
						<svg xmlns="http://www.w3.org/2000/svg" class="size-3.5 mr-2 text-muted-foreground shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><rect width="18" height="18" x="3" y="3" rx="2"/><path d="M3 9h18"/><path d="M9 21V9"/></svg>
						<span class="truncate text-left">
							{selectedSchemaName}
						</span>
					</div>
				</Select.Trigger>
				<Select.Content side="bottom" class="max-h-60 overflow-y-auto">
					<Select.Item value="all" class="text-xs">{m.documents_all_schemas()}</Select.Item>
					{#each schemas as schema (schema.id)}
						<Select.Item value={schema.id} class="text-xs">{schema.name}</Select.Item>
					{/each}
				</Select.Content>
			</Select.Root>
		</div>
	</div>

	{#if deleteMutation.isError}
		<Alert.Root variant="destructive" class="mt-2">
			<Alert.Description>{deleteMutation.error.message}</Alert.Description>
		</Alert.Root>
	{/if}

	{#if moveMutation.isError}
		<Alert.Root variant="destructive" class="mt-2">
			<Alert.Description>{moveMutation.error.message}</Alert.Description>
		</Alert.Root>
	{/if}

	{#if updateMutation.isError}
		<Alert.Root variant="destructive" class="mt-2">
			<Alert.Description>{updateMutation.error.message}</Alert.Description>
		</Alert.Root>
	{/if}

	<!-- Data Table Container -->
	<div class="overflow-x-auto rounded-xl border bg-background shadow-xs">
		<Table.Root>
			<Table.Header class="bg-muted/40 sticky top-0 z-10 border-b">
				{#each table.getHeaderGroups() as headerGroup (headerGroup.id)}
					<Table.Row class="hover:bg-transparent">
						{#each headerGroup.headers as header (header.id)}
							<Table.Head colspan={header.colSpan} class="h-10 text-xs font-semibold uppercase tracking-wider text-muted-foreground/90 py-2.5">
								{#if !header.isPlaceholder}
									<FlexRender
										content={header.column.columnDef.header}
										context={header.getContext()}
									/>
								{/if}
							</Table.Head>
						{/each}
					</Table.Row>
				{/each}
			</Table.Header>
			<Table.Body>
				{#if documentsQuery.isLoading}
					<Table.Row>
						<Table.Cell colspan={columns.length} class="h-56 p-0">
							<div class="flex h-56 w-full items-center justify-center">
								<Spinner class="size-18 text-foreground dark:text-blue-500" />
							</div>
						</Table.Cell>
					</Table.Row>
				{:else if collectionNotFound}
					<Table.Row>
						<Table.Cell colspan={columns.length} class="h-56 text-center">
							<div class="mx-auto flex max-w-[340px] flex-col items-center justify-center gap-2 p-6">
								<div class="rounded-full bg-muted/50 p-3.5 text-muted-foreground/80">
									<FolderIcon class="size-6" aria-hidden="true" />
								</div>
								<h3 class="mt-2 text-sm font-semibold text-foreground">
									{m.documents_collection_not_found_title()}
								</h3>
								<p class="px-2 text-xs leading-normal text-muted-foreground">
									{m.documents_collection_not_found_body()}
								</p>
								<Button
									type="button"
									variant="outline"
									size="sm"
									class="mt-3.5 h-8 text-xs font-medium"
									onclick={() => setFilterParam("collection", undefined)}
								>
									{m.documents_view_all_documents()}
								</Button>
							</div>
						</Table.Cell>
					</Table.Row>
				{:else if documentsQuery.isError}
					<Table.Row>
						<Table.Cell colspan={columns.length} class="h-24">
							<div
								class="mx-auto flex max-w-xl flex-wrap items-center justify-center gap-3 text-center text-sm text-destructive"
							>
								<span>{documentsQuery.error.message}</span>
								<Button
									type="button"
									variant="outline"
									size="sm"
									onclick={() => documentsQuery.refetch()}
								>
									{m.documents_retry()}
								</Button>
							</div>
						</Table.Cell>
					</Table.Row>
				{:else if table.getRowModel().rows.length}
					{#each table.getRowModel().rows as row (row.id)}
						<Table.Row 
							data-state={selectedIds.has(row.original.id) && "selected"}
							class="transition-colors hover:bg-muted/40 duration-150"
						>
							{#each row.getVisibleCells() as cell (cell.id)}
								<Table.Cell class="py-3 px-4">
									<FlexRender content={cell.column.columnDef.cell} context={cell.getContext()} />
								</Table.Cell>
							{/each}
						</Table.Row>
					{/each}
				{:else if documentsQuery.isSuccess && documents.length === 0}
					<Table.Row>
						<Table.Cell colspan={columns.length} class="h-56 text-center">
							<div class="mx-auto flex max-w-[340px] flex-col items-center justify-center gap-2 p-6">
								<div class="rounded-full bg-muted/50 p-3.5 text-muted-foreground/80">
									<SearchIcon class="size-6" aria-hidden="true" />
								</div>
								<h3 class="text-sm font-semibold mt-2 text-foreground">
									{m.documents_no_documents_found()}
								</h3>
								<p class="text-xs text-muted-foreground leading-normal px-2">
									{#if filenameFilter || appliedDateRange || collectionId || schemaId}
										{m.documents_empty_filtered_body()}
									{:else}
										{m.documents_empty_unfiltered_body()}
									{/if}
								</p>
								{#if filenameFilter || appliedDateRange || collectionId || schemaId}
									<Button
										type="button"
										variant="outline"
										size="sm"
										class="mt-3.5 h-8 text-xs font-medium"
										onclick={() => {
											filenameFilter = "";
											clearDateRange();
											const newParams = new URLSearchParams(page.url.searchParams);
											newParams.delete("collection");
											newParams.delete("schema");
											newParams.delete("cursor");
											resetPagination();
											goto(`?${newParams.toString()}`, { keepFocus: true, noScroll: true });
										}}
									>
										{m.documents_clear_filters()}
									</Button>
								{:else}
									<Button href="/app/new-job" size="sm" class="mt-3.5 h-8 text-xs font-medium gap-1.5">
										<UploadIcon class="size-3.5" />
										{m.documents_process_first_document()}
									</Button>
								{/if}
							</div>
						</Table.Cell>
					</Table.Row>
				{/if}
			</Table.Body>
		</Table.Root>
	</div>

	<!-- Pagination Footer -->
	<div
		class="flex flex-col gap-3 text-xs text-muted-foreground sm:flex-row sm:items-center sm:justify-between px-1.5"
	>
		<div>
			{#if documents.length > 0}
				{documentsPageCountLabel}
			{:else}
				{m.documents_no_documents_to_show()}
			{/if}
		</div>
		<div class="flex flex-wrap items-center gap-3">
			<Select.Root type="single" bind:value={() => String(pageSize), setPageSize}>
				<Select.Trigger size="sm" class="w-24 bg-background/50 h-8 text-xs" aria-label={m.documents_rows_per_page()}>
					{pageSize}
				</Select.Trigger>
				<Select.Content side="top">
					{#each PAGE_SIZE_OPTIONS as option (option)}
						<Select.Item value={String(option)} class="text-xs">{option}</Select.Item>
					{/each}
				</Select.Content>
			</Select.Root>

			<div class="flex items-center gap-2">
				<Button
					type="button"
					variant="outline"
					size="sm"
					class="h-8 text-xs"
					onclick={goPrevious}
					disabled={cursorState.history.length === 0 || documentsQuery.isFetching}
				>
					<ChevronLeftIcon class="size-4 mr-1" aria-hidden="true" />
					{m.documents_previous()}
				</Button>
				<Button
					type="button"
					variant="outline"
					size="sm"
					class="h-8 text-xs"
					onclick={goNext}
					disabled={!nextCursor || documentsQuery.isFetching}
				>
					{m.documents_next()}
					<ChevronRightIcon class="size-4 ml-1" aria-hidden="true" />
				</Button>
			</div>
		</div>
	</div>

	<DocumentPreviewDialog
		bind:open={previewOpen}
		filename={previewFilename}
		markdown={previewQuery.data?.markdown}
		annotationJson={previewQuery.data?.annotation_json}
		isLoading={previewQuery.isLoading}
		error={previewQuery.error}
		onRetry={() => previewQuery.refetch()}
		onRename={renamePreviewDocument}
		renamePending={previewDocumentId !== null && updatingDocumentId === previewDocumentId}
	/>

	<MoveCollectionsDialog
		bind:open={moveDialogOpen}
		collections={collections}
		collectionsLoading={collectionsQuery.isLoading}
		collectionsError={collectionsQuery.error ?? null}
		selectedCount={selectedIds.size}
		pending={moveMutation.isPending}
		error={moveMutation.error ?? null}
		onSubmit={moveSelectedDocuments}
	/>

	<DownloadDialog
		bind:open={downloadDialogOpen}
		documents={downloadTargets}
		pending={downloadPending}
		error={downloadError}
		onDownload={downloadSelectedFormat}
	/>

	<!-- Floating Bulk Actions Center Bar -->
	{#if selectedIds.size > 0}
		<div class="fixed bottom-6 left-1/2 z-50 -translate-x-1/2 flex items-center gap-3.5 px-4.5 py-2.5 rounded-full border border-border bg-background/95 shadow-xl backdrop-blur-md animate-in fade-in slide-in-from-bottom-4 duration-300">
			<span class="text-xs font-semibold text-foreground px-1 select-none whitespace-nowrap">
				{selectedDocumentsCountLabel}
			</span>
			<div class="h-4 w-[1px] bg-border"></div>
			<div class="flex items-center gap-1.5">
				<Button
					type="button"
					variant="ghost"
					size="sm"
					class="h-8 rounded-full px-3 text-xs text-muted-foreground/80 hover:text-foreground transition-all"
					disabled={downloadPending || selectedDocuments.length === 0}
					aria-label={m.documents_download_selected()}
					onclick={() => openDownloadDialog(selectedDocuments)}
				>
					{#if downloadPending}
						<LoaderIcon class="size-3.5 animate-spin mr-1.5" aria-hidden="true" />
						{m.documents_downloading()}
					{:else}
						<DownloadIcon class="size-3.5 mr-1.5" aria-hidden="true" />
						{m.documents_download()}
					{/if}
				</Button>

				<Button
					type="button"
					variant="ghost"
					size="sm"
					class="h-8 rounded-full px-3 text-xs text-muted-foreground/80 hover:text-foreground transition-all"
					disabled={moveMutation.isPending}
					onclick={openMoveDialog}
				>
					{#if moveMutation.isPending}
						<LoaderIcon class="size-3.5 animate-spin mr-1.5" aria-hidden="true" />
						{m.documents_moving()}
					{:else}
						<FolderIcon class="size-3.5 mr-1.5" aria-hidden="true" />
						{m.documents_move()}
					{/if}
				</Button>

				<Button
					type="button"
					variant="destructive"
					size="sm"
					class="h-8 rounded-full px-4.5 text-xs font-medium shadow-xs"
					disabled={deleteMutation.isPending}
					onclick={confirmBulkDelete}
				>
					{#if deleteMutation.isPending}
						<LoaderIcon class="size-3.5 animate-spin mr-1.5" aria-hidden="true" />
						{m.documents_deleting()}
					{:else}
						<Trash2Icon class="size-3.5 mr-1.5" aria-hidden="true" />
						{m.documents_delete()}
					{/if}
				</Button>
			</div>
		</div>
	{/if}
</div>
