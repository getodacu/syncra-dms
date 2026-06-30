import { env } from '$env/dynamic/private';

export const INTERNAL_API_HEADER = 'X-Syncra-Internal-Token';

export function apiBaseUrl() {
	return (privateEnv('SYNCRA_API_BASE_URL') || 'http://localhost:8080').replace(/\/+$/, '');
}

export function privateEnv(key: string) {
	return process.env[key] || env[key];
}

export function internalAPIHeaders(headers: HeadersInit = {}) {
	const token = (privateEnv('SYNCRA_INTERNAL_API_TOKEN') || '').trim();
	if (!token) return null;

	const output = new Headers(headers);
	output.set(INTERNAL_API_HEADER, token);
	return output;
}
