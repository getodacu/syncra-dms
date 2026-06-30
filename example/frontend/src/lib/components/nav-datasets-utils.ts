import type { CreateDatasetInput, DatasetField, DatasetResponse } from "$lib/client/datasets";

export const DATASETS_QUERY_KEY = ["datasets"] as const;
export const DATASET_DETAIL_QUERY_KEY = ["dataset"] as const;
export const DATASET_ROWS_QUERY_KEY = ["dataset-rows"] as const;
export const DATASET_LOADING_ROWS = [0, 1, 2] as const;
export const DATASET_DIALOG_SCHEMA_PAGE_SIZE = 100;
export const DATASET_OVERFLOW_THRESHOLD = 10;

export type DatasetQueryKey =
	| typeof DATASETS_QUERY_KEY
	| readonly [...typeof DATASET_DETAIL_QUERY_KEY, string]
	| readonly [...typeof DATASET_ROWS_QUERY_KEY, string];
export type DatasetFieldTreeSelectionNode = DatasetField & {
	children: DatasetFieldTreeSelectionNode[];
};
export type DatasetDialogSchema = {
	id: string;
	name: string;
	description: string;
	schema: unknown;
};
export type DatasetDialogSchemaPage<TSchema extends DatasetDialogSchema = DatasetDialogSchema> = {
	schemas: TSchema[];
	next_cursor: string | null;
};
export type DatasetDialogMode = "create" | "edit";
export type DatasetDialogValue = CreateDatasetInput;
export type DatasetDialogState = {
	open: boolean;
	mode: DatasetDialogMode;
	editingDataset: DatasetResponse | null;
};
export type DatasetListStatus = "loading" | "error" | "empty" | "ready";
export type SubmitDatasetAction =
	| { type: "create"; input: DatasetDialogValue }
	| { type: "update"; id: string; input: DatasetDialogValue }
	| { type: "none" };

type DatasetClickEvent = Pick<Event, "stopPropagation">;
type DatasetClickHandler = (event: DatasetClickEvent) => void;

type DatasetListStatusInput = {
	isLoading: boolean;
	isError: boolean;
	datasetCount: number;
};

type DatasetUpdateSuccessInput = {
	pathname: string;
	updatedDatasetId: string;
};

type DatasetDeleteSuccessInput = {
	pathname: string;
	deletedDatasetId: string;
};

type FetchTypedDatasetDialogSchemaPage<TSchema extends DatasetDialogSchema> = (
	cursor: string | null
) => Promise<DatasetDialogSchemaPage<TSchema>>;

export function datasetHref(datasetId: string) {
	return `/app/datasets/${encodeURIComponent(datasetId)}`;
}

export function datasetDetailQueryKey(datasetId: string) {
	return [...DATASET_DETAIL_QUERY_KEY, datasetId] as const;
}

export function datasetRowsQueryKey(datasetId: string) {
	return [...DATASET_ROWS_QUERY_KEY, datasetId] as const;
}

export function isAllDatasetsActive(pathname: string) {
	return pathname === "/app/datasets";
}

export function isDatasetActive(pathname: string, datasetId: string) {
	return pathname === datasetHref(datasetId);
}

export function datasetListStatus(input: DatasetListStatusInput): DatasetListStatus {
	if (input.isLoading) return "loading";
	if (input.isError) return "error";
	if (input.datasetCount === 0) return "empty";
	return "ready";
}

export function datasetListOverflows(datasetCount: number) {
	return datasetCount > DATASET_OVERFLOW_THRESHOLD;
}

export function retryDatasets<T>(refetch: () => T): T {
	return refetch();
}

export function openCreateDatasetDialogState(): DatasetDialogState {
	return { open: true, mode: "create", editingDataset: null };
}

export function openEditDatasetDialogState(dataset: DatasetResponse): DatasetDialogState {
	return { open: true, mode: "edit", editingDataset: dataset };
}

export function closeDatasetDialogState(mode: DatasetDialogMode): DatasetDialogState {
	return { open: false, mode, editingDataset: null };
}

export function datasetDialogInitialValue(
	editingDataset: DatasetResponse | null
): DatasetDialogValue | undefined {
	if (!editingDataset) return undefined;

	return {
		name: editingDataset.name,
		schema_id: editingDataset.schema_id,
		selected_fields: editingDataset.selected_fields.map((field) => ({ ...field })),
	};
}

export function datasetDialogPending(createPending: boolean, updatePending: boolean) {
	return createPending || updatePending;
}

export function datasetDialogError(
	mode: DatasetDialogMode,
	createError: Error | null,
	updateError: Error | null
) {
	return mode === "create" ? createError : updateError;
}

export function datasetFieldsAfterSchemaChange(
	currentSchemaId: string,
	nextSchemaId: string,
	selectedFields: DatasetField[]
) {
	return currentSchemaId === nextSchemaId ? selectedFields : [];
}

export function datasetFieldNodePathMap(fieldTree: DatasetFieldTreeSelectionNode[]) {
	const nodesByPath = new Map<string, DatasetField>();

	for (const node of flattenDatasetFieldTree(fieldTree)) {
		if (!nodesByPath.has(node.path)) {
			nodesByPath.set(node.path, {
				path: node.path,
				key: node.key,
				label: node.label,
			});
		}
	}

	return nodesByPath;
}

export function validDatasetSelectedFields(
	selectedFields: DatasetField[],
	validFieldNodesByPath: ReadonlyMap<string, DatasetField>
) {
	const selectedPaths = new Set<string>();
	const validFields: DatasetField[] = [];

	for (const field of selectedFields) {
		if (selectedPaths.has(field.path)) continue;

		const currentField = validFieldNodesByPath.get(field.path);
		if (!currentField) continue;

		selectedPaths.add(field.path);
		validFields.push({ ...currentField });
	}

	return validFields;
}

export function canSubmitDatasetDialog(input: {
	pending: boolean;
	name: string;
	selectedSchemaExists: boolean;
	fieldTreeHasFields: boolean;
	validSelectedFieldCount: number;
}) {
	return (
		!input.pending &&
		input.name.trim().length > 0 &&
		input.selectedSchemaExists &&
		input.fieldTreeHasFields &&
		input.validSelectedFieldCount > 0
	);
}

export function datasetSubmitAction(
	mode: DatasetDialogMode,
	editingDataset: DatasetResponse | null,
	input: DatasetDialogValue
): SubmitDatasetAction {
	if (mode === "create") return { type: "create", input };
	if (!editingDataset) return { type: "none" };

	return { type: "update", id: editingDataset.id, input };
}

export function datasetCreateSuccessInvalidationKeys(): DatasetQueryKey[] {
	return [DATASETS_QUERY_KEY];
}

export function datasetUpdateSuccessInvalidationKeys(
	input: DatasetUpdateSuccessInput
): DatasetQueryKey[] {
	const keys: DatasetQueryKey[] = [DATASETS_QUERY_KEY];

	if (isDatasetActive(input.pathname, input.updatedDatasetId)) {
		keys.push(datasetDetailQueryKey(input.updatedDatasetId));
		keys.push(datasetRowsQueryKey(input.updatedDatasetId));
	}

	return keys;
}

export function datasetDeleteSuccessEffects(input: DatasetDeleteSuccessInput): {
	invalidateQueryKeys: DatasetQueryKey[];
	navigateTo: string | null;
} {
	const active = isDatasetActive(input.pathname, input.deletedDatasetId);
	const invalidateQueryKeys: DatasetQueryKey[] = [DATASETS_QUERY_KEY];

	if (active) {
		invalidateQueryKeys.push(datasetDetailQueryKey(input.deletedDatasetId));
		invalidateQueryKeys.push(datasetRowsQueryKey(input.deletedDatasetId));
	}

	return {
		invalidateQueryKeys,
		navigateTo: active ? "/app/datasets" : null,
	};
}

export function stopDatasetMenuNavigation(event: Pick<Event, "stopPropagation">) {
	event.stopPropagation();
}

export function composeDatasetMenuTriggerClick(triggerClick: unknown): DatasetClickHandler {
	return (event) => {
		if (typeof triggerClick === "function") {
			(triggerClick as DatasetClickHandler)(event);
		}

		stopDatasetMenuNavigation(event);
	};
}

export function runDatasetMenuAction(event: Pick<Event, "stopPropagation">, action: () => void) {
	stopDatasetMenuNavigation(event);
	action();
}

export async function fetchAllDatasetDialogSchemas<TSchema extends DatasetDialogSchema>(
	fetchPage: FetchTypedDatasetDialogSchemaPage<TSchema>
): Promise<TSchema[]> {
	const schemas: TSchema[] = [];
	const seenCursors = new Set<string>();
	let cursor: string | null = null;

	for (;;) {
		const page = await fetchPage(cursor);
		schemas.push(...page.schemas);

		if (!page.next_cursor) return schemas;
		if (seenCursors.has(page.next_cursor)) {
			throw new Error("Schema pagination loop detected");
		}

		seenCursors.add(page.next_cursor);
		cursor = page.next_cursor;
	}
}

function flattenDatasetFieldTree(
	nodes: DatasetFieldTreeSelectionNode[]
): DatasetFieldTreeSelectionNode[] {
	return nodes.flatMap((node) => [node, ...flattenDatasetFieldTree(node.children)]);
}
