import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";

const source = () => readFileSync(new URL("./document-preview-dialog.svelte", import.meta.url), "utf8");

function normalizeSource(value: string) {
	return value.replace(/\s+/g, " ").trim();
}

describe("document preview dialog source", () => {
	it("supports optional inline filename renaming", () => {
		const dialog = source();
		const normalized = normalizeSource(dialog);

		expect(dialog).toContain('import PencilIcon from "@lucide/svelte/icons/pencil";');
		expect(dialog).toContain('import { m } from "$lib/paraglide/messages.js";');
		expect(normalized).toContain("onRename?: (originalFilename: string) => void | Promise<void>;");
		expect(normalized).toContain("renamePending?: boolean;");
		expect(normalized).toContain("function startRenaming()");
		expect(normalized).toContain("function handleRenameKeydown(event: KeyboardEvent)");
		expect(normalized).toContain('event.key === "Enter"');
		expect(normalized).toContain("await onRename?.(nextFilename);");
		expect(normalized).toContain("{#if canRename}");
		expect(normalized).toContain("m.documents_preview_fallback_title()");
		expect(normalized).toContain("m.documents_preview_description()");
		expect(normalized).toContain("m.documents_rename_document_title()");
		expect(normalized).toContain("m.documents_loading_document()");
		expect(normalized).toContain("m.documents_copy_markdown()");
		expect(normalized).toContain("m.documents_no_preview_available()");
	});
});
