import { getTeamMembers } from '$lib/api/team';
import type { TeamInvite, TeamMemberWithUser } from '../$lib/types';
import type { PageLoad } from './$types';

export type MembersPageData = {
	members: TeamMemberWithUser[];
	invites: TeamInvite[];
};

export const load: PageLoad<MembersPageData> = async ({ params }) => {
	const { teamID } = params;

	const membersResponse = await getTeamMembers(teamID);

	return {
		members: membersResponse.data,
		invites: []
	};
};
