<script lang="ts">
	import Icon from '@iconify/svelte';
	import { goto } from '$app/navigation';
	import { toast } from 'svelte-sonner';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { updateTeamInvite } from '$lib/api/team';
	import { listUserInvites } from '$lib/api/user';
	import type { MemberRole, TeamInviteWithTeam } from '$lib/types';

	let { open = $bindable(false) } = $props();

	let invites = $state<TeamInviteWithTeam[]>([]);
	let invitesLoading = $state(false);
	let invitesError = $state('');
	let inviteActionLoading = $state<string | null>(null);

	const roleLabels: Record<MemberRole, string> = {
		owner: 'Owner',
		admin: 'Admin',
		member: 'Member',
		viewer: 'Viewer'
	};

	const roleBadgeVariant: Record<MemberRole, 'default' | 'secondary' | 'outline'> = {
		owner: 'default',
		admin: 'secondary',
		member: 'outline',
		viewer: 'outline'
	};

	const handleOpenChange = (nextOpen: boolean) => {
		open = nextOpen;
		if (!nextOpen) return;
		if (!invitesLoading) {
			void loadInvites();
		}
	};

	async function loadInvites() {
		invitesLoading = true;
		invitesError = '';
		try {
			const response = await listUserInvites();
			invites = response.data ?? [];
		} catch (error) {
			invitesError = error instanceof Error ? error.message : 'Failed to load invites.';
		} finally {
			invitesLoading = false;
		}
	}

	async function handleInviteAction(invite: TeamInviteWithTeam, status: 'accepted' | 'rejected') {
		if (inviteActionLoading) return;
		inviteActionLoading = invite.id;
		try {
			await updateTeamInvite(invite.teamId, invite.id, status);
			invites = invites.filter((item) => item.id !== invite.id);
			if (status === 'accepted') {
				toast.success(`Joined ${invite.teamName}`);
				open = false;
				await goto(`/${invite.teamId}`);
			} else {
				toast.success('Invite rejected');
			}
		} catch (error) {
			const message = error instanceof Error ? error.message : 'Failed to update invite.';
			toast.error(message);
		} finally {
			inviteActionLoading = null;
		}
	}
</script>

<Dialog.Root bind:open onOpenChange={handleOpenChange}>
	<Dialog.Content class="sm:max-w-lg">
		<Dialog.Header>
			<Dialog.Title>Team invitations</Dialog.Title>
			<Dialog.Description>
				Review pending invites. Accepting an invite switches you to that team.
			</Dialog.Description>
		</Dialog.Header>

		<div class="mt-4 space-y-3">
			{#if invitesLoading}
				<p class="text-sm text-muted-foreground">Loading invites…</p>
			{:else if invitesError}
				<p class="text-sm text-destructive">{invitesError}</p>
			{:else if invites.length}
				{#each invites as invite (invite.id)}
					<div class="rounded-lg border bg-card p-4 shadow-sm">
						<div class="flex flex-wrap items-start justify-between gap-3">
							<div class="space-y-1">
								<p class="text-sm font-semibold">{invite.teamName}</p>
								<p class="text-xs text-muted-foreground">
									Invited as {roleLabels[invite.role]} · Expires
									{new Date(invite.expiresAt).toLocaleDateString()}
								</p>
								<p class="text-xs text-muted-foreground">{invite.invitedEmail}</p>
							</div>
							<Badge variant={roleBadgeVariant[invite.role]}>
								{roleLabels[invite.role]}
							</Badge>
						</div>
						<div class="mt-4 flex flex-wrap gap-2">
							<Button
								size="sm"
								disabled={inviteActionLoading === invite.id}
								onclick={() => handleInviteAction(invite, 'accepted')}
							>
								<Icon icon="lucide:check" />
								Accept
							</Button>
							<Button
								size="sm"
								variant="outline"
								disabled={inviteActionLoading === invite.id}
								onclick={() => handleInviteAction(invite, 'rejected')}
							>
								<Icon icon="lucide:x" />
								Reject
							</Button>
						</div>
					</div>
				{/each}
			{:else}
				<div class="rounded-lg border border-dashed px-4 py-8 text-center">
					<p class="text-sm font-medium">No pending invites</p>
					<p class="mt-1 text-xs text-muted-foreground">You're all caught up.</p>
				</div>
			{/if}
		</div>
	</Dialog.Content>
</Dialog.Root>
