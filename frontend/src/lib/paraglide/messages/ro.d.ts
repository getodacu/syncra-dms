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
export const app_name: (inputs: App_NameInputs) => LocalizedString;
export const app_status: (inputs: App_StatusInputs) => LocalizedString;
export const backend_api: (inputs: Backend_ApiInputs) => LocalizedString;
export const database: (inputs: DatabaseInputs) => LocalizedString;
export const version: (inputs: VersionInputs) => LocalizedString;
export const ready: (inputs: ReadyInputs) => LocalizedString;
export const not_ready: (inputs: Not_ReadyInputs) => LocalizedString;
export const connected: (inputs: ConnectedInputs) => LocalizedString;
export const unavailable: (inputs: UnavailableInputs) => LocalizedString;
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
