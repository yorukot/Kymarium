"use client";

import { useMemo } from "react";
import { format } from "date-fns";
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
import type { PublicIncident } from "@/lib/schemas/public-status-page";
import { humanizeIdentifier } from "@/lib/parsers/strings";
import { formatRelativeTime } from "@/lib/parsers/datetime";
import { IncidentStatusDot } from "@/components/incident/status-dot";

export function PublicIncidentList({ incidents }: { incidents: PublicIncident[] }) {
  const grouped = useMemo(() => {
    type Aggregated = PublicIncident & { monitorIds: Set<string> };
    const byKey = new Map<string, Aggregated>();

    incidents.forEach((incident) => {
      const key = `${(incident.title ?? "").trim().toLowerCase() || incident.id}-${incident.startedAt}`;
      const existing = byKey.get(key);
      if (!existing) {
        byKey.set(key, { ...incident, monitorIds: new Set([incident.monitorId]) });
        return;
      }

      existing.monitorIds.add(incident.monitorId);
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
        <IncidentCard key={incident.id} incident={incident} />
      ))}
    </div>
  );
}

type AggregatedIncident = PublicIncident & { monitorIds?: string[] };

function IncidentCard({ incident }: { incident: AggregatedIncident }) {
  const timeline = useMemo(
    () =>
      [...(incident.timeline ?? [])].sort((a, b) => {
        const aTime = new Date(a.createdAt).getTime();
        const bTime = new Date(b.createdAt).getTime();
        return aTime - bTime;
      }),
    [incident.timeline],
  );

  const startedLabel = formatRelativeTime(incident.startedAt);
  const resolvedLabel = incident.resolvedAt ? formatRelativeTime(incident.resolvedAt) : null;

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
        <div className="space-y-2">
          <div className="flex items-center gap-2 text-sm font-medium text-foreground">
            <Ticket className="h-4 w-4" />
            Updates
          </div>
          {!timeline.length ? (
            <div className="text-sm text-muted-foreground">No public updates yet.</div>
          ) : (
            <ul className="space-y-3">
              {timeline.map((event) => {
                const at = new Date(event.createdAt);
                const atLabel = Number.isNaN(at.getTime())
                  ? event.createdAt
                  : `${format(at, "MMM d, yyyy HH:mm")} (${formatRelativeTime(event.createdAt)})`;
                return (
                  <li key={event.id} className="rounded-lg border border-border/70 bg-muted/30 p-3">
                    <div className="flex items-center justify-between gap-2 text-xs text-muted-foreground">
                      <span className="capitalize">{humanizeIdentifier(event.eventType)}</span>
                      <span>{atLabel}</span>
                    </div>
                    <p className="text-sm leading-relaxed text-foreground">{event.message}</p>
                  </li>
                );
              })}
            </ul>
          )}
        </div>

        {incident.monitorIds?.length ? (
          <div className="space-y-2">
            <Separator />
            <div className="flex flex-wrap items-center gap-2 text-xs text-muted-foreground">
              <span className="font-semibold text-foreground/80">Related monitors:</span>
              {incident.monitorIds.map((id) => (
                <Badge key={id} variant="outline" className="bg-muted text-[11px]">
                  {id}
                </Badge>
              ))}
            </div>
          </div>
        ) : null}
      </CardContent>
    </Card>
  );
}
