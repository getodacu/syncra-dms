/* eslint-disable */
import { getLocale, experimentalStaticLocale } from "../runtime.js"

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
import * as __en from "./en.js"
import * as __ro from "./ro.js"
/**
* | output |
* | --- |
* | "Syncra DMS" |
*
* @param {App_NameInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const app_name = /** @type {((inputs?: App_NameInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<App_NameInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.app_name(inputs)
	return __ro.app_name(inputs)
});
/**
* | output |
* | --- |
* | "Environment status" |
*
* @param {App_StatusInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const app_status = /** @type {((inputs?: App_StatusInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<App_StatusInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.app_status(inputs)
	return __ro.app_status(inputs)
});
/**
* | output |
* | --- |
* | "Backend API" |
*
* @param {Backend_ApiInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const backend_api = /** @type {((inputs?: Backend_ApiInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Backend_ApiInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.backend_api(inputs)
	return __ro.backend_api(inputs)
});
/**
* | output |
* | --- |
* | "Database" |
*
* @param {DatabaseInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const database = /** @type {((inputs?: DatabaseInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<DatabaseInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.database(inputs)
	return __ro.database(inputs)
});
/**
* | output |
* | --- |
* | "Version" |
*
* @param {VersionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const version = /** @type {((inputs?: VersionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<VersionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.version(inputs)
	return __ro.version(inputs)
});
/**
* | output |
* | --- |
* | "Ready" |
*
* @param {ReadyInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const ready = /** @type {((inputs?: ReadyInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<ReadyInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.ready(inputs)
	return __ro.ready(inputs)
});
/**
* | output |
* | --- |
* | "Not ready" |
*
* @param {Not_ReadyInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const not_ready = /** @type {((inputs?: Not_ReadyInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Not_ReadyInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.not_ready(inputs)
	return __ro.not_ready(inputs)
});
/**
* | output |
* | --- |
* | "Connected" |
*
* @param {ConnectedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const connected = /** @type {((inputs?: ConnectedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<ConnectedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.connected(inputs)
	return __ro.connected(inputs)
});
/**
* | output |
* | --- |
* | "Unavailable" |
*
* @param {UnavailableInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const unavailable = /** @type {((inputs?: UnavailableInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<UnavailableInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.unavailable(inputs)
	return __ro.unavailable(inputs)
});
