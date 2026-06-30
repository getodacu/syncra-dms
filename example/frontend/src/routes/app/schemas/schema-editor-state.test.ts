import { existsSync, readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";

const editorSource = () => readFileSync(new URL("./schema-editor.svelte", import.meta.url), "utf8");
const editPageSource = () =>
	readFileSync(new URL("./edit/[id]/+page.svelte", import.meta.url), "utf8");
const listPageSource = () => readFileSync(new URL("./+page.svelte", import.meta.url), "utf8");
const actionsCellSource = () =>
	readFileSync(new URL("./actions-cell.svelte", import.meta.url), "utf8");
const schemaIdCopyModuleUrl = new URL("./schema-id-copy.svelte", import.meta.url);
const schemaIdCopySource = () =>
	existsSync(schemaIdCopyModuleUrl) ? readFileSync(schemaIdCopyModuleUrl, "utf8") : "";
const newPageSource = () => readFileSync(new URL("./new/+page.svelte", import.meta.url), "utf8");
const newJobPageSource = () =>
	readFileSync(new URL("../new-job/+page.svelte", import.meta.url), "utf8");
const appCssSource = () => readFileSync(new URL("../../../app.css", import.meta.url), "utf8");
const pnpmWorkspaceSource = () =>
	readFileSync(new URL("../../../../pnpm-workspace.yaml", import.meta.url), "utf8");
const pnpmLockSource = () =>
	readFileSync(new URL("../../../../pnpm-lock.yaml", import.meta.url), "utf8");
const jsonJoyPatchSource = () =>
	readFileSync(new URL("../../../../patches/jsonjoy-builder@1.0.3.patch", import.meta.url), "utf8");
const componentSource = (path: string) =>
	readFileSync(new URL(`../../../lib/components/${path}`, import.meta.url), "utf8");
const jsonJoyLocaleModuleUrl = new URL(
	"../../../lib/components/json-schema-builder-locale.ts",
	import.meta.url
);
const jsonJoyLocaleModuleSource = () =>
	existsSync(jsonJoyLocaleModuleUrl) ? readFileSync(jsonJoyLocaleModuleUrl, "utf8") : "";
const jsonJoyEnglishLocaleSource = () =>
	readFileSync(
		new URL("../../../../node_modules/jsonjoy-builder/dist/i18n/locales/en.js", import.meta.url),
		"utf8"
	);
const translationObjectSource = (source: string, objectName?: string) => {
	if (!objectName) return source;
	return new RegExp(`${objectName}\\s*=\\s*\\{([\\s\\S]*?)\\}\\s*as const`).exec(source)?.[1] ?? "";
};
const translationKeysFromSource = (source: string, objectName?: string) =>
	Array.from(
		translationObjectSource(source, objectName).matchAll(/^\s*([a-zA-Z][A-Za-z0-9]*):\s*"/gm),
		(match) => match[1]
	);
const translationValueFromSource = (source: string, key: string, objectName?: string) => {
	const match = new RegExp(`^\\s*${key}:\\s*"((?:[^"\\\\]|\\\\.)*)"`, "m").exec(
		translationObjectSource(source, objectName)
	);
	return match ? JSON.parse(`"${match[1]}"`) as string : undefined;
};
const placeholderTokens = (value: string | undefined) => value?.match(/\{[^}]+}/g) ?? [];

describe("schema editor state", () => {
	it("keeps the schema as raw state so JSONJoy can structuredClone it", () => {
		const source = editorSource();

		expect(source).toContain("let schema = $state.raw<JsonSchemaValue>");
		expect(source).not.toContain("let schema = $state<JsonSchemaValue>");
	});

	it("syncs JSONJoy and Monaco themes from mode-watcher", () => {
		const source = componentSource("json-schema-builder.svelte");

		expect(source).toContain('import { mode } from "mode-watcher";');
		expect(source).toContain('const jsonJoyTheme = $derived(mode.current === "dark" ? "dark" : "light");');
		expect(source).toContain('className: schemaBuilderClassName');
		expect(source).toContain("jsonJoyTheme;");
		expect(source).toContain('import("monaco-editor")');
		expect(source).toContain("monaco.editor.setTheme(monacoThemeName);");
		expect(source).toContain("new MutationObserver");
		expect(source).toContain("stopThemeSync();");
	});

	it("limits JSONJoy visual field type controls to extraction-friendly types", () => {
		const builder = componentSource("json-schema-builder.svelte");
		const locale = jsonJoyLocaleModuleSource();
		const workspace = pnpmWorkspaceSource();
		const lockfile = pnpmLockSource();
		const patch = jsonJoyPatchSource();

		expect(locale).toContain('schemaTypeObject: "Group"');
		expect(locale).toContain("englishJsonJoyBuilderMessages");
		expect(builder).toContain("messages: jsonJoyBuilderConfig.messages");
		expect(workspace).toContain("patchedDependencies:");
		expect(workspace).toContain("jsonjoy-builder@1.0.3: patches/jsonjoy-builder@1.0.3.patch");
		expect(lockfile).toContain("patchedDependencies:");
		expect(lockfile).toContain("jsonjoy-builder@1.0.3:");
		expect(patch).toContain("dist/components/SchemaEditor/SchemaTypeSelector.js");
		expect(patch).toContain("dist/components/SchemaEditor/TypeDropdown.js");
		expect(patch).toContain('-        id: "anyOf",');
		expect(patch).toContain('-        id: "oneOf",');
		expect(patch).toContain('-        id: "allOf",');
		expect(patch).toContain('-    "null",');
	});

	it("wires JSONJoy localization from the Paraglide locale", () => {
		const source = componentSource("json-schema-builder.svelte");

		expect(source).toContain('import { getLocale } from "$lib/paraglide/runtime.js";');
		expect(source).toContain(
			'import { jsonJoyBuilderConfigForLocale } from "./json-schema-builder-locale";'
		);
		expect(source).toContain(
			"const jsonJoyBuilderConfig = $derived(jsonJoyBuilderConfigForLocale(getLocale()));"
		);
		expect(source).toContain("locale: jsonJoyBuilderConfig.locale");
		expect(source).toContain("messages: jsonJoyBuilderConfig.messages");
	});

	it("keeps the Romanian JSONJoy locale aligned with the installed English locale keys", () => {
		const englishKeys = translationKeysFromSource(jsonJoyEnglishLocaleSource());
		const romanianKeys = translationKeysFromSource(
			jsonJoyLocaleModuleSource(),
			"jsonJoyRomanianLocale"
		);

		expect(romanianKeys).toEqual(englishKeys);
	});

	it("uses Syncra Romanian wording for representative JSONJoy builder strings", () => {
		const locale = jsonJoyLocaleModuleSource();

		expect(translationValueFromSource(locale, "schemaEditorTitle", "jsonJoyRomanianLocale")).toBe(
			"Editor schemă JSON"
		);
		expect(translationValueFromSource(locale, "fieldAddNewButton", "jsonJoyRomanianLocale")).toBe(
			"Adaugă câmp"
		);
		expect(translationValueFromSource(locale, "propertyRequired", "jsonJoyRomanianLocale")).toBe(
			"Obligatoriu"
		);
		expect(translationValueFromSource(locale, "propertyOptional", "jsonJoyRomanianLocale")).toBe(
			"Opțional"
		);
		expect(translationValueFromSource(locale, "schemaTypeObject", "jsonJoyRomanianLocale")).toBe(
			"Grup"
		);
		expect(translationValueFromSource(locale, "schemaTypeArray", "jsonJoyRomanianLocale")).toBe(
			"Listă"
		);
		expect(translationValueFromSource(locale, "schemaTypeString", "jsonJoyRomanianLocale")).toBe(
			"Text"
		);
		expect(translationValueFromSource(locale, "schemaTypeBoolean", "jsonJoyRomanianLocale")).toBe(
			"Da/Nu"
		);
		expect(translationValueFromSource(locale, "validatorValid", "jsonJoyRomanianLocale")).toBe(
			"JSON-ul respectă schema."
		);
		expect(translationValueFromSource(locale, "validatorErrorInvalidSyntax", "jsonJoyRomanianLocale")).toBe(
			"Sintaxă JSON invalidă"
		);
		expect(translationValueFromSource(locale, "visualEditorNoFieldsHint1", "jsonJoyRomanianLocale")).toBe(
			"Nu există câmpuri definite"
		);
		expect(translationValueFromSource(locale, "visualEditorNoFieldsHint2", "jsonJoyRomanianLocale")).toBe(
			"Adaugă primul câmp pentru a începe"
		);
	});

	it("preserves JSONJoy interpolation placeholders in Romanian strings", () => {
		const english = jsonJoyEnglishLocaleSource();
		const romanian = jsonJoyLocaleModuleSource();

		for (const key of [
			"validatorErrorCount",
			"validatorErrorLocationLineAndColumn",
			"validatorErrorLocationLineOnly"
		]) {
			expect(
				placeholderTokens(translationValueFromSource(romanian, key, "jsonJoyRomanianLocale"))
			).toEqual(placeholderTokens(translationValueFromSource(english, key)));
		}
	});

	it("maps JSONJoy variables for the schema builder and portal dialogs", () => {
		const source = appCssSource();

		expect(source).toContain(":root .jsonjoy,");
		expect(source).toContain(":root .jsonjoy.dark,");
		expect(source).toContain(":root.dark .jsonjoy,");
		expect(source).toContain("[data-json-schema-builder] .jsonjoy,");
		expect(source).toContain("[data-json-schema-builder] .jsonjoy.dark");
		expect(source).toContain("--jsonjoy-background: var(--background);");
		expect(source).toContain("--jsonjoy-foreground: var(--foreground);");
		expect(source).toContain("--jsonjoy-card: var(--card);");
		expect(source).toContain("--jsonjoy-popover: var(--popover);");
		expect(source).toContain("--jsonjoy-primary: var(--primary);");
		expect(source).toContain("--jsonjoy-secondary: var(--secondary);");
		expect(source).toContain("--jsonjoy-muted: var(--muted);");
		expect(source).toContain("--jsonjoy-accent: var(--accent);");
		expect(source).toContain("--jsonjoy-destructive: var(--destructive);");
		expect(source).toContain("--jsonjoy-border: var(--border);");
		expect(source).toContain("--jsonjoy-input: var(--input);");
		expect(source).toContain("--jsonjoy-ring: var(--ring);");
		expect(source).toContain("--jsonjoy-radius: var(--radius);");
		expect(source).toContain('--jsonjoy-font-sans: var(--font-sans, "Inter Variable", system-ui, sans-serif);');
		expect(source).toContain("color: var(--jsonjoy-color-foreground);");
		expect(source).toContain("[data-json-schema-builder] .jsonjoy.json-editor-container,");
		expect(source).toContain("background-color: var(--jsonjoy-color-background);");
		expect(source).toContain("[data-json-schema-builder] .jsonjoy.json-editor-container > .h-\\[600px\\]");
		expect(source).toContain("height: 100%;");
		expect(source).toContain("min-height: 0;");
	});

	it("lets the structure designer fill available page height", () => {
		const editor = editorSource();
		const builder = componentSource("json-schema-builder.svelte");

		expect(editor).toContain("@container/main flex min-h-0 flex-1 flex-col");
		expect(editor).toContain("flex min-h-0 flex-1 flex-col gap-6");
		expect(editor).toContain("flex min-h-[640px] flex-1 flex-col");
		expect(editor).toContain("class=\"min-h-0 flex-1\"");
		expect(editor).toContain('class="h-full min-h-0 border-none rounded-none"');
		expect(builder).toContain('jsonJoyTheme === "dark" ? "h-full min-h-0 dark" : "h-full min-h-0"');
		expect(builder).toContain("relative h-full min-h-[480px]");
		expect(builder).toContain('class="h-full min-h-0 w-full"');
	});

	it("lets parent pages clear save and server-error feedback when fields change", () => {
		const editor = editorSource();

		expect(editor).toContain("onDirty = () => undefined");
		expect(editor).toContain("onDirty();");
		expect(newPageSource()).toContain("onDirty={clearFeedback}");
		expect(editPageSource()).toContain("onDirty={markDirty}");
	});

	it("blocks new schema saves when the schema has no fields", () => {
		const source = newPageSource();
		const submitBody = source.match(/function submit\([\s\S]*?\n\t}/)?.[0] ?? "";
		const emptyErrorIndex = submitBody.indexOf(
			"displayedServerError = m.schemas_empty_schema_error();"
		);
		const mutateIndex = submitBody.indexOf("mutation.mutate");

		expect(source).toContain('import { schemaHasFields } from "./schema-validation";');
		expect(source).not.toContain("EMPTY_SCHEMA_ERROR");
		expect(submitBody).toContain("if (!schemaHasFields(input.schema))");
		expect(emptyErrorIndex).toBeGreaterThan(-1);
		expect(mutateIndex).toBeGreaterThan(-1);
		expect(emptyErrorIndex).toBeLessThan(mutateIndex);
	});

	it("does not reset mutation lifecycle from dirty feedback handlers", () => {
		expect(newPageSource()).not.toContain("mutation.reset()");
		expect(editPageSource()).not.toContain("mutation.reset()");
	});

	it("clears edit feedback when SvelteKit reuses the route for a different schema id", () => {
		const source = editPageSource();

		expect(source).toContain("$effect(() =>");
		expect(source).toContain("schemaId;");
		expect(source).toContain("clearFeedback();");
		expect(source).toContain("dirty = false;");
		expect(source).toContain("pendingSchemaId = null;");
	});

	it("treats missing edit schemas as terminal not-found errors", () => {
		const source = editPageSource();
		const normalized = source.replace(/\s+/g, " ").trim();

		expect(source).toContain("isSchemaNotFoundError");
		expect(source).toContain("shouldRetrySchemaQuery");
		expect(normalized).toContain("retry: shouldRetrySchemaQuery");
		expect(normalized).toContain(
			"const schemaNotFound = $derived(isSchemaNotFoundError(schemaQuery.error));"
		);
		expect(normalized).toContain("{#if schemaNotFound}");
		expect(source).toContain("m.schemas_not_found_title()");
		expect(source).toContain("m.schemas_not_found_body()");
		expect(source).toContain('href="/app/schemas"');
	});

	it("announces saved feedback to assistive technologies", () => {
		const source = editorSource();

		expect(source).toContain('role="status"');
		expect(source).toContain('aria-live="polite"');
	});

	it("recreates the edit editor when refetched schema content changes", () => {
		const source = editPageSource();

		expect(source).toContain("editorKey");
		expect(source).toContain("cleanEditorKey");
		expect(source).toContain("activeEditorKey");
		expect(source).toContain("schemaQuery.data.updated_at");
		expect(source).toContain("if (!dirty) activeEditorKey = cleanEditorKey;");
		expect(source).toContain("const editorKey = $derived(activeEditorKey || cleanEditorKey);");
		expect(source).not.toContain('dirty ? `${schemaId}:dirty` : cleanEditorKey');
		expect(source).toContain("{#key editorKey}");
		expect(source).not.toContain("{#key schemaQuery.data.id}");
	});

	it("scopes edit pending state to the schema being saved", () => {
		const source = editPageSource();

		expect(source).toContain("pendingSchemaId");
		expect(source).toContain("pendingSchemaId === schemaId");
		expect(source).toContain("variables.schemaId !== schemaId");
		expect(source).toContain("pending={isCurrentSchemaSavePending}");
		expect(source).not.toContain("pending={mutation.isPending}");
	});

	it("keeps the editor key pinned while an edit save is pending", () => {
		const source = editPageSource();
		const submitBody = source.match(/function submit\([\s\S]*?\n\t}/)?.[0] ?? "";

		expect(submitBody).not.toContain("dirty = false;");
		expect(submitBody).toContain("mutation.mutate");
		expect(source).toContain("saved = result;");
		expect(source).toContain("dirty = false;");
	});

	it("keeps successful save invalidation outside stale feedback guards", () => {
		const newPage = newPageSource();
		const editPage = editPageSource();
		const newGuardIndex = newPage.indexOf("if (variables.feedbackVersion !== feedbackVersion) return;");
		const newInvalidationIndex = newPage.indexOf(
			'void queryClient.invalidateQueries({ queryKey: ["schemas"] });'
		);
		const editGuardIndex = editPage.indexOf("if (variables.feedbackVersion !== feedbackVersion");
		const editListInvalidationIndex = editPage.indexOf(
			'void queryClient.invalidateQueries({ queryKey: ["schemas"] });'
		);
		const editItemInvalidationIndex = editPage.indexOf(
			'void queryClient.invalidateQueries({ queryKey: ["schema", variables.schemaId] });'
		);

		expect(newInvalidationIndex).toBeGreaterThan(-1);
		expect(newInvalidationIndex).toBeLessThan(newGuardIndex);
		expect(editListInvalidationIndex).toBeGreaterThan(-1);
		expect(editListInvalidationIndex).toBeLessThan(editGuardIndex);
		expect(editItemInvalidationIndex).toBeGreaterThan(-1);
		expect(editItemInvalidationIndex).toBeLessThan(editGuardIndex);
	});

	it("seeds the edit query cache before unpinning after successful save", () => {
		const source = editPageSource();
		const cacheSeedIndex = source.indexOf(
			'queryClient.setQueryData<SchemaResponse>(["schema", variables.schemaId], result);'
		);
		const guardIndex = source.indexOf("if (variables.feedbackVersion !== feedbackVersion");
		const dirtyClearIndex = source.indexOf("dirty = false;", guardIndex);

		expect(cacheSeedIndex).toBeGreaterThan(-1);
		expect(cacheSeedIndex).toBeLessThan(guardIndex);
		expect(cacheSeedIndex).toBeLessThan(dirtyClearIndex);
	});

	it("renders an optional clone button from the shared schema editor", () => {
		const source = editorSource();

		expect(source).toContain("clonePending = false");
		expect(source).toContain("onClone = null");
		expect(source).toContain("{#if onClone}");
		expect(source).toContain("disabled={clonePending || pending}");
		expect(source).toContain("onclick={onClone}");
		expect(source).toContain("m.schemas_cloning()");
		expect(source).toContain("m.schemas_clone()");
	});

	it("renders schema IDs with copy controls in the list and edit header", () => {
		const listPage = listPageSource();
		const editor = editorSource();
		const editPage = editPageSource();
		const schemaIdCopy = schemaIdCopySource();

		const nameColumnIndex = listPage.indexOf("header: m.schemas_name_column()");
		const idColumnIndex = listPage.indexOf("header: m.schemas_id_column()");
		const strictColumnIndex = listPage.indexOf("header: m.schemas_strict_mode_column()");

		expect(nameColumnIndex).toBeGreaterThan(-1);
		expect(idColumnIndex).toBeGreaterThan(nameColumnIndex);
		expect(strictColumnIndex).toBeGreaterThan(idColumnIndex);
		expect(listPage).toContain('import SchemaIdCopy from "./schema-id-copy.svelte";');
		expect(listPage).toContain("renderComponent(SchemaIdCopy, { schemaId: row.original.id, compact: true })");
		expect(listPage).toContain('cell.column.id === "id" && "min-w-[220px] max-w-[320px]"');

		expect(editor).toContain("schemaId = null");
		expect(editor).toContain("schemaId?: string | null");
		expect(editor).toContain("{#if schemaId}");
		expect(editor).toContain("<SchemaIdCopy schemaId={schemaId} showLabel />");
		expect(editPage).toContain("schemaId={schemaQuery.data.id}");

		expect(schemaIdCopy).toContain('import CopyIcon from "@lucide/svelte/icons/copy";');
		expect(schemaIdCopy).toContain('import { toast } from "svelte-sonner";');
		expect(schemaIdCopy).toContain("navigator.clipboard.writeText(schemaId)");
		expect(schemaIdCopy).toContain("m.schemas_copy_id_success()");
		expect(schemaIdCopy).toContain("m.schemas_copy_id_error()");
		expect(schemaIdCopy).toContain("aria-label={m.schemas_copy_id_aria({ id: schemaId })}");
	});

	it("uses Paraglide messages for schema route labels and actions", () => {
		const listPage = listPageSource();
		const editor = editorSource();
		const newPage = newPageSource();
		const editPage = editPageSource();
		const actionsCell = actionsCellSource();
		const schemaIdCopy = schemaIdCopySource();
		const nameCell = readFileSync(new URL("./name-cell.svelte", import.meta.url), "utf8");
		const createdHeader = readFileSync(
			new URL("./created-date-header.svelte", import.meta.url),
			"utf8"
		);

		for (const source of [
			listPage,
			editor,
			newPage,
			editPage,
			actionsCell,
			schemaIdCopy,
			nameCell,
			createdHeader
		]) {
			expect(source).toContain('import { m } from "$lib/paraglide/messages.js";');
		}

		for (const messageCall of [
			"m.schemas_delete_single_title()",
			"m.schemas_delete_single_description({ name: schema.name })",
			"m.schemas_select_all_on_page()",
			"m.schemas_name_column()",
			"m.schemas_id_column()",
			"m.schemas_strict_mode_column()",
			"m.schemas_updated_column()",
			"m.schemas_no_schemas_found()",
			"m.schemas_empty_body()",
			"m.schemas_showing_schemas_one({ count: schemas.length })",
			"m.schemas_selected_count_other({ count: selectedIds.size })"
		]) {
			expect(listPage).toContain(messageCall);
		}

		for (const messageCall of [
			"m.schemas_editor_badge()",
			"m.schemas_general_settings()",
			"m.schemas_schema_name_label()",
			"m.schemas_schema_name_placeholder()",
			"m.schemas_structure_designer()",
			"m.schemas_saving()",
			"m.schemas_cloning()"
		]) {
			expect(editor).toContain(messageCall);
		}

		for (const messageCall of [
			"m.schemas_id_label()",
			"m.schemas_copy_id()",
			"m.schemas_copy_id_aria({ id: schemaId })",
			"m.schemas_copy_id_success()",
			"m.schemas_copy_id_error()"
		]) {
			expect(schemaIdCopy).toContain(messageCall);
		}

		expect(newPage).toContain("m.schemas_new_title()");
		expect(newPage).toContain("m.schemas_saved_success_with_id({ name: result.name, id: result.id })");
		expect(editPage).toContain("m.schemas_edit_title()");
		expect(editPage).toContain("m.schemas_not_found_title()");
		expect(actionsCell).toContain("m.schemas_create_job_with({ name: schema.name })");
		expect(nameCell).toContain("m.schemas_no_description()");
		expect(createdHeader).toContain("m.schemas_sort_created_ascending()");
	});

	it("wires edit-page clone mutation to redirect to the cloned schema", () => {
		const source = editPageSource();

		expect(source).toContain('import { goto } from "$app/navigation";');
		expect(source).toContain("cloneSchema,");
		expect(source).toContain(
			"const cloneMutation = createMutation<SchemaResponse, Error, { schema: SchemaResponse }>"
		);
		expect(source).toContain("mutationFn: ({ schema }) => cloneSchema(fetch, schema)");
		expect(source).toContain(
			'queryClient.setQueryData<SchemaResponse>(["schema", result.id], result);'
		);
		expect(source).toContain('void queryClient.invalidateQueries({ queryKey: ["schemas"] });');
		expect(source).toContain("void goto(`/app/schemas/edit/${result.id}`);");
		expect(source).toContain("cloneMutation.mutate({ schema: schemaQuery.data });");
		expect(source).toContain("clonePending={cloneMutation.isPending}");
		expect(source).toContain("onClone={cloneCurrentSchema}");
	});

	it("adds clone icon wiring to the schemas action column", () => {
		const list = listPageSource();
		const actions = actionsCellSource();

		expect(actions).toContain('import CopyPlusIcon from "@lucide/svelte/icons/copy-plus";');
		expect(actions).toContain('import RocketIcon from "@lucide/svelte/icons/rocket";');
		expect(actions).toContain('import { buildNewJobPath } from "../new-job/schema-query";');
		expect(actions).toContain("onClone");
		expect(actions).toContain("clonePending = false");
		expect(actions).toContain("aria-label={m.schemas_clone_aria({ name: schema.name })}");
		expect(actions).toContain("onclick={() => onClone(schema)}");
		expect(actions).toContain("<CopyPlusIcon");
		expect(actions).toContain("href={buildNewJobPath(schema.id)}");
		expect(actions).toContain("aria-label={m.schemas_create_job_with({ name: schema.name })}");
		expect(actions).toContain("title={m.schemas_create_job_with({ name: schema.name })}");
		expect(actions).toContain("<RocketIcon");
		expect(list).toContain("cloneSchema,");
		expect(list).toContain("const cloneMutation = createMutation<SchemaResponse, Error, CloneSchemaVariables>");
		expect(list).toContain("mutationFn: ({ schema }) => cloneSchema(fetch, schema)");
		expect(list).toContain("void goto(`/app/schemas/edit/${result.id}`);");
		expect(list).toContain("onClone: cloneSingleSchema");
		expect(list).toContain("clonePending: cloneMutation.isPending");
		expect(list).toContain('cell.column.id === "actions" && "w-40"');
	});

	it("links to the schema recipe library before the create schema action", () => {
		const source = listPageSource();
		const libraryIndex = source.indexOf('href="/app/schemas/library"');
		const newSchemaIndex = source.indexOf('href="/app/schemas/new"');

		expect(source).toContain('import BookOpenIcon from "@lucide/svelte/icons/book-open";');
		expect(source).toContain("m.schemas_library()");
		expect(libraryIndex).toBeGreaterThan(-1);
		expect(newSchemaIndex).toBeGreaterThan(-1);
		expect(libraryIndex).toBeLessThan(newSchemaIndex);
	});

	it("points navigation chrome at canonical schema routes", () => {
		const sidebar = componentSource("app-sidebar.svelte");
		const header = componentSource("site-header.svelte");

		expect(sidebar).toContain('url: "/app/schemas/new"');
		expect(header).toContain('pathname === "/app/schemas/library"');
		expect(header).toContain("m.schemas_library()");
		expect(header).toContain('pathname === "/app/schemas/new"');
		expect(header).toContain('pathname.startsWith("/app/schemas/edit/")');
	});

	it("points new-job empty schema guidance at the canonical create route", () => {
		const source = newJobPageSource();

		expect(source).toContain('href="/app/schemas/new"');
		expect(source).toContain('import { page } from "$app/state";');
		expect(source).toContain("schemaIdFromSearchParams(page.url.searchParams)");
		expect(source).toContain("isSelectedSchemaIdUnavailable({");
	});
});
