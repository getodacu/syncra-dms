import { json } from "@sveltejs/kit";

const DEFAULT_PUBLIC_SERVER_ERROR = "A server error occurred. Please try again.";

export function publicErrorMessage(
	status: number,
	message: string,
	fallback = DEFAULT_PUBLIC_SERVER_ERROR
) {
	if (status >= 500) return fallback;

	return message;
}

export function publicErrorStatus(status: number) {
	return status >= 500 ? 502 : status;
}

export function jsonPublicErrorResponse(status: number, message: string, fallback?: string) {
	return json(
		{ error: publicErrorMessage(status, message, fallback) },
		{ status: publicErrorStatus(status) }
	);
}
