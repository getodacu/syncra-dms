import { fetchOperationalStatus } from '$lib/server/syncra-api';

export async function load({ fetch }) {
	return {
		backend: await fetchOperationalStatus(fetch)
	};
}
