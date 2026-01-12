import type {
  MonitorAnalyticsRawData,
  MonitorAnalyticsSeriesPoint,
} from "@/lib/schemas/monitor-analytics";

function parseBucketToMs(bucket: string): number | null {
  const trimmed = bucket.trim().toLowerCase();
  const match = /^(\d+)(min|h)$/.exec(trimmed);
  if (!match) return null;

  const value = Number(match[1]);
  if (!Number.isFinite(value) || value <= 0) return null;

  const unit = match[2];
  switch (unit) {
    case "min":
      return value * 60_000;
    case "h":
      return value * 60 * 60_000;
    default:
      return null;
  }
}

function floorToBucket(ms: number, bucketMs: number): number {
  return Math.floor(ms / bucketMs) * bucketMs;
}

function ceilToBucket(ms: number, bucketMs: number): number {
  const floored = floorToBucket(ms, bucketMs);
  return ms === floored ? floored : floored + bucketMs;
}

function uniqueStrings(values: string[]): string[] {
  return Array.from(new Set(values));
}

function toSeriesKey(timestampMs: number, regionID: string): string {
  return `${timestampMs}|${regionID}`;
}

function parseTimestampMs(value: string): number | null {
  const ms = new Date(value).getTime();
  return Number.isNaN(ms) ? null : ms;
}

function toBucketTimestamp(timestampMs: number): string {
  return new Date(timestampMs).toISOString();
}

export function fillMonitorAnalyticsSeries(
  analytics: MonitorAnalyticsRawData,
): MonitorAnalyticsSeriesPoint[] {
  const series = analytics.series ?? [];
  const window = analytics.window;

  if (!window?.start || !window?.end || !window?.bucket) {
    return series;
  }

  const bucketMs = parseBucketToMs(window.bucket);
  if (!bucketMs) {
    return series;
  }

  const startMsRaw = parseTimestampMs(window.start);
  const endMsRaw = parseTimestampMs(window.end);
  if (startMsRaw === null || endMsRaw === null) {
    return series;
  }

  if (startMsRaw >= endMsRaw) {
    return series;
  }

  const regionIDs =
    analytics.monitor?.regions?.length
      ? uniqueStrings(analytics.monitor.regions)
      : uniqueStrings(series.map((p) => p.region_id));

  if (regionIDs.length === 0) {
    return series;
  }

  // Backend query uses `bucket >= start AND bucket < end`, so we must start from
  // the first bucket boundary at/after `start`, not the floor bucket.
  const startBucketMs = ceilToBucket(startMsRaw, bucketMs);
  const endExclusiveMs = endMsRaw;

  const existing = new Map<string, MonitorAnalyticsSeriesPoint>();
  for (const point of series) {
    const tsMs = parseTimestampMs(point.timestamp);
    if (tsMs === null) continue;
    const bucketStartMs = floorToBucket(tsMs, bucketMs);
    existing.set(toSeriesKey(bucketStartMs, point.region_id), {
      ...point,
      // Normalize timestamps so we don't create duplicate rows due to different
      // RFC3339 representations of the same instant (e.g. Z vs +08:00).
      timestamp: toBucketTimestamp(bucketStartMs),
    });
  }

  const filled: MonitorAnalyticsSeriesPoint[] = [];
  for (let t = startBucketMs; t < endExclusiveMs; t += bucketMs) {
    const timestamp = toBucketTimestamp(t);

    for (const regionID of regionIDs) {
      const key = toSeriesKey(t, regionID);
      const found = existing.get(key);

      if (found) {
        filled.push(found);
        continue;
      }

      filled.push({
        timestamp,
        region_id: regionID,
        total_count: 0,
        good_count: 0,
        uptime_pct: 0,
        p50_ms: 0,
        p75_ms: 0,
        p90_ms: 0,
        p95_ms: 0,
        p99_ms: 0,
      });
    }
  }

  filled.sort((a, b) => {
    const ta = parseTimestampMs(a.timestamp) ?? 0;
    const tb = parseTimestampMs(b.timestamp) ?? 0;
    if (ta !== tb) return ta - tb;
    return a.region_id.localeCompare(b.region_id);
  });

  return filled;
}

export type MonitorAnalyticsMetric =
  | "total_count"
  | "good_count"
  | "uptime_pct"
  | "p50_ms"
  | "p75_ms"
  | "p90_ms"
  | "p95_ms"
  | "p99_ms";

export function toStackedBarChartData(
  analytics: MonitorAnalyticsRawData,
  metric: MonitorAnalyticsMetric = "total_count",
): Array<Record<string, number | string>> {
  const filled = fillMonitorAnalyticsSeries(analytics);

  const regionIDs =
    analytics.monitor?.regions?.length
      ? uniqueStrings(analytics.monitor.regions)
      : uniqueStrings(filled.map((p) => p.region_id));

  const byTimestamp = new Map<string, Record<string, number | string>>();
  for (const point of filled) {
    const tsMs = parseTimestampMs(point.timestamp);
    if (tsMs === null) continue;
    const key = toBucketTimestamp(tsMs);
    const row = byTimestamp.get(key) ?? { timestamp: key };
    row[point.region_id] = point[metric];
    byTimestamp.set(key, row);
  }

  const rows = Array.from(byTimestamp.values());
  rows.sort((a, b) => {
    const ta = parseTimestampMs(String(a.timestamp)) ?? 0;
    const tb = parseTimestampMs(String(b.timestamp)) ?? 0;
    return ta - tb;
  });

  // Ensure stable keys even if a region has no samples in the window.
  for (const row of rows) {
    for (const regionID of regionIDs) {
      if (typeof row[regionID] !== "number") row[regionID] = 0;
    }
  }

  return rows;
}
