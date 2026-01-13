import { apiRequest } from "@/lib/api/client";
import type { TeamPayload } from "@/lib/schemas/team";

export type TeamResponse = {
  id: string;
  name: string;
  updatedAt: string;
  createdAt: string;
};

export type CreateTeamResponse = {
  message: string;
  data?: TeamResponse;
};

export function createTeam(payload: TeamPayload) {
  return apiRequest<CreateTeamResponse>("/api/teams", {
    method: "POST",
    body: payload,
    defaultError: "Create team failed",
    redirectOn401: true,
  });
}

type MessageResponse = {
  message?: string;
};

export function updateTeam(teamID: string, payload: TeamPayload) {
  return apiRequest<{ message?: string; data?: TeamResponse }>(`/api/teams/${teamID}`, {
    method: "PUT",
    body: payload,
    defaultError: "Update team failed",
    redirectOn401: true,
  });
}

export function deleteTeam(teamID: string) {
  return apiRequest<MessageResponse>(`/api/teams/${teamID}`, {
    method: "DELETE",
    defaultError: "Delete team failed",
    redirectOn401: true,
  });
}

export function leaveTeam(teamID: string) {
  return apiRequest<MessageResponse>(`/api/teams/${teamID}/leave`, {
    method: "POST",
    defaultError: "Leave team failed",
    redirectOn401: true,
  });
}
