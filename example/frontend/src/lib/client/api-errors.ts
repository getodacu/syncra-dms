export function publicApiErrorMessage(status: number, value: unknown, fallback: string) {
	if (status >= 500) return fallback;

	return apiErrorMessage(value, fallback);
}

export function apiErrorMessage(value: unknown, fallback: string) {
	if (isRecord(value) && typeof value.error === "string" && value.error.trim()) {
		return value.error;
	}

	return fallback;
}

function isRecord(value: unknown): value is Record<string, unknown> {
	return typeof value === "object" && value !== null && !Array.isArray(value);
}
