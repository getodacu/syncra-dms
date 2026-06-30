<script lang="ts">
	import { setContext, type Snippet } from "svelte";

	interface Props {
		value: number;
		step?: number;
		min?: number;
		max?: number;
		children?: Snippet;
	}

	let {
		value = $bindable(0),
		step = 1,
		min = -Infinity,
		max = Infinity,
		children
	}: Props = $props();

	// Create a stateful context object using getters/setters for Svelte 5 runes reactivity
	const context = {
		get value() {
			return value;
		},
		set value(v: number) {
			value = Math.min(max, Math.max(min, v));
		},
		get step() {
			return step;
		},
		get min() {
			return min;
		},
		get max() {
			return max;
		}
	};

	setContext("number-field", context);
</script>

{@render children?.()}
