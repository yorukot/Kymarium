import { getNotifications } from '$lib/api/notification';
import type { Notification } from '../../../types';
import type { PageLoad } from './$types';

export type NotificationsPageData = {
	notifications: Notification[];
};

export const load: PageLoad<NotificationsPageData> = async ({ params }) => {
	const { teamID } = params;
	const response = await getNotifications(teamID);

	return {
		notifications: response.data
	};
};
