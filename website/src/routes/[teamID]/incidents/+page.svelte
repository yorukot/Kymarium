<script lang="ts">
  import Icon from '@iconify/svelte';
  import * as Card from '$lib/components/ui/card/index.js';
  import Badge from '$lib/components/ui/badge/badge.svelte';
  import { Button } from '$lib/components/ui/button';
  import { Separator } from '$lib/components/ui/separator/index.js';
  import type { Incident } from '../../../lib/types/index.js';
  import {
    formatIncidentDate,
    formatIncidentDuration,
    severityBadgeClass,
    severityLabel,
    statusLabel,
    statusTone
  } from '$lib/styles/incident';
  import type { IncidentWithMonitors } from './+page';

  /** @type {import('./$types').PageProps} */
  let { data } = $props();

  const incidents = $derived<IncidentWithMonitors[]>(data.incidents ?? []);
  const openIncidents = $derived(sortIncidents(incidents.filter((i) => i.status !== 'resolved')));
  const resolvedIncidents = $derived(sortIncidents(incidents.filter((i) => i.status === 'resolved')));

  function sortIncidents<T extends Incident>(list: T[]): T[] {
    return [...list].sort((a, b) => {
      const aTime = new Date(a.resolvedAt ?? a.startedAt).getTime();
      const bTime = new Date(b.resolvedAt ?? b.startedAt).getTime();
      return bTime - aTime;
    });
  }

</script>

<div class="flex flex-col gap-6">
  <header class="flex flex-col gap-3">
    <div class="flex items-start justify-between gap-3 flex-wrap">
      <div class="flex flex-col gap-1">
        <p class="text-sm text-muted-foreground">Team incidents</p>
        <div class="flex items-center gap-2 flex-wrap">
          <h1 class="text-2xl font-bold">Incidents</h1>
          <Badge variant="secondary">Manual</Badge>
        </div>
      </div>
      <Button href="incidents/new">
        <Icon icon="lucide:plus" />
        Create incident
      </Button>
    </div>
    <p class="text-sm text-muted-foreground">
      Open incidents are shown first. Resolved ones appear below with their duration.
    </p>
  </header>

  <section class="flex flex-col gap-3">
    {@render SectionHeading({ label: 'Open', count: openIncidents.length, icon: 'lucide:alert-octagon' })}
    {#if openIncidents.length === 0}
      {@render EmptyState({ message: 'No open incidents' })}
    {:else}
      {@render IncidentList({ items: openIncidents })}
    {/if}
  </section>

  <Separator class="my-2" />

  <section class="flex flex-col gap-3">
    {@render SectionHeading({ label: 'Resolved', count: resolvedIncidents.length, icon: 'lucide:check-circle-2' })}
    {#if resolvedIncidents.length === 0}
      {@render EmptyState({ message: 'No resolved incidents yet' })}
    {:else}
      {@render IncidentList({ items: resolvedIncidents, subdued: true })}
    {/if}
  </section>
</div>

<!-- Components -->
{#snippet SectionHeading({ label, count, icon }: { label: string; count: number; icon: string })}
  <div class="flex items-center gap-2">
    <Icon icon={icon} class="size-5" />
    <h2 class="text-lg font-semibold">{label}</h2>
    <Badge variant="secondary">{count}</Badge>
  </div>
{/snippet}

{#snippet EmptyState({ message }: { message: string })}
  <Card.Root class="p-6 text-center text-muted-foreground border-dashed">
    {message}
  </Card.Root>
{/snippet}

{#snippet IncidentList({ items, subdued = false }: { items: IncidentWithMonitors[]; subdued?: boolean })}
  <div class="flex flex-col gap-3">
    {#each items as incident (incident.id)}
      {@render IncidentCard({ incident, subdued })}
    {/each}
  </div>
{/snippet}

{#snippet IncidentCard({ incident, subdued = false }: { incident: IncidentWithMonitors; subdued?: boolean })}
  <Card.Root class={`p-4 flex flex-col gap-3 ${subdued ? 'opacity-80' : ''}`}>
    <div class="flex items-start gap-3 flex-wrap">
      <div class="flex-1 min-w-0 flex flex-col gap-1">
        <div class="flex items-center gap-2 flex-wrap">
          <p class="text-lg font-semibold truncate">
            {incident.title ?? `Incident #${incident.id}`}
          </p>
          {#if incident.title}
            <span class="text-sm text-muted-foreground">#{incident.id}</span>
          {/if}
          <Badge class={severityBadgeClass(incident.severity)}>
            {severityLabel(incident.severity)}
          </Badge>
          <span class={`text-sm font-medium ${statusTone(incident.status)}`}>
            {statusLabel(incident.status)}
          </span>
          {#if incident.isPublic}
            <Badge variant="outline" class="gap-1">
              <Icon icon="lucide:globe-2" class="size-3.5" /> Public
            </Badge>
          {/if}
          {#if incident.autoResolve}
            <Badge variant="outline" class="gap-1">
              <Icon icon="lucide:clock-3" class="size-3.5" /> Auto-resolve
            </Badge>
          {/if}
        </div>
        <div class="text-sm text-muted-foreground flex gap-3 flex-wrap">
          <span class="flex items-center gap-1">
            <Icon icon="lucide:play-circle" class="size-4" />
            Started {formatIncidentDate(incident.startedAt)}
          </span>
          <span class="flex items-center gap-1">
            <Icon icon="lucide:flag" class="size-4" />
            Resolved {formatIncidentDate(incident.resolvedAt)}
          </span>
          <span class="flex items-center gap-1">
            <Icon icon="lucide:timer" class="size-4" />
            Duration {formatIncidentDuration(incident)}
          </span>
          <span class="flex items-center gap-1">
            <Icon icon="lucide:clock" class="size-4" />
            Updated {formatIncidentDate(incident.updatedAt)}
          </span>
        </div>
      </div>
      <div class="flex items-center gap-2">
        <Button size="sm" variant="ghost" href={`incidents/${incident.id}`}>
          View incident
          <Icon icon="lucide:arrow-right" />
        </Button>
      </div>
    </div>

    <Card.Content class="p-0 flex flex-col gap-2">
      <div class="flex items-center gap-2 text-sm font-medium">
        <Icon icon="lucide:monitor" class="size-4" />
        Monitors ({incident.monitorNames.length})
      </div>
      {#if incident.monitorNames.length === 0}
        <p class="text-sm text-muted-foreground">No monitor links found for this incident.</p>
      {:else}
        <div class="flex flex-wrap gap-2">
          {#each incident.monitorNames as name (name)}
            <Badge variant="secondary" class="truncate max-w-[14rem]">
              {name}
            </Badge>
          {/each}
        </div>
      {/if}
    </Card.Content>
  </Card.Root>
{/snippet}
