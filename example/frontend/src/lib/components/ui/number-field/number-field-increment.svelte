<script lang="ts">
	import { getContext, type Snippet } from "svelte";
	import { Button } from "$lib/components/ui/button/index.js";
	import PlusIcon from "@lucide/svelte/icons/plus";

	interface Props {
		class?: string;
		variant?: any;
		tabindex?: number | null;
		children?: Snippet;
		[key: string]: any;
	}

	let { class: className, variant = "outline", tabindex = null, children, ...restProps }: Props = $props();

	const context = getContext<{
		value: number;
		step: number;
		min: number;
		max: number;
	}>("number-field");

	function increment() {
		if (context) {
			context.value = context.value + context.step;
		}
	}
</script>

<Button
	type="button"
	onclick={increment}
	disabled={context ? context.value >= context.max : false}
	class={className}
	variant={variant}
	tabindex={tabindex === null ? -1 : tabindex}
	{...restProps}
>
	{#if children}
		{@render children()}
	{:else}
		<PlusIcon class="size-4" />
	{/if}
</Button>
