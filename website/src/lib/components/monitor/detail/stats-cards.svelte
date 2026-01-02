<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import type { MonitorAnalytics } from '../$lib/types';

	type Props = {
		analytics: MonitorAnalytics;
	};

	let { analytics }: Props = $props();

	const summary = $derived(analytics.summary);
	const incidentsCount = $derived(analytics.incidents?.length ?? 0);
	const totalRequests = $derived(summary?.totalCount ?? 0);
	const failedRequests = $derived(Math.max(totalRequests - (summary?.goodCount ?? 0), 0));

	const latestCheck = $derived(analytics.monitor.lastChecked ?? analytics.window.end);

	const formatPercent = (value: number) =>
		Number.isFinite(value) ? `${value.toFixed(2)}%` : '–';

	const formatNumber = (value: number) =>
		Number.isFinite(value) ? value.toLocaleString() : '–';

	const formatLatency = (value: number) =>
		Number.isFinite(value) && value > 0 ? `${Math.round(value)}ms` : '–';

	const formatRelativeTime = (isoDate: string | undefined) => {
		if (!isoDate) return '–';
		const parsed = new Date(isoDate);
		if (Number.isNaN(parsed.getTime())) return '–';

		const diffSeconds = Math.floor((Date.now() - parsed.getTime()) / 1000);
		if (diffSeconds < 0) return 'just now';
		if (diffSeconds < 60) return `${diffSeconds}s ago`;
		const diffMinutes = Math.floor(diffSeconds / 60);
		if (diffMinutes < 60) return `${diffMinutes}m ago`;
		const diffHours = Math.floor(diffMinutes / 60);
		if (diffHours < 24) return `${diffHours}h ago`;
		const diffDays = Math.floor(diffHours / 24);
		return `${diffDays}d ago`;
	};
</script>

<div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4">
	<Card.Root class="py-2 px-4 gap-0 border-success bg-success/20">
		<div class="text-md text-success">UPTIME</div>
		<div class="text-lg">{formatPercent(summary?.uptimePct ?? 0)}</div>
	</Card.Root>

	<Card.Root class="py-2 px-4 gap-0 border-destructive bg-destructive/20">
		<div class="text-md text-destructive">Failed request</div>
		<div class="text-lg">{formatNumber(failedRequests)}</div>
	</Card.Root>
	<Card.Root class="py-2 px-4 gap-0">
		<div class="text-md text-foreground/50">Total request</div>
		<div class="text-lg">{formatNumber(totalRequests)}</div>
	</Card.Root>

	<Card.Root class="py-2 px-4 gap-0">
		<div class="text-md text-foreground/50">Total incident</div>
		<div class="text-lg">{formatNumber(incidentsCount)}</div>
	</Card.Root>

	<Card.Root class="py-2 px-4 gap-0 border-success bg-transparent border-none">
		<div class="text-md text-card-foreground/50">Latest check</div>
		<div class="text-lg">{formatRelativeTime(latestCheck)}</div>
	</Card.Root>
	<Card.Root class="py-2 px-4 gap-0">
		<div class="text-md text-card-foreground/50">P50</div>
		<div class="text-lg">{formatLatency(summary?.p50Ms ?? 0)}</div>
	</Card.Root>

	<Card.Root class="py-2 px-4 gap-0">
		<div class="text-md text-card-foreground/50">P75</div>
		<div class="text-lg">{formatLatency(summary?.p75Ms ?? 0)}</div>
	</Card.Root>

	<Card.Root class="py-2 px-4 gap-0">
		<div class="text-md text-card-foreground/50">P90</div>
		<div class="text-lg">{formatLatency(summary?.p90Ms ?? 0)}</div>
	</Card.Root>

	<Card.Root class="py-2 px-4 gap-0">
		<div class="text-md text-card-foreground/50">P95</div>
		<div class="text-lg">{formatLatency(summary?.p95Ms ?? 0)}</div>
	</Card.Root>

	<Card.Root class="py-2 px-4 gap-0">
		<div class="text-md text-card-foreground/50">P99</div>
		<div class="text-lg">{formatLatency(summary?.p99Ms ?? 0)}</div>
	</Card.Root>
</div>
