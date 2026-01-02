import { getNotifications } from '$lib/api/notification';
import { getRegions } from '$lib/api/region';
import type { PageLoad } from './$types';
import type { Notification, Region } from '$lib/types';

export type NewMonitorPageData = {
	notifications: Notification[];
	regions: Region[];
};

export const load: PageLoad<NewMonitorPageData> = async ({ params }) => {
	const { teamID } = params;
	const [notificationsRes, regionsRes] = await Promise.all([
		getNotifications(teamID),
		getRegions()
	]);

	return {
		notifications: notificationsRes.data,
		regions: regionsRes.data
	};
};
