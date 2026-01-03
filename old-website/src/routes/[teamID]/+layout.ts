import { getTeams } from '$lib/api/team';
import { getUser } from '$lib/api/user';
import type { Team, User } from '$lib/types';
import type { LayoutLoad } from './$types';

export type SidebarData = {
  teams: Team[];
  user: User;
};

export const load: LayoutLoad<SidebarData> = async () => {
  const [teamsResponse, userResponse] = await Promise.all([
    getTeams(),
    getUser()
  ]);

  return {
    teams: teamsResponse.data,
    user: userResponse.data
  };
};

export const ssr = false;