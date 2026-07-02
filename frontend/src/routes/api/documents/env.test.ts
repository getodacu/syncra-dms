import { readFileSync } from 'node:fs';
import { describe, expect, it } from 'vitest';

const minimumDocumentUploadLimitBytes = 26 * 1024 * 1024;

describe('document upload environment defaults', () => {
	it('configures adapter-node body limit above the Go document upload limit', () => {
		const envExample = readFileSync(new URL('../../../../.env.example', import.meta.url), 'utf8');
		const bodySizeLimit = envExample
			.split('\n')
			.find((line) => line.trim().startsWith('BODY_SIZE_LIMIT='))
			?.split('=')[1]
			?.trim();

		expect(bodySizeLimit).toBeTruthy();
		expect(parseBodySizeLimit(bodySizeLimit ?? '')).toBeGreaterThanOrEqual(
			minimumDocumentUploadLimitBytes
		);
	});
});

function parseBodySizeLimit(value: string) {
	const match = value.match(/^(\d+)([KMG])?$/i);
	if (!match) return Number.NaN;

	const amount = Number(match[1]);
	const unit = match[2]?.toUpperCase();
	if (unit === 'K') return amount * 1024;
	if (unit === 'M') return amount * 1024 * 1024;
	if (unit === 'G') return amount * 1024 * 1024 * 1024;
	return amount;
}
