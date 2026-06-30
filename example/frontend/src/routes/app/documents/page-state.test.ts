import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";

const source = () => readFileSync(new URL("./+page.svelte", import.meta.url), "utf8");

function normalizeSource(value: string) {
	return value.replace(/\s+/g, " ").trim();
}

function extractBalanced(sourceText: string, startIndex: number, open: string, close: string) {
	let depth = 0;

	for (let index = startIndex; index < sourceText.length; index += 1) {
		const character = sourceText[index];
		if (character === open) depth += 1;
		if (character === close) {
			depth -= 1;
			if (depth === 0) return sourceText.slice(startIndex, index + 1);
		}
	}

	throw new Error(`Unable to find matching "${close}" for "${open}"`);
}

function extractBalancedAfter(sourceText: string, marker: string, open: string, close: string) {
	const markerIndex = sourceText.indexOf(marker);
	if (markerIndex === -1) throw new Error(`Missing source marker: ${marker}`);

	const startIndex = sourceText.indexOf(open, markerIndex + marker.length);
	if (startIndex === -1) throw new Error(`Missing "${open}" after source marker: ${marker}`);

	return extractBalanced(sourceText, startIndex, open, close);
}

function extractDocumentsQuery(sourceText: string) {
	const startIndex = sourceText.indexOf("const documentsQuery = createQuery");
	const endIndex = sourceText.indexOf("const previewQuery =", startIndex);

	if (startIndex === -1) throw new Error("Missing documents query");
	if (endIndex === -1) throw new Error("Missing documents query boundary");

	return sourceText.slice(startIndex, endIndex);
}

function extractEffectContaining(sourceText: string, marker: string) {
	const markerIndex = sourceText.indexOf(marker);
	if (markerIndex === -1) throw new Error(`Missing effect marker: ${marker}`);

	const effectIndex = sourceText.lastIndexOf("$effect(() =>", markerIndex);
	if (effectIndex === -1) throw new Error(`Missing $effect for marker: ${marker}`);

	const startIndex = sourceText.indexOf("{", effectIndex);
	if (startIndex === -1) throw new Error(`Missing effect body for marker: ${marker}`);

	return extractBalanced(sourceText, startIndex, "{", "}");
}

function expectDocumentsFilterObject(objectSource: string) {
	const object = normalizeSource(objectSource);

	expect(object).toMatch(/^\{ collectionId,/);
	expect(object).toContain("filename: normalizedFilename");
	expect(object).toContain("createdFrom: dateBounds.createdFrom");
	expect(object).toContain("createdTo: dateBounds.createdTo");
	expect(object).toContain("cursor: effectiveCursor");
	expect(object).toContain("size: pageSize");
	expect(object).toContain("sort: sortDirection");
}

describe("documents page collection filter wiring", () => {
	it("derives collection state from the page URL", () => {
		const page = source();
		expect(page).toContain('import { page } from "$app/state";');
		expect(page).toContain("let previousCollectionId = $state<string | undefined>();");
		expect(normalizeSource(page)).toContain(
			'const collectionId = $derived(page.url.searchParams.get("collection")?.trim() || undefined);'
		);
		expect(normalizeSource(page)).toContain(
			"const effectiveCursor = $derived( previousCollectionId === collectionId && previousSchemaId === schemaId ? cursorState.currentCursor : null );"
		);
	});

	it("uses the effective cursor in the documents query key and request variables", () => {
		const page = source();
		const documentsQuery = extractDocumentsQuery(page);
		const queryKey = extractBalancedAfter(documentsQuery, "queryKey:", "[", "]");
		const queryKeyFilters = extractBalancedAfter(queryKey, '"ocr-documents"', "{", "}");
		const requestVariables = extractBalancedAfter(
			documentsQuery,
			"fetchOCRDocuments(",
			"{",
			"}"
		);

		expectDocumentsFilterObject(queryKeyFilters);
		expectDocumentsFilterObject(requestVariables);
		expect(normalizeSource(documentsQuery)).not.toContain("cursor: cursorState.currentCursor");
	});

	it("treats missing collection filters as terminal not-found errors", () => {
		const page = source();
		const normalized = normalizeSource(page);

		expect(page).toContain("isOCRDocumentNotFoundError");
		expect(page).toContain("shouldRetryOCRDocumentsQuery");
		expect(normalized).toContain(
			"retry: (failureCount, error) => shouldRetryOCRDocumentsQuery(failureCount, error, { collectionId }),"
		);
		expect(normalized).toContain(
			"const collectionNotFound = $derived( Boolean(collectionId) && isOCRDocumentNotFoundError(documentsQuery.error) );"
		);
		expect(normalized).toContain("{:else if collectionNotFound}");
		expect(page).toContain('import { m } from "$lib/paraglide/messages.js";');
		expect(page).toContain("m.documents_collection_not_found_title()");
		expect(page).toContain("m.documents_collection_not_found_body()");
		expect(page).toContain("m.documents_view_all_documents()");
		expect(page).toContain('onclick={() => setFilterParam("collection", undefined)}');
		expect(page).toContain("selectedCollectionName");
		expect(page).toContain("m.documents_unknown_collection()");
		expect(page).toContain("m.documents_all_collections()");
	});

	it("resets pagination when collection changes", () => {
		const page = source();
		const effect = normalizeSource(
			extractEffectContaining(page, "if (previousCollectionId !== collectionId || previousSchemaId !== schemaId)")
		);

		expect(effect).toContain(
			"if (previousCollectionId !== collectionId || previousSchemaId !== schemaId) { previousCollectionId = collectionId; previousSchemaId = schemaId; resetPagination(); selectedIds = new Set(); selectedDocumentsById = {}; renamingDocumentId = null; }"
		);
	});

	it("adds a collections column to the documents table", () => {
		const page = source();
		const normalized = normalizeSource(page);

		expect(page).toContain('import CollectionsCell from "./collections-cell.svelte";');
		expect(normalized).toContain(
			"id: \"collections\", header: m.documents_collections_column(), cell: ({ row }) => renderComponent(CollectionsCell, { document: row.original }), enableSorting: false,"
		);
	});

	it("wires the bulk move dialog to collections and OCR document queries", () => {
		const page = source();
		const normalized = normalizeSource(page);

		expect(page).toContain('import MoveCollectionsDialog from "./move-collections-dialog.svelte";');
		expect(page).toContain("fetchCollections");
		expect(page).toContain("moveOCRDocuments");
		expect(normalized).toContain('queryKey: ["collections"]');
		expect(normalized).toContain('mutationKey: ["ocr-documents", "move"]');
		expect(normalized).toContain('void queryClient.invalidateQueries({ queryKey: ["ocr-documents"] });');
		expect(normalized).toContain('void queryClient.invalidateQueries({ queryKey: ["collections"] });');
		expect(normalized).toContain("function openMoveDialog()");
		expect(normalized).toContain("function moveSelectedDocuments(collectionIds: string[])");
		expect(normalized).toContain("onSubmit={moveSelectedDocuments}");
		expect(normalized).toContain("m.documents_move()");
	});

	it("uses shared personal schema options for schema filters", () => {
		const page = source();
		const normalized = normalizeSource(page);

		expect(page).toContain("PERSONAL_SCHEMA_OPTIONS_QUERY_KEY");
		expect(page).toContain("fetchPersonalSchemaOptions");
		expect(normalized).toContain("queryKey: PERSONAL_SCHEMA_OPTIONS_QUERY_KEY");
		expect(normalized).toContain("queryFn: () => fetchPersonalSchemaOptions(fetch)");
		expect(normalized).toContain("const schemas = $derived(schemasQuery.data ?? []);");
	});

	it("wires document download state, modal, and API calls", () => {
		const page = source();
		const normalized = normalizeSource(page);

		expect(page).toContain('import DownloadDialog from "./download-dialog.svelte";');
		expect(page).toContain("downloadOCRDocuments");
		expect(page).toContain("type DownloadFormat");
		expect(page).toContain('import { toast } from "svelte-sonner";');
		expect(normalized).toContain("let selectedDocumentsById = $state.raw<SelectedDocumentsById>({});");
		expect(normalized).toContain("let downloadDialogOpen = $state(false);");
		expect(normalized).toContain("let downloadPending = $state(false);");
		expect(normalized).toContain("function openDownloadDialog(documents: OCRDocumentListItemResponse[])");
		expect(normalized).toContain("async function downloadSelectedFormat(format: DownloadFormat)");
		expect(normalized).toContain(
			"const result = await downloadOCRDocuments(fetch, ids, format, m.documents_failed_download());"
		);
		expect(normalized).toContain("m.documents_failed_download()");
		expect(normalized).toContain("toast.error(message);");
		expect(normalized).toContain("<DownloadDialog bind:open={downloadDialogOpen} documents={downloadTargets} pending={downloadPending} error={downloadError} onDownload={downloadSelectedFormat} />");
	});

	it("opens downloads from row and selected bulk actions", () => {
		const page = source();
		const normalized = normalizeSource(page);

		expect(normalized).toContain("onDownload: (document) => openDownloadDialog([document])");
		expect(normalized).toContain("downloadPending, deletePending: deleteMutation.isPending");
		expect(normalized).toContain("onclick={() => openDownloadDialog(selectedDocuments)}");
		expect(normalized).toContain("disabled={downloadPending || selectedDocuments.length === 0}");
		expect(normalized).toContain("m.documents_downloading()");
	});

	it("stores selected document metadata across page selections", () => {
		const page = source();
		const normalized = normalizeSource(page);

		expect(normalized).toContain(
			"const selectedDocuments = $derived( [...selectedIds] .map((id) => selectedDocumentsById[id]) .filter((document): document is OCRDocumentListItemResponse => Boolean(document)) );"
		);
		expect(normalized).toContain("nextSelectedDocumentsById[document.id] = document;");
		expect(normalized).not.toContain("selectedDocumentsById[document.id] === document");
		expect(normalized).not.toContain("currentDocument === document");
		expect(normalized).toContain("delete nextSelectedDocumentsById[document.id];");
		expect(normalized).toContain("delete nextSelectedDocumentsById[id];");
		expect(normalized).toContain("selectedDocumentsById = {};");
	});

	it("wires document rename activation through actions and preview dialog", () => {
		const page = source();
		const normalized = normalizeSource(page);

		expect(normalized).toContain("let renamingDocumentId = $state<string | null>(null);");
		expect(normalized).toContain("function renamePreviewDocument(originalFilename: string)");
		expect(normalized).toContain("editing: () => renamingDocumentId === row.original.id");
		expect(normalized).toContain("onEditingChange: (editing) => { renamingDocumentId = editing ? row.original.id : null; }");
		expect(normalized).toContain("onPreview: openPreview");
		expect(normalized).toContain("onRename: (document) => { renamingDocumentId = document.id; }");
		expect(normalized).toContain("onRename={renamePreviewDocument}");
		expect(normalized).toContain("renamePending={previewDocumentId !== null && updatingDocumentId === previewDocumentId}");
	});
});
