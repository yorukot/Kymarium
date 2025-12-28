<script lang="ts">
	import { onMount } from 'svelte';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { goto } from '$app/navigation';

	const inBrowser = typeof window !== 'undefined';
	let email = $state('');
	let countdown = $state(3);

	if (inBrowser) {
		const url = new URL(window.location.href);
		email = url.searchParams.get('email') ?? '';
	}

	onMount(() => {
		const interval = setInterval(() => {
			countdown -= 1;
			if (countdown <= 0) {
				clearInterval(interval);
				goto('/auth/login');
			}
		}, 1000);

		return () => clearInterval(interval);
	});
</script>

<div class="flex min-h-screen w-full items-center justify-center px-4">
	<Card.Root class="w-full max-w-lg">
		<Card.Header>
			<Card.Title>Email verified</Card.Title>
			<Card.Description>
				{#if email}
					Your email <span class="font-medium">{email}</span> is verified.
				{:else}
					Your email is verified.
				{/if}
				Please log in to continue.
			</Card.Description>
		</Card.Header>
		<Card.Content class="space-y-4">
			<p class="text-sm text-muted-foreground">
				Redirecting to login in {countdown}s...
			</p>
			<Button class="w-full" href="/auth/login">Go to login</Button>
		</Card.Content>
	</Card.Root>
</div>
