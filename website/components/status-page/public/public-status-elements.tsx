import { Fragment } from "react";

import { StatusDot } from "@/components/monitor/status-dot";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import type {
  PublicStatusPageElement,
  PublicStatusPageMonitor,
} from "@/lib/schemas/public-status-page";
import { StatusTimeline } from "./status-timeline";

function MonitorRow({ monitor }: { monitor: PublicStatusPageMonitor }) {
  const isHistorical = monitor.type === "historical_timeline";
  const hasTimeline = isHistorical && (monitor.timeline?.length ?? 0) > 0;

  return (
    <div className="space-y-2 rounded-lg border border-border/70 px-3 py-2">
      <div className="flex flex-wrap items-center gap-2">
        <StatusDot status={monitor.status} />
        <span className="truncate text-sm font-medium">{monitor.name}</span>
      </div>

      {isHistorical ? (
        hasTimeline ? (
          <StatusTimeline timeline={monitor.timeline} />
        ) : (
          <div className="text-xs text-muted-foreground">
            No historical timeline found for this monitor.
          </div>
        )
      ) : null}
    </div>
  );
}

function StatusElementCard({ element }: { element: PublicStatusPageElement }) {
  const isHistorical = element.type === "historical_timeline";
  const hasTimeline = isHistorical && (element.timeline?.length ?? 0) > 0;
  const hasMonitors = (element.monitors?.length ?? 0) > 0;

  return (
    <Card className="border-border/80 shadow-sm py-1.5 gap-0">
      <CardHeader className="gap-1">
        <div className="flex flex-wrap items-center gap-2">
          <StatusDot status={element.status} />
          <CardTitle className="text-lg font-semibold">
            {element.name}
          </CardTitle>
        </div>
      </CardHeader>
      <CardContent className="space-y-4">
        <Fragment>
          {hasTimeline && (
            <StatusTimeline
              timeline={element.timeline}
              label="Daily uptime (UTC)"
            />
          )}
        </Fragment>
        {hasMonitors ? (
          <div className="space-y-2">
            <div className="text-xs font-semibold uppercase tracking-wide text-muted-foreground">
              Monitors
            </div>
            <div className="space-y-2">
              {element.monitors.map((monitor) => (
                <MonitorRow key={monitor.id} monitor={monitor} />
              ))}
            </div>
          </div>
        ) : null}
      </CardContent>
    </Card>
  );
}

export function PublicStatusElements({
  elements,
}: {
  elements: PublicStatusPageElement[];
}) {
  if (!elements.length) {
    return (
      <div className="rounded-xl border border-dashed px-4 py-3 text-sm text-muted-foreground">
        No status page elements configured yet.
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {elements.map((element) => (
        <StatusElementCard key={element.id} element={element} />
      ))}
    </div>
  );
}
