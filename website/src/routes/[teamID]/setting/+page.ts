import { getTeam } from '$lib/api/team';
import type { TeamWithRole } from '../../../types';
import type { PageLoad } from './$types';

export type TeamSettingPageData = {
	team: TeamWithRole;
};

export const load: PageLoad<TeamSettingPageData> = async ({ params }) => {
	const { teamID } = params;

	const teamResponse = await getTeam(teamID);

	return {
		team: teamResponse.data
	};
};
