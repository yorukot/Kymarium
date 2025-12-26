<script lang="ts">
	import Button from '$lib/components/ui/button/button.svelte';
	import * as InputGroup from '$lib/components/ui/input-group';
	import * as Select from '$lib/components/ui/select';
	import Icon from '@iconify/svelte';
	import type { StatusPageElement, StatusPageMonitor, StatusPageElementType } from '$lib/types';

	const {
		monitor,
		namePrefix,
		isElement = false,
		onDelete
	}: {
		monitor: StatusPageMonitor | StatusPageElement;
		namePrefix: string;
		isElement?: boolean;
		onDelete?: (monitorId: string) => void;
	} = $props();

	let typeValue = $derived<StatusPageElementType>(monitor.type);
	let typeInput = $state<HTMLInputElement | null>(null);

	$effect(() => {
		typeValue = monitor.type;
	});

	$effect(() => {
		if (!typeInput) return;
		typeInput.value = typeValue;
		typeInput.dispatchEvent(new Event('input', { bubbles: true }));
		typeInput.dispatchEvent(new Event('change', { bubbles: true }));
	});

	const typeLabel = (t: StatusPageElementType) =>
		t === 'historical_timeline' ? 'Historical Timeline' : 'Only Current Status';

	const typeIcon = (t: StatusPageElementType) =>
		t === 'historical_timeline' ? 'lucide:chart-line' : 'lucide:circle-small';
</script>

<div class="flex justify-between items-center p-2 gap-2">
  <Icon icon="lucide:grip-vertical" class="size-4 shrink-0" />

	<InputGroup.Root class="w-full">
		<InputGroup.Input
			name={`${namePrefix}.name`}
			placeholder="Please enter element name"
		/>
		<input type="hidden" name={`${namePrefix}.sortOrder`} value={monitor.sortOrder} />
		{#if isElement}
			<input type="hidden" name={`${namePrefix}.monitor`} value="true" />
		{/if}
		{#if monitor.monitorId}
			<input type="hidden" name={`${namePrefix}.monitorId`} value={monitor.monitorId} />
		{/if}
		<InputGroup.Addon class="hidden sm:block">
			<Icon icon="lucide:activity" />
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

		<Button size="icon" variant="destructive" onclick={() => onDelete?.(monitor.id)}>
			<Icon icon="lucide:trash" />
		</Button>
	</div>
</div>
