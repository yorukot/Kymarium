import { PUBLIC_API_BASE } from '$env/static/public';

const API_BASE = (PUBLIC_API_BASE ?? '').replace(/\/+$/, '');

export type ApiResponse<T> = {
	data?: T;
	message?: string;
};

const normalizePath = (path: string) => {
	if (!path) return '';
	return path.startsWith('/') ? path : `/${path}`;
};

export const buildApiUrl = (path: string) => `${API_BASE}${normalizePath(path)}`;

const parseResponse = async <T>(response: Response): Promise<ApiResponse<T>> => {
	try {
		return (await response.json()) as ApiResponse<T>;
	} catch {
		return {};
	}
};

export const apiFetch = async <T>(path: string, init: RequestInit = {}) => {
	const response = await fetch(buildApiUrl(path), {
		credentials: 'include',
		...init
	});

	const body = await parseResponse<T>(response);

	return { response, body };
};
