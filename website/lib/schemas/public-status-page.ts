import type { StatusPageElementType } from "@/lib/schemas/status-page";

export type PublicTimelinePoint = {
  day: string;
  success: number;
  fail: number;
};

export type PublicStatusPageMonitorRaw = {
  id: string;
  monitor_id: string;
  group_id?: string | null;
  name: string;
  type: StatusPageElementType;
  sort_order: number;
  status?: string;
  uptime_sli_30?: number;
  uptime_sli_60?: number;
  uptime_sli_90?: number;
  timeline?: PublicTimelinePoint[];
};

export type PublicStatusPageElementRaw = {
  id: string;
  name: string;
  type: StatusPageElementType;
  sort_order: number;
  status?: string;
  monitor: boolean;
  monitor_id?: string | null;
  uptime_sli_30?: number;
  uptime_sli_60?: number;
  uptime_sli_90?: number;
  timeline?: PublicTimelinePoint[];
  monitors: PublicStatusPageMonitorRaw[];
};

export type PublicStatusPageModelRaw = {
  id: string;
  team_id: string;
  title: string;
  slug: string;
  icon?: string | null;
  created_at: string;
  updated_at: string;
};

export type PublicIncidentTimelineRaw = {
  id: string;
  incident_id: string;
  created_by?: string | null;
  message: string;
  event_type: string;
  created_at: string;
  updated_at: string;
};

export type PublicIncidentRaw = {
  id: string;
  title?: string | null;
  status: string;
  severity: string;
  is_public: boolean;
  auto_resolve: boolean;
  started_at: string;
  resolved_at?: string | null;
  created_at: string;
  updated_at: string;
  timeline: PublicIncidentTimelineRaw[];
  monitor_id: string;
};

export type PublicStatusPageResponseRaw = {
  StatusPage?: PublicStatusPageModelRaw;
  status_page?: PublicStatusPageModelRaw;
  Elements?: PublicStatusPageElementRaw[];
  elements?: PublicStatusPageElementRaw[];
  Incidents?: PublicIncidentRaw[];
  incidents?: PublicIncidentRaw[];
};

export type PublicStatusPageModel = {
  id: string;
  teamId: string;
  title: string;
  slug: string;
  icon?: string | null;
  createdAt: string;
  updatedAt: string;
};

export type PublicStatusPageMonitor = {
  id: string;
  monitorId: string;
  groupId?: string | null;
  name: string;
  type: StatusPageElementType;
  sortOrder: number;
  status?: string;
  uptimeSLI30?: number;
  uptimeSLI60?: number;
  uptimeSLI90?: number;
  timeline?: PublicTimelinePoint[];
};

export type PublicStatusPageElement = {
  id: string;
  name: string;
  type: StatusPageElementType;
  sortOrder: number;
  status?: string;
  monitor: boolean;
  monitorId?: string | null;
  uptimeSLI30?: number;
  uptimeSLI60?: number;
  uptimeSLI90?: number;
  timeline?: PublicTimelinePoint[];
  monitors: PublicStatusPageMonitor[];
};

export type PublicIncidentTimeline = {
  id: string;
  incidentId: string;
  createdBy?: string | null;
  message: string;
  eventType: string;
  createdAt: string;
  updatedAt: string;
};

export type PublicIncident = {
  id: string;
  title?: string | null;
  status: string;
  severity: string;
  isPublic: boolean;
  autoResolve: boolean;
  startedAt: string;
  resolvedAt?: string | null;
  createdAt: string;
  updatedAt: string;
  timeline: PublicIncidentTimeline[];
  monitorId: string;
};

export type PublicStatusPageData = {
  statusPage: PublicStatusPageModel;
  elements: PublicStatusPageElement[];
  incidents: PublicIncident[];
};
