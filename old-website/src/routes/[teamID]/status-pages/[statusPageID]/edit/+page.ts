import { getMonitors } from '$lib/api/monitor';
import { getStatusPage } from '$lib/api/status-page';
import type { Monitor } from '$lib/types';
import type { StatusPageWithElements } from '$lib/types/status-page';
import type { PageLoad } from './$types';

export type StatusPageEditData = {
	statusPage: StatusPageWithElements;
	monitors: Monitor[];
};

export const load: PageLoad<StatusPageEditData> = async ({ params }) => {
	const { teamID, statusPageID } = params;

	const [statusPageResponse, monitorsResponse] = await Promise.all([
		getStatusPage(teamID, statusPageID),
		getMonitors(teamID)
	]);

	return {
		statusPage: statusPageResponse.data,
		monitors: monitorsResponse.data
	};
};
