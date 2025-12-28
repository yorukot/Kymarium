<script lang="ts">
	import { onMount } from 'svelte';
	import { PUBLIC_API_BASE } from '$env/static/public';
	import * as Card from '$lib/components/ui/card/index.js';

	const inBrowser = typeof window !== 'undefined';
	let status: 'loading' | 'missing' = $state('loading');

	onMount(async () => {
		if (!inBrowser) return;
		const url = new URL(window.location.href);
		const token = url.searchParams.get('token');
		if (!token) {
			status = 'missing';
			return;
		}

		window.location.href = `${PUBLIC_API_BASE}/auth/verify?token=${encodeURIComponent(token)}`;
	});
</script>

<div class="flex min-h-screen w-full items-center justify-center px-4">
	<Card.Root class="w-full max-w-lg">
		<Card.Header>
			<Card.Title>
				{status === 'missing' ? 'Missing verification link' : 'Verifying your email...'}
			</Card.Title>
			<Card.Description>
				{#if status === 'missing'}
					The verification link is missing. Please log in and request a new one.
				{:else}
					Please wait while we redirect you to complete verification.
				{/if}
			</Card.Description>
		</Card.Header>
	</Card.Root>
</div>
