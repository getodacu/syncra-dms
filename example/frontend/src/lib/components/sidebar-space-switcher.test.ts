import { existsSync, readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";

const componentUrl = new URL("./sidebar-space-switcher.svelte", import.meta.url);
const source = () => (existsSync(componentUrl) ? readFileSync(componentUrl, "utf8") : "");
const appSidebarSource = () => readFileSync(new URL("./app-sidebar.svelte", import.meta.url), "utf8");
const adminSidebarSource = () => readFileSync(new URL("./admin-sidebar.svelte", import.meta.url), "utf8");
const normalizeSource = (value: string) => value.replace(/\s+/g, " ").trim();

describe("sidebar space switcher", () => {
	it("offers both spaces to admin users", () => {
		const component = normalizeSource(source());

		expect(component).toContain('import { m } from "$lib/paraglide/messages.js";');
		expect(component).toContain('user?.role === "admin"');
		expect(component).toContain("m.sidebar_user_space()");
		expect(component).toContain('href="/app"');
		expect(component).toContain("m.sidebar_admin_portal()");
		expect(component).toContain('href="/admin-portal"');
		expect(component).toContain("m.sidebar_switch_space()");
	});

	it("marks the current admin destination as active", () => {
		const component = normalizeSource(source());

		expect(component).toContain('activeSpace === "app" ? "page" : undefined');
		expect(component).toContain('activeSpace === "admin" ? "page" : undefined');
	});

	it("preserves a static brand link for non-admin users", () => {
		const component = normalizeSource(source());

		expect(component).toContain("{:else}");
		expect(component).toContain("brandHref");
		expect(component).toContain("brandTitle");
		expect(component).toContain("m.sidebar_syncra_admin()");
		expect(component).toContain("m.sidebar_syncra()");
		expect(component).toContain('<span class="text-base font-semibold">{brandTitle}</span>');
	});

	it("is mounted by both sidebar variants", () => {
		const appSidebar = normalizeSource(appSidebarSource());
		const adminSidebar = normalizeSource(adminSidebarSource());

		expect(appSidebar).toContain('import SidebarSpaceSwitcher from "./sidebar-space-switcher.svelte";');
		expect(appSidebar).toContain('<SidebarSpaceSwitcher user={user} activeSpace="app" />');
		expect(adminSidebar).toContain('import SidebarSpaceSwitcher from "./sidebar-space-switcher.svelte";');
		expect(adminSidebar).toContain('<SidebarSpaceSwitcher user={user} activeSpace="admin" />');
	});
});
