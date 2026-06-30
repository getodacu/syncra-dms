export function publicErrorStatus(status: number) {
	if (status >= 500) return 502;
	return status;
}

export function publicErrorMessage(status: number, message: string, fallback: string) {
	if (status >= 500) return fallback;
	return message || fallback;
}
