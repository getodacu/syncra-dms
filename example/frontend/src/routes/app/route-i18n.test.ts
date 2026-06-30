import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";

const enMessages = JSON.parse(readFileSync(new URL("../../../messages/en.json", import.meta.url), "utf8"));
const roMessages = JSON.parse(readFileSync(new URL("../../../messages/ro.json", import.meta.url), "utf8"));

const routeKeys = [
	"common_delete",
	"common_flexible",
	"common_next",
	"common_previous",
	"common_retry",
	"common_rows_per_page",
	"common_strict",
	"schemas_new_title",
	"schemas_library",
	"schemas_edit_title",
	"schemas_delete_single_title",
	"schemas_no_schemas_found",
	"schemas_editor_badge",
	"schemas_schema_name_label",
	"schemas_structure_designer",
	"jobs_page_title",
	"jobs_delete_single_title",
	"jobs_no_jobs_found",
	"jobs_inline_schema",
	"jobs_saved_extraction_schema",
	"new_job_select_schema",
	"new_job_upload_documents",
	"new_job_run_monitor",
	"new_job_ocr_only_mode_active",
	"new_job_pending_upload_queue",
	"new_job_run_extraction_one",
	"new_job_run_extraction_other",
	"billing_purchase_credits",
	"billing_secure_checkout",
	"billing_orders_all_orders",
	"billing_orders_no_orders_found",
	"billing_order_status_paid",
	"credit_usage_all_activity",
	"credit_usage_no_usage_found",
	"account_settings_title",
	"account_settings_save_password",
	"billing_profile_title",
	"billing_profile_save_button"
];

describe("app route i18n messages", () => {
	it("defines route messages in English and Romanian", () => {
		for (const key of routeKeys) {
			expect(enMessages[key], `en ${key}`).toEqual(expect.any(String));
			expect(roMessages[key], `ro ${key}`).toEqual(expect.any(String));
		}
	});
});
