import type { PublicIncident } from './incident';

// ============================================================================
// Status Page Types
// ============================================================================

export type StatusPageElementType = 'historical_timeline' | 'current_status_indicator';

export interface StatusPage {
	id: string;
	teamId: string;
	title: string;
	slug: string;
	icon?: string; // base64 or URL depending on API usage
	createdAt: string;
	updatedAt: string;
}

export interface StatusPageGroup {
	id: string;
	statusPageId: string;
	name: string;
	type: StatusPageElementType;
	sortOrder: number;
}

export interface StatusPageMonitor {
	id: string;
	statusPageId: string;
	monitorId: string;
	groupId?: string | null;
	name: string;
	type: StatusPageElementType;
	sortOrder: number;
}

export interface StatusPageElement {
	id: string;
	statusPageId: string;
	name: string;
	type: StatusPageElementType;
	sortOrder: number;
	monitor: boolean;
	monitorId?: string | null;
	monitors: StatusPageMonitor[];
}

export interface StatusPageWithElements {
	statusPage: StatusPage;
	elements: StatusPageElement[];
}

export interface PublicTimelinePoint {
	day: string;
	success: number;
	fail: number;
}

export interface PublicStatusPageGroup {
	id: string;
	name: string;
	type: StatusPageElementType;
	sortOrder: number;
	status?: 'up' | 'down';
	uptimeSli30?: number;
	uptimeSli60?: number;
	uptimeSli90?: number;
	timeline?: PublicTimelinePoint[];
}

export interface PublicStatusPageMonitor {
	id: string;
	monitorId: string;
	groupId?: string | null;
	name: string;
	type: StatusPageElementType;
	sortOrder: number;
	status?: 'up' | 'down';
	uptimeSli30?: number;
	uptimeSli60?: number;
	uptimeSli90?: number;
	timeline?: PublicTimelinePoint[];
}

export interface PublicStatusPageElement {
	id: string;
	name: string;
	type: StatusPageElementType;
	sortOrder: number;
	status?: 'up' | 'down';
	monitor: boolean;
	monitorId?: string | null;
	uptimeSli30?: number;
	uptimeSli60?: number;
	uptimeSli90?: number;
	timeline?: PublicTimelinePoint[];
	monitors: PublicStatusPageMonitor[];
}

export interface PublicStatusPageData {
	statusPage: StatusPage;
	elements: PublicStatusPageElement[];
	incidents: PublicIncident[];
}
