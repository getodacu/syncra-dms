import { CalendarDateTime, getLocalTimeZone, type DateValue } from "@internationalized/date";

export type SortDirection = "asc" | "desc";
export type CreditUsageHistoryEntryType = "purchase" | "debit";

export type DateRangeValue = {
	start?: DateValue;
	end?: DateValue;
};

export type CreditUsageHistoryListQuery = {
	type?: CreditUsageHistoryEntryType;
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

export function buildCreditUsageHistoryQueryPath(query: CreditUsageHistoryListQuery = {}) {
	const params = new URLSearchParams();

	setNonEmptyParam(params, "type", query.type);
	setNonEmptyParam(params, "created_from", query.createdFrom);
	setNonEmptyParam(params, "created_to", query.createdTo);
	setNonEmptyParam(params, "cursor", query.cursor);
	setNonEmptyParam(params, "size", query.size);
	setNonEmptyParam(params, "sort", query.sort);

	const queryString = params.toString();
	return queryString
		? `/api/billing/credit-usage-history?${queryString}`
		: "/api/billing/credit-usage-history";
}

export function dateRangeToQueryBounds(
	range?: DateRangeValue,
	timeZone = getLocalTimeZone()
): Pick<CreditUsageHistoryListQuery, "createdFrom" | "createdTo"> {
	return {
		createdFrom: range?.start ? dateValueToISOString(range.start, timeZone, 0, 0, 0, 0) : undefined,
		createdTo: range?.end ? dateValueToISOString(range.end, timeZone, 23, 59, 59, 999) : undefined,
	};
}

export function formatCreditUsageHistoryDate(value: string) {
	const date = new Date(value);
	if (Number.isNaN(date.getTime())) return "Invalid date";

	return date.toISOString().slice(0, 10);
}

export function formatCreditsDelta(value: number) {
	if (!Number.isFinite(value)) return "0";
	if (value === 0) return "0";

	const formatted = Math.abs(value).toLocaleString("en-US");
	return value > 0 ? `+${formatted}` : `-${formatted}`;
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
