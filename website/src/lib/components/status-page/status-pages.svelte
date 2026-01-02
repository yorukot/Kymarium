<script lang="ts">
	import type { StatusPageWithElements } from '$lib/types';
	import { Card } from '$lib/components/ui/card';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import * as AlertDialog from '$lib/components/ui/alert-dialog';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import Icon from '@iconify/svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import { goto, invalidateAll } from '$app/navigation';
	import { toast } from 'svelte-sonner';
	import { deleteStatusPage } from '$lib/api/status-page';
	import DropdownMenuSeparator from '../ui/dropdown-menu/dropdown-menu-separator.svelte';

	let { statusPages, teamID }: { statusPages: StatusPageWithElements[]; teamID: string } = $props();

	let deleteOpen = $state(false);
	let deleteTarget = $state<StatusPageWithElements | null>(null);
	let deleteLoading = $state(false);

	function countMonitors(item: StatusPageWithElements) {
		return item.elements.reduce((acc, element) => {
			if (element.monitor) return acc + 1;
			return acc + (element.monitors?.length ?? 0);
		}, 0);
	}

	function requestDelete(item: StatusPageWithElements) {
		deleteTarget = item;
		deleteOpen = true;
	}

	async function confirmDelete() {
		if (!deleteTarget) return;
		deleteLoading = true;
		try {
			await deleteStatusPage(teamID, deleteTarget.statusPage.id);
			await invalidateAll();
			toast.success('Status page deleted');
		} catch (error) {
			const message = error instanceof Error ? error.message : 'Failed to delete status page';
			toast.error(message);
		} finally {
			deleteLoading = false;
			deleteOpen = false;
			deleteTarget = null;
		}
	}
</script>

<div>
	<div class="flex flex-col gap-2">
		{#each statusPages as statusPageItem (statusPageItem.statusPage.id)}
			<Card class="py-2 px-4">
				<div class="flex items-center justify-between gap-2">
					<div class="flex-1 min-w-0">
						<div class="flex items-center gap-2 min-w-0">
							<div class="min-w-0 flex-1 flex gap-2 items-center">
								<h2 class="text-md font-semibold truncate">
									{statusPageItem.statusPage.title}
								</h2>
								<a href={`/s/${statusPageItem.statusPage.slug}`} target="_blank" rel="noopener noreferrer" class="flex items-center">
									<Icon
										icon="lucide:external-link"
										class="inline size-4 text-muted-foreground cursor-pointer"
									/>
								</a>
							</div>
							<Badge variant="secondary" class="rounded-sm shrink-0">
								{countMonitors(statusPageItem)} monitors
							</Badge>
						</div>
					</div>
					<div>
						<DropdownMenu.Root>
							<DropdownMenu.Trigger class="shrink-0">
								<Button variant="ghost" size="icon">
									<Icon icon="lucide:more-vertical" />
								</Button>
							</DropdownMenu.Trigger>
							<DropdownMenu.Content>
								<DropdownMenu.Group>
									<DropdownMenu.Item
										onclick={() => goto(`status-pages/${statusPageItem.statusPage.id}/edit`)}
									>
										<Icon icon="lucide:edit" /> Edit
									</DropdownMenu.Item>
									<DropdownMenuSeparator />
									<DropdownMenu.Item
										variant="destructive"
										onclick={() => requestDelete(statusPageItem)}
									>
										<Icon icon="lucide:trash-2" /> Delete
									</DropdownMenu.Item>
								</DropdownMenu.Group>
							</DropdownMenu.Content>
						</DropdownMenu.Root>
					</div>
				</div>
			</Card>
		{:else}
			<Card class="py-6 px-4 text-sm text-muted-foreground">No status pages yet.</Card>
		{/each}
	</div>
</div>

<AlertDialog.Root bind:open={deleteOpen}>
	<AlertDialog.Content class="max-w-md">
		<AlertDialog.Header>
			<AlertDialog.Title>Delete status page</AlertDialog.Title>
			<AlertDialog.Description>
				This will remove {deleteTarget?.statusPage.title ?? 'this status page'} and its elements.
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel disabled={deleteLoading}>Cancel</AlertDialog.Cancel>
			<AlertDialog.Action onclick={confirmDelete} disabled={deleteLoading}>
				{deleteLoading ? 'Deleting…' : 'Delete'}
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
