import { apiFetch, type ApiResponse } from './client';

type Team = {
	id: string;
	name: string;
	role?: string;
};

export type TeamResponse = ApiResponse<Team>;

export const getTeamById = async (teamId: string) => {
	const { response, body } = await apiFetch<Team>(`/api/teams/${teamId}`);

	return { response, body };
};
