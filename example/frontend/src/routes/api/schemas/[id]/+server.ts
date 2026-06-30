import { json } from "@sveltejs/kit";

import { deleteSchema, getSchema, isSchemaApiError, updateSchema } from "$lib/server/schemas";
import { jsonPublicErrorResponse } from "$lib/server/public-errors";
import type { RequestHandler } from "./$types";
import { readSchemaRequest } from "../schema-request";

function schemaErrorResponse(error: unknown) {
	if (isSchemaApiError(error)) {
		return jsonPublicErrorResponse(error.status, error.message);
	}

	throw error;
}

export const GET: RequestHandler = async ({ fetch, locals, params }) => {
	if (!locals.user) {
		return json({ error: "authentication required" }, { status: 401 });
	}

	try {
		const schema = await getSchema(fetch, params.id, { userId: locals.user.id });
		return json(schema);
	} catch (error) {
		return schemaErrorResponse(error);
	}
};

export const PUT: RequestHandler = async ({ request, fetch, locals, params }) => {
	if (!locals.user) {
		return json({ error: "authentication required" }, { status: 401 });
	}

	const bodyResult = await readSchemaRequest(request);
	if (!bodyResult.ok) {
		return json({ error: bodyResult.error }, { status: 400 });
	}

	try {
		const schema = await updateSchema(fetch, params.id, bodyResult.input, {
			userId: locals.user.id
		});
		return json(schema);
	} catch (error) {
		return schemaErrorResponse(error);
	}
};

export const DELETE: RequestHandler = async ({ fetch, locals, params }) => {
	if (!locals.user) {
		return json({ error: "authentication required" }, { status: 401 });
	}

	try {
		const deleted = await deleteSchema(fetch, params.id, { userId: locals.user.id });
		return json(deleted);
	} catch (error) {
		return schemaErrorResponse(error);
	}
};
