<script lang="ts">
	import * as Chart from '$lib/components/ui/chart/index.js';
	import { scaleBand } from 'd3-scale';
	import { BarChart } from 'layerchart';
	import { cubicInOut } from 'svelte/easing';
	import { SvelteMap } from 'svelte/reactivity';
	import type { MonitorAnalytics } from '../$lib/types';

	type Props = {
		analytics: MonitorAnalytics;
	};

	type AggregatedBucket = {
		timestamp: Date;
		success: number;
		failed: number;
		total: number;
	};

	let { analytics }: Props = $props();

	// Aggregate across regions per time bucket, then stack success vs failed.
	const aggregated = $derived.by(() => {
		const bucketMs = parseBucketDuration(analytics.window?.bucket);
		const buckets = new SvelteMap<number, { total: number; good: number }>();

		for (const point of analytics.series ?? []) {
			const bucketTime = parseTimestamp(point.timestamp);
			if (bucketTime === undefined) continue;

			const bucket = buckets.get(bucketTime) ?? { total: 0, good: 0 };
			bucket.total += point.totalCount;
			bucket.good += point.goodCount;
			buckets.set(bucketTime, bucket);
		}

		const bucketTimes = Array.from(buckets.keys()).sort((a, b) => a - b);

		if (bucketTimes.length === 0) return [];

		const startTimeRaw = parseTimestamp(analytics.window?.start);
		const endTimeRaw = parseTimestamp(analytics.window?.end);

		const startTime =
			startTimeRaw !== undefined ? Math.floor(startTimeRaw / bucketMs) * bucketMs : bucketTimes[0];

		const endTime =
			endTimeRaw !== undefined
				? Math.ceil(endTimeRaw / bucketMs) * bucketMs
				: bucketTimes[bucketTimes.length - 1] + bucketMs;

		if (bucketMs <= 0 || Number.isNaN(startTime) || Number.isNaN(endTime) || startTime > endTime) {
			return bucketTimes.map((ts) =>
				buildBucketEntry(ts, buckets.get(ts) ?? { total: 0, good: 0 })
			);
		}

		const filled: AggregatedBucket[] = [];

		for (let ts = startTime; ts < endTime; ts += bucketMs) {
			const bucket = buckets.get(ts) ?? { total: 0, good: 0 };
			filled.push(buildBucketEntry(ts, bucket));
		}

		return filled;
	});

	const xAxisTicks = $derived.by(() => {
		const dates = aggregated.map((b) => b.timestamp);
		if (dates.length <= 1) return dates;
		return [dates[0], dates[dates.length - 1]];
	});

	function buildBucketEntry(
		timestamp: number,
		bucket: { total: number; good: number }
	): AggregatedBucket {
		return {
			timestamp: new Date(timestamp),
			success: bucket.good,
			failed: Math.max(bucket.total - bucket.good, 0),
			total: bucket.total
		};
	}

	function parseBucketDuration(bucket?: string): number {
		const match = /^(\d+)([smhd])$/.exec(bucket ?? '');
		if (!match) {
			return 30 * 60 * 1000;
		}

		const value = Number(match[1]);
		const unit = match[2] as 's' | 'm' | 'h' | 'd';
		const unitMsMap: Record<'s' | 'm' | 'h' | 'd', number> = {
			s: 1_000,
			m: 60_000,
			h: 3_600_000,
			d: 86_400_000
		};

		if (!Number.isFinite(value)) {
			return 30 * 60 * 1000;
		}

		return (unitMsMap[unit] ?? 30 * 60 * 1000) * value;
	}

	function parseTimestamp(value?: string): number | undefined {
		if (!value) {
			return undefined;
		}

		const time = new Date(value).getTime();
		return Number.isNaN(time) ? undefined : time;
	}

	const chartConfig = {
		success: { label: 'Success', color: 'var(--success)' },
		failed: { label: 'Failed', color: 'var(--destructive)' }
	} satisfies Chart.ChartConfig;
</script>

<Chart.Container config={chartConfig} class="h-20 w-full">
	<BarChart
		data={aggregated}
		xScale={scaleBand().padding(0.2)}
		x="timestamp"
		axis="x"
		rule={true}
		grid={false}
		series={[
			{
				key: 'success',
				label: 'Success',
				color: chartConfig.success.color,
				props: { rounded: 'none' }
			},
			{
				key: 'failed',
				label: 'Failed',
				color: chartConfig.failed.color,
				props: { rounded: 'none' }
			}
		]}
		seriesLayout="stack"
		props={{
			bars: {
				stroke: 'none',
				motion: {
					y: { type: 'tween', duration: 400, easing: cubicInOut },
					height: { type: 'tween', duration: 400, easing: cubicInOut }
				}
			},
			xAxis: {
				ticks: xAxisTicks.length ? xAxisTicks : undefined,
				tickSpacing: 64,
				format: (value) =>
					value.toLocaleTimeString('en', {
						hour: '2-digit',
						minute: '2-digit'
					})
			}
		}}
	>
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
	</BarChart>
</Chart.Container>
