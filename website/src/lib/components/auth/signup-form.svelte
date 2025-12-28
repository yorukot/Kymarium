<script lang="ts">
	import { goto } from '$app/navigation';
	import { buildOAuthUrl, registerUser } from '$lib/api/auth.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Field from '$lib/components/ui/field/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import Icon from '@iconify/svelte';

	import { createForm } from 'felte';
	import { validator } from '@felte/validator-zod';
	import { z } from 'zod';

	const inBrowser = typeof window !== 'undefined';

	let redirectTo = '/';
	if (inBrowser) {
		const url = new URL(window.location.href);
		redirectTo = url.searchParams.get('next') ?? '/';
	}

	const signupSchema = z
		.object({
			displayName: z.string().min(3).max(255),
			email: z.email().max(255),
			password: z.string().min(8).max(255),
			confirmPassword: z.string().min(8).max(255)
		})
		.refine((data) => data.password === data.confirmPassword, {
			path: ['confirmPassword'],
			message: 'Passwords do not match.'
		});

	const { form, errors, isSubmitting } = createForm({
		extend: validator({ schema: signupSchema }),
		onSubmit: async (values) => {
			try {
				const response = await registerUser(values.displayName, values.email, values.password);
				if (response?.message?.toLowerCase().includes('verification email sent')) {
					await goto(`/auth/verify/sent?email=${encodeURIComponent(values.email)}`);
					return;
				}
				await goto(redirectTo);
			} catch (error) {
				return {
					FORM_ERROR:
						error instanceof Error
							? error.message
							: 'Unable to sign up right now. Please try again.'
				};
			}
		}
	});

	const handleGoogle = () => {
		if (!inBrowser) return;
		window.location.href = buildOAuthUrl('google', redirectTo);
	};
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>Create an account</Card.Title>
		<Card.Description>Enter your information below to create your account</Card.Description>
	</Card.Header>

	<Card.Content>
		<form use:form class="space-y-4">
			<Field.Group>
				<Field.Field>
					<Field.Label for="displayName">Display Name</Field.Label>
					<Input
						id="displayName"
						name="displayName"
						type="text"
						placeholder="John Doe"
						autocomplete="name"
						required
					/>
					{#if $errors.displayName}
						<Field.Description class="text-destructive">
							{$errors.displayName[0]}
						</Field.Description>
					{/if}
				</Field.Field>

				<Field.Field>
					<Field.Label for="email">Email</Field.Label>
					<Input
						id="email"
						name="email"
						type="email"
						placeholder="m@example.com"
						autocomplete="email"
						required
					/>
					<Field.Description>
						We'll use this to contact you. We will not share your email.
					</Field.Description>
					{#if $errors.email}
						<Field.Description class="text-destructive">
							{$errors.email[0]}
						</Field.Description>
					{/if}
				</Field.Field>

				<Field.Field>
					<Field.Label for="password">Password</Field.Label>
					<Input
						id="password"
						name="password"
						type="password"
						autocomplete="new-password"
						required
					/>
					<Field.Description>Must be at least 8 characters long.</Field.Description>
					{#if $errors.password}
						<Field.Description class="text-destructive">
							{$errors.password[0]}
						</Field.Description>
					{/if}
				</Field.Field>

				<Field.Field>
					<Field.Label for="confirmPassword">Confirm Password</Field.Label>
					<Input
						id="confirmPassword"
						name="confirmPassword"
						type="password"
						autocomplete="new-password"
						required
					/>
					{#if $errors.confirmPassword}
						<Field.Description class="text-destructive">
							{$errors.confirmPassword[0]}
						</Field.Description>
					{/if}
				</Field.Field>

				{#if $errors.FORM_ERROR}
					<Field.Description class="text-destructive text-center">
						{$errors.FORM_ERROR}
					</Field.Description>
				{/if}

				<Field.Field>
					<Button type="submit" disabled={$isSubmitting}>
						{$isSubmitting ? 'Creating account...' : 'Create Account'}
					</Button>

					<Button type="button" variant="outline" onclick={handleGoogle}>
						<Icon icon="ri:google-fill" class="size-5" />
						Sign up with Google
					</Button>

					<Field.Description class="px-6 text-center">
						Already have an account?
						<a href="/auth/login">Sign in</a>
					</Field.Description>
				</Field.Field>
			</Field.Group>
		</form>
	</Card.Content>
</Card.Root>
