"use client";

import Link from "next/link";
import { useMemo } from "react";
import { Clock3, Ticket } from "lucide-react";

import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { Timeline, type TimelineItem } from "@/components/ui/timeline";
import type { PublicIncident } from "@/lib/schemas/public-status-page";
import { humanizeIdentifier } from "@/lib/parsers/strings";
import { formatRelativeTime } from "@/lib/parsers/datetime";
import { IncidentStatusDot } from "@/components/incident/status-dot";

type Props = {
  incidents: PublicIncident[];
  monitorNames?: Record<string, string>;
  slug?: string;
};

export function PublicIncidentList({
  incidents,
  monitorNames = {},
  slug,
}: Props) {
  const grouped = useMemo(() => {
    type Aggregated = Omit<PublicIncident, "monitorIds"> & {
      monitorIds: Set<string>;
    };
    const byKey = new Map<string, Aggregated>();

    incidents.forEach((incident) => {
      const key = `${(incident.title ?? "").trim().toLowerCase() || incident.id}-${incident.startedAt}`;
      const ids = incident.monitorIds?.length ? incident.monitorIds : [];
      const existing = byKey.get(key);
      if (!existing) {
        byKey.set(key, { ...incident, monitorIds: new Set(ids) });
        return;
      }

      ids.forEach((id) => existing.monitorIds.add(id));
      // Prefer a resolved timestamp if any copy has it.
      if (!existing.resolvedAt && incident.resolvedAt) {
        existing.resolvedAt = incident.resolvedAt;
      }
      // Merge timelines, deduping by event id.
      const combinedTimeline = [...(existing.timeline ?? []), ...(incident.timeline ?? [])];
      const deduped = Array.from(
        new Map(combinedTimeline.map((e) => [e.id, e])).values(),
      );
      existing.timeline = deduped;
    });

    return Array.from(byKey.values())
      .map((item) => ({
        ...item,
        monitorIds: Array.from(item.monitorIds),
      }))
      .sort((a, b) => {
        const aTime = new Date(a.startedAt).getTime();
        const bTime = new Date(b.startedAt).getTime();
        return Number.isNaN(bTime) || Number.isNaN(aTime) ? 0 : bTime - aTime;
      });
  }, [incidents]);

  if (!incidents.length) {
    return (
      <div className="rounded-xl border border-dashed px-4 py-3 text-sm text-muted-foreground">
        No public incidents have been published.
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {grouped.map((incident) => (
        <IncidentCard
          key={incident.id}
          incident={incident}
          monitorNames={monitorNames}
          slug={slug}
        />
      ))}
    </div>
  );
}

type AggregatedIncident = PublicIncident & { monitorIds?: string[] };

function IncidentCard({
  incident,
  monitorNames,
  slug,
}: {
  incident: AggregatedIncident;
  monitorNames: Record<string, string>;
  slug?: string;
}) {
  const isResolved =
    (incident.status ?? "").toLowerCase() === "resolved" || Boolean(incident.resolvedAt);
  const timeline = useMemo(
    () =>
      [...(incident.timeline ?? [])].sort((a, b) => {
        const aTime = new Date(a.createdAt).getTime();
        const bTime = new Date(b.createdAt).getTime();
        if (Number.isNaN(aTime) || Number.isNaN(bTime)) return 0;
        // Newest first
        return bTime - aTime;
      }),
    [incident.timeline],
  );

  const timelineItems: TimelineItem[] = timeline.map((event) => {
    const at = new Date(event.createdAt);
    const absoluteLabel = Number.isNaN(at.getTime())
      ? event.createdAt
      : at.toLocaleString();

    return {
      id: event.id,
      title: (
        <div className="flex flex-wrap items-center gap-2">
          <span className="capitalize">{humanizeIdentifier(event.eventType)}</span>
          <span className="text-xs text-muted-foreground">
            {formatRelativeTime(event.createdAt)}
          </span>
        </div>
      ),
      description: (
        <p className="text-sm leading-relaxed text-foreground">{event.message}</p>
      ),
      children: (
        <div className="text-xs text-muted-foreground" suppressHydrationWarning>
          {absoluteLabel}
        </div>
      ),
    };
  });

  const startedLabel = formatRelativeTime(incident.startedAt);
  const resolvedLabel = incident.resolvedAt ? formatRelativeTime(incident.resolvedAt) : null;

  if (isResolved) {
    const body = (
      <Card className="border-border/80">
        <CardHeader className="gap-1">
          <div className="flex flex-wrap items-center gap-2">
            <IncidentStatusDot status={incident.status} />
            <CardTitle className="text-lg font-semibold">
              {incident.title || `Incident ${incident.id}`}
            </CardTitle>
          </div>
          <CardDescription className="flex flex-wrap items-center gap-3 text-sm">
            <span className="flex items-center gap-1">
              <Clock3 className="h-4 w-4" />
              Started {startedLabel}
            </span>
            {resolvedLabel ? (
              <span className="flex items-center gap-1">
                <Clock3 className="h-4 w-4" />
                Resolved {resolvedLabel}
              </span>
            ) : null}
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-3">
          {incident.monitorIds?.length ? (
            <div className="space-y-2">
              <Separator />
              <div className="flex flex-wrap items-center gap-2 text-xs text-muted-foreground">
                <span className="font-semibold text-foreground/80">Related monitors:</span>
                {incident.monitorIds.map((id) => {
                  const label = monitorNames[id] ?? id;
                  return (
                    <Badge key={id} variant="outline" className="bg-muted text-[11px]">
                      {label}
                    </Badge>
                  );
                })}
              </div>
            </div>
          ) : null}
        </CardContent>
      </Card>
    );

    if (slug) {
      return (
        <Link key={incident.id} href={`/s/${slug}/${incident.id}`} className="block">
          {body}
        </Link>
      );
    }

    return body;
  }

  return (
    <Card className="border-border/80">
      <CardHeader className="gap-1">
        <div className="flex flex-wrap items-center gap-2">
          <IncidentStatusDot status={incident.status} />
          <CardTitle className="text-lg font-semibold">
            {incident.title || `Incident ${incident.id}`}
          </CardTitle>
        </div>
        <CardDescription className="flex flex-wrap items-center gap-3 text-sm">
          <span className="flex items-center gap-1">
            <Clock3 className="h-4 w-4" />
            Started {startedLabel}
          </span>
          {resolvedLabel ? (
            <span className="flex items-center gap-1">
              <Clock3 className="h-4 w-4" />
              Resolved {resolvedLabel}
            </span>
          ) : null}
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-3">
        <Separator />
        {!timeline.length ? (
          <div className="text-sm text-muted-foreground">No public updates yet.</div>
        ) : (
          <Timeline items={timelineItems} compact />
        )}

        {incident.monitorIds?.length ? (
          <div className="space-y-2">
            <Separator />
            <div className="flex flex-wrap items-center gap-2 text-xs text-muted-foreground">
              <span className="font-semibold text-foreground/80">Related monitors:</span>
              {incident.monitorIds.map((id) => {
                const label = monitorNames[id] ?? id;
                return (
                  <Badge key={id} variant="outline" className="bg-muted text-[11px]">
                    {label}
                  </Badge>
                );
              })}
            </div>
          </div>
        ) : null}
      </CardContent>
    </Card>
  );
}
