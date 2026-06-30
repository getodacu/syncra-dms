export type JsonSchemaObject = Record<string, unknown>;

export type SchemaTreeNode = {
	id: string;
	name: string;
	path: string;
	type: string;
	required: boolean;
	description?: string;
	details: string[];
	children: SchemaTreeNode[];
};

function isRecord(value: unknown): value is Record<string, unknown> {
	return typeof value === "object" && value !== null && !Array.isArray(value);
}

function stringArray(value: unknown) {
	return Array.isArray(value) ? value.filter((item): item is string => typeof item === "string") : [];
}

function pointerSegment(value: string) {
	return value.replaceAll("~", "~0").replaceAll("/", "~1");
}

function propertyPath(parentPath: string, propertyName: string) {
	return `${parentPath}/properties/${pointerSegment(propertyName)}`;
}

function typeLabel(schema: Record<string, unknown>) {
	const type = schema.type;
	if (typeof type === "string") return type;
	if (Array.isArray(type)) return stringArray(type).join(" | ") || "unknown";
	if ("properties" in schema) return "object";
	if ("items" in schema) return "array";
	return "unknown";
}

function detailList(schema: Record<string, unknown>) {
	const details: string[] = [];
	if (typeof schema.format === "string") details.push(`format: ${schema.format}`);
	if (Array.isArray(schema.enum)) {
		const values = schema.enum
			.filter((value) => value === null || ["string", "number", "boolean"].includes(typeof value))
			.map(String)
			.slice(0, 6);
		if (values.length) details.push(`enum: ${values.join(", ")}`);
	}
	if (typeof schema.minimum === "number") details.push(`min: ${schema.minimum}`);
	if (typeof schema.maximum === "number") details.push(`max: ${schema.maximum}`);
	if (typeof schema.minLength === "number") details.push(`min length: ${schema.minLength}`);
	if (typeof schema.maxLength === "number") details.push(`max length: ${schema.maxLength}`);
	return details;
}

function nodeFor(name: string, path: string, schema: Record<string, unknown>, required: boolean) {
	const node: SchemaTreeNode = {
		id: path,
		name,
		path,
		type: typeLabel(schema),
		required,
		description: typeof schema.description === "string" ? schema.description : undefined,
		details: detailList(schema),
		children: []
	};

	const requiredChildren = new Set(stringArray(schema.required));
	if (isRecord(schema.properties)) {
		node.children = Object.entries(schema.properties)
			.filter((entry): entry is [string, Record<string, unknown>] => isRecord(entry[1]))
			.map(([childName, childSchema]) =>
				nodeFor(childName, propertyPath(path, childName), childSchema, requiredChildren.has(childName))
			);
	}

	if (isRecord(schema.items)) {
		node.children = [nodeFor("items", `${path}/items`, schema.items, false), ...node.children];
	}

	return node;
}

export function buildSchemaTree(schema: JsonSchemaObject) {
	const requiredRoot = new Set(stringArray(schema.required));
	if (isRecord(schema.properties)) {
		const nodes = Object.entries(schema.properties)
			.filter((entry): entry is [string, Record<string, unknown>] => isRecord(entry[1]))
			.map(([name, childSchema]) => nodeFor(name, propertyPath("", name), childSchema, requiredRoot.has(name)));

		if (nodes.length) return nodes;
	}

	return [nodeFor("schema", "schema", schema, false)];
}
