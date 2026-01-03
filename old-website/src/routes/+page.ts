import { redirect } from '@sveltejs/kit';
import { getTeams } from '$lib/api/team';
import type { PageLoad } from './$types';

export const ssr = false;

export const load: PageLoad = async () => {
	const response = await getTeams();
	const teams = response.data ?? [];

	if (!teams.length) {
		redirect(307, '/new-team');
	}

	redirect(307, `/${teams[0].id}`);
};
