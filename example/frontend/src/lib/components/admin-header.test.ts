import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";

const adminHeaderSource = () => readFileSync(new URL("./admin-header.svelte", import.meta.url), "utf8");
const normalizeSource = (value: string) => value.replace(/\s+/g, " ").trim();

describe("admin header behavior", () => {
	it("uses Paraglide messages for route titles and theme toggle", () => {
		const source = normalizeSource(adminHeaderSource());

		expect(source).toContain('import { m } from "$lib/paraglide/messages.js";');
		for (const messageCall of [
			"m.admin_nav_invoices()",
			"m.admin_nav_orders()",
			"m.admin_nav_json_recipes()",
			"m.admin_nav_users()",
			"m.admin_nav_user()",
			"m.admin_nav_admin()",
			"m.common_toggle_theme()"
		]) {
			expect(source).toContain(messageCall);
		}
	});
});
