import type { Metadata } from "next";
import IncidentList from "@/components/incident/incident";
import { Button } from "@/components/ui/button";
import { buildCookieHeader } from "@/lib/api/cookies";
import { parseIncidents } from "@/lib/parsers/incidents";
import type { IncidentListItem, IncidentRawData } from "@/lib/schemas/incident";
import { Plus } from "lucide-react";
import Link from "next/link";
import { redirect } from "next/navigation";

export const metadata: Metadata = {
  title: "Incidents",
};

export default async function IncidentsPage({
  params,
}: {
  params: Promise<{ teamID: string }>;
}) {
  const { teamID } = await params;
  const incidents = await fetchIncidents(teamID);

  return (
    <div className="flex flex-col gap-4">
      <div className="flex justify-between items-center">
        <div className="flex flex-col gap-1">
          <span className="text-xl font-bold">Incident List</span>
          <span className="text-sm text-muted-foreground">
            All the incidents are down below
          </span>
        </div>
        <Link href={`/teams/${teamID}/incidents/new`}>
          <Button>
            <Plus />
            Add Incident
          </Button>
        </Link>
      </div>
      <IncidentList teamID={teamID} incidents={incidents} />
    </div>
  );
}

type IncidentsResponse = {
  message?: string;
  data?: IncidentRawData[];
};

async function fetchIncidents(teamID: string): Promise<IncidentListItem[]> {
  const apiBase = process.env.NEXT_PUBLIC_BACKEND_URL;
  if (!apiBase) {
    throw new Error("Missing NEXT_PUBLIC_BACKEND_URL");
  }

  const cookieHeader = await buildCookieHeader();
  const res = await fetch(`${apiBase}/api/teams/${teamID}/incidents`, {
    method: "GET",
    headers: cookieHeader ? { cookie: cookieHeader } : undefined,
    cache: "no-store",
  });

  if (res.status === 401) {
    redirect(`/login?next=/teams/${teamID}/incidents`);
  }

  if (!res.ok) {
    throw new Error("Failed to load incidents");
  }

  const body = (await res.json()) as IncidentsResponse;
  if (!Array.isArray(body?.data)) {
    return [];
  }

  return parseIncidents(body.data);
}
