import JSZip from "jszip";

import { renderMarkdown } from "$lib/components/ui/document-preview/document-preview-utils";
import type { OCRDocumentResponse } from "$lib/server/ocr";

export type DownloadFormat = "markdown" | "html" | "json";

export type DownloadFile = {
	filename: string;
	content: string;
	contentType: string;
};

const FORMAT_CONFIG: Record<DownloadFormat, { extension: string; contentType: string }> = {
	markdown: { extension: ".md", contentType: "text/markdown; charset=utf-8" },
	html: { extension: ".html", contentType: "text/html; charset=utf-8" },
	json: { extension: ".json", contentType: "application/json; charset=utf-8" }
};

const UNSAFE_FILENAME_CHARS = /[\u0000-\u001f\u007f<>:"/\\|?*]+/g;

export function isDownloadFormat(value: unknown): value is DownloadFormat {
	return value === "markdown" || value === "html" || value === "json";
}

export function stripFileExtension(filename: string) {
	const basename = filename.replaceAll("\\", "/").split("/").pop()?.trim() ?? "";
	const extensionIndex = basename.lastIndexOf(".");

	if (extensionIndex <= 0) return basename;

	return basename.slice(0, extensionIndex);
}

export function sanitizeDownloadBasename(filename: string, fallback = "document") {
	const basename = stripFileExtension(filename)
		.replace(UNSAFE_FILENAME_CHARS, "-")
		.replace(/\s+/g, " ")
		.replace(/^[. ]+|[. ]+$/g, "")
		.trim();

	return basename || fallback;
}

export function documentHasSchema(document: Pick<OCRDocumentResponse, "schema_id" | "has_inline_schema">) {
	return Boolean(document.schema_id || document.has_inline_schema);
}

export function renderStandaloneHTML(document: Pick<OCRDocumentResponse, "original_filename" | "markdown">) {
	const title = sanitizeDownloadBasename(document.original_filename);
	const body = renderMarkdown(document.markdown);

	return [
		"<!doctype html>",
		'<html lang="en">',
		"<head>",
		'<meta charset="utf-8">',
		'<meta name="viewport" content="width=device-width, initial-scale=1">',
		`<title>${escapeHTML(title)}</title>`,
		"</head>",
		"<body>",
		body,
		"</body>",
		"</html>"
	].join("\n");
}

export function buildDownloadFiles(documents: OCRDocumentResponse[], format: DownloadFormat) {
	const { extension, contentType } = FORMAT_CONFIG[format];
	const seenFilenames = new Map<string, number>();
	const files: DownloadFile[] = [];

	for (const document of documents) {
		const content = fileContent(document, format);
		if (content === null) continue;

		const basename = sanitizeDownloadBasename(document.original_filename);
		const filename = uniqueFilename(basename, extension, seenFilenames);
		files.push({ filename, content, contentType });
	}

	return files;
}

export async function buildZip(files: DownloadFile[]) {
	const zip = new JSZip();

	for (const file of files) {
		zip.file(file.filename, file.content);
	}

	return zip.generateAsync({ type: "uint8array" });
}

export function zipDownloadFilename(date = new Date()) {
	return `syncra-${downloadTimestamp(date)}.zip`;
}

export function contentDispositionAttachment(filename: string) {
	const fallback = filename
		.replace(/[^\x20-\x7e]/g, "_")
		.replace(/["\\]/g, "_")
		.trim() || "download";

	return `attachment; filename="${fallback}"; filename*=UTF-8''${encodeURIComponent(filename)}`;
}

function fileContent(document: OCRDocumentResponse, format: DownloadFormat) {
	if (format === "markdown") return document.markdown;
	if (format === "html") return renderStandaloneHTML(document);

	if (!documentHasSchema(document) || document.annotation_json === undefined) return null;

	return JSON.stringify(document.annotation_json, null, 2);
}

function uniqueFilename(basename: string, extension: string, seenFilenames: Map<string, number>) {
	const key = `${basename}${extension}`.toLowerCase();
	const count = seenFilenames.get(key) ?? 0;
	seenFilenames.set(key, count + 1);

	if (count === 0) return `${basename}${extension}`;

	return `${basename}-${count + 1}${extension}`;
}

function downloadTimestamp(date: Date) {
	return date.toISOString().replace(/\.\d{3}Z$/, "Z").replaceAll(":", "-");
}

function escapeHTML(value: string) {
	return value
		.replaceAll("&", "&amp;")
		.replaceAll("<", "&lt;")
		.replaceAll(">", "&gt;")
		.replaceAll('"', "&quot;")
		.replaceAll("'", "&#039;");
}
