import { apiRequest } from "@/lib/api/client";

type MessageResponse = {
  message?: string;
};

export function updateTeamMemberRole(
  teamID: string,
  userID: string,
  role: string,
) {
  return apiRequest<MessageResponse>(`/api/teams/${teamID}/members/${userID}`, {
    method: "PATCH",
    body: { role },
    defaultError: "Failed to update member role",
    redirectOn401: true,
  });
}

export function removeTeamMember(teamID: string, userID: string) {
  return apiRequest<MessageResponse>(`/api/teams/${teamID}/members/${userID}`, {
    method: "DELETE",
    defaultError: "Failed to remove member",
    redirectOn401: true,
  });
}
