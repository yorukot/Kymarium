<script lang="ts">
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Card from '$lib/components/ui/card';
	import * as Select from '$lib/components/ui/select';
	import Icon from '@iconify/svelte';
	import type { createForm } from 'felte';
	import { flip } from 'svelte/animate';
	import { dndzone, SHADOW_ITEM_MARKER_PROPERTY_NAME } from 'svelte-dnd-action';

	import type { MonitorWithIncidents, StatusPageElement, StatusPageMonitor } from '$lib/types';

	import StatusPageMonitorCard from './monitor-raw-status-page.svelte';
	import StatusPageGroupCard from './group-raw-status-page.svelte';
	import type { StatusPageUpsertValues } from './schema';

	type FelteReturn = ReturnType<typeof createForm<StatusPageUpsertValues>>;
	type SetFields = FelteReturn['setFields'];

	type DndElement = StatusPageElement & {
		[SHADOW_ITEM_MARKER_PROPERTY_NAME]?: boolean;
	};
	type DndMonitor = StatusPageMonitor & {
		[SHADOW_ITEM_MARKER_PROPERTY_NAME]?: boolean;
	};
	type DndItem = DndElement | DndMonitor;

	const flipDurationMs = 150;
	const dndType = 'status-page-elements';

	let monitorID = $state('');
	let groupID = $state('');

	let {
		availableMonitors,
		statusPageId,
		elements,
		setFields
	}: {
		availableMonitors: MonitorWithIncidents[];
		statusPageId: string;
		elements: StatusPageElement[];
		setFields: SetFields;
	} = $props();

	const sortByOrder = (
		a: { sortOrder: number; id: string },
		b: { sortOrder: number; id: string }
	) => {
		const diff = a.sortOrder - b.sortOrder;
		if (diff !== 0) return diff;
		return a.id.localeCompare(b.id);
	};

	function normalizeElements(list: StatusPageElement[]) {
		return list
			.map((element) => ({
				...element,
				monitors: (element.monitors ?? []).slice().sort(sortByOrder)
			}))
			.sort(sortByOrder);
	}

	let editableElements = $derived<StatusPageElement[]>(
		normalizeElements(structuredClone(elements ?? []))
	);

	const dndElements = $derived(editableElements as DndElement[]);

	const monitorTriggerContent = $derived(
		availableMonitors.find((m) => m.id === monitorID)?.name ?? 'Select a monitor'
	);

	const groupTriggerContent = $derived(
		editableElements.find((g) => !g.monitor && g.id === groupID)?.name ?? 'Ungrouped'
	);

	const isShadowItem = (item: DndItem) => Boolean(item?.[SHADOW_ITEM_MARKER_PROPERTY_NAME]);
	const int64Pattern = /^-?\d+$/;

	const isInt64String = (value?: string | null): value is string =>
		typeof value === 'string' && int64Pattern.test(value);

	function stripShadowItems(list: StatusPageElement[]) {
		return list
			.filter((item) => !isShadowItem(item))
			.map((element) => ({
				...element,
				monitors: (element.monitors ?? []).filter((m) => !isShadowItem(m))
			}));
	}

	function applySortOrders(list: StatusPageElement[]) {
		return list.map((element, index) => {
			const base = { ...element, sortOrder: index + 1 };
			if (element.monitor) {
				return { ...base, monitors: [] };
			}
			const monitors = (element.monitors ?? []).map((monitor, monitorIndex) => ({
				...monitor,
				groupId: element.id,
				sortOrder: monitorIndex + 1
			}));
			return { ...base, monitors };
		});
	}

	function previewElements(next: StatusPageElement[]) {
		editableElements = next;
	}

	function commitElements(next: StatusPageElement[]) {
		const committed = applySortOrders(stripShadowItems(next));
		editableElements = committed;
		setFields('elements', sanitizeElementsForForm(committed), true);
	}

	function sanitizeElementsForForm(list: StatusPageElement[]) {
		return list.map((element) => {
			const monitorId = isInt64String(element.monitorId ?? null) ? element.monitorId : undefined;
			const base: StatusPageElement = {
				...element,
				id: isInt64String(element.id) ? element.id : (undefined as unknown as string),
				monitorId: element.monitor ? monitorId : undefined,
				monitors: []
			};
			if (element.monitor) {
				return base;
			}
			const monitors = (element.monitors ?? []).map((monitor) => ({
				...monitor,
				id: isInt64String(monitor.id) ? monitor.id : (undefined as unknown as string),
				groupId: isInt64String(monitor.groupId ?? null) ? monitor.groupId : undefined
			}));
			return { ...base, monitors };
		});
	}

	function mergeGroupMonitors(next: StatusPageElement[]) {
		return next.map((element) => {
			if (element.monitor) return element;
			const current = editableElements.find((entry) => !entry.monitor && entry.id === element.id);
			return {
				...element,
				monitors: current?.monitors ?? element.monitors ?? []
			};
		});
	}

	function toTopLevelItems(items: DndItem[]): StatusPageElement[] {
		return items.map((item) => {
			if (isShadowItem(item)) return item as StatusPageElement;
			if ('monitor' in item) {
				return {
					...item,
					monitors: item.monitors ?? []
				};
			}
			return {
				id: item.id,
				statusPageId,
				name: item.name ?? 'New Monitor',
				type: item.type ?? 'historical_timeline',
				sortOrder: item.sortOrder ?? 0,
				monitor: true,
				monitorId: item.monitorId ?? null,
				monitors: []
			};
		});
	}

	function toGroupItems(groupId: string, items: DndItem[]): StatusPageMonitor[] {
		return items.flatMap((item) => {
			if (isShadowItem(item)) return [item as StatusPageMonitor];
			if ('monitor' in item && item.monitor === false) return [];
			const monitorId = item.monitorId ?? null;
			if (!monitorId) return [];
			return [
				{
					id: item.id,
					statusPageId,
					monitorId,
					groupId,
					name: item.name ?? 'New Monitor',
					type: item.type ?? 'historical_timeline',
					sortOrder: item.sortOrder ?? 0
				}
			];
		});
	}

	function updateGroupMonitors(groupId: string, monitors: StatusPageMonitor[], commit: boolean) {
		const nextElements = editableElements.map((element) => {
			if (element.monitor || element.id !== groupId) return element;
			return {
				...element,
				monitors
			};
		});
		if (commit) {
			commitElements(nextElements);
		} else {
			previewElements(nextElements);
		}
	}

	const onDeleteMonitor = (id: string) => {
		commitElements(editableElements.filter((element) => element.id !== id));
	};

	const onDeleteGroup = (id: string) => {
		commitElements(editableElements.filter((element) => element.id !== id));
	};

	const onDeleteMonitorInGroup = (id: string) => {
		const nextElements = editableElements.map((element) => {
			if (element.monitor) return element;
			return {
				...element,
				monitors: (element.monitors ?? []).filter((monitor) => monitor.id !== id)
			};
		});
		commitElements(nextElements);
	};

	function nextTopLevelSortOrder() {
		const max = editableElements.reduce((acc, element) => Math.max(acc, element.sortOrder ?? 0), 0);
		return max + 1;
	}

	function nextSortOrderInGroup(groupId: string) {
		const group = editableElements.find((element) => !element.monitor && element.id === groupId);
		const list = group?.monitors ?? [];
		const max = list.reduce((acc, m) => Math.max(acc, m.sortOrder ?? 0), 0);
		return max + 1;
	}

	const onAddMonitor = () => {
		if (!monitorID) return;

		const targetGroupId = groupID === '0' || groupID === '' ? null : groupID;

		if (!targetGroupId) {
			const newElement: StatusPageElement = {
				id: crypto.randomUUID(),
				statusPageId,
				name: availableMonitors.find((m) => m.id === monitorID)?.name ?? 'New Monitor',
				type: 'historical_timeline',
				sortOrder: nextTopLevelSortOrder(),
				monitor: true,
				monitorId: monitorID,
				monitors: []
			};

			commitElements([...editableElements, newElement]);
		} else {
			commitElements(
				editableElements.map((element) => {
					if (element.monitor || element.id !== targetGroupId) return element;
					const nextMonitor: StatusPageMonitor = {
						id: crypto.randomUUID(),
						statusPageId,
						monitorId: monitorID,
						groupId: targetGroupId,
						name: availableMonitors.find((m) => m.id === monitorID)?.name ?? 'New Monitor',
						type: 'historical_timeline',
						sortOrder: nextSortOrderInGroup(targetGroupId)
					};
					return {
						...element,
						monitors: [...(element.monitors ?? []), nextMonitor]
					};
				})
			);
		}

		monitorID = '';
		groupID = '';
	};

	const onAddGroup = () => {
		const newGroup: StatusPageElement = {
			id: crypto.randomUUID(),
			statusPageId,
			name: 'New Group',
			type: 'historical_timeline',
			sortOrder: nextTopLevelSortOrder(),
			monitor: false,
			monitorId: null,
			monitors: []
		};

		commitElements([...editableElements, newGroup]);
	};

	const onTopLevelConsider = (event: CustomEvent) => {
		const items = mergeGroupMonitors(toTopLevelItems(event.detail.items as DndItem[]));
		previewElements(items);
	};

	const onTopLevelFinalize = (event: CustomEvent) => {
		const items = mergeGroupMonitors(toTopLevelItems(event.detail.items as DndItem[]));
		commitElements(items);
	};

	const onGroupConsider = (groupId: string, event: CustomEvent) => {
		const items = toGroupItems(groupId, event.detail.items as DndItem[]);
		updateGroupMonitors(groupId, items, false);
	};

	const onGroupFinalize = (groupId: string, event: CustomEvent) => {
		const items = toGroupItems(groupId, event.detail.items as DndItem[]);
		updateGroupMonitors(groupId, items, true);
	};
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>Status page elements</Card.Title>
		<Card.Description>Configure the status page's elements.</Card.Description>
	</Card.Header>

	<Card.Content class="flex flex-col gap-4">
		<!-- top selectors -->
		<div class="flex gap-2 flex-col sm:flex-row warp">
			<Select.Root type="single" name="monitor" bind:value={monitorID}>
				<Select.Trigger class="max-w-lg w-full">{monitorTriggerContent}</Select.Trigger>
				<Select.Content>
					<Select.Group>
						<Select.Label>Monitors</Select.Label>
						{#each availableMonitors as monitor (monitor.id)}
							<Select.Item value={monitor.id} label={monitor.name}>
								{monitor.name}
							</Select.Item>
						{/each}
					</Select.Group>
				</Select.Content>
			</Select.Root>

			<div class="flex gap-2 w-full">
				<Select.Root type="single" name="group" bind:value={groupID}>
					<Select.Trigger class="max-w-lg w-full truncate">{groupTriggerContent}</Select.Trigger>
					<Select.Content>
						<Select.Group>
							<Select.Label>Groups</Select.Label>
							<Select.Item value="0" label="Ungrouped">Ungrouped</Select.Item>
							{#each editableElements.filter((element) => !element.monitor) as group (group.id)}
								<Select.Item value={group.id} label={group.name}>
									{group.name}
								</Select.Item>
							{/each}
						</Select.Group>
					</Select.Content>
				</Select.Root>

				<Button size="icon" onclick={onAddMonitor}>
					<Icon icon="lucide:plus" />
				</Button>
			</div>
		</div>

		<!-- elements list -->
		<div
			class="flex flex-col gap-2"
			use:dndzone={{ items: dndElements, flipDurationMs, type: dndType }}
			onconsider={onTopLevelConsider}
			onfinalize={onTopLevelFinalize}
		>
			{#each dndElements as element, i (element.id)}
				<div
					animate:flip={{ duration: flipDurationMs }}
					data-is-dnd-shadow-item-hint={element[SHADOW_ITEM_MARKER_PROPERTY_NAME]}
				>
					{#if element[SHADOW_ITEM_MARKER_PROPERTY_NAME]}
						<div class="h-12 rounded-md border border-dashed bg-muted/40"></div>
					{:else if element.monitor}
						<Card.Root class="bg-muted p-0">
							<StatusPageMonitorCard
								monitor={element}
								namePrefix={`elements[${i}]`}
								isElement
								onDelete={onDeleteMonitor}
							/>
						</Card.Root>
					{:else}
						<StatusPageGroupCard
							group={element}
							namePrefix={`elements[${i}]`}
							{onDeleteGroup}
							onDeleteMonitor={onDeleteMonitorInGroup}
							{dndType}
							{flipDurationMs}
							onMonitorsDndConsider={onGroupConsider}
							onMonitorsDndFinalize={onGroupFinalize}
						/>
					{/if}
				</div>
			{/each}
		</div>

		<div class="flex justify-end">
			<Button onclick={onAddGroup}>
				<Icon icon="lucide:plus" />
				New Group
			</Button>
		</div>
	</Card.Content>
</Card.Root>
