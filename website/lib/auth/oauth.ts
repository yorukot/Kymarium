import { normalizeNextPath } from "@/lib/auth/next-path";

const API_BASE = process.env.NEXT_PUBLIC_API_BASE_URL;

export type OAuthProvider = "google";

export function buildOAuthUrl(provider: OAuthProvider, nextPath?: string) {
  if (!API_BASE) {
    throw new Error("Missing NEXT_PUBLIC_API_BASE_URL");
  }

  const url = new URL(`/api/auth/oauth/${provider}`, API_BASE);
  const safeNext = normalizeNextPath(nextPath) ?? "/";

  url.searchParams.set("next", safeNext);

  return url.toString();
}
