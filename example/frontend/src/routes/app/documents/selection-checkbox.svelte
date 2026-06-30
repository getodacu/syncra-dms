<script lang="ts">
	import { Checkbox } from "$lib/components/ui/checkbox/index.js";

	let {
		checked = false,
		indeterminate = false,
		ariaLabel,
		onCheckedChange,
	}: {
		checked?: boolean;
		indeterminate?: boolean;
		ariaLabel: string;
		onCheckedChange: (checked: boolean) => void;
	} = $props();

	function handleCheckedChange(value: boolean | "indeterminate") {
		onCheckedChange(value === true);
	}
</script>

<Checkbox
	bind:checked={() => checked, handleCheckedChange}
	{indeterminate}
	aria-label={ariaLabel}
/>
