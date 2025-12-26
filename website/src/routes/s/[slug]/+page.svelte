<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import type {
		PublicStatusPageData,
		PublicStatusPageElement,
		PublicStatusPageMonitor
	} from '$lib/types';
	import HistoricalMonitor from '$lib/components/status-page/public/historical-monitor.svelte';
	import Group from '$lib/components/status-page/public/group.svelte';
	import { statusMeta } from '$lib/styles/status';

	type Props = {
		data: {
			statusPage: PublicStatusPageData;
		};
	};

	let { data }: Props = $props();
	const statusPage = $derived(data.statusPage);

	const sortedElements = $derived.by(() =>
		[...(statusPage.elements ?? [])].sort((a, b) => a.sortOrder - b.sortOrder)
	);
	const groupedMonitors = $derived.by(() =>
		sortedElements.filter((element) => !element.monitor).flatMap((group) => group.monitors ?? [])
	);

	const topLevelMonitors = $derived.by(() =>
		sortedElements.filter((element) => element.monitor).map((element) => toPublicMonitor(element))
	);

	const allMonitors = $derived.by(() => [...topLevelMonitors, ...groupedMonitors]);

	const openIncidents = $derived.by(() =>
		statusPage.incidents.filter((incident) => incident.status !== 'resolved')
	);

	const overallStatus = $derived.by(() => {
		if (openIncidents.length) return 'down';
		return allMonitors.some((monitor) => monitor.status === 'down') ? 'down' : 'up';
	});

	let days = $state(60);

	const MOBILE = '(max-width: 639px)'; // < 640
	const TABLET = '(min-width: 640px) and (max-width: 1023px)'; // 640~1023
	const DESKTOP = '(min-width: 1024px)'; // >= 1024

	function computeDays() {
		const mobile = window.matchMedia(MOBILE).matches;
		const tablet = window.matchMedia(TABLET).matches;

		if (mobile) return 30;
		if (tablet) return 60;
		return 90;
	}

	$effect(() => {
		if (typeof window === 'undefined') return;

		const mqs = [window.matchMedia(MOBILE), window.matchMedia(TABLET), window.matchMedia(DESKTOP)];

		const update = () => {
			days = computeDays();
		};

		update();

		mqs.forEach((mq) => mq.addEventListener('change', update));

		return () => {
			mqs.forEach((mq) => mq.removeEventListener('change', update));
		};
	});

	function toPublicMonitor(element: PublicStatusPageElement): PublicStatusPageMonitor {
		return {
			id: element.id,
			monitorId: element.monitorId ?? element.id,
			groupId: null,
			name: element.name,
			type: element.type,
			sortOrder: element.sortOrder,
			status: element.status,
			uptimeSli30: element.uptimeSli30,
			uptimeSli60: element.uptimeSli60,
			uptimeSli90: element.uptimeSli90,
			timeline: element.timeline
		};
	}
</script>
<div class="min-h-screen bg-muted/30">
	<div class="mx-auto flex w-full max-w-5xl flex-col gap-6 px-4 py-10 md:px-8">
		<header class="flex flex-col gap-4">
			<div class="flex flex-col gap-2">
				<h1 class="text-3xl font-bold md:text-4xl">{statusPage.statusPage.title}</h1>
				<p class="text-sm text-muted-foreground">
					Status updates for {statusPage.statusPage.slug}
				</p>
			</div>
			<div class="flex items-center gap-3 w-full p-4 {statusMeta[overallStatus].dot} rounded">
				<p class="text-sm font-semibold uppercase">
					{#if overallStatus === 'up'}
						All systems operational
					{:else if overallStatus === 'down'}
						Some systems are experiencing issues
					{/if}
				</p>
			</div>
		</header>
		
		<section class="flex flex-col gap-4">
			{#if sortedElements.length === 0}
				<Card.Root class="p-6 text-sm text-muted-foreground">
					No monitors are configured for this status page yet.
				</Card.Root>
			{/if}

			{#each sortedElements as element (element.id)}
				{#if element.monitor}
					<HistoricalMonitor monitor={toPublicMonitor(element)} {days} />
				{:else}
					<Group group={element} {days} />
				{/if}
			{/each}
		</section>
	</div>
</div>
