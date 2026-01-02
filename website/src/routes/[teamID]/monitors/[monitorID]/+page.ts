import { getMonitorAnalytics } from '$lib/api/monitor';
import { getRegions } from '$lib/api/region';
import type { PageLoad } from './$types';

export type MonitorDetailData = {
	analytics: MonitorAnalytics;
	regions: Region[];
};

export const load: PageLoad<MonitorDetailData> = async ({ params }) => {
	const { teamID, monitorID } = params;
	const [regionRes, monitorAnalyticsRes] = await Promise.all([
		getRegions(),
		getMonitorAnalytics(teamID, monitorID)
	]);
	return {
		analytics: monitorAnalyticsRes.data,
		regions: regionRes.data
	};
};
