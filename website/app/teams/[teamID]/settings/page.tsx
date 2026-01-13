import { redirect } from "next/navigation";

import TeamSettingsPageClient from "@/components/team/team-settings-page";
import { buildCookieHeader } from "@/lib/api/cookies";
import type { TeamRole } from "@/lib/schemas/team-member";

type TeamResponse = {
  message?: string;
  data?: {
    id: string;
    name: string;
    role: TeamRole;
  };
};

async function fetchTeam(teamID: string): Promise<{ name: string; role: TeamRole }> {
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
    redirect(`/login?next=/teams/${teamID}/settings`);
  }

  if (!res.ok) {
    throw new Error("Failed to load team");
  }

  const body = (await res.json()) as TeamResponse;
  const team = body?.data;
  if (!team) {
    throw new Error("Failed to load team");
  }

  return {
    name: team.name ?? "",
    role: team.role ?? "member",
  };
}

export default async function TeamSettingsPage({
  params,
}: {
  params: Promise<{ teamID: string }>;
}) {
  const { teamID } = await params;
  const team = await fetchTeam(teamID);

  return (
    <TeamSettingsPageClient
      teamID={teamID}
      teamName={team.name}
      teamRole={team.role}
    />
  );
}

