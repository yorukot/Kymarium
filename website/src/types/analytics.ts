// ============================================================================
// Monitor Analytics Types
// ============================================================================

import type { Monitor } from './monitor';
import type { Incident } from './incident';

export interface AnalyticsWindow {
	start: string;
	end: string;
	bucket: string;
}

export interface AnalyticsSummary {
	totalCount: number;
	goodCount: number;
	uptimePct: number;
	p50Ms: number;
	p75Ms: number;
	p90Ms: number;
	p95Ms: number;
	p99Ms: number;
}

export interface AnalyticsRegionSummary extends AnalyticsSummary {
	regionId: string;
}

export interface AnalyticsSeriesPoint {
	timestamp: string;
	regionId: string;
	totalCount: number;
	goodCount: number;
	uptimePct: number;
	p50Ms: number;
	p75Ms: number;
	p90Ms: number;
	p95Ms: number;
	p99Ms: number;
}

export interface MonitorAnalytics {
	monitor: Monitor;
	window: AnalyticsWindow;
	summary: AnalyticsSummary;
	regions: AnalyticsRegionSummary[];
	series: AnalyticsSeriesPoint[];
	incidents: Incident[];
}
