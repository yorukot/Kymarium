import type {
  IncidentDetailItem,
  IncidentEventItem,
  IncidentEventRawData,
  IncidentListItem,
  IncidentRawData,
} from "@/lib/schemas/incident";

export function parseIncident(raw: IncidentRawData): IncidentListItem {
  return {
    id: raw.id,
    title: raw.title ?? undefined,
    status: raw.status,
    severity: raw.severity,
    isPublic: raw.is_public,
    startedAt: raw.started_at,
    resolvedAt: raw.resolved_at ?? undefined,
  };
}

export function parseIncidents(rawList: IncidentRawData[]): IncidentListItem[] {
  return rawList.map(parseIncident);
}

export function parseIncidentDetail(raw: IncidentRawData): IncidentDetailItem {
  return {
    id: raw.id,
    title: raw.title ?? undefined,
    status: raw.status,
    severity: raw.severity,
    isPublic: raw.is_public,
    autoResolve: raw.auto_resolve,
    startedAt: raw.started_at,
    resolvedAt: raw.resolved_at ?? undefined,
    createdAt: raw.created_at,
    updatedAt: raw.updated_at,
  };
}

export function parseIncidentEvent(raw: IncidentEventRawData): IncidentEventItem {
  return {
    id: raw.id,
    incidentId: raw.incident_id,
    message: raw.message,
    eventType: raw.event_type,
    isPublic: raw.is_public,
    createdAt: raw.created_at,
    updatedAt: raw.updated_at,
  };
}

export function parseIncidentEvents(
  rawList: IncidentEventRawData[],
): IncidentEventItem[] {
  return rawList.map(parseIncidentEvent);
}
