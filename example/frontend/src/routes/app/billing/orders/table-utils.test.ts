import { CalendarDate } from "@internationalized/date";
import { describe, expect, it } from "vitest";

import {
	buildBillingOrdersQueryPath,
	cursorNextState,
	cursorPreviousState,
	dateRangeToQueryBounds,
	formatCredits,
	formatCurrencyAmount,
	formatOrderDate,
	formatPaymentDateTime,
	formatRawCents,
	resetCursorState,
} from "./table-utils";

describe("billing orders table utils", () => {
	it("builds billing orders paths with only non-empty query parameters", () => {
		expect(buildBillingOrdersQueryPath({ size: 20, sort: "desc" })).toBe(
			"/api/billing/orders?size=20&sort=desc"
		);
		expect(
			buildBillingOrdersQueryPath({
				status: "paid",
				cursor: "cursor-1",
				size: "50",
				sort: "asc",
			})
		).toBe("/api/billing/orders?status=paid&cursor=cursor-1&size=50&sort=asc");
		expect(
			buildBillingOrdersQueryPath({
				status: "pending",
				createdFrom: "2026-06-04T00:00:00.000Z",
				createdTo: "2026-06-04T23:59:59.999Z",
			})
		).toBe(
			"/api/billing/orders?status=pending&created_from=2026-06-04T00%3A00%3A00.000Z&created_to=2026-06-04T23%3A59%3A59.999Z"
		);
		expect(buildBillingOrdersQueryPath({ cursor: "", size: " 50 ", sort: "asc" })).toBe(
			"/api/billing/orders?size=50&sort=asc"
		);
		expect(buildBillingOrdersQueryPath()).toBe("/api/billing/orders");
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
		const paidAt = "2026-06-04T12:30:00Z";
		const expectedPaidAt = new Intl.DateTimeFormat("en-US", {
			year: "numeric",
			month: "short",
			day: "numeric",
			hour: "numeric",
			minute: "2-digit",
		}).format(new Date(paidAt));

		expect(formatOrderDate("2026-06-04T12:30:00Z")).toBe("2026-06-04");
		expect(formatOrderDate("bad")).toBe("Invalid date");
		expect(formatPaymentDateTime(paidAt)).toBe(expectedPaidAt);
		expect(formatPaymentDateTime()).toBe("-");
		expect(formatPaymentDateTime("bad")).toBe("Invalid date");
		expect(formatRawCents(1000)).toBe("1,000");
		expect(formatRawCents(Number.NaN)).toBe("0");
		expect(formatCurrencyAmount(1000, "EUR")).toBe("€10.00");
		expect(formatCurrencyAmount(1000, "not-a-currency")).toBe("€10.00");
		expect(formatCredits(5000)).toBe("5,000");
		expect(formatCredits(Number.NaN)).toBe("0");
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
