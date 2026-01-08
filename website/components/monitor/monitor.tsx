"use client";

import { Link, MoreVertical } from "lucide-react";
import { Card, CardContent } from "../ui/card";
import { Button } from "../ui/button";
import type { MonitorListItem } from "@/lib/schemas/monitor";
import { formatDistanceStrict } from "date-fns";

function formatRelativeTime(value: string): string {
  if (!value) return "--";

  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return "--";

  return formatDistanceStrict(date, new Date(Date.now()), {
    addSuffix: true,
    roundingMethod: "round",
  });
}

function formatUptime(value?: number) {
  if (value === undefined || value === null || Number.isNaN(value)) return "--";
  return `${value.toFixed(2)}%`;
}

export default function Monitor({ monitors }: { monitors: MonitorListItem[] }) {
  if (!monitors.length) return null;

  return (
    <>
      {monitors.map((monitor) => (
        <Card key={monitor.id} className="p-3.5">
          <CardContent className="flex items-center justify-between gap-3 p-0">
            <div className="flex min-w-0 flex-col gap-2">
              <div className="flex items-center gap-2 text-md min-w-0">
                <div className="inline-flex h-4 w-4 flex-none items-center justify-center rounded-full bg-destructive/40">
                  <div className="h-2 w-2 rounded-full bg-destructive/70" />
                </div>
                <span className="min-w-0 flex-1 truncate">{monitor.name}</span>
              </div>

              <div className="flex items-center gap-2 text-sm text-muted-foreground min-w-0">
                <Link size={14} className="flex-none shrink-0" aria-label={monitor.targetLabel} />
                <span className="min-w-0 flex-1 truncate">{monitor.targetValue || "--"}</span>
              </div>
            </div>

            <div className="flex flex-none gap-3">
              <div className="hidden md:flex flex-col items-center justify-start gap-1">
                <span className="text-muted-foreground text-sm">Last checked</span>
                <span suppressHydrationWarning>{formatRelativeTime(monitor.lastChecked)}</span>
              </div>

              <div className="flex flex-col items-center justify-start gap-1">
                <span className="text-muted-foreground text-sm">30d uptime</span>
                <span>{formatUptime(monitor.uptimeSLI30)}</span>
              </div>

              <div className="flex items-center justify-center">
                <Button size="icon" variant="ghost">
                  <MoreVertical size={16} />
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>
      ))}
    </>
  );
}
