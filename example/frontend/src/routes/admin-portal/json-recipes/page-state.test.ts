import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";

const source = (path: string) => readFileSync(new URL(path, import.meta.url), "utf8");

describe("admin JSON recipes pages", () => {
	it("renders a recipe list with pagination, empty, loading, and error states", () => {
		const page = source("./+page.svelte");

		expect(page).toContain('import { m } from "$lib/paraglide/messages.js";');
		expect(page).toContain("fetchAdminJSONRecipes");
		expect(page).toContain("json_recipes_title");
		expect(page).toContain("json_recipes_new_recipe");
		expect(page).toContain("json_recipes_manage_categories");
		expect(page).toContain("fetchAdminJSONRecipeCategories");
		expect(page).toContain("json_recipes_no_recipes_found");
		expect(page).toContain("json_recipes_loading");
		expect(page).toContain("cursorNextState");
		expect(page).toContain("cursorPreviousState");
		expect(page).toContain("counter");
		expect(page).toContain("categoryLabel");
		expect(page).toContain('import * as Dialog from "$lib/components/ui/dialog/index.js";');
		expect(page).toContain("<Dialog.Root");
		expect(page).not.toContain("window.confirm");
	});

	it("uses the visual JSON schema builder in the admin recipe editor", () => {
		const editor = source("./json-recipe-editor.svelte");

		expect(editor).toContain("JsonSchemaBuilder");
		expect(editor).toContain("json_recipes_editor_badge");
		expect(editor).toContain("json_recipes_title_label");
		expect(editor).toContain("json_recipes_description_label");
		expect(editor).toContain("json_recipes_category_label");
		expect(editor).toContain("json_recipes_others");
		expect(editor).toContain("json_recipes_validation_title_required");
		expect(editor).toContain("json_recipes_validation_json_object");
		expect(editor).not.toContain("strict");
	});

	it("covers create, update, delete, loading, and error states across editor routes", () => {
		const newPage = source("./new/+page.svelte");
		const detailPage = source("./[id]/+page.svelte");

		expect(newPage).toContain("createAdminJSONRecipe");
		expect(newPage).toContain("fetchAdminJSONRecipeCategories");
		expect(newPage).toContain("json_recipes_created_success");
		expect(detailPage).toContain("fetchAdminJSONRecipe");
		expect(detailPage).toContain("fetchAdminJSONRecipeCategories");
		expect(detailPage).toContain("updateAdminJSONRecipe");
		expect(detailPage).toContain("deleteAdminJSONRecipe");
		expect(detailPage).toContain("json_recipes_delete_confirm");
		expect(detailPage).toContain('import * as Dialog from "$lib/components/ui/dialog/index.js";');
		expect(detailPage).toContain("<Dialog.Root");
		expect(detailPage).not.toContain("window.confirm");
		expect(detailPage).toContain("json_recipes_could_not_load");
		expect(detailPage).toContain("json_recipes_not_found_title");
	});

	it("provides a dedicated admin category management subpage", () => {
		const page = source("./categories/+page.svelte");

		expect(page).toContain("fetchAdminJSONRecipeCategories");
		expect(page).toContain("createAdminJSONRecipeCategory");
		expect(page).toContain("updateAdminJSONRecipeCategory");
		expect(page).toContain("deleteAdminJSONRecipeCategory");
		expect(page).toContain("json_recipe_categories_title");
		expect(page).toContain("json_recipe_categories_title_en_label");
		expect(page).toContain("json_recipe_categories_title_ro_label");
		expect(page).toContain("<Dialog.Root");
	});
});
