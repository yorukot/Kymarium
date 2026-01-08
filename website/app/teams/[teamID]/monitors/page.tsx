import Monitor from "@/components/monitor/monitor";
import { Button } from "@/components/ui/button";
import { buildCookieHeader } from "@/lib/api/cookies";
import { parseMonitors } from "@/lib/parsers/monitors";
import type { MonitorListItem, MonitorRawData } from "@/lib/schemas/monitor";
import { Plus } from "lucide-react";
import Link from "next/link";
import { redirect } from "next/navigation";

export default async function MonitorsPage({
  params,
}: {
  params: Promise<{ teamID: string }>;
}) {
  const { teamID } = await params;
  const monitors = await fetchMonitors(teamID);

  return (
    <div className="flex flex-col gap-4">
      <div className="flex justify-between items-center">
        <div className="flex flex-col gap-1">
          <span className="text-xl font-bold">Monitor List</span>
          <span className="text-sm text-muted-foreground">
            All the monitor are down below
          </span>
        </div>
        <Link href={`/teams/${teamID}/monitors/new`}>
          <Button>
            <Plus />
            Add Monitor
          </Button>
        </Link>
      </div>
      <Monitor monitors={monitors} />
    </div>
  );
}

type MonitorsResponse = {
  message?: string;
  data?: MonitorRawData[];
};

async function fetchMonitors(teamID: string): Promise<MonitorListItem[]> {
  const apiBase = process.env.NEXT_PUBLIC_API_BASE_URL;
  if (!apiBase) {
    throw new Error("Missing NEXT_PUBLIC_API_BASE_URL");
  }

  const cookieHeader = await buildCookieHeader();
  const res = await fetch(`${apiBase}/api/teams/${teamID}/monitors`, {
    method: "GET",
    headers: cookieHeader ? { cookie: cookieHeader } : undefined,
    cache: "no-store",
  });

  if (res.status === 401) {
    redirect(`/login?next=/teams/${teamID}/monitors`);
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
