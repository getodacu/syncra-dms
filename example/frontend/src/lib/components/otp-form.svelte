<script lang="ts">
	import { enhance } from "$app/forms";
	import { Button } from "$lib/components/ui/button/index.js";
	import * as Card from "$lib/components/ui/card/index.js";
	import * as Field from "$lib/components/ui/field/index.js";
	import { Input } from "$lib/components/ui/input/index.js";
	import * as InputOTP from "$lib/components/ui/input-otp/index.js";
	import type { ComponentProps } from "svelte";
	import type { SubmitFunction } from "@sveltejs/kit";

	type OTPFormState = {
		values?: { email?: string; otp?: string };
		fieldErrors?: Record<string, string>;
		error?: string;
		success?: string;
	};

	let {
		email = "",
		form,
		...props
	}: ComponentProps<typeof Card.Root> & {
		email?: string;
		form?: OTPFormState | null;
	} = $props();

	let otp = $state("");
	let submitting = $state(false);

	const currentEmail = $derived(form?.values?.email ?? email);

	const submit: SubmitFunction = () => {
		submitting = true;
		return async ({ update }) => {
			submitting = false;
			await update({ reset: false });
		};
	};
</script>

<Card.Root {...props}>
	<Card.Header>
		<Card.Title>Enter verification code</Card.Title>
		<Card.Description>
			We sent a 6-digit code to your email. Check your inbox and spam folder.
		</Card.Description>
	</Card.Header>
	<Card.Content>
		<form method="POST" action="?/verify" use:enhance={submit}>
			<Field.Group>
				{#if form?.error}
					<Field.Error>{form.error}</Field.Error>
				{/if}
				{#if form?.success}
					<Field.Description>{form.success}</Field.Description>
				{/if}
				<Field.Field>
					<Field.Label for="email">Email</Field.Label>
					<Input
						id="email"
						name="email"
						type="email"
						value={currentEmail}
						placeholder="m@example.com"
						required
					/>
					{#if form?.fieldErrors?.email}
						<Field.Error>{form.fieldErrors.email}</Field.Error>
					{/if}
				</Field.Field>
				<Field.Field>
					<Field.Label for="otp">Verification code</Field.Label>
					<input type="hidden" name="otp" value={otp} />
					<InputOTP.Root maxlength={6} id="otp" required bind:value={otp}>
						{#snippet children({ cells })}
							<InputOTP.Group
								class="gap-2.5 *:data-[slot=input-otp-slot]:rounded-md *:data-[slot=input-otp-slot]:border"
							>
								{#each cells as cell (cell)}
									<InputOTP.Slot {cell} />
								{/each}
							</InputOTP.Group>
						{/snippet}
					</InputOTP.Root>
					<Field.Description>
						Enter the 6-digit code sent to your email.
					</Field.Description>
					{#if form?.fieldErrors?.otp}
						<Field.Error>{form.fieldErrors.otp}</Field.Error>
					{/if}
				</Field.Field>
				<Field.Group>
					<Button type="submit" disabled={submitting}>
						{submitting ? "Verifying..." : "Verify"}
					</Button>
					<Button
						variant="outline"
						type="submit"
						formaction="?/resend"
						formnovalidate
						disabled={submitting}
					>
						Resend code
					</Button>
				</Field.Group>
			</Field.Group>
		</form>
	</Card.Content>
</Card.Root>
