import type { CreateDatasetInput, DatasetField, UpdateDatasetInput } from "$lib/server/datasets";
import { readJsonObjectRequest } from "../schemas/schema-request";

type CreateDatasetRequestResult =
	| { ok: true; input: CreateDatasetInput }
	| { ok: false; error: string };

type UpdateDatasetRequestResult =
	| { ok: true; input: UpdateDatasetInput }
	| { ok: false; error: string };

type ParsedDatasetRequestResult =
	| { ok: true; name: string; schemaId?: string; selectedFields: DatasetField[] }
	| { ok: false; error: string };

const MAX_DATASET_NAME_CHARACTERS = 160;

function normalizedString(value: unknown) {
	return typeof value === "string" ? value.trim() : "";
}

function parseDatasetFields(value: unknown) {
	if (!Array.isArray(value)) {
		return {
			ok: false,
			error: "selected_fields must be an array of dataset fields"
		} as const;
	}

	if (value.length === 0) {
		return { ok: false, error: "selected_fields is required" } as const;
	}

	const selectedFields: DatasetField[] = [];
	for (const field of value) {
		if (typeof field !== "object" || field === null || Array.isArray(field)) {
			return {
				ok: false,
				error: "selected_fields must contain path, key, and label strings"
			} as const;
		}

		const candidate = field as Record<string, unknown>;
		if (
			typeof candidate.path !== "string" ||
			typeof candidate.key !== "string" ||
			typeof candidate.label !== "string"
		) {
			return {
				ok: false,
				error: "selected_fields must contain path, key, and label strings"
			} as const;
		}

		if (!candidate.path.trim() || !candidate.key.trim() || !candidate.label.trim()) {
			return {
				ok: false,
				error: "selected_fields must contain non-empty path, key, and label strings"
			} as const;
		}

		selectedFields.push({
			path: candidate.path,
			key: candidate.key,
			label: candidate.label
		});
	}

	return { ok: true, selectedFields } as const;
}

async function readDatasetRequest(request: Request): Promise<ParsedDatasetRequestResult> {
	const bodyResult = await readJsonObjectRequest(request);
	if (!bodyResult.ok) return bodyResult;

	const { body } = bodyResult;
	const name = normalizedString(body.name);
	if (!name) {
		return { ok: false, error: "name is required" };
	}

	if (Array.from(name).length > MAX_DATASET_NAME_CHARACTERS) {
		return { ok: false, error: "name must be at most 160 characters" };
	}

	const fieldsResult = parseDatasetFields(body.selected_fields);
	if (!fieldsResult.ok) return fieldsResult;

	if ("schema_id" in body && typeof body.schema_id !== "string") {
		return { ok: false, error: "schema_id must be a string" };
	}

	const schemaId = normalizedString(body.schema_id);
	return {
		ok: true,
		name,
		...(schemaId ? { schemaId } : {}),
		selectedFields: fieldsResult.selectedFields
	};
}

export async function readCreateDatasetRequest(
	request: Request
): Promise<CreateDatasetRequestResult> {
	const result = await readDatasetRequest(request);
	if (!result.ok) return result;

	if (!result.schemaId) {
		return { ok: false, error: "schema_id is required" };
	}

	return {
		ok: true,
		input: {
			name: result.name,
			schema_id: result.schemaId,
			selected_fields: result.selectedFields
		}
	};
}

export async function readUpdateDatasetRequest(
	request: Request
): Promise<UpdateDatasetRequestResult> {
	const result = await readDatasetRequest(request);
	if (!result.ok) return result;

	if (!result.schemaId) {
		return { ok: false, error: "schema_id is required" };
	}

	return {
		ok: true,
		input: {
			name: result.name,
			schema_id: result.schemaId,
			selected_fields: result.selectedFields
		}
	};
}
