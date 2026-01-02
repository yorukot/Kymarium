import type { Account, Session, TeamInviteWithTeam, User } from '$lib/types';
import { apiRequest } from './utils';

export type UserResponse = {
	message: string;
	data: User;
};

export function getUser(): Promise<UserResponse> {
	return apiRequest<UserResponse>("/users/me", { defaultError: 'Failed to fetch user' });
}

export type AccountResponse = {
	message: string;
	data: Account[];
};

export type SessionResponse = {
	message: string;
	data: Session[];
};

export type UpdateUserRequest = Partial<{
	displayName: string;
	avatar: string | null;
}>;

export type UpdateUserResponse = {
	message: string;
	data: User;
};

export type UpdatePasswordRequest = {
	currentPassword: string;
	newPassword: string;
};

export type UpdatePasswordResponse = {
	message: string;
};

export type RevokeSessionsResponse = {
	message: string;
};

export function updateUser(payload: UpdateUserRequest): Promise<UpdateUserResponse> {
	return apiRequest<UpdateUserResponse>("/users/me", {
		method: 'PATCH',
		body: payload,
		defaultError: 'Failed to update user'
	});
}

export function getAccounts(): Promise<AccountResponse> {
	return apiRequest<AccountResponse>("/users/me/account", { defaultError: 'Failed to fetch account info' });
}

export function updatePassword(payload: UpdatePasswordRequest): Promise<UpdatePasswordResponse> {
	return apiRequest<UpdatePasswordResponse>("/users/me/password", {
		method: 'PATCH',
		body: payload,
		defaultError: 'Failed to update password'
	});
}

export function listSessions(): Promise<SessionResponse> {
	return apiRequest<SessionResponse>("/users/me/sessions", { defaultError: 'Failed to fetch sessions' });
}

export type UserInvitesResponse = {
	message: string;
	data: TeamInviteWithTeam[];
};

export function listUserInvites(): Promise<UserInvitesResponse> {
	return apiRequest<UserInvitesResponse>("/users/me/invites", {
		defaultError: 'Failed to fetch invites'
	});
}

export function revokeSession(sessionId: string): Promise<RevokeSessionsResponse> {
	return apiRequest<RevokeSessionsResponse>(`/users/me/sessions/${sessionId}/revoke`, {
		method: 'POST',
		defaultError: 'Failed to revoke session'
	});
}

export function revokeOtherSessions(): Promise<RevokeSessionsResponse> {
	return apiRequest<RevokeSessionsResponse>("/users/me/sessions/revoke-others", {
		method: 'POST',
		defaultError: 'Failed to revoke other sessions'
	});
}
