import { redirect } from "next/navigation";
import camelcaseKeys from "camelcase-keys";

import { buildCookieHeader } from "@/lib/api/cookies";
import { parseRegions } from "@/lib/parsers/regions";
import { parseNotifications } from "@/lib/parsers/notifications";
import type { Region } from "@/lib/schemas/region";
import type {
  Notification,
  NotificationRawData,
} from "@/lib/schemas/notification";
import type { MonitorFormValues, MonitorRawData } from "@/lib/schemas/monitor";
import { DEFAULT_HTTP, DEFAULT_PING } from "@/lib/schemas/monitor";
import NewMonitorForm from "@/components/monitor/new/new-monitor";

type HttpFormValues = Extract<MonitorFormValues, { type: "http" }>;
type PingFormValues = Extract<MonitorFormValues, { type: "ping" }>;

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

type MonitorResponse = {
  message?: string;
  data?: MonitorRawData;
};

function buildLoginRedirectPath(teamID: string, monitorID: string) {
  return `/login?next=/teams/${teamID}/monitors/${monitorID}/edit`;
}

async function fetchRegions(
  teamID: string,
  monitorID: string,
): Promise<Region[]> {
  const apiBase = process.env.NEXT_PUBLIC_BACKEND_URL;
  if (!apiBase) {
    throw new Error("Missing NEXT_PUBLIC_BACKEND_URL");
  }

  const cookieHeader = await buildCookieHeader();
  const res = await fetch(`${apiBase}/api/regions`, {
    method: "GET",
    headers: cookieHeader ? { cookie: cookieHeader } : undefined,
    cache: "no-store",
  });

  if (res.status === 401) {
    redirect(buildLoginRedirectPath(teamID, monitorID));
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

async function fetchNotifications(
  teamID: string,
  monitorID: string,
): Promise<Notification[]> {
  const apiBase = process.env.NEXT_PUBLIC_BACKEND_URL;
  if (!apiBase) {
    throw new Error("Missing NEXT_PUBLIC_BACKEND_URL");
  }

  const cookieHeader = await buildCookieHeader();
  const res = await fetch(`${apiBase}/api/teams/${teamID}/notifications`, {
    method: "GET",
    headers: cookieHeader ? { cookie: cookieHeader } : undefined,
    cache: "no-store",
  });

  if (res.status === 401) {
    redirect(buildLoginRedirectPath(teamID, monitorID));
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

async function fetchMonitor(
  teamID: string,
  monitorID: string,
): Promise<MonitorRawData | null> {
  const apiBase = process.env.NEXT_PUBLIC_BACKEND_URL;
  if (!apiBase) {
    throw new Error("Missing NEXT_PUBLIC_BACKEND_URL");
  }

  const cookieHeader = await buildCookieHeader();
  const res = await fetch(`${apiBase}/api/teams/${teamID}/monitors/${monitorID}`, {
    method: "GET",
    headers: cookieHeader ? { cookie: cookieHeader } : undefined,
    cache: "no-store",
  });

  if (res.status === 401) {
    redirect(buildLoginRedirectPath(teamID, monitorID));
  }

  if (!res.ok) {
    throw new Error("Failed to load monitor");
  }

  const body = (await res.json()) as MonitorResponse;
  return body?.data ?? null;
}

type MonitorConfig = Record<string, unknown>;

function buildInitialValues(
  monitor: MonitorRawData,
  fallbackRegionIDs: string[],
): MonitorFormValues {
  const base = {
    name: monitor.name ?? "",
    interval: monitor.interval ?? 60,
    failureThreshold: monitor.failure_threshold ?? 2,
    recoveryThreshold: monitor.recovery_threshold ?? 2,
    regions: Array.isArray(monitor.regions) ? monitor.regions : fallbackRegionIDs,
    notifications: Array.isArray(monitor.notification) ? monitor.notification : [],
  };

  const configRaw =
    monitor.config && typeof monitor.config === "object"
      ? (monitor.config as MonitorConfig)
      : {};

  const config = camelcaseKeys(configRaw, { deep: true }) as Record<
    string,
    unknown
  >;

  if (monitor.type !== "ping") {
    const http = {
      ...DEFAULT_HTTP,
      ...(config as Partial<HttpFormValues["http"]>),
    } as HttpFormValues["http"];
    const values: HttpFormValues = {
      ...base,
      type: "http",
      http,
      ping: undefined,
    };
    return values;
  }

  const ping = {
    ...DEFAULT_PING,
    ...(config as Partial<PingFormValues["ping"]>),
  } as PingFormValues["ping"];
  const values: PingFormValues = {
    ...base,
    type: "ping",
    ping,
    http: undefined,
  };
  return values;
}

export default async function EditMonitorPage({
  params,
}: {
  params: Promise<{ teamID: string; monitorID: string }>;
}) {
  const { teamID, monitorID } = await params;
  const [regions, notifications, monitor] = await Promise.all([
    fetchRegions(teamID, monitorID),
    fetchNotifications(teamID, monitorID),
    fetchMonitor(teamID, monitorID),
  ]);

  if (!monitor) {
    throw new Error("Monitor not found");
  }

  const defaultRegions = regions.map((region) => region.id);
  const initialValues = buildInitialValues(monitor, defaultRegions);

  return (
    <NewMonitorForm
      teamID={teamID}
      regions={regions}
      notifications={notifications}
      monitorID={monitorID}
      initialValues={initialValues}
    />
  );
}
