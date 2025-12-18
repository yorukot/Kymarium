<script lang="ts">
	import * as Chart from '$lib/components/ui/chart/index.js';
	import { scaleBand } from 'd3-scale';
	import { BarChart, type ChartContextValue } from 'layerchart';
	import { cubicInOut } from 'svelte/easing';
	import MonitorStatsCards from '$lib/components/monitor/detail/stats-cards.svelte';

	/** @type {import('./$types').PageProps} */
	let { data } = $props();

	const chartData = [
		{ month: 'January', desktop: 186, mobile: 80 },
		{ month: 'February', desktop: 305, mobile: 200 },
		{ month: 'March', desktop: 237, mobile: 120 },
		{ month: 'April', desktop: 73, mobile: 190 },
		{ month: 'May', desktop: 209, mobile: 130 },
		{ month: 'Test', desktop: 214, mobile: 140 },
		{ month: 'A', desktop: 214, mobile: 140 },
		{ month: 'B', desktop: 214, mobile: 140 },
		{ month: 'C', desktop: 214, mobile: 140 },
		{ month: 'D', desktop: 214, mobile: 140 },
		{ month: 'F', desktop: 214, mobile: 140 },
		{ month: 'G', desktop: 214, mobile: 140 }
	];

	const chartConfig = {
		desktop: { label: 'Desktop', color: 'var(--success)' },
		mobile: { label: 'Mobile', color: 'var(--destructive)' }
	} satisfies Chart.ChartConfig;

	let context = $state<ChartContextValue>();
</script>

<div class="flex flex-col gap-4">
	<h1 class="text-2xl font-bold">Monitor Analytics</h1>
	<!-- There should be a button that can select time for example last 7 days last 24 hours last 14 days last 30 days last 90 days -->
	<MonitorStatsCards analytics={data.analytics} />

	<div>
		<h1 class="text-2xl font-bold">Uptime</h1>
		<Chart.Container config={chartConfig} class="h-36 w-full">
			<BarChart
				bind:context
				data={chartData}
				xScale={scaleBand().padding(0.2)}
				x="month"
				axis="x"
				rule={true}
				grid={false}
				series={[
					{
						key: 'desktop',
						label: 'Desktop',
						color: chartConfig.desktop.color,
						props: { rounded: 'bottom' }
					},
					{
						key: 'mobile',
						label: 'Mobile',
						color: chartConfig.mobile.color
					}
				]}
				seriesLayout="stack"
				props={{
					bars: {
						stroke: 'none',
						initialY: context?.height,
						initialHeight: 0,
						motion: {
							y: { type: 'tween', duration: 500, easing: cubicInOut },
							height: { type: 'tween', duration: 500, easing: cubicInOut }
						}
					}
				}}
			>
				{#snippet tooltip()}
					<Chart.Tooltip />
				{/snippet}
			</BarChart>
		</Chart.Container>
	</div>
	<pre>{JSON.stringify(data.analytics, null, 2)}</pre>
</div>
