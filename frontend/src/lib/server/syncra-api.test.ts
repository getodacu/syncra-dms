import { describe, expect, test, vi } from 'vitest';
import { fetchOperationalStatus, normalizeApiBaseUrl } from './syncra-api';

describe('normalizeApiBaseUrl', () => {
	test('uses the local Go API when unset', () => {
		expect(normalizeApiBaseUrl('')).toBe('http://localhost:8080');
	});

	test('removes trailing slashes from configured URLs', () => {
		expect(normalizeApiBaseUrl('http://localhost:9090///')).toBe('http://localhost:9090');
	});
});

describe('fetchOperationalStatus', () => {
	test('loads backend version and readiness through server-side fetch', async () => {
		const fetch = vi.fn(async (input: string | URL | Request) => {
			const url = String(input);
			if (url.endsWith('/version')) {
				return jsonResponse(200, {
					app: 'Syncra DMS',
					module: 'ai.ro/syncra/dms',
					version: 'test'
				});
			}
			if (url.endsWith('/readyz')) {
				return jsonResponse(200, {
					status: 'ready'
				});
			}
			throw new Error(`unexpected URL ${url}`);
		});

		const status = await fetchOperationalStatus(fetch, 'http://localhost:8080/');

		expect(fetch).toHaveBeenCalledWith('http://localhost:8080/version');
		expect(fetch).toHaveBeenCalledWith('http://localhost:8080/readyz');
		expect(status.apiBaseUrl).toBe('http://localhost:8080');
		expect(status.version).toEqual({
			ok: true,
			status: 200,
			data: {
				app: 'Syncra DMS',
				module: 'ai.ro/syncra/dms',
				version: 'test'
			}
		});
		expect(status.readiness).toEqual({
			ok: true,
			status: 200,
			data: {
				status: 'ready'
			}
		});
	});

	test('returns endpoint errors without throwing from the page load path', async () => {
		const fetch = vi.fn(async (input: string | URL | Request) => {
			const url = String(input);
			if (url.endsWith('/version')) {
				throw new Error('connection refused');
			}
			return jsonResponse(503, {
				status: 'not_ready',
				error: 'database unavailable'
			});
		});

		const status = await fetchOperationalStatus(fetch, 'http://localhost:8080');

		expect(status.version).toEqual({
			ok: false,
			error: 'connection refused'
		});
		expect(status.readiness).toEqual({
			ok: false,
			status: 503,
			data: {
				status: 'not_ready',
				error: 'database unavailable'
			},
			error: 'database unavailable'
		});
	});
});

function jsonResponse(status: number, body: unknown): Response {
	return new Response(JSON.stringify(body), {
		status,
		headers: {
			'content-type': 'application/json'
		}
	});
}
