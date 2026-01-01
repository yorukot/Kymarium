<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import { ModeWatcher } from 'mode-watcher';
	import { Toaster } from '$lib/components/ui/sonner/index.js';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { page } from '$app/state';

	let { children } = $props();

	const appName = 'Knocker';

	const title = $derived(buildTitle(page));

	function buildTitle(currentPage: typeof page): string {
		const path = currentPage.url.pathname;
		const authTitles: Record<string, string> = {
			'/auth/login': 'Log in',
			'/auth/register': 'Create account',
			'/auth/verify': 'Verify email',
			'/auth/verify/sent': 'Verification email sent',
			'/auth/verify/success': 'Email verified',
			'/auth/verify/expired': 'Verification expired'
		};

		if (path.startsWith('/s/')) {
			const statusTitle = (currentPage.data as { statusPage?: { statusPage?: { title?: string } } })
				?.statusPage?.statusPage?.title;
			if (statusTitle) return `${statusTitle} - Status - ${appName}`;
			return `Status - ${appName}`;
		}

		if (authTitles[path]) return `${authTitles[path]} - ${appName}`;
		if (path.startsWith('/auth')) return `Authentication - ${appName}`;
		if (path === '/' || path === '') return appName;

		const lastSegment = path.split('/').filter(Boolean).pop();
		if (!lastSegment) return appName;
		return `${formatSegment(lastSegment)} - ${appName}`;
	}

	function formatSegment(segment: string): string {
		return segment
			.split('-')
			.filter(Boolean)
			.map((part) => part.charAt(0).toUpperCase() + part.slice(1))
			.join(' ');
	}
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<title>{title}</title>
</svelte:head>

<ModeWatcher />
<Toaster />
<Tooltip.Provider>
	{@render children()}
</Tooltip.Provider>
