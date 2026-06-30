import { describe, expect, it } from "vitest";

import {
	avatarFileError,
	clearedAccountSettingsPendingState,
	passwordConfirmationError,
	shouldApplyAccountSettingsResult,
	trimOptionalValue
} from "./account-settings-utils";

describe("account settings utils", () => {
	it("accepts supported avatar images up to 5 MB", () => {
		const file = new File([new Uint8Array([1])], "avatar.png", { type: "image/png" });
		expect(avatarFileError(file)).toBeNull();
	});

	it("rejects unsupported avatar files", () => {
		const file = new File([new Uint8Array([1])], "avatar.txt", { type: "text/plain" });
		expect(avatarFileError(file)).toBe("Choose a PNG, JPG, GIF, AVIF, APNG, SVG, or WEBP image.");
	});

	it("rejects oversized avatar files", () => {
		const file = new File([new Uint8Array((5 << 20) + 1)], "avatar.png", { type: "image/png" });
		expect(avatarFileError(file)).toBe("Avatar must be 5 MB or smaller.");
	});

	it("validates password confirmation", () => {
		expect(passwordConfirmationError("short", "short")).toBe("Password must be at least 8 characters.");
		expect(passwordConfirmationError("👍👍👍👍", "👍👍👍👍")).toBe(
			"Password must be at least 8 characters."
		);
		expect(passwordConfirmationError("password123", "different")).toBe("Passwords do not match.");
		expect(passwordConfirmationError("password123", "password123")).toBeNull();
	});

	it("guards account settings results by request user identity", () => {
		expect(shouldApplyAccountSettingsResult("user-1", "user-1")).toBe(true);
		expect(shouldApplyAccountSettingsResult("user-1", "user-2")).toBe(false);
		expect(shouldApplyAccountSettingsResult(null, null)).toBe(false);
	});

	it("provides cleared pending state for account resets", () => {
		expect(clearedAccountSettingsPendingState()).toEqual({
			avatarPending: false,
			namePending: false,
			emailPending: false,
			passwordPending: false
		});
	});

	it("trims optional values", () => {
		expect(trimOptionalValue(" Ada ")).toBe("Ada");
		expect(trimOptionalValue(" ")).toBe("");
	});
});
