import type { MonitorRawData } from "@/lib/schemas/monitor";

export type MonitorAnalyticsWindow = {
  start: string;
  end: string;
  bucket: string;
};

export type MonitorAnalyticsSummary = {
  total_count: number;
  good_count: number;
  uptime_pct: number;
  p50_ms: number;
  p75_ms: number;
  p90_ms: number;
  p95_ms: number;
  p99_ms: number;
};

export type MonitorAnalyticsRegionSummary = {
  region_id: string;
  total_count: number;
  good_count: number;
  uptime_pct: number;
  p50_ms: number;
  p75_ms: number;
  p90_ms: number;
  p95_ms: number;
  p99_ms: number;
};

export type MonitorAnalyticsSeriesPoint = {
  timestamp: string;
  region_id: string;
  total_count: number;
  good_count: number;
  uptime_pct: number;
  p50_ms: number;
  p75_ms: number;
  p90_ms: number;
  p95_ms: number;
  p99_ms: number;
};

export type MonitorAnalyticsIncident = {
  id: string;
  monitor_id: string;
  title?: string | null;
  status: string;
  started_at: string;
  resolved_at?: string | null;
  created_at: string;
  updated_at: string;
};

export type MonitorAnalyticsRawData = {
  monitor: MonitorRawData;
  window?: MonitorAnalyticsWindow;
  summary?: MonitorAnalyticsSummary;
  regions?: MonitorAnalyticsRegionSummary[];
  series?: MonitorAnalyticsSeriesPoint[];
  incidents?: MonitorAnalyticsIncident[];
};
