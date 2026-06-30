import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";

const enMessages = JSON.parse(readFileSync(new URL("../../../messages/en.json", import.meta.url), "utf8"));
const roMessages = JSON.parse(readFileSync(new URL("../../../messages/ro.json", import.meta.url), "utf8"));

const navigationKeys = [
	"common_cancel",
	"header_credit_balance_unavailable",
	"header_credits",
	"header_credits_unavailable",
	"nav_account",
	"nav_account_link_conflict",
	"nav_account_link_denied",
	"nav_account_link_failed",
	"nav_account_link_not_configured",
	"nav_account_link_sign_in_again",
	"nav_account_linked",
	"nav_billing_orders",
	"nav_create_job",
	"nav_create_quick_ocr_job",
	"nav_create_schema",
	"nav_credit_usage_history",
	"nav_developer_settings",
	"nav_edit_schema",
	"nav_get_help",
	"nav_jobs",
	"nav_log_out",
	"nav_logout_description",
	"nav_logout_failed",
	"nav_logout_title",
	"nav_new_job",
	"nav_new_schema",
	"nav_no_email_address",
	"nav_notifications",
	"nav_quick_ocr",
	"nav_schemas",
	"admin_nav_invoices",
	"admin_nav_json_recipes",
	"admin_nav_admin",
	"admin_nav_orders",
	"admin_nav_user",
	"admin_nav_users",
	"admin_user_fallback",
	"sidebar_admin_portal",
	"sidebar_switch_space",
	"sidebar_syncra",
	"sidebar_syncra_admin",
	"sidebar_user_space"
];

describe("navigation i18n messages", () => {
	it("defines every navigation message in English and Romanian", () => {
		for (const key of navigationKeys) {
			expect(enMessages[key], `en ${key}`).toEqual(expect.any(String));
			expect(roMessages[key], `ro ${key}`).toEqual(expect.any(String));
		}
	});
});
