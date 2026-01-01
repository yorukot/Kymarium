<script lang="ts">
	import type { Monitor } from '../../types';
	import { Card } from '$lib/components/ui/card';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import Icon from '@iconify/svelte';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import { deleteMonitor } from '$lib/api/monitor';
	import { toast } from 'svelte-sonner';
	import { page } from '$app/state';
	import DeleteMonitorDialog from './delete-monitor-dialog.svelte';
	import { goto } from '$app/navigation';
	import Button from '../ui/button/button.svelte';
	import { monitorTarget } from '$lib/utils/monitor';
	import { statusMeta } from '$lib/styles/status';

	let { monitors }: { monitors: Monitor[] } = $props();

	let confirmOpen = $state(false);
	let deleting = $state(false);
	let targetMonitor: Monitor | null = $state(null);

	function askDelete(monitor: Monitor) {
		targetMonitor = monitor;
		confirmOpen = true;
	}

	async function handleDelete() {
		if (!targetMonitor) return;
		const teamID = page.params.teamID;
		if (!teamID) {
			toast.error('Missing team id');
			return;
		}

		deleting = true;
		try {
			await deleteMonitor(teamID, targetMonitor.id);
			monitors = monitors.filter((m) => m.id !== targetMonitor?.id);
			toast.success('Monitor deleted');
		} catch (err) {
			const message = err instanceof Error ? err.message : 'Failed to delete monitor';
			toast.error(message);
		} finally {
			deleting = false;
			confirmOpen = false;
			targetMonitor = null;
		}
	}
</script>

<div>
	<div class="flex flex-col gap-2">
		{#if monitors.length === 0}
			<Card class="py-6 px-4 text-sm text-muted-foreground">No monitors yet.</Card>
		{/if}
		{#each monitors as monitor, monitorIndex (monitor.id)}
			<Card id={'monitor-card-' + monitorIndex} class="py-2 px-4">
				<div class="flex justify-between items-center gap-2">
					<div class="flex-1 min-w-0">
						<div class="flex items-center gap-2 min-w-0">
							<span
								class={`w-3 h-3 rounded-full shrink-0 ${statusMeta[monitor.status]?.dot ?? 'bg-muted'}`}
							></span>

							<h2 class="text-md font-semibold truncate min-w-0 flex-1">
								{monitor.name}
							</h2>
						</div>

						<div class="min-w-0">
							<p class="text-sm text-muted-foreground truncate">
								<Badge variant="secondary" class="rounded-sm shrink-0">
									{monitor.type.toUpperCase()}
								</Badge>
								{monitorTarget(monitor)}
							</p>
						</div>
					</div>
					<div>
						<div></div>
						<DropdownMenu.Root>
							<DropdownMenu.Trigger class="shrink-0">
								<Button variant="ghost" size="icon">
									<Icon icon="lucide:more-vertical" />
								</Button>
							</DropdownMenu.Trigger>
							<DropdownMenu.Content>
								<DropdownMenu.Group>
									<DropdownMenu.Item onclick={() => goto(`monitors/${monitor.id}`)}>
										<Icon icon="lucide:eye" /> View details
									</DropdownMenu.Item>
									<DropdownMenu.Item onclick={() => goto(`monitors/${monitor.id}/edit`)}>
										<Icon icon="lucide:edit" /> Edit
									</DropdownMenu.Item>
									<DropdownMenu.Separator />
									<DropdownMenu.Item variant="destructive" onclick={() => askDelete(monitor)}>
										<Icon icon="lucide:trash" />
										Delete
									</DropdownMenu.Item>
								</DropdownMenu.Group>
							</DropdownMenu.Content>
						</DropdownMenu.Root>
					</div>
				</div>
			</Card>
		{/each}
	</div>

	<DeleteMonitorDialog
		bind:open={confirmOpen}
		monitor={targetMonitor}
		onConfirm={handleDelete}
		loading={deleting}
	/>
</div>
