// ============================================================================
// Monitor Types
// ============================================================================

import type { HTTPMonitorConfig, PingMonitorConfig } from './monitor-config';
import type { Incident } from './incident';

export type MonitorType = 'http' | 'ping';
export type MonitorStatus = 'up' | 'down';

export interface Monitor {
	id: string;
	teamId: string;
	name: string;
	type: MonitorType;
	config: HTTPMonitorConfig | PingMonitorConfig | Record<string, unknown>;
	interval: number;
	status: MonitorStatus;
	uptimeSli30?: number;
	lastChecked: string;
	nextCheck: string;
	failureThreshold: number;
	recoveryThreshold: number;
	regions: string[];
	notification: string[];
	updatedAt: string;
	createdAt: string;
}

export interface MonitorWithIncidents extends Monitor {
	incidents?: Incident[];
}

export interface MonitorNotification {
	id: string;
	monitorId: string;
	notificationId: string;
}
