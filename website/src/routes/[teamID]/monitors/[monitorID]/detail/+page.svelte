<script lang="ts">
	import MonitorStatsCards from '$lib/components/monitor/detail/stats-cards.svelte';
	import UptimeChart from '$lib/components/monitor/detail/uptime-chart.svelte';

	/** @type {import('./$types').PageProps} */
	let { data } = $props();
</script>

<div class="flex flex-col gap-4">
	<h1 class="text-2xl font-bold">Monitor Analytics</h1>
	<!-- There should be a button that can select time for example last 7 days last 24 hours last 14 days last 30 days last 90 days -->
	<MonitorStatsCards analytics={data.analytics} />

	<div>
		<h1 class="text-2xl font-bold">Uptime</h1>
		<UptimeChart analytics={data.analytics} />
	</div>
</div>
