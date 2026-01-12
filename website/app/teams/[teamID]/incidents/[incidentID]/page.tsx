import { notFound, redirect } from "next/navigation";

import { buildCookieHeader } from "@/lib/api/cookies";
import { parseIncidentDetail, parseIncidentEvents } from "@/lib/parsers/incidents";
import type {
  IncidentDetailItem,
  IncidentEventItem,
  IncidentEventRawData,
  IncidentRawData,
} from "@/lib/schemas/incident";
import { IncidentDetail } from "@/components/incident/detail/incident-detail";
import { IncidentHydrator } from "./incident-hydrator";

type IncidentResponse = {
  message?: string;
  data?: IncidentRawData;
};

type IncidentEventsResponse = {
  message?: string;
  data?: IncidentEventRawData[];
};

async function fetchIncident(
  teamID: string,
  incidentID: string,
): Promise<IncidentDetailItem | null> {
  const apiBase = process.env.NEXT_PUBLIC_API_BASE_URL;
  if (!apiBase) {
    throw new Error("Missing NEXT_PUBLIC_API_BASE_URL");
  }

  const cookieHeader = await buildCookieHeader();
  const res = await fetch(`${apiBase}/api/teams/${teamID}/incidents/${incidentID}`, {
    method: "GET",
    headers: cookieHeader ? { cookie: cookieHeader } : undefined,
    cache: "no-store",
  });

  if (res.status === 401) {
    redirect(`/login?next=/teams/${teamID}/incidents/${incidentID}`);
  }

  if (res.status === 404) {
    return null;
  }

  if (!res.ok) {
    throw new Error("Failed to load incident");
  }

  const body = (await res.json()) as IncidentResponse;
  if (!body?.data) return null;

  return parseIncidentDetail(body.data);
}

async function fetchIncidentEvents(
  teamID: string,
  incidentID: string,
): Promise<IncidentEventItem[]> {
  const apiBase = process.env.NEXT_PUBLIC_API_BASE_URL;
  if (!apiBase) {
    throw new Error("Missing NEXT_PUBLIC_API_BASE_URL");
  }

  const cookieHeader = await buildCookieHeader();
  const res = await fetch(
    `${apiBase}/api/teams/${teamID}/incidents/${incidentID}/events`,
    {
      method: "GET",
      headers: cookieHeader ? { cookie: cookieHeader } : undefined,
      cache: "no-store",
    },
  );

  if (res.status === 401) {
    redirect(`/login?next=/teams/${teamID}/incidents/${incidentID}`);
  }

  if (res.status === 404) {
    return [];
  }

  if (!res.ok) {
    throw new Error("Failed to load incident events");
  }

  const body = (await res.json()) as IncidentEventsResponse;
  if (!Array.isArray(body?.data)) {
    return [];
  }

  return parseIncidentEvents(body.data);
}

export default async function IncidentDetailPage({
  params,
}: {
  params: Promise<{ teamID: string; incidentID: string }>;
}) {
  const { teamID, incidentID } = await params;
  const [incident, events] = await Promise.all([
    fetchIncident(teamID, incidentID),
    fetchIncidentEvents(teamID, incidentID),
  ]);

  if (!incident) {
    notFound();
  }

  return (
    <div className="flex flex-col gap-4">
      <IncidentHydrator incidentID={incidentID} title={incident.title ?? ""} />
      <IncidentDetail teamID={teamID} incident={incident} events={events} />
    </div>
  );
}
