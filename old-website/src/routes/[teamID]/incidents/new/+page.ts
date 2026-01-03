import { getMonitors } from '$lib/api/monitor';
import type { Monitor } from '$lib/types';
import type { PageLoad } from './$types';

export type NewIncidentPageData = {
	monitors: Monitor[];
};

export const load: PageLoad<NewIncidentPageData> = async ({ params }) => {
	const { teamID } = params;
	const monitorsRes = await getMonitors(teamID);

	return {
		monitors: monitorsRes.data ?? []
	};
};
