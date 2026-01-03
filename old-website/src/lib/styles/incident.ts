import type { Incident, IncidentEventType } from '$lib/types';

export const incidentSeverityMeta: Record<
	Incident['severity'],
	{ label: string; badge: string; dot: string }
> = {
	emergency: {
		label: 'Emergency',
		badge: '!bg-destructive !text-destructive-foreground border-transparent',
		dot: 'bg-destructive'
	},
	critical: {
		label: 'Critical',
		badge: '!bg-destructive/80 !text-destructive-foreground border-transparent',
		dot: 'bg-destructive/80'
	},
	major: {
		label: 'Major',
		badge: '!bg-foreground/10 !text-foreground border-transparent',
		dot: 'bg-foreground/50'
	},
	minor: {
		label: 'Minor',
		badge: '!bg-secondary !text-secondary-foreground border-transparent',
		dot: 'bg-secondary'
	},
	info: {
		label: 'Info',
		badge: '!bg-muted !text-muted-foreground border-transparent',
		dot: 'bg-muted'
	}
};

export const incidentStatusMeta: Record<Incident['status'], { label: string; tone: string }> = {
	detected: { label: 'Detected', tone: 'text-destructive' },
	investigating: { label: 'Investigating', tone: 'text-destructive' },
	identified: { label: 'Identified', tone: 'text-destructive' },
	monitoring: { label: 'Monitoring', tone: 'text-foreground' },
	resolved: { label: 'Resolved', tone: 'text-success' }
};

export const incidentEventTypeMeta: Record<IncidentEventType, { label: string; badge: string }> = {
	detected: { label: 'Detected', badge: '!bg-destructive/10 !text-destructive border-transparent' },
	investigating: {
		label: 'Investigating',
		badge: '!bg-destructive/10 !text-destructive border-transparent'
	},
	identified: {
		label: 'Identified',
		badge: '!bg-destructive/10 !text-destructive border-transparent'
	},
	monitoring: {
		label: 'Monitoring',
		badge: '!bg-secondary !text-secondary-foreground border-transparent'
	},
	update: { label: 'Update', badge: '!bg-muted !text-foreground border-transparent' },
	notification_sent: {
		label: 'Notification sent',
		badge: '!bg-muted !text-foreground border-transparent'
	},
	published: { label: 'Published', badge: '!bg-success/10 !text-success border-transparent' },
	unpublished: {
		label: 'Unpublished',
		badge: '!bg-muted !text-muted-foreground border-transparent'
	},
	auto_resolved: {
		label: 'Auto resolved',
		badge: '!bg-success/10 !text-success border-transparent'
	},
	manually_resolved: {
		label: 'Manually resolved',
		badge: '!bg-success/10 !text-success border-transparent'
	}
};

export const severityLabel = (severity: Incident['severity']) =>
	incidentSeverityMeta[severity]?.label ?? severity;
export const severityBadgeClass = (severity: Incident['severity']) =>
	incidentSeverityMeta[severity]?.badge ?? '!bg-muted !text-foreground border-transparent';
export const severityDotClass = (severity: Incident['severity']) =>
	incidentSeverityMeta[severity]?.dot ?? 'bg-muted';

export const statusLabel = (status: Incident['status']) =>
	incidentStatusMeta[status]?.label ?? status;
export const statusTone = (status: Incident['status']) =>
	incidentStatusMeta[status]?.tone ?? 'text-foreground';

export const eventTypeLabel = (eventType: IncidentEventType) =>
	incidentEventTypeMeta[eventType]?.label ?? eventType;
export const eventTypeBadgeClass = (eventType: IncidentEventType) =>
	incidentEventTypeMeta[eventType]?.badge ?? '!bg-muted !text-foreground border-transparent';

export function formatIncidentDate(value?: string) {
	if (!value) return '—';
	const date = new Date(value);
	return Number.isNaN(date.getTime()) ? '—' : date.toLocaleString();
}

export function formatIncidentDuration(target: Pick<Incident, 'startedAt' | 'resolvedAt'>) {
	const toLocalMs = (value: string | Date) => {
		const date = value instanceof Date ? value : new Date(value);
		if (Number.isNaN(date.getTime())) return Number.NaN;
		return date.getTime() - date.getTimezoneOffset() * 60_000;
	};

	const start = toLocalMs(target.startedAt);
	const end = toLocalMs(target.resolvedAt ?? new Date());
	if (!Number.isFinite(start) || !Number.isFinite(end) || end < start) return '—';
	const diffMs = end - start;
	const minutes = Math.floor(diffMs / 60000);
	if (minutes < 60) return `${minutes}m`;
	const hours = Math.floor(minutes / 60);
	const mins = minutes % 60;
	if (hours < 24) return `${hours}h ${mins}m`;
	const days = Math.floor(hours / 24);
	return `${days}d ${hours % 24}h`;
}
