import { CalendarDate } from "@internationalized/date";
import { describe, expect, it } from "vitest";

import { DatasetClientError } from "$lib/client/datasets";
import {
	DATASET_ROWS_QUERY_RETRY_LIMIT,
	cursorNextState,
	cursorPreviousState,
	datasetCellText,
	datasetExportFilename,
	dateRangeToQueryBounds,
	formatDatasetDate,
	resetCursorState,
	shouldRetryDatasetRowsQuery,
} from "./table-utils";

describe("datasets table utils", () => {
	it("formats dataset cell values for table display", () => {
		expect(datasetCellText(null)).toBe("");
		expect(datasetCellText(undefined)).toBe("");
		expect(datasetCellText("Invoice total")).toBe("Invoice total");
		expect(datasetCellText('[{"amount":100}]')).toBe('[{"amount":100}]');
		expect(datasetCellText(42.5)).toBe("42.5");
		expect(datasetCellText(true)).toBe("true");
		expect(datasetCellText(false)).toBe("false");
		expect(datasetCellText(["ocr", "reviewed"])).toBe('["ocr","reviewed"]');
		expect(datasetCellText({ currency: "USD", tags: ["ocr", "reviewed"] })).toBe(
			'{"currency":"USD","tags":["ocr","reviewed"]}'
		);
	});

	it("builds safe readable export filenames with collision-resistant timestamps", () => {
		const date = new Date("2026-06-07T13:14:15.000Z");

		expect(datasetExportFilename("Invoice dataset", "csv", date)).toBe(
			"Invoice dataset 2026-06-07 13-14-15.csv"
		);
		expect(datasetExportFilename("Q2 / invoices: east\\west", "xlsx", date)).toBe(
			"Q2 invoices east west 2026-06-07 13-14-15.xlsx"
		);
		expect(datasetExportFilename("Client - A / Batch", "csv", date)).toBe(
			"Client - A Batch 2026-06-07 13-14-15.csv"
		);
		expect(datasetExportFilename("\u0000<>:\"/\\|?*  ", "xlsx", date)).toBe(
			"dataset 2026-06-07 13-14-15.xlsx"
		);
	});

	it("uses a stable fallback timestamp for invalid dates", () => {
		expect(datasetExportFilename("Dataset", "csv", new Date("not-a-date"))).toBe(
			"Dataset unknown-date.csv"
		);
	});

	it("converts range calendar values into inclusive RFC3339 bounds", () => {
		expect(
			dateRangeToQueryBounds(
				{
					start: new CalendarDate(2026, 6, 1),
					end: new CalendarDate(2026, 6, 7),
				},
				"UTC"
			)
		).toEqual({
			createdFrom: "2026-06-01T00:00:00.000Z",
			createdTo: "2026-06-07T23:59:59.999Z",
		});
		expect(dateRangeToQueryBounds({ start: new CalendarDate(2026, 6, 1) }, "UTC")).toEqual({
			createdFrom: "2026-06-01T00:00:00.000Z",
			createdTo: undefined,
		});
		expect(dateRangeToQueryBounds({ end: new CalendarDate(2026, 6, 7) }, "UTC")).toEqual({
			createdFrom: undefined,
			createdTo: "2026-06-07T23:59:59.999Z",
		});
		expect(dateRangeToQueryBounds(undefined, "UTC")).toEqual({
			createdFrom: undefined,
			createdTo: undefined,
		});
	});

	it("formats dataset dates for compact table display", () => {
		expect(formatDatasetDate("2026-05-27T15:20:30Z")).toBe("2026-05-27");
		expect(formatDatasetDate("not-a-date")).toBe("Invalid date");
		expect(formatDatasetDate("not-a-date", "Dată invalidă")).toBe("Dată invalidă");
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

	it("does not advance cursor state when the next cursor is empty", () => {
		const state = { currentCursor: "cursor-1", history: [null] };

		expect(cursorNextState(state)).toBe(state);
		expect(cursorNextState(state, null)).toBe(state);
		expect(cursorNextState(state, "")).toBe(state);
		expect(cursorNextState(state, "   ")).toBe(state);
	});

	it("does not retry missing dataset row queries", () => {
		expect(DATASET_ROWS_QUERY_RETRY_LIMIT).toBe(3);
		expect(shouldRetryDatasetRowsQuery(0, new DatasetClientError(404, "dataset not found"))).toBe(
			false
		);
		expect(
			shouldRetryDatasetRowsQuery(0, new DatasetClientError(500, "dataset service unavailable"))
		).toBe(true);
		expect(shouldRetryDatasetRowsQuery(3, new Error("network unavailable"))).toBe(false);
	});
});
