import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";

const source = () => readFileSync(new URL("./filename-cell.svelte", import.meta.url), "utf8");

function normalizeSource(value: string) {
	return value.replace(/\s+/g, " ").trim();
}

describe("documents filename cell source", () => {
	it("receives controlled editing and preview callbacks", () => {
		const cell = source();
		const normalized = normalizeSource(cell);

		expect(cell).toContain('import { m } from "$lib/paraglide/messages.js";');
		expect(cell).toContain("onPreview");
		expect(cell).toContain("onEditingChange");
		expect(normalized).toContain("editing = false,");
		expect(normalized).toContain("editing?: boolean | (() => boolean);");
		expect(normalized).toContain('const isEditing = $derived.by(() => (typeof editing === "function" ? editing() : editing));');
		expect(normalized).toContain("onPreview: (document: OCRDocumentListItemResponse) => void;");
		expect(normalized).toContain("onEditingChange?: (editing: boolean) => void;");
	});

	it("opens preview from the filename display instead of starting rename", () => {
		const cell = source();
		const normalized = normalizeSource(cell);

		expect(normalized).toContain("onclick={() => onPreview(document)}");
		expect(normalized).not.toContain("onclick={startEditing}");
		expect(normalized).toContain(
			"aria-label={m.documents_preview_file({ name: document.original_filename })}"
		);
		expect(normalized).toContain(
			"aria-label={m.documents_rename_file({ name: document.original_filename })}"
		);
		expect(normalized).toContain("m.documents_failed_rename()");
	});
});
