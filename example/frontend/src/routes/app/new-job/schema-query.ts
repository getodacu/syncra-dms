export const NEW_JOB_SCHEMA_QUERY_PARAM = "schema_id";

type SchemaIdItem = {
	id: string;
};

export function buildNewJobPath(schemaId: string) {
	const params = new URLSearchParams();
	const normalizedSchemaId = schemaId.trim();

	if (normalizedSchemaId) {
		params.set(NEW_JOB_SCHEMA_QUERY_PARAM, normalizedSchemaId);
	}

	const query = params.toString();
	return query ? `/app/new-job?${query}` : "/app/new-job";
}

export function schemaIdFromSearchParams(params: URLSearchParams) {
	return params.get(NEW_JOB_SCHEMA_QUERY_PARAM)?.trim() ?? "";
}

export function isSelectedSchemaIdUnavailable(input: {
	selectedSchemaId: string;
	schemas: SchemaIdItem[];
	schemasLoaded: boolean;
}) {
	const schemaId = input.selectedSchemaId.trim();
	if (!input.schemasLoaded || !schemaId) return false;

	return !input.schemas.some((schema) => schema.id === schemaId);
}
