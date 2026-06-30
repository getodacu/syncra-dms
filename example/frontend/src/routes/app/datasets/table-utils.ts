import { CalendarDateTime, getLocalTimeZone, type DateValue } from "@internationalized/date";
import { isDatasetNotFoundError } from "$lib/client/datasets";

export type SortDirection = "asc" | "desc";

export const DATASET_ROWS_QUERY_RETRY_LIMIT = 3;

export type CursorState = {
	currentCursor: string | null;
	history: (string | null)[];
};

export type DatasetExportFormat = "csv" | "xlsx";

export type DateRangeValue = {
	start?: DateValue;
	end?: DateValue;
};

export function dateRangeToQueryBounds(range?: DateRangeValue, timeZone = getLocalTimeZone()) {
	return {
		createdFrom: range?.start ? dateValueToISOString(range.start, timeZone, 0, 0, 0, 0) : undefined,
		createdTo: range?.end ? dateValueToISOString(range.end, timeZone, 23, 59, 59, 999) : undefined,
	};
}

export function datasetCellText(value: unknown): string {
	if (value === null || value === undefined) return "";
	if (typeof value === "string") return value;
	if (typeof value === "number" || typeof value === "boolean") return String(value);

	if (Array.isArray(value) || isRecord(value)) {
		try {
			return JSON.stringify(value);
		} catch {
			return String(value);
		}
	}

	return String(value);
}

export function datasetExportFilename(
	name: string,
	format: DatasetExportFormat,
	date: Date = new Date()
): string {
	const safeName = sanitizeFilenameBase(name) || "dataset";
	const timestamp = exportTimestamp(date);

	return `${safeName} ${timestamp}.${format}`;
}

export function formatDatasetDate(value: string, invalidDateLabel = "Invalid date"): string {
	const date = new Date(value);
	if (Number.isNaN(date.getTime())) return invalidDateLabel;

	return date.toISOString().slice(0, 10);
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

export function shouldRetryDatasetRowsQuery(failureCount: number, error: unknown) {
	return !isDatasetNotFoundError(error) && failureCount < DATASET_ROWS_QUERY_RETRY_LIMIT;
}

function sanitizeFilenameBase(value: string): string {
	return value
		.replace(/[\x00-\x1f\x7f]/g, " ")
		.replace(/[<>:"/\\|?*]+/g, " ")
		.replace(/\s+/g, " ")
		.replace(/[. ]+$/g, "")
		.trim();
}

function exportTimestamp(date: Date): string {
	if (Number.isNaN(date.getTime())) return "unknown-date";

	return date.toISOString().slice(0, 19).replace("T", " ").replace(/:/g, "-");
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

function isRecord(value: unknown): value is Record<string, unknown> {
	return typeof value === "object" && value !== null && !Array.isArray(value);
}
