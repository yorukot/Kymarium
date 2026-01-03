<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { deleteStatusPage } from '$lib/api/status-page';
	import * as AlertDialog from '$lib/components/ui/alert-dialog';
	import { buttonVariants } from '$lib/components/ui/button/index.js';
	import { Spinner } from '$lib/components/ui/spinner';
	import { toast } from 'svelte-sonner';

	let deleteLoading = $state(false);
	async function handleDelete() {
		const teamID = page.params.teamID;
		const statusPageID = page.params.statusPageID;
		if (!teamID || !statusPageID) {
			toast.error('Missing team or status page id');
			return;
		}

		deleteLoading = true;
		try {
			await deleteStatusPage(teamID, statusPageID);
			toast.success('Status page deleted');
			goto(`/${teamID}/status-pages`);
		} catch (error) {
			const message = error instanceof Error ? error.message : 'Failed to delete status page';
			toast.error(message);
		} finally {
			deleteLoading = false;
		}
	}
</script>

<AlertDialog.Root>
	<AlertDialog.Trigger class={buttonVariants({ variant: 'destructive' })}>
		Delete
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>Delete status page</AlertDialog.Title>
			<AlertDialog.Description>
				This will remove the status page and its groups and monitors.
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel>Cancel</AlertDialog.Cancel>
			<AlertDialog.Action onclick={handleDelete} disabled={deleteLoading}>
				{#if deleteLoading}
					<Spinner />
				{/if}
				Delete
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
