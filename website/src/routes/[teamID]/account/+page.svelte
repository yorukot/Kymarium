<script lang="ts">
	import { createForm } from 'felte';
	import { validator } from '@felte/validator-zod';
	import { reporter, ValidationMessage } from '@felte/reporter-svelte';
	import { z } from 'zod';
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

	let passwordSuccess = $state('');

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

	const passwordSchema = z
		.object({
			currentPassword: z
				.string()
				.min(8, 'Current password must be at least 8 characters.')
				.max(255, 'Current password is too long.'),
			newPassword: z
				.string()
				.min(8, 'New password must be at least 8 characters.')
				.max(255, 'New password is too long.'),
			confirmPassword: z
				.string()
				.min(8, 'Confirm password must be at least 8 characters.')
				.max(255, 'Confirm password is too long.')
		})
		.superRefine((values, ctx) => {
			if (values.currentPassword === values.newPassword) {
				ctx.addIssue({
					code: z.ZodIssueCode.custom,
					path: ['newPassword'],
					message: 'New password must be different from the current password.'
				});
			}
			if (values.newPassword !== values.confirmPassword) {
				ctx.addIssue({
					code: z.ZodIssueCode.custom,
					path: ['confirmPassword'],
					message: 'Passwords do not match.'
				});
			}
		});

	type PasswordFormValues = z.infer<typeof passwordSchema>;

	const {
		form: passwordForm,
		isSubmitting: isPasswordSubmitting,
		reset: resetPasswordForm
	} = createForm<PasswordFormValues>({
		extend: [validator({ schema: passwordSchema }), reporter()],
		onSubmit: async (values) => {
			passwordSuccess = '';
			try {
				await updatePassword({
					currentPassword: values.currentPassword,
					newPassword: values.newPassword
				});
				resetPasswordForm();
				passwordSuccess = 'Password updated successfully.';
				toast.success('Password updated');
			} catch (err) {
				const message = err instanceof Error ? err.message : 'Failed to update password.';
				toast.error(message);
				return { FORM_ERROR: message };
			}
		}
	});

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
					<p class="text-xs text-muted-foreground">Leave blank to use the default avatar.</p>
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
							<span class="text-xs text-muted-foreground"> Email verification status. </span>
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
						<div
							class="flex flex-col gap-2 rounded-lg border border-border/60 px-4 py-3 md:flex-row md:items-center md:justify-between"
						>
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
			<form class="grid gap-4" use:passwordForm>
				<div class="grid gap-2">
					<label class="text-sm font-medium" for="current-password">Current password</label>
					<Input
						id="current-password"
						name="currentPassword"
						type="password"
						placeholder="Enter current password"
						autocomplete="current-password"
						disabled={$isPasswordSubmitting}
					/>
					<ValidationMessage for="currentPassword" let:messages>
						{#if messages?.length}
							<p class="text-sm text-destructive">{messages[0]}</p>
						{/if}
					</ValidationMessage>
				</div>
				<div class="grid gap-2">
					<label class="text-sm font-medium" for="new-password">New password</label>
					<Input
						id="new-password"
						name="newPassword"
						type="password"
						placeholder="Enter new password"
						autocomplete="new-password"
						disabled={$isPasswordSubmitting}
					/>
					<ValidationMessage for="newPassword" let:messages>
						{#if messages?.length}
							<p class="text-sm text-destructive">{messages[0]}</p>
						{/if}
					</ValidationMessage>
				</div>
				<div class="grid gap-2">
					<label class="text-sm font-medium" for="confirm-password">Confirm new password</label>
					<Input
						id="confirm-password"
						name="confirmPassword"
						type="password"
						placeholder="Confirm new password"
						autocomplete="new-password"
						disabled={$isPasswordSubmitting}
					/>
					<ValidationMessage for="confirmPassword" let:messages>
						{#if messages?.length}
							<p class="text-sm text-destructive">{messages[0]}</p>
						{/if}
					</ValidationMessage>
				</div>

				<ValidationMessage for="FORM_ERROR" let:messages>
					{#if messages?.length}
						<p class="text-sm text-destructive">{messages[0]}</p>
					{/if}
				</ValidationMessage>
				{#if passwordSuccess}
					<p class="text-sm text-success">{passwordSuccess}</p>
				{/if}

				<div class="flex items-center gap-2">
					<Button type="submit" disabled={$isPasswordSubmitting}>
						{#if $isPasswordSubmitting}
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
			<div class="flex items-center justify-between gap-2">
				<Card.Title>Active sessions</Card.Title>
				<Button size="icon-sm" onclick={refreshSessions} disabled={sessionsLoading}>
					<Icon icon="lucide:refresh-ccw"/>
				</Button>
			</div>
			<Card.Description>Manage where you are currently signed in.</Card.Description>
		</Card.Header>
		<Card.Content class="grid gap-4">
			{#if hasSessions}
				<div class="grid gap-3">
					{#each sessions as session (session.id)}
						<div
							class="flex flex-col gap-3 rounded-lg border border-border/60 px-4 py-3 md:flex-row md:items-center md:justify-between"
						>
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
								<Button
									variant="outline"
									size="sm"
									onclick={() => handleRevokeSession(session.id)}
									disabled={revokeSessionId === session.id}
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
				<div class="flex items-center justify-end gap-2">
					<Button
						variant="destructive"
						onclick={handleRevokeOtherSessions}
						disabled={revokeOtherLoading || sessions.length < 2}
					>
						{#if revokeOtherLoading}
							<Icon icon="lucide:loader-2" class="animate-spin" />
						{/if}
						Revoke all sessions
					</Button>
				</div>
			{:else}
				<p class="text-sm text-muted-foreground">No active sessions found.</p>
			{/if}
		</Card.Content>
	</Card.Root>
</div>
