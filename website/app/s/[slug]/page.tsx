import { notFound } from "next/navigation";
import { StatusDot } from "@/components/monitor/status-dot";
import { PublicIncidentList } from "@/components/status-page/public/public-incident-list";
import { PublicStatusElements } from "@/components/status-page/public/public-status-elements";
import { CardTitle } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { parsePublicStatusPage } from "@/lib/parsers/public-status-page";
import { cn } from "@/lib/utils";
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

function statusLabel(status: string): string {
  if (status === "up") return "All Systems Operational";
  if (status === "down") return "Service Disruption";
  return "Status Unknown";
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
  const monitorNames = data.elements
    .flatMap((element) => element.monitors)
    .reduce<Record<string, string>>((acc, monitor) => {
      acc[monitor.monitorId] = monitor.name;
      return acc;
    }, {});
  const openIncidents = data.incidents.filter(
    (incident) => (incident.status ?? "").toLowerCase() !== "resolved",
  );
  const pastIncidents = data.incidents.filter(
    (incident) => (incident.status ?? "").toLowerCase() === "resolved",
  );

  return (
    <main className="bg-background">
      <div className="mx-auto flex max-w-6xl flex-col gap-8 px-4 py-8 sm:px-6 lg:px-8">
        <div className="flex w-full flex-col gap-3">
          <div className="flex flex-wrap items-center gap-3">
            <CardTitle className="text-2xl font-bold">
              {data.statusPage.title}
            </CardTitle>
          </div>
          <div
            className={cn(
              "flex w-full items-center gap-2 rounded-lg border px-4 py-2 text-sm font-medium",
              status === "up"
                ? "border-successed/40 bg-successed/10 text-successed"
                : "border-destructive/40 bg-destructive/10 text-destructive",
            )}
          >
            <StatusDot status={status} className="h-3 w-3" />
            <span>{statusLabel(status)}</span>
          </div>
        </div>
        {openIncidents.length > 0 && (
          <section className="space-y-4">
            <div className="flex flex-col gap-1">
              <h2 className="text-lg font-semibold">Open incidents</h2>
              <p className="text-sm text-muted-foreground">
                Public incidents that are currently active.
              </p>
            </div>
            <PublicIncidentList
              incidents={openIncidents}
              monitorNames={monitorNames}
              slug={slug}
            />
            <Separator />
          </section>
        )}

        <section className="space-y-4">
          <PublicStatusElements elements={data.elements} incidents={data.incidents} />
        </section>

        <Separator />

        <section className="space-y-4">
          <div className="flex flex-col gap-1">
            <h2 className="text-lg font-semibold">Public incidents</h2>
            <p className="text-sm text-muted-foreground">
              Previously resolved and historical incidents.
            </p>
          </div>
          {pastIncidents.length ? (
            <PublicIncidentList
              incidents={pastIncidents}
              monitorNames={monitorNames}
              slug={slug}
            />
          ) : (
            <div className="rounded-xl border border-dashed px-4 py-3 text-sm text-muted-foreground">
              No past public incidents yet.
            </div>
          )}
        </section>
      </div>
    </main>
  );
}
