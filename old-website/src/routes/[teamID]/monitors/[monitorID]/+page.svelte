<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { onMount } from 'svelte';
	import LatencyChart from '$lib/components/monitor/detail/latency-chart.svelte';
	import MonitorStatsCards from '$lib/components/monitor/detail/stats-cards.svelte';
	import UptimeChart from '$lib/components/monitor/detail/uptime-chart.svelte';
	import { Button } from '$lib/components/ui/button';
	import Icon from '@iconify/svelte';

	/** @type {import('./$types').PageProps} */
	let { data } = $props();

	onMount(() => {
		const interval = setInterval(() => {
			invalidateAll();
		}, 10_000);

		return () => clearInterval(interval);
	});
</script>

<div class="flex flex-col gap-10">
	<div class="space-y-4">
		<div class="flex items-center justify-between gap-4">
			<h1 class="text-2xl font-bold">{data.analytics.monitor.name}</h1>
			<Button href={`${data.analytics.monitor.id}/edit`}>
				<Icon icon="lucide:pencil" />
				Go to edit
			</Button>
		</div>
		<!-- There should be a button that can select time for example last 7 days last 24 hours last 14 days last 30 days last 90 days -->
		<MonitorStatsCards analytics={data.analytics} />
	</div>

	<div class="space-y-4">
		<h1 class="text-2xl font-bold">Uptime</h1>
		<UptimeChart analytics={data.analytics} />
	</div>

	<div>
		<LatencyChart analytics={data.analytics} regions={data.regions}/>
	</div>
</div>
