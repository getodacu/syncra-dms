import { readFileSync } from "node:fs";
import { describe, expect, it, vi } from "vitest";

import type { CollectionInput, CollectionResponse } from "$lib/client/collections";
import {
	COLLECTION_LOADING_ROWS,
	COLLECTION_OVERFLOW_THRESHOLD,
	COLLECTIONS_QUERY_KEY,
	OCR_DOCUMENTS_QUERY_KEY,
	closeCollectionDialogState,
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
	isAllDocumentsActive,
	isCollectionActive,
	openCreateCollectionDialogState,
	openEditCollectionDialogState,
	retryCollections,
	runCollectionMenuAction,
	stopCollectionMenuNavigation,
	type CollectionQueryKey,
} from "./nav-collections-utils";

const collection = (patch: Partial<CollectionResponse> = {}): CollectionResponse => ({
	id: "collection-1",
	created_at: "2026-05-01T00:00:00.000Z",
	updated_at: "2026-05-02T00:00:00.000Z",
	user_id: "user-1",
	name: "Invoices",
	schema_ids: ["schema-1"],
	schema_count: 1,
	document_count: 2,
	...patch,
});

const input: CollectionInput = {
	name: "Receipts",
	schema_ids: ["schema-2", "schema-3"],
};

const queryKeys = (keys: CollectionQueryKey[]) => keys.map((key) => [...key]);
const navSource = () =>
	readFileSync(new URL("./nav-collections.svelte", import.meta.url), "utf8");
const collectionDialogSource = () =>
	readFileSync(new URL("./collection-dialog.svelte", import.meta.url), "utf8");
const normalizeSource = (value: string) => value.replace(/\s+/g, " ").trim();

describe("collections sidebar behavior", () => {
	it("derives active document and collection links from route state", () => {
		expect(isAllDocumentsActive("/app/documents", null)).toBe(true);
		expect(isAllDocumentsActive("/app/documents", "")).toBe(true);
		expect(isAllDocumentsActive("/app/documents", "collection-1")).toBe(false);
		expect(isAllDocumentsActive("/app/schemas", null)).toBe(false);

		expect(isCollectionActive("/app/documents", "collection-1", "collection-1")).toBe(true);
		expect(isCollectionActive("/app/documents", "collection-2", "collection-1")).toBe(false);
		expect(isCollectionActive("/app/schemas", "collection-1", "collection-1")).toBe(false);
		expect(collectionHref("collection with space")).toBe(
			"/app/documents?collection=collection+with+space"
		);
	});

	it("models the add/create dialog flow", () => {
		expect(openCreateCollectionDialogState()).toEqual({
			open: true,
			mode: "create",
			editingCollection: null,
		});
		expect(collectionDialogInitialValue(null)).toBeUndefined();
		expect(collectionSubmitAction("create", null, input)).toEqual({ type: "create", input });
		expect(closeCollectionDialogState("create")).toEqual({
			open: false,
			mode: "create",
			editingCollection: null,
		});
	});

	it("models edit dialog state and update submit actions", () => {
		const item = collection({ id: "collection-7", name: "Contracts", schema_ids: ["schema-a"] });

		expect(openEditCollectionDialogState(item)).toEqual({
			open: true,
			mode: "edit",
			editingCollection: item,
		});
		expect(collectionDialogInitialValue(item)).toEqual({
			name: "Contracts",
			schema_ids: ["schema-a"],
		});
		expect(collectionSubmitAction("edit", item, input)).toEqual({
			type: "update",
			id: "collection-7",
			input,
		});
		expect(collectionSubmitAction("edit", null, input)).toEqual({ type: "none" });
	});

	it("selects dialog pending and error state by mutation mode", () => {
		const createError = new Error("create failed");
		const updateError = new Error("update failed");

		expect(collectionDialogPending(false, false)).toBe(false);
		expect(collectionDialogPending(true, false)).toBe(true);
		expect(collectionDialogPending(false, true)).toBe(true);
		expect(collectionDialogError("create", createError, updateError)).toBe(createError);
		expect(collectionDialogError("edit", createError, updateError)).toBe(updateError);
	});

	it("models loading, retry, error, empty, and ready list states", async () => {
		const refetch = vi.fn(() => Promise.resolve("refetched"));

		expect(COLLECTION_LOADING_ROWS).toEqual([0, 1, 2]);
		expect(
			collectionListStatus({ isLoading: true, isError: true, collectionCount: 1 })
		).toBe("loading");
		expect(
			collectionListStatus({ isLoading: false, isError: true, collectionCount: 1 })
		).toBe("error");
		expect(
			collectionListStatus({ isLoading: false, isError: false, collectionCount: 0 })
		).toBe("empty");
		expect(
			collectionListStatus({ isLoading: false, isError: false, collectionCount: 1 })
		).toBe("ready");
		await expect(retryCollections(refetch)).resolves.toBe("refetched");
		expect(refetch).toHaveBeenCalledOnce();
	});

	it("only overflows collection navigation after ten collections", () => {
		expect(COLLECTION_OVERFLOW_THRESHOLD).toBe(10);
		expect(collectionListOverflows(0)).toBe(false);
		expect(collectionListOverflows(10)).toBe(false);
		expect(collectionListOverflows(11)).toBe(true);
	});

	it("invalidates documents after updating the selected collection only", () => {
		expect(
			queryKeys(
				collectionUpdateSuccessInvalidationKeys({
					pathname: "/app/documents",
					selectedCollectionId: "collection-1",
					updatedCollectionId: "collection-1",
				})
			)
		).toEqual([COLLECTIONS_QUERY_KEY, OCR_DOCUMENTS_QUERY_KEY]);

		expect(
			queryKeys(
				collectionUpdateSuccessInvalidationKeys({
					pathname: "/app/documents",
					selectedCollectionId: "collection-2",
					updatedCollectionId: "collection-1",
				})
			)
		).toEqual([COLLECTIONS_QUERY_KEY]);

		expect(
			queryKeys(
				collectionUpdateSuccessInvalidationKeys({
					pathname: "/app/schemas",
					selectedCollectionId: "collection-1",
					updatedCollectionId: "collection-1",
				})
			)
		).toEqual([COLLECTIONS_QUERY_KEY]);
	});

	it("plans delete invalidation and selected-collection navigation", () => {
		expect(
			collectionDeleteSuccessEffects({
				pathname: "/app/documents",
				selectedCollectionId: "collection-1",
				deletedCollectionId: "collection-1",
			})
		).toEqual({
			invalidateQueryKeys: [COLLECTIONS_QUERY_KEY, OCR_DOCUMENTS_QUERY_KEY],
			navigateTo: "/app/documents",
		});

		expect(
			collectionDeleteSuccessEffects({
				pathname: "/app/documents",
				selectedCollectionId: "collection-2",
				deletedCollectionId: "collection-1",
			})
		).toEqual({
			invalidateQueryKeys: [COLLECTIONS_QUERY_KEY, OCR_DOCUMENTS_QUERY_KEY],
			navigateTo: null,
		});
	});

	it("runs action menu commands without invoking navigation", () => {
		const stopPropagation = vi.fn();
		const triggerEvent = { stopPropagation };
		const actionEvent = { stopPropagation: vi.fn() };
		const action = vi.fn();
		const navigate = vi.fn();

		stopCollectionMenuNavigation(triggerEvent);
		runCollectionMenuAction(actionEvent, action);

		expect(stopPropagation).toHaveBeenCalledOnce();
		expect(actionEvent.stopPropagation).toHaveBeenCalledOnce();
		expect(action).toHaveBeenCalledOnce();
		expect(navigate).not.toHaveBeenCalled();
	});

	it("composes the dropdown trigger click with propagation stopping", () => {
		const calls: string[] = [];
		const event = {
			stopPropagation: vi.fn(() => calls.push("stopPropagation")),
		};
		const triggerClick = vi.fn(() => calls.push("triggerClick"));

		composeCollectionMenuTriggerClick(triggerClick)(event);

		expect(triggerClick).toHaveBeenCalledOnce();
		expect(triggerClick).toHaveBeenCalledWith(event);
		expect(event.stopPropagation).toHaveBeenCalledOnce();
		expect(calls).toEqual(["triggerClick", "stopPropagation"]);
	});

	it("still stops propagation when the trigger has no click handler", () => {
		const event = { stopPropagation: vi.fn() };

		composeCollectionMenuTriggerClick(undefined)(event);

		expect(event.stopPropagation).toHaveBeenCalledOnce();
	});

	it("wires behavior helpers into the Svelte component", () => {
		const source = normalizeSource(navSource());

		expect(source).toContain("PERSONAL_SCHEMA_OPTIONS_QUERY_KEY");
		expect(source).toContain("fetchPersonalSchemaOptions(fetch)");
		expect(source).toContain("collectionUpdateSuccessInvalidationKeys({");
		expect(source).toContain("collectionDeleteSuccessEffects({");
		expect(source).toContain(
			"const collectionsOverflow = $derived(collectionListOverflows(collections.length))"
		);
		expect(source).toContain(
			'class={["mt-1", collectionsOverflow && "max-h-[22.25rem] overflow-y-auto pr-1"]}'
		);
		expect(source).toContain(
			"onclick={() => void retryCollections(() => collectionsQuery.refetch())}"
		);
		expect(source).toContain("href={collectionHref(collection.id)}");
		expect(source).toContain("onclick={composeCollectionMenuTriggerClick(props.onclick)}");
		expect(source).toContain(
			"onclick={(event) => runCollectionMenuAction(event, () => openEditDialog(collection))}"
		);
		expect(source).toContain(
			"onclick={(event) => runCollectionMenuAction(event, () => confirmCollectionDelete(collection))}"
		);
	});

	it("translates collection sidebar and dialog copy with Paraglide messages", () => {
		const nav = normalizeSource(navSource());
		const dialog = normalizeSource(collectionDialogSource());

		expect(nav).toContain('import { m } from "$lib/paraglide/messages.js";');
		expect(nav).toContain("m.documents_collections_nav_label()");
		expect(nav).toContain("m.documents_add_collection()");
		expect(nav).toContain("m.documents_all_documents()");
		expect(nav).toContain("m.documents_retry_collections()");
		expect(nav).toContain("m.documents_no_collections()");
		expect(nav).toContain("m.documents_collection_actions()");
		expect(nav).toContain("m.documents_delete_collection_title()");
		expect(nav).toContain("m.documents_delete_collection_description({ name: collection.name })");

		expect(dialog).toContain('import { m } from "$lib/paraglide/messages.js";');
		expect(dialog).toContain("m.documents_collection_dialog_title_new()");
		expect(dialog).toContain("m.documents_collection_dialog_title_edit()");
		expect(dialog).toContain("m.documents_collection_name_placeholder()");
		expect(dialog).toContain("m.documents_search_schemas()");
		expect(dialog).toContain("m.documents_no_schemas_found()");
	});
});
