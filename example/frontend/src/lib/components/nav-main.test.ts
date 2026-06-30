import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";

const navMainSource = () => readFileSync(new URL("./nav-main.svelte", import.meta.url), "utf8");
const appSidebarSource = () => readFileSync(new URL("./app-sidebar.svelte", import.meta.url), "utf8");
const navSecondarySource = () => readFileSync(new URL("./nav-secondary.svelte", import.meta.url), "utf8");
const normalizeSource = (value: string) => value.replace(/\s+/g, " ").trim();

describe("main navigation behavior", () => {
	it("places developer settings in the bottom secondary navigation", () => {
		const navMain = normalizeSource(navMainSource());
		const appSidebar = normalizeSource(appSidebarSource());

		expect(navMain).not.toContain("Developer Settings");
		expect(navMain).not.toContain("/app/developer-settings");
		expect(appSidebar).toContain("title: m.nav_developer_settings()");
		expect(appSidebar).toContain('url: "/app/developer-settings"');
	});

	it("uses Paraglide messages for app navigation labels", () => {
		const navMain = normalizeSource(navMainSource());
		const appSidebar = normalizeSource(appSidebarSource());
		const navSecondary = normalizeSource(navSecondarySource());

		expect(appSidebar).toContain('import { m } from "$lib/paraglide/messages.js";');
		for (const messageCall of [
			"m.nav_dashboard()",
			"m.nav_schemas()",
			"m.nav_jobs()",
			"m.nav_billing()",
			"m.nav_developer_settings()",
			"m.nav_get_help()"
		]) {
			expect(appSidebar).toContain(messageCall);
		}

		expect(navMain).toContain('import { m } from "$lib/paraglide/messages.js";');
		expect(navMain).toContain("const quickOCRLabel = m.nav_quick_ocr();");
		expect(navMain).toContain("tooltipContent={quickOCRLabel}");
		expect(navMain).toContain("<span>{quickOCRLabel}</span>");
		expect(navMain).toContain("m.nav_create_quick_ocr_job()");
		expect(appSidebar).toContain("m.nav_create_schema()");
		expect(appSidebar).toContain("m.nav_create_job()");
		expect(navMain).toContain("title={plus.title}");
		expect(navSecondary).not.toContain('"Developer Settings"');
		expect(navSecondary).not.toContain('"Get Help"');
	});
});
