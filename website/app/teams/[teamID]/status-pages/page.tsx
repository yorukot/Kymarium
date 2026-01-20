import type { Metadata } from "next";
import Link from "next/link";
import { redirect } from "next/navigation";
import { Plus } from "lucide-react";

import StatusPageList from "@/components/status-page/status-page";
import { Button } from "@/components/ui/button";
import { buildCookieHeader } from "@/lib/api/cookies";
import { parseStatusPages } from "@/lib/parsers/status-pages";
import type { StatusPageListItem, StatusPageRawData } from "@/lib/schemas/status-page";

export const metadata: Metadata = {
  title: "Status Pages",
};

type StatusPagesResponse = {
  message?: string;
  data?: StatusPageRawData[];
};

async function fetchStatusPages(teamID: string): Promise<StatusPageListItem[]> {
  const apiBase = process.env.NEXT_PUBLIC_BACKEND_URL;
  if (!apiBase) {
    throw new Error("Missing NEXT_PUBLIC_BACKEND_URL");
  }

  const cookieHeader = await buildCookieHeader();
  const res = await fetch(`${apiBase}/api/teams/${teamID}/status-pages`, {
    method: "GET",
    headers: cookieHeader ? { cookie: cookieHeader } : undefined,
    cache: "no-store",
  });

  if (res.status === 401) {
    redirect(`/login?next=/teams/${teamID}/status-pages`);
  }

  if (!res.ok) {
    throw new Error("Failed to load status pages");
  }

  const body = (await res.json()) as StatusPagesResponse;
  if (!Array.isArray(body?.data)) {
    return [];
  }

  return parseStatusPages(body.data);
}

export default async function StatusPagesPage({
  params,
}: {
  params: Promise<{ teamID: string }>;
}) {
  const { teamID } = await params;
  const statusPages = await fetchStatusPages(teamID);

  return (
    <div className="flex flex-col gap-4">
      <div className="flex justify-between items-center">
        <div className="flex flex-col gap-1">
          <span className="text-xl font-bold">Status Page List</span>
          <span className="text-sm text-muted-foreground">
            All the status pages are down below
          </span>
        </div>
        <Link href={`/teams/${teamID}/status-pages/new`}>
          <Button>
            <Plus />
            Add Status Page
          </Button>
        </Link>
      </div>

      {statusPages.length === 0 ? (
        <div className="rounded-md border border-dashed px-3 py-2 text-sm text-muted-foreground">
          No status pages yet. Create one to publish a public status page.
        </div>
      ) : (
        <StatusPageList teamID={teamID} statusPages={statusPages} />
      )}
    </div>
  );
}
