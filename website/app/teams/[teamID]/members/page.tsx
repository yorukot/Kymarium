import { redirect } from "next/navigation";

import { buildCookieHeader } from "@/lib/api/cookies";
import MembersPageClient from "@/components/member/members-page";
import { parseTeamInvites } from "@/lib/parsers/team-invites";
import { parseTeamMembers } from "@/lib/parsers/team-members";
import type { TeamInvite, TeamInviteRawData } from "@/lib/schemas/team-invite";
import type { TeamMember, TeamMemberRawData, TeamRole } from "@/lib/schemas/team-member";

type TeamResponse = {
  message?: string;
  data?: {
    role?: TeamRole;
  };
};

type MembersResponse = {
  message?: string;
  data?: TeamMemberRawData[];
};

type InvitesResponse = {
  message?: string;
  data?: TeamInviteRawData[];
};

async function fetchTeamRole(teamID: string): Promise<TeamRole> {
  const apiBase = process.env.NEXT_PUBLIC_API_BASE_URL;
  if (!apiBase) {
    throw new Error("Missing NEXT_PUBLIC_API_BASE_URL");
  }

  const cookieHeader = await buildCookieHeader();
  const res = await fetch(`${apiBase}/api/teams/${teamID}`, {
    method: "GET",
    headers: cookieHeader ? { cookie: cookieHeader } : undefined,
    cache: "no-store",
  });

  if (res.status === 401) {
    redirect(`/login?next=/teams/${teamID}/members`);
  }

  if (!res.ok) {
    throw new Error("Failed to load team");
  }

  const body = (await res.json()) as TeamResponse;
  return body?.data?.role ?? "member";
}

async function fetchMembers(teamID: string): Promise<TeamMember[]> {
  const apiBase = process.env.NEXT_PUBLIC_API_BASE_URL;
  if (!apiBase) {
    throw new Error("Missing NEXT_PUBLIC_API_BASE_URL");
  }

  const cookieHeader = await buildCookieHeader();
  const res = await fetch(`${apiBase}/api/teams/${teamID}/members`, {
    method: "GET",
    headers: cookieHeader ? { cookie: cookieHeader } : undefined,
    cache: "no-store",
  });

  if (res.status === 401) {
    redirect(`/login?next=/teams/${teamID}/members`);
  }

  if (!res.ok) {
    throw new Error("Failed to load members");
  }

  const body = (await res.json()) as MembersResponse;
  if (!Array.isArray(body?.data)) {
    return [];
  }

  return parseTeamMembers(body.data);
}

async function fetchInvites(teamID: string): Promise<TeamInvite[]> {
  const apiBase = process.env.NEXT_PUBLIC_API_BASE_URL;
  if (!apiBase) {
    throw new Error("Missing NEXT_PUBLIC_API_BASE_URL");
  }

  const cookieHeader = await buildCookieHeader();
  const res = await fetch(`${apiBase}/api/teams/${teamID}/invites`, {
    method: "GET",
    headers: cookieHeader ? { cookie: cookieHeader } : undefined,
    cache: "no-store",
  });

  if (res.status === 401) {
    redirect(`/login?next=/teams/${teamID}/members`);
  }

  if (!res.ok) {
    return [];
  }

  const body = (await res.json()) as InvitesResponse;
  if (!Array.isArray(body?.data)) {
    return [];
  }

  return parseTeamInvites(body.data);
}

export default async function MembersPage({
  params,
}: {
  params: Promise<{ teamID: string }>;
}) {
  const { teamID } = await params;

  const [teamRole, members] = await Promise.all([
    fetchTeamRole(teamID),
    fetchMembers(teamID),
  ]);

  const invites =
    teamRole === "owner" || teamRole === "admin" ? await fetchInvites(teamID) : [];

  return (
    <MembersPageClient
      teamID={teamID}
      teamRole={teamRole}
      members={members}
      invites={invites}
    />
  );
}

