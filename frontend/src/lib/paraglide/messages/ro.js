/* eslint-disable */
/** @typedef {import('../runtime.js').LocalizedString} LocalizedString */
/** @typedef {{}} App_NameInputs */
/** @typedef {{}} App_StatusInputs */
/** @typedef {{}} Backend_ApiInputs */
/** @typedef {{}} DatabaseInputs */
/** @typedef {{}} VersionInputs */
/** @typedef {{}} ReadyInputs */
/** @typedef {{}} Not_ReadyInputs */
/** @typedef {{}} ConnectedInputs */
/** @typedef {{}} UnavailableInputs */


export const app_name = /** @type {(inputs: App_NameInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Syncra DMS`)
};

export const app_status = /** @type {(inputs: App_StatusInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Stare mediu`)
};

export const backend_api = /** @type {(inputs: Backend_ApiInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`API backend`)
};

export const database = /** @type {(inputs: DatabaseInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Bază de date`)
};

export const version = /** @type {(inputs: VersionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Versiune`)
};

export const ready = /** @type {(inputs: ReadyInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Pregătit`)
};

export const not_ready = /** @type {(inputs: Not_ReadyInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nepregătit`)
};

export const connected = /** @type {(inputs: ConnectedInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Conectat`)
};

export const unavailable = /** @type {(inputs: UnavailableInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Indisponibil`)
};
