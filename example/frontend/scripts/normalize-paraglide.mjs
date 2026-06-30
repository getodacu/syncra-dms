import { readdirSync, readFileSync, statSync, writeFileSync } from "node:fs";
import { join } from "node:path";
import { fileURLToPath } from "node:url";

const root = join(fileURLToPath(new URL("..", import.meta.url)), "src/lib/paraglide");
const textFilePattern = /\.(?:d\.ts|js|json|md|ts)$/;

export function normalizeParaglide() {
	for (const file of generatedFiles(root)) {
		if (!textFilePattern.test(file)) continue;

		const original = readFileSync(file, "utf8");
		const normalized = ensureFinalNewline(original.replace(/[ \t]+$/gm, ""));
		if (normalized !== original) {
			writeFileSync(file, normalized);
		}
	}
}

/**
 * @param {string} directory
 * @returns {Generator<string>}
 */
function* generatedFiles(directory) {
	for (const entry of readdirSync(directory)) {
		const path = join(directory, entry);
		const stat = statSync(path);

		if (stat.isDirectory()) {
			yield* generatedFiles(path);
		} else if (stat.isFile()) {
			yield path;
		}
	}
}

/**
 * @param {string} value
 * @returns {string}
 */
function ensureFinalNewline(value) {
	return value.length > 0 && !value.endsWith("\n") ? `${value}\n` : value;
}

if (process.argv[1] === fileURLToPath(import.meta.url)) {
	normalizeParaglide();
}
