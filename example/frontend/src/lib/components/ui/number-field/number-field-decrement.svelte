<script lang="ts">
	import { getContext, type Snippet } from "svelte";
	import { Button } from "$lib/components/ui/button/index.js";
	import MinusIcon from "@lucide/svelte/icons/minus";

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

	function decrement() {
		if (context) {
			context.value = context.value - context.step;
		}
	}
</script>

<Button
	type="button"
	onclick={decrement}
	disabled={context ? context.value <= context.min : false}
	class={className}
	variant={variant}
	tabindex={tabindex === null ? -1 : tabindex}
	{...restProps}
>
	{#if children}
		{@render children()}
	{:else}
		<MinusIcon class="size-4" />
	{/if}
</Button>
