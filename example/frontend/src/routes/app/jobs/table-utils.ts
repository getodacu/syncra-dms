export type SortDirection = "asc" | "desc";
export type JobStatus = "queued" | "pending" | "processing" | "completed" | "failed" | "unknown";

export type JobListQuery = {
	status?: string;
	cursor?: string | null;
	size?: number | string;
	sort?: SortDirection;
};

export type CursorState = {
	currentCursor: string | null;
	history: (string | null)[];
};

export function buildJobsQueryPath(query: JobListQuery = {}) {
	const params = new URLSearchParams();

	setNonEmptyParam(params, "status", query.status);
	setNonEmptyParam(params, "cursor", query.cursor);
	setNonEmptyParam(params, "size", query.size);
	setNonEmptyParam(params, "sort", query.sort);

	const queryString = params.toString();
	return queryString ? `/api/ocr/jobs?${queryString}` : "/api/ocr/jobs";
}

export function formatCreatedDate(value: string) {
	const date = new Date(value);
	if (Number.isNaN(date.getTime())) return "Invalid date";

	return date.toISOString().slice(0, 10);
}

export function formatFileSize(bytes: number) {
	if (!Number.isFinite(bytes) || bytes <= 0) return "0 B";

	const units = ["B", "KB", "MB", "GB"] as const;
	let size = bytes;
	let unitIndex = 0;

	while (size >= 1024 && unitIndex < units.length - 1) {
		size /= 1024;
		unitIndex += 1;
	}

	const formattedSize =
		unitIndex === 0 ? String(Math.round(size)) : size.toFixed(1).replace(/\.0$/, "");

	return `${formattedSize} ${units[unitIndex]}`;
}

export function normalizeJobStatus(status: string): JobStatus {
	if (
		status === "queued" ||
		status === "pending" ||
		status === "processing" ||
		status === "completed" ||
		status === "failed"
	) {
		return status;
	}

	return "unknown";
}

export function isTerminalJobStatus(status: string) {
	return status === "completed" || status === "failed";
}

export function shouldPollJobStatus(status: string) {
	const normalized = normalizeJobStatus(status);
	return normalized === "queued" || normalized === "pending" || normalized === "processing";
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
