import { redirect } from "next/navigation";

import { buildCookieHeader } from "@/lib/api/cookies";
import { Region } from "@/lib/schemas/region";
import { parseRegions } from "@/lib/parsers/regions";
import { Notification, NotificationRawData } from "@/lib/schemas/notification";
import { parseNotifications } from "@/lib/parsers/notifications";
import NewMonitorForm from "@/components/monitor/new/new-monitor";

type RegionsResponse = {
  message?: string;
  data?: Array<{
    id: string;
    name: string;
  }>;
};

type NotificationsResponse = {
  message?: string;
  data?: NotificationRawData[];
};

async function fetchRegions(teamID: string): Promise<Region[]> {
  const apiBase = process.env.NEXT_PUBLIC_API_BASE_URL;
  if (!apiBase) {
    throw new Error("Missing NEXT_PUBLIC_API_BASE_URL");
  }

  const cookieHeader = await buildCookieHeader();
  const res = await fetch(`${apiBase}/api/regions`, {
    method: "GET",
    headers: cookieHeader ? { cookie: cookieHeader } : undefined,
    cache: "no-store",
  });

  if (res.status === 401) {
    redirect(`/login?next=/teams/${teamID}/monitors/new`);
  }

  if (!res.ok) {
    throw new Error("Failed to load regions");
  }

  const body = (await res.json()) as RegionsResponse;
  if (!Array.isArray(body?.data)) {
    return [];
  }

  return parseRegions(body.data);
}

async function fetchNotifications(teamID: string): Promise<Notification[]> {
  const apiBase = process.env.NEXT_PUBLIC_API_BASE_URL;
  if (!apiBase) {
    throw new Error("Missing NEXT_PUBLIC_API_BASE_URL");
  }

  const cookieHeader = await buildCookieHeader();
  const res = await fetch(`${apiBase}/api/teams/${teamID}/notifications`, {
    method: "GET",
    headers: cookieHeader ? { cookie: cookieHeader } : undefined,
    cache: "no-store",
  });

  if (res.status === 401) {
    redirect(`/login?next=/teams/${teamID}/monitors/new`);
  }

  if (!res.ok) {
    throw new Error("Failed to load notifications");
  }

  const body = (await res.json()) as NotificationsResponse;
  if (!Array.isArray(body?.data)) {
    return [];
  }

  return parseNotifications(body.data);
}

export default async function NewMonitorPage({
  params,
}: {
  params: Promise<{ teamID: string }>;
}) {
  const { teamID } = await params;
  const [regions, notifications] = await Promise.all([
    fetchRegions(teamID),
    fetchNotifications(teamID),
  ]);

  return (
    <NewMonitorForm
      teamID={teamID}
      regions={regions}
      notifications={notifications}
    />
  );
}
