export type DatasetFieldTreeNode = {
	id: string;
	name: string;
	path: string;
	key: string;
	label: string;
	type: string;
	jsonCell: boolean;
	description?: string;
	children: DatasetFieldTreeNode[];
};

type JsonSchemaObject = Record<string, unknown>;

const MAX_DATASET_FIELD_KEY_CHARACTERS = 120;
const MAX_DATASET_FIELD_LABEL_CHARACTERS = 160;
const COLLISION_KEY_SEPARATOR = "__";

export function buildDatasetFieldTree(schema: unknown): DatasetFieldTreeNode[] {
	if (!isObjectSchema(schema)) return [];

	const properties = schema.properties;
	if (!isRecord(properties)) return [];

	const nodes = Object.entries(properties).map(([name, childSchema]) =>
		buildFieldTreeNode(name, [name], childSchema)
	);
	assignUniqueKeys(nodes);
	return nodes;
}

function buildFieldTreeNode(
	name: string,
	pathSegments: string[],
	schema: unknown
): DatasetFieldTreeNode {
	const type = schemaType(schema);
	const path = jsonPointerPath(pathSegments);
	const description =
		isRecord(schema) && typeof schema.description === "string"
			? { description: schema.description }
			: {};
	const children =
		type === "object" && isRecord(schema) && isRecord(schema.properties)
			? Object.entries(schema.properties).map(([childName, childSchema]) =>
					buildFieldTreeNode(childName, [...pathSegments, childName], childSchema)
				)
			: [];

	return {
		id: path,
		name,
		path,
		key: keyFromPathSegments(pathSegments),
		label: labelFromName(name),
		type,
		jsonCell: type === "object" || type === "array",
		...description,
		children,
	};
}

function isObjectSchema(value: unknown): value is JsonSchemaObject & { properties?: unknown } {
	if (!isRecord(value)) return false;

	const types = schemaTypes(value);
	if (types.length > 0) return types.includes("object");

	return isRecord(value.properties);
}

function schemaType(schema: unknown) {
	if (!isRecord(schema)) return "unknown";

	const types = schemaTypes(schema);
	if (types.includes("object")) return "object";
	if (types.includes("array")) return "array";
	const firstConcreteType = types.find((type) => type !== "null");
	if (firstConcreteType) return firstConcreteType;
	if (isRecord(schema.properties)) return "object";
	if ("items" in schema) return "array";

	return "unknown";
}

function schemaTypes(schema: JsonSchemaObject) {
	const { type } = schema;

	if (typeof type === "string" && type.trim()) return [type.trim()];
	if (Array.isArray(type)) {
		return type
			.filter((item): item is string => typeof item === "string" && item.trim() !== "")
			.map((item) => item.trim());
	}

	return [];
}

function jsonPointerPath(segments: string[]) {
	return `/${segments.map(jsonPointerEscape).join("/")}`;
}

function jsonPointerEscape(segment: string) {
	return segment.replace(/~/g, "~0").replace(/\//g, "~1");
}

function keyFromPathSegments(segments: string[]) {
	const key = segments.map(sanitizeKeySegment).join("_").replace(/_+/g, "_");
	return truncateKey(key || "field", MAX_DATASET_FIELD_KEY_CHARACTERS);
}

function sanitizeKeySegment(segment: string) {
	return segment.trim().replace(/[^A-Za-z0-9]+/g, "_").replace(/^_+|_+$/g, "") || "field";
}

function labelFromName(name: string) {
	const label = name.trim() ? name : "(empty)";
	return truncateText(label, MAX_DATASET_FIELD_LABEL_CHARACTERS);
}

function assignUniqueKeys(nodes: DatasetFieldTreeNode[]) {
	const flattenedNodes = flattenNodes(nodes);
	const nodesByBaseKey = new Map<string, DatasetFieldTreeNode[]>();

	for (const node of flattenedNodes) {
		const bucket = nodesByBaseKey.get(node.key);
		if (bucket) {
			bucket.push(node);
		} else {
			nodesByBaseKey.set(node.key, [node]);
		}
	}

	const usedKeys = new Set<string>();
	for (const [baseKey, bucket] of nodesByBaseKey) {
		if (bucket.length === 1) {
			bucket[0].key = baseKey;
			usedKeys.add(baseKey);
		}
	}

	for (const [baseKey, bucket] of nodesByBaseKey) {
		if (bucket.length === 1) continue;

		for (const node of bucket) {
			node.key = uniqueCollisionKey(baseKey, node.path, usedKeys);
		}
	}
}

function flattenNodes(nodes: DatasetFieldTreeNode[]): DatasetFieldTreeNode[] {
	return nodes.flatMap((node) => [node, ...flattenNodes(node.children)]);
}

function uniqueCollisionKey(baseKey: string, path: string, usedKeys: Set<string>) {
	const pathSuffix = hashJsonPointerPath(path);

	for (let attempt = 0; ; attempt += 1) {
		const suffix = attempt === 0 ? pathSuffix : `${pathSuffix}_${attempt}`;
		const key = keyWithSuffix(baseKey, suffix);

		if (!usedKeys.has(key)) {
			usedKeys.add(key);
			return key;
		}
	}
}

function keyWithSuffix(baseKey: string, suffix: string) {
	const suffixPart = `${COLLISION_KEY_SEPARATOR}${suffix}`;
	const baseLength = MAX_DATASET_FIELD_KEY_CHARACTERS - Array.from(suffixPart).length;
	const safeBaseKey = truncateKey(baseKey, Math.max(1, baseLength));

	return `${safeBaseKey}${suffixPart}`;
}

function hashJsonPointerPath(path: string) {
	let hash = 0x811c9dc5;

	for (const char of path) {
		hash ^= char.codePointAt(0) ?? 0;
		hash = Math.imul(hash, 0x01000193);
	}

	return (hash >>> 0).toString(36);
}

function truncateKey(key: string, maxLength: number) {
	const truncated = truncateText(key, maxLength);
	return truncated || "field";
}

function truncateText(value: string, maxLength: number) {
	return Array.from(value).slice(0, maxLength).join("");
}

function isRecord(value: unknown): value is Record<string, unknown> {
	return typeof value === "object" && value !== null && !Array.isArray(value);
}
