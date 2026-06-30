import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";

const source = () => readFileSync(new URL("./otp-form.svelte", import.meta.url), "utf8");

describe("OTPForm", () => {
	it("lets resend submit without native OTP validation", () => {
		const resendButton = source().match(
			/<Button[^>]*formaction="\?\/resend"[^>]*>[\s\S]*?Resend code[\s\S]*?<\/Button>/
		)?.[0];

		expect(resendButton).toContain("formnovalidate");
	});
});
