<script lang="ts" module>
	export type MultiSelectOption = {
		label: string;
		value: string;
		keywords?: string[];
		disabled?: boolean;
		icon?: string;
	};
</script>

<script lang="ts">
	import { cn } from '$lib/utils.js';
	import { Popover, PopoverContent, PopoverTrigger } from '../popover';
	import { Command, CommandEmpty, CommandGroup, CommandItem, CommandList } from '../command';
	import Badge from '../badge/badge.svelte';
	import ChevronsUpDownIcon from '@lucide/svelte/icons/chevrons-up-down';
	import XIcon from '@lucide/svelte/icons/x';
	import Icon from '@iconify/svelte';
	import Checkbox from '../checkbox/checkbox.svelte';

	type MultiSelectProps = {
		value?: string[];
		open?: boolean;
		options?: MultiSelectOption[];
		name?: string;
		placeholder?: string;
		emptyMessage?: string;
		class?: string;
		disabled?: boolean;
		clearable?: boolean;
		maxBadges?: number;
		closeOnSelect?: boolean;
	};

	let {
		value = $bindable<string[]>([]),
		open = $bindable(false),
		options = [] as MultiSelectOption[],
		name,
		placeholder = 'Select options',
		emptyMessage = 'No results found.',
		class: className,
		disabled = false,
		clearable = true,
		maxBadges = 3,
		closeOnSelect = false
	}: MultiSelectProps = $props();

	const selectedOptions = $derived(
		options.filter((option: MultiSelectOption) => value?.includes(option.value))
	);
	const remainingCount = $derived(Math.max(0, (value?.length || 0) - maxBadges));

	const optionOrder = $derived(
		new Map(options.map((option: MultiSelectOption, index: number) => [option.value, index]))
	);

	function toggle(optionValue: string) {
		if (disabled) return;

		if (!value) value = [];

		const exists = value.includes(optionValue);
		const next = exists
			? value.filter((item: string) => item !== optionValue)
			: [...value, optionValue];

		value = next
			.slice()
			.sort(
				(a: string, b: string) =>
					(optionOrder.get(a) ?? Number.POSITIVE_INFINITY) -
					(optionOrder.get(b) ?? Number.POSITIVE_INFINITY)
			);

		if (closeOnSelect && !exists) {
			open = false;
		}
	}

	function clearSelection(event?: Event) {
		event?.preventDefault();
		event?.stopPropagation();
		if (disabled) return;
		value = [];
	}
</script>

{#if name}
	{#each value ?? [] as selected (selected)}
		<input type="hidden" {name} value={selected} />
	{/each}
{/if}

<Popover bind:open>
	<PopoverTrigger
		{disabled}
		data-slot="multi-select-trigger"
		class={cn(
			'border-input shadow-xs focus-visible:border-ring focus-visible:ring-ring/50 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive flex min-h-10 w-full items-center justify-between gap-2 rounded-md border bg-transparent px-3 py-2 text-left text-sm outline-none transition-[color,box-shadow] focus-visible:ring-[3px] disabled:cursor-not-allowed disabled:opacity-50',
			selectedOptions.length === 0 && 'text-muted-foreground',
			className
		)}
	>
		<div class="flex flex-1 flex-wrap items-center gap-2">
			{#if selectedOptions.length === 0}
				<span class="text-muted-foreground">{placeholder}</span>
			{:else}
				{#each selectedOptions.slice(0, maxBadges) as option (option.value)}
					<Badge variant="secondary" class="flex items-center gap-1">
						{#if option.icon}
							<Icon icon={option.icon} class="size-3.5 text-muted-foreground" aria-hidden="true" />
						{/if}
						{option.label}
					</Badge>
				{/each}
				{#if remainingCount > 0}
					<Badge variant="outline" class="text-xs text-muted-foreground">
						+{remainingCount} more
					</Badge>
				{/if}
			{/if}
		</div>
		<div class="flex items-center gap-1">
			{#if clearable && (value?.length ?? 0) > 0}
				<button
					type="button"
					onclick={clearSelection}
					class="hover:bg-muted text-muted-foreground focus-visible:border-ring focus-visible:ring-ring/50 rounded-sm p-1 outline-none transition focus-visible:ring-[2px]"
					aria-label="Clear selection"
				>
					<XIcon class="size-4" />
				</button>
			{/if}
			<ChevronsUpDownIcon class="text-muted-foreground size-4" aria-hidden="true" />
		</div>
	</PopoverTrigger>

	<PopoverContent class="p-0 w-(--bits-popover-anchor-width) min-w-(--bits-popover-anchor-width)">
		<Command class="border-0 shadow-none">
			<CommandList>
				<CommandEmpty>{emptyMessage}</CommandEmpty>
				<CommandGroup>
					{#each options as option (option.value)}
						<CommandItem
							value={option.label}
							keywords={option.keywords ?? [option.value]}
							disabled={option.disabled}
							onSelect={() => toggle(option.value)}
							data-selected={value?.includes(option.value)}
							class="gap-2"
						>
							<Checkbox
								checked={value?.includes(option.value)}
								disabled={option.disabled}
								class="pointer-events-none size-5"
							/>
							<div class="flex items-center gap-2 leading-tight">
								{#if option.icon}
									<Icon
										icon={option.icon}
										class="size-4 text-muted-foreground"
										aria-hidden="true"
									/>
								{/if}
								<span class="leading-tight">{option.label}</span>
							</div>
						</CommandItem>
					{/each}
				</CommandGroup>
			</CommandList>
		</Command>
	</PopoverContent>
</Popover>
