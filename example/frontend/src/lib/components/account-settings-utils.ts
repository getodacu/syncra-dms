const MAX_AVATAR_SIZE = 5 << 20;

const SUPPORTED_AVATAR_TYPES = new Set([
	"image/png",
	"image/jpeg",
	"image/jpg",
	"image/gif",
	"image/avif",
	"image/apng",
	"image/svg+xml",
	"image/webp"
]);

export type AccountSettingsPendingState = {
	avatarPending: boolean;
	namePending: boolean;
	emailPending: boolean;
	passwordPending: boolean;
};

export function avatarFileError(file: File) {
	if (!SUPPORTED_AVATAR_TYPES.has(file.type)) {
		return "Choose a PNG, JPG, GIF, AVIF, APNG, SVG, or WEBP image.";
	}

	if (file.size > MAX_AVATAR_SIZE) {
		return "Avatar must be 5 MB or smaller.";
	}

	return null;
}

export function passwordConfirmationError(password: string, confirmation: string) {
	if ([...password].length < 8) {
		return "Password must be at least 8 characters.";
	}

	if (password !== confirmation) {
		return "Passwords do not match.";
	}

	return null;
}

export function trimOptionalValue(value: string) {
	return value.trim();
}

export function shouldApplyAccountSettingsResult(
	requestUserId: string | null | undefined,
	currentUserId: string | null | undefined
) {
	return Boolean(requestUserId) && requestUserId === currentUserId;
}

export function clearedAccountSettingsPendingState(): AccountSettingsPendingState {
	return {
		avatarPending: false,
		namePending: false,
		emailPending: false,
		passwordPending: false
	};
}

export function fileToDataURL(file: File) {
	return new Promise<string>((resolve, reject) => {
		const reader = new FileReader();
		reader.onload = () => resolve(String(reader.result));
		reader.onerror = () => reject(reader.error);
		reader.readAsDataURL(file);
	});
}
