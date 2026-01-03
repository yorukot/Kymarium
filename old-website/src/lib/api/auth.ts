import { PUBLIC_API_BASE } from '$env/static/public';
import { apiRequest, type ApiDefaultBody } from './utils';
import { publicApiRequest } from './utils';

let refreshPromise: Promise<boolean> | null = null;

export async function refreshToken(): Promise<boolean> {
	if (!refreshPromise) {
		refreshPromise = fetch(`${PUBLIC_API_BASE}/auth/refresh`, {
			method: 'POST',
			credentials: 'include'
		})
			.then((res) => res.ok)
			.finally(() => {
				refreshPromise = null;
			});
	}

	return refreshPromise;
}

export function buildOAuthUrl(provider: string, next: string = '/'): string {
	const params = new URLSearchParams();
	if (next) {
		params.set('next', next);
	}

	const queryString = params.toString();
	return `${PUBLIC_API_BASE}/auth/oauth/${provider}${queryString ? `?${queryString}` : ''}`;
}

export async function login(email: string, password: string): Promise<ApiDefaultBody> {
	return apiRequest<ApiDefaultBody>('/auth/login', {
		method: 'POST',
		body: { email, password },
		defaultError: 'Login failed'
	});
}

export async function registerUser(
	displayName: string,
	email: string,
	password: string
): Promise<ApiDefaultBody> {
	return apiRequest<ApiDefaultBody>('/auth/register', {
		method: 'POST',
		body: { displayName, email, password },
		defaultError: 'Registration failed'
	});
}

export async function logout(): Promise<ApiDefaultBody> {
	return apiRequest<ApiDefaultBody>('/auth/logout', {
		method: 'POST',
		defaultError: 'Logout failed'
	});
}

export async function resendVerification(email: string): Promise<ApiDefaultBody> {
	return publicApiRequest<ApiDefaultBody>('/auth/verify/resend', {
		method: 'POST',
		body: { email },
		defaultError: 'Resend verification failed'
	});
}
