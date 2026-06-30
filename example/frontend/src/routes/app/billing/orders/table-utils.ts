import { CalendarDateTime, getLocalTimeZone, type DateValue } from "@internationalized/date";

import { formatCents } from "$lib/billing/pricing";

export type SortDirection = "asc" | "desc";
export type BillingOrderStatus = "pending" | "paid" | "failed" | "refunded" | "canceled";

export type DateRangeValue = {
	start?: DateValue;
	end?: DateValue;
};

export type BillingOrdersListQuery = {
	status?: BillingOrderStatus;
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

export function buildBillingOrdersQueryPath(query: BillingOrdersListQuery = {}) {
	const params = new URLSearchParams();

	setNonEmptyParam(params, "status", query.status);
	setNonEmptyParam(params, "created_from", query.createdFrom);
	setNonEmptyParam(params, "created_to", query.createdTo);
	setNonEmptyParam(params, "cursor", query.cursor);
	setNonEmptyParam(params, "size", query.size);
	setNonEmptyParam(params, "sort", query.sort);

	const queryString = params.toString();
	return queryString ? `/api/billing/orders?${queryString}` : "/api/billing/orders";
}

export function dateRangeToQueryBounds(
	range?: DateRangeValue,
	timeZone = getLocalTimeZone()
): Pick<BillingOrdersListQuery, "createdFrom" | "createdTo"> {
	return {
		createdFrom: range?.start ? dateValueToISOString(range.start, timeZone, 0, 0, 0, 0) : undefined,
		createdTo: range?.end ? dateValueToISOString(range.end, timeZone, 23, 59, 59, 999) : undefined,
	};
}

export function formatOrderDate(value: string) {
	const date = new Date(value);
	if (Number.isNaN(date.getTime())) return "Invalid date";

	return date.toISOString().slice(0, 10);
}

export function formatPaymentDateTime(value?: string) {
	if (!value) return "-";

	const date = new Date(value);
	if (Number.isNaN(date.getTime())) return "Invalid date";

	return new Intl.DateTimeFormat("en-US", {
		year: "numeric",
		month: "short",
		day: "numeric",
		hour: "numeric",
		minute: "2-digit"
	}).format(date);
}

export function formatRawCents(value: number) {
	if (!Number.isFinite(value)) return "0";

	return Math.trunc(value).toLocaleString("en-US");
}

export function formatCredits(value: number) {
	if (!Number.isFinite(value)) return "0";

	return Math.trunc(value).toLocaleString("en-US");
}

export function formatCurrencyAmount(amountCents: number, currency: string) {
	const normalizedCurrency = currency.trim().toUpperCase() || "EUR";
	try {
		return formatCents(amountCents, normalizedCurrency);
	} catch {
		return formatCents(amountCents);
	}
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
