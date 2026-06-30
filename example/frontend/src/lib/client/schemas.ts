import { publicApiErrorMessage } from "$lib/client/api-errors";

export type PersonalSchemaOption = {
	id: string;
	name: string;
	description: string;
};

type ClientFetch = (input: RequestInfo | URL, init?: RequestInit) => Promise<Response>;
type SchemaOptionSource = {
	id: string;
	name: string;
	description: string;
};
type SchemaListPayload = {
	schemas: PersonalSchemaOption[];
	next_cursor: string | null;
};

export const PERSONAL_SCHEMA_OPTIONS_LIMIT = 100;
export const PERSONAL_SCHEMA_OPTIONS_QUERY_KEY = [
	"schemas",
	"mine",
	"options",
	{ size: PERSONAL_SCHEMA_OPTIONS_LIMIT },
] as const;

export async function fetchPersonalSchemaOptions(
	fetchFn: ClientFetch
): Promise<PersonalSchemaOption[]> {
	const response = await fetchFn(
		`/api/schemas?scope=mine&size=${PERSONAL_SCHEMA_OPTIONS_LIMIT}`
	);
	const json = await readResponseJSON(response);

	if (!response.ok) {
		throw new Error(publicApiErrorMessage(response.status, json, "Failed to load schemas"));
	}
	if (!isSchemaListPayload(json)) {
		throw new Error("Invalid schema response");
	}

	return json.schemas.map(toPersonalSchemaOption);
}

export function upsertPersonalSchemaOption(
	current: PersonalSchemaOption[] | undefined,
	schema: SchemaOptionSource
) {
	const option = toPersonalSchemaOption(schema);
	const options = current ?? [];
	const existingIndex = options.findIndex((item) => item.id === option.id);

	if (existingIndex === -1) return [option, ...options];

	const next = [...options];
	next[existingIndex] = option;
	return next;
}

export function removePersonalSchemaOptions(
	current: PersonalSchemaOption[] | undefined,
	ids: string[]
) {
	if (!current) return current;

	const deletedIds = new Set(ids);
	return current.filter((option) => !deletedIds.has(option.id));
}

async function readResponseJSON(response: Response): Promise<unknown> {
	const text = await response.text();
	if (!text.trim()) return null;

	try {
		return JSON.parse(text);
	} catch {
		return null;
	}
}

function isSchemaListPayload(value: unknown): value is SchemaListPayload {
	return (
		isRecord(value) &&
		Array.isArray(value.schemas) &&
		value.schemas.every(isPersonalSchemaOption) &&
		(typeof value.next_cursor === "string" || value.next_cursor === null)
	);
}

function isPersonalSchemaOption(value: unknown): value is PersonalSchemaOption {
	return (
		isRecord(value) &&
		typeof value.id === "string" &&
		typeof value.name === "string" &&
		typeof value.description === "string"
	);
}

function toPersonalSchemaOption(schema: SchemaOptionSource): PersonalSchemaOption {
	return {
		id: schema.id,
		name: schema.name,
		description: schema.description,
	};
}

function isRecord(value: unknown): value is Record<string, unknown> {
	return typeof value === "object" && value !== null && !Array.isArray(value);
}
