<script lang="ts">
	import { enhance } from "$app/forms";
	import { Button } from "$lib/components/ui/button/index.js";
	import * as Card from "$lib/components/ui/card/index.js";
	import * as Field from "$lib/components/ui/field/index.js";
	import { Input } from "$lib/components/ui/input/index.js";
	import type { SubmitFunction } from "@sveltejs/kit";
	import type { PageProps } from "./$types";

	type RecoverPasswordFormState = {
		values?: { email?: string; token?: string };
		fieldErrors?: Record<string, string>;
		error?: string;
		success?: string;
	};

	let { data, form }: PageProps & { form?: RecoverPasswordFormState | null } = $props();

	let submitting = $state<"request" | "reset" | null>(null);
	const emailValue = $derived(form?.values?.email ?? data.email);
	const tokenValue = $derived(form?.values?.token ?? data.token);
	const resetMode = $derived(Boolean(data.email && data.token));

	function submitFor(action: "request" | "reset"): SubmitFunction {
		return () => {
			submitting = action;
			return async ({ update }) => {
				submitting = null;
				await update({ reset: false });
			};
		};
	}
</script>

<div class="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
	<Card.Root class="w-full max-w-sm">
		<Card.Header>
			<Card.Title>{resetMode ? "Reset password" : "Recover password"}</Card.Title>
			<Card.Description>
				{resetMode
					? "Enter a new password for your account."
					: "Enter your email and we'll send a password reset link."}
			</Card.Description>
		</Card.Header>
		<Card.Content>
			{#if resetMode}
				<form method="POST" action="?/reset" use:enhance={submitFor("reset")}>
					<Field.Group>
						{#if form?.error}
							<Field.Error>{form.error}</Field.Error>
						{/if}
						<input type="hidden" name="email" value={emailValue} />
						<input type="hidden" name="token" value={tokenValue} />
						<Field.Field>
							<Field.Label for="recover-password-email">Email</Field.Label>
							<Input id="recover-password-email" type="email" value={emailValue} disabled />
							{#if form?.fieldErrors?.email}
								<Field.Error>{form.fieldErrors.email}</Field.Error>
							{/if}
							{#if form?.fieldErrors?.token}
								<Field.Error>{form.fieldErrors.token}</Field.Error>
							{/if}
						</Field.Field>
						<Field.Field>
							<Field.Label for="recover-password-new">New password</Field.Label>
							<Input id="recover-password-new" name="password" type="password" required />
							<Field.Description>Must be at least 8 characters long.</Field.Description>
							{#if form?.fieldErrors?.password}
								<Field.Error>{form.fieldErrors.password}</Field.Error>
							{/if}
						</Field.Field>
						<Field.Field>
							<Field.Label for="recover-password-confirm">Confirm password</Field.Label>
							<Input
								id="recover-password-confirm"
								name="confirmPassword"
								type="password"
								required
							/>
							{#if form?.fieldErrors?.confirmPassword}
								<Field.Error>{form.fieldErrors.confirmPassword}</Field.Error>
							{/if}
						</Field.Field>
						<Field.Field>
							<Button type="submit" class="w-full" disabled={submitting === "reset"}>
								{submitting === "reset" ? "Resetting..." : "Reset password"}
							</Button>
							<Field.Description class="text-center">
								<a href="/login">Back to login</a>
							</Field.Description>
						</Field.Field>
					</Field.Group>
				</form>
			{:else}
				<form method="POST" action="?/request" use:enhance={submitFor("request")}>
					<Field.Group>
						{#if form?.success}
							<Field.Description>{form.success}</Field.Description>
						{/if}
						{#if form?.error}
							<Field.Error>{form.error}</Field.Error>
						{/if}
						<Field.Field>
							<Field.Label for="recover-password-request-email">Email</Field.Label>
							<Input
								id="recover-password-request-email"
								name="email"
								type="email"
								placeholder="m@example.com"
								value={emailValue}
								required
							/>
							{#if form?.fieldErrors?.email}
								<Field.Error>{form.fieldErrors.email}</Field.Error>
							{/if}
						</Field.Field>
						<Field.Field>
							<Button type="submit" class="w-full" disabled={submitting === "request"}>
								{submitting === "request" ? "Sending..." : "Send reset link"}
							</Button>
							<Field.Description class="text-center">
								Remembered your password? <a href="/login">Sign in</a>
							</Field.Description>
						</Field.Field>
					</Field.Group>
				</form>
			{/if}
		</Card.Content>
	</Card.Root>
</div>
