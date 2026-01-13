import { apiRequest } from "@/lib/api/client";
import type { CreateTeamInviteValues } from "@/lib/schemas/team-invite";

type MessageResponse = {
  message?: string;
};

export function createTeamInvite(teamID: string, payload: CreateTeamInviteValues) {
  return apiRequest<MessageResponse>(`/api/teams/${teamID}/invites`, {
    method: "POST",
    body: payload,
    defaultError: "Failed to send invite",
    redirectOn401: true,
  });
}

export function cancelTeamInvite(teamID: string, inviteID: string) {
  return apiRequest<MessageResponse>(`/api/teams/${teamID}/invites/${inviteID}`, {
    method: "DELETE",
    defaultError: "Failed to cancel invite",
    redirectOn401: true,
  });
}

