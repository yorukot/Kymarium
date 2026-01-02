import { getIncident, getIncidentEvents } from '$lib/api/incident';
import { getMonitors } from '$lib/api/monitor';
import type { Incident, IncidentEvent, MonitorWithIncidents } from '$lib/types';
import type { PageLoad } from './$types';

export type IncidentDetailPageData = {
	incident: Incident;
	events: IncidentEvent[];
	monitorNames: string[];
};

export const load: PageLoad<IncidentDetailPageData> = async ({ params }) => {
	const { teamID, incidentID } = params;

	const [incidentRes, eventsRes, monitorsRes] = await Promise.all([
		getIncident(teamID, incidentID),
		getIncidentEvents(teamID, incidentID),
		getMonitors(teamID)
	]);

	const monitors = monitorsRes.data ?? [];
	const monitorNamesByIncident = buildMonitorLookup(monitors);

	return {
		incident: incidentRes.data,
		events: eventsRes.data ?? [],
		monitorNames: Array.from(monitorNamesByIncident[incidentID] ?? [])
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
