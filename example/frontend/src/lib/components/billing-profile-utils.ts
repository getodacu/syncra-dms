export type BillingProfileForm = {
	entity_type: "individual" | "company";
	billing_name: string;
	billing_email: string;
	country_code: string;
	address_line1: string;
	address_line2: string;
	city: string;
	region: string;
	postal_code: string;
	fiscal_code: string;
	registration_number: string;
};

export type BillingProfileAPIValue = {
	entity_type?: "individual" | "company";
	billing_name?: string;
	billing_email?: string;
	country_code?: string;
	address_line1?: string;
	address_line2?: string;
	city?: string;
	region?: string;
	postal_code?: string;
	fiscal_code?: string;
	registration_number?: string;
};

type BillingProfileFormValue = Omit<BillingProfileAPIValue, "entity_type"> & {
	entity_type?: unknown;
};

type AccountUserValue = {
	name?: string | null;
	email?: string | null;
};

const ISO_COUNTRY_CODES = new Set(
	"AD AE AF AG AI AL AM AO AQ AR AS AT AU AW AX AZ BA BB BD BE BF BG BH BI BJ BL BM BN BO BQ BR BS BT BV BW BY BZ CA CC CD CF CG CH CI CK CL CM CN CO CR CU CV CW CX CY CZ DE DJ DK DM DO DZ EC EE EG EH ER ES ET FI FJ FK FM FO FR GA GB GD GE GF GG GH GI GL GM GN GP GQ GR GS GT GU GW GY HK HM HN HR HT HU ID IE IL IM IN IO IQ IR IS IT JE JM JO JP KE KG KH KI KM KN KP KR KW KY KZ LA LB LC LI LK LR LS LT LU LV LY MA MC MD ME MF MG MH MK ML MM MN MO MP MQ MR MS MT MU MV MW MX MY MZ NA NC NE NF NG NI NL NO NP NR NU NZ OM PA PE PF PG PH PK PL PM PN PR PS PT PW PY QA RE RO RS RU RW SA SB SC SD SE SG SH SI SJ SK SL SM SN SO SR SS ST SV SX SY SZ TC TD TF TG TH TJ TK TL TM TN TO TR TT TV TW TZ UA UG UM US UY UZ VA VC VE VG VI VN VU WF WS YE YT ZA ZM ZW".split(
		" "
	)
);

const MAX_LENGTHS = {
	billing_name: 255,
	billing_email: 320,
	address_line1: 255,
	address_line2: 255,
	city: 160,
	region: 160,
	postal_code: 40,
	fiscal_code: 80,
	registration_number: 120
};

function trimmed(value: unknown) {
	return typeof value === "string" ? value.trim() : "";
}

function exceedsMaxLength(value: string, maxLength: number) {
	return Array.from(value).length > maxLength;
}

function maxLengthError(label: string, value: string, maxLength: number) {
	return exceedsMaxLength(value, maxLength)
		? `${label} must be ${maxLength} characters or fewer.`
		: null;
}

function validCountryCode(value: string) {
	return ISO_COUNTRY_CODES.has(value);
}

export function defaultBillingProfileForm(user: AccountUserValue = {}): BillingProfileForm {
	return {
		entity_type: "individual",
		billing_name: trimmed(user.name),
		billing_email: trimmed(user.email),
		country_code: "RO",
		address_line1: "",
		address_line2: "",
		city: "",
		region: "",
		postal_code: "",
		fiscal_code: "",
		registration_number: ""
	};
}

export function normalizeBillingProfileForm(form: BillingProfileFormValue): BillingProfileForm {
	return {
		entity_type: form.entity_type === "company" ? "company" : "individual",
		billing_name: trimmed(form.billing_name),
		billing_email: trimmed(form.billing_email),
		country_code: trimmed(form.country_code).toUpperCase(),
		address_line1: trimmed(form.address_line1),
		address_line2: trimmed(form.address_line2),
		city: trimmed(form.city),
		region: trimmed(form.region),
		postal_code: trimmed(form.postal_code),
		fiscal_code: trimmed(form.fiscal_code),
		registration_number: trimmed(form.registration_number)
	};
}

export function validateBillingProfileForm(form: BillingProfileFormValue): string | null {
	const normalized = normalizeBillingProfileForm(form);

	if (!normalized.billing_name) {
		return normalized.entity_type === "company"
			? "Company name is required."
			: "Full name is required.";
	}
	const nameLengthError = maxLengthError(
		normalized.entity_type === "company" ? "Company name" : "Full name",
		normalized.billing_name,
		MAX_LENGTHS.billing_name
	);
	if (nameLengthError) return nameLengthError;

	const emailLengthError = maxLengthError(
		"Billing email",
		normalized.billing_email,
		MAX_LENGTHS.billing_email
	);
	if (emailLengthError) return emailLengthError;
	if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(normalized.billing_email)) {
		return "Billing email is invalid.";
	}

	if (!normalized.country_code) {
		return "Country is required.";
	}

	if (!validCountryCode(normalized.country_code)) {
		return "Country code is invalid.";
	}

	if (!normalized.address_line1) {
		return "Address line 1 is required.";
	}
	const addressLine1LengthError = maxLengthError(
		"Address line 1",
		normalized.address_line1,
		MAX_LENGTHS.address_line1
	);
	if (addressLine1LengthError) return addressLine1LengthError;

	const addressLine2LengthError = maxLengthError(
		"Address line 2",
		normalized.address_line2,
		MAX_LENGTHS.address_line2
	);
	if (addressLine2LengthError) return addressLine2LengthError;

	if (!normalized.city) {
		return "City is required.";
	}
	const cityLengthError = maxLengthError("City", normalized.city, MAX_LENGTHS.city);
	if (cityLengthError) return cityLengthError;

	const regionLengthError = maxLengthError("Region/state", normalized.region, MAX_LENGTHS.region);
	if (regionLengthError) return regionLengthError;

	if (!normalized.postal_code) {
		return "Postal code is required.";
	}
	const postalCodeLengthError = maxLengthError(
		"Postal code",
		normalized.postal_code,
		MAX_LENGTHS.postal_code
	);
	if (postalCodeLengthError) return postalCodeLengthError;

	if (
		normalized.entity_type === "company" &&
		normalized.country_code === "RO" &&
		!normalized.fiscal_code
	) {
		return "Fiscal code is required for Romanian companies.";
	}
	const fiscalCodeLengthError = maxLengthError(
		"Fiscal code",
		normalized.fiscal_code,
		MAX_LENGTHS.fiscal_code
	);
	if (fiscalCodeLengthError) return fiscalCodeLengthError;

	const registrationNumberLengthError = maxLengthError(
		"Registration number",
		normalized.registration_number,
		MAX_LENGTHS.registration_number
	);
	if (registrationNumberLengthError) return registrationNumberLengthError;

	return null;
}

export function formFromBillingProfile(
	profile: BillingProfileAPIValue & Record<string, unknown>
): BillingProfileForm {
	return normalizeBillingProfileForm(profile);
}
