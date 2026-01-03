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
	import { reporter, ValidationMessage } from '@felte/reporter-svelte';
	import { z } from 'zod';
	import Spinner from '../ui/spinner/spinner.svelte';

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

	const { form, isSubmitting, setErrors } = createForm({
		extend: [validator({ schema: signupSchema }), reporter()],
		onSubmit: async (values) => {
			try {
				const response = await registerUser(values.displayName, values.email, values.password);
				if (response?.message?.toLowerCase().includes('verification email sent')) {
					await goto(`/auth/verify/sent?email=${encodeURIComponent(values.email)}`);
					return;
				}
				await goto(redirectTo);
			} catch (error) {
				const message =
					error instanceof Error
						? error.message
						: 'Unable to sign up right now. Please try again.';
				const normalizedMessage = message.toLowerCase();
				if (normalizedMessage.includes('email is already in use')) {
					setErrors({ email: 'This email is already in use.' });
					return;
				}
				setErrors({ FORM_ERROR: message });
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
					<ValidationMessage for="displayName" let:messages>
						{#if messages?.length}
							<Field.Description class="text-destructive">
								{messages[0]}
							</Field.Description>
						{/if}
					</ValidationMessage>
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
					<ValidationMessage for="email" let:messages>
						{#if messages?.length}
							<Field.Description class="text-destructive">
								{messages[0]}
							</Field.Description>
						{/if}
					</ValidationMessage>
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
					<ValidationMessage for="password" let:messages>
						{#if messages?.length}
							<Field.Description class="text-destructive">
								{messages[0]}
							</Field.Description>
						{/if}
					</ValidationMessage>
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
					<ValidationMessage for="confirmPassword" let:messages>
						{#if messages?.length}
							<Field.Description class="text-destructive">
								{messages[0]}
							</Field.Description>
						{/if}
					</ValidationMessage>
				</Field.Field>

				<ValidationMessage for="FORM_ERROR" let:messages>
					{#if messages?.length}
						<Field.Description class="text-destructive text-center">
							{messages[0]}
						</Field.Description>
					{/if}
				</ValidationMessage>

				<div class="flex flex-col gap-3">
					<Button type="submit" disabled={$isSubmitting}>
						{#if $isSubmitting}
							<Spinner /> Creating...
						{:else}
							Create Account
						{/if}
					</Button>

					<Button type="button" variant="outline" onclick={handleGoogle}>
						<Icon icon="ri:google-fill" class="size-5" />
						Sign up with Google
					</Button>

					<Field.Description class="px-6 text-center">
						Already have an account?
						<a href="/auth/login">Sign in</a>
					</Field.Description>
				</div>
			</Field.Group>
		</form>
	</Card.Content>
</Card.Root>
