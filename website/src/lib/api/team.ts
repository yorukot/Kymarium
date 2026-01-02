import type { MemberRole, Team, TeamInvite, TeamMemberWithUser, TeamWithRole } from '$lib/types';
import { apiRequest } from './utils';

export type TeamsResponse = {
  message: string;
  data: Team[];
};

export function getTeams(): Promise<TeamsResponse> {
	return apiRequest<TeamsResponse>("/teams", { defaultError: 'Failed to fetch teams' });
}

export type TeamResponse = {
	message: string;
	data: TeamWithRole;
};

export function getTeam(teamID: string): Promise<TeamResponse> {
	return apiRequest<TeamResponse>(`/teams/${teamID}`, {
		defaultError: 'Failed to fetch team'
	});
}

export type CreateTeamResponse = {
	message: string;
	data: Team;
};

export function createTeam(name: string): Promise<CreateTeamResponse> {
	return apiRequest<CreateTeamResponse>('/teams', {
		method: 'POST',
		body: { name },
		defaultError: 'Failed to create team'
	});
}

export type UpdateTeamResponse = {
	message: string;
	data: Team;
};

export function updateTeam(teamID: string, name: string): Promise<UpdateTeamResponse> {
	return apiRequest<UpdateTeamResponse>(`/teams/${teamID}`, {
		method: 'PUT',
		body: { name },
		defaultError: 'Failed to update team'
	});
}

export type DeleteTeamResponse = {
	message: string;
};

export function deleteTeam(teamID: string): Promise<DeleteTeamResponse> {
	return apiRequest<DeleteTeamResponse>(`/teams/${teamID}`, {
		method: 'DELETE',
		defaultError: 'Failed to delete team'
	});
}

export type LeaveTeamResponse = {
	message: string;
};

export function leaveTeam(teamID: string): Promise<LeaveTeamResponse> {
	return apiRequest<LeaveTeamResponse>(`/teams/${teamID}/leave`, {
		method: 'POST',
		defaultError: 'Failed to leave team'
	});
}

export type TeamMembersResponse = {
	message: string;
	data: TeamMemberWithUser[];
};

export function getTeamMembers(teamID: string): Promise<TeamMembersResponse> {
	return apiRequest<TeamMembersResponse>(`/teams/${teamID}/members`, {
		defaultError: 'Failed to fetch team members'
	});
}

export type TeamInvitesResponse = {
	message: string;
	data: TeamInvite[];
};

export function getTeamInvites(teamID: string): Promise<TeamInvitesResponse> {
	return apiRequest<TeamInvitesResponse>(`/teams/${teamID}/invites`, {
		defaultError: 'Failed to fetch team invites'
	});
}

export type CreateTeamInviteResponse = {
	message: string;
	data: TeamInvite;
};

export function createTeamInvite(teamID: string, email: string, role?: MemberRole): Promise<CreateTeamInviteResponse> {
	return apiRequest<CreateTeamInviteResponse>(`/teams/${teamID}/invites`, {
		method: 'POST',
		body: { email, role },
		defaultError: 'Failed to send invite'
	});
}

export type CancelTeamInviteResponse = {
	message: string;
	data: TeamInvite;
};

export function cancelTeamInvite(teamID: string, inviteID: string): Promise<CancelTeamInviteResponse> {
	return apiRequest<CancelTeamInviteResponse>(`/teams/${teamID}/invites/${inviteID}`, {
		method: 'DELETE',
		defaultError: 'Failed to cancel invite'
	});
}

export type UpdateTeamInviteResponse = {
	message: string;
	data: TeamInvite;
};

export function updateTeamInvite(
	teamID: string,
	inviteID: string,
	status: 'accepted' | 'rejected'
): Promise<UpdateTeamInviteResponse> {
	return apiRequest<UpdateTeamInviteResponse>(`/teams/${teamID}/invites/${inviteID}`, {
		method: 'PATCH',
		body: { status },
		defaultError: 'Failed to update invite'
	});
}

export type RemoveTeamMemberResponse = {
	message: string;
};

export function removeTeamMember(teamID: string, userID: string): Promise<RemoveTeamMemberResponse> {
	return apiRequest<RemoveTeamMemberResponse>(`/teams/${teamID}/members/${userID}`, {
		method: 'DELETE',
		defaultError: 'Failed to remove team member'
	});
}
