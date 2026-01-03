<script lang="ts">
	import { formatSli, statusMeta } from '$lib/styles/status';
	import type { PublicStatusPageElement } from '$lib/types';
	import Icon from '@iconify/svelte';
	import HistoricalTimeline from './historical-timeline.svelte';
	import HistoricalMonitor from './historical-monitor.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import { slide } from 'svelte/transition';

	let {
		group,
		days = 90
	}: {
		group: PublicStatusPageElement;
		days?: number;
	} = $props();

	let open = $state(false);

	function selectSli() {
		if (days <= 30) return formatSli(group.uptimeSli30);
		if (days <= 60) return formatSli(group.uptimeSli60);
		return formatSli(group.uptimeSli90);
	}
</script>

<div class="rounded-md border bg-background/70 p-3">
	<div class="flex items-center justify-between gap-3">
		<h2 class="text-lg font-semibold flex items-center gap-1">
			<Button variant="ghost" size="icon" class="p-0 size-6" onclick={() => (open = !open)}>
				<Icon icon={open ? 'lucide:minus' : 'lucide:plus'} class="text-foreground/70 size-5" />
			</Button>
			{group.name}
		</h2>
			<span class={`text-sm font-medium ${statusMeta[group.status ?? 'up'].tone}`}>
				{statusMeta[group.status ?? 'up'].label}
			</span>
	</div>
	{#if group.type === 'historical_timeline' && group.timeline?.length}
		<HistoricalTimeline sli={selectSli()} timeline={group.timeline} {days} />
	{/if}

	{#if open}
		<div
			class="mt-3 flex flex-col gap-3"
			in:slide={{ duration: 200 }}
			out:slide={{ duration: 200 }}
		>
			{#each group.monitors as monitor (monitor.id)}
				<HistoricalMonitor {monitor} {days} />
			{/each}
		</div>
	{/if}
</div>
