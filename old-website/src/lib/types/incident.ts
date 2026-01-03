// ============================================================================
// Incident Types
// ============================================================================

export type IncidentStatus =
	| 'detected'
	| 'investigating'
	| 'identified'
	| 'monitoring'
	| 'resolved';

export type IncidentSeverity = 'emergency' | 'critical' | 'major' | 'minor' | 'info';

export type IncidentEventType =
	| 'detected'
	| 'notification_sent'
	| 'manually_resolved'
	| 'auto_resolved'
	| 'unpublished'
	| 'published'
	| 'investigating'
	| 'identified'
	| 'update'
	| 'monitoring';

export interface Incident {
	id: string;
	title?: string | null;
	status: IncidentStatus;
	severity: IncidentSeverity;
	isPublic: boolean;
	autoResolve: boolean;
	startedAt: string;
	resolvedAt?: string;
	createdAt: string;
	updatedAt: string;
}

export interface PublicIncident extends Incident {
	monitorId: string;
	timeline?: IncidentEvent[];
}

export interface IncidentMonitor {
	id: string;
	incidentId: string;
	monitorId: string;
}

export interface IncidentEvent {
	id: string;
	incidentId: string;
	createdBy?: string | null;
	message: string;
	eventType: IncidentEventType;
	createdAt: string;
	updatedAt: string;
}
