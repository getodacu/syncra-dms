import type { CreateSchemaInput } from "$lib/server/schemas";

const MAX_SCHEMA_REQUEST_BYTES = 1 << 20;
const MAX_SCHEMA_JSON_BYTES = 1 << 20;
const textEncoder = new TextEncoder();

export type SchemaRequestResult =
	| { ok: true; input: CreateSchemaInput }
	| { ok: false; error: string };

export type JsonObjectRequestResult =
	| { ok: true; body: Record<string, unknown> }
	| { ok: false; error: string };

export function isJsonObject(value: unknown): value is Record<string, unknown> {
	return typeof value === "object" && value !== null && !Array.isArray(value);
}

function bodySizeExceedsLimit(request: Request) {
	const contentLength = request.headers.get("content-length");
	if (!contentLength) return false;

	const size = Number(contentLength);
	return Number.isFinite(size) && size > MAX_SCHEMA_REQUEST_BYTES;
}

async function readLimitedBody(request: Request) {
	if (bodySizeExceedsLimit(request)) return null;
	if (!request.body) return "";

	const reader = request.body.getReader();
	const textDecoder = new TextDecoder();
	let bytes = 0;
	let body = "";

	try {
		while (true) {
			const { done, value } = await reader.read();
			if (done) break;

			bytes += value.byteLength;
			if (bytes > MAX_SCHEMA_REQUEST_BYTES) {
				await reader.cancel();
				return null;
			}
			body += textDecoder.decode(value, { stream: true });
		}

		body += textDecoder.decode();
		return body;
	} finally {
		reader.releaseLock();
	}
}

function jsonByteLength(value: unknown) {
	return textEncoder.encode(JSON.stringify(value)).byteLength;
}

export async function readJsonObjectRequest(request: Request): Promise<JsonObjectRequestResult> {
	let body: unknown;
	try {
		const text = await readLimitedBody(request);
		if (text === null) {
			return { ok: false, error: "request body too large" };
		}
		body = JSON.parse(text);
	} catch {
		return { ok: false, error: "invalid JSON body" };
	}

	if (!isJsonObject(body)) {
		return { ok: false, error: "invalid JSON body" };
	}

	return { ok: true, body };
}

export async function readSchemaRequest(request: Request): Promise<SchemaRequestResult> {
	const bodyResult = await readJsonObjectRequest(request);
	if (!bodyResult.ok) return bodyResult;

	const { body } = bodyResult;
	const name = typeof body.name === "string" ? body.name.trim() : "";
	if (!name) {
		return { ok: false, error: "name is required" };
	}

	if (Array.from(name).length > 160) {
		return { ok: false, error: "name must be at most 160 characters" };
	}

	if (!isJsonObject(body.schema)) {
		return { ok: false, error: "schema must be a JSON object" };
	}

	if (jsonByteLength(body.schema) > MAX_SCHEMA_JSON_BYTES) {
		return { ok: false, error: "schema is too large" };
	}

	if ("description" in body && typeof body.description !== "string") {
		return { ok: false, error: "description must be a string" };
	}

	if ("strict" in body && typeof body.strict !== "boolean") {
		return { ok: false, error: "strict must be a boolean" };
	}

	return {
		ok: true,
		input: {
			name,
			description: typeof body.description === "string" ? body.description : "",
			strict: typeof body.strict === "boolean" ? body.strict : true,
			schema: body.schema
		}
	};
}
