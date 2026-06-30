import { describe, expect, test, vi } from 'vitest';

vi.mock('$lib/server/syncra-api', () => ({
	fetchOperationalStatus: vi.fn(async () => ({
		apiBaseUrl: 'http://localhost:8080',
		version: {
			ok: true,
			status: 200,
			data: {
				app: 'Syncra DMS',
				module: 'ai.ro/syncra/dms',
				version: 'test'
			}
		},
		readiness: {
			ok: true,
			status: 200,
			data: {
				status: 'ready'
			}
		}
	}))
}));

import { fetchOperationalStatus } from '$lib/server/syncra-api';
import { load } from './+page.server';

describe('home page server load', () => {
	test('returns backend operational status', async () => {
		const fetch = vi.fn();

		const data = await load({ fetch } as unknown as Parameters<typeof load>[0]);

		expect(fetchOperationalStatus).toHaveBeenCalledWith(fetch);
		expect(data.backend.apiBaseUrl).toBe('http://localhost:8080');
		expect(data.backend.version.ok).toBe(true);
		expect(data.backend.readiness.ok).toBe(true);
	});
});
