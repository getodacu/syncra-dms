import { env } from '$env/dynamic/private';

const defaultApiBaseUrl = 'http://localhost:8080';

export type BackendVersion = {
	app: string;
	module: string;
	version: string;
};

export type BackendReadiness = {
	status: string;
	error?: string;
};

export type EndpointResult<T> =
	| {
			ok: true;
			status: number;
			data: T;
	  }
	| {
			ok: false;
			status?: number;
			data?: T;
			error: string;
	  };

export type OperationalStatus = {
	apiBaseUrl: string;
	version: EndpointResult<BackendVersion>;
	readiness: EndpointResult<BackendReadiness>;
};

type ServerFetch = (input: string | URL | Request, init?: RequestInit) => Promise<Response>;

export function normalizeApiBaseUrl(value: string | undefined): string {
	const trimmed = value?.trim() ?? '';
	return (trimmed || defaultApiBaseUrl).replace(/\/+$/, '');
}

export async function fetchOperationalStatus(
	fetchFn: ServerFetch,
	baseUrl = env.SYNCRA_API_BASE_URL
): Promise<OperationalStatus> {
	const apiBaseUrl = normalizeApiBaseUrl(baseUrl);
	const [version, readiness] = await Promise.all([
		fetchJSON<BackendVersion>(fetchFn, `${apiBaseUrl}/version`),
		fetchJSON<BackendReadiness>(fetchFn, `${apiBaseUrl}/readyz`)
	]);

	return {
		apiBaseUrl,
		version,
		readiness
	};
}

async function fetchJSON<T>(fetchFn: ServerFetch, url: string): Promise<EndpointResult<T>> {
	try {
		const response = await fetchFn(url);
		const data = (await response.json().catch(() => undefined)) as T | undefined;

		if (response.ok && data !== undefined) {
			return {
				ok: true,
				status: response.status,
				data
			};
		}

		return {
			ok: false,
			status: response.status,
			data,
			error: responseError(response, data)
		};
	} catch (error) {
		return {
			ok: false,
			error: error instanceof Error ? error.message : String(error)
		};
	}
}

function responseError(response: Response, data: unknown): string {
	if (data && typeof data === 'object') {
		const record = data as Record<string, unknown>;
		if (typeof record.error === 'string' && record.error !== '') {
			return record.error;
		}
		if (typeof record.message === 'string' && record.message !== '') {
			return record.message;
		}
	}
	return response.statusText || `HTTP ${response.status}`;
}
