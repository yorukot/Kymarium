<script lang="ts">
	import { Tooltip } from 'bits-ui';
	import { formatDate, formatUpTo2Decimals, timelineTone } from '$lib/styles/status';
	import type { PublicTimelinePoint } from '../../../types';

	let { point }: {
		point: PublicTimelinePoint;
	} = $props();

	const title = $derived.by(
		() => `${formatDate(point.day)} · ${point.success} success / ${point.fail} fail`
	);
</script>

<Tooltip.Root delayDuration={0}>
	<Tooltip.Trigger>
		<div class={`h-6 rounded-none ${timelineTone(point)}`} {title}></div>
	</Tooltip.Trigger>
	<Tooltip.Content
		side="top"
		sideOffset={8}
		class="bg-popover text-popover-foreground w-64 max-w-[calc(100vw-1rem)] rounded-md border p-4 text-sm shadow-md"
	>
		<div class="flex flex-col space-y-2">
			<div class="flex justify-between space-x-4">
				<span class="font-medium text-foreground">{formatDate(point.day)}</span>
				<span class="font-medium text-foreground">
					{#if point.fail + point.success === 0}
						<span class="font-medium text-foreground/50">No data</span>
					{:else}
						<span class="font-medium text-foreground/50">
							{formatUpTo2Decimals((point.success / (point.success + point.fail)) * 100)}%
							uptime
						</span>
					{/if}
				</span>
			</div>
			<div>
				{#if point.fail === 0}
					<span class="font-medium text-foreground/50">No down recoard found</span>
				{:else}
					<span class="font-medium text-foreground/50">
						{Math.round((point.success / (point.success + point.fail)) * 100)}% uptime
					</span>
				{/if}
			</div>
		</div>
	</Tooltip.Content>
</Tooltip.Root>
