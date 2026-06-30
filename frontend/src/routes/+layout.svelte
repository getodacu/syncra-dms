<script lang="ts">
	import { browser, dev } from '$app/environment';
	import { ModeWatcher } from 'mode-watcher';
	import { QueryClient, QueryClientProvider } from '@tanstack/svelte-query';
	import { onMount, type Component, type Snippet } from 'svelte';
	import '../app.css';

	let { children }: { children: Snippet } = $props();
	let QueryDevtools: Component | undefined = $state();

	const queryClient = new QueryClient({
		defaultOptions: {
			queries: {
				enabled: browser
			}
		}
	});

	if (import.meta.env.DEV && dev && browser) {
		onMount(async () => {
			const { SvelteQueryDevtools } = await import('@tanstack/svelte-query-devtools');
			QueryDevtools = SvelteQueryDevtools;
		});
	}
</script>

<svelte:head>
	<title>Syncra DMS</title>
	<meta
		name="description"
		content="Syncra DMS work environment for document management workflows."
	/>
</svelte:head>

<QueryClientProvider client={queryClient}>
	<ModeWatcher />
	{@render children()}
	{#if QueryDevtools}
		<QueryDevtools />
	{/if}
</QueryClientProvider>
