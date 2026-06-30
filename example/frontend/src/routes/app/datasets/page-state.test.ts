import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";

const detailSource = () => readFileSync(new URL("./[id]/+page.svelte", import.meta.url), "utf8");
const normalizeSource = (value: string) => value.replace(/\s+/g, " ").trim();

describe("dataset detail page state", () => {
	it("treats missing datasets as terminal not-found errors", () => {
		const page = detailSource();
		const normalized = normalizeSource(page);

		expect(page).toContain('import { isDatasetNotFoundError } from "$lib/client/datasets";');
		expect(page).toContain("shouldRetryDatasetRowsQuery");
		expect(normalized).toContain("retry: shouldRetryDatasetRowsQuery");
		expect(normalized).toContain(
			"const datasetNotFound = $derived(isDatasetNotFoundError(rowsQuery.error));"
		);
		expect(normalized).toContain("{#if rowsQuery.isError && !datasetNotFound}");
		expect(normalized).toContain("{:else if datasetNotFound}");
		expect(page).toContain('import { m } from "$lib/paraglide/messages.js";');
		expect(page).toContain("m.datasets_not_found_body()");
		expect(page).toContain('href="/app/datasets"');
	});
});
