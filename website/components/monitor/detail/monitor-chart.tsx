"use client";
import "flag-icons/css/flag-icons.min.css";

import { format } from "date-fns";
import { Area, AreaChart, CartesianGrid, XAxis } from "recharts";
import * as React from "react";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  ChartContainer,
  ChartLegend,
  ChartLegendContent,
  ChartTooltip,
  ChartTooltipContent,
  type ChartConfig,
} from "@/components/ui/chart";
import type { MonitorAnalyticsRawData } from "@/lib/schemas/monitor-analytics";
import { toStackedBarChartData } from "@/lib/parsers/monitor-analytics";
import type { Region } from "@/lib/schemas/region";
import { formatLatencyMs } from "@/lib/parsers/format";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

type LatencyMetric = "p50_ms" | "p75_ms" | "p90_ms" | "p95_ms" | "p99_ms";

const latencyMetricOptions: Array<{ value: LatencyMetric; label: string }> = [
  { value: "p50_ms", label: "P50" },
  { value: "p75_ms", label: "P75" },
  { value: "p90_ms", label: "P90" },
  { value: "p95_ms", label: "P95" },
  { value: "p99_ms", label: "P99" },
];

const chartColors = [
  "var(--chart-1)",
  "var(--chart-2)",
  "var(--chart-3)",
  "var(--chart-4)",
  "var(--chart-5)",
];

function safeParseMs(value?: string): number | null {
  if (!value) return null;
  const ms = new Date(value).getTime();
  return Number.isNaN(ms) ? null : ms;
}

function buildTimeFormatters(window?: MonitorAnalyticsRawData["window"]) {
  const startMs = safeParseMs(window?.start);
  const endMs = safeParseMs(window?.end);

  const windowMs =
    startMs !== null && endMs !== null && endMs > startMs ? endMs - startMs : 0;
  const windowDays = windowMs / (24 * 60 * 60 * 1000);

  const showDate = windowDays > 3;
  const showTime = windowDays <= 7;
  const includeYear =
    startMs !== null &&
    endMs !== null &&
    new Date(startMs).getFullYear() !== new Date(endMs).getFullYear();

  const xAxisFormat = (value: string) => {
    const date = new Date(value);
    if (!showDate) return format(date, "HH:mm");
    if (showTime) return format(date, "MMM d HH:mm");
    return format(date, includeYear ? "MMM d, yyyy" : "MMM d");
  };

  const tooltipLabelFormat = (value: unknown) => {
    const date = new Date(String(value));
    if (!showDate) return format(date, "PP p");
    if (showTime) return format(date, "PP p");
    return format(date, "PP");
  };

  return { xAxisFormat, tooltipLabelFormat };
}

function buildChartConfig(regionIDs: string[], regions: Region[]): ChartConfig {
  const nameByID = new Map(regions.map((r) => [r.id, r.displayName]));
  const flagByID = new Map(regions.map((r) => [r.id, r.flag]));

  return regionIDs.reduce<ChartConfig>((acc, regionID, index) => {
    const flag = flagByID.get(regionID) ?? "un";
    acc[regionID] = {
      label: nameByID.get(regionID) ?? `Region ${regionID}`,
      color: chartColors[index % chartColors.length],
      icon: () => (
        <span className="flex items-center gap-1">
          <span
            className="h-2 w-2 shrink-0 rounded-[2px]"
            style={{ backgroundColor: `var(--color-${regionID})` }}
          />
          <span className={`fi fi-${flag}`}></span>
        </span>
      ),
    };
    return acc;
  }, {});
}

export default function MonitorChartBar({
  analytics,
  regions,
}: {
  analytics: MonitorAnalyticsRawData;
  regions: Region[];
}) {
  const [metric, setMetric] = React.useState<LatencyMetric>("p75_ms");

  const regionIDs = analytics.monitor?.regions ?? [];
  const chartConfig = buildChartConfig(regionIDs, regions);
  const chartData = toStackedBarChartData(analytics, metric);
  const regionByID = new Map(regions.map((r) => [r.id, r]));
  const metricLabel =
    latencyMetricOptions.find((opt) => opt.value === metric)?.label ?? "P75";
  const { xAxisFormat, tooltipLabelFormat } = buildTimeFormatters(
    analytics.window,
  );

  return (
    <Card>
      <CardHeader>
        <div className="flex flex-wrap items-center justify-between gap-2">
          <CardTitle>{metricLabel} latency (ms)</CardTitle>
          <Select
            value={metric}
            onValueChange={(v) => setMetric(v as LatencyMetric)}
          >
            <SelectTrigger className="w-28" size="sm">
              <SelectValue />
            </SelectTrigger>
            <SelectContent align="end">
              {latencyMetricOptions.map((opt) => (
                <SelectItem key={opt.value} value={opt.value}>
                  {opt.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      </CardHeader>
      <CardContent>
        <ChartContainer config={chartConfig} className="w-full">
          <AreaChart accessibilityLayer data={chartData}>
            <CartesianGrid vertical={false} />
            <XAxis
              dataKey="timestamp"
              tickLine={false}
              tickMargin={10}
              axisLine={false}
              tickFormatter={xAxisFormat}
            />
            <ChartTooltip
              content={
                <ChartTooltipContent
                  labelFormatter={tooltipLabelFormat}
                  className="gap-1"
                  formatter={(value, _name, item) => {
                    const regionID = String(item.dataKey ?? item.name ?? "");
                    const region = regionByID.get(regionID);
                    const city = region?.city ?? `Region ${regionID}`;
                    const country = region?.country ?? "";
                    const flag = region?.flag ?? "un";
                    const indicatorColor =
                      (item.payload as { fill?: string } | undefined)?.fill ??
                      item.color ??
                      "var(--muted)";

                    return (
                      <div className="flex w-full items-center justify-between gap-2">
                        <div className="flex min-w-0 items-center gap-2">
                          <div
                            className="h-2.5 w-2.5 shrink-0 rounded-[2px]"
                            style={{ backgroundColor: indicatorColor }}
                          />
                          <div className="grid min-w-0 gap-1 leading-tight">
                            <span className="truncate">{city}</span>
                            {country ? (
                              <span className="flex items-center gap-1 truncate text-muted-foreground">
                                <span className={`fi fi-${flag}`}></span>
                                <span className="truncate">{country}</span>
                              </span>
                            ) : null}
                          </div>
                        </div>
                        <span className="text-foreground font-mono font-medium tabular-nums">
                          {formatLatencyMs(
                            typeof value === "number" ? value : Number(value),
                          )}
                        </span>
                      </div>
                    );
                  }}
                />
              }
            />
            <ChartLegend content={<ChartLegendContent />} />
            {regionIDs.map((regionID) => (
              <Area
                key={regionID}
                dataKey={regionID}
                type="monotone"
                stroke={`var(--color-${regionID})`}
                fill={`var(--color-${regionID})`}
                fillOpacity={0.35}
              />
            ))}
          </AreaChart>
        </ChartContainer>
      </CardContent>
    </Card>
  );
}
