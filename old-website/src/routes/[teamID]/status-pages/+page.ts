import { getStatusPages } from '$lib/api/status-page';
import type { StatusPageWithElements } from '../$lib/types';
import type { PageLoad } from './$types';

export type StatusPagesData = {
	statusPages: StatusPageWithElements[];
};

export const load: PageLoad<StatusPagesData> = async ({ params }) => {
	const { teamID } = params;
	const response = await getStatusPages(teamID);

	return {
		statusPages: response.data
	};
};
