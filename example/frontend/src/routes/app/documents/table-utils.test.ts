import { CalendarDate } from "@internationalized/date";
import { describe, expect, it } from "vitest";

import {
	buildDocumentsQueryPath,
	cursorNextState,
	cursorPreviousState,
	dateRangeToQueryBounds,
	fileIconKind,
	formatCreatedDate,
	formatFileSize,
	headerSelectionState,
	resetCursorState,
	togglePageSelection,
	toggleSelection,
	truncateFilename,
} from "./table-utils";

describe("documents table utils", () => {
	it("builds document list paths with only non-empty query parameters", () => {
		expect(buildDocumentsQueryPath({ size: 20, sort: "desc" })).toBe(
			"/api/ocr/documents?size=20&sort=desc"
		);
		expect(buildDocumentsQueryPath({ collectionId: "collection-1", size: 20 })).toBe(
			"/api/ocr/documents?collection=collection-1&size=20"
		);
		expect(buildDocumentsQueryPath({ filename: " ", cursor: "", size: " 50 ", sort: "asc" })).toBe(
			"/api/ocr/documents?size=50&sort=asc"
		);
		expect(
			buildDocumentsQueryPath({
				filename: " invoice ",
				createdFrom: "2026-05-01T00:00:00.000Z",
				createdTo: "2026-05-29T23:59:59.999Z",
				cursor: "cursor-1",
				size: 50,
				sort: "asc",
			})
		).toBe(
			"/api/ocr/documents?filename=invoice&created_from=2026-05-01T00%3A00%3A00.000Z&created_to=2026-05-29T23%3A59%3A59.999Z&cursor=cursor-1&size=50&sort=asc"
		);
	});

	it("converts range calendar values into inclusive RFC3339 bounds", () => {
		expect(
			dateRangeToQueryBounds(
				{
					start: new CalendarDate(2026, 5, 1),
					end: new CalendarDate(2026, 5, 29),
				},
				"UTC"
			)
		).toEqual({
			createdFrom: "2026-05-01T00:00:00.000Z",
			createdTo: "2026-05-29T23:59:59.999Z",
		});
		expect(dateRangeToQueryBounds({ start: new CalendarDate(2026, 5, 1) }, "UTC")).toEqual({
			createdFrom: "2026-05-01T00:00:00.000Z",
			createdTo: undefined,
		});
		expect(dateRangeToQueryBounds(undefined, "UTC")).toEqual({
			createdFrom: undefined,
			createdTo: undefined,
		});
		expect(dateRangeToQueryBounds({ end: new CalendarDate(2026, 5, 29) }, "UTC")).toEqual({
			createdFrom: undefined,
			createdTo: "2026-05-29T23:59:59.999Z",
		});
	});

	it("formats table display values", () => {
		expect(truncateFilename("short.pdf")).toBe("short.pdf");
		expect(truncateFilename("averyverylongfilename.pdf")).toBe("averyverylongfile...");
		expect(truncateFilename("averyverylongfilename.pdf")).toHaveLength(20);
		expect(formatCreatedDate("2026-05-27T15:20:30Z")).toBe("2026-05-27");
		expect(formatCreatedDate("not-a-date")).toBe("Invalid date");
		expect(formatCreatedDate("not-a-date", "Dată invalidă")).toBe("Dată invalidă");
		expect(formatFileSize(12)).toBe("12 B");
		expect(formatFileSize(1024)).toBe("1 KB");
		expect(formatFileSize(1536)).toBe("1.5 KB");
		expect(formatFileSize(1048576)).toBe("1 MB");
		expect(formatFileSize(0)).toBe("0 B");
		expect(formatFileSize(-1)).toBe("0 B");
		expect(formatFileSize(Number.NaN)).toBe("0 B");
		expect(formatFileSize(Number.POSITIVE_INFINITY)).toBe("0 B");
	});

	it("classifies file icon kinds from mime type", () => {
		expect(fileIconKind("application/pdf")).toBe("pdf");
		expect(fileIconKind("image/png")).toBe("image");
		expect(fileIconKind("text/plain")).toBe("file");
		expect(fileIconKind("")).toBe("file");
	});

	it("tracks header checkbox state for visible rows", () => {
		expect(headerSelectionState(["a", "b"], new Set())).toEqual({
			checked: false,
			indeterminate: false,
			selectedCount: 0,
		});
		expect(headerSelectionState(["a", "b"], new Set(["a"]))).toEqual({
			checked: false,
			indeterminate: true,
			selectedCount: 1,
		});
		expect(headerSelectionState(["a", "b"], new Set(["a", "b", "hidden"]))).toEqual({
			checked: true,
			indeterminate: false,
			selectedCount: 2,
		});
	});

	it("toggles individual and visible-page selections without dropping hidden selections", () => {
		expect([...toggleSelection(new Set(["a"]), "b", true)].sort()).toEqual(["a", "b"]);
		expect([...toggleSelection(new Set(["a", "b"]), "a", false)].sort()).toEqual(["b"]);
		expect([...togglePageSelection(["a", "b"], new Set(["hidden"]), true)].sort()).toEqual([
			"a",
			"b",
			"hidden",
		]);
		expect([...togglePageSelection(["a", "b"], new Set(["a", "b", "hidden"]), false)]).toEqual([
			"hidden",
		]);

		const originalSelected = new Set(["a"]);
		const nextSelected = toggleSelection(originalSelected, "b", true);
		expect(nextSelected).not.toBe(originalSelected);
		expect([...originalSelected]).toEqual(["a"]);

		const originalPageSelected = new Set(["hidden"]);
		const nextPageSelected = togglePageSelection(["a", "b"], originalPageSelected, true);
		expect(nextPageSelected).not.toBe(originalPageSelected);
		expect([...originalPageSelected]).toEqual(["hidden"]);
	});

	it("maintains cursor history for next, previous, and reset", () => {
		expect(cursorNextState({ currentCursor: null, history: [] }, "cursor-1")).toEqual({
			currentCursor: "cursor-1",
			history: [null],
		});
		expect(
			cursorNextState({ currentCursor: "cursor-1", history: [null] }, "cursor-2")
		).toEqual({
			currentCursor: "cursor-2",
			history: [null, "cursor-1"],
		});
		expect(cursorPreviousState({ currentCursor: "cursor-2", history: [null, "cursor-1"] })).toEqual({
			currentCursor: "cursor-1",
			history: [null],
		});
		expect(resetCursorState()).toEqual({ currentCursor: null, history: [] });
	});

	it("does not advance cursor state when the next cursor is unusable", () => {
		const state = { currentCursor: "cursor-1", history: [null] };

		expect(cursorNextState(state)).toBe(state);
		expect(cursorNextState(state, null)).toBe(state);
		expect(cursorNextState(state, "")).toBe(state);
		expect(cursorNextState(state, "   ")).toBe(state);
	});
});
