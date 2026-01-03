<script lang="ts">
  import Button from '$lib/components/ui/button/button.svelte';
  import * as Card from '$lib/components/ui/card';
  import * as InputGroup from '$lib/components/ui/input-group';
  import * as Select from '$lib/components/ui/select';
  import { Separator } from '$lib/components/ui/separator';
  import Icon from '@iconify/svelte';
  import { flip } from 'svelte/animate';
  import { dndzone, SHADOW_ITEM_MARKER_PROPERTY_NAME } from 'svelte-dnd-action';
  import type { StatusPageElement, StatusPageElementType } from '$lib/types';
  import StatusPageMonitorRow from './monitor-raw-status-page.svelte';

  const {
    group,
    namePrefix,
    onDeleteGroup,
    onDeleteMonitor,
    dndType,
    flipDurationMs,
    onMonitorsDndConsider,
    onMonitorsDndFinalize
  }: {
    group: StatusPageElement;
    namePrefix: string;
    onDeleteGroup?: (groupId: string) => void;
    onDeleteMonitor?: (monitorId: string) => void;
    dndType: string;
    flipDurationMs: number;
    onMonitorsDndConsider?: (groupId: string, event: CustomEvent) => void;
    onMonitorsDndFinalize?: (groupId: string, event: CustomEvent) => void;
  } = $props();

  let typeValue = $derived<StatusPageElementType>(group.type);
  let typeInput = $state<HTMLInputElement | null>(null);

  $effect(() => {
    typeValue = group.type;
  });

  $effect(() => {
    if (!typeInput) return;
    typeInput.value = typeValue;
    typeInput.dispatchEvent(new Event('input', { bubbles: true }));
    typeInput.dispatchEvent(new Event('change', { bubbles: true }));
  });

  type DndMonitor = (StatusPageElement['monitors'][number]) & {
    [SHADOW_ITEM_MARKER_PROPERTY_NAME]?: boolean;
  };

  const dndMonitors = $derived((group.monitors ?? []) as DndMonitor[]);

  const typeLabel = (t: StatusPageElementType) =>
    t === 'historical_timeline' ? 'Historical Timeline' : 'Only Current Status';

  const typeIcon = (t: StatusPageElementType) =>
    t === 'historical_timeline' ? 'lucide:chart-line' : 'lucide:circle-small';
</script>

<Card.Root class="bg-muted p-0">
  <Card.Content class="p-0">
    <div class="flex justify-between items-center p-2 gap-2">
      <Icon icon="lucide:grip-vertical" class="size-4 shrink-0" />
      <InputGroup.Root class="w-full">
        <InputGroup.Input
          name={`${namePrefix}.name`}
          placeholder="Please enter element name"
        />
        <input type="hidden" name={`${namePrefix}.monitor`} value="false" />
        <input type="hidden" name={`${namePrefix}.sortOrder`} value={group.sortOrder} />
        <InputGroup.Addon class="hidden sm:block">
          <Icon icon="lucide:layers" />
        </InputGroup.Addon>
      </InputGroup.Root>

      <div class="flex items-center gap-2">
        <Select.Root type="single" bind:value={typeValue}>
          <Select.Trigger class="lg:w-51">
            <Icon icon={typeIcon(typeValue)} />
            <p class="hidden lg:block">{typeLabel(typeValue)}</p>
          </Select.Trigger>
          <Select.Content>
            <Select.Group>
              <Select.Item value="historical_timeline" label="Historical Timeline">
                <Icon icon="lucide:chart-line" /> Historical Timeline
              </Select.Item>
              <Select.Item value="current_status_indicator" label="Only Current Status">
                <Icon icon="lucide:circle-small" /> Only Current Status
              </Select.Item>
            </Select.Group>
          </Select.Content>
        </Select.Root>
        <input type="hidden" name={`${namePrefix}.type`} bind:this={typeInput} />

        <Button
          size="icon"
          variant="destructive"
          onclick={() => onDeleteGroup?.(group.id)}
        >
          <Icon icon="lucide:trash" />
        </Button>
      </div>
    </div>

    <Separator />

    <div class="relative">
      <div
        class="flex flex-col gap-2 min-h-10"
        use:dndzone={{ items: dndMonitors, flipDurationMs, type: dndType }}
        onconsider={(event) => onMonitorsDndConsider?.(group.id, event)}
        onfinalize={(event) => onMonitorsDndFinalize?.(group.id, event)}
      >
        {#if (group.monitors ?? []).length === 0}
          <div class="p-2 pl-6 text-sm opacity-70">No monitor yet.</div>
        {/if}
        {#each dndMonitors as m, monitorIndex (m.id)}
          <div
            animate:flip={{ duration: flipDurationMs }}
            data-is-dnd-shadow-item-hint={m[SHADOW_ITEM_MARKER_PROPERTY_NAME]}
          >
            {#if m[SHADOW_ITEM_MARKER_PROPERTY_NAME]}
              <div class="h-10 rounded-md border border-dashed bg-muted/40"></div>
            {:else}
              <StatusPageMonitorRow
                monitor={m}
                namePrefix={`${namePrefix}.monitors[${monitorIndex}]`}
                onDelete={onDeleteMonitor}
              />
            {/if}
          </div>
        {/each}
      </div>
    </div>
  </Card.Content>
</Card.Root>
