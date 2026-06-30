import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";

const source = () => readFileSync(new URL("./actions-cell.svelte", import.meta.url), "utf8");

function normalizeSource(value: string) {
	return value.replace(/\s+/g, " ").trim();
}

describe("documents actions cell source", () => {
	it("renders row actions in a dropdown menu", () => {
		const cell = source();
		const normalized = normalizeSource(cell);

		expect(cell).toContain('import DotsVerticalIcon from "@tabler/icons-svelte/icons/dots-vertical";');
		expect(cell).toContain('import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index.js";');
		expect(cell).toContain('import { m } from "$lib/paraglide/messages.js";');
		expect(normalized).toContain("onRename: (document: OCRDocumentListItemResponse) => void;");
		expect(normalized).toContain("let menuOpen = $state(false);");
		expect(normalized).toContain("function selectRename(event: Event)");
		expect(normalized).toContain("event.preventDefault();");
		expect(normalized).toContain("menuOpen = false;");
		expect(normalized).toContain("setTimeout(() => { onRename(document); }, 0);");
		expect(normalized).toContain("bind:open={menuOpen}");
		expect(normalized).toContain('<DropdownMenu.Content align="end" class="w-36">');
		expect(normalized).toContain(
			"aria-label={m.documents_open_actions_for({ name: document.original_filename })}"
		);
		expect(normalized).toContain("<DotsVerticalIcon");
		expect(normalized).toContain("onSelect={() => onPreview(document)}");
		expect(normalized).toContain("{m.documents_preview()}");
		expect(normalized).toContain("disabled={renamePending} onSelect={selectRename}");
		expect(normalized).toContain("{m.documents_rename()}");
		expect(normalized).toContain("disabled={downloadPending} onSelect={() => onDownload(document)}");
		expect(normalized).toContain("{m.documents_download()}");
		expect(normalized).toContain('variant="destructive" disabled={deletePending} onSelect={() => onDelete(document)}');
		expect(normalized).toContain("{m.documents_delete()}");
		expect(cell).not.toContain("EyeIcon");
		expect(cell).not.toContain("PencilIcon");
		expect(cell).not.toContain("DownloadIcon");
		expect(cell).not.toContain("Trash2Icon");
	});
});
