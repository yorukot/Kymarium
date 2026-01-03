import { PUBLIC_API_BASE } from '$env/static/public';
import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import camelcaseKeys from 'camelcase-keys';
import snakecaseKeys from 'snakecase-keys';
import { refreshToken } from './auth';

type HttpMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';

type ApiErrorBody = {
	message?: string;
};

export type ApiDefaultBody = {
	message: string;
};

type ApiOptions = {
	method?: HttpMethod;
	body?: unknown;
	defaultError?: string;
	headers?: HeadersInit;
};

let redirectingToLogin = false;

function isAuthExpiredResponse(res: Response, errorBody: ApiErrorBody | null): boolean {
	if (res.status === 401) return true;
	const message = typeof errorBody?.message === 'string' ? errorBody.message : '';
	return message.toLowerCase().includes('refresh token not found');
}

async function redirectToLogin() {
	if (!browser || redirectingToLogin) return;

	redirectingToLogin = true;
	const next = window.location.pathname + window.location.search + window.location.hash;
	await goto(`/auth/login?next=${encodeURIComponent(next)}`, { replaceState: true });
}

async function parseJson(res: Response): Promise<unknown> {
	try {
		return await res.json();
	} catch {
		return null;
	}
}

export function isAuthExpiredError(err: unknown): boolean {
	return err instanceof Error && err.message === 'AUTH_EXPIRED';
}

function normalizeResponse<T>(data: unknown): T {
	if (data && typeof data === 'object') {
		return camelcaseKeys(data as object, { deep: true }) as T;
	}
	return data as T;
}

function normalizeRequestBody(value: unknown): unknown {
	if (!value || typeof value !== 'object') return value;

	if (
		(typeof FormData !== 'undefined' && value instanceof FormData) ||
		(typeof URLSearchParams !== 'undefined' && value instanceof URLSearchParams) ||
		(typeof Blob !== 'undefined' && value instanceof Blob) ||
		(typeof File !== 'undefined' && value instanceof File) ||
		value instanceof Date
	) {
		return value;
	}

	return snakecaseKeys(value as Record<string, unknown>, { deep: true });
}

export async function apiRequest<T>(url: string, options: ApiOptions = {}): Promise<T> {
	const { method = 'GET', body, defaultError = 'Request failed', headers } = options;
	const shouldNormalizeBody = body && (method === 'POST' || method === 'PUT' || method === 'PATCH');
	const normalizedBody = shouldNormalizeBody ? normalizeRequestBody(body) : body;

	const requestInit: RequestInit = {
		method,
		credentials: 'include',
		headers: {
			...(normalizedBody ? { 'Content-Type': 'application/json' } : {}),
			...headers
		},
		body: normalizedBody ? JSON.stringify(normalizedBody) : undefined
	};

	const res = await fetch(`${PUBLIC_API_BASE}${url}`, requestInit);

	const responseBody = await parseJson(res);

	/* ---------- auth retry ---------- */

	if (res.status === 401) {
		const refreshed = await refreshToken();
		if (!refreshed) {
			await redirectToLogin();
			throw new Error('AUTH_EXPIRED');
		}

		const retryRes = await fetch(`${PUBLIC_API_BASE}${url}`, requestInit);

		if (!retryRes.ok) {
			await redirectToLogin();
			throw new Error('AUTH_EXPIRED');
		}

		const retryBody = await parseJson(retryRes);
		return normalizeResponse<T>(retryBody);
	}

	/* ---------- error handling ---------- */

	if (!res.ok) {
		const errorBody = responseBody as ApiErrorBody | null;
		if (isAuthExpiredResponse(res, errorBody)) {
			await redirectToLogin();
			throw new Error('AUTH_EXPIRED');
		}
		const message =
			typeof errorBody?.message === 'string'
				? `${defaultError}: ${errorBody.message}`
				: defaultError;

		throw new Error(message);
	}

	/* ---------- success ---------- */

	return normalizeResponse<T>(responseBody);
}

export async function publicApiRequest<T>(url: string, options: ApiOptions = {}): Promise<T> {
	const { method = 'GET', body, defaultError = 'Request failed', headers } = options;
	const shouldNormalizeBody = body && (method === 'POST' || method === 'PUT' || method === 'PATCH');
	const normalizedBody = shouldNormalizeBody ? normalizeRequestBody(body) : body;

	const requestInit: RequestInit = {
		method,
		credentials: 'omit',
		headers: {
			...(normalizedBody ? { 'Content-Type': 'application/json' } : {}),
			...headers
		},
		body: normalizedBody ? JSON.stringify(normalizedBody) : undefined
	};

	const res = await fetch(`${PUBLIC_API_BASE}${url}`, requestInit);
	const responseBody = await parseJson(res);

	if (!res.ok) {
		const errorBody = responseBody as ApiErrorBody | null;
		const message =
			typeof errorBody?.message === 'string'
				? `${defaultError}: ${errorBody.message}`
				: defaultError;
		throw new Error(message);
	}

	return normalizeResponse<T>(responseBody);
}
