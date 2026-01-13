import type { TeamMember, TeamMemberRawData } from "@/lib/schemas/team-member";

export function parseTeamMember(raw: TeamMemberRawData): TeamMember {
  return {
    id: String(raw.id),
    teamID: String(raw.team_id),
    userID: String(raw.user_id),
    role: raw.role,
    displayName: raw.display_name ?? "",
    email: raw.email ?? "",
    avatar: raw.avatar ?? null,
    createdAt: raw.created_at,
    updatedAt: raw.updated_at,
  };
}

export function parseTeamMembers(rawList: TeamMemberRawData[]): TeamMember[] {
  return rawList.map(parseTeamMember);
}

