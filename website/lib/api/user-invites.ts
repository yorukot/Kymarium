import { apiRequest } from "@/lib/api/client";
import { parseTeamInvitesWithTeam } from "@/lib/parsers/team-invites";
import type {
  TeamInviteRawData,
  TeamInviteWithTeam,
  TeamInviteWithTeamRawData,
} from "@/lib/schemas/team-invite";

type ListInvitesResponse = {
  message?: string;
  data?: TeamInviteWithTeamRawData[];
};

export async function listMyPendingInvites(): Promise<TeamInviteWithTeam[]> {
  const res = await apiRequest<ListInvitesResponse>("/api/users/me/invites", {
    defaultError: "Failed to fetch invites",
    redirectOn401: true,
    camelcaseResponse: false,
  });

  return parseTeamInvitesWithTeam(res.data.data ?? []);
}

type UpdateInviteResponse = {
  message?: string;
  data?: TeamInviteRawData;
};

export function respondToTeamInvite(
  teamID: string,
  inviteID: string,
  status: "accepted" | "rejected",
) {
  return apiRequest<UpdateInviteResponse>(`/api/teams/${teamID}/invites/${inviteID}`, {
    method: "PATCH",
    body: { status },
    defaultError: "Failed to update invite",
    redirectOn401: true,
    camelcaseResponse: false,
  });
}
