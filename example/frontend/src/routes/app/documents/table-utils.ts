import { CalendarDateTime, getLocalTimeZone, type DateValue } from "@internationalized/date";

export type SortDirection = "asc" | "desc";
export type FileIconKind = "pdf" | "image" | "file";

export type DateRangeValue = {
	start?: DateValue;
	end?: DateValue;
};

export type DocumentListQuery = {
	collectionId?: string;
	schemaId?: string;
	filename?: string;
	createdFrom?: string;
	createdTo?: string;
	cursor?: string | null;
	size?: number | string;
	sort?: SortDirection;
};

export type CursorState = {
	currentCursor: string | null;
	history: (string | null)[];
};

export function buildDocumentsQueryPath(query: DocumentListQuery = {}) {
	const params = new URLSearchParams();

	setNonEmptyParam(params, "collection", query.collectionId);
	setNonEmptyParam(params, "schema_id", query.schemaId);
	setNonEmptyParam(params, "filename", query.filename);
	setNonEmptyParam(params, "created_from", query.createdFrom);
	setNonEmptyParam(params, "created_to", query.createdTo);
	setNonEmptyParam(params, "cursor", query.cursor);
	setNonEmptyParam(params, "size", query.size);
	setNonEmptyParam(params, "sort", query.sort);

	const queryString = params.toString();
	return queryString ? `/api/ocr/documents?${queryString}` : "/api/ocr/documents";
}

export function dateRangeToQueryBounds(
	range?: DateRangeValue,
	timeZone = getLocalTimeZone()
): Pick<DocumentListQuery, "createdFrom" | "createdTo"> {
	return {
		createdFrom: range?.start ? dateValueToISOString(range.start, timeZone, 0, 0, 0, 0) : undefined,
		createdTo: range?.end ? dateValueToISOString(range.end, timeZone, 23, 59, 59, 999) : undefined,
	};
}

export function truncateFilename(filename: string, maxLength = 20) {
	if (filename.length <= maxLength) return filename;
	if (maxLength <= 3) return ".".repeat(maxLength);

	return `${filename.slice(0, maxLength - 3)}...`;
}

export function formatCreatedDate(value: string, invalidDateLabel = "Invalid date") {
	const date = new Date(value);
	if (Number.isNaN(date.getTime())) return invalidDateLabel;

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

export function fileIconKind(mimeType: string): FileIconKind {
	const normalizedMimeType = mimeType.toLowerCase();

	if (normalizedMimeType === "application/pdf") return "pdf";
	if (normalizedMimeType.startsWith("image/")) return "image";

	return "file";
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

function dateValueToISOString(
	date: DateValue,
	timeZone: string,
	hour: number,
	minute: number,
	second: number,
	millisecond: number
) {
	return new CalendarDateTime(
		date.calendar,
		date.era,
		date.year,
		date.month,
		date.day,
		hour,
		minute,
		second,
		millisecond
	)
		.toDate(timeZone)
		.toISOString();
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
