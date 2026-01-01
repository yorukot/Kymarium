<script lang="ts">
	import Icon from '@iconify/svelte';
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import * as Table from '$lib/components/ui/table/index.js';
	import * as Avatar from '$lib/components/ui/avatar/index.js';
	import { NativeSelect, NativeSelectOption } from '$lib/components/ui/native-select/index.js';
	import * as Empty from '$lib/components/ui/empty/index.js';
	import {
		cancelTeamInvite,
		createTeamInvite,
		getTeamInvites,
		removeTeamMember
	} from '$lib/api/team';
	import type { MemberRole, TeamInvite, TeamMemberWithUser } from '$lib/types';

	/** @type {import('./$types').PageProps} */
	let { data } = $props();

	const teamID = $derived(page.params.teamID ?? '');
	let members = $derived<TeamMemberWithUser[]>(data.members ?? []);
	let invites = $derived<TeamInvite[]>(data.invites ?? []);

	let inviteEmail = $state('');
	let inviteRole = $state<MemberRole>('member');
	let inviteError = $state('');
	let inviteSubmitting = $state(false);

	let invitesLoaded = $state(false);
	let invitesLoading = $state(false);
	let invitesError = $state('');

	let memberActionError = $state('');
	let memberActionLoading = $state<string | null>(null);

	const currentMember = $derived(members.find((member) => member.userId === data.user?.id));

	const canManageInvites = $derived(
		currentMember?.role === 'owner' || currentMember?.role === 'admin'
	);

	const pendingInvites = $derived((invites ?? []).filter((invite) => invite.status === 'pending'));

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

	onMount(() => {
		if (!canManageInvites || invitesLoaded || invitesLoading) return;
		void loadInvites();
	});

	async function loadInvites() {
		invitesLoading = true;
		invitesError = '';
		try {
			const response = await getTeamInvites(teamID);
			invites = response.data ?? [];
			invitesLoaded = true;
		} catch (error) {
			invitesError = error instanceof Error ? error.message : 'Failed to load invites.';
		} finally {
			invitesLoading = false;
		}
	}

	async function handleInviteSubmit(event: SubmitEvent) {
		event.preventDefault();
		inviteError = '';
		memberActionError = '';

		if (!inviteEmail.trim()) {
			inviteError = 'Email is required.';
			return;
		}

		inviteSubmitting = true;
		try {
			const response = await createTeamInvite(teamID, inviteEmail.trim(), inviteRole);
			invites = [response.data, ...invites];
			inviteEmail = '';
		} catch (error) {
			inviteError = error instanceof Error ? error.message : 'Failed to send invite.';
		} finally {
			inviteSubmitting = false;
		}
	}

	async function handleCancelInvite(invite: TeamInvite) {
		memberActionError = '';
		memberActionLoading = invite.id;

		try {
			const response = await cancelTeamInvite(teamID, invite.id);
			invites = invites.map((item) => (item.id === invite.id ? response.data : item));
		} catch (error) {
			memberActionError = error instanceof Error ? error.message : 'Failed to cancel invite.';
		} finally {
			memberActionLoading = null;
		}
	}

	async function handleRemoveMember(member: TeamMemberWithUser) {
		memberActionError = '';
		memberActionLoading = member.userId;

		if (!confirm(`Remove ${member.displayName} from this team?`)) {
			memberActionLoading = null;
			return;
		}

		try {
			await removeTeamMember(teamID, member.userId);
			members = members.filter((item) => item.userId !== member.userId);
		} catch (error) {
			memberActionError = error instanceof Error ? error.message : 'Failed to remove member.';
		} finally {
			memberActionLoading = null;
		}
	}
</script>

<header class="flex flex-col gap-2">
	<h1 class="text-2xl font-semibold">Members</h1>
	<p class="text-sm text-muted-foreground">
		Manage who can access this team and send invites to new teammates.
	</p>
</header>

<div class="grid gap-6 mt-6">
  <Card.Root class="overflow-hidden">
		<Card.Header>
			<Card.Title>Team members</Card.Title>
			<Card.Description>Everyone with access to this team.</Card.Description>
		</Card.Header>
		<Card.Content>
			{#if memberActionError}
				<p class="text-sm text-destructive mb-3">{memberActionError}</p>
			{/if}
			{#if members.length}
					<Table.Root>
						<Table.Header>
							<Table.Row>
								<Table.Head>Member</Table.Head>
								<Table.Head>Email</Table.Head>
								<Table.Head>Role</Table.Head>
								<Table.Head class="text-right">Actions</Table.Head>
							</Table.Row>
						</Table.Header>
						<Table.Body>
							{#each members as member (member.id)}
								<Table.Row>
									<Table.Cell>
										<div class="flex items-center gap-3">
											<Avatar.Root class="size-9">
												{#if member.avatar}
													<Avatar.Image src={member.avatar} alt={member.displayName} />
												{/if}
												<Avatar.Fallback>
													{member.displayName.slice(0, 1).toUpperCase()}
												</Avatar.Fallback>
											</Avatar.Root>
											<div>
												<p class="font-medium">{member.displayName}</p>
												<p class="text-xs text-muted-foreground">ID {member.userId}</p>
											</div>
										</div>
									</Table.Cell>
									<Table.Cell>{member.email}</Table.Cell>
									<Table.Cell>
										<Badge variant={roleBadgeVariant[member.role]}>{roleLabels[member.role]}</Badge>
									</Table.Cell>
									<Table.Cell class="text-right">
										{#if canManageInvites && member.userId !== data.user?.id}
											<Button
												variant="destructive"
												size="sm"
												disabled={memberActionLoading === member.userId}
												onclick={() => handleRemoveMember(member)}
											>
												<Icon icon="lucide:user-x" />
												Remove
											</Button>
										{:else}
											<span class="text-xs text-muted-foreground">—</span>
										{/if}
									</Table.Cell>
								</Table.Row>
							{/each}
						</Table.Body>
					</Table.Root>
			{:else}
				<Empty.Root>
					<Empty.Title>No members yet</Empty.Title>
					<Empty.Description>Invite someone to join this team.</Empty.Description>
				</Empty.Root>
			{/if}
		</Card.Content>
	</Card.Root>

	{#if canManageInvites}
		<Card.Root>
			<Card.Header>
				<Card.Title>Invite teammate</Card.Title>
				<Card.Description>Send an invite via email. Invites expire in 7 days.</Card.Description>
			</Card.Header>
			<Card.Content>
				<form class="grid gap-4 md:grid-cols-[1fr_auto_auto]" onsubmit={handleInviteSubmit}>
					<Input
						type="email"
						placeholder="teammate@company.com"
						bind:value={inviteEmail}
						required
					/>
					<NativeSelect bind:value={inviteRole}>
						<NativeSelectOption value="member">Member</NativeSelectOption>
						<NativeSelectOption value="admin">Admin</NativeSelectOption>
						<NativeSelectOption value="viewer">Viewer</NativeSelectOption>
					</NativeSelect>
					<Button type="submit" disabled={inviteSubmitting}>
						<Icon icon="lucide:send" />
						Invite
					</Button>
				</form>
				{#if inviteError}
					<p class="text-sm text-destructive mt-3">{inviteError}</p>
				{/if}
			</Card.Content>
		</Card.Root>

		<Card.Root>
			<Card.Header>
				<Card.Title>Pending invites</Card.Title>
				<Card.Description>Keep track of outstanding invites.</Card.Description>
			</Card.Header>
			<Card.Content>
				{#if invitesLoading}
					<p class="text-sm text-muted-foreground">Loading invites…</p>
				{:else if invitesError}
					<p class="text-sm text-destructive">{invitesError}</p>
				{:else if pendingInvites.length}
					<Table.Root>
						<Table.Header>
							<Table.Row>
								<Table.Head>Email</Table.Head>
								<Table.Head>Role</Table.Head>
								<Table.Head>Expires</Table.Head>
								<Table.Head class="text-right">Action</Table.Head>
							</Table.Row>
						</Table.Header>
						<Table.Body>
							{#each pendingInvites as invite (invite.id)}
								<Table.Row>
									<Table.Cell>{invite.invitedEmail}</Table.Cell>
									<Table.Cell>
										<Badge variant={roleBadgeVariant[invite.role]}>
											{roleLabels[invite.role]}
										</Badge>
									</Table.Cell>
									<Table.Cell>
										{new Date(invite.expiresAt).toLocaleDateString()}
									</Table.Cell>
									<Table.Cell class="text-right">
										<Button
											variant="destructive"
											size="sm"
											disabled={memberActionLoading === invite.id}
											onclick={() => handleCancelInvite(invite)}
										>
											<Icon icon="lucide:x" />
											Cancel
										</Button>
									</Table.Cell>
								</Table.Row>
							{/each}
						</Table.Body>
					</Table.Root>
				{:else}
					<Empty.Root>
						<Empty.Title>No pending invites</Empty.Title>
						<Empty.Description>Invites you send will show up here.</Empty.Description>
					</Empty.Root>
				{/if}
			</Card.Content>
		</Card.Root>
	{/if}
</div>
