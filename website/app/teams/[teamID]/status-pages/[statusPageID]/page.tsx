import Link from "next/link";
import { redirect } from "next/navigation";

import { Button } from "@/components/ui/button";
import { buildCookieHeader } from "@/lib/api/cookies";
import { parseStatusPageDetail } from "@/lib/parsers/status-pages";
import type { StatusPageDetailItem, StatusPageRawData } from "@/lib/schemas/status-page";

type StatusPageResponse = {
  message?: string;
  data?: StatusPageRawData;
};

async function fetchStatusPage(
  teamID: string,
  statusPageID: string,
): Promise<StatusPageDetailItem> {
  const apiBase = process.env.NEXT_PUBLIC_API_BASE_URL;
  if (!apiBase) {
    throw new Error("Missing NEXT_PUBLIC_API_BASE_URL");
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

export default async function StatusPageDetail({
  params,
}: {
  params: Promise<{ teamID: string; statusPageID: string }>;
}) {
  const { teamID, statusPageID } = await params;
  const statusPage = await fetchStatusPage(teamID, statusPageID);

  return (
    <div className="flex flex-col gap-4">
      <div className="flex justify-between items-center">
        <div className="flex flex-col gap-1">
          <span className="text-xl font-bold">{statusPage.title}</span>
          <span className="text-sm text-muted-foreground">/{statusPage.slug}</span>
        </div>

        <Link href={`/teams/${teamID}/status-pages`}>
          <Button variant="outline">Back to list</Button>
        </Link>
      </div>

      <div className="rounded-md border border-dashed px-3 py-2 text-sm text-muted-foreground">
        Status page editor UI is coming soon. This page currently shows basic details
        only.
      </div>

      <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
        <div className="rounded-md border px-3 py-2">
          <div className="text-sm text-muted-foreground">Elements</div>
          <div className="text-lg font-semibold">{statusPage.elementCount}</div>
        </div>
        <div className="rounded-md border px-3 py-2">
          <div className="text-sm text-muted-foreground">Monitors</div>
          <div className="text-lg font-semibold">{statusPage.monitorCount}</div>
        </div>
      </div>
    </div>
  );
}

