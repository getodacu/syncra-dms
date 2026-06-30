<script lang="ts">
	import OTPForm from "$lib/components/otp-form.svelte";
	import type { PageProps } from "./$types";

	let { data, form }: PageProps = $props();
</script>

<div class="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
	<div class="w-full max-w-xs">
		<OTPForm email={data.email} {form} />
	</div>
</div>
