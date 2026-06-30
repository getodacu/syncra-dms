<script lang="ts">
	import { browser, dev } from "$app/environment";
	import favicon from "$lib/assets/favicon.png";
	import { QueryClient, QueryClientProvider } from "@tanstack/svelte-query";
	import { onMount, type Component, type Snippet } from "svelte";
	import { ModeWatcher } from "mode-watcher";
	import "../app.css";

	let { children }: { children: Snippet } = $props();
	let QueryDevtools: Component | undefined = $state();

	const queryClient = new QueryClient({
		defaultOptions: {
			queries: {
				enabled: browser,
			},
		},
	});

	if (import.meta.env.DEV && dev && browser) {
		onMount(async () => {
			const { SvelteQueryDevtools } = await import(
				"@tanstack/svelte-query-devtools"
			);
			QueryDevtools = SvelteQueryDevtools;
		});
	}
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>

<QueryClientProvider client={queryClient}>
	<ModeWatcher />
	{@render children()}
	{#if QueryDevtools}
		<QueryDevtools />
	{/if}
</QueryClientProvider>
