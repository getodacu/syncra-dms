<script lang="ts">
	import { ModeWatcher } from 'mode-watcher';
	import type { Snippet } from 'svelte';
	import '../app.css';

	let { children }: { children: Snippet } = $props();
</script>

<svelte:head>
	<title>Syncra DMS</title>
	<meta
		name="description"
		content="Syncra DMS work environment for document management workflows."
	/>
</svelte:head>

<ModeWatcher />
{@render children()}
