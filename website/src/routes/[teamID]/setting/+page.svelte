<script lang="ts">
	import Icon from '@iconify/svelte';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import { deleteTeam, leaveTeam, updateTeam } from '$lib/api/team';
	import type { TeamWithRole } from '../../../types';

	/** @type {import('./$types').PageProps} */
	let { data } = $props();

	const teamID = $derived(page.params.teamID ?? '');
	let team = $state<TeamWithRole>(data.team);
	let name = $state(data.team?.name ?? '');
	let saveError = $state('');
	let saveSuccess = $state('');
	let saveLoading = $state(false);
	let deleteError = $state('');
	let deleteLoading = $state(false);
	let deleteOpen = $state(false);
		let leaveError = $state('');
		let leaveLoading = $state(false);
		let leaveOpen = $state(false);

	const canEditTeam = $derived(team?.role === 'owner' || team?.role === 'admin');
	const canDeleteTeam = $derived(team?.role === 'owner');
	const isNameDirty = $derived(name.trim() !== (team?.name ?? ''));

	const roleLabel: Record<TeamWithRole['role'], string> = {
		owner: 'Owner',
		admin: 'Admin',
		member: 'Member',
		viewer: 'Viewer'
	};

	const roleBadgeVariant: Record<TeamWithRole['role'], 'default' | 'secondary' | 'outline'> = {
		owner: 'default',
		admin: 'secondary',
		member: 'outline',
		viewer: 'outline'
	};

	async function handleSave(event: SubmitEvent) {
		event.preventDefault();
		saveError = '';
		saveSuccess = '';

		if (!canEditTeam) {
			saveError = 'You do not have permission to update this team.';
			return;
		}

		if (!name.trim()) {
			saveError = 'Team name is required.';
			return;
		}

		if (!isNameDirty) {
			return;
		}

		saveLoading = true;
		try {
			const response = await updateTeam(teamID, name.trim());
			team = { ...team, ...response.data };
			name = response.data.name;
			saveSuccess = 'Team updated successfully.';
		} catch (error) {
			saveError = error instanceof Error ? error.message : 'Failed to update team.';
		} finally {
			saveLoading = false;
		}
	}

	async function handleDelete() {
		deleteError = '';

		if (!canDeleteTeam) {
			deleteError = 'Only team owners can delete this team.';
			return;
		}

		deleteLoading = true;
		try {
			await deleteTeam(teamID);
			deleteOpen = false;
			await goto('/');
		} catch (error) {
			deleteError = error instanceof Error ? error.message : 'Failed to delete team.';
		} finally {
			deleteLoading = false;
		}
	}

		async function handleLeave() {
			leaveError = '';

			if (canDeleteTeam) {
				leaveError = 'Owners must transfer ownership before leaving.';
				return;
			}

		leaveLoading = true;
		try {
			await leaveTeam(teamID);
			leaveOpen = false;
			await goto('/');
		} catch (error) {
			leaveError = error instanceof Error ? error.message : 'Failed to leave team.';
		} finally {
			leaveLoading = false;
		}
	}
</script>

<header class="flex flex-col gap-2">
	<h1 class="text-2xl font-semibold">Team settings</h1>
	<p class="text-sm text-muted-foreground">
		Update your team profile details and manage ownership settings.
	</p>
</header>

<div class="grid gap-6 mt-6">
	<Card.Root>
		<Card.Header>
			<Card.Title>Profile</Card.Title>
			<Card.Description>Basic details about your team.</Card.Description>
		</Card.Header>
		<Card.Content>
			<form class="grid gap-4" onsubmit={handleSave}>
				<div class="grid gap-2">
					<label class="text-sm font-medium" for="team-name">Team name</label>
					<Input
						id="team-name"
						placeholder="Team name"
						bind:value={name}
						disabled={!canEditTeam || saveLoading}
					/>
					<p class="text-xs text-muted-foreground">
						This name appears in the team switcher and notifications.
					</p>
				</div>

				<div class="grid gap-2 md:grid-cols-2">
					<div class="grid gap-2">
						<label class="text-sm font-medium" for="team-id">Team ID</label>
						<Input id="team-id" value={team?.id ?? ''} readonly disabled />
						<p class="text-xs text-muted-foreground">Used in API requests and integrations.</p>
					</div>
					<div class="grid gap-2">
						<span class="text-sm font-medium">Your role</span>
						<div class="flex items-center gap-2">
							<Badge variant={roleBadgeVariant[team.role]}>{roleLabel[team.role]}</Badge>
							<span class="text-xs text-muted-foreground">Permissions are based on role.</span>
						</div>
					</div>
				</div>

				{#if saveError}
					<p class="text-sm text-destructive">{saveError}</p>
				{/if}
				{#if saveSuccess}
					<p class="text-sm text-success">{saveSuccess}</p>
				{/if}

				<div class="flex flex-wrap items-center gap-2">
					<Button type="submit" disabled={!canEditTeam || saveLoading || !isNameDirty}>
						{#if saveLoading}
							<Icon icon="lucide:loader-2" class="animate-spin" />
						{/if}
						Save changes
					</Button>
					{#if !canEditTeam}
						<span class="text-xs text-muted-foreground">
							Only owners and admins can edit team details.
						</span>
					{/if}
				</div>
			</form>
		</Card.Content>
	</Card.Root>

	<Card.Root>
		<Card.Header>
			<Card.Title>Danger zone</Card.Title>
			<Card.Description>Deleting a team removes all monitors, incidents, and data.</Card.Description>
		</Card.Header>
		<Card.Content class="grid gap-3">
			{#if deleteError}
				<p class="text-sm text-destructive">{deleteError}</p>
			{/if}
			{#if leaveError}
				<p class="text-sm text-destructive">{leaveError}</p>
			{/if}
			<AlertDialog.Root bind:open={deleteOpen}>
				<AlertDialog.Trigger class="w-fit">
					<Button variant="destructive" disabled={!canDeleteTeam || deleteLoading}>
						Delete team
					</Button>
				</AlertDialog.Trigger>
				<AlertDialog.Portal>
					<AlertDialog.Overlay />
					<AlertDialog.Content class="max-w-md">
						<AlertDialog.Header>
							<AlertDialog.Title>Delete this team?</AlertDialog.Title>
							<AlertDialog.Description>
								This action is permanent. All team data will be removed.
							</AlertDialog.Description>
						</AlertDialog.Header>
						<AlertDialog.Footer>
							<AlertDialog.Cancel disabled={deleteLoading}>Cancel</AlertDialog.Cancel>
							<AlertDialog.Action onclick={handleDelete} disabled={deleteLoading}>
								{#if deleteLoading}
									<Icon icon="lucide:loader-2" class="animate-spin" />
								{/if}
								Confirm delete
							</AlertDialog.Action>
						</AlertDialog.Footer>
					</AlertDialog.Content>
				</AlertDialog.Portal>
			</AlertDialog.Root>
			{#if !canDeleteTeam}
				<AlertDialog.Root bind:open={leaveOpen}>
					<AlertDialog.Trigger class="w-fit">
						<Button variant="destructive" disabled={leaveLoading}>
							Leave team
						</Button>
					</AlertDialog.Trigger>
					<AlertDialog.Portal>
						<AlertDialog.Overlay />
						<AlertDialog.Content class="max-w-md">
							<AlertDialog.Header>
								<AlertDialog.Title>Leave this team?</AlertDialog.Title>
								<AlertDialog.Description>
									You will lose access to this team and its monitors.
								</AlertDialog.Description>
							</AlertDialog.Header>
							<AlertDialog.Footer>
								<AlertDialog.Cancel disabled={leaveLoading}>Cancel</AlertDialog.Cancel>
								<AlertDialog.Action onclick={handleLeave} disabled={leaveLoading}>
									{#if leaveLoading}
										<Icon icon="lucide:loader-2" class="animate-spin" />
									{/if}
									Confirm leave
								</AlertDialog.Action>
							</AlertDialog.Footer>
						</AlertDialog.Content>
					</AlertDialog.Portal>
				</AlertDialog.Root>
			{/if}
			{#if !canDeleteTeam}
				<p class="text-xs text-muted-foreground">
					Only team owners can delete a team.
				</p>
			{/if}
		</Card.Content>
	</Card.Root>
</div>
