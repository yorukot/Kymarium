import type { Metadata } from "next";
import { notFound } from "next/navigation";
import { Clock3, Flame, Shield, TimerReset } from "lucide-react";

import { StatusDot } from "@/components/monitor/status-dot";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { Timeline, type TimelineItem } from "@/components/ui/timeline";
import { parsePublicStatusPage } from "@/lib/parsers/public-status-page";
import { formatRelativeTime } from "@/lib/parsers/datetime";
import type {
  PublicStatusPageData,
  PublicStatusPageResponseRaw,
} from "@/lib/schemas/public-status-page";
import { humanizeIdentifier } from "@/lib/parsers/strings";
import { cn } from "@/lib/utils";

export const dynamic = "force-dynamic";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ slug: string; incidentId: string }>;
}): Promise<Metadata> {
  const { slug, incidentId } = await params;
  const data = await fetchPublicStatusPage(slug);
  if (!data) {
    return { title: "Incident" };
  }

  const incident = data.incidents.find((i) => i.id === incidentId);
  const incidentTitle = incident?.title || `Incident ${incidentId}`;
  return {
    title: `${incidentTitle} | ${data.statusPage.title}`,
  };
}

type ApiResponse = {
  message?: string;
  data?: PublicStatusPageResponseRaw;
};

async function fetchPublicStatusPage(
  slug: string,
): Promise<PublicStatusPageData | null> {
  const apiBase = process.env.NEXT_PUBLIC_BACKEND_URL;
  if (!apiBase) {
    throw new Error("Missing NEXT_PUBLIC_BACKEND_URL");
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
  return humanizeIdentifier(status || "status unknown");
}

function statusToneClasses(status: string) {
  const isUp = (status ?? "").toLowerCase() === "up";
  const isResolved = (status ?? "").toLowerCase() === "resolved";
  if (isUp || isResolved) {
    return "border-successed/40 bg-successed/10 text-successed";
  }
  return "border-destructive/40 bg-destructive/10 text-destructive";
}

function severityToneClasses(severity: string) {
  const value = (severity ?? "").toLowerCase();
  if (value === "critical" || value === "high") {
    return "bg-destructive/10 text-destructive border-destructive/30";
  }
  if (value === "medium") {
    return "bg-amber-100 text-amber-900 border-amber-200";
  }
  return "bg-muted text-muted-foreground border-muted-foreground/30";
}

function mapMonitorNames(data: PublicStatusPageData): Record<string, string> {
  return data.elements
    .flatMap((element) => element.monitors)
    .reduce<Record<string, string>>((acc, monitor) => {
      acc[monitor.monitorId] = monitor.name;
      return acc;
    }, {});
}

export default async function PublicIncidentDetailPage({
  params,
}: {
  params: Promise<{ slug: string; incidentId: string }>;
}) {
  const { slug, incidentId } = await params;
  const data = await fetchPublicStatusPage(slug);

  if (!data) {
    notFound();
  }

  const incident = data.incidents.find((i) => i.id === incidentId);
  if (!incident) {
    notFound();
  }

  const monitorNames = mapMonitorNames(data);
  const status = overallStatus(data);
  const startedLabel = formatRelativeTime(incident.startedAt);
  const resolvedLabel = incident.resolvedAt
    ? formatRelativeTime(incident.resolvedAt)
    : null;
  const timeline = [...(incident.timeline ?? [])].sort((a, b) => {
    const aTime = new Date(a.createdAt).getTime();
    const bTime = new Date(b.createdAt).getTime();
    if (Number.isNaN(aTime) || Number.isNaN(bTime)) return 0;
    return bTime - aTime; // newest first
  });

  const timelineItems: TimelineItem[] = timeline.map((event) => {
    const at = new Date(event.createdAt);
    const absoluteLabel = Number.isNaN(at.getTime())
      ? event.createdAt
      : at.toLocaleString();

    return {
      id: event.id,
      title: (
        <div className="flex flex-wrap items-center gap-2">
          <span>{humanizeIdentifier(event.eventType)}</span>
          <span className="text-xs text-muted-foreground">
            {formatRelativeTime(event.createdAt)}
          </span>
        </div>
      ),
      description: <p className="text-sm leading-relaxed text-foreground">{event.message}</p>,
      children: (
        <div className="text-xs text-muted-foreground" suppressHydrationWarning>
          {absoluteLabel}
        </div>
      ),
    };
  });

  const relatedMonitorNames =
    incident.monitorIds?.map((id) => monitorNames[id] ?? id) ?? [];

  return (
    <main className="bg-background">
      <div className="mx-auto flex max-w-5xl flex-col gap-8 px-4 py-10 sm:px-6 lg:px-8">
        <div className="flex flex-col gap-4">
          <div className="flex flex-wrap items-center gap-3">
            <StatusDot status={status} className="h-5 w-5" />
            <CardTitle className="text-3xl font-bold">
              {data.statusPage.title}
            </CardTitle>
          </div>
          <div
            className={cn(
              "flex w-full items-center gap-2 rounded-xl border px-4 py-3 text-sm font-medium shadow-sm",
              statusToneClasses(status),
            )}
          >
            <StatusDot status={status} className="h-3 w-3" />
            <span>{statusLabel(status)}</span>
          </div>
        </div>

        <Card className="border-border/70 shadow-sm">
          <CardHeader className="space-y-3">
            <div className="flex flex-wrap items-center gap-3">
              <StatusDot status={incident.status} className="h-5 w-5" />
              <CardTitle className="text-2xl font-semibold">
                {incident.title || `Incident ${incident.id}`}
              </CardTitle>
              <Badge variant="outline" className={severityToneClasses(incident.severity)}>
                <Flame className="mr-1 h-3.5 w-3.5" />
                {humanizeIdentifier(incident.severity)}
              </Badge>
              <Badge
                variant="outline"
                className={cn(
                  "gap-1 border px-2 py-1 text-xs",
                  statusToneClasses(incident.status),
                )}
              >
                <Shield className="h-3.5 w-3.5" />
                {humanizeIdentifier(incident.status)}
              </Badge>
            </div>
            <div className="flex flex-wrap items-center gap-4 text-sm text-muted-foreground">
              <span className="flex items-center gap-2">
                <Clock3 className="h-4 w-4" />
                Started {startedLabel}
              </span>
              {resolvedLabel ? (
                <span className="flex items-center gap-2">
                  <TimerReset className="h-4 w-4" />
                  Resolved {resolvedLabel}
                </span>
              ) : null}
            </div>
          </CardHeader>
          <CardContent className="space-y-6">
            <div className="space-y-2">
              <h3 className="text-sm font-semibold uppercase tracking-wide text-muted-foreground">
                Related monitors
              </h3>
              {relatedMonitorNames.length ? (
                <div className="flex flex-wrap gap-2">
                  {relatedMonitorNames.map((label) => (
                    <Badge key={label} variant="secondary" className="bg-muted text-foreground">
                      {label}
                    </Badge>
                  ))}
                </div>
              ) : (
                <div className="text-sm text-muted-foreground">No related monitors listed.</div>
              )}
            </div>

            <Separator />

            {!timeline.length ? (
              <div className="rounded-lg border border-dashed px-4 py-3 text-sm text-muted-foreground">
                No public updates yet.
              </div>
            ) : (
              <Timeline items={timelineItems} compact />
            )}
          </CardContent>
        </Card>
      </div>
    </main>
  );
}
