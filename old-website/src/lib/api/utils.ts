import { PUBLIC_API_BASE } from '$env/static/public';
import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import camelcaseKeys from 'camelcase-keys';
import snakecaseKeys from 'snakecase-keys';
import { refreshToken } from './auth';
import { StatusCodes } from 'http-status-codes';

type HttpMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';

type ApiErrorBody = {
	message?: string;
};

export type ApiDefaultBody = {
	message: string;
	statusCode: number;
};

export type ApiResponse<T> = ApiDefaultBody & {
	data: T;
};

export class ApiError extends Error {
	statusCode: number;
	body: ApiErrorBody | null;

	constructor(message: string, statusCode: number, body: ApiErrorBody | null = null) {
		super(message);
		this.name = 'ApiError';
		this.statusCode = statusCode;
		this.body = body;
		Object.setPrototypeOf(this, new.target.prototype);
	}
}

type ApiOptions = {
	method?: HttpMethod;
	body?: unknown;
	defaultError?: string;
	headers?: HeadersInit;
};

export async function apiRequest<T>(url: string, options: ApiOptions = {}): Promise<T> {
	const { method = 'GET', defaultError = 'Request failed', headers } = options;
	let { body } = options;
	const shouldNormalizeBody = body && (method === 'POST' || method === 'PUT' || method === 'PATCH');

	// Make sure body is in snake_case
	if (shouldNormalizeBody && body && typeof body === 'object') {
		body = snakecaseKeys(body as Record<string, unknown>, { deep: true });
	}

	const requestInit: RequestInit = {
		method,
		credentials: 'include',
		headers: {
			...(body ? { 'Content-Type': 'application/json' } : {}),
			...headers
		},
		body: body ? JSON.stringify(body) : undefined
	};

	const res = await fetch(`${PUBLIC_API_BASE}/${url}`, requestInit);

	let responseBody = await parseJson(res);

	/* ---------- auth retry ---------- */

	if (res.status === StatusCodes.UNAUTHORIZED) {
		const refreshed = await refreshToken();
		if (!refreshed) {
			await redirectToLogin();
			throw new ApiError('AUTH_EXPIRED', res.status);
		}

		const retryRes = await fetch(`${PUBLIC_API_BASE}/${url}`, requestInit);

		responseBody = await parseJson(retryRes);
	}
	/* ---------- error handling ---------- */

	if (!res.ok) {
		const errorBody = normalizeResponse<ApiErrorBody | null>(responseBody);
		const message =
			typeof errorBody?.message === 'string'
				? `${defaultError}: ${errorBody.message}`
				: defaultError;
		throw new ApiError(message, res.status, errorBody);
	}

	/* ---------- success ---------- */

	const normalized = normalizeResponse<T>(responseBody);
	return attachStatusCode<T>(normalized, res.status);
}

export async function publicApiRequest<T>(url: string, options: ApiOptions = {}): Promise<T> {
	const { method = 'GET', defaultError = 'Request failed', headers } = options;
	let { body } = options;
	const shouldNormalizeBody = body && (method === 'POST' || method === 'PUT' || method === 'PATCH');

	// Make sure body is in snake_case
	if (shouldNormalizeBody && body && typeof body === 'object') {
		body = snakecaseKeys(body as Record<string, unknown>, { deep: true });
	}

	const requestInit: RequestInit = {
		method,
		credentials: 'omit',
		headers: {
			...(body ? { 'Content-Type': 'application/json' } : {}),
			...headers
		},
		body: body ? JSON.stringify(body) : undefined
	};

	const res = await fetch(`${PUBLIC_API_BASE}/${url}`, requestInit);
	const responseBody = await parseJson(res);

	if (!res.ok) {
		const errorBody = normalizeResponse<ApiErrorBody | null>(responseBody);
		const message =
			typeof errorBody?.message === 'string'
				? `${defaultError}: ${errorBody.message}`
				: defaultError;
		throw new ApiError(message, res.status, errorBody);
	}

	const normalized = normalizeResponse<T>(responseBody);
	return attachStatusCode<T>(normalized, res.status);
}

async function redirectToLogin() {
	if (!browser) return;

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

function normalizeResponse<T>(data: unknown): T {
	if (data && typeof data === 'object') {
		return camelcaseKeys(data as object, { deep: true }) as T;
	}
	return data as T;
}

function attachStatusCode<T>(data: unknown, statusCode: number): T {
	if (data && typeof data === 'object') {
		const record = data as Record<string, unknown>;
		if (!('statusCode' in record)) {
			return { ...record, statusCode } as T;
		}
	}
	return data as T;
}
