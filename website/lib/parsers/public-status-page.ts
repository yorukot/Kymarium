import {
  type PublicIncident,
  type PublicIncidentRaw,
  type PublicIncidentTimeline,
  type PublicIncidentTimelineRaw,
  type PublicStatusPageData,
  type PublicStatusPageElement,
  type PublicStatusPageElementRaw,
  type PublicStatusPageModel,
  type PublicStatusPageModelRaw,
  type PublicStatusPageMonitor,
  type PublicStatusPageMonitorRaw,
  type PublicStatusPageResponseRaw,
} from "@/lib/schemas/public-status-page";

function toStatusPageModel(raw: PublicStatusPageModelRaw): PublicStatusPageModel {
  return {
    id: raw.id,
    teamId: raw.team_id,
    title: raw.title,
    slug: raw.slug,
    icon: raw.icon ?? null,
    createdAt: raw.created_at,
    updatedAt: raw.updated_at,
  };
}

function toMonitor(raw: PublicStatusPageMonitorRaw): PublicStatusPageMonitor {
  return {
    id: raw.id,
    monitorId: raw.monitor_id,
    groupId: raw.group_id ?? null,
    name: raw.name,
    type: raw.type,
    sortOrder: raw.sort_order ?? 0,
    status: raw.status,
    uptimeSLI30: raw.uptime_sli_30,
    uptimeSLI60: raw.uptime_sli_60,
    uptimeSLI90: raw.uptime_sli_90,
    timeline: raw.timeline ?? [],
  };
}

function toElement(raw: PublicStatusPageElementRaw): PublicStatusPageElement {
  const monitors = (raw.monitors ?? []).map(toMonitor).sort(sortBySortOrder);

  return {
    id: raw.id,
    name: raw.name,
    type: raw.type,
    sortOrder: raw.sort_order ?? 0,
    status: raw.status,
    monitor: raw.monitor,
    monitorId: raw.monitor_id ?? null,
    uptimeSLI30: raw.uptime_sli_30,
    uptimeSLI60: raw.uptime_sli_60,
    uptimeSLI90: raw.uptime_sli_90,
    timeline: raw.timeline ?? [],
    monitors,
  };
}

function toIncidentTimeline(
  raw: PublicIncidentTimelineRaw,
): PublicIncidentTimeline {
  return {
    id: raw.id,
    incidentId: raw.incident_id,
    createdBy: raw.created_by ?? null,
    message: raw.message,
    eventType: raw.event_type,
    createdAt: raw.created_at,
    updatedAt: raw.updated_at,
  };
}

function toIncident(raw: PublicIncidentRaw): PublicIncident {
  const monitorIds = Array.isArray(raw.monitor_id)
    ? raw.monitor_id.filter(Boolean)
    : raw.monitor_id
      ? [raw.monitor_id]
      : [];

  return {
    id: raw.id,
    title: raw.title ?? null,
    status: raw.status,
    severity: raw.severity,
    isPublic: raw.is_public,
    autoResolve: raw.auto_resolve,
    startedAt: raw.started_at,
    resolvedAt: raw.resolved_at ?? null,
    createdAt: raw.created_at,
    updatedAt: raw.updated_at,
    monitorIds,
    timeline: (raw.timeline ?? []).map(toIncidentTimeline),
    // Filled later when status page context is known
    statusPageSlug: undefined,
  };
}

function sortBySortOrder<T extends { sortOrder: number }>(a: T, b: T) {
  return a.sortOrder - b.sortOrder;
}

function sortIncidentsByStart(a: PublicIncident, b: PublicIncident) {
  const aTime = new Date(a.startedAt).getTime();
  const bTime = new Date(b.startedAt).getTime();
  return Number.isNaN(bTime) || Number.isNaN(aTime) ? 0 : bTime - aTime;
}

export function parsePublicStatusPage(
  raw?: PublicStatusPageResponseRaw | null,
): PublicStatusPageData | null {
  if (!raw) return null;

  const statusPageRaw = raw.StatusPage ?? raw.status_page;
  if (!statusPageRaw) return null;

  const elementRawList = raw.Elements ?? raw.elements ?? [];
  const incidentRawList = raw.Incidents ?? raw.incidents ?? [];

  const statusPage = toStatusPageModel(statusPageRaw);
  const elements = elementRawList.map(toElement).sort(sortBySortOrder);
  const incidents = incidentRawList
    .map(toIncident)
    .map((incident) => ({ ...incident, statusPageSlug: statusPage.slug }))
    .sort(sortIncidentsByStart);

  return { statusPage, elements, incidents };
}
