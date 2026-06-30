import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";

const source = () => readFileSync(new URL("./download-dialog.svelte", import.meta.url), "utf8");

function normalizeSource(value: string) {
	return value.replace(/\s+/g, " ").trim();
}

describe("document download dialog", () => {
	it("shows JSON only when at least one target has a saved or inline schema", () => {
		const dialog = normalizeSource(source());

		expect(dialog).toContain(
			"documents.some((document) => Boolean(document.schema_id || document.has_inline_schema))"
		);
		expect(dialog).toContain("{#if hasJSONOption}");
		expect(dialog).toContain('onclick={() => onDownload("json")}');
	});

	it("disables format actions and renders a spinner while pending", () => {
		const rawDialog = source();
		const dialog = normalizeSource(rawDialog);

		expect(rawDialog).toContain('import { m } from "$lib/paraglide/messages.js";');
		expect(dialog).toContain("m.documents_download_dialog_title_one()");
		expect(dialog).toContain("m.documents_download_dialog_title_other({ count: selectedCount })");
		expect(dialog).toContain("m.documents_selected_documents()");
		expect(dialog).toContain("disabled={pending}");
		expect(dialog).toContain("{#if pending}");
		expect(dialog).toContain('<Spinner class="size-3.5" />');
		expect(dialog).toContain("m.documents_preparing_download()");
	});
});
