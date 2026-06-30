import { json } from "@sveltejs/kit";

import {
	isAuthApiError,
	isSupportedPreferredLanguage,
	setPreferredLanguageCookie,
	updateAuthUser
} from "$lib/server/auth";
import { jsonPublicErrorResponse } from "$lib/server/public-errors";
import type { UpdateAuthUserInput } from "$lib/server/auth";
import type { RequestHandler } from "./$types";

const MAX_AUTH_USER_AVATAR_BYTES = 5 << 20;
const MAX_AUTH_USER_REQUEST_BYTES = Math.ceil(MAX_AUTH_USER_AVATAR_BYTES / 3) * 4 + 4096;

export const PATCH: RequestHandler = async ({ request, fetch, locals, cookies }) => {
	if (!locals.user) {
		return json({ error: "authentication required" }, { status: 401 });
	}

	let body: unknown;
	try {
		const text = await readLimitedBody(request);
		if (text === null) {
			return json({ error: "request body too large" }, { status: 400 });
		}
		body = JSON.parse(text);
	} catch {
		return json({ error: "invalid JSON body" }, { status: 400 });
	}

	if (!isUpdateAuthUserInput(body)) {
		return json({ error: "invalid user update payload" }, { status: 400 });
	}

	try {
		const user = await updateAuthUser(fetch, request.headers.get("cookie"), body);
		setPreferredLanguageCookie(cookies, user.preferredLanguage);
		return json(user);
	} catch (error) {
		if (isAuthApiError(error)) {
			return jsonPublicErrorResponse(error.status, error.message);
		}
		throw error;
	}
};

function isUpdateAuthUserInput(value: unknown): value is UpdateAuthUserInput {
	if (typeof value !== "object" || value === null || Array.isArray(value)) return false;
	if (!("preferredLanguage" in value)) return true;
	return isSupportedPreferredLanguage(value.preferredLanguage);
}

function bodySizeExceedsLimit(request: Request) {
	const contentLength = request.headers.get("content-length");
	if (!contentLength) return false;

	const size = Number(contentLength);
	return Number.isFinite(size) && size > MAX_AUTH_USER_REQUEST_BYTES;
}

async function readLimitedBody(request: Request) {
	if (bodySizeExceedsLimit(request)) return null;
	if (!request.body) return "";

	const reader = request.body.getReader();
	const decoder = new TextDecoder();
	let bytes = 0;
	let body = "";

	try {
		while (true) {
			const { done, value } = await reader.read();
			if (done) break;

			bytes += value.byteLength;
			if (bytes > MAX_AUTH_USER_REQUEST_BYTES) {
				await reader.cancel();
				return null;
			}
			body += decoder.decode(value, { stream: true });
		}

		body += decoder.decode();
		return body;
	} finally {
		reader.releaseLock();
	}
}
