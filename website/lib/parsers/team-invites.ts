import type {
  TeamInvite,
  TeamInviteRawData,
  TeamInviteWithTeam,
  TeamInviteWithTeamRawData,
} from "@/lib/schemas/team-invite";

export function parseTeamInvite(raw: TeamInviteRawData): TeamInvite {
  return {
    id: String(raw.id),
    teamID: String(raw.team_id),
    invitedBy: String(raw.invited_by),
    invitedTo: String(raw.invited_to),
    invitedEmail: raw.invited_email,
    role: raw.role,
    status: raw.status,
    expiresAt: raw.expires_at,
    acceptedAt: raw.accepted_at ?? null,
    rejectedAt: raw.rejected_at ?? null,
    canceledAt: raw.canceled_at ?? null,
    createdAt: raw.created_at,
    updatedAt: raw.updated_at,
  };
}

export function parseTeamInvites(rawList: TeamInviteRawData[]): TeamInvite[] {
  return rawList.map(parseTeamInvite);
}

export function parseTeamInviteWithTeam(raw: TeamInviteWithTeamRawData): TeamInviteWithTeam {
  return {
    ...parseTeamInvite(raw),
    teamName: raw.team_name,
  };
}

export function parseTeamInvitesWithTeam(
  rawList: TeamInviteWithTeamRawData[],
): TeamInviteWithTeam[] {
  return rawList.map(parseTeamInviteWithTeam);
}
