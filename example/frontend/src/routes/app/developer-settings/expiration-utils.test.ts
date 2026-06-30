import { CalendarDate } from "@internationalized/date";
import { describe, expect, it } from "vitest";

import {
	apiKeyExpirationDateToISOString,
	apiKeyExpirationPresetDate,
	isSameCalendarDate
} from "./expiration-utils";

describe("API key expiration utilities", () => {
	it("calculates supported expiration preset dates from a base date", () => {
		const baseDate = new CalendarDate(2026, 6, 9);

		expect(apiKeyExpirationPresetDate("week", baseDate).toString()).toBe("2026-06-16");
		expect(apiKeyExpirationPresetDate("month", baseDate).toString()).toBe("2026-07-09");
		expect(apiKeyExpirationPresetDate("quarter", baseDate).toString()).toBe("2026-09-09");
	});

	it("serializes selected expiration dates at the end of the local day", () => {
		const date = new CalendarDate(2026, 6, 9);

		expect(apiKeyExpirationDateToISOString(date, "UTC")).toBe("2026-06-09T23:59:59.999Z");
		expect(apiKeyExpirationDateToISOString(date, "Europe/Bucharest")).toBe(
			"2026-06-09T20:59:59.999Z"
		);
	});

	it("compares calendar dates without comparing object identity", () => {
		expect(isSameCalendarDate(new CalendarDate(2026, 6, 9), new CalendarDate(2026, 6, 9))).toBe(
			true
		);
		expect(isSameCalendarDate(new CalendarDate(2026, 6, 9), new CalendarDate(2026, 6, 10))).toBe(
			false
		);
		expect(isSameCalendarDate(undefined, new CalendarDate(2026, 6, 9))).toBe(false);
	});
});
