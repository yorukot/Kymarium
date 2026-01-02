import { getMonitor } from '$lib/api/monitor';
import { getNotifications } from '$lib/api/notification';
import { getRegions } from '$lib/api/region';
import type { Monitor, Notification, Region } from '$lib/types';
import type { PageLoad } from './$types';
import { error } from '@sveltejs/kit';

export type EditMonitorData = {
	monitor: Monitor;
	notifications: Notification[];
	regions: Region[];
};

export const load: PageLoad<EditMonitorData> = async ({ params }) => {
	const { teamID, monitorID } = params;

	if (!monitorID) {
		error(400, 'monitorId query parameter is required');
	}

	const [monitorRes, notificationsRes, regionsRes] = await Promise.all([
		getMonitor(teamID, monitorID),
		getNotifications(teamID),
		getRegions()
	]);

	const monitorRaw = monitorRes.data;
	const monitor: Monitor = {
		...monitorRaw,
		teamId: monitorRaw.teamId,
		lastChecked: monitorRaw.lastChecked,
		nextCheck: monitorRaw.nextCheck,
		failureThreshold: monitorRaw.failureThreshold,
		recoveryThreshold: monitorRaw.recoveryThreshold,
		updatedAt: monitorRaw.updatedAt,
		createdAt: monitorRaw.createdAt
	};

	return {
		monitor,
		notifications: notificationsRes.data,
		regions: regionsRes.data
	};
};
