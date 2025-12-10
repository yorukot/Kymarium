import { apiFetch, buildApiUrl, type ApiResponse } from './client';

const AUTH_BASE_PATH = '/api/auth';

export const refreshAccessToken = async () =>
	apiFetch<never>(`${AUTH_BASE_PATH}/refresh`, {
		method: 'POST'
	});

export const login = async (email: string, password: string) => {
	const { body } = await apiFetch<never>(`${AUTH_BASE_PATH}/login`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ email, password })
	});

	return body;
};

export const registerUser = async (displayName: string, email: string, password: string) => {
	const { body } = await apiFetch<never>(`${AUTH_BASE_PATH}/register`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ display_name: displayName, email, password })
	});

	return body;
};

export const buildOAuthUrl = (provider: string, nextPath = '/') => {
	const params = new URLSearchParams({ next: nextPath });
	return `${buildApiUrl(`${AUTH_BASE_PATH}/oauth/${provider}`)}?${params.toString()}`;
};

export type AuthResponse<T> = ApiResponse<T>;
