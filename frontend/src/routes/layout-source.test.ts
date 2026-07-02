import { readFileSync } from 'node:fs';
import { describe, expect, it } from 'vitest';

describe('root layout query provider', () => {
	it('wraps the app with TanStack QueryClientProvider and browser-only queries', () => {
		const source = readFileSync(new URL('./+layout.svelte', import.meta.url), 'utf8');

		expect(source).toContain("import { QueryClient, QueryClientProvider }");
		expect(source).toContain('const queryClient = new QueryClient');
		expect(source).toContain('enabled: browser');
		expect(source).toContain('<QueryClientProvider client={queryClient}>');
		expect(source).toContain('SvelteQueryDevtools');
	});
});

describe('app sidebar layout shell', () => {
	it('uses the adapted shadcn sidebar shell and support overlays', () => {
		const source = readFileSync(new URL('./app/+layout.svelte', import.meta.url), 'utf8');

		expect(source).toContain('<Sidebar.Provider');
		expect(source).toContain(
			'<AppSidebar variant="inset" user={data.user} permissions={data.permissions} />'
		);
		expect(source).toContain('<Sidebar.Inset>');
		expect(source).toContain('<SiteHeader />');
		expect(source).toContain('<ConfirmDeleteDialog />');
		expect(source).toContain('<Toaster position="top-right"');
	});

	it('wires only valid Syncra DMS app navigation routes', () => {
		const source = readFileSync(
			new URL('../lib/components/app-sidebar.svelte', import.meta.url),
			'utf8'
		);

		expect(source).toContain("title: 'Dashboard'");
		expect(source).toContain("url: '/app'");
		expect(source).toContain("title: 'Organization Units'");
		expect(source).toContain("url: '/app/organization-units'");
		expect(source).toContain("title: 'Documents'");
		expect(source).toContain("url: '/app/documents'");
		expect(source).not.toContain('/app/billing');
		expect(source).not.toContain('/app/jobs');
		expect(source).not.toContain('/app/datasets');
	});

	it('includes a route title and theme toggle in the site header', () => {
		const source = readFileSync(
			new URL('../lib/components/site-header.svelte', import.meta.url),
			'utf8'
		);

		expect(source).toContain('Sidebar.Trigger');
		expect(source).toContain('toggleMode');
		expect(source).toContain('Toggle theme');
		expect(source).toContain("if (pathname === '/app/organization-units')");
		expect(source).toContain("if (pathname === '/app/documents')");
	});
});
