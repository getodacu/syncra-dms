import { describe, expect, it } from "vitest";

import {
	buildSchemasQueryPath,
	cursorNextState,
	cursorPreviousState,
	formatDate,
	headerSelectionState,
	resetCursorState,
	togglePageSelection,
	toggleSelection,
} from "./table-utils";

describe("schemas table utils", () => {
	it("builds mine-scoped schema list paths with non-empty query parameters", () => {
		expect(buildSchemasQueryPath({})).toBe("/api/schemas?scope=mine&size=20");
		expect(buildSchemasQueryPath({ size: 20, sort: "desc" })).toBe(
			"/api/schemas?scope=mine&size=20&sort=desc"
		);
		expect(buildSchemasQueryPath({ cursor: " ", size: " 50 ", sort: "asc" })).toBe(
			"/api/schemas?scope=mine&size=50&sort=asc"
		);
	});

	it("encodes cursor values as query parameters", () => {
		expect(buildSchemasQueryPath({ cursor: "cursor:/next page", size: 25 })).toBe(
			"/api/schemas?scope=mine&cursor=cursor%3A%2Fnext+page&size=25"
		);
	});

	it("tracks header checkbox state for visible schema rows", () => {
		expect(headerSelectionState(["schema-1", "schema-2"], new Set())).toEqual({
			checked: false,
			indeterminate: false,
			selectedCount: 0,
		});
		expect(headerSelectionState(["schema-1", "schema-2"], new Set(["schema-1"]))).toEqual({
			checked: false,
			indeterminate: true,
			selectedCount: 1,
		});
		expect(
			headerSelectionState(
				["schema-1", "schema-2"],
				new Set(["schema-1", "schema-2", "hidden"])
			)
		).toEqual({
			checked: true,
			indeterminate: false,
			selectedCount: 2,
		});
	});

	it("toggles individual and visible-page selections without dropping hidden selections", () => {
		expect([...toggleSelection(new Set(["schema-1"]), "schema-2", true)].sort()).toEqual([
			"schema-1",
			"schema-2",
		]);
		expect([...toggleSelection(new Set(["schema-1", "schema-2"]), "schema-1", false)]).toEqual([
			"schema-2",
		]);
		expect(
			[...togglePageSelection(["schema-1", "schema-2"], new Set(["hidden"]), true)].sort()
		).toEqual(["hidden", "schema-1", "schema-2"]);
		expect(
			[
				...togglePageSelection(
					["schema-1", "schema-2"],
					new Set(["schema-1", "schema-2", "hidden"]),
					false
				),
			]
		).toEqual(["hidden"]);

		const originalSelected = new Set(["schema-1"]);
		const nextSelected = toggleSelection(originalSelected, "schema-2", true);
		expect(nextSelected).not.toBe(originalSelected);
		expect([...originalSelected]).toEqual(["schema-1"]);

		const originalPageSelected = new Set(["hidden"]);
		const nextPageSelected = togglePageSelection(
			["schema-1", "schema-2"],
			originalPageSelected,
			true
		);
		expect(nextPageSelected).not.toBe(originalPageSelected);
		expect([...originalPageSelected]).toEqual(["hidden"]);
	});

	it("maintains cursor history for next, previous, and reset", () => {
		expect(cursorNextState({ currentCursor: null, history: [] }, "cursor-1")).toEqual({
			currentCursor: "cursor-1",
			history: [null],
		});
		expect(cursorNextState({ currentCursor: "cursor-1", history: [null] }, "cursor-2")).toEqual({
			currentCursor: "cursor-2",
			history: [null, "cursor-1"],
		});
		expect(
			cursorPreviousState({ currentCursor: "cursor-2", history: [null, "cursor-1"] })
		).toEqual({
			currentCursor: "cursor-1",
			history: [null],
		});
		expect(cursorPreviousState({ currentCursor: "cursor-1", history: [] })).toEqual({
			currentCursor: "cursor-1",
			history: [],
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

	it("formats schema dates", () => {
		expect(formatDate("2026-05-27T15:20:30Z")).toBe("2026-05-27");
		expect(formatDate("not-a-date")).toBe("Invalid date");
	});
});
