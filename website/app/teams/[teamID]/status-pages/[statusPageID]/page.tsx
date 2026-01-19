import { redirect } from "next/navigation";

import StatusPageEditor from "@/components/status-page/status-page-editor";
import { buildCookieHeader } from "@/lib/api/cookies";
import { parseMonitors } from "@/lib/parsers/monitors";
import { parseStatusPageDetail } from "@/lib/parsers/status-pages";
import type { MonitorListItem, MonitorRawData } from "@/lib/schemas/monitor";
import type {
  StatusPageDetailItem,
  StatusPageRawData,
} from "@/lib/schemas/status-page";

type StatusPageResponse = {
  message?: string;
  data?: StatusPageRawData;
};

type MonitorsResponse = {
  message?: string;
  data?: MonitorRawData[];
};

async function fetchStatusPage(
  teamID: string,
  statusPageID: string,
): Promise<StatusPageDetailItem> {
  const apiBase = process.env.NEXT_PUBLIC_BACKEND_URL;
  if (!apiBase) {
    throw new Error("Missing NEXT_PUBLIC_BACKEND_URL");
  }

  const cookieHeader = await buildCookieHeader();
  const res = await fetch(
    `${apiBase}/api/teams/${teamID}/status-pages/${statusPageID}`,
    {
      method: "GET",
      headers: cookieHeader ? { cookie: cookieHeader } : undefined,
      cache: "no-store",
    },
  );

  if (res.status === 401) {
    redirect(`/login?next=/teams/${teamID}/status-pages/${statusPageID}`);
  }

  if (!res.ok) {
    throw new Error("Failed to load status page");
  }

  const body = (await res.json()) as StatusPageResponse;
  if (!body?.data) {
    throw new Error("Status page missing data");
  }

  return parseStatusPageDetail(body.data);
}

async function fetchMonitors(
  teamID: string,
  statusPageID: string,
): Promise<MonitorListItem[]> {
  const apiBase = process.env.NEXT_PUBLIC_BACKEND_URL;
  if (!apiBase) {
    throw new Error("Missing NEXT_PUBLIC_BACKEND_URL");
  }

  const cookieHeader = await buildCookieHeader();
  const res = await fetch(`${apiBase}/api/teams/${teamID}/monitors`, {
    method: "GET",
    headers: cookieHeader ? { cookie: cookieHeader } : undefined,
    cache: "no-store",
  });

  if (res.status === 401) {
    redirect(`/login?next=/teams/${teamID}/status-pages/${statusPageID}`);
  }

  if (!res.ok) {
    throw new Error("Failed to load monitors");
  }

  const body = (await res.json()) as MonitorsResponse;
  if (!Array.isArray(body?.data)) {
    return [];
  }

  return parseMonitors(body.data);
}

export default async function StatusPageDetail({
  params,
}: {
  params: Promise<{ teamID: string; statusPageID: string }>;
}) {
  const { teamID, statusPageID } = await params;
  const [statusPage, monitors] = await Promise.all([
    fetchStatusPage(teamID, statusPageID),
    fetchMonitors(teamID, statusPageID),
  ]);

  return (
    <StatusPageEditor
      teamID={teamID}
      statusPage={statusPage}
      monitorOptions={monitors}
    />
  );
}
