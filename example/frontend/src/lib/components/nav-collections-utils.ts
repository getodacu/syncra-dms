import type { CollectionInput, CollectionResponse } from "$lib/client/collections";

export const COLLECTIONS_QUERY_KEY = ["collections"] as const;
export const OCR_DOCUMENTS_QUERY_KEY = ["ocr-documents"] as const;
export const COLLECTION_LOADING_ROWS = [0, 1, 2] as const;
export const COLLECTION_OVERFLOW_THRESHOLD = 10;

export type CollectionQueryKey =
	| typeof COLLECTIONS_QUERY_KEY
	| typeof OCR_DOCUMENTS_QUERY_KEY;
export type CollectionDialogMode = "create" | "edit";
export type CollectionDialogState = {
	open: boolean;
	mode: CollectionDialogMode;
	editingCollection: CollectionResponse | null;
};
export type CollectionListStatus = "loading" | "error" | "empty" | "ready";
export type SubmitCollectionAction =
	| { type: "create"; input: CollectionInput }
	| { type: "update"; id: string; input: CollectionInput }
	| { type: "none" };
type CollectionClickEvent = Pick<Event, "stopPropagation">;
type CollectionClickHandler = (event: CollectionClickEvent) => void;

type CollectionRouteState = {
	pathname: string;
	selectedCollectionId: string | null;
};

type CollectionListStatusInput = {
	isLoading: boolean;
	isError: boolean;
	collectionCount: number;
};

type CollectionUpdateSuccessInput = CollectionRouteState & {
	updatedCollectionId: string;
};

type CollectionDeleteSuccessInput = CollectionRouteState & {
	deletedCollectionId: string;
};

export function collectionHref(collectionId: string) {
	const params = new URLSearchParams({ collection: collectionId });
	return `/app/documents?${params.toString()}`;
}

export function isAllDocumentsActive(pathname: string, selectedCollectionId: string | null) {
	return isDocumentsPage(pathname) && !selectedCollectionId;
}

export function isCollectionActive(
	pathname: string,
	selectedCollectionId: string | null,
	collectionId: string
) {
	return isDocumentsPage(pathname) && selectedCollectionId === collectionId;
}

export function collectionListStatus(input: CollectionListStatusInput): CollectionListStatus {
	if (input.isLoading) return "loading";
	if (input.isError) return "error";
	if (input.collectionCount === 0) return "empty";
	return "ready";
}

export function collectionListOverflows(collectionCount: number) {
	return collectionCount > COLLECTION_OVERFLOW_THRESHOLD;
}

export function retryCollections<T>(refetch: () => T): T {
	return refetch();
}

export function openCreateCollectionDialogState(): CollectionDialogState {
	return { open: true, mode: "create", editingCollection: null };
}

export function openEditCollectionDialogState(
	collection: CollectionResponse
): CollectionDialogState {
	return { open: true, mode: "edit", editingCollection: collection };
}

export function closeCollectionDialogState(mode: CollectionDialogMode): CollectionDialogState {
	return { open: false, mode, editingCollection: null };
}

export function collectionDialogInitialValue(
	editingCollection: CollectionResponse | null
): CollectionInput | undefined {
	if (!editingCollection) return undefined;

	return {
		name: editingCollection.name,
		schema_ids: editingCollection.schema_ids,
	};
}

export function collectionDialogPending(createPending: boolean, updatePending: boolean) {
	return createPending || updatePending;
}

export function collectionDialogError(
	mode: CollectionDialogMode,
	createError: Error | null,
	updateError: Error | null
) {
	return mode === "create" ? createError : updateError;
}

export function collectionSubmitAction(
	mode: CollectionDialogMode,
	editingCollection: CollectionResponse | null,
	input: CollectionInput
): SubmitCollectionAction {
	if (mode === "create") return { type: "create", input };
	if (!editingCollection) return { type: "none" };

	return { type: "update", id: editingCollection.id, input };
}

export function collectionUpdateSuccessInvalidationKeys(
	input: CollectionUpdateSuccessInput
): CollectionQueryKey[] {
	const keys: CollectionQueryKey[] = [COLLECTIONS_QUERY_KEY];

	if (
		isCollectionActive(input.pathname, input.selectedCollectionId, input.updatedCollectionId)
	) {
		keys.push(OCR_DOCUMENTS_QUERY_KEY);
	}

	return keys;
}

export function collectionDeleteSuccessEffects(input: CollectionDeleteSuccessInput): {
	invalidateQueryKeys: CollectionQueryKey[];
	navigateTo: string | null;
} {
	return {
		invalidateQueryKeys: [COLLECTIONS_QUERY_KEY, OCR_DOCUMENTS_QUERY_KEY],
		navigateTo: isCollectionActive(
			input.pathname,
			input.selectedCollectionId,
			input.deletedCollectionId
		)
			? "/app/documents"
			: null,
	};
}

export function stopCollectionMenuNavigation(event: Pick<Event, "stopPropagation">) {
	event.stopPropagation();
}

export function composeCollectionMenuTriggerClick(
	triggerClick: unknown
): CollectionClickHandler {
	return (event) => {
		if (typeof triggerClick === "function") {
			(triggerClick as CollectionClickHandler)(event);
		}

		stopCollectionMenuNavigation(event);
	};
}

export function runCollectionMenuAction(
	event: Pick<Event, "stopPropagation">,
	action: () => void
) {
	stopCollectionMenuNavigation(event);
	action();
}

function isDocumentsPage(pathname: string) {
	return pathname === "/app/documents";
}
