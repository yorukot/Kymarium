import snakecaseKeys from "snakecase-keys";
import camelcaseKeys from "camelcase-keys";

export type ApiOptions = {
  method?: "GET" | "POST" | "PUT" | "PATCH" | "DELETE";
  headers?: Record<string, string>;
  body?: unknown;

  /** fetch init */
  credentials?: RequestCredentials;

  /** error message */
  defaultError?: string;

  /** abort timeout */
  timeoutMs?: number;

  /** transform */
  snakecaseBody?: boolean; // request: camelCase -> snake_case
  camelcaseResponse?: boolean; // response: snake_case -> camelCase

  /** auth UX */
  redirectOn401?: boolean; // if true, client-side redirect to /login?next=...
  loginPath?: string; // default "/login"
};

export class ApiError extends Error {
  constructor(
    message: string,
    public status: number,
    public body: unknown,
  ) {
    super(message);
    this.name = "ApiError";
  }
}

export type ApiResult<T> = {
  statusCode: number;
  data: T;
};

const API_BASE = process.env.NEXT_PUBLIC_API_BASE_URL;

function isPlainObject(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

function hasStringMessage(x: unknown): x is { message: string } {
  return isPlainObject(x) && typeof x.message === "string";
}

async function parseJson(res: Response): Promise<unknown> {
  const text = await res.text();
  if (!text) return null;
  try {
    return JSON.parse(text);
  } catch {
    return text;
  }
}

function camelizeUnknown(data: unknown): unknown {
  if (Array.isArray(data) || isPlainObject(data)) {
    return camelcaseKeys(data as Record<string, unknown>, { deep: true });
  }
  return data;
}

function redirectToLogin(loginPath: string) {
  if (typeof window === "undefined") return;

  if (window.location.pathname === loginPath) return;

  const currentPath =
    window.location.pathname + window.location.search + window.location.hash;

  const loginUrl = new URL(loginPath, window.location.origin);
  loginUrl.searchParams.set("next", currentPath);

  window.location.replace(loginUrl.toString());
}

export async function apiRequest<T>(
  path: string,
  options: ApiOptions = {},
): Promise<ApiResult<T>> {
  if (!API_BASE) throw new Error("Missing NEXT_PUBLIC_API_BASE_URL");

  const {
    method = "GET",
    headers,
    body,
    credentials = "include",
    defaultError = "Request failed",
    timeoutMs = 15000,
    snakecaseBody = true,
    camelcaseResponse = true,
    redirectOn401 = false,
    loginPath = "/login",
  } = options;

  const shouldHaveBody = body !== undefined && method !== "GET";

  let finalBody: unknown = body;

  // Request normalization: camelCase -> snake_case
  if (snakecaseBody && shouldHaveBody && isPlainObject(body)) {
    finalBody = snakecaseKeys(body, { deep: true });
  }

  const controller = new AbortController();
  const timer = setTimeout(() => controller.abort(), timeoutMs);

  try {
    const res = await fetch(`${API_BASE}${path}`, {
      method,
      credentials,
      signal: controller.signal,
      headers: {
        ...(shouldHaveBody ? { "Content-Type": "application/json" } : {}),
        ...(headers ?? {}),
      },
      body: shouldHaveBody ? JSON.stringify(finalBody) : undefined,
    });

    let parsed = await parseJson(res);

    // Response normalization: snake_case -> camelCase (also for error body)
    if (camelcaseResponse) {
      parsed = camelizeUnknown(parsed);
    }

    // 401 handling
    if (res.status === 401) {
      if (redirectOn401) redirectToLogin(loginPath);
      throw new ApiError("Unauthorized", res.status, parsed);
    }

    // Non-2xx handling
    if (!res.ok) {
      const message = hasStringMessage(parsed)
        ? `${defaultError}: ${parsed.message}`
        : `${defaultError} (${res.status})`;

      throw new ApiError(message, res.status, parsed);
    }

    return {
      statusCode: res.status,
      data: parsed as T,
    };
  } finally {
    clearTimeout(timer);
  }
}
