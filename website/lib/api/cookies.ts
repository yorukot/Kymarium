import { headers } from "next/headers"

export async function buildCookieHeader() {
  try {
    const headerStore = await headers()
    return headerStore.get("cookie") ?? ""
  } catch {
    return ""
  }
}
