import { json } from "@sveltejs/kit";

import {
	createSchema,
	deleteSchemas,
	isSchemaApiError,
	listSchemas,
	listSchemasPage
} from "$lib/server/schemas";
import { jsonPublicErrorResponse } from "$lib/server/public-errors";
import type { RequestHandler } from "./$types";
import { readJsonObjectRequest, readSchemaRequest } from "./schema-request";

function schemaErrorResponse(error: unknown) {
	if (isSchemaApiError(error)) {
		return jsonPublicErrorResponse(error.status, error.message);
	}

	throw error;
}

export const GET: RequestHandler = async ({ url, fetch, locals }) => {
	if (!locals.user) {
		return json({ error: "authentication required" }, { status: 401 });
	}

	const scope = url.searchParams.get("scope") ?? "system";
	if (scope !== "system" && scope !== "mine") {
		return json({ error: "invalid schema scope" }, { status: 400 });
	}

	try {
		const userOptions = scope === "mine" ? { userId: locals.user.id } : {};
		const wantsPage =
			url.searchParams.has("cursor") ||
			url.searchParams.has("size") ||
			url.searchParams.has("sort");

		if (wantsPage) {
			const page = await listSchemasPage(fetch, {
				...userOptions,
				...(url.searchParams.has("cursor")
					? { cursor: url.searchParams.get("cursor") ?? "" }
					: {}),
				...(url.searchParams.has("size") ? { size: url.searchParams.get("size") ?? "" } : {}),
				...(url.searchParams.has("sort") ? { sort: url.searchParams.get("sort") ?? "" } : {})
			});
			return json(page);
		}

		const schemas =
			scope === "mine" ? await listSchemas(fetch, userOptions) : await listSchemas(fetch);
		return json(schemas);
	} catch (error) {
		return schemaErrorResponse(error);
	}
};

export const POST: RequestHandler = async ({ request, fetch, locals }) => {
	if (!locals.user) {
		return json({ error: "authentication required" }, { status: 401 });
	}

	const bodyResult = await readSchemaRequest(request);
	if (!bodyResult.ok) {
		return json({ error: bodyResult.error }, { status: 400 });
	}

	try {
		const saved = await createSchema(
			fetch,
			bodyResult.input,
			{
				userId: locals.user.id
			}
		);

		return json(saved, { status: 201 });
	} catch (error) {
		return schemaErrorResponse(error);
	}
};

export const DELETE: RequestHandler = async ({ request, fetch, locals }) => {
	if (!locals.user) {
		return json({ error: "authentication required" }, { status: 401 });
	}

	const bodyResult = await readJsonObjectRequest(request);
	if (!bodyResult.ok) {
		return json({ error: bodyResult.error }, { status: 400 });
	}

	const { ids } = bodyResult.body;
	if (!Array.isArray(ids) || !ids.every((id) => typeof id === "string")) {
		return json({ error: "ids must be an array of strings" }, { status: 400 });
	}

	try {
		const deleted = await deleteSchemas(fetch, ids, { userId: locals.user.id });
		return json(deleted);
	} catch (error) {
		return schemaErrorResponse(error);
	}
};
