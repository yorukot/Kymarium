<script lang="ts">
	import type { PublicIncident, IncidentEvent } from '$lib/types';
	import {
		eventTypeLabel,
		formatIncidentDate,
		statusLabel,
		statusTone
	} from '$lib/styles/incident';
	import Icon from '@iconify/svelte';

	let { incidents = [], monitorNameById = {} }: {
		incidents?: PublicIncident[];
		monitorNameById?: Record<string, string>;
	} = $props();

	const openIncidents = $derived.by(() =>
		incidents.filter((incident) => incident.status !== 'resolved')
	);

	const orderedIncidents = $derived.by(() =>
		[...openIncidents].sort(
			(a, b) => new Date(b.startedAt).getTime() - new Date(a.startedAt).getTime()
		)
	);

	function orderTimeline(timeline?: IncidentEvent[]) {
		if (!timeline) return [];
		return [...timeline]
			.filter((event) => Boolean(event.createdBy) && event.message?.trim().length)
			.sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime());
	}
</script>

<section class="flex flex-col gap-4">
	{#if orderedIncidents.length >= 0}
		<div class="flex flex-col gap-4">
			{#each orderedIncidents as incident (incident.id)}
				<div class="p-0 bg-background/70 border-destructive border rounded">
					<div class="flex flex-col">
						<div
							class="flex flex-wrap items-start justify-between gap-3 bg-destructive rounded-t-xs p-2"
						>
							<div class="flex flex-col gap-1">
								<div class="flex flex-wrap items-center">
									<p class="text-xl font-semibold">
										{incident.title ?? `Incident #${incident.id}`}
									</p>
									<span class={`text-xs font-medium ${statusTone(incident.status)}`}>
										{statusLabel(incident.status)}
									</span>
								</div>
							</div>
						</div>
						{#if orderTimeline(incident.timeline).length === 0}
							<p class="mt-2 text-sm text-muted-foreground">No timeline updates yet.</p>
						{:else}
							<div class="space-y-4 p-4">
								{#each orderTimeline(incident.timeline) as event (event.id)}
									<div class="relative">
										<div class="flex flex-col gap-1">
											<div class="flex flex-col gap-2">
												<div class="flex flex-wrap items-center gap-2">
														<p class="font-semibold">
															{eventTypeLabel(event.eventType)}:
														</p>
														<span class="text-sm font-medium">
															{event.message || '—'}
														</span>
												</div>
												<span class="text-xs text-muted-foreground">
													{formatIncidentDate(event.createdAt)}
												</span>
											</div>
										</div>
									</div>
								{/each}
							</div>
						{/if}
						<div class="border-t border-border/60 px-4 py-3">
							<div class="flex items-center gap-2 text-xs text-muted-foreground">
								<Icon icon="lucide:monitor" class="size-3.5" />
								<span>Related monitor</span>
								<span class="text-foreground">
									{monitorNameById[incident.monitorId] ?? `#${incident.monitorId}`}
								</span>
							</div>
						</div>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</section>
