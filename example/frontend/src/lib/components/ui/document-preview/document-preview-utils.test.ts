import { describe, expect, it } from "vitest";

import { formatAnnotationJSON, renderHighlightedJSON, renderMarkdown } from "./document-preview-utils";

describe("document preview helpers", () => {
	it("formats JSON annotation objects with two-space indentation", () => {
		expect(formatAnnotationJSON({ paid: true, total: 10 })).toBe(
			JSON.stringify({ paid: true, total: 10 }, null, 2)
		);
	});

	it("uses the existing fallback message for missing JSON annotations", () => {
		expect(formatAnnotationJSON(undefined)).toBe("No JSON annotation available.");
		expect(formatAnnotationJSON(undefined, "Nicio adnotare JSON.")).toBe("Nicio adnotare JSON.");
	});

	it("highlights valid JSON with Highlight.js markup", () => {
		const html = renderHighlightedJSON({ total: 10 });

		expect(html).toContain("hljs-attr");
		expect(html).toContain("&quot;total&quot;");
		expect(html).toContain("hljs-number");
	});

	it("does not emit unsafe raw HTML from JSON string values", () => {
		const html = renderHighlightedJSON({ note: "</span><script>alert(1)</script>" });

		expect(html).not.toContain("<script");
		expect(html).not.toContain("</script>");
		expect(html).toContain("&lt;script&gt;");
	});

	it("renders markdown as HTML", () => {
		const html = renderMarkdown("# Invoice\n\nPaid");

		expect(html).toContain("<h1>Invoice</h1>");
		expect(html).toContain("<p>Paid</p>");
	});

	it("escapes custom empty markdown fallback text", () => {
		expect(renderMarkdown(null, "<empty>")).toBe("<p>&lt;empty&gt;</p>");
	});

	it("renders markdown line breaks with GitHub-flavored markdown settings", () => {
		const html = renderMarkdown("Line one\nLine two");

		expect(html).toContain("Line one<br>Line two");
	});

	it("preserves base64 markdown images in generated HTML", () => {
		const html = renderMarkdown(
			"![Receipt](data:image/png;base64,iVBORw0KGgo= \"Scanned receipt\")"
		);

		expect(html).toContain(
			'<img src="data:image/png;base64,iVBORw0KGgo=" alt="Receipt" title="Scanned receipt">'
		);
	});

	it("preserves safe raw base64 image tags in generated HTML", () => {
		const html = renderMarkdown(
			'<img src="data:image/jpeg;base64,/9j/4AAQSkZJRg==" alt="Receipt" width="640" height="480">'
		);

		expect(html).toBe(
			'<img src="data:image/jpeg;base64,/9j/4AAQSkZJRg==" alt="Receipt" width="640" height="480">'
		);
	});

	it("does not emit unsafe raw HTML from markdown", () => {
		const html = renderMarkdown('<img src="javascript:alert(1)" onerror="alert(2)">');

		expect(html).not.toContain("<img");
		expect(html).not.toContain('src="javascript:');
		expect(html).toContain("&lt;img");
	});

	it("does not emit unsafe markdown link URLs", () => {
		const html = renderMarkdown("[Run](javascript:alert(1))");

		expect(html).not.toContain("<a");
		expect(html).not.toContain("javascript:");
		expect(html).toContain("<p>Run</p>");
	});
});
