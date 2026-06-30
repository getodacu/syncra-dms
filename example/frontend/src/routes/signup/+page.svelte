<script lang="ts">
	import SignupForm from "$lib/components/signup-form.svelte";
	import type { PageProps } from "./$types";

	let { form }: PageProps = $props();
</script>

<div class="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
	<div class="w-full max-w-sm">
		<SignupForm {form} />
	</div>
</div>
