import { CalendarDate } from "@internationalized/date";
import { describe, expect, it } from "vitest";

import {
	buildCreditUsageHistoryQueryPath,
	cursorNextState,
	cursorPreviousState,
	dateRangeToQueryBounds,
	formatCreditsDelta,
	formatCreditUsageHistoryDate,
	resetCursorState,
} from "./table-utils";

describe("credit usage history table utils", () => {
	it("builds credit usage history paths with only non-empty query parameters", () => {
		expect(buildCreditUsageHistoryQueryPath({ size: 20, sort: "desc" })).toBe(
			"/api/billing/credit-usage-history?size=20&sort=desc"
		);
		expect(
			buildCreditUsageHistoryQueryPath({
				type: "debit",
				cursor: "cursor-1",
				size: "50",
				sort: "asc",
			})
		).toBe("/api/billing/credit-usage-history?type=debit&cursor=cursor-1&size=50&sort=asc");
		expect(
			buildCreditUsageHistoryQueryPath({
				type: "purchase",
				createdFrom: "2026-06-04T00:00:00.000Z",
				createdTo: "2026-06-04T23:59:59.999Z",
			})
		).toBe(
			"/api/billing/credit-usage-history?type=purchase&created_from=2026-06-04T00%3A00%3A00.000Z&created_to=2026-06-04T23%3A59%3A59.999Z"
		);
		expect(buildCreditUsageHistoryQueryPath({ cursor: "", size: " 50 ", sort: "asc" })).toBe(
			"/api/billing/credit-usage-history?size=50&sort=asc"
		);
		expect(buildCreditUsageHistoryQueryPath()).toBe("/api/billing/credit-usage-history");
	});

	it("converts local range calendar values into inclusive RFC3339 bounds", () => {
		expect(
			dateRangeToQueryBounds(
				{ start: new CalendarDate(2026, 6, 4), end: new CalendarDate(2026, 6, 5) },
				"UTC"
			)
		).toEqual({
			createdFrom: "2026-06-04T00:00:00.000Z",
			createdTo: "2026-06-05T23:59:59.999Z",
		});
		expect(dateRangeToQueryBounds({ start: new CalendarDate(2026, 6, 4) }, "UTC")).toEqual({
			createdFrom: "2026-06-04T00:00:00.000Z",
			createdTo: undefined,
		});
		expect(dateRangeToQueryBounds({ end: new CalendarDate(2026, 6, 5) }, "UTC")).toEqual({
			createdFrom: undefined,
			createdTo: "2026-06-05T23:59:59.999Z",
		});
		expect(dateRangeToQueryBounds(undefined, "UTC")).toEqual({
			createdFrom: undefined,
			createdTo: undefined,
		});
	});

	it("formats display values", () => {
		expect(formatCreditUsageHistoryDate("2026-06-04T12:30:00Z")).toBe("2026-06-04");
		expect(formatCreditUsageHistoryDate("bad")).toBe("Invalid date");
		expect(formatCreditsDelta(1000)).toBe("+1,000");
		expect(formatCreditsDelta(-25)).toBe("-25");
		expect(formatCreditsDelta(0)).toBe("0");
		expect(formatCreditsDelta(Number.NaN)).toBe("0");
	});

	it("maintains cursor history", () => {
		const first = cursorNextState(resetCursorState(), "cursor-1");
		const second = cursorNextState(first, "cursor-2");

		expect(first).toEqual({ currentCursor: "cursor-1", history: [null] });
		expect(second).toEqual({ currentCursor: "cursor-2", history: [null, "cursor-1"] });
		expect(cursorPreviousState(second)).toEqual({ currentCursor: "cursor-1", history: [null] });
		expect(cursorPreviousState(first)).toEqual({ currentCursor: null, history: [] });
		expect(cursorNextState(first, "")).toBe(first);
		expect(cursorNextState(first, "   ")).toBe(first);
	});
});
