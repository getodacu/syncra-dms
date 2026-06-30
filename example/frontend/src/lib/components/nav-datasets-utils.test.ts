import { readFileSync } from "node:fs";
import { describe, expect, it, vi } from "vitest";

import type { CreateDatasetInput, DatasetResponse } from "$lib/client/datasets";
import {
	DATASET_DETAIL_QUERY_KEY,
	DATASET_DIALOG_SCHEMA_PAGE_SIZE,
	DATASET_LOADING_ROWS,
	DATASET_OVERFLOW_THRESHOLD,
	DATASET_ROWS_QUERY_KEY,
	DATASETS_QUERY_KEY,
	canSubmitDatasetDialog,
	closeDatasetDialogState,
	composeDatasetMenuTriggerClick,
	datasetCreateSuccessInvalidationKeys,
	datasetDeleteSuccessEffects,
	datasetDialogError,
	datasetDialogInitialValue,
	datasetDialogPending,
	datasetFieldNodePathMap,
	datasetFieldsAfterSchemaChange,
	datasetHref,
	datasetListOverflows,
	datasetListStatus,
	datasetSubmitAction,
	datasetUpdateSuccessInvalidationKeys,
	fetchAllDatasetDialogSchemas,
	isAllDatasetsActive,
	isDatasetActive,
	openCreateDatasetDialogState,
	openEditDatasetDialogState,
	retryDatasets,
	runDatasetMenuAction,
	stopDatasetMenuNavigation,
	validDatasetSelectedFields,
	type DatasetQueryKey,
} from "./nav-datasets-utils";

const dataset = (patch: Partial<DatasetResponse> = {}): DatasetResponse => ({
	id: "dataset-1",
	created_at: "2026-06-01T00:00:00.000Z",
	updated_at: "2026-06-02T00:00:00.000Z",
	user_id: "user-1",
	name: "Invoice dataset",
	schema_id: "schema-1",
	schema_name: "Invoice",
	selected_fields: [{ path: "/total", key: "total", label: "total" }],
	field_count: 1,
	...patch,
});

const input: CreateDatasetInput = {
	name: "Receipts",
	schema_id: "schema-2",
	selected_fields: [{ path: "/merchant", key: "merchant", label: "merchant" }],
};

const queryKeys = (keys: DatasetQueryKey[]) => keys.map((key) => [...key]);
const navSource = () =>
	readFileSync(new URL("./nav-datasets.svelte", import.meta.url), "utf8");
const dialogSource = () =>
	readFileSync(new URL("./dataset-dialog.svelte", import.meta.url), "utf8");
const normalizeSource = (value: string) => value.replace(/\s+/g, " ").trim();

describe("datasets sidebar behavior", () => {
	it("derives dataset hrefs and active links from encoded route state", () => {
		expect(isAllDatasetsActive("/app/datasets")).toBe(true);
		expect(isAllDatasetsActive("/app/datasets/")).toBe(false);
		expect(isAllDatasetsActive("/app/datasets/dataset-1")).toBe(false);
		expect(isAllDatasetsActive("/app/documents")).toBe(false);

		expect(datasetHref("dataset with space")).toBe("/app/datasets/dataset%20with%20space");
		expect(datasetHref("dataset/1")).toBe("/app/datasets/dataset%2F1");
		expect(isDatasetActive("/app/datasets/dataset%2F1", "dataset/1")).toBe(true);
		expect(isDatasetActive("/app/datasets/dataset/1", "dataset/1")).toBe(false);
		expect(isDatasetActive("/app/datasets/dataset%2F1/rows", "dataset/1")).toBe(false);
	});

	it("models the add/create dialog flow", () => {
		expect(openCreateDatasetDialogState()).toEqual({
			open: true,
			mode: "create",
			editingDataset: null,
		});
		expect(datasetDialogInitialValue(null)).toBeUndefined();
		expect(datasetSubmitAction("create", null, input)).toEqual({ type: "create", input });
		expect(datasetCreateSuccessInvalidationKeys()).toEqual([DATASETS_QUERY_KEY]);
		expect(closeDatasetDialogState("create")).toEqual({
			open: false,
			mode: "create",
			editingDataset: null,
		});
	});

	it("models edit dialog state and update submit actions", () => {
		const item = dataset({
			id: "dataset-7",
			name: "Contracts",
			schema_id: "schema-a",
			selected_fields: [{ path: "/vendor/name", key: "vendor_name", label: "name" }],
		});

		expect(openEditDatasetDialogState(item)).toEqual({
			open: true,
			mode: "edit",
			editingDataset: item,
		});
		expect(datasetDialogInitialValue(item)).toEqual({
			name: "Contracts",
			schema_id: "schema-a",
			selected_fields: [{ path: "/vendor/name", key: "vendor_name", label: "name" }],
		});
		expect(datasetDialogInitialValue(item)?.selected_fields).not.toBe(item.selected_fields);
		expect(datasetSubmitAction("edit", item, input)).toEqual({
			type: "update",
			id: "dataset-7",
			input,
		});
		expect(datasetSubmitAction("edit", null, input)).toEqual({ type: "none" });
	});

	it("selects dialog pending and error state by mutation mode", () => {
		const createError = new Error("create failed");
		const updateError = new Error("update failed");

		expect(datasetDialogPending(false, false)).toBe(false);
		expect(datasetDialogPending(true, false)).toBe(true);
		expect(datasetDialogPending(false, true)).toBe(true);
		expect(datasetDialogError("create", createError, updateError)).toBe(createError);
		expect(datasetDialogError("edit", createError, updateError)).toBe(updateError);
	});

	it("clears selected fields when the schema changes", () => {
		const selectedFields = [{ path: "/total", key: "total", label: "total" }];

		expect(datasetFieldsAfterSchemaChange("schema-1", "schema-1", selectedFields)).toBe(
			selectedFields
		);
		expect(datasetFieldsAfterSchemaChange("schema-1", "schema-2", selectedFields)).toEqual([]);
	});

	it("normalizes selected fields against the current field tree", () => {
		const fieldTree = [
			{
				path: "/total",
				key: "current_total",
				label: "Total",
				children: [],
			},
			{
				path: "/vendor",
				key: "vendor",
				label: "Vendor",
				children: [
					{
						path: "/vendor/name",
						key: "current_vendor_name",
						label: "Vendor name",
						children: [],
					},
				],
			},
		];
		const currentFields = datasetFieldNodePathMap(fieldTree);

		expect([...currentFields.keys()]).toEqual(["/total", "/vendor", "/vendor/name"]);
		expect(
			validDatasetSelectedFields(
				[
					{ path: "/total", key: "stale_total", label: "stale total" },
					{ path: "/missing", key: "missing", label: "Missing" },
					{ path: "/total", key: "duplicate_total", label: "Duplicate total" },
					{ path: "/vendor/name", key: "stale_vendor_name", label: "stale vendor" },
				],
				currentFields
			)
		).toEqual([
			{ path: "/total", key: "current_total", label: "Total" },
			{ path: "/vendor/name", key: "current_vendor_name", label: "Vendor name" },
		]);
	});

	it("blocks dataset dialog submit until the current schema and valid fields are ready", () => {
		const validInput = {
			pending: false,
			name: "Dataset",
			selectedSchemaExists: true,
			fieldTreeHasFields: true,
			validSelectedFieldCount: 1,
		};

		expect(canSubmitDatasetDialog(validInput)).toBe(true);
		expect(canSubmitDatasetDialog({ ...validInput, pending: true })).toBe(false);
		expect(canSubmitDatasetDialog({ ...validInput, name: " " })).toBe(false);
		expect(canSubmitDatasetDialog({ ...validInput, selectedSchemaExists: false })).toBe(false);
		expect(canSubmitDatasetDialog({ ...validInput, fieldTreeHasFields: false })).toBe(false);
		expect(canSubmitDatasetDialog({ ...validInput, validSelectedFieldCount: 0 })).toBe(false);
	});

	it("loads all dataset dialog schemas across pages", async () => {
		const firstSchema = {
			id: "schema-1",
			name: "Schema 1",
			description: "",
			schema: { type: "object", properties: {} },
		};
		const secondSchema = {
			id: "schema-2",
			name: "Schema 2",
			description: "",
			schema: { type: "object", properties: {} },
		};
		const fetchPage = vi
			.fn()
			.mockResolvedValueOnce({ schemas: [firstSchema], next_cursor: "cursor-2" })
			.mockResolvedValueOnce({ schemas: [secondSchema], next_cursor: null });

		await expect(fetchAllDatasetDialogSchemas(fetchPage)).resolves.toEqual([
			firstSchema,
			secondSchema,
		]);
		expect(fetchPage).toHaveBeenNthCalledWith(1, null);
		expect(fetchPage).toHaveBeenNthCalledWith(2, "cursor-2");
		expect(DATASET_DIALOG_SCHEMA_PAGE_SIZE).toBe(100);
	});

	it("fails schema pagination loops instead of polling forever", async () => {
		const fetchPage = vi
			.fn()
			.mockResolvedValueOnce({ schemas: [], next_cursor: "cursor-1" })
			.mockResolvedValueOnce({ schemas: [], next_cursor: "cursor-1" });

		await expect(fetchAllDatasetDialogSchemas(fetchPage)).rejects.toThrow(
			"Schema pagination loop detected"
		);
	});

	it("models loading, retry, error, empty, and ready list states", async () => {
		const refetch = vi.fn(() => Promise.resolve("refetched"));

		expect(DATASET_LOADING_ROWS).toEqual([0, 1, 2]);
		expect(datasetListStatus({ isLoading: true, isError: true, datasetCount: 1 })).toBe(
			"loading"
		);
		expect(datasetListStatus({ isLoading: false, isError: true, datasetCount: 1 })).toBe(
			"error"
		);
		expect(datasetListStatus({ isLoading: false, isError: false, datasetCount: 0 })).toBe(
			"empty"
		);
		expect(datasetListStatus({ isLoading: false, isError: false, datasetCount: 1 })).toBe(
			"ready"
		);
		await expect(retryDatasets(refetch)).resolves.toBe("refetched");
		expect(refetch).toHaveBeenCalledOnce();
	});

	it("only overflows dataset navigation after ten datasets", () => {
		expect(DATASET_OVERFLOW_THRESHOLD).toBe(10);
		expect(datasetListOverflows(0)).toBe(false);
		expect(datasetListOverflows(10)).toBe(false);
		expect(datasetListOverflows(11)).toBe(true);
	});

	it("invalidates rows and detail after updating the active dataset only", () => {
		expect(
			queryKeys(
				datasetUpdateSuccessInvalidationKeys({
					pathname: "/app/datasets/dataset-1",
					updatedDatasetId: "dataset-1",
				})
			)
		).toEqual([
			DATASETS_QUERY_KEY,
			[...DATASET_DETAIL_QUERY_KEY, "dataset-1"],
			[...DATASET_ROWS_QUERY_KEY, "dataset-1"],
		]);

		expect(
			queryKeys(
				datasetUpdateSuccessInvalidationKeys({
					pathname: "/app/datasets/dataset-2",
					updatedDatasetId: "dataset-1",
				})
			)
		).toEqual([DATASETS_QUERY_KEY]);
	});

	it("plans delete invalidation and active-dataset navigation", () => {
		expect(
			datasetDeleteSuccessEffects({
				pathname: "/app/datasets/dataset-1",
				deletedDatasetId: "dataset-1",
			})
		).toEqual({
			invalidateQueryKeys: [
				DATASETS_QUERY_KEY,
				[...DATASET_DETAIL_QUERY_KEY, "dataset-1"],
				[...DATASET_ROWS_QUERY_KEY, "dataset-1"],
			],
			navigateTo: "/app/datasets",
		});

		expect(
			datasetDeleteSuccessEffects({
				pathname: "/app/datasets/dataset-2",
				deletedDatasetId: "dataset-1",
			})
		).toEqual({
			invalidateQueryKeys: [DATASETS_QUERY_KEY],
			navigateTo: null,
		});
	});

	it("runs action menu commands without invoking navigation", () => {
		const stopPropagation = vi.fn();
		const triggerEvent = { stopPropagation };
		const actionEvent = { stopPropagation: vi.fn() };
		const action = vi.fn();
		const navigate = vi.fn();

		stopDatasetMenuNavigation(triggerEvent);
		runDatasetMenuAction(actionEvent, action);

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

		composeDatasetMenuTriggerClick(triggerClick)(event);

		expect(triggerClick).toHaveBeenCalledOnce();
		expect(triggerClick).toHaveBeenCalledWith(event);
		expect(event.stopPropagation).toHaveBeenCalledOnce();
		expect(calls).toEqual(["triggerClick", "stopPropagation"]);
	});

	it("still stops propagation when the trigger has no click handler", () => {
		const event = { stopPropagation: vi.fn() };

		composeDatasetMenuTriggerClick(undefined)(event);

		expect(event.stopPropagation).toHaveBeenCalledOnce();
	});

	it("wires behavior helpers into the Svelte component sources", () => {
		const nav = normalizeSource(navSource());
		const dialog = normalizeSource(dialogSource());

		expect(nav).toContain('import { m } from "$lib/paraglide/messages.js";');
		expect(nav).toContain("fetchDatasets(fetch)");
		expect(nav).toContain("fetchAllDatasetDialogSchemas((cursor) => fetchSchemas(fetch, { cursor, size: DATASET_DIALOG_SCHEMA_PAGE_SIZE })");
		expect(nav).toContain("datasetCreateSuccessInvalidationKeys()");
		expect(nav).toContain("datasetUpdateSuccessInvalidationKeys({");
		expect(nav).toContain("datasetDeleteSuccessEffects({");
		expect(nav).toContain(
			"const datasetsOverflow = $derived(datasetListOverflows(datasets.length))"
		);
		expect(nav).toContain(
			'class={["mt-1", datasetsOverflow && "max-h-[22.25rem] overflow-y-auto pr-1"]}'
		);
		expect(nav).toContain("href={datasetHref(dataset.id)}");
		expect(nav).toContain("onclick={composeDatasetMenuTriggerClick(props.onclick)}");
		expect(nav).toContain(
			"onclick={(event) => runDatasetMenuAction(event, () => openEditDialog(dataset))}"
		);
		expect(nav).toContain(
			"onclick={(event) => runDatasetMenuAction(event, () => confirmDatasetDelete(dataset))}"
		);
		expect(nav).toContain("m.datasets_add_dataset()");
		expect(nav).toContain("m.datasets_delete_failed()");

		expect(dialog).toContain('import { m } from "$lib/paraglide/messages.js";');
		expect(dialog).toContain("buildDatasetFieldTree(selectedSchema?.schema)");
		expect(dialog).toContain("datasetFieldNodePathMap(fieldTree)");
		expect(dialog).toContain("validDatasetSelectedFields(selectedFields, validFieldNodesByPath)");
		expect(dialog).toContain("datasetFieldsAfterSchemaChange");
		expect(dialog).toContain("selected_fields: validSelectedFields.map");
		expect(dialog).toContain("m.datasets_dialog_title_new()");
		expect(dialog).toContain("m.datasets_no_fields_selected()");
	});
});
