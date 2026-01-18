import { StatusDot } from "@/components/monitor/status-dot";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import type {
  PublicIncident,
  PublicStatusPageElement,
  PublicStatusPageMonitor,
} from "@/lib/schemas/public-status-page";
import { StatusTimeline } from "./status-timeline";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";

function MonitorRow({
  monitor,
  incidents,
}: {
  monitor: PublicStatusPageMonitor;
  incidents?: PublicIncident[];
}) {
  const isHistorical = monitor.type === "historical_timeline";
  const hasTimeline = isHistorical && (monitor.timeline?.length ?? 0) > 0;
  const monitorIncidents = incidents?.filter(
    (incident) => incident.monitorIds?.includes(monitor.monitorId),
  );

  return (
    <div className="space-y-2 rounded-lg border border-border/70 px-3 py-2">
      <div className="flex flex-wrap items-center gap-2">
        <StatusDot status={monitor.status} />
        <span className="truncate text-sm font-medium">{monitor.name}</span>
      </div>

      {isHistorical ? (
        hasTimeline ? (
          <StatusTimeline timeline={monitor.timeline} incidents={monitorIncidents} />
        ) : (
          <div className="text-xs text-muted-foreground">
            No historical timeline found for this monitor.
          </div>
        )
      ) : null}
    </div>
  );
}

function StatusElementCard({
  element,
  incidents,
}: {
  element: PublicStatusPageElement;
  incidents?: PublicIncident[];
}) {
  const isHistorical = element.type === "historical_timeline";
  const hasTimeline = isHistorical && (element.timeline?.length ?? 0) > 0;
  const hasMonitors = (element.monitors?.length ?? 0) > 0;
  const isGroup = !element.monitor;

  const timelineSection = hasTimeline && (
    <StatusTimeline timeline={element.timeline} incidents={incidents} label="Daily uptime (UTC)" />
  );

  const monitorsSection = hasMonitors && (
    <div className="space-y-2">
      {element.monitors.map((monitor) => (
        <MonitorRow key={monitor.id} monitor={monitor} incidents={incidents} />
      ))}
    </div>
  );

  if (!isGroup) {
    return (
      <Card className="border-border/80 shadow-sm py-3 gap-0">
        <CardHeader className="gap-1">
          <div className="flex flex-wrap items-center gap-2">
            <StatusDot status={element.status} />
            <CardTitle className="text-lg font-semibold">
              {element.name}
            </CardTitle>
          </div>
        </CardHeader>
        <CardContent className="flex flex-col gap-4 text-balance">
          {timelineSection}
          {monitorsSection}
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="border-border/80 shadow-sm py-3 gap-0">
      <Accordion type="single" collapsible>
        <AccordionItem value="group" className="group">
          <CardHeader className="gap-1">
            <AccordionTrigger className="py-0 flex items-center">
              <div className="flex flex-wrap items-center gap-2">
                <StatusDot status={element.status} />
                <CardTitle className="text-lg font-semibold">
                  {element.name}
                </CardTitle>
              </div>
            </AccordionTrigger>
          </CardHeader>
          <CardContent className="space-y-0 ">
            {timelineSection}
            <AccordionContent className="flex flex-col gap-4 text-balance mt-4">
              {monitorsSection}
            </AccordionContent>
          </CardContent>
        </AccordionItem>
      </Accordion>
    </Card>
  );
}

export function PublicStatusElements({
  elements,
  incidents,
}: {
  elements: PublicStatusPageElement[];
  incidents?: PublicIncident[];
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
        <StatusElementCard key={element.id} element={element} incidents={incidents} />
      ))}
    </div>
  );
}
