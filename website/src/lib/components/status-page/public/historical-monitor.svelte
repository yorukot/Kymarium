<script lang="ts">
	import { formatSli, statusMeta } from '$lib/styles/status';
	import type { PublicStatusPageMonitor } from '../../../types';
	import HistoricalTimeline from './historical-timeline.svelte';

	let { monitor, days = 90 }: { monitor: PublicStatusPageMonitor; days?: number } = $props();

	function selectSli() {
		if (days <= 30) return formatSli(monitor.uptimeSli30);
		if (days <= 60) return formatSli(monitor.uptimeSli60);
		return formatSli(monitor.uptimeSli90);
	}
</script>

<div class="rounded-md border bg-background/70 p-3">
	<div class="flex flex-wrap items-center justify-between gap-3">
		<div class="flex items-center gap-2">
			<span class="text-lg font-semibold">{monitor.name}</span>
		</div>
		<span class={`text-sm font-medium ${statusMeta[monitor.status ?? 'up'].tone}`}>
			{statusMeta[monitor.status ?? 'up'].label}
		</span>
	</div>
	{#if monitor.type === 'historical_timeline' && monitor.timeline?.length}
		<HistoricalTimeline sli={selectSli()} timeline={monitor.timeline} {days} />
	{/if}
</div>
