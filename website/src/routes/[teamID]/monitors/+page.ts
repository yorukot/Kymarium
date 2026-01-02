import { getMonitors } from '$lib/api/monitor';
import type { PageLoad } from './$types';
import type { MonitorWithIncidents } from '../$lib/types';

export type MonitorsData = {
  monitors: MonitorWithIncidents[];
};

export const load: PageLoad<MonitorsData> = async ({ params }) => {
  const { teamID } = params;
  const response = await getMonitors(teamID);

  return {
    monitors: response.data
  };
};
