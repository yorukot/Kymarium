import { apiFetch, type ApiResponse } from './client';

type User = {
	id: string;
	display_name: string;
	avatar?: string | null;
};

export type UserResponse = ApiResponse<User>;

export const getCurrentUser = async () => {
	const { response, body } = await apiFetch<User>('/api/users/me');
	return { response, body };
};
