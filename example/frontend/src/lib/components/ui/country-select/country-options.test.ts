import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";

import { countryOptions, getCountryOption } from "./country-options";

const countryOptionsSource = () =>
	readFileSync(new URL("./country-options.ts", import.meta.url), "utf8");

describe("country options", () => {
	it("maps Romania to English and Romanian display data", () => {
		expect(getCountryOption("RO")).toMatchObject({
			code: "RO",
			en: "Romania",
			ro: "România",
			flagSrc: "/flags/24x24/ro.png"
		});
	});

	it("exposes uppercase codes and the expected selector fields", () => {
		const country = getCountryOption("ro");

		expect(country).not.toBeNull();
		expect(countryOptions.every((option) => option.code === option.code.toUpperCase())).toBe(true);
		expect(Object.keys(country ?? {}).sort()).toEqual(
			["code", "en", "flagSrc", "keywords", "ro"].sort()
		);
		expect(country).not.toHaveProperty("alpha2");
		expect(country).not.toHaveProperty("alpha3");
		expect(country).not.toHaveProperty("fr");
	});

	it("includes country names and ISO code in search keywords", () => {
		const country = getCountryOption("RO");

		expect(country?.keywords).toEqual(expect.arrayContaining(["Romania", "România", "RO", "ro"]));
	});

	it("hardcodes the English and Romanian country data in the helper", () => {
		const source = countryOptionsSource();

		expect(source).toContain("const COUNTRIES: CountryData[] = [");
		expect(source).toContain('{ code: "RO", en: "Romania", ro: "România" }');
		expect(source).not.toContain("countries.json");
	});
});
