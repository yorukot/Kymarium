export function normalizeNextPath(nextPath?: string) {
  if (!nextPath) return undefined;

  const trimmed = nextPath.trim();
  if (!trimmed.startsWith("/") || trimmed.startsWith("//")) return undefined;

  return trimmed;
}
