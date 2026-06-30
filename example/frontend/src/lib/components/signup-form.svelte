<script lang="ts">
	import { enhance } from "$app/forms";
	import { Button } from "$lib/components/ui/button/index.js";
	import * as Card from "$lib/components/ui/card/index.js";
	import * as Field from "$lib/components/ui/field/index.js";
	import { Input } from "$lib/components/ui/input/index.js";
	import type { ComponentProps } from "svelte";
	import type { SubmitFunction } from "@sveltejs/kit";

	type SignupFormState = {
		values?: { name?: string; email?: string };
		fieldErrors?: Record<string, string>;
		error?: string;
	};

	let {
		form,
		...restProps
	}: ComponentProps<typeof Card.Root> & { form?: SignupFormState | null } = $props();

	let submitting = $state(false);

	const submit: SubmitFunction = () => {
		submitting = true;
		return async ({ update }) => {
			submitting = false;
			await update({ reset: false });
		};
	};
</script>

<Card.Root {...restProps}>
	<Card.Header>
		<Card.Title>Create an account</Card.Title>
		<Card.Description>Enter your information below to create your account</Card.Description>
	</Card.Header>
	<Card.Content>
		<form method="POST" use:enhance={submit}>
			<Field.Group>
				{#if form?.error}
					<Field.Error>{form.error}</Field.Error>
				{/if}
				<Field.Field>
					<Field.Label for="name">Full Name</Field.Label>
					<Input
						id="name"
						name="name"
						type="text"
						placeholder="John Doe"
						value={form?.values?.name ?? ""}
						required
					/>
					{#if form?.fieldErrors?.name}
						<Field.Error>{form.fieldErrors.name}</Field.Error>
					{/if}
				</Field.Field>
				<Field.Field>
					<Field.Label for="email">Email</Field.Label>
					<Input
						id="email"
						name="email"
						type="email"
						placeholder="m@example.com"
						value={form?.values?.email ?? ""}
						required
					/>
					<Field.Description>
						We'll use this to contact you. We will not share your email with anyone
						else.
					</Field.Description>
					{#if form?.fieldErrors?.email}
						<Field.Error>{form.fieldErrors.email}</Field.Error>
					{/if}
				</Field.Field>
				<Field.Field>
					<Field.Label for="password">Password</Field.Label>
					<Input id="password" name="password" type="password" required />
					<Field.Description>Must be at least 8 characters long.</Field.Description>
					{#if form?.fieldErrors?.password}
						<Field.Error>{form.fieldErrors.password}</Field.Error>
					{/if}
				</Field.Field>
				<Field.Field>
					<Field.Label for="confirm-password">Confirm Password</Field.Label>
					<Input id="confirm-password" name="confirmPassword" type="password" required />
					{#if form?.fieldErrors?.confirmPassword}
						<Field.Error>{form.fieldErrors.confirmPassword}</Field.Error>
					{/if}
				</Field.Field>
				<Field.Group>
					<Field.Field>
						<Button type="submit" disabled={submitting}>
							{submitting ? "Creating account..." : "Create Account"}
						</Button>
						<Field.Description class="px-6 text-center">
							Already have an account? <a href="/login">Sign in</a>
						</Field.Description>
					</Field.Field>
				</Field.Group>
			</Field.Group>
		</form>
	</Card.Content>
</Card.Root>
