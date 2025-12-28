<script lang="ts">
	import { onMount } from 'svelte';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { resendVerification } from '$lib/api/auth.js';

	const inBrowser = typeof window !== 'undefined';
	let email = $state('');
	if (inBrowser) {
		const url = new URL(window.location.href);
		email = url.searchParams.get('email') ?? '';
	}

	let cooldown = $state(0);
	let isSending = $state(false);
	let feedback = $state('');

	const formatCooldown = (seconds: number) => {
		const padded = String(seconds).padStart(2, '0');
		return `0:${padded}`;
	};

	const tickCooldown = () => {
		if (cooldown <= 0) return;
		cooldown -= 1;
		if (cooldown === 0) {
			feedback = '';
		}
	};

	onMount(() => {
		const interval = setInterval(tickCooldown, 1000);
		return () => clearInterval(interval);
	});

	const handleResend = async () => {
		if (!email || cooldown > 0 || isSending) return;
		isSending = true;
		feedback = '';
		try {
			await resendVerification(email);
			cooldown = 60;
			feedback = 'Verification email sent again.';
		} catch (error) {
			feedback =
				error instanceof Error
					? error.message
					: 'Unable to resend verification email. Please try again.';
		} finally {
			isSending = false;
		}
	};
</script>

<div class="flex min-h-screen w-full items-center justify-center px-4">
	<Card.Root class="w-full max-w-lg">
		<Card.Header>
			<Card.Title>Check your inbox</Card.Title>
			<Card.Description>
				{#if email}
					We sent a verification link to <span class="font-medium">{email}</span>.
				{:else}
					We sent a verification link to your email address.
				{/if}
				Please verify to continue.
			</Card.Description>
		</Card.Header>
		<Card.Content class="space-y-4">
			<p class="text-sm text-muted-foreground">
				Didn’t receive the email? Check your spam folder or try again in a few minutes.
			</p>
			<Button
				type="button"
				variant="outline"
				class="w-full"
				disabled={!email || isSending || cooldown > 0}
				onclick={handleResend}
			>
				{#if cooldown > 0}
					Resend available in {formatCooldown(cooldown)}
				{:else}
					{isSending ? 'Sending...' : 'Resend verification email'}
				{/if}
			</Button>
			{#if feedback}
				<p class="text-sm text-muted-foreground">{feedback}</p>
			{/if}
			<Button variant="outline" class="w-full" href="/auth/login">Back to login</Button>
		</Card.Content>
	</Card.Root>
</div>
