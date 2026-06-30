import type { OCRRecipeResponse } from "./api";

export type RecipeFieldSummary = {
	key: string;
	label: string;
	type: string;
	description: string;
	required: boolean;
};

export type RecipeSummary = {
	fieldCount: number;
	requiredCount: number;
	fields: RecipeFieldSummary[];
	searchText: string;
	prettyJson: string;
};

export function summarizeRecipe(recipe: OCRRecipeResponse): RecipeSummary {
	const properties = objectValue(recipe.json.properties);
	const required = stringSet(recipe.json.required);
	const fields = Object.entries(properties)
		.map(([key, schema]) => summarizeField(key, schema, required.has(key)))
		.sort((left, right) => Number(right.required) - Number(left.required) || left.label.localeCompare(right.label));

	const fieldSearchText = fields
		.map((field) => `${field.key} ${field.label} ${field.type} ${field.description}`)
		.join(" ");
	const categorySearchText = recipe.category
		? `${recipe.category.title.en} ${recipe.category.title.ro}`
		: "";

	return {
		fieldCount: fields.length,
		requiredCount: fields.filter((field) => field.required).length,
		fields,
		searchText: normalizeSearchText(`${recipe.title} ${recipe.description} ${categorySearchText} ${fieldSearchText}`),
		prettyJson: JSON.stringify(recipe.json, null, 2)
	};
}

export function normalizeSearchText(value: string) {
	return value
		.toLocaleLowerCase()
		.normalize("NFD")
		.replace(/\p{Diacritic}/gu, "")
		.replace(/[_-]+/g, " ")
		.replace(/\s+/g, " ")
		.trim();
}

function summarizeField(key: string, value: unknown, required: boolean): RecipeFieldSummary {
	const schema = objectValue(value);
	const description = typeof schema.description === "string" ? schema.description : "";
	const label = key.replace(/[_-]+/g, " ");

	return {
		key,
		label,
		type: fieldType(schema),
		description,
		required
	};
}

function fieldType(schema: Record<string, unknown>) {
	const type = schema.type;
	if (typeof type === "string") {
		if (type === "array") {
			const items = objectValue(schema.items);
			const itemType = typeof items.type === "string" ? items.type : "item";
			return `${itemType}[]`;
		}
		return type;
	}
	if (Array.isArray(type)) {
		return type.filter((item) => typeof item === "string").join(" | ") || "unknown";
	}
	if (Array.isArray(schema.enum)) return "enum";
	if ("properties" in schema) return "object";
	return "unknown";
}

function objectValue(value: unknown): Record<string, unknown> {
	return typeof value === "object" && value !== null && !Array.isArray(value)
		? (value as Record<string, unknown>)
		: {};
}

function stringSet(value: unknown) {
	if (!Array.isArray(value)) return new Set<string>();
	return new Set(value.filter((item): item is string => typeof item === "string"));
}
