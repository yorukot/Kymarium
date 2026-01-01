<script lang="ts">
	import '../layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import AppSidebar from '$lib/components/sidebar/app-sidebar.svelte';
	import * as Breadcrumb from '$lib/components/ui/breadcrumb/index.js';
	import { Separator } from '$lib/components/ui/separator/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import { ModeWatcher } from 'mode-watcher';
	import { page } from '$app/state';
	import type { Page } from '@sveltejs/kit';
	import type { SidebarData } from './+layout';
	import type { Monitor } from '../../lib/types/monitor.js';
	import type { MonitorAnalytics } from '../../lib/types/analytics.js';

	type Crumb = {
		label: string;
		href: string;
	};

	type BreadcrumbPageData = {
		monitor?: Monitor;
		monitors?: Monitor[];
		analytics?: MonitorAnalytics;
	};

	/** @type {import('./$types').PageProps} */
	let { children, data } = $props();

	const crumbs: Crumb[] = $derived(buildBreadcrumbs(page, data));
	const title = $derived(buildTitle(crumbs));

	function buildBreadcrumbs(currentPage: Page, layoutData: SidebarData): Crumb[] {
		const segments = currentPage.url.pathname.split('/').filter(Boolean);
		const pageData = currentPage.data as BreadcrumbPageData | undefined;
		const monitors = pageData?.monitors ?? [];
		const monitor = pageData?.monitor ?? pageData?.analytics?.monitor;

		if (!segments.length) return [];

		const [teamId, ...rest] = segments;
		const teamName = layoutData?.teams?.find((team) => team.id === teamId)?.name ?? 'Team';
		let href = `/${teamId}`;

		const trail: Crumb[] = [
			{
				label: teamName,
				href
			}
		];

		for (let index = 0; index < rest.length; index++) {
			const segment = rest[index];
			href = `${href}/${segment}`;

			let label = formatSegment(segment);

			if (rest[index - 1] === 'monitors') {
				label =
					(monitor && monitor.id === segment && monitor.name) ||
					monitors.find((item) => item.id === segment)?.name ||
					label;
			}

			trail.push({
				label,
				href
			});
		}

		return trail;
	}

	function formatSegment(segment: string): string {
		return segment
			.split('-')
			.filter(Boolean)
			.map((part) => part.charAt(0).toUpperCase() + part.slice(1))
			.join(' ');
	}

	function buildTitle(segments: Crumb[]): string {
		const appName = 'Knocker';
		if (!segments.length) return appName;
		const teamName = segments[0]?.label;
		const last = segments[segments.length - 1]?.label ?? teamName ?? appName;
		if (!teamName || last === teamName) return `${last} - ${appName}`;
		return `${last} - ${teamName} - ${appName}`;
	}
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<title>{title}</title>
</svelte:head>

<ModeWatcher />

<Sidebar.Provider>
	<AppSidebar user={data.user} teams={data.teams} />
	<Sidebar.Inset>
		<header
			class="group-has-data-[collapsible=icon]/sidebar-wrapper:h-12 flex h-16 shrink-0 items-center gap-2 transition-[width,height] ease-linear"
		>
			<div class="flex items-center gap-2 px-4">
				<Sidebar.Trigger class="-ms-1" />
				<Separator orientation="vertical" class="me-2 data-[orientation=vertical]:h-4" />
				{#if crumbs.length}
					<Breadcrumb.Root>
						<Breadcrumb.List>
							{#each crumbs as crumb, index (crumb.href)}
								<Breadcrumb.Item>
									{#if index < crumbs.length - 1}
										<Breadcrumb.Link href={crumb.href}>{crumb.label}</Breadcrumb.Link>
									{:else}
										<Breadcrumb.Page>{crumb.label}</Breadcrumb.Page>
									{/if}
								</Breadcrumb.Item>
								{#if index < crumbs.length - 1}
									<Breadcrumb.Separator class="hidden md:block" />
								{/if}
							{/each}
						</Breadcrumb.List>
					</Breadcrumb.Root>
				{/if}
			</div>
		</header>
		<main class="flex flex-1 flex-col gap-4 p-4 pt-0 w-full">
			<div class="max-w-5xl mx-auto w-full">
				{@render children()}
			</div>
		</main>
	</Sidebar.Inset>
</Sidebar.Provider>
