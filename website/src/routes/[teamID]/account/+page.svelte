<script lang="ts">
	import Icon from '@iconify/svelte';
	import * as Avatar from '$lib/components/ui/avatar/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { toast } from 'svelte-sonner';
	import {
		listSessions,
		revokeOtherSessions,
		revokeSession,
		updatePassword,
		updateUser
	} from '$lib/api/user';
	import { isAuthExpiredError } from '$lib/api/utils';
	import type { Account, Session, User } from '$lib/types';

	/** @type {import('./$types').PageProps} */
	let { data } = $props();

	let user = $derived<User>(data.user);
	let accounts = $derived<Account[]>(data.accounts ?? []);
	let sessions = $derived<Session[]>(data.sessions ?? []);

	let displayName = $derived(data.user.displayName ?? '');
	let avatar = $derived(data.user.avatar ?? '');

	let profileError = $state('');
	let profileSuccess = $state('');
	let profileLoading = $state(false);

	let currentPassword = $state('');
	let newPassword = $state('');
	let confirmPassword = $state('');
	let passwordError = $state('');
	let passwordSuccess = $state('');
	let passwordLoading = $state(false);

	let sessionsLoading = $state(false);
	let revokeOtherLoading = $state(false);
	let revokeSessionId = $state<string | null>(null);

	const isProfileDirty = $derived(
		displayName.trim() !== (user.displayName ?? '') || (avatar ?? '') !== (user.avatar ?? '')
	);

	const hasAccounts = $derived(accounts.length > 0);
	const hasSessions = $derived(sessions.length > 0);

	function formatTimestamp(value: string) {
		const date = new Date(value);
		if (Number.isNaN(date.getTime())) return value;
		return date.toLocaleString();
	}

	function normalizeAvatar(value: string) {
		const trimmed = value.trim();
		return trimmed === '' ? null : trimmed;
	}

	async function handleProfileSave(event: SubmitEvent) {
		event.preventDefault();
		profileError = '';
		profileSuccess = '';

		if (!displayName.trim()) {
			profileError = 'Display name is required.';
			return;
		}

		if (!isProfileDirty) return;

		profileLoading = true;
		try {
			const res = await updateUser({
				displayName: displayName.trim(),
				avatar: normalizeAvatar(avatar ?? '')
			});

			user = res.data;
			displayName = res.data.displayName ?? '';
			avatar = res.data.avatar ?? '';
			profileSuccess = 'Profile updated successfully.';
			toast.success('Profile updated');
		} catch (err) {
			profileError = err instanceof Error ? err.message : 'Failed to update profile.';
			toast.error(profileError);
		} finally {
			profileLoading = false;
		}
	}

	async function handlePasswordSave(event: SubmitEvent) {
		event.preventDefault();
		passwordError = '';
		passwordSuccess = '';

		if (!currentPassword.trim() || !newPassword.trim()) {
			passwordError = 'Both current and new password are required.';
			return;
		}

		if (newPassword !== confirmPassword) {
			passwordError = 'Passwords do not match.';
			return;
		}

		passwordLoading = true;
		try {
			await updatePassword({ currentPassword, newPassword });
			currentPassword = '';
			newPassword = '';
			confirmPassword = '';
			passwordSuccess = 'Password updated successfully.';
			toast.success('Password updated');
		} catch (err) {
			passwordError = err instanceof Error ? err.message : 'Failed to update password.';
			toast.error(passwordError);
		} finally {
			passwordLoading = false;
		}
	}

	async function refreshSessions() {
		sessionsLoading = true;
		try {
			const res = await listSessions();
			sessions = res.data;
		} catch (err) {
			if (isAuthExpiredError(err)) return;
			const message = err instanceof Error ? err.message : 'Failed to refresh sessions.';
			toast.error(message);
		} finally {
			sessionsLoading = false;
		}
	}

	async function handleRevokeSession(sessionId: string) {
		revokeSessionId = sessionId;
		try {
			await revokeSession(sessionId);
			toast.success('Session revoked');
			await refreshSessions();
		} catch (err) {
			if (isAuthExpiredError(err)) return;
			const message = err instanceof Error ? err.message : 'Failed to revoke session.';
			toast.error(message);
		} finally {
			revokeSessionId = null;
		}
	}

	async function handleRevokeOtherSessions() {
		revokeOtherLoading = true;
		try {
			await revokeOtherSessions();
			toast.success('Other sessions revoked');
			await refreshSessions();
		} catch (err) {
			if (isAuthExpiredError(err)) return;
			const message = err instanceof Error ? err.message : 'Failed to revoke other sessions.';
			toast.error(message);
		} finally {
			revokeOtherLoading = false;
		}
	}
</script>

<header class="flex flex-col gap-2">
	<h1 class="text-2xl font-semibold">Account</h1>
	<p class="text-sm text-muted-foreground">Manage your profile, account info, and sessions.</p>
</header>

<div class="grid gap-6 mt-6">
	<Card.Root>
		<Card.Header>
			<Card.Title>Profile</Card.Title>
			<Card.Description>Update your name and avatar for the workspace.</Card.Description>
		</Card.Header>
		<Card.Content>
			<form class="grid gap-4" onsubmit={handleProfileSave}>
				<div class="flex flex-col gap-4 md:flex-row md:items-center md:gap-6">
					<Avatar.Root class="size-16 rounded-lg">
						<Avatar.Image src={avatar} alt={displayName} />
						<Avatar.Fallback class="rounded-lg">{displayName?.[0] ?? 'U'}</Avatar.Fallback>
					</Avatar.Root>
					<div class="grid gap-2 flex-1">
						<label class="text-sm font-medium" for="display-name">Display name</label>
						<Input
							id="display-name"
							placeholder="Your name"
							bind:value={displayName}
							disabled={profileLoading}
						/>
					</div>
				</div>

				<div class="grid gap-2">
					<label class="text-sm font-medium" for="avatar-url">Avatar URL</label>
					<Input
						id="avatar-url"
						placeholder="https://example.com/avatar.png"
						bind:value={avatar}
						disabled={profileLoading}
					/>
					<p class="text-xs text-muted-foreground">
						Leave blank to use the default avatar.
					</p>
				</div>

				<div class="grid gap-2 md:grid-cols-2">
					<div class="grid gap-2">
						<label class="text-sm font-medium" for="user-id">User ID</label>
						<Input id="user-id" value={user.id} readonly disabled />
					</div>
					<div class="grid gap-2">
						<span class="text-sm font-medium">Verification</span>
						<div class="flex items-center gap-2">
							<Badge variant={user.verified ? 'default' : 'outline'}>
								{user.verified ? 'Verified' : 'Unverified'}
							</Badge>
							<span class="text-xs text-muted-foreground">
								Email verification status.
							</span>
						</div>
					</div>
				</div>

				{#if profileError}
					<p class="text-sm text-destructive">{profileError}</p>
				{/if}
				{#if profileSuccess}
					<p class="text-sm text-success">{profileSuccess}</p>
				{/if}

				<div class="flex items-center gap-2">
					<Button type="submit" disabled={profileLoading || !isProfileDirty}>
						{#if profileLoading}
							<Icon icon="lucide:loader-2" class="animate-spin" />
						{/if}
						Save changes
					</Button>
				</div>
			</form>
		</Card.Content>
	</Card.Root>

	<Card.Root>
		<Card.Header>
			<Card.Title>Account info</Card.Title>
			<Card.Description>Login providers connected to your account.</Card.Description>
		</Card.Header>
		<Card.Content class="grid gap-4">
			{#if hasAccounts}
				<div class="grid gap-3">
					{#each accounts as account (account.id)}
						<div class="flex flex-col gap-2 rounded-lg border border-border/60 px-4 py-3 md:flex-row md:items-center md:justify-between">
							<div class="grid gap-1">
								<span class="text-sm font-medium">{account.email}</span>
								<span class="text-xs text-muted-foreground">
									Provider: {account.provider}
								</span>
							</div>
							<div class="flex items-center gap-2">
								{#if account.isPrimary}
									<Badge variant="secondary">Primary</Badge>
								{/if}
								<span class="text-xs text-muted-foreground">
									Added {formatTimestamp(account.createdAt)}
								</span>
							</div>
						</div>
					{/each}
				</div>
			{:else}
				<p class="text-sm text-muted-foreground">No account details available.</p>
			{/if}
		</Card.Content>
	</Card.Root>

	<Card.Root>
		<Card.Header>
			<Card.Title>Password</Card.Title>
			<Card.Description>Update your password for email sign-in.</Card.Description>
		</Card.Header>
		<Card.Content>
			<form class="grid gap-4" onsubmit={handlePasswordSave}>
				<div class="grid gap-2">
					<label class="text-sm font-medium" for="current-password">Current password</label>
					<Input
						id="current-password"
						type="password"
						placeholder="Enter current password"
						bind:value={currentPassword}
						disabled={passwordLoading}
					/>
				</div>
				<div class="grid gap-2">
					<label class="text-sm font-medium" for="new-password">New password</label>
					<Input
						id="new-password"
						type="password"
						placeholder="Enter new password"
						bind:value={newPassword}
						disabled={passwordLoading}
					/>
				</div>
				<div class="grid gap-2">
					<label class="text-sm font-medium" for="confirm-password">Confirm new password</label>
					<Input
						id="confirm-password"
						type="password"
						placeholder="Confirm new password"
						bind:value={confirmPassword}
						disabled={passwordLoading}
					/>
				</div>

				{#if passwordError}
					<p class="text-sm text-destructive">{passwordError}</p>
				{/if}
				{#if passwordSuccess}
					<p class="text-sm text-success">{passwordSuccess}</p>
				{/if}

				<div class="flex items-center gap-2">
					<Button type="submit" disabled={passwordLoading}>
						{#if passwordLoading}
							<Icon icon="lucide:loader-2" class="animate-spin" />
						{/if}
						Update password
					</Button>
				</div>
			</form>
		</Card.Content>
	</Card.Root>

	<Card.Root>
		<Card.Header>
			<Card.Title>Active sessions</Card.Title>
			<Card.Description>Manage where you are currently signed in.</Card.Description>
		</Card.Header>
		<Card.Content class="grid gap-4">
			<div class="flex flex-wrap items-center justify-between gap-3">
				<div class="text-sm text-muted-foreground">
					{sessionsLoading ? 'Refreshing sessions…' : 'Active sessions are updated in real time.'}
				</div>
				<Button
					variant="outline"
					onclick={refreshSessions}
					disabled={sessionsLoading}
				>
					Refresh
				</Button>
			</div>

			{#if hasSessions}
				<div class="grid gap-3">
					{#each sessions as session (session.id)}
						<div class="flex flex-col gap-3 rounded-lg border border-border/60 px-4 py-3 md:flex-row md:items-center md:justify-between">
							<div class="grid gap-1">
								<span class="text-sm font-medium">
									{session.userAgent ?? 'Unknown device'}
								</span>
								<span class="text-xs text-muted-foreground">
									{session.ip ? `IP ${session.ip}` : 'IP unavailable'} ·
									{formatTimestamp(session.createdAt)}
								</span>
							</div>
							<div class="flex items-center gap-2">
								{#if session.current}
									<Badge>Current</Badge>
								{/if}
								<Button
									variant="outline"
									size="sm"
									onclick={() => handleRevokeSession(session.id)}
									disabled={session.current || revokeSessionId === session.id}
								>
									{#if revokeSessionId === session.id}
										<Icon icon="lucide:loader-2" class="animate-spin" />
									{/if}
									Revoke
								</Button>
							</div>
						</div>
					{/each}
				</div>
				<div class="flex items-center justify-between gap-2">
					<p class="text-xs text-muted-foreground">
						Revoking a session logs it out on the next refresh.
					</p>
					<Button
						variant="destructive"
						onclick={handleRevokeOtherSessions}
						disabled={revokeOtherLoading || sessions.length < 2}
					>
						{#if revokeOtherLoading}
							<Icon icon="lucide:loader-2" class="animate-spin" />
						{/if}
						Revoke other sessions
					</Button>
				</div>
			{:else}
				<p class="text-sm text-muted-foreground">No active sessions found.</p>
			{/if}
		</Card.Content>
	</Card.Root>
</div>
