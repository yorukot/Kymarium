import { notFound } from "next/navigation";
import { Globe2, RefreshCcw } from "lucide-react";

import { StatusDot } from "@/components/monitor/status-dot";
import { PublicIncidentList } from "@/components/status-page/public/public-incident-list";
import { PublicStatusElements } from "@/components/status-page/public/public-status-elements";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { formatRelativeTime } from "@/lib/parsers/datetime";
import { parsePublicStatusPage } from "@/lib/parsers/public-status-page";
import type {
  PublicStatusPageData,
  PublicStatusPageResponseRaw,
} from "@/lib/schemas/public-status-page";

export const dynamic = "force-dynamic";

type ApiResponse = {
  message?: string;
  data?: PublicStatusPageResponseRaw;
};

async function fetchPublicStatusPage(
  slug: string,
): Promise<PublicStatusPageData | null> {
  const apiBase = process.env.NEXT_PUBLIC_API_BASE_URL;
  if (!apiBase) {
    throw new Error("Missing NEXT_PUBLIC_API_BASE_URL");
  }

  const res = await fetch(`${apiBase}/api/status-pages/${slug}`, {
    method: "GET",
    cache: "no-store",
  });

  if (res.status === 404) {
    return null;
  }

  if (!res.ok) {
    throw new Error("Failed to load public status page");
  }

  const body = (await res.json()) as ApiResponse;
  return parsePublicStatusPage(body?.data);
}

function overallStatus(data: PublicStatusPageData): string {
  const hasOpenIncident = data.incidents.some(
    (incident) => (incident.status ?? "").toLowerCase() !== "resolved",
  );
  if (hasOpenIncident) return "down";

  const elementDown = data.elements.some(
    (element) => (element.status ?? "").toLowerCase() === "down",
  );
  if (elementDown) return "down";

  return "up";
}

export default async function PublicStatusPage({
  params,
}: {
  params: Promise<{ slug: string }>;
}) {
  const { slug } = await params;
  const data = await fetchPublicStatusPage(slug);

  if (!data) {
    notFound();
  }

  const status = overallStatus(data);
  const updatedLabel = formatRelativeTime(data.statusPage.updatedAt);

  return (
    <main className="bg-background">
      <div className="mx-auto flex max-w-6xl flex-col gap-8 px-4 py-8 sm:px-6 lg:px-8">
        <Card className="border-border/80">
          <CardHeader className="gap-2">
            <div className="flex flex-wrap items-center gap-3">
              <StatusDot status={status} className="h-5 w-5" />
              <CardTitle className="text-2xl font-bold">
                {data.statusPage.title}
              </CardTitle>
            </div>
            <CardDescription className="flex flex-wrap items-center gap-4 text-sm">
              <span className="flex items-center gap-1">
                <Globe2 className="h-4 w-4" />
                /s/{data.statusPage.slug}
              </span>
              <span className="flex items-center gap-1">
                <RefreshCcw className="h-4 w-4" />
                Updated {updatedLabel}
              </span>
            </CardDescription>
          </CardHeader>
          <CardContent className="text-sm text-muted-foreground">
            Public overview of monitors, incidents, and historical uptime. Day
            counts adapt to your device so mobile users see a shorter window,
            while desktop users get the full 90-day view.
          </CardContent>
        </Card>

        <section className="space-y-4">
          <PublicStatusElements elements={data.elements} />
        </section>

        <Separator />

        <section className="space-y-4">
          <div className="flex flex-col gap-1">
            <h2 className="text-lg font-semibold">Public incidents</h2>
          </div>
          <PublicIncidentList incidents={data.incidents} />
        </section>
      </div>
    </main>
  );
}
