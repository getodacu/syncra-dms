/**
* | output |
* | --- |
* | "Syncra DMS" |
*
* @param {App_NameInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const app_name: ((inputs?: App_NameInputs, options?: {
    locale?: "en" | "ro";
}) => LocalizedString) & import("../runtime.js").MessageMetadata<App_NameInputs, {
    locale?: "en" | "ro";
}, {}>;
/**
* | output |
* | --- |
* | "Environment status" |
*
* @param {App_StatusInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const app_status: ((inputs?: App_StatusInputs, options?: {
    locale?: "en" | "ro";
}) => LocalizedString) & import("../runtime.js").MessageMetadata<App_StatusInputs, {
    locale?: "en" | "ro";
}, {}>;
/**
* | output |
* | --- |
* | "Backend API" |
*
* @param {Backend_ApiInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const backend_api: ((inputs?: Backend_ApiInputs, options?: {
    locale?: "en" | "ro";
}) => LocalizedString) & import("../runtime.js").MessageMetadata<Backend_ApiInputs, {
    locale?: "en" | "ro";
}, {}>;
/**
* | output |
* | --- |
* | "Database" |
*
* @param {DatabaseInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const database: ((inputs?: DatabaseInputs, options?: {
    locale?: "en" | "ro";
}) => LocalizedString) & import("../runtime.js").MessageMetadata<DatabaseInputs, {
    locale?: "en" | "ro";
}, {}>;
/**
* | output |
* | --- |
* | "Version" |
*
* @param {VersionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const version: ((inputs?: VersionInputs, options?: {
    locale?: "en" | "ro";
}) => LocalizedString) & import("../runtime.js").MessageMetadata<VersionInputs, {
    locale?: "en" | "ro";
}, {}>;
/**
* | output |
* | --- |
* | "Ready" |
*
* @param {ReadyInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const ready: ((inputs?: ReadyInputs, options?: {
    locale?: "en" | "ro";
}) => LocalizedString) & import("../runtime.js").MessageMetadata<ReadyInputs, {
    locale?: "en" | "ro";
}, {}>;
/**
* | output |
* | --- |
* | "Not ready" |
*
* @param {Not_ReadyInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const not_ready: ((inputs?: Not_ReadyInputs, options?: {
    locale?: "en" | "ro";
}) => LocalizedString) & import("../runtime.js").MessageMetadata<Not_ReadyInputs, {
    locale?: "en" | "ro";
}, {}>;
/**
* | output |
* | --- |
* | "Connected" |
*
* @param {ConnectedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const connected: ((inputs?: ConnectedInputs, options?: {
    locale?: "en" | "ro";
}) => LocalizedString) & import("../runtime.js").MessageMetadata<ConnectedInputs, {
    locale?: "en" | "ro";
}, {}>;
/**
* | output |
* | --- |
* | "Unavailable" |
*
* @param {UnavailableInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const unavailable: ((inputs?: UnavailableInputs, options?: {
    locale?: "en" | "ro";
}) => LocalizedString) & import("../runtime.js").MessageMetadata<UnavailableInputs, {
    locale?: "en" | "ro";
}, {}>;
export type LocalizedString = import("../runtime.js").LocalizedString;
export type App_NameInputs = {};
export type App_StatusInputs = {};
export type Backend_ApiInputs = {};
export type DatabaseInputs = {};
export type VersionInputs = {};
export type ReadyInputs = {};
export type Not_ReadyInputs = {};
export type ConnectedInputs = {};
export type UnavailableInputs = {};
