import { getIncidents } from '$lib/api/incident';
import { getMonitors } from '$lib/api/monitor';
import type { Incident, MonitorWithIncidents } from '../$lib/types';
import type { PageLoad } from './$types';

export type IncidentWithMonitors = Incident & { monitorNames: string[] };

export type IncidentsPageData = {
  incidents: IncidentWithMonitors[];
};

export const load: PageLoad<IncidentsPageData> = async ({ params }) => {
  const { teamID } = params;

  const [incidentsRes, monitorsRes] = await Promise.all([
    getIncidents(teamID),
    getMonitors(teamID)
  ]);

  const monitors = monitorsRes.data ?? [];
  const monitorNamesByIncident: Record<string, Set<string>> = buildMonitorLookup(monitors);

  const incidents = (incidentsRes.data ?? []).map((incident) => ({
    ...incident,
    monitorNames: Array.from(monitorNamesByIncident[incident.id] ?? [])
  }));

  return {
    incidents
  };
};

function buildMonitorLookup(monitors: MonitorWithIncidents[]): Record<string, Set<string>> {
  return monitors.reduce<Record<string, Set<string>>>((acc, monitor) => {
    monitor.incidents?.forEach((incident) => {
      const key = incident.id;
      if (!acc[key]) acc[key] = new Set<string>();
      acc[key].add(monitor.name);
    });
    return acc;
  }, {});
}
