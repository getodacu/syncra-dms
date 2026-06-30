import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";

const source = readFileSync(new URL("./admin-sidebar.svelte", import.meta.url), "utf8");

describe("admin sidebar", () => {
	it("contains only admin navigation entries", () => {
		expect(source).toContain('import { m } from "$lib/paraglide/messages.js";');
		expect(source).toContain("m.admin_nav_users()");
		expect(source).toContain("/admin-portal/users");
		expect(source).toContain("m.admin_nav_invoices()");
		expect(source).toContain("/admin-portal/invoices");
		expect(source).toContain("m.admin_nav_orders()");
		expect(source).toContain("/admin-portal/orders");
		expect(source).toContain("m.admin_nav_json_recipes()");
		expect(source).toContain("/admin-portal/json-recipes");
		expect(source).toContain("m.nav_log_out()");

		for (const userAppEntry of [
			"Quick Create"
		]) {
			expect(source).not.toContain(userAppEntry);
		}
	});
});
