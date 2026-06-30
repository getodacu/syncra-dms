import type { JsonSchemaObject } from "../api";

export function schemaHasFields(schema: JsonSchemaObject) {
	const properties = schema.properties;

	return isRecord(properties) && Object.keys(properties).length > 0;
}

function isRecord(value: unknown): value is Record<string, unknown> {
	return typeof value === "object" && value !== null && !Array.isArray(value);
}
