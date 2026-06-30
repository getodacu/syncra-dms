import { describe, expect, it } from "vitest";

import {
	buildJobsQueryPath,
	cursorNextState,
	cursorPreviousState,
	formatCreatedDate,
	formatFileSize,
	headerSelectionState,
	isTerminalJobStatus,
	normalizeJobStatus,
	resetCursorState,
	shouldPollJobStatus,
	togglePageSelection,
	toggleSelection,
} from "./table-utils";

describe("jobs table utils", () => {
	it("builds job list paths with only non-empty query parameters", () => {
		expect(buildJobsQueryPath({ size: 20, sort: "desc" })).toBe(
			"/api/ocr/jobs?size=20&sort=desc"
		);
		expect(
			buildJobsQueryPath({
				status: "processing",
				cursor: "cursor-1",
				size: " 50 ",
				sort: "asc",
			})
		).toBe("/api/ocr/jobs?status=processing&cursor=cursor-1&size=50&sort=asc");
		expect(buildJobsQueryPath({ status: " ", cursor: "", size: 20 })).toBe(
			"/api/ocr/jobs?size=20"
		);
	});

	it("formats table display values", () => {
		expect(formatCreatedDate("2026-05-27T15:20:30Z")).toBe("2026-05-27");
		expect(formatCreatedDate("not-a-date")).toBe("Invalid date");
		expect(formatFileSize(12)).toBe("12 B");
		expect(formatFileSize(1024)).toBe("1 KB");
		expect(formatFileSize(1536)).toBe("1.5 KB");
		expect(formatFileSize(1048576)).toBe("1 MB");
		expect(formatFileSize(0)).toBe("0 B");
		expect(formatFileSize(Number.NaN)).toBe("0 B");
	});

	it("classifies terminal and polling statuses", () => {
		expect(normalizeJobStatus("queued")).toBe("queued");
		expect(normalizeJobStatus("pending")).toBe("pending");
		expect(normalizeJobStatus("processing")).toBe("processing");
		expect(normalizeJobStatus("completed")).toBe("completed");
		expect(normalizeJobStatus("failed")).toBe("failed");
		expect(normalizeJobStatus("unexpected")).toBe("unknown");
		expect(isTerminalJobStatus("completed")).toBe(true);
		expect(isTerminalJobStatus("failed")).toBe(true);
		expect(shouldPollJobStatus("queued")).toBe(true);
		expect(shouldPollJobStatus("pending")).toBe(true);
		expect(shouldPollJobStatus("processing")).toBe(true);
		expect(shouldPollJobStatus("completed")).toBe(false);
		expect(shouldPollJobStatus("failed")).toBe(false);
	});

	it("tracks visible row selection without dropping hidden selections", () => {
		expect(headerSelectionState(["a", "b"], new Set(["a"]))).toEqual({
			checked: false,
			indeterminate: true,
			selectedCount: 1,
		});
		expect([...toggleSelection(new Set(["a"]), "b", true)].sort()).toEqual(["a", "b"]);
		expect([...toggleSelection(new Set(["a", "b"]), "a", false)]).toEqual(["b"]);
		expect([...togglePageSelection(["a", "b"], new Set(["hidden"]), true)].sort()).toEqual([
			"a",
			"b",
			"hidden",
		]);
	});

	it("maintains cursor history for next, previous, and reset", () => {
		expect(cursorNextState({ currentCursor: null, history: [] }, "cursor-1")).toEqual({
			currentCursor: "cursor-1",
			history: [null],
		});
		expect(cursorPreviousState({ currentCursor: "cursor-1", history: [null] })).toEqual({
			currentCursor: null,
			history: [],
		});
		expect(cursorNextState({ currentCursor: "cursor-1", history: [null] }, "")).toEqual({
			currentCursor: "cursor-1",
			history: [null],
		});
		expect(resetCursorState()).toEqual({ currentCursor: null, history: [] });
	});
});
