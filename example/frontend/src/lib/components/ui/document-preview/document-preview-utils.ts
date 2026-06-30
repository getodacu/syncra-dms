import DOMPurify from "dompurify";
import hljs from "highlight.js/lib/core";
import json from "highlight.js/lib/languages/json";
import { marked } from "marked";

const NO_JSON_ANNOTATION_MESSAGE = "No JSON annotation available.";
const NO_MARKDOWN_CONTENT_MESSAGE = "No markdown content available.";
const SAFE_URL_PROTOCOL_PATTERN = /^(https?|ftps?|mailto|tel|callto|sms|cid|xmpp|matrix):/i;
const URL_PROTOCOL_PATTERN = /^[a-z][a-z0-9+.-]*:/i;
const BASE64_IMAGE_SRC_PATTERN =
	/^data:image\/(?:png|jpe?g|gif|webp|bmp|avif);base64,[a-z0-9+/=\s]+$/i;
const RAW_IMAGE_TAG_PATTERN = /^<img\s+([^>]*?)\/?>$/i;
const RAW_IMAGE_ATTRIBUTE_PATTERN =
	/([a-zA-Z][\w:-]*)\s*=\s*(?:"([^"]*)"|'([^']*)'|([^\s"'>]+))/g;
const IMAGE_SIZE_PATTERN = /^\d{1,5}(?:\.\d+)?%?$/;
const MARKDOWN_SANITIZE_CONFIG = {
	ADD_ATTR: ["src", "alt", "title", "width", "height"],
	ADD_DATA_URI_TAGS: ["img"],
	ADD_TAGS: ["img"],
	USE_PROFILES: { html: true },
};

hljs.registerLanguage("json", json);

export function formatAnnotationJSON(value: unknown, emptyJSONLabel = NO_JSON_ANNOTATION_MESSAGE) {
	if (value === undefined) return emptyJSONLabel;

	try {
		return JSON.stringify(value, null, 2) ?? String(value);
	} catch {
		return String(value);
	}
}

export function renderHighlightedJSON(value: unknown, emptyJSONLabel = NO_JSON_ANNOTATION_MESSAGE) {
	const formattedJSON = formatAnnotationJSON(value, emptyJSONLabel);

	try {
		const highlightedHTML = hljs.highlight(formattedJSON, { language: "json" }).value;

		return sanitizeHTML(highlightedHTML, {
			ALLOWED_ATTR: ["class"],
			ALLOWED_TAGS: ["span"],
		});
	} catch {
		return escapeHTML(formattedJSON);
	}
}

export function renderMarkdown(value: string | null | undefined, emptyMarkdownLabel = NO_MARKDOWN_CONTENT_MESSAGE) {
	if (!value) return `<p>${escapeHTML(emptyMarkdownLabel)}</p>`;

	try {
		const sanitizer = getHTMLSanitizer();

		const html = marked.parse(value, {
			async: false,
			breaks: true,
			gfm: true,
			renderer: sanitizer ? undefined : createFallbackMarkdownRenderer(),
		}) as string;

		return sanitizer ? sanitizer(html, MARKDOWN_SANITIZE_CONFIG) : html;
	} catch {
		return `<pre>${escapeHTML(value)}</pre>`;
	}
}

function sanitizeHTML(html: string, config?: unknown) {
	const sanitizer = getHTMLSanitizer();
	if (!sanitizer) return html;

	return sanitizer(html, config);
}

function getHTMLSanitizer() {
	const purifier = DOMPurify as unknown as {
		sanitize?: (dirty: string, config?: unknown) => string;
	};

	return typeof purifier.sanitize === "function" ? purifier.sanitize.bind(purifier) : null;
}

function createFallbackMarkdownRenderer() {
	const renderer = new marked.Renderer();

	renderer.html = ({ text }) => renderSafeRawImageTag(text) ?? escapeHTML(text);
	renderer.link = ({ href, title, tokens }) => {
		const text = renderer.parser.parseInline(tokens);
		if (!isSafeURL(href)) return text;

		return `<a href="${escapeHTMLAttribute(href)}"${formatOptionalAttribute("title", title)}>${text}</a>`;
	};
	renderer.image = ({ href, title, text }) => {
		if (!isSafeImageURL(href)) return escapeHTML(text);

		return `<img src="${escapeHTMLAttribute(href)}" alt="${escapeHTMLAttribute(text)}"${formatOptionalAttribute("title", title)}>`;
	};

	return renderer;
}

function renderSafeRawImageTag(value: string) {
	const match = value.trim().match(RAW_IMAGE_TAG_PATTERN);
	if (!match) return null;

	const attributes = parseRawHTMLAttributes(match[1]);
	const src = attributes.get("src");
	if (!src || !isSafeImageURL(src)) return null;

	const parts = [`src="${escapeHTMLAttribute(src)}"`];
	parts.push(`alt="${escapeHTMLAttribute(attributes.get("alt") ?? "")}"`);

	const title = attributes.get("title");
	if (title) parts.push(`title="${escapeHTMLAttribute(title)}"`);

	for (const attributeName of ["width", "height"]) {
		const attributeValue = attributes.get(attributeName);
		if (attributeValue && IMAGE_SIZE_PATTERN.test(attributeValue)) {
			parts.push(`${attributeName}="${escapeHTMLAttribute(attributeValue)}"`);
		}
	}

	return `<img ${parts.join(" ")}>`;
}

function parseRawHTMLAttributes(value: string) {
	const attributes = new Map<string, string>();

	for (const match of value.matchAll(RAW_IMAGE_ATTRIBUTE_PATTERN)) {
		attributes.set(match[1].toLowerCase(), match[2] ?? match[3] ?? match[4] ?? "");
	}

	return attributes;
}

function formatOptionalAttribute(name: string, value: string | null | undefined) {
	return value ? ` ${name}="${escapeHTMLAttribute(value)}"` : "";
}

function isSafeImageURL(value: string) {
	return BASE64_IMAGE_SRC_PATTERN.test(value.trim()) || isSafeURL(value);
}

function isSafeURL(value: string) {
	const trimmedValue = value.trim();

	if (SAFE_URL_PROTOCOL_PATTERN.test(trimmedValue)) return true;
	if (URL_PROTOCOL_PATTERN.test(trimmedValue)) return false;

	return true;
}

function escapeHTML(value: string) {
	return value
		.replaceAll("&", "&amp;")
		.replaceAll("<", "&lt;")
		.replaceAll(">", "&gt;")
		.replaceAll('"', "&quot;")
		.replaceAll("'", "&#039;");
}

function escapeHTMLAttribute(value: string) {
	return escapeHTML(value).replaceAll("`", "&#096;");
}
