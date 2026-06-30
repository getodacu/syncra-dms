import type { CollectionInput } from "$lib/server/collections";
import { readJsonObjectRequest } from "../schemas/schema-request";

type CollectionRequestResult =
	| { ok: true; input: CollectionInput }
	| { ok: false; error: string };

export async function readCollectionRequest(
	request: Request
): Promise<CollectionRequestResult> {
	const bodyResult = await readJsonObjectRequest(request);
	if (!bodyResult.ok) return bodyResult;

	const { body } = bodyResult;
	const name = typeof body.name === "string" ? body.name.trim() : "";
	if (!name) {
		return { ok: false, error: "name is required" };
	}

	if (
		!Array.isArray(body.schema_ids) ||
		!body.schema_ids.every((schemaId) => typeof schemaId === "string")
	) {
		return { ok: false, error: "schema_ids must be an array of strings" };
	}

	return {
		ok: true,
		input: {
			name,
			schema_ids: body.schema_ids
		}
	};
}
