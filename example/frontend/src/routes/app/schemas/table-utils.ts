export type SortDirection = "asc" | "desc";

export type SchemaListQuery = {
	cursor?: string | null;
	size?: number | string;
	sort?: SortDirection;
};

export type CursorState = {
	currentCursor: string | null;
	history: (string | null)[];
};

export function buildSchemasQueryPath(query: SchemaListQuery = {}): string {
	const params = new URLSearchParams();

	params.set("scope", "mine");
	setNonEmptyParam(params, "cursor", query.cursor);
	setNonEmptyParam(params, "size", query.size ?? 20);
	setNonEmptyParam(params, "sort", query.sort);

	return `/api/schemas?${params.toString()}`;
}

export function formatDate(value: string): string {
	const date = new Date(value);
	if (Number.isNaN(date.getTime())) return "Invalid date";

	return date.toISOString().slice(0, 10);
}

export function headerSelectionState(visibleIds: string[], selectedIds: Set<string>) {
	const selectedCount = visibleIds.filter((id) => selectedIds.has(id)).length;

	return {
		checked: visibleIds.length > 0 && selectedCount === visibleIds.length,
		indeterminate: selectedCount > 0 && selectedCount < visibleIds.length,
		selectedCount,
	};
}

export function toggleSelection(selectedIds: Set<string>, id: string, checked: boolean) {
	const nextSelectedIds = new Set(selectedIds);

	if (checked) {
		nextSelectedIds.add(id);
	} else {
		nextSelectedIds.delete(id);
	}

	return nextSelectedIds;
}

export function togglePageSelection(
	visibleIds: string[],
	selectedIds: Set<string>,
	checked: boolean
) {
	const nextSelectedIds = new Set(selectedIds);

	for (const id of visibleIds) {
		if (checked) {
			nextSelectedIds.add(id);
		} else {
			nextSelectedIds.delete(id);
		}
	}

	return nextSelectedIds;
}

export function resetCursorState(): CursorState {
	return { currentCursor: null, history: [] };
}

export function cursorNextState(state: CursorState, nextCursor?: string | null): CursorState {
	const cursor = nextCursor?.trim();

	if (!cursor) return state;

	return {
		currentCursor: cursor,
		history: [...state.history, state.currentCursor],
	};
}

export function cursorPreviousState(state: CursorState): CursorState {
	if (state.history.length === 0) {
		return {
			currentCursor: state.currentCursor,
			history: [],
		};
	}

	const history = state.history.slice(0, -1);
	const currentCursor = state.history[state.history.length - 1] ?? null;

	return { currentCursor, history };
}

function setNonEmptyParam(
	params: URLSearchParams,
	key: string,
	value: string | number | null | undefined
) {
	if (value === null || value === undefined) return;

	const stringValue = String(value).trim();
	if (!stringValue) return;

	params.set(key, stringValue);
}
