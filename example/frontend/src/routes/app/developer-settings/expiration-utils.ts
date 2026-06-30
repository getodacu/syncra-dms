import { CalendarDateTime, getLocalTimeZone, type DateValue } from "@internationalized/date";

export type APIKeyExpirationPreset = "week" | "month" | "quarter";

export function apiKeyExpirationPresetDate(preset: APIKeyExpirationPreset, baseDate: DateValue) {
	switch (preset) {
		case "week":
			return baseDate.add({ weeks: 1 });
		case "month":
			return baseDate.add({ months: 1 });
		case "quarter":
			return baseDate.add({ months: 3 });
	}
}

export function apiKeyExpirationDateToISOString(
	date: DateValue,
	timeZone = getLocalTimeZone()
) {
	return new CalendarDateTime(
		date.calendar,
		date.era,
		date.year,
		date.month,
		date.day,
		23,
		59,
		59,
		999
	)
		.toDate(timeZone)
		.toISOString();
}

export function isSameCalendarDate(left: DateValue | undefined, right: DateValue | undefined) {
	return Boolean(left && right && left.toString() === right.toString());
}
