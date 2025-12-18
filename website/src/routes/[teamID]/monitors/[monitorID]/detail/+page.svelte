<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Chart from '$lib/components/ui/chart/index.js';
	import { scaleBand } from 'd3-scale';
	import { BarChart, Highlight, type ChartContextValue } from 'layerchart';
	import TrendingUpIcon from '@lucide/svelte/icons/trending-up';
	import { cubicInOut } from 'svelte/easing';

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
	<div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4">
		<Card.Root class="py-2 px-4 gap-0 border-success bg-success/20">
			<div class="text-md text-success">UPTIME</div>
			<div class="text-lg">99.9%</div>
		</Card.Root>

		<Card.Root class="py-2 px-4 gap-0 border-destructive bg-destructive/20">
			<div class="text-md text-destructive">Failed request</div>
			<div class="text-lg">1k</div>
		</Card.Root>
		<Card.Root class="py-2 px-4 gap-0">
			<div class="text-md text-foreground/50">Total request</div>
			<div class="text-lg">100k</div>
		</Card.Root>

		<Card.Root class="py-2 px-4 gap-0">
			<div class="text-md text-foreground/50">Total incident</div>
			<div class="text-lg">0</div>
		</Card.Root>

		<Card.Root class="py-2 px-4 gap-0 border-success bg-transparent border-none">
			<div class="text-md text-card-foreground/50">Latest check</div>
			<div class="text-lg">10s ago</div>
		</Card.Root>
		<Card.Root class="py-2 px-4 gap-0">
			<div class="text-md text-card-foreground/50">P50</div>
			<div class="text-lg">50ms</div>
		</Card.Root>

		<Card.Root class="py-2 px-4 gap-0">
			<div class="text-md text-card-foreground/50">P75</div>
			<div class="text-lg">50ms</div>
		</Card.Root>

		<Card.Root class="py-2 px-4 gap-0">
			<div class="text-md text-card-foreground/50">P90</div>
			<div class="text-lg">50ms</div>
		</Card.Root>

		<Card.Root class="py-2 px-4 gap-0">
			<div class="text-md text-card-foreground/50">P95</div>
			<div class="text-lg">50ms</div>
		</Card.Root>

		<Card.Root class="py-2 px-4 gap-0">
			<div class="text-md text-card-foreground/50">P99</div>
			<div class="text-lg">50ms</div>
		</Card.Root>
	</div>

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
