<script lang="ts">
	import type { PublicTimelinePoint } from '../../../types';
	import HistoricalTimelinePoint from './historical-timeline-point.svelte';

	let {
		sli,
		timeline,
		days = 90
	}: {
		sli: string;
		timeline: PublicTimelinePoint[];
		days?: number;
	} = $props();

	const visibleTimeline = $derived.by(() => (days ? timeline.slice(-days) : timeline));
</script>

<div class="flex flex-col gap-2">
	<div class="mt-3 grid grid-flow-col auto-cols-fr gap-1">
		{#each visibleTimeline as point (point.day)}
			<HistoricalTimelinePoint {point} />
		{/each}
	</div>
	<div class="relative flex justify-between gap-2 text-xs text-foreground/70
            before:content-[''] before:absolute before:left-0 before:right-0
            before:top-1/2 before:h-px before:bg-foreground/30">
  <p class="relative z-10 bg-background pr-2">{days} days</p>
  <p class="relative z-10 bg-background px-2">{sli}</p>
  <p class="relative z-10 bg-background pl-2">now</p>
</div>
</div>
