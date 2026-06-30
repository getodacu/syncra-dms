<script lang="ts">
	import OcrRecipeLibrary from "$lib/ocr-recipes/ocr-recipe-library.svelte";
	import { m } from "$lib/paraglide/messages.js";

	type PageData = import("./$types").PageData;

	let { data }: { data: PageData } = $props();
</script>

<svelte:head>
	<title>{m.schemas_library()} | Syncra</title>
</svelte:head>

<div class="@container/main flex flex-1 flex-col gap-6 p-4 lg:p-6">
	<OcrRecipeLibrary
		recipes={data.recipes}
		loadError={data.loadError}
		isLoggedIn={data.isLoggedIn}
		userId={data.userId}
	/>
</div>
