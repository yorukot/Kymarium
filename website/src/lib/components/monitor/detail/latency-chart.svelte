<script lang="ts">
	import * as Chart from '$lib/components/ui/chart/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import ChartContainer from '$lib/components/ui/chart/chart-container.svelte';
	import { scaleUtc } from 'd3-scale';
	import { Area, AreaChart, ChartClipPath } from 'layerchart';
	import { curveLinear } from "d3-shape";
	import { cubicInOut } from 'svelte/easing';
	import { SvelteMap } from 'svelte/reactivity';
	import type { MonitorAnalytics, Region } from '../$lib/types';
	import { regionFlagEmoji } from '$lib/utils/region';

	type PercentileKey = 'p50Ms' | 'p75Ms' | 'p90Ms' | 'p95Ms' | 'p99Ms';

	type ChartPoint = {
		date: Date;
		latency: number;
	} & Record<string, number | Date>;

	type Props = {
		analytics: MonitorAnalytics;
		regions: Region[];
	};

	const chartColors = [
		'var(--chart-1)',
		'var(--chart-2)',
		'var(--chart-3)',
		'var(--chart-4)',
		'var(--chart-5)'
	];

	const percentileOptions: { value: PercentileKey; label: string }[] = [
		{ value: 'p50Ms', label: 'P50' },
		{ value: 'p75Ms', label: 'P75' },
		{ value: 'p90Ms', label: 'P90' },
		{ value: 'p95Ms', label: 'P95' },
		{ value: 'p99Ms', label: 'P99' }
	];

	let { analytics, regions }: Props = $props();

	let percentile: PercentileKey = $state('p95Ms');

	const percentileLabel = $derived(
		percentileOptions.find((option) => option.value === percentile)?.label ?? 'P95'
	);

	const regionIds = $derived.by(() => (analytics?.regions ?? []).map((region) => region.regionId));

	const regionNameMap = $derived.by(() => {
		const map = new SvelteMap<string, string>();

		for (const region of regions) {
			map.set(region.id, region.displayName);
		}

		return map;
	});

	const regionById = $derived.by(() => {
		const map = new SvelteMap<string, Region>();
		for (const region of regions) {
			map.set(region.id, region);
		}
		return map;
	});

	const chartConfig = $derived.by(() =>
		regionIds.reduce((config, regionId, index) => {
			const color = chartColors[index % chartColors.length];
			const region = regionById.get(regionId);
			const flag = region ? regionFlagEmoji(region) : null;

			config[regionId] = {
				label: `${flag ?? ''}${flag ? ' ' : ''}${regionNameMap.get(regionId) ?? regionId}`,
				color
			};
			return config;
		}, {} as Chart.ChartConfig)
	);

	const chartSeries = $derived.by(() =>
		regionIds.map((regionId, index) => {
			const region = regionById.get(regionId);
			const flag = region ? regionFlagEmoji(region) : null;

			return {
				key: regionId,
				label: `${flag ?? ''}${flag ? ' ' : ''}${regionNameMap.get(regionId) ?? regionId}`,
				color: chartColors[index % chartColors.length]
			};
		})
	);

	const latencySeries = $derived.by(() => buildLatencySeries(analytics, percentile));

	const xAxisTicks = $derived.by(() => {
		const dates = latencySeries.map((p) => p.date);
		if (dates.length === 0) return [];
		if (dates.length === 1) return [dates[0]];
		return [dates[0], dates[dates.length - 1]];
	});

	function buildLatencySeries(data: MonitorAnalytics, pKey: PercentileKey): ChartPoint[] {
		if (!data?.series?.length) return [];

		const buckets = new SvelteMap<
			number,
			SvelteMap<
				string,
				{
					sum: number;
					weight: number;
				}
			>
		>();

		for (const point of data.series) {
			const timestamp = parseTimestamp(point.timestamp);
			const latency = point[pKey];
			if (timestamp === undefined || !Number.isFinite(latency)) continue;

			const weight = point.totalCount > 0 ? point.totalCount : 1;
			const regions = buckets.get(timestamp) ?? new SvelteMap();
			const bucket = regions.get(point.regionId) ?? { sum: 0, weight: 0 };

			bucket.sum += latency * weight;
			bucket.weight += weight;
			regions.set(point.regionId, bucket);
			buckets.set(timestamp, regions);
		}

		return Array.from(buckets.entries())
			.sort((a, b) => a[0] - b[0])
			.map(([timestamp, regionMap]) => {
				const entry: ChartPoint = {
					date: new Date(timestamp),
					latency: 0
				};

				for (const [regionId, bucket] of regionMap.entries()) {
					entry[regionId] = bucket.weight > 0 ? bucket.sum / bucket.weight : 0;
				}

				return entry;
			});
	}

	function parseTimestamp(value?: string): number | undefined {
		if (!value) return undefined;
		const time = new Date(value).getTime();
		return Number.isNaN(time) ? undefined : time;
	}

	function formatWindowRange(window?: { start?: string; end?: string }): string {
		if (!window?.start || !window?.end) return 'Latency over time';
		const start = new Date(window.start);
		const end = new Date(window.end);
		if (Number.isNaN(start.getTime()) || Number.isNaN(end.getTime())) return 'Latency over time';

		const formatter = new Intl.DateTimeFormat('en', {
			month: 'short',
			day: 'numeric'
		});

		return `${formatter.format(start)} - ${formatter.format(end)}`;
	}
</script>

<div class="space-y-4">
	<div class="flex flex-wrap items-center justify-between gap-4 mb-4">
		<div class="flex gap-2 items-baseline">
			<h1 class="text-2xl font-bold">Latency</h1>

			<span class="text-sm text-foreground/50">{formatWindowRange(analytics.window)}</span>
		</div>
		<Select.Root type="single" bind:value={percentile}>
			<Select.Trigger class="w-28 rounded-lg sm:ms-auto" aria-label="Select percentile">
				{percentileLabel}
			</Select.Trigger>
			<Select.Content class="rounded-xl">
				{#each percentileOptions as option (option.value)}
					<Select.Item value={option.value} class="rounded-lg">
						{option.label}
					</Select.Item>
				{/each}
			</Select.Content>
		</Select.Root>
	</div>
	<div class="m-5">
		<ChartContainer config={chartConfig} class="aspect-auto h-[250px] w-full">
			<AreaChart
				legend
				data={latencySeries}
				x="date"
				xScale={scaleUtc()}
				series={chartSeries}
				seriesLayout="overlap"
				props={{
					area: {
						curve: curveLinear,
						'fill-opacity': 0.4,
						line: { class: 'stroke-1' },
						motion: 'tween'
					},
					xAxis: {
						ticks: xAxisTicks,
						format: (value) =>
							new Intl.DateTimeFormat('en', {
								hour: '2-digit',
								minute: '2-digit'
							}).format(value as Date)
					},
					yAxis: { format: () => '' }
				}}
			>
				{#snippet marks({ series, getAreaProps })}
					<ChartClipPath
						initialWidth={0}
						motion={{
							width: { type: 'tween', duration: 1000, easing: cubicInOut }
						}}
					>
						{#each series as s, i (s.key)}
							<Area {...getAreaProps(s, i)} fill={s.color} />
						{/each}
					</ChartClipPath>
				{/snippet}
				{#snippet tooltip()}
					<Chart.Tooltip
						labelFormatter={(v: Date) => {
							return v.toLocaleDateString('en-US', {
								month: 'short',
								day: 'numeric',
								hour: '2-digit',
								minute: '2-digit'
							});
						}}
						indicator="line"
					/>
				{/snippet}
			</AreaChart>
		</ChartContainer>
	</div>
</div>
