import { getMonitorAnalytics } from '$lib/api/monitor';
import type { MonitorAnalytics } from '../../../../../types';
import type { PageLoad } from './$types';

export type MonitorDetailData = {
	analytics: MonitorAnalytics;
};

export const load: PageLoad<MonitorDetailData> = async ({ params }) => {
	const { teamID, monitorID } = params;
	const response = await getMonitorAnalytics(teamID, monitorID);

	return {
		analytics: response.data
	};
};
