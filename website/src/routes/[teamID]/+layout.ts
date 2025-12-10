import { error, redirect } from '@sveltejs/kit';
import { refreshAccessToken } from '$lib/api/auth';
import { getTeamById } from '$lib/api/team';
import { getCurrentUser } from '$lib/api/user';

export const ssr = false;

const ensureAuthenticatedUser = async (url: URL) => {

	let { response, body } = await getCurrentUser();

	if (response.status === 401) {
		const refreshResult = await refreshAccessToken();
		if (!refreshResult.response.ok) {
			const next = encodeURIComponent(url.pathname + url.search);
			redirect(302, `/auth/login?next=${next}`);
		}

		({ response, body } = await getCurrentUser());
	}

	if (!response.ok || !body.data) {
		throw error(response.status || 500, body.message ?? 'Failed to fetch current user');
	}

	return body.data;
};

export const load: import('./$types').LayoutLoad = async ({ params, url }) => {
	const user = await ensureAuthenticatedUser(url);

	const { response: teamResponse, body: teamBody } = await getTeamById(params.teamID);

	if (teamResponse.status === 401) {
		const next = encodeURIComponent(url.pathname + url.search);
		redirect(302, `/auth/login?next=${next}`);
	}

	if (teamResponse.status === 404) {
		throw error(404, 'Team not found');
	}

	if (!teamResponse.ok || !teamBody.data) {
		throw error(teamResponse.status || 500, teamBody.message ?? 'Failed to fetch team');
	}

	return {
		user: {
			id: user.id,
			name: user.display_name,
			avatar: user.avatar ?? undefined
		},
		team: {
			id: teamBody.data.id,
			name: teamBody.data.name,
			role: teamBody.data.role
		},
		navItems: [
			{ title: 'Monitors', url: `/${params.teamID}/monitors`, icon: 'lucide:monitor' },
			{ title: 'Incidents', url: `/${params.teamID}/incidents`, icon: 'lucide:alert-triangle' },
			{ title: 'Notifications', url: `/${params.teamID}/notifications`, icon: 'lucide:bell' }
		]
	};
};
