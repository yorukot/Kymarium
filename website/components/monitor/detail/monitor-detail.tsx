import Link from "next/link";
import {
  CircleX,
  Clock,
  Hash,
  Pencil,
  Percent,
  Timer,
  TriangleAlert,
} from "lucide-react";

import { parseMonitor } from "@/lib/parsers/monitors";
import { formatRelativeTime } from "@/lib/parsers/datetime";
import { formatCount, formatLatencyMs, formatUptime } from "@/lib/parsers/format";
import type { MonitorAnalyticsRawData } from "@/lib/schemas/monitor-analytics";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { StatusDot } from "../status-dot";
import { MonitorStatusCard } from "./monitor-status-card";
import { MonitorTimeRangeSelect } from "./monitor-time-range-select";

export function MonitorDetail({
  teamID,
  monitorID,
  analytics,
}: {
  teamID: string;
  monitorID: string;
  analytics: MonitorAnalyticsRawData;
}) {
  const monitor = analytics.monitor;
  const parsed = parseMonitor(monitor);
  const summary = analytics.summary;
  const incidents = analytics.incidents ?? [];

  const totalCount = summary?.total_count ?? null;
  const goodCount = summary?.good_count ?? null;
  const failureCount =
    totalCount !== null && goodCount !== null
      ? Math.max(0, totalCount - goodCount)
      : null;

  const uptimeTone = monitor.status === "up" ? "successed" : "destructive";
  const failTone =
    failureCount !== null && failureCount > 0 ? "destructive" : "successed";

  return (
    <>
      <div className="flex gap-3 justify-between">
        <div className="flex min-w-0 flex-col gap-1">
          <div className="flex items-center gap-2 min-w-0">
            <StatusDot status={monitor.status} />
            <h1 className="text-xl font-bold truncate">{monitor.name}</h1>
            <Badge variant="outline">{monitor.type}</Badge>
          </div>
          <div className="text-sm text-muted-foreground">
            {parsed.targetLabel}:{" "}
            <span className="text-foreground">
              {parsed.targetValue || "--"}
            </span>
          </div>
        </div>

        <div className="flex gap-2">
          <Link href={`/teams/${teamID}/monitors/${monitorID}/edit`}>
            <Button variant="secondary">
              <Pencil />
              Edit
            </Button>
          </Link>
        </div>
      </div>
      <MonitorTimeRangeSelect />
      <div className="grid gap-4 grid-cols-2 sm:grid-cols-3 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5">
        <MonitorStatusCard
          title="UPTIME"
          value={formatUptime(summary?.uptime_pct)}
          tone={uptimeTone}
          icon={<Percent className="h-4 w-4" />}
        />
        <MonitorStatusCard
          title="INCIDENT"
          value={formatCount(incidents.length)}
          tone={"destructive"}
          icon={<TriangleAlert className="h-4 w-4" />}
        />
        <MonitorStatusCard
          title="FAIL"
          value={formatCount(failureCount)}
          tone={failTone}
          icon={<CircleX className="h-4 w-4" />}
        />
        <MonitorStatusCard
          title="TOTAL"
          value={formatCount(totalCount)}
          tone="default"
          icon={<Hash className="h-4 w-4" />}
        />
        <MonitorStatusCard
          title="LAST CHECKED"
          value={formatRelativeTime(monitor.last_checked)}
          tone="default"
          icon={<Clock className="h-4 w-4" />}
        />

        <MonitorStatusCard
          title="P50"
          value={formatLatencyMs(summary?.p50_ms)}
          tone="default"
          icon={<Timer className="h-4 w-4" />}
        />
        <MonitorStatusCard
          title="P75"
          value={formatLatencyMs(summary?.p75_ms)}
          tone="default"
          icon={<Timer className="h-4 w-4" />}
        />
        <MonitorStatusCard
          title="P90"
          value={formatLatencyMs(summary?.p90_ms)}
          tone="default"
          icon={<Timer className="h-4 w-4" />}
        />
        <MonitorStatusCard
          title="P95"
          value={formatLatencyMs(summary?.p95_ms)}
          tone="default"
          icon={<Timer className="h-4 w-4" />}
        />
        <MonitorStatusCard
          title="P99"
          value={formatLatencyMs(summary?.p99_ms)}
          tone="default"
          icon={<Timer className="h-4 w-4" />}
        />
      </div>
    </>
  );
}
