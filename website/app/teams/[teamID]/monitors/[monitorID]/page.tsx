import type { Metadata } from "next";
import { notFound, redirect } from "next/navigation";

import { buildCookieHeader } from "@/lib/api/cookies";
import { parseRegions } from "@/lib/parsers/regions";
import type { MonitorAnalyticsRawData } from "@/lib/schemas/monitor-analytics";
import { isRFC3339Seconds } from "@/lib/parsers/datetime";
import type { Region, RegionRawData } from "@/lib/schemas/region";
import { MonitorHydrator } from "./monitor-hydrator";
import { MonitorDetail } from "@/components/monitor/detail/monitor-detail";
import MonitorChartBar from "@/components/monitor/detail/monitor-chart";

export const metadata: Metadata = {
  title: "Monitor Details",
};

type MonitorAnalyticsResponse = {
  message?: string;
  data?: MonitorAnalyticsRawData;
};

type RegionsResponse = {
  message?: string;
  data?: RegionRawData[];
};

async function fetchMonitorAnalytics(
  teamID: string,
  monitorID: string,
  window?: {
    start?: string;
    end?: string;
  },
): Promise<MonitorAnalyticsRawData | null> {
  const apiBase = process.env.NEXT_PUBLIC_BACKEND_URL;
  if (!apiBase) {
    throw new Error("Missing NEXT_PUBLIC_BACKEND_URL");
  }

  const cookieHeader = await buildCookieHeader();
  const params = new URLSearchParams();
  if (window?.start) params.set("start", window.start);
  if (window?.end) params.set("end", window.end);

  const res = await fetch(
    `${apiBase}/api/teams/${teamID}/monitors/${monitorID}/analytics${
      params.size ? `?${params.toString()}` : ""
    }`,
    {
      method: "GET",
      headers: cookieHeader ? { cookie: cookieHeader } : undefined,
      cache: "no-store",
    },
  );

  if (res.status === 401) {
    redirect(`/login?next=/teams/${teamID}/monitors/${monitorID}`);
  }

  if (res.status === 404) {
    return null;
  }

  if (!res.ok) {
    throw new Error("Failed to load monitor");
  }

  const body = (await res.json()) as MonitorAnalyticsResponse;
  return body?.data ?? null;
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
    redirect(`/login?next=/teams/${teamID}/monitors/${monitorID}`);
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

export default async function MonitorSummaryPage({
  params,
  searchParams,
}: {
  params: Promise<{ teamID: string; monitorID: string }>;
  searchParams?: Promise<{ start?: string; end?: string; preset?: string }>;
}) {
  const { teamID, monitorID } = await params;
  const qp = (await searchParams) ?? {};
  const start =
    typeof qp.start === "string" && isRFC3339Seconds(qp.start)
      ? qp.start
      : undefined;
  const end =
    typeof qp.end === "string" && isRFC3339Seconds(qp.end) ? qp.end : undefined;

  const [analytics, regions] = await Promise.all([
    fetchMonitorAnalytics(teamID, monitorID, { start, end }),
    fetchRegions(teamID, monitorID),
  ]);

  if (!analytics?.monitor) {
    notFound();
  }

  return (
    <div className="flex flex-col gap-4">
      <MonitorHydrator monitorID={monitorID} name={analytics.monitor.name} />
      <MonitorDetail
        teamID={teamID}
        monitorID={monitorID}
        analytics={analytics}
      />
      <MonitorChartBar analytics={analytics} regions={regions} />
    </div>
  );
}
