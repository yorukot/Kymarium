"use client";

import { useEffect, useMemo, useState } from "react";
import { format } from "date-fns";

import { cn } from "@/lib/utils";
import type { PublicTimelinePoint } from "@/lib/schemas/public-status-page";
import {
  HoverCard,
  HoverCardContent,
  HoverCardTrigger,
} from "@/components/ui/hover-card";

const DEVICE_DAY_WINDOWS = [
  { minWidth: 1280, days: 90 },
  { minWidth: 700, days: 60 },
  { minWidth: 0, days: 30 },
] as const;

function pickDaysForWidth(width: number): number {
  for (const window of DEVICE_DAY_WINDOWS) {
    if (width >= window.minWidth) {
      return window.days;
    }
  }
  return DEVICE_DAY_WINDOWS[DEVICE_DAY_WINDOWS.length - 1].days;
}

function useResponsiveTimelineDays(defaultDays = 30) {
  const [days, setDays] = useState(defaultDays);

  useEffect(() => {
    const update = () => setDays(pickDaysForWidth(window.innerWidth));
    update();
    window.addEventListener("resize", update);
    return () => window.removeEventListener("resize", update);
  }, []);

  return days;
}

function uptimePercent(point: PublicTimelinePoint): number | null {
  const total = point.success + point.fail;
  if (total <= 0) return null;
  return (point.success / total) * 100;
}

type TimelineTone = "success" | "warning" | "danger" | "muted";

function toneForUptime(pct: number | null): TimelineTone {
  if (pct === null) return "muted";
  if (pct >= 100) return "success";
  if (pct >= 99.9) return "warning";
  return "danger";
}

const toneClassMap: Record<TimelineTone, string> = {
  success: "bg-successed",
  warning: "bg-amber-400 dark:bg-amber-500",
  danger: "bg-destructive",
  muted: "bg-border",
};

function dayLabel(value: string): string {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return format(date, "MMM d");
}

export function StatusTimeline({
  timeline,
}: {
  timeline?: PublicTimelinePoint[];
  label?: string;
}) {
  const visibleDays = useResponsiveTimelineDays();
  const points = useMemo(
    () => (timeline ?? []).slice(-(visibleDays ?? 0)),
    [timeline, visibleDays],
  );

  const summaryUptime = useMemo(() => {
    const totalSuccess = points.reduce((acc, p) => acc + p.success, 0);
    const totalFail = points.reduce((acc, p) => acc + p.fail, 0);
    const total = totalSuccess + totalFail;
    if (total <= 0) return null;
    return (totalSuccess / total) * 100;
  }, [points]);

  const summaryUptimeLabel =
    summaryUptime === null ? "No data" : `${summaryUptime.toFixed(2)}% uptime`;

  if (!points.length) {
    return (
      <div className="text-xs text-muted-foreground">
        No historical checks yet.
      </div>
    );
  }

  return (
    <div className="space-y-0">
      <div className="flex items-end gap-0.5 overflow-x-auto pb-2 w-full">
        {points.map((point) => {
          const pct = uptimePercent(point);
          const tone = toneForUptime(pct);
          const uptimeLabel = pct === null ? "No data" : `${pct.toFixed(2)}%`;
          const date = dayLabel(point.day);

          return (
            <HoverCard
              key={`${point.day}-${point.success}-${point.fail}`}
              openDelay={50}
              closeDelay={0}
            >
              <HoverCardTrigger asChild>
                <div
                  className={cn(
                    "border border-border/60 h-8 transition-[transform,opacity] duration-150 ease-out hover:border-foreground/40 w-full",
                    toneClassMap[tone],
                  )}
                  aria-label={`${date}: ${uptimeLabel}`}
                />
              </HoverCardTrigger>
              <HoverCardContent className="text-xs" side="top" align="center">
                <div className="font-semibold">{date}</div>
                <div className="text-foreground">{uptimeLabel} uptime</div>
                <div className="text-muted-foreground text-[11px]">
                  Success {point.success} · Fail {point.fail}
                </div>
              </HoverCardContent>
            </HoverCard>
          );
        })}
      </div>
      <div className="flex items-center gap-2 text-xs text-muted-foreground">
        <span className="shrink-0">{points.length}days ago</span>

        <div className="flex-1 flex items-center min-w-0">
          <div className="h-px bg-border/70 w-full" />
          <span className="px-2 shrink-0 text-foreground/80">
            {summaryUptimeLabel}
          </span>
          <div className="h-px bg-border/70 w-full" />
        </div>

        <span className="shrink-0">now</span>
      </div>
    </div>
  );
}
