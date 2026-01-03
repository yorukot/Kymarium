<script lang="ts">
	import * as AlertDialog from '$lib/components/ui/alert-dialog';
	import { buttonVariants } from '$lib/components/ui/button';
	import { cn } from '$lib/utils';
	import type { Monitor } from '$lib/types';

	let {
		open = $bindable(false),
		monitor,
		onConfirm,
		loading = false
	}: {
		open?: boolean;
		monitor: Monitor | null;
		onConfirm: () => void;
		loading?: boolean;
	} = $props();
</script>

<AlertDialog.Root bind:open={open}>
	<AlertDialog.Portal>
		<AlertDialog.Overlay />
		<AlertDialog.Content>
			<AlertDialog.Header>
				<AlertDialog.Title>Delete monitor</AlertDialog.Title>
				<AlertDialog.Description>
					Are you sure you want to delete
					<strong>{monitor?.name}</strong>? This action cannot be undone.
				</AlertDialog.Description>
			</AlertDialog.Header>
			<AlertDialog.Footer>
				<button
					class={cn(buttonVariants({ variant: 'outline' }))}
					disabled={loading}
					onclick={() => (open = false)}
				>
					Cancel
				</button>
				<button
					class={cn(buttonVariants({ variant: 'destructive' }))}
					onclick={onConfirm}
					disabled={loading}
				>
					{loading ? 'Deleting…' : 'Delete'}
				</button>
			</AlertDialog.Footer>
		</AlertDialog.Content>
	</AlertDialog.Portal>
</AlertDialog.Root>
